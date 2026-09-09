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

Implements the [core investigation loop](../features/core-investigation-loop/README.md)
across `datatug-cli` (server, resolver, access control), `datatug-apps` (web UI),
`datatug-core` (two model additions) and `datatug-demo-projects` (the demo knowledge
project). Tasks 1–4 are the Phase-0 spine repair from [ROADMAP.md](../../ROADMAP.md)
(build, `serve`, agent path contract, demo project); tasks 5–10 are the loop. A fresh
coding session starts at task 1 with no product decisions outstanding except the three
assumptions (A1–A3) listed in the Feature, which default to the recommendations there.

Current-state references a coding agent needs first: [current-state audit](../research/2026-09-09-current-state-audit.md) §"Reuse inventory", [Web UI audit](../research/2026-09-09-web-ui-audit.md) §"Reuse pointers".

## Approach

- **Repair before build.** Nothing in tasks 5–10 is demonstrable until `datatug serve`
  starts and the web client reaches it; tasks 1–3 make the existing prototype real
  without adding concepts.
- **Server owns meaning and enforcement.** The resolver (mappings, related lookups,
  applicable queries) and the executor (`accesspolicies.Run` on every read) live in the
  CLI/server. The web UI only renders what the server returns. This keeps the TUI and a
  future `datatug context` verb able to reuse the same code.
- **Reuse over rewrite.** Store wiring copies `project_base_command.go`'s working path;
  execution reuses `pkg/accesspolicies`; the grid stays `@sneat/datagrid`; related cards
  extend `env-db-table`'s FK cards; parameter inputs reuse `input-parameters`.
- **Model changes only where required.** Two additions in `datatug-core/pkg/datatug`:
  `EntityField.Mappings []PhysicalRef` and `QueryDef.Type = "DTQL"`. Because
  `datatug-cli` carries a byte-identical vendored copy of `pkg/datatug`, task 1 first
  restores the module dependency (the commented-out `replace` in `datatug-cli/go.mod`
  points the way); if that proves disproportionate, the change is mirrored into both
  copies and the drift is recorded in the task's PR.
- **Security is tested, not asserted.** Task 6 ships the negative tests; task 10 walks
  the journey with a restricted principal.
- **Parallel waves.** Wave A: tasks 1, 4, 5 (three repos, independent). Wave B: tasks
  2, 3, 6, 7, 11 (after 1; 7 also after 5; 11 after 4 and 6). Wave C: tasks 8, 9 (after 2, 6, 7). Wave D:
  task 10. One WB worktree per task; the landing owner batches compatible branches
  onto one integration branch per repository and validates once.
- **Integration order.** Land `datatug-core` model change → `datatug-cli` (tasks 1, 5,
  6, 7) → `datatug-demo-projects` (task 4, needs the model) → `datatug-apps` (tasks 2,
  3, 8, 9) → task 10 e2e against pinned versions of all four.
- **Definition of done.** All acceptance criteria of the Feature pass in the journey
  e2e against the real agent; `main` of each repository builds and tests green; the
  demo repository validates; the Implementation Matrix rows read *Implemented*; a
  release of `datatug-cli` and a deploy of `datatug.app` carry the work.

## Tasks

### Task 1: Build, serve and store wiring in datatug-cli

**Id:** task-1
**Verifies:** core-investigation-loop#ac:agent-path-contract-server
**Depends-On:** —
**Status:** complete
**Implemented-by:** datatug/datatug-cli@940902b65a4c965c3b524d9be0ebf8e2b9125d2c (main)

