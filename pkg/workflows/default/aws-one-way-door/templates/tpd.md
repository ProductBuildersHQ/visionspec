# Test Plan Document (TPD) - Product

## Overview

**Product Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

The TPD proves the customer promise (Press Release/FAQ), the product behavior
(PRD/UXD), and the technical guarantees (TRD/IRD) all actually hold — and it
proves them **before launch**. For the one-way-door commitments there is no
rollback: validation that would normally happen safely in production must
complete pre-launch, at launch scale, or the commitment is not ready to make.

## 1. References

| Document | Link |
|----------|------|
| Press Release | |
| FAQ | |
| PRD | |
| UXD | |
| TRD | |
| IRD | *or "Not applicable — see IRD declaration"* |

## 2. Promise and Requirement Traceability

*The full chain, one table. Every Press Release promise traces forward to an
automated test or eval; every row without one is a verification gap, not an
oversight to fix later.*

| Chain | ID | Verified By (TC ID) |
|-------|-----|----------------------|
| PRESS promise | | |
| → PRD behavior | FR-001 | |
| → UXD behavior | | |
| → TRD/IRD guarantee | TRD-001 | |
| → TPD scenario | | |

## 3. Pre-Launch Validation of Irreversible Commitments

*One row per PRD OWD-N / TRD TWD-N commitment. Each proof is executed before
the commitment becomes irrevocable — "we'll watch it in production" is not a
plan when there is no way back.*

| Commitment (OWD/TWD ID) | Pre-Launch Proof | Executed At Scale? | Test Case IDs | Status |
|--------------------------|-------------------|---------------------|---------------|--------|
| OWD-1 / TWD-1 | *e.g., contract tests against frozen v1 spec; load test at launch projection; migration rehearsal on production-like data* | | | |

**Launch-scale definition:** <!-- the think-big projection these proofs run
against, from TRD §9 — state the number, not "high load" -->

## 4. Reversibility Containment Testing

*The reversible scope must actually be reversible. Flags, rollbacks, and
kill criteria for §3.2/TRD-RSD decisions are tested, not assumed — an
untested rollback is a one-way door you didn't mean to build.*

| Reversible Decision (RSD ID) | Test | Rollback Executed? | Test Case IDs |
|-------------------------------|------|---------------------|---------------|
| RSD-1 | Feature flag on/off; state consistent after off | | |

## 5. Coverage Matrices

*ID-based. A requirement, guarantee, or user-visible state with no test case
is a gap — surface it here, don't let it hide.*

### 5.1 PRD Coverage

| PRD ID | Test Case ID | Description | Priority |
|--------|--------------|-------------|----------|
| FR-001 | TC-001 | | High |

### 5.2 TRD Coverage

| TRD ID | Test Case ID | Description | Priority |
|--------|--------------|-------------|----------|
| TRD-001 | TC-NFR-001 | | High |

### 5.3 IRD Coverage

**Applicability:** {{ "Applicable" | "Not applicable — <rationale>" }}

| IRD ID | Test Case ID | Description | Priority |
|--------|--------------|-------------|----------|
| | | | |

### 5.4 UXD Coverage

| UXD State / Flow | Test Case ID | E2E or UAT | Priority |
|--------------------|--------------|------------|----------|
| | | | |

## 6. Functional Test Cases

#### TC-001: [Test Case Name]

| Attribute | Value |
|-----------|-------|
| **Requirement** | FR-001 |
| **Priority** | High |
| **Type** | Functional |

**Preconditions:**

**Test Steps:**

1. Step 1
2. Step 2

**Expected Results:**

**Status:** <!-- Pass / Fail / Blocked / Not Run -->

---

## 7. Non-Functional Test Cases

#### TC-NFR-001: Performance at Launch Projection

| Attribute | Value |
|-----------|-------|
| **Requirement** | TRD-xxx (Operational Guarantees) |
| **Priority** | High |
| **Type** | Performance |

**Test Configuration:**

| Parameter | Value |
|-----------|-------|
| Load level | <!-- the launch projection from §3, not current traffic --> |
| Duration | |
| Ramp-up Period | |

