# User Experience Design (UXD) - Product

## Document Information

| Field | Value |
|-------|-------|
| Product | {project_name} |
| Version | 1.0 |
| Author | {author} |
| Reviewers | |
| Last Updated | {date} |

The UXD is the product-level **experience contract**: normative user-visible
interaction behavior for the experience the Press Release announced. Backend
implementation details belong in the TRD. Every interaction below has a
stated resulting state, every journey traces to the PRD story it serves, and
every piece of copy is marked final or explicitly provisional — an agent
should never have to guess whether placeholder text is real.

## 1. References

| Document | Link |
|----------|------|
| Press Release | |
| FAQ | |
| PRD | |
| MRD | *or "Not produced — see PRD §5.2"* |

## 2. Press Promise Alignment

*Working backwards made concrete: each headline promise in the Press Release
is delivered by a specific journey below. A promise with no journey is an
announced experience nobody designed; a journey with no promise is scope the
customer was never told about — justify it or cut it.*

| Press Release Promise | Delivered By (Journey) | PRD Story |
|------------------------|--------------------------|-----------|
| | Journey 1 | US-001 |

## 3. User Research Summary

### 3.1 Research Methods

- [ ] User interviews: [count] participants
- [ ] Surveys: [count] responses
- [ ] Usability testing: [count] sessions
- [ ] Analytics review: [data sources]

### 3.2 Key Insights

1. [Insight 1]
2. [Insight 2]

## 4. First-Run Experience

*A new product's day-1 experience is part of the announced promise — the
Press Release's "how to get started" paragraph is a commitment this section
delivers. Trace it like any other journey.*

**Traces to:** US-001 / Press Release getting-started paragraph

| Step | User Action | System Response | Resulting State |
|------|-------------|-------------------|--------------------|
| 1 | *e.g., signs up* | | |
| 2 | *first meaningful outcome ("aha" moment)* | | |

**Time to first value:** <!-- target, measurable -->

## 5. User Journeys

*Every journey traces to the PRD story it serves. Every step has a resulting
state — "then what?" should never be unanswered.*

### Journey 1: [Journey Name]

**Traces to:** US-001

**Persona:** [Persona name]

**Goal:** [What the user wants to accomplish]

| Step | User Action | System Response | Resulting State |
|------|-------------|-------------------|--------------------|
| 1 | | | |
| 2 | | | |

**Success Criteria:**

- [ ] [Criteria 1]

---

## 6. Permission-Differentiated Views

*If any part of this experience differs by role or permission level, state
the difference explicitly — do not leave it to be inferred from the primary
flow.*

| View / Element | Role / Permission | Difference from Primary Flow |
|-------------------|----------------------|----------------------------------|
| | | *or "Not applicable — single-permission experience"* |

## 7. Information Architecture

### 7.1 Navigation Structure

```
[Primary Nav Item 1]
├── [Sub-item 1.1]
└── [Sub-item 1.2]
```

### 7.2 Content Hierarchy

<!-- How content is organized and prioritized -->

## 8. Wireframes

### 8.1 Screen 1: [Name]

**Purpose:**

```
┌─────────────────────────────────────────┐
│  Header                                  │
├─────────────────────────────────────────┤
│  [Wireframe content here]               │
│  ┌─────────┐  ┌─────────┐               │
│  │ Action  │  │ Action  │               │
│  └─────────┘  └─────────┘               │
└─────────────────────────────────────────┘
```

**Interactions:**

- Button A: Does X
- Button B: Does Y

## 9. Interaction Design

*Every interaction states its resulting state or outcome — a button with no
defined result is incomplete, not implied.*

#### Interaction 1: [Name]

**Traces to:** US-001 / FR-001

**Trigger:** [What initiates this interaction]

**Flow:**

1. [Step 1]
2. [Step 2]

**Resulting State / Outcome:** [What state the UI is in after this completes]

**Feedback:** [How the system communicates state during/after]

**Error Handling:** [How errors are communicated]

---

## 10. User-Visible States

*Every flow's full state set — not only the happy path. A state left out here
is a state an implementer will invent.*

| Flow | Loading | Empty | Error | Timeout | Recovery |
|------|---------|-------|-------|---------|----------|
| | | | | | |

### State Detail Template

#### [State Name]

**Trigger:** [What causes this state]

**Message / Copy:** [Exact text]

**Copy status:** {{ Final | Provisional — <what's pending> }}

**Recovery / Next Action:** [How the user proceeds, if applicable]

---

## 11. Content

### 11.1 Copy

| Element | Copy | Copy Status |
|---------|------|-------------------|
| Page title | | {{ Final \| Provisional }} |
| CTA button | | {{ Final \| Provisional }} |
| Empty state | | {{ Final \| Provisional }} |
| Error message | | {{ Final \| Provisional }} |

### 11.2 Voice and Tone

<!-- Voice and tone guidance for this product. For a new product line this
     defines the voice, it doesn't inherit one — be specific enough to guide
     real copy decisions. -->

## 12. Accessibility Requirements (WCAG 2.1 AA)

### Perceivable

- [ ] All non-text content has text alternatives
- [ ] Color contrast ratio: 4.5:1 for normal text, 3:1 for large text
- [ ] Text can be resized to 200% without loss of content
- [ ] Color is not the only means of conveying information

### Operable

- [ ] All functionality available via keyboard
- [ ] No keyboard traps
- [ ] Focus order is logical and focus is visible at all times
- [ ] Touch targets minimum 44x44 CSS pixels
- [ ] Timing adjustable or can be turned off

### Understandable

- [ ] Language of page/parts identified
- [ ] Focus and input don't change context unexpectedly
- [ ] Errors identified, described, and paired with a suggestion

### Robust

- [ ] ARIA used correctly
- [ ] Status messages announced to assistive technology

### Testing Plan

| Test Type | Tool/Method | Frequency |
|-----------|-------------|-----------|
| Automated | axe-core, Lighthouse | Every build |
| Manual | Keyboard-only navigation | Before release |
| Screen reader | NVDA, VoiceOver | Before release |

## 13. Responsive Design

| Breakpoint | Width | Layout Changes |
|------------|-------|-------------------|
| Mobile | 320-767px | |
| Tablet | 768-1023px | |
| Desktop | 1024px+ | |

## 14. Internationalization

- [ ] RTL support required: Yes / No
- [ ] Languages supported: [List]
- [ ] Text expansion buffer: [e.g., 20-30%]

## 15. Prototype Links

| Platform | Link |
|----------|------|
| Figma | |

## 16. Open Questions

| # | Question | Owner | Resolution |
|---|----------|-------|------------|
| 1 | | | |

---

**Approval:**

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Design | | | |
| Product | | | |
| Accessibility | | | |
| Engineering | | | |
