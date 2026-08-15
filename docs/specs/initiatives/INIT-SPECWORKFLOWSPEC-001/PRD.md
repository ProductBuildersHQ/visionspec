# PRD — Spec Hardening — Normative Contracts, Layered Rubrics, and Content Provenance

**Initiative:** `INIT-SPECWORKFLOWSPEC-001`
**Repository:** `github.com/ProductBuildersHQ/specification-workflow-spec`
**Status:** Draft — pending stakeholder review (pbhq-lite gate: after PRD)

## 1. Problem

The specification-workflow-spec library is the definition-side contract for the
ProductBuildersHQ ecosystem (VisionStudio, visionspec, the website). Three
problems undermine it as an input to AI-assisted development:

1. **Templates are design outlines, not contracts.** The enterprise TRD (and
   its siblings) asks for a proposed design but never states what the
   implemented system is obligated to guarantee: no stable requirement IDs, no
   traceability, no invariants, no MUST-language. An AI agent implementing from
   these documents must invent product decisions. (Verified: the enterprise TRD
   contains zero instances of traceability tables, invariants, or normative
   keywords.)
2. **Templates and rubrics disagree.** Rubrics require content the paired
   template never asks for (PRD rubric demands decision-reversibility analysis
   the PRD template has no section for; TPD rubric demands PRESS-promise
   validation the TPD template never mentions). An author can follow a template
   perfectly and still fail its rubric.
3. **Content had no provenance.** Templates/rubrics were duplicated across
   repos (prism-roadmap ↔ this library) with no declared source of truth, and
   workflows could silently reference specs whose templates resolved to
   nothing — pbhq-lite shipped with rubrics for trd/plan/roadmap but no
   templates anywhere in its resolution chain, and nothing failed.

Problem 3 is solved (Phase 1, complete). Problems 1–2 are this initiative's
remaining scope.

## 2. Goals

- **G1.** Every required spec in every workflow resolves to a template and a
  rubric, enforced by a test with no allowlist. *(done, Phase 1)*
- **G2.** Every template/rubric file has one declared, machine-enforced source
  of truth. *(done, Phase 1)*
- **G3.** The enterprise TRD/PRD/TPD/IRD/UXD family reads as normative
  contracts: stable requirement IDs, bidirectional traceability, invariants,
  measurable non-functionals.
- **G4.** Template–rubric alignment: an author who fully completes a template
  can pass its rubric; every required rubric criterion has a template section
  that satisfies it, or an explicit N/A-with-rationale path.
- **G5.** Rubrics separate judgment layers — AWS Leadership-Principle style
  criteria (advisory, semantic) never gate implementation readiness
  (blocking, largely deterministic) — via first-class schema metadata in
  structured-evaluation.

## 3. Non-Goals

- **NG1.** No changes to the workflow set itself (no new/removed workflows,
  no changes to which specs are required per workflow).
- **NG2.** No execution logic — the library stays declarative; judging remains
  in consumers (`visionstudio spec judge`).
- **NG3.** No migration of PRD/TRD/MRD template ownership to prism-roadmap in
  this initiative (recorded in ROADMAP backlog).
- **NG4.** No re-labeling of rubrics outside the enterprise tree and
  aws-two-way-door (big-tech, google, etc. follow in a later pass).

## 4. Users and Stories

| ID | Story | Acceptance criterion |
|----|-------|----------------------|
| US-001 | As a **spec author**, when I complete every section of a template, its rubric passes, so I never fail on content the template didn't ask for. | AC-001: For each hardened pair, every required rubric criterion maps to ≥1 template section (or an N/A-with-rationale mechanism). |
| US-002 | As an **AI implementation agent**, I can implement from a TRD without inventing product decisions, because guarantees are enumerated with stable IDs, invariants, and testable language. | AC-002: Hardened TRD contains ID'd requirements (MUST/MUST NOT), an invariants section, and a PRD↔TRD↔TPD traceability matrix. |
| US-003 | As an **LLM judge (or its operator)**, I can distinguish advisory Leadership-Principle findings from blocking implementation-readiness gates. | AC-003: Rubric criteria carry `class`, `blocking`, and `evaluation` metadata; a leadership_principle criterion is never blocking. |
| US-004 | As a **library consumer (VisionStudio)**, I get the hardened content by upgrading the module — resolution, inheritance, and the catalog export keep working unchanged. | AC-004: Full test suite (incl. `TestRequiredSpecsResolve`, catalog xref) green; no breaking API changes outside the announced rubric-schema minor bump. |
| US-005 | As the **library maintainer**, I can regenerate any synced content deterministically and CI tells me when copies drift from their source. | AC-005: `prism-sync -check` green; provenance headers on all synced files. *(done, Phase 1)* |

