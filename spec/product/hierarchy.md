# DataTug capability hierarchy

## Reading order and classification

The product outcome is repeatable operational investigation. The project holds durable
knowledge; the investigation is the unit of use. Start with the core loop, then inspect
its capabilities, the architecture enabling them, and the surfaces exposing them.
P1–P15 are stable reference IDs, not fifteen equally urgent product priorities.

1. **Core loop:** useful result → meaningful selection → related data/applicable query
   → confirmed bindings → authorized result → retained context → continuation.
2. **Primary capabilities:** Investigation Context (P6) and Shared Executable Knowledge
   (P7), supported by semantics (P3), relationships (P5) and provenance (P12).
3. **Enablers:** query representation, adapters, policy boundary, source registry and
   Git-native project assets (P1/P2/P9/P10).
4. **Surfaces/distribution:** web, CLI/TUI, AI consumers and Dashboardius (P11/P14/P15).
5. **Extensions/experiments:** probabilistic detection, infrastructure exploration,
   hosted collaboration, richer semantics and automated playbooks after evidence.

An Idea records a problem or hypothesis. A Feature owns observable behavior and ACs;
subfeatures refine it. A plan task names implementation work and Verifies those ACs.
Experiments have measured outcomes, not an Implemented label. Non-goals bound a named
phase unless the founder explicitly made a permanent product decision.

## Capability map and current evidence

The table summarizes the 9 September independent review, not continuously refreshed
remote state. Refresh its receipts before implementation; component presence is not
journey acceptance. [Roadmap](../../ROADMAP.md) owns sequencing; [Phase 1 plan](../plans/2026-09-09-phase-1-core-investigation-loop.md)
owns corrective tasks. Earlier detailed audit tables remain historical evidence in
[research](../research/README.md), not an alternative live state matrix.

| ID | Owning Feature | Kind | Evidence boundary | Sequence |
|---|---|---|---|---|
| P1 | [Portable query model](../features/portable-query-model/README.md) | enabler | DTQL and saved-query models exist; binding/transport acceptance pending | 1 |
| P2 | [Heterogeneous sources](../features/heterogeneous-sources/README.md) | capability | SQLite/inGitDB paths and HTTP adapter exist; all-three browser acceptance pending | 1 |
| P3 | [Semantic layer](../features/semantic-layer/README.md) | capability | Declared mappings and resolver exist; integrated grid/endpoint contract pending | 1 |
| P4 | [Meaning inference](../features/meaning-inference/README.md) | capability/experiment | Deterministic inferred mappings exist; review/correction later; probabilistic inference unproved | 2 / extensions |
| P5 | [Related-data navigation](../features/related-data-navigation/README.md) | core interaction | Engine/endpoints/components exist; shared registry, counts and wire compatibility pending | 1 |
| P6 | [Investigation Context](../features/investigation-context/README.md) | core capability | Context components exist; scoped typed state/conflict acceptance pending | 1 |
| P7 | [Shared executable knowledge](../features/executable-knowledge-library/README.md) | primary capability | Queries/folders/resolver exist; complete discovery and reuse journey not accepted | 1–3 |
| P8 | [Project grants and sharing](../features/sharing-scopes/README.md) | collaboration capability | Standalone project with multi-space grants is explicit direction; acceptance pending | 4 |
| P9 | [Git-native projects](../features/git-native-projects/README.md) | storage enabler | Filestore and credential validation exist; OpenVaultDB store/self-exploration future | 1–2 / extensions |
| P10 | [Server-enforced access](../features/server-enforced-access/README.md) | security enabler | Secure-read components landed; stronger layered VM chain pending convergence | 1 prerequisite / later UI |
| P11 | [Surfaces](../features/surfaces/README.md) | distribution/surface | CLI/TUI and browser exist; core browser journey incomplete; hosted path needs receipt | 1 / separate parity |
| P12 | [Provenance](../features/provenance/README.md) | explainability capability | Basic mode/policy metadata exists; transport alignment and richer capture history pending | 1–2 |
| P13 | [Infrastructure inventory](../features/infrastructure-inventory/README.md) | extension | Environment/server/catalog models exist; linked service/record exploration unproved | extensions |
| P14 | [AI consumes the same model](../features/agent-context/README.md) | consumer/extension | CLI skills exist; agent context and brokered AI surfaces require bounded plans | extensions |
| P15 | [Dashboardius dashboards](../features/dashboards/README.md) | presentation surface | Demo interactions and board models exist; real shared persistence not accepted | 5 |
| P16 | [Incidents (Incidentius core)](../features/incidents/README.md) | product capability | Nothing exists; specified 2026-09-11 from the founder's vision | B1 |
| P17 | [Evidence records and result snapshots](../features/evidence-records/README.md) | evidence enabler | Seed only: `Result.provenance` and HTTP fixture snapshots; no execution record or store | B1 |
| P18 | [Annotations](../features/annotations/README.md) | knowledge capability | Nothing exists | B1 |
| P19 | [Executable checks](../features/checks/README.md) | executable-knowledge capability | Dataset value checks exist in `datatug-core` (`checks.go`); no purposes, predicates or runs | B1 |
| P20 | [Compare](../features/compare/README.md) | execution primitive | Schema compare exists (`pkg/comparator`); no recordset or cohort compare | B1 (API) / B2 (UI) |
| P11+ | [Product profiles](../features/product-profiles/README.md) | surface | Nothing exists; founder ruling 2026-09-11 puts all apps in one Nx workspace | B1 (profiles) / C3 (cutover) |

