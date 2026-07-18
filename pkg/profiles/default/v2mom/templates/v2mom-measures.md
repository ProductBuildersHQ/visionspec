# Measures

<!--
V2MOM Measures Template
Measures are quantifiable metrics that track progress on methods.
Each measure should be SMART (Specific, Measurable, Achievable, Relevant, Time-bound).
-->

## Overview

**V2MOM Level**: {{ .Level }}
**Fiscal Period**: {{ .FiscalPeriod }}
**Owner**: {{ .Owner }}

## Measures Summary

| # | Measure | Method | Baseline | Target | Stretch | Current | Status |
|---|---------|--------|----------|--------|---------|---------|--------|
| 1 | {{ .Measure1Name }} | {{ .Measure1Method }} | {{ .Measure1Baseline }} | {{ .Measure1Target }} | {{ .Measure1Stretch }} | {{ .Measure1Current }} | {{ .Measure1Status }} |
| 2 | {{ .Measure2Name }} | {{ .Measure2Method }} | {{ .Measure2Baseline }} | {{ .Measure2Target }} | {{ .Measure2Stretch }} | {{ .Measure2Current }} | {{ .Measure2Status }} |
| 3 | {{ .Measure3Name }} | {{ .Measure3Method }} | {{ .Measure3Baseline }} | {{ .Measure3Target }} | {{ .Measure3Stretch }} | {{ .Measure3Current }} | {{ .Measure3Status }} |

## Measure Status Legend

- 🟢 **On Track**: Current ≥ Expected at this point
- 🟡 **At Risk**: Current < Expected but recoverable
- 🔴 **Off Track**: Significant gap, intervention needed
- ⚪ **Not Started**: Measurement period hasn't begun

## Measure Details

### MS1: {{ .Measure1Name }}

**For Method**: {{ .Measure1Method }}
**Owner**: {{ .Measure1Owner }}
**Frequency**: {{ .Measure1Frequency }} <!-- daily, weekly, monthly, quarterly -->

#### Definition
<!-- What exactly are we measuring? Be precise -->

#### Goals

| Level | Value | Description |
|-------|-------|-------------|
| Baseline | {{ .Measure1Baseline }} | Starting point / current state |
| Target | {{ .Measure1Target }} | Expected achievement |
| Stretch | {{ .Measure1Stretch }} | Aspirational goal |

#### Measurement Method
<!-- How is this measured? Data source, calculation method -->

**Data Source**: {{ .Measure1DataSource }}
**Calculation**: {{ .Measure1Calculation }}
**Reporting Cadence**: {{ .Measure1Frequency }}

#### Progress Tracking

| Period | Target | Actual | Status | Notes |
|--------|--------|--------|--------|-------|
| {{ .Measure1Period1 }} | {{ .Measure1PeriodTarget1 }} | {{ .Measure1PeriodActual1 }} | {{ .Measure1PeriodStatus1 }} | {{ .Measure1PeriodNotes1 }} |
| {{ .Measure1Period2 }} | {{ .Measure1PeriodTarget2 }} | {{ .Measure1PeriodActual2 }} | {{ .Measure1PeriodStatus2 }} | {{ .Measure1PeriodNotes2 }} |

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| KPI | {{ .Measure1KPIID }} | Tracks |
| Outcome | {{ .Measure1OutcomeID }} | Measures |

#### Parent Measure Contribution
<!-- For department/team levels: How does this roll up to parent measures? -->

**Contributes To**: {{ .Measure1ParentMeasure }}
**Contribution Weight**: {{ .Measure1ContributionWeight }}

---

### MS2: {{ .Measure2Name }}

**For Method**: {{ .Measure2Method }}
**Owner**: {{ .Measure2Owner }}
**Frequency**: {{ .Measure2Frequency }}

#### Definition

#### Goals

| Level | Value | Description |
|-------|-------|-------------|
| Baseline | | |
| Target | | |
| Stretch | | |

#### Measurement Method

**Data Source**:
**Calculation**:
**Reporting Cadence**:

#### Progress Tracking

| Period | Target | Actual | Status | Notes |
|--------|--------|--------|--------|-------|
| | | | | |

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| KPI | | |

#### Parent Measure Contribution

**Contributes To**:
**Contribution Weight**:

---

### MS3: {{ .Measure3Name }}

**For Method**: {{ .Measure3Method }}
**Owner**: {{ .Measure3Owner }}
**Frequency**: {{ .Measure3Frequency }}

#### Definition

#### Goals

| Level | Value | Description |
|-------|-------|-------------|
| Baseline | | |
| Target | | |
| Stretch | | |

#### Measurement Method

**Data Source**:
**Calculation**:
**Reporting Cadence**:

#### Progress Tracking

| Period | Target | Actual | Status | Notes |
|--------|--------|--------|--------|-------|
| | | | | |

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| KPI | | |

#### Parent Measure Contribution

**Contributes To**:
**Contribution Weight**:

---

## Dashboard

```
Method Achievement Progress
─────────────────────────────────────────────────

M1: {{ .Method1Name }}
  MS1: {{ .Measure1Name }}
  [████████░░░░░░░░░░░░] 40% → Target: {{ .Measure1Target }}

M2: {{ .Method2Name }}
  MS2: {{ .Measure2Name }}
  [██████████████░░░░░░] 70% → Target: {{ .Measure2Target }}

M3: {{ .Method3Name }}
  MS3: {{ .Measure3Name }}
  [████░░░░░░░░░░░░░░░░] 20% → Target: {{ .Measure3Target }}

─────────────────────────────────────────────────
Overall V2MOM Progress: 43%
```

---

## Metadata

```yaml
v2mom:
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  measures:
    - id: ms1
      name: "{{ .Measure1Name }}"
      method: "{{ .Measure1Method }}"
      baseline: {{ .Measure1Baseline }}
      target: {{ .Measure1Target }}
      stretch: {{ .Measure1Stretch }}
      frequency: "{{ .Measure1Frequency }}"
      productcontext:
        kpis: []
        outcomes: []
      parent_measure: "{{ .Measure1ParentMeasure }}"
```
