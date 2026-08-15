package workflows

import (
	"slices"
	"testing"
)

// TestPBHQLiteExecutionParsed guards against the silent-drop regression where
// pbhq-lite's execution block used the unrecognized top-level key "workflow:"
// and parsed to nothing.
func TestPBHQLiteExecutionParsed(t *testing.T) {
	w, err := DefaultLoader().Load("pbhq-lite")
	if err != nil {
		t.Fatalf("Load(pbhq-lite) error: %v", err)
	}

	exec := w.Workflow.Execution
	if exec == nil {
		t.Fatal("pbhq-lite Execution is nil — execution block failed to parse")
	}
	wantSeq := []string{"prd", "trd", "plan", "roadmap"}
	if !slices.Equal(exec.Sequence, wantSeq) {
		t.Errorf("pbhq-lite sequence = %v, want %v", exec.Sequence, wantSeq)
	}
	if exec.IterationTrigger != "roadmap" {
		t.Errorf("pbhq-lite iteration_trigger = %q, want %q", exec.IterationTrigger, "roadmap")
	}
	if len(exec.ReviewGates) != 2 {
		t.Errorf("pbhq-lite review gates = %d, want 2", len(exec.ReviewGates))
	}
}

// TestAWSProductMatchesD2 verifies aws-one-way-door against the authoritative
// visionspec aws-one-way-door-flow.d2 diagram: execution phases, sequence, and
// the synthesis DAG edges.
func TestAWSProductMatchesD2(t *testing.T) {
	w, err := DefaultLoader().Load("aws-one-way-door")
	if err != nil {
		t.Fatalf("Load(aws-one-way-door) error: %v", err)
	}

	exec := w.Workflow.Execution
	if exec == nil {
		t.Fatal("aws-one-way-door Execution is nil")
	}
	wantSeq := []string{"press", "faq", "mrd", "prd", "narrative-6p", "uxd", "trd", "tpd", "ird"}
	if !slices.Equal(exec.Sequence, wantSeq) {
		t.Errorf("aws-one-way-door sequence = %v, want %v", exec.Sequence, wantSeq)
	}
	if len(exec.Phases) != 6 {
		t.Fatalf("aws-one-way-door phases = %d, want 6", len(exec.Phases))
	}
	wantPhases := map[string][]string{
		"vision":       {"press", "faq"},
		"validation":   {"mrd"},
		"requirements": {"prd"},
		"decision":     {"narrative-6p"},
		"experience":   {"uxd"},
		"technical":    {"trd", "tpd", "ird"},
	}
	for _, p := range exec.Phases {
		want, ok := wantPhases[p.ID]
		if !ok {
			t.Errorf("aws-one-way-door unexpected phase %q", p.ID)
			continue
		}
		if !slices.Equal(p.Specs, want) {
			t.Errorf("aws-one-way-door phase %q specs = %v, want %v", p.ID, p.Specs, want)
		}
	}

	// The press release is the human-authored founding artifact — its merged
	// synthesis rule must have empty sources (the explicit override of
	// enterprise's press←mrd; a PR synthesized from research is "working
	// forwards").
	if rule, ok := w.Workflow.Synthesis["press"]; !ok {
		t.Error("aws-one-way-door press should carry the explicit human-authored override rule")
	} else if len(rule.Sources) != 0 {
		t.Errorf("aws-one-way-door press should be human-authored (empty sources), got %v", rule.Sources)
	}

	// Synthesis edges per the D2 diagram.
	wantSynthesis := map[string][]string{
		"faq":          {"press"},
		"mrd":          {"press", "faq"},
		"prd":          {"press", "faq", "mrd", "opportunity-spec"},
		"narrative-6p": {"press", "faq", "prd"},
		"trd":          {"prd", "uxd", "mrd"},
		"tpd":          {"prd", "trd", "uxd"},
		"ird":          {"trd"},
	}
	for spec, wantSources := range wantSynthesis {
		rule, ok := w.Workflow.Synthesis[spec]
		if !ok {
			t.Errorf("aws-one-way-door missing synthesis rule for %q", spec)
			continue
		}
		if !slices.Equal(rule.Sources, wantSources) {
			t.Errorf("aws-one-way-door synthesis %q sources = %v, want %v", spec, rule.Sources, wantSources)
		}
	}

	// MRD is optional post-FAQ validation (Amazon has no MRD); the unified
	// output spec is not part of the Working Backwards flow.
	if req, ok := w.Workflow.SpecConfig["mrd"]; !ok || req.Required {
		t.Error("aws-one-way-door spec_config[mrd] should be present and not required")
	}
	if req, ok := w.Workflow.SpecConfig["spec"]; !ok || req.Required {
		t.Error("aws-one-way-door spec_config[spec] should be present and not required")
	}
	// Scale picks the deepening tool: both MRD and OpportunitySpec available.
	if req, ok := w.Workflow.SpecConfig["opportunity-spec"]; !ok || req.Required {
		t.Error("aws-one-way-door opportunity-spec should be present and not required")
	}

	// Iteration is the essence of PR/FAQ practice; the decision meeting is
	// the one-way-door gate.
	if exec.IterationTrigger != "faq" {
		t.Errorf("aws-one-way-door iteration_trigger = %q, want faq", exec.IterationTrigger)
	}
	if len(exec.ReviewGates) != 2 {
		t.Fatalf("aws-one-way-door review gates = %d, want 2", len(exec.ReviewGates))
	}
	if g := exec.ReviewGates[0]; g.After != "faq" || g.Action != "prfaq_review" {
		t.Errorf("aws-one-way-door gate[0] = %+v, want after=faq action=prfaq_review", g)
	}
	if g := exec.ReviewGates[1]; g.After != "narrative-6p" || g.Action != "decision_meeting" || !g.Required {
		t.Errorf("aws-one-way-door gate[1] = %+v, want after=narrative-6p action=decision_meeting required", g)
	}
}

