# Independent review: DataTug and Dashboardius, 9 September 2026

## Executive assessment

**Verdict: REVISE PHASE 1 BEFORE IMPLEMENTATION.** This means revise before assigning another autonomous implementation wave; implementation already exists and should be reconciled, not restarted.

Fable 5.1 did a strong job identifying and explaining the product thesis, a useful job organizing a large research corpus, and a materially weaker job turning that understanding into a consistent, executable integration plan. The volume is substantial, but completion labels, competitive certainty, and interface specificity outrun the evidence. This is not a reason for major product or architecture rework.

Preserve the central proposition: a project contains semantic mappings and reusable executable knowledge; an investigator accumulates meaningful context, discovers applicable knowledge, and continues across heterogeneous sources through a secure server. Investigation Context and Shared Executable Knowledge are correctly prominent. The weak point is the bridge from this proposition to one working browser journey.

Before further implementation, correct the Phase 1 wire contracts, execution target and binding rules, context scope, and protected-read acceptance. Then integrate and verify the substantial merged and pending code against those contracts. Do not start the present Phase 2 merely because several Phase 1 task rows say Complete.

This review made no implementation changes, commits, merges, branch changes, or worktree cleanup. Recommendations below are not actions already performed and do not authorize discarding pending work. Review documentation is the only repository content changed by this exercise.

## Scope, method, and evidence limits

The source brief was read first: `/Users/alex/projects/datatug-frontier-research-planning-prompt.md`. Review order was vision, documentation inventory, product/hierarchy, competitor methodology and spot checks, roadmap and Phase 1, Dashboardius, then implementation verification. The review did not substitute a code audit for the requested product review.

Discovery covered Mac repositories and their worktrees, plus the development VM at `ai@178.104.41.143`, principally under `/home/ai/projects`. Git history and status distinguish committed, staged, unstaged, untracked, and already-integrated work. Filesystem dates are discovery hints, not proof of authorship or completion. “Fable work” below means the deliverables under review, not a forensic assertion that a particular model authored every line.

Remote main refs and selected PR conclusions were checked directly during the review, around 14:06 UTC. This is a dated snapshot: another active session can change them. Ancestry differences against a stale local main are not treated as proof that a change is unmerged. Pending ACL acceptance claims are identified as recorded evidence rather than tests independently rerun here.

The review inventories the broad corpus and reads the central documents deeply. It samples individual competitor reports and verifies selected important claims; it does not independently re-research all 84 ranked entries. Authentication and real persisted Dashboardius sharing were not completed. No production data or user account was modified.

## Artefacts reviewed and canonical ownership

| Artefact family | Location | Assessment and canonical role |
|---|---|---|
| Source exercise and vision | `/Users/alex/projects/datatug-frontier-research-planning-prompt.md` | Governing brief for this review; not a second implementation spec. |
| Product definition and hierarchy | [product README](../product/README.md), [hierarchy](../product/hierarchy.md), hub README and spec indexes | Canonical product map. Keep and simplify its classifications. |
| Capability specifications | `spec/features/`, especially core-investigation-loop, investigation-context, executable-knowledge-library, heterogeneous-sources, sharing-scopes, server-enforced-access, dashboards | Canonical behavior and ownership; correct contradictions in place. |
| Roadmap | [ROADMAP](../../ROADMAP.md) | Canonical sequence, but not an authoritative delivery receipt. |
| Phase 1 | [Feature](../features/core-investigation-loop/README.md), [implementation plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-09-phase-1-core-investigation-loop/README.md) | Correct these before the next implementation wave. |
| Phase 2 | [library depth and provenance plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-09-phase-2-library-depth-and-provenance/README.md) | Useful backlog, not yet an executable successor phase. |
| Current state and browser research | [current-state audit](2026-09-09-current-state-audit.md), [Web UI audit](2026-09-09-web-ui-audit.md) | Dated evidence. Reconcile later landings without rewriting historical observations as current truth. |
| Competitor research | `datatug/backstage/competitors/`: README, methodology, feature matrix, conceptual hunt, positioning, opportunities, categorized individual reports | Canonical private research. Preserve evidence; downgrade unsupported certainty. |
| Dashboardius product material | `datatug/backstage/marketing/dashboardius/README.md`; hub dashboards Feature; `sneat-co/backstage/spec/ideas/dashboardius.md` | Hub Feature owns the integration contract. Dashboardius should own its own user-facing feature specifications. |
| Dashboardius implementation record | `sneat-co/dashboardius/README.md`, `AGENTS.md`, `docs/IMPLEMENTATION-RECORD.md` | Useful prototype evidence; repository code/history is from 8 September, not a new 9 September implementation. |
| Implementation-local specifications | CLI/apps spec trees, DALgo/DTQL/provider specs, HTTP adapter examples | Own local contracts and verification, with links to the hub Feature. Promised CLI/apps core-investigation-loop Feature files are absent. |
| Pending VM delivery evidence | OpenVaultDB SDK layered-ACL worktree `docs/layered-acl-delivery.md`, `layered-acl-acceptance.md`, `layered-acl-review-checkpoints.md` | Essential implementation/dependency evidence; not a replacement for the DataTug product spec. |

