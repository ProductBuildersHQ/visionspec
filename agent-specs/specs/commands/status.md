---
name: status
description: Show the current status of the VisionSpec project
arguments: []
dependencies: []
---

# Project Status

Display the current status of specifications in the project.

## Usage

```
/status
```

## Output

```
Project: {project-name}
Profile: {profile}
Path: docs/specs/{project}/

Specifications:
  Source (Authored):
    ✓ MRD: approved (2026-06-01)
    ○ UXD: not created

  GTM (Synthesized):
    ✓ Press: approved (2026-06-01)
    ✓ FAQ: approved (2026-06-01)

  Technical (Synthesized):
    ✓ PRD: approved (2026-06-01)
    ⋯ TRD: pending evaluation
    ○ TPD: not created
    ○ IRD: not created

Reconciliation:
  ✗ spec.md: not generated
  Missing: TPD, IRD

Export:
  ○ No exports yet

Legend:
  ✓ = approved
  ⋯ = exists, pending approval
  ○ = not created
  ✗ = blocked

Next steps:
  1. Run `/synthesize tpd` to create Test Plan
  2. Run `/synthesize ird` to create Infrastructure Requirements
  3. Run `/eval trd` and `/approve trd`
  4. Run `/reconcile` when all specs approved
```

## Status Checks

### Spec Existence
Check for files at expected locations:
- `source/mrd.md`
- `source/uxd.md`
- `gtm/press.md`
- `gtm/faq.md`
- `technical/prd.md`
- `technical/trd.md`
- `technical/tpd.md`
- `technical/ird.md`
- `spec.md`

### Approval Status
Check for approval records:
- `approval/*.approval.json`

### Evaluation Status
Check for evaluation reports:
- `eval/*.eval.json`

### Export Status
Check for export directories:
- `.aidlc/`
- `.specify/`
- `gsd/`
- `gastown/`
- `gascity/`

## Profile Requirements

### Startup Profile
Required: MRD, Press, FAQ, PRD, TRD, TPD
Optional: UXD, IRD

### Enterprise Profile
Required: MRD, Press, FAQ, PRD, UXD, TRD, TPD, IRD

### Minimal Profile
Required: MRD, PRD, TRD
Optional: All others
