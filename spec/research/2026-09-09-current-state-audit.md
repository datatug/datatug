# DataTug current-state audit (2026-09-09)

Method: read-only archaeology of five prior audit passes (`datatug-cli`, `datatug-core` + demo/validator repos, `datatug-apps`, `dal-go/dalgo` + adapters + DTQL + ACL, ecosystem/marketing/Sneat) plus a lead-verified sixth pass — building the CLI at tag `v0.13.0`, running `datatug serve`/`demo`/`init` live, and a hands-on browser walk of `datatug-apps` on 2026-09-08/09 — consolidated and reconciled here. Every claim below cites a `path:line`, a command actually run, or a URL actually visited; nothing is invented. Classification vocabulary (one per capability row): **usable** (works as built, end to end) · **implemented-but-poorly-exposed** (real, but no UI/CLI path reaches it, or the path is broken) · **partial** (some paths/backends/scopes work, others don't) · **prototype** (code exists, does not function or is a stub) · **backend-only** (server/library exists, no consuming UI) · **ui-only/mock** (UI exists, no real data behind it) · **specified/documented** (spec/design exists, zero shipped code) · **planned** (roadmap intent, nothing built) · **obsolete** (superseded or dead) · **unclear** (not resolved by evidence gathered).

Facts newly established by the lead after the five source reports were written are folded in and, where they conflict with a source report, they override it — flagged inline as "(lead, 2026-09-09)".

---

> **Update, later on 2026-09-09:** the layered-ACL branches named below as `vm/*` refs were pushed to origin on the founder's approval (`datatug/dtql` `layered-acl-design` c2d0487; `layered-acl-query` in dal-go/dalgo 65c4824, dalgo2sql 45fe36c, dalgo2ingitdb de6f64c, openvaultdb-go 7998895). The `vm/*` names remain valid as the local evidence refs this audit read. The design-packet-versus-plan status question is carried in the hub `ROADMAP.md`, open question 4.

## Verdict in ten lines

1. **The single most functional, most secured path in the whole ecosystem is `datatug query run`**: DTQL parsing (`dal-go/dalgo/dtql`) → `pkg/accesspolicies` → `access.SecureReadSession` (`dal-go/dalgo/access`), live-verified against a real SQLite Chinook DB with both `--no-policies` and policy-enforced runs (`datatug-cli/apps/datatugapp/commands/cmd_query.go`). Nothing else in the stack is this solid.
2. **Every "first contact" path a new user would try is broken.** `datatug init` panics unconditionally (`pkg/datatug-core/storage/vars.go:19`); `datatug demo` downloads real data then fails to save the project (`cmd_demo.go:308-344`); `datatug serve` panics before it opens a socket (`cmd_serve.go:43`, confirmed live by the lead on `v0.13.0`), and even past that panic, route registration itself panics (`handler is not provided`, `pkg/server/endpoints/register.go:31`, because `http_server.go` passes a nil handler) and `storage.NewDatatugStore` is a closure that ends in `panic("implement me")`.
3. **The web UI cannot reach the CLI agent it's built to talk to.** Every web-client call to `datatug serve` is missing the server's hard-coded `/datatug` URL prefix (`datatug-apps/.../services/repo/agent.service.ts` etc. vs `datatug-cli/pkg/server/endpoints/routes.go:12-107`) — a mechanical, total mismatch across 8+ call sites, not a typo in one place.
4. **The project/query/entity data model is broad, coherent, and well-tested** (`datatug-core/pkg/datatug`, 99%+ coverage on core types) — Project, Environment/Server/DbCatalog, Entity/EntityField, Query/Parameters, Recordset, Board/Widget, FK/ReferencedBy — but it has **no runtime "related records" lookup**, **no semantic-classification engine** (naming patterns are data, never executed against a live schema), and **no query-result cache**.
5. **Server-enforced ACL is real and adapter-independent but wired into exactly one call site in the entire ecosystem** (`datatug-cli/apps/datatugapp/commands/cmd_query.go`) — `datatug serve`, the TUI, and every `datatug-apps` write path enforce zero access policy.
6. **A layered-ACL MVP is materially further along than a fresh read of shipped code would suggest.** A design packet (3 review rounds, founder-approved decisions H1-H7) and 19 commits of implementation across `dal-go/dalgo` (`layered-acl-query`, 65c4824, 9 commits: policy loader, scoped field/collection masks, typed principals, execution gates), `dalgo2sql` (45fe36c), `dalgo2ingitdb` (de6f64c), `openvaultdb-go` (7998895), and `datatug/dtql` (`layered-acl-design`, c2d0487) sit on the Hetzner VM, fetched locally as `vm/*` refs, **not pushed to GitHub** (blocked by the auto-mode classifier). The SpecScore plan (`spec/plans/layered-acl-mvp.md`) records tasks 1-10 complete, 11-24 queued, status **Executing** — Task 21, "DataTug query/write and Explain Access integration," is queued, not started.
7. **Marketing runs materially ahead of code on `datatug.io`**: an MIT-license claim contradicts the CLI's actual Apache-2.0 license (`datatug-io/src/pages/pricing.astro:11` vs `datatug-cli/LICENSE`); a documented `datatug mcp` server, `datatug serve --rpc` flag, and four top-level commands (`snapshot`/`diff`/`schema`/`transform`) do not exist in the CLI's command tree at all.
8. **`datatug-apps` main moved to `5d54b40` on 2026-09-08** ("enable enterprise OIDC/SAML SSO"), after the web-ui.md source report was written; the lead's 2026-09-09 browser walk of the current build confirms a `/login` page with Sign up / Sign in / "Sign in with company SSO" / Quick login / Google / Apple / Microsoft / Facebook / GitHub — real UI, depth not yet independently verified.
9. **Dashboardius (`sneat-co/dashboardius`) is a separate, live Sneat product** that re-exports DataTug's own `Board`/`BoardRow`/`BoardCard`/`BoardWidget` schema field-for-field into a PrimeNG app with Chart.js (`p-chart`) and switchable grid engines (plain table, Tabulator via `@sneat/datagrid`, AG Grid) — live at dashboardius.com on sample data only, while DataTug's own web UI board rendering remains unimplemented (`DatatugBoardService.getBoard()` → not implemented).
10. **The blocking defects cluster in a small number of places**: the storage-indirection layer (`storage.NewDatatugStore`/`storage.GetStore`) behind `init`/`demo`/`serve`; the web client's missing `/datatug` path prefix; six empty-shell web pages (Widgets/Tags/Resources/Variables/DB-Models/Diff); and the fact that ACL enforcement stops at the CLI's edge. Fixing the storage indirection and the path prefix would very likely unblock the majority of "broken" rows in the capability table below without any new feature work.

---

## Capability classification

| Capability | Pillar | State | Evidence | Note |
|---|---|---|---|---|
| `dal.StructuredQuery` AST (From/Where/GroupBy/Having/OrderBy/Columns/Limit/Offset/cursors) | P1 | partial | `dal-go/dalgo/dal/query.go:30-61` | Joins type-only; only `dalgo2memory` executes them; no adapter renders joins |
| DTQL YAML serialization | P1 | usable | `dal-go/dalgo/dtql/README.md`, roundtrip tests | Lossless relational read-only subset; joins/GroupBy/functions explicitly out of scope |
| DTQL as a saved query file format | P1 | prototype | `find /Users/alex/projects/datatug -iname '*.dtql.yaml'` → 0 hits | Only fixtures exist inside `dal-go/dalgo/dtql/examples/` |
| `datatug query run` (DTQL + ACL execution) | P1, P10 | usable | `datatug-cli/apps/datatugapp/commands/cmd_query.go`, live-verified | Only ACL-secured path in the whole CLI |
| `datatug/dtql`, `dtql-go`, `dtql-backstage` standalone repos | P1 | obsolete | 0 `.go` files, single 2026-06-09 commit | `dtql-org` landing page markets these as canonical; real impl is `dal-go/dalgo/dtql` |
| Layered ACL MVP policy format v2 (`dtql.org/access/v1`) | P1, P10 | specified/documented, implementation underway | `vm/layered-acl-design` c2d0487; `dal-go/dalgo` `layered-acl-query` 65c4824 (9 commits) | Not on GitHub; SpecScore plan tasks 1-10 complete, 11-24 queued, status Executing (lead, 2026-09-09) |
| `sqlite://` query/copy backend | P2 | usable | `datatug-cli/pkg/dbcopy/url.go` | Fully wired; used in every live CLI verification this audit ran |
| `ingitdb://` query/copy backend | P2 | usable (local only) | `pkg/dbcopy/url.go:106-117` | Remote-looking URLs explicitly rejected |
| `postgres://` query/copy backend | P2 | planned | `pkg/dbcopy/url.go:34-37,164-165` | Parses for shape only; `Open()` always `ErrPostgresNotWired` |
| Firestore as a DTQL/StructuredQuery source | P2 | partial | `dal-go/dalgo2firestore/query.go:11-52,135-138` | No GroupBy/Having/Columns; WHERE `Or` groups rejected; recordset reads unsupported |
| Firestore as a schema-scan source (not persistence) | P2 | usable | `datatug-cli/pkg/schemers/firestoreschema` | Different use of Firestore than a DataTug store |
| MSSQL/generic-SQL schema drivers | P2 | usable, relocated | moved `datatug-core`→`datatug-cli` in `a524b61` (2025-12-10) | Deliberate relocation, not a gap |
| HTTP/REST query execution | P2, P7 | broken | `datatug-apps/.../http-executor.ts:20` throws "Not implemented yet" | Model (`QueryDefTarget`, `ApiService`) is real; only the lower `executeCommand()` path works |
| `dalgo2ingitdb` adapter | P2 | usable | `ingitdb/dalgo2ingitdb/query.go:70-140` | Most feature-complete non-SQL query engine surveyed: real GroupBy+Having+aggregation+columns |
| `dalgo2openvaultdb` adapter | P2 | usable, small honest MVP | `dalgo2openvaultdb/query.go:39-53,99-100` | Explicit `ErrNotSupported` for Offset/cursor/GroupBy/Having/Columns/OR |
| `dalgo2ingitdb4local` adapter | P2 | prototype | `db_local.go:47-83` | DB-level methods `panic("implement me")`; conformance suite deliberately not run |
| Entity / EntityField model + `KnownTypes` | P3 | usable | `datatug-core/pkg/datatug/entities.go`, 99%+ coverage | |
| `datatug entity {add,list,show,field}` CLI | P3 | usable, spec Stable | live-verified against cloned demo project | |
| `EntityFieldRef` semantic tagging (params/columns) | P3 | usable | `parameters.go:20`, `recordset_def.go:65` | The actual "semantic tag" mechanism in use |
| Web UI entities pages (list/detail/edit) | P3 | usable | `datatug-apps` `EntitiesPageComponent` + `EntityService`, real CRUD | |
| `NamePatterns` (regex/exact naming classification) | P4 | prototype (data-only) | `entities.go:67-104` | No code in `datatug-core` executes pattern matching against a scanned schema |
| "AI-built knowledge base, auto-classify" marketing claim | P4 | unclear/likely aspirational | `datatug.io/knowledge-base.astro:18` | No corresponding classifier code found in any audited repo |
| FK / ReferencedBy schema graph | P5 | usable | `db_objects.go:251-341`, `pkg/schemer/scan_constraints.go:57-100` | Bidirectional, populated from a single constraint scan |
| TUI FK/referrers panels | P5 | implemented-but-poorly-exposed | `table_fks.go`, `table_referrers.go` | Real, live schema data; read-only, no click-to-navigate found |
| Web UI FK-driven related-record cards | P5 | usable | `env-db-table.page.html` "Record" tab, `foreign-key-card`/`ref-by-card` components, `goTable()` | Real click-through navigation — but depends on the broken agent path (P11 row below) for live data |
| Runtime "fetch related rows for this value" query execution | P5 | not implemented | grep for "related" across `datatug-core/pkg/**` (excl. tests) → 0 hits beyond FK/schema plumbing | Only `ParameterLookup{DB,SQL,KeyFields}` (a value-picker shape) exists |
| `datatug2md` JOIN-snippet doc generator | P5 | usable (doc-only) | `pkg/datatug2md/table_to_readme.go` | Not executable; the only place "JOIN" SQL text appears in the codebase |
| Variables / Tags data model | P6 | usable (model only) | `variables.go`, `tag.go` | |
| Web UI Variables/Tags/Resources pages | P6 | ui-only/mock | `VariablesPageComponent`/`TagsPageComponent`/`ResourcesPageComponent` — header + empty `<ion-content>` | `VariableService`/`RecordsetService` have zero methods |
| "Investigation context / basket" concept | P6 | unclear / not found | no matching commits (`git log -S"basket"`) in `datatug-cli`; no dedicated page beyond the empty Variables shell | |
| Query / QueryDef model + folders | P7 | usable | `datatug-core/pkg/datatug/query.go` | |
| Saved query text format on disk (`.sql` + `.sql.json` sidecar) | P7 | usable | `datatug-demo-projects/demo-project-1/queries/**` | Not DTQL; DTQL lives in a separate repo family |
| Web UI query editor + folders/tabs (Active/Bookmarked/Personal/Shared) | P7 | usable UI, partial execution | `QueryPageComponent`, `QueriesPageComponent`; SQL path via agent (broken pending P11 fix), HTTP path throws not-implemented | Tabs shell confirmed rendering live (lead, 2026-09-09) at `/store/localhost:8989/project/demo-project-1/queries` |
| `datatug queries` CLI command | P7 | obsolete/stub | `panic("not implemented")`, `cmd_queries.go:16-19`, reproduced live | Superseded in intent by `query run` |
| AI query-builder / ai-terminal-query-builder | P7, P14 | specified/documented, zero code | 3 specs `Status: Approved`; zero matching code in `datatug-apps` or `datatug-cli` | |
| MCP "serve-brokered query builder" | P7, P11, P14 | specified/documented, zero code | `spec/features/serve-brokered-query-builder/README.md` (Approved); `grep -rli mcp --include=*.go .` → 0 real hits | Depends on a `datatug serve` that cannot currently start |
| `datatug_projects` Firestore collection (own `userIDs[]`) | P8 | usable, not space-bounded | `sneat-firebase/firestore.rules:71-114` | Parallel to, not integrated with, Sneat's `/spaces/{spaceID}/ext/...` model |
| `@sneat/space-models` / `@sneat/space-services` deps | P8 | unused | 0 import sites, `grep -rn "@sneat/space"` across `libs/`/`apps/` | Declared dependency, dead |
| Space Explorer (`/explore/:spaceId`) | P8 | usable | `space-explorer-page.component.ts`, read-only, "READ-ONLY BY CONSTRUCTION" | Real transparency viewer PoC |
| KB sharing scopes (private/team/department/company) | P8 | unclear/likely aspirational | `datatug.io/knowledge-base.astro:7-11,55` | No implementation found in `datatug-cli` or `datatug-apps` |
| `spaceID` field on `ICreateProjectData` | P8 | unwired | `project.service.ts:193-197` | Typed, never populated or read |
| Filestore persistence (local dir project layout) | P9 | usable | `datatug-core/pkg/storage/filestore` | Only implemented `Store` backend |
| GitHub-repo-as-project-store | P9 | partial (write-only) | `pkg/dtgithub/github_storage.go:50-114` real; `FileExists`/`OpenFile` `panic("implement me")` (lines 40-48) | Write path real, cannot read back through the same interface |
| `--git=stage` mutation staging | P9 | usable, scoped | `git_integration.go`; wired only into `entity add` | `--git=commit` explicitly rejected, exit 2 |
| Project registry (`~/datatug/.datatug.yaml`) | P9 | usable | `dtconfig/settings.go` | |
| DALgo `access` package (default-deny, hierarchical, row/field/principal) | P10 | usable, spec Stable | `dal-go/dalgo/access/*`, `go test ./access/...` passes | |
| `access.SecureReadSession`/Write/Readwrite wrappers | P10 | usable, adapter-independent | `access/session.go:241-257` | Wraps any `dal.DB`; genuinely reused by `datatug-cli` |
| `datatug-cli` `pkg/accesspolicies` + `query run` enforcement | P10 | usable, narrowly scoped | `cmd_query.go` is the only call site (`grep -rln "accesspolicies\."`) | |
| `datatug serve` / `pkg/api` ACL enforcement | P10 | none | no auth/token/ACL check found in `pkg/server`/`pkg/api` | Legacy `pkg/sqlexecute` exec endpoints enforce nothing |
| OpenVaultDB own capability auth (`records:read`) | P10 | usable, coarse, uncoordinated | `openvaultdb-go/pkg/auth/auth.go:1-45` | Pinned to stale `dal-go/dalgo v0.64.4`, predates the Sept `access` work |
| InGitDB own ACL (shipped code) | P10 | none in shipped code | no policy code in `ingitdb-go`/`ingitdb`/`ingitdb-specs` | Design + first implementation commits exist on the VM (see row below) |
| Layered ACL MVP implementation (branch `layered-acl-query`, on origin since 2026-09-09, unmerged) | P10 | implementation underway, unmerged | `dal-go/dalgo` `layered-acl-query` 65c4824 (9 commits: policy loader, scoped field/collection masks, typed principals, execution gates); `dalgo2sql` 45fe36c (2 commits); `dalgo2ingitdb` de6f64c (3 commits); `openvaultdb-go` 7998895 (5 commits) | Not pushed to GitHub (auto-mode classifier blocked remote push, per lead) — foundation more advanced than a shipped-code-only read suggests |
| TUI screens (databases/environments/entities/settings/viewers) | P11 | usable | `pkg/sneatview/sneatnav` framework, live schema browsing | |
| TUI Queries/Dashboards screens | P7, P11, P15 | ui-only/mock | `dtproject/queries_screen.go:9-22`, `dashboards_screen.go:9-21` — hard-coded strings | Real backend models/routes exist for both |
| Ctrl+W TUI→web hand-off | P11 | partial (single route mapped) | `webui_handoff.go:26-28` — only `"projects": "/my"` | Every other screen falls back to the bare web origin |
| `datatug serve` HTTP/agent API | P11 | prototype (broken) | `cmd_serve.go:43` nil-deref; `register.go:31` panics on nil handler; `storage.NewDatatugStore` closure ends `panic("implement me")` (`http_server.go:33-39`) | Reproduced live by the lead on `v0.13.0`; 0.0% test coverage on `pkg/server` |
| `datatug-apps` local-agent HTTP client | P11 | broken (path mismatch) | client calls e.g. `{agent}/agent-info`; server only serves `{agent}/datatug/agent-info` (`routes.go:12-107`) | Total, mechanical mismatch across 8+ call sites, not isolated typos |
| Enterprise OIDC/SAML SSO login | P11 | newly landed, depth unverified | `datatug-apps` main `5d54b40` (2026-09-08); `/login` shows SSO/Quick login/Google/Apple/Microsoft/Facebook/GitHub (lead browser walk, 2026-09-09) | web-ui.md source report predates this commit |
| Live demo-project navigation (home → project) | P9, P11 | broken | lead's browser walk: clicking "DataTug Demo Project @ GitHub" → `No provider found for _DatatugFoldersService` + 404 on `github.com/datatug/datatug-demo-project` (real repo is `datatug-demo-projects`, file at `demo-project-1/datatug-project.json`) | Confirms `ecosystem-claims.md`'s independently-found broken demo link, with the root cause now identified |
| Explain-Access reason-code/`blockers[]` taxonomy | P12 | specified/documented | `vm/layered-acl-design` `04-authorization-contract.md` | Not present in any shipped code |
| `datatug-cli`'s current access-denial explanation | P12 | implemented-but-poorly-exposed | `pkg/accesspolicies/explain.go` derives attribution by splitting the English `Decision.Explanation` string on `" via "` | Flagged as a gap by the design packet itself |
| Environments / servers / DbCatalog model | P13 | usable | `environment.go`, `server.go`, `dbcatalog.go` | |
| DB schema scan (`pkg/schemer`, driver-agnostic) | P13 | usable | `scanner.go`; concrete drivers live in `datatug-cli` | |
| DB schema comparison (`DbModel` vs. scanned catalog) | P13 | usable | `datatug-core/pkg/comparator/compare_databases.go` | |
| `datatug scan` CLI (FK/index/constraint extraction) | P13 | usable | `cmd_scan_db.go`; commits `92a4140`/`21ddde4`/`032872a` | |
| AI skills plugin (`ai-plugin`, 9 `SKILL.md`s) | P14 | usable, mostly verified | `ai-plugin/skills/*` | 8 of 9 wrap real, present CLI subcommands |
| `query-builder` skill | P14 | unverified/not implementable | `ai-plugin/skills/query-builder/SKILL.md` | Depends on the unbuilt MCP daemon |
| Board/BoardRow/BoardCard/BoardWidget model | P15 | usable (model only) | `datatug-core/pkg/datatug/boards.go` | Widget-kind conflated with query-type — acknowledged code smell (`boards.go:153`) |
| Board rendering in `datatug-core` | P15 | backend-only, by design | no HTML/template/JS anywhere in `datatug-core` | Deliberate — rendering is a consumer's job |
| Web UI boards system (board-page, sql/grid/tabs widgets) | P15 | partial | real components exist; `DatatugBoardService.getBoard()` throws not-implemented on some paths; boards hard-code `localhost:8989` | |
| TUI dashboards screen | P15 | ui-only/mock | `dashboards_screen.go:9-21` hard-coded string | |
| Dashboardius (`sneat-co/dashboardius`) | P15 | usable demo / prototype product | live at dashboardius.com, PrimeNG + Chart.js `p-chart`, switchable grid engines (table/Tabulator via `@sneat/datagrid`/AG Grid), copies DataTug's board schema field-for-field (per `docs/IMPLEMENTATION-RECORD.md`) | Sample data only; brand-tier reconciliation with DataTug's Tier-3 status is an open question |

---

## Per-repository findings

### datatug-cli
The one place all the pieces are wired: real command tree with `entity`, `db copy`, `scan`, `query run` all live-verified against real data; `pkg/accesspolicies` (90.9% coverage) and `query run` are the best-tested, best-built surface in the ecosystem. But `init`/`demo`/`serve`/`show`/`dataset*`/`queries` (6+ commands) are broken or stubbed, all with zero or near-zero test coverage — a sharp, test-coverage-explained split between "new, spec-driven work" (229 commits since 2026-05-01) and untouched pre-cobra-migration leftovers. HEAD (`a913de4`) does not build (`go-github` v90/v91 import mismatch); the lead built and live-tested tag `v0.13.0` instead and reproduced the `serve` panics described above plus a `--help` crash on the legacy `updateUrlConfig` command (`-h` shorthand collision with cobra's own help flag). Two incompatible project-storage abstractions coexist — `filestore.New*` (works, used by most commands) and `storage.NewDatatugStore`/`GetStore` (never wired outside tests, used by `init`/`demo`/`serve`) — and the bug is self-acknowledged in a test that wraps the known panic in `recover()` (`pkg/server/endpoints/handlers_test.go:387-389`) rather than fixing it.

### datatug-core (+ demo repos, validators, small libs)
A pure engine+model library since a deliberate 2025-12-10 extraction (`a524b61` moved `pkg/api`, server endpoints, and the MSSQL driver out to `datatug-cli`) — zero HTTP surface, zero concrete DB drivers today, by design. The project model is broad, coherent, and ~99% tested. No related-records runtime lookup, no pattern-execution/classification code, no join/pipeline engine anywhere — these are all data shapes without a consumer in this repo. `datatug-cli` vendors a byte-identical copy of `pkg/datatug` instead of depending on the module (a commented-out `replace` directive in its `go.mod`) — a real drift risk. Demo repos are wildly uneven: `datatug-demo-projects/demo-project-1` (Chinook-based) is real and current; `datatug-demo-northwind`/`-adventureworks` are stubs/empty shells frozen since 2021-2023; the org's `chinook-database` fork is ~8.5 years behind its own upstream. `datatug-go-models` is a near-empty, zero-consumer repo. Secrets modeling is a real gap: `QueryDefTarget.Credentials` can land in git-tracked query files, and `Project.Validate()` never validates `Queries` at all.

### datatug-apps
Builds and unit-tests cleanly (`pnpm nx build datatug-app`; 115/115 tests pass, 28.96% statement coverage — mostly 1-assertion smoke tests). Three data backends per project (local CLI agent, Firestore, public GitHub raw fetches); the agent path is comprehensively broken by a missing `/datatug` URL prefix (see capability table), and Firestore-backed writes call `sneat-go` endpoints (`datatug/projects/create_project` etc.) that don't exist anywhere in `sneat-go` (confirmed by repo-wide grep — only an unrelated Telegram bot profile matches "datatug"). Six routed pages are empty shells (Widgets/Tags/Resources/Variables/DB-Models/Diff). "New Project" creation is doubly broken: the backend endpoint doesn't exist, and the client hard-codes `userIDs: []`, which would lock the creator out of their own project under Firestore rules even if the write succeeded. Git history is a 100% lossless migration from `sneat-co/sneat-apps` — nothing was gutted; real, substantial features (Space Explorer, Vault Explorer, the SQL query editor, entity CRUD, the Tabulator result grid with FK-driven related-record cards) coexist with dead-end navigation (`DbserverPageComponent` orphaned; a literal `"serverType"` path segment that should be `:serverType`) and hard-coded placeholder UI (a fake "SomeTable" references tab, a mock 3-person grid on the home page). Main advanced to `5d54b40` on 2026-09-08 with enterprise OIDC/SAML SSO, after this report's own source pass — see the web-ui companion audit for the full route matrix and the lead's hands-on walk of the current build.

### dal-go/dalgo + adapters + DTQL + ACL (incl. VM branches)
`dal-go/dalgo` main already contains the `dtql` package (Approved, Grade A: lossless YAML serialization of the relational read-only subset of `dal.StructuredQuery`) and the `access` package (Stable: default-deny, hierarchical, path-scoped, row-level `where`/`check` conditions, field allow-lists, principal bindings) — both landed 2026-06-05 through 2026-09-03. `access.SecureReadSession` wraps any `dal.DB` adapter-independently and is real code, not aspiration. Of 12 surveyed adapters, only `dalgo2ingitdb` executes GroupBy/Having/columns fully; the SQL family shares one text renderer; Firestore/Datastore/OpenVaultDB reject GroupBy/joins/Columns outright, and Firestore additionally rejects OR-groups in WHERE — a live gap given `access`'s residual-condition rewriter unconditionally pushes possibly-OR conditions into the query with no post-filter fallback (`access/condition.go:182-210`). DataTug's actual consumption is narrow: `datatug-cli` wires both packages into `query run`; `datatug-core` and `datatug-apps` import neither. **The layered-ACL MVP is materially further along than the DALgo report alone would suggest**: a design packet (134 files, 26,598 lines, 3 review rounds, founder-approved H1-H7) and real implementation — `dal-go/dalgo` branch `layered-acl-query` (65c4824, 9 commits: portable ACL policy loader, scoped field/collection masks, typed principals, execution gates), `dalgo2sql` (45fe36c, 2 commits), `dalgo2ingitdb` (de6f64c, 3 commits), `openvaultdb-go` (7998895, 5 commits), and `datatug/dtql` branch `layered-acl-design` (c2d0487) — all exist on the Hetzner VM, fetched locally as `vm/*` remote refs, and are **not pushed to GitHub** (the lead's auto-mode classifier blocked the remote push). The SpecScore plan `spec/plans/layered-acl-mvp.md` records tasks 1-10 complete, 11-24 queued, status **Executing**, with the note "User approval covers implementation... avoids frequent remote pushes/merges." Task 21, "DataTug query/write and Explain Access integration," is queued, not started — so Phase-1 DataTug integration work has a design and a partial DALgo-side foundation to build on, but no DataTug-facing code yet, and the whole packet is one disk away from being lost until it lands on `origin`.

### Ecosystem claims (marketing/datatug.io/ai-plugin/Sneat integration/Dashboardius)
`datatug.io` and the CLI's own internal risk log are unusually candid with each other, but two public claims are flatly false: MIT licensing (CLI is Apache-2.0) and a working `datatug mcp`/`--rpc` server (zero matching code). Four top-level commands advertised on the docs/homepage (`snapshot`/`diff`/`schema`/`transform`) don't exist in the command tree. The AI-agent story splits into two real maturity tiers: the skills plugin (`ai-plugin`, 9 skills, 8 wrapping real commands) is genuine; the MCP-brokered query builder the site markets is 100% unbuilt. `datatug.app` is live in production (Ionic/Angular shell, shared `sneat-eur3-1` Firebase auth) but its one canonical "click here for the demo" path is broken (confirmed independently by both the ecosystem audit and the lead's live browser walk — see capability table). DTQL is marketed on a polished `dtql.org` site as a working query language; on disk the `dtql`/`dtql-go`/`dtql-backstage` repos are one-line-README stubs — the real implementation lives in `dal-go/dalgo/dtql` under a different module path than the one the site tells visitors to `go get`. DataTug shares Sneat's production Firebase project but keeps its own parallel `datatug_projects` membership model rather than Sneat's space-bounded convention; `@sneat/space-models`/`@sneat/space-services` are declared dependencies with zero import sites. Dashboardius (`sneat-co/dashboardius`, public repo) is a brand-new (2026-09-08), separate Sneat product — PrimeNG standalone Angular, Chart.js via `p-chart`, switchable grid engines (plain table / Tabulator via `@sneat/datagrid` / AG Grid), deployed at dashboardius.com, reusing DataTug's Board schema field-for-field per its own `docs/IMPLEMENTATION-RECORD.md` — while structurally depending on DataTug's Tier-3-independent data model, a tier tension its own spec (`spec/ideas/dashboardius.md`) flags as unresolved. DataTug itself stays Tier-3 independent by explicit, repeatedly-restated brand-architecture rule.

---

## Journey to first useful result (measured conceptually)

| Stage | Status today | Blocking defect(s) |
|---|---|---|
| **Time-to-first-connected-source** | Blocked via the documented onboarding path; works via manual CLI | `datatug init` panics unconditionally; `datatug demo` downloads real data (Chinook DB + `datatug-demo-project` clone) but fails to save the project (`storage.GetStore` "no store configured"). Working path exists only if the user manually `git clone`s a project directory and runs `entity list -d <dir>` — undocumented, CLI-only. |
| **First query/result** | Works, CLI-only | `datatug query run --db sqlite://<path> --from <table> --no-policies` is live-verified end-to-end. No web/agent path reaches it: `datatug serve` panics before listening, and even fixed, the web client's path-prefix mismatch would still 404 every call. |
| **First related-data discovery** | Broken everywhere it would matter | FK/referrer data is real and complete (schema-level graph), but the TUI panel is read-only with no click-to-navigate, and the web UI's real click-through FK cards (`env-db-table.page.html`) depend on the same broken agent connection as the query path above. No runtime "fetch rows related to this value" execution exists in `datatug-core` at all — only a value-picker shape (`ParameterLookup`). |
| **First shared-query reuse** | UI shell exists, nothing behind it | The web UI's Queries page renders ACTIVE/BOOKMARKED/PERSONAL/SHARED tabs (confirmed live by the lead, 2026-09-09) but the underlying agent connection is broken so no real query list loads; the CLI has no "shared" query concept beyond folders; the MCP query-builder meant to give an agent shared reuse is unbuilt. |
| **Understanding current context** | Partial, real navigation state exists but surfaces errors | TUI project-home menu (databases/environments/queries/dashboards) is real and works. Web UI has a real nav-context service layer (`DatatugNavContextService`, `ProjectTracker`/`StoreTracker`), but the lead's live walk hit `No provider found for _DatatugNavService` and an `NgModule` assertion error navigating into a project's environment/DB page — the context machinery exists but is not fully wired for the routes exercised. |

---

## Contradictions between docs, specs, marketing and code

| Claim | Where | Reality | Evidence |
|---|---|---|---|
| "Open source. **MIT**-licensed." | `datatug-io/src/pages/pricing.astro:11`; Homebrew cask `license: MIT` | CLI is Apache License 2.0 | `datatug-cli/LICENSE`; parked, not resolved, per owner (`marketing/distribution/datatug-app-landing.md:68`) |
| `datatug mcp` server; `datatug serve --rpc --addr :7842` | `datatug-io/docs/integrations/mcp.md`, `integrations.astro:22` | Zero `MCP` hits in CLI Go source; `serveCommandArgs()` registers no flags at all | `grep -rn "MCP" --include=*.go` datatug-cli → 0; `cmd_serve.go` |
| Top-level commands `snapshot`, `diff`, `schema`, `transform` | `datatug-io/docs/reference/cli.md`, homepage capabilities grid | None registered in `DatatugCommand()`; no corresponding `cmd_*.go` files | `apps/datatugapp/commands/command.go` |
| `datatug query --project demo --sql "SELECT 1"` | `datatug-io/docs/quickstart.md` | Actual command is `query run` with `--db`/`--file`/`--from`/`--as`/`--role`; no `--project`/`--sql` flags | `cmd_query.go` |
| MONETIZATION.md: "The MCP server" is "free, today" | `datatug-cli/marketing/MONETIZATION.md:18` | Contradicts the same repo's own spec status (Approved, unbuilt) — an internal-doc-vs-code contradiction, not just a site bug | Same repo |
| DTQL: "read it, hand-edit it, diff it, reload it losslessly" (live dtql.org site) | `datatug.io` homepage; dtql.org | `dtql`/`dtql-go`/`dtql-backstage` are one-line-README stubs; the real relational-query work lives in `dal-go/dalgo/dtql`, a different repo/org | repo listings; `dal-go/dalgo/dtql/README.md` |
| Sharing scopes: private/team/department/company for KB entries | `datatug-io/knowledge-base.astro` | No code path found implementing a 4-tier sharing model in either `datatug-cli` or `datatug-apps` | grep across both repos |
| "Dashboards, curated queries in the web UI" | `datatug.io` homepage, `platform.astro` | `DatatugBoardService.getBoard()` returns not-implemented; no chart library (`echarts`/`ngx-charts`/`chart.js`/`d3`) found anywhere in `datatug-apps` | `spec/ideas/dashboardius.md:25`; grep |
| Demo project opens at `datatug.app/pwa/agent/github/project/datatug-demo-project@datatug` | `datatug-demo-projects/README.md` | Route doesn't exist in current `datatug-apps` router; live-browsed 2026-09-09 renders a Sentry crash screen; separately, the lead's own walk hit `No provider found for _DatatugFoldersService` clicking the demo item on the current build, plus a 404 against the wrong repo name (`datatug-demo-project` vs. real `datatug-demo-projects`) | `datatug-routing.module.ts`; both audits' live browsing |
| `cli/init` spec status "Implementing" claims init "produces a lint-clean baseline... immediately openable by show, serve, and the TUI" | `datatug-cli/spec/features/cli/init/README.md` | `init` panics unconditionally; `show` cannot succeed under any invocation; `serve` panics before listening | `cmd_init_project.go:112`, `project_base_command.go:13-23`, `cmd_serve.go:43` |
| Layered-ACL MVP is "design only, no implementation authorized" | `vm/layered-acl-design` packet's own self-reported status | SpecScore plan shows 10 of 24 tasks already complete on implementation branches, status Executing (lead, 2026-09-09) — the design-packet's status line is about *that specific packet's own sign-off gate*, not the parallel implementation track that had already started | `spec/plans/layered-acl-mvp.md` |

---

## Reuse inventory for implementation

**Serve/agent + ACL**
- `datatug-cli/pkg/accesspolicies/{run,dir,explain}.go` + `cmd_query.go` — the one place policy loading, principal construction, and `access.SecureReadSession` are correctly wired end to end; a server-side "run this query" endpoint should call the exact same `accesspolicies.Run(ctx, db, query, Options{...})`.
- `pkg/dbcopy/url.go`'s `BackendRef.Parse`/`Open` — reuse for opening a `dal.DB` from any URL a server request names, rather than a new scheme dispatcher.
- Fix `datatug-apps/.../services/repo/agent.service.ts`, `agent-state.service.ts`, `project.service.ts`, `project-item-service.ts` to add the server's `/datatug` prefix (`pkg/server/endpoints/routes.go:12-107`) — the single highest-leverage fix, likely unblocking most agent-backed web features at once.
- `dal-go/dalgo/access/{document,policy,session}.go` — do not write a new evaluator; extend this one. `vm/layered-acl-design` (branch `layered-acl-design` on `datatug/dtql`, fetch via the VM) is the design to read before designing anything new — already through two independent reviews and founder sign-off on H1-H7.

**Parameterized queries + semantic tags**
- `cmd_query.go`'s `readQueryOptions`/`buildQuery`, `accesspolicies.ParseVariables`/`substituteParams` — the reference implementation for `--var name=value` parameter binding.
- `EntityFieldRef` (`entities.go:145-149`) — the existing "tag this parameter/column with semantic meaning" mechanism; reuse rather than invent a new one.
- Web-side: `models/definition/parameter.ts`, `parameter-lookup.service.ts`, `input-parameters.component.ts` — real, working, already used by `board-page.component.ts` and `sql-query-editor.component.ts`.

**Related-data lookup**
- `schemer.ForeignKeysProvider`/`ReferrersProvider` (`pkg/schemer/interfaces.go`) — enumerate the relationship graph; already implemented for SQLite/SQL-info-schema/MSSQL.
- Web: `pages/signed-in/env-db-table/foreign-key-card/`, `ref-by-card/`, `record-values-card/` — a real, working FK-driven navigation pattern to replicate elsewhere.
- TUI: `foreignKeysBox`/`referrersBox` (`table_fks.go`, `table_referrers.go`) need a `SetSelectedFunc` — they currently only implement display.

**Result grid + cell selection**
- `@sneat/datagrid` (wraps Tabulator.js, `SelectRowModule`/`InteractionModule`/`MenuModule` already registered) — row select/deselect, context menu, group-by are real; cell-range/multi-cell selection is not exposed (Tabulator's `range` module isn't registered in the published package).
- The FK-popover cell-formatter pattern in `env-db-table.page.ts` (`cellFormatter`, `CellPopoverComponent`) is the template to extend for other cell-level interactions.
- Dashboardius already demonstrates a switchable-grid-engine pattern (plain table / Tabulator via `@sneat/datagrid` / AG Grid) worth reading before choosing a second grid engine for DataTug itself.

**Context / basket**
- No working code to reuse — `VariablesPageComponent`/`VariableService` and `pkg/datatug/variables.go` are the closest data shapes (both currently inert); a Phase-1 "context" feature would be new work on top of the `Variables{Vars map[string]VarInfo}` model.

**Demo dataset**
- `datatug-demo-projects/demo-project-1` (Chinook-based, real, current) is the only living, usable demo — 7 entities, 3 queries, 5 environments; `Customer`/`Employee` tables exist but have no `Entity` wrapper yet (gap for a customers/orders-style Phase-1 pitch).
- `pkg/dbcopy/testdata/`, `pkg/dbcopy/filter/testdata/` — real schema/data fixtures for copy+filter testing.

**Tests / fixtures**
- `pkg/accesspolicies/*_test.go`, `cmd_query_test.go` (90.9% coverage) — the CLI-side ACL reference tests.
- `dal-go/dalgo/access/*_test.go` (~130 KB), `dalgotest.RunConformance` (shared adapter conformance suite, real and passing for `dalgo2sqlite`, `dalgo2ingitdb`, `dalgo2openvaultdb`) — reuse the conformance harness for any new adapter rather than writing bespoke tests.

---

## Open questions for the founder

1. Is the `a913de4` build break (go-github v90→v91) fixed anywhere, or does every fresh `git clone && go build` of `datatug-cli` main currently fail? Blocks any new contributor or CI run from main.
2. Is `storage.NewDatatugStore`/`storage.GetStore` dead scaffolding that should be deleted in favor of `filestore.New*` (the path that actually works), or was it built for a specific future consumer (the MCP daemon?) and simply never finished? Determines whether fixing `serve` means "wire the existing indirection" or "remove it."
3. Are `show`/`dataset*`/`datasets`/`queries` intentionally deprioritized in favor of `entity`/`query run`, or simply unnoticed? Determines whether Phase-1 effort should fix or delete them.
4. ~~Should the layered-ACL MVP packet and its implementation branches (`vm/*` refs) be pushed to `origin` now?~~ Answered: pushed on 2026-09-09 (see the update note at the top). Open instead: whether implementation sign-off was given after the packet (ROADMAP open question 4).
5. Given Task 21 of the layered-ACL plan ("DataTug query/write and Explain Access integration") is queued and not started, should Phase-1 prioritize landing tasks 11-24 of that plan before building new DataTug-side ACL UI from scratch?
6. Is Firestore-as-a-project-store a live near-term goal, or were the `firestore:` struct tags on every `datatug-core` type added speculatively and since abandoned?
7. Is `datatug-cli`'s vendored, byte-identical copy of `datatug-core/pkg/datatug` intentional, or drift that should be fixed by restoring the module dependency? Affects whether "the project model" is one thing or two slowly-diverging things.
8. Is there an actual plan for a `datatug`-specific `sneat-go` module (writes for environments/schemas/projects), or should Firestore-backed projects stay read-only/agent-only for now, given `sneat-go` has zero routes for any `datatug/...` path today?
9. Should `datatug_projects` move to the space-bounded model (`/spaces/{spaceId}/ext/datatug/...`), given `ICreateProjectData.spaceID` already exists as an unused field and the `@sneat/space-*` packages are already (unused) dependencies?
10. Are Widgets/Tags/Resources/Variables/DB-Models/Diff genuinely next on a roadmap, or should they be hidden from web nav until built — they are currently one click away and totally empty?
11. Does Dashboardius belong to a DataTug project or a Sneat space, and what is its brand tier (Tier-1/2 `-ius` naming vs. its dependency on DataTug's Tier-3 data model)? Its own spec flags this as unresolved and it blocks persistence design.
12. Should `datatug/dtql`, `dtql-go`, `dtql-backstage` be retired/redirected and `dtql-org`'s landing page corrected, given the real implementation lives in `dal-go/dalgo/dtql` under a different module path than what the site advertises?
13. Should the parked Apache-vs-MIT license inconsistency be resolved now, given it is now baked into the release pipeline (Homebrew cask), not just marketing copy?
14. What is the disposition of the fresh `dalgo2ingitdb` `feat(access)`/`fix(access)` commits noted mid-audit — a first slice of the layered-ACL plan's InGitDB work packages already underway, or unrelated?

---

## Open Questions

1. Is the `a913de4` build break (go-github v90→v91) fixed anywhere, or does every fresh `git clone && go build` of `datatug-cli` main currently fail?
2. Is `storage.NewDatatugStore`/`storage.GetStore` dead scaffolding to delete, or built for a specific future consumer and never finished?
3. Are `show`/`dataset*`/`datasets`/`queries` intentionally deprioritized, or simply unnoticed?
4. Should the layered-ACL MVP packet and its implementation branches be pushed to `origin` now, independent of the Phase-1 decision?
5. Given Task 21 of the layered-ACL plan is queued, should Phase-1 prioritize landing tasks 11-24 before building new DataTug-side ACL UI from scratch?
6. Is Firestore-as-a-project-store a live near-term goal, or abandoned speculative tagging?
7. Is `datatug-cli`'s vendored copy of `datatug-core/pkg/datatug` intentional or drift?
8. Is there an actual plan for a `datatug`-specific `sneat-go` module, or should Firestore-backed projects stay read-only/agent-only for now?
9. Should `datatug_projects` move to the space-bounded model, given `spaceID` and the `@sneat/space-*` packages already exist unused?
10. Are Widgets/Tags/Resources/Variables/DB-Models/Diff genuinely next on a roadmap, or should they be hidden from nav until built?
11. Does Dashboardius belong to a DataTug project or a Sneat space, and what is its brand tier?
12. Should `datatug/dtql`, `dtql-go`, `dtql-backstage` be retired/redirected and `dtql-org`'s landing page corrected?
13. Should the parked Apache-vs-MIT license inconsistency be resolved now that it's baked into the release pipeline?
14. What is the disposition of the fresh `dalgo2ingitdb` `feat(access)`/`fix(access)` commits noted mid-audit?
