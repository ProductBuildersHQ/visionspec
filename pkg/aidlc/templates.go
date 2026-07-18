// Package aidlc provides document templates for AIDLC workflows.
package aidlc

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// Template represents a document template.
type Template struct {
	// DocType is the document type this template is for.
	DocType DocType `json:"doc_type" yaml:"doc_type"`

	// Name is the template display name.
	Name string `json:"name" yaml:"name"`

	// Description describes the document's purpose.
	Description string `json:"description" yaml:"description"`

	// Content is the markdown template content.
	Content string `json:"content" yaml:"content"`

	// Sections lists the expected sections.
	Sections []TemplateSection `json:"sections" yaml:"sections"`
}

// TemplateSection describes a section in a document template.
type TemplateSection struct {
	// ID is the section identifier.
	ID string `json:"id" yaml:"id"`

	// Title is the section heading.
	Title string `json:"title" yaml:"title"`

	// Required indicates if this section is mandatory.
	Required bool `json:"required" yaml:"required"`

	// Description explains what goes in this section.
	Description string `json:"description" yaml:"description"`
}

// TemplateData contains variables for template rendering.
type TemplateData struct {
	// ProjectName is the project name.
	ProjectName string

	// Title is the document title.
	Title string

	// Author is the document author.
	Author string

	// Date is the creation date.
	Date string

	// Version is the document version.
	Version string

	// Description is a brief description.
	Description string

	// Custom contains additional custom fields.
	Custom map[string]string
}

// DefaultTemplateData returns template data with defaults.
func DefaultTemplateData(projectName string) TemplateData {
	return TemplateData{
		ProjectName: projectName,
		Title:       projectName,
		Author:      "",
		Date:        time.Now().Format("2006-01-02"),
		Version:     "1.0",
		Description: "",
		Custom:      make(map[string]string),
	}
}

// templateRegistry holds all document templates.
var templateRegistry = map[DocType]*Template{
	DocVisionDocument:     visionDocumentTemplate,
	DocRequirementsSpec:   requirementsSpecTemplate,
	DocTechnicalSpec:      technicalSpecTemplate,
	DocArchitectureSpec:   architectureSpecTemplate,
	DocImplementationPlan: implementationPlanTemplate,
	DocTestPlan:           testPlanTemplate,
	DocIntegrationPlan:    integrationPlanTemplate,
	DocSecurityReview:     securityReviewTemplate,
	DocRunbook:            runbookTemplate,
	DocMonitoringPlan:     monitoringPlanTemplate,
	DocDisasterPlan:       disasterPlanTemplate,
	DocSLODocument:        sloDocumentTemplate,
}

// GetTemplate returns the template for a document type.
func GetTemplate(docType DocType) (*Template, bool) {
	tmpl, ok := templateRegistry[docType]
	return tmpl, ok
}

// AllTemplates returns all registered templates.
func AllTemplates() map[DocType]*Template {
	return templateRegistry
}

