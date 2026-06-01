# Phase 4: Technical

The Technical phase synthesizes implementation specifications from source specs and codebase context.

## Objective

Create approved TRD, TPD, and IRD that:

- Define system architecture grounded in reality
- Plan comprehensive testing
- Specify infrastructure requirements

## Entry Criteria

- MRD, PRD, and UXD are approved
- Codebase context has been gathered (if applicable)

## Technical Spec Flow

```
MRD + PRD + UXD + Context
    ↓
TRD (Technical Requirements Document)
    ↓
TPD (Test Plan Document)
    ↓
IRD (Infrastructure Requirements Document)
```

## Workflow

### Step 1: Gather Context

Before synthesis, gather codebase context to ground decisions:

```bash
visionspec context gather -p <project>
visionspec context show -p <project>
```

Context sources include:

- Repository structure
- Existing APIs and interfaces
- Technology stack
- Deployment patterns

### Step 2: Synthesize TRD

```bash
visionspec synthesize trd -p <project>
```

The TRD should define:

**Architecture**

- System components
- Data flow
- Integration points

**APIs**

- Endpoints and contracts
- Request/response formats
- Authentication

**Data Model**

- Entities and relationships
- Storage requirements
- Migration strategy

**Non-Functional Requirements**

- Performance targets
- Security requirements
- Scalability considerations

### Step 3: Review TRD

- [ ] Architecture aligns with existing system
- [ ] APIs are consistent with current patterns
- [ ] Security is addressed
- [ ] Performance targets are realistic

### Step 4: Synthesize TPD

```bash
visionspec synthesize tpd -p <project>
```

The TPD should cover:

**Test Strategy**

- Unit test approach
- Integration test scope
- E2E test coverage

**Test Cases from PRD**

- Acceptance criteria tests
- User story validation

**Technical Tests from TRD**

- API contract tests
- Performance tests
- Security tests

**User Journey Tests from UXD**

- Happy path E2E
- Error handling
- Accessibility tests

### Step 5: Review TPD

- [ ] All PRD acceptance criteria have tests
- [ ] TRD APIs have contract tests
- [ ] UXD journeys have E2E tests
- [ ] CI/CD integration is specified

### Step 6: Synthesize IRD

```bash
visionspec synthesize ird -p <project>
```

The IRD should specify:

**Infrastructure**

- Compute resources
- Storage systems
- Network topology

**Deployment**

- Environments (dev, staging, prod)
- Deployment strategy (blue-green, canary)
- Rollback procedures

**Operations**

- Monitoring and alerting
- Logging strategy
- Incident response

### Step 7: Evaluate and Approve

```bash
visionspec eval trd -p <project>
visionspec eval tpd -p <project>
visionspec eval ird -p <project>

visionspec approve trd -p <project>
visionspec approve tpd -p <project>
visionspec approve ird -p <project>
```

## Exit Criteria

- TRD approved
- TPD approved
- IRD approved
- All evaluation scores >= 7.0

## Next Phase

→ [Phase 5: Reconciliation](05-reconciliation.md)

## Anti-Patterns

- **Greenfield syndrome**: Ignoring existing system constraints. Always gather context first.
- **Architecture astronaut**: Over-engineering for hypothetical scale. Design for current needs.
- **Test plan as afterthought**: TPD should drive implementation, not document it after.
- **Skipping IRD**: "We'll figure out deployment later" leads to production surprises.
