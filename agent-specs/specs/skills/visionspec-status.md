---
name: visionspec-status
description: Check project readiness and specification status
triggers: [status, project status, check status, what's done, progress]
---

# VisionSpec Status

Check the current status of specifications in a VisionSpec project.

## Purpose

Provides a quick overview of:

- Which specifications exist
- Approval and evaluation status
- What's blocking reconciliation
- Recommended next steps

## When to Use

- Starting a session to understand current state
- Before reconciliation to identify blockers
- After synthesis to see what's next
- When user asks "what's done?" or "what's next?"

## Invocation

```
visionspec status
```

Or via Claude Code:

```
/status
```

## Output Interpretation

| Symbol | Meaning |
|--------|---------|
| `✓` | Approved and ready |
| `⋯` | Exists, pending approval |
| `○` | Not created yet |
| `✗` | Blocked by dependencies |

## Integration Points

- **Profiles**: Status respects profile requirements (startup vs enterprise)
- **Evaluations**: Shows latest eval scores
- **Approvals**: Shows approval timestamps
- **Exports**: Shows which targets have exports

## Next Actions

Based on status output:

- Missing specs → `visionspec-synthesize`
- Pending evaluation → `visionspec-eval`
- All approved → `visionspec-reconcile`
- Ready to ship → `visionspec-export`
