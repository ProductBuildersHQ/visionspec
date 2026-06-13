// Package search provides full-text search across specification files.
package search

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Searcher provides full-text search over spec files.
type Searcher struct {
	specsDir string
}

// NewSearcher creates a new searcher for a specs directory.
func NewSearcher(specsDir string) *Searcher {
	return &Searcher{specsDir: specsDir}
}

// SearchResult represents a single search hit.
type SearchResult struct {
	Project    string  `json:"project"`
	SpecType   string  `json:"spec_type"`
	FilePath   string  `json:"file_path"`
	Line       int     `json:"line"`
	Score      float64 `json:"score"`
	Snippet    string  `json:"snippet"`
	Context    string  `json:"context,omitempty"`
	MatchCount int     `json:"match_count"`
}

// SearchOutput contains the full search response.
type SearchOutput struct {
	Query       string          `json:"query"`
	TotalHits   int             `json:"total_hits"`
	Results     []*SearchResult `json:"results"`
	Took        string          `json:"took"`
	ByProject   map[string]int  `json:"by_project,omitempty"`
	BySpecType  map[string]int  `json:"by_spec_type,omitempty"`
	Truncated   bool            `json:"truncated,omitempty"`
}

// SearchOptions configures search behavior.
type SearchOptions struct {
	Limit       int      // Maximum results to return (default 20)
	Projects    []string // Filter by project names (empty = all)
	SpecTypes   []string // Filter by spec types (mrd, prd, trd, etc.)
	CaseSensitive bool   // Case sensitive matching
	Regex       bool     // Treat query as regex
	ContextLines int     // Lines of context around match (default 1)
}

// DefaultSearchOptions returns sensible defaults.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Limit:        20,
		ContextLines: 1,
	}
}

// Search performs a full-text search across specs.
func (s *Searcher) Search(query string, opts SearchOptions) (*SearchOutput, error) {
	startTime := time.Now()

	if opts.Limit <= 0 {
		opts.Limit = 20
	}

	output := &SearchOutput{
		Query:      query,
		Results:    []*SearchResult{},
		ByProject:  make(map[string]int),
		BySpecType: make(map[string]int),
	}

	// Build search pattern
	var pattern *regexp.Regexp
	var err error
	if opts.Regex {
		if opts.CaseSensitive {
			pattern, err = regexp.Compile(query)
		} else {
			pattern, err = regexp.Compile("(?i)" + query)
		}
		if err != nil {
			return nil, err
		}
	}

	// Walk specs directory
	err = filepath.Walk(s.specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip hidden and special files
		if strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(info.Name(), "_") {
			return nil
		}

		// Parse project and spec type from path
		relPath, _ := filepath.Rel(s.specsDir, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 1 {
			return nil
		}

		project := parts[0]
		specType := strings.TrimSuffix(info.Name(), ".md")

		// Apply filters
		if len(opts.Projects) > 0 && !contains(opts.Projects, project) {
			return nil
		}
		if len(opts.SpecTypes) > 0 && !contains(opts.SpecTypes, specType) {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Search for matches
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			var matches bool
			var matchCount int

			if pattern != nil {
				matches = pattern.MatchString(line)
				if matches {
					matchCount = len(pattern.FindAllString(line, -1))
				}
			} else {
				searchLine := line
				searchQuery := query
				if !opts.CaseSensitive {
					searchLine = strings.ToLower(line)
					searchQuery = strings.ToLower(query)
				}
				matchCount = strings.Count(searchLine, searchQuery)
				matches = matchCount > 0
			}

			if matches {
				// Build context
				contextStart := i - opts.ContextLines
				if contextStart < 0 {
					contextStart = 0
				}
				contextEnd := i + opts.ContextLines + 1
				if contextEnd > len(lines) {
					contextEnd = len(lines)
				}
				context := strings.Join(lines[contextStart:contextEnd], "\n")

				result := &SearchResult{
					Project:    project,
					SpecType:   specType,
					FilePath:   relPath,
					Line:       i + 1,
					Score:      float64(matchCount),
					Snippet:    truncate(strings.TrimSpace(line), 150),
					Context:    context,
					MatchCount: matchCount,
				}

				output.Results = append(output.Results, result)
				output.ByProject[project]++
				output.BySpecType[specType]++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by score (match count) descending
	sort.Slice(output.Results, func(i, j int) bool {
		return output.Results[i].Score > output.Results[j].Score
	})

	output.TotalHits = len(output.Results)

	// Apply limit
	if len(output.Results) > opts.Limit {
		output.Results = output.Results[:opts.Limit]
		output.Truncated = true
	}

	output.Took = time.Since(startTime).String()

	return output, nil
}

// contains checks if a string is in a slice.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// truncate truncates a string to maxLen with ellipsis.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
