# workflows

List available workflows from a spec-workflows repository.

## Synopsis

```bash
visionspec workflows [flags]
```

## Description

The `workflows` command lists all available workflow methodologies from a spec-workflows repository. It auto-discovers the repository using a search order similar to how Claude Code finds `CLAUDE.md` files.

### Auto-Discovery Search Order

1. `--workflows-repo` flag (explicit override)
2. `VISIONSPEC_WORKFLOWS_REPO` environment variable
3. Walk up from current directory looking for `spec-workflows/` or `.spec-workflows/`
4. `~/.config/visionspec/spec-workflows/` (user default)

## Flags

| Flag | Description |
|------|-------------|
| `--workflows-repo` | Path to spec-workflows repository (global flag) |

## Examples

### List Available Workflows

```bash
# Auto-discover and list workflows
visionspec workflows
```

Output:

```
Repository: /home/user/.config/visionspec/spec-workflows

Available workflows (6):

  aws-working-backwards:
    - aws-working-backwards/product
    - aws-working-backwards/feature

  big-tech:
    - big-tech/product
    - big-tech/feature

  lean-startup:
    - lean-startup/product
    - lean-startup/feature

Use with: visionspec init <project> --workflow=<workflow-id>
```

### With Explicit Repository Path

```bash
visionspec workflows --workflows-repo=/path/to/my-org-workflows
```

### When No Repository Found

```bash
visionspec workflows
```

Output:

```
No spec-workflows repository found.

Search locations (in order):
  1. --workflows-repo flag
  2. VISIONSPEC_WORKFLOWS_REPO environment variable
  3. spec-workflows/ or .spec-workflows/ in current or parent directories
  4. ~/.config/visionspec/spec-workflows/

To get started:
  git clone https://github.com/ProductBuildersHQ/spec-workflows ~/.config/visionspec/spec-workflows
```

## Workflow Structure

Workflows follow the format `<methodology>/<level>`:

| Component | Description | Examples |
|-----------|-------------|----------|
| Methodology | The product development approach | `aws-working-backwards`, `big-tech`, `lean-startup` |
| Level | Product or feature scope | `product`, `feature` |

### Product vs Feature Level

| Level | Starting Spec | Use Case |
|-------|---------------|----------|
| `product` | MRD + 6-pager narrative | New product lines, major initiatives |
| `feature` | OpportunitySpec (12-box) | Features on existing products |

## Setting Up spec-workflows

### User-Level Installation

```bash
# Clone to user config directory (auto-discovered)
git clone https://github.com/ProductBuildersHQ/spec-workflows ~/.config/visionspec/spec-workflows
```

### Project-Level Installation

```bash
# Clone alongside your project
git clone https://github.com/ProductBuildersHQ/spec-workflows ./spec-workflows
```

### Organization Fork

```bash
# Fork and customize for your organization
git clone https://github.com/YourOrg/spec-workflows ~/.config/visionspec/spec-workflows
```

## Environment Variable

```bash
# Set in your shell profile
export VISIONSPEC_WORKFLOWS_REPO=/path/to/spec-workflows
```

## Related Commands

- [`init`](init.md) - Initialize a project with `--workflow` flag
- [`profiles`](profiles.md) - List configuration profiles

## See Also

- [spec-workflows Repository](https://github.com/ProductBuildersHQ/spec-workflows)
- [Custom Profiles Guide](../guides/custom-profiles.md)
