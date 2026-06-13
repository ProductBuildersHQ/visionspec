// Package testmap provides test coverage mapping between specs and code.
//
// This package analyzes test coverage data and maps it to specification
// requirements, helping identify which requirements have adequate test
// coverage and which need additional testing.
package testmap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// CoverageMapping represents the relationship between specs and tests.
type CoverageMapping struct {
	Project      string             `json:"project"`
	GeneratedAt  time.Time          `json:"generated_at"`
	Requirements []RequirementCover `json:"requirements"`
	Tests        []TestInfo         `json:"tests"`
	Summary      CoverageSummary    `json:"summary"`
	Unmapped     UnmappedItems      `json:"unmapped"`
}

// RequirementCover shows test coverage for a requirement.
type RequirementCover struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	SpecFile    string   `json:"spec_file"`
	SpecSection string   `json:"spec_section,omitempty"`
	Tests       []string `json:"tests,omitempty"`      // Test function names
	TestFiles   []string `json:"test_files,omitempty"` // Test file paths
	Coverage    float64  `json:"coverage"`             // 0-100
	Status      string   `json:"status"`               // covered, partial, uncovered
	Priority    string   `json:"priority,omitempty"`
}

// TestInfo contains information about a test.
type TestInfo struct {
	Name         string   `json:"name"`
	File         string   `json:"file"`
	Package      string   `json:"package,omitempty"`
	Requirements []string `json:"requirements,omitempty"` // Requirement IDs this test covers
	Passed       bool     `json:"passed"`
	Duration     float64  `json:"duration_ms,omitempty"`
}

// CoverageSummary provides aggregate coverage statistics.
type CoverageSummary struct {
	TotalRequirements     int     `json:"total_requirements"`
	CoveredRequirements   int     `json:"covered_requirements"`
	PartialRequirements   int     `json:"partial_requirements"`
	UncoveredRequirements int     `json:"uncovered_requirements"`
	TotalTests            int     `json:"total_tests"`
	MappedTests           int     `json:"mapped_tests"`
	OverallCoverage       float64 `json:"overall_coverage"` // 0-100
}

// UnmappedItems tracks items without proper mapping.
type UnmappedItems struct {
	Requirements []string `json:"requirements,omitempty"` // Requirements without tests
	Tests        []string `json:"tests,omitempty"`        // Tests without requirement refs
}

// Mapper performs test coverage mapping.
type Mapper struct {
	projectPath  string
	specPatterns []string
	testPatterns []string
}

// NewMapper creates a new test coverage mapper.
func NewMapper(projectPath string) *Mapper {
	return &Mapper{
		projectPath:  projectPath,
		specPatterns: []string{"*.md", "*.yaml", "*.yml"},
		testPatterns: []string{"*_test.go", "test_*.py", "*_test.ts", "*.test.js"},
	}
}

// WithSpecPatterns sets custom spec file patterns.
func (m *Mapper) WithSpecPatterns(patterns []string) *Mapper {
	m.specPatterns = patterns
	return m
}

// WithTestPatterns sets custom test file patterns.
func (m *Mapper) WithTestPatterns(patterns []string) *Mapper {
	m.testPatterns = patterns
	return m
}

// Map generates a coverage mapping from specs and tests.
func (m *Mapper) Map() (*CoverageMapping, error) {
	mapping := &CoverageMapping{
		Project:      filepath.Base(m.projectPath),
		GeneratedAt:  time.Now(),
		Requirements: []RequirementCover{},
		Tests:        []TestInfo{},
	}

	// Extract requirements from spec files
	requirements := m.extractRequirements()

	// Find test files and extract test info
	tests, err := m.extractTests()
	if err != nil {
		return nil, fmt.Errorf("extracting tests: %w", err)
	}

	// Map tests to requirements
	m.mapTestsToRequirements(requirements, tests)

	mapping.Requirements = requirements
	mapping.Tests = tests
	mapping.Summary = m.calculateSummary(requirements, tests)
	mapping.Unmapped = m.findUnmapped(requirements, tests)

	return mapping, nil
}

