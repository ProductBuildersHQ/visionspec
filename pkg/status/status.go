// Package status generates project status reports.
package status

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Report represents a project status report.
type Report struct {
	Project      string                `json:"project"`
	Path         string                `json:"path"`
	GeneratedAt  time.Time             `json:"generated_at"`
	Readiness    types.ReadinessStatus `json:"readiness"`
	Specs        []SpecStatus          `json:"specs"`
	Summary      Summary               `json:"summary"`
	GraphMetrics *GraphMetrics         `json:"graph_metrics,omitempty"`
}

// GraphMetrics contains traceability and graph statistics.
type GraphMetrics struct {
	TotalRequirements int     `json:"total_requirements"`
	TotalUserStories  int     `json:"total_user_stories"`
	TotalConstraints  int     `json:"total_constraints"`
	TotalDecisions    int     `json:"total_decisions"`
	TraceCoverage     float64 `json:"trace_coverage"` // Percentage of requirements traced to TRD
	ConflictCount     int     `json:"conflict_count"`
	GraphPath         string  `json:"graph_path,omitempty"`
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
		Info     int `json:"info"`
	} `json:"findings"`
	// Categories breakdown from rubric evaluation
	Categories *CategoryBreakdown `json:"categories,omitempty"`
}

// CategoryBreakdown represents category pass/partial/fail counts.
type CategoryBreakdown struct {
	Pass    int `json:"pass"`
	Partial int `json:"partial"`
	Fail    int `json:"fail"`
	Total   int `json:"total"`
}

// Summary provides aggregate statistics.
type Summary struct {
	TotalSpecs     int `json:"total_specs"`
	PresentSpecs   int `json:"present_specs"`
	EvaluatedSpecs int `json:"evaluated_specs"`
	ApprovedSpecs  int `json:"approved_specs"`
	BlockingIssues int `json:"blocking_issues"`
}

// RichReport extends Report with LLM-filled semantic fields for enhanced display.
// The deterministic fields (Pipeline, AggregateCategories) are computed from specs.
// The semantic fields (Highlights, MediumFindings, NextSteps, KeyDecisions) are LLM-generated.
type RichReport struct {
	*Report

	// Pipeline is the ordered list of specs with their statuses for visualization.
	Pipeline []PipelineStage `json:"pipeline"`

	// AggregateCategories is the total category breakdown across all specs.
	AggregateCategories CategoryBreakdown `json:"aggregate_categories"`

	// CompletionPercent is the percentage of specs that are complete.
	CompletionPercent int `json:"completion_percent"`

	// --- LLM-filled semantic fields ---

	// Highlights maps spec type to a brief highlight string (LLM-generated).
	// Example: "MRD" -> "Market analysis, competitive positioning"
	Highlights map[types.SpecType]string `json:"highlights,omitempty"`

	// MediumFindings is a list of non-blocking medium-severity findings (LLM-generated).
	MediumFindings []Finding `json:"medium_findings,omitempty"`

	// NextSteps is a prioritized list of recommended next actions (LLM-generated).
	NextSteps []string `json:"next_steps,omitempty"`

	// KeyDecisions are important architectural/design decisions extracted from specs (LLM-generated).
	KeyDecisions []KeyDecision `json:"key_decisions,omitempty"`

	// ReadySummary is a one-line summary of readiness state (LLM-generated).
	ReadySummary string `json:"ready_summary,omitempty"`
}

// PipelineStage represents a single stage in the pipeline visualization.
type PipelineStage struct {
	Type   types.SpecType `json:"type"`
	Status StageStatus    `json:"status"` // complete, pending, missing
	Label  string         `json:"label"`  // Display label (e.g., "MRD", "spec.md")
}

// StageStatus represents the status of a pipeline stage.
type StageStatus string

const (
	StageComplete StageStatus = "complete"
	StagePending  StageStatus = "pending"
	StageMissing  StageStatus = "missing"
)

// Finding represents a single evaluation finding with description.
type Finding struct {
	Spec        types.SpecType `json:"spec"`
	Severity    string         `json:"severity"` // critical, high, medium, low, info
	Description string         `json:"description"`
}

