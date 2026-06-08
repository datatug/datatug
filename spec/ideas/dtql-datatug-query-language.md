---
format: https://specscore.md/idea-specification
status: Approved
---

# Idea: DTQL — DataTug Query Language

**Status:** Approved
**Date:** 2026-06-05
**Owner:** alexander.trakhimenok
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we define DTQL — a 1:1, human-readable YAML serialization of the dalgo query AST — as a first-class part of the dalgo library, so DataTug and any dalgo consumer can save, hand-edit, diff, and reload queries and deterministically render them to each backend's native language, without reinventing GraphQL?

## Context

Promotes the datatug sidekick seed at spec/ideas/seeds/dtql-datatug-query-lanaguage-needs-definition-and.md. The just-merged serve-brokered query builder (Capability in datatug plus Implementations in datatug-cli/datatug-apps/datatug-ai-skills) already depends on DTQL-YAML as the textual form of the dalgo AST: REQ:dtql-yaml-form mandates a lossless AST↔YAML round-trip, and REQ:canonical-ast renders that AST to each connection's native query language. But DTQL is referenced, not defined — no schema, no validator, no owning artifact. During the dual-mode work we settled that GraphQL is a native language (a target/source), not the meta layer; DTQL stays the relational+document AST serialization. Crucially, `github.com/dal-go/dalgo` is **already SpecScore-managed** and its `dal` package **already models the query**: `dal.StructuredQuery` (`From`/`Where`/`GroupBy`/`OrderBy`/`Columns`/`StartFrom`) is the structured AST DTQL serializes, and `dal.TextQuery` (`Text` + `Args`) already represents the native-text form — a near-exact match to the serve-brokered native/DTQL mode split. So this Idea serializes an existing model rather than inventing one.

## Recommended Direction

Define DTQL in the dalgo repo as the canonical, human-readable, 1:1 lossless YAML serialization of dalgo's query AST — structural, hierarchical, plain off-the-shelf YAML with no bespoke grammar or parser. dalgo **already owns the query AST** — its `dal.StructuredQuery` interface (`From`/`Where`/`GroupBy`/`OrderBy`/`Columns`/`StartFrom`); this Idea adds the DTQL YAML schema, (de)serialization, and a round-trip validator proving `StructuredQuery`→YAML→`StructuredQuery` identity. dalgo's existing `dal.TextQuery` (`Text` + `Args`) already represents the native-text form, so DTQL serializes the *structured* side of a split dalgo already makes. DataTug consumes DTQL two ways: as the serve-brokered builder's internal exchange form (the agent and web UI inspect and edit it) and as a user-facing saved query format (.dtql.yaml stored in projects, version-controlled, diffable, shared). Because the AST spans relational and nested/document (Firestore-style) shapes, DTQL is not SQL-shaped; it renders to each connection driver's native query language on execution. GraphQL is recorded as settled: a native language DTQL can target, not a competing meta representation — a GraphQL-as-meta approach would still have to build the dalgo AST underneath. Advantages over the alternatives: structure, hierarchy, and human-readability for free via YAML; deterministic point-and-click editability in the web UI because DTQL is the AST; one model rendered to many dialects; saved, versioned, diffable queries as plain text; and zero new parser to design or maintain.

## Alternatives Considered

**Define DTQL in `datatug-go-models` (or `datatug`) instead of dalgo.** Keeps the schema next to DataTug's usage and avoids touching the dalgo library, but confines a generally useful query serialization to DataTug and splits the AST (the model) from its serialization across repos. Rejected in favour of dalgo so the query model and its DTQL form live together and any dalgo consumer benefits. Trade-off, surfaced as a Must-be-true assumption: dalgo must own a query AST, which it does not today.

**Invent a bespoke DTQL textual grammar (PRQL/Malloy-style) instead of YAML.** A purpose-built syntax reads more fluently and is terser, but costs a parser and grammar to design, test, version, and keep in lock-step with the AST, and is harder to edit point-and-click. Rejected for the MVP — YAML gives structure plus human-readability for free and round-trips with the AST via a standard serializer. Adopting PRQL/Malloy as an alternate textual surface remains a later option, not this cycle.

