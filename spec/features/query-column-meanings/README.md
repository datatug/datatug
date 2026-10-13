---
format: https://specscore.md/feature-specification
status: Stable
---

# Feature: Query column meanings

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-column-meanings?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-column-meanings?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-column-meanings?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-column-meanings?op=request-change) |
**Status:** Stable
**Source Ideas:** —
**Date:** 2026-10-13
**Owner:** DataTug
**Supersedes:** —

## Summary

Map query result columns to existing or newly created MeaningGraph entities and properties while DataTug creates the corresponding ModelSpec result structure. The pilot succeeds when adding or changing a meaning is easy and a reader can see and inspect the meaning directly in subsequent query results.

Authority: the user requested design, SpecScore specification, planning and implementation, then instructed: "Work autonomously. I’ll judge the implemented pilot." The acceptance priority is: "how flawlessly and hassle free user can add/change meanings and mappings and how it’s reused in query results (like how well it clear to user that column has meaning)". This authorizes implementation without another design approval checkpoint; it is not a claim that the user approved each implementation choice below.

## Problem

A column name and database type do not tell a reader what the values represent. Users should capture that knowledge where they encounter it, without leaving their query to author several metadata files or configure a registry. Meaning must survive saving and rerunning a query, and must not become misleading when the query changes.

