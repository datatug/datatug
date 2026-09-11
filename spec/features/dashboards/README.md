---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Dashboards as a surface

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Dashboards are a presentation surface over DataTug projects. DataTug owns durable
board/query contracts, storage, execution and ACL; Dashboardius owns presentation and
editing and can be sold separately. The latest Go board model already preserves
parameters, but Dashboardius's row IDs/widget payloads still differ and its sample
interactions do not establish persisted round-trip compatibility. The bounded real
persistence journey is Phase 5. The query-binding decision it depended on is now made:
the founder ruled 2026-09-09 that a SQL widget references a library query rather than
carrying query text or an execution target of its own.
No dashboards are added to Phase 1; later embedding of the shared renderer is not banned.

Founder ruling 2026-09-11 (verbatim): *"Not sure if it's part of the spec your read but I
think DataTug, Incidentius, Dashbaordius apps code should be in same repo (nx workspace with
shared libs and either separate or single configurable app)"*. The Dashboardius app therefore
moves into the `datatug/datatug-apps` workspace (cutover owned by
[product-profiles](../product-profiles/README.md)); the board renderer becomes a shared library
that DataTug and Incidentius pages embed. Incident and investigation dashboards are specified
by Dashboardius in its `incident-dashboards` Feature; this Feature owns only the DataTug-side
contract they need (REQ:incident-boards-and-shared-renderer).

## Problem

A useful sample-data dashboard does not yet prove a query-backed, shared durable
board. Copied TypeScript types can differ from Go serialization: the current Go row
model does not retain Dashboardius row IDs, and sample SQL widget data is not the
canonical durable query definition. One contract must preserve identifiers, bindings,
parameters, layout and widget configuration through the actual server store.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Board schema (Go + TS) | partial; cross-language acceptance pending | — | enabler |
| Board persistence and query binding through the DataTug server | planned | 5 | strategic |
| Board emerging from exploration ("pin this result") | planned | 5 | strategic |
| Shared board-model package (`@datatug/board-models`) | added in apps; Dashboardius adoption pending | 5 | enabler |
| Embedding shared Dashboardius presentation | undecided; outside Phase 1 | later | surface |

#### REQ: board-schema-shared-source-of-truth

The durable board schema MUST have one declared authority and real server/Go/browser
conformance tests. Preserve canonical card IDs, parameters, requiredParams, layout and widget
configuration. Row identity must remain stable for supported edits after reload, either
through canonical persisted IDs or explicitly defined derived view keys. Migrate the
prototype row IDs deliberately; never silently drop supported durable identity. Runtime query results/provenance are separate from durable query bindings;
saving a board must not serialize cached restricted rows as a portable asset. Unknown
widget compatibility/migration must be explicit, never silent data loss.

#### REQ: datatug-owns-schema-persistence-execution

DataTug MUST own the board schema, board persistence in the project (as a Git file, matching
`git-native-projects`), query execution for bound queries, and access control on every board read
or write. Dashboardius MUST NOT execute a query or read a data source directly; every board
operation Dashboardius performs MUST go through DataTug's server API. This is the split the
2026-09-08 founder direction states and this roadmap section restates as binding: *"DataTug owns
the board schema, board persistence in the project (Git file), query execution and access control,
and the 'pin this result to a board' entry point from exploration."*

#### REQ: dashboardius-owns-presentation-and-editing

Dashboardius owns widget kinds, layout, charts and sharing links, and consumes DataTug only through
the server API: list projects/queries, execute a bound query with parameters, read/write a board.
This Feature does not specify Dashboardius's own UI; it specifies the boundary DataTug's server API
MUST hold so Dashboardius never needs a direct data connection.

#### REQ: board-persistence-and-query-binding

A board MUST be persistable through the DataTug server as a project file, with an
explicit executable binding and parameters for each query-backed widget. That binding
is a reference to a library query, ruled by the founder 2026-09-09 (verbatim): *"I'm
Ok with the suggested option 1"* — a SQL widget MUST NOT carry query text or an
execution target of its own. The query's text and target come from the referenced
`QueryDef`; the environment stays the viewer's choice at view time. `GET/PUT /datatug/boards/*` MUST authorize the project asset and metadata, validate
writes and reject cached protected rows in the durable board. Each widget execution
separately applies source/row/column policies through the common executor. Reading
board JSON does not filter data rows that are not stored in that asset.

