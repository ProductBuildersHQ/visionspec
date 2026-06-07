# Examples

Learn VisionSpec through complete, working examples.

## Available Examples

### [PetStore API](petstore/index.md)

A classic API example demonstrating the full VisionSpec workflow:

- Working Backwards specification flow
- All spec types (MRD, Press, FAQ, PRD, TRD, TPD, IRD)
- Export to multiple execution targets
- Implementation with AI-DLC

**Best for**: Learning the complete VisionSpec workflow

```bash
# Quick start
visionspec init petstore-api --profile startup
```

## Example Structure

Each example includes:

1. **Overview** - What we're building and why
2. **Specifications** - Complete spec files
3. **Working Backwards** - Step-by-step ideation flow
4. **Export** - Exporting to execution targets
5. **Execution** - Running with AI coding agents

## Running Examples Locally

Clone the VisionSpec repository and navigate to examples:

```bash
git clone https://github.com/ProductBuildersHQ/visionspec
cd visionspec/examples
```

Or initialize your own version:

```bash
visionspec init petstore-api --profile startup
```

## Creating Your Own Examples

Use the PetStore pattern as a template:

1. **Start with MRD** - Define the market problem
2. **Synthesize Working Backwards** - Generate Press, FAQ, PRD
3. **Author UXD** - Define user experience
4. **Synthesize Technical** - Generate TRD, TPD, IRD
5. **Reconcile** - Create unified spec.md
6. **Export** - Choose your execution target
7. **Execute** - Implement with AI agents

## Profiles for Examples

| Example Type | Profile | Framework |
|--------------|---------|-----------|
| Quick prototype | `startup` | Lean Startup |
| API service | `growth` | Stripe API-First |
| Enterprise app | `enterprise` | AWS Working Backwards |
| Design-focused | `startup` | Design Thinking |

## Next Steps

- [PetStore API Example](petstore/index.md) - Complete walkthrough
- [Quick Start](../getting-started/quickstart.md) - Get started in 5 minutes
- [Execution Targets](../execution-targets/index.md) - Choose your target
