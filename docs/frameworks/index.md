# Methodology Frameworks

VisionSpec supports multiple product development methodologies through its profile system. Each profile provides customized templates and evaluation rubrics tailored to the methodology's artifacts and practices.

## Available Frameworks

| Profile | Methodology | Best For |
|---------|-------------|----------|
| [Big Tech Product](big-tech.md) | AWS + Google + Stripe + Netflix + more (MRD start) | New product lines with comprehensive practices |
| [Big Tech Feature](big-tech.md#big-tech-feature) | AWS + Google + Stripe + Netflix + more (OpportunitySpec start) | Features with comprehensive practices |
| [AWS Product](aws.md) | Working Backwards (MRD start) | New product lines, major initiatives |
| [AWS Feature](opportunity-spec.md) | Working Backwards (OpportunitySpec start) | Features on existing products |
| [Google](google.md) | Design Docs + RFC | Engineering-heavy, data-driven organizations |
| [Stripe](stripe.md) | API-First | Platform and API products |
| [Shape Up](shapeup.md) | Basecamp (Ryan Singer) | Fixed time, variable scope, betting on pitches |
| [Continuous Discovery](continuous-discovery.md) | Teresa Torres | Weekly touchpoints, assumption testing |
| [Lean Startup](lean-startup.md) | Build-Measure-Learn | Early-stage validation |
| [Design Thinking](design-thinking.md) | Stanford d.school | Human-centered design |
| [JTBD](jtbd.md) | Jobs to be Done | Understanding customer motivations |

## Choosing a Framework

### By Scope

| Scope | Starting Document | Recommended Profile |
|-------|-------------------|---------------------|
| New product line | MRD | Big Tech Product, AWS Product, Lean Startup |
| New feature | OpportunitySpec / OST | Big Tech Feature, AWS Feature, Continuous Discovery |
| API/Platform feature | OpportunitySpec | Big Tech Feature, Stripe |
| Fixed-time project | Pitch | Shape Up |
| Continuous improvement | OST | Continuous Discovery |
| Experiment | Hypothesis | Lean Startup, Continuous Discovery |

### By Company Stage

| Stage | Recommended Profiles |
|-------|---------------------|
| Pre-PMF (0-1) | Lean Startup, Continuous Discovery, Design Thinking, JTBD |
| Startup | Lean Startup, Shape Up, Continuous Discovery |
| Growth | Big Tech Product/Feature, Google, Stripe, Continuous Discovery |
| Enterprise | Big Tech Product/Feature, Google, AWS Product, Stripe |

### By Product Type

| Product Type | Recommended Profiles |
|--------------|---------------------|
| Consumer App | Big Tech Product, Design Thinking, JTBD, AWS Product |
| B2B SaaS | Big Tech Product, AWS Product, Google, Lean Startup |
| Platform/API | Big Tech Feature, Stripe, Google |
| Internal Tools | Google, Big Tech Feature |
| Hardware | Design Thinking, JTBD |

### By Team Culture

| Culture | Recommended Profiles |
|---------|---------------------|
| Writing-heavy | Big Tech Product, AWS Product, Stripe, Google |
| Engineering-led | Big Tech Product/Feature, Google, Stripe |
| Design-led | Design Thinking, JTBD, Continuous Discovery |
| Data-driven | Big Tech Product/Feature, Google, Lean Startup |
| Customer-obsessed | Big Tech Product, AWS Product, JTBD, Continuous Discovery |
| High autonomy | Shape Up, Big Tech Product/Feature (Netflix/Spotify) |
| Research-integrated | Continuous Discovery |

## Using Profiles

### List Available Profiles

```bash
multispec profiles list
```

### Initialize Project with Profile

```bash
multispec init my-project --profile aws
```

### Show Profile Details

```bash
multispec profiles show google
```

### Export Profile for Customization

```bash
multispec profiles export stripe ./my-profiles/stripe-custom
```

## Combining Frameworks

Profiles can be extended and combined. For example:

```yaml
# multispec.yaml
profile: aws
extends:
  - lean-startup  # Add hypothesis validation

# Or create a custom profile
profile: ./profiles/my-hybrid.yaml
```

## VisionSpec Document Mapping

Each framework maps its artifacts to VisionSpec's document types:

| VisionSpec | Big Tech | AWS | Google | Shape Up | Discovery | Lean | Design | JTBD |
|------------|----------|-----|--------|----------|-----------|------|--------|------|
| **MRD** | MRD + OKRs | Business Case | OKR Doc | — | OST | Hypothesis | Empathy Map | Job Statement |
| **PRD** | PRD + API | PRFAQ | RFC | Pitch | Assumption Map | MVP PRD | Prototype Spec | Outcome Reqs |
| **UXD** | UXD + DX | — | Experiment | — | Experience Map | Experiment | Journey Map | Job Map |
| **TRD** | Design Doc | Design Doc | Design Doc | — | — | — | — | Solution Arch |
| **Narrative** | 6-Pager | 6-Pager | — | — | — | Pivot Doc | — | — |
| **Tracking** | — | — | — | Scope/Hill | Snapshot | — | — | — |