This is the manual query-result authoring slice. The separate draft Semantic Advisor proposal (DataTug PR #132) concerns physical-schema inference and is neither adopted nor made a dependency by this feature.

## Behavior

### User journey

1. I run a query and open a column's visible menu using pointer or keyboard. The selected header remains identifiable while a compact mapping editor opens.
2. I choose an existing property or create an entity and property inline. I can read definitions before choosing. Model details are available without obscuring the main decision.
3. I save. The query and all its semantic metadata are persisted together; success closes the editor and the same grid displays a visible meaning badge without another click or rerun.
4. I click the badge to inspect or change the meaning. A failed save leaves my edits in place, with a useful error and a retry action.
5. I reopen the saved query and run it. Its mapped columns carry the same badges. If the query or returned shape changed, the previous meaning is visibly marked for review rather than silently confirmed.
6. I select a previously defined property from another personal query saved on this device without exporting a file first. Its source query is visible, and editing the copied definition leaves the source unchanged. I can also import MeaningGraph definitions or export MeaningGraph and ModelSpec documents for use outside DataTug.

### REQ: column-entry-and-feedback

The grid MUST offer a discoverable column-header action with an accessible name, usable without right-click and on narrow screens. Mapped columns MUST retain their original column name and show a secondary meaning label in the header. A label MUST open the mapping editor in one action and show both the property and its owning entity. Confirmed, unsaved and Needs review states MUST be visually and textually distinguishable, not color-only. No-data results with column metadata SHOULD remain mappable.

### REQ: inline-authoring

The editor MUST support selecting an existing entity/property, creating both inline, changing a binding, changing project-owned labels/definitions, and removing a binding. Existing properties MUST show definitions and scope. An optional context label distinguishes roles such as billing and shipping; it is not a new MeaningGraph binding-role enum. Duplicate ids and incompatible selections MUST produce field-level feedback. For this pilot, compatible means an existing property whose `of` resolves to the selected entity and whose target model field exists: a MeaningGraph property does not declare a scalar type, so DataTug MUST NOT invent a type-compatibility restriction. The user judges semantic suitability from its definition; each target ModelSpec field retains its own declared type. Cancel MUST leave saved state unchanged; removing a binding MUST retain reusable definitions. Explicit definition changes MUST update every use of that concept in the saved query.

### REQ: typed-result-bundle

The saved query MUST carry one versioned `semanticModel` bundle containing its MeaningGraph document, ModelSpec document, output bindings and the execution/schema basis used to confirm them. MeaningGraph owns entities and properties (`meaning/draft-2`); ModelSpec owns a `QueryResult` record's typed fields (`1.0-draft-2`). Every returned column MUST have a deterministic unique field, including repeated/invalid SQL column names. Unknown driver types MUST become `any`; result values MUST NOT be sampled to infer identity, requiredness, uniqueness or types, nor copied into metadata. Joined result rows MUST NOT automatically be declared instances of one semantic entity. Property bindings use the MeaningGraph `value` role.

### REQ: atomic-save-and-compatibility

Cloud, GitHub and local-agent saves MUST use the existing authenticated/capability-checked query API with revision conflict checks and operation replay. Servers MUST explicitly advertise `semanticModel` support; clients MUST refuse to send semantic metadata to older servers that could drop it. The full query and bundle MUST share one save boundary. A device-personal copy MAY provide the no-setup pilot path, but MUST persist the complete executable query and bundle together, be discoverable as a saved query, restore on a cold route, and be labeled "Saved on this device". It MUST NOT be a metadata overlay on a read-only query. Device concurrent updates MUST detect stale revisions. Retry MUST preserve the pending edits and must not create duplicate copies or definitions. Permission, quota or revision failures MUST never show Saved.

### REQ: reuse-and-review

Confirmed mappings MUST be restored after rerun, route reopen and full browser reload. A SHA-256 fingerprint of query type, exact text, connection id and federation configuration plus result index and ordered declared column names/types scopes each mapping. A changed basis MUST mark affected mappings Needs review and exclude them from confirmed labels until the user explicitly remaps or reconfirms against a newly executed result. While fingerprint validation is pending, the UI MUST show a neutral checking state instead of a false review warning or confirmed state; mapping mutations MUST wait for readiness. An unavailable basis MUST never be presented as confirmed. The pilot MAY conservatively invalidate all mappings after any query text edit, including aliases. A column index or alias alone MUST NOT establish identity. Previous definitions and proposed mappings MUST remain available during review. Existing query-owned meanings MUST be reusable. The property picker MUST offer meanings from other device-personal queries in the same project and identify their source query. Selecting one MUST copy its required definitions without its old field bindings or any mutation of the source query; conflicting ids MUST remain distinct. Loading saved choices MUST NOT block inline creation. Explicit import of supported exported documents MUST also remain available.

### REQ: validation-and-export

Both client and server MUST validate bounded bundle shape and references: unique concepts, properties belonging to existing entities, supported types, existing model fields, unique mapped output positions, consistent MeaningGraph field bindings and no unsupported lossy members. Existing credential screening MUST traverse newly persisted text. Exports MUST be valid standalone MeaningGraph and ModelSpec documents with resolvable sibling model references. The export UI MUST distinguish definition snapshots from public Registry publication and preserve standard format version names. Exporting or saving metadata MUST not export result rows.

### Pilot scope and design decisions

- One query owns its semantic bundle. Query-wide atomic storage reuses existing save semantics instead of introducing a multi-document project transaction.
- The easy trial path saves a complete device-personal query; existing server stores use their real API when upgraded. No direct client Firestore writes.
- The default graph and model belong to this query; shared public concepts and registries are not modified. A context role is descriptive metadata, not structural foreign-key inference.
- In scope: manual single-column mapping, inline authoring/edit/removal, reuse, visible meaning labels, durable save/reload, conservative invalidation, standard exports and an exercisable pilot.
- Deferred: automatic inference, source-wide mapping, public Registry search/publication, bulk selection, project-wide mutable catalogs, semantic row-instance import and a SQL lineage parser.

## Acceptance Criteria

### AC: visible-column-meaning (verifies REQ:column-entry-and-feedback)

**Given** a query result with a column called `amount`, **When** the user maps it to Invoice / Total, **Then** the grid immediately retains `amount` and shows Invoice / Total as a clickable header label; keyboard activation opens that mapping and the same actions are usable at a narrow viewport.

### AC: create-change-remove (verifies REQ:inline-authoring)

**Given** no matching meaning, **When** the user creates Invoice and Total in the mapping editor, saves, edits their definition, changes the column to another property, and removes the binding, **Then** every successful step updates the visible result without navigation away, the definitions remain available for reuse, and canceling an intermediate edit changes no saved state.

### AC: safe-result-shape (verifies REQ:typed-result-bundle)

**Given** duplicate column labels, an unknown driver type, a computed column and joined rows, **When** a mapping is saved, **Then** ModelSpec fields are unique and valid, the unknown type is `any`, no keys/required/unique constraints or row-instance claims are invented, and no result values occur in the saved metadata.

### AC: revisioned-store-roundtrip (verifies REQ:atomic-save-and-compatibility)

**Given** a supported server store and current query revision, **When** the query and semantic bundle are saved and read again, **Then** the full bundle is identical, stale writes leave the record unchanged, an identical retry returns the same operation result, and an older server lacking capability refuses the semantic save in the UI before any lossy request.

### AC: complete-device-copy (verifies REQ:atomic-save-and-compatibility)

**Given** a read-only demo query, **When** the user saves a personal copy and its first mapping, then cold-opens its saved route, **Then** both the executable query and bundle are restored, the copy is discoverable and labeled as device storage, and the original query is unchanged. A concurrent stale edit is refused without losing the editor draft.

### AC: save-failure-keeps-work (verifies REQ:atomic-save-and-compatibility)

**Given** an open mapping editor with unsaved definitions, **When** persistence fails or conflicts, **Then** the editor stays open with its values, an actionable error appears, no confirmed badge or success notice appears, and retry cannot create duplicate saved copies or orphan definitions. If the user leaves the query while a lookup or save is pending, its late response MUST NOT overwrite the destination query or navigate the user back.

### AC: rerun-and-change (verifies REQ:reuse-and-review)

**Given** a saved mapping, **When** the unchanged query is reopened and rerun, **Then** the meaning is restored without a false Needs review warning while validation is pending; **When** text, source, result index or declared output schema changes, **Then** its prior mapping is Needs review until a user confirms it against a fresh execution, including when the column label remains unchanged.

### AC: reuse-existing-meaning (verifies REQ:reuse-and-review)

**Given** a property created in a prior mapping, **When** the user maps another compatible column and selects that meaning directly in a separate query, **Then** the existing definition and its source query are selectable without retyping or exporting a file, its entity is visible, and changing the new query does not alter the source query. Importing an exported meaning remains an additional supported path.

### AC: interoperable-exports (verifies REQ:validation-and-export)

**Given** the Invoice Total mapping, **When** its documents are exported, **Then** the MeaningGraph checker and ModelSpec checker accept the exported documents together, the property binding resolves to a field in the result model, and neither output contains query rows.

### AC: invalid-bundle-rejected (verifies REQ:validation-and-export)

**Given** a malformed, oversized, duplicate or dangling mapping bundle, unsupported fields or credential-bearing metadata, **When** the client or server is asked to save it, **Then** validation rejects it before persistence and the prior query revision remains intact.

### AC: pilot-journey (verifies REQ:column-entry-and-feedback, REQ:inline-authoring, REQ:typed-result-bundle, REQ:atomic-save-and-compatibility, REQ:reuse-and-review, REQ:validation-and-export)

**Given** the pilot opened through the real DataTug app, **When** a user runs a query, creates a meaning, saves, reopens, reruns, changes the meaning, reuses it, reloads, and exports, **Then** the journey works without manually injecting storage or forcing UI change detection and its visible labels and persisted documents agree. Automated browser coverage records the same journey, plus focused failure and stale-result cases. Browser network fixtures may isolate provider availability but MUST NOT be reported as proof of a live cloud deployment. This browser journey is the Rehearse equivalent for the pilot.

## Verification evidence

- Core contract: [PR #366](https://github.com/datatug/datatug-core/pull/366), released as Core `v0.53.0` at `6a958361d5561e47bc597a7d215c3b9d6c039d9f`. Full suite coverage is 100.0%; all seven observed post-target checks passed.
- Server persistence: [backend PR #44](https://github.com/datatug/backend/pull/44), released as backend `v0.29.0` at `cc71dce177d18c39e80445d5960551537cc8e8b5`. All ten observed post-target checks passed, including Firestore emulator roundtrips. This is library and storage verification; deployment of a consuming cloud service is not claimed.
- Local API: [CLI PR #376](https://github.com/datatug/datatug-cli/pull/376) landed at `43e621b8c0239da056b71e68be49066bf738286e`, which is also the remote CLI `v0.76.0` tag. All eighteen observed post-target checks passed, including release artifacts. The adapter advertises semantic metadata support and preserves the bundle through query read/save.
- Export compatibility: ModelSpec CLI accepts the result document; the current MeaningGraph Core reference checker accepts the exported graph with the exported ModelSpec document supplied as its model. The inspected MeaningGraph CLI rejects `meaning/draft-2`; released CLI support and its ordinary file-loading path are not verified consumers of these pilot exports.
- A live, unmocked local DataTug browser run verified first-copy creation, retained query parameter and results, immediate confirmed label, one-click reopening and renaming. A second run directly selected Invoice / Invoice reference from an earlier personal query, saved a new copy without a file exchange, retained seven rows and Customer ID 1, and immediately displayed the meaning under InvoiceId. After reopening, its label was changed to Invoice identifier; rerunning the original query still showed Invoice reference. Automated persistence, conflict, validation and browser checks complement these receipts; this does not claim a deployed cloud-store save.
- UI implementation: [PR #300](https://github.com/datatug/datatug-apps/pull/300). The browser journey exercises keyboard and 390px entry, first-copy save, cold reload and rerun, reuse across columns, rename propagation, exports, malformed import rejection, changed-query review, stale removal protection, an actual device revision conflict preserving the old label and draft, and valid import into a separate saved query. GitHub metadata is fixture-backed in this automated test; query execution uses the actual SQLite path.
- The final local browser journey also verifies right-click entry and that the primary Save button remains fully visible at 390px after optional context is collapsed. ContextMenu and Shift+F10 use the same editor; ordinary typing remains available to the grid.
- UI revision `d17f581a52ea5b0be5a6c8ea6cb9fba85cf18f9e`: 3,197 main unit tests passed (three skipped), typecheck passed, focused lint reported zero findings, and production build passed. Four browser cases passed, including the complete mapping/reuse journey and three affected existing scenarios; the real-agent federated cleanup test passed separately. Deferred-response tests cover capability lookup, save failure, copy creation and leaving an Ionic-cached page. Independent review approved this revision with no unresolved findings.
- All four jobs in [source CI run 38456195695](https://github.com/datatug/datatug-apps/actions/runs/38456195695) passed on that exact revision, including the general browser suite, mobile WebKit and production-build demo handoff stages.
- Integration revision `26f36914e7bf53808585e1bac421951d4de8a95d` includes app main `023c18de75455803ea74e54096718c9440d9ceed`. Its production build and four targeted browser cases passed again, and independent review found no lost pilot behavior or unresolved integration findings.
- The first integrated CI attempt passed the mapping journey, unit suite, build, checkout journeys and real CLI journey. One existing join-query assertion timed out after five seconds while its external SQLite download was still incomplete (HTTP 200 headers received, body transfer unfinished in the trace). That focused test passed unchanged locally; the failed browser job was rerun without code or assertion changes.
- The unchanged retry passed all four required jobs in [integration CI run 38463321557](https://github.com/datatug/datatug-apps/actions/runs/38463321557). WB verified two stable observations, landed UI PR #300 as `b24983bdab73944af1d17276f01bfbfe9b7621c4`, and synchronized the clean canonical checkout to that exact remote main revision.

- Production build `0.1.0+586` reports exact merge `b24983bdab73944af1d17276f01bfbfe9b7621c4`; [Cloudflare deployment 38469133898](https://github.com/datatug/datatug-apps/actions/runs/38469133898) passed. An unmocked production browser journey created Invoice / Invoice number, renamed it to Invoice reference, cold-reloaded and reran seven rows, selected that property directly in another query, and changed its copied label to Invoice identifier. Rerunning the source still showed Invoice reference. Actual production exports passed ModelSpec CLI and the MeaningGraph Core reference checker with their sibling model.
- All five observed exact-main checks (four CI jobs and Cloudflare deployment) passed for `b24983bdab73944af1d17276f01bfbfe9b7621c4`, with two stable WB observations. [Main CI run 38469132733](https://github.com/datatug/datatug-apps/actions/runs/38469132733) is green.
- That live journey exposed a brief false Needs review label before asynchronous fingerprints settled. [Follow-up PR #303](https://github.com/datatug/datatug-apps/pull/303) distinguishes neutral validation-in-progress from a confirmed mismatch and retains mutation guards. Revision `1845a1b2d950561c7fe421730f4ad70b8b6102da` passed 3,198 unit tests (three skipped), focused regression tests, typecheck, lint and the complete browser journey; independent review found no unresolved findings.

- Follow-up [CI run 38470428491](https://github.com/datatug/datatug-apps/actions/runs/38470428491) passed all required jobs. PR #303 landed as `249ce6629c49cd15feb46959d9a887f6b5c63c5d`; [Cloudflare deployment 38474220923](https://github.com/datatug/datatug-apps/actions/runs/38474220923) passed, and the public build reports `0.1.0+587` at that exact revision. The previously saved production query cold-reloaded and reran seven rows with Invoice / Invoice reference, no review warning, and validation no longer busy.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
