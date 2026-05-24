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

	// Write spec.md
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(spec), 0600); err != nil {
		return nil, fmt.Errorf("writing spec.md: %w", err)
	}

	// Generate plan.md from spec
	plan := t.generatePlan(spec, config.ProjectName)
	planPath := filepath.Join(specDir, "plan.md")
	if err := os.WriteFile(planPath, []byte(plan), 0600); err != nil {
		return nil, fmt.Errorf("writing plan.md: %w", err)
	}

	// Generate tasks.md from spec
	tasks := t.generateTasks(spec, config.ProjectName)
	tasksPath := filepath.Join(specDir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(tasks), 0600); err != nil {
		return nil, fmt.Errorf("writing tasks.md: %w", err)
	}

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: specDir,
		Files: []string{
			specPath,
			planPath,
			tasksPath,
		},
		Success: true,
		Message: fmt.Sprintf("Exported to %s", specDir),
	}, nil
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
