package testmap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMapper(t *testing.T) {
	mapper := NewMapper("/test/project")

	if mapper.projectPath != "/test/project" {
		t.Errorf("projectPath = %s, want /test/project", mapper.projectPath)
	}

	if len(mapper.specPatterns) == 0 {
		t.Error("Should have default spec patterns")
	}

	if len(mapper.testPatterns) == 0 {
		t.Error("Should have default test patterns")
	}
}

func TestMapper_WithPatterns(t *testing.T) {
	mapper := NewMapper("/test").
		WithSpecPatterns([]string{"*.spec"}).
		WithTestPatterns([]string{"*_test.py"})

	if len(mapper.specPatterns) != 1 || mapper.specPatterns[0] != "*.spec" {
		t.Error("Should set custom spec patterns")
	}

	if len(mapper.testPatterns) != 1 || mapper.testPatterns[0] != "*_test.py" {
		t.Error("Should set custom test patterns")
	}
}

func TestMapper_Map(t *testing.T) {
	tmpDir := t.TempDir()

	// Create spec directory with requirements
	specDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Requirements

## FR-001 User Login

Users should be able to log in.

## FR-002 User Logout

Users should be able to log out.

## REQ-100 Performance

System should be fast.
`
	if err := os.WriteFile(filepath.Join(specDir, "prd.md"), []byte(specContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create test file with requirement references
	testContent := `package auth

// TestUserLogin covers FR-001
func TestUserLogin(t *testing.T) {
	// test implementation
}

// Tests: FR-001
func TestLoginValidation(t *testing.T) {
}

func TestOther(t *testing.T) {
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "auth_test.go"), []byte(testContent), 0600); err != nil {
		t.Fatal(err)
	}

	mapper := NewMapper(tmpDir)
	mapping, err := mapper.Map()
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if mapping.Project == "" {
		t.Error("Project should be set")
	}

	if len(mapping.Requirements) == 0 {
		t.Error("Should find requirements")
	}

	if len(mapping.Tests) == 0 {
		t.Error("Should find tests")
	}
}

func TestMapper_ExtractRequirementTitle(t *testing.T) {
	mapper := NewMapper("/test")

	tests := []struct {
		content  string
		reqID    string
		expected string
	}{
		{
			content:  "## FR-001: User Login Feature\nDescription",
			reqID:    "FR-001",
			expected: "User Login Feature",
		},
		{
			content:  "### REQ-42 - Performance Optimization\nDetails",
			reqID:    "REQ-42",
			expected: "- Performance Optimization",
		},
		{
			content:  "No title here",
			reqID:    "FR-999",
			expected: "",
		},
	}

	for _, tt := range tests {
		result := mapper.extractRequirementTitle(tt.content, tt.reqID)
		if result != tt.expected {
			t.Errorf("extractRequirementTitle(%q, %q) = %q, want %q",
				tt.content, tt.reqID, result, tt.expected)
		}
	}
}

func TestMapper_ExtractTestFunctions(t *testing.T) {
	mapper := NewMapper("/test")

	goContent := `package test

// Covers: FR-001
func TestLogin(t *testing.T) {}

func TestLogout(t *testing.T) {}

func Helper() {}
`

	tests := mapper.extractTestFunctions(goContent, "auth_test.go")

	if len(tests) != 2 {
		t.Errorf("Should find 2 tests, found %d", len(tests))
	}

	// First test should have requirement reference
	hasReq := false
	for _, test := range tests {
		if len(test.Requirements) > 0 {
			hasReq = true
			break
		}
	}
	if !hasReq {
		t.Error("Should find at least one test with requirement reference")
	}
}

