// Package target provides a Jira export adapter.
package target

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	Register(&JiraTarget{})
}

// JiraTarget exports specs as Jira issues.
type JiraTarget struct{}

// JiraIssue represents a Jira issue structure.
type JiraIssue struct {
	Key          string            `json:"key,omitempty"`
	Summary      string            `json:"summary"`
	Description  string            `json:"description"`
	IssueType    string            `json:"issuetype"`
	Priority     string            `json:"priority,omitempty"`
	Labels       []string          `json:"labels,omitempty"`
	Components   []string          `json:"components,omitempty"`
	FixVersions  []string          `json:"fixVersions,omitempty"`
	Sprint       string            `json:"sprint,omitempty"`
	Epic         string            `json:"epic,omitempty"`
	Assignee     string            `json:"assignee,omitempty"`
	Reporter     string            `json:"reporter,omitempty"`
	Status       string            `json:"status,omitempty"`
	StoryPoints  float64           `json:"storyPoints,omitempty"`
	AcceptCrit   string            `json:"acceptanceCriteria,omitempty"`
	CustomFields map[string]string `json:"customFields,omitempty"`
	SubTasks     []JiraIssue       `json:"subtasks,omitempty"`
	RequireID    string            `json:"requirementId,omitempty"`
}

// JiraExport contains the full Jira export structure.
type JiraExport struct {
	Project    JiraProject `json:"project"`
	ExportedAt time.Time   `json:"exportedAt"`
	Source     string      `json:"source"`
	Issues     []JiraIssue `json:"issues"`
	Epics      []JiraIssue `json:"epics,omitempty"`
	Components []string    `json:"components,omitempty"`
	Versions   []string    `json:"versions,omitempty"`
}

// JiraProject contains project metadata.
type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Name returns the target name.
func (t *JiraTarget) Name() string {
	return "jira"
}

// Description returns a description of the target.
func (t *JiraTarget) Description() string {
	return "Export tasks as Jira issues (JSON for Jira API or import)"
}

// Capabilities returns what this target supports.
func (t *JiraTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: false,
		MultiAgent:        false,
		Verification:      false,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *JiraTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("empty spec content")
	}
	return nil
}

// Export exports the spec to this target.
func (t *JiraTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	result := &ExportResult{
		Target:    t.Name(),
		OutputDir: config.OutputDir,
		Files:     []string{},
	}

	export := t.Convert(spec, config)

	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	// Write main issues JSON
	issuesPath := filepath.Join(config.OutputDir, "jira-issues.json")
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling export: %w", err)
	}

	if err := os.WriteFile(issuesPath, data, 0600); err != nil {
		return nil, fmt.Errorf("writing issues file: %w", err)
	}
	result.Files = append(result.Files, issuesPath)

	// Write CSV for import
	csvPath := filepath.Join(config.OutputDir, "jira-import.csv")
	csvContent := export.GenerateCSV()
	if err := os.WriteFile(csvPath, []byte(csvContent), 0600); err != nil {
		return nil, fmt.Errorf("writing CSV file: %w", err)
	}
	result.Files = append(result.Files, csvPath)

	result.Success = true
	result.Message = fmt.Sprintf("Exported %d issues and %d epics", len(export.Issues), len(export.Epics))

	return result, nil
}

// Convert transforms spec content to JiraExport.
func (t *JiraTarget) Convert(spec string, config ExportConfig) *JiraExport {
	export := &JiraExport{
		ExportedAt: time.Now(),
		Issues:     []JiraIssue{},
		Epics:      []JiraIssue{},
	}

	if key, ok := config.Options["project_key"].(string); ok {
		export.Project.Key = key
	}
	if name, ok := config.Options["project_name"].(string); ok {
		export.Project.Name = name
	}

	export.Source = config.ProjectName

	// Extract issues
	epics, issues := t.extractIssues(spec)
	export.Epics = epics
	export.Issues = issues

	// Collect unique components and versions
	componentSet := make(map[string]bool)
	versionSet := make(map[string]bool)

	for _, issue := range append(epics, issues...) {
		for _, comp := range issue.Components {
			componentSet[comp] = true
		}
		for _, ver := range issue.FixVersions {
			versionSet[ver] = true
		}
	}

	for comp := range componentSet {
		export.Components = append(export.Components, comp)
	}
	for ver := range versionSet {
		export.Versions = append(export.Versions, ver)
	}

	return export
}

