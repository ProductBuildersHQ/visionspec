# Context Sources TRD

## Overview

Technical design for context aggregation from multiple sources to ground spec synthesis in reality.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              CLI Layer                                   │
│  multispec context gather|show|snapshot                                  │
│  multispec synthesize --with-context                                     │
│  multispec align --with-context                                          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           pkg/context/                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │ Aggregator  │  │   Cache     │  │  Snapshot   │  │   Config    │    │
│  │             │  │             │  │             │  │   Loader    │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
        ┌───────────────┬───────────┼───────────┬───────────────┐
        ▼               ▼           ▼           ▼               ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────┐
│ pkg/context │ │ pkg/context │ │ pkg/context │ │ pkg/context │ │   MCP   │
│   /git      │ │  /graphize  │ │   /mcp      │ │   /file     │ │ Servers │
│             │ │             │ │             │ │             │ │(external)│
│ - Structure │ │ - Nodes     │ │ - Client    │ │ - Markdown  │ │         │
│ - Deps      │ │ - Edges     │ │ - Jira      │ │ - OpenAPI   │ │         │
│ - APIs      │ │ - Stats     │ │ - Confluence│ │ - Diagrams  │ │         │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ └─────────┘
        │               │               │               │
        ▼               ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        External Resources                                │
│  Git Repos    .graphize/    Jira/Confluence    Local Files              │
└─────────────────────────────────────────────────────────────────────────┘
```

## Package Structure

```
pkg/context/
├── context.go          # Core types: Source, ContextData, AggregatedContext
├── aggregator.go       # Aggregator: combines multiple sources
├── config.go           # Configuration parsing from multispec.yaml
├── cache.go            # Caching layer with TTL
├── snapshot.go         # Snapshot save/load
├── git/
│   ├── source.go       # GitSource implementation
│   ├── structure.go    # Directory tree analysis
│   ├── deps.go         # Dependency extraction (go.mod, package.json)
│   ├── api.go          # API schema detection (OpenAPI, GraphQL, Proto)
│   └── git_test.go
├── graphize/
│   ├── source.go       # GraphizeSource implementation
│   ├── extract.go      # Node/edge extraction from .graphize/
│   └── graphize_test.go
├── mcp/
│   ├── client.go       # MCP client (subprocess management)
│   ├── source.go       # MCPSource implementation
│   ├── jira.go         # Jira-specific extraction
│   ├── confluence.go   # Confluence-specific extraction
│   └── mcp_test.go
├── file/
│   ├── source.go       # FileSource implementation
│   └── file_test.go
└── context_test.go
```

## Core Types

### Source Interface

```go
// pkg/context/context.go

// Source represents a context source
type Source interface {
    // Name returns the source identifier
    Name() string

    // Type returns the source type
    Type() SourceType

    // Fetch retrieves context data from the source
    Fetch(ctx context.Context) (*ContextData, error)

    // String returns a human-readable description
    String() string
}

type SourceType string

const (
    SourceTypeGit      SourceType = "git"
    SourceTypeGraphize SourceType = "graphize"
    SourceTypeMCP      SourceType = "mcp"
    SourceTypeFile     SourceType = "file"
)
```

### Context Data

```go
// ContextData holds data from a single source
type ContextData struct {
    Source    string       `json:"source"`
    Type      SourceType   `json:"type"`
    FetchedAt time.Time    `json:"fetched_at"`
    Duration  Duration     `json:"duration"`

    // Type-specific data (one populated based on Type)
    Code     *CodeContext     `json:"code,omitempty"`
    Graph    *GraphContext    `json:"graph,omitempty"`
    External *ExternalContext `json:"external,omitempty"`
    File     *FileContext     `json:"file,omitempty"`

    // LLM-friendly summary (generated)
    Summary  string `json:"summary"`

    // Errors during fetch (partial results)
    Errors   []string `json:"errors,omitempty"`
}

// CodeContext holds git repository analysis
type CodeContext struct {
    RepoPath     string            `json:"repo_path"`
    Branch       string            `json:"branch"`
    Commit       string            `json:"commit"`
    Structure    *TreeNode         `json:"structure"`
    Dependencies []Dependency      `json:"dependencies"`
    APIs         []APISchema       `json:"apis"`
    Languages    map[string]int    `json:"languages"`  // lang -> LOC
    README       string            `json:"readme,omitempty"`
}

type TreeNode struct {
    Name     string      `json:"name"`
    Type     string      `json:"type"`  // "file" or "dir"
    Children []*TreeNode `json:"children,omitempty"`
}

