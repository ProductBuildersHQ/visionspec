# {project_name} - Experiment Plan

**Author:** {author}
**Experiment ID:** [EXP-XXXX]
**Status:** Proposed | Approved | Running | Analyzing | Complete
**Last Updated:** {date}

---

<!-- EXPERIMENT PURPOSE:
     Data-driven decision making requires structured experimentation.
     This template ensures experiments are well-designed, statistically
     valid, and produce actionable insights.

     Key principle: Define success criteria BEFORE running the experiment. -->

## Experiment Overview

### Hypothesis

<!-- State your hypothesis clearly and specifically.
     Format: "If we [change], then [metric] will [direction] by [amount]
     because [rationale]" -->

**If** we [specific change/treatment],
**then** [metric] will [increase/decrease] by [expected magnitude],
**because** [causal reasoning].

### One-Liner

[Single sentence describing the experiment]

### Success Criteria

<!-- Define success BEFORE running. This prevents moving goalposts. -->

| Metric | Baseline | Minimum Success | Target | Stretch |
|--------|----------|-----------------|--------|---------|
| Primary: [Metric] | [Current value] | [Minimum to ship] | [Goal] | [Ideal] |
| Secondary: [Metric] | | | | |
| Guardrail: [Metric] | | Must not degrade | | |

### Decision Framework

<!-- What will you do based on results? -->

| Outcome | Action |
|---------|--------|
| Primary metric hits target | Ship to 100% |
| Primary metric hits minimum | Ship with iteration plan |
| Primary metric misses minimum | Do not ship; iterate or abandon |
| Guardrail metric degrades | Do not ship regardless of primary |

## Background

### Context

[Why are we running this experiment? What do we hope to learn?]

### Previous Experiments

<!-- What have we tried before? What did we learn? -->

| Experiment | Hypothesis | Result | Learning |
|------------|------------|--------|----------|
| [EXP-XXX] | [Hypothesis] | [Win/Loss/Neutral] | [What we learned] |

### Related Work

[Links to relevant docs, research, or industry examples]

## Experiment Design

### Treatment Description

**Control (A):** [Current experience - describe specifically]

**Treatment (B):** [New experience - describe specifically]

<!-- If multivariate, add more treatments -->

**Treatment (C):** [Alternative treatment if testing multiple]

### Audience

**Target Population:** [Who is eligible for the experiment]

**Exclusions:**
- [Who is excluded and why]
- [E.g., new users, specific markets, etc.]

**Sample Size:**
- **Required:** [Statistical calculation - see power analysis]
- **Expected duration:** [How long to reach sample size]

### Randomization

**Unit of randomization:** [User, session, device, etc.]
**Allocation:** [50/50, 90/10, etc. and rationale]
**Stratification:** [If stratifying, by what dimensions]

### Power Analysis

<!-- Ensure experiment is properly powered -->

| Parameter | Value |
|-----------|-------|
| Baseline conversion/metric | [Current rate] |
| Minimum detectable effect (MDE) | [Smallest change worth detecting] |
| Statistical significance (α) | [Usually 0.05] |
| Statistical power (1-β) | [Usually 0.80] |
| Required sample size per arm | [Calculated] |
| Estimated runtime | [Days/weeks] |

## Metrics

### Primary Metric

