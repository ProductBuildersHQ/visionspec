# AIDLC Workflow Integration Guide

This guide covers integrating AWS AI DLC (AI Development Lifecycle) workflows with VisionSpec and VisionStudio.

## Overview

AIDLC provides a structured approach to AI/ML development with three main phases:

1. **Inception** - Requirements gathering, vision documents, technical specs
2. **Construction** - Implementation plans, test plans, architecture
3. **Operations** - Runbooks, monitoring, maintenance

VisionSpec provides Go types and sync utilities for working with AIDLC workflows, while VisionStudio offers visual management and LLM-based evaluation.

## Getting Started

### Directory Structure

AIDLC projects use two synchronized directories:

```
your-project/
├── .visionspec/           # VisionSpec managed specs
│   └── specs.json         # Spec metadata and sync state
├── aidlc-docs/            # AIDLC workflow documents
│   ├── inception/
│   │   ├── vision_document.md
│   │   ├── requirements_spec.md
│   │   └── technical_spec.md
│   ├── construction/
│   │   ├── implementation_plan.md
│   │   └── test_plan.md
│   ├── operations/
│   │   ├── runbook.md
│   │   └── monitoring_plan.md
│   └── aidlc-state.md     # Workflow state tracking
└── visionspec.yaml        # Project configuration
```

### Configuration

Add AIDLC to your `visionspec.yaml`:

```yaml
name: my-project
profile: aidlc

aidlc:
  enabled: true
  sync:
    direction: bidirectional  # bidirectional, to_aidlc, from_aidlc
    auto_sync: true
```

## Document Types

### Inception Phase

| Document Type | Description |
|--------------|-------------|
| `vision_document` | Project vision, goals, success criteria |
| `requirements_spec` | Functional and non-functional requirements |
| `technical_spec` | Technical architecture and design |
| `architecture_spec` | System architecture decisions |
| `api_contract` | API definitions and contracts |
| `data_model` | Data models and schemas |

### Construction Phase

| Document Type | Description |
|--------------|-------------|
| `implementation_plan` | Implementation tasks and timeline |
| `test_plan` | Testing strategy and test cases |
| `security_review` | Security analysis and mitigations |
| `performance_plan` | Performance requirements and benchmarks |
| `deployment_plan` | Deployment procedures |

### Operations Phase

| Document Type | Description |
|--------------|-------------|
| `runbook` | Operational procedures |
| `monitoring_plan` | Metrics, alerts, SLOs |
| `incident_response` | Incident handling procedures |
| `maintenance_guide` | Maintenance and updates |

## Using the Go Package

### Types

```go
import "github.com/ProductBuildersHQ/visionspec/pkg/aidlc"

// Check document type
docType := aidlc.DocVisionDocument
fmt.Println(docType.Phase())       // inception
fmt.Println(docType.DisplayName()) // Vision Document

// Quality ratings
rating := aidlc.RatingGood
fmt.Println(rating.Score()) // 0.75
```

### Workflow Management

```go
// Create default workflow
workflow := aidlc.DefaultWorkflow()

// Check current phase
phase := workflow.CurrentPhase()

// Get phase requirements
reqs := workflow.GetPhaseRequirements(aidlc.PhaseConstruction)
fmt.Printf("Progress: %.0f%%\n", reqs.ProgressPercent)
fmt.Printf("Can advance: %v\n", reqs.CanAdvance)
```

### Sync Engine

```go
engine := aidlc.NewSyncEngine(
    ".visionspec",      // VisionSpec directory
    "aidlc-docs",       // AIDLC docs directory
)

// Check for differences
ctx := context.Background()
diff, err := engine.DiffState(ctx)

// Sync changes
result, err := engine.Sync(ctx)
```

### Phase Transitions

```go
workflow := aidlc.DefaultWorkflow()
rules := aidlc.DefaultTransitionRules()

// Check if transition is possible
result := workflow.CanTransitionTo(aidlc.PhaseConstruction, rules)
if !result.Success {
    fmt.Println("Blocked by:", result.BlockingIssues)
}

// Perform transition
result, err := workflow.TransitionTo(aidlc.PhaseConstruction, rules)
```

### Document Templates

```go
// Get template
template, ok := aidlc.GetTemplate(aidlc.DocVisionDocument)

// Render with data
data := aidlc.TemplateData{
    ProjectName: "My Project",
    Title:       "Vision Document",
    Author:      "Team Lead",
}
content, err := aidlc.RenderTemplate(aidlc.DocVisionDocument, data)
```

