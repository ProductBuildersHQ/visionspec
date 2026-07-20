<!--
  AI-Generated Plan
  Source: Claude Code planning session
  Date: 2025-05-24

  This plan.md is a copy of the AI-assisted planning output.
  It documents the implementation approach before coding begins.
-->

# VisionSpec v0.3.0: Composability

## Overview

Enable organizations (companies, open source projects, non-profits) to import visionspec as a library and build custom CLI tools with:

- Organization-specific spec templates
- Organization-specific evaluation rubrics
- Configurable mandatory/optional specs
- Composable Cobra commands

## Requirements Summary

| Requirement | Description |
|-------------|-------------|
| CLI as Library | Move CLI commands to `pkg/cli` for import/composition |
| Custom Templates | Load templates from project/org directory |
| Custom Rubrics | Load rubrics from YAML files, not just Go code |
| Configurable Specs | Make mandatory/optional configurable per project |

---

## Architecture

### 1. CLI as Library (`pkg/cli`)

Move from `internal/cli` to `pkg/cli` with composition API:

```go
// pkg/cli/cli.go
package cli

import "github.com/spf13/cobra"

// Config allows customization of CLI behavior
type Config struct {
    // TemplateLoader for custom templates
    TemplateLoader templates.Loader

    // RubricLoader for custom rubrics
    RubricLoader rubrics.Loader

    // SpecConfig for mandatory/optional settings
    SpecConfig *SpecConfig
}

// DefaultConfig returns visionspec defaults
func DefaultConfig() *Config

// AddCommandsTo adds all visionspec commands to a root command
func AddCommandsTo(root *cobra.Command, cfg *Config)

// Commands returns individual commands for selective addition
func Commands(cfg *Config) struct {
    Init      *cobra.Command
    Lint      *cobra.Command
    Status    *cobra.Command
    Graph     *cobra.Command
    Eval      *cobra.Command
    Approve   *cobra.Command
    Synthesize *cobra.Command
    Reconcile *cobra.Command
    Export    *cobra.Command
    Serve     *cobra.Command
}
```

**Usage by organization:**

```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/ProductBuildersHQ/visionspec/pkg/cli"
    "github.com/ProductBuildersHQ/visionspec/pkg/templates"
    "github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
)

func main() {
    root := &cobra.Command{
        Use:   "org-spec",
        Short: "Organization specification tool",
    }

    cfg := cli.DefaultConfig()

    // Override with org templates
    cfg.TemplateLoader = templates.NewChainLoader(
        templates.NewFileLoader(".orgspec/templates"),
        templates.EmbeddedLoader(), // fallback
    )

    // Override with org rubrics
    cfg.RubricLoader = rubrics.NewChainLoader(
        rubrics.NewFileLoader(".orgspec/rubrics"),
        rubrics.EmbeddedLoader(), // fallback
    )

    // Add all visionspec commands
    cli.AddCommandsTo(root, cfg)

    // Add org-specific commands
    root.AddCommand(orgAuditCmd, orgComplianceCmd)

    root.Execute()
}
```

### 2. Custom Templates (`pkg/templates`)

Add loader interface for template sources:

```go
// pkg/templates/loader.go

// Loader loads templates from various sources
type Loader interface {
    // Load returns template content for a spec type
    Load(specType types.SpecType) (*Template, error)

    // Available returns all available spec types
    Available() []types.SpecType
}

// EmbeddedLoader returns the built-in embedded templates
func EmbeddedLoader() Loader

// FileLoader loads templates from a directory
// Templates are named: {spec-type}.md (e.g., prd.md, custom-spec.md)
func NewFileLoader(dir string) Loader

// ChainLoader tries loaders in order, returning first match
func NewChainLoader(loaders ...Loader) Loader

// RegisterCustomType registers a custom spec type for the loader
func (l *FileLoader) RegisterCustomType(name string, category types.SpecCategory)
```

**Directory structure for org templates:**

```
.orgspec/templates/
├── mrd.md           # Override default MRD
├── prd.md           # Override default PRD
├── security.md      # Custom org spec type
└── compliance.md    # Custom org spec type
```

