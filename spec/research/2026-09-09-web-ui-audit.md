# DataTug Web UI audit (2026-09-09)

Repo: `/Users/alex/projects/datatug/datatug-apps`, commit `5d54b40` (2026-09-08, "enable enterprise OIDC SSO"; main advanced past this to that commit after the source `web-ui.md` report was written on an earlier commit). Method: source audit (routes, services, package.json, grep for `ion-*`/PrimeNG/Tabulator usage) plus a hands-on dev-server walk by the lead on 2026-09-09 — `pnpm nx serve datatug-app`, navigated by URL, read console/network state. Read-only throughout; nothing in the repo was modified.

---

## What works end to end today

- `pnpm nx build datatug-app` succeeds (warnings only); `pnpm nx test datatug-main` — 85 files / 115 tests, all pass (28.96% statement coverage, mostly 1-assertion smoke tests).
- Home page (`/`) renders live: store list, demo-project item, and an honest "no local agent detected at localhost:8989" info message (lead, 2026-09-09) — this matches the source finding that `AgentStateService` correctly detects a connection-refused agent.
- Entities pages (list/detail/edit) — real CRUD against the local agent or a GitHub raw-file fallback.
- Space Explorer (`/explore/:spaceId`) and Vault/inGitDB Explorer (`/explore-vault`) — real, read-only-by-construction transparency viewers.
- The SQL query editor (`sql-query-editor.component.ts`) and its execution pipeline through `AgentService`, when the agent's path prefix is fixed (see "What is broken").
- FK-driven related-record navigation in the table viewer (`env-db-table.page.html` "Record" tab) — real click-through (`goTable()`), contingent on the same agent connection.
- `/login` renders a real, populated auth surface — Sign up / Sign in / "Sign in with company SSO" / Quick login / Google / Apple / Microsoft / Facebook / GitHub (lead, 2026-09-09) — the new enterprise OIDC/SAML SSO landed 2026-09-08.

## What is hidden

- Widgets/Tags/Resources/Variables/DB-Models/Diff pages are one click away from the project menu but are totally empty — nothing signals to a user that they're unfinished before clicking in.
- `DbserverPageComponent` is fully coded but unreachable — its routing module is never imported anywhere.
- The GitHub OAuth device-flow code (`pkg/auth/ghauth` on the CLI side) has no equivalent web-visible entry beyond the login page button; depth not verified this pass.

## What is confusing

- Two different "servers" pages (`ServersPageComponent`) are reachable from two different routes but both dead-end on server-detail navigation.
- `routingParamServerType` is interpolated as the literal string `"serverType"` instead of `:serverType` in `datatug-routing-proj.ts:159` — the resulting route can never match a real generated URL.
- `/agent` just redirects to `/`, but the CLI's `datatug serve` auto-opens a browser to `<clientURL>/agent/<host:port>` — a URL shape this app has no route for at all.
- The home page has a dead "GitHub repositories" search box with a pre-filled value and no click handler.

## What is broken

