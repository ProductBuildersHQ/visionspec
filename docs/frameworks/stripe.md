# Stripe API-First

Stripe's methodology treats APIs as products and documentation as a first-class deliverable. This approach emphasizes contract-first design, developer empathy, and precision in specifications.

## The Flow

![Stripe API-First Flow](../diagrams/stripe-flow.svg)

## Key Principles

1. **Contract-First**: API design precedes implementation
2. **Docs as Product**: Documentation is a first-class product
3. **Developer Empathy**: Specs from the API consumer's perspective
4. **Precision**: No ambiguity in inputs, outputs, errors
5. **Consistency**: Predictable patterns across the API surface
6. **Externally Legible**: Readable by external developers

## VisionSpec Mapping

| Stripe Artifact | VisionSpec Type | Purpose |
|-----------------|-----------------|---------|
| Platform Strategy | MRD | Developer market and ecosystem |
| API Product Requirements | PRD | Developer journeys and use cases |
| API Specification | TRD | Contract definition (OpenAPI-style) |
| Developer Experience Spec | UXD | Onboarding, SDKs, documentation |

## Using the Stripe Profile

### Initialize a Project

```bash
multispec init payments-api --profile stripe
```

### Create Platform Strategy (MRD)

```bash
multispec draft mrd -p payments-api
```

The Platform Strategy template includes:

- Developer market sizing (TAM/SAM/SOM)
- Developer personas and segments
- Competitive landscape
- Business model and monetization
- Go-to-market strategy

### Create API PRD

```bash
multispec draft prd -p payments-api
```

The API PRD template includes:

- Developer personas and priorities
- Developer journeys (Time to Hello World)
- Functional requirements
- Integration patterns
- Error handling requirements
- SDK requirements
- Success metrics

### Synthesize API Specification (TRD)

```bash
multispec synthesize trd -p payments-api
```

The API Specification template includes:

- Resource definitions
- Endpoint documentation
- Request/response schemas
- Error codes with remediation
- Pagination
- Idempotency
- Versioning strategy
- Webhook events

### Create Developer Experience Spec (UXD)

```bash
multispec draft uxd -p payments-api
```

The DX Spec template includes:

- Onboarding journey (< 15 min to first call)
- Documentation architecture
- SDK design principles
- Error experience quality
- Interactive tools (API explorer, webhook tester)
- DX metrics

## Rubric Categories

### Platform Strategy Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Developer Market | 20% | Market clearly defined and sized |
| Competitive Position | 15% | Differentiation clear and defensible |
| Platform Strategy | 20% | Ecosystem thinking present |
| Business Model | 15% | Monetization and unit economics |
| Go-to-Market | 15% | Developer acquisition strategy |
| Success Metrics | 10% | Platform success measurable |
| Roadmap/Investment | 5% | Phased approach justified |

### API PRD Evaluation

| Category | Weight | Description |
|----------|--------|-------------|
| Developer Personas | 15% | Audience clearly segmented |
| Developer Journeys | 20% | Integration paths mapped |
| Use Case Clarity | 20% | Specific and prioritized |
| Integration Patterns | 15% | Sync, async, webhook defined |
| Error Experience | 15% | Error handling from dev perspective |
| SDK Requirements | 10% | Languages and features defined |
| Success Metrics | 5% | Developer success measurable |

### API Specification Evaluation (TRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Contract Clarity | 25% | Precise, unambiguous contracts |
| Consistency | 20% | Predictable patterns throughout |
| Error Documentation | 20% | Every error with remediation |
| Resource Modeling | 15% | Clear relationships |
| Versioning Strategy | 10% | Breaking change policy |
| Idempotency/Safety | 5% | Safe retry behavior |
| Rate Limiting | 5% | Limits documented |

### DX Spec Evaluation (UXD)

| Category | Weight | Description |
|----------|--------|-------------|
| Onboarding Journey | 25% | < 15 min to first call |
| Documentation Architecture | 20% | Well-structured information |
| Documentation Standards | 20% | Quality requirements defined |
| SDK Design | 15% | Principles and coverage |
| Error Experience | 10% | Errors help developers succeed |
| Interactive Tools | 5% | API explorer, webhook tester |
| DX Metrics | 5% | Experience measurable |

## Example Workflow

```bash
# 1. Initialize project
multispec init billing-api --profile stripe

# 2. Define Platform Strategy
multispec draft mrd -p billing-api
multispec eval mrd -p billing-api
multispec approve mrd -p billing-api

# 3. Define API Product Requirements
multispec draft prd -p billing-api
multispec eval prd -p billing-api
multispec approve prd -p billing-api

# 4. Synthesize API Specification
multispec synthesize trd -p billing-api
multispec eval trd -p billing-api
multispec approve trd -p billing-api

# 5. Define Developer Experience
multispec draft uxd -p billing-api
multispec eval uxd -p billing-api
multispec approve uxd -p billing-api

# 6. Check status
multispec status -p billing-api

# 7. Export for implementation
multispec export speckit -p billing-api
```

## API Design Patterns

The Stripe profile encourages these patterns:

### Human-Readable IDs

```
ch_1234567890   (charge)
cus_abcdefghij  (customer)
pi_xyz123       (payment intent)
```

### Idempotency Keys

```bash
curl https://api.example.com/v1/charges \
  -H "Idempotency-Key: unique-key-12345"
```

### Expandable Objects

```bash
curl https://api.example.com/v1/charges/ch_xxx \
  -d "expand[]"=customer
```

### Consistent Error Format

```json
{
  "error": {
    "type": "invalid_request_error",
    "code": "parameter_missing",
    "message": "Required parameter 'amount' is missing",
    "param": "amount",
    "doc_url": "https://docs.example.com/errors#parameter_missing"
  }
}
```

## Reference Materials

For deeper understanding of Stripe's API-first approach, see:

- [Stripe API Reference](https://docs.stripe.com/api)
- [Stripe's payments APIs: The first 10 years](https://stripe.dev/blog/payment-api-design)
- [Inside Stripe's Engineering Culture](https://newsletter.pragmaticengineer.com/p/stripe-part-2)
- Internal reference: `frameworks-internal/stripe-api-first/`
