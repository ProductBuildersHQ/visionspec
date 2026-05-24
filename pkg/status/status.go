// Package status generates project status reports.
package status

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/plexusone/multispec/pkg/config"
	"github.com/plexusone/multispec/pkg/types"
)

// Report represents a project status report.
type Report struct {
	Project     string                `json:"project"`
	Path        string                `json:"path"`
	GeneratedAt time.Time             `json:"generated_at"`
	Readiness   types.ReadinessStatus `json:"readiness"`
	Specs       []SpecStatus          `json:"specs"`
	Summary     Summary               `json:"summary"`
}

// SpecStatus represents the status of a single spec.
type SpecStatus struct {
	Type       types.SpecType     `json:"type"`
	Category   types.SpecCategory `json:"category"`
	Filename   string             `json:"filename"`
	Exists     bool               `json:"exists"`
	Required   bool               `json:"required"`
	Status     types.SpecStatus   `json:"status"`
	EvalStatus *EvalStatus        `json:"eval_status,omitempty"`
	Approval   *types.Approval    `json:"approval,omitempty"`
}

// EvalStatus represents evaluation results.
type EvalStatus struct {
	Exists   bool   `json:"exists"`
	Decision string `json:"decision"` // pass, conditional, fail
	Findings struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
	} `json:"findings"`
}

// Summary provides aggregate statistics.
type Summary struct {
	TotalSpecs     int `json:"total_specs"`
	PresentSpecs   int `json:"present_specs"`
	EvaluatedSpecs int `json:"evaluated_specs"`
	ApprovedSpecs  int `json:"approved_specs"`
	BlockingIssues int `json:"blocking_issues"`
}

// Generate generates a status report for a project.
func Generate(project *types.Project) (*Report, error) {
	return GenerateWithConfig(project, project.GetSpecConfig())
}

// GenerateWithConfig generates a status report using a custom SpecConfig.
func GenerateWithConfig(project *types.Project, specConfig *types.SpecConfig) (*Report, error) {
	if specConfig == nil {
		specConfig = types.DefaultSpecConfig()
	}

	report := &Report{
		Project:     project.Name,
		Path:        project.Path,
		GeneratedAt: time.Now(),
	}

	// Check each spec type from the config
	for _, specName := range specConfig.AllSpecs() {
		specType := types.SpecType(specName)
		specPath := config.SpecPath(project.Path, specType)
		evalPath := config.EvalPath(project.Path, specType)

		ss := SpecStatus{
			Type:     specType,
			Category: specConfig.GetCategory(specName),
			Filename: specType.Filename(),
			Required: specConfig.IsRequired(specName),
		}

		// Check if spec exists
		if _, err := os.Stat(specPath); err == nil {
			ss.Exists = true
			ss.Status = types.StatusDraft
			report.Summary.PresentSpecs++
		} else {
			ss.Status = types.StatusMissing
		}

		// Check if eval exists
		if _, err := os.Stat(evalPath); err == nil {
			ss.EvalStatus = &EvalStatus{Exists: true}
			ss.Status = types.StatusEvaluated
			report.Summary.EvaluatedSpecs++
			// TODO: Load eval file and parse findings
		}

		// Check approval
		if project.Approvals != nil {
			if approval, ok := project.Approvals[specType]; ok {
				ss.Approval = approval
				ss.Status = types.StatusApproved
				report.Summary.ApprovedSpecs++
			}
		}

		report.Specs = append(report.Specs, ss)
		report.Summary.TotalSpecs++
	}

	// Calculate readiness
	report.Readiness = calculateReadiness(report)

	return report, nil
}

func calculateReadiness(report *Report) types.ReadinessStatus {
	status := types.ReadinessStatus{
		Ready: true,
	}

	// Gate 1: All required source specs present
	requiredPresent := true
	for _, spec := range report.Specs {
		if spec.Required && !spec.Exists {
			requiredPresent = false
			break
		}
	}
	status.Gates = append(status.Gates, types.ReadinessGate{
		Name:    "Required specs present",
		Passed:  requiredPresent,
		Message: gateMessage(requiredPresent, "All required specs exist", "Missing required specs"),
	})
	if !requiredPresent {
		status.Ready = false
	}

	// Gate 2: All evals passing
	evalsPass := true
	for _, spec := range report.Specs {
		if spec.EvalStatus != nil && spec.EvalStatus.Decision == "fail" {
			evalsPass = false
			break
		}
	}
	status.Gates = append(status.Gates, types.ReadinessGate{
		Name:    "Evaluations passing",
		Passed:  evalsPass,
		Message: gateMessage(evalsPass, "No blocking eval findings", "Blocking eval findings exist"),
	})
	if !evalsPass {
		status.Ready = false
	}

	// Gate 3: Required approvals obtained
	approvalsObtained := true
	for _, spec := range report.Specs {
		if spec.Required && spec.Exists && spec.Approval == nil {
			approvalsObtained = false
			break
		}
	}
	status.Gates = append(status.Gates, types.ReadinessGate{
		Name:    "Approvals obtained",
		Passed:  approvalsObtained,
		Message: gateMessage(approvalsObtained, "All required specs approved", "Pending approvals"),
	})
	if !approvalsObtained {
		status.Ready = false
	}

	// Gate 4: spec.md exists
	specExists := false
	specPath := filepath.Join(report.Path, "spec.md")
	if _, err := os.Stat(specPath); err == nil {
		specExists = true
	}
	status.Gates = append(status.Gates, types.ReadinessGate{
		Name:    "Execution spec generated",
		Passed:  specExists,
		Message: gateMessage(specExists, "spec.md exists", "spec.md not generated"),
	})
	if !specExists {
		status.Ready = false
	}

	// Summary
	if status.Ready {
		status.Summary = "Ready for AI-assisted development"
	} else {
		failedCount := 0
		for _, gate := range status.Gates {
			if !gate.Passed {
				failedCount++
			}
		}
		status.Summary = "Not ready: " + pluralize(failedCount, "blocker", "blockers")
	}

	return status
}

