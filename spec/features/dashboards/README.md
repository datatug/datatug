---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Dashboards as a surface

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dashboards?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P15 is dashboards as a surface: boards emerge from exploration and are preserved as project files, with DataTug owning the board schema, persistence and query execution while the separate, standalone Dashboardius product owns presentation and editing.

Founder direction (2026-09-08 Idea, verbatim): *"A Dashboardius board IS a DataTug board; what Dashboardius adds is widget kinds and an editing experience."* Audited state on 2026-09-09: the `Board`/`BoardRow`/`BoardCard`/`BoardWidget` schema is usable as a Go+TS model, copied field-for-field into Dashboardius; board persistence and query binding through the DataTug server is unimplemented (`DatatugBoardService.getBoard()` is not implemented); Dashboardius itself is a live, separate Sneat product on sample data only. Phase 3 builds board persistence, the "pin this result" entry point, and the shared `@datatug/board-models` package; DataTug's own web UI deliberately never grows chart rendering. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`, `sneat-co/dashboardius`) link to it.

## Problem

A board that cannot be traced back to the queries and access policy that produced it is just a
picture; DataTug's thesis requires a board to stay a first-class, executable, Git-native asset like
every other library item. The schema for this already exists and is genuinely shared —
Dashboardius "re-exports DataTug's own `Board`/`BoardRow`/`BoardCard`/`BoardWidget` schema
field-for-field into a PrimeNG app with Chart.js" (current-state audit, verdict line 9) — but the
DataTug side of the contract, board persistence through the server, does not exist:
`DatatugBoardService.getBoard()` throws not-implemented, and boards in the web UI hard-code
`localhost:8989` (current-state audit, capability table, P15 rows). Dashboardius is live at
dashboardius.com on sample data only, "while DataTug's own web UI board rendering remains
unimplemented" — so the product that is supposed to present DataTug boards has nothing real to
present yet. Without a published `@datatug/board-models` package, the schema also risks drifting
between the two repos as a documented copy rather than a single source of truth.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Board schema (Go + TS) | usable (model) | — | enabler |
| Board persistence and query binding through the DataTug server | planned | 3 | strategic |
| Board emerging from exploration ("pin this result") | planned | 3 | strategic |
| Shared board-model package (`@datatug/board-models`) | planned | 3 | enabler |
| Charts in DataTug's own web UI | planned | — | — |

#### REQ: board-schema-shared-source-of-truth

The `Board`/`BoardRow`/`BoardCard`/`BoardWidget` schema MUST have exactly one source of truth. This
capability is already usable as a Go+TS model; the risk this REQ guards against is drift between
DataTug's copy and Dashboardius's field-for-field re-export, which REQ:shared-board-model-package
below resolves.

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

A board MUST be persistable through the DataTug server as a project file, with each widget bound to
a library query and its parameters. `GET/PUT /datatug/boards/*` and an execute endpoint MUST apply
the same access control as every other server read (`server-enforced-access` REQ:serve-acl-all-
reads) — a board never bypasses row/column policy because it is "just a widget."

#### REQ: board-from-exploration

A user MUST be able to "pin this result" from an exploration screen (a query result, a related-
record lookup) directly to a board, without leaving the investigation to hand-author a board file.
This is DataTug's entry point into Dashboardius, not a DataTug-side charting feature.

#### REQ: shared-board-model-package

The board schema MUST be published as `@datatug/board-models`, consumed by both `datatug-apps` and
`sneat-co/dashboardius`, removing Dashboardius's current documented copy. Until this package is
published, any schema change MUST be applied to both repos in the same session to avoid drift.

#### REQ: no-charts-in-datatug-web-ui

