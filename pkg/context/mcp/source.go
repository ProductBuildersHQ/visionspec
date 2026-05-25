package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ctx "github.com/plexusone/multispec/pkg/context"
)

// Source fetches context from an MCP server.
type Source struct {
	name   string
	config ctx.MCPServerConfig
	client *Client
}

// NewSource creates a new MCP source from configuration.
func NewSource(name string, cfg ctx.MCPServerConfig) (*Source, error) {
	if cfg.Command == "" {
		return nil, fmt.Errorf("MCP server command is required")
	}

	return &Source{
		name:   name,
		config: cfg,
	}, nil
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return fmt.Sprintf("mcp:%s", s.name)
}

// Type returns the source type.
func (s *Source) Type() ctx.SourceType {
	return ctx.SourceTypeMCP
}

// String returns a human-readable description.
func (s *Source) String() string {
	return fmt.Sprintf("MCP Server: %s (%s)", s.name, s.config.Command)
}

// Fetch retrieves context from the MCP server.
func (s *Source) Fetch(c context.Context) (*ctx.ContextData, error) {
	start := time.Now()

	// Start client if not already running
	if s.client == nil {
		client, err := NewClient(s.config.Command, s.config.Args, s.config.Env)
		if err != nil {
			return nil, fmt.Errorf("starting MCP server: %w", err)
		}
		s.client = client

		// Initialize the connection
		timeout := s.config.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		initCtx, cancel := context.WithTimeout(c, timeout)
		defer cancel()

		if err := s.client.Initialize(initCtx); err != nil {
			s.client.Close()
			s.client = nil
			return nil, fmt.Errorf("initializing MCP server: %w", err)
		}
	}

	// Get available tools
	tools, err := s.client.ListTools(c)
	if err != nil {
		return nil, fmt.Errorf("listing tools: %w", err)
	}

	external := &ctx.ExternalContext{
		ServerName: s.name,
		ServerType: detectServerType(s.name, tools),
	}

	// Fetch context based on available tools
	if err := s.fetchContext(c, tools, external); err != nil {
		return &ctx.ContextData{
			Source:    s.Name(),
			Type:      ctx.SourceTypeMCP,
			FetchedAt: time.Now(),
			Duration:  time.Since(start),
			External:  external,
			Errors:    []string{err.Error()},
		}, nil
	}

	return &ctx.ContextData{
		Source:    s.Name(),
		Type:      ctx.SourceTypeMCP,
		FetchedAt: time.Now(),
		Duration:  time.Since(start),
		External:  external,
		Summary:   generateMCPSummary(external),
	}, nil
}

// Close shuts down the MCP client.
func (s *Source) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func detectServerType(name string, tools []Tool) string {
	nameLower := strings.ToLower(name)

	// Check name hints
	if strings.Contains(nameLower, "jira") {
		return "jira"
	}
	if strings.Contains(nameLower, "confluence") {
		return "confluence"
	}
	if strings.Contains(nameLower, "google") {
		return "google"
	}
	if strings.Contains(nameLower, "linear") {
		return "linear"
	}
	if strings.Contains(nameLower, "notion") {
		return "notion"
	}

	// Check tool names for hints
	for _, tool := range tools {
		toolLower := strings.ToLower(tool.Name)
		if strings.Contains(toolLower, "jira") || strings.Contains(toolLower, "issue") {
			return "jira"
		}
		if strings.Contains(toolLower, "confluence") || strings.Contains(toolLower, "page") {
			return "confluence"
		}
	}

	return "generic"
}

func (s *Source) fetchContext(c context.Context, tools []Tool, external *ctx.ExternalContext) error {
	// Build tool map for quick lookup
	toolMap := make(map[string]Tool)
	for _, t := range tools {
		toolMap[t.Name] = t
	}

	// Try to fetch issues (Jira, Linear, etc.)
	if err := s.fetchIssues(c, toolMap, external); err != nil {
		// Non-fatal, continue with other fetches
		fmt.Printf("  Warning: fetching issues: %v\n", err)
	}

	// Try to fetch pages (Confluence, Notion, etc.)
	if err := s.fetchPages(c, toolMap, external); err != nil {
		fmt.Printf("  Warning: fetching pages: %v\n", err)
	}

	// Try to fetch documents (Google Docs, etc.)
	if err := s.fetchDocuments(c, toolMap, external); err != nil {
		fmt.Printf("  Warning: fetching documents: %v\n", err)
	}

	return nil
}

func (s *Source) fetchIssues(c context.Context, tools map[string]Tool, external *ctx.ExternalContext) error {
	// Look for issue-fetching tools
	issueTools := []string{
		"search_issues",
		"list_issues",
		"get_issues",
		"jira_search",
		"linear_issues",
	}

	for _, toolName := range issueTools {
		if _, ok := tools[toolName]; ok {
			// Get config for query
			args := make(map[string]any)
			if s.config.Config != nil {
				if jql, ok := s.config.Config["jql"]; ok {
					args["jql"] = jql
				}
				if query, ok := s.config.Config["query"]; ok {
					args["query"] = query
				}
				if project, ok := s.config.Config["project"]; ok {
					args["project"] = project
				}
			}

			// Default to recent issues if no query specified
			if len(args) == 0 {
				args["maxResults"] = 50
			}

			content, err := s.client.CallTool(c, toolName, args)
			if err != nil {
				continue
			}

			// Parse issues from response
			for _, item := range content {
				if item.Type == "text" {
					issues := parseIssuesFromText(item.Text)
					external.Issues = append(external.Issues, issues...)
				}
			}

			if len(external.Issues) > 0 {
				return nil
			}
		}
	}

	return nil
}

