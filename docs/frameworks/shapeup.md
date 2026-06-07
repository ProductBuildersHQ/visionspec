# Shape Up

Shape Up is Basecamp's product development methodology created by Ryan Singer. It emphasizes fixed time with variable scope, betting on shaped pitches rather than managing backlogs.

## When to Use

Use Shape Up when:

- You want fixed timelines with flexible scope
- You're tired of endless backlogs and estimation sessions
- You want small, autonomous teams
- You prefer betting on shaped work over planning sprints

## Core Concepts

| Concept | Description |
|---------|-------------|
| **Appetite** | Fixed time budget (2 or 6 weeks), not estimates |
| **Shaping** | Define problem and solution at the right abstraction level |
| **Betting** | Leadership bets on pitches, not backlogs |
| **Cycles** | 6-week building cycles with 2-week cool-down |
| **Hill Charts** | Track uncertainty, not completion |
| **Circuit Breaker** | Stop if it's not working |

## The Shape Up Flow

```
1. SHAPING (senior people)
   Shape raw ideas into pitches
       ↓
2. BETTING TABLE (leadership)
   Bet on pitches for next cycle
       ↓
3. BUILDING (small teams)
   6-week cycle with full autonomy
       ↓
4. COOL-DOWN (everyone)
   2-week break for bugs, exploration, recovery
       ↓
   (repeat)
```

## Key Artifacts

### Pitch

The pitch is the core artifact - a shaped problem and solution ready for betting.

| Section | Purpose |
|---------|---------|
| **Problem** | Raw idea, problem statement, evidence |
| **Appetite** | Time budget (2 or 6 weeks) with rationale |
| **Solution** | Breadboards, fat marker sketches, approach |
| **Rabbit Holes** | Things to avoid that could blow up the project |
| **No-Gos** | Explicitly out of scope |

### Hill Chart

Track progress during building by plotting scopes on a hill:

```
                    ▲ Top (figured out)
                   /|\
    Uphill        / | \        Downhill
  (uncertain)    /  |  \     (executing)
                /   |   \
───────────────/────┼────\───────────────
   0%         25%  50%   75%         100%
```

## Using the Shape Up Profile

### Initialize a Project

```bash
multispec init my-feature --profile shapeup
```

### Create a Pitch

```bash
multispec draft shapeup-pitch -p my-feature
```

### Evaluate the Pitch

```bash
multispec eval shapeup-pitch -p my-feature
```

### Track Scopes During Building

```bash
multispec draft shapeup-scope -p my-feature
```

## Comparison with Other Methodologies

| Aspect | Shape Up | Scrum | AWS Working Backwards |
|--------|----------|-------|----------------------|
| Time | Fixed (6 weeks) | Fixed (2-week sprints) | Variable |
| Scope | Variable | Fixed (sprint commitment) | Variable |
| Planning | Betting table | Sprint planning | PR/FAQ approval |
| Backlog | None | Yes | Roadmap |
| Estimation | Appetite (time worth) | Story points | None |
| Autonomy | High | Medium | Medium |

## Example Workflow

```bash
# 1. Initialize project
multispec init checkout-improvement --profile shapeup

# 2. Shape and write pitch
multispec draft shapeup-pitch -p checkout-improvement
# ... shape the problem and solution ...

# 3. Evaluate pitch readiness
multispec eval shapeup-pitch -p checkout-improvement
# Score: 85% (PASS)

# 4. Pitch goes to betting table (human decision)
# If bet, assign team and cycle

# 5. During building, track scopes
multispec draft shapeup-scope -p checkout-improvement
# Update hill positions as work progresses

# 6. Check status
multispec status -p checkout-improvement
```

## Rubric Categories

### Pitch Evaluation

| Category | Weight | Description |
|----------|--------|-------------|
| Problem Definition | 20% | Is the problem clear and worth solving? |
| Appetite | 20% | Is the time budget appropriate and reasoned? |
| Solution Shaping | 25% | Is the solution at the right abstraction level? |
| Risk Management | 20% | Are rabbit holes and no-gos identified? |
| Document Quality | 15% | Is the pitch clear and actionable? |

## Principles

1. **Fixed Time, Variable Scope** - Appetite sets the time; scope flexes to fit
2. **Shaping** - Define at the right level of abstraction
3. **Betting Not Planning** - No backlogs, bet on shaped pitches
4. **Small Teams** - One designer, one or two programmers
5. **Full Responsibility** - Teams have full autonomy during the cycle
6. **Circuit Breaker** - If it's not working, stop
7. **Cool-Down** - Two-week break between cycles
8. **No Backlogs** - If it's important, it will come back
9. **Appetite Not Estimates** - How much time is this worth?
10. **Hill Charts** - Track uncertainty, not completion

## References

- [Shape Up (free book)](https://basecamp.com/shapeup)
- [Ryan Singer's talks](https://www.feltpresence.com/)
- [Basecamp](https://basecamp.com/)
