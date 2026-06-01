# Jobs to be Done Framework Phases

Outcome-driven innovation focused on customer jobs.

## Overview

Jobs to be Done (JTBD) focuses on what customers are trying to accomplish, not what they say they want. Customers "hire" products to do jobs. Understanding the job leads to better solutions.

```
Job Identification (what are customers trying to accomplish?)
    ↓
Job Mapping (steps in getting the job done)
    ↓
Outcome Definition (how customers measure success)
    ↓
Opportunity Analysis (underserved outcomes)
    ↓
Solution Design (address high-opportunity outcomes)
    ↓
spec.md
```

## Core Concept: The Job Statement

A job statement has a specific structure:

```
[Context] + [Motivation] + [Desired Outcome]

When [situation], I want to [motivation], so I can [outcome].
```

Example:
> "When I'm rushing to leave for work, I want to eat something nutritious, so I can have energy throughout the morning without feeling sluggish."

**Key insight**: The job is NOT "buy breakfast cereal" - that's a solution. The job is "get energy for the morning."

## Phase 1: Job Identification (MRD)

**Goal**: Identify the core functional job customers are trying to get done.

### Job Types

| Type | Description | Example |
|------|-------------|---------|
| Functional | The task to accomplish | Get energy for the morning |
| Emotional | How they want to feel | Feel good about health choices |
| Social | How they want to be perceived | Be seen as a healthy parent |

### Discovery Questions

1. What are you ultimately trying to accomplish?
2. Why is that important?
3. What's the bigger goal this serves?
4. What would success look like?

### Workflow

```bash
visionspec create mrd -p <project>
```

MRD should capture:

- [ ] Core functional job statement
- [ ] Related emotional jobs
- [ ] Related social jobs
- [ ] Job context (when, where, why triggered)

### Job Statement Validation

Good job statements are:

- [ ] Solution-agnostic (no product references)
- [ ] Stable over time (job exists before/after your product)
- [ ] Universal (applies to customer segment, not individuals)

## Phase 2: Job Mapping

**Goal**: Break down the job into discrete steps.

### Universal Job Map

Most jobs follow this pattern:

| Step | Description | Questions |
|------|-------------|-----------|
| 1. Define | Determine goals and plan | What do I want to achieve? |
| 2. Locate | Gather inputs needed | What do I need to get started? |
| 3. Prepare | Set up for execution | How do I get ready? |
| 4. Confirm | Verify readiness | Am I set up correctly? |
| 5. Execute | Perform the job | How do I do the main task? |
| 6. Monitor | Track progress | Is it working? |
| 7. Modify | Make adjustments | How do I correct course? |
| 8. Conclude | Finish the job | How do I complete it? |

### Workflow

Map each step for the core job:

```markdown
## Job Map: [Job Statement]

### 1. Define
- What goals are customers setting?
- What triggers the job?

### 2. Locate
- What inputs do they need?
- Where do they get them?

[Continue for all 8 steps...]
```

### Output

Include job map in MRD or separate document.

## Phase 3: Outcome Definition

**Goal**: Identify how customers measure success at each job step.

### Outcome Statement Structure

```
[Direction] + [Metric] + [Context]

Minimize the time it takes to [job step action].
Minimize the likelihood of [negative event].
Increase the ability to [positive capability].
```

### Outcome Types

| Direction | Use When |
|-----------|----------|
| Minimize time | Speed matters |
| Minimize likelihood | Risk/error matters |
| Minimize effort | Ease matters |
| Increase ability | Capability matters |

### Workflow

For each job step, identify 3-5 outcomes:

```markdown
## Outcomes: [Job Step]

1. Minimize the time it takes to [action]
2. Minimize the likelihood of [negative outcome]
3. Increase the ability to [positive capability]
4. Minimize the effort required to [action]
```

### Outcome Importance Rating

Have customers rate outcomes:

| Rating | Meaning |
|--------|---------|
| 1-3 | Not important |
| 4-6 | Somewhat important |
| 7-8 | Important |
| 9-10 | Critical |

