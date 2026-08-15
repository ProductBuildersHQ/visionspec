package integrations

import (
	"testing"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/spectype"
)

func specTypeIDs() map[string]bool {
	ids := map[string]bool{}
	for _, st := range spectype.CoreSpecTypes() {
		ids[st.ID] = true
	}
	return ids
}

func TestDefaultsLoad(t *testing.T) {
	want := []string{"ai-dlc", "kiro", "openspec", "spec-kit"}
	got := List()
	if len(got) != len(want) {
		t.Fatalf("expected %d integrations %v, got %d %v", len(want), want, len(got), got)
	}
	for i, id := range want {
		if got[i] != id {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], id)
		}
	}
}

func TestAllValidate(t *testing.T) {
	for _, in := range All() {
		if err := in.Validate(); err != nil {
			t.Errorf("integration %q failed validation: %v", in.ID, err)
		}
	}
}

// TestSpecTypeReferencesResolve ensures every spec_type referenced by an
// artifact or input port names a real spec type in the registry.
func TestSpecTypeReferencesResolve(t *testing.T) {
	valid := specTypeIDs()
	for _, in := range All() {
		for _, a := range in.Artifacts {
			if a.SpecType != "" && !valid[a.SpecType] {
				t.Errorf("integration %q artifact %q references unknown spec_type %q", in.ID, a.ID, a.SpecType)
			}
		}
		for _, p := range in.Inputs {
			for _, st := range p.SpecTypes {
				if !valid[st] {
					t.Errorf("integration %q input %q references unknown spec_type %q", in.ID, p.ID, st)
				}
			}
		}
	}
}

func TestGetUnknown(t *testing.T) {
	if _, err := Get("does-not-exist"); err == nil {
		t.Fatal("expected error for unknown integration, got nil")
	}
}
