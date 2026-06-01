// Package drift provides spec-to-code drift detection.
//
// Drift detection compares the reconciled spec.md with the actual codebase
// to identify requirements that haven't been implemented, code that exists
// without spec coverage, and mismatches between spec and implementation.
package drift

import (
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

// DriftType categorizes the type of drift detected.
type DriftType string

const (
	// DriftUnimplemented indicates a spec requirement not found in code.
	DriftUnimplemented DriftType = "unimplemented"

	// DriftUndocumented indicates code that exists without spec coverage.
	DriftUndocumented DriftType = "undocumented"

	// DriftMismatch indicates both spec and code exist but differ.
	DriftMismatch DriftType = "mismatch"
)

// Severity indicates how critical the drift is.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Category groups drift items by type.
type Category string

const (
	CategoryAPI    Category = "api"
	CategoryData   Category = "data"
	CategoryUI     Category = "ui"
	CategoryInfra  Category = "infra"
	CategoryConfig Category = "config"
	CategoryOther  Category = "other"
)

// DriftItem represents a single drift finding.
type DriftItem struct {
	ID          string    `json:"id"`
	Type        DriftType `json:"type"`
	Severity    Severity  `json:"severity"`
	Category    Category  `json:"category"`
	Description string    `json:"description"`
	SpecRef     string    `json:"spec_ref,omitempty"` // Reference to spec section
	CodeRef     string    `json:"code_ref,omitempty"` // Reference to code file:line
	Suggestion  string    `json:"suggestion,omitempty"`
}

// DriftReport contains all drift findings for a project.
type DriftReport struct {
	Project     string       `json:"project"`
	GeneratedAt time.Time    `json:"generated_at"`
	Items       []DriftItem  `json:"items"`
	Summary     DriftSummary `json:"summary"`
}

// DriftSummary provides aggregate drift statistics.
type DriftSummary struct {
	TotalItems    int               `json:"total_items"`
	ByType        map[DriftType]int `json:"by_type"`
	BySeverity    map[Severity]int  `json:"by_severity"`
	ByCategory    map[Category]int  `json:"by_category"`
	CriticalCount int               `json:"critical_count"`
	HighCount     int               `json:"high_count"`
	HasBlockers   bool              `json:"has_blockers"`
}

// Detector performs drift detection.
type Detector struct{}

// NewDetector creates a new drift detector.
func NewDetector() *Detector {
	return &Detector{}
}

// DetectOptions configures drift detection.
type DetectOptions struct {
	MinSeverity  Severity    // Only report items at or above this severity
	Categories   []Category  // Only report items in these categories (empty = all)
	ExcludeTypes []DriftType // Exclude these drift types
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() DetectOptions {
	return DetectOptions{
		MinSeverity: SeverityLow,
	}
}

// Detect compares spec content against codebase context.
func (d *Detector) Detect(specContent string, ctx *context.AggregatedContext, opts DetectOptions) (*DriftReport, error) {
	analyzer := &Analyzer{}

	// Extract requirements from spec
	requirements := analyzer.ExtractRequirements(specContent)

	// Analyze context for implementations
	implementations := analyzer.ExtractImplementations(ctx)

	// Compare and find drift
	items := analyzer.Compare(requirements, implementations)

	// Filter based on options
	filtered := filterItems(items, opts)

	// Calculate summary
	summary := calculateSummary(filtered)

	projectName := ""
	if ctx != nil {
		projectName = ctx.Project
	}

	return &DriftReport{
		Project:     projectName,
		GeneratedAt: time.Now(),
		Items:       filtered,
		Summary:     summary,
	}, nil
}

// filterItems filters drift items based on options.
func filterItems(items []DriftItem, opts DetectOptions) []DriftItem {
	var result []DriftItem

	for _, item := range items {
		// Check severity threshold
		if !meetsMinSeverity(item.Severity, opts.MinSeverity) {
			continue
		}

		// Check category filter
		if len(opts.Categories) > 0 && !containsCategory(opts.Categories, item.Category) {
			continue
		}

		// Check exclusions
		if containsDriftType(opts.ExcludeTypes, item.Type) {
			continue
		}

		result = append(result, item)
	}

	return result
}

// meetsMinSeverity checks if severity meets the minimum threshold.
func meetsMinSeverity(actual, minimum Severity) bool {
	severityOrder := map[Severity]int{
		SeverityCritical: 4,
		SeverityHigh:     3,
		SeverityMedium:   2,
		SeverityLow:      1,
	}
	return severityOrder[actual] >= severityOrder[minimum]
}

func containsCategory(categories []Category, target Category) bool {
	for _, c := range categories {
		if c == target {
			return true
		}
	}
	return false
}

func containsDriftType(types []DriftType, target DriftType) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}

// calculateSummary computes aggregate statistics.
func calculateSummary(items []DriftItem) DriftSummary {
	summary := DriftSummary{
		TotalItems: len(items),
		ByType:     make(map[DriftType]int),
		BySeverity: make(map[Severity]int),
		ByCategory: make(map[Category]int),
	}

	for _, item := range items {
		summary.ByType[item.Type]++
		summary.BySeverity[item.Severity]++
		summary.ByCategory[item.Category]++

		switch item.Severity {
		case SeverityCritical:
			summary.CriticalCount++
			summary.HasBlockers = true
		case SeverityHigh:
			summary.HighCount++
			summary.HasBlockers = true
		}
	}

	return summary
}

// HasDrift returns true if any drift was detected.
func (r *DriftReport) HasDrift() bool {
	return len(r.Items) > 0
}

// HasBlockers returns true if there are critical or high severity items.
func (r *DriftReport) HasBlockers() bool {
	return r.Summary.HasBlockers
}

// FilterBySeverity returns items at or above the given severity.
func (r *DriftReport) FilterBySeverity(minSeverity Severity) []DriftItem {
	return filterItems(r.Items, DetectOptions{MinSeverity: minSeverity})
}

// FilterByCategory returns items in the given categories.
func (r *DriftReport) FilterByCategory(categories ...Category) []DriftItem {
	return filterItems(r.Items, DetectOptions{Categories: categories})
}
