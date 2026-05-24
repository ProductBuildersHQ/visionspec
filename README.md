# MultiSpec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/multispec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/multispec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/multispec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/multispec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/multispec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/multispec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/multispec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/multispec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/multispec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/multispec
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/multispec
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fmultispec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/multispec
 [repo-url]: https://github.com/plexusone/multispec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/multispec/blob/main/LICENSE

Multi-domain specification orchestration for humans and AI agents.

## Overview

MultiSpec bridges the gap between organizational intent (MRD, PRD, UXD) and executable specifications for AI coding agents. It provides:

- ✍️ **Domain-specific authoring** - Separate specs for PM, UX, Engineering
- 📣 **GTM synthesis** - LLM-generated press releases, FAQs, narratives (Working Backwards)
- ⚙️ **Technical synthesis** - LLM-generated TRD, IRD from source specs
- 📊 **Structured evaluation** - Per-domain LLM judges with customizable rubrics
- 🔄 **Reconciliation** - Conflict detection and tradeoff resolution
- 📦 **Target adapters** - Export to SpecKit, GSD, GasTown, GasCity, OpenSpec

## Installation

```bash
go install github.com/plexusone/multispec/cmd/multispec@v0.3.0
```

## Quick Start

```bash
# Initialize a new project
multispec init user-onboarding

# Validate project structure
multispec lint

# Check project status
multispec status
multispec status --format json
multispec status --format html > status.html
```

## Directory Structure

```
docs/specs/
├── CONSTITUTION.md                    # Repo-level governance
├── ROADMAP.md                         # Cross-project priorities
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
    ├── eval/                          # All evaluations
    │   ├── mrd.eval.json
    │   ├── prd.eval.json
    │   └── ...
    ├── .graphize/                     # Requirement graph
    ├── spec.md                        # Reconciled execution spec
    ├── current-truth.md               # Post-ship state
    ├── status.html                    # Readiness report
    ├── index.md                       # MkDocs page
    └── multispec.yaml                 # Configuration
```

## Document Lifecycle

```
HUMAN-AUTHORED (Source)
  mrd.md → prd.md → uxd.md
           ↓
LLM-GENERATED (GTM) ← Working Backwards
  press.md → faq.md → narrative.md
           ↓
LLM-GENERATED (Technical)
  trd.md → ird.md
           ↓
RECONCILIATION
  All approved specs → spec.md
           ↓
TARGET EXPORT
  spec.md → SpecKit | GSD | GasTown | GasCity | OpenSpec
           ↓
POST-SHIP ALIGNMENT
  spec.md + reality → current-truth.md
```

## CLI Commands

| Command | Description | Status |
|---------|-------------|--------|
| `init <project>` | Initialize a new project with standard directory structure | Implemented |
| `lint [project]` | Validate directory structure and naming conventions | Implemented |
| `status` | Show project status and readiness gates | Implemented |
| `targets` | List available export targets | Implemented |
| `eval [type]` | Evaluate specs using LLM judges | Implemented |
| `synthesize <type>` | Generate specs from source docs | Implemented |
| `reconcile` | Generate unified execution spec | Implemented |
| `approve <type>` | Approve a spec for reconciliation | Implemented |
| `export <target>` | Export to target execution system | Implemented (SpecKit) |
| `graph <cmd>` | Manage requirement graphs (extract, export, query) | Implemented |
| `profiles <cmd>` | Manage configuration profiles (list, show) | Implemented |
| `serve` | Start MCP server for AI integration | Implemented |

## Status Command

The `status` command shows project readiness with multiple output formats:

```bash
# Terminal output with readiness gates
multispec status -p myproject

# JSON for programmatic use
multispec status -p myproject --format json

# HTML report with traffic light indicator
multispec status -p myproject --format html > status.html

# Markdown for embedding in docs
multispec status -p myproject --format markdown

# CI mode - exits non-zero if not ready
multispec status -p myproject --ci
```

### Readiness Gates

