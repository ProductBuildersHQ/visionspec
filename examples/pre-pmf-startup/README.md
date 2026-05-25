# Pre-PMF Startup Example

A lightweight visionspec configuration for pre-product-market-fit startups.

## Philosophy

Early-stage startups need to move fast. This configuration:

- **Minimizes required documentation** - Only PRD is required
- **Focuses on essentials** - Simple templates without bureaucracy
- **Enables AI-assisted development** - Just enough structure for AI coding assistants
- **Supports iteration** - Add more specs as you grow

## What's Required

| Spec | Required | When to Add |
|------|----------|-------------|
| PRD | Yes | Always - defines what you're building |
| MRD | No | When seeking funding or pivoting |
| UXD | No | When UX becomes a competitive advantage |
| TRD | No | When architecture decisions need documentation |
| GTM | No | Wait until post-PMF |

## Building

```bash
go build -o startup-spec ./examples/pre-pmf-startup
```

## Usage

```bash
# Initialize a new feature
startup-spec init user-auth

# Check your PRD
startup-spec eval prd

# Velocity check - are you over-documenting?
startup-spec velocity check
```

## Customizing Further

You can adjust the spec requirements by modifying `startupSpecConfig()` in `main.go`.

To add more templates or rubrics:

1. Add files to `templates/` or `rubrics/`
2. They will be embedded in the binary via `//go:embed`

## Graduating to Post-PMF

Once you find product-market fit, consider migrating to the post-PMF enterprise
configuration which adds:

- Required MRD and UXD
- Security and compliance requirements
- GTM documentation
- Formal approval workflows
