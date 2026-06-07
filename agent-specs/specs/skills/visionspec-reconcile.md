---
name: visionspec-reconcile
description: Generate unified specification from all approved specs
triggers: [reconcile, unify, consolidate, final spec, spec.md]
---

# VisionSpec Reconcile

Generate a unified specification document from all approved specs.

## Purpose

Creates a single `spec.md` that:

- Integrates all approved specifications
- Resolves cross-document references
- Provides complete implementation guidance
- Serves as the source of truth

## When to Use

- All required specs are approved
- Ready to hand off to implementation
- Need a consolidated view
- Preparing for export

## Invocation

```
visionspec reconcile
visionspec reconcile --dry-run
```

Or via Claude Code:

```
/reconcile
```

## Prerequisites

Reconciliation requires:

1. **Profile Requirements Met** - All required specs exist
2. **All Specs Approved** - No pending evaluations
3. **Consistency** - No unresolved conflicts

## Process

1. **Check Readiness** - Verify all requirements met
2. **Load Specs** - Read all approved specifications
3. **Detect Conflicts** - Identify inconsistencies
4. **Resolve** - Apply resolution rules or prompt
5. **Generate** - Create unified spec.md
6. **Validate** - Final consistency check

## Output

```
Reconciling specifications...

Loading approved specs:
  ✓ MRD (approved 2026-06-01)
  ✓ Press (approved 2026-06-01)
  ✓ FAQ (approved 2026-06-01)
  ✓ PRD (approved 2026-06-02)
  ✓ TRD (approved 2026-06-03)
  ✓ TPD (approved 2026-06-03)

Conflict detection:
  ⚠ Timeline mismatch: MRD says Q3, TRD says Q4
  → Resolved: Using TRD timeline (more specific)

Generating spec.md...

✓ Reconciliation complete: spec.md

Contents:
  - 15 features from PRD
  - 8 architecture decisions from TRD
  - 42 test cases from TPD
  - Complete traceability matrix

Next steps:
  - Review spec.md
  - Run /export <target> to deploy
```

## Conflict Resolution

When conflicts are detected:

| Conflict Type | Resolution |
|---------------|------------|
| Timeline | Use most specific |
| Scope | Use PRD as authoritative |
| Technical | Use TRD as authoritative |
| Requirements | Prompt user |

## Dry Run

Preview reconciliation without generating:

```
visionspec reconcile --dry-run

Would reconcile:
  - 6 specifications
  - 2 potential conflicts
  - Estimated output: ~8,500 words
```

## spec.md Structure

The generated spec.md includes:

1. **Executive Summary** - From Press Release
2. **Background** - From MRD problem statement
3. **Requirements** - From PRD
4. **Technical Design** - From TRD
5. **Test Strategy** - From TPD
6. **Appendices** - FAQ, traceability matrix

## Post-Reconciliation

After spec.md is generated:

- Review for accuracy
- Export to targets as needed
- Archive approved specs
- Begin implementation
