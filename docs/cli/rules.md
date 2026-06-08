# rules

Manage workflow rules for AI assistant orchestration.

## Usage

```bash
visionspec rules <subcommand>
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `list` | List available workflow rules |
| `export` | Export workflow rules to a directory |

## Overview

Workflow rules guide AI assistants (Claude Code, Kiro, Cursor) through the VisionSpec specification workflow.

Rules provide:

- Trigger patterns for activating VisionSpec workflows
- Phase-by-phase guidance for spec creation
- Framework-specific flows (AWS, Lean Startup, Design Thinking, etc.)
- Evaluation and approval gates

## rules list

List all available workflow rules.

```bash
visionspec rules list
```

**Output:**

```
Available workflow rules:

  core-workflow.md
  phases/discovery.md
  phases/authoring.md
  gates/evaluation.md
  frameworks/aws.md
  ...

Export rules to your project with: visionspec rules export
```

## rules export

Export VisionSpec workflow rules to your project directory.

```bash
visionspec rules export [output-dir]
```

**Arguments:**

| Argument | Default | Description |
|----------|---------|-------------|
| `output-dir` | `.visionspec-rules` | Directory to export rules to |

**Examples:**

```bash
# Export to default location (.visionspec-rules)
visionspec rules export

# Export to custom directory
visionspec rules export ./my-rules
```

**Output:**

```
✓ Exported 15 rule files to .visionspec-rules

Contents:
  core-workflow.md         - Main orchestration rules
  phases/                  - Phase-by-phase guidance
  gates/                   - Evaluation and approval gates
  frameworks/              - Framework-specific flows

To use with Claude Code, add to your CLAUDE.md:
  See .visionspec-rules/ for VisionSpec workflow guidance.
```

## Rule Structure

After export, your project will contain:

```
.visionspec-rules/
├── core-workflow.md         # Main orchestration rules
├── phases/                  # Phase-by-phase guidance
│   ├── discovery.md         # Discovery phase
│   ├── authoring.md         # Authoring phase
│   └── synthesis.md         # Synthesis phase
├── gates/                   # Evaluation and approval gates
│   ├── evaluation.md        # Evaluation criteria
│   └── approval.md          # Approval workflow
└── frameworks/              # Framework-specific flows
    ├── aws.md               # Amazon Working Backwards
    ├── lean-startup.md      # Lean Startup
    └── design-thinking.md   # Stanford Design Thinking
```

## AI Assistant Integration

### Claude Code

Reference the rules in your `CLAUDE.md`:

```markdown
## Workflow Rules

For detailed orchestration rules, see `.visionspec-rules/`:

- [core-workflow.md](.visionspec-rules/core-workflow.md) - Main orchestration rules
- [phases/](.visionspec-rules/phases/) - Phase-by-phase guidance
- [gates/](.visionspec-rules/gates/) - Evaluation and approval gates

**Trigger Pattern**: Activate the workflow when user says "Using VisionSpec, [intent]".
```

### AWS Kiro

Copy rules to `.kiro/steering/`:

```bash
visionspec rules export .kiro/steering
```

### Cursor

Copy rules to `.cursor/rules/`:

```bash
visionspec rules export .cursor/rules
```

## See Also

- [profiles](profiles.md) - Configuration profiles for different project types
- [init](init.md) - Initialize a new project with rules