// extractRequirements finds requirements in spec files.
func (m *Mapper) extractRequirements() []RequirementCover {
	var requirements []RequirementCover

	// Look in common spec locations
	specDirs := []string{
		filepath.Join(m.projectPath, "source"),
		filepath.Join(m.projectPath, "spec"),
		filepath.Join(m.projectPath, "specs"),
		m.projectPath,
	}

	// Regex patterns for requirement IDs
	reqPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)REQ-\d+`),
		regexp.MustCompile(`(?i)FR-\d+`),  // Functional requirement
		regexp.MustCompile(`(?i)NFR-\d+`), // Non-functional requirement
		regexp.MustCompile(`(?i)US-\d+`),  // User story
		regexp.MustCompile(`(?i)STORY-\d+`),
		regexp.MustCompile(`(?i)FEAT-\d+`),        // Feature
		regexp.MustCompile(`\[REQ:\s*([^\]]+)\]`), // [REQ: description]
	}

	seen := make(map[string]bool)

	for _, dir := range specDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Skip directories that don't exist
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			for _, pattern := range m.specPatterns {
				matched, _ := filepath.Match(pattern, entry.Name())
				if !matched {
					continue
				}

				filePath := filepath.Join(dir, entry.Name())
				content, err := os.ReadFile(filePath)
				if err != nil {
					continue
				}

				// Find requirements in the content
				for _, re := range reqPatterns {
					matches := re.FindAllString(string(content), -1)
					for _, match := range matches {
						reqID := strings.ToUpper(strings.TrimSpace(match))
						if seen[reqID] {
							continue
						}
						seen[reqID] = true

						req := RequirementCover{
							ID:       reqID,
							SpecFile: filePath,
							Status:   "uncovered",
							Coverage: 0,
						}

						// Try to extract title from context
						req.Title = m.extractRequirementTitle(string(content), match)

						requirements = append(requirements, req)
					}
				}
				break // Only match once per file
			}
		}
	}

	// Sort by ID for consistent output
	sort.Slice(requirements, func(i, j int) bool {
		return requirements[i].ID < requirements[j].ID
	})

	return requirements
}

// extractRequirementTitle attempts to extract the title for a requirement.
func (m *Mapper) extractRequirementTitle(content, reqID string) string {
	// Look for patterns like "REQ-001: Title" or "REQ-001 - Title"
	patterns := []string{
		reqID + `[:\-]\s*([^\n]+)`,
		`#[#]*\s*` + reqID + `[:\-]?\s*([^\n]+)`,
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		if match := re.FindStringSubmatch(content); len(match) > 1 {
			return strings.TrimSpace(match[1])
		}
	}

	return ""
}

