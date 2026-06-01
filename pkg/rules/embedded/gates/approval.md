# Approval Gate

Approval marks a spec as ready for use in downstream phases. Only approved specs are used in synthesis and reconciliation.

## Approval Command

```bash
visionspec approve <type> -p <project>
```

Or via MCP:

```
approve(project, specType)
```

## Approval Prerequisites

Before approving, verify:

| Prerequisite | Check |
|--------------|-------|
| Spec exists | File is at expected path |
| Evaluation passed | Score >= 7.0, no critical/high |
| User reviewed | Human has read and accepted |
| Dependencies met | Prerequisite specs approved |

## Dependency Chain

Specs must be approved in order:

```
MRD (no dependencies)
    ↓
Press (requires: MRD)
    ↓
FAQ (requires: MRD, Press)
    ↓
PRD (requires: MRD, Press, FAQ)
    ↓
UXD (no synthesis dependencies, but should follow PRD)
    ↓
TRD (requires: MRD, PRD)
    ↓
TPD (requires: PRD, TRD, UXD)
    ↓
IRD (requires: TRD)
```

## Approval Workflow

### Step 1: Verify Evaluation

```bash
visionspec eval prd -p myproject

# Must see:
# Score: >= 7.0
# No critical findings
# No high findings
```

### Step 2: User Confirmation

Before approving, confirm with user:

```
The PRD has passed evaluation with score 8.5.

Findings:
  [MEDIUM] User Stories: Story #3 could be more specific
  [LOW] Scope: Consider adding explicit exclusions

Do you want to:
1. Approve as-is (acknowledge findings)
2. Address findings first
3. Review the spec content
```

### Step 3: Record Approval

```bash
visionspec approve prd -p myproject

# Records:
# - Approver (from git config or env)
# - Timestamp
# - Any comments
```

### Step 4: Verify Status

```bash
visionspec status -p myproject

# Should show:
# [+] PRD: Approved (2024-01-15 by jsmith)
```

## Approval Records

Approvals are stored in `visionspec.yaml`:

```yaml
approvals:
  mrd:
    approver: jsmith
    approved_at: 2024-01-15T10:30:00Z
    comment: "Ready for Working Backwards phase"
  prd:
    approver: jsmith
    approved_at: 2024-01-15T14:45:00Z
    comment: "Approved with known scope question in story #3"
```

## Revoking Approval

If a spec needs changes after approval:

```bash
# Edit the spec
# This automatically revokes approval

# Re-evaluate
visionspec eval prd -p myproject

# Re-approve
visionspec approve prd -p myproject
```

## Approval Gates in Status

```bash
visionspec status -p myproject

Project Status: myproject
Path: docs/specs/myproject

Readiness: NOT READY

Gates:
  [+] Required specs present
  [+] Evaluations passing
  [X] Approvals obtained ← 2 specs need approval
  [ ] Execution spec generated

Specs:
  MRD: Present, Evaluated (8.5), Approved
  Press: Present, Evaluated (7.8), Approved
  FAQ: Present, Evaluated (8.0), Approved
  PRD: Present, Evaluated (8.5), NEEDS APPROVAL ←
  UXD: Present, Evaluated (7.5), NEEDS APPROVAL ←
  TRD: Not present
```

## Conditional Approval

For specs with medium/low findings, approval can include acknowledgment:

```bash
visionspec approve prd -p myproject --comment "Approved with known limitation in story #3"
```

## Approval Policies

Organizations can configure approval requirements:

```yaml
# visionspec.yaml
approval_policy:
  require_evaluation: true
  min_score: 7.0
  allow_medium_findings: true
  allow_low_findings: true
  require_comment: false
```

## See Also

- [evaluation.md](evaluation.md) - Evaluation before approval
- [../core-workflow.md](../core-workflow.md) - Overall workflow
