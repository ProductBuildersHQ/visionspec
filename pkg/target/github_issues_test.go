package target

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGitHubIssuesTarget_Interface(t *testing.T) {
	target := &GitHubIssuesTarget{}

	// Test Name
	if target.Name() != "github-issues" {
		t.Errorf("Name() = %s, want github-issues", target.Name())
	}

	// Test Description
	if target.Description() == "" {
		t.Error("Description() should not be empty")
	}

	// Test Capabilities
	caps := target.Capabilities()
	if !caps.SequentialTasks {
		t.Error("Capabilities().SequentialTasks should be true")
	}
}

func TestGitHubIssuesTarget_Validate(t *testing.T) {
	target := &GitHubIssuesTarget{}

	// Empty spec should fail
	if err := target.Validate(""); err == nil {
		t.Error("Validate should fail for empty spec")
	}

	// Non-empty spec should pass
	if err := target.Validate("# Test Spec"); err != nil {
		t.Errorf("Validate should pass for non-empty spec: %v", err)
	}
}

func TestGitHubIssuesTarget_Export(t *testing.T) {
	target := &GitHubIssuesTarget{}
	tmpDir := t.TempDir()

	spec := `# Test Project Spec

## Backend

- [ ] Design API endpoints
- [ ] Implement authentication
- [x] Set up database

## Frontend

### FR-001 User Dashboard

Build the user dashboard with key metrics.

- Given a user is logged in
- When they navigate to the dashboard
- Then they see their stats

- [ ] Create dashboard component [high]
- [ ] Add chart widgets
`

	config := ExportConfig{
		ProjectName: "test-project",
		OutputDir:   tmpDir,
		Options:     map[string]any{"repository": "owner/repo"},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed")
	}

	if len(result.Files) < 2 {
		t.Error("Export should create at least 2 files (JSON and shell script)")
	}

	// Verify the JSON output file exists and is valid
	jsonFile := filepath.Join(tmpDir, "github-issues.json")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var export GitHubIssuesExport
	if err := json.Unmarshal(data, &export); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if export.Repository != "owner/repo" {
		t.Errorf("Repository = %s, want owner/repo", export.Repository)
	}

	if len(export.Issues) == 0 {
		t.Error("Should have extracted issues")
	}

	// Verify CLI script exists
	cliFile := filepath.Join(tmpDir, "create-issues.sh")
	if _, err := os.Stat(cliFile); os.IsNotExist(err) {
		t.Error("CLI script should exist")
	}
}

func TestGitHubIssuesTarget_ExtractIssues(t *testing.T) {
	target := &GitHubIssuesTarget{}

	spec := `## Section One

- [ ] Open task
- [x] Completed task

### FR-001 Feature One

Feature description.

- [ ] Sub task [critical]
`

	config := ExportConfig{
		ProjectName: "test",
		OutputDir:   t.TempDir(),
	}

	issues := target.extractIssues(spec, config)

	if len(issues) == 0 {
		t.Error("Should extract at least one issue")
	}

	// Check for completed task state
	hasCompleted := false
	hasCritical := false
	for _, issue := range issues {
		if issue.State == "closed" {
			hasCompleted = true
		}
		if issue.Priority == "critical" {
			hasCritical = true
		}
	}

	if !hasCompleted {
		t.Error("Should have a closed issue from completed task")
	}

	if !hasCritical {
		t.Error("Should have extracted critical priority")
	}
}

func TestGitHubIssuesTarget_GenerateCLICommands(t *testing.T) {
	export := &GitHubIssuesExport{
		Repository: "owner/repo",
		Issues: []GitHubIssue{
			{
				Title:  "Test Issue",
				Body:   "Issue body",
				Labels: []string{"bug", "urgent"},
				State:  "open",
			},
			{
				Title: "Closed Issue",
				State: "closed",
			},
		},
	}

	script := export.GenerateCLICommands()

	if script == "" {
		t.Error("Should generate CLI commands")
	}

	// Should have shebang
	if script[:11] != "#!/bin/bash" {
		t.Error("Script should start with shebang")
	}

	// Should have gh issue create
	if !contains(script, "gh issue create") {
		t.Error("Script should contain gh issue create")
	}

	// Should not create closed issues
	if contains(script, "Closed Issue") {
		t.Error("Script should not create closed issues")
	}
}

func TestGitHubIssuesTarget_Registration(t *testing.T) {
	target, err := Get("github-issues")
	if err != nil {
		t.Fatalf("GitHubIssues target should be registered: %v", err)
	}

	if target.Name() != "github-issues" {
		t.Errorf("Got wrong target: %s", target.Name())
	}
}

func TestGHSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"UPPER_CASE", "upper-case"},
		{"already-slug", "already-slug"},
	}

	for _, tt := range tests {
		result := ghSlugify(tt.input)
		if result != tt.expected {
			t.Errorf("ghSlugify(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractGHRequirementID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FR-001 User Login", "FR-001"},
		{"This is NFR-123 here", "NFR-123"},
		{"US-42 Story", "US-42"},
		{"No requirement", ""},
	}

	for _, tt := range tests {
		result := extractGHRequirementID(tt.input)
		if result != tt.expected {
			t.Errorf("extractGHRequirementID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
