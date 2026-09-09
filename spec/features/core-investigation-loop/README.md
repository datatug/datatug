---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Core investigation loop (Phase 1)

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/core-investigation-loop?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/core-investigation-loop?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/core-investigation-loop?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/core-investigation-loop?op=request-change) |

**Status:** Approved
**Date:** 2026-09-09
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Kind:** Capability (implementations in `datatug-cli` and `datatug-apps`)
**Plan:** [2026-09-09-phase-1-core-investigation-loop](../../plans/2026-09-09-phase-1-core-investigation-loop.md)
**Pillars:** P3 semantic layer · P5 related-data navigation · P6 investigation context · P7 executable knowledge library · P10 server-enforced access · P11 surfaces (see [hierarchy](../../product/hierarchy.md))

## Summary

Phase 1 makes DataTug's defining loop work end to end, at a basic but coherent level, in
the web UI over the local agent (`datatug serve`), on one demo project with three sources
(SQLite, inGitDB and a public HTTP endpoint — founder ruling, see A1):

```text
open DataTug → open the demo project → open a table or saved query
→ rows in a dense grid, mapped columns carry a semantic marker
→ select a cell → context panel: meaning · related records (all three sources) · applicable queries (parameters bound)
→ run an applicable query through the server → results, ACL-limited if the principal is restricted
→ add the value to the Investigation Context → open another query: parameter bound from context
→ follow a related record → continue
```

Everything the loop needs already has a partial implementation; this Feature specifies
the missing resolver, endpoints and screens, and forbids everything else. It is a
**vertical slice**, not feature completeness.

### Assumptions made by the lead (not founder rulings — challenge freely)

- **A1 (superseded by a founder ruling, 2026-09-09, verbatim):** *"I want to have in
  MVP: inGitDB, SQLite, HTTP endpoint. Regards HTTP - some public sources like countries
  or currency rates or something that is common and is understandable and can looks good
  in demos."* The demo therefore has three sources. What remains the lead's assumption
  is the concrete choice: **CountriesNow** (`https://countriesnow.space/api/v0.1/countries/currency/q?country={name}`,
  no key; REST Countries was found deprecated on 2026-09-09 — every route redirects to a notice) joined from Chinook's `Customer.Country` through a `Country.Name` semantic
  field, and **Frankfurter** (`https://api.frankfurter.dev/v1/latest?base=USD&symbols={currency}`,
  no key; the old `.app` host now redirects) for the country's currency rate. Both are read-only JSON GET; the demo ships a
  recorded snapshot fixture so the loop still works offline. Reason for these two: common,
  understandable, and they give the loop a third hop (customer → country → currency).
- **A2.** Opaque SQL text is executed under source-level ACL in a read-only session only;
  row and column policies apply to structured (DTQL) queries and to every server-built
  lookup. Reason: DALgo cannot rewrite text it does not parse; the layered-ACL design
  takes the same position.
- **A3.** The internal name `basket` stays; the user-facing label is "Investigation
  Context".

## Problem

Today the web UI cannot reach the agent, the agent cannot start, the demo project link is
broken, and no screen shows a semantic meaning, a related record across sources, an
applicable query, or an investigation context. The concepts exist as data types
(`Entity`, `EntityField.NamePatterns`, `EntityFieldRef` on parameters and columns, the
FK graph, typed `ParameterDef`, `accesspolicies.Run`) with no runtime that connects
them. A knowledgeable observer cannot see why DataTug is different.

## User journeys (the Feature must support all four)

**J1 — First useful result (the null-action path).** The user runs `datatug demo`,
the browser opens the demo project, they click *Album* and see rows. Nobody does anything
else. Good result: rows are visible within two clicks; the server log shows the request.

**J2 — From a value to related knowledge.** On the *Customer* table the user clicks a
`CustomerId` cell. The context panel shows *Customer.ID = 5*, *Related: Invoices (12)
from chinook.sqlite, Support notes (2) from support-notes (inGitDB), Country facts (1)
from CountriesNow (HTTP)*, and *Applicable queries: "Customer invoices" (Customer.ID
bound → 5), "Customer purchases by genre", "Exchange rate for the customer's currency"
(Country.Currency bound → CAD)*.
They click *Customer invoices*; the query page opens with the parameter already filled
and runs through the server. Good result: the results grid shows invoices for customer 5,
and the query page states "Customer.ID bound from selection".

**J3 — Carrying context.** The user presses *Add to context* on `Customer.ID = 5`, then
opens the saved query *Customer purchases by genre*. Its `Customer.ID` parameter is
bound from context, shown as a chip with the origin. They disable the chip; the
parameter empties and the run button asks for a value. Good result: the binding is
visible, reversible, and never silent.

