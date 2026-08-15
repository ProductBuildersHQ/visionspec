# Infrastructure Requirements Document (IRD)

## Overview

**Project Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Draft

The IRD is the **infrastructure contract**: what the infrastructure MUST
guarantee, tool-agnostic (Terraform, CloudFormation, Pulumi, or other IaC).
The TRD states what the application must guarantee; this document states what
the infrastructure underneath it must guarantee. If this change has no
infrastructure impact, do not skip this document — complete §1–2 and the
declaration below instead.

### No Infrastructure Changes (if applicable)

**Infrastructure impact:** {{ None | See below }}

<!-- If None: state the rationale and the evidence it rests on, so "no
     infrastructure changes" is a reviewed decision, not a skipped section. -->

**Rationale:**
**Evidence:** <!-- e.g., TRD capacity analysis, existing resource headroom -->

## 1. Introduction

### 1.1 Purpose

<!-- What infrastructure is being defined? What system does it support? -->

### 1.2 Scope

<!-- What is in scope and out of scope for this infrastructure design? -->

### 1.3 References

| Document | Link |
|----------|------|
| TRD | |
| PRD | |
| Security Policy | |

## 2. Infrastructure Requirements and Traceability

| IRD ID | Requirement | Source (TRD ID) | Verification (TPD ID) |
|--------|-------------|-------------------|--------------------------|
| IRD-001 | *e.g., "The API workload MUST have no public ingress."* | TRD-xxx | TPD-INFRA-xxx |

## 3. Environment Model

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

## 4. Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Infrastructure Overview                   │
├─────────────────────────────────────────────────────────────┤
│   [Load Balancer] ──► [App Servers] ──► [Database]          │
└─────────────────────────────────────────────────────────────┘
```

## 5. Compute Resources

### 5.1 Application Servers

| Component | Instance Type | Count (Min) | Count (Max) | Auto-scale |
|-----------|---------------|-------------|-------------|------------|
| Web Server | | | | Yes/No |
| App Server | | | | Yes/No |

### 5.2 Container Orchestration / Serverless

| Attribute | Value |
|-----------|-------|
| Platform | <!-- Kubernetes, ECS, Lambda, etc. --> |
| Cluster / Function Config | |

## 6. Data Storage

### 6.1 Databases and Storage

| Resource | Type | Size | Replicas | Backup |
|----------|------|------|----------|--------|
| Primary DB | | | | Daily |

### 6.2 Backup and Recovery

| Data Type | Backup Frequency | Retention | RTO | RPO |
|-----------|------------------|-----------|-----|-----|
| Database | | | | |

## 7. Networking

| Component | CIDR / Config | Purpose |
|-----------|----------------|---------|
| VPC | 10.0.0.0/16 | Main network |
| Public Subnet | 10.0.1.0/24 | Load balancers |
| Private Subnet | 10.0.2.0/24 | Application |
| Data Subnet | 10.0.3.0/24 | Databases |

## 8. Resource Lifecycle

*A valid deployment can still produce an unsafe replacement. This section is
what prevents that — for every resource category that holds state or data.*

| Resource Category | Create | Update | Replace | Delete | Protected? |
|---------------------|--------|--------|---------|--------|-------------|
| Data-bearing (DB, storage) | | | Backup before replace | Deny by default | Yes |
| Disposable (compute, cache) | | | | | No |

- **Expected replacements:** <!-- which changes are known to force a resource replacement, documented in advance -->
- **Import / adoption of existing resources:** <!-- how pre-existing resources enter management -->
- **Behavior during full stack destruction:** <!-- what survives, what doesn't -->

## 9. Security

### 9.1 Identity and Access

| Role | Permissions | Principal |
|------|-------------|-----------|
| Admin | Full access | |
| Deploy | Apply/plan | |
| Read-only | Read | |

### 9.2 Network Security

| Security Group | Inbound | Outbound | Associated |
|----------------|---------|----------|------------|
| Web | 443 from 0.0.0.0/0 | All | Load Balancer |
| Data | 5432 from App SG | None | Database |

### 9.3 Secrets and Encryption

| Data State / Item | Method | Owner / Rotation |
|-----------------------|--------|---------------------|
| At rest | AES-256 | KMS |
| In transit | TLS 1.3 | ACM |
| Secrets | Secrets Manager / equivalent | <!-- rotation policy --> |

**Prohibited:** Plaintext secrets in stack outputs, state files committed to
version control, or long-lived static credentials where a workload identity
is available.

### 9.4 Compliance

| Requirement | Implementation | Validation |
|-------------|----------------|------------|
| SOC 2 / GDPR / PCI-DSS (as applicable) | | |

## 10. Observability

| Type | Destination | Retention | Alert |
|------|--------------|-----------|-------|
| Logs | | 30 days | |
| Metrics | | | > threshold |
| Traces | | | |
| Audit | | 1 year | |

## 11. High Availability and Disaster Recovery

| Metric | Target |
|--------|--------|
| Uptime SLA | |
| RTO | |
| RPO | |

| Scenario | Recovery Procedure | Last Tested |
|----------|-----------------------|--------------|
| AZ / node failure | | |
| Region failure | | |
| Data corruption | Point-in-time restore | |

*"Last Tested" is not optional for a Pass — a documented recovery procedure
that has never been executed is a gap, not a guarantee.*

## 12. Cost

| Resource | Quantity | Monthly Cost |
|----------|----------|---------------|
| | | |
| **Total** | | |

## 13. Deployment Guarantees

*The mechanics of getting a change safely into production, tool-agnostic.*

- **Preview / plan required before production apply:** <!-- yes/no, enforced how -->
- **Policy validation:** <!-- what's checked before apply (e.g., no public buckets) -->
- **Drift detection:** <!-- how divergence between desired and actual state is caught -->
- **Deployment ordering:** <!-- dependency order across stacks/modules -->
- **Approval requirements:** <!-- who must approve a production apply -->
- **Failure and retry behavior:** <!-- what happens when an apply fails partway -->
- **Rollback vs. roll-forward:** <!-- for data-bearing resources, rollback is often not
     safe — state the actual recovery strategy -->

## 14. Provisioning

| Tool | Repository | Coverage |
|------|------------|----------|
| | | |

## 15. Operations

### 15.1 Runbooks

| Runbook | Scenario | Location |
|---------|----------|----------|
| | | |

### 15.2 On-Call

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
