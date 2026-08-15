# Developer Experience Spec: {{.Title}}

> **Docs as Product**: Documentation is a first-class product, not an afterthought.

## Overview

Brief description of the developer experience goals for this API.

**Primary DX Goal**: A developer should go from zero to first successful API call in under 15 minutes.

---

## Developer Onboarding Journey

### Step 1: Discovery

**Goal**: Developer finds and understands what the API does.

| Touchpoint | Experience | Quality Bar |
|------------|------------|-------------|
| Landing page | Clear value prop, "Get Started" prominent | < 5 seconds to understand |
| Pricing page | Transparent, free tier visible | No hidden costs |
| Docs home | Organized, searchable | < 3 clicks to any page |

**Required Elements**:

- [ ] One-sentence API description
- [ ] 3-5 key use cases
- [ ] Clear "Get Started" CTA
- [ ] Pricing transparency

---

### Step 2: Signup & API Key

**Goal**: Developer gets API credentials quickly.

| Flow | Experience | Quality Bar |
|------|------------|-------------|
| Signup | Email/OAuth, minimal fields | < 2 minutes |
| Verification | Email confirmation | Instant send |
| API Key | Visible immediately, easy copy | One click to copy |

**Required Elements**:

- [ ] Test mode key by default (safe to experiment)
- [ ] Clear test vs. live mode distinction
- [ ] Key management dashboard
- [ ] Key rotation capability

---

### Step 3: First API Call

**Goal**: Developer makes successful API call and understands response.

| Resource | Experience | Quality Bar |
|----------|------------|-------------|
| Quickstart | Copy-paste code that works | < 5 minutes to success |
| API Explorer | Interactive, pre-filled with test key | Works in browser |
| Error handling | First error is educational | Links to resolution |

**Quickstart Requirements**:

```markdown
## Quickstart Structure

1. Install (one command)
2. Initialize (one line, with API key)
3. Make call (3-5 lines)
4. See result (formatted output)
```

**Example Quickstart (Node.js)**:

```javascript
// 1. Install
// npm install @example/api

// 2. Initialize
const Example = require('@example/api');
const client = new Example('sk_test_xxx');

// 3. Create a resource
const resource = await client.resources.create({
  name: 'My First Resource'
});

// 4. See the result
console.log(resource.id); // res_1234567890
```

---

### Step 4: Understanding the API

**Goal**: Developer understands API patterns and can build confidently.

| Resource | Purpose | Quality Bar |
|----------|---------|-------------|
| API Reference | Every endpoint documented | Complete, accurate |
| Guides | Common tasks explained | Step-by-step with code |
| Examples | Real-world use cases | Copy-paste ready |
| SDKs | Language-native access | Idiomatic, typed |

---

## Documentation Architecture

### Information Hierarchy

```
Documentation
├── Quickstart (< 5 min to first call)
├── Guides
│   ├── Authentication
│   ├── Error Handling
│   ├── Webhooks
│   ├── Pagination
│   └── [Use Case Guides]
├── API Reference
│   ├── Overview
│   ├── [Resource A]
│   │   ├── The Resource Object
│   │   ├── Create
│   │   ├── Retrieve
│   │   ├── Update
│   │   ├── Delete
│   │   └── List
│   └── [Resource B]
├── SDKs
│   ├── Node.js
│   ├── Python
│   └── [Others]
├── Webhooks
│   ├── Setup
│   ├── Event Types
│   └── Security
├── Errors
│   ├── Error Types
│   └── Error Codes
└── Changelog
```

### Documentation Standards

| Element | Standard | Example |
|---------|----------|---------|
| **Endpoint** | Verb + Path | `POST /v1/resources` |
| **Parameter** | Name, Type, Required, Description | `name` · string · required · The resource name |
| **Code Sample** | Working, copy-paste ready | Includes all imports |
| **Response** | Actual JSON, formatted | Full object shown |
| **Error** | Code, message, remediation | What to do |

---

## SDK Design

### SDK Principles

1. **Idiomatic**: Feels native to each language
2. **Type-Safe**: Full type definitions where supported
3. **Predictable**: Same patterns across all SDKs
4. **Resilient**: Automatic retries, timeouts, error handling
5. **Debuggable**: Clear error messages, logging options

### SDK Feature Matrix

| Feature | Node.js | Python | Go | Java |
|---------|---------|--------|-----|------|
| Full API coverage | ✓ | ✓ | ✓ | ✓ |
| TypeScript types | ✓ | N/A | ✓ | ✓ |
| Automatic retries | ✓ | ✓ | ✓ | ✓ |
| Webhook verification | ✓ | ✓ | ✓ | ✓ |
| Async support | ✓ | ✓ | ✓ | ✓ |
| Streaming | ✓ | ✓ | ✓ | ✓ |
| Pagination helpers | ✓ | ✓ | ✓ | ✓ |

