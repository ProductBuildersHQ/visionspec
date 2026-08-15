package pipelines

import (
	"testing"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/integrations"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/spectype"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
)

func TestDefaultsLoad(t *testing.T) {
	want := []string{
		"aws-one-way-door-to-ai-dlc",
		"big-tech-product-to-openspec",
		"pbhq-lite-to-spec-kit",
	}
	got := List()
	if len(got) != len(want) {
		t.Fatalf("expected %d pipelines %v, got %d %v", len(want), want, len(got), got)
	}
	for i, id := range want {
		if got[i] != id {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], id)
		}
	}
}

func TestAllValidate(t *testing.T) {
	for _, p := range All() {
		if err := p.Validate(); err != nil {
			t.Errorf("pipeline %q failed validation: %v", p.ID, err)
		}
	}
}

// TestReferencesResolve is the key integrity check: every pipeline must point at
// a real definition workflow, a real execution integration, a real input port on
// that integration, and its handoffs must name real spec types and real
// execution-side artifacts.
func TestReferencesResolve(t *testing.T) {
	specTypes := map[string]bool{}
	for _, st := range spectype.CoreSpecTypes() {
		specTypes[st.ID] = true
	}

	loader := workflows.DefaultLoader()
	workflowNames := map[string]bool{}
	for _, name := range loader.Available() {
		workflowNames[name] = true
	}

	for _, p := range All() {
		// Definition workflow must exist and load.
		if !workflowNames[p.Definition.Workflow] {
			t.Errorf("pipeline %q references unknown workflow %q", p.ID, p.Definition.Workflow)
		} else if _, err := loader.Load(p.Definition.Workflow); err != nil {
			t.Errorf("pipeline %q workflow %q failed to load: %v", p.ID, p.Definition.Workflow, err)
		}

		// Definition outputs must be real spec types.
		for _, out := range p.Definition.Outputs {
			if !specTypes[out] {
				t.Errorf("pipeline %q definition output %q is not a known spec type", p.ID, out)
			}
		}

		// Execution integration must exist.
		in, err := integrations.Get(p.Execution.Integration)
		if err != nil {
			t.Errorf("pipeline %q references unknown integration %q", p.ID, p.Execution.Integration)
			continue
		}

		// Input port must resolve on that integration.
		if p.Execution.InputPort != "" {
			found := false
			for _, port := range in.Inputs {
				if port.ID == p.Execution.InputPort {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("pipeline %q references unknown input port %q on integration %q", p.ID, p.Execution.InputPort, in.ID)
			}
		}

		// Handoff spec types must be real; handoff targets must resolve to an
		// artifact or input port on the execution integration.
		artifactIDs := map[string]bool{}
		for _, a := range in.Artifacts {
			artifactIDs[a.ID] = true
		}
		portIDs := map[string]bool{}
		for _, port := range in.Inputs {
			portIDs[port.ID] = true
		}
		for _, h := range p.Handoffs {
			if !specTypes[h.SpecType] {
				t.Errorf("pipeline %q handoff references unknown spec type %q", p.ID, h.SpecType)
			}
			if h.ToArtifact != "" && !artifactIDs[h.ToArtifact] && !portIDs[h.ToArtifact] {
				t.Errorf("pipeline %q handoff target %q is not an artifact or input port on integration %q", p.ID, h.ToArtifact, in.ID)
			}
		}
	}
}
