# GSD (Get Shit Done)

GSD is a pragmatic execution methodology focused on rapid iteration with wave-based parallel execution and artifact verification.

## Overview

GSD provides:

- **Wave-based execution**: Parallel tasks within dependency waves
- **Artifact verification**: Explicit verification of deliverables
- **Progress tracking**: `STATE.md` for execution state
- **Fast iteration**: Minimal ceremony, maximum velocity

## When to Use GSD

GSD is ideal for:

- ✅ Rapid prototyping and MVPs
- ✅ Hackathons and time-boxed projects
- ✅ Teams that value speed over ceremony
- ✅ Projects with parallelizable work streams
- ✅ Lean Startup experimentation cycles

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
│         visionspec export gsd                                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                        GSD                                       │
│                                                                  │
│  PLAN.md       (plan with YAML frontmatter)                     │
│  STATE.md      (progress tracking)                              │
│  config.json   (GSD configuration)                              │
└─────────────────────────────────────────────────────────────────┘
```

## Export Format

When you run `visionspec export gsd`, VisionSpec creates:

```
gsd/
├── PLAN.md       # Plan with must_haves, truths, artifacts
├── STATE.md      # Progress tracking
└── config.json   # GSD configuration
```

### PLAN.md

The plan with YAML frontmatter defining constraints:

```markdown
---
must_haves:
  - Pet CRUD API endpoints
  - OpenAPI specification
  - Unit test coverage > 80%

truths:
  - Using Go 1.22+
  - PostgreSQL for persistence
  - RESTful API design

artifacts:
  - cmd/server/main.go
  - api/openapi.yaml
  - internal/pet/handler.go
  - internal/pet/handler_test.go
---

# PetStore API Plan

## Wave 1: Foundation
> Parallel: No dependencies

- [ ] Initialize Go module
- [ ] Set up project structure
- [ ] Configure PostgreSQL connection

## Wave 2: Core
> Parallel: Depends on Wave 1

- [ ] Implement Pet model
- [ ] Create Pet repository
- [ ] Build Pet handler

## Wave 3: API
> Parallel: Depends on Wave 2

- [ ] Define OpenAPI spec
- [ ] Add input validation
- [ ] Implement error handling

## Wave 4: Quality
> Parallel: Depends on Wave 3

- [ ] Write unit tests
- [ ] Add integration tests
- [ ] Generate API documentation
```

### STATE.md

Progress tracking:

```markdown
# Execution State

## Current Wave: 2

## Completed
- [x] Wave 1: Foundation (3/3 tasks)

## In Progress
- [ ] Wave 2: Core (1/3 tasks)
  - [x] Implement Pet model
  - [ ] Create Pet repository
  - [ ] Build Pet handler

## Blocked
None

## Artifacts Verified
- [x] cmd/server/main.go
- [x] go.mod
- [ ] api/openapi.yaml
```

### config.json

GSD configuration:

```json
{
  "project": "petstore-api",
  "model_profile": "balanced",
  "verification": {
    "enabled": true,
    "strict": false
  },
  "waves": {
    "parallel_tasks": 3,
    "timeout_minutes": 30
  }
}
```

## Complete Workflow

### Step 1: Create Specifications

```bash
# Initialize with lean profile for speed
visionspec init petstore-api --profile startup

# Quick spec creation
visionspec create mrd -p petstore-api
visionspec synthesize prd -p petstore-api
visionspec synthesize trd -p petstore-api

# Fast reconciliation
visionspec reconcile -p petstore-api
```

### Step 2: Export to GSD

```bash
visionspec export gsd -p petstore-api
```

Output:

```
⋯ Exporting to gsd...
✓ Exported to GSD format
  Output: gsd/
  Files:
    - PLAN.md (4 waves, 12 tasks)
    - STATE.md
    - config.json
```

### Step 3: Execute with GSD

In your AI coding agent:

```
Using GSD, execute PLAN.md
Start with Wave 1: Foundation
```

The agent will:

1. Read `PLAN.md` for context and constraints
2. Execute tasks within each wave in parallel
3. Update `STATE.md` after each task
4. Verify artifacts exist
5. Proceed to next wave when current completes

### Step 4: Track Progress

Monitor `STATE.md` for real-time progress:

```bash
# Watch progress
watch cat gsd/STATE.md

