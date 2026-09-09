# DataTug feature hierarchy — pillars → capabilities → features

**Status:** Draft for founder review · **Date:** 2026-09-09 · **Owner:** alex
**Reads with:** [product thesis](README.md) · [current-state audit](../research/2026-09-09-current-state-audit.md) · [ROADMAP.md](../../ROADMAP.md)

## How to read this page

- **Pillar** — one of the fifteen defining pillars from the thesis. Each is a top-level
  Capability Feature under [`spec/features/`](../features/README.md).
- **Capability** — a coherent unit of behaviour under a pillar; the sub-feature level.
  Big enough to demo, small enough to own. Not a task.
- **State** — the audited state on 2026-09-09 using the audit vocabulary
  (usable · implemented-but-poorly-exposed · partial · prototype · backend-only ·
  ui-only/mock · specified/documented · planned · obsolete · unclear).
- **Phase** — where the roadmap places the capability: 0 (repair), 1 (core loop),
  2 (library depth), 3 (playbooks, scopes, dashboards), 4 (layered ACL, agents),
  5 (inference, discovery, infra), X (experiment), — (not building).
- **Kind** — *strategic* (defines DataTug), *enabler* (needed so strategic work can
  land), *hygiene* (must be true but differentiates nothing).

Existing Features that predate this hierarchy keep their slugs and are linked from the
pillar they belong to: `ai-query-builder` and `serve-brokered-query-builder` (both
Capabilities with cross-repo Implementations) sit under P14.

## P1 — Portable query model (DTQL + DALgo) · Feature `portable-query-model`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| DTQL YAML serialization of the DALgo query AST | usable (external, Stable in `dal-go/dalgo/dtql`) | — | enabler | Lossless round-trip tests; JSON Schema published. Joins, group-by, cursors out of DTQL scope. |
| DTQL as a saved-query type in projects (`queries/**/*.dtql.yaml`, `QueryDef.Type = DTQL`) | planned | 1 | strategic | `QueryDef.Type` enum has SQL/HTTP/StructuredSQL only (`datatug-core/pkg/datatug/query.go:124-137`). |
| Native query escape hatch (SQL text, HTTP request) | usable (SQL) / backend-only (HTTP) | 1 (SQL; read-only HTTP GET per founder ruling 2026-09-09), 2 (rest of HTTP) | enabler | SQL runs through legacy `pkg/sqlexecute`; HTTP executor throws in web. |
| Rendering to backends through DALgo adapters | partial | 0–1 | enabler | sqlite + ingitdb wired in `pkg/dbcopy/url.go`; postgres explicitly stubbed; mssql/mysql/firestore not in the dispatcher. |
| Standalone `dtql*` repositories and dtql.org claims | obsolete / contradicted | 0 | hygiene | Real code is `dal-go/dalgo/dtql`; retire or redirect the stubs. |

## P2 — Heterogeneous operational sources · Feature `heterogeneous-sources`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| SQL sources (SQLite, SQL Server, MySQL, Postgres) — scan + query | partial | 0–1 | enabler | Scan via `pkg/schemers/*`; released binaries are `CGO_ENABLED=0` so SQLite via mattn is absent; use `dalgo2sqlite` (pure Go) for structured queries. |
| inGitDB (Git-backed data) as a source | usable (local paths) | 1 | strategic | `ingitdb://` scheme; proven in `query run` tests. The cheapest second source for a cross-source demo. |
| Firestore as a source | partial | 3 | enabler | Scan display only in TUI; `dalgo2firestore` rejects OR groups and GroupBy. |
| HTTP / GraphQL API requests as queries | backend-only | 1 (read-only JSON GET reference data, founder ruling 2026-09-09) / 2 (rest) | strategic | `QueryDef` HTTP type + `HTTPWidgetDef` modeled; no executor in web, untested in CLI. |
| Logs as a source | planned | 5 | strategic | No model. |
| Source registry and connection secrets outside project files | partial (gap) | 0–1 | hygiene | `QueryDefTarget.Credentials` can be committed; `AuthCredential` stores no secret. |

