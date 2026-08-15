# Release v0.2.0

**Release Date:** 2026-08-02

This release makes specification-workflow-spec the canonical home for workflow definitions: the complete default workflow library (24 methodologies) is now embedded and loadable as Go structs with no filesystem access. It also renames the top-level entity from Profile to Workflow and consolidates rubric definitions on structured-evaluation's canonical type.

## Breaking Changes

### `pkg/profile` renamed to `pkg/workflow`

The top-level entity for a methodology configuration is now `Workflow`; "profile" was too generic. The inner phase/gate ordering field is renamed from `Workflow` to `Execution` to avoid stuttering.

| Before | After |
|--------|-------|
| `pkg/profile` | `pkg/workflow` |
| `profile.Profile` | `workflow.Workflow` |
| `Profile.Workflow` | `Workflow.Execution` |
| `schema/profile.schema.json` | `schema/workflow.schema.json` |

### `pkg/rubricdef` removed

`RubricDefinition` duplicated structured-evaluation's `rubric.RubricSet` schema and forced a lossy conversion at the evaluation boundary. Rubric definitions now use `github.com/plexusone/structured-evaluation/rubric.RubricSet` directly — the same type LLM-as-Judge evaluation consumes.

## Features

### Embedded Workflow Library (`pkg/workflows`)

All 24 default workflow definitions, previously maintained in visionspec, are now embedded here as canonical data:

- **Stage-based:** 0-1, startup, growth, enterprise
- **Methodology-based:** aws-product, aws-feature, google, stripe, lean-startup, design-thinking, jtbd, shapeup, continuous-discovery
- **Big Tech syntheses:** big-tech, big-tech-product, big-tech-feature, big-tech-essentials, big-tech-essentials-product, big-tech-essentials-feature
- **Strategic:** v2mom, v2mom-company, v2mom-department, v2mom-team
- **ProductBuildersHQ:** pbhq-lite

Each workflow bundles its configuration, markdown templates, and structured-evaluation rubrics.

### Composable Loaders (`pkg/workflows`)

```go
// Embedded defaults with inheritance resolution
w, err := workflows.DefaultLoader().Load("aws-feature")

// Organization overrides falling back to embedded defaults
loader := workflows.NewResolvingLoader(workflows.NewChainLoader(
    workflows.NewFileLoader("./custom-workflows"),
    workflows.DefaultLoader(),
))
```

- `Loader` interface with embedded, file, `fs.FS`, chain, and resolving implementations
- `ResolvingLoader` resolves `extends` inheritance with circular-reference detection
- `LoadedWorkflow` bundles the parsed `Workflow` with its templates and rubrics

### Template Content (`pkg/template`)

`Template` gains a `Content` field carrying raw markdown, so embedded templates render with no filesystem access.

### Methodology YAML Shorthand (`pkg/workflow`)

Principles and artifacts parse from the shorthand used throughout the workflow library (previously dropped or rejected):

```yaml
methodology:
  principles:
    - customer_obsession: "Start with the customer and work backwards"
  artifacts:
    primary:
      - api_spec: "API Specification - OpenAPI/contract defining all endpoints"
```

- `Principle` accepts shorthand and canonical forms; names derive from IDs
- `Artifacts` accepts flat sequences, category-grouped mappings, and bare strings
- `Methodology` gains a `Source` attribution field

## Migration

```go
// Before (v0.1.0)
import "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/profile"
p, err := profile.ParseYAMLFile("profile.yaml")
phases := p.Workflow.Phases

// After (v0.2.0)
import "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
w, err := workflows.DefaultLoader().Load("aws-feature")
phases := w.Workflow.Execution.Phases
rubric := w.Rubrics["press"] // *rubric.RubricSet (structured-evaluation)
```

## Installation

```bash
go get github.com/ProductBuildersHQ/specification-workflow-spec@v0.2.0
```

## Requirements

- Go 1.24 or later

## Links

- [Repository](https://github.com/ProductBuildersHQ/specification-workflow-spec)
- [Changelog](https://github.com/ProductBuildersHQ/specification-workflow-spec/blob/main/CHANGELOG.md)
- [structured-evaluation](https://github.com/plexusone/structured-evaluation)
