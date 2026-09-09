---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Shared executable knowledge library

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/executable-knowledge-library?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Shared Executable Knowledge is the primary product capability: meaningful parameterized
queries become discoverable from context, explain their bindings and can be reused by
another investigator. Collections and human-guided scenarios compose those same assets.
The expected compounding value is a hypothesis, not established retention or a moat.

Queries, folders, semantic metadata and resolver code exist. The independent review
found mismatched browser/server contracts and incomplete journey acceptance. Phase 1
converges discovery and binding; Phase 2 adds project-file capture and useful provenance;
Phase 3 guides a bounded sequence; Phase 4 adds hosted project grants/promotion.

## Problem

Teams often repeat investigations but cannot readily find and reuse the query and
context that answered the previous question. DataTug should reduce that reconstruction
cost. Existing query libraries, notebooks and trusted AI query assets overlap this
capability; measure whether semantic discovery and explicit bindings improve second-user
reuse. Historical CLI panic/empty-tab findings are dated audit evidence, not proof of
current code state.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Parameterized queries with typed parameters | usable (model, CLI) / partial (web) | 1 | strategic |
| Semantic parameters (`requires Customer.ID`) | implemented-but-poorly-exposed | 1 | strategic |
| Query applicability ("useful right now") from selection/context | partial; acceptance pending | 1 | strategic |
| Explainability of a suggestion (selected column → maps to → query requires) | partial; acceptance pending | 1 | strategic |
| Query collections (folders with meaning) | usable (folders) | 2 | strategic |
| Capture from exploration (save the current query/lookup as a library asset) | planned | 2 | strategic |
| Investigation playbooks (entry points, branches, HTTP steps) | specified here only | 3 | strategic |
| Knowledge promotion lifecycle (personal → team → project → org) | planned | 4 | strategic |
| Templates / public knowledge packs | planned | X | strategic |

#### REQ: parameterized-queries-typed-parameters

Parameters MUST use declared types and identical server/CLI/browser validation.
Required/default/null semantics, conflicts, target selection and explicit overrides
follow the [core transport contract](../core-investigation-loop/api-contract.md).
A missing non-semantic required parameter also blocks Run. Values must actually bind
at execution, not merely appear in an echoed result header.

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

A user with project-write permission MUST be able to save an executed query/lookup as
an ordinary project query, with editable semantic parameters and no captured result rows
or secrets. The server validates authorization, path/ID and content, writes atomically,
and rejects a conflicting revision rather than overwriting another author's query.
On denial/failure the browser preserves the draft; it reports success only after a
persisted response. Capture defaults to a draft Git diff, not an implicit commit or
hosted team grant. Query author, purpose and source/binding provenance are inspectable.

#### REQ: investigation-playbooks

The Phase 3 slice MUST guide a person through named entry points and explicit ordered
or branching steps referencing existing queries, including configured HTTP queries.
Each step requires a deliberate action and has an observable next/terminal state.
No separate execution or ACL engine is introduced. Unattended conditional execution,
scheduling and retries are outside this slice; future automation requires a separate
plan and evidence, not a permanent product prohibition.

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
**When** the user clicks "Save as project query"
**Then** a new `<id>.query.json`/`<id>.query.dtql` pair appears as an uncommitted Git change with the semantic parameter already populated.

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

### AC: captured-query-explains-provenance (verifies REQ:capture-from-exploration, REQ:semantic-parameters)

**Given** a captured query with semantic parameters and a suggested mapping
**When** another authorized investigator inspects its source and accepts or corrects
that mapping, then reopens the project
**Then** source/query/capture provenance remains inspectable, declared/inferred origin
is accurate, and the persisted correction binds subsequent runs explicitly. Read-only
principals cannot modify it; conflict/failure retains the draft. Verify the detailed
provenance and meaning-inference ACs alongside this library outcome.

## Not Doing / Out of Scope

- A generic workflow engine in the initial playbook slice; future automation is not decided by this scope limit.
- A no-code app builder: the library holds queries and playbooks, not application screens.
- A data catalog or governance suite: applicability and capture make reuse fast, not "approved."
- Public knowledge packs mixed silently into an organization's private library.

## Open Questions

- Rich collection metadata beyond folders and README descriptions requires evidence;
  Phase 2 uses those existing primitives rather than inventing domain/owner tag schemas.
- Public knowledge packs and later automated playbooks require separately scoped plans.
Neither question blocks Phase 1 discovery or Phase 2 file-based capture.

---
*This document follows the https://specscore.md/feature-specification*