## P3 — Semantic entities and fields · Feature `semantic-layer`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Entities with typed fields and key fields | usable | 1 | strategic | `datatug entity` CLI (Stable spec); web entities CRUD. |
| Entity → tables mapping | usable (model) | 1 | strategic | `Entity.Tables []DBCollectionKey`. |
| Field → column mapping (explicit, per source) | planned | 1 | strategic | No per-field physical mapping type; needed for `Customer.ID → chinook.Customer.CustomerId`. |
| Semantic tags on parameters and result columns (`EntityFieldRef`) | implemented-but-poorly-exposed | 1 | strategic | `ParameterDef.Meta`, `RecordsetColumnDef.Meta` exist; nothing reads them at runtime. |
| Semantic field types (money, email, GUID… ; "currency" seed) | partial | 2 | enabler | `KnownTypes`; seed idea `semantic-field-types-currency…` in datatug-cli ideas. |
| Entity summary in the demo project covers Customer/Employee | planned | 1 | hygiene | 11 tables, 7 entities; Customer/Employee missing. |

## P4 — Meaning inference · Feature `meaning-inference`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Deterministic name-pattern resolver (`NamePatterns` → column classification) | implemented-but-poorly-exposed (data only) | 1 | strategic | `EntityField.NamePatterns` regexp/exact exist; no code applies them. |
| Suggestions from type, neighbours, conventions, sampling | planned | 5 | strategic | Probabilistic; must stay inspectable and correctable. |
| Review, accept, correct, and record provenance of a mapping | planned | 2 | strategic | Ties to P12. |

## P5 — Relationships and related-data navigation · Feature `related-data-navigation`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Foreign-key graph scanned bidirectionally | usable | 1 | enabler | `ForeignKey`/`ReferencedBy`, `pkg/schemer/scan_constraints.go`. |
| FK-driven related cards on the table page (web) | usable UI / broken data path | 0–1 | strategic | `env-db-table` Record tab with `fk-card`/`ref-by-card`; agent path mismatch. |
| Semantic relationships across sources (same semantic field in two sources) | planned | 1 | strategic | Derived from P3 mappings: no new graph type needed for the first cut. |
| Related-record lookup executed by the server for a selected value | planned | 1 | strategic | Build a DALgo `StructuredQuery` from FK/semantic relationship + value; run through `accesspolicies.Run`. |
| Cell/value selection → context panel (meaning, related, applicable, actions) | planned | 1 | strategic | Tabulator grid supports cell click; `CellPopoverComponent` is hard-coded sample text today. |
| Automatic relationship discovery (beyond declared FKs) | planned | 5 / X | strategic | Value-overlap sampling is an experiment before it is a feature. |

## P6 — Investigation Context · Feature `investigation-context`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Context basket: add a semantic value from a cell, remove, temporarily disable | planned | 1 | strategic | Web "Variables" page is an empty shell; `Variables` model exists in core. Name decision: keep the internal `basket`; user-facing label "Investigation Context". |
| Automatic binding of semantic parameters from context and selection | planned | 1 | strategic | Resolver: `ParameterDef.Meta == basket value's EntityFieldRef`. |
| Visible "why this query is constrained" and no hidden filters | planned | 1 | strategic | Applied bindings shown on the query page before run. |
| Auto-narrowing compatible queries by environment/time | planned | 2 | strategic | Requires Environment/Time semantic fields. |
| Conflicting values, propagation between related queries | planned | 2 | strategic | Conceptual model in the Phase-1 Feature; behaviour later. |
| Saved / shareable investigations and history | planned | 3 | strategic | Trail of queries, selections, pivots → convertible into a playbook. |

## P7 — Shared executable knowledge library · Feature `executable-knowledge-library`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Parameterized queries with typed parameters | usable (model, CLI) / partial (web) | 1 | strategic | `parameters.go`; web `input-parameters` + `ParameterLookupService` exist. |
| Semantic parameters (`requires Customer.ID`) | implemented-but-poorly-exposed | 1 | strategic | `ParameterDef.Meta`. |
| Query applicability ("useful right now") from selection/context | planned | 1 | strategic | Server endpoint: applicable queries for a set of semantic values. |
| Explainability of a suggestion (selected column → maps to → query requires) | planned | 1 | strategic | Return the resolution chain with each suggestion. |
| Query collections (folders with meaning) | usable (folders) | 2 | strategic | `QueriesFolder` recursion; decide whether collections need more than folders. |
| Capture from exploration (save the current query/lookup as a library asset) | planned | 2 | strategic | Turns ad-hoc success into a team asset. |
| Investigation playbooks (entry points, branches, HTTP steps) | specified here only | 3 | strategic | Not a workflow engine: human-guided. |
| Knowledge promotion lifecycle (personal → team → project → org) | planned | 3 | strategic | Needs P8. |
| Templates / public knowledge packs | planned | X | strategic | Separate from organization-private knowledge. |

