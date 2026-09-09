---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Serve-Brokered AI Query Builder

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=request-change) |
**Status:** Approved
**Date:** 2026-06-05
**Owner:** alexander.trakhimenok
**Source Ideas:** serve-brokered-query-builder
**Supersedes:** —
**Grade:** A

## Summary

A platform-agnostic **Capability** for an AI query builder whose conversation lives in the user's existing AI-agent terminal while a local DataTug CLI daemon (`serve`) holds the query and a hosted Web UI shows and deterministically edits it.

The daemon brokers a single canonical query per named-query *tab* between an MCP face (the agent) and an HTTP/WS face (the Web UI), keeping both in sync. Each tab is in one of **two modes**: **DTQL mode** (default) — the query is a dalgo query AST whose textual form is a 1:1 **DTQL (DataTug Query Language) YAML** serialization, deterministically editable point-and-click in the Web UI and rendered to the connection's native query language on execution; or **native mode** — the query is the connection's own native text (SQL, GraphQL, a document-store query, etc.) plus parameters, held verbatim, edited as text, and executed as-is. The dalgo AST spans relational and nested/document (e.g. Firestore) shapes, so DTQL mode is not SQL-only. A tab may start in either mode and may convert DTQL→native one-way; native→DTQL is not supported. This is an alternative delivery architecture to the in-app `ai-query-builder` Capability: no AI-chat panel is built in any surface. The Capability defines only the cross-surface contract and owns the Implementation Matrix tracking the per-repo realizations.

## Problem

DataTug needs query building bound to live connections, the current result set, and a turn-by-turn loop — so it must be in-app — but building a conversational chat panel into both a CLI and a Web UI duplicates work and drifts. This Capability removes the chat panel entirely: the agent terminal the user already runs *is* the conversation, a single daemon owns the query state, and the Web UI is a deterministic editor plus live results. Without one platform-agnostic contract, the daemon, the Web UI, and the agent skill would each re-invent the protocol and diverge. This pins the contract once: what "one query, edited from a terminal agent and a web UI, brokered by a local daemon" must mean on every surface.

## Behavior

### Session state and representation

The daemon is the single source of truth for the live query; the conversation's history is not its concern.

#### REQ: tab-current-version

The daemon MUST hold the current version of each named query (a *tab*) as the single source of truth, and MUST NOT retain query history — history is owned by the Web UI. A tab is addressed by an explicit tab id on every operation.

#### REQ: query-mode

Each tab MUST have a query mode — either **`dtql`** (the default) or **`native`** — that determines how its current query is represented and edited. The mode is part of the tab's state and is reported on inspection. DTQL mode holds a dalgo AST (REQ:canonical-ast); native mode holds verbatim native text (REQ:native-mode-representation).

#### REQ: canonical-ast

