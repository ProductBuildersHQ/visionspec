# visionspec init

Initialize a new VisionSpec project.

## Synopsis

```bash
visionspec init <project-name>
```

## Description

Creates a new project with the standard VisionSpec directory structure:

```
docs/specs/{project}/
├── source/
│   ├── mrd.md
│   ├── prd.md
│   └── uxd.md
├── gtm/
├── technical/
├── eval/
└── visionspec.yaml
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `project-name` | Yes | Name of the project (kebab-case) |

## Examples

```bash
# Create a new project
visionspec init user-onboarding

# Create another project
visionspec init payment-integration
```

## Created Files

### `visionspec.yaml`

Project configuration file:

```yaml
name: user-onboarding
description: ""
version: 0.1.0
```

### Source Spec Placeholders

The `source/` directory contains placeholder markdown files:

- `mrd.md` - Market Requirements Document
- `prd.md` - Product Requirements Document
- `uxd.md` - User Experience Design

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
