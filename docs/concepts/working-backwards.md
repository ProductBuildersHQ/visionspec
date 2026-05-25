# Working Backwards

VisionSpec implements Amazon's Working Backwards methodology as the foundation for specification development. This document explains the approach, why it matters for AI-assisted execution, and how to apply it effectively.

## What is Working Backwards?

Working Backwards is a product development methodology pioneered at Amazon. Instead of starting with requirements and hoping they lead to a good customer experience, you start with the end state and work backwards:

1. **Write the Press Release first** - Describe the product as if it's already shipped
2. **Challenge it with FAQs** - Ask hard questions about scope, feasibility, and gaps
3. **Derive requirements** - Only then write detailed requirements grounded in the vision

This inversion ensures that every requirement traces back to a specific customer outcome.

## The VisionSpec Flow

```
MRD (human-authored)
    │
    │  "What market problem are we solving?"
    ↓
Press Release (synthesized from MRD)
    │
    │  "How will we announce this to customers?"
    ↓
FAQ (synthesized from MRD + Press)
    │
    │  "What questions will stakeholders ask?"
    ↓
PRD (synthesized from MRD + Press + FAQ)
    │
    │  "What detailed requirements follow from this vision?"
    ↓
Narrative 1P/6P (synthesized, for stakeholder review)
    │
    │  "Does leadership align on this vision?"
    ↓
UXD (human-authored)
    │
    │  "How will users interact with this?"
    ↓
TRD (synthesized from MRD + PRD + UXD + context)
    │
    │  "How will we build this technically?"
    ↓
IRD (synthesized from TRD + context)
    │
    │  "How will we deploy and operate this?"
    ↓
spec.md (reconciled from all approved specs)
    │
    │  "What should the AI agent execute?"
    ↓
AI-Assisted Execution
```

## Why This Matters for AI Execution

When an AI coding agent executes from specifications, it benefits from:

### 1. Vision Anchoring

The Press Release creates an unambiguous target state. When the agent encounters conflicting requirements or ambiguous instructions, it can ask: "Which interpretation better serves the announced customer experience?"

### 2. Early Gap Detection

The FAQ forces explicit consideration of edge cases, scope boundaries, and potential objections. Questions like "What happens when..." are answered *before* implementation, not discovered mid-coding.

### 3. Testable Requirements

PRD derived from Press + FAQ tends to produce concrete, testable requirements rather than abstract feature lists. "Users can export reports in PDF format" is more actionable than "the system should support exports."

### 4. Traceability Chain

Every technical decision can trace back through the chain:

```
IRD decision → TRD requirement → PRD feature → FAQ clarification → Press vision → MRD problem
```

This enables principled conflict resolution during reconciliation.

## Synthesized Documents in Git

All synthesized documents (Press, FAQ, PRD, TRD, IRD, Narratives) are committed to git and can be:

- **Reviewed** by stakeholders before approval
- **Edited** by humans to add nuance or correct errors
- **Refined** collaboratively with AI assistants like Claude Code
- **Versioned** alongside the codebase

The initial synthesis provides a strong starting point, but human judgment remains in the loop at every stage.

## Where Narratives Fit

The 1-page and 6-page narratives serve as **stakeholder alignment documents**:

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **1-Page Narrative** | Executive summary for quick alignment | Leadership reviews, status updates |
| **6-Page Narrative** | Deep-dive for thorough review | Team meetings, architecture reviews |

Synthesize narratives after PRD but before starting technical work:

```bash
# Generate narratives for review
visionspec synthesize narrative-1p
visionspec synthesize narrative-6p

# Share with stakeholders, gather feedback
# Edit narratives as needed based on feedback

# Then proceed to technical synthesis
visionspec synthesize trd
```

## Human vs. Synthesized Documents

| Document | Default Source | Can Be Synthesized? | Typical Workflow |
|----------|----------------|---------------------|------------------|
| MRD | Human-authored | No | Product/business owner writes |
| Press | Synthesized | Yes | Synthesize, review, refine |
| FAQ | Synthesized | Yes | Synthesize, add questions |
| PRD | Either | Yes | Synthesize or human-author |
| UXD | Human-authored | No | Designer creates |
| Narrative | Synthesized | Yes | Synthesize for review |
| TRD | Synthesized | Yes | Synthesize, architect reviews |
| IRD | Synthesized | Yes | Synthesize, SRE reviews |

## Recommended Workflow

### For New Projects

```bash
# 1. Initialize project
visionspec init my-feature

# 2. Write MRD (human-authored)
# Define market problem, audience, business goals

# 3. Synthesize Working Backwards chain
visionspec synthesize press    # Vision document
visionspec synthesize faq      # Challenge assumptions
visionspec synthesize prd      # Detailed requirements

# 4. Generate narratives for stakeholder review
visionspec synthesize narrative-1p
visionspec synthesize narrative-6p

# 5. Review and refine all documents
# Edit in git, collaborate with AI, gather feedback

# 6. Write UXD (human-authored)
# Define user journeys and interactions

# 7. Synthesize technical specs
visionspec synthesize trd --eval
visionspec synthesize ird --eval

# 8. Approve all specs
visionspec approve mrd
visionspec approve prd
# ... approve all

# 9. Reconcile into execution spec
visionspec reconcile

# 10. Export to target system
visionspec export speckit
```

### For Existing PRDs

If you already have a human-authored PRD, you can still use Working Backwards for validation:

```bash
# Synthesize Press from existing MRD + PRD
visionspec synthesize press

# Generate FAQ to challenge the PRD
visionspec synthesize faq

# Compare FAQ questions with PRD coverage
# Update PRD to address gaps
```

## Configuration by Profile

Different organizational stages need different levels of ceremony:

| Profile | Working Backwards Documents |
|---------|---------------------------|
| `0-1` | Just hypothesis (no full flow) |
| `startup` | PRD only (minimal ceremony) |
| `growth` | PRD, FAQ (validate scope) |
| `enterprise` | Full flow (MRD → Press → FAQ → PRD → Narratives → TRD → IRD) |

```bash
# Initialize with appropriate profile
visionspec init my-feature --profile growth
```

## References

The Working Backwards methodology was developed at Amazon and is documented in:

- **Bryar, Colin and Bill Carr.** *[Working Backwards: Insights, Stories, and Secrets from Inside Amazon](https://www.amazon.com/Working-Backwards-Insights-Stories-Secrets/dp/1250267595/)*. St. Martin's Press, 2021.

  The definitive guide by two former Amazon executives (Bryar was Jeff Bezos's "shadow" for two years; Carr launched Amazon Music and Prime Video). Covers the PR/FAQ process, single-threaded leadership, and other Amazon mechanisms.

- **AWS Blog.** *[Working Backwards to Drive Customer Experience and SMB Innovation Forward](https://aws.amazon.com/blogs/smb/working-backwards-to-drive-customer-experience-and-smb-innovation-forward/)*

  Practical application of Working Backwards for small and medium businesses.

- **Amazon Jobs.** *[Working Backwards](https://www.amazon.jobs/en/landing_pages/working-backwards)*

  Amazon's official overview of the methodology.

## Related Documentation

- [Synthesize Command](../cli/synthesize.md)
- [Reconcile Command](../cli/reconcile.md)
