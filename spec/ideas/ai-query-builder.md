---
format: https://specscore.md/idea-specification
status: Specified
---

# Idea: AI Query Builder

**Status:** Specified
**Date:** 2026-06-04
**Owner:** alexander.trakhimenok
**Promotes To:** ai-query-builder
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we let DataTug users build and progressively refine a query through a chain of plain-English requests — each one transforming the *current* query (add a column, add a filter, change ordering) against live connections and results — as one in-app capability defined once and realized per-surface (CLI and Web) so the surfaces stay honestly tracked instead of silently drifting?

## Context

DataTug ships the same product on two surfaces in separate repos — a CLI (`datatug-cli`, mature SpecScore spec tree) and a Web UI (`datatug-apps`, an Nx workspace with no SpecScore tree yet). An AI-assisted query-builder feature (originally framed as "AI chat") is wanted on both. SpecScore's `capability-and-platform-implementations` Feature (`github.com/specscore/specscore`, Status: Implementing, Grade A) exists for exactly this: a platform-agnostic Capability owns an `## Implementation Matrix`, and each per-surface Implementation links upward via `**Implements:** specscore:feature/<slug>@github.com/datatug/datatug`. `specscore` CLI 0.6.0 is on PATH and lints these references. This Idea is the first real dogfood of that pattern across three repos: Capability in `datatug/`, CLI Implementation in `datatug-cli/`, Web Implementation in `datatug-apps/`.

**Why this must be in-app.** Authoring DataTug project metadata (entities, datasets, environments) can be delegated to *any* external AI agent wielding DataTug skills, so it does not need to live in DataTug's own UI. Building a query is different: it is bound to live connections, the current result set, and a turn-by-turn feedback loop — so it must be an in-app experience. That is the line that earns this feature its place in the product, and it narrows the scope to an **AI-assisted query builder**.

**Why progressive refinement is the crux.** The job is not one prompt → one query. It is a chain of requests that each transform the query the user is already looking at:

1. "Show me last 10 users who logged into the app" — seed query.
2. "Add a column that shows exchange rate for the user's primary currency to EUR." — add a (computed/joined) column.
3. "Hide users who made an order." — add a filter (an anti-join).

The current query is **mutable session state**; each utterance is a delta against it, and the user sees the updated query and refreshed results before the next turn. The examples also span easy→hard: steps 2 and 3 imply cross-source/computed columns and anti-joins, which is where MVP limitations will honestly land.

## Recommended Direction

Model the AI query builder as ONE Capability with TWO Implementations using SpecScore's capability/implementation pattern — not as two independent features. The Capability lives in `datatug/` and defines the platform-agnostic job: an **in-app conversational query builder** that maintains a *current query* as session state and applies each natural-language request as a transformation against it (add/remove/modify columns, filters, joins, ordering, limits), grounded in the project's connections and schema, with the query and its results visible and re-runnable each turn. The Capability's acceptance criteria describe only what must hold on every surface, and it owns the Implementation Matrix tracking CLI and Web parity. The two Implementations specify only their per-surface deltas and limitations.

The refinement model is the part worth specifying carefully. Each turn maps an utterance to a *delta on the current query* — not a regeneration from scratch — so "add a column" and "hide users who made an order" edit the intended part of the query and leave the rest intact. This stateful loop, not one-shot NL→SQL, is what makes it a query *builder* rather than a prompt box, and it is the behavior the dogfood must prove.

The honest MVP is CLI-first internal dogfood. The Capability declares the full refinement vocabulary; the **CLI Implementation** in `datatug-cli/` nails the core loop on one connection — seed a query from NL, render results in the existing CLI table UI, then apply follow-up requests as transformations (at minimum add-column and add-filter, the simple cases) re-running and re-rendering each time. Operations that require cross-source or computed columns (the EUR exchange-rate example) and anti-joins (the order filter) may ship as a documented `Partial`. The Web Implementation is `Planned`. The cycle-one success signal is that the example refinement chain (or a comparable one) runs end-to-end in the CLI *and* that the cross-repo capability/implementation plumbing holds and lints clean — proving both the in-app builder and the SpecScore pattern at lowest cost. The CLI already has a mature spec tree and renders result tables, so it validates both cheaply; Web leads later, once the Capability's ACs have survived one real Implementation.

## Alternatives Considered

**Metadata-authoring copilot as the in-app feature** (the previously-considered direction). Rejected: project-metadata authoring can be delegated to any external AI agent using DataTug skills, so it does not justify an in-app surface. Query building is bound to live connections, results, and iteration, and does — so the in-app bet is the query builder, not a general copilot.

**One-shot NL→SQL with no session state.** Rejected: a stateless prompt box cannot express the progressive refinement in the motivating examples, where each request edits the query the user is already looking at. The mutable-current-query model is the whole point; without it this is just SQL autocomplete.

