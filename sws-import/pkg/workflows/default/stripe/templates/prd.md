# API Product Requirements: {{.Title}}

> **Developer-Centric**: Every requirement is written from the API consumer's perspective.

## Executive Summary

One paragraph describing what this API enables developers to do and why it matters.

---

## Developer Audience

### Primary Personas

| Persona | Description | Integration Complexity | Volume |
|---------|-------------|------------------------|--------|
| **Startup Developer** | Solo developer building MVP. Values simplicity and fast time-to-integration. | Simple | Low |
| **Platform Engineer** | Building internal tooling. Needs reliability and comprehensive error handling. | Moderate | Medium |
| **Enterprise Integrator** | Large-scale implementation. Requires compliance, auditing, and bulk operations. | Complex | High |

### Use Case Priority

| Priority | Use Case | Persona | Value |
|----------|----------|---------|-------|
| P0 | {{primary_use_case}} | All | Core functionality |
| P1 | {{secondary_use_case}} | Platform Engineer | Scale operations |
| P2 | {{tertiary_use_case}} | Enterprise | Compliance/audit |

---

## Developer Journeys

### Journey 1: First Integration (Time to Hello World)

**Goal**: Developer goes from zero to first successful API call.

**Success Metric**: < 15 minutes from signup to first successful response.

| Step | Developer Action | Our Response | Friction Points |
|------|------------------|--------------|-----------------|
| 1 | Finds documentation | Clear landing page | None |
| 2 | Gets API key | Self-service, instant | Must verify email |
| 3 | Makes test call | Returns sample data | Auth format unclear? |
| 4 | Understands response | Well-documented schema | Field meanings? |
| 5 | Handles first error | Clear error message + remediation | None |

**Required API Capabilities**:

- [ ] Test mode with realistic sample data
- [ ] Copy-paste curl commands in docs
- [ ] Clear error messages with `doc_url` links
- [ ] Interactive API explorer

---

### Journey 2: Production Integration

**Goal**: Developer moves from test to production with confidence.

**Success Metric**: Zero support tickets during production go-live.

| Step | Developer Action | Our Response | Friction Points |
|------|------------------|--------------|-----------------|
| 1 | Reviews production checklist | Checklist in docs | Checklist discoverable? |
| 2 | Configures webhooks | Easy setup + testing tool | Webhook debugging |
| 3 | Handles errors gracefully | All error codes documented | Edge cases |
| 4 | Implements idempotency | Clear idempotency docs | Pattern unclear? |
| 5 | Monitors integration | Dashboard + alerts | Metrics available? |

**Required API Capabilities**:

- [ ] Production readiness checklist
- [ ] Webhook testing/replay tool
- [ ] Complete error code documentation
- [ ] Idempotency support with examples
- [ ] API usage dashboard

---

### Journey 3: Scaling Integration

**Goal**: Developer handles growth without re-architecting.

| Step | Developer Action | Our Response |
|------|------------------|--------------|
| 1 | Needs bulk operations | Batch endpoints available |
| 2 | Hits rate limits | Clear limits + upgrade path |
| 3 | Needs async processing | Webhook-based patterns |
| 4 | Wants to reduce latency | Edge caching guidance |

**Required API Capabilities**:

- [ ] Batch/bulk endpoints
- [ ] Transparent rate limits with headers
- [ ] Async operation support
- [ ] Performance optimization docs

---

## Functional Requirements

### Core Operations

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| FR-001 | Create {{resource}} via POST | P0 | Returns created object |
| FR-002 | Retrieve {{resource}} via GET | P0 | By ID |
| FR-003 | Update {{resource}} via POST | P0 | Partial updates supported |
| FR-004 | Delete {{resource}} via DELETE | P0 | Soft delete vs hard delete? |
| FR-005 | List {{resource}}s with pagination | P0 | Cursor-based, max 100 per page |

### Query & Filtering

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| FR-010 | Filter by status | P1 | Exact match |
| FR-011 | Filter by created date | P1 | gt, gte, lt, lte operators |
| FR-012 | Filter by metadata | P2 | Key existence + value match |
| FR-013 | Search by name | P2 | Prefix matching |

