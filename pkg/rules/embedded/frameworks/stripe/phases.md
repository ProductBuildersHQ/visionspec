# Stripe Framework Phases

API-first, developer experience focused methodology.

## Overview

Stripe's approach prioritizes developer experience (DX) and API design. The API contract is defined before implementation, documentation is a product feature, and the bar for developer friction is extremely high.

```
API Contract (define the interface first)
    ↓
MRD (developer pain points)
    ↓
DX Review (developer experience critique)
    ↓
PRD (API requirements)
    ↓
UXD (docs, examples, error messages)
    ↓
TRD (implementation behind API)
    ↓
TPD (API contract tests, integration tests)
    ↓
IRD
    ↓
spec.md
```

## Core Principle: API First

> "The API is the product. The implementation is a detail."

Everything flows from the API contract:
- Documentation is written from the API
- Tests validate the API contract
- Internal implementation serves the API
- Error messages are designed, not accidents

## Phase 1: API Contract

**Goal**: Design the ideal API before any implementation.

### API Design Principles

| Principle | Description |
|-----------|-------------|
| Predictable | Consistent patterns, no surprises |
| Idiomatic | Feels natural in target languages |
| Complete | Does everything needed, nothing more |
| Documented | Self-explanatory, well-documented |
| Recoverable | Errors guide to resolution |

### API Contract Template

```markdown
# API: [Resource Name]

## Overview
[What this API does, who it's for]

## Authentication
[How requests are authenticated]

## Endpoints

### Create [Resource]

`POST /v1/resources`

**Request:**
```json
{
  "name": "string (required)",
  "description": "string (optional)",
  "metadata": "object (optional)"
}
```

**Response:**
```json
{
  "id": "res_xxx",
  "object": "resource",
  "name": "string",
  "created": 1234567890
}
```

**Errors:**
| Code | Meaning | Resolution |
|------|---------|------------|
| 400 | Invalid name | Name must be 1-100 characters |
| 401 | Invalid API key | Check API key in dashboard |
| 429 | Rate limited | Implement exponential backoff |

### List [Resources]

`GET /v1/resources`

[Same structure...]
```

### API Design Review Questions

Before proceeding, answer:

- [ ] Can a developer use this in < 15 minutes?
- [ ] Is every parameter name self-explanatory?
- [ ] Are errors actionable?
- [ ] Is the pagination obvious?
- [ ] Does it match our existing API patterns?

## Phase 2: Developer Pain Points (MRD)

**Goal**: Understand what developers struggle with today.

### Workflow

```bash
visionspec create mrd -p <project>
```

MRD for Stripe-style should focus on:

- [ ] Current developer friction points
- [ ] Support ticket patterns
- [ ] Community feedback (forums, GitHub issues)
- [ ] Competitive analysis (other APIs)
- [ ] Integration complexity

### Pain Point Categories

| Category | Example |
|----------|---------|
| Onboarding | "Took 3 hours to make first request" |
| Documentation | "Couldn't find webhooks docs" |
| Error handling | "Error message didn't help" |
| Edge cases | "Worked in test, failed in production" |
| Integration | "Hard to combine with our stack" |

### Developer Persona

```markdown
## Persona: [Name]

**Role**: [Backend engineer, Full-stack, etc.]
**Experience**: [Junior, Mid, Senior]
**Context**: [Startup, Enterprise, etc.]

**Goals:**
- Get to "hello world" fast
- Understand the API intuitively
- Handle errors gracefully

**Frustrations:**
- Poor error messages
- Incomplete examples
- Inconsistent patterns
```

## Phase 3: DX Review

**Goal**: Ruthlessly critique the developer experience.

### DX Review Checklist

#### Time to Hello World
- [ ] Can install SDK in one command?
- [ ] Can make first request in < 5 minutes?
- [ ] Is authentication straightforward?

#### Documentation
- [ ] Is every endpoint documented?
- [ ] Are there working code examples?
- [ ] Is there a quick start guide?
- [ ] Can you copy-paste examples that work?

#### Error Experience
- [ ] Do errors explain what went wrong?
- [ ] Do errors tell how to fix it?
- [ ] Are error codes consistent?
- [ ] Is there an error reference?

#### SDK Experience
- [ ] Are SDKs idiomatic for each language?
- [ ] Do SDKs handle retries?
- [ ] Do SDKs validate inputs client-side?
- [ ] Are types/autocomplete helpful?

### DX Review Process

1. **Fresh eyes test** - Have someone unfamiliar try to integrate
2. **Time tracking** - Measure time to first success
3. **Friction log** - Document every moment of confusion
4. **Competitive comparison** - How do competitors do this?