**Two independent per-repo "AI chat" features, or Web-first / both-surfaces-parallel.** Rejected: independent features re-create the cross-surface drift the Capability pattern exists to kill, with no single source of truth. Web-first or parallel costs more (datatug-apps needs SpecScore bootstrap and a richer UX) and dilutes the dogfood goal. CLI leads at lowest cost; Web is `Planned`. (Single-repo co-location of all three docs — the pattern's special case — is also rejected because the surfaces genuinely live in separate repos.)

## MVP Scope

One cycle, dogfooded internally. (1) Author the AI Query Builder **Capability** in `datatug/` — platform-agnostic ACs for a stateful conversational query builder (a *current query* held as session state, per-turn natural-language transformations applied to it, and visible/re-runnable results) plus an Implementation Matrix with CLI and Web rows. (2) Author and build the **CLI Implementation** in `datatug-cli/` (`**Implements:** specscore:feature/ai-query-builder@github.com/datatug/datatug`) that seeds a query from a natural-language prompt against one project connection, renders results in the existing CLI table UI, and applies follow-up requests as transformations — at minimum add-column and add-filter — re-running and re-rendering each turn; operations needing cross-source/computed columns or anti-joins may ship `Partial`. (3) Mark the Web row `Planned`. Success = the motivating example chain (or a comparable one) runs end-to-end in the CLI, AND the cross-repo `**Implements:**` link resolves and `specscore spec lint` passes in both repos. Embarrassing-but-shippable beats a polished two-surface build.

## Not Doing (and Why)

- Web Implementation build-out this cycle — marked `Planned` in the matrix; CLI leads to validate the pattern at lowest cost.
- Metadata/project authoring as part of this feature — explicitly excluded; it can be delegated to any external AI agent using DataTug skills and so needs no in-app surface. This Idea is the in-app query builder.
- Guaranteeing complex operations in the MVP — cross-source/computed columns (the EUR exchange-rate example) and anti-join filters (the order example) may ship `Partial`; the MVP guarantees only single-source add-column / add-filter / ordering refinements.
- Executing mutating or DDL statements from chat — the builder produces read-only queries; write-path safety is its own design problem.
- A multi-LLM-provider abstraction layer — pick one provider path for MVP; provider-neutrality is premature until the refinement loop proves out.
- New SpecScore Document Kinds — reuse the existing Capability/Implementation roles on the Feature Kind.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | The progressive-refinement model is reliable — the assistant can apply a follow-up request ("add column X", "hide users who made an order") as a delta on the current query without mangling the rest of it. | Dogfood the example chain plus variations. Acceptance: each follow-up modifies the intended part of the query and leaves the rest intact, more often than not, on a real connection. |
| Must-be-true | The cross-repo `**Implements:**` reference from `datatug-cli` to the `datatug/` Capability is authorable and lints clean with `specscore` 0.6.0 (the capability pattern is still `Implementing` upstream). | Author the Capability + CLI Implementation; run `specscore spec lint` in both repos. Acceptance: the same-repo half resolves, the cross-repo notation is accepted, and a target rename produces a clear error. |
| Should-be-true | The "current query" can be represented at the Capability level in a way both surfaces share (a query model/AST or canonical SQL) without leaking CLI-specific UX. | Define the Capability's notion of *current query* abstractly; confirm the CLI Implementation and a hypothetical Web Implementation could both realize it from the same ACs. |
| Should-be-true | An in-app builder beats delegating query-building to an external AI agent + DataTug skills — the live connection/result feedback loop is worth building in-product. | Dogfood for one cycle. Acceptance: the team reaches for the in-app refinement loop over copy-pasting between an external agent and a query runner. |
| Might-be-true | Complex operations (computed/joined columns like the FX-rate example, anti-join filters like the order example) can be added to the same refinement loop later without re-architecting the simple-case MVP. | Sketch those two transformations against the MVP's current-query representation; confirm they are additive, not a redesign. |

## SpecScore Integration

- **New Features this would create:** an `ai-query-builder` **Capability** in `datatug/` (platform-agnostic in-app query builder; owns the Implementation Matrix); an `ai-query-builder` **CLI Implementation** in `datatug-cli/` (`**Implements:** specscore:feature/ai-query-builder@github.com/datatug/datatug`, complex joins/anti-joins possibly `Partial`); a future `ai-query-builder` **Web Implementation** in `datatug-apps/` (initially `Planned`).
- **Existing Features affected:** `datatug-cli`'s query/execute/console/UI features (`cli/queries`, `cli/execute`, `cli/console`, `cli/ui`) — the builder reuses the existing query-run + result-render machinery as a new conversational mode; DataTug's query model and connection/schema models (the grounding and the current-query representation).
- **Dependencies:** SpecScore's `capability-and-platform-implementations` Feature (`github.com/specscore/specscore`, Implementing) — this Idea is its first cross-repo dogfood; `specscore` CLI ≥ 0.6.0 for `**Implements:**` lint; the CLI's existing query execution + result rendering, reused rather than reinvented.

## Open Questions

- What is the canonical representation of the "current query" the refinement loop edits — generated SQL re-parsed each turn, a DataTug query model/AST, or accumulated NL instructions re-synthesized? This determines how reliably deltas apply and how the state shares across surfaces.
- Where is the MVP line between a guaranteed and a `Partial` operation — single-source column/filter/order (guaranteed) vs cross-source joins and computed columns (Partial)?
- Which LLM provider/path does the CLI MVP use (BYO cloud API key vs local model), and where do credentials live relative to the DataTug project config?
- Which repo owns the Capability long-term — `datatug/` is chosen here — and how is that ownership made discoverable from the Implementations? (Mirrors an open question in the upstream capability Feature.)

---
*This document follows the https://specscore.md/idea-specification*
