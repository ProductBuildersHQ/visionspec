# VisionSpec Lint

Validate project structure and specification format.

## Overview

The lint command ensures specifications adhere to expected structure, naming conventions, and format requirements. Use it before commits, after manual edits, or as part of CI/CD.

## Workflow

### 1. Run Lint

```
visionspec lint
```

Or with auto-fix:

```
visionspec lint --fix
```

### 2. Review Results

Output is organized by category:

**Structure**

- visionspec.yaml validity
- Required directories present
- Profile requirements met

**Format**

- Valid YAML frontmatter
- Required frontmatter fields
- Markdown heading hierarchy
- Link validity

**Consistency**

- Spec types match categories
- Cross-references resolve
- No orphaned files

### 3. Address Issues

For each issue:

| Severity | Action |
|----------|--------|
| Error (✗) | Must fix before proceeding |
| Warning (⚠) | Should fix, but not blocking |

### 4. Auto-Fix

With `--fix`, lint automatically corrects:

- Missing frontmatter fields
- Heading hierarchy issues
- Trailing whitespace
- Missing newlines at EOF

## CI Integration

Add to your pipeline:

```yaml
- name: Lint specs
  run: visionspec lint --strict
```

Exit codes:

- `0`: All checks pass
- `1`: Warnings only
- `2`: Errors found

## Common Issues

### Missing visionspec.yaml

Create with:

```
visionspec init
```

### Invalid Frontmatter

Ensure YAML is valid:

```yaml
---
title: My PRD
version: 1.0
status: draft
---
```

### Heading Hierarchy

Don't skip levels:

```markdown
# H1
## H2  (correct)
### H3 (correct)

# H1
### H3 (incorrect - skipped H2)
```

## Related Skills

- `visionspec-status`: See overall project status
- `visionspec-eval`: Quality evaluation beyond linting
