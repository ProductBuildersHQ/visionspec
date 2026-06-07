# Execution Targets

VisionSpec generates specifications. Execution targets **consume** those specifications and turn them into working software.

## What Are Execution Targets?

An execution target is an AI coding agent system that takes your reconciled `spec.md` and orchestrates the actual implementation. VisionSpec supports multiple targets, each with different strengths:

```
┌─────────────────────────────────────────────────────────────────┐
│                        VISIONSPEC                                │
│                                                                  │
│  MRD → Press → FAQ → PRD → TRD → TPD → IRD                      │
│                        ↓                                         │
│                    spec.md                                       │
│                        ↓                                         │
│              visionspec export <target>                          │
└─────────────────────────────────────────────────────────────────┘
                              ↓
         ┌────────────────────┼────────────────────┐
         ↓                    ↓                    ↓
    ┌─────────┐         ┌─────────┐         ┌─────────┐
    │ AI-DLC  │         │ SpecKit │         │   GSD   │
    │ (AWS)   │         │(GitHub) │         │         │
    └─────────┘         └─────────┘         └─────────┘
         ↓                    ↓                    ↓
    ┌─────────┐         ┌─────────┐         ┌─────────┐
    │GasTown  │         │GasCity  │         │OpenSpec │
    │         │         │         │         │(future) │
    └─────────┘         └─────────┘         └─────────┘
```

## Available Targets

| Target | Best For | Execution Model |
|--------|----------|-----------------|
| [AWS AI-DLC](aws-aidlc.md) | Enterprise workflows, multi-phase development | Three-phase lifecycle with approval gates |
| [GitHub SpecKit](speckit.md) | GitHub-native development, PR workflows | Sequential task execution |
| [GSD](gsd.md) | Fast iteration, parallel execution | Wave-based parallel tasks |
| [GasTown](gastown.md) | Complex multi-agent orchestration | DAG-based bead execution |
| [GasCity](gascity.md) | Role-based agent orchestration | Agent-order architecture |

## The Export Workflow

### Step 1: Generate Specifications

Use VisionSpec's Working Backwards methodology to create your specs:

```bash
# Initialize project
visionspec init my-project --profile enterprise

# Author source specs (human-written)
visionspec create mrd -p my-project
visionspec create prd -p my-project
visionspec create uxd -p my-project

# Synthesize derived specs (LLM-generated)
visionspec synthesize press -p my-project
visionspec synthesize faq -p my-project
visionspec synthesize trd -p my-project

# Evaluate and approve
visionspec eval all -p my-project
visionspec approve all -p my-project
```

### Step 2: Reconcile

Generate the unified execution specification:

```bash
visionspec reconcile -p my-project
# Output: docs/specs/my-project/spec.md
```

### Step 3: Export to Target

Transform `spec.md` into the target's format:

```bash
# Export to your chosen target
visionspec export aidlc -p my-project
visionspec export speckit -p my-project
visionspec export gsd -p my-project
visionspec export gastown -p my-project
visionspec export gascity -p my-project
```

### Step 4: Execute

Use the target system to implement the specification. Each target has its own trigger pattern:

```bash
# AWS AI-DLC
"Using AI-DLC, implement the project based on the vision document."

# GitHub SpecKit
"Using SpecKit, execute the plan in .specify/"

# GSD
"Using GSD, execute PLAN.md"
```

## Choosing a Target

See [Choosing a Target](choosing-a-target.md) for a detailed comparison to help you decide which execution target fits your workflow.

## Key Concepts

### Specification vs Execution

VisionSpec handles **specification** (what to build):

- Market requirements (MRD)
- Product vision (Press Release)
- Functional requirements (PRD)
- Technical architecture (TRD)
- Test strategy (TPD)
- Infrastructure needs (IRD)

Execution targets handle **implementation** (how to build it):

- Code generation
- Test creation
- Build orchestration
- Deployment

### Framework vs Target

**Frameworks** (part of VisionSpec) define how you create specifications:

- AWS Working Backwards
- Google Design Docs
- Stripe API-First
- Lean Startup
- Design Thinking
- Jobs to be Done

**Targets** (external systems) define how specifications become code:

- AWS AI-DLC
- GitHub SpecKit
- GSD
- GasTown
- GasCity

You can mix any framework with any target. For example:

- Use **Lean Startup** framework → export to **GSD** for fast iteration
- Use **AWS Working Backwards** → export to **AWS AI-DLC** for enterprise rigor
- Use **Stripe API-First** → export to **SpecKit** for API-driven development

## Next Steps

- [AWS AI-DLC](aws-aidlc.md) - Enterprise-grade development lifecycle
- [GitHub SpecKit](speckit.md) - GitHub-native spec-driven development
- [GSD](gsd.md) - Fast, parallel execution
- [GasTown](gastown.md) - Multi-agent DAG orchestration
- [GasCity](gascity.md) - Role-based agent coordination
- [Choosing a Target](choosing-a-target.md) - Decision guide
