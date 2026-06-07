---
name: spec-synthesis
description: Techniques for synthesizing downstream specifications from upstream documents
triggers: [synthesize, generate spec, derive, transform]
---

# Specification Synthesis

Techniques for generating downstream specifications from upstream documents while maintaining consistency and traceability.

## Synthesis Principles

1. **Upstream Authority** - Downstream specs must not contradict upstream
2. **Traceability** - Every downstream element traces to upstream source
3. **Completeness** - All upstream requirements flow downstream
4. **Consistency** - Terminology and IDs are consistent across specs

## Dependency Graph

```
MRD
 ├── Press Release
 │    └── FAQ
 │         └── PRD
 │              ├── UXD (may be authored)
 │              └── TRD
 │                   ├── TPD
 │                   └── IRD
 └── spec.md (reconciled)
```

## Synthesis Techniques

### Information Extraction

Extract key elements from source documents:

```
MRD → Extract:
  - Problem statement → Press headline
  - Target market → Press "Who it's for"
  - Business goals → Press success metrics
  - Constraints → FAQ scope questions
```

### Transformation Patterns

| From | To | Transformation |
|------|-----|----------------|
| MRD Market Requirement | PRD Functional Requirement | Add implementation details |
| MRD Business Goal | PRD Success Metric | Make measurable |
| Press Feature | PRD User Story | Add persona and criteria |
| PRD Requirement | TRD Component | Map to technical design |
| TRD Component | TPD Test Case | Define verification |

### ID Tracing

Maintain ID relationships across specs:

```
MR-1 (Market Requirement)
  → FR-1 (Functional Requirement)
    → US-001 (User Story)
      → TC-001 (Test Case)
```

### Cross-Reference Format

When referencing upstream:
```markdown
<!-- Traces to: MR-1, MR-2 -->
**FR-1**: The system shall provide pet CRUD operations.
```

## Synthesis Process

### Step 1: Load Sources
```python
sources = [
    "source/mrd.md",
    "gtm/press.md",  # if exists
    "gtm/faq.md",    # if exists
]
```

### Step 2: Extract Elements
```python
elements = {
    "requirements": extract_requirements(mrd),
    "features": extract_features(press),
    "constraints": extract_constraints(faq),
}
```

### Step 3: Transform
```python
for req in elements["requirements"]:
    user_story = transform_to_story(req)
    add_acceptance_criteria(user_story)
    link_to_source(user_story, req.id)
```

### Step 4: Validate
```python
validate_completeness(output, sources)
validate_consistency(output, sources)
validate_traceability(output)
```

### Step 5: Write Output
```python
write_spec(output, target_path)
write_metadata(output, eval_path)
```

## Content Mapping

### MRD → Press Release

| MRD Section | Press Section |
|-------------|---------------|
| Problem Statement | "The Problem" |
| Target Market Primary | "Who It's For" |
| Market Requirements (MR-1 to MR-3) | Key Features |
| Business Goals | Success Metrics |
| Constraints | Availability Notes |

### Press → FAQ

| Press Section | FAQ Category |
|---------------|--------------|
| Key Features | Customer Questions |
| Technical Details | Technical Questions |
| Availability | Scope Questions |
| Claims | Challenging Questions |

### MRD + Press + FAQ → PRD

| Source | PRD Section |
|--------|-------------|
| MR-* requirements | FR-* requirements |
| Press features | User Stories |
| FAQ scope | Release Criteria |
| FAQ challenges | Non-Functional Requirements |

### PRD → TRD

| PRD Section | TRD Section |
|-------------|-------------|
| FR-* requirements | Component Design |
| User Stories | API Endpoints |
| NFR-* requirements | Architecture Decisions |
| Data entities | Data Models |

## Quality Checks

### Completeness Check
```
For each MR-* in MRD:
  Assert exists FR-* in PRD that traces to MR-*
  Assert exists US-* in PRD that traces to MR-*
```

### Consistency Check
```
For each term in Glossary:
  Assert consistent usage across all specs
  Assert no contradictory definitions
```

### Traceability Check
```
For each FR-* in PRD:
  Assert traces_to field references valid MR-*
  Assert MR-* exists in MRD
```

## Conflict Resolution

When sources conflict:

1. **Identify** - Document the conflict
2. **Analyze** - Determine which source is authoritative
3. **Resolve** - Choose resolution, document rationale
4. **Update** - Update upstream if needed, or note in Decision Log

## Output Metadata

Each synthesized spec includes metadata:

```yaml
---
synthesized_from:
  - source/mrd.md
  - gtm/press.md
synthesized_at: 2026-06-01T12:00:00Z
synthesized_by: synthesizer-agent
traceability:
  FR-1: [MR-1, MR-2]
  FR-2: [MR-3]
  US-001: [FR-1]
---
```
