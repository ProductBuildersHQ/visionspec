# Execution: PetStore API

This walkthrough demonstrates exporting the PetStore API specification and implementing it with different AI coding agents.

## Prerequisites

Before execution, ensure you have:

1. Completed specifications (MRD, PRD, UXD, TRD, TPD, IRD)
2. Reconciled `spec.md`
3. All specs approved

```bash
# Verify project status
visionspec status -p petstore-api

# Expected output:
# Project: petstore-api
# Profile: startup
#
# Specifications:
#   ✓ MRD: approved
#   ✓ PRD: approved (synthesized)
#   ✓ UXD: approved
#   ✓ TRD: approved (synthesized)
#   ✓ TPD: approved (synthesized)
#   ✓ IRD: approved (synthesized)
#
# Reconciliation:
#   ✓ spec.md: generated
#
# Ready for export: Yes
```

## Option 1: AWS AI-DLC

Best for enterprise implementations with approval gates.

### Export

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

### Execute with Claude Code

In Claude Code:

```
Using AI-DLC, implement the PetStore API based on the vision
document in .aidlc/

Start with the Inception phase to analyze requirements and
design the architecture.
```

AI-DLC will:

1. **Workspace Detection** - Recognize this as a greenfield project
2. **Requirements Analysis** - Load from `imported-requirements.md`
3. **Application Design** - Design components based on TRD
4. **Units Generation** - Decompose into:
   - `pet-service` - Pet CRUD operations
   - `store-service` - Inventory and orders
   - `user-service` - Authentication
   - `api-gateway` - HTTP routing

5. **Construction** - For each unit:
   - Create functional design
   - Generate code and tests
   - Document decisions

6. **Build & Test** - Generate:
   - Build instructions
   - Test execution plan
   - Deployment guide

### Monitor Progress

```bash
# View AI-DLC state
cat aidlc-docs/aidlc-state.md

# View audit trail
cat aidlc-docs/audit.md

# Check generated code
ls -la internal/pet/
ls -la cmd/server/
```

## Option 2: GitHub SpecKit

Best for GitHub-native PR workflows.

### Export

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

### Execute with Claude Code

```
Using SpecKit, execute the plan in .specify/

Start with Task 1: Project Setup
```

### Review Generated Files

**.specify/plan.md:**

```markdown
# Implementation Plan

## Phase 1: Foundation
- [ ] Initialize Go module
- [ ] Set up directory structure
- [ ] Configure database connection
- [ ] Add health check endpoint

## Phase 2: Pet API
- [ ] Create Pet model
- [ ] Implement Pet repository
- [ ] Build Pet handler
- [ ] Add Pet service tests

## Phase 3: Store API
- [ ] Create Order model
- [ ] Implement inventory endpoints
- [ ] Build order processing

## Phase 4: User API
- [ ] Implement authentication
- [ ] Add user management
- [ ] Set up JWT middleware
```

**.specify/tasks.md:**

```markdown
# Tasks

## Task 1: Project Setup
**Branch**: `001-project-setup`
**Dependencies**: None
**Acceptance Criteria**:
- [ ] Go module initialized with `go.mod`
- [ ] Directory structure matches TRD
- [ ] Makefile with build targets

## Task 2: Pet Model
**Branch**: `002-pet-model`
**Dependencies**: Task 1
**Acceptance Criteria**:
- [ ] Pet struct in `internal/pet/model.go`
- [ ] JSON and DB tags
- [ ] Validation methods
- [ ] Unit tests
```

### PR Workflow

Each task creates a branch and PR:

```bash
# View branches
git branch -a

# Expected:
# 001-project-setup
# 002-pet-model
# 003-pet-repository
# ...
```

## Option 3: GSD (Get Shit Done)

Best for fast, parallel implementation.

### Export

```bash
visionspec export gsd -p petstore-api
```

Output:

```
⋯ Exporting to gsd...
✓ Exported to GSD format
  Output: gsd/
  Files:
    - PLAN.md (4 waves, 16 tasks)
    - STATE.md
    - config.json
```

### Execute with Claude Code

```
Using GSD, execute PLAN.md in gsd/

Start with Wave 1: Foundation
```

### Review Generated Files

**gsd/PLAN.md:**

