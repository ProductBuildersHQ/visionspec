# synthesize

Generate specifications from source documents using LLM synthesis.

## Usage

```bash
visionspec synthesize <type> [flags]
```

## Description

The `synthesize` command generates a specification document from its source
specs using an LLM. Which sources a type needs, and what guidance drives its
prompt, come from the project's configured profile (its `synthesis:` rules)
— not a single fixed flow. Run `visionspec workflow` to see the exact
sequence and sources for the current project, or [`profiles show
<name>`](profiles.md) for another profile.

A type with no synthesis rule under the configured profile — or an
explicit empty source list, e.g. Press under `aws-one-way-door` — is
human-authored instead; use [`create`](create.md) for those.

## Example: aws-one-way-door

```
Press (human-authored — write this first)
    ↓
FAQ (synthesized from Press)
    ↓
MRD or OpportunitySpec (optional deepening, synthesized from Press + FAQ)
    ↓
PRD (synthesized from Press + FAQ + MRD/OpportunitySpec)
    ↓
Narrative 6P (synthesized from Press + FAQ + PRD)
    ↓
UXD (human-authored)
    ↓
TRD (synthesized from PRD + UXD + MRD, where present)
    ↓
TPD, IRD (synthesized from PRD/TRD/UXD)
```

This is one profile's sequence, not a universal one — `enterprise`, for
example, has no Press/FAQ at all and derives PRD from MRD directly. See
[Working Backwards](../concepts/working-backwards.md) for the methodology
behind the AWS door profiles specifically.

## Arguments

| Argument | Description |
|----------|-------------|
| `type` | Spec type to synthesize. Valid values depend on the project's configured profile — see `visionspec workflow`. |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--eval` | bool | `false` | Run evaluation after synthesis |
| `--no-context` | bool | `false` | Skip context gathering for technical synthesis |

## Without a Configured Profile

If the project has no `--profile` set (bare `visionspec init` with no
profile), synthesis falls back to the original Working Backwards flow for
these types only:

| Target | Required Sources |
|--------|------------------|
| `press` | mrd |
| `faq` | mrd, press |
| `prd` | mrd, press, faq |
| `trd` | mrd, prd |
| `tpd` | prd, trd, uxd |
| `ird` | trd |
| `narrative-1p` | mrd, prd |
| `narrative-6p` | mrd, prd |

## Examples

```bash
# See this project's actual sequence and sources first
visionspec workflow

# Then synthesize in that order, e.g. under aws-one-way-door:
visionspec create press            # Human-authored — write it first
visionspec synthesize faq          # From press
visionspec synthesize prd          # From press + faq (+ mrd/opportunity-spec if present)

# Technical synthesis
visionspec synthesize trd --eval
visionspec synthesize tpd
visionspec synthesize ird --no-context
```

## Context Grounding

For TRD, TPD, and IRD synthesis, the command automatically gathers codebase context if configured in `visionspec.yaml`:

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

Example under `aws-one-way-door` (press already created via `create`):

```
⋯ Synthesizing faq from [press]...
✓ Generated docs/specs/my-project/gtm/faq.md

⋯ Synthesizing prd from [press faq]...
✓ Generated docs/specs/my-project/source/prd.md

⋯ Gathering codebase context for grounding...
  Gathered context from 2 sources
⋯ Synthesizing trd from [prd uxd]...
✓ Generated docs/specs/my-project/technical/trd.md

⋯ Evaluating trd...
✓ trd: 8.2/10 PASS
```

The bracketed list after "from" is always the actual sources found on disk
for that type under the project's configured profile — not a fixed set.

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
