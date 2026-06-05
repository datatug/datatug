# Feature: Serve-Brokered AI Query Builder

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/serve-brokered-query-builder?op=request-change) |
**Status:** Approved
**Date:** 2026-06-05
**Owner:** alexander.trakhimenok
**Source Ideas:** serve-brokered-query-builder
**Supersedes:** —
**Grade:** A

## Summary

A platform-agnostic **Capability** for an AI query builder whose conversation lives in the user's existing AI-agent terminal while a local DataTug CLI daemon (`serve`) holds the query and a hosted Web UI shows and deterministically edits it. The daemon brokers a single canonical query — represented as one dalgo AST per named-query *tab* — between an MCP face (the agent) and an HTTP/WS face (the Web UI), keeping both in sync. This is an alternative delivery architecture to the in-app `ai-query-builder` Capability: no AI-chat panel is built in any surface. The Capability defines only the cross-surface contract and owns the Implementation Matrix tracking the per-repo realizations.

## Problem

DataTug needs query building bound to live connections, the current result set, and a turn-by-turn loop — so it must be in-app — but building a conversational chat panel into both a CLI and a Web UI duplicates work and drifts. This Capability removes the chat panel entirely: the agent terminal the user already runs *is* the conversation, a single daemon owns the query state, and the Web UI is a deterministic editor plus live results. Without one platform-agnostic contract, the daemon, the Web UI, and the agent skill would each re-invent the protocol and diverge. This pins the contract once: what "one query, edited from a terminal agent and a web UI, brokered by a local daemon" must mean on every surface.

## Behavior

### Session state and representation

The daemon is the single source of truth for the live query; the conversation's history is not its concern.

#### REQ: tab-current-version

The daemon MUST hold the current version of each named query (a *tab*) as the single source of truth, and MUST NOT retain query history — history is owned by the Web UI. A tab is addressed by an explicit tab id on every operation.

#### REQ: canonical-ast

The current query MUST be represented as one canonical dalgo AST per tab. This single representation is shared by all writers; it is rendered to executable SQL only on execution, in the dialect of the tab's connection.

#### REQ: tab-connection-binding

Each tab MUST be bound to exactly one DataTug connection, selected when the tab is created; that connection's driver fixes the tab's SQL dialect for the life of the tab.

### Agent interface (MCP)

The agent drives the query through structured tools — the LLM stays in the terminal.

#### REQ: mcp-agent-face

The daemon MUST expose an MCP face offering, at minimum, tools to create a tab, apply a change to a tab, inspect a tab's query, and run a tab's query.

#### REQ: apply-change-payload

An apply-change operation MUST target a tab by id and carry three parts: the prose description of the change, the structured delta, and the full resulting query. The daemon MUST adopt the supplied full query as that tab's new current version and MUST NOT re-interpret the prose to derive the query.

#### REQ: terminal-quiet-by-default

By default the agent MUST NOT echo the query in the terminal; it sends each change to the daemon and the user reads the query in the Web UI. On explicit request the agent MUST be able to show the query in dalgo form or native (dialect SQL) form, once or on every change. Refusals and unresolvable requests MUST be reported in the terminal.

### Web interface (HTTP/WS)

The Web UI is a deterministic editor and live view — never a prose/chat surface.

#### REQ: web-deterministic-edits

The Web UI MUST apply deterministic, structured edits (such as filters, column add/remove/select, ordering) to the same canonical dalgo AST, expressed as structured operations and never as prose.

#### REQ: web-history-and-revert

The Web UI MUST hold the change history (the running list of prose change-descriptions) and MUST allow reverting to any prior step; reverting sends that prior version to the daemon as the new current query. History is client-side and need not survive a page reload.

#### REQ: results-presentation

When a tab's query is executed, its results MUST be presented with a Web-UI-configurable row limit (default 1000) and "load next N" / "load all" controls rather than server paging, and the executed query MUST be re-runnable without re-stating the conversation.

### Synchronization

#### REQ: two-way-sync

Changes from either face MUST reach the other and the agent: agent-applied changes MUST stream to the Web UI, and a Web-UI edit MUST signal the agent so it reconciles by re-reading the tab's current query. Both faces operate on the one shared canonical AST, so neither surface nor the agent drifts.

### Conversational disambiguation

#### REQ: candidate-options

When the agent is not confident which interpretation to apply, it MUST be able to submit multiple candidate options for a tab; the Web UI MUST let the user compare them by previewing each option's results and commit one as the new current query, after which the unchosen options are discarded.

### Onboarding and access

#### REQ: deep-link-onboarding

