# Test Plan Document (TPD) - Feature

## Overview

**Feature Name:** {feature_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

The TPD proves the customer promise (Press Release/FAQ), the product behavior
(PRD/UXD), and the technical guarantees (TRD) all actually hold.

## 1. References

| Document | Link |
|----------|------|
| Press Release | |
| FAQ | |
| PRD | |
| UXD | |
| TRD | |

## 2. Scope

### 2.1 In Scope

<!-- What is being tested -->

### 2.2 Out of Scope

<!-- What is not being tested (covered by existing tests) -->

## 3. Requirements Traceability

### 3.1 Customer Promise Coverage

*Every Press Release / FAQ promise traces forward to a test.*

| Promise | PRD Requirement | Test Cases |
|-------------|-----------------------|------------|
| | FR-1 | TC-1, TC-2 |

### 3.2 PRD Coverage

| Requirement | Test Cases | Priority |
|-------------|------------|----------|
| FR-1 | TC-1, TC-2 | P0 |

### 3.3 TRD Coverage

| TRD ID | Test Cases | Priority |
|--------|------------|----------|
| TRD-1 | TC-3 | P0 |

### 3.4 UXD Coverage

| User Journey | Test Cases |
|--------------|------------|
| | |

## 4. Test Cases

### 4.1 Happy Path

| ID | Scenario | Steps | Expected Result |
|----|----------|-------|-----------------------|
| TC-1 | | | |

### 4.2 Error Handling

| ID | Scenario | Steps | Expected Result |
|----|----------|-------|-----------------------|
| TC-E1 | | | |

### 4.3 Edge Cases

| ID | Scenario | Steps | Expected Result |
|----|----------|-------|-----------------------|
| TC-EC1 | | | |

## 5. Test Types

### 5.1 Unit Tests

| Component | Test File | Status |
|-----------|-----------|--------|
| | | |

### 5.2 Integration Tests

| Integration | Test Approach | Status |
|-------------|---------------------|--------|
| API | | |
| Database | | |

### 5.3 E2E Tests

| User Journey | Test File | Status |
|--------------|-----------|--------|
| | | |

### 5.4 Performance Tests

| Scenario | Target | Test |
|----------|--------|------|
| Response time | | |

## 6. Reversibility and Operational Excellence Testing

*Risk-based, not unconditional — chaos testing scales with blast radius.*

**Risk tier:** {{ high | medium | low }} — <!-- rationale -->

| Test Type | Applicable? | Scenario |
|-----------|--------------|----------|
| Rollback / feature-flag off | | |
| Chaos / fault injection | Yes if high risk tier — otherwise "Not applicable — <rationale>" | |
| Observability validation (alerts fire, not just configured) | | |
| Runbook validation (executed, not just reviewed) | | |

## 7. Test Environment

| Environment | Purpose | Data |
|-------------|---------|------|
| Local | Unit tests | Mocks |
| CI | Integration | Fixtures |
| Staging | E2E | Anonymized |

## 8. Automation

### 8.1 CI Pipeline

- [ ] Unit tests on PR
- [ ] Integration tests on merge
- [ ] E2E tests on deploy to staging

### 8.2 Quality Gates

| Gate | Criteria | Blocking |
|------|----------|----------|
| PR merge | Tests pass | Yes |
| Deploy prod | E2E pass | Yes |

## 9. Regression

| Suite | Tests | Risk |
|-------|-------|------|
| | | |

## 10. Verification Sequencing

*Product-level sequencing only — dates and staffing belong in PLAN/ROADMAP.*

| Phase | Depends On |
|-------|------------------|
| Unit/Integration | TRD approved |
| E2E | Staging deploy |
| Sign-off | Exit criteria met |

## 11. Sign-Off

### 11.1 Exit Criteria

- [ ] Every row in §3 has a passing test
- [ ] No P0/P1 bugs open
- [ ] Feature flag tested on/off
- [ ] Rollback tested

### 11.2 Approvers

| Role | Name | Sign-Off |
|------|------|----------|
| QA | | |
| Engineering | | |
| Product | | |

## 12. Known Issues

| Issue | Severity | Mitigation |
|-------|----------|------------|
| | | |

## 13. Open Questions

| Question | Owner | Resolution |
|----------|-------|------------|
| | | |