// TestAWSFeatureMatchesD2 verifies aws-two-way-door against the authoritative
// visionspec aws-two-way-door-flow.d2 diagram.
func TestAWSFeatureMatchesD2(t *testing.T) {
	w, err := DefaultLoader().Load("aws-two-way-door")
	if err != nil {
		t.Fatalf("Load(aws-two-way-door) error: %v", err)
	}

	exec := w.Workflow.Execution
	if exec == nil {
		t.Fatal("aws-two-way-door Execution is nil")
	}
	wantSeq := []string{"press", "faq", "opportunity-spec", "prd", "uxd", "trd", "tpd"}
	if !slices.Equal(exec.Sequence, wantSeq) {
		t.Errorf("aws-two-way-door sequence = %v, want %v", exec.Sequence, wantSeq)
	}

	// Same founding artifact as products: human-authored press release
	// (explicit empty-sources override of enterprise's press←mrd).
	if rule, ok := w.Workflow.Synthesis["press"]; !ok {
		t.Error("aws-two-way-door press should carry the explicit human-authored override rule")
	} else if len(rule.Sources) != 0 {
		t.Errorf("aws-two-way-door press should be human-authored (empty sources), got %v", rule.Sources)
	}

	// Two-way-door ceremony: 6-pager optional; OpportunitySpec is optional
	// post-FAQ discovery deepening; unified output spec not in the flow.
	if req, ok := w.Workflow.SpecConfig["narrative-6p"]; !ok || req.Required {
		t.Error("aws-two-way-door narrative-6p should be present and not required")
	}
	if req, ok := w.Workflow.SpecConfig["opportunity-spec"]; !ok || req.Required {
		t.Error("aws-two-way-door opportunity-spec should be present and not required")
	}
	if req, ok := w.Workflow.SpecConfig["spec"]; !ok || req.Required {
		t.Error("aws-two-way-door spec_config[spec] should be present and not required")
	}
	// Scale picks the deepening tool: both MRD and OpportunitySpec available.
	if req, ok := w.Workflow.SpecConfig["mrd"]; !ok || req.Required {
		t.Error("aws-two-way-door mrd should be present and not required")
	}

	wantSynthesis := map[string][]string{
		"faq":              {"press"},
		"opportunity-spec": {"press", "faq"},
		"prd":              {"press", "faq", "opportunity-spec", "mrd"},
		"trd":              {"prd", "uxd"}, // inherited from enterprise
		"tpd":              {"prd", "trd", "uxd"},
	}
	for spec, wantSources := range wantSynthesis {
		rule, ok := w.Workflow.Synthesis[spec]
		if !ok {
			t.Errorf("aws-two-way-door missing synthesis rule for %q", spec)
			continue
		}
		if !slices.Equal(rule.Sources, wantSources) {
			t.Errorf("aws-two-way-door synthesis %q sources = %v, want %v", spec, rule.Sources, wantSources)
		}
	}

	// Two-way-door ceremony: iterative PR/FAQ review is the only formal gate.
	if exec.IterationTrigger != "faq" {
		t.Errorf("aws-two-way-door iteration_trigger = %q, want faq", exec.IterationTrigger)
	}
	if len(exec.ReviewGates) != 1 {
		t.Fatalf("aws-two-way-door review gates = %d, want 1", len(exec.ReviewGates))
	}
	if g := exec.ReviewGates[0]; g.After != "faq" || g.Action != "prfaq_review" {
		t.Errorf("aws-two-way-door gate[0] = %+v, want after=faq action=prfaq_review", g)
	}
}

