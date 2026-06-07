# sync

Synchronize execution state with target systems.

## Synopsis

```bash
visionspec sync <target> [flags]
```

## Description

The `sync` command performs bidirectional synchronization between VisionSpec projects and execution targets. It reads the current state from the target system and updates the project's execution state.

## Arguments

| Argument | Description |
|----------|-------------|
| `target` | Target system to sync: speckit, gsd, gastown |

## Flags

| Flag | Description |
|------|-------------|
| `-p, --project` | Project name |
| `--dry-run` | Show what would be synced without making changes |
| `--format` | Output format: text, json (default: text) |

## Supported Targets

### SpecKit

Syncs with `.specify/tasks.md` file.

```bash
visionspec sync speckit -p myproject
```

Parses task checkboxes:

- `[ ]` - todo
- `[x]` - done

### GSD

Syncs with `STATE.md` file.

```bash
visionspec sync gsd -p myproject
```

Parses task markers:

- `[ ]` - todo
- `[~]` - in_progress
- `[x]` - done

### GasTown

Syncs with `beads/*.toml` files.

```bash
visionspec sync gastown -p myproject
```

Parses TOML status field:

- `status = "ready"` - todo
- `status = "blocked"` - blocked
- `status = "done"` - done

## Examples

### Basic Sync

```bash
# Sync with SpecKit
visionspec sync speckit -p user-onboarding

# Sync with GSD
visionspec sync gsd -p api-refactor

# Sync with GasTown
visionspec sync gastown -p multi-agent-workflow
```

### Dry Run

```bash
# Preview sync without changes
visionspec sync speckit -p myproject --dry-run
```

### JSON Output

```bash
# Get sync results as JSON
visionspec sync gsd -p myproject --format json
```

## Execution State

After sync, the project's `visionspec.yaml` includes execution state:

```yaml
execution:
  target: speckit
  synced_at: 2024-01-15T10:30:00Z
  tasks:
    - id: TASK-001
      title: Implement user login
      status: done
    - id: TASK-002
      title: Add password reset
      status: in_progress
  summary:
    total_tasks: 10
    todo_count: 3
    in_progress: 2
    done_count: 5
```

## Status Integration

Use with `status` command to see execution progress:

```bash
# Sync then check status
visionspec sync speckit -p myproject
visionspec status -p myproject
```

## See Also

- [export](export.md) - Export to target systems
- [status](status.md) - Check project status
- [drift](drift.md) - Detect spec-to-code drift
