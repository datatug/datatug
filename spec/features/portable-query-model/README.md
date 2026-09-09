---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Portable query model (DTQL + DALgo)

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P1 is the portable query model: one logical query AST (DALgo `StructuredQuery`), serialized losslessly as DTQL YAML, rendered to each backend's native query language, with native text kept as an explicit escape hatch.

It is the pillar every other pillar leans on: P10's server-enforced ACL can only rewrite a query it can parse as structure, and P7's shared library can only diff and review a query in Git if it is data, not opaque text. Audited state on 2026-09-09: DTQL serialization and the DALgo `access` policy engine are Stable and well tested in `dal-go/dalgo`, but DataTug itself consumes them at exactly one call site, `datatug query run` (`datatug-cli/apps/datatugapp/commands/cmd_query.go`); `QueryDef.Type` has no `DTQL` value yet (`datatug-core/pkg/datatug/query.go:124-137`); and the standalone `dtql`/`dtql-go`/`dtql-backstage` repositories and the `dtql.org` site are obsolete stubs that contradict the real implementation. Phase 0 retires the obsolete repos; Phase 1 (`core-investigation-loop`) adds `DTQL` as a saved-query type; Phase 2 completes the HTTP escape hatch and the remaining backend adapters. This is a Capability; per-repo Implementations (`datatug-cli`, `dal-go/dalgo`) link to it.

## Problem

Query text is source-language-specific: T-SQL, a Firestore query, an inGitDB filter and an HTTP
request body share no common shape. A query written for one source cannot move to another, cannot
be structurally rewritten by the server to enforce row/column policy, and cannot be diffed
meaningfully in Git the way a data file can. DataTug's thesis depends on queries being executable
knowledge assets, not opaque strings. Today that promise is only kept at one call site: DTQL
serialization is "usable" and the DALgo `access` package is "usable, spec Stable"
(`dal-go/dalgo/dtql/README.md`; `dal-go/dalgo/access/*`), but `QueryDef.Type` in DataTug's own
project model still only accepts `SQL`/`HTTP`/`StructuredSQL`
(`datatug-core/pkg/datatug/query.go:124-137`), so no project file can hold a DTQL document. Backend
rendering is likewise partial: `sqlite://` and `ingitdb://` are fully wired but `postgres://`
`Open()` always returns `ErrPostgresNotWired` (`pkg/dbcopy/url.go:34-37,164-165`), and `mssql`/
`mysql`/`firestore` are not in the dispatcher at all. Marketing on `datatug.io` and the standalone
`dtql.org` site present DTQL as an installable product independent of `dal-go/dalgo`; the audit
found the `dtql`/`dtql-go`/`dtql-backstage` repos are one-line-README stubs with a single 2026-06-09
commit and zero `.go` files, while the real implementation lives under a different module path
entirely (current-state audit, "Contradictions" table, DTQL row).

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| DTQL YAML serialization of the DALgo query AST | usable (external, Stable in `dal-go/dalgo/dtql`) | — | enabler |
| DTQL as a saved-query type in projects (`queries/**/*.dtql.yaml`, `QueryDef.Type = DTQL`) | planned | 1 | strategic |
| Native query escape hatch (SQL text, HTTP request) | usable (SQL) / backend-only (HTTP) | 1 (SQL), 2 (HTTP) | enabler |
| Rendering to backends through DALgo adapters | partial | 0–1 | enabler |
| Standalone `dtql*` repositories and dtql.org claims | obsolete / contradicted | 0 | hygiene |

#### REQ: dtql-ast-consumption

DataTug MUST treat `dal-go/dalgo`'s `StructuredQuery` AST and its DTQL YAML serialization as the
single source of truth for a portable query; DataTug repos MUST NOT fork or reimplement the
parser, serializer or the DALgo `access` policy evaluator. Any DataTug-side query representation
(saved query, related-record lookup, applicable-query resolution) SHOULD be expressible as this AST
so the server can apply row/column policy uniformly (see `server-enforced-access`).

