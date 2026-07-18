# V2MOM Cascade Alignment

<!--
V2MOM Alignment Template
Shows how this V2MOM aligns to and supports parent V2MOM methods.
This document visualizes the cascade relationship.
-->

## Overview

**This V2MOM**: {{ .ThisV2MOM }}
**Level**: {{ .Level }}
**Fiscal Period**: {{ .FiscalPeriod }}
**Owner**: {{ .Owner }}

**Parent V2MOM**: {{ .ParentV2MOM }}
**Parent Level**: {{ .ParentLevel }}
**Parent Owner**: {{ .ParentOwner }}

---

## Cascade Visualization

```
┌─────────────────────────────────────────────────────────────────┐
│  {{ .ParentLevel | upper }} V2MOM: {{ .ParentV2MOM }}           │
│  Vision: "{{ .ParentVision | truncate 50 }}"                    │
├─────────────────────────────────────────────────────────────────┤
│  Methods:                                                        │
│  ├── {{ .ParentMethod1 }}                                       │
│  │   └── ✓ Supported by this V2MOM                              │
│  ├── {{ .ParentMethod2 }}                                       │
│  │   └── ✓ Supported by this V2MOM                              │
│  └── {{ .ParentMethod3 }}                                       │
│      └── ○ Not directly supported                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ CASCADE
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  {{ .Level | upper }} V2MOM: {{ .ThisV2MOM }}                   │
│  Vision: "{{ .ThisVision | truncate 50 }}"                      │
├─────────────────────────────────────────────────────────────────┤
│  Methods:                                                        │
│  ├── {{ .Method1 }} → supports {{ .Method1Parent }}             │
│  ├── {{ .Method2 }} → supports {{ .Method2Parent }}             │
│  └── {{ .Method3 }} → supports {{ .Method3Parent }}             │
└─────────────────────────────────────────────────────────────────┘
```

---

## Alignment Matrix

### Vision Alignment

| This Vision | Parent Methods Supported | Alignment Score |
|-------------|-------------------------|-----------------|
| {{ .ThisVision }} | {{ .SupportedParentMethods }} | {{ .VisionAlignmentScore }}/10 |

**Alignment Rationale**:
<!-- Explain how this vision advances the parent's methods -->

### Method-to-Method Mapping

| This Method | Supports Parent Method | Contribution | Coverage |
|-------------|----------------------|--------------|----------|
| {{ .Method1 }} | {{ .Method1Parent }} | {{ .Method1Contribution }} | {{ .Method1Coverage }}% |
| {{ .Method2 }} | {{ .Method2Parent }} | {{ .Method2Contribution }} | {{ .Method2Coverage }}% |
| {{ .Method3 }} | {{ .Method3Parent }} | {{ .Method3Contribution }} | {{ .Method3Coverage }}% |

### Parent Method Coverage

| Parent Method | Covered By | Total Coverage | Gap Analysis |
|--------------|-----------|----------------|--------------|
| {{ .ParentMethod1 }} | {{ .ParentMethod1CoveredBy }} | {{ .ParentMethod1Coverage }}% | {{ .ParentMethod1Gap }} |
| {{ .ParentMethod2 }} | {{ .ParentMethod2CoveredBy }} | {{ .ParentMethod2Coverage }}% | {{ .ParentMethod2Gap }} |
| {{ .ParentMethod3 }} | {{ .ParentMethod3CoveredBy }} | {{ .ParentMethod3Coverage }}% | {{ .ParentMethod3Gap }} |

### Coverage Gaps

<!-- Parent methods not adequately covered by this V2MOM -->

| Gap | Parent Method | Current Coverage | Recommendation |
|-----|--------------|------------------|----------------|
| {{ .Gap1 }} | {{ .Gap1ParentMethod }} | {{ .Gap1Coverage }}% | {{ .Gap1Recommendation }} |

---

## Measure Roll-up

### How Measures Contribute to Parent

| This Measure | Parent Measure | Roll-up Type | Weight |
|--------------|---------------|--------------|--------|
| {{ .Measure1 }} | {{ .Measure1Parent }} | {{ .Measure1RollupType }} | {{ .Measure1Weight }}% |
| {{ .Measure2 }} | {{ .Measure2Parent }} | {{ .Measure2RollupType }} | {{ .Measure2Weight }}% |

