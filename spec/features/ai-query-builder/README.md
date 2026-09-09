---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: AI Query Builder

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/ai-query-builder?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/ai-query-builder?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/ai-query-builder?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/ai-query-builder?op=request-change) |

**Status:** Approved
**Date:** 2026-06-04
**Owner:** alexander.trakhimenok
**Source Ideas:** ai-query-builder
**Supersedes:** —
**Grade:** A

## Summary

The AI Query Builder is a platform-agnostic **Capability**: an in-app conversational query builder where a user seeds a query in natural language and then progressively refines it through a chain of plain-English requests.

The *current query* is held as inspectable session state; each request is applied as a transformation against it. This Feature is a Capability per SpecScore's [capability-and-platform-implementations](https://specscore.md/feature-specification) pattern — it defines only behavior that must hold on every surface and owns the Implementation Matrix tracking the CLI and Web realizations.

## Problem

DataTug ships on multiple surfaces (a CLI and a Web UI) that drift apart. Building a query is inherently in-app — it is bound to live connections, the current result set, and a turn-by-turn feedback loop, so it cannot be delegated to an external AI agent the way metadata authoring can. Without a single platform-agnostic definition, each surface would re-invent the query-builder behavior and diverge silently. This Capability pins the cross-surface contract once: what "a current query refined by conversation" means, independent of how any one surface renders it.

## Behavior

### Session grounding

A builder session is the unit of interaction: one data source, one evolving query, many turns.

#### REQ: session-source-binding

A builder session MUST be bound to exactly one DataTug data source (a connection / dataset / environment) that is selected before the first request. Every query produced in the session is grounded in that source and its schema; the source does not change for the life of the session.

#### REQ: current-query-state

The session MUST maintain exactly one *current query* as structured, persistent state. The current query is created by the first request and carried across turns; each later request reads and updates this same state rather than starting over.

### Conversational refinement

#### REQ: nl-seed

The first natural-language request in a session MUST produce the current query, grounded in the session source's schema (the seed).

#### REQ: nl-transformation-delta

Each subsequent natural-language request MUST be applied as a transformation (a delta) to the current query — such as adding, removing, or modifying columns, filters, ordering, limits, or joins — while preserving the parts of the query the request does not address. The enumerated kinds are illustrative of what a delta may touch; which kinds a given surface guarantees is recorded per-surface in the Implementation Matrix, not mandated uniformly here.

#### REQ: query-inspectable

The current query's definition (not only its results) MUST be viewable by the user on any turn, so the user can see what the conversation has built.

### Execution

#### REQ: execution-mode-toggle

The builder MUST support two execution modes, **auto-run** and **manual-apply**, switchable by the user during a session. In auto-run, a changed current query is executed and its results presented without further user action. In manual-apply, a changed current query MUST NOT execute until the user issues an explicit user-initiated run action. The default mode is Implementation-defined and MUST be documented by each Implementation.

#### REQ: results-visible

When the current query is executed, its results MUST be presented to the user in a form the surface can render, and the executed query MUST be re-runnable without re-stating the conversation.

### Safety and failure modes

#### REQ: read-only-queries

The builder MUST produce only read-only queries. A request to mutate data or schema (insert/update/delete/DDL) MUST be refused, leaving the current query unchanged.

#### REQ: unresolvable-request

When a request cannot be applied to the current query — because it is ambiguous, references schema that does not exist in the session source, or asks for an operation the Implementation does not support — the builder MUST report the reason and leave the current query unchanged (no silent partial or incorrect edit).

### Cross-surface parity

#### REQ: parity-matrix

This Capability MUST own an `## Implementation Matrix` recording every platform Implementation's parity (`Platform`, `Status`, one-line `Brief`, `Link`), so that a reader of the Capability alone can see which surface is ahead and where each surface falls short, without opening the per-surface repositories.

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI | Partial | Single-source add-column/add-filter/order guaranteed; cross-source joins and computed columns not yet supported | specscore:feature/ai-query-builder@github.com/datatug/datatug-cli |
| Web | Planned | In-app query builder for the DataTug web UI; not yet started | specscore:feature/ai-query-builder@github.com/datatug/datatug-apps |

## Acceptance Criteria

### AC: source-bound-before-seed (verifies REQ:session-source-binding)

**Given** a new builder session whose data source has been selected
**When** the user issues the first request and any later request
**Then** every produced query is grounded in that one source's schema, and no turn switches the session to a different source.

### AC: state-persists-across-turns (verifies REQ:current-query-state)

**Given** a session with a current query built over two prior turns
**When** the user issues a third request
**Then** the third request updates the same current query, and the columns/filters from the prior turns that the request does not mention remain present.

