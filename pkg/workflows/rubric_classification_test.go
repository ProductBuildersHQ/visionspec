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
