package drift

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

// Analyzer performs spec-to-code comparison.
type Analyzer struct{}

// SpecRequirement represents a requirement extracted from spec.md.
type SpecRequirement struct {
	ID          string
	Description string
	Type        string // api, data, feature, etc.
	Section     string // Section of spec where found
}

// CodeImplementation represents a feature found in the codebase.
type CodeImplementation struct {
	ID          string
	Description string
	Type        string // api, data, feature
	File        string
	Line        int
}

// ExtractRequirements parses spec.md and extracts requirements.
func (a *Analyzer) ExtractRequirements(specContent string) []SpecRequirement {
	var requirements []SpecRequirement

	// Extract functional requirements
	requirements = append(requirements, a.extractFunctionalRequirements(specContent)...)

	// Extract API requirements
	requirements = append(requirements, a.extractAPIRequirements(specContent)...)

	// Extract data requirements
	requirements = append(requirements, a.extractDataRequirements(specContent)...)

	// Extract task items
	requirements = append(requirements, a.extractTasks(specContent)...)

	return requirements
}

// extractFunctionalRequirements finds FR-XXX or REQ-XXX patterns.
func (a *Analyzer) extractFunctionalRequirements(content string) []SpecRequirement {
	var reqs []SpecRequirement

	// Match requirement IDs with descriptions
	reqRE := regexp.MustCompile(`(?m)^[-*]\s*\*?\*?(FR|REQ|FUNC)-?(\d{3})\*?\*?:?\s*(.+)$`)
	matches := reqRE.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			reqs = append(reqs, SpecRequirement{
				ID:          match[1] + "-" + match[2],
				Description: strings.TrimSpace(match[3]),
				Type:        "feature",
				Section:     "Functional Requirements",
			})
		}
	}

	// Also look for numbered requirements in sections
	sectionRE := regexp.MustCompile(`(?s)##\s*(?:Functional\s+)?Requirements?\s*\n(.*?)(?:\n##|\z)`)
	sectionMatches := sectionRE.FindAllStringSubmatch(content, -1)

	for _, section := range sectionMatches {
		if len(section) < 2 {
			continue
		}

		// Extract numbered items
		itemRE := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
		items := itemRE.FindAllStringSubmatch(section[1], -1)

		for i, item := range items {
			if len(item) >= 2 {
				id := fmt.Sprintf("REQ-%03d", i+1)
				// Skip if we already have this from the ID pattern
				found := false
				for _, r := range reqs {
					if strings.Contains(item[1], r.ID) {
						found = true
						break
					}
				}
				if !found {
					reqs = append(reqs, SpecRequirement{
						ID:          id,
						Description: strings.TrimSpace(item[1]),
						Type:        "feature",
						Section:     "Requirements",
					})
				}
			}
		}
	}

	return reqs
}

// extractAPIRequirements finds API endpoints mentioned in spec.
func (a *Analyzer) extractAPIRequirements(content string) []SpecRequirement {
	var reqs []SpecRequirement

	// Look for endpoint patterns: GET /path, POST /path, etc.
	apiRE := regexp.MustCompile(`(?m)(GET|POST|PUT|PATCH|DELETE)\s+(/[^\s\n]+)`)
	matches := apiRE.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 3 {
			key := match[1] + " " + match[2]
			if seen[key] {
				continue
			}
			seen[key] = true

			reqs = append(reqs, SpecRequirement{
				ID:          fmt.Sprintf("API-%s-%s", match[1], sanitizePath(match[2])),
				Description: key,
				Type:        "api",
				Section:     "API",
			})
		}
	}

	return reqs
}

// extractDataRequirements finds data model/entity mentions.
func (a *Analyzer) extractDataRequirements(content string) []SpecRequirement {
	var reqs []SpecRequirement

	// Look for entity/model definitions in tables or lists
	entityRE := regexp.MustCompile(`(?i)(?:entity|model|table):\s*(\w+)`)
	matches := entityRE.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 2 {
			name := match[1]
			if seen[name] {
				continue
			}
			seen[name] = true

			reqs = append(reqs, SpecRequirement{
				ID:          "DATA-" + name,
				Description: name + " entity/model",
				Type:        "data",
				Section:     "Data Model",
			})
		}
	}

	return reqs
}

// extractTasks finds task items from the spec.
func (a *Analyzer) extractTasks(content string) []SpecRequirement {
	var reqs []SpecRequirement

	// Match task checkboxes
	taskRE := regexp.MustCompile(`(?m)^[-*]\s*\[[ x]\]\s*(.+)$`)
	matches := taskRE.FindAllStringSubmatch(content, -1)

	for i, match := range matches {
		if len(match) >= 2 {
			reqs = append(reqs, SpecRequirement{
				ID:          fmt.Sprintf("TASK-%03d", i+1),
				Description: strings.TrimSpace(match[1]),
				Type:        "task",
				Section:     "Tasks",
			})
		}
	}

	return reqs
}

