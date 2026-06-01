# Phase 2: Vision (Working Backwards)

The Vision phase implements Amazon's Working Backwards methodology. Start with the customer announcement (Press Release), challenge it (FAQ), then derive requirements (PRD).

## Objective

Create approved Press Release, FAQ, and PRD that:

- Define the customer experience before implementation
- Surface gaps and concerns early
- Translate vision into testable requirements

## Entry Criteria

- MRD is approved
- Problem and audience are clearly defined

## Working Backwards Flow

```
MRD (approved)
    ↓
Press Release (synthesize, review, approve)
    ↓
FAQ (synthesize, review, approve)
    ↓
PRD (synthesize, review, approve)
```

**Why this order?**

The Press Release forces clarity on the customer experience. The FAQ challenges that vision and surfaces gaps. Only then is the PRD derived—grounded in validated vision rather than abstract feature lists.

## Workflow

### Step 1: Synthesize Press Release

```bash
visionspec synthesize press -p <project>
```

The Press Release should answer:

- What is the product/feature called?
- Who is it for?
- What problem does it solve?
- How does the customer benefit?
- What does the customer say about it?

### Step 2: Review Press Release

Review with user:

- [ ] Is this the announcement you'd want to read?
- [ ] Does it capture the essence of the value?
- [ ] Would this excite the target customer?

If not, iterate:

1. Identify what's wrong
2. Edit the press.md directly
3. Re-evaluate

### Step 3: Synthesize FAQ

```bash
visionspec synthesize faq -p <project>
```

The FAQ should include:

**Customer Questions**

- How does this work?
- What does it cost?
- How is this different from alternatives?

**Stakeholder Questions**

- What are the technical risks?
- What resources are needed?
- What's the timeline?

**Skeptic Questions**

- Why will this succeed when others failed?
- What could go wrong?
- What are we NOT doing?

### Step 4: Review FAQ

The FAQ is where bad ideas die. Review critically:

- [ ] Are hard questions answered honestly?
- [ ] Do answers reveal scope creep?
- [ ] Are risks acknowledged?

Use FAQ answers to refine the vision.

### Step 5: Synthesize PRD

```bash
visionspec synthesize prd -p <project>
```

The PRD translates vision into requirements:

- User stories from Press Release scenarios
- Acceptance criteria from FAQ answers
- Non-functional requirements from stakeholder concerns

### Step 6: Review PRD

- [ ] Are requirements testable?
- [ ] Is scope aligned with Press Release?
- [ ] Are edge cases covered (from FAQ)?

### Step 7: Evaluate and Approve

```bash
visionspec eval press -p <project>
visionspec eval faq -p <project>
visionspec eval prd -p <project>

visionspec approve press -p <project>
visionspec approve faq -p <project>
visionspec approve prd -p <project>
```

## Exit Criteria

- Press Release approved
- FAQ approved
- PRD approved
- All evaluation scores >= 7.0

## Next Phase

→ [Phase 3: Experience](03-requirements.md) (UXD)

## Anti-Patterns

- **Skipping the Press Release**: "We already know what we want." The Press Release forces clarity; don't skip it.
- **Softball FAQ**: Questions that are easy to answer. Push for hard questions.
- **PRD before Press**: Writing requirements before defining the vision leads to feature creep.