func gateMessage(passed bool, passMsg, failMsg string) string {
	if passed {
		return passMsg
	}
	return failMsg
}

func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", n, plural)
}

// RenderText renders the report as terminal text.
func RenderText(w io.Writer, report *Report) error {
	// Header
	fmt.Fprintf(w, "Project: %s\n", report.Project)
	fmt.Fprintf(w, "Path: %s\n", report.Path)
	fmt.Fprintf(w, "\n")

	// Readiness summary
	if report.Readiness.Ready {
		fmt.Fprintf(w, "Status: READY\n")
	} else {
		fmt.Fprintf(w, "Status: NOT READY\n")
	}
	fmt.Fprintf(w, "%s\n\n", report.Readiness.Summary)

	// Gates
	fmt.Fprintf(w, "Readiness Gates:\n")
	for _, gate := range report.Readiness.Gates {
		icon := "X"
		if gate.Passed {
			icon = "+"
		}
		fmt.Fprintf(w, "  [%s] %s: %s\n", icon, gate.Name, gate.Message)
	}
	fmt.Fprintf(w, "\n")

	// Specs by category
	fmt.Fprintf(w, "Specifications:\n")
	fmt.Fprintf(w, "  %-12s %-10s %-8s %-10s %-10s\n", "TYPE", "CATEGORY", "EXISTS", "EVAL", "APPROVED")
	fmt.Fprintf(w, "  %-12s %-10s %-8s %-10s %-10s\n", "----", "--------", "------", "----", "--------")

	for _, spec := range report.Specs {
		exists := "-"
		if spec.Exists {
			exists = "yes"
		}

		eval := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Exists {
			if spec.EvalStatus.Decision != "" {
				eval = spec.EvalStatus.Decision
			} else {
				eval = "yes"
			}
		}

		approved := "-"
		if spec.Approval != nil {
			approved = "yes"
		}

		required := ""
		if spec.Required {
			required = "*"
		}

		fmt.Fprintf(w, "  %-12s %-10s %-8s %-10s %-10s%s\n",
			spec.Type, spec.Category, exists, eval, approved, required)
	}

	fmt.Fprintf(w, "\n  * = required\n")

	// Summary
	fmt.Fprintf(w, "\nSummary:\n")
	fmt.Fprintf(w, "  Total: %d, Present: %d, Evaluated: %d, Approved: %d\n",
		report.Summary.TotalSpecs,
		report.Summary.PresentSpecs,
		report.Summary.EvaluatedSpecs,
		report.Summary.ApprovedSpecs)

	return nil
}

