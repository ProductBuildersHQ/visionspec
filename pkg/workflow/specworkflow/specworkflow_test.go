package specworkflow

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
)

func TestFromProfile(t *testing.T) {
	loader := profiles.DefaultLoader()
	profile, err := loader.Load("big-tech-product")
	if err != nil {
		t.Skipf("profile not available: %v", err)
	}

	w, err := FromProfile(profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Name != profile.Name {
		t.Errorf("expected workflow name %q, got %q", profile.Name, w.Name)
	}

	// Should have phases
	if len(w.Phases) == 0 {
		t.Error("expected phases in workflow")
	}

	// Should have nodes from profile's spec config
	if len(w.Nodes) == 0 {
		t.Error("expected nodes in workflow")
	}
}

func TestUpdateFromProject(t *testing.T) {
	// Create a simple workflow
	w, _ := workflow.NewBuilder("test").
		Phase("discovery", "Discovery", 1).
		Node("mrd", "MRD").Add().
		Node("prd", "PRD").DependsOn("mrd").Add().
		Build()

	// Create a project with spec statuses
	project := &types.Project{
		Specs: map[types.SpecType]*types.Spec{
			types.SpecTypeMRD: {
				Type:   types.SpecTypeMRD,
				Status: types.StatusApproved,
			},
			types.SpecTypePRD: {
				Type:   types.SpecTypePRD,
				Status: types.StatusDraft,
			},
		},
	}

	UpdateFromProject(w, project)

	mrd, _ := w.GetNode("mrd")
	if mrd.Status != workflow.StatusCompleted {
		t.Errorf("expected mrd to be completed, got %s", mrd.Status)
	}

	prd, _ := w.GetNode("prd")
	if prd.Status != workflow.StatusInProgress {
		t.Errorf("expected prd to be in_progress, got %s", prd.Status)
	}
}

func TestUpdateFromProjectNil(t *testing.T) {
	w, _ := workflow.NewBuilder("test").
		Phase("p1", "Phase 1", 1).
		Node("n1", "Node 1").Add().
		Build()

	// Should not panic with nil project
	UpdateFromProject(w, nil)

	// Should not panic with project without specs
	UpdateFromProject(w, &types.Project{})
}

func TestCategoryToPhase(t *testing.T) {
	tests := []struct {
		category types.SpecCategory
		expected string
	}{
		{types.CategorySource, "discovery"},
		{types.CategoryGTM, "vision"},
		{types.CategoryTechnical, "technical"},
		{types.CategoryOutput, "reconcile"},
	}

	for _, tc := range tests {
		result := categoryToPhase(tc.category)
		if result != tc.expected {
			t.Errorf("categoryToPhase(%s) = %s, expected %s", tc.category, result, tc.expected)
		}
	}
}

func TestCategoryToNodeType(t *testing.T) {
	tests := []struct {
		category types.SpecCategory
		expected workflow.NodeType
	}{
		{types.CategorySource, TypeSource},
		{types.CategoryGTM, TypeGTM},
		{types.CategoryTechnical, TypeTechnical},
		{types.CategoryOutput, TypeOutput},
	}

	for _, tc := range tests {
		result := categoryToNodeType(tc.category)
		if result != tc.expected {
			t.Errorf("categoryToNodeType(%s) = %s, expected %s", tc.category, result, tc.expected)
		}
	}
}

func TestIsAutomated(t *testing.T) {
	tests := []struct {
		category types.SpecCategory
		expected bool
	}{
		{types.CategorySource, false},
		{types.CategoryGTM, true},
		{types.CategoryTechnical, true},
		{types.CategoryOutput, true},
	}

	for _, tc := range tests {
		result := isAutomated(tc.category)
		if result != tc.expected {
			t.Errorf("isAutomated(%s) = %v, expected %v", tc.category, result, tc.expected)
		}
	}
}