### 3. Custom Rubrics (`pkg/rubrics`)

Add loader interface for rubric sources:

```go
// pkg/rubrics/loader.go

// Loader loads rubrics from various sources
type Loader interface {
    // Load returns rubric for a spec type
    Load(specType types.SpecType) (*RubricSet, error)

    // Available returns all available spec types with rubrics
    Available() []types.SpecType
}

// EmbeddedLoader returns the built-in Go-defined rubrics
func EmbeddedLoader() Loader

// FileLoader loads rubrics from YAML files in a directory
// Rubrics are named: {spec-type}.rubric.yaml
func NewFileLoader(dir string) Loader

// ChainLoader tries loaders in order, returning first match
func NewChainLoader(loaders ...Loader) Loader
```

**Rubric YAML format (compatible with structured-evaluation):**

```yaml
# .orgspec/rubrics/prd.rubric.yaml
spec_type: prd
name: "Organization PRD Rubric"
description: "Evaluation rubric for Product Requirements Documents"
version: "1.0"

categories:
  - id: security-requirements
    name: "Security Requirements"
    description: "Are security requirements explicitly stated?"
    weight: 2.0
    required: true
    criteria:
      pass: "Security requirements are comprehensive with threat models"
      partial: "Basic security requirements present but incomplete"
      fail: "No security requirements or major gaps"

  - id: accessibility
    name: "Accessibility"
    description: "Does spec address accessibility requirements?"
    weight: 1.5
    required: true
    criteria:
      pass: "WCAG 2.1 AA compliance addressed comprehensively"
      partial: "Some accessibility considerations but incomplete"
      fail: "Accessibility not addressed"

pass_criteria:
  require_all_pass: false
  max_critical: 0
  max_high: 0
  max_medium: 3
```

### 4. Configurable Spec Requirements

Add spec configuration to `visionspec.yaml`:

```yaml
# visionspec.yaml
name: user-onboarding
description: "User onboarding feature"

# Spec requirements configuration
specs:
  # Source specs
  mrd:
    required: false          # Optional for this org
    template: mrd            # Use default template
  prd:
    required: true           # Mandatory
  uxd:
    required: true

  # Custom org specs
  security:
    required: true
    category: source         # source | gtm | technical
    template: security       # Maps to security.md template
    rubric: security         # Maps to security.rubric.yaml

  accessibility:
    required: true
    category: source

  # Technical specs
  trd:
    required: true
  ird:
    required: false

# LLM configuration (existing)
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
```

**Go types:**

```go
// pkg/types/spec_config.go

// SpecRequirement defines requirements for a spec type
type SpecRequirement struct {
    Required bool         `yaml:"required"`
    Category SpecCategory `yaml:"category,omitempty"`
    Template string       `yaml:"template,omitempty"`
    Rubric   string       `yaml:"rubric,omitempty"`
}

// SpecConfig maps spec types to their requirements
type SpecConfig struct {
    Specs map[string]*SpecRequirement `yaml:"specs"`
}

// IsRequired returns whether a spec is required (with fallback to defaults)
func (sc *SpecConfig) IsRequired(specType SpecType) bool

// GetCategory returns the category for a spec type
func (sc *SpecConfig) GetCategory(specType string) SpecCategory

// CustomSpecs returns all custom (non-standard) spec types
func (sc *SpecConfig) CustomSpecs() []string
```

---

## Implementation Plan

### Phase 1: Spec Configuration

1. **Add SpecConfig types** (`pkg/types/spec_config.go`)
   - `SpecRequirement` struct
   - `SpecConfig` with methods
   - Update `Project` to include `SpecConfig`

2. **Update config loading** (`pkg/config/config.go`)
   - Parse `specs:` section from visionspec.yaml
   - Merge with defaults

3. **Update IsRequired logic** (`pkg/types/spec.go`)
   - Use project config instead of hardcoded values

### Phase 2: Template Loader

1. **Create Loader interface** (`pkg/templates/loader.go`)
   - `Loader` interface
   - `EmbeddedLoader` (wraps current behavior)
   - `FileLoader` (reads from directory)
   - `ChainLoader` (tries multiple loaders)

