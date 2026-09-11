---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Server-enforced access control

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/server-enforced-access?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

The server enforces project/source/row/column access through DALgo. Browsers and AI
agents consume authorized results and cannot elevate their identity. The independent
review found landed secure-read/server components and a stronger pending VM layered
ACL chain; integration and full browser negative acceptance remain incomplete.
Consume that work in dependency order rather than treating DataTug integration as a
future greenfield Phase 4 task.

## Problem

A visible limitation or client-side filter cannot secure data. The same enforcement
boundary must cover direct select, saved query, related lookup, counts and metadata
sampling. Native SQL that cannot be constrained must fail closed under a protected
profile. The prior label-only exception was inconsistent with that guarantee.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| DALgo policies: default-deny, hierarchical, row conditions, field allow-lists, principals | usable (external, Stable) | — | enabler |
| CLI `query run` enforcement | usable | — | strategic |
| `datatug serve` enforcement: every read passes `accesspolicies.Run`; principal from the serving user or an authenticated caller | partial; acceptance pending | 1 | strategic |
| Source-level allow/deny per principal | partial; acceptance pending | 1 | strategic |
| Native SQL under protected access | refusal required; privileged opaque grant separate | 1 | security |
| Layered multi-owner policies and Explain Access | pending VM stack; refresh delivery receipts | prerequisite when used; broader UI later | enabler |

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

Use the protected-versus-privileged rule in core-investigation-loop
REQ:opaque-sql-limitation and its transport appendix. Protected native SQL MUST be
refused before dispatch; source read alone does not grant opaque execution. Explicit
privileged read-only execution must disclose its distinct profile. Every legacy route
must obey this boundary, and the support demo must not have an opaque grant.

#### REQ: layered-acl-track

Consume the approved DALgo/OpenVaultDB layered-ACL contract and its existing provider,
SDK and DataTug work. Refresh the delivery packet and remote receipts before assuming
which tasks are pending. The review found DataTug UI/server implementation already on
VM branches, tested with local dependencies; it is not proof of a published chain.
Policy format, evaluator and Explain Access taxonomy stay upstream. Reconcile and
publish providers before release acceptance in consumers; no parallel evaluator.

#### REQ: recorded-evidence-rechecked

Recorded executions, result snapshots, annotations, and incident events, projections and watch
streams ([evidence-records](../evidence-records/README.md), [annotations](../annotations/README.md),
[incidents](../incidents/README.md) REQ:events-and-projection-policy-rechecked) MUST be read
through the same policy path as live data: the current principal's source/row/column policy is
re-evaluated on every read, hidden columns are never returned from a snapshot, a fact value in
an event whose column the principal may not read is redacted at read time, and a principal
without current access receives a non-disclosing denial. Historical evidence is never a
backdoor (vision §15).

#### REQ: incident-scoped-grants

NEXT (Track B2), not MVP: a temporary investigation grant MAY widen a principal's read access
for the duration of an assigned incident — named sources and collections, masked sensitive
fields, an expiry tied to the incident's resolution or a time limit — issued and revoked
through the server and recorded as incident events (vision §35). It composes with the layered
ACL contract; it MUST NOT bypass masking, and sensitive-field escalation MUST be narrow,
justified and temporary. Controlled production modifications with approval (vision §36) are
FUTURE and require their own Feature.

## Acceptance Criteria

### AC: snapshot-read-under-current-policy (verifies REQ:recorded-evidence-rechecked)

**Given** a snapshot recorded by `admin` that includes `Customer.Email`
**When** `support` reads it
**Then** `Email` is absent, the limitation names the policy, and a request for the raw snapshot bytes is refused without data.

### AC: incident-grant-expires (verifies REQ:incident-scoped-grants)

**Given** an investigation grant for `support` on `Invoice` tied to INC-1
**When** INC-1 is resolved
**Then** the next read under `support` is limited as before the grant, and the grant's issue and expiry are events on INC-1.

### AC: dalgo-access-is-the-only-evaluator (verifies REQ:dalgo-access-policies-external)

**Given** any DataTug repository's access-control code
**When** it is reviewed
**Then** it calls into `dal-go/dalgo/access`, with no parallel row/column evaluator found elsewhere in the codebase.

### AC: query-run-still-enforces (verifies REQ:cli-query-run-enforcement)

**Given** the `support` principal and a policy hiding `Customer.Email`
**When** `datatug query run --as support --role support` selects `Customer`
**Then** `Email` is absent from the result, matching the CLI's existing, live-verified behaviour.

### AC: serve-enforces-every-read (verifies REQ:serve-acl-all-reads)

**Given** `datatug serve --as support --role support`
**When** any read reaches the server, including a related-record lookup and an applicable-query sample
**Then** all of them pass through `accesspolicies.Run`, matching `core-investigation-loop` AC:restricted-rows-and-columns-server.

### AC: source-denied-entirely (verifies REQ:source-level-allow-deny)

**Given** a principal with no allow-listed access to the `support-notes` inGitDB source
**When** they request any record from it
**Then** the server refuses the whole source, not merely specific rows or columns within it.

### AC: native-sql-labelled-not-silently-filtered (verifies REQ:opaque-sql-limitation)

**Given** support uses the protected profile without opaque-query permission
**When** any native SQL route is requested
**Then** it is refused before execution, without data. An explicit privileged grant is
tested independently and is labeled as lacking row/column enforcement. The AC slug is
retained for traceability; its earlier label-only behavior is superseded.

### AC: layered-acl-task-21-is-the-integration-point (verifies REQ:layered-acl-track)

**Given** the pending layered provider/SDK/DataTug implementation
**When** it is integrated against current main and published dependency revisions
**Then** actual protected journeys pass using the upstream policy contract, and exact
remote receipts distinguish landed components from the complete dependency chain.

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
