// Package context provides context aggregation from multiple sources
// for grounding spec synthesis in reality.
//
// Sources include git repositories, graphize graphs, MCP servers, and local files.
// The aggregated context is used by synthesis and alignment commands to generate
// specs that reflect actual codebases and external tool state.
package context

import (
	"context"
	"time"
)

// SourceType identifies the type of context source.
type SourceType string

const (
	// SourceTypeGit represents a git repository source.
	SourceTypeGit SourceType = "git"

	// SourceTypeGraphize represents a graphize graph source.
	SourceTypeGraphize SourceType = "graphize"

	// SourceTypeMCP represents an MCP server source.
	SourceTypeMCP SourceType = "mcp"

	// SourceTypeFile represents a local file source.
	SourceTypeFile SourceType = "file"
)

// Source represents a context source that can be fetched.
type Source interface {
	// Name returns a unique identifier for this source.
	Name() string

	// Type returns the source type.
	Type() SourceType

	// Fetch retrieves context data from the source.
	Fetch(ctx context.Context) (*ContextData, error)

	// String returns a human-readable description.
	String() string
}

// ContextData holds data from a single source.
type ContextData struct {
	// Source identifier
	Source string `json:"source"`

	// Type of source
	Type SourceType `json:"type"`

	// When the data was fetched
	FetchedAt time.Time `json:"fetched_at"`

	// How long the fetch took
	Duration time.Duration `json:"duration"`

	// Type-specific data (one populated based on Type)
	Code     *CodeContext     `json:"code,omitempty"`
	Graph    *GraphContext    `json:"graph,omitempty"`
	External *ExternalContext `json:"external,omitempty"`
	File     *FileContext     `json:"file,omitempty"`

	// LLM-friendly summary (generated from the data)
	Summary string `json:"summary"`

	// Errors encountered during fetch (partial results may still be present)
	Errors []string `json:"errors,omitempty"`
}

// CodeContext holds git repository analysis results.
type CodeContext struct {
	// Repository path (local or URL)
	RepoPath string `json:"repo_path"`

	// Current branch name
	Branch string `json:"branch,omitempty"`

	// Current commit hash (short)
	Commit string `json:"commit,omitempty"`

	// Directory structure tree
	Structure *TreeNode `json:"structure,omitempty"`

	// Extracted dependencies
	Dependencies []Dependency `json:"dependencies,omitempty"`

	// Detected API schemas
	APIs []APISchema `json:"apis,omitempty"`

	// Lines of code by language
	Languages map[string]int `json:"languages,omitempty"`

	// README content (truncated if large)
	README string `json:"readme,omitempty"`
}

// TreeNode represents a file or directory in the code structure.
type TreeNode struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"` // "file" or "dir"
	Size     int64       `json:"size,omitempty"`
	Children []*TreeNode `json:"children,omitempty"`
}

// Dependency represents a project dependency.
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`   // "direct" or "indirect"
	Source  string `json:"source"` // "go.mod", "package.json", etc.
}

// APISchema represents a detected API schema file.
type APISchema struct {
	Path   string     `json:"path"`
	Format string     `json:"format"` // "openapi", "graphql", "proto"
	Title  string     `json:"title,omitempty"`
	Routes []APIRoute `json:"routes,omitempty"`
}

// APIRoute represents an API endpoint.
type APIRoute struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Summary     string `json:"summary,omitempty"`
	OperationID string `json:"operation_id,omitempty"`
}

// GraphContext holds graphize graph analysis results.
type GraphContext struct {
	// Path to the .graphize directory
	GraphPath string `json:"graph_path"`

	// Raw graph data
	Nodes []GraphNode `json:"nodes,omitempty"`
	Edges []GraphEdge `json:"edges,omitempty"`

	// Node counts by type
	NodeCount        int `json:"node_count"`
	EdgeCount        int `json:"edge_count"`
	RequirementCount int `json:"requirement_count"`
	CodeCount        int `json:"code_count"`
	TestCount        int `json:"test_count"`

	// Coverage metrics
	LinkedRequirements int     `json:"linked_requirements"`
	TestedRequirements int     `json:"tested_requirements"`
	CodeCoverage       float64 `json:"code_coverage"`
	TestCoverage       float64 `json:"test_coverage"`

	// Metadata
	Version string `json:"version,omitempty"`
	Tool    string `json:"tool,omitempty"`

	// Extracted typed nodes (optional, for higher-level analysis)
	Requirements []Requirement `json:"requirements,omitempty"`
	Decisions    []Decision    `json:"decisions,omitempty"`
	Constraints  []Constraint  `json:"constraints,omitempty"`
	UserStories  []UserStory   `json:"user_stories,omitempty"`

	// Extracted edges
	Traceability []TraceLink `json:"traceability,omitempty"`

	// Graph statistics
	Stats GraphStats `json:"stats"`
}

// GraphNode represents a generic node in the graphize graph.
type GraphNode struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Label    string            `json:"label"`
	Path     string            `json:"path,omitempty"`
	Line     int               `json:"line,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// GraphEdge represents an edge in the graphize graph.
type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// Requirement represents a requirement node from graphize.
type Requirement struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Status      string   `json:"status,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	SourceSpec  string   `json:"source_spec,omitempty"` // e.g., "prd"
}

// Decision represents an architectural decision node.
type Decision struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"` // "proposed", "accepted", "deprecated"
	Date        string   `json:"date,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Constraint represents a constraint node.
type Constraint struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"` // "technical", "business", "regulatory"
	SourceSpec  string `json:"source_spec,omitempty"`
}

