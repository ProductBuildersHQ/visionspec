# {project_name} - Product Requirements

**Date:** {date}
**Author:** {author}
**Status:** Draft

---

<!-- JTBD REMINDER:
     "When we buy a product, we essentially 'hire' something to get a job done.
     If it does the job well, when we are confronted with the same job, we hire
     that same product again."
     - Clayton Christensen -->

## Job Alignment

<!-- Ensure solution is aligned with the job, not just features -->

### The Job We're Solving

**Job Statement:**
When [situation/context],
I want to [motivation/goal],
So I can [expected outcome].

### Why Customers Will Hire This

**Push:** [What makes current situation painful enough to switch?]

**Pull:** [What makes our solution attractive?]

### Why Customers Might Not Hire This

**Anxiety:** [What concerns might prevent adoption?]
**Mitigation:** [How we address those concerns]

**Habit:** [What current behaviors we must overcome?]
**Strategy:** [How we make switching easy]

## Outcome-Based Requirements

<!-- Requirements framed as outcomes, not features
     "Minimize the time it takes to..." not "Add a button that..." -->

### Primary Outcomes to Enable

<!-- Map to underserved outcomes from MRD -->

| Outcome | Target Improvement | How Solution Achieves It |
|---------|-------------------|-------------------------|
| Minimize time to [action] | [Current: X → Target: Y] | [Mechanism] |
| Minimize likelihood of [negative] | [Current: X% → Target: Y%] | [Mechanism] |
| Increase ability to [positive] | [Current: X → Target: Y] | [Mechanism] |

### Outcome Metrics

<!-- How we'll measure if we're achieving these outcomes -->

| Outcome | Metric | Current Baseline | Target | Measurement Method |
|---------|--------|------------------|--------|-------------------|
| | | | | |

## Functional Requirements

<!-- The functional capabilities needed to achieve outcomes -->

### Must Enable (P0)

<!-- Without these, customers can't get the job done -->

**Requirement 1:** [Name]
- Outcome enabled: [Which outcome this serves]
- Acceptance criteria:
  - [ ] [Criterion]
  - [ ] [Criterion]

**Requirement 2:** [Name]
- Outcome enabled: [Which outcome this serves]
- Acceptance criteria:
  - [ ] [Criterion]

### Should Enable (P1)

<!-- Improves how well the job gets done -->

**Requirement 3:** [Name]
- Outcome enabled: [Which outcome this serves]
- Acceptance criteria:
  - [ ] [Criterion]

### Could Enable (P2)

<!-- Nice to have but not essential for job completion -->

- [Requirement]: [Outcome it would enable]

## Emotional Requirements

<!-- Jobs have emotional dimensions—don't ignore them -->

### How Customers Should Feel

| Moment | Desired Emotion | How We Create It |
|--------|-----------------|------------------|
| First use | [e.g., Confident] | [Design/copy/experience approach] |
| During use | [e.g., In control] | |
| After completion | [e.g., Accomplished] | |

### Emotions to Avoid

| Emotion to Avoid | Why It's a Risk | How We Prevent It |
|-----------------|-----------------|-------------------|
| [e.g., Anxious] | [When it might occur] | [Mitigation] |
| [e.g., Stupid] | [When it might occur] | |

## Social Requirements

<!-- How does using this solution affect how customers are perceived? -->

### Social Outcomes

| Desired Perception | How Solution Enables It |
|-------------------|------------------------|
| [e.g., Competent] | [What about the solution supports this?] |
| [e.g., Modern] | |

### Social Risks

| Risk | When It Arises | Mitigation |
|------|---------------|------------|
| [e.g., Looks lazy] | [If using automation] | [Frame as efficient] |

## Job Stories (User Stories)

<!-- Stories in job format for development -->

### Core Job Stories

**Story 1:**
When [situation with emotional/social context],
I want to [action with motivation],
So I can [outcome including feeling/perception].

Acceptance criteria:
- [ ] [Criterion tied to outcome]
- [ ] [Criterion tied to outcome]

**Story 2:**
When [situation with emotional/social context],
I want to [action with motivation],
So I can [outcome including feeling/perception].

### Supporting Job Stories

**Story 3:**
When [situation],
I want to [action],
So I can [outcome].

## Hiring Criteria

<!-- What must be true for customers to "hire" this solution? -->

### Minimum Hiring Bar

<!-- Without these, customers won't hire us -->

- [ ] [Must be true for initial hire]
- [ ] [Must be true for initial hire]
- [ ] [Must be true for initial hire]

### Repeat Hiring Criteria

<!-- What keeps customers hiring us again? -->

- [ ] [Must be true for repeat hire]
- [ ] [Must be true for repeat hire]

### Firing Triggers

<!-- What would cause customers to fire us? -->

- [Trigger]: [How we prevent it]
- [Trigger]: [How we prevent it]

## Switching Experience

<!-- Make it easy to switch FROM current solutions -->

### From [Current Solution 1]

| Switching Barrier | How We Lower It |
|-------------------|-----------------|
| [Data migration] | [Automatic import] |
| [Learning curve] | [Familiar patterns] |
| [Sunk cost feeling] | [Immediate value] |

### From "Doing Nothing"

<!-- Sometimes the hardest switch is from inaction -->

| Inertia Factor | How We Overcome It |
|----------------|-------------------|
| [Good enough] | [Show hidden costs] |
| [Not priority] | [Trigger-based activation] |

## Non-Requirements

<!-- What we're explicitly NOT building -->

### Out of Scope

| Feature/Capability | Why Out of Scope |
|-------------------|------------------|
| [Feature X] | [Not part of core job] |
| [Feature Y] | [Overserves the outcome] |

### Adjacent Jobs We Won't Solve

<!-- Related jobs that we're intentionally not addressing -->

- [Adjacent job]: [Why separate / future consideration]

---

## Success Metrics

### Job Completion Rate

**Primary metric:** [What % of users complete the job?]
**Target:** [X%]
**Measurement:** [How tracked]

### Outcome Achievement

| Outcome | Metric | Target |
|---------|--------|--------|
| | | |

### Hiring Signals

| Signal | Indicates | Target |
|--------|-----------|--------|
| [Repeat usage] | [Job done well] | [X% return rate] |
| [Referral] | [Social job done well] | [X% refer] |

---

*This document follows Jobs-to-be-Done methodology. Requirements should trace to outcomes, not feature requests.*
