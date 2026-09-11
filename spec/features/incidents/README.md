---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Incidents (Incidentius core)

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/incidents?op=request-change) |
**Status:** Approved
**Source Ideas:** [Incidentius](https://github.com/sneat-co/backstage/blob/main/spec/ideas/incidentius.md)
**Kind:** Capability (implementations in `datatug-core`, `datatug-cli`, `datatug-apps`, `datatug-demo-projects`)
**Family:** [product-family](../../product/product-family.md) · P16 in the [hierarchy](../../product/hierarchy.md)
**Plan:** [2026-09-11-incidentius-mvp](../../plans/2026-09-11-incidentius-mvp.md)
**Vision:** [incident-resolution vision](../../research/2026-09-11-incident-resolution-vision.md) §1–§5, §11, §17–§25, §30–§31, §38–§39, §41–§43, §50–§54

## Summary

An incident is a DataTug record in a configurable incident store that answers, while people work, what is happening, what have we learned, what should we investigate next, and how do we resolve it; Incidentius is an incident-resolution workspace built inside DataTug as a product profile, with an append-only event stream, hypotheses, evidence by reference, metrics, recurrence links, back-links from DataTug assets and a Resolution Record.

> **Incidentius is an incident-resolution workspace.**

It starts when somebody says **"Houston, we've got a problem"** (vision §1) and it attaches to a
DataTug Investigation Context rather than reimplementing one (vision §5). Everything that
happens is an append-only event; the current state is a projection; replay reconstructs what
the team believed at any moment. An incident may span several DataTug projects maintained by
different teams, and a project may be used by several Incidentius instances. Hypotheses,
evidence, cohorts, checks, metrics, recurrence links, status views and the Resolution Record
are all views of one stream. The MVP lives inside the DataTug application; Incidentius is a
product profile, not a second application (vision §4).

Not an incident-management platform: no paging, on-call, tickets, status pages, SLA machinery
(vision §3.1). Those systems answer *who should respond*; this Feature answers *what is happening
and how we resolve it*.

```text
Incidentius Incident  (store id + incident id)
│
├── identity / lifecycle / outcome
├── participants / incident roles
├── referenced DataTug projects (one or many)
├── attached Investigation Context   ← owned by DataTug (P6); roles + overlays added there
├── hypotheses (proposed … confirmed/rejected) ── evidence links
├── evidence: execution records · snapshots · annotations · check runs · compare runs   ← DataTug artefacts
├── metrics and resolution criteria
├── recurrence / related links
└── event stream → projection → timeline · watch · status views · Resolution Record
```

## Problem

During an incident the team's shared memory is chat scrollback. Nobody can answer *who already
checked this, with which query, against which records, and what it showed at 14:42*; wrong
theories vanish; the same investigation costs the same effort on recurrence; and the
postmortem is rewritten from scratch afterwards. DataTug already holds the pieces: typed
facts with scope, library queries with semantic parameters, server-enforced execution with
provenance and limitations, Git-native project assets and grants. What is missing is the
incident object that ties them into resolution semantics and a durable record, without
turning a curated knowledge project into an operational log.

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Incident store (configurable location), event stream, projection, list/detail, "Houston" creation | new | MVP | strategic |
| Multi-project incidents; back-links from DataTug assets by store id + incident id | new | MVP (model, API, CLI) · NEXT (UI badges) | strategic |
| Lifecycle, five outcomes, duplicate merge | new | MVP | strategic |
| Attached Investigation Context with cohort roles and hypothesis overlays | new on P6 | MVP | strategic |
| Hypotheses with states and evidence links | new | MVP | strategic |
| Evidence via execution records, snapshots, annotations, check runs, compare runs | new on P17–P20 | MVP | strategic |
| Policy re-check on every event, projection and watch read | new on P10 | MVP | security |
| Checks listed and run at incident start | new on P19 | MVP (manual) · NEXT (auto-run safe checks) | strategic |
| Metrics as check series; explicit resolution criteria | new | MVP (numeric table) · NEXT (Dashboardius boards) | strategic |
| Recurrence and related links; search; deterministic similarity; TTRI | new | MVP · NEXT (ranking) | strategic |
| Troubleshooter status projection | new | MVP (deterministic) · NEXT (AI narratives per audience) | strategic |
| API + `datatug incident` + `watch` event stream | new | MVP | enabler |
| Resolution Record page; timeline; export into the project | new | MVP (page) · NEXT (export, replay UI, RCA generation) | strategic |
| Canonical demo incident with post-credits recurrence | new | MVP | acceptance |
| Capture channels: pasted text extraction, Slack/Teams | new | NEXT (deterministic extraction, read-only Slack) · FUTURE (voice, commands) | extension |
| Appreciation credits, deployment correlation, SpecScore/Codegrapher links | new | NEXT / FUTURE | extension |
| AI Coordinator / Troubleshooter, self-learning skills, generated videos, training replay | new | FUTURE | extension |

### Identity and storage

#### REQ: incident-store

Incidents MUST live in an **incident store** owned and written only by the DataTug server,
never hand-edited. Its location is configured per served setup, in the founder's words
(2026-09-11, verbatim): *"Agree on ability of incidents to be stored in dedicated repo. But it
should be able to store in same repo as datatug or even in the application repo if users want
(can be OK for small projects)"*. The store MUST therefore be routable to (a) the DataTug
project's own repository, under `incidents/` (the default for local and small setups), (b) a
dedicated incidents repository, (c) the application's own repository, or (d) the hosted store,
space-bounded like every other Sneat record (lead assumption: `/spaces/{spaceID}/ext/incidentius/…`).
All four sit behind one store interface backed by DALgo, with inGitDB for Git-backed locations;
the MVP delivers (a)–(c), and (d) is NEXT (Track B2, hosted multi-user incidents) behind the
same interface. The layout is identical: `incidents/<id>/events.jsonl` (append-only, monotonic `seq`) and
`incidents/<id>/incident.json` (the projection, regenerable). A store has a `storeId` following
DataTug's existing store-id scheme; the routing key lives in machine-local or served-project
configuration, never as a credential in a committed file (lead assumption for the key name:
`stores.incidents`). Execution records and snapshots ([evidence-records](../evidence-records/README.md))
follow the same routing unless `stores.evidence` (lead assumption for the key name) names
another configured store. Writes to the incident store do not use, and do not wait
for, the Phase 2 project-write path for knowledge assets; checks and annotations created from an
incident are knowledge assets and do use it. Lead assumption: the human id is `INC-<n>` from a
store-scoped sequence assigned by the server, plus a `uid` (UUID).

#### REQ: multi-project-incidents

Founder ruling (2026-09-11, verbatim, answering the lead's question whether one incident
may reference several DataTug projects): *"yes it can."* An incident MAY reference one or more DataTug projects
(`projects: [{storeId, projectId, environment?}]`), because a real incident spans systems
maintained by different teams. The converse relation, *"A DataTug project can be used by
multiple incidentius instances"* (founder, 2026-09-11), is carried by REQ:incident-store and
REQ:incident-references-and-back-links. Every fact, execution record, check run and compare run carries its own project and
environment scope; because the transport appendix `Scope` is `{project, environment,
securityContextId}` with no store, the plan's model task adds an optional `Scope.storeId`
(lead assumption) so evidence from a second store stays attributable and each reader sees
each project's evidence under that project's grants. The incident itself belongs to one store
(one team or space); people from other teams join as invited participants. Entities with the
same name in two projects are not assumed identical; cross-project semantic identity is a
later capability of the semantic layer. The MVP journey exercises one project; the model MUST
NOT hard-code one.

#### REQ: incident-references-and-back-links

An incident is addressed globally by `IncidentRef = {storeId, incidentId}` (founder,
2026-09-11: *"link queries and resources of DataTug to incidents (storeID+incidentID)"*). The
canonical string form is `<storeId>/<incidentId>` (lead assumption): store ids MUST NOT
contain `/`, which the plan's model task validates. The wire spelling of the two fields is
reconciled by that task with the existing `datatug-core` `dto.ProjectRef` (JSON `storage`,
`project`), which today has no string form or parser; the logical names in every spec stay
`storeId` and `incidentId`. Back-links MUST exist in two forms: **derived** — the server answers "which
incidents touched query, check, board, execution record or entity value X" from the events of
every configured store (`GET /datatug/incidents?query=<id>|check=<id>|board=<id>`, `datatug
incident list --query <id>`); and **recorded provenance** — a knowledge asset created from an
incident (a promoted check, a captured query, an annotation, a board) carries
`origin.incident: IncidentRef` in its own file. A board's `incidentRef`
([dashboards](../dashboards/README.md) REQ:incident-boards-and-shared-renderer) is an
`IncidentRef`. UI badges ("used in 3 incidents") are NEXT; the API and CLI are MVP.

#### REQ: houston-creation

A user MUST be able to create an incident from a title and free text ("Houston, we've got a
problem"), from a running investigation ("make this an incident", which attaches the current
Investigation Context as the canonical layer without copying or losing it — vision §5), and
from the CLI/API. The free text is kept as the incident's `description`. Creation records the
reporter as a participant and emits `incident.created`.

### Event stream

#### REQ: event-stream-append-only

Every change to an incident MUST be an event appended by the server. Events are never edited
or deleted; corrections are new events. Lead assumption for the shape (an engineering decision
for the approved scope, not a founder ruling):

```text
IncidentEvent {
  id, seq, at (RFC3339 UTC), incident (IncidentRef),
  actor   {kind: human|agent|system, id, via?: web|cli|api|slack},
  type    (see list),
  assertion {kind: observation|claim|question|hypothesis|inference|deterministic-result,
             confidence?: speculative|likely|confirmed},
  refs    [{kind: execution|snapshot|annotation|check|compare|query|board|fact|hypothesis|event|incident|project, id}],
  payload (per type)
}
types: incident.created · incident.status · incident.outcome · incident.merged ·
       participant.joined · participant.role · project.linked ·
       context.layer.created · context.fact.added · context.fact.removed ·
       context.fact.promoted · context.fact.rejected ·
       hypothesis.proposed · hypothesis.state ·
       evidence.attached · execution.recorded · check.run · compare.run · annotation.added ·
       metric.observed · criterion.added · criterion.evaluated ·
       note.added · message.captured · link.added · deployment.noted · credit.given ·
       access.requested · access.approved · access.declined · access.revoked · access.expired
```

`assertion.kind` separates fact from belief (vision §12, §31, §57 Q12): a recorded execution
is an `observation`; a check or compare result is a `deterministic-result`; a human statement
is a `claim` and a human question a `question` (vision §30 "what remains unanswered"); a
hypothesis is `hypothesis`; anything produced by AI or a heuristic is `inference` and MUST
carry `refs` to the events it was derived from. The server MUST reject an `inference` event
without refs and an `observation` or `deterministic-result` from an `agent` actor that carries
no execution, check or compare reference.

#### REQ: projection-and-replay

The server MUST derive `incident.json` by a deterministic fold over the events and MUST be able
to produce the projection at any time `T` by folding only events with `at <= T`
(`GET …?at=`, `datatug incident show <id> --at <ts>`). The projection includes identity,
description, status, outcome, participants, referenced projects, the context layers,
hypotheses with their evidence, evidence list, check and compare runs, metrics series,
criteria, links, TTRI timestamps and the status view. Two folds over the same events MUST
produce byte-identical JSON.

#### REQ: events-and-projection-policy-rechecked

Every read of events, of the projection and of the watch stream MUST be evaluated for the
current principal, because event payloads carry typed fact values and human text. A fact whose
`entity.field` resolves to a column the principal may not read MUST be returned redacted
(value replaced by a redaction marker; the field name only if it is authorized metadata); an
annotation body follows [annotations](../annotations/README.md) REQ:annotation-visibility-follows-target;
evidence references resolve under the reader's policy when opened. Redaction is applied at
read time and never stored, so a later grant reveals the original. This extends
[server-enforced-access](../server-enforced-access/README.md) REQ:recorded-evidence-rechecked
to the incident stream; a project grant alone never discloses what a policy hides.

#### REQ: watch-event-cursor

The server MUST expose the event stream for one incident and for a whole store
(`GET /datatug/incidents/{id}/events?since=<cursor>`, `GET /datatug/incidents/events?since=`)
as NDJSON with a resumable opaque cursor, policy-rechecked per REQ:events-and-projection-policy-rechecked.
`datatug incident watch [<id>] [--json] [--since <cursor>]` renders one line per event for
humans (vision §51 example) on stdout, raw events with `--json`, and prints the last cursor
on stderr on exit. The same stream powers the timeline, status views, AI coordination,
connectors and replay; there is no second notification mechanism.

### Lifecycle

#### REQ: lifecycle-and-outcomes

Status MUST be one of `open → investigating → mitigating → recovering → resolved → watching →
closed`. Lead assumption: the names follow vision §22's sequence with `open` added in front,
and only `open`, `resolved` and `closed` are mandatory stops, which is the lead's reading of
vision §22 (*"These do not all need to become bureaucratic mandatory statuses"*). Resolving
MUST record exactly one outcome: `resolved`, `false-alarm`, `accepted` (known issue),
`handed-off`, or `unresolved` (investigation stopped) (vision §23). Duplicate is not an outcome:
`POST …/merge` and `datatug incident merge` MUST append the duplicate's events into the
surviving incident with `incident.merged` refs and close the duplicate with `mergedInto` set;
merging an already-merged incident is refused with the transport appendix's `INVALID_REQUEST`
(no new error code). A recurrence stays a separate incident linked to the earlier one
(REQ:recurrence-and-links).

#### REQ: access-requests-in-incident

A participant MUST be able to request, inside the incident, temporary read access to named
sources and collections of a referenced project, with a reason and an optional time limit,
and an approver — the holder of the project's incident-approver role (founder, 2026-09-11,
verbatim: *"Someone with appropriate role. By default project admin has this role but can
move it to someone else. It's part of Sneat contactius role management"*) — MUST be able to
approve, decline or revoke it there (founder, 2026-09-11: privileged-access request and
approval in its simplest form is part of the MVP). The
request, decision, use (the grant id on every execution record produced under it), expiry
and revocation are events (`access.*`), so the audit question "who could read what, when,
and who allowed it" is answered by the same stream as everything else. The projection MUST
list pending requests and active grants with their scope and expiry; expiry MUST be applied
by the server when the incident is closed or at the time limit, without a human action;
merging an incident revokes the duplicate's grants. The scope named in `access.*` events and
in the projection's grant list is redacted for a reader who may not read that source, like
any other payload (REQ:events-and-projection-policy-rechecked). What a grant may and may not
widen, and the MVP honesty about the approval step under a single served principal, are
fixed by [server-enforced-access](../server-enforced-access/README.md)
REQ:incident-scoped-grants; this Feature owns only the workflow and its events. Exposed as
`datatug incident access request|approve|decline|revoke|list` and on the incident page; a
read that wants the grant carries the incident (`ExecutionRequest.incident`).

### Investigation Context

#### REQ: attached-investigation-context

An incident MUST reference its canonical context and any overlays through the
`investigation-context` Feature's layered model (REQ:cohort-roles,
REQ:context-layers-and-promotion there). Facts keep the transport shape of
[api-contract.md](../core-investigation-loop/api-contract.md) with the optional `role` and
`layer` fields and their project scope. Adding, promoting and rejecting facts are incident
events; promotion moves a fact from a hypothesis layer to the canonical layer without deleting
its history; rejection keeps it attached to the hypothesis (vision §11.1). An incident MUST
NOT hold a second copy of context semantics, binding rules or conflict handling.

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
annotation, a check run or a compare run ([evidence-records](../evidence-records/README.md),
[annotations](../annotations/README.md), [checks](../checks/README.md),
[compare](../compare/README.md)). Attaching evidence is an `evidence.attached` event with refs;
the incident MUST NOT store result rows itself. An execution run from within an incident MUST
be recorded (`record:true`) and tagged with the `IncidentRef` in its provenance; the same tag is
how derived back-links are answered.

