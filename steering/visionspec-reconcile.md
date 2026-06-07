# Steering: visionspec-reconcile

## Command

```
visionspec reconcile
visionspec reconcile --dry-run
```

## Purpose

Generate a unified `spec.md` from all approved specifications.

## When to Invoke

- All required specs are approved
- Ready for implementation handoff
- User asks to consolidate or unify
- User wants final spec

## Triggers

- "reconcile"
- "generate spec.md"
- "unify specs"
- "consolidate"
- "final spec"
- "ready for implementation"

## Prerequisites

| Requirement | Check |
|-------------|-------|
| Profile requirements met | All required specs exist |
| All specs approved | No pending evaluations |
| Consistency | No unresolved conflicts |

## Expected Output

```
Reconciling specifications...

Loading approved specs:
  ✓ MRD (approved 2026-06-01)
  ✓ Press (approved 2026-06-01)
  ✓ FAQ (approved 2026-06-01)
  ✓ PRD (approved 2026-06-02)
  ✓ TRD (approved 2026-06-03)

Conflict detection:
  ⚠ Timeline mismatch: MRD says Q3, TRD says Q4
  → Resolved: Using TRD timeline (more specific)

Generating spec.md...

✓ Reconciliation complete: spec.md
```

## Conflict Resolution

| Conflict Type | Resolution |
|---------------|------------|
| Timeline | Use most specific |
| Scope | PRD is authoritative |
| Technical | TRD is authoritative |
| Other | Prompt for decision |

## spec.md Structure

Generated file includes:

1. Executive Summary (from Press)
2. Background (from MRD)
3. Requirements (from PRD)
4. Technical Design (from TRD)
5. Test Strategy (from TPD)
6. Appendices (FAQ, traceability)

## Follow-up Actions

After reconciliation:

1. Review spec.md
2. Run `visionspec export <target>`
3. Begin implementation

## Dry Run

Preview without generating:

```
visionspec reconcile --dry-run
```

## Error Handling

- **Missing specs**: Run `visionspec status` to identify
- **Unapproved specs**: Approve pending specs first
- **Conflicts**: Manual resolution may be needed
