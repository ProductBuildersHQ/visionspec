# watch

Watch spec files and auto-run evaluation on changes.

## Usage

```bash
visionspec watch [project] [flags]
```

## Description

Monitor spec files for changes and automatically run evaluation when spec files are modified. Useful during spec authoring for continuous feedback.

## Features

- Watches all `.md` files in the project's source directory
- Debounces rapid changes (500ms default)
- Runs lint and eval on changes
- Shows real-time status updates

## Arguments

| Argument | Description |
|----------|-------------|
| `project` | Project name (optional, defaults to current project) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--debounce` | `500ms` | Debounce interval for file changes |

## Examples

```bash
# Watch current project
visionspec watch

# Watch specific project
visionspec watch user-onboarding

# Custom debounce interval
visionspec watch --debounce 1s

# Watch with verbose output
visionspec watch -v
```

## Output

The watch command displays real-time updates:

```
Watching specs/user-onboarding/source for changes...

[14:32:05] Detected change: prd.md
[14:32:05] Running lint...
[14:32:05] Lint passed
[14:32:06] Running eval prd...
[14:32:08] PRD Score: 8.2/10 (PASS)
```

## See Also

- [eval](eval.md) - Run evaluations
- [lint](lint.md) - Validate directory structure
- [status](status.md) - Show project status
