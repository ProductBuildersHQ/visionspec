package loops

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/integrations"
	"github.com/ProductBuildersHQ/visionspec/pkg/spectype"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

func TestDefaultsLoad(t *testing.T) {
	want := []string{"aws-ai-dlc", "pbhq-two-loop"}
	got := List()
	if len(got) != len(want) {
		t.Fatalf("expected %d loop systems %v, got %d %v", len(want), want, len(got), got)
	}
	for i, id := range want {
		if got[i] != id {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], id)
		}
	}
}

func TestAllValidate(t *testing.T) {
	for _, s := range All() {
		if err := s.Validate(); err != nil {
			t.Errorf("loop system %q failed validation: %v", s.ID, err)
		}
	}
}

func TestTwoLoopShape(t *testing.T) {
	s, err := Get("pbhq-two-loop")
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Loops) != 2 {
		t.Fatalf("pbhq-two-loop should have 2 loops, got %d", len(s.Loops))
	}
	if s.Loops[0].ID != "product" || s.Loops[1].ID != "builder" {
		t.Fatalf("expected product then builder, got %q, %q", s.Loops[0].ID, s.Loops[1].ID)
	}
	if len(s.Seams) != 3 {
		t.Fatalf("expected 3 seams (handoff, telemetry, escalation), got %d", len(s.Seams))
	}
}

// TestCrossReferencesResolve ensures every station's spec type, workflow, and
// integration reference resolves against the rest of the library.
func TestCrossReferencesResolve(t *testing.T) {
	specTypes := map[string]bool{}
	for _, st := range spectype.CoreSpecTypes() {
		specTypes[st.ID] = true
	}
	workflowNames := map[string]bool{}
	for _, name := range workflows.DefaultLoader().Available() {
		workflowNames[name] = true
	}

	for _, s := range All() {
		for _, l := range s.Loops {
			for _, st := range l.Stations {
				for _, id := range st.SpecTypes {
					if !specTypes[id] {
						t.Errorf("%s: station %s.%s references unknown spec type %q", s.ID, l.ID, st.ID, id)
					}
				}
				for _, name := range st.Workflows {
					if !workflowNames[name] {
						t.Errorf("%s: station %s.%s references unknown workflow %q", s.ID, l.ID, st.ID, name)
					}
				}
				for _, id := range st.Integrations {
					if _, err := integrations.Get(id); err != nil {
						t.Errorf("%s: station %s.%s references unknown integration %q", s.ID, l.ID, st.ID, id)
					}
				}
			}
		}
	}
}