// UserStory represents a user story node.
type UserStory struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	AsA            string   `json:"as_a,omitempty"`
	IWant          string   `json:"i_want,omitempty"`
	SoThat         string   `json:"so_that,omitempty"`
	AcceptanceCrit []string `json:"acceptance_criteria,omitempty"`
}

// TraceLink represents a traceability edge between nodes.
type TraceLink struct {
	FromID   string `json:"from_id"`
	FromType string `json:"from_type"`
	ToID     string `json:"to_id"`
	ToType   string `json:"to_type"`
	Relation string `json:"relation"` // traces_to, derived_from, depends_on, conflicts_with
}

// GraphStats holds statistics about the graph.
type GraphStats struct {
	RequirementCount int            `json:"requirement_count"`
	DecisionCount    int            `json:"decision_count"`
	ConstraintCount  int            `json:"constraint_count"`
	UserStoryCount   int            `json:"user_story_count"`
	TraceabilityPct  float64        `json:"traceability_pct"` // % of requirements with traces
	NodesByType      map[string]int `json:"nodes_by_type"`
	EdgesByType      map[string]int `json:"edges_by_type"`
}

// ExternalContext holds data from MCP servers (Jira, Confluence, etc.).
type ExternalContext struct {
	// Server identifier
	ServerName string `json:"server_name"`

	// Server type (jira, confluence, google, etc.)
	ServerType string `json:"server_type"`

	// Issue tracker data
	Issues []Issue `json:"issues,omitempty"`
	Epics  []Epic  `json:"epics,omitempty"`

	// Documentation data
	Pages     []Page     `json:"pages,omitempty"`
	Documents []Document `json:"documents,omitempty"`
}

// Issue represents an issue from Jira or similar.
type Issue struct {
	Key         string   `json:"key"`
	Type        string   `json:"type"` // epic, story, task, bug
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Created     string   `json:"created,omitempty"`
	Updated     string   `json:"updated,omitempty"`
}

// Epic represents an epic with child issues.
type Epic struct {
	Issue
	ChildKeys []string `json:"child_keys,omitempty"`
}

// Page represents a wiki/documentation page.
type Page struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Space   string   `json:"space,omitempty"`
	Content string   `json:"content,omitempty"` // Truncated if large
	Labels  []string `json:"labels,omitempty"`
	Version int      `json:"version,omitempty"`
	Updated string   `json:"updated,omitempty"`
	URL     string   `json:"url,omitempty"`
}

// Document represents a document from Google Docs, Office 365, etc.
type Document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type,omitempty"` // doc, sheet, slide
	Content string `json:"content,omitempty"`
	Updated string `json:"updated,omitempty"`
	URL     string `json:"url,omitempty"`
}

// FileContext holds data from local files.
type FileContext struct {
	Path    string `json:"path"`
	Type    string `json:"type"` // architecture, api_spec, readme, diagram
	Content string `json:"content"`
	Format  string `json:"format,omitempty"` // markdown, yaml, json
}

// AggregatedContext combines data from all sources.
type AggregatedContext struct {
	// Project name
	Project string `json:"project"`

	// When gathering completed
	GatheredAt time.Time `json:"gathered_at"`

	// Total time to gather all sources
	Duration time.Duration `json:"duration"`

	// Individual source data
	Sources []*ContextData `json:"sources"`

	// Combined summary for LLM consumption
	Summary string `json:"summary"`

	// Quick access flags
	HasCode     bool `json:"has_code"`
	HasGraph    bool `json:"has_graph"`
	HasExternal bool `json:"has_external"`
	HasFiles    bool `json:"has_files"`

	// Error count across all sources
	ErrorCount int `json:"error_count"`
}

// CodeContexts returns all code contexts from sources.
func (ac *AggregatedContext) CodeContexts() []*CodeContext {
	var result []*CodeContext
	for _, src := range ac.Sources {
		if src.Code != nil {
			result = append(result, src.Code)
		}
	}
	return result
}

// GraphContexts returns all graph contexts from sources.
func (ac *AggregatedContext) GraphContexts() []*GraphContext {
	var result []*GraphContext
	for _, src := range ac.Sources {
		if src.Graph != nil {
			result = append(result, src.Graph)
		}
	}
	return result
}

// ExternalContexts returns all external contexts from sources.
func (ac *AggregatedContext) ExternalContexts() []*ExternalContext {
	var result []*ExternalContext
	for _, src := range ac.Sources {
		if src.External != nil {
			result = append(result, src.External)
		}
	}
	return result
}

// FileContexts returns all file contexts from sources.
func (ac *AggregatedContext) FileContexts() []*FileContext {
	var result []*FileContext
	for _, src := range ac.Sources {
		if src.File != nil {
			result = append(result, src.File)
		}
	}
	return result
}

// AllRequirements returns requirements from all graph contexts.
func (ac *AggregatedContext) AllRequirements() []Requirement {
	var result []Requirement
	for _, g := range ac.GraphContexts() {
		result = append(result, g.Requirements...)
	}
	return result
}

// AllDecisions returns decisions from all graph contexts.
func (ac *AggregatedContext) AllDecisions() []Decision {
	var result []Decision
	for _, g := range ac.GraphContexts() {
		result = append(result, g.Decisions...)
	}
	return result
}

// AllAPIs returns API schemas from all code contexts.
func (ac *AggregatedContext) AllAPIs() []APISchema {
	var result []APISchema
	for _, c := range ac.CodeContexts() {
		result = append(result, c.APIs...)
	}
	return result
}

// AllDependencies returns dependencies from all code contexts.
func (ac *AggregatedContext) AllDependencies() []Dependency {
	var result []Dependency
	for _, c := range ac.CodeContexts() {
		result = append(result, c.Dependencies...)
	}
	return result
}
