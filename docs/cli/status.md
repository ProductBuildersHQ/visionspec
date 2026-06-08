# visionspec status

Show project status and readiness.

## Synopsis

```bash
visionspec status [flags]
```

## Description

Displays the current status of a VisionSpec project, including:

- Pipeline progress visualization
- Spec existence and status
- Category breakdown and findings
- Readiness gates
- Overall readiness summary

The default output is optimized for AI agents with pipeline visualization and box-drawing tables.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Project name (required) |
| `--format` | | `text` | Output format: `text`, `json`, `markdown` |
| `--basic` | | `false` | Basic output without pipeline visualization |
| `--ci` | | `false` | Exit non-zero if not ready |

## Examples

```bash
# Rich terminal output (default, optimized for AI agents)
visionspec status -p user-onboarding

# JSON for programmatic use
visionspec status -p user-onboarding --format json

# Markdown with pipeline visualization
visionspec status -p user-onboarding --format markdown

# Basic output (legacy format)
visionspec status -p user-onboarding --basic

# CI mode
visionspec status -p user-onboarding --ci
```

## Output Formats

### Text Format (Rich - Default)

The default rich output includes pipeline visualization with status icons and box-drawing tables:

```
⏺ VisionSpec Status

  Pipeline Progress

  MRD → Press → FAQ → PRD → UXD → TRD → TPD → IRD → spec.md
   ✅     ✅     ✅    ✅    ✅    ✅    ✅    ❌      ❌

  Summary
  ┌─────────┬────────────┬───────────────────┬─────────────────┐
  │  Spec   │   Status   │    Categories     │    Findings     │
  ├─────────┼────────────┼───────────────────┼─────────────────┤
  │ MRD     │ ✅ Pass    │ 6/6 pass          │ 2 info          │
  ├─────────┼────────────┼───────────────────┼─────────────────┤
  │ PRD     │ ✅ Pass    │ 6 pass, 1 partial │ 1 low, 1 info   │
  ├─────────┼────────────┼───────────────────┼─────────────────┤
  │ IRD     │ ❌ Missing │ -                 │ -               │
  └─────────┴────────────┴───────────────────┴─────────────────┘
  Overall: 7/9 specs complete (78%)

  Aggregate: 41 pass, 5 partial, 0 fail across 46 categories
```

Status icons:

- ✅ Complete (evaluated and passing, or approved)
- 🔄 Pending (draft, needs evaluation)
- ❌ Missing (file does not exist)

### Text Format (Basic)

Use `--basic` for simplified legacy output:

```
Project: user-onboarding
Path: docs/specs/user-onboarding

Status: NOT READY
Not ready: 2 blockers

Readiness Gates:
  [+] Required specs present: All required specs exist
  [+] Evaluations passing: No blocking eval findings
  [X] Approvals obtained: Pending approvals
  [X] Execution spec generated: spec.md not generated

Specifications:
  TYPE         CATEGORY   EXISTS   EVAL       APPROVED
  ----         --------   ------   ----       --------
  mrd          source     yes      pass       yes*
  prd          source     yes      pass       yes*
  uxd          source     yes      -          -
  trd          technical  -        -          -*

  * = required

Summary:
  Total: 10, Present: 3, Evaluated: 2, Approved: 2
```

### JSON Format

The JSON format includes the full rich report structure with pipeline data:

```json
{
  "project": "user-onboarding",
  "path": "docs/specs/user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "readiness": {
    "ready": false,
    "summary": "Not ready: 2 blockers",
    "gates": [...]
  },
  "specs": [...],
  "summary": {...},
  "pipeline": [
    {"type": "mrd", "status": "complete", "label": "MRD"},
    {"type": "prd", "status": "complete", "label": "PRD"},
    {"type": "ird", "status": "missing", "label": "IRD"}
  ],
  "aggregate_categories": {
    "pass": 41,
    "partial": 5,
    "fail": 0,
    "total": 46
  },
  "completion_percent": 78
}
```

### Markdown Format

Generates GitHub-flavored markdown with pipeline visualization:

```markdown
# VisionSpec Status: user-onboarding

## Pipeline Progress

\`\`\`
MRD → Press → FAQ → PRD → UXD → TRD → TPD → IRD → spec.md
 ✅     ✅     ✅    ✅    ✅    ✅    ✅    ❌      ❌
\`\`\`

**Completion:** 7/9 specs (78%)

## Summary

| Spec | Status | Categories | Findings |
|------|--------|------------|----------|
| MRD | :white_check_mark: Pass | 6/6 pass | 2 info |
| PRD | :white_check_mark: Pass | 6p/1pt | 1 low, 1 info |
| IRD | :x: Missing | - | - |
```

## Readiness Gates

| Gate | Description |
|------|-------------|
| Required specs present | mrd, prd, uxd, trd must exist |
| Evaluations passing | No eval files with `fail` decision |
| Approvals obtained | All required specs have approval in config |
| Execution spec generated | `spec.md` file exists |

## Spec Types

| Type | Category | Required |
|------|----------|----------|
| mrd | source | Yes |
| prd | source | Yes |
| uxd | source | No |
| press | gtm | No |
| faq | gtm | No |
| narrative | gtm | No |
| trd | technical | Yes |
| ird | technical | No |
| sec | technical | No |
| spec | reconciled | No |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Command succeeded (or project is ready with `--ci`) |
| 1 | Project not ready (with `--ci` flag) |

## See Also

- [lint](lint.md) - Validate project structure
- [init](init.md) - Create a new project
