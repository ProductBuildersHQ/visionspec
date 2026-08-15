# Technical Requirements Document (TRD) - Feature

## Overview

**Feature Name:** {feature_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

This TRD is the feature-level technical contract: what the change MUST
guarantee, scoped to a two-way-door feature landing inside an existing
system. Use MUST / MUST NOT / SHOULD / MAY consistently.

## 1. References

| Document | Link |
|----------|------|
| OpportunitySpec | |
| PRD | |

## 2. Technical Requirements and Traceability

| TRD ID | Requirement | Traces To (PRD) | Verified By (TPD) |
|--------|-------------|----------------------|-------------------------|
| TRD-1 | | FR-1 | TC-1 |

## 3. Context and Invariants

### 3.1 System Context

<!-- How does this feature fit into the existing system? -->

```
┌─────────────────────────────────────────────────────────────┐
│                    Existing System                           │
├─────────────────────────────────────────────────────────────┤
│   [Existing A] ──────► [NEW FEATURE] ──────► [Existing B]   │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 Invariants

*What must hold true regardless of implementation detail — this is what
lets an implementer fill gaps safely instead of guessing.*

- <!-- e.g., "Existing API consumers MUST see no behavior change unless they opt in." -->

## 4. Component Design

**New Components:**

| Component | Purpose | Change |
|-----------|---------|--------|
| | | New |
| | | Modified |

### 4.1 Data Flow

<!-- How does data flow through the new feature? -->

## 5. API Changes

### 5.1 New / Modified Endpoints

| Endpoint | Method | Backward Compatible |
|----------|--------|---------------------------|
| | | Yes/No |

### 5.2 Request/Response

```json
// POST /api/v1/feature
// Request
{ "field": "value" }
// Response
{ "id": "123", "status": "success" }
```

### 5.3 Idempotency and Error Semantics

<!-- Retry safety, error codes, and their retryability -->

## 6. Data Changes

### 6.1 Schema Changes

| Entity | Change | Migration Required |
|--------|--------|---------------------------|
| | | |

### 6.2 Migration and Compatibility

**Applicability:** {{ "Applicable" | "Not applicable — no schema/data changes" }}

<!-- If migrating: phases, validation, and the point past which rollback is
     no longer safe (use roll-forward instead). -->

## 7. Non-Functional Requirements

### 7.1 Performance Impact

| Metric | Current | Expected | Acceptable |
|--------|---------|----------|------------|
| Latency (p50 / p99) | | | |

### 7.2 Security Considerations

- [ ] Authentication changes
- [ ] Authorization changes
- [ ] Data handling changes

## 8. Feature Flags

| Flag | Purpose | Default |
|------|---------|---------|
| | | off |

## 9. Rollout Plan

### 9.1 Phases

| Phase | Audience | Criteria to Proceed |
|-------|----------|-----------------------------|
| 1. Internal | Dogfood | No P0 bugs |
| 2. Beta | 5% users | Success metrics met |
| 3. GA | 100% | |

### 9.2 Rollback Plan

<!-- How to rollback if issues found; if data changes make rollback unsafe
     past a point, state the roll-forward alternative (see §6.2). -->

## 10. Dependencies

| Dependency | Type | Risk |
|------------|------|------|
| | | |

## 11. Testing Strategy

<!-- Unit, integration, and E2E approach. TPD carries the full verification plan. -->

## 12. Open Questions

*No question here may block implementation at approval.*

| Question | Owner | Resolution |
|----------|-------|------------|
| | | |