#### REQ: checks-at-incident-start

When an incident's canonical context names entities or sources, the incident page and
`datatug incident show` MUST list the applicable health and diagnostic checks
(`checks` REQ:checks-listed-by-purpose-and-scope) with their latest known result, and let an
authorized participant run them; each run is a `check.run` event. Automatic execution of safe
checks on creation is NEXT and MUST remain opt-in per store.

### Metrics and resolution

#### REQ: metrics-and-resolution-criteria

An incident metric MUST be a check or library query with a numeric projection
(`evidence-records` REQ:series-from-records) placed in a category `impact | diagnostic |
recovery | watch` (vision §20). Each recorded run appends a `metric.observed` point; the
projection exposes the series. Resolution criteria MUST be explicit predicates over metrics or
check results ("no new stuck orders for 30 minutes", "backlog < 10"), each evaluated from
recorded evidence (`criterion.evaluated`), never from a free-text claim; `resolved` MAY be
recorded while criteria are unmet, but the projection MUST show which are unmet. The MVP
renders series as numeric tables in DataTug (no chart code until the shared renderer lands,
[dashboards](../dashboards/README.md) REQ:no-charts-in-datatug-web-ui); boards and charts are
Dashboardius's [incident-dashboards](https://github.com/sneat-co/dashboardius/blob/main/spec/features/incident-dashboards/README.md).
AI turning a sentence into a proposed criterion is NEXT.

