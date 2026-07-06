package workflow

import "testing"

func TestBuilder(t *testing.T) {
	w, err := NewBuilder("test-workflow").
		Description("A test workflow").
		Phase("phase1", "Phase 1", 1).
		Node("node1", "Node 1").
		Description("First node").
		Type("source").
		Add().
		Node("node2", "Node 2").
		DependsOn("node1").
		Automated().
		Add().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Name != "test-workflow" {
		t.Errorf("expected name 'test-workflow', got %q", w.Name)
	}

	if w.Description != "A test workflow" {
		t.Errorf("expected description, got %q", w.Description)
	}

	if len(w.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(w.Nodes))
	}

	node1, ok := w.GetNode("node1")
	if !ok {
		t.Fatal("expected node1")
	}
	if node1.Type != "source" {
		t.Errorf("expected type 'source', got %q", node1.Type)
	}
	if node1.Automated {
		t.Error("expected node1 to not be automated")
	}

	node2, ok := w.GetNode("node2")
	if !ok {
		t.Fatal("expected node2")
	}
	if len(node2.DependsOn) != 1 || node2.DependsOn[0] != "node1" {
		t.Errorf("expected dependency on node1, got %v", node2.DependsOn)
	}
	if !node2.Automated {
		t.Error("expected node2 to be automated")
	}
}

func TestBuilderMetadata(t *testing.T) {
	w, err := NewBuilder("test").
		Phase("p1", "Phase 1", 1).
		Node("n1", "Node 1").
		Metadata("key1", "value1").
		Metadata("key2", 42).
		Add().
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	node, _ := w.GetNode("n1")
	if node.Metadata["key1"] != "value1" {
		t.Errorf("expected key1='value1', got %v", node.Metadata["key1"])
	}
	if node.Metadata["key2"] != 42 {
		t.Errorf("expected key2=42, got %v", node.Metadata["key2"])
	}
}

func TestBuilderValidation(t *testing.T) {
	_, err := NewBuilder("test").
		Phase("p1", "Phase 1", 1).
		Node("n1", "Node 1").
		DependsOn("nonexistent").
		Add().
		Build()

	if err == nil {
		t.Error("expected validation error for missing dependency")
	}
}

func TestMustBuild(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustBuild with invalid workflow")
		}
	}()

	NewBuilder("test").
		Phase("p1", "Phase 1", 1).
		Node("n1", "Node 1").
		DependsOn("nonexistent").
		Add().
		MustBuild()
}
