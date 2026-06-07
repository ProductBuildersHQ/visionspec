---
name: visionspec-synthesize
description: Generate downstream specifications from upstream inputs
triggers: [synthesize, generate, create from, derive, auto-generate]
---

# VisionSpec Synthesize

Generate downstream specifications from upstream inputs using AI synthesis.

## Purpose

Automatically creates specifications by:

- Analyzing upstream documents
- Applying Working Backwards methodology
- Following spec-type templates
- Ensuring consistency and traceability

## When to Use

- After MRD is approved → synthesize Press, FAQ
- After Press/FAQ approved → synthesize PRD
- After PRD approved → synthesize TRD, TPD, IRD
- To accelerate spec creation

## Invocation

```
visionspec synthesize <spec-type>
visionspec synthesize press --from mrd
```

Or via Claude Code:

```
/synthesize press
/synthesize prd
```

## Synthesis Flow

```
MRD (authored)
  ↓
Press Release (synthesized)
FAQ (synthesized)
  ↓
PRD (synthesized or authored)
  ↓
UXD (synthesized or authored)
TRD (synthesized)
TPD (synthesized)
IRD (synthesized)
```

## Supported Synthesis

| Target | Sources | Description |
|--------|---------|-------------|
| press | mrd | Customer-facing announcement |
| faq | mrd, press | Customer and internal questions |
| prd | mrd, press, faq | Product requirements |
| uxd | prd | User experience design |
| trd | prd, uxd | Technical requirements |
| tpd | prd, trd | Test plan |
| ird | trd | Infrastructure requirements |

## Process

1. **Load Sources** - Read upstream specs
2. **Apply Template** - Use profile template
3. **Generate Draft** - AI synthesis with traceability
4. **Save Draft** - Write to drafts directory
5. **Suggest Eval** - Prompt for evaluation

## Output

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
  2. Run /eval trd
  3. Run /approve trd
```

## Customization

### Profile Templates

Synthesis uses profile-specific templates:

```
profiles/enterprise/templates/trd.md
profiles/startup/templates/trd.md
```

### Constitution

Organization constitutions influence tone:

```yaml
# .visionspec/constitution.yaml
voice:
  technical: formal
  documentation: comprehensive
```

## Quality Assurance

Synthesized specs include:

- Source document references
- Requirement traceability
- Consistency markers
- AI synthesis disclosure

All synthesized specs should be reviewed and evaluated before approval.
