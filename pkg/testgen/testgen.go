// Package testgen generates executable test stubs from TPD (Test Plan Document).
package testgen

import (
	"fmt"
	"time"
)

// TestType represents the type of test case.
type TestType string

const (
	TestTypeFunctional  TestType = "functional"
	TestTypeAPI         TestType = "api"
	TestTypeIntegration TestType = "integration"
	TestTypeE2E         TestType = "e2e"
	TestTypePerformance TestType = "performance"
	TestTypeUAT         TestType = "uat"
)

// TestCase represents a single test case extracted from TPD.
type TestCase struct {
	ID        string   `json:"id"`         // TC-001, API-001
	Title     string   `json:"title"`      // Test case title
	Type      TestType `json:"type"`       // functional, api, integration, e2e, performance
	Input     string   `json:"input"`      // Test input description
	Expected  string   `json:"expected"`   // Expected output/result
	Priority  string   `json:"priority"`   // P0, P1, P2
	Steps     []string `json:"steps"`      // Test steps
	SourceRef string   `json:"source_ref"` // REQ-XXX traceability
}

// APITestCase represents an API-specific test case.
type APITestCase struct {
	Endpoint         string `json:"endpoint"`
	Method           string `json:"method"`
	Scenario         string `json:"scenario"`
	ExpectedResponse string `json:"expected_response"`
	Priority         string `json:"priority"`
}

// JourneyTestCase represents a user journey test case.
type JourneyTestCase struct {
	Journey    string   `json:"journey"`
	Steps      []string `json:"steps"`
	Assertions []string `json:"assertions"`
	Priority   string   `json:"priority"`
}

// Generator defines the interface for test stub generators.
type Generator interface {
	// Generate generates test stubs from parsed test cases.
	Generate(cases []TestCase, opts GenerateOptions) (*GenerateResult, error)

	// Language returns the target language name.
	Language() string
}

// GenerateOptions configures test generation.
type GenerateOptions struct {
	OutputDir     string // Output directory for generated files
	PackageName   string // Package/module name
	TestFramework string // testing, testify, jest, pytest
	GroupBy       string // type, file, priority
}

// DefaultOptions returns sensible defaults for generation.
func DefaultOptions() GenerateOptions {
	return GenerateOptions{
		OutputDir:   ".",
		PackageName: "tests",
		GroupBy:     "type",
	}
}

// GenerateResult contains the output of test generation.
type GenerateResult struct {
	Language    string          `json:"language"`
	Framework   string          `json:"framework"`
	OutputDir   string          `json:"output_dir"`
	Files       []GeneratedFile `json:"files"`
	TotalTests  int             `json:"total_tests"`
	GeneratedAt time.Time       `json:"generated_at"`
}

// GeneratedFile represents a generated test file.
type GeneratedFile struct {
	Path      string `json:"path"`
	TestCount int    `json:"test_count"`
	Content   string `json:"-"` // Content not included in JSON
}

// ParsedTPD contains all test cases parsed from a TPD document.
type ParsedTPD struct {
	FunctionalTests []TestCase        `json:"functional_tests"`
	APITests        []APITestCase     `json:"api_tests"`
	JourneyTests    []JourneyTestCase `json:"journey_tests"`
}

// AllTestCases returns all test cases as a flat list.
func (p *ParsedTPD) AllTestCases() []TestCase {
	var all []TestCase
	all = append(all, p.FunctionalTests...)

	// Convert API tests to TestCase
	for i, api := range p.APITests {
		all = append(all, TestCase{
			ID:       fmt.Sprintf("API-%03d", i+1),
			Title:    fmt.Sprintf("%s %s: %s", api.Method, api.Endpoint, api.Scenario),
			Type:     TestTypeAPI,
			Input:    fmt.Sprintf("%s %s", api.Method, api.Endpoint),
			Expected: api.ExpectedResponse,
			Priority: api.Priority,
		})
	}

	// Convert journey tests to TestCase
	for i, journey := range p.JourneyTests {
		all = append(all, TestCase{
			ID:       fmt.Sprintf("E2E-%03d", i+1),
			Title:    journey.Journey,
			Type:     TestTypeE2E,
			Steps:    journey.Steps,
			Expected: joinStrings(journey.Assertions, "; "),
			Priority: journey.Priority,
		})
	}

	return all
}

// TestsByType groups test cases by type.
func (p *ParsedTPD) TestsByType() map[TestType][]TestCase {
	result := make(map[TestType][]TestCase)
	for _, tc := range p.AllTestCases() {
		result[tc.Type] = append(result[tc.Type], tc)
	}
	return result
}

// TestsByPriority groups test cases by priority.
func (p *ParsedTPD) TestsByPriority() map[string][]TestCase {
	result := make(map[string][]TestCase)
	for _, tc := range p.AllTestCases() {
		priority := tc.Priority
		if priority == "" {
			priority = "P2"
		}
		result[priority] = append(result[priority], tc)
	}
	return result
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for i := 1; i < len(ss); i++ {
		result += sep + ss[i]
	}
	return result
}

// Registry holds registered generators.
var registry = make(map[string]Generator)

// Register adds a generator to the registry.
func Register(gen Generator) {
	registry[gen.Language()] = gen
}

// Get retrieves a generator by language name.
func Get(language string) (Generator, error) {
	gen, ok := registry[language]
	if !ok {
		return nil, fmt.Errorf("unknown generator language: %s", language)
	}
	return gen, nil
}

// Available returns all registered generator language names.
func Available() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
