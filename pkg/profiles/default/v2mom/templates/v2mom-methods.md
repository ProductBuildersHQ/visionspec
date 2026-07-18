# Methods

<!--
V2MOM Methods Template
Methods are the specific initiatives and actions to achieve the vision.
Each method should be concrete, measurable, and linked to ProductContext entities.
-->

## Overview

**V2MOM Level**: {{ .Level }}
**Fiscal Period**: {{ .FiscalPeriod }}
**Owner**: {{ .Owner }}

## Methods Summary

| # | Method | Owner | Capabilities | Target Date |
|---|--------|-------|--------------|-------------|
| 1 | {{ .Method1Name }} | {{ .Method1Owner }} | {{ .Method1Capabilities }} | {{ .Method1Target }} |
| 2 | {{ .Method2Name }} | {{ .Method2Owner }} | {{ .Method2Capabilities }} | {{ .Method2Target }} |
| 3 | {{ .Method3Name }} | {{ .Method3Owner }} | {{ .Method3Capabilities }} | {{ .Method3Target }} |

## Method Details

### M1: {{ .Method1Name }}

**Owner**: {{ .Method1Owner }}
**Priority**: {{ .Method1Priority }} <!-- P0, P1, P2 -->
**Status**: {{ .Method1Status }} <!-- not_started, in_progress, at_risk, completed -->

#### Description
<!-- What will we do? Be specific and actionable -->

#### Success Criteria
<!-- How will we know this method succeeded? -->

- [ ] {{ .Method1Criterion1 }}
- [ ] {{ .Method1Criterion2 }}
- [ ] {{ .Method1Criterion3 }}

#### Value Alignment
<!-- Which values does this method embody? -->

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Capability | {{ .Method1CapabilityID }} | Implements |
| Project | {{ .Method1ProjectID }} | Delivers |
| Customer Request | {{ .Method1RequestID }} | Addresses |

#### Parent Method Alignment
<!-- For department/team levels: Which parent method does this support? -->

**Supports**: {{ .Method1ParentMethod }}
**Contribution**: {{ .Method1Contribution }}

---

### M2: {{ .Method2Name }}

**Owner**: {{ .Method2Owner }}
**Priority**: {{ .Method2Priority }}
**Status**: {{ .Method2Status }}

#### Description

#### Success Criteria

- [ ]
- [ ]
- [ ]

#### Value Alignment

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Capability | | |
| Project | | |

#### Parent Method Alignment

**Supports**:
**Contribution**:

---

### M3: {{ .Method3Name }}

**Owner**: {{ .Method3Owner }}
**Priority**: {{ .Method3Priority }}
**Status**: {{ .Method3Status }}

#### Description

#### Success Criteria

- [ ]
- [ ]
- [ ]

#### Value Alignment

#### ProductContext Links

| Link Type | ID | Relationship |
|-----------|-----|-------------|
| Capability | | |
| Project | | |

#### Parent Method Alignment

**Supports**:
**Contribution**:

---

## Dependencies

### Cross-Team Dependencies

| Method | Depends On | Team | Status |
|--------|-----------|------|--------|
| {{ .Method1Name }} | {{ .Dependency1 }} | {{ .DependencyTeam1 }} | {{ .DependencyStatus1 }} |

### Sequencing

```mermaid
gantt
    title Methods Timeline
    dateFormat YYYY-MM-DD
    section Methods
    {{ .Method1Name }} :m1, {{ .Method1Start }}, {{ .Method1Duration }}
    {{ .Method2Name }} :m2, after m1, {{ .Method2Duration }}
    {{ .Method3Name }} :m3, {{ .Method3Start }}, {{ .Method3Duration }}
```

---

## Metadata

```yaml
v2mom:
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  methods:
    - id: m1
      name: "{{ .Method1Name }}"
      owner: "{{ .Method1Owner }}"
      priority: {{ .Method1Priority }}
      productcontext:
        capabilities: []
        projects: []
        requests: []
      parent_method: "{{ .Method1ParentMethod }}"
```