Today's hub changes include the initial hierarchy/research/roadmap package, approved Phase 1 Feature and plan, founder source/licensing/scope decisions, Phase 2 drafting, task-status reconciliations, and board-schema decisions. The most recent board decision exists in the Mac `spec-discoverability` worktree and remote main, ahead of the canonical Mac checkout. It was included in this review.

The VM hub and marketing main checkouts are older than the Mac corpus. The VM's distinctive contribution is pending ACL integration and CLI follow-ups, not a competing newer research corpus. This matters: searching only VM main would miss the valuable work; searching only Mac main would miss it too.

## Ideas/features hierarchy

The hierarchy is a meaningful improvement. Someone new can identify the intended product, its core concepts, major surfaces, and proposed direction. Context and executable knowledge have first-class prominence rather than being buried under a query editor or AI feature.

The remaining problem is classification. Fifteen coequal pillars and eighteen top-level features mix user value, infrastructure, distribution, and surfaces. DTQL/DALgo, Git storage, cloud deployment, AI consumption, investigation behavior, and dashboards should not all imply equal product priority. The map needs a short primary journey above the capability inventory.

Recommended reading structure, without creating another parallel hierarchy:

1. **Product outcome:** repeat and continue operational investigations using shared executable knowledge.
2. **Core loop:** obtain a result, understand it, capture context, discover and bind a next query, execute safely, continue.
3. **Capabilities:** semantic mapping, relationships, query collections, context, provenance, and guided scenarios.
4. **Architectural enablers:** source adapters, DTQL/DALgo, policy enforcement, durable project assets, and versioning.
5. **Surfaces and distribution:** DataTug web/CLI/server, Dashboardius presentation, self-hosted and hosted offerings.
6. **Experiments and future opportunities:** assisted semantic detection, AI consumers, broader workflow automation, infrastructure self-exploration.

Retain existing IDs and links. Add a reconciliation table for old CLI ideas such as shared modules, semantic field types/currency, copy workflows, and snapshots. State whether each is incorporated, deferred, or superseded. Similar AI/query-builder ideas across repositories need explicit ownership and relationships, not an implied new canon that leaves old promises intact.

Implementation tasks belong in plans, not as additional product pillars. A capability should state an observable user result, then link to its enabler and tasks. The current hierarchy's implementation matrix also needs to stop pointing at local Feature files that do not exist.

## Product-definition assessment

Fable understood the product's defining combination. The best organizing unit is an **investigation**, with a project as the durable knowledge boundary. A database browser, query editor, semantic catalog, dashboard, and AI assistant are ways into or through that investigation, not five independent products inside DataTug.

The concepts can form a clean causal chain:

> Physical data is mapped to semantic fields → a selected record contributes typed facts → facts make saved queries applicable → explicit bindings select a safe execution → results add evidence and new context → the next investigator can reuse the knowledge.

This chain is stronger than a feature checklist. Automatic relationship discovery and semantic detection should produce explainable candidates; they should not silently invent business identity. Manual mappings are a sensible first slice. Browser joins may compose only authorized results; server-side predicates, joins, aggregates and source access require their own enforcement. A browser join is not an ACL implementation.

Infrastructure inventory and Git-native metadata are useful inputs and operational advantages. They should support the investigation loop before becoming a parallel infrastructure-management roadmap. AI agents should eventually use the same query/context/policy contracts; an AI-specific execution path would undo the architectural advantage.

The standalone DataTug project and grants to multiple team spaces are explicit product decisions. The roadmap's later space-rooted storage sketch must be reconciled with the sharing Feature rather than overriding this exception with a generic platform convention. A board belongs to a DataTug project. DataTug CLI Apache 2.0 and a private web UI are also explicit decisions, not open research questions.

## Shared Executable Knowledge assessment

The prominence is right. Parameterized queries, discoverability from context, collections, semantic parameter meaning, and guided reuse are the strongest parts of the strategy. The query library should be the first reusable asset; a scenario then sequences those assets rather than copying queries into a second execution system.

The specification is not yet precise enough about what makes knowledge executable:

- Required parameters, defaults, explicit overrides, ambiguous semantic matches, missing manual inputs, and type validation need deterministic behavior.
- A query needs a resolvable source/target and execution profile. A text body and a list of semantic parameters are insufficient.
- “Applicable” should distinguish runnable now, needs input, unavailable source, and denied, without revealing protected query or data details.
- The UI should explain which selected fact supplies each binding and permit correction before execution.
- Reuse needs a small quality loop: author/owner, description, example purpose, and a reproducible run. Rich governance can follow later.

Do not postpone every mechanism of sharing while presenting “Save as team query” as an earlier completed workflow. File/Git project sharing is a valid initial boundary if explicitly stated; hosted membership is a separate later boundary.

Human-guided playbooks are a good scope limit. A permanent prohibition on future unattended conditional execution is not established by the vision. Record the present non-goal as a phase limitation, leaving future automation subject to evidence and a later decision.

## Investigation Context assessment

Context is correctly treated as more than a query's parameter form. However, its current specification leaves important behavior to implementers, and implementation has taken incompatible shortcuts.

