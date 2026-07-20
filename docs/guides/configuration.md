# Configuration Reference

Complete reference for `visionspec.yaml` configuration options.

## Overview

Each VisionSpec project contains a `visionspec.yaml` file in its root directory. This file configures LLM providers, export targets, context sources, rubrics, and spec requirements.

## Minimal Configuration

```yaml
name: my-project
```

Most settings have sensible defaults. Only the project name is required.

## Full Configuration

```yaml
name: my-project

# Constitution file path (relative or absolute)
constitution: ../CONSTITUTION.md

# LLM provider configuration
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  temperature: 0.0
  max_tokens: 8192

# Spec requirements configuration
spec_config:
  mrd:
    required: true
  prd:
    required: true
  uxd:
    required: false
  trd:
    required: true
  ird:
    required: false

# Custom rubrics configuration
rubrics:
  directory: .rubrics
  strict_mode: false
  max_critical: 0
  max_high: 0
  max_medium: -1
  overrides:
    prd: custom/prd-v2.rubric.yaml

# Context sources for grounding
context:
  cache_ttl: 30m
  repositories:
    - path: "."
      analyze: [structure, deps, apis]
  files:
    - path: docs/architecture.md
      type: architecture

# Export target configuration
targets:
  default: speckit
  speckit:
    enabled: true
    output_dir: export/speckit
```

## Configuration Sections

### name

The project identifier. Must be kebab-case.

```yaml
name: user-onboarding
```

### constitution

Path to the constitution file, either relative to the project or absolute.

```yaml
constitution: ../CONSTITUTION.md
```

### llm

Configure the LLM provider for evaluations and synthesis.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `provider` | string | `anthropic` | LLM provider (anthropic, openai, gemini) |
| `model` | string | varies by provider | Model identifier |
| `temperature` | float | `0.0` | Randomness (0.0 = deterministic) |
| `max_tokens` | int | `8192` | Maximum response length |

**Example:**

```yaml
llm:
  provider: anthropic
  model: claude-sonnet-4-20250514
  temperature: 0.0
  max_tokens: 8192
```

**Environment variables:**

LLM credentials are configured via environment variables:

- `ANTHROPIC_API_KEY` - For Anthropic/Claude
- `OPENAI_API_KEY` - For OpenAI
- `GOOGLE_API_KEY` - For Google/Gemini

### spec_config

Configure which specs are required and their settings.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `required` | bool | varies | Whether spec is mandatory |
| `category` | string | auto | Spec category (source, gtm, technical) |
| `template` | string | spec type | Template to use |
| `rubric` | string | spec type | Rubric to use |

**Built-in spec types and defaults:**

| Spec | Category | Default Required |
|------|----------|------------------|
| `mrd` | source | true |
| `prd` | source | true |
| `uxd` | source | false |
| `press` | gtm | false |
| `faq` | gtm | false |
| `narrative` | gtm | false |
| `trd` | technical | true |
| `ird` | technical | false |

**Example:**

```yaml
spec_config:
  mrd:
    required: true
  prd:
    required: true
  uxd:
    required: true
    rubric: uxd-mobile  # Use custom rubric
  trd:
    required: true
  # Custom spec type
  api-spec:
    required: true
    category: technical
    template: api
    rubric: api
```

### rubrics

Configure custom rubric loading for evaluations.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `directory` | string | - | Directory containing `.rubric.yaml` files |
| `overrides` | map | - | Map spec types to specific rubric files |
| `strict_mode` | bool | `false` | Require all categories to pass |
| `max_critical` | int | `0` | Maximum critical findings allowed |
| `max_high` | int | `0` | Maximum high findings allowed |
| `max_medium` | int | `-1` | Maximum medium findings (-1 = unlimited) |

**Example:**

```yaml
rubrics:
  directory: .rubrics
  strict_mode: true
  max_critical: 0
  max_high: 0
  max_medium: 3
  overrides:
    prd: enterprise/prd.rubric.yaml
```

Rubric files in the directory should be named `{spec-type}.rubric.yaml` (e.g., `prd.rubric.yaml`).

Rubrics resolve in order: `overrides` (per spec type) → `directory` → the
profile's built-in rubrics. Relative paths are resolved against the project
directory. `strict_mode` and the `max_*` finding limits apply to every
evaluation (`eval`, `synthesize --eval`, and the MCP `run_eval` / `eval_draft`
tools). Rubric files use the structured-evaluation native format.

### context

