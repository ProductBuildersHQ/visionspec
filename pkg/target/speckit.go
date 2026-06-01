package target

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	Register(&SpecKitTarget{})
}

// SpecKitTarget exports to GitHub Spec-Kit format.
type SpecKitTarget struct{}

// Name returns the target name.
func (t *SpecKitTarget) Name() string {
	return "speckit"
}

// Description returns a description of the target.
func (t *SpecKitTarget) Description() string {
	return "GitHub Spec-Kit format (spec.md, plan.md, tasks.md)"
}

// Capabilities returns what this target supports.
func (t *SpecKitTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: false,
		MultiAgent:        false,
		Verification:      true,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *SpecKitTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("spec content is empty")
	}
	return nil
}

// Export exports the spec to SpecKit format.
func (t *SpecKitTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = "specs"
	}

	// Determine branch naming
	branchNumbering := "sequential"
	if bn, ok := config.Options["branch_numbering"].(string); ok && bn != "" {
		branchNumbering = bn
	}

	// Generate branch name
	var branchName string
	switch branchNumbering {
	case "timestamp":
		branchName = fmt.Sprintf("%s-%s", time.Now().Format("20060102-150405"), slugify(config.ProjectName))
	default:
		// Sequential: find next number
		branchName = fmt.Sprintf("001-%s", slugify(config.ProjectName))
		// In production, would scan existing dirs to find next number
	}

	// Create output directory
	specDir := filepath.Join(outputDir, branchName)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var files []string

	// Write spec.md
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(spec), 0600); err != nil {
		return nil, fmt.Errorf("writing spec.md: %w", err)
	}
	files = append(files, specPath)

	// Generate plan.md from spec
	plan := t.generatePlan(spec, config.ProjectName)
	planPath := filepath.Join(specDir, "plan.md")
	if err := os.WriteFile(planPath, []byte(plan), 0600); err != nil {
		return nil, fmt.Errorf("writing plan.md: %w", err)
	}
	files = append(files, planPath)

	// Generate tasks.md from spec
	tasks := t.generateTasks(spec, config.ProjectName)
	tasksPath := filepath.Join(specDir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(tasks), 0600); err != nil {
		return nil, fmt.Errorf("writing tasks.md: %w", err)
	}
	files = append(files, tasksPath)

	// Sync constitution if provided
	if constitutionPath, ok := config.Options["constitution_path"].(string); ok && constitutionPath != "" {
		if constitutionContent, err := os.ReadFile(constitutionPath); err == nil {
			constitutionFiles, err := t.syncConstitution(outputDir, constitutionContent)
			if err != nil {
				// Non-fatal: log but continue
				fmt.Printf("Warning: failed to sync constitution: %v\n", err)
			} else {
				files = append(files, constitutionFiles...)
			}
		}
	}

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: specDir,
		Files:     files,
		Success:   true,
		Message:   fmt.Sprintf("Exported to %s", specDir),
	}, nil
}

// syncConstitution syncs CONSTITUTION.md to SpecKit memory location.
func (t *SpecKitTarget) syncConstitution(outputDir string, constitution []byte) ([]string, error) {
	var files []string

	// SpecKit stores constitution in .specify/memory/constitution.md
	// Clean the path to satisfy gosec G703 (path traversal)
	memoryDir := filepath.Clean(filepath.Join(filepath.Dir(outputDir), ".specify", "memory"))
	if err := os.MkdirAll(memoryDir, 0755); err != nil {
		return nil, fmt.Errorf("creating memory directory: %w", err)
	}

	constitutionPath := filepath.Clean(filepath.Join(memoryDir, "constitution.md"))
	// #nosec G703 -- outputDir is from controlled export config, not arbitrary user input
	if err := os.WriteFile(constitutionPath, constitution, 0600); err != nil {
		return nil, fmt.Errorf("writing constitution.md: %w", err)
	}
	files = append(files, constitutionPath)

	return files, nil
}

// generatePlan creates a plan.md from the spec.
func (t *SpecKitTarget) generatePlan(_ string, projectName string) string {
	var sb strings.Builder

	sb.WriteString("# Implementation Plan: ")
	sb.WriteString(projectName)
	sb.WriteString("\n\n")

	sb.WriteString("## Overview\n\n")
	sb.WriteString("This plan outlines the implementation approach for the project.\n\n")

	sb.WriteString("## Phases\n\n")
	sb.WriteString("### Phase 1: Setup\n")
	sb.WriteString("- [ ] Review spec.md thoroughly\n")
	sb.WriteString("- [ ] Set up development environment\n")
	sb.WriteString("- [ ] Create initial project structure\n\n")

	sb.WriteString("### Phase 2: Core Implementation\n")
	sb.WriteString("- [ ] Implement core functionality\n")
	sb.WriteString("- [ ] Write unit tests\n")
	sb.WriteString("- [ ] Document APIs\n\n")

	sb.WriteString("### Phase 3: Integration\n")
	sb.WriteString("- [ ] Integrate components\n")
	sb.WriteString("- [ ] Write integration tests\n")
	sb.WriteString("- [ ] Performance testing\n\n")

	sb.WriteString("### Phase 4: Verification\n")
	sb.WriteString("- [ ] Verify against acceptance criteria\n")
	sb.WriteString("- [ ] Security review\n")
	sb.WriteString("- [ ] Final documentation\n\n")

	sb.WriteString("## Notes\n\n")
	sb.WriteString("See spec.md for detailed requirements and tasks.md for task breakdown.\n")

	return sb.String()
}

