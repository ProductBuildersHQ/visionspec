package types

import (
	"testing"

	swf "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
)

func TestSpecConfigFromWorkflow(t *testing.T) {
	w := &swf.Workflow{
		Name: "test",
		SpecConfig: map[string]*swf.SpecRequirement{
			"prd":   {Required: true, Category: "source"},
			"press": {Required: false, Category: "gtm"},
		},
	}

	sc := SpecConfigFromWorkflow(w)
	if sc == nil {
		t.Fatal("SpecConfigFromWorkflow returned nil")
	}

	if !sc.IsRequired("prd") {
		t.Error("prd should be required")
	}
	if sc.GetCategory("press") != CategoryGTM {
		t.Errorf("press category = %q, want %q", sc.GetCategory("press"), CategoryGTM)
	}
	// Template/Rubric name overrides no longer come through
	// SpecConfigFromWorkflow: sws's Template/Rubric fields are provenance
	// pointers (*workflow.SpecSource), not name overrides. GetTemplate
	// falls back to the spec type's own name.
	if sc.GetTemplate("press") != "press" {
		t.Errorf("press template = %q, want %q", sc.GetTemplate("press"), "press")
	}
}

func TestSpecConfigFromWorkflowNil(t *testing.T) {
	if sc := SpecConfigFromWorkflow(nil); sc != nil {
		t.Error("expected nil for nil workflow")
	}
	if sc := SpecConfigFromWorkflow(&swf.Workflow{Name: "empty"}); sc != nil {
		t.Error("expected nil for workflow without spec config")
	}
}