- **The local-agent HTTP client is missing the CLI server's `/datatug` URL prefix on every call** — client calls `GET {agent}/agent-info`; the server (`datatug-cli/pkg/server/endpoints/routes.go:12-107`) only serves `GET {agent}/datatug/agent-info`. This is total and mechanical across 8+ call sites (`agent-state.service.ts:42`, `project.service.ts:122,173`, `agent.service.ts:48,83`, `project-item-service.ts:47,76,135,164,176,194`), not isolated typos — it is the single highest-value thing to fix first.
- **Firestore-backed writes call `sneat-go` endpoints that don't exist.** `createEnvironment()`, `createSchema()`, `createNewProject()` all POST to `https://api.sneat.ws/v0/datatug/...`; `sneat-go` has zero routes for any `datatug/...` path (only an unrelated Telegram bot profile matches "datatug" in that repo).
- **New-project creation is broken two independent ways**: the backend endpoint above doesn't exist, and the client hard-codes `userIDs: []` — even if the write succeeded, Firestore rules require the creator's uid in `userIDs` to read the project back, locking the creator out of their own new project.
- **Firestore-backed project summaries never resolve** — `ProjectService.watchProjectSummary()` for `storeId==='firestore'` subscribes to `firestoreChanges()`, which is a hard `throwError('Not implemented')`.
- **Clicking the demo-project link fails live** (lead, 2026-09-09): `No provider found for _DatatugFoldersService`, plus a 404 fetching `https://raw.githubusercontent.com/datatug/datatug-demo-project/main/datatug/datatug-project.json` — the real repo is `datatug-demo-projects` and the file is at `demo-project-1/datatug-project.json`. Two independent bugs (a missing DI provider and a wrong hard-coded URL) compound on the one path a new user is most likely to try first.
- Navigating to `/store/localhost:8989/project/demo-project-1/queries` renders the queries tabs shell (ACTIVE/BOOKMARKED/PERSONAL/SHARED, Tags, Order by) but throws `No provider found for _DatatugNavService`, an `ASSERTION ERROR` in `_EnvDbPageComponent`, and `unsupported format of store id:localhost:8989` in the console (lead, 2026-09-09).
- `/entities` and the table route render nothing but the app title under that same store-id path (lead, 2026-09-09) — consistent with the broken agent connection above.
- HTTP-type queries: `Coordinator.execute()` routes to `HttpExecutor.execute()`, which throws "Not implemented yet" — only the lower-level `executeCommand()` path works.
- `EnvironmentService.putConnection()` throws an empty-message error; `IndexedDbExecutor.execute()` throws "Not implemented yet" (appears entirely unwired elsewhere); `EntityService.deleteEntity()` has a `throw` inside an error-recovery branch.

## What is mocked

- Home page: a commented-out `<gui-grid>` bound to 3 fake people, from an abandoned `@generic-ui/ngx-grid` dependency.
- `env-db-table.page.html`'s "References" tab: hard-coded `SomeTable` row and a fake badge count "3"; empty Views/Stored-procedures cards.
- `diff-page.component.html`: literally `It's @if (true) { True }`.
- `table.service.ts:32-52`: hard-coded placeholder GitLab URL and a literal `'<TO_BE_SUPPLIEUD>>'` (sic) private-token string — never wired to anything real.
- `@sneat/datagrid`'s cell formatter: `formatterParams: { url: 'test-url' }` for any column ending in "Id"; its bundled `CellPopoverComponent` shows hard-coded "Id: 1" / "Code: Something" sample rows regardless of the actual cell clicked.

## Backend without coherent UI

- Board/dashboard data model (`Board`/`BoardRow`/`BoardCard`/`BoardWidget`) has full server CRUD routes; the web UI's `DatatugBoardService.getBoard()`-shaped calls throw "not implemented" in some paths and boards hard-code `localhost:8989` instead of the current store id.
- `RecordsetService` and `VariableService` are `@Injectable()` classes with zero methods, backing pages that render nothing.
- Folders/query-sets have full server CRUD (`routes.go`); UI is real for queries but the folder concept isn't separately exposed beyond the query tree.

## UI without real backing

- Widgets, Tags, Resources, Variables, DB-Models, Diff — six routed pages, each a header + empty `<ion-content>`, zero injected services beyond Ionic imports.
- `DbModelPageComponent` reads its data off `history.state['dbmodel']` (router navigation state) rather than fetching by ID — a direct link or page refresh shows nothing.

## Where the journey dead-ends

- **New user, home page → demo project**: dies immediately on a missing DI provider + wrong hard-coded GitHub URL (see "What is broken").
- **New user, home page → local agent**: `datatug serve` doesn't start (CLI-side bug, see companion current-state audit) — and even if it did, the client's missing `/datatug` prefix would 404 every call.
- **Signed-in user → New Project**: dead on both the missing backend endpoint and the `userIDs: []` client bug.
- **Any user → DB server detail**: `ServersPageComponent.goDbServer()` navigates to a URL nothing matches; `DbserverPageComponent` is orphaned dead code.

---

## Route and page matrix

Guard column: `SNEAT_AUTH_GUARDS` = Firebase-authenticated only; blank = unguarded.

