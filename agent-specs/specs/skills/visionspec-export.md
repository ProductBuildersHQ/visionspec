---
name: visionspec-export
description: Export specifications to execution targets
triggers: [export, deploy, publish, ai-dlc, speckit, gsd]
---

# VisionSpec Export

Export reconciled specifications to execution targets.

## Purpose

Transforms specs into formats for:

- AI coding assistants (AI-DLC, SpecKit)
- Development workflows (GSD, GasTown)
- Documentation systems
- Custom targets

## When to Use

- spec.md is generated
- Ready for implementation handoff
- Targeting specific execution platforms
- Automating deployment

## Invocation

```
visionspec export <target>
visionspec export ai-dlc
visionspec export all
```

Or via Claude Code:

```
/export ai-dlc
/export speckit
```

## Supported Targets

| Target | Output | Purpose |
|--------|--------|---------|
| ai-dlc | `.aidlc/` | AWS AI Developer Loop Companion |
| speckit | `.specify/` | GitHub SpecKit format |
| gsd | `gsd/` | GSD specification format |
| gastown | `gastown/` | GasTown deployment |
| gascity | `gascity/` | GasCity orchestration |
| markdown | `docs/` | Documentation site |
| json | `export/` | Machine-readable JSON |

## Process

1. **Load spec.md** - Read reconciled specification
2. **Apply Transform** - Target-specific conversion
3. **Generate Files** - Create target structure
4. **Validate Output** - Check format requirements
5. **Report** - Show exported artifacts

## Output

```
Exporting to ai-dlc...

Source: spec.md (reconciled 2026-06-03)

Generating AI-DLC artifacts:
  ✓ .aidlc/context.md
  ✓ .aidlc/requirements.json
  ✓ .aidlc/architecture.md
  ✓ .aidlc/test-plan.md

Validation:
  ✓ All required fields present
  ✓ JSON schemas valid
  ✓ References resolved

Export complete: .aidlc/

Next steps:
  - Commit .aidlc/ to repository
  - AI coding assistant ready to use
```

## Target Details

### AI-DLC

AWS AI Developer Loop Companion format:

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
├── spec.yaml        # Specification manifest
├── requirements/    # Individual requirements
├── decisions/       # Architecture Decision Records
└── tests/           # Test specifications
```

### GSD

GSD specification format:

```
gsd/
├── spec.gsd.yaml    # GSD manifest
├── features/        # Feature specifications
└── tasks/           # Implementation tasks
```

## Custom Targets

Define custom export targets:

```yaml
# .visionspec/exports/custom.yaml
name: custom
output_dir: custom-output
transforms:
  - extract_requirements
  - generate_tasks
  - create_manifest
template: templates/custom/
```

## Export All

Export to all configured targets:

```
visionspec export all

Exporting to all targets...
  ✓ ai-dlc: .aidlc/
  ✓ speckit: .specify/
  ✓ gsd: gsd/

Summary: 3 targets exported
```

## CI Integration

```yaml
- name: Export specs
  run: |
    visionspec reconcile
    visionspec export ai-dlc
    git add .aidlc/
```
