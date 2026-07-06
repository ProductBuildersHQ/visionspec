# Reason Codes

Reason codes are standardized identifiers for evaluation findings that enable automated repair workflows with AI assistance.

## Overview

When VisionSpec evaluates a specification, it identifies issues and assigns each a **reason code**. These codes follow a structured format that enables:

- **Categorization** - Quickly understand what type of issue was found
- **Prioritization** - Route issues to the right team or repair strategy
- **Automation** - AI agents can use repair prompts to fix issues automatically
- **Analytics** - Track common issues across projects and teams

## Code Format

Reason codes use a `{CATEGORY}-{ISSUE}` format:

```
REQ-AMBIGUOUS       # Requirements: Ambiguous requirement
SEC-NO_AUTH         # Security: Missing authentication
UX-NO_ARIA          # UX: Missing ARIA labels
METRIC-UNMEASURABLE # Metrics: Unmeasurable metric
```

## Categories

| Prefix | Domain | Description |
|--------|--------|-------------|
| `REQ-` | Requirements | Issues with requirement clarity, completeness, testability |
| `METRIC-` | Metrics | Issues with success metrics, KPIs, measurement |
| `USER-` | User/Persona | Issues with user personas, journeys, problem statements |
| `ARCH-` | Architecture | Issues with technical design, APIs, data models |
| `SEC-` | Security | Security gaps, auth, privacy, validation |
| `SCALE-` | Scalability | Performance, capacity, single points of failure |
| `INFRA-` | Infrastructure | Deployment, monitoring, recovery, operations |
| `DOC-` | Documentation | Missing diagrams, examples, inconsistencies |
| `SCOPE-` | Scope | Scope creep, unbounded scope, missing constraints |
| `UX-` | UX/Accessibility | Accessibility, error states, responsive design |

## Requirements Codes (REQ-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `REQ-AMBIGUOUS` | Requirement lacks specificity or has multiple interpretations | Medium |
| `REQ-NO_CRITERIA` | Requirement lacks acceptance criteria | High |
| `REQ-CONFLICT` | Two or more requirements contradict | Critical |
| `REQ-INCOMPLETE` | Requirement missing essential details | Medium |
| `REQ-UNTESTABLE` | Requirement cannot be objectively verified | Medium |
| `REQ-MISSING_REASON` | Requirement lacks business rationale | Low |

## Metrics Codes (METRIC-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `METRIC-UNMEASURABLE` | Success metric cannot be objectively measured | Medium |
| `METRIC-NO_BASELINE` | Metric lacks baseline value | Low |
| `METRIC-NO_TARGET` | Metric lacks target value | Medium |
| `METRIC-UNREALISTIC` | Target appears unrealistic | Medium |
| `METRIC-NO_TRACKING` | No plan for how metrics will be tracked | Low |
| `METRIC-MISSING_KPI` | Critical feature missing KPI | Medium |
| `METRIC-VANITY` | Metric doesn't correlate with business value | Low |

## User Codes (USER-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `USER-NO_PERSONA` | Target user persona not defined | High |
| `USER-INCOMPLETE` | Persona lacks essential details | Medium |
| `USER-NO_JOURNEY` | User journey not documented | Medium |
| `USER-UNCLEAR_PROBLEM` | Problem statement vague or unclear | High |
| `USER-NO_GOALS` | User goals not articulated | Medium |
| `USER-NO_PAIN_POINTS` | Pain points not documented | Medium |

## Architecture Codes (ARCH-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `ARCH-NO_ERROR_HANDLING` | Error handling strategy incomplete | Medium |
| `ARCH-NO_API` | API contract not specified | High |
| `ARCH-NO_DATA_MODEL` | Data model not defined | High |
| `ARCH-MISSING_DEP` | Required dependency not documented | Medium |
| `ARCH-GAP` | Architecture has unexplained gaps | High |
| `ARCH-NO_INTERFACE` | Interface between components not defined | Medium |
| `ARCH-CIRCULAR_DEP` | Circular dependency detected | High |
| `ARCH-TIGHT_COUPLING` | Components too tightly coupled | Medium |

## Security Codes (SEC-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `SEC-GAP` | Security consideration not addressed | High |
| `SEC-NO_AUTH` | Authentication mechanism not specified | Critical |
| `SEC-NO_AUTHZ` | Authorization model not defined | High |
| `SEC-PRIVACY` | Data privacy requirements not addressed | High |
| `SEC-NO_VALIDATION` | Input validation not specified | Medium |
| `SEC-NO_ENCRYPTION` | Encryption requirements not specified | High |
| `SEC-HARDCODED_SECRET` | Hardcoded secrets detected | Critical |
| `SEC-INJECTION_RISK` | Potential injection vulnerability | Critical |

## Scalability Codes (SCALE-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `SCALE-CONCERN` | Scalability concern not addressed | Medium |
| `SCALE-PERFORMANCE` | Performance risk not mitigated | Medium |
| `SCALE-NO_CAPACITY` | Capacity planning not documented | Low |
| `SCALE-SPOF` | Single point of failure identified | High |
| `SCALE-NO_RATE_LIMIT` | Rate limiting not specified | Medium |
| `SCALE-NO_CACHE` | Caching strategy not defined | Low |
| `SCALE-BLOCKING_OP` | Blocking operation in critical path | Medium |

