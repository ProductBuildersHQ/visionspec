// Package pipeline links a definition-side workflow to an execution-side
// integration: the output specs of the workflow become the input to the tool.
//
// A Pipeline is the modeled form of a handoff such as "AWS Working Backwards
// (aws-one-way-door) feeds AWS AI-DLC Workflows (ai-dlc)", where the six-pager,
// PRD, and TRD produced on the definition side become the input the execution
// tool consumes. Modeling the handoff as data (rather than prose) lets a viewer
// render the full definition-to-execution flow, and lets execution engines drive
// the export deterministically.
//
// Pipeline holds no execution logic. It references a workflow by name (see
// pkg/workflow / pkg/workflows) and an integration by ID (see pkg/integration /
// pkg/integrations), and declares how spec types map across the seam.
package pipeline

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Pipeline connects a definition workflow to an execution integration.
type Pipeline struct {
	// ID is the canonical identifier (e.g., "aws-one-way-door-to-ai-dlc").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Canonical identifier for the pipeline"`

	// Name is the human-readable name.
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Human-readable name"`

	// Description explains the pipeline and when to use it.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Pipeline purpose and use case"`

	// Definition is the upstream definition-side workflow.
	Definition DefinitionRef `json:"definition" yaml:"definition" jsonschema:"required,description=Upstream definition workflow"`

	// Execution is the downstream execution-side integration.
	Execution ExecutionRef `json:"execution" yaml:"execution" jsonschema:"required,description=Downstream execution integration"`

	// Handoffs map individual definition-side spec types onto execution-side
	// artifacts or input ports.
	Handoffs []Handoff `json:"handoffs,omitempty" yaml:"handoffs,omitempty" jsonschema:"description=Per-spec-type handoff mappings"`

	// Reference is a URL to documentation describing this pipeline.
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty" jsonschema:"format=uri,description=Documentation URL"`
}

// DefinitionRef identifies the upstream workflow and the specs it hands off.
type DefinitionRef struct {
	// Workflow is the workflow name in specification-workflow-spec
	// (e.g., "aws-one-way-door").
	Workflow string `json:"workflow" yaml:"workflow" jsonschema:"required,description=Definition workflow name"`

	// Outputs lists the spec type IDs handed off downstream
	// (e.g., ["narrative-6p", "prd", "trd"]).
	Outputs []string `json:"outputs,omitempty" yaml:"outputs,omitempty" jsonschema:"description=Spec type IDs handed off downstream"`
}

// ExecutionRef identifies the downstream integration and the port that receives
// the handoff.
type ExecutionRef struct {
	// Integration is the integration ID (e.g., "ai-dlc").
	Integration string `json:"integration" yaml:"integration" jsonschema:"required,description=Execution integration ID"`

	// InputPort is the ID of the integration input Port that receives the
	// handoff, if a specific one applies.
	InputPort string `json:"input_port,omitempty" yaml:"input_port,omitempty" jsonschema:"description=Integration input port ID that receives the handoff"`
}

// Handoff maps one definition-side spec type onto an execution-side destination.
type Handoff struct {
	// SpecType is the definition-side spec type ID (e.g., "trd").
	SpecType string `json:"spec_type" yaml:"spec_type" jsonschema:"required,description=Definition-side spec type ID"`

	// ToArtifact is the execution-side artifact ID (or input port ID) the spec
	// feeds into.
	ToArtifact string `json:"to_artifact,omitempty" yaml:"to_artifact,omitempty" jsonschema:"description=Execution-side artifact or port the spec feeds into"`

	// Transform notes any mapping or transformation applied at the seam
	// (e.g., "flatten to a single requirements.md").
	Transform string `json:"transform,omitempty" yaml:"transform,omitempty" jsonschema:"description=Note on any transform applied at the seam"`
}

// ParseYAML parses a Pipeline from YAML bytes.
func ParseYAML(data []byte) (*Pipeline, error) {
	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("unmarshaling pipeline: %w", err)
	}
	return &p, nil
}

// Validate checks that required fields are present.
func (p *Pipeline) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("pipeline: id is required")
	}
	if p.Name == "" {
		return fmt.Errorf("pipeline %q: name is required", p.ID)
	}
	if p.Definition.Workflow == "" {
		return fmt.Errorf("pipeline %q: definition.workflow is required", p.ID)
	}
	if p.Execution.Integration == "" {
		return fmt.Errorf("pipeline %q: execution.integration is required", p.ID)
	}
	for i, h := range p.Handoffs {
		if h.SpecType == "" {
			return fmt.Errorf("pipeline %q: handoff %d has empty spec_type", p.ID, i)
		}
	}
	return nil
}
