// Package reuse provides requirements reuse tracking across projects.
package reuse

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Tracker identifies similar requirements across projects.
type Tracker struct {
	specsDir string
}

// NewTracker creates a new reuse tracker.
func NewTracker(specsDir string) *Tracker {
	return &Tracker{specsDir: specsDir}
}

// ReuseReport contains the analysis results.
type ReuseReport struct {
	SimilarRequirements []SimilarGroup   `json:"similar_requirements"`
	DuplicatePatterns   []DuplicateGroup `json:"duplicate_patterns"`
	ReuseCandidates     []ReuseCandidate `json:"reuse_candidates"`
	Summary             ReuseSummary     `json:"summary"`
}

// SimilarGroup represents requirements with similar wording.
type SimilarGroup struct {
	Pattern     string           `json:"pattern"`
	Description string           `json:"description"`
	Items       []RequirementRef `json:"items"`
	Similarity  float64          `json:"similarity"`
}

// RequirementRef references a specific requirement in a spec.
type RequirementRef struct {
	Project  string `json:"project"`
	SpecType string `json:"spec_type"`
	Line     int    `json:"line"`
	Text     string `json:"text"`
}

// DuplicateGroup represents exact or near-exact duplicates.
type DuplicateGroup struct {
	Text     string           `json:"text"`
	Count    int              `json:"count"`
	Projects []string         `json:"projects"`
	Refs     []RequirementRef `json:"refs"`
}

// ReuseCandidate suggests potential reuse opportunities.
type ReuseCandidate struct {
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Items       []RequirementRef `json:"items"`
	Suggestion  string           `json:"suggestion"`
	Priority    string           `json:"priority"` // high, medium, low
}

// ReuseSummary provides aggregate statistics.
type ReuseSummary struct {
	TotalRequirements   int `json:"total_requirements"`
	UniquePatterns      int `json:"unique_patterns"`
	DuplicateCount      int `json:"duplicate_count"`
	ReuseCandidateCount int `json:"reuse_candidate_count"`
	ProjectsAnalyzed    int `json:"projects_analyzed"`
}

// Requirement represents an extracted requirement.
type Requirement struct {
	Project  string
	SpecType string
	Line     int
	Text     string
	Keywords []string
}

// Analyze performs the reuse analysis.
func (t *Tracker) Analyze() (*ReuseReport, error) {
	// Extract all requirements
	requirements, err := t.extractRequirements()
	if err != nil {
		return nil, err
	}

	report := &ReuseReport{
		SimilarRequirements: []SimilarGroup{},
		DuplicatePatterns:   []DuplicateGroup{},
		ReuseCandidates:     []ReuseCandidate{},
	}

	// Find duplicates and similar requirements
	report.DuplicatePatterns = t.findDuplicates(requirements)
	report.SimilarRequirements = t.findSimilar(requirements)
	report.ReuseCandidates = t.identifyReuseCandidates(report)

	// Calculate summary
	projects := make(map[string]bool)
	for _, r := range requirements {
		projects[r.Project] = true
	}

	report.Summary = ReuseSummary{
		TotalRequirements:   len(requirements),
		UniquePatterns:      len(report.SimilarRequirements),
		DuplicateCount:      len(report.DuplicatePatterns),
		ReuseCandidateCount: len(report.ReuseCandidates),
		ProjectsAnalyzed:    len(projects),
	}

	return report, nil
}

