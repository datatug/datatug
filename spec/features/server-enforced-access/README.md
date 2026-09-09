---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Server-enforced access control

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P10 is server-enforced access control: source, row and column policies decided by the DataTug server through DALgo's `access` package, with the browser never a trusted enforcement boundary.

It is the pillar that lets DataTug make a credible self-hosting claim distinct from browser-centric tools (hypothesis H8, "confirmed in DALgo, unproven in the web path"). Audited state on 2026-09-09: DALgo's `access` package and `access.SecureReadSession` are usable and Stable, and `datatug query run` is the one call site in the whole ecosystem that enforces them — but `datatug serve` and every `datatug-apps` write path enforce zero access policy today. A layered multi-owner ACL MVP is materially further along in design and DALgo-side implementation than shipped code alone would suggest, tracked in a separate, already-approved plan this Feature does not redesign. Phase 1 (`core-investigation-loop`) wires the CLI's proven enforcement pattern into `serve`; Phase 4 consumes the layered-ACL MVP once its DataTug-integration task lands. This is a Capability; per-repo Implementations (`datatug-cli`, `dal-go/dalgo`) link to it.

## Problem

A browser cannot be trusted to enforce a row or column restriction on itself — any client-side
filter is a suggestion, not a control. DataTug's architecture already has the right shape for a
real answer: DALgo's `access` package supports default-deny, hierarchical, row-condition and
field-allow-list policies bound to principals (`dal-go/dalgo/access/*`, Stable), and
`access.SecureReadSession` wraps any adapter without adapter-specific code
(`access/session.go:241-257`). The gap is where it is wired: `pkg/accesspolicies` and `query run`
are "the only call site in the entire ecosystem"
(`datatug-cli/apps/datatugapp/commands/cmd_query.go`) — `datatug serve`/`pkg/api` have "no auth/
token/ACL check found," and the legacy `pkg/sqlexecute` exec endpoints "enforce nothing" (current-
state audit, capability table, P10 rows). Every write path in `datatug-apps` is equally unenforced.
Opaque SQL text is a structural limitation, not a bug to fix inside Phase 1: DALgo cannot rewrite
text it does not parse, so row/column policy can only apply to structured (DTQL) queries — Phase
1's own assumption A2 states this explicitly rather than promising something the architecture
cannot deliver yet.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| DALgo policies: default-deny, hierarchical, row conditions, field allow-lists, principals | usable (external, Stable) | — | enabler |
| CLI `query run` enforcement | usable | — | strategic |
| `datatug serve` enforcement: every read passes `accesspolicies.Run`; principal from the serving user or an authenticated caller | planned | 1 | strategic |
| Source-level allow/deny per principal | planned | 1 | strategic |
| Opaque SQL text under ACL: source-level only, read-only session; row/column restrictions apply to structured (DTQL) queries | planned (explicit limitation) | 1 | hygiene |
| Layered multi-owner policies, structured blockers, Explain Access, policy UI | specified/documented + implementing on the VM branch (dalgo plan tasks 1–10 complete) | 4 | strategic |

#### REQ: dalgo-access-policies-external

DataTug MUST treat `dal-go/dalgo/access` as the single policy evaluator; no DataTug repository MUST
implement a parallel row/column/principal check. This capability is already Stable and external;
this REQ pins DataTug's obligation to consume it rather than reinvent it, matching P1's equivalent
rule for the query AST.

#### REQ: cli-query-run-enforcement

`datatug query run` MUST continue to run every read through `pkg/accesspolicies` and
`access.SecureReadSession`, as it does today. This is the reference implementation every other
surface's enforcement (this Feature's remaining REQs) is required to match, not replace.

#### REQ: serve-acl-all-reads

Every read the web UI can trigger through `datatug serve` MUST pass through `accesspolicies.Run`
with the project's policy set and a fixed serving principal. This is specified in full by
`core-investigation-loop` REQ:server-acl-all-reads and REQ:principal-selection; this Feature owns
it as the standing P10 guarantee that must keep holding as `serve` grows past the Phase 1 demo.

#### REQ: source-level-allow-deny