// extractTests finds test functions and their requirement references.
func (m *Mapper) extractTests() ([]TestInfo, error) {
	var tests []TestInfo

	err := filepath.Walk(m.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			// Skip common non-test directories
			base := filepath.Base(path)
			if base == "vendor" || base == "node_modules" || base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this is a test file
		isTestFile := false
		for _, pattern := range m.testPatterns {
			if matched, _ := filepath.Match(pattern, info.Name()); matched {
				isTestFile = true
				break
			}
		}

		if !isTestFile {
			return nil
		}

		content, err := os.ReadFile(path) //nolint:gosec // G122: filepath.Walk callback operates on trusted project paths
		if err != nil {
			return nil
		}

		// Extract test functions based on file type
		fileTests := m.extractTestFunctions(string(content), path)
		tests = append(tests, fileTests...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tests, nil
}

// extractTestFunctions extracts test function info from file content.
func (m *Mapper) extractTestFunctions(content, filePath string) []TestInfo {
	var tests []TestInfo

	ext := filepath.Ext(filePath)

	// Different patterns for different languages
	var testPattern *regexp.Regexp
	switch ext {
	case ".go":
		testPattern = regexp.MustCompile(`func\s+(Test\w+)\s*\(`)
	case ".py":
		testPattern = regexp.MustCompile(`def\s+(test_\w+)\s*\(`)
	case ".ts", ".js":
		testPattern = regexp.MustCompile(`(?:it|test|describe)\s*\(\s*['"]([^'"]+)['"]`)
	default:
		return tests
	}

	// Find requirement references in comments
	reqRefPattern := regexp.MustCompile(`(?i)(?:covers?|tests?|verifies?|implements?)\s*:?\s*((?:REQ|FR|NFR|US|STORY|FEAT)-\d+)`)

	matches := testPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		test := TestInfo{
			Name:   match[1],
			File:   filePath,
			Passed: true, // Default to true; actual status would come from test results
		}

		// Look for requirement references near the test
		testIdx := strings.Index(content, match[0])
		if testIdx >= 0 {
			// Check preceding 500 chars for requirement refs (comments)
			start := testIdx - 500
			if start < 0 {
				start = 0
			}
			context := content[start:testIdx]

			reqRefs := reqRefPattern.FindAllStringSubmatch(context, -1)
			for _, ref := range reqRefs {
				if len(ref) > 1 {
					test.Requirements = append(test.Requirements, strings.ToUpper(ref[1]))
				}
			}
		}

		tests = append(tests, test)
	}

	return tests
}

// mapTestsToRequirements links tests to their requirements.
func (m *Mapper) mapTestsToRequirements(requirements []RequirementCover, tests []TestInfo) {
	reqIndex := make(map[string]*RequirementCover)
	for i := range requirements {
		reqIndex[requirements[i].ID] = &requirements[i]
	}

	for _, test := range tests {
		for _, reqID := range test.Requirements {
			if req, ok := reqIndex[reqID]; ok {
				req.Tests = append(req.Tests, test.Name)
				req.TestFiles = appendUnique(req.TestFiles, test.File)
			}
		}
	}

	// Update coverage status
	for i := range requirements {
		req := &requirements[i]
		if len(req.Tests) == 0 {
			req.Status = "uncovered"
			req.Coverage = 0
		} else if len(req.Tests) >= 3 {
			req.Status = "covered"
			req.Coverage = 100
		} else {
			req.Status = "partial"
			req.Coverage = float64(len(req.Tests)) * 33.3
		}
	}
}

// calculateSummary computes aggregate statistics.
func (m *Mapper) calculateSummary(requirements []RequirementCover, tests []TestInfo) CoverageSummary {
	summary := CoverageSummary{
		TotalRequirements: len(requirements),
		TotalTests:        len(tests),
	}

	for _, req := range requirements {
		switch req.Status {
		case "covered":
			summary.CoveredRequirements++
		case "partial":
			summary.PartialRequirements++
		case "uncovered":
			summary.UncoveredRequirements++
		}
	}

	for _, test := range tests {
		if len(test.Requirements) > 0 {
			summary.MappedTests++
		}
	}

	if summary.TotalRequirements > 0 {
		covered := float64(summary.CoveredRequirements) + float64(summary.PartialRequirements)*0.5
		summary.OverallCoverage = (covered / float64(summary.TotalRequirements)) * 100
	}

	return summary
}

// findUnmapped identifies items without proper mapping.
func (m *Mapper) findUnmapped(requirements []RequirementCover, tests []TestInfo) UnmappedItems {
	unmapped := UnmappedItems{}

	for _, req := range requirements {
		if req.Status == "uncovered" {
			unmapped.Requirements = append(unmapped.Requirements, req.ID)
		}
	}

	for _, test := range tests {
		if len(test.Requirements) == 0 {
			unmapped.Tests = append(unmapped.Tests, test.Name)
		}
	}

	return unmapped
}

// appendUnique appends a string to a slice if not already present.
func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// RenderMarkdown outputs the coverage mapping as Markdown.
func (cm *CoverageMapping) RenderMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Test Coverage Mapping\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** %s\n\n", cm.Project))
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", cm.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Total Requirements | %d |\n", cm.Summary.TotalRequirements))
	sb.WriteString(fmt.Sprintf("| Covered | %d |\n", cm.Summary.CoveredRequirements))
	sb.WriteString(fmt.Sprintf("| Partial | %d |\n", cm.Summary.PartialRequirements))
	sb.WriteString(fmt.Sprintf("| Uncovered | %d |\n", cm.Summary.UncoveredRequirements))
	sb.WriteString(fmt.Sprintf("| Overall Coverage | %.1f%% |\n", cm.Summary.OverallCoverage))
	sb.WriteString(fmt.Sprintf("| Total Tests | %d |\n", cm.Summary.TotalTests))
	sb.WriteString(fmt.Sprintf("| Mapped Tests | %d |\n", cm.Summary.MappedTests))
	sb.WriteString("\n")

	// Coverage by requirement
	sb.WriteString("## Requirements Coverage\n\n")
	sb.WriteString("| ID | Title | Status | Tests | Coverage |\n")
	sb.WriteString("|----|-------|--------|-------|----------|\n")

	for _, req := range cm.Requirements {
		status := statusIcon(req.Status)
		title := req.Title
		if title == "" {
			title = "-"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %.0f%% |\n",
			req.ID, truncate(title, 40), status, len(req.Tests), req.Coverage))
	}
	sb.WriteString("\n")

	// Uncovered requirements
	if len(cm.Unmapped.Requirements) > 0 {
		sb.WriteString("## Uncovered Requirements\n\n")
		sb.WriteString("The following requirements have no test coverage:\n\n")
		for _, reqID := range cm.Unmapped.Requirements {
			sb.WriteString(fmt.Sprintf("- ⚠️ %s\n", reqID))
		}
		sb.WriteString("\n")
	}

	// Unmapped tests
	if len(cm.Unmapped.Tests) > 0 && len(cm.Unmapped.Tests) <= 20 {
		sb.WriteString("## Unmapped Tests\n\n")
		sb.WriteString("The following tests are not linked to any requirement:\n\n")
		for _, test := range cm.Unmapped.Tests {
			sb.WriteString(fmt.Sprintf("- %s\n", test))
		}
		sb.WriteString("\n")
	} else if len(cm.Unmapped.Tests) > 20 {
		sb.WriteString(fmt.Sprintf("## Unmapped Tests\n\n%d tests are not linked to any requirement.\n\n",
			len(cm.Unmapped.Tests)))
	}

	return sb.String()
}

// ExportJSON exports the coverage mapping as JSON.
func (cm *CoverageMapping) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(cm, "", "  ")
}

// statusIcon returns an icon for a coverage status.
func statusIcon(status string) string {
	switch status {
	case "covered":
		return "✅ Covered"
	case "partial":
		return "🟡 Partial"
	case "uncovered":
		return "❌ Uncovered"
	default:
		return status
	}
}

// truncate shortens a string.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