// extractRequirements extracts requirements from all spec files.
func (t *Tracker) extractRequirements() ([]Requirement, error) {
	var requirements []Requirement

	// Patterns that indicate requirements
	reqPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^[-*]\s*(MUST|SHALL|SHOULD|WILL)\s+`),
		regexp.MustCompile(`(?i)^[-*]\s*\[REQ-`),
		regexp.MustCompile(`(?i)^[-*]\s*The system (must|shall|should|will)`),
		regexp.MustCompile(`(?i)^[-*]\s*Users (must|shall|should|will|can)`),
	}

	err := filepath.Walk(t.specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		relPath, _ := filepath.Rel(t.specsDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 1 {
			return nil
		}

		project := parts[0]
		specType := strings.TrimSuffix(info.Name(), ".md")

		content, err := os.ReadFile(path) //nolint:gosec // G122: User-provided specs directory for internal analysis
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)

			// Check if line looks like a requirement
			isReq := false
			for _, pattern := range reqPatterns {
				if pattern.MatchString(line) {
					isReq = true
					break
				}
			}

			// Also check for bullet points in requirements sections
			if !isReq && (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ")) {
				if i > 0 {
					// Check if we're in a requirements section
					for j := i - 1; j >= 0 && j > i-10; j-- {
						prevLine := strings.ToLower(lines[j])
						if strings.Contains(prevLine, "requirement") ||
							strings.Contains(prevLine, "user stor") ||
							strings.Contains(prevLine, "acceptance") {
							isReq = true
							break
						}
						if strings.HasPrefix(lines[j], "##") {
							break
						}
					}
				}
			}

			if isReq {
				text := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
				requirements = append(requirements, Requirement{
					Project:  project,
					SpecType: specType,
					Line:     i + 1,
					Text:     text,
					Keywords: extractKeywords(text),
				})
			}
		}

		return nil
	})

	return requirements, err
}

// extractKeywords extracts significant words from text.
func extractKeywords(text string) []string {
	// Remove common words and extract meaningful terms
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "as": true, "by": true, "at": true, "from": true,
		"must": true, "shall": true, "should": true, "will": true, "can": true,
		"that": true, "this": true, "it": true, "its": true,
	}

	words := regexp.MustCompile(`\w+`).FindAllString(strings.ToLower(text), -1)
	var keywords []string
	seen := make(map[string]bool)

	for _, word := range words {
		if len(word) > 2 && !stopWords[word] && !seen[word] {
			keywords = append(keywords, word)
			seen[word] = true
		}
	}

	return keywords
}

// findDuplicates finds exact or near-exact duplicate requirements.
func (t *Tracker) findDuplicates(requirements []Requirement) []DuplicateGroup {
	// Normalize and group by text
	normalized := make(map[string][]Requirement)

	for _, req := range requirements {
		// Normalize: lowercase, remove extra whitespace
		norm := strings.ToLower(strings.Join(strings.Fields(req.Text), " "))
		// Remove common prefixes
		norm = regexp.MustCompile(`^(the system |users? |the application )`).ReplaceAllString(norm, "")
		normalized[norm] = append(normalized[norm], req)
	}

	var groups []DuplicateGroup
	for text, reqs := range normalized {
		if len(reqs) < 2 {
			continue
		}

		// Check if from multiple projects
		projects := make(map[string]bool)
		var refs []RequirementRef
		for _, r := range reqs {
			projects[r.Project] = true
			refs = append(refs, RequirementRef{
				Project:  r.Project,
				SpecType: r.SpecType,
				Line:     r.Line,
				Text:     r.Text,
			})
		}

		if len(projects) >= 2 {
			projectList := make([]string, 0, len(projects))
			for p := range projects {
				projectList = append(projectList, p)
			}
			sort.Strings(projectList)

			groups = append(groups, DuplicateGroup{
				Text:     text,
				Count:    len(reqs),
				Projects: projectList,
				Refs:     refs,
			})
		}
	}

	// Sort by count descending
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Count > groups[j].Count
	})

	// Limit to top 20
	if len(groups) > 20 {
		groups = groups[:20]
	}

	return groups
}

// findSimilar finds requirements with similar keywords.
func (t *Tracker) findSimilar(requirements []Requirement) []SimilarGroup {
	// Group by keyword overlap
	var groups []SimilarGroup

	// Build keyword -> requirements index
	keywordIndex := make(map[string][]int)
	for i, req := range requirements {
		for _, kw := range req.Keywords {
			keywordIndex[kw] = append(keywordIndex[kw], i)
		}
	}

	// Find requirements that share multiple keywords
	pairScores := make(map[string]float64)
	pairRefs := make(map[string][]int)

	for _, indices := range keywordIndex {
		if len(indices) < 2 {
			continue
		}
		for i := 0; i < len(indices); i++ {
			for j := i + 1; j < len(indices); j++ {
				// Only count if from different projects
				if requirements[indices[i]].Project == requirements[indices[j]].Project {
					continue
				}
				key := pairKey(indices[i], indices[j])
				pairScores[key]++
				if len(pairRefs[key]) == 0 {
					pairRefs[key] = []int{indices[i], indices[j]}
				}
			}
		}
	}

	// Create groups for high-scoring pairs
	for key, score := range pairScores {
		if score < 3 {
			continue
		}

		refs := pairRefs[key]
		if len(refs) != 2 {
			continue
		}

		r1 := requirements[refs[0]]
		r2 := requirements[refs[1]]

		// Calculate similarity
		similarity := score / float64(max(len(r1.Keywords), len(r2.Keywords)))

		if similarity < 0.3 {
			continue
		}

		groups = append(groups, SimilarGroup{
			Pattern:     strings.Join(commonKeywords(r1.Keywords, r2.Keywords), ", "),
			Description: "Requirements with similar terminology",
			Items: []RequirementRef{
				{Project: r1.Project, SpecType: r1.SpecType, Line: r1.Line, Text: r1.Text},
				{Project: r2.Project, SpecType: r2.SpecType, Line: r2.Line, Text: r2.Text},
			},
			Similarity: similarity,
		})
	}

	// Sort by similarity descending
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Similarity > groups[j].Similarity
	})

	// Limit to top 15
	if len(groups) > 15 {
		groups = groups[:15]
	}

	return groups
}

func pairKey(i, j int) string {
	if i > j {
		i, j = j, i
	}
	return strconv.Itoa(i) + "|" + strconv.Itoa(j)
}

func commonKeywords(a, b []string) []string {
	set := make(map[string]bool)
	for _, k := range a {
		set[k] = true
	}
	var common []string
	for _, k := range b {
		if set[k] {
			common = append(common, k)
		}
	}
	return common
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// identifyReuseCandidates suggests reuse opportunities.
func (t *Tracker) identifyReuseCandidates(report *ReuseReport) []ReuseCandidate {
	var candidates []ReuseCandidate

	// Duplicates with 3+ occurrences should be consolidated
	for _, dup := range report.DuplicatePatterns {
		if dup.Count >= 3 {
			candidates = append(candidates, ReuseCandidate{
				Type:        "consolidate",
				Description: "Identical requirement appears in " + strings.Join(dup.Projects, ", "),
				Items:       dup.Refs,
				Suggestion:  "Consider extracting to a shared requirements library or CONSTITUTION.md",
				Priority:    "high",
			})
		}
	}

	// High-similarity pairs suggest shared patterns
	for _, sim := range report.SimilarRequirements {
		if sim.Similarity >= 0.6 {
			candidates = append(candidates, ReuseCandidate{
				Type:        "template",
				Description: "Similar requirements could use a common template",
				Items:       sim.Items,
				Suggestion:  "Consider creating a parameterized requirement template",
				Priority:    "medium",
			})
		}
	}

	return candidates
}
