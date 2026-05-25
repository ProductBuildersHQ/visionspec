package target

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	Register(&GasTownTarget{})
}

// GasTownTarget exports to GasTown TOML formula format.
// GasTown uses convoy, workflow, and expansion formulas with Bead definitions.
type GasTownTarget struct{}

// Name returns the target name.
func (t *GasTownTarget) Name() string {
	return "gastown"
}

// Description returns a description of the target.
func (t *GasTownTarget) Description() string {
	return "GasTown TOML formulas (convoy, workflow, expansion)"
}

// Capabilities returns what this target supports.
func (t *GasTownTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: true,
		MultiAgent:        true,
		Verification:      true,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *GasTownTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("spec content is empty")
	}
	return nil
}

// Export exports the spec to GasTown format.
func (t *GasTownTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	// Determine formula type
	formulaType := "workflow"
	if ft, ok := config.Options["formula_type"].(string); ok && ft != "" {
		formulaType = ft
	}

	// Determine rig (execution environment)
	rig := "default"
	if r, ok := config.Options["rig"].(string); ok && r != "" {
		rig = r
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = ".gastown"
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var files []string

	// Extract beads (tasks) from spec
	beads := t.extractBeads(spec)

	// Generate main formula file
	formula := t.generateFormula(config.ProjectName, formulaType, rig, beads)
	formulaPath := filepath.Join(outputDir, "formula.toml")
	if err := os.WriteFile(formulaPath, []byte(formula), 0600); err != nil {
		return nil, fmt.Errorf("writing formula.toml: %w", err)
	}
	files = append(files, formulaPath)

	// Generate beads directory with individual bead files
	beadsDir := filepath.Join(outputDir, "beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		return nil, fmt.Errorf("creating beads directory: %w", err)
	}

	for _, bead := range beads {
		beadContent := t.generateBead(bead)
		beadPath := filepath.Join(beadsDir, bead.ID+".toml")
		if err := os.WriteFile(beadPath, []byte(beadContent), 0600); err != nil {
			return nil, fmt.Errorf("writing bead %s: %w", bead.ID, err)
		}
		files = append(files, beadPath)
	}

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: outputDir,
		Files:     files,
		Success:   true,
		Message:   fmt.Sprintf("Exported to %s with %d beads", outputDir, len(beads)),
	}, nil
}

// Bead represents a GasTown bead (unit of work).
type Bead struct {
	ID          string
	Name        string
	Description string
	Status      string   // ready, blocked, done
	BlockedBy   []string // IDs of beads that must complete first
	Agent       string   // Which agent handles this
	Priority    int
}

// extractBeads extracts beads from the spec.
func (t *GasTownTarget) extractBeads(spec string) []Bead {
	var beads []Bead

	// Look for task lists
	taskRE := regexp.MustCompile(`(?m)^[-*]\s*\[[ x]\]\s*(.+)$`)
	matches := taskRE.FindAllStringSubmatch(spec, -1)

	for i, match := range matches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			status := "ready"
			if strings.Contains(match[0], "[x]") || strings.Contains(match[0], "[X]") {
				status = "done"
			}

			beads = append(beads, Bead{
				ID:          fmt.Sprintf("bead-%03d", i+1),
				Name:        truncate(desc, 50),
				Description: desc,
				Status:      status,
				Agent:       "default",
				Priority:    i + 1,
			})
		}
	}

	// Look for numbered tasks
	numberedRE := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	numberedMatches := numberedRE.FindAllStringSubmatch(spec, -1)
	startIdx := len(beads)
	for i, match := range numberedMatches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			if strings.HasPrefix(desc, "#") || len(desc) < 5 {
				continue
			}

			bead := Bead{
				ID:          fmt.Sprintf("bead-%03d", startIdx+i+1),
				Name:        truncate(desc, 50),
				Description: desc,
				Status:      "ready",
				Agent:       "default",
				Priority:    startIdx + i + 1,
			}

			// Simple dependency: each task depends on the previous one
			if len(beads) > 0 {
				bead.BlockedBy = []string{beads[len(beads)-1].ID}
				bead.Status = "blocked"
			}

			beads = append(beads, bead)
		}
	}

	// Default beads if none found
	if len(beads) == 0 {
		beads = []Bead{
			{ID: "bead-001", Name: "Setup", Description: "Project setup and configuration", Status: "ready", Agent: "setup", Priority: 1},
			{ID: "bead-002", Name: "Implement", Description: "Core implementation", Status: "blocked", Agent: "dev", Priority: 2, BlockedBy: []string{"bead-001"}},
			{ID: "bead-003", Name: "Test", Description: "Testing and validation", Status: "blocked", Agent: "qa", Priority: 3, BlockedBy: []string{"bead-002"}},
			{ID: "bead-004", Name: "Deploy", Description: "Deployment", Status: "blocked", Agent: "ops", Priority: 4, BlockedBy: []string{"bead-003"}},
		}
	}

	return beads
}

