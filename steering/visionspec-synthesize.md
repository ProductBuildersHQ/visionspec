# Steering: visionspec-synthesize

## Command

```
visionspec synthesize <type>
visionspec synthesize press --from mrd
```

## Purpose

Generate downstream specifications from upstream inputs using AI synthesis.

## When to Invoke

- After upstream specs are approved
- User wants to auto-generate specs
- Accelerating spec creation
- User asks to derive or create from existing

## Triggers

- "synthesize press"
- "generate prd from mrd"
- "create trd"
- "auto-generate"
- "derive from"
- "generate downstream"

## Synthesis Flow

```
MRD (authored)
  ↓
Press Release, FAQ
  ↓
PRD
  ↓
UXD, TRD, TPD, IRD
```

## Expected Output

```
Synthesizing trd from prd, uxd...

Sources loaded:
  ✓ source/mrd.md (approved)
  ✓ technical/prd.md (approved)
  ✓ technical/uxd.md (approved)

Generating draft...

✓ Draft created: drafts/trd.md

Traceability:
  - 12 requirements traced to PRD
  - 5 UX considerations incorporated
  - 3 non-functional requirements added

Next steps:
  1. Review drafts/trd.md
  2. Run visionspec eval trd
  3. Run visionspec approve trd
```

## Valid Synthesis Targets

| Target | Required Sources |
|--------|------------------|
| press | mrd (approved) |
| faq | mrd, press (approved) |
| prd | mrd, press, faq (approved) |
| uxd | prd (approved) |
| trd | prd (approved) |
| tpd | prd, trd (approved) |
| ird | trd (approved) |

## Follow-up Actions

After synthesis:

1. Review the draft
2. Run `visionspec eval <type>`
3. Address findings
4. Run `visionspec approve <type>`

## Prerequisites

- Upstream specs must be approved
- Profile must allow the spec type
- visionspec.yaml must exist

## Error Handling

- **Missing upstream**: Synthesize or author upstream first
- **Not approved**: Approve upstream specs first
- **Profile restriction**: Check profile requirements
