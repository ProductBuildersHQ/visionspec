package context

import "time"

// Config holds the context configuration from visionspec.yaml.
type Config struct {
	// Project name (from parent config)
	ProjectName string `yaml:"-"`

	// Git repositories to analyze
	Repositories []RepositoryConfig `yaml:"repositories,omitempty"`

	// Standalone graphize graphs
	Graphize []GraphizeConfig `yaml:"graphize,omitempty"`

	// MCP servers for external context
	MCPServers map[string]MCPServerConfig `yaml:"mcp_servers,omitempty"`

	// Local files to include
	Files []FileConfig `yaml:"files,omitempty"`

	// Cache TTL for context data (default: 5m)
	CacheTTL time.Duration `yaml:"cache_ttl,omitempty"`
}

// RepositoryConfig configures a git repository source.
type RepositoryConfig struct {
	// Path to local repository
	Path string `yaml:"path,omitempty"`

	// URL for remote repository (clone if not exists)
	URL string `yaml:"url,omitempty"`

	// Branch to checkout (default: default branch)
	Branch string `yaml:"branch,omitempty"`

	// Sparse checkout paths (for large repos)
	Sparse []string `yaml:"sparse,omitempty"`

	// Include patterns for file analysis
	Include []string `yaml:"include,omitempty"`

	// Exclude patterns
	Exclude []string `yaml:"exclude,omitempty"`

	// What to analyze: structure, dependencies, api_schemas, readme
	Analyze []string `yaml:"analyze,omitempty"`

	// Graphize detection: "auto", "true", "false"
	Graphize string `yaml:"graphize,omitempty"`

	// Maximum depth for directory tree (default: 5)
	MaxDepth int `yaml:"max_depth,omitempty"`
}

// GraphizeConfig configures a standalone graphize source.
type GraphizeConfig struct {
	// Path to .graphize directory
	Path string `yaml:"path"`

	// Human-readable name
	Name string `yaml:"name,omitempty"`

	// Node types to include (empty = all)
	IncludeNodes []string `yaml:"include_nodes,omitempty"`

	// Edge types to include (empty = all)
	IncludeEdges []string `yaml:"include_edges,omitempty"`
}

// MCPServerConfig configures an MCP server source.
type MCPServerConfig struct {
	// Command to start the server
	Command string `yaml:"command"`

	// Arguments to the command
	Args []string `yaml:"args,omitempty"`

	// Environment variables
	Env map[string]string `yaml:"env,omitempty"`

	// Server-specific configuration
	Config map[string]any `yaml:"config,omitempty"`

	// Timeout for server operations (default: 30s)
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

// FileConfig configures a local file source.
type FileConfig struct {
	// Path to the file
	Path string `yaml:"path"`

	// Type of content: architecture, api_spec, readme, diagram
	Type string `yaml:"type,omitempty"`

	// Maximum content size to include (default: 50KB)
	MaxSize int64 `yaml:"max_size,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		CacheTTL: 5 * time.Minute,
	}
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	// TODO: Validate paths exist, URLs are valid, etc.
	return nil
}

// HasSources returns true if any context sources are configured.
func (c *Config) HasSources() bool {
	return len(c.Repositories) > 0 ||
		len(c.Graphize) > 0 ||
		len(c.MCPServers) > 0 ||
		len(c.Files) > 0
}

// SourceCount returns the total number of configured sources.
func (c *Config) SourceCount() int {
	count := len(c.Repositories) + len(c.Graphize) + len(c.MCPServers) + len(c.Files)

	// Count auto-detected graphize sources
	for _, repo := range c.Repositories {
		if repo.Graphize == "auto" || repo.Graphize == "true" {
			count++ // Potential graphize source
		}
	}

	return count
}
