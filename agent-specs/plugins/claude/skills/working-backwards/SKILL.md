---
name: working-backwards
description: Working Backwards methodology for product development starting from customer outcomes
triggers: [working backwards, amazon, press release, mrd, customer-first]
---

# Working Backwards Methodology

The Working Backwards approach starts with the desired customer outcome and works backward to define what needs to be built.

## Core Principles

1. **Start with the Customer** - Define success from the customer's perspective first
2. **Press Release First** - Write the announcement before building
3. **FAQ Challenges** - Anticipate hard questions and scope issues
4. **Iterative Refinement** - Specs inform each other in cycles

## Document Flow

```
Market Problem (MRD)
       ↓
  Customer Vision (Press Release)
       ↓
  Challenges & Scope (FAQ)
       ↓
  Product Definition (PRD)
       ↓
  Technical Design (TRD/TPD/IRD)
       ↓
  Implementation (spec.md)
```

## MRD: Market Requirements Document

**Purpose**: Define the market problem, not the solution

**Key Sections**:
- Problem Statement (Core problem, Why now, Cost of inaction)
- Target Market (Primary/Secondary segments, TAM/SAM/SOM)
- Competitive Landscape (Direct/Indirect competitors, Differentiation)
- Market Requirements (Must-have, Should-have, Nice-to-have)
- Business Goals (Success metrics, Strategic alignment)
- Constraints and Assumptions
- Timeline and Milestones
- Risks with Mitigations

## Press Release

**Purpose**: Articulate the customer-facing vision

**Key Elements**:
- Headline with clear value proposition
- Problem being solved (customer perspective)
- Solution overview (features as benefits)
- Customer quotes (hypothetical but realistic)
- Call to action (availability, pricing)

**Writing Tips**:
- Write as if the product is launching today
- Use customer language, not technical jargon
- Make quotes believable and specific
- Include concrete numbers and benefits

## FAQ

**Purpose**: Challenge assumptions and clarify scope

**Question Categories**:
- **Customer Questions** - What real customers would ask
- **Technical Questions** - Developer/integration concerns
- **Scope Questions** - What's in/out, why
- **Challenging Questions** - Hard questions that test the vision

**Good FAQ Questions**:
- "Why would I switch from [competitor]?"
- "What happens if you go out of business?"
- "How do you handle [edge case]?"
- "Why don't you support [feature]?"

## PRD: Product Requirements Document

**Purpose**: Define what to build with testable criteria

**Key Elements**:
- User Stories with acceptance criteria
- Functional Requirements (numbered, traceable)
- Non-Functional Requirements (performance, security, etc.)
- Release scope (MVP vs future)

**User Story Format**:
```
As a [persona],
I want to [action],
So that [benefit].

Acceptance Criteria:
- [ ] Criterion 1
- [ ] Criterion 2
```

## Technical Specifications

### TRD (Technical Requirements)
- Architecture decisions
- Technology stack
- Component design
- API contracts
- Data models

### TPD (Test Plan)
- Test strategy
- Coverage requirements
- Test types (unit, integration, e2e)
- Performance thresholds

### IRD (Infrastructure Requirements)
- Deployment architecture
- Resource specifications
- Security requirements
- Cost estimation

## Quality Criteria

A good Working Backwards document set:
- [ ] Has clear traceability between documents
- [ ] Uses consistent terminology
- [ ] Addresses stakeholder concerns
- [ ] Is specific enough to implement
- [ ] Challenges assumptions explicitly