# Or check artifacts
ls -la cmd/server/ internal/pet/
```

## Configuration

Configure GSD export in `visionspec.yaml`:

```yaml
targets:
  gsd:
    enabled: true
    output_dir: gsd

    # Execution profile
    model_profile: balanced  # fast, balanced, thorough

    # Wave configuration
    max_parallel_tasks: 3
    wave_timeout_minutes: 30

    # Verification
    verify_artifacts: true
    strict_verification: false
```

## Wave-Based Execution

### How Waves Work

```
Wave 1          Wave 2          Wave 3          Wave 4
┌─────┐        ┌─────┐        ┌─────┐        ┌─────┐
│ T1  │───────→│ T4  │───────→│ T7  │───────→│ T10 │
├─────┤        ├─────┤        ├─────┤        ├─────┤
│ T2  │ parallel│ T5  │ parallel│ T8  │ parallel│ T11 │
├─────┤        ├─────┤        ├─────┤        ├─────┤
│ T3  │        │ T6  │        │ T9  │        │ T12 │
└─────┘        └─────┘        └─────┘        └─────┘
   │              │              │              │
   └──────────────┴──────────────┴──────────────┘
              Sequential between waves
```

### Defining Waves

Waves are determined by dependencies:

| Task | Dependencies | Wave |
|------|--------------|------|
| Go module init | None | 1 |
| Project structure | None | 1 |
| DB connection | None | 1 |
| Pet model | Go module | 2 |
| Pet repository | Pet model | 2 |
| Pet handler | Pet model | 2 |
| OpenAPI spec | Pet handler | 3 |
| Unit tests | Pet handler | 3 |

## Artifact Verification

### Must-Have Artifacts

Define required deliverables:

```yaml
---
artifacts:
  - cmd/server/main.go      # Entry point
  - api/openapi.yaml        # API specification
  - internal/pet/handler.go # Core handler
  - Makefile                # Build automation
---
```

### Verification Process

After each wave:

1. GSD checks if artifacts exist
2. Optionally validates file content
3. Updates `STATE.md` with verification status
4. Blocks next wave if strict verification fails

## Best Practices

### 1. Keep Waves Small

Each wave should complete in 15-30 minutes:

```markdown
## Wave 1: Foundation (15 min)
- [ ] Init module
- [ ] Create directories
- [ ] Add Makefile

## Wave 2: Models (20 min)
- [ ] Pet struct
- [ ] Store struct
- [ ] User struct
```

### 2. Parallelize Aggressively

Tasks without dependencies should be in the same wave:

```markdown
## Wave 2: Core (parallel)
- [ ] Pet CRUD        # No interdependency
- [ ] Store CRUD      # No interdependency
- [ ] User CRUD       # No interdependency
```

### 3. Define Clear Artifacts

Be specific about deliverables:

```yaml
artifacts:
  # ✅ Good: Specific files
  - internal/pet/handler.go
  - internal/pet/handler_test.go

  # ❌ Avoid: Vague patterns
  - internal/**/*.go
```

### 4. Use Truths for Constraints

Lock down technical decisions:

```yaml
truths:
  - Go 1.22+ required
  - PostgreSQL 15+
  - No ORM (use sqlx)
  - Error wrapping with fmt.Errorf
```

## Troubleshooting

### "Wave blocked"

Check for failed verifications:

```bash
# View state
cat gsd/STATE.md

# Check missing artifacts
grep "^\- \[ \]" gsd/STATE.md
```

### "Artifacts not found"

Verify paths in PLAN.md match actual structure:

```bash
# List generated files
find . -name "*.go" -type f

# Update artifacts list if needed
```

### "Parallel tasks conflicting"

Reduce parallelism:

```yaml
targets:
  gsd:
    max_parallel_tasks: 1  # Sequential within wave
```

## See Also

- [Choosing a Target](choosing-a-target.md) - Compare with other targets
- [Lean Startup Framework](../frameworks/lean-startup.md) - Natural pairing with GSD
- [CLI: export](../cli/export.md) - Export command reference
