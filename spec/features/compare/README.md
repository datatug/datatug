---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Compare

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=request-change) |
**Status:** Approved
**Source Ideas:** —

## Summary

Compare executes one library query under two scopes and diffs the two recordsets by key.
Executing checks through Compare is explicitly deferred until the Checks Feature and
Incidentius Task 5 have shipped their executable contract; Task 8 does not expose `checkId`.

It is the execution primitive behind the most common troubleshooting move there is — put
something that works next to something that does not and look at the difference — whether the
two sides are two environments, two cohorts of records, or the same query at two points in time.

## Problem

The founder names this as a central troubleshooting question and then generalises it
(vision §13):

> **Why does this work in UAT but fail in PROD?**

> **Compare something that works with something that does not.**

The examples span bad vs good records, PROD vs UAT, before vs after a deployment, FedEx vs UPS,
affected vs unaffected customers, and the current incident vs a previous occurrence. Today none
of that exists as a primitive. `datatug-core/pkg/comparator` compares **schemas** across
environments — tables, columns, differences in structure. It is a useful sibling and it answers
a different question: it cannot tell you that 93% of affected orders ship with one provider
while only 31% of healthy ones do.

Without the primitive, the comparison happens in a spreadsheet, off the platform, with no
provenance and no policy — which is exactly where a production data leak comes from.

```text
          one library query
                 │
        ┌────────┴────────┐
        ▼                 ▼
      LEFT scope        RIGHT scope
   env=prod          env=uat
   facts=affected    facts=control          each side executed SEPARATELY,
   execution=ops/billing/exec-7f3   execution=ops/billing/exec-9ab
                                              under ITS OWN policy, as its own
        │                 │                 execution record
        └────────┬────────┘
                 ▼
          keyed recordset diff
      added · removed · changed(per column) · summary
      optional categorical distribution per side
                 │
                 └── limitations reported PER SIDE; a policy-filtered side
                     is labelled, never presented as a data difference
```

## Behavior

The public-fixture browser session subtype below is separate from this durable, citable Compare
contract and does not narrow or satisfy its existing requirements.

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Two-scope execution: env A vs env B | planned | MVP (B1) | strategic |
| Two-scope execution: fact set A vs fact set B (cohorts with roles) | planned | MVP (B1) | strategic |
| Two-scope execution: execution record A vs execution record B | planned | MVP (B1) | strategic |
| Keyed recordset diff with per-column differences | planned | MVP (B1) | strategic |
| Cohort distribution comparison for a categorical column | planned | MVP (B1) | strategic |
| Per-side limitations; policy-limited diffs labelled | planned | MVP (B1) | strategic |
| Explicit, ephemeral comparison of two admitted public browser fixtures | planned | bounded browser session | surface |
| `POST /datatug/compare` and `datatug compare` | planned | MVP (B1) | surface |
| Compare result referenced as incident evidence | planned | MVP (B1) | strategic |
| Ordered-stream keyed diff backed by DALgo `recordops` | planned | MVP (B1) | strategic |
| Optional evictable local comparison-row cache for paged UI | planned | MVP (B1), Task 15 | surface |
| Compare UI: one "Compare affected vs control" action on the incident page, rendering `CompareResult` as plain tables (summary; changed/added/removed rows; per-side limitations; distribution) | planned | MVP (B1) | surface |
| Compare UI: side-by-side grids, diff highlighting | planned | NEXT | surface |
| Comparison boards in Dashboardius | planned | NEXT | surface |
| Comparing configuration, deployments and infrastructure (vision §13) | planned | FUTURE | strategic |

#### REQ: two-scope-execution

Compare MUST take one required library query id and two side specifications. Check execution is
not part of this Task 8 contract. A `scope` or `facts` side MUST execute the query live, separately under that
side's own scope; a `record` side MUST NOT re-execute anything and instead reuses the named
past execution record. The three side kinds are:

