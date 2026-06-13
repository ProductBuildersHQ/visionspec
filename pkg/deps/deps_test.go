package deps

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.projectPath != tmpDir {
		t.Errorf("projectPath = %s, want %s", manager.projectPath, tmpDir)
	}
}

func TestManager_GetDependencies_NoDeps(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	deps, err := manager.GetDependencies()
	if err != nil {
		t.Fatalf("GetDependencies failed: %v", err)
	}

	if deps != nil && len(deps) > 0 {
		t.Errorf("Expected no dependencies, got %d", len(deps))
	}
}

func TestManager_AddDependency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial config
	configPath := filepath.Join(tmpDir, "visionspec.yaml")
	if err := os.WriteFile(configPath, []byte("name: test-project\n"), 0600); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	manager := NewManager(tmpDir)

	dep := Dependency{
		Project:     "other-project",
		Type:        DepTypeRequires,
		Required:    true,
		Description: "Required for auth",
	}

	if err := manager.AddDependency(dep); err != nil {
		t.Fatalf("AddDependency failed: %v", err)
	}

	// Verify dependency was added
	deps, err := manager.GetDependencies()
	if err != nil {
		t.Fatalf("GetDependencies failed: %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dependency, got %d", len(deps))
	}

	if deps[0].Project != "other-project" {
		t.Errorf("Project = %s, want other-project", deps[0].Project)
	}
}

func TestManager_RemoveDependency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with a dependency
	configPath := filepath.Join(tmpDir, "visionspec.yaml")
	config := `name: test-project
dependencies:
  - project: dep-to-remove
    type: requires
    required: true
  - project: dep-to-keep
    type: related
    required: false
`
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	manager := NewManager(tmpDir)

	if err := manager.RemoveDependency("dep-to-remove"); err != nil {
		t.Fatalf("RemoveDependency failed: %v", err)
	}

	deps, err := manager.GetDependencies()
	if err != nil {
		t.Fatalf("GetDependencies failed: %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dependency after removal, got %d", len(deps))
	}

	if deps[0].Project != "dep-to-keep" {
		t.Errorf("Remaining project = %s, want dep-to-keep", deps[0].Project)
	}
}

func TestDependencyGraph_DetectCycles(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"a": {Name: "a", Dependencies: []Dependency{{Project: "b", Type: DepTypeRequires}}},
			"b": {Name: "b", Dependencies: []Dependency{{Project: "c", Type: DepTypeRequires}}},
			"c": {Name: "c", Dependencies: []Dependency{{Project: "a", Type: DepTypeRequires}}},
		},
		Edges: []DependencyEdge{
			{From: "a", To: "b", Type: DepTypeRequires},
			{From: "b", To: "c", Type: DepTypeRequires},
			{From: "c", To: "a", Type: DepTypeRequires},
		},
	}

	cycles := graph.DetectCycles()
	if len(cycles) == 0 {
		t.Error("Expected to detect a cycle")
	}
}

func TestDependencyGraph_TopologicalSort(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"a": {Name: "a", Dependencies: []Dependency{}},
			"b": {Name: "b", Dependencies: []Dependency{{Project: "a", Type: DepTypeRequires}}},
			"c": {Name: "c", Dependencies: []Dependency{{Project: "b", Type: DepTypeRequires}}},
		},
		Edges: []DependencyEdge{
			{From: "b", To: "a", Type: DepTypeRequires},
			{From: "c", To: "b", Type: DepTypeRequires},
		},
	}

	sorted, err := graph.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	// a should come before b, b should come before c
	aIdx, bIdx, cIdx := -1, -1, -1
	for i, name := range sorted {
		switch name {
		case "a":
			aIdx = i
		case "b":
			bIdx = i
		case "c":
			cIdx = i
		}
	}

	if aIdx > bIdx || bIdx > cIdx {
		t.Errorf("Invalid sort order: %v", sorted)
	}
}

func TestDependencyGraph_Validate(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"a": {Name: "a", Dependencies: []Dependency{
				{Project: "missing", Type: DepTypeRequires, Required: true},
			}},
		},
		Edges: []DependencyEdge{
			{From: "a", To: "missing", Type: DepTypeRequires, Required: true},
		},
	}

	result := graph.Validate()
	if result.Valid {
		t.Error("Expected validation to fail for missing required dependency")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected at least one error")
	}
}

func TestDependencyGraph_ExportMermaid(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"project-a": {Name: "project-a", Status: "completed"},
			"project-b": {Name: "project-b", Status: "draft"},
		},
		Edges: []DependencyEdge{
			{From: "project-b", To: "project-a", Type: DepTypeRequires},
		},
	}

	mermaid := graph.ExportMermaid()

	if mermaid == "" {
		t.Error("ExportMermaid returned empty string")
	}

	if len(mermaid) < 50 {
		t.Error("Mermaid output seems too short")
	}

	// Check for key elements
	if !contains(mermaid, "graph TD") {
		t.Error("Missing graph header")
	}
}

func TestDependencyGraph_GetAffectedProjects(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"core":   {Name: "core", Dependents: []string{"api", "web"}},
			"api":    {Name: "api", Dependencies: []Dependency{{Project: "core"}}, Dependents: []string{"mobile"}},
			"web":    {Name: "web", Dependencies: []Dependency{{Project: "core"}}},
			"mobile": {Name: "mobile", Dependencies: []Dependency{{Project: "api"}}},
		},
	}

	affected := graph.GetAffectedProjects("core")

	// Should include api, web, and mobile (transitively)
	if len(affected) < 2 {
		t.Errorf("Expected at least 2 affected projects, got %d", len(affected))
	}
}

func TestDependencyGraph_GenerateReport(t *testing.T) {
	graph := &DependencyGraph{
		Projects: map[string]*ProjectNode{
			"a": {Name: "a", Dependencies: []Dependency{}},
			"b": {Name: "b", Dependencies: []Dependency{{Project: "a"}}},
			"c": {Name: "c", Dependencies: []Dependency{}}, // orphan (no deps, no dependents)
		},
		Edges: []DependencyEdge{
			{From: "b", To: "a", Type: DepTypeRequires},
		},
	}

	// Set dependents for proper categorization
	graph.Projects["a"].Dependents = []string{"b"}

	report := graph.GenerateReport()

	if report.TotalProjects != 3 {
		t.Errorf("TotalProjects = %d, want 3", report.TotalProjects)
	}

	if report.ValidationResult == nil {
		t.Error("ValidationResult should not be nil")
	}
}

// Helper
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
