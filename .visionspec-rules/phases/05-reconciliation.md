# Phase 5: Reconciliation

The Reconciliation phase combines all approved specs into a unified execution spec (`spec.md`) suitable for AI coding agents.

## Objective

Generate a `spec.md` that:

- Consolidates all requirements from source specs
- Resolves conflicts between specs
- Provides clear, actionable tasks for implementation

## Entry Criteria

- All required specs are approved:
  - MRD (source)
  - PRD (source or synthesized)
  - UXD (source)
  - TRD (synthesized)
  - TPD (synthesized)
  - IRD (synthesized)

## Workflow

### Step 1: Verify Readiness

```bash
visionspec status -p <project>
```

All gates should pass:

- [ ] Required specs present
- [ ] Evaluations passing
- [ ] Approvals obtained

### Step 2: Reconcile

```bash
visionspec reconcile -p <project>
```

The reconciliation process:

1. **Extracts requirements** from each spec
2. **Builds dependency graph** of requirements
3. **Detects conflicts** between specs
4. **Resolves tradeoffs** (may require user input)
5. **Generates unified spec.md**

### Step 3: Handle Conflicts

If conflicts are detected:

**Example Conflict: Performance vs. Feature**

```
PRD: "System must support 1000 concurrent users"
TRD: "Architecture supports 100 concurrent users with current design"
```

Resolution options:

1. Revise PRD requirement
2. Revise TRD architecture
3. Document as known limitation with roadmap

Ask user for decision and document in spec.md.

### Step 4: Review spec.md

The generated spec.md should contain:

**Header**

- Project name
- Version
- Generation timestamp
- Source specs used

**Requirements Summary**

- Functional requirements (from PRD)
- Non-functional requirements (from TRD)
- User experience requirements (from UXD)

**Implementation Tasks**

- Ordered by dependency
- Each task is actionable
- Tests are specified per task

**Acceptance Criteria**

- Derived from PRD acceptance criteria
- Maps to TPD test cases

### Step 5: Export to Target

```bash
# Choose target based on execution system
visionspec export speckit -p <project>
visionspec export gsd -p <project>
visionspec export aidlc -p <project>
```

## Exit Criteria

- spec.md exists
- No unresolved conflicts
- Ready for export to target system

## Export Targets

| Target | Use When |
|--------|----------|
| `speckit` | Using GitHub Spec-Kit for execution |
| `gsd` | Using Get Shit Done methodology |
| `gastown` | Multi-agent formula execution |
| `gascity` | Multi-agent orchestration |
| `aidlc` | AWS AI-DLC Workflows |

## Post-Reconciliation

After implementation begins, maintain alignment:

```bash
# Check implementation against spec
visionspec graph diff -p <project>

# Update current truth after shipping
visionspec graph snapshot -p <project>
```

## Anti-Patterns

- **Reconciling too early**: All specs must be approved first. Partial reconciliation leads to rework.
- **Ignoring conflicts**: Conflicts must be resolved, not papered over.
- **Spec drift**: Update spec.md when requirements change during implementation.
- **Export without reconcile**: Always reconcile before export to ensure consistency.
