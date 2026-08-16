package workflows

import (
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
)

// syncedRubricSpecs are owned upstream in prism-roadmap (see tools/prism-sync);
// their layering is prism-roadmap's concern, so the classification guard skips
// them here to avoid diverging synced copies.
var syncedRubricSpecs = map[string]bool{
	"bmc": true, "opportunity-spec": true, "ost": true, "assumption-map": true,
	"discovery-snapshot": true, "experience-map": true,
	"v2mom-vision": true, "v2mom-values": true, "v2mom-methods": true,
	"v2mom-obstacles": true, "v2mom-measures": true, "v2mom-alignment": true,
	"v2mom-summary": true,
}

// TestAllLocalRubricsClassified guards RMI-024: every locally-owned rubric
// carries layered evaluation metadata (a class on every category), so consumers
// can separate advisory Leadership-Principle judgment from blocking
// implementation-readiness gates. prism-synced rubrics are exempt (owned
// upstream). INV-3 is asserted universally below.
func TestAllLocalRubricsClassified(t *testing.T) {
	loader := DefaultLoader()
	classified := 0
	for _, name := range loader.Available() {
		w, err := loader.Load(name)
		if err != nil {
			t.Fatalf("load %s: %v", name, err)
		}
		for spec, r := range w.Rubrics {
			if syncedRubricSpecs[spec] {
				continue
			}
			for _, c := range r.Categories {
				if c.Class == "" {
					t.Errorf("%s/%s: category %q has no class (RMI-024 layering)", name, spec, c.ID)
					continue
				}
				classified++
			}
		}
	}
	if classified == 0 {
		t.Fatal("no classified categories found — layering regressed")
	}
}

// TestConditionalBlockingCategoriesAllowNA guards RMI-029: the generic
// lightweight rubrics (pbhq-lite, inherited by quick-fix) get applied to
// arbitrary initiatives, so their conditional blocking gates — the ones that
// assume a service/feature — must offer a not_applicable option, or they
// false-fail on maintenance/refactor/compliance/library work (the RMI-026
// finding). Universally-applicable gates (goals, requirements) intentionally
// have no escape hatch and are not listed here.
func TestConditionalBlockingCategoriesAllowNA(t *testing.T) {
	want := map[string][]string{
		"trd": {"api_design", "data_models", "security_considerations", "dependencies"},
		"prd": {"user_stories", "success_metrics"},
	}
	w, err := DefaultLoader().Load("pbhq-lite")
	if err != nil {
		t.Fatalf("load pbhq-lite: %v", err)
	}
	for spec, cats := range want {
		r, ok := w.Rubrics[spec]
		if !ok {
			t.Fatalf("pbhq-lite: no %s rubric", spec)
		}
		for _, catID := range cats {
			var cat *rubric.Category
			for i := range r.Categories {
				if r.Categories[i].ID == catID {
					cat = &r.Categories[i]
				}
			}
			if cat == nil {
				t.Errorf("pbhq-lite/%s: category %q not found", spec, catID)
				continue
			}
			hasNA := false
			for _, o := range cat.Scale.Options {
				if o.Value == "not_applicable" {
					hasNA = true
				}
			}
			if !hasNA {
				t.Errorf("pbhq-lite/%s: conditional blocking category %q lacks a not_applicable option (RMI-029)", spec, catID)
			}
		}
	}
}

// TestNoLeadershipPrincipleIsBlocking guards INV-3 across the entire embedded
// library (synced rubrics included): advisory principle-based judgment can
// never be a hard implementation gate.
func TestNoLeadershipPrincipleIsBlocking(t *testing.T) {
	loader := DefaultLoader()
	for _, name := range loader.Available() {
		w, err := loader.Load(name)
		if err != nil {
			t.Fatalf("load %s: %v", name, err)
		}
		for spec, r := range w.Rubrics {
			for _, c := range r.Categories {
				if c.Class == rubric.ClassLeadershipPrinciple && c.Blocking {
					t.Errorf("%s/%s: category %q is leadership_principle and blocking (violates INV-3)", name, spec, c.ID)
				}
			}
		}
	}
}