## Infrastructure Codes (INFRA-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `INFRA-NO_DEPLOY` | Deployment strategy not defined | Medium |
| `INFRA-NO_MONITOR` | Monitoring strategy not defined | Medium |
| `INFRA-NO_ALERT` | Alerting strategy not defined | Low |
| `INFRA-NO_RECOVERY` | Disaster recovery plan not documented | Medium |
| `INFRA-NO_BACKUP` | Backup strategy not specified | Medium |
| `INFRA-NO_RUNBOOK` | Operational runbook not provided | Low |
| `INFRA-NO_ROLLBACK` | Rollback procedure not defined | Medium |
| `INFRA-ENV_MISMATCH` | Environment configuration mismatch risk | Medium |

## Documentation Codes (DOC-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `DOC-INSUFFICIENT` | Documentation insufficient for implementation | Medium |
| `DOC-OUTDATED` | Reference to outdated information | Low |
| `DOC-NO_DIAGRAM` | Visual diagram would improve clarity | Low |
| `DOC-NO_EXAMPLES` | Missing examples to clarify usage | Low |
| `DOC-INCONSISTENT` | Documentation inconsistent with other sections | Medium |

## Scope Codes (SCOPE-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `SCOPE-CREEP` | Scope includes unnecessary features | Medium |
| `SCOPE-UNBOUNDED` | Scope is unbounded or unclear | High |
| `SCOPE-NO_CONSTRAINTS` | Constraints not documented | Medium |
| `SCOPE-NO_NON_GOALS` | Non-goals not explicitly stated | Low |
| `SCOPE-MVP_UNCLEAR` | MVP scope not clearly defined | Medium |
| `SCOPE-NO_TIMELINE` | Timeline or milestones not specified | Low |

## UX Codes (UX-*)

| Code | Description | Default Severity |
|------|-------------|------------------|
| `UX-NO_ARIA` | ARIA labels not specified | High |
| `UX-NO_ERROR_STATE` | Error state UI not designed | High |
| `UX-NO_LOADING` | Loading state not specified | Medium |
| `UX-NO_EMPTY` | Empty state not designed | Low |
| `UX-NO_RESPONSIVE` | Responsive behavior not specified | Medium |
| `UX-NO_KEYBOARD` | Keyboard navigation not specified | Medium |
| `UX-INCOMPLETE_NAV` | Navigation flow incomplete | Medium |
| `UX-NO_FEEDBACK` | User feedback mechanism not designed | Low |

## AI-Assisted Repair

Each reason code includes a **repair prompt** - a detailed instruction that AI agents can use to automatically fix the issue.

### Example: REQ-NO_CRITERIA

**Finding:**
```json
{
  "code": "REQ-NO_CRITERIA",
  "message": "FR-2.3 lacks acceptance criteria",
  "location": "FR-2.3"
}
```

**Repair Prompt:**
> Add acceptance criteria using Given/When/Then format. Include: (1) preconditions (Given), (2) action trigger (When), (3) expected outcome (Then). Cover both success and failure scenarios.

### Human Review Required

Some codes are flagged as requiring human review after AI repair:

| Category | Requires Human Review |
|----------|----------------------|
| Security codes (`SEC-*`) | Yes - security decisions need expert review |
| Privacy (`SEC-PRIVACY`) | Yes - compliance implications |
| Scope decisions (`SCOPE-CREEP`, `SCOPE-UNBOUNDED`) | Yes - business judgment needed |
| Architecture gaps (`ARCH-GAP`) | Yes - design decisions |
| Timelines (`SCOPE-NO_TIMELINE`) | Yes - PM input needed |

Codes that can be fixed autonomously:
- `REQ-AMBIGUOUS` - AI can add specificity
- `DOC-NO_DIAGRAM` - AI can generate diagrams
- `UX-NO_LOADING` - AI can add loading state specs

## Evaluation Result Format (v2)

Evaluation results include reason codes in findings:

```json
{
  "schemaVersion": "v2",
  "scoreV2": 4,
  "decision": "conditional",
  "pass": false,
  "confidence": 0.78,
  "blocking": ["REQ-NO_CRITERIA", "METRIC-UNMEASURABLE"],
  "dimensions": [
    {
      "id": "requirements",
      "name": "Requirements Clarity",
      "score": 3,
      "confidence": 0.75,
      "reasonCodes": ["REQ-NO_CRITERIA"],
      "findings": [...]
    }
  ],
  "findings": [
    {
      "category": "requirements",
      "severity": "medium",
      "code": "REQ-NO_CRITERIA",
      "message": "FR-2.3 lacks acceptance criteria",
      "location": "FR-2.3"
    }
  ]
}
```

## Backwards Compatibility

Legacy codes (e.g., `AMBIGUOUS_REQUIREMENT`) are automatically normalized to the new format (`REQ-AMBIGUOUS`) for backwards compatibility.

## See Also

- [eval](../cli/eval.md) - Evaluate specifications
- [Concepts Overview](./index.md) - Core VisionSpec concepts
