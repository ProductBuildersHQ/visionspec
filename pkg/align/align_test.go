package align

import (
	"testing"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

func TestNewAligner(t *testing.T) {
	aligner := NewAligner()
	if aligner == nil {
		t.Error("NewAligner returned nil")
	}
	if aligner.comparator == nil {
		t.Error("Aligner.comparator is nil")
	}
}

func TestAlign_EmptySpec(t *testing.T) {
	aligner := NewAligner()
	ctx := &context.AggregatedContext{
		Project:    "test-project",
		GatheredAt: time.Now(),
	}

	result, err := aligner.Align("", ctx, DefaultAlignOptions())
	if err != nil {
		t.Fatalf("Align returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Align returned nil result")
	}

	if result.Project != "test-project" {
		t.Errorf("Expected project test-project, got %s", result.Project)
	}

	// Empty spec should have no discrepancies
	if len(result.Discrepancies) != 0 {
		t.Errorf("Expected 0 discrepancies for empty spec, got %d", len(result.Discrepancies))
	}
}

func TestAlign_WithRequirements(t *testing.T) {
	aligner := NewAligner()

	specContent := `# Product Requirements

## Features

- The system MUST support user authentication
- The system SHOULD display error messages
- Users can create new accounts

## API

- GET /api/users - List all users
- POST /api/users - Create a new user
`

	ctx := &context.AggregatedContext{
		Project:    "test-project",
		GatheredAt: time.Now(),
	}

	result, err := aligner.Align(specContent, ctx, DefaultAlignOptions())
	if err != nil {
		t.Fatalf("Align returned error: %v", err)
	}

	// Should have discrepancies since no implementations provided
	if len(result.Discrepancies) == 0 {
		t.Error("Expected discrepancies for unimplemented requirements")
	}

	// Check that requirements were extracted
	if result.Coverage.TotalRequirements == 0 {
		t.Error("Expected non-zero total requirements")
	}
}

func TestFilterDiscrepancies(t *testing.T) {
	items := []Discrepancy{
		{ID: "1", Severity: SeverityCritical, Category: CategoryAPI},
		{ID: "2", Severity: SeverityHigh, Category: CategoryData},
		{ID: "3", Severity: SeverityMedium, Category: CategoryAPI},
		{ID: "4", Severity: SeverityLow, Category: CategoryUI},
	}

	tests := []struct {
		name     string
		opts     AlignOptions
		expected int
	}{
		{
			name:     "no filter",
			opts:     AlignOptions{MinSeverity: SeverityLow},
			expected: 4,
		},
		{
			name:     "high severity filter",
			opts:     AlignOptions{MinSeverity: SeverityHigh},
			expected: 2, // critical and high
		},
		{
			name:     "critical only",
			opts:     AlignOptions{MinSeverity: SeverityCritical},
			expected: 1,
		},
		{
			name:     "category filter",
			opts:     AlignOptions{MinSeverity: SeverityLow, Categories: []Category{CategoryAPI}},
			expected: 2, // items 1 and 3
		},
		{
			name:     "max limit",
			opts:     AlignOptions{MinSeverity: SeverityLow, MaxDiscrepancies: 2},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterDiscrepancies(items, tt.opts)
			if len(result) != tt.expected {
				t.Errorf("Expected %d items, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestMeetsMinSeverity(t *testing.T) {
	tests := []struct {
		actual   Severity
		minimum  Severity
		expected bool
	}{
		{SeverityCritical, SeverityLow, true},
		{SeverityCritical, SeverityCritical, true},
		{SeverityHigh, SeverityCritical, false},
		{SeverityLow, SeverityMedium, false},
		{SeverityMedium, SeverityMedium, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.actual)+">="+string(tt.minimum), func(t *testing.T) {
			result := meetsMinSeverity(tt.actual, tt.minimum)
			if result != tt.expected {
				t.Errorf("meetsMinSeverity(%s, %s) = %v, want %v", tt.actual, tt.minimum, result, tt.expected)
			}
		})
	}
}

func TestAlignmentResult_Methods(t *testing.T) {
	result := &AlignmentResult{
		Discrepancies: []Discrepancy{
			{ID: "1", Severity: SeverityCritical, Category: CategoryAPI},
			{ID: "2", Severity: SeverityLow, Category: CategoryData},
		},
		Summary: AlignmentSummary{
			TotalDiscrepancies: 2,
			CriticalCount:      1,
			IsAligned:          false,
		},
	}

	if !result.HasDiscrepancies() {
		t.Error("HasDiscrepancies should return true")
	}

	if !result.HasBlockers() {
		t.Error("HasBlockers should return true when IsAligned is false")
	}

	filtered := result.FilterBySeverity(SeverityCritical)
	if len(filtered) != 1 {
		t.Errorf("FilterBySeverity(Critical) returned %d items, want 1", len(filtered))
	}

	filtered = result.FilterByCategory(CategoryAPI)
	if len(filtered) != 1 {
		t.Errorf("FilterByCategory(API) returned %d items, want 1", len(filtered))
	}
}

func TestCalculateAlignmentSummary(t *testing.T) {
	items := []Discrepancy{
		{Severity: SeverityCritical, Type: DiscrepancyMissingFeature},
		{Severity: SeverityHigh, Type: DiscrepancyDiverged},
		{Severity: SeverityMedium, Type: DiscrepancyPartialImplementation},
	}

	summary := calculateAlignmentSummary(items)

	if summary.TotalDiscrepancies != 3 {
		t.Errorf("TotalDiscrepancies = %d, want 3", summary.TotalDiscrepancies)
	}

	if summary.CriticalCount != 1 {
		t.Errorf("CriticalCount = %d, want 1", summary.CriticalCount)
	}

	if summary.HighCount != 1 {
		t.Errorf("HighCount = %d, want 1", summary.HighCount)
	}

	if summary.IsAligned {
		t.Error("IsAligned should be false with critical/high issues")
	}

	if summary.AlignmentScore >= 1.0 {
		t.Error("AlignmentScore should be less than 1.0 with discrepancies")
	}
}