```typescript
type SideSpec =
  | {kind: 'scope'; storeId: string; project: string; environment: string; parameters?: {[id: string]: TypedValue}}
  | {kind: 'facts'; storeId: string; project: string; environment: string; cohortRole: 'affected' | 'control'}
  | {kind: 'record'; execution: ExecutionRef};                                // a past execution
```

The union is strict: exactly one discriminator shape is accepted, fields belonging to another
shape are rejected, and unknown fields (including `checkId`) are `INVALID_REQUEST`.

- **environment A vs B** — the same query, two environments (`{kind:'scope'}` on each side).
- **fact set A vs B** — the same query and environment, two incident-bound sets: the affected
  cohort and the healthy control. A facts side never accepts client-provided `Fact[]`; the
  server reads the named incident under the current principal and selects its current-policy-
  visible, enabled, canonical facts for `cohortRole` whose scope exactly matches the side's
  `storeId`, `project` and `environment`. Facts from another project or environment are excluded
  before same-field and bounded-set validation. Cohort roles are the
  [investigation-context](../investigation-context/README.md) extension; compare consumes them
  and does not reimplement them (product-family §3, the firm boundary). If either side is
  `facts`, `CompareRequest.incident` is required and is the sole fact authority.
- **execution record A vs B** — two past runs, i.e. snapshots or time. A `record` side MUST first
  have a readable row snapshot. A terminal snapshot state (`expired` or `deleted`) MUST fail the
  whole compare with `SNAPSHOT_EXPIRED`; a record for which no retained readable row snapshot
  exists returns the same code. Once that lifecycle check succeeds, the record MUST also have
  `resultComplete:true`. `resultComplete:false` or an absent legacy value cannot prove a complete
  input and MUST fail closed with `SOURCE_UNAVAILABLE`; a snapshot read that reports
  `truncated:true` fails with the same code. Neither case becomes an empty side or a fingerprint-only comparison
  (AC:snapshot-vs-live-shows-recovery). The server MUST load each qualified execution record,
  require its non-empty `queryId` to equal `CompareRequest.queryId`, and reject a mismatch before
  diffing; records from different revisions of that same query remain comparable.

A live (`scope` or `facts`) side MUST produce its own new execution record (`evidence-records`
REQ:execution-record-captured); a `record` side reuses its existing execution record without
re-execution. Either way the two halves of a comparison are independently citable evidence. A
live side is executed as one ordinary authorized run; there is no combined query, no cross-source
join and no second executor — this Feature adds a caller and a differ, not an engine.
Where policy and retention allow, live executions request and persist row snapshots so the
comparison is reproducible. If a side cannot be reproduced, its receipt and limitations say so;
the response never implies that a fingerprint alone can recreate rows.

A `facts` side resolves its incident cohort and binds through the transport appendix's bounded
`TypedValueSet`: exactly one same-field equality cohort of 1..500 unique, non-null, same-typed values. The query parameter
must declare set support and the source must execute one native set/array/`IN` binding. Any
other cohort shape is `INVALID_REQUEST`; a source without safe set binding is
`SOURCE_UNAVAILABLE`. Compare never loops one execution per fact or silently truncates a
cohort.

#### REQ: public-fixture-browser-session

The browser MAY expose a separate, explicitly selected comparison subtype for exactly two
currently admitted public SQLite fixture targets: a pinned baseline and a separately published
synthetic candidate. This subtype compares one shared browser query draft; it does not extend the
server `SideSpec` union, `CompareRequest`, or `CompareResult`, and it MUST NOT be represented as
durable, independently citable, replayable, incident-linked, or shareable evidence. Agent-backed
Compare continues to use the three `SideSpec` kinds above and all existing persistence, policy,
receipt, and incident requirements in this Feature.

