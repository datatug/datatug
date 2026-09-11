---
format: https://specscore.md/feature-specification
status: Draft
---
# Feature: Executable checks

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/checks?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/checks?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/checks?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/checks?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

A check is a library query plus an expected condition: executable knowledge that answers one
yes/no question about a system, carries a purpose and a scope, and whose every run is an
execution record.

Checks are how a useful investigation query stops being a thing one person
remembers — they turn into the list an incident starts from, the way a hypothesis is validated,
the proof that a remediation worked, and the regression guard that runs in CI.

## Problem

Two things exist today and neither is this. `datatug-core/pkg/datatug/checks.go` defines
`Check{ID, Title, Type: options|regexp|sql, Data}` — *dataset value* checks that `datatug-action`
runs in CI to validate data shape. They have no purpose, no scope, no expected condition beyond
the value constraint, and no recorded runs. Separately, the query library holds useful queries
that nobody can assert anything about.

So an incident starts from a blank page, which is precisely what the founder rules out
(vision §19):

> An incident should not start from a blank page.

The direction is equally explicit about where checks come from and where they end up
(vision §18):

> AI can discover/propose a check, but mature checks should become deterministic where
> practical.

```text
  a query that proved useful during an incident
              │  "promote to check"
              ▼
   Check { queryId, expect, purpose, scope, cost/risk/privilege }
              │
     ┌────────┼──────────────┬────────────────┬──────────────────┐
     ▼        ▼              ▼                ▼                  ▼
   health   diagnostic  hypothesis-     remediation-        regression
     │        │          validation      verification            │
     │        │              │                │                  │
     └────────┴──────────────┴────────────────┘                  └── datatug-action in CI
                     │                                                (non-zero exit = fail)
                     ▼
          every run is an ExecutionRecord
          with outcome pass | fail | error
          (assertion.kind = deterministic-result)
```

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Check = library query + `expect` predicate + purpose/scope/metadata | planned | MVP (B1) | strategic |
| Existing `options`/`regexp`/`sql` dataset checks keep validating | exists (`checks.go`) | MVP (B1) | enabler |
| A run is an execution record with `pass`/`fail`/`error` | planned | MVP (B1) | strategic |
| Promote an executed query into a check | planned | MVP (B1) | strategic |
| List applicable checks by purpose and scope | planned | MVP (B1) | strategic |
| `datatug check list\|show\|run` and `/datatug/checks/*` | planned | MVP (B1) | surface |
| Regression checks in CI via `datatug-action` | runner exists; purpose filter planned | MVP (B1) | surface |
| Auto-run the safe health checks when an incident opens | planned | NEXT | strategic |
| Check status widgets on a board | planned | NEXT | surface |
| AI proposes a check from an investigation pattern | planned | FUTURE | strategic |

#### REQ: check-definition

A check MUST reference a **library query** by id and MUST NOT carry query text of its own. This
is the same rule the founder chose for board widgets, ruled 2026-09-09 (verbatim): *"I'm Ok
with the suggested option 1"* — a widget references a library query rather than carrying query
text or an execution target. A check inherits the query's text, targets and typed parameters,
so one definition of every query survives and a check's access control is identical to its
query's access control.

A check MUST carry an `expect` predicate, a `purpose`, and MAY carry a `scope` and cost/risk/
privilege metadata:

```typescript
type CheckSpec = {
  queryId: string;                     // library query; never inline text
  parameters?: {[parameterId: string]: TypedValue};   // fixed bindings; the rest bind from scope/context
  expect: Expectation;
  purpose: 'health' | 'diagnostic' | 'hypothesis-validation'
         | 'remediation-verification' | 'regression';
  scope?: {entities?: string[]; sources?: string[]; environments?: string[]};
  cost?: 'cheap' | 'moderate' | 'expensive';          // lead assumption: three bands
  risk?: 'safe' | 'sensitive' | 'privileged';         // lead assumption: three bands
  requiredPrivilege?: string;                         // policy name, not a grant
};

type Expectation =
  | {kind: 'row-count'; op: '=='|'!='|'<'|'<='|'>'|'>='; value: number}
  | {kind: 'single-value'; column: string; value: ValueConstraint}   // first row of a 1-row result
  | {kind: 'every-row'; column: string; value: ValueConstraint};     // all rows satisfy

type ValueConstraint =
  | {kind: 'options'; type: string; options: unknown[]}  // OptionsValueCheck payload, verbatim
  | {kind: 'regexp'; regexp: string}                     // RegexpValueCheck payload, verbatim
  | {kind: 'compare'; op: '=='|'!='|'<'|'<='|'>'|'>='; value: TypedValue};
```

The `options` and `regexp` constraints MUST reuse the existing evaluation kinds' payloads
field-for-field, so the code that evaluates a dataset column check evaluates a query column
check — no second evaluator. `requiredPrivilege` names the policy a principal must satisfy; it
never grants anything and never appears on the wire as an authorization claim. The operator set
matches the optional `Fact.condition` proposed in the transport appendix (a pending lead
assumption, not yet in `datatug-core`), so one comparison vocabulary serves context, checks
and compare once both land.

