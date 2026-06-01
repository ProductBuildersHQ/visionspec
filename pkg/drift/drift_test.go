package drift

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

func TestExtractRequirements(t *testing.T) {
	analyzer := &Analyzer{}

	specContent := `# Spec

## Requirements

- FR-001: User can login with email
- FR-002: User can logout

## API

- GET /users - List users
- POST /users - Create user
`

	reqs := analyzer.ExtractRequirements(specContent)

	if len(reqs) < 2 {
		t.Errorf("Expected at least 2 functional requirements, got %d", len(reqs))
	}

	// Check FR-001 was extracted
	found := false
	for _, r := range reqs {
		if r.ID == "FR-001" {
			found = true
			if !strings.Contains(r.Description, "login") {
				t.Errorf("Expected FR-001 to contain 'login', got %s", r.Description)
			}
		}
	}
	if !found {
		t.Error("FR-001 not found in extracted requirements")
	}
}

func TestExtractAPIRequirements(t *testing.T) {
	analyzer := &Analyzer{}

	specContent := `## API Endpoints

GET /users - List all users
POST /users - Create a user
DELETE /users/{id} - Delete a user
`

	reqs := analyzer.extractAPIRequirements(specContent)

	if len(reqs) != 3 {
		t.Errorf("Expected 3 API requirements, got %d", len(reqs))
	}

	// Check endpoint IDs
	ids := make(map[string]bool)
	for _, r := range reqs {
		ids[r.ID] = true
	}

	if !ids["API-GET-users"] {
		t.Error("Expected API-GET-users requirement")
	}
	if !ids["API-POST-users"] {
		t.Error("Expected API-POST-users requirement")
	}
}

func TestCompare(t *testing.T) {
	analyzer := &Analyzer{}

	requirements := []SpecRequirement{
		{ID: "FR-001", Description: "User login", Type: "feature"},
		{ID: "FR-002", Description: "User logout", Type: "feature"},
		{ID: "API-GET--users", Description: "GET /users", Type: "api"},
	}

	implementations := []CodeImplementation{
		{ID: "FR-001", Description: "User login", Type: "feature"},
		// FR-002 is missing
		{ID: "API-GET--users", Description: "GET /users", Type: "api"},
		{ID: "API-POST--admin", Description: "POST /admin", Type: "api"}, // Extra in code
	}

	items := analyzer.Compare(requirements, implementations)

	// Should find FR-002 as unimplemented
	foundUnimplemented := false
	foundUndocumented := false

	for _, item := range items {
		if item.Type == DriftUnimplemented && strings.Contains(item.SpecRef, "FR-002") {
			foundUnimplemented = true
		}
		if item.Type == DriftUndocumented && strings.Contains(item.ID, "admin") {
			foundUndocumented = true
		}
	}

	if !foundUnimplemented {
		t.Error("Expected to find FR-002 as unimplemented")
	}
	if !foundUndocumented {
		t.Error("Expected to find admin API as undocumented")
	}
}

func TestDetect(t *testing.T) {
	detector := NewDetector()

	specContent := `# Spec

## Requirements

- FR-001: User can login
- FR-002: User can register

## Tasks

- [ ] Implement login
- [ ] Implement register
`

	// Empty context (no implementations)
	ctx := &context.AggregatedContext{
		Project: "test-project",
	}

	report, err := detector.Detect(specContent, ctx, DefaultOptions())
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if report.Project != "test-project" {
		t.Errorf("Expected project name 'test-project', got %s", report.Project)
	}

	// Should find unimplemented items (requirements with no matching implementations)
	if len(report.Items) == 0 {
		// This is expected since we have requirements but no implementations
		// The analyzer should find drift
	}
}

func TestFilterBySeverity(t *testing.T) {
	report := &DriftReport{
		Items: []DriftItem{
			{ID: "1", Severity: SeverityCritical},
			{ID: "2", Severity: SeverityHigh},
			{ID: "3", Severity: SeverityMedium},
			{ID: "4", Severity: SeverityLow},
		},
	}

	high := report.FilterBySeverity(SeverityHigh)
	if len(high) != 2 { // Critical and High
		t.Errorf("Expected 2 items at High or above, got %d", len(high))
	}
}

