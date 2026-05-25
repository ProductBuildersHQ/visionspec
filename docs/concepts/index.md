# Concepts

Core concepts and methodologies behind VisionSpec.

## Working Backwards

VisionSpec implements Amazon's Working Backwards methodology as the foundation for specification development. Instead of starting with requirements, you start with the customer experience and work backwards to derive requirements.

- [Working Backwards Guide](working-backwards.md) - Full methodology and workflow

## Key Ideas

### Vision-First Development

The Press Release is written first, describing the product as if it's already shipped. This creates an unambiguous target state that grounds all subsequent requirements.

### Synthesized but Editable

All LLM-generated documents (Press, FAQ, PRD, TRD, IRD, Narratives) are:

- Committed to git alongside source code
- Reviewable by stakeholders before approval
- Editable by humans to add nuance
- Refinable collaboratively with AI assistants

### Traceability Chain

Every decision traces back through the chain:

```
IRD → TRD → PRD → FAQ → Press → MRD
```

This enables principled conflict resolution during reconciliation.

### Profile-Based Configuration

Different organizational stages need different levels of ceremony:

| Profile | Documents | Use Case |
|---------|-----------|----------|
| `0-1` | Hypothesis only | Idea validation |
| `startup` | PRD | Minimum viable spec |
| `growth` | PRD, FAQ | Validate scope |
| `enterprise` | Full Working Backwards | Complete documentation |
