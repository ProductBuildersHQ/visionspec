# Steering: visionspec-eval

## Command

```
visionspec eval <type>
visionspec eval --all
```

## Purpose

Evaluate specifications for quality, completeness, and consistency using LLM-as-a-Judge.

## When to Invoke

- After drafting or synthesizing specs
- Before approval
- User asks to assess quality
- User wants feedback on spec
- When iteration is needed

## Triggers

- "evaluate prd"
- "assess quality"
- "score the spec"
- "check quality"
- "is this good enough"
- "review spec"

## Expected Output

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

| Score | Status | Action |
|-------|--------|--------|
| >= 80 | PASS | Ready for approval |
| 60-79 | CONDITIONAL | Address recommendations |
| < 60 | FAIL | Significant revision needed |

## Evaluation Criteria

| Criterion | Weight | Description |
|-----------|--------|-------------|
| Completeness | 25% | All required sections |
| Consistency | 25% | Aligns with upstream |
| Clarity | 20% | Clear language |
| Traceability | 15% | References upstream |
| Quality | 15% | Professional format |

## Follow-up Actions

| Result | Next Step |
|--------|-----------|
| PASS | `visionspec approve <type>` |
| CONDITIONAL | Update spec, re-eval |
| FAIL | Major revision, re-eval |

## Spec Types

Valid types: mrd, press, faq, prd, uxd, trd, tpd, ird

## Error Handling

- **Spec not found**: Create or synthesize first
- **Upstream missing**: Eval upstream specs first
- **Invalid format**: Run `visionspec lint` first
