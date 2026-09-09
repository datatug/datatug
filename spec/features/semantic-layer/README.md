---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Semantic entities and fields

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-layer?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-layer?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-layer?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-layer?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P3 is DataTug's semantic layer: entities and typed fields distinct from tables and columns, with explicit, inspectable mappings from physical locations, and semantic tags carried on query parameters and result columns.

It is the pillar that turns "a column named `CustomerId`" into "`Customer.ID`, the same concept everywhere it appears" — the foundation P4 (inference), P5 (related-data navigation), P6 (investigation context) and P7 (the knowledge library) all build on. Audited state on 2026-09-09: the `Entity`/`EntityField` model is usable and well tested (`datatug-core/pkg/datatug/entities.go`, 99%+ coverage), `datatug entity` CLI and the web entities CRUD pages are usable, and `EntityFieldRef` semantic tagging exists on parameters and result columns — but nothing at runtime resolves a mapping or reads those tags; the model is data without a resolver. Phase 1 (`core-investigation-loop`) builds the missing field→column mapping resolver and the endpoint that reads `EntityFieldRef`; Phase 2 adds richer semantic field types. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`, `datatug-apps`) link to it.

## Problem

A physical schema (tables, columns, foreign keys) only tells you what storage looks like; it does
not tell you that `CustomerId` in one table and `cust_id` in another mean the same business
concept, or that a query parameter named `id` requires a `Customer.ID` value. Without an explicit
semantic layer, every screen and every query has to re-derive meaning from naming conventions, and
an AI agent has nothing more reliable than the same guesswork. DataTug's `Entity`/`EntityField`
model already captures typed fields and entity→table mappings (`datatug-core/pkg/datatug/
entities.go`; `Entity.Tables []DBCollectionKey`), and `ParameterDef.Meta`/`RecordsetColumnDef.Meta`
already carry an `EntityFieldRef` tag on parameters and result columns
(`datatug-core/pkg/datatug/parameters.go:20`; `recordset_def.go:65`) — but the audit found "no code
in `datatug-core` executes pattern matching against a scanned schema" and nothing reads those
`Meta` tags at runtime (current-state audit, verdict line 4). The demo project itself is short two
entities a customer-support pitch needs: `Customer` and `Employee` have tables but no `Entity`
wrapper yet (hierarchy.md P3 row 6).

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Entities with typed fields and key fields | usable | 1 | strategic |
| Entity → tables mapping | usable (model) | 1 | strategic |
| Field → column mapping (explicit, per source) | planned | 1 | strategic |
| Semantic tags on parameters and result columns (`EntityFieldRef`) | implemented-but-poorly-exposed | 1 | strategic |
| Semantic field types (money, email, GUID… ; "currency" seed) | partial | 2 | enabler |
| Entity summary in the demo project covers Customer/Employee | planned | 1 | hygiene |

#### REQ: entities-typed-fields

A project MUST be able to declare an entity with typed fields and designated key fields, editable
through both the CLI (`datatug entity`) and the web UI, staying consistent between the two. This
capability is already usable; this REQ pins it as a standing guarantee, not new behaviour.

#### REQ: entity-table-mapping

An entity MUST be able to declare which physical tables/collections it draws from
(`Entity.Tables`), across more than one source. This mapping is the entry point for P5's
cross-source related-data navigation: two entities of the same field in different sources are how a
lookup crosses a source boundary.

#### REQ: field-column-mapping

An entity field MUST be able to declare an explicit mapping to a physical column, per source. This
is specified in full by `core-investigation-loop` REQ:field-mapping-model (declared mappings are
authoritative; `NamePatterns` may derive an inferred mapping and must label it as such); this
Feature asserts it as the P3 capability without restating the resolver mechanics.

#### REQ: semantic-tags-resolved-at-runtime

`EntityFieldRef` tags on `ParameterDef.Meta` and `RecordsetColumnDef.Meta` MUST be read by a runtime
resolver, not merely stored. This is specified by `core-investigation-loop` REQ:semantic-resolution-
endpoint (the server-side endpoint every surface must use) and REQ:parameter-auto-binding (how a
tagged parameter gets filled); this Feature asserts the tags themselves are a P3 concern shared by
every surface, CLI and TUI included, not only the web endpoint Phase 1 builds first.

#### REQ: semantic-field-types

DataTug SHOULD offer semantic field types beyond a bare entity-field name — money, email, GUID and
similar `KnownTypes` — so a mapping can carry not just "this is `Customer.ID`" but "this is
formatted and validated as a GUID". Field types MUST remain inspectable and editable, never a
silent coercion of stored values.

#### REQ: demo-entity-coverage

The demo project's entity summary MUST cover the tables a customer-support investigation actually
touches. `core-investigation-loop` REQ:demo-knowledge-project adds `Customer` and `Invoice`;
`Employee` remains a gap this Feature records but does not resolve in Phase 1 — a future task MUST
add an `Employee` entity with a declared mapping before any investigation journey uses it.

## Acceptance Criteria

### AC: entity-crud-consistent (verifies REQ:entities-typed-fields)

**Given** the demo project
**When** an entity field is added via `datatug entity field` and then viewed in the web entities page
**Then** the field's type and key-field status match in both surfaces.

### AC: entity-spans-multiple-tables (verifies REQ:entity-table-mapping)

**Given** the `Customer` entity
**When** its `Tables` list is inspected
**Then** it names the `chinook.Customer` table, and a second source's equivalent collection can be added to the same list without creating a second entity.

### AC: declared-mapping-wins-over-pattern (verifies REQ:field-column-mapping)

**Given** `Customer.ID` has a declared mapping to `chinook.Customer.CustomerId` and also matches a `NamePatterns` rule
**When** the semantic endpoint resolves the column
**Then** the response is `provenance: declared`, matching `core-investigation-loop` AC:mapped-columns-server.

### AC: tags-drive-parameter-binding (verifies REQ:semantic-tags-resolved-at-runtime)

**Given** `customer-invoices`'s `Customer.ID` parameter carries `Meta: EntityFieldRef{Customer,ID}`
**When** `Customer.ID = 5` is selected on the grid
**Then** the parameter is pre-filled, matching `core-investigation-loop` AC:bound-from-selection.

### AC: semantic-type-visible-not-coercive (verifies REQ:semantic-field-types)

**Given** `Customer.Email` is tagged with the `email` semantic type
**When** the field is shown in the entities page
**Then** the type is visible as metadata and the underlying stored value is unchanged.

### AC: employee-entity-gap-recorded (verifies REQ:demo-entity-coverage)

**Given** the demo project after Phase 1 ships
**When** `datatug entity list` runs
**Then** `Customer` and `Invoice` appear with declared mappings and `Employee` is absent, matching the recorded gap rather than a silent claim of completeness.

## Not Doing / Out of Scope

- A data catalog or governance suite: entities capture meaning to navigate and execute, not to
  certify or govern data.
- Probabilistic or automatic meaning inference: P3 only covers explicit, declared structure; see
  `meaning-inference` (P4) for suggestions and correction workflows.
- Connector breadth as a growth lever: entity mapping depth is proven on the sources already
  wired (SQL, inGitDB) before it is extended to every possible backend.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
