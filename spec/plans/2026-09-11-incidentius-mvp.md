---
format: https://specscore.md/plan-specification
status: Draft
---
# Plan: Incidentius MVP inside DataTug

**Status:** Draft
**Source Feature:** incidents
**Related Features:** evidence-records, checks, annotations, compare, product-profiles, investigation-context, dashboards, server-enforced-access, provenance, executable-knowledge-library
**Prerequisite Plans:** 2026-09-09-phase-1-core-investigation-loop, 2026-09-09-phase-2-library-depth-and-provenance
**Date:** 2026-09-11
**Owner:** alex
**Supersedes:** —

## Summary

Deliver roadmap Track B1: incidents as a DataTug capability, inside the existing application
and CLI, on a local or Git project under the single `serve` principal. The plan proves the
thesis with the smallest coherent slice (vision §59): the incident store and its event stream,
attachment to Investigation Context with cohorts and hypothesis overlays, evidence by
reference to recorded executions, checks, compare runs and annotations, metrics with explicit
criteria, recurrence and search, back-links from DataTug assets, API/CLI/`watch`, a
deterministic status view, and the canonical demo incident with its boring recurrence, walked
end to end by a real e2e. It deliberately excludes hosted multi-user collaboration,
Dashboardius boards, AI narratives, capture channels and privileged production access; those
are Track B2+ and stay named in their owning Features.

**Dependencies, stated precisely.** Phase 1 is accepted (2026-09-11). Incidents, execution
records and snapshots are written to the incident/evidence store through the server's own
store path (`incidents` REQ:incident-store) and do not wait for anything. Checks and
annotations are knowledge assets written into the project, so tasks 5 and 7 (and therefore
the parts of tasks 9, 11 and 12 that use them) need the authorized project-write path that
`2026-09-09-phase-2-library-depth-and-provenance` task 2 introduces for captured queries;
tasks 5 and 7 generalise it to `checks/` and `annotations/` files (a lead assumption: that
task writes query files only). Phase 2 is therefore a prerequisite plan, and tasks 5 and 7
name its task 2 as an external dependency. No task may add a second context model, executor, policy evaluator,
renderer or storage path.

## Journey

Written in the user's own sentences; each stage names its observable good result. Nobody
performs hidden setup between stages.

**Start, null action.** "I run `datatug demo`, the browser opens the demo project, and I see an
*Incidents* entry in the project menu; under the Incidentius profile it is the home page."
Good result: the incident list renders (empty or with the fixture) with no error and no
hidden URL.

**Houston.** "Support says Canadian customers cannot see invoice totals. I click *Houston,
we've got a problem*, type a title and paste the message with customer 5 and invoice 42." Good
result: INC-1 exists in the project's own `incidents/` store, the reporter is a participant,
the canonical context holds the two facts I typed with their project scope, and
`datatug incident show INC-1` prints the same.

**Cohorts.** "I open the Customer table from the incident, select Canadian customers, mark
them *affected*, and mark three Brazilian customers *healthy control*." Good result: the
canonical layer shows both roles; a query opened from the incident binds from the affected
facts and lists the bindings in its results header (no hidden filters).

**Checks first.** "The incident page already lists the health checks for Invoice and Customer;
I run them." Good result: each run is recorded as an execution record and a `check.run` event;
one check fails ("invoice totals present") and shows its predicate and rows count.

**Access.** "The invoice query I need says *Invoice: not visible to you*. I click *Request
access*, name the collection and the reason, and the lead approves it from the incident
page." Good result: `access.requested` and `access.approved` are events with their actors,
the query now returns rows with `Customer.Email` still masked, its execution record names
the grant, and `datatug incident show INC-1` lists the grant as active with its expiry.

**Hypotheses.** "John proposes 'the FX-rate provider is down' with the five Canadian invoices as
its scope; Jack proposes 'the 14:17 deployment changed rounding'." Good result: two hypotheses
in `proposed`, each with an overlay layer, each visible in the timeline with its actor.

**Evidence.** "Jack runs the currency-rate query for CAD, records it, and attaches the record
to John's hypothesis as weakening evidence: the live rate came back fine. He annotates
`Invoice.ID = 42`: 'total missing, rate present'." Good result: the execution record carries
provenance, limitations and INC-1's reference; the annotation resurfaces when invoice 42
appears in any other result; John's hypothesis moves to `weakened` with the reason and record
linked.

**Compare.** "Maria runs the same invoice query for the affected set and the control set." Good
result: the compare result shows the affected side missing totals and the control side complete,
each side under its own limitations, recorded as a `compare.run` and attached as evidence.

