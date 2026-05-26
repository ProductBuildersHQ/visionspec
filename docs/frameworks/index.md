# Methodology Frameworks

VisionSpec supports multiple product development methodologies through its profile system. Each profile provides customized templates and evaluation rubrics tailored to the methodology's artifacts and practices.

## Available Frameworks

| Profile | Methodology | Best For |
|---------|-------------|----------|
| [AWS](aws.md) | Working Backwards | Customer-centric product development |
| [Google](google.md) | Design Docs + RFC | Engineering-heavy, data-driven organizations |
| [Stripe](stripe.md) | API-First | Platform and API products |
| [Lean Startup](lean-startup.md) | Build-Measure-Learn | Early-stage validation |
| [Design Thinking](design-thinking.md) | Stanford d.school | Human-centered design |
| [JTBD](jtbd.md) | Jobs to be Done | Understanding customer motivations |

## Choosing a Framework

### By Company Stage

| Stage | Recommended Profiles |
|-------|---------------------|
| Pre-PMF (0-1) | Lean Startup, Design Thinking, JTBD |
| Startup | Lean Startup, AWS |
| Growth | Google, Stripe, AWS |
| Enterprise | Google, AWS, Stripe |

### By Product Type

| Product Type | Recommended Profiles |
|--------------|---------------------|
| Consumer App | Design Thinking, JTBD, AWS |
| B2B SaaS | AWS, Google, Lean Startup |
| Platform/API | Stripe, Google |
| Internal Tools | Google |
| Hardware | Design Thinking, JTBD |

### By Team Culture

| Culture | Recommended Profiles |
|---------|---------------------|
| Writing-heavy | AWS, Stripe, Google |
| Engineering-led | Google, Stripe |
| Design-led | Design Thinking, JTBD |
| Data-driven | Google, Lean Startup |
| Customer-obsessed | AWS, JTBD, Design Thinking |

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

| VisionSpec | AWS | Google | Stripe | Lean | Design | JTBD |
|------------|-----|--------|--------|------|--------|------|
| **MRD** | Business Case | OKR Doc | Platform Strategy | Hypothesis | Empathy Map | Job Statement |
| **PRD** | PRFAQ | RFC | API PRD | MVP PRD | Prototype Spec | Outcome Reqs |
| **UXD** | — | Experiment | DX Spec | Experiment | Journey Map | Job Map |
| **TRD** | Design Doc | Design Doc | API Spec | — | — | Solution Arch |
| **Narrative** | 6-Pager | — | — | Pivot Doc | — | — |
