# VisionSpec Synthesize

Generate downstream specifications from upstream inputs.

## Overview

The synthesize command uses AI to create specifications by analyzing upstream documents, applying Working Backwards methodology, and following templates.

## Workflow

### 1. Verify Prerequisites

Synthesis requires approved upstream specs. Check with:

```
visionspec status
```

### 2. Run Synthesis

Generate a specific spec:

```
visionspec synthesize press
visionspec synthesize trd
```

### 3. Review Draft

Synthesized specs are saved as drafts:

```
drafts/press.md
drafts/trd.md
```

Review for:

- Accuracy of derived content
- Appropriate level of detail
- Correct traceability

### 4. Evaluate

Run evaluation on the draft:

```
visionspec eval press
```

### 5. Iterate or Approve

If evaluation passes:

```
visionspec approve press
```

If not, revise and re-evaluate.

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
| faq | mrd, press | Questions and answers |
| prd | mrd, press, faq | Product requirements |
| uxd | prd | User experience design |
| trd | prd, uxd | Technical requirements |
| tpd | prd, trd | Test plan |
| ird | trd | Infrastructure requirements |

## Customization

### Templates

Synthesis uses profile templates:

```
profiles/enterprise/templates/trd.md
profiles/startup/templates/trd.md
```

### Constitution

Organization constitutions influence tone and style:

```yaml
# .visionspec/constitution.yaml
voice:
  technical: formal
  documentation: comprehensive
```

## Quality Notes

Synthesized specs include:

- Source document references
- Requirement traceability
- AI synthesis disclosure

Always review synthesized content before approval.

## Tips

- Start with MRD to establish strong foundation
- Synthesize in dependency order
- Use eval feedback to improve synthesis quality
- Consider authoring critical specs (MRD, PRD) manually

## Related Skills

- `author-*`: Manually author specs
- `visionspec-eval`: Evaluate synthesized specs
- `visionspec-status`: Check synthesis prerequisites
