# MCP Server

VisionSpec includes a Model Context Protocol (MCP) server for integration with AI coding assistants.

## Overview

The MCP server provides a standardized interface for AI assistants like Claude Code and Kiro CLI to interact with VisionSpec projects.

## Installation

The MCP server is built as a separate binary:

```bash
# Build
make build

# The binary is at:
bin/visionspec-mcp
```

Or install directly:

```bash
go install github.com/ProductBuildersHQ/visionspec/cmd/mcp-server@latest
```

## Configuration

### Claude Code

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "visionspec": {
      "command": "visionspec-mcp",
      "args": []
    }
  }
}
```

### Kiro CLI

Add to your Kiro steering configuration:

```yaml
mcp:
  servers:
    - name: visionspec
      command: visionspec-mcp
```

## Transport

The MCP server uses stdio transport by default, communicating via stdin/stdout with JSON-RPC 2.0 messages.

## Working Directory

The MCP server operates relative to the current working directory. It looks for:

1. `docs/specs/` directory containing VisionSpec projects
2. Individual project directories with `visionspec.yaml`

## Available Tools

See [MCP Tools](tools.md) for the complete list of available tools.

### Project Tools

| Tool | Description |
|------|-------------|
| `list_projects` | List all projects |
| `get_project_status` | Get project readiness |

### Spec Tools

| Tool | Description |
|------|-------------|
| `get_spec` | Get spec content |
| `get_eval` | Get evaluation results |
| `run_eval` | Run evaluation against rubric |
| `synthesize` | Generate specs from sources |
| `reconcile` | Generate unified execution spec |
| `approve` | Approve a spec |
| `export` | Export to target system |

### Draft Authoring Tools

| Tool | Description |
|------|-------------|
| `start_draft` | Initialize a new draft |
| `get_draft` | Get current draft content |
| `update_draft` | Save draft content |
| `eval_draft` | Evaluate draft against rubric |
| `finalize_draft` | Promote draft to final spec |
| `discard_draft` | Delete a draft |
| `list_drafts` | List all drafts in a project |

## Example Usage

Once configured, you can interact with VisionSpec through your AI assistant:

```
User: What visionspec projects are available?
AI: [calls list_projects tool]
    Found 2 projects: user-onboarding, payment-flow

User: What's the status of user-onboarding?
AI: [calls get_project_status tool]
    Project user-onboarding is NOT READY.
    - Required specs present: PASS
    - Evaluations passing: PASS
    - Approvals obtained: FAIL (pending: prd, trd)
    - Execution spec generated: FAIL
```
