---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Portable query model (DTQL + DALgo)

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/portable-query-model?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P1 is the portable query model: one logical query AST (DALgo `StructuredQuery`), serialized losslessly as DTQL YAML, rendered to each backend's native query language, with native text kept as an explicit escape hatch. The user-approved 2026-10-09 [Query Workspace](../query-workspace/README.md) adds TugQL, a readable FROM-first authoring syntax with optional SELECT last, and TugQTree, the structured YAML form of the same model. This makes the parser implementation mandatory in Go and TypeScript; the pending spike chooses bounded grammar and parser architecture, not whether the text form exists. TugQTree/YAML compatibility remains, with no file migration. New Query Workspace specs/UI use TugQL/TugQTree, while DTQL remains the existing library and website name; this Feature does not authorize an umbrella rename, retirement, redirect, or migration.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

It is the pillar every other pillar leans on: P10's server-enforced ACL can only rewrite a query it can parse as structure, and P7's shared library can only diff and review a query in Git if it is data, not opaque text. Initial audit snapshot (before later 9 September landings): DTQL serialization and the DALgo `access` policy engine are Stable and well tested in `dal-go/dalgo`, but DataTug itself consumes them at exactly one call site, `datatug query run` (`datatug-cli/apps/datatugapp/commands/cmd_query.go`); `QueryDef.Type` has no `DTQL` value yet (`datatug-core/pkg/datatug/query.go:124-137`); and the standalone `dtql`/`dtql-go`/`dtql-backstage` repositories and the `dtql.org` site are obsolete stubs that contradict the real implementation. Phase 0 retires the obsolete repos; Phase 1 (`core-investigation-loop`) adds `DTQL` as a saved-query type; Phase 2 completes the HTTP escape hatch and the remaining backend adapters. This is a Capability; per-repo Implementations (`datatug-cli`, `dal-go/dalgo`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> Query text is source-language-specific: T-SQL, a Firestore query, an inGitDB filter and an HTTP
> request body share no common shape. A query written for one source cannot move to another, cannot
> be structurally rewritten by the server to enforce row/column policy, and cannot be diffed
> meaningfully in Git the way a data file can. DataTug's thesis depends on queries being executable
> knowledge assets, not opaque strings. Today that promise is only kept at one call site: DTQL
> serialization is "usable" and the DALgo `access` package is "usable, spec Stable"
> (`dal-go/dalgo/dtql/README.md`; `dal-go/dalgo/access/*`), but `QueryDef.Type` in DataTug's own
> project model still only accepts `SQL`/`HTTP`/`StructuredSQL`
> (`datatug-core/pkg/datatug/query.go:124-137`), so no project file can hold a DTQL document. Backend
> rendering is likewise partial: `sqlite://` and `ingitdb://` are fully wired but `postgres://`
> `Open()` always returns `ErrPostgresNotWired` (`pkg/dbcopy/url.go:34-37,164-165`), and `mssql`/
> `mysql`/`firestore` are not in the dispatcher at all. Marketing on `datatug.io` and the standalone
> `dtql.org` site present DTQL as an installable product independent of `dal-go/dalgo`; the audit
> found the `dtql`/`dtql-go`/`dtql-backstage` repos are one-line-README stubs with a single 2026-06-09
> commit and zero `.go` files, while the real implementation lives under a different module path
> entirely (current-state audit, "Contradictions" table, DTQL row).

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| DTQL YAML serialization of the DALgo query AST | usable (external, Stable in `dal-go/dalgo/dtql`) | — | enabler |
| DTQL as a saved-query type in projects (`queries/**/*.query.dtql`, `QueryDef.Type = DTQL`) | planned | 1 | strategic |
| Native query escape hatch (SQL text, HTTP request) | usable (SQL) / backend-only (HTTP) | 1 (SQL), 2 (HTTP) | enabler |
| Rendering to backends through DALgo adapters | partial | 0–1 | enabler |
| Standalone `dtql*` repositories and dtql.org claims | obsolete / contradicted | 0 | hygiene |

#### REQ: dtql-ast-consumption

DataTug MUST treat `dal-go/dalgo`'s `StructuredQuery` AST and its DTQL YAML serialization as the
single source of truth for a portable query; DataTug repos MUST NOT fork or reimplement the
parser, serializer or the DALgo `access` policy evaluator. Any DataTug-side query representation
(saved query, related-record lookup, applicable-query resolution, or an adopted alternate authoring
frontend) SHOULD be expressible as this AST so the server can apply row/column policy uniformly
(see `server-enforced-access`). An alternate frontend MUST map into the same AST and retain DTQL
YAML compatibility for saved-query forms that use it; considering a grammar does not change existing
YAML documents or claim new parsing support.

For the required TugQL frontend, the [Query Workspace authoring contract](../query-workspace/README.md#req-named-parameter-block) makes human readability and Git-friendly, line-local edits primary. Compact single-line SELECT lists keep their comma syntax. Multiline projections require `select (` with one selected item per indented logical line, no item commas, and a closing `)` alone aligned with the header. When formatting expands a compact list to multiple lines, it uses that block form and preserves commas nested in expressions, function arguments, strings, and quoted identifiers. Parenthesized compact expressions remain legal; the parser identifies a SELECT block from the newline after its header parenthesis, not from `(` alone. The optional `PARAMETERS (` header is on one line before `FROM`, followed by typed `@Name` declarations, one per logical line with newline separators and no commas, with `(` on that header line and `)` on its own aligned line; `WITH <Name> AS (` and `USING (` block headers follow the same opening/closing placement. Each query clause also starts on its own line. Structural indent is one tab or two spaces per level, consistent per query: header/closing parenthesis/top-level clauses have zero indent, declarations and expression continuations use one indent level, and ON is indented one level under JOIN. Each CTE query body adds one level beneath its WITH declaration, and its closing parenthesis returns to the declaration level. Mixed tabs/spaces and missing required indentation are diagnosed; the formatter can repair indentation with project/team, user, preserve-existing, then application-default precedence. Indent preferences are tab, two spaces, or preserve existing; default is two spaces and the current workspace user preference is two spaces. A table or selected-field alias may be omitted, but when provided it requires same-line `AS`. Executable queries use reserved words consistently all lowercase or all uppercase; mixed case is diagnosed but remains editable and can be repaired by formatting. The formatter's style precedence is project/team preference, user preference, first wholly cased reserved word in source order, then application default. Preferences are lowercase, uppercase, or preserve existing; preserve existing uses the first valid reserved word and falls back to the application default. The current Query Workspace user preference is lowercase. Formatting and autoformat-on-save normalize reserved-word casing and structural indentation; if a SELECT list is expanded into multiple lines, they use the SELECT block rules above. Merely reading/parsing a consistent query does not reformat it. Formatting, autocomplete and inferred syntax honor the resolved style, while identifiers, parameters, functions, strings, comments, quoted identifiers and user-defined type names retain authored spelling (including `FROM "Order"`), and formatting never runs a query. Raw parse/save without formatting preserves the source; formatting is stable and idempotent and preserves all non-keyword source identity. Day-one nonrecursive WITH CTEs and multiple named CTEs map to the same canonical query model, share the preceding PARAMETERS bindings, and resolve outputs with source lineage; exact multiple-CTE syntax and recursive support remain open. Parameter declarations and defaults map to that same model, with runtime values carried separately from source text. They MUST NOT create a second AST or break existing DTQL YAML. Saved-query imports use `with <Name> from "<project-relative-path>"` plus an optional `using (` block with the opening parenthesis on the header line and closing parenthesis on its own aligned line; they resolve to the same canonical query model at the importing project revision, with normalized paths/revisions recorded as dependencies. They validate imported REQUIRED/default declarations and never implicitly capture caller parameters. The complete type catalog, assignment rules and required-NULL behavior remain open for the implementation spike; imported native SQL continues to follow its source dialect.

TugQL named scalar subqueries in SELECT require a `<Name> as (` header with the opening parenthesis on the same line, a nested FROM-first body, and an aligned closing `)` on its own line; the name precedes the body, with no alias after the close. `WITH` and `SCALAR` are not selected-item introducers; the nested body follows the ordinary supported query grammar. They lower to the same canonical query model, including qualified outer correlations and lexical alias scopes. A scalar result must have one output column; runtime cardinality is zero rows → NULL, one row → value, more than one row → runtime cardinality error. Multi-column scalar results are invalid. Do not claim that parsing or schema resolution statically predicts row cardinality. The nested body may omit SELECT only when the existing optional-SELECT rule resolves its authorized projection to one column; an outer query containing this selected item has SELECT. Ordinary `expression as Alias`, mandatory `WITH Name as (` and `with Name from "path"` imports retain their own syntax.

#### REQ: dtql-saved-query-type

A project MUST be able to store a query as a `<id>.query.dtql` sidecar with `QueryDef.Type = DTQL`,
loaded and executed the same way as a `SQL`-type query. This is specified in full by
`core-investigation-loop` REQ:dtql-query-type; this Feature only asserts it as a P1 capability and
does not restate the mechanics.

#### REQ: native-query-escape-hatch

Where a source's query language cannot be fully expressed as a `StructuredQuery` (a stored
procedure call, a complex HTTP request, a hand-tuned index hint), a project MUST be able to save
the native text verbatim as a `QueryDef` with a non-DTQL type. Native queries MUST be labelled as
such everywhere they are shown, and MUST NOT be silently treated as structured for ACL purposes
(see `server-enforced-access` REQ referencing the opaque-SQL limitation).

#### REQ: backend-adapter-rendering

Every source DataTug claims to support (P2) MUST have a DALgo adapter that can either execute the
`StructuredQuery` AST directly or explicitly reject the unsupported shape (never silently drop a
clause). `sqlite://` and `ingitdb://` MUST stay usable as the reference adapters; `postgres://` and
the remaining SQL/document adapters MUST either render correctly or fail loudly, never partially.

#### REQ: retire-obsolete-dtql-repos

The standalone `dtql`, `dtql-go`, `dtql-backstage` repositories and the `dtql.org` site MUST NOT be
presented as the canonical DTQL implementation while they contain no working code. They MUST either
be retired, redirected to `dal-go/dalgo/dtql`, or brought current before any public claim repeats
their URLs. This is a separate historical public-claims hygiene concern; the approved Query Workspace
does not activate retirement or redirect work. DTQL remains the current DALgo library and website
name, and no umbrella rename, site retirement, redirect, or migration is authorized here.

## Acceptance Criteria

### AC: dtql-is-source-of-truth (verifies REQ:dtql-ast-consumption)

**Given** a DataTug repository adds new query behaviour, including an adopted text-authoring frontend
**When** the change touches query parsing, serialization or policy evaluation
**Then** it maps structured authoring to DALgo's `StructuredQuery` and reuses `dal-go/dalgo/dtql` and
`dal-go/dalgo/access` rather than adding a parallel AST, serializer or policy evaluator; code review
verifies the boundary. A TugQL parameter block and its typed declarations map to that same model,
while runtime values remain separate from source text. Existing YAML DTQL serialization remains its
own compatibility contract, and imported native SQL remains source-dialect text.
Named scalar subqueries also lower to this same model with matching output names and outer-reference identities; lexical local aliases do not escape. Validate one-column shape, and keep zero/one/many-row cardinality in separate execution fixtures.

### AC: dtql-query-type-saves (verifies REQ:dtql-saved-query-type)

**Given** the demo project and a `customer-invoices.query.dtql` document
**When** `datatug entity list` and `datatug query run` load the project
**Then** the query loads with `Type: DTQL` and runs through the same policy path as a `SQL` query (see `core-investigation-loop` AC:dtql-query-runs).

### AC: native-query-labelled (verifies REQ:native-query-escape-hatch)

**Given** the demo project's `albums_by_title` SQL-text query
**When** it is listed in the web UI or `datatug query run`
**Then** its type is shown as native SQL, not DTQL, and its result header carries the opaque-SQL policy limitation when one applies.

### AC: sqlite-and-ingitdb-render (verifies REQ:backend-adapter-rendering)

**Given** a `StructuredQuery` built for `Customer` in `chinook.sqlite` and for `support-notes` in the demo project's inGitDB collection
**When** each is executed through its DALgo adapter
**Then** both return rows without dropping any WHERE clause, and a `postgres://` target explicitly returns "not wired" rather than partial results.

### AC: no-canonical-claim-for-stub-repos (verifies REQ:retire-obsolete-dtql-repos)

**Given** `datatug.io` or `dtql.org` copy referencing an installable `dtql` package
**When** the marketing content is reviewed against the current-state audit
**Then** either the copy is corrected to point at `dal-go/dalgo/dtql` or the standalone repos are retired before the claim is repeated.

## Not Doing / Out of Scope

- Connector breadth as a growth lever: rendering to every possible backend is not a race; the
  minimum set that proves the architecture (SQL, inGitDB, Firestore, HTTP) comes first.
- A BI/warehouse tool: DTQL exists to make operational queries portable and policy-checkable, not
  to centralize data into a warehouse.
- Client-direct database access for protected sources: the native-text escape hatch still executes
  through the trusted server, never a direct browser connection.
- A text-to-SQL chat box: DTQL is a structured serialization format, not a natural-language query
  interface; AI-assisted query building is P14's concern, not P1's.

## Open Questions

1. Which bounded FROM-first grammar details and parser architecture should implement the required TugQL form while preserving DALgo's canonical model and TugQTree/YAML compatibility?
2. Which initial parameter types and assignment rules can DALgo's canonical model represent, and should a narrowly scoped additive contract be proposed before the TugQL adapter is implemented? Keep required-NULL behavior distinct from a missing binding until explicitly decided.
3. The older public-claims hygiene question about the standalone `dtql*` sites/repos remains separate future work; Query Workspace does not activate the prior optional retire/redirect clause, and the DTQL library and website names remain unchanged.

---
*This document follows the https://specscore.md/feature-specification*
