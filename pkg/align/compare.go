package align

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/context"
)

// Requirement represents an extracted spec requirement.
type Requirement struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` // "feature", "api", "data", "behavior", etc.
	Description string   `json:"description"`
	Section     string   `json:"section"`     // Spec section path
	LineNumber  int      `json:"line_number"` // Line in spec
	Priority    string   `json:"priority"`    // "must", "should", "could"
	Keywords    []string `json:"keywords"`    // Key terms for matching
}

// Implementation represents a code implementation.
type Implementation struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` // "function", "endpoint", "model", etc.
	Name        string   `json:"name"`
	Description string   `json:"description"`
	FilePath    string   `json:"file_path"`
	LineNumber  int      `json:"line_number"`
	Keywords    []string `json:"keywords"`
}

// Comparator compares specs against implementations.
type Comparator struct{}

// NewComparator creates a new comparator.
func NewComparator() *Comparator {
	return &Comparator{}
}

// ExtractRequirements parses spec content to extract requirements.
func (c *Comparator) ExtractRequirements(specContent string) []Requirement {
	var requirements []Requirement

	lines := strings.Split(specContent, "\n")
	currentSection := ""
	sectionStack := []string{}

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track section hierarchy via headers
		if strings.HasPrefix(trimmed, "#") {
			level := countLeadingHashes(trimmed)
			title := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))

			// Adjust section stack based on level
			for len(sectionStack) >= level {
				sectionStack = sectionStack[:len(sectionStack)-1]
			}
			sectionStack = append(sectionStack, title)
			currentSection = strings.Join(sectionStack, " > ")
			continue
		}

		// Extract requirements from various patterns
		req := c.extractRequirementFromLine(trimmed, currentSection, lineNum+1)
		if req != nil {
			requirements = append(requirements, *req)
		}
	}

	return requirements
}

// extractRequirementFromLine attempts to extract a requirement from a line.
func (c *Comparator) extractRequirementFromLine(line, section string, lineNum int) *Requirement {
	// Skip empty lines and non-content
	if line == "" || strings.HasPrefix(line, "---") {
		return nil
	}

	// Pattern: bullet points with requirements
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
		content := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
		return c.parseRequirement(content, section, lineNum)
	}

	// Pattern: numbered lists
	numberedRegex := regexp.MustCompile(`^\d+\.\s+(.+)`)
	if matches := numberedRegex.FindStringSubmatch(line); len(matches) > 1 {
		return c.parseRequirement(matches[1], section, lineNum)
	}

	// Pattern: "MUST", "SHALL", "SHOULD" requirements (RFC 2119)
	rfcKeywords := []string{"MUST", "SHALL", "SHOULD", "MAY", "REQUIRED"}
	for _, kw := range rfcKeywords {
		if strings.Contains(strings.ToUpper(line), kw) {
			return c.parseRequirement(line, section, lineNum)
		}
	}

	// Pattern: API endpoint definitions
	apiRegex := regexp.MustCompile(`(?i)(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s+(/\S+)`)
	if apiRegex.MatchString(line) {
		return &Requirement{
			ID:          generateID("api", section, lineNum),
			Type:        "api",
			Description: line,
			Section:     section,
			LineNumber:  lineNum,
			Priority:    "must",
			Keywords:    extractKeywords(line),
		}
	}

	// Pattern: Feature descriptions with action verbs
	featureVerbs := []string{"provides", "enables", "allows", "supports", "implements", "displays", "shows", "creates", "updates", "deletes"}
	lineLower := strings.ToLower(line)
	for _, verb := range featureVerbs {
		if strings.Contains(lineLower, verb) {
			return c.parseRequirement(line, section, lineNum)
		}
	}

	return nil
}

// parseRequirement creates a requirement from content.
func (c *Comparator) parseRequirement(content, section string, lineNum int) *Requirement {
	priority := "should"
	contentUpper := strings.ToUpper(content)

	if strings.Contains(contentUpper, "MUST") || strings.Contains(contentUpper, "REQUIRED") || strings.Contains(contentUpper, "SHALL") {
		priority = "must"
	} else if strings.Contains(contentUpper, "COULD") || strings.Contains(contentUpper, "OPTIONAL") || strings.Contains(contentUpper, "MAY") {
		priority = "could"
	}

	reqType := inferRequirementType(content, section)

	return &Requirement{
		ID:          generateID(reqType, section, lineNum),
		Type:        reqType,
		Description: content,
		Section:     section,
		LineNumber:  lineNum,
		Priority:    priority,
		Keywords:    extractKeywords(content),
	}
}

