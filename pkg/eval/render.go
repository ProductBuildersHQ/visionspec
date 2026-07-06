// Package eval provides evaluation orchestration for spec documents.
package eval

import (
	"fmt"
	"io"
	"strings"

	"github.com/plexusone/structured-evaluation/rubric"
)

// Renderer renders evaluation results to various formats.
type Renderer interface {
	Render(w io.Writer, result *Result) error
}

// MarkdownRenderer renders results to markdown format.
type MarkdownRenderer struct{}

// NewMarkdownRenderer creates a new markdown renderer.
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Render writes the result as markdown to the writer.
func (r *MarkdownRenderer) Render(w io.Writer, result *Result) error {
	// Header
	fmt.Fprintf(w, "# %s Evaluation Report\n\n", strings.ToUpper(string(result.SpecType)))

	// Summary
	fmt.Fprintf(w, "## Summary\n\n")
	if result.Summary != "" {
		fmt.Fprintf(w, "%s\n\n", result.Summary)
	}

	// Decision with v2 format
	decision := "FAIL"
	if result.Passed {
		decision = "PASS"
	}
	fmt.Fprintf(w, "**Decision:** %s %s\n", decision, result.Decision)

	// Score display: v2 (1-5) or legacy (0-10)
	if result.IntScore > 0 {
		fmt.Fprintf(w, "**Score:** %d/5 (%s)\n", result.IntScore, result.IntScore.String())
	} else {
		fmt.Fprintf(w, "**Score:** %.1f/10\n", result.Score)
	}

	// Confidence
	if result.Confidence > 0 {
		confidenceLabel := "High"
		if result.Confidence < 0.7 {
			confidenceLabel = "Low"
		} else if result.Confidence < 0.9 {
			confidenceLabel = "Medium"
		}
		fmt.Fprintf(w, "**Confidence:** %.0f%% (%s)\n", result.Confidence*100, confidenceLabel)
	}

	fmt.Fprintf(w, "**Timestamp:** %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05 MST"))

	// Blocking codes
	if len(result.Blocking) > 0 {
		fmt.Fprintf(w, "**Blocking Issues:**\n\n")
		for _, code := range result.Blocking {
			info := rubric.GetReasonCodeInfo(code)
			if info != nil {
				fmt.Fprintf(w, "- `%s`: %s\n", code, info.Description)
			} else {
				fmt.Fprintf(w, "- `%s`\n", code)
			}
		}
		fmt.Fprintln(w)
	}

	// Categories with v2 format
	fmt.Fprintf(w, "## Category Scores\n\n")
	fmt.Fprintf(w, "| Category | Score | Confidence | Weight | Status |\n")
	fmt.Fprintf(w, "|----------|-------|------------|--------|--------|\n")
	for _, cat := range result.Categories {
		status := "Fail"
		if cat.IntScore >= 4 {
			status = "Pass"
		} else if cat.IntScore >= 3 {
			status = "Partial"
		}

		scoreStr := fmt.Sprintf("%d (%s)", cat.IntScore, cat.IntScore.String())
		confStr := fmt.Sprintf("%.0f%%", cat.Confidence*100)
		if cat.Confidence == 0 {
			confStr = "-"
		}
		fmt.Fprintf(w, "| %s | %s | %s | %.0f%% | %s |\n", cat.Name, scoreStr, confStr, cat.Weight*100, status)
	}
	fmt.Fprintln(w)

	// Category details with reason codes
	fmt.Fprintf(w, "## Category Details\n\n")
	for _, cat := range result.Categories {
		fmt.Fprintf(w, "### %s (%d/5 - %s)\n\n", cat.Name, cat.IntScore, cat.IntScore.String())
		if cat.Explanation != "" {
			fmt.Fprintf(w, "%s\n\n", cat.Explanation)
		}
		if len(cat.ReasonCodes) > 0 {
			fmt.Fprintf(w, "**Reason Codes:** ")
			for i, code := range cat.ReasonCodes {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "`%s`", code)
			}
			fmt.Fprintf(w, "\n\n")
		}
	}

	// Findings with v2 fields
	if len(result.Findings) > 0 {
		fmt.Fprintf(w, "## Findings\n\n")

		// Group by severity
		bySeverity := make(map[string][]Finding)
		for _, f := range result.Findings {
			bySeverity[f.Severity] = append(bySeverity[f.Severity], f)
		}

		// Order: critical, high, medium, low, info
		severityOrder := []string{"critical", "high", "medium", "low", "info"}
		for _, sev := range severityOrder {
			findings := bySeverity[sev]
			if len(findings) == 0 {
				continue
			}
			fmt.Fprintf(w, "### %s\n\n", strings.ToUpper(sev))
			for _, f := range findings {
				fmt.Fprintf(w, "#### %s\n\n", f.Title)
				fmt.Fprintf(w, "**Category:** %s\n", f.Category)
				if f.Code != "" {
					fmt.Fprintf(w, "**Code:** `%s`\n", f.Code)
				}
				if f.Location != "" {
					fmt.Fprintf(w, "**Location:** %s\n", f.Location)
				}
				fmt.Fprintln(w)
				fmt.Fprintf(w, "%s\n\n", f.Description)
				if f.Recommendation != "" {
					fmt.Fprintf(w, "**Recommendation:** %s\n\n", f.Recommendation)
				}
				if f.Evidence != "" {
					fmt.Fprintf(w, "**Evidence:** %s\n\n", f.Evidence)
				}
			}
		}
	}

	// Judge metadata
	fmt.Fprintf(w, "## Evaluation Metadata\n\n")
	fmt.Fprintf(w, "- **Model:** %s\n", result.Judge.Model)
	fmt.Fprintf(w, "- **Provider:** %s\n", result.Judge.Provider)
	fmt.Fprintf(w, "- **Temperature:** %.2f\n", result.Judge.Temperature)
	fmt.Fprintf(w, "- **Tokens:** %d\n", result.Judge.Tokens)

	return nil
}

