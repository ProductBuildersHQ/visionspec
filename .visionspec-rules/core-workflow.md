# VisionSpec Core Workflow

This document defines the orchestration rules for AI assistants working with VisionSpec projects. Use this workflow to guide users through spec creation, evaluation, and reconciliation.

## Trigger Pattern

Activate this workflow when the user says:

- "Using VisionSpec, [intent]"
- "Help me create specs for [project]"
- "I need to write a PRD/MRD/TRD"
- "What specs do I need for [project]?"

## Workflow Phases

VisionSpec implements Amazon's Working Backwards methodology in five phases:

```
Phase 1: DISCOVERY
    MRD (human-authored)
        ↓
Phase 2: VISION (Working Backwards)
    Press → FAQ → PRD (synthesized, reviewable)
        ↓
Phase 3: EXPERIENCE
    UXD (human-authored)
        ↓
Phase 4: TECHNICAL
    TRD → TPD → IRD (synthesized, reviewable)
        ↓
Phase 5: RECONCILIATION
    spec.md (unified execution spec)
```

## Phase Gates

Each phase has approval gates. Do NOT proceed to the next phase until gates pass.

| Phase | Gate | Criteria |
|-------|------|----------|
| Discovery | MRD Approved | Problem is clear, audience defined, business goals stated |
| Vision | Press+FAQ+PRD Approved | Vision is compelling, FAQ addresses concerns, PRD is testable |
| Experience | UXD Approved | User journeys defined, interactions specified |
| Technical | TRD+TPD+IRD Approved | Architecture is sound, tests planned, infra specified |
| Reconciliation | spec.md Generated | All specs reconciled, conflicts resolved |

## Orchestration Rules

### Rule 1: Always Check Status First

Before any operation, check project status:

```
visionspec status -p <project>
```

This tells you:
- Which specs exist
- Which specs are approved
- Which gates are passing
- What the next step should be

### Rule 2: Follow the Phase Order

Never skip phases. If user asks to write TRD but MRD doesn't exist:

1. Explain the Working Backwards methodology
2. Offer to start with MRD instead
3. Show the dependency chain

### Rule 3: Evaluate Before Approving

Every spec must be evaluated before approval:

```
visionspec eval <type> -p <project>
```

Only approve if:
- Score >= 7.0
- No critical findings
- No high findings (or user acknowledges them)

### Rule 4: Synthesize Don't Write

Technical specs (TRD, TPD, IRD) and GTM specs (Press, FAQ) should be synthesized:

```
visionspec synthesize <type> -p <project>
```

Human review and refinement follows synthesis.

### Rule 5: Gather Context for Technical Specs

Before synthesizing TRD/TPD/IRD, gather codebase context:

```
visionspec context gather -p <project>
```

This grounds technical decisions in the actual system.

## Common User Intents

### "I have an idea for a feature"

→ Start Phase 1 (Discovery)
→ Guide user through MRD creation
→ Use `skills/author-mrd/` workflow

### "I need to write requirements"

→ Check if MRD exists and is approved
→ If yes: Start Phase 2, synthesize Press/FAQ, then PRD
→ If no: Explain Working Backwards, start with MRD

### "Generate the technical spec"

→ Check if PRD and UXD are approved
→ If yes: Gather context, synthesize TRD
→ If no: Guide user to complete prerequisites

### "Export to [target]"

→ Check if spec.md exists
→ If no: Run reconciliation first
→ Export to target (speckit, gsd, gastown, gascity, aidlc)

## Error Recovery

### Evaluation Fails

1. Show findings to user
2. Identify critical/high issues
3. Propose specific fixes
4. Re-evaluate after changes

### Synthesis Produces Poor Output

1. Check if source specs are complete
2. Verify context was gathered (for technical specs)
3. Regenerate with user guidance
4. Manual editing is acceptable

### Conflict During Reconciliation

1. Show conflicting requirements
2. Present tradeoff options
3. Get user decision
4. Document resolution in spec.md

## See Also

- [phases/01-discovery.md](phases/01-discovery.md) - MRD authoring
- [phases/02-vision.md](phases/02-vision.md) - Working Backwards flow
- [phases/03-requirements.md](phases/03-requirements.md) - PRD/UXD authoring
- [phases/04-technical.md](phases/04-technical.md) - Technical synthesis
- [phases/05-reconciliation.md](phases/05-reconciliation.md) - Final spec generation
- [gates/evaluation.md](gates/evaluation.md) - Evaluation criteria
- [gates/approval.md](gates/approval.md) - Approval process
