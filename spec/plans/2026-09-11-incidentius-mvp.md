---
format: https://specscore.md/plan-specification
status: Draft
---
# Plan: Incidentius MVP inside DataTug

**Status:** Draft
**Source Feature:** incidents
**Related Features:** evidence-records, checks, annotations, compare, product-profiles, investigation-context, dashboards
**Prerequisite Plans:** 2026-09-09-phase-1-core-investigation-loop
**Date:** 2026-09-11
**Owner:** alex
**Supersedes:** —

## Summary

Deliver roadmap Track B1: incidents as a DataTug capability, inside the existing application
and CLI, on a local or Git project under the single `serve` principal. The plan proves the
thesis with the smallest coherent slice (vision §59): the incident object and its event
stream, attachment to Investigation Context with cohorts and hypothesis overlays, evidence by
reference to recorded executions, checks, metrics with explicit criteria, recurrence and
search, API/CLI/`watch`, a deterministic status view, and the canonical demo incident with its
boring recurrence, walked end to end by a real e2e. It deliberately excludes hosted multi-user
collaboration, Dashboardius boards, AI narratives, capture channels and privileged production
access; those are Track B2+ and stay named in their owning Features.

Dependencies: Phase 1 is accepted (2026-09-11). Task 5's "promote query → check" needs the
authorized project-write path from Phase 2 task 2; every other task can start now. No task may
add a second context model, executor, policy evaluator, renderer or storage path.

## Journey

Written in the user's own sentences; each stage names its observable good result. Nobody
performs hidden setup between stages.

**Start, null action.** "I run `datatug demo`, the browser opens the demo project, and I see an
*Incidents* entry in the project menu; under the Incidentius profile it is the home page."
Good result: the incident list renders (empty or with the fixture) with no error and no
hidden URL.

**Houston.** "Support says Canadian customers cannot see invoice totals. I click *Houston,
we've got a problem*, type a title and paste the message with customer 5 and invoice 42." Good
result: INC-1 exists as `incidents/INC-1/`, the reporter is a participant, the canonical
context holds the two facts I typed, and `datatug incident show INC-1` prints the same.

**Cohorts.** "I open the Customer table from the incident, select Canadian customers, mark
them *affected*, and mark three Brazilian customers *healthy control*." Good result: the
canonical layer shows both roles; a query opened from the incident binds from the affected
facts and lists the bindings in its results header (no hidden filters).

**Checks first.** "The incident page already lists the health checks for Invoice and Customer;
I run them." Good result: each run is recorded as an execution record and a `check.run` event;
one check fails ("invoice totals present") and shows its predicate and rows count.

**Hypotheses.** "John proposes 'the FX-rate provider is down' with the five Canadian invoices as
its scope; Jack proposes 'the 14:17 deployment changed rounding'." Good result: two hypotheses
in `proposed`, each with an overlay layer, each visible in the timeline with its actor.

**Evidence.** "Jack runs the currency-rate query for CAD, records it, and attaches the record
to John's hypothesis as weakening evidence: the live rate came back fine. He annotates
`Invoice.ID = 42`: 'total missing, rate present'." Good result: the execution record carries
provenance and limitations, the annotation resurfaces when invoice 42 appears in any other
result, and John's hypothesis moves to `weakened` with the reason and record linked.

**Compare.** "Maria runs the same invoice query for the affected set and the control set." Good
result: the compare result shows the affected side missing totals and the control side complete,
each side under its own limitations, attached as evidence.

**Metrics and criteria.** "We add the metric *stuck invoices* (impact) and the criterion *stuck
invoices == 0 for 30 minutes*." Good result: the series starts with the first recorded point;
the criterion reads unmet.

**Mitigation and recovery, nobody touches the incident for a while.** "Anna applies the fix;
the metric's check is run again twice." Good result: the series reads 37 → 19 → 0, the
criterion flips to satisfied on the recorded evidence, status moves through `mitigating`,
`recovering`, `resolved` with outcome `resolved`, then `watching` and `closed`.

**Record.** "A newcomer opens the Resolution Record." Good result: timeline, hypotheses with
their evidence including the rejected one, metrics, criteria, outcome, and every evidence item
opens under the newcomer's own policy; `show --at 14:42` shows what was believed then.

**Epilogue A, three months later.** "The same symptom appears; I create INC-2." Good result:
`similar` names INC-1 with the matching signals; I link `recurrence-of`; the previously created
checks run; the recurrence resolves in minutes with a recorded TTRI.

**Epilogue B, an agent.** "A script runs `datatug incident watch --json --since <cursor>` while
all of the above happens." Good result: it receives every event exactly once as NDJSON and can
resume from the printed cursor.

## Approach

- Reconcile before implementing: refresh `origin/main` of `datatug-core`, `datatug-cli`,
  `datatug-apps`, `datatug-demo-projects`; read the transport appendix and the sibling Features
  (`evidence-records`, `checks`, `annotations`, `compare`, `product-profiles`) at their
  current revisions. Component presence is not journey acceptance.
