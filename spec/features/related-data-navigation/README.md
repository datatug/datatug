---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Relationships and related-data navigation

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/related-data-navigation?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/related-data-navigation?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/related-data-navigation?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/related-data-navigation?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P5 is related-data navigation: a selected value becomes a navigation point to related records, across sources, using declared foreign keys and semantic mappings, beyond what a single table's schema can show on its own.

It is one half of the concrete evidence for the thesis's central claim that "a value on screen is a navigation point" — the other half is P6's Investigation Context, which carries a selected value forward. Audited state on 2026-09-09: the foreign-key graph is usable and scanned bidirectionally; the web UI's FK-driven related cards are real, working navigation components, but they depend on the agent connection that is broken end to end (P11); there is no runtime "fetch rows related to this value" execution anywhere in `datatug-core` — only a value-picker shape, `ParameterLookup`. Phase 1 (`core-investigation-loop`) builds the missing server-executed lookup and the cell-to-context-panel flow for the demo project's two sources; Phase 5/X considers automatic relationship discovery beyond declared foreign keys. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`, `datatug-apps`) link to it.

## Problem

A record rarely stands alone: a customer has invoices, an invoice has line items, a customer might
have support notes in an entirely different, non-relational source. Seeing only the current table
forces the user to already know which other tables to check and to write the join or lookup by
hand. DataTug's schema scan already builds a complete FK graph (`ForeignKey`/`ReferencedBy`,
`pkg/schemer/scan_constraints.go:57-100`) and the web UI already has real click-through FK cards
(`env-db-table.page.html` "Record" tab, `foreign-key-card`/`ref-by-card` components) — but the audit
found "no runtime 'fetch related rows for this value' query execution" anywhere in `datatug-core`
(a `grep` for "related" across `pkg/**` returned zero hits beyond FK/schema plumbing), and the FK
cards themselves "depend on the broken agent path" (current-state audit, P5 capability rows). The
TUI's equivalent panels are read-only, with no click-to-navigate (`table_fks.go`, `table_referrers.go`
— real data, no `SetSelectedFunc`). None of this reaches across a source boundary at all today.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Foreign-key graph scanned bidirectionally | usable | 1 | enabler |
| FK-driven related cards on the table page (web) | usable UI / broken data path | 0–1 | strategic |
| Semantic relationships across sources (same semantic field in two sources) | planned | 1 | strategic |
| Related-record lookup executed by the server for a selected value | planned | 1 | strategic |
| Cell/value selection → context panel (meaning, related, applicable, actions) | planned | 1 | strategic |
| Automatic relationship discovery (beyond declared FKs) | planned | 5 / X | strategic |

#### REQ: fk-graph-bidirectional

Every scanned SQL source MUST expose both directions of a foreign key: the referencing side and the
referenced-by side, so a related-data lookup can be built without re-deriving direction from raw
constraint metadata each time. This capability is already usable; this REQ pins it as a standing
guarantee every new adapter or surface can rely on.

#### REQ: fk-driven-related-cards

The web UI's table page MUST show related-record cards driven by the FK graph once the underlying
agent connection is fixed (see `surfaces` REQ:agent-connection-fixed). The card components
themselves (`foreign-key-card`, `ref-by-card`) are real; this REQ requires them to be backed by
live server data, not the placeholder text they show today.

#### REQ: semantic-relationships-across-sources

A relationship MUST also be derivable from two entity fields sharing the same semantic identity in
different sources, not only from a declared SQL foreign key — this is how a `Customer.ID` in
`chinook.sqlite` reaches a `support-notes` record in an inGitDB collection with no FK between the
two. This is specified by `core-investigation-loop` REQ:related-lookup-model ("semantic mappings of
the same field in other sources"); no new graph type is needed for the first cut.

#### REQ: related-lookup-server-executed

A related-record lookup MUST be executed by the server as a DALgo `StructuredQuery`, never
constructed or executed by the browser. This is specified in full by `core-investigation-loop`
REQ:related-lookup-execution, including that it runs through the same access-policy path as every
other read.

#### REQ: cell-selection-context-panel

Selecting a value (a grid cell, a TUI list item) MUST open a context panel showing the value's
meaning (P3/P4), its related records (this Feature) and the library queries it makes applicable
(P7) — the single interaction that ties the pillars together for the user. Phase 1 builds the web
version of this panel over the demo project's two sources; parity on the TUI is Phase 2.

#### REQ: automatic-relationship-discovery

Discovering relationships beyond declared foreign keys or semantic mappings — for example by
sampling value overlap between two columns — MUST stay an experiment validated for precision before
it becomes a shipped feature, and MUST always be distinguishable from a declared or FK-derived
relationship. It MUST NOT silently merge into the declared relationship graph.

## Acceptance Criteria

### AC: fk-graph-both-directions (verifies REQ:fk-graph-bidirectional)

**Given** the demo project's `chinook.sqlite` schema, where `Invoice.CustomerId` references `Customer.CustomerId`
**When** the FK graph is scanned
**Then** `Customer` shows `Invoice` under referenced-by and `Invoice` shows `Customer` under foreign keys.

### AC: related-cards-show-live-data (verifies REQ:fk-driven-related-cards)

**Given** the `Customer` table page in the web UI, agent connection fixed
**When** a customer row is opened
**Then** the related-record cards show real invoice counts from the server, not the current placeholder text.

### AC: cross-source-relationship-resolved (verifies REQ:semantic-relationships-across-sources)

**Given** `Customer.ID` is mapped in both `chinook.sqlite` and the `support-notes` inGitDB collection
**When** `Customer.ID = 5` is selected
**Then** the related-lookup list includes `Support notes` from `support-notes` alongside `Invoices` from `chinook`, matching `core-investigation-loop` AC:related-across-sources.

### AC: lookup-runs-through-policy-path (verifies REQ:related-lookup-server-executed)

**Given** the `support` principal with a Canada-only row policy on `Customer`
**When** the related *Invoices* lookup runs for a Brazilian customer
**Then** it returns zero rows with the limitation shown, matching `core-investigation-loop` AC:related-respects-policy — never a browser-side filter of a full result.

### AC: context-panel-shows-all-three (verifies REQ:cell-selection-context-panel)

**Given** a `CustomerId` cell with value 5 is selected on the *Customer* table
**When** the context panel opens
**Then** it shows the semantic marker (`Customer.ID · declared`), related records from both sources, and the applicable queries with their resolution chain, in one panel.

### AC: discovered-relationship-labelled (verifies REQ:automatic-relationship-discovery)

**Given** an experimental value-overlap scan proposes a relationship with no declared FK or semantic mapping backing it
**When** the proposal is shown anywhere
**Then** it is labelled `discovered` and distinguishable from `declared`/`fk`/`semantic`, and it is never auto-merged into the declared relationship graph.

## Not Doing / Out of Scope

- Cross-source joins: only lookups (fetch related rows for a value) are in scope; a real join
  engine is deferred until a real investigation needs one (Phase 5, and only if justified).
- A generic workflow engine: navigating related data is a human-guided click, not an automated
  traversal.
- Automatic relationship discovery reaching production before its precision is validated on the
  demo project as an explicit experiment.
- A data catalog's relationship diagram as a governance artifact; DataTug's relationship graph
  exists to be navigated, not documented for its own sake.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