func TestMapper_CalculateSummary(t *testing.T) {
	mapper := NewMapper("/test")

	requirements := []RequirementCover{
		{ID: "FR-001", Status: "covered", Coverage: 100},
		{ID: "FR-002", Status: "partial", Coverage: 50},
		{ID: "FR-003", Status: "uncovered", Coverage: 0},
	}

	tests := []TestInfo{
		{Name: "TestA", Requirements: []string{"FR-001"}},
		{Name: "TestB", Requirements: []string{}},
		{Name: "TestC", Requirements: []string{"FR-002"}},
	}

	summary := mapper.calculateSummary(requirements, tests)

	if summary.TotalRequirements != 3 {
		t.Errorf("TotalRequirements = %d, want 3", summary.TotalRequirements)
	}

	if summary.CoveredRequirements != 1 {
		t.Errorf("CoveredRequirements = %d, want 1", summary.CoveredRequirements)
	}

	if summary.PartialRequirements != 1 {
		t.Errorf("PartialRequirements = %d, want 1", summary.PartialRequirements)
	}

	if summary.UncoveredRequirements != 1 {
		t.Errorf("UncoveredRequirements = %d, want 1", summary.UncoveredRequirements)
	}

	if summary.TotalTests != 3 {
		t.Errorf("TotalTests = %d, want 3", summary.TotalTests)
	}

	if summary.MappedTests != 2 {
		t.Errorf("MappedTests = %d, want 2", summary.MappedTests)
	}
}

func TestMapper_FindUnmapped(t *testing.T) {
	mapper := NewMapper("/test")

	requirements := []RequirementCover{
		{ID: "FR-001", Status: "covered"},
		{ID: "FR-002", Status: "uncovered"},
	}

	tests := []TestInfo{
		{Name: "TestA", Requirements: []string{"FR-001"}},
		{Name: "TestB", Requirements: []string{}},
	}

	unmapped := mapper.findUnmapped(requirements, tests)

	if len(unmapped.Requirements) != 1 || unmapped.Requirements[0] != "FR-002" {
		t.Error("Should find FR-002 as unmapped requirement")
	}

	if len(unmapped.Tests) != 1 || unmapped.Tests[0] != "TestB" {
		t.Error("Should find TestB as unmapped test")
	}
}

func TestCoverageMapping_RenderMarkdown(t *testing.T) {
	mapping := &CoverageMapping{
		Project: "test-project",
		Summary: CoverageSummary{
			TotalRequirements:     5,
			CoveredRequirements:   2,
			PartialRequirements:   1,
			UncoveredRequirements: 2,
			TotalTests:            10,
			MappedTests:           6,
			OverallCoverage:       50.0,
		},
		Requirements: []RequirementCover{
			{ID: "FR-001", Title: "Login", Status: "covered", Coverage: 100},
			{ID: "FR-002", Title: "Logout", Status: "uncovered", Coverage: 0},
		},
		Unmapped: UnmappedItems{
			Requirements: []string{"FR-002"},
			Tests:        []string{"TestOrphan"},
		},
	}

	md := mapping.RenderMarkdown()

	if md == "" {
		t.Error("Should generate Markdown")
	}

	if !strings.Contains(md, "# Test Coverage Mapping") {
		t.Error("Should have title")
	}

	if !strings.Contains(md, "test-project") {
		t.Error("Should contain project name")
	}

	if !strings.Contains(md, "FR-001") {
		t.Error("Should contain requirement ID")
	}

	if !strings.Contains(md, "Uncovered Requirements") {
		t.Error("Should have uncovered section")
	}
}

func TestCoverageMapping_ExportJSON(t *testing.T) {
	mapping := &CoverageMapping{
		Project: "test",
		Summary: CoverageSummary{TotalRequirements: 1},
	}

	data, err := mapping.ExportJSON()
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Should produce JSON output")
	}

	if !strings.Contains(string(data), "test") {
		t.Error("JSON should contain project name")
	}
}

func TestAppendUnique(t *testing.T) {
	slice := []string{"a", "b"}

	// Add new item
	result := appendUnique(slice, "c")
	if len(result) != 3 {
		t.Error("Should add new item")
	}

	// Add existing item
	result = appendUnique(result, "a")
	if len(result) != 3 {
		t.Error("Should not add duplicate")
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		status   string
		contains string
	}{
		{"covered", "Covered"},
		{"partial", "Partial"},
		{"uncovered", "Uncovered"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := statusIcon(tt.status)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("statusIcon(%q) = %q, should contain %q", tt.status, result, tt.contains)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exact", 5, "exact"},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}
