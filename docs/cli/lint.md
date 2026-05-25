# visionspec lint

Validate project directory structure and naming conventions.

## Synopsis

```bash
visionspec lint [project] [flags]
```

## Description

Validates that the project follows VisionSpec conventions:

- Directory structure matches canonical layout
- File naming follows conventions (lowercase specs, kebab-case projects)
- Required directories exist
- Config file is valid

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `project` | No | Project name to lint (lints all if omitted) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `text` | Output format: `text` or `json` |
| `--ci` | `false` | Exit with non-zero code if lint fails |

## Examples

```bash
# Lint all projects
visionspec lint

# Lint specific project
visionspec lint user-onboarding

# JSON output for CI
visionspec lint --format json

# CI mode (exit 1 on failure)
visionspec lint --ci
```

## Validation Rules

### Directory Structure

The linter checks for the expected directory layout:

```
docs/specs/{project}/
├── source/       # Required
├── gtm/          # Required
├── technical/    # Required
├── eval/         # Required
└── visionspec.yaml  # Required
```

### Naming Conventions

| Element | Convention | Valid | Invalid |
|---------|------------|-------|---------|
| Project directory | kebab-case | `user-onboarding` | `UserOnboarding` |
| Spec files | lowercase.md | `mrd.md` | `MRD.md` |
| Eval files | lowercase.eval.json | `mrd.eval.json` | `MRD_eval.json` |

### Config File

The `visionspec.yaml` file must exist and be valid YAML with required fields:

```yaml
name: project-name      # Required
description: "..."      # Optional
version: 0.1.0          # Optional
```

## Output Formats

### Text Output

```
Linting project: user-onboarding
  [+] Directory structure valid
  [+] Config file valid
  [+] Naming conventions followed

Result: PASSED
```

### JSON Output

```json
{
  "project": "user-onboarding",
  "passed": true,
  "errors": [],
  "warnings": []
}
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | All checks passed |
| 1 | Lint failures (with `--ci` flag) |

## See Also

- [init](init.md) - Create a new project
- [status](status.md) - Check project readiness
