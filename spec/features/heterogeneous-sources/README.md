---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Heterogeneous operational sources

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/heterogeneous-sources?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/heterogeneous-sources?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/heterogeneous-sources?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/heterogeneous-sources?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

The approved Phase 1 source set is SQLite, inGitDB and public keyless HTTP JSON GET.
The HTTP adapter and CLI path exist; the independent review found browser dispatch
and shared source-resolution gaps. Complete all three through one registry and the
same protected server boundary. GraphQL, authenticated HTTP, Firestore query depth
and logs need separate later plans, not implicit Phase 1 support.

## Problem

An organization's operational data is never in one place: a relational database for transactions, a
document store for flexible records, a Git-backed collection for team-authored notes, an HTTP API
for a third-party service. Tools that assume a single connection (most SQL clients, most BI tools)
force the user to pick one source and lose the others; DataTug's thesis is that the same navigation
and investigation loop should work regardless of where a value lives. Today that is only half true:
`sqlite://` and `ingitdb://` are fully wired (`pkg/dbcopy/url.go`), which is why the Phase 1 demo
uses them, but `postgres://` `Open()` always returns `ErrPostgresNotWired`
(`pkg/dbcopy/url.go:34-37,164-165`), Firestore's query adapter "rejects OR groups and GroupBy"
(`dal-go/dalgo2firestore/query.go:11-52,135-138`), and HTTP execution literally throws "Not
implemented yet" in the web client (`datatug-apps/.../http-executor.ts:20`) even though the
`QueryDefTarget`/`ApiService` model is real. Released CLI binaries also silently drop SQLite support
because they are built `CGO_ENABLED=0` against the `mattn` driver — a gap the audit calls out
explicitly as fixable by switching to the pure-Go `dalgo2sqlite`.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| SQL sources (SQLite, SQL Server, MySQL, Postgres) — scan + query | partial | 0–1 | enabler |
| inGitDB (Git-backed data) as a source | usable (local paths) | 1 | strategic |
| Firestore as a source | partial | 3 | enabler |
| HTTP JSON GET reference queries | adapter/CLI exist; browser acceptance pending | 1 | strategic |
| Authenticated HTTP / GraphQL | future, separately scoped | extension | strategic |
| Logs as a source | planned | 5 | strategic |
| Source registry and connection secrets outside project files | partial (gap) | 0–1 | hygiene |

#### REQ: sql-sources-scan-and-query

DataTug MUST scan and query SQLite, SQL Server, MySQL and Postgres sources through `pkg/schemers/*`
for scan and DALgo adapters for structured execution. Release binaries MUST use the pure-Go
`dalgo2sqlite` adapter so SQLite support does not depend on CGO. Postgres MUST either execute
structured queries or fail with an explicit "not wired" error, never a silent partial result.

#### REQ: ingitdb-source

inGitDB collections addressed by an `ingitdb://` URL MUST be usable as both a scan target and a
query target for local paths, proven end to end by `query run`. This is the second source Phase 1
relies on for the demo journey; see `core-investigation-loop` REQ:demo-knowledge-project for the
concrete `support-notes` collection wired there. Remote inGitDB URLs are out of scope here (see
`git-native-projects` for the remote OpenVaultDB project-store direction).

#### REQ: firestore-source

Firestore MUST be usable as a schema-scan source in every surface that scans SQL sources today.
Firestore as a query-execution target SHOULD support the same WHERE/GroupBy shape SQL sources
support; until it does, the server MUST reject an unsupported clause explicitly rather than drop it
silently, matching the DALgo adapter's own `ErrNotSupported` behaviour.

#### REQ: http-graphql-source

Phase 1 MUST execute configured keyless JSON GET queries through dalgo2http in both
CLI and browser, with common typed recordsets and source identity. Follow the core
transport appendix for target allowlisting, policy capability checks, time/size bounds,
SSRF/redirect handling, explicit snapshot selection and errors. Authenticated HTTP and
GraphQL remain future capabilities requiring a credential authority and their own plan.

#### REQ: logs-source

Logs are a planned source with no model yet. When designed, a log source MUST fit the same
recordset/query shape as other sources rather than introducing a parallel search-only concept, so
the investigation loop (P5, P6) applies unchanged.

#### REQ: secrets-outside-project-files

A source's connection secret (password, API key, token) MUST NOT be storable inside a
Git-committed project file. `datatug validate` MUST reject a `QueryDefTarget` that carries
credentials; credentials belong in a local, non-committed store (for example `~/.datatug`). This
REQ is shared ground with `git-native-projects`' file-layout guarantees; the two Features MUST stay
consistent rather than defining the rule twice.

## Acceptance Criteria

### AC: sqlite-runs-in-release-binary (verifies REQ:sql-sources-scan-and-query)

**Given** a `CGO_ENABLED=0` release build of `datatug`
**When** `datatug query run --db sqlite://chinook.sqlite --from Customer` runs
**Then** rows return without a missing-driver error, because the pure-Go `dalgo2sqlite` adapter is used.

### AC: ingitdb-demo-collection-loads (verifies REQ:ingitdb-source)

**Given** the demo project's `support-notes` inGitDB collection keyed by `CustomerId`
**When** it is opened via `ingitdb://` in `query run` or the web agent
**Then** records load and filter by `CustomerId`, matching `core-investigation-loop` AC:demo-project-shape.

### AC: firestore-scan-works-query-rejects-cleanly (verifies REQ:firestore-source)

**Given** a Firestore-backed source with an OR-group WHERE clause
**When** the source is scanned for schema versus queried for data
**Then** the scan succeeds and the query is explicitly rejected with a clear "unsupported" error, never a silently narrowed result.

### AC: http-query-executes (verifies REQ:http-graphql-source)

**Given** a `QueryDef` of type HTTP pointed at a shipping-status endpoint
**When** it is run from the CLI and from the web UI
**Then** both return a recordset shaped like a table result, and no credential appears in the saved `.query.json` file.

### AC: secrets-guard-rejects-committed-credentials (verifies REQ:secrets-outside-project-files)

**Given** a `QueryDefTarget` with an inline password committed to a project file
**When** `datatug validate` runs
**Then** validation fails naming the offending file and field, consistent with `git-native-projects` AC for the same guard.

## Not Doing / Out of Scope

- Connector breadth as a growth lever: the minimum set (SQL via DALgo, inGitDB, Firestore, HTTP)
  proves the architecture; new connectors are not added to chase parity with any single competitor.
- A data catalog or governance suite: sources are scanned and mapped to be navigated and executed
  against, not certified or governed.
- A BI/warehouse tool: DataTug explores operational data where it lives; it does not centralize
  sources into one store.
- Client-direct database access for protected sources: every source, including HTTP, is reached
  through the trusted server, never a direct browser connection.

## Open Questions

- Authenticated HTTP: which existing server credential authority and descriptor
  policy should a later slice consume? Decide here before implementing secret-bearing
  requests; never put credentials in project files.
The three Phase 1 source types are already approved and are not an either/or choice.

---
*This document follows the https://specscore.md/feature-specification*