// inferRequirementType guesses the requirement type from content.
func inferRequirementType(content, section string) string {
	contentLower := strings.ToLower(content)
	sectionLower := strings.ToLower(section)

	switch {
	case strings.Contains(sectionLower, "api") || strings.Contains(contentLower, "endpoint") || strings.Contains(contentLower, "route"):
		return "api"
	case strings.Contains(sectionLower, "data") || strings.Contains(contentLower, "schema") || strings.Contains(contentLower, "model"):
		return "data"
	case strings.Contains(sectionLower, "ui") || strings.Contains(sectionLower, "interface") || strings.Contains(contentLower, "display"):
		return "ui"
	case strings.Contains(sectionLower, "security") || strings.Contains(contentLower, "auth") || strings.Contains(contentLower, "permission"):
		return "security"
	case strings.Contains(sectionLower, "performance") || strings.Contains(contentLower, "latency") || strings.Contains(contentLower, "throughput"):
		return "performance"
	case strings.Contains(sectionLower, "infrastructure") || strings.Contains(contentLower, "deploy") || strings.Contains(contentLower, "scale"):
		return "infrastructure"
	default:
		return "feature"
	}
}

// ExtractImplementations parses context to extract implementations.
func (c *Comparator) ExtractImplementations(ctx *context.AggregatedContext) []Implementation {
	var implementations []Implementation

	if ctx == nil {
		return implementations
	}

	// Extract from code context
	for _, data := range ctx.Sources {
		if data.Code != nil {
			impls := c.extractFromCodeContext(data.Code)
			implementations = append(implementations, impls...)
		}
	}

	return implementations
}

// extractFromCodeContext extracts implementations from code analysis.
func (c *Comparator) extractFromCodeContext(code *context.CodeContext) []Implementation {
	var implementations []Implementation

	if code == nil {
		return implementations
	}

	// Extract from directory structure
	if code.Structure != nil {
		c.extractFromTreeNode(code.Structure, "", &implementations)
	}

	// Extract from API schemas if available
	for _, api := range code.APIs {
		for _, route := range api.Routes {
			implementations = append(implementations, Implementation{
				ID:          generateID("endpoint", route.Path, 0),
				Type:        "endpoint",
				Name:        route.Method + " " + route.Path,
				Description: route.Summary,
				FilePath:    api.Path,
				LineNumber:  0,
				Keywords:    extractKeywords(route.Path + " " + route.Summary + " " + route.OperationID),
			})
		}
	}

	return implementations
}

// extractFromTreeNode recursively extracts implementations from the file tree.
func (c *Comparator) extractFromTreeNode(node *context.TreeNode, path string, implementations *[]Implementation) {
	if node == nil {
		return
	}

	var currentPath string
	if path == "" {
		currentPath = node.Name
	} else {
		currentPath = path + "/" + node.Name
	}

	// Extract implementation hints from file names
	if node.Type == "file" {
		// Infer implementations from file names
		name := node.Name

		// Common patterns that indicate implementation
		if strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, ".ts") ||
			strings.HasSuffix(name, ".js") ||
			strings.HasSuffix(name, ".py") ||
			strings.HasSuffix(name, ".java") ||
			strings.HasSuffix(name, ".rs") {
			// Skip test files
			if !strings.Contains(name, "_test") && !strings.Contains(name, ".test.") && !strings.Contains(name, ".spec.") {
				baseName := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(
					strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".go"), ".ts"), ".js"),
					".py"), ".java"), ".rs")

				*implementations = append(*implementations, Implementation{
					ID:          generateID("file", currentPath, 0),
					Type:        "file",
					Name:        baseName,
					Description: currentPath,
					FilePath:    currentPath,
					LineNumber:  0,
					Keywords:    extractKeywords(baseName + " " + currentPath),
				})
			}
		}
	}

	// Recurse into children
	for _, child := range node.Children {
		c.extractFromTreeNode(child, currentPath, implementations)
	}
}

