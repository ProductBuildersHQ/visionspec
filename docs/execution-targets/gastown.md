# GasTown

GasTown is a multi-agent orchestration system using formulas and beads for DAG-based task execution.

## Overview

GasTown provides:

- **DAG-based execution**: Directed acyclic graph task dependencies
- **Bead abstraction**: Modular, reusable task units
- **Formula types**: Convoy, workflow, and expansion patterns
- **Multi-agent support**: Distribute work across specialized agents
- **Rig integration**: Configurable execution environments

## When to Use GasTown

GasTown is ideal for:

- ✅ Complex multi-agent orchestration
- ✅ Projects with intricate task dependencies
- ✅ Teams running multiple AI agents in parallel
- ✅ Large-scale code generation projects
- ✅ Systems requiring task reusability (beads)

## Integration with VisionSpec

### The Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                     VISIONSPEC                                   │
│                                                                  │
│  MRD → Press → FAQ → PRD → TRD                                  │
│                  ↓                                               │
│              spec.md                                             │
│                  ↓                                               │
│         visionspec export gastown                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     GASTOWN                                      │
│                                                                  │
│  formula.toml     (formula definition)                          │
│  beads/                                                          │
│  ├── foundation.toml                                             │
│  ├── pet-crud.toml                                               │
│  ├── store-crud.toml                                             │
│  └── api-docs.toml                                               │
└─────────────────────────────────────────────────────────────────┘
```

## Core Concepts

### Formulas

A formula defines the overall execution plan:

```toml
# formula.toml
[formula]
name = "petstore-api"
type = "convoy"  # convoy, workflow, or expansion
rig = "go-backend"

[formula.context]
spec = "docs/specs/petstore-api/spec.md"
```

### Beads

Beads are atomic task units:

```toml
# beads/pet-crud.toml
[bead]
id = "pet-crud"
name = "Pet CRUD Operations"
priority = 2

[bead.dependencies]
requires = ["foundation"]

[bead.task]
description = "Implement Pet CRUD API endpoints"
artifacts = [
  "internal/pet/handler.go",
  "internal/pet/repository.go",
  "internal/pet/handler_test.go"
]

[bead.acceptance]
criteria = [
  "All CRUD operations working",
  "Unit tests passing",
  "Integration tests passing"
]
```

### Formula Types

| Type | Description | Use Case |
|------|-------------|----------|
| **Convoy** | Parallel execution with sync points | Multi-agent projects |
| **Workflow** | Sequential with branching | Complex dependencies |
| **Expansion** | Dynamic bead generation | Large-scale generation |

## Export Format

When you run `visionspec export gastown`, VisionSpec creates:

```
gastown/
├── formula.toml
└── beads/
    ├── foundation.toml
    ├── pet-crud.toml
    ├── store-crud.toml
    ├── user-auth.toml
    └── api-docs.toml
```

### formula.toml

```toml
[formula]
name = "petstore-api"
type = "convoy"
rig = "default"
version = "1.0.0"

[formula.context]
spec = "docs/specs/petstore-api/spec.md"
mrd = "docs/specs/petstore-api/source/mrd.md"
trd = "docs/specs/petstore-api/technical/trd.md"

[formula.execution]
parallel_beads = 3
timeout_minutes = 60

[formula.beads]
order = [
  "foundation",
  "pet-crud",
  "store-crud",
  "user-auth",
  "api-docs"
]
```

### Bead Structure

Each bead in `beads/` follows this structure:

```toml
[bead]
id = "foundation"
name = "Project Foundation"
priority = 1

[bead.dependencies]
requires = []  # No dependencies - runs first

[bead.task]
description = """
Set up the project foundation including:
- Go module initialization
- Directory structure
- Database connection
- Configuration management
"""

artifacts = [
  "go.mod",
  "go.sum",
  "cmd/server/main.go",
  "internal/config/config.go",
  "internal/database/connection.go"
]

[bead.acceptance]
criteria = [
  "Go module initializes successfully",
  "Server starts without errors",
  "Database connection established"
]

[bead.context]
from_spec = ["architecture", "technology_stack"]
```

## Complete Workflow

### Step 1: Create Specifications

```bash
visionspec init petstore-api --profile enterprise

# Create comprehensive specs for multi-agent work
visionspec create mrd -p petstore-api
visionspec create prd -p petstore-api
visionspec synthesize trd -p petstore-api
visionspec synthesize tpd -p petstore-api

