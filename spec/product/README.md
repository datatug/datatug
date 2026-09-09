# DataTug — product thesis and pillars

**Status:** Draft for founder review · **Date:** 2026-09-09 · **Owner:** alex
**Evidence:** [current-state audit](../research/2026-09-09-current-state-audit.md), [Web UI audit](../research/2026-09-09-web-ui-audit.md), [competitor research](https://github.com/datatug/marketing/tree/main/competitors)
**Hierarchy:** [pillars → capabilities → features](hierarchy.md) · **Roadmap:** [ROADMAP.md](../../ROADMAP.md)

This document is the one place that says what DataTug is. Every Feature under
[`spec/features/`](../features/README.md) hangs off one of the pillars below. When a
Feature and this page disagree, fix one of them the same day; do not let both stand.

## Thesis

> **DataTug is a context-aware, semantic data exploration and troubleshooting
> environment over heterogeneous operational data sources.** It turns an
> organization's distributed operational data into a **navigable graph of meaning,
> relationships, reusable knowledge and executable investigation paths** — stored as
> Git-native project files, executed by a trusted server, usable identically by people
> in a terminal, people in a browser, and AI agents.

The claim that distinguishes DataTug from the products it overlaps with:

> In most tools the user builds an application (or a dashboard, or a catalog) *over*
> the data. In DataTug the data and the organization's accumulated knowledge about it
> *are* the application. A value on screen is a navigation point; a saved query is a
> reusable, parameterized asset that knows when it applies; an investigation carries
> its context with it.

And the claim about what compounds:

> DataTug does not merely help people ask arbitrary questions. It accumulates and
> surfaces the organization's **proven ways of answering recurring questions and
> investigating recurring problems**, so the tool gets more useful the longer a team
> uses it — and so an AI agent can follow the same proven paths instead of guessing
> SQL from a physical schema.

## Fundamental abstraction

Competitors are built around one abstraction: Retool around the *app*, Metabase around
the *question/dashboard*, DBeaver around the *connection*, Atlan around the *catalog
asset*, Grafana around the *time series panel*, Postman around the *request collection*,
Databricks Genie around the *agent*.

DataTug's fundamental abstraction is the **executable knowledge project**:

| Layer | What it holds | Where it lives |
|---|---|---|
| **Physical** | Sources, environments, servers, scanned schemas, foreign keys | project files (`environments/`, `servers/`, `dbmodels/`) |
| **Semantic** | Entities, fields with meaning (`Customer.ID`), physical↔semantic mappings, relationships that may cross sources | project files (`entities/`) |
| **Executable** | Parameterized queries with semantic parameters, query collections, investigation playbooks | project files (`queries/`, playbooks) |
| **Session** | Investigation context (the values, filters and time window the user is currently carrying) | the surface (web/CLI/agent), optionally saved to the project |

The **project** is the unit of accumulation and sharing (Git-native, promotable from
personal → team → project → organization). The **investigation** is the unit of use: a
context-carrying session that moves through the graph. The **trusted server** (the
DataTug CLI in `serve` mode, or a hosted DataTug service) is the unit of enforcement:
every read and write passes through it, and access policies are decided there.

## Pillars

Each pillar is a top-level Capability Feature. The state column is the audited state on
2026-09-09 (see the [hierarchy](hierarchy.md) for capability-level detail).

| # | Pillar | One-line definition | Audited state |
|---|---|---|---|
| P1 | **Portable query model (DTQL + DALgo)** | One logical query AST (DALgo `StructuredQuery`), serialized as DTQL YAML, rendered to each backend's native language; native text as an escape hatch. | DTQL + DALgo access are Stable in `dal-go/dalgo`; DataTug consumes them only in `datatug query run`. |
| P2 | **Heterogeneous operational sources** | Explore data where it lives: SQL, document stores, Git-backed data (inGitDB), HTTP/GraphQL APIs, logs. | SQL Server/SQLite/MySQL/Postgres scan via CLI; Firestore and inGitDB via DALgo; HTTP modeled, not executed in web. |
| P3 | **Semantic entities and fields** | Entities and typed fields distinct from tables/columns; explicit, inspectable mappings from physical locations. | Model exists (`Entity`, `EntityField.NamePatterns`, `EntityFieldRef` tags); no runtime resolves mappings. |
| P4 | **Meaning inference** | Suggest field meaning from names, types, neighbours and conventions; suggestions are inspectable and correctable, never silently authoritative. | Modeled as data; no inference engine. |
| P5 | **Relationships and related-data navigation** | A selected value becomes a navigation point to related records, across sources, beyond declared foreign keys. | FK graph scanned bidirectionally; web table page shows FK-driven related cards; no cross-source or semantic lookup. |
| P6 | **Investigation Context** | Portable investigation state (`Customer.ID=C123`, `Environment=prod`, time window) that follows the user and binds compatible query parameters automatically. | Not implemented (web "Variables" page is an empty shell; `Variables` model exists in core). |
| P7 | **Shared executable knowledge library** (primary) | Parameterized queries with semantic parameters, query applicability ("useful right now"), collections, playbooks, entry points; the explore→capture→share→reuse loop. | Queries + typed parameters + folders exist and are Git files; applicability, auto-binding and playbooks do not exist. |
| P8 | **Sharing scopes (Sneat spaces)** | Personal → team → project → organization scopes and promotion, reusing Sneat's space membership model. | Not implemented; `datatug_projects` uses a flat `userIDs[]`; Sneat space packages are declared but unused. |
| P9 | **Git-native projects** | Every knowledge asset is a readable file: versioned, diffable, reviewable, portable; no secrets in project files. Project stores: local directory, GitHub repository, remote OpenVaultDB vault, DataTug Cloud — and the project store itself is explorable inside DataTug. | Filestore is the only storage; `--git` flag on mutating CLI commands; credentials can still land in query files; OpenVaultDB store and self-exploration are founder additions of 2026-09-09, not started. |
| P10 | **Server-enforced access control** | Source, row and column policies decided by the DataTug server through DALgo; the browser is never trusted. | Real in the CLI (`query run` through `access.SecureReadSession`); absent from `serve` and the web UI; layered multi-owner ACL is designed and implementing in DALgo/OpenVaultDB/InGitDB. |
| P11 | **Three surfaces, one project** | TUI, web (via `datatug serve` or hosted), AI agent (skills, MCP) over the same project with screen-level hand-off. | TUI shipped; web builds but the agent link is broken at both ends; skills plugin shipped; MCP unbuilt though marketed. |
| P12 | **Provenance** | Every value can say where it came from: source, table/endpoint, field, query, transformation, mapping, policy. | Not implemented beyond schema docs. |
| P13 | **Infrastructure inventory** | Servers, environments, databases, endpoints and services linked to sources so context can climb from a record to its infrastructure. | Environments/servers/catalogs modeled and scanned; no services/endpoints inventory. |
| P14 | **AI uses the same model** | Agents consume the same entities, mappings, library, playbooks, context and ACL as humans; no disconnected chat architecture. | Skills plugin drives the CLI; `serve-brokered-query-builder` Capability approved, unbuilt. |
| P15 | **Dashboards as a surface** (Dashboardius) | Boards emerge from exploration and are preserved as project files; the standalone Dashboardius product presents DataTug boards. | Board schema modeled in Go and TS; web board service unimplemented; Dashboardius demo live on its own domain. |

## Deliberate non-goals

- A generic workflow engine. Playbooks are human-guided investigation paths; automation
  is a later consequence, not the design centre.
- A no-code app builder (Retool/Appsmith). DataTug never asks the user to build an app
  before they can look at their data.
- A data catalog or governance suite (Atlan/DataHub). DataTug records meaning in order to
  *navigate and execute*, not to certify and govern.
- A BI/warehouse tool. DataTug explores operational data where it lives; it does not
  centralize.
- Connector breadth as a growth lever. The minimum set that proves the architecture
  (SQL via DALgo, inGitDB, Firestore, HTTP) comes first.
- Client-direct database access for protected sources. All reads go through the server.
- A text-to-SQL chat box. AI enters through the same knowledge model as people.

## Hypotheses tested against evidence

| # | Hypothesis (from the planning brief) | Verdict | Evidence |
|---|---|---|---|
| H1 | The strongest differentiation is the *combination* of semantic understanding, contextual exploration and reusable executable team knowledge — not any single query technology. | **Confirmed.** No product researched combines P3–P7; each pillar individually has a stronger incumbent. | competitor research: closest per-pillar competitors differ (Datasette/Azimutt for P5, Splunk/Kibana for P6/P7, Atlan for P3, Bruno/Postman for P9-style files). |
| H2 | The context basket can be a first-class Investigation Context rather than a filter bar. | **Confirmed, with a scope rule.** It is portable state with explicit constraints; it must never become a hidden filter. | Datadog saved views, Splunk workflow actions and Grafana variables all stop at "filter"; none carries semantic values across query assets. |
| H3 | Parameterized queries become substantially more valuable when parameters have semantic meaning and bind automatically. | **Confirmed and cheap.** `ParameterDef.Meta` (an `EntityFieldRef`) already exists in the model; only the resolver and UI are missing. | `datatug-core/pkg/datatug/parameters.go`. |
| H4 | Collections can evolve into human-guided playbooks without a workflow engine. | **Plausible; unproven.** Keep playbooks as a later phase; the folder tree already exists. | `QueriesFolder` recursion; no product researched offers human-guided branching playbooks over data. |
| H5 | Sneat team/project organization gives a natural sharing structure. | **Confirmed for the model, not for the code.** Space membership/roles exist in Sneat; DataTug does not use them yet. | ecosystem audit §5. |
| H6 | Git-native storage attracts technical teams. | **Confirmed by competitor behaviour.** Bruno, DataGrip, dbt, Grafana-as-code and Lightdash all moved toward files; Retool/Superblocks gate Git behind Enterprise tiers. | competitor research. |
| H7 | The web UI should expose the data itself as navigable, not require an internal app. | **Confirmed.** This is the line against the internal-tool category. | internal-tools research. |
| H8 | A trusted server enforcing source/row/column ACL differentiates from browser-centric tools and enables self-hosting. | **Confirmed in DALgo, unproven in the web path.** The CLI already enforces; `serve` does not. | DALgo/ACL audit. |
| H9 | DataTug can bridge DB exploration, catalog concepts, troubleshooting, internal-tool UX and agent context without replacing those categories. | **Confirmed as positioning; risk is scope.** The roadmap holds the line with explicit non-goals. | this page. |
| H10 | The compounding moat is the organization-specific graph of mappings + relationships + contextual queries + playbooks. | **Confirmed as the bet; not yet demonstrable.** Nothing today captures or reuses that graph at runtime. | audits. |

## Brand and product boundaries

- DataTug remains a Tier-3 independent brand (founder ruling, restated 2026-07-01); it
  integrates with Sneat Work as `/data`, and reuses Sneat identity (`sneat-eur3-1`).
- Dashboardius is the presentation product over DataTug boards; DataTug keeps the board
  schema and the query/knowledge model. See the roadmap's Dashboardius section for the
  contract between them.

## Open Questions

- **UI kit for DataTug and Dashboardius (Ionic, PrimeNG, or both).** The fleet standard (founder, 2026-08-19) makes Ionic primary with PrimeNG islands for data surfaces in Sneat apps; Dashboardius shipped as a standalone PrimeNG app. DataTug is a Tier-3 desktop-first data tool whose current Ionic shell is thin. Recommendation and trade-offs are in [ROADMAP.md](../../ROADMAP.md#ui-kit-decision); the founder decides.
- **Home of the DataTug↔Dashboardius board contract** — publish `@datatug/board-models` from `datatug-apps`, or keep the documented copy until Phase 3?
- **Does a Dashboardius board belong to a DataTug project or to a Sneat space?** (raised in the Dashboardius Idea; must be settled before board persistence is built.)
