# GitHub SpecKit

GitHub SpecKit is a spec-driven development system designed for GitHub-native workflows with sequential task execution and PR-based collaboration.

## Overview

SpecKit provides:

- **Sequential execution**: Tasks execute in defined order
- **GitHub integration**: Native PR and branch workflows
- **Memory system**: Persistent context via `.specify/memory/`
- **Constitution sync**: Governance rules propagate to agents

## When to Use SpecKit

SpecKit is ideal for:

- ✅ GitHub-centric development workflows
- ✅ Teams using PR-based code review
- ✅ Projects with clear sequential dependencies
- ✅ Open source projects with contributor guidelines
- ✅ Single-agent execution patterns

## Integration with VisionSpec

### The Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                     VISIONSPEC                                   │
│                                                                  │
│  MRD → Press → FAQ → PRD → TRD                                  │
│                  ↓                                               │
│              spec.md                                             │
│                  ↓                                               │
│         visionspec export speckit                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     SPECKIT                                      │
│                                                                  │
│  .specify/                                                       │
│  ├── spec.md      (unified specification)                       │
│  ├── plan.md      (phased implementation plan)                  │
│  ├── tasks.md     (sequential task breakdown)                   │
│  └── memory/                                                     │
│      └── constitution.md (governance rules)                      │
└─────────────────────────────────────────────────────────────────┘
```

## Export Format

When you run `visionspec export speckit`, VisionSpec creates:

```
.specify/
├── spec.md           # Unified specification from VisionSpec
├── plan.md           # Implementation plan with phases
├── tasks.md          # Task breakdown with dependencies
└── memory/
    └── constitution.md  # From CONSTITUTION.md if present
```

### spec.md

The unified specification, directly from VisionSpec reconciliation:

```markdown
# PetStore API Specification

## Overview
[Executive summary]

## Requirements
### Functional Requirements
[From PRD]

### Non-Functional Requirements
[From TRD]

## Architecture
[From TRD]

## Test Strategy
[From TPD]
```

### plan.md

Phased implementation plan:

```markdown
# Implementation Plan

## Phase 1: Foundation
- [ ] Set up project structure
- [ ] Configure build system
- [ ] Establish testing framework

## Phase 2: Core Features
- [ ] Implement Pet CRUD operations
- [ ] Add Store inventory management
- [ ] Create User authentication

## Phase 3: Integration
- [ ] Connect to database
- [ ] Add API validation
- [ ] Implement error handling

## Phase 4: Polish
- [ ] Add documentation
- [ ] Performance optimization
- [ ] Security hardening
```

### tasks.md

Granular task breakdown:

```markdown
# Tasks

## Task 1: Project Setup
**Branch**: `001-project-setup`
**Dependencies**: None
**Acceptance Criteria**:
- Go module initialized
- Directory structure created
- CI pipeline configured

## Task 2: Pet Model
**Branch**: `002-pet-model`
**Dependencies**: Task 1
**Acceptance Criteria**:
- Pet struct defined
- JSON serialization working
- Validation rules applied
```

## Complete Workflow

### Step 1: Create Specifications

```bash
# Initialize and create specs
visionspec init petstore-api --profile startup

# Quick Working Backwards flow
visionspec create mrd -p petstore-api
visionspec synthesize press -p petstore-api
visionspec synthesize prd -p petstore-api
visionspec synthesize trd -p petstore-api

# Reconcile
visionspec reconcile -p petstore-api
```

### Step 2: Export to SpecKit

```bash
visionspec export speckit -p petstore-api
```

Output:

```
⋯ Exporting to speckit...
✓ Exported to SpecKit format
  Output: .specify/
  Files:
    - spec.md
    - plan.md
    - tasks.md
    - memory/constitution.md
```

### Step 3: Execute with SpecKit

In your AI coding agent:

```
Using SpecKit, execute the plan in .specify/
Start with Task 1: Project Setup
```

The agent will:

1. Read `spec.md` for context
2. Follow `plan.md` phases
3. Execute `tasks.md` sequentially
4. Create branches per task
5. Open PRs for review

### Step 4: Review PRs

Each task creates a PR:

```
PR #1: 001-project-setup
PR #2: 002-pet-model
PR #3: 003-pet-repository
...
```

## Configuration

Configure SpecKit export in `visionspec.yaml`:

```yaml
targets:
  speckit:
    enabled: true
    output_dir: .specify

    # Branch naming
    branch_numbering: sequential  # or: semantic

    # Task granularity
    task_size: medium  # small, medium, large

    # Include memory directory
    include_memory: true
```

## Branch Numbering

### Sequential (Default)

```
001-project-setup
002-pet-model
003-pet-repository
```

### Semantic

```
feat/project-setup
feat/pet-model
feat/pet-repository
```

## Constitution Sync

If your repository has `CONSTITUTION.md`, it's synced to `.specify/memory/constitution.md`:

```bash
# Repository root
CONSTITUTION.md

# After export
.specify/memory/constitution.md  # Copy of CONSTITUTION.md
```

This ensures agents follow your governance rules.

## Best Practices

### 1. Keep Tasks Atomic

Each task should be:

- Independently reviewable
- Small enough for one PR
- Clear acceptance criteria

### 2. Use Branch Protection

Configure GitHub branch protection:

```yaml
# .github/settings.yml
branches:
  - name: main
    protection:
      required_pull_request_reviews:
        required_approving_review_count: 1
```

### 3. Include Acceptance Criteria

Tasks should have testable criteria:

```markdown
## Task 3: Pet Repository
**Acceptance Criteria**:
- [ ] CRUD operations implemented
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] Code coverage > 80%
```

## Troubleshooting

### "Tasks not generated"

Ensure your PRD has user stories:

```bash
# Check PRD content
cat docs/specs/my-project/source/prd.md
# Should contain user stories or requirements
```

### "Branch conflicts"

Reset task numbering:

```yaml
# visionspec.yaml
targets:
  speckit:
    branch_start: 100  # Start from 100
```

### "Constitution not synced"

Verify CONSTITUTION.md exists:

```bash
ls -la CONSTITUTION.md
visionspec export speckit -p my-project
ls -la .specify/memory/
```

## See Also

- [Choosing a Target](choosing-a-target.md) - Compare with other targets
- [CLI: export](../cli/export.md) - Export command reference
- [Google Design Docs Framework](../frameworks/google.md) - Good pairing for SpecKit
