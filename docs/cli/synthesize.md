# synthesize

Generate specifications from source documents using LLM synthesis.

## Usage

```bash
multispec synthesize <type> [flags]
```

## Description

The `synthesize` command generates specification documents from existing source specs using an LLM. It supports two synthesis paths:

**GTM Synthesis (Working Backwards)**

- `press` - Press Release from MRD + PRD
- `faq` - FAQ from Press Release
- `narrative-1p` - 1-Page Narrative from MRD + PRD
- `narrative-6p` - 6-Page Narrative from MRD + PRD + UXD

**Technical Synthesis**

- `trd` - Technical Requirements from MRD + PRD + UXD + CONSTITUTION + CONTEXT
- `ird` - Infrastructure Requirements from TRD + CONSTITUTION + CONTEXT

## Arguments

| Argument | Description |
|----------|-------------|
| `type` | Spec type to synthesize: `trd`, `ird`, `press`, `faq`, `narrative-1p`, `narrative-6p` |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--eval` | bool | `false` | Run evaluation after synthesis |
| `--no-context` | bool | `false` | Skip context gathering for technical synthesis |

## Required Sources

| Target | Required Sources |
|--------|------------------|
| `press` | MRD, PRD |
| `faq` | Press Release |
| `narrative-1p` | MRD, PRD |
| `narrative-6p` | MRD, PRD, UXD |
| `trd` | MRD, PRD, UXD |
| `ird` | TRD |

## Examples

```bash
# Generate Press Release
multispec synthesize press

# Generate TRD with evaluation
multispec synthesize trd --eval

# Generate IRD without context gathering
multispec synthesize ird --no-context

# Generate FAQ document
multispec synthesize faq
```

## Context Grounding

For TRD and IRD synthesis, the command automatically gathers codebase context if configured in `multispec.yaml`:

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
⋯ Gathering codebase context for grounding...
  Gathered context from 2 sources
⋯ Synthesizing trd from [mrd prd uxd]...
✓ Generated docs/specs/my-project/technical/trd.md

⋯ Evaluating trd...
✓ trd: 8.2/10 PASS
```

## LLM Configuration

Configure the LLM in `multispec.yaml`:

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
