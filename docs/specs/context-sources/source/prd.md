# Context Sources PRD

## Overview

Enable visionspec to aggregate context from multiple sources—git repositories, graphize graphs, MCP servers, and local files—to ground spec synthesis and alignment in reality.

**Marketing Name:** Grounding

## Problem Statement

Current visionspec synthesis operates in a vacuum:

1. **Greenfield assumption**: TRD/IRD synthesis assumes no existing code
2. **Spec-only context**: Synthesis uses only spec documents as input
3. **Manual alignment**: Post-ship alignment requires manual "shipped reality" input
4. **Siloed tools**: Context in Jira, Confluence, Google Docs isn't accessible

For brownfield projects, this creates drift between specs and reality.

## Target Users

| User | Need |
|------|------|
| Engineering Lead | Generate TRD that reflects existing architecture |
| Product Manager | Align PRD with what's actually shipped |
| Tech Writer | Pull context from Jira/Confluence for docs |
| Platform Team | Automate spec generation from codebase |

## User Stories

### US-1: Git Repository Context

As an engineering lead, I want to synthesize a TRD from my existing codebase, so that the technical spec reflects what's actually built.

**Acceptance Criteria:**

- [ ] Configure one or more git repository paths in visionspec.yaml
- [ ] Analyze code structure (directory tree, modules)
- [ ] Extract dependencies (go.mod, package.json, requirements.txt)
- [ ] Detect API schemas (OpenAPI, GraphQL, Proto)
- [ ] Include README and inline documentation
- [ ] Support remote repos via URL with sparse checkout

### US-2: Graphize Graph Context

As a requirements engineer, I want synthesis to use existing graphize traceability graphs, so that generated specs maintain requirement relationships.

**Acceptance Criteria:**

- [ ] Auto-detect .graphize/ directories in configured repos
- [ ] Extract nodes: requirements, decisions, constraints, user stories
- [ ] Extract edges: traces_to, derived_from, depends_on, conflicts_with
- [ ] Include traceability statistics in context
- [ ] Support filtering by node/edge type

### US-3: MCP Server Context

As a product manager, I want to pull context from Jira and Confluence, so that specs reflect current project state.

**Acceptance Criteria:**

- [ ] Configure MCP servers in visionspec.yaml
- [ ] Connect to MCP servers as a client
- [ ] Fetch issues/epics from Jira
- [ ] Fetch pages from Confluence
- [ ] Support Google Docs, Office 365, Aha, Productboard
- [ ] Handle authentication via MCP server config

### US-4: Context-Aware Synthesis

As an engineer, I want `visionspec synthesize trd --with-context` to use aggregated context, so that generated specs are grounded in reality.

**Acceptance Criteria:**

- [ ] `--with-context` flag on synthesize command
- [ ] Aggregator combines all configured sources
- [ ] Synthesis prompt includes context summary
- [ ] Generated TRD references actual code/APIs
- [ ] Generated IRD references actual infrastructure

### US-5: Context-Aware Alignment

As a tech lead, I want `visionspec align --with-context` to compare specs against codebase, so that I can detect drift automatically.

**Acceptance Criteria:**

- [ ] `--with-context` flag on align command
- [ ] Compare spec.md against actual codebase
- [ ] Detect unimplemented requirements
- [ ] Detect undocumented features
- [ ] Generate current-truth.md with findings

### US-6: Context Snapshots

As a DevOps engineer, I want to snapshot context for CI reproducibility, so that synthesis is deterministic.

**Acceptance Criteria:**

- [ ] `visionspec context gather` fetches all sources
- [ ] `visionspec context snapshot` saves to JSON
- [ ] `visionspec synthesize --context-file=...` uses snapshot
- [ ] Snapshots are cacheable and diffable

## Functional Requirements

### FR-1: Context Source Interface

```go
type Source interface {
    Name() string
    Type() SourceType
    Fetch(ctx context.Context) (*ContextData, error)
}

type SourceType string

const (
    SourceTypeGit      SourceType = "git"
    SourceTypeGraphize SourceType = "graphize"
    SourceTypeMCP      SourceType = "mcp"
    SourceTypeFile     SourceType = "file"
)
```

### FR-2: Context Data Schema

```go
type ContextData struct {
    Source    string     `json:"source"`
    Type      SourceType `json:"type"`
    FetchedAt time.Time  `json:"fetched_at"`

    // Code context (git)
    Code      *CodeContext     `json:"code,omitempty"`

    // Graph context (graphize)
    Graph     *GraphContext    `json:"graph,omitempty"`

    // External context (MCP)
    External  *ExternalContext `json:"external,omitempty"`

    // LLM-friendly summary
    Summary   string           `json:"summary"`
}
```

### FR-3: Configuration Schema

```yaml
context:
  repositories:
    - path: ../backend
      include: ["**/*.go"]
      exclude: ["vendor/"]
      analyze: [structure, dependencies, api_schemas]
      graphize: auto

  graphize:
    - path: ./.graphize
      include_nodes: [requirement, decision]

  mcp_servers:
    jira:
      command: "npx"
      args: ["-y", "@anthropic/mcp-jira"]
      config:
        project: "PROJ"

  files:
    - path: ./docs/architecture.md
      type: architecture
```

### FR-4: CLI Commands

| Command | Description |
|---------|-------------|
| `visionspec context gather` | Fetch context from all sources |
| `visionspec context show` | Display aggregated context |
| `visionspec context snapshot` | Save context to JSON file |
| `visionspec context refresh` | Refresh cached context |
| `visionspec synthesize --with-context` | Synthesize using context |
| `visionspec align --with-context` | Align using context |

### FR-5: MCP Client

Multispec must act as an MCP client to connect to external MCP servers:

- Start MCP server subprocess
- Send tool calls to fetch data
- Parse responses into ContextData
- Handle errors and timeouts

## Non-Functional Requirements

### NFR-1: Performance

- Context gathering should complete in < 30 seconds for typical repos
- Support incremental updates (only fetch changed data)
- Cache context with configurable TTL

### NFR-2: Security

- MCP server credentials stored securely (not in visionspec.yaml)
- Support environment variable substitution
- Audit log of external data access

### NFR-3: Extensibility

- Plugin architecture for custom context sources
- Custom analyzers for language-specific code analysis
- Custom MCP server integrations

## Out of Scope (v0.4.0)

- Real-time sync with external tools
- Bidirectional sync (writing back to Jira/Confluence)
- Visual graph editor for context
- Multi-tenant context sharing

## Dependencies

| Dependency | Purpose |
|------------|---------|
| graphize | Graph loading and analysis |
| graphfs | Graph data structures |
| go-git | Git repository operations |
| mcp-go | MCP client implementation |

## Success Metrics

| Metric | Target |
|--------|--------|
| Context gather time | < 30s for 100k LOC repo |
| Synthesis accuracy | TRD matches codebase 90%+ |
| Drift detection | Identifies 95%+ of spec-code drift |
| MCP server support | 5+ integrations at launch |

## Milestones

| Phase | Scope | Target |
|-------|-------|--------|
| v0.4.0-alpha | Git + Graphize context | Week 1-2 |
| v0.4.0-beta | MCP client + synthesis | Week 3-4 |
| v0.4.0-rc | Snapshots + caching | Week 5 |
| v0.4.0 | Release | Week 6 |
