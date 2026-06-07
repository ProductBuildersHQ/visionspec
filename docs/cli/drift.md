# drift

Detect divergence between specifications and implementation.

## Synopsis

```bash
visionspec drift [flags]
```

## Description

The `drift` command compares specifications against the actual codebase to identify divergence. It detects:

- **Unimplemented**: Requirements in spec but not in code
- **Undocumented**: Features in code but not in spec
- **Mismatch**: Both exist but differ

## Flags

| Flag | Description |
|------|-------------|
| `-p, --project` | Project name |
| `--severity` | Filter by severity: critical, high, medium, low |
| `--format` | Output format: text, json, markdown (default: text) |
| `--ci` | CI mode: exit non-zero if drift detected |

## Examples

### Basic Drift Detection

```bash
# Detect drift in project
visionspec drift -p myproject
```

### Filter by Severity

```bash
# Show only critical and high severity drift
visionspec drift -p myproject --severity high

# Show only critical drift
visionspec drift -p myproject --severity critical
```

### Output Formats

```bash
# Text output (default)
visionspec drift -p myproject

# JSON output for programmatic use
visionspec drift -p myproject --format json

# Markdown output for documentation
visionspec drift -p myproject --format markdown > drift-report.md
```

### CI Integration

```bash
# Exit non-zero if drift detected (for CI pipelines)
visionspec drift -p myproject --ci
```

## Drift Types

### Unimplemented

Requirements specified but not found in code:

```
DRIFT-001 [high] unimplemented
  Description: User authentication endpoint not implemented
  Spec Ref: REQ-AUTH-001
  Suggestion: Implement POST /api/auth/login endpoint
```

### Undocumented

Code features not covered by specifications:

```
DRIFT-002 [medium] undocumented
  Description: Rate limiting middleware exists but not specified
  Code Ref: pkg/middleware/ratelimit.go:45
  Suggestion: Add rate limiting requirements to TRD
```

### Mismatch

Both exist but differ in behavior:

```
DRIFT-003 [high] mismatch
  Description: API returns 401 but spec says 403
  Spec Ref: API-AUTH-002
  Code Ref: pkg/handlers/auth.go:78
  Suggestion: Align error codes between spec and implementation
```

## Severity Levels

| Severity | Description |
|----------|-------------|
| `critical` | Blocking issues that must be resolved |
| `high` | Important divergence affecting functionality |
| `medium` | Notable differences that should be addressed |
| `low` | Minor inconsistencies |

## Drift Report

Example JSON output:

```json
{
  "project": "user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "items": [
    {
      "id": "DRIFT-001",
      "type": "unimplemented",
      "severity": "high",
      "category": "api",
      "description": "Missing /users/{id} DELETE endpoint",
      "spec_ref": "REQ-USER-005",
      "suggestion": "Implement user deletion endpoint"
    }
  ],
  "summary": {
    "total": 5,
    "by_type": {
      "unimplemented": 2,
      "undocumented": 2,
      "mismatch": 1
    },
    "by_severity": {
      "critical": 0,
      "high": 2,
      "medium": 2,
      "low": 1
    }
  }
}
```

## Context Integration

Drift detection uses the context infrastructure to analyze the codebase:

```bash
# Gather context first for better analysis
visionspec context gather -p myproject

# Then run drift detection
visionspec drift -p myproject
```

## See Also

- [context](context.md) - Gather codebase context
- [sync](sync.md) - Synchronize execution state
- [status](status.md) - Check project status
