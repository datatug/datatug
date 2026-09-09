---
format: https://specscore.md/plan-specification
status: Draft
---
# Plan: Phase 2 — Library depth and provenance

**Status:** Draft
**Source:** none
**Date:** 2026-09-09
**Owner:** alex
**Supersedes:** —

## Summary

Implements [ROADMAP.md](../../ROADMAP.md) "Phase 2 — Library depth and provenance"
across `datatug-core` (semantic field types), `datatug-cli` (server endpoints, CLI
verbs, TUI), `datatug-apps` (web UI), `dal-go/dalgo2http` (HTTP adapter auth) and
`datatug-demo-projects` (fixtures for the new capabilities). It draws on seven Draft
Features whose capability tables mark a Phase-2 row: `executable-knowledge-library`
(capture, collections), `provenance` (result and mapping provenance), `meaning-inference`
(review/accept/correct), `heterogeneous-sources` (HTTP execution "(rest)"), `surfaces`
(hand-off, TUI), `agent-context` (`datatug context`), and `investigation-context`
(environment/time auto-narrowing — read in full because ROADMAP.md names it explicitly
even though it is not in this stream's originally-listed seven Features; see Open
Question 1). None of these Features is Approved yet; this Plan does not wait for that —
it cites their existing AC slugs as the traceable target, per this stream's brief.

**Preconditions already landed (state, not tasks; verified against the repositories,
not restated from any earlier summary — see Risks for one correction):**

