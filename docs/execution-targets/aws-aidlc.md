# AWS AI-DLC

AWS AI-DLC (AI-Driven Development Lifecycle) is an adaptive software development methodology that guides AI coding agents through structured implementation phases.

## Overview

AI-DLC provides:

- **Three-phase lifecycle**: Inception → Construction → Operations
- **Adaptive execution**: Stages execute based on project complexity
- **Human-in-the-loop**: Approval gates at critical decision points
- **Complete audit trail**: Full traceability from spec to code
- **Multi-agent support**: Specialized agents for different phases

## When to Use AI-DLC

AI-DLC is ideal for:

- ✅ Enterprise software development
- ✅ Projects requiring audit trails and compliance
- ✅ Complex multi-phase implementations
- ✅ Teams using Claude Code, Amazon Q Developer, or Kiro
- ✅ Projects with multiple stakeholders requiring approval gates

## Integration with VisionSpec

### The Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                     VISIONSPEC (Specification)                   │
│                                                                  │
│  MRD → Press Release → FAQ → PRD → UXD                          │
│                  ↓                                               │
│  TRD → TPD → IRD (grounded in codebase context)                 │
│                  ↓                                               │
│  Conflict Detection → Reconciliation → spec.md                  │
│                  ↓                                               │
│         visionspec export aidlc                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     AWS AI-DLC (Execution)                       │
│                                                                  │
│  🔵 INCEPTION                                                    │
│     Workspace Detection → Requirements → Application Design      │
│                  ↓                                               │
│  🟢 CONSTRUCTION                                                 │
│     Per-Unit: Functional Design → NFR → Code Generation          │
│                  ↓                                               │
│  Build & Test Instructions                                       │
│                  ↓                                               │
│  🟡 OPERATIONS (future)                                         │
└─────────────────────────────────────────────────────────────────┘
```

### What VisionSpec Provides

VisionSpec creates the strategic foundation:

| VisionSpec Artifact | AI-DLC Usage |
|---------------------|--------------|
| MRD (Market Requirements) | Vision document context |
| Press Release | Customer-facing vision |
| FAQ | Scope clarification, edge cases |
| PRD (Product Requirements) | Imported requirements |
| UXD (User Experience) | Design constraints |
| TRD (Technical Requirements) | Technical environment |
| IRD (Infrastructure) | Infrastructure context |
| spec.md | Unified execution specification |

### What AI-DLC Produces

AI-DLC handles implementation:

| AI-DLC Artifact | Purpose |
|-----------------|---------|
| `aidlc-state.md` | Central state tracking |
| `audit.md` | Complete audit trail |
| Architecture docs | System design |
| Functional designs | Per-unit business logic |
| Application code | Generated implementation |
| Test suites | Unit, integration, E2E tests |
| Build instructions | CI/CD guidance |

## Export Format

When you run `visionspec export aidlc`, VisionSpec creates:

```
.aidlc/
├── vision-document.md          # From MRD + Press Release + PRD
├── technical-environment.md    # From TRD + IRD + context
└── imported-requirements.md    # From spec.md requirements section
```

### vision-document.md

Combines your Working Backwards artifacts into a single vision document:

```markdown
# Vision Document

## Executive Summary
[From Press Release]

## Problem Statement
[From MRD market problem]

## Customer Benefits
[From Press Release + FAQ]

## Success Metrics
[From PRD success criteria]

## Scope
[From FAQ scope clarification]
```

### technical-environment.md

Provides technical context for implementation:

```markdown
# Technical Environment

## Architecture Overview
[From TRD architecture section]

## Technology Stack
[From TRD + IRD technology choices]

## Infrastructure
[From IRD deployment requirements]

## Constraints
[From TRD constraints]
```

### imported-requirements.md

Contains the actionable requirements:

```markdown
# Imported Requirements

## Functional Requirements
[From spec.md functional section]

## Non-Functional Requirements
[From spec.md NFR section]

## Acceptance Criteria
[From PRD user stories]
```

## Complete Workflow

### Step 1: Create Specifications with VisionSpec

```bash
# Initialize project with enterprise profile
visionspec init petstore-api --profile enterprise

# Author source specs
visionspec create mrd -p petstore-api
# Edit docs/specs/petstore-api/source/mrd.md with market context

visionspec create prd -p petstore-api
# Edit docs/specs/petstore-api/source/prd.md with requirements

visionspec create uxd -p petstore-api
# Edit docs/specs/petstore-api/source/uxd.md with user experience

# Synthesize Working Backwards artifacts
visionspec synthesize press -p petstore-api
visionspec synthesize faq -p petstore-api

# Synthesize technical specs
visionspec synthesize trd -p petstore-api
visionspec synthesize tpd -p petstore-api
visionspec synthesize ird -p petstore-api

# Evaluate all specs
visionspec eval all -p petstore-api

