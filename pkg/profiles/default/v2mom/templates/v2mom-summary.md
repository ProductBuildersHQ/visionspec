# V2MOM Summary

<!--
V2MOM Summary Template
A unified view of the complete V2MOM with cross-references and ProductContext links.
-->

## Header

| Attribute | Value |
|-----------|-------|
| **V2MOM ID** | {{ .V2MOMID }} |
| **Name** | {{ .V2MOMName }} |
| **Level** | {{ .Level }} |
| **Fiscal Period** | {{ .FiscalPeriod }} |
| **Owner** | {{ .Owner }} |
| **Status** | {{ .Status }} |
| **Last Updated** | {{ .UpdatedAt }} |
| **Parent V2MOM** | {{ .ParentV2MOM }} |

---

## 🎯 Vision

> {{ .VisionStatement }}

**Time Horizon**: {{ .VisionTimeHorizon }}
**Alignment Score**: {{ .VisionAlignmentScore }}/10

---

## 💎 Values (Ranked)

| Rank | Value | Trade-off Guidance |
|------|-------|-------------------|
| 1 | **{{ .Value1 }}** | {{ .Value1Tradeoff }} |
| 2 | **{{ .Value2 }}** | {{ .Value2Tradeoff }} |
| 3 | **{{ .Value3 }}** | {{ .Value3Tradeoff }} |
| 4 | **{{ .Value4 }}** | {{ .Value4Tradeoff }} |
| 5 | **{{ .Value5 }}** | {{ .Value5Tradeoff }} |

---

## 🚀 Methods

### Method Summary

| # | Method | Owner | Priority | Status | Target | Capabilities |
|---|--------|-------|----------|--------|--------|--------------|
| 1 | {{ .Method1 }} | {{ .Method1Owner }} | {{ .Method1Priority }} | {{ .Method1Status }} | {{ .Method1Target }} | {{ .Method1Capabilities }} |
| 2 | {{ .Method2 }} | {{ .Method2Owner }} | {{ .Method2Priority }} | {{ .Method2Status }} | {{ .Method2Target }} | {{ .Method2Capabilities }} |
| 3 | {{ .Method3 }} | {{ .Method3Owner }} | {{ .Method3Priority }} | {{ .Method3Status }} | {{ .Method3Target }} | {{ .Method3Capabilities }} |

### Method → Parent Alignment

| This Method | Supports Parent Method |
|-------------|----------------------|
| {{ .Method1 }} | {{ .Method1Parent }} |
| {{ .Method2 }} | {{ .Method2Parent }} |
| {{ .Method3 }} | {{ .Method3Parent }} |

---

## 🚧 Obstacles

| # | Obstacle | Category | Severity | Status | Mitigation |
|---|----------|----------|----------|--------|------------|
| 1 | {{ .Obstacle1 }} | {{ .Obstacle1Category }} | {{ .Obstacle1Severity }} | {{ .Obstacle1Status }} | {{ .Obstacle1Mitigation }} |
| 2 | {{ .Obstacle2 }} | {{ .Obstacle2Category }} | {{ .Obstacle2Severity }} | {{ .Obstacle2Status }} | {{ .Obstacle2Mitigation }} |
| 3 | {{ .Obstacle3 }} | {{ .Obstacle3Category }} | {{ .Obstacle3Severity }} | {{ .Obstacle3Status }} | {{ .Obstacle3Mitigation }} |

---

## 📊 Measures

| # | Measure | Method | Baseline | Target | Current | Status |
|---|---------|--------|----------|--------|---------|--------|
| 1 | {{ .Measure1 }} | {{ .Measure1Method }} | {{ .Measure1Baseline }} | {{ .Measure1Target }} | {{ .Measure1Current }} | {{ .Measure1Status }} |
| 2 | {{ .Measure2 }} | {{ .Measure2Method }} | {{ .Measure2Baseline }} | {{ .Measure2Target }} | {{ .Measure2Current }} | {{ .Measure2Status }} |
| 3 | {{ .Measure3 }} | {{ .Measure3Method }} | {{ .Measure3Baseline }} | {{ .Measure3Target }} | {{ .Measure3Current }} | {{ .Measure3Status }} |

### Progress Overview

```
Overall V2MOM Progress: {{ .OverallProgress }}%

Vision Clarity:      [{{ .VisionProgress | progressBar }}] {{ .VisionScore }}/10
Method Execution:    [{{ .MethodProgress | progressBar }}] {{ .MethodsComplete }}/{{ .MethodsTotal }}
Obstacle Resolution: [{{ .ObstacleProgress | progressBar }}] {{ .ObstaclesResolved }}/{{ .ObstaclesTotal }}
Measure Achievement: [{{ .MeasureProgress | progressBar }}] {{ .MeasuresOnTrack }}/{{ .MeasuresTotal }}
```

