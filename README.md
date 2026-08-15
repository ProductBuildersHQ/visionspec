# Specification Workflow Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/ProductBuildersHQ/specification-workflow-spec/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/ProductBuildersHQ/specification-workflow-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/ProductBuildersHQ/specification-workflow-spec
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://productbuildershq.com/specification-workflow-spec
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=ProductBuildersHQ%2Fspecification-workflow-spec
 [loc-svg]: https://tokei.rs/b1/github/ProductBuildersHQ/specification-workflow-spec
 [repo-url]: https://github.com/ProductBuildersHQ/specification-workflow-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/ProductBuildersHQ/specification-workflow-spec/blob/main/LICENSE

A formal specification for defining product specification workflows.

## Overview

`specification-workflow-spec` provides standardized types for defining:

- **Spec Types** - Registry of specification document types (PRD, MRD, Press Release, FAQ, 6-Pager, etc.)
- **Workflows** - Methodology configurations bundling spec requirements, synthesis rules, and evaluation criteria
- **Templates** - Document structure definitions with required/optional sections and embedded content
- **Rubrics** - LLM-as-Judge evaluation criteria using structured-evaluation's `rubric.RubricSet`
- **Synthesis Rules** - Dependency graphs for generating specs from other specs
- **Phase Gates** - Approval checkpoints and workflow control

## Architecture

```
┌───────────────────────────────────────────────────────────────────────────┐
│                       Workflow (Methodology Configuration)                │
├───────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  SpecConfig  │  │  Synthesis   │  │  Templates   │  │   Rubrics    │   │
│  │  (required/  │  │  (DAG of     │  │  (document   │  │  (evaluation │   │
│  │   optional)  │  │   sources)   │  │   structure) │  │   criteria)  │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
└───────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                                Execution                                  │
├───────────────────────────────────────────────────────────────────────────┤
│  Phase 1: Discovery  →  Gate  →  Phase 2: Vision  →  Gate  →  Phase 3...  │
│  (MRD)                          (Press, FAQ)                  (PRD, UXD)  │
└───────────────────────────────────────────────────────────────────────────┘
```

## Installation

```bash
go get github.com/ProductBuildersHQ/specification-workflow-spec
```

## Packages

| Package | Description |
|---------|-------------|
| `pkg/spectype` | Spec type registry and category definitions |
| `pkg/workflow` | Workflow configuration (spec requirements, synthesis, execution, evaluation) |
| `pkg/workflows` | Embedded default workflows with loaders (embedded, file, chain, resolving) |
| `pkg/template` | Spec template structure definitions |
| `pkg/synthesis` | Synthesis rule DAG for spec generation |
| `pkg/gate` | Phase gates and approval checkpoints |
| `pkg/layout` | Filesystem layout conventions for spec projects |
| `pkg/diagram` | D2 and Mermaid diagram generation from workflows |
| `pkg/integration` | Descriptor types for external execution-side SDD tools |
| `pkg/integrations` | Embedded default tool integrations (spec-kit, ai-dlc, openspec, kiro) |
| `pkg/pipeline` | Types linking a definition workflow to an execution integration |
| `pkg/pipelines` | Embedded default definition→execution pipelines |
| `schema` | Generated JSON Schema files |

