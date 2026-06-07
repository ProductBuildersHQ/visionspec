---
name: synthesize
description: Synthesize a specification from upstream documents
arguments: [type]
dependencies: [working-backwards, spec-synthesis]
---

# Synthesize Specification

Generate a downstream specification from upstream documents.

## Usage

```
/synthesize <type>
```

Where `<type>` is one of:
- `press` - Press Release (requires: MRD)
- `faq` - FAQ (requires: MRD, Press)
- `prd` - Product Requirements (requires: MRD, Press, FAQ)
- `trd` - Technical Requirements (requires: PRD, UXD optional)
- `tpd` - Test Plan (requires: TRD)
- `ird` - Infrastructure Requirements (requires: TRD)

## Dependencies

```
press ← mrd
faq ← mrd, press
prd ← mrd, press, faq
trd ← prd, [uxd]
tpd ← trd
ird ← trd
```

## Process

1. **Verify Prerequisites**
   - Check required upstream specs exist
   - Verify upstream specs are approved (or warn if not)

2. **Load Upstream**
   - Read all required upstream documents
   - Parse and extract relevant sections

3. **Synthesize**
   - Apply transformation rules
   - Generate downstream content
   - Maintain traceability

4. **Write Output**
   - Write to appropriate location
   - Include synthesis metadata

5. **Display Summary**
   - Show what was generated
   - Note any warnings or issues

## Output Locations

| Type | Location |
|------|----------|
| press | `docs/specs/{project}/gtm/press.md` |
| faq | `docs/specs/{project}/gtm/faq.md` |
| prd | `docs/specs/{project}/technical/prd.md` |
| trd | `docs/specs/{project}/technical/trd.md` |
| tpd | `docs/specs/{project}/technical/tpd.md` |
| ird | `docs/specs/{project}/technical/ird.md` |

## Output

```
⋯ Synthesizing {type}...
  Loading: source/mrd.md
  Loading: gtm/press.md (if applicable)

✓ Synthesized {type}
  Output: docs/specs/{project}/{dir}/{type}.md

  Summary:
    - {N} requirements mapped
    - {M} user stories generated
    - {K} traceability links created

Next steps:
1. Review the generated spec
2. Run `/eval {type}` to check quality
3. Run `/approve {type}` when ready
```

## Error Cases

### Missing Prerequisites
```
✗ Cannot synthesize {type}
  Missing required specs:
    - source/mrd.md (not found)

  Run `/create mrd` first.
```

### Unapproved Prerequisites
```
⚠ Synthesizing {type} with unapproved prerequisites
  Warning: The following specs are not approved:
    - source/mrd.md (pending)

  Proceeding anyway. Re-synthesize after approval for consistency.
```