---

## 🔗 ProductContext Integration

### Linked Capabilities

| Capability | Methods | Status |
|------------|---------|--------|
| {{ .Capability1 }} | {{ .Capability1Methods }} | {{ .Capability1Status }} |
| {{ .Capability2 }} | {{ .Capability2Methods }} | {{ .Capability2Status }} |

### Linked Projects

| Project | Methods | Owner | Status |
|---------|---------|-------|--------|
| {{ .Project1 }} | {{ .Project1Methods }} | {{ .Project1Owner }} | {{ .Project1Status }} |
| {{ .Project2 }} | {{ .Project2Methods }} | {{ .Project2Owner }} | {{ .Project2Status }} |

### Linked KPIs

| KPI | Measure | Current | Target |
|-----|---------|---------|--------|
| {{ .KPI1 }} | {{ .KPI1Measure }} | {{ .KPI1Current }} | {{ .KPI1Target }} |
| {{ .KPI2 }} | {{ .KPI2Measure }} | {{ .KPI2Current }} | {{ .KPI2Target }} |

### Linked Decisions

| Decision | Obstacle | Status |
|----------|----------|--------|
| {{ .Decision1 }} | {{ .Decision1Obstacle }} | {{ .Decision1Status }} |

---

## 📊 Cascade Context

### Full Cascade Path

```
┌────────────────────────────────────────────┐
│ COMPANY: {{ .CompanyV2MOM }}               │
│ Vision: {{ .CompanyVision | truncate 40 }} │
└────────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────┐
│ DEPARTMENT: {{ .DepartmentV2MOM }}         │
│ Vision: {{ .DepartmentVision | truncate 40 }} │
└────────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────┐
│ TEAM: {{ .TeamV2MOM }}                     │
│ Vision: {{ .TeamVision | truncate 40 }}    │  ← YOU ARE HERE
└────────────────────────────────────────────┘
```

### Sibling V2MOMs (Same Parent)

| V2MOM | Owner | Methods | Alignment |
|-------|-------|---------|-----------|
| {{ .Sibling1 }} | {{ .Sibling1Owner }} | {{ .Sibling1Methods }} | {{ .Sibling1Alignment }}% |
| {{ .Sibling2 }} | {{ .Sibling2Owner }} | {{ .Sibling2Methods }} | {{ .Sibling2Alignment }}% |

### Child V2MOMs (Cascade Down)

| V2MOM | Owner | Level | Alignment |
|-------|-------|-------|-----------|
| {{ .Child1 }} | {{ .Child1Owner }} | {{ .Child1Level }} | {{ .Child1Alignment }}% |
| {{ .Child2 }} | {{ .Child2Owner }} | {{ .Child2Level }} | {{ .Child2Alignment }}% |

---

## 📈 Health Dashboard

| Metric | Score | Status |
|--------|-------|--------|
| Vision Clarity | {{ .VisionScore }}/10 | {{ .VisionHealth }} |
| Values Alignment | {{ .ValuesScore }}/10 | {{ .ValuesHealth }} |
| Method Progress | {{ .MethodScore }}/10 | {{ .MethodHealth }} |
| Obstacle Management | {{ .ObstacleScore }}/10 | {{ .ObstacleHealth }} |
| Measure Achievement | {{ .MeasureScore }}/10 | {{ .MeasureHealth }} |
| Cascade Alignment | {{ .AlignmentScore }}/10 | {{ .AlignmentHealth }} |
| **Overall V2MOM Health** | **{{ .OverallScore }}/10** | **{{ .OverallHealth }}** |

---

## Metadata

```yaml
v2mom:
  id: {{ .V2MOMID }}
  name: {{ .V2MOMName }}
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  owner: {{ .Owner }}

  cascade:
    parent: {{ .ParentV2MOM }}
    children: {{ .ChildV2MOMs }}
    siblings: {{ .SiblingV2MOMs }}

  health:
    overall: {{ .OverallScore }}
    vision: {{ .VisionScore }}
    values: {{ .ValuesScore }}
    methods: {{ .MethodScore }}
    obstacles: {{ .ObstacleScore }}
    measures: {{ .MeasureScore }}
    alignment: {{ .AlignmentScore }}

  productcontext:
    capabilities: {{ .LinkedCapabilities }}
    projects: {{ .LinkedProjects }}
    kpis: {{ .LinkedKPIs }}
    decisions: {{ .LinkedDecisions }}
```
