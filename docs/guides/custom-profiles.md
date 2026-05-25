# Custom Profiles Guide

This guide explains how to create custom profiles with organization-specific templates and rubrics.

## Overview

Profiles bundle three things:

1. **Spec Configuration** - Which specs are required and their categories
2. **Templates** - Markdown templates for each spec type
3. **Rubrics** - Evaluation criteria for each spec type

You can customize any or all of these to match your organization's needs.

## Quick Start

The fastest way to create a custom profile is to export an existing one:

```bash
# Export the enterprise profile as a starting point
visionspec profiles export enterprise ./my-profile

# Customize the files
vim ./my-profile/profile.yaml
vim ./my-profile/templates/prd.md
vim ./my-profile/rubrics/prd.rubric.yaml

# Use your custom profile
visionspec init my-project --profile-dir ./my-profile
```

## Profile Directory Structure

A profile directory must follow this structure:

```
my-profile/
├── profile.yaml           # Required: spec configuration
├── templates/             # Optional: custom templates
│   ├── prd.md
│   ├── mrd.md
│   └── {spec-type}.md
└── rubrics/               # Optional: custom rubrics
    ├── prd.rubric.yaml
    ├── mrd.rubric.yaml
    └── {spec-type}.rubric.yaml
```

## File Naming Conventions

**Templates and rubrics must use specific naming patterns:**

| Resource | Pattern | Examples |
|----------|---------|----------|
| Templates | `{spec-type}.md` | `prd.md`, `mrd.md`, `security.md` |
| Rubrics | `{spec-type}.rubric.yaml` | `prd.rubric.yaml`, `security.rubric.yaml` |

The spec-type in the filename must match:

- The `spec_type` field in rubric YAML files
- The spec name in `profile.yaml`
- The spec type used in visionspec commands

## profile.yaml Reference

```yaml
# Required: unique identifier
name: my-enterprise

# Required: human-readable description
description: "Custom enterprise profile with security requirements"

# Optional: inherit from another profile
extends: enterprise

# Spec configuration
spec_config:
  # Source specs
  prd:
    required: true
    category: source
  mrd:
    required: true
    category: source
  uxd:
    required: false    # Optional spec
    category: source

  # GTM specs
  press:
    required: true
    category: gtm
  faq:
    required: true
    category: gtm

  # Technical specs
  trd:
    required: true
    category: technical

  # Custom spec types
  security-review:
    required: true
    category: technical
```

### Spec Categories

| Category | Purpose | Examples |
|----------|---------|----------|
| `source` | Human-authored input specs | prd, mrd, uxd |
| `gtm` | Go-to-market specs | press, faq, narrative |
| `technical` | Technical design specs | trd, ird |
| `output` | Generated output | spec |

## Creating Custom Templates

Templates are Markdown files with optional placeholders:

```markdown
# Product Requirements Document (PRD)

**Project:** {project_name}
**Author:** {author}
**Date:** {date}

## 1. Problem Statement

<!-- Describe the problem being solved -->

## 2. User Stories

<!-- ACME REQUIREMENT: All user stories must include security acceptance criteria -->

### US-1: [Title]

**As a** [user type]
**I want** [capability]
**So that** [benefit]

**Acceptance Criteria:**
- [ ] Functional criteria
- [ ] Security criteria (required)

## 3. Security Requirements

<!-- ACME POLICY: This section is mandatory -->

### Authentication
### Authorization
### Data Protection
```

### Available Placeholders

| Placeholder | Replaced With |
|-------------|---------------|
| `{project_name}` | Project name from init |
| `{author}` | Current user (if available) |
| `{date}` | Current date (YYYY-MM-DD) |

## Creating Custom Rubrics

Rubrics define evaluation criteria in YAML format:

```yaml
# rubrics/prd.rubric.yaml
spec_type: prd
name: "Acme PRD Rubric"
description: "PRD evaluation with Acme security requirements"
version: "1.0"

categories:
  - id: problem_definition
    name: "Problem Definition"
    description: "Is the problem clearly articulated?"
    weight: 0.15
    required: true
    criteria:
      pass: "Problem is specific, measurable, and tied to user needs"
      partial: "Problem is stated but lacks specificity"
      fail: "Problem is unclear or missing"

  - id: user_stories
    name: "User Stories"
    description: "Are user stories complete with acceptance criteria?"
    weight: 0.20
    required: true
    criteria:
      pass: "All stories follow format with testable acceptance criteria"
      partial: "Some stories missing criteria"
      fail: "User stories missing or inadequate"

  # Custom category for your organization
  - id: security_requirements
    name: "Security Requirements"
    description: "Are security requirements documented? (ACME POLICY)"
    weight: 0.25
    required: true
    criteria:
      pass: "Authentication, authorization, and data protection addressed"
      partial: "Some security considerations but gaps exist"
      fail: "Security requirements missing"

pass_criteria:
  require_all_pass: false
  max_critical: 0
  max_high: 0
  max_medium: 2
```

