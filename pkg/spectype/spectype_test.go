package spectype_test

import (
	"testing"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/spectype"
)

// validPDLCStages must match pdlc.Stage* string values exactly. Kept as a
// local literal set (not an import) to avoid the import cycle documented on
// spectype.PDLCStage: specification-workflow-spec sits upstream of pdlc in
// the dependency chain (pdlc -> visionspec -> specification-workflow-spec).
var validPDLCStages = map[spectype.PDLCStage]bool{
	spectype.PDLCStageProductDefinition: true,
	spectype.PDLCStageBuilderDefinition: true,
}

func TestEveryPDLCStageIsValid(t *testing.T) {
	for _, st := range spectype.CoreSpecTypes() {
		if st.PDLCStage == "" {
			continue // execution-tracking spec types (plan, roadmap) intentionally carry none
		}
		if !validPDLCStages[st.PDLCStage] {
			t.Errorf("spec type %q has unknown PDLCStage %q", st.ID, st.PDLCStage)
		}
	}
}

func TestExecutionSpecTypesCarryNoPDLCStage(t *testing.T) {
	// plan and roadmap track work across the whole lifecycle, not one
	// content-producing stage — they must stay unset, not default to a
	// stage that would misrepresent them as belonging to it.
	for _, id := range []string{"plan", "roadmap"} {
		found := false
		for _, st := range spectype.CoreSpecTypes() {
			if st.ID != id {
				continue
			}
			found = true
			if st.PDLCStage != "" {
				t.Errorf("spec type %q should have no PDLCStage, got %q", id, st.PDLCStage)
			}
		}
		if !found {
			t.Errorf("expected spec type %q in CoreSpecTypes()", id)
		}
	}
}

func TestNonExecutionSpecTypesHaveAPDLCStage(t *testing.T) {
	// Every spec type outside CategoryExecution should be categorizable into
	// a PDLC stage — an uncategorized non-execution spec type is a gap a
	// downstream stage-based consumer (e.g. Threat Model Spec) can't bucket.
	for _, st := range spectype.CoreSpecTypes() {
		if st.Category == spectype.CategoryExecution {
			continue
		}
		if st.PDLCStage == "" {
			t.Errorf("spec type %q (category %q) has no PDLCStage", st.ID, st.Category)
		}
	}
}

func TestTechnicalCategorySpecTypesMapToBuilderDefinition(t *testing.T) {
	// TRD, TPD, and IRD are all technical planning documents produced before
	// implementation begins — they belong to Builder Definition, not
	// Implementation/Deployment, which consume code/IaC artifacts, never
	// workflow specs.
	for _, st := range spectype.CoreSpecTypes() {
		if st.Category != spectype.CategoryTechnical {
			continue
		}
		if st.PDLCStage != spectype.PDLCStageBuilderDefinition {
			t.Errorf("technical spec type %q has PDLCStage %q, want %q", st.ID, st.PDLCStage, spectype.PDLCStageBuilderDefinition)
		}
	}
}