**Metrics and criteria.** "We add the metric *stuck invoices* (impact) and the criterion *stuck
invoices == 0 for 30 minutes*." Good result: the series starts with the first recorded point;
the criterion reads unmet.

**Nothing happens for an hour, null action.** "Everybody is in a meeting. Nobody opens the
incident, runs anything or types anything." Good result: the projection is unchanged, the
`watch` stream emits nothing, the last metric point still shows its recorded time, and
reopening the page shows exactly the earlier state.

**Mitigation and recovery.** "Anna applies the fix; the metric's check is run again twice."
Good result: the series reads 37 → 19 → 0, the criterion flips to satisfied on the recorded
evidence, status moves through `mitigating`, `recovering`, `resolved` with outcome `resolved`,
then `watching` and `closed`.

**Record.** "A newcomer opens the Resolution Record. A support colleague, whose policy hides
customer emails, opens it too." Good result: timeline, hypotheses with their evidence including
the rejected one, metrics, criteria, outcome; every evidence item opens under the reader's own
policy; the support colleague sees the email fact redacted in the timeline and the projection;
`show --at 14:42` shows what was believed then.

**Epilogue A, three months later.** "The same symptom appears; I create INC-2." Good result:
`similar` names INC-1 with the matching signals; I link `recurrence-of`; the previously created
check runs; the recurrence resolves in minutes with a recorded TTRI; the `customer-invoices`
query lists both incidents.

**Epilogue B, an agent.** "A script runs `datatug incident watch --json --since <cursor>` while
all of the above happens." Good result: it receives every event exactly once as NDJSON and can
resume from the cursor printed on stderr.

## Approach

- Reconcile before implementing: refresh `origin/main` of `datatug-core`, `datatug-cli`,
  `datatug-apps`, `datatug-demo-projects`; read the transport appendix and the sibling Features
  (`evidence-records`, `checks`, `annotations`, `compare`, `product-profiles`) at their
  current revisions. Component presence is not journey acceptance.
- Model first, in `datatug-core`: the incident record, `IncidentRef`, store routing and
  interface, event schema, projection and fold, the execution record. Server second, in
  `datatug-cli`: the DALgo-backed store implementations (project repository, dedicated
  repository, application repository), endpoints, recording, checks runner, compare, incident
  stream, policy re-check on reads, CLI verbs. Browser third, in `datatug-apps`: incident pages,
  context roles/overlays, evidence and metrics tables, the `incidentius` profile. Demo last:
  fixtures and the whole-journey e2e that gates the release candidate.
- Contract amendments (optional `Fact.role`, `Fact.layer`; optional `Scope.storeId`;
  optional `ExecutionRequest.incident`;
  `record`/`snapshot` flags on execution; `Result.provenance.incident` as an `IncidentRef`;
  the incident endpoints; `Board.incidentRef` as an `IncidentRef`; the opened board widget
  name boundary) are frozen as fixtures in `datatug-core/pkg/apicontract` before server
  or client work, as Phase 1 task 12 did.
- Independent tasks run in separate WB worktrees with one landing owner per repository; each
  repository lands one candidate per task, never independent declarations of a finished journey.
- Track status with `specscore task change-status --plan 2026-09-11-incidentius-mvp`; a task
  completes only with implementation evidence and passing ACs.
- Every task's output goes through an adversarial reviewer before landing (founder standing
  rule, 2026-09-11).

## Tasks

### Task 1: Incident model, store routing, event schema, projection

**Id:** task-1
**Verifies:** incidents#ac:event-stream-append-only-and-replayable, incidents#ac:show-at-reconstructs-belief, incidents#ac:inference-requires-refs, incidents#ac:outcomes-and-duplicate-merge, incidents#ac:incident-store-routes-to-configured-location, incidents#ac:incident-spans-two-projects
**Depends-On:** —
**Status:** planning

In `datatug-core` and `datatug-cli`: `Incident`, `IncidentRef{storeId, incidentId}` and its
canonical string form, referenced projects, `IncidentEvent` with actor/assertion (including
`question`)/refs, the typed payloads, the deterministic fold to the projection (with
`at <= T`), lifecycle/outcome validation, merge semantics, the incident-store interface with
the store routing key and the three DALgo-backed locations (project repository under
`incidents/`, dedicated repository, application repository) sharing one layout, and the
`apicontract` fixtures for `/datatug/incidents/*`. 100% covered; byte-identical fold test.

### Task 2: Execution records and opt-in snapshots in the evidence store

