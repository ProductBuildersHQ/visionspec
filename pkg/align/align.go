// Package align provides post-ship spec-to-reality alignment checking.
//
// Alignment compares the reconciled spec.md against the actual implementation
// to identify discrepancies between what was specified and what was built.
// This enables continuous validation that the shipped product matches the spec.
package align

import (
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

// DiscrepancyType categorizes the type of discrepancy detected.
type DiscrepancyType string

const (
	// DiscrepancyMissingFeature indicates a specified feature not found in code.
	DiscrepancyMissingFeature DiscrepancyType = "missing_feature"

	// DiscrepancyUndocumentedCode indicates code functionality without spec coverage.
	DiscrepancyUndocumentedCode DiscrepancyType = "undocumented_code"

	// DiscrepancyDiverged indicates both spec and code exist but have diverged.
	DiscrepancyDiverged DiscrepancyType = "diverged"

	// DiscrepancyPartialImplementation indicates feature is partially implemented.
	DiscrepancyPartialImplementation DiscrepancyType = "partial_implementation"

	// DiscrepancyBehaviorMismatch indicates behavior differs from spec.
	DiscrepancyBehaviorMismatch DiscrepancyType = "behavior_mismatch"
)

// Severity indicates how critical the discrepancy is.
type Severity string

const (
	SeverityCritical Severity = "critical" // Blocking issue, must fix
	SeverityHigh     Severity = "high"     // Significant deviation
	SeverityMedium   Severity = "medium"   // Notable but not blocking
	SeverityLow      Severity = "low"      // Minor discrepancy
	SeverityInfo     Severity = "info"     // Informational only
)

// Category groups discrepancies by functional area.
type Category string

const (
	CategoryAPI      Category = "api"      // API endpoints, contracts
	CategoryData     Category = "data"     // Data models, schemas
	CategoryUI       Category = "ui"       // User interface
	CategoryBehavior Category = "behavior" // Business logic
	CategoryInfra    Category = "infra"    // Infrastructure
	CategorySecurity Category = "security" // Security features
	CategoryPerf     Category = "perf"     // Performance characteristics
	CategoryOther    Category = "other"    // Uncategorized
)

// Discrepancy represents a single alignment finding.
type Discrepancy struct {
	ID          string          `json:"id"`
	Type        DiscrepancyType `json:"type"`
	Severity    Severity        `json:"severity"`
	Category    Category        `json:"category"`
	Description string          `json:"description"`
	SpecRef     string          `json:"spec_ref,omitempty"` // Reference to spec section/line
	CodeRef     string          `json:"code_ref,omitempty"` // Reference to code file:line
	Expected    string          `json:"expected,omitempty"` // What spec says
	Actual      string          `json:"actual,omitempty"`   // What code does
	Suggestion  string          `json:"suggestion,omitempty"`
	Evidence    []Evidence      `json:"evidence,omitempty"`
}

// Evidence provides supporting details for a discrepancy.
type Evidence struct {
	Type    string `json:"type"`    // "spec_excerpt", "code_snippet", "test_result"
	Content string `json:"content"` // The evidence content
	Source  string `json:"source"`  // Where it came from
}

// AlignmentResult contains the full alignment analysis.
type AlignmentResult struct {
	Project       string            `json:"project"`
	GeneratedAt   time.Time         `json:"generated_at"`
	SpecPath      string            `json:"spec_path"`
	ContextSource string            `json:"context_source,omitempty"`
	Discrepancies []Discrepancy     `json:"discrepancies"`
	Summary       AlignmentSummary  `json:"summary"`
	Coverage      AlignmentCoverage `json:"coverage"`
	Metadata      map[string]any    `json:"metadata,omitempty"`
}

// AlignmentSummary provides aggregate statistics.
type AlignmentSummary struct {
	TotalDiscrepancies int                     `json:"total_discrepancies"`
	ByType             map[DiscrepancyType]int `json:"by_type"`
	BySeverity         map[Severity]int        `json:"by_severity"`
	ByCategory         map[Category]int        `json:"by_category"`
	CriticalCount      int                     `json:"critical_count"`
	HighCount          int                     `json:"high_count"`
	AlignmentScore     float64                 `json:"alignment_score"` // 0.0 to 1.0
	IsAligned          bool                    `json:"is_aligned"`      // True if no critical/high issues
}

// AlignmentCoverage tracks spec coverage metrics.
type AlignmentCoverage struct {
	TotalRequirements  int     `json:"total_requirements"`
	ImplementedCount   int     `json:"implemented_count"`
	PartialCount       int     `json:"partial_count"`
	MissingCount       int     `json:"missing_count"`
	CoveragePercentage float64 `json:"coverage_percentage"`
	UndocumentedCount  int     `json:"undocumented_count"` // Code without spec
}

// Aligner performs alignment checking.
type Aligner struct {
	comparator *Comparator
}

// NewAligner creates a new alignment checker.
func NewAligner() *Aligner {
	return &Aligner{
		comparator: NewComparator(),
	}
}

// AlignOptions configures alignment checking.
type AlignOptions struct {
	MinSeverity      Severity   // Only report items at or above this severity
	Categories       []Category // Only report items in these categories (empty = all)
	ExcludeTypes     []DiscrepancyType
	IncludeEvidence  bool // Include evidence snippets
	MaxDiscrepancies int  // Limit results (0 = unlimited)
}

// DefaultAlignOptions returns sensible defaults.
func DefaultAlignOptions() AlignOptions {
	return AlignOptions{
		MinSeverity:     SeverityLow,
		IncludeEvidence: true,
	}
}

// Align compares spec content against codebase context.
func (a *Aligner) Align(specContent string, ctx *context.AggregatedContext, opts AlignOptions) (*AlignmentResult, error) {
	// Extract requirements from spec
	requirements := a.comparator.ExtractRequirements(specContent)

	// Extract implementations from context
	implementations := a.comparator.ExtractImplementations(ctx)

	// Compare and find discrepancies
	discrepancies := a.comparator.Compare(requirements, implementations, opts.IncludeEvidence)

	// Filter based on options
	filtered := filterDiscrepancies(discrepancies, opts)

	// Calculate summary and coverage
	summary := calculateAlignmentSummary(filtered)
	coverage := calculateCoverage(requirements, implementations, discrepancies)

	projectName := ""
	if ctx != nil {
		projectName = ctx.Project
	}

	return &AlignmentResult{
		Project:       projectName,
		GeneratedAt:   time.Now(),
		Discrepancies: filtered,
		Summary:       summary,
		Coverage:      coverage,
	}, nil
}

// filterDiscrepancies filters based on options.
func filterDiscrepancies(items []Discrepancy, opts AlignOptions) []Discrepancy {
	var result []Discrepancy

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
		if containsDiscrepancyType(opts.ExcludeTypes, item.Type) {
			continue
		}

		result = append(result, item)

		// Check limit
		if opts.MaxDiscrepancies > 0 && len(result) >= opts.MaxDiscrepancies {
			break
		}
	}

	return result
}

