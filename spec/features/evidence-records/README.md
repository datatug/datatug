---
format: https://specscore.md/feature-specification
status: Draft
---
# Feature: Evidence records and result snapshots

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

An execution record is the immutable receipt of one query execution through the DataTug
server: what was asked, under which scope and principal, which bindings and limitations
applied, how many rows came back and a fingerprint of those rows.

A result snapshot is the optional, bounded, policy-re-checked copy of the rows themselves.
Together they are the evidence layer every other incident artefact points at — hypotheses, annotations, checks,
compare and incident metrics all reference a record by id rather than restating what somebody
saw.

## Problem

Production data changes while people investigate it. An engineer pastes "847 stuck orders"
into a chat window at 14:52 and nobody can later tell which query produced it, against which
environment, under whose permissions, or what "stuck" meant that day. Today DataTug has only a
seed of this: `Result.provenance` names the source, query and observed time of the response
currently on screen, and nothing survives the page reload. There is no per-execution audit
record (reconciliation §1, "Auditing — partial"), no stored result, and no fingerprint, so
"is this the same result as an hour ago?" cannot be answered without rerunning and eyeballing.

The founder's requirement (vision §15) is that query-result snapshots become *first-class
immutable evidence* — with an explicit constraint in the same section:

> Historical snapshots must not become a backdoor around current permissions. Retention,
> masking and access semantics need explicit design.

This Feature is that explicit design.

```text
  run a query through the server
            │
            ▼
   ExecutionRecord  (immutable: scope, principal, bindings, limitations, rowCount, fingerprint)
            │
            ├── snapshot? ──► stored rows (opt-in, bounded, masked, retained N days)
            │                         │
            │                         └── read ──► policy re-evaluated for the CURRENT principal
            │
            ├── same query + same bindings over time ──► series ──► incident metrics
            ├── referenced by ──► incident · hypothesis · annotation · board · compare
            └── never ──► a way to see data you are not allowed to see now
```

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Execution record written for a recorded run | planned | MVP (B1) | strategic |
| Result fingerprint on every recorded run, snapshot or not | planned | MVP (B1) | strategic |
| Opt-in bounded result snapshot | planned | MVP (B1) | strategic |
| Access re-checked when a snapshot is read | planned | MVP (B1) | strategic |
| Retention and masking policy per project | planned | MVP (B1) | enabler |
| Records referenced by incidents, annotations, boards, compare | planned | MVP (B1) | strategic |
| Numeric series from repeated runs | planned | MVP (B1) | strategic |
| `datatug execution` CLI and `/datatug/executions/*` API | planned | MVP (B1) | surface |
| Snapshot reused as a shared cache instead of a live run | planned | NEXT | strategic |
| Records of a **write** operation (before/after state, approval) — vision §36 | planned | FUTURE | strategic |
| Record export as part of a portable Resolution Record | planned | NEXT | surface |

#### REQ: execution-record-captured

An execution that goes through the server MUST be recorded when the request asks for it
(`record: true`) or when project policy says record always. Recording MUST NOT depend on the
caller: the web UI, the CLI and an AI agent all produce the same record through the same
executor. A record is immutable once written — no field may be edited or removed, and a
correction is a new record, never an overwrite.

A record carries exactly these fields (field names are the wire contract; types reuse
[`core-investigation-loop/api-contract.md`](../core-investigation-loop/api-contract.md)):

```typescript
type ExecutionRecord = {
  id: string;                       // server-assigned, project-unique, stable
  scope: {project: string; environment: string};
  queryId?: string;                 // folder-qualified library query id
  queryRevision?: string;           // the revision of that query definition at run time
  dtqlHash?: string;                // ad-hoc DTQL instead of a library query
  parameters: {[parameterId: string]: TypedValue};
  bindingsApplied: Binding[];       // exactly the bindings the executor applied
  principal: {id: string; roles: string[]; groups: string[]};
  executedAt: string;               // RFC3339, UTC
  durationMs: number;
  limitations: Limitation[];        // what policy restricted, never what it hid
  provenance: Result['provenance']; // source, collection, mode, observedAt, executionProfile
  rowCount: number;                 // rows actually returned to the caller
  resultFingerprint: string;
  snapshotRef?: string;             // present only when rows were stored
};
```

