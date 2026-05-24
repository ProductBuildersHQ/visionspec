# Core Concepts

## Specification Types

MultiSpec organizes specifications into three categories:

### Source Specs (Human-Authored)

| Spec | Name | Purpose |
|------|------|---------|
| MRD | Market Requirements | Market problem, audience, business goals |
| PRD | Product Requirements | User stories, functional requirements |
| UXD | User Experience Design | User journeys, interaction flows |

### GTM Specs (LLM-Generated)

Generated using the Working Backwards methodology:

| Spec | Name | Input | Purpose |
|------|------|-------|---------|
| PRESS | Press Release | MRD + PRD | Customer announcement format |
| FAQ | FAQ Document | PRESS | Challenges claims, surfaces gaps |
| NARRATIVE | Narrative | MRD + PRD + FAQ | Internal vision document |

### Technical Specs (LLM-Generated)

| Spec | Name | Input | Purpose |
|------|------|-------|---------|
| TRD | Technical Requirements | MRD + PRD + UXD + CONSTITUTION | Architecture, APIs, data models |
| IRD | Infrastructure Requirements | TRD + CONSTITUTION | Deployment, scaling, operations |

## Directory Structure

```
docs/specs/
├── CONSTITUTION.md           # Repo-level governance (required)
├── ROADMAP.md                # Cross-project priorities
└── {project}/                # kebab-case project name
    ├── source/               # Human-authored
    │   ├── mrd.md
    │   ├── prd.md
    │   └── uxd.md
    ├── gtm/                  # LLM-generated GTM
    │   ├── press.md
    │   ├── faq.md
    │   └── narrative.md
    ├── technical/            # LLM-generated technical
    │   ├── trd.md
    │   └── ird.md
    ├── eval/                 # Evaluation results
    │   ├── mrd.eval.json
    │   ├── prd.eval.json
    │   └── ...
    ├── spec.md               # Reconciled execution spec
    ├── current-truth.md      # Post-ship maintained state
    └── multispec.yaml        # Project configuration
```

## Readiness Gates

Projects progress through readiness gates:

1. **Required specs present** - All required source specs exist
2. **Evaluations passing** - No critical/high findings
3. **Approvals obtained** - Required specs have approvals
4. **Execution spec generated** - spec.md created via reconciliation

## Evaluation System

Each spec can be evaluated using LLM-as-a-Judge:

- **Rubrics** define evaluation criteria per spec type
- **Findings** are categorized by severity (critical, high, medium, low)
- **Decisions** are pass, conditional, or fail

Evaluation results are stored as `{spec}.eval.json` in the `eval/` directory.

## Reconciliation

Reconciliation combines all approved specs into a unified execution spec:

1. Loads all approved source, GTM, and technical specs
2. Detects conflicts and missing traceability
3. Generates `spec.md` (unified execution spec)
4. Generates `spec.eval.json` (reconciliation evaluation)

## Target Adapters

The reconciled `spec.md` can be exported to various execution systems:

| Target | Format | Use Case |
|--------|--------|----------|
| SpecKit | spec.md, plan.md, tasks.md | GitHub-based workflows |
| GSD | PLAN.md, STATE.md | Get Shit Done methodology |
| GasTown | TOML formulas, Beads | Multi-agent orchestration |
| GasCity | city.toml | Agent city configuration |
| OpenSpec | JSON/YAML | Portable interchange format |

## CONSTITUTION.md

The `CONSTITUTION.md` file at `docs/specs/CONSTITUTION.md` defines repo-level constraints and patterns:

- Technology choices
- Coding standards
- Architecture patterns
- Security requirements
- Performance targets

All generated specs must adhere to the constitution.

## Requirement Graphs

MultiSpec extracts requirement graphs from specs for traceability analysis:

- **Nodes**: requirements, user stories, constraints, decisions, sections
- **Edges**: traces_to, derived_from, contains

Graphs enable visualization and traceability analysis:

```bash
# Extract graph from specs
multispec graph extract

# Export as interactive HTML
multispec graph export --format html

# Query specific node types
multispec graph query --type requirement --spec prd
```

See [Graph Command](../cli/graph.md) for full documentation.

## Post-Ship Alignment

After shipping, `current-truth.md` maintains alignment between specs and reality:

- Documents actual capabilities
- Notes divergences from spec
- Tracks limitations discovered in production
- Updates GTM docs with alignment notes
