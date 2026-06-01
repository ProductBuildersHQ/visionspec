# Google Framework Phases

OKRs, Design Docs, and data-driven experimentation.

## Overview

Google's approach combines structured goal-setting (OKRs), thorough technical planning (Design Docs), collaborative feedback (RFC), and data-driven experimentation. It's well-suited for engineering-heavy organizations.

```
OKRs (objectives and key results)
    ↓
MRD (problem and opportunity sizing)
    ↓
Design Doc (technical approach with alternatives)
    ↓
RFC (request for comments)
    ↓
PRD + UXD (refined requirements)
    ↓
Experiment Design (A/B tests, launch criteria)
    ↓
TRD → TPD → IRD
    ↓
spec.md
```

## Phase 1: OKRs

**Goal**: Connect the initiative to measurable organizational objectives.

### OKR Structure

```markdown
## Objective: [Qualitative goal - inspiring, memorable]

### Key Results:
1. [Quantitative measure] from [X] to [Y]
2. [Quantitative measure] from [X] to [Y]
3. [Quantitative measure] from [X] to [Y]
```

### OKR Principles

| Principle | Description |
|-----------|-------------|
| Ambitious | 70% achievement = success |
| Measurable | Numbers, not activities |
| Time-bound | Quarterly or annual |
| Aligned | Connects to company/team OKRs |

### Example

```markdown
## Objective: Make search results instantly useful

### Key Results:
1. Increase direct answer rate from 35% to 50%
2. Reduce time-to-information from 15s to 8s
3. Improve user satisfaction (CSAT) from 4.1 to 4.5
```

### Workflow

Document in MRD or separate OKR document:

- [ ] Which team/company OKR does this support?
- [ ] What are the 2-4 key results for this initiative?
- [ ] How will we measure each key result?

## Phase 2: Problem Statement (MRD)

**Goal**: Define the problem with data and opportunity sizing.

### Workflow

```bash
visionspec create mrd -p <project>
```

Google MRDs emphasize:

- [ ] **Data-backed problem** - Metrics showing the issue
- [ ] **Opportunity size** - TAM/SAM/SOM or impact estimate
- [ ] **User segments** - Who is affected, how much
- [ ] **Current state** - What exists today, why it's insufficient

### Data Requirements

| Element | Data Needed |
|---------|-------------|
| Problem severity | Error rates, support tickets, abandonment |
| User impact | Number of affected users, frequency |
| Business impact | Revenue, engagement, retention metrics |
| Opportunity | Potential improvement, market size |

## Phase 3: Design Doc

**Goal**: Thorough technical proposal with alternatives considered.

### Design Doc Structure

```markdown
# Design Doc: [Feature Name]

## Status
[Draft | In Review | Approved | Implemented | Deprecated]

## Authors
[Names]

## Reviewers
[Names]

## Last Updated
[Date]

---

## Overview
[1-2 paragraph summary]

## Context and Scope
- What problem are we solving?
- What is NOT in scope?
- Related work and prior art

## Goals and Non-Goals
### Goals
- [Goal 1]
- [Goal 2]

### Non-Goals
- [Explicit non-goal 1]
- [Explicit non-goal 2]

## Design

### System Architecture
[Diagrams and description]

### APIs
[Interface definitions]

### Data Model
[Schema, storage]

### Alternatives Considered
#### Alternative 1: [Name]
- Pros: [...]
- Cons: [...]
- Why not chosen: [...]

#### Alternative 2: [Name]
[Same structure]

## Security Considerations
[Threat model, mitigations]

## Privacy Considerations
[Data handling, PII, compliance]

## Metrics and Monitoring
[What we'll measure, alerts]

## Rollout Plan
[Phases, feature flags, rollback]

## Open Questions
[Unresolved issues for RFC]
```

### Alternatives Section

**Critical**: Always document alternatives considered. This shows thoughtful analysis and helps reviewers understand the decision space.

| Alternative | Pros | Cons | Verdict |
|-------------|------|------|---------|
| Option A | Fast, simple | Doesn't scale | Rejected |
| Option B | Scales well | Complex | **Selected** |
| Option C | Cheapest | Vendor lock-in | Rejected |

## Phase 4: RFC (Request for Comments)

**Goal**: Gather feedback from stakeholders before committing.