## Phase 4: Opportunity Analysis

**Goal**: Find underserved outcomes with high opportunity scores.

### Opportunity Score

```
Opportunity = Importance + (Importance - Satisfaction)
```

Where Importance and Satisfaction are both rated 1-10.

| Score | Interpretation |
|-------|---------------|
| < 10 | Over-served (don't focus here) |
| 10-12 | Appropriately served |
| 12-15 | Underserved (opportunity) |
| > 15 | Highly underserved (big opportunity) |

### Opportunity Landscape

Plot outcomes on a matrix:

```
        High Importance
             |
  Over-      |     Under-
  served     |     served
             |    (FOCUS)
  -----------+------------
             |
  Don't      |   Low
  bother     |   priority
             |
        Low Importance

        High ← Satisfaction → Low
```

### Workflow

```markdown
## Opportunity Analysis

| Outcome | Importance | Satisfaction | Score |
|---------|------------|--------------|-------|
| Minimize time to X | 9 | 4 | 14 |
| Minimize risk of Y | 7 | 7 | 7 |
| Increase ability to Z | 8 | 3 | 13 |

**Top Opportunities:**
1. Minimize time to X (Score: 14)
2. Increase ability to Z (Score: 13)
```

## Phase 5: Solution Design (PRD + UXD)

**Goal**: Design solutions that address high-opportunity outcomes.

### Solution Strategy

| Strategy | When to Use |
|----------|-------------|
| Discrete | Address one underserved outcome |
| Dominant | Address multiple underserved outcomes |
| Disruptive | Address overserved outcomes at lower cost |
| Sustaining | Improve on existing solutions |

### Workflow

```bash
visionspec create prd -p <project>
```

PRD should:

- [ ] List target outcomes with opportunity scores
- [ ] Map features to outcomes they address
- [ ] Prioritize features by outcome impact
- [ ] Define success metrics tied to outcomes

### Outcome-Based Requirements

For each feature:

```markdown
## Feature: [Name]

**Target Outcomes:**
- Minimize time to [X] (currently 9 importance, 4 satisfaction)
- Increase ability to [Y] (currently 8 importance, 3 satisfaction)

**Success Criteria:**
- Reduce time from [current] to [target]
- Enable [capability] that [percentage] of users can achieve

**Measurement:**
- Time tracking in app
- User survey on satisfaction
```

### UXD for Jobs

```bash
visionspec create uxd -p <project>
```

UXD should:

- [ ] Map user journey to job steps
- [ ] Show how each interaction addresses outcomes
- [ ] Minimize friction in job execution
- [ ] Support job completion, not feature discovery

## Phase 6: Technical and Build

**Goal**: Build solution that delivers on outcome promises.

### Workflow

```bash
visionspec context gather -p <project>
visionspec synthesize trd -p <project>
visionspec synthesize tpd -p <project>
visionspec synthesize ird -p <project>
visionspec reconcile -p <project>
```

### Outcome-Driven Testing (TPD)

Test cases should validate outcomes:

```markdown
## Test: Outcome Validation

**Outcome**: Minimize time to [X]

**Test Case**:
1. User attempts to [job step]
2. Measure time taken
3. Compare to target

**Pass Criteria**: Time <= [target]
```

## JTBD Gates

| Gate | Criteria |
|------|----------|
| Job identified | Clear job statement, validated with customers |
| Job mapped | All 8 steps documented with outcomes |
| Opportunities scored | Survey data, opportunity scores calculated |
| Solution addresses outcomes | Features mapped to high-opportunity outcomes |
| Outcomes validated | Post-launch measurement confirms improvement |

## See Also

- [Jobs to be Done](https://jobs-to-be-done-book.com/) - Anthony Ulwick
- [Competing Against Luck](https://www.christenseninstitute.org/books/competing-against-luck/) - Clayton Christensen
- [Outcome-Driven Innovation](https://strategyn.com/outcome-driven-innovation-process/)