**Roll-up Types**:
- **Sum**: Values add together (e.g., revenue from multiple teams)
- **Average**: Values average (e.g., customer satisfaction across regions)
- **Min**: Lowest value determines parent (e.g., SLA compliance)
- **Milestone**: Binary contribution to parent milestone

---

## Sibling V2MOM Dependencies

### Peer V2MOMs at Same Level

| Sibling V2MOM | Owner | Shared Parent Methods | Dependencies |
|--------------|-------|----------------------|--------------|
| {{ .Sibling1 }} | {{ .Sibling1Owner }} | {{ .Sibling1SharedMethods }} | {{ .Sibling1Dependencies }} |
| {{ .Sibling2 }} | {{ .Sibling2Owner }} | {{ .Sibling2SharedMethods }} | {{ .Sibling2Dependencies }} |

### Cross-V2MOM Dependencies

| My Method | Depends On | Sibling V2MOM | Status |
|-----------|-----------|---------------|--------|
| {{ .Method1 }} | {{ .Method1SiblingDep }} | {{ .Method1SiblingV2MOM }} | {{ .Method1SiblingDepStatus }} |

---

## Alignment Health Score

```
Overall Alignment Score: {{ .OverallAlignmentScore }}/100

Breakdown:
├── Vision-to-Methods Alignment:  {{ .VisionMethodsScore }}/30
├── Method Coverage:              {{ .MethodCoverageScore }}/30
├── Measure Roll-up:              {{ .MeasureRollupScore }}/20
└── Dependency Management:        {{ .DependencyScore }}/20

{{ if gt .OverallAlignmentScore 80 }}
✅ STRONG ALIGNMENT - This V2MOM effectively supports parent objectives
{{ else if gt .OverallAlignmentScore 60 }}
🟡 MODERATE ALIGNMENT - Some gaps need attention
{{ else }}
🔴 WEAK ALIGNMENT - Significant alignment issues require resolution
{{ end }}
```

---

## Visualization for VisionStudio

<!-- This section provides data for VisionStudio cascade visualization -->

```mermaid
graph TD
    subgraph "{{ .ParentLevel }} Level"
        PV[Vision: {{ .ParentVision | truncate 30 }}]
        PM1[M1: {{ .ParentMethod1 | truncate 20 }}]
        PM2[M2: {{ .ParentMethod2 | truncate 20 }}]
        PM3[M3: {{ .ParentMethod3 | truncate 20 }}]
        PV --> PM1
        PV --> PM2
        PV --> PM3
    end

    subgraph "{{ .Level }} Level"
        TV[Vision: {{ .ThisVision | truncate 30 }}]
        TM1[M1: {{ .Method1 | truncate 20 }}]
        TM2[M2: {{ .Method2 | truncate 20 }}]
        TM3[M3: {{ .Method3 | truncate 20 }}]
        TV --> TM1
        TV --> TM2
        TV --> TM3
    end

    PM1 -.-> TV
    PM2 -.-> TV
    TM1 --> PM1
    TM2 --> PM1
    TM3 --> PM2

    style PM1 fill:#90EE90
    style PM2 fill:#90EE90
    style PM3 fill:#FFB6C1
```

---

## Metadata

```yaml
alignment:
  this_v2mom: {{ .ThisV2MOM }}
  level: {{ .Level }}
  parent_v2mom: {{ .ParentV2MOM }}
  parent_level: {{ .ParentLevel }}

  method_mapping:
    - this_method: "{{ .Method1 }}"
      parent_method: "{{ .Method1Parent }}"
      contribution: "{{ .Method1Contribution }}"
    - this_method: "{{ .Method2 }}"
      parent_method: "{{ .Method2Parent }}"
      contribution: "{{ .Method2Contribution }}"

  scores:
    overall: {{ .OverallAlignmentScore }}
    vision_methods: {{ .VisionMethodsScore }}
    method_coverage: {{ .MethodCoverageScore }}
    measure_rollup: {{ .MeasureRollupScore }}
    dependencies: {{ .DependencyScore }}
```
