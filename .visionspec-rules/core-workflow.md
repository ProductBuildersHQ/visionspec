# VisionSpec Core Workflow

This document defines the orchestration rules for AI assistants working with VisionSpec projects. Use this workflow to guide users through spec creation, evaluation, and reconciliation.

## Trigger Pattern

Activate this workflow when the user says:

- "Using VisionSpec, [intent]"
- "Help me create specs for [project]"
- "I need to write a PRD/MRD/TRD"
- "What specs do I need for [project]?"

## Step 1: Framework Selection

When starting a new project, ask the user which methodology framework to use:

```
Which methodology framework would you like to use?

1. Big Tech - Product (recommended) - Best of all worlds (10 methodologies), MRD start
2. Big Tech - Feature - Best of all worlds (10 methodologies), OpportunitySpec start
3. AWS Working Backwards - Product - Start with MRD for new product lines
4. AWS Working Backwards - Feature - Start with OpportunitySpec for features
5. Google - OKRs, design docs, and experimentation
6. Stripe - API-first, developer experience focused
7. Shape Up - Fixed time, variable scope, betting on pitches (Basecamp)
8. Continuous Discovery - Weekly touchpoints, OST, assumption testing (Teresa Torres)
9. Lean Startup - Hypothesis-driven, build-measure-learn cycles
10. Design Thinking - Human-centered, empathy-first approach
11. Jobs to be Done - Focus on customer job statements and outcomes
```

Note: Big Tech profiles include Shape Up and Continuous Discovery practices as optional artifacts.
When using Big Tech, you can add pitches, hill charts, OSTs, and assumption maps as needed.

**Default**: If user doesn't specify or says "just start":
- For **new products**: use `big-tech-product` (comprehensive best practices, MRD start)
- For **features on existing products**: use `big-tech-feature` (comprehensive best practices, OpportunitySpec start)
- For **API/platform products**: use `big-tech-feature` or `stripe`

Initialize project with the selected profile:

```bash
visionspec init <project> --profile <framework>
# Examples:
visionspec init myproject --profile big-tech-product       # Best practices, MRD start (recommended)
visionspec init myproject --profile big-tech-feature       # Best practices, OpportunitySpec start
visionspec init myproject --profile aws-product            # Amazon Working Backwards, MRD start
visionspec init myproject --profile aws-feature            # Amazon Working Backwards, OpportunitySpec start
visionspec init myproject --profile google
visionspec init myproject --profile stripe
visionspec init myproject --profile shapeup                # Basecamp Shape Up, Pitch-based
visionspec init myproject --profile continuous-discovery   # Teresa Torres, OST + weekly touchpoints
visionspec init myproject --profile lean-startup
visionspec init myproject --profile design-thinking
visionspec init myproject --profile jtbd
```

## Framework-Specific Flows

### Big Tech Product (Recommended for New Products)

Combined best practices from AWS, Google, Stripe, Netflix, Spotify, and more. Starts with MRD.

```
MRD + OKRs (market problem with measurable objectives)
    ↓
Press Release (customer announcement - Amazon)
    ↓
FAQ (challenge assumptions, include developer FAQ - Amazon/Stripe)
    ↓
6-Pager (narrative stakeholder alignment - Amazon)
    ↓
PRD + API Contracts (requirements with non-goals and alternatives - Google/Stripe)
    ↓
UXD + DX (user experience with developer experience and accessibility - Stripe/Microsoft)
    ↓
Design Doc (TRD with explicit tradeoffs and reversibility - Google)
    ↓
Test Plan + Experiments (TPD with hypothesis-driven validation - Google/Spotify)
    ↓
spec.md (reconciled execution spec)
```

### Big Tech Feature (Recommended for Features)

Same practices as Big Tech Product, but starts with OpportunitySpec for feature-level opportunities.

```
OpportunitySpec + OKRs (12-box canvas with measurable objectives)
    ↓
Press Release (feature announcement - Amazon)
    ↓
FAQ (challenge assumptions, include developer FAQ - Amazon/Stripe)
    ↓
PRD + API Contracts (requirements with non-goals and alternatives - Google/Stripe)
    ↓
UXD + DX (user experience with developer experience and accessibility - Stripe/Microsoft)
    ↓
Design Doc (TRD with explicit tradeoffs and reversibility - Google)
    ↓
Test Plan + Experiments (TPD with hypothesis-driven validation - Google/Spotify)
    ↓
spec.md (reconciled execution spec)
```

Key practices integrated (both profiles):
- Customer obsession + Working backwards (Amazon)
- OKRs + 10x thinking + Explicit tradeoffs (Google)
- API-first + Developer experience (Stripe)
- Freedom & Responsibility + Context not control (Netflix)
- Bets not projects + Fail fast (Spotify)
- DRI (Apple) + Growth mindset (Microsoft)

### AWS Working Backwards

Customer-centric, vision-first approach. Start with the press release.

```
MRD (market problem)
    ↓
Press Release (customer announcement)
    ↓
FAQ (challenge assumptions)
    ↓
PRD (derive requirements)
    ↓
UXD (user experience)
    ↓
TRD → TPD → IRD (technical specs)
    ↓
spec.md (reconciled execution spec)
```

### Lean Startup

Hypothesis-driven, rapid validation cycles.

```
Hypothesis (what we believe)
    ↓
MRD (problem validation)
    ↓
MVP PRD (minimum viable product)
    ↓
Experiment Design (how to test)
    ↓
UXD (lean prototype)
    ↓
TRD (technical approach)
    ↓
spec.md → Build → Measure → Learn → Iterate
```

### Design Thinking

Human-centered, empathy-first approach.

