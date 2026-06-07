# Working Backwards: PetStore API

This walkthrough demonstrates the Working Backwards methodology using the PetStore API example.

## The Working Backwards Flow

```
MRD (Market Problem)
        ↓
Press Release (Vision)
        ↓
FAQ (Challenges & Scope)
        ↓
PRD (Requirements)
```

## Step 1: Market Requirements Document (MRD)

The MRD captures the market problem, not the solution.

### Create the MRD

```bash
visionspec create mrd -p petstore-api
```

### Author the MRD

Edit `docs/specs/petstore-api/source/mrd.md`:

```markdown
# PetStore API - Market Requirements Document

## Executive Summary

Pet retailers face fragmented technology solutions that don't
integrate well, leading to inefficient operations and poor
customer experiences.

## Market Problem

### Current State
- Pet stores use disconnected systems for inventory, sales, and customers
- No modern API-first solutions designed for pet retail
- Integration with third-party services requires custom development
- Mobile and online presence is difficult to maintain

### Impact
- Store owners spend 40% more time on administrative tasks
- Lost sales due to inventory inaccuracies
- Inability to compete with large online retailers
- Customer data scattered across multiple systems

## Target Audience

### Primary: Small Pet Store Owners
- 1-3 locations
- $100K - $5M annual revenue
- Limited IT resources
- Need simple, affordable solutions

### Secondary: Pet Store Chains
- 4-20 locations
- $5M - $50M annual revenue
- Have IT staff
- Need scalable, customizable solutions

### Tertiary: Third-Party Developers
- Building pet-related applications
- Need reliable pet data APIs
- Willing to pay for quality integrations

## Business Goals

| Goal | Metric | Target |
|------|--------|--------|
| Market Penetration | Small store adoption | 10% in Year 1 |
| Ecosystem | Third-party integrations | 1,000+ |
| Efficiency | Store operation cost reduction | 25% |
| Revenue | API subscription revenue | $2M ARR |

## Competitive Landscape

| Competitor | Strength | Weakness |
|------------|----------|----------|
| Generic POS | Market presence | Not pet-specific |
| Legacy systems | Existing installs | No API, outdated |
| Custom builds | Tailored | Expensive, fragile |

## Success Criteria

1. API response time < 200ms for 95th percentile
2. 99.9% uptime SLA
3. Support 10,000 concurrent users
4. OpenAPI specification for all endpoints
5. SDK support for 3+ languages
```

### Evaluate the MRD

```bash
visionspec eval mrd -p petstore-api
```

## Step 2: Press Release (Vision)

The Press Release defines the customer-facing vision.

### Synthesize the Press Release

```bash
visionspec synthesize press -p petstore-api
```

### Review Generated Press Release

`docs/specs/petstore-api/gtm/press.md`:

```markdown
# PetStore API Launches Modern Platform for Pet Retailers

**FOR IMMEDIATE RELEASE**

**San Francisco, CA — [Date]** — Today, PetStore API announces
the launch of a modern, API-first platform designed specifically
for pet retailers, enabling store owners to manage their entire
operation through a single, powerful interface.

## The Problem We're Solving

Pet store owners have long struggled with disconnected systems
for inventory, sales, and customer management. Existing solutions
either lack modern API capabilities or weren't designed with pet
retail in mind.

"I was spending hours every week reconciling inventory across
different systems," said Jane Smith, owner of Happy Paws Pet Shop.
"There had to be a better way."

## Our Solution

PetStore API provides a complete, API-first platform that handles:

- **Pet Management**: Track pets with detailed profiles including
  breed, age, health records, and availability
- **Inventory Control**: Real-time inventory across all locations
  with automatic reorder alerts
- **Order Processing**: Seamless order management from cart to
  delivery with payment integration
- **Customer Management**: Unified customer profiles with purchase
  history and preferences

## Customer Benefits

> "PetStore API reduced our inventory management time by 60% and
> enabled us to launch our mobile app in weeks instead of months.
> The API documentation is excellent, and our developers love
> working with it."
>
> — Jane Smith, Owner, Happy Paws Pet Shop

> "As a developer building a pet adoption app, PetStore API's
> comprehensive endpoints and SDKs saved us months of development
> time. The data quality is exceptional."
>
> — Alex Chen, CTO, PetMatch

## Availability

PetStore API is available today with plans starting at $49/month
for small stores. Enterprise pricing is available for chains.
Visit [petstore-api.example.com](https://petstore-api.example.com)
to get started with a free trial.

## About PetStore API

PetStore API is the leading API-first platform for pet retailers,
serving over 500 stores and powering 100+ third-party applications.

**Media Contact:**
press@petstore-api.example.com
```

