package target

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	Register(&GSDTarget{})
}

// GSDTarget exports to GSD (Get Shit Done) format.
// GSD uses PLAN.md with YAML frontmatter, STATE.md for tracking,
// and .planning/config.json for configuration.
type GSDTarget struct{}

// Name returns the target name.
func (t *GSDTarget) Name() string {
	return "gsd"
}

// Description returns a description of the target.
func (t *GSDTarget) Description() string {
	return "GSD planning format (PLAN.md, STATE.md, config.json)"
}

// Capabilities returns what this target supports.
func (t *GSDTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: true,
		MultiAgent:        false,
		Verification:      true,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *GSDTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("spec content is empty")
	}
	return nil
}

// Export exports the spec to GSD format.
func (t *GSDTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = ".planning"
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var files []string

	// Extract information from spec
	mustHaves := t.extractMustHaves(spec)
	phases := t.extractPhases(spec)
	tasks := t.extractTasks(spec)

	// Generate PLAN.md
	plan := t.generatePlan(spec, config.ProjectName, mustHaves, phases, tasks)
	planPath := filepath.Join(outputDir, "PLAN.md")
	if err := os.WriteFile(planPath, []byte(plan), 0600); err != nil {
		return nil, fmt.Errorf("writing PLAN.md: %w", err)
	}
	files = append(files, planPath)

	// Generate STATE.md
	state := t.generateState(config.ProjectName, phases, tasks)
	statePath := filepath.Join(outputDir, "STATE.md")
	if err := os.WriteFile(statePath, []byte(state), 0600); err != nil {
		return nil, fmt.Errorf("writing STATE.md: %w", err)
	}
	files = append(files, statePath)

	// Generate config.json
	configJSON := t.generateConfig(config)
	configPath := filepath.Join(outputDir, "config.json")
	if err := os.WriteFile(configPath, configJSON, 0600); err != nil {
		return nil, fmt.Errorf("writing config.json: %w", err)
	}
	files = append(files, configPath)

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: outputDir,
		Files:     files,
		Success:   true,
		Message:   fmt.Sprintf("Exported to %s", outputDir),
	}, nil
}

// MustHaves represents the GSD must_haves structure.
type MustHaves struct {
	Truths    []string `json:"truths"`
	Artifacts []string `json:"artifacts"`
	KeyLinks  []string `json:"key_links"`
}

// Phase represents a GSD phase.
type Phase struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tasks       []string `json:"tasks"`
	Wave        int      `json:"wave"`
}

// Task represents a GSD task.
type Task struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Status      string   `json:"status"` // todo, in_progress, done
	Priority    string   `json:"priority"`
	DependsOn   []string `json:"depends_on,omitempty"`
}

// extractMustHaves extracts must_haves from the spec.
func (t *GSDTarget) extractMustHaves(spec string) MustHaves {
	mustHaves := MustHaves{
		Truths:    []string{},
		Artifacts: []string{},
		KeyLinks:  []string{},
	}

	// Extract acceptance criteria as truths
	criteriaRE := regexp.MustCompile(`(?i)(?:acceptance criteria|success criteria|requirements?)[\s:]*\n((?:[-*]\s+.+\n?)+)`)
	if matches := criteriaRE.FindStringSubmatch(spec); len(matches) > 1 {
		lines := strings.Split(matches[1], "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				truth := strings.TrimLeft(line, "-* ")
				if truth != "" {
					mustHaves.Truths = append(mustHaves.Truths, truth)
				}
			}
		}
	}

	// Extract deliverables as artifacts
	deliverableRE := regexp.MustCompile(`(?i)(?:deliverables?|outputs?|artifacts?)[\s:]*\n((?:[-*]\s+.+\n?)+)`)
	if matches := deliverableRE.FindStringSubmatch(spec); len(matches) > 1 {
		lines := strings.Split(matches[1], "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				artifact := strings.TrimLeft(line, "-* ")
				if artifact != "" {
					mustHaves.Artifacts = append(mustHaves.Artifacts, artifact)
				}
			}
		}
	}

	// Extract dependencies as key_links
	depsRE := regexp.MustCompile(`(?i)(?:dependencies|prerequisites|depends on)[\s:]*\n((?:[-*]\s+.+\n?)+)`)
	if matches := depsRE.FindStringSubmatch(spec); len(matches) > 1 {
		lines := strings.Split(matches[1], "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				link := strings.TrimLeft(line, "-* ")
				if link != "" {
					mustHaves.KeyLinks = append(mustHaves.KeyLinks, link)
				}
			}
		}
	}

	// Default truths if none found
	if len(mustHaves.Truths) == 0 {
		mustHaves.Truths = []string{"All requirements implemented", "Tests passing", "Documentation complete"}
	}

	return mustHaves
}

