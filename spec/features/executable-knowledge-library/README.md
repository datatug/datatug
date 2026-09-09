---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Shared executable knowledge library

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P7 is the primary pillar of the thesis: a shared, executable knowledge library — parameterized queries with semantic parameters, query applicability ("useful right now"), collections, capture from exploration, investigation playbooks, and a promotion lifecycle — the explore → capture → share → reuse loop meant to make DataTug get more useful the longer a team uses it (hypothesis H10, "the compounding moat").

Audited state on 2026-09-09: queries, typed parameters and folders exist and are Git files (usable); `ParameterDef.Meta` semantic tagging exists but nothing reads it at runtime; applicability, explainability, capture-from-exploration, playbooks and promotion do not exist at all. Phase 1 (`core-investigation-loop`) builds applicability, its resolution chain, and parameter auto-binding; Phase 2 adds capture and richer collections; Phase 3 adds playbooks and promotion (needs P8's sharing scopes). This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`, `datatug-apps`) link to it.

## Problem

Most tools help a user ask an arbitrary question; almost none help an organization accumulate and
resurface its proven answers to recurring questions. DataTug's differentiating bet (hypothesis H1,
"confirmed — no product researched combines P3–P7") depends on this pillar actually compounding: a
query someone wrote last month should surface itself again when today's investigation matches its
shape, not sit unfound in a folder. Today the `QueryDef`/folder model is real and Git-native
(`datatug-core/pkg/datatug/query.go`; `datatug-demo-projects/demo-project-1/queries/**`), and the
web UI's Active/Bookmarked/Personal/Shared tabs render live — but the audit found zero applicability
logic and confirmed the tabs shell has "nothing behind it" because the agent connection is broken
(current-state audit, "Journey to first useful result" table, "First shared-query reuse" row).
`datatug queries` itself is a stub that panics (`cmd_queries.go:16-19`, reproduced live). Playbooks
and promotion have no code at all.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Parameterized queries with typed parameters | usable (model, CLI) / partial (web) | 1 | strategic |
| Semantic parameters (`requires Customer.ID`) | implemented-but-poorly-exposed | 1 | strategic |
| Query applicability ("useful right now") from selection/context | planned | 1 | strategic |
| Explainability of a suggestion (selected column → maps to → query requires) | planned | 1 | strategic |
| Query collections (folders with meaning) | usable (folders) | 2 | strategic |
| Capture from exploration (save the current query/lookup as a library asset) | planned | 2 | strategic |
| Investigation playbooks (entry points, branches, HTTP steps) | specified here only | 3 | strategic |
| Knowledge promotion lifecycle (personal → team → project → org) | planned | 3 | strategic |
| Templates / public knowledge packs | planned | X | strategic |

#### REQ: parameterized-queries-typed-parameters

A query MUST be able to declare typed parameters, editable and runnable from the CLI and the web UI
consistently. Already usable in the model and CLI, partial in the web (`input-parameters`,
`ParameterLookupService`); this REQ pins cross-surface consistency as a standing guarantee.

#### REQ: semantic-parameters

A parameter MUST be able to declare that it `requires` a semantic field (`ParameterDef.Meta` as an
`EntityFieldRef`), not just a name and type. This feeds `core-investigation-loop`
REQ:parameter-auto-binding and REQ:applicable-queries; this Feature owns `Meta` as a first-class,
inspectable declaration on every library parameter.

#### REQ: query-applicability

Given a set of active semantic values, the library MUST say which queries are "useful right now"
versus "not yet applicable" and why. Specified in full by `core-investigation-loop`
REQ:applicable-queries; this Feature asserts it as a standing capability that must keep working as
the library grows past Phase 1's three demo queries.

#### REQ: applicability-explainability

Every applicable or not-yet-applicable entry MUST carry the resolution chain from the selected
value to the requirement it satisfies or misses (`column → maps to field → query requires`). This
is the same chain REQ:applicable-queries already requires — explainability is applicability's
contract for showing its own reasoning, not a separate mechanism.

#### REQ: query-collections-folders

Queries MUST be organizable into recursive folders, and a folder SHOULD carry more meaning than its
name — owner, domain, description — once proven necessary. Folders are already usable
(`QueriesFolder` recursion); no new metadata is mandated until the Open Question below resolves.

#### REQ: capture-from-exploration