### Recurrence and discovery

#### REQ: recurrence-and-links

An incident MUST support typed links `recurrence-of` and `related` to other incidents,
including incidents in other stores by `IncidentRef`, with the link as an event. TTRI (time to
recurrence identification, vision §24) MUST be derivable from events: `suggestedAt` (first
`inference` naming a candidate recurrence) and `confirmedAt` (the human `link.added` with
`recurrence-of`).

#### REQ: search-and-similarity

`GET /datatug/incidents/search` and `datatug incident search` MUST find incidents across the
configured stores by text and by facts (`Customer.ID=5`, `entity=Invoice`). `…/{id}/similar`
MUST rank other incidents by a deterministic, explainable overlap of signals — entities,
fields, context values, query and check ids, failed checks, environments, referenced projects,
outcome — and MUST show which signals matched (vision §25). Semantic/embedding ranking and
knowledge decay are NEXT/FUTURE and MUST NOT replace the explainable score.

### Projections and people

#### REQ: status-projection

The projection MUST include a troubleshooter status view generated deterministically from state
(counts of affected facts, hypotheses by state with their latest evidence, open questions and
unowned work, latest metric points) where every sentence carries refs to the events it
summarises (vision §31). AI-written narratives for troubleshooters, leadership and support are
NEXT; they MUST be `inference` events with refs and MUST never restate a hypothesis as a fact.

