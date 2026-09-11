# DataTug · Incidentius · Dashboardius — one product family

**Status:** Draft for founder review · **Date:** 2026-09-11 · **Owner:** alex
**Source:** the founder's [incident-resolution vision](../research/2026-09-11-incident-resolution-vision.md)
(cited below as *vision §n*) reconciled against the hub in
[2026-09-11 reconciliation](../research/2026-09-11-incident-vision-reconciliation.md).
**Authority:** the [product thesis](README.md) still says what DataTug is; this page says how
three product profiles share it. When they disagree, fix one the same day.

## 1. What the family is

DataTug is the technical backbone. Incidentius is the incident-resolution experience.
Dashboardius is the dashboard experience. They are **three product profiles of one
application**, not three implementations (vision §1, §4).

```text
                       ONE APPLICATION FAMILY
                                │
             shared domain model (datatug-core) + one trusted server
                                │
                   product profile (hostname / configuration)
                                │
          ┌─────────────────────┼─────────────────────┐
          ▼                     ▼                     ▼
       DataTug             Incidentius            Dashboardius
   "Interrogate the      "Why are we            "Show the evidence,
    system and data."     investigating,          impact, progress,
                          what have we learned,   recovery and the
                          what happens next?"     useful views."
          │                     │                     │
   home = project /      home = incidents        home = boards
   investigation
```

Positioning statements to preserve verbatim (vision §61):

> **Incidentius is an incident-resolution workspace.**
> **Build capabilities once. Expose them through DataTug, Incidentius and Dashboardius product profiles.**
> **Integrate deeply, depend loosely.**

Incidentius is *not* an incident-management platform: no paging, on-call, ticket workflow,
status pages or SLA machinery. It integrates with those (vision §3.1). Every feature must
answer: *does this help somebody understand or resolve the problem faster?*

## 2. What "one application" means

| Layer | One of them | Where |
|---|---|---|
| Domain model | one Go model: projects, entities, queries, boards, **incidents**, execution records, annotations, checks | `datatug-core` |
| Trusted server | one execution and access boundary (`datatug serve`, later hosted); every read, write, snapshot and incident event passes through it | `datatug-cli` `pkg/server` |
| Contracts | one transport appendix per capability; `@datatug/board-models`; incident event schema | hub Features |
| Apps code | **one Nx workspace** holding every app and the shared libraries | `datatug/datatug-apps` |
| Product shells | hostname/configuration selects home, brand, navigation and entry points | [product-profiles](../features/product-profiles/README.md) |

**Founder ruling, 2026-09-11 (verbatim):** *"Not sure if it's part of the spec your read but I
think DataTug, Incidentius, Dashbaordius apps code should be in same repo (nx workspace with
shared libs and either separate or single configurable app)"*.

This supersedes the repository half of the 2026-09-08 Dashboardius Idea, which chose a
standalone repository for the Dashboardius app. It does not decide "separate apps" versus
"one configurable app"; that remains the founder's choice. The lead's recommendation
(assumption, not a ruling):

- **DataTug and Incidentius are one configurable app** (`apps/datatug-app`) selected by
  profile. They share the Ionic frame, the agent client, the Investigation Context and every
  project screen; an incident page is a project page with a different home.
- **Dashboardius is a separate app in the same workspace** (`apps/dashboardius`). It is
  PrimeNG with Chart.js and CDK drag-drop, deploys to its own domain on its own cadence, and
  keeps the "instrument panel" design. Folding it into the Ionic shell would put two UI kits
  and two bundles behind one entry point for no shared code beyond the libraries below.
- **Shared libraries** carry the reuse: `@datatug/board-models` (exists), a Dashboardius-owned
  **board renderer library** that DataTug/Incidentius pages embed as an island (the founder's
  UI-kit rule, 2026-09-09: Ionic frame plus PrimeNG/AG Grid islands only where Ionic lacks
  components), the agent client, auth, build-info.
- Moving `sneat-co/dashboardius/web` into the workspace is a **full cutover**
  (rule `cutover-verbs-mean-full-cutover`): its Cloudflare worker, `wrangler.jsonc`, the shared
  `cf-deploy.yml` call, Playwright e2e, Renovate config, PrimeNG theme and licence scripts all
  move or are re-pointed; the old app path is removed; `sneat-co/dashboardius` stays the
  Dashboardius **spec hub** (like this repository is for DataTug). Track C3 in the
  [roadmap](../../ROADMAP.md) owns the inventory.

