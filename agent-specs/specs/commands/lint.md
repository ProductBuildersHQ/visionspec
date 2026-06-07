---
name: lint
description: Validate project structure and specification format
arguments: [--fix, --strict]
dependencies: []
---

# Lint Project

Validate the structure and format of a VisionSpec project.

## Usage

```
/lint
/lint --fix
/lint --strict
```

## Process

1. **Check Structure** - Verify visionspec.yaml and directories
2. **Validate Format** - Check YAML frontmatter and Markdown
3. **Check Consistency** - Verify cross-references resolve
4. **Report Issues** - Display errors and warnings
5. **Auto-Fix** - Optionally fix correctable issues

## Validation Checks

### Structure

- `visionspec.yaml` exists and is valid
- Required directories present
- Profile requirements met

### Format

- Valid YAML frontmatter
- Required frontmatter fields
- Markdown heading hierarchy
- Link validity

### Consistency

- Spec types match categories
- Cross-references resolve
- No orphaned files

## Output

```
Linting project: {project-name}

Structure:
  ✓ visionspec.yaml valid
  ✓ source/ directory present
  ✗ technical/ missing ird.md (required by profile)

Format:
  ✓ All specs have valid frontmatter
  ⚠ prd.md: Heading hierarchy skip (h1 → h3)

Consistency:
  ✓ All internal references resolve
  ⚠ trd.md: References FR-9 not found in PRD

Summary: 1 error, 2 warnings
```

## Flags

| Flag | Description |
|------|-------------|
| `--fix` | Automatically fix correctable issues |
| `--strict` | Treat warnings as errors |

## Auto-Fix

With `--fix`, lint automatically corrects:

- Missing frontmatter fields
- Heading hierarchy issues
- Trailing whitespace
- Missing newlines at EOF

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks pass |
| 1 | Warnings only |
| 2 | Errors found |

## CI Integration

```yaml
- name: Lint specs
  run: visionspec lint --strict
```

## Next Steps

After lint passes:

- Run `/eval <type>` to evaluate quality
- Run `/status` to see overall progress
