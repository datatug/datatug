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

Promotes the datatug sidekick seed at spec/ideas/seeds/dtql-datatug-query-lanaguage-needs-definition-and.md. The following serve-brokered references and implementation assumptions record the context of this Idea's 2026-06-05 approval; they are not a current implementation receipt. At that time, the serve-brokered query builder (Capability in datatug plus Implementations in datatug-cli/datatug-apps/datatug-ai-skills) referenced DTQL-YAML as the textual form of the dalgo AST: REQ:dtql-yaml-form mandates a lossless AST↔YAML round-trip, and REQ:canonical-ast renders that AST to each connection's native query language. But DTQL is referenced, not defined — no schema, no validator, no owning artifact. The 2026-06-05 recommendation below selected YAML serialization for the original format/MVP scope. The user-approved 2026-10-09 [Query Workspace](../features/query-workspace/README.md) adds the required TugQL readable text form and TugQTree structured YAML form over the same query model. Go and TypeScript parser implementations remain planned, not delivered, and this direction neither changes existing YAML documents nor authorizes an umbrella name change, website retirement, redirect, or migration: DTQL remains the library and website name, with any broader naming or old-site hygiene question reserved for separate future work. The earlier dual-mode conclusion that GraphQL is a native query language rather than the structured AST/meta representation remains about that representation boundary; it does not prohibit another syntax for authoring the structured AST. Crucially, `github.com/dal-go/dalgo` is **already SpecScore-managed** and its `dal` package **already models the query**: `dal.StructuredQuery` (`From`/`Where`/`GroupBy`/`OrderBy`/`Columns`/`StartFrom`) is the structured AST DTQL serializes, and `dal.TextQuery` (`Text` + `Args`) already represents the native-text form — a near-exact match to the serve-brokered native/DTQL mode split. So this Idea serializes an existing model rather than inventing one.

## Recommended Direction

The original 2026-06-05 recommendation was to define DTQL in the dalgo repo as the canonical, human-readable, 1:1 lossless YAML serialization of dalgo's query AST — structural, hierarchical, plain off-the-shelf YAML with no bespoke grammar or parser. dalgo **already owns the query AST** — its `dal.StructuredQuery` interface (`From`/`Where`/`GroupBy`/`OrderBy`/`Columns`/`StartFrom`); this Idea adds the DTQL YAML schema, (de)serialization, and a round-trip validator proving `StructuredQuery`→YAML→`StructuredQuery` identity. dalgo's existing `dal.TextQuery` (`Text` + `Args`) already represents the native-text form, so DTQL serializes the *structured* side of a split dalgo already makes. DataTug consumes DTQL two ways: as the serve-brokered builder's internal exchange form (the agent and web UI inspect and edit it) and as a user-facing saved query format (.dtql.yaml stored in projects, version-controlled, diffable, shared). Because the AST spans relational and nested/document (Firestore-style) shapes, DTQL is not SQL-shaped; it renders to each connection driver's native query language on execution. The 2026-06-05 decision treated GraphQL as a native language DTQL can target, not a competing structured meta representation — a GraphQL-as-meta approach would still have to build the dalgo AST underneath. The later [Query Workspace](../features/query-workspace/README.md) may support GraphQL import/native mode under its own schema, resolver, authorization and execution gates. Advantages over the alternatives at the serialization layer: structure, hierarchy, and human-readability for free via YAML; deterministic point-and-click editability in the web UI because DTQL is the AST; one model rendered to many dialects; saved, versioned, diffable queries as plain text; and zero new parser to design or maintain.

## Alternatives Considered

**Define DTQL in `datatug-go-models` (or `datatug`) instead of dalgo.** Keeps the schema next to DataTug's usage and avoids touching the dalgo library, but confines a generally useful query serialization to DataTug and splits the AST (the model) from its serialization across repos. Rejected in favour of dalgo so the query model and its DTQL form live together and any dalgo consumer benefits. The original Must-be-true assumption was that dalgo owns a query AST; DALgo's verified `dal.StructuredQuery` contract satisfies that assumption.

**Invent a bespoke DTQL textual grammar (PRQL/Malloy-style) instead of YAML.** A purpose-built syntax reads more fluently and is terser, but costs a parser and grammar to design, test, version, and keep in lock-step with the AST. The 2026-06-05 MVP recommendation set YAML as the structured serialization format; it did not decide whether a separate authoring frontend should exist. The 2026-10-09 [Query Workspace](../features/query-workspace/README.md) now requires TugQL as a readable text form and TugQTree as the structured YAML form for the same model. Its spike selects the bounded grammar and parser architecture, not whether a text surface ships. Existing DTQL YAML remains compatible, and DTQL stays the library and website name; no umbrella rename, site retirement, redirect, or migration is authorized.

