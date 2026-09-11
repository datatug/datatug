# DataTug roadmap

**Evidence baseline:** 9 September 2026 [independent review](spec/research/2026-09-09-independent-review.md).
**Authority:** [product thesis](spec/product/README.md), owning [Features](spec/features/README.md), and the founder decisions below.

## Journey and sequencing rule

An investigator opens a project, gets a useful result, understands a selected fact,
discovers reusable knowledge, confirms bindings and executes safely across sources.
The next investigator can reuse that knowledge. Start with this loop, then deepen
reuse, guide a short sequence, share it under grants, and present persisted dashboards.
Every phase names an observable result and its own end-to-end gate. A component,
merged branch or local test is not a completed journey or release.

## Tracks (2026-09-11)

The roadmap now carries three tracks for one product family
([product-family](spec/product/product-family.md)). Phases 0–5 below are **Track A, the
DataTug spine**, unchanged. **Track B, Incidentius** and **Track C, Dashboardius** follow
Phase 5 in this document and depend on Track A as drawn. Track B no longer starts beside
Phase 2: it starts only after the **Core coherent-experience review** gate below (founder,
2026-09-11: *"complete library refinement - e.g. DataTug improvements first. complete Phase
2 and review the core as a coherent experience before prioritizing Incidentius."*). The one exception is the `incidentius` profile scaffold, which
runs now beside Phase 2 (founder, 2026-09-11: *"I actually would prefer if you do in parallel scaffolding of the incidentius app"*):

```text
Track A  DataTug spine        0 ── 1 ── 2 ── 3 ── 4 ── 5
                                         │    │    │    │
Track B  Incidentius          B1 ────────┘    │    B2 ─┘    B3 ── B4
Track C  Dashboardius         C1 = Phase 5 ───┘    C2 ────── C3 (shell rebuild; full cutover)
```

The one part of B1 running beside Phase 2 is the `incidentius` profile scaffold (Incidentius plan
Task 16).

## Phase 0 — existing spine and delivery repair

**Outcome:** `datatug demo` opens the actual project; source → Album gives real rows
without hidden setup. The CLI builds, validates project credentials and serves a
registered project. Historical build/store/path/auth fixes are implementation evidence,
not remaining greenfield tasks. The review built current CLI and opened the local app;
project summary returned 200. Do not continue blaming the old auth-hook panic without
reproducing it on the candidate.

**Gate:** real entry-point browser test plus an exact client/server/demo revision
manifest. Phase 1 is complete (2026-09-11): corrective tasks 12–17 landed; the release candidate is datatug-cli v0.23.0 + datatug-apps 8653bd9, live at https://datatug.app, journey CI required on datatug-apps main.
**Scope:** no new product concepts, no unrelated command cleanup or UI-kit migration.

## Phase 1 — converge the core investigation loop

[Feature](spec/features/core-investigation-loop/README.md) ·
[transport contract](spec/features/core-investigation-loop/api-contract.md) ·
[plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-09-phase-1-core-investigation-loop/README.md).

**Outcome:** browse → semantic selection → related data/applicable query → explicit
bindings → authorized execution → context → continued investigation. SQLite, inGitDB
and keyless HTTP JSON GET all participate. Country/currency facts must actually be
selected or retrieved; they are not guessed from an ID. The support principal gets
only allowed rows/columns on every participating source; native SQL fails closed for
that protected profile. Closing/reopening and safe share/replay have separate ACs.

**Included:** exact transport, one source resolver, typed scoped context, missing and
ambiguous inputs, policy-safe lookups/counts, actual SQL/DTQL binding, HTTP browser
execution, explicit live/snapshot choice and J1–J4 real-agent acceptance. Preserve the
semantic engine, adapters and stronger VM ACL implementation.

**State:** the whole journey now walks end to end. Corrective task 12 (freeze the
transport, one source resolver) landed on all three sides — server (datatug-cli),
client (datatug-apps), schema authority (datatug-core `pkg/apicontract`, v0.26.0 and
later) — and task 10 (journey e2e J1–J4) is complete as of 2026-09-10: J1–J4 pass
single-worker against a datatug-cli built from main and the real web client, J4
against its own `--as support --role support` agent, and the `journey e2e` CI job runs
them on every datatug-apps pull request against a pinned datatug-cli release. The
review's wire mismatches, the parameter echo without application, the hardcoded
store default, the schema-qualified policy path and the missing table-page
limitation header are all fixed. Tasks 13–17 landed 2026-09-11 (the plan's status is Implemented; the release candidate is
named in the Phase 0 gate line above).

**Dependencies:** core/provider revisions first, then server and browser convergence,
then pinned demo/transport/journey acceptance. Layered security is integrated where
needed now; it is not deferred past protected execution.
**Non-goals:** boards/charts, hosted tenancy, GraphQL/authenticated HTTP, TUI parity,
workflow engine, cross-source joins and new UI shell.
**Autonomy:** high after reading the corrected contracts; engineering defaults are
explicit and no Dashboardius binding decision blocks this slice.
**Delivery gate:** exact candidate CI, all real journeys, dependency convergence,
authoritative remote receipts, a named DataTug CLI release and app deployment with
smoke evidence. Record each separately; required CI cannot succeed by skipping setup.

## Phase 2 — executable library quality and provenance

[plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-09-phase-2-library-depth-and-provenance/README.md).

**Outcome:** one engineer captures a useful query; another authorized user of the same
Git/local project can find, understand and rerun it with correct bindings and provenance.
The first delivery writes ordinary project query files; it does not claim hosted team
membership or a promotion lifecycle. A read-only principal cannot save.
**Included:** save as project query, folders plus README descriptions, result source/
query provenance, explainable mappings and a bounded accept/correct flow.
**Dependencies:** Phase 1 accepted; explicit server project-write authorization and
atomic conflict-aware persistence. No credential-bearing query is committed.
**Gate:** capture → diff → close → reopen → rerun, with read-only denial and write
failure preserving the draft. Measure reuse in an actual team experiment.
**Excluded:** credential vault design, new HTTP verbs, TUI/agent parity, rich collection
metadata and time-window semantics. These remain named future work in owning Features,
not silently deleted and not required to fit an arbitrary task count.
**Autonomy:** high for the bounded file/Git journey; hosted sharing has its own phase.

## Core coherent-experience review (gate)

**Trigger:** Phase 2's gate is met (capture → diff → close → reopen → rerun, with read-only
denial and write failure preserving the draft).
**Process:** the founder reviews the DataTug core — Phases 0–2 — as one coherent experience,
not as separate merged components.
**Outcome:** the review decides what comes next. Track B (Incidentius) does not start before
this gate passes (founder, 2026-09-11: *"complete library refinement - e.g. DataTug
improvements first. complete Phase 2 and review the core as a coherent experience before
prioritizing Incidentius."*). This supersedes Track B's earlier "runs beside Phase 2"
sequencing. The one exception is the `incidentius` profile scaffold (Incidentius plan Task 16),
which runs beside Phase 2 (founder, 2026-09-11: *"I actually would prefer if you do in parallel scaffolding of the incidentius app"*).

## Phase 3 — human-guided investigations

**Outcome:** from an entry ID, a person follows a short sequence of existing library
queries, chooses a branch, and reaches an explicit terminal summary. A second person
can replay the sequence with their own values and permissions.
**Included:** named entry points, query references, explicit transitions and context
handoff; preserve provenance and authorization at every step.
**Dependencies:** Phase 2 library reuse; an owning executable-knowledge Feature plan
must define the concrete sequence before implementation.
**Gate:** all branches and both close/replay endings run end to end. No duplicated SQL
inside a second workflow executor.
**Non-goal:** unattended execution for this phase. Future automation remains a later
evidence-based decision, not a permanent product prohibition.
**Autonomy:** high for the bounded human-guided slice; do not bundle hosting or boards.

## Phase 4 — project grants and distribution

**Outcome:** an authorized colleague opens the same standalone DataTug project through
a granted team space; revocation removes access. Project identity survives sharing to
multiple spaces, and a denied space cannot recover data through a shared link.
**Included:** project grants using space membership/roles and explicit identity binding.
A hosted or remote-store lane requires its own bounded plan and authenticated server
contract; deployment topology is not chosen implicitly by this roadmap.
**Dependencies:** accepted protected execution, published/converged layered ACL where
required, and sharing-scopes acceptance. There is one project, not duplicate
`/spaces/{spaceID}/ext/datatug/` project records masquerading as shared ownership.
**Gate:** grant, second actor read, revoke, denied replay and cross-project isolation.
**Autonomy:** medium until hosting/storage choices are specified in their owning
Features. OpenVaultDB project storage is retained as approved future capability, not
silently removed or treated as already shipped.

## Phase 5 — Dashboardius real persistence

**Outcome:** real query → card → save → reload → authorized colleague. DataTug owns
sources/execution/security/project assets; Dashboardius owns dense presentation,
widgets, editing, Undo/Reset and links.
**Dependencies:** [dashboards Feature](spec/features/dashboards/README.md), lossless
shared board contract, server write and concurrent-edit policy, and recipient
authorization. The SQL-widget binding is already ruled (founder, 2026-09-09: *"I'm Ok
with the suggested option 1"*) — a widget references a library query by `queryId` plus
bindings.
**Gate:** a genuine cross-language board round trip and real persisted two-actor
journey; scripted sample cards do not satisfy it.
**Autonomy:** high; no unresolved Feature decision blocks this phase. Presentation-only
refinements may proceed independently.

## Track B — Incidentius (incident-resolution workspace inside DataTug)

Owning Feature: [incidents](spec/features/incidents/README.md) with
[evidence-records](spec/features/evidence-records/README.md),
[checks](spec/features/checks/README.md), [annotations](spec/features/annotations/README.md),
[compare](spec/features/compare/README.md) and
[product-profiles](spec/features/product-profiles/README.md). Vision:
[research record](spec/research/2026-09-11-incident-resolution-vision.md); mapping:
[reconciliation](spec/research/2026-09-11-incident-vision-reconciliation.md).

### B1 — Incidentius MVP

[Plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-11-incidentius-mvp/README.md).

**Outcome:** on a local or Git project, "Houston, we've got a problem" creates an incident
attached to an Investigation Context with cohorts and hypothesis overlays; checks run and are
recorded; hypotheses carry evidence by reference to execution records, snapshots and
annotations; metrics form series with explicit resolution criteria; recurrence and search work;
`datatug incident` and `watch` expose the same event stream to agents; a participant requests
temporary read access inside the incident and an approver grants it, recorded as events and
expired at closure (MVP inclusion is the founder's, 2026-09-11; the mechanism is the lead's);
the canonical demo incident and its boring
recurrence walk end to end.
**Dependencies:** Phase 1 accepted (done). Incidents and evidence write to the routed incident
store through the server; checks and annotations are project knowledge assets and need
Phase 2 task 2 (the authorized project-write path), sequenced first. Track B does not start
until Phase 2 completes and the [Core coherent-experience review](#core-coherent-experience-review-gate)
gate passes and prioritizes it (founder, 2026-09-11). The profile scaffold, plan Task 16, is
exempt and runs now (founder, 2026-09-11: *"I actually would prefer if you do in parallel scaffolding of the incidentius app"*).
**Non-goals:** hosted multi-user incidents, boards/charts, AI narratives, capture channels,
production writes (approved mutating actions are FUTURE).
**Gate:** the whole-journey e2e with both epilogues, pinned revisions, a named CLI release and
a deployment receipt.

### B2 — collaboration and incident dashboards

**Outcome:** several authenticated people work one incident live under project grants;
incident dashboards (impact, recovery, watch) render in Dashboardius and embedded in the
incident page; compare has a UI; AI status projections for troubleshooters, leadership and
support are generated as `inference` events with refs; appreciation credits; replay UI.
**Dependencies:** Phase 4 (grants, authenticated server) for live collaboration; Phase 5 (board
persistence) and Track C2 for dashboards. Incident-scoped temporary read grants are B1
(`server-enforced-access` REQ:incident-scoped-grants); B2 adds them for authenticated
multi-user sessions.

### B3 — capture channels and coordination

**Outcome:** pasted text and messages become typed facts with provenance and confidence
(deterministic extractors first); read-only Slack/Teams capture into the same stream; manual
deployment events and deployment correlation; SpecScore and Codegrapher links where present.
Autonomous agents watch new incidents and propose hypotheses and safe checks through the same
API. **Dependencies:** B2; agent-context parity.

### B4 — the knowledge flywheel

**Outcome:** AI-proposed checks tested against history and promoted cautiously; self-learning
extraction skills; AI Coordinator; ownership and expertise from artefacts; controlled
production modifications with approval; masking policy for AI-visible data; generated
Resolution Record documents and stories. Every item requires evidence from B1–B3 before a plan.

## Track C — Dashboardius

- **C1 = Phase 5** (real persistence), unchanged.
- **C2 — incident and investigation dashboards:** `metric-series`, comparison, criteria and
  replay widgets over the incident event stream and execution records; the board renderer as a
  workspace library embedded by DataTug/Incidentius pages. Owning Feature:
  `sneat-co/dashboardius` `incident-dashboards` and its plan. Depends on C1, B1 and the
  Incidentius MVP plan's task 13 (board widget extension boundary).
- **C3 — one-app cutover:** Dashboardius becomes the `dashboardius` profile of the single
  configurable app in `datatug/datatug-apps` (founder, 2026-09-11: *"one configurable app."*).
  Its shell — cards, lists, toolbars, menus, forms — is rebuilt with Ionic components; PrimeNG,
  AG Grid, Chart.js and CDK drag-drop survive only as leaf widgets in the board renderer
  library; the PrimeNG shell in `sneat-co/dashboardius` is retired. A complete cutover under
  the migrate-means-migrate rule (inventory in [product-profiles](spec/features/product-profiles/README.md)).
  Can start after the renderer library exists (Dashboardius plan task 1); not gated by C1.
  **No executable owner plan exists yet (F8):** the inventory names the work, but deployment
  transfer, old-shell removal and domain smoke have no task chain. C3 needs its own plan,
  with deployment-before-removal acceptance, before it starts.

## Independently planned extensions and experiments

Retain these opportunities without pretending they are one autonomous phase:

- OpenVaultDB project-store round trip and project-store self-exploration: owning
  git-native-projects Feature; bounded after its policy and record contract is ready.
- Agent context, MCP/brokered builder and TUI handoff: same semantic/execution boundary;
  no alternative authorization path. Separate parity plans.
- Authenticated HTTP/GraphQL: credential authority, descriptors and SSRF contract first.
- Time-window semantics and richer collections: user evidence before new metadata.
- Relationship/meaning suggestions: measure precision and correction cost before UI;
  manual declarations remain authoritative.
- Infrastructure inventory and cross-source joins: demonstrate a concrete investigation
  that needs them; reject unsupported query shapes rather than weakening policy.
- AI-assisted capture and public knowledge packs: reuse existing library contracts.

## Dashboardius and DataTug

Founder words: "A Dashboardius board IS a DataTug board; what Dashboardius adds is
widget kinds and an editing experience." The board belongs to the DataTug project.
Dashboardius may be sold separately; it must not become a second data platform.
The shared renderer may be embedded later if useful; no independent second renderer
is planned. No charts or boards are added to the Phase 1 investigation task.

Founder words, 2026-09-11: "Not sure if it's part of the spec your read but I think DataTug,
Incidentius, Dashbaordius apps code should be in same repo (nx workspace with shared libs and
either separate or single configurable app)". The three apps' code therefore lives in the
`datatug/datatug-apps` Nx workspace, and the founder settled the second half the same day:
*"one configurable app."* — all three profiles in one Ionic-framed app, PrimeNG and friends
as leaf components only ([product-family](spec/product/product-family.md) §2). This supersedes
the standalone-repository part of the 2026-09-08 Dashboardius Idea, not its schema decisions.
DataTug CLI is Apache 2.0; the web UI is private for now, with future licensing not
selected. Do not market all DataTug surfaces as already open source.

## UI kit decision

Founder words, 9 September: "UI-kit rule - confirmed". DataTug keeps the Ionic frame
and bounded islands for data surfaces it lacks, reusing the existing shared grid.
Amended 2026-09-11 (founder: *"one configurable app."*; *"We should migrate PrimeNG cards
and other components that we can migrate to Ionic and use PrimeNG just as components if
Ionic misses functionality"*): Dashboardius's shell is rebuilt with Ionic components inside
the single app, and PrimeNG, AG Grid, Chart.js and CDK drag-drop remain only as leaf widgets
in the board renderer library — the one UI-kit migration the family makes (Track C3).
Shared presentation libraries serve all three profiles.

## Recorded founder decisions

- MVP sources: "I want to have in MVP: inGitDB, SQLite, HTTP endpoint. Regards HTTP -
  some public sources like countries or currency rates or something that is common
  and is understandable and can looks good in demos."
- Scope approval: "I approve all artefacts - including Phase 1 + addition of HTTP data sources".
- HTTP adapter: "dalgo2http - yes".
- Layered ACL: "It was authorised and is being implement on VM by codex agent. I've
  asked it to drive it to completion." Authorization is not a remote landing receipt.
- Board ownership: "board belongs to DataTug project. Project can be shared across Sneat spaces".
- Project scope: "DataTug project is standalone entity that can be shared (granted)
  between multiple spaces (teams)". This is an explicit exception to the platform's
  usual space-bounded record default.
- Licensing: "CLI - Apache 2.0. Web UI private for now, will be AGPL or something like that".
- Ahead-of-code marketing claims: "keep with planned label".
- Apps workspace (2026-09-11): "DataTug, Incidentius, Dashbaordius apps code should be in same
  repo (nx workspace with shared libs and either separate or single configurable app)".
- One app (2026-09-11): "one configurable app." — all three profiles; on the frame: "We should
  migrate PrimeNG cards and other components that we can migrate to Ionic and use PrimeNG just
  as components if Ionic misses functionality".
- Multi-project incidents (2026-09-11): "yes it can."
- Demo story (2026-09-11): "Chinook first, our own inGitDB second, maybe a copy of Chinook or
  our own use case".
- Incident access approver (2026-09-11): "Someone with appropriate role. By default project
  admin has this role but can move it to someone else. It's part of Sneat contactius role
  management".
- Incidentius positioning (vision document, 2026-09-11): "Incidentius is an incident-resolution
  workspace." MVP: "do not build three applications. Build Incidents inside DataTug".
- Pilot order, external roadmap/plans review (2026-09-11): "complete library refinement -
  e.g. DataTug improvements first. complete Phase 2 and review the core as a coherent
  experience before prioritizing Incidentius."
- Snapshots in Git, external roadmap/plans review (2026-09-11): "Not in MVP".
- Web compare scope, external roadmap/plans review (2026-09-11): "compare - A".
- Incidentius scaffold timing (2026-09-11): "I actually would prefer if you do in parallel
  scaffolding of the incidentius app".

## Open Questions

Product choices live beside their owning Feature: hosted identity and sharing in
[sharing-scopes](spec/features/sharing-scopes/README.md#open-questions), and
distribution/privacy interpretation in
[surfaces](spec/features/surfaces/README.md#open-questions). No release version, provider
change or deployment success is inferred from this plan.

- **C3 owner plan (F8):** Track C3 has a cutover inventory but no executable owner plan yet
  (deployment transfer, old-shell removal, domain smoke). See Track C3 above.
- **Personal-query storage acceptance owner:** personal-query storage has no acceptance
  owner and is outside Phase 2's project-query-only scope; it must not silently expand
  that scope. An owning Feature/plan needs to claim it explicitly.