**GraphQL as the meta query representation (Hasura/PostGraphile-style).** Rejected, and recorded as settled during the dual-mode work: GraphQL is hierarchical field-selection over a typed resolver graph, lacks standard predicate/sort/aggregate semantics, and needs a schema+resolver server — translating it to each SQL or document backend means building the dalgo AST anyway. GraphQL is a native language DTQL can target, not the abstraction.

## MVP Scope

A timeboxed spike in dalgo that (1) covers the relational, read-only subset of the existing `dal.StructuredQuery` the serve-brokered builder needs first (`From`→source, `Columns`→select, `Where`→filter, `OrderBy`→order, `Args`→parameters), (2) defines the DTQL YAML schema for that subset plus a short purpose/advantages/scope doc, and (3) ships a round-trip test proving `StructuredQuery`→DTQL-YAML→`StructuredQuery` is lossless on that subset. Nested/document shapes, per-driver native rendering, and the saved-.dtql.yaml project-file UX are explicitly deferred to follow-ons once the relational core round-trips. Embarrassing-but-shippable: a tiny AST plus schema plus one green round-trip test beats a complete grammar.

## Not Doing (and Why)

- Bespoke DTQL grammar/parser — YAML is the serialization; PRQL/Malloy adoption is a later option, not MVP
- native→DTQL parsing (lifting arbitrary native SQL/GraphQL back into the AST) — render is one-way only, mirroring the serve-brokered contract
- Nested/document (Firestore-style) AST and DTQL coverage — deferred until the relational core round-trips
- Per-driver native rendering (AST→SQL dialect / document query) — owned by the serve-brokered daemon and dalgo drivers, referenced not built here
- Saved-query project-file UX in DataTug (storing/loading .dtql.yaml) — a DataTug follow-on once the format is stable
- GraphQL-as-meta evaluation — settled (native only); not reopened

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true (confirmed) | dalgo owns a query AST for DTQL to serialize. **Confirmed**: dalgo is SpecScore-managed and its `dal` package already defines `dal.StructuredQuery` (and `dal.TextQuery` for native). | Verified by inspecting `dal-go/dalgo/dal/query.go`; the MVP serializes the existing `StructuredQuery` — no new AST required. |
| Must-be-true | A 1:1 lossless AST↔YAML round-trip is achievable with off-the-shelf YAML for the relational subset (no semantic loss, stable key ordering). | The MVP round-trip test (AST→DTQL-YAML→AST) passes byte-for-AST on a representative relational query set. |
| Should-be-true | DTQL-YAML is human-readable and hand-editable enough that users edit saved queries by hand rather than only through the UI. | Dogfood: hand-edit a saved `.dtql.yaml`, reload it, and confirm the team finds it tolerable versus raw AST/JSON. |
| Should-be-true | One AST renders faithfully to multiple native dialects — the serve-brokered daemon's premise. | Render the MVP subset to ≥2 dialects (e.g. SQLite + Postgres) and compare result sets on a sample database. |
| Might-be-true | DTQL is useful to dalgo consumers beyond DataTug, justifying its home in dalgo. | Sketch one non-DataTug dalgo use case; do not validate this cycle. |


## SpecScore Integration

- **New Features this would create:** in `dal-go/dalgo` (already SpecScore-managed) — a **DTQL** Feature: DTQL YAML schema + (de)serialization of the existing `dal.StructuredQuery` + an AST↔YAML round-trip validator.
- **Existing Features affected:** `datatug`'s `serve-brokered-query-builder` Capability (`REQ:dtql-yaml-form`, `REQ:canonical-ast`) and its `datatug-cli` / `datatug-apps` / `datatug-ai-skills` Implementations — they reference DTQL-YAML and would consume dalgo's definition rather than restating it.
- **Dependencies:** the just-merged serve-brokered dual-mode specs; dalgo's existing `dal` query model (`StructuredQuery` / `TextQuery`).

## Open Questions

- Is this Idea artifact correctly homed in `datatug` (next to the seed and the consuming Capability), or should it be relocated to `dal-go/dalgo` (confirmed SpecScore-managed) where the Feature will live?
- How much of the existing `dal.StructuredQuery` does the relational MVP cover — `GroupBy`, joins/`From` composition, computed columns — before nested/document work begins?
- What is the versioning/compatibility policy for DTQL once queries are saved to disk (schema evolution, forward/backward compatibility)?

---
*This document follows the https://specscore.md/idea-specification*
