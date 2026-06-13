# VisionSpec Roadmap

Multi-domain specification orchestration for humans and AI agents.

## Vision

VisionSpec bridges the gap between organizational intent (MRD, PRD, UXD) and executable specifications for AI coding agents. It provides:

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
    └── visionspec.yaml                 # Project configuration
```

### Naming Conventions (Enforced)

| Element | Convention | Example |
|---------|------------|---------|
| Project directory | `kebab-case` | `user-onboarding`, `user-onboarding` |
| Spec files | `lowercase.md` | `mrd.md`, `prd.md`, `spec.md` |
| Eval files | `{spec}.eval.json` | `mrd.eval.json`, `press.eval.json` |
| Config file | `visionspec.yaml` | Fixed name |
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

- [x] RMI-001: Initialize Go module (`github.com/ProductBuildersHQ/visionspec`)
- [x] RMI-002: Create CLI skeleton with Cobra (`visionspec` command)
- [x] RMI-003: Define core types package (`pkg/types/`)
- [x] RMI-004: Add configuration loading (`visionspec.yaml`)
- [x] RMI-005: Set up CI (lint, test, build)
  - `.github/workflows/go-ci.yaml` - build and test
  - `.github/workflows/go-lint.yaml` - golangci-lint
  - `.github/workflows/go-sast-codeql.yaml` - security analysis
- [x] RMI-006: Create project README

- [x] RMI-007: Implement `visionspec lint` command
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

- [x] RMI-010: Implement `visionspec init` command
  - Create `docs/specs/{project}/` structure
  - Create `source/`, `gtm/`, `technical/`, `eval/` subdirectories
  - Generate `visionspec.yaml` project config

- [x] RMI-011: Support CONSTITUTION.md at `docs/specs/CONSTITUTION.md`
  - Repo-level governance document
  - Optional org-level at `~/.config/visionspec/CONSTITUTION.md`
  - `pkg/config/config.go` - FindConstitution, LoadConstitution functions
  - Used in synth, reconcile, and export commands

### MkDocs Integration

- [x] RMI-016: Generate `{project}/index.md` for each project
  - Spec overview with status badges
  - Links to all specs (source, gtm, technical)
  - Eval summary (pass/fail counts, open findings)
  - Last updated timestamps
  - `pkg/mkdocs/mkdocs.go` - GenerateProjectIndex, WriteProjectIndex
  - `visionspec docs project` command

- [x] RMI-017: Generate `docs/specs/index.md` (specs landing page)
  - List all projects with status
  - Link to CONSTITUTION.md and ROADMAP.md
  - Cross-project metrics
  - `pkg/mkdocs/mkdocs.go` - GenerateSpecsLanding, WriteSpecsLanding
  - `visionspec docs generate` command

- [x] RMI-018: Generate MkDocs navigation structure
  - Auto-update `mkdocs.yml` nav section
  - Or generate `nav.yml` partial for include
  - Support `mkdocs-awesome-pages-plugin` `.pages` files

- [x] RMI-019: Render eval JSON to markdown for MkDocs
  - `visionspec render-evals {project}`
  - Generate `eval/index.md` with rendered findings
  - Collapsible sections per spec
  - Severity badges and status indicators

### Project Status Report

- [x] RMI-019a: Implement `visionspec status` core logic
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
  - `visionspec status --ci` exits non-zero if not ready
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

- [x] RMI-015: Implement `visionspec create {spec-type}` command
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

- [x] RMI-023c: Implement `visionspec eval {spec-type}` CLI command
  - Load spec and rubric
  - Call LLM judge
  - Write `{spec}.eval.json` output

- [x] RMI-024: Implement `visionspec eval --all` command
  - Evaluate all source specs, GTM docs, and technical docs
  - Generate all `*.eval.json` files
  - Support filtering: `--source`, `--gtm`, `--technical`

- [x] RMI-025: Implement `visionspec render {eval-file}` command
  - Render JSON eval to Markdown for human review
  - Use `structured-evaluation/render/markdown`

- [x] RMI-026: Implement `visionspec status` command
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
  - Support project-level LLM config in visionspec.yaml

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

- [x] RMI-027: Implement `visionspec synthesize press` command
  - Input: MRD + PRD
  - Output: `gtm/press.md` (press release format)
  - Template: Hook → Problem → Solution → Quote → CTA → Benefits
  - Generate PRESS_EVAL.json

- [x] RMI-028: Implement `visionspec synthesize faq` command
  - Input: press.md
  - Output: `gtm/faq.md`
  - Structure: External FAQs + Internal FAQs
  - Challenge claims in press release
  - Generate FAQ_EVAL.json

- [x] RMI-029: Implement `visionspec synthesize narrative` command
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
  - `visionspec synthesize press --eval` generates press.md + press.eval.json
  - Auto-evaluate after generation

### Technical Document Generation

### TRD Generation

- [x] RMI-030: Implement `visionspec synthesize trd` command
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

- [x] RMI-032: Implement `visionspec synthesize ird` command
  - Input: TRD + CONSTITUTION
  - Output: `technical/ird.md`
  - Generate IRD_EVAL.json

- [x] RMI-033: Define IRD template structure
  - Infrastructure requirements
  - Deployment architecture
  - Scaling considerations
  - Operational requirements

### Approval Workflow

- [x] RMI-034: Implement `visionspec approve {spec-type}` command
  - Record approval in `visionspec.yaml`
  - Track approver, timestamp
  - Gate for reconciliation

- [x] RMI-035: Support approval status in `visionspec status`
  - Show pending approvals
  - Show approval history

### Post-Ship Alignment

- [x] RMI-036: Implement `visionspec align` command
  - Input: spec.md + shipped reality (from engineering)
  - Output: `current-truth.md`
  - Detect: ungrounded claims, missed opportunities, drift
  - `pkg/align/align.go` - Aligner, AlignmentResult, Discrepancy types
  - `pkg/align/compare.go` - Comparator for spec vs context analysis
  - `pkg/cli/commands.go` - alignCmd with --with-context, --context-file flags

- [x] RMI-037: Define current-truth.md structure
  - Product summary (current state)
  - Active capabilities table
  - Known boundaries/limitations
  - Source specs and evidence
  - Recent alignment notes
  - `pkg/align/align.go` - AlignmentCoverage, AlignmentSummary types

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

- [x] RMI-042: Implement `visionspec reconcile` command
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

- [x] RMI-052: Implement `visionspec targets` command
  - List available targets
  - Show capabilities

- [x] RMI-053: Implement `visionspec export {target}` command
  - Route to appropriate adapter
  - Support multiple targets: `visionspec export speckit,gsd`

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

### OpenSpec Adapter

- [x] RMI-090: Define OpenSpec export format
  - Portable JSON/YAML structure
  - Agent-agnostic representation
  - `pkg/target/openspec.go` - OpenSpecDocument, OpenSpecFeature, OpenSpecTask, etc.

- [x] RMI-091: Implement OpenSpec adapter
  - Standards-compliant export to JSON or YAML
  - Parses spec.md sections into structured format
  - Supports separate files for features/tasks
  - `pkg/target/openspec.go` - OpenSpecTarget with Export method

---

## Phase 6: Claude Code / Kiro CLI Integration

Seamless integration with AI coding assistant workflows via multi-agent-spec and assistantkit.

### Skill Definitions (multi-agent-spec)

- [x] RMI-098: Add Skill schema to multi-agent-spec
  - `sdk/go/skill.go` - Skill struct with builder methods
  - `schema/skill/skill.schema.json` - JSON Schema
  - Loader functions for skill directories
  - Matches assistantkit canonical type

- [x] RMI-099: Define visionspec skills in multi-agent-spec format
  - `visionspec-status` - Check project readiness
  - `visionspec-lint` - Validate project structure
  - `visionspec-eval` - Run evaluations
  - `visionspec-synthesize` - Generate specs
  - `visionspec-reconcile` - Generate unified spec
  - `visionspec-export` - Export to targets

### Skill Generation (assistantkit)

- [x] RMI-100: Generate Claude Code skills via assistantkit
  - `skills/visionspec-status/SKILL.md`
  - `skills/visionspec-lint/SKILL.md`
  - etc.

- [x] RMI-101: Generate Kiro CLI steering files via assistantkit
  - `steering/visionspec-status.md`
  - `steering/visionspec-lint.md`
  - etc.

### Automation

- [x] RMI-102: Implement `visionspec watch` command
  - File watcher for spec changes
  - Auto-run lint on change
  - Debounce support for rapid changes

- [x] RMI-103: Support git hooks
  - Pre-commit: lint changed spec files
  - Pre-push: evaluate specs and check for blockers
  - `pkg/hooks/hooks.go` - Manager with Install, Uninstall, Status methods
  - `pkg/hooks/templates.go` - Hook script templates
  - `visionspec hooks install/uninstall/status` commands

---

## Phase 7: Graphize Integration

Requirement graph visualization via `github.com/plexusone/graphize`.

### Spec Extractor

- [x] RMI-140: Create spec extractor for graphize
  - `pkg/specgraph/specgraph.go` - SpecExtractor
  - Parse markdown specs (mrd.md, prd.md, trd.md, etc.)
  - Extract requirements, constraints, decisions as nodes
  - Infer relationships as edges

- [x] RMI-141: Define spec node types
  - `requirement` - Functional requirements from PRD
  - `user_story` - User stories from PRD
  - `constraint` - Constraints from CONSTITUTION, TRD
  - `acceptance_criteria` - Testable criteria
  - `decision` - Architectural decisions from TRD
  - `tradeoff` - Explicit tradeoffs from reconciliation
  - `capability` - Current capabilities from CURRENT-TRUTH
  - `section`, `spec` - Structural nodes

- [x] RMI-142: Define spec edge types
  - `traces_to` - Requirement traceability (PRD → TRD)
  - `derived_from` - Synthesis source (TRD → MRD + PRD)
  - `conflicts_with` - Detected conflicts
  - `satisfies` - Implementation satisfies requirement
  - `depends_on` - Requirement dependencies
  - `blocks` - Blocking relationships
  - `supersedes` - Decision replacement
  - `contains` - Section/spec containment

### Graph Storage

- [x] RMI-143: Store spec graph in project directory
  - `.graphize/` directory under `docs/specs/{project}/`
  - Version controlled with project
  - `SaveJSON`/`LoadJSON` for spec-graph.json

- [x] RMI-144: Implement `visionspec graph` commands
  - `visionspec graph extract` - Build graph from specs
  - `visionspec graph query` - Query relationships with filters
  - `visionspec graph export` - Export to HTML/JSON/GraphML

### Traceability Analysis

- [x] RMI-145: Implement traceability reports
  - `ComputeMetrics()` - TraceCoverage, ConflictCount
  - Requirements without TRD coverage
  - TRD tasks without PRD traceability
  - Integrated in `visionspec status` via GraphMetrics

- [x] RMI-146: Conflict detection via graph
  - Query `conflicts_with` edges
  - Highlight in reconciliation
  - Surface in SPEC_EVAL

### Visualization

- [x] RMI-147: Generate spec graph HTML visualization
  - Export via `graphize/pkg/exporters/htmlsite`
  - GraphML export via `graphize/pkg/exporters/graphml`
  - JSON export for custom visualization

- [x] RMI-148: MkDocs graph integration
  - Embed graphize visualization in project index.md
  - Or link to standalone HTML export from graphize
  - Uses graphize pkg/exporters/htmlsite for visualization
  - `pkg/mkdocs/mkdocs.go` - DetectGraphizeOutput, GenerateProjectIndexWithGraph, EmbedGraphizeInIndex

---

## Phase 8: Advanced Features

Future enhancements.

### Multi-Project Support

- [x] RMI-110: Support cross-project dependencies
  - Project references in spec.md
  - Cross-project reconciliation
  - `pkg/deps/deps.go` - DependencyManager, ProjectDep types

- [x] RMI-111: Implement `docs/specs/ROADMAP.md` generation
  - Aggregate project statuses
  - Prioritization tracking
  - `pkg/roadmap/roadmap.go` - Generator with template-based output

### Organizational Memory

- [x] RMI-120: Decision log persistence
  - Track tradeoffs across projects
  - Searchable decision history
  - `pkg/decisions/log.go` - DecisionLog, Decision types

- [x] RMI-121: Rationale graphs
  - Link decisions to requirements
  - Impact analysis
  - `pkg/decisions/rationale.go` - RationaleGraph, RationaleNode types

### Analytics

- [x] RMI-130: Evaluation metrics dashboard
  - Spec quality trends
  - Common failure patterns
  - `pkg/metrics/metrics.go` - Collector, EvalMetrics, ProjectMetrics
  - `visionspec metrics` command with text/json/html output

- [x] RMI-131: Reconciliation metrics
  - Conflict frequency
  - Resolution time
  - `pkg/metrics/metrics.go` - ReconcileMetrics type

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
| v0.4.0 | 4, 5, 11 | Reconciliation with conflicts, target adapters (GSD, GasTown, GasCity), Context Sources |
| v0.5.0 | 12 | Methodology profiles (AWS, Google, Stripe, Lean Startup, Design Thinking, JTBD), Working Backwards flow |
| v0.6.0 | 13 | TPD spec type, AWS AI-DLC export, workflow rules for AI assistants |
| v0.7.0 | 14 | Execution integration (status sync, drift detection, test generation, issue export) |
| v1.0.0 | 8, 10 | Production release with full feature set |

---

## Phase 9: Composability (v0.3.0)

Enable organizations (companies, open source projects, non-profits) to compose custom CLI tools with visionspec as a library.

### CLI as Library

- [x] RMI-200: Move CLI commands from `internal/cli` to `pkg/cli`
  - Create `pkg/cli/cli.go` with composition API
  - `AddCommandsTo(root *cobra.Command, cfg *Config)`
  - `Commands(cfg *Config)` for selective command access
  - `DefaultConfig()` for visionspec defaults

- [x] RMI-201: Update `cmd/visionspec/main.go` to use `pkg/cli`
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

- [x] RMI-231: Update `visionspec.yaml` schema
  - Add `specs:` section for per-spec configuration
  - Parse and merge with defaults

- [x] RMI-232: Update `SpecType.IsRequired()` to use config
  - Check project config first, then defaults

### Documentation & Examples

- [x] RMI-240: Create example org CLI (`examples/org-cli/`)
  - Sample CLI importing visionspec as library
  - Custom templates and rubrics
  - Custom spec types

- [x] RMI-241: Add custom profiles documentation (`docs/guides/`)
  - Custom profiles guide (`docs/guides/custom-profiles.md`)
  - CLI profiles reference (`docs/cli/profiles.md`)
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
  - `profiles export <name> <dir>` - Export profile for customization
  - `--profile` flag on init command
  - `--profile-dir` flag for custom profile directories

- [x] RMI-253: Update example CLIs to use profiles
  - `examples/0-1-product/` - Uses "0-1" profile
  - `examples/pre-pmf-startup/` - Uses "startup" profile
  - `examples/1-n-growth/` - Uses "growth" profile
  - `examples/post-pmf-enterprise/` - Uses "enterprise" profile

- [x] RMI-254: Add profile tests (`pkg/profiles/*_test.go`)
  - Test profile loading and inheritance
  - Test profile merging
  - Test template/rubric loader creation

---

## Phase 10: Platform Enhancements

Future enhancements for testing, integrations, and developer experience.

### Testing & Quality

- [x] RMI-300: Add comprehensive profile tests
  - Unit tests for `pkg/profiles/`
  - Profile inheritance testing
  - Loader chain testing

- [x] RMI-301: Add MCP integration tests
  - Test MCP resource handlers (templates, rubrics, profiles)
  - Test resource listing endpoints
  - URI scheme validation

- [x] RMI-302: Add end-to-end authoring workflow tests
  - start_draft → update_draft → eval_draft → finalize_draft
  - Test with real project structure (temp directories)
  - Verify file system operations (draft_test.go, session_test.go)

### Profile CLI Enhancements

- [x] RMI-310: Implement `profiles create <name>` command
  - Interactive profile creation wizard
  - Select base profile to extend
  - Choose required spec types
  - Generate profile.yaml

- [x] RMI-311: Implement `profiles extend <base> <name>` command
  - Create profile extending another
  - Override specific settings
  - Custom templates/rubrics directory

- [x] RMI-312: Implement profile validation
  - `profiles validate <path>` command
  - Check profile.yaml schema
  - Verify referenced templates/rubrics exist

### MCP Resources

- [x] RMI-320: Expose templates as MCP resources
  - `templates://` URI scheme
  - List available templates
  - Read template content

- [x] RMI-321: Expose rubrics as MCP resources
  - `rubrics://` URI scheme
  - List available rubrics
  - Read rubric definitions

- [x] RMI-322: Expose profiles as MCP resources
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

- [x] RMI-340: Implement spec version tracking
  - Track spec versions with git-like history
  - Store version metadata in eval/versions/
  - SHA256 content hashing for change detection

- [x] RMI-341: Implement `visionspec version diff <spec> [version]`
  - Compare current spec with previous version
  - LCS-based diff algorithm
  - Compact and full output modes

- [x] RMI-342: Implement `visionspec version list <spec>`
  - Show version history for spec
  - Display timestamps, hashes, and messages
  - Alias: `visionspec version history`

- [x] RMI-343: Implement `visionspec version revert <spec> <version>`
  - Restore spec to previous version
  - Creates new version for audit trail
  - Custom revert message support

### Cross-Project Analysis

Note: These features use graphize (`github.com/plexusone/graphize`) which provides the underlying graph-based search, reuse, and pattern detection capabilities.

- [x] RMI-350: Implement `visionspec search <query>`
  - Full-text search across all projects
  - Filter by spec type, project, date
  - Return ranked results
  - `pkg/search/search.go` - Searcher, SearchResult, SearchOutput
  - `visionspec search` command with --project, --type, --limit, --regex flags

- [x] RMI-351: Implement requirements reuse tracking
  - Detect similar requirements across projects
  - Suggest reuse opportunities
  - Track requirement lineage
  - `pkg/reuse/reuse.go` - Tracker, ReuseReport, DuplicateGroup, SimilarGroup
  - `visionspec reuse` command

- [x] RMI-352: Implement pattern detection
  - Identify common patterns across specs
  - Suggest templates from patterns
  - Generate pattern reports
  - `pkg/patterns/patterns.go` - Detector, PatternReport, StructuralPattern, ContentPattern, AntiPattern
  - `visionspec patterns` command

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
  - `visionspec-lint.yml` - Validate on PR
  - `visionspec-eval.yml` - Evaluate changed specs
  - `visionspec-status.yml` - Post status comment

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

- [x] RMI-402: Configuration schema in visionspec.yaml
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

- [x] RMI-440: Implement `visionspec context` command group
  - `context gather` - fetch from all sources
  - `context show` - display aggregated context
  - `context refresh` - clear cache and re-fetch
  - `context snapshot` - save to JSON file

- [x] RMI-441: Add `--with-context` flag to synthesize
  - Load context before synthesis
  - Pass to Synthesizer
  - Include in prompts

- [x] RMI-442: Add `--with-context` flag to align
  - Compare spec against codebase context
  - Detect drift and unimplemented features
  - Generate current-truth.md
  - `pkg/cli/commands.go` - alignCmd flags

- [x] RMI-443: Add `--context-file` flag
  - Load context from snapshot file
  - For CI reproducibility
  - `pkg/cli/commands.go` - alignCmd, driftCmd flags

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

---

## Phase 12: Methodology Profiles (v0.5.0)

Comprehensive product management methodology frameworks with templates, rubrics, and Go structs.

### Profile System

- [x] RMI-480: Create profile inheritance system
  - `extends:` field for profile inheritance
  - `abstract: true` for base-only profiles
  - Profile merging with override support

- [x] RMI-481: Create Big Tech composite profile
  - Combined practices from AWS, Google, Stripe, Netflix, Spotify, Meta, Apple, Microsoft
  - ~30 unified principles across all companies
  - `big-tech/profile.yaml` (abstract base)
  - `big-tech-product/profile.yaml` (MRD start)
  - `big-tech-feature/profile.yaml` (OpportunitySpec start)

- [x] RMI-482: Create Shape Up profile (Basecamp)
  - Pitch-based development with appetite not estimates
  - Hill charts for progress tracking
  - Fixed time, variable scope
  - `shapeup/profile.yaml`
  - Templates: shapeup-pitch.md, shapeup-scope.md
  - Rubrics: shapeup-pitch.rubric.yaml

- [x] RMI-483: Create Continuous Discovery profile (Teresa Torres)
  - Weekly touchpoints and story-based interviews
  - Opportunity Solution Trees (OST)
  - Assumption testing by type (desirability, viability, feasibility, usability, ethical)
  - `continuous-discovery/profile.yaml`
  - Templates: discovery-snapshot.md, assumption-map.md, ost.md
  - Rubrics: discovery-snapshot.rubric.yaml, assumption-map.rubric.yaml

### prism-roadmap Canvas Types

Go structs for strategic planning canvases in `github.com/grokify/prism-roadmap`.

- [x] RMI-484: Create Shape Up canvas types (`canvas/shapeup.go`)
  - `ShapeUpPitch` - Problem, appetite, solution, rabbit holes, no-gos
  - `ShapeUpBet` - Betting table decisions
  - `ShapeUpScope` - Hill chart tracking during building
  - Supporting types: `SUProblem`, `SUAppetite`, `SUSolution`, `SURabbitHole`
  - Unit tests: `canvas/shapeup_test.go`

- [x] RMI-485: Create Continuous Discovery canvas types (`canvas/discovery.go`)
  - `DiscoverySnapshot` - Weekly discovery summary
  - `CDInterview`, `CDStory` - Story-based interview data
  - `CDAssumptionTest`, `CDAssumption` - Assumption testing
  - `AssumptionMap` - Risk matrix by type
  - `ExperienceMap` - Customer journey mapping
  - Unit tests: `canvas/discovery_test.go`

- [x] RMI-486: Update Canvas wrapper discriminated union
  - Add `ShapeUpPitch`, `ShapeUpBet`, `ShapeUpScope` to Canvas wrapper
  - Add `DiscoverySnapshot`, `AssumptionMap`, `ExperienceMap` to wrapper
  - Update `CanvasType` enum

- [x] RMI-487: Generate JSON schemas for new canvas types
  - `schema/shapeup-pitch.schema.json`
  - `schema/shapeup-bet.schema.json`
  - `schema/shapeup-scope.schema.json`
  - `schema/discovery-snapshot.schema.json`
  - `schema/assumption-map.schema.json`
  - `schema/experience-map.schema.json`

- [x] RMI-488: Create renderers for new canvas types
  - D2 renderer for Shape Up hill chart
  - D2 renderer for OST tree structure
  - Mermaid renderer for both
  - Markdown table renderer

### Big Tech Profile Enhancement

Enhance Big Tech profile to be "best of all worlds" integrating Shape Up and Continuous Discovery.

- [x] RMI-489: Integrate Shape Up practices into Big Tech profile
  - Add appetite-based scoping as alternative to story points
  - Add pitch-based workflow for major features
  - Add hill chart tracking for execution visibility
  - Add circuit breaker principle

- [x] RMI-490: Integrate Continuous Discovery practices into Big Tech profile
  - Add weekly touchpoints as standard practice
  - Add OST for opportunity mapping
  - Add assumption testing framework (DVFUE matrix)
  - Add story-based interview guidelines

- [x] RMI-491: Create Big Tech "best of all worlds" documentation
  - Update `docs/frameworks/big-tech.md` with integrated practices
  - Add practice selection guidelines by context
  - Add conflict resolution when practices overlap

### Remaining Framework Profiles

- [x] RMI-492: Create Lean Startup profile Go structs
  - `LeanStartupCanvas` - Build-Measure-Learn cycle
  - `MVP`, `LSExperiment`, `Pivot` - MVP and pivot tracking
  - `canvas/leanstartup.go`
  - Renderers: D2, Mermaid, Markdown
  - Schema: `leanstartup.schema.json`

- [x] RMI-493: Create Design Thinking profile Go structs
  - `DesignThinkingCanvas` - Five phases: Empathize, Define, Ideate, Prototype, Test
  - `EmpathyMap`, `DTIdea`, `DTPrototype`, `DTTest` - All phase types
  - `canvas/designthinking.go`
  - Renderers: D2, Mermaid, Markdown
  - Schema: `designthinking.schema.json`

- [x] RMI-494: Create JTBD profile Go structs
  - `JobStatement` - Job-to-be-done definition
  - `OutcomeExpectation` - Success metrics
  - `canvas/jtbd.go` with full JTBD methodology types
  - `canvas/jtbd_test.go` with comprehensive unit tests

### Documentation

- [x] RMI-495: Create framework documentation
  - `docs/frameworks/shapeup.md`
  - `docs/frameworks/continuous-discovery.md`
  - Update `docs/frameworks/index.md` with all frameworks

- [x] RMI-496: Update core workflow documentation
  - Added Shape Up flow (pitch → betting → scope → hill chart → build)
  - Added Continuous Discovery flow (touchpoints → OST → assumptions → testing)
  - Added combined Shape Up + Continuous Discovery workflow
  - `.visionspec-rules/core-workflow.md` - Framework-Specific Flows section

### Profile Template/Rubric Completeness

Ensure all profiles have complete template and rubric coverage for LLM-as-a-Judge evaluation.

- [x] RMI-497: Add missing templates/rubrics to growth profile
  - `templates/trd.md` - Lightweight technical requirements
  - `templates/ird.md` - Lightweight infrastructure requirements
  - `templates/tpd.md` - Lightweight test plan
  - `rubrics/trd.rubric.yaml`, `rubrics/ird.rubric.yaml`, `rubrics/tpd.rubric.yaml`
  - Update profile.yaml to include optional technical specs

- [x] RMI-498: Add missing templates/rubrics to aws-feature profile
  - `templates/prd.md` - Feature-level PRD
  - `templates/trd.md` - Feature-level TRD
  - `templates/uxd.md` - Feature-level UXD
  - `templates/ird.md` - Feature-level IRD
  - `templates/tpd.md` - Feature-level TPD
  - `rubrics/trd.rubric.yaml`, `rubrics/ird.rubric.yaml`, `rubrics/uxd.rubric.yaml`

- [x] RMI-499: Add missing templates/rubrics to big-tech-feature profile
  - Same templates/rubrics as aws-feature
  - Naming updated to "Big Tech Feature"

- [x] RMI-500a: Add missing templates/rubrics to startup profile
  - `templates/uxd.md` - Lightweight UX design
  - `templates/trd.md` - Lightweight technical design
  - `rubrics/uxd.rubric.yaml`, `rubrics/trd.rubric.yaml`
  - Focus on pragmatism and MVP scope

- [x] RMI-501: Add missing templates/rubrics to 0-1 profile
  - `templates/lean-canvas.md` - Lean Canvas for business model
  - `templates/experiment.md` - Experiment card for hypothesis testing
  - `rubrics/lean-canvas.rubric.yaml`, `rubrics/experiment.rubric.yaml`
  - Update profile.yaml to include validation specs

- [x] RMI-508: Add missing templates to big-tech-essentials-feature profile
  - `templates/opportunity-spec.md` - 12-box OpportunitySpec (required by profile)
  - Profile extends big-tech-essentials → enterprise, but enterprise lacks opportunity-spec template

- [x] RMI-509: Add missing templates to big-tech-essentials-product profile
  - `templates/narrative-6p.md` - Amazon-style 6-pager (required by profile)
  - Profile extends big-tech-essentials → enterprise, but enterprise lacks narrative-6p template

### Workflow Diagrams

Visual documentation of profile selection and reconciliation workflows.

- [x] RMI-502: Create profile selection decision tree diagram
  - `docs/diagrams/profile-selection-decision-tree.d2`
  - `docs/diagrams/profile-selection-decision-tree.svg`
  - Decision flow: Scope → Stage → Methodology → Profile
  - Covers all 19 profiles with color-coded recommendations

- [x] RMI-503: Create reconciliation workflow diagram
  - `docs/diagrams/reconciliation-workflow.d2`
  - `docs/diagrams/reconciliation-workflow.svg`
  - Flow: Source Docs → Evaluation → Conflict Detection → Resolution → spec.md → Approval → Export

### prism-roadmap Prioritization Frameworks

RICE Scoring and Kano Model integration with OpportunitySpec.

- [x] RMI-504: Add RICE scoring types to prism-roadmap
  - `prioritization/rice.go` - RICEScore, ImpactLevel, ConfidenceLevel
  - `prioritization/rice_test.go` - Unit tests
  - Formula: Score = (Reach × Impact × Confidence) / Effort

- [x] RMI-505: Add Kano Model types to prism-roadmap
  - `prioritization/kano.go` - KanoCategory, KanoFeature, KanoAnalysis
  - `prioritization/kano_test.go` - Unit tests
  - Categories: Must-Be, Performance, Attractive, Indifferent, Reverse

- [x] RMI-506: Integrate prioritization into OpportunitySpec
  - `canvas/opportunity_spec.go` - Add RICE and Kano fields
  - Helper methods: SetRICE(), SetKano(), IsMustHave(), IsDelighter()
  - GetPrioritizationSummary() for combined analysis

- [x] RMI-507: Add prioritization documentation
  - `docs/canvas/prioritization.md` - RICE and Kano comprehensive guide
  - Update `docs/canvas/opportunity-spec.md` with prioritization section
  - Update multispec `docs/frameworks/opportunity-spec.md`

---

## Phase 13: AI Workflow Orchestration (v0.6.0)

AI assistant integration with workflow rules for Claude Code, Kiro, and Cursor.

### Test Plan Document (TPD)

- [x] RMI-500: Add TPD spec type
  - TPD template with 14 sections (strategy, cases, automation, CI/CD)
  - Synthesis from PRD + TRD + UXD
  - TPD rubrics for enterprise, AWS, and Google profiles

### AWS AI-DLC Export Target

- [x] RMI-510: Implement AIDLC target adapter
  - Generate vision-document.md from MRD + Press
  - Generate technical-environment.md from TRD + IRD + context
  - Generate imported-requirements.md from spec.md
  - `pkg/target/aidlc.go`

### Workflow Rules

- [x] RMI-520: Create workflow rules directory structure
  - `.visionspec-rules/core-workflow.md` main orchestration
  - Phase-specific rules: Discovery, Vision, Experience, Technical, Reconciliation
  - Gate definitions for evaluation and approval

- [x] RMI-521: Add framework selection to workflow
  - 6 methodology flows (AWS, Lean Startup, Design Thinking, JTBD, Google, Stripe)
  - Default to AWS Working Backwards
  - Framework-specific phase rules in `frameworks/` subdirectory

- [x] RMI-522: Implement rules CLI commands
  - `rules list` - List available workflow rules
  - `rules export [dir]` - Export rules to project
  - Embedded rules via go:embed for distribution

### Documentation

- [x] RMI-530: TPD documentation
  - TPD in concepts.md and synthesize.md
  - CONSTITUTION Quality and Testing Requirements section
  - AWS framework doc updated with TPD mapping

- [x] RMI-531: Export targets documentation
  - AIDLC target in export.md and targets.md
  - CLAUDE.md updated with workflow rules reference

---

## Phase 14: Execution Integration (v0.7.0)

Bidirectional integration with AI coding agent execution systems.

### Status Synchronization

- [x] RMI-600: Bidirectional status sync (interface only)
  - `pkg/target/target.go` - Syncer interface, SyncResult, TaskState
  - `visionspec sync <target>` command implemented
  - Note: No targets implement Syncer yet (pending target adapter updates)

- [x] RMI-601: Execution status tracking in MCP
  - MCP tools: `get_execution_status`, `track_requirement`
  - Execution status persistence in `.visionspec/execution-status.json`
  - Requirement progress tracking (pending, in_progress, implemented, blocked)
  - `internal/mcp/server.go` - ExecutionStatus, RequirementStatus types

### Spec Drift Detection

- [x] RMI-610: Implement drift detection
  - `pkg/drift/drift.go` - Detector, DriftReport, DriftItem
  - `pkg/drift/analyze.go` - Analyzer with requirement/implementation extraction
  - `pkg/drift/render.go` - Text, JSON, Markdown renderers
  - `visionspec drift` command with --severity, --format, --ci flags

- [x] RMI-611: Drift resolution workflow
  - MCP tool: `get_resolution_plan` generates prioritized resolution actions
  - MCP prompt: `resolve_drift` for guided drift resolution
  - Categorized resolution strategies (missing_feature, undocumented_code, diverged)
  - `pkg/align/resolution.go` - ResolutionEngine, ResolutionPlan types

### Executable Test Generation

- [x] RMI-620: Generate test cases from TPD
  - `pkg/testgen/testgen.go` - Generator interface, TestCase, ParsedTPD
  - `pkg/testgen/parser.go` - TPD markdown parser
  - `pkg/testgen/go.go` - Go test generator (testing, testify)
  - `pkg/testgen/typescript.go` - TypeScript/Jest generator
  - `pkg/testgen/python.go` - Python/pytest generator
  - `visionspec generate tests` command with --lang, --framework, --group-by

- [x] RMI-621: Test coverage mapping
  - Map TPD test cases to actual tests
  - Track coverage by requirement
  - Report uncovered requirements
  - `pkg/testmap/testmap.go` - Mapper, CoverageReport types

### Issue Export

- [x] RMI-630: GitHub Issues export
  - Export requirements as GitHub issues
  - Create milestones from phases
  - Link issues to specs via labels
  - `pkg/target/github_issues.go` - GitHubIssuesTarget
  - `visionspec export github` command

- [x] RMI-631: Jira export
  - Export requirements as Jira epics/stories
  - Map priorities and labels
  - Create project boards
  - `pkg/target/jira.go` - JiraTarget
  - `visionspec export jira` command

### MCP Execution Tracking

- [x] RMI-640: MCP execution context
  - MCP tool: `get_execution_context` provides spec summary, requirements, guidance
  - Includes project status, spec content, TRD, and codebase context
  - `internal/mcp/server.go` - handleGetExecutionContext

- [x] RMI-641: Execution prompts
  - MCP prompt: `implement_requirement` - Guides implementation of specific requirement
  - MCP prompt: `verify_acceptance` - Guides acceptance criteria verification
  - MCP prompt: `resolve_drift` - Guides drift resolution
  - `internal/mcp/server.go` - registerPrompts, handleImplementRequirementPrompt, etc.