```markdown
---
must_haves:
  - Pet CRUD API functional
  - OpenAPI 3.0 specification
  - Unit test coverage > 80%
  - Docker deployment ready

truths:
  - Go 1.22+
  - PostgreSQL 15+
  - Chi router
  - JWT authentication

artifacts:
  - cmd/server/main.go
  - internal/pet/handler.go
  - internal/pet/repository.go
  - api/openapi.yaml
  - Dockerfile
---

# PetStore API Implementation Plan

## Wave 1: Foundation (Parallel)
- [ ] Initialize Go module and dependencies
- [ ] Create directory structure
- [ ] Set up PostgreSQL connection
- [ ] Configure Chi router base

## Wave 2: Pet API (Parallel)
- [ ] Implement Pet model and validation
- [ ] Create Pet repository (CRUD)
- [ ] Build Pet HTTP handler
- [ ] Add Pet service tests

## Wave 3: Store & User API (Parallel)
- [ ] Implement inventory endpoints
- [ ] Build order processing
- [ ] Add user authentication
- [ ] Implement JWT middleware

## Wave 4: Quality & Deploy (Parallel)
- [ ] Generate OpenAPI specification
- [ ] Add integration tests
- [ ] Create Dockerfile
- [ ] Write deployment guide
```

### Track Progress

```bash
# Watch state file
watch cat gsd/STATE.md

# Verify artifacts
ls -la cmd/server/ internal/pet/ api/
```

## Option 4: GasTown

Best for complex multi-agent orchestration.

### Export

```bash
visionspec export gastown -p petstore-api
```

### Review Generated Files

**gastown/formula.toml:**

```toml
[formula]
name = "petstore-api"
type = "convoy"
rig = "go-backend"

[formula.beads]
order = [
  "foundation",
  "pet-model",
  "pet-repository",
  "pet-handler",
  "store-inventory",
  "store-order",
  "user-auth",
  "api-docs"
]
```

### Execute with Multiple Agents

GasTown supports running multiple agents in parallel:

```bash
# Agent 1: Foundation + Pet
gastown run --bead foundation --bead pet-model

# Agent 2: Store
gastown run --bead store-inventory --bead store-order

# Agent 3: User + Docs
gastown run --bead user-auth --bead api-docs
```

## Option 5: GasCity

Best for role-based agent coordination.

### Export

```bash
visionspec export gascity -p petstore-api
```

### Review Generated Files

**gascity/city.toml:**

```toml
[city]
name = "petstore-api"
mode = "orchestrated"

[[agents]]
role = "backend"
capabilities = ["go", "api", "database"]

[[agents]]
role = "devops"
capabilities = ["docker", "ci_cd"]

[[agents]]
role = "qa"
capabilities = ["testing"]

[[orders]]
id = "pet-api"
agent = "backend"
priority = 1
dependencies = []

[[orders]]
id = "store-api"
agent = "backend"
priority = 2
dependencies = ["pet-api"]

[[orders]]
id = "deployment"
agent = "devops"
priority = 3
dependencies = ["pet-api", "store-api"]

[[orders]]
id = "testing"
agent = "qa"
priority = 4
dependencies = ["deployment"]
```

### Execute with Role Assignment

```
Using GasCity, coordinate the petstore-api city.

I am the backend agent. Assign my orders.
```

## Comparing Execution Results

After implementation, all targets should produce:

```
petstore-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── pet/
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── store/
│   │   ├── handler.go
│   │   ├── inventory.go
│   │   └── order.go
│   └── user/
│       ├── handler.go
│       └── auth.go
├── api/
│   └── openapi.yaml
├── migrations/
│   └── 001_initial.sql
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Verification

After execution, verify the implementation:

```bash
# Run tests
make test

# Check coverage
make coverage

# Build
make build

# Run locally
make run

# Test API
curl http://localhost:8080/v1/pets
```

## Iteration

If changes are needed:

1. Update VisionSpec specs
2. Re-reconcile: `visionspec reconcile -p petstore-api`
3. Re-export: `visionspec export <target> -p petstore-api`
4. Continue execution with context preserved

## Summary

| Target | Best For | Execution Style |
|--------|----------|-----------------|
| AI-DLC | Enterprise | Phased with approvals |
| SpecKit | GitHub | PR-based sequential |
| GSD | Speed | Wave-parallel |
| GasTown | Complex | DAG multi-agent |
| GasCity | Roles | Agent-coordinated |

All targets consume the same VisionSpec artifacts (`spec.md`) and produce working software. Choose based on your team structure and workflow preferences.