type Dependency struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Type    string `json:"type"`  // "direct" or "indirect"
    Source  string `json:"source"` // "go.mod", "package.json", etc.
}

type APISchema struct {
    Path   string `json:"path"`
    Format string `json:"format"` // "openapi", "graphql", "proto"
    Title  string `json:"title,omitempty"`
    Routes []APIRoute `json:"routes,omitempty"`
}

type APIRoute struct {
    Method string `json:"method"`
    Path   string `json:"path"`
    Summary string `json:"summary,omitempty"`
}

// GraphContext holds graphize graph analysis
type GraphContext struct {
    GraphPath    string         `json:"graph_path"`
    NodeCount    int            `json:"node_count"`
    EdgeCount    int            `json:"edge_count"`
    Requirements []Requirement  `json:"requirements"`
    Decisions    []Decision     `json:"decisions"`
    Constraints  []Constraint   `json:"constraints"`
    UserStories  []UserStory    `json:"user_stories"`
    Traceability []TraceLink    `json:"traceability"`
    Stats        GraphStats     `json:"stats"`
}

type Requirement struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description,omitempty"`
    Priority    string   `json:"priority,omitempty"`
    Tags        []string `json:"tags,omitempty"`
}

type TraceLink struct {
    FromID   string `json:"from_id"`
    FromType string `json:"from_type"`
    ToID     string `json:"to_id"`
    ToType   string `json:"to_type"`
    Relation string `json:"relation"` // traces_to, derived_from, etc.
}

// ExternalContext holds MCP server data
type ExternalContext struct {
    ServerName string    `json:"server_name"`
    ServerType string    `json:"server_type"` // jira, confluence, google, etc.
    Issues     []Issue   `json:"issues,omitempty"`
    Pages      []Page    `json:"pages,omitempty"`
    Documents  []Document `json:"documents,omitempty"`
}

type Issue struct {
    Key         string   `json:"key"`
    Type        string   `json:"type"`  // epic, story, task, bug
    Summary     string   `json:"summary"`
    Description string   `json:"description,omitempty"`
    Status      string   `json:"status"`
    Priority    string   `json:"priority,omitempty"`
    Labels      []string `json:"labels,omitempty"`
}
```

### Aggregated Context

```go
// AggregatedContext combines data from all sources
type AggregatedContext struct {
    Project   string         `json:"project"`
    GatheredAt time.Time     `json:"gathered_at"`
    Duration  time.Duration  `json:"duration"`
    Sources   []*ContextData `json:"sources"`

    // Merged summaries for LLM
    Summary   string `json:"summary"`

    // Quick access helpers (derived from Sources)
    HasCode     bool `json:"has_code"`
    HasGraph    bool `json:"has_graph"`
    HasExternal bool `json:"has_external"`
}

func (ac *AggregatedContext) CodeContexts() []*CodeContext
func (ac *AggregatedContext) GraphContexts() []*GraphContext
func (ac *AggregatedContext) ExternalContexts() []*ExternalContext
func (ac *AggregatedContext) GenerateSummary() string
```

## Aggregator

```go
// pkg/context/aggregator.go

// Aggregator fetches and combines context from multiple sources
type Aggregator struct {
    sources []Source
    cache   *Cache
    config  *ContextConfig
}

// NewAggregator creates an aggregator from configuration
func NewAggregator(cfg *ContextConfig) (*Aggregator, error) {
    sources := make([]Source, 0)

    // Create git sources
    for _, repo := range cfg.Repositories {
        src, err := git.NewSource(repo)
        if err != nil {
            return nil, fmt.Errorf("git source %s: %w", repo.Path, err)
        }
        sources = append(sources, src)

        // Auto-detect graphize if enabled
        if repo.Graphize == "auto" || repo.Graphize == "true" {
            graphizePath := filepath.Join(repo.Path, ".graphize")
            if _, err := os.Stat(graphizePath); err == nil {
                gSrc := graphize.NewSource(graphizePath, repo.Path)
                sources = append(sources, gSrc)
            }
        }
    }

    // Create explicit graphize sources
    for _, g := range cfg.Graphize {
        src := graphize.NewSource(g.Path, g.Name)
        sources = append(sources, src)
    }

    // Create MCP sources
    for name, server := range cfg.MCPServers {
        src, err := mcp.NewSource(name, server)
        if err != nil {
            return nil, fmt.Errorf("mcp source %s: %w", name, err)
        }
        sources = append(sources, src)
    }

    // Create file sources
    for _, f := range cfg.Files {
        src := file.NewSource(f.Path, f.Type)
        sources = append(sources, src)
    }

    return &Aggregator{
        sources: sources,
        cache:   NewCache(cfg.CacheTTL),
        config:  cfg,
    }, nil
}

