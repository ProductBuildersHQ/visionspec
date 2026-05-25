# MultiSpec Roadmap

Multi-domain specification orchestration for humans and AI agents.

## Vision

MultiSpec bridges the gap between organizational intent (MRD, PRD, UXD) and executable specifications for AI coding agents. It provides:

- **Domain-specific authoring** - Separate specs for PM, UX, Engineering
- **GTM synthesis** - LLM-generated press releases, FAQs, narratives (Working Backwards)
- **Technical synthesis** - LLM-generated TRD, IRD from source specs
- **Structured evaluation** - Per-domain LLM judges with customizable rubrics
- **Reconciliation** - Conflict detection and tradeoff resolution
- **Target adapters** - Export to SpecKit, GSD, GasTown, GasCity, OpenSpec
- **Post-ship alignment** - Maintain current-truth after shipping

### Document Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────┐
│ HUMAN-AUTHORED (Source)                                                 │
│   mrd.md → prd.md → uxd.md                                              │
└─────────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────────┐
│ LLM-GENERATED (GTM) ← Working Backwards methodology                     │
│   press.md → faq.md → narrative.md                                      │
└─────────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────────┐
│ LLM-GENERATED (Technical)                                               │
│   trd.md → ird.md                                                       │
└─────────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────────┐
│ RECONCILIATION                                                          │
│   All approved specs → spec.md (execution spec)                         │
└─────────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────────┐
│ TARGET EXPORT                                                           │
│   spec.md → SpecKit | GSD | GasTown | GasCity | OpenSpec                │
└─────────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────────┐
│ POST-SHIP ALIGNMENT                                                     │
│   spec.md + shipped reality → current-truth.md                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Directory Structure (Canonical)

```
docs/specs/
├── CONSTITUTION.md                    # Repo-level governance (CAPS)
├── ROADMAP.md                         # Cross-project priorities (CAPS)
└── {project}/                         # kebab-case project name
    ├── source/                        # Human-authored specs
    │   ├── mrd.md
    │   ├── prd.md
    │   └── uxd.md
    ├── gtm/                           # LLM-generated GTM docs
    │   ├── press.md
    │   ├── faq.md
    │   └── narrative.md
    ├── technical/                     # LLM-generated technical docs
    │   ├── trd.md
    │   └── ird.md
    ├── eval/                          # All evaluations (centralized)
    │   ├── mrd.eval.json
    │   ├── prd.eval.json
    │   ├── uxd.eval.json
    │   ├── press.eval.json
    │   ├── faq.eval.json
    │   ├── narrative.eval.json
    │   ├── trd.eval.json
    │   ├── ird.eval.json
    │   └── spec.eval.json
    ├── .graphize/                     # Requirement graph (via graphize)
    ├── spec.md                        # Reconciled execution spec
    ├── current-truth.md               # Post-ship maintained state
    ├── status.html                    # Project readiness report
    ├── index.md                       # MkDocs project page (generated)
    └── multispec.yaml                 # Project configuration
```

### Naming Conventions (Enforced)

| Element | Convention | Example |
|---------|------------|---------|
| Project directory | `kebab-case` | `user-onboarding`, `user-onboarding` |
| Spec files | `lowercase.md` | `mrd.md`, `prd.md`, `spec.md` |
| Eval files | `{spec}.eval.json` | `mrd.eval.json`, `press.eval.json` |
| Config file | `multispec.yaml` | Fixed name |
| Repo-level docs | `CAPS.md` | `CONSTITUTION.md`, `ROADMAP.md` |

**Design principles:**
- Specs (markdown) for humans, evals (JSON) for machines
- Centralized evals enable easy status aggregation
- Fixed naming enables automation without configuration
- CAPS for repo-level canonical docs (like README.md)
- `docs/` directory integrates with MkDocs for documentation sites

---

## Phase 0: Project Foundation

Core project setup and CLI scaffolding.

- [x] RMI-001: Initialize Go module (`github.com/plexusone/multispec`)
- [x] RMI-002: Create CLI skeleton with Cobra (`multispec` command)
- [x] RMI-003: Define core types package (`pkg/types/`)
- [x] RMI-004: Add configuration loading (`multispec.yaml`)
- [x] RMI-005: Set up CI (lint, test, build)
  - `.github/workflows/go-ci.yaml` - build and test
  - `.github/workflows/go-lint.yaml` - golangci-lint
  - `.github/workflows/go-sast-codeql.yaml` - security analysis
- [x] RMI-006: Create project README

- [x] RMI-007: Implement `multispec lint` command
  - Validate directory structure matches canonical layout
  - Validate file naming conventions (lowercase specs, kebab-case projects)
  - Report errors for non-standard names
  - Exit non-zero for CI integration

- [x] RMI-008: Implement MCP server skeleton
  - MCP tools: list_projects, get_project_status, get_spec, get_eval
  - MCP tools: run_eval, synthesize, reconcile, approve, export
  - Stdio transport support

- [x] RMI-009: Connect MCP handlers to library code
  - list_projects → scan docs/specs/ directory
  - get_project_status → pkg/status.Generate()
  - get_spec → read spec file content (stub)
  - Other handlers remain stubs until Phase 2-4

