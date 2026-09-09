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
the remaining integration, contracts and acceptance, without expanding the product scope. It is a
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
- **A2 (corrected engineering rule):** protected execution refuses native SQL it cannot constrain. A separately granted privileged opaque-query profile may run native SQL read-only, without a row/column enforcement claim. The support principal never receives that grant. See REQ:opaque-sql-limitation.

- **A3.** The internal name `basket` stays; the user-facing label is "Investigation
  Context".

## Problem

The initial audit found a disconnected prototype. The later independent review found
landed models, server execution, semantic endpoints and browser components, but
incompatible serialized contracts and incomplete protected-user/browser acceptance.
The remaining problem is to make those pieces one coherent investigation. Consult
[the dated review](../../research/2026-09-09-independent-review.md) for evidence;
component existence is not proof that this Feature's acceptance criteria pass.

## User journeys (the Feature must support all four)

**J1 — First useful result (null-action path).** Run `datatug demo`; the browser
opens the registered project. Click the source and Album table, at most two clicks,
and see real rows. Nobody supplies a hidden URL or performs another setup action.
Good result: the grid contains rows and the server records the authorized request.

**J2 — From facts to related knowledge.** Open Customer and select CustomerId 5.
The panel shows Customer.ID, authorized invoice/support-note lookups and applicable
queries with the chosen value and explanation. Open Customer invoices, inspect the
binding, and run. Good result: only that customer's invoices return. Separately select
a Customer row's Country value, retrieve its country facts over live HTTP, then add
an actually returned currency value to context and run its reference-rate query.
Good result: the three-source path uses facts actually selected or retrieved, never
CAD inferred from an unrelated customer ID. Counts come from the pinned demo dataset,
not conflicting handwritten expectations; restricted counts may be unavailable.

**J3 — Carrying context.** Add Customer.ID to context and open another saved query.
The compatible parameter is visibly bound from context. Disable the fact: an untouched
automatic binding clears and required input blocks Run. Select a conflicting customer:
show a choice instead of silently rebinding. Good result: only confirmed bindings
constrain execution. Switch project/environment/principal: no previous scope's facts
or late response enter the new context.

**J4 — Restricted investigator.** Restart the agent as `--as support --role support`.
Browse Customer with an allowed projection: only Canadian rows return, Email is absent,
and the header states the applied limitations. Run customer-invoices for a Brazilian
and a Canadian fixture ID: the former yields no rows and the latter only authorized
invoices. Test support notes under their own policy. Good result: policies cover every
source/table used; there is no assumed Customer→Invoice policy inheritance. Explicit
selection of Email and opaque native SQL are refused server-side without row data.

**Epilogue A — close/reopen.** Closing a tab does not stop the agent. Reopen its
project URL to resume browsing; a new browser session starts with empty context.
Persisted project assets remain, but unsaved context is not promised durable.

**Epilogue B — share/replay.** Copy a navigation link containing only authorized
project/screen identifiers by default. A colleague opens it under their own agent or
authenticated access, supplies/reselects needed values and re-executes. Good result:
no credentials, hidden values or prior authorized results travel in the URL. Optional
value sharing requires explicit preview/consent and recipient reauthorization.

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

