---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Meaning inference

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/meaning-inference?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/meaning-inference?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/meaning-inference?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/meaning-inference?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P4 is meaning inference: suggesting a field's semantic meaning from column names, types, neighbouring columns and conventions (and, later, safe data sampling), with every suggestion inspectable and correctable — never silently authoritative.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

It exists so that mapping a physical schema onto the semantic layer (P3) does not require hand-declaring every field before the loop becomes useful; it is what lets a new project reach "mostly mapped" quickly while keeping a human (or an audit trail, see P12) in charge of what is actually declared. Initial audit snapshot (before later 9 September landings): `EntityField.NamePatterns` (regexp/exact rules) exist as data, but the audit found no code in `datatug-core` that executes pattern matching against a scanned schema — it is "implemented-but-poorly-exposed (data only)". Probabilistic suggestion from type/neighbours/ sampling and a review-and-correct workflow are both planned, unbuilt. Phase 1 (`core-investigation-loop`) takes the first, deliberately narrow step: the field-mapping resolver MAY derive an inferred mapping from `NamePatterns` and MUST label it. Phase 2 adds review/correct/ provenance; Phase 5 adds probabilistic suggestion. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> Hand-declaring a mapping for every field in every table does not scale past a handful of entities,
> and a new project with an unmapped schema cannot yet benefit from P5's related-data navigation or
> P7's applicable-query resolution. The gap is not the rule format — `EntityField.NamePatterns`
> already exists and is well-typed (`datatug-core/pkg/datatug/entities.go:67-104`) — it is that
> nothing runs those rules against a live scanned schema (current-state audit, capability table,
> `NamePatterns` row: "No code in `datatug-core` executes pattern matching against a scanned
> schema"). Marketing on `datatug.io` claims an "AI-built knowledge base, auto-classify" capability
> (`datatug.io/knowledge-base.astro:18`) that the audit found "unclear/likely aspirational" with no
> corresponding classifier code in any audited repo. Without an inference step that is honest about
> its confidence and always correctable, a "meaning suggestion engine" risks becoming a black box that
> silently mislabels a field — the opposite of DataTug's inspectability thesis.

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| Deterministic name-pattern resolver (`NamePatterns` → column classification) | implemented-but-poorly-exposed (data only) | 1 | strategic |
| Suggestions from type, neighbours, conventions, sampling | planned | 5 | strategic |
| Review, accept, correct, and record provenance of a mapping | planned | 2 | strategic |

#### REQ: deterministic-name-pattern-resolver

The resolver MUST be able to evaluate `EntityField.NamePatterns` regexp/exact rules against a
scanned schema's columns and produce a candidate mapping when exactly one column matches. This is
the deterministic first slice `core-investigation-loop` REQ:field-mapping-model already specifies
("the resolver MAY derive one from `EntityField.NamePatterns`... and MUST label such a mapping as
`inferred`"); this Feature owns extending it to every source and surface, not only the Phase 1
endpoint.

#### REQ: probabilistic-meaning-suggestions

Beyond exact-pattern matches, the resolver SHOULD be able to suggest a likely mapping from column
type, neighbouring columns and naming conventions, and MAY use safe data sampling (never sampling a
field a principal cannot read). Every such suggestion MUST carry a confidence signal and MUST be
distinguishable in every response from a declared or deterministically-inferred mapping. This
capability is an experiment before it is a feature: it MUST be validated for precision on the demo
project before it reaches any UI.

#### REQ: review-accept-correct-provenance

A user MUST be able to review a suggested or inferred mapping, accept it (promoting it to declared),
or correct it (replacing the suggestion with an explicit mapping). Every acceptance or correction
MUST be recorded with who or what made the change and when, feeding `provenance`'s mapping-
provenance record (see `provenance` Feature, P12) rather than silently overwriting history.

## Acceptance Criteria

### AC: name-pattern-resolves-single-match (verifies REQ:deterministic-name-pattern-resolver)

**Given** an `EntityField.NamePatterns` rule for `Customer.ID` matching `CustomerId`/`CustId`/`customer_id`
**When** the resolver runs against the `chinook.Customer` scanned schema
**Then** it proposes `CustomerId → Customer.ID` labelled `inferred`, matching `core-investigation-loop` AC:inferred-mapping-labelled.

### AC: ambiguous-pattern-does-not-guess (verifies REQ:deterministic-name-pattern-resolver)

**Given** a `NamePatterns` rule that matches two columns in the same table
**When** the resolver runs
**Then** it MUST NOT silently pick one; it returns no mapping (or both as candidates for review) rather than an unlabelled guess.

### AC: suggestion-confidence-shown (verifies REQ:probabilistic-meaning-suggestions)

**Given** a column with no declared or pattern-matched mapping but a plausible neighbour-based suggestion
**When** the suggestion is returned
**Then** it carries a confidence value and is visually distinct from a `declared` or `inferred` (deterministic) mapping.

### AC: accept-promotes-to-declared (verifies REQ:review-accept-correct-provenance)

**Given** an inferred mapping `CustomerId → Customer.ID`
**When** a user accepts it in the review UI
**Then** the entity file is updated with a declared mapping and a provenance record naming the user and timestamp.

### AC: correction-overrides-suggestion (verifies REQ:review-accept-correct-provenance)

**Given** a suggested mapping the user judges wrong
**When** the user corrects it to a different column
**Then** the corrected mapping becomes declared, the rejected suggestion is not silently reused later, and the correction is recorded in provenance.

## Not Doing / Out of Scope

- A data catalog or governance suite: inference exists to make navigation and execution faster, not
  to certify meaning organization-wide.
- A text-to-SQL chat box: suggestions are structural (name/type/neighbour based), not a
  natural-language interpretation of intent.
- Automatic relationship discovery by value-overlap sampling: that is a P5 experiment (Phase 5/X),
  kept separate from field-meaning suggestion.
- Any suggestion becoming authoritative without an explicit accept/correct step; DataTug never
  silently promotes a guess to declared fact.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