// TestQuickFixProfile verifies the quick-fix workflow: only ROADMAP required,
// everything else optional, inheriting pbhq-lite's templates and rubrics.
func TestQuickFixProfile(t *testing.T) {
	w, err := DefaultLoader().Load("quick-fix")
	if err != nil {
		t.Fatalf("Load(quick-fix) error: %v", err)
	}

	required := w.Workflow.RequiredSpecs()
	if !slices.Equal(required, []string{"roadmap"}) {
		t.Errorf("quick-fix required specs = %v, want [roadmap]", required)
	}
	for _, opt := range []string{"prd", "trd", "plan"} {
		req, ok := w.Workflow.SpecConfig[opt]
		if !ok {
			t.Errorf("quick-fix missing optional spec %q", opt)
			continue
		}
		if req.Required {
			t.Errorf("quick-fix %q should be optional", opt)
		}
	}

	// Inherits pbhq-lite's roadmap rubric via extends.
	if _, ok := w.Rubrics["roadmap"]; !ok {
		t.Error("quick-fix missing inherited roadmap rubric")
	}
}

// TestAWSProductTemplateInheritance verifies that aws-one-way-door's sparse
// templates dir (press/faq/narrative-6p only) is completed by enterprise's
// templates through the extends chain.
func TestAWSProductTemplateInheritance(t *testing.T) {
	w, err := DefaultLoader().Load("aws-one-way-door")
	if err != nil {
		t.Fatalf("Load(aws-one-way-door) error: %v", err)
	}

	for _, specType := range []string{"press", "faq", "narrative-6p", "mrd", "prd", "trd", "tpd", "uxd"} {
		tmpl, ok := w.Templates[specType]
		if !ok {
			t.Errorf("aws-one-way-door missing template %q (inheritance from enterprise?)", specType)
			continue
		}
		if tmpl.Content == "" {
			t.Errorf("aws-one-way-door template %q has no content", specType)
		}
	}
}

// TestLeadershipPrincipleGrounding verifies both door profiles carry the full
// 16 Amazon Leadership Principles with the canonical source reference, and
// that the LP-derived rubric categories exist.
func TestLeadershipPrincipleGrounding(t *testing.T) {
	for _, id := range []string{"aws-one-way-door", "aws-two-way-door"} {
		w, err := DefaultLoader().Load(id)
		if err != nil {
			t.Fatalf("Load(%s) error: %v", id, err)
		}
		m := w.Workflow.Methodology
		if m == nil {
			t.Fatalf("%s: methodology missing", id)
		}
		if len(m.Principles) != 16 {
			t.Errorf("%s: principles = %d, want 16 (full Amazon LP set)", id, len(m.Principles))
		}
		if m.Reference == "" {
			t.Errorf("%s: methodology reference URL missing", id)
		}

		wantCategories := map[string]string{
			"press":        "economic_sustainability",
			"faq":          "disconfirmation_rigor",
			"narrative-6p": "ownership_horizon",
		}
		for specType, catID := range wantCategories {
			rs, ok := w.Rubrics[specType]
			if !ok {
				t.Errorf("%s: missing %s rubric", id, specType)
				continue
			}
			found := false
			for _, c := range rs.Categories {
				if c.ID == catID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s: %s rubric missing LP category %q", id, specType, catID)
			}
		}
	}
}