Each explicit side Run freezes the same authored source/revision, relationship binding, typed
CustomerId binding and project/access context. It starts an independent worker with its own
generation, abort/cancel state, result, timestamp, target/schema/fixture pin, row count,
limitations and actual worker receipt. A side observation is eligible only after the worker has
exhausted its entire bounded SQL statement successfully within row and byte budgets; only then
may it set `resultComplete:true`. `LIMIT 1–100` means complete for that query result, never the
whole database. A failed, cancelled, capped, partial, expired, denied, missing, changed or stale
side is ineligible. One side's failure MUST preserve the other side's last committed result.
Retries use the frozen pair; a new generation requires both sides to run again. Old results may
remain visible as Previous run, but a source, binding, target, project/principal context or
metadata change immediately invalidates comparison eligibility. Late responses cannot commit
into a newer generation.

Side observations and the comparison result exist only in page memory and expire on refresh or
close. Correlation IDs are local bookkeeping, not authorization. Neither object contains an
`ExecutionRef`, incident identity, durable `ComparisonRef`, or shareable evidence URL. The UI
MUST label the result `Session comparison · expires when this page closes or refreshes`, show both
actual targets, fixture pins, times, row counts, receipts and limitations, and state that this
result does not satisfy `REQ:compare-produces-evidence` or the citable/replayable ACs.

Immediately before Compare, the application re-reads both currently authorized public catalog,
schema and relationship entries and verifies both fixture bytes/pins, the frozen draft and
bindings, current project/access context, and two live eligible generations. It then explicitly
compares by unique non-null `InvoiceId` through the shared Go Core typed diff kernel in a published,
verified WASM adapter. The browser does not implement a JavaScript diff. Any missing/expired/
changed/denied side, null or duplicate key, unsupported cell or comparator, incompatible shared
type, conflicting known unit, or unequal policy visibility blocks comparison. The first browser
transport admits safe integer keys and raw text only; it makes no unsafe-int64 or exact-DECIMAL
claim. `InvoiceDate` stays raw stored SQLite text under the versioned joined representation rule;
there is no temporal normalization. Unknown units permit only an explicitly labelled raw typed
comparison, with no semantic-equivalence or tolerance claim. The public fixture profile has no
row filtering or hidden columns; any policy-limited target requires a separate reviewed contract.

Show changed raw before/after values and baseline-only/candidate-only keys in these bounded query
results. They are not claims of database-row deletion or cause. Core's `Diff.truncated` may report
a cap on displayed differences after both complete inputs; it is a distinct limitation and never
makes a partial side eligible. The display remains session-only even if both underlying worker
receipts are visible.

#### REQ: keyed-recordset-diff

The two recordsets MUST be diffed by key, producing `added`, `removed` and `changed` rows plus
an aggregate summary. If either side is a historical `record`, the caller MUST provide the
explicit `key`; no semantic, obvious-column, positional, whole-row or fingerprint-only fallback
is allowed. For two live sides, caller-supplied key columns take precedence and otherwise the
key fields of the mapped entity may be used. If the resolved key is absent or not unique on
both sides, the request MUST fail with an explicit error naming the ambiguity.

```typescript
type CompareSideReceipt = {
  execution: ExecutionRef;
  executedAt: string;
  rowCount: number;
  limitations: Limitation[];
  reproducible: boolean;
};

type CompareResult = {
  left: CompareSideReceipt;
  right: CompareSideReceipt;
  columns: Column[];                                      // visible shared intersection, stable order
  key: string[];
  added:   {key: TypedValue[]; row: TypedValue[]}[];     // present right, absent left
  removed: {key: TypedValue[]; row: TypedValue[]}[];     // present left, absent right
  changed: {key: TypedValue[]; columns: {column: string; left: TypedValue; right: TypedValue}[]}[];
  summary: {
    added: number; removed: number; changed: number; unchanged: number;
    columnsOnlyOnOneSide: {column: string; side: 'left' | 'right'}[];
  };
  distribution?: Distribution;                            // REQ:cohort-distribution-comparison
  policyLimited: boolean;                                 // REQ:compare-is-policy-bound
  truncated: boolean;
};
```

