# graph

Manage requirement graphs using Graphize integration.

## Overview

The `graph` command extracts requirement graphs from spec files, enabling traceability analysis and visualization. It integrates with [graphize](https://github.com/plexusone/graphize) for graph operations.

## Subcommands

### extract

Extract a requirement graph from all specs in the project.

```bash
visionspec graph extract
```

**Output:**

- Creates `.graphize/spec-graph.json` in the project directory
- Prints node count summary by type

**Example:**

```
Extracted graph with 45 nodes and 32 edges
Saved to: docs/specs/user-onboarding/.graphize/spec-graph.json

Node types:
  spec: 4
  section: 18
  requirement: 12
  user_story: 6
  constraint: 3
  decision: 2
```

### export

Export the graph to various formats for visualization or analysis.

```bash
visionspec graph export [--format FORMAT] [--output DIR]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `html` | Export format: `html`, `graphml`, `json` |
| `--output` | `.graphize` | Output directory |

**Formats:**

| Format | Output File | Description |
|--------|-------------|-------------|
| `html` | `graph.html` | Interactive HTML visualization |
| `graphml` | `spec-graph.graphml` | GraphML for graph tools (yEd, Gephi) |
| `json` | `spec-graph.json` | Raw graph data |

**Examples:**

```bash
# Export as interactive HTML
visionspec graph export --format html

# Export as GraphML for yEd
visionspec graph export --format graphml

# Export to custom directory
visionspec graph export --format html --output ./reports
```

### query

Query nodes in the graph by type or spec.

```bash
visionspec graph query [--type TYPE] [--spec SPEC]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--type` | Filter by node type (requirement, user_story, constraint, decision) |
| `--spec` | Filter by spec type (mrd, prd, uxd, trd) |

**Examples:**

```bash
# List all requirements
visionspec graph query --type requirement

# List all PRD nodes
visionspec graph query --spec prd

# List constraints from TRD
visionspec graph query --type constraint --spec trd
```

**Output:**

```
Found 12 nodes

[requirement] The system shall support SSO login
  ID: prd_req_req_prd_1
  Spec: prd
  Text: The system shall support SSO login via SAML 2.0...

[requirement] Users must be able to reset passwords
  ID: prd_req_req_prd_2
  Spec: prd
  Text: Users must be able to reset passwords via email...
```

## Node Types

The graph extractor identifies the following node types from specs:

| Type | Source | Pattern |
|------|--------|---------|
| `requirement` | PRD | Statements with "shall", "must", "should" |
| `user_story` | PRD | "As a... I want... so that..." |
| `acceptance_criteria` | PRD | "Given/When/Then" patterns |
| `constraint` | MRD, TRD | Sections titled "Constraints", "Limitations" |
| `decision` | TRD | Sections titled "Decisions", "Architecture" |
| `section` | All | Markdown headings |
| `spec` | All | Root node for each spec file |

## Edge Types

Relationships extracted between nodes:

| Type | Description |
|------|-------------|
| `contains` | Section contains requirements/decisions |
| `traces_to` | Requirement traces to technical decision |
| `derived_from` | TRD derived from PRD and MRD |

## Graph Storage

Graphs are stored in the `.graphize/` directory within each project:

```
docs/specs/{project}/
└── .graphize/
    ├── spec-graph.json      # Primary graph data
    ├── spec-graph.graphml   # GraphML export (if generated)
    └── graph.html           # HTML visualization (if generated)
```

## Integration with Graphize

The graph command uses:

- [graphize](https://github.com/plexusone/graphize) v0.3.0 - Graph extraction and export
- [graphfs](https://github.com/plexusone/graphfs) v0.2.0 - Graph data structures

For advanced graph analysis, you can use the graphize CLI directly on the exported JSON:

```bash
# Install graphize
go install github.com/plexusone/graphize/cmd/graphize@latest

# Analyze the spec graph
graphize analyze .graphize/spec-graph.json

# Generate detailed report
graphize report .graphize/spec-graph.json
```