| Route | Page component | Data source | State | Reachable from nav? | Hands-on observed (lead, 2026-09-09) | Notes |
|---|---|---|---|---|---|---|
| `/` | `DatatugHomePageComponent` | none (shell) + children | usable (shell) | yes — app root | **Yes** — stores list, demo-project item, "no local agent detected" message rendered | Dead hard-coded GitHub search box; commented-out mock grid |
| `/` (child `my-projects`) | `MyDatatugProjectsComponent` | `DatatugUserService` | partial | yes (embedded) | not separately observed | `my-projects.component.ts` file is empty; real logic elsewhere |
| `/` (child `my-stores`) | `MyStoresComponent` | `DatatugUserService` + `AgentStateService` | partial/mismatched | yes (embedded) | consistent with observed "no agent detected" message | Hard-codes `localhost:8989`; path-mismatched agent call |
| `/hello-world` | `HelloWorldPageComponent` | none | prototype | no | not observed | Nx scaffold leftover |
| `/debug` | `DebugPageComponent` | none | usable (as designed) | no (direct hit only) | not observed | Intentional Sentry smoke-test page |
| `/my` | `DatatugMyPageComponent` | none | usable | yes, via auth menu item | not observed | Entry point for the two explorers |
| `/explore/:spaceId` | `SpaceExplorerPageComponent` | Firestore, read-only | usable | yes, via `/my` | not observed | Space id typed by hand |
| `/explore-vault` | `VaultExplorerPageComponent` | GitHub via `@ingitdb/client-github` | usable | yes, via `/my` | not observed | Token kept in memory only |
| `/signed-out` | — | — | n/a | redirect to `/` | not observed | |
| `/store/:storeId` | `DatatugStorePageComponent` | agent HTTP or Firestore/GitHub | partial | yes, from My Stores | partially — see `/store/localhost:8989/...` rows below | No `SNEAT_AUTH_GUARDS` on this subtree |
| `.../diff` | `DiffPageComponent` | none | prototype (empty) | via store nav | not observed | Literal placeholder text |
| `.../environment` | `EnvironmentPageComponent` | agent/Firestore | partial | yes | not observed | |
| `.../project/:projectId` | `ProjectPageComponent` | `ProjectService.watchProjectSummary` | partial/broken (Firestore) | yes | not observed directly | Firestore path never resolves |
| `.../board/:board` | `BoardPageComponent` | agent, hard-coded `localhost:8989` | partial | via `/boards` | not observed | Only works against default local agent port |
| `.../boards` | `BoardsPageComponent` | `DatatugBoardService` + `DatatugFoldersService` | usable UI / backend partial | yes | not observed | `getBoard()`-shaped calls throw not-implemented in some paths |
| `.../entities` | `EntitiesPageComponent` | `EntityService` / GitHub raw | usable | yes | **Yes, under `/store/localhost:8989/project/demo-project-1/`** — rendered nothing but the app title | Real CRUD per source, but observed empty live under this store id |
| `.../new-entity` | `EntityEditPageComponent` | `EntityService` | partial | yes | not observed | |
| `.../entity/:entityId` | `EntityPageComponent` | `EntityService` | usable | yes | not observed | |
| `.../environments` | `EnvironmentsPageComponent` | `EnvironmentService` | partial | yes | not observed | |
| `.../widgets` | `WidgetsPageComponent` | none | prototype (empty shell) | yes | **Yes** — rendered an empty "Widgets" shell | Confirms source finding live |
| `.../tags` | `TagsPageComponent` | none | prototype (empty shell) | yes | not observed | |
| `.../resources` | `ResourcesPageComponent` | none | prototype (empty shell) | yes | not observed | |
| `.../variables` | `VariablesPageComponent` | `VariableService` (empty class) | prototype (empty shell) | yes | not observed | |
| `.../dbmodel/:dbModelId` | `DbModelPageComponent` | `history.state['dbmodel']` | fragile/partial | via `/dbmodels` | not observed | No fallback fetch by id |
| `.../dbmodels` | `DbModelsPageComponent` | none | prototype (empty shell) | yes | not observed | |
| `.../env/:environmentId` | `EnvironmentPageComponent` | `EnvironmentService`/`ProjectService.getFull` | partial | yes | not observed | |
| `.../db/:dbCatalogId` | `EnvDbPageComponent` | `ProjectService.getFull()` | partial | yes | **Yes** — `ASSERTION ERROR: NgModule '_EnvDbPageComponent'...` thrown navigating this path | `loadChildren` points at a plain component class, not a `Routes` array |
| `.../table/:tableType` (also `view/:tableType`) | `EnvDbTablePageComponent` | `AgentService.select()` (path-mismatched) | usable UI / broken data | yes | **Yes** — table route rendered nothing but the app title | Real Tabulator grid + FK tabs; blocked by the agent path bug |
| `.../servers` | `ServersPageComponent` | `DbServerService` | usable list / dead-end detail | yes | not observed | `goDbServer()` targets a URL nothing matches |
| `.../query/:queryId` | `QueryPageComponent` | `QueriesService`, `Coordinator` | usable UI / partial execution | yes | not observed | HTTP-type queries throw not-implemented |
| `.../queries` | `QueriesPageComponent` | `QueriesService` (agent) | usable | yes | **Yes** — ACTIVE/BOOKMARKED/PERSONAL/SHARED tabs, Tags, Order by all rendered | Also threw `No provider found for _DatatugNavService`, `unsupported format of store id` in console |
| `.../server/serverType/:dbDriver/:dbServerId/:dbCatalogId` | → `EnvDbTablePageComponent` | — | likely broken route | not reachable by normal nav | not observed | Literal `"serverType"` segment, should be `:serverType` |
| — (orphaned) | `DbserverPageComponent` | not reached | orphaned, unreachable | no | not observed | Routing module never imported |
| `/store/:storeId/project` | — | — | n/a | redirect to `''` | not observed | |
| `/agent` | — | — | n/a | redirect to `/` | not observed | CLI's `datatug serve` auto-opens a URL shape this app has no route for |
| `/login`, auth routes | `@sneat/app-auth`/`@sneat/auth-ui` (external lib) | Firebase Auth | usable | yes | **Yes** — Sign up / Sign in / "Sign in with company SSO" / Quick login / Google / Apple / Microsoft / Facebook / GitHub all rendered | New enterprise OIDC/SAML SSO (main `5d54b40`, 2026-09-08); depth beyond render not verified |

