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

Rubrics use the shared [structured-evaluation](https://github.com/plexusone/structured-evaluation)
rubric format — the same definition used across the ecosystem. The spec type is
taken from the filename (`prd.rubric.yaml` → `prd`).

**Flat rubric** (categorical pass/partial/fail per category):

```yaml
# rubrics/prd.rubric.yaml
id: prd-rubric
name: "Acme PRD Rubric"
description: "PRD evaluation with Acme security requirements"
version: "1.0"
passCriteria:
  minCategoriesPassing: all_required
  maxFindingsSeverity: {critical: 0, high: 0, medium: 2, low: -1}
categories:
  - id: problem_definition
    name: "Problem Definition"
    description: "Is the problem clearly articulated?"
    weight: 0.15
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["Problem is specific, measurable, and tied to user needs"]}
        - {value: partial, criteria: ["Problem is stated but lacks specificity"]}
        - {value: fail, criteria: ["Problem is unclear or missing"]}

  # Custom category for your organization
  - id: security_requirements
    name: "Security Requirements"
    description: "Are security requirements documented? (ACME POLICY)"
    weight: 0.25
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["Authentication, authorization, and data protection addressed"]}
        - {value: partial, criteria: ["Some security considerations but gaps exist"]}
        - {value: fail, criteria: ["Security requirements missing"]}
```

**Rich rubric** (weighted sub-criteria with indicators, rolled up to a score):

```yaml
# rubrics/discovery.rubric.yaml
id: discovery-rubric
name: "Discovery Rubric"
version: "1.0"
passCriteria:
  scoreThresholds: {pass: 80, partial: 60}
categories:
  - id: assumption_coverage
    name: "Assumption Coverage"
    weight: 25
    criteria:
      - id: desirability
        name: "Desirability"
        weight: 25
        pass:
          description: "Desirability assumptions are identified"
          indicators: ["customer demand cited", "willingness-to-pay evidence"]
```

### Rubric Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Rubric identifier (convention: `<spec-type>-rubric`) |
| `name` | string | Display name for the rubric |
| `description` | string | Purpose of this rubric |
| `version` | string | Rubric version |
| `evaluationType` | string | `analytic` (per-category, default) or `holistic` |
| `passCriteria` | object | What constitutes passing |
| `categories` | array | Evaluation categories |

The spec type comes from the filename, not a field (`prd.rubric.yaml` → `prd`).

### Category Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (snake_case) |
| `name` | string | Display name |
| `description` | string | What this category evaluates |
| `weight` | float | Relative importance (any positive scale; normalized against the total) |
| `required` | bool | Must this category pass? |
| `scale` | object | Categorical scale with pass/partial/fail options (flat rubrics) |
| `criteria` | array | Weighted sub-criteria, each with `pass`/`partial`/`fail` bands (rich rubrics) |

Use **either** `scale` (flat) **or** `criteria` (rich) per category.

### Scale and Pass Criteria

| Field | Type | Description |
|-------|------|-------------|
| `scale.type` | string | `categorical` (recommended), `checklist`, `binary`, or `likert` |
| `scale.options[].value` | string | `pass`, `partial`, or `fail` |
| `scale.options[].criteria` | array | What that band requires |
| `passCriteria.minCategoriesPassing` | string | `all`, `all_required`, or a number |
| `passCriteria.maxFindingsSeverity` | object | Max findings per severity (`-1` = unlimited) |
| `passCriteria.scoreThresholds` | object | Numeric `pass`/`partial` cutoffs (0-100) for rich rubrics |

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
