# Steering: visionspec-status

## Command

```
visionspec status
```

## Purpose

Display the current status of specifications in a VisionSpec project.

## When to Invoke

- User asks about project status or progress
- Beginning a session to understand state
- Before recommending next actions
- User asks "what's done?" or "what's next?"

## Triggers

- "what's the status"
- "show me progress"
- "what specs exist"
- "what's done"
- "what's next"
- "check project"

## Expected Output

```
Project: {project-name}
Profile: {profile}

Specifications:
  Source (Authored):
    ✓ MRD: approved
    ○ UXD: not created

  GTM (Synthesized):
    ✓ Press: approved
    ✓ FAQ: approved

  Technical (Synthesized):
    ✓ PRD: approved
    ⋯ TRD: pending evaluation

Reconciliation:
  ✗ spec.md: not generated

Next steps:
  1. Run visionspec eval trd
  2. Run visionspec approve trd
  3. Run visionspec reconcile
```

## Follow-up Actions

Based on status:

| Status | Recommended Command |
|--------|---------------------|
| Missing specs | `visionspec synthesize <type>` |
| Pending eval | `visionspec eval <type>` |
| Pending approval | `visionspec approve <type>` |
| All approved | `visionspec reconcile` |
| Reconciled | `visionspec export <target>` |

## Profile Awareness

Status respects active profile requirements:

- **startup**: MRD, PRD, TRD required
- **enterprise**: All specs required
- **custom**: Per profile definition

## Error Handling

If no project found:

```
Error: No visionspec.yaml found
Run: visionspec init
```