// KeyDecision represents an architectural or design decision.
type KeyDecision struct {
	Area   string `json:"area"`   // e.g., "Local Dev", "Cloud", "IaC"
	Choice string `json:"choice"` // e.g., "AWS with Pulumi Go SDK"
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

// WithGraphMetrics adds graph metrics to a report.
// This is a separate call because graph extraction can be expensive.
func (r *Report) WithGraphMetrics(metrics *GraphMetrics) *Report {
	r.GraphMetrics = metrics
	return r
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

	// Graph metrics if available
	if report.GraphMetrics != nil {
		fmt.Fprintf(w, "\nTraceability Metrics:\n")
		fmt.Fprintf(w, "  Requirements: %d, User Stories: %d, Constraints: %d, Decisions: %d\n",
			report.GraphMetrics.TotalRequirements,
			report.GraphMetrics.TotalUserStories,
			report.GraphMetrics.TotalConstraints,
			report.GraphMetrics.TotalDecisions)
		fmt.Fprintf(w, "  Trace Coverage: %.0f%%, Conflicts: %d\n",
			report.GraphMetrics.TraceCoverage*100,
			report.GraphMetrics.ConflictCount)
	}

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
  <title>VisionSpec Status: %s</title>
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
    Generated at %s by VisionSpec
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

	// Graph metrics if available
	if report.GraphMetrics != nil {
		fmt.Fprintf(w, "\n## Traceability Metrics\n\n")
		fmt.Fprintf(w, "- **Requirements:** %d extracted\n", report.GraphMetrics.TotalRequirements)
		fmt.Fprintf(w, "- **User Stories:** %d extracted\n", report.GraphMetrics.TotalUserStories)
		fmt.Fprintf(w, "- **Constraints:** %d extracted\n", report.GraphMetrics.TotalConstraints)
		fmt.Fprintf(w, "- **Decisions:** %d extracted\n", report.GraphMetrics.TotalDecisions)
		fmt.Fprintf(w, "- **Trace Coverage:** %.0f%%\n", report.GraphMetrics.TraceCoverage*100)
		fmt.Fprintf(w, "- **Conflicts:** %d detected\n", report.GraphMetrics.ConflictCount)
		if report.GraphMetrics.GraphPath != "" {
			fmt.Fprintf(w, "\n[View Graph →](%s)\n", report.GraphMetrics.GraphPath)
		}
	}

	fmt.Fprintf(w, "\n---\n")
	fmt.Fprintf(w, "*Generated at %s by VisionSpec*\n", report.GeneratedAt.Format(time.RFC3339))

	return nil
}

// NewRichReport creates a RichReport from a Report, computing deterministic fields.
// LLM-filled fields (Highlights, MediumFindings, NextSteps, KeyDecisions) are left empty
// and should be populated separately via the Set* methods.
func NewRichReport(report *Report) *RichReport {
	rr := &RichReport{
		Report:     report,
		Pipeline:   buildPipeline(report),
		Highlights: make(map[types.SpecType]string),
	}

	// Compute aggregate categories
	for _, spec := range report.Specs {
		if spec.EvalStatus != nil && spec.EvalStatus.Categories != nil {
			rr.AggregateCategories.Pass += spec.EvalStatus.Categories.Pass
			rr.AggregateCategories.Partial += spec.EvalStatus.Categories.Partial
			rr.AggregateCategories.Fail += spec.EvalStatus.Categories.Fail
			rr.AggregateCategories.Total += spec.EvalStatus.Categories.Total
		}
	}

	// Compute completion percent
	complete := 0
	for _, stage := range rr.Pipeline {
		if stage.Status == StageComplete {
			complete++
		}
	}
	if len(rr.Pipeline) > 0 {
		rr.CompletionPercent = (complete * 100) / len(rr.Pipeline)
	}

	return rr
}

// buildPipeline creates an ordered pipeline from the report specs.
func buildPipeline(report *Report) []PipelineStage {
	// Define the standard pipeline order
	pipelineOrder := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		types.SpecTypeTRD,
		types.SpecTypeTPD,
		types.SpecTypeIRD,
		"spec", // Special case for spec.md
	}

	// Build lookup map
	specMap := make(map[types.SpecType]SpecStatus)
	for _, spec := range report.Specs {
		specMap[spec.Type] = spec
	}

	stages := make([]PipelineStage, 0, len(pipelineOrder))
	for _, specType := range pipelineOrder {
		stage := PipelineStage{
			Type:  specType,
			Label: string(specType),
		}

		// Special handling for spec.md
		if specType == "spec" {
			stage.Label = "spec.md"
			specPath := filepath.Join(report.Path, "spec.md")
			if _, err := os.Stat(specPath); err == nil {
				stage.Status = StageComplete
			} else {
				stage.Status = StageMissing
			}
			stages = append(stages, stage)
			continue
		}

		// Check spec status
		if spec, ok := specMap[specType]; ok {
			if spec.Exists {
				if spec.EvalStatus != nil && spec.EvalStatus.Decision == "pass" {
					stage.Status = StageComplete
				} else if spec.Approval != nil {
					stage.Status = StageComplete
				} else {
					stage.Status = StagePending
				}
			} else {
				stage.Status = StageMissing
			}
		} else {
			stage.Status = StageMissing
		}

		stages = append(stages, stage)
	}

	return stages
}

