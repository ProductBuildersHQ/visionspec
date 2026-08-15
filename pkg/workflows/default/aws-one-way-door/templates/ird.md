# Infrastructure Requirements Document (IRD) - Product

## Overview

**Product Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Draft

The IRD is the product-level **infrastructure contract** for a one-way-door
launch: what the infrastructure MUST guarantee, tool-agnostic (Terraform,
CloudFormation, Pulumi, or other IaC). The TRD states what the application
must guarantee; this document states what the infrastructure underneath it
must guarantee. Some infrastructure choices are themselves one-way doors —
region and data residency, account/tenancy structure, externally visible
endpoints — name them in §3 rather than discovering them after launch.

### No Infrastructure Changes (if applicable)

**Infrastructure impact:** {{ None | See below }}

<!-- "None" is implausible for a product-scale launch — if claimed, the
     rationale and evidence must be strong enough to survive review as a
     deliberate decision, not a skipped section. -->

**Rationale:**
**Evidence:** <!-- e.g., TRD capacity analysis, existing platform headroom -->

## 1. Introduction

### 1.1 Purpose

<!-- What infrastructure is being defined? What system does it support? -->

### 1.2 Scope

<!-- In scope and out of scope for this infrastructure design. -->

### 1.3 References

| Document | Link |
|----------|------|
| TRD | |
| PRD | |
| Security Policy | |

## 2. Infrastructure Requirements and Traceability

| IRD ID | Requirement | Source (TRD ID) | Verification (TPD ID) |
|--------|-------------|-------------------|--------------------------|
| IRD-001 | *e.g., "The API workload MUST have no public ingress except the API gateway."* | TRD-001 | TC-INFRA-001 |

## 3. Irreversible Infrastructure Commitments

*The infrastructure mirror of PRD §3 / TRD §3. Choices that cannot be walked
back after a public launch — typically region and data residency, account or
tenancy structure, and externally visible endpoints/domains. Each row links
to the PRD/TRD door it carries; an infrastructure one-way door with no
product-level counterpart is an escalation, not a footnote.*

| ID | Commitment | Carries (OWD/TWD ID) | Why Irreversible | Pre-Launch Proof (TPD) |
|----|------------|-----------------------|-------------------|--------------------------|
| IWD-1 | *e.g., EU data residency for customer data* | OWD-2 | Announced compliance posture; migration would breach it | |

## 4. Environment Model

| Environment | Account / Project / Tenant | Region | Purpose |
|-------------|------------------------------|--------|---------|
| Development | | | Dev/testing |
| Staging | | | Pre-production |
| Production | | | Live traffic |
| DR | | | Disaster recovery |

- **Isolation boundaries:** <!-- how environments/tenants are separated -->
- **Naming and tagging conventions:** <!-- resource naming, cost-allocation tags -->
- **Configuration ownership:** <!-- what's shared vs. environment-specific -->
- **Promotion model:** <!-- how a change moves dev → staging → prod -->