The founder's earlier rulings stand and compose with this one: *"A Dashboardius board IS a
DataTug board; what Dashboardius adds is widget kinds and an editing experience"*, *"board
belongs to DataTug project. Project can be shared across Sneat spaces"*, *"UI-kit rule -
confirmed"*, and the SQL-widget binding ruling (*"I'm Ok with the suggested option 1"*).

## 3. Ownership boundaries

| Concern | Owner | Rule |
|---|---|---|
| Sources, environments, semantic layer, queries, execution, access control, provenance | DataTug | Every profile reads and writes through the server; nobody else executes a query or evaluates a policy. |
| **Investigation Context** (typed facts, conditions, scope isolation, no hidden filters) | DataTug ([investigation-context](../features/investigation-context/README.md)) | **Firm boundary (vision §5):** Incidentius attaches an incident to a context; it never reimplements it. Cohort roles and overlays are extensions of this Feature. |
| Execution records and result snapshots | DataTug ([evidence-records](../features/evidence-records/README.md)) | Evidence is a DataTug artefact; incidents reference it. |
| Annotations | DataTug ([annotations](../features/annotations/README.md)) | Data-addressable knowledge; an incident is one provenance dimension. |
| Checks (health, diagnostic, hypothesis-validation, remediation-verification, regression) | DataTug ([checks](../features/checks/README.md)) | Library queries with expected conditions; runs are execution records. |
| Compare (same query across two scopes) | DataTug ([compare](../features/compare/README.md)) | An execution primitive; incidents and dashboards consume it. |
| Boards: schema, persistence, query binding, execution, ACL | DataTug ([dashboards](../features/dashboards/README.md)) | Unchanged. |
| Board presentation, widget kinds, editing, `metric-series`, incident/investigation dashboards, the embeddable renderer | Dashboardius (`sneat-co/dashboardius` specs) | Never a second data platform. |
| Incident identity, lifecycle, participants, hypotheses, evidence relationships, recurrence, outcomes, metrics semantics, event stream, status projections | Incidentius ([incidents](../features/incidents/README.md)) | Built on the rows above; owns resolution semantics only. |
| Product shells (home, brand, navigation, entry points) | [product-profiles](../features/product-profiles/README.md) | Presentation of the same capabilities; no profile-private data. |

## 4. Shared capability versus profile-specific experience

| Capability | Owning Feature | DataTug profile | Incidentius profile | Dashboardius |
|---|---|---|---|---|
| Browse, select, related data, applicable queries | core-investigation-loop, related-data-navigation | home loop | reached from an incident's context | picks queries for cards |
| Investigation Context with roles and overlays | investigation-context | context panel | canonical + hypothesis overlays, promote/reject | board parameters bound from context (later) |
| Execution records, snapshots, fingerprints | evidence-records | "record this run" on any result | evidence attached to hypotheses; recovery series | series behind `metric-series` |
| Annotations | annotations | on any value/row/result | incident-scoped notes with provenance | shown on cards (later) |
| Checks | checks | library entries with expected conditions; `datatug-action` in CI | run at incident start; hypothesis validation; resolution criteria | check status widgets |
| Compare | compare | env/env, snapshot/snapshot | affected vs control, PROD vs UAT, before vs after | comparison boards |
| Incident model, event stream, `datatug incident` | incidents | an investigation may become an incident | **home**: list, "Houston, we've got a problem", timeline, hypotheses, metrics, status | incident boards read the same events/records |
| Boards | dashboards (+ Dashboardius Features) | "pin this result"; embedded renderer island | incident dashboard: impact, diagnostic, recovery, watch, replay | **home**: authoring and editing |
| Access control, grants | server-enforced-access, sharing-scopes | per project | incident participants are project members; incident-scoped grants later | board readers follow project grants |
| CLI/API/agents | surfaces, agent-context | `datatug …` | `datatug incident …`, `watch` event stream | board API only |