#### REQ: participants-and-roles

Persona, incident role and permission MUST stay separate (vision §42–§43). Incident roles are
`reporter | investigator | coordinator | observer`, recorded as events; they confer no data
access. Permission comes from each referenced project's grants and policies; the incident
store's own read access follows its location (project grant when co-located, the dedicated
repository's or the space's membership otherwise). Incident-scoped temporary read grants with
request and approval are MVP in `server-enforced-access` REQ:incident-scoped-grants
(REQ:access-requests-in-incident below); approval authority comes from the project's policy,
not from an incident role. Business users and external counterparties appear as participants
of kind `external` with no access.

#### REQ: resolution-record

Every incident MUST have a Resolution Record page (web) and document (`incident.json` plus the
rendered timeline) that shows the timeline, participants, referenced projects, hypotheses with
evidence, metrics, criteria, outcome and links, with every evidence item opening its execution
record or snapshot under the reader's current policy (vision §38). Exporting the record into a
referenced project or an inGitDB dataset as a reviewable knowledge artefact, generated
RCA/postmortem/executive/support documents, appreciation credits (vision §41) and the replay UI
(vision §39) are NEXT projections over the same record.

### Surfaces

#### REQ: api-and-cli

Every capability above MUST be API-first under `/datatug/incidents/*` (list with asset
filters, create, show with `?at=`, append event, events stream, search, similar, merge, access
request/approve/decline/revoke/list) using
the transport appendix's envelopes and error codes, and exposed through `datatug incident`
verbs specified in `datatug-cli` `cli/incident` (vision §50: *humans, scripts and AI agents
operate on the same model*). Agents are principals through the same server; there is no
agent-only path.

