---
format: https://specscore.md/plan-specification
status: Executing
---
# Plan: Phase 1 — core investigation loop

**Status:** Executing
**Source Feature:** core-investigation-loop
**Date:** 2026-09-09
**Owner:** alex
**Supersedes:** —

## Summary

Complete the approved browser investigation over SQLite, inGitDB and HTTP by
converging existing implementations. Do not restart completed components or begin
Phase 2 before acceptance. The [Feature](../features/core-investigation-loop/README.md)
and its [normative transport appendix](../features/core-investigation-loop/api-contract.md)
now specify the missing boundaries.

The [independent review](../research/2026-09-09-independent-review.md) records the
9 September remote refs, API failures, disabled J4 and pending VM work. Those are a
dated baseline; refresh refs and authoritative CI conclusions before editing code.
The auth-hook panic cited in earlier progress notes was no longer reproduced by the
review's current CLI: project summary returned 200. Remaining contract failures must
not be attributed to that historical cause.

Completed task records below retain their original lifecycle/provenance. They prove
landed components, not acceptance of the revised Feature. SpecScore complete is
terminal; corrective tasks 12–17 cover the additional work rather than rewriting
history. No completed component or review finding is a waiver of a journey AC.

## Approach

### Journey first

Walk J1–J4 and both epilogues in the Feature. The null-action start must reach real
rows from the real entry point. The middle must preserve explicit typed facts and
bindings across all three sources. The restricted investigator must see only their
authorized result. Closing/reopening and sharing/replay each have their own terminal
outcome; neither promises durable unsaved context or portable authorization.

### Reconcile before implementing

Inventory current remote main/PR refs, staged/unstaged deltas and worktree ownership
on Mac and VM. Preserve CLI follow-ups, semantic demo tests, the external core module,
layered ACL provider/SDK/server/UI work, and journey/deploy PRs. Use the review's
reconciliation table as leads, not fresh receipts. No wholesale merge of stale
integration branches or duplicate protected-query product surface.

### Dependency order

Task 12 freezes schemas and source resolution, using landed task 5 models. Then
13 integrates executor/policy dependencies and 14 completes HTTP; 15 adapts the
browser/context. Task 16 corrects fixtures and performs real consumer/provider checks.
Task 17 proves the whole journey and publication candidate. Tasks 13 and 15 may run in
parallel; task 14 follows 13. Parallel work is allowed only in separate worktrees with one declared contract and landing owner.
Their final integration is one candidate per repository, not independent declarations
of a finished journey.

The old model dependency is now external. Never revive the deleted vendored core
copy or mirror fixes into two authorities. Port pending configuration changes upstream
and converge consumers to published provider revisions before release acceptance.

### Execution tracking and claims

Task status below is the shared progress record, maintained with `specscore task
change-status --plan 2026-09-09-phase-1-core-investigation-loop`. It is not updated
automatically from agent activity or GitHub merges. Completed tasks must carry
implementation evidence; an active task's Note identifies its owner, bounded scope
and remaining dependencies. Publish status changes so another machine can read them.

Before starting, inspect current Mac/VM worktrees and open PRs, then acquire the shared
WB claim `datatug-phase1-task<N>` without force or takeover. Different historical
worktree names can cover the same task, so an absent claim alone does not prove that
the work is free. Preserve existing owners and their pending changes. WB Work Logs
record worktree/session ownership; the plan Note connects that ownership to this
plan's task ID. Release the claim only after the lane's work is handed off or closed.

Independent preparation may start while dependencies remain open, but it must name
that boundary in its Note. It does not close the task's dependent acceptance. When
only predecessor work remains, record `blocked` with the precise missing evidence;
resume through the lifecycle command after the dependency is available. Never count
a policy fixture or isolated test as the complete browser/agent journey.

### Definition of done

Every revised Feature AC passes against recorded immutable client/server/core/provider/
demo revisions. Required CI fails on missing prerequisites, runs J4 and starts from
the actual entry point without intercepted responses. Distinguish local tests,
integrated candidate CI, remote landing, CLI release, hosted deploy and human smoke.
Record the exact CLI release owner/tag and deployed app SHA; do not invent a release
version as part of this plan. A locally linked GOWORK test is useful evidence but not
published dependency convergence. One deployment owner/path gates the same candidate
on required journey checks.

Follow current user authorization and repository protections for landing. No unrelated
branch/worktree cleanup is part of this plan. Finished source cleanup follows verified
remote receipts, never an assumption based on branch names or task status.

## Tasks

Tasks 1–11 preserve the original work-package IDs and component descriptions for
traceability. If their older implementation detail conflicts with the revised Feature,
the Feature and its transport appendix govern; tasks 12–17 perform the reconciliation.
No task may remove an approved source, skip a required journey, or weaken protected
execution to satisfy a numeric task/test budget.

### Task 1: Build, serve and store wiring in datatug-cli