**Id:** task-2
**Verifies:** incidents#ac:hypothesis-rejected-keeps-evidence
**Supporting ACs:** evidence-records#ac:record-created-on-request, evidence-records#ac:fingerprint-without-snapshot, evidence-records#ac:snapshot-read-rechecks-policy, evidence-records#ac:retention-expires-snapshot-keeps-record, evidence-records#ac:series-from-repeated-runs, evidence-records#ac:snapshot-read-is-not-an-execution, server-enforced-access#ac:snapshot-read-under-current-policy, provenance#ac:recorded-execution-keeps-chain
**Depends-On:** 1
**Status:** planning

In `datatug-core` + `datatug-cli`: `ExecutionRecord` with canonical fingerprint and optional
`IncidentRef` provenance, `record`/`snapshot` flags on `exec/run_query` and related-rows,
records and snapshots written to the evidence store that follows the incident-store routing,
retention and masking policy, snapshot reads re-checked through `accesspolicies.Run`,
`GET /datatug/executions*`, `datatug execution list|show`. Negative tests: a principal
without current access to the recorded columns gets a non-disclosing denial.

### Task 3: Incident endpoints, policy-rechecked reads, CLI verbs and watch stream

**Id:** task-3
**Verifies:** incidents#ac:houston-creates-incident-with-context, incidents#ac:watch-streams-events-with-cursor, incidents#ac:cli-json-matches-api, incidents#ac:events-redacted-under-policy
**Supporting ACs:** evidence-records#ac:record-referenced-from-incident
**Depends-On:** 1, 2
**Status:** planning

In `datatug-cli`: list (with query/check/board filters), create, show (`?at=`), append,
events (NDJSON, cursor), search, similar, merge, under the serve principal with store-write
capability checks; every read of events, projection and stream applies the current
principal's policies (fact redaction, annotation visibility); `datatug incident` verbs per
`cli/incident` including `merge`; `watch` renders one human line per event on stdout or raw
NDJSON, cursor on stderr. Contract test against the real server; no polling loop in the CLI.

### Task 4: Context cohort roles and hypothesis overlays

**Id:** task-4
**Verifies:** incidents#ac:investigation-becomes-incident, incidents#ac:overlay-promotion-updates-canonical
**Supporting ACs:** investigation-context#ac:cohort-role-visible-and-bindable, investigation-context#ac:promotion-moves-fact-with-history
**Depends-On:** 1
**Status:** planning

Amend the transport appendix (optional `Fact.role`, `Fact.layer`) and `apicontract`; extend
the Investigation Context page and service in `datatug-apps` with roles, layers and
promote/reject that emit incident events; "make this an incident" from any investigation.
An open query re-binding from a promotion shows the explicit conflict choice, never a silent
rebind.

### Task 5: Checks with purposes and expected conditions

**Id:** task-5
**Verifies:** incidents#ac:checks-run-at-start-recorded
**Supporting ACs:** checks#ac:check-run-records-result, checks#ac:health-checks-listed-for-entity, checks#ac:legacy-check-types-still-validate, checks#ac:promote-query-creates-check-file, checks#ac:check-definition-rejects-inline-sql, checks#ac:regression-check-fails-ci, executable-knowledge-library#ac:useful-query-becomes-check
**Depends-On:** 2
**Status:** planning

In `datatug-core` + `datatug-cli`: check definitions referencing library queries with `expect`
predicates, purposes and scope; the runner records an execution with a deterministic result;
`datatug check list|show|run` with the regression exit code; listing by purpose and scope for
the incident page; "promote query → check" writing the check file with `origin.incident`.
Blocked on the project-write path from `2026-09-09-phase-2-library-depth-and-provenance`
task 2 for every write of a check file; the runner and listing do not wait.

### Task 6: Metrics series and resolution criteria

**Id:** task-6
**Verifies:** incidents#ac:metric-series-and-criteria
**Depends-On:** 1, 2, 5
**Status:** planning

Metric definitions (check or query + numeric projection + category), `metric.observed` points
from recorded runs, criteria as predicates evaluated from evidence, projection fields, and the
numeric series/criteria tables in the incident page. No chart code in DataTug until the shared
renderer lands (dashboards REQ:no-charts-in-datatug-web-ui).

### Task 7: Annotations on values and execution records

**Id:** task-7
**Verifies:** incidents#ac:hypothesis-rejected-keeps-evidence
**Supporting ACs:** annotations#ac:annotate-value-resurfaces-in-other-query, annotations#ac:annotation-carries-provenance, annotations#ac:annotation-hidden-when-target-denied, annotations#ac:annotation-on-execution-record, annotations#ac:annotation-never-stores-restricted-rows
**Depends-On:** 2
**Status:** planning