// Gather fetches context from all sources
func (a *Aggregator) Gather(ctx context.Context) (*AggregatedContext, error) {
    start := time.Now()
    results := make([]*ContextData, 0, len(a.sources))
    var mu sync.Mutex
    var wg sync.WaitGroup

    // Fetch from all sources concurrently
    for _, src := range a.sources {
        wg.Add(1)
        go func(s Source) {
            defer wg.Done()

            // Check cache first
            if cached := a.cache.Get(s.Name()); cached != nil {
                mu.Lock()
                results = append(results, cached)
                mu.Unlock()
                return
            }

            // Fetch from source
            data, err := s.Fetch(ctx)
            if err != nil {
                data = &ContextData{
                    Source: s.Name(),
                    Type:   s.Type(),
                    Errors: []string{err.Error()},
                }
            }

            // Cache result
            a.cache.Set(s.Name(), data)

            mu.Lock()
            results = append(results, data)
            mu.Unlock()
        }(src)
    }

    wg.Wait()

    ac := &AggregatedContext{
        Project:    a.config.ProjectName,
        GatheredAt: time.Now(),
        Duration:   time.Since(start),
        Sources:    results,
    }

    // Set helper flags
    for _, src := range results {
        switch src.Type {
        case SourceTypeGit:
            ac.HasCode = true
        case SourceTypeGraphize:
            ac.HasGraph = true
        case SourceTypeMCP:
            ac.HasExternal = true
        }
    }

    // Generate combined summary
    ac.Summary = ac.GenerateSummary()

    return ac, nil
}
```

## Git Source Implementation

```go
// pkg/context/git/source.go

type GitSource struct {
    path     string
    config   RepositoryConfig
}

func NewSource(cfg RepositoryConfig) (*GitSource, error) {
    // Validate path exists or clone URL
    if cfg.URL != "" {
        // Clone to temp directory with sparse checkout
        path, err := cloneSpare(cfg.URL, cfg.Branch, cfg.Sparse)
        if err != nil {
            return nil, err
        }
        cfg.Path = path
    }

    return &GitSource{path: cfg.Path, config: cfg}, nil
}

func (g *GitSource) Fetch(ctx context.Context) (*ContextData, error) {
    start := time.Now()

    code := &CodeContext{
        RepoPath: g.path,
    }

    // Get git info
    repo, err := git.PlainOpen(g.path)
    if err == nil {
        head, _ := repo.Head()
        code.Branch = head.Name().Short()
        code.Commit = head.Hash().String()[:7]
    }

    // Analyze based on config
    for _, analysis := range g.config.Analyze {
        switch analysis {
        case "structure":
            code.Structure = g.analyzeStructure()
        case "dependencies":
            code.Dependencies = g.analyzeDependencies()
        case "api_schemas":
            code.APIs = g.analyzeAPIs()
        case "readme":
            code.README = g.readREADME()
        }
    }

    // Count languages
    code.Languages = g.countLanguages()

    return &ContextData{
        Source:    g.Name(),
        Type:      SourceTypeGit,
        FetchedAt: time.Now(),
        Duration:  time.Since(start),
        Code:      code,
        Summary:   g.generateSummary(code),
    }, nil
}
```

## MCP Client Implementation

```go
// pkg/context/mcp/client.go

// Client manages connection to an MCP server subprocess
type Client struct {
    name    string
    cmd     *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
    decoder *json.Decoder
    encoder *json.Encoder
    mu      sync.Mutex
}

func NewClient(name string, cfg MCPServerConfig) (*Client, error) {
    cmd := exec.Command(cfg.Command, cfg.Args...)

    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }

    // Set environment from config
    cmd.Env = os.Environ()
    for k, v := range cfg.Env {
        cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
    }

    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("starting MCP server: %w", err)
    }

    return &Client{
        name:    name,
        cmd:     cmd,
        stdin:   stdin,
        stdout:  stdout,
        decoder: json.NewDecoder(stdout),
        encoder: json.NewEncoder(stdin),
    }, nil
}

func (c *Client) CallTool(ctx context.Context, tool string, args map[string]any) (any, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Send JSON-RPC request
    req := map[string]any{
        "jsonrpc": "2.0",
        "id":      1,
        "method":  "tools/call",
        "params": map[string]any{
            "name":      tool,
            "arguments": args,
        },
    }

    if err := c.encoder.Encode(req); err != nil {
        return nil, fmt.Errorf("sending request: %w", err)
    }

    // Read response
    var resp map[string]any
    if err := c.decoder.Decode(&resp); err != nil {
        return nil, fmt.Errorf("reading response: %w", err)
    }

    if errObj, ok := resp["error"]; ok {
        return nil, fmt.Errorf("MCP error: %v", errObj)
    }

    return resp["result"], nil
}

