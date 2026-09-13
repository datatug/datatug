---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Evidence records and result snapshots

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/evidence-records?op=request-change) |
**Status:** Approved
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
  ref: ExecutionRef;                // {storeId, projectId, executionId}; globally unambiguous
  scope: {storeId: string; project: string; environment: string};
  queryId?: string;                 // folder-qualified library query id
  queryRevision?: string;           // the revision of that query definition at run time
  dtqlHash?: string;                // ad-hoc DTQL instead of a library query
  parameters: {[parameterId: string]: TypedValue | TypedValueSet};
  bindingsApplied: Binding[];       // exactly the bindings the executor applied
  principal: {id: string; roles: string[]; groups: string[]};
  policyFingerprint: string;        // canonical current security-policy scope at execution
  executedAt: string;               // RFC3339, UTC
  durationMs: number;
  limitations: Limitation[];        // what policy restricted, never what it hid
  provenance: Result['provenance']; // source, collection, mode, observedAt, executionProfile
  authorizedFields: FieldAccessRef[]; // visible column identities, never values
  rowCount: number;                 // rows actually returned to the caller
  resultFingerprint: string;
  snapshotRef?: string;             // immutable identity; availability is separate state
  incident?: IncidentRef;            // canonical incident association; server-normalized
  grantUses?: GrantRef[];            // server-attested, 1..n when incident grants affected access
  measurements: ScalarMeasurement[];
};
type ExecutionRef = {storeId: string; projectId: string; executionId: string};
type GrantRef = {
  incident: IncidentRef;             // qualifies the incident-store-local grant id
  grantId: string;
  approvalMutationId: string;        // stable identity of the access.approved mutation
};
type FieldAccessRef = {
  storeId: string; project: string; environment: string; source: string;
  collection?: string; column: string; entity?: string; field?: string;
};
type NumericTypedValue =
  | {type: 'number'; value: number}
  | {type: 'integer'; value: string}
  | {type: 'decimal'; value: string};
