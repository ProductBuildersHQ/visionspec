package synth

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
)

func TestResolveRule(t *testing.T) {
	w := &workflow.Workflow{
		Synthesis: map[string]*workflow.SynthesisRule{
			"faq":   {Sources: []string{"press"}, Guidance: "challenge it"},
			"press": {Sources: []string{}}, // deliberately human-authored
		},
	}

	tests := []struct {
		name        string
		specType    types.SpecType
		wantOK      bool
		wantSources int
	}{
		{"declared rule with sources", types.SpecTypeFAQ, true, 1},
		{"declared empty-sources rule", types.SpecTypePress, true, 0},
		{"undeclared type", types.SpecTypePRD, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := ResolveRule(w, tt.specType)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && len(rule.Sources) != tt.wantSources {
				t.Errorf("len(Sources) = %d, want %d", len(rule.Sources), tt.wantSources)
			}
		})
	}
}

func TestResolveRule_NilWorkflow(t *testing.T) {
	rule, ok := ResolveRule(nil, types.SpecTypeFAQ)
	if ok || rule != nil {
		t.Errorf("expected ok=false, rule=nil for nil workflow, got ok=%v rule=%v", ok, rule)
	}
}

func TestCanSynthesizeWithRule(t *testing.T) {
	tests := []struct {
		name     string
		rule     *workflow.SynthesisRule
		ok       bool
		specType types.SpecType
		want     bool
	}{
		{"rule with sources", &workflow.SynthesisRule{Sources: []string{"press"}}, true, types.SpecTypeFAQ, true},
		{"rule with empty sources is not synthesizable", &workflow.SynthesisRule{Sources: []string{}}, true, types.SpecTypePress, false},
		{"no rule falls back to legacy: bmc not in legacy list", nil, false, types.SpecType("bmc"), false},
		{"no rule falls back to legacy: faq is in legacy list", nil, false, types.SpecTypeFAQ, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanSynthesizeWithRule(tt.rule, tt.ok, tt.specType)
			if got != tt.want {
				t.Errorf("CanSynthesizeWithRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSourcesForRule(t *testing.T) {
	rule := &workflow.SynthesisRule{Sources: []string{"press", "faq"}}
	got := SourcesForRule(rule, true, types.SpecTypePRD)
	want := []types.SpecType{types.SpecTypePress, types.SpecTypeFAQ}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %v, want %v", i, got[i], want[i])
		}
	}

	// Falls back to the legacy table when the workflow declares no rule.
	got = SourcesForRule(nil, false, types.SpecTypeFAQ)
	if len(got) != 2 || got[0] != types.SpecTypeMRD || got[1] != types.SpecTypePress {
		t.Errorf("legacy fallback for faq = %v, want [mrd press]", got)
	}
}

func TestBuildPromptFromRule_IncludesGuidanceAndPromptContext(t *testing.T) {
	s := NewSynthesizer(nil) // client is unused by buildPrompt
	input := SynthesisInput{}
	input.Set(types.SpecTypePress, "press content")

	rule := &workflow.SynthesisRule{
		Sources:       []string{"press"},
		Guidance:      "GUIDANCE_MARKER",
		PromptContext: "PROMPT_CONTEXT_MARKER",
	}

	prompt, sources := s.buildPrompt(types.SpecTypeFAQ, input, "TEMPLATE_MARKER", rule)

	for _, marker := range []string{"GUIDANCE_MARKER", "PROMPT_CONTEXT_MARKER", "TEMPLATE_MARKER", "press content"} {
		if !strings.Contains(prompt, marker) {
			t.Errorf("prompt missing %q", marker)
		}
	}
	if len(sources) != 1 || sources[0] != types.SpecTypePress {
		t.Errorf("sources = %v, want [press]", sources)
	}
}

func TestBuildLegacyPrompt_Unchanged(t *testing.T) {
	s := NewSynthesizer(nil)
	input := SynthesisInput{}
	input.Set(types.SpecTypeMRD, "mrd content")

	// nil rule -> legacy path, exercised by every pre-RMI-031 caller with
	// no configured workflow.
	prompt, sources := s.buildPrompt(types.SpecTypePress, input, "TEMPLATE", nil)

	if !strings.Contains(prompt, "Working Backwards") {
		t.Error("legacy press prompt should mention Working Backwards")
	}
	if len(sources) != 1 || sources[0] != types.SpecTypeMRD {
		t.Errorf("sources = %v, want [mrd]", sources)
	}
}

// TestLoadWorkflowRule_RealCatalog exercises the exact scenario that was
// broken before RMI-031: under aws-one-way-door, press is human-authored
// (no sources) and faq requires only press, not mrd; under enterprise,
// bmc is synthesizable even though it isn't in the legacy hardcoded list.
func TestLoadWorkflowRule_RealCatalog(t *testing.T) {
	loaded, rule, ok := LoadWorkflowRule(nil, "aws-one-way-door", types.SpecTypePress)
	if loaded == nil {
		t.Fatal("expected a loaded workflow for aws-one-way-door")
	}
	if !ok {
		t.Fatal("expected aws-one-way-door to declare a press rule")
	}
	if len(rule.Sources) != 0 {
		t.Errorf("aws-one-way-door press sources = %v, want none (human-authored)", rule.Sources)
	}
	if CanSynthesizeWithRule(rule, ok, types.SpecTypePress) {
		t.Error("press should not be synthesizable under aws-one-way-door")
	}

	_, faqRule, faqOK := LoadWorkflowRule(nil, "aws-one-way-door", types.SpecTypeFAQ)
	if !faqOK {
		t.Fatal("expected aws-one-way-door to declare a faq rule")
	}
	sources := SourcesForRule(faqRule, faqOK, types.SpecTypeFAQ)
	if len(sources) != 1 || sources[0] != types.SpecTypePress {
		t.Errorf("aws-one-way-door faq sources = %v, want [press] (not mrd)", sources)
	}

	_, bmcRule, bmcOK := LoadWorkflowRule(nil, "enterprise", types.SpecType("bmc"))
	if !bmcOK {
		t.Fatal("expected enterprise to declare a bmc rule")
	}
	if !CanSynthesizeWithRule(bmcRule, bmcOK, types.SpecType("bmc")) {
		t.Error("bmc should be synthesizable under enterprise, even though it's not in the legacy hardcoded list")
	}

	if loaded, _, ok := LoadWorkflowRule(nil, "not-a-real-profile", types.SpecTypeFAQ); loaded != nil || ok {
		t.Errorf("expected nil/false for a nonexistent profile, got loaded=%v ok=%v", loaded, ok)
	}

	if loaded, _, ok := LoadWorkflowRule(nil, "", types.SpecTypeFAQ); loaded != nil || ok {
		t.Errorf("expected nil/false for an empty profile name (no --profile at init), got loaded=%v ok=%v", loaded, ok)
	}
}
