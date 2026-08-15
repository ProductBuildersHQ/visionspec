# Roadmap (ROADMAP)

## Overview

**Project Name:** {project_name}
**Author:** {author}
**Date:** {date}
**Version:** 1.0
**Status:** Active

The ROADMAP owns *whether and when* work happens: a machine-readable ledger of
Roadmap Items (RMIs) with stable IDs, status, and completion evidence. It is
non-normative — it records execution state and must not own requirements or
design decisions. Keeping it separate from the PLAN means per-session status
churn never touches the plan's narrative.

## 1. Phases

Themed phases group related RMIs. Phase status is derived from its members'
statuses, never set directly.

| Phase | Theme | RMIs | Status |
|-------|-------|------|--------|
| Phase 1 | | RMI-{repo}-001 … | |
| Phase 2 | | | |

## 2. Roadmap Items (RMIs)

Each RMI has a stable ID (`RMI-<REPOSLUG>-<NNN>`), the requirements it satisfies,
its dependencies, and — when done — a link to the evidence that proves it. A
checked box is not evidence; the completion reference is.

### RMI-{repo}-001: {deliverable}

- **Status:** planned <!-- planned | in-progress | blocked | done -->
- **Satisfies:** PRD-... , TRD-... <!-- requirement IDs this RMI delivers -->
- **Depends on:** —
- **Verification:** TPD-... <!-- test/eval that proves completion -->
- **Completion evidence:** <!-- commit, PR, eval output, or approved artifact -->
- **Blocking decisions:** <!-- open decisions that gate this RMI, if any -->
- **Owner:** <!-- optional -->

### RMI-{repo}-002: {deliverable}

- **Status:** planned
- **Satisfies:**
- **Depends on:** RMI-{repo}-001
- **Verification:**
- **Completion evidence:**
- **Blocking decisions:**
- **Owner:**

## 3. Traceability

The roadmap is a traceability index, not a second requirements document. Every
normative requirement should trace forward to at least one RMI.

| Requirement | RMI(s) | Status |
|-------------|--------|--------|
| PRD-... | RMI-{repo}-001 | |
| TRD-... | | |
