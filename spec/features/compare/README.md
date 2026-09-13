---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Compare

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/compare?op=request-change) |
**Status:** Approved
**Source Ideas:** —

## Summary

Compare executes one library query or check under two scopes and diffs the two recordsets by
key.

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
      one library query or check
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

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Two-scope execution: env A vs env B | planned | MVP (B1) | strategic |
| Two-scope execution: fact set A vs fact set B (cohorts with roles) | planned | MVP (B1) | strategic |
| Two-scope execution: execution record A vs execution record B | planned | MVP (B1) | strategic |
| Keyed recordset diff with per-column differences | planned | MVP (B1) | strategic |
| Cohort distribution comparison for a categorical column | planned | MVP (B1) | strategic |
| Per-side limitations; policy-limited diffs labelled | planned | MVP (B1) | strategic |
| `POST /datatug/compare` and `datatug compare` | planned | MVP (B1) | surface |
| Compare result referenced as incident evidence | planned | MVP (B1) | strategic |
| Compare UI: one "Compare affected vs control" action on the incident page, rendering `CompareResult` as plain tables (summary; changed/added/removed rows; per-side limitations; distribution) | planned | MVP (B1) | surface |
| Compare UI: side-by-side grids, diff highlighting | planned | NEXT | surface |
| Comparison boards in Dashboardius | planned | NEXT | surface |
| Comparing configuration, deployments and infrastructure (vision §13) | planned | FUTURE | strategic |

#### REQ: two-scope-execution

Compare MUST take exactly one executable — a library query id or a check id — and two side
specifications. A `scope` or `facts` side MUST execute the query live, separately under that
side's own scope; a `record` side MUST NOT re-execute anything and instead reuses the named
past execution record. The three side kinds are:

```typescript
type SideSpec =
  | {kind: 'scope'; storeId: string; project: string; environment: string; parameters?: {[id: string]: TypedValue}}
  | {kind: 'facts'; storeId: string; project: string; environment: string; facts: Fact[]} // one cohort
  | {kind: 'record'; execution: ExecutionRef};                                // a past execution
```

- **environment A vs B** — the same query, two environments (`{kind:'scope'}` on each side).
- **fact set A vs B** — the same query and environment, two sets of facts: the affected cohort
  and the control. Cohort roles are the [investigation-context](../investigation-context/README.md)
  extension that names which set a fact belongs to; compare consumes them and does not
  reimplement them (product-family §3, the firm boundary).
- **execution record A vs B** — two past runs, i.e. snapshots or time. A `record` side MUST have
  a readable result snapshot; a record with only a fingerprint and a row count supports a
  fingerprint-equality answer and MUST NOT be silently expanded into a row diff. A `record` side
  whose snapshot has expired MUST fail the whole compare with error code `SNAPSHOT_EXPIRED`,
  never an empty side (AC:snapshot-vs-live-shows-recovery).

A live (`scope` or `facts`) side MUST produce its own new execution record (`evidence-records`
REQ:execution-record-captured); a `record` side reuses its existing execution record without
re-execution. Either way the two halves of a comparison are independently citable evidence. A
live side is executed as one ordinary authorized run; there is no combined query, no cross-source
join and no second executor — this Feature adds a caller and a differ, not an engine.

A `facts` side binds through the transport appendix's bounded `TypedValueSet`: exactly one
same-field equality cohort of 1..500 unique, non-null, same-typed values. The query parameter
must declare set support and the source must execute one native set/array/`IN` binding. Any
other cohort shape is `INVALID_REQUEST`; a source without safe set binding is
`SOURCE_UNAVAILABLE`. Compare never loops one execution per fact or silently truncates a
cohort.

#### REQ: keyed-recordset-diff

The two recordsets MUST be diffed by key, producing `added`, `removed` and `changed` rows plus
an aggregate summary. Key columns are resolved in this order: the caller's explicit `key`
columns; otherwise the key fields of the entity the result's columns map to, via the semantic
layer's declared mappings. If neither resolves to a key that is unique on both sides, the
request MUST fail with an explicit error naming the ambiguity — never fall back to positional
or whole-row comparison, which produces a plausible and wrong answer.

```typescript
type CompareResult = {
  left:  {execution: ExecutionRef; executedAt: string; rowCount: number; limitations: Limitation[]};
  right: {execution: ExecutionRef; executedAt: string; rowCount: number; limitations: Limitation[]};
  key: string[];
  added:   {key: TypedValue[]; row: TypedValue[]}[];     // present right, absent left
  removed: {key: TypedValue[]; row: TypedValue[]}[];     // present left, absent right
  changed: {key: TypedValue[]; columns: {column: string; left: TypedValue; right: TypedValue}[]}[];
  summary: {added: number; removed: number; changed: number; unchanged: number; columnsOnlyOnOneSide: string[]};
  distribution?: Distribution;                            // REQ:cohort-distribution-comparison
  policyLimited: boolean;                                 // REQ:compare-is-policy-bound
  truncated: boolean;
};
```

