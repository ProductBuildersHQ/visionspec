package rubrics

import (
	"testing"

	swf "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
	sws "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/types"
)

func TestLoaderForWorkflow(t *testing.T) {
	w, err := sws.DefaultLoader().Load("aws-two-way-door")
	if err != nil {
		t.Fatalf("loading aws-two-way-door workflow: %v", err)
	}

	loader := LoaderForWorkflow(w)

	t.Run("loads workflow rubric", func(t *testing.T) {
		rs, err := loader.Load(types.SpecType("press"))
		if err != nil {
			t.Fatalf("loading press rubric: %v", err)
		}
		if rs.ID == "" {
			t.Error("press rubric has no ID")
		}
		if len(rs.Categories) == 0 {
			t.Error("press rubric has no categories")
		}
	})

	t.Run("unknown rubric errors", func(t *testing.T) {
		if _, err := loader.Load(types.SpecType("does-not-exist")); err == nil {
			t.Error("expected error for unknown rubric")
		}
	})
}

func TestLoaderForWorkflowNoRubrics(t *testing.T) {
	// A workflow with no rubrics of its own falls back to the default loader.
	w := &sws.LoadedWorkflow{Workflow: &swf.Workflow{Name: "empty"}}

	loader := LoaderForWorkflow(w)
	if _, err := loader.Load(types.SpecTypePRD); err != nil {
		t.Errorf("expected default rubric fallback for prd: %v", err)
	}
}