#### REQ: profile-home

Under the `incidentius` product profile the home route MUST be the incident list with the
"Houston" entry point; under the `datatug` profile incidents are reachable from the project
menu and any investigation ([product-profiles](../product-profiles/README.md)).

### Canonical demo

#### REQ: canonical-demo-incident

`datatug-demo-projects/demo-project-1` MUST ship a fixture incident (in the project's own
`incidents/` store) whose events involve several human actors, useful wrong theories, red
herrings, affected and control cohorts, an environment comparison, an external counterparty,
recorded executions and snapshots, annotations, checks, metrics, mitigation, verified recovery
and a resolution — and a second, "post-credits" fixture incident that is a recurrence resolved
in minutes (vision §53–§54). The same fixtures MUST serve the whole-journey e2e, the
interactive demo, onboarding and the marketing story. The story's setting is ruled (founder,
2026-09-11, verbatim: *"Chinook first, our own inGitDB second, maybe a copy of Chinook or our
own use case"*; lead reading: the canonical story is built on `demo-project-1` — Chinook plus
the inGitDB `support-notes` and the HTTP sources — a second demo project on our own inGitDB
data follows, and a Chinook copy or our own use case is a possible third). The story text
itself is authored during the plan's fixture task and reviewed by the founder before fixtures
are written.