A principal MUST be allow/deny-scoped per source and collection, not only per row/column within an
already-accessible source. This is the policy-path-scope mechanism `core-investigation-loop`
REQ:server-acl-all-reads relies on; this Feature asserts it must extend to every source
`heterogeneous-sources` adds, not only SQLite and inGitDB.

#### REQ: opaque-sql-limitation

Opaque SQL text executed as-is MUST be restricted to source-level policy only, in a read-only
session, and MUST be labelled when a row or column rule exists for that source that it cannot
apply. This is Phase 1's explicit, named limitation — `core-investigation-loop` assumption A2 and
REQ:opaque-sql-limitation state it exactly this way — because DALgo cannot rewrite text it does not
parse; the layered-ACL design takes the same position, and this Feature does not propose a
different one.

#### REQ: layered-acl-track

Layered multi-owner policies, structured blockers, Explain Access and a policy-authoring UI MUST
follow the already-approved layered-ACL MVP plan (`dal-go/dalgo` `spec/plans/layered-acl-mvp.md`)
rather than being redesigned here. That plan records tasks 1–10 complete on the VM branch
`layered-acl-query` (policy loader, scoped field/collection masks, typed principals, execution
gates, mirrored into `dalgo2sql`, `dalgo2ingitdb` and `openvaultdb-go`), with task 21, "DataTug
query/write and Explain Access integration," as the point where this Feature's Phase 4 work begins.
This Feature's obligation is to consume that plan's contract when task 21 is scoped, not to define
a competing ACL design.

## Acceptance Criteria

### AC: dalgo-access-is-the-only-evaluator (verifies REQ:dalgo-access-policies-external)

**Given** any DataTug repository's access-control code
**When** it is reviewed
**Then** it calls into `dal-go/dalgo/access`, with no parallel row/column evaluator found elsewhere in the codebase.

### AC: query-run-still-enforces (verifies REQ:cli-query-run-enforcement)

**Given** the `support` principal and a policy hiding `Customer.Email`
**When** `datatug query run --as support` selects `Customer`
**Then** `Email` is absent from the result, matching the CLI's existing, live-verified behaviour.

### AC: serve-enforces-every-read (verifies REQ:serve-acl-all-reads)

**Given** `datatug serve --as support`
**When** any read reaches the server, including a related-record lookup and an applicable-query sample
**Then** all of them pass through `accesspolicies.Run`, matching `core-investigation-loop` AC:restricted-rows-and-columns-server.

### AC: source-denied-entirely (verifies REQ:source-level-allow-deny)

**Given** a principal with no allow-listed access to the `support-notes` inGitDB source
**When** they request any record from it
**Then** the server refuses the whole source, not merely specific rows or columns within it.

### AC: native-sql-labelled-not-silently-filtered (verifies REQ:opaque-sql-limitation)

**Given** the `support` principal and a row/column policy on `Customer`
**When** the SQL-text query `albums_by_title` runs
**Then** it executes read-only under source-level policy only and the result states "row/column policies not applied to native SQL," matching `core-investigation-loop` AC:native-sql-labelled — never a claim that row/column policy was applied when it was not.

### AC: layered-acl-task-21-is-the-integration-point (verifies REQ:layered-acl-track)

**Given** the layered-ACL MVP plan's task 21 begins
**When** DataTug's query/write and Explain Access integration is implemented
**Then** it consumes the plan's already-approved contract (policy format, typed principals, execution gates) rather than introducing a new one, and any deviation is raised against that plan, not decided unilaterally in this Feature.

## Not Doing / Out of Scope

- Client-direct database access for protected sources: every read and write reaches a source
  through the trusted server, never a direct browser or agent connection.
- Redesigning the layered-ACL policy format, typed-principal model, or Explain Access taxonomy —
  that work belongs to the dalgo plan this Feature consumes, not to this Feature.
- Multi-user authenticated principals in Phase 1: `serve` fixes one principal per session via
  `--as`; per-request authenticated principals are Phase 3/4.
- Row/column enforcement on opaque SQL text: an explicit, stated architectural limitation, not a
  gap this Feature promises to close without a parser.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