// RenderTemplate renders a template with the given data.
func RenderTemplate(docType DocType, data TemplateData) (string, error) {
	tmpl, ok := GetTemplate(docType)
	if !ok {
		return "", fmt.Errorf("no template found for document type: %s", docType)
	}

	t, err := template.New(string(docType)).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

// --- Inception Phase Templates ---

var visionDocumentTemplate = &Template{
	DocType:     DocVisionDocument,
	Name:        "Vision Document",
	Description: "Defines the high-level product vision, goals, and success criteria.",
	Content: `---
title: "{{.Title}} - Vision Document"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Vision Document: {{.Title}}

## Executive Summary

{{.Description}}

<!-- Provide a concise overview of the product vision, key objectives, and expected outcomes. -->

## Problem Statement

### Current State
<!-- Describe the current situation, pain points, and challenges that users or the organization face. -->

### Impact
<!-- Quantify the impact of these problems (time lost, revenue impact, user frustration). -->

## Vision Statement

<!-- A clear, inspiring statement of what success looks like. Example: "Enable users to accomplish X in half the time with twice the confidence." -->

## Goals and Objectives

### Primary Goals
1. <!-- Goal 1: Specific, measurable outcome -->
2. <!-- Goal 2: Specific, measurable outcome -->
3. <!-- Goal 3: Specific, measurable outcome -->

### Success Metrics
| Metric | Current State | Target | Timeline |
|--------|---------------|--------|----------|
| <!-- Metric 1 --> | <!-- value --> | <!-- target --> | <!-- date --> |
| <!-- Metric 2 --> | <!-- value --> | <!-- target --> | <!-- date --> |

## Target Users

### User Personas
<!-- Define 2-3 primary user personas with their goals, pain points, and needs. -->

#### Persona 1: [Name]
- **Role**:
- **Goals**:
- **Pain Points**:
- **Needs**:

## Scope

### In Scope
<!-- List features and capabilities that ARE included in this initiative. -->

### Out of Scope
<!-- List features and capabilities that are explicitly NOT included. -->

## Constraints and Assumptions

### Constraints
<!-- Technical, business, regulatory, or resource constraints. -->

### Assumptions
<!-- Key assumptions that must hold true for success. -->

## Risks and Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| <!-- Risk 1 --> | High/Medium/Low | High/Medium/Low | <!-- strategy --> |

## Stakeholders

| Stakeholder | Role | Interest |
|-------------|------|----------|
| <!-- name --> | <!-- role --> | <!-- their interest in this project --> |

## Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | | | |
| Engineering Lead | | | |
| Executive Sponsor | | | |
`,
	Sections: []TemplateSection{
		{ID: "executive-summary", Title: "Executive Summary", Required: true, Description: "High-level overview"},
		{ID: "problem-statement", Title: "Problem Statement", Required: true, Description: "Current state and impact"},
		{ID: "vision-statement", Title: "Vision Statement", Required: true, Description: "What success looks like"},
		{ID: "goals-objectives", Title: "Goals and Objectives", Required: true, Description: "Measurable outcomes"},
		{ID: "target-users", Title: "Target Users", Required: true, Description: "User personas"},
		{ID: "scope", Title: "Scope", Required: true, Description: "In/out of scope"},
		{ID: "risks", Title: "Risks and Mitigations", Required: false, Description: "Risk assessment"},
	},
}

var requirementsSpecTemplate = &Template{
	DocType:     DocRequirementsSpec,
	Name:        "Requirements Specification",
	Description: "Detailed functional and non-functional requirements with acceptance criteria.",
	Content: `---
title: "{{.Title}} - Requirements Specification"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Requirements Specification: {{.Title}}

## Overview

{{.Description}}

## Functional Requirements

### FR-001: [Requirement Name]
- **Priority**: High/Medium/Low
- **Description**:
- **Rationale**:
- **Acceptance Criteria**:
  - [ ] Criterion 1
  - [ ] Criterion 2
- **Dependencies**:

### FR-002: [Requirement Name]
- **Priority**: High/Medium/Low
- **Description**:
- **Rationale**:
- **Acceptance Criteria**:
  - [ ] Criterion 1
  - [ ] Criterion 2
- **Dependencies**:

## Non-Functional Requirements

### Performance Requirements
| Requirement | Metric | Target | Measurement |
|-------------|--------|--------|-------------|
| Response Time | P95 latency | < 200ms | Load testing |
| Throughput | Requests/sec | > 1000 | Load testing |

### Scalability Requirements
<!-- How the system should scale with load. -->

### Security Requirements
<!-- Authentication, authorization, data protection requirements. -->

### Availability Requirements
<!-- Uptime targets, failover requirements. -->

### Usability Requirements
<!-- User experience standards, accessibility requirements. -->

## User Stories

### Epic: [Epic Name]

#### US-001: [Story Title]
**As a** [user type]
**I want to** [action]
**So that** [benefit]

**Acceptance Criteria**:
- Given [context], when [action], then [outcome]
- Given [context], when [action], then [outcome]

## Data Requirements

### Data Entities
<!-- Key data entities and their attributes. -->

### Data Flows
<!-- How data moves through the system. -->

### Data Retention
<!-- Data lifecycle and retention policies. -->

## Integration Requirements

### External Systems
| System | Integration Type | Purpose |
|--------|-----------------|---------|
| <!-- system --> | <!-- API/Event/File --> | <!-- purpose --> |

## Constraints

### Technical Constraints
<!-- Technology, platform, or infrastructure constraints. -->

### Business Constraints
<!-- Timeline, budget, or regulatory constraints. -->

## Traceability Matrix

| Requirement ID | User Story | Test Case | Status |
|----------------|------------|-----------|--------|
| FR-001 | US-001 | TC-001 | Pending |
`,
	Sections: []TemplateSection{
		{ID: "functional-requirements", Title: "Functional Requirements", Required: true, Description: "What the system must do"},
		{ID: "non-functional-requirements", Title: "Non-Functional Requirements", Required: true, Description: "Quality attributes"},
		{ID: "user-stories", Title: "User Stories", Required: true, Description: "User-centric requirements"},
		{ID: "data-requirements", Title: "Data Requirements", Required: false, Description: "Data structures and flows"},
		{ID: "integration-requirements", Title: "Integration Requirements", Required: false, Description: "External integrations"},
	},
}

var technicalSpecTemplate = &Template{
	DocType:     DocTechnicalSpec,
	Name:        "Technical Specification",
	Description: "Detailed technical design including APIs, data models, and algorithms.",
	Content: `---
title: "{{.Title}} - Technical Specification"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Technical Specification: {{.Title}}

## Overview

{{.Description}}

## System Context

### Context Diagram
<!-- High-level diagram showing the system and its external interactions. -->

### External Interfaces
| Interface | Type | Description |
|-----------|------|-------------|
| <!-- interface --> | <!-- REST/gRPC/Event --> | <!-- description --> |

## Technical Design

### Component Architecture

#### Component: [Component Name]
- **Purpose**:
- **Technology**:
- **Dependencies**:
- **Interfaces**:

### Data Model

#### Entity: [Entity Name]
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique identifier |
| <!-- field --> | <!-- type --> | <!-- constraints --> | <!-- description --> |

### API Specification

#### Endpoint: [Method] /api/v1/[resource]
- **Description**:
- **Authentication**: Required/Optional
- **Request**:
` + "```json" + `
{
  "field": "value"
}
` + "```" + `
- **Response (200)**:
` + "```json" + `
{
  "data": {}
}
` + "```" + `
- **Error Responses**:
  - 400: Bad Request
  - 401: Unauthorized
  - 404: Not Found

### Algorithms and Logic

#### Algorithm: [Name]
- **Purpose**:
- **Input**:
- **Output**:
- **Complexity**: O(n)
- **Pseudocode**:
` + "```" + `
<!-- pseudocode here -->
` + "```" + `

## Infrastructure

### Deployment Architecture
<!-- Describe deployment topology, cloud resources, networking. -->

### Resource Requirements
| Component | CPU | Memory | Storage | Instances |
|-----------|-----|--------|---------|-----------|
| <!-- component --> | <!-- cpu --> | <!-- mem --> | <!-- storage --> | <!-- count --> |

## Security Considerations

### Authentication
<!-- How users/services authenticate. -->

### Authorization
<!-- Access control model and policies. -->

### Data Protection
<!-- Encryption at rest and in transit. -->

## Performance Considerations

### Caching Strategy
<!-- What is cached, TTLs, invalidation. -->

### Database Optimization
<!-- Indexes, query optimization, connection pooling. -->

### Scaling Strategy
<!-- Horizontal/vertical scaling approaches. -->

## Error Handling

### Error Codes
| Code | HTTP Status | Description | Recovery |
|------|-------------|-------------|----------|
| <!-- code --> | <!-- status --> | <!-- description --> | <!-- how to recover --> |

### Retry Policy
<!-- Retry strategies for different failure types. -->

## Testing Strategy

### Unit Testing
<!-- Unit test coverage requirements and approach. -->

### Integration Testing
<!-- Integration test approach and environments. -->

### Performance Testing
<!-- Load testing requirements and benchmarks. -->

## Migration Plan

### Database Migrations
<!-- Schema migration approach, backward compatibility. -->

### Data Migration
<!-- Data migration steps if applicable. -->

## Appendix

### Glossary
| Term | Definition |
|------|------------|
| <!-- term --> | <!-- definition --> |

### References
- [Reference 1](url)
`,
	Sections: []TemplateSection{
		{ID: "system-context", Title: "System Context", Required: true, Description: "External interactions"},
		{ID: "technical-design", Title: "Technical Design", Required: true, Description: "Core design"},
		{ID: "api-specification", Title: "API Specification", Required: true, Description: "API contracts"},
		{ID: "infrastructure", Title: "Infrastructure", Required: true, Description: "Deployment details"},
		{ID: "security", Title: "Security Considerations", Required: true, Description: "Security design"},
		{ID: "testing", Title: "Testing Strategy", Required: false, Description: "Test approach"},
	},
}

var architectureSpecTemplate = &Template{
	DocType:     DocArchitectureSpec,
	Name:        "Architecture Specification",
	Description: "System architecture including components, patterns, and decisions.",
	Content: `---
title: "{{.Title}} - Architecture Specification"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Architecture Specification: {{.Title}}

## Overview

{{.Description}}

## Architecture Principles

1. **[Principle Name]**: [Description and rationale]
2. **[Principle Name]**: [Description and rationale]
3. **[Principle Name]**: [Description and rationale]

## System Architecture

### High-Level Architecture
<!-- C4 Context or similar high-level diagram. -->

### Component View
<!-- Component diagram showing major building blocks. -->

### Deployment View
<!-- Infrastructure and deployment topology. -->

## Architecture Decisions

### ADR-001: [Decision Title]
- **Status**: Proposed/Accepted/Deprecated
- **Context**: [What is the issue?]
- **Decision**: [What was decided?]
- **Rationale**: [Why this decision?]
- **Consequences**: [What are the implications?]
- **Alternatives Considered**:
  - Option A: [Description, pros/cons]
  - Option B: [Description, pros/cons]

### ADR-002: [Decision Title]
- **Status**: Proposed/Accepted/Deprecated
- **Context**:
- **Decision**:
- **Rationale**:
- **Consequences**:

## Component Catalog

### Component: [Name]
- **Type**: Service/Library/Infrastructure
- **Responsibility**:
- **Technology Stack**:
- **Dependencies**:
- **Interfaces**:
- **Data Stores**:

## Integration Architecture

### Internal Communication
| From | To | Protocol | Pattern |
|------|-----|----------|---------|
| <!-- service --> | <!-- service --> | <!-- REST/gRPC/Events --> | <!-- sync/async --> |

### External Integrations
<!-- Third-party services and APIs. -->

## Data Architecture

### Data Stores
| Store | Type | Purpose | Data Retention |
|-------|------|---------|----------------|
| <!-- store --> | <!-- PostgreSQL/Redis/S3 --> | <!-- purpose --> | <!-- retention --> |

### Data Flow
<!-- How data flows through the system. -->

## Security Architecture

### Trust Boundaries
<!-- Security zones and trust relationships. -->

### Identity and Access
<!-- Authentication and authorization architecture. -->

## Observability Architecture

### Logging Strategy
<!-- Logging approach, levels, retention. -->

### Metrics and Monitoring
<!-- Key metrics, alerting thresholds. -->

### Tracing
<!-- Distributed tracing approach. -->

## Scalability and Resilience

### Scaling Strategy
<!-- How the system scales horizontally/vertically. -->

### Fault Tolerance
<!-- How the system handles failures. -->

### Disaster Recovery
<!-- DR strategy and RTO/RPO. -->

## Technology Stack

| Category | Technology | Version | Rationale |
|----------|------------|---------|-----------|
| Language | <!-- lang --> | <!-- version --> | <!-- why --> |
| Framework | <!-- framework --> | <!-- version --> | <!-- why --> |
| Database | <!-- db --> | <!-- version --> | <!-- why --> |
| Messaging | <!-- mq --> | <!-- version --> | <!-- why --> |

## Constraints and Trade-offs

### Constraints
<!-- Architectural constraints from business or technical requirements. -->

### Trade-offs
<!-- Key trade-offs made and their rationale. -->

## Future Considerations

<!-- Planned architectural evolution and roadmap. -->
`,
	Sections: []TemplateSection{
		{ID: "principles", Title: "Architecture Principles", Required: true, Description: "Guiding principles"},
		{ID: "system-architecture", Title: "System Architecture", Required: true, Description: "Architecture views"},
		{ID: "decisions", Title: "Architecture Decisions", Required: true, Description: "ADRs"},
		{ID: "components", Title: "Component Catalog", Required: true, Description: "Component details"},
		{ID: "data-architecture", Title: "Data Architecture", Required: true, Description: "Data design"},
		{ID: "security", Title: "Security Architecture", Required: true, Description: "Security design"},
	},
}

// --- Construction Phase Templates ---

var implementationPlanTemplate = &Template{
	DocType:     DocImplementationPlan,
	Name:        "Implementation Plan",
	Description: "Detailed implementation timeline, tasks, and resource allocation.",
	Content: `---
title: "{{.Title}} - Implementation Plan"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Implementation Plan: {{.Title}}

## Overview

{{.Description}}

## Implementation Approach

### Development Methodology
<!-- Agile/Scrum/Kanban approach, sprint length, ceremonies. -->

### Release Strategy
<!-- Phased release, feature flags, rollout approach. -->

## Milestones

| Milestone | Description | Target Date | Dependencies |
|-----------|-------------|-------------|--------------|
| M1: Foundation | Core infrastructure | <!-- date --> | None |
| M2: MVP | Core features | <!-- date --> | M1 |
| M3: Beta | Full feature set | <!-- date --> | M2 |
| M4: GA | Production release | <!-- date --> | M3 |

## Work Breakdown

### Phase 1: Foundation
| Task | Description | Owner | Estimate | Dependencies |
|------|-------------|-------|----------|--------------|
| 1.1 | Set up development environment | <!-- owner --> | <!-- days --> | None |
| 1.2 | Infrastructure provisioning | <!-- owner --> | <!-- days --> | 1.1 |
| 1.3 | CI/CD pipeline setup | <!-- owner --> | <!-- days --> | 1.2 |

### Phase 2: Core Development
| Task | Description | Owner | Estimate | Dependencies |
|------|-------------|-------|----------|--------------|
| 2.1 | <!-- task --> | <!-- owner --> | <!-- days --> | <!-- deps --> |

### Phase 3: Testing and Hardening
| Task | Description | Owner | Estimate | Dependencies |
|------|-------------|-------|----------|--------------|
| 3.1 | <!-- task --> | <!-- owner --> | <!-- days --> | <!-- deps --> |

## Resource Allocation

### Team Structure
| Role | Name | Allocation | Responsibilities |
|------|------|------------|------------------|
| Tech Lead | <!-- name --> | 100% | Architecture, code review |
| Backend Engineer | <!-- name --> | 100% | API development |
| Frontend Engineer | <!-- name --> | 100% | UI development |
| QA Engineer | <!-- name --> | 50% | Testing |

### External Dependencies
| Dependency | Provider | Status | Contact |
|------------|----------|--------|---------|
| <!-- service --> | <!-- provider --> | <!-- status --> | <!-- contact --> |

## Risk Management

| Risk | Probability | Impact | Mitigation | Contingency |
|------|-------------|--------|------------|-------------|
| <!-- risk --> | H/M/L | H/M/L | <!-- strategy --> | <!-- backup plan --> |

## Quality Gates

### Code Quality
- [ ] Code review approval (2 reviewers)
- [ ] Unit test coverage > 80%
- [ ] No critical linting errors
- [ ] Documentation updated

### Integration Quality
- [ ] Integration tests passing
- [ ] API contract tests passing
- [ ] Performance benchmarks met

### Release Quality
- [ ] Security scan passed
- [ ] Load testing completed
- [ ] Runbook reviewed
- [ ] Rollback tested

## Communication Plan

### Status Updates
- Daily standups: [time, channel]
- Weekly status: [day, format]
- Milestone reviews: [format]

### Escalation Path
1. Team Lead
2. Engineering Manager
3. Director of Engineering

## Dependencies Timeline

` + "```mermaid" + `
gantt
    title Implementation Timeline
    dateFormat  YYYY-MM-DD
    section Foundation
    Dev Environment    :a1, 2024-01-01, 5d
    Infrastructure     :a2, after a1, 10d
    CI/CD Setup        :a3, after a1, 5d
    section Development
    Core APIs          :b1, after a2, 15d
    Frontend           :b2, after b1, 10d
    section Testing
    Integration Tests  :c1, after b2, 5d
    Load Testing       :c2, after c1, 3d
` + "```" + `

## Appendix

### Definitions
| Term | Definition |
|------|------------|
| <!-- term --> | <!-- definition --> |
`,
	Sections: []TemplateSection{
		{ID: "approach", Title: "Implementation Approach", Required: true, Description: "Methodology and strategy"},
		{ID: "milestones", Title: "Milestones", Required: true, Description: "Key deliverables"},
		{ID: "work-breakdown", Title: "Work Breakdown", Required: true, Description: "Detailed tasks"},
		{ID: "resources", Title: "Resource Allocation", Required: true, Description: "Team and dependencies"},
		{ID: "risks", Title: "Risk Management", Required: true, Description: "Risk identification"},
		{ID: "quality-gates", Title: "Quality Gates", Required: true, Description: "Acceptance criteria"},
	},
}

var testPlanTemplate = &Template{
	DocType:     DocTestPlan,
	Name:        "Test Plan",
	Description: "Comprehensive testing strategy including test cases and coverage requirements.",
	Content: `---
title: "{{.Title}} - Test Plan"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Test Plan: {{.Title}}

## Overview

{{.Description}}

## Test Strategy

### Testing Levels
| Level | Scope | Ownership | Automation |
|-------|-------|-----------|------------|
| Unit | Individual components | Developers | 100% |
| Integration | Component interactions | Developers/QA | 80% |
| System | End-to-end flows | QA | 60% |
| Acceptance | Business requirements | QA/Product | 40% |

### Testing Types
- [ ] Functional Testing
- [ ] Performance Testing
- [ ] Security Testing
- [ ] Accessibility Testing
- [ ] Usability Testing

## Test Environment

### Environments
| Environment | Purpose | Data | URL |
|-------------|---------|------|-----|
| Dev | Development testing | Synthetic | <!-- url --> |
| QA | QA testing | Synthetic | <!-- url --> |
| Staging | Pre-production | Anonymized prod | <!-- url --> |
| Prod | Production | Live | <!-- url --> |

### Test Data
<!-- Test data requirements, generation, and management. -->

## Test Cases

### TC-001: [Test Case Name]
- **Priority**: High/Medium/Low
- **Type**: Functional/Performance/Security
- **Preconditions**:
- **Steps**:
  1. Step 1
  2. Step 2
  3. Step 3
- **Expected Result**:
- **Postconditions**:

### TC-002: [Test Case Name]
- **Priority**: High/Medium/Low
- **Type**: Functional/Performance/Security
- **Preconditions**:
- **Steps**:
  1. Step 1
- **Expected Result**:
- **Postconditions**:

## Coverage Requirements

### Code Coverage
| Component | Target | Measurement |
|-----------|--------|-------------|
| Core Services | 80% | Line coverage |
| API Layer | 90% | Line coverage |
| Utilities | 70% | Line coverage |

### Requirement Coverage
| Requirement | Test Cases | Status |
|-------------|------------|--------|
| FR-001 | TC-001, TC-002 | Pending |
| FR-002 | TC-003 | Pending |

## Performance Testing

### Load Test Scenarios
| Scenario | Users | Duration | Target |
|----------|-------|----------|--------|
| Normal Load | 100 | 30min | P95 < 200ms |
| Peak Load | 500 | 15min | P95 < 500ms |
| Stress Test | 1000 | 10min | No errors |

### Performance Acceptance Criteria
- Response time P95 < 200ms under normal load
- Error rate < 0.1%
- Throughput > 1000 req/sec

## Security Testing

### Security Test Cases
| Test | Description | Tool | Status |
|------|-------------|------|--------|
| OWASP Top 10 | Common vulnerabilities | OWASP ZAP | Pending |
| Authentication | Auth bypass attempts | Manual | Pending |
| Authorization | Access control testing | Manual | Pending |
| Injection | SQL/NoSQL/Command injection | SQLMap | Pending |

## Defect Management

### Severity Levels
| Severity | Description | SLA |
|----------|-------------|-----|
| Critical | System unusable | 4 hours |
| High | Major feature broken | 1 day |
| Medium | Minor feature affected | 3 days |
| Low | Cosmetic issue | 1 week |

### Defect Workflow
1. Open -> In Progress -> Fixed -> Verified -> Closed

## Test Schedule

| Phase | Start | End | Deliverables |
|-------|-------|-----|--------------|
| Test Planning | <!-- date --> | <!-- date --> | Test plan approved |
| Test Development | <!-- date --> | <!-- date --> | Test cases ready |
| Test Execution | <!-- date --> | <!-- date --> | Test results |
| Regression | <!-- date --> | <!-- date --> | Regression report |

## Exit Criteria

### Release Criteria
- [ ] All critical/high defects resolved
- [ ] Test coverage targets met
- [ ] Performance benchmarks passed
- [ ] Security scan passed
- [ ] No regression in existing functionality

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Test environment instability | Dedicated test environment with monitoring |
| Insufficient test data | Test data generation scripts |
| Resource constraints | Prioritize critical test cases |
`,
	Sections: []TemplateSection{
		{ID: "strategy", Title: "Test Strategy", Required: true, Description: "Overall testing approach"},
		{ID: "environment", Title: "Test Environment", Required: true, Description: "Environment setup"},
		{ID: "test-cases", Title: "Test Cases", Required: true, Description: "Detailed test cases"},
		{ID: "coverage", Title: "Coverage Requirements", Required: true, Description: "Coverage targets"},
		{ID: "performance", Title: "Performance Testing", Required: true, Description: "Load testing"},
		{ID: "security", Title: "Security Testing", Required: true, Description: "Security tests"},
	},
}

var integrationPlanTemplate = &Template{
	DocType:     DocIntegrationPlan,
	Name:        "Integration Plan",
	Description: "Plan for integrating with external systems and services.",
	Content: `---
title: "{{.Title}} - Integration Plan"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Integration Plan: {{.Title}}

## Overview

{{.Description}}

## Integration Landscape

### Systems Map
<!-- Diagram showing all integrated systems and data flows. -->

### Integration Summary
| System | Direction | Protocol | Frequency | Criticality |
|--------|-----------|----------|-----------|-------------|
| <!-- system --> | Inbound/Outbound/Bidirectional | REST/gRPC/Events | Real-time/Batch | High/Medium/Low |

## Integration Details

### Integration: [System Name]

#### Overview
- **Purpose**:
- **Owner**:
- **Contact**:

#### Technical Details
- **Protocol**: REST/gRPC/GraphQL/Events
- **Authentication**: OAuth/API Key/mTLS
- **Endpoint**:
- **Rate Limits**:

#### Data Contract
**Request**:
` + "```json" + `
{
  "field": "value"
}
` + "```" + `

**Response**:
` + "```json" + `
{
  "data": {}
}
` + "```" + `

#### Error Handling
| Error | Recovery | Retry |
|-------|----------|-------|
| Timeout | Exponential backoff | 3 attempts |
| 429 Rate Limited | Wait and retry | With delay |
| 5xx Server Error | Alert and retry | 3 attempts |

#### Testing
- **Test Environment**:
- **Test Credentials**:
- **Test Data**:

## Event Architecture

### Event Catalog
| Event | Publisher | Consumers | Schema |
|-------|-----------|-----------|--------|
| <!-- event --> | <!-- service --> | <!-- services --> | <!-- schema link --> |

### Event Schema: [EventName]
` + "```json" + `
{
  "type": "event.type",
  "version": "1.0",
  "data": {}
}
` + "```" + `

## Dependency Management

### External Dependencies
| Dependency | SLA | Fallback | Owner |
|------------|-----|----------|-------|
| <!-- service --> | 99.9% | <!-- fallback strategy --> | <!-- team --> |

### Circuit Breaker Configuration
| Service | Threshold | Timeout | Reset |
|---------|-----------|---------|-------|
| <!-- service --> | 50% failures | 30s | 60s |

## Migration Plan

### Phase 1: Preparation
- [ ] API credentials provisioned
- [ ] Test environment configured
- [ ] Integration code developed
- [ ] Contract tests passing

### Phase 2: Testing
- [ ] Integration testing complete
- [ ] Performance testing complete
- [ ] Failover testing complete

### Phase 3: Rollout
- [ ] Canary deployment
- [ ] Gradual traffic shift
- [ ] Full production traffic

## Monitoring and Alerting

### Integration Health
| Metric | Threshold | Alert |
|--------|-----------|-------|
| Latency P95 | > 500ms | Warning |
| Error Rate | > 1% | Critical |
| Availability | < 99.9% | Critical |

### Dashboards
- Integration Overview: [link]
- Per-System Detail: [link]

## Runbook

### Common Issues
| Issue | Symptoms | Resolution |
|-------|----------|------------|
| Auth failure | 401 errors | Refresh credentials |
| Rate limiting | 429 errors | Implement backoff |
| Timeout | Slow responses | Check network/service |

## Appendix

### API Documentation Links
- [System A API Docs](url)
- [System B API Docs](url)

### Contact List
| System | Team | Email | Slack |
|--------|------|-------|-------|
| <!-- system --> | <!-- team --> | <!-- email --> | <!-- channel --> |
`,
	Sections: []TemplateSection{
		{ID: "landscape", Title: "Integration Landscape", Required: true, Description: "Systems overview"},
		{ID: "details", Title: "Integration Details", Required: true, Description: "Per-system details"},
		{ID: "events", Title: "Event Architecture", Required: false, Description: "Event-driven integrations"},
		{ID: "dependencies", Title: "Dependency Management", Required: true, Description: "External dependencies"},
		{ID: "migration", Title: "Migration Plan", Required: true, Description: "Rollout strategy"},
		{ID: "monitoring", Title: "Monitoring and Alerting", Required: true, Description: "Health checks"},
	},
}

var securityReviewTemplate = &Template{
	DocType:     DocSecurityReview,
	Name:        "Security Review",
	Description: "Security assessment including threat modeling and controls.",
	Content: `---
title: "{{.Title}} - Security Review"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Security Review: {{.Title}}

## Overview

{{.Description}}

## Security Classification

| Classification | Value |
|----------------|-------|
| Data Sensitivity | Public/Internal/Confidential/Restricted |
| System Criticality | Low/Medium/High/Critical |
| Compliance Scope | SOC2/GDPR/HIPAA/PCI-DSS |

## Threat Model

### Assets
| Asset | Sensitivity | Location | Protection |
|-------|-------------|----------|------------|
| User PII | High | Database | Encrypted at rest |
| API Keys | Critical | Vault | HSM-backed |
| Session Data | Medium | Redis | Encrypted |

### Threat Actors
| Actor | Motivation | Capability | Likelihood |
|-------|------------|------------|------------|
| External Attacker | Financial | High | Medium |
| Malicious Insider | Espionage | Medium | Low |
| Automated Bot | Disruption | Low | High |

### STRIDE Analysis

#### Spoofing
- **Threat**: Attacker impersonates legitimate user
- **Mitigation**: MFA, session management
- **Status**: Mitigated

#### Tampering
- **Threat**: Modification of data in transit/at rest
- **Mitigation**: TLS, integrity checks
- **Status**: Mitigated

#### Repudiation
- **Threat**: User denies actions
- **Mitigation**: Audit logging
- **Status**: Mitigated

#### Information Disclosure
- **Threat**: Unauthorized data access
- **Mitigation**: Encryption, access control
- **Status**: Mitigated

#### Denial of Service
- **Threat**: Service unavailability
- **Mitigation**: Rate limiting, WAF
- **Status**: Mitigated

#### Elevation of Privilege
- **Threat**: Unauthorized access escalation
- **Mitigation**: RBAC, least privilege
- **Status**: Mitigated

## Security Controls

### Authentication
- [ ] MFA enforced for privileged access
- [ ] Strong password policy
- [ ] Session timeout configured
- [ ] Brute force protection

### Authorization
- [ ] RBAC implemented
- [ ] Least privilege principle
- [ ] Resource-level permissions
- [ ] Authorization audit logging

### Data Protection
- [ ] TLS 1.3 for transit
- [ ] AES-256 for rest
- [ ] Key rotation policy
- [ ] Secrets in vault

### Network Security
- [ ] WAF configured
- [ ] DDoS protection
- [ ] Network segmentation
- [ ] Firewall rules

### Application Security
- [ ] Input validation
- [ ] Output encoding
- [ ] CSRF protection
- [ ] Security headers

## Vulnerability Assessment

### SAST Results
| Finding | Severity | Status | Remediation |
|---------|----------|--------|-------------|
| <!-- finding --> | Critical/High/Medium/Low | Open/Fixed | <!-- action --> |

### DAST Results
| Finding | Severity | Status | Remediation |
|---------|----------|--------|-------------|
| <!-- finding --> | Critical/High/Medium/Low | Open/Fixed | <!-- action --> |

### Dependency Scan
| Dependency | CVE | Severity | Action |
|------------|-----|----------|--------|
| <!-- dep --> | <!-- CVE --> | <!-- severity --> | Upgrade/Accept |

## Compliance Checklist

### SOC 2 Controls
- [ ] Access control policies documented
- [ ] Security awareness training
- [ ] Incident response plan
- [ ] Change management process

### GDPR Requirements
- [ ] Data processing agreements
- [ ] Privacy policy updated
- [ ] Consent management
- [ ] Data subject rights

## Incident Response

### Response Team
| Role | Contact | Responsibilities |
|------|---------|------------------|
| Security Lead | <!-- contact --> | Coordinate response |
| Engineering Lead | <!-- contact --> | Technical remediation |
| Legal | <!-- contact --> | Regulatory notification |

### Response Procedures
1. Identify and contain
2. Assess impact
3. Remediate
4. Notify stakeholders
5. Post-mortem

## Security Testing

### Penetration Testing
- **Scope**: [Define scope]
- **Schedule**: [Testing dates]
- **Provider**: [Internal/External]

### Bug Bounty
- **Program**: [Link to program]
- **Scope**: [Covered assets]

## Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Security Engineer | | | |
| Security Lead | | | |
| Engineering Lead | | | |
`,
	Sections: []TemplateSection{
		{ID: "classification", Title: "Security Classification", Required: true, Description: "Data and system classification"},
		{ID: "threat-model", Title: "Threat Model", Required: true, Description: "Threats and mitigations"},
		{ID: "controls", Title: "Security Controls", Required: true, Description: "Security measures"},
		{ID: "vulnerabilities", Title: "Vulnerability Assessment", Required: true, Description: "Scan results"},
		{ID: "compliance", Title: "Compliance Checklist", Required: false, Description: "Regulatory compliance"},
		{ID: "incident-response", Title: "Incident Response", Required: true, Description: "Response procedures"},
	},
}

// --- Operations Phase Templates ---

var runbookTemplate = &Template{
	DocType:     DocRunbook,
	Name:        "Runbook",
	Description: "Operational procedures for deployment, maintenance, and incident response.",
	Content: `---
title: "{{.Title}} - Runbook"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Runbook: {{.Title}}

## Overview

{{.Description}}

## Service Information

| Property | Value |
|----------|-------|
| Service Name | {{.Title}} |
| Repository | [link] |
| CI/CD Pipeline | [link] |
| Documentation | [link] |
| On-Call Rotation | [link] |

## Architecture Summary

<!-- High-level diagram of the service and its dependencies. -->

### Components
| Component | Purpose | Location |
|-----------|---------|----------|
| <!-- component --> | <!-- purpose --> | <!-- cluster/region --> |

### Dependencies
| Dependency | Type | Criticality | Fallback |
|------------|------|-------------|----------|
| <!-- service --> | Internal/External | Critical/Important/Optional | <!-- strategy --> |

## Deployment

### Prerequisites
- [ ] Access to deployment environment
- [ ] Approved deployment ticket
- [ ] Monitoring dashboards open
- [ ] Rollback plan reviewed

### Deployment Steps
1. **Pre-deployment checks**
` + "```bash" + `
# Verify service health
kubectl get pods -n production
` + "```" + `

2. **Deploy new version**
` + "```bash" + `
# Deploy using CI/CD or kubectl
kubectl set image deployment/service service=image:tag -n production
` + "```" + `

3. **Verify deployment**
` + "```bash" + `
# Check rollout status
kubectl rollout status deployment/service -n production
` + "```" + `

4. **Post-deployment validation**
- [ ] Health check endpoints responding
- [ ] Metrics flowing to dashboards
- [ ] No error spike in logs

### Rollback Procedure
` + "```bash" + `
# Rollback to previous version
kubectl rollout undo deployment/service -n production
` + "```" + `

## Common Operations

### Scaling

#### Manual Scaling
` + "```bash" + `
# Scale to N replicas
kubectl scale deployment/service --replicas=N -n production
` + "```" + `

#### Auto-scaling Configuration
<!-- HPA settings and tuning parameters. -->

### Configuration Changes
1. Update configuration in [config repo]
2. Create PR and get approval
3. Merge triggers automated deployment
4. Verify configuration applied

### Log Access
` + "```bash" + `
# View recent logs
kubectl logs -f deployment/service -n production

# Search logs in log aggregator
<!-- logging platform query -->
` + "```" + `

## Incident Response

### Severity Levels
| Severity | Definition | Response Time | Escalation |
|----------|------------|---------------|------------|
| P1 | Service down | 5 min | Immediate |
| P2 | Major degradation | 15 min | 30 min |
| P3 | Minor impact | 1 hour | 4 hours |
| P4 | No user impact | 1 day | N/A |

### Common Issues

#### High Error Rate
**Symptoms**: Error rate > 1% in dashboards
**Investigation**:
1. Check error logs for patterns
2. Identify affected endpoints
3. Check recent deployments

**Resolution**:
- If deployment related: Rollback
- If dependency related: Check dependency health
- If traffic related: Scale up

#### High Latency
**Symptoms**: P95 latency > 500ms
**Investigation**:
1. Check database query latency
2. Check external dependency latency
3. Check CPU/memory utilization

**Resolution**:
- If database: Optimize queries or add read replicas
- If dependency: Implement caching or circuit breaker
- If resource: Scale up instances

#### Out of Memory
**Symptoms**: OOM kills, pods restarting
**Investigation**:
1. Check memory metrics
2. Identify memory leaks with profiler
3. Review recent code changes

**Resolution**:
- Immediate: Increase memory limits
- Long-term: Fix memory leak

### Escalation Contacts

| Level | Team | Contact | When |
|-------|------|---------|------|
| L1 | On-call | PagerDuty | First response |
| L2 | Engineering Lead | <!-- contact --> | 15 min no resolution |
| L3 | Director | <!-- contact --> | P1 > 30 min |

## Monitoring

### Dashboards
- [Service Overview](url)
- [Error Analysis](url)
- [Performance Metrics](url)

### Key Metrics
| Metric | Normal | Warning | Critical |
|--------|--------|---------|----------|
| Error Rate | < 0.1% | > 0.5% | > 1% |
| Latency P95 | < 200ms | > 300ms | > 500ms |
| CPU Usage | < 60% | > 75% | > 90% |
| Memory Usage | < 70% | > 80% | > 90% |

### Alerts
| Alert | Condition | Severity | Runbook Section |
|-------|-----------|----------|-----------------|
| HighErrorRate | error_rate > 1% for 5m | P2 | High Error Rate |
| HighLatency | p95_latency > 500ms for 5m | P2 | High Latency |

## Maintenance

### Scheduled Maintenance
<!-- Regular maintenance tasks and schedules. -->

### Database Operations
` + "```bash" + `
# Connect to database (read-only)
# WARNING: Production access requires approval
` + "```" + `

## Appendix

### Useful Commands
` + "```bash" + `
# Quick health check
curl -s http://service/health | jq

# Get pod details
kubectl describe pod <pod-name> -n production
` + "```" + `

### Links
- [Architecture Diagram](url)
- [API Documentation](url)
- [On-call Playbook](url)
`,
	Sections: []TemplateSection{
		{ID: "service-info", Title: "Service Information", Required: true, Description: "Basic service details"},
		{ID: "deployment", Title: "Deployment", Required: true, Description: "Deployment procedures"},
		{ID: "operations", Title: "Common Operations", Required: true, Description: "Day-to-day tasks"},
		{ID: "incident-response", Title: "Incident Response", Required: true, Description: "Troubleshooting"},
		{ID: "monitoring", Title: "Monitoring", Required: true, Description: "Observability"},
		{ID: "maintenance", Title: "Maintenance", Required: false, Description: "Regular maintenance"},
	},
}

var monitoringPlanTemplate = &Template{
	DocType:     DocMonitoringPlan,
	Name:        "Monitoring Plan",
	Description: "Observability strategy including metrics, logs, and alerting.",
	Content: `---
title: "{{.Title}} - Monitoring Plan"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Monitoring Plan: {{.Title}}

## Overview

{{.Description}}

## Observability Strategy

### Three Pillars
| Pillar | Tool | Retention | Purpose |
|--------|------|-----------|---------|
| Metrics | Prometheus/Datadog | 30 days | Performance tracking |
| Logs | Loki/Elasticsearch | 14 days | Debugging |
| Traces | Jaeger/Tempo | 7 days | Request flows |

## Metrics

### Golden Signals
| Signal | Metric | Query | Dashboard |
|--------|--------|-------|-----------|
| Latency | http_request_duration_seconds | histogram_quantile(0.95, ...) | [link] |
| Traffic | http_requests_total | rate(...[5m]) | [link] |
| Errors | http_requests_total{status=~"5.."} | rate(...[5m]) | [link] |
| Saturation | container_memory_usage_bytes | <!-- query --> | [link] |

### Custom Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| <!-- metric_name --> | Counter/Gauge/Histogram | <!-- labels --> | <!-- description --> |

### Dashboard Catalog
| Dashboard | Purpose | Owner | Link |
|-----------|---------|-------|------|
| Service Overview | High-level health | SRE | [link] |
| Performance Deep Dive | Latency analysis | Engineering | [link] |
| Business Metrics | KPI tracking | Product | [link] |

## Logging

### Log Levels
| Level | Usage | Example |
|-------|-------|---------|
| ERROR | Failures requiring attention | Database connection failed |
| WARN | Unusual but handled conditions | Retry attempt |
| INFO | Significant events | Request processed |
| DEBUG | Development/troubleshooting | Query parameters |

### Structured Log Format
` + "```json" + `
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "service": "{{.Title}}",
  "trace_id": "abc123",
  "message": "Request processed",
  "duration_ms": 42,
  "user_id": "user123"
}
` + "```" + `

### Log Retention
| Log Type | Retention | Storage | Cost |
|----------|-----------|---------|------|
| Application | 14 days | Hot | $X/GB |
| Access | 30 days | Warm | $X/GB |
| Audit | 1 year | Cold | $X/GB |

## Tracing

### Instrumentation
| Component | Library | Sampling |
|-----------|---------|----------|
| HTTP Server | OpenTelemetry | 10% |
| Database | OpenTelemetry | 100% |
| External API | OpenTelemetry | 100% |

### Trace Context
- Propagation: W3C Trace Context
- Span attributes: user_id, tenant_id, request_id

## Alerting

### Alert Rules
| Alert | Condition | Severity | Runbook |
|-------|-----------|----------|---------|
| HighErrorRate | error_rate > 1% for 5m | P2 | [link] |
| HighLatency | p95_latency > 500ms for 5m | P2 | [link] |
| ServiceDown | up == 0 for 1m | P1 | [link] |
| HighMemory | memory > 90% for 10m | P3 | [link] |

### Alert Routing
| Severity | Channel | Recipients | Timing |
|----------|---------|------------|--------|
| P1 | PagerDuty | On-call | 24/7 |
| P2 | Slack + PagerDuty | Team | Business hours |
| P3 | Slack | Team | Best effort |
| P4 | Email | Team | Weekly digest |

### On-Call Schedule
- Rotation: Weekly
- Primary: [link to schedule]
- Secondary: [link to schedule]

## SLIs and SLOs

### Service Level Indicators
| SLI | Definition | Measurement |
|-----|------------|-------------|
| Availability | % of successful requests | success_rate |
| Latency | P95 request latency | histogram |
| Correctness | % of valid responses | validation checks |

### Service Level Objectives
| SLO | Target | Measurement Window |
|-----|--------|-------------------|
| Availability | 99.9% | 30 days |
| Latency P95 | < 200ms | 30 days |
| Error Budget | 0.1% / month | Monthly |

## Health Checks

### Endpoints
| Endpoint | Purpose | Interval |
|----------|---------|----------|
| /health | Liveness | 10s |
| /ready | Readiness | 5s |
| /metrics | Prometheus scrape | 15s |

### Synthetic Monitoring
| Check | Interval | Locations | Alert |
|-------|----------|-----------|-------|
| Homepage | 1 min | 3 regions | Yes |
| API Health | 30 sec | 3 regions | Yes |
| Critical Flow | 5 min | 1 region | Yes |

## Capacity Planning

### Growth Projections
| Resource | Current | +6 months | +12 months |
|----------|---------|-----------|------------|
| Requests/sec | <!-- value --> | <!-- projection --> | <!-- projection --> |
| Storage | <!-- value --> | <!-- projection --> | <!-- projection --> |
| Compute | <!-- value --> | <!-- projection --> | <!-- projection --> |

### Scaling Triggers
| Metric | Scale Out | Scale In |
|--------|-----------|----------|
| CPU | > 70% avg | < 30% avg |
| Memory | > 80% avg | < 40% avg |
| Queue Length | > 100 | < 10 |

## Runbook Integration

### Alert to Runbook Mapping
| Alert | Runbook Section | Escalation |
|-------|-----------------|------------|
| HighErrorRate | [Runbook: High Error Rate] | L1 -> L2 (15m) |
| ServiceDown | [Runbook: Service Down] | Immediate L2 |
`,
	Sections: []TemplateSection{
		{ID: "strategy", Title: "Observability Strategy", Required: true, Description: "Overall approach"},
		{ID: "metrics", Title: "Metrics", Required: true, Description: "Key metrics"},
		{ID: "logging", Title: "Logging", Required: true, Description: "Log management"},
		{ID: "alerting", Title: "Alerting", Required: true, Description: "Alert configuration"},
		{ID: "slos", Title: "SLIs and SLOs", Required: true, Description: "Service levels"},
		{ID: "capacity", Title: "Capacity Planning", Required: false, Description: "Growth planning"},
	},
}

var disasterPlanTemplate = &Template{
	DocType:     DocDisasterPlan,
	Name:        "Disaster Recovery Plan",
	Description: "Business continuity and disaster recovery procedures.",
	Content: `---
title: "{{.Title}} - Disaster Recovery Plan"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# Disaster Recovery Plan: {{.Title}}

## Overview

{{.Description}}

## Recovery Objectives

| Objective | Target | Justification |
|-----------|--------|---------------|
| RTO (Recovery Time Objective) | <!-- hours --> | <!-- business reason --> |
| RPO (Recovery Point Objective) | <!-- hours --> | <!-- data loss tolerance --> |
| MTTR (Mean Time To Recovery) | <!-- hours --> | <!-- operational capability --> |

## Business Impact Analysis

### Critical Services
| Service | Priority | RTO | RPO | Impact/Hour |
|---------|----------|-----|-----|-------------|
| <!-- service --> | P1 | <!-- time --> | <!-- time --> | $<!-- amount --> |

### Dependencies
| Dependency | Type | Recovery Priority | Alternative |
|------------|------|-------------------|-------------|
| <!-- dep --> | Internal/External | P1/P2/P3 | <!-- fallback --> |

## Disaster Scenarios

### Scenario 1: Data Center Outage
- **Probability**: Low
- **Impact**: Critical
- **Recovery Strategy**: Failover to DR region
- **RTO**: 4 hours

### Scenario 2: Database Corruption
- **Probability**: Medium
- **Impact**: High
- **Recovery Strategy**: Point-in-time recovery
- **RTO**: 2 hours

### Scenario 3: Ransomware Attack
- **Probability**: Medium
- **Impact**: Critical
- **Recovery Strategy**: Clean room recovery
- **RTO**: 24 hours

## Recovery Procedures

### Regional Failover

#### Prerequisites
- [ ] DR region validated
- [ ] DNS failover configured
- [ ] Data replication verified

#### Procedure
1. **Assess situation** (5 min)
   - Confirm primary region is unavailable
   - Verify DR region health

2. **Initiate failover** (15 min)
` + "```bash" + `
# Update DNS to point to DR region
# Or trigger automated failover
` + "```" + `

3. **Verify services** (30 min)
   - [ ] All services healthy
   - [ ] Data integrity verified
   - [ ] External integrations working

4. **Notify stakeholders** (10 min)
   - Engineering teams
   - Business stakeholders
   - Customers (if needed)

### Database Recovery

#### Point-in-Time Recovery
` + "```bash" + `
# Restore database to specific timestamp
# Steps depend on database platform
` + "```" + `

#### From Backup
` + "```bash" + `
# Restore from latest backup
# Verify data integrity
` + "```" + `

### Clean Room Recovery
For ransomware or security incidents:
1. Isolate affected systems
2. Provision clean infrastructure
3. Restore from verified clean backups
4. Rebuild affected services
5. Security validation before reconnection

## Backup Strategy

### Backup Schedule
| Data | Frequency | Retention | Location | Encryption |
|------|-----------|-----------|----------|------------|
| Database | Hourly | 7 days | S3 (different region) | AES-256 |
| Files | Daily | 30 days | S3 (different region) | AES-256 |
| Config | On change | 90 days | Git + S3 | Yes |

### Backup Verification
| Test | Frequency | Last Tested | Next Test |
|------|-----------|-------------|-----------|
| Backup completion | Daily (automated) | <!-- date --> | Daily |
| Restore to test env | Monthly | <!-- date --> | <!-- date --> |
| Full DR simulation | Quarterly | <!-- date --> | <!-- date --> |

## Communication Plan

### Notification Chain
| Event | Notified | Method | Template |
|-------|----------|--------|----------|
| DR declared | Exec team | Phone | [link] |
| Recovery started | Engineering | Slack | [link] |
| Recovery complete | All hands | Email | [link] |
| Customer impact | Support | Email | [link] |

### Stakeholder Contacts
| Role | Name | Phone | Email |
|------|------|-------|-------|
| DR Coordinator | <!-- name --> | <!-- phone --> | <!-- email --> |
| Engineering Lead | <!-- name --> | <!-- phone --> | <!-- email --> |
| Executive Sponsor | <!-- name --> | <!-- phone --> | <!-- email --> |

## DR Testing

### Test Schedule
| Test Type | Frequency | Duration | Participants |
|-----------|-----------|----------|--------------|
| Tabletop exercise | Quarterly | 2 hours | DR team |
| Partial failover | Semi-annual | 4 hours | Engineering |
| Full DR drill | Annual | 8 hours | All teams |

### Test Checklist
- [ ] Backup restoration works
- [ ] Failover completes within RTO
- [ ] Data loss within RPO
- [ ] All critical services operational
- [ ] External integrations verified
- [ ] Communication plan executed

## Post-Incident

### Recovery Validation
- [ ] All services healthy
- [ ] Data integrity confirmed
- [ ] Performance baseline restored
- [ ] Security controls verified

### Post-Mortem
1. Timeline documentation
2. Root cause analysis
3. Recovery effectiveness review
4. Improvement action items

## Maintenance

### Plan Updates
- Review: Quarterly
- Update: On infrastructure changes
- Approval: DR Coordinator + Engineering Lead

### Training
| Audience | Training | Frequency |
|----------|----------|-----------|
| Engineering | DR procedures | Quarterly |
| On-call | Failover runbook | Monthly |
| Leadership | Business continuity | Annually |
`,
	Sections: []TemplateSection{
		{ID: "objectives", Title: "Recovery Objectives", Required: true, Description: "RTO/RPO targets"},
		{ID: "impact", Title: "Business Impact Analysis", Required: true, Description: "Critical services"},
		{ID: "scenarios", Title: "Disaster Scenarios", Required: true, Description: "Risk scenarios"},
		{ID: "procedures", Title: "Recovery Procedures", Required: true, Description: "Step-by-step recovery"},
		{ID: "backup", Title: "Backup Strategy", Required: true, Description: "Backup configuration"},
		{ID: "testing", Title: "DR Testing", Required: true, Description: "Validation testing"},
	},
}

var sloDocumentTemplate = &Template{
	DocType:     DocSLODocument,
	Name:        "SLO Document",
	Description: "Service Level Objectives defining reliability targets and error budgets.",
	Content: `---
title: "{{.Title}} - SLO Document"
author: "{{.Author}}"
date: "{{.Date}}"
version: "{{.Version}}"
status: draft
---

# SLO Document: {{.Title}}

## Overview

{{.Description}}

## Service Description

| Property | Value |
|----------|-------|
| Service | {{.Title}} |
| Owner | <!-- team --> |
| Tier | <!-- 1/2/3 --> |
| Dependencies | <!-- critical dependencies --> |

## SLI Definitions

### Availability SLI
- **Definition**: Proportion of successful HTTP requests
- **Good Event**: HTTP status code < 500
- **Total Events**: All HTTP requests to the service
- **Measurement**:
` + "```promql" + `
sum(rate(http_requests_total{status!~"5.."}[5m])) /
sum(rate(http_requests_total[5m]))
` + "```" + `

### Latency SLI
- **Definition**: Proportion of requests faster than threshold
- **Good Event**: Request completed in < 200ms
- **Total Events**: All HTTP requests
- **Measurement**:
` + "```promql" + `
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
` + "```" + `

### Correctness SLI
- **Definition**: Proportion of valid responses
- **Good Event**: Response passes validation
- **Total Events**: All responses
- **Measurement**: Application-specific validation checks

## SLO Targets

### Primary SLOs
| SLI | Target | Window | Error Budget |
|-----|--------|--------|--------------|
| Availability | 99.9% | 30 days | 43.2 minutes/month |
| Latency (P95) | < 200ms | 30 days | N/A |
| Latency (P99) | < 500ms | 30 days | N/A |

### Secondary SLOs
| SLI | Target | Window | Notes |
|-----|--------|--------|-------|
| Correctness | 99.99% | 30 days | Data validation |
| Throughput | > 1000 rps | Peak hour | Capacity baseline |

## Error Budget Policy

### Error Budget Calculation
` + "```" + `
Error Budget = (1 - SLO) × Window
For 99.9% SLO over 30 days:
  Error Budget = 0.001 × 30 days = 43.2 minutes
` + "```" + `

### Budget Consumption Tiers
| Budget Remaining | Status | Actions |
|------------------|--------|---------|
| > 50% | Healthy | Normal operations |
| 25-50% | Warning | Reduce deployment velocity |
| 10-25% | Critical | Focus on reliability |
| < 10% | Exhausted | Freeze non-critical changes |

### Budget Exhaustion Response
When error budget is exhausted:
1. Halt non-essential deployments
2. Focus engineering on reliability
3. Post-mortem for major incidents
4. Executive review if budget exhausted

## Alerting

### SLO-Based Alerts
| Alert | Condition | Severity | Window |
|-------|-----------|----------|--------|
| BurnRateCritical | 14.4x burn rate | P1 | 1 hour |
| BurnRateHigh | 6x burn rate | P2 | 6 hours |
| BurnRateElevated | 1x burn rate | P3 | 24 hours |
| BudgetLow | < 25% remaining | P3 | 30 days |

### Multi-Window Alert Example
` + "```promql" + `
# Fast burn: 14.4x burn rate for 1 hour
(1 - (sum(rate(http_requests_total{status!~"5.."}[1h])) /
      sum(rate(http_requests_total[1h])))) > (14.4 * 0.001)
` + "```" + `

## Dashboards

### SLO Dashboard Components
1. **Current SLI values** - Real-time SLI metrics
2. **Error budget status** - Remaining budget %
3. **Burn rate** - Current vs allowed consumption
4. **Historical trend** - SLI over time
5. **Top error contributors** - Error breakdown

### Dashboard Links
- [SLO Overview](url)
- [Error Budget Tracking](url)
- [Incident Impact](url)

## Reporting

### Weekly Report
- SLI performance vs targets
- Error budget consumed
- Top incident contributors
- Trend analysis

### Monthly Review
- SLO attainment summary
- Error budget reconciliation
- Improvement initiatives
- Target adjustments (if needed)

## SLO Review Process

### Quarterly Review
1. Analyze SLO performance data
2. Review customer impact correlation
3. Assess target appropriateness
4. Propose adjustments if needed
5. Document changes

### Change Criteria
Consider adjusting SLOs when:
- Customer expectations change
- Infrastructure capabilities change
- Business requirements evolve
- Historical data suggests different targets

## Appendix

### Glossary
| Term | Definition |
|------|------------|
| SLI | Service Level Indicator - a metric measuring service behavior |
| SLO | Service Level Objective - target value for an SLI |
| SLA | Service Level Agreement - contractual commitment |
| Error Budget | Allowed unreliability within SLO target |
| Burn Rate | Rate of error budget consumption |

### References
- [Google SRE Book - SLOs](https://sre.google/sre-book/service-level-objectives/)
- [Implementing SLOs](https://sre.google/workbook/implementing-slos/)
`,
	Sections: []TemplateSection{
		{ID: "sli-definitions", Title: "SLI Definitions", Required: true, Description: "What we measure"},
		{ID: "slo-targets", Title: "SLO Targets", Required: true, Description: "Target values"},
		{ID: "error-budget", Title: "Error Budget Policy", Required: true, Description: "Budget management"},
		{ID: "alerting", Title: "Alerting", Required: true, Description: "SLO-based alerts"},
		{ID: "reporting", Title: "Reporting", Required: false, Description: "Regular reviews"},
	},
}