---

## Phase 1: Directory Structure & Source Specs

Establish conventions for spec organization and authoring.

### Directory Structure

- [x] RMI-010: Implement `multispec init` command
  - Create `docs/specs/{project}/` structure
  - Create `source/`, `gtm/`, `technical/`, `eval/` subdirectories
  - Generate `multispec.yaml` project config

- [x] RMI-011: Support CONSTITUTION.md at `docs/specs/CONSTITUTION.md`
  - Repo-level governance document
  - Optional org-level at `~/.config/multispec/CONSTITUTION.md`
  - `pkg/config/config.go` - FindConstitution, LoadConstitution functions
  - Used in synth, reconcile, and export commands

### MkDocs Integration

- [x] RMI-016: Generate `{project}/index.md` for each project
  - Spec overview with status badges
  - Links to all specs (source, gtm, technical)
  - Eval summary (pass/fail counts, open findings)
  - Last updated timestamps
  - `pkg/mkdocs/mkdocs.go` - GenerateProjectIndex, WriteProjectIndex
  - `multispec docs project` command

- [x] RMI-017: Generate `docs/specs/index.md` (specs landing page)
  - List all projects with status
  - Link to CONSTITUTION.md and ROADMAP.md
  - Cross-project metrics
  - `pkg/mkdocs/mkdocs.go` - GenerateSpecsLanding, WriteSpecsLanding
  - `multispec docs generate` command

- [ ] RMI-018: Generate MkDocs navigation structure
  - Auto-update `mkdocs.yml` nav section
  - Or generate `nav.yml` partial for include
  - Support `mkdocs-awesome-pages-plugin` `.pages` files

- [ ] RMI-019: Render eval JSON to markdown for MkDocs
  - `multispec render-evals {project}`
  - Generate `eval/index.md` with rendered findings
  - Collapsible sections per spec
  - Severity badges and status indicators

### Project Status Report

- [x] RMI-019a: Implement `multispec status` core logic
  - `pkg/status/status.go` - Generate() function
  - Check spec existence per type
  - Check eval file existence
  - Check approval status
  - Calculate readiness gates

- [x] RMI-019b: Implement status renderers
  - `RenderText()` - Terminal output with colors/icons
  - `RenderHTML()` - Browser/MkDocs report with traffic light
  - `RenderMarkdown()` - For embedding in index.md
  - JSON output already works via CLI

- [x] RMI-019c: Define readiness gates
  - All required source specs present (mrd, prd, uxd, trd)
  - All evals passing (no critical/high findings)
  - All required approvals obtained
  - spec.md generated

- [x] RMI-019d: Integrate graphize metrics in status report
  - Traceability coverage percentage
  - Requirements without TRD coverage
  - Conflict count
  - Link to graph visualization
  - `pkg/status/status.go` - GraphMetrics struct, RenderText/RenderMarkdown
  - `pkg/specgraph/specgraph.go` - ComputeMetrics function

- [x] RMI-019e: CI exit codes for readiness
  - `multispec status --ci` exits non-zero if not ready
  - CLI flag wired up, needs renderer to output before exit

### Source Spec Templates

- [x] RMI-012: Create mrd.md template (Market Requirements)
  - Market problem, target audience, competitive landscape
  - Business metrics, success criteria

- [x] RMI-013: Create prd.md template (Product Requirements)
  - User stories, functional requirements
  - Acceptance criteria, priorities

- [x] RMI-014: Create uxd.md template (User Experience Design)
  - User journeys, interaction flows
  - Accessibility requirements

- [x] RMI-014a: Create trd.md template (Technical Requirements)
  - Architecture overview, API contracts
  - Data models, technical constraints

- [x] RMI-014b: Create ird.md template (Infrastructure Requirements)
  - Infrastructure architecture, compute, storage
  - Security, observability, DR planning

- [x] RMI-014c: Create Press Release template (Working Backwards)
  - Headline, customer problem, solution
  - Customer quote, call to action

- [x] RMI-014d: Create FAQ template
  - Question coverage across categories
  - Pricing, getting started, objection handling

- [x] RMI-014e: Create Narrative templates (1-pager and 6-pager)
  - Executive summary format (1-pager)
  - AWS 6-pager format with appendices

- [x] RMI-015: Implement `multispec create {spec-type}` command
  - Scaffold new spec from template
  - Support: mrd, prd, uxd, trd, ird, press, faq, narrative-1p, narrative-6p
  - `pkg/cli/commands.go` - createCmd function
  - `pkg/cli/cli.go` - registered in CommandSet

---

## Phase 2: Evaluation Engine

Integrate with `structured-evaluation` for per-spec evaluation.

### Rubric System

- [x] RMI-020: Define rubric file format (Go structs, leveraging `structured-evaluation`)
  - Categories, weights, scales (categorical with range anchors)
  - Pass criteria, severity thresholds

