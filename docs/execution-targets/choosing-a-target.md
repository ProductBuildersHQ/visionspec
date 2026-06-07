# Choosing a Target

This guide helps you select the right execution target for your project.

## Quick Decision Matrix

| If you need... | Choose |
|----------------|--------|
| Enterprise rigor with approval gates | [AWS AI-DLC](aws-aidlc.md) |
| GitHub-native PR workflows | [SpecKit](speckit.md) |
| Fast iteration with parallel execution | [GSD](gsd.md) |
| Complex DAG-based task orchestration | [GasTown](gastown.md) |
| Role-based multi-agent coordination | [GasCity](gascity.md) |

## Feature Comparison

| Feature | AI-DLC | SpecKit | GSD | GasTown | GasCity |
|---------|--------|---------|-----|---------|---------|
| **Execution Model** | 3-phase lifecycle | Sequential | Wave-parallel | DAG | Role-based |
| **Parallel Execution** | Per-phase | No | Yes (waves) | Yes (DAG) | Yes (roles) |
| **Multi-Agent** | Specialists | No | No | Yes (beads) | Yes (agents) |
| **Approval Gates** | Yes | Implicit (PRs) | No | No | Optional |
| **Audit Trail** | Complete | Git history | STATE.md | Formula state | City state |
| **Human-in-Loop** | Required | PR review | Optional | Optional | Optional |
| **GitHub Integration** | External | Native | External | External | External |

## By Project Type

### Startup MVP

**Recommended**: [GSD](gsd.md)

```
Speed > Ceremony
Parallel execution
Fast iteration
Minimal overhead
```

```bash
visionspec init my-mvp --profile startup
visionspec export gsd
```

### Enterprise Feature

**Recommended**: [AWS AI-DLC](aws-aidlc.md)

```
Compliance required
Audit trails
Approval gates
Multi-stakeholder
```

```bash
visionspec init my-feature --profile enterprise
visionspec export aidlc
```

### Open Source Project

**Recommended**: [SpecKit](speckit.md)

```
GitHub-native
PR-based review
Sequential tasks
Community contributors
```

```bash
visionspec init my-oss --profile startup
visionspec export speckit
```

### Large-Scale System

**Recommended**: [GasTown](gastown.md) or [GasCity](gascity.md)

```
Complex dependencies
Multiple agents
Parallelizable work
Specialized roles
```

```bash
visionspec init my-platform --profile enterprise
visionspec export gastown  # or gascity
```

## By Team Structure

### Solo Developer

| Scenario | Target |
|----------|--------|
| Quick prototype | GSD |
| Careful implementation | SpecKit |
| Complex project | GasTown |

### Small Team (2-5)

| Scenario | Target |
|----------|--------|
| Collaborative PRs | SpecKit |
| Parallel workstreams | GSD |
| Shared agents | GasTown |

### Large Team (5+)

| Scenario | Target |
|----------|--------|
| Role specialization | GasCity |
| Compliance needs | AI-DLC |
| Mixed complexity | GasTown |

### Enterprise

| Scenario | Target |
|----------|--------|
| Regulated industry | AI-DLC |
| Multi-team coordination | GasCity |
| Complex orchestration | GasTown + GasCity |

## By Methodology Pairing

Different VisionSpec frameworks pair naturally with different targets:

### AWS Working Backwards + AI-DLC

```
Natural pairing
Same philosophy
Enterprise-grade
Full lifecycle
```

```bash
visionspec init project --profile enterprise
# Use AWS framework
visionspec export aidlc
```

### Lean Startup + GSD

```
Build-measure-learn
Fast iteration
Minimal ceremony
Quick pivots
```

```bash
visionspec init experiment --profile startup
# Use Lean Startup framework
visionspec export gsd
```

### Google Design Docs + SpecKit

```
RFC-style specs
Sequential implementation
PR-based review
Design doc tradition
```

```bash
visionspec init design --profile growth
# Use Google framework
visionspec export speckit
```

### Stripe API-First + GasCity

```
API contracts first
Role specialization
Backend/Frontend split
Clear boundaries
```

```bash
visionspec init api --profile enterprise
# Use Stripe framework
visionspec export gascity
```

## Complexity Analysis

### Simple Projects

**Characteristics**:

- Single component
- < 10 tasks
- No parallelization needed
- Single developer

**Best Targets**: SpecKit, GSD

### Medium Projects

**Characteristics**:

- Multiple components
- 10-50 tasks
- Some parallelization
- Small team

**Best Targets**: GSD, GasTown

### Complex Projects

**Characteristics**:

- Many components
- 50+ tasks
- High parallelization
- Large team
- Compliance needs

**Best Targets**: GasTown, GasCity, AI-DLC

## Migration Paths

### Starting Simple, Growing Complex

```
Start: GSD
  ↓ (need structure)
Migrate: SpecKit
  ↓ (need parallelization)
Migrate: GasTown
  ↓ (need roles)
Migrate: GasCity
```

### Starting Enterprise

```
Start: AI-DLC
  ↓ (need multi-agent)
Add: GasTown/GasCity for implementation
Keep: AI-DLC for governance
```

## Hybrid Approaches

### AI-DLC + GasCity

Use AI-DLC for lifecycle management, GasCity for implementation:

```bash
# Inception with AI-DLC
visionspec export aidlc

# Construction with GasCity
visionspec export gascity
```

### GasTown + GasCity

Use GasTown for task DAG, GasCity for agent coordination:

```bash
# Task breakdown
visionspec export gastown

# Agent assignment
visionspec export gascity
```

## Cost Considerations

| Target | Token Usage | Agent Count | Typical Cost |
|--------|-------------|-------------|--------------|
| GSD | Low | 1 | $ |
| SpecKit | Low | 1 | $ |
| AI-DLC | Medium | 1-3 | $$ |
| GasTown | High | 3-5 | $$$ |
| GasCity | High | 3-10 | $$$$ |

## Decision Flowchart

```
Start
  │
  ▼
Need compliance/audit?
  │
  ├─ Yes → AI-DLC
  │
  No
  │
  ▼
Need multi-agent?
  │
  ├─ Yes → Need role specialization?
  │          │
  │          ├─ Yes → GasCity
  │          │
  │          └─ No → GasTown
  │
  No
  │
  ▼
Need parallel execution?
  │
  ├─ Yes → GSD
  │
  No
  │
  ▼
GitHub-native workflow?
  │
  ├─ Yes → SpecKit
  │
  No
  │
  ▼
Default: GSD (simplest)
```

## Summary

| Target | Best For | Avoid When |
|--------|----------|------------|
| **AI-DLC** | Enterprise, compliance | Speed is critical |
| **SpecKit** | GitHub projects, PRs | Need parallelization |
| **GSD** | MVPs, fast iteration | Complex dependencies |
| **GasTown** | Complex DAGs | Simple projects |
| **GasCity** | Role-based teams | Solo development |

## Next Steps

- [AWS AI-DLC](aws-aidlc.md) - Enterprise lifecycle
- [GitHub SpecKit](speckit.md) - GitHub-native workflows
- [GSD](gsd.md) - Fast parallel execution
- [GasTown](gastown.md) - DAG orchestration
- [GasCity](gascity.md) - Role-based agents