#### REQ: board-from-exploration

A user MUST be able to "pin this result" from an exploration screen (a query result, a related-
record lookup) directly to a board, without leaving the investigation to hand-author a board file.
This is DataTug's entry point into Dashboardius, not a DataTug-side charting feature.

#### REQ: shared-board-model-package

`@datatug/board-models` MUST be consumed by both apps, replacing independently evolving
copies. Existing package work is preserved, but adoption alone does not prove Go
compatibility. A test must load a browser-written board through the actual Go store,
save it and reopen in Dashboardius with all supported durable fields intact and deterministic row view identity. Test migration from the prototype shape explicitly.

#### REQ: no-charts-in-datatug-web-ui

Phase 1 MUST NOT add board rendering/layout work. Dashboardius owns the reusable
presentation implementation. A later embedded shared renderer may be specified when
useful; DataTug MUST NOT create an independent second renderer. The REQ ID is retained,
but its earlier permanent embedding prohibition was an unsupported scope inference.

#### REQ: incident-boards-and-shared-renderer

A board MAY carry an optional `incidentRef` (a project-scoped incident id) so that an incident
can open, create and list its boards; access still follows the project's grants, never the
incident. Widget kinds needed by incident dashboards (`metric-series`, comparison, criteria,
replay) are Dashboardius widget *names* inside the existing `{name, data}` envelope, whose
`data` references execution records, compare results and checks by id through the server API
(`evidence-records`, `compare`, `checks`, `incidents`); a board never stores result rows. The
Go validator MUST preserve unknown Dashboardius widget names losslessly (an explicit extension
boundary, not silent acceptance of arbitrary fields). The board renderer used by the
Dashboardius app MUST be the same workspace library a DataTug/Incidentius page embeds as an
island, so one board renders identically in both. Lead assumption: `incidentRef` is a
string field on `Board`, mirrored in `@datatug/board-models`.

## Acceptance Criteria

### AC: schema-published-once-consumed-twice (verifies REQ:board-schema-shared-source-of-truth, REQ:shared-board-model-package)

**Given** `@datatug/board-models` is published
**When** `datatug-apps` and `sneat-co/dashboardius` both import it
**Then** neither repo carries its own copy of the `Board`/`BoardRow`/`BoardCard`/`BoardWidget` types.

### AC: dashboardius-never-queries-directly (verifies REQ:datatug-owns-schema-persistence-execution, REQ:dashboardius-owns-presentation-and-editing)

**Given** a Dashboardius board bound to a DataTug query
**When** the board renders
**Then** Dashboardius calls DataTug's execute endpoint for the data and performs no direct database or Firestore query of its own.

### AC: board-round-trips-through-server (verifies REQ:board-persistence-and-query-binding)

**Given** a board with two widgets bound to `customer-invoices` and `customer-purchases-by-genre`
**When** it is saved via `PUT /datatug/boards/*` and reopened via `GET`
**Then** both widgets' query bindings and parameters are restored exactly, board metadata access is authorized independently, and each widget's re-execution applies that recipient's source/row/column policy. No sender-authorized result rows are restored from board JSON.

### AC: pin-result-creates-board-entry (verifies REQ:board-from-exploration)

**Given** the *Customer invoices* result for `Customer.ID = 5`
**When** the user clicks "Pin to board"
**Then** a new board widget is created bound to that query and parameter value, without leaving the current investigation screen.

### AC: datatug-ui-has-no-chart-widget (verifies REQ:no-charts-in-datatug-web-ui)

**Given** the Phase 1 implementation candidate
**When** its scope is reviewed
**Then** it adds no board/chart implementation; later embedding decisions reuse Dashboardius presentation instead of introducing a second renderer.

### AC: incident-board-round-trips-widget-names (verifies REQ:incident-boards-and-shared-renderer)