---

## Core concepts vs current UI

Phase-1 relevance judged against the core loop: project → source → table/query → grid → cell selection → related/applicable → parameterized run → context → ACL-aware result.

| Concept | State | Evidence | Phase-1 relevance |
|---|---|---|---|
| Project/source navigation | usable | `ProjectTracker`/`StoreTracker`, `DatatugNavContextService` | **High** — first step of the core loop |
| Environments | partial | `EnvironmentService.getEnvSummary()` path-mismatched; `putConnection()` throws | Medium — needed for real source connection, not the demo path |
| DB schemas/tables | usable UI / broken data fetch | Real Tabulator grid + FK/keys/indexes; `AgentService.select()` mismatched | **High** — table/query and grid steps |
| Entities | usable | Full CRUD, GitHub fallback | Medium — semantic layer, not blocking the core loop |
| Queries (list/editor/execution/params/folders) | partial | Real editor + folders; SQL path mismatched, HTTP path unimplemented | **High** — parameterized-run step |
| Result grid | usable component, thin config | `@sneat/datagrid`/Tabulator; no cell-range selection exposed | **High** — grid + cell-selection steps |
| Related records / related-data panel | usable | FK-driven `sneat-datatug-fk-card` + `referencedBy`, click-through | **High** — related/applicable step |
| Context / basket / variables page | prototype (empty) | `VariablesPageComponent` — empty shell | **High** — context step has no UI at all today |
| Widgets / dashboards (boards) | split: boards=usable, widgets-page=empty | Real board system; unrelated empty "Widgets" nav page | Low — not on the core exploration loop |
| Diff page | prototype (empty) | Placeholder text | Low |
| Servers/dbserver pages | partial, dead-end detail | List works; detail navigation broken | Medium — source-connection adjacent |
| Resources | prototype (empty) | Header + empty content | Low |
| Store (top-level per-store page) | partial | Own "Not implemented yet" throw | Medium |
| Tags | prototype (empty) | Header + empty content | Low |
| Space-explorer | usable | Read-only Firestore browser | Low — transparency tool, not the core loop |
| Vault explorer | usable | Read-only GitHub/inGitDB browser | Low |
| AI/query-builder | specified/documented only, zero code | 3 specs Approved; zero matching code | Medium — would sit on top of the parameterized-run step once built |
| Auth/login (providers incl. new SSO) | usable (external lib) | e2e-verified render; SSO landed 2026-09-08 | Medium — gates everything, but not itself the loop |
| Menus/side-nav | usable, partial reach | Auth item, store/project/env selectors | High — the connective tissue of the whole loop |
| Empty/error states | mixed | Deliberate ones exist; several pages are accidentally always-empty | Medium |