In `dtql` mode the current query MUST be represented as one canonical dalgo query AST per tab. This single representation is shared by all writers; it is rendered to the connection driver's native query language only on execution. The AST is source-shape-agnostic — it expresses both relational queries (rendered to the tab connection's SQL dialect) and nested/document queries (rendered to the document store's native query form). In `native` mode there is no AST; see REQ:native-mode-representation.

#### REQ: dtql-yaml-form

DTQL (DataTug Query Language) is the canonical **textual form of the dalgo AST**: a 1:1, lossless YAML serialization. It MUST round-trip with the AST (parse YAML → AST → serialize YAML with no semantic loss) using a standard YAML parser — DTQL is a serialization of the model, not a separate grammar requiring a bespoke parser. The "dalgo form" the agent and Web UI show or hand-edit is this DTQL-YAML. (Adopting an existing relational/analytical text language — e.g. PRQL or Malloy — as an alternate textual surface is out of scope this cycle; see Not Doing.)

#### REQ: native-mode-representation

In `native` mode the tab MUST hold the connection's own native query text plus its parameter set as the source of truth, verbatim. The daemon MUST NOT maintain a dalgo AST for a native-mode tab, MUST NOT parse or rewrite the native text, and MUST execute it as-is against the tab's connection. Structured point-and-click editing is unavailable in native mode; the query is edited as text (REQ:web-native-edits).

#### REQ: mode-transition

A tab MAY be created in either mode. A `dtql` tab MAY be converted to `native` once and irreversibly: the daemon renders the current AST to native text in the tab's connection dialect, adopts that text as the native-mode current query, and discards the AST. Conversion from `native` to `dtql` MUST NOT be offered — native text is never parsed back into the AST.

#### REQ: tab-connection-binding

Each tab MUST be bound to exactly one DataTug connection, selected when the tab is created; that connection's driver fixes the tab's native query language (SQL dialect, or a document store's native query form) for the life of the tab.

### Agent interface (MCP)

The agent drives the query through structured tools — the LLM stays in the terminal.

#### REQ: mcp-agent-face

The daemon MUST expose an MCP face offering, at minimum, tools to create a tab (in either mode), apply a change to a tab, inspect a tab's query, convert a `dtql` tab to `native`, and run a tab's query.

#### REQ: apply-change-payload

An apply-change operation MUST target a tab by id and carry the prose description of the change plus the full resulting query. In `dtql` mode the full query is the dalgo AST (as DTQL-YAML) and the operation SHOULD also carry the structured delta; in `native` mode the full query is the native text plus parameters and the structured delta is omitted (there is no AST to diff). In both modes the daemon MUST adopt the supplied full query as that tab's new current version and MUST NOT re-interpret the prose to derive the query. An apply-change MUST match the tab's mode.

#### REQ: terminal-quiet-by-default

By default the agent MUST NOT echo the query in the terminal; it sends each change to the daemon and the user reads the query in the Web UI. On explicit request the agent MUST be able to show the query — in `dtql` mode as DTQL-YAML (the dalgo form) or as the rendered native query, and in `native` mode as the native text — once or on every change. Refusals and unresolvable requests MUST be reported in the terminal.

### Web interface (HTTP/WS)

The Web UI is a deterministic editor and live view — never a prose/chat surface.

#### REQ: web-deterministic-edits

In `dtql` mode the Web UI MUST apply deterministic, structured edits (such as filters, field/column add/remove/select, ordering, and nested/sub-collection selection) to the same canonical dalgo AST, expressed as structured operations and never as natural-language prose.

#### REQ: web-native-edits

In `native` mode the Web UI MUST present the tab's native query text in an editable text area and let the user edit and submit it verbatim. This is query-language text, not chat prose: the Web UI still never sends natural-language instructions to be interpreted — it sends the native query the engine executes as-is.

#### REQ: query-parameters

A tab's query MAY declare named parameters. The Web UI MUST let the user view and set parameter values in both modes, and parameters MUST be bound at execution. In `dtql` mode parameters are nodes of the dalgo AST; in `native` mode they accompany the native text.

#### REQ: web-history-and-revert

The Web UI MUST hold the change history (the running list of prose change-descriptions) and MUST allow reverting to any prior step; reverting sends that prior version to the daemon as the new current query. History is client-side and need not survive a page reload.

#### REQ: results-presentation

When a tab's query is executed, its results MUST be presented with a Web-UI-configurable row limit (default 1000) and "load next N" / "load all" controls rather than server paging, and the executed query MUST be re-runnable without re-stating the conversation.

### Synchronization

#### REQ: two-way-sync

Changes from either face MUST reach the other and the agent: agent-applied changes MUST stream to the Web UI, and a Web-UI edit MUST signal the agent so it reconciles by re-reading the tab's current query. Both faces operate on the tab's single current query (the canonical dalgo AST in `dtql` mode, the verbatim native text in `native` mode), so neither surface nor the agent drifts.

### Conversational disambiguation

#### REQ: candidate-options

When the agent is not confident which interpretation to apply, it MUST be able to submit multiple candidate options for a tab (each a full query in the tab's mode — a dalgo AST in `dtql` mode, native text in `native` mode); the Web UI MUST let the user compare them by previewing each option's results and commit one as the new current query, after which the unchosen options are discarded.

### Onboarding and access

#### REQ: deep-link-onboarding

To open the Web UI, the agent MUST create the tab, select its connection, and return a deep link carrying the query id, the local backend host, and a one-time code (in the link fragment) that the Web UI exchanges for the daemon's session token; opening the link MUST attach the hosted Web UI to that daemon and tab.

#### REQ: local-loopback-access

The daemon MUST bind to loopback and require a shared token; the Web UI (the hosted app at https://datatug.app by default, overridable to a custom URL) connects from the user's browser over a CORS-allowlisted origin.

### Safety

#### REQ: read-only-queries

The builder MUST produce only read-only queries, enforced per mode. In `dtql` mode the dalgo AST is read-only by construction — it cannot express data or schema mutation. In `native` mode, where the text is opaque and executed as-is, the daemon MUST execute every query through a read-only session/transaction on the connection so that any mutating statement (insert/update/delete/DDL) fails at the engine; this enforcement is driver-agnostic and does not rely on parsing the native text. In both modes a request whose evident intent is to mutate data or schema MUST be refused where detectable and MUST leave the tab's current query unchanged.

### Cross-surface parity

#### REQ: parity-matrix

This Capability MUST own an `## Implementation Matrix` recording every realization's parity (`Platform`, `Status`, one-line `Brief`, `Link`), so a reader of the Capability alone can see each surface's state without opening the per-repo repositories.

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI (daemon) | Planned | `serve` daemon: per-tab dual-mode state (DTQL dalgo-AST ⇄ DTQL-YAML, or native text+params), native renderer + read-only execution, MCP + HTTP/WS faces, two-way sync, disambiguation | specscore:feature/serve-brokered-query-builder@github.com/datatug/datatug-cli |
| Web | Planned | `query-builder` screen (DTQL point-and-click + native textarea + params) + `ai-terminal-query-builder` comms to the DAL backend | specscore:feature/query-builder@github.com/datatug/datatug-apps |
| Agent | Planned | `query-builder` skill driving the daemon's MCP tools from the terminal: DTQL-YAML vs native mode, one-way DTQL→native conversion | specscore:feature/query-builder@github.com/datatug/datatug-ai-skills |

## Acceptance Criteria

### AC: daemon-holds-current-no-history (verifies REQ:tab-current-version)

**Given** a tab whose query has been changed several times
**When** the daemon is asked for that tab's query by its id
**Then** it returns only the latest version, and the daemon exposes no prior-version history of its own.

### AC: tab-has-mode (verifies REQ:query-mode)

**Given** a tab created in `dtql` mode and a tab created in `native` mode
**When** each tab is inspected
**Then** each reports its mode, and a `dtql` tab exposes a dalgo AST while a `native` tab exposes verbatim native text — never both.

### AC: ast-is-canonical (verifies REQ:canonical-ast)

**Given** a `dtql`-mode tab with a non-empty current query
**When** the query is executed
**Then** the same canonical dalgo AST is rendered to the tab connection's native query language (the SQL dialect for a SQL connection, the native query form for a document store), and that AST is the representation every writer updated.

### AC: dtql-yaml-round-trips (verifies REQ:dtql-yaml-form)

**Given** a `dtql`-mode tab's dalgo AST
**When** it is serialized to DTQL-YAML and that YAML is parsed back
**Then** the resulting AST is identical to the original (lossless 1:1 round-trip), and the YAML is what the agent and Web UI show as the dalgo form.

### AC: native-tab-held-verbatim (verifies REQ:native-mode-representation)

**Given** a `native`-mode tab seeded with native query text and parameters
**When** the tab is inspected and run
**Then** the daemon returns and executes the native text unchanged, maintains no dalgo AST for it, and does not rewrite it.

### AC: dtql-to-native-one-way (verifies REQ:mode-transition)

**Given** a `dtql`-mode tab
**When** it is converted to `native`
**Then** the AST is rendered to native text in the tab's dialect, the tab becomes `native` holding that text, the AST is discarded, and no operation converts it back to `dtql`.

### AC: parameters-bound-both-modes (verifies REQ:query-parameters)

**Given** a tab declaring a named parameter, in either mode
**When** the user sets the parameter value in the Web UI and runs the query
**Then** the value is bound at execution, and the parameter is an AST node in `dtql` mode or accompanies the native text in `native` mode.

### AC: tab-bound-to-one-connection (verifies REQ:tab-connection-binding)

**Given** a tab created against a selected connection
**When** any later operation runs on that tab
**Then** the tab stays bound to that one connection and its dialect does not change for the life of the tab.

### AC: mcp-tools-present (verifies REQ:mcp-agent-face)

**Given** a running daemon
**When** an MCP client lists the available tools
**Then** tools to create a tab, apply a change, inspect, and run are available.

### AC: apply-change-adopts-full-query (verifies REQ:apply-change-payload)

**Given** a `dtql`-mode tab with a current query
**When** the agent calls apply-change for that tab id with prose, a structured delta, and a full resulting query (the dalgo AST)
**Then** the daemon makes the supplied full query the tab's new current version without re-deriving it from the prose; and for a `native`-mode tab the same holds with the full native text (no delta required).

### AC: web-native-edit-is-verbatim (verifies REQ:web-native-edits)

**Given** a `native`-mode tab shown in the Web UI
**When** the user edits the native query text and submits it
**Then** the edited text is sent and adopted verbatim as the tab's current query, with no natural-language interpretation step.

### AC: agent-quiet-then-shows-on-request (verifies REQ:terminal-quiet-by-default)

**Given** a session where the agent has applied a change
**When** the user has not asked to see the query
**Then** the query is not echoed in the terminal; and when the user asks to see it, the agent shows the dalgo or native form.

### AC: web-edit-is-structured (verifies REQ:web-deterministic-edits)

**Given** a `dtql`-mode tab with a current query shown in the Web UI
**When** the user adds a filter via the Web UI controls
**Then** the change is sent as a structured operation on the canonical dalgo AST, not as prose.

### AC: web-revert-sets-current (verifies REQ:web-history-and-revert)

**Given** a Web UI showing a change history of several steps
**When** the user reverts to an earlier step
**Then** that earlier version is sent to the daemon as the new current query for the tab.

### AC: results-limited-and-rerunnable (verifies REQ:results-presentation)

**Given** an executed query whose result set exceeds the configured row limit
**When** the results are presented
**Then** at most the configured limit (default 1000) is shown with "load next N" / "load all" controls, and the same query can be re-run without re-stating the conversation.

### AC: edit-syncs-both-ways (verifies REQ:two-way-sync)

**Given** an attached Web UI and an agent on the same tab
**When** the user makes a deterministic edit in the Web UI
**Then** the agent is signalled and, on re-reading the tab, sees the updated current query — and an agent-applied change likewise appears in the Web UI.

### AC: unsure-offers-options (verifies REQ:candidate-options)

**Given** a request the agent cannot confidently resolve to one query
**When** the agent submits multiple candidate options for the tab
**Then** the Web UI lets the user preview each option's results and commit one, after which the unchosen options are discarded.

### AC: link-attaches-web-ui (verifies REQ:deep-link-onboarding)

**Given** the user starts in the agent terminal
**When** the agent creates a tab and returns its deep link, and the user opens it
**Then** the hosted Web UI attaches to that daemon and tab using the query id, backend host, and a one-time code exchanged for the session token.

### AC: loopback-token-required (verifies REQ:local-loopback-access)

**Given** the daemon is running on loopback
**When** a client connects without the valid shared token
**Then** access is refused; and the hosted Web UI origin is permitted via CORS when it presents the token.

### AC: mutation-refused (verifies REQ:read-only-queries)

**Given** a `dtql`-mode tab with a current query
**When** the user asks to "delete all inactive users"
**Then** the request is refused, no mutating query is produced, and the tab's current query is unchanged.

### AC: native-mutation-blocked-at-engine (verifies REQ:read-only-queries)

**Given** a `native`-mode tab whose text contains a mutating statement (e.g. `DELETE`/`DROP`)
**When** the query is executed
**Then** it runs through a read-only session/transaction and the engine rejects the mutation, so no data or schema change occurs — without the daemon parsing the native text.

### AC: matrix-tracks-surfaces (verifies REQ:parity-matrix)

**Given** the Capability README
**When** a reader inspects it
**Then** an `## Implementation Matrix` is present with one row per surface, each carrying a Status, a one-line Brief, and a Link to that surface's realization.

## Rehearse Integration

Every AC has a concrete behavioral surface (daemon state, MCP tools, structured web edits, sync, refusal paths, document structure) and is testable in principle. However, this Capability defines cross-surface *contract* behavior with no executable surface of its own — each AC is exercised against a concrete Implementation (the CLI daemon in `datatug-cli`, the Web features in `datatug-apps`), where the Rehearse stubs will live. Stub scaffolding is therefore deferred to the Implementation Features so the stubs track real, runnable surfaces rather than an abstract Capability.

## Not Doing / Out of Scope

- In-app AI-chat panel on any surface — the agent terminal is the only conversational surface; the Web UI is deterministic controls (DTQL mode) or a native-text area (native mode), not chat prose.
- `native`→`dtql` conversion (parsing native text back into the dalgo AST) — DTQL→native is one-way only this cycle.
- A bespoke DTQL grammar/parser — DTQL is a 1:1 YAML serialization of the dalgo AST parsed by a standard YAML parser; adopting an existing text language (PRQL/Malloy) as an alternate surface is a post-MVP option, not this cycle.
- Query history in the daemon — the daemon holds only current versions; history is Web-UI-held and need not survive a reload.
- Session persistence across daemon restarts — in-memory only; a restart starts fresh.
- Remote multi-user serve and link-sharing — local single-session (loopback + token); reverse-proxy sharing (ngrok/Cloudflare) is a post-MVP follow-on the token model anticipates.
- Seeding a tab from a saved DataTug query or a table/entity pick — natural-language seeding only this cycle.
- Mutating or DDL queries — read-only only (REQ:read-only-queries).
- A bespoke DataTug agent/LLM loop — the conversation is driven by external agents via the `query-builder` skill (`datatug-ai-skills`), not a conversation engine built into DataTug.
- Mandating concurrency control for simultaneous agent + Web-UI edits — last-write-wins this cycle.

## Assumption Carryover

From the source Idea `serve-brokered-query-builder`:

- **Agent can emit a faithful prose + delta + full-query triple the daemon adopts as canonical (Must-be-true)** — encoded as REQ:apply-change-payload + REQ:canonical-ast, validated by AC:apply-change-adopts-full-query against the CLI daemon Implementation.
- **The daemon keeps both faces and the agent in sync two-way (Must-be-true)** — encoded as REQ:two-way-sync, validated by AC:edit-syncs-both-ways.
- **Deterministic web editing is worth building alongside conversational editing (Should-be-true)** — encoded as REQ:web-deterministic-edits + REQ:web-history-and-revert; confirmed only by dogfooding an Implementation.
- **The agent can produce distinct, meaningful candidate options (Should-be-true)** — encoded as REQ:candidate-options; confirmed only by dogfooding.

## Open Questions

- Should the CLI daemon and the two `datatug-apps` features (`query-builder`, `ai-terminal-query-builder`) be specified next, and in what order, given the daemon contract is the dependency for both?
- Does the agent target a tab purely by explicit id, or also via an implicit "active tab" convenience at the skill layer? (The contract here mandates explicit id.)
- For nested/document (e.g. Firestore) sources in `dtql` mode, which structured Web-UI controls and renderer coverage ship in the first cycle vs. follow Partial — and is `native` mode the pragmatic fallback for document-query shapes the AST does not yet model?
- Is DTQL-YAML defined inline by the dalgo AST schema, or does it warrant its own small schema/spec artifact (and where — `datatug-go-models`)?

## Sidekick Seeds Generated

- [dtql-datatug-query-lanaguage-needs-definition-and](../../ideas/seeds/dtql-datatug-query-lanaguage-needs-definition-and.md) — captured 2026-06-05 by user

---
*This document follows the https://specscore.md/feature-specification*