// extractPhases extracts phases from the spec.
func (t *GSDTarget) extractPhases(spec string) []Phase {
	var phases []Phase

	// Look for phase/section headers
	phaseRE := regexp.MustCompile(`(?m)^##\s*(?:Phase\s*\d+[:\s]*)?(.+)$`)
	matches := phaseRE.FindAllStringSubmatch(spec, -1)

	wave := 1
	for _, match := range matches {
		if len(match) > 1 {
			name := strings.TrimSpace(match[1])
			// Skip certain headers
			if strings.EqualFold(name, "Overview") || strings.EqualFold(name, "Introduction") ||
				strings.EqualFold(name, "Decision Log") || strings.EqualFold(name, "Traceability Matrix") {
				continue
			}
			phases = append(phases, Phase{
				Name:        name,
				Description: fmt.Sprintf("Implement %s", name),
				Wave:        wave,
			})
			wave++
		}
	}

	// Default phases if none found
	if len(phases) == 0 {
		phases = []Phase{
			{Name: "Setup", Description: "Project setup and configuration", Wave: 1},
			{Name: "Implementation", Description: "Core implementation", Wave: 2},
			{Name: "Testing", Description: "Testing and validation", Wave: 3},
			{Name: "Deployment", Description: "Deployment and release", Wave: 4},
		}
	}

	return phases
}

// extractTasks extracts tasks from the spec.
func (t *GSDTarget) extractTasks(spec string) []Task {
	var tasks []Task

	// Look for task lists
	taskRE := regexp.MustCompile(`(?m)^[-*]\s*\[[ x]\]\s*(.+)$`)
	matches := taskRE.FindAllStringSubmatch(spec, -1)

	for i, match := range matches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			status := "todo"
			if strings.Contains(match[0], "[x]") || strings.Contains(match[0], "[X]") {
				status = "done"
			}

			tasks = append(tasks, Task{
				ID:          fmt.Sprintf("TASK-%03d", i+1),
				Description: desc,
				Status:      status,
				Priority:    "medium",
			})
		}
	}

	// Look for numbered tasks
	numberedRE := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	numberedMatches := numberedRE.FindAllStringSubmatch(spec, -1)
	for i, match := range numberedMatches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			// Skip if it looks like a header or section
			if strings.HasPrefix(desc, "#") || len(desc) < 5 {
				continue
			}
			tasks = append(tasks, Task{
				ID:          fmt.Sprintf("TASK-%03d", len(tasks)+i+1),
				Description: desc,
				Status:      "todo",
				Priority:    "medium",
			})
		}
	}

	// Default tasks if none found
	if len(tasks) == 0 {
		tasks = []Task{
			{ID: "TASK-001", Description: "Review specification", Status: "todo", Priority: "high"},
			{ID: "TASK-002", Description: "Set up development environment", Status: "todo", Priority: "high"},
			{ID: "TASK-003", Description: "Implement core functionality", Status: "todo", Priority: "high", DependsOn: []string{"TASK-002"}},
			{ID: "TASK-004", Description: "Write unit tests", Status: "todo", Priority: "medium", DependsOn: []string{"TASK-003"}},
			{ID: "TASK-005", Description: "Integration testing", Status: "todo", Priority: "medium", DependsOn: []string{"TASK-004"}},
			{ID: "TASK-006", Description: "Documentation", Status: "todo", Priority: "medium"},
			{ID: "TASK-007", Description: "Deploy to production", Status: "todo", Priority: "high", DependsOn: []string{"TASK-005"}},
		}
	}

	return tasks
}

