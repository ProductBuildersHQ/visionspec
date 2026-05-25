# synthesize

Generate specifications from source documents using LLM synthesis.

## Usage

```bash
visionspec synthesize <type> [flags]
```

## Description

The `synthesize` command generates specification documents from existing source specs using an LLM. It implements Amazon's Working Backwards methodology where the Press Release defines the vision, and requirements flow from that vision.

## Working Backwards Flow

VisionSpec implements Amazon's Working Backwards methodology:

1. **MRD** - Define the market problem (human-authored)
2. **Press** - Write the press release announcing the solution
3. **FAQ** - Anticipate customer and stakeholder questions
4. **PRD** - Derive detailed requirements from the vision

```
MRD (human-authored)
    ↓
Press (synthesized from MRD)
    ↓
FAQ (synthesized from MRD + Press)
    ↓
PRD (synthesized from MRD + Press + FAQ)
    ↓
UXD (human-authored)
    ↓
TRD (synthesized from MRD + PRD + UXD + context)
    ↓
IRD (synthesized from TRD + context)
```

## Synthesis Types

**Working Backwards Flow**

- `press` - Press Release from MRD (vision document)
- `faq` - FAQ from MRD + Press (scope clarification)
- `prd` - PRD from MRD + Press + FAQ (detailed requirements)

**Technical Synthesis**

- `trd` - Technical Requirements from MRD + PRD + UXD + CONSTITUTION + CONTEXT
- `ird` - Infrastructure Requirements from TRD + CONSTITUTION + CONTEXT

**Narrative Documents**

- `narrative-1p` - 1-Page Narrative from MRD + PRD
- `narrative-6p` - 6-Page Narrative from MRD + PRD + UXD

## Arguments

| Argument | Description |
|----------|-------------|
| `type` | Spec type to synthesize: `press`, `faq`, `prd`, `trd`, `ird`, `narrative-1p`, `narrative-6p` |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--eval` | bool | `false` | Run evaluation after synthesis |
| `--no-context` | bool | `false` | Skip context gathering for technical synthesis |

## Synthesis Dependencies

| Target | Required Sources | Description |
|--------|------------------|-------------|
| `press` | MRD | Vision document from market requirements |
| `faq` | MRD, Press | Scope clarification from vision |
| `prd` | MRD, Press, FAQ | Detailed requirements from Working Backwards artifacts |
| `trd` | MRD, PRD | Technical requirements (UXD optional) |
| `ird` | TRD | Infrastructure requirements |
| `narrative-1p` | MRD, PRD | 1-page executive narrative |
| `narrative-6p` | MRD, PRD | 6-page detailed narrative (UXD optional) |

## Examples

```bash
# Working Backwards flow
visionspec synthesize press        # Generate Press Release from MRD
visionspec synthesize faq          # Generate FAQ from MRD + Press
visionspec synthesize prd          # Generate PRD from MRD + Press + FAQ

# Technical synthesis
visionspec synthesize trd --eval   # Generate TRD with evaluation
visionspec synthesize ird --no-context  # Generate IRD without context gathering

# Narrative documents
visionspec synthesize narrative-1p
visionspec synthesize narrative-6p
```

## Context Grounding

For TRD and IRD synthesis, the command automatically gathers codebase context if configured in `visionspec.yaml`:

```yaml
context:
  repositories:
    - path: "."
      include_structure: true
      include_deps: true
      include_apis: true
```

This grounds technical decisions in the reality of existing code. Use `--no-context` to skip this step.

## Output

```
⋯ Synthesizing press from [mrd]...
✓ Generated docs/specs/my-project/gtm/press.md

⋯ Synthesizing faq from [mrd press]...
✓ Generated docs/specs/my-project/gtm/faq.md

⋯ Synthesizing prd from [mrd press faq]...
✓ Generated docs/specs/my-project/source/prd.md

⋯ Gathering codebase context for grounding...
  Gathered context from 2 sources
⋯ Synthesizing trd from [mrd prd]...
✓ Generated docs/specs/my-project/technical/trd.md

⋯ Evaluating trd...
✓ trd: 8.2/10 PASS
```

## LLM Configuration

Configure the LLM in `visionspec.yaml`:

```yaml
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  temperature: 0.7
  max_tokens: 8192
```

## See Also

- [eval](eval.md) - Evaluate synthesized specs
- [reconcile](reconcile.md) - Combine specs into execution spec
- [context](context.md) - Manage context sources