// SetHighlight sets the LLM-generated highlight for a spec type.
func (rr *RichReport) SetHighlight(specType types.SpecType, highlight string) {
	rr.Highlights[specType] = highlight
}

// SetMediumFindings sets the LLM-generated medium findings.
func (rr *RichReport) SetMediumFindings(findings []Finding) {
	rr.MediumFindings = findings
}

// SetNextSteps sets the LLM-generated next steps.
func (rr *RichReport) SetNextSteps(steps []string) {
	rr.NextSteps = steps
}

// SetKeyDecisions sets the LLM-generated key decisions.
func (rr *RichReport) SetKeyDecisions(decisions []KeyDecision) {
	rr.KeyDecisions = decisions
}

// SetReadySummary sets the LLM-generated readiness summary.
func (rr *RichReport) SetReadySummary(summary string) {
	rr.ReadySummary = summary
}

// RenderRichText renders the rich report as terminal text with box-drawing.
func RenderRichText(w io.Writer, rr *RichReport) error {
	fmt.Fprintf(w, "⏺ VisionSpec Status\n\n")

	// Pipeline Progress
	if err := RenderPipelineProgress(w, rr.Pipeline); err != nil {
		return err
	}

	// Summary Table
	fmt.Fprintf(w, "\n  Summary\n")
	if err := RenderSummaryTable(w, rr); err != nil {
		return err
	}

	// Aggregate stats
	fmt.Fprintf(w, "  Overall: %d/%d specs complete (%d%%)\n\n",
		countComplete(rr.Pipeline), len(rr.Pipeline), rr.CompletionPercent)
	fmt.Fprintf(w, "  Aggregate: %d pass, %d partial, %d fail across %d categories\n\n",
		rr.AggregateCategories.Pass,
		rr.AggregateCategories.Partial,
		rr.AggregateCategories.Fail,
		rr.AggregateCategories.Total)

	// Medium Findings (if any)
	if len(rr.MediumFindings) > 0 {
		if err := RenderFindings(w, rr.MediumFindings); err != nil {
			return err
		}
	}

	// Next Steps (if any)
	if len(rr.NextSteps) > 0 {
		if err := RenderNextSteps(w, rr.NextSteps); err != nil {
			return err
		}
	}

	// Key Decisions (if any)
	if len(rr.KeyDecisions) > 0 {
		if err := RenderKeyDecisions(w, rr.KeyDecisions); err != nil {
			return err
		}
	}

	// Ready summary
	if rr.ReadySummary != "" {
		fmt.Fprintf(w, "  Ready for: %s\n", rr.ReadySummary)
	}

	return nil
}

// RenderPipelineProgress renders the pipeline visualization with arrows and status icons.
func RenderPipelineProgress(w io.Writer, pipeline []PipelineStage) error {
	fmt.Fprintf(w, "  Pipeline Progress\n\n")

	// Build label line with arrows
	var labels []string
	for _, stage := range pipeline {
		labels = append(labels, stage.Label)
	}

	// Print labels with arrows
	fmt.Fprintf(w, "  ")
	for i, label := range labels {
		if i > 0 {
			fmt.Fprintf(w, " → ")
		}
		fmt.Fprintf(w, "%s", label)
	}
	fmt.Fprintf(w, "\n")

	// Print status icons aligned under labels
	fmt.Fprintf(w, "  ")
	for i, stage := range pipeline {
		if i > 0 {
			fmt.Fprintf(w, "   ") // Space for " → "
		}
		icon := statusIcon(stage.Status)
		// Pad to match label width
		padding := len(stage.Label) - runeWidth(icon)
		if padding > 0 {
			fmt.Fprintf(w, "%s%s", icon, spaces(padding))
		} else {
			fmt.Fprintf(w, "%s", icon)
		}
	}
	fmt.Fprintf(w, "\n")

	return nil
}

