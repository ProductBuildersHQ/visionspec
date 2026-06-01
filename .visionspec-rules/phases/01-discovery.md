# Phase 1: Discovery

The Discovery phase establishes the market problem, target audience, and business goals through the Market Requirements Document (MRD).

## Objective

Create an approved MRD that clearly defines:

- The problem being solved
- Who experiences this problem
- Why solving it matters to the business

## Entry Criteria

- User has expressed intent to create a new product/feature
- Project has been initialized (`visionspec init <project>`)

## Workflow

### Step 1: Initialize Draft

```bash
visionspec create mrd -p <project>
# Or via MCP: start_draft(project, "mrd")
```

### Step 2: Discovery Questions

Ask the user these questions to populate the MRD:

**Problem Space**

1. What problem are you trying to solve?
2. How do users currently work around this problem?
3. What is the cost of not solving this problem?

**Target Audience**

4. Who experiences this problem most acutely?
5. How many users/customers are affected?
6. What are their characteristics (role, industry, size)?

**Business Goals**

7. Why is this important to your organization?
8. What metrics would indicate success?
9. What is the timeline for delivery?

**Constraints**

10. What technical/business constraints exist?
11. What is explicitly out of scope?
12. Are there regulatory or compliance requirements?

### Step 3: Fill Template Sections

Map answers to MRD template sections:

| Question | MRD Section |
|----------|-------------|
| 1-3 | Problem Statement |
| 4-6 | Target Audience |
| 7-8 | Business Goals |
| 9 | Timeline |
| 10-12 | Constraints |

### Step 4: Evaluate

```bash
visionspec eval mrd -p <project>
```

Check for:

- [ ] Problem is specific and measurable
- [ ] Audience is clearly defined
- [ ] Business goals are quantifiable
- [ ] Constraints are explicit
- [ ] Score >= 7.0

### Step 5: Iterate

If evaluation fails:

1. Review findings with user
2. Clarify ambiguous areas
3. Add missing information
4. Re-evaluate

### Step 6: Approve

```bash
visionspec approve mrd -p <project>
```

## Exit Criteria

- MRD exists at `source/mrd.md`
- Evaluation score >= 7.0
- No critical or high findings
- MRD is approved

## Next Phase

→ [Phase 2: Vision](02-vision.md)

## Anti-Patterns

- **Solution-first thinking**: User describes solution before problem. Redirect to problem definition.
- **Vague audience**: "Everyone" is not a target audience. Push for specifics.
- **Missing metrics**: Business goals without measurable outcomes. Ask "How will you know this succeeded?"
