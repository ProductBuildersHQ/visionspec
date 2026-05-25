# serve

Start the MCP server for AI assistant integration.

## Usage

```bash
visionspec serve [flags]
```

## Description

The `serve` command starts a Model Context Protocol (MCP) server that enables AI coding assistants like Claude Code and Kiro CLI to interact with VisionSpec projects programmatically.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--port` | int | `0` | HTTP port (0 for stdio transport) |
| `--transport` | string | `stdio` | Transport mode: `stdio`, `http`, `sse` |

## Transport Modes

| Mode | Use Case |
|------|----------|
| `stdio` | Default. Communication via stdin/stdout. Used by Claude Code. |
| `http` | HTTP server on specified port. For custom integrations. |
| `sse` | Server-Sent Events. For browser-based clients. |

## MCP Tools

The server exposes these tools to AI assistants:

### Project Management

| Tool | Description |
|------|-------------|
| `list_projects` | List all VisionSpec projects |
| `get_project_status` | Get project readiness status |

### Spec Operations

| Tool | Description |
|------|-------------|
| `get_spec` | Get specification content |
| `get_eval` | Get evaluation results |
| `run_eval` | Run evaluation against rubric |
| `synthesize` | Generate specs from sources |
| `reconcile` | Generate unified spec.md |
| `approve` | Approve a specification |
| `export` | Export to target system |

### Draft Authoring

| Tool | Description |
|------|-------------|
| `start_draft` | Initialize a new draft |
| `get_draft` | Get current draft content |
| `update_draft` | Save draft content |
| `eval_draft` | Evaluate draft against rubric |
| `finalize_draft` | Promote draft to final spec |
| `discard_draft` | Delete a draft |
| `list_drafts` | List all drafts in project |

## Claude Code Configuration

Add to `~/.claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "visionspec": {
      "command": "visionspec",
      "args": ["serve"]
    }
  }
}
```

Or if using the separate MCP binary:

```json
{
  "mcpServers": {
    "visionspec": {
      "command": "visionspec-mcp"
    }
  }
}
```

## Kiro CLI Configuration

Add to your Kiro steering configuration:

```yaml
mcp:
  servers:
    - name: visionspec
      command: visionspec
      args: ["serve"]
```

## Examples

```bash
# Start with stdio transport (for Claude Code)
visionspec serve

# Start HTTP server on port 8080
visionspec serve --port 8080 --transport http

# Start SSE server
visionspec serve --port 3000 --transport sse
```

## Working Directory

The MCP server operates relative to the current working directory. It discovers projects by looking for:

1. `docs/specs/` directory containing VisionSpec projects
2. Individual project directories with `visionspec.yaml`

## See Also

- [MCP Server Overview](../mcp/index.md) - Detailed MCP documentation
- [MCP Tools Reference](../mcp/tools.md) - Complete tool documentation
