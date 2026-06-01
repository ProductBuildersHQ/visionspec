# Design Thinking Framework Phases

Stanford d.school's human-centered design methodology.

## Overview

Design Thinking is a human-centered approach to innovation. It emphasizes empathy with users, rapid prototyping, and iterative refinement based on feedback.

```
Empathize (understand users)
    ↓
Define (frame the problem)
    ↓
Ideate (generate solutions)
    ↓
Prototype (build to think)
    ↓
Test (learn from users)
    ↓
Iterate
```

## Phase 1: Empathize

**Goal**: Deeply understand users, their needs, and their context.

### Research Methods

| Method | When to Use |
|--------|-------------|
| Interviews | Understanding motivations, emotions |
| Observation | Seeing actual behavior vs. reported |
| Shadowing | Day-in-the-life understanding |
| Journey mapping | End-to-end experience |
| Empathy mapping | Synthesizing user perspective |

### Workflow

1. Plan research (who, what, where)
2. Conduct interviews (5-8 users minimum)
3. Observe users in context
4. Document findings

### Empathy Map Template

```markdown
## User: [Persona Name]

### Says
- Direct quotes from interviews
- What they tell you they do

### Thinks
- What might they be thinking?
- What does this tell us about beliefs?

### Does
- Observable behaviors
- Actions and activities

### Feels
- Emotional state
- Frustrations, delights
```

### Output

Capture in MRD "Target Audience" section or separate empathy research document.

## Phase 2: Define (MRD)

**Goal**: Synthesize research into actionable problem statement.

### Point of View (POV)

```markdown
[User] needs [need] because [insight].
```

Example:
> "Busy parents need a way to quickly prepare healthy meals because they feel guilty about relying on fast food but don't have time for elaborate cooking."

### How Might We (HMW)

Convert POV into opportunity questions:

- How might we make healthy cooking faster?
- How might we reduce guilt about meal choices?
- How might we fit cooking into busy schedules?

### Workflow

```bash
visionspec create mrd -p <project>
```

MRD should include:

- [ ] User personas with empathy insights
- [ ] POV statement
- [ ] HMW questions
- [ ] Success criteria from user perspective

## Phase 3: Ideate

**Goal**: Generate many possible solutions without judgment.

### Ideation Rules

1. **Defer judgment** - No criticism during ideation
2. **Go for quantity** - More ideas = better ideas
3. **Build on ideas** - "Yes, and..." not "No, but..."
4. **Be visual** - Sketch, don't just describe
5. **Stay focused** - On the HMW question

### Ideation Techniques

| Technique | Description |
|-----------|-------------|
| Brainstorming | Classic group ideation |
| Worst possible idea | Invert to find good ideas |
| SCAMPER | Substitute, Combine, Adapt, Modify, Put to other use, Eliminate, Reverse |
| Analogous inspiration | How do other domains solve this? |
| Constraint removal | What if money/time were unlimited? |

### Workflow

1. Generate 50+ ideas (quantity over quality)
2. Cluster similar ideas
3. Vote on promising clusters
4. Select 3-5 concepts to prototype

### Output

Capture selected concepts in PRD "Proposed Solutions" section.

## Phase 4: Prototype (PRD + UXD)

**Goal**: Build to think, not to ship.

### Prototype Principles

- **Low fidelity first** - Paper before pixels
- **Fail fast** - Quick builds, quick learning
- **Multiple options** - Test 2-3 concepts
- **User-testable** - Real enough to get feedback

### Prototype Spectrum

| Fidelity | Time | Purpose |
|----------|------|---------|
| Paper sketches | Minutes | Concept direction |
| Storyboards | Hours | User journey |
| Clickable mockup | Day | Interaction flow |
| Looks-like prototype | Days | Visual/emotional response |
| Works-like prototype | Week | Functional feedback |
| Integrated prototype | Weeks | Full experience |

### Workflow

```bash
visionspec create prd -p <project>
visionspec create uxd -p <project>
```

PRD focuses on:

- [ ] Core value proposition
- [ ] Key features to test
- [ ] What we're learning (not building)

UXD focuses on:

- [ ] User flows for testing
- [ ] Prototype specifications
- [ ] Test script and scenarios

## Phase 5: Test

**Goal**: Learn from real users interacting with prototypes.

### Testing Guidelines

1. **Show, don't tell** - Let users experience, then react
2. **Create experiences** - Test in realistic context
3. **Observe and listen** - Don't defend or explain
4. **Ask "why"** - Dig into reactions
5. **Test to learn** - Not to validate ego

### Test Protocol

```markdown
## Test Session: [Date]

### Setup
- Prototype version: [X]
- User profile: [Persona]
- Context: [Where/when testing]

### Tasks
1. [Ask user to accomplish goal]
2. [Observe without helping]
3. [Note where they struggle]

### Questions
- What was that like for you?
- What was confusing?
- What would you expect to happen?
- Would you use this? Why/why not?

### Observations
- [What did they do?]
- [What surprised you?]
- [What did they miss?]
```

### Synthesis

After testing sessions:

- [ ] Identify patterns across users
- [ ] Distinguish "nice to have" from critical
- [ ] Update problem understanding
- [ ] Decide: iterate, pivot, or proceed

## Phase 6: Iterate

**Goal**: Refine based on learning, then build.

### Iteration Decision Matrix

| Feedback | Action |
|----------|--------|
| Core value unclear | Return to Define |
| Solution direction wrong | Return to Ideate |
| Interaction problems | Refine Prototype |
| Positive validation | Proceed to Technical |

### Proceeding to Build

When user testing validates the solution:

```bash
visionspec context gather -p <project>
visionspec synthesize trd -p <project>
visionspec synthesize tpd -p <project>
visionspec reconcile -p <project>
```

## Design Thinking Gates

| Gate | Criteria |
|------|----------|
| Empathy complete | 5+ user interviews, empathy maps |
| Problem defined | Clear POV and HMW statements |
| Ideas generated | 50+ ideas, 3-5 selected concepts |
| Prototype tested | 5+ users tested, patterns identified |
| Validation achieved | Users demonstrate need satisfaction |

## See Also

- [Stanford d.school](https://dschool.stanford.edu/)
- [IDEO Design Thinking](https://designthinking.ideo.com/)
- [The Design of Everyday Things](https://www.nngroup.com/books/design-everyday-things-revised/) - Don Norman