#### REQ: existing-dataset-checks-retained

The current `Check{ID, Title, Type: options|regexp|sql, Data}` in
`datatug-core/pkg/datatug/checks.go` MUST stay valid. Every project file that validates today
MUST validate unchanged after this Feature lands, and `datatug validate` / `datatug-action`
MUST keep running those checks exactly as they do now.

The new form is an additive fourth evaluation kind in the same `Check.Type` switch (lead
assumption): `type: "query"`, whose `data` holds `CheckSpec`'s `queryId`, `parameters` and
`expect`. `purpose`, `scope`, `cost`, `risk` and `requiredPrivilege` are new **optional**
top-level fields on `Check`, so an existing `options`/`regexp`/`sql` check gains a purpose
without changing its `data`. An unknown `type` MUST keep failing validation as it does today,
rather than being accepted as a forward-compatible extension.

#### REQ: check-run-is-execution-record

Running a check MUST execute its query through the same executor as any other run, under the
caller's own policy, and MUST produce an [execution record](../evidence-records/README.md)
(`evidence-records` REQ:execution-record-captured) plus a check outcome:

```typescript
type CheckRun = {
  checkId: string;
  recordId: string;                       // the execution record this run produced
  outcome: 'pass' | 'fail' | 'error';
  evaluated: {expect: Expectation; actual: TypedValue | number; explanation: string};
  evaluatedAt: string;
};
```

`error` MUST be distinct from `fail`: a source that is unavailable, a policy that denied the
read, or a timeout is an `error`, never a failing check. A run whose result was restricted —
any `Limitation` with `rowsFiltered: true` — MUST NOT be reported as `pass` or `fail` on a
row-count expectation, because the count the check saw is not the count that exists; it MUST be
`error` with the limitation named (lead assumption, and the safe direction).

A check run under an incident MUST enter the incident event stream with
`assertion.kind = deterministic-result` (product-family §5, D12), carrying the record id and the
evaluated predicate, so the timeline can distinguish a measurement from an opinion.

#### REQ: promote-query-to-check

From an executed query result a user MUST be able to create a check with an expected condition,
pre-filled from what they just saw: the query id and its bindings, a suggested expectation
derived from the actual result (row count `== n`, or `every-row` over a selected column), and a
purpose the user picks. The created check MUST be a project file, reviewable as a diff like any
other knowledge asset.

This needs the project-write path from the hub's **Phase 2 plan task 2**; until that path
exists, the promote action MUST be absent or explicitly disabled with the reason shown, never
offered and then failing. It is the same entry point `dashboards` REQ:board-from-exploration
establishes for boards ("pin this result") applied to assertions, and it is the mechanism by
which vision §18's progression from human investigation to deterministic check actually happens.

#### REQ: checks-listed-by-purpose-and-scope

Given the entities, sources and environment an investigation is currently touching, the system
MUST be able to list the applicable checks, filtered by purpose. Applicability MUST be
deterministic and explainable: a check is applicable when its `scope` names an entity, source or
environment present in the current scope, or when it declares no scope and its query's
parameters can be bound from the active facts — the same reasoning `executable-knowledge-library`
REQ:query-applicability already applies to queries, with the same explanation chain.

An incident therefore opens with a list of relevant `health` and `diagnostic` checks rather than
a blank page (vision §19). MVP lists them and runs them on demand; automatically running the
`risk: safe`, `cost: cheap` ones when an incident opens is NEXT. Checks the current principal
may not run MUST be omitted entirely, exactly as unauthorized queries are.

#### REQ: ci-runs-regression-checks

`datatug check run --purpose regression` MUST run every regression-purpose check in scope and
exit non-zero when any of them is `fail` or `error`, so a shell or a CI job can gate on it. The
output MUST name each failing check, its expectation and its actual value, and MUST support
`--json` for machine consumption.

`datatug-action` MUST be able to invoke it in a GitHub workflow the same way it already runs the
validator, so a check created during an incident guards the next release — the flywheel
product-family §7 lists under "datatug-action / CI". A regression check MUST NOT require a
snapshot or an incident to run; it needs only a project, an environment and a principal.

#### REQ: ai-proposed-check-needs-human-review

An AI-proposed check MUST NOT become a project check without an explicit human action, and its
proposal MUST be recorded as `assertion.kind = inference` with the events it was derived from.
This is the candidate-then-validate pattern product-family §7 names, applied here: the LLM may
draft a `CheckSpec`, the deterministic run and the human decide. No check ever enters the
library, is run automatically, or is cited as evidence on the strength of a model's output
alone. This capability is FUTURE; the REQ exists so nothing implements the easy half first.

#### REQ: api-and-cli

The server MUST expose, under `/datatug`, following the existing scope, staleness and error
envelope rules:

| Method/path | Request | Success envelope |
|---|---|---|
| GET `checks` | Scope + optional `purpose`, `entity`, `source` | `{checks: CheckBrief[]}` — authorized checks only |
| GET `checks/{id}` | Scope | `{check: Check}` |
| POST `checks/{id}/run` | Scope + `{parameters?, record?: boolean, snapshot?: boolean, incidentId?}` | `{run: CheckRun, result?: Result}` |
| POST `checks` | Scope + `{fromRecordId, expect, purpose, …}` | `{check: Check}` — the promote path, project-write gated |

