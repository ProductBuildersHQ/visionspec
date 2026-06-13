package align

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// CurrentTruth represents the current state of implementation relative to spec.
type CurrentTruth struct {
	Project        string         `json:"project"`
	GeneratedAt    time.Time      `json:"generated_at"`
	SpecVersion    string         `json:"spec_version,omitempty"`
	AlignmentScore float64        `json:"alignment_score"`
	Status         TruthStatus    `json:"status"`
	Sections       []TruthSection `json:"sections"`
	Summary        TruthSummary   `json:"summary"`
	NextActions    []string       `json:"next_actions,omitempty"`
}

// TruthStatus indicates overall alignment status.
type TruthStatus string

const (
	TruthStatusAligned  TruthStatus = "aligned"  // Spec and code match
	TruthStatusDrifted  TruthStatus = "drifted"  // Minor deviations
	TruthStatusDiverged TruthStatus = "diverged" // Significant misalignment
	TruthStatusUnknown  TruthStatus = "unknown"  // Unable to determine
)

// TruthSection represents a spec section's truth status.
type TruthSection struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"` // Section path in spec
	Status      TruthStatus `json:"status"`
	Score       float64     `json:"score"`
	Implemented []TruthItem `json:"implemented,omitempty"`
	Partial     []TruthItem `json:"partial,omitempty"`
	Missing     []TruthItem `json:"missing,omitempty"`
	Notes       string      `json:"notes,omitempty"`
}

// TruthItem represents a single requirement's truth status.
type TruthItem struct {
	Requirement string   `json:"requirement"`
	Status      string   `json:"status"` // "done", "partial", "missing", "diverged"
	CodeRefs    []string `json:"code_refs,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// TruthSummary provides aggregate truth metrics.
type TruthSummary struct {
	TotalRequirements int     `json:"total_requirements"`
	Implemented       int     `json:"implemented"`
	Partial           int     `json:"partial"`
	Missing           int     `json:"missing"`
	CoveragePercent   float64 `json:"coverage_percent"`
	AlignmentPercent  float64 `json:"alignment_percent"`
}

// GenerateCurrentTruth creates a CurrentTruth document from an alignment result.
func GenerateCurrentTruth(result *AlignmentResult) *CurrentTruth {
	truth := &CurrentTruth{
		Project:        result.Project,
		GeneratedAt:    time.Now(),
		AlignmentScore: result.Summary.AlignmentScore,
	}

	// Determine overall status
	if result.Summary.IsAligned {
		truth.Status = TruthStatusAligned
	} else if result.Summary.CriticalCount > 0 {
		truth.Status = TruthStatusDiverged
	} else {
		truth.Status = TruthStatusDrifted
	}

	// Group discrepancies by section/category
	sectionMap := make(map[Category]*TruthSection)

	for _, d := range result.Discrepancies {
		section, ok := sectionMap[d.Category]
		if !ok {
			section = &TruthSection{
				Name:   string(d.Category),
				Path:   string(d.Category),
				Status: TruthStatusAligned,
			}
			sectionMap[d.Category] = section
		}

		item := TruthItem{
			Requirement: d.Description,
			Notes:       d.Suggestion,
		}

		if d.CodeRef != "" {
			item.CodeRefs = []string{d.CodeRef}
		}

		switch d.Type {
		case DiscrepancyMissingFeature:
			item.Status = "missing"
			section.Missing = append(section.Missing, item)
			section.Status = worseStatus(section.Status, TruthStatusDrifted)
		case DiscrepancyPartialImplementation:
			item.Status = "partial"
			section.Partial = append(section.Partial, item)
			section.Status = worseStatus(section.Status, TruthStatusDrifted)
		case DiscrepancyDiverged, DiscrepancyBehaviorMismatch:
			item.Status = "diverged"
			section.Missing = append(section.Missing, item)
			section.Status = TruthStatusDiverged
		case DiscrepancyUndocumentedCode:
			// Track but don't affect status negatively
			item.Status = "undocumented"
			section.Notes = addNote(section.Notes, fmt.Sprintf("Undocumented: %s", d.Description))
		}
	}

	// Convert map to slice
	for _, section := range sectionMap {
		// Calculate section score
		total := len(section.Implemented) + len(section.Partial) + len(section.Missing)
		if total > 0 {
			section.Score = float64(len(section.Implemented)+len(section.Partial)/2) / float64(total)
		} else {
			section.Score = 1.0
		}
		truth.Sections = append(truth.Sections, *section)
	}

	// Calculate summary
	truth.Summary = TruthSummary{
		TotalRequirements: result.Coverage.TotalRequirements,
		Implemented:       result.Coverage.ImplementedCount,
		Partial:           result.Coverage.PartialCount,
		Missing:           result.Coverage.MissingCount,
		CoveragePercent:   result.Coverage.CoveragePercentage,
		AlignmentPercent:  result.Summary.AlignmentScore * 100,
	}

	// Generate next actions
	truth.NextActions = generateNextActions(result)

	return truth
}

// generateNextActions creates actionable recommendations.
func generateNextActions(result *AlignmentResult) []string {
	var actions []string

	// Priority: Critical discrepancies first
	criticals := result.FilterBySeverity(SeverityCritical)
	for _, d := range criticals {
		if len(actions) < 3 {
			actions = append(actions, fmt.Sprintf("[CRITICAL] %s", d.Suggestion))
		}
	}

	// Then high severity
	highs := result.FilterBySeverity(SeverityHigh)
	for _, d := range highs {
		if len(actions) < 5 && d.Severity == SeverityHigh {
			actions = append(actions, fmt.Sprintf("[HIGH] %s", d.Suggestion))
		}
	}

	// General recommendations based on coverage
	if result.Coverage.CoveragePercentage < 50 {
		actions = append(actions, "Consider prioritizing feature implementation - coverage is below 50%")
	}

	if result.Coverage.UndocumentedCount > 5 {
		actions = append(actions, fmt.Sprintf("Review %d undocumented code items for spec coverage", result.Coverage.UndocumentedCount))
	}

	if len(actions) == 0 {
		actions = append(actions, "No immediate actions required - implementation aligns with spec")
	}

	return actions
}

func worseStatus(a, b TruthStatus) TruthStatus {
	order := map[TruthStatus]int{
		TruthStatusAligned:  0,
		TruthStatusDrifted:  1,
		TruthStatusDiverged: 2,
		TruthStatusUnknown:  3,
	}
	if order[b] > order[a] {
		return b
	}
	return a
}

func addNote(existing, note string) string {
	if existing == "" {
		return note
	}
	return existing + "; " + note
}

// RenderMarkdown renders the CurrentTruth as a markdown document.
func (t *CurrentTruth) RenderMarkdown() (string, error) {
	tmpl, err := template.New("truth").Funcs(template.FuncMap{
		"statusEmoji": func(s TruthStatus) string {
			switch s {
			case TruthStatusAligned:
				return "+"
			case TruthStatusDrifted:
				return "~"
			case TruthStatusDiverged:
				return "!"
			default:
				return "?"
			}
		},
		"itemStatusEmoji": func(s string) string {
			switch s {
			case "done":
				return "[x]"
			case "partial":
				return "[~]"
			case "missing":
				return "[ ]"
			case "diverged":
				return "[!]"
			default:
				return "[-]"
			}
		},
		"percent": func(f float64) string {
			return fmt.Sprintf("%.1f%%", f)
		},
	}).Parse(currentTruthTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

const currentTruthTemplate = `# Current Truth: {{.Project}}

> Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}
> Status: {{statusEmoji .Status}} {{.Status}}
> Alignment Score: {{percent .AlignmentScore}}

