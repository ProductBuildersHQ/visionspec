# Context Sources MRD

## Market Problem

### The Spec-Reality Gap

Organizations face a fundamental disconnect between documentation and implementation:

1. **Brownfield Blindness**: Most spec tools assume greenfield development, ignoring existing codebases
2. **Context Silos**: Relevant information scattered across Jira, Confluence, Google Docs, code repos
3. **Manual Alignment**: Keeping specs updated with shipped code requires constant manual effort
4. **AI Hallucination**: LLM-based tools generate specs that don't reflect actual system state

**Impact**: Engineering teams spend 20-30% of time reconciling specs with reality, or abandon specs entirely.

### Market Segments

| Segment | Pain Point | Opportunity |
|---------|------------|-------------|
| Enterprise | Legacy systems with outdated docs | Auto-generate TRDs from existing code |
| Scale-ups | Fast iteration outpaces documentation | Real-time spec-code alignment |
| Platform teams | Multi-repo complexity | Cross-repo context aggregation |
| Agencies | Client handoffs with poor docs | Generate specs from inherited codebases |

## Competitive Landscape

| Competitor | Approach | Limitation |
|------------|----------|------------|
| Traditional spec tools | Manual authoring | No code awareness |
| AI doc generators | Generate from prompts | Hallucinate without context |
| Code documentation | Extract from code only | No business context |
| ADR tools | Decisions only | No traceability to code |

**Multispec Differentiation**: Unified context aggregation + structured synthesis + reality grounding.

## Target Audience

### Primary: Engineering Leaders

- VP Engineering, Principal Engineers, Staff Engineers
- Pain: Technical specs drift from implementation
- Need: Auto-generated TRDs that reflect actual architecture

### Secondary: Product Managers

- Senior PMs, Technical PMs
- Pain: Can't verify if specs match shipped product
- Need: Automated alignment reports

### Tertiary: Platform/DevOps

- Platform Engineers, DevOps Leads
- Pain: Documentation gaps in CI/CD pipelines
- Need: Automated spec generation in pipelines

## Business Goals

### G1: Increase Adoption in Brownfield Contexts

Current multispec adoption limited to greenfield projects. Context sources enables brownfield use cases, expanding TAM 5x.

### G2: Differentiate from AI Doc Tools

Pure AI tools hallucinate. Context-grounded synthesis produces accurate, verifiable specs.

### G3: Create Platform Lock-in via Integrations

MCP server integrations (Jira, Confluence, etc.) create switching costs and ecosystem value.

## Success Criteria

| Metric | Current | Target |
|--------|---------|--------|
| Brownfield project adoption | ~10% | 50%+ |
| Spec accuracy (user reported) | N/A | 90%+ match reality |
| Integration connections | 0 | 5+ MCP servers |
| Context gather usage | 0 | 40% of synthesize calls |

## Market Positioning

### Tagline Options

- "Grounding" - Specs grounded in reality
- "Reality-Aware Specs" - Documentation that knows your code
- "Context-First Synthesis" - From codebase to spec, automatically

### Messaging

**For Engineering Leaders:**
> "Generate technical specs that actually reflect your architecture. Multispec analyzes your codebase, requirement graphs, and project tools to synthesize TRDs grounded in reality—not hallucinations."

**For Product Managers:**
> "Finally know if your specs match what shipped. Multispec's Grounding feature compares your PRD against the actual codebase and flags drift automatically."

## Risks

| Risk | Mitigation |
|------|------------|
| MCP server ecosystem immature | Start with well-supported servers (Jira, Confluence) |
| Context gathering too slow | Incremental updates, caching, snapshots |
| Privacy concerns with code analysis | Local-only analysis, no data leaves machine |
| Complexity overwhelms users | Progressive disclosure, sensible defaults |

## Timeline

| Phase | Date | Deliverable |
|-------|------|-------------|
| Alpha | Week 2 | Git + Graphize context |
| Beta | Week 4 | MCP client, synthesis integration |
| RC | Week 5 | Snapshots, caching, docs |
| GA | Week 6 | v0.4.0 release |
