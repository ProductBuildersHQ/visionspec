# approve

Mark a specification as approved for reconciliation.

## Usage

```bash
visionspec approve <spec-type> [flags]
```

## Description

The `approve` command marks a specification as approved and ready for reconciliation. Approvals are required before running `visionspec reconcile` to generate the unified execution spec.

Approvals are recorded in `visionspec.yaml` with:

- Approver identifier
- Timestamp
- Optional comment

## Arguments

| Argument | Description |
|----------|-------------|
| `spec-type` | Type of spec to approve (mrd, prd, uxd, trd, ird, etc.) |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--approver` | string | `$USER` | Approver email or identifier |
| `--comment` | string | `""` | Approval comment |

## Examples

```bash
# Approve PRD with default approver
visionspec approve prd

# Approve TRD with specific approver
visionspec approve trd --approver=eng@company.com

# Approve with comment
visionspec approve mrd --approver=pm@company.com --comment="Reviewed in spec review meeting"
```

## Output

```
✓ Approved prd by john
  Comment: Reviewed in spec review meeting
```

## Approval Workflow

1. **Write spec** - Create or synthesize the specification
2. **Evaluate** - Run `visionspec eval <type>` to check quality
3. **Review** - Human review of spec content
4. **Approve** - Mark as approved with this command
5. **Reconcile** - Run `visionspec reconcile` after all approvals

## Required Approvals

The specs required for reconciliation depend on your profile:

| Profile | Required Approvals |
|---------|-------------------|
| `startup` | prd |
| `growth` | prd, uxd |
| `enterprise` | mrd, prd, uxd, trd |

## Checking Approval Status

Use `visionspec status` to see which specs are approved:

```bash
visionspec status

# Approvals:
#   ✓ mrd - approved by pm@co on 2024-01-10
#   ✓ prd - approved by pm@co on 2024-01-11
#   ✗ uxd - not approved
#   ✗ trd - not approved
```

## Storage

Approvals are stored in `visionspec.yaml`:

```yaml
approvals:
  mrd:
    approver: pm@company.com
    approved_at: 2024-01-10T14:30:00Z
    comment: "Reviewed in planning meeting"
  prd:
    approver: pm@company.com
    approved_at: 2024-01-11T09:00:00Z
```

## See Also

- [eval](eval.md) - Evaluate specs before approval
- [reconcile](reconcile.md) - Reconcile after all approvals
- [status](status.md) - Check approval status
