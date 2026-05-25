# Post-PMF Enterprise Example

A comprehensive visionspec configuration for post-product-market-fit enterprises
building SaaS products with microservices, web UI, and mobile apps.

## Philosophy

Mature companies need rigorous documentation:

- **Complete coverage** - All source specs required
- **Security by design** - Security requirements in every spec
- **Accessibility first** - WCAG 2.1 AA compliance
- **Compliance ready** - SOC 2, GDPR considerations
- **GTM documentation** - Ready for product launches

## What's Required

| Spec | Required | Purpose |
|------|----------|---------|
| MRD | Yes | Market context and business case |
| PRD | Yes | Product requirements with security |
| UXD | Yes | UX with accessibility requirements |
| TRD | Yes | Technical architecture |
| Press | Yes | GTM: Press release |
| FAQ | Yes | GTM: Customer FAQs |
| spec.md | Yes | Reconciled execution spec |

### Optional but Recommended

| Spec | Purpose |
|------|---------|
| IRD | Infrastructure requirements |
| Narrative 1-pager | Executive summary |
| Narrative 6-pager | Detailed narrative |

## Key Features

### Security Requirements

All PRDs must include:

- Authentication & authorization requirements
- Data protection and encryption
- Audit logging requirements
- Compliance framework mappings

### Accessibility (WCAG 2.1 AA)

All UXDs must document:

- Perceivable requirements (contrast, text alternatives)
- Operable requirements (keyboard, focus management)
- Understandable requirements (predictability, error handling)
- Testing plan

### Platform Support

Templates include sections for:

- Web browser support matrix
- Mobile OS version requirements
- API versioning strategy
- Integration requirements

## Building

```bash
go build -o enterprise-spec ./examples/post-pmf-enterprise
```

## Usage

```bash
# Initialize a new feature
enterprise-spec init payment-processing

# Validate all specs
enterprise-spec lint

# Evaluate all specs (strict mode)
enterprise-spec eval --all

# Check compliance requirements
enterprise-spec compliance check

# View security checklist
enterprise-spec security checklist

# View compliance frameworks
enterprise-spec compliance frameworks
```

## Evaluation Rubrics

The enterprise rubrics are stricter than defaults:

| Spec | Required Categories | Max Medium Findings |
|------|---------------------|---------------------|
| PRD | 7 (including Security) | 2 |
| UXD | 8 (including 3 WCAG categories) | 2 |

### PRD Evaluation Categories

1. Problem Definition (required)
2. Requirements Coverage (required)
3. User Stories & Acceptance (required)
4. **Security Requirements** (required)
5. Platform & Integration (required)
6. Scope Definition (required)
7. Success Metrics (required)
8. Rollout Plan (optional)

### UXD Evaluation Categories

1. User Research (required)
2. User Journeys (required)
3. Information Architecture (required)
4. **Accessibility: Perceivable** (required)
5. **Accessibility: Operable** (required)
6. **Accessibility: Understandable** (required)
7. Responsive Design (required)
8. Error & Empty States (required)
9. Internationalization (optional)

## Customizing

### Adding Custom Specs

Add entries to `enterpriseSpecConfig()` in `main.go`:

```go
"security": {Required: true, Category: types.CategorySource},
```

### Modifying Rubrics

Edit files in `rubrics/` directory. They are embedded at build time.

### Adding Templates

Add files to `templates/` directory. They are embedded at build time.

## Migration from Pre-PMF

When graduating from pre-PMF to enterprise:

1. Add MRD and UXD to existing projects
2. Add security sections to existing PRDs
3. Add accessibility sections to existing UXDs
4. Create GTM documentation
5. Update visionspec.yaml with new requirements