Profile-specific by design: the **home route**, **brand** (name, logo, colours, copy, domain),
**navigation order**, **entry points** ("Houston" button vs "open a project" vs "new board"),
demo defaults, and marketing story. Everything else is shared.

## 5. Architecture

```text
Browser apps (datatug-apps Nx workspace)
  datatug-app  ── profile: datatug | incidentius ──┐
  dashboardius ── PrimeNG, board renderer lib ─────┤  HTTP /datatug/*  (never a DB connection)
                                                    ▼
datatug serve  (datatug-cli pkg/server, later a hosted profile)
  ├─ project store        → files (Git-native) · GitHub · OpenVaultDB (future) · hosted
  ├─ semantic resolver    → mappings, related lookups, applicable queries
  ├─ executor             → accesspolicies.Run → SecureReadSession → DALgo adapters (SQLite · inGitDB · HTTP · …)
  ├─ evidence store       → execution records, result snapshots (retention, masking, re-checked on read)
  ├─ checks runner        → library query + expected condition → execution record
  ├─ compare              → two scopes → keyed recordset diff
  └─ incidents            → append-only event stream per incident → projection → watch cursor
```

Incident data flow (vision §3.2, §57):

```text
people investigate (queries · context · hypotheses · notes · check runs)
        ↓ every action is an IncidentEvent {actor, assertion kind, provenance refs}
append-only events.jsonl  (source of truth; replayable at any time T)
        ↓ deterministic fold
incident.json projection  (current state; what the list/detail/API show)
        ↓ same stream
watch · timeline · status views · AI coordination · connectors · Resolution Record
```

`assertion.kind` is one of `observation` (a deterministic result or a recorded fact),
`claim` (a human statement), `hypothesis`, `inference` (AI or heuristic output, always with the
events it was derived from), `deterministic-result` (a check evaluation). Generated views keep
the links; nothing AI-produced is ever recorded as an observation (vision §12, §26, §31).

## 6. Staged delivery

The [roadmap](../../ROADMAP.md) owns sequencing. Summary of the tracks and their gates:

```text
Track A  DataTug spine        0 ── 1 ── 2 ── 3 ── 4 ── 5
                                    │         │    │    │
Track B  Incidentius          B1 ───┘         │    B2 ─┘    B3 ── B4
                              (MVP inside     │  (hosted collaboration,
                               DataTug)       │   incident dashboards)
Track C  Dashboardius         C1 = Phase 5 ───┘    C2 (incident boards) ── C3 (workspace cutover)
```

- **B1 MVP** starts after Phase 1 (accepted 2026-09-11) and runs beside Phase 2; it needs
  Phase 2 task 2 (the project-write path) only for "promote query → check".
- **B2** needs Phase 4 (project grants, authenticated server) for live multi-user incidents and
  Phase 5 (board persistence) for incident dashboards.
- **C3** (moving the Dashboardius app into the workspace) can start any time; it is
  mechanical but must be a complete cutover.

## 7. Ecosystem integration map — integrate deeply, depend loosely

Every row is optional amplification surfaced at the moment of demonstrated value; none is a
hard requirement of any profile (vision §3.6, §45, §46).