// ExtractImplementations analyzes codebase context for implementations.
func (a *Analyzer) ExtractImplementations(ctx *context.AggregatedContext) []CodeImplementation {
	if ctx == nil {
		return nil
	}

	var implementations []CodeImplementation

	// Extract from code contexts
	for _, code := range ctx.CodeContexts() {
		// APIs
		for _, api := range code.APIs {
			for _, route := range api.Routes {
				implementations = append(implementations, CodeImplementation{
					ID:          fmt.Sprintf("API-%s-%s", route.Method, sanitizePath(route.Path)),
					Description: fmt.Sprintf("%s %s", route.Method, route.Path),
					Type:        "api",
					File:        api.Path,
				})
			}
		}
	}

	// Extract from graph contexts
	for _, graph := range ctx.GraphContexts() {
		// Requirements
		for _, req := range graph.Requirements {
			implementations = append(implementations, CodeImplementation{
				ID:          req.ID,
				Description: req.Title,
				Type:        "feature",
			})
		}
	}

	return implementations
}

// Compare finds drift between requirements and implementations.
func (a *Analyzer) Compare(requirements []SpecRequirement, implementations []CodeImplementation) []DriftItem {
	var items []DriftItem

	// Build lookup maps
	reqByID := make(map[string]SpecRequirement)
	for _, r := range requirements {
		reqByID[r.ID] = r
	}

	implByID := make(map[string]CodeImplementation)
	for _, i := range implementations {
		implByID[i.ID] = i
	}

	// Find unimplemented requirements
	for _, req := range requirements {
		if _, found := implByID[req.ID]; !found {
			// Also check by description similarity
			if !hasMatchingImpl(req, implementations) {
				items = append(items, DriftItem{
					ID:          fmt.Sprintf("DRIFT-%s", req.ID),
					Type:        DriftUnimplemented,
					Severity:    determineSeverity(req),
					Category:    categoryFromType(req.Type),
					Description: fmt.Sprintf("Requirement not implemented: %s", req.Description),
					SpecRef:     fmt.Sprintf("%s (%s)", req.ID, req.Section),
					Suggestion:  fmt.Sprintf("Implement %s as specified in %s", req.ID, req.Section),
				})
			}
		}
	}

	// Find undocumented implementations
	for _, impl := range implementations {
		if _, found := reqByID[impl.ID]; !found {
			// Also check by description similarity
			if !hasMatchingReq(impl, requirements) {
				items = append(items, DriftItem{
					ID:          fmt.Sprintf("DRIFT-UNDOC-%s", impl.ID),
					Type:        DriftUndocumented,
					Severity:    SeverityLow,
					Category:    categoryFromType(impl.Type),
					Description: fmt.Sprintf("Implementation not in spec: %s", impl.Description),
					CodeRef:     impl.File,
					Suggestion:  "Add to spec.md or remove if not needed",
				})
			}
		}
	}

	return items
}

// hasMatchingImpl checks if a requirement has a matching implementation by description.
func hasMatchingImpl(req SpecRequirement, implementations []CodeImplementation) bool {
	reqWords := extractKeywords(req.Description)

	for _, impl := range implementations {
		implWords := extractKeywords(impl.Description)
		if similarity(reqWords, implWords) > 0.5 {
			return true
		}
	}
	return false
}

// hasMatchingReq checks if an implementation has a matching requirement by description.
func hasMatchingReq(impl CodeImplementation, requirements []SpecRequirement) bool {
	implWords := extractKeywords(impl.Description)

	for _, req := range requirements {
		reqWords := extractKeywords(req.Description)
		if similarity(reqWords, implWords) > 0.5 {
			return true
		}
	}
	return false
}

// extractKeywords returns significant words from a string.
func extractKeywords(s string) map[string]bool {
	words := make(map[string]bool)
	wordRE := regexp.MustCompile(`\b\w{3,}\b`)
	matches := wordRE.FindAllString(strings.ToLower(s), -1)

	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "that": true,
		"with": true, "from": true, "this": true, "have": true,
		"should": true, "must": true, "will": true, "can": true,
	}

	for _, word := range matches {
		if !stopWords[word] {
			words[word] = true
		}
	}
	return words
}

// similarity calculates Jaccard similarity between two word sets.
func similarity(a, b map[string]bool) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	intersection := 0
	for word := range a {
		if b[word] {
			intersection++
		}
	}

	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// sanitizePath creates a safe ID from a URL path.
func sanitizePath(path string) string {
	// Replace / with - and remove special chars
	result := strings.ReplaceAll(path, "/", "-")
	result = strings.ReplaceAll(result, "{", "")
	result = strings.ReplaceAll(result, "}", "")
	result = strings.Trim(result, "-")
	return result
}

// determineSeverity determines severity based on requirement type.
func determineSeverity(req SpecRequirement) Severity {
	// High severity for core features
	if strings.Contains(strings.ToLower(req.Description), "must") ||
		strings.Contains(strings.ToLower(req.Description), "critical") ||
		strings.Contains(strings.ToLower(req.Description), "security") {
		return SeverityHigh
	}

	// Medium for API requirements
	if req.Type == "api" {
		return SeverityMedium
	}

	// Low for tasks and data
	return SeverityLow
}

// categoryFromType maps requirement type to drift category.
func categoryFromType(t string) Category {
	switch t {
	case "api":
		return CategoryAPI
	case "data":
		return CategoryData
	case "feature":
		return CategoryOther
	default:
		return CategoryOther
	}
}
