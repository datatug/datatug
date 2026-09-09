---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Investigation Context

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Investigation Context is the typed, scoped set of facts an investigator carries
between queries. It gives Shared Executable Knowledge its inputs while keeping every
binding visible and overridable. This is a product hypothesis to test with reuse, not
proven competitive white space.

The initial empty Variables page is historical. The independent review found context
components and navigation on current main, with isolation and transport gaps. Phase 1
must complete add/remove/disable, typed conflicts and project/environment/principal
isolation. Time-window semantics and saved shared investigations are later extensions.

## Problem

An investigator repeatedly copies the same customer or invoice values between tools.
DataTug should carry meaningful facts without silently changing the question being run.
Other tools have substantial notebook variables and context reuse; the opportunity to
test is whether semantic meaning plus reusable cross-source queries reduces effort.
See the competitor review rather than claiming that every alternative stops at filters.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Context basket: add a semantic value from a cell, remove, temporarily disable | partial; acceptance pending | 1 | strategic |
| Automatic binding of semantic parameters from context and selection | partial; acceptance pending | 1 | strategic |
| Visible "why this query is constrained" and no hidden filters | partial; acceptance pending | 1 | strategic |
| Auto-narrowing compatible queries by environment/time | planned | 2 | strategic |
| Typed conflicts and scope isolation | partial; acceptance pending | 1 | strategic |
| Saved / shareable investigations and history | planned | 4 | strategic |

#### REQ: context-basket-add-remove-disable

The user MUST be able to add a semantic value to the Investigation Context from any cell
selection, remove it, and disable it temporarily without losing it. This is specified in full by
`core-investigation-loop` REQ:context-basket (an ordered set of `entity.field = value`, origin, and
enabled flag, visible on every project screen); this Feature owns the concept across every surface
the context basket must eventually reach, not only the web UI Phase 1 ships first.

#### REQ: automatic-parameter-binding-from-context

A query parameter whose semantic tag matches an active context value MUST be pre-filled when the
query page opens, distinct in origin from a value bound from the current selection. This is
specified in full by `core-investigation-loop` REQ:parameter-auto-binding.

#### REQ: no-hidden-filters

No query MUST be silently constrained by the Investigation Context; every binding applied MUST be
shown on the query page before the run and in the results header after it. This is specified in
full by `core-investigation-loop` REQ:no-hidden-filters and is the scope rule hypothesis H2 exists
to protect — the Investigation Context earns its name only if it is never mistaken for an invisible
filter bar.

#### REQ: auto-narrowing-by-environment-and-time

The context SHOULD carry an active environment and, once environment/time semantic fields exist
(P13), a time window, and SHOULD narrow which queries are shown as applicable to those that make
sense in that environment or window. Narrowing MUST remain visible and overridable, following the
same no-hidden-filters rule as value bindings.

#### REQ: conflicting-values-and-propagation

Conflict handling and store/project/environment/principal isolation MUST ship in
Phase 1. Follow the [normative binding contract](../core-investigation-loop/api-contract.md#binding-and-context-behavior)
for typed equality, precedence and explicit confirmation. No first-match heuristic or
late response may silently rebind an open query. Rich time-window semantics can follow
later; the security/context boundary cannot.

#### REQ: saved-shareable-investigations

A saved investigation MAY carry an explicitly selected safe subset of facts and query
references. Default links carry navigation identifiers only. Before values are shared,
the user MUST preview/approve them; credentials, inaccessible values and cached
results MUST NOT travel in the URL. The recipient independently authorizes and
re-executes references under their own principal. Sharing confers no execution rights.

## Acceptance Criteria

### AC: add-remove-disable-context-value (verifies REQ:context-basket-add-remove-disable)

**Given** `Customer.ID = 5` selected on the *Customer* table
**When** the user clicks *Add to context*, then later disables the chip
**Then** the value appears in the context list with its origin, and disabling it removes its effect on open queries without deleting it, matching `core-investigation-loop` AC:context-carries.

### AC: context-binds-a-different-query (verifies REQ:automatic-parameter-binding-from-context)

**Given** `Customer.ID = 5` is active in the context
**When** the user opens *Customer purchases by genre* from the queries list, not from a selection
**Then** its `Customer.ID` parameter is pre-filled and labelled `from context`, distinct from `from selection`.

### AC: bindings-always-shown (verifies REQ:no-hidden-filters)

**Given** an active context with two values
**When** any query runs
**Then** the results header lists exactly the bindings that were applied, and a query that used none of the context values shows none.

### AC: environment-narrows-applicable-list (verifies REQ:auto-narrowing-by-environment-and-time)

**Given** the context's active environment is `prod` and a query is tagged for `staging` only
**When** the applicable-queries panel renders
**Then** the `staging`-only query is not listed as applicable, and switching the context environment changes the list visibly.

### AC: conflicting-customer-values-flagged (verifies REQ:conflicting-values-and-propagation)

**Given** a query page already bound to `Customer.ID = 5` from context
**When** the user selects a different customer, `Customer.ID = 9`, on another open screen
**Then** the system shows a conflict rather than silently rebinding the open query, and propagation to that query is an explicit action.

### AC: shared-investigation-link-carries-state-not-secrets (verifies REQ:saved-shareable-investigations)

**Given** an investigation containing customer facts and restricted result data
**When** a navigation link is shared and opened by another user
**Then** the link contains no raw facts/results by default; optional explicit fact
sharing requires preview, and the recipient reauthorizes and re-executes rather than
inheriting the sender's access.

## Not Doing / Out of Scope

- A generic workflow engine: the context basket carries state a human is actively using; it does
  not automate a sequence of actions on its own.
- Hidden or hard-coded filters of any kind: this Feature's entire purpose is the opposite of a
  silent filter bar.
- Multi-user shared live context (two people editing the same context basket simultaneously) —
  sharing is link-based and asynchronous, not a real-time collaborative session.
- A BI/warehouse tool's saved-view concept: an investigation carries semantic values and query
  trail, not a stored SQL snapshot.

## Open Questions

None blocking Phase 1. Internal `basket` naming is an engineering choice; the public
label remains Investigation Context. Define saved-investigation retention and sharing
controls in this Feature before activating that later work.

---
*This document follows the https://specscore.md/feature-specification*
