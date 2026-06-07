---
name: approve
description: Approve a specification for use in reconciliation
arguments: [type, approver]
dependencies: []
---

# Approve Specification

Mark a specification as approved for reconciliation.

## Usage

```
/approve <type>
/approve <type> --approver "name@example.com"
```

## Prerequisites

- Spec must exist
- Spec should have passing evaluation (warning if not)

## Process

1. **Verify Spec Exists** - Check spec file is present
2. **Check Evaluation** - Warn if no evaluation or failing
3. **Record Approval** - Create approval record
4. **Display Confirmation** - Show approval status

## Output

```
⋯ Approving {type}...

✓ Approved {type}
  Spec: docs/specs/{project}/source/{type}.md
  Approver: {approver}
  Timestamp: 2026-06-01T12:00:00Z

  Approval record: docs/specs/{project}/approval/{type}.approval.json
```

## Approval Record Format

```json
{
  "spec_type": "mrd",
  "spec_path": "docs/specs/project/source/mrd.md",
  "spec_hash": "sha256:abc123...",
  "approved_at": "2026-06-01T12:00:00Z",
  "approved_by": "name@example.com",
  "evaluation": {
    "score": 85,
    "status": "pass",
    "report_path": "eval/mrd.eval.json"
  }
}
```

## Warnings

### No Evaluation
```
⚠ Approving {type} without evaluation
  Run `/eval {type}` first for quality assurance.
  Proceeding with approval anyway.
```

### Failing Evaluation
```
⚠ Approving {type} with failing evaluation
  Evaluation score: 55/100 (FAIL)
  Consider addressing issues before approval.
  Proceeding with approval anyway.
```

## Revoke Approval

To revoke an approval, delete the approval file:
```bash
rm docs/specs/{project}/approval/{type}.approval.json
```

Then re-approve after making changes.
