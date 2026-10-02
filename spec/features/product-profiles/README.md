---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Product profiles

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/product-profiles?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/product-profiles?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/product-profiles?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/product-profiles?op=request-change) |
**Status:** Approved
**Source Ideas:** —

## Summary

A product profile is the presentation layer that turns an application in the workspace into
DataTug, Incidentius or Dashboardius: home route, brand, navigation order, entry points and
demo defaults, selected by hostname or configuration. One build serves all three profiles
(founder, 2026-09-11: *"one configurable app."*); the frame is Ionic everywhere, and PrimeNG,
AG Grid, Chart.js and CDK drag-drop appear only as leaf components where Ionic lacks the
widget.

A profile changes what a user sees first and nothing else — every route stays reachable, every
deep link works across profiles, and no profile holds private data or private authority.

## Problem

The founder's requirement (vision §4: *"Build capabilities once. Expose them through
different product shells."*, restated in vision §61) is one application with three shells:

> **Build capabilities once. Expose them through DataTug, Incidentius and Dashboardius product
> profiles.** (vision §61)

> Do not create a separate Incidentius frontend/repository.

The obvious failure mode is three implementations that drift: three navigation models, three
board renderers, three auth flows, three build pipelines, and eventually three data models.
The second failure mode is subtler — a "profile" that quietly becomes a tenant, where the
Incidentius shell can see something the DataTug shell cannot, and a deep link sent by a
colleague opens a 404 because the recipient's hostname is different.

Today the repository half is already split: `datatug/datatug-apps` is an Nx workspace, and
`sneat-co/dashboardius` is a standalone Angular/PrimeNG app that deploys itself to
`dashboardius.com` through its own `web/worker.js`, `wrangler.jsonc` and the shared
`sneat-co/cicd` `cf-deploy.yml` workflow. The founder ruled that split away for code, and the
move is a cutover with a real inventory, not a file copy.

```text
  https://datatug.app         https://incidentius.…         https://dashboardius.com
          │                            │                             │
          └──────── profile resolution (hostname / configuration) ───┘
                                  │
                                  ▼
                    apps/datatug-app — one Ionic-framed build
              profile: datatug | incidentius | dashboardius
              home    = project | incidents   | boards
                                  │
        shared libs: @datatug/board-models, board renderer (chart, grid and
        drag-drop leaf widgets), agent client, auth, build-info
                                  │
                        one server, one model, one policy
```

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Profile table: home, brand, navigation, entry points, demo defaults | planned | MVP (B1) | strategic |
| Hostname selection with config/env override and `?profile=` for local dev | planned | MVP (B1) | enabler |
| Every route reachable in every profile; deep links cross profiles | planned | MVP (B1) | strategic |
| One build serves every profile | planned | MVP (B1) | enabler |
| Declarative profile configuration (adding a profile is a data change) | planned | MVP (B1) | enabler |
| Profile selection confers no authority | planned | MVP (B1) | strategic |
| Brand and domain honesty (no unregistered domain claimed) | planned | MVP (B1) | surface |
| Build-info footer per profile (`@sneat/components` AppVersionComponent) | exists in Sneat apps; adoption planned | MVP (B1) | enabler |
| Dashboardius shell rebuilt in Ionic inside the single app (cutover) (full cutover) | planned | NEXT | enabler |
| Dashboardius-owned board renderer library embedded as an island | planned | NEXT | strategic |
| Per-profile marketing sites and demo content | planned | NEXT | surface |

#### REQ: profile-table

