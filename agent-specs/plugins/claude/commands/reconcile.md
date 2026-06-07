---
name: reconcile
description: Reconcile all approved specifications into unified spec.md
arguments: []
dependencies: [spec-synthesis]
---

# Reconcile Specifications

Combine all approved specs into a unified implementation-ready spec.md.

## Usage

```
/reconcile
```

## Prerequisites

Required approved specs (profile-dependent):
- MRD (always required)
- Press Release (always required)
- FAQ (always required)
- PRD (always required)
- TRD (always required)
- TPD (always required)
- UXD (enterprise profile)
- IRD (enterprise profile)

## Process

1. **Verify Prerequisites**
   - Check all required specs exist
   - Check all required specs are approved

2. **Load Specs**
   - Read all approved specifications
   - Parse structure and content

3. **Detect Conflicts**
   - Cross-reference requirements
   - Identify contradictions
   - Flag unresolved issues

4. **Generate Decision Log**
   - Document tradeoffs
   - Record conflict resolutions

5. **Generate spec.md**
   - Unified document
   - Traceability matrix
   - Implementation guidance

## Output

```
⋯ Reconciling specifications...

  Loading approved specs:
    ✓ mrd (approved 2026-06-01)
    ✓ press (approved 2026-06-01)
    ✓ faq (approved 2026-06-01)
    ✓ prd (approved 2026-06-01)
    ✓ trd (approved 2026-06-01)
    ✓ tpd (approved 2026-06-01)

  Checking consistency...
    ✓ No conflicts detected

  Generating spec.md...

✓ Reconciliation complete
  Output: docs/specs/{project}/spec.md

  Summary:
    - 6 specs reconciled
    - 12 requirements traced
    - 8 user stories included
    - 4 implementation phases

Next steps:
  - Review spec.md
  - Run `/export <target>` to export for implementation
```

## Error Cases

### Missing Specs
```
✗ Cannot reconcile
  Missing required specs:
    - technical/trd.md (not found)
    - technical/tpd.md (not found)

  Run `/synthesize trd` and `/synthesize tpd` first.
```

### Unapproved Specs
```
✗ Cannot reconcile
  Unapproved specs:
    - source/mrd.md (pending)
    - technical/prd.md (pending)

  Run `/approve mrd` and `/approve prd` first.
```

### Conflicts Detected
```
⚠ Reconciliation completed with conflicts
  Output: docs/specs/{project}/spec.md

  Conflicts detected:
    - MR-3 requires feature X, but TRD excludes it
    - PRD targets 100ms latency, TRD targets 200ms

  See Decision Log in spec.md for resolutions.
  Review and update upstream specs if needed.
```
