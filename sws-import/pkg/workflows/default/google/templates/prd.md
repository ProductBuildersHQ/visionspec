# {project_name} - Product Requirements + RFC

**Author:** {author}
**Reviewers:** [Engineering, PM, UX, Legal, etc.]
**RFC Status:** Draft | Open for Comment | Final Comment | Approved | Withdrawn
**Comment Period Ends:** [Date]
**Last Updated:** {date}

---

<!-- RFC PURPOSE:
     An RFC (Request for Comments) solicits feedback before committing to
     a direction. It's a peer review process that ensures diverse perspectives
     are heard and incorporated.

     Key principle: Seek dissent early. It's cheaper to change a doc than code. -->

## RFC Summary

### One-Liner

[Single sentence describing the proposal]

### Abstract

<!-- 1 paragraph summary of what this RFC proposes -->

[Concise summary of the proposal, why it's needed, and what it changes]

### Comment Solicitation

<!-- What kind of feedback are you specifically seeking? -->

**Feedback requested on:**
- [ ] Technical approach
- [ ] Product direction
- [ ] User experience
- [ ] Privacy/Security implications
- [ ] Resource requirements
- [ ] Timeline feasibility
- [ ] [Other specific questions]

## Motivation

### Problem Statement

<!-- What problem does this solve? Why does it matter? -->

[Clear description of the problem and its impact]

### User Impact

<!-- Who is affected and how? -->

| User Segment | Current Pain | Proposed Improvement |
|--------------|--------------|---------------------|
| [Segment 1] | [Pain point] | [How this helps] |
| [Segment 2] | | |

### Business Impact

<!-- Why does this matter to the business? -->

- **OKR Alignment:** [Which OKRs does this support?]
- **Success Metrics:** [How will we measure success?]

## Proposal

### Goals

<!-- What are we trying to achieve? -->

| ID | Goal | Success Metric | Target |
|----|------|----------------|--------|
| G1 | [Goal] | [Metric] | [Target] |
| G2 | | | |
| G3 | | | |

### Non-Goals

<!-- What are we explicitly NOT doing? -->

| ID | Non-Goal | Rationale |
|----|----------|-----------|
| NG1 | [Non-goal] | [Why out of scope] |
| NG2 | | |

### User Stories

<!-- How will users interact with this? -->

**Story 1:**
As a [user type], I want to [action], so that [outcome].
- Acceptance: [Criteria]

**Story 2:**
As a [user type], I want to [action], so that [outcome].
- Acceptance: [Criteria]

### Functional Requirements

| ID | Requirement | Priority | Rationale |
|----|-------------|----------|-----------|
| FR1 | [Requirement] | P0/P1/P2 | [Why needed] |
| FR2 | | | |

### Non-Functional Requirements

| Category | Requirement | Target |
|----------|-------------|--------|
| Performance | [Latency, throughput] | [Specific target] |
| Scalability | [Scale requirements] | |
| Reliability | [Availability, durability] | |
| Security | [Security requirements] | |

## Alternatives Considered

<!-- What other approaches did you consider? -->

### Alternative 1: [Name]

[Description]

| Pros | Cons |
|------|------|
| [Pro 1] | [Con 1] |

**Disposition:** Rejected because [reason]

### Alternative 2: [Name]

[Description]

**Disposition:** Rejected because [reason]

### Do Nothing

**Why unacceptable:** [Reason status quo doesn't work]

## Tradeoffs and Risks

### Key Tradeoffs

| Decision | We Get | We Give Up | Acceptable Because |
|----------|--------|------------|-------------------|
| [Decision] | [Benefit] | [Cost] | [Rationale] |

### Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| [Risk 1] | H/M/L | H/M/L | [How addressed] |
| [Risk 2] | | | |

### Dependencies

| Dependency | Owner | Status | Risk if Delayed |
|------------|-------|--------|-----------------|
| [Dependency] | [Team] | [Status] | [Impact] |

## Implementation

### High-Level Approach

[Brief description of technical approach - details in Design Doc]

### Phases

| Phase | Scope | Estimated Effort | Dependencies |
|-------|-------|------------------|--------------|
| Phase 1 | [What's included] | [T-shirt size] | [What's needed] |
| Phase 2 | | | |

### Success Metrics

<!-- How will we know this worked? -->

| Metric | Baseline | Target | Measurement Method |
|--------|----------|--------|-------------------|
| [Metric 1] | [Current] | [Goal] | [How measured] |
| [Metric 2] | | | |

### Experiment Plan

<!-- If this requires experimentation/A/B testing -->

**Hypothesis:** [What we believe will happen]
**Experiment:** [How we'll test]
**Success Criteria:** [What would validate hypothesis]
**Duration:** [How long to run]

## Launch Plan

### Launch Criteria

- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

### Rollout Strategy

| Stage | Audience | Duration | Success Criteria |
|-------|----------|----------|------------------|
| Dogfood | Internal | 1 week | [Criteria] |
| Beta | 1% users | 2 weeks | |
| GA | 100% | - | |

### Rollback Plan

[How we roll back if problems arise]

## Open Questions

<!-- Questions for reviewers to address -->

| # | Question | Context | Proposed Answer |
|---|----------|---------|-----------------|
| 1 | [Question] | [Why asking] | [Your suggestion] |
| 2 | | | |

---

## RFC Process

### Timeline

| Date | Milestone |
|------|-----------|
| [Date] | RFC published |
| [Date] | Comment period ends |
| [Date] | Final disposition |
| [Date] | Implementation begins |

### Stakeholders

| Role | Name | Review Status |
|------|------|---------------|
| Author | {author} | - |
| Engineering Lead | | Pending |
| PM Lead | | Pending |
| UX Lead | | Pending |
| [Other] | | |

### Discussion Summary

<!-- Updated as comments come in -->

**Key themes from feedback:**
1. [Theme 1]
2. [Theme 2]

**Changes made based on feedback:**
1. [Change 1]
2. [Change 2]

---

## Review Comments

### Comment Thread 1: [Topic]

**[Reviewer Name]** ([Date]):
> [Comment text]

**[Author]** ([Date]):
> [Response]

**Resolution:** [How resolved]

---

### Comment Thread 2: [Topic]

**[Reviewer Name]** ([Date]):
> [Comment text]

**Resolution:** [How resolved]

---

## Approval

| Reviewer | Role | Decision | Date |
|----------|------|----------|------|
| [Name] | [Role] | Approve/Block/Abstain | [Date] |

**Final Status:** [Approved / Rejected / Withdrawn]
**Rationale:** [If rejected or withdrawn, why]

---

*This document follows the Google RFC format. It should complete the comment period before implementation.*