# Approve specs
visionspec approve mrd -p petstore-api --approver "product@example.com"
visionspec approve prd -p petstore-api --approver "product@example.com"
visionspec approve trd -p petstore-api --approver "tech@example.com"

# Reconcile into unified spec
visionspec reconcile -p petstore-api
```

### Step 2: Export to AI-DLC

```bash
visionspec export aidlc -p petstore-api
```

Output:

```
⋯ Exporting to aidlc...
✓ Exported to AWS AI-DLC format
  Output: .aidlc/
  Files:
    - vision-document.md
    - technical-environment.md
    - imported-requirements.md
```

### Step 3: Execute with AI-DLC

In your AI coding agent (Claude Code, Amazon Q, Kiro):

```
Using AI-DLC, implement the PetStore API based on the vision document in .aidlc/
```

AI-DLC will:

1. **Workspace Detection** - Analyze existing codebase (if any)
2. **Requirements Analysis** - Load imported requirements, ask clarifying questions
3. **Application Design** - Design components and services
4. **Units Generation** - Decompose into implementable units
5. **Construction** - For each unit:
   - Functional design
   - NFR design
   - Code generation
   - Test generation
6. **Build & Test** - Create build and test instructions

### Step 4: Review and Iterate

AI-DLC maintains human-in-the-loop at critical points:

- Requirements approval before design
- Design approval before code generation
- Code review before completion

All decisions are logged in `aidlc-docs/audit.md`.

## Configuration

Configure AI-DLC export in `visionspec.yaml`:

```yaml
targets:
  aidlc:
    enabled: true
    output_dir: .aidlc

    # Include additional context
    include_context: true

    # Map VisionSpec specs to AI-DLC documents
    vision_sources:
      - mrd
      - press
      - prd
    technical_sources:
      - trd
      - ird
```

## AI-DLC Phases Explained

### 🔵 INCEPTION Phase

**Purpose**: Determine WHAT to build and WHY

| Stage | Condition | VisionSpec Source |
|-------|-----------|-------------------|
| Workspace Detection | Always | N/A (analyzes codebase) |
| Reverse Engineering | Brownfield only | Existing code |
| Requirements Analysis | Always | `imported-requirements.md` |
| User Stories | User-facing features | PRD user stories |
| Application Design | New components | TRD architecture |
| Units Generation | Complex systems | TRD components |
| Workflow Planning | Always | All specs |

### 🟢 CONSTRUCTION Phase

**Purpose**: Determine HOW to build it

For each unit of work:

| Stage | Condition | Output |
|-------|-----------|--------|
| Functional Design | New business logic | Detailed design |
| NFR Requirements | Performance/security needs | NFR spec |
| NFR Design | NFR patterns needed | Implementation approach |
| Infrastructure Design | Infra changes | Service mapping |
| Code Generation | Always | Application code + tests |

After all units:

| Stage | Output |
|-------|--------|
| Build and Test | Build instructions, test suites |

### 🟡 OPERATIONS Phase

**Purpose**: Deploy and run it (future expansion)

Currently a placeholder for:

- Deployment automation
- Monitoring setup
- Incident response

## Best Practices

### 1. Complete Your Specs First

AI-DLC works best with comprehensive specifications:

```bash
# Ensure all required specs are approved
visionspec status -p my-project

# Status should show:
# ✓ MRD: approved
# ✓ PRD: approved
# ✓ TRD: approved
# ✓ spec.md: generated
```

### 2. Include Context Sources

Ground your specs in reality:

```yaml
# visionspec.yaml
context:
  git:
    - path: .
      include_patterns:
        - "**/*.go"
        - "**/*.ts"
  files:
    - path: docs/architecture.md
```

### 3. Use Appropriate Profile

Match VisionSpec profile to project stage:

| Stage | VisionSpec Profile | AI-DLC Depth |
|-------|-------------------|--------------|
| Prototype | `startup` | Minimal |
| MVP | `growth` | Standard |
| Enterprise | `enterprise` | Comprehensive |

### 4. Review AI-DLC Artifacts

After AI-DLC execution, review:

- `aidlc-docs/audit.md` - Decision trail
- `aidlc-docs/aidlc-state.md` - Execution state
- Generated code - Implementation quality

## Troubleshooting

### "Missing vision document"

Ensure export completed:

```bash
visionspec export aidlc -p my-project
ls -la .aidlc/
```

### "Requirements not found"

Verify reconciliation:

```bash
visionspec status -p my-project
# Check that spec.md exists and has requirements
```

### "Context mismatch"

Re-export with fresh context:

```bash
visionspec context gather -p my-project
visionspec reconcile -p my-project
visionspec export aidlc -p my-project
```

## See Also

- [Choosing a Target](choosing-a-target.md) - Compare AI-DLC with other targets
- [AWS Working Backwards Framework](../frameworks/aws.md) - Best paired with AI-DLC
- [CLI: export](../cli/export.md) - Export command reference