Comparison of values MUST be typed, using the same `TypedValue` equality the binding contract
fixes: no string/number coercion, `null` equals only `null`, and a `decimal` compares by its
canonical decimal form, not as a float. `columns` is the explicit metadata for the stable,
visible intersection and every `added`/`removed` row is aligned to it. Only those columns may
appear in `changed`; a column present on one side only MUST be reported in
`summary.columnsOnlyOnOneSide` with its side, not as every row having changed. Key, row, column
and distribution value ordering MUST be deterministic. Results MUST be bounded by the same default and
maximum result limits the transport contract sets, with `truncated` set when the cap is hit.

The first-class Core comparison path MUST accept two streams already ordered by the resolved
composite typed key. It MUST delegate key alignment to DALgo `recordops.DiffFunc` — the shared
lagging-cursor/K-way merge primitive — rather than maintaining a second map-and-union merge loop.
Each input cursor advances only when it is at the smallest current key. The stream path MUST
retain O(1) alignment state (one current row per side), reject a duplicate key or a stream that is
not strictly ascending, and preserve deterministic composite-key and typed-value semantics. The
existing materialized `Recordset` comparer remains a compatibility adapter: it MAY copy and sort
its already-materialized rows before feeding the same ordered-stream path, but it MUST NOT carry
an independent diff algorithm.

The portable typed-key order is lexicographic by composite-key component. Numbers, integers and
decimals order by numeric value; exact decimal spellings with equal numeric value use their
canonical transport text as a tie-break. `false` precedes `true`; dates and datetimes order
chronologically, with exact datetime transport text breaking equal-instant ties; and strings use
UTF-8 binary lexical order, not a locale collation. Key column types MUST match on both sides and
null key components are invalid. Core MUST expose this comparator and a canonical sort-key encoder
so ordered providers can implement or verify the contract. The encoded `sort_key` MUST have the
same order under binary collation; SQLite cache indexes therefore use binary collation for it.

Core MAY accept a synchronous, error-returning observer for each aligned record. The observer
receives a zero-based deterministic ordinal, canonical sort key, typed display key, state
(`matched`, `added`, `removed`, or `changed`), and post-policy shared-visible data only. A matched,
added, or removed record carries one full row; a changed record carries the full left baseline row
once plus right-side column deltas. Returning an observer error MUST stop both input streams and
fail the comparison so a transactional consumer can roll back. The observer does not make Core a
storage owner and does not relax the bounded `CompareResult` response.

The complete input bound is separate from that response bound. MVP compare may materialize at most
10,000 authorized rows per side before diffing. A live side MUST detect a row beyond that bound,
persist its execution receipt with `resultComplete:false`, and fail the comparison with
`SOURCE_UNAVAILABLE`; it MUST NOT compare or snapshot the first 10,000 rows as if they were complete.
A record side whose complete snapshot exceeds 10,000 rows fails with the same code. Within that bound,
the comparer consumes the complete side, then applies the default 100 and maximum 500 only to emitted
diff rows; the input bound MUST NOT truncate a complete 847-row replay to the transport maximum.

#### REQ: cohort-distribution-comparison

For a chosen categorical column, compare MUST report the counts and percentages of each distinct
value on each side, and MUST surface over-representation — but MUST NOT state a conclusion.

```typescript
type Distribution = {
  column: string;
  values: {
    value: TypedValue;
    left:  {count: number; pct: number};
    right: {count: number; pct: number};
    ratio: number | null;    // left.pct / right.pct; null when right.pct is 0
  }[];
  truncated: boolean;        // distinct values capped
};
```

This is vision §14's example — affected orders 93% FedEx against a healthy control at 31% — and
the rule the founder states with it is binding on the output:

> Incidentius should not conclude that FedEx is broken, but should surface the overrepresentation
> and immediately offer FedEx-specific diagnostics.

