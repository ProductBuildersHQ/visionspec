# Organization Customization

VisionSpec is designed for organizations to build their own prescriptive CLI tools. This guide explains how to create a custom CLI with organization-specific templates, rubrics, constitutions, and app type constraints.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     Organization CLI                             │
│  (e.g., plexus-spec, acme-spec)                                  │
│                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Org         │  │ Org         │  │ Org         │              │
│  │ Templates   │  │ Rubrics     │  │ Constitutions│             │
│  │ (prescriptive)│ │ (strict)    │  │ (defaults)  │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│         │                │                │                      │
│         ▼                ▼                ▼                      │
│  ┌─────────────────────────────────────────────────┐            │
│  │              ChainLoader (fallback)              │            │
│  └─────────────────────────────────────────────────┘            │
└─────────┬────────────────┬────────────────┬─────────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    visionspec (open source)                      │
│                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Default     │  │ Default     │  │ Built-in    │              │
│  │ Templates   │  │ Rubrics     │  │ App Types   │              │
│  │ (flexible)  │  │ (choices)   │  │ (generic)   │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

## Open Source vs. Organization

| Aspect | Open Source VisionSpec | Organization CLI |
|--------|------------------------|------------------|
| **Templates** | Choices and placeholders | Prescriptive with pre-filled values |
| **Rubrics** | "Database choice documented" | "MUST use PostgreSQL with RLS" |
| **Constitutions** | None provided | Built-in org/team/project defaults |
| **App Types** | Generic constraints | Stricter organization requirements |
| **Commands** | Standard visionspec commands | Standard + org-specific commands |
| **Distribution** | `go install` from public repo | Single binary with embedded resources |

## Project Structure

```
org-visionspec/
├── cmd/
│   └── org-spec/
│       └── main.go           # CLI entry point
├── templates/                # Organization templates
│   ├── mrd.md
│   ├── prd.md
│   ├── ird.md               # Pre-filled with Pulumi, PostgreSQL
│   └── ...
├── rubrics/                  # Organization rubrics
│   ├── prd.rubric.yaml
│   ├── trd.rubric.yaml      # Stricter criteria
│   └── ...
├── constitutions/            # Organization defaults
│   ├── organization/
│   │   └── acme.yaml        # Org-wide defaults
│   ├── team/
│   │   └── platform.yaml    # Team overrides
│   └── project/             # Project-specific (optional)
├── apptypes/                 # App type constraints
│   ├── microservice.yaml
│   ├── website.yaml
│   └── ...
├── go.mod
└── README.md
```

## Building an Organization CLI

### Step 1: Create the Project

```bash
mkdir org-visionspec && cd org-visionspec
go mod init github.com/myorg/org-visionspec
go get github.com/ProductBuildersHQ/visionspec
```

### Step 2: Embed Organization Resources

```go
// cmd/org-spec/main.go
package main

import (
    "embed"
    "os"

    "github.com/ProductBuildersHQ/visionspec/pkg/apptypes"
    "github.com/ProductBuildersHQ/visionspec/pkg/cli"
    "github.com/ProductBuildersHQ/visionspec/pkg/constitution"
    "github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
    "github.com/ProductBuildersHQ/visionspec/pkg/templates"
    "github.com/spf13/cobra"
)

//go:embed templates/*.md
var orgTemplates embed.FS

//go:embed rubrics/*.yaml
var orgRubrics embed.FS

//go:embed constitutions/**/*.yaml
var orgConstitutions embed.FS

//go:embed apptypes/*.yaml
var orgAppTypes embed.FS

func main() {
    root := &cobra.Command{
        Use:   "org-spec",
        Short: "Organization VisionSpec CLI",
        Long: `Custom visionspec CLI with organization standards.

