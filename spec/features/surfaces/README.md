---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Three surfaces, one project

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/surfaces?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/surfaces?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/surfaces?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/surfaces?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P11 is the "three surfaces, one project" claim: TUI, web (local agent or hosted) and AI agents all work over the same project files, with screen-level hand-off between them, and no surface holds private state the others cannot see.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

Initial audit snapshot (before later 9 September landings): the TUI is usable for schema browsing (Queries/Dashboards panels remain placeholders); the web UI through the local agent is "prototype (broken both ends)" — `datatug serve` panics before it listens and the web client is missing the server's `/datatug` URL prefix at 8+ call sites; the hosted web UI just landed enterprise OIDC/SAML SSO but the one canonical demo-project link is broken; Ctrl+W hand-off maps exactly one route; the agent skills plugin is usable (8 of 9 skills verified); MCP is specified, zero code. Phase 0 repairs the agent connection and hides six empty-shell pages; Phase 2 extends hand-off and TUI parity; Phase 4 depends on a working `serve` for MCP. This is a Capability; per-repo Implementations (`datatug-cli`, `datatug-apps`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> DataTug's differentiation depends on every surface reaching the same project, the same semantic
> layer, the same library — not three separate products that happen to read the same files. Today
> the connection between the two most-used surfaces does not exist: "every web-client call to
> `datatug serve` is missing the server's hard-coded `/datatug` URL prefix... a mechanical, total
> mismatch across 8+ call sites, not a typo in one place" (current-state audit, verdict line 3), and
> `datatug serve` itself panics before it opens a socket (`cmd_serve.go:43`, confirmed live). The
> "click here for the demo" path the hosted web UI markets as the on-ramp is broken independently: a
> `No provider found for _DatatugFoldersService` error and a 404 against the wrong repository name
> (current-state audit, capability table, "Live demo-project navigation" row). Six routed web pages
> (Widgets/Tags/Resources/Variables/DB-Models/Diff) are one click away and completely empty. Ctrl+W
> hand-off from the TUI to the web only maps the `"projects"` screen; every other TUI screen falls
> back to the bare web origin (`webui_handoff.go:26-28`).

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| TUI (tview) screens | usable (Queries/Dashboards panels are placeholders) | 2 | strategic |
| Web UI through the local agent (`datatug serve`) | prototype (broken both ends) | 0 | enabler |
| Hosted web UI (datatug.app) with Sneat identity and enterprise SSO | partial | 3 | strategic |
| Screen-level hand-off CLI → web (Ctrl+W) and web → CLI | partial (one route mapped) | 2 | strategic |
| Agent skills plugin | usable (8 of 9 verified) | 2 | strategic |
| MCP face / serve-brokered query builder | specified/documented (Approved), zero code | 4 | strategic |
| Dead pages hidden until built (Widgets, Tags, Resources, Variables, DB models, Diff) | prototype shells | 0 | hygiene |

#### REQ: tui-schema-browsing

The TUI MUST let a user browse databases, environments, entities and settings against real, live
schema data, as it does today. Its Queries and Dashboards panels remain hard-coded placeholders
until Phase 2 binds them to real data; this REQ does not claim they are done.

#### REQ: agent-connection-fixed

The web UI and `datatug serve` MUST agree on one route table (the `/datatug` prefix and endpoint
names in `pkg/server/endpoints/routes.go`), verified by a contract test against the real server.
This is specified in full by `core-investigation-loop` REQ:agent-path-contract; this Feature owns
it as the standing P11 guarantee that the connection, once fixed, must not regress.

#### REQ: hosted-web-ui-sneat-identity

The hosted web UI on `datatug.app` MUST authenticate through Sneat's shared identity, including
enterprise OIDC/SAML SSO where configured, and its canonical demo-project link MUST resolve to a
real, current project rather than a 404 or a missing-provider error.

#### REQ: screen-level-handoff

A user MUST be able to move from a TUI screen to the equivalent web screen (Ctrl+W) and from a web
screen back to the equivalent CLI/TUI context, for every screen the loop covers — not only the
`"projects"` screen the current mapping supports. Hand-off MUST preserve project, table/query and
selection state, matching the URL-carries-state, never-credentials guarantee in `core-investigation-
loop`'s epilogue branches.

#### REQ: agent-skills-plugin

The AI agent skills plugin MUST keep wrapping real, present CLI subcommands; a skill that depends on
unbuilt infrastructure (the `query-builder` skill, which needs the unbuilt MCP daemon) MUST be
marked as such rather than presented as equally ready. 8 of 9 skills are already verified against
real commands; this REQ is the standing bar new skills must also clear.

#### REQ: mcp-serve-brokered-query-builder