#### REQ: dtql-saved-query-type

A project MUST be able to store a query as a `.dtql.yaml` sidecar with `QueryDef.Type = DTQL`,
loaded and executed the same way as a `SQL`-type query. This is specified in full by
`core-investigation-loop` REQ:dtql-query-type; this Feature only asserts it as a P1 capability and
does not restate the mechanics.

#### REQ: native-query-escape-hatch

Where a source's query language cannot be fully expressed as a `StructuredQuery` (a stored
procedure call, a complex HTTP request, a hand-tuned index hint), a project MUST be able to save
the native text verbatim as a `QueryDef` with a non-DTQL type. Native queries MUST be labelled as
such everywhere they are shown, and MUST NOT be silently treated as structured for ACL purposes
(see `server-enforced-access` REQ referencing the opaque-SQL limitation).

#### REQ: backend-adapter-rendering

Every source DataTug claims to support (P2) MUST have a DALgo adapter that can either execute the
`StructuredQuery` AST directly or explicitly reject the unsupported shape (never silently drop a
clause). `sqlite://` and `ingitdb://` MUST stay usable as the reference adapters; `postgres://` and
the remaining SQL/document adapters MUST either render correctly or fail loudly, never partially.

#### REQ: retire-obsolete-dtql-repos

The standalone `dtql`, `dtql-go`, `dtql-backstage` repositories and the `dtql.org` site MUST NOT be
presented as the canonical DTQL implementation while they contain no working code. They MUST either
be retired, redirected to `dal-go/dalgo/dtql`, or brought current before any public claim repeats
their URLs.

## Acceptance Criteria

### AC: dtql-is-source-of-truth (verifies REQ:dtql-ast-consumption)

**Given** a DataTug repository adds new query behaviour
**When** the change touches query parsing, serialization or policy evaluation
**Then** it imports `dal-go/dalgo/dtql` and `dal-go/dalgo/access` rather than adding parallel logic, verified by code review against this REQ.

### AC: dtql-query-type-saves (verifies REQ:dtql-saved-query-type)

**Given** the demo project and a `customer-invoices.dtql.yaml` document
**When** `datatug entity list` and `datatug query run` load the project
**Then** the query loads with `Type: DTQL` and runs through the same policy path as a `SQL` query (see `core-investigation-loop` AC:dtql-query-runs).

### AC: native-query-labelled (verifies REQ:native-query-escape-hatch)

**Given** the demo project's `albums_by_title` SQL-text query
**When** it is listed in the web UI or `datatug query run`
**Then** its type is shown as native SQL, not DTQL, and its result header carries the opaque-SQL policy limitation when one applies.

### AC: sqlite-and-ingitdb-render (verifies REQ:backend-adapter-rendering)

**Given** a `StructuredQuery` built for `Customer` in `chinook.sqlite` and for `support-notes` in the demo project's inGitDB collection
**When** each is executed through its DALgo adapter
**Then** both return rows without dropping any WHERE clause, and a `postgres://` target explicitly returns "not wired" rather than partial results.

### AC: no-canonical-claim-for-stub-repos (verifies REQ:retire-obsolete-dtql-repos)

**Given** `datatug.io` or `dtql.org` copy referencing an installable `dtql` package
**When** the marketing content is reviewed against the current-state audit
**Then** either the copy is corrected to point at `dal-go/dalgo/dtql` or the standalone repos are retired before the claim is repeated.

## Not Doing / Out of Scope

- Connector breadth as a growth lever: rendering to every possible backend is not a race; the
  minimum set that proves the architecture (SQL, inGitDB, Firestore, HTTP) comes first.
- A BI/warehouse tool: DTQL exists to make operational queries portable and policy-checkable, not
  to centralize data into a warehouse.
- Client-direct database access for protected sources: the native-text escape hatch still executes
  through the trusted server, never a direct browser connection.
- A text-to-SQL chat box: DTQL is a structured serialization format, not a natural-language query
  interface; AI-assisted query building is P14's concern, not P1's.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
