# GasCity

GasCity is a role-based multi-agent orchestration system using cities, agents, and orders for coordinated execution.

## Overview

GasCity provides:

- **Role-based agents**: Specialized agents with defined capabilities
- **Order system**: Tasks assigned to specific agent roles
- **Orchestration modes**: Orchestrated, autonomous, or hybrid execution
- **City abstraction**: Container for agents and their work

## When to Use GasCity

GasCity is ideal for:

- ✅ Complex projects requiring specialized roles
- ✅ Teams with distinct frontend/backend/DevOps agents
- ✅ Projects needing coordination between specialists
- ✅ Large-scale systems with role-based decomposition
- ✅ Organizations with defined engineering roles

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
│         visionspec export gascity                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     GASCITY                                      │
│                                                                  │
│  city.toml                                                       │
│  ├── [city] metadata                                             │
│  ├── [[agents]] role definitions                                │
│  └── [[orders]] task assignments                                │
└─────────────────────────────────────────────────────────────────┘
```

## Core Concepts

### Cities

A city contains agents and their orders:

```toml
[city]
name = "petstore-api"
mode = "orchestrated"
```

### Agents

Agents are specialized roles:

```toml
[[agents]]
role = "backend"
capabilities = ["go", "api", "database"]
```

### Orders

Orders are tasks assigned to agents:

```toml
[[orders]]
id = "pet-api"
agent = "backend"
description = "Implement Pet API"
```

## Export Format

When you run `visionspec export gascity`, VisionSpec creates:

```
gascity/
└── city.toml
```

### city.toml

Complete city configuration:

```toml
[city]
name = "petstore-api"
version = "1.0.0"
mode = "orchestrated"  # orchestrated, autonomous, hybrid

[city.context]
spec = "docs/specs/petstore-api/spec.md"
mrd = "docs/specs/petstore-api/source/mrd.md"
trd = "docs/specs/petstore-api/technical/trd.md"

# Agent Definitions
[[agents]]
role = "architect"
capabilities = ["system_design", "api_design", "documentation"]
model = "claude-3-opus"

[[agents]]
role = "backend"
capabilities = ["go", "api", "database", "testing"]
model = "claude-3-sonnet"

[[agents]]
role = "frontend"
capabilities = ["react", "typescript", "css", "testing"]
model = "claude-3-sonnet"

[[agents]]
role = "devops"
capabilities = ["docker", "kubernetes", "ci_cd", "monitoring"]
model = "claude-3-haiku"

[[agents]]
role = "qa"
capabilities = ["testing", "automation", "security"]
model = "claude-3-sonnet"

# Orders (Tasks)
[[orders]]
id = "architecture"
agent = "architect"
priority = 1
description = "Design system architecture and API contracts"
dependencies = []
artifacts = [
  "docs/architecture.md",
  "api/openapi.yaml"
]

[[orders]]
id = "pet-api"
agent = "backend"
priority = 2
description = "Implement Pet CRUD API endpoints"
dependencies = ["architecture"]
artifacts = [
  "internal/pet/handler.go",
  "internal/pet/repository.go",
  "internal/pet/service.go"
]

[[orders]]
id = "pet-ui"
agent = "frontend"
priority = 2
description = "Build Pet management UI components"
dependencies = ["architecture"]
artifacts = [
  "web/src/components/PetList.tsx",
  "web/src/components/PetForm.tsx"
]

[[orders]]
id = "pet-tests"
agent = "qa"
priority = 3
description = "Create comprehensive test suite for Pet features"
dependencies = ["pet-api", "pet-ui"]
artifacts = [
  "internal/pet/handler_test.go",
  "web/src/components/PetList.test.tsx"
]

[[orders]]
id = "deployment"
agent = "devops"
priority = 4
description = "Configure CI/CD and deployment"
dependencies = ["pet-api", "pet-ui"]
artifacts = [
  "Dockerfile",
  ".github/workflows/ci.yml",
  "k8s/deployment.yaml"
]
```

## Complete Workflow

### Step 1: Create Specifications

```bash
visionspec init petstore-api --profile enterprise

# Create specs for role-based decomposition
visionspec create mrd -p petstore-api
visionspec create prd -p petstore-api
visionspec create uxd -p petstore-api
visionspec synthesize trd -p petstore-api
visionspec synthesize ird -p petstore-api

visionspec reconcile -p petstore-api
```

### Step 2: Export to GasCity

```bash
visionspec export gascity -p petstore-api
```

Output:

```
⋯ Exporting to gascity...
✓ Exported to GasCity format
  Output: gascity/
  Files:
    - city.toml (5 agents, 5 orders)
```

### Step 3: Execute with GasCity

Using GasCity CLI:

```bash
# Initialize the city
gascity init gascity/city.toml

# Run all orders
gascity run --city gascity/city.toml

# Run specific agent's orders
gascity run --agent backend

# Run specific order
gascity run --order pet-api
```

With AI agents:

```
Using GasCity, coordinate the city in gascity/city.toml
Assign orders to agents based on their roles
```

### Step 4: Monitor Progress

```bash
# View city status
gascity status --city gascity/city.toml

# View agent status
gascity status --agent backend