- [x] RMI-021: Create default rubrics
  - `pkg/rubrics/mrd.go` - Market requirements evaluation
  - `pkg/rubrics/prd.go` - Product requirements evaluation
  - `pkg/rubrics/uxd.go` - UX design evaluation
  - `pkg/rubrics/trd.go` - Technical requirements evaluation
  - `pkg/rubrics/ird.go` - Infrastructure requirements evaluation
  - `pkg/rubrics/press.go` - Press release evaluation
  - `pkg/rubrics/faq.go` - FAQ evaluation
  - `pkg/rubrics/narrative1p.go` - 1-pager narrative evaluation
  - `pkg/rubrics/narrative6p.go` - 6-pager narrative evaluation

- [x] RMI-022: Support custom rubrics in project config
  - Override default rubrics per project
  - Rubric inheritance/extension
  - `pkg/types/project.go` - RubricsConfig struct
  - `pkg/rubrics/loader.go` - FileLoader, ChainLoader for custom rubrics

### Evaluation Commands

- [x] RMI-023a: Implement MCP `run_eval` tool
  - Load spec and rubric
  - Call LLM judge via omnillm-core
  - Return evaluation results with findings

- [x] RMI-023b: Implement MCP `eval_draft` tool
  - Evaluate draft content before finalization
  - Track eval history in draft metadata

- [ ] RMI-023c: Implement `multispec eval {spec-type}` CLI command
  - Load spec and rubric
  - Call LLM judge
  - Write `{spec}.eval.json` output

- [ ] RMI-024: Implement `multispec eval --all` command
  - Evaluate all source specs, GTM docs, and technical docs
  - Generate all `*.eval.json` files
  - Support filtering: `--source`, `--gtm`, `--technical`

- [ ] RMI-025: Implement `multispec render {eval-file}` command
  - Render JSON eval to Markdown for human review
  - Use `structured-evaluation/render/markdown`

- [x] RMI-026: Implement `multispec status` command
  - Summary of open items across all evals
  - Severity counts, blocking issues

### AI Co-Authoring (Draft Workflow)

- [x] RMI-026a: Implement draft package (`pkg/draft/`)
  - Draft CRUD operations (Start, Get, Update, Discard, Finalize)
  - Session management with status tracking
  - Eval history persistence

- [x] RMI-026b: Implement MCP draft tools
  - `start_draft` - Initialize draft from template
  - `get_draft` - Retrieve draft content and metadata
  - `update_draft` - Save draft content with versioning
  - `eval_draft` - Evaluate draft against rubric
  - `finalize_draft` - Promote draft to final spec
  - `discard_draft` - Delete draft
  - `list_drafts` - List all drafts in project

- [x] RMI-026c: Implement LLM evaluation integration
  - `pkg/eval/eval.go` - Evaluation orchestration
  - `pkg/eval/llm.go` - LLM provider integration via omnillm-core
  - Support project-level LLM config in multispec.yaml

- [x] RMI-026d: Create authoring skills
  - `skills/author-mrd/SKILL.md`
  - `skills/author-prd/SKILL.md`
  - `skills/author-uxd/SKILL.md`
  - `skills/author-trd/SKILL.md`
  - `skills/author-ird/SKILL.md`
  - `skills/author-press/SKILL.md`
  - `skills/author-faq/SKILL.md`
  - `skills/author-narrative-1p/SKILL.md`
  - `skills/author-narrative-6p/SKILL.md`

---

## Phase 3: GTM & Technical Synthesis

LLM-generated documents from source specs + constitution.

### GTM Document Generation (Working Backwards)

- [x] RMI-027: Implement `multispec synthesize press` command
  - Input: MRD + PRD
  - Output: `gtm/press.md` (press release format)
  - Template: Hook → Problem → Solution → Quote → CTA → Benefits
  - Generate PRESS_EVAL.json

- [x] RMI-028: Implement `multispec synthesize faq` command
  - Input: press.md
  - Output: `gtm/faq.md`
  - Structure: External FAQs + Internal FAQs
  - Challenge claims in press release
  - Generate FAQ_EVAL.json

- [x] RMI-029: Implement `multispec synthesize narrative` command
  - Input: MRD + PRD + FAQ
  - Output: `gtm/narrative.md`
  - Structure: Customer → Tension → Future State → Promise → Principles → Non-Goals
  - Generate NARRATIVE_EVAL.json

### GTM Evaluation Rubrics

- [x] RMI-029a: Create press release rubric (`pkg/rubrics/press.go`)
  - Categories: headline-impact, customer-problem, solution-clarity, customer-validation, call-to-action, readability

- [x] RMI-029b: Create FAQ rubric (`pkg/rubrics/faq.go`)
  - Categories: question-coverage, answer-clarity, customer-language, pricing-transparency, getting-started, objection-handling

- [x] RMI-029c: Create narrative rubrics
  - `pkg/rubrics/narrative1p.go` - 1-pager evaluation
  - `pkg/rubrics/narrative6p.go` - 6-pager evaluation (AWS format)

- [x] RMI-029d: Support `--eval` flag on synthesize commands
  - `multispec synthesize press --eval` generates press.md + press.eval.json
  - Auto-evaluate after generation

### Technical Document Generation

### TRD Generation

