package target

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJiraTarget_Interface(t *testing.T) {
	target := &JiraTarget{}

	// Test Name
	if target.Name() != "jira" {
		t.Errorf("Name() = %s, want jira", target.Name())
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
	if !caps.DependencyGraph {
		t.Error("Capabilities().DependencyGraph should be true")
	}
}

func TestJiraTarget_Validate(t *testing.T) {
	target := &JiraTarget{}

	// Empty spec should fail
	if err := target.Validate(""); err == nil {
		t.Error("Validate should fail for empty spec")
	}

	// Non-empty spec should pass
	if err := target.Validate("# Test Spec"); err != nil {
		t.Errorf("Validate should pass for non-empty spec: %v", err)
	}
}

func TestJiraTarget_Export(t *testing.T) {
	target := &JiraTarget{}
	tmpDir := t.TempDir()

	spec := `# Test Project Spec

## Backend Epic

### FR-001 User Authentication

Implement user authentication.

- [ ] Design auth flow
- [ ] Implement JWT tokens [high]
- [x] Set up OAuth

## Frontend Epic

### US-002 Dashboard

Build user dashboard.

- Given a user is logged in
- When they view dashboard
- Then they see metrics

- [ ] Create components [critical]
`

	config := ExportConfig{
		ProjectName: "test-project",
		OutputDir:   tmpDir,
		Options: map[string]any{
			"project_key":  "TEST",
			"project_name": "Test Project",
		},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed")
	}

	if len(result.Files) < 2 {
		t.Error("Export should create at least 2 files (JSON and CSV)")
	}

	// Verify the JSON output file exists and is valid
	jsonFile := filepath.Join(tmpDir, "jira-issues.json")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var export JiraExport
	if err := json.Unmarshal(data, &export); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if export.Project.Key != "TEST" {
		t.Errorf("Project.Key = %s, want TEST", export.Project.Key)
	}

	if len(export.Epics) == 0 {
		t.Error("Should have extracted epics")
	}

	if len(export.Issues) == 0 {
		t.Error("Should have extracted issues")
	}

	// Verify CSV file exists
	csvFile := filepath.Join(tmpDir, "jira-import.csv")
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Error("CSV file should exist")
	}
}

func TestJiraTarget_ExtractIssues(t *testing.T) {
	target := &JiraTarget{}

	spec := `## Epic One

### Story One

Description here.

- [ ] Task A
- [x] Task B [high]
- [ ] Task C [critical]

### US-001 User Story

Story description.

- Given condition
- When action
- Then result
`

	epics, issues := target.extractIssues(spec)

	if len(epics) == 0 {
		t.Error("Should extract at least one epic")
	}

	if len(issues) == 0 {
		t.Error("Should extract at least one issue")
	}

	// Check for story type
	hasStory := false
	hasTask := false
	hasDoneTask := false
	hasHighPriority := false

	for _, issue := range issues {
		if issue.IssueType == "Story" {
			hasStory = true
			// Tasks under a story become subtasks
			for _, subtask := range issue.SubTasks {
				if subtask.IssueType == "Task" {
					hasTask = true
					if subtask.Status == "Done" {
						hasDoneTask = true
					}
					if subtask.Priority == "High" {
						hasHighPriority = true
					}
				}
			}
		}
		// Also check top-level tasks
		if issue.IssueType == "Task" {
			hasTask = true
			if issue.Status == "Done" {
				hasDoneTask = true
			}
			if issue.Priority == "High" {
				hasHighPriority = true
			}
		}
	}

	if !hasStory {
		t.Error("Should have a Story issue type")
	}

	if !hasTask {
		t.Error("Should have Task issue types (as subtasks)")
	}

	if !hasDoneTask {
		t.Error("Should have a Done task from [x] marker")
	}

	if !hasHighPriority {
		t.Error("Should have extracted high priority")
	}
}

func TestJiraExport_GenerateCSV(t *testing.T) {
	export := &JiraExport{
		Epics: []JiraIssue{
			{
				Summary:   "Test Epic",
				IssueType: "Epic",
			},
		},
		Issues: []JiraIssue{
			{
				Summary:     "Test Story",
				Description: "Story description",
				IssueType:   "Story",
				Priority:    "High",
				Labels:      []string{"backend"},
				Epic:        "Test Epic",
				StoryPoints: 5,
			},
			{
				Summary:    "Quote \"test\"",
				IssueType:  "Task",
				AcceptCrit: "- Given X\n- Then Y",
			},
		},
	}

	csv := export.GenerateCSV()

	if csv == "" {
		t.Error("Should generate CSV content")
	}

	// Check header
	if !strings.HasPrefix(csv, "Summary,Description,Issue Type") {
		t.Error("CSV should start with header row")
	}

	// Check epic row
	if !strings.Contains(csv, "Test Epic") {
		t.Error("CSV should contain epic")
	}

	// Check story row
	if !strings.Contains(csv, "Test Story") {
		t.Error("CSV should contain story")
	}

	// Check story points
	if !strings.Contains(csv, "5") {
		t.Error("CSV should contain story points")
	}
}

func TestJiraTarget_Registration(t *testing.T) {
	target, err := Get("jira")
	if err != nil {
		t.Fatalf("Jira target should be registered: %v", err)
	}

	if target.Name() != "jira" {
		t.Errorf("Got wrong target: %s", target.Name())
	}
}

func TestExtractJiraReqID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FR-001 Feature", "FR-001"},
		{"NFR-123 Performance", "NFR-123"},
		{"REQ-42 Requirement", "REQ-42"},
		{"US-999 User Story", "US-999"},
		{"STORY-1 Epic story", "STORY-1"},
		{"No requirement here", ""},
	}

	for _, tt := range tests {
		result := extractJiraReqID(tt.input)
		if result != tt.expected {
			t.Errorf("extractJiraReqID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestJiraCSVEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"has,comma", "\"has,comma\""},
		{"has\"quote", "\"has\"\"quote\""},
		{"has\nnewline", "\"has\nnewline\""},
		{"normal text", "normal text"},
	}

	for _, tt := range tests {
		result := jiraCSVEscape(tt.input)
		if result != tt.expected {
			t.Errorf("jiraCSVEscape(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
