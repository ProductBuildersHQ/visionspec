# VisionSpec

Multi-domain specification orchestration for humans and AI agents.

## What is VisionSpec?

VisionSpec bridges the gap between organizational intent (MRD, PRD, UXD) and executable specifications for AI coding agents. It provides a structured workflow for:

- **Dual-methodology system** - Separate requirements (WHAT) from implementation (HOW) methodologies
- **Methodology profiles** - AWS, Google, Stripe, Lean Startup, Design Thinking, JTBD
- **AIDLC workflow** - Full AWS AI-DLC integration with 12 document types across 3 phases
- **V2MOM strategic planning** - Company, department, and team-level strategic alignment
- **Domain-specific authoring** - Separate specs for PM, UX, Engineering
- **GTM synthesis** - LLM-generated press releases, FAQs, narratives (Working Backwards)
- **Technical synthesis** - LLM-generated TRD, TPD, IRD from source specs
- **Structured evaluation** - Per-domain LLM judges with 1-5 numeric scoring
- **Reconciliation** - Conflict detection and tradeoff resolution
- **Target adapters** - Export to SpecKit, GSD, GasTown, GasCity, AWS AI-DLC, OpenSpec

## Working Backwards Flow

VisionSpec implements Amazon's Working Backwards methodology:

```
1. MARKET PROBLEM (human-authored)
   mrd.md
       ↓
2. WORKING BACKWARDS (synthesized, editable)
   press.md  →  faq.md  →  prd.md
   (vision)     (scope)    (requirements)
       ↓
3. STAKEHOLDER REVIEW (synthesized, editable)
   narrative-1p.md / narrative-6p.md
       ↓
4. USER EXPERIENCE (human-authored)
   uxd.md
       ↓
5. TECHNICAL SPECS (synthesized, editable)
   trd.md  →  tpd.md  →  ird.md
   (design)   (tests)    (infra)
       ↓
6. RECONCILIATION
   All approved specs → spec.md
       ↓
7. AI EXECUTION
   spec.md → SpecKit | GSD | GasTown | GasCity
       ↓
8. POST-SHIP ALIGNMENT
   spec.md + reality → current-truth.md
```

All synthesized documents are committed to git and can be reviewed, edited, and refined by humans or collaboratively with AI assistants.

See the [Working Backwards Guide](concepts/working-backwards.md) for the full methodology.

## Quick Start

```bash
# Install
go install github.com/ProductBuildersHQ/visionspec/cmd/visionspec@v0.13.0

# Initialize a new project
visionspec init user-onboarding

# Initialize with a methodology profile
visionspec init my-product --profile aws

# Initialize with a stage profile
visionspec init my-feature --profile startup

# Validate project structure
visionspec lint

# Check project status
visionspec status
```

## Key Features

### Readiness Gates

VisionSpec tracks project readiness through configurable gates:

| Gate | Description |
|------|-------------|
| Required specs present | All required source specs exist |
| Evaluations passing | No blocking evaluation findings |
| Approvals obtained | All required specs have approvals |
| Execution spec generated | `spec.md` has been created |

### Multiple Output Formats

Status reports can be generated in multiple formats for different use cases:

- **Text** - Terminal output with icons
- **JSON** - Programmatic access
- **HTML** - Browser-viewable reports with traffic light indicators
- **Markdown** - For embedding in documentation

### MCP Integration

VisionSpec includes an MCP (Model Context Protocol) server for integration with AI coding assistants like Claude Code and Kiro CLI.

## Configuration Profiles

VisionSpec includes two types of profiles:

### Stage-Based Profiles

| Profile | Required Specs | Use Case |
|---------|---------------|----------|
| `0-1` | hypothesis | Idea validation phase |
| `startup` | prd | Pre-PMF startups |
| `growth` | prd, uxd, faq | 1-N scaling phase |
| `enterprise` | mrd, prd, uxd, trd, press, faq, spec | Post-PMF enterprises |

### Methodology Profiles

| Profile | Methodology | Best For |
|---------|-------------|----------|
| `aws` | Working Backwards | Customer-centric products |
| `google` | Design Docs + RFC | Engineering-heavy orgs |
| `stripe` | API-First | Platform/API products |
| `lean-startup` | Build-Measure-Learn | Early validation |
| `design-thinking` | Stanford d.school | Human-centered design |
| `jtbd` | Jobs to be Done | Customer motivations |

See [Frameworks](frameworks/index.md) for detailed methodology documentation.

```bash
# List available profiles
visionspec profiles list

# Initialize with a methodology profile
visionspec init my-product --profile aws

# Export a profile for customization
visionspec profiles export enterprise ./my-profile

# Use a custom profile
visionspec init my-project --profile-dir ./my-profile
```

Organizations can create custom profiles with their own templates and rubrics. See the [Custom Profiles Guide](guides/custom-profiles.md) for details.

## Project Status

See the [ROADMAP](specs/ROADMAP.md) for detailed implementation status and [Release Notes](releases/v0.13.0.md) for the latest release.

**Current Version:** v0.13.0

| Component | Status |
|-----------|--------|
| CLI (init, lint, status, eval, synthesize, reconcile) | Complete |
| MCP Server (draft workflow, eval) | Complete |
| Evaluation Engine (v2 schema with 1-5 scoring) | Complete |
| GTM & Technical Synthesis (Press, FAQ, PRD, TRD, TPD, IRD) | Complete |
| Reconciliation | Complete |
| Export (SpecKit, GSD, GasTown, GasCity, AIDLC) | Complete |
| Graphize Integration | Complete |
| Profiles & Composability | Complete |
| Methodology Profiles (AWS, Google, Stripe, Lean, DT, JTBD) | Complete |
| AI Workflow Orchestration (rules) | Complete |
| AIDLC Workflow Package (pkg/aidlc/) | Complete |
| Dual-Methodology System | Complete |
| V2MOM Strategic Planning Profiles | Complete |