- [x] RMI-030: Implement `multispec synthesize trd` command
  - Input: MRD + PRD + UXD + CONSTITUTION
  - Output: `technical/trd.md`
  - Generate TRD_EVAL.json

- [x] RMI-031: Define TRD template structure
  - Architecture overview
  - API contracts
  - Data models
  - Technical constraints
  - Traceability to source requirements

### IRD Generation

- [x] RMI-032: Implement `multispec synthesize ird` command
  - Input: TRD + CONSTITUTION
  - Output: `technical/ird.md`
  - Generate IRD_EVAL.json

- [x] RMI-033: Define IRD template structure
  - Infrastructure requirements
  - Deployment architecture
  - Scaling considerations
  - Operational requirements

### Approval Workflow

- [x] RMI-034: Implement `multispec approve {spec-type}` command
  - Record approval in `multispec.yaml`
  - Track approver, timestamp
  - Gate for reconciliation

- [x] RMI-035: Support approval status in `multispec status`
  - Show pending approvals
  - Show approval history

### Post-Ship Alignment

- [ ] RMI-036: Implement `multispec align` command (moved to Phase 11 with context)
  - Input: spec.md + shipped reality (from engineering)
  - Output: `current-truth.md`
  - Detect: ungrounded claims, missed opportunities, drift
  - Update GTM docs with alignment notes

- [ ] RMI-037: Define current-truth.md structure
  - Product summary (current state)
  - Active capabilities table
  - Known boundaries/limitations
  - Source specs and evidence
  - Recent alignment notes

---

## Phase 4: Reconciliation Engine

Conflict detection and unified spec generation.

### Conflict Detection

- [x] RMI-040: Implement conflict detection algorithm
  - Cross-spec requirement conflicts
  - Constraint violations
  - Missing traceability
  - `pkg/reconcile/conflicts.go` - ConflictDetector with pattern-based detection

- [x] RMI-041: Define conflict representation
  - Conflict type (requirement, constraint, tradeoff, missing)
  - Source specs involved
  - Severity level (high, medium, low)
  - Suggested resolution
  - Confidence score for detected conflicts

### spec.md Generation

- [x] RMI-042: Implement `multispec reconcile` command
  - Input: All approved specs
  - Output: `spec.md` (unified execution spec)
  - Output: `spec.eval.json` (reconciliation evaluation)
  - Pre-reconciliation conflict detection included in LLM prompt

- [x] RMI-043: Define spec.md structure
  - Resolved requirements
  - Consolidated constraints
  - Task decomposition
  - Dependency graph
  - Decision log (tradeoffs made)
  - Traceability matrix

- [x] RMI-044: Support unresolved conflicts in spec.eval.json
  - Conflicts requiring human decision
  - Status: reconciled, reconciled_with_tradeoffs, needs_review
  - Decision log with resolutions

---

## Phase 5: Target Adapters

Export reconciled specs to downstream execution systems.

### Adapter Framework

- [x] RMI-050: Define `Target` interface
  - `Name()`, `Description()`, `Capabilities()`
  - `Validate()`, `Export()`
  - `pkg/target/target.go`

- [x] RMI-051: Implement target registry
  - Register adapters by name
  - List available targets
  - `Get()`, `Available()`, `ListTargets()`

- [x] RMI-052: Implement `multispec targets` command
  - List available targets
  - Show capabilities

- [x] RMI-053: Implement `multispec export {target}` command
  - Route to appropriate adapter
  - Support multiple targets: `multispec export speckit,gsd`

### SpecKit Adapter (Priority 1)

- [x] RMI-060: Implement SpecKit adapter
  - Generate `specs/{seq}-{name}/spec.md`
  - Generate `specs/{seq}-{name}/plan.md`
  - Generate `specs/{seq}-{name}/tasks.md`
  - `pkg/target/speckit.go`

- [x] RMI-061: Support SpecKit constitution sync
  - Update `.specify/memory/constitution.md` from CONSTITUTION.md
  - `pkg/target/speckit.go` - syncConstitution method
  - `pkg/cli/commands.go` - pass constitution path to export config

- [x] RMI-062: Support SpecKit branch conventions
  - Sequential (`001-feature`) or timestamp naming

### GSD Adapter (Priority 2)

- [x] RMI-070: Implement GSD adapter
  - Generate `PLAN.md` files with YAML frontmatter + XML tasks
  - Generate initial `STATE.md`
  - Generate `.planning/config.json`
  - `pkg/target/gsd.go`

- [x] RMI-071: Map requirements to `must_haves`
  - `must_haves.truths` from acceptance criteria
  - `must_haves.artifacts` from deliverables
  - `must_haves.key_links` from dependencies

- [x] RMI-072: Support GSD phases
  - Map spec phases to GSD phase structure
  - Generate wave dependencies

### GasTown Adapter (Priority 3)

- [x] RMI-080: Implement GasTown adapter
  - Generate TOML formulas (convoy/workflow/expansion)
  - Generate Bead definitions
  - `pkg/target/gastown.go`

- [x] RMI-081: Support formula types
  - Convoy for parallel review
  - Workflow for sequential execution
  - Expansion for template-based generation