The CLI follows the singular-namespace convention: `datatug check list [--purpose …]`,
`datatug check show <id>`, `datatug check run <id>|--purpose <purpose>`, each with `--json`.
CLI and web MUST use these endpoints; neither evaluates an expectation locally, so "did it
pass?" has exactly one answer everywhere.

## Acceptance Criteria

### AC: check-run-records-result (verifies REQ:check-run-is-execution-record)

**Given** a check `finalizer-backlog` bound to the *Finalizer queue depth* library query with `expect: {kind: 'row-count', op: '<', value: 100}`
**When** it is run against `prod` with `record: true`
**Then** the response carries `outcome: 'fail'` with the evaluated expectation and the actual count, an execution record exists for the underlying run with matching `rowCount` and `limitations`, and running it against a source that is down returns `outcome: 'error'` rather than `fail`.

### AC: promote-query-creates-check-file (verifies REQ:promote-query-to-check)

**Given** an executed *Stuck orders* result with 847 rows and the project-write path available
**When** the user promotes it to a `health` check with the suggested expectation `row-count == 0`
**Then** a reviewable check file is written referencing the library query id and its bindings, holding no query text of its own; and on a read-only project the promote action is disabled with its reason shown rather than offered and failing.

### AC: health-checks-listed-for-entity (verifies REQ:checks-listed-by-purpose-and-scope)

**Given** an investigation whose active scope touches the `Order` entity in `prod`, a `health` check scoped to `Order`, a `diagnostic` check scoped to `Payment`, and a check the current principal may not run
**When** applicable checks are listed with `purpose=health`
**Then** the `Order` health check is listed with the explanation of why it applies, the `Payment` check is not, and the unauthorized check is absent with no count of omissions.

### AC: regression-check-fails-ci (verifies REQ:ci-runs-regression-checks)

**Given** two regression checks, one passing and one failing, in a project used by a CI job
**When** `datatug check run --purpose regression --json` runs in that job
**Then** the process exits non-zero, the output names the failing check with its expectation and actual value, and the passing check is reported as `pass` — and with both passing, the process exits zero.

### AC: legacy-check-types-still-validate (verifies REQ:existing-dataset-checks-retained)

**Given** a project containing existing `options`, `regexp` and `sql` checks that validate on today's `datatug-core`
**When** `datatug validate` runs after this Feature lands
**Then** all three still validate unchanged with no file edits, a new `type: "query"` check validates alongside them, and an unknown check type is still rejected.

### AC: check-definition-rejects-inline-sql (verifies REQ:check-definition)

**Given** a candidate check file whose `data` carries SQL text and no `queryId`
**When** the project is validated
**Then** it is rejected with an error naming the library-query rule, and the same check expressed as `{queryId, expect}` validates.

## Not Doing / Out of Scope

- A monitoring or alerting system: checks run when a person, a CI job or (later) an incident
  asks. There is no scheduler, no alert routing, no paging — product-family §1 keeps incident
  *management* out of the family entirely.
- A second query engine or expression language. An expectation is a small closed predicate over
  a result the existing executor produced; it is not a rules DSL and does not embed SQL.
- Automatic remediation. A `remediation-verification` check confirms a fix; it never applies one.
  Controlled production modifications (vision §36) are a separate FUTURE track.
- Replacing `datatug-core`'s dataset checks or migrating existing project files.
- Check *scheduling*, history retention policy or trend alerting: the history of a check is its
  execution records, governed by [evidence-records](../evidence-records/README.md).
- Incident semantics — which checks constitute resolution criteria, and what a failing check
  means for a hypothesis — belong to [incidents](../incidents/README.md).

## Open Questions

- **Where a check file lives.** Lead assumption: `checks/<id>.json` as a project asset, beside
  queries and boards, so it is granted and diffed like every other knowledge asset. The existing
  dataset checks live inside dataset definitions; whether the two should ever share a directory
  is undecided.
- **Cost and risk bands.** Three values each (`cheap|moderate|expensive`, `safe|sensitive|privileged`)
  are a lead assumption chosen to be usable for the NEXT "auto-run the safe cheap ones" rule.
  Numeric budgets would be more precise and less likely to be filled in honestly.
- **Parameter binding at run time.** A check with unbound parameters can bind them from the
  Investigation Context, which makes "the same check" mean different things in two runs. Lead
  assumption: a check declares its fixed `parameters`, and any remaining required parameter must
  be supplied explicitly at run time — never auto-bound from context — so a check's identity is
  stable. Confirm before hypothesis-validation checks are used as evidence.
- **`sql` type overlap.** The existing `SQLCheck{Query}` holds SQL text, which the new form
  forbids. Lead assumption: it is retained as-is for dataset validation and is not extended;
  new work uses `type: "query"`. Whether `sql` should eventually be deprecated is a separate
  decision with a migration cost.

---
*This document follows the https://specscore.md/feature-specification*