func (c *Client) Close() error {
    c.stdin.Close()
    return c.cmd.Wait()
}
```

## Configuration

```go
// pkg/context/config.go

type ContextConfig struct {
    ProjectName  string                     `yaml:"-"`
    Repositories []RepositoryConfig         `yaml:"repositories"`
    Graphize     []GraphizeConfig           `yaml:"graphize"`
    MCPServers   map[string]MCPServerConfig `yaml:"mcp_servers"`
    Files        []FileConfig               `yaml:"files"`
    CacheTTL     time.Duration              `yaml:"cache_ttl"`
}

type RepositoryConfig struct {
    Path     string   `yaml:"path"`
    URL      string   `yaml:"url,omitempty"`
    Branch   string   `yaml:"branch,omitempty"`
    Sparse   []string `yaml:"sparse,omitempty"`
    Include  []string `yaml:"include"`
    Exclude  []string `yaml:"exclude"`
    Analyze  []string `yaml:"analyze"`
    Graphize string   `yaml:"graphize"` // "auto", "true", "false"
}

type GraphizeConfig struct {
    Path         string   `yaml:"path"`
    Name         string   `yaml:"name"`
    IncludeNodes []string `yaml:"include_nodes"`
    IncludeEdges []string `yaml:"include_edges"`
}

type MCPServerConfig struct {
    Command string            `yaml:"command"`
    Args    []string          `yaml:"args"`
    Env     map[string]string `yaml:"env"`
    Config  map[string]any    `yaml:"config"`
}

type FileConfig struct {
    Path string `yaml:"path"`
    Type string `yaml:"type"`
}

func LoadContextConfig(project *types.Project) (*ContextConfig, error) {
    // Parse context section from multispec.yaml
}
```

## CLI Commands

```go
// pkg/cli/context.go

func contextCmd(cfg *Config) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "context",
        Short: "Manage context sources",
        Long: `Gather, view, and snapshot context from configured sources.

Sources include git repositories, graphize graphs, MCP servers, and local files.`,
    }

    cmd.AddCommand(
        contextGatherCmd(cfg),
        contextShowCmd(cfg),
        contextSnapshotCmd(cfg),
        contextRefreshCmd(cfg),
    )

    return cmd
}

func contextGatherCmd(cfg *Config) *cobra.Command {
    return &cobra.Command{
        Use:   "gather",
        Short: "Fetch context from all configured sources",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Load project config
            // Create aggregator
            // Gather context
            // Display summary
        },
    }
}
```

## Synthesis Integration

```go
// pkg/synth/synth.go (updated)

func (s *Synthesizer) SynthesizeWithContext(
    ctx context.Context,
    targetType types.SpecType,
    input SynthesisInput,
    context *context.AggregatedContext,
) (*SynthesisResult, error) {
    // Get template
    tmpl, err := templates.Get(targetType)
    if err != nil {
        return nil, err
    }

    // Build prompt with context
    prompt := s.buildPromptWithContext(targetType, input, context, tmpl.Content)

    // Call LLM
    content, err := s.client.Complete(ctx, prompt)
    if err != nil {
        return nil, err
    }

    return &SynthesisResult{
        SpecType: targetType,
        Content:  content,
        Sources:  extractSources(input),
        Context:  context.Summary,
    }, nil
}