- Model first, in `datatug-core`: the incident asset, event schema, projection and fold; the
  execution record. Server second, in `datatug-cli`: endpoints, recording, checks runner,
  compare, incident stream, CLI verbs. Browser third, in `datatug-apps`: incident pages,
  context roles/overlays, evidence and metrics tables, the `incidentius` profile. Demo last:
  fixtures and the whole-journey e2e that gates the release candidate.
- Contract amendments (optional `Fact.role`, `Fact.layer`; `record`/`snapshot` flags on
  execution; the incident endpoints) are frozen as fixtures in `datatug-core/pkg/apicontract`
  before server or client work, as Phase 1 task 12 did.
- Independent tasks run in separate WB worktrees with one landing owner per repository; each
  repository lands one candidate per task, never independent declarations of a finished journey.
- Track status with `specscore task change-status --plan 2026-09-11-incidentius-mvp`; a task
  completes only with implementation evidence and passing ACs.

## Tasks

### Task 1: Incident model, event schema, projection

**Id:** task-1
**Verifies:** incidents#ac:event-stream-append-only-and-replayable, incidents#ac:show-at-reconstructs-belief, incidents#ac:inference-requires-refs, incidents#ac:outcomes-and-duplicate-merge
**Depends-On:** —
**Status:** planning

In `datatug-core`: `Incident`, `IncidentEvent` with actor/assertion/refs, the typed payloads,
the deterministic fold to the projection (with `at <= T`), lifecycle/outcome validation, merge
semantics, `incidents/<id>/events.jsonl` + `incident.json` in the filestore, and the
`apicontract` fixtures for `/datatug/incidents/*`. 100% covered; byte-identical fold test.

### Task 2: Execution records and opt-in snapshots

**Id:** task-2
**Verifies:** incidents#ac:hypothesis-rejected-keeps-evidence
**Supporting ACs:** evidence-records#ac:record-created-on-request, evidence-records#ac:fingerprint-without-snapshot, evidence-records#ac:snapshot-read-rechecks-policy, evidence-records#ac:retention-expires-snapshot-keeps-record, evidence-records#ac:series-from-repeated-runs
**Depends-On:** —
**Status:** planning

In `datatug-core` + `datatug-cli`: `ExecutionRecord` with canonical fingerprint, `record`/
`snapshot` flags on `exec/run_query` and related-rows, server-side store with retention and
masking policy, snapshot reads re-checked through `accesspolicies.Run`, `GET /datatug/executions*`,
`datatug execution list|show`. Negative tests: a principal without current access to the
recorded columns gets a non-disclosing denial.

### Task 3: Incident endpoints, CLI verbs and watch stream

**Id:** task-3
**Verifies:** incidents#ac:houston-creates-incident-with-context, incidents#ac:watch-streams-events-with-cursor, incidents#ac:cli-json-matches-api
**Depends-On:** 1
**Status:** planning

In `datatug-cli`: list/create/show(`?at=`)/append/events(NDJSON, cursor)/search/similar/merge
under the serve principal with project-write capability checks; `datatug incident` verbs per
`cli/incident`; `watch` renders one human line per event or raw NDJSON. Contract test against
the real server; no polling loop in the CLI.

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
**Supporting ACs:** checks#ac:check-run-records-result, checks#ac:health-checks-listed-for-entity, checks#ac:legacy-check-types-still-validate, checks#ac:promote-query-creates-check-file
**Depends-On:** 2
**Status:** planning

In `datatug-core` + `datatug-cli`: check definitions referencing library queries with `expect`
predicates, purposes and scope; the runner records an execution with a deterministic result;
`datatug check list|show|run`; listing by purpose and scope for the incident page. "Promote
query → check" lands only after Phase 2 task 2's write path exists; the rest does not wait.

### Task 6: Metrics series and resolution criteria

**Id:** task-6
**Verifies:** incidents#ac:metric-series-and-criteria
**Depends-On:** 1, 2, 5
**Status:** planning

Metric definitions (check or query + numeric projection + category), `metric.observed` points
from recorded runs, criteria as predicates evaluated from evidence, projection fields, and the
numeric series/criteria tables in the incident page. No chart code in DataTug (dashboards
REQ:no-charts-in-datatug-web-ui).

### Task 7: Annotations on values and execution records

**Id:** task-7
**Verifies:** incidents#ac:hypothesis-rejected-keeps-evidence
**Supporting ACs:** annotations#ac:annotate-value-resurfaces-in-other-query, annotations#ac:annotation-carries-provenance, annotations#ac:annotation-hidden-when-target-denied, annotations#ac:annotation-on-execution-record
**Depends-On:** 2
**Status:** planning

`annotations/` project asset written through the server, `POST/GET /datatug/annotations`,
`datatug annotation add|list`, grid marker and context-panel resurfacing in `datatug-apps`,
visibility filtered by the target's current policy.

### Task 8: Compare API and CLI (two scopes)