## Existing ideas and overlapping feature reconciliation

Retain source IDs/history. This is a scope/ownership mapping, not an unauthorized
lifecycle transition or deletion of an older idea.

| Existing artefact | Owning capability | Disposition |
|---|---|---|
| CLI `spec/ideas/shared-module-system.md` | P7 library + P9 project assets | Future packaging/reuse input; no new module engine in Phase 1. |
| CLI `spec/ideas/semantic-metadata-cli.md` | P3 semantic layer | Per-surface authoring of the shared model, not a second semantic authority. |
| CLI semantic-field-types currency seed | P3/P4 semantic types | Reuse type declarations; do not infer a currency from CustomerId. Rich types remain an extension. |
| CLI `spec/ideas/cross-engine-db-copy.md`, `db-copy-filtering.md` | P1/P2 source/query contracts | Separate data-movement workflow; never implied by read-only investigation acceptance. |
| CLI `spec/ideas/db-snapshot-command.md` | P9/P12 capture/provenance | Separate snapshot product capability; explicit HTTP demo snapshots do not implement it. |
| Hub `ai-query-builder` | P14 | In-app assisted query UX consuming the shared server contract. |
| Hub `serve-brokered-query-builder` | P14 | Broker/transport contract; not an alternative query or policy engine. |
| Apps `query-builder`, `ai-query-builder`, `ai-terminal-query-builder` Features | P11/P14 | Consumer/surface specifications: inventory overlap before an AI task, retain distinct UI modes, share execution/model contracts. |
| Dashboardius backstage Idea | P15 dashboards | Product provenance; execution/storage boundary lives in the hub Feature and presentation journey in Dashboardius specs. |

## Deliberate scope boundaries

Phase 1 delivers all three approved source types and a protected browser investigation.
It does not add a workflow engine, hosted tenancy, chart renderer, broad connector
catalog or AI-specific execution path. Later shared Dashboardius embedding is compatible
with separate presentation ownership. A standalone project can be granted to multiple
spaces; it is not duplicated beneath each space.

## Open Questions

Questions belong to the owning Feature. Follow the links above for board query binding,
hosted identity, future HTTP credentials and collection metadata. Internal basket naming
and the already-approved source/UI-kit choices do not require another founder decision.