2. **Update template functions** (`pkg/templates/templates.go`)
   - Accept `Loader` parameter or use default
   - Support custom spec types from file

3. **Add custom type registration**
   - Allow non-standard spec types from files

### Phase 3: Rubric Loader

1. **Create Loader interface** (`pkg/rubrics/loader.go`)
   - `Loader` interface
   - `EmbeddedLoader` (wraps current registry)
   - `FileLoader` (reads YAML files)
   - `ChainLoader`

2. **Define YAML schema** (`pkg/rubrics/yaml.go`)
   - Parse/write structured-evaluation's `rubric.RubricSet` directly (canonical format)
   - The legacy `RubricYAML` flat struct was retired in v0.14.0

3. **Update rubric functions** (`pkg/rubrics/rubrics.go`)
   - Accept `Loader` parameter or use default

### Phase 4: CLI Library

1. **Create pkg/cli package**
   - Move commands from `internal/cli/commands.go`
   - Add `Config` struct
   - Add `AddCommandsTo()` function
   - Add `Commands()` for selective access

2. **Update cmd/visionspec/main.go**
   - Use `pkg/cli` with default config

3. **Add integration tests**
   - Test composing custom CLI
   - Test custom loaders

### Phase 5: Documentation & Examples

1. **Create example** (`examples/org-cli/`)
   - Sample org CLI using visionspec as library
   - Custom templates and rubrics
   - Custom spec types

2. **Update documentation**
   - `docs/composability/` section
   - Configuration reference for specs
   - Template and rubric customization guide

---

## File Changes

### New Files

| Path | Description |
|------|-------------|
| `pkg/types/spec_config.go` | SpecRequirement and SpecConfig types |
| `pkg/templates/loader.go` | Loader interface and implementations |
| `pkg/rubrics/loader.go` | Loader interface and implementations |
| `pkg/rubrics/yaml.go` | YAML parsing for rubrics |
| `pkg/cli/cli.go` | CLI composition API |
| `pkg/cli/config.go` | CLI configuration |
| `pkg/cli/commands.go` | Moved from internal/cli |
| `examples/org-cli/` | Example organization integration |
| `docs/composability/integration.md` | Composability guide |

### Modified Files

| Path | Changes |
|------|---------|
| `pkg/types/project.go` | Add SpecConfig field |
| `pkg/types/spec.go` | Update IsRequired to use config |
| `pkg/config/config.go` | Parse specs section |
| `pkg/templates/templates.go` | Use Loader interface |
| `pkg/rubrics/rubrics.go` | Use Loader interface |
| `cmd/visionspec/main.go` | Use pkg/cli |
| `internal/cli/` | Remove (moved to pkg/cli) |

---

## Testing Strategy

1. **Unit tests for loaders**
   - FileLoader with temp directories
   - ChainLoader precedence
   - YAML rubric parsing

2. **Integration tests**
   - Custom CLI composition
   - Project with custom specs
   - Eval with custom rubrics

3. **Example validation**
   - Build and run org-cli example
   - Verify custom templates work
   - Verify custom rubrics work

---

## Compatibility

- **Backward compatible**: Projects without `specs:` config use defaults
- **No breaking changes**: Existing code continues to work
- **Gradual adoption**: Orgs can start with one custom spec and expand

---

## Use Cases

### Open Source Project

```yaml
specs:
  mrd:
    required: false    # OSS often skips market analysis
  prd:
    required: true
  uxd:
    required: false    # Optional for CLI tools
  contributing:        # Custom spec
    required: true
    category: source
```

### Non-Profit Organization

```yaml
specs:
  mrd:
    required: true
  prd:
    required: true
  impact:              # Custom: social impact assessment
    required: true
    category: source
  accessibility:       # Custom: accessibility requirements
    required: true
    category: source
```

### Enterprise

```yaml
specs:
  security:
    required: true
    category: source
  compliance:
    required: true
    category: source
  privacy:
    required: true
    category: technical
```
