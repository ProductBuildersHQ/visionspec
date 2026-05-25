# targets

List available export targets and their capabilities.

## Usage

```bash
multispec targets
```

## Description

The `targets` command displays all available export targets that can be used with `multispec export`.

## Output

```
Available targets:

  speckit   GitHub Spec-Kit format (spec.md, plan.md, tasks.md)
  gsd       Get Shit Done format (PLAN.md, STATE.md)
  gastown   GasTown formulas and beads
  gascity   GasCity city.toml configuration
  openspec  OpenSpec portable format (future)
```

## Target Comparison

| Feature | SpecKit | GSD | GasTown | GasCity |
|---------|---------|-----|---------|---------|
| Sequential tasks | Yes | Yes | Yes | Yes |
| Parallel execution | No | Yes (waves) | Yes (convoy) | Yes |
| Multi-agent | No | No | Yes | Yes |
| Verification | Implicit | Yes | Yes | Yes |
| Dependency graph | Yes | Yes | Yes (Beads) | Yes |

## See Also

- [export](export.md) - Export to a target system
