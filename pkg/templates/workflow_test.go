package templates

import (
	"testing"

	swf "github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	sws "github.com/ProductBuildersHQ/visionspec/pkg/workflows"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestLoaderForWorkflow(t *testing.T) {
	w, err := sws.DefaultLoader().Load("aws-two-way-door")
	if err != nil {
		t.Fatalf("loading aws-two-way-door workflow: %v", err)
	}

	loader := LoaderForWorkflow(w)

	t.Run("loads workflow template", func(t *testing.T) {
		tmpl, err := loader.Load(types.SpecType("press"))
		if err != nil {
			t.Fatalf("loading press template: %v", err)
		}
		if tmpl.Content == "" {
			t.Error("press template has no content")
		}
		if tmpl.SpecType != types.SpecType("press") {
			t.Errorf("SpecType = %q, want %q", tmpl.SpecType, "press")
		}
	})

	t.Run("falls back to canvas templates", func(t *testing.T) {
		// opportunity-spec canvas comes from prism-roadmap
		if _, err := loader.Load(types.SpecType("opportunity-spec")); err != nil {
			t.Errorf("loading opportunity-spec canvas: %v", err)
		}
	})

	t.Run("unknown template errors", func(t *testing.T) {
		if _, err := loader.Load(types.SpecType("does-not-exist")); err == nil {
			t.Error("expected error for unknown template")
		}
	})
}

func TestLoaderForWorkflowNoTemplates(t *testing.T) {
	// A workflow with no templates of its own falls back to the default loader.
	w := &sws.LoadedWorkflow{Workflow: &swf.Workflow{Name: "empty"}}

	loader := LoaderForWorkflow(w)
	if _, err := loader.Load(types.SpecTypePRD); err != nil {
		t.Errorf("expected default template fallback for prd: %v", err)
	}
}

func TestMapLoaderAvailable(t *testing.T) {
	w, err := sws.DefaultLoader().Load("aws-two-way-door")
	if err != nil {
		t.Fatalf("loading aws-two-way-door workflow: %v", err)
	}

	loader := NewMapLoader(w.Templates)
	available := loader.Available()
	if len(available) != len(w.Templates) {
		t.Errorf("Available() = %d entries, want %d", len(available), len(w.Templates))
	}
}
