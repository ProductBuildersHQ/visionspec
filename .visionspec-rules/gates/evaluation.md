# Evaluation Gate

Evaluation ensures specs meet quality standards before approval. Every spec must pass evaluation.

## Evaluation Command

```bash
visionspec eval <type> -p <project>
```

Or via MCP:

```
eval_draft(project, specType)  # For drafts
run_eval(project, specType)    # For finalized specs
```

## Passing Criteria

A spec passes evaluation when:

| Criterion | Threshold |
|-----------|-----------|
| Overall score | >= 7.0 |
| Critical findings | 0 |
| High findings | 0 |
| Medium findings | Acknowledged |
| Low findings | Optional |

## Rubric Structure

Each spec type has a domain-specific rubric with categories:

### MRD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| Problem Definition | 25% | Clarity, specificity, measurability |
| Target Audience | 25% | Segmentation, characteristics, size |
| Business Goals | 25% | Alignment, metrics, timeline |
| Constraints | 25% | Technical, business, regulatory |

### PRD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| User Stories | 25% | Completeness, testability, priority |
| Requirements | 25% | Clarity, feasibility, traceability |
| Acceptance Criteria | 25% | Measurability, coverage |
| Scope Definition | 25% | In/out scope clarity |

### UXD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| User Journeys | 25% | Completeness, error handling |
| Interaction Design | 25% | Clarity, consistency |
| Accessibility | 25% | WCAG compliance, keyboard nav |
| Responsive Design | 25% | Mobile considerations |

### TRD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| Architecture | 25% | Clarity, scalability, integration |
| API Design | 25% | Consistency, documentation |
| Security | 25% | Auth, data protection, threats |
| Performance | 25% | Targets, benchmarks |

### TPD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| PRD Coverage | 25% | Acceptance criteria tests |
| TRD Coverage | 25% | API and integration tests |
| UXD Coverage | 25% | E2E and journey tests |
| Automation | 25% | CI/CD integration |

### IRD Rubric

| Category | Weight | Focus |
|----------|--------|-------|
| Infrastructure | 25% | Resources, redundancy |
| Deployment | 25% | Strategy, rollback |
| Operations | 25% | Monitoring, alerting |
| Security | 25% | Network, access control |

## Finding Severity

| Severity | Definition | Action |
|----------|------------|--------|
| Critical | Fundamental flaw, blocks progress | Must fix |
| High | Significant issue, risk to success | Must fix |
| Medium | Notable concern, should address | Should fix |
| Low | Minor improvement opportunity | Optional |

## Evaluation Workflow

### Step 1: Run Evaluation

```bash
visionspec eval mrd -p myproject
```

### Step 2: Review Results

```
Evaluation Results: mrd
Score: 6.5/10.0 (NEEDS IMPROVEMENT)

Findings:
  [HIGH] Problem Definition: Problem statement is vague
    → "Users struggle with X" - what specifically do they struggle with?

  [MEDIUM] Business Goals: Missing success metrics
    → Add quantifiable goals (e.g., "reduce time by 50%")

  [LOW] Constraints: Could add more detail on timeline
```

### Step 3: Address Findings

For each finding:

1. Understand the issue
2. Propose fix to user
3. Update spec
4. Re-evaluate

### Step 4: Confirm Pass

```bash
visionspec eval mrd -p myproject

Evaluation Results: mrd
Score: 8.2/10.0 (PASS)

Findings:
  [LOW] Constraints: Could add more detail on timeline
```

## Custom Rubrics

Organizations can customize rubrics:

```bash
# Export default rubric for customization
visionspec profiles export enterprise ./my-rubrics

# Use custom rubrics
visionspec eval mrd -p myproject --rubric-dir ./my-rubrics
```

## See Also

- [approval.md](approval.md) - Approval process after evaluation
- [../core-workflow.md](../core-workflow.md) - Overall workflow
