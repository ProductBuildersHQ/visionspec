# metrics

View evaluation and reconciliation metrics dashboard.

## Synopsis

```bash
visionspec metrics [flags]
```

## Description

The `metrics` command displays a dashboard of evaluation and reconciliation metrics for a project. It aggregates scores, findings, and trends across all specification types.

## Flags

| Flag | Description |
|------|-------------|
| `-p, --project` | Project name |
| `--format` | Output format: text, json, html (default: text) |
| `--since` | Show metrics since date (YYYY-MM-DD) |
| `--compare` | Compare with previous evaluation |

## Examples

### Basic Metrics

```bash
# Show metrics for project
visionspec metrics -p myproject
```

### Output Formats

```bash
# Text output (default)
visionspec metrics -p myproject

# JSON output for programmatic use
visionspec metrics -p myproject --format json

# HTML output for reports
visionspec metrics -p myproject --format html > metrics.html
```

### Historical Comparison

```bash
# Show metrics since specific date
visionspec metrics -p myproject --since 2024-01-01

# Compare with previous evaluation
visionspec metrics -p myproject --compare
```

## Dashboard Output

### Text Format

```
╭─────────────────────────────────────────╮
│         Evaluation Metrics              │
│         user-onboarding                 │
├─────────────────────────────────────────┤
│ Spec       Score    Findings   Status   │
│ ──────────────────────────────────────  │
│ MRD        8.5      2 low      ✓ Pass   │
│ PRD        7.8      1 medium   ✓ Pass   │
│ UXD        7.2      3 medium   ✓ Pass   │
│ TRD        8.1      1 low      ✓ Pass   │
│ IRD        7.5      2 low      ✓ Pass   │
├─────────────────────────────────────────┤
│ Overall    7.8      9 total    ✓ Ready  │
╰─────────────────────────────────────────╯

Reconciliation Status:
  Last reconciled: 2024-01-15T10:30:00Z
  Conflicts resolved: 3
  Manual interventions: 1
```

## Metrics Report

Example JSON output:

```json
{
  "project": "user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "evaluations": {
    "mrd": {
      "score": 8.5,
      "findings": {"low": 2, "medium": 0, "high": 0, "critical": 0},
      "passing": true
    },
    "prd": {
      "score": 7.8,
      "findings": {"low": 0, "medium": 1, "high": 0, "critical": 0},
      "passing": true
    }
  },
  "reconciliation": {
    "last_run": "2024-01-15T10:30:00Z",
    "conflicts_resolved": 3,
    "manual_interventions": 1
  },
  "summary": {
    "average_score": 7.8,
    "total_findings": 9,
    "specs_passing": 5,
    "specs_failing": 0,
    "ready_for_export": true
  }
}
```

## Metric Categories

### Evaluation Metrics

| Metric | Description |
|--------|-------------|
| Score | Overall quality score (0-10) |
| Findings | Issues found by severity |
| Pass/Fail | Whether spec meets quality gate |

### Reconciliation Metrics

| Metric | Description |
|--------|-------------|
| Conflicts | Number of conflicts between specs |
| Resolved | Automatically resolved conflicts |
| Manual | Conflicts requiring manual intervention |

### Trend Metrics

| Metric | Description |
|--------|-------------|
| Score Delta | Change from previous evaluation |
| New Findings | Issues introduced since last run |
| Fixed | Issues resolved since last run |

## See Also

- [eval](eval.md) - Run spec evaluation
- [reconcile](reconcile.md) - Generate unified spec
- [status](status.md) - Check project status
