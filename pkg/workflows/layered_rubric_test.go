package workflows

import (
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
)

// hardenedRubrics lists every workflow/spec rubric hardened under
// INIT-SPECWORKFLOWSPEC-001: normative categories with layered class/blocking
// metadata (RMI-008..013). Extend this list as more rubrics are hardened.
var hardenedRubrics = []struct {
	workflow string
	spec     string
}{
	{"enterprise", "trd"},
	{"enterprise", "prd"},
	{"enterprise", "tpd"},
	{"enterprise", "ird"},
	{"enterprise", "uxd"},
	{"aws-two-way-door", "trd"},
	{"aws-two-way-door", "prd"},
	{"aws-two-way-door", "tpd"},
	{"aws-two-way-door", "ird"},
	{"aws-two-way-door", "uxd"},
	{"aws-one-way-door", "trd"},
	{"aws-one-way-door", "prd"},
	{"aws-one-way-door", "tpd"},
	{"aws-one-way-door", "ird"},
	{"aws-one-way-door", "uxd"},
}

// TestHardenedRubricWeightsSumToOne guards TRD-011: category weights in every
// hardened rubric must sum to 1.0, so weighted scoring is well-defined.
func TestHardenedRubricWeightsSumToOne(t *testing.T) {
	loader := DefaultLoader()
	for _, hr := range hardenedRubrics {
		w, err := loader.Load(hr.workflow)
		if err != nil {
			t.Fatalf("load %s: %v", hr.workflow, err)
		}
		r, ok := w.Rubrics[hr.spec]
		if !ok {
			t.Fatalf("%s/%s: rubric not found", hr.workflow, hr.spec)
		}
		var sum float64
		for _, c := range r.Categories {
			sum += c.Weight
		}
		if sum < 0.999 || sum > 1.001 {
			t.Errorf("%s/%s: category weights sum to %v, want 1.0", hr.workflow, hr.spec, sum)
		}
	}
}

// TestHardenedRubricLayering guards TRD-009: every category in a hardened
// rubric is classified, and INV-3 — a leadership_principle category is never
// blocking, since advisory judgment cannot gate implementation.
func TestHardenedRubricLayering(t *testing.T) {
	loader := DefaultLoader()
	for _, hr := range hardenedRubrics {
		w, err := loader.Load(hr.workflow)
		if err != nil {
			t.Fatalf("load %s: %v", hr.workflow, err)
		}
		r, ok := w.Rubrics[hr.spec]
		if !ok {
			t.Fatalf("%s/%s: rubric not found", hr.workflow, hr.spec)
		}

		blockingCount, lpCount := 0, 0
		for _, c := range r.Categories {
			if c.Class == "" {
				t.Errorf("%s/%s: category %q has no class", hr.workflow, hr.spec, c.ID)
			}
			if c.Blocking {
				blockingCount++
			}
			if c.Class == rubric.ClassLeadershipPrinciple {
				lpCount++
				if c.Blocking {
					t.Errorf("%s/%s: category %q is leadership_principle and blocking (violates INV-3)", hr.workflow, hr.spec, c.ID)
				}
			}
		}
		if blockingCount == 0 {
			t.Errorf("%s/%s: no blocking category found — a hardened rubric should have at least one implementation_readiness gate", hr.workflow, hr.spec)
		}

		// Validate() itself must also reject an LP+blocking combination — the
		// invariant is enforced at the schema layer, not only asserted here.
		if issues := r.Validate(); hasBlockingLPIssue(issues) {
			t.Errorf("%s/%s: Validate() reported an LP-blocking issue on content that should be clean: %v", hr.workflow, hr.spec, issues)
		}
	}
}

func hasBlockingLPIssue(issues []string) bool {
	for _, issue := range issues {
		if len(issue) > 0 && containsMustNotBlocking(issue) {
			return true
		}
	}
	return false
}

func containsMustNotBlocking(s string) bool {
	const substr = "must not be blocking"
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