// TerminalRenderer renders results for terminal output.
type TerminalRenderer struct {
	Verbose bool
}

// NewTerminalRenderer creates a new terminal renderer.
func NewTerminalRenderer(verbose bool) *TerminalRenderer {
	return &TerminalRenderer{Verbose: verbose}
}

// Render writes the result to the terminal.
func (r *TerminalRenderer) Render(w io.Writer, result *Result) error {
	// Header with v2 score format
	decision := "FAIL"
	if result.Passed {
		decision = "PASS"
	}

	// Score display
	scoreStr := fmt.Sprintf("%.1f/10", result.Score)
	if result.IntScore > 0 {
		scoreStr = fmt.Sprintf("%d/5 (%s)", result.IntScore, result.IntScore.String())
	}

	fmt.Fprintf(w, "\n%s %s: %s\n", strings.ToUpper(string(result.SpecType)), decision, scoreStr)

	// Confidence indicator
	if result.Confidence > 0 {
		confIcon := "High"
		if result.Confidence < 0.7 {
			confIcon = "Low"
		} else if result.Confidence < 0.9 {
			confIcon = "Med"
		}
		fmt.Fprintf(w, "Confidence: %.0f%% (%s)\n", result.Confidence*100, confIcon)
	}
	fmt.Fprintln(w)

	// Summary
	if result.Summary != "" {
		fmt.Fprintf(w, "%s\n\n", result.Summary)
	}

	// Blocking codes
	if len(result.Blocking) > 0 {
		fmt.Fprintf(w, "Blocking: ")
		for i, code := range result.Blocking {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "%s", code)
		}
		fmt.Fprintf(w, "\n\n")
	}

	// Categories with v2 format
	fmt.Fprintf(w, "Categories:\n")
	for _, cat := range result.Categories {
		status := "X"
		if cat.IntScore >= 4 {
			status = "+"
		} else if cat.IntScore >= 3 {
			status = "~"
		}

		scoreDisplay := fmt.Sprintf("%d/5", cat.IntScore)
		if cat.IntScore == 0 {
			scoreDisplay = fmt.Sprintf("%.1f/10", cat.Score)
		}

		fmt.Fprintf(w, "  %s %-20s %s (%.0f%%)\n", status, cat.Name, scoreDisplay, cat.Weight*100)

		if r.Verbose {
			if cat.Explanation != "" {
				fmt.Fprintf(w, "    %s\n", cat.Explanation)
			}
			if len(cat.ReasonCodes) > 0 {
				fmt.Fprintf(w, "    Codes: ")
				for i, code := range cat.ReasonCodes {
					if i > 0 {
						fmt.Fprintf(w, ", ")
					}
					fmt.Fprintf(w, "%s", code)
				}
				fmt.Fprintln(w)
			}
		}
	}
	fmt.Fprintln(w)

	// Findings summary
	if len(result.Findings) > 0 {
		// Count by severity
		counts := make(map[string]int)
		for _, f := range result.Findings {
			counts[f.Severity]++
		}

		fmt.Fprintf(w, "Findings: ")
		parts := []string{}
		if n := counts["critical"]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d critical", n))
		}
		if n := counts["high"]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d high", n))
		}
		if n := counts["medium"]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d medium", n))
		}
		if n := counts["low"]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d low", n))
		}
		if n := counts["info"]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d info", n))
		}
		fmt.Fprintf(w, "%s\n\n", strings.Join(parts, ", "))

		// Detailed findings in verbose mode
		if r.Verbose {
			for _, f := range result.Findings {
				codeStr := ""
				if f.Code != "" {
					codeStr = fmt.Sprintf(" [%s]", f.Code)
				}
				fmt.Fprintf(w, "  [%s]%s %s: %s\n", strings.ToUpper(f.Severity), codeStr, f.Title, f.Description)
				if f.Location != "" {
					fmt.Fprintf(w, "    Location: %s\n", f.Location)
				}
				if f.Recommendation != "" {
					fmt.Fprintf(w, "    Fix: %s\n", f.Recommendation)
				}
			}
			fmt.Fprintln(w)
		}
	}

	return nil
}

