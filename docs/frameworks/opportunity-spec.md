# OpportunitySpec Framework

OpportunitySpec is a merged 12-box framework from [prism-roadmap](https://github.com/grokify/prism-roadmap) that combines Jeff Patton's Opportunity Canvas (discovery-focused) with Marty Cagan's SVPG Opportunity Assessment (business case-focused).

## When to Use

Use OpportunitySpec as optional **post-FAQ deepening** when the risk warrants
discovery depth — typically at feature scale. Both AWS door profiles offer it
alongside MRD: the door sets the ceremony; the scale picks the tool.

| Deepening Tool | Typical Scale | Available In |
|----------------|---------------|--------------|
| OpportunitySpec (discovery canvas) | Feature | `aws-one-way-door`, `aws-two-way-door` |
| MRD (market validation) | Product | `aws-one-way-door`, `aws-two-way-door` |

## The OpportunitySpec Flow

Working Backwards starts with the human-authored Press Release — the founding
artifact — and the FAQ that stress-tests it. The OpportunitySpec deepens the
FAQ's discovery and business-case answers, and feeds PRD synthesis when
present:

```
1. WORKING BACKWARDS (human-authored PR, iterated with FAQ)
   press.md  ⇄  faq.md
   (vision)     (challenge assumptions)
       ↓
2. OPPORTUNITY DEEPENING (synthesized from press + faq, optional)
   opportunity-spec.md
   (12-box discovery canvas; flag unanswered boxes as discovery gaps)
       ↓
3. REQUIREMENTS (synthesized, editable)
   prd.md  ←  press + faq + whichever deepening is present
       ↓
4. USER EXPERIENCE (human-authored)
   uxd.md
       ↓
5. TECHNICAL SPECS (synthesized, editable)
   trd.md  →  tpd.md
```

## OpportunitySpec Structure

The 12-box canvas is organized in a 3×4 grid:

### Row 1: Discovery

| Box | Question | Source |
|-----|----------|--------|
| 1. Users & Problem | Who has the problem and what is it? | Patton + Cagan |
| 2. Current Solutions | How do people solve this today? | Patton + Cagan |
| 3. Solution Ideas | What are our solution concepts? | Patton |

### Row 2: Value

| Box | Question | Source |
|-----|----------|--------|
| 4. User Value | What value does this provide to users? | Patton + Cagan |
| 5. Business Value | What value does this provide to the business? | Patton + Cagan |
| 6. Competitive Edge | Why are we best suited to pursue this? | Cagan |

### Row 3: Market

| Box | Question | Source |
|-----|----------|--------|
| 7. Market & Timing | Who's the market and why now? | Cagan |
| 8. Go-to-Market | How will we get this to users? | Patton + Cagan |
| 9. Success Metrics | How will we measure success? | Patton + Cagan |

### Row 4: Validation

| Box | Question | Source |
|-----|----------|--------|
| 10. Critical Requirements | What must be true for this to succeed? | Cagan |
| 11. Risks & Assumptions | What are we betting on and what could go wrong? | Patton + Cagan |
| 12. Recommendation | Given all the above, what's the recommendation? | Cagan + Patton |

## Using the aws-two-way-door Profile

### Initialize a Feature Project

```bash
multispec init mobile-app-feature --profile aws-two-way-door
```

### Draft the OpportunitySpec

```bash
multispec draft opportunity-spec -p mobile-app-feature
```

The template includes all 12 boxes with guidance for each:

```markdown
## Row 1: Discovery

### Box 1: Users & Problem

*"Who has the problem and what problem are we solving?"*

**Primary Users:**
<!-- List primary user types affected by this problem -->

**Problem Statement:**
<!-- Clear, specific statement of the problem -->

**Evidence:**
<!-- What evidence supports this problem exists? -->
```

### Evaluate the OpportunitySpec

```bash
multispec eval opportunity-spec -p mobile-app-feature
```

The rubric evaluates 5 categories with weighted scoring:

| Category | Weight | Criteria |
|----------|--------|----------|
| Discovery Clarity | 20% | Problem definition, user identification, competitive analysis |
| Value Proposition | 25% | User value, business value, differentiation |
| Market & Timing | 20% | Market definition, timing rationale, success metrics |
| Validation Readiness | 20% | Requirements, risks, recommendation quality |
| Document Quality | 15% | Coherence, completeness, actionability |

**Thresholds:**

- PASS: ≥80%
- PARTIAL: ≥60%
- FAIL: <60%

### Proceed to Working Backwards

Once the OpportunitySpec passes evaluation:

```bash
# Approve OpportunitySpec
multispec approve opportunity-spec -p mobile-app-feature

# Synthesize Press Release from OpportunitySpec
multispec synthesize press -p mobile-app-feature

# Continue with FAQ and PRD
multispec synthesize faq -p mobile-app-feature
multispec synthesize prd -p mobile-app-feature
```

## Prioritization Frameworks

OpportunitySpec includes two prioritization frameworks from prism-roadmap to help evaluate and rank features.

### RICE Scoring

RICE provides quantitative prioritization based on four factors:

```
Score = (Reach × Impact × Confidence) / Effort
```

| Component | Description | Values |
|-----------|-------------|--------|
| **Reach** | Users affected per time period | Integer (e.g., 1000 users/quarter) |
| **Impact** | Effect on each user | massive (3x), high (2x), medium (1x), low (0.5x), minimal (0.25x) |
| **Confidence** | Certainty in estimates | high (100%), medium (80%), low (50%) |
| **Effort** | Resources required | Float (person-months) |

**When to use RICE:**

- Comparing multiple features for roadmap prioritization
- Quantitative justification for investment decisions
- Effort and reach estimates are available

**Example in OpportunitySpec:**

Add RICE scoring to Box 12 (Recommendation) or as a summary section:

```markdown
## Prioritization

### RICE Score

| Component | Value | Justification |
|-----------|-------|---------------|
| Reach | 1,000 users/quarter | Based on current active user count |
| Impact | High (2x) | Solves critical workflow blocker |
| Confidence | Medium (80%) | Customer interviews support this |
| Effort | 2 person-months | Engineering estimate |

**Score:** (1000 × 2.0 × 0.8) / 2 = **800**
```

### Kano Model

Kano classifies features by customer satisfaction impact using paired questions:

1. **Functional:** "If the product HAS this feature, how do you feel?"
2. **Dysfunctional:** "If the product DOES NOT HAVE this feature, how do you feel?"

| Category | Description | Priority |
|----------|-------------|----------|
| **Must-Be** | Basic expectation - absence causes dissatisfaction | 5 (highest) |
| **Performance** | Linear satisfaction - more is better | 4 |
| **Attractive** | Delighter - unexpected positive surprise | 3 |
| **Indifferent** | No significant impact on satisfaction | 1 |
| **Reverse** | Unwanted - presence causes dissatisfaction | 0 |

**When to use Kano:**

- Understanding customer expectations vs. delighters
- Deciding between "must-have" and "nice-to-have"
- Product-market fit validation

**Example in OpportunitySpec:**

Add Kano classification to Box 4 (User Value) or as a summary:

```markdown
### Kano Classification

| Question | Response |
|----------|----------|
| If users HAVE this feature | Like it |
| If users DON'T HAVE it | Dislike it |

**Category:** Performance

Users expect improved performance with this feature and will be dissatisfied without it.
```

### Combined Prioritization

Use both frameworks for comprehensive evaluation:

1. **Kano first:** Classify by customer impact type
2. **RICE for ranking:** Within each Kano category, rank by RICE score

| Kano Category | RICE Score | Action |
|---------------|------------|--------|
| Must-Be | Any | Implement first |
| Performance | High (>1000) | High priority |
| Performance | Low (<500) | Medium priority |
| Attractive | High | Differentiation opportunity |
| Attractive | Low | Nice to have |
| Indifferent | Any | Deprioritize |
| Reverse | Any | Do not implement |

## Integration with prism-roadmap

OpportunitySpec types, templates, and rubrics are maintained in [prism-roadmap](https://github.com/grokify/prism-roadmap):

| Asset | prism-roadmap Path | Description |
|-------|-------------------|-------------|
| Go Types | `canvas/opportunity_spec.go` | OpportunitySpec struct with RICE/Kano |
| Prioritization | `prioritization/rice.go`, `prioritization/kano.go` | RICE and Kano types |
| Template | `templates/opportunity-spec.md` | Markdown template with placeholders |
| Rubric | `rubrics/opportunity-spec.rubric.yaml` | LLM-as-a-Judge evaluation criteria |

multispec imports these canonical assets and can customize them via profile configuration.

## Comparison with MRD

| Aspect | MRD | OpportunitySpec |
|--------|-----|-----------------|
| Scope | New product/market | Feature on existing product |
| Market Analysis | Deep (TAM/SAM/SOM, buyer personas) | Focused (relevant segment, sizing) |
| Discovery | Assumes research done | Includes discovery boxes |
| Business Case | Full investment case | Incremental investment |
| Length | 5-10 pages | 2-4 pages |

## Example Workflow

```bash
# 1. Initialize feature project
multispec init checkout-optimization --profile aws-two-way-door

# 2. Human authors the Press Release — the founding artifact
multispec draft press -p checkout-optimization

# 3. Synthesize FAQ to stress-test the PR; iterate the pair through drafts
multispec synthesize faq -p checkout-optimization
multispec eval faq -p checkout-optimization
multispec approve faq -p checkout-optimization

# 4. Optional deepening: populate the OpportunitySpec from press + faq
multispec synthesize opportunity-spec -p checkout-optimization

# 5. Collaborate on OpportunitySpec — fill the discovery gaps it flags
# ... edit opportunity-spec.md ...

# 6. Evaluate and approve
multispec eval opportunity-spec -p checkout-optimization
# Score: 85% (PASS)
multispec approve opportunity-spec -p checkout-optimization

# 7. Generate PRD (grounded in press + faq + opportunity-spec)
multispec synthesize prd -p checkout-optimization

# 8. Human authors UXD
multispec draft uxd -p checkout-optimization

# 9. Generate technical specs
multispec synthesize trd -p checkout-optimization

# 10. Check status
multispec status -p checkout-optimization
```

## References

- [prism-roadmap OpportunitySpec](https://grokify.github.io/prism-roadmap/canvas/opportunity-spec/)
- [Jeff Patton - Opportunity Canvas](https://www.jpattonassociates.com/opportunity-canvas/)
- [Marty Cagan - Assessing Product Opportunities](https://www.svpg.com/assessing-product-opportunities/)