## 5. Functional Requirements

| ID | Requirement | Story | RMI |
|----|-------------|-------|-----|
| FR-001 | The rubric schema MUST support per-criterion `class` (leadership_principle, specification_quality, implementation_readiness, deterministic_integrity), `blocking`, and `evaluation` (deterministic, semantic, human) metadata, additively and back-compatible with v0.13.0 YAML. | US-003 | 006 |
| FR-002 | The rubric schema MUST support evidence-based judge instructions (cite sections/IDs, distinguish missing vs negative evidence, report confidence). | US-003 | 006 |
| FR-003 | The enterprise TRD template MUST include: stable TRD-IDs, requirements/traceability matrix, system invariants, state & lifecycle, interface contracts, operational guarantees, compatibility & migration guarantees. | US-002 | 008 |
| FR-004 | The enterprise PRD template MUST include: Assumptions & Constraints, Decision Register (reversibility), stable story/AC IDs, FR→story traceability; and MUST NOT carry execution timeline detail. | US-001 | 009 |
| FR-005 | The enterprise TPD template MUST include TRD/IRD coverage matrices and the PRESS→PRD→UXD→TRD/IRD→TPD promise-traceability chain, with risk-based applicability for resiliency/chaos content. | US-001 | 010 |
| FR-006 | The enterprise IRD template MUST include IRD-ID requirements, environment model, Pulumi organization constraints, resource-lifecycle safety, deployment guarantees, DR (RTO/RPO), and an explicit "Infrastructure impact: None" reviewed-decision path. | US-002 | 011 |
| FR-007 | The enterprise UXD template MUST define all user-visible states, permission-differentiated views, copy status, measurable accessibility, and UXD→PRD traceability. | US-001 | 012 |
| FR-008 | Every hardened rubric MUST be layered per FR-001 with unconditional criteria replaced by applicability-aware ones (`N/A with rationale` never scored as missing). | US-001, US-003 | 008–012 |
| FR-009 | The aws-two-way-door local set MUST receive the same hardening in feature-level framing, with Leadership-Principle criteria retained as `class: leadership_principle`, non-blocking. | US-001–US-003 | 013 |
| FR-010 | All three repos MUST pass full verification, and release MUST follow dependency order (structured-evaluation → prism-roadmap → specification-workflow-spec) with no replace directives surviving. | US-004 | 014 |

## 6. Success Metrics

| Metric | Baseline (2026-08-12) | Target | Measurement |
|--------|------------------------|--------|-------------|
| Required specs resolving to template+rubric | 28 gaps | 0 gaps, test-enforced | `TestRequiredSpecsResolve` *(achieved 2026-08-13)* |
| Hardened template families (enterprise + aws-two-way-door) | 0 / 10 | 10 / 10 | RMI-008–013 complete |
| Required rubric criteria with a satisfying template section | unmeasured (known mismatches in PRD/TRD/TPD) | 100% per hardened pair | Per-pair alignment review recorded in TPD-style checklist; spot-checked by `visionstudio spec judge` |
| Rubric criteria carrying layer metadata (hardened set) | 0% | 100% | Schema-validated parse of hardened rubrics |
| Blocking criteria that are Leadership-Principle class | n/a | 0 | Deterministic check over parsed rubrics |

## 7. Assumptions and Constraints

- structured-evaluation accepts an additive minor release (v0.14.0); its
  maintainer is the same stakeholder (single-person, multi-repo initiative).
- No commits/tags until explicitly authorized; work staged in three working
  trees; release order is a hard constraint (FR-010).
- enterprise changes propagate by inheritance to 7 workflow families — this is
  intended (it is the mechanism that makes the fix systemic), but reviewers
  must evaluate the diff against inheritors, not just enterprise.

## 8. Open Questions

- OQ-1: Should `spec` (the unified specification) template also be hardened in
  this pass, or after the five core families prove the pattern? *(default: after)*
- OQ-2: Does `narrative-6p` warrant Leadership-Principle labeling in the
  aws-two-way-door pass, or only prd/trd/tpd/ird/uxd? *(default: only the five)*