A user MUST be able to save the query or lookup they just ran as a library asset directly from the
result screen. Captured queries MUST land as ordinary project files (a visible Git diff) with
semantic parameters and applicability already populated from the context that produced them.

#### REQ: investigation-playbooks

A playbook MUST represent a human-guided investigation path with named entry points and ordered or
branching steps, including HTTP steps alongside query steps. It MUST NOT introduce automation
semantics (unattended conditional execution, scheduled runs, retries) — it guides a person through
steps they still choose, distinguishing it from a generic workflow engine.

#### REQ: knowledge-promotion-lifecycle

A query, collection or playbook MUST be promotable from personal → team → project → organization
scope without duplicating the underlying file. Promotion depends on `sharing-scopes` (P8) for the
scope model; this Feature owns the promotion action and its effect on discoverability.

#### REQ: templates-public-knowledge-packs

DataTug MAY offer public, non-organization-specific knowledge packs as a separate, opt-in layer
from an organization's private library. Public packs MUST NOT be mixed into private knowledge
without an explicit import step.

## Acceptance Criteria

### AC: parameters-typed-consistently (verifies REQ:parameterized-queries-typed-parameters)

**Given** `customer-invoices` with a typed `Customer.ID` parameter
**When** it is run from `datatug query run --var` and from the web query page
**Then** both accept the same type and reject the same malformed input.

### AC: semantic-parameter-declared (verifies REQ:semantic-parameters)

**Given** `customer-invoices`'s parameter definition
**When** it is inspected
**Then** `Meta` names `Customer.ID` as an `EntityFieldRef`, matching `core-investigation-loop` AC:applicable-with-chain-server.

### AC: applicability-updates-with-library-growth (verifies REQ:query-applicability)

**Given** a fourth query is added to the demo project requiring `Invoice.ID`
**When** `Customer.ID = 5` is selected (no `Invoice.ID` known yet)
**Then** the new query appears under "not yet applicable — needs Invoice.ID", alongside the existing `invoice-lines` entry, matching `core-investigation-loop` AC:applicable-with-chain.

### AC: chain-explains-every-entry (verifies REQ:applicability-explainability)

**Given** the applicable-queries panel for `Customer.ID = 5`
**When** each entry is inspected
**Then** every applicable and not-yet-applicable entry carries a resolution chain, none is unexplained.

### AC: folder-recursion-works (verifies REQ:query-collections-folders)

**Given** a nested folder `customer-support/escalations`
**When** the queries page renders
**Then** queries in the nested folder are reachable and the folder tree matches the on-disk layout.

### AC: capture-lands-as-git-diff (verifies REQ:capture-from-exploration)

**Given** an ad-hoc related-record lookup the user just ran for `Customer.ID = 5`
**When** the user clicks "Save as team query"
**Then** a new `.query.json`/`.dtql.yaml` pair appears as an uncommitted Git change with the semantic parameter already populated.

### AC: playbook-guides-without-automating (verifies REQ:investigation-playbooks)

**Given** a "Checkout failure investigation" playbook with three branching steps
**When** a user runs it
**Then** each step requires the user's action to advance, and no step executes unattended on a schedule.

### AC: promotion-changes-visibility-not-content (verifies REQ:knowledge-promotion-lifecycle)

**Given** a personal query promoted to team scope
**When** a teammate in the same Sneat space browses the library
**Then** the query becomes visible to them with identical content, and the underlying file is not duplicated.

### AC: public-pack-import-is-explicit (verifies REQ:templates-public-knowledge-packs)

**Given** a public "customer-support starter pack"
**When** an organization imports it
**Then** the imported queries land as a distinct, labelled set, never silently merged into the organization's own private library entries.

## Not Doing / Out of Scope

- A generic workflow engine: playbooks guide a human through chosen steps; they never execute
  unattended, on a schedule, or with conditional branching a human did not choose.
- A no-code app builder: the library holds queries and playbooks, not application screens.
- A data catalog or governance suite: applicability and capture make reuse fast, not "approved."
- Public knowledge packs mixed silently into an organization's private library.

## Open Questions

- Should `queries/`-level collections gain semantics beyond folders (owner, domain, entry points),
  or is a folder plus a README the simplest useful model? Recommendation: keep folders simple until
  Phase 2 proves the need for more.

---
*This document follows the https://specscore.md/feature-specification*