---

## UI kit coupling (Ionic vs PrimeNG)

**Package facts** (`datatug-apps/package.json`, verified by reading, 2026-09-09): `@ionic/angular` 9.0.0, `@ionic/core` 9.0.0, `@ionic-native/*` 5.36.0, `@sneat/datagrid` 0.27.6, `@sneat/grid` 0.27.6, `tabulator-tables` 6.3.1. **No PrimeNG, no AG Grid, no `@primeuix/*` dependency exists anywhere in this repo** (`grep -riE '"(primeng|ag-grid|@primeuix)'` across `package.json` → 0 hits). `@sneat/datagrid` (Tabulator-based) is imported in exactly two files: `pages/signed-in/env-db-table/env-db-table.page.ts` and `pages/signed-in/env-db/env-db-page.component.ts`.

**`ion-*` element usage per page directory** (`grep -roE '<ion-[a-z-]+' <dir> | wc -l`, `libs/datatug/main/src/lib/pages/**`):

| Page directory | `<ion-*>` tag count | HTML files |
|---|---|---|
| `signed-in/env-db-table` | 106 | 5 |
| `signed-in/vault-explorer` | 71 | 1 |
| `signed-in/db-schema` | 64 | 1 |
| `signed-in/entity-edit` | 53 | 2 |
| `signed-in/dbserver` | 52 | 1 |
| `signed-in/entities` | 45 | 1 |
| `signed-in/space-explorer` | 44 | 1 |
| `signed-in/servers` | 39 | 1 |
| `signed-in/store` | 33 | 1 |
| `signed-in/entity` | 32 | 1 |
| `signed-in/environment` | 28 | 1 |
| `board/ui/pages/board` | 27 | 1 |
| `signed-in/project` | 21 | 1 |
| `board/ui/pages/boards` | 17 | 1 |
| `signed-in/env-db` | 18 | 1 |
| `signed-in/environments` | 11 | 1 |
| `signed-in/db-schemas` | 7 | 1 |
| `signed-in/resources`, `widgets`, `variables`, `tags`, `diff` | 4-6 each | 1 each |
| `queries/query` | 296 | 5 |
| `queries/queries` | 99 | 2 |
| top-level `home` | 67 | 4 |
| top-level `my` | 19 | 1 |

Every page in the app is built from native Ionic elements; zero PrimeNG markup exists anywhere in the current codebase.

**Migration cost, if the shell moved to PrimeNG (facts, not a recommendation):** Per the fleet standard (`sneat-co/sneat-specs/standards/frontend-ux/ionic-and-primeng.md`, founder decision 2026-08-19), Ionic is primary fleet-wide and PrimeNG is admitted only as bounded islands for data surfaces Ionic has no component for — page chrome, navigation, overlays, forms, and the submit affordance must always stay Ionic. Under that standard, a wholesale shell migration to PrimeNG is out of policy; the realistic path is islands. The heaviest, most PrimeNG-plausible surfaces by `ion-*` density above are `env-db-table` (106 tags across 5 templates — the result-grid page, already backed by Tabulator via `@sneat/datagrid`, which the standard names as the preferred grid) and the query editor (`queries/query`, 296 tags across 5 templates — almost certainly page chrome + forms, which the standard reserves for Ionic regardless). Everything under `signed-in/*` that is currently an empty shell (Widgets/Tags/Resources/Variables/DB-Models/Diff) is cheap to rebuild in either library, since there is no existing logic to port. Dashboardius (`sneat-co/dashboardius`) is a working reference for what a PrimeNG-primary data surface with a switchable grid engine looks like in this same ecosystem, if an island approach is chosen for a future data-grid-heavy DataTug screen. This audit does not make the shell-vs-island call — that is a founder decision per the fleet standard.

---

## Time to first useful result