Define context identity as at least project, store/agent, environment, and principal/security session, with typed values and provenance. Number `5` and string `"5"` must not collapse merely because they stringify alike. A context entry should distinguish a selected fact from an inferred fact, with an explanation for the latter.

Two values for the same semantic field must produce a visible choice or conflict, not “first value wins.” A change of project, server, environment, or principal must have specified reset/isolation behavior. Session persistence alone does not establish that boundary. Time/expiry can be modest initially, but cross-project isolation cannot be postponed because a plan has reached a task-count budget.

The customer journey currently assumes that selecting CustomerId 5 establishes Country and then currency CAD. An ID alone does not contain those facts. Either select a row containing Country, explicitly retrieve the authorized customer record, or ask for the missing country. Then show the HTTP-derived currency fact and its provenance. Correct the fixture's inconsistent expected query counts instead of allowing tests to choose whichever number is convenient.

URL sharing requires more than excluding credentials. Context can contain customer identifiers, emails and operational values. Define explicit safe sharing, recipient authorization, and re-execution; do not turn a link into a portable authorized result or silently serialize the whole basket.

## Architecture/specification assessment

The broad division is sound: web interaction, server execution/security, DALgo/DTQL and adapters beneath, project assets above. Three gaps prevent autonomous implementation readiness.

**First, publish an exact transport contract.** Endpoint names and prose such as “recordset + limitations” have already led the browser and server to implement different shapes. Choose one canonical wire schema, query/body parameter placement, identifiers, error taxonomy, limitation representation, and compatible client representation. Validate real serialization across the boundary. Do not add a third abstraction that leaves both current implementations intact.

**Second, define one source registry and target-resolution contract.** A semantic endpoint cannot discover sources by a private filesystem convention while execution uses project environments/catalogs. SQLite, inGitDB and HTTP all need the same project source identity. HTTP does not naturally have a SQL database/environment pair; the transport must express the actual target without pretending otherwise.

**Third, define trust profiles explicitly.** The Feature's native-SQL carve-out permits source-level access without row/column enforcement, whereas the journey promises ACL-filtered results and the pending protected profile refuses native SQL. A label saying that filtering was skipped is not enforcement. Protected principals must fail closed for execution the policy engine cannot constrain. If an opaque source-level execution grant exists, it must be a deliberately separate privileged capability, with no claim that row/column rules still apply.

This is a specification contradiction, not a demonstrated native-SQL exploit for the current support fixture. The inspected support policy does not grant opaque execution. Preserve that distinction.

Also specify:

- Policies for every source/table used by the support journey, including invoices and support notes; a Customer filter does not automatically become an Invoice policy.
- Bounded, authorized relationship lookups and counts. A client-supplied lookup ID is not authority; counts and metadata can leak inaccessible information.
- Fixed local principal versus authenticated remote principal, role assignment, and authorization of write endpoints. CORS is not authentication, and `--as support` alone does not mean `--role support`.
- HTTP target allowlists and redirect policy, time/size limits, credential handling, and explicit live-versus-snapshot choice. Snapshot-backed acceptance must not masquerade as live HTTP acceptance.

## Competitor-research assessment

The research has real value: broad coverage across application builders, catalogs, SQL tools, BI, API tooling, observability and AI; evidence labels and uncertainty notes; and a useful separation of overall, conceptual, commercial, pillar-specific, strategic-threat and UX-learning competitors. This is substantially better than ranking familiar brands by popularity.

The closeness scoring is not reproducible enough to support a precise league table. D13 appears under two weight groups in the methodology while the matrix chooses one. Unknown/N/A treatment, confidence and sensitivity are not adequately specified. Analyst judgments and matrix scores diverge markedly—for example, DataGrip and Coginiti rank closer in the qualitative view while Metabase and Hasura score much higher in the matrix—without a clear reconciliation. These differences may be legitimate, but the reader needs the reason. Intended DataTug capabilities must be distinguished from current shipping capability when scoring closeness.

Normalize overlapping entries such as Steampipe/Powerpipe and Databricks/Genie. Repair the unfinished D01 label. Keep the generator and inputs reproducible rather than referring to disposable scratch tooling. A smaller calibrated shortlist is more decision-useful than decimal precision across 84 entries.

Most seriously, synthesis sometimes turns an unverified combination hypothesis into certainty: “Nobody else sells that,” “genuine white space,” and hypotheses marked Confirmed. The absence of an exact checklist match does not establish lack of substitutes, customer preference, or a defensible advantage.

### Primary-source spot checks