# View order status
gascity status --order pet-api
```

## Configuration

Configure GasCity export in `visionspec.yaml`:

```yaml
targets:
  gascity:
    enabled: true
    output_dir: gascity

    # Orchestration mode
    mode: orchestrated  # orchestrated, autonomous, hybrid

    # Agent configuration
    agents:
      - role: backend
        capabilities: [go, api, database]
        model: claude-3-sonnet
      - role: frontend
        capabilities: [react, typescript]
        model: claude-3-sonnet
      - role: devops
        capabilities: [docker, kubernetes]
        model: claude-3-haiku

    # Order generation
    order_granularity: medium
    auto_dependencies: true
```

## Orchestration Modes

### Orchestrated

Central coordinator manages all agents:

```toml
[city]
mode = "orchestrated"

# Coordinator assigns orders
# Agents wait for instructions
# Strict dependency enforcement
```

**Best for**: Critical projects, complex dependencies, strict ordering

### Autonomous

Agents self-coordinate:

```toml
[city]
mode = "autonomous"

# Agents claim available orders
# Self-manage dependencies
# Maximum parallelism
```

**Best for**: Independent workstreams, experienced teams

### Hybrid

Mix of coordinated and autonomous:

```toml
[city]
mode = "hybrid"

# Critical paths orchestrated
# Independent work autonomous
# Balanced control
```

**Best for**: Large projects with mixed complexity

## Agent Roles

### Common Role Patterns

| Role | Capabilities | Typical Orders |
|------|--------------|----------------|
| **architect** | system_design, api_design | Architecture, contracts |
| **backend** | go, python, api, database | API implementation |
| **frontend** | react, typescript, css | UI components |
| **devops** | docker, k8s, ci_cd | Infrastructure |
| **qa** | testing, automation | Test suites |
| **security** | auth, encryption | Security review |
| **docs** | documentation, api_docs | Documentation |

### Custom Roles

Define project-specific roles:

```toml
[[agents]]
role = "ml_engineer"
capabilities = ["python", "pytorch", "mlops"]
model = "claude-3-opus"

[[agents]]
role = "data_engineer"
capabilities = ["sql", "spark", "airflow"]
model = "claude-3-sonnet"
```

## Order Dependencies

### Dependency Graph

Orders form a dependency DAG:

```
┌────────────┐
│architecture│
└─────┬──────┘
      │
 ┌────┴────┐
 ↓         ↓
┌───────┐ ┌───────┐
│pet-api│ │pet-ui │
└───┬───┘ └───┬───┘
    │         │
    └────┬────┘
         ↓
   ┌──────────┐
   │pet-tests │
   └────┬─────┘
        │
        ↓
  ┌──────────┐
  │deployment│
  └──────────┘
```

### Dependency Rules

```toml
[[orders]]
id = "pet-tests"
dependencies = ["pet-api", "pet-ui"]  # Both must complete

[[orders]]
id = "deployment"
dependencies = ["pet-tests"]  # Must pass tests first
```

## Best Practices

### 1. Match Agents to TRD Components

Align agents with your TRD architecture:

```
TRD Components          →  GasCity Agents
├── API Layer           →  backend
├── UI Layer            →  frontend
├── Data Layer          →  backend, data_engineer
├── Infrastructure      →  devops
└── Testing             →  qa
```

### 2. Use Appropriate Models

Match model capability to task complexity:

```toml
[[agents]]
role = "architect"
model = "claude-3-opus"  # Complex reasoning

[[agents]]
role = "devops"
model = "claude-3-haiku"  # Routine tasks
```

### 3. Define Clear Boundaries

Agents should have non-overlapping responsibilities:

```toml
# ✅ Good: Clear boundaries
[[agents]]
role = "backend"
capabilities = ["go", "api"]

[[agents]]
role = "frontend"
capabilities = ["react", "typescript"]

# ❌ Avoid: Overlapping
[[agents]]
role = "fullstack"
capabilities = ["go", "api", "react", "typescript"]
```

### 4. Order Granularity

Right-size your orders:

```toml
# ✅ Good: Focused orders
[[orders]]
id = "pet-handler"
agent = "backend"
artifacts = ["internal/pet/handler.go"]

# ❌ Avoid: Mega-orders
[[orders]]
id = "entire-backend"
artifacts = ["internal/**/*.go"]
```

## Troubleshooting

### "Agent not found"

Verify agent is defined:

```bash
grep -A3 "role = \"backend\"" gascity/city.toml
```

### "Order stuck"

Check dependencies:

```bash
# View order dependencies
grep -A5 "id = \"pet-tests\"" gascity/city.toml

# Check if dependencies completed
gascity status --order pet-api
gascity status --order pet-ui
```

### "Capability mismatch"

Ensure agent has required capabilities:

```toml
[[orders]]
id = "kubernetes-deploy"
agent = "devops"
# Requires: kubernetes capability

[[agents]]
role = "devops"
capabilities = ["docker", "kubernetes"]  # ✓ Has kubernetes
```

## See Also

- [GasTown](gastown.md) - Lower-level bead orchestration
- [Choosing a Target](choosing-a-target.md) - Compare with other targets
- [CLI: export](../cli/export.md) - Export command reference