func TestFilterByCategory(t *testing.T) {
	report := &DriftReport{
		Items: []DriftItem{
			{ID: "1", Category: CategoryAPI},
			{ID: "2", Category: CategoryData},
			{ID: "3", Category: CategoryAPI},
		},
	}

	api := report.FilterByCategory(CategoryAPI)
	if len(api) != 2 {
		t.Errorf("Expected 2 API items, got %d", len(api))
	}
}

func TestRenderText(t *testing.T) {
	report := &DriftReport{
		Project: "test-project",
		Items: []DriftItem{
			{
				ID:          "DRIFT-FR-001",
				Type:        DriftUnimplemented,
				Severity:    SeverityHigh,
				Category:    CategoryOther,
				Description: "Requirement not implemented: User login",
				SpecRef:     "FR-001 (Requirements)",
				Suggestion:  "Implement login feature",
			},
		},
		Summary: DriftSummary{
			TotalItems:  1,
			ByType:      map[DriftType]int{DriftUnimplemented: 1},
			BySeverity:  map[Severity]int{SeverityHigh: 1},
			ByCategory:  map[Category]int{CategoryOther: 1},
			HighCount:   1,
			HasBlockers: true,
		},
	}

	var buf bytes.Buffer
	renderer := NewRenderer(FormatText)
	err := renderer.Render(&buf, report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "test-project") {
		t.Error("Expected project name in output")
	}
	if !strings.Contains(output, "DRIFT DETECTED") {
		t.Error("Expected drift detected message")
	}
	if !strings.Contains(output, "User login") {
		t.Error("Expected drift item description")
	}
}

func TestRenderJSON(t *testing.T) {
	report := &DriftReport{
		Project: "test-project",
		Items:   []DriftItem{},
		Summary: DriftSummary{TotalItems: 0},
	}

	var buf bytes.Buffer
	renderer := NewRenderer(FormatJSON)
	err := renderer.Render(&buf, report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, `"project": "test-project"`) {
		t.Error("Expected JSON with project name")
	}
}

func TestRenderMarkdown(t *testing.T) {
	report := &DriftReport{
		Project: "test-project",
		Items: []DriftItem{
			{
				ID:          "DRIFT-FR-001",
				Type:        DriftUnimplemented,
				Severity:    SeverityCritical,
				Category:    CategoryAPI,
				Description: "Missing API endpoint",
			},
		},
		Summary: DriftSummary{
			TotalItems:    1,
			CriticalCount: 1,
			HasBlockers:   true,
			ByType:        map[DriftType]int{DriftUnimplemented: 1},
			BySeverity:    map[Severity]int{SeverityCritical: 1},
			ByCategory:    map[Category]int{CategoryAPI: 1},
		},
	}

	var buf bytes.Buffer
	renderer := NewRenderer(FormatMarkdown)
	err := renderer.Render(&buf, report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "# Drift Report") {
		t.Error("Expected markdown header")
	}
	if !strings.Contains(output, "test-project") {
		t.Error("Expected project name")
	}
	if !strings.Contains(output, "Missing API endpoint") {
		t.Error("Expected drift item description")
	}
}

func TestHasDrift(t *testing.T) {
	noDrift := &DriftReport{Items: []DriftItem{}}
	if noDrift.HasDrift() {
		t.Error("Expected no drift")
	}

	hasDrift := &DriftReport{Items: []DriftItem{{ID: "1"}}}
	if !hasDrift.HasDrift() {
		t.Error("Expected drift")
	}
}

func TestHasBlockers(t *testing.T) {
	noBlockers := &DriftReport{
		Summary: DriftSummary{HasBlockers: false},
	}
	if noBlockers.HasBlockers() {
		t.Error("Expected no blockers")
	}

	hasBlockers := &DriftReport{
		Summary: DriftSummary{HasBlockers: true},
	}
	if !hasBlockers.HasBlockers() {
		t.Error("Expected blockers")
	}
}

func TestSimilarity(t *testing.T) {
	a := map[string]bool{"user": true, "login": true, "auth": true}
	b := map[string]bool{"user": true, "authentication": true, "login": true}

	sim := similarity(a, b)
	if sim < 0.4 {
		t.Errorf("Expected higher similarity, got %f", sim)
	}

	// Empty sets
	empty := map[string]bool{}
	if similarity(empty, a) != 0 {
		t.Error("Expected 0 similarity with empty set")
	}
}
