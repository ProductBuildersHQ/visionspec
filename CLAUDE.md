# CLAUDE.md

Instructions for AI assistants working with VisionSpec.

## Overview

VisionSpec is a multi-domain specification orchestration tool. It helps teams create, evaluate, and reconcile specifications (MRD, PRD, UXD, TRD, IRD) before exporting to AI coding agent execution systems.

## Project Structure

```
docs/specs/
├── CONSTITUTION.md           # Repo-level governance constraints
├── ROADMAP.md                # Cross-project priorities
└── {project}/                # Individual projects (kebab-case)
    ├── source/               # Human-authored specs
    │   ├── mrd.md            # Market Requirements
    │   ├── prd.md            # Product Requirements
    │   └── uxd.md            # User Experience Design
    ├── gtm/                  # LLM-generated GTM docs
    │   ├── press.md          # Press Release (Working Backwards)
    │   ├── faq.md            # FAQ Document
    │   └── narrative.md      # Internal Narrative
    ├── technical/            # LLM-generated technical docs
    │   ├── trd.md            # Technical Requirements
    │   └── ird.md            # Infrastructure Requirements
    ├── eval/                 # Evaluation results
    ├── spec.md               # Reconciled execution spec
    └── visionspec.yaml        # Project configuration
```

## MCP Tools

When the visionspec MCP server is configured, use these tools:

### Project Management

| Tool | Use When |
|------|----------|
| `list_projects` | User asks "what projects exist?" or you need to find projects |
| `get_project_status` | User asks about readiness or progress of a project |

### Spec Authoring (Draft Workflow)

| Tool | Use When |
|------|----------|
| `start_draft` | User wants to create/write a new spec |
| `get_draft` | Need to see current draft content |
| `update_draft` | Save progress on a draft |
| `eval_draft` | Check if draft meets quality criteria |
| `finalize_draft` | Draft is complete and passes evaluation |
| `discard_draft` | User wants to abandon a draft |
| `list_drafts` | See all in-progress drafts |

### Spec Operations

| Tool | Use When |
|------|----------|
| `get_spec` | User asks to see a spec's content |
| `get_eval` | User asks about evaluation results |
| `run_eval` | User wants to evaluate a finalized spec |
| `synthesize` | Generate TRD/IRD from source specs, or GTM docs |
| `reconcile` | Combine all specs into unified spec.md |
| `approve` | Mark a spec as approved for reconciliation |
| `export` | Export to SpecKit, GSD, GasTown, or GasCity |

## CLI Commands

If MCP is not available, use CLI via Bash:

```bash
# Project management
visionspec init <project>              # Create new project
visionspec lint [project]              # Validate structure
visionspec status -p <project>         # Check readiness

# Spec operations
visionspec create <type> -p <project>  # Scaffold spec from template
visionspec eval <type> -p <project>    # Evaluate a spec
visionspec synthesize <type> -p <project>  # Generate TRD/IRD/GTM
visionspec reconcile -p <project>      # Generate spec.md
visionspec approve <type> -p <project> # Approve for reconciliation
visionspec export <target> -p <project> # Export to target system

# Context (for grounding synthesis)
visionspec context gather -p <project> # Fetch codebase context
visionspec context show -p <project>   # Display context summary

# Profiles
visionspec profiles list               # List available profiles
visionspec profiles show <name>        # Show profile details
visionspec profiles export <name> <dir> # Export for customization
```

## Authoring Workflows

When users want to author specs, follow the skill workflows in `skills/`:

### Source Specs (Human-Authored)

1. **MRD** (`skills/author-mrd/`): Market problem, audience, business goals
2. **PRD** (`skills/author-prd/`): User stories, functional requirements
3. **UXD** (`skills/author-uxd/`): User journeys, interaction flows

### GTM Specs (LLM-Generated)

4. **Press Release** (`skills/author-press/`): Customer announcement
5. **FAQ** (`skills/author-faq/`): Challenges claims, surfaces gaps
6. **Narrative** (`skills/author-narrative-1p/`, `author-narrative-6p/`): Internal vision

### Technical Specs (LLM-Generated)

7. **TRD** (`skills/author-trd/`): Architecture, APIs, data models
8. **IRD** (`skills/author-ird/`): Deployment, scaling, operations

## Common Tasks

### "Help me write a PRD"

1. Use `start_draft(project, "prd")` to initialize
2. Ask discovery questions (problem, users, goals)
3. Collaboratively fill sections
4. Use `update_draft` to save progress
5. Use `eval_draft` to check quality
6. Iterate until passing
7. Use `finalize_draft` to complete

### "What's the status of project X?"

Use `get_project_status(project)` to see:
- Required specs present
- Evaluations passing
- Approvals obtained
- Execution spec generated

### "Generate the technical spec"

1. Ensure PRD and UXD are approved
2. Use `synthesize(project, "trd")` to generate TRD
3. Review and refine with user
4. Use `run_eval(project, "trd")` to evaluate
5. Use `approve(project, "trd")` when ready

### "Export to SpecKit"

1. Ensure all required specs are approved
2. Use `reconcile(project)` to generate spec.md
3. Use `export(project, "speckit")` to create .specify/ structure

## Evaluation Criteria

Specs are evaluated on domain-specific rubrics:

- **PRD**: Problem definition, user stories, requirements, metrics
- **TRD**: Architecture clarity, API design, security, scalability
- **UXD**: User flows, accessibility, error handling

Findings have severity levels: critical, high, medium, low.
A spec passes with score >= 7.0 and no critical/high findings.

## Export Targets

| Target | Output | Use Case |
|--------|--------|----------|
| `speckit` | .specify/ directory | GitHub Spec-Kit workflows |
| `gsd` | PLAN.md, STATE.md | Get Shit Done methodology |
| `gastown` | formula.toml | Multi-agent formulas |
| `gascity` | city.toml | Multi-agent orchestration |

## Tips

- Always check project status before suggesting next steps
- Use `eval_draft` frequently during authoring to catch issues early
- Technical specs (TRD, IRD) should be synthesized, not manually written
- The CONSTITUTION.md constrains all generated specs
- Context sources ground technical synthesis in actual codebase data
