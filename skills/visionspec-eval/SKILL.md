# VisionSpec Eval

Run quality evaluations on specifications.

## Overview

The eval command uses LLM-as-a-Judge methodology to score specifications against rubrics, identify gaps, and generate actionable recommendations.

## Workflow

### 1. Run Evaluation

Evaluate a single spec:

```
visionspec eval prd
```

Or all specs:

```
visionspec eval --all
```

### 2. Review Score

Scores are calculated across weighted criteria:

| Criterion | Weight | What We Look For |
|-----------|--------|------------------|
| Completeness | 25% | All required sections present |
| Consistency | 25% | Aligns with upstream specs |
| Clarity | 20% | Clear, unambiguous language |
| Traceability | 15% | References upstream requirements |
| Quality | 15% | Professional formatting |

### 3. Interpret Results

| Score | Status | Action |
|-------|--------|--------|
| >= 80 | PASS | Ready for approval |
| 60-79 | CONDITIONAL | Address recommendations |
| < 60 | FAIL | Significant revision needed |

### 4. Address Findings

For each finding:

1. Review the specific issue identified
2. Update the spec to address it
3. Re-run evaluation
4. Repeat until passing

### 5. Approve

Once evaluation passes:

```
visionspec approve prd
```

## Spec-Specific Criteria

Each spec type has additional criteria:

**MRD**

- Market sizing accuracy
- Competitive analysis depth

**PRD**

- User stories have acceptance criteria
- Requirements are testable

**TRD**

- Architecture decisions documented
- API contracts defined

**UXD**

- User flows complete
- Accessibility addressed

## Custom Rubrics

Profiles can provide custom rubrics:

```yaml
# profiles/enterprise/rubrics/prd.rubric.yaml
categories:
  - id: security
    name: Security Requirements
    weight: 0.20
```

## Tips

- Run eval early and often during drafting
- Address high-severity findings first
- Use eval feedback to guide revisions
- Custom rubrics ensure org-specific quality

## Related Skills

- `visionspec-synthesize`: Generate specs to evaluate
- `author-*`: Author specs to evaluate
- `visionspec-status`: Check which specs need evaluation
