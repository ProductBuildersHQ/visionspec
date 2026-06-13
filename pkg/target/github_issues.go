// Package target provides a GitHub Issues export adapter.
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
	Register(&GitHubIssuesTarget{})
}

// GitHubIssuesTarget exports specs as GitHub Issues.
type GitHubIssuesTarget struct{}

// GitHubIssue represents a GitHub issue structure.
type GitHubIssue struct {
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Labels     []string `json:"labels,omitempty"`
	Assignees  []string `json:"assignees,omitempty"`
	Milestone  string   `json:"milestone,omitempty"`
	Project    string   `json:"project,omitempty"`
	State      string   `json:"state,omitempty"`
	Priority   string   `json:"priority,omitempty"`
	TaskID     string   `json:"task_id,omitempty"`
	RequireID  string   `json:"requirement_id,omitempty"`
	AcceptCrit []string `json:"acceptance_criteria,omitempty"`
}

// GitHubIssuesExport contains the full export structure.
type GitHubIssuesExport struct {
	Repository string        `json:"repository"`
	ExportedAt time.Time     `json:"exported_at"`
	Source     string        `json:"source"`
	Issues     []GitHubIssue `json:"issues"`
	Milestones []string      `json:"milestones,omitempty"`
	Labels     []string      `json:"labels,omitempty"`
}

// Name returns the target name.
func (t *GitHubIssuesTarget) Name() string {
	return "github-issues"
}

// Description returns a description of the target.
func (t *GitHubIssuesTarget) Description() string {
	return "Export tasks as GitHub Issues (JSON for gh CLI or API)"
}

// Capabilities returns what this target supports.
func (t *GitHubIssuesTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: false,
		MultiAgent:        false,
		Verification:      false,
		DependencyGraph:   false,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *GitHubIssuesTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("empty spec content")
	}
	return nil
}

// Export exports the spec to this target.
func (t *GitHubIssuesTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
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
	issuesPath := filepath.Join(config.OutputDir, "github-issues.json")
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling export: %w", err)
	}

	if err := os.WriteFile(issuesPath, data, 0600); err != nil {
		return nil, fmt.Errorf("writing issues file: %w", err)
	}
	result.Files = append(result.Files, issuesPath)

	// Write CLI script
	cliPath := filepath.Join(config.OutputDir, "create-issues.sh")
	cliScript := export.GenerateCLICommands()
	if err := os.WriteFile(cliPath, []byte(cliScript), 0600); err != nil {
		return nil, fmt.Errorf("writing CLI script: %w", err)
	}
	result.Files = append(result.Files, cliPath)

	result.Success = true
	result.Message = fmt.Sprintf("Exported %d issues", len(export.Issues))

	return result, nil
}

// Convert transforms spec content to GitHubIssuesExport.
func (t *GitHubIssuesTarget) Convert(spec string, config ExportConfig) *GitHubIssuesExport {
	export := &GitHubIssuesExport{
		ExportedAt: time.Now(),
		Issues:     []GitHubIssue{},
	}

	if repo, ok := config.Options["repository"].(string); ok {
		export.Repository = repo
	}

	export.Source = config.ProjectName

	// Extract issues from spec content
	issues := t.extractIssues(spec, config)
	export.Issues = issues

	// Collect unique labels and milestones
	labelSet := make(map[string]bool)
	milestoneSet := make(map[string]bool)

	for _, issue := range issues {
		for _, label := range issue.Labels {
			labelSet[label] = true
		}
		if issue.Milestone != "" {
			milestoneSet[issue.Milestone] = true
		}
	}

	for label := range labelSet {
		export.Labels = append(export.Labels, label)
	}
	for ms := range milestoneSet {
		export.Milestones = append(export.Milestones, ms)
	}

	return export
}

