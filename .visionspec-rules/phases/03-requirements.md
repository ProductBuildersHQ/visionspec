# Phase 3: Experience

The Experience phase defines how users interact with the product through the User Experience Design (UXD) document.

## Objective

Create an approved UXD that:

- Maps user journeys from PRD user stories
- Defines interaction flows
- Specifies UI/UX requirements

## Entry Criteria

- PRD is approved
- User stories are clear
- Target audience is understood

## Workflow

### Step 1: Initialize UXD

```bash
visionspec create uxd -p <project>
# Or via MCP: start_draft(project, "uxd")
```

### Step 2: Map User Journeys

For each user story in the PRD, define:

1. **Entry point**: How does the user arrive?
2. **Goal**: What is the user trying to accomplish?
3. **Steps**: What actions does the user take?
4. **Exit**: What indicates success?
5. **Errors**: What could go wrong?

### Step 3: Discovery Questions

**Primary Flows**

1. What is the most common user journey?
2. What is the "happy path" for each user story?
3. What triggers the user to start this journey?

**Edge Cases**

4. What happens when things go wrong?
5. How does the user recover from errors?
6. What are the boundary conditions?

**Accessibility**

7. How will users with disabilities interact?
8. What are the keyboard navigation patterns?
9. What screen reader considerations exist?

**Responsive Design**

10. How does the experience differ on mobile?
11. What features are desktop-only?
12. How do touch interactions work?

### Step 4: Fill UXD Sections

| Area | UXD Section |
|------|-------------|
| User journeys | User Journeys |
| Wireframes/mockups | Visual Design |
| Error handling | Error States |
| Accessibility | Accessibility |
| Mobile | Responsive Design |

### Step 5: Evaluate

```bash
visionspec eval uxd -p <project>
```

Check for:

- [ ] All PRD user stories have journeys
- [ ] Error states are defined
- [ ] Accessibility is addressed
- [ ] Score >= 7.0

### Step 6: Approve

```bash
visionspec approve uxd -p <project>
```

## Exit Criteria

- UXD exists at `source/uxd.md`
- Evaluation score >= 7.0
- No critical or high findings
- UXD is approved

## Next Phase

→ [Phase 4: Technical](04-technical.md)

## Anti-Patterns

- **Pixel-perfect too early**: Focus on flows and interactions, not visual polish.
- **Ignoring errors**: Happy path only. Always define error states.
- **Desktop-first blindness**: Consider mobile from the start.
- **Accessibility as afterthought**: Build it in, don't bolt it on.