DataTug's own web UI MUST NOT grow chart rendering or board-layout editing; Dashboardius presents,
DataTug executes and persists. This is an explicit non-goal, not a temporary gap — the guard rail
the roadmap states directly: *"Phase 1 does not build boards. Dashboardius work must not pull chart
or layout work into DataTug's web UI."*

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
**Then** both widgets' query bindings and parameters are restored exactly, and a restricted principal's board read is limited the same way `core-investigation-loop` AC:restricted-rows-and-columns-server limits a direct query run.

### AC: pin-result-creates-board-entry (verifies REQ:board-from-exploration)

**Given** the *Customer invoices* result for `Customer.ID = 5`
**When** the user clicks "Pin to board"
**Then** a new board widget is created bound to that query and parameter value, without leaving the current investigation screen.

### AC: datatug-ui-has-no-chart-widget (verifies REQ:no-charts-in-datatug-web-ui)

**Given** DataTug's own web UI codebase
**When** it is checked for a charting library (`echarts`/`ngx-charts`/`chart.js`/`d3`) or a board-layout editor
**Then** none is found, matching the current-state audit's finding and the guard rail this REQ states.

## Not Doing / Out of Scope

- Charts, layout editing, or widget-kind design inside DataTug's own web UI — permanently, not just
  in Phase 1; that is Dashboardius's scope.
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
- **How is a board SQL widget executed?** *Open. Raised by the coding agent on 2026-09-09; the
  founder asked for more context before deciding, which this entry supplies. Nothing below is a
  founder ruling.*

  **The gap.** A board SQL widget stores `SQLWidgetSettings{ Query string }` and nothing else. It
  holds the SQL text but names no database and no environment, so nothing can decide where to run
  it. `ParameterLookup` right beside it in `boards.go` does carry a `db`, which makes the omission
  look accidental rather than deliberate. The TypeScript copy had invented `db` and `env` on the
  widget; both were dead code and were dropped in 2026-09-09's reconciliation rather than added to
  Go, precisely because inventing a wire field ahead of this decision is what caused the drift.

  **What already exists, and why it matters.** A saved library query is not missing this. `QueryDef`
  (`datatug-core/pkg/datatug/query.go`) carries `Type`, `Text`, `Parameters` and
  `Targets []QueryDefTarget`, where a target is driver, catalog, protocol, host, port and
  credentials. So DataTug already has a modelled answer to "where does this query run" — just not
  on boards. Meanwhile the execute API (`SelectRequest`) takes `Project`, `Environment`, `Database`
  and `SQL` per call, and the board page already renders an environment selector, so environment is
  in practice chosen by the viewer at view time, not stored.

  **Why this is really a binding question.** REQ:board-persistence-and-query-binding already says
  each widget binds to a library query plus parameters. A widget holding inline SQL text is that
  requirement's opposite: it duplicates the query instead of referencing it, which is why it has
  nowhere to inherit a target from. Deciding the execution target and deciding whether a widget
  references a query are the same decision.

  **Options.**
  1. *Widget references a library query.* `SQLWidgetSettings` becomes a query reference plus
     parameter bindings; the target comes from the referenced `QueryDef.Targets`, and environment
     stays a view-time choice. Satisfies the existing REQ, keeps one definition of every query, and
     makes board access control identical to query access control. Costs a migration for any
     inline-SQL widget and makes a board unable to hold an ad-hoc query.
  2. *Board declares the target, widgets inherit it.* One db/env per board, alongside the board's
     parameters. Simple, matches how a viewer thinks about a dashboard, and needs no query library.
     Prevents one board from drawing on two databases.
  3. *Widget declares its own db, environment stays view-time.* Closest to what the TypeScript copy
     had invented. Most flexible, but re-introduces the duplication the reconciliation removed and
     leaves each widget's target unvalidated against the project's environments.

  **What is blocked until this is answered.** DataTug's own web UI renders a SQL widget's query
  text and never executes it. AC:board-round-trips-through-server cannot be written against a real
  binding, and AC:pin-result-creates-board-entry has nowhere to record which database the pinned
  result came from.

---
*This document follows the https://specscore.md/feature-specification*