Exactly one of `queryId` or `dtqlHash` MUST be present, mirroring the `ExecutionRequest` rule
that exactly one of `queryId`/`dtql` is required. `bindingsApplied` and `limitations` are
copied from the `Result` the executor produced, so a record can never claim a binding or an
absence of restriction the run did not actually have. A record MUST NOT contain query text for
an ad-hoc DTQL run (only its hash), MUST NOT contain credentials, and MUST NOT contain result
rows — rows live only in a snapshot, under REQ:result-snapshot-opt-in.

Storage (lead assumption): records are operational records, not knowledge assets. They live in
the **evidence store**, which follows the incident-store routing of
[incidents](../incidents/README.md) REQ:incident-store — the project's own repository under
`executions/<YYYY>/<MM>/<id>.json` by default (Git-native and diffable for a small setup), or
a dedicated repository, the application's repository, or the hosted space-bounded store — and
are written only by the server through the store's own DALgo-backed path. This path does not
depend on the Phase 2 project-write path for knowledge assets. A record created under an
incident carries that incident's `IncidentRef` in its provenance, which is how derived
back-links are answered.

#### REQ: result-snapshot-opt-in

Result rows MUST be stored server-side only when the caller explicitly asks (`snapshot: true`)
**and** project policy allows snapshots for that source. A snapshot is bounded: the server
MUST refuse to store more than the configured byte cap and MUST record the refusal on the
record (`snapshotRef` absent) rather than storing a truncated result that looks complete.

The fingerprint is different: `resultFingerprint` MUST be computed on every recorded run, with
or without a snapshot. It is a hash over the canonical typed rows — SHA-256 over the JSON
encoding of `{columns:[{name,type}], rows: TypedValue[][]}` with object keys in the order
declared in the contract and rows in the order the executor returned them (lead assumption).
Because the hash is order-sensitive, two runs of a query with no deterministic ordering MAY
differ in fingerprint while holding the same rows; a query used as a metric or a check SHOULD
therefore declare an ordering. See Open Questions.

A snapshot MUST store typed values, not rendered text, so a snapshot read produces the same
`TypedValue` shapes as a live run and no second serialization path appears.

#### REQ: snapshot-access-rechecked

Reading a snapshot MUST re-evaluate the **current** principal's policy against the source,
collection and columns recorded on the record, through the same evaluator that authorizes a
live read (`server-enforced-access` REQ:dalgo-access-policies-external). Columns the current
principal may not see MUST be dropped at read time even though they were stored. A principal
with no current access to the recorded source or collection MUST receive a non-disclosing
denial — NOT_FOUND or ACCESS_DENIED per the error envelope, with no row count, no column
names, no fingerprint and no hint that a snapshot exists.

A snapshot is never a backdoor. The rule holds for every reader including the principal who
created the record: access is evaluated now, not at capture time. Row-level policy MUST also
be re-applied where the policy is expressible over the stored rows; where it is not, the
server MUST deny the snapshot read rather than serve rows it cannot re-filter (lead
assumption, and the safer of the two failure modes).

#### REQ: retention-and-masking-policy

