# PetStore API Example

This example demonstrates the complete VisionSpec workflow using the classic PetStore API concept.

## What We're Building

A RESTful API for a pet store that allows:

- Managing pets (CRUD operations)
- Managing inventory
- Processing orders
- User authentication

This is based on the [OpenAPI PetStore](https://petstore.swagger.io/) specification, reimagined through VisionSpec's Working Backwards methodology.

## The Complete Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. IDEATION (Working Backwards)                                  │
│    MRD → Press Release → FAQ → PRD                              │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. EXPERIENCE DESIGN                                             │
│    UXD (User journeys, API contracts)                           │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. TECHNICAL SPECIFICATION                                       │
│    TRD → TPD → IRD                                              │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 4. RECONCILIATION                                                │
│    Conflict detection → spec.md                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 5. EXECUTION                                                     │
│    Export to target → AI agent implementation                   │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Initialize the project
visionspec init petstore-api --profile startup

# Follow the Working Backwards flow
visionspec create mrd -p petstore-api
visionspec synthesize press -p petstore-api
visionspec synthesize faq -p petstore-api
visionspec synthesize prd -p petstore-api

# Add user experience
visionspec create uxd -p petstore-api

# Generate technical specs
visionspec synthesize trd -p petstore-api
visionspec synthesize tpd -p petstore-api
visionspec synthesize ird -p petstore-api

# Evaluate and approve
visionspec eval all -p petstore-api
visionspec approve all -p petstore-api

# Reconcile
visionspec reconcile -p petstore-api

# Export to your chosen target
visionspec export aidlc -p petstore-api
```

## Project Structure

After completing the workflow:

```
docs/specs/petstore-api/
├── source/
│   ├── mrd.md          # Market Requirements
│   ├── prd.md          # Product Requirements
│   └── uxd.md          # User Experience Design
├── gtm/
│   ├── press.md        # Press Release
│   ├── faq.md          # FAQ
│   └── narrative-1p.md # Executive Summary
├── technical/
│   ├── trd.md          # Technical Requirements
│   ├── tpd.md          # Test Plan
│   └── ird.md          # Infrastructure Requirements
├── eval/
│   ├── mrd.eval.json
│   ├── prd.eval.json
│   └── trd.eval.json
├── spec.md             # Reconciled specification
└── visionspec.yaml     # Project configuration
```

## Detailed Walkthrough

- [Working Backwards](working-backwards.md) - MRD → Press → FAQ → PRD flow
- [Technical Specification](technical-spec.md) - TRD, TPD, IRD creation
- [Execution](execution.md) - Exporting and running with AI agents

## Sample Specifications

### MRD Summary

```markdown
# PetStore API - Market Requirements

## Market Problem
Pet store owners need a modern, API-first platform to manage
their inventory, process orders, and serve customers across
multiple channels (web, mobile, in-store).

## Target Audience
- Small to medium pet store owners
- Pet store chains seeking modernization
- Third-party developers building pet-related applications

## Business Goals
- Capture 10% of small pet store market in Year 1
- Enable 1000+ third-party integrations
- Reduce store operation costs by 25%
```

### Press Release Summary

```markdown
# PetStore API Launches Modern Platform for Pet Retailers

**San Francisco, CA** — Today we announce PetStore API, the
first truly API-first platform designed specifically for pet
retailers.

Pet store owners can now manage their entire operation through
a single, powerful API that integrates with any system...

**Customer Quote:**
"PetStore API reduced our inventory management time by 60% and
enabled us to launch our mobile app in weeks instead of months."
— Jane Smith, Owner, Happy Paws Pet Shop
```

### TRD Summary

```markdown
# Technical Requirements Document

## Architecture
- RESTful API with OpenAPI 3.0 specification
- Go backend with Chi router
- PostgreSQL database
- Redis caching layer
- JWT authentication

## API Endpoints
- `/pets` - Pet management (CRUD)
- `/store/inventory` - Inventory management
- `/store/order` - Order processing
- `/user` - User authentication

## Non-Functional Requirements
- Response time < 200ms (p95)
- 99.9% availability
- Support 10,000 concurrent users
```

## Execution Targets

This example can be exported to any target:

### AWS AI-DLC

```bash
visionspec export aidlc -p petstore-api
```

Best for: Enterprise implementation with approval gates

### SpecKit

```bash
visionspec export speckit -p petstore-api
```

Best for: GitHub-native PR workflow

### GSD

```bash
visionspec export gsd -p petstore-api
```

Best for: Fast parallel implementation

### GasTown/GasCity

```bash
visionspec export gastown -p petstore-api
visionspec export gascity -p petstore-api
```

Best for: Multi-agent orchestration

## Key Learnings

1. **Start with the customer** - The Press Release forces customer-first thinking
2. **FAQ surfaces gaps** - Challenging questions reveal missing requirements
3. **TRD grounds reality** - Technical specs prevent over-engineering
4. **Reconciliation catches conflicts** - Automated conflict detection saves time
5. **Export flexibility** - Same specs work with any execution target

## Next Steps

- [Working Backwards Walkthrough](working-backwards.md)
- [Technical Specification](technical-spec.md)
- [Execution Guide](execution.md)
