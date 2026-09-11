---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Incidents (Incidentius core)

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=request-change) |
**Status:** Draft
**Source Ideas:** [Incidentius](https://github.com/sneat-co/backstage/blob/main/spec/ideas/incidentius.md)
**Kind:** Capability (implementations in `datatug-core`, `datatug-cli`, `datatug-apps`, `datatug-demo-projects`)
**Family:** [product-family](../../product/product-family.md) · P16 in the [hierarchy](../../product/hierarchy.md)
**Plan:** [2026-09-11-incidentius-mvp](../../plans/2026-09-11-incidentius-mvp.md)
**Vision:** [incident-resolution vision](../../research/2026-09-11-incident-resolution-vision.md) §1–§5, §11, §17–§25, §30–§31, §38–§39, §41–§43, §50–§54

## Summary

An incident is a DataTug project asset that answers, while people work, what is happening, what have we learned, what should we investigate next, and how do we resolve it; Incidentius is an incident-resolution workspace built inside DataTug as a product profile, with an append-only event stream, hypotheses, evidence by reference, metrics, recurrence links and a Resolution Record.

> **Incidentius is an incident-resolution workspace.**

It starts when somebody says **"Houston, we've got a problem"** (vision §1) and it attaches to a DataTug
Investigation Context rather than reimplementing one (vision §5). Everything that happens is an
append-only event; the current state is a projection; replay reconstructs what the team believed
at any moment. Hypotheses, evidence, cohorts, checks, metrics, recurrence links, status views and
the Resolution Record are all views of that one stream. The MVP lives inside the DataTug
application; Incidentius is a product profile, not a second application (vision §4).

Not an incident-management platform: no paging, on-call, tickets, status pages, SLA machinery
(vision §3.1). Those systems answer *who should respond*; this Feature answers *what is happening
and how we resolve it*.

```text
Incidentius Incident
│
├── identity / lifecycle / outcome
├── participants / incident roles
├── attached Investigation Context   ← owned by DataTug (P6); roles + overlays added there
├── hypotheses (proposed … confirmed/rejected) ── evidence links
├── evidence: execution records · snapshots · annotations · check runs   ← DataTug artefacts
├── metrics and resolution criteria
├── recurrence / related links
└── event stream → projection → timeline · watch · status views · Resolution Record
```

## Problem

During an incident the team's shared memory is chat scrollback. Nobody can answer *who already
checked this, with which query, against which records, and what it showed at 14:42*; wrong
theories vanish; the same investigation costs the same effort on recurrence; and the
postmortem is rewritten from scratch afterwards. DataTug already holds the pieces: typed
facts with scope and conditions, library queries with semantic parameters, server-enforced
execution with provenance and limitations, Git-native project assets and grants. What is
missing is the incident object that ties them into resolution semantics and a durable record.

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Incident asset, event stream, projection, list/detail, "Houston" creation | new | MVP | strategic |
| Lifecycle, five outcomes, duplicate merge | new | MVP | strategic |
| Attached Investigation Context with cohort roles and hypothesis overlays | new on P6 | MVP | strategic |
| Hypotheses with states and evidence links | new | MVP | strategic |
| Evidence via execution records, snapshots, annotations, check runs | new on P17–P19 | MVP | strategic |
| Checks listed and run at incident start | new on P19 | MVP (manual) · NEXT (auto-run safe checks) | strategic |
| Metrics as check series; explicit resolution criteria | new | MVP (numeric table) · NEXT (Dashboardius boards) | strategic |
| Recurrence and related links; search; deterministic similarity; TTRI | new | MVP · NEXT (ranking) | strategic |
| Troubleshooter status projection | new | MVP (deterministic) · NEXT (AI narratives per audience) | strategic |
| API + `datatug incident` + `watch` event stream | new | MVP | enabler |
| Resolution Record page; timeline | new | MVP · NEXT (replay UI, RCA/postmortem generation) | strategic |
| Canonical demo incident with post-credits recurrence | new | MVP | acceptance |
| Capture channels: pasted text extraction, Slack/Teams | new | NEXT (deterministic extraction, read-only Slack) · FUTURE (voice, commands) | extension |
| Appreciation credits, deployment correlation, SpecScore/Codegrapher links | new | NEXT / FUTURE | extension |
| AI Coordinator / Troubleshooter, self-learning skills, generated videos, training replay | new | FUTURE | extension |

### Identity and storage

#### REQ: incident-is-project-asset

An incident MUST be stored as a project asset under `incidents/<id>/` with two files:
`events.jsonl` (append-only, one event per line, monotonic `seq`) and `incident.json` (the
projection, regenerable from the events). It is written only through the server. Access
follows the project's grants and server-enforced access control, with no per-incident owner
field, exactly as boards do (founder, 2026-09-09: *"board belongs to DataTug project. Project
can be shared across Sneat spaces"*). Lead assumption: the human id is `INC-<n>` from a
project-scoped sequence assigned by the server, and a `uid` (UUID) exists for cross-project
links.

#### REQ: houston-creation

A user MUST be able to create an incident from a title and free text ("Houston, we've got a
problem"), from a running investigation ("make this an incident", which attaches the current
Investigation Context as the canonical layer without copying or losing it — vision §5), and
from the CLI/API. Creation records the reporter as a participant and emits `incident.created`.

### Event stream

#### REQ: event-stream-append-only

Every change to an incident MUST be an event appended by the server. Events are never edited
or deleted; corrections are new events. Lead assumption for the shape (an engineering decision
for the approved scope, not a founder ruling):

```text
IncidentEvent {
  id, seq, at (RFC3339 UTC), incident,
  actor   {kind: human|agent|system, id, via?: web|cli|api|slack},
  type    (see list),
  assertion {kind: observation|claim|hypothesis|inference|deterministic-result,
             confidence?: speculative|likely|confirmed},
  refs    [{kind: execution|snapshot|annotation|check|query|fact|hypothesis|event|incident|board, id}],
  payload (per type)
}
types: incident.created · incident.status · incident.outcome · incident.merged ·
       participant.joined · participant.role ·
       context.layer.created · context.fact.added · context.fact.removed ·
       context.fact.promoted · context.fact.rejected ·
       hypothesis.proposed · hypothesis.state ·
       evidence.attached · execution.recorded · check.run · annotation.added ·
       metric.observed · criterion.added · criterion.evaluated ·
       note.added · message.captured · link.added · deployment.noted · credit.given
```

`assertion.kind` separates fact from belief (vision §12, §31, §57 Q12): a deterministic check
result or recorded execution is `observation`/`deterministic-result`; a human statement is a
`claim`; a hypothesis is `hypothesis`; anything produced by AI or a heuristic is `inference`
and MUST carry `refs` to the events it was derived from. The server MUST reject an
`inference` event without refs and an `observation` event from an `agent` actor that carries
no execution/check reference.

#### REQ: projection-and-replay

The server MUST derive `incident.json` by a deterministic fold over the events and MUST be able
to produce the projection at any time `T` by folding only events with `at <= T`
(`GET …?at=`, `datatug incident show <id> --at <ts>`). The projection includes identity,
status, outcome, participants, the context layers, hypotheses with their evidence, evidence
list, check runs, metrics series, criteria, links, TTRI timestamps and the status view. Two
folds over the same events MUST produce byte-identical JSON.

#### REQ: watch-event-cursor

The server MUST expose the event stream for one incident and for the whole project
(`GET /datatug/incidents/{id}/events?since=<cursor>`, `GET /datatug/incidents/events?since=`)
as NDJSON with a resumable opaque cursor. `datatug incident watch [<id>] [--json] [--since
<cursor>]` renders one line per event for humans (vision §51 example) and raw events with
`--json`. The same stream powers the timeline, status views, AI coordination, connectors and
replay; there is no second notification mechanism.

### Lifecycle

#### REQ: lifecycle-and-outcomes

Status MUST be one of `open → investigating → mitigating → recovering → resolved → watching →
closed`; only `open`, `resolved` and `closed` are mandatory stops, the others are useful
semantics, never bureaucracy (vision §22–§23). Resolving MUST record exactly one outcome:
`resolved`, `false-alarm`, `accepted` (known issue), `handed-off`, or `unresolved`
(investigation stopped). Duplicate is not an outcome: `POST …/merge` MUST append the
duplicate's events into the surviving incident with `incident.merged` refs and close the
duplicate with `mergedInto` set. A recurrence stays a separate incident linked to the earlier
one (REQ:recurrence-and-links).

### Investigation Context

#### REQ: attached-investigation-context

An incident MUST reference its canonical context and any overlays through the
`investigation-context` Feature's layered model (REQ:cohort-roles,
REQ:context-layers-and-promotion there). Facts keep the transport shape of
[api-contract.md](../core-investigation-loop/api-contract.md) with the optional `role` and
`layer` fields. Adding, promoting and rejecting facts are incident events; promotion moves a
fact from a hypothesis layer to the canonical layer without deleting its history; rejection
keeps it attached to the hypothesis (vision §11.1). An incident MUST NOT hold a second copy of
context semantics, binding rules or conflict handling.

### Hypotheses and evidence

#### REQ: hypotheses

A hypothesis MUST have a title, a state (`proposed | investigating | supported | weakened |
rejected | confirmed`), an optional scope (facts, typically a cohort), the queries/checks used,
and evidence links with a reason per link. State changes are events; a rejected hypothesis and
its evidence remain visible forever (vision §3.3: *negative evidence is positive progress*). A
later reader MUST be able to see what was checked, against which records, with which queries,
what was observed and why the theory was rejected (vision §17).

#### REQ: evidence-links

Evidence MUST be a reference to a DataTug artefact: an execution record, a result snapshot, an
annotation or a check run ([evidence-records](../evidence-records/README.md),
[annotations](../annotations/README.md), [checks](../checks/README.md)). Attaching evidence is
an `evidence.attached` event with refs; the incident MUST NOT store result rows itself. An
execution run from within an incident MUST be recorded (`record:true`) and tagged with the
incident id in its provenance.

#### REQ: checks-at-incident-start

When an incident's canonical context names entities or sources, the incident page and
`datatug incident show` MUST list the applicable health and diagnostic checks
(`checks` REQ:checks-listed-by-purpose-and-scope) with their latest known result, and let an
authorized participant run them; each run is a `check.run` event. Automatic execution of safe
checks on creation is NEXT and MUST remain opt-in per project.

### Metrics and resolution

#### REQ: metrics-and-resolution-criteria

An incident metric MUST be a check or library query with a numeric projection
(`evidence-records` REQ:series-from-records) placed in a category `impact | diagnostic |
recovery | watch` (vision §20). Each recorded run appends a `metric.observed` point; the
projection exposes the series. Resolution criteria MUST be explicit predicates over metrics or
check results ("no new stuck orders for 30 minutes", "backlog < 10"), each evaluated from
recorded evidence (`criterion.evaluated`), never from a free-text claim; `resolved` MAY be
recorded while criteria are unmet, but the projection MUST show which are unmet. The MVP
renders series as numeric tables in DataTug; boards and charts are Dashboardius's
[incident-dashboards](https://github.com/sneat-co/dashboardius/blob/main/spec/features/incident-dashboards/README.md).
AI turning a sentence into a proposed criterion is NEXT.

### Recurrence and discovery

#### REQ: recurrence-and-links

An incident MUST support typed links `recurrence-of` and `related` to other incidents, with the
link as an event. TTRI (time to recurrence identification, vision §24) MUST be derivable from
events: `suggestedAt` (first `inference` naming a candidate recurrence) and `confirmedAt`
(the human `link.added` with `recurrence-of`).

#### REQ: search-and-similarity

`GET /datatug/incidents/search` and `datatug incident search` MUST find incidents by text and
by facts (`Customer.ID=5`, `entity=Invoice`). `…/{id}/similar` MUST rank other incidents by a
deterministic, explainable overlap of signals — entities, fields, context values, query and
check ids, failed checks, environments, outcome — and MUST show which signals matched
(vision §25). Semantic/embedding ranking and knowledge decay are NEXT/FUTURE and MUST NOT
replace the explainable score.

### Projections and people

#### REQ: status-projection

The projection MUST include a troubleshooter status view generated deterministically from state
(counts of affected facts, hypotheses by state with their latest evidence, unowned questions,
latest metric points) where every sentence carries refs to the events it summarises
(vision §31). AI-written narratives for troubleshooters, leadership and support are NEXT; they
MUST be `inference` events with refs and MUST never restate a hypothesis as a fact.

#### REQ: participants-and-roles

Persona, incident role and permission MUST stay separate (vision §42–§43). Incident roles are
`reporter | investigator | coordinator | observer`, recorded as events; they confer no data
access. Permission comes from the project's grants and policies; incident-scoped temporary
grants are NEXT in `server-enforced-access`. Business users and external counterparties appear
as participants of kind `external` with no access.

#### REQ: resolution-record

Every incident MUST have a Resolution Record page (web) and document (`incident.json` plus the
rendered timeline) that shows the timeline, participants, hypotheses with evidence, metrics,
criteria, outcome and links, with every evidence item opening its execution record or snapshot
under the reader's current policy (vision §38). Generated RCA/postmortem/executive/support
documents, appreciation credits (vision §41) and replay UI (vision §39) are NEXT projections
over the same record.

### Surfaces

#### REQ: api-and-cli

Every capability above MUST be API-first under `/datatug/incidents/*` (list, create, show
with `?at=`, append event, events stream, search, similar, merge) using the transport
appendix's envelopes and error codes, and exposed through `datatug incident` verbs specified
in `datatug-cli` `cli/incident` (vision §50: *humans, scripts and AI agents operate on the
same model*). Agents are principals through the same server; there is no agent-only path.

#### REQ: profile-home

Under the `incidentius` product profile the home route MUST be the incident list with the
"Houston" entry point; under the `datatug` profile incidents are reachable from the project
menu and any investigation ([product-profiles](../product-profiles/README.md)).

### Canonical demo

#### REQ: canonical-demo-incident

`datatug-demo-projects/demo-project-1` MUST ship a fixture incident whose events involve
several human actors, useful wrong theories, red herrings, affected and control cohorts, an
environment comparison, an external counterparty, recorded executions and snapshots,
annotations, checks, metrics, mitigation, verified recovery and a resolution — and a second,
"post-credits" fixture incident that is a recurrence resolved in minutes (vision §53–§54). The
same fixtures MUST serve the whole-journey e2e, the interactive demo, onboarding and the
marketing story. The story itself is not locked by this Feature (open question below).

## Not Doing / Out of Scope

- Paging, on-call, escalation, ticketing, status pages, SLA workflow (integrate, never build).
- A second context model, query engine, policy evaluator, renderer or storage path.
- Storing result rows inside the incident; evidence is by reference.
- Live multi-user editing under authenticated principals in the MVP (needs Phase 4 grants and an
  authenticated server profile — Track B2).
- Voice capture, Slack/Teams commands, self-learning extractors, AI Coordinator, generated
  videos, training mode, controlled production writes (FUTURE tiers; see the reconciliation).

## Acceptance Criteria

### AC: houston-creates-incident-with-context (verifies REQ:houston-creation, REQ:incident-is-project-asset)

**Given** the demo project served with `--as admin`
**When** the user creates an incident titled "Canadian invoices stuck" with the text "customer 5 and invoice 42 affected"
**Then** `incidents/INC-1/events.jsonl` holds an `incident.created` event, the reporter is a participant, and the canonical context layer exists.

### AC: investigation-becomes-incident (verifies REQ:houston-creation, REQ:attached-investigation-context)

**Given** an Investigation Context holding `Customer.ID = 5` and `Customer.Country = Canada`
**When** the user chooses "make this an incident"
**Then** the new incident's canonical layer references those facts with their origins intact, and the original context is unchanged.

### AC: event-stream-append-only-and-replayable (verifies REQ:event-stream-append-only, REQ:projection-and-replay)

**Given** the fixture incident with 40 events
**When** the projection is folded twice, and a request tries to edit event 12
**Then** both folds are byte-identical, and the edit is refused with `INVALID_REQUEST`.

### AC: show-at-reconstructs-belief (verifies REQ:projection-and-replay)

**Given** the fixture incident in which hypothesis H1 is rejected at 14:43
**When** `datatug incident show INC-1 --at 14:42:00Z` runs
**Then** H1 is `investigating`, its later rejection and evidence are absent, and the status view matches the 14:42 state.

### AC: inference-requires-refs (verifies REQ:event-stream-append-only)

**Given** an agent principal
**When** it appends an `inference` event without refs, or an `observation` without an execution reference
**Then** both are refused with a structured error and nothing is appended.

### AC: watch-streams-events-with-cursor (verifies REQ:watch-event-cursor)

**Given** `datatug incident watch INC-1 --json` running
**When** another actor appends two events and the watcher restarts with `--since <cursor>`
**Then** exactly the two events were streamed as NDJSON before the restart and nothing is duplicated after it.

### AC: outcomes-and-duplicate-merge (verifies REQ:lifecycle-and-outcomes)

**Given** INC-3 is a duplicate of INC-1
**When** INC-3 is merged into INC-1
**Then** INC-1 gains INC-3's events with `incident.merged` refs, INC-3 is closed with `mergedInto: INC-1` and no outcome, and resolving INC-1 later requires exactly one of the five outcomes.

### AC: overlay-promotion-updates-canonical (verifies REQ:attached-investigation-context)

**Given** hypothesis H17 with an overlay layer holding `Customer.ID = 11..15` with role `suspected`
**When** H17 is confirmed and its facts promoted
**Then** the canonical layer contains them with role `affected`, the overlay still records the original facts and the promotion event, and an open query bound from context shows a visible re-binding choice rather than silently changing.

### AC: hypothesis-rejected-keeps-evidence (verifies REQ:hypotheses, REQ:evidence-links)

**Given** H12 "retry malfunction" with two execution records attached
**When** H12 is set to `rejected` with reason "retries completed for affected orders"
**Then** the projection still lists H12, both records, the reason and the queries used, and each record opens under the reader's current policy.

### AC: checks-run-at-start-recorded (verifies REQ:checks-at-incident-start)

**Given** the canonical context names `Invoice` and the demo's health checks
**When** the incident page lists checks and the user runs "invoice totals present"
**Then** the list came from the checks API filtered by purpose and scope, and the run appended a `check.run` event referencing a new execution record.

### AC: metric-series-and-criteria (verifies REQ:metrics-and-resolution-criteria)

**Given** an impact metric "stuck invoices" and a criterion "stuck invoices == 0"
**When** the metric's check is run three times returning 37, 19 and 0
**Then** the projection shows the series with three `metric.observed` points and the criterion evaluated `satisfied` only after the third run.

### AC: recurrence-link-records-ttri (verifies REQ:recurrence-and-links)

**Given** INC-2 created at 09:17 and an `inference` event at 09:17:21 naming INC-1 as a candidate recurrence
**When** a human links INC-2 `recurrence-of` INC-1 at 09:18
**Then** the projection reports `ttri.suggestedAt = 09:17:21` and `ttri.confirmedAt = 09:18`.

### AC: search-finds-similar-by-entities (verifies REQ:search-and-similarity)

**Given** INC-1 resolved and INC-2 open with overlapping `Invoice` facts, the same failed check and the same environment
**When** `datatug incident similar INC-2` runs
**Then** INC-1 ranks first with the matched signals listed, and an unrelated incident with only a shared word in its title ranks below it.

### AC: status-projection-carries-provenance (verifies REQ:status-projection)

**Given** the fixture incident at its peak
**When** the status view renders
**Then** every sentence carries refs resolving to existing events and no hypothesis is worded as a fact.

### AC: participant-role-separate-from-permission (verifies REQ:participants-and-roles)

**Given** a `support` principal added as `coordinator`
**When** they open an execution record containing `Customer.Email`
**Then** the column is absent under their policy regardless of the role, matching `core-investigation-loop` AC:hidden-column-refused.

### AC: cli-json-matches-api (verifies REQ:api-and-cli)

**Given** the incident API and `datatug incident show INC-1 --json`
**When** both are called for the same incident
**Then** the CLI output equals the API envelope byte for byte, and every CLI verb maps to one API operation.

### AC: incidentius-profile-home (verifies REQ:profile-home)

**Given** the app served under the `incidentius` profile
**When** it opens
**Then** the incident list with the "Houston, we've got a problem" entry point is the home route, and the same route is reachable under the `datatug` profile from the project menu.

### AC: demo-incident-journey-e2e (verifies REQ:canonical-demo-incident)

**Given** the demo project and a running agent
**When** the whole-journey e2e runs
**Then** it walks create → context with cohorts → checks → hypotheses → evidence → metrics → mitigation → recovery → resolved → watching → closed → Resolution Record against the real server with no intercepted responses, asserting events and records, not rendered text alone.

### AC: recurrence-epilogue-boring (verifies REQ:canonical-demo-incident, REQ:search-and-similarity)

**Given** the resolved fixture incident and the recurrence fixture
**When** the recurrence is opened
**Then** `similar` names the first incident with matching signals, its checks run and the recurrence resolves with a recorded TTRI, demonstrating *"every incident should make the next one more boring"*.

## Open Questions

- **Canonical story.** Lead recommendation: build it on `demo-project-1` (Chinook + inGitDB
  `support-notes` + the CountriesNow/Frankfurter HTTP sources): "Canadian invoices stuck",
  cohorts by `Customer.Country`, the FX-rate provider as the external counterparty, a
  deployment as the red herring. If Chinook cannot carry a detective story, a second
  orders/shipping demo project is needed. Founder to confirm the direction; the story text is
  authored during the plan's fixture task, not here.
- **Id scheme** `INC-<n>` per project plus `uid` — lead assumption.
- **Snapshot defaults for incident-recorded runs** (opt-in vs always) belong to
  [evidence-records](../evidence-records/README.md).
- **Incidentius brand and domain** — [product-profiles](../product-profiles/README.md).

---
*This document follows the https://specscore.md/feature-specification*
