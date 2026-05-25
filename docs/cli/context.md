# context

Gather and manage codebase context for grounding spec synthesis.

## Usage

```bash
multispec context <subcommand> [flags]
```

## Description

The `context` command manages codebase context that grounds technical spec synthesis in reality. Context is gathered from git repositories, graphize requirement graphs, local files, and MCP servers.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `gather` | Collect context from all configured sources |
| `show` | Display current context summary |
| `save` | Save context snapshot to file |
| `sources` | List configured context sources |

## context gather

Gather context from all configured sources.

```bash
multispec context gather [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--timeout` | duration | `2m` | Timeout for gathering |
| `--format` | string | `text` | Output format: `text`, `json` |
| `--refresh` | bool | `false` | Refresh cache before gathering |

**Example:**

```bash
# Gather context
multispec context gather

# Output as JSON
multispec context gather --format json

# Force refresh cached data
multispec context gather --refresh

# Custom timeout
multispec context gather --timeout 5m
```

## context show

Display the current context summary from the last gather or saved snapshot.

```bash
multispec context show
```

If no snapshot exists, prompts to run `gather` first.

## context save

Save gathered context to a snapshot file for CI reproducibility.

```bash
multispec context save [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-o`, `--output` | string | `.context-snapshot.json` | Output file path |

**Example:**

```bash
# Save to default location
multispec context save

# Save to custom path
multispec context save -o artifacts/context.json
```

## context sources

List all configured context sources.

```bash
multispec context sources
```

**Output:**

```
Configured context sources:

Git Repositories:
  - /path/to/repo
    (graphize: auto)

Graphize Graphs:
  - spec-graph (/path/to/.graphize)

Local Files:
  - architecture.md (architecture)
  - docs/adr/*.md (adr)

MCP Servers:
  - jira (not yet implemented)
```

## Configuration

Configure context sources in `multispec.yaml`:

```yaml
context:
  # Git repositories to analyze
  repositories:
    - path: "."
      include_structure: true
      include_deps: true
      include_apis: true
      graphize: auto
    - url: "https://github.com/org/repo"
      branch: main
      shallow: true

  # Graphize requirement graphs
  graphize:
    - path: ".graphize"
      name: "spec-graph"
      include_nodes: true
      include_edges: true

  # Local documentation files
  files:
    - path: "docs/architecture.md"
      type: architecture
    - path: "docs/adr/*.md"
      type: adr
    - path: "api/openapi.yaml"
      type: api_spec

  # MCP servers for external tools
  mcp_servers:
    jira:
      command: "npx"
      args: ["-y", "@anthropic/mcp-jira"]
      tools: ["get_issue", "search_issues"]

  # Cache settings
  cache_ttl: 30m
```

## Repository Analysis

When analyzing git repositories, the context includes:

- **Structure** - Directory tree and file organization
- **Dependencies** - Package dependencies (go.mod, package.json, etc.)
- **APIs** - OpenAPI specs, GraphQL schemas, protobuf definitions
- **README** - Repository documentation

## Auto-Detection

If no repositories are configured, multispec auto-detects the repository at the project root.

## Snapshots

Context snapshots enable reproducible synthesis in CI:

```bash
# In CI pipeline
multispec context save -o context-snapshot.json

# Later, use saved context
multispec synthesize trd --context-file context-snapshot.json
```

## See Also

- [Context Sources Guide](../context-sources.md) - Comprehensive configuration guide
- [synthesize](synthesize.md) - Uses context for TRD/IRD synthesis