To open the Web UI, the agent MUST create the tab, select its connection, and return a deep link carrying the query id, the local backend host, and a one-time code (in the link fragment) that the Web UI exchanges for the daemon's session token; opening the link MUST attach the hosted Web UI to that daemon and tab.

#### REQ: local-loopback-access

The daemon MUST bind to loopback and require a shared token; the Web UI (the hosted app at https://datatug.app by default, overridable to a custom URL) connects from the user's browser over a CORS-allowlisted origin.

### Safety

#### REQ: read-only-queries

The builder MUST produce only read-only queries. A request to mutate data or schema (insert/update/delete/DDL) MUST be refused and MUST leave the tab's current query unchanged.

### Cross-surface parity

#### REQ: parity-matrix

This Capability MUST own an `## Implementation Matrix` recording every realization's parity (`Platform`, `Status`, one-line `Brief`, `Link`), so a reader of the Capability alone can see each surface's state without opening the per-repo repositories.

## Implementation Matrix

| Platform | Status | Brief | Link |
|----------|--------|-------|------|
| CLI (daemon) | Planned | `serve` daemon: per-tab dalgo-AST state, MCP + HTTP/WS faces, two-way sync, disambiguation | specscore:feature/serve-brokered-query-builder@github.com/datatug/datatug-cli |
| Web | Planned | `query-builder` screen + `ai-terminal-query-builder` comms to the DAL backend | specscore:feature/query-builder@github.com/datatug/datatug-apps |
| Agent | Planned | `query-builder` skill driving the daemon's MCP tools from the terminal | specscore:feature/query-builder@github.com/datatug/datatug-ai-skills |

## Acceptance Criteria

### AC: daemon-holds-current-no-history (verifies REQ:tab-current-version)

**Given** a tab whose query has been changed several times
**When** the daemon is asked for that tab's query by its id
**Then** it returns only the latest version, and the daemon exposes no prior-version history of its own.

### AC: ast-is-canonical (verifies REQ:canonical-ast)

**Given** a tab with a non-empty current query
**When** the query is executed
**Then** the same canonical dalgo AST is rendered to SQL in the tab connection's dialect, and that AST is the representation every writer updated.

### AC: tab-bound-to-one-connection (verifies REQ:tab-connection-binding)

**Given** a tab created against a selected connection
**When** any later operation runs on that tab
**Then** the tab stays bound to that one connection and its dialect does not change for the life of the tab.

### AC: mcp-tools-present (verifies REQ:mcp-agent-face)

**Given** a running daemon
**When** an MCP client lists the available tools
**Then** tools to create a tab, apply a change, inspect, and run are available.

### AC: apply-change-adopts-full-query (verifies REQ:apply-change-payload)

**Given** a tab with a current query
**When** the agent calls apply-change for that tab id with prose, a structured delta, and a full resulting query
**Then** the daemon makes the supplied full query the tab's new current version without re-deriving it from the prose.

### AC: agent-quiet-then-shows-on-request (verifies REQ:terminal-quiet-by-default)

**Given** a session where the agent has applied a change
**When** the user has not asked to see the query
**Then** the query is not echoed in the terminal; and when the user asks to see it, the agent shows the dalgo or native form.

### AC: web-edit-is-structured (verifies REQ:web-deterministic-edits)

**Given** a tab with a current query shown in the Web UI
**When** the user adds a filter via the Web UI controls
**Then** the change is sent as a structured operation on the canonical AST, not as prose.

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

**Given** a tab with a current query
**When** the user asks to "delete all inactive users"
**Then** the request is refused, no mutating query is produced, and the tab's current query is unchanged.

### AC: matrix-tracks-surfaces (verifies REQ:parity-matrix)

**Given** the Capability README
**When** a reader inspects it
**Then** an `## Implementation Matrix` is present with one row per surface, each carrying a Status, a one-line Brief, and a Link to that surface's realization.

## Rehearse Integration

Every AC has a concrete behavioral surface (daemon state, MCP tools, structured web edits, sync, refusal paths, document structure) and is testable in principle. However, this Capability defines cross-surface *contract* behavior with no executable surface of its own — each AC is exercised against a concrete Implementation (the CLI daemon in `datatug-cli`, the Web features in `datatug-apps`), where the Rehearse stubs will live. Stub scaffolding is therefore deferred to the Implementation Features so the stubs track real, runnable surfaces rather than an abstract Capability.

## Not Doing / Out of Scope

- In-app AI-chat panel on any surface — the agent terminal is the only conversational surface; the Web UI is deterministic controls, not prose.
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

---
*This document follows the https://specscore.md/feature-specification*
