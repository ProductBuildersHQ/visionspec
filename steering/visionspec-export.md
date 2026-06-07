# Steering: visionspec-export

## Command

```
visionspec export <target>
visionspec export all
```

## Purpose

Export reconciled specifications to execution targets.

## When to Invoke

- spec.md is generated
- Ready for implementation
- User wants to publish or deploy
- User asks to export to specific platform

## Triggers

- "export to ai-dlc"
- "publish specs"
- "deploy to speckit"
- "generate for coding assistant"
- "export all"

## Prerequisites

- `spec.md` must exist (reconciled)
- Target must be supported or configured

## Supported Targets

| Target | Output Directory | Purpose |
|--------|------------------|---------|
| ai-dlc | `.aidlc/` | AWS AI Developer Loop |
| speckit | `.specify/` | GitHub SpecKit |
| gsd | `gsd/` | GSD format |
| gastown | `gastown/` | GasTown deployment |
| gascity | `gascity/` | GasCity orchestration |
| markdown | `docs/` | Documentation |
| json | `export/` | Machine-readable |

## Expected Output

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

Export complete: .aidlc/
```

## Target Output Structures

### AI-DLC

```
.aidlc/
├── context.md
├── requirements.json
├── architecture.md
└── test-plan.md
```

### SpecKit

```
.specify/
├── spec.yaml
├── requirements/
├── decisions/
└── tests/
```

## Follow-up Actions

After export:

1. Review generated files
2. Commit to version control
3. Begin implementation with AI assistant

## Export All

```
visionspec export all
```

Exports to all configured targets.

## Error Handling

- **No spec.md**: Run `visionspec reconcile` first
- **Invalid target**: Check supported targets
- **Validation failure**: Fix source spec issues