Configure context sources for grounding technical spec synthesis.

#### context.cache_ttl

How long to cache gathered context data.

```yaml
context:
  cache_ttl: 30m
```

#### context.repositories

Git repositories to analyze for context.

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Local path to repository |
| `url` | string | Remote URL (alternative to path) |
| `branch` | string | Branch to analyze (for remote repos) |
| `include` | []string | Glob patterns to include |
| `exclude` | []string | Glob patterns to exclude |
| `analyze` | []string | What to analyze: structure, deps, apis |
| `graphize` | string | Path to graphize graph or "auto" |
| `max_depth` | int | Maximum directory depth |

**Example:**

```yaml
context:
  repositories:
    - path: "."
      analyze: [structure, deps, apis]
      exclude: ["vendor/**", "node_modules/**"]
    - url: "https://github.com/org/shared-lib"
      branch: main
      include: ["pkg/**"]
```

#### context.graphize

Standalone graphize graphs for requirement tracing.

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Path to graphize directory |
| `name` | string | Display name for the graph |
| `include_nodes` | []string | Node types to include |
| `include_edges` | []string | Edge types to include |

**Example:**

```yaml
context:
  graphize:
    - path: .graphize
      name: spec-graph
      include_nodes: [requirement, user_story]
      include_edges: [traces_to, depends_on]
```

#### context.files

Local files to include as context.

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | File path (supports globs) |
| `type` | string | Context type: architecture, adr, api_spec |
| `max_size` | int | Maximum file size in bytes |

**Example:**

```yaml
context:
  files:
    - path: docs/architecture.md
      type: architecture
    - path: docs/adr/*.md
      type: adr
    - path: api/openapi.yaml
      type: api_spec
```

#### context.mcp_servers

MCP servers for external context (not yet implemented).

| Field | Type | Description |
|-------|------|-------------|
| `command` | string | Command to run |
| `args` | []string | Command arguments |
| `env` | map | Environment variables |
| `config` | map | Server-specific config |
| `timeout` | duration | Connection timeout |

**Example:**

```yaml
context:
  mcp_servers:
    jira:
      command: npx
      args: ["-y", "@anthropic/mcp-jira"]
      timeout: 30s
```

### targets

Configure export target settings.

#### targets.default

Default export target when none specified.

```yaml
targets:
  default: speckit
```

#### targets.speckit

GitHub Spec-Kit format configuration.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable this target |
| `output_dir` | string | `export/speckit` | Output directory |
| `branch_numbering` | string | `sequential` | Numbering: sequential, timestamp |

```yaml
targets:
  speckit:
    enabled: true
    output_dir: export/speckit
    branch_numbering: sequential
```

#### targets.gsd

Get Shit Done format configuration.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable this target |
| `output_dir` | string | `export/gsd` | Output directory |
| `model_profile` | string | `balanced` | Profile: balanced, quality, budget |

```yaml
targets:
  gsd:
    enabled: true
    model_profile: quality
```

#### targets.gastown

GasTown formula format configuration.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable this target |
| `formula_type` | string | `convoy` | Type: convoy, workflow, expansion |
| `rig` | string | - | Rig name to use |

```yaml
targets:
  gastown:
    enabled: true
    formula_type: convoy
    rig: my-rig
```

#### targets.gascity

GasCity orchestration format configuration.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable this target |
| `city_dir` | string | `export/gascity` | Output directory |

```yaml
targets:
  gascity:
    enabled: true
    city_dir: export/gascity
```

## Configuration Profiles

VisionSpec supports configuration profiles that provide pre-configured settings for different organizational stages.

**Available profiles:**

| Profile | Description |
|---------|-------------|
| `0-1` | Early-stage, minimal requirements |
| `startup` | Growing team, balanced requirements |
| `growth` | Scaling organization, comprehensive coverage |
| `enterprise` | Full compliance and traceability |

Export a profile to use as a starting point:

```bash
visionspec profiles export startup
```

See the [Profiles Guide](profiles.md) for details.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `OPENAI_API_KEY` | OpenAI API key |
| `GOOGLE_API_KEY` | Google AI API key |
| `MULTISPEC_CONFIG` | Override config file path |
| `MULTISPEC_PROFILE` | Default profile to use |

## See Also

- [Profiles Guide](profiles.md) - Pre-configured settings
- [Custom Rubrics](../rubrics/custom-rubrics.md) - Rubric customization
- [Context Sources](../context-sources.md) - Context configuration