// RenderEvaluationReportMarkdown renders a structured-evaluation report to markdown.
func RenderEvaluationReportMarkdown(w io.Writer, report *rubric.Rubric) error {
	// Header
	fmt.Fprintf(w, "# Evaluation Report: %s\n\n", report.ReviewType)

	// Metadata
	fmt.Fprintf(w, "## Metadata\n\n")
	fmt.Fprintf(w, "- **Document:** %s\n", report.Metadata.Document)
	fmt.Fprintf(w, "- **Generated:** %s\n", report.Metadata.GeneratedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(w, "- **Generated By:** %s\n", report.Metadata.GeneratedBy)
	if report.SchemaVersion != "" {
		fmt.Fprintf(w, "- **Schema Version:** %s\n", report.SchemaVersion)
	}
	fmt.Fprintln(w)

	// Decision with v2 format
	fmt.Fprintf(w, "## Decision\n\n")
	decision := "FAIL"
	if report.Pass {
		decision = "PASS"
	}
	fmt.Fprintf(w, "**Status:** %s\n", decision)

	// v2 score
	if report.IntScore > 0 {
		fmt.Fprintf(w, "**Score:** %d/5 (%s)\n", report.IntScore, report.IntScore.String())
	}
	if report.Confidence > 0 {
		fmt.Fprintf(w, "**Confidence:** %.0f%%\n", report.Confidence*100)
	}
	fmt.Fprintln(w)

	// Blocking codes
	if len(report.Blocking) > 0 {
		fmt.Fprintf(w, "**Blocking Issues:**\n\n")
		for _, code := range report.Blocking {
			info := rubric.GetReasonCodeInfo(code)
			if info != nil {
				fmt.Fprintf(w, "- `%s`: %s\n", code, info.Description)
			} else {
				fmt.Fprintf(w, "- `%s`\n", code)
			}
		}
		fmt.Fprintln(w)
	}

	if report.Decision.Rationale != "" {
		fmt.Fprintf(w, "%s\n\n", report.Decision.Rationale)
	}

	// Categories with v2 format
	if len(report.Categories) > 0 {
		fmt.Fprintf(w, "## Categories\n\n")
		fmt.Fprintf(w, "| Category | Score | Confidence |\n")
		fmt.Fprintf(w, "|----------|-------|------------|\n")
		for _, cat := range report.Categories {
			scoreStr := string(cat.Score)
			if cat.IntScore > 0 {
				scoreStr = fmt.Sprintf("%d (%s)", cat.IntScore, cat.IntScore.String())
			}
			confStr := "-"
			if cat.Confidence > 0 {
				confStr = fmt.Sprintf("%.0f%%", cat.Confidence*100)
			}
			fmt.Fprintf(w, "| %s | %s | %s |\n", cat.Category, scoreStr, confStr)
		}
		fmt.Fprintln(w)

		// Category reasoning details
		fmt.Fprintf(w, "### Category Details\n\n")
		for _, cat := range report.Categories {
			scoreStr := string(cat.Score)
			if cat.IntScore > 0 {
				scoreStr = fmt.Sprintf("%d/5 - %s", cat.IntScore, cat.IntScore.String())
			}
			fmt.Fprintf(w, "#### %s (%s)\n\n", cat.Category, scoreStr)
			if cat.Reasoning != "" {
				fmt.Fprintf(w, "%s\n\n", cat.Reasoning)
			}
			if len(cat.ReasonCodes) > 0 {
				fmt.Fprintf(w, "**Reason Codes:** ")
				for i, code := range cat.ReasonCodes {
					if i > 0 {
						fmt.Fprintf(w, ", ")
					}
					fmt.Fprintf(w, "`%s`", code)
				}
				fmt.Fprintf(w, "\n\n")
			}
		}
	}

	// Findings with v2 fields
	if len(report.Findings) > 0 {
		fmt.Fprintf(w, "## Findings\n\n")
		for _, f := range report.Findings {
			codeStr := ""
			if f.Code != "" {
				codeStr = fmt.Sprintf(" `%s`", f.Code)
			}
			fmt.Fprintf(w, "### [%s]%s %s\n\n", f.Severity, codeStr, f.Title)
			if f.Location != "" {
				fmt.Fprintf(w, "**Location:** %s\n\n", f.Location)
			}
			fmt.Fprintf(w, "%s\n\n", f.Description)
			if f.Recommendation != "" {
				fmt.Fprintf(w, "**Recommendation:** %s\n\n", f.Recommendation)
			}
		}
	}

	return nil
}