Enforced standards:
- Go backend with Huma+Chi for REST APIs
- PostgreSQL with RLS for multi-tenancy
- Pulumi (Go SDK) for Infrastructure as Code
- 99.9% availability minimum for microservices`,
    }

    cfg := orgConfig()
    cli.AddCommandsTo(root, cfg)

    // Add organization-specific commands
    root.AddCommand(complianceCmd())

    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### Step 3: Configure Loaders with ChainLoader

The key pattern is `ChainLoader`: organization resources first, then fall back to visionspec defaults.

```go
func orgConfig() *cli.Config {
    cfg := cli.DefaultConfig()
    cfg.Version = "1.0.0-acme"

    // Templates: Org first, then visionspec defaults
    cfg.TemplateLoader = templates.NewChainLoader(
        templates.NewEmbedFSLoader(orgTemplates, "templates"),
        templates.EmbeddedLoader(),
    )

    // Rubrics: Org first, then visionspec defaults
    cfg.RubricLoader = rubrics.NewChainLoader(
        rubrics.NewEmbedFSLoader(orgRubrics, "rubrics"),
        rubrics.EmbeddedLoader(),
    )

    // Constitutions: Org-specific only (no fallback)
    cfg.ConstitutionLoader = constitution.NewEmbeddedLoader(
        orgConstitutions, "constitutions",
    )

    // App Types: Org first, then visionspec defaults
    cfg.AppTypeLoader = apptypes.NewChainLoader(
        apptypes.NewEmbeddedLoader(orgAppTypes, "apptypes"),
        apptypes.DefaultLoader(),
    )

    return cfg
}
```

### Step 4: Create Organization Templates

Organization templates are more prescriptive. Instead of placeholders, they have pre-filled values:

```markdown
<!-- templates/ird.md -->
# Infrastructure Requirements Document (IRD)

## Required Declarations

### Infrastructure as Code (IaC) Declaration

| Choice | Tool | Justification |
|--------|------|---------------|
| [x] **Pulumi** | Go SDK | Organization standard: type-safe, Go-native |

### Observability Declaration

| Pillar | Declaration | Tool/Platform | Justification |
|--------|-------------|---------------|---------------|
| **Metrics** | [x] Implementing | Prometheus + Grafana | Organization standard |
| **Traces** | [x] Implementing | OpenTelemetry | Organization standard |
| **Logging** | [x] Implementing | Loki + Grafana | Organization standard |

## Database

- **Type:** PostgreSQL (organization standard)
- **Multi-tenancy:** Row Level Security (RLS)
- **Tenant context:** `SET app.current_tenant = :tenant_id`
```

### Step 5: Create Organization Rubrics

Organization rubrics have stricter evaluation criteria:

Rubrics use the shared structured-evaluation format; the spec type comes from
the filename (`trd.rubric.yaml` → `trd`):

```yaml
# rubrics/trd.rubric.yaml
id: trd-rubric
name: "Organization TRD Evaluation"
categories:
  - id: language_choice
    name: "Language Choice"
    weight: 0.15
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["Uses Go for backend services. Exceptions documented with approval."]}
        - {value: partial, criteria: ["Uses Go but missing exception documentation for non-Go components."]}
        - {value: fail, criteria: ["Uses non-Go backend without documented exception approval."]}

  - id: database_choice
    name: "Database Choice"
    weight: 0.15
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["Uses PostgreSQL with RLS for multi-tenancy."]}
        - {value: partial, criteria: ["Uses PostgreSQL but RLS not configured."]}
        - {value: fail, criteria: ["Uses non-PostgreSQL database without exception approval."]}
```

## Constitution System

Constitutions define organizational defaults that flow through a hierarchy.

### Hierarchy

```
org/acme.yaml              # Organization-wide defaults
    ↓ inherits
team/platform.yaml         # Team overrides
    ↓ inherits
project/myservice.yaml     # Project-specific choices
```

### Organization Constitution

```yaml
# constitutions/organization/acme.yaml
apiVersion: visionspec/v1
kind: Constitution
metadata:
  name: acme
  level: organization
  version: "1.0"

technical:
  languages:
    backend:
      primary: go
      allowed: [go, rust]
      exceptionsRequire: approval
    frontend:
      primary: typescript
    wasm: rust

  apis:
    rest:
      framework: huma-chi
      specFormat: openapi-3.1
      styleGuide: google-api-design-guide

  database:
    relational: postgresql
    multiTenancy: rls

  tenancy:
    model: multi-tenant
    isolation: rls
    tenantIdHeader: X-Tenant-ID