So the response carries counts, percentages and a ratio, and the surface presenting it MUST
present the finding as an observation with its two execution records attached, never as a cause.
Offering the relevant diagnostics is the [checks](../checks/README.md) REQ:checks-listed-by-purpose-and-scope
path, driven by the entity the over-represented value belongs to. Distinct values MUST be capped
(lead assumption: 50, matching the related-discovery cap) with `truncated` set — a distribution
over a high-cardinality column is a mistake, not a feature.
Values use deterministic typed ordering, and counts, percentages and ratios are computed from
the full authorized side row counts before the 50-value output cap is applied.

#### REQ: compare-is-policy-bound

Neither side may disclose more than its own authorized execution would. Concretely:

- Each side MUST be executed under the **current** principal against that side's scope. Compare
  never runs one side with the other side's authority, and there is no "compare principal".
- A column hidden on either side MUST NOT appear anywhere in the diff — not in `changed`, not in
  `added`/`removed` row payloads, not in `distribution`. The comparable column set is the
  intersection of the two sides' visible columns (lead assumption), and the dropped columns MUST
  be named in that side's `limitations` exactly as the transport contract already allows.
- `limitations` MUST be reported **per side**, never merged, so a reader can tell which half was
  restricted.
- When either side reports `rowsFiltered: true`, `policyLimited` MUST be true and the surface
  MUST label the diff: rows missing because a policy filtered them are indistinguishable from
  rows missing because the data differs, and presenting them as `removed` without that label is
  a false finding. A `record` side whose snapshot the current principal may not read MUST fail
  the whole compare with a non-disclosing denial (`evidence-records` REQ:snapshot-access-rechecked),
  never a one-sided result.

Masking happens in the trusted execution path before the rows reach the differ (vision §37), so
the differ never holds a value the caller could not have retrieved directly.

#### REQ: compare-produces-evidence

A compare MUST be citable. The response MUST carry both qualified `ExecutionRef`s, and a compare invoked under
an incident MUST enter the event stream as a `compare.run` event with
`assertion.kind = deterministic-result` referencing those two records and the key used — the
comparison is a deterministic computation over two recorded executions, not a claim
([incidents](../incidents/README.md) REQ:event-stream-append-only). Re-running the same compare later produces two new records when both sides are live; record sides reuse their named records. The earlier comparison
stays valid evidence of what was true then.

For a compare that executes either side under `CompareRequest.incident`, the server normalizes
that ref onto each new execution record's top-level `incident` field. It MUST NOT place the
incident inside `Result.provenance`; each result and the compare response expose only the
qualified `ExecutionRef`s of the recorded sides.

An incident compare requires a caller-supplied, non-empty `mutationId`. The server appends
`compare.run` only after both sides, the diff and every required execution/snapshot persistence
step have succeeded. If a later step fails, any live-side execution receipt already persisted
remains valid evidence, the endpoint returns an incomplete typed error carrying those receipts,
and no `compare.run` is appended. Retrying an already completed `mutationId` does not re-execute
or manufacture a cached `CompareResult`: because the result itself is not stored, the server
may return an already-applied conflict with the prior qualified `ComparisonRef`.

MVP does not persist a `CompareResult` as its own project asset (lead assumption). An optional
evictable local row cache described by REQ:local-comparison-cache is operational acceleration,
not a portable project asset or evidence authority. The two
execution records plus the request parameters reproduce it only while both receipts retain
readable row snapshots; otherwise the receipt states honestly that the result is not
reproducible. Persisting the diff would create a third copy of data already governed by retention
and masking. An incident references the pair.

#### REQ: local-comparison-cache

Task 15 MAY ask the local DataTug server/CLI to cache every aligned comparison record so the web
UI can page through matching as well as non-matching rows after the bounded `CompareResult` has
been returned. This is a consumer of Core's ordered-stream observer, not Core storage code and not
part of Task 9's incident-page model.

The cache MUST be local, optional, evictable, non-authoritative and schema-versioned. Its identity
MUST bind the store, project, environment, principal, policy revision, query id and canonical
bindings, plus both qualified snapshot/execution refs. A cache hit MUST re-check current
authorization and MUST NOT substitute for execution-record or snapshot authorization. Only rows
already reduced to the shared post-policy visible columns may reach the cache.

