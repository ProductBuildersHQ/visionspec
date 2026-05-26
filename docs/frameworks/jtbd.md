# Jobs to be Done (JTBD)

Clayton Christensen's Jobs to be Done framework focuses on understanding the underlying motivations behind customer behavior. People don't buy products—they "hire" them to get a job done.

## The Flow

![Jobs to be Done Flow](../diagrams/jtbd-flow.svg)

## Key Principles

1. **Jobs Are Stable**: Customer objectives persist despite changing products
2. **Hire and Fire**: Customers "hire" products to do jobs
3. **Three Dimensions**: Functional, emotional, and social jobs
4. **Outcome-Focused**: Measure success by desired outcomes
5. **Context Matters**: Same person, different situations, different jobs

## Core Concept

> "People don't want a quarter-inch drill. They want a quarter-inch hole."
> — Theodore Levitt

People buy products and services to get a **job** done. A job is the task, goal, or objective a customer is trying to accomplish in a specific situation.

## VisionSpec Mapping

| JTBD Artifact | VisionSpec Type | Purpose |
|---------------|-----------------|---------|
| Job Statement & Context | MRD | The job to be done |
| Outcome-Driven Requirements | PRD | Desired outcomes |
| Job Map | UXD | Stages of the job |
| Solution Architecture | TRD | How we address outcomes |

## Using the JTBD Profile

### Initialize a Project

```bash
multispec init home-organization --profile jtbd
```

### Create Job Statement (MRD)

```bash
multispec draft mrd -p home-organization
```

The Job Statement template includes:

**Job Executor**

- Who is performing the job?
- What is their role/context?

**Core Functional Job**

- What are they trying to accomplish?
- Statement in solution-neutral language

**Job Context**

- When does this job arise?
- What triggers the need?
- What constraints exist?

**Related Jobs**

- Emotional jobs (how they want to feel)
- Social jobs (how they want to be perceived)
- Related functional jobs

**Current Solutions**

- What do they "hire" today?
- Why do they "fire" existing solutions?

### Create Outcome-Driven Requirements (PRD)

```bash
multispec draft prd -p home-organization
```

The Outcome-Driven Requirements template includes:

**Desired Outcomes**

Format: [Direction] + [Metric] + [Object of Control] + [Context]

Examples:

- Minimize the time it takes to find items when needed
- Minimize the likelihood of forgetting where something is stored
- Increase the confidence that items are stored safely

**Outcome Prioritization**

| Outcome | Importance | Satisfaction | Opportunity |
|---------|------------|--------------|-------------|
| ... | 1-10 | 1-10 | I + (I - S) |

**Underserved Outcomes**

- High importance, low satisfaction = opportunity

**Overserved Outcomes**

- Lower importance, high satisfaction = cost reduction opportunity

### Create Job Map (UXD)

```bash
multispec draft uxd -p home-organization
```

The Job Map template follows the Universal Job Map:

1. **Define**: What aspects must be defined before starting?
2. **Locate**: What inputs are needed?
3. **Prepare**: What preparation is required?
4. **Confirm**: What must be verified before proceeding?
5. **Execute**: What is the core action?
6. **Monitor**: What must be watched during execution?
7. **Modify**: What adjustments might be needed?
8. **Conclude**: How is the job completed?

For each stage:

- Actions taken
- Desired outcomes
- Current pain points
- Opportunity areas

## Rubric Categories

### Job Statement Evaluation (MRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Job Executor Clarity | 15% | Specific person defined |
| Job Definition | 25% | Solution-neutral, clear job |
| Context Definition | 20% | When/where job arises |
| Three Dimensions | 15% | Functional, emotional, social |
| Current Solutions | 15% | What they hire today |
| Competition | 10% | All job competitors identified |

### Outcome Requirements Evaluation (PRD)

| Category | Weight | Description |
|----------|--------|-------------|
| Outcome Format | 20% | Proper outcome statements |
| Outcome Coverage | 20% | All stages of job covered |
| Quantification | 20% | Importance/satisfaction measured |
| Opportunity Identification | 20% | Underserved outcomes found |
| Prioritization | 10% | Clear priority order |
| Actionability | 10% | Guides solution design |

### Job Map Evaluation (UXD)

| Category | Weight | Description |
|----------|--------|-------------|
| Stage Completeness | 20% | All 8 stages considered |
| Outcome Mapping | 25% | Outcomes at each stage |
| Pain Point Identification | 20% | Frustrations documented |
| Opportunity Areas | 20% | Clear design opportunities |
| User Validation | 10% | Based on real user data |
| Actionability | 5% | Guides solution design |

## The Opportunity Algorithm

Calculate opportunity score for each desired outcome:

```
Opportunity = Importance + (Importance - Satisfaction)
```

| Score | Meaning |
|-------|---------|
| > 15 | Extreme opportunity |
| 12-15 | High opportunity |
| 10-12 | Moderate opportunity |
| < 10 | Low opportunity or overserved |

## Example: The Milkshake Story

The famous McDonald's milkshake study illustrates JTBD:

**Morning Commuters**

- **Job**: Make my boring commute more interesting
- **Context**: Alone, driving, before 8am
- **Competition**: Bagels (crumbs), bananas (too quick), donuts (sticky)
- **Why milkshake wins**: Thick (lasts long), one-handed, entertaining

**Afternoon Parents**

- **Job**: Be a good parent (say yes to something)
- **Context**: After school, with kids
- **Competition**: Ice cream, cookies, toys
- **Why milkshake loses**: Too big, takes too long

**Same product, different jobs, different requirements.**

## Example Workflow

```bash
# 1. Initialize project
multispec init laundry-service --profile jtbd

# 2. Define the job
multispec draft mrd -p laundry-service
# Conduct Switch interviews
# Identify job executor, context, competing solutions
multispec eval mrd -p laundry-service
multispec approve mrd -p laundry-service

# 3. Document desired outcomes
multispec draft prd -p laundry-service
# Survey importance and satisfaction
# Calculate opportunity scores
multispec eval prd -p laundry-service
multispec approve prd -p laundry-service

# 4. Map the job
multispec draft uxd -p laundry-service
# Universal Job Map for laundry
# Outcomes at each stage
multispec eval uxd -p laundry-service
multispec approve uxd -p laundry-service

# 5. Synthesize solution architecture
multispec synthesize trd -p laundry-service

# 6. Check status
multispec status -p laundry-service
```

## Switch Interviews

To understand jobs, conduct "Switch" interviews:

**Timeline**: When did you start thinking about switching?

**Push**: What wasn't working with your old solution?

**Pull**: What attracted you to the new solution?

**Anxiety**: What concerns did you have about switching?

**Habit**: What made it hard to change from the old way?

## Reference Materials

For deeper understanding of Jobs to be Done, see:

- [Christensen Institute JTBD](https://www.christenseninstitute.org/theory/jobs-to-be-done/)
- [Strategyn ODI](https://strategyn.com/jobs-to-be-done/)
- *Competing Against Luck* by Clayton Christensen
- *Jobs to be Done* by Anthony Ulwick
- Internal reference: `frameworks-internal/jobs-to-be-done/`