// Compare finds discrepancies between requirements and implementations.
func (c *Comparator) Compare(requirements []Requirement, implementations []Implementation, includeEvidence bool) []Discrepancy {
	var discrepancies []Discrepancy

	// Track which implementations are matched
	matchedImpls := make(map[string]bool)

	// Check each requirement for implementation
	for _, req := range requirements {
		matches := c.findMatches(req, implementations)

		if len(matches) == 0 {
			// Missing feature
			discrepancy := Discrepancy{
				ID:          "DISC-" + req.ID,
				Type:        DiscrepancyMissingFeature,
				Severity:    c.priorityToSeverity(req.Priority),
				Category:    c.typeToCategory(req.Type),
				Description: fmt.Sprintf("Requirement not implemented: %s", truncate(req.Description, 100)),
				SpecRef:     fmt.Sprintf("spec.md:%d (%s)", req.LineNumber, req.Section),
				Expected:    req.Description,
				Suggestion:  fmt.Sprintf("Implement: %s", req.Description),
			}
			if includeEvidence {
				discrepancy.Evidence = []Evidence{
					{Type: "spec_excerpt", Content: req.Description, Source: fmt.Sprintf("Line %d", req.LineNumber)},
				}
			}
			discrepancies = append(discrepancies, discrepancy)
		} else {
			// Mark implementations as matched
			for _, impl := range matches {
				matchedImpls[impl.ID] = true
			}

			// Check for partial implementation or divergence
			// This is a simplified check - a full implementation would do deeper analysis
			bestMatch := matches[0]
			if c.isPartialMatch(req, bestMatch) {
				discrepancies = append(discrepancies, Discrepancy{
					ID:          "DISC-PARTIAL-" + req.ID,
					Type:        DiscrepancyPartialImplementation,
					Severity:    SeverityMedium,
					Category:    c.typeToCategory(req.Type),
					Description: fmt.Sprintf("Partially implemented: %s", truncate(req.Description, 80)),
					SpecRef:     fmt.Sprintf("spec.md:%d", req.LineNumber),
					CodeRef:     fmt.Sprintf("%s:%d", bestMatch.FilePath, bestMatch.LineNumber),
					Expected:    req.Description,
					Actual:      bestMatch.Description,
				})
			}
		}
	}

	// Find undocumented implementations
	for _, impl := range implementations {
		if !matchedImpls[impl.ID] && c.isSignificantImplementation(impl) {
			discrepancies = append(discrepancies, Discrepancy{
				ID:          "DISC-UNDOC-" + impl.ID,
				Type:        DiscrepancyUndocumentedCode,
				Severity:    SeverityLow,
				Category:    c.typeToCategory(impl.Type),
				Description: fmt.Sprintf("Code without spec coverage: %s", impl.Name),
				CodeRef:     fmt.Sprintf("%s:%d", impl.FilePath, impl.LineNumber),
				Actual:      impl.Description,
				Suggestion:  "Add to spec or verify this is intentional",
			})
		}
	}

	return discrepancies
}

// findMatches finds implementations that match a requirement.
func (c *Comparator) findMatches(req Requirement, implementations []Implementation) []Implementation {
	var matches []Implementation

	for _, impl := range implementations {
		if c.matches(req, impl) {
			matches = append(matches, impl)
		}
	}

	return matches
}

// matches checks if a requirement and implementation match.
func (c *Comparator) matches(req Requirement, impl Implementation) bool {
	// Check keyword overlap
	reqKeywords := make(map[string]bool)
	for _, kw := range req.Keywords {
		reqKeywords[strings.ToLower(kw)] = true
	}

	matchCount := 0
	for _, kw := range impl.Keywords {
		if reqKeywords[strings.ToLower(kw)] {
			matchCount++
		}
	}

	// Require at least some keyword overlap
	minMatches := 1
	if len(req.Keywords) > 3 {
		minMatches = 2
	}

	if matchCount >= minMatches {
		return true
	}

	// Also check for substring matches
	reqLower := strings.ToLower(req.Description)
	implLower := strings.ToLower(impl.Name + " " + impl.Description)

	// Check if key parts of the requirement appear in implementation
	for _, kw := range req.Keywords {
		if len(kw) > 3 && strings.Contains(implLower, strings.ToLower(kw)) {
			return true
		}
	}

	// Check if implementation name appears in requirement
	if len(impl.Name) > 3 && strings.Contains(reqLower, strings.ToLower(impl.Name)) {
		return true
	}

	return false
}