// extractIssues parses spec content and extracts GitHub issues.
func (t *GitHubIssuesTarget) extractIssues(content string, _ ExportConfig) []GitHubIssue {
	var issues []GitHubIssue

	lines := strings.Split(content, "\n")
	var currentSection string
	var currentIssue *GitHubIssue
	var bodyLines []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Track section headers
		if strings.HasPrefix(trimmed, "## ") {
			if currentIssue != nil {
				currentIssue.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
				issues = append(issues, *currentIssue)
				currentIssue = nil
				bodyLines = nil
			}
			currentSection = strings.TrimPrefix(trimmed, "## ")
			continue
		}

		// Look for task items
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") {
			if currentIssue != nil {
				currentIssue.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
				issues = append(issues, *currentIssue)
				bodyLines = nil
			}

			isCompleted := strings.HasPrefix(trimmed, "- [x] ")
			taskText := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- [ ] "), "- [x] ")

			issue := GitHubIssue{
				Title:  taskText,
				Labels: []string{},
				State:  "open",
			}

			if isCompleted {
				issue.State = "closed"
			}

			if currentSection != "" {
				issue.Labels = append(issue.Labels, ghSlugify(currentSection))
			}

			// Extract priority
			if strings.Contains(strings.ToLower(taskText), "[critical]") {
				issue.Priority = "critical"
				issue.Labels = append(issue.Labels, "priority:critical")
			} else if strings.Contains(strings.ToLower(taskText), "[high]") {
				issue.Priority = "high"
				issue.Labels = append(issue.Labels, "priority:high")
			}

			currentIssue = &issue
			continue
		}

		// Look for requirement headers
		if strings.HasPrefix(trimmed, "### ") {
			headerText := strings.TrimPrefix(trimmed, "### ")

			if isGHRequirementHeader(headerText) {
				if currentIssue != nil {
					currentIssue.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
					issues = append(issues, *currentIssue)
					bodyLines = nil
				}

				issue := GitHubIssue{
					Title:     headerText,
					Labels:    []string{"feature"},
					State:     "open",
					RequireID: extractGHRequirementID(headerText),
				}

				if currentSection != "" {
					issue.Labels = append(issue.Labels, ghSlugify(currentSection))
				}

				currentIssue = &issue
				continue
			}
		}

		// Collect body/acceptance criteria
		if currentIssue != nil && trimmed != "" {
			if strings.HasPrefix(trimmed, "- Given ") ||
				strings.HasPrefix(trimmed, "- When ") ||
				strings.HasPrefix(trimmed, "- Then ") {
				currentIssue.AcceptCrit = append(currentIssue.AcceptCrit, trimmed)
			} else {
				bodyLines = append(bodyLines, trimmed)
			}
		}
	}

	if currentIssue != nil {
		currentIssue.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
		issues = append(issues, *currentIssue)
	}

	return issues
}

// GenerateCLICommands generates gh CLI commands to create issues.
func (e *GitHubIssuesExport) GenerateCLICommands() string {
	var sb strings.Builder

	sb.WriteString("#!/bin/bash\n")
	sb.WriteString("# Generated GitHub Issues CLI commands\n")
	sb.WriteString(fmt.Sprintf("# Repository: %s\n", e.Repository))
	sb.WriteString(fmt.Sprintf("# Generated: %s\n\n", e.ExportedAt.Format("2006-01-02 15:04:05")))

	if len(e.Labels) > 0 {
		sb.WriteString("# Create labels\n")
		for _, label := range e.Labels {
			sb.WriteString(fmt.Sprintf("gh label create \"%s\" --force || true\n", escapeGHShell(label)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("# Create issues\n")
	for _, issue := range e.Issues {
		if issue.State == "closed" {
			continue
		}

		cmd := fmt.Sprintf("gh issue create --title \"%s\"", escapeGHShell(issue.Title))

		if issue.Body != "" {
			cmd += fmt.Sprintf(" --body \"%s\"", escapeGHShell(issue.Body))
		}

		for _, label := range issue.Labels {
			cmd += fmt.Sprintf(" --label \"%s\"", escapeGHShell(label))
		}

		sb.WriteString(cmd + "\n")
	}

	return sb.String()
}

// Helper functions

func ghSlugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func isGHRequirementHeader(s string) bool {
	// Order matters: longer prefixes first to avoid partial matches (NFR- before FR-)
	prefixes := []string{"STORY-", "NFR-", "FEAT-", "REQ-", "US-", "FR-"}
	upperS := strings.ToUpper(s)
	for _, prefix := range prefixes {
		if strings.HasPrefix(upperS, prefix) {
			return true
		}
	}
	return false
}

func extractGHRequirementID(s string) string {
	// Order matters: longer prefixes first to avoid partial matches (NFR- before FR-)
	patterns := []string{"STORY-", "NFR-", "FEAT-", "REQ-", "US-", "FR-"}
	for _, prefix := range patterns {
		if idx := strings.Index(strings.ToUpper(s), prefix); idx >= 0 {
			end := idx + len(prefix)
			for end < len(s) && (s[end] >= '0' && s[end] <= '9') {
				end++
			}
			return strings.ToUpper(s[idx:end])
		}
	}
	return ""
}

func escapeGHShell(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "$", "\\$")
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