The MCP face of the surfaces story is delivered by the separate, already-Approved Feature
`serve-brokered-query-builder` (a daemon brokering one canonical query per tab between an MCP face
and an HTTP/WS face) — this Feature links to it as a related Capability rather than restating its
contract, and requires only that it depend on a working `datatug serve` (REQ:agent-connection-fixed)
before it can ship.

#### REQ: hide-dead-pages

A web page with no real data behind it (Widgets, Tags, Resources, Variables, DB models, Diff) MUST
be hidden from navigation, behind a feature flag, until it is built. A page reachable by a click
MUST NOT be an empty shell a user can stumble into.

#### REQ: product-profiles-and-workspace

The web surface MUST be one build that presents itself as the `datatug`, `incidentius` or
`dashboardius` product profile by hostname or configuration — one configurable app for the
whole family (founder, 2026-09-11: *"one configurable app."*) — and that app MUST live in the
`datatug/datatug-apps` Nx workspace (founder, 2026-09-11, quoted in
[product-profiles](../product-profiles/README.md)). Profiles change home, brand, navigation and
entry points only; every screen, the TUI hand-off and the agent skills reach the same routes
under either profile. The mechanism itself is specified in `product-profiles`; this REQ pins
the P11 guarantee that a profile never becomes a private surface.

## Acceptance Criteria

### AC: profile-does-not-hide-a-surface (verifies REQ:product-profiles-and-workspace)

**Given** the app served under the `incidentius` profile
**When** a TUI hand-off link for the *Customer* table and an agent skill deep link are opened
**Then** both resolve to the same screens they reach under the `datatug` profile, with the Incidentius home and branding.

### AC: tui-browses-real-schema (verifies REQ:tui-schema-browsing)

**Given** the demo project opened in the TUI
**When** the user browses `Customer`'s columns and foreign keys
**Then** real, live-scanned schema data renders, while the Queries and Dashboards panels remain visibly placeholder until Phase 2.

### AC: agent-path-resolves (verifies REQ:agent-connection-fixed)

**Given** `datatug serve` is running
**When** the web client calls any endpoint it uses
**Then** every call resolves under `/datatug/`, matching `core-investigation-loop` AC:agent-path-contract-client.

### AC: hosted-demo-link-opens-the-project (verifies REQ:hosted-web-ui-sneat-identity)

**Given** the hosted web UI's "DataTug Demo Project" entry point
**When** a user clicks it
**Then** the current demo project opens, not a `No provider found` error or a 404 against the wrong repository name.

### AC: handoff-preserves-selection (verifies REQ:screen-level-handoff)

**Given** the TUI's *Customer* table screen with row 5 highlighted
**When** the user presses Ctrl+W
**Then** the web UI opens the same table with the same table opens; selecting a sensitive row requires explicit safe-state transfer and authorization rather than silently serializing its value into a URL.

### AC: skill-marks-unready-dependency (verifies REQ:agent-skills-plugin)

**Given** the `query-builder` skill
**When** its `SKILL.md` is inspected
**Then** it states its dependency on the unbuilt MCP daemon rather than presenting itself as equally usable to the 8 verified skills.

### AC: mcp-face-depends-on-working-serve (verifies REQ:mcp-serve-brokered-query-builder)

**Given** `serve-brokered-query-builder`'s Implementation Matrix
**When** its dependency is checked
**Then** it names a working `datatug serve` (REQ:agent-connection-fixed) as a precondition, and this Feature does not restate the daemon's own REQs.

### AC: empty-pages-hidden (verifies REQ:hide-dead-pages)

**Given** a fresh web UI build
**When** the main navigation is inspected
**Then** Widgets, Tags, Resources, Variables, DB models and Diff are not reachable by any click until each has real data behind it.

## Not Doing / Out of Scope

- Any UI-kit migration: this Feature does not decide Ionic-vs-PrimeNG; that decision lives in
  `ROADMAP.md`'s UI-kit section and applies uniformly once made.
- Building the MCP daemon itself: that is `serve-brokered-query-builder`'s and `ai-query-builder`'s
  scope, linked here, not restated.
- Real-time multi-user collaboration between surfaces; hand-off is a one-user, one-session
  transition, not shared live state.
- A no-code app builder inside any surface; every surface presents the same data/knowledge model,
  not a screen-building tool.

## Open Questions

- The founder selected Ionic frame plus bounded islands for DataTug; Dashboardius
  remains Angular/PrimeNG. This is not awaiting another UI-kit confirmation.
- Clarify repository visibility versus license restrictions for "Web UI private for
  now" before changing visibility or selecting a future open-source license.
- Before hosted publication, verify one Cloudflare deploy owner/path and authenticated
  principal contract, with candidate-specific journey gates.

---
*This document follows the https://specscore.md/feature-specification*
