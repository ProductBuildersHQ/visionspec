# PLAN — Spec Hardening — Normative Contracts, Layered Rubrics, and Content Provenance

**Initiative:** `INIT-SPECWORKFLOWSPEC-001`
**Repository:** `github.com/ProductBuildersHQ/specification-workflow-spec`
**Status:** Draft

The PLAN owns sequencing and rationale. Requirements live in PRD/TRD; execution
state lives in ROADMAP. Phase 1 (provenance + source of truth) is complete;
this plan covers Phases 2–4.

## 1. Approach

Schema first, then content, then adaptation. The layered-rubric metadata
(structured-evaluation) is the foundation every hardened rubric is authored
against — authoring rubrics first and retrofitting metadata would mean touching
all ten files twice. Within content, enterprise precedes aws-two-way-door
because inheritance propagates the enterprise fix to seven workflow families
for free, and the adaptation pass can then diff against a settled canonical
form. TRD leads the template work because it is the ideation doc's centerpiece
(design outline → technical contract) and establishes the ID/traceability
vocabulary (TRD-xxx ↔ PRD-xxx ↔ TPD-xxx) every other template references.

## 2. Affected Components

| Component | Change |
|-----------|--------|
| plexusone/structured-evaluation `rubric/` | Additive schema fields + Validate() + tests (v0.14.0) |
| specification-workflow-spec `go.mod` | Bump structured-evaluation (temp replace until tag) |
| `pkg/workflows/default/enterprise/{templates,rubrics}/` | trd, prd, tpd, ird, uxd hardened pairs |
| `pkg/workflows/default/aws-two-way-door/{templates,rubrics}/` | same, feature-level framing |
| `pkg/workflows/` tests | layered-rubric assertions (class/blocking/weights) |
| grokify/prism-roadmap | none in Phases 2–4 (release sequencing only) |

## 3. Work Breakdown

| Step | Work | Satisfies | Verification | Depends on |
|------|------|-----------|--------------|------------|
| P-01 | structured-evaluation: Class/Blocking/Evaluation + JudgeInstructions, Validate(), back-compat corpus test | TRD-001..003 | SE unit tests | — |
| P-02 | Wire library to extended schema (replace + bump), suite green | TRD-012 | full suite | P-01 |
| P-03 | enterprise TRD template rewrite + layered rubric | TRD-004, 009, 011 | section↔criterion map; weight test | P-02 |
| P-04 | enterprise PRD template + rubric | TRD-005, 009, 011 | same | P-03 (uses TRD ID vocabulary) |
| P-05 | enterprise TPD template + rubric (coverage matrices reference TRD/IRD IDs) | TRD-006, 009, 011 | same | P-03 |
| P-06 | enterprise IRD template + rubric | TRD-007, 009, 011 | same | P-03 |
| P-07 | enterprise UXD template + rubric | TRD-008, 009, 011 | same | P-04 (traces to PRD IDs) |
| P-08 | Layered-rubric test: parse hardened sets, assert class/blocking/N-A invariants (INV-3) | TRD-009 | new test green | P-03..P-07 |
| P-09 | aws-two-way-door adaptation (5 template+rubric pairs, LP criteria relabeled) | TRD-010 | scoped tests | P-08 |
| P-10 | Three-repo sweep; drop replace directives; release-order checklist | TRD-012, PRD FR-010 | build/test/vet/lint × 3; prism-sync -check | P-09 |

## 4. Sequencing Notes

- P-04..P-07 are parallelizable in principle; executed serially to keep each
  template's ID vocabulary consistent with the TRD anchor (P-03).
- Review gates (pbhq-lite): stakeholder review after PRD (before P-01 —
  covers scope); tech-lead review after this PLAN (before implementation).

## 5. Testing Strategy

Deterministic checks carry the load: schema round-trip and back-compat corpus
in structured-evaluation; layered-rubric invariants and weight arithmetic here;
the existing resolution/provenance/xref suite as regression floor. Semantic
validation (`visionstudio spec judge` against hardened rubrics) is a
calibration exercise, deferred to backlog — not a release gate.

## 6. Risks and Open Questions

| ID | Risk | Mitigation |
|----|------|------------|
| R-1 | Replace directives leak into a push | TODO(release) markers; P-10 checklist item; pre-push checklist already forbids them |
| R-2 | Enterprise changes surprise inheriting workflows (7 families) | Full-suite gate at every step; review diff explicitly lists inheritors |
| R-3 | Hardened templates grow so long they stop being used | Keep every added section table-first and skimmable; N/A-with-rationale keeps small features cheap |
| R-4 | Rubric layer semantics drift between repos | Validate() lives in structured-evaluation (TQ-1 default), a single owner |
| R-5 | Three-repo release ordering error | Explicit order in TRD §5; single release checklist in P-10 |
