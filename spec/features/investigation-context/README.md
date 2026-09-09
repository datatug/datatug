---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Investigation Context

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/investigation-context?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P6 is the Investigation Context: portable state (semantic values, environment, time window) that follows the user across screens and automatically binds compatible query parameters, always visibly and never as a hidden filter.

It is the thesis's second differentiating claim made concrete — "the context basket can be a first-class Investigation Context rather than a filter bar" (hypothesis H2, confirmed with a scope rule: it must never become a hidden filter). Audited state on 2026-09-09: not implemented. The web "Variables" page is an empty shell (`VariablesPageComponent`/`VariableService` — header plus empty `<ion-content>`, zero methods), and the audit found no dedicated context/basket page and no matching commits (`git log -S"basket"` returns nothing in `datatug-cli`); the closest data shape is the inert `Variables{Vars map[string] VarInfo}` model. Phase 1 (`core-investigation-loop`) builds the first real slice: add/remove/ disable a value, visible bindings, no hidden filters. Phase 2 adds environment/time auto-narrowing and conflict handling; Phase 3 adds saved/shareable investigations. This is a Capability; per-repo Implementations (`datatug-apps`, `datatug-core`) link to it.

## Problem

An investigation is rarely one query: the user picks a customer, checks their invoices, checks
their support notes, then runs a purchases-by-genre query — all about the same customer. Without a
carried context, the user re-enters `Customer.ID = 5` into every screen by hand, and there is no
way to see, at a glance, "everything I'm currently looking at is scoped to this customer, this
environment." Competitor research the thesis leans on found the closest analogues all stop short:
"Datadog saved views, Splunk workflow actions and Grafana variables all stop at 'filter'; none
carries semantic values across query assets" (product thesis, H2 evidence). DataTug's own web UI
has the placeholder for this — the Variables page — but it is UI-only with nothing behind it
(current-state audit, P6 capability rows), and no runtime concept of a context basket exists in any
audited repo today.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Context basket: add a semantic value from a cell, remove, temporarily disable | planned | 1 | strategic |
| Automatic binding of semantic parameters from context and selection | planned | 1 | strategic |
| Visible "why this query is constrained" and no hidden filters | planned | 1 | strategic |
| Auto-narrowing compatible queries by environment/time | planned | 2 | strategic |
| Conflicting values, propagation between related queries | planned | 2 | strategic |
| Saved / shareable investigations and history | planned | 3 | strategic |

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

When two related queries would receive incompatible bound values for the same semantic field (for
example a query already scoped to one customer receiving a second, different customer's value from
a fresh selection), the user MUST be shown the conflict rather than have one value silently
overwrite the other. Propagation of a new value to related open queries MUST be an explicit,
visible action, not an automatic side effect.

#### REQ: saved-shareable-investigations

A trail of queries, selections and pivots (an investigation) SHOULD be saveable and shareable, and
SHOULD be convertible into a playbook (P7) once playbooks exist. A shared investigation link MUST
carry project, screen and selection state, and MUST NOT carry credentials, matching the
epilogue-branch guarantee in `core-investigation-loop`'s user journeys.

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

**Given** a saved investigation covering two queries and a context value
**When** its link is opened by a different signed-in user
**Then** project, screen and selection state are restored and no credential travels in the URL, matching `core-investigation-loop`'s epilogue branch (b).

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

- Is the internal name `basket` kept for the Investigation Context, with only the user-facing label
  changed? Recommendation from the hierarchy: yes — matches `core-investigation-loop` assumption A3.

---
*This document follows the https://specscore.md/feature-specification*
