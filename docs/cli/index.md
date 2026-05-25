# CLI Reference

MultiSpec provides a command-line interface for managing specifications.

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--project` | `-p` | Project name or path |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help |
| `--version` | | Show version |

## Commands

| Command | Description |
|---------|-------------|
| [init](init.md) | Initialize a new project |
| [lint](lint.md) | Validate directory structure |
| [status](status.md) | Show project status |
| [profiles](profiles.md) | Manage configuration profiles |
| [graph](graph.md) | Manage requirement graphs |
| `targets` | List available export targets |
| `eval` | Evaluate specs using LLM judges |
| `synthesize` | Generate specs from source docs |
| `reconcile` | Generate unified execution spec |
| `approve` | Approve a spec for reconciliation |
| `export` | Export to target execution system |
| `serve` | Start MCP server |
| `docs` | Generate MkDocs documentation |
| `context` | Gather codebase context |

## Usage Examples

```bash
# Initialize a project
multispec init user-onboarding

# Lint all projects
multispec lint

# Lint specific project
multispec lint user-onboarding

# Check status
multispec status -p user-onboarding

# JSON output
multispec status -p user-onboarding --format json

# Generate HTML report
multispec status -p user-onboarding --format html > status.html

# CI mode (exit non-zero if not ready)
multispec status -p user-onboarding --ci

# List export targets
multispec targets
```
