# align

Check alignment between specifications and current implementation truth.

## Synopsis

```bash
visionspec align [flags]
```

## Description

The `align` command verifies that specifications match the current state of the implementation. Unlike `drift` which detects divergence, `align` performs a structured comparison against a documented "current truth" to identify discrepancies that need resolution.

## Flags

| Flag | Description |
|------|-------------|
| `-p, --project` | Project name |
| `--with-context` | Include codebase context in analysis |
| `--context-file` | Path to pre-gathered context file |
| `--format` | Output format: text, json, markdown (default: text) |
| `--ci` | CI mode: exit non-zero if misaligned |

## Examples

### Basic Alignment Check

```bash
# Check alignment for project
visionspec align -p myproject
```

### With Codebase Context

```bash
# Include live codebase context
visionspec align -p myproject --with-context

# Use pre-gathered context file
visionspec align -p myproject --context-file context.json
```

### Output Formats

```bash
# Text output (default)
visionspec align -p myproject

# JSON output for programmatic use
visionspec align -p myproject --format json

# Markdown output for documentation
visionspec align -p myproject --format markdown > alignment-report.md
```

### CI Integration

```bash
# Exit non-zero if misaligned (for CI pipelines)
visionspec align -p myproject --ci
```

## Discrepancy Types

### Missing Feature

Specified in requirements but not present:

```
ALIGN-001 [high] missing_feature
  Description: OAuth2 login not implemented
  Spec Ref: REQ-AUTH-003
  Suggestion: Implement OAuth2 authentication flow
```

### Undocumented Code

Implementation exists but not specified:

```
ALIGN-002 [medium] undocumented_code
  Description: Session timeout logic exists but not specified
  Code Ref: pkg/auth/session.go:120
  Suggestion: Add session management requirements to TRD
```

### Diverged

Both exist but behavior differs:

```
ALIGN-003 [high] diverged
  Description: Password requirements differ from spec
  Spec Ref: REQ-AUTH-001
  Code Ref: pkg/auth/validation.go:45
  Suggestion: Align password validation rules
```

## Alignment Report

Example JSON output:

```json
{
  "project": "user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "discrepancies": [
    {
      "type": "missing_feature",
      "spec_ref": "REQ-AUTH-003",
      "severity": "high",
      "description": "OAuth2 login not implemented"
    }
  ],
  "summary": {
    "total": 3,
    "aligned": 45,
    "discrepancies": 3,
    "alignment_percentage": 93.75
  }
}
```

## Current Truth

The alignment check compares against a `current-truth.md` file that documents the actual implementation state:

```markdown
# Current Truth

## Implemented Features
- [x] Basic authentication
- [x] Password reset
- [ ] OAuth2 login (in progress)

## API Endpoints
- POST /api/auth/login - implemented
- POST /api/auth/logout - implemented
- POST /api/auth/oauth - not implemented
```

## See Also

- [drift](drift.md) - Detect spec-to-code drift
- [context](context.md) - Gather codebase context
- [status](status.md) - Check project status
