# Lean Startup

Eric Ries' Lean Startup methodology emphasizes rapid experimentation and validated learning. The Build-Measure-Learn feedback loop is the core engine for turning ideas into sustainable businesses.

## The Flow

![Lean Startup Build-Measure-Learn Flow](../diagrams/lean-startup-flow.svg)

## Key Principles

1. **Validated Learning**: Progress through scientific testing of hypotheses
2. **Build-Measure-Learn**: Rapid feedback loops
3. **Minimum Viable Product (MVP)**: Learn with least effort
4. **Pivot or Persevere**: Data-driven decision to change direction
5. **Innovation Accounting**: Measure progress in uncertain conditions

## VisionSpec Mapping

| Lean Artifact | VisionSpec Type | Purpose |
|---------------|-----------------|---------|
| Hypothesis Document | MRD | Assumptions to validate |
| MVP PRD | PRD | Minimum feature set for learning |
| Experiment Design | UXD | Test methodology |
| Pivot/Persevere Decision | Narrative | Learning documentation |

## Using the Lean Startup Profile

### Initialize a Project

```bash
multispec init my-hypothesis --profile lean-startup
```

### Create Hypothesis Document (MRD)

```bash
multispec draft mrd -p my-hypothesis
```

The Hypothesis Document template includes:

- Problem hypothesis
- Solution hypothesis
- Customer segment hypothesis
- Riskiest assumptions (ordered by risk)
- Leap of faith assumptions
- Success/failure criteria

### Create MVP PRD

```bash
multispec draft prd -p my-hypothesis
```

The MVP PRD template includes:

- Learning objectives (what we want to learn)
- MVP scope (minimum features)
- What's explicitly excluded
- Build vs. fake (Wizard of Oz options)
- Timeline constraints

### Create Experiment Design (UXD)

```bash
multispec draft uxd -p my-hypothesis
```

The Experiment Design template includes:

- Hypothesis being tested
- Experiment type (A/B, landing page, concierge, etc.)
- Success metrics (actionable, not vanity)
- Sample size and duration
- Data collection method
- Analysis plan

### Document Pivot/Persevere Decision

```bash
multispec draft narrative-1p -p my-hypothesis
```

The Pivot Document includes:

- What we learned
- Data supporting conclusions
- Pivot type (if applicable)
- New hypotheses to test
- Resources needed

## Rubric Categories

### Hypothesis Document Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Hypothesis Clarity | 25% | Clear, falsifiable statements |
| Assumption Identification | 20% | Key assumptions explicit |
| Risk Prioritization | 20% | Riskiest assumptions first |
| Testability | 15% | Hypotheses can be tested |
| Success Criteria | 15% | Clear pass/fail thresholds |
| Customer Definition | 5% | Specific customer segment |

### MVP PRD Evaluation

| Category | Weight | Description |
|----------|--------|-------------|
| Learning Focus | 25% | Learning objectives clear |
| Minimum Scope | 25% | Truly minimum for learning |
| Exclusions Clear | 15% | What's NOT included |
| Build vs. Fake | 15% | Considered Wizard of Oz |
| Time-Boxed | 10% | Iteration has deadline |
| Measurable | 10% | Can measure outcomes |

### Experiment Design Evaluation (UXD)

| Category | Weight | Description |
|----------|--------|-------------|
| Hypothesis Clarity | 20% | What exactly are we testing |
| Actionable Metrics | 20% | Not vanity metrics |
| Experiment Validity | 20% | Will results be meaningful |
| Sample Size | 15% | Statistically valid |
| Data Collection | 15% | How we'll gather data |
| Analysis Plan | 10% | How we'll interpret |

## Experiment Types

The Lean Startup profile supports various experiment types:

| Type | When to Use | Effort |
|------|-------------|--------|
| **Smoke Test** | Test demand before building | Low |
| **Landing Page** | Test value proposition | Low |
| **Concierge** | Manual service, learn process | Medium |
| **Wizard of Oz** | Fake automation, real experience | Medium |
| **A/B Test** | Compare specific variations | Medium |
| **Prototype** | Test usability, not demand | Medium |
| **Cohort Analysis** | Measure behavior over time | High |

## Example Workflow

```bash
# 1. Initialize project
multispec init food-delivery --profile lean-startup

# 2. Document hypotheses
multispec draft mrd -p food-delivery
# Focus on riskiest assumptions
multispec eval mrd -p food-delivery
multispec approve mrd -p food-delivery

# 3. Define MVP
multispec draft prd -p food-delivery
# Minimum features for learning
multispec eval prd -p food-delivery
multispec approve prd -p food-delivery

# 4. Design experiment
multispec draft uxd -p food-delivery
# How we'll test the hypothesis
multispec eval uxd -p food-delivery
multispec approve uxd -p food-delivery

# 5. Run experiment, collect data...

# 6. Document pivot/persevere
multispec draft narrative-1p -p food-delivery

# 7. If pivot, start new iteration
multispec init food-delivery-v2 --profile lean-startup
```

## Metrics That Matter

### Avoid Vanity Metrics

| Vanity Metric | Actionable Alternative |
|---------------|----------------------|
| Total users | Active users (DAU/MAU) |
| Page views | Conversion rate |
| Downloads | Retention rate |
| Registered users | Paying customers |
| Followers | Engagement rate |

### The One Metric That Matters (OMTM)

Choose one metric that best represents current learning focus:

- **Problem validation**: Customer interview conversion
- **Solution validation**: Sign-up rate
- **Product validation**: Activation rate
- **Revenue validation**: Paying conversion

## Innovation Accounting

Track progress through three phases:

1. **Establish baseline**: First MVP data
2. **Tune the engine**: Improve toward ideal
3. **Pivot or persevere**: Decide based on data

## Reference Materials

For deeper understanding of Lean Startup methodology, see:

- [The Lean Startup](https://theleanstartup.com/)
- *The Lean Startup* by Eric Ries
- *Running Lean* by Ash Maurya
- Internal reference: `frameworks-internal/lean-startup/`