`annotations/` project asset written through the server, `POST/GET /datatug/annotations`,
`datatug annotation add|list`, grid marker and context-panel resurfacing in `datatug-apps`,
visibility filtered by the target's current policy, `origin.incident` when created under an
incident. Blocked on the project-write path from `2026-09-09-phase-2-library-depth-and-provenance`
task 2 for the write endpoint; read-side resurfacing over fixture annotations does not wait.

### Task 8: Compare API and CLI (two scopes)

**Id:** task-8
**Verifies:** incidents#ac:demo-incident-journey-e2e
**Supporting ACs:** compare#ac:env-diff-shows-changed-rows, compare#ac:cohort-distribution-highlights-overrepresentation, compare#ac:compare-respects-limitations, compare#ac:compare-sides-are-execution-records, compare#ac:snapshot-vs-live-shows-recovery
**Depends-On:** 2
**Status:** planning

`POST /datatug/compare` and `datatug compare` for environment-vs-environment, fact-set-vs-
fact-set (cohorts) and record-vs-record, keyed recordset diff and categorical distribution,
each side its own recorded execution under its own policy; a `compare.run` event with
`deterministic-result` when invoked under an incident; attachable as evidence. Compare UI is
NEXT; this task's contribution to the journey is verified by the whole-journey AC because the
compare step is part of it.

### Task 9: Incident pages and the incidentius profile

**Id:** task-9
**Verifies:** incidents#ac:status-projection-carries-provenance, incidents#ac:participant-role-separate-from-permission, incidents#ac:incidentius-profile-home
**Supporting ACs:** product-profiles#ac:hostname-selects-home, product-profiles#ac:deep-link-works-in-other-profile, product-profiles#ac:one-build-serves-both-profiles, product-profiles#ac:profile-override-grants-nothing
**Depends-On:** 3, 4, 6
**Status:** planning

In `datatug-apps`: incident list/detail/create pages, timeline, hypotheses and evidence
panels, participants with roles, referenced projects, status view with refs, Resolution Record
page; the profile table and hostname/config selection for the two profiles one build serves;
menu entry under the `datatug` profile. Ionic frame; no chart islands in this task.

### Task 10: Recurrence, related links, search, similarity and back-links

**Id:** task-10
**Verifies:** incidents#ac:recurrence-link-records-ttri, incidents#ac:search-finds-similar-by-entities, incidents#ac:query-lists-its-incidents
**Depends-On:** 3, 5
**Status:** planning

Typed links as events (including cross-store `IncidentRef`s), TTRI derivation, text + fact
search across configured stores, deterministic explainable similarity over entities/fields/
values/queries/checks/failed checks/environments/projects/outcome, derived back-links
(`incident list --query|--check|--board`), in API, CLI and the incident page.

### Task 11: Canonical demo incident fixtures

**Id:** task-11
**Verifies:** incidents#ac:demo-incident-journey-e2e
**Depends-On:** 1, 2, 4, 5, 6, 7, 8
**Status:** planning

In `datatug-demo-projects`: author the story (open question in the Feature: Chinook first;
founder reviews the written story before fixture authoring), the fixture
`incidents/INC-1/events.jsonl` in the project's own store with several actors, useful wrong
theories, red herrings, cohorts, an environment comparison, the FX-rate provider as external
counterparty, recorded executions and snapshots, annotations, checks, compare runs, metrics,
mitigation, recovery, resolution; the recurrence fixture `INC-2`; validator coverage for the
incident store layout.

### Task 12: Whole-journey e2e and release candidate

**Id:** task-12
**Verifies:** incidents#ac:demo-incident-journey-e2e, incidents#ac:recurrence-epilogue-boring, incidents#ac:houston-creates-incident-with-context, incidents#ac:investigation-becomes-incident, incidents#ac:watch-streams-events-with-cursor, incidents#ac:events-redacted-under-policy
**Depends-On:** 3, 4, 6, 7, 8, 9, 10, 11, 14
**Status:** planning

Playwright + CLI journey from `datatug demo` through the Resolution Record, the null-action
hour and both epilogues, against the real server with no intercepted responses, asserting
events, records, redaction and cursors; required CI job on `datatug-apps` and `datatug-cli`;
pinned revision set (core, cli, apps, demo, providers); a named datatug-cli release and
datatug.app deployment receipt.

### Task 13: Board widget extension boundary and incident reference

