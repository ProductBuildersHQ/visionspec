# Lean Startup Framework Phases

Build-Measure-Learn methodology for rapid validation.

## Overview

The Lean Startup methodology focuses on validated learning through rapid experimentation. Instead of elaborate planning, build the minimum viable product (MVP) and learn from real customer feedback.

```
Hypothesis (what we believe)
    ↓
MRD (problem validation)
    ↓
MVP PRD (minimum viable product)
    ↓
Experiment Design (how to test)
    ↓
UXD (lean prototype)
    ↓
TRD (technical approach)
    ↓
spec.md → Build → Measure → Learn → Iterate
```

## Phase 1: Hypothesis

**Goal**: Define what you believe to be true about the problem and solution.

### Hypothesis Template

```markdown
We believe that [target customer]
has a problem [achieving goal / doing task]
because [obstacle / pain point].

We believe that [solution approach]
will solve this problem by [mechanism].

We will know we are right when we see [measurable signal].
```

### Workflow

1. Ask user to articulate their hypothesis
2. Challenge assumptions - what must be true?
3. Identify riskiest assumption
4. Document in hypothesis.md or MRD preamble

### Key Questions

- What's the riskiest assumption?
- How could we be wrong?
- What's the minimum we need to learn?

## Phase 2: Problem Validation (MRD)

**Goal**: Validate that the problem exists and is worth solving.

### Workflow

```bash
visionspec create mrd -p <project>
```

Focus MRD on:

- [ ] Evidence of problem (interviews, data, observation)
- [ ] Size of opportunity (TAM, frequency, severity)
- [ ] Willingness to pay or change behavior
- [ ] Current alternatives and their shortcomings

### Validation Criteria

Before proceeding, validate:

| Signal | Evidence |
|--------|----------|
| Problem exists | Customer interviews (5+) |
| Problem is painful | Frequency and severity data |
| Worth solving | Willingness to pay/change |

## Phase 3: MVP Definition (PRD)

**Goal**: Define the smallest product that tests the hypothesis.

### MVP Principles

1. **Minimum** - Strip to essential learning
2. **Viable** - Must actually work
3. **Product** - Delivers value, not just a demo

### Workflow

```bash
visionspec create prd -p <project>
```

PRD should answer:

- What's the ONE thing this MVP must do?
- What can we cut and still learn?
- What's the fastest path to customer hands?

### Anti-Patterns

- **Feature creep**: "While we're at it..." NO. Ship minimal.
- **Perfect polish**: Good enough to learn, not to impress.
- **Analytics overload**: One key metric, not a dashboard.

## Phase 4: Experiment Design

**Goal**: Define how you will measure success.

### Key Metric

Choose ONE metric that indicates hypothesis validation:

| Metric Type | Example |
|-------------|---------|
| Activation | % users who complete core action |
| Retention | % users who return in week 2 |
| Revenue | % users who pay |
| Referral | NPS score or viral coefficient |

### Experiment Structure

```markdown
## Experiment: [Name]

**Hypothesis**: [From Phase 1]

**Metric**: [Single key metric]

**Target**: [What number validates the hypothesis?]

**Duration**: [How long to run?]

**Sample size**: [How many users needed?]

**Pass criteria**: [Metric >= X]
**Fail criteria**: [Metric < Y]
```

## Phase 5: Lean UXD

**Goal**: Design the minimal user experience for learning.

### Workflow

```bash
visionspec create uxd -p <project>
```

Focus on:

- [ ] Core user journey only
- [ ] Fast path to value
- [ ] Clear call to action
- [ ] Feedback mechanism (how users tell you it's working)

### Prototyping Options

| Fidelity | When to Use |
|----------|-------------|
| Paper/Sketch | Very early, testing concepts |
| Clickable prototype | Testing flows, no code |
| Concierge | Manual backend, real frontend |
| Wizard of Oz | Fake automation, appears real |
| MVP | Minimal real implementation |

## Phase 6: Technical (TRD)

**Goal**: Technical approach for building MVP fast.

### Workflow

```bash
visionspec context gather -p <project>
visionspec synthesize trd -p <project>
```

Technical decisions for MVP:

- [ ] What's fastest to build?
- [ ] What can we fake/manual?
- [ ] What technical debt is acceptable?
- [ ] What's the throwaway plan?

### Build vs. Fake Matrix

| Feature | Build if... | Fake if... |
|---------|-------------|------------|
| Core value | Central to hypothesis | Never fake core value |
| Edge cases | High frequency | Low frequency |
| Scale | Testing scale hypothesis | Not testing scale |
| Polish | Retention is metric | Activation is metric |

## Phase 7: Build-Measure-Learn Loop

**Goal**: Execute experiment and learn.

### Build

```bash
visionspec reconcile -p <project>
visionspec export speckit -p <project>
# Implement MVP
```

### Measure

- Track key metric
- Collect qualitative feedback
- Watch for unexpected signals

### Learn

After experiment completes:

| Result | Action |
|--------|--------|
| Pass | Invest more, expand scope |
| Fail (metric) | Pivot or iterate hypothesis |
| Fail (engagement) | Wrong solution or wrong audience |
| Inconclusive | Refine experiment, larger sample |

### Iteration

Update specs based on learning:

```bash
# Revise hypothesis
# Update MRD with learnings
# Modify PRD for next iteration
# Repeat
```

## Lean Startup Gates

| Gate | Criteria |
|------|----------|
| Problem validated | 5+ customer interviews, evidence |
| MVP scoped | Single hypothesis, minimal features |
| Experiment designed | Key metric, pass/fail criteria |
| Learning documented | Results, insights, next steps |

## See Also

- [The Lean Startup](http://theleanstartup.com/) - Eric Ries
- [Running Lean](https://leanstack.com/running-lean-book) - Ash Maurya
