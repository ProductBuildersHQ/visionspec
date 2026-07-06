package workflow

import (
	"testing"
)

func TestNewWorkflow(t *testing.T) {
	w := New("test-workflow")
	if w.Name != "test-workflow" {
		t.Errorf("expected name 'test-workflow', got %q", w.Name)
	}
	if len(w.Nodes) != 0 {
		t.Errorf("expected empty nodes, got %d", len(w.Nodes))
	}
	if len(w.Phases) != 0 {
		t.Errorf("expected empty phases, got %d", len(w.Phases))
	}
}

func TestAddPhase(t *testing.T) {
	w := New("test")
	w.AddPhase("phase1", "Phase 1", 1)
	w.AddPhase("phase2", "Phase 2", 2)

	if len(w.Phases) != 2 {
		t.Errorf("expected 2 phases, got %d", len(w.Phases))
	}
	if w.Phases[0].ID != "phase1" {
		t.Errorf("expected first phase 'phase1', got %q", w.Phases[0].ID)
	}
}

func TestAddNode(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	err := w.AddNode(&Node{
		ID:    "node1",
		Name:  "Node 1",
		Phase: "p1",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	node, ok := w.GetNode("node1")
	if !ok {
		t.Error("expected to find node1")
	}
	if node.Status != StatusPending {
		t.Errorf("expected status pending, got %q", node.Status)
	}
}

func TestDependencies(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1", DependsOn: []string{"a"}})
	w.AddNode(&Node{ID: "c", Name: "C", Phase: "p1", DependsOn: []string{"a", "b"}})

	deps := w.Dependencies("c")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(deps))
	}

	dependents := w.Dependents("a")
	if len(dependents) != 2 {
		t.Errorf("expected 2 dependents, got %d", len(dependents))
	}
}

func TestIsReady(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1", DependsOn: []string{"a"}})

	// Initially, a is ready (no deps), b is not
	if !w.IsReady("a") {
		t.Error("expected node a to be ready")
	}
	if w.IsReady("b") {
		t.Error("expected node b to not be ready")
	}

	// Complete a, then b should be ready
	w.UpdateStatus("a", StatusCompleted)
	if !w.IsReady("b") {
		t.Error("expected node b to be ready after a completed")
	}
}

func TestReadyNodes(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1"})
	w.AddNode(&Node{ID: "c", Name: "C", Phase: "p1", DependsOn: []string{"a"}})

	ready := w.ReadyNodes()
	if len(ready) != 2 {
		t.Errorf("expected 2 ready nodes, got %d", len(ready))
	}
}

func TestTopologicalSort(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1", DependsOn: []string{"a"}})
	w.AddNode(&Node{ID: "c", Name: "C", Phase: "p1", DependsOn: []string{"b"}})

	sorted, err := w.TopologicalSort()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should be ordered a -> b -> c
	ids := make([]string, len(sorted))
	for i, n := range sorted {
		ids[i] = n.ID
	}

	if ids[0] != "a" || ids[1] != "b" || ids[2] != "c" {
		t.Errorf("expected order [a, b, c], got %v", ids)
	}
}

func TestCycleDetection(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1", DependsOn: []string{"c"}})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1", DependsOn: []string{"a"}})
	w.AddNode(&Node{ID: "c", Name: "C", Phase: "p1", DependsOn: []string{"b"}})

	_, err := w.TopologicalSort()
	if err == nil {
		t.Error("expected cycle detection error")
	}
}

func TestProgress(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1"})
	w.AddNode(&Node{ID: "c", Name: "C", Phase: "p1"})
	w.AddNode(&Node{ID: "d", Name: "D", Phase: "p1"})

	completed, total, percent := w.Progress()
	if completed != 0 || total != 4 || percent != 0 {
		t.Errorf("expected 0/4 (0%%), got %d/%d (%.1f%%)", completed, total, percent)
	}

	w.UpdateStatus("a", StatusCompleted)
	w.UpdateStatus("b", StatusSkipped)

	completed, total, percent = w.Progress()
	if completed != 2 || total != 4 || percent != 50 {
		t.Errorf("expected 2/4 (50%%), got %d/%d (%.1f%%)", completed, total, percent)
	}
}

func TestValidate(t *testing.T) {
	w := New("test")
	w.AddPhase("p1", "Phase 1", 1)

	w.AddNode(&Node{ID: "a", Name: "A", Phase: "p1"})
	w.AddNode(&Node{ID: "b", Name: "B", Phase: "p1", DependsOn: []string{"nonexistent"}})

	err := w.Validate()
	if err == nil {
		t.Error("expected validation error for missing dependency")
	}
}
