# Technical Requirements Document (TRD) - Feature

## Overview

**Feature Name:** {feature_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

## 1. References

| Document | Link |
|----------|------|
| OpportunitySpec | |
| PRD | |

## 2. Architecture

### 2.1 Context

<!-- How does this feature fit into the existing system? -->

```
┌─────────────────────────────────────────────────────────────┐
│                    Existing System                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   [Existing A] ──────► [NEW FEATURE] ──────► [Existing B]   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Component Design

**New Components:**

| Component | Purpose | Changes |
|-----------|---------|---------|
| | | New |
| | | Modified |

### 2.3 Data Flow

<!-- How does data flow through the new feature? -->

## 3. API Changes

### 3.1 New Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| | | |

### 3.2 Request/Response

```json
// POST /api/v1/feature
// Request
{
  "field": "value"
}

// Response
{
  "id": "123",
  "status": "success"
}
```

### 3.3 Modified Endpoints

| Endpoint | Change | Backward Compatible |
|----------|--------|---------------------|
| | | Yes/No |

## 4. Data Changes

### 4.1 Schema Changes

| Entity | Change | Migration Required |
|--------|--------|-------------------|
| | | |

### 4.2 Data Migration

<!-- If migration needed, describe approach -->

## 5. Non-Functional Requirements

### 5.1 Performance Impact

| Metric | Current | Expected | Acceptable |
|--------|---------|----------|------------|
| Latency (p50) | | | |
| Latency (p99) | | | |

### 5.2 Security Considerations

- [ ] Authentication changes
- [ ] Authorization changes
- [ ] Data handling changes

## 6. Feature Flags

| Flag | Purpose | Default |
|------|---------|---------|
| | | off |

## 7. Rollout Plan

### 7.1 Phases

| Phase | Audience | Criteria to Proceed |
|-------|----------|---------------------|
| 1. Internal | Dogfood | No P0 bugs |
| 2. Beta | 5% users | Success metrics met |
| 3. GA | 100% | |

### 7.2 Rollback Plan

<!-- How to rollback if issues found -->

## 8. Dependencies

| Dependency | Type | Risk |
|------------|------|------|
| | | |

## 9. Testing Strategy

### 9.1 Unit Tests

<!-- Key unit test scenarios -->

### 9.2 Integration Tests

<!-- Key integration points to test -->

### 9.3 E2E Tests

<!-- Critical user journeys to test -->

## 10. Open Questions

| Question | Owner | Status |
|----------|-------|--------|
| | | |
