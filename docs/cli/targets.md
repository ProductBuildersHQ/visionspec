# targets

List available export targets and their capabilities.

## Usage

```bash
visionspec targets
```

## Description

The `targets` command displays all available export targets that can be used with `visionspec export`.

## Output

```
Available targets:

  speckit   GitHub Spec-Kit format (spec.md, plan.md, tasks.md)
  gsd       Get Shit Done format (PLAN.md, STATE.md)
  gastown   GasTown formulas and beads
  gascity   GasCity city.toml configuration
  aidlc     AWS AI-DLC Workflows format (vision-document.md, technical-environment.md, imported-requirements.md)
  openspec  OpenSpec portable format (future)
```

## Target Comparison

| Feature | SpecKit | GSD | GasTown | GasCity | AIDLC |
|---------|---------|-----|---------|---------|-------|
| Sequential tasks | Yes | Yes | Yes | Yes | Yes |
| Parallel execution | No | Yes (waves) | Yes (convoy) | Yes | No |
| Multi-agent | No | No | Yes | Yes | Yes (specialists) |
| Verification | Implicit | Yes | Yes | Yes | Yes (gates) |
| Dependency graph | Yes | Yes | Yes (Beads) | Yes | Yes |

## See Also

- [export](export.md) - Export to a target system
