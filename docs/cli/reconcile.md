# reconcile

Generate a unified execution spec from approved specifications.

## Usage

```bash
multispec reconcile [flags]
```

## Description

The `reconcile` command combines all approved specifications into a unified execution spec (`spec.md`). This is the final step before exporting to AI coding agent execution systems.

The reconciliation process:

1. Verifies all required specs are approved
2. Detects conflicts between specs
3. Resolves conflicts using LLM-assisted analysis
4. Generates unified `spec.md`
5. Creates `spec.eval.json` with reconciliation metadata

## Prerequisites

All required specs must be approved before reconciliation:

```bash
multispec approve mrd
multispec approve prd
multispec approve uxd
multispec approve trd
```

## Conflict Detection

The reconciler detects conflicts across specs:

| Conflict Type | Specs Involved | Example |
|---------------|----------------|---------|
| Performance vs Feature | PRD, TRD | Feature requires real-time sync but TRD specifies batch processing |
| Security vs Usability | UXD, TRD | UX flow bypasses security controls |
| Scope vs Timeline | MRD, PRD | PRD includes features not in MRD scope |
| Missing Traceability | PRD, TRD | Requirements not covered in technical design |

## Output Files

**spec.md** - Unified execution specification containing:

- Executive summary
- Consolidated requirements
- Technical architecture
- Implementation phases
- Acceptance criteria

**eval/spec.eval.json** - Reconciliation metadata:

```json
{
  "spec_type": "spec",
  "generated_at": "2024-01-15T10:30:00Z",
  "sources": ["mrd", "prd", "uxd", "trd", "ird"],
  "conflicts": [
    {
      "id": "C001",
      "type": "performance_vs_feature",
      "description": "Real-time requirement conflicts with batch architecture",
      "severity": "high",
      "resolution": "Implement hybrid approach with real-time for critical paths"
    }
  ],
  "decision_log": [...],
  "status": "reconciled_with_tradeoffs"
}
```

## Reconciliation Status

| Status | Description |
|--------|-------------|
| `reconciled` | No conflicts detected |
| `reconciled_with_tradeoffs` | Conflicts resolved with documented tradeoffs |
| `needs_review` | Unresolved high-severity conflicts |

## Examples

```bash
# Run reconciliation
multispec reconcile

# Output:
# ✓ All required specs approved: [mrd prd uxd trd ird]
# ⋯ Reconciling specs...
# ✓ Generated docs/specs/my-project/spec.md
#   Sources: [mrd prd uxd trd ird]
#   Conflicts detected: 2
#     ✓ C001: Performance vs feature trade-off
#     ✓ C002: Security model alignment
# ✓ Generated docs/specs/my-project/eval/spec.eval.json
```

## Handling Failures

If reconciliation fails due to missing approvals:

```
Missing approvals:
  ✗ trd
  ✗ ird

Approve specs with: multispec approve <spec-type>
```

If high-severity conflicts cannot be resolved, the status will be `needs_review` and manual intervention is required.

## See Also

- [approve](approve.md) - Approve specs for reconciliation
- [export](export.md) - Export reconciled spec to target systems
- [status](status.md) - Check project readiness