- [x] RMI-082: Map task dependencies to Bead DAG
  - Blocked/ready relationships
  - Convoy coordination

### GasCity Adapter (Priority 3)

- [x] RMI-085: Implement GasCity adapter
  - Generate `city.toml` agent configuration
  - Generate agent definitions
  - Generate orders
  - `pkg/target/gascity.go`

### OpenSpec Adapter (Future)

- [ ] RMI-090: Define OpenSpec export format
  - Portable JSON/YAML structure
  - Agent-agnostic representation

- [ ] RMI-091: Implement OpenSpec adapter
  - Standards-compliant export
  - Interoperability focus

---

## Phase 6: Claude Code / Kiro CLI Integration

Seamless integration with AI coding assistant workflows via multi-agent-spec and assistantkit.

### Skill Definitions (multi-agent-spec)

- [x] RMI-098: Add Skill schema to multi-agent-spec
  - `sdk/go/skill.go` - Skill struct with builder methods
  - `schema/skill/skill.schema.json` - JSON Schema
  - Loader functions for skill directories
  - Matches assistantkit canonical type

- [ ] RMI-099: Define multispec skills in multi-agent-spec format
  - `multispec-status` - Check project readiness
  - `multispec-lint` - Validate project structure
  - `multispec-eval` - Run evaluations
  - `multispec-synthesize` - Generate specs
  - `multispec-reconcile` - Generate unified spec
  - `multispec-export` - Export to targets

### Skill Generation (assistantkit)

- [ ] RMI-100: Generate Claude Code skills via assistantkit
  - `skills/multispec-status/SKILL.md`
  - `skills/multispec-lint/SKILL.md`
  - etc.

- [ ] RMI-101: Generate Kiro CLI steering files via assistantkit
  - `steering/multispec-status.md`
  - `steering/multispec-lint.md`
  - etc.

### Automation

- [ ] RMI-102: Implement `multispec watch` command
  - File watcher for spec changes
  - Auto-run eval on change
  - Auto-run reconcile when all approved

- [ ] RMI-103: Support git hooks
  - Pre-commit: validate specs
  - Post-commit: trigger eval

---

## Phase 7: Graphize Integration

Requirement graph visualization via `github.com/plexusone/graphize`.

### Spec Extractor

- [ ] RMI-140: Create spec extractor for graphize
  - New extractor in `graphize/pkg/extract/spec/`
  - Parse markdown specs (mrd.md, prd.md, trd.md, etc.)
  - Extract requirements, constraints, decisions as nodes
  - Infer relationships as edges

- [ ] RMI-141: Define spec node types
  - `requirement` - Functional requirements from PRD
  - `user_story` - User stories from PRD
  - `constraint` - Constraints from CONSTITUTION, TRD
  - `acceptance_criteria` - Testable criteria
  - `decision` - Architectural decisions from TRD
  - `tradeoff` - Explicit tradeoffs from reconciliation
  - `capability` - Current capabilities from CURRENT-TRUTH

- [ ] RMI-142: Define spec edge types
  - `traces_to` - Requirement traceability (PRD → TRD)
  - `derived_from` - Synthesis source (TRD → MRD + PRD)
  - `conflicts_with` - Detected conflicts
  - `satisfies` - Implementation satisfies requirement
  - `depends_on` - Requirement dependencies
  - `blocks` - Blocking relationships
  - `supersedes` - Decision replacement

### Graph Storage

- [ ] RMI-143: Store spec graph in project directory
  - `.graphize/` directory under `docs/specs/{project}/`
  - Version controlled with project
  - One file per node/edge (git-friendly)

- [ ] RMI-144: Implement `multispec graph` commands
  - `multispec graph extract` - Build graph from specs
  - `multispec graph query` - Query relationships
  - `multispec graph export` - Export to HTML/JSON/GraphML

### Traceability Analysis

- [ ] RMI-145: Implement traceability reports
  - Requirements without TRD coverage
  - TRD tasks without PRD traceability
  - Orphaned constraints
  - Missing acceptance criteria

- [ ] RMI-146: Conflict detection via graph
  - Query `conflicts_with` edges
  - Highlight in reconciliation
  - Surface in SPEC_EVAL

### Visualization

- [ ] RMI-147: Generate spec graph HTML visualization
  - Interactive Cytoscape.js graph
  - Color-coded by spec type (PRD=blue, TRD=green, CONSTITUTION=red)
  - Filter by relationship type
  - Search by requirement ID

- [ ] RMI-148: MkDocs graph integration
  - Embed graph visualization in project index.md
  - Or link to standalone HTML export

---

## Phase 8: Advanced Features

Future enhancements.

### Multi-Project Support

- [ ] RMI-110: Support cross-project dependencies
  - Project references in spec.md
  - Cross-project reconciliation

- [ ] RMI-111: Implement `docs/specs/ROADMAP.md` generation
  - Aggregate project statuses
  - Prioritization tracking

### Organizational Memory

- [ ] RMI-120: Decision log persistence
  - Track tradeoffs across projects
  - Searchable decision history

- [ ] RMI-121: Rationale graphs
  - Link decisions to requirements
  - Impact analysis

