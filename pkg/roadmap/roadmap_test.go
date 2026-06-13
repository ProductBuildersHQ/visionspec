package roadmap

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewGenerator(t *testing.T) {
	tmpDir := t.TempDir()
	gen := NewGenerator(tmpDir)

	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}

	if gen.specsRoot != tmpDir {
		t.Errorf("specsRoot = %s, want %s", gen.specsRoot, tmpDir)
	}
}

func TestGenerator_Generate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a project
	projectDir := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	config := `project:
  name: test-project
  description: A test project
  phase: phase-1
  priority: high
`
	if err := os.WriteFile(filepath.Join(projectDir, "visionspec.yaml"), []byte(config), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	gen := NewGenerator(tmpDir)
	roadmap, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if roadmap == nil {
		t.Fatal("Generate returned nil roadmap")
	}

	// Should have at least one phase
	if len(roadmap.Phases) == 0 {
		t.Error("Expected at least one phase")
	}
}

func TestRoadmap_RenderMarkdown(t *testing.T) {
	now := time.Now()
	roadmap := &Roadmap{
		Title:       "Test Roadmap",
		Description: "A test roadmap",
		Version:     "1.0",
		UpdatedAt:   now,
		Phases: []Phase{
			{
				ID:          "phase-1",
				Name:        "Phase One",
				Description: "First phase",
				Status:      PhaseStatusInProgress,
				Projects:    []string{"project-a", "project-b"},
				Milestones: []Milestone{
					{
						ID:       "ms-1",
						Name:     "Milestone 1",
						Status:   MilestoneStatusCompleted,
						Priority: PriorityHigh,
						RMI:      "RMI-001",
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := roadmap.RenderMarkdown(&buf); err != nil {
		t.Fatalf("RenderMarkdown failed: %v", err)
	}

	output := buf.String()

	// Check for key elements
	if len(output) < 100 {
		t.Error("Output seems too short")
	}

	expectations := []string{
		"# Test Roadmap",
		"Phase One",
		"Milestone 1",
		"project-a",
		"RMI-001",
	}

	for _, exp := range expectations {
		if !containsString(output, exp) {
			t.Errorf("Output should contain %q", exp)
		}
	}
}

func TestRoadmap_GetProgress(t *testing.T) {
	roadmap := &Roadmap{
		Phases: []Phase{
			{
				Milestones: []Milestone{
					{Status: MilestoneStatusCompleted},
					{Status: MilestoneStatusCompleted},
					{Status: MilestoneStatusPending},
					{Status: MilestoneStatusInProgress},
				},
			},
		},
	}

	progress := roadmap.GetProgress()

	// 2 out of 4 = 50%
	if progress != 50.0 {
		t.Errorf("Progress = %f, want 50", progress)
	}
}

func TestRoadmap_AddMilestone(t *testing.T) {
	roadmap := &Roadmap{
		Phases: []Phase{
			{ID: "phase-1", Name: "Phase One"},
		},
	}

	ms := Milestone{
		ID:       "ms-new",
		Name:     "New Milestone",
		Status:   MilestoneStatusPending,
		Priority: PriorityMedium,
	}

	if err := roadmap.AddMilestone("phase-1", ms); err != nil {
		t.Fatalf("AddMilestone failed: %v", err)
	}

	if len(roadmap.Phases[0].Milestones) != 1 {
		t.Errorf("Expected 1 milestone, got %d", len(roadmap.Phases[0].Milestones))
	}

	// Try adding to non-existent phase
	err := roadmap.AddMilestone("nonexistent", ms)
	if err == nil {
		t.Error("Expected error for non-existent phase")
	}
}

func TestRoadmap_UpdateMilestoneStatus(t *testing.T) {
	roadmap := &Roadmap{
		Phases: []Phase{
			{
				ID: "phase-1",
				Milestones: []Milestone{
					{ID: "ms-1", Status: MilestoneStatusPending},
				},
			},
		},
	}

	if err := roadmap.UpdateMilestoneStatus("ms-1", MilestoneStatusCompleted); err != nil {
		t.Fatalf("UpdateMilestoneStatus failed: %v", err)
	}

	if roadmap.Phases[0].Milestones[0].Status != MilestoneStatusCompleted {
		t.Error("Status was not updated")
	}

	// Try updating non-existent milestone
	err := roadmap.UpdateMilestoneStatus("nonexistent", MilestoneStatusCompleted)
	if err == nil {
		t.Error("Expected error for non-existent milestone")
	}
}

func TestRoadmap_Save(t *testing.T) {
	tmpDir := t.TempDir()

	roadmap := &Roadmap{
		Title:     "Test Roadmap",
		UpdatedAt: time.Now(),
		Phases: []Phase{
			{ID: "phase-1", Name: "Phase One"},
		},
	}

	if err := roadmap.Save(tmpDir); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Check YAML file exists
	yamlPath := filepath.Join(tmpDir, "ROADMAP.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Errorf("ROADMAP.yaml not created: %v", err)
	}

	// Check MD file exists
	mdPath := filepath.Join(tmpDir, "ROADMAP.md")
	if _, err := os.Stat(mdPath); err != nil {
		t.Errorf("ROADMAP.md not created: %v", err)
	}
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a roadmap YAML
	content := `title: Test Roadmap
version: "1.0"
phases:
  - id: phase-1
    name: Phase One
    status: in_progress
`
	yamlPath := filepath.Join(tmpDir, "ROADMAP.yaml")
	if err := os.WriteFile(yamlPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write YAML: %v", err)
	}

	roadmap, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if roadmap.Title != "Test Roadmap" {
		t.Errorf("Title = %s, want Test Roadmap", roadmap.Title)
	}

	if len(roadmap.Phases) != 1 {
		t.Errorf("Expected 1 phase, got %d", len(roadmap.Phases))
	}
}

func TestFormatPhaseName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"phase-one", "Phase One"},
		{"backlog", "Backlog"},
		{"q1-2024", "Q1 2024"},
	}

	for _, tt := range tests {
		got := formatPhaseName(tt.input)
		if got != tt.want {
			t.Errorf("formatPhaseName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetStatusBadge(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{"completed", "✅ Completed"},
		{"in_progress", "🚧 In Progress"},
		{"blocked", "🚫 Blocked"},
		{"planned", "📋 Planned"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		got := getStatusBadge(tt.status)
		if got != tt.want {
			t.Errorf("getStatusBadge(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

// Helper
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
