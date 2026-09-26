package speckit

import (
	"fmt"
	"path"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

// WorkflowManifest mirrors spec-kit's workflow.yml (plugin schema 1.0).
type WorkflowManifest struct {
	SchemaVersion string                   `yaml:"schema_version"`
	Workflow      WorkflowMeta             `yaml:"workflow"`
	Requires      ManifestRequires         `yaml:"requires"`
	Inputs        map[string]WorkflowInput `yaml:"inputs,omitempty"`
	Steps         []any                    `yaml:"steps"`
}

// WorkflowMeta is the workflow identity block.
type WorkflowMeta struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
}

// WorkflowInput is one entry in a workflow's inputs map.
type WorkflowInput struct {
	Type     string   `yaml:"type"`
	Required bool     `yaml:"required,omitempty"`
	Default  string   `yaml:"default,omitempty"`
	Prompt   string   `yaml:"prompt,omitempty"`
	Enum     []string `yaml:"enum,omitempty"`
}

// WFCommandStep runs one spec-kit command.
type WFCommandStep struct {
	ID          string      `yaml:"id"`
	Command     string      `yaml:"command"`
	Integration string      `yaml:"integration"`
	Input       WFStepInput `yaml:"input"`
}

// WFStepInput is a command step's argument string.
type WFStepInput struct {
	Args string `yaml:"args"`
}

// WFGateStep pauses for human review.
type WFGateStep struct {
	ID      string `yaml:"id"`
	Type    string `yaml:"type"`
	Message string `yaml:"message"`
}

// WFIfStep conditionally runs nested steps.
type WFIfStep struct {
	ID        string `yaml:"id"`
	Type      string `yaml:"type"`
	Condition string `yaml:"condition"`
	Then      []any  `yaml:"then"`
}

// BuildWorkflow exports a loaded workflow family's full execution sequence
// (product specs, then technical specs) as a single spec-kit workflow: one
// command step per spec type, invoking whichever extension
// (BuildExtensions) owns that spec, with review gates placed exactly where
// the family's Execution.ReviewGates declare them. Optional spec types
// (SpecConfig[spec].Required == false) are wrapped in an `if` step gated by
// a generated include_<spec> input, so the chain can be run either minimally
// or with every deepening enabled.
//
// This is the data-driven counterpart to a hand-authored triage workflow
// (e.g. risk-adjusted): it does not branch on work type or reversibility —
// it simply runs one family's sequence end to end, in the node order and
// gates the family already declares.
func BuildWorkflow(lw *workflows.LoadedWorkflow, opts Options) (File, error) {
	if lw == nil || lw.Workflow == nil {
		return File{}, fmt.Errorf("nil workflow")
	}
	w := lw.Workflow
	if w.Execution == nil || len(w.Execution.Sequence) == 0 {
		return File{}, fmt.Errorf("workflow %q has no execution sequence to export", w.Name)
	}

	opts, err := opts.withDefaults(w.Name)
	if err != nil {
		return File{}, err
	}

	product, engineering := partitionSequence(w)

	gates := map[string]*workflow.ReviewGate{}
	if w.Execution != nil {
		for i := range w.Execution.ReviewGates {
			g := &w.Execution.ReviewGates[i]
			gates[g.After] = g
		}
	}

	inputs := map[string]WorkflowInput{
		"work": {
			Type:     "string",
			Required: true,
			Prompt:   "Describe the work this chain produces documents for",
		},
		"slug": {
			Type:     "string",
			Required: true,
			Prompt:   "Short kebab-case slug for this work",
		},
		"integration": {
			Type:    "string",
			Default: "auto",
			Prompt:  "Integration to use (e.g. claude, copilot; 'auto' uses the project's initialized integration)",
		},
	}

	var steps []any
	for i, spec := range w.Execution.Sequence {
		extID := opts.ProductID
		if slices.Contains(engineering, spec) {
			extID = opts.EngineeringID
		}

		args := "slug={{ inputs.slug }}"
		if i == 0 {
			args = "{{ inputs.work }} slug={{ inputs.slug }}"
		}
		cmdStep := WFCommandStep{
			ID:          spec,
			Command:     commandName(extID, spec),
			Integration: "{{ inputs.integration }}",
			Input:       WFStepInput{Args: args},
		}

		if isRequiredSpec(w, spec) {
			steps = append(steps, cmdStep)
		} else {
			includeInput := "include_" + spec
			cmdStep.ID = spec + "-run"
			inputs[includeInput] = WorkflowInput{
				Type:    "string",
				Default: "skip",
				Enum:    []string{"skip", "include"},
				Prompt:  fmt.Sprintf("Optional deepening: run %s? skip | include", specTitle(spec)),
			}
			steps = append(steps, WFIfStep{
				ID:        spec,
				Type:      "if",
				Condition: fmt.Sprintf("{{ inputs.%s == 'include' }}", includeInput),
				Then:      []any{cmdStep},
			})
		}

		if g := gates[spec]; g != nil {
			message := fmt.Sprintf("%s review gate.", actionTitle(g.Action))
			if g.Required {
				message = fmt.Sprintf("%s (required): do not proceed without an explicit human decision.", actionTitle(g.Action))
			}
			steps = append(steps, WFGateStep{ID: g.StepID(), Type: "gate", Message: message})
		}
	}

	manifest := WorkflowManifest{
		SchemaVersion: "1.0",
		Workflow: WorkflowMeta{
			ID:          opts.WorkflowID,
			Name:        specTitle(opts.ProductID) + " Chain",
			Version:     opts.Version,
			Author:      opts.Author,
			Description: chainDescription(w, w.Execution.Sequence, "end-to-end chain, generated from the family's execution sequence and review gates"),
		},
		Requires: ManifestRequires{SpecKitVersion: opts.SpecKitVersion},
		Inputs:   inputs,
		Steps:    steps,
	}

	if !extensionIDRe.MatchString(manifest.Workflow.ID) {
		return File{}, fmt.Errorf("workflow id %q: must match %s", manifest.Workflow.ID, extensionIDRe)
	}
	if len(product) == 0 && len(engineering) == 0 {
		return File{}, fmt.Errorf("workflow %q has no specs to chain", w.Name)
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return File{}, fmt.Errorf("marshaling workflow.yml for %s: %w", opts.WorkflowID, err)
	}
	content := append([]byte(yamlHeader(opts.Family)), data...)

	return File{
		Path:    path.Join("workflows", opts.WorkflowID, "workflow.yml"),
		Content: content,
	}, nil
}
