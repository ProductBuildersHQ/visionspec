---
name: create
description: Create a new specification file (mrd or uxd)
arguments: [type]
dependencies: []
---

# Create Specification

Create a new specification file that requires manual authoring.

## Usage

```
/create <type>
```

Where `<type>` is one of:
- `mrd` - Market Requirements Document
- `uxd` - User Experience Design

## Process

### For MRD

1. Check if project is initialized (visionspec.yaml exists)
2. Create `docs/specs/{project}/source/mrd.md` with template
3. If source materials exist (IDEATION.md, etc.), pre-populate from them
4. Display guidance for completing the MRD

### For UXD

1. Verify MRD and PRD exist
2. Create `docs/specs/{project}/source/uxd.md` with template
3. Pre-populate API contracts from PRD if available
4. Display guidance for completing the UXD

## Output

```
✓ Created {type} at docs/specs/{project}/source/{type}.md

Next steps:
1. Edit the file to complete all sections
2. Run `/eval {type}` to check quality
3. Run `/approve {type}` when ready
```

## Templates

### MRD Template
```markdown
# {Project} - Market Requirements Document

## Overview
**Project Name:** {project}
**Author:** {author}
**Date:** {date}
**Version:** 1.0

## 1. Problem Statement
### The Core Problem
<!-- What fundamental problem are we solving? -->

### Why Now?
<!-- What makes this the right time? -->

### Cost of Inaction
<!-- What happens if we don't solve this? -->

## 2. Target Market
### Primary Segment
<!-- Profile, pain points, market size -->

### Secondary Segments
<!-- Additional opportunities -->

## 3. Competitive Landscape
| Competitor | Strength | Weakness |
|------------|----------|----------|
| | | |

## 4. Market Requirements
### Must-Have
| ID | Requirement | Rationale |
|----|-------------|-----------|
| MR-1 | | |

### Should-Have
| ID | Requirement | Rationale |
|----|-------------|-----------|
| MR-N | | |

## 5. Business Goals
| Metric | Target | Timeline |
|--------|--------|----------|
| | | |

## 6. Constraints and Assumptions
### Constraints
-

### Assumptions
-

## 7. Timeline
| Milestone | Date | Description |
|-----------|------|-------------|
| | | |

## 8. Risks
| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| | | | |
```

### UXD Template
```markdown
# {Project} - User Experience Design

## Design Principles
1.
2.
3.

## User Personas
### Persona 1: {Name}
- **Role**:
- **Goals**:
- **Pain Points**:

## User Journeys
### Journey 1: {Name}
```
1. User action
2. System response
3. ...
```

## API Contracts
### Resource: {Name}
```yaml
{Resource}:
  type: object
  properties:
    id:
      type: integer
```

## Error Handling
| Code | Type | Description |
|------|------|-------------|
| 400 | | |

## Pagination
<!-- Pagination format -->

## Authentication
<!-- Auth methods -->
```
