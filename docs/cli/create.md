# create

Create a new spec file from a template.

## Usage

```bash
multispec create <spec-type> [flags]
```

## Description

The `create` command scaffolds a new specification file from an embedded template. This provides a starting point for writing specs with the correct structure and sections.

The command must be run from within a multispec project directory (where `multispec.yaml` exists).

## Arguments

| Argument | Description |
|----------|-------------|
| `spec-type` | Type of spec to create |

## Supported Spec Types

| Type | Category | Description |
|------|----------|-------------|
| `mrd` | Source | Market Requirements Document |
| `prd` | Source | Product Requirements Document |
| `uxd` | Source | User Experience Design |
| `press` | GTM | Press Release (Working Backwards) |
| `faq` | GTM | FAQ Document |
| `narrative-1p` | GTM | 1-Page Narrative |
| `narrative-6p` | GTM | 6-Page Narrative |
| `trd` | Technical | Technical Requirements Document |
| `ird` | Technical | Infrastructure Requirements Document |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | `false` | Overwrite existing file |

## Examples

```bash
# Create a new MRD
multispec create mrd

# Create a PRD
multispec create prd

# Create a Press Release (Working Backwards)
multispec create press

# Overwrite existing TRD
multispec create trd --force
```

## Output

The command creates a file in the appropriate directory:

- Source specs: `source/{type}.md`
- GTM specs: `gtm/{type}.md`
- Technical specs: `technical/{type}.md`

## See Also

- [init](init.md) - Initialize a new project with `--with-templates`
- [synthesize](synthesize.md) - Generate specs from source documents
