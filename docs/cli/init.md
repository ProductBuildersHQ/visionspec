# visionspec init

Initialize a new VisionSpec project.

## Synopsis

```bash
visionspec init <project-name> [flags]
```

## Description

Creates a new project with the standard VisionSpec directory structure:

```
docs/specs/{project}/
├── source/
├── gtm/
├── technical/
├── eval/
└── visionspec.yaml
```

The `source/`/`gtm`/`technical`/`eval` directories are created empty —
pass `--with-templates` to scaffold every required spec's template file
into them immediately, or use [`create <type>`](create.md) one at a time.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `project-name` | Yes | Name of the project (kebab-case) |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--profile` | string | none | Configuration profile — see `visionspec profiles list` for all 25 |
| `--with-templates` | bool | `false` | Scaffold every required spec's template file immediately |
| `--constitution` | string | `../CONSTITUTION.md` | Path to a constitution file (relative or absolute) |
| `--workflow` | string | none | External spec-workflows repository workflow ID — a separate, older mechanism from `--profile`; most projects only need `--profile` |

## Examples

```bash
# Create a new project with no profile (legacy Working Backwards fallback)
visionspec init user-onboarding

# Create a project using an embedded profile, with templates scaffolded
visionspec init payment-integration --profile aws-two-way-door --with-templates

# See the full profile catalog first
visionspec profiles list
```

## Created Files

### `visionspec.yaml`

Project configuration file:

```yaml
# visionspec project configuration
# See: https://github.com/ProductBuildersHQ/visionspec

name: payment-integration
path: /path/to/docs/specs/payment-integration
constitution: ../CONSTITUTION.md
workflow: aws-two-way-door
targets:
    default: speckit
    speckit:
        enabled: true
        branch_numbering: sequential
created_at: 2026-08-16T22:07:32.172906-07:00
updated_at: 2026-08-16T22:07:32.172906-07:00
```

`workflow` records the resolved `--profile` — every later command
([`synthesize`](synthesize.md), [`create`](create.md),
[`workflow`](workflow.md)) reads it back to know which profile's rules
apply. It's empty when `--profile` wasn't given, and later commands fall
back to the legacy Working Backwards defaults.

### With `--with-templates`

Scaffolds a template file for every spec the profile marks required,
routed into `source/`, `gtm/`, or `technical/` by category. Without a
`--profile`, this scaffolds the legacy default set: `mrd.md`, `prd.md`,
`uxd.md`.

## Project Naming

Project names should be:

- Lowercase
- Hyphen-separated (kebab-case)
- Descriptive

**Good:** `user-onboarding`, `payment-flow`, `admin-dashboard`

**Bad:** `UserOnboarding`, `payment_flow`, `ADMIN`

## See Also

- [lint](lint.md) - Validate project structure
- [status](status.md) - Check project readiness
