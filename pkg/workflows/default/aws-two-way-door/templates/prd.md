# Product Requirements Document (PRD) - Feature

## Overview

**Feature Name:** {feature_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Draft

This PRD is the feature-level product contract, working backwards from the
OpportunitySpec, Press Release, and FAQ. Execution dates and status belong in
ROADMAP — this document states product sequencing constraints only.

## 1. OpportunitySpec Reference

| Document | Link |
|----------|------|
| OpportunitySpec | |
| Press Release | |
| FAQ | |

## 2. Problem Statement

### 2.1 User Problem

<!-- From OpportunitySpec Box 1: Users & Problem -->

### 2.2 Business Problem

<!-- From OpportunitySpec Box 5: Business Value -->

## 3. Goals and Success Metrics

### 3.1 Goals

<!-- From OpportunitySpec Box 9: Success Metrics -->

| Goal | Metric | Target |
|------|--------|--------|
| | | |

### 3.2 Non-Goals

<!-- What this feature explicitly does NOT do -->

## 4. Decision Framing

*Two-way-door decisions are made fast, at the team level, with an explicit
walk-back path — not studied for weeks. If any part of this feature is
actually hard to reverse, flag it here for escalation rather than letting it
hide inside a fast decision.*

**Reversibility:** Two-way-door <!-- or: contains a one-way-door element — escalate -->

**Walk-back path:** <!-- rollback mechanism, kill criteria, exit cost -->

**Simplifying assumptions:**

| ID | Assumption | Why It's Safe to Assume |
|----|------------|----------------------------|
| A-1 | | |

## 5. User Stories

*Every story and acceptance criterion has a stable ID.*

### 5.1 Primary User Journey

**US-1: [Story Title]**

**As a** [user type]
**I want to** [action]
**So that** [benefit]

**Acceptance Criteria:**

- [ ] **AC-1-1:** Given [precondition], when [action], then [result]
- [ ] **AC-1-2:** Given [error/boundary condition], when [action], then [result]

### 5.2 Additional User Stories

| ID | Story | Priority |
|----|-------|----------|
| US-2 | | Must Have |
| US-3 | | Should Have |

## 6. Functional Requirements

### 6.1 Core Requirements

*Every requirement traces to the story it serves.*

| ID | Requirement | Traces To | Source |
|----|-------------|-----------|--------|
| FR-1 | | US-1 | FAQ Q1 |
| FR-2 | | US-2 | Press Release |

### 6.2 Edge Cases

| Scenario | Expected Behavior |
|----------|-----------------------|
| | |

## 7. Non-Functional Requirements

| Category | Requirement |
|----------|-------------|
| Performance | |
| Security | |
| Accessibility | |

## 8. Dependencies

| Dependency | Type | Status |
|------------|------|--------|
| | Internal | |
| | External | |

## 9. Risks and Mitigations

<!-- From OpportunitySpec Box 11: Risks & Assumptions -->

| Risk | Impact | Mitigation |
|------|--------|------------|
| | | |

## 10. Release Readiness

*Product-level release conditions only — scheduling belongs in PLAN/ROADMAP.*

- **Launch gating criteria:** <!-- what must be true before this ships -->
- **Rollback trigger:** <!-- product-level condition that pulls the feature back -->

## 11. Open Questions

*No question here may block implementation at approval.*

| Question | Owner | Resolution |
|----------|-------|------------|
| | | |

## Appendix

### A. Wireframes/Mockups

### B. Technical Notes

### C. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | | | Initial draft |
