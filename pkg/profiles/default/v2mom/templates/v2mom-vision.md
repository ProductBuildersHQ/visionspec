# Vision

<!--
V2MOM Vision Template
The Vision is your north star - what you want to achieve.
It should be inspiring, clear, memorable, and achievable.
-->

## Overview

**V2MOM Level**: {{ .Level }} <!-- company | department | team | individual -->
**Fiscal Period**: {{ .FiscalPeriod }} <!-- e.g., FY2025-H2 -->
**Owner**: {{ .Owner }}
**Last Updated**: {{ .UpdatedAt }}

## Vision Statement

<!--
Write a single, compelling statement that describes the future state you're working toward.
Good visions are:
- Inspiring: Motivates the team to action
- Clear: Easy to understand and remember
- Achievable: Ambitious but realistic
- Aligned: Supports organizational strategy
-->

> {{ .VisionStatement }}

## Vision Rationale

### Why This Vision?

<!-- Explain why this vision matters and how it supports organizational goals -->

### Success Picture

<!-- Describe what success looks like when the vision is achieved -->

### Time Horizon

<!-- When do you expect to achieve this vision? -->

## Parent Alignment

<!-- For department/team/individual V2MOMs only -->

### Parent Vision
<!-- Reference the parent level's vision -->

### Supporting Methods
<!-- Which parent methods does this vision support? -->

| Parent Method | How This Vision Supports It |
|--------------|----------------------------|
| {{ .ParentMethod1 }} | {{ .Support1 }} |
| {{ .ParentMethod2 }} | {{ .Support2 }} |

## ProductContext Links

<!-- Link to relevant ProductContext entities -->

- **Strategies**: <!-- List linked strategy IDs -->
- **Capabilities**: <!-- List related capability IDs -->

---

## Metadata

```yaml
v2mom:
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  parent_v2mom: {{ .ParentV2MOMID }}
  productcontext:
    strategies: []
    capabilities: []
```
