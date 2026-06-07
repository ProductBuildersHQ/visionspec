# VisionSpec Export

Export specifications to execution targets.

## Overview

The export command transforms reconciled specifications into formats suitable for AI coding assistants, development workflows, and documentation systems.

## Workflow

### 1. Verify Prerequisites

Export requires a reconciled `spec.md`. Check with:

```
visionspec status
```

### 2. Choose Target

Select an export target:

| Target | Output | Purpose |
|--------|--------|---------|
| ai-dlc | `.aidlc/` | AWS AI Developer Loop Companion |
| speckit | `.specify/` | GitHub SpecKit format |
| gsd | `gsd/` | GSD specification format |
| gastown | `gastown/` | GasTown deployment |
| gascity | `gascity/` | GasCity orchestration |
| markdown | `docs/` | Documentation site |
| json | `export/` | Machine-readable JSON |

### 3. Run Export

Single target:

```
visionspec export ai-dlc
```

All targets:

```
visionspec export all
```

### 4. Verify Output

Check generated files:

```
ls -la .aidlc/
```

### 5. Commit

Add exported files to version control:

```bash
git add .aidlc/
git commit -m "chore: export specs to AI-DLC"
```

## Target Details

### AI-DLC

AWS AI Developer Loop Companion:

```
.aidlc/
├── context.md       # Project context
├── requirements.json # Structured requirements
├── architecture.md  # Technical design
└── test-plan.md     # Test strategy
```

### SpecKit

GitHub SpecKit format:

```
.specify/
├── spec.yaml        # Manifest
├── requirements/    # Individual requirements
├── decisions/       # ADRs
└── tests/           # Test specifications
```

### GSD

GSD specification format:

```
gsd/
├── spec.gsd.yaml    # Manifest
├── features/        # Feature specs
└── tasks/           # Implementation tasks
```

## CI Integration

Add to your pipeline:

```yaml
- name: Export specs
  run: |
    visionspec reconcile
    visionspec export ai-dlc
    git add .aidlc/
```

## Custom Targets

Define custom export in `.visionspec/exports/`:

```yaml
# custom.yaml
name: custom
output_dir: custom-output
transforms:
  - extract_requirements
  - generate_tasks
template: templates/custom/
```

## Tips

- Export after each reconciliation
- Commit exports to version control
- Use AI-DLC for AI coding assistants
- Use SpecKit for GitHub integration

## Related Skills

- `visionspec-reconcile`: Generate spec.md first
- `visionspec-status`: Check export status
