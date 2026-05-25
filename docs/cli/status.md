# visionspec status

Show project status and readiness.

## Synopsis

```bash
visionspec status [flags]
```

## Description

Displays the current status of a VisionSpec project, including:

- Spec existence and status
- Evaluation results
- Approval status
- Readiness gates
- Overall readiness summary

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Project name (required) |
| `--format` | | `text` | Output format: `text`, `json`, `html`, `markdown` |
| `--ci` | | `false` | Exit non-zero if not ready |

## Examples

```bash
# Terminal output
visionspec status -p user-onboarding

# JSON for programmatic use
visionspec status -p user-onboarding --format json

# HTML report
visionspec status -p user-onboarding --format html > status.html

# Markdown for embedding
visionspec status -p user-onboarding --format markdown

# CI mode
visionspec status -p user-onboarding --ci
```

## Output Formats

### Text Format

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
  mrd          source     yes      -          -*
  prd          source     yes      -          -*
  uxd          source     yes      -          -
  trd          technical  -        -          -*

  * = required

Summary:
  Total: 10, Present: 3, Evaluated: 0, Approved: 0
```

### JSON Format

```json
{
  "project": "user-onboarding",
  "path": "docs/specs/user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "readiness": {
    "ready": false,
    "summary": "Not ready: 2 blockers",
    "gates": [
      {"name": "Required specs present", "passed": true, "message": "All required specs exist"},
      {"name": "Evaluations passing", "passed": true, "message": "No blocking eval findings"},
      {"name": "Approvals obtained", "passed": false, "message": "Pending approvals"},
      {"name": "Execution spec generated", "passed": false, "message": "spec.md not generated"}
    ]
  },
  "specs": [...],
  "summary": {
    "total_specs": 10,
    "present_specs": 3,
    "evaluated_specs": 0,
    "approved_specs": 0
  }
}
```

### HTML Format

Generates a self-contained HTML page with:

- Traffic light status indicator (green/red)
- Styled readiness gates
- Spec status table
- Summary statistics

### Markdown Format

Generates GitHub-flavored markdown with:

- Emoji status indicators (:white_check_mark:, :x:)
- Formatted tables
- Suitable for embedding in README or docs

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
