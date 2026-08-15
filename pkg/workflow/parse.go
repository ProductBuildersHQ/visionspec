package workflow

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParseYAML parses a Workflow from YAML bytes.
func ParseYAML(data []byte) (*Workflow, error) {
	var w Workflow
	if err := yaml.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("parsing workflow YAML: %w", err)
	}
	return &w, nil
}

// ParseYAMLFile parses a Workflow from a YAML file path.
func ParseYAMLFile(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading workflow file: %w", err)
	}
	return ParseYAML(data)
}

// ParseYAMLFromFS parses a Workflow from an fs.FS at the given path.
func ParseYAMLFromFS(fsys fs.FS, path string) (*Workflow, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("reading workflow from FS: %w", err)
	}
	return ParseYAML(data)
}

// LoadFromFS loads a workflow directory from an fs.FS.
// Expects profile.yaml at the root of the directory.
func LoadFromFS(fsys fs.FS, dir string) (*Workflow, error) {
	workflowPath := filepath.Join(dir, "profile.yaml")
	return ParseYAMLFromFS(fsys, workflowPath)
}

// ToYAML serializes a Workflow to YAML bytes.
func (w *Workflow) ToYAML() ([]byte, error) {
	data, err := yaml.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("marshaling workflow: %w", err)
	}
	return data, nil
}

// Clone creates a deep copy of the Workflow.
func (w *Workflow) Clone() *Workflow {
	if w == nil {
		return nil
	}

	clone := &Workflow{
		Name:        w.Name,
		Description: w.Description,
		Extends:     w.Extends,
		Abstract:    w.Abstract,
	}

	// Clone Methodology
	if w.Methodology != nil {
		clone.Methodology = &Methodology{
			Name:        w.Methodology.Name,
			Description: w.Methodology.Description,
			Creator:     w.Methodology.Creator,
			Source:      w.Methodology.Source,
			Reference:   w.Methodology.Reference,
		}
		if w.Methodology.Principles != nil {
			clone.Methodology.Principles = make([]Principle, len(w.Methodology.Principles))
			copy(clone.Methodology.Principles, w.Methodology.Principles)
		}
		if w.Methodology.Artifacts != nil {
			clone.Methodology.Artifacts = make(Artifacts, len(w.Methodology.Artifacts))
			copy(clone.Methodology.Artifacts, w.Methodology.Artifacts)
		}
	}

	// Clone SpecConfig
	if w.SpecConfig != nil {
		clone.SpecConfig = make(map[string]*SpecRequirement)
		for k, v := range w.SpecConfig {
			clone.SpecConfig[k] = &SpecRequirement{
				Required:    v.Required,
				Category:    v.Category,
				Description: v.Description,
				Template:    v.Template.clone(),
				Rubric:      v.Rubric.clone(),
			}
		}
	}

	// Clone Synthesis
	if w.Synthesis != nil {
		clone.Synthesis = make(map[string]*SynthesisRule)
		for k, v := range w.Synthesis {
			rule := &SynthesisRule{
				Guidance:      v.Guidance,
				PromptContext: v.PromptContext,
				Required:      v.Required,
				Priority:      v.Priority,
			}
			if v.Sources != nil {
				rule.Sources = make([]string, len(v.Sources))
				copy(rule.Sources, v.Sources)
			}
			clone.Synthesis[k] = rule
		}
	}

	// Clone Execution
	if w.Execution != nil {
		clone.Execution = &Execution{
			IterationTrigger: w.Execution.IterationTrigger,
		}
		if w.Execution.Sequence != nil {
			clone.Execution.Sequence = make([]string, len(w.Execution.Sequence))
			copy(clone.Execution.Sequence, w.Execution.Sequence)
		}
		if w.Execution.Phases != nil {
			clone.Execution.Phases = make([]Phase, len(w.Execution.Phases))
			copy(clone.Execution.Phases, w.Execution.Phases)
		}
		if w.Execution.ReviewGates != nil {
			clone.Execution.ReviewGates = make([]ReviewGate, len(w.Execution.ReviewGates))
			copy(clone.Execution.ReviewGates, w.Execution.ReviewGates)
		}
	}

	// Clone Evaluation
	if w.Evaluation != nil {
		clone.Evaluation = &EvaluationConfig{
			PassThreshold:    w.Evaluation.PassThreshold,
			PartialThreshold: w.Evaluation.PartialThreshold,
		}
		if w.Evaluation.MaxFindingsSeverity != nil {
			clone.Evaluation.MaxFindingsSeverity = &FindingSeverityLimits{
				Critical: w.Evaluation.MaxFindingsSeverity.Critical,
				High:     w.Evaluation.MaxFindingsSeverity.High,
				Medium:   w.Evaluation.MaxFindingsSeverity.Medium,
				Low:      w.Evaluation.MaxFindingsSeverity.Low,
			}
		}
	}

	return clone
}

// Merge combines this workflow with a parent workflow.
// Settings from this workflow override the parent.
func (w *Workflow) Merge(parent *Workflow) *Workflow {
	if parent == nil {
		return w.Clone()
	}

	// Start with a clone of the parent
	merged := parent.Clone()
	merged.Name = w.Name
	merged.Description = w.Description
	merged.Extends = "" // Clear extends since we've resolved it

	// Override with child's methodology if present
	if w.Methodology != nil {
		merged.Methodology = w.Methodology
	}

	// Merge SpecConfig (child overrides parent)
	if w.SpecConfig != nil {
		if merged.SpecConfig == nil {
			merged.SpecConfig = make(map[string]*SpecRequirement)
		}
		maps.Copy(merged.SpecConfig, w.SpecConfig)
	}

	// Merge Synthesis (child overrides parent)
	if w.Synthesis != nil {
		if merged.Synthesis == nil {
			merged.Synthesis = make(map[string]*SynthesisRule)
		}
		maps.Copy(merged.Synthesis, w.Synthesis)
	}

	// Override execution if child has one
	if w.Execution != nil {
		merged.Execution = w.Execution
	}

	// Override evaluation if child has one
	if w.Evaluation != nil {
		merged.Evaluation = w.Evaluation
	}

	return merged
}

// Validate checks if the workflow is valid.
func (w *Workflow) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	return nil
}

// RequiredSpecs returns the list of required spec type IDs.
func (w *Workflow) RequiredSpecs() []string {
	if w.SpecConfig == nil {
		return nil
	}

	var required []string
	for id, req := range w.SpecConfig {
		if req.Required {
			required = append(required, id)
		}
	}
	return required
}

// IsRequired returns whether a spec type is required.
func (w *Workflow) IsRequired(specType string) bool {
	if w.SpecConfig == nil {
		return false
	}
	if req, ok := w.SpecConfig[specType]; ok {
		return req.Required
	}
	return false
}

// GetCategory returns the category for a spec type.
func (w *Workflow) GetCategory(specType string) string {
	if w.SpecConfig == nil {
		return ""
	}
	if req, ok := w.SpecConfig[specType]; ok {
		return req.Category
	}
	return ""
}
