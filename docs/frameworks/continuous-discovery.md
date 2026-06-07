# Continuous Discovery

Continuous Discovery is Teresa Torres's framework for integrating customer research into the daily product development process through weekly touchpoints, opportunity solution trees, and assumption testing.

## When to Use

Use Continuous Discovery when:

- You want to stay connected to customers continuously, not just during research phases
- You need to reduce risk before building by testing assumptions
- You want a structured approach to discovery that integrates with delivery
- You're working with a product trio (PM, Design, Engineering)

## Core Concepts

| Concept | Description |
|---------|-------------|
| **Weekly Touchpoints** | Talk to customers every week |
| **Opportunity Solution Tree** | Visualize path from outcome to solutions |
| **Assumption Testing** | Test risky assumptions before building |
| **Story-Based Interviews** | Collect stories (past behavior), not opinions |
| **Compare and Contrast** | Test multiple solutions in parallel |

## The Continuous Discovery Flow

```
1. OUTCOME
   Define a clear, measurable outcome
       ↓
2. OPPORTUNITIES
   Discover customer needs through weekly interviews
       ↓
3. SOLUTIONS
   Generate multiple solution ideas per opportunity
       ↓
4. ASSUMPTIONS
   Map assumptions by type, prioritize by risk
       ↓
5. EXPERIMENTS
   Test riskiest assumptions with small experiments
       ↓
6. LEARNINGS
   Update OST based on what you learn
       ↓
   (repeat weekly)
```

## Key Artifacts

### Opportunity Solution Tree (OST)

Visual map from outcome to opportunities to solutions to experiments.

```
[Outcome: Increase activation to 60%]
         |
         +-- [Opportunity: Users don't understand value]
         |       |
         |       +-- [Solution: Interactive tutorial]
         |       |       +-- [Experiment: Prototype test]
         |       |
         |       +-- [Solution: Personalized onboarding]
         |               +-- [Experiment: A/B test]
         |
         +-- [Opportunity: Setup is too complex]
                 |
                 +-- [Solution: Wizard flow]
```

### Discovery Snapshot

Weekly summary of discovery activities:

- Interviews conducted
- Stories collected
- Opportunities discovered/updated
- Assumptions tested
- Key learnings
- Decisions made

### Assumption Map

Assumptions organized by type:

| Type | Question |
|------|----------|
| **Desirability** | Will customers want this? |
| **Viability** | Will it work for the business? |
| **Feasibility** | Can we build this? |
| **Usability** | Can users figure it out? |
| **Ethical** | Should we build this? |

Prioritize: High importance + Low confidence = Test first

## Using the Continuous Discovery Profile

### Initialize a Project

```bash
multispec init my-feature --profile continuous-discovery
```

### Create an Opportunity Solution Tree

```bash
multispec draft ost -p my-feature
```

### Create Weekly Discovery Snapshot

```bash
multispec draft discovery-snapshot -p my-feature
```

### Create Assumption Map

```bash
multispec draft assumption-map -p my-feature
```

### Evaluate Discovery Snapshot

```bash
multispec eval discovery-snapshot -p my-feature
```

## Story-Based Interviews

In Continuous Discovery, we collect stories, not opinions.

| Bad (Opinion) | Good (Story) |
|---------------|--------------|
| "What do you want?" | "Tell me about the last time you..." |
| "Would you use this?" | "What happened when you tried to..." |
| "Is this important to you?" | "Walk me through how you..." |

### Story Structure

| Element | Question |
|---------|----------|
| **Situation** | What was the context? |
| **Behavior** | What did you do? |
| **Outcome** | What happened? |
| **Emotions** | How did you feel? |
| **Pain Points** | What was frustrating? |
| **Workaround** | How did you work around it? |

## Assumption Testing

### Risk Matrix

```
                    CONFIDENCE
                Low         High
           ┌───────────┬───────────┐
     High  │  TEST     │  MONITOR  │
IMPORTANCE │  FIRST    │           │
           ├───────────┼───────────┤
     Low   │  TEST     │  SKIP     │
           │  LATER    │           │
           └───────────┴───────────┘
```

### Test Methods

| Method | Best For | Duration |
|--------|----------|----------|
| Prototype test | Usability | 1 week |
| Fake door | Desirability | 1-2 weeks |
| Survey | Quantifying | 1 week |
| Interview | Exploring | Ongoing |
| Data analysis | Validation | 1 week |
| A/B test | Comparing | 2+ weeks |

## Example Workflow

```bash
# 1. Initialize project
multispec init user-activation --profile continuous-discovery

# 2. Create OST with outcome
multispec draft ost -p user-activation
# Set outcome: "Increase user activation rate to 60%"

# 3. Week 1: Start discovery
multispec draft discovery-snapshot -p user-activation
# Conduct interviews, capture stories, identify opportunities

# 4. Update OST with opportunities
# Add opportunities discovered from interviews

# 5. Generate solutions for top opportunity
# Add 3+ solution ideas to the OST

# 6. Map assumptions for top solution
multispec draft assumption-map -p user-activation

# 7. Test riskiest assumption
# Run small experiment, capture results

# 8. Week 2: Continue discovery
multispec draft discovery-snapshot -p user-activation
# More interviews, more learning, update OST

# 9. Evaluate discovery cadence
multispec eval discovery-snapshot -p user-activation
# Ensure weekly touchpoints are happening

# 10. Once validated, synthesize PRD
multispec synthesize prd -p user-activation
```

## Rubric Categories

### Discovery Snapshot Evaluation

| Category | Weight | Description |
|----------|--------|-------------|
| Weekly Touchpoints | 25% | At least one interview per week |
| Story Capture | 25% | Stories with situation-behavior-outcome |
| Opportunity Tracking | 20% | Opportunities discovered and prioritized |
| Assumption Testing | 20% | Tests running with clear hypotheses |
| Decision Making | 10% | Evidence-driven decisions |

### Assumption Map Evaluation

| Category | Weight | Description |
|----------|--------|-------------|
| Assumption Coverage | 25% | All types considered |
| Risk Assessment | 25% | Importance and confidence rated |
| Prioritization | 20% | Riskiest assumptions identified |
| Test Planning | 20% | Tests designed with success criteria |
| Document Quality | 10% | Clear and actionable |

## Principles

1. **Weekly Touchpoints** - Talk to customers every week
2. **Outcome-Driven** - Start with a clear, measurable outcome
3. **Opportunity Mapping** - Map customer needs, pains, desires
4. **Solution Trees** - Visualize the path from outcome to solutions
5. **Assumption Testing** - Test assumptions before building
6. **Small Experiments** - Run small, fast experiments
7. **Story-Based Interviews** - Collect stories, not opinions
8. **Collaborative Discovery** - Product trio works together
9. **Compare and Contrast** - Test multiple solutions in parallel
10. **Iterate on Learning** - Let learnings drive next experiment

## References

- [Continuous Discovery Habits (book)](https://www.producttalk.org/continuous-discovery-habits/)
- [Product Talk](https://www.producttalk.org/)
- [Teresa Torres](https://www.linkedin.com/in/teresatorres/)
- [Opportunity Solution Tree template](https://www.producttalk.org/opportunity-solution-tree/)