*Progress 2026-09-09:* Audit evidence: datatug-cli PR [#192](https://github.com/datatug/datatug-cli/pull/192) and [#195](https://github.com/datatug/datatug-cli/pull/195) merged; related CI/dependency fixes are in [#191](https://github.com/datatug/datatug-cli/pull/191), [#188](https://github.com/datatug/datatug-cli/pull/188), and [#189](https://github.com/datatug/datatug-cli/pull/189). The required module-dependency restoration remains open in [#197](https://github.com/datatug/datatug-cli/pull/197), so this task is not complete.

*Progress 2026-09-09 (14:30 UTC):* Reconciliation evidence: datatug-cli PR [#197](https://github.com/datatug/datatug-cli/pull/197) (module dependency on `datatug-core` restored, vendored `pkg/datatug-core` copy deleted) and [#198](https://github.com/datatug/datatug-cli/pull/198) (serve link carries scheme/host/port) are merged. The `recover()`-around-the-panic test is already replaced by `TestNewDatatugStoreFactory` (`pkg/server/http_server_test.go`, landed via #192); the `go build ./...`/`go test ./...` CI gate runs on every pull request via the required `strongo_workflow` check (deliberately no `paths:` filter, so it also gates dependency/Renovate PRs) and is green on `main` HEAD (`c1d1e30`, run 34350654432, conclusion `success`). Every literal AC item is now evidenced; task is complete.

Fix the go-github v91 imports so `main` builds; add a CI gate that runs `go build ./...`
on dependency PRs. Make `datatug serve` start: nil-check `config.Server`, bind
`--host/--port/--project/--as/--role/--group` flags, pass a real `Handler` to
`endpoints.RegisterDatatugHandlers`, and implement `storage.NewDatatugStore`/`GetStore`
over `filestore` (mirroring `project_base_command.go`). Restore the module dependency on
`datatug-core` (delete the vendored `pkg/datatug-core` copy) or mirror-and-record.
Replace the `recover()`-around-the-panic test with a real handler test. Tests: unit
tests for store wiring and flag binding; `go test ./...` green. Non-goals: ACL, new
endpoints, TUI.

### Task 2: Web agent client contract and demo entry

**Id:** task-2
**Verifies:** core-investigation-loop#ac:agent-path-contract-client, core-investigation-loop#ac:first-result-two-clicks
**Depends-On:** 1
**Status:** in_progress

*Progress 2026-09-09:* datatug-apps PR [#43](https://github.com/datatug/datatug-apps/pull/43) merged the `/datatug` client paths, host:port parsing, providers, demo link, route fix, and shell feature flag. The real-agent contract/journey verification is tracked separately in task 3, so this task remains in progress pending the complete acceptance path.

*Progress 2026-09-09 (14:30 UTC):* datatug-apps PR [#63](https://github.com/datatug/datatug-apps/pull/63) fixed `parseDatatugStoreRef()` so `http-`/`https-`-prefixed `host:port` store ids round-trip correctly (root cause: the bare `host:port` regex short-circuit swallowed the scheme prefix before delegating to `@sneat/core`'s `parseStoreRef()`), and added this task's own contract test (`e2e/journey/store-id-scheme.spec.ts`, real `datatug serve`, no interception). That new spec could not get a clean pass in this environment: PR #63 reports the *pre-existing, unmodified* `journey.spec.ts` J1 test also failing with an unrelated network error, reproduced 3x — the same server panic task 3 root-caused in #57 (`GetAuthTokenFromHttpRequest is nil`), still open (lane S36, uncommitted in `datatug-cli/.worktrees/phase1-cli-auth-hook`). Task remains in progress: the contract fix is merged but its own required test is not yet a passing, verified run.

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

*Progress 2026-09-09:* datatug-apps PR [#46](https://github.com/datatug/datatug-apps/pull/46) merged the real-agent Playwright harness and J1 skeleton, but its validation reports J1 failing against current CLI `main`; follow-up PR [#48](https://github.com/datatug/datatug-apps/pull/48) is open. Task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-apps PR [#57](https://github.com/datatug/datatug-apps/pull/57) root-caused J1's failure to two `datatug-cli`/`datatug-demo-projects` server defects, not a frontend gap: (1) every `handle(...)`-routed endpoint (including `project_summary`, J1's first request) panics with `GetAuthTokenFromHttpRequest is nil` because `datatug serve` never sets that hosted-backend auth hook; (2) `exec/select` failed to load `demo-project-1`'s `<id>.env.json` environment files (fixed upstream in datatug-core [#308](https://github.com/datatug/datatug-core/pull/308), v0.22.0, but datatug-cli's `go.mod` is still pinned to v0.21.0 pending the open [#202](https://github.com/datatug/datatug-cli/pull/202)). Follow-up PR #48 remains open. The auth-hook fix is in progress but not yet pushed (lane S36, worktree `datatug-cli/.worktrees/phase1-cli-auth-hook`, branch `fix/phase1-cli-auth-hook`, uncommitted as of this reconciliation). Task remains in progress; J1 does not pass yet.

Playwright fixture that builds/launches `datatug serve --project <demo> --as admin` on
a free port, starts the web dev server pointed at it, and tears both down; a first test
walks J1 (open demo → Album → rows) asserting the agent log line. Non-goals: J2–J4
(task 10).

### Task 4: Demo knowledge project and credential guard

**Id:** task-4
**Verifies:** core-investigation-loop#ac:demo-project-shape
**Depends-On:** —
**Status:** complete
**Implemented-by:** datatug/datatug-demo-projects@6ca9ddd1452ab1179a6c1e019b400b18bcf4a844 (main)

*Progress 2026-09-09:* datatug-demo-projects PR #8 landed the project (entities incl. `Country`, queries, support-notes, HTTP query defs + snapshot fixtures, policies); the credential rule landed in datatug-core PR #302; declared `Mappings` on the demo's fields are a follow-up once the CLI consumes that model version.

*Progress 2026-09-09:* Reconciliation evidence: datatug-demo-projects PR [#8](https://github.com/datatug/datatug-demo-projects/pull/8) and mapping/DTQL follow-up [#9](https://github.com/datatug/datatug-demo-projects/pull/9) are merged, and datatug-core PR [#302](https://github.com/datatug/datatug-core/pull/302) is merged. The project shape, mappings, queries, support-notes collection, policies, HTTP definitions, fixtures, and validation tests are evidenced; task is complete.

In `datatug-demo-projects/demo-project-1`: add `Customer` and `Invoice` entities with
declared mappings (`Customer.ID → chinook.Customer.CustomerId`, `Customer.Email →
chinook.Customer.Email`, `Invoice.ID → chinook.Invoice.InvoiceId`, plus
`Customer.ID → support-notes.CustomerId`), semantic types, and `NamePatterns` for
`CustomerId`/`customer_id`; three parameterized queries (`customer-invoices` as DTQL,
`invoice-lines` as SQL, `customer-purchases-by-genre` as SQL) with `Parameters[].Meta`
set; an inGitDB `support-notes` collection (8 records keyed by `CustomerId`, 2 for
customer 5); `policies/` with `admin` and `support` (Canada rows on `Customer`, `Email`
hidden); a `Country` entity with `Country.Name → chinook.Customer.Country` and
`Country.Currency`; an HTTP `QueryDef` for CountriesNow (`name` parameter tagged
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

*Progress 2026-09-09:* datatug-core PRs [#302](https://github.com/datatug/datatug-core/pull/302) and [#303](https://github.com/datatug/datatug-core/pull/303) are merged for the model and pure resolver foundations. The CLI semantic endpoint and module integration are not evidenced by merged work yet; task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-cli PR [#200](https://github.com/datatug/datatug-cli/pull/200) merged `GET /datatug/semantic/columns` against the module's real `datatug.Entity`/`QueryDef` types, with table-driven resolver tests (`semantic_columns_test.go`, `semantic_fixture_test.go`) proving declared-over-inferred provenance and a real-demo-project integration test (`TestApplicable_RealDemoProject`, gated on `DATATUG_DEMO_PROJECT_DIR`, skips cleanly in CI). Every literal AC item — model, resolver, endpoint, and both named test shapes — is now evidenced; task is complete.

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
**Depends-On:** 1
**Status:** complete
**Implemented-by:** datatug/datatug-cli@a30132d9fee1f8f670b06bb61f49c84fa119c90e (main)

*Progress 2026-09-09:* datatug-cli PR [#193](https://github.com/datatug/datatug-cli/pull/193) merged the policy-enforced executor and security tests. Serve-path wiring/principal integration is still open in [#199](https://github.com/datatug/datatug-cli/pull/199); task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-cli PR [#199](https://github.com/datatug/datatug-cli/pull/199) merged the serve-path wiring: `resolveServeSession` builds the `secureread.Session` from `--as/--role/--group` and refuses to start with a policy set and no principal; `GET /datatug/exec/select` and the new `POST /datatug/exec/run_query` dispatch through `RunStructured`/`RunDTQL`/`RunNativeSQL`; `agent-info` reports the principal and no endpoint accepts one from the request. `pkg/server/security_matrix_test.go` covers `restricted-rows-and-columns-server`, `hidden-column-refused`, `native-sql-labelled`, and `dtql-query-runs`. Every literal AC item is now evidenced; task is complete.

Route `exec/select`, saved-query runs and the new `POST /datatug/exec/run_query`
through `accesspolicies.Run` with the project policy set (`<project>/policies/` then
`~/.datatug/policies`) and the principal fixed by `serve --as/--role/--group`; open
sources through `pkg/dbcopy/url.go` (`sqlite://` via `dalgo2sqlite`, `ingitdb://`);
execute `DTQL` queries via `dtql.Deserialize`; execute `SQL` text in a read-only
session under source-level policy only and label it; return structured `limitations`
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

*Progress 2026-09-09:* datatug-core PR [#304](https://github.com/datatug/datatug-core/pull/304) merged related/applicable model logic and tests. The datatug-cli HTTP endpoint wiring is still in progress with no merged PR evidenced; task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-cli PR [#200](https://github.com/datatug/datatug-cli/pull/200) merged `GET /datatug/semantic/related`, `GET /datatug/semantic/related/rows`, and `POST /datatug/queries/applicable`, routed through task 6's `secureread` executor. Its fixture tests match the plan's own named scenarios exactly: Customer 5 → Invoices count 7 (`referencedBy`) and Support notes count 2 (`sameField`); `customer-invoices` applicable with a resolution chain, `invoice-lines` not-yet (missing `Invoice.ID`); a restricted principal sees zero rows and no true count (`AC:related-respects-policy`). Every literal AC item is now evidenced; task is complete.

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

*Progress 2026-09-09:* datatug-apps PR [#47](https://github.com/datatug/datatug-apps/pull/47) merged the semantic library, context-panel primitives, context service, and limitation-header foundation. Existing page/grid wiring remains a follow-up; task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-apps PR [#59](https://github.com/datatug/datatug-apps/pull/59) wires `env-db-table.page.ts` to `GET /datatug/semantic/columns` (column markers) and adds a new `ContextPanelComponent` side panel on cell click, with unit tests against fixture responses. Two literal AC items are still unmet: the hard-coded `CellPopoverComponent` (`@sneat/datagrid`) sample content was left in place for the pre-existing FK popup rather than replaced — the new context panel was added alongside it, not instead of it; and the wiring was built and tested against a synthetic fixture before `/datatug/semantic/columns` existed on `datatug-cli` main (PR #59 merged 17 seconds before PR #200), so its `source` parameter mapping carries its own unverified `ASSUMPTION` comment. J2, the e2e proof of this task's grid/panel behaviour, fails today (blocked by the same server panic as task 3). Task remains in progress.

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

*Progress 2026-09-09:* datatug-apps PR [#47](https://github.com/datatug/datatug-apps/pull/47) merged the Investigation Context service and reusable UI foundation, but query-page binding and project/page integration are not evidenced as merged. Task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* datatug-apps PR [#59](https://github.com/datatug/datatug-apps/pull/59) wires the query page to `InvestigationContextService.bindingsFor()` and `POST /datatug/exec/run_query` (merged in datatug-cli #199), renders `LimitationHeaderComponent` and the applied-bindings list, and adds `InvestigationContextBarComponent` to the project shell — with unit-test coverage for the binding-precedence/limitation behaviour this task's own AC calls for. One literal AC item is unmet: the "Variables" route still resolves to `VariablesPageComponent` (`datatug-routing-proj.ts:60`), not the Investigation Context screen. J3, the e2e proof of context-carries behaviour, fails today (same server panic as task 3). Task remains in progress.

An `InvestigationContextService` (signals, session-scoped) holding
`{entity, field, value, origin, enabled}` items and the environment; the "Variables"
route becomes the Investigation Context screen; every project page shows the collapsed
context with a count; the query page binds `ParameterDef.Meta` parameters from the
selection or context (chip with origin, clear/override), and the results header lists
bindings applied and policy limitations (from task 6). Run uses `POST
/datatug/exec/run_query`. Tests: binding precedence (selection over context), disabled
chip empties the value, header renders limitations. Non-goals: persistence of context
to the project.

### Task 10: Journey e2e J1–J4 and release

**Id:** task-10
**Verifies:** core-investigation-loop#ac:journey-walks
**Depends-On:** 3, 8, 9
**Status:** in_progress

*Progress 2026-09-09:* No merged PR or passing J1–J4 journey evidence was found. PR [#46](https://github.com/datatug/datatug-apps/pull/46) leaves J2–J4 as fixme and reports J1 failing; the implementation matrix remains Planned. Task remains in progress.

*Progress 2026-09-09 (14:30 UTC):* No change to the journey verdict since C4's note. datatug-apps PR [#59](https://github.com/datatug/datatug-apps/pull/59) filled in J2/J3 (deliberately left un-skipped despite failing, per its own brief) and PR [#57](https://github.com/datatug/datatug-apps/pull/57) root-caused the shared blocker as a `datatug-cli` server panic (`GetAuthTokenFromHttpRequest is nil`) plus an environment-file naming mismatch, neither fixable from `datatug-apps`. The auth-hook fix is in progress but unpushed (lane S36, `datatug-cli/.worktrees/phase1-cli-auth-hook`). J1–J4 do not pass; the Implementation Matrix stays Planned. Task remains in progress.

Extend the harness with J2 (cell → related + applicable → run bound), J3 (context
carries, chip disable), J4 (`--as support`: rows filtered, `Email` hidden, header
states limitation, hand-crafted request refused), asserting mechanism (server log
entries, `limitations`, chip origins). Update the Implementation Matrix rows, release
`datatug-cli`, deploy `datatug.app`, and update `datatug-demo-projects` README.

### Task 11: HTTP reference-data source executor

**Id:** task-11
**Verifies:** core-investigation-loop#ac:http-source-in-demo
**Depends-On:** 4, 6
**Status:** complete
**Implemented-by:** datatug/datatug-cli@2dd8afee2e44b66df0a49dedfbb4239856965357 (main)

*Progress 2026-09-09:* dal-go/dalgo2http PR [#2](https://github.com/dal-go/dalgo2http/pull/2), datatug-cli PR [#196](https://github.com/datatug/datatug-cli/pull/196), and the demo definitions/fixtures in datatug-demo-projects PR [#8](https://github.com/datatug/datatug-demo-projects/pull/8) are merged. The adapter, snapshot fallback, HTTP source registration, and unit/demo fixture validation satisfy this task; task is complete.

Founder ruling 2026-09-09 (verbatim): *"I want to have in MVP: inGitDB, SQLite, HTTP
endpoint."* Executor: the generic DALgo adapter **`dal-go/dalgo2http`** (founder 2026-09-09:
*"dalgo2http - yes"*): declarative per-collection descriptors, GET only, equality pushdown on
declared parameters, fail-closed on anything else, recorded snapshots with `live`/`snapshot`
provenance. In `datatug-cli`: register `http://`/`https://` in `pkg/dbcopy/url.go` to open a
`dalgo2http` database from the project's HTTP `QueryDef`s, and surface the adapter's
provenance in the result header. Related lookups (task 7) and applicable queries treat
the HTTP source like any other. Tests: unit tests with a local `httptest` server and the
snapshot path; the demo's Country cell → Country facts and the exchange-rate query pass
in the journey harness with the network disabled. Non-goals: POST, auth headers, GraphQL,
HTTP steps in playbooks (Phase 2–3). Wave B, runs beside task 7.

## Deferred AC Coverage

None — every acceptance criterion of the Feature is covered by a task above.

## Risks

- **Storage indirection may be cheaper to delete than to fix** (task 1). Mitigation:
  the implementer records the choice in the PR; either path satisfies the tests.
- **Vendored model copy in `datatug-cli`** (task 1/5). Mitigation: restore the module
  dependency first; mirror only as a recorded fallback.
- **Ionic density for the context panel** (task 8). Mitigation: a plain `ion-list`
  first; PrimeNG island only if a tree is needed.
- **DALgo cannot filter native SQL rows** (task 6). Mitigation: explicit label; DTQL for
  the demo's headline query.
- **`dalgo2sqlite` vs `mattn/go-sqlite3` behaviour differences** (task 6). Mitigation:
  structured queries go through `dalgo2sqlite`; native SQL keeps the legacy driver in
  dev builds; release builds document the CGO limitation until Phase 2 removes it.
- **Public HTTP endpoints in a demo** (task 11). Mitigation: recorded snapshot fixture
  served when unreachable, labelled as such; the journey harness runs with the network
  disabled so the demo never depends on a third party being up.
- **Founder assumptions A2–A3** default to the Feature's recommendations; the concrete
  public endpoints (CountriesNow, Frankfurter) are the lead's choice and swap freely —
  REST Countries was already dead when task 4 recorded fixtures.

## Open Questions

- Should task 1 also delete `updateUrlConfig` and the unbound-flag commands
  (`show`, `dataset*`, `datasets`, `queries`), or leave them for Phase 2's
  fix-or-delete pass? Recommendation: leave; keep task 1 mechanical.

---
*This document follows the https://specscore.md/plan-specification*
