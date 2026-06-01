package testgen

import (
	"regexp"
	"strings"
)

// Parser extracts test cases from TPD markdown documents.
type Parser struct{}

// NewParser creates a new TPD parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse extracts test cases from TPD markdown content.
func (p *Parser) Parse(content string) (*ParsedTPD, error) {
	result := &ParsedTPD{
		FunctionalTests: []TestCase{},
		APITests:        []APITestCase{},
		JourneyTests:    []JourneyTestCase{},
	}

	// Parse functional test tables (Section 3)
	result.FunctionalTests = append(result.FunctionalTests, p.parseFunctionalTests(content)...)

	// Parse API test tables (Section 4.1)
	result.APITests = append(result.APITests, p.parseAPITests(content)...)

	// Parse journey test tables (Section 5.1)
	result.JourneyTests = append(result.JourneyTests, p.parseJourneyTests(content)...)

	return result, nil
}

// parseFunctionalTests extracts functional test cases from TPD.
// Looks for tables with format: | ID | Test Case | Input | Expected Output | Priority |
func (p *Parser) parseFunctionalTests(content string) []TestCase {
	var tests []TestCase

	// Find all tables with the functional test format
	tableRE := regexp.MustCompile(`(?s)\|\s*ID\s*\|\s*Test Case\s*\|\s*Input\s*\|\s*Expected Output\s*\|\s*Priority\s*\|[^\n]*\n\|[-\s|]+\|\n((?:\|[^\n]+\n)+)`)
	matches := tableRE.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		rows := strings.Split(strings.TrimSpace(match[1]), "\n")
		for _, row := range rows {
			tc := p.parseFunctionalRow(row)
			if tc.ID != "" {
				tests = append(tests, tc)
			}
		}
	}

	return tests
}

// parseFunctionalRow parses a single functional test table row.
func (p *Parser) parseFunctionalRow(row string) TestCase {
	cells := parseTableRow(row)
	if len(cells) < 5 {
		return TestCase{}
	}

	return TestCase{
		ID:       strings.TrimSpace(cells[0]),
		Title:    strings.TrimSpace(cells[1]),
		Type:     TestTypeFunctional,
		Input:    strings.TrimSpace(cells[2]),
		Expected: strings.TrimSpace(cells[3]),
		Priority: strings.TrimSpace(cells[4]),
	}
}

// parseAPITests extracts API test cases from TPD.
// Looks for tables with format: | Endpoint | Method | Test Scenario | Expected Response | Priority |
func (p *Parser) parseAPITests(content string) []APITestCase {
	var tests []APITestCase

	// Find all tables with the API test format
	tableRE := regexp.MustCompile(`(?s)\|\s*Endpoint\s*\|\s*Method\s*\|\s*Test Scenario\s*\|\s*Expected Response\s*\|\s*Priority\s*\|[^\n]*\n\|[-\s|]+\|\n((?:\|[^\n]+\n)+)`)
	matches := tableRE.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		rows := strings.Split(strings.TrimSpace(match[1]), "\n")
		for _, row := range rows {
			tc := p.parseAPIRow(row)
			if tc.Endpoint != "" || tc.Scenario != "" {
				tests = append(tests, tc)
			}
		}
	}

	return tests
}

// parseAPIRow parses a single API test table row.
func (p *Parser) parseAPIRow(row string) APITestCase {
	cells := parseTableRow(row)
	if len(cells) < 5 {
		return APITestCase{}
	}

	return APITestCase{
		Endpoint:         strings.TrimSpace(cells[0]),
		Method:           strings.TrimSpace(cells[1]),
		Scenario:         strings.TrimSpace(cells[2]),
		ExpectedResponse: strings.TrimSpace(cells[3]),
		Priority:         strings.TrimSpace(cells[4]),
	}
}

// parseJourneyTests extracts user journey test cases from TPD.
// Looks for tables with format: | Journey | Steps | Assertions | Priority |
func (p *Parser) parseJourneyTests(content string) []JourneyTestCase {
	var tests []JourneyTestCase

	// Find all tables with the journey test format
	tableRE := regexp.MustCompile(`(?s)\|\s*Journey\s*\|\s*Steps\s*\|\s*Assertions\s*\|\s*Priority\s*\|[^\n]*\n\|[-\s|]+\|\n((?:\|[^\n]+\n)+)`)
	matches := tableRE.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		rows := strings.Split(strings.TrimSpace(match[1]), "\n")
		for _, row := range rows {
			tc := p.parseJourneyRow(row)
			if tc.Journey != "" {
				tests = append(tests, tc)
			}
		}
	}

	return tests
}

// parseJourneyRow parses a single journey test table row.
func (p *Parser) parseJourneyRow(row string) JourneyTestCase {
	cells := parseTableRow(row)
	if len(cells) < 4 {
		return JourneyTestCase{}
	}

	// Parse steps and assertions (may be comma or semicolon separated)
	steps := parseList(cells[1])
	assertions := parseList(cells[2])

	return JourneyTestCase{
		Journey:    strings.TrimSpace(cells[0]),
		Steps:      steps,
		Assertions: assertions,
		Priority:   strings.TrimSpace(cells[3]),
	}
}

// parseTableRow splits a markdown table row into cells.
func parseTableRow(row string) []string {
	row = strings.TrimSpace(row)
	if !strings.HasPrefix(row, "|") {
		return nil
	}

	// Remove leading and trailing pipes
	row = strings.Trim(row, "|")

	// Split by pipe
	cells := strings.Split(row, "|")
	result := make([]string, len(cells))
	for i, cell := range cells {
		result[i] = strings.TrimSpace(cell)
	}
	return result
}

// parseList splits a string into a list by common separators.
func parseList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	// Try numbered list format: 1. item 2. item
	if regexp.MustCompile(`^\d+\.`).MatchString(s) {
		re := regexp.MustCompile(`\d+\.\s*`)
		parts := re.Split(s, -1)
		var result []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}

	// Try semicolon separation
	if strings.Contains(s, ";") {
		return splitAndTrim(s, ";")
	}

	// Try comma separation
	if strings.Contains(s, ",") {
		return splitAndTrim(s, ",")
	}

	// Return as single item
	return []string{s}
}

// splitAndTrim splits a string and trims whitespace from each part.
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