// RenderHTML renders the report as HTML.
func RenderHTML(w io.Writer, report *Report) error {
	// Traffic light color
	statusColor := "#dc3545" // red
	statusText := "NOT READY"
	if report.Readiness.Ready {
		statusColor = "#28a745" // green
		statusText = "READY"
	}

	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>MultiSpec Status: %s</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 2rem; }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: #666; margin-bottom: 2rem; }
    .status-badge { display: inline-block; padding: 0.5rem 1rem; border-radius: 4px; color: white; font-weight: bold; }
    .summary { font-size: 1.1rem; margin: 1rem 0; }
    table { border-collapse: collapse; width: 100%%; margin: 1rem 0; }
    th, td { border: 1px solid #ddd; padding: 0.75rem; text-align: left; }
    th { background: #f5f5f5; }
    .gate { padding: 0.5rem 0; }
    .gate-pass { color: #28a745; }
    .gate-fail { color: #dc3545; }
    .required { color: #dc3545; font-weight: bold; }
    .yes { color: #28a745; }
    .no { color: #999; }
    footer { margin-top: 2rem; color: #999; font-size: 0.9rem; }
  </style>
</head>
<body>
  <h1>%s</h1>
  <div class="subtitle">%s</div>

  <div>
    <span class="status-badge" style="background: %s">%s</span>
  </div>
  <p class="summary">%s</p>

  <h2>Readiness Gates</h2>
  <div>
`, report.Project, report.Project, report.Path, statusColor, statusText, report.Readiness.Summary)

	for _, gate := range report.Readiness.Gates {
		class := "gate-fail"
		icon := "&#10007;"
		if gate.Passed {
			class = "gate-pass"
			icon = "&#10003;"
		}
		fmt.Fprintf(w, `    <div class="gate %s">%s %s: %s</div>
`, class, icon, gate.Name, gate.Message)
	}

	fmt.Fprintf(w, `  </div>

  <h2>Specifications</h2>
  <table>
    <thead>
      <tr>
        <th>Type</th>
        <th>Category</th>
        <th>Exists</th>
        <th>Evaluated</th>
        <th>Approved</th>
      </tr>
    </thead>
    <tbody>
`)

	for _, spec := range report.Specs {
		existsClass := "no"
		existsText := "-"
		if spec.Exists {
			existsClass = "yes"
			existsText = "Yes"
		}

		evalClass := "no"
		evalText := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Exists {
			evalClass = "yes"
			if spec.EvalStatus.Decision != "" {
				evalText = spec.EvalStatus.Decision
			} else {
				evalText = "Yes"
			}
		}

		approvedClass := "no"
		approvedText := "-"
		if spec.Approval != nil {
			approvedClass = "yes"
			approvedText = "Yes"
		}

		typeText := string(spec.Type)
		if spec.Required {
			typeText = fmt.Sprintf(`<span class="required">%s*</span>`, spec.Type)
		}

		fmt.Fprintf(w, `      <tr>
        <td>%s</td>
        <td>%s</td>
        <td class="%s">%s</td>
        <td class="%s">%s</td>
        <td class="%s">%s</td>
      </tr>
`, typeText, spec.Category, existsClass, existsText, evalClass, evalText, approvedClass, approvedText)
	}

	fmt.Fprintf(w, `    </tbody>
  </table>
  <p><span class="required">*</span> = required</p>

  <h2>Summary</h2>
  <ul>
    <li>Total specs: %d</li>
    <li>Present: %d</li>
    <li>Evaluated: %d</li>
    <li>Approved: %d</li>
  </ul>

  <footer>
    Generated at %s by MultiSpec
  </footer>
</body>
</html>
`, report.Summary.TotalSpecs, report.Summary.PresentSpecs, report.Summary.EvaluatedSpecs, report.Summary.ApprovedSpecs, report.GeneratedAt.Format(time.RFC3339))

	return nil
}

// RenderMarkdown renders the report as Markdown.
func RenderMarkdown(w io.Writer, report *Report) error {
	// Header
	fmt.Fprintf(w, "# Project Status: %s\n\n", report.Project)
	fmt.Fprintf(w, "**Path:** `%s`\n\n", report.Path)

	// Status badge
	if report.Readiness.Ready {
		fmt.Fprintf(w, "**Status:** :white_check_mark: READY\n\n")
	} else {
		fmt.Fprintf(w, "**Status:** :x: NOT READY\n\n")
	}
	fmt.Fprintf(w, "%s\n\n", report.Readiness.Summary)

	// Gates
	fmt.Fprintf(w, "## Readiness Gates\n\n")
	for _, gate := range report.Readiness.Gates {
		icon := ":x:"
		if gate.Passed {
			icon = ":white_check_mark:"
		}
		fmt.Fprintf(w, "- %s **%s**: %s\n", icon, gate.Name, gate.Message)
	}
	fmt.Fprintf(w, "\n")

	// Specs table
	fmt.Fprintf(w, "## Specifications\n\n")
	fmt.Fprintf(w, "| Type | Category | Exists | Evaluated | Approved |\n")
	fmt.Fprintf(w, "|------|----------|--------|-----------|----------|\n")

	for _, spec := range report.Specs {
		exists := "-"
		if spec.Exists {
			exists = ":white_check_mark:"
		}

		eval := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Exists {
			if spec.EvalStatus.Decision != "" {
				eval = spec.EvalStatus.Decision
			} else {
				eval = ":white_check_mark:"
			}
		}

		approved := "-"
		if spec.Approval != nil {
			approved = ":white_check_mark:"
		}

		typeText := string(spec.Type)
		if spec.Required {
			typeText = fmt.Sprintf("**%s***", spec.Type)
		}

		fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
			typeText, spec.Category, exists, eval, approved)
	}

	fmt.Fprintf(w, "\n\\* = required\n\n")

	// Summary
	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "- **Total:** %d\n", report.Summary.TotalSpecs)
	fmt.Fprintf(w, "- **Present:** %d\n", report.Summary.PresentSpecs)
	fmt.Fprintf(w, "- **Evaluated:** %d\n", report.Summary.EvaluatedSpecs)
	fmt.Fprintf(w, "- **Approved:** %d\n", report.Summary.ApprovedSpecs)
	fmt.Fprintf(w, "\n---\n")
	fmt.Fprintf(w, "*Generated at %s by MultiSpec*\n", report.GeneratedAt.Format(time.RFC3339))

	return nil
}