visionspec reconcile -p petstore-api
```

### Step 2: Export to GasTown

```bash
visionspec export gastown -p petstore-api
```

Output:

```
⋯ Exporting to gastown...
✓ Exported to GasTown format
  Output: gastown/
  Files:
    - formula.toml (convoy, 5 beads)
    - beads/foundation.toml
    - beads/pet-crud.toml
    - beads/store-crud.toml
    - beads/user-auth.toml
    - beads/api-docs.toml
```

### Step 3: Execute with GasTown

Using GasTown CLI:

```bash
# Initialize the formula
gastown init gastown/formula.toml

# Execute all beads
gastown run --formula gastown/formula.toml

# Or execute specific bead
gastown run --bead pet-crud
```

With AI agents:

```
Using GasTown, execute the formula in gastown/formula.toml
Start with the foundation bead
```

### Step 4: Monitor Progress

```bash
# View formula status
gastown status --formula gastown/formula.toml

# View bead status
gastown status --bead pet-crud
```

## Configuration

Configure GasTown export in `visionspec.yaml`:

```yaml
targets:
  gastown:
    enabled: true
    output_dir: gastown

    # Formula configuration
    formula_type: convoy  # convoy, workflow, expansion
    rig: go-backend

    # Bead generation
    bead_granularity: medium  # fine, medium, coarse
    max_parallel_beads: 3

    # Dependencies
    auto_dependencies: true
```

## DAG Execution

### Dependency Resolution

GasTown builds a DAG from bead dependencies:

```
          ┌─────────────┐
          │ foundation  │
          └──────┬──────┘
                 │
       ┌─────────┼─────────┐
       ↓         ↓         ↓
┌──────────┐ ┌──────────┐ ┌──────────┐
│ pet-crud │ │store-crud│ │user-auth │
└────┬─────┘ └────┬─────┘ └────┬─────┘
     │            │            │
     └────────────┼────────────┘
                  ↓
          ┌───────────┐
          │ api-docs  │
          └───────────┘
```

### Execution Order

1. **Level 0**: `foundation` (no dependencies)
2. **Level 1**: `pet-crud`, `store-crud`, `user-auth` (parallel)
3. **Level 2**: `api-docs` (depends on all Level 1)

## Formula Types Explained

### Convoy

Parallel execution with synchronization:

```toml
[formula]
type = "convoy"

# Multiple agents work in parallel
# Sync points between dependency levels
# Best for: Multi-agent, parallelizable work
```

### Workflow

Sequential with conditional branching:

```toml
[formula]
type = "workflow"

# Strict sequential execution
# Supports conditional paths
# Best for: Complex dependencies, conditionals
```

### Expansion

Dynamic bead generation:

```toml
[formula]
type = "expansion"

# Beads can spawn sub-beads
# Dynamic work discovery
# Best for: Large-scale generation, exploration
```

## Best Practices

### 1. Design Atomic Beads

Each bead should be independently executable:

```toml
# ✅ Good: Self-contained
[bead]
id = "pet-handler"
artifacts = ["internal/pet/handler.go", "internal/pet/handler_test.go"]

# ❌ Avoid: Too coupled
[bead]
id = "all-handlers"
artifacts = ["internal/**/*.go"]
```

### 2. Minimize Dependencies

Fewer dependencies = more parallelism:

```toml
# ✅ Good: Minimal dependencies
[bead.dependencies]
requires = ["foundation"]

# ❌ Avoid: Unnecessary dependencies
[bead.dependencies]
requires = ["foundation", "config", "utils", "models"]
```

### 3. Use Priority Ordering

Guide execution order within parallel sets:

```toml
# Higher priority = runs first when parallelism is limited
[bead]
id = "critical-path"
priority = 1  # Highest

[bead]
id = "nice-to-have"
priority = 10  # Lower
```

### 4. Include Context References

Link beads to relevant spec sections:

```toml
[bead.context]
from_spec = ["requirements.pet_management", "architecture.api_layer"]
from_trd = ["api_design", "data_models"]
```

## Troubleshooting

### "Circular dependency detected"

Check bead dependencies for cycles:

```bash
# Visualize dependencies
gastown graph --formula gastown/formula.toml

# Look for A → B → C → A patterns
```

### "Bead execution failed"

Check bead artifacts and criteria:

```bash
# View bead details
cat gastown/beads/pet-crud.toml

# Verify expected artifacts exist
ls -la internal/pet/
```

### "Rig not found"

Define or use default rig:

```toml
[formula]
rig = "default"  # Use default rig
```

## See Also

- [GasCity](gascity.md) - Higher-level agent orchestration
- [Choosing a Target](choosing-a-target.md) - Compare with other targets
- [CLI: export](../cli/export.md) - Export command reference