The MVP SQLite shape is a fixed envelope, not dynamic source columns:

```sql
comparison_rows(comparison_id, sort_key, state, payload_json)
INDEX (comparison_id, sort_key)
INDEX (comparison_id, state, sort_key)
```

`payload_json` MUST be canonical/minified JSON stored as SQLite `TEXT` (not SQLite JSONB), with no
compression or codec negotiation. A matched record stores `{key,row}` once; added and removed
store their one full row once; changed stores the left baseline row plus right-side deltas. The UI
pages through an opaque comparison-cache ref and never treats a cache row as incident evidence.
Extracting explicitly selected JSON fields into indexable columns is a separate future Idea; MVP
does not add arbitrary JSON-field indexes, dynamic source columns or columnar spill.

#### REQ: schema-compare-remains-separate

`datatug-core/pkg/comparator` (schema and structural differences across environments) MUST
remain as it is. This Feature MUST NOT absorb, rename or reimplement it, and MUST NOT be
implemented on top of it: one compares structure, the other compares records, and the two answer
different questions on the same screen. "Compare configuration, deployments, reference data and
infrastructure" (vision §13) is the FUTURE convergence point where both feed one comparison view;
nothing in MVP depends on that.

#### REQ: api-and-cli

The server MUST expose `POST /datatug/compare`, following the existing scope, staleness and
error envelope rules:

```typescript
type CompareRequest = {
  securityContextId: string;
  queryId: string;
  left: SideSpec; right: SideSpec;
  incident?: IncidentRef;
  mutationId?: string;                     // required iff incident is present
  key?: string[];
  distributionColumn?: string;
  limit?: number;
};

type CompareErrorResponse = {
  error: ErrorBody;
  left?:  CompareSideReceipt;              // present only after this side persisted
  right?: CompareSideReceipt;
  comparison?: ComparisonRef;              // prior completed mutation, when known
};
```

It is a POST because sides carry semantic values, which must never reach a URL, browser history
or an access log — the same reason `semantic/related` is a POST. An absent/blank `queryId`, a
`checkId`, an incident without a non-empty `mutationId`, or a non-empty `mutationId` without an
incident is `INVALID_REQUEST`. An omitted and explicitly empty `mutationId` are the same semantic
absence when there is no incident. Check-backed comparisons will amend this union only after Task
5 freezes the check executable contract.

The CLI takes one executable and two sides, with the side syntax naming its kind:

```text
datatug compare --query stuck-orders --left env=prod       --right env=uat
datatug compare --query stuck-orders --left facts=affected --right facts=control
datatug compare --query stuck-orders --left execution=ops/billing/exec-7f3 --right execution=ops/billing/exec-9ab
```

plus `--key`, `--distribution <column>` and `--json`. Mixing kinds across sides (a `record` left
and an `env` right) MUST be allowed — that is the live-vs-snapshot recovery case — and the two
sides MUST NOT be required to share an environment.

## Acceptance Criteria

### AC: env-diff-shows-changed-rows (verifies REQ:two-scope-execution, REQ:keyed-recordset-diff)

**Given** the *Shipping rate configuration* library query and the same reference rows in `prod` and `uat`, differing in one row's `rate` column and with one row present only in `prod`
**When** `datatug compare --query shipping-rate-config --left env=prod --right env=uat` runs
**Then** the result names the key column, reports one `changed` row with the `rate` values from each side, one `removed` row, a summary counting the rest as unchanged, and two distinct qualified execution refs — and a request whose key is not unique on both sides fails with an explicit error rather than diffing positionally.

### AC: cohort-distribution-highlights-overrepresentation (verifies REQ:cohort-distribution-comparison)

**Given** an affected cohort and a healthy control cohort of orders, and `shipping_provider` as the distribution column
**When** the two fact sets are compared
**Then** the response reports per-value counts and percentages for both sides with a ratio, showing the affected side's dominant provider against the control's, the output states no cause, and the surface offers that provider's diagnostic checks rather than a conclusion.