## Summary

| Metric | Value |
|--------|-------|
| Total Requirements | {{.Summary.TotalRequirements}} |
| Implemented | {{.Summary.Implemented}} |
| Partial | {{.Summary.Partial}} |
| Missing | {{.Summary.Missing}} |
| Coverage | {{percent .Summary.CoveragePercent}} |
| Alignment | {{percent .Summary.AlignmentPercent}} |

## Next Actions

{{range .NextActions}}
- {{.}}
{{end}}

{{range .Sections}}
## {{.Name}}

Status: {{statusEmoji .Status}} {{.Status}} (Score: {{percent .Score}})

{{if .Implemented}}
### Implemented

{{range .Implemented}}
- {{itemStatusEmoji .Status}} {{.Requirement}}{{if .CodeRefs}} ({{range $i, $ref := .CodeRefs}}{{if $i}}, {{end}}` + "`{{$ref}}`" + `{{end}}){{end}}
{{end}}
{{end}}

{{if .Partial}}
### Partial

{{range .Partial}}
- {{itemStatusEmoji .Status}} {{.Requirement}}{{if .Notes}} - {{.Notes}}{{end}}
{{end}}
{{end}}

{{if .Missing}}
### Missing

{{range .Missing}}
- {{itemStatusEmoji .Status}} {{.Requirement}}{{if .Notes}} - {{.Notes}}{{end}}
{{end}}
{{end}}

{{if .Notes}}
> Note: {{.Notes}}
{{end}}

{{end}}
---

*This document reflects the current state of implementation relative to the specification.
Use ` + "`visionspec align`" + ` to regenerate.*
`

// ParseCurrentTruth parses a current-truth.md file back into a struct.
// This is useful for tracking changes over time.
func ParseCurrentTruth(content string) (*CurrentTruth, error) {
	// This is a simplified parser - a full implementation would use
	// a proper markdown parser
	truth := &CurrentTruth{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse project name from title
		if strings.HasPrefix(line, "# Current Truth:") {
			truth.Project = strings.TrimSpace(strings.TrimPrefix(line, "# Current Truth:"))
		}

		// Parse status
		if strings.Contains(line, "Status:") && strings.Contains(line, "aligned") {
			truth.Status = TruthStatusAligned
		} else if strings.Contains(line, "Status:") && strings.Contains(line, "drifted") {
			truth.Status = TruthStatusDrifted
		} else if strings.Contains(line, "Status:") && strings.Contains(line, "diverged") {
			truth.Status = TruthStatusDiverged
		}
	}

	return truth, nil
}
