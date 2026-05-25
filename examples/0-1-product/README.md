# 0-1 Product Example

A minimalist visionspec configuration for the 0-1 phase: going from idea to first working product.

## Philosophy

> "If you're not embarrassed by the first version of your product, you've launched too late." - Reid Hoffman

The 0-1 phase is about **learning**, not building. This configuration:

- **Replaces traditional specs with a single hypothesis document**
- **Focuses on testability over completeness**
- **Optimizes for learning velocity**
- **No bureaucracy, maximum iteration speed**

## What's Required

| Spec | Required | Purpose |
|------|----------|---------|
| hypothesis | Yes | Your testable bet |
| Everything else | No | Add later if needed |

## The Hypothesis Document

Instead of MRD/PRD/UXD, you write a single hypothesis:

```
We believe that [specific target users]
have a problem with [concrete problem].

We believe that [proposed solution]
will solve this problem.

We will know we're right when [measurable success metric].
```

## Building

```bash
go build -o zero-to-one ./examples/0-1-product
```

## Usage

```bash
# Start a new experiment
zero-to-one init meal-planner

# See hypothesis template
zero-to-one hypothesis

# Check if hypothesis is testable
zero-to-one validate

# Guidance on pivoting
zero-to-one pivot

# Evaluate your hypothesis
zero-to-one eval hypothesis
```

## Evaluation Criteria

The hypothesis rubric checks only three things:

1. **Specificity** - Are users and problem specific enough to test?
2. **Testability** - Can you test this in under 2 weeks?
3. **Measurability** - Do you have a number to track?

## When to Graduate

Move to **pre-PMF startup** configuration when:

- [ ] Hypothesis validated with real users
- [ ] Ready to build a more complete product
- [ ] Need to communicate with others (investors, team)

Don't graduate too early. Stay in 0-1 until you have signal.

## Anti-Patterns

**Don't do this in 0-1:**

- Write a 10-page PRD
- Design detailed UX flows
- Plan technical architecture
- Create GTM materials
- Build for scale

**Do this instead:**

- One hypothesis
- One metric
- One week to build
- One week to test
- Decide: proceed, pivot, or kill
