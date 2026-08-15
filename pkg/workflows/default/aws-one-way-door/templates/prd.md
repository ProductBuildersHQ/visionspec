# Product Requirements Document (PRD) - Product

## Overview

**Product Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Draft

This PRD is the product-level contract for a one-way-door launch, working
backwards from the human-authored Press Release, stress-tested by the FAQ, and
deepened by the MRD where produced. Build spend is gated on the 6-pager
decision meeting — this document must be complete enough to support that
decision. Execution dates and staffing belong in PLAN/ROADMAP; this document
states product requirements and sequencing constraints only.

## 1. Working Backwards References

| Document | Link |
|----------|------|
| Press Release | |
| FAQ | |
| MRD | *or "Not produced — <rationale>" (MRD is optional; say why the FAQ's business case suffices)* |
| 6-Pager (decision narrative) | |

## 2. Problem Statement

### 2.1 Customer Problem

<!-- From the Press Release problem paragraph: who is the customer and what
     pain does the announced experience remove? Write it so the customer would
     recognize themselves. -->

### 2.2 Business Problem

<!-- From the FAQ business-case answers: market, economics, why now. Cite the
     FAQ question IDs that carry the evidence. -->

## 3. One-Way-Door Commitments

*The reason this ceremony exists. Name every commitment that cannot be walked
back once launched — a public API contract, a pricing promise, a brand or
product name, a data model or residency choice, a partner agreement. For each,
the deliberation must be proportionate to the stakes, and the evidence of that
deliberation must be linked, not asserted. Keep this list as small as
possible: a good one-way-door product design makes the irreversible core
narrow and keeps everything else reversible.*

### 3.1 Irreversible Commitments

| ID | Commitment | Why It Cannot Be Walked Back | Blast Radius If Wrong | Exit Cost If Forced | Deliberation Evidence |
|----|------------|------------------------------|------------------------|----------------------|------------------------|
| OWD-1 | | | | | 6-pager §6 / decision meeting |

**Personal-money test:** <!-- Would you make these commitments with your own
money, at this scope, on this timeline? Answer honestly — this is the
Ownership check, not a formality. -->

### 3.2 Reversible Sub-Decisions

*Everything inside the launch that IS reversible gets decided fast at team
level — do not spend one-way-door scrutiny on two-way-door choices.*

| ID | Decision | Owner | Walk-Back Path |
|----|----------|-------|----------------|
| RSD-1 | | | |

## 4. Goals and Success Metrics

*Launch is the starting line, not the finish line. Metrics focus on inputs
the team controls, with baselines and targets, over a horizon long enough to
matter — a one-way-door bet is not judged at day 30.*

### 4.1 Goals

| Goal | Input Metric | Baseline | 90-Day Target | Year-1 Target |
|------|--------------|----------|----------------|----------------|
| | | | | |

### 4.2 Non-Goals

<!-- What this product explicitly does NOT do. At product scale, non-goals
     prevent the launch from absorbing every adjacent idea. -->

## 5. Target Users and Market Grounding

### 5.1 Primary Persona

<!-- From the Press Release customer and MRD segments where present. Be
     specific enough that this customer would recognize themselves. -->

### 5.2 Market Evidence

**MRD status:** {{ "Produced — see references" | "Not produced — <why the FAQ's business case suffices at this risk level>" }}

| Claim | Evidence Source (MRD § / FAQ Q) |
|-------|----------------------------------|
| | |

## 6. User Stories

*Every story and acceptance criterion has a stable ID. Acceptance criteria
cover happy path, error, empty, permission, and boundary behavior — not only
the primary flow.*

### Epic: [Epic Name]

#### US-001: [Story Title]

**As a** [user type]
**I want** [goal]
**So that** [benefit]

**Acceptance Criteria:**

- [ ] **AC-001-1:** Given [precondition], when [action], then [result]
- [ ] **AC-001-2:** Given [error/boundary condition], when [action], then [result]

**Dependencies:** None

---

#### US-002: [Story Title]

**As a** [user type]
**I want** [goal]
**So that** [benefit]

**Acceptance Criteria:**

- [ ] **AC-002-1:** Given [precondition], when [action], then [result]

**Dependencies:** [e.g., US-001]

---

## 7. Functional Requirements

*Every requirement traces to the story it serves and to the Working Backwards
source that justifies it (a Press Release promise or an FAQ answer). A
requirement with neither is scope creep — justify it or cut it.*

| ID | Requirement | Priority | Traces To | Source |
|----|-------------|----------|-----------|--------|
| FR-001 | | P0 | US-001 | Press promise / FAQ Q# |

## 8. Non-Functional Requirements

| Category | Requirement |
|----------|-------------|
| Performance | |
| Scalability | <!-- at the Press Release's think-big scale, not current traffic --> |
| Security | |
| Accessibility | |
| Reliability | |

## 9. Dependencies and Risks

### 9.1 Dependencies

| Dependency | Type | Owner | Status |
|------------|------|-------|--------|
| | Internal / External | | |

### 9.2 Risks

<!-- From the FAQ's hardest questions. The risks the FAQ surfaced do not
     disappear because the PRD is written — carry them forward with
     mitigations. -->

| Risk | Impact | Mitigation | Source (FAQ Q#) |
|------|--------|------------|------------------|
| | | | |

## 10. Release Strategy

*Product-level release conditions only — scheduling belongs in PLAN/ROADMAP.
Build spend is gated on the decision meeting; nothing in this section
authorizes work before that gate passes.*

- **Decision gate:** 6-pager decision meeting approved on <!-- date/link -->
- **Launch gating criteria:** <!-- what must be true before the public commitment is made -->
- **Pre-launch validation:** <!-- for the OWD-N commitments there is no rollback — what proof completes before launch (see TPD) -->
- **Phasing:** <!-- private preview → limited GA → GA, if applicable; which OWD-N commitments become irrevocable at which phase -->

## 11. Open Questions

*No question here may block the decision meeting. A question that would
change an OWD-N row is not an open question — it is an unfinished PRD.*

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
