# speckit

Export a VisionSpec workflow family as native GitHub Spec Kit extensions and a workflow.

## Usage

```bash
visionspec speckit export [workflow] [flags]
```

## Description

The `speckit` command group packages a whole VisionSpec workflow methodology — templates, rubrics, review gates, and phase ordering — as installable [Spec Kit](https://github.com/github/spec-kit) plugins, for use in **any** Spec Kit project. This is a different mechanism from [`export speckit`](export.md#speckit): that target turns an already-reconciled VisionSpec project's `spec.md` into `plan.md`/`tasks.md`; `speckit export` turns a workflow *family* (templates, rubrics, gates — not a specific project's specs) into a Spec Kit extension anyone can install, whether or not they use VisionSpec at all.

`export` splits the family's execution sequence by spec category into a product extension (`source`/`gtm` specs) and, when present, an engineering add-on (`technical` specs). Each spec type becomes one Spec Kit command that authors its document from the family's template and self-evaluates against the family's rubric. It also compiles the family's full sequence and review gates into a single workflow.yml.

## Arguments

| Argument | Description |
|----------|-------------|
| `workflow` | Workflow family to export (default: `aws-one-way-door`) |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output`, `-o` | string | `speckit` | Output directory for generated artifacts |
| `--check` | bool | `false` | Verify committed artifacts match the generator instead of writing (CI drift gate) |
| `--product-id` | string | `""` | Extension id for the product chain (default: family name without `aws-` prefix) |
| `--engineering-id` | string | `""` | Extension id for the engineering add-on (default: `<product-id>-engineering`) |
| `--ext-version` | string | `""` | Extension version `X.Y.Z` (default: `0.1.0`) |

## Examples

```bash
# Export the default family into ./speckit
visionspec speckit export

# Verify committed artifacts are in sync (CI drift gate)
visionspec speckit export --check

# Export another family with a custom extension id
visionspec speckit export big-tech -o ./speckit --product-id big-tech
```

## Output

```
Created speckit/extensions/one-way-door/extension.yml
Created speckit/extensions/one-way-door/commands/speckit.one-way-door.press.md
Created speckit/extensions/one-way-door/templates/press.md
...
Created speckit/workflows/one-way-door-chain/workflow.yml
```

## Installing the Result

The generated directory installs into any Spec Kit project with the `specify` CLI:

```bash
specify extension add speckit/extensions/one-way-door --dev
specify workflow add speckit/workflows/one-way-door-chain --dev
```

See the [Spec Kit Plugins guide](../guides/speckit-plugins.md) for a full walkthrough, including the companion hand-authored `risk-adjusted` triage workflow and role bundles that aren't produced by this command (they're maintained directly under `speckit/`).

## See Also

- [Spec Kit Plugins guide](../guides/speckit-plugins.md) - full walkthrough: install, run, and hand off to Spec Kit's core flow
- [export](export.md) - the `spec.md`/`plan.md`/`tasks.md` exporter for a reconciled project
- [profiles](profiles.md) - list and inspect workflow families available to export
