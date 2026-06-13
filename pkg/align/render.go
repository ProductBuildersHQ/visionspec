package align

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

// OutputFormat specifies the output format.
type OutputFormat string

const (
	OutputFormatText     OutputFormat = "text"
	OutputFormatJSON     OutputFormat = "json"
	OutputFormatMarkdown OutputFormat = "markdown"
)

// RenderResult renders an alignment result in the specified format.
func RenderResult(result *AlignmentResult, format OutputFormat) (string, error) {
	switch format {
	case OutputFormatJSON:
		return renderJSON(result)
	case OutputFormatMarkdown:
		return renderMarkdown(result)
	case OutputFormatText:
		return renderText(result)
	default:
		return renderText(result)
	}
}

func renderJSON(result *AlignmentResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func renderText(result *AlignmentResult) (string, error) {
	var buf strings.Builder

	// Header
	fmt.Fprintf(&buf, "Alignment Report: %s\n", result.Project)
	fmt.Fprintf(&buf, "Generated: %s\n", result.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&buf, "Status: %s\n\n", statusLabel(result.Summary.IsAligned))

	// Summary
	fmt.Fprintf(&buf, "Summary:\n")
	fmt.Fprintf(&buf, "  Alignment Score: %.1f%%\n", result.Summary.AlignmentScore*100)
	fmt.Fprintf(&buf, "  Total Discrepancies: %d\n", result.Summary.TotalDiscrepancies)
	fmt.Fprintf(&buf, "  Critical: %d | High: %d\n", result.Summary.CriticalCount, result.Summary.HighCount)
	fmt.Fprintln(&buf)

	// Coverage
	fmt.Fprintf(&buf, "Coverage:\n")
	fmt.Fprintf(&buf, "  Requirements: %d total\n", result.Coverage.TotalRequirements)
	fmt.Fprintf(&buf, "  Implemented: %d (%.1f%%)\n", result.Coverage.ImplementedCount, result.Coverage.CoveragePercentage)
	fmt.Fprintf(&buf, "  Partial: %d | Missing: %d\n", result.Coverage.PartialCount, result.Coverage.MissingCount)
	if result.Coverage.UndocumentedCount > 0 {
		fmt.Fprintf(&buf, "  Undocumented Code: %d items\n", result.Coverage.UndocumentedCount)
	}
	fmt.Fprintln(&buf)

	// Discrepancies by severity
	if len(result.Discrepancies) > 0 {
		fmt.Fprintf(&buf, "Discrepancies:\n")

		// Group by severity
		bySeverity := map[Severity][]Discrepancy{}
		for _, d := range result.Discrepancies {
			bySeverity[d.Severity] = append(bySeverity[d.Severity], d)
		}

		// Print in order of severity
		for _, sev := range []Severity{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo} {
			items := bySeverity[sev]
			if len(items) == 0 {
				continue
			}

			fmt.Fprintf(&buf, "\n  [%s] (%d items)\n", strings.ToUpper(string(sev)), len(items))
			for _, d := range items {
				fmt.Fprintf(&buf, "    %s %s\n", severityIcon(d.Severity), truncate(d.Description, 70))
				if d.SpecRef != "" {
					fmt.Fprintf(&buf, "      Spec: %s\n", d.SpecRef)
				}
				if d.CodeRef != "" {
					fmt.Fprintf(&buf, "      Code: %s\n", d.CodeRef)
				}
				if d.Suggestion != "" {
					fmt.Fprintf(&buf, "      Fix: %s\n", d.Suggestion)
				}
			}
		}
	} else {
		fmt.Fprintf(&buf, "No discrepancies found - implementation aligns with spec!\n")
	}

	return buf.String(), nil
}

func renderMarkdown(result *AlignmentResult) (string, error) {
	tmpl, err := template.New("align").Funcs(template.FuncMap{
		"severityBadge": func(s Severity) string {
			switch s {
			case SeverityCritical:
				return "**CRITICAL**"
			case SeverityHigh:
				return "**HIGH**"
			case SeverityMedium:
				return "*MEDIUM*"
			case SeverityLow:
				return "LOW"
			default:
				return "INFO"
			}
		},
		"statusBadge": func(aligned bool) string {
			if aligned {
				return "ALIGNED"
			}
			return "MISALIGNED"
		},
		"percent": func(f float64) string {
			return fmt.Sprintf("%.1f%%", f*100)
		},
		"typeLabel": func(t DiscrepancyType) string {
			switch t {
			case DiscrepancyMissingFeature:
				return "Missing"
			case DiscrepancyUndocumentedCode:
				return "Undocumented"
			case DiscrepancyDiverged:
				return "Diverged"
			case DiscrepancyPartialImplementation:
				return "Partial"
			case DiscrepancyBehaviorMismatch:
				return "Mismatch"
			default:
				return string(t)
			}
		},
	}).Parse(alignmentMarkdownTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, result); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

const alignmentMarkdownTemplate = `# Alignment Report: {{.Project}}

> Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}

## Status: {{statusBadge .Summary.IsAligned}}

| Metric | Value |
|--------|-------|
| Alignment Score | {{percent .Summary.AlignmentScore}} |
| Total Discrepancies | {{.Summary.TotalDiscrepancies}} |
| Critical Issues | {{.Summary.CriticalCount}} |
| High Issues | {{.Summary.HighCount}} |

## Coverage

| Metric | Value |
|--------|-------|
| Total Requirements | {{.Coverage.TotalRequirements}} |
| Implemented | {{.Coverage.ImplementedCount}} |
| Partial | {{.Coverage.PartialCount}} |
| Missing | {{.Coverage.MissingCount}} |
| Coverage | {{printf "%.1f%%" .Coverage.CoveragePercentage}} |
{{if .Coverage.UndocumentedCount}}| Undocumented Code | {{.Coverage.UndocumentedCount}} |{{end}}

{{if .Discrepancies}}
## Discrepancies

{{range .Discrepancies}}
### {{severityBadge .Severity}}: {{typeLabel .Type}}

{{.Description}}

{{if .SpecRef}}- **Spec Reference:** ` + "`{{.SpecRef}}`" + `{{end}}
{{if .CodeRef}}- **Code Reference:** ` + "`{{.CodeRef}}`" + `{{end}}
{{if .Expected}}- **Expected:** {{.Expected}}{{end}}
{{if .Actual}}- **Actual:** {{.Actual}}{{end}}
{{if .Suggestion}}- **Suggestion:** {{.Suggestion}}{{end}}

{{end}}
{{else}}
## No Discrepancies

Implementation aligns with specification.
{{end}}

---
*Generated by VisionSpec align*
`

func statusLabel(aligned bool) string {
	if aligned {
		return "ALIGNED"
	}
	return "MISALIGNED"
}

func severityIcon(s Severity) string {
	switch s {
	case SeverityCritical:
		return "[!]"
	case SeverityHigh:
		return "[H]"
	case SeverityMedium:
		return "[M]"
	case SeverityLow:
		return "[L]"
	default:
		return "[i]"
	}
}
