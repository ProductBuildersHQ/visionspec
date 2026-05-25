# docs

Generate MkDocs-compatible documentation files.

## Usage

```bash
visionspec docs <subcommand> [flags]
```

## Description

The `docs` command generates markdown files for MkDocs integration. It creates index pages for projects and the main specs landing page.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `generate` | Generate all index.md files |
| `project` | Generate index.md for a specific project |

## docs generate

Generate index.md files for all projects and the specs landing page.

```bash
visionspec docs generate [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--with-graph` | bool | `false` | Include traceability metrics from spec graph |

**Output files:**

- `docs/specs/index.md` - Specs landing page with project listing
- `docs/specs/{project}/index.md` - Project index with status and specs table

**Example:**

```bash
# Generate all documentation
visionspec docs generate

# Include graph metrics
visionspec docs generate --with-graph
```

## docs project

Generate index.md for a specific project.

```bash
visionspec docs project <project-name> [flags]
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--with-graph` | bool | `false` | Include traceability metrics from spec graph |

**Example:**

```bash
# Generate index for specific project
visionspec docs project user-onboarding

# With graph metrics
visionspec docs project user-onboarding --with-graph
```

## Generated Content

### Project Index (`{project}/index.md`)

- Status badge (ready/not ready)
- Project description
- Specs table with links
- Traceability metrics (if `--with-graph`)

### Specs Landing (`specs/index.md`)

- Projects table with status and progress
- Links to CONSTITUTION.md and ROADMAP.md
- Cross-project metrics summary

## Graph Metrics

When using `--with-graph`, includes:

- Total requirements count
- User stories count
- Trace coverage percentage
- Conflict count

Requires running `visionspec graph extract` first.

## MkDocs Integration

The generated files integrate with your MkDocs site. Add to `mkdocs.yml`:

```yaml
nav:
  - Specifications:
      - Overview: specs/index.md
      - user-onboarding: specs/user-onboarding/index.md
      - payment-flow: specs/payment-flow/index.md
```

## Example Output

**specs/index.md:**

```markdown
# Specifications

| Project | Status | Progress |
|---------|--------|----------|
| [user-onboarding](user-onboarding/index.md) | Ready | 100% |
| [payment-flow](payment-flow/index.md) | In Progress | 60% |

## Repository Documents

- [CONSTITUTION.md](CONSTITUTION.md)
- [ROADMAP.md](ROADMAP.md)
```

**{project}/index.md:**

```markdown
# user-onboarding

Status: Ready

## Specifications

| Spec | Status | Last Eval |
|------|--------|-----------|
| [MRD](source/mrd.md) | Approved | 9.2/10 |
| [PRD](source/prd.md) | Approved | 8.5/10 |
| [TRD](technical/trd.md) | Approved | 8.8/10 |

## Traceability

- Requirements: 12
- User Stories: 8
- Trace Coverage: 92%
```

## See Also

- [status](status.md) - Check project status
- [graph](graph.md) - Extract and query requirement graphs
