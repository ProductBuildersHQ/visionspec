---
name: eval
description: Evaluate a specification for quality and consistency
arguments: [type]
dependencies: [working-backwards]
---

# Evaluate Specification

Run quality evaluation on a specification.

## Usage

```
/eval <type>
/eval all
```

Where `<type>` is: mrd, press, faq, prd, uxd, trd, tpd, ird

## Process

1. **Load Spec** - Read the specification to evaluate
2. **Load Upstream** - Read upstream specs for consistency checking
3. **Apply Criteria** - Score each evaluation criterion
4. **Generate Report** - Create detailed evaluation report
5. **Display Results** - Show summary with pass/fail status

## Evaluation Criteria

| Criterion | Weight | Description |
|-----------|--------|-------------|
| Completeness | 25% | All required sections present |
| Consistency | 25% | No contradictions with upstream |
| Clarity | 20% | Clear, unambiguous language |
| Traceability | 15% | References to upstream requirements |
| Quality | 15% | Professional formatting |

## Output

```
⋯ Evaluating {type}...

✓ Evaluation complete: PASS (85/100)

  Completeness:  90/100  ✓
  Consistency:   85/100  ✓
  Clarity:       80/100  ✓
  Traceability:  75/100  ⚠
  Quality:       90/100  ✓

  Findings:
    - All sections present and complete
    - Consistent with upstream MRD
    - Some traceability links missing

  Recommendations:
    - Add traces for FR-3, FR-4 to MR requirements
    - Consider expanding NFR section

  Report: docs/specs/{project}/eval/{type}.eval.json

Next steps:
  - Address recommendations (optional)
  - Run `/approve {type}` to approve
```

## Pass/Fail Thresholds

| Score | Status |
|-------|--------|
| >= 80 | PASS |
| 60-79 | CONDITIONAL (requires review) |
| < 60 | FAIL |

## Evaluate All

Running `/eval all` evaluates all existing specs in dependency order:

```
⋯ Evaluating all specs...

  mrd:   PASS (85/100)
  press: PASS (82/100)
  faq:   PASS (78/100) ⚠ Conditional
  prd:   PASS (88/100)
  trd:   PASS (90/100)

  Overall: 4 PASS, 1 CONDITIONAL, 0 FAIL

  See individual reports in docs/specs/{project}/eval/
```