**J4 — Restricted principal.** The user switches the serving principal to `support` (a
policy allows `Customer` rows where `Country == 'Canada'` and hides `Email`). The same
*Customer invoices* run returns only Canadian customers' invoices, the `Email` column is
absent, and the results header shows "Limited by policy customers-support: rows
filtered, 1 column hidden". A direct request from the browser for the hidden column is
refused by the server. Good result: the browser can display limitations but cannot lift
them.

Epilogue branches: (a) the user closes the browser; the CLI agent keeps running and
`datatug serve` prints the URL to reopen the same screen; (b) the user copies the
screen URL to a colleague who runs the same demo — the URL carries project, table and
selection, never credentials.

## Behavior

### Semantic resolution

#### REQ: field-mapping-model

A project MUST be able to declare, per entity field, the physical columns it maps to:
`EntityField.Mappings []PhysicalRef{Source, Collection, Column}` stored in the entity
file. Declared mappings are authoritative. When no mapping is declared, the resolver
MAY derive one from `EntityField.NamePatterns` against the scanned schema, and MUST
label such a mapping as `inferred` in every response that carries it.

#### REQ: semantic-resolution-endpoint

The server MUST expose a read endpoint that, for a source and collection (table), returns
for each column the semantic field it maps to (`entity`, `field`, `provenance:
declared|inferred`) or nothing. The web UI MUST use only this endpoint to mark columns;
it MUST NOT evaluate name patterns itself.

#### REQ: semantic-markers-in-grid

The results grid MUST mark mapped columns (header adornment plus tooltip
`Customer.ID · declared`) and MUST NOT alter values or order.

### Related records

#### REQ: related-lookup-model

Given a semantic value (`entity.field = value`), the server MUST compute the related
lookups from (a) foreign keys in the scanned schema that reference or are referenced by a
column mapped to that field, and (b) semantic mappings of the same field in other
sources. Each lookup MUST carry a label, the target source and collection, and the count
or "count unavailable".

#### REQ: related-lookup-execution

Related records MUST be fetched by the server, which builds a DALgo `StructuredQuery`
(`where target.column == value`, limit N) and executes it through the access-policy path
(REQ:server-acl-all-reads). The browser MUST NOT construct queries for related records.

### Applicable queries and parameter binding

#### REQ: applicable-queries

Given a set of semantic values (from the selection and the active Investigation
Context), the server MUST return the library queries whose semantic parameters
(`ParameterDef.Meta` = `EntityFieldRef`) can all be bound, and separately those that
cannot yet be bound, each with the missing fields. Each applicable entry MUST include
the resolution chain: `selected column CustomerId → maps to Customer.ID (declared) →
query requires Customer.ID`.

#### REQ: parameter-auto-binding

When the query page opens with a selection or context available, parameters whose
`Meta` matches a semantic value MUST be pre-filled and shown with their origin
(`from selection`, `from context`). The user MUST be able to clear or override a bound
value before running. A run MUST send the bound values to the server as ordinary
parameters; the server MUST NOT bind values on its own.

### Investigation Context

#### REQ: context-basket

The web UI MUST keep an Investigation Context per browser session: an ordered set of
semantic values (`entity.field = value`, origin, enabled flag) plus the current
environment. Users MUST be able to add a value from a cell, remove it, and disable it
temporarily without removing it. The context MUST be visible on every project screen
(collapsed by default, count shown).

#### REQ: no-hidden-filters

A query MUST NOT be constrained by context unless a binding is shown on the query page
before the run. The results header MUST list the bindings that were applied.

### Saved queries

#### REQ: dtql-query-type