// RenderSummaryTable renders the summary table with box-drawing characters.
func RenderSummaryTable(w io.Writer, rr *RichReport) error {
	// Determine column widths
	specW := 7 // "Spec" column
	statusW := 10
	catW := 17
	findingsW := 15

	// Adjust findings column for highlights if present
	hasHighlights := len(rr.Highlights) > 0
	highlightsW := 40
	if hasHighlights {
		findingsW = highlightsW
	}

	// Header
	if hasHighlights {
		fmt.Fprintf(w, "  ┌%s┬%s┬%s┬%s┐\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", highlightsW+2))
		fmt.Fprintf(w, "  │ %-*s │ %-*s │ %-*s │ %-*s │\n",
			specW, "Spec", statusW, "Status", catW, "Categories", highlightsW, "Key Highlights")
		fmt.Fprintf(w, "  ├%s┼%s┼%s┼%s┤\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", highlightsW+2))
	} else {
		fmt.Fprintf(w, "  ┌%s┬%s┬%s┬%s┐\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", findingsW+2))
		fmt.Fprintf(w, "  │ %-*s │ %-*s │ %-*s │ %-*s │\n",
			specW, "Spec", statusW, "Status", catW, "Categories", findingsW, "Findings")
		fmt.Fprintf(w, "  ├%s┼%s┼%s┼%s┤\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", findingsW+2))
	}

	// Rows
	for i, spec := range rr.Report.Specs {
		specLabel := string(spec.Type)
		if spec.Type == "spec" {
			specLabel = "spec.md"
		}

		// Status
		statusStr := formatSpecStatus(spec)

		// Categories
		catStr := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Categories != nil {
			cat := spec.EvalStatus.Categories
			if cat.Partial == 0 && cat.Fail == 0 {
				catStr = fmt.Sprintf("%d/%d pass", cat.Pass, cat.Total)
			} else if cat.Fail == 0 {
				catStr = fmt.Sprintf("%d pass, %d partial", cat.Pass, cat.Partial)
			} else {
				catStr = fmt.Sprintf("%dp/%dpt/%df", cat.Pass, cat.Partial, cat.Fail)
			}
		}

		// Findings or Highlights
		lastCol := "-"
		if hasHighlights {
			if hl, ok := rr.Highlights[spec.Type]; ok {
				lastCol = truncate(hl, highlightsW)
			}
		} else {
			lastCol = formatFindings(spec.EvalStatus)
		}

		// Print row
		if hasHighlights {
			fmt.Fprintf(w, "  │ %-*s │ %-*s │ %-*s │ %-*s │\n",
				specW, specLabel, statusW, statusStr, catW, catStr, highlightsW, lastCol)
		} else {
			fmt.Fprintf(w, "  │ %-*s │ %-*s │ %-*s │ %-*s │\n",
				specW, specLabel, statusW, statusStr, catW, catStr, findingsW, lastCol)
		}

		// Row separator (except for last row)
		if i < len(rr.Report.Specs)-1 {
			if hasHighlights {
				fmt.Fprintf(w, "  ├%s┼%s┼%s┼%s┤\n",
					repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", highlightsW+2))
			} else {
				fmt.Fprintf(w, "  ├%s┼%s┼%s┼%s┤\n",
					repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", findingsW+2))
			}
		}
	}

	// Footer
	if hasHighlights {
		fmt.Fprintf(w, "  └%s┴%s┴%s┴%s┘\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", highlightsW+2))
	} else {
		fmt.Fprintf(w, "  └%s┴%s┴%s┴%s┘\n",
			repeat("─", specW+2), repeat("─", statusW+2), repeat("─", catW+2), repeat("─", findingsW+2))
	}

	return nil
}