| Claim or gap | Evidence checked | Review conclusion |
|---|---|---|
| Observability tools lack meaningful cross-step context | [Datadog notebook variables](https://docs.datadoghq.com/notebooks/guide/template_variables_analysis_notebooks/) and [notebooks](https://docs.datadoghq.com/notebooks/) describe reusable variables across cells and data-source/transform/visualization workflows, with documented limits. | “Single-hop saved views” is too dismissive, and the individual report already contains contrary nuance. This does not prove Datadog has DataTug's exact semantic basket. |
| Executable operational knowledge peers adequately covered | [Runme](https://runme.dev/), [session environment behavior](https://docs.runme.dev/configuration/dotenv-direnv/), and [walkthrough](https://docs.runme.dev/resources/walkthrough/) expose executable Markdown and cross-cell operational context. | Add Runme as a conceptual/UX comparator. Its context differs from a semantic data model, but dismissing it would artificially enlarge the proposed white space. |
| AI-data products are mainly generic chat over data | [Databricks Genie quality guidance](https://docs.databricks.com/aws/en/genie-agents/tune-quality) and [Genie interactions](https://docs.databricks.com/aws/en/genie/talk-to-genie) document trusted parameterized assets and reuse. | The research is right to treat Genie as a strategic threat. The synthesis should acknowledge its overlap in reusable executable knowledge. |
| Semantic-layer category is sufficient | dbt/MetricFlow appears in the corpus without a comparable dedicated assessment. | Add a focused semantic-layer comparator rather than implying catalogs alone cover business semantics and metric/query reuse. |

Windmill's screening should also be upgraded from prior-knowledge dismissal to an evidence-backed conceptual comparison. These are targeted corrections, not a request to redo the whole research project.

## Differentiation assessment

The combination hypothesis is persuasive as a product direction and unproven as a commercial advantage. Heterogeneous operational sources plus explainable semantic context and reusable queries can reduce repeated investigation work. Secure execution and Git-native assets make this credible for teams that cannot centralize all their data in a warehouse or application builder.

The asymmetric opportunity is not “implement every missing competitor capability.” It is making a small, repeated investigation reusable by the next engineer without forcing that engineer to learn schemas, copy IDs repeatedly, or reconstruct the previous investigator's reasoning.

Recommended first demand experiment: recruit teams doing recurring customer/order/support investigations; have one engineer capture a useful query and context mapping, then a second engineer repeat and extend the investigation. Measure time to first useful result, manual copying, incorrect bindings, knowledge reuse, and authoring overhead. Compare with the team's existing SQL/API/notebook/chat process. Establish thresholds before running the experiment; there are no validated thresholds or demand outcomes in the reviewed artefacts.

A defensible position would come from accumulated useful mappings and executable team knowledge, trustworthy execution, and low-friction reuse. The feature combination by itself is copyable. Avoid “confirmed” until observed user behavior supports it.

## Roadmap assessment

The initial browser-first vertical slice is the right strategy. The later sequencing is too broad for the requested fast start and first three-to-five phases with little owner involvement. Phase 3 bundles scenarios, sharing/space organization, hosted execution, OpenVaultDB project storage, infrastructure exploration, dashboards and UI consolidation. Those are different product decisions and dependency chains.

The pending ACL stack makes part of the later security work much cheaper than a greenfield roadmap suggests, but it is not yet a generally available provider chain. Conversely, today's landing of individual semantic components does not make the browser slice complete.

Recommended sequence, preserving the three approved MVP source types:

| Sequence | Observable good result | Scope discipline |
|---|---|---|
| Phase 1 repair and convergence | A browser user completes the core investigation with SQLite, inGitDB and HTTP; a restricted user receives only authorized results and an honest limitation/denial. | Freeze contracts, integrate existing work, prove the journey. Do not add a new product surface. |
| Phase 2 executable library quality | Another engineer can find, understand, bind and rerun a saved query with provenance and missing-input guidance. | Collections and reuse before a generalized workflow engine. |
| Phase 3 guided investigation | A small human-driven sequence reuses existing queries/context and gives an explicit end state. | No hosted tenancy or automatic orchestration required to prove this value. |
| Phase 4 collaboration and distribution | A recipient accesses the same standalone project under the correct grant, with revocation and isolation verified. | Choose the hosting/storage lane separately; do not implicitly rewrite project ownership. |
| Phase 5 Dashboardius persistence | A real query-backed card is saved to the DataTug project and reopened by an authorized colleague. | Reuse execution and storage; resolve board binding before promising this autonomous phase. |

Security convergence is a prerequisite to any protected hosted journey, not a late optional feature. UI kit migration and infrastructure self-exploration should not be hidden inside the collaboration phase. They can have bounded independent plans when justified.

## Phase 1 assessment

A new Astra session can use the documents to understand the product without repeating research. It cannot safely execute them as written without making consequential contract and trust decisions. The corrections are larger than editorial cleanup, hence the verdict.

The intended happy path is appropriate: open project, choose source, browse/query, inspect meaningful data, select a fact or record, discover related data or an applicable saved query, inspect bindings, execute through the server, retain context, and continue. The spec must additionally walk the missing-parameter, ambiguous-binding, denied, unavailable-source, and project-switch paths.

Make the early task graph follow actual dependencies: model/contract before endpoints and consumers, source resolution before semantic lookup and execution, protected policy behavior before restricted-user acceptance, and all three source adapters before the complete journey. The existing task graph omits dependencies that task prose actually uses. A numbered task budget cannot justify deferring context isolation or removing an approved source from acceptance.

Some task completion claims refer to component-level work. In particular, the HTTP adapter exists but the browser execution endpoint still refuses HTTP saved queries. ACL unit tests do not establish J4. These tasks should record “component landed; journey acceptance pending” until the exact user outcome passes.

The test plan needs three distinct evidence layers:

1. Engine/policy tests for typing, bindings, denial and source behavior.
2. Real transport tests proving the browser's request and response contract against the server.
3. Browser journeys beginning at the specified entry point, without intercepted fake responses, direct deep-link shortcuts, or disabled security scenarios.

Missing infrastructure must fail the required CI journey job, not skip it successfully. Pin or record the CLI, core/provider and demo revisions used by acceptance. Floating companion main branches make a result hard to reproduce. Deployment should depend on the applicable journey gate for the candidate being deployed.

## Dashboardius assessment

The product split is coherent when stated as ownership rather than permanent UI exclusion.

| Shared contract | DataTug owns | Dashboardius owns |
|---|---|---|
| Board/query references, parameter bindings, authorized result metadata, source provenance, durable serialization | Projects, sources, semantic/query assets, execution, ACL, persistence, grants and Git-backed storage | Dense dashboard presentation, cards, chart/table widgets, rearrangement/removal, hidden metrics, interaction reducers, Undo/Reset and prompt-to-action UI |

Angular/PrimeNG and the reducer-driven interactions are a reasonable presentational foundation. Preserve stable metric slots, chart/table alternatives, reversible actions, and lazy grid-engine boundaries. Prompt-driven interactions must retain the distinction between deterministic suggestions and actual model-generated operations.

The live Dashboardius demo clearly labels its scripted interaction; adding sign-ins to the sign-ups chart worked, and Undo restored it. Save explicitly said saving was not built and changes existed only in the tab. That is honest prototype behavior, not evidence of persistence. The visible invalid PrimeUI license notice needs resolution before a polished public product claim. Authentication and a colleague reopening saved state were not independently demonstrated in this review.

DataTug's own board UI must not become a second independently maintained renderer. Equally, the source brief and separate Dashboardius positioning do not establish a permanent ban on embedding the shared renderer within DataTug. “No boards in Phase 1” is a scope decision; “no chart rendering anywhere in DataTug forever” is a much larger inference. Record an embedding decision only when needed.

Schema convergence is unfinished. The latest core commit already adds `Board.Parameters` and `RequiredParams`; that particular loss has been fixed. Dashboardius still has row IDs and demo widget data that differ from the canonical Go model. A TypeScript JSON stringify/parse round trip does not prove Go load/save compatibility. Go's typed `BoardRow` has no row ID; decode/re-encode can lose the identifier on which the interaction reducer depends. SQL-widget target/binding is explicitly open in the latest dashboards Feature, and must remain so rather than being invented by an agent.

The next Dashboardius milestone should be one real query → one card → save → reload → authorized colleague, with a real cross-language persistence test. Future GitHub storage should reuse DataTug's project storage mechanism. Dashboardius needs an owning Feature specification for this journey; a backstage idea and marketing README are not sufficient acceptance contracts.

## Implementation reality check

The following checks verify planning claims; they are not an exhaustive security or code audit.

### Confirmed integration mismatches

| Boundary | Browser implementation | Server implementation / implication |
|---|---|---|
| Semantic columns and related lists | `libs/datatug/semantic/src/lib/models/models.ts` and `services/semantic-api.service.ts` consume arrays directly. | CLI `pkg/server/endpoints/semantic_columns.go` and related endpoint return envelopes such as `{columns: [...]}` and `{related: [...]}`. |
| Applicable saved queries | Browser posts project in body and expects `queryId` / `parameterId` fields. | `semantic_applicable.go` reads project from URL query; response uses `query` / binding `parameter`, with different provenance requirements. A real browser-shaped request returned HTTP 400, project required. |
| Related rows | Browser expects a nested recordset with column objects and row arrays. | Server returns column names and `{key,data}` rows. Limitations also differ. |
| Saved-query execution | Query page sends project, query ID and parameters. | `pkg/api/run_query_api.go` requires environment and database before dispatch. Real request returned HTTP 400 for missing environment. Response result and binding shapes differ too. |
| SQL parameter application | UI can show supplied bindings. | Server passes SQL text to `RunNativeSQL` without those parameters, then echoes request parameters as `BindingsApplied`. That echo is not proof they were applied. |
| HTTP query execution | Part of approved browser journey and task completion narrative. | Server code explicitly refuses HTTP saved queries; the CLI-specific HTTP path exists. |
| Semantic source discovery | Intended to work on the same sources as queries. | `semantic_schema.go` resolves its own project filesystem layout; normal execution resolves environment/catalog configuration. Against the existing demo, a semantic-columns request returned unknown source. |
| Context and agent selection | Root context service persists one session key; agent URL factory snapshots the initial route. | No project/environment/principal namespace in that context key, and store switches are explicitly deferred in the factory comment. This violates the required investigation boundary unless navigation is constrained and tested. |

These are contract failures visible in source and, for the indicated requests, in actual transport responses. They explain why individual package tests can pass while the product journey remains unproved. See [CLI execution source at reviewed main](https://github.com/datatug/datatug-cli/blob/38b677b28c5a0e7083c9716524d7ffeafd04a087/pkg/api/run_query_api.go) and [browser semantic models at reviewed main](https://github.com/datatug/datatug-apps/blob/7775f53211ba30430befe032ff86d10d44850e84/libs/datatug/semantic/src/lib/models/models.ts).

### Checks actually performed

- CLI `38b677b`: focused `go test ./pkg/secureread ./pkg/server/... ./pkg/api` passed. A CLI binary built successfully. These checks do not establish browser interoperability.
- Apps `7775f53`: resolved Nx configuration inspected and development build/serve succeeded. This review did not independently run a full production build or the complete browser suite.
- The built CLI ran on a temporary loopback port against the existing demo project. Ping and project summary returned HTTP 200; the semantic/applicable/run-query checks above returned real errors. No fixture was rewritten to hide them.
- Current browser journey source leaves J4 as `test.fixme` and skips on missing CLI/demo prerequisites. J1 jumps directly to a table URL rather than proving all earlier navigation. Therefore “all journeys passed” is not justified by this suite's mere completion.
- Local current DataTug opened and displayed the project shell, including Investigation Context navigation. The older audit's implication that this route is wholly absent is stale. Production DataTug still exposes legacy homepage copy; clicking the demo entry did not produce useful navigation in this inspection. The cause of that production behavior was not established.
- Dashboardius live scripted action, Undo and Save messaging were inspected as described above. No claim is made that persistence, sign-in or real query execution works there.

## Mac versus VM; merged versus pending reconciliation

### Verified remote main snapshot

| Repository | Remote main SHA |
|---|---|
| datatug/datatug | `054a7f47903851fad8b70c13d8b8e663837b6b50` |
| datatug/backstage (formerly datatug/marketing) | `160afec6a08229fe83f3dd22040aa90a165d6e79` |
| datatug/datatug-cli | `38b677b28c5a0e7083c9716524d7ffeafd04a087` |
| datatug/datatug-core | `33a2eeb2730a6bf5d1cf56b0e7ae0e8752bf3e20` |
| datatug/datatug-apps | `7775f53211ba30430befe032ff86d10d44850e84` |
| datatug/datatug-demo-projects | `e36b9ddddbdaf2a462f09b97cdd7c18b02b9f45c` |
| datatug/datatug-site | `435ac7d5320e347ca946c1602599d110b6109113` |
| sneat-co/dashboardius | `c6d8f9d7ece1e6fc7e6279a0ec842e07cd6bd9b6` |

Mac hub main was `6ff6d15`, while its `spec-discoverability` worktree contained current remote main. VM hub `e0bfe3a`, marketing `cb80bae`, CLI main `9cef1e2`, apps main `9c2c2d8`, demo `77dabec`, core `b49b211`, and site `86b5af7` were stale relative to relevant newer work. No checkout was pulled or switched during review.

### Significant pending work and recommended disposition

VM relative paths below are under `/home/ai/projects` unless stated otherwise. “Rebase” and “integrate” are recommendations for a later authorized session, not work performed here.

| Work | Observed state | Recommendation |
|---|---|---|
| CLI `.worktrees/phase1-cli-followups-2` | HEAD `38b677b`, matching remote main, with nine staged file modifications covering saved-query/native SQL bindings, dependency changes and HTTP/auth/CORS/security tests. | **Keep and finish.** Review staged contents independently of HEAD. Carry actual SQL argument binding into the server execution path too; do not treat CLI-only binding as endpoint completion. |
| CLI `.worktrees/phase1-cli-semantic-endpoints` | HEAD `e8549ef`; committed endpoint work incorporated, but an unstaged `semantic_demo_project_test.go` delta remains. | **Preserve the test delta**, compare it with current acceptance, integrate useful coverage. Do not delete as an “already merged worktree.” |
| CLI `.worktrees/phase1-core-module-dependency` | HEAD `2674bd7`, no unique commits against inspected VM main. | **Likely completed residue**, subject to remote/patch and dirty-state verification before later cleanup. |
| CLI `.worktrees/datatug-layered-acl-query` | HEAD `a7e1f62`, including protected proxy work and configuration changes; references the old vendored core configuration path. | **Keep, rebase and partially integrate.** Move shared configuration changes to the external core module and reconcile overlapping server/auth/origin code. Preserve the stronger fail-closed profile. |
| Older VM WB `datatug-query-policies` worktree | HEAD `f6d6d17`; ancestry differences do not establish unique functionality after squash merges. | **Compare patches**, extract any missing tests; do not merge wholesale or discard by commit count. |
| Apps `.worktrees/datatug-apps-layered-acl-query` | HEAD `5ebb264`, four commits adding protected query UI, real acceptance, context binding and constrained-result validation. | **Keep and integrate selectively.** `/openvaultdb` proves a protected path, not the semantic core journey. Reuse its security behavior without creating a permanent parallel product. |
| Apps integration worktree | HEAD `e45adcd`, integrates layered UI onto an earlier main. | **Rebase current candidate**, reconcile later main changes, validate once after integration. Existing integration success does not cover today's newest main. |
| Apps SSO branches/worktrees | Enterprise SSO branch `38155e6`, another integration `8bcfb1d`; local modified indexes and untracked specification files remain. | **Preserve and compare** against now-canonical SSO specs. Keep outside the core-loop candidate unless required; no blanket cleanup. |
| Site `.worktrees/phase1-access-policies-native-sql` | HEAD `8e38c69`, pending opaque-query/policy documentation. | **Keep and revise** to match explicit protected versus privileged execution profiles. |
| Apps PR #48 | `ci/phase1-journey-ci`, head `b6f1685`, OPEN; journey-e2e FAILURE, other build/test/default E2E checks SUCCESS in run `34337467977`. | **Finish/rebase.** Pin companion revisions, remove successful skip paths for required CI, implement J4. Do not call it green because other jobs pass. |
| Apps PR #54 | `ci/phase0-app-deploy`, head `323d643`, OPEN; observed build/test/E2E checks SUCCESS. | **Keep, reconcile deployment ownership and gating.** Proposed main deployment does not depend on the new journey gate; possible overlap with existing Cloudflare Git integration needs configuration verification. |
| DALgo layered ACL | Worktree `dal-go/dalgo/.worktrees/layered-acl-query`, head `c8adfb6781756938250573e7da531cee02ce9418`; PR #158 OPEN. Relevant lint/build/test/schema/spec conclusions SUCCESS, run `34360498361`. | **Preserve provider-first delivery.** Green PR checks are not a merged/released dependency. Delivery notes record a branch-protection/WB blocker; resolve through its owner, not by assuming release. |
| ACL provider/SDK chain | SQL adapter `4bb140a70937ad35d2d4a328e19459b7cb9d7491`; inGitDB adapter under `ingitdb/dalgo2ingitdb`, `3b86e58341b55212af8aefa109af59894e45c5b8`; OpenVaultDB SDK `89f704feb30a528e2b85be000636d6e35c45757f`. | **Keep together as a dependency wave.** Inspect/release providers before converging consumers to published versions. Local GOWORK success is not published dependency convergence. |
| Already integrated adjacent ACL work | DTQL spec packet `c2d048784ef982a0533b1cdfb912f13182ba6b53`; dalgo2openvaultdb mapping fix `cacd422564fc27844b989e1c80eb96320ce4af4b`. | **Preserve; do not redo.** Neither implies the complete layered stack is landed. Distinguish the DTQL specification repository from DALgo's executable DTQL implementation. |
| Dashboardius prototype | Clean inspected repository at `c6d8f9d`; no today's implementation delta found. | **Preserve reducers/widgets/UX.** Replace duplicate durable model definitions through a real contract migration, not a rewrite of the useful presentation layer. |

The VM delivery packet reports protected DataTug browser acceptance 2/2, 26 unit tests, build and lint passing with locally linked dependencies. This review inspected the implementation and record but did not independently rerun those VM tests. Preserve that evidence and rerun the relevant acceptance against the final integrated, published dependency candidate. Other notarization, CLI-framework and CI-hygiene residue was inventoried but is not evidence of unfinished core-loop implementation; do not absorb unrelated work into this remediation.

## Ranked problems

| Rank | Problem | Consequence and required correction |
|---|---|---|
| **Critical** | Contradictory protected-read/native-SQL contract and incomplete support policies | Freeze trust profiles, refuse unconstrainable protected execution, and specify policy coverage for every journey source. No current support native-SQL exploit is asserted. |
| **High** | Browser/server wire and source-resolution contracts disagree | Core loop cannot be accepted as an integrated product. Publish and test one exact contract. |
| **High** | Bindings are echoed without being applied on the saved SQL server path; HTTP browser dispatch absent | Apparent execution success can misrepresent the query actually run. Verify actual parameter effect and complete all three approved sources. |
| **High** | Context lacks project/environment/principal isolation and deterministic ambiguity behavior | Wrong facts can bind the next query. Define scoped typed context and test transitions. |
| **High** | Phase 1 completion outruns actual E2E evidence | J4 disabled, optional prerequisite skips and deep-link shortcuts leave required acceptance unproved. Correct task evidence and CI gates. |
| **High** | Pending stronger ACL implementation not incorporated into planning baseline | Risk of duplicate work or weaker replacement. Reconcile provider chain and consumers before the next wave. |
| **High** | Dashboardius durable schema and execution binding incomplete | Demo interactions cannot yet round-trip as the promised shared board. Complete canonical contract and real persistence journey. |
| **Medium** | Competitive synthesis overstates white space and certainty | Strategic confidence exceeds research evidence. Calibrate scoring, add targeted peers, run a demand test. |
| **Medium** | Later roadmap phases bundle unrelated product decisions | Low-autonomy, high-conflict implementation waves. Split collaboration, hosting, playbooks and dashboards. |
| **Medium** | Scope rulings and inferred permanent prohibitions intermixed | Separate product ownership from embedding, and initial playbook scope from all future automation. |
| **Medium** | Old ideas, missing local Feature links, ownership and status sketches conflict | New contributors cannot determine authority or current state. Reconcile canonical indexes and dated evidence. |
| **Low** | Public prototype polish and stale homepage copy | Invalid license notice and legacy copy weaken credibility. Fix later without distracting from the core contract/acceptance work. |

## Recommended documentation/spec corrections

Update existing canonical artefacts rather than generating a competing spec set:

1. Core-investigation-loop Feature: exact wire schemas, target/source identity, applied binding behavior, errors, limitation states, principal profiles, HTTP boundary and all four journeys.
2. Investigation Context Feature: typed scoped state, conflict/missing-input behavior, provenance, navigation isolation and safe sharing.
3. Sharing-scopes Feature and roadmap: standalone project with grants to multiple spaces; explain the platform exception consistently.
4. Phase 1 plan: reconcile delivered components versus accepted journeys, fix dependencies, fixture assumptions/counts and required CI behavior. Link each completion to the candidate and observable acceptance.
5. Dashboards Feature: retain the already-recorded parameters decision; keep the SQL-widget binding decision open with concrete alternatives. Specify cross-language row/widget compatibility and separate Phase 1 exclusion from permanent embedding policy.
6. Research methodology and synthesis: reproducible weights/Unknown handling, intended-versus-current comparison, overlapping product normalization and evidence-qualified claims.
7. Product hierarchy and local spec indexes: classify outcomes/enablers/surfaces, reconcile old ideas, and repair missing implementation Feature links.

Product questions should remain in the owning Feature's Open Questions. This review identifies them; it does not silently answer the explicitly pending SQL-widget decision or route it into an unrelated project.

## Recommended implementation/reconciliation work

Preserve the shared semantic engine, DALgo HTTP adapter and example descriptors, protected-read tests, external core module, current board-parameter fix, context/grid/query UI components, Dashboardius reducer/widgets, and VM ACL chain. The central remediation is integration, not replacement.

Use a single reviewed dependency/candidate manifest, one landing owner per repository/target, and isolated WB worktrees for edits in the later implementation session. Review the staged VM follow-ups and old core import path before cherry-picking or merging. Resolve the canonical contract once, adapt both sides, and run a real consumer/provider test before broader UI acceptance. No source should be cleaned merely because its branch name appears in an earlier merge.

Landability must be established from current remote refs and authoritative run conclusions, including the exact integrated candidate. This review deliberately leaves all source branches, staged changes and worktrees untouched.

## Exact recommended next autonomous coding task

The next task should be **“Converge and prove the DataTug Phase 1 core investigation journey.”** It includes bounded canonical spec repair before code; it does not reopen product research or begin Phase 2. Give Astra this prompt:

> Read `datatug/datatug/spec/research/2026-09-09-independent-review.md`, the original `datatug-frontier-research-planning-prompt.md`, and the current canonical core-investigation-loop Feature and Phase 1 plan. Complete the existing Phase 1 vertical browser slice by reconciling merged and pending work; do not restart the implementation or repeat competitor research.
>
> First refresh the read-only inventory against current remote refs on Mac and VM. Preserve staged/unstaged work. Inspect the pending CLI follow-ups, layered ACL provider/SDK/server/UI chain, and journey/deployment PRs identified in the review. Name one landing owner per repository and use WB-managed isolated worktrees for changes.
>
> Before code, correct the owning canonical Features and plan in place: publish one exact browser/server contract for semantic columns, related data, applicable queries and execution; unify project source/target resolution; define typed context scoped to store/project/environment/principal; define missing/ambiguous parameter behavior; and make protected execution fail closed when row/column policies cannot be enforced. Treat an explicit privileged opaque-query grant as a separate trust profile. State implementation assumptions as assumptions, not founder rulings. Preserve SQLite, inGitDB and HTTP in approved scope.
>
> Reuse the existing semantic engine, HTTP adapter, browser components and stronger pending ACL work. Reconcile the external core-module migration and publish/converge the provider chain in dependency order where required. Make saved-query parameters actually affect server execution, support HTTP saved queries through the browser, and align real serialized request/result/limitation shapes. Do not substitute fake responses or local fixture patches for integration.
>
> Verify the complete browser journey from its real entry point: project → source → useful result → semantic selection → related data/applicable query → visible correct bindings → authorized server execution → retained context → continued investigation. Verify restricted access on every participating source, including denial of forbidden native SQL, hidden fields, filtered rows, and policy-safe related lookups/counts. Verify missing parameters, ambiguous values, unavailable HTTP, explicit snapshot selection, and isolation when project/environment/principal changes. Fix the support demo policies and commands so roles and expected rows are explicit.
>
> Required CI must fail when prerequisites are missing, must execute J4, and must record/pin companion revisions. Add real transport and parameter-effect checks; then run meaningful end-to-end acceptance without request interception or deep-link shortcuts. Reconcile the deploy path so the same candidate's required journey gates protect publication. Do not add dashboards, a UI-kit migration, hosted tenancy, or Phase 2 features to this task.
>
> If an actual unresolved founder decision blocks the core slice, record a concrete choice in the owning Feature and ask only for that decision while continuing independent work. The open Dashboardius SQL-widget binding decision is outside this task. Follow current authorization and repository protections for landing; do not infer new permission from this review. Report exact remote SHAs, authoritative CI outcomes, full journey evidence, remaining blockers and precise cleanup state. Do not call a component-level pass or a locally linked dependency wave a completed product journey.

## Final verdict

**REVISE PHASE 1 BEFORE IMPLEMENTATION.** Keep Fable's product thesis, most of the hierarchy and research evidence, and the substantial implementation. Correct the contracts and acceptance boundary, then converge the existing code into one demonstrably working investigation. The strategy is promising; the current plan is not yet safe to hand to an autonomous implementation wave unchanged.