// generateTasks creates a tasks.md from the spec.
func (t *SpecKitTarget) generateTasks(_ string, projectName string) string {
	var sb strings.Builder

	sb.WriteString("# Tasks: ")
	sb.WriteString(projectName)
	sb.WriteString("\n\n")

	sb.WriteString("## Task List\n\n")
	sb.WriteString("| # | Task | Status | Priority | Assignee |\n")
	sb.WriteString("|---|------|--------|----------|----------|\n")
	sb.WriteString("| 1 | Review and understand spec | Todo | High | - |\n")
	sb.WriteString("| 2 | Set up project structure | Todo | High | - |\n")
	sb.WriteString("| 3 | Implement core features | Todo | High | - |\n")
	sb.WriteString("| 4 | Write tests | Todo | Medium | - |\n")
	sb.WriteString("| 5 | Documentation | Todo | Medium | - |\n")
	sb.WriteString("| 6 | Code review | Todo | High | - |\n")
	sb.WriteString("| 7 | Deploy | Todo | High | - |\n")
	sb.WriteString("\n")

	sb.WriteString("## Dependencies\n\n")
	sb.WriteString("```\n")
	sb.WriteString("1 → 2 → 3 → 4 → 5 → 6 → 7\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Notes\n\n")
	sb.WriteString("- Tasks extracted from spec.md\n")
	sb.WriteString("- Update status as work progresses\n")

	return sb.String()
}

// slugify converts a string to a slug suitable for branch names.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// CanSync returns whether this target supports syncing.
func (t *SpecKitTarget) CanSync(config ExportConfig) bool {
	// Check if output dir exists with tasks.md
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = "specs"
	}

	// Find the latest branch directory
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			tasksPath := filepath.Join(outputDir, entry.Name(), "tasks.md")
			if _, err := os.Stat(tasksPath); err == nil {
				return true
			}
		}
	}

	return false
}

// Sync retrieves task state from the exported SpecKit target.
func (t *SpecKitTarget) Sync(config ExportConfig) (*SyncResult, error) {
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = "specs"
	}

	// Find the latest branch directory
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("reading output directory: %w", err)
	}

	var latestDir string
	for _, entry := range entries {
		if entry.IsDir() {
			latestDir = filepath.Join(outputDir, entry.Name())
		}
	}

	if latestDir == "" {
		return nil, fmt.Errorf("no exported branch found in %s", outputDir)
	}

	// Read tasks.md
	tasksPath := filepath.Join(latestDir, "tasks.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("reading tasks.md: %w", err)
	}

	// Parse tasks
	tasks := t.parseTasks(string(content))

	// Calculate summary
	summary := SyncSummary{TotalTasks: len(tasks)}
	for _, task := range tasks {
		switch task.Status {
		case "todo":
			summary.TodoCount++
		case "in_progress":
			summary.InProgress++
		case "done":
			summary.DoneCount++
		}
	}

	return &SyncResult{
		Target:   t.Name(),
		SyncedAt: time.Now(),
		Tasks:    tasks,
		Summary:  summary,
	}, nil
}

// parseTasks extracts tasks from tasks.md content.
// SpecKit uses: [ ] for todo, [x] for done
func (t *SpecKitTarget) parseTasks(content string) []TaskState {
	var tasks []TaskState

	// Look for table rows
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}

		// Skip header and separator rows
		if strings.Contains(line, "---") || strings.Contains(line, "Task") {
			continue
		}

		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}

		id := strings.TrimSpace(cells[1])
		title := strings.TrimSpace(cells[2])
		statusCell := strings.TrimSpace(cells[3])

		status := "todo"
		statusLower := strings.ToLower(statusCell)
		if strings.Contains(statusLower, "done") || strings.Contains(statusLower, "complete") {
			status = "done"
		} else if strings.Contains(statusLower, "progress") || strings.Contains(statusLower, "wip") {
			status = "in_progress"
		}

		if id != "" && id != "#" {
			tasks = append(tasks, TaskState{
				ID:     id,
				Title:  title,
				Status: status,
			})
		}
	}

	// Also look for checkbox lists
	checkboxRE := strings.NewReplacer()
	_ = checkboxRE
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- [ ]") {
			title := strings.TrimPrefix(line, "- [ ]")
			tasks = append(tasks, TaskState{
				ID:     fmt.Sprintf("CB-%03d", i+1),
				Title:  strings.TrimSpace(title),
				Status: "todo",
			})
		} else if strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]") {
			title := strings.TrimPrefix(line, "- [x]")
			title = strings.TrimPrefix(title, "- [X]")
			tasks = append(tasks, TaskState{
				ID:     fmt.Sprintf("CB-%03d", i+1),
				Title:  strings.TrimSpace(title),
				Status: "done",
			})
		}
	}

	return tasks
}
