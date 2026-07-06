# eval

Evaluate specifications using LLM-as-a-Judge.

## Usage

```bash
visionspec eval [spec-type] [flags]
```

## Description

The `eval` command evaluates specification documents against rubrics using an LLM judge. It checks for quality, completeness, and adherence to best practices.

Evaluation results are saved to `eval/{spec-type}.json` and include:

- Overall score (1-5 integer scale)
- Pass/fail decision with blocking reason codes
- Per-dimension scores with confidence values
- Findings with [reason codes](../concepts/reason-codes.md) for automated repair

## Arguments

| Argument | Description |
|----------|-------------|
| `spec-type` | Specific spec type to evaluate (optional if using category flags) |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | `false` | Evaluate all specs |
| `--source` | bool | `false` | Evaluate source specs only (mrd, prd, uxd) |
| `--gtm` | bool | `false` | Evaluate GTM docs only (press, faq, narrative) |
| `--technical` | bool | `false` | Evaluate technical docs only (trd, ird) |

## Examples

```bash
# Evaluate a specific spec
visionspec eval prd

# Evaluate all specs
visionspec eval --all

# Evaluate only source specs
visionspec eval --source

# Evaluate GTM documents
visionspec eval --gtm

# Evaluate technical specs
visionspec eval --technical
```

## Output

The command prints evaluation status for each spec:

```
⋯ Evaluating prd...
✓ prd: 4/5 (Good) PASS [conf: 88%]

⋯ Evaluating trd...
✗ trd: 2/5 (Major Revisions) FAIL [conf: 75%]
  Blocking: SEC-NO_AUTH, ARCH-NO_API
```

## Score Scale (1-5)

| Score | Label | Description |
|-------|-------|-------------|
| 5 | Excellent | Exceeds expectations |
| 4 | Good | Meets expectations well |
| 3 | Acceptable | Meets minimum requirements |
| 2 | Major Revisions | Significant work needed |
| 1 | Unacceptable | Does not meet requirements |

## Evaluation Results (v2)

Results are saved to `eval/{spec-type}.json` in v2 format:

```json
{
  "schemaVersion": "v2",
  "scoreV2": 4,
  "decision": "conditional",
  "pass": false,
  "confidence": 0.78,
  "blocking": ["REQ-NO_CRITERIA", "METRIC-UNMEASURABLE"],
  "dimensions": [
    {
      "id": "requirements",
      "name": "Requirements Clarity",
      "score": 3,
      "severity": "minor",
      "confidence": 0.75,
      "reasonCodes": ["REQ-NO_CRITERIA"],
      "findings": [
        {
          "category": "requirements",
          "severity": "medium",
          "code": "REQ-NO_CRITERIA",
          "message": "FR-2.3 lacks acceptance criteria",
          "location": "FR-2.3"
        }
      ]
    }
  ],
  "findings": [...]
}
```

### Key Fields

| Field | Description |
|-------|-------------|
| `schemaVersion` | Format version ("v2") |
| `scoreV2` | Overall score (1-5 integer) |
| `pass` | Boolean pass/fail gate |
| `confidence` | Evaluation confidence (0.0-1.0) |
| `blocking` | Reason codes that caused failure |
| `dimensions` | Per-dimension breakdown with scores |

## Reason Codes

Findings include standardized [reason codes](../concepts/reason-codes.md) that enable automated repair:

| Prefix | Domain | Example |
|--------|--------|---------|
| `REQ-` | Requirements | `REQ-AMBIGUOUS`, `REQ-NO_CRITERIA` |
| `SEC-` | Security | `SEC-NO_AUTH`, `SEC-PRIVACY` |
| `UX-` | UX/Accessibility | `UX-NO_ARIA`, `UX-NO_ERROR_STATE` |
| `ARCH-` | Architecture | `ARCH-NO_API`, `ARCH-GAP` |

See [Reason Codes Reference](../concepts/reason-codes.md) for the complete list.

## Severity Levels

| Severity | Description | Impact |
|----------|-------------|--------|
| `critical` | Spec cannot be used | Automatic fail, added to `blocking` |
| `high` | Major issues | Likely fail, added to `blocking` |
| `medium` | Should be addressed | May affect score |
| `low` | Minor improvements | Informational |
| `info` | Observations | No impact on score |

## Pass Criteria

A spec passes evaluation when:

- Score >= 3 (Acceptable)
- No critical findings
- No high-severity findings (or within threshold)
- No blocking reason codes

## Confidence Values

The evaluation includes confidence scores (0.0-1.0):

| Confidence | Interpretation |
|------------|----------------|
| 0.9+ | Very confident in assessment |
| 0.7-0.9 | Confident |
| 0.5-0.7 | Somewhat confident, consider review |
| < 0.5 | Low confidence, needs human review |

Low confidence evaluations are flagged with "Needs Review" in VisionStudio.

## LLM Configuration

The LLM used for evaluation is configured in `visionspec.yaml`:

```yaml
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  temperature: 0.3
```

## See Also

- [Reason Codes Reference](../concepts/reason-codes.md) - Complete reason code documentation
- [synthesize](synthesize.md) - Generate specs with optional evaluation
- [approve](approve.md) - Approve specs that pass evaluation