Rubric definitions use [structured-evaluation](https://github.com/plexusone/structured-evaluation)'s canonical `rubric.RubricSet` type.

## Spec Types

The registry defines canonical spec types across methodologies:

### Source Specs (Human-Authored)

| ID | Name | Category | Origins |
|----|------|----------|---------|
| `mrd` | Market Requirements Document | source | enterprise, aws-one-way-door, big-tech-product |
| `prd` | Product Requirements Document | source | startup, enterprise, big-tech |
| `uxd` | User Experience Design | source | design-thinking, big-tech |
| `opportunity-spec` | Opportunity Specification | source | aws-two-way-door, big-tech-feature |
| `hypothesis` | Hypothesis Document | source | lean-startup, 0-1 |
| `shapeup-pitch` | Shape Up Pitch | source | shapeup |
| `ost` | Opportunity Solution Tree | source | continuous-discovery |

### GTM Specs (Synthesized)

| ID | Name | Category | Origins |
|----|------|----------|---------|
| `press` | Press Release | gtm | aws-one-way-door, big-tech |
| `faq` | Frequently Asked Questions | gtm | aws-one-way-door, big-tech |
| `narrative-6p` | Six-Pager Narrative | gtm | aws-one-way-door, big-tech-product |
| `narrative-1p` | One-Pager Executive Summary | gtm | enterprise, big-tech |
| `bmc` | Business Model Canvas | gtm | enterprise, lean-startup |

### Technical Specs (Synthesized)

| ID | Name | Category | Origins |
|----|------|----------|---------|
| `trd` | Technical Requirements Document | technical | enterprise, google, big-tech |
| `tpd` | Test Plan Document | technical | enterprise, big-tech |
| `ird` | Infrastructure Requirements Document | technical | enterprise, big-tech-product |

### Execution Specs

| ID | Name | Category | Origins |
|----|------|----------|---------|
| `plan` | Implementation Plan | execution | pbhq-lite |
| `roadmap` | Roadmap | execution | pbhq-lite |
| `spec` | Reconciled Specification | output | enterprise |

See `pkg/spectype/spectype.go` for the full registry.

## Workflows

Workflows bundle spec requirements, synthesis rules, templates, and rubrics for
specific methodologies. Default workflows (aws-one-way-door, big-tech-feature,
lean-startup, etc.) are embedded and load with no filesystem access:

```go
import "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"

// Load with inheritance resolution (aws-two-way-door extends enterprise)
w, err := workflows.DefaultLoader().Load("aws-two-way-door")
if err != nil {
    // handle error
}

w.Workflow.Name              // "aws-two-way-door"
w.Workflow.RequiredSpecs()   // required spec type IDs
w.Templates["press"].Content // raw markdown template
w.Rubrics["press"].Categories // structured-evaluation rubric categories
```

Loaders compose for customization:

```go
// Organization overrides from a directory, falling back to embedded defaults
loader := workflows.NewResolvingLoader(workflows.NewChainLoader(
    workflows.NewFileLoader("./custom-workflows"),
    workflows.DefaultLoader(),
))
```

### Inheritance and Provenance

Workflows inherit via `extends:` (e.g. both AWS door profiles extend
`enterprise`), with child entries overriding parent entries per key. Beyond
implicit inheritance, a spec can declare exactly where its template or rubric
comes from:

```yaml
spec_config:
  mrd:
    required: false
    template: {from: enterprise}   # resolved from the enterprise workflow
    rubric: enterprise             # bare-string shorthand for the same
  press:
    template: local                # this workflow's own templates/ dir
```

Declared provenance is **loader-enforced**: a source that does not actually
provide the file is a load-time error, and a source whose resolution path leads
back to the declaring workflow (e.g. naming a descendant) fails as a circular
reference instead of recursing. Point sources at ancestors or unrelated
workflows; descendant-owned content must be copied, not referenced.

Two library-wide guarantees are enforced by tests:

- **Every declared spec resolves.** Each spec in a workflow's `spec_config` —
  required *or* optional — must resolve to both a template and a rubric.
  "Optional" means optional to *use* in a workflow run, not optional for the
  library to support.
- **Every synthesis source is declared.** A synthesis rule may only consume
  specs the workflow declares in its `spec_config`.

### Layered Rubrics

Rubrics separate three kinds of judgment via a per-category `class`, so an
LLM-as-Judge evaluation can distinguish cultural fit from contract quality:

| Class | Meaning | Blocking? |
|-------|---------|-----------|
| `leadership_principle` | Methodology/culture judgment (e.g. customer obsession, frugality) | Never |
| `specification_quality` | Is the document a complete, testable contract? | Often |
| `implementation_readiness` | Can a team build from this without guessing? | Often |

Categories also carry `blocking` (a failing blocking category fails the
document regardless of weighted score) and an `evaluation` mode
(`deterministic`, `semantic`, or `human`). The hardened workflow families
(`enterprise`, `aws-one-way-door`, `aws-two-way-door`) ship fully layered
PRD/TRD/TPD/IRD/UXD rubric sets, guarded by regression tests.

### Domain Content Sync

Roadmap-domain spec content (MRD, OpportunitySpec, V2MOM, BMC, …) is owned
upstream by [prism-roadmap](https://github.com/grokify/prism-roadmap) and
synced into the embedded library by `tools/prism-sync` (a nested Go module).
Synced files carry provenance headers; `prism-sync -check` detects drift.

## Definition vs. Execution

Workflows describe the **definition** side of spec-driven development — the
methodologies that produce *what to build* (PRD, six-pager, TRD). External tools
such as GitHub Spec-Kit and AWS AI-DLC own the **execution** side — consuming
those specs to produce code. These are the "Definition" and "Execution" halves
surfaced in the VisionStudio viewer.

- **Integrations** (`pkg/integration`, `pkg/integrations`) are declarative
  descriptors of external tools: how to detect a project on disk, its artifacts
  and lifecycle, how to compute status from those artifacts, and where its
  inputs and outputs connect. They hold **no scanning or execution logic** —
  consumers implement detection against the contract.
- **Pipelines** (`pkg/pipeline`, `pkg/pipelines`) link a definition workflow to
  an execution integration, declaring how the workflow's output specs feed the
  tool (e.g., `aws-one-way-door` → `ai-dlc`).

Default execution integrations: **spec-kit**, **ai-dlc**, **openspec**, **kiro**.

```go
import (
    "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/integrations"
    "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/pipelines"
)

// An execution tool descriptor: detection, artifacts, lifecycle, status.
in, _ := integrations.Get("spec-kit")
in.Detection.RootMarkers        // [".specify/", ...] — how to recognize a project
in.Status.Method                // "task-checkboxes" — how a viewer reads progress
in.Artifacts[0].SpecType        // maps a tool file to a spec type ID where applicable

// A definition→execution handoff: aws-one-way-door's specs feed AI-DLC.
p, _ := pipelines.Get("aws-one-way-door-to-ai-dlc")
p.Definition.Workflow           // "aws-one-way-door"
p.Execution.Integration         // "ai-dlc"
p.Handoffs                      // per-spec-type mappings across the seam
```

## Diagram Generation

Generate D2 or Mermaid diagrams from workflows:

```go
import "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/diagram"

// Generate D2 diagram
opts := diagram.DefaultOptions()
opts.Title = "AWS Product Flow"
d2, _ := diagram.Generate(w.Workflow, diagram.FormatD2, opts)

// Generate Mermaid diagram
mermaid, _ := diagram.Generate(w.Workflow, diagram.FormatMermaid, opts)
```

Output formats:

- **D2** - [D2 language](https://d2lang.com/) for SVG generation via `d2` CLI
- **Mermaid** - [Mermaid](https://mermaid.js.org/) for embedding in Markdown

## Schema Generation

JSON Schema files are generated from Go types:

```bash
go generate ./schema/...
```

This produces:

- `schema/spectype.schema.json`
- `schema/workflow.schema.json`
- `schema/template.schema.json`
- `schema/synthesis.schema.json`
- `schema/gate.schema.json`
- `schema/layout.schema.json`

## Ecosystem

This repository is the contract layer of the ProductBuildersHQ spec stack
(`visionstudio → visionspec → specification-workflow-spec`): it defines
workflow types, schemas, and the embedded default workflow library, and holds
no execution logic. The layers above act on it.

| Project | Role |
|---------|------|
| [visionspec](https://github.com/ProductBuildersHQ/visionspec) | The engine: CLI, MCP server, and importable SDK executing these workflows (scaffolding, LLM synthesis, LLM-as-Judge evaluation, lint/drift/status) |
| [visionstudio](https://github.com/ProductBuildersHQ/visionstudio) | The studio: LLM-powered app loading workflow data from this library directly and executing via visionspec |
| [structured-evaluation](https://github.com/plexusone/structured-evaluation) | Canonical rubric and evaluation-report types; rubrics here are its `rubric.RubricSet` |
| [multi-agent-spec](https://github.com/plexusone/multi-agent-spec) | Multi-agent system definitions |

## License

MIT
