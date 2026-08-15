# TRD — Spec Hardening — Normative Contracts, Layered Rubrics, and Content Provenance

**Initiative:** `INIT-SPECWORKFLOWSPEC-001`
**Repository:** `github.com/ProductBuildersHQ/specification-workflow-spec`
**Status:** Draft

This TRD is written in the normative style it prescribes: stable IDs,
MUST-language, traceability, invariants.

## 1. Architecture

Three-repo topology, dependency order left to right:

```
plexusone/structured-evaluation     grokify/prism-roadmap          ProductBuildersHQ/specification-workflow-spec
  rubric.RubricSet (schema)   ←──     rubrics/*.rubric.yaml          pkg/workflows (embedded library)
  + layer metadata (v0.14.0)          templates/*.md                 tools/prism-sync (nested module) ──imports──▶ prism-roadmap
                                      (domain source of truth)       docs/specs/... (this initiative)
```

- The library embeds all content (`//go:embed`); consumers get no filesystem
  or network access requirements. *(existing, unchanged)*
- `tools/prism-sync` is a separate Go module so prism-roadmap's dependency
  graph never enters the library's `go.mod`. *(shipped, Phase 1)*
- structured-evaluation is the only shared schema dependency; both other repos
  parse rubrics into `rubric.RubricSet`.

## 2. Technical Requirements and Traceability

| TRD ID | Requirement | Source | Verification |
|--------|-------------|--------|--------------|
| TRD-001 | `rubric.Criterion` and `rubric.Category` MUST gain optional fields `Class` (enum: leadership_principle, specification_quality, implementation_readiness, deterministic_integrity), `Blocking` (bool), and `Evaluation` (enum: deterministic, semantic, human), serialized as `class`, `blocking`, `evaluation`. | PRD FR-001 | structured-evaluation unit tests: v0.13.0 YAML parses unchanged; new fields round-trip |
| TRD-002 | `rubric.RubricSet` MUST gain an optional `JudgeInstructions []string` field (evidence rules for LLM judges); existing `JudgePromptTemplate` unchanged. | PRD FR-002 | round-trip test |
| TRD-003 | All schema additions MUST be additive: zero changes to existing field names, types, or YAML tags; v0.13.0 documents MUST parse identically under v0.14.0. | PRD FR-001 | parse-corpus test over all rubrics in both consuming repos |
| TRD-004 | The enterprise TRD template MUST contain, as first-class sections: Technical Requirements & Traceability (table: TRD-ID / requirement / source PRD-ID / verification TPD-ID), System Invariants, State & Lifecycle, Interface Contracts, Operational Guarantees (measurable), Compatibility & Migration. Requirement rows MUST use MUST/MUST NOT/SHOULD/MAY. | PRD FR-003 | template section checklist in review; rubric criterion↔section map |
| TRD-005 | The enterprise PRD template MUST add Assumptions & Constraints and Decision Register (reversibility: one-way/two-way-door) sections; every user story and acceptance criterion MUST carry a stable ID; every FR row MUST reference its story; execution dates MUST NOT appear. | PRD FR-004 | same |
| TRD-006 | The enterprise TPD template MUST add coverage matrices for TRD and IRD IDs and a promise-traceability chain (PRESS→PRD→UXD→TRD/IRD→TPD); resiliency/chaos/runbook sections MUST carry risk-based applicability. | PRD FR-005 | same |
| TRD-007 | The enterprise IRD template MUST add: IRD-ID table (requirement/source/verification), Environment Model, Pulumi Organization, Resource Lifecycle, Deployment Guarantees (preview/policy/drift), Security & Secrets, DR (RTO/RPO), and a "No Infrastructure Changes" declaration block. | PRD FR-006 | same |
| TRD-008 | The enterprise UXD template MUST enumerate all user-visible states (loading, empty, partial, error, timeout, recovery) per flow, permission-differentiated behavior, copy finality, measurable accessibility, and a UXD→PRD trace table. | PRD FR-007 | same |
| TRD-009 | Every hardened rubric MUST: (a) label every criterion with `class` and `evaluation`; (b) set `blocking: true` only on non-leadership_principle criteria; (c) express conditional requirements via applicability (N/A-with-rationale) instead of unconditional demands. | PRD FR-008 | deterministic test parsing hardened rubrics and asserting (a)–(c) |
| TRD-010 | aws-two-way-door local templates/rubrics MUST implement TRD-004..TRD-009 in feature-level framing without losing existing Leadership-Principle criteria (relabeled, non-blocking). | PRD FR-009 | same tests scoped to aws-two-way-door |
| TRD-011 | Weight arithmetic MUST hold for every hardened rubric: category weights sum to 100; criterion weights sum to their category weight. | PRD FR-008 | extend the existing weight test pattern (prism-roadmap `TestV2MOMRubricWeights`) to hardened sets |
| TRD-012 | The library MUST load all hardened content through the existing loaders with `TestRequiredSpecsResolve` and the catalog xref test green; no public API changes beyond the rubric version bump. | PRD FR-010 | full suite, all repos |

## 3. System Invariants

- INV-1: A `local` spec source always resolves against the workflow that
  declared it, never against an inheriting child. *(shipped, Phase 1;
  guarded by TestQuickFixInheritsResolvedProvenance)*
- INV-2: A declared template/rubric source that does not provide the file is a
  load-time error, never a silent fallback.
- INV-3: A rubric criterion with `class: leadership_principle` is never
  `blocking: true` (TRD-009b) — advisory judgment cannot gate implementation.
- INV-4: Synced files differ from their prism-roadmap source only in the
  provenance header (`prism-sync -check` enforced).
- INV-5: v0.13.0 rubric YAML parses under v0.14.0 with identical semantics
  (TRD-003).

## 4. Schema Design Notes (TRD-001/002)

- Enums are string types with constants, matching the existing
  `EvaluationType` pattern in the rubric package; unknown values pass through
  parsing (YAML is permissive) but the package gains a `Validate()` check.
- `Blocking` defaults false; absence of `class` means unclassified (legacy) —
  consumers treat unclassified as specification_quality.
- JSON tags use camelCase (`judgeInstructions`) matching the package's
  existing convention (`evaluationType`, `passCriteria`) rather than the
  PBHQ snake_case default — consistency within the owning module wins.

## 5. Compatibility and Migration

- structured-evaluation: minor bump v0.13.0 → v0.14.0 (additive only).
- prism-roadmap and this library: dependency bumps; no content migration
  required for non-hardened rubrics (fields optional).
- Release order (hard requirement): structured-evaluation → prism-roadmap →
  specification-workflow-spec. Replace directives exist only in working trees
  and MUST be dropped before any push (tracked by TODO(release) comments).

## 6. Verification Strategy

- Unit: schema round-trip + back-compat corpus (structured-evaluation);
  layered-rubric assertions (TRD-009, TRD-011) in this library.
- Integration: existing suite — `TestRequiredSpecsResolve`, provenance tests,
  catalog xref, `prism-sync -check`.
- Semantic: `visionstudio spec judge` calibration run against hardened rubrics
  (ROADMAP backlog; not a release gate for this initiative).

## 7. Open Questions

- TQ-1: Should `Validate()` on RubricSet reject blocking leadership_principle
  criteria (hard error) or should that live only in this library's tests?
  *(default: both — cheap invariant, belongs in the schema owner)*
