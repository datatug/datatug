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
**Status:** in_progress

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

Playwright fixture that builds/launches `datatug serve --project <demo> --as admin --role admin` on
a free port, starts the web dev server pointed at it, and tears both down; a first test
walks J1 (open demo → Album → rows) asserting the agent log line. Non-goals: J2–J4
(task 10).

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
**Status:** in_progress

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
**Status:** in_progress

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
**Note:** Complete 2026-09-10: J1, J2, J3 and J4 pass single-worker against datatug-cli built from main aba31f6 (PR #221) and datatug-apps main 438f1dc (PR #82) — lanes S103/S104. J4 runs against its own --as support --role support agent (no intercepted responses, no deep-link shortcuts): only Canadian rows, Email hidden, limitation header on the table and query pages, Brazilian id → no rows, Canadian id → only its 7 invoices, hidden-column selection refused. Journey harness: apps/datatug-app/e2e/journey (README documents DATATUG_BIN/DATATUG_CLI_DIR/DATATUG_DEMO_DIR). Follow-up: the journey project is not yet a CI job (PR #48 being revived, lane S105); task 17 owns the release candidate and publication.
**Evidence:** datatug/datatug-apps#82@438f1dc, datatug/datatug-cli#221@aba31f6, datatug/datatug-apps#80@5bee04e, datatug/datatug-apps#81@814c104, datatug/datatug-cli#219@743b8d9, datatug/datatug-apps#79@f5873b2

Complete J1–J4 in the real-agent harness after corrective transport, executor,
HTTP, context and fixture integration. Cover the revised Feature mechanisms without
interception, missing-prerequisite skips or direct-entry shortcuts. This task produces
journey evidence only; it does not release or deploy. Task 17 owns final candidate,
epilogue, publication and remote receipt checks.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T20:16:49Z; reason=Record the first J1 end-to-end pass; before_sha256=de609446bc90d9cd53833aedc398394f8c3f56608698f28a2e94af2f6bd5bc8e
**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T21:24:53Z; reason=Record the J2 end-to-end pass and the J3 state; before_sha256=f16e13bd37c65c99c84fd98ef99bb1b4eedce6c9cdd1389ab7fe3706d18c6e84
**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T22:18:55Z; reason=Record the J3 pass, the canonical query-id convention and the J4 harness with its first failure; before_sha256=a1148755786cd6a853fca8379ca5f2babc4d25cf33059dd9e42b354abc6ffde7
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
**Status:** in_progress
**Note:** In progress under the lead session claim datatug-phase1-task13 (2026-09-10). Founder ruling, verbatim: "Yes, claim Task 13 and reconcile the ACL branches yourself". Lanes: S106 lands the provider layered-ACL branches (dalgo2sql, dalgo2ingitdb; openvaultdb-go if reachable) onto their mains; S107 ports the datatug-cli layered-ACL worktree (a7e1f62) onto today's main, resolving the five conflicts against the Task 12 transport freeze with the contract as authority; the datatug-apps route conflict follows. dal-go/dalgo #158 is merged (v0.80.0 of dal-go/dalgo). Reconciliation inventory: spec/research/2026-09-09-layered-acl-reconciliation.md.
**Evidence:** spec/research/2026-09-09-layered-acl-reconciliation.md, dal-go/dalgo#158@8dece8b

Reconcile pending VM provider/SDK/server changes and staged CLI SQL arguments. Port configuration to external core. Protected reads fail closed for unsupported shapes; explicit privileged opaque grants remain separate. Apply driver-bound parameters on the server path and return only confirmed bindings. Test current and legacy routes, principal tampering, forbidden source/field access and actual row changes between two parameter values. Publish/converge provider dependencies before final acceptance.

**Annotation Amendment:** actor=S70 (Claude Sonnet 5); at=2026-09-09T17:08:31Z; reason=Link the layered-ACL reconciliation inventory/port-plan note for Task 13; before_sha256=0b52e35f4c56f8bd7dc0bead2c27c1f61b33608d91c8b6d37c03ce96803f6739
**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:04Z; reason=Record what Task 12 already delivered toward Task 13's scope and what remains.; before_sha256=afecb9b4164f7b4274050b03785fb4a741f9d19ce8656df35a01f00caaa47932
**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:21Z; reason=Preserve the prior reconciliation-inventory note's content while adding what Task 12 delivered toward Task 13's scope.; before_sha256=548f1762fa5d0f5f2ea684956be5fe43b7cdef9109ceedbd6c47e08afca221cc
### Task 14: Complete bounded HTTP browser execution

**Id:** task-14
**Verifies:** core-investigation-loop#ac:http-source-in-demo, core-investigation-loop#ac:safe-http-and-lookups
**Depends-On:** 12, 13
**Status:** planning
**Note:** Provider-side HTTP bounds landed in dal-go/dalgo2http PR #3 (lane S72): https-only descriptors, guarded dialer (private/loopback/link-local/metadata-service addresses + DNS rebinding refused), redirects refused, 2 MiB response body cap with a RESPONSE_TOO_LARGE-style error, projection applied at the adapter boundary, no secrets sent in query strings. Status unchanged: this is the provider/adapter half only; the browser-side bounded HTTP execution (task-14's own Depends-On 12, 13) still needs wiring through the exec/run_query endpoint.
**Evidence:** dal-go/dalgo2http#3@ac48296

Reuse dalgo2http and descriptor/snapshot assets. Execute HTTP saved queries through the same browser endpoint and registry without dummy SQL target fields. Enforce configured hosts, address/redirect limits, time/size budgets, capability checks and cancellation. Live failure is an error; snapshot is a distinct explicit user request with recorded provenance. Verify a real public response and deterministic isolated negative tests separately.

**Annotation Amendment:** actor=lead-session; at=2026-09-09T17:45:00Z; reason=Record provider-side HTTP bounds landed upstream (lane S72), status unchanged pending browser-side wiring.; before_sha256=4808250aa8258a78217199c1476f383928887c7141d5a26d5b13647169718925
### Task 15: Converge browser semantics and isolated context

**Id:** task-15
**Verifies:** core-investigation-loop#ac:typed-context-isolation, core-investigation-loop#ac:bound-from-selection, core-investigation-loop#ac:context-carries
**Depends-On:** 12
**Status:** in_progress
**Note:** In progress: lane S75 landed datatug-apps#73 (typed scope-keyed context, binding resolver with origin tiers/conflict/ambiguous/missing blocking, late-response discard, golden fixtures from datatug-core v0.26.0 with manifest drift test). Status stays in_progress until the J1 journey e2e walks bound-from-selection and context-carries without crutches (lane S79 fixes the J1 deep-link blocker; Task 17 owns the journey run).
**Evidence:** datatug/datatug-apps#73@1bb8406

Adapt grid/panel/query-page components to exact envelopes. Scope context and reactive agent selection by agent/project/environment/securityContextId. Implement typed equality, conflicts, missing inputs, visible overrides and late-response rejection. Selection cannot silently replace an open query binding. Test no facts cross project/principal boundaries and a new tab does not inherit prior context implicitly.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-09T18:32:01Z; reason=Record the Task 15 client implementation landing; before_sha256=88ec7fc81dbdd89b5cf7e6b151288cc2204a15bb3336bf0d76ce1545dd763fb8
### Task 16: Correct fixture policies and real transport acceptance

**Id:** task-16
**Verifies:** core-investigation-loop#ac:demo-project-shape, core-investigation-loop#ac:related-respects-policy, core-investigation-loop#ac:real-transport-and-parameter-effect
**Depends-On:** 13, 14, 15
**Status:** in_progress
**Note:** Reclaimed by the lead session on 2026-09-10 — founder, verbatim: "Codex run out of limits so Task 16 got stuck. Can you check worktrees on this machine and VM find the WIP, reclaim and continue it work?" No uncommitted Codex WIP existed on either machine; Codex's landed Task 16 work is datatug-demo-projects 09407c3 (per-source support policy). Lane S109 continues: explicit per-source policies, pinned Brazilian/Canadian ids and counts, selected-country/retrieved-currency facts, forged-lookup and cross-scope transport acceptance against the real server, and the Open Question 1 decision (support journey uses pinned Canadian/Brazilian ids). Final browser-shaped acceptance is the J1–J4 journey CI job (datatug-apps#83).
**Evidence:** wb-claim:datatug-phase1-task16 (reclaimed), datatug/datatug-demo-projects@09407c3

Keep the real demo assets. Declare policies for Customer, Invoice and support-notes independently; choose Brazilian/Canadian IDs and expected counts from the pinned dataset. Use selected Country and retrieved currency facts rather than guessing from CustomerId. Preserve pending semantic demo tests. Exercise actual browser-shaped requests against the real server, including forged lookup IDs, policy-safe counts and unauthorized query metadata omission.

**Annotation Amendment:** actor=projects-c5 (Claude Fable 5.1); at=2026-09-10T07:03:54Z; reason=Reclaim Task 16 from Codex (out of limits) for the lead session; before_sha256=0e0e210e02ce3c21a6e4c4c2ee04063cde7a2a197a7a4a4b0fffb64b351e53c5
### Task 17: Prove J1–J4 and finish one reproducible release candidate

**Id:** task-17
**Verifies:** core-investigation-loop#ac:journey-walks
**Depends-On:** 3, 8, 9, 10, 11, 16
**Status:** planning

Refresh/rebase the journey and deploy candidates, pin or record every companion revision, fail required CI on missing setup, enable J4, remove direct-entry shortcuts and prove both epilogues. Obtain authoritative candidate CI conclusions. Reconcile one Cloudflare deploy owner/path with the required journey gate. Report exact remote/release/deploy receipts and required smoke evidence separately; follow authorization for landing and receipt-aware cleanup.

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