type MeasurementProjection = {
  id: string; column?: string;
  aggregate: 'first'|'sum'|'min'|'max'|'avg'|'count'|'rowCount';
};
type ScalarMeasurement = {
  projection: MeasurementProjection;
} & (
  | {completeness: 'complete'; value: NumericTypedValue}
  | {completeness: 'unavailable'; reason: 'no-rows'|'non-numeric'|'policy-limited'|'truncated'|'source-refused'}
);
```

Exactly one of `queryId` or `dtqlHash` MUST be present, mirroring the `ExecutionRequest` rule
that exactly one of `queryId`/`dtql` is required. `bindingsApplied` and `limitations` are
copied from the `Result` the executor produced, so a record can never claim a binding or an
absence of restriction the run did not actually have. A record MUST NOT contain query text for
an ad-hoc DTQL run (only its hash), MUST NOT contain credentials, and MUST NOT contain result
rows — rows live only in a snapshot, under REQ:result-snapshot-opt-in.

`authorizedFields` retains one identity for every column returned to the caller: its qualified
physical source/collection/column and, when declared, semantic entity/field mapping. It stores
no values. This manifest is immutable evidence of which fields were present and lets later
annotation writes derive conservative sensitivity even when no result snapshot was retained.
The server normalizes every accepted request to explicit `storeId` in `ref`, `scope`,
`authorizedFields` and the top-level `incident` before writing the record; ambiguity is
rejected. `incident` has exactly this one wire location on an execution record: it MUST NOT be
copied into `Result.provenance`. `Result.execution`, when recording was requested and
succeeded, is only the qualified `ExecutionRef` of this record.

`grantUses` is server-attested authorization audit, not client provenance. It is absent when
ordinary policy alone authorized the result and contains one or more qualified `GrantRef`s
only when those incident-scoped grants changed the authorization decision. The server derives
it from the grants actually applied by the policy evaluator; a client cannot supply, echo or
upgrade it. Each referenced approval mutation MUST resolve through `GrantRef.incident` to an
`access.approved` event for the same `grantId`, the record's top-level `incident` MUST equal
that incident, and the active grant MUST cover the record's exact store/project/source scope.
Correctly qualified references may cross project, incident and evidence stores. An unresolved
reference, qualifier mismatch, or cross-incident/cross-store substitution fails recording
rather than producing an unauditable record. Reading an execution record rechecks current
authorization for each grant reference: when the reader cannot see the incident or grant
scope, `grantUses` is omitted without a count or placeholder, and direct grant resolution
returns `NOT_FOUND`. The immutable reference remains in storage.

`ScalarMeasurement` is computed in the trusted execution path while the authorized result is
available, never reconstructed later from a fingerprint or an expired snapshot. A requested
projection always leaves a retained entry: either one typed numeric `value` with
`completeness:'complete'`, or `completeness:'unavailable'` with a reason. A truncated or
policy-limited result cannot produce a complete aggregate except `rowCount`; refusing the
measurement is preferable to recording a plausible partial number. The projection id and
shape are immutable evidence alongside the query revision, source scope, applied bindings and
limitations already captured by the record.

Storage (lead assumption): records are operational records, not knowledge assets. They live in
the **evidence store**, which follows the incident-store routing of
[incidents](../incidents/README.md) REQ:incident-store — the project's own repository under
`executions/<YYYY>/<MM>/<id>.json` by default (Git-native and diffable for a small setup), or
a dedicated repository, the application's repository, or the hosted space-bounded store — and
are written only by the server through the store's own DALgo-backed path. This path does not
depend on the Phase 2 project-write path for knowledge assets. A record created under an
incident carries that server-normalized `IncidentRef` in its top-level `incident` field, which
is how derived back-links are answered.

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
the snapshot bytes only. The immutable execution record and its `snapshotRef` are never edited.
A separate snapshot-state record keyed by that ref carries
`{availability:'available'|'expired'|'deleted', changedAt, reason?}` and may advance only from
`available` to a terminal state. A reader asking for an expired snapshot gets the explicit
state and timestamp, never a silent empty recordset; the execution record permanently keeps
its `rowCount`, `limitations`, `resultFingerprint` and scalar measurements.

#### REQ: record-reference-and-provenance

An execution record MUST be addressed from an incident, hypothesis, annotation, board widget
or compare result by `ExecutionRef = {storeId, projectId, executionId}`. A bare execution id
is accepted only inside an already-qualified store and project route; it is never persisted as
a cross-artefact reference. No artefact copies a record's rows, fingerprint or limitations
into its own file. A record MUST carry
`Result.provenance` verbatim from the run, so `provenance.executionProfile` tells a later
reader whether the evidence came from a `protected` or an `opaque-privileged` execution
(`core-investigation-loop` REQ:opaque-sql-limitation) without re-deriving it.

Deleting the artefact that referenced a record MUST NOT delete the record. A dangling
reference MUST resolve to an explicit "record not visible to you or removed" response, so
evidence links degrade honestly rather than silently disappearing from a timeline.

#### REQ: series-from-records

Execution records form one **series** only when all of these identities match: evidence-store
id, source store/project/environment, security-context policy fingerprint, query id and query
revision (or ad-hoc DTQL hash), applied bindings, and scalar projection id/shape. Different
environments, revisions or `policyFingerprint` values MUST NOT be joined into one line. The retained
`ScalarMeasurement` yields `{executedAt, value, execution: ExecutionRef}` without reopening a
snapshot. An unavailable measurement yields an explicit unavailable point with its reason;
series APIs report it and never omit it as though no run occurred.

A series MUST be derived, never stored as a second artefact: it is a query over records. A
point MUST carry the qualified execution ref it came from so a number on an incident metrics
table is one click from the evidence that produced it, and a point whose record the current
principal may not read MUST be omitted from the series with the omission reported, not silently
skipped.
This is what the recovery view in vision §15 and §20 needs — `37 → 19 → 3`, each number
traceable.

A duration criterion such as `value == 0 for 30 minutes` is satisfied only when its metric
declares an observation cadence and maximum tolerated gap, the window begins with a complete
qualifying observation at or before its start, every subsequent observation through the end is
complete and qualifying, and no gap exceeds the tolerance. One zero point, an unavailable
point, a policy-hidden point or a gap cannot prove the duration; the
criterion remains `unknown` with the missing interval reported.

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
| GET `executions` | Scope + optional `queryId`, `incident: IncidentRef`, `since`, `until`, `limit` | `{executions: ExecutionRecordBrief[], truncated: boolean}` |
| GET `executions/{id}` | Scope | `ExecutionRecord` |
| GET `executions/{id}/snapshot` | Scope + optional `limit` | `Result`-shaped recordset with `limitations` re-derived now |
| POST `exec/run_query` | `ExecutionRequest` + `record?: boolean, snapshot?: boolean, measurementProjections?: MeasurementProjection[]` | `Result`, whose top-level `execution?: ExecutionRef` is present only when recording succeeds |

`record` and `snapshot` are additive optional fields on the existing `ExecutionRequest`;
`snapshot: true` without `record: true` is INVALID_REQUEST rather than an implied record.
`measurementProjections` names the immutable scalar projections to retain in the record and
requires `record: true`; each entry uses the `MeasurementProjection` shape and its id
must be unique within the request.
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
**Then** the response carries a qualified `execution`, `GET /datatug/executions/{executionId}` under its qualified Scope returns a record whose `queryId`, `parameters`, `bindingsApplied`, `principal`, `limitations` and `rowCount` match the run exactly, and a second `PUT`-style attempt to modify any field of that record is rejected.

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
**Then** the snapshot bytes are gone and the snapshot read returns an explicit "expired" response from the separate snapshot-state record, while the immutable execution record still returns its original `snapshotRef`, `rowCount`, `limitations`, `resultFingerprint` and scalar measurements unchanged.

### AC: series-from-repeated-runs (verifies REQ:series-from-records)

**Given** three records for the same query and bindings with `rowCount` 37, 19 and 3
**When** a series is requested with the `rowCount` projection
**Then** three points are returned in `executedAt` order with values 37, 19, 3, each carrying the qualified `ExecutionRef` it came from, and a fourth record the current principal may not read is omitted with the omission reported rather than silently dropped.

### AC: record-referenced-from-incident (verifies REQ:record-reference-and-provenance)

**Given** an incident with a hypothesis that cites a qualified `ExecutionRef`
**When** the incident is read
**Then** the incident file holds the qualified reference only — no rows, fingerprint or limitations are copied into it — and when the record is unreadable for the current principal the reference renders as an explicit "not visible or removed" entry instead of vanishing from the timeline.

### AC: grant-use-is-server-attested (verifies REQ:execution-record-captured, REQ:record-reference-and-provenance)

**Given** `support` can read an incident but ordinary project policy denies the incident's
`Invoice` source, and a server-approved incident grant widens row access without unmasking
fields
**When** `support` records a query carrying that incident and then attempts another request
whose unknown client-supplied `grantUses` claims a grant from another store
**Then** the record's top-level `incident` is the server-normalized `IncidentRef`, its
`grantUses` contains only the server-applied grant qualified by that incident and its
`access.approved` mutation id, `Result.execution` contains only the qualified `ExecutionRef`,
the forged request is rejected before execution, and an otherwise identical run authorized by
ordinary policy has no `grantUses` field.

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

- ~~**Snapshot storage location and Git.**~~ Records are specified above as operational records
  in the evidence store, which follows the incident-store routing and so may be the project's
  own repository, a dedicated or application repository, or the hosted store. Snapshot **bytes**
  were the awkward half: committing production rows to a Git project is a data-exfiltration
  surface. **Ruled by the founder, 2026-09-11 (verbatim): *"Not in MVP"*** — in MVP, snapshot
  bytes never enter a Git-tracked path; no per-project opt-in to commit them exists this phase.
- **Snapshot byte storage backend.** The ruling above says only that snapshot bytes never enter
  Git in MVP; it does not specify where they actually live (a dedicated server-side store, a
  blob store, or something else), referenced by `snapshotRef`. That design is not settled here
  and is left as an explicit open question for the Incidentius rework rather than invented in
  this evidence-records draft.
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
