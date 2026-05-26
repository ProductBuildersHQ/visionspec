# Design Thinking

Stanford d.school's Design Thinking methodology is a human-centered approach to innovation. It emphasizes empathy, experimentation, and iteration through five interconnected stages.

## The Flow

![Design Thinking Flow](../diagrams/design-thinking-flow.svg)

## Key Principles

1. **Human-Centered**: Deep empathy for users
2. **Bias Toward Action**: Prototype early and often
3. **Radical Collaboration**: Diverse perspectives
4. **Embrace Ambiguity**: Navigate uncertainty
5. **Iterate**: Learn through making

## The Five Stages

### 1. Empathize

Understand users through observation, engagement, and immersion.

**Methods**: User interviews, shadowing, journey mapping

### 2. Define

Synthesize findings into a clear point of view (POV).

**Methods**: POV statements, "How Might We" questions, empathy maps

### 3. Ideate

Generate a broad range of possible solutions.

**Methods**: Brainstorming, mind mapping, sketching

### 4. Prototype

Build quick, rough representations of ideas.

**Methods**: Paper prototypes, storyboards, role-playing

### 5. Test

Gather feedback and refine solutions.

**Methods**: User testing, A/B testing, feedback capture

## VisionSpec Mapping

| Design Thinking Artifact | VisionSpec Type | Purpose |
|--------------------------|-----------------|---------|
| Empathy Map | MRD | User understanding |
| HMW Questions | MRD | Problem framing |
| Journey Map | UXD | User experience flow |
| Prototype Spec | PRD | What to build and test |
| Test Plan | UXD | Validation methodology |

## Using the Design Thinking Profile

### Initialize a Project

```bash
multispec init user-onboarding --profile design-thinking
```

### Create Empathy Map (MRD)

```bash
multispec draft mrd -p user-onboarding
```

The Empathy Map template includes:

- **Says**: Quotes from users
- **Thinks**: What users are thinking
- **Does**: Observable behaviors
- **Feels**: Emotional state
- **Pains**: Frustrations and challenges
- **Gains**: Goals and desires

### Define "How Might We" Questions

The MRD template also includes:

- Problem synthesis
- POV statement (User... Needs... Because...)
- HMW questions (How Might We...)
- Prioritized design challenges

### Create Journey Map (UXD)

```bash
multispec draft uxd -p user-onboarding
```

The Journey Map template includes:

- User persona
- Journey stages
- Actions at each stage
- Thoughts and feelings
- Pain points
- Opportunities
- Touchpoints

### Create Prototype Spec (PRD)

```bash
multispec draft prd -p user-onboarding
```

The Prototype Spec template includes:

- HMW question being addressed
- Prototype type (paper, digital, physical)
- Fidelity level (low, medium, high)
- Key interactions to test
- What we're NOT testing
- Materials needed

### Create Test Plan (UXD)

The UXD can also include a Test Plan:

- Testing objectives
- Participant criteria
- Test scenarios
- Questions to ask
- Observation guide
- Feedback capture method

## Rubric Categories

### Empathy Map Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| User Clarity | 20% | Specific user defined |
| Empathy Depth | 25% | Rich understanding of user |
| Says/Thinks/Does/Feels | 20% | All quadrants covered |
| Pains Identified | 15% | Real frustrations found |
| Gains Identified | 10% | User goals understood |
| Insights | 10% | Non-obvious findings |

### HMW Questions Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| POV Clarity | 25% | Clear User/Needs/Because |
| HMW Quality | 30% | Questions open possibilities |
| Problem Scope | 20% | Neither too broad nor narrow |
| User-Centered | 15% | Focused on user needs |
| Actionability | 10% | Can guide ideation |

### Journey Map Evaluation (UXD)

| Category | Weight | Description |
|----------|--------|-------------|
| Journey Completeness | 20% | Full experience mapped |
| Emotional Mapping | 20% | Feelings at each stage |
| Pain Points | 20% | Frustrations identified |
| Opportunities | 20% | Design opportunities clear |
| Touchpoints | 10% | All interactions noted |
| Actionability | 10% | Guides design decisions |

### Prototype Spec Evaluation (PRD)

| Category | Weight | Description |
|----------|--------|-------------|
| HMW Alignment | 20% | Addresses specific question |
| Appropriate Fidelity | 20% | Right level for stage |
| Testable | 25% | Clear what to test |
| Focused | 15% | Not testing everything |
| Feasible | 10% | Can be built quickly |
| Iteration Ready | 10% | Easy to modify |

## Iteration Patterns

Design Thinking is non-linear. Common iteration patterns:

```
Test → Prototype → Test (refine)
Test → Define → Ideate (reframe)
Test → Empathize → Define (deeper understanding)
Ideate → Prototype → Ideate (expand possibilities)
```

## Example Workflow

```bash
# 1. Initialize project
multispec init checkout-flow --profile design-thinking

# 2. Empathize: Create empathy map from research
multispec draft mrd -p checkout-flow
# Include user interviews, observations
multispec eval mrd -p checkout-flow
multispec approve mrd -p checkout-flow

# 3. Define: HMW questions are in MRD
# Review and refine POV statement

# 4. Ideate: (happens outside VisionSpec - brainstorming)

# 5. Prototype: Create spec for what to build
multispec draft prd -p checkout-flow
multispec eval prd -p checkout-flow
multispec approve prd -p checkout-flow

# 6. Test: Create test plan
multispec draft uxd -p checkout-flow
multispec eval uxd -p checkout-flow
multispec approve uxd -p checkout-flow

# 7. Run tests, gather feedback...

# 8. Iterate: Return to earlier stage as needed
# If reframing: update MRD
# If refining: update PRD
```

## Brainstorming Rules

When facilitating ideation:

1. **Defer judgment** - No criticism during ideation
2. **Encourage wild ideas** - The unusual sparks innovation
3. **Build on others' ideas** - "Yes, and..."
4. **Stay focused on topic** - One challenge at a time
5. **One conversation at a time** - Listen actively
6. **Be visual** - Sketch ideas
7. **Go for quantity** - More ideas = more options

## Reference Materials

For deeper understanding of Design Thinking, see:

- [Stanford d.school](https://dschool.stanford.edu/)
- [Design Thinking Bootleg](https://dschool.stanford.edu/tools/design-thinking-bootleg)
- *Creative Confidence* by Tom and David Kelley
- Internal reference: `frameworks-internal/stanford-design-thinking/`
