---
name: visionspec-lint
description: Validate project structure and specification format
triggers: [lint, validate, check structure, verify format]
---

# VisionSpec Lint

Validate the structure and format of a VisionSpec project.

## Purpose

Ensures specifications adhere to:

- Expected directory structure
- Required file naming conventions
- Valid YAML frontmatter
- Correct Markdown formatting
- Profile requirements

## When to Use

- Before committing changes
- After manual spec edits
- To troubleshoot eval/reconcile failures
- As part of CI/CD pipeline

## Invocation

```
visionspec lint
visionspec lint --fix
```

Or via Claude Code:

```
/lint
```

## Validation Checks

### Structure

- `visionspec.yaml` exists and is valid
- Profile directory exists if specified
- All required spec directories present

### Format

- Frontmatter is valid YAML
- Required frontmatter fields present
- Markdown headings follow hierarchy
- Links are valid (internal and external)

### Consistency

- Spec types match expected categories
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

## Auto-Fix

With `--fix`, lint can automatically correct:

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

## Integration

Use in CI:

```yaml
- name: Lint specs
  run: visionspec lint --strict
```