**Id:** task-8
**Verifies:** incidents#ac:demo-incident-journey-e2e
**Supporting ACs:** compare#ac:env-diff-shows-changed-rows, compare#ac:cohort-distribution-highlights-overrepresentation, compare#ac:compare-respects-limitations
**Depends-On:** 2
**Status:** planning

`POST /datatug/compare` and `datatug compare` for environment-vs-environment and fact-set-vs-
fact-set (cohorts), keyed recordset diff and categorical distribution, each side its own
recorded execution under its own policy; attachable as evidence. Compare UI is NEXT.

### Task 9: Incident pages and the incidentius profile

**Id:** task-9
**Verifies:** incidents#ac:status-projection-carries-provenance, incidents#ac:participant-role-separate-from-permission, incidents#ac:incidentius-profile-home
**Supporting ACs:** product-profiles#ac:hostname-selects-home, product-profiles#ac:deep-link-works-in-other-profile, product-profiles#ac:one-build-serves-both-profiles
**Depends-On:** 3, 4, 6
**Status:** planning

In `datatug-apps`: incident list/detail/create pages, timeline, hypotheses and evidence
panels, participants with roles, status view with refs, Resolution Record page; the profile
table and hostname/config selection; menu entry under the `datatug` profile. Ionic frame; no
chart islands in this task.

### Task 10: Recurrence, related links, search and similarity

**Id:** task-10
**Verifies:** incidents#ac:recurrence-link-records-ttri, incidents#ac:search-finds-similar-by-entities
**Depends-On:** 3
**Status:** planning

Typed links as events, TTRI derivation, text + fact search, deterministic explainable
similarity over entities/fields/values/queries/checks/failed checks/environments/outcome, in
API, CLI and the incident page.

### Task 11: Canonical demo incident fixtures

**Id:** task-11
**Verifies:** incidents#ac:demo-incident-journey-e2e
**Depends-On:** 1, 5
**Status:** planning

In `datatug-demo-projects`: author the story (open question in the Feature: Chinook first),
the fixture `incidents/INC-1/events.jsonl` with several actors, useful wrong theories, red
herrings, cohorts, an environment comparison, the FX-rate provider as external counterparty,
recorded executions and snapshots, annotations, checks, metrics, mitigation, recovery,
resolution; the recurrence fixture `INC-2`; validator coverage for the incident asset layout.

### Task 12: Whole-journey e2e and release candidate

**Id:** task-12
**Verifies:** incidents#ac:demo-incident-journey-e2e, incidents#ac:recurrence-epilogue-boring, incidents#ac:houston-creates-incident-with-context, incidents#ac:investigation-becomes-incident, incidents#ac:watch-streams-events-with-cursor
**Depends-On:** 3, 4, 6, 7, 8, 9, 10, 11
**Status:** planning

Playwright + CLI journey from `datatug demo` through the Resolution Record and both epilogues
against the real server with no intercepted responses, asserting events, records and cursors;
required CI job on `datatug-apps` and `datatug-cli`; pinned revision set (core, cli, apps,
demo, providers); a named datatug-cli release and datatug.app deployment receipt.

## Deferred AC Coverage

None: every acceptance criterion of `incidents` is verified by a task above.

Sibling-Feature work this plan deliberately leaves to other plans (not acceptance of this
plan): the compare UI (`compare` AC snapshot-vs-live-shows-recovery, NEXT), the Dashboardius
app cutover into the workspace (`product-profiles` AC dashboardius-app-in-workspace-cutover-complete,
Track C3), Dashboardius `incident-dashboards` (its own plan in `sneat-co/dashboardius`), and
the Track B2/B3 incident capabilities (hosted multi-user incidents, AI narratives, replay UI,
credits, capture channels, incident-scoped grants) after Phases 4 and 5. Task 5's promote path
additionally waits for `2026-09-09-phase-2-library-depth-and-provenance` task 2.

## Risks

- **Two context models by accident.** Overlays and roles are fields on the existing `Fact`
  contract and layers in the existing service; reviewers must reject any incident-side copy.
- **Evidence becomes a permission backdoor.** Snapshot reads re-check current policy; negative
  tests are part of task 2's acceptance, not optional.
- **Event-stream write concurrency.** `events.jsonl` appends are serialised by the server per
  incident; the MVP is single-server. Multi-writer/hosted semantics are Track B2.
- **Demo story mismatch.** Chinook may not carry a convincing detective story; task 11 starts
  with a written story reviewed by the founder before fixture authoring.
- **Scope creep from the vision.** Anything not in a task above stays in its owning Feature's
  NEXT/FUTURE rows; the plan's completion claim covers only the ACs listed.

## Open Questions

- Founder: confirm the Chinook-based canonical story direction before task 11 (owning
  Feature `incidents`).
- Founder: separate apps or one configurable app in the workspace (owning Feature
  `product-profiles`); task 9 assumes one configurable app for DataTug/Incidentius.

---
*This document follows the https://specscore.md/plan-specification*