### AC: facts-cohort-is-side-scoped (verifies REQ:two-scope-execution)

**Given** one incident contains visible canonical affected facts from two projects or environments
**When** a facts side names one exact `storeId`, `project` and `environment`
**Then** only facts matching all three side qualifiers reach same-field/set validation and the native query binding, while foreign-scope facts are excluded and cannot alter or broaden the comparison.

### AC: snapshot-vs-live-shows-recovery (verifies REQ:two-scope-execution)

**Given** execution record `exec-7f3` captured at peak impact with `resultComplete:true` and a stored snapshot of 847 stuck orders
**When** it is compared against a live run of the same query in the same environment after mitigation
**Then** the diff reports the rows that are no longer stuck as `removed` with a summary reflecting the recovery, both sides carry their own qualified execution ref and timestamp, and comparing against a record whose snapshot has expired fails with an explicit `SNAPSHOT_EXPIRED` error rather than an empty left side.

### AC: compare-respects-limitations (verifies REQ:compare-is-policy-bound)

**Given** the `support` principal, for whom `prod` hides the `email` column and filters some rows, and `uat` hides nothing
**When** the two environments are compared
**Then** `email` appears nowhere in the diff or distribution, each side's `limitations` are reported separately with `email` named on the `prod` side, `policyLimited` is true because rows were filtered, and the diff is labelled so filtered rows are not presented as data differences.

### AC: public-fixture-session-compare-is-ephemeral (verifies REQ:public-fixture-browser-session)

**Given** the shared Invoice–Customer draft, a pinned public baseline fixture, a separately
published synthetic candidate fixture, and an integer CustomerId binding
**When** the analyst explicitly runs both sides and then chooses Compare by `InvoiceId`
**Then** each independent worker must exhaust its bounded statement and retain its own actual
receipt, target/hash, generation, time and complete result; the Core-backed diff shows only the
bounded query's changed raw-text value and baseline-only/candidate-only keys. The result is
labelled `Session comparison · expires when this page closes or refreshes`, stays in page memory,
and is not a server `CompareResult`, `ExecutionRef`, incident event, durable/replayable citation,
or deletion/cause claim. Parameter-free Preview, target selection and view switching start no
worker. A failed/cancelled side preserves the successful side; retry uses the frozen pair; any
changed pin/context/generation, incomplete or capped side, denied/missing side, null/duplicate key,
unsupported type, incompatible type/conflicting known units, or unequal policy visibility blocks
Compare. An output-difference cap is shown separately after complete inputs. Refresh/close clears
the session observations. Existing agent/server Compare behavior and its durable evidence ACs
remain unchanged.

### AC: compare-sides-are-execution-records (verifies REQ:compare-produces-evidence)

**Given** a compare run under incident `INC-2841`
**When** the incident timeline is read
**Then** it holds one event with `assertion.kind = deterministic-result` referencing both qualified execution refs and the key used, no copy of the diffed rows is stored in the incident, and both records remain independently readable under their own policy.

### AC: incomplete-compare-keeps-side-receipts (verifies REQ:compare-produces-evidence)

**Given** the left live side persisted its execution receipt and the right side or diff later fails
**When** the compare endpoint returns its typed error
**Then** the left receipt remains readable and is returned by the error response, no fabricated right receipt or partial `CompareResult` is returned, and no `compare.run` event is appended.

### AC: colliding-execution-ids-do-not-cross-resolve (verifies REQ:two-scope-execution, REQ:compare-produces-evidence)

**Given** evidence stores `store-a` and `store-b` each contain project `billing` execution id
`exec-7f3`, with different rows and policy, and an incident in `store-a`
**When** a compare request or incident evidence link names
`{storeId:'store-b', projectId:'billing', executionId:'exec-7f3'}`
**Then** only the `store-b` record is authorized and read, the `store-a` collision is never
considered, and omitting any qualifier is rejected rather than resolved by search order.