func (s *Synthesizer) buildPromptWithContext(
    targetType types.SpecType,
    input SynthesisInput,
    ctx *context.AggregatedContext,
    template string,
) string {
    var sb strings.Builder

    sb.WriteString("You are a technical writer synthesizing specification documents.\n\n")

    // Add spec inputs (existing logic)
    // ...

    // Add context section
    if ctx != nil && ctx.HasCode {
        sb.WriteString("## Codebase Context\n\n")
        sb.WriteString("The following represents the current implementation:\n\n")

        for _, code := range ctx.CodeContexts() {
            sb.WriteString(fmt.Sprintf("### Repository: %s (%s)\n\n", code.RepoPath, code.Commit))

            if code.Structure != nil {
                sb.WriteString("**Structure:**\n```\n")
                sb.WriteString(renderTree(code.Structure, 0))
                sb.WriteString("```\n\n")
            }

            if len(code.APIs) > 0 {
                sb.WriteString("**APIs:**\n")
                for _, api := range code.APIs {
                    sb.WriteString(fmt.Sprintf("- %s (%s)\n", api.Title, api.Format))
                    for _, route := range api.Routes {
                        sb.WriteString(fmt.Sprintf("  - %s %s\n", route.Method, route.Path))
                    }
                }
                sb.WriteString("\n")
            }

            if len(code.Dependencies) > 0 {
                sb.WriteString("**Key Dependencies:**\n")
                for _, dep := range code.Dependencies[:min(10, len(code.Dependencies))] {
                    sb.WriteString(fmt.Sprintf("- %s@%s\n", dep.Name, dep.Version))
                }
                sb.WriteString("\n")
            }
        }
    }

    if ctx != nil && ctx.HasGraph {
        sb.WriteString("## Requirement Traceability\n\n")

        for _, graph := range ctx.GraphContexts() {
            if len(graph.Requirements) > 0 {
                sb.WriteString("**Requirements:**\n")
                for _, req := range graph.Requirements {
                    sb.WriteString(fmt.Sprintf("- %s: %s\n", req.ID, req.Title))
                }
                sb.WriteString("\n")
            }

            if len(graph.Traceability) > 0 {
                sb.WriteString("**Traceability:**\n")
                for _, link := range graph.Traceability {
                    sb.WriteString(fmt.Sprintf("- %s → %s (%s)\n", link.FromID, link.ToID, link.Relation))
                }
                sb.WriteString("\n")
            }
        }
    }

    // Add template and instructions
    sb.WriteString("## Template\n\n")
    sb.WriteString(template)
    sb.WriteString("\n\n## Instructions\n\n")
    sb.WriteString("1. Generate the document following the template structure\n")
    sb.WriteString("2. Ground your output in the codebase context provided\n")
    sb.WriteString("3. Reference actual APIs, dependencies, and structure\n")
    sb.WriteString("4. Maintain traceability to requirements where applicable\n")
    sb.WriteString("5. Output ONLY the completed document\n")

    return sb.String()
}
```

## Testing Strategy

### Unit Tests

- `pkg/context/git/`: Mock filesystem, test structure/dep/API extraction
- `pkg/context/graphize/`: Mock graphfs, test node/edge extraction
- `pkg/context/mcp/`: Mock subprocess, test JSON-RPC protocol
- `pkg/context/aggregator_test.go`: Test concurrent gathering, caching

### Integration Tests

- Real git repos (use testdata/ or clone small public repos)
- Real graphize graphs (use fixtures)
- Mock MCP servers (create test server binary)

### E2E Tests

```bash
# Test context gather
multispec context gather

# Test synthesis with context
multispec synthesize trd --with-context

# Test snapshot round-trip
multispec context snapshot > ctx.json
multispec synthesize trd --context-file=ctx.json
```

## Implementation Phases

### Phase 1: Foundation (v0.4.0-alpha)

- [ ] Core types: Source, ContextData, AggregatedContext
- [ ] Aggregator skeleton
- [ ] Configuration parsing
- [ ] CLI: `context gather`, `context show`

### Phase 2: Git + Graphize (v0.4.0-alpha)

- [ ] Git source: structure analysis
- [ ] Git source: dependency extraction
- [ ] Git source: API schema detection
- [ ] Graphize source: node/edge extraction
- [ ] Auto-detect .graphize/ in repos

### Phase 3: Synthesis Integration (v0.4.0-beta)

- [ ] Update Synthesizer with context support
- [ ] `--with-context` flag on synthesize
- [ ] Context-aware prompt building
- [ ] Test TRD/IRD synthesis with context

### Phase 4: MCP Client (v0.4.0-beta)

- [ ] MCP client: subprocess management
- [ ] MCP client: JSON-RPC protocol
- [ ] Jira source implementation
- [ ] Confluence source implementation

### Phase 5: Polish (v0.4.0-rc)

- [ ] Caching with TTL
- [ ] Snapshot save/load
- [ ] `--context-file` flag
- [ ] Align command with context
- [ ] Documentation

## Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/go-git/go-git/v5` | Git operations | v5.x |
| `github.com/plexusone/graphize` | Graph loading | existing |
| `github.com/plexusone/graphfs` | Graph structures | existing |
| `gopkg.in/yaml.v3` | Config parsing | existing |
| `github.com/bmatcuk/doublestar/v4` | Glob matching | v4.x |

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Large repos slow to analyze | Incremental analysis, caching, depth limits |
| MCP server crashes | Timeout handling, graceful degradation |
| Context too large for LLM | Summarization, prioritization, truncation |
| API schema parsing failures | Graceful fallback, report unparseable files |