| Gate | Description |
|------|-------------|
| Required specs present | All required source specs (mrd, prd, uxd, trd) exist |
| Evaluations passing | No blocking evaluation findings |
| Approvals obtained | All required specs have approvals |
| Execution spec generated | `spec.md` has been created |

## MCP Server

MultiSpec includes an MCP (Model Context Protocol) server for integration with AI coding assistants like Claude Code and Kiro CLI.

```bash
# Run MCP server directly
multispec-mcp
```

### MCP Tools

| Tool | Status | Description |
|------|--------|-------------|
| `list_projects` | Implemented | List all multispec projects |
| `get_project_status` | Implemented | Get project readiness status |
| `start_draft` | Implemented | Initialize a new draft |
| `update_draft` | Implemented | Save draft content |
| `eval_draft` | Implemented | Evaluate draft against rubric |
| `finalize_draft` | Implemented | Promote draft to final spec |
| `get_draft` | Implemented | Retrieve current draft |
| `discard_draft` | Implemented | Delete a draft |
| `get_spec` | Implemented | Get specification content |
| `get_eval` | Implemented | Get evaluation results |
| `run_eval` | Implemented | Run evaluation against rubric |
| `synthesize` | Implemented | Generate a spec |
| `reconcile` | Implemented | Generate execution spec |
| `approve` | Implemented | Approve a specification |
| `export` | Implemented | Export to target system |

## Export Targets

| Target | Description |
|--------|-------------|
| `speckit` | GitHub Spec-Kit format |
| `gsd` | Get Shit Done (PLAN.md, STATE.md) |
| `gastown` | GasTown formulas and beads |
| `gascity` | GasCity city.toml configuration |
| `openspec` | OpenSpec portable format (future) |

## Configuration Profiles

Profiles define which specs are required for different product lifecycle stages:

```bash
# List available profiles
multispec profiles list

# Show profile details
multispec profiles show startup

# Initialize with a profile
multispec init my-project --profile startup
```

| Profile | Required Specs | Use Case |
|---------|---------------|----------|
| `0-1` | hypothesis | Idea validation phase |
| `startup` | prd | Pre-PMF startups |
| `growth` | prd, uxd, faq | 1-N scaling phase |
| `enterprise` | mrd, prd, uxd, trd, press, faq, spec | Post-PMF enterprises |

## CLI as Library

Organizations can build custom CLI tools using multispec as a library:

```go
import (
    "github.com/plexusone/multispec/pkg/cli"
    "github.com/plexusone/multispec/pkg/profiles"
)

func main() {
    root := &cobra.Command{Use: "org-spec"}

    // Load a profile
    profile, _ := profiles.DefaultLoader().Load("enterprise")
    cfg := cli.ConfigFromProfile(profile)

    // Add multispec commands
    cli.AddCommandsTo(root, cfg)

    root.Execute()
}
```

See `examples/` for complete organization CLI examples.

## Dependencies

- [structured-evaluation](https://github.com/plexusone/structured-evaluation) - Rubric and evaluation types
- [graphize](https://github.com/plexusone/graphize) - Requirement graph extraction

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint

# Install locally
make install
```

## Project Status

See [ROADMAP.md](docs/specs/ROADMAP.md) for detailed implementation status and [CHANGELOG.md](CHANGELOG.md) for release history.

**Current Version:** v0.3.0

| Phase | Status |
|-------|--------|
| Phase 0: Foundation | Complete |
| Phase 1: Directory Structure | Complete |
| Phase 2: Evaluation Engine | Complete |
| Phase 3: GTM & Technical Synthesis | Complete |
| Phase 4: Reconciliation Engine | Complete |
| Phase 5: Target Adapters | Partial (SpecKit) |
| Phase 6: Claude Code Integration | Complete |
| Phase 7: Graphize Integration | Complete |
| Phase 8: Advanced Features | Not Started |
| Phase 9: Composability | Complete |
| Phase 10: Platform Enhancements | Not Started |

## License

MIT
