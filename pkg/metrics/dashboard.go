// Package metrics provides dashboard rendering for project metrics.
package metrics

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
)

// OutputFormat specifies the dashboard output format.
type OutputFormat string

const (
	FormatJSON     OutputFormat = "json"
	FormatHTML     OutputFormat = "html"
	FormatTerminal OutputFormat = "terminal"
	FormatMarkdown OutputFormat = "markdown"
)

// Dashboard renders project metrics in various formats.
type Dashboard struct {
	metrics *ProjectMetrics
}

// NewDashboard creates a new dashboard from project metrics.
func NewDashboard(metrics *ProjectMetrics) *Dashboard {
	return &Dashboard{metrics: metrics}
}

// Render outputs the dashboard in the specified format.
func (d *Dashboard) Render(w io.Writer, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return d.renderJSON(w)
	case FormatHTML:
		return d.renderHTML(w)
	case FormatTerminal:
		return d.renderTerminal(w)
	case FormatMarkdown:
		return d.renderMarkdown(w)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// renderJSON outputs metrics as formatted JSON.
func (d *Dashboard) renderJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(d.metrics)
}

// renderTerminal outputs metrics for terminal display.
func (d *Dashboard) renderTerminal(w io.Writer) error {
	m := d.metrics

	// Header
	fmt.Fprintf(w, "╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(w, "║              VisionSpec Project Metrics                      ║\n")
	fmt.Fprintf(w, "╠══════════════════════════════════════════════════════════════╣\n")
	fmt.Fprintf(w, "║  Project: %-50s  ║\n", truncate(m.Project, 50))
	fmt.Fprintf(w, "║  Generated: %-48s  ║\n", m.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "╚══════════════════════════════════════════════════════════════╝\n\n")

	// Health Score
	healthBar := d.healthBar(m.HealthScore)
	fmt.Fprintf(w, "Health Score: %s %.1f/100 [%s]\n\n", healthBar, m.HealthScore, m.HealthStatus)

	// Summary
	fmt.Fprintf(w, "┌─────────────────────────────────────────────────────────────┐\n")
	fmt.Fprintf(w, "│ SUMMARY                                                     │\n")
	fmt.Fprintf(w, "├─────────────────────────────────────────────────────────────┤\n")
	fmt.Fprintf(w, "│  Total Specs:     %-6d    Overall Score:    %6.1f%%        │\n",
		m.Summary.TotalSpecs, m.Summary.OverallScore)
	fmt.Fprintf(w, "│  Evaluated:       %-6d    Quality Score:    %6.1f%%        │\n",
		m.Summary.EvaluatedSpecs, m.Summary.QualityScore)
	fmt.Fprintf(w, "│  Passing:         %-6d    Readiness Score:  %6.1f%%        │\n",
		m.Summary.PassingEvals, m.Summary.ReadinessScore)
	fmt.Fprintf(w, "└─────────────────────────────────────────────────────────────┘\n\n")

	// Evaluation Metrics
	if m.Eval != nil {
		fmt.Fprintf(w, "┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Fprintf(w, "│ EVALUATION METRICS                                          │\n")
		fmt.Fprintf(w, "├─────────────────────────────────────────────────────────────┤\n")
		fmt.Fprintf(w, "│  Total Evaluations: %-6d    Average Score: %5.2f          │\n",
			m.Eval.TotalEvaluations, m.Eval.AverageScore)
		fmt.Fprintf(w, "│  Passed:            %-6d    Failed:        %-6d          │\n",
			m.Eval.PassCount, m.Eval.FailCount)
		fmt.Fprintf(w, "│  Findings:          %-6d                                   │\n",
			m.Eval.FindingsCount)

		if len(m.Eval.BySeverity) > 0 {
			fmt.Fprintf(w, "│                                                             │\n")
			fmt.Fprintf(w, "│  By Severity:                                               │\n")
			for sev, count := range m.Eval.BySeverity {
				icon := severityIcon(sev)
				fmt.Fprintf(w, "│    %s %-10s: %-6d                                    │\n",
					icon, sev, count)
			}
		}
		fmt.Fprintf(w, "└─────────────────────────────────────────────────────────────┘\n\n")
	}

	// Reconciliation Metrics
	if m.Reconcile != nil {
		fmt.Fprintf(w, "┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Fprintf(w, "│ RECONCILIATION METRICS                                      │\n")
		fmt.Fprintf(w, "├─────────────────────────────────────────────────────────────┤\n")
		fmt.Fprintf(w, "│  Total Reconciliations: %-6d                               │\n",
			m.Reconcile.TotalReconciliations)
		fmt.Fprintf(w, "│  Successful:            %-6d    Conflicts: %-6d          │\n",
			m.Reconcile.SuccessCount, m.Reconcile.ConflictCount)
		fmt.Fprintf(w, "│  Specs Included:        %-6d    Tasks Gen: %-6d          │\n",
			m.Reconcile.SpecsIncluded, m.Reconcile.TasksGenerated)
		if m.Reconcile.LastReconcile != nil {
			fmt.Fprintf(w, "│  Last Reconcile: %-42s │\n",
				m.Reconcile.LastReconcile.Format("2006-01-02 15:04:05"))
		}
		fmt.Fprintf(w, "└─────────────────────────────────────────────────────────────┘\n\n")
	}

	// Alignment Metrics
	if m.Align != nil {
		fmt.Fprintf(w, "┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Fprintf(w, "│ ALIGNMENT METRICS                                           │\n")
		fmt.Fprintf(w, "├─────────────────────────────────────────────────────────────┤\n")
		fmt.Fprintf(w, "│  Alignment Score:  %6.1f%%    Coverage: %6.1f%%             │\n",
			m.Align.AlignmentScore, m.Align.CoveragePercent)
		fmt.Fprintf(w, "│  Discrepancies:    %-6d                                    │\n",
			m.Align.DiscrepancyCount)
		fmt.Fprintf(w, "│                                                             │\n")
		fmt.Fprintf(w, "│  By Severity:                                               │\n")
		fmt.Fprintf(w, "│    🔴 Critical: %-6d    🟠 High:   %-6d                   │\n",
			m.Align.CriticalCount, m.Align.HighCount)
		fmt.Fprintf(w, "│    🟡 Medium:   %-6d    🟢 Low:    %-6d                   │\n",
			m.Align.MediumCount, m.Align.LowCount)
		fmt.Fprintf(w, "│                                                             │\n")
		fmt.Fprintf(w, "│  Missing Features:    %-6d                                 │\n",
			m.Align.MissingFeatures)
		fmt.Fprintf(w, "│  Undocumented Code:   %-6d                                 │\n",
			m.Align.UndocumentedCode)
		fmt.Fprintf(w, "└─────────────────────────────────────────────────────────────┘\n\n")
	}

	// Drift Metrics
	if m.Drift != nil {
		fmt.Fprintf(w, "┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Fprintf(w, "│ DRIFT METRICS                                               │\n")
		fmt.Fprintf(w, "├─────────────────────────────────────────────────────────────┤\n")
		driftStatus := "No drift detected"
		if m.Drift.HasDrift {
			driftStatus = "Drift detected"
		}
		fmt.Fprintf(w, "│  Status: %-51s │\n", driftStatus)
		fmt.Fprintf(w, "│  Drift Score:  %6.1f    Items: %-6d                       │\n",
			m.Drift.DriftScore, m.Drift.ItemCount)
		fmt.Fprintf(w, "│  Critical:     %-6d    High:  %-6d                       │\n",
			m.Drift.CriticalCount, m.Drift.HighCount)
		fmt.Fprintf(w, "│  Trend: %-52s │\n", m.Drift.TrendDirection)
		fmt.Fprintf(w, "└─────────────────────────────────────────────────────────────┘\n")
	}

	return nil
}

// healthBar creates a visual health bar.
func (d *Dashboard) healthBar(score float64) string {
	filled := int(score / 10)
	empty := 10 - filled

	var bar strings.Builder
	bar.WriteString("[")
	for i := 0; i < filled; i++ {
		if score >= 75 {
			bar.WriteString("█")
		} else if score >= 50 {
			bar.WriteString("▓")
		} else {
			bar.WriteString("░")
		}
	}
	for i := 0; i < empty; i++ {
		bar.WriteString("·")
	}
	bar.WriteString("]")
	return bar.String()
}

// severityIcon returns an icon for a severity level.
func severityIcon(severity string) string {
	switch strings.ToLower(severity) {
	case "critical":
		return "🔴"
	case "high":
		return "🟠"
	case "medium":
		return "🟡"
	case "low":
		return "🟢"
	default:
		return "⚪"
	}
}

// truncate shortens a string to maxLen.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// renderMarkdown outputs metrics as Markdown.
func (d *Dashboard) renderMarkdown(w io.Writer) error {
	m := d.metrics

	fmt.Fprintf(w, "# VisionSpec Project Metrics\n\n")
	fmt.Fprintf(w, "**Project:** %s  \n", m.Project)
	fmt.Fprintf(w, "**Generated:** %s\n\n", m.GeneratedAt.Format("2006-01-02 15:04:05"))

	// Health Score
	fmt.Fprintf(w, "## Health Score\n\n")
	fmt.Fprintf(w, "| Score | Status |\n")
	fmt.Fprintf(w, "|-------|--------|\n")
	fmt.Fprintf(w, "| %.1f/100 | %s |\n\n", m.HealthScore, m.HealthStatus)

	// Summary
	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "| Metric | Value |\n")
	fmt.Fprintf(w, "|--------|-------|\n")
	fmt.Fprintf(w, "| Total Specs | %d |\n", m.Summary.TotalSpecs)
	fmt.Fprintf(w, "| Evaluated Specs | %d |\n", m.Summary.EvaluatedSpecs)
	fmt.Fprintf(w, "| Passing Evals | %d |\n", m.Summary.PassingEvals)
	fmt.Fprintf(w, "| Overall Score | %.1f%% |\n", m.Summary.OverallScore)
	fmt.Fprintf(w, "| Quality Score | %.1f%% |\n", m.Summary.QualityScore)
	fmt.Fprintf(w, "| Readiness Score | %.1f%% |\n\n", m.Summary.ReadinessScore)

	// Evaluation
	if m.Eval != nil {
		fmt.Fprintf(w, "## Evaluation Metrics\n\n")
		fmt.Fprintf(w, "| Metric | Value |\n")
		fmt.Fprintf(w, "|--------|-------|\n")
		fmt.Fprintf(w, "| Total Evaluations | %d |\n", m.Eval.TotalEvaluations)
		fmt.Fprintf(w, "| Pass Count | %d |\n", m.Eval.PassCount)
		fmt.Fprintf(w, "| Fail Count | %d |\n", m.Eval.FailCount)
		fmt.Fprintf(w, "| Average Score | %.2f |\n", m.Eval.AverageScore)
		fmt.Fprintf(w, "| Total Findings | %d |\n\n", m.Eval.FindingsCount)

		if len(m.Eval.BySeverity) > 0 {
			fmt.Fprintf(w, "### Findings by Severity\n\n")
			fmt.Fprintf(w, "| Severity | Count |\n")
			fmt.Fprintf(w, "|----------|-------|\n")
			for sev, count := range m.Eval.BySeverity {
				fmt.Fprintf(w, "| %s | %d |\n", sev, count)
			}
			fmt.Fprintf(w, "\n")
		}
	}

	// Reconciliation
	if m.Reconcile != nil {
		fmt.Fprintf(w, "## Reconciliation Metrics\n\n")
		fmt.Fprintf(w, "| Metric | Value |\n")
		fmt.Fprintf(w, "|--------|-------|\n")
		fmt.Fprintf(w, "| Total Reconciliations | %d |\n", m.Reconcile.TotalReconciliations)
		fmt.Fprintf(w, "| Success Count | %d |\n", m.Reconcile.SuccessCount)
		fmt.Fprintf(w, "| Conflict Count | %d |\n", m.Reconcile.ConflictCount)
		fmt.Fprintf(w, "| Specs Included | %d |\n", m.Reconcile.SpecsIncluded)
		fmt.Fprintf(w, "| Tasks Generated | %d |\n", m.Reconcile.TasksGenerated)
		if m.Reconcile.LastReconcile != nil {
			fmt.Fprintf(w, "| Last Reconcile | %s |\n",
				m.Reconcile.LastReconcile.Format("2006-01-02 15:04:05"))
		}
		fmt.Fprintf(w, "\n")
	}

	// Alignment
	if m.Align != nil {
		fmt.Fprintf(w, "## Alignment Metrics\n\n")
		fmt.Fprintf(w, "| Metric | Value |\n")
		fmt.Fprintf(w, "|--------|-------|\n")
		fmt.Fprintf(w, "| Alignment Score | %.1f%% |\n", m.Align.AlignmentScore)
		fmt.Fprintf(w, "| Coverage | %.1f%% |\n", m.Align.CoveragePercent)
		fmt.Fprintf(w, "| Total Discrepancies | %d |\n", m.Align.DiscrepancyCount)
		fmt.Fprintf(w, "| Critical | %d |\n", m.Align.CriticalCount)
		fmt.Fprintf(w, "| High | %d |\n", m.Align.HighCount)
		fmt.Fprintf(w, "| Medium | %d |\n", m.Align.MediumCount)
		fmt.Fprintf(w, "| Low | %d |\n", m.Align.LowCount)
		fmt.Fprintf(w, "| Missing Features | %d |\n", m.Align.MissingFeatures)
		fmt.Fprintf(w, "| Undocumented Code | %d |\n\n", m.Align.UndocumentedCode)
	}

	// Drift
	if m.Drift != nil {
		fmt.Fprintf(w, "## Drift Metrics\n\n")
		fmt.Fprintf(w, "| Metric | Value |\n")
		fmt.Fprintf(w, "|--------|-------|\n")
		fmt.Fprintf(w, "| Has Drift | %v |\n", m.Drift.HasDrift)
		fmt.Fprintf(w, "| Drift Score | %.1f |\n", m.Drift.DriftScore)
		fmt.Fprintf(w, "| Item Count | %d |\n", m.Drift.ItemCount)
		fmt.Fprintf(w, "| Critical | %d |\n", m.Drift.CriticalCount)
		fmt.Fprintf(w, "| High | %d |\n", m.Drift.HighCount)
		fmt.Fprintf(w, "| Trend | %s |\n\n", m.Drift.TrendDirection)
	}

	return nil
}

// HTML template for dashboard.
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VisionSpec Metrics - {{.Project}}</title>
    <style>
        :root {
            --bg: #1a1b26;
            --surface: #24283b;
            --primary: #7aa2f7;
            --success: #9ece6a;
            --warning: #e0af68;
            --error: #f7768e;
            --text: #c0caf5;
            --text-muted: #565f89;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            padding: 2rem;
            line-height: 1.6;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        header {
            text-align: center;
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid var(--surface);
        }
        h1 { color: var(--primary); font-size: 2rem; }
        .meta { color: var(--text-muted); margin-top: 0.5rem; }
        .health-score {
            background: var(--surface);
            border-radius: 12px;
            padding: 2rem;
            text-align: center;
            margin-bottom: 2rem;
        }
        .health-value {
            font-size: 4rem;
            font-weight: bold;
        }
        .health-value.healthy { color: var(--success); }
        .health-value.warning { color: var(--warning); }
        .health-value.critical { color: var(--error); }
        .health-bar {
            width: 100%;
            height: 12px;
            background: var(--bg);
            border-radius: 6px;
            margin-top: 1rem;
            overflow: hidden;
        }
        .health-bar-fill {
            height: 100%;
            border-radius: 6px;
            transition: width 0.3s ease;
        }
        .health-bar-fill.healthy { background: var(--success); }
        .health-bar-fill.warning { background: var(--warning); }
        .health-bar-fill.critical { background: var(--error); }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 1.5rem;
        }
        .card {
            background: var(--surface);
            border-radius: 12px;
            padding: 1.5rem;
        }
        .card h2 {
            color: var(--primary);
            font-size: 1.1rem;
            margin-bottom: 1rem;
            padding-bottom: 0.5rem;
            border-bottom: 1px solid var(--bg);
        }
        .metric-row {
            display: flex;
            justify-content: space-between;
            padding: 0.5rem 0;
            border-bottom: 1px solid var(--bg);
        }
        .metric-row:last-child { border-bottom: none; }
        .metric-label { color: var(--text-muted); }
        .metric-value { font-weight: 600; }
        .severity-grid {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 0.5rem;
            margin-top: 1rem;
        }
        .severity-item {
            text-align: center;
            padding: 0.75rem;
            background: var(--bg);
            border-radius: 8px;
        }
        .severity-item.critical { border-left: 3px solid var(--error); }
        .severity-item.high { border-left: 3px solid var(--warning); }
        .severity-item.medium { border-left: 3px solid #bb9af7; }
        .severity-item.low { border-left: 3px solid var(--success); }
        .severity-count { font-size: 1.5rem; font-weight: bold; }
        .severity-label { font-size: 0.75rem; color: var(--text-muted); }
        footer {
            text-align: center;
            margin-top: 2rem;
            padding-top: 1rem;
            border-top: 1px solid var(--surface);
            color: var(--text-muted);
            font-size: 0.875rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>VisionSpec Metrics</h1>
            <div class="meta">
                <strong>{{.Project}}</strong> &middot;
                Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}
            </div>
        </header>

        <div class="health-score">
            <div class="health-value {{.HealthStatus}}">{{printf "%.0f" .HealthScore}}</div>
            <div>Health Score</div>
            <div class="health-bar">
                <div class="health-bar-fill {{.HealthStatus}}" style="width: {{.HealthScore}}%"></div>
            </div>
        </div>

        <div class="grid">
            <div class="card">
                <h2>Summary</h2>
                <div class="metric-row">
                    <span class="metric-label">Total Specs</span>
                    <span class="metric-value">{{.Summary.TotalSpecs}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Evaluated Specs</span>
                    <span class="metric-value">{{.Summary.EvaluatedSpecs}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Passing Evaluations</span>
                    <span class="metric-value">{{.Summary.PassingEvals}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Overall Score</span>
                    <span class="metric-value">{{printf "%.1f%%" .Summary.OverallScore}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Quality Score</span>
                    <span class="metric-value">{{printf "%.1f%%" .Summary.QualityScore}}</span>
                </div>
            </div>

            {{if .Eval}}
            <div class="card">
                <h2>Evaluations</h2>
                <div class="metric-row">
                    <span class="metric-label">Total Evaluations</span>
                    <span class="metric-value">{{.Eval.TotalEvaluations}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Passed</span>
                    <span class="metric-value">{{.Eval.PassCount}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Failed</span>
                    <span class="metric-value">{{.Eval.FailCount}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Average Score</span>
                    <span class="metric-value">{{printf "%.2f" .Eval.AverageScore}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Total Findings</span>
                    <span class="metric-value">{{.Eval.FindingsCount}}</span>
                </div>
            </div>
            {{end}}

            {{if .Reconcile}}
            <div class="card">
                <h2>Reconciliation</h2>
                <div class="metric-row">
                    <span class="metric-label">Total Reconciliations</span>
                    <span class="metric-value">{{.Reconcile.TotalReconciliations}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Successful</span>
                    <span class="metric-value">{{.Reconcile.SuccessCount}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Conflicts</span>
                    <span class="metric-value">{{.Reconcile.ConflictCount}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Specs Included</span>
                    <span class="metric-value">{{.Reconcile.SpecsIncluded}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Tasks Generated</span>
                    <span class="metric-value">{{.Reconcile.TasksGenerated}}</span>
                </div>
            </div>
            {{end}}

            {{if .Align}}
            <div class="card">
                <h2>Alignment</h2>
                <div class="metric-row">
                    <span class="metric-label">Alignment Score</span>
                    <span class="metric-value">{{printf "%.1f%%" .Align.AlignmentScore}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Coverage</span>
                    <span class="metric-value">{{printf "%.1f%%" .Align.CoveragePercent}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Total Discrepancies</span>
                    <span class="metric-value">{{.Align.DiscrepancyCount}}</span>
                </div>
                <div class="severity-grid">
                    <div class="severity-item critical">
                        <div class="severity-count">{{.Align.CriticalCount}}</div>
                        <div class="severity-label">Critical</div>
                    </div>
                    <div class="severity-item high">
                        <div class="severity-count">{{.Align.HighCount}}</div>
                        <div class="severity-label">High</div>
                    </div>
                    <div class="severity-item medium">
                        <div class="severity-count">{{.Align.MediumCount}}</div>
                        <div class="severity-label">Medium</div>
                    </div>
                    <div class="severity-item low">
                        <div class="severity-count">{{.Align.LowCount}}</div>
                        <div class="severity-label">Low</div>
                    </div>
                </div>
            </div>
            {{end}}

            {{if .Drift}}
            <div class="card">
                <h2>Drift Detection</h2>
                <div class="metric-row">
                    <span class="metric-label">Status</span>
                    <span class="metric-value">{{if .Drift.HasDrift}}Drift Detected{{else}}No Drift{{end}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Drift Score</span>
                    <span class="metric-value">{{printf "%.1f" .Drift.DriftScore}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Drift Items</span>
                    <span class="metric-value">{{.Drift.ItemCount}}</span>
                </div>
                <div class="metric-row">
                    <span class="metric-label">Trend</span>
                    <span class="metric-value">{{.Drift.TrendDirection}}</span>
                </div>
            </div>
            {{end}}
        </div>

        <footer>
            Generated by VisionSpec
        </footer>
    </div>
</body>
</html>`

// renderHTML outputs metrics as an HTML dashboard.
func (d *Dashboard) renderHTML(w io.Writer) error {
	tmpl, err := template.New("dashboard").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	return tmpl.Execute(w, d.metrics)
}