A project MUST declare, in project configuration written through the server, a snapshot policy
with at least: a retention duration, a per-snapshot byte cap, and the columns that are masked
by policy. Columns masked by policy MUST never be written to a snapshot at all — masking
happens in the trusted execution path before the rows reach storage, matching vision §37
("Masking must happen in the trusted execution path before prohibited data reaches the client
or AI"), not as a filter applied on the way out.

Defaults, if a project declares nothing (lead assumptions, not founder rulings): retention
**30 days**, per-snapshot cap **2 MiB** — the same bound the transport contract already
applies to an HTTP response body, so one number governs both paths. Retention expiry deletes
the snapshot bytes only; the execution record itself is permanent and keeps its `rowCount`,
`limitations` and `resultFingerprint`, with `snapshotRef` cleared and an `snapshotExpiredAt`
timestamp added. A reader asking for an expired snapshot gets an explicit "expired" response,
never a silent empty recordset.

#### REQ: record-reference-and-provenance

An execution record MUST be addressable by id from an incident, a hypothesis, an annotation, a
board widget and a compare result, and those references MUST be plain ids — no artefact
copies a record's rows, fingerprint or limitations into its own file. A record MUST carry
`Result.provenance` verbatim from the run, so `provenance.executionProfile` tells a later
reader whether the evidence came from a `protected` or an `opaque-privileged` execution
(`core-investigation-loop` REQ:opaque-sql-limitation) without re-deriving it.

Deleting the artefact that referenced a record MUST NOT delete the record. A dangling
reference MUST resolve to an explicit "record not visible to you or removed" response, so
evidence links degrade honestly rather than silently disappearing from a timeline.

#### REQ: series-from-records

Execution records for the same `queryId` with the same `bindingsApplied` values, ordered by
`executedAt`, form a **series**. A numeric projection over a series — `{column, aggregate}`,
where `aggregate` is one of `first`, `sum`, `min`, `max`, `avg`, `count` (lead assumption), or
the implicit `rowCount` — MUST yield metric points `{executedAt, value: TypedValue, recordId}`.

A series MUST be derived, never stored as a second artefact: it is a query over records. A
point MUST carry the `recordId` it came from so a number on an incident metrics table is one
click from the evidence that produced it, and a point whose record the current principal may
not read MUST be omitted from the series with the omission reported, not silently skipped.
This is what the recovery view in vision §15 and §20 needs — `37 → 19 → 3`, each number
traceable.

#### REQ: snapshot-is-evidence-not-an-execution-mode

Reading a stored result snapshot MUST NOT be an execution. `GET /datatug/executions/{id}/snapshot`
returns recorded evidence; it does not produce a new `Result`, does not consume a query target,
and does not set `Result.provenance.mode`. The existing `mode: 'live' | 'snapshot'` and
`snapshotId` on `ExecutionRequest`/`Result.provenance` keep their current meaning — a recorded
**HTTP fixture** under `fixtures/http/<queryId>.json`, addressed as `<queryId>@<recordedAt>` —
and MUST NOT be overloaded to address an execution record's snapshot (lead assumption). Using a
result snapshot in place of a live run as a shared cache (vision §15 benefit 2) is NEXT and
needs its own explicit, labelled request; MVP never substitutes stored rows for a live run.

This is also the boundary against `datatug db snapshot`
(`datatug-cli/spec/ideas/db-snapshot-command.md`, Approved): that command copies whole tables
into an inGitDB tree. A result snapshot is the output of one execution. Both keep their names;
neither is implemented in terms of the other.

#### REQ: api-and-cli

The server MUST expose, under `/datatug`, following the existing transport contract for scope,
staleness and error envelope:

| Method/path | Request | Success envelope |
|---|---|---|
| GET `executions` | Scope + optional `queryId`, `incidentId`, `since`, `until`, `limit` | `{executions: ExecutionRecordBrief[], truncated: boolean}` |
| GET `executions/{id}` | Scope | `ExecutionRecord` |
| GET `executions/{id}/snapshot` | Scope + optional `limit` | `Result`-shaped recordset with `limitations` re-derived now |
| POST `exec/run_query` | `ExecutionRequest` + `record?: boolean, snapshot?: boolean` | `Result` plus `recordId?: string` |

`record` and `snapshot` are additive optional fields on the existing `ExecutionRequest`;
`snapshot: true` without `record: true` is INVALID_REQUEST rather than an implied record.
Listing MUST omit records the current principal may not read, exactly as the query library
omits unauthorized queries, and MUST NOT report how many were omitted.

The CLI follows the established singular-namespace convention (`datatug entity|query|serve`):
`datatug execution list`, `datatug execution show <id> [--snapshot]`, each with `--json`. A
`--snapshot` read obeys REQ:snapshot-access-rechecked identically to the HTTP path, because it
goes through the same server.

## Acceptance Criteria

### AC: record-created-on-request (verifies REQ:execution-record-captured)

**Given** a project whose write path is available and the *Customer invoices* library query
**When** the client posts `exec/run_query` with `record: true` and `Customer.ID = 5`
**Then** the response carries a `recordId`, `GET /datatug/executions/{recordId}` returns a record whose `queryId`, `parameters`, `bindingsApplied`, `principal`, `limitations` and `rowCount` match the run exactly, and a second `PUT`-style attempt to modify any field of that record is rejected.

### AC: fingerprint-without-snapshot (verifies REQ:result-snapshot-opt-in)

**Given** the same query run twice with `record: true` and `snapshot: false`, with no data change between the runs
**When** both records are read
**Then** both carry a non-empty `resultFingerprint`, the two fingerprints are equal, neither record has a `snapshotRef`, and no result rows are stored anywhere for either run.

### AC: snapshot-read-rechecks-policy (verifies REQ:snapshot-access-rechecked)

**Given** a record whose snapshot was captured by a principal who could see the `email` column, and the `support` principal who cannot
**When** `support` reads `GET /datatug/executions/{id}/snapshot`
**Then** the `email` column is absent from the returned recordset and named in `limitations.hiddenColumns`; and when `support`'s access to the whole source is revoked, the same read returns a non-disclosing denial with no row count, column names or fingerprint.

### AC: retention-expires-snapshot-keeps-record (verifies REQ:retention-and-masking-policy)

**Given** a project with the default 30-day retention and a record whose snapshot is 31 days old
**When** retention runs and the record is read
**Then** the snapshot bytes are gone and the snapshot read returns an explicit "expired" response, while the record still returns its `rowCount`, `limitations` and `resultFingerprint`, with `snapshotRef` cleared.

### AC: series-from-repeated-runs (verifies REQ:series-from-records)

**Given** three records for the same query and bindings with `rowCount` 37, 19 and 3
**When** a series is requested with the `rowCount` projection
**Then** three points are returned in `executedAt` order with values 37, 19, 3, each carrying the `recordId` it came from, and a fourth record the current principal may not read is omitted with the omission reported rather than silently dropped.

### AC: record-referenced-from-incident (verifies REQ:record-reference-and-provenance)

**Given** an incident with a hypothesis that cites an execution record by id
**When** the incident is read
**Then** the incident file holds the record id only — no rows, fingerprint or limitations are copied into it — and when the record is unreadable for the current principal the reference renders as an explicit "not visible or removed" entry instead of vanishing from the timeline.

### AC: snapshot-read-is-not-an-execution (verifies REQ:snapshot-is-evidence-not-an-execution-mode)

**Given** a record with a stored snapshot and an HTTP-backed query with a recorded fixture
**When** the snapshot endpoint is called and, separately, `exec/run_query` is called with `mode: 'snapshot'`
**Then** the snapshot endpoint returns recorded evidence without creating a new execution record or touching a query target, and the `mode: 'snapshot'` run still resolves the HTTP fixture `<queryId>@<recordedAt>` exactly as the transport contract specifies, with no cross-addressing between the two.

## Not Doing / Out of Scope

- A second query engine or result store: records and snapshots are written by the same executor
  that already runs every query, behind the same policy evaluator.
- A general-purpose time-series database or metrics backend. A series is a derived read over
  records, not a stored aggregation, and incident **semantics** for metrics belong to
  [incidents](../incidents/README.md).
- Bulk table snapshots. `datatug db snapshot` copies tables into inGitDB; that Idea keeps its
  own name, scope and plan.
- Write/mutation records (who changed Order 72841, before/after state, approval) — vision §36
  is a separate FUTURE track under `server-enforced-access`.
- Charting a series. DataTug renders a numeric table; charts are Dashboardius's, per
  `dashboards` REQ:no-charts-in-datatug-web-ui.
- Re-authorizing the *capturing* principal's rights to a later reader in any form: no delegation,
  no "shared evidence" grant, no signed link that carries access.

## Open Questions

- **Snapshot storage location and Git.** Records are specified above as project assets so they
  travel with the Resolution Record. Snapshot **bytes** are the awkward half: committing
  production rows to a Git project is a data-exfiltration surface, so the lead assumption is a
  server-side store outside the Git working tree, referenced by `snapshotRef`, with an explicit
  per-project opt-in to commit them. Founder decision needed on whether snapshots may ever be
  committed at all.
- **Retention and masking defaults** (30 days, 2 MiB) are engineering assumptions carried over
  from the transport contract's own bounds, not founder rulings — also listed in
  [product-family §9](../../product/product-family.md#9-open-questions).
- **Fingerprint ordering.** Order-sensitive hashing is simple and cheap but makes an unordered
  query's fingerprint unstable. The alternative — sorting canonical rows before hashing —
  costs memory on large results and hides genuine ordering regressions. Decide before the first
  check compares fingerprints across runs.
- **Record-always policy granularity.** Per project, per source, or per query? MVP assumes per
  project with a per-source override; no evidence yet that per-query is needed.

---
*This document follows the https://specscore.md/feature-specification*
