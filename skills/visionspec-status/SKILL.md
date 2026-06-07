# VisionSpec Status

Check the current status of a VisionSpec project.

## Overview

The status command provides a quick overview of specification progress, showing which specs exist, their approval status, and recommended next steps.

## Workflow

### 1. Check Project Status

Run status to understand current state:

```
visionspec status
```

### 2. Interpret Output

The output shows:

- **Source Specs**: Authored specifications (MRD, UXD)
- **GTM Specs**: Go-to-market documents (Press, FAQ)
- **Technical Specs**: Implementation specs (PRD, TRD, TPD, IRD)
- **Reconciliation**: Whether spec.md is generated
- **Exports**: Which targets have been exported to

### 3. Status Symbols

| Symbol | Meaning | Action Needed |
|--------|---------|---------------|
| `✓` | Approved | None |
| `⋯` | Exists, pending | Run eval, then approve |
| `○` | Not created | Synthesize or author |
| `✗` | Blocked | Resolve dependencies |

### 4. Follow Recommendations

The status output includes "Next steps" with prioritized actions.

## Integration

Status respects the active profile's requirements:

- **startup**: MRD, PRD, TRD required
- **enterprise**: All specs required
- **custom**: Profile-specific requirements

## When to Use

- Starting a session to understand current state
- Before reconciliation to identify blockers
- After synthesis to see what's next
- When user asks "what's done?" or "what's next?"

## Related Skills

After checking status:

- Missing specs → `author-*` or `visionspec-synthesize`
- Pending evaluation → `visionspec-eval`
- All approved → `visionspec-reconcile`
- Ready to ship → `visionspec-export`
