# version

Manage spec version history.

## Usage

```bash
visionspec version <subcommand> [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `create` | Create a new version from current spec |
| `list` | List all versions for a spec |
| `show` | Show a specific version |
| `diff` | Compare versions |
| `revert` | Revert to a previous version |

## version create

Create a new version snapshot of a spec.

```bash
visionspec version create <spec-type> -p <project> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--message` | `-m` | Version message describing changes |
| `--project` | `-p` | Project name |

**Examples:**

```bash
# Create version with message
visionspec version create prd -p user-onboarding -m "Added authentication requirements"

# Create version for TRD
visionspec version create trd -p user-onboarding -m "Updated API design"
```

## version list

List all versions for a spec.

```bash
visionspec version list <spec-type> -p <project>
```

**Examples:**

```bash
# List PRD versions
visionspec version list prd -p user-onboarding
```

**Output:**

```
Versions for prd (user-onboarding):

  v3  2024-01-15 14:32  Added authentication requirements
  v2  2024-01-10 09:15  Updated user stories
  v1  2024-01-08 11:00  Initial version
```

## version show

Show a specific version of a spec.

```bash
visionspec version show <spec-type> <version> -p <project>
```

**Examples:**

```bash
# Show version 2 of PRD
visionspec version show prd v2 -p user-onboarding

# Show latest version
visionspec version show prd latest -p user-onboarding
```

## version diff

Compare two versions of a spec.

```bash
visionspec version diff <spec-type> <version1> <version2> -p <project>
```

**Examples:**

```bash
# Compare v1 and v2
visionspec version diff prd v1 v2 -p user-onboarding

# Compare current with previous version
visionspec version diff prd v2 current -p user-onboarding
```

**Output:**

Shows unified diff format highlighting additions, deletions, and changes.

## version revert

Revert a spec to a previous version.

```bash
visionspec version revert <spec-type> <version> -p <project>
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--force` | Revert without confirmation |

**Examples:**

```bash
# Revert to version 2
visionspec version revert prd v2 -p user-onboarding

# Force revert without confirmation
visionspec version revert prd v2 -p user-onboarding --force
```

## Version Storage

Versions are stored in the project's `.versions/` directory:

```
docs/specs/user-onboarding/
├── source/
│   └── prd.md          # Current version
└── .versions/
    └── prd/
        ├── v1.md
        ├── v2.md
        └── versions.json
```

## See Also

- [status](status.md) - Show project status
- [approve](approve.md) - Approve specs for reconciliation
