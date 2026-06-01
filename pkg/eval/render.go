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

	// Decision
	decision := "❌ FAIL"
	if result.Passed {
		decision = "✅ PASS"
	}
	fmt.Fprintf(w, "**Decision:** %s\n", decision)
	fmt.Fprintf(w, "**Score:** %.1f/10\n", result.Score)
	fmt.Fprintf(w, "**Timestamp:** %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05 MST"))

	// Categories
	fmt.Fprintf(w, "## Category Scores\n\n")
	fmt.Fprintf(w, "| Category | Score | Weight | Status |\n")
	fmt.Fprintf(w, "|----------|-------|--------|--------|\n")
	for _, cat := range result.Categories {
		status := "❌"
		if cat.Score >= 7.0 {
			status = "✅"
		} else if cat.Score >= 5.0 {
			status = "⚠️"
		}
		fmt.Fprintf(w, "| %s | %.1f | %.0f%% | %s |\n", cat.Name, cat.Score, cat.Weight*100, status)
	}
	fmt.Fprintln(w)

	// Category details
	fmt.Fprintf(w, "## Category Details\n\n")
	for _, cat := range result.Categories {
		fmt.Fprintf(w, "### %s (%.1f/10)\n\n", cat.Name, cat.Score)
		if cat.Explanation != "" {
			fmt.Fprintf(w, "%s\n\n", cat.Explanation)
		}
	}

	// Findings
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
				fmt.Fprintf(w, "**Category:** %s\n\n", f.Category)
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
	// Header
	decision := "✗ FAIL"
	if result.Passed {
		decision = "✓ PASS"
	}
	fmt.Fprintf(w, "\n%s %s: %.1f/10 %s\n\n", string(result.SpecType), decision, result.Score, result.Decision)

	// Summary
	if result.Summary != "" {
		fmt.Fprintf(w, "%s\n\n", result.Summary)
	}

	// Categories
	fmt.Fprintf(w, "Categories:\n")
	for _, cat := range result.Categories {
		status := "✗"
		if cat.Score >= 7.0 {
			status = "✓"
		} else if cat.Score >= 5.0 {
			status = "~"
		}
		fmt.Fprintf(w, "  %s %-20s %.1f/10 (%.0f%%)\n", status, cat.Name, cat.Score, cat.Weight*100)
		if r.Verbose && cat.Explanation != "" {
			fmt.Fprintf(w, "    %s\n", cat.Explanation)
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
				fmt.Fprintf(w, "  [%s] %s: %s\n", strings.ToUpper(f.Severity), f.Title, f.Description)
				if f.Recommendation != "" {
					fmt.Fprintf(w, "    → %s\n", f.Recommendation)
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
	fmt.Fprintf(w, "- **Generated By:** %s\n\n", report.Metadata.GeneratedBy)

	// Decision (Decision is a value type, not pointer - check if Status is set)
	if report.Decision.Status != "" {
		decision := "❌ FAIL"
		if report.Decision.Passed {
			decision = "✅ PASS"
		}
		fmt.Fprintf(w, "## Decision\n\n")
		fmt.Fprintf(w, "**Status:** %s (%s)\n\n", decision, report.Decision.Status)
		if report.Decision.Rationale != "" {
			fmt.Fprintf(w, "%s\n\n", report.Decision.Rationale)
		}
	}

	// Categories
	if len(report.Categories) > 0 {
		fmt.Fprintf(w, "## Categories\n\n")
		fmt.Fprintf(w, "| Category | Score |\n")
		fmt.Fprintf(w, "|----------|-------|\n")
		for _, cat := range report.Categories {
			fmt.Fprintf(w, "| %s | %s |\n", cat.Category, cat.Score)
		}
		fmt.Fprintln(w)

		// Category reasoning details
		fmt.Fprintf(w, "### Category Details\n\n")
		for _, cat := range report.Categories {
			fmt.Fprintf(w, "#### %s (%s)\n\n", cat.Category, cat.Score)
			if cat.Reasoning != "" {
				fmt.Fprintf(w, "%s\n\n", cat.Reasoning)
			}
		}
	}

	// Findings
	if len(report.Findings) > 0 {
		fmt.Fprintf(w, "## Findings\n\n")
		for _, f := range report.Findings {
			fmt.Fprintf(w, "### [%s] %s\n\n", f.Severity, f.Title)
			fmt.Fprintf(w, "%s\n\n", f.Description)
			if f.Recommendation != "" {
				fmt.Fprintf(w, "**Recommendation:** %s\n\n", f.Recommendation)
			}
		}
	}

	return nil
}
