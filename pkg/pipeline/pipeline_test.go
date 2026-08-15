package pipeline

import "testing"

func validPipeline() *Pipeline {
	return &Pipeline{
		ID:   "aws-one-way-door-to-ai-dlc",
		Name: "AWS Working Backwards to AI-DLC",
		Definition: DefinitionRef{
			Workflow: "aws-one-way-door",
			Outputs:  []string{"narrative-6p", "prd", "trd"},
		},
		Execution: ExecutionRef{
			Integration: "ai-dlc",
			InputPort:   "requirements-in",
		},
		Handoffs: []Handoff{
			{SpecType: "trd", ToArtifact: "requirements", Transform: "seed construction intent"},
		},
	}
}

func TestValidate_OK(t *testing.T) {
	if err := validPipeline().Validate(); err != nil {
		t.Fatalf("expected valid pipeline, got: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Pipeline)
	}{
		{"missing id", func(p *Pipeline) { p.ID = "" }},
		{"missing name", func(p *Pipeline) { p.Name = "" }},
		{"missing workflow", func(p *Pipeline) { p.Definition.Workflow = "" }},
		{"missing integration", func(p *Pipeline) { p.Execution.Integration = "" }},
		{"handoff empty spec_type", func(p *Pipeline) {
			p.Handoffs = append(p.Handoffs, Handoff{ToArtifact: "x"})
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := validPipeline()
			tc.mutate(p)
			if err := p.Validate(); err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
id: pbhq-lite-to-spec-kit
name: PBHQ Lite to Spec-Kit
definition:
  workflow: pbhq-lite
  outputs:
    - spec
execution:
  integration: spec-kit
  input_port: feature-description
handoffs:
  - spec_type: spec
    to_artifact: spec
`)
	p, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if p.Definition.Workflow != "pbhq-lite" {
		t.Fatalf("expected workflow pbhq-lite, got %q", p.Definition.Workflow)
	}
	if p.Execution.Integration != "spec-kit" {
		t.Fatalf("expected integration spec-kit, got %q", p.Execution.Integration)
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("parsed pipeline should validate: %v", err)
	}
}
