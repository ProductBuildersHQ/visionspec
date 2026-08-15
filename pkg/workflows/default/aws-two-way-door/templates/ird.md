# Infrastructure Requirements Document (IRD) - Feature

## Overview

**Feature Name:** {feature_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

This is the feature-level infrastructure contract, scoped to the incremental
change this feature requires. If this feature has no infrastructure impact,
complete the declaration below and skip the rest.

### No Infrastructure Changes (if applicable)

**Infrastructure impact:** {{ None | See below }}

**Rationale:**
**Evidence:** <!-- e.g., existing resource headroom, TRD capacity analysis -->

## 1. References

| Document | Link |
|----------|------|
| TRD | |
| PRD | |

## 2. Infrastructure Requirements and Traceability

| IRD ID | Requirement | Traces To (TRD) | Verified By (TPD) |
|--------|-------------|-------------------------|--------------------------|
| IRD-1 | | TRD-1 | TC-INFRA-1 |

## 3. Current Infrastructure

### 3.1 Existing Components

<!-- What infrastructure currently exists that this feature will use? -->

| Component | Service | Current Capacity |
|-----------|---------|------------------------|
| | | |

## 4. Infrastructure Changes

### 4.1 New / Modified Resources

| Resource | Type | Change | Purpose |
|----------|------|--------|---------|
| | | New / Modified | |

### 4.2 Resource Lifecycle Safety

*A valid deployment can still produce an unsafe replacement. State protection
for anything that holds state or data.*

| Resource | Data-Bearing? | Replace/Delete Protection |
|----------|---------------------|---------------------------------|
| | Yes/No | |

**Expected replacements:** <!-- changes known to force a replacement, documented in advance -->

## 5. Scaling

### 5.1 Expected Load

| Metric | Current | With Feature | Peak |
|--------|---------|---------------------|------|
| Requests/sec | | | |

## 6. Security

- [ ] New security groups / firewall rules
- [ ] Access control changes (IAM)
- [ ] Encryption at rest / in transit

## 7. Monitoring

| Metric | Source | Alert Threshold |
|--------|--------|-----------------------|
| | | |

## 8. Cost Impact

| Resource | Monthly Cost |
|----------|--------------------|
| | |
| **Total** | |

## 9. Deployment Guarantees

- **Preview / plan required before production apply:** <!-- yes/no -->
- **Deployment order:** <!-- sequence of infrastructure changes -->
- **Rollback vs. roll-forward:** <!-- for data-bearing resources, rollback is often unsafe — state the real strategy -->

## 10. Testing

| Test | Purpose | Automated |
|------|---------|-----------------|
| | | Yes/No |

## 11. Documentation

- [ ] Architecture diagrams
- [ ] Runbooks

## 12. Open Questions

| Question | Owner | Resolution |
|----------|-------|------------|
| | | |