## API Reference

VisionStudio exposes REST endpoints for AIDLC workflow management:

### State and Workflow

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects/{project}/aidlc/state` | GET | Get workflow state |
| `/api/projects/{project}/aidlc/workflow` | GET | Get workflow DAG |

### Documents

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects/{project}/aidlc/documents` | GET | List documents |
| `/api/projects/{project}/aidlc/documents/{docId}` | GET | Get document |
| `/api/projects/{project}/aidlc/documents/create` | POST | Create from template |

### Sync

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects/{project}/aidlc/sync/diff` | GET | Get sync diff |
| `/api/projects/{project}/aidlc/sync` | POST | Trigger sync |

### Phase Management

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects/{project}/aidlc/phase/requirements` | GET | Get phase requirements |
| `/api/projects/{project}/aidlc/phase/transition` | POST | Transition phase |

### Templates

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects/{project}/aidlc/templates` | GET | List templates |
| `/api/projects/{project}/aidlc/templates/{docType}` | GET | Get template |

## Quality Evaluation

AIDLC documents are evaluated using LLM-based quality scoring:

### Ratings

| Rating | Score Range | Description |
|--------|-------------|-------------|
| EXCELLENT | 0.90+ | Exceeds requirements |
| GOOD | 0.70-0.89 | Meets requirements |
| NEEDS_IMPROVEMENT | 0.50-0.69 | Partial compliance |
| POOR | Below 0.50 | Significant issues |

### Evaluation Dimensions

Documents are scored across multiple dimensions:

- **Clarity** - Writing quality and understandability
- **Completeness** - Coverage of required sections
- **Consistency** - Alignment with other documents
- **Feasibility** - Technical and business viability
- **Testability** - Can requirements be verified?

### Issue Severities

| Severity | Weight | Action Required |
|----------|--------|-----------------|
| Critical | 1.0 | Must fix before proceeding |
| High | 0.8 | Should fix before proceeding |
| Medium | 0.5 | Should address in next iteration |
| Low | 0.3 | Nice to have |
| Info | 0.1 | Informational only |

## Frontend Components

VisionStudio provides React components for AIDLC visualization:

```tsx
import {
  AIDLCWorkflowView,
  AIDLCDocumentView,
  AIDLCSyncPanel,
  PhaseRequirementsPanel,
  TransitionButton,
  TemplateSelector,
  EvaluationResultsPanel,
} from './components/aidlc'

// Workflow visualization
<AIDLCWorkflowView
  projectName="my-project"
  onNodeClick={(nodeId) => console.log(nodeId)}
/>

// Phase requirements with progress
<PhaseRequirementsPanel
  projectName="my-project"
  onDocumentClick={(docType) => console.log(docType)}
/>

// Phase transition controls
<TransitionButton
  projectName="my-project"
  currentPhase="inception"
  targetPhase="construction"
  canAdvance={true}
  onTransitionComplete={(result) => console.log(result)}
/>

// Template selection modal
<TemplateSelector
  projectName="my-project"
  isOpen={true}
  onClose={() => {}}
  onDocumentCreated={(doc) => console.log(doc)}
/>

// Evaluation results display
<EvaluationResultsPanel score={qualityScore} compact={false} />
```

## Best Practices

1. **Complete phases sequentially** - Don't skip ahead; each phase builds on the previous
2. **Maintain document quality** - Address evaluation issues before advancing phases
3. **Use bidirectional sync** - Keep VisionSpec and AIDLC directories in sync
4. **Track state** - Use `aidlc-state.md` to monitor workflow progress
5. **Review transitions** - Document reasons for phase transitions in the transition log

## Troubleshooting

### Sync Conflicts

If you encounter sync conflicts:

1. Check `aidlc-state.md` for current state
2. Run diff to see pending changes: `GET /api/projects/{project}/aidlc/sync/diff`
3. Resolve conflicts manually or use force sync

### Phase Transition Blocked

If phase transition is blocked:

1. Check blocking documents: view transition result's `blocking_docs`
2. Complete required documents
3. Ensure minimum quality scores are met
4. Review blocking issues for guidance

### Missing Templates

If templates aren't loading:

1. Verify daemon is running
2. Check API connection
3. Ensure visionspec package is up to date
