package rubrics

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
)

func TestGet(t *testing.T) {
	tests := []struct {
		specType types.SpecType
		wantErr  bool
		wantName string
	}{
		{types.SpecTypeMRD, false, "MRD Evaluation"},
		{types.SpecTypePRD, false, "PRD Evaluation"},
		{types.SpecTypeUXD, false, "UXD Evaluation"},
		{types.SpecTypePress, false, "Press Release Evaluation"},
		{types.SpecTypeFAQ, false, "FAQ Evaluation"},
		{types.SpecTypeNarrative1P, false, "Narrative 1-Pager Evaluation"},
		{types.SpecTypeNarrative6P, false, "Narrative 6-Pager Evaluation"},
		{types.SpecTypeTRD, false, "TRD Evaluation"},
		{types.SpecTypeIRD, false, "IRD Evaluation"},
		{"unknown", true, ""}, // No rubric for unknown type
	}

	for _, tt := range tests {
		t.Run(string(tt.specType), func(t *testing.T) {
			rs, err := Get(tt.specType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Get(%s) expected error, got nil", tt.specType)
				}
				return
			}
			if err != nil {
				t.Errorf("Get(%s) unexpected error: %v", tt.specType, err)
				return
			}
			if rs.Name != tt.wantName {
				t.Errorf("Get(%s).Name = %q, want %q", tt.specType, rs.Name, tt.wantName)
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	// Should not panic for valid types
	rs := MustGet(types.SpecTypeMRD)
	if rs == nil {
		t.Error("MustGet(MRD) returned nil")
	}

	// Should panic for invalid types
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet(unknown) did not panic")
		}
	}()
	MustGet("unknown")
}

func TestAvailable(t *testing.T) {
	available := Available()
	if len(available) < 3 {
		t.Errorf("Available() returned %d types, expected at least 3", len(available))
	}
}

func TestHasRubric(t *testing.T) {
	hasRubricTypes := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypeNarrative1P,
		types.SpecTypeNarrative6P,
		types.SpecTypeTRD,
		types.SpecTypeIRD,
	}
	for _, st := range hasRubricTypes {
		if !HasRubric(st) {
			t.Errorf("HasRubric(%s) returned false", st)
		}
	}
	if HasRubric("unknown") {
		t.Error("HasRubric(unknown) returned true")
	}
}

func TestRubricSetTotalWeight(t *testing.T) {
	specTypes := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypeNarrative1P,
		types.SpecTypeNarrative6P,
		types.SpecTypeTRD,
		types.SpecTypeIRD,
	}

	for _, st := range specTypes {
		rs := MustGet(st)
		var total float64
		for _, cat := range rs.Categories {
			total += cat.Weight
		}

		// Total weight should be 1.0 (100%)
		if total < 0.99 || total > 1.01 {
			t.Errorf("%s total weight = %f, expected ~1.0", st, total)
		}
	}
}

func TestRubricSetGetCategory(t *testing.T) {
	rs := MustGet(types.SpecTypePRD)

	cat := rs.GetCategory("user_stories")
	if cat == nil {
		t.Fatal("GetCategory('user_stories') not found")
	}
	if cat.Name != "User Stories" {
		t.Errorf("Category name = %q, want 'User Stories'", cat.Name)
	}
	if cat.Weight != 0.20 {
		t.Errorf("Category weight = %f, want 0.20", cat.Weight)
	}

	if rs.GetCategory("nonexistent") != nil {
		t.Error("GetCategory('nonexistent') should return nil")
	}
}

// passOptionCriteria returns the criteria strings for the "pass" scale option.
func passOptionCriteria(cat *rubric.Category) []string {
	for _, opt := range cat.Scale.Options {
		if opt.Value == "pass" {
			return opt.Criteria
		}
	}
	return nil
}

func TestCategoryCriteria(t *testing.T) {
	rs := MustGet(types.SpecTypeMRD)

	for i := range rs.Categories {
		cat := &rs.Categories[i]
		if len(cat.Scale.Options) == 0 {
			t.Errorf("Category %q has no scale options", cat.ID)
			continue
		}
		if len(passOptionCriteria(cat)) == 0 {
			t.Errorf("Category %q has no pass criteria", cat.ID)
		}
	}
}

func TestPassCriteria(t *testing.T) {
	defaultCriteria := DefaultPassCriteria()
	if defaultCriteria.MinCategoriesPassing != "all_required" {
		t.Errorf("DefaultPassCriteria().MinCategoriesPassing = %q, want all_required", defaultCriteria.MinCategoriesPassing)
	}
	if defaultCriteria.MaxFindings == nil {
		t.Fatal("DefaultPassCriteria().MaxFindings is nil")
	}
	if defaultCriteria.MaxFindings.Critical != 0 {
		t.Errorf("DefaultPassCriteria() critical = %d, want 0", defaultCriteria.MaxFindings.Critical)
	}
	if defaultCriteria.MaxFindings.High != 0 {
		t.Errorf("DefaultPassCriteria() high = %d, want 0", defaultCriteria.MaxFindings.High)
	}
	if defaultCriteria.MaxFindings.Medium != -1 {
		t.Errorf("DefaultPassCriteria() medium = %d, want -1 (unlimited)", defaultCriteria.MaxFindings.Medium)
	}

	strictCriteria := StrictPassCriteria()
	if strictCriteria.MinCategoriesPassing != "all" {
		t.Errorf("StrictPassCriteria().MinCategoriesPassing = %q, want all", strictCriteria.MinCategoriesPassing)
	}
	if strictCriteria.MaxFindings == nil || strictCriteria.MaxFindings.Medium != 3 {
		t.Errorf("StrictPassCriteria() medium = %v, want 3", strictCriteria.MaxFindings)
	}
}
