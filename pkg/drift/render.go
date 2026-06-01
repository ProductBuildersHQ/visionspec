package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// RenderFormat specifies the output format.
type RenderFormat string

const (
	FormatText     RenderFormat = "text"
	FormatJSON     RenderFormat = "json"
	FormatMarkdown RenderFormat = "markdown"
)

// Renderer renders drift reports.
type Renderer struct {
	format RenderFormat
}

// NewRenderer creates a renderer with the specified format.
func NewRenderer(format RenderFormat) *Renderer {
	return &Renderer{format: format}
}

// Render writes the drift report to the writer.
func (r *Renderer) Render(w io.Writer, report *DriftReport) error {
	switch r.format {
	case FormatJSON:
		return r.renderJSON(w, report)
	case FormatMarkdown:
		return r.renderMarkdown(w, report)
	default:
		return r.renderText(w, report)
	}
}

// renderText renders the report as terminal text.
func (r *Renderer) renderText(w io.Writer, report *DriftReport) error {
	// Header
	fmt.Fprintf(w, "Drift Report: %s\n", report.Project)
	fmt.Fprintf(w, "Generated: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(w)

	// Summary
	if report.Summary.TotalItems == 0 {
		fmt.Fprintln(w, "✓ No drift detected - spec and code are aligned")
		return nil
	}

	// Status
	if report.HasBlockers() {
		fmt.Fprintln(w, "✗ DRIFT DETECTED - blockers found")
	} else {
		fmt.Fprintln(w, "⚠ DRIFT DETECTED - no blockers")
	}
	fmt.Fprintln(w)

	// Summary stats
	fmt.Fprintf(w, "Summary:\n")
	fmt.Fprintf(w, "  Total items:     %d\n", report.Summary.TotalItems)
	fmt.Fprintf(w, "  Critical:        %d\n", report.Summary.CriticalCount)
	fmt.Fprintf(w, "  High:            %d\n", report.Summary.HighCount)
	fmt.Fprintf(w, "  Unimplemented:   %d\n", report.Summary.ByType[DriftUnimplemented])
	fmt.Fprintf(w, "  Undocumented:    %d\n", report.Summary.ByType[DriftUndocumented])
	fmt.Fprintf(w, "  Mismatched:      %d\n", report.Summary.ByType[DriftMismatch])
	fmt.Fprintln(w)

	// Group items by severity
	fmt.Fprintln(w, "Findings:")
	fmt.Fprintln(w)

	// Critical items first
	for _, item := range report.Items {
		if item.Severity == SeverityCritical {
			r.renderTextItem(w, item)
		}
	}

	// High
	for _, item := range report.Items {
		if item.Severity == SeverityHigh {
			r.renderTextItem(w, item)
		}
	}

	// Medium
	for _, item := range report.Items {
		if item.Severity == SeverityMedium {
			r.renderTextItem(w, item)
		}
	}

	// Low
	for _, item := range report.Items {
		if item.Severity == SeverityLow {
			r.renderTextItem(w, item)
		}
	}

	return nil
}

func (r *Renderer) renderTextItem(w io.Writer, item DriftItem) {
	// Severity icon
	icon := "○"
	switch item.Severity {
	case SeverityCritical:
		icon = "●" // filled
	case SeverityHigh:
		icon = "◐" // half
	case SeverityMedium:
		icon = "◔" // quarter
	}

	// Type indicator
	typeStr := string(item.Type)
	switch item.Type {
	case DriftUnimplemented:
		typeStr = "MISSING"
	case DriftUndocumented:
		typeStr = "UNDOC"
	case DriftMismatch:
		typeStr = "MISMATCH"
	}

	fmt.Fprintf(w, "  %s [%s] %s\n", icon, strings.ToUpper(string(item.Severity)), item.Description)
	fmt.Fprintf(w, "    Type: %s, Category: %s\n", typeStr, item.Category)
	if item.SpecRef != "" {
		fmt.Fprintf(w, "    Spec: %s\n", item.SpecRef)
	}
	if item.CodeRef != "" {
		fmt.Fprintf(w, "    Code: %s\n", item.CodeRef)
	}
	if item.Suggestion != "" {
		fmt.Fprintf(w, "    → %s\n", item.Suggestion)
	}
	fmt.Fprintln(w)
}

// renderJSON renders the report as JSON.
func (r *Renderer) renderJSON(w io.Writer, report *DriftReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// renderMarkdown renders the report as Markdown.
func (r *Renderer) renderMarkdown(w io.Writer, report *DriftReport) error {
	// Header
	fmt.Fprintf(w, "# Drift Report: %s\n\n", report.Project)
	fmt.Fprintf(w, "_Generated: %s_\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

	// Summary
	if report.Summary.TotalItems == 0 {
		fmt.Fprintln(w, "✅ **No drift detected** - spec and code are aligned")
		return nil
	}

	// Status badge
	if report.HasBlockers() {
		fmt.Fprintln(w, "❌ **DRIFT DETECTED** - blockers found")
	} else {
		fmt.Fprintln(w, "⚠️ **DRIFT DETECTED** - no blockers")
	}

	// Summary table
	fmt.Fprintln(w, "## Summary")
	fmt.Fprintln(w, "| Metric | Count |")
	fmt.Fprintln(w, "|--------|-------|")
	fmt.Fprintf(w, "| Total items | %d |\n", report.Summary.TotalItems)
	fmt.Fprintf(w, "| Critical | %d |\n", report.Summary.CriticalCount)
	fmt.Fprintf(w, "| High | %d |\n", report.Summary.HighCount)
	fmt.Fprintf(w, "| Unimplemented | %d |\n", report.Summary.ByType[DriftUnimplemented])
	fmt.Fprintf(w, "| Undocumented | %d |\n", report.Summary.ByType[DriftUndocumented])
	fmt.Fprintf(w, "| Mismatched | %d |\n", report.Summary.ByType[DriftMismatch])
	fmt.Fprintln(w)

	// By category
	if len(report.Summary.ByCategory) > 0 {
		fmt.Fprintln(w, "### By Category")
		fmt.Fprintln(w, "| Category | Count |")
		fmt.Fprintln(w, "|----------|-------|")
		for cat, count := range report.Summary.ByCategory {
			fmt.Fprintf(w, "| %s | %d |\n", cat, count)
		}
		fmt.Fprintln(w)
	}

	// Findings
	fmt.Fprintln(w, "## Findings")

	// Group by severity
	for _, severity := range []Severity{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow} {
		items := report.FilterBySeverity(severity)
		filteredItems := []DriftItem{}
		for _, item := range items {
			if item.Severity == severity {
				filteredItems = append(filteredItems, item)
			}
		}

		if len(filteredItems) == 0 {
			continue
		}

		// Section header
		severityEmoji := "ℹ️"
		switch severity {
		case SeverityCritical:
			severityEmoji = "🔴"
		case SeverityHigh:
			severityEmoji = "🟠"
		case SeverityMedium:
			severityEmoji = "🟡"
		case SeverityLow:
			severityEmoji = "🟢"
		}

		fmt.Fprintf(w, "### %s %s\n\n", severityEmoji, strings.Title(string(severity)))

		for _, item := range filteredItems {
			r.renderMarkdownItem(w, item)
		}
	}

	return nil
}

func (r *Renderer) renderMarkdownItem(w io.Writer, item DriftItem) {
	// Type badge
	typeBadge := "📋"
	switch item.Type {
	case DriftUnimplemented:
		typeBadge = "❓ Missing"
	case DriftUndocumented:
		typeBadge = "📝 Undocumented"
	case DriftMismatch:
		typeBadge = "⚡ Mismatch"
	}

	fmt.Fprintf(w, "#### %s `%s`\n\n", typeBadge, item.ID)
	fmt.Fprintf(w, "%s\n\n", item.Description)

	fmt.Fprintf(w, "- **Category:** %s\n", item.Category)
	if item.SpecRef != "" {
		fmt.Fprintf(w, "- **Spec Reference:** %s\n", item.SpecRef)
	}
	if item.CodeRef != "" {
		fmt.Fprintf(w, "- **Code Reference:** `%s`\n", item.CodeRef)
	}
	if item.Suggestion != "" {
		fmt.Fprintf(w, "- **Suggestion:** %s\n", item.Suggestion)
	}
	fmt.Fprintln(w)
}
