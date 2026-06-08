# render

Render an evaluation file to markdown format.

## Usage

```bash
visionspec render <eval-file> [flags]
```

## Description

Render an existing evaluation JSON file to markdown format. This is useful for viewing evaluation results in a readable format or for generating documentation from past evaluations.

## Arguments

| Argument | Description |
|----------|-------------|
| `eval-file` | Path to evaluation JSON file(s) |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | stdout | Output file (default: stdout or `<input>.md`) |
| `--evaluation` | | false | Render structured evaluation report format |

## Examples

```bash
# Render PRD evaluation to stdout
visionspec render eval/prd.eval.json

# Output to file
visionspec render eval/prd.eval.json -o report.md

# Render all evaluations
visionspec render eval/*.eval.json

# Render as structured evaluation report
visionspec render eval/prd.eval.json --evaluation
```

## Output Format

The rendered markdown includes:

- Overall score and pass/fail status
- Category-by-category breakdown
- Detailed findings with severity levels
- Recommendations for improvement

## See Also

- [eval](eval.md) - Run evaluations
- [status](status.md) - Show project status including eval results
