# AIDLC Demo Project

This sample project demonstrates the integration between VisionSpec and AWS AI-DLC Workflows.

## Overview

The project showcases:
- Bidirectional sync between `.visionspec/` and `aidlc-docs/` directories
- AIDLC three-phase workflow (Inception, Construction, Operations)
- Document evaluation with LLM-as-judge rubrics
- Phase transition validation

## Directory Structure

```
aidlc-demo/
├── README.md                    # This file
├── visionspec.yaml              # Project configuration
├── .visionspec/                 # VisionSpec documents (JSON)
│   └── specs.json               # Synced spec data
├── aidlc-docs/                  # AIDLC documents (Markdown)
│   ├── inception/               # Phase 1: Discovery & Planning
│   │   ├── vision_document.md
│   │   ├── requirements_spec.md
│   │   ├── technical_spec.md
│   │   └── architecture_spec.md
│   ├── construction/            # Phase 2: Implementation & Testing
│   │   ├── implementation_plan.md
│   │   ├── test_plan.md
│   │   ├── integration_plan.md
│   │   └── security_review.md
│   └── operations/              # Phase 3: Deployment & Ops
│       ├── runbook.md
│       ├── monitoring_plan.md
│       ├── disaster_recovery_plan.md
│       └── slo_document.md
└── aidlc-state.md               # Workflow state tracking
```

## Getting Started

### 1. View Workflow Status

```bash
visionspec status --project examples/aidlc-demo
```

### 2. Evaluate Documents

```bash
visionspec eval vision_document --project examples/aidlc-demo
```

### 3. Check Phase Transition Readiness

```bash
visionspec aidlc phase-status --project examples/aidlc-demo
```

### 4. Sync Changes

```bash
# Export VisionSpec to AIDLC format
visionspec aidlc sync --to-aidlc --project examples/aidlc-demo

# Import AIDLC to VisionSpec format
visionspec aidlc sync --from-aidlc --project examples/aidlc-demo
```

## Document Types by Phase

### Inception Phase
| Document | Purpose | Required |
|----------|---------|----------|
| Vision Document | High-level product vision and goals | Yes |
| Requirements Spec | Functional and non-functional requirements | Yes |
| Technical Spec | Detailed technical design | Yes |
| Architecture Spec | System architecture and ADRs | No |

### Construction Phase
| Document | Purpose | Required |
|----------|---------|----------|
| Implementation Plan | Task breakdown and timeline | Yes |
| Test Plan | Testing strategy and cases | Yes |
| Integration Plan | External system integrations | No |
| Security Review | Threat model and controls | Yes |

### Operations Phase
| Document | Purpose | Required |
|----------|---------|----------|
| Runbook | Operational procedures | Yes |
| Monitoring Plan | Observability strategy | Yes |
| Disaster Recovery | Business continuity | No |
| SLO Document | Service level objectives | Yes |

## Quality Evaluation

Documents are evaluated against rubrics with the following rating scale:
- **EXCELLENT** (0.9-1.0): Production ready
- **GOOD** (0.7-0.9): Minor improvements needed
- **NEEDS_IMPROVEMENT** (0.5-0.7): Significant gaps
- **POOR** (0.0-0.5): Major revision required

## Phase Transition Rules

To transition to the next phase:
1. All required documents must be completed
2. Minimum quality score of 0.7 required
3. No critical issues in evaluations
4. Explicit approval may be required

## Sample Scenario

This demo simulates a "Task Management API" project:
- **Inception**: Vision and requirements defined
- **Construction**: In progress (implementation plan started)
- **Operations**: Not started

Use this project to test the full AIDLC workflow integration.