### Output

Document in PRD "Developer Experience Requirements" section.

## Phase 4: Requirements (PRD)

**Goal**: Define requirements from API and DX perspective.

### Workflow

```bash
visionspec create prd -p <project>
```

PRD should include:

- [ ] API contract (from Phase 1)
- [ ] DX requirements (from Phase 3)
- [ ] SDK requirements by language
- [ ] Documentation requirements
- [ ] Error message requirements

### Requirement Format

```markdown
## Requirement: [ID]

**User Story:**
As a developer, I want to [action], so I can [outcome].

**API Contract:**
[Endpoint, request, response]

**DX Criteria:**
- Time to success: < [X] minutes
- Error case: [Specific scenario] → [Helpful message]
- Documentation: [What must be documented]

**Acceptance Criteria:**
- [ ] Copy-paste example works
- [ ] SDK method feels natural
- [ ] Error messages are actionable
```

## Phase 5: UXD (Documentation as Product)

**Goal**: Design documentation with the same rigor as UI.

### Workflow

```bash
visionspec create uxd -p <project>
```

UXD for API products covers:

- [ ] Documentation structure
- [ ] Code example design
- [ ] Error message UX
- [ ] Interactive elements (API explorer)

### Documentation Structure

```markdown
## Documentation Plan

### Getting Started
1. Create account (< 2 min)
2. Get API key (< 1 min)
3. Make first request (< 5 min)
4. Understand response

### Reference
- Every endpoint
- Every parameter
- Every error code
- Every event type

### Guides
- Common use cases
- Integration patterns
- Migration guides
- Troubleshooting

### Examples
- Working code in every language
- Real scenarios, not "foo bar"
- Error handling included
```

### Error Message Design

Errors should be designed, not just coded:

| Component | Purpose |
|-----------|---------|
| Error code | Machine-readable identifier |
| Type | Category of error |
| Message | Human-readable explanation |
| Doc URL | Link to more information |
| Suggestion | How to fix it |

**Example:**

```json
{
  "error": {
    "code": "invalid_card_number",
    "type": "card_error",
    "message": "The card number is not a valid credit card number.",
    "doc_url": "https://stripe.com/docs/error-codes/invalid-card-number",
    "suggestion": "Check that the card number was entered correctly."
  }
}
```

## Phase 6: Technical (TRD)

**Goal**: Implementation that serves the API contract.

### Workflow

```bash
visionspec context gather -p <project>
visionspec synthesize trd -p <project>
```

TRD priorities for API-first:

1. **API contract compliance** - Implementation must match contract exactly
2. **Backward compatibility** - Never break existing integrations
3. **Performance** - API response times are product SLAs
4. **Idempotency** - Support safe retries

### Versioning Strategy

```markdown
## API Versioning

**Pattern**: Date-based versions (2024-01-15)

**Rules:**
- Breaking changes require new version
- Old versions supported for 2 years
- Version passed via header or URL

**What's breaking:**
- Removing a field
- Changing a field type
- Changing error codes
- Changing required parameters

**What's NOT breaking:**
- Adding new fields
- Adding new endpoints
- Adding new optional parameters
```

## Phase 7: Testing (TPD)

**Goal**: Test the API contract, not just implementation.

### Workflow

```bash
visionspec synthesize tpd -p <project>
```

### Contract Testing

```markdown
## Contract Tests

### Endpoint: POST /v1/resources

**Success cases:**
- Valid request → 200 + correct response shape
- Optional fields omitted → 200 + defaults applied

**Error cases:**
- Missing required field → 400 + helpful message
- Invalid field type → 400 + field identified
- Invalid auth → 401 + clear message

**Edge cases:**
- Max length inputs → Handled gracefully
- Unicode → Supported correctly
- Empty strings → Validated appropriately
```

### SDK Testing

- [ ] Every SDK method works as documented
- [ ] Error handling works across SDKs
- [ ] Retries work correctly
- [ ] Types are accurate

## Stripe Framework Gates

| Gate | Criteria |
|------|----------|
| API contract approved | Reviewed by API team, consistent with patterns |
| DX review passed | Fresh eyes test < 15 min to success |
| Documentation complete | Every endpoint, example, error documented |
| Contract tests passing | All API behaviors verified |

## See Also

- [Stripe API Design](https://stripe.com/blog/api-design)
- [Stripe Engineering Blog](https://stripe.com/blog/engineering)
- [API Design Patterns](https://www.oreilly.com/library/view/api-design-patterns/9781617295850/) - JJ Geewax