### AC: first-request-seeds (verifies REQ:nl-seed)

**Given** an empty session bound to a source
**When** the user asks "show me the last 10 users who logged in"
**Then** a current query grounded in that source's schema is created and becomes the session's current query.

### AC: follow-up-is-delta (verifies REQ:nl-transformation-delta)

**Given** a current query selecting users with a filter applied
**When** the user says "add a column for the user's country" and then "order by signup date"
**Then** each request updates the same current query — the country column is added and the ordering is applied — while the existing selection and filter remain intact.

### AC: query-can-be-viewed (verifies REQ:query-inspectable)

**Given** a session with a non-empty current query
**When** the user asks to see the current query
**Then** the query's definition is presented, distinct from its result rows.

### AC: manual-mode-defers-execution (verifies REQ:execution-mode-toggle)

**Given** a session set to manual-apply mode
**When** a request changes the current query
**Then** the query is updated but NOT executed until the user issues an explicit run action.

### AC: auto-run-executes-on-change (verifies REQ:execution-mode-toggle)

**Given** a session set to auto-run mode
**When** a request changes the current query
**Then** the updated query is executed and its results are presented without any additional user action.

### AC: results-presented-and-rerunnable (verifies REQ:results-visible)

**Given** a current query that has been executed
**When** the results are returned
**Then** they are presented in a renderable form and the same query can be re-run without re-entering the conversation.

### AC: mutation-refused (verifies REQ:read-only-queries)

**Given** a session with a current query
**When** the user asks to "delete all inactive users"
**Then** the request is refused, no mutating query is produced, and the current query is unchanged.

### AC: unknown-column-reported (verifies REQ:unresolvable-request)

**Given** a current query whose source has no `nickname` column
**When** the user says "add a column for nickname"
**Then** the builder reports that the column cannot be resolved and leaves the current query unchanged.

### AC: matrix-tracks-surfaces (verifies REQ:parity-matrix)

**Given** the Capability README
**When** a reader inspects it
**Then** an `## Implementation Matrix` is present with one row per surface (CLI, Web), each carrying a Status, a one-line Brief, and a Link to that surface's Implementation.

## Rehearse Integration

All ten ACs have a concrete behavioral surface (session state, query transformation, execution gating, refusal paths, document structure) and are testable in principle. However, this Capability defines cross-surface *contract* behavior with no executable surface of its own — every AC is exercised against a concrete Implementation (the CLI Implementation in `datatug-cli`, where the Rehearse stubs live). Stub scaffolding is therefore deferred to the Implementation Features rather than authored here, so the stubs track real, runnable surfaces instead of an abstract Capability.

## Not Doing / Out of Scope

- Platform-specific UX (terminal keys, web checkboxes/buttons, rendering details) — those belong to each Implementation, not the Capability.
- Metadata/project authoring by chat — explicitly excluded in the source Idea; that need is served by external AI agents using DataTug skills.
- Mutating or DDL queries — read-only only (REQ:read-only-queries).
- Mandating the internal representation of the current query (SQL string vs query AST) — left to Implementations; the Capability requires only that the state be structured, persistent, and inspectable.
- Cross-source joins and computed columns as a guaranteed cross-surface behavior — surfaces may support them, but the Capability does not require them; their parity is tracked per-surface in the matrix.

## Assumption Carryover

From the source Idea `ai-query-builder`:

- **Progressive refinement is reliable (Must-be-true)** — encoded as REQ:nl-transformation-delta + REQ:current-query-state and validated by AC:follow-up-is-delta / AC:state-persists-across-turns against the CLI Implementation.
- **A current query can be represented platform-agnostically (Should-be-true)** — addressed by keeping the representation unspecified (Not Doing) while requiring it be structured/persistent/inspectable (REQ:current-query-state, REQ:query-inspectable).
- **Cross-repo `Implements:` lints clean (Must-be-true)** — exercised by the CLI Implementation's `**Implements:**` reference to this Capability and this Capability's matrix links; remains open until both repos lint together.

## Open Questions

- Should the execution-mode default be fixed by the Capability, or remain Implementation-defined (current choice)?
- When a request is partially applicable (one clause resolves, another does not), should the builder apply the resolvable part or reject the whole request? Currently REQ:unresolvable-request rejects the whole request; a future iteration may allow partial application with disclosure.
- How is the current query represented so it can be shared/serialized across surfaces if a session ever migrates between CLI and Web? Out of scope now; tracked for the Web Implementation.

---
*This document follows the https://specscore.md/feature-specification*