// generatePlan creates a PLAN.md in GSD format.
func (t *GSDTarget) generatePlan(_, projectName string, mustHaves MustHaves, phases []Phase, tasks []Task) string {
	var sb strings.Builder

	// YAML frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", projectName))
	sb.WriteString(fmt.Sprintf("created: %s\n", time.Now().Format("2006-01-02")))
	sb.WriteString("status: planning\n")
	sb.WriteString("must_haves:\n")
	sb.WriteString("  truths:\n")
	for _, truth := range mustHaves.Truths {
		sb.WriteString(fmt.Sprintf("    - %q\n", truth))
	}
	sb.WriteString("  artifacts:\n")
	for _, artifact := range mustHaves.Artifacts {
		sb.WriteString(fmt.Sprintf("    - %q\n", artifact))
	}
	sb.WriteString("  key_links:\n")
	for _, link := range mustHaves.KeyLinks {
		sb.WriteString(fmt.Sprintf("    - %q\n", link))
	}
	sb.WriteString("---\n\n")

	// Main content
	sb.WriteString("# ")
	sb.WriteString(projectName)
	sb.WriteString("\n\n")

	sb.WriteString("## Overview\n\n")
	sb.WriteString("This plan was generated from a visionspec reconciled specification.\n\n")

	// Phases
	sb.WriteString("## Phases\n\n")
	for _, phase := range phases {
		sb.WriteString(fmt.Sprintf("### Wave %d: %s\n\n", phase.Wave, phase.Name))
		sb.WriteString(phase.Description)
		sb.WriteString("\n\n")
	}

	// Tasks in XML format (GSD convention)
	sb.WriteString("## Tasks\n\n")
	sb.WriteString("```xml\n")
	sb.WriteString("<tasks>\n")
	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("  <task id=%q status=%q priority=%q", task.ID, task.Status, task.Priority))
		if len(task.DependsOn) > 0 {
			sb.WriteString(fmt.Sprintf(" depends_on=%q", strings.Join(task.DependsOn, ",")))
		}
		sb.WriteString(">\n")
		sb.WriteString(fmt.Sprintf("    %s\n", task.Description))
		sb.WriteString("  </task>\n")
	}
	sb.WriteString("</tasks>\n")
	sb.WriteString("```\n\n")

	// Include original spec reference
	sb.WriteString("## Source Specification\n\n")
	sb.WriteString("See the original spec.md for detailed requirements.\n")

	return sb.String()
}

// generateState creates a STATE.md for tracking progress.
func (t *GSDTarget) generateState(projectName string, _ []Phase, tasks []Task) string {
	var sb strings.Builder

	sb.WriteString("# State: ")
	sb.WriteString(projectName)
	sb.WriteString("\n\n")

	sb.WriteString("## Current Status\n\n")
	sb.WriteString("- **Phase**: Wave 1\n")
	sb.WriteString("- **Progress**: 0%\n")
	sb.WriteString(fmt.Sprintf("- **Last Updated**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	// Task summary
	todo, inProgress, done := 0, 0, 0
	for _, task := range tasks {
		switch task.Status {
		case "todo":
			todo++
		case "in_progress":
			inProgress++
		case "done":
			done++
		}
	}

	sb.WriteString("## Task Summary\n\n")
	sb.WriteString(fmt.Sprintf("| Status | Count |\n"))
	sb.WriteString("| ------ | ----- |\n")
	sb.WriteString(fmt.Sprintf("| Todo | %d |\n", todo))
	sb.WriteString(fmt.Sprintf("| In Progress | %d |\n", inProgress))
	sb.WriteString(fmt.Sprintf("| Done | %d |\n\n", done))

	// Active tasks
	sb.WriteString("## Active Tasks\n\n")
	hasActive := false
	for _, task := range tasks {
		if task.Status == "in_progress" {
			sb.WriteString(fmt.Sprintf("- [~] %s: %s\n", task.ID, task.Description))
			hasActive = true
		}
	}
	if !hasActive {
		sb.WriteString("No tasks in progress.\n")
	}
	sb.WriteString("\n")

	// Next up
	sb.WriteString("## Next Up\n\n")
	count := 0
	for _, task := range tasks {
		if task.Status == "todo" && count < 3 {
			sb.WriteString(fmt.Sprintf("- [ ] %s: %s\n", task.ID, task.Description))
			count++
		}
	}
	if count == 0 {
		sb.WriteString("All tasks complete!\n")
	}
	sb.WriteString("\n")

	// Blockers
	sb.WriteString("## Blockers\n\n")
	sb.WriteString("None identified.\n\n")

	// Notes
	sb.WriteString("## Notes\n\n")
	sb.WriteString("- State file auto-generated by visionspec\n")
	sb.WriteString("- Update this file as work progresses\n")

	return sb.String()
}

// generateConfig creates the config.json file.
func (t *GSDTarget) generateConfig(config ExportConfig) []byte {
	modelProfile := "default"
	if mp, ok := config.Options["model_profile"].(string); ok && mp != "" {
		modelProfile = mp
	}

	cfg := map[string]any{
		"project": config.ProjectName,
		"version": "1.0",
		"gsd": map[string]any{
			"model_profile":    modelProfile,
			"auto_update":      true,
			"parallel_waves":   false,
			"verification":     true,
			"checkpoint_waves": true,
		},
		"paths": map[string]string{
			"plan":  "PLAN.md",
			"state": "STATE.md",
		},
		"generated_at": time.Now().Format(time.RFC3339),
		"source":       "visionspec",
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	return data
}
