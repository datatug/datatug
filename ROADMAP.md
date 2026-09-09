# DataTug roadmap — fast start, quick wins, compounding capability

**Status:** Draft for founder review · **Date:** 2026-09-09 · **Owner:** alex
**Reads with:** [product thesis](spec/product/README.md) · [feature hierarchy](spec/product/hierarchy.md) · [current-state audit](spec/research/2026-09-09-current-state-audit.md) · [Phase 1 Feature](spec/features/core-investigation-loop/README.md) · [Phase 1 plan](spec/plans/2026-09-09-phase-1-core-investigation-loop.md) · [competitor research](https://github.com/datatug/marketing/tree/main/competitors)

The roadmap optimises for **product value per unit of implementation complexity** and
for **visible progress an autonomous coding agent can deliver without product
decisions**. Phases 0–2 are almost entirely completion of partial work already in the
repositories. Every item traces to a pillar in the hierarchy; nothing is here because a
competitor has it.

## The ladder

| Rung | What it holds | Phases |
|---|---|---|
| **No-brainers** | Broken things in the spine that block every demo: the CLI build, `datatug serve`, the web↔agent path prefix, the demo project link, hiding empty pages, secrets guard. | 0 |
| **Quick differentiators** | The core investigation loop end to end on one demo project with two sources: semantic mapping, cell → context panel, applicable queries with auto-bound parameters, Investigation Context, related lookups, server-enforced ACL. | 1 |
| **Strategic foundations** | Library depth (collections, capture, explainability), provenance, HTTP queries, agent-readable context, CLI/web hand-off registry, TUI parity for the loop. | 2 |
| **Larger strategic bets** | Playbooks, sharing scopes on Sneat spaces, hosted team workspace, board persistence + Dashboardius round trip, UI-kit consolidation. | 3 |
| **Bets that need other tracks** | Layered multi-owner ACL with Explain Access (DALgo plan task 21), serve-brokered query builder / MCP. | 4 |
| **Experiments** | Automatic relationship discovery by value sampling, probabilistic meaning suggestions, AI-assisted playbook capture, public knowledge packs, in-app AI builder. | 5 / X |
| **Explicitly not building** | Workflow engine, app builder, catalog/governance suite, warehouse, connector race, client-direct DB access, disconnected text-to-SQL, charts inside DataTug's own UI. | — |

## Phase 0 — Repair the spine (no-brainers)

- **Objective.** A fresh clone builds, `datatug serve` starts, the web UI talks to it, and the demo project opens. Nothing new is designed.
- **User outcome.** `datatug demo` → browser opens the demo project → a table renders rows. Today none of these three steps works.
- **Why now.** Every later phase runs through `serve`; three "first contact" commands (`init`, `demo`, `serve`) currently panic and `main` does not compile (go-github v90→v91 bump without import rewrite).
- **Included.** (1) Fix imports for go-github v91 and add a build gate that fails a dependency PR when it breaks `go build ./...`. (2) Wire `storage.NewDatatugStore`/`GetStore` to the working `filestore` path already used by `project_base_command.go`; pass a real handler to `RegisterDatatugHandlers`; nil-check `config.Server`; bind `--host/--port/--project` flags on `serve`. (3) Add the `/datatug` prefix (and the real endpoint names) to the web agent client; accept `localhost:8989`-style store ids. (4) Point the demo link at `datatug/datatug-demo-projects/main/demo-project-1/datatug-project.json`; provide `DatatugFoldersService`. (5) Hide the six empty shell pages (Widgets, Tags, Resources, Variables, DB models, Diff) behind a feature flag. (6) A Playwright e2e that boots the real agent (no `localhost:8989` interception) and asserts one grid renders. (7) Validator rule: a `QueryDefTarget` with credentials fails `datatug validate`. (8) Use `dalgo2sqlite` (pure Go) for structured queries so release binaries keep SQLite without CGO.
- **Excluded.** Any new concept. TUI changes. Fixing the six `show`/`dataset*`/`queries` unbound-flag commands (deprecate or fix in Phase 2).
- **Hygiene (founder rulings 2026-09-09).** datatug.io keeps the ahead-of-code claims (MCP, `snapshot`/`diff`/`schema`/`transform`, DTQL repos) *"with planned label"*; the site and cask licence text changes to Apache-2.0 for the CLI; the web UI (`datatug-apps`) states its licence position per the founder: *"Web UI private for now, will be AGPL or something like that"* (the repository is public today with no LICENSE file — whether "private" means the repository visibility or only the licence is still to confirm).
- **Dependencies.** None.
- **Risks.** `pkg/api`'s storage indirection may be simpler to delete than to fix; the implementer decides and records it.
- **Impact.** Removes the single biggest credibility gap (nothing demoable). No competitive movement by itself.
- **Complexity.** Low–medium. **Autonomy.** Very high — no product decisions.
- **Acceptance.** `go build ./... && go test ./...` green on `main`; `datatug serve` serves `/datatug/ping`; e2e boots agent + web and sees Chinook rows; `datatug demo` completes; `datatug validate` rejects committed credentials.
- **Test strategy.** Unit tests for the store wiring (replace the `recover()`-around-the-panic test); e2e against the real agent.
- **Demo.** `brew install datatug && datatug demo` opens the browser on the demo project's Album table.

## Phase 1 — The core investigation loop (quick differentiators)

Specified in full in the [Phase 1 Feature](spec/features/core-investigation-loop/README.md) and [plan](spec/plans/2026-09-09-phase-1-core-investigation-loop.md).

- **Objective.** Make DataTug's defining concepts work at a basic, coherent, demonstrable level in the web UI, backed by the server for execution and access control.
- **User outcome.** Open the demo project → open a table or query → see rows in a dense grid with semantic markers on mapped columns → select a `CustomerId` cell → the context panel shows *Customer.ID = 5*, related records (invoices from SQLite, support notes from an inGitDB collection, country facts from a public HTTP endpoint), and applicable library queries with parameters already bound → run one through the server → add the value to the Investigation Context → open another query and see the parameter bound from context → switch to a restricted principal and see rows filtered and a column hidden, with the limitation explained.
- **Why now.** Every piece has a partial implementation: `EntityFieldRef` tags, `NamePatterns`, the FK graph, typed parameters, `input-parameters`, the Tabulator grid with FK cards, `accesspolicies.Run`. Wiring them is cheaper than any new subsystem and it is the demo that shows *why DataTug is different*.
- **Included.** Field→column mapping resolver (explicit mappings + `NamePatterns`); semantic markers in the grid; context panel; related lookups via server-built DALgo queries; applicable-query resolution with explanation; parameter auto-binding; Investigation Context (add/remove/disable, visible bindings); `serve` running every read through `accesspolicies.Run` with a project policy set and a chosen principal; DTQL saved-query type; demo project extended (Customer/Invoice entities, semantic fields, three parameterized queries, an inGitDB `support-notes` collection, a policy file with two principals); journey e2e.
- **Excluded.** Playbooks, collections semantics, HTTP execution beyond read-only JSON GET reference data, provenance UI beyond ACL limitations, inference beyond deterministic patterns, cross-source joins, Sneat scopes, MCP, TUI parity, any UI-kit migration.
- **Dependencies.** Phase 0.
- **Risks.** Opaque SQL text cannot be row/column-filtered by DALgo → Phase 1 states the limitation explicitly (source-level ACL + read-only session for SQL text; row/column ACL for structured/DTQL queries and for all server-built lookups). Ionic density for the context panel → keep it a plain list first.
- **Impact.** After Phase 1 a knowledgeable observer can see the semantic + contextual + reusable-knowledge combination that no researched competitor offers in one place.
- **Sources (founder ruling 2026-09-09, verbatim).** *"I want to have in MVP: inGitDB, SQLite, HTTP endpoint. Regards HTTP - some public sources like countries or currency rates or something that is common and is understandable and can looks good in demos."* Phase 1 therefore carries three sources; read-only HTTP JSON execution moves from Phase 2 into Phase 1 (plan task 11) through a generic DALgo adapter, `dal-go/dalgo2http` (founder 2026-09-09, verbatim: *"dalgo2http - yes"*; repository bootstrapped the same day), and the concrete public endpoints are the lead's assumption in the Feature.
- **Complexity.** Medium. **Autonomy.** High — the Feature fixes every product choice; the remaining lead assumptions are flagged as the lead's, not the founder's.
- **Acceptance / tests / demo.** See the Feature's acceptance criteria and the plan's journey e2e.

## Phase 2 — Library depth and provenance (strategic foundations)

- **Objective.** Turn the loop into a compounding library: capture from exploration, collections with meaning, explainable suggestions, provenance, HTTP queries, and the agent reading the same context.
- **User outcome.** "Save this lookup as a team query", "why was this suggested", "where did this value come from", "run the shipping-status HTTP request for this order", "`datatug context` prints what the agent should know".
- **Why now.** These are the second half of the explore → capture → share loop, and they make the AI story honest before the MCP face exists.
- **Included.** Capture-from-exploration; collection README + domain tags; explanation chain UI; provenance record on every result (source, collection, query, policy limitations); HTTP `QueryDef` execution in CLI and web; `datatug context` CLI verb (agent-readable entities, mappings, applicable queries for a value); Ctrl+W registry covering the loop's screens; TUI "Queries" panel bound to real data; fix-or-delete the unbound-flag commands; environment/time semantic fields for auto-narrowing.
- **Excluded.** Playbooks, scopes, boards, layered ACL.
- **Dependencies.** Phase 1.
- **Risks.** HTTP execution needs secret handling (Phase 0's guard + `~/.datatug` credentials).
- **Impact.** Closes the "shared query library" baseline gap against DataGrip/DBeaver/Navicat while keeping the semantic edge.
- **Complexity.** Medium. **Autonomy.** High.
- **Acceptance.** A captured query appears in Git as a diff; provenance visible on a result; HTTP query runs from the demo; `datatug context --value Customer.ID=5` lists applicable queries.

## Phase 3 — Playbooks, scopes, dashboards (larger strategic bets)

- **Objective.** Knowledge that is organised, shared and presented: investigation playbooks, Sneat-space sharing scopes, board persistence with the Dashboardius round trip, hosted team workspace, UI-kit consolidation.
- **User outcome.** Follow "Checkout failure investigation" from an Order ID with branches; promote a query from personal to team; pin a result to a board and open it on dashboardius.com.
- **Why now.** Requires a working loop and library; also the first monetizable surface (hosted workspace).
- **Included.** Playbook model (entry points, steps, branches, HTTP steps) + guided runner; project store backends — a remote **OpenVaultDB** vault as a first-class project store (read and write through the OVDB API, OVDB's own policies composing with DataTug's) and the **project store explorable inside DataTug** (entities, queries, mappings and policies browsable as a source through the same loop; founder additions 2026-09-09); scope model on Sneat spaces (space membership → project access; `datatug_projects` under `/spaces/{spaceID}/ext/datatug/` per the space-bounded rule — founder decision required); promotion lifecycle; board persistence in `serve` + `@datatug/board-models` package + Dashboardius query binding; execute the UI-kit decision (see below).
- **Excluded.** Layered ACL, MCP.
- **Dependencies.** Phase 2; scopes and board ownership are now founder-ruled (see Open Questions 5–6, resolved 2026-09-09); for the OpenVaultDB project store, the `layered-acl-query` branches landing (OVDB HTTP query integration is task 5 of that plan, already complete; the codex agent on the VM is driving the plan to completion per the founder's 2026-09-09 ruling).
- **Risks.** Playbooks drifting toward a workflow engine — the Feature must forbid automation semantics.
- **Impact.** The moat starts compounding; Dashboardius gets real data.
- **Complexity.** High. **Autonomy.** Medium — needs founder decisions listed in the Open Questions.

## Phase 4 — Layered access control and agents

- **Objective.** DataTug consumes the layered-ACL MVP (DALgo plan task 21: query/write + Explain Access in DataTug) and ships the serve-brokered query builder's MCP face.
- **Dependencies.** The DALgo/OpenVaultDB/InGitDB `layered-acl-query` branches (tasks 11–20) landing; Phase 0's working `serve`.
- **Autonomy.** Medium — contracts frozen by the ACL packet; implementation follows its work packages.

## Phase 5 / experiments

- Automatic relationship discovery by value sampling (experiment: precision on the demo project before any UI).
- Probabilistic meaning suggestions (types, neighbours, sampling) with review UI.
- Infrastructure inventory: services/endpoints, record → infrastructure climb.
- Cross-source joins on the server (only if a real investigation needs a join rather than a lookup).
- AI-assisted playbook capture from investigation history; public knowledge packs; in-app AI builder (behind the serve-brokered path).

## Dashboardius and DataTug

Founder direction (2026-09-08 Idea, verbatim): *"A Dashboardius board IS a DataTug board; what Dashboardius adds is widget kinds and an editing experience."* This roadmap keeps that split:

- **DataTug owns** the board schema, board persistence in the project (Git file), query execution and access control, and the "pin this result to a board" entry point from exploration (Phase 3).
- **Dashboardius owns** presentation and editing (widgets, layout, charts, sharing links) and consumes DataTug through the server API: list projects/queries, execute a bound query with parameters, read/write a board. It never queries a database directly.
- **Contract work** (Phase 3): publish `@datatug/board-models`; add `GET/PUT /datatug/boards/*` + execute endpoints with ACL to `serve`/hosted API; board ownership is founder-ruled (2026-09-09, verbatim): *"board belongs to DataTug project. Project can be shared across Sneat spaces"*.
- **Selling it separately** is compatible with this split: Dashboardius is the presentation product with its own domain and pricing; DataTug remains free and open source (marketing/MONETIZATION.md). The competitor research records how Metabase, Superset/Preset, Lightdash and Evidence price viewers and embedding — viewer seats and embedding are where those products charge; a DataTug-backed board with link sharing for signed-in colleagues is the cheapest "share a live board" path found.
- **Guard rail.** Phase 1 does not build boards. Dashboardius work must not pull chart or layout work into DataTug's web UI.

## UI kit decision

**Founder leaning (2026-09-09, verbatim):** *"We probably should PrimeNG as base without ionic only if we think ionic is limiting us. Otherwise we probably should use ionic as app frame and base components and use PrimeNG or other libs like ag grid where ionic lack components/functionality."*

**Confirmed by the founder 2026-09-09, verbatim:** *"UI-kit rule - confirmed"* — Ionic frame + islands is the standing rule for DataTug's web UI.

This matches the fleet standard (`sneat-specs/standards/frontend-ux/ionic-and-primeng.md`, founder decision 2026-08-19): Ionic is primary; PrimeNG only inside bounded islands for data surfaces Ionic lacks; `@sneat/datagrid` (Tabulator) before `p-table`.

**Assessment — is Ionic limiting DataTug?** Evidence from the Web UI audit:

| Need on the core loop | Ionic has it? | What DataTug uses / would use |
|---|---|---|
| App frame, menu, header, tabs, lists, modals, popovers, toasts | Yes | Ionic (already in place) |
| Dense result grid with cell selection, sort, filter | No | `@sneat/datagrid` (Tabulator) — already used by the table page; add Tabulator's range/cell modules in the published package |
| Schema/collection tree | No | PrimeNG `p-tree` island |
| Resizable side panel for the context panel | No (layout) | CSS/CDK resize inside `ion-split-pane`; no library needed |
| Charts | No | Not in DataTug (Dashboardius) |
| Keyboard command palette | No | Custom component |

Nothing on the loop needs a component Ionic lacks *and* PrimeNG islands cannot provide. The Ionic shell in `datatug-apps` is thin (menu + headers + lists); the risk the standard names — mixing overlays — does not arise because islands never own overlays. **Recommendation:** keep Ionic as the frame, adopt the island rule for the grid and tree in Phase 1, and revisit only if a Phase-2/3 screen (playbook runner, board editor) needs desktop density Ionic cannot give. Dashboardius stays a standalone PrimeNG app (already shipped); the two products share the board schema and, later, published islands (`@sneat/prime-*` in `sneat-libs/libs/prime/*` per the standard) rather than a shell. This is an assumption by the lead until the founder confirms.

## Open Questions

Numbered for the founder; answers unblock the phases named.

1. ~~Confirm the Phase 1 scope as specified.~~ Ruled 2026-09-09 (verbatim): *"I approve all artefacts - including Phase 1 + addition of HTTP data sources"*. Feature `core-investigation-loop` and its plan are Approved; implementation streams started the same day (Phase 0 spine/web/demo/site, Phase 1 model). Lead assumptions A2–A3 and the concrete HTTP endpoints stand unless the founder objects. — Phase 1
2. ~~Second source for the demo: inGitDB collection or HTTP endpoint?~~ Ruled 2026-09-09 (verbatim): *"I want to have in MVP: inGitDB, SQLite, HTTP endpoint. Regards HTTP - some public sources like countries or currency rates or something that is common and is understandable and can looks good in demos."* Recorded in the Phase 1 Feature and plan (task 11). — Phase 1
3. ~~UI kit: confirm "Ionic frame + islands".~~ Ruled 2026-09-09 (verbatim): *"UI-kit rule - confirmed"*. — Phase 1
4. ~~Layered-ACL status: which record is true?~~ Ruled 2026-09-09 (verbatim): *"It was authorised and is being implement on VM by codex agent. I've asked it to drive it to completion."* The dalgo plan (`spec/plans/layered-acl-mvp.md`, Executing, tasks 1–10 complete) is the current record; the design packet's "no implementation authorized" line is superseded. Branches are on origin (see the audit's update note). — Phase 4 dependency, in progress
5. ~~Board ownership: DataTug project or Sneat space?~~ Ruled 2026-09-09 (verbatim): *"board belongs to DataTug project. Project can be shared across Sneat spaces"*. Recorded in Feature `dashboards`. — Phase 3
6. ~~Scopes: space-bounded or project-level?~~ Ruled 2026-09-09 (verbatim): *"DataTug project is standalone entity that can be shared (granted) between multiple spaces (teams)"*. This is an explicit exception to the space-bounded default; recorded in Feature `sharing-scopes`. — Phase 3
7. ~~Licence: Apache-2.0 or MIT?~~ Ruled 2026-09-09 (verbatim): *"CLI - Apache 2.0. Web UI private for now, will be AGPL or something like that"*. Still to confirm: whether "private" means making `datatug-apps` a private repository or only withholding an open-source licence (it is public today with no LICENSE file). — Phase 0 hygiene
8. ~~Marketing claims ahead of code: retract or label?~~ Ruled 2026-09-09 (verbatim): *"keep with planned label"*. — Phase 0 hygiene