## Not Doing / Out of Scope

- Paging, on-call, escalation, ticketing, status pages, SLA workflow (integrate, never build).
- A second context model, query engine, policy evaluator, renderer or storage path.
- Storing result rows inside the incident; evidence is by reference.
- Live multi-user editing under authenticated principals in the MVP (needs Phase 4 grants and an
  authenticated server profile — Track B2). Incident-scoped read grants are MVP as a recorded
  request/approval workflow; their enforcement against a principal who can choose `--as`
  arrives with the same Phase 4 work.
- Cross-project semantic identity (one `Order.id` across projects); facts stay project-qualified.
- Voice capture, Slack/Teams commands, self-learning extractors, AI Coordinator, generated
  videos, training mode, controlled production writes (FUTURE tiers; see the reconciliation).

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI / server (`datatug-cli`, `datatug serve`) | Planned | incident store routing, event stream and projection, policy-rechecked reads, `datatug incident` verbs, `watch`, access requests | specscore://github.com/datatug/datatug-cli/feature/cli/incident |
| Web (`datatug-apps`) | Planned | incident pages, `incidentius` profile home, context roles and overlays, evidence and metrics tables, access request and approval controls | — |
| Core (`datatug-core`) | Planned | `IncidentRef`, event and projection types, `ExecutionRequest.incident`, `Scope.storeId`, board widget extension boundary and round-trip fixtures | — |
| Demo projects (`datatug-demo-projects`) | Planned | canonical demo incident and its recurrence fixture in `demo-project-1` | — |

Several acceptance criteria below are phrased through `datatug incident …` or the incident page. Under REQ:capability-role they belong in the CLI and web Implementations or must be reworded platform-neutral; that move is an Amending round tracked as a follow-up, not part of this matrix.

## Acceptance Criteria

### AC: houston-creates-incident-with-context (verifies REQ:houston-creation, REQ:incident-store)

**Given** the demo project served with `--as admin` and the default store routing
**When** the user creates an incident titled "Canadian invoices stuck" with the text "customer 5 and invoice 42 affected"
**Then** `incidents/INC-1/events.jsonl` in the project's own repository holds an `incident.created` event, the description is kept, the reporter is a participant, and the canonical context layer exists.

### AC: incident-store-routes-to-configured-location (verifies REQ:incident-store)

**Given** the served project configured with a dedicated incidents repository as its store
**When** an incident is created and the project's repository is inspected
**Then** the events and projection exist only in the dedicated repository, the project repository is untouched, `datatug incident show` resolves the incident by `storeId` + `INC-1`, and the same test passes with the application-repository location.

### AC: incident-spans-two-projects (verifies REQ:multi-project-incidents)

**Given** an incident referencing the demo project and a second project under the `support` principal, whose grants cover only the second project
**When** the projection is read
**Then** facts and evidence from the demo project appear redacted or denied under that principal while the second project's evidence is complete, and both project references remain listed.

### AC: query-lists-its-incidents (verifies REQ:incident-references-and-back-links)

**Given** `customer-invoices` was recorded during INC-1 and a check was promoted from it during INC-1
**When** `datatug incident list --query customers/customer-invoices` runs and the promoted check file is inspected
**Then** INC-1 is listed with its `storeId`, and the check file carries `origin.incident` equal to INC-1's `IncidentRef`.

### AC: board-lists-its-incidents-by-ref (verifies REQ:incident-references-and-back-links)

**Given** a board file whose `incidentRef` is `{storeId: <demo store>, incidentId: INC-1}` and an `INC-1` event that references that board id
**When** `datatug incident list --board <boardId> --json` runs, and the board is round-tripped through the Go validator and the `@datatug/board-models` types
**Then** the list contains `INC-1` and no other incident, and the `incidentRef` survives the round trip unchanged in both directions.

### AC: investigation-becomes-incident (verifies REQ:houston-creation, REQ:attached-investigation-context)

