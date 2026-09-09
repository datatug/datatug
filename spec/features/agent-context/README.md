---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Agent-readable project context

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/agent-context?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/agent-context?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/agent-context?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/agent-context?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P14 is the claim that AI agents use the same model people do: the same entities, mappings, library queries, investigation context and access decisions, exposed through CLI verbs and MCP resources rather than a disconnected chat architecture (deliberate non-goal: "a text-to-SQL chat box...

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

AI enters through the same knowledge model as people"). This Feature, `agent-context`, is P14's own capability — a resolver an agent can call to learn what a human would see in the context panel — distinct from the two existing, separately-specified Approved Features under this same pillar: `ai-query-builder` (in-app conversational query building) and `serve-brokered-query-builder` (a daemon brokering a query between an MCP face and a web UI). Initial audit snapshot (before later 9 September landings): the skills plugin genuinely drives the CLI (8 of 9 skills verified against real commands); agent-readable project context as a CLI/MCP resource is planned, unbuilt; both sibling Features are specified/documented with zero shipped code. Phase 2 builds `datatug context`, reusing the same resolver `core-investigation-loop` builds for the web UI. This is a Capability; per-repo Implementations (`datatug-cli`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> An AI agent guessing SQL from a raw physical schema repeats every mistake a human would make
> without the semantic layer, and cannot benefit from a team's accumulated library or from an active
> Investigation Context a human started. DataTug's skills plugin already proves the model works for
> CLI-wrapping skills — "the skills plugin (`ai-plugin`, 9 skills, 8 wrapping real commands) is
> genuine" (current-state audit, ecosystem findings) — but nothing today lets an agent ask "what does
> this value mean, what's related, what queries apply" the way `core-investigation-loop`'s context
> panel answers that for a human. Without this Feature, an agent's only options are the same
> unstructured schema dump every text-to-SQL tool already offers, or waiting for the still-unbuilt MCP
> daemon (`serve-brokered-query-builder`) — which itself "depends on a `datatug serve` that cannot
> currently start" (current-state audit, capability table, MCP row).

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| Skills plugin driving the CLI | usable | — | strategic |
| Agent-readable project context (entities, mappings, applicable queries) as a CLI/MCP resource | planned | 2 | strategic |
| Serve-brokered query builder (MCP + HTTP/WS faces) | specified/documented | 4 | strategic |
| In-app AI query builder | specified/documented | X | strategic |

#### REQ: skills-plugin-drives-cli

The AI agent skills plugin MUST keep wrapping real, present CLI subcommands, staying at or above
the 8-of-9-verified bar the audit found. A skill MUST NOT be published against a command that does
not exist or does not work (matching `surfaces` REQ:agent-skills-plugin, which owns this same
guarantee from the surfaces side).

#### REQ: agent-readable-context-resource

An agent MUST be able to ask, for a semantic value or a set of values, what a human would see in
the context panel: the entity/field meaning, related records, and applicable library queries with
their resolution chain. This MUST reuse the exact resolver `core-investigation-loop` builds for the
web UI (REQ:semantic-resolution-endpoint, REQ:related-lookup-model, REQ:applicable-queries) —
exposed as a `datatug context` CLI verb and, later, an MCP resource — not a separate, agent-specific
reimplementation that could drift from what a human sees.

#### REQ: agent-context-respects-access-control

Every read `datatug context` performs on behalf of an agent MUST pass through the same
`accesspolicies.Run` path as a human's read (`server-enforced-access` REQ:serve-acl-all-reads). An
agent MUST NOT see a mapping, related record, or applicable query a human under the same principal
could not see.

#### REQ: link-serve-brokered-query-builder

The MCP-brokered, in-terminal query-building experience is owned by the separate, Approved Feature
`serve-brokered-query-builder` — this Feature links it as a related P14 Capability and does not
restate its contract. `serve-brokered-query-builder` depends on `surfaces` REQ:agent-connection-
fixed before it can ship, and MAY consume this Feature's context resource once both exist.

#### REQ: link-ai-query-builder

The in-app conversational query-building experience is owned by the separate, Approved Feature
`ai-query-builder` — this Feature links it as a related P14 Capability and does not restate its
contract. It remains an experiment behind the serve-brokered path per the roadmap, not a competing
default.

## Acceptance Criteria

### AC: skills-still-wrap-real-commands (verifies REQ:skills-plugin-drives-cli)

**Given** the `ai-plugin` skills directory
**When** each `SKILL.md` is checked against the current CLI command tree
**Then** at least 8 of 9 skills wrap a command that actually exists and works, matching the audited baseline.

### AC: context-verb-matches-web-panel (verifies REQ:agent-readable-context-resource)

**Given** `Customer.ID = 5` in the demo project
**When** `datatug context --value Customer.ID=5` runs
**Then** it lists the same related records and applicable queries (with resolution chain) that `core-investigation-loop` AC:related-across-sources-server and AC:applicable-with-chain-server return for the same value through the web endpoint.

### AC: agent-context-blocked-by-policy (verifies REQ:agent-context-respects-access-control)

**Given** `datatug context` invoked under the `support` principal
**When** it resolves context for a Brazilian customer under a Canada-only row policy
**Then** it returns the same "rows filtered" limitation a human `support` session would see, never an unfiltered agent-only view.

### AC: serve-brokered-builder-not-restated (verifies REQ:link-serve-brokered-query-builder)

**Given** this Feature's Behavior section
**When** it is compared against `serve-brokered-query-builder`'s REQs
**Then** no REQ here duplicates the daemon's tab/mode/session contract; the two Features remain separately owned and cross-linked.

### AC: ai-query-builder-not-restated (verifies REQ:link-ai-query-builder)

**Given** this Feature's Behavior section
**When** it is compared against `ai-query-builder`'s REQs
**Then** no REQ here duplicates the in-app conversational session contract; the two Features remain separately owned and cross-linked.

## Not Doing / Out of Scope

- A text-to-SQL chat box: agent access to context is structured (entities, mappings, applicable
  queries), not a free-form natural-language-to-SQL translator.
- Redefining `ai-query-builder` or `serve-brokered-query-builder`'s own contracts; this Feature
  links to both without restating them.
- Agent-only access paths that bypass server-enforced access control; an agent is a principal like
  any other.
- A disconnected agent architecture that reads the physical schema directly instead of the semantic
  layer.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