**GraphQL as the meta query representation (Hasura/PostGraphile-style).** The 2026-06-05 dual-mode decision set aside GraphQL as the structured AST/meta representation: translating it to each SQL or document backend still needs a canonical query model. This does not rule out GraphQL as a gated native query/import mode, nor does it rule out a distinct textual authoring frontend that maps to the DALgo AST.

## MVP Scope

The original 2026-06-05 MVP recommendation was a timeboxed spike in dalgo covering the relational, read-only subset of the existing `dal.StructuredQuery` needed first, defining its DTQL YAML schema and round-trip validation. That history documents the original format scope; it is not evidence that each proposed capability shipped. The 2026-10-09 [Query Workspace Plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/query-workspace/README.md) separately requires a TugQL authoring frontend and TugQTree YAML compatibility over the canonical model; implementation remains pending. Nested/document shapes, per-driver native rendering, and saved-query UX remain owned by their current Features and Plans.

## Not Doing (and Why)

- Replacing DTQL YAML serialization with a bespoke grammar — the 2026-06-05 recommendation retained YAML; Query Workspace adds TugQL alongside compatible TugQTree YAML rather than replacing it
- General native→DTQL parsing (lifting arbitrary native SQL/GraphQL back into the AST) is outside this Idea's serialization scope; the approved [Query Workspace](../features/query-workspace/README.md) separately evaluates constrained import with loss diagnostics and no automatic execution grant
- Nested/document (Firestore-style) AST and DTQL coverage — deferred until the relational core round-trips
- Per-driver native rendering (AST→SQL dialect / document query) — owned by the serve-brokered daemon and dalgo drivers, referenced not built here
- Saved-query project-file UX in DataTug (storing/loading .dtql.yaml) was outside this original Idea; the [Query Workspace](../features/query-workspace/README.md) handles its own saved-query contract while preserving the established serialization format
- Replacing the canonical structured AST with GraphQL — the 2026-06-05 decision set that aside; Query Workspace may still evaluate gated GraphQL native/import support

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true (confirmed) | dalgo owns a query AST for DTQL to serialize. **Confirmed in the 2026-06-05 review**: the `dal` package defines `dal.StructuredQuery` and `dal.TextQuery`. | The original recommendation serialized the existing `StructuredQuery`; a later authoring frontend must continue to map to that canonical model. This historical statement does not attest to current implementation status. |
| Must-be-true | A 1:1 lossless AST↔YAML round-trip is achievable with off-the-shelf YAML for the relational subset (no semantic loss, stable key ordering). | The MVP round-trip test (AST→DTQL-YAML→AST) passes byte-for-AST on a representative relational query set. |
| Should-be-true | DTQL-YAML is human-readable and hand-editable enough that users edit saved queries by hand rather than only through the UI. | Dogfood: hand-edit a saved `.dtql.yaml`, reload it, and confirm the team finds it tolerable versus raw AST/JSON. |
| Should-be-true | One AST renders faithfully to multiple native dialects — the serve-brokered daemon's premise. | Render the MVP subset to ≥2 dialects (e.g. SQLite + Postgres) and compare result sets on a sample database. |
| Might-be-true | DTQL is useful to dalgo consumers beyond DataTug, justifying its home in dalgo. | Sketch one non-DataTug dalgo use case; do not validate this cycle. |


## SpecScore Integration

- **New Features this would create:** in `dal-go/dalgo` (already SpecScore-managed) — a **DTQL** Feature: DTQL YAML schema + (de)serialization of the existing `dal.StructuredQuery` + an AST↔YAML round-trip validator.
- **Existing Features affected:** `datatug`'s `serve-brokered-query-builder` Capability (`REQ:dtql-yaml-form`, `REQ:canonical-ast`) and its `datatug-cli` / `datatug-apps` / `datatug-ai-skills` Implementations — they reference DTQL-YAML and would consume dalgo's definition rather than restating it.
- **Dependencies:** the just-merged serve-brokered dual-mode specs; dalgo's existing `dal` query model (`StructuredQuery` / `TextQuery`).

## Open Questions

- Is this Idea artifact correctly homed in `datatug` (next to the seed and consuming Features), or should it be relocated to `dal-go/dalgo` (confirmed SpecScore-managed) where serialization implementation would live?
- How much of the existing `dal.StructuredQuery` does the relational MVP cover — `GroupBy`, joins/`From` composition, computed columns — before nested/document work begins?
- What is the versioning/compatibility policy for DTQL once queries are saved to disk (schema evolution, forward/backward compatibility)?

---
*This document follows the https://specscore.md/idea-specification*