**Metric:** [Name]
**Definition:** [Exactly how it's calculated]
**Why primary:** [Why this is the best measure of success]
**Expected direction:** [Increase/Decrease]
**Expected magnitude:** [Based on hypothesis]

### Secondary Metrics

| Metric | Definition | Expected Change | Why Included |
|--------|------------|-----------------|--------------|
| [Metric 1] | [Calculation] | [Direction] | [Rationale] |
| [Metric 2] | | | |

### Guardrail Metrics

<!-- Metrics that must NOT degrade -->

| Metric | Definition | Acceptable Range | Why Guardrail |
|--------|------------|------------------|---------------|
| [Metric 1] | [Calculation] | [No more than X% degradation] | [What harm it prevents] |
| [Metric 2] | | | |

### Counter Metrics

<!-- Metrics that help detect unintended consequences -->

| Metric | What It Catches |
|--------|-----------------|
| [Metric] | [Unintended consequence it would reveal] |

## Implementation

### Technical Requirements

- [ ] [Requirement 1]
- [ ] [Requirement 2]
- [ ] [Logging/tracking requirements]

### Feature Flags

| Flag Name | Description | Default |
|-----------|-------------|---------|
| [flag_name] | [What it controls] | [Off/On] |

### Rollout Plan

| Stage | % Allocation | Duration | Gate Criteria |
|-------|--------------|----------|---------------|
| Ramp 1 | 1% | 1 day | No crashes, logging works |
| Ramp 2 | 10% | 3 days | Metrics tracking correctly |
| Full | 50% | Until significant | Power achieved |

### Monitoring

**Dashboard:** [Link to experiment dashboard]
**Alerts:** [What alerts are set up]
**Check-in cadence:** [Daily/Weekly review]

## Analysis Plan

### Pre-Analysis

<!-- Planned before experiment starts -->

- [ ] Verify randomization is balanced
- [ ] Confirm metrics are logging correctly
- [ ] Check for novelty effects (if applicable)
- [ ] Document any anomalies during experiment period

### Statistical Methods

**Primary analysis:** [T-test, chi-square, regression, etc.]
**Corrections:** [Bonferroni, FDR, etc. if multiple comparisons]
**Confidence level:** [95%, etc.]

### Segmentation Analysis

<!-- What segments will you analyze? Define in advance. -->

| Segment | Hypothesis |
|---------|------------|
| [New vs returning users] | [Expected difference] |
| [Mobile vs desktop] | |
| [Market/region] | |

### Interpretation Guidelines

| Scenario | Interpretation | Action |
|----------|----------------|--------|
| Statistically significant positive | Treatment works | Consider shipping |
| Statistically significant negative | Treatment hurts | Do not ship |
| Not statistically significant | Can't conclude | Need more data or abandon |

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| [Risk 1] | H/M/L | H/M/L | [Plan] |
| [Risk 2] | | | |

### Stopping Criteria

<!-- When should we stop early? -->

- [ ] Guardrail metric degrades by more than [X%]
- [ ] [Specific harm indicator]
- [ ] Reached statistical significance early (with proper correction)

## Timeline

| Milestone | Date | Owner |
|-----------|------|-------|
| Experiment design approved | | |
| Implementation complete | | |
| Experiment start | | |
| Expected power reached | | |
| Analysis complete | | |
| Decision made | | |

## Approvals

| Role | Name | Status | Date |
|------|------|--------|------|
| Experiment Owner | {author} | - | |
| Engineering Lead | | Pending | |
| Data Science | | Pending | |
| Product Lead | | Pending | |

---

## Results (Post-Experiment)

### Summary

**Experiment ran:** [Start date] to [End date]
**Sample size achieved:** [N per arm]
**Result:** [Win / Loss / Neutral / Inconclusive]

### Primary Metric Results

| Arm | N | Metric Value | 95% CI | vs Control |
|-----|---|--------------|--------|------------|
| Control | | | | - |
| Treatment | | | | [+/-X%] |

**Statistical significance:** [p = X.XX]
**Practical significance:** [Is the effect meaningful?]

### Secondary Metrics

| Metric | Control | Treatment | Change | Significant? |
|--------|---------|-----------|--------|--------------|
| | | | | |

### Guardrail Metrics

| Metric | Control | Treatment | Status |
|--------|---------|-----------|--------|
| | | | [OK / Degraded] |

### Segment Analysis

[Key segment findings]

### Learnings

1. [Learning 1]
2. [Learning 2]
3. [Learning 3]

### Decision

**Decision:** [Ship / Iterate / Abandon]
**Rationale:** [Why this decision]
**Next steps:** [What happens next]

---

*This document follows Google's experiment methodology. Results should be reviewed with statistical rigor.*
