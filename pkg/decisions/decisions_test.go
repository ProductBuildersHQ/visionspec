package decisions

import (
	"path/filepath"
	"testing"
	"time"
)

func TestNewDecisionLog(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	if log == nil {
		t.Fatal("NewDecisionLog returned nil")
	}
}

func TestDecisionLog_Add(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	decision := &Decision{
		Title:     "Use PostgreSQL",
		Status:    StatusAccepted,
		Context:   "We need a database",
		Decision:  "Use PostgreSQL",
		Rationale: "Mature, reliable, feature-rich",
		Tags:      []string{"database", "infrastructure"},
	}

	if err := log.Add(decision); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Check that ID was assigned
	if decision.ID == "" {
		t.Error("Decision ID should be assigned")
	}

	// Check that date was set
	if decision.Date.IsZero() {
		t.Error("Decision date should be set")
	}
}

func TestDecisionLog_Get(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	decision := &Decision{
		ID:        "ADR-0001",
		Title:     "Test Decision",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
	}

	if err := log.Add(decision); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	retrieved, err := log.Get("ADR-0001")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Title != "Test Decision" {
		t.Errorf("Title = %s, want Test Decision", retrieved.Title)
	}

	// Try getting non-existent
	_, err = log.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent decision")
	}
}

func TestDecisionLog_List(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	// Add multiple decisions
	for i := 0; i < 3; i++ {
		if err := log.Add(&Decision{
			Title:     "Decision",
			Status:    StatusAccepted,
			Context:   "Context",
			Decision:  "Decision",
			Rationale: "Rationale",
		}); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	list := log.List()
	if len(list) != 3 {
		t.Errorf("List returned %d decisions, want 3", len(list))
	}
}

func TestDecisionLog_ListByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	// Add decisions with different statuses
	statuses := []DecisionStatus{StatusAccepted, StatusAccepted, StatusProposed, StatusRejected}
	for _, status := range statuses {
		if err := log.Add(&Decision{
			Title:     "Decision",
			Status:    status,
			Context:   "Context",
			Decision:  "Decision",
			Rationale: "Rationale",
		}); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	accepted := log.ListByStatus(StatusAccepted)
	if len(accepted) != 2 {
		t.Errorf("ListByStatus(Accepted) = %d, want 2", len(accepted))
	}
}

func TestDecisionLog_ListByTag(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	if err := log.Add(&Decision{
		Title:     "Database Choice",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
		Tags:      []string{"database", "infrastructure"},
	}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := log.Add(&Decision{
		Title:     "API Framework",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
		Tags:      []string{"api", "infrastructure"},
	}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	infra := log.ListByTag("infrastructure")
	if len(infra) != 2 {
		t.Errorf("ListByTag(infrastructure) = %d, want 2", len(infra))
	}

	database := log.ListByTag("database")
	if len(database) != 1 {
		t.Errorf("ListByTag(database) = %d, want 1", len(database))
	}
}

func TestDecisionLog_Supersede(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	old := &Decision{
		ID:        "ADR-0001",
		Title:     "Old Decision",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
	}

	newD := &Decision{
		ID:        "ADR-0002",
		Title:     "New Decision",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "New approach",
		Rationale: "Better rationale",
	}

	if err := log.Add(old); err != nil {
		t.Fatalf("Add old failed: %v", err)
	}

	if err := log.Add(newD); err != nil {
		t.Fatalf("Add new failed: %v", err)
	}

	if err := log.Supersede("ADR-0001", "ADR-0002"); err != nil {
		t.Fatalf("Supersede failed: %v", err)
	}

	oldD, _ := log.Get("ADR-0001")
	if oldD.Status != StatusSuperseded {
		t.Errorf("Old status = %s, want superseded", oldD.Status)
	}

	if oldD.SupersededBy != "ADR-0002" {
		t.Errorf("SupersededBy = %s, want ADR-0002", oldD.SupersededBy)
	}

	newD2, _ := log.Get("ADR-0002")
	if newD2.Supersedes != "ADR-0001" {
		t.Errorf("Supersedes = %s, want ADR-0001", newD2.Supersedes)
	}
}

func TestDecisionLog_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	decision := &Decision{
		ID:        "ADR-0001",
		Title:     "Test Decision",
		Status:    StatusAccepted,
		Date:      time.Now(),
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
		Tags:      []string{"test"},
	}

	if err := log.Add(decision); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := log.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load in new instance
	log2, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog (reload) failed: %v", err)
	}

	loaded, err := log2.Get("ADR-0001")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if loaded.Title != "Test Decision" {
		t.Errorf("Loaded title = %s, want Test Decision", loaded.Title)
	}
}

func TestDecision_RenderMarkdown(t *testing.T) {
	decision := &Decision{
		ID:             "ADR-0001",
		Title:          "Use PostgreSQL",
		Status:         StatusAccepted,
		Date:           time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		DecisionMakers: []string{"Alice", "Bob"},
		Context:        "We need a database for our application.",
		Decision:       "Use PostgreSQL as our primary database.",
		Rationale:      "PostgreSQL is mature, reliable, and feature-rich.",
		Consequences:   []string{"Need PostgreSQL expertise", "Good for complex queries"},
		Alternatives: []Alternative{
			{
				Title:       "MySQL",
				Description: "Popular open-source database",
				Pros:        []string{"Popular", "Good tooling"},
				Cons:        []string{"Less feature-rich"},
				RejectedFor: "Fewer advanced features",
			},
		},
		Related: []string{"ADR-0002"},
		Tags:    []string{"database", "infrastructure"},
	}

	md := decision.RenderMarkdown()

	expectations := []string{
		"# ADR-0001: Use PostgreSQL",
		"✅ Accepted",
		"2024-01-15",
		"Alice, Bob",
		"## Context",
		"## Decision",
		"## Rationale",
		"## Consequences",
		"## Alternatives Considered",
		"MySQL",
		"## Related Decisions",
		"database, infrastructure",
	}

	for _, exp := range expectations {
		if !containsString(md, exp) {
			t.Errorf("Markdown should contain %q", exp)
		}
	}
}

func TestDecisionLog_RenderIndex(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	statuses := []DecisionStatus{StatusAccepted, StatusProposed, StatusRejected}
	for i, status := range statuses {
		if err := log.Add(&Decision{
			ID:        "ADR-000" + string(rune('1'+i)),
			Title:     "Decision " + string(rune('1'+i)),
			Status:    status,
			Date:      time.Now(),
			Context:   "Context",
			Decision:  "Decision",
			Rationale: "Rationale",
		}); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	index := log.RenderIndex()

	if len(index) < 100 {
		t.Error("Index seems too short")
	}

	expectations := []string{
		"# Architecture Decision Records",
		"Accepted",
		"Proposed",
		"Rejected",
	}

	for _, exp := range expectations {
		if !containsString(index, exp) {
			t.Errorf("Index should contain %q", exp)
		}
	}
}

func TestDecisionLog_GetStatistics(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	// Add decisions
	decisions := []struct {
		status  DecisionStatus
		tags    []string
		project string
	}{
		{StatusAccepted, []string{"database"}, "project-a"},
		{StatusAccepted, []string{"database", "infrastructure"}, "project-a"},
		{StatusProposed, []string{"api"}, "project-b"},
	}

	for _, d := range decisions {
		if err := log.Add(&Decision{
			Title:     "Decision",
			Status:    d.status,
			Date:      time.Now(),
			Context:   "Context",
			Decision:  "Decision",
			Rationale: "Rationale",
			Tags:      d.tags,
			Project:   d.project,
		}); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	stats := log.GetStatistics()

	if stats.Total != 3 {
		t.Errorf("Total = %d, want 3", stats.Total)
	}

	if stats.ByStatus["accepted"] != 2 {
		t.Errorf("ByStatus[accepted] = %d, want 2", stats.ByStatus["accepted"])
	}

	if stats.ByTag["database"] != 2 {
		t.Errorf("ByTag[database] = %d, want 2", stats.ByTag["database"])
	}

	if stats.ByProject["project-a"] != 2 {
		t.Errorf("ByProject[project-a] = %d, want 2", stats.ByProject["project-a"])
	}
}

func TestBuildRationaleGraph(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "decisions.yaml")

	log, err := NewDecisionLog(logPath)
	if err != nil {
		t.Fatalf("NewDecisionLog failed: %v", err)
	}

	// Add decisions with relationships
	if err := log.Add(&Decision{
		ID:        "ADR-0001",
		Title:     "Decision 1",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
	}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := log.Add(&Decision{
		ID:        "ADR-0002",
		Title:     "Decision 2",
		Status:    StatusAccepted,
		Context:   "Context",
		Decision:  "Decision",
		Rationale: "Rationale",
		Related:   []string{"ADR-0001"},
	}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	graph := BuildRationaleGraph(log)

	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}
}

func TestRationaleGraph_RenderMermaid(t *testing.T) {
	graph := &RationaleGraph{
		Nodes: []RationaleNode{
			{ID: "ADR-0001", Title: "Decision 1", Status: StatusAccepted},
			{ID: "ADR-0002", Title: "Decision 2", Status: StatusProposed},
		},
		Edges: []RationaleEdge{
			{From: "ADR-0002", To: "ADR-0001", Type: EdgeTypeRelated},
		},
	}

	mermaid := graph.RenderMermaid()

	if len(mermaid) < 50 {
		t.Error("Mermaid output seems too short")
	}

	if !containsString(mermaid, "graph TD") {
		t.Error("Missing graph header")
	}
}

func TestRationaleGraph_GetConnectedComponents(t *testing.T) {
	graph := &RationaleGraph{
		Nodes: []RationaleNode{
			{ID: "a"},
			{ID: "b"},
			{ID: "c"},
			{ID: "d"}, // Disconnected
		},
		Edges: []RationaleEdge{
			{From: "a", To: "b"},
			{From: "b", To: "c"},
		},
	}

	components := graph.GetConnectedComponents()

	// Should have 2 components: {a,b,c} and {d}
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}
}

func TestRationaleGraph_FindPath(t *testing.T) {
	graph := &RationaleGraph{
		Nodes: []RationaleNode{
			{ID: "a"},
			{ID: "b"},
			{ID: "c"},
		},
		Edges: []RationaleEdge{
			{From: "a", To: "b"},
			{From: "b", To: "c"},
		},
	}

	path := graph.FindPath("a", "c")
	if path == nil {
		t.Fatal("Expected to find path")
	}

	if len(path) != 3 {
		t.Errorf("Expected path length 3, got %d", len(path))
	}

	// Test no path
	noPath := graph.FindPath("c", "a")
	if noPath != nil {
		t.Error("Should not find path from c to a")
	}
}

func TestRationaleGraph_AnalyzeImpact(t *testing.T) {
	graph := &RationaleGraph{
		Nodes: []RationaleNode{
			{ID: "core"},
			{ID: "api"},
			{ID: "web"},
		},
		Edges: []RationaleEdge{
			{From: "api", To: "core", Type: EdgeTypeDependsOn},
			{From: "web", To: "api", Type: EdgeTypeDependsOn},
		},
	}

	impact := graph.AnalyzeImpact("core")

	if len(impact.DirectlyAffected) == 0 {
		t.Error("Expected directly affected decisions")
	}
}

// Helper
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
