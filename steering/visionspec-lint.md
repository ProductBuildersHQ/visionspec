# Steering: visionspec-lint

## Command

```
visionspec lint
visionspec lint --fix
```

## Purpose

Validate project structure and specification format.

## When to Invoke

- Before committing changes
- After manual spec edits
- When eval or reconcile fails unexpectedly
- User asks to validate or check format
- As part of pre-commit workflow

## Triggers

- "validate specs"
- "check format"
- "lint"
- "verify structure"
- "fix formatting"

## Expected Output

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

## Auto-Fix

With `--fix`, automatically corrects:

- Missing frontmatter fields
- Heading hierarchy issues
- Trailing whitespace
- Missing newlines at EOF

## Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| 0 | Pass | Proceed |
| 1 | Warnings | Review, can proceed |
| 2 | Errors | Must fix |

## Follow-up Actions

| Issue | Resolution |
|-------|------------|
| Missing visionspec.yaml | Run `visionspec init` |
| Invalid frontmatter | Fix YAML syntax |
| Heading skip | Use `--fix` or manual edit |
| Missing references | Update cross-references |

## CI Integration

```yaml
- name: Lint specs
  run: visionspec lint --strict
```

## Error Handling

Common issues:

- **No visionspec.yaml**: Project not initialized
- **YAML parse error**: Invalid frontmatter syntax
- **Missing directory**: Create required directories