## P8 — Sharing scopes on Sneat spaces · Feature `sharing-scopes`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Personal / team / project / organization scopes | planned | 3 | strategic | Marketed on datatug.io; nothing implemented. |
| Project membership through Sneat space membership and roles | planned | 3 | enabler | `datatug_projects` uses flat `userIDs[]`; `@sneat/space-*` declared, unused; `spaceID` typed, unwired. Founder ruling 2026-09-09: a project is a standalone entity granted to multiple spaces. |
| Hosted team workspace on datatug.app | partial | 3 | strategic | Nearest monetizable wedge per backstage B2B strategy (2026-07-01). |

## P9 — Git-native projects · Feature `git-native-projects`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Readable file layout for every asset | usable | — | strategic | `pkg/storage/filestore`; demo project on GitHub. |
| Mutating commands stage their own files (`--git=stage`) | usable (entity add only) | 2 | enabler | `--git=commit` deferred by design. |
| Secrets never in project files (validation + guidance) | planned | 0–1 | hygiene | Add a validator rule; document `~/.datatug` credentials. |
| GitHub-hosted projects: read and write | partial (write-only) | 2 | strategic | `GhStorage.FileExists/OpenFile` panic. |
| Project store abstraction with pluggable backends (local files · GitHub · **remote OpenVaultDB** · DataTug Cloud) | partial (interface exists; filestore only) | 3 | strategic | Founder, 2026-09-09: *"It should be possible to store and use datatug project from remote OpenVaultDB."* `storage.Store`/`ProjectStore` interfaces exist; `dalgo2openvaultdb` exists; OpenVaultDB enforces its own policies, composing with DataTug's (layered-ACL design). |
| **Project store explorable inside DataTug** (the project's own entities, queries, mappings, policies browsable and queryable as a source) | planned | 2–3 | strategic | Founder, 2026-09-09: *"The project store itself should be available for exploration in DataTug itself."* Natural when project records are inGitDB/OpenVaultDB collections: the same grid, related-data and context loop applies to the knowledge graph itself. |
| Diff/review affordances for knowledge changes | planned | 4 | strategic | Rides on Git; UI later. |

## P10 — Server-enforced access control · Feature `server-enforced-access`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| DALgo policies: default-deny, hierarchical, row conditions, field allow-lists, principals | usable (external, Stable) | — | enabler | `dal-go/dalgo/access`. |
| CLI `query run` enforcement | usable | — | strategic | `pkg/accesspolicies`, exit code 5, tests against inGitDB. |
| `datatug serve` enforcement: every read passes `accesspolicies.Run`; principal from the serving user or an authenticated caller | planned | 1 | strategic | `serve` enforces nothing today; legacy `pkg/sqlexecute` path bypasses policies. |
| Source-level allow/deny per principal | planned | 1 | strategic | Policy path scope per source/collection. |
| Opaque SQL text under ACL: source-level only, read-only session; row/column restrictions apply to structured (DTQL) queries | planned (explicit limitation) | 1 | hygiene | Consistent with the layered-ACL design ("native routes unsupported under the protected profile"). |
| Layered multi-owner policies, structured blockers, Explain Access, policy UI | specified/documented + implementing on branch `layered-acl-query` (on origin since 2026-09-09) (dalgo plan tasks 1–10 complete) | 4 | strategic | `datatug/dtql` `layered-acl-design`; dalgo `spec/plans/layered-acl-mvp.md` task 21 is the DataTug integration. |

## P11 — Three surfaces, one project · Feature `surfaces`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| TUI (tview) screens | usable (Queries/Dashboards panels are placeholders) | 2 | strategic | `apps/datatugapp/datatugui`. |
| Web UI through the local agent (`datatug serve`) | prototype (broken both ends) | 0 | enabler | Server panics on start; client omits the `/datatug` prefix; `unsupported format of store id`. |
| Hosted web UI (datatug.app) with Sneat identity and enterprise SSO | partial | 3 | strategic | Live; demo link broken; SSO landed 2026-09-08. |
| Screen-level hand-off CLI → web (Ctrl+W) and web → CLI | partial (one route mapped) | 2 | strategic | `datatugui.screenPathToWebPath`. |
| Agent skills plugin | usable (8 of 9 verified) | 2 | strategic | `datatug/ai-plugin`. |
| MCP face / serve-brokered query builder | specified/documented (Approved), zero code | 4 | strategic | Depends on a working `serve`. |
| Dead pages hidden until built (Widgets, Tags, Resources, Variables, DB models, Diff) | prototype shells | 0 | hygiene | Six empty shells one click away. |

## P12 — Provenance · Feature `provenance`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Result provenance: source, table/collection/endpoint, query, policy limitations | planned | 2 | strategic | `accesspolicies.Explain` is the seed; structured provenance replaces English parsing (layered-ACL design). |
| Mapping provenance: who/what declared or inferred a mapping | planned | 2 | strategic | Ties to P4 review. |

## P13 — Infrastructure inventory · Feature `infrastructure-inventory`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Environments, servers, catalogs | usable | — | enabler | `environment.go`, `server.go`; scanned by `scan`. |
| Services and endpoints linked to sources | prototype (vestigial `ApiService`/`ApiEndpoint`) | 5 | strategic | Not persisted by any store. |
| Climb from a record to its infrastructure | planned | 5 | strategic | Order → source → service → environment → server. |

## P14 — AI uses the same model · Features `ai-query-builder`, `serve-brokered-query-builder`, `agent-context`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Skills plugin driving the CLI | usable | — | strategic | Verified 8/9. |
| Agent-readable project context (entities, mappings, applicable queries) as a CLI/MCP resource | planned | 2 | strategic | The same resolver the web uses, exposed as `datatug context` / MCP resource. |
| Serve-brokered query builder (MCP + HTTP/WS faces) | specified/documented | 4 | strategic | Grade-A Capability; implementations Planned. |
| In-app AI query builder | specified/documented | X | strategic | Alternative delivery; keep as an experiment behind the serve-brokered path. |

## P15 — Dashboards as a surface · Feature `dashboards`

| Capability | State | Phase | Kind | Evidence / note |
|---|---|---|---|---|
| Board schema (Go + TS) | usable (model) | — | enabler | `datatug-core/pkg/datatug/boards.go`; copied field-for-field into Dashboardius. |
| Board persistence and query binding through the DataTug server | planned | 3 | strategic | `DatatugBoardService.getBoard()` not implemented; Dashboardius MVP needs it. Founder ruling 2026-09-09: boards belong to the DataTug project. |
| Board emerging from exploration ("pin this result") | planned | 3 | strategic | The DataTug-side entry point into Dashboardius. |
| Shared board-model package (`@datatug/board-models`) | planned | 3 | enabler | Removes the documented copy in `sneat-co/dashboardius`. |
| Charts in DataTug's own web UI | planned | — | — | Not building; Dashboardius presents. |

## Explicitly not building

Generic workflow engine · no-code app builder · catalog/governance suite · warehouse or
centralization · connector breadth race · client-direct DB access for protected sources ·
text-to-SQL chat box disconnected from the knowledge model · charts inside DataTug's own
web UI beyond what Dashboardius provides.

## Open Questions

- Should `queries/`-level collections gain semantics beyond folders (owner, domain,
  entry points), or is a folder plus a README the simplest useful model? (P7, Phase 2)
- Is the internal name `basket` kept for the Investigation Context, with only the
  user-facing label changed? (P6) — recommendation: yes.
- Which second source proves heterogeneity in Phase 1: an inGitDB collection (cheapest,
  already wired) or an HTTP endpoint (closer to the shipping-status example, needs the
  executor)? — recommendation: inGitDB in Phase 1, HTTP in Phase 2.