### AC: record-query-mismatch-is-rejected (verifies REQ:two-scope-execution)

**Given** a record side names a qualified execution whose stored `queryId` is empty or differs from `CompareRequest.queryId`
**When** the server resolves the comparison request
**Then** it rejects the request before diffing or appending `compare.run`, while records using different revisions of the same query remain comparable.

### AC: schema-comparator-remains-independent (verifies REQ:schema-compare-remains-separate)

**Given** the existing `datatug-core/pkg/comparator` schema-comparison package and the new
recordset comparison implementation
**When** package dependencies and tests are checked
**Then** the recordset comparer neither imports nor changes `pkg/comparator`, the existing
schema-comparator tests remain green, and recordset comparison has its own typed fixtures and
tests rather than adapting schema-diff output.

### AC: ordered-stream-uses-shared-cursor-merge (verifies REQ:keyed-recordset-diff)

**Given** two streams strictly ordered by a composite typed key, where either side alternates as
the lagging cursor and some equal keys are matched or changed
**When** Core compares them and a consumer observes every aligned record
**Then** DALgo `recordops.DiffFunc` supplies the alignment, the observer receives deterministic
ordinals and states in key order, only the cursor or cursors at the smallest key advance, early
observer failure closes both sources, and duplicate or descending keys fail explicitly; and the
materialized `Compare` adapter returns the same result for the same rows in arbitrary input order.
Natural-order fixtures MUST include numeric `1, 2, 10`, negative and decimal values, chronological
date/datetime values, UTF-8 binary strings and composite keys; exact-equal decimal or datetime
meanings with distinct canonical transport spellings MUST remain distinct and strictly ordered.

### AC: local-cache-is-policy-bound-and-pageable (verifies REQ:local-comparison-cache, REQ:compare-is-policy-bound)

**Given** Task 15 requests a local cache for a comparison containing matched, changed, added and
removed records
**When** the local server consumes Core's ordered observer
**Then** it writes schema-versioned canonical/minified JSON `TEXT` envelopes indexed by comparison,
state and canonical sort key; stores a matched/full one-sided row once and a changed baseline plus
right delta; returns only an opaque cache ref; and every later page read re-checks current
authorization rather than trusting the cache identity or its earlier policy revision.

## Not Doing / Out of Scope

- Schema comparison. `pkg/comparator` keeps that job; see REQ:schema-compare-remains-separate.
- Cross-source joins. Each side is one authorized execution of one query; compare never joins
  two sources into a single query — an explicit Phase 1 non-goal that this Feature does not
  quietly relax.
- Comparing more than two sides. "Compare something that works with something that does not" is
  binary; N-way comparison is not specified and not implied by the request shape.
- Rich diff *rendering* — side-by-side grids and highlighting are NEXT for the web UI; the MVP
  incident-page action renders `CompareResult` as plain tables only. Comparison boards belong to
  Dashboardius.
- Causal inference, significance testing or anomaly scoring on a distribution. Counts,
  percentages and a ratio; the human and the diagnostics do the rest.
- Configuration, deployment and infrastructure comparison (vision §13's fuller list) — FUTURE,
  and dependent on [infrastructure-inventory](../infrastructure-inventory/README.md).
- Persisting a compare result as a project asset (see REQ:compare-produces-evidence).
- Compression, SQLite JSONB, arbitrary JSON-field indexes, dynamic source columns or columnar
  storage for the local comparison cache. A future Idea covers selected-field extraction.
- Check-backed compare execution. `checkId` belongs to the Checks executable contract and is
  deferred to Incidentius Task 5; Task 8 is query-only.

## Open Questions

- **Distribution over a computed bucket.** Comparing `created_at` distributions needs bucketing
  (per hour, per server), which is a projection this Feature does not specify. Lead assumption:
  the query does the bucketing and exposes it as a column; no bucketing logic in compare.

---
*This document follows the https://specscore.md/feature-specification*