**Id:** task-13
**Verifies:** incidents#ac:board-lists-its-incidents-by-ref
**Supporting ACs:** dashboards#ac:incident-board-round-trips-widget-names
**Depends-On:** 1, 10
**Status:** planning
**Activation:** Track C enabler; not on the MVP journey path (task 12 does not depend on it)

Open the board contract for Track C without building any renderer. The AC needs the derived
`--board` back-link filter from task 10, so this task follows it. The supporting Dashboardius
AC's `PUT /datatug/boards/*` round trip needs the boards API that Track C1 delivers; this task
proves the contract half (Go ↔ TypeScript fixture round trip) and that AC completes when C1
lands. `Board.incidentRef` as an
`IncidentRef`; the Go validator (`BoardWidget.Validate` in `datatug-core/pkg/datatug/boards.go`
today accepts only `SQL`, `HTTP`, `tabs`) and the `@datatug/board-models` `BoardWidgetName`
union gain an explicit extension boundary that preserves registered Dashboardius widget names
(`metric-series`, `comparison`, `resolution-criteria`, `replay`) losslessly; a Go ↔ TypeScript
round-trip fixture; back-links from boards to incidents through the derived index.

### Task 14: Incident-scoped read grants — request, approval, expiry

**Id:** task-14
**Verifies:** incidents#ac:access-request-approval-and-expiry-are-events
**Supporting ACs:** server-enforced-access#ac:incident-grant-expires, server-enforced-access#ac:incident-grant-needs-approval, server-enforced-access#ac:incident-grant-never-unmasks
**Depends-On:** 3, 9
**Status:** planning

The lead's simplest form, under the founder's 2026-09-11 MVP inclusion: `access.requested` / `access.approved` /
`access.declined` / `access.revoked` / `access.expired` events; the projection lists pending
requests and active grants; the server merges active grants for the calling principal and
incident into the policy set handed to `accesspolicies.Run` (no second evaluator), never
widening field allow-lists, masks, native SQL or writes; the grant id on every execution
record produced under it; expiry applied by the server at closure or at the time limit; merge revokes the duplicate's
grants; `ExecutionRequest.incident` on reads (`datatug query run --incident`);
approver = project owner or policy-named incident approver (lead assumption); endpoints
`POST /datatug/incidents/{id}/access/*`, verbs `datatug incident access request|approve|
decline|revoke|list`, and the request/approve controls on the incident page; negative tests
for unapproved, expired and out-of-incident reads.

## Deferred AC Coverage

None: every acceptance criterion of `incidents` is verified by a task above.

Sibling-Feature work this plan deliberately leaves to other plans (not acceptance of this
plan): the compare UI (NEXT), the Dashboardius app cutover into the workspace
(`product-profiles` AC dashboardius-app-in-workspace-cutover-complete, Track C3), Dashboardius
`incident-dashboards` (its own plan in `sneat-co/dashboardius`), and the Track B2/B3 incident
capabilities (hosted multi-user incidents, AI narratives, replay UI, credits, capture channels,
back-link UI badges) after Phases 4 and 5.

## Risks

- **Two context models by accident.** Overlays and roles are fields on the existing `Fact`
  contract and layers in the existing service; reviewers must reject any incident-side copy.
- **Evidence becomes a permission backdoor.** Snapshot reads, event reads and the projection
  re-check current policy; the negative tests in tasks 2 and 3 are acceptance, not optional.
- **The project-write path arrives late.** Tasks 5 and 7 depend on Phase 2 task 2; if it slips,
  the journey can be walked with fixture checks and annotations but the plan is not complete.
- **A knowledge repository turns into a log.** The default store is the project's own
  repository; the routing to a dedicated or application repository must be exercised by task 1's
  tests so the default is a choice, not a lock-in.
- **Event-stream write concurrency.** Appends are serialised by the server per incident; the
  MVP is single-server. Multi-writer/hosted semantics are Track B2.
- **Demo story mismatch.** Chinook may not carry a convincing detective story; task 11 starts
  with a written story reviewed by the founder before fixture authoring.
- **Scope creep from the vision.** Anything not in a task above stays in its owning Feature's
  NEXT/FUTURE rows; the plan's completion claim covers only the ACs listed.

## Open Questions

- Founder: confirm the Chinook-based canonical story direction before task 11 (owning
  Feature `incidents`).
- Founder: separate apps or one configurable app in the workspace (owning Feature
  `product-profiles`); task 9 assumes one configurable app for DataTug/Incidentius.
- Founder: the store routing key name and the hosted default (owning Feature `incidents`).

---
*This document follows the https://specscore.md/plan-specification*