The family MUST define three profiles — `datatug`, `incidentius`, `dashboardius` — each
declaring: home route, brand (name, logo, colours, copy), navigation order, entry points, and
demo defaults. One build of `apps/datatug-app` serves all three (founder, 2026-09-11: *"one
configurable app."*). The Dashboardius home is an Ionic page; its cards, lists, toolbars,
menus and forms are Ionic components; PrimeNG, AG Grid, Chart.js and CDK drag-drop remain
only as leaf widgets inside the board renderer library, with their theme CSS scoped to those
components and never loaded globally.

| | `datatug` | `incidentius` | `dashboardius` |
|---|---|---|---|
| Home route | the project / investigation workspace | the incidents list | the boards list |
| Brand name | DataTug | Incidentius | Dashboardius |
| Primary entry point | "open a project" | **"Houston, we've got a problem"** | "new board" |
| Navigation order | sources · queries · context · incidents · boards | incidents · context · queries · sources · boards | boards · queries · projects |
| Demo default | the demo knowledge project | the canonical demo incident | the demo board |

The entry-point wording for `incidentius` is the founder's own (vision §23: *"Houston, we've got
a problem"*), and the *Incidents* item in the `datatug` side menu is fixed by the founder
(2026-09-11, verbatim: *"DataTug app should always have Incidents menu item in side menu."*):
it MUST be present under the `datatug` profile. Lead assumptions, not the founder's words:
it is present whether or not a project is open, and no feature flag or project setting removes
it. The rest of the navigation orders and the demo defaults are lead assumptions. Brand colours, logo
and copy per profile are lead assumptions pending the brand decision in Open Questions.

A profile MUST be selected by **hostname**, with an environment/configuration override for
deployments that do not have a distinct hostname, and a `?profile=` query parameter accepted
**only** when the app is running in a local development configuration. Precedence (lead
assumption): explicit configuration override, then hostname match, then the `datatug` default.
An unknown hostname MUST fall back to `datatug` and MUST NOT fail to boot.

#### REQ: no-profile-private-data

A profile changes presentation only. Every route MUST remain reachable in every profile, and a
deep link MUST work when opened under a different profile than the one that produced it — the
recipient may see a different chrome, never a 404 and never a redirect that loses the target.

No data, project, capability or API is profile-private: the server does not know which profile
a request came from, and the profile MUST NOT be a parameter of any `/datatug/*` call. Two
users looking at the same incident under two profiles are looking at the same asset under the
same grants. The consequence to protect deliberately: an incident opened from the Incidentius
home and the same incident reached from a DataTug project page are the same page, not two
implementations of one screen.

#### REQ: one-workspace-shared-libs

All three apps MUST live in the `datatug/datatug-apps` Nx workspace. Founder ruling,
2026-09-11 (verbatim):

> *"Not sure if it's part of the spec your read but I think DataTug, Incidentius, Dashbaordius
> apps code should be in same repo (nx workspace with shared libs and either separate or single
> configurable app)"*

This settles the repository; the founder settled the second half on the same day (*"one
configurable app."*, and on the frame: *"We should migrate PrimeNG cards and other components
that we can migrate to Ionic and use PrimeNG just as components if Ionic misses
functionality"*). The shared libraries MUST carry the reuse: `@datatug/board-models`
(exists, in `libs/datatug/board-models`), a **Dashboardius-owned board renderer library** that
DataTug/Incidentius pages embed as an island under the founder's UI-kit rule (ruled 2026-09-09,
verbatim: *"UI-kit rule - confirmed"* — Ionic frame with PrimeNG/AG Grid islands only where
Ionic lacks components), the agent client, auth, and build-info.

No app may fork a shared library, and no app may introduce a second board renderer
(`dashboards` REQ:no-charts-in-datatug-web-ui), a second agent client or a second auth path.
Dashboardius's shell is rebuilt with Ionic components, and PrimeNG, AG Grid, Chart.js and CDK
drag-drop remain only as leaf widgets in the renderer library (founder, 2026-09-11); the
roadmap's UI-kit section carries the same amendment.

#### REQ: dashboardius-app-cutover

Moving `sneat-co/dashboardius/web` into the workspace MUST be a **full cutover**, not an
addition (rule `cutover-verbs-mean-full-cutover`), and — since the founder ruled *"one
configurable app."* (2026-09-11) — a **rebuild of the shell, not a move**: Dashboardius
becomes the `dashboardius` profile of `apps/datatug-app`; its layout, cards, lists, toolbars,
menus and forms are rebuilt with Ionic components (founder, 2026-09-11: *"We should migrate
PrimeNG cards and other components that we can migrate to Ionic and use PrimeNG just as
components if Ionic misses functionality"*); PrimeNG, AG Grid, Chart.js and CDK drag-drop
survive only as leaf widgets inside the board renderer library, with their theme CSS scoped
to those components and never loaded globally; the PrimeNG shell in `sneat-co/dashboardius`
is retired. Scroll and gesture behaviour of AG Grid and CDK drag-drop inside `ion-content`
is an acceptance criterion, not an assumption. The deployment inventory, read from that
repository on 2026-09-11, is:

| Artefact | Path today | Cutover action |
|---|---|---|
| Cloudflare Worker | `web/worker.js` | move; re-point at the workspace build output |
| Wrangler config | `web/wrangler.jsonc` | move; keep the `dashboardius.com` binding and the dry-run validation step |
| Deploy call | `.github/workflows/ci.yml` → `sneat-co/cicd/.github/workflows/cf-deploy.yml@main` with `working-directory: web`, `build-command`, `wrangler-config`, `smoke-command` | re-point `working-directory` and `build-command` at the workspace app; keep the org `CLOUDFLARE_API_TOKEN` secret and `CLOUDFLARE_ACCOUNT_ID` variable |
| Playwright e2e | `web/playwright.config.ts`, `web/e2e/{auth,landing,prerender,responsive}.spec.ts` | move all four specs; keep the prerender assertion running against the built HTML |
| Post-deploy smoke | `web/scripts/post-deploy-smoke.mjs https://dashboardius.com` | move; keep the live-domain smoke |
| PrimeNG palette and licence | `web/scripts/generate-primeng-palette.mjs`, `web/scripts/write-primeui-license.mjs`, `PRIMEUI_LICENSE` secret | move both scripts; keep the licence optional so its absence is a visible banner, not a build failure |
| PrimeNG shell: pages, layout, cards, menus, forms | `web/src/app/**` | **rebuild** with Ionic components in `apps/datatug-app` under the `dashboardius` profile (founder, 2026-09-11); retire the PrimeNG shell |
| PrimeNG theme CSS | `web/src/styles*` (loaded globally today) | scope to the renderer's leaf components; never loaded globally |
| OG image generation | `web/scripts/generate-og-image.mjs` | move |
| Renovate | `.github/renovate.json` | fold into the workspace's Renovate configuration |

The old app path MUST be removed once the workspace app is live, and the live
`https://dashboardius.com` deployment MUST come from the workspace before removal — a cutover
is finished when nothing is served from the old path, not when the new one builds. Whether
`sneat-co/dashboardius` remains the Dashboardius **spec hub** (as this repository is for
DataTug) or is archived is a founder decision; the lead recommendation is to keep it as a spec
hub, and no spec content moves in this cutover either way.

#### REQ: brand-and-domain-honesty

The approved Incidentius layout uses `incidentius.com` for the public landing page and
`app.incidentius.com` for the application. The apex MUST remain a marketing surface and MUST
NOT select the application profile. The exact app hostname MUST select the `incidentius`
profile; `datatug.app` and `dashboardius.com` continue to select their own profiles. Public
copy MUST describe the app as available only once the app hostname resolves with valid TLS
and serves the Incidentius profile. Configuration override and `?profile=incidentius` remain
local-development mechanisms, not production aliases. Incidentius's brand tier relative to
DataTug remains open.

Every profile MUST show build information — version, git hash, build timestamp — reusing the
existing `@sneat/components` `AppVersionComponent` rather than a per-profile implementation,
so "which build am I looking at?" is answerable in all three shells the same way.

#### REQ: profile-config-is-declarative

The profile table MUST live in one declarative configuration consumed by the apps, not in
scattered conditionals. Adding or changing a profile MUST be a data change to that
configuration plus, at most, registering a new home route — never a new branch in shared
component code. A component MAY read a profile token (home route, brand name, entry-point
label) but MUST NOT branch on the profile identity to change behaviour, because that is how
three implementations grow back inside one build.

Profile resolution MUST happen once at application bootstrap and be exposed as an injectable
token (lead assumption), so tests can instantiate any profile without a hostname and a screen
can be rendered under every profile token that build serves in one test run.

#### REQ: profile-selection-is-not-authorization

Selecting a profile MUST NOT grant anything. The hostname, the configuration override and the
`?profile=` development parameter change presentation only: they never change the principal,
never change which projects are listed, never change policy evaluation, and never appear in an
authorization decision. The server MUST reject a profile field on any request that could
influence execution or access, the same way the transport contract rejects a client-supplied
principal or role.

This is stated separately from REQ:no-profile-private-data because the two failure modes
differ: that one is about reachability, this one is about a URL parameter becoming a privilege.
`?profile=` in particular MUST be inert outside local development — not because it is dangerous
today, but so that it never becomes the thing somebody later tries to gate a feature on.

## Acceptance Criteria

### AC: hostname-selects-home (verifies REQ:profile-table)

**Given** one deployed build served on `datatug.app` and `app.incidentius.com`, while
`incidentius.com` serves the public landing page
**When** an authenticated user opens each application root URL
**Then** `datatug.app` lands on the project/investigation workspace and
`app.incidentius.com` lands on the incidents list with the "Houston, we've got a problem"
entry point, each with its own brand name and navigation order; `incidentius.com` remains the
public landing surface, and an unrecognised hostname lands on the DataTug home rather than
failing to boot.

### AC: incidents-always-in-datatug-menu (verifies REQ:profile-table)

**Given** the build served under the `datatug` profile, once with a project open and once with no project open
**When** the side menu renders in each state
**Then** an *Incidents* item is present both times and opens the incident list, and no feature flag or project setting offered by the app can remove it.

### AC: deep-link-works-in-other-profile (verifies REQ:no-profile-private-data)

**Given** a link to an incident detail page produced under the `incidentius` profile
**When** an authorized user opens it under the `datatug` profile
**Then** the same incident renders with the same data and the same grants, with DataTug chrome — no 404, no redirect to a profile home, and no capability missing from the page.

### AC: one-build-serves-both-profiles (verifies REQ:profile-config-is-declarative)

**Given** the CI pipeline for `apps/datatug-app`
**When** the DataTug, Incidentius and Dashboardius deployments are produced (the AC slug predates the one-app ruling and is kept for traceability)
**Then** all three are served from the same build artefact with different configuration, the pipeline runs no second build for a second profile, and a test renders a shared screen under each of the three profile tokens without a hostname.

### AC: dashboardius-shell-is-ionic (verifies REQ:dashboardius-app-cutover)

**Given** the single app served under the `dashboardius` profile
**When** the boards list and a board page render
**Then** the frame, cards, lists, toolbars, menus and forms are Ionic components, PrimeNG, AG Grid, Chart.js and CDK drag-drop appear only inside board renderer widgets, no global PrimeNG theme stylesheet is loaded, and a board with a grid widget scrolls and a card drags correctly inside `ion-content`.

### AC: dashboardius-app-in-workspace-cutover-complete (verifies REQ:dashboardius-app-cutover)

**Given** the cutover of `sneat-co/dashboardius/web` into `datatug/datatug-apps`
**When** the cutover is declared complete
**Then** `https://dashboardius.com` is served by the workspace build through the moved Worker and `wrangler.jsonc` with the `cf-deploy.yml` call re-pointed, all four Playwright specs and the post-deploy smoke run from the workspace, the PrimeNG palette and licence scripts are in the workspace with `PRIMEUI_LICENSE` still optional, Renovate covers the moved app, and the old `web/` app path no longer exists or deploys. The Dashboardius shell renders as Ionic components with PrimeNG, AG Grid, Chart.js and CDK drag-drop only inside renderer widgets, and no global PrimeNG theme stylesheet is loaded (AC:dashboardius-shell-is-ionic).

### AC: profile-override-grants-nothing (verifies REQ:profile-selection-is-not-authorization)

**Given** a user whose principal cannot read project `acme`
**When** they open the app with `?profile=incidentius` and with each hostname in turn
**Then** the listed projects, the evaluated policy and every server response are identical in all cases, `acme` is absent throughout, no request carries a profile field that reaches execution or access, and `?profile=` is ignored outside local development.

## Not Doing / Out of Scope

- Tenancy. A profile is not an organisation, a workspace or a data boundary; project grants
  (`sharing-scopes`) are the only boundary.
- Per-profile feature flags, entitlements or paywalls. Pricing separates viewing from
  contributing (product-family D13), not one profile from another.
- Resolving Incidentius's brand tier.
- Marketing sites. `datatug/backstage` (`marketing/`) owns the public story; a profile owns in-app copy.
- Moving Dashboardius **specs**. Only the app moves in the cutover; the spec-hub question is a
  separate founder decision.

## Open Questions

- **Incidentius brand tier.** The split-domain layout is decided. Tier-3 independent like
  DataTug, or a DataTug capability with its own public identity?
- **Whether `sneat-co/dashboardius` is kept as the Dashboardius spec hub or archived** after the
  app cutover. Lead recommendation: keep it.
- **Profile-aware demo content.** Each profile declares a demo default, but the canonical demo
  incident and the demo knowledge project may or may not be the same project. Lead assumption:
  one demo project carries both stories; if Chinook cannot carry the incident story
  ([incidents](../incidents/README.md) Open Questions), the Incidentius profile's demo default
  points at a second project and the profile table absorbs that with no code change.

---
*This document follows the https://specscore.md/feature-specification*