The browser MUST maintain typed context scoped by normalized agent URL, project,
environment and server-issued securityContextId, following the
[transport and context contract](api-contract.md#binding-and-context-behavior).
Users add/remove/disable facts; every project screen exposes the context count and
panel. Scope changes, conflicts, missing values, manual overrides and late responses
MUST follow that contract. Session storage alone is not an isolation boundary.

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

`datatug serve` fixes its principal using `--as <user> --role <r> --group <g>`;
a user ID does not imply a role. With a policy set, startup MUST require an explicit
principal. The browser displays but cannot change it. The Phase 1 fixed-principal
profile binds loopback; non-loopback requires a separately authenticated profile.
CORS is not authentication. Unneeded write endpoints MUST fail closed.

#### REQ: opaque-sql-limitation

The protected profile MUST refuse opaque native SQL before execution with
UNSUPPORTED_PROTECTED_EXECUTION. An explicit privileged opaque-query grant MAY allow
read-only native SQL in a separate profile; the result MUST disclose that row/column
enforcement does not apply. Source read permission alone MUST NOT imply that grant.
The browser cannot choose a profile, and the support demo has no opaque grant.
All legacy execution routes MUST obey the same boundary.

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
MUST carry a recorded snapshot fixture offered after a live failure; the user explicitly
selects snapshot mode in a separate request. The result MUST label mode and recorded time. No write, no auth header,
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
  │                     └─ access.SecureReadSession → DALgo adapters (SQLite, inGitDB, HTTP)
  └─ native SQL → refused for protected profile; explicit privileged opaque grant only
```

- The resolver is pure over project data + scanned schema; it is reusable by the TUI and
  by a future `datatug context` verb / MCP resource (Phase 2).
- No new storage. No WebSocket. No auth beyond the fixed principal (loopback).
- Secrets: the demo uses a local SQLite file, a local inGitDB directory and keyless
  public HTTP endpoints; the validator refuses committed credentials (Phase 0).

## API contracts (behaviour, not code)

The normative [transport appendix](api-contract.md) defines exact JSON envelopes,
typed values, request placement, source resolution, binding precedence, error codes,
budgets and migration. Server/client implementations MUST conform to it together.
This supersedes the earlier sketch table; tests must not accept both incompatible
shapes as equally valid.

#### REQ: exact-transport-and-source-contract

Semantic and execution endpoints MUST use the same authorized project source registry,
exact request/response schemas and target-resolution rules in the appendix. Real
consumer/provider tests MUST prove serialization and actual parameter effects.

#### REQ: context-isolation-and-input-resolution

All Phase 1 context and binding behavior MUST follow the appendix, including typed
identity, scope isolation, ambiguity, explicit overrides and required-input blocking.

#### REQ: bounded-policy-safe-http-and-lookups

HTTP and relationship operations MUST obey the appendix's authorization, SSRF,
redirect, size/time/count bounds, cancellation and explicit snapshot rules.

## Security

- Enforcement point: the server, inside `accesspolicies.Run`. The browser only displays.
- Tests must prove: unauthorized source inaccessible; row restriction survives an
  alternate query on the same table; hidden column never present in any response
  (select, related rows, saved query, DTQL); related lookup and applicable-query
  sampling respect the same policies; a hand-crafted browser request for the hidden
  column is refused; no endpoint accepts a principal from the request.
- Not implemented in Phase 1 (recorded, not claimed): execution of native SQL in the protected profile; multi-user authenticated principals; layered multi-owner policies; Explain Access.

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI / server (`datatug serve`) | Partial | store wiring, resolver, ACL on all reads, DTQL query type, principal flags, demo project | [corrective plan](../../plans/2026-09-09-phase-1-core-investigation-loop.md) |
| Web (`datatug-apps`) | Partial | agent client contract, semantic markers, context panel, applicable queries, auto-binding, Investigation Context, limitation header, journey e2e | [corrective plan](../../plans/2026-09-09-phase-1-core-investigation-loop.md) |
| Demo project (`datatug-demo-projects`) | Partial | Customer/Invoice/Country entities, mappings, parameterized queries, inGitDB support-notes, HTTP reference source with snapshot fixture, policy set | — |

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
**Then** it lists *Invoices* (source `chinook`, count from pinned fixture) and *Support notes* (source `support-notes`, count 2), and expanding either fetches rows from the server, not the browser.

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

**Given** `datatug serve --as support --role support` and explicit policies for Customer,
Invoice and support-notes
**When** Customer is browsed with its allowed projection and customer-invoices runs for
a Brazilian then Canadian fixture customer
**Then** Customer exposes only Canadian rows without Email; the first invoice run
returns no rows and the second only authorized invoices; each result names its applied
limitations. No test expects an Invoice result to hide a Customer column it never selected.

### AC: hidden-column-refused (verifies REQ:server-acl-all-reads)

**Given** the `support` principal
**When** a hand-crafted request selects `Email` from *Customer*
**Then** the server refuses with `ACCESS_DENIED` naming policy `customers-support` and no row data.

### AC: related-respects-policy (verifies REQ:related-lookup-execution, REQ:server-acl-all-reads)

**Given** the `support` principal and a Brazilian customer selected
**When** the related *Invoices* lookup runs
**Then** it returns zero rows with the limitation shown, and the count reads "count unavailable" rather than the true count.

### AC: native-sql-labelled (verifies REQ:opaque-sql-limitation)

**Given** support has the protected profile and no opaque-query grant
**When** native SQL is requested through any current or legacy endpoint
**Then** execution is refused before adapter dispatch with UNSUPPORTED_PROTECTED_EXECUTION.
**And** an independently tested explicit privileged grant can execute read-only SQL,
returning `opaque-privileged` provenance and the row/column-enforcement warning.

The existing AC ID is retained for traceability; its earlier label-only behavior is
superseded by fail-closed protected execution.

### AC: demo-project-shape (verifies REQ:demo-knowledge-project)

**Given** the demo project on `main`
**When** `datatug validate` runs
**Then** it passes, and `datatug entity list` shows `Customer`, `Invoice` and `Country` with declared mappings, and the inGitDB collection loads through `ingitdb://`.

### AC: http-source-in-demo (verifies REQ:http-reference-source)

**Given** the demo project and a running agent
**When** the user selects a `Country` cell on the Customer table, or runs the "Exchange rate for the customer's currency" query
**Then** the server executes the HTTP reference source read-only, returns rows shaped like any other source, and the result header states "live" or "snapshot"; with the network disabled the live request fails honestly; only an explicit subsequent snapshot request succeeds, retaining its recorded timestamp and policy checks.

### AC: journey-walks (verifies REQ:journey-e2e)

**Given** the agent and web dev server are running
**When** the Playwright journey runs
**Then** J1–J4 pass with mechanism assertions and no request interception.

### AC: agent-path-contract-server (verifies REQ:agent-path-contract)

**Given** `datatug serve --project <demo> --as admin --role admin` is running
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

**Given** `datatug serve --as support --role support` and explicit policies for each source
**When** real execution requests run customer-invoices for Brazilian and Canadian
fixture customers, and an allowed Customer projection is requested
**Then** unauthorized invoice rows are absent; Customer excludes Email; applied
`limitations[]` correspond to each result's actual source. agent-info reports support.
No request can elevate the serving principal.

### AC: related-across-sources-server (verifies REQ:related-lookup-model, REQ:related-lookup-execution)

**Given** the demo project and principal `admin`
**When** `POST /datatug/semantic/related` is called for `Customer.ID = 5` and the *Invoices* lookup rows are fetched
**Then** the lookups list *Invoices* (`chinook`, count from pinned fixture) and *Support notes* (`support-notes`, count 2), and the rows come from a server-built structured query executed through the policy path.

### AC: applicable-with-chain-server (verifies REQ:applicable-queries)

**Given** the demo project
**When** `POST /datatug/queries/applicable` is called with `Customer.ID = 5`
**Then** `customer-invoices` and `customer-purchases-by-genre` are applicable with a resolution chain, and `invoice-lines` is listed as not yet applicable with `missing` containing the query’s invoice parameter ID and a chain explaining that Invoice.ID is unavailable.

### AC: typed-context-isolation (verifies REQ:context-isolation-and-input-resolution)

**Given** distinct typed values 5 and "5", two conflicting Customer.ID facts, and a
request in flight
**When** binding is attempted and the user changes project/environment/principal
**Then** types stay distinct, conflict blocks an unconfirmed run, old scope facts do
not bind and its late response is discarded; false/zero/null follow declared types.

### AC: real-transport-and-parameter-effect (verifies REQ:exact-transport-and-source-contract)

**Given** pinned server/client/demo revisions and two authorized customer IDs
**When** actual browser-shaped requests exercise each endpoint, empty/error cases and
saved SQL/DTQL parameters
**Then** exact envelopes deserialize, changed parameters change the selected rows,
HTTP targets resolve without a fake database, and only actually applied bindings return.

### AC: safe-http-and-lookups (verifies REQ:bounded-policy-safe-http-and-lookups)

**Given** protected sources, a forged/cross-scope lookup, and disallowed HTTP targets
**When** requests attempt forbidden addresses, redirects, excessive size/time or counts
**Then** they fail at the server without protected data; authorized results remain
bounded and cancellation releases resources. Unavailable exact counts return null.

## Not Doing / Out of Scope

- Playbooks, collections beyond folders, capture-from-exploration, knowledge promotion (Phase 2–3).
- HTTP beyond configured keyless JSON GET reference data; GraphQL, authenticated HTTP, logs and Firestore query depth (separate later plans).
- Probabilistic inference, automatic relationship discovery (separate experiments).
- Cross-source joins; only lookups.
- Hosted multi-user principals, collaboration grants and policy-authoring UI (later plans). Existing layered ACL implementation is reused where required for this protected slice, not redesigned or prohibited.
- MCP / serve-brokered builder (separate extension plan).
- TUI parity for the loop (separate surfaces plan).
- Any UI-kit migration; boards/charts.

## Open Questions

1. **Which Customer ID do the J1 ACs use for the support principal?** The ACs
   `bound-from-selection` and `context-carries` bind `Customer.ID = 5`, but in the
   pinned Chinook data customer 5 is in the Czech Republic, and the support fixture
   policy (`demo-project-1/policies/customers.yaml`) restricts Customer and Invoice
   rows to Canada, so the support principal sees zero rows for it; customer 3
   (Canada) returns 7 invoices (lane S77 live smoke, 2026-09-09, datatug-cli 7fc6efb).
   Either the ACs name a Canadian customer, or the fixture policy admits customer 5,
   or J1 runs as a different principal. Owner: Phase 1 task 16 (fixture policies and
   pinned-data expectations). Not a founder ruling; recorded by the lead session.

None of the remaining items blocks the corrected local Phase 1 slice. Internal `basket` naming is an
engineering detail; the user-facing label is Investigation Context. Concrete HTTP
descriptors are maintained in the demo project and acceptance manifest, with no
claim that a specific provider was founder-selected. Broader hosting, HTTP credentials
and Dashboardius query-binding choices remain in their owning Features.

---
*This document follows the https://specscore.md/feature-specification*
