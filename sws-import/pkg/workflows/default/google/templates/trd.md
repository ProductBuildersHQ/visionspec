# {project_name} - Design Doc

**Author:** {author}
**Reviewers:** [List reviewers]
**Status:** Draft | In Review | Approved | Implemented
**Last Updated:** {date}

---

<!-- DESIGN DOC PURPOSE:
     A design doc is a technical specification that describes HOW you plan
     to solve a problem. It forces you to think through the design before
     writing code, and enables peer review of the approach.

     Key principle: Be explicit about tradeoffs. Every decision has costs. -->

## Overview

### TL;DR

<!-- 2-3 sentences summarizing what this doc proposes -->

[Brief summary of the proposal]

### Context

<!-- Why are we doing this? What's the background? -->

[Background context that reviewers need to understand the proposal]

### Goals

<!-- What are we trying to achieve? Be specific and measurable. -->

- **G1:** [Goal 1 - specific and measurable]
- **G2:** [Goal 2 - specific and measurable]
- **G3:** [Goal 3 - specific and measurable]

### Non-Goals

<!-- What are we explicitly NOT trying to achieve?
     This is as important as goals - it bounds the scope and prevents creep. -->

- **NG1:** [Non-goal 1 - what we won't do and why]
- **NG2:** [Non-goal 2 - what we won't do and why]
- **NG3:** [Non-goal 3 - explicitly out of scope]

## Background

<!-- Deeper context for reviewers unfamiliar with the problem space -->

### Current State

[Describe how things work today]

### Problem Statement

[What specific problem does this design solve?]

### Requirements

<!-- What must the solution do? Distinguish hard requirements from nice-to-haves -->

**Hard Requirements:**
- [Must do X]
- [Must support Y]
- [Must not break Z]

**Soft Requirements:**
- [Should do X if possible]
- [Nice to have Y]

## Design

### System Overview

<!-- High-level architecture diagram or description -->

```
[ASCII diagram or description of system architecture]
```

### Detailed Design

#### Component 1: [Name]

[Detailed description of this component]

**Interface:**
```
[API, interface definition, or contract]
```

**Behavior:**
- [How it works]
- [Edge cases]

#### Component 2: [Name]

[Detailed description]

### Data Model

<!-- If applicable, describe the data model -->

```
[Schema, data structures, or ERD]
```

### API Design

<!-- If applicable, describe the API -->

| Endpoint/Method | Description | Request | Response |
|-----------------|-------------|---------|----------|
| | | | |

## Alternatives Considered

<!-- This section is critical. What other approaches did you consider?
     Why did you reject them? This shows due diligence and helps reviewers
     understand the decision space. -->

### Alternative 1: [Name]

**Description:** [How this approach would work]

**Pros:**
- [Advantage 1]
- [Advantage 2]

**Cons:**
- [Disadvantage 1]
- [Disadvantage 2]

**Why rejected:** [Specific reason this wasn't chosen]

### Alternative 2: [Name]

**Description:** [How this approach would work]

**Pros:**
- [Advantage 1]

**Cons:**
- [Disadvantage 1]

**Why rejected:** [Specific reason]

### Alternative 3: Do Nothing

**Description:** Keep the current system as-is

**Pros:**
- No development cost
- No migration risk

**Cons:**
- [Why status quo is unacceptable]

**Why rejected:** [Why we must act]

## Tradeoffs

<!-- Be explicit about the tradeoffs in your chosen design.
     Every design decision has costs. Acknowledge them. -->

| Decision | Benefit | Cost | Why Acceptable |
|----------|---------|------|----------------|
| [Decision 1] | [What we gain] | [What we give up] | [Why this tradeoff is right] |
| [Decision 2] | | | |
| [Decision 3] | | | |

### Key Tradeoff: [Most Important One]

[Detailed discussion of the most significant tradeoff]

## Cross-Cutting Concerns

### Scalability

<!-- How does this design scale? What are the limits? -->

- **Current scale:** [What it needs to handle today]
- **Target scale:** [What it should handle]
- **Bottlenecks:** [Known limitations]
- **Scaling strategy:** [How to scale when needed]

### Reliability

<!-- How does this design handle failures? -->

- **Failure modes:** [What can go wrong]
- **Recovery:** [How system recovers]
- **SLO targets:** [Availability, latency targets]

### Security

<!-- Security considerations -->

- **Threat model:** [What threats are we defending against]
- **Mitigations:** [How we address them]
- **Data handling:** [Sensitive data considerations]

### Privacy

<!-- Privacy considerations, especially for user data -->

- **Data collected:** [What data]
- **Retention:** [How long kept]
- **Access controls:** [Who can access]

### Observability

<!-- How will we know if this is working? -->

- **Metrics:** [Key metrics to track]
- **Logging:** [What to log]
- **Alerting:** [When to alert]

## Implementation Plan

### Phases

<!-- Break implementation into phases if appropriate -->

**Phase 1: [Name]**
- [Deliverable 1]
- [Deliverable 2]
- **Milestone:** [What marks completion]

**Phase 2: [Name]**
- [Deliverable 1]
- **Milestone:** [What marks completion]

### Migration Strategy

<!-- If changing existing system, how do we migrate? -->

- **Approach:** [Big bang, gradual, shadow mode, etc.]
- **Rollback plan:** [How to undo if problems]
- **Feature flags:** [What flags control rollout]

### Testing Strategy

- **Unit tests:** [Coverage expectations]
- **Integration tests:** [What to test]
- **Load tests:** [Performance validation]
- **Canary/Shadow:** [Pre-production validation]

## Launch Plan

### Launch Criteria

<!-- What must be true before we launch? -->

- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

### Rollout Plan

| Stage | % Traffic | Duration | Success Criteria |
|-------|-----------|----------|------------------|
| Canary | 1% | 1 day | [Metrics within bounds] |
| Limited | 10% | 3 days | |
| Full | 100% | - | |

### Rollback Triggers

<!-- When do we roll back? -->

- [Metric 1] exceeds [threshold]
- [Error rate] above [X%]
- [Customer-reported issue type]

## Open Questions

<!-- What haven't you figured out yet? Be honest. -->

| Question | Owner | Due Date | Status |
|----------|-------|----------|--------|
| [Question 1] | [Who] | [When] | Open/Resolved |
| [Question 2] | | | |

## Future Work

<!-- What's explicitly deferred to later? -->

- [Future enhancement 1]
- [Future enhancement 2]

---

## Appendix

### A. Glossary

| Term | Definition |
|------|------------|
| | |

### B. References

- [Link to related doc 1]
- [Link to related doc 2]

### C. Revision History

| Date | Author | Changes |
|------|--------|---------|
| {date} | {author} | Initial draft |

---

## Review Comments

<!-- Reviewers add comments here or use doc commenting -->

### Reviewer: [Name]
**Date:** [Date]
**Status:** Approved / Needs Changes / Blocking

**Comments:**
- [Comment 1]
- [Comment 2]

---

*This document follows the Google Design Doc format. It should be reviewed and approved before implementation begins.*