## Step 3: FAQ (Challenges & Scope)

The FAQ challenges assumptions and clarifies scope.

### Synthesize the FAQ

```bash
visionspec synthesize faq -p petstore-api
```

### Review Generated FAQ

`docs/specs/petstore-api/gtm/faq.md`:

```markdown
# PetStore API - Frequently Asked Questions

## Customer Questions

### Q: How is this different from generic POS systems?
**A:** PetStore API is built specifically for pet retail with
features like health record tracking, breed databases, and
adoption workflow support that generic systems lack.

### Q: Can I use my existing hardware?
**A:** Yes, PetStore API works with any internet-connected device.
We provide integrations for common barcode scanners and receipt
printers.

### Q: What happens to my existing data?
**A:** We provide free data migration from common systems including
Square, Lightspeed, and spreadsheets. Our team handles the
migration with zero downtime.

## Technical Questions

### Q: What API format do you use?
**A:** RESTful API with OpenAPI 3.0 specification. We also provide
GraphQL for complex queries.

### Q: How do you handle authentication?
**A:** JWT-based authentication with OAuth 2.0 support for
third-party integrations. API keys available for server-to-server
communication.

### Q: What are the rate limits?
**A:** Free tier: 1,000 requests/hour. Paid plans: 10,000-100,000
requests/hour depending on plan.

## Scope Clarification

### Q: Do you handle payment processing?
**A:** We integrate with Stripe and Square for payment processing.
We don't process payments directly to avoid PCI compliance burden
on our platform.

### Q: Do you support international stores?
**A:** Initially US and Canada only. International expansion
planned for Year 2.

### Q: Do you handle shipping?
**A:** We integrate with ShipStation and major carriers. We don't
operate our own shipping network.

## Challenging Questions

### Q: Why would existing stores switch?
**A:** Stores switch because:
1. API-first means they can build custom solutions
2. Modern interface reduces training time
3. Real-time inventory prevents stockouts
4. Third-party integrations enable new revenue streams

### Q: What if you go out of business?
**A:** All data is exportable at any time in standard formats.
We provide 90-day notice for any service discontinuation with
full data portability.

### Q: How do you compete with free solutions?
**A:** Free solutions lack:
- Reliable uptime guarantees
- Professional support
- Pet-specific features
- Third-party ecosystem
- Data security compliance
```

## Step 4: Product Requirements (PRD)

The PRD synthesizes MRD, Press, and FAQ into actionable requirements.

### Synthesize the PRD

```bash
visionspec synthesize prd -p petstore-api
```

### Review Generated PRD

`docs/specs/petstore-api/source/prd.md`:

```markdown
# PetStore API - Product Requirements Document

## Overview

Based on the market requirements and Working Backwards artifacts,
this PRD defines the functional requirements for PetStore API v1.0.

## User Stories

### Pet Management

**US-001: Add a Pet**
As a store owner, I want to add a new pet to inventory so that
I can track and sell it.

Acceptance Criteria:
- [ ] Can enter pet details (name, species, breed, age, price)
- [ ] Can upload pet photos
- [ ] Can set availability status
- [ ] Can add health records
- [ ] Pet appears in inventory immediately

**US-002: Search Pets**
As a customer, I want to search for pets by criteria so that
I can find the perfect pet.

Acceptance Criteria:
- [ ] Can search by species (dog, cat, bird, etc.)
- [ ] Can filter by breed
- [ ] Can filter by age range
- [ ] Can filter by price range
- [ ] Can sort by relevance, price, or date added

**US-003: Update Pet Status**
As a store owner, I want to update pet status so that inventory
is accurate.

Acceptance Criteria:
- [ ] Can mark pet as sold
- [ ] Can mark pet as reserved
- [ ] Can mark pet as unavailable (illness, etc.)
- [ ] Status changes reflect immediately in API

### Inventory Management

**US-004: View Inventory**
As a store owner, I want to view current inventory so that I
know what's in stock.

Acceptance Criteria:
- [ ] Can see all products with quantities
- [ ] Can filter by category
- [ ] Can see low-stock alerts
- [ ] Can export inventory report

**US-005: Restock Alert**
As a store owner, I want automatic restock alerts so that I
don't run out of products.

Acceptance Criteria:
- [ ] Can set minimum quantity threshold per product
- [ ] Receive alert when quantity falls below threshold
- [ ] Can configure alert method (email, webhook)

### Order Processing

**US-006: Create Order**
As a customer, I want to create an order so that I can purchase
products.

Acceptance Criteria:
- [ ] Can add multiple items to order
- [ ] Can apply discount codes
- [ ] Can select payment method
- [ ] Can choose shipping or pickup
- [ ] Receive order confirmation

**US-007: Track Order**
As a customer, I want to track my order so that I know when
it will arrive.

Acceptance Criteria:
- [ ] Can view order status
- [ ] Can see shipping tracking number
- [ ] Receive status update notifications

## Functional Requirements

### FR-001: Pet CRUD API
The system shall provide RESTful endpoints for creating, reading,
updating, and deleting pet records.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/pets` | GET | List pets with filtering |
| `/pets` | POST | Create new pet |
| `/pets/{id}` | GET | Get pet details |
| `/pets/{id}` | PUT | Update pet |
| `/pets/{id}` | DELETE | Delete pet |

### FR-002: Inventory API
The system shall provide endpoints for inventory management.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/store/inventory` | GET | Get inventory summary |
| `/store/inventory/{id}` | PUT | Update quantity |

### FR-003: Order API
The system shall provide endpoints for order processing.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/store/order` | POST | Create order |
| `/store/order/{id}` | GET | Get order status |
| `/store/order/{id}` | DELETE | Cancel order |

### FR-004: User API
The system shall provide endpoints for user management.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/user` | POST | Create user |
| `/user/login` | POST | Authenticate |
| `/user/{username}` | GET | Get user profile |
| `/user/{username}` | PUT | Update user |

## Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-001 | API response time | < 200ms (p95) |
| NFR-002 | Availability | 99.9% uptime |
| NFR-003 | Concurrent users | 10,000 |
| NFR-004 | Data retention | 7 years |
| NFR-005 | Security | SOC 2 Type II |

## Release Criteria

### MVP (v1.0)
- [ ] Pet CRUD API complete
- [ ] Basic inventory management
- [ ] Order creation and tracking
- [ ] User authentication
- [ ] OpenAPI documentation
- [ ] Go SDK

### v1.1
- [ ] Advanced search
- [ ] Webhooks
- [ ] JavaScript SDK

### v1.2
- [ ] GraphQL API
- [ ] Batch operations
- [ ] Analytics dashboard
```

## Step 5: Approve Specifications

Approve the completed Working Backwards artifacts:

```bash
# Evaluate all specs
visionspec eval all -p petstore-api

# Review evaluation results
cat docs/specs/petstore-api/eval/mrd.eval.json

# Approve specs
visionspec approve mrd -p petstore-api --approver "product@example.com"
visionspec approve prd -p petstore-api --approver "product@example.com"
```

## Next Steps

- [Technical Specification](technical-spec.md) - Create TRD, TPD, IRD
- [Execution](execution.md) - Export and implement
