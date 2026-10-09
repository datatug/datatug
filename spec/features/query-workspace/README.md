---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: DataTug Query Workspace

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-workspace?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-workspace?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-workspace?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/query-workspace?op=request-change) |

**Status:** Approved
**Source Ideas:** —
**Date:** 2026-10-09
**Owner:** DataTug
**Supersedes:** —

## Summary

One browser workspace lets analysts and developers author a saved query, preview its execution, run it on authorized sources, investigate related data, compare results, and place a governed reference on a dashboard. Its shared shell has Author, Compare, and Explore layouts; Author has Compose and Code views. The user approved this Feature and its implementation plan on 2026-10-09, including the join-authoring and projection rules below. The accompanying prototype establishes the interaction direction. Current demo/browser receipts and prerequisites are reconciled during planning; screenshots or this specification alone do not prove the end-to-end runtime journey.

## Problem

The observed query page makes a long text editor and stacked setup cards dominate the first viewport. Target, bindings, Run, and results can be out of view. Analysts need a legible source-to-result path; developers need code, emitted requests, and precise execution evidence. A unified workspace must serve both without losing state, silently running work, implying unsupported dialects, or hiding policy limits.

## Behavior

This Feature is the browser integration seam for the owning DataTug [Compare](https://github.com/datatug/datatug/blob/main/spec/features/compare/README.md), [core investigation loop](https://github.com/datatug/datatug/blob/main/spec/features/core-investigation-loop/README.md), Investigation Context, query storage, dashboards, and served query-builder contracts. Those artifacts retain authority for execution, comparison, policy, and storage. Reconcile their contracts before implementation; this approved design does not create a second comparer or ACL engine.

### REQ: shared-workspace-state

Author, Compare, and Explore MUST preserve one query draft, authorized target, bindings with origins, result references, provenance, selection, and unsaved state. Compose and Code are views of that draft. Switching views, targets, or layouts MUST NOT execute, discard edits, or relabel old results. The first viewport shows query intent/code, target, binding summary, Run, and results. Advanced controls expand on request. At 1280 px an AI/context drawer overlays rather than permanently narrowing the workspace.

### REQ: authoring-language-boundary

Existing YAML DTQL and supported native SQL remain distinct editable representations with original text preserved. Compose edits only a semantic subset it can round-trip; unsupported syntax remains in Code with an explanation. An SQL or GraphQL import MUST preserve the original text as a draft. A parser MUST accept only complete constructs in its supported profile; for an accepted parse it MUST diagnose every clause, argument, directive, or semantic loss that is not represented. An unsupported or unparsed statement MUST be rejected as a whole with the source preserved; partial conversion MUST NOT produce an executable query or grant native execution. Failed conversion MUST NOT bypass policy. GraphQL is an optional native, schema-backed mode only after resolver, field, authorization, and execution capabilities are established. Adoption of a FROM-first grammar remains spike-gated; illustrative syntax is not executable DTQL. Any adopted FROM-first grammar MUST permit omitting SELECT and MUST place optional SELECT last, as approved for this authoring experience. The spike MUST establish that behavior or record an infeasibility for an explicit product decision before adoption; it cannot silently move SELECT into another position.

The user-approved forms for new specs and UI are TugQL (pronounced “tug cue ell”), a readable FROM-first query syntax with optional SELECT last, and TugQTree (pronounced “tug cue tree”), structured YAML for the same query model. The UI labels are `TugQL · Query` and `TugQTree · YAML`. Existing YAML remains compatible; this decision requires no file migration and does not change delivered DTQL behavior. DTQL remains the existing library and website name. No umbrella rename, retirement, or migration is authorized; consider that only in a separate brainstorming session. If implemented, Go and TypeScript MUST implement the same versioned TugQL grammar and semantic contract. Both implementations MUST run the same checked-in golden fixture corpus: source text, expected canonical TugQTree YAML where serialization applies, structured diagnostic fixtures, and optional schema snapshots for relationship resolution. Both CI lanes MUST execute the identical corpus, extending DALgo's existing byte-identical Go/TypeScript fixture convention rather than creating a new engine or fixture convention. Syntax parsing MUST be separate from schema-dependent resolution: omitted ON or relationship shorthand may resolve to an executable explicit predicate only with an authorized declared schema; the parser MUST NOT guess relationships. The TypeScript implementation MUST use DALgo-JS's existing structured-query model and DTQL YAML parsing/serialization contracts; Go MUST lower to DALgo's existing structured query model. These are two authoring/parser adapters over one query model, not separate execution or policy engines. Diagnostic spans use 1-based line and column positions, with columns measured in Unicode scalar values, and an exclusive end position; Go and TypeScript MUST emit the same spans. Distinct representation identifiers and version policy for TugQL and TugQTree remain an implementation task; this Feature does not claim those types or version contracts already exist.

### REQ: optional-select-projection

If an adopted authoring language permits omitted SELECT, a single-source ungrouped query resolves only authorized source fields and displays that projection before Run. A grouped query defaults to its grouping keys only and never invents aggregates. A joined query without SELECT resolves an authorized projection across its sources while retaining qualified field identity; it does not require an explicit SELECT solely because a join exists. The projection's column set and order MUST be fixed before Run and MUST NOT vary with returned row values. Before execution, the authoring resolver MUST lower an omitted-SELECT join to an explicit ordered projection of qualified fields with unique output aliases; it MUST NOT pass an empty Columns list to an existing join executor. Source-field identity and lineage MUST remain separate from output aliases. Joined key fields MAY form one output column, attributed to both source fields, only when declared join comparison semantics guarantee lossless exact typed-value equality (including type and representation); its representative is the first present source value in authored join order. If collation, coercion, or unknown comparison semantics can match different representations, both qualified fields MUST remain separate output columns for the whole result, even when some rows happen to be byte-equal. A lossy merged value is forbidden, and qualified selection of either source field remains available. Outer-join absence remains distinct from a present null value and is preserved per source. In an adopted FROM-first grammar, SELECT is optional and, when present, MUST appear last. Alias binding and execution order MUST be specified; ambiguity is an error and textual order MUST NOT change grouping, sorting, or LIMIT semantics.

### REQ: source-and-join-discovery

The context rail searches authorized tables/views and saved queries. Join suggestions show key columns, direction, cardinality or unknown, provenance (foreign key, declared semantic mapping, inferred), and available cross-source cost/limit. Name similarity alone grants no authority. A declared foreign key is eligible for automatic ON completion only when its complete mapping and target key are uniquely resolved across all authorized source aliases already present in the query. Applying a suggestion edits the draft but does not run. Every field selected or referenced in a join, filter, group, order, or expression MUST be authorized and visible before it enters the draft or preview, and authorization MUST be revalidated at Run. Unauthorized metadata is not revealed.

### REQ: relationship-backed-join-authoring

When the author writes `FROM Invoice JOIN Customer` without ON, or writes `JOIN Customer ON`, Compose MAY autocomplete a predicate only when the authorized schema declares exactly one complete foreign-key relationship between the new target and all eligible source aliases already present in the query, to the referenced primary or unique key. An unqualified `current source` means this full set of prior source aliases, not only the most recently joined table; if more than one relationship qualifies, the shorthand is ambiguous. A qualified alias narrows resolution only when one complete relationship remains. This does not require source and target columns to share a name; for example, `Invoice.CustomerId -> Customer.Id` resolves to `Invoice.CustomerId = Customer.Id`. On Run, the editor MUST first insert and display the fully qualified, fully mapped predicate in the draft, then execute only if that exact relationship remains unambiguous and authorized. If zero or multiple relationships qualify, Run MUST stop with candidates or a diagnostic and MUST NOT guess. Candidate ordering and generated predicates MUST be deterministic. `JOIN Customer ON CustomerId` is shorthand only when `CustomerId` resolves to one declared relationship from the current source to a target primary or unique key; a qualified field or explicit complete predicate is required when it is ambiguous. A composite foreign key is one relationship and MUST expand to every ordered field pair; a partial component MUST NOT be treated as the whole relationship. The draft and saved query retain the relationship identity, ordered field pairs, source-field references, and schema version so edits, reopen, receipts, and provenance remain attributable. Before execution, the saved binding MUST be revalidated against the authorized current schema. If the relationship is absent, changed, or ambiguous, Run MUST stop and request review; it MUST NOT silently substitute another relationship.

### REQ: preview-and-execution-receipt

A provisional compiler/adapter preview is distinct from an immutable executed-request receipt. Preview binds draft revision, principal/policy context, target mapping, dialect, schema version, and parameter types/placeholders; any change marks it stale. Run revalidates/recompiles on the same execution path. A receipt shows actual emitted SQL or native source requests, local federation stages, target, binding metadata, and execution identity without credentials. A preview is never labelled actual SQL. A native non-SQL path shows why SQL is unavailable. Database EXPLAIN and EXPLAIN ANALYZE are separate capability-gated actions whose possible execution effects are disclosed.

### REQ: explicit-run-and-results

Run is explicit and names target, revision, and bindings. Missing, ambiguous, or conflicting required bindings block it near the field. Running, cancellation, error, empty, partial, and successful states are distinct. Editing or switching target retains the last committed result as Previous run with its original provenance. Late responses cannot attach to a newer draft or target; one target's failure cannot erase another's success. Headers show query revision, target, time, completeness, limitations, live/recorded status, applied nonsecret binding values and origins, row count when known, and timings. Secret bindings show only type and set/unset state; their values never enter URLs, telemetry, persisted drafts, result receipts, dashboard cards, or browser result cache. Nonsecret ad hoc values remain session scoped unless the user explicitly saves a query default through an authorized path.

### REQ: related-query-exploration

Selecting a cell can offer authorized related tables, records, and applicable saved queries through the existing semantic mappings and Investigation Context. Look up this value previews the saved query, typed binding and origin, target, and live/snapshot choice before execution; selection alone never runs. Use as data is a separate composition operation subject to output schema, access, boundedness, and cycle checks. The governed read path remains authoritative; name matches alone are not relationships.

### REQ: multi-target-results

The user can explicitly run one frozen query revision and visible binding set on several authorized compatible environments. Each target has independent progress, budget, cancellation, error, receipt, and result status. Results are retained independently. The UI discloses unavailable/setup-required targets and makes no atomic cross-database snapshot or identical-observation-time claim.

### REQ: honest-comparison

Compare consumes authorized result references and extends the existing typed recordset comparer where necessary. Current Core Options do not provide numeric tolerance or meaning/unit compatibility, and current side receipts do not carry completeness/expiry metadata; these are explicit contract extensions before the UI may offer those guarantees. Default comparison uses exact equality for types the current comparer supports. A missing/expired/unauthorized side, null or non-unique selected key, unsupported comparator, or incompatible schema or declared units blocks a keyed diff. Unknown meaning/units MAY permit a raw typed-value diff only when labelled as such, without a semantic-equivalence claim. Under the current Compare contract, partial/truncated results fail closed and cannot produce a CompareResult. Query Workspace MUST preserve that behavior until the owning Compare Feature/API contract is explicitly amended and that amendment is implemented. Only after that gate may an explicitly requested partial sample diff be offered. For such an authorized partial sample **or unequal policy visibility**, one-sided keys are labelled visible only in baseline/candidate, never Added or Removed; definitive removal/completeness totals are suppressed. Unequal policy visibility is prominently labelled and cannot support a causal or deletion assertion. The UI reports per-side target, timestamp, policy and completeness limits, row counts when known, and only differences justified by those limits. Reading persisted results rechecks access. A difference is not a causal conclusion.

### REQ: ai-action-proposal

AI uses the same typed authorized command handlers as manual controls and shows operation scope. An edit request produces a reversible diff and never runs by itself. An explicit request to run may execute after showing revision, target, bindings, and effect, without an extra approval tap solely because it came through AI. Opening the drawer alone sends no result cells to an external model. Applied edits mark the draft unsaved.

### REQ: dashboard-reference

Add to dashboard previews board/card destination and stores a reference to a saved query with revision, parameter, target/environment, visualization, and refresh policies. Public demo browsing and drafting MAY work without sign-in; committing a saved query or shared dashboard card requires a real signed-in principal with project write permission. At that submit point, the UI preserves the nonsecret draft across sign-in and resumes the intended save without silently publishing or using anonymous authentication. Each viewer's query/data access is rechecked; saving a card copies no credentials and grants no access. The card shows freshness and snapshot source.

### REQ: demo-and-accessibility

The maintained datatug-demo-project provides a prerequisite-labelled browser-runnable first example and a path through filter, preview, Run, invoice lookup, Compare, and dashboard card. Deterministic synthetic same-schema comparison fixtures include typed added/changed rows plus duplicate-key, partial, missing-field, and permission-limited cases. Public canonical DemoDB data is not altered to manufacture a difference. Controls have labels, visible focus, keyboard operation, non-color-only differences, and usable layouts at 1280/1440 px and 200% zoom.

## Acceptance Criteria

### AC: state-survives-layouts (verifies REQ:shared-workspace-state)

**Given** an unsaved customer query, a selected target, a CustomerId binding from a selected cell, and a previous result
**When** the user moves Author → Compare → Explore and Compose → Code, then changes target without running
**Then** draft, binding value and origin, result reference, and selection remain intact; the old result keeps its original target and Previous run label; no execution occurs; Run and results remain discoverable at 1280 px.

### AC: safe-authoring-round-trip (verifies REQ:authoring-language-boundary)

**Given** supported YAML DTQL, native SQL with an unsupported vendor function, and GraphQL with unmapped arguments
**When** each is opened or imported through Compose and Code
**Then** supported semantics round-trip; every original source remains a draft even when Run is denied; diagnostics identify loss; Compose refuses destructive conversion; GraphQL execution appears only for an authorized native source.

### AC: authoring-parity-fixtures (verifies REQ:authoring-language-boundary)

**Given** the versioned TugQL profile, its matching TugQTree serialization contract, and one checked-in corpus containing valid syntax, malformed syntax, unsupported constructs, diagnostic spans, quoted identifiers, aliases, omitted SELECT, final SELECT, composite and ambiguous relationships, and optional schema snapshots
**When** both Go and TypeScript implementations parse and resolve every corpus case in their respective CI lanes
**Then** both produce the same accepted/rejected outcome, canonical structured diagnostics and Unicode source coordinates, resolved semantic model, and expected canonical TugQTree YAML where specified. Syntax-only cases do not infer schema relationships; schema-dependent omitted-ON and shorthand fixtures resolve only from the supplied declared schema and emit the full explicit predicate. Changing an input, expected YAML, diagnostic, or schema fixture requires the same corpus to fail in both lanes until both implementations agree. Existing YAML remains compatible without migration, and the UI labels are exactly `TugQL · Query` and `TugQTree · YAML`.

### AC: projection-is-explicit (verifies REQ:optional-select-projection)

**Given** authorized Invoice fields, one restricted field, a grouped query, and `Invoice.CustomerId -> Customer.Id` where case-insensitive collation or coercion may join distinct string representations
**When** the query omits SELECT in an adopted authoring language that supports omission, with `LEFT JOIN` and `WHERE i.CustomerId IS NULL` also available as an explicit query
**Then** single-source output lists only authorized fields, grouped output lists grouping keys without invented aggregates, and each joined projection's output columns and order are fixed before Run. The FK columns appear once, attributed to both sources and represented by the first present source in authored join order, only when declared comparison semantics guarantee exact typed equality; when collation or coercion can match different representations, both qualified columns stay in the schema for every row. A manual qualified projection such as SELECT i.CustomerId, c.Id remains available and emits both fields separately. After a left join, `i.CustomerId IS NULL` still addresses the left source field, and the result preserves right-side absence separately from a present null. The adapter receives explicit ordered projections with unique output aliases, never an empty joined Columns list; lineage remains attached to the two source fields. No output identity overwrites another. An ambiguous last-SELECT alias is diagnosed before Run, and SELECT omission/last position is honored by the adopted grammar.

### AC: joins-are-explained (verifies REQ:source-and-join-discovery)

**Given** a declared Customer–Invoice foreign key, an inferred name match, a cross-source relationship, and an unauthorized table
**When** the user searches joins and selects a suggestion
**Then** authorized candidates show complete key mappings, provenance, cardinality or unknown, and available cost/limit; applying one changes only the draft and retains the relationship binding; unauthorized metadata and execution remain unavailable. Name similarity alone cannot produce a runnable predicate.

### AC: relationship-join-authors-safely (verifies REQ:relationship-backed-join-authoring)

**Given** one unique declared FK `Invoice.CustomerId -> Customer.Id`, two possible Customer FKs, a composite FK, a name-only match, an unauthorized target, and a saved FK binding whose schema mapping later changes
**When** the author writes `FROM Invoice JOIN Customer` without ON, writes `JOIN Customer ON CustomerId`, selects a suggestion, or runs the saved query after schema change
**Then** the unique relationship inserts the fully qualified predicate `Invoice.CustomerId = Customer.Id` before execution, regardless of differing field names, and its relationship identity, ordered pairs, source references, and schema version remain in the saved draft. The shorthand resolves only when unique. Composite completion includes every declared field pair. Ambiguity, partial composite mapping, name-only inference, denied access, or a changed/removed saved relationship blocks execution with a diagnostic and never silently binds a replacement. Suggestions and generated predicates have deterministic ordering, and choosing a suggestion alone never runs.

### AC: preview-matches-execution-path (verifies REQ:preview-and-execution-receipt)

**Given** DTQL with a typed CustomerId binding and two environment mappings
**When** the user previews one target, changes target, and runs
**Then** the old preview becomes stale, the new preview comes from the executing adapter/planner, and the immutable receipt identifies actual requests/local stages, target, placeholders, redacted binding metadata, and execution without credentials or secret values. A non-SQL target never displays fabricated actual SQL.

### AC: run-states-are-truthful (verifies REQ:explicit-run-and-results)

**Given** a prior success followed by an edit with a missing required binding
**When** the user attempts Run, resolves the binding, then runs, cancels or receives a partial/error response, and changes target while a late response is in flight
**Then** Run blocks until resolution; previous results retain their original attribution; terminal states show reason/completeness; cancellation/failure preserves committed output; the late response cannot attach to the new target. Secret binding values appear only as set/unset and are absent from URLs, telemetry, persisted drafts, receipts, cards, and result caches.

### AC: lookup-is-reviewable (verifies REQ:related-query-exploration)

**Given** an authorized selected CustomerId cell and a saved Customer invoices query with a typed parameter
**When** the user chooses Look up this value
**Then** a preview names query, target, bound value and origin before explicit execution, and only that customer's authorized invoices return. Use as data is separately named and rejects cycles, incompatible outputs, and unauthorized inputs.

### AC: target-runs-are-independent (verifies REQ:multi-target-results)

**Given** two authorized same-schema targets and one frozen revision and binding set
**When** the user explicitly runs both and the candidate target fails
**Then** baseline receipt and result remain accessible, candidate shows its own error/retry, and neither side claims an atomic shared snapshot or invented observation time.

### AC: comparison-explains-limits (verifies REQ:honest-comparison)

**Given** synthetic baseline/candidate results with one added keyed row and one changed decimal value, plus a partial/truncated candidate
**When** the user compares them with a unique typed key
**Then** an exact supported comparison reports precise values and each side's receipt, time, policy, and completeness; null/duplicate keys, incompatible schema or declared units, unsupported comparer, and missing/expired/unauthorized sides block keyed diff. Unknown units permit only a labelled raw typed-value comparison. Under the current Compare contract the partial/truncated candidate fails closed with no CompareResult. Only after a reviewed Compare Feature/API amendment is implemented may the user explicitly request a partial sample; then one-sided keys are labelled visible only in baseline/candidate and Added/Removed and completeness totals are suppressed. Unequal policy visibility has the same visible-only rule. A warning rejects causal/deletion assertions; numeric tolerance is unavailable until its Core contract extension passes.

### AC: ai-actions-are-scoped (verifies REQ:ai-action-proposal)

**Given** an authorized draft, a request to change its filter, and a separate explicit request to run on a named target
**When** the composer handles both requests
**Then** the edit is a reversible diff and does not execute; the run action shows revision, target, and bindings before the same governed execution as manual Run; opening the drawer alone sends no result cells to a model.

### AC: dashboard-card-keeps-authority (verifies REQ:dashboard-reference)

**Given** a successful public demo result, a guest's draft edit, and a second principal
**When** the guest chooses Add to dashboard, signs in at Save submit, and the second principal later opens the card
**Then** the nonsecret draft survives sign-in; only an authorized project writer can commit a saved query/card; the UI requires saving or choosing the prior revision and previews destination and policies; it persists a query reference without credentials or copied rows and reauthorizes the second principal before display or refresh. Anonymous sign-in is not used.

### AC: browser-demo-journey (verifies REQ:demo-and-accessibility)

**Given** a fresh browser session and the maintained demo project
**When** an analyst opens the curated Chinook example, changes a filter, previews and runs, selects a customer for invoice lookup, compares deterministic demo targets, signs in at dashboard Save if needed, and reopens the card without hidden setup
**Then** every stage has a visible good result with source/target/parameter evidence; sign-in retains the nonsecret draft and a project write check governs commit; restricted or setup-required examples name their prerequisite. The journey is keyboard operable at 1280/1440 px and 200% zoom with labeled controls and non-color-only diff cues.

## Open Questions

1. Which parser implementation and FROM-first syntax details best fit the approved TugQL semantics and parity constraints in the timeboxed spike? In a separate branding brainstorm, should the alternative pair DTQL/DTQT be considered? This Feature makes no umbrella rename or retirement decision: DTQL remains the library and website name.
2. Which compiler/adapter surfaces can provide both preview and executed receipt for browser SQLite, local-agent SQL, and federation without a second SQL approximation?
3. Which supported browser connection provides the curated Customer-invoices execution, and which current browser receipts prove its stages? Hosted BigQuery is not assumed browser-executable.
4. Which dashboard revision/refresh policies should default for a first card while preserving viewer authorization?

---
*This document follows the https://specscore.md/feature-specification*