**Given** a board with `incidentRef: INC-1` and a `metric-series` widget referencing check `stuck-invoices`
**When** it is saved through `PUT /datatug/boards/*`, loaded by the Go store and reopened in Dashboardius and in the embedded island on the incident page
**Then** `incidentRef`, the widget name and its `data` survive unchanged, no result rows were stored, and both surfaces render the same points fetched from `/datatug/executions`.

## Not Doing / Out of Scope

- Chart/layout implementation in Phase 1. Later embedding may reuse Dashboardius presentation; an independent second renderer remains outside scope.
- A BI/warehouse tool: boards bind to DataTug's own queries and access control, not to a
  centralized warehouse Dashboardius or DataTug maintains separately.
- Dashboardius's own sharing-link or embedding pricing model; this Feature specifies the DataTug-
  side API contract only.
- Resolving Dashboardius's brand-tier question relative to DataTug's Tier-3 status; that is a
  separate, still-open branding decision.

## Open Questions

- ~~Does a Dashboardius board belong to a DataTug project or to a Sneat space?~~ Ruled by the
  founder 2026-09-09 (verbatim): *"board belongs to DataTug project. Project can be shared across
  Sneat spaces"*. Boards are project files; who can see a board follows the project's grants to
  spaces (Feature `sharing-scopes`), never a board-level owner.
- ~~Home of the DataTug↔Dashboardius board contract: publish `@datatug/board-models` from
  `datatug-apps`, or keep the documented copy until Phase 3?~~ (`product/README.md` Open
  Questions.) Resolved by delivery 2026-09-09: `@datatug/board-models` lives in
  `datatug/datatug-apps` under `libs/datatug/board-models` (PR datatug/datatug-apps#49, merged
  2026-09-09), with types generated field-for-field from `datatug-core/pkg/datatug/boards.go`.
  `datatug-apps`' own pre-existing copy of the schema in `libs/datatug/main` is being removed on
  branch `feat/board-models-reconcile` of the same repo. `sneat-co/dashboardius` still has to
  switch to the package (see that package's README, "What Dashboardius must change"), so
  AC:schema-published-once-consumed-twice is not yet satisfied.
- ~~**Does `Board` need `parameters`?**~~ `boards.go`'s `Board` was `ProjectItem` + `Rows` and
  carried no parameters, while `ProjBoardBrief` in the same file carries `Parameters` and
  `RequiredParams`. A `board.json` holding `parameters` was therefore silently dropped on a
  `LoadBoard`/`SaveBoard` round trip: Go's JSON decoder ignores the unknown field and the encoder
  never writes it back. **Ruled by the founder 2026-09-09 (verbatim): *"1 - yes"*** — `Board` MUST
  carry them. `Parameters` and `RequiredParams` are added to `Board`, mirroring `ProjBoardBrief`
  field for field, so a board and its brief stay in step and a saved board keeps its parameters.
- ~~**How is a board SQL widget executed?**~~ A board SQL widget used to store
  `SQLWidgetSettings{ Query string }` and nothing else: it held the SQL text but named no database
  and no environment, so nothing could decide where to run it. The TypeScript copy had invented
  `db` and `env` on the widget; both were dead code and were dropped in 2026-09-09's
  reconciliation rather than added to Go.

  **Ruled by the founder 2026-09-09 (verbatim): *"I'm Ok with the suggested option 1"*** — the
  widget references a library query. `SQLWidgetSettings` carries a `queryId` plus parameter
  bindings and no query text or execution target of its own; the text and target come from the
  referenced `QueryDef` (which already carries `Targets`: driver, catalog, protocol, host, port),
  and the environment stays the viewer's choice at view time, as the board page's environment
  selector and the per-call execute API already assume. This keeps one definition of every query
  and makes a board's access control identical to its queries' access control.

  Accepted costs, stated when the option was chosen: any inline-SQL widget needs migrating, and a
  board cannot hold an ad-hoc query that is not in the library. No migration is required inside
  this org's own data — the only board in `datatug-demo-projects` has no widgets at all.

  The wire shape (`queryId`, plus bindings that take either a constant `value` or a
  `boardParameterId`) is the coding agent's design expressing that ruling, not the founder's own
  words, and stays open to review.

---
*This document follows the https://specscore.md/feature-specification*