func (s *Source) fetchPages(c context.Context, tools map[string]Tool, external *ctx.ExternalContext) error {
	// Look for page-fetching tools
	pageTools := []string{
		"search_pages",
		"list_pages",
		"get_pages",
		"confluence_search",
		"notion_search",
	}

	for _, toolName := range pageTools {
		if _, ok := tools[toolName]; ok {
			args := make(map[string]any)
			if s.config.Config != nil {
				if space, ok := s.config.Config["space"]; ok {
					args["space"] = space
				}
				if query, ok := s.config.Config["query"]; ok {
					args["query"] = query
				}
			}

			if len(args) == 0 {
				args["limit"] = 20
			}

			content, err := s.client.CallTool(c, toolName, args)
			if err != nil {
				continue
			}

			for _, item := range content {
				if item.Type == "text" {
					pages := parsePagesFromText(item.Text)
					external.Pages = append(external.Pages, pages...)
				}
			}

			if len(external.Pages) > 0 {
				return nil
			}
		}
	}

	return nil
}

func (s *Source) fetchDocuments(c context.Context, tools map[string]Tool, external *ctx.ExternalContext) error {
	// Look for document-fetching tools
	docTools := []string{
		"list_documents",
		"search_documents",
		"google_docs_list",
	}

	for _, toolName := range docTools {
		if _, ok := tools[toolName]; ok {
			args := make(map[string]any)
			if s.config.Config != nil {
				if folder, ok := s.config.Config["folder"]; ok {
					args["folder"] = folder
				}
				if query, ok := s.config.Config["query"]; ok {
					args["query"] = query
				}
			}

			content, err := s.client.CallTool(c, toolName, args)
			if err != nil {
				continue
			}

			for _, item := range content {
				if item.Type == "text" {
					docs := parseDocumentsFromText(item.Text)
					external.Documents = append(external.Documents, docs...)
				}
			}

			if len(external.Documents) > 0 {
				return nil
			}
		}
	}

	return nil
}

// parseIssuesFromText extracts issues from tool response text.
// This is a simplified parser - real implementation would handle JSON responses.
func parseIssuesFromText(text string) []ctx.Issue {
	var issues []ctx.Issue

	// Try to parse as JSON array
	var jsonIssues []struct {
		Key         string   `json:"key"`
		Type        string   `json:"type"`
		Summary     string   `json:"summary"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		Priority    string   `json:"priority"`
		Assignee    string   `json:"assignee"`
		Labels      []string `json:"labels"`
	}

	if err := parseJSON(text, &jsonIssues); err == nil {
		for _, ji := range jsonIssues {
			issues = append(issues, ctx.Issue{
				Key:         ji.Key,
				Type:        ji.Type,
				Summary:     ji.Summary,
				Description: truncate(ji.Description, 500),
				Status:      ji.Status,
				Priority:    ji.Priority,
				Assignee:    ji.Assignee,
				Labels:      ji.Labels,
			})
		}
	}

	return issues
}

func parsePagesFromText(text string) []ctx.Page {
	var pages []ctx.Page

	var jsonPages []struct {
		ID      string   `json:"id"`
		Title   string   `json:"title"`
		Space   string   `json:"space"`
		Content string   `json:"content"`
		Labels  []string `json:"labels"`
		URL     string   `json:"url"`
	}

	if err := parseJSON(text, &jsonPages); err == nil {
		for _, jp := range jsonPages {
			pages = append(pages, ctx.Page{
				ID:      jp.ID,
				Title:   jp.Title,
				Space:   jp.Space,
				Content: truncate(jp.Content, 2000),
				Labels:  jp.Labels,
				URL:     jp.URL,
			})
		}
	}

	return pages
}

func parseDocumentsFromText(text string) []ctx.Document {
	var docs []ctx.Document

	var jsonDocs []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Type    string `json:"type"`
		Content string `json:"content"`
		URL     string `json:"url"`
	}

	if err := parseJSON(text, &jsonDocs); err == nil {
		for _, jd := range jsonDocs {
			docs = append(docs, ctx.Document{
				ID:      jd.ID,
				Title:   jd.Title,
				Type:    jd.Type,
				Content: truncate(jd.Content, 2000),
				URL:     jd.URL,
			})
		}
	}

	return docs
}

func parseJSON(text string, v any) error {
	// Try to find JSON in the text
	start := strings.Index(text, "[")
	if start == -1 {
		start = strings.Index(text, "{")
	}
	if start == -1 {
		return fmt.Errorf("no JSON found")
	}

	// Find matching end
	depth := 0
	end := -1
	startChar := text[start]
	endChar := byte(']')
	if startChar == '{' {
		endChar = '}'
	}

	for i := start; i < len(text); i++ {
		if text[i] == startChar {
			depth++
		} else if text[i] == endChar {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}

	if end == -1 {
		return fmt.Errorf("no matching JSON end")
	}

	return json.Unmarshal([]byte(text[start:end]), v)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func generateMCPSummary(external *ctx.ExternalContext) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("MCP Server: %s (%s)", external.ServerName, external.ServerType))

	if len(external.Issues) > 0 {
		parts = append(parts, fmt.Sprintf("Issues: %d", len(external.Issues)))
	}
	if len(external.Pages) > 0 {
		parts = append(parts, fmt.Sprintf("Pages: %d", len(external.Pages)))
	}
	if len(external.Documents) > 0 {
		parts = append(parts, fmt.Sprintf("Documents: %d", len(external.Documents)))
	}

	return strings.Join(parts, ", ")
}