### SDK Code Patterns

**Creating a Resource**:

```javascript
// Node.js - Async/await (preferred)
const resource = await client.resources.create({
  name: 'Example'
});

// Node.js - Callbacks (supported)
client.resources.create({ name: 'Example' }, (err, resource) => {
  if (err) { /* handle */ }
});
```

```python
# Python - Sync
resource = client.resources.create(name='Example')

# Python - Async
resource = await client.resources.create_async(name='Example')
```

```go
// Go
resource, err := client.Resources.Create(&ResourceParams{
    Name: stripe.String("Example"),
})
```

---

## Error Experience

### Error Message Quality Checklist

Every error must answer:

- [ ] **What happened?** (error type)
- [ ] **Where?** (which parameter/field)
- [ ] **Why?** (specific reason)
- [ ] **How to fix?** (remediation steps)
- [ ] **Learn more?** (documentation link)

### Error Examples

**Good Error**:

```json
{
  "error": {
    "type": "invalid_request_error",
    "code": "parameter_invalid",
    "message": "The 'email' parameter must be a valid email address. You provided 'not-an-email'.",
    "param": "email",
    "doc_url": "https://docs.example.com/errors#parameter_invalid"
  }
}
```

**Bad Error** (what to avoid):

```json
{
  "error": "Invalid request"
}
```

---

## Interactive Tools

### API Explorer

**Purpose**: Try API calls without writing code.

**Requirements**:

- [ ] Pre-populated with test API key
- [ ] All endpoints available
- [ ] Editable request bodies
- [ ] Formatted response display
- [ ] Copy as curl/SDK code

### Webhook Tester

**Purpose**: Test webhook integrations without deploying.

**Requirements**:

- [ ] Unique test endpoint per developer
- [ ] Send test events on demand
- [ ] View delivery attempts and payloads
- [ ] Replay failed deliveries

### Request Log

**Purpose**: Debug integration issues.

**Requirements**:

- [ ] Last 100 requests visible
- [ ] Filter by endpoint, status, time
- [ ] View full request/response
- [ ] Identify errors quickly

---

## Onboarding Metrics

| Metric | Definition | Target | Measurement |
|--------|------------|--------|-------------|
| **Time to First Call** | Signup → first 200 response | < 15 min | Analytics |
| **Quickstart Completion** | Started → finished quickstart | > 70% | Analytics |
| **Documentation Bounce** | Leave docs without action | < 30% | Analytics |
| **SDK Adoption** | API calls via SDK vs. raw HTTP | > 70% | User-agent |
| **Support Tickets** | Tickets per integration | < 2 | Support system |

---

## Content Requirements

### Every API Endpoint Needs

| Content | Required | Example |
|---------|----------|---------|
| Description | ✓ | "Creates a new resource" |
| Parameters table | ✓ | Name, type, required, description |
| Request example (curl) | ✓ | Copy-paste ready |
| Request example (SDK) | ✓ | At least 2 languages |
| Response example | ✓ | Full JSON |
| Error codes | ✓ | Possible errors for this endpoint |

### Every Guide Needs

| Content | Required | Example |
|---------|----------|---------|
| Goal statement | ✓ | "In this guide, you'll learn..." |
| Prerequisites | ✓ | What you need before starting |
| Step-by-step | ✓ | Numbered steps with code |
| Complete code | ✓ | Full working example |
| Next steps | ✓ | Related guides |

---

## Testing the DX

### DX Review Checklist

Before launch, test with developers who haven't seen the API:

- [ ] Can they complete quickstart in < 15 minutes?
- [ ] Do they understand what the API does from the homepage?
- [ ] Can they find the endpoint they need in < 3 clicks?
- [ ] Can they debug their first error without support?
- [ ] Do code samples work when copy-pasted?

### Ongoing DX Monitoring

| Signal | Source | Action Threshold |
|--------|--------|------------------|
| Time to first call increasing | Analytics | > 20 min |
| Support tickets about docs | Support | > 10/week |
| Negative doc feedback | Feedback widget | < 70% helpful |
| SDK bug reports | GitHub | Any P0 |

---

## Appendix: Documentation Style Guide

### Voice and Tone

- **Clear**: Use simple language, avoid jargon
- **Direct**: "Do X" not "You might want to consider doing X"
- **Respectful**: Assume competence, don't over-explain basics
- **Consistent**: Same patterns throughout

### Code Samples

- **Working**: All samples must execute successfully
- **Complete**: Include all imports and setup
- **Realistic**: Use realistic data, not "foo" and "bar"
- **Commented**: Explain non-obvious lines

### Formatting

| Element | Format |
|---------|--------|
| API endpoints | `POST /v1/resources` (code format) |
| Parameters | `parameter_name` (code format) |
| Values | `"string value"` (quoted) |
| Emphasis | **bold** for key terms |
