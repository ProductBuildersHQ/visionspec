# Google Design Docs + RFC

Google's engineering culture emphasizes written communication through Design Docs and RFCs. This approach focuses on explicit tradeoffs, alternatives considered, and peer review before building.

## The Flow

![Google Design Docs + RFC Flow](../diagrams/google-flow.svg)

## Key Principles

1. **Explicit Tradeoffs**: Every decision has costs; make them visible
2. **Alternatives Considered**: Always evaluate multiple approaches
3. **Peer Review**: RFCs ensure diverse perspectives before building
4. **Data-Driven**: OKRs with measurable key results
5. **Experiment Culture**: A/B testing and statistical rigor

## VisionSpec Mapping

| Google Artifact | VisionSpec Type | Purpose |
|-----------------|-----------------|---------|
| OKR Document | MRD | Strategic objectives and key results |
| RFC | PRD | Product proposal for peer review |
| Design Doc | TRD | Technical specification with tradeoffs |
| Experiment Plan | UXD | A/B test design and success criteria |

## Using the Google Profile

### Initialize a Project

```bash
multispec init my-feature --profile google
```

### Create OKRs (MRD)

```bash
multispec draft mrd -p my-feature
```

The MRD template becomes an OKR document with:

- Objectives (qualitative, inspiring)
- Key Results (specific, measurable)
- Strategic alignment
- Resource allocation

### Create RFC (PRD)

```bash
multispec draft prd -p my-feature
```

The RFC template includes:

- Problem statement
- Goals and Non-Goals
- Proposed solution
- Alternatives considered
- Reviewers and comment period

### Synthesize Design Doc (TRD)

```bash
multispec synthesize trd -p my-feature
```

The Design Doc includes:

- Context and scope
- Goals and Non-Goals
- Alternatives analysis
- Explicit tradeoffs
- Cross-cutting concerns
- Implementation plan

### Create Experiment Plan (UXD)

```bash
multispec draft uxd -p my-feature
```

The Experiment Plan includes:

- Hypothesis (If/Then/Because format)
- Success criteria (pre-defined)
- Power analysis
- Metrics design
- Analysis plan

## Rubric Categories

### OKR Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Objective Quality | 20% | Inspiring and appropriately ambitious |
| Key Result Measurability | 25% | Specific with baselines and targets |
| Strategic Alignment | 15% | Connected to higher-level strategy |
| Ambition Level | 15% | Stretch targets, not sandbagged |
| Resource Alignment | 10% | Team allocation mapped |
| Dependencies/Risks | 10% | Dependencies mapped with fallbacks |
| Check-In Structure | 5% | Tracking cadence defined |

### RFC Evaluation (PRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Problem Clarity | 15% | Problem and motivation articulated |
| Goals/Non-Goals | 15% | Scope well-bounded with exclusions |
| Alternatives/Tradeoffs | 15% | Multiple options with honest analysis |
| Peer Review Setup | 15% | Reviewers listed, feedback requested |
| Risk Assessment | 10% | Risks identified with mitigations |
| Success Metrics | 15% | Clear, measurable criteria |
| Launch Plan | 10% | Phased rollout with rollback |
| Open Questions | 5% | Unknowns surfaced for discussion |

### Design Doc Evaluation (TRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Goals/Non-Goals | 15% | Specific, measurable, bounded |
| Alternatives Analysis | 20% | Genuine consideration of options |
| Explicit Tradeoffs | 20% | What we gain, what we give up |
| Technical Depth | 15% | Detailed enough for implementation |
| Cross-Cutting Concerns | 15% | Security, scalability, observability |
| Implementation Plan | 10% | Milestones, migration, rollback |
| Open Questions | 5% | Unknowns acknowledged |

### Experiment Evaluation (UXD)

| Category | Weight | Description |
|----------|--------|-------------|
| Hypothesis Clarity | 20% | If/Then/Because format |
| Pre-Defined Success Criteria | 20% | Criteria before experiment runs |
| Statistical Rigor | 20% | Power analysis, appropriate methods |
| Metrics Design | 15% | Primary, secondary, guardrails |
| Implementation Plan | 10% | Phased rollout, stopping criteria |
| Analysis Plan | 10% | Pre-analysis checks defined |
| Risk Mitigation | 5% | Rollback plan exists |

## Example Workflow

```bash
# 1. Initialize project
multispec init search-ranking --profile google

# 2. Define OKRs
multispec draft mrd -p search-ranking
multispec eval mrd -p search-ranking
multispec approve mrd -p search-ranking

# 3. Write RFC for peer review
multispec draft prd -p search-ranking
multispec eval prd -p search-ranking
# ... gather peer feedback ...
multispec approve prd -p search-ranking

# 4. Synthesize Design Doc
multispec synthesize trd -p search-ranking
multispec eval trd -p search-ranking
multispec approve trd -p search-ranking

# 5. Define Experiment Plan
multispec draft uxd -p search-ranking
multispec eval uxd -p search-ranking
multispec approve uxd -p search-ranking

# 6. Check status
multispec status -p search-ranking
```

## Reference Materials

For deeper understanding of Google's engineering practices, see:

- [Google Engineering Practices](https://google.github.io/eng-practices/)
- [Design Docs at Google](https://www.industrialempathy.com/posts/design-docs-at-google/)
- *Software Engineering at Google* (O'Reilly)
- Internal reference: `frameworks-internal/google-design-docs/`
