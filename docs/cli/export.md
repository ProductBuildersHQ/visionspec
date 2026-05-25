# export

Export the reconciled spec to target execution systems.

## Usage

```bash
visionspec export <target> [flags]
```

## Description

The `export` command transforms the unified `spec.md` into formats compatible with AI coding agent execution systems. Each target produces different output files optimized for specific workflows.

## Prerequisites

Before exporting, you must have a reconciled `spec.md`:

```bash
visionspec reconcile
```

## Arguments

| Argument | Description |
|----------|-------------|
| `target` | Target system: `speckit`, `gsd`, `gastown`, `gascity` |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dry-run` | bool | `false` | Show what would be exported without writing |
| `--output` | string | `""` | Output directory (default: `export/<target>`) |

## Available Targets

### speckit

GitHub Spec-Kit format for spec-driven development.

**Output files:**

- `spec.md` - Unified specification
- `plan.md` - Implementation plan with phases
- `tasks.md` - Task breakdown

**Features:**

- Sequential task execution
- Constitution sync to `.specify/memory/`
- Branch numbering support

### gsd

Get Shit Done methodology format.

**Output files:**

- `PLAN.md` - Plan with YAML frontmatter (must_haves, truths, artifacts)
- `STATE.md` - Progress tracking
- `config.json` - GSD configuration

**Features:**

- Wave-based parallel execution
- Artifact verification
- Progress state tracking

### gastown

GasTown multi-agent formula format.

**Output files:**

- `formula.toml` - Formula definition (convoy/workflow/expansion)
- `beads/*.toml` - Individual task beads

**Features:**

- DAG-based execution
- Bead dependencies
- Priority ordering

### gascity

GasCity multi-agent orchestration format.

**Output files:**

- `city.toml` - City configuration with agents and orders

**Features:**

- Agent role definitions
- Order (task) dependencies
- Orchestration modes (orchestrated, autonomous, hybrid)

## Examples

```bash
# Export to SpecKit
visionspec export speckit

# Export to GSD with custom output directory
visionspec export gsd --output ./my-gsd-export

# Dry run to see what would be created
visionspec export gastown --dry-run

# Export to GasCity
visionspec export gascity
```

## Output

```
⋯ Exporting to speckit...
✓ Exported to SpecKit format
  Output: docs/specs/my-project/export/speckit
  Files:
    - spec.md
    - plan.md
    - tasks.md
```

## Configuration

Target-specific configuration in `visionspec.yaml`:

```yaml
targets:
  default: speckit
  speckit:
    enabled: true
    output_dir: export/speckit
    branch_numbering: sequential
  gsd:
    model_profile: balanced
  gastown:
    formula_type: convoy
    rig: my-rig
  gascity:
    city_dir: export/gascity
```

## Constitution Sync

For SpecKit exports, if `CONSTITUTION.md` exists, it's synced to `.specify/memory/constitution.md` for agent memory integration.

## See Also

- [reconcile](reconcile.md) - Generate spec.md before export
- [targets](targets.md) - List available targets