Comparison of values MUST be typed, using the same `TypedValue` equality the binding contract
fixes: no string/number coercion, `null` equals only `null`, and a `decimal` compares by its
canonical decimal form, not as a float. Only columns present on **both** sides may appear in
`changed`; a column present on one side only MUST be reported in `summary.columnsOnlyOnOneSide`
as a shape difference, not as every row having changed. Results MUST be bounded by the same default and
maximum result limits the transport contract sets, with `truncated` set when the cap is hit.

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
([incidents](../incidents/README.md) REQ:event-stream-append-only). Re-running the same compare later produces two new records; the earlier comparison
stays valid evidence of what was true then.

For a compare that executes either side under `CompareRequest.incident`, the server normalizes
that ref onto each new execution record's top-level `incident` field. It MUST NOT place the
incident inside `Result.provenance`; each result and the compare response expose only the
qualified `ExecutionRef`s of the recorded sides.

MVP does not persist a `CompareResult` as its own project asset (lead assumption): the two
execution records plus the request parameters reproduce it, and storing the diff would create a
third copy of data already governed by retention and masking. An incident references the pair.

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
  queryId?: string; checkId?: string;      // exactly one required
  left: SideSpec; right: SideSpec;
  incident?: IncidentRef;
  key?: string[];
  distributionColumn?: string;
  limit?: number;
};
```

It is a POST because both sides carry semantic values, which must never reach a URL, browser
history or an access log — the same reason `semantic/related` is a POST. Supplying neither or
both of `queryId`/`checkId` is INVALID_REQUEST.

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

### AC: snapshot-vs-live-shows-recovery (verifies REQ:two-scope-execution)

**Given** execution record `exec-7f3` captured at peak impact with a stored snapshot of 847 stuck orders
**When** it is compared against a live run of the same query in the same environment after mitigation
**Then** the diff reports the rows that are no longer stuck as `removed` with a summary reflecting the recovery, both sides carry their own qualified execution ref and timestamp, and comparing against a record whose snapshot has expired fails with an explicit `SNAPSHOT_EXPIRED` error rather than an empty left side.

### AC: compare-respects-limitations (verifies REQ:compare-is-policy-bound)

**Given** the `support` principal, for whom `prod` hides the `email` column and filters some rows, and `uat` hides nothing
**When** the two environments are compared
**Then** `email` appears nowhere in the diff or distribution, each side's `limitations` are reported separately with `email` named on the `prod` side, `policyLimited` is true because rows were filtered, and the diff is labelled so filtered rows are not presented as data differences.

### AC: compare-sides-are-execution-records (verifies REQ:compare-produces-evidence)

**Given** a compare run under incident `INC-2841`
**When** the incident timeline is read
**Then** it holds one event with `assertion.kind = deterministic-result` referencing both qualified execution refs and the key used, no copy of the diffed rows is stored in the incident, and both records remain independently readable under their own policy.

### AC: colliding-execution-ids-do-not-cross-resolve (verifies REQ:two-scope-execution, REQ:compare-produces-evidence)

**Given** evidence stores `store-a` and `store-b` each contain project `billing` execution id
`exec-7f3`, with different rows and policy, and an incident in `store-a`
**When** a compare request or incident evidence link names
`{storeId:'store-b', projectId:'billing', executionId:'exec-7f3'}`
**Then** only the `store-b` record is authorized and read, the `store-a` collision is never
considered, and omitting any qualifier is rejected rather than resolved by search order.

### AC: schema-comparator-remains-independent (verifies REQ:schema-compare-remains-separate)

**Given** the existing `datatug-core/pkg/comparator` schema-comparison package and the new
recordset comparison implementation
**When** package dependencies and tests are checked
**Then** the recordset comparer neither imports nor changes `pkg/comparator`, the existing
schema-comparator tests remain green, and recordset comparison has its own typed fixtures and
tests rather than adapting schema-diff output.

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

## Open Questions

- **Key resolution when a query's columns are unmapped.** Lead assumption: explicit `--key` is
  always accepted, semantic entity keys are used when available, and anything else is a hard
  error. Whether a single obviously-unique column should be auto-adopted is undecided; the safe
  default is to refuse.
- **Column-set intersection versus failure.** MVP intersects the two sides' visible columns and
  reports the difference. The stricter alternative — refuse to compare when the visible column
  sets differ — is less useful during an incident and arguably more honest. Confirm before the
  first cross-environment comparison with asymmetric policies.
- **Distribution over a computed bucket.** Comparing `created_at` distributions needs bucketing
  (per hour, per server), which is a projection this Feature does not specify. Lead assumption:
  the query does the bucketing and exposes it as a column; no bucketing logic in compare.
- **Fingerprint-only comparison.** Two records with no snapshots can still answer "same result or
  not" from `resultFingerprint`. Whether that is a compare mode or a separate cheap call is
  undecided; it is not in MVP.

---
*This document follows the https://specscore.md/feature-specification*