### Analytics

- [ ] RMI-130: Evaluation metrics dashboard
  - Spec quality trends
  - Common failure patterns

- [ ] RMI-131: Reconciliation metrics
  - Conflict frequency
  - Resolution time

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/plexusone/structured-evaluation` | Rubric and evaluation types |
| `github.com/plexusone/omnillm-core` | LLM provider abstraction |
| `github.com/plexusone/graphize` | Requirement graph extraction and visualization |
| `github.com/modelcontextprotocol/go-sdk` | MCP server implementation |
| `github.com/spf13/cobra` | CLI framework |
| `github.com/spf13/viper` | Configuration |
| `github.com/fsnotify/fsnotify` | File watching |
| `gopkg.in/yaml.v3` | YAML parsing for profiles and rubrics |
| `github.com/gorilla/websocket` | Real-time collaboration (future) |

---

## Target Compatibility Matrix

| Feature | SpecKit | GSD | GasTown | GasCity | OpenSpec |
|---------|---------|-----|---------|---------|----------|
| Sequential tasks | Yes | Yes | Yes | Yes | Yes |
| Parallel execution | No | Yes (waves) | Yes (convoy) | Yes | TBD |
| Multi-agent | No | No | Yes | Yes | TBD |
| Verification | Implicit | Yes | Yes | Yes | TBD |
| Dependency graph | Yes | Yes | Yes (Beads) | Yes | TBD |

---

## Version Milestones

| Version | Phase | Key Deliverables |
|---------|-------|------------------|
| v0.1.0 | 0-1 | CLI skeleton, directory structure, templates |
| v0.2.0 | 2, 7 | Evaluation engine, rubrics, graphize integration |
| v0.3.0 | 9 | Composability (custom templates, rubrics, profiles, CLI as library) |
| v0.4.0 | 11 | Context Sources / Grounding (git repos, graphize, MCP servers) |
| v0.5.0 | 4 | Reconciliation engine, spec.md generation |
| v0.6.0 | 5a | SpecKit adapter |
| v0.7.0 | 5b | GSD adapter |
| v0.8.0 | 5c | GasTown/GasCity adapters |
| v0.9.0 | 6 | Claude Code skills, Kiro integration |
| v0.10.0 | 10a | Testing, profile CLI enhancements, MCP resources |
| v0.11.0 | 10b | Export targets (Linear, Jira, Notion, Confluence) |
| v0.12.0 | 10c | Spec versioning, cross-project analysis |
| v0.13.0 | 10d | Real-time collaboration |
| v0.14.0 | 10e | CI/CD integration (GitHub Actions, pre-commit) |
| v1.0.0 | 8 | Production release with full feature set |

---

## Phase 9: Composability (v0.3.0)

Enable organizations (companies, open source projects, non-profits) to compose custom CLI tools with multispec as a library.

### CLI as Library

- [x] RMI-200: Move CLI commands from `internal/cli` to `pkg/cli`
  - Create `pkg/cli/cli.go` with composition API
  - `AddCommandsTo(root *cobra.Command, cfg *Config)`
  - `Commands(cfg *Config)` for selective command access
  - `DefaultConfig()` for multispec defaults

- [x] RMI-201: Update `cmd/multispec/main.go` to use `pkg/cli`
  - `internal/cli/root.go` now uses `pkg/cli.AddCommandsTo()`

### Custom Templates

- [x] RMI-210: Create template Loader interface (`pkg/templates/loader.go`)
  - `Loader` interface with `Load()` and `Available()`
  - `EmbeddedLoader()` - wraps current embedded templates
  - `NewFileLoader(dir)` - loads from directory
  - `NewChainLoader(loaders...)` - tries loaders in order

- [x] RMI-211: Support custom spec types from templates
  - Allow non-standard spec types (e.g., `security.md`, `compliance.md`)
  - Register custom types with category

### Custom Rubrics

- [x] RMI-220: Create rubric Loader interface (`pkg/rubrics/loader.go`)
  - `Loader` interface with `Load()` and `Available()`
  - `EmbeddedLoader()` - wraps current Go-defined rubrics
  - `NewFileLoader(dir)` - loads YAML rubrics
  - `NewChainLoader(loaders...)`

- [x] RMI-221: Define rubric YAML schema (`pkg/rubrics/yaml.go`)
  - `RubricYAML` struct for parsing
  - Validation and conversion to `RubricSet`
  - Compatible with structured-evaluation

### Configurable Spec Requirements

- [x] RMI-230: Add `SpecConfig` types (`pkg/types/spec_config.go`)
  - `SpecRequirement` struct (required, category, template, rubric)
  - `SpecConfig` with helper methods
  - `IsRequired()` with fallback to defaults

- [ ] RMI-231: Update `multispec.yaml` schema
  - Add `specs:` section for per-spec configuration
  - Parse and merge with defaults

- [ ] RMI-232: Update `SpecType.IsRequired()` to use config
  - Check project config first, then defaults

### Documentation & Examples

- [x] RMI-240: Create example org CLI (`examples/org-cli/`)
  - Sample CLI importing multispec as library
  - Custom templates and rubrics
  - Custom spec types

- [ ] RMI-241: Add organization integration docs (`docs/organization/`)
  - Integration guide
  - Configuration reference
  - Template and rubric customization

### Configuration Profiles

- [x] RMI-250: Create profile system (`pkg/profiles/`)
  - `Profile` type with Name, Description, Extends, SpecConfig
  - `ProfileLoader` interface with `Load()` and `Available()`
  - `EmbedFSLoader` - loads from embedded filesystem
  - `FileLoader` - loads from directory
  - `ChainLoader` - tries loaders in order
  - `ResolvingLoader` - resolves profile inheritance

- [x] RMI-251: Create default profiles
  - `0-1` - Minimal profile with hypothesis document only
  - `startup` - PRD only for pre-PMF startups
  - `growth` - PRD + UXD + FAQ for 1-N scaling (extends startup)
  - `enterprise` - Full spec suite with security/compliance

- [x] RMI-252: Add profile CLI commands
  - `profiles list` - List available profiles
  - `profiles show <name>` - Show profile details
  - `--profile` flag on init command

- [x] RMI-253: Update example CLIs to use profiles
  - `examples/0-1-product/` - Uses "0-1" profile
  - `examples/pre-pmf-startup/` - Uses "startup" profile
  - `examples/1-n-growth/` - Uses "growth" profile
  - `examples/post-pmf-enterprise/` - Uses "enterprise" profile

- [ ] RMI-254: Add profile tests (`pkg/profiles/*_test.go`)
  - Test profile loading and inheritance
  - Test profile merging
  - Test template/rubric loader creation

---

## Phase 10: Platform Enhancements

Future enhancements for testing, integrations, and developer experience.

### Testing & Quality

- [ ] RMI-300: Add comprehensive profile tests
  - Unit tests for `pkg/profiles/`
  - Profile inheritance testing
  - Loader chain testing

- [ ] RMI-301: Add MCP integration tests
  - Test all 17 MCP tools
  - Mock LLM responses for eval testing
  - Draft workflow end-to-end tests

- [ ] RMI-302: Add end-to-end authoring workflow tests
  - start_draft → update_draft → eval_draft → finalize_draft
  - Test with real project structure
  - Verify file system operations

### Profile CLI Enhancements

- [ ] RMI-310: Implement `profiles create <name>` command
  - Interactive profile creation wizard
  - Select base profile to extend
  - Choose required spec types
  - Generate profile.yaml

- [ ] RMI-311: Implement `profiles extend <base> <name>` command
  - Create profile extending another
  - Override specific settings
  - Custom templates/rubrics directory

- [ ] RMI-312: Implement profile validation
  - `profiles validate <path>` command
  - Check profile.yaml schema
  - Verify referenced templates/rubrics exist

### MCP Resources

- [ ] RMI-320: Expose templates as MCP resources
  - `templates://` URI scheme
  - List available templates
  - Read template content

- [ ] RMI-321: Expose rubrics as MCP resources
  - `rubrics://` URI scheme
  - List available rubrics
  - Read rubric definitions

- [ ] RMI-322: Expose profiles as MCP resources
  - `profiles://` URI scheme
  - List available profiles
  - Read profile configuration

### Export Target Integrations

- [ ] RMI-330: Implement Linear adapter
  - Export requirements as Linear issues
  - Create projects from specs
  - Sync status updates

- [ ] RMI-331: Implement Jira adapter
  - Export requirements as Jira epics/stories
  - Map priorities and labels
  - Create project boards

- [ ] RMI-332: Implement Notion adapter
  - Export specs to Notion pages
  - Create linked databases
  - Sync bidirectionally (optional)

- [ ] RMI-333: Implement Confluence adapter
  - Export specs to Confluence pages
  - Create space structure
  - Link requirements to pages

### Spec Versioning

- [ ] RMI-340: Implement spec version tracking
  - Track spec versions with git-like history
  - Store version metadata in multispec.yaml
  - Immutable version snapshots

- [ ] RMI-341: Implement `multispec diff <spec> [version]`
  - Compare current spec with previous version
  - Show changes by section
  - Highlight requirement changes

- [ ] RMI-342: Implement `multispec history <spec>`
  - Show version history for spec
  - Display change summaries
  - Link to full diffs

- [ ] RMI-343: Implement `multispec revert <spec> <version>`
  - Restore spec to previous version
  - Create new version for revert
  - Preserve audit trail

### Cross-Project Analysis

- [ ] RMI-350: Implement `multispec search <query>`
  - Full-text search across all projects
  - Filter by spec type, project, date
  - Return ranked results

- [ ] RMI-351: Implement requirements reuse tracking
  - Detect similar requirements across projects
  - Suggest reuse opportunities
  - Track requirement lineage

- [ ] RMI-352: Implement pattern detection
  - Identify common patterns across specs
  - Suggest templates from patterns
  - Generate pattern reports

### Real-time Collaboration

- [ ] RMI-360: Implement WebSocket server
  - Real-time spec editing
  - Multiple concurrent editors
  - Operational transformation

- [ ] RMI-361: Implement presence indicators
  - Show who is editing which spec
  - Cursor positions
  - Edit activity feed

- [ ] RMI-362: Implement conflict resolution
  - Detect concurrent edits
  - Merge non-conflicting changes
  - Prompt for conflict resolution

### CI/CD Integration

- [ ] RMI-370: Create GitHub Actions workflows
  - `multispec-lint.yml` - Validate on PR
  - `multispec-eval.yml` - Evaluate changed specs
  - `multispec-status.yml` - Post status comment

- [ ] RMI-371: Create pre-commit hooks
  - `pre-commit-lint` - Run lint before commit
  - `pre-commit-format` - Format specs
  - Integration with pre-commit framework

- [ ] RMI-372: Implement PR comment integration
  - Post eval results as PR comments
  - Show status badge in PR
  - Link to detailed report

- [ ] RMI-373: Create GitLab CI templates
  - `.gitlab-ci.yml` templates
  - Parallel evaluation jobs
  - Artifact publishing

---

## Phase 11: Context Sources / Grounding (v0.4.0)

Aggregate context from multiple sources to ground spec synthesis in reality.

**Project Spec:** [docs/specs/context-sources/](context-sources/)

**Marketing Name:** Grounding

### Context Source Interface

- [x] RMI-400: Define Source interface (`pkg/context/`)
  - `Source` interface with `Name()`, `Type()`, `Fetch()`
  - `ContextData` unified data structure
  - `AggregatedContext` combined results
  - Source types: git, graphize, mcp, file

- [x] RMI-401: Implement Aggregator
  - Concurrent fetching from multiple sources
  - Caching with configurable TTL
  - Error handling and partial results

- [x] RMI-402: Configuration schema in multispec.yaml
  - `context.repositories[]` - git repo configs
  - `context.graphize[]` - graphize graph paths
  - `context.mcp_servers{}` - MCP server configs
  - `context.files[]` - local file configs

### Git Repository Analysis

- [x] RMI-410: Implement GitSource (`pkg/context/git/`)
  - Structure analysis (directory tree)
  - Dependency extraction (go.mod, package.json, etc.)
  - API schema detection (OpenAPI, GraphQL, Proto)
  - README and documentation extraction
  - Language statistics (LOC by language)

- [x] RMI-411: Support remote repositories
  - Clone via URL with sparse checkout
  - Branch selection
  - Shallow clone for performance

### Graphize Integration

- [x] RMI-420: Implement GraphizeSource (`pkg/context/graphize/`)
  - Load graphs from .graphize/ directories
  - Extract nodes: requirement, decision, constraint, user_story
  - Extract edges: traces_to, derived_from, depends_on
  - Traceability statistics

- [x] RMI-421: Auto-detect graphize in git repos
  - `graphize: auto` config option
  - Discover .graphize/ in repo root

### MCP Client

- [x] RMI-430: Implement MCP client (`pkg/context/mcp/`)
  - Subprocess management for MCP servers
  - JSON-RPC protocol implementation
  - Tool call interface

- [x] RMI-431: Jira integration
  - Fetch issues by JQL
  - Extract epics, stories, tasks
  - Include descriptions, status, labels

- [x] RMI-432: Confluence integration
  - Fetch pages by space/label
  - Extract page content
  - Include metadata

- [x] RMI-433: Additional MCP servers
  - Google Docs
  - Office 365
  - Aha
  - Productboard

### CLI Commands

- [x] RMI-440: Implement `multispec context` command group
  - `context gather` - fetch from all sources
  - `context show` - display aggregated context
  - `context refresh` - clear cache and re-fetch
  - `context snapshot` - save to JSON file

- [x] RMI-441: Add `--with-context` flag to synthesize
  - Load context before synthesis
  - Pass to Synthesizer
  - Include in prompts

- [ ] RMI-442: Add `--with-context` flag to align
  - Compare spec against codebase context
  - Detect drift and unimplemented features
  - Generate current-truth.md

- [ ] RMI-443: Add `--context-file` flag
  - Load context from snapshot file
  - For CI reproducibility

### Context-Aware Synthesis

- [x] RMI-450: Update Synthesizer for context
  - `SynthesizeWithContext()` method
  - Context-aware prompt building
  - Include code structure, APIs, dependencies
  - Include graphize traceability

- [x] RMI-451: Context-aware TRD synthesis
  - Reference actual codebase structure
  - Include existing API contracts
  - Trace to graphize requirements

- [x] RMI-452: Context-aware IRD synthesis
  - Reference actual infrastructure
  - Include deployment configs
  - Trace to TRD architecture

### Caching and Snapshots

- [x] RMI-460: Implement context cache
  - In-memory cache with TTL
  - Invalidation on config change

- [x] RMI-461: Implement context snapshots
  - JSON serialization of AggregatedContext
  - Load from file for offline/CI use
  - Diff between snapshots

### Documentation

- [x] RMI-470: Context sources user guide
  - Configuration reference
  - Git repo setup
  - MCP server configuration
  - Graphize integration

- [x] RMI-471: Context sources API documentation
  - Source interface
  - Writing custom sources
  - Extending MCP integrations