// meetsMinSeverity checks if severity meets the minimum threshold.
func meetsMinSeverity(actual, minimum Severity) bool {
	severityOrder := map[Severity]int{
		SeverityCritical: 5,
		SeverityHigh:     4,
		SeverityMedium:   3,
		SeverityLow:      2,
		SeverityInfo:     1,
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

func containsDiscrepancyType(types []DiscrepancyType, target DiscrepancyType) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}

// calculateAlignmentSummary computes aggregate statistics.
func calculateAlignmentSummary(items []Discrepancy) AlignmentSummary {
	summary := AlignmentSummary{
		TotalDiscrepancies: len(items),
		ByType:             make(map[DiscrepancyType]int),
		BySeverity:         make(map[Severity]int),
		ByCategory:         make(map[Category]int),
		IsAligned:          true,
	}

	for _, item := range items {
		summary.ByType[item.Type]++
		summary.BySeverity[item.Severity]++
		summary.ByCategory[item.Category]++

		switch item.Severity {
		case SeverityCritical:
			summary.CriticalCount++
			summary.IsAligned = false
		case SeverityHigh:
			summary.HighCount++
			summary.IsAligned = false
		}
	}

	// Calculate alignment score (inverse of weighted discrepancies)
	if len(items) == 0 {
		summary.AlignmentScore = 1.0
	} else {
		// Weight: critical=10, high=5, medium=2, low=1, info=0.5
		weights := map[Severity]float64{
			SeverityCritical: 10,
			SeverityHigh:     5,
			SeverityMedium:   2,
			SeverityLow:      1,
			SeverityInfo:     0.5,
		}
		totalWeight := 0.0
		for _, item := range items {
			totalWeight += weights[item.Severity]
		}
		// Score decreases with more/worse discrepancies
		// Using a decay formula: score = 1 / (1 + totalWeight/10)
		summary.AlignmentScore = 1.0 / (1.0 + totalWeight/10.0)
	}

	return summary
}

// calculateCoverage computes coverage metrics.
func calculateCoverage(requirements []Requirement, _ []Implementation, discrepancies []Discrepancy) AlignmentCoverage {
	coverage := AlignmentCoverage{
		TotalRequirements: len(requirements),
	}

	// Count missing features
	missingCount := 0
	partialCount := 0
	for _, d := range discrepancies {
		switch d.Type {
		case DiscrepancyMissingFeature:
			missingCount++
		case DiscrepancyPartialImplementation:
			partialCount++
		case DiscrepancyUndocumentedCode:
			coverage.UndocumentedCount++
		}
	}

	coverage.MissingCount = missingCount
	coverage.PartialCount = partialCount
	coverage.ImplementedCount = coverage.TotalRequirements - missingCount - partialCount

	if coverage.TotalRequirements > 0 {
		coverage.CoveragePercentage = float64(coverage.ImplementedCount) / float64(coverage.TotalRequirements) * 100
	}

	return coverage
}

// HasDiscrepancies returns true if any discrepancies were found.
func (r *AlignmentResult) HasDiscrepancies() bool {
	return len(r.Discrepancies) > 0
}

// HasBlockers returns true if there are critical or high severity issues.
func (r *AlignmentResult) HasBlockers() bool {
	return !r.Summary.IsAligned
}

// FilterBySeverity returns items at or above the given severity.
func (r *AlignmentResult) FilterBySeverity(minSeverity Severity) []Discrepancy {
	return filterDiscrepancies(r.Discrepancies, AlignOptions{MinSeverity: minSeverity})
}

// FilterByCategory returns items in the given categories.
func (r *AlignmentResult) FilterByCategory(categories ...Category) []Discrepancy {
	return filterDiscrepancies(r.Discrepancies, AlignOptions{Categories: categories})
}

// FilterByType returns items of the given types.
func (r *AlignmentResult) FilterByType(types ...DiscrepancyType) []Discrepancy {
	var result []Discrepancy
	typeSet := make(map[DiscrepancyType]bool)
	for _, t := range types {
		typeSet[t] = true
	}
	for _, d := range r.Discrepancies {
		if typeSet[d.Type] {
			result = append(result, d)
		}
	}
	return result
}