### Advanced Operations

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| FR-020 | Batch create (up to 100) | P1 | Returns array of results |
| FR-021 | Batch update | P2 | Atomic or partial success? |
| FR-022 | Export all data | P2 | Async job with webhook |

---

## Integration Patterns

### Pattern 1: Synchronous Request/Response

**When to use**: Real-time operations, user-facing flows.

```
Client → POST /v1/{{resource}}s → 201 Created → Continue
```

**Requirements**:

- Response time < 500ms p99
- Clear error responses for immediate handling
- Idempotency support for safe retries

### Pattern 2: Webhook-Driven

**When to use**: Async operations, event-driven architectures.

```
Client → POST /v1/{{resource}}s → 202 Accepted
         ← webhook: {{resource}}.created
```

**Requirements**:

- Webhook delivery within 30 seconds
- Retry with exponential backoff
- Webhook signature verification
- Event replay capability

### Pattern 3: Polling

**When to use**: When webhooks aren't possible.

```
Client → POST /v1/{{resource}}s → 202 Accepted
         → GET /v1/{{resource}}s/:id (poll until status != pending)
```

**Requirements**:

- Status field clearly indicates completion
- Polling guidance in documentation
- Rate limit friendly polling intervals

---

## Error Handling Requirements

### Developer Expectations

| Scenario | Developer Expects | Our Response |
|----------|-------------------|--------------|
| Invalid parameter | Know which parameter and why | `param` field + specific message |
| Missing required field | Know what's missing | `parameter_missing` code |
| Resource not found | Confirm it doesn't exist | 404 with clear message |
| Rate limited | Know when to retry | 429 + `Retry-After` header |
| Server error | Know we're aware | 500 + incident ID for support |

### Error Message Quality

Every error message must answer:

1. **What happened?** (error type)
2. **Why?** (specific reason)
3. **What to do?** (remediation)
4. **Where to learn more?** (doc_url)

---

## SDK Requirements

### Supported Languages (Priority Order)

| Language | Priority | Notes |
|----------|----------|-------|
| Node.js / TypeScript | P0 | Most common |
| Python | P0 | Data/ML use cases |
| Ruby | P1 | Rails ecosystem |
| Go | P1 | Infrastructure use cases |
| Java | P2 | Enterprise |
| PHP | P2 | WordPress/Laravel |

### SDK Principles

- [ ] Idiomatic to each language
- [ ] Type-safe where possible
- [ ] Automatic retries with backoff
- [ ] Configurable timeouts
- [ ] Easy error handling
- [ ] Webhook signature verification built-in

---

## Documentation Requirements

### Structure

| Section | Purpose | Must Include |
|---------|---------|--------------|
| Quickstart | First API call | Copy-paste example, < 5 min |
| Authentication | How to auth | Key types, security best practices |
| API Reference | Every endpoint | Request/response, errors, examples |
| Guides | Common tasks | Step-by-step with code |
| Changelog | What changed | Breaking changes highlighted |

### Every Endpoint Must Have

- [ ] Description of what it does
- [ ] All parameters with types and constraints
- [ ] Example request (curl + SDK)
- [ ] Example response
- [ ] Possible error codes
- [ ] Rate limit information

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Time to First Call | < 15 min | Analytics: signup → first 200 |
| API Error Rate | < 0.1% | 4xx/5xx responses |
| Documentation NPS | > 50 | Survey |
| Support Tickets per Integration | < 2 | Support system |
| SDK Adoption | > 70% | User agent tracking |

---

## Open Questions

| Question | Options | Decision Owner | Due Date |
|----------|---------|----------------|----------|
| Soft delete vs hard delete? | Soft (recoverable) / Hard (permanent) | API Lead | |
| Batch failure mode? | Atomic / Partial success | API Lead | |
| Webhook retry policy? | 3 retries / 5 retries / Configurable | Platform | |

---

## Appendix: Competitive Analysis

| Competitor | Strength | Weakness | Our Differentiation |
|------------|----------|----------|---------------------|
| | | | |
