# hooks

Manage Git hooks for automatic spec validation.

## Synopsis

```bash
visionspec hooks <subcommand> [flags]
```

## Description

The `hooks` command manages Git hooks that automatically validate specifications during git operations. This ensures spec quality is maintained throughout the development workflow.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `install` | Install visionspec Git hooks |
| `uninstall` | Remove visionspec Git hooks |
| `status` | Show installed hooks status |

## Supported Hooks

| Hook | Trigger | Action |
|------|---------|--------|
| `pre-commit` | Before commit | Lint changed spec files |
| `pre-push` | Before push | Evaluate specs, check for blockers |

## Examples

### Install Hooks

```bash
# Install all hooks
visionspec hooks install

# Install specific hooks
visionspec hooks install --hooks pre-commit,pre-push
```

### Uninstall Hooks

```bash
# Remove all visionspec hooks
visionspec hooks uninstall

# Remove specific hooks
visionspec hooks uninstall --hooks pre-commit
```

### Check Status

```bash
# Show status of all hooks
visionspec hooks status
```

## Hook Behavior

### pre-commit

Runs before each commit to validate changed spec files:

```bash
# Automatically runs on git commit
git commit -m "Update PRD"

# Hook output:
# [visionspec] Linting changed specs...
# [visionspec] ✓ docs/specs/myproject/source/prd.md
# [visionspec] All specs valid
```

If validation fails:

```bash
# [visionspec] Linting changed specs...
# [visionspec] ✗ docs/specs/myproject/source/prd.md
#   - Missing required section: User Stories
# [visionspec] Commit blocked: fix spec errors first
```

### pre-push

Runs before push to evaluate specs and check for blockers:

```bash
# Automatically runs on git push
git push origin main

# Hook output:
# [visionspec] Evaluating specs before push...
# [visionspec] ✓ PRD score: 8.2 (passing)
# [visionspec] ✓ TRD score: 7.8 (passing)
# [visionspec] Push allowed
```

If evaluation fails:

```bash
# [visionspec] Evaluating specs before push...
# [visionspec] ✗ PRD score: 5.5 (failing)
#   - Critical: Missing acceptance criteria
# [visionspec] Push blocked: resolve critical findings
```

## Status Output

```bash
$ visionspec hooks status

Git Hooks Status
────────────────────────────────────
Hook          Installed  VisionSpec
────────────────────────────────────
pre-commit    ✓          ✓
pre-push      ✓          ✓
commit-msg    ✗          -
post-commit   ✗          -
────────────────────────────────────

Hooks directory: /path/to/repo/.git/hooks
```

## Backup and Restore

When installing hooks, existing hooks are backed up:

```bash
$ visionspec hooks install

Installing hooks...
  pre-commit: backed up existing → pre-commit.backup
  pre-commit: installed
  pre-push: installed

Hooks installed successfully.
```

When uninstalling, backups are restored:

```bash
$ visionspec hooks uninstall

Uninstalling hooks...
  pre-commit: removed
  pre-commit: restored from backup
  pre-push: removed

Hooks uninstalled successfully.
```

## Configuration

Hooks can be configured in `visionspec.yaml`:

```yaml
hooks:
  enabled: true
  pre_commit:
    enabled: true
    lint_only: true  # Only lint, don't evaluate
  pre_push:
    enabled: true
    require_passing: true  # Block push if evaluation fails
    min_score: 7.0  # Minimum score to pass
```

## CI Integration

For CI environments where hooks should be skipped:

```bash
# Skip hooks with git flag
git commit --no-verify -m "CI commit"
git push --no-verify

# Or set environment variable
VISIONSPEC_SKIP_HOOKS=1 git commit -m "CI commit"
```

## See Also

- [lint](lint.md) - Validate spec structure
- [eval](eval.md) - Evaluate spec quality
- [status](status.md) - Check project status
