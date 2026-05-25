# Context Sources

Context sources enable visionspec to gather information from your existing codebases, requirement graphs, and external tools to ground spec synthesis in reality.

## Overview

When synthesizing technical specs (TRD, IRD), visionspec can automatically gather context from:

- **Git Repositories** - Code structure, dependencies, API schemas
- **Graphize Graphs** - Requirement traceability data
- **MCP Servers** - External tools like Jira, Confluence, Google Docs
- **Local Files** - Architecture docs, design decisions

This context is included in LLM prompts to ensure generated specs reflect your actual codebase and organizational knowledge.

## Configuration

Add context sources to your `visionspec.yaml`:

```yaml
name: my-feature
context:
  repositories:
    - path: /path/to/backend
      graphize: auto
    - url: https://github.com/org/frontend.git
      branch: main
      sparse:
        - src/components
        - src/api

  graphize:
    - path: /path/to/requirements/.graphize
      name: product-requirements

  files:
    - path: docs/architecture.md
      type: architecture
    - path: docs/adr/
      type: decisions

  mcp_servers:
    jira:
      command: npx
      args: ["-y", "@anthropic/jira-mcp"]
      config:
        project: PROJ
        jql: "status != Done"

    confluence:
      command: npx
      args: ["-y", "@anthropic/confluence-mcp"]
      config:
        space: ENGINEERING

  cache_ttl: 5m
```

## Git Repositories

### Local Repository

```yaml
context:
  repositories:
    - path: /absolute/path/to/repo
      analyze:
        - structure      # Directory tree
        - dependencies   # go.mod, package.json, etc.
        - api_schemas    # OpenAPI, GraphQL, Proto
        - readme         # README content
      max_depth: 5       # Directory tree depth
      include:
        - "**/*.go"
        - "**/*.ts"
      exclude:
        - "**/vendor/**"
        - "**/node_modules/**"
```

### Remote Repository

```yaml
context:
  repositories:
    - url: https://github.com/org/repo.git
      branch: main
      sparse:
        - src/
        - pkg/
```

Remote repositories are cloned to a cache directory (`/tmp/visionspec-repos/`). Shallow clones are used for performance.

### Extracted Information

| Analysis | Data Extracted |
|----------|----------------|
| `structure` | Directory tree, file sizes |
| `dependencies` | Package names, versions (go.mod, package.json, requirements.txt, Cargo.toml) |
| `api_schemas` | OpenAPI/Swagger routes, GraphQL queries/mutations, Proto services |
| `readme` | README content (truncated to 5KB) |

## Graphize Integration

[Graphize](https://github.com/plexusone/graphize) provides requirement traceability graphs.

### Auto-Detection

When `graphize: auto` is set on a repository, visionspec automatically loads graphs from `.graphize/` directories.

```yaml
context:
  repositories:
    - path: /path/to/repo
      graphize: auto  # Detect .graphize/ in repo
```

### Standalone Graph

```yaml
context:
  graphize:
    - path: /path/to/project
      name: requirements
      include_nodes:
        - requirement
        - decision
        - constraint
      include_edges:
        - traces_to
        - depends_on
```

### Extracted Information

| Data | Description |
|------|-------------|
| Nodes | Requirements, decisions, constraints, user stories |
| Edges | Traceability links between nodes |
| Coverage | Percentage of requirements linked to code/tests |

## MCP Servers

MCP (Model Context Protocol) servers provide access to external tools.

### Configuration

```yaml
context:
  mcp_servers:
    server-name:
      command: npx       # Command to start server
      args: ["-y", "@anthropic/jira-mcp"]
      env:               # Environment variables
        JIRA_URL: https://company.atlassian.net
        JIRA_TOKEN: ${JIRA_TOKEN}
      config:            # Server-specific config
        project: PROJ
        jql: "status != Done ORDER BY updated DESC"
      timeout: 30s
```

### Supported Servers

| Server | Package | Data |
|--------|---------|------|
| Jira | `@anthropic/jira-mcp` | Issues, epics, stories |
| Confluence | `@anthropic/confluence-mcp` | Pages, spaces |
| Google Docs | `@anthropic/gdocs-mcp` | Documents |
| Linear | `@anthropic/linear-mcp` | Issues, projects |
| Notion | `@anthropic/notion-mcp` | Pages, databases |

### Environment Variables

Use `${VAR}` syntax to reference environment variables:

```yaml
mcp_servers:
  jira:
    env:
      JIRA_TOKEN: ${JIRA_API_TOKEN}
```

## Local Files

Include individual files as context:

```yaml
context:
  files:
    - path: docs/architecture.md
      type: architecture

    - path: docs/api-spec.yaml
      type: api_spec

    - path: DECISIONS.md
      type: decisions

    - path: diagrams/system.puml
      type: diagram
```

### File Types

| Type | Description |
|------|-------------|
| `architecture` | System architecture documentation |
| `api_spec` | API specifications (OpenAPI, GraphQL) |
| `readme` | README files |
| `decisions` | ADRs, decision logs |
| `diagram` | PlantUML, Mermaid diagrams |
| `requirements` | Requirements documents |
| `config` | Configuration files |
| `document` | Generic documents |

## CLI Commands

### Gather Context

```bash
# Gather from all configured sources
visionspec context gather

# Gather with JSON output
visionspec context gather --format json

# Force refresh (ignore cache)
visionspec context gather --refresh
```

### Show Context

```bash
# Display context summary
visionspec context show
```

### Save Snapshot

```bash
# Save context to file for CI/reproducibility
visionspec context save --output context-snapshot.json
```

### List Sources

```bash
# Show configured sources
visionspec context sources
```

## Context in Synthesis

Context is automatically used when synthesizing TRD and IRD:

```bash
# TRD synthesis includes context
visionspec synthesize trd

# Skip context gathering
visionspec synthesize trd --no-context
```

The context summary is included in the LLM prompt to ground technical decisions.

## Caching

Context is cached in memory with configurable TTL:

```yaml
context:
  cache_ttl: 5m  # Default: 5 minutes
```

Use `--refresh` to bypass cache:

```bash
visionspec context gather --refresh
```

## Snapshots

Save context snapshots for reproducibility:

```bash
# Save snapshot
visionspec context save -o context.json

# Use snapshot in CI
visionspec synthesize trd --context-file context.json
```

Snapshots contain all gathered context and can be version controlled.

## Writing Custom Sources

Implement the `Source` interface:

```go
package custom

import (
    "context"
    ctx "github.com/ProductBuildersHQ/visionspec/pkg/context"
)

type MySource struct {
    // ...
}

func (s *MySource) Name() string {
    return "custom:my-source"
}

func (s *MySource) Type() ctx.SourceType {
    return ctx.SourceTypeFile // or custom type
}

func (s *MySource) Fetch(c context.Context) (*ctx.ContextData, error) {
    // Fetch and return context data
    return &ctx.ContextData{
        Source: s.Name(),
        Type:   s.Type(),
        // ...
    }, nil
}

func (s *MySource) String() string {
    return "My Custom Source"
}
```

Register with the aggregator:

```go
agg := ctx.NewAggregator("project", cfg)
agg.AddSource(&MySource{})
```

## Best Practices

1. **Use sparse checkout** for large repositories to reduce clone time
2. **Set appropriate TTL** - longer for stable codebases, shorter for active development
3. **Use snapshots in CI** for reproducible builds
4. **Filter nodes/edges** in graphize to reduce context size
5. **Truncate large files** - context is automatically truncated for LLM token limits
