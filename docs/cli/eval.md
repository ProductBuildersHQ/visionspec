# eval

Evaluate specifications using LLM-as-a-Judge.

## Usage

```bash
visionspec eval [spec-type] [flags]
```

## Description

The `eval` command evaluates specification documents against rubrics using an LLM judge. It checks for quality, completeness, and adherence to best practices.

Evaluation results are saved to `eval/{spec-type}.eval.json` and include:

- Overall score (0-10)
- Pass/fail decision
- Findings with severity levels
- Category scores

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
✓ prd: 8.5/10 PASS (3 findings)

⋯ Evaluating trd...
✗ trd: 5.2/10 FAIL (7 findings)
```

## Evaluation Results

Results are saved to `eval/{spec-type}.eval.json`:

```json
{
  "spec_type": "prd",
  "score": 8.5,
  "passed": true,
  "findings": [
    {
      "category": "user_stories",
      "severity": "medium",
      "message": "User story US-3 missing acceptance criteria"
    }
  ],
  "category_scores": {
    "problem_definition": 9.0,
    "user_stories": 7.5,
    "requirements": 8.5
  }
}
```

## Severity Levels

| Severity | Description | Impact |
|----------|-------------|--------|
| `critical` | Spec cannot be used | Automatic fail |
| `high` | Major issues | Likely fail |
| `medium` | Should be addressed | May affect score |
| `low` | Minor improvements | Informational |

## Pass Criteria

A spec passes evaluation when:

- Score >= 7.0
- No critical findings
- No more than 2 high findings

## LLM Configuration

The LLM used for evaluation is configured in `visionspec.yaml`:

```yaml
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  temperature: 0.3
```

## See Also

- [synthesize](synthesize.md) - Generate specs with optional evaluation
- [approve](approve.md) - Approve specs that pass evaluation
