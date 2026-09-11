# Reconciling the incident-resolution vision with the hub (2026-09-11)

**Input:** the founder's [vision document](2026-09-11-incident-resolution-vision.md) (sections
cited as §n). **Method:** every concept in it was mapped to the hub's Features, the transport
contract, `datatug-core`, `datatug-cli`, `datatug-apps`, `sneat-co/dashboardius`,
`sneat-co/backstage` Ideas and `datatug/marketing`, read at these revisions: hub `main`
55608f9, dashboardius `main` e560229, datatug-cli `main` 72bb00a, datatug-core v0.27.x
(`pkg/datatug/*.go`, `pkg/apicontract`, `pkg/comparator`). This record is evidence; the
[product-family](../product/product-family.md) page and the owning Features carry the decisions.

## 1. Required audit (vision §58) — what exists

| Area | Evidence | State |
|---|---|---|
| Investigation Context | Feature `investigation-context`; `api-contract.md` `Fact{entity,field,value,origin,physical?,mapping?,enabled}` with `condition?` specified as a pending lead assumption (not yet in `datatug-core/pkg/apicontract/fact.go`); scope keyed by agent URL/project/environment/securityContextId; conflict and no-hidden-filters rules; `datatug-apps` Investigation Context page (PR #65, #119) | usable; conditions pending; roles and overlays absent |
| Logical model | Feature `semantic-layer`; `EntityField.Mappings []PhysicalRef{Source,Collection,Column}`; `pkg/semantic.Resolve`; declared/inferred provenance | usable; per-source mappings cover per-environment physical differences because a source is environment-scoped in the registry |
| Environments / infrastructure | `environment.go`, `server.go`, `dbcatalog.go`, `datatug scan`; Feature `infrastructure-inventory` (services/endpoints vestigial; climb planned) | partial |
| DALgo / data sources | Feature `heterogeneous-sources`: SQLite, inGitDB, keyless HTTP GET through `dalgo2http`; logs = REQ:logs-source (planned, no model) | usable for three sources; logs planned |
| Query library | Feature `executable-knowledge-library`; folders, semantic parameters, applicability with chain; capture (Phase 2), playbooks (Phase 3) | usable / partial |
| Health checks | `datatug-core/pkg/datatug/checks.go` `Check{ID,Title,Type: options\|regexp\|sql,Data}` (dataset value checks); `datatug-action` runs the validator in CI; `datatug-cli/spec/ideas/db-snapshot-command.md` (Approved) is a bulk data snapshot, not a health check | data-validation checks only; no purpose/expected-condition model |
| Snapshots / caching | `Result.provenance{mode: live\|snapshot, snapshotId, observedAt}` for recorded HTTP fixtures; no execution record, no result store, no fingerprint | seed only |
| Annotations | none in any repo | absent |
| inGitDB | `dalgo2ingitdb`; demo `support-notes` collection; Feature `git-native-projects` (project store explorable) | usable |
| Dashboardius integration | hub Feature `dashboards`; `boards.go` with `SQLWidgetSettings{queryId, parameters}`; `@datatug/board-models` in `datatug-apps/libs/datatug/board-models`; `sneat-co/dashboardius` standalone PrimeNG app, scripted demo, Feature `query-backed-board-persistence` + plan (Draft) | model converged; persistence not built |
| Access control | Feature `server-enforced-access`; `pkg/accesspolicies`, `pkg/secureread`, DALgo `access`; protected vs opaque-privileged profiles; layered ACL track; fixed principal per `serve` | usable locally; multi-user hosted pending |
| Auditing | server logs record correlation ID and operation, not raw facts (`api-contract.md`); no per-execution audit record | partial |
| CLI conventions | singular namespaces `datatug entity\|query\|serve\|scan\|demo\|validate`; `--json`/format flags; specs per command under `spec/features/cli/<cmd>/` | usable |
| AI/agent surfaces | `ai-plugin` skills (8/9 verified); Features `agent-context` (planned `datatug context`), `ai-query-builder` and `serve-brokered-query-builder` (Approved, no code); ecosystem pattern `sneat-ai-backend` (candidate-then-validate) | partial |
| Compare | `pkg/comparator` compares **schemas** across environments; `diffs.go` schema/table/column diffs | schema compare only |

## 2. Concept map: existing · extend · new

| Vision | Concept | Disposition | Owner Feature | Tier |
|---|---|---|---|---|
| §4 | Product profiles by hostname/config | **new** | product-profiles | MVP |
| §4, ruling 2026-09-11 | One Nx workspace for the three apps | **new (cutover)** | product-profiles; roadmap C3 | NEXT |
| §5 | Incident attaches to Investigation Context | **new on existing** | incidents ← investigation-context | MVP |
| §6 | Operational map: environments, infrastructure, sources | existing / extend later | infrastructure-inventory | FUTURE (climb) |
| §7 | Logs as first-class sources; cross-source correlation | planned (REQ:logs-source); joins are a non-goal today | heterogeneous-sources | FUTURE |
| §8 | Logical model with physical mappings per source | existing | semantic-layer | — |
| §9 | Improve the model while investigating | existing (Phase 2 accept/correct) + evidence-backed suggestion | meaning-inference, provenance | NEXT (evidence count) |
| §10 | Data flow and lineage; divergence point | **new** | infrastructure-inventory (later) | FUTURE |
| §11 | Cohort roles on context values | **extend** | investigation-context | MVP |
| §11.1 | Layered/scoped context, promote/reject | **extend** | investigation-context + incidents events | MVP |
| §12 | Source-agnostic record sets, inGitDB persistence of cohorts | existing sources + **extend** (save a cohort as an inGitDB collection) | git-native-projects, heterogeneous-sources | NEXT |
| §13 | Cross-environment comparison; the compare primitive | **new** | compare | MVP (API/CLI, two scopes) · NEXT (UI) |
| §14 | External agents / counterparties | **new** semantic entity kind | semantic-layer | FUTURE |
| §15 | Query-result snapshots as immutable evidence | **new** (seed: `Result.provenance`) | evidence-records | MVP |
| §16 | Row/value/entity annotations, data-addressable | **new** | annotations | MVP (value, execution record) · NEXT (row sets, relationships) |
| §17 | Hypotheses, evidence, uncertainty reduction | **new** | incidents | MVP |
| §18 | Human → AI-assisted → deterministic checks | **extend** (`Check` type) | checks | MVP (deterministic) · FUTURE (AI proposal) |
| §19 | Health checks at incident start | **extend** | checks, incidents | MVP (run listed checks manually) · NEXT (auto-run safe checks) |
| §20 | Incident metrics and resolution criteria | **new** | incidents (semantics), evidence-records (series) | MVP (table) |
| §21 | Dashboardius during incidents | **new** | Dashboardius `incident-dashboards`, hub `dashboards` | NEXT (needs Phase 5) |
| §22 | Recovery and watch period | **new** lifecycle semantics | incidents | MVP (statuses) |
| §23 | Lifecycle and five outcomes; duplicates merge | **new** | incidents | MVP |
| §24 | Recurrence relationship, TTRI | **new** | incidents | MVP (link + timestamps) · NEXT (AI TTRI) |
| §25 | Search and similarity; knowledge decay | **new** | incidents (search), provenance (staleness) | MVP (deterministic overlap) · NEXT (ranking) · FUTURE (decay) |
| §26 | Voice capture | **new** | incidents (capture channels) | FUTURE |
| §27 | Slack/Teams connectors | **new** | incidents (capture channels) | NEXT (read-only) · FUTURE (commands) |
| §28 | Parse text into typed context | **new** (seeds: `VarSetting.ValuePattern`, `EntityField.NamePatterns`) | investigation-context, incidents | NEXT |
| §29 | Self-learning extraction skills | **new** | agent-context | FUTURE |
| §30 | AI Troubleshooter vs Coordinator | **new** | agent-context, incidents | FUTURE |
| §31 | Generated status views | **new**; deterministic projection first | incidents | MVP (troubleshooter text) · NEXT (AI, leadership, support) |
| §32 | Deployment / code / SpecScore correlation | **new** | incidents (deployment events), infrastructure-inventory | NEXT (manual deployment events) · FUTURE (correlation) |
| §33 | Codegrapher | **new** optional | incidents integration | FUTURE |
| §34 | Ownership and expertise | **new** | infrastructure-inventory / sharing-scopes | FUTURE |
| §35 | Incident-scoped privileged access | **extend** | server-enforced-access, sharing-scopes | NEXT |
| §36 | Controlled production modifications | **new** | server-enforced-access | FUTURE |
| §37 | Masking and AI-visible data | **extend** | server-enforced-access | FUTURE (policy dimensions) |
| §38 | Incident Resolution Record as a mini-site | **new** projection over the event stream | incidents | MVP (record page) · NEXT (RCA/postmortem generation) |
| §39 | Timeline replay and training mode | **new** | incidents | MVP (`show --at`, timeline) · NEXT (replay UI) · FUTURE (training) |
| §40 | Generated videos | **new** | — | FUTURE |
| §41 | Appreciation/credits | **new** projection | incidents | NEXT |
| §42–43 | Personas, incident role, permission | spec content | incidents, backstage Idea | MVP (roles list) |
| §44 | Knowledge flywheel | thesis | product/README | — |
| §45 | SpecScore links | **new** optional link fields | incidents | NEXT |
| §46 | OpenVaultDB / inGitDB | existing direction | git-native-projects | — |
| §47–48 | Pricing; value beside the invoice | principle | backstage Idea → `datatug/marketing/MONETIZATION.md` | — |
| §49 | Effectiveness metrics | measurable from the event stream | incidents (AC) | MVP (MTTR, TTRI raw), NEXT (reuse counters) |
| §50–51 | `datatug incident` CLI, `watch`, event stream | **new** | incidents; `datatug-cli` `cli/incident` | MVP |
| §52 | Autonomous agent workflow | **new** | agent-context | NEXT/FUTURE |
| §53–54 | Canonical demo incident + recurrence | **new** deliverable | incidents; `datatug-demo-projects` | MVP |
| §55–56 | GTM: lead with Incidentius; backbone ≠ entry point | product strategy | backstage Idea | — |

## 3. Conflicts found and how they were resolved

| Conflict | Resolution |
|---|---|
| Vision §1/§4 "one application/repository" vs the 2026-09-08 Dashboardius Idea (standalone repo, *"Build dashboards inside DataTug … Rejected"*) and the 2026-09-09 UI-kit ruling | Founder ruling 2026-09-11 (one Nx workspace) settles the repository; the UI-kit ruling settles the kits. Resolution: profiles for DataTug/Incidentius in one app, Dashboardius a sibling app in the same workspace, shared libs and an embeddable renderer. The Idea's repo argument is superseded; its schema argument stands. |
| Vision §4 "Dashboardius capabilities / dashboards" inside the MVP app vs hub REQ:no-charts-in-datatug-web-ui and Phase 5 not built | MVP incident metrics are numeric series tables in DataTug; boards and charts arrive with Track C (Phase 5 persistence + renderer library). The REQ is kept: no *independent* second renderer. |
| Vision §57 Q2 immutable events vs mutable state | Event-sourced with a derived projection (D3). Mutable state exists only as a cache of the fold. |
| Where incidents live: project files vs the platform's space-bounded default | The first draft made incidents project assets by analogy with boards. Founder ruling 2026-09-11 (*"Agree on ability of incidents to be stored in dedicated repo. But it should be able to store in same repo as datatug or even in the application repo if users want (can be OK for small projects)"*) makes the store location configurable: project repository by default, dedicated repository, application repository, or the hosted space-bounded store; one server-owned write path; incidents reference one or many projects (*"A DataTug project can be used by multiple incidentius instances"*); assets link back by store id + incident id. Execution records and snapshots follow the same routing. |
| Vision §59 "basic annotations/context linkage using existing capabilities" vs a full `annotations` Feature | No annotation capability exists to reuse, so the MVP slice of the new Feature (entity values and execution records) is the "basic" form; row sets, relationships and fields stay NEXT. |
| Vision §18 `Check` abstraction vs existing `checks.go` dataset checks | Extend, do not replace: the existing `Check{Type: options\|regexp\|sql}` becomes one evaluation kind; a check gains purpose, expected condition, scope, cost/risk, privilege and recorded runs. |
| Vision §15 snapshots vs `db snapshot` Idea (Approved, bulk copy to inGitDB) | Different things: a result snapshot is the recorded output of one execution; `db snapshot` copies tables. Both keep their names; evidence-records links to the Idea. |
| Vision §13 compare vs `pkg/comparator` | Different things: schema compare stays; a new recordset compare primitive is specified. |
| Vision §11 roles vs existing `Fact.origin` | Orthogonal fields: `origin` says where the fact came from; `role` says what cohort it belongs to. Both optional on the wire. |
| Vision §59 "AI-generated status summary if existing AI infrastructure makes it cheap" | No AI service exists in this family; `sneat-ai-backend` offers a pattern, not a service. Deterministic projection in MVP; AI narrative NEXT. |
| Vision §47 active-contributor pricing vs Dashboardius packaging Option A (viewing free) | Compatible: viewing is free in both; contributors are billed. Recorded in the Idea; marketing owns numbers. |

## 4. MVP boundary (what B1 proves)

Inside DataTug, on a local or Git project, under the single `serve` principal, hosted
multi-user deferred:

- incident store with configurable location (project repository by default), event stream,
  projection, list/detail/create ("Houston, we've got a problem"); incidents referencing one or
  many projects; back-links from queries, checks and boards by store id + incident id;
- every read of events, projection and watch stream re-checked against the reader's policy;
- lifecycle open → investigating → mitigating → recovering → resolved → watching → closed with
  only open/resolved/closed mandatory; the five outcomes; duplicates merge;
- canonical Investigation Context with cohort roles; hypothesis overlays with promote/reject;
- hypotheses with states and evidence links;
- execution records with fingerprints; opt-in result snapshots; re-checked on read;
- annotations on entity values and execution records;
- checks = library query + expected condition; manual run; promote query → check (needs the
  Phase 2 project-write path);
- metrics as check series; manual resolution criteria with evidence; numeric table;
- related/recurrence links; search by text and entity facts; deterministic similarity;
- `datatug incident …` CLI, `/datatug/incidents/*` API, `watch --json --since`, `show --at`;
- deterministic troubleshooter status projection;
- canonical demo incident fixture with a recurrence epilogue and a whole-journey e2e.

NEXT and FUTURE tiers are in §2 above and in the [roadmap](../../ROADMAP.md).

## Open Questions

None here; questions live in the owning Features and in
[product-family](../product/product-family.md#9-open-questions).