### RFC Process

1. **Publish** Design Doc to review group
2. **Announce** RFC with deadline (typically 1-2 weeks)
3. **Collect** comments and questions
4. **Respond** to all feedback
5. **Revise** Design Doc based on feedback
6. **Approve** with reviewer sign-off

### Reviewer Selection

| Reviewer Type | Why Include |
|---------------|-------------|
| Tech Lead | Architectural alignment |
| Security | Security review |
| SRE/Ops | Operational concerns |
| Adjacent teams | Integration impact |
| Senior engineer | Quality and patterns |

### RFC Etiquette

**Commenters:**
- Be specific and constructive
- Distinguish blocking vs. non-blocking
- Suggest alternatives, not just problems

**Authors:**
- Respond to all comments
- Mark resolved vs. acknowledged
- Update doc with changes

## Phase 5: Requirements (PRD + UXD)

**Goal**: Refine requirements based on Design Doc decisions.

### Workflow

```bash
visionspec create prd -p <project>
visionspec create uxd -p <project>
```

PRD should:

- [ ] Reference OKRs and key results
- [ ] Align with Design Doc decisions
- [ ] Include measurable acceptance criteria
- [ ] Define experiment parameters

UXD should:

- [ ] Support A/B test variations
- [ ] Include metrics instrumentation
- [ ] Plan for gradual rollout UX

## Phase 6: Experiment Design

**Goal**: Plan how to validate the feature with data.

### Experiment Framework

```markdown
## Experiment: [Name]

### Hypothesis
[What we believe will happen]

### Metrics
**Primary:** [The ONE metric we're optimizing]
**Secondary:** [2-3 supporting metrics]
**Guardrails:** [Metrics that must NOT degrade]

### Design
- **Type:** A/B test | Holdback | Staged rollout
- **Variants:** Control, Treatment(s)
- **Population:** [Who is included/excluded]
- **Duration:** [Statistical power calculation]
- **Sample size:** [Required for significance]

### Launch Criteria
**Green:** Primary metric +[X]%, no guardrail regression
**Yellow:** Primary metric +[Y]%, minor guardrail regression
**Red:** Any guardrail regression > [Z]%

### Rollout Plan
- Week 1: 1% traffic
- Week 2: 10% traffic (if green)
- Week 3: 50% traffic (if green)
- Week 4: 100% (if green)
```

### Statistical Rigor

| Concept | Standard |
|---------|----------|
| Significance level (α) | 0.05 (95% confidence) |
| Power (1-β) | 0.80 (80% chance to detect effect) |
| Minimum detectable effect | Business-relevant change size |
| Duration | Until statistical significance |

## Phase 7: Technical Specs

**Goal**: Detailed implementation specs from Design Doc.

### Workflow

```bash
visionspec context gather -p <project>
visionspec synthesize trd -p <project>
visionspec synthesize tpd -p <project>
visionspec synthesize ird -p <project>
```

TRD should elaborate on Design Doc:

- [ ] Detailed component design
- [ ] API specifications (Protocol Buffers, gRPC)
- [ ] Data pipeline definitions
- [ ] Feature flag configuration

TPD should include:

- [ ] Unit test requirements
- [ ] Integration test plan
- [ ] A/B test validation
- [ ] Performance benchmarks

## Phase 8: Build and Launch

**Goal**: Execute with experiment rigor.

### Workflow

```bash
visionspec reconcile -p <project>
visionspec export speckit -p <project>
```

### Launch Checklist

- [ ] Feature flags configured
- [ ] Monitoring dashboards ready
- [ ] Experiment tracking enabled
- [ ] Rollback procedure documented
- [ ] On-call coverage confirmed
- [ ] Launch review approved

## Google Framework Gates

| Gate | Criteria |
|------|----------|
| OKRs defined | Clear objectives, measurable key results |
| Design Doc approved | RFC complete, reviewers signed off |
| Experiment designed | Metrics, power analysis, launch criteria |
| Launch criteria met | Statistical significance, no guardrail regression |

## See Also

- [How Google Works](https://www.howgoogleworks.net/) - Eric Schmidt
- [Google Design Docs](https://www.industrialempathy.com/posts/design-docs-at-google/)
- [Measure What Matters](https://www.whatmatters.com/) - John Doerr (OKRs)
