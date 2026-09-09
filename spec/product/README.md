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

Useful comparison abstractions (not exhaustive descriptions of these products) include: Retool around the *app*, Metabase around
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

The **project** is the unit of accumulation and sharing (Git-native, standalone and shareable through explicit grants to multiple spaces). The **investigation** is the unit of use: a
context-carrying session that moves through the graph. The **trusted server** (the
DataTug CLI in `serve` mode, or a hosted DataTug service) is the unit of enforcement:
every read and write passes through it, and access policies are decided there.

## Pillars

P1–P15 are stable capability IDs, not equal-priority product pillars. P6 and P7
organize the primary user value; other capabilities, enablers and surfaces support it.
See the [hierarchy](hierarchy.md) for current evidence, classification and old-idea
reconciliation. The initial-audit absence claims are not current implementation state.

| ID | Capability | Role |
|---|---|---|
| P1 | [Portable query model](../features/portable-query-model/README.md) | enabler |
| P2 | [Heterogeneous sources](../features/heterogeneous-sources/README.md) | capability |
| P3 | [Semantic layer](../features/semantic-layer/README.md) | capability |
| P4 | [Meaning inference](../features/meaning-inference/README.md) | capability/experiment |
| P5 | [Related-data navigation](../features/related-data-navigation/README.md) | core interaction |
| P6 | [Investigation Context](../features/investigation-context/README.md) | core capability |
| P7 | [Shared executable knowledge](../features/executable-knowledge-library/README.md) | primary capability |
| P8 | [Project grants and sharing](../features/sharing-scopes/README.md) | collaboration capability |
| P9 | [Git-native projects](../features/git-native-projects/README.md) | storage enabler |
| P10 | [Server-enforced access](../features/server-enforced-access/README.md) | security enabler |
| P11 | [Surfaces](../features/surfaces/README.md) | distribution/surface |
| P12 | [Provenance](../features/provenance/README.md) | explainability capability |
| P13 | [Infrastructure inventory](../features/infrastructure-inventory/README.md) | extension |
| P14 | [AI consumes the same model](../features/agent-context/README.md) | consumer/extension |
| P15 | [Dashboardius dashboards](../features/dashboards/README.md) | presentation surface |

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
  (SQLite via DALgo, inGitDB, keyless HTTP JSON GET) comes first.
- Client-direct database access for protected sources. All reads go through the server.
- A text-to-SQL chat box. AI enters through the same knowledge model as people.

## Hypotheses tested against evidence

Desk research tests plausibility and overlap, not demand or defensibility. None of
these rows is a completed customer experiment.

| ID | Hypothesis | Evidence-qualified assessment |
|---|---|---|
| H1 | The combination is more valuable than isolated features | Plausible; overlapping competitors exist. Test complete repeated investigations against current workflows. |
| H2 | Semantic context improves on repeated manual filtering | Plausible; notebook variables and query reuse are meaningful comparators. No hidden filters is a product rule, not market proof. |
| H3 | Semantic parameters make reuse easier | Model/resolver infrastructure exists; integration and second-user reuse must demonstrate value. |
| H4 | Collections can support human-guided playbooks | Plausible extension; prove a short sequence before adding orchestration. |
| H5 | Space membership can support project sharing | Existing platform primitives help; standalone-project grants and revocation still require acceptance. |
| H6 | Git-native storage attracts technical teams | Competitor adoption is supporting evidence; it does not establish DataTug demand. |
| H7 | Direct data navigation avoids app-building overhead | Deliberate UX direction to test with time to first useful result. |
| H8 | A trusted execution boundary enables safe operational access | Architectural requirement; partial implementation and pending ACL convergence need end-to-end proof. |
| H9 | DataTug can bridge adjacent categories without replacing them | Positioning choice with significant scope risk; roadmap keeps a bounded investigation first. |
| H10 | Team mappings and executable knowledge compound | Retention/defensibility hypothesis; measure capture cost and another engineer's actual reuse. |

Run a bounded customer/support investigation experiment. One engineer captures a
useful query/mapping; a second repeats and extends the investigation. Record time to
useful result, manual copying, incorrect bindings, capture effort and repeat reuse.
Compare with the team's present SQL/API/notebook/chat workflow. Agree success/failure
thresholds before the experiment; no thresholds or results are currently established.

## Brand and product boundaries

- DataTug remains a Tier-3 independent brand (founder ruling, restated 2026-07-01); it
  integrates with Sneat Work as `/data`, and reuses Sneat identity (`sneat-eur3-1`).
- Dashboardius is the presentation product over DataTug boards; DataTug keeps the board
  schema and the query/knowledge model. See the roadmap's Dashboardius section for the
  contract between them.

## Open Questions

- Board query-target/binding choice: [dashboards Feature](../features/dashboards/README.md#open-questions).
- Hosted grant and identity details: [sharing-scopes](../features/sharing-scopes/README.md#open-questions).
- Web privacy/license interpretation before visibility changes: [surfaces](../features/surfaces/README.md#open-questions).
The Ionic frame/islands rule, standalone project ownership, three MVP sources and CLI
Apache 2.0 direction are already recorded in the roadmap; do not reopen them implicitly.