- `datatug-core`: `EntityField.Mappings []PhysicalRef`, `QueryDef.Type = DTQL`, and the
  credential-validation rule (PR #302); the pure column→entity-field resolver in
  `pkg/semantic` (PR #303); related lookups and query applicability with resolution
  chains, also in `pkg/semantic` (PR #304).
- `datatug-cli`: `pkg/secureread`, the policy-enforced executor for structured, DTQL
  and native-SQL queries that every task below extends rather than re-derives (PR #193,
  merged after the Phase 0 `serve` spine repair).
- `dal-go/dalgo2http`: the read-only HTTP/JSON adapter (GET only, declarative
  descriptors, snapshot fixtures with `live`/`snapshot` provenance) — MVP merged.
- `datatug-apps`: `libs/datatug/semantic`, the web foundation library consuming the
  `pkg/semantic` endpoints.
- `datatug-demo-projects`: `demo-project-1` already carries HTTP `QueryDef`s
  (`queries/reference/country-facts.query.http`, `.../currency-rate.query.http`) with
  recorded fixtures, and `policies/customers.yaml` (admin/support principals).

## Journey

Phase 1 already delivered the loop's first pass: open the demo project, select
`Customer.ID = 5`, see related records across three sources and applicable queries with
resolution chains, run one, add it to context, and see it limited under a restricted
principal. Phase 2 is what the *next* investigation looks like once a team has used the
loop for a while, and what an agent or a TUI user — not only the web UI — can now do.

- **Start — the null-action baseline.** The user opens `datatug serve`'s web UI on the
  demo project and browses `Customer` exactly as in Phase 1. Nobody does anything else.
  Good result: nothing regresses; the Phase 1 loop is unchanged.
- **An ad-hoc lookup becomes a team asset.** While investigating customer 5 the user
  runs a related-record lookup that has no saved query behind it yet, then clicks *Save
  as team query*. Good result: a new `.query.json`/`.dtql.yaml` pair appears as an
  uncommitted Git change under `queries/customers/`, its semantic parameter and
  applicability already populated from the context that produced it — reusable next
  time without anyone re-typing it (task 2).
- **A suggested mapping is reviewed, not silently trusted.** On the entities page the
  user sees `Employee.ManagerId` marked `inferred` (a `NamePatterns` match). They open
  the review panel, see the resolution chain, and accept it. Good result: the entity
  file now carries a declared mapping, and its provenance record names the user and the
  timestamp — traceable later without grepping commit history (tasks 6–7).
- **A result explains itself.** The user reopens the *Customer invoices* result and asks
  "where did this come from." Good result: the answer names the source (`chinook`), the
  table (`Invoice`), the query that produced it and the policy limitation already shown
  in Phase 1 — the full chain, not only the limitation half (tasks 4–5).
- **The HTTP source authenticates.** The user (or the demo project's maintainer) adds a
  second HTTP reference query that needs an API key. Good result: it runs from both the
  CLI and the web UI, the key lives in `~/.datatug`, and `datatug validate` still finds
  no credential in any committed project file (task 8).
- **The context narrows by environment.** The user sets the active environment to
  `prod`; a `staging`-only query drops out of the applicable-queries panel. Good result:
  the panel visibly changes, and switching back restores the query — never a silent
  filter (task 11).
- **An agent asks the same question a human would.** Separately, an agent (or the user
  from a terminal) runs `datatug context --value Customer.ID=5` under the `support`
  principal. Good result: it lists the same related records and applicable queries the
  web context panel would show that principal — rows filtered, `Email` still hidden,
  never an unfiltered agent-only view (task 9).
- **Hand-off follows the user.** The user presses Ctrl+W on the TUI's *Customer
  invoices* query screen. Good result: the web UI opens the same query with the same
  selection, not the bare origin — and the TUI's own Queries panel, a placeholder
  through Phase 1, already listed that query against real project data (task 10).
- **Epilogue (a) — closes the browser.** The captured query and the accepted mapping
  are already ordinary Git-tracked or uncommitted files; closing the browser loses no
  state, and `datatug serve` still prints the URL to reopen the same screen — Phase 1's
  guarantee, unchanged.
- **Epilogue (b) — shares the result link.** A colleague opens the same "where did this
  come from" link; they see the identical provenance chain filtered by their own
  principal, and the URL carries no credential and no API key.

## Approach

- **Extend, don't reopen.** Phase 1 built the resolver (`pkg/semantic`), the
  policy-enforced executor (`pkg/secureread`), the applicable-queries chain and the
  three-source demo. Every task below adds a capability on that surface — capture,
  review, provenance, HTTP auth, environment, hand-off, TUI, the agent verb — rather
  than re-deriving it.
- **Server still owns meaning and enforcement.** Capture, mapping review, provenance and
  `datatug context` are new server endpoints or CLI verbs that reuse
  `accesspolicies.Run` and `pkg/secureread`; the web UI and the TUI only render what the
  server returns, matching Phase 1's rule so a future MCP face can reuse the same code
  (`agent-context` REQ:link-serve-brokered-query-builder).
- **Model changes stay additive.** `EntityField.Mappings` and `QueryDef.Type` already
  exist (Phase 1); Phase 2 adds a provenance sub-record to a mapping and to every
  result, plus new `KnownTypes` entries (`environment`, `time`, alongside
  `money`/`email`/`GUID`) — no existing field's shape changes.
- **Reuse over rewrite.** Capture reuses the applicable-queries resolution chain to
  pre-populate the new file; mapping review reuses the same `pkg/semantic` resolver
  Phase 1 built; `datatug context` reuses the exact endpoints task 7/8/9 of Phase 1
  expose to the web UI rather than a parallel agent-only path
  (`agent-context` REQ:agent-readable-context-resource forbids the alternative).
- **Parallel waves.** Wave A: tasks 1–4, 6, 8–10 (eight tasks, largely independent of
  each other; several land in `datatug-cli` and `datatug-apps` concurrently — the
  landing owner batches compatible same-repo branches onto one integration branch per
  repository and validates once, as Phase 1's Approach also states). Wave B: tasks 5, 7,
  11 (each depends on its Wave-A sibling: 5 on 4, 7 on 6, 11 on 1). Wave C: task 12 (the
  journey e2e, after every task whose UI or CLI surface it walks). One WB worktree per
  task.
- **Integration order.** `datatug-core` (task 1's field types) → `dal-go/dalgo2http`
  (task 8's auth headers) → `datatug-cli` (tasks 2, 3, 4, 6, 8, 9, 10's server/CLI/TUI
  halves) → `datatug-demo-projects` (task 1's `Employee` entity, task 8's second HTTP
  descriptor) → `datatug-apps` (tasks 2, 3, 5, 7, 10, 11's web halves) → task 12 e2e
  against pinned versions of all four.
- **Definition of done.** Every acceptance criterion named in a task's `Verifies:`
  passes against the real agent; `main` of each touched repository builds and tests
  green; the Implementation Matrix rows these tasks touch read *Implemented*; the demo
  project validates.

## Tasks

### Task 1: Semantic field types, environment/time kinds, and the Employee entity gap

**Id:** task-1
**Verifies:** semantic-layer#ac:semantic-type-visible-not-coercive, semantic-layer#ac:employee-entity-gap-recorded
**Depends-On:** —
**Status:** planning

In `datatug-core/pkg/datatug`: extend `KnownTypes` beyond `money`/`email`/`GUID` with
`environment` and `time`, inspectable and editable metadata only — never a silent
coercion of the stored value (`semantic-layer` REQ:semantic-field-types). In
`datatug-demo-projects/demo-project-1`: add the `Employee` entity with a declared
mapping, closing the gap `semantic-layer` REQ:demo-entity-coverage recorded but did not
resolve in Phase 1. Tests: type round-trip in the entities page and `datatug entity`
CLI; `datatug validate` passes with the new entity. Non-goals: probabilistic type
inference (Phase 5); auto-narrowing itself (task 11 consumes this task's `environment`
kind).

### Task 2: Capture-from-exploration

**Id:** task-2
**Verifies:** executable-knowledge-library#ac:capture-lands-as-git-diff
**Depends-On:** —
**Status:** planning

In `datatug-cli`: a `POST /datatug/library/capture` endpoint that takes the
related-lookup or applicable-query context that produced a result and writes a new
`.query.json`/`.dtql.yaml` pair into the project's `queries/` tree, with semantic
parameters and applicability already populated from that context
(`executable-knowledge-library` REQ:capture-from-exploration) — the write goes through
the same filestore path Phase 1 used for reads, never bypassing it. In `datatug-apps`:
a *Save as team query* action on the result screen and the related-lookup panel calling
this endpoint. Tests: the new file round-trips through `datatug entity list`/`query
run`; a component test asserts the action appears only when a lookup has no saved query
behind it yet. Non-goals: promotion scope (personal/team/org — Phase 3, needs
`sharing-scopes`); playbooks.

### Task 3: Query collection metadata — folder README and domain tag

**Id:** task-3
**Verifies:** executable-knowledge-library#ac:folder-recursion-works
**Depends-On:** —
**Status:** planning

In `datatug-core`: an optional `README.md` and a `domain` string per `QueriesFolder`,
additive to the existing recursive-folder model — the Feature's own Open Question
recommends keeping folders simple until proven otherwise, so this task adds only the two
fields ROADMAP.md names (README, domain tag), not owner/entry-point metadata. In
`datatug-cli`: read/write support; in `datatug-apps`: the folder tree shows the domain
tag and links the README. Tests: nested-folder fixture from the demo project still
resolves with the new optional fields present or absent. Non-goals: any folder metadata
beyond README + domain tag.

### Task 4: Result provenance — server

**Id:** task-4
**Verifies:** provenance#ac:result-shows-full-chain
**Depends-On:** —
**Status:** planning

In `datatug-cli/pkg/secureread`: extend `Result` with a `Provenance` record alongside
the existing `Limitations` (source, table/collection/endpoint, the query or lookup that
produced it, any transformation applied) on every read path — `exec/select`, saved-query
runs, `/datatug/exec/run_query`, related lookups and applicable-query runs
(`provenance` REQ:result-provenance, extending `core-investigation-loop`
REQ:limitation-visible's policy half to the source/query/transformation half). The HTTP
adapter's own `live`/`snapshot` provenance (`dal-go/dalgo2http/provenance.go`) becomes
one input to this chain, not a separate concept. Tests: table-driven — SQLite, inGitDB
and HTTP results each carry a complete, correctly-shaped chain; a restricted principal's
chain still names the source even when rows are filtered.

### Task 5: Result provenance — web UI

**Id:** task-5
**Verifies:** provenance#ac:result-shows-full-chain
**Depends-On:** 4
**Status:** planning

In `datatug-apps`: a "where did this come from" affordance on a result (grid header or
row) that renders task 4's `Provenance` chain — source, table/endpoint, query,
transformation, policy limitation — reusing the same limitation display Phase 1 task 9
already ships. Tests: component test with a fixture response covering all three source
kinds. Non-goals: a lineage-diagram visualization (out of scope per the Feature's "Not
Doing" section).

### Task 6: Mapping review, accept, correct and provenance — server

**Id:** task-6
**Verifies:** meaning-inference#ac:accept-promotes-to-declared, meaning-inference#ac:correction-overrides-suggestion, provenance#ac:mapping-shows-declarer-and-time
**Depends-On:** —
**Status:** planning

In `datatug-core`: a mapping-provenance record shape (who/what, declared/inferred/
corrected, timestamp) attached to `EntityField.Mappings`
(`provenance` REQ:mapping-provenance). In `datatug-cli`: `POST
/datatug/semantic/mappings/accept` (promotes an inferred mapping to declared) and `POST
/datatug/semantic/mappings/correct` (replaces a suggestion with an explicit mapping,
recording the correction so the rejected suggestion is not silently reused —
`meaning-inference` REQ:review-accept-correct-provenance); both write the entity file
and the provenance record together. Tests: accept promotes and records provenance;
correct overrides and does not resurface the same wrong suggestion on the next resolver
run. Non-goals: probabilistic suggestions (Phase 5).

### Task 7: Mapping review UI

**Id:** task-7
**Verifies:** meaning-inference#ac:accept-promotes-to-declared
**Depends-On:** 6
**Status:** planning

In `datatug-apps`: the entities page marks an `inferred` mapping distinctly, opens a
review panel showing the resolution chain, and offers *Accept*/*Correct* calling task
6's endpoints; an accepted or corrected mapping shows its provenance (who, when).
Tests: component test — accept updates the badge to `declared`; correct opens a column
picker and does not silently guess. Non-goals: probabilistic-suggestion confidence UI
(Phase 5).

### Task 8: HTTP QueryDef authentication and arbitrary descriptors

**Id:** task-8
**Verifies:** heterogeneous-sources#ac:http-query-executes, heterogeneous-sources#ac:secrets-guard-rejects-committed-credentials
**Depends-On:** —
**Status:** planning

Phase 1 task 11 shipped two hard-coded, key-less descriptors (CountriesNow,
Frankfurter). This task closes the "(rest)" half of `heterogeneous-sources`'
HTTP-source row: in `dal-go/dalgo2http`, support an auth-header field on a descriptor,
resolved from a local, non-committed credential store, never from the project file
(extends `heterogeneous-sources` REQ:secrets-outside-project-files to HTTP); in
`datatug-cli`, any project-declared HTTP `QueryDef` — not only the two demo
descriptors — resolves and executes the same way through `pkg/dbcopy/url.go`. In
`datatug-demo-projects`: one additional HTTP reference query that needs a key (fixture
snapshot recorded so the journey harness stays network-independent). Tests: a local
`httptest` server asserting the header is sent and never logged or echoed; `datatug
validate` still rejects a key committed inline. Non-goals: POST, GraphQL, HTTP steps in
playbooks (Phase 2–3 per the Feature's own non-goals).

### Task 9: `datatug context` CLI verb

**Id:** task-9
**Verifies:** agent-context#ac:context-verb-matches-web-panel, agent-context#ac:agent-context-blocked-by-policy
**Depends-On:** —
**Status:** planning

In `datatug-cli`: a `datatug context --value <Entity.Field>=<value> [--as principal]`
verb that calls the exact resolver endpoints the web context panel uses (`pkg/semantic`
relations and applicable-queries, executed through `pkg/secureread`) — no
agent-specific reimplementation (`agent-context` REQ:agent-readable-context-resource).
Every read goes through `accesspolicies.Run` under the given principal, so a `support`
run sees the same row/column limitations a human `support` session would
(`agent-context` REQ:agent-context-respects-access-control). Tests: fixture comparison
— `datatug context --value Customer.ID=5` against `admin` matches the web endpoint's
JSON shape; the same call under `support` shows the Canada-only filter and hidden
`Email`, matching Phase 1's `AC:hidden-column-refused` behaviour for a hand-crafted
request. Non-goals: an MCP resource (Phase 4, `serve-brokered-query-builder`'s scope).

### Task 10: Surfaces parity — hand-off, TUI Queries panel, command hygiene

**Id:** task-10
**Verifies:** surfaces#ac:handoff-preserves-selection, surfaces#ac:tui-browses-real-schema, executable-knowledge-library#ac:parameters-typed-consistently
**Depends-On:** —
**Status:** planning

In `datatug-cli`: extend `screenPathToWebPath` (`webui_handoff.go`) to cover every loop
screen — table, saved query, context — carrying project/table-or-query/selection state,
and add the reverse web→CLI hand-off (`surfaces` REQ:screen-level-handoff); bind the TUI
Queries panel to real project queries, running through `pkg/secureread` like every other
surface (`surfaces` REQ:tui-schema-browsing); delete or reimplement the unbound-flag
commands (`show`, `dataset*`, `datasets`, and the panicking `queries` stub,
`cmd_queries.go:16-19`) so the CLI's typed-parameter query behaviour stays consistent
between `query run`, `query list`, and whatever replaces the stubs
(`executable-knowledge-library` REQ:parameterized-queries-typed-parameters) — deleting
without a working replacement is an acceptable outcome per the audit's own
recommendation; the implementer records which it chose. Tests: a hand-off contract test
for every mapped screen, not only `"projects"`; a TUI test against the demo project's
real queries; `datatug --help` no longer lists a command that panics. Non-goals: the
Dashboards TUI panel (still a placeholder — no Phase 2 REQ names it); the MCP daemon.

### Task 11: Environment/time auto-narrowing of applicable queries

**Id:** task-11
**Verifies:** investigation-context#ac:environment-narrows-applicable-list
**Depends-On:** 1
**Status:** planning

In `datatug-cli`: the applicable-queries resolver (Phase 1's `pkg/semantic`) accepts the
context's active environment (task 1's `environment` semantic kind) and narrows the
applicable list to queries not tagged for a different environment, never silently — the
narrowing is visible and overridable, the same rule Phase 1 already applies to value
bindings (`investigation-context` REQ:auto-narrowing-by-environment-and-time). In
`datatug-apps`: the Investigation Context screen exposes an active-environment selector;
the applicable-queries panel visibly changes when it is switched. Tests: a
`staging`-only demo query drops out under `prod` and reappears when the environment
switches back. Non-goals: a time-window filter (the REQ scopes it to "environment and,
once environment/time semantic fields exist, a time window" — this task ships the
environment half; the time-window half is deferred, see Open Questions).

### Task 12: Phase 2 journey e2e

**Id:** task-12
**Verifies:** executable-knowledge-library#ac:capture-lands-as-git-diff, provenance#ac:result-shows-full-chain, meaning-inference#ac:accept-promotes-to-declared, heterogeneous-sources#ac:http-query-executes, agent-context#ac:context-verb-matches-web-panel, surfaces#ac:handoff-preserves-selection, investigation-context#ac:environment-narrows-applicable-list
**Depends-On:** 2, 5, 7, 8, 9, 10, 11
**Status:** planning

Extends Phase 1's Playwright harness (`datatug-apps`, task 3 of the Phase 1 plan) with
the journey this Plan opens with: capture a lookup as a team query → review and accept
an inferred mapping → ask a result "where did this come from" → run the key-guarded
HTTP query → switch the active environment and see the applicable list narrow →
Ctrl+W to the TUI and see the same query listed in its Queries panel. A separate,
non-Playwright test asserts `datatug context --value Customer.ID=5 --as support` against
the same demo project fixtures, matching the web panel. Tests: mechanism assertions
(server log entries, the captured file's Git status, the provenance chain's fields, the
narrowed list), not only outcome, per this plan's Approach. Update the Implementation
Matrix rows these twelve tasks touch.

## Deferred AC Coverage

Not every AC of every Draft Feature this Plan draws on belongs to Phase 2; ROADMAP.md's
own Included/Excluded lists and each Feature's capability-table Phase column decide
that, and this Plan follows them rather than covering every AC a Feature happens to
define. Deferred, with why:

- **`executable-knowledge-library`** — `semantic-parameter-declared`,
  `query-applicability`, `applicability-explainability` (Phase 1, already delivered);
  `investigation-playbooks`, `knowledge-promotion-lifecycle` (Phase 3, needs
  `sharing-scopes`); `templates-public-knowledge-packs` (Phase X). None are in
  ROADMAP.md's Phase 2 Included list.
- **`provenance`** — `access-refusal-carries-reason-code`: the REQ text itself says this
  "MUST" waits on the layered-ACL track's task 21 (Phase 4); today's `ACCESS_DENIED`
  response stays the interim form, unchanged by this Plan.
- **`meaning-inference`** — `deterministic-name-pattern-resolver`,
  `ambiguous-pattern-does-not-guess` (Phase 1, already delivered);
  `probabilistic-meaning-suggestions`, `suggestion-confidence-shown` (Phase 5
  experiment, explicitly gated on precision validation before any UI).
- **`heterogeneous-sources`** — `sql-sources-scan-and-query`,
  `sqlite-runs-in-release-binary` (Phase 0–1); `firestore-source`,
  `firestore-scan-works-query-rejects-cleanly` (Phase 3 per the capability table);
  `logs-source` (Phase 5).
- **`portable-query-model`** — every AC here is Phase 0–1 work already delivered
  (`dtql-is-source-of-truth`, `dtql-query-type-saves`, `sqlite-and-ingitdb-render`,
  `no-canonical-claim-for-stub-repos`); `native-query-labelled`'s HTTP half is what task
  8 extends, not a separate AC to re-verify.
- **`surfaces`** — `agent-connection-fixed`/`agent-path-resolves`,
  `hide-dead-pages`/`empty-pages-hidden` (Phase 0, already delivered);
  `hosted-web-ui-sneat-identity`/`hosted-demo-link-opens-the-project` (Phase 3);
  `mcp-serve-brokered-query-builder`/`mcp-face-depends-on-working-serve` (Phase 4, owned
  by the separate `serve-brokered-query-builder` Feature); `agent-skills-plugin`'s AC is
  a standing bar (8-of-9) this Plan does not regress but does not re-verify absent a new
  skill.
- **`agent-context`** — `skills-plugin-drives-cli` (standing bar, same reasoning as
  above); `link-serve-brokered-query-builder`/`link-ai-query-builder` are satisfied
  structurally by this Plan's tasks not restating either daemon's contract — no task
  needed.
- **`investigation-context`** — `context-basket-add-remove-disable`,
  `automatic-parameter-binding-from-context`, `no-hidden-filters` (Phase 1, already
  delivered); `saved-shareable-investigations` (Phase 3); `conflicting-values-and-propagation`
  is marked Phase 2 in the Feature's own capability table but is not named in
  ROADMAP.md's Phase 2 Included list — deferred here, see Open Question 1.
- **`semantic-layer`** — `entities-typed-fields`, `entity-table-mapping`,
  `field-column-mapping`, `semantic-tags-resolved-at-runtime` (Phase 1, already
  delivered).

## Risks

- **Correction to this stream's brief.** The brief that seeded this Plan named
  "datatug-core v0.19.0" as a landed precondition. The repository's newest tag is
  `v0.16.53` and `datatug-cli/go.mod` carries no `require` line for `datatug-core` at
  all (only a commented-out `replace`, matching Phase 1 task 1's still-open item to
  restore the module dependency or record the mirror). This Plan states preconditions
  by PR/package instead of a version number for that reason
  (`always-qualify-version-numbers` — there is no verified version to qualify here).
- **Phase 1's own plan tracker is stale.** `specscore plan info
  2026-09-09-phase-1-core-investigation-loop` reports 0 of 11 tasks complete, while the
  repositories show tasks 1, 4, 5, 11 and part of 6 already merged (PRs #302–304 in
  `datatug-core`; #192–193 in `datatug-cli`; the `dalgo2http` MVP; #8 in
  `datatug-demo-projects`). A `specscore plan reconcile` on that plan looks warranted so
  its record matches reality, but that plan is not this stream's to touch — flagged
  here, not acted on.
- **Concurrent writers to `EntityField.Mappings`.** Tasks 2 (capture), 3 (folder
  metadata) and 6 (mapping accept/correct) all touch entity/query files under
  `datatug-core`'s model. Mitigation: land task 1's model addition first per the
  Integration order; the landing owner batches same-repo branches rather than merging
  them independently.
- **HTTP auth headers are the demo's first real secret.** Mitigation: reuse Phase 1's
  `~/.datatug` store and `datatug validate` guard rather than opening a new credential
  path; task 8's test asserts the header is never logged or echoed.
- **`environment` as a semantic field kind could be confused with the existing
  `Environment`/`Server` infrastructure model** (`infrastructure-inventory`,
  `environment.go`). Mitigation: task 1 scopes the new `KnownTypes` entry to the
  existing project `Environment` concept `investigation-context`'s REQ already names,
  not a new per-row concept.
- **Fix-or-delete for `show`/`dataset*`/`datasets`/`queries`** (task 10) was left open
  by Phase 1's own Open Questions with a recommendation to defer to Phase 2; this Plan
  is that Phase 2 and picks it up, but the choice between "fix" and "delete" per command
  stays the task 10 implementer's, recorded in its PR, matching the Phase 1 pattern for
  the storage-indirection decision.

## Open Questions

1. `investigation-context`'s capability table marks "Conflicting values, propagation
   between related queries" as Phase 2, but ROADMAP.md's Phase 2 Included list does not
   name it, and it did not fit this stream's 8–12 task budget alongside the items
   ROADMAP.md does name. Recommendation: leave it for a follow-up task within Phase 2
   (before Phase 3) rather than pulling it into Phase 3 by default, since its own
   Feature already scopes it here — confirm.
2. Task 11 ships environment-based narrowing only; the REQ also names a time-window
   half once "environment/time semantic fields exist." Task 1 adds the `time` semantic
   kind as a type, but no task builds the time-window narrowing UI/logic itself.
   Recommendation: defer the time-window half to Phase 3 explicitly (update
   `investigation-context`'s capability table when this Plan is approved) rather than
   leaving it silently unscoped — confirm.
3. Task 10 folds "fix-or-delete the unbound-flag commands" together with hand-off and
   the TUI Queries panel because all three are `datatug-cli` surface-hygiene work with
   no dedicated Feature AC of their own beyond `parameters-typed-consistently`.
   Recommendation: keep them combined (matches Phase 1 task 2's precedent of bundling
   several client-contract fixes under one task) rather than splitting for the sake of
   a cleaner Verifies mapping — confirm, or split if the founder wants independent
   landing/rollback for command deletion versus hand-off.

---
*This document follows the https://specscore.md/plan-specification*
