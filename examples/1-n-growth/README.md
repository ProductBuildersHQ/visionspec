# 1-N Growth Example

A metrics-driven visionspec configuration for the 1-N growth phase: scaling what works after finding product-market fit.

## Philosophy

You've found product-market fit. Now scale it:

- **Metrics-driven development** - Every feature tied to growth metrics
- **Experiment-first culture** - A/B test before full rollout
- **Funnel optimization** - Know where you are in AARRR
- **Speed with guardrails** - Move fast but protect what works

## What's Required

| Spec | Required | Purpose |
|------|----------|---------|
| PRD | Yes | With growth metrics and experiment design |
| UXD | Yes | Consistent experience at scale |
| FAQ | Yes | Support scaling users |
| MRD | No | For new market expansion |
| TRD | No | As architecture complexity grows |

## Key Features

### Growth Metrics Integration

PRDs include:

- Current metrics and targets
- Hypothesis format
- Funnel stage mapping (AARRR)
- Guardrail metrics

### Experiment Support

- A/B test design in PRD
- Sample size and duration
- Feature flag planning
- Rollout phases

### Analytics Requirements

- Events to track
- Dashboard needs
- Measurement plan

## Building

```bash
go build -o growth-spec ./examples/1-n-growth
```

## Usage

```bash
# Start a growth project
growth-spec init checkout-optimization

# See AARRR metrics framework
growth-spec metrics pirate

# North Star Metric guidance
growth-spec metrics north-star

# Experimentation framework
growth-spec experiment framework

# Pre-launch checklist
growth-spec experiment checklist

# Funnel optimization guide
growth-spec funnel
```

## Evaluation Criteria

The growth PRD rubric evaluates:

1. **Growth Context** - Current metrics and opportunity
2. **Hypothesis** - Testable prediction with expected impact
3. **Funnel Mapping** - Features tied to AARRR stages
4. **Experiment Design** - A/B test setup (if applicable)
5. **Analytics** - Tracking and measurement
6. **Guardrails** - Protection against regressions
7. **Rollout Plan** - Phased deployment

## AARRR Funnel Stages

Tag each user story with its funnel stage:

| Stage | Focus | Example Metrics |
|-------|-------|-----------------|
| Acquisition | How users find you | CAC, traffic, signups |
| Activation | First value experience | Onboarding %, time to value |
| Retention | Users coming back | DAU/MAU, churn, frequency |
| Revenue | Monetization | ARPU, LTV, conversion |
| Referral | Viral growth | NPS, viral coefficient |

## When to Graduate

Move to **post-PMF enterprise** when:

- [ ] Team growing beyond 20-30 people
- [ ] Compliance requirements emerge
- [ ] Enterprise customers with security needs
- [ ] Multiple products/platforms

## Anti-Patterns

**Don't do this in 1-N:**

- Ship without metrics
- Skip A/B testing "because we know it's good"
- Ignore guardrail metrics
- Optimize vanity metrics
- Build without funnel context

**Do this instead:**

- Hypothesis for every feature
- Measure before and after
- Protect retention while growing acquisition
- Focus on North Star Metric
- Ship → Measure → Learn → Repeat
