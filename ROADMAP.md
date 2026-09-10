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

## Phase 0 — existing spine and delivery repair

**Outcome:** `datatug demo` opens the actual project; source → Album gives real rows
without hidden setup. The CLI builds, validates project credentials and serves a
registered project. Historical build/store/path/auth fixes are implementation evidence,
not remaining greenfield tasks. The review built current CLI and opened the local app;
project summary returned 200. Do not continue blaming the old auth-hook panic without
reproducing it on the candidate.

**Gate:** real entry-point browser test plus an exact client/server/demo revision
manifest. Phase 1 corrective tasks 12–16 are complete (2026-09-11); Task 17 holds the pinned release candidate (datatug-cli v0.23.0, datatug-apps 8653bd9) pending the live-site receipt.
**Scope:** no new product concepts, no unrelated command cleanup or UI-kit migration.

## Phase 1 — converge the core investigation loop

[Feature](spec/features/core-investigation-loop/README.md) ·
[transport contract](spec/features/core-investigation-loop/api-contract.md) ·
[plan](spec/plans/2026-09-09-phase-1-core-investigation-loop.md).

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
limitation header are all fixed. Still open: task 13 (protected-execution
convergence with the implementer agent's layered-ACL chain — ownership to be
settled), task 14 (browser-side bounded HTTP execution; the provider-side bounds
landed in dalgo2http v0.2.0), task 16 (fixture policies, Codex) and task 17 (one
reproducible release candidate and publication, after 16). Refresh remote state
before implementing.

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

[Draft plan](spec/plans/2026-09-09-phase-2-library-depth-and-provenance.md).

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
shared board contract, explicit query-target binding decision, server write and
concurrent-edit policy, and recipient authorization.
**Gate:** a genuine cross-language board round trip and real persisted two-actor
journey; scripted sample cards do not satisfy it.
**Autonomy:** blocked only for dependent work by the explicit SQL-widget binding
question in the Feature. Presentation-only refinements may proceed independently.

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
DataTug CLI is Apache 2.0; the web UI is private for now, with future licensing not
selected. Do not market all DataTug surfaces as already open source.

## UI kit decision

Founder words, 9 September: "UI-kit rule - confirmed". DataTug keeps the Ionic frame
and bounded islands for data surfaces it lacks, reusing the existing shared grid.
Dashboardius keeps its standalone Angular/PrimeNG presentation. Neither journey
requires a UI-kit migration. Shared presentation libraries may serve both products;
product ownership does not dictate a permanent ban on embedding.

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

## Open Questions

Product choices live beside their owning Feature: board SQL-widget binding in
[dashboards](spec/features/dashboards/README.md#open-questions), hosted identity and
sharing in [sharing-scopes](spec/features/sharing-scopes/README.md#open-questions),
and distribution/privacy interpretation in [surfaces](spec/features/surfaces/README.md#open-questions).
No release version, provider change or deployment success is inferred from this plan.