## 5. Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Infrastructure Overview                   │
├─────────────────────────────────────────────────────────────┤
│   [Load Balancer] ──► [App Servers] ──► [Database]          │
└─────────────────────────────────────────────────────────────┘
```

## 6. Launch Capacity

*Sized to the Press Release's think-big projection (TRD §9), not current
traffic — and provisioned in advance. Service quotas discovered at launch are
outages.*

### 6.1 Load Projections

| Metric | Steady State | Launch-Day Peak | Year-1 Projection |
|--------|---------------|------------------|--------------------|
| Requests/sec | | | |
| Storage | | | |
| Concurrent users | | | |

### 6.2 Quotas and Limits

| Service Quota / Limit | Current | Required for Launch Peak | Raised? |
|------------------------|---------|---------------------------|---------|
| | | | |

### 6.3 Compute and Scaling

| Component | Instance/Platform | Min | Max | Auto-scale | Pre-warmed for Launch? |
|-----------|-------------------|-----|-----|------------|--------------------------|
| | | | | Yes/No | |

## 7. Data Storage

### 7.1 Databases and Storage

| Resource | Type | Size | Replicas | Backup |
|----------|------|------|----------|--------|
| Primary DB | | | | Daily |

### 7.2 Backup and Recovery

| Data Type | Backup Frequency | Retention | RTO | RPO |
|-----------|------------------|-----------|-----|-----|
| Database | | | | |

## 8. Networking

| Component | CIDR / Config | Purpose |
|-----------|----------------|---------|
| VPC | 10.0.0.0/16 | Main network |
| Public Subnet | 10.0.1.0/24 | Load balancers |
| Private Subnet | 10.0.2.0/24 | Application |
| Data Subnet | 10.0.3.0/24 | Databases |

**Externally visible endpoints / domains:** <!-- these are public commitments
once launched — cross-reference §3 -->

## 9. Resource Lifecycle

*A valid deployment can still produce an unsafe replacement. This section is
what prevents that — for every resource category that holds state or data.*

| Resource Category | Create | Update | Replace | Delete | Protected? |
|---------------------|--------|--------|---------|--------|-------------|
| Data-bearing (DB, storage) | | | Backup before replace | Deny by default | Yes |
| Disposable (compute, cache) | | | | | No |

- **Expected replacements:** <!-- which changes are known to force a resource replacement, documented in advance -->
- **Import / adoption of existing resources:** <!-- how pre-existing resources enter management -->
- **Behavior during full stack destruction:** <!-- what survives, what doesn't -->

## 10. Security

### 10.1 Identity and Access

| Role | Permissions | Principal |
|------|-------------|-----------|
| Admin | Full access | |
| Deploy | Apply/plan | |
| Read-only | Read | |

### 10.2 Network Security

| Security Group | Inbound | Outbound | Associated |
|----------------|---------|----------|------------|
| Web | 443 from 0.0.0.0/0 | All | Load Balancer |
| Data | 5432 from App SG | None | Database |

### 10.3 Secrets and Encryption

| Data State / Item | Method | Owner / Rotation |
|-----------------------|--------|---------------------|
| At rest | AES-256 | KMS |
| In transit | TLS 1.3 | ACM |
| Secrets | Secrets Manager / equivalent | <!-- rotation policy --> |

**Prohibited:** Plaintext secrets in stack outputs, state files committed to
version control, or long-lived static credentials where a workload identity
is available.

### 10.4 Compliance

| Requirement | Implementation | Validation |
|-------------|----------------|------------|
| SOC 2 / GDPR / PCI-DSS (as applicable) | | |

## 11. High Availability and Disaster Recovery

*For a one-way-door launch, DR is proven before the public commitment is
made — "Last Tested" must be a pre-launch date, and the test must be an
executed procedure, not a tabletop review.*

| Metric | Target |
|--------|--------|
| Uptime SLA | |
| RTO | |
| RPO | |

| Scenario | Recovery Procedure | Last Tested (pre-launch) |
|----------|-----------------------|---------------------------|
| AZ / node failure | | |
| Region failure | | |
| Data corruption | Point-in-time restore | |

## 12. Observability

| Type | Destination | Retention | Alert |
|------|--------------|-----------|-------|
| Logs | | 30 days | |
| Metrics | | | > threshold |
| Traces | | | |
| Audit | | 1 year | |

## 13. Cost

*Frugality with a long horizon: cost at launch and at the year-1 projection,
so the bet's economics are visible before it becomes irrevocable.*

| Resource | Quantity | Monthly Cost (Launch) | Monthly Cost (Year-1 Projection) |
|----------|----------|------------------------|-----------------------------------|
| | | | |
| **Total** | | | |

## 14. Deployment Guarantees

*The mechanics of getting a change safely into production, tool-agnostic.*

- **Preview / plan required before production apply:** <!-- yes/no, enforced how -->
- **Policy validation:** <!-- what's checked before apply (e.g., no public buckets) -->
- **Drift detection:** <!-- how divergence between desired and actual state is caught -->
- **Deployment ordering:** <!-- dependency order across stacks/modules -->
- **Approval requirements:** <!-- who must approve a production apply -->
- **Failure and retry behavior:** <!-- what happens when an apply fails partway -->
- **Rollback vs. roll-forward:** <!-- for data-bearing resources, rollback is often not
     safe — state the actual recovery strategy; for §3 commitments there is no
     rollback, only the pre-launch proof -->

## 15. Provisioning

| Tool | Repository | Coverage |
|------|------------|----------|
| | | |

## 16. Operations

### 16.1 Runbooks

| Runbook | Scenario | Location | Executed (not just reviewed)? |
|---------|----------|----------|-------------------------------|
| | | | |

### 16.2 On-Call

| Tier | Response Time | Escalation |
|------|---------------|------------|
| L1 | | |

## Appendix

### A. Resource Inventory

| Resource ID | Type | Environment | Owner |
|-------------|------|-------------|-------|

### B. Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | | | Initial draft |