| Sibling | Integration | Depth | Coupling |
|---|---|---|---|
| **sneat.work** | DataTug powers `/data` (Tier-3 brand, ruling restated 2026-07-01); Incidentius is the incident-resolution capability of the work pillar; shared Sneat identity (`sneat-eur3-1`); project grants to Sneat spaces (Phase 4) | identity, grants, deep links | link out; no Sneat module owns incident data |
| **OpenVaultDB** | a DataTug project (incidents included) storable in and usable from a remote vault (founder 2026-09-09: *"It should be possible to store and use datatug project from remote OpenVaultDB"*); vault policies compose under the server's enforcement | project store backend | adapter behind `storage.Store`; absent vault changes nothing |
| **inGitDB** | contextual/reference record sets (cohorts, healthy controls) persisted as versioned collections; support-notes in the demo; the natural home of a Resolution Record's datasets (vision §12) | source + dataset persistence | one source among many via `dalgo2ingitdb` |
| **WB / Workbench** | incident evidence about a **deployment** can reference WB's operations journal and work logs (which branch landed when); `wb` remains the delivery tool for building this family | deployment correlation (NEXT) | read-only import of receipts; no WB dependency at runtime |
| **SpecScore** | features link expected behaviour to DataTug queries/checks; an incident finding becomes a requirement or regression scenario linked back to the incident (vision §32, §45); this family is specified in SpecScore | knowledge flywheel | optional link fields; absent SpecScore, findings stay in the incident |
| **Synchestra** | dispatches AI agents that consume `datatug incident watch --json` (vision §52); a persistent agent proposes hypotheses and runs safe checks | agent orchestration (NEXT/FUTURE) | agents are principals through the same server; Synchestra is one dispatcher |
| **Codegrapher** | from a table/column/stack-trace pattern to repositories, symbols, callers, recent changes, ownership (vision §33); MCP `serve` and `impact` already exist | code correlation (FUTURE) | absent index, the investigation continues |
| **Sneat.ai / conversational runtime** | the candidate-then-validate pattern (LLM proposes a semantic candidate, deterministic validation decides) is the template for AI status views, extraction and hypothesis suggestions | AI trust boundary | pattern reuse, not a service dependency |
| **datatug-action / CI** | checks run in CI (regression purpose); a check created during an incident guards the next release | flywheel | already shipped runner |
| **DALgo / DTQL** | every source, including logs later, is a DALgo adapter; DTQL is the portable query shape for checks and compare | data access | inherited from P1/P2 |
| **Slack/Teams** | alternate capture surfaces feeding the same event stream (vision §27) | capture (NEXT read-only, FUTURE commands) | connector, never a second incident system |

## 8. Decisions and their reasons

| ID | Decision | Why |
|---|---|---|
| D1 | One workspace, profiles for DataTug/Incidentius, Dashboardius a sibling app sharing libs | founder ruling above; UI-kit rule; no second renderer or data platform |
| D2 | Incidents are project assets (`incidents/<id>/`) | same ownership as boards (founder ruling); grants and ACL apply unchanged; Git-native portability of the Resolution Record |
| D3 | Event-sourced incident, projection derived, replay at T | vision §57 Q2–Q3, Q11; one stream powers watch/timeline/AI/connectors; wrong theories are never lost (§3.3) |
| D4 | Context roles and overlays extend investigation-context | vision §5 firm boundary; useful speculation usable without becoming truth (§11.1) |
| D5 | Execution records + snapshots are DataTug evidence | historical evidence, shared cache, collaboration, comparison over time (§15); no permission backdoor |
| D6 | Data-addressable annotations | annotate where evidence exists; resurface across queries (§16) |
| D7 | Checks generalise library queries | executable knowledge; deterministic before AI (§18–19) |
| D8 | Compare as an execution primitive | "compare something that works with something that does not" (§13) |
| D9 | Metrics are check series; criteria are explicit predicates | measure the problem and the recovery (§20, §22) |
| D10 | Profiles by hostname/config, one build | simplest mechanism (§4) |
| D11 | CLI/API-first, structured `watch` | humans, scripts and agents on one model (§50–52) |
| D12 | Assertion kinds separate fact, claim, hypothesis, inference | AI trust boundary (§12, §31, §57 Q12) |
| D13 | Knowledge spreads for free; collaboration creates the bill | growth mechanism; never a paywall during an outage (§47) |
| D14 | Canonical demo incident with a boring recurrence is an MVP deliverable | E2E fixture, demo, onboarding, marketing, regression in one (§53–54) |

## 9. Open Questions

Owner-routed; each also appears in its owning Feature.

- **Separate apps or one configurable app** for the workspace (founder; §2 recommendation).
- **Incidentius domain and brand tier** — no domain is registered (checked `sneat-ops/data/domains.json`); Tier-3 independent like DataTug, or a DataTug sub-brand? (founder)
- **Does Chinook carry the canonical demo story**, or does the family need a second demo project (orders/shipping)? ([incidents](../features/incidents/README.md) open questions; lead recommends trying Chinook first.)
- **Snapshot retention and masking defaults** ([evidence-records](../features/evidence-records/README.md)).
- **Whether `sneat-co/dashboardius` is kept as spec hub or archived after the cutover** (founder; recommendation: keep).
