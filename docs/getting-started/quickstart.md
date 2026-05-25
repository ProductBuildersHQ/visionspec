# Quick Start

This guide walks you through creating your first VisionSpec project.

## Initialize a Project

Create a new project with the standard directory structure:

```bash
visionspec init user-onboarding
```

This creates:

```
docs/specs/user-onboarding/
├── source/
│   ├── mrd.md      # Market Requirements (placeholder)
│   ├── prd.md      # Product Requirements (placeholder)
│   └── uxd.md      # User Experience Design (placeholder)
├── gtm/            # For LLM-generated GTM docs
├── technical/      # For LLM-generated technical docs
├── eval/           # For evaluation results
└── visionspec.yaml  # Project configuration
```

## Author Source Specs

Edit the source specifications in the `source/` directory:

### `source/mrd.md` - Market Requirements

```markdown
# Market Requirements Document

## Problem Statement

What problem does this solve?

## Target Audience

Who benefits from this?

## Business Goals

What business metrics will improve?
```

### `source/prd.md` - Product Requirements

```markdown
# Product Requirements Document

## User Stories

- As a user, I want to...

## Functional Requirements

1. The system shall...

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
```

### `source/uxd.md` - User Experience Design

```markdown
# User Experience Design

## User Journey

1. User opens app
2. User sees...

## Interaction Flows

Describe key interactions...
```

## Validate Your Project

Check that your project follows VisionSpec conventions:

```bash
visionspec lint user-onboarding
```

## Check Project Status

View the current status and readiness:

```bash
# Terminal output
visionspec status -p user-onboarding

# JSON format
visionspec status -p user-onboarding --format json

# Generate HTML report
visionspec status -p user-onboarding --format html > status.html
```

## Readiness Gates

The status command shows readiness gates:

| Gate | Requirement |
|------|-------------|
| Required specs present | mrd.md, prd.md, uxd.md, trd.md exist |
| Evaluations passing | No critical/high findings in evals |
| Approvals obtained | Required specs have approvals in visionspec.yaml |
| Execution spec generated | spec.md exists |

## Next Steps

Once your source specs are complete:

1. **Synthesize GTM docs**
   ```bash
   visionspec synthesize press
   visionspec synthesize faq
   visionspec synthesize narrative
   ```

2. **Synthesize technical docs**
   ```bash
   visionspec synthesize trd
   visionspec synthesize ird
   ```

3. **Run evaluations**
   ```bash
   visionspec eval --all
   ```

4. **Get approvals**
   ```bash
   visionspec approve prd
   visionspec approve trd
   ```

5. **Reconcile to execution spec**
   ```bash
   visionspec reconcile
   ```

6. **Export to target system**
   ```bash
   visionspec export speckit
   ```

7. **Extract and visualize requirement graph**
   ```bash
   visionspec graph extract
   visionspec graph export --format html
   visionspec graph query --type requirement
   ```