**Id:** task-1
**Verifies:** core-investigation-loop#ac:agent-path-contract-server
**Depends-On:** —
**Status:** complete
**Implemented-by:** datatug/datatug-cli@940902b65a4c965c3b524d9be0ebf8e2b9125d2c (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

Fix the go-github v91 imports so `main` builds; add a CI gate that runs `go build ./...`
on dependency PRs. Make `datatug serve` start: nil-check `config.Server`, bind
`--host/--port/--project/--as/--role/--group` flags, pass a real `Handler` to
`endpoints.RegisterDatatugHandlers`, and implement `storage.NewDatatugStore`/`GetStore`
over `filestore` (mirroring `project_base_command.go`). Restore the module dependency on
`datatug-core` (delete the vendored `pkg/datatug-core` copy) without keeping a second vendored copy.
Replace the `recover()`-around-the-panic test with a real handler test. Tests: unit
tests for store wiring and flag binding; `go test ./...` green. Non-goals: ACL, new
endpoints, TUI.

### Task 2: Web agent client contract and demo entry

**Id:** task-2
**Verifies:** core-investigation-loop#ac:agent-path-contract-client, core-investigation-loop#ac:first-result-two-clicks
**Depends-On:** 1
**Status:** complete
**Implemented-by:** datatug/datatug-apps@8653bd9 (main)
**Note:** Complete 2026-09-11 (lead session, historical task closed on final evidence). The web client speaks the frozen transport (datatug-core v0.27.3 apicontract envelopes, golden fixtures with manifest drift test, datatug-apps#73 and later) and the demo entry is the project URL: the required journey CI job on datatug-apps main (run 34477823015, datatug-cli v0.23.0) walks J1 from the project page to a first result by clicking through project → environments → catalog → Album (agent-path-contract-client, first-result-two-clicks).
**Evidence:** datatug/datatug-apps#91@8653bd9, datatug/datatug-apps#83@3c13cb0, datatug/datatug-apps#73@1bb8406, datatug/datatug-cli@v0.23.0

In `datatug-apps`: prefix every agent call with `/datatug` and use the registered
endpoint names; accept `localhost:8989`-style store ids; provide `DatatugFoldersService`
and `DatatugNavService` on the routes that lack them; fix the demo project link to
`datatug-demo-projects/main/demo-project-1/datatug-project.json`; hide the six empty
shell pages behind a feature flag; make the `EnvDbPageComponent` route resolve. Tests: a
contract test that starts the real `datatug serve` binary and asserts each client call
resolves (no `localhost:8989` interception). Non-goals: new screens.

### Task 3: Journey e2e harness against the real agent

**Id:** task-3
**Verifies:** core-investigation-loop#ac:journey-harness
**Depends-On:** 1, 4
**Status:** in_progress
**Note:** Open defect, 2026-09-10: the journey title-load assertion, the Title textbox value that J1 and Epilogue A both check, intermittently stays "Loading..." until its 15 s timeout. It failed main's journey CI job, single worker, in Epilogue A in datatug-apps runs 34453639702 (4d242ad) and 34457418273 attempt 1 (03b6f16, which passed on re-run), and J1 once locally under fullyParallel against datatug-cli main 5be73b7. In that local failure the Playwright trace shows both GET /datatug/projects/project_summary calls answered 200 in about 5 ms and 1 ms with the correct CORS header and no console error: the agent answers promptly and the page never commits the title, so the fault is client-side in datatug-apps. Host load does not explain it; root cause not established. Until it is fixed, the J1 walk this task delivers is not reliably green. This Note is the single account of the defect: datatug-apps journey.spec.ts's header and journey/README.md point here, so record each finding here.

Playwright fixture that builds/launches `datatug serve --project <demo> --as admin --role admin` on
a free port, starts the web dev server pointed at it, and tears both down; a first test
walks J1 (open demo → Album → rows) asserting the agent log line. Non-goals: J2–J4
(task 10).

**Annotation Amendment:** actor=claude-code (Claude Opus 5); at=2026-09-10T10:46:14Z; reason=Record the open journey title-load flake and its evidence as its single account; datatug-apps docs point here; before_sha256=58d044f81d34d92db015aaffff122ddaf91c855d3ef4bb4c3505ef9afee9a2de
### Task 4: Demo knowledge project and credential guard

**Id:** task-4
**Verifies:** core-investigation-loop#ac:demo-project-shape
**Depends-On:** 5
**Status:** complete
**Implemented-by:** datatug/datatug-demo-projects@6ca9ddd1452ab1179a6c1e019b400b18bcf4a844 (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

In `datatug-demo-projects/demo-project-1`: add `Customer` and `Invoice` entities with
declared mappings (`Customer.ID → chinook.Customer.CustomerId`, `Customer.Email →
chinook.Customer.Email`, `Invoice.ID → chinook.Invoice.InvoiceId`, plus
`Customer.ID → support-notes.CustomerId`), semantic types, and `NamePatterns` for
`CustomerId`/`customer_id`; three parameterized queries (use DTQL for every query in the protected support journey; retain SQL examples only under the explicit privileged profile) with `Parameters[].Meta`
set; an inGitDB `support-notes` collection (8 records keyed by `CustomerId`, 2 for
customer 5); `policies/` with explicit admin/support roles and separate Customer, Invoice and support-notes policies (Canadian scope, Customer.Email hidden); a `Country` entity with `Country.Name → chinook.Customer.Country` and
`Country.Currency`; an HTTP `QueryDef` for the demo-selected country provider (`name` parameter tagged
`Country.Name`) and one for Frankfurter (`to` parameter tagged `Country.Currency`), each
with a recorded snapshot fixture under `fixtures/http/`; README updated with the real
datatug.app link. In `datatug-core` +
`datatug-project-validator`: a validation rule that fails when any `QueryDefTarget`
carries credentials. Tests: `datatug validate` green; a fixture test loads the inGitDB
collection through `ingitdb://`. Uses the model from task 5 (coordinate branches).

### Task 5: Field mappings model and semantic resolver

**Id:** task-5
**Verifies:** core-investigation-loop#ac:mapped-columns-server, core-investigation-loop#ac:inferred-mapping-labelled
**Depends-On:** —
**Status:** complete
**Implemented-by:** datatug/datatug-cli@5f8881e5e0a9b82716bb46156b68f542da0e91f0 (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

In `datatug-core/pkg/datatug`: `EntityField.Mappings []PhysicalRef{Source, Collection,
Column}` with validation and filestore round-trip; `QueryDef.Type` gains `DTQL`. In
`datatug-cli`: a `pkg/semantic` resolver that, for a source/collection, returns per
column `{entity, field, provenance: declared|inferred}` using declared mappings first
and `NamePatterns` second; `GET /datatug/semantic/columns`. Tests: table-driven resolver
tests (declared wins over inferred; no match → nothing; case-insensitive patterns);
endpoint test with the demo project. Non-goals: suggestions beyond name patterns.

### Task 6: Access control on every server read, principal flags, DTQL execution

**Id:** task-6
**Verifies:** core-investigation-loop#ac:restricted-rows-and-columns-server, core-investigation-loop#ac:hidden-column-refused, core-investigation-loop#ac:native-sql-labelled, core-investigation-loop#ac:dtql-query-runs
**Depends-On:** 1, 5
**Status:** complete
**Implemented-by:** datatug/datatug-cli@a30132d9fee1f8f670b06bb61f49c84fa119c90e (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

Route `exec/select`, saved-query runs and the new `POST /datatug/exec/run_query`
through `accesspolicies.Run` with the project policy set (`<project>/policies/` then
`~/.datatug/policies`) and the principal fixed by `serve --as/--role/--group`; open
sources through `pkg/dbcopy/url.go` (`sqlite://` via `dalgo2sqlite`, `ingitdb://`);
execute `DTQL` queries via `dtql.Deserialize`; execute `SQL` text in a read-only
session only with an explicit privileged opaque grant; refuse it in the protected profile; return structured `limitations`
(policy, rowsFiltered, hiddenColumns) and `ACCESS_DENIED` refusals that never echo
hidden values; `agent-info` reports the principal. Tests: the security matrix from the
Feature (unauthorized source; row restriction survives an alternate query; hidden column
absent from select/related/saved/DTQL; hand-crafted hidden-column request refused; no
principal accepted from the request); `serve` refuses to start with a policy set and no
`--as`. Non-goals: multi-user auth, layered policies, Explain Access.

### Task 7: Related lookups and applicable queries endpoints

**Id:** task-7
**Verifies:** core-investigation-loop#ac:related-across-sources-server, core-investigation-loop#ac:applicable-with-chain-server, core-investigation-loop#ac:related-respects-policy
**Depends-On:** 5, 6
**Status:** complete
**Implemented-by:** datatug/datatug-cli@5f8881e5e0a9b82716bb46156b68f542da0e91f0 (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

`GET /datatug/semantic/related` computes lookups from the FK graph (`ForeignKeys`,
`ReferencedBy`) and from mappings of the same field in other sources; `GET
/datatug/semantic/related/rows` builds a DALgo `StructuredQuery` and runs it through
task 6's executor (counts come from the policy path too, else "count unavailable");
`POST /datatug/queries/applicable` matches `ParameterDef.Meta` against the supplied
values and returns applicable/not-yet lists with resolution chains. Tests: demo project
fixtures (customer 5 → Invoices 7, Support notes 2; Invoice lines not yet applicable);
restricted principal sees zero rows and no true count.

### Task 8: Web semantic markers and context panel

**Id:** task-8
**Verifies:** core-investigation-loop#ac:mapped-columns-marked, core-investigation-loop#ac:related-across-sources, core-investigation-loop#ac:applicable-with-chain
**Depends-On:** 2, 7
**Status:** complete
**Implemented-by:** datatug/datatug-apps@8653bd9 (main)
**Note:** Complete 2026-09-11 (lead session, historical task closed on final evidence). Semantic column markers, the context panel with related lookups across sources and applicable queries with binding chains are live in datatug-apps (#73, #79, #85) and proven by J2 in the required journey CI job on main (run 34477823015): clicking a CustomerId cell shows related records (Customer/Invoice/support-notes) and applicable queries, and running one binds the parameter (mapped-columns-marked, related-across-sources, applicable-with-chain).
**Evidence:** datatug/datatug-apps#91@8653bd9, datatug/datatug-apps#83@3c13cb0, datatug/datatug-apps#73@1bb8406, datatug/datatug-cli@v0.23.0, datatug/datatug-apps#79@f5873b2, datatug/datatug-apps#85@fc84ada

Grid column headers consume `semantic/columns` (marker + tooltip with provenance); a
cell click opens the context panel (Ionic side panel in `ion-split-pane`; grid remains
the `@sneat/datagrid` island) showing the semantic value, related lookups with counts and
expandable rows, applicable queries with the chain, and "Add to context". Replace the
hard-coded `CellPopoverComponent` sample content. Tests: component tests with fixture
responses; keyboard reachability. Non-goals: hover mechanics beyond the tooltip.

### Task 9: Query page auto-binding, Investigation Context and limitation header

**Id:** task-9
**Verifies:** core-investigation-loop#ac:bound-from-selection, core-investigation-loop#ac:context-carries, core-investigation-loop#ac:restricted-rows-and-columns
**Depends-On:** 2, 6
**Status:** complete
**Implemented-by:** datatug/datatug-apps@8653bd9 (main)
**Note:** Complete 2026-09-11 (lead session, historical task closed on final evidence). Query page auto-binding from selection, the Investigation Context basket with from-context binding and chip clearing, and the limitation header are live (datatug-apps#73, #81, #82) and proven by J1/J3/J4 in the required journey CI job on main (run 34477823015): Customer.ID binds 'from selection' and runs as a plain parameter, a context value binds 'from context' on another query and disabling the chip empties it, and the support principal sees Canada-only rows with Email hidden and the limitation stated (bound-from-selection, context-carries, restricted-rows-and-columns).
**Evidence:** datatug/datatug-apps#91@8653bd9, datatug/datatug-apps#83@3c13cb0, datatug/datatug-apps#73@1bb8406, datatug/datatug-cli@v0.23.0, datatug/datatug-apps#81@814c104, datatug/datatug-apps#82@438f1dc

An `InvestigationContextService` (signals, session-scoped) holding
`{entity, field, value, origin, enabled}` items and the environment; the "Variables"
route becomes the Investigation Context screen; every project page shows the collapsed
context with a count; the query page binds `ParameterDef.Meta` parameters from the
selection or context (chip with origin, clear/override), and the results header lists
bindings applied and policy limitations (from task 6). Run uses `POST
/datatug/exec/run_query`. Tests: binding precedence (selection over context), disabled
chip empties the value, header renders limitations. Non-goals: persistence of context
to the project.

### Task 10: Journey e2e J1–J4 integration

**Id:** task-10
**Verifies:** core-investigation-loop#ac:journey-walks
**Depends-On:** 3, 8, 9, 11, 12, 13, 14, 15, 16
**Status:** complete
**Implemented-by:** datatug/datatug-apps@438f1dc4 (main)
**Evidence:** datatug/datatug-apps#82@438f1dc, datatug/datatug-cli#221@aba31f6, datatug/datatug-apps#80@5bee04e, datatug/datatug-apps#81@814c104, datatug/datatug-cli#219@743b8d9, datatug/datatug-apps#79@f5873b2
**Note:** Complete 2026-09-10: J1, J2, J3 and J4 pass single-worker against datatug-cli built from main aba31f6 (PR #221) and datatug-apps main 438f1dc (PR #82) — lanes S103/S104. J4 runs against its own --as support --role support agent (no intercepted responses, no deep-link shortcuts): only Canadian rows, Email hidden, limitation header on the table and query pages, Brazilian id → no rows, Canadian id → only its 7 invoices, hidden-column selection refused. Journey harness: apps/datatug-app/e2e/journey (README documents DATATUG_BIN/DATATUG_CLI_DIR/DATATUG_DEMO_DIR). Follow-up: the journey project is not yet a CI job (PR #48 being revived, lane S105); task 17 owns the release candidate and publication. Update 2026-09-10: that follow-up is done. datatug-apps PR #83 (merged 2026-09-09) added the "journey e2e (real datatug-cli agent)" CI job, pinned to datatug-cli v0.22.0 on main as of 2026-09-10; PR #48 was closed unmerged.

Complete J1–J4 in the real-agent harness after corrective transport, executor,
HTTP, context and fixture integration. Cover the revised Feature mechanisms without
interception, missing-prerequisite skips or direct-entry shortcuts. This task produces
journey evidence only; it does not release or deploy. Task 17 owns final candidate,
epilogue, publication and remote receipt checks.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T20:16:49Z; reason=Record the first J1 end-to-end pass; before_sha256=de609446bc90d9cd53833aedc398394f8c3f56608698f28a2e94af2f6bd5bc8e
**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T21:24:53Z; reason=Record the J2 end-to-end pass and the J3 state; before_sha256=f16e13bd37c65c99c84fd98ef99bb1b4eedce6c9cdd1389ab7fe3706d18c6e84
**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T22:18:55Z; reason=Record the J3 pass, the canonical query-id convention and the J4 harness with its first failure; before_sha256=a1148755786cd6a853fca8379ca5f2babc4d25cf33059dd9e42b354abc6ffde7
**Annotation Amendment:** actor=claude-code (Claude Opus 5); at=2026-09-10T10:46:14Z; reason=Mark the journey-CI follow-up done: datatug-apps PR #83 added the journey CI job; before_sha256=c6e9f74183f94c634d28fa4079b9c4cda68f7d9327caaf6c153891e1ccfa5c24
### Task 11: HTTP reference-data source executor

**Id:** task-11
**Verifies:** core-investigation-loop#ac:http-source-in-demo
**Depends-On:** 4, 6
**Status:** complete
**Implemented-by:** datatug/datatug-cli@2dd8afee2e44b66df0a49dedfbb4239856965357 (main)

**Acceptance note:** Historical component completion; revised journey/contract acceptance is tracked in tasks 12–17.

Founder ruling 2026-09-09 (verbatim): *"I want to have in MVP: inGitDB, SQLite, HTTP
endpoint."* Executor: the generic DALgo adapter **`dal-go/dalgo2http`** (founder 2026-09-09:
*"dalgo2http - yes"*): declarative per-collection descriptors, GET only, equality pushdown on
declared parameters, fail-closed on anything else, recorded snapshots with `live`/`snapshot`
provenance. In `datatug-cli`: register `http://`/`https://` in `pkg/dbcopy/url.go` to open a
`dalgo2http` database from the project's HTTP `QueryDef`s, and surface the adapter's
provenance in the result header. Related lookups (task 7) and applicable queries treat
the HTTP source like any other. Tests: unit tests with a local `httptest` server and the
snapshot path; the demo's Country cell → Country facts and the exchange-rate query pass
in explicit snapshot-mode negative/offline checks; a separately identified live check proves HTTP acceptance. Non-goals: POST, auth headers, GraphQL,
HTTP steps in playbooks (Phase 2–3). Wave B, runs beside task 7.

### Task 12: Freeze transport and unify source resolution

**Id:** task-12
**Verifies:** core-investigation-loop#ac:real-transport-and-parameter-effect
**Depends-On:** 5
**Status:** complete
**Note:** Landed on all three sides: server (datatug-cli PR #212, exact envelopes for agent-info/exec/run_query/semantic/columns/semantic/related(GET->POST)/semantic/related/rows(GET->POST)/queries/applicable, one source resolver pkg/api/resolver.go, deleted the query-param session, --allow-writes gate off by default, --allow-opaque-sql operator grant, driver-bound SQL parameters proven with two CustomerId values), client (datatug-apps PR #72, libs/datatug/semantic/src/contract types+strict decoders+fixtures, SemanticApiService on POST with Scope, AgentContextService owns securityContextId, target selector, snapshot-retry, STALE_CONTEXT handling), and schema authority (datatug-core PR #312, tag v0.26.0, pkg/apicontract). Named follow-up debt, both already dispatched: datatug-cli still carries a provisional pkg/apicontract_local copy (lane S78 swaps it for core v0.26.0); datatug-apps fixtures are a hand-authored copy (lane S75 swaps them for the golden set).
**Evidence:** datatug/datatug-cli#212@7fc6efb, datatug/datatug-apps#72@8d27338, datatug/datatug-core#312@v0.26.0

Inventory every semantic/execute caller and choose one schema authority. Implement the exact Feature appendix on server and client contract fixtures; unify existing project environment/catalog targets behind a source resolver shared by semantic lookup and execution. Cut over all callers and remove superseded adapters. Do not create a competing source store. Verify empty/error envelopes and source resolution for SQLite, inGitDB and HTTP.

### Task 13: Converge protected execution and actual parameter binding

**Id:** task-13
**Verifies:** core-investigation-loop#ac:restricted-rows-and-columns-server, core-investigation-loop#ac:native-sql-labelled, core-investigation-loop#ac:hidden-column-refused
**Depends-On:** 12
**Status:** complete
**Implemented-by:** datatug/datatug-cli@e0e4c17 (main)
**Note:** Complete 2026-09-10 (lead session, per the founder's verbatim 'Yes, claim Task 13 and reconcile the ACL branches yourself'). Provider chain landed and tagged: dal-go/dalgo v0.80.0 (#158, layered ACL), dal-go/dalgo2sql v0.12.0 (#178, bounded structured SQLite queries, validated sources, bound values), ingitdb/dalgo2ingitdb v0.4.0 (#8, owner policies, protected reads suppress derived values, identity realms), openvaultdb/openvaultdb-go v0.4.0 (#17, DTQL-over-HTTP ACL; OVDB itself stays parked per Feature Open Question). datatug-cli: #222 acceptance tests (Task 13 security matrix), #223 pins dalgo2sql v0.12.0 + dalgo2ingitdb v0.4.0 and wires protected reads through dbcopy.BackendRef.OpenProtected (ingitdb WithStoredOnlyReads so formula columns never leak hidden fields; pkg/accesspolicies stays the sole enforcement point); native SQL runs with driver-bound arguments (dal.QueryArg → sql.Named, pkg/secureread/native_sql.go), opaque SQL needs --allow-opaque-sql AND a per-principal opaque-query policy grant (S109 matrix rows 5–8). datatug-apps #85 guards the client. Follow-up recorded, not blocking: the structured DTQL→SQL path still uses dalgo2sql's legacy emitSQL (quoted constants) because dalgo2sql v0.12.0's StructuredQueryDialect rejects FROM aliases (customer-invoices.query.dtql uses alias i); lane S114 adds alias support upstream, then datatug-cli enables the dialect.
**Evidence:** dal-go/dalgo#158@8dece8b, dal-go/dalgo2sql#178@a7649d1, ingitdb/dalgo2ingitdb#8@17581d0, openvaultdb/openvaultdb-go#17@5e2a46e, datatug/datatug-cli#222@1220105, datatug/datatug-cli#223@e0e4c17, datatug/datatug-apps#85@fc84ada

Reconcile pending VM provider/SDK/server changes and staged CLI SQL arguments. Port configuration to external core. Protected reads fail closed for unsupported shapes; explicit privileged opaque grants remain separate. Apply driver-bound parameters on the server path and return only confirmed bindings. Test current and legacy routes, principal tampering, forbidden source/field access and actual row changes between two parameter values. Publish/converge provider dependencies before final acceptance.

**Annotation Amendment:** actor=S70 (Claude Sonnet 5); at=2026-09-09T17:08:31Z; reason=Link the layered-ACL reconciliation inventory/port-plan note for Task 13; before_sha256=0b52e35f4c56f8bd7dc0bead2c27c1f61b33608d91c8b6d37c03ce96803f6739
**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:04Z; reason=Record what Task 12 already delivered toward Task 13's scope and what remains.; before_sha256=afecb9b4164f7b4274050b03785fb4a741f9d19ce8656df35a01f00caaa47932
**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:21Z; reason=Preserve the prior reconciliation-inventory note's content while adding what Task 12 delivered toward Task 13's scope.; before_sha256=548f1762fa5d0f5f2ea684956be5fe43b7cdef9109ceedbd6c47e08afca221cc
### Task 14: Complete bounded HTTP browser execution

**Id:** task-14
**Verifies:** core-investigation-loop#ac:http-source-in-demo, core-investigation-loop#ac:safe-http-and-lookups
**Depends-On:** 12, 13
**Status:** complete
**Implemented-by:** datatug/datatug-cli@5be73b7 (main)
**Note:** Complete 2026-09-10 (lane S111). Server: pkg/httpsource now defaults to dalgo2http ModeLive (the silent ModeLiveThenSnapshot fallback is gone; the two legacy CLI commands opt into it explicitly); exec/run_query drives ModeLive/ModeSnapshot from the request, validates a requested snapshotId (<queryId>@<recordedAt RFC3339>, derived from fixtures/http/<queryId>.json) before dispatch, maps upstream/DNS/timeout failures to SOURCE_UNAVAILABLE, snapshot miss to NOT_FOUND, snapshot on a non-HTTP query to INVALID_REQUEST, an unenforceable protected predicate to ACCESS_DENIED without leaking field names; snapshot provenance keeps the fixture's recorded time as observedAt; serve gains --exec-timeout (default 10s, clamped to the 30s ceiling) and --http-offline. Browser: result header states live or snapshot · recorded <date>; SOURCE_UNAVAILABLE with details.availableSnapshots renders an explicit 'use recorded snapshot' action that sends a second request with mode:snapshot; decoder accepts the sibling details key. Journey J2b (live: Country cell → country-facts → live + CAD; offline agent: honest failure → explicit snapshot → recorded 2026-09-09) passes in the real-agent CI job pinned to datatug-cli v0.22.0. Lead assumptions recorded in api-contract.md pending founder confirmation: the sibling top-level details.availableSnapshots shape (datatug-core ErrorBody has no extension field) and the --http-offline switch.
**Evidence:** datatug/datatug-cli#224@5be73b7, datatug/datatug-apps#87@620b4fa, datatug/datatug#39@334b580, dal-go/dalgo2http#3@ac48296

Reuse dalgo2http and descriptor/snapshot assets. Execute HTTP saved queries through the same browser endpoint and registry without dummy SQL target fields. Enforce configured hosts, address/redirect limits, time/size budgets, capability checks and cancellation. Live failure is an error; snapshot is a distinct explicit user request with recorded provenance. Verify a real public response and deterministic isolated negative tests separately.

**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:00Z; reason=Record provider-side HTTP bounds landed upstream (lane S72), status unchanged pending browser-side wiring.; before_sha256=4808250aa8258a78217199c1476f383928887c7141d5a26d5b13647169718925
### Task 15: Converge browser semantics and isolated context

**Id:** task-15
**Verifies:** core-investigation-loop#ac:typed-context-isolation, core-investigation-loop#ac:bound-from-selection, core-investigation-loop#ac:context-carries
**Depends-On:** 12
**Status:** complete
**Implemented-by:** datatug/datatug-apps@3c13cb0 (main)
**Note:** Complete 2026-09-10: exit condition met — the real-agent journey e2e (apps#83 CI job 'journey e2e (real datatug-cli agent)', datatug-cli v0.20.6) walks J1 bound-from-selection ('from selection' binding, clear, run sends the plain parameter) and J3 context-carries ('from context', chip disable empties, bindings listed in the results header) without intercepted responses or deep links; typed-context-isolation proven by datatug-apps#73 unit tests (typed 5 vs "5", conflict blocks unconfirmed run, scope switch opens an empty basket, isCurrentScope late-response-discard gate) plus J4 against a separate support principal. Task 17 still owns the full candidate/epilogue run.
**Evidence:** datatug/datatug-apps#73@1bb8406, datatug/datatug-apps#79@f5873b2, datatug/datatug-apps#82@438f1dc, datatug/datatug-apps#83@3c13cb0

Adapt grid/panel/query-page components to exact envelopes. Scope context and reactive agent selection by agent/project/environment/securityContextId. Implement typed equality, conflicts, missing inputs, visible overrides and late-response rejection. Selection cannot silently replace an open query binding. Test no facts cross project/principal boundaries and a new tab does not inherit prior context implicitly.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T18:32:01Z; reason=Record the Task 15 client implementation landing; before_sha256=88ec7fc81dbdd89b5cf7e6b151288cc2204a15bb3336bf0d76ce1545dd763fb8
### Task 16: Correct fixture policies and real transport acceptance

**Id:** task-16
**Verifies:** core-investigation-loop#ac:demo-project-shape, core-investigation-loop#ac:related-respects-policy, core-investigation-loop#ac:real-transport-and-parameter-effect
**Depends-On:** 13, 14, 15
**Status:** complete
**Implemented-by:** datatug/datatug-apps@620b4fa (main)
**Note:** Complete 2026-09-10. Lane S109: per-source support policies, pinned Brazilian customer 1 / Canadian customer 3 / investigation customer 5 counts verified against the pinned chinook-local.sqlite (sha256 matches phase1-acceptance.json), tests/go.mod on datatug-core v0.27.3 (datatug-demo-projects#15), support ACs name customers 1 and 3 and Open Question 1 resolved (datatug/datatug#36), 22-request browser-shaped acceptance matrix against real admin and support agents matched api-contract.md (Canada-only rows, Email hidden, hidden-column selection → ACCESS_DENIED, native SQL → UNSUPPORTED_PROTECTED_EXECUTION without --allow-opaque-sql and per-principal grant with it, forged as=/role= ignored, null counts for support vs real counts for admin, forged/cross-scope lookupIds rejected without leakage, two-value parameter effect, cross-scope securityContextId → STALE_CONTEXT). Real transport acceptance on the final candidate: the datatug-apps#87 journey CI job (datatug-cli v0.22.0) walks J1–J4, J2b (HTTP live and explicit snapshot) and Epilogues A/B (datatug-apps#86) end to end without interception.
**Evidence:** datatug/datatug-demo-projects#15@1c98f2c, datatug/datatug#36@a21f464, datatug/datatug-apps#87@620b4fa, datatug/datatug-apps#86@4d242ad, datatug/datatug-cli#224@5be73b7

Keep the real demo assets. Declare policies for Customer, Invoice and support-notes independently; choose Brazilian/Canadian IDs and expected counts from the pinned dataset. Use selected Country and retrieved currency facts rather than guessing from CustomerId. Preserve pending semantic demo tests. Exercise actual browser-shaped requests against the real server, including forged lookup IDs, policy-safe counts and unauthorized query metadata omission.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-10T07:03:54Z; reason=Reclaim Task 16 from Codex (out of limits) for the lead session; before_sha256=0e0e210e02ce3c21a6e4c4c2ee04063cde7a2a197a7a4a4b0fffb64b351e53c5
**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-10T07:30:48Z; reason=Record lane S109's landed Task 16 work; status stays in_progress until Task 14 (bounded HTTP browser execution, lane S111) lands; before_sha256=e1743968e58e8ff384de20931df90bc9cd63f76fbaa25785f2d9b895bf4473ca
### Task 17: Prove J1–J4 and finish one reproducible release candidate

**Id:** task-17
**Verifies:** core-investigation-loop#ac:journey-walks
**Depends-On:** 3, 8, 9, 10, 11, 16
**Status:** complete
**Implemented-by:** datatug/datatug-apps@8653bd9 (main)
**Note:** Complete 2026-09-11 (lead session). Candidate revision set, pinned: datatug/datatug-core v0.27.3 (a5f10fe); datatug/datatug-cli v0.23.0 (6f0b2d6; GitHub release published 2026-09-10T12:07:50Z, 6 assets; the published darwin_arm64 binary smoke-tested against demo-project-1: agent-info, queries/all_queries, catalog-tables); datatug/datatug-apps main 8653bd9 (#91 real in-app navigation with click-through journeys, #89 direct-nav fixes, #87 HTTP live/snapshot, #86 epilogues); datatug/datatug-demo-projects 1c98f2c; dal-go/dalgo v0.80.0; dal-go/dalgo2sql v0.13.0; ingitdb/dalgo2ingitdb v0.4.0; dal-go/dalgo2http v0.2.0. Authoritative CI on that exact set: datatug-apps PR #91 run 34477338700 and post-merge main run 34477823015 (test, build, e2e, journey e2e against datatug-cli v0.23.0 walking J1–J4, J2b, Epilogues A/B and direct-nav from the project URL only — all success); datatug-cli Go CI success on 6f0b2d6. Deploy receipt: Cloudflare Workers Builds builds https://datatug.app from datatug-apps main on every push; the GitHub commit status 'Workers Builds: datatug-apps' on 8653bd9 is success at 2026-09-10T12:40:23Z (no workflow, no token in CI; apps#54 closed as redundant). Journey gate reconciled: on 2026-09-11 the founder set branch protection on datatug-apps main requiring build, test, e2e and 'journey e2e (real datatug-cli agent)' (verified via the branches/main/protection API). Follow-ups, not blocking: agent-info reports version 0.0.1 instead of the release version; build info (version, git hash, build timestamp) in the app menu plus /build-info.json is lane S124.
**Evidence:** datatug/datatug-apps#91@8653bd9, datatug/datatug-apps#89@fd89eb1, datatug/datatug-apps#87@620b4fa, datatug/datatug-apps#86@4d242ad, datatug/datatug-cli#227@6f0b2d6, datatug/datatug-cli@v0.23.0, datatug/datatug-core@v0.27.3, datatug/datatug-demo-projects@1c98f2c

Refresh/rebase the journey and deploy candidates, pin or record every companion revision, fail required CI on missing setup, enable J4, remove direct-entry shortcuts and prove both epilogues. Obtain authoritative candidate CI conclusions. Reconcile one Cloudflare deploy owner/path with the required journey gate. Report exact remote/release/deploy receipts and required smoke evidence separately; follow authorization for landing and receipt-aware cleanup.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-10T08:39:10Z; reason=Record Task 17 prerequisites landed and what remains; before_sha256=9bc2d673c5aa100cc344676926194be8bd8ca5f48f0fa2dcc4856e23f2125a69
## Deferred AC Coverage

None — every acceptance criterion of the Feature is covered by a task above.

## Risks

- A completed component can hide incompatible contracts. Mitigation: exact serialized
  consumer/provider tests followed by whole-journey acceptance.
- Pending ACL work references the deleted vendored core configuration path. Port it
  upstream; preserve tests and converge provider releases rather than mirroring code.
- Opaque SQL cannot satisfy protected row/column guarantees. Refuse it in that profile;
  use structured/DTQL for restricted journeys, not a warning-only bypass.
- Public endpoints may fail. Verify live success separately; deterministic failure and
  explicit snapshot tests must not be reported as live acceptance.
- Target branches may move while integrating. Use current refs and exact candidate CI;
  revalidate changed candidates and preserve receipts/worktrees until remote proof.
- Numeric budgets never waive source coverage, context isolation or a required AC.

## Open Questions

None blocking the specified local investigation slice. Unrelated unbound CLI commands
belong to a separate surfaces maintenance task. Broader credential/hosting/board
choices remain in their owning Features and are not assigned to this plan.

---
*This document follows the https://specscore.md/plan-specification*