| Stage | Status today | Blocking defect(s) |
|---|---|---|
| First connected source (web) | Broken | No `datatug serve` reachable path works (server-side panics, client-side missing `/datatug` prefix); Firestore-backed project creation is broken two ways; demo-project link 404s on a missing DI provider and a wrong hard-coded GitHub URL |
| First query/result (web) | Broken | Same agent-path mismatch blocks the query editor's SQL execution path; HTTP-query path throws not-implemented regardless |
| First related-data discovery (web) | Blocked by data fetch, UI is real | `env-db-table` FK/referenced-by cards are real and click-through, but the underlying table-select call is agent-path-mismatched |
| First shared-query reuse (web) | UI shell renders, no data | Queries page tabs (Active/Bookmarked/Personal/Shared) render live but the query list itself depends on the broken agent connection |
| Understanding current context (web) | Partial | Nav-context services exist but threw `No provider found for _DatatugNavService` and an `NgModule` assertion error during the live walk on the env/db route |

---

## Reuse pointers

- **Local-agent client**: `services/repo/agent.service.ts` (select/execute), `agent-state.service.ts` (polling/health), `store-api.service.ts` (generic GET/POST/PUT/DELETE). Fix the `/datatug` path prefix here first (`pkg/server/endpoints/routes.go:12-107` on the CLI side) — this one change likely unblocks most agent-backed features simultaneously.
- **Result grid with cell selection**: extend `@sneat/datagrid`'s `DataGridComponent` (Tabulator-based, `SelectRowModule`/`InteractionModule`/`MenuModule` already registered). Cell-range selection needs Tabulator's `range` module registered in the published `@sneat/datagrid` package (not this repo). The FK-popover `cellFormatter`/`CellPopoverComponent` pattern in `env-db-table.page.ts` is the template for new cell-level interactions.
- **Related-data/context panel**: `pages/signed-in/env-db-table/foreign-key-card/`, `ref-by-card/`, `record-values-card/` — a working, real FK-driven navigation pattern to replicate for any other "related records" surface.
- **Parameterized query invocation**: `models/definition/parameter.ts`, `parameter-lookup.service.ts`, `input-parameters.component.ts` — already used by `board-page.component.ts` and `sql-query-editor.component.ts`.
- **Query execution pipeline seam**: `executor/coordinator.ts` dispatches to `AgentService` or `HttpExecutor` by command shape — a reasonable seam to add new executor types without touching call sites.
- **Nav/context state**: `services/nav/datatug-nav-context.service.ts`, `services/nav/datatug-nav.service.ts`, `services/nav/contexts/project.tracker.ts`/`store.tracker.ts` — fix the missing DI providers surfaced by the live walk before building new context UI on top.

---

## Open Questions

1. Is `localhost:8989` still the intended CLI-agent integration point, or has the web↔CLI protocol moved since these client paths were written? The `/datatug`-prefix mismatch is total and mechanical across 8+ call sites, which reads more like the CLI's route prefix changed after the web client was last touched than like independent typos.
2. Is there an actual plan for a `datatug`-specific `sneat-go` module (writes for environments/schemas/projects), or should Firestore-backed projects be read-only/agent-only for now?
3. Should `datatug_projects` move to the space-bounded model, given `ICreateProjectData.spaceID` already exists unused and `@sneat/space-services`/`@sneat/space-models` are already dependencies?
4. Are Widgets/Tags/Resources/Variables/DB-Models/Diff genuinely next on a roadmap, or should they be hidden from nav until built?
5. Is `DbserverPageComponent` intentionally retired, or meant to be reconnected — and if reconnected, does the `server/${routingParamServerType}/...` route need the `:` fix?
6. Has project-creation ever been exercised end-to-end since the `sneat-apps` migration, given the `userIDs: []` bug, or has all real usage been agent/GitHub-store-only?
7. Now that SSO landed on 2026-09-08 (`5d54b40`), is enterprise OIDC/SAML in scope for Phase-1 at all, or is it unrelated forward work that should be left alone while Phase-1 focuses on the core exploration loop?
8. If a data-grid-heavy screen ever needs a PrimeNG island (per the fleet's Ionic-primary/PrimeNG-island standard), should it follow Dashboardius's switchable-grid-engine pattern, or a purpose-built `@sneat/datagrid` extension?
