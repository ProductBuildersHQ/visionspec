// Package patterns provides specification pattern detection across projects.
package patterns

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Detector identifies common patterns in specification files.
type Detector struct {
	specsDir string
}

// NewDetector creates a new pattern detector.
func NewDetector(specsDir string) *Detector {
	return &Detector{specsDir: specsDir}
}

// PatternReport contains all detected patterns.
type PatternReport struct {
	StructuralPatterns []StructuralPattern `json:"structural_patterns"`
	ContentPatterns    []ContentPattern    `json:"content_patterns"`
	AntiPatterns       []AntiPattern       `json:"anti_patterns"`
	Summary            PatternSummary      `json:"summary"`
}

// StructuralPattern represents a common document structure pattern.
type StructuralPattern struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Projects    []string `json:"projects"`
	Occurrences int      `json:"occurrences"`
	Example     string   `json:"example,omitempty"`
}

// ContentPattern represents a common content pattern.
type ContentPattern struct {
	Type        string            `json:"type"`
	Pattern     string            `json:"pattern"`
	Description string            `json:"description"`
	Occurrences int               `json:"occurrences"`
	Examples    []PatternInstance `json:"examples"`
}

// PatternInstance is a specific occurrence of a pattern.
type PatternInstance struct {
	Project  string `json:"project"`
	SpecType string `json:"spec_type"`
	Line     int    `json:"line"`
	Text     string `json:"text"`
}

// AntiPattern represents a detected specification anti-pattern.
type AntiPattern struct {
	Type        string            `json:"type"`
	Severity    string            `json:"severity"` // high, medium, low
	Description string            `json:"description"`
	Instances   []PatternInstance `json:"instances"`
	Suggestion  string            `json:"suggestion"`
}

// PatternSummary provides aggregate statistics.
type PatternSummary struct {
	TotalSpecs       int            `json:"total_specs"`
	StructuralCount  int            `json:"structural_pattern_count"`
	ContentCount     int            `json:"content_pattern_count"`
	AntiPatternCount int            `json:"anti_pattern_count"`
	CommonSections   []string       `json:"common_sections"`
	PatternsByType   map[string]int `json:"patterns_by_type"`
	QualityScore     float64        `json:"quality_score"` // 0-100
}

// SpecInfo holds extracted information about a spec file.
type SpecInfo struct {
	Project  string
	SpecType string
	Path     string
	Sections []string
	Content  string
	Lines    []string
}

// Detect runs all pattern detection algorithms.
func (d *Detector) Detect() (*PatternReport, error) {
	// Load all specs
	specs, err := d.loadSpecs()
	if err != nil {
		return nil, err
	}

	report := &PatternReport{
		StructuralPatterns: []StructuralPattern{},
		ContentPatterns:    []ContentPattern{},
		AntiPatterns:       []AntiPattern{},
	}

	// Detect structural patterns
	report.StructuralPatterns = d.detectStructuralPatterns(specs)

	// Detect content patterns
	report.ContentPatterns = d.detectContentPatterns(specs)

	// Detect anti-patterns
	report.AntiPatterns = d.detectAntiPatterns(specs)

	// Calculate summary
	report.Summary = d.calculateSummary(specs, report)

	return report, nil
}

// loadSpecs loads all specification files.
func (d *Detector) loadSpecs() ([]SpecInfo, error) {
	var specs []SpecInfo

	err := filepath.Walk(d.specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		relPath, _ := filepath.Rel(d.specsDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 1 {
			return nil
		}

		content, err := os.ReadFile(path) //nolint:gosec // G122: User-provided specs directory for internal analysis
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		sections := extractSections(lines)

		specs = append(specs, SpecInfo{
			Project:  parts[0],
			SpecType: strings.TrimSuffix(info.Name(), ".md"),
			Path:     relPath,
			Sections: sections,
			Content:  string(content),
			Lines:    lines,
		})

		return nil
	})

	return specs, err
}

// extractSections extracts section headers from lines.
func extractSections(lines []string) []string {
	var sections []string
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			section := strings.TrimPrefix(line, "## ")
			sections = append(sections, strings.TrimSpace(section))
		}
	}
	return sections
}