```
Empathy Research (user interviews, observation)
    ↓
MRD (define the problem)
    ↓
Ideation (brainstorm solutions)
    ↓
PRD (selected solution)
    ↓
UXD (prototype)
    ↓
User Testing → Iterate
    ↓
TRD → TPD → IRD
    ↓
spec.md
```

### Jobs to be Done (JTBD)

Focus on customer jobs and desired outcomes.

```
Job Statements (what customers are trying to accomplish)
    ↓
MRD (job context and constraints)
    ↓
Outcome Expectations (how customers measure success)
    ↓
PRD (solution that addresses job)
    ↓
UXD (experience that enables job completion)
    ↓
TRD → TPD → IRD
    ↓
spec.md
```

### Google

OKRs, design docs, and data-driven experimentation.

```
OKRs (objectives and key results)
    ↓
MRD (problem and opportunity)
    ↓
Design Doc (technical approach with alternatives)
    ↓
RFC (request for comments, gather feedback)
    ↓
PRD + UXD (refined requirements)
    ↓
Experiment Design (A/B tests, metrics)
    ↓
TRD → TPD → IRD
    ↓
spec.md
```

### Stripe

API-first, developer experience focused.

```
API Contract (define the interface first)
    ↓
MRD (developer pain points)
    ↓
DX Review (developer experience critique)
    ↓
PRD (API requirements)
    ↓
UXD (docs, examples, error messages)
    ↓
TRD (implementation behind API)
    ↓
TPD (API contract tests, integration tests)
    ↓
IRD
    ↓
spec.md
```

### Shape Up (Basecamp)

Fixed time, variable scope. Bet on pitches, build in 6-week cycles.

```
Pitch (problem + appetite + solution shape)
    ↓
Betting Table (evaluate and bet on pitches)
    ↓
Scope Mapping (break pitch into scopes)
    ↓
Hill Chart Tracking (uphill = figuring out, downhill = execution)
    ↓
PRD (distilled requirements from pitch)
    ↓
UXD (scope-level wireframes)
    ↓
TRD (technical approach per scope)
    ↓
spec.md → Build → Ship (within 6-week cycle)
```

Key principles:
- **Appetite over estimates**: Set time budget first, shape work to fit
- **Pitches not backlogs**: Evaluate opportunities fresh each cycle
- **Hill charts**: Track work as "figuring out" vs "making it happen"
- **Circuit breaker**: If not done by deadline, it gets cut (no overruns)
- **Cool-down**: 2 weeks after each 6-week cycle for cleanup and exploration

### Continuous Discovery (Teresa Torres)

Weekly touchpoints, Opportunity Solution Trees, assumption testing.

```
Weekly Customer Touchpoints (interviews, observations)
    ↓
Opportunity Solution Tree (OST) (outcomes → opportunities → solutions)
    ↓
Discovery Snapshot (weekly summary of learnings)
    ↓
Assumption Mapping (identify and prioritize assumptions)
    ↓
Assumption Testing (experiments to validate/invalidate)
    ↓
PRD (validated requirements from tested assumptions)
    ↓
UXD (experience map, user flows)
    ↓
TRD (technical approach)
    ↓
spec.md → Build → Continuous learning
```

Key principles:
- **Weekly touchpoints**: Talk to customers every week, not just at project start
- **Story-based interviews**: Ask about recent experiences, not hypotheticals
- **Opportunity Solution Trees**: Map outcomes to opportunities to solutions
- **Assumption testing by type**: Test desirability, viability, feasibility, usability, ethical
- **Integrate with delivery**: Discovery happens continuously, not as a phase

### Combining Shape Up + Continuous Discovery

These frameworks complement each other:
- Use **Continuous Discovery** for ongoing opportunity identification
- Use **Shape Up** for execution cycles on validated opportunities

```
Weekly Customer Touchpoints (Continuous Discovery)
    ↓
Opportunity Solution Tree (Continuous Discovery)
    ↓
Pitch (Shape Up - shape validated opportunities)
    ↓
Betting Table (Shape Up - bet on top opportunities)
    ↓
6-Week Cycle (Shape Up - build with hill chart tracking)
    ↓
Weekly Touchpoints continue during build
    ↓
Ship → Measure → Feed back to OST
```

## Detecting Existing Framework

If the project already exists, check which framework is configured:

```bash
visionspec status -p <project>
# Look for "Profile:" in output
```

Or check `visionspec.yaml`:

```yaml
profile: aws-product  # or aws-feature, lean-startup, design-thinking, etc.
```

Continue with the configured framework's flow.

## Workflow Phases (AWS Working Backwards Default)

The default flow implements Amazon's Working Backwards methodology in five phases:

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

### Default Phases (AWS Working Backwards)

- [phases/01-discovery.md](phases/01-discovery.md) - MRD authoring
- [phases/02-vision.md](phases/02-vision.md) - Working Backwards flow
- [phases/03-requirements.md](phases/03-requirements.md) - PRD/UXD authoring
- [phases/04-technical.md](phases/04-technical.md) - Technical synthesis
- [phases/05-reconciliation.md](phases/05-reconciliation.md) - Final spec generation

### Framework-Specific Phases

- [frameworks/lean-startup/phases.md](frameworks/lean-startup/phases.md) - Build-Measure-Learn cycles
- [frameworks/design-thinking/phases.md](frameworks/design-thinking/phases.md) - Stanford d.school 5-stage process
- [frameworks/jtbd/phases.md](frameworks/jtbd/phases.md) - Jobs to be Done outcome-driven approach
- [frameworks/google/phases.md](frameworks/google/phases.md) - OKRs, Design Docs, experiments
- [frameworks/stripe/phases.md](frameworks/stripe/phases.md) - API-first, developer experience

### Gates

- [gates/evaluation.md](gates/evaluation.md) - Evaluation criteria
- [gates/approval.md](gates/approval.md) - Approval process
