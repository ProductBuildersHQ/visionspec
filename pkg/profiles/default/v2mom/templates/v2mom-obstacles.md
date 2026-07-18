# Obstacles

<!--
V2MOM Obstacles Template
Obstacles are challenges that could prevent success.
Each obstacle should have mitigation strategies and an owner.
-->

## Overview

**V2MOM Level**: {{ .Level }}
**Fiscal Period**: {{ .FiscalPeriod }}
**Owner**: {{ .Owner }}

## Obstacles Summary

| # | Obstacle | Category | Severity | Owner | Status |
|---|----------|----------|----------|-------|--------|
| 1 | {{ .Obstacle1Name }} | {{ .Obstacle1Category }} | {{ .Obstacle1Severity }} | {{ .Obstacle1Owner }} | {{ .Obstacle1Status }} |
| 2 | {{ .Obstacle2Name }} | {{ .Obstacle2Category }} | {{ .Obstacle2Severity }} | {{ .Obstacle2Owner }} | {{ .Obstacle2Status }} |
| 3 | {{ .Obstacle3Name }} | {{ .Obstacle3Category }} | {{ .Obstacle3Severity }} | {{ .Obstacle3Owner }} | {{ .Obstacle3Status }} |

## Obstacle Categories

- **Technical**: Engineering, architecture, infrastructure challenges
- **Organizational**: Process, alignment, communication issues
- **Resource**: People, budget, time constraints
- **Market**: Competition, customer, regulatory challenges
- **External**: Dependencies on third parties, economic factors

## Obstacle Details

### O1: {{ .Obstacle1Name }}

**Category**: {{ .Obstacle1Category }}
**Severity**: {{ .Obstacle1Severity }} <!-- critical, high, medium, low -->
**Owner**: {{ .Obstacle1Owner }}
**Status**: {{ .Obstacle1Status }} <!-- identified, mitigating, resolved, accepted -->

#### Description
<!-- What is the obstacle? Be specific about the challenge -->

#### Impact
<!-- What methods does this obstacle threaten? -->

| Affected Method | Impact Level | Description |
|----------------|--------------|-------------|
| {{ .Obstacle1AffectedMethod1 }} | {{ .Obstacle1Impact1 }} | {{ .Obstacle1ImpactDesc1 }} |

#### Root Cause Analysis
<!-- What's causing this obstacle? -->

#### Mitigation Strategies

| Strategy | Effort | Effectiveness | Owner | Status |
|----------|--------|---------------|-------|--------|
| {{ .Obstacle1Strategy1 }} | {{ .Obstacle1Effort1 }} | {{ .Obstacle1Effectiveness1 }} | {{ .Obstacle1StrategyOwner1 }} | {{ .Obstacle1StrategyStatus1 }} |
| {{ .Obstacle1Strategy2 }} | {{ .Obstacle1Effort2 }} | {{ .Obstacle1Effectiveness2 }} | {{ .Obstacle1StrategyOwner2 }} | {{ .Obstacle1StrategyStatus2 }} |

#### Contingency Plan
<!-- If mitigation fails, what's the fallback? -->

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Decision | {{ .Obstacle1DecisionID }} | Addresses |
| Customer Request | {{ .Obstacle1RequestID }} | Blocks |

---

### O2: {{ .Obstacle2Name }}

**Category**: {{ .Obstacle2Category }}
**Severity**: {{ .Obstacle2Severity }}
**Owner**: {{ .Obstacle2Owner }}
**Status**: {{ .Obstacle2Status }}

#### Description

#### Impact

| Affected Method | Impact Level | Description |
|----------------|--------------|-------------|
| | | |

#### Root Cause Analysis

#### Mitigation Strategies

| Strategy | Effort | Effectiveness | Owner | Status |
|----------|--------|---------------|-------|--------|
| | | | | |

#### Contingency Plan

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Decision | | |

---

### O3: {{ .Obstacle3Name }}

**Category**: {{ .Obstacle3Category }}
**Severity**: {{ .Obstacle3Severity }}
**Owner**: {{ .Obstacle3Owner }}
**Status**: {{ .Obstacle3Status }}

#### Description

#### Impact

| Affected Method | Impact Level | Description |
|----------------|--------------|-------------|
| | | |

#### Root Cause Analysis

#### Mitigation Strategies

| Strategy | Effort | Effectiveness | Owner | Status |
|----------|--------|---------------|-------|--------|
| | | | | |

#### Contingency Plan

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Decision | | |

---

## Risk Matrix

```
Severity →
  ↑ Critical │  O1  │     │     │
    High     │      │  O2 │     │
    Medium   │      │     │  O3 │
    Low      │      │     │     │
             └──────┴─────┴─────┴────
               High   Med   Low
                   ← Likelihood
```

## Escalation Path

For critical/high severity obstacles that cannot be resolved at this level:

1. Document blocker with full context
2. Escalate to parent V2MOM owner
3. Request resources/decisions needed
4. Track resolution timeline

---

## Metadata

```yaml
v2mom:
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  obstacles:
    - id: o1
      name: "{{ .Obstacle1Name }}"
      category: "{{ .Obstacle1Category }}"
      severity: "{{ .Obstacle1Severity }}"
      owner: "{{ .Obstacle1Owner }}"
      affects_methods: []
      productcontext:
        decisions: []
        requests: []
```
