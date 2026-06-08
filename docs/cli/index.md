# CLI Reference

VisionSpec provides a command-line interface for managing specifications.

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--project` | `-p` | Project name or path |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help |
| `--version` | | Show version |

## Commands

### Project Setup

| Command | Description |
|---------|-------------|
| [init](init.md) | Initialize a new project |
| [create](create.md) | Create specs from templates |
| [lint](lint.md) | Validate directory structure |
| [status](status.md) | Show project status |
| [profiles](profiles.md) | Manage configuration profiles |

### Spec Workflow

| Command | Description |
|---------|-------------|
| [eval](eval.md) | Evaluate specs using LLM judges |
| [render](render.md) | Render evaluation files to markdown |
| [synthesize](synthesize.md) | Generate GTM/technical specs from sources |
| [reconcile](reconcile.md) | Generate unified execution spec |
| [approve](approve.md) | Approve a spec for reconciliation |
| [watch](watch.md) | Watch spec files and auto-run eval |
| [version](version.md) | Manage spec version history |

### Export & Integration

| Command | Description |
|---------|-------------|
| [export](export.md) | Export to target execution system |
| [targets](targets.md) | List available export targets |
| [serve](serve.md) | Start MCP server |
| [docs](docs.md) | Generate MkDocs documentation |
| [rules](rules.md) | Manage workflow rules for AI assistants |

### Execution Integration

| Command | Description |
|---------|-------------|
| [generate](generate.md) | Generate test stubs from TPD |
| [sync](sync.md) | Sync status with execution targets |
| [drift](drift.md) | Detect spec-to-code drift |

### Context & Traceability

| Command | Description |
|---------|-------------|
| [context](context.md) | Gather codebase context |
| [graph](graph.md) | Manage requirement graphs |

## Usage Examples

```bash
# Initialize a project
visionspec init user-onboarding

# Lint all projects
visionspec lint

# Lint specific project
visionspec lint user-onboarding

# Check status
visionspec status -p user-onboarding

# JSON output
visionspec status -p user-onboarding --format json

# Generate HTML report
visionspec status -p user-onboarding --format html > status.html

# CI mode (exit non-zero if not ready)
visionspec status -p user-onboarding --ci

# List export targets
visionspec targets
```
