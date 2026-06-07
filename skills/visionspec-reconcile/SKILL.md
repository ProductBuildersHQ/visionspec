# VisionSpec Reconcile

Generate a unified specification from all approved specs.

## Overview

The reconcile command creates a single `spec.md` that integrates all approved specifications, resolves cross-document references, and provides complete implementation guidance.

## Workflow

### 1. Verify Prerequisites

Reconciliation requires all required specs to be approved. Check with:

```
visionspec status
```

Look for:

- All required specs exist (per profile)
- All specs show `✓` (approved)
- No blocking issues

### 2. Preview (Optional)

Dry-run to see what will be reconciled:

```
visionspec reconcile --dry-run
```

This shows:

- Specs to be included
- Potential conflicts
- Estimated output size

### 3. Run Reconcile

```
visionspec reconcile
```

### 4. Review Conflicts

If conflicts are detected, reconcile handles them:

| Conflict Type | Resolution |
|---------------|------------|
| Timeline | Use most specific source |
| Scope | PRD is authoritative |
| Technical | TRD is authoritative |
| Other | Prompt for decision |

### 5. Review spec.md

The generated `spec.md` includes:

1. **Executive Summary** - From Press Release
2. **Background** - From MRD problem statement
3. **Requirements** - From PRD
4. **Technical Design** - From TRD
5. **Test Strategy** - From TPD
6. **Appendices** - FAQ, traceability matrix

### 6. Export

Once reconciled, export to targets:

```
visionspec export ai-dlc
visionspec export speckit
```

## Prerequisites

| Requirement | Description |
|-------------|-------------|
| Profile Met | All required specs exist |
| All Approved | No pending evaluations |
| Consistency | No unresolved conflicts |

## Conflict Examples

**Timeline Mismatch**

```
MRD says: "Launch Q3 2026"
TRD says: "Launch Q4 2026"
→ Resolved: Using TRD (more specific technical estimate)
```

**Scope Conflict**

```
Press says: "Supports mobile"
PRD says: "Web-only MVP"
→ Prompt: Which scope is correct?
```

## Tips

- Run eval on all specs before reconcile
- Use dry-run to preview conflicts
- Keep PRD as the authoritative scope source
- Archive approved specs after reconcile

## Related Skills

- `visionspec-status`: Check reconcile readiness
- `visionspec-eval`: Ensure specs are ready
- `visionspec-export`: Export reconciled spec