// generateFormula creates the main formula.toml file.
func (t *GasTownTarget) generateFormula(projectName, formulaType, rig string, beads []Bead) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# GasTown Formula\n")
	sb.WriteString(fmt.Sprintf("# Generated by visionspec on %s\n\n", time.Now().Format("2006-01-02")))

	// Formula metadata
	sb.WriteString("[formula]\n")
	sb.WriteString(fmt.Sprintf("name = %q\n", projectName))
	sb.WriteString(fmt.Sprintf("type = %q\n", formulaType))
	sb.WriteString(fmt.Sprintf("rig = %q\n", rig))
	sb.WriteString("version = \"1.0\"\n")
	sb.WriteString(fmt.Sprintf("created = %q\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("source = \"multispec\"\n\n")

	// Formula type specific configuration
	switch formulaType {
	case "convoy":
		sb.WriteString("[convoy]\n")
		sb.WriteString("# Parallel review configuration\n")
		sb.WriteString("parallel = true\n")
		sb.WriteString("quorum = \"all\"  # all, majority, any\n")
		sb.WriteString("timeout = \"24h\"\n\n")
	case "workflow":
		sb.WriteString("[workflow]\n")
		sb.WriteString("# Sequential execution configuration\n")
		sb.WriteString("parallel = false\n")
		sb.WriteString("checkpoint_on_complete = true\n")
		sb.WriteString("rollback_on_failure = true\n\n")
	case "expansion":
		sb.WriteString("[expansion]\n")
		sb.WriteString("# Template expansion configuration\n")
		sb.WriteString("template = \"default\"\n")
		sb.WriteString("variables = {}\n\n")
	}

	// Beads reference
	sb.WriteString("[beads]\n")
	sb.WriteString("# Bead definitions are in ./beads/*.toml\n")
	sb.WriteString("directory = \"beads\"\n")
	sb.WriteString(fmt.Sprintf("count = %d\n\n", len(beads)))

	// Bead list
	sb.WriteString("[[beads.list]]\n")
	for _, bead := range beads {
		sb.WriteString(fmt.Sprintf("# %s\n", bead.Name))
		sb.WriteString(fmt.Sprintf("[[beads.list]]\n"))
		sb.WriteString(fmt.Sprintf("id = %q\n", bead.ID))
		sb.WriteString(fmt.Sprintf("status = %q\n", bead.Status))
		if len(bead.BlockedBy) > 0 {
			sb.WriteString(fmt.Sprintf("blocked_by = [%s]\n", quoteStrings(bead.BlockedBy)))
		}
		sb.WriteString("\n")
	}

	// Execution order
	sb.WriteString("[execution]\n")
	sb.WriteString("# DAG-based execution order\n")
	readyBeads := []string{}
	for _, bead := range beads {
		if bead.Status == "ready" {
			readyBeads = append(readyBeads, bead.ID)
		}
	}
	sb.WriteString(fmt.Sprintf("entry_points = [%s]\n", quoteStrings(readyBeads)))

	return sb.String()
}

// generateBead creates a single bead TOML file.
func (t *GasTownTarget) generateBead(bead Bead) string {
	var sb strings.Builder

	sb.WriteString("# Bead Definition\n\n")

	sb.WriteString("[bead]\n")
	sb.WriteString(fmt.Sprintf("id = %q\n", bead.ID))
	sb.WriteString(fmt.Sprintf("name = %q\n", bead.Name))
	sb.WriteString(fmt.Sprintf("description = %q\n", bead.Description))
	sb.WriteString(fmt.Sprintf("status = %q\n", bead.Status))
	sb.WriteString(fmt.Sprintf("priority = %d\n\n", bead.Priority))

	sb.WriteString("[bead.agent]\n")
	sb.WriteString(fmt.Sprintf("type = %q\n", bead.Agent))
	sb.WriteString("capabilities = [\"execute\", \"report\"]\n\n")

	if len(bead.BlockedBy) > 0 {
		sb.WriteString("[bead.dependencies]\n")
		sb.WriteString(fmt.Sprintf("blocked_by = [%s]\n", quoteStrings(bead.BlockedBy)))
		sb.WriteString("wait_for_all = true\n\n")
	}

	sb.WriteString("[bead.execution]\n")
	sb.WriteString("timeout = \"1h\"\n")
	sb.WriteString("retries = 3\n")
	sb.WriteString("on_failure = \"pause\"  # pause, skip, abort\n")

	return sb.String()
}

// truncate shortens a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// quoteStrings formats a string slice for TOML.
func quoteStrings(ss []string) string {
	quoted := make([]string, len(ss))
	for i, s := range ss {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ", ")
}