// isPartialMatch checks if the match is partial.
func (c *Comparator) isPartialMatch(req Requirement, impl Implementation) bool {
	// Count keyword matches
	reqKeywords := make(map[string]bool)
	for _, kw := range req.Keywords {
		reqKeywords[strings.ToLower(kw)] = true
	}

	matchCount := 0
	for _, kw := range impl.Keywords {
		if reqKeywords[strings.ToLower(kw)] {
			matchCount++
		}
	}

	// If less than half of keywords match, consider it partial
	if len(req.Keywords) > 0 {
		matchRatio := float64(matchCount) / float64(len(req.Keywords))
		return matchRatio < 0.5
	}

	return false
}

// isSignificantImplementation checks if an implementation is significant enough to report.
func (c *Comparator) isSignificantImplementation(impl Implementation) bool {
	// Skip test files
	if strings.Contains(impl.FilePath, "_test.go") || strings.Contains(impl.FilePath, ".test.") {
		return false
	}

	// Skip private/internal helpers
	if impl.Type == "function" {
		firstChar := ""
		if len(impl.Name) > 0 {
			firstChar = string(impl.Name[0])
		}
		if firstChar == strings.ToLower(firstChar) {
			// unexported function - skip unless it's significant
			if len(impl.Name) < 10 {
				return false
			}
		}
	}

	return true
}

// priorityToSeverity converts requirement priority to discrepancy severity.
func (c *Comparator) priorityToSeverity(priority string) Severity {
	switch priority {
	case "must":
		return SeverityHigh
	case "should":
		return SeverityMedium
	case "could":
		return SeverityLow
	default:
		return SeverityMedium
	}
}

// typeToCategory converts requirement/implementation type to category.
func (c *Comparator) typeToCategory(t string) Category {
	switch t {
	case "api", "endpoint":
		return CategoryAPI
	case "data", "type", "model", "schema":
		return CategoryData
	case "ui", "interface":
		return CategoryUI
	case "security", "auth":
		return CategorySecurity
	case "performance", "perf":
		return CategoryPerf
	case "infrastructure", "infra":
		return CategoryInfra
	case "behavior", "function":
		return CategoryBehavior
	default:
		return CategoryOther
	}
}

// Helper functions

func countLeadingHashes(s string) int {
	count := 0
	for _, c := range s {
		if c == '#' {
			count++
		} else {
			break
		}
	}
	return count
}

func generateID(prefix, context string, lineNum int) string {
	data := fmt.Sprintf("%s:%s:%d", prefix, context, lineNum)
	hash := sha256.Sum256([]byte(data))
	// Use at most 3 characters from prefix
	prefixLen := len(prefix)
	if prefixLen > 3 {
		prefixLen = 3
	}
	return fmt.Sprintf("%s-%x", strings.ToUpper(prefix[:prefixLen]), hash[:4])
}

func extractKeywords(text string) []string {
	// Remove punctuation and split into words
	cleaned := regexp.MustCompile(`[^\w\s]`).ReplaceAllString(text, " ")
	words := strings.Fields(cleaned)

	// Filter out common stop words and short words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true, "is": true, "was": true,
		"are": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true, "may": true,
		"might": true, "must": true, "shall": true, "can": true,
		"this": true, "that": true, "these": true, "those": true,
		"it": true, "its": true, "they": true, "them": true, "their": true,
	}

	var keywords []string
	seen := make(map[string]bool)
	for _, word := range words {
		lower := strings.ToLower(word)
		if len(word) > 2 && !stopWords[lower] && !seen[lower] {
			keywords = append(keywords, lower)
			seen[lower] = true
		}
	}

	return keywords
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