**Success Criteria:**

| Metric | Target | Actual |
|--------|--------|--------|
| Response Time (p50) | | |
| Response Time (p99) | | |
| Error Rate | | |
| Throughput | | |

---

#### TC-NFR-002: Public Contract Conformance

| Attribute | Value |
|-----------|-------|
| **Requirement** | TRD-xxx (Interface Contracts §6) |
| **Priority** | High |
| **Type** | Contract |

- [ ] Responses conform to the frozen published schema
- [ ] Undocumented fields absent from responses
- [ ] Error semantics match the documented table
- [ ] Rate limits behave as documented

---

## 8. Resiliency, Observability, and Runbook Testing

*Risk-based — but a one-way-door launch is rarely low-tier. State the risk
assessment either way, and remember DR that has never been executed is a
gap, not a guarantee.*

**Risk tier:** {{ high | medium | low }} — <!-- rationale; justify anything below high for a public product launch -->

| Test Type | Applicable? | Scenario | Test Case ID |
|-----------|--------------|----------|--------------|
| Dependency failure / degraded mode | | | |
| Chaos / fault injection | Yes, if high risk tier — otherwise "Not applicable — <rationale>" | | |
| GameDay (launch-day operational rehearsal) | | | |
| DR recovery procedure executed (see IRD §11) | | | |
| Observability validation (alerts actually fire, dashboards actually populate — not just configured) | | | |
| Runbook validation (the documented recovery procedure is executed, not just reviewed) | | | |

## 9. Test Data

| Data Type | Source | Sensitivity | Handling |
|-----------|--------|-------------|----------|
| | | | |

<!-- Reproducibility matters as much as realism — flaky or manual-only
     checks need an owner and disposition. Migration rehearsals need
     production-like volume and shape, not toy fixtures. -->

## 10. Test Automation

### 10.1 Automation Scope

| Test Type | Automation Coverage | Tool | Execution Layer |
|-----------|---------------------|------|-------------------|
| Unit | | | |
| Integration | | | |
| E2E | | | |
| Contract | | | |

### 10.2 CI/CD Integration

| Stage | Tests Run | Gate Criteria |
|-------|-----------|---------------|
| Commit | Unit | 100% pass |
| Build | Unit + Integration | 100% pass |
| Deploy (Staging) | System + Contract | 100% pass |
| Pre-launch | §3 proofs complete | All rows Pass |

## 11. Entry and Exit Criteria

### 11.1 Entry Criteria

- [ ] 6-pager decision meeting approved (build spend gate)
- [ ] TRD completed and approved
- [ ] Test environment available at launch-scale capacity
- [ ] Test data prepared

### 11.2 Exit Criteria (Launch Gate)

- [ ] Every row in §2 (Promise and Requirement Traceability) has a passing verification
- [ ] Every row in §3 (Pre-Launch Validation) is Pass — no irreversible commitment launches unproven
- [ ] Every row in §4 (Reversibility Containment) has an executed rollback
- [ ] No critical or high defects open
- [ ] Performance and security targets met at launch projection

## 12. Verification Sequencing

*Product-level sequencing only — which verification must complete before
which gate. Dates and staffing belong in PLAN/ROADMAP, not here.*

| Phase | Depends On | Gate |
|-------|------------|------|
| Test Planning | TRD approved | |
| Test Execution | Environment at launch scale | |
| Pre-Launch Proofs (§3) | All functional/NFR passing | Launch |
| Test Closure | Exit criteria met | |

## 13. Sign-Off

| Role | Name | Sign-Off |
|------|------|----------|
| QA | | |
| Engineering | | |
| Product | | |
| Bar raiser | | |

## 14. Known Issues and Open Questions

| Item | Type | Severity / Owner | Disposition |
|------|------|-------------------|-------------|
| | Issue / Question | | |

## Appendix

### A. Test Case Inventory

| ID | Name | Type | Priority | Status |
|----|------|------|----------|--------|

### B. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | | | Initial draft |