// extractIssues parses spec content and extracts Jira issues.
func (t *JiraTarget) extractIssues(content string) ([]JiraIssue, []JiraIssue) {
	var epics []JiraIssue
	var issues []JiraIssue

	lines := strings.Split(content, "\n")
	var currentSection string
	var currentEpic *JiraIssue
	var currentIssue *JiraIssue
	var descLines []string
	var acceptCrit []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Track top-level headers as Epics
		if strings.HasPrefix(trimmed, "## ") {
			if currentIssue != nil {
				currentIssue.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
				if len(acceptCrit) > 0 {
					currentIssue.AcceptCrit = strings.Join(acceptCrit, "\n")
				}
				issues = append(issues, *currentIssue)
				currentIssue = nil
				descLines = nil
				acceptCrit = nil
			}

			currentSection = strings.TrimPrefix(trimmed, "## ")

			epic := JiraIssue{
				Summary:   currentSection,
				IssueType: "Epic",
				Status:    "To Do",
			}
			currentEpic = &epic
			epics = append(epics, epic)
			continue
		}

		// Sub-headers as Stories
		if strings.HasPrefix(trimmed, "### ") {
			if currentIssue != nil {
				currentIssue.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
				if len(acceptCrit) > 0 {
					currentIssue.AcceptCrit = strings.Join(acceptCrit, "\n")
				}
				issues = append(issues, *currentIssue)
				descLines = nil
				acceptCrit = nil
			}

			header := strings.TrimPrefix(trimmed, "### ")

			issue := JiraIssue{
				Summary:   header,
				IssueType: "Story",
				Status:    "To Do",
				Labels:    []string{},
			}

			if currentEpic != nil {
				issue.Epic = currentEpic.Summary
			}

			if reqID := extractJiraReqID(header); reqID != "" {
				issue.RequireID = reqID
				issue.Labels = append(issue.Labels, reqID)
			}

			currentIssue = &issue
			continue
		}

		// Task items
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") {
			isCompleted := strings.HasPrefix(trimmed, "- [x] ")
			taskText := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- [ ] "), "- [x] ")

			task := JiraIssue{
				Summary:   taskText,
				IssueType: "Task",
				Status:    "To Do",
				Labels:    []string{},
			}

			if isCompleted {
				task.Status = "Done"
			}

			// Priority extraction
			if strings.Contains(strings.ToLower(taskText), "[critical]") {
				task.Priority = "Highest"
			} else if strings.Contains(strings.ToLower(taskText), "[high]") {
				task.Priority = "High"
			} else if strings.Contains(strings.ToLower(taskText), "[medium]") {
				task.Priority = "Medium"
			} else if strings.Contains(strings.ToLower(taskText), "[low]") {
				task.Priority = "Low"
			}

			if currentIssue != nil {
				currentIssue.SubTasks = append(currentIssue.SubTasks, task)
			} else {
				if currentEpic != nil {
					task.Epic = currentEpic.Summary
				}
				issues = append(issues, task)
			}
			continue
		}

		// Acceptance criteria
		if strings.HasPrefix(trimmed, "- Given ") ||
			strings.HasPrefix(trimmed, "- When ") ||
			strings.HasPrefix(trimmed, "- Then ") {
			if currentIssue != nil && trimmed != "" {
				acceptCrit = append(acceptCrit, trimmed)
			}
			continue
		}

		// Description lines
		if currentIssue != nil && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			descLines = append(descLines, trimmed)
		}
	}

	if currentIssue != nil {
		currentIssue.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
		if len(acceptCrit) > 0 {
			currentIssue.AcceptCrit = strings.Join(acceptCrit, "\n")
		}
		issues = append(issues, *currentIssue)
	}

	return epics, issues
}

// GenerateCSV generates CSV for Jira import.
func (e *JiraExport) GenerateCSV() string {
	var sb strings.Builder

	sb.WriteString("Summary,Description,Issue Type,Priority,Labels,Epic Link,Story Points\n")

	for _, epic := range e.Epics {
		sb.WriteString(jiraCSVEscape(epic.Summary) + ",")
		sb.WriteString(jiraCSVEscape(epic.Description) + ",")
		sb.WriteString("Epic,")
		sb.WriteString(",")
		sb.WriteString(jiraCSVEscape(strings.Join(epic.Labels, " ")) + ",")
		sb.WriteString(",")
		sb.WriteString("\n")
	}

	for _, issue := range e.Issues {
		sb.WriteString(jiraCSVEscape(issue.Summary) + ",")

		desc := issue.Description
		if issue.AcceptCrit != "" {
			desc += "\n\nAcceptance Criteria:\n" + issue.AcceptCrit
		}
		sb.WriteString(jiraCSVEscape(desc) + ",")

		sb.WriteString(issue.IssueType + ",")
		sb.WriteString(issue.Priority + ",")
		sb.WriteString(jiraCSVEscape(strings.Join(issue.Labels, " ")) + ",")
		sb.WriteString(jiraCSVEscape(issue.Epic) + ",")
		if issue.StoryPoints > 0 {
			sb.WriteString(fmt.Sprintf("%.0f", issue.StoryPoints))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Helper functions

func extractJiraReqID(s string) string {
	// Order matters: longer prefixes first to avoid partial matches (NFR- before FR-)
	patterns := []string{"STORY-", "NFR-", "REQ-", "US-", "FR-"}
	for _, prefix := range patterns {
		if idx := strings.Index(strings.ToUpper(s), prefix); idx >= 0 {
			end := idx + len(prefix)
			for end < len(s) && s[end] >= '0' && s[end] <= '9' {
				end++
			}
			return strings.ToUpper(s[idx:end])
		}
	}
	return ""
}

func jiraCSVEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		s = "\"" + s + "\""
	}
	return s
}