**Given** an Investigation Context holding `Customer.ID = 5` and `Customer.Country = Canada`
**When** the user chooses "make this an incident"
**Then** the new incident's canonical layer references those facts with their origins and project scope intact, and the original context is unchanged.

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

### AC: events-redacted-under-policy (verifies REQ:events-and-projection-policy-rechecked)

**Given** INC-1 whose canonical context holds `Customer.Email = a@b.com` and a note quoting it, and the `support` principal for whom `Customer.Email` is hidden
**When** `support` reads the projection, the event stream and `watch`
**Then** the fact value is redacted in all three, the note is withheld per the annotation visibility rule, `admin` still sees both, and nothing was rewritten in the store.

### AC: watch-streams-events-with-cursor (verifies REQ:watch-event-cursor)

**Given** `datatug incident watch INC-1 --json` running
**When** another actor appends two events and the watcher restarts with `--since <cursor>` taken from its stderr
**Then** exactly the two events were streamed as NDJSON on stdout before the restart and nothing is duplicated after it.

### AC: outcomes-and-duplicate-merge (verifies REQ:lifecycle-and-outcomes)

**Given** INC-3 is a duplicate of INC-1
**When** `datatug incident merge INC-3 --into INC-1` runs
**Then** INC-1 gains INC-3's events with `incident.merged` refs, INC-3 is closed with `mergedInto: INC-1` and no outcome, and resolving INC-1 later requires exactly one of the five outcomes.

### AC: overlay-promotion-updates-canonical (verifies REQ:attached-investigation-context)

**Given** hypothesis H17 with an overlay layer holding `Customer.ID = 11..15` with role `suspected`
**When** H17 is confirmed and its facts promoted
**Then** the canonical layer contains them with role `affected`, the overlay still records the original facts and the promotion event, and an open query bound from context shows a visible re-binding choice rather than silently changing.

### AC: hypothesis-rejected-keeps-evidence (verifies REQ:hypotheses, REQ:evidence-links)

**Given** H12 "retry malfunction" with two execution records attached
**When** H12 is set to `rejected` with reason "retries completed for affected orders"
**Then** the projection still lists H12, both records, the reason and the queries used, each record carries INC-1's `IncidentRef` in its provenance, and each opens under the reader's current policy.

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

### AC: access-request-approval-and-expiry-are-events (verifies REQ:access-requests-in-incident)

**Given** `support` is an investigator on INC-1, which references the demo project, and may not read `Invoice`
**When** `support` runs `datatug incident access request INC-1 --source chinook --collection Invoice --reason "stuck invoices"`, `admin` runs `datatug incident access approve INC-1 <requestId>`, `support` reads `Invoice` with `--incident INC-1`, and INC-1 is later closed
**Then** the stream carries `access.requested`, `access.approved` and `access.expired` in that order with their actors, the execution record of the read (made with `--incident INC-1`) names the grant id, `datatug incident show INC-1` listed the grant as active before closure and absent after it, and a read after closure is limited as before.

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
**Then** it walks create → context with cohorts → checks → hypotheses → evidence → access request and approval → compare → metrics → mitigation → recovery → resolved → watching → closed → Resolution Record against the real server with no intercepted responses, asserting events and records, not rendered text alone.

### AC: recurrence-epilogue-boring (verifies REQ:canonical-demo-incident, REQ:search-and-similarity)

**Given** the resolved fixture incident and the recurrence fixture
**When** the recurrence is opened
**Then** `similar` names the first incident with matching signals, its checks run and the recurrence resolves with a recorded TTRI, demonstrating *"every incident should make the next one more boring"*.

## Open Questions

- **Canonical story text.** Setting ruled (REQ:canonical-demo-incident); the lead's candidate
  on `demo-project-1` is "Canadian invoices stuck", cohorts by `Customer.Country`, the FX-rate
  provider as the external counterparty, a deployment as the red herring. The founder reviews
  the written story before fixtures are authored (plan task 11).
- **Store routing key and default.** Lead assumption: `stores.incidents` in machine-local or
  served-project configuration, defaulting to the project's own repository; evidence records
  follow it. Confirm the key name and whether the hosted default is the space store.
- **Id scheme** `INC-<n>` per store plus `uid` — lead assumption.
- **Incidentius brand and domain** — [product-profiles](../product-profiles/README.md).

---
*This document follows the https://specscore.md/feature-specification*
