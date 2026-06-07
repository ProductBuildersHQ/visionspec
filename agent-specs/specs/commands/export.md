---
name: export
description: Export reconciled specification to an execution target format
arguments: [target]
dependencies: [spec-synthesis]
---

# Export Specification

Export the reconciled spec.md to an execution target format.

## Usage

```
/export <target>
```

Where `<target>` is one of:
- `aidlc` - AWS AI-DLC format
- `speckit` - GitHub SpecKit format
- `gsd` - GSD (Get Shit Done) format
- `gastown` - GasTown format
- `gascity` - GasCity format

## Prerequisites

- Reconciled spec.md must exist

## Target Descriptions

| Target | Best For | Output |
|--------|----------|--------|
| aidlc | Enterprise with approval gates | `.aidlc/` |
| speckit | GitHub PR workflows | `.specify/` |
| gsd | Fast parallel execution | `gsd/` |
| gastown | DAG multi-agent | `gastown/` |
| gascity | Role-based agents | `gascity/` |

## Output Examples

### AWS AI-DLC
```
.aidlc/
├── vision-document.md
├── technical-environment.md
└── imported-requirements.md
```

### GitHub SpecKit
```
.specify/
├── spec.md
├── plan.md
├── tasks.md
└── memory/
    └── constitution.md
```

### GSD
```
gsd/
├── PLAN.md
├── STATE.md
└── config.json
```

### GasTown
```
gastown/
├── formula.toml
└── beads/
    ├── foundation.toml
    └── core.toml
```

### GasCity
```
gascity/
├── city.toml
└── districts/
    ├── backend.toml
    └── devops.toml
```

## Output

```
⋯ Exporting to {target}...

✓ Exported to {target} format
  Output: {directory}/
  Files:
    - {file1}
    - {file2}
    - {file3}

Next steps for {target}:
  {Target-specific instructions}
```

## Target-Specific Instructions

### AI-DLC
```
Next steps:
  In Claude Code, run:
    "Using AI-DLC, implement the project based on .aidlc/"
```

### SpecKit
```
Next steps:
  In Claude Code, run:
    "Using SpecKit, execute the plan in .specify/"
```

### GSD
```
Next steps:
  In Claude Code, run:
    "Using GSD, execute PLAN.md in gsd/"
```

### GasTown
```
Next steps:
  Run: gastown run --formula gastown/formula.toml
```

### GasCity
```
Next steps:
  In Claude Code, run:
    "Using GasCity, coordinate the city in gascity/"
```
