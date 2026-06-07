---
name: visionspec-eval
description: Run quality evaluations on specifications
triggers: [eval, evaluate, assess quality, score, check quality]
---

# VisionSpec Eval

Evaluate specifications for quality, completeness, and consistency.

## Purpose

Uses LLM-as-a-Judge methodology to:

- Score specifications against rubrics
- Identify gaps and inconsistencies
- Ensure upstream traceability
- Generate actionable recommendations

## When to Use

- After drafting or synthesizing specs
- Before requesting approval
- When iteration is needed
- To measure quality over time

## Invocation

```
visionspec eval <spec-type>
visionspec eval --all
```

Or via Claude Code:

```
/eval prd
/eval all
```

## Evaluation Criteria

Each spec type has specific rubric categories:

### Common Criteria

| Criterion | Weight | Description |
|-----------|--------|-------------|
| Completeness | 25% | All required sections present |
| Consistency | 25% | Aligns with upstream specs |
| Clarity | 20% | Clear, unambiguous language |
| Traceability | 15% | References upstream requirements |
| Quality | 15% | Professional formatting |

### Spec-Specific

- **MRD**: Market sizing, competitive analysis
- **PRD**: User stories, acceptance criteria
- **TRD**: Architecture decisions, API contracts
- **UXD**: User flows, accessibility

## Output

```
Evaluating prd...

Score: 85/100 (PASS)

  Completeness:  90/100  ✓
  Consistency:   85/100  ✓
  Clarity:       80/100  ✓
  Traceability:  75/100  ⚠
  Quality:       90/100  ✓

Findings:
  ✓ All required sections present
  ✓ User stories have acceptance criteria
  ⚠ FR-3, FR-4 missing trace to MRD

Recommendations:
  - Add MR references to FR-3, FR-4
  - Consider NFR for rate limiting

Report: eval/prd.eval.json
```

## Pass Thresholds

| Score | Status |
|-------|--------|
| >= 80 | PASS |
| 60-79 | CONDITIONAL |
| < 60 | FAIL |

## Next Actions

- PASS → Run `visionspec approve <spec-type>`
- CONDITIONAL → Address recommendations, re-eval
- FAIL → Significant revision needed

## Custom Rubrics

Profiles can provide custom rubrics:

```yaml
# profiles/enterprise/rubrics/prd.rubric.yaml
categories:
  - id: security
    name: Security Requirements
    weight: 0.20
    criteria:
      - Authentication addressed
      - Authorization model defined
      - Data protection specified
```
