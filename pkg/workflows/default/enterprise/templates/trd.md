# Technical Requirements Document (TRD)

## Overview

**Project Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Draft

This TRD is a **technical contract**: what the system MUST guarantee, not only
how it is currently designed. Use MUST / MUST NOT / SHOULD / MAY consistently.
A condition that does not apply to this change is answered with
`Not applicable — <rationale>`, never left blank — an unanswered row reads as
missing, not exempt.

## 1. Introduction

### 1.1 Purpose

<!-- What system or feature does this TRD describe, and what decision does it enable? -->

### 1.2 Scope

<!-- What is in scope and out of scope for this technical design? -->

### 1.3 References

| Document | Link |
|----------|------|
| PRD | |
| UXD | |
| MRD | |

## 2. Technical Requirements and Traceability

*The contract itself. Every applicable PRD/UXD requirement maps to one or more
TRD requirements — or is explicitly marked "Not applicable" with rationale.
Every TRD requirement maps to a TPD verification method.*

| TRD ID | Requirement | Source | Verification |
|--------|-------------|--------|---------------|
| TRD-001 | *e.g., "The API MUST reject requests without a valid bearer token."* | PRD-FR-003 | TPD-AUTH-014 |
| TRD-002 | | | |

## 3. System Invariants

*Conditions that MUST hold across every code path, including ones the rest of
this document did not anticipate. Invariants are what let an implementer (human
or AI) fill gaps safely instead of guessing.*

- *e.g., A user MUST never access an object outside an authorized tenant.*
- *e.g., Reprocessing an event MUST NOT create duplicate records.*
- *e.g., A failed migration MUST leave the previous schema readable.*

## 4. Architecture

### 4.1 High-Level Architecture

<!-- Overall system architecture. Include a diagram if helpful (ASCII or link). -->

```
┌─────────────────────────────────────────────────────────────┐
│                     System Overview                          │
├─────────────────────────────────────────────────────────────┤
│   [Component A] ──────► [Component B] ──────► [Component C] │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Component Responsibilities

| Component | Purpose | Responsibilities |
|-----------|---------|-------------------|
| | | |

### 4.3 Data Flow

<!-- How data moves through the system. Sequence diagrams for key flows if helpful. -->

## 5. State and Lifecycle

*Applies to any entity with a status, phase, or lifecycle. If nothing in this
change has meaningful state, answer "Not applicable" with rationale.*

**Applicability:** {{ "Applicable" | "Not applicable — <rationale>" }}

| Entity | States | Trigger | Invalid Transitions |
|--------|--------|---------|----------------------|
| | | | |

- **Idempotency:** <!-- Is the operation safe to retry / replay? -->
- **Concurrency:** <!-- Behavior under concurrent operations on the same entity -->
- **Retention / deletion:** <!-- Data lifecycle and deletion behavior -->
- **Terminal and recovery states:** <!-- Which states are terminal; how recovery happens -->

## 6. Interface Contracts

*Prefer a machine-readable contract (OpenAPI, JSON Schema, Protobuf) linked
here over hand-written tables where practical.*

### 6.1 Endpoints / Operations

| Operation | Method | Description |
|-----------|--------|--------------|
| | | |

### 6.2 Request / Response Contract

```
POST /api/v1/resource
Request:  { "field": "value" }
Response: { "id": "123", "status": "created" }
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| | | | |

### 6.3 Authentication and Authorization

<!-- How are requests authenticated? What authorization model? -->

### 6.4 Error Semantics

| Error Code | HTTP Status | Meaning | Retryable |
|------------|-------------|---------|-----------|
| | | | |

### 6.5 Idempotency, Pagination, Versioning

- **Idempotency:** <!-- idempotency keys, safe-retry guarantees -->
- **Pagination and ordering:** <!-- or "Not applicable — <rationale>" -->
- **Versioning and backward compatibility:** <!-- or "Not applicable — <rationale>" -->
- **Rate limits:** <!-- or "Not applicable — <rationale>" -->

### 6.6 Events Emitted and Consumed

| Event | Emitted By | Consumed By | Schema Version |
|-------|-----------|-------------|-----------------|
| | | | |

## 7. Data Design

### 7.1 Data Models

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique identifier |

### 7.2 Data Storage

| Data Type | Storage | Rationale |
|-----------|---------|-----------|
| | | |

## 8. Operational Guarantees

*Measurable, not aspirational. Every row is a number, a method, or
"Not applicable — <rationale>" — never a blank or a word like "fast."*

| Category | Requirement | Measurement | Applicability |
|----------|-------------|-------------|-----------------|
| Performance (p50 / p99) | | | |
| Throughput | | | |
| Availability (SLA, RTO, RPO) | | | |
| Reliability | | | |
| Capacity limits | | | |
| Degraded behavior | | | |
| Observability (logs / metrics / traces / alerts) | | | |
| Auditability | | | |

## 9. Security and Privacy

- [ ] Authentication mechanism
- [ ] Authorization model
- [ ] Data encryption (at rest)
- [ ] Data encryption (in transit)
- [ ] Audit logging
- [ ] Vulnerability scanning / dependency policy

<!-- For each unchecked item that does not apply, state why. -->

## 10. Compatibility and Migration Guarantees

**Applicability:** {{ "Applicable" | "Not applicable — <rationale>" }}

- **Supported versions:** <!-- what must keep working -->
- **Forward / backward compatibility:** <!-- old clients vs new server, and vice versa -->
- **Mixed-version deployment behavior:** <!-- rolling deploys, dual-read/write windows -->
- **Migration phases:** <!-- ordered steps, each independently safe -->
- **Data validation:** <!-- how migrated data is verified -->
- **Abort conditions and rollback boundary:** <!-- past which point rollback is no longer safe -->

## 11. Dependencies

| Dependency | Type | Purpose | Criticality |
|------------|------|---------|-------------|
| | External service / library / infra | | |

## 12. Testing Strategy

<!-- Unit, integration, load, and security testing approach. Reference TPD for
     the full verification plan — this section states the strategy, TPD proves it. -->

## 13. Deployment

### 13.1 Rollout and Feature Flags

| Flag | Purpose | Default |
|------|---------|---------|
| | | |

### 13.2 Rollback Plan

<!-- How to roll back if deployment fails. If rollback is unsafe past a point
     (see §10), state the roll-forward alternative. -->

## 14. Risks and Open Questions

| ID | Risk / Question | Impact | Mitigation / Resolution |
|----|------------------|--------|--------------------------|
| | | | |

*A TRD entering review should carry no open question that blocks
implementation. Unresolved questions here are non-blocking by construction.*

## Appendix

### A. Detailed Sequence Diagrams

### B. Database Schema

### C. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | | | Initial draft |