### Rubric Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec_type` | string | Must match filename (e.g., `prd` for `prd.rubric.yaml`) |
| `name` | string | Display name for the rubric |
| `description` | string | Purpose of this rubric |
| `version` | string | Rubric version |
| `categories` | array | Evaluation categories |
| `pass_criteria` | object | What constitutes passing |

### Category Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (snake_case) |
| `name` | string | Display name |
| `description` | string | What this category evaluates |
| `weight` | float | Relative importance (0.0-1.0) |
| `required` | bool | Must this category pass? |
| `criteria.pass` | string | What constitutes passing |
| `criteria.partial` | string | What constitutes partial pass |
| `criteria.fail` | string | What constitutes failure |

### Pass Criteria

| Field | Type | Description |
|-------|------|-------------|
| `require_all_pass` | bool | All categories must pass |
| `max_critical` | int | Maximum critical findings allowed |
| `max_high` | int | Maximum high findings allowed |
| `max_medium` | int | Maximum medium findings allowed (-1 = unlimited) |

## Profile Inheritance

Profiles can extend other profiles:

```yaml
name: acme-startup
description: "Acme startup profile with security baseline"
extends: startup

spec_config:
  # Inherits prd: required from startup
  # Add security requirement
  security-review:
    required: true
    category: technical
```

When extending:

- Spec config is merged (child overrides parent)
- Templates are chained (child checked first, then parent)
- Rubrics are chained (child checked first, then parent)

## Building a Custom CLI

For distribution, you can compile profiles into a custom binary:

```go
package main

import (
    "embed"
    "github.com/ProductBuildersHQ/visionspec/pkg/cli"
    "github.com/ProductBuildersHQ/visionspec/pkg/templates"
    "github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
    "github.com/spf13/cobra"
)

//go:embed templates/*.md
var orgTemplates embed.FS

//go:embed rubrics/*.rubric.yaml
var orgRubrics embed.FS

func main() {
    root := &cobra.Command{Use: "acme-spec"}

    cfg := cli.DefaultConfig()
    cfg.TemplateLoader = templates.NewChainLoader(
        templates.NewEmbedFSLoader(orgTemplates, "templates"),
        templates.EmbeddedLoader(),  // Fallback to defaults
    )
    cfg.RubricLoader = rubrics.NewChainLoader(
        rubrics.NewEmbedFSLoader(orgRubrics, "rubrics"),
        rubrics.EmbeddedLoader(),  // Fallback to defaults
    )

    cli.AddCommandsTo(root, cfg)
    root.Execute()
}
```

Build with:

```bash
go build -o acme-spec
```

The resulting binary contains all templates and rubrics - no external files needed.

## Loader Types

| Loader | Source | Use Case |
|--------|--------|----------|
| `EmbeddedLoader()` | Built-in defaults | Fallback to visionspec defaults |
| `NewEmbedFSLoader(fs, dir)` | `embed.FS` | Compile into binary |
| `NewFileLoader(dir)` | Filesystem | Runtime loading |
| `NewChainLoader(...)` | Multiple | Try loaders in order |

## Best Practices

1. **Start from an existing profile** - Export and modify rather than starting from scratch

2. **Document your additions** - Add comments explaining organization-specific requirements

3. **Use meaningful weights** - Higher weights for more important criteria

4. **Test your rubrics** - Run evaluations on sample specs to verify criteria

5. **Version your profiles** - Use git to track profile changes

6. **Chain loaders for fallback** - Always fall back to defaults for specs you haven't customized

## See Also

- [profiles command](../cli/profiles.md) - CLI reference
- [Evaluation System](../getting-started/concepts.md#evaluation-system) - How evaluation works
- [examples/org-cli](https://github.com/ProductBuildersHQ/visionspec/tree/main/examples/org-cli) - Complete example
