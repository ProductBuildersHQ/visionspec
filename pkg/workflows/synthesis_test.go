package workflows

import (
	"slices"
	"sort"
	"strings"
	"testing"
)

// TestSynthesisSourcesDeclared guards the whole embedded workflow library:
// every same-document source named by a synthesis rule must be a spec the
// workflow declares in spec_config. A synthesis rule consuming an undeclared
// spec is drift — a document the workflow can neither render nor evaluate
// feeding one it can. Qualified sources ("parent:v2mom-methods",
// "grandparent:v2mom-vision") address documents in another instance of the
// cascade hierarchy at runtime, not this workflow's spec_config, and are
// out of scope here.
func TestSynthesisSourcesDeclared(t *testing.T) {
	loader := DefaultLoader()
	names := loader.Available()
	sort.Strings(names)

	for _, name := range names {
		w, err := loader.Load(name)
		if err != nil {
			t.Errorf("workflow %q: load error: %v", name, err)
			continue
		}
		for target, rule := range w.Workflow.Synthesis {
			if rule == nil {
				continue
			}
			for _, src := range rule.Sources {
				if strings.Contains(src, ":") {
					continue
				}
				if w.Workflow.SpecConfig[src] == nil {
					t.Errorf("workflow %q: synthesis %q names source %q, which is not declared in spec_config", name, target, src)
				}
			}
		}
	}
}

// TestAWSDoorProfilesSymmetricDeepening pins the symmetry between the two AWS
// door profiles: both document MRD (product-scale) and OpportunitySpec
// (feature-scale) as the optional post-FAQ deepening paths, so whichever one a
// team chooses must feed PRD synthesis in both profiles — the door sets the
// ceremony; the scale picks the tool.
func TestAWSDoorProfilesSymmetricDeepening(t *testing.T) {
	loader := DefaultLoader()
	for _, name := range []string{"aws-one-way-door", "aws-two-way-door"} {
		w, err := loader.Load(name)
		if err != nil {
			t.Fatalf("workflow %q: load error: %v", name, err)
		}
		rule := w.Workflow.Synthesis["prd"]
		if rule == nil {
			t.Fatalf("workflow %q: no prd synthesis rule", name)
		}
		for _, src := range []string{"press", "faq", "mrd", "opportunity-spec"} {
			if !slices.Contains(rule.Sources, src) {
				t.Errorf("workflow %q: prd synthesis sources %v missing %q", name, rule.Sources, src)
			}
		}
	}
}