// detectStructuralPatterns finds common document structures.
func (d *Detector) detectStructuralPatterns(specs []SpecInfo) []StructuralPattern {
	var patterns []StructuralPattern

	// Count section occurrences across specs
	sectionCounts := make(map[string][]string)          // section -> projects
	specTypeSections := make(map[string]map[string]int) // specType -> section -> count

	for _, spec := range specs {
		for _, section := range spec.Sections {
			sectionNorm := strings.ToLower(section)
			sectionCounts[sectionNorm] = append(sectionCounts[sectionNorm], spec.Project)

			if specTypeSections[spec.SpecType] == nil {
				specTypeSections[spec.SpecType] = make(map[string]int)
			}
			specTypeSections[spec.SpecType][sectionNorm]++
		}
	}

	// Find commonly used sections
	for section, projects := range sectionCounts {
		if len(projects) >= 3 {
			// Get unique projects
			uniqueProjects := unique(projects)
			patterns = append(patterns, StructuralPattern{
				Type:        "common_section",
				Name:        section,
				Description: "Section appears across multiple projects",
				Projects:    uniqueProjects,
				Occurrences: len(projects),
			})
		}
	}

	// Find spec-type specific patterns
	for specType, sections := range specTypeSections {
		for section, count := range sections {
			if count >= 2 {
				patterns = append(patterns, StructuralPattern{
					Type:        "spec_type_section",
					Name:        section,
					Description: "Standard section for " + specType + " specs",
					Occurrences: count,
				})
			}
		}
	}

	// Sort by occurrences
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Occurrences > patterns[j].Occurrences
	})

	// Limit
	if len(patterns) > 25 {
		patterns = patterns[:25]
	}

	return patterns
}

// detectContentPatterns finds common content patterns.
func (d *Detector) detectContentPatterns(specs []SpecInfo) []ContentPattern {
	var patterns []ContentPattern

	// Pattern definitions
	patternDefs := []struct {
		name  string
		regex *regexp.Regexp
		desc  string
	}{
		{"user_story", regexp.MustCompile(`(?i)^[-*]\s*As an? .+, I want .+, so that`), "User story format (As a... I want... So that...)"},
		{"acceptance_criteria", regexp.MustCompile(`(?i)^[-*]\s*Given .+, when .+, then`), "Gherkin-style acceptance criteria"},
		{"requirement_id", regexp.MustCompile(`\[REQ-\w+-\d+\]`), "Requirement ID references"},
		{"priority_marker", regexp.MustCompile(`(?i)\b(MUST|SHALL|SHOULD|MAY|WILL NOT|SHALL NOT)\b`), "RFC 2119 priority keywords"},
		{"metric_definition", regexp.MustCompile(`(?i)(KPI|metric|measure|target):.+\d+`), "Quantitative metrics"},
		{"api_reference", regexp.MustCompile(`(?i)(GET|POST|PUT|DELETE|PATCH)\s+/\w+`), "REST API endpoint references"},
		{"data_model", regexp.MustCompile(`(?i)^\s*[-*]\s*\w+:\s*(string|int|bool|uuid|timestamp)`), "Data model field definitions"},
	}

	for _, def := range patternDefs {
		var instances []PatternInstance
		for _, spec := range specs {
			for i, line := range spec.Lines {
				if def.regex.MatchString(line) {
					instances = append(instances, PatternInstance{
						Project:  spec.Project,
						SpecType: spec.SpecType,
						Line:     i + 1,
						Text:     truncate(strings.TrimSpace(line), 100),
					})
				}
			}
		}

		if len(instances) >= 2 {
			// Limit examples
			examples := instances
			if len(examples) > 5 {
				examples = examples[:5]
			}

			patterns = append(patterns, ContentPattern{
				Type:        def.name,
				Pattern:     def.regex.String(),
				Description: def.desc,
				Occurrences: len(instances),
				Examples:    examples,
			})
		}
	}

	// Sort by occurrences
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Occurrences > patterns[j].Occurrences
	})

	return patterns
}