`QueryDef.Type` MUST accept `DTQL`, stored as `<id>.query.dtql` beside `<id>.query.json` (the filestore's existing `<id>.query.<type>` sidecar rule, landed in datatug-core PR #302)
sidecar; the server MUST execute it through DALgo (`dtql.Deserialize` → policy path).
Existing `SQL` queries keep working.

### Access control

#### REQ: server-acl-all-reads

Every read the web UI can trigger through `datatug serve` (select, saved query run,
related lookup, applicable-query resolution that samples data) MUST pass through
`accesspolicies.Run` with the project's policy set and the serving principal. A request
that references a field outside the principal's allow-list MUST be refused by the
server with a structured limitation, not silently emptied.

#### REQ: principal-selection

`datatug serve` MUST accept `--as <user> --role <r> --group <g>` (the same flags as
`query run`) to fix the serving principal for the session; the web UI MUST display the
principal and MUST NOT be able to change it. (Authenticated multi-user principals are
Phase 3/4.)

#### REQ: opaque-sql-limitation

For `SQL`-type queries executed as text, the server MUST apply source-level policy
(allow/deny for the source) and execute in a read-only session; it MUST label the result
"row/column policies not applied to native SQL" when any row or column rule exists for
that source. (Assumption A2.)

#### REQ: limitation-visible

Every result MUST carry the policy limitations applied (policy name, rows filtered
yes/no, hidden columns), and the web UI MUST show them in the results header.

### Demo project

#### REQ: demo-knowledge-project

`datatug-demo-projects/demo-project-1` MUST gain: entities `Customer`, `Invoice` with
declared mappings and semantic fields `Customer.ID`, `Customer.Email`, `Invoice.ID`;
three parameterized queries (`customer-invoices` requires `Customer.ID`; `invoice-lines`
requires `Invoice.ID`; `customer-purchases-by-genre` requires `Customer.ID`), at least
one as DTQL; an inGitDB `support-notes` collection (5–10 records keyed by `CustomerId`)
with a mapping to `Customer.ID`; an HTTP reference source (REQ:http-reference-source)
mapped to `Country.Name`; and a policy set with principals `admin` (all) and
`support` (Canada rows, `Email` hidden).

#### REQ: http-reference-source

The demo project MUST include at least one public, keyless, read-only HTTP JSON endpoint
as a source (founder ruling 2026-09-09: *"HTTP endpoint ... some public sources like
countries or currency rates"*), declared as a `QueryDef` of type HTTP with a parameter
carrying a semantic tag, so that the server can execute it for related lookups and
applicable queries exactly like the other two sources. Execution MUST be server-side
(`datatug serve`) through the DALgo adapter `dal-go/dalgo2http` (founder ruling 2026-09-09:
*"dalgo2http - yes"*), read-only GET, with response rows shaped as a recordset; the project
MUST carry a recorded snapshot fixture used when the endpoint is unreachable, and the
result header MUST say which of live or snapshot was served. No write, no auth header,
no secret in project files.

### Surfaces and repair

#### REQ: agent-path-contract

The web agent client and `datatug serve` MUST agree on one route table (prefix
`/datatug`, endpoint names as registered in `pkg/server/endpoints/routes.go`), verified
by a contract test that runs the real server.

#### REQ: journey-e2e

A Playwright test MUST walk J1–J4 against the real agent and the demo project with no
network interception, asserting mechanism (server log entries, limitation headers, chip
origins), not only rendered text.

## Architecture

```text
Browser (datatug-apps, Ionic shell + Tabulator grid island)
  │  HTTP  /datatug/*  (loopback, CORS allow-list)             — never a DB connection
  ▼
datatug serve (datatug-cli pkg/server + pkg/api)
  ├─ project store  → filestore (project files; Git-native)
  ├─ resolver       → semantic mappings, related lookups, applicable queries   (new package, e.g. pkg/semantic)
  ├─ executor       → accesspolicies.Run(ctx, db, query, opts)  ← policy set + principal
  │                     └─ access.SecureReadSession → DALgo adapters (dalgo2sqlite, dalgo2ingitdb)
  └─ legacy text SQL → read-only session, source-level policy only (labelled)
```

- The resolver is pure over project data + scanned schema; it is reusable by the TUI and
  by a future `datatug context` verb / MCP resource (Phase 2).
- No new storage. No WebSocket. No auth beyond the fixed principal (loopback).
- Secrets: the demo uses a local SQLite file, a local inGitDB directory and keyless
  public HTTP endpoints; the validator refuses committed credentials (Phase 0).

## API contracts (behaviour, not code)

| Endpoint | Input | Output |
|---|---|---|
| `GET /datatug/semantic/columns` | project, source, collection | `[{column, entity, field, provenance}]` |
| `GET /datatug/semantic/related` | project, entity, field, value, limit | `[{label, source, collection, count|null, lookupId}]` |
| `GET /datatug/semantic/related/rows` | project, lookupId, value, limit | recordset + limitations |
| `POST /datatug/queries/applicable` | project, values `[{entity, field, value, origin}]` | `{applicable:[{queryId, bindings, chain}], notYet:[{queryId, missing}]}` |
| `POST /datatug/exec/run_query` | project, queryId or DTQL document, parameters | recordset + `limitations` + `bindingsApplied` |
| `GET /datatug/agent-info` | — | version, principal, project ids |

Errors are structured (`code`, `message`, `field`); access refusals use `code: ACCESS_DENIED`
with the policy name and never the hidden value.

## Security

- Enforcement point: the server, inside `accesspolicies.Run`. The browser only displays.
- Tests must prove: unauthorized source inaccessible; row restriction survives an
  alternate query on the same table; hidden column never present in any response
  (select, related rows, saved query, DTQL); related lookup and applicable-query
  sampling respect the same policies; a hand-crafted browser request for the hidden
  column is refused; no endpoint accepts a principal from the request.
- Not implemented in Phase 1 (recorded, not claimed): row/column rules on native SQL
  text; multi-user authenticated principals; layered multi-owner policies; Explain Access.

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI / server (`datatug serve`) | Planned | store wiring, resolver, ACL on all reads, DTQL query type, principal flags, demo project | specscore:feature/core-investigation-loop@github.com/datatug/datatug-cli |
| Web (`datatug-apps`) | Planned | agent client contract, semantic markers, context panel, applicable queries, auto-binding, Investigation Context, limitation header, journey e2e | specscore:feature/core-investigation-loop@github.com/datatug/datatug-apps |
| Demo project (`datatug-demo-projects`) | Planned | Customer/Invoice/Country entities, mappings, parameterized queries, inGitDB support-notes, HTTP reference source with snapshot fixture, policy set | — |

## Acceptance Criteria

### AC: first-result-two-clicks (verifies REQ:agent-path-contract)

**Given** `datatug demo` has run and the browser opened the demo project
**When** the user opens the *Album* table
**Then** rows render in the grid and the agent log shows the select request under `/datatug/`.

### AC: mapped-columns-marked (verifies REQ:semantic-resolution-endpoint, REQ:semantic-markers-in-grid)

**Given** the *Customer* table is open
**When** the grid renders
**Then** `CustomerId` and `Email` carry semantic markers reading `Customer.ID · declared` and `Customer.Email · declared`, `FirstName` carries none, and the markers came from the semantic endpoint.

### AC: inferred-mapping-labelled (verifies REQ:field-mapping-model)

**Given** an entity field with a `NamePatterns` rule and no declared mapping
**When** the semantic endpoint resolves a matching column
**Then** the response carries `provenance: inferred` and the grid marker says `inferred`.

### AC: related-across-sources (verifies REQ:related-lookup-model, REQ:related-lookup-execution)

**Given** a `CustomerId` cell with value 5 is selected
**When** the context panel loads
**Then** it lists *Invoices* (source `chinook`, count 7) and *Support notes* (source `support-notes`, count 2), and expanding either fetches rows from the server, not the browser.

### AC: applicable-with-chain (verifies REQ:applicable-queries)

**Given** `Customer.ID = 5` is selected
**When** the panel lists applicable queries
**Then** *Customer invoices* and *Customer purchases by genre* are applicable with the chain `CustomerId → Customer.ID (declared) → requires Customer.ID`, and *Invoice lines* is listed under "not yet applicable — needs Invoice.ID".

### AC: bound-from-selection (verifies REQ:parameter-auto-binding)

**Given** the user opens *Customer invoices* from the panel
**When** the query page renders
**Then** `Customer.ID` shows `5 · from selection`, can be cleared, and a run sends `5` as an ordinary parameter.

### AC: context-carries (verifies REQ:context-basket, REQ:no-hidden-filters)

**Given** `Customer.ID = 5` was added to the Investigation Context
**When** the user opens *Customer purchases by genre* from the queries list
**Then** the parameter is bound `from context`; disabling the chip empties it; the results header lists the bindings applied.

### AC: dtql-query-runs (verifies REQ:dtql-query-type)

**Given** `customer-invoices.query.dtql` exists in the demo project
**When** it is run with `Customer.ID = 5`
**Then** rows return through the policy path and the query is listed with type `DTQL`.

### AC: restricted-rows-and-columns (verifies REQ:server-acl-all-reads, REQ:principal-selection, REQ:limitation-visible)

**Given** `datatug serve --as support`
**When** *Customer invoices* runs for `Customer.ID = 5` (a Brazilian customer) and then for a Canadian customer
**Then** the first returns no rows with limitation `customers-support: rows filtered`, the second returns rows without `Email`, and the header reads `1 column hidden`.

### AC: hidden-column-refused (verifies REQ:server-acl-all-reads)

**Given** the `support` principal
**When** a hand-crafted request selects `Email` from *Customer*
**Then** the server refuses with `ACCESS_DENIED` naming policy `customers-support` and no row data.

### AC: related-respects-policy (verifies REQ:related-lookup-execution, REQ:server-acl-all-reads)

**Given** the `support` principal and a Brazilian customer selected
**When** the related *Invoices* lookup runs
**Then** it returns zero rows with the limitation shown, and the count reads "count unavailable" rather than the true count.

### AC: native-sql-labelled (verifies REQ:opaque-sql-limitation)

**Given** the `support` principal
**When** the SQL-text query *albums_by_title* runs
**Then** it executes read-only and the header states "row/column policies not applied to native SQL".

### AC: demo-project-shape (verifies REQ:demo-knowledge-project)

**Given** the demo project on `main`
**When** `datatug validate` runs
**Then** it passes, and `datatug entity list` shows `Customer`, `Invoice` and `Country` with declared mappings, and the inGitDB collection loads through `ingitdb://`.

### AC: http-source-in-demo (verifies REQ:http-reference-source)

**Given** the demo project and a running agent
**When** the user selects a `Country` cell on the Customer table, or runs the "Exchange rate for the customer's currency" query
**Then** the server executes the HTTP reference source read-only, returns rows shaped like any other source, and the result header states "live" or "snapshot"; with the network disabled the same request succeeds from the snapshot fixture.

### AC: journey-walks (verifies REQ:journey-e2e)

**Given** the agent and web dev server are running
**When** the Playwright journey runs
**Then** J1–J4 pass with mechanism assertions and no request interception.

### AC: agent-path-contract-server (verifies REQ:agent-path-contract)

**Given** `datatug serve --project <demo> --as admin` is running
**When** the route table is listed and `GET /datatug/agent-info` is called
**Then** every route the web client uses exists under `/datatug/` and the response carries version, principal and project ids.

### AC: agent-path-contract-client (verifies REQ:agent-path-contract)

**Given** the real server is running
**When** the web client's contract test calls each agent endpoint it uses
**Then** every call resolves with a success status and no request was intercepted or mocked.

### AC: journey-harness (verifies REQ:journey-e2e)

**Given** the Playwright fixture
**When** the suite starts
**Then** it boots the real agent on a free port and the web dev server against it, J1 passes, and both processes are torn down.

### AC: mapped-columns-server (verifies REQ:semantic-resolution-endpoint, REQ:field-mapping-model)

**Given** the demo project
**When** `GET /datatug/semantic/columns` is called for source `chinook`, collection `Customer`
**Then** `CustomerId` resolves to `Customer.ID` and `Email` to `Customer.Email`, both `declared`, and `FirstName` resolves to nothing.

### AC: restricted-rows-and-columns-server (verifies REQ:server-acl-all-reads, REQ:principal-selection, REQ:limitation-visible)

**Given** `datatug serve --as support`
**When** `POST /datatug/exec/run_query` runs `customer-invoices` for a Brazilian and then a Canadian customer
**Then** the first returns zero rows with `limitations.rowsFiltered = true`, the second returns rows without `Email` and `limitations.hiddenColumns = ["Email"]`, and `agent-info` reports principal `support`.

### AC: related-across-sources-server (verifies REQ:related-lookup-model, REQ:related-lookup-execution)

**Given** the demo project and principal `admin`
**When** `GET /datatug/semantic/related` is called for `Customer.ID = 5` and the *Invoices* lookup rows are fetched
**Then** the lookups list *Invoices* (`chinook`, count 7) and *Support notes* (`support-notes`, count 2), and the rows come from a server-built structured query executed through the policy path.

### AC: applicable-with-chain-server (verifies REQ:applicable-queries)

**Given** the demo project
**When** `POST /datatug/queries/applicable` is called with `Customer.ID = 5`
**Then** `customer-invoices` and `customer-purchases-by-genre` are applicable with a resolution chain, and `invoice-lines` is listed as not yet applicable with `missing: [Invoice.ID]`.

## Not Doing / Out of Scope

- Playbooks, collections beyond folders, capture-from-exploration, knowledge promotion (Phase 2–3).
- HTTP/GraphQL execution, logs, Firestore as a source (Phase 2–3).
- Probabilistic inference, automatic relationship discovery (Phase 5 / experiments).
- Cross-source joins; only lookups.
- Sneat spaces, hosted multi-user principals, layered ACL, Explain Access (Phase 3–4).
- MCP / serve-brokered builder (Phase 4).
- TUI parity for the loop (Phase 2).
- Any UI-kit migration; boards/charts.

## Open Questions

- A1 is ruled (three sources); confirm the lead's concrete endpoint choice (REST Countries + Frankfurter) and assumptions A2–A3 (opaque-SQL limitation, `basket` name).
- Should the serving principal default to `admin` when `--as` is omitted, or should `serve` refuse to start without a principal when a policy set exists? Recommendation: default `admin` only when no policy set exists; otherwise require `--as`.

---
*This document follows the https://specscore.md/feature-specification*