infrastructure:
  iac:
    tool: pulumi
    language: go

  availability:
    target: "99.9"
    rto: "1h"
    rpo: "15m"

  observability:
    metrics:
      platform: prometheus
      visualization: grafana
    traces:
      collector: opentelemetry-collector
    logging:
      platform: loki

prompts:
  languageChoice: |
    Go is the primary backend language.

    Consider Rust only for:
    - CPU-intensive processing
    - WASM modules
    - Memory-critical components

    Exceptions require tech lead approval.
```

### Team Constitution

```yaml
# constitutions/team/platform.yaml
apiVersion: visionspec/v1
kind: Constitution
metadata:
  name: platform
  level: team
  inherits: organization/acme  # Inherits from org

infrastructure:
  # Team has higher availability requirement
  availability:
    target: "99.99"
    multiAZ: true
```

### Project Constitution

```yaml
# constitutions/project/myservice.yaml
apiVersion: visionspec/v1
kind: Constitution
metadata:
  name: myservice
  level: project
  inherits: team/platform

# Project-specific deployment targets
infrastructure:
  localDev:
    priority:
      - binaries
      - podman
      - localstack  # Need AWS emulation

# Documented exceptions
exceptions:
  - field: technical.languages.backend
    value: rust
    component: wasm-inference
    justification: "WASM performance requirements"
    approvedBy: tech-lead
    approvedAt: "2026-05-15"
```

### Resolving Constitutions

```go
// Resolve inheritance chain
resolver := constitution.NewResolver(cfg.ConstitutionLoader)
resolved, err := resolver.ResolveChain("project/myservice")

// Result has all inherited values merged
fmt.Println(resolved.Technical.Languages.Backend.Primary) // "go"
fmt.Println(resolved.Infrastructure.Availability.Target)   // "99.99"
```

## App Type Specifications

App types define constraints per application category:

```yaml
# apptypes/microservice.yaml
apiVersion: visionspec/v1
kind: AppTypeSpec
metadata:
  name: microservice
  version: "1.0"

artifacts:
  required:
    - binary
    - container-image
  optional:
    - openapi-spec
    - helm-chart

defaults:
  technical:
    apiStyles: [rest, grpc]
    embeddedDb: false
    statefulAllowed: false
  infrastructure:
    containerized: true
    horizontalScaling: true
    minAvailabilityTarget: "99.9"

specs:
  required: [mrd, prd, trd, ird]
  optional: [uxd, press, faq]

constraints:
  tenancy:
    allowed: [single-tenant, multi-tenant]
  availability:
    minimum: "99.9"
```

## Adding Organization Commands

```go
func complianceCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "compliance",
        Short: "Check project compliance with organization standards",
        RunE: func(cmd *cobra.Command, args []string) error {
            // 1. Load project constitution
            // 2. Resolve inheritance chain
            // 3. Validate against org requirements
            // 4. Report deviations
            return nil
        },
    }
}
```

## Distribution

Build a single binary with all resources embedded:

```bash
# Build for current platform
go build -o org-spec ./cmd/org-spec

# Cross-compile for distribution
GOOS=linux GOARCH=amd64 go build -o org-spec-linux ./cmd/org-spec
GOOS=darwin GOARCH=arm64 go build -o org-spec-darwin ./cmd/org-spec
```

The binary is self-contained—no external files needed.

## Best Practices

1. **Version constitutions** - Use semantic versioning for constitution changes
2. **Document exceptions** - Require justification and approval for deviations
3. **Gradual rollout** - Use ChainLoader to fall back to defaults during migration
4. **Test rubrics** - Ensure evaluation criteria are unambiguous
5. **Embed everything** - Single binary distribution simplifies deployment

## See Also

- [Configuration Guide](configuration.md)
- [Custom Profiles](custom-profiles.md)
- [`examples/org-cli/`](https://github.com/ProductBuildersHQ/visionspec/tree/main/examples/org-cli)