// RenderFindings renders the medium findings list.
func RenderFindings(w io.Writer, findings []Finding) error {
	fmt.Fprintf(w, "  Medium Findings (non-blocking)\n\n")
	for i, f := range findings {
		fmt.Fprintf(w, "  %d. %s: %s\n", i+1, f.Spec, f.Description)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

// RenderNextSteps renders the next steps list.
func RenderNextSteps(w io.Writer, steps []string) error {
	fmt.Fprintf(w, "  Next Steps\n\n")
	for i, step := range steps {
		fmt.Fprintf(w, "  %d. %s\n", i+1, step)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

// RenderKeyDecisions renders the key decisions table with box-drawing.
func RenderKeyDecisions(w io.Writer, decisions []KeyDecision) error {
	// Determine column widths
	areaW := 13
	choiceW := 40

	for _, d := range decisions {
		if len(d.Area) > areaW {
			areaW = len(d.Area)
		}
		if len(d.Choice) > choiceW {
			choiceW = len(d.Choice)
		}
	}

	fmt.Fprintf(w, "  Key Decisions\n")
	fmt.Fprintf(w, "  ┌%s┬%s┐\n", repeat("─", areaW+2), repeat("─", choiceW+2))
	fmt.Fprintf(w, "  │ %-*s │ %-*s │\n", areaW, "Decision", choiceW, "Choice")
	fmt.Fprintf(w, "  ├%s┼%s┤\n", repeat("─", areaW+2), repeat("─", choiceW+2))

	for i, d := range decisions {
		fmt.Fprintf(w, "  │ %-*s │ %-*s │\n", areaW, d.Area, choiceW, d.Choice)
		if i < len(decisions)-1 {
			fmt.Fprintf(w, "  ├%s┼%s┤\n", repeat("─", areaW+2), repeat("─", choiceW+2))
		}
	}

	fmt.Fprintf(w, "  └%s┴%s┘\n", repeat("─", areaW+2), repeat("─", choiceW+2))
	return nil
}

// RenderRichMarkdown renders the rich report as Markdown with code blocks for tables.
func RenderRichMarkdown(w io.Writer, rr *RichReport) error {
	fmt.Fprintf(w, "# VisionSpec Status: %s\n\n", rr.Report.Project)

	// Pipeline Progress
	fmt.Fprintf(w, "## Pipeline Progress\n\n")
	fmt.Fprintf(w, "```\n")
	var labels []string
	for _, stage := range rr.Pipeline {
		labels = append(labels, stage.Label)
	}
	for i, label := range labels {
		if i > 0 {
			fmt.Fprintf(w, " → ")
		}
		fmt.Fprintf(w, "%s", label)
	}
	fmt.Fprintf(w, "\n")
	for i, stage := range rr.Pipeline {
		if i > 0 {
			fmt.Fprintf(w, "   ")
		}
		icon := statusIcon(stage.Status)
		padding := len(stage.Label) - runeWidth(icon)
		if padding > 0 {
			fmt.Fprintf(w, "%s%s", icon, spaces(padding))
		} else {
			fmt.Fprintf(w, "%s", icon)
		}
	}
	fmt.Fprintf(w, "\n```\n\n")

	// Completion
	fmt.Fprintf(w, "**Completion:** %d/%d specs (%d%%)\n\n",
		countComplete(rr.Pipeline), len(rr.Pipeline), rr.CompletionPercent)

	// Summary Table
	fmt.Fprintf(w, "## Summary\n\n")
	hasHighlights := len(rr.Highlights) > 0

	if hasHighlights {
		fmt.Fprintf(w, "| Spec | Status | Categories | Key Highlights |\n")
		fmt.Fprintf(w, "|------|--------|------------|----------------|\n")
	} else {
		fmt.Fprintf(w, "| Spec | Status | Categories | Findings |\n")
		fmt.Fprintf(w, "|------|--------|------------|----------|\n")
	}

	for _, spec := range rr.Report.Specs {
		specLabel := string(spec.Type)
		if spec.Type == "spec" {
			specLabel = "spec.md"
		}

		statusStr := formatSpecStatusMd(spec)
		catStr := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Categories != nil {
			cat := spec.EvalStatus.Categories
			if cat.Partial == 0 && cat.Fail == 0 {
				catStr = fmt.Sprintf("%d/%d pass", cat.Pass, cat.Total)
			} else if cat.Fail == 0 {
				catStr = fmt.Sprintf("%d pass, %d partial", cat.Pass, cat.Partial)
			} else {
				catStr = fmt.Sprintf("%dp/%dpt/%df", cat.Pass, cat.Partial, cat.Fail)
			}
		}

		lastCol := "-"
		if hasHighlights {
			if hl, ok := rr.Highlights[spec.Type]; ok {
				lastCol = hl
			}
		} else {
			lastCol = formatFindings(spec.EvalStatus)
		}

		fmt.Fprintf(w, "| %s | %s | %s | %s |\n", specLabel, statusStr, catStr, lastCol)
	}
	fmt.Fprintf(w, "\n")

	// Aggregate
	fmt.Fprintf(w, "**Aggregate:** %d pass, %d partial, %d fail across %d categories\n\n",
		rr.AggregateCategories.Pass,
		rr.AggregateCategories.Partial,
		rr.AggregateCategories.Fail,
		rr.AggregateCategories.Total)

	// Medium Findings
	if len(rr.MediumFindings) > 0 {
		fmt.Fprintf(w, "## Medium Findings (non-blocking)\n\n")
		for i, f := range rr.MediumFindings {
			fmt.Fprintf(w, "%d. **%s:** %s\n", i+1, f.Spec, f.Description)
		}
		fmt.Fprintf(w, "\n")
	}

	// Next Steps
	if len(rr.NextSteps) > 0 {
		fmt.Fprintf(w, "## Next Steps\n\n")
		for i, step := range rr.NextSteps {
			fmt.Fprintf(w, "%d. %s\n", i+1, step)
		}
		fmt.Fprintf(w, "\n")
	}

	// Key Decisions
	if len(rr.KeyDecisions) > 0 {
		fmt.Fprintf(w, "## Key Decisions\n\n")
		fmt.Fprintf(w, "| Decision | Choice |\n")
		fmt.Fprintf(w, "|----------|--------|\n")
		for _, d := range rr.KeyDecisions {
			fmt.Fprintf(w, "| %s | %s |\n", d.Area, d.Choice)
		}
		fmt.Fprintf(w, "\n")
	}

	// Ready summary
	if rr.ReadySummary != "" {
		fmt.Fprintf(w, "**Ready for:** %s\n\n", rr.ReadySummary)
	}

	fmt.Fprintf(w, "---\n")
	fmt.Fprintf(w, "*Generated at %s by VisionSpec*\n", rr.Report.GeneratedAt.Format(time.RFC3339))

	return nil
}

func formatSpecStatusMd(spec SpecStatus) string {
	if !spec.Exists {
		return ":x: Missing"
	}
	if spec.EvalStatus != nil {
		switch spec.EvalStatus.Decision {
		case "pass":
			return ":white_check_mark: Pass"
		case "conditional":
			return ":warning: Cond"
		case "fail":
			return ":x: Fail"
		}
	}
	if spec.Approval != nil {
		return ":white_check_mark: Pass"
	}
	return ":hourglass: Draft"
}

// Helper functions

func statusIcon(status StageStatus) string {
	switch status {
	case StageComplete:
		return "✅"
	case StagePending:
		return "🔄"
	case StageMissing:
		return "❌"
	default:
		return "?"
	}
}

func runeWidth(s string) int {
	// Emojis are typically 2 characters wide in terminal, but we count runes
	return 2 // Assume emoji width
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	return repeat(" ", n)
}

func repeat(s string, n int) string {
	return strings.Repeat(s, n)
}

func countComplete(pipeline []PipelineStage) int {
	count := 0
	for _, stage := range pipeline {
		if stage.Status == StageComplete {
			count++
		}
	}
	return count
}

func formatSpecStatus(spec SpecStatus) string {
	if !spec.Exists {
		return "❌ Missing"
	}
	if spec.EvalStatus != nil {
		switch spec.EvalStatus.Decision {
		case "pass":
			return "✅ Pass"
		case "conditional":
			return "⚠️ Cond"
		case "fail":
			return "❌ Fail"
		}
	}
	if spec.Approval != nil {
		return "✅ Pass"
	}
	return "🔄 Draft"
}

func formatFindings(eval *EvalStatus) string {
	if eval == nil {
		return "-"
	}
	parts := []string{}
	if eval.Findings.Critical > 0 {
		parts = append(parts, fmt.Sprintf("%d critical", eval.Findings.Critical))
	}
	if eval.Findings.High > 0 {
		parts = append(parts, fmt.Sprintf("%d high", eval.Findings.High))
	}
	if eval.Findings.Medium > 0 {
		parts = append(parts, fmt.Sprintf("%d medium", eval.Findings.Medium))
	}
	if eval.Findings.Low > 0 {
		parts = append(parts, fmt.Sprintf("%d low", eval.Findings.Low))
	}
	if eval.Findings.Info > 0 {
		parts = append(parts, fmt.Sprintf("%d info", eval.Findings.Info))
	}
	if len(parts) == 0 {
		return "-"
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
