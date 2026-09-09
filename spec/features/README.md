---
format: https://specscore.md/features-index-specification
---

# Features

Feature specifications for this project.

## Index

| Feature | Status | Description |
|---------|--------|-------------|
| [AI Query Builder](ai-query-builder/README.md) | Approved | The AI Query Builder is a platform-agnostic **Capability**: an in-app conversational query builder where a user seeds a query in natural language and then progressively refines it through a chain of plain-English requests. |
| [Serve-Brokered AI Query Builder](serve-brokered-query-builder/README.md) | Approved | A platform-agnostic **Capability** for an AI query builder whose conversation lives in the user's existing AI-agent terminal while a local DataTug CLI daemon (`serve`) holds the query and a hosted Web UI shows and deterministically edits it. |
| [Portable query model (DTQL + DALgo)](portable-query-model/README.md) | Draft | P1 is the portable query model: one logical query AST (DALgo `StructuredQuery`), serialized losslessly as DTQL YAML, rendered to each backend's native query language, with native text kept as an explicit escape hatch. |
| [Heterogeneous operational sources](heterogeneous-sources/README.md) | Draft | P2 is the promise that DataTug explores data where it lives: SQL databases, Git-backed data (inGitDB), document stores (Firestore), HTTP/GraphQL APIs and, eventually, logs — all through DALgo adapters, with connection secrets kept outside project files. |
| [Semantic entities and fields](semantic-layer/README.md) | Draft | P3 is DataTug's semantic layer: entities and typed fields distinct from tables and columns, with explicit, inspectable mappings from physical locations, and semantic tags carried on query parameters and result columns. |
| [Meaning inference](meaning-inference/README.md) | Draft | P4 is meaning inference: suggesting a field's semantic meaning from column names, types, neighbouring columns and conventions (and, later, safe data sampling), with every suggestion inspectable and correctable — never silently authoritative. |
| [Relationships and related-data navigation](related-data-navigation/README.md) | Draft | P5 is related-data navigation: a selected value becomes a navigation point to related records, across sources, using declared foreign keys and semantic mappings, beyond what a single table's schema can show on its own. |
| [Investigation Context](investigation-context/README.md) | Draft | P6 is the Investigation Context: portable state (semantic values, environment, time window) that follows the user across screens and automatically binds compatible query parameters, always visibly and never as a hidden filter. |
| [Shared executable knowledge library](executable-knowledge-library/README.md) | Draft | P7 is the primary pillar of the thesis: a shared, executable knowledge library — parameterized queries with semantic parameters, query applicability ("useful right now"), collections, capture from exploration, investigation playbooks, and a promotion lifecycle — the explore → capture → share → reuse loop meant to make DataTug get more useful the longer a team uses it (hypothesis H10, "the compounding moat"). |
| [Sharing scopes on Sneat spaces](sharing-scopes/README.md) | Draft | P8 is sharing scopes: personal → team → project → organization scopes for knowledge assets, reusing Sneat's existing space membership and role model rather than inventing a parallel one. |
| [Git-native projects](git-native-projects/README.md) | Draft | P9 is Git-native projects: every knowledge asset is a readable, versioned, diffable, portable file — no secrets in project files — and the project store itself, wherever it lives, is explorable inside DataTug. |
| [Server-enforced access control](server-enforced-access/README.md) | Draft | P10 is server-enforced access control: source, row and column policies decided by the DataTug server through DALgo's `access` package, with the browser never a trusted enforcement boundary. |
| [Three surfaces, one project](surfaces/README.md) | Draft | P11 is the "three surfaces, one project" claim: TUI, web (local agent or hosted) and AI agents all work over the same project files, with screen-level hand-off between them, and no surface holds private state the others cannot see. |
| [Provenance](provenance/README.md) | Draft | P12 is provenance: every value and every mapping can say where it came from — source, table or endpoint, field, query, transformation, mapping and any policy limitation applied. |
| [Infrastructure inventory](infrastructure-inventory/README.md) | Draft | P13 is infrastructure inventory: environments, servers, databases, endpoints and services linked to sources, so an investigation can climb from a record to the infrastructure that produced it — "this order came from `prod`'s `orders-service` on `db-primary-3`." It is a later-phase extension of the navigation loop P5 establishes for data records, applied one level up to the systems that host the data. |
| [Agent-readable project context](agent-context/README.md) | Draft | P14 is the claim that AI agents use the same model people do: the same entities, mappings, library queries, investigation context and access decisions, exposed through CLI verbs and MCP resources rather than a disconnected chat architecture (deliberate non-goal: "a text-to-SQL chat box... |
| [Dashboards as a surface](dashboards/README.md) | Draft | P15 is dashboards as a surface: boards emerge from exploration and are preserved as project files, with DataTug owning the board schema, persistence and query execution while the separate, standalone Dashboardius product owns presentation and editing. |
| [Core investigation loop (Phase 1)](core-investigation-loop/README.md) | Draft | Phase 1 makes DataTug's defining loop work end to end, at a basic but coherent level, in the web UI over the local agent (`datatug serve`), on one demo project with two sources: |

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/features-index-specification*