// detectAntiPatterns finds specification anti-patterns.
func (d *Detector) detectAntiPatterns(specs []SpecInfo) []AntiPattern {
	var antiPatterns []AntiPattern

	// Anti-pattern definitions
	antiPatternDefs := []struct {
		name     string
		check    func(spec SpecInfo) []PatternInstance
		severity string
		desc     string
		suggest  string
	}{
		{
			"empty_section",
			func(spec SpecInfo) []PatternInstance {
				var instances []PatternInstance
				for i, line := range spec.Lines {
					if strings.HasPrefix(line, "## ") {
						// Check if next non-empty line is another header
						for j := i + 1; j < len(spec.Lines); j++ {
							nextLine := strings.TrimSpace(spec.Lines[j])
							if nextLine == "" {
								continue
							}
							if strings.HasPrefix(nextLine, "#") {
								instances = append(instances, PatternInstance{
									Project:  spec.Project,
									SpecType: spec.SpecType,
									Line:     i + 1,
									Text:     line,
								})
							}
							break
						}
					}
				}
				return instances
			},
			"medium",
			"Empty section with no content",
			"Add content or remove the section header",
		},
		{
			"vague_requirement",
			func(spec SpecInfo) []PatternInstance {
				var instances []PatternInstance
				vagueWords := regexp.MustCompile(`(?i)\b(fast|good|nice|easy|simple|better|improved|enhanced|etc|various|several)\b`)
				for i, line := range spec.Lines {
					if (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ")) && vagueWords.MatchString(line) {
						instances = append(instances, PatternInstance{
							Project:  spec.Project,
							SpecType: spec.SpecType,
							Line:     i + 1,
							Text:     truncate(line, 80),
						})
					}
				}
				return instances
			},
			"high",
			"Vague or unmeasurable requirement",
			"Use specific, measurable criteria instead of subjective terms",
		},
		{
			"missing_acceptance_criteria",
			func(spec SpecInfo) []PatternInstance {
				// Check if PRD/UXD has user stories but no acceptance criteria
				if spec.SpecType != "prd" && spec.SpecType != "uxd" {
					return nil
				}
				hasUserStories := false
				hasAcceptance := false
				for _, line := range spec.Lines {
					lineLower := strings.ToLower(line)
					if strings.Contains(lineLower, "as a") || strings.Contains(lineLower, "user stor") {
						hasUserStories = true
					}
					if strings.Contains(lineLower, "acceptance") || strings.Contains(lineLower, "given") {
						hasAcceptance = true
					}
				}
				if hasUserStories && !hasAcceptance {
					return []PatternInstance{{
						Project:  spec.Project,
						SpecType: spec.SpecType,
						Line:     0,
						Text:     "Spec has user stories but no acceptance criteria",
					}}
				}
				return nil
			},
			"high",
			"User stories without acceptance criteria",
			"Add Given/When/Then acceptance criteria for each user story",
		},
		{
			"todo_marker",
			func(spec SpecInfo) []PatternInstance {
				var instances []PatternInstance
				todoPattern := regexp.MustCompile(`(?i)\b(TODO|FIXME|TBD|XXX)\b`)
				for i, line := range spec.Lines {
					if todoPattern.MatchString(line) {
						instances = append(instances, PatternInstance{
							Project:  spec.Project,
							SpecType: spec.SpecType,
							Line:     i + 1,
							Text:     truncate(line, 80),
						})
					}
				}
				return instances
			},
			"low",
			"Incomplete content markers",
			"Complete the TODO items before finalizing the spec",
		},
	}

	for _, def := range antiPatternDefs {
		var allInstances []PatternInstance
		for _, spec := range specs {
			instances := def.check(spec)
			allInstances = append(allInstances, instances...)
		}

		if len(allInstances) > 0 {
			// Limit instances
			if len(allInstances) > 10 {
				allInstances = allInstances[:10]
			}

			antiPatterns = append(antiPatterns, AntiPattern{
				Type:        def.name,
				Severity:    def.severity,
				Description: def.desc,
				Instances:   allInstances,
				Suggestion:  def.suggest,
			})
		}
	}

	return antiPatterns
}

// calculateSummary computes aggregate statistics.
func (d *Detector) calculateSummary(specs []SpecInfo, report *PatternReport) PatternSummary {
	// Find most common sections
	sectionCounts := make(map[string]int)
	for _, spec := range specs {
		for _, section := range spec.Sections {
			sectionCounts[strings.ToLower(section)]++
		}
	}

	type sectionCount struct {
		name  string
		count int
	}
	var sortedSections []sectionCount
	for name, count := range sectionCounts {
		sortedSections = append(sortedSections, sectionCount{name, count})
	}
	sort.Slice(sortedSections, func(i, j int) bool {
		return sortedSections[i].count > sortedSections[j].count
	})

	var commonSections []string
	for i, sc := range sortedSections {
		if i >= 10 {
			break
		}
		commonSections = append(commonSections, sc.name)
	}

	// Pattern counts by type
	patternsByType := make(map[string]int)
	for _, p := range report.StructuralPatterns {
		patternsByType[p.Type]++
	}
	for _, p := range report.ContentPatterns {
		patternsByType[p.Type]++
	}

	// Calculate quality score (fewer anti-patterns = higher score)
	qualityScore := 100.0
	for _, ap := range report.AntiPatterns {
		count := len(ap.Instances)
		switch ap.Severity {
		case "high":
			qualityScore -= float64(count) * 5
		case "medium":
			qualityScore -= float64(count) * 2
		case "low":
			qualityScore -= float64(count) * 0.5
		}
	}
	if qualityScore < 0 {
		qualityScore = 0
	}

	return PatternSummary{
		TotalSpecs:       len(specs),
		StructuralCount:  len(report.StructuralPatterns),
		ContentCount:     len(report.ContentPatterns),
		AntiPatternCount: len(report.AntiPatterns),
		CommonSections:   commonSections,
		PatternsByType:   patternsByType,
		QualityScore:     qualityScore,
	}
}

func unique(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	sort.Strings(result)
	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
