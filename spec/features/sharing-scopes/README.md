---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Sharing scopes on Sneat spaces

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P8 is sharing scopes: personal → team → project → organization scopes for knowledge assets, reusing Sneat's existing linkage system to assign Space contacts to standalone projects with project roles. For shared/protected nonpublic projects, Space membership does not automatically confer access or supply a project role; current platform guards still apply. Public-project and owner-only private access retain their existing policies.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

It is the enabler P7's knowledge-promotion lifecycle depends on (hypothesis H5, plausible integration with an explicit standalone-project boundary — Sneat's space membership/roles exist, DataTug does not use them yet). Initial audit snapshot (before later 9 September landings): not implemented. `datatug_projects` uses a flat `userIDs[]` field rather than Sneat's space-bounded convention, `@sneat/space-models`/`@sneat/space-services` are declared dependencies with zero import sites, and `ICreateProjectData.spaceID` is a typed field that is never populated or read. Phase 4 builds the scope model and project-membership wiring; it is also the nearest monetizable wedge (a hosted team workspace) per the backstage B2B strategy. This is a Capability; per-repo Implementations (`datatug-apps`, `sneat-go`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> A knowledge library only compounds if the right people can find what other people captured, at the
> right granularity — a personal scratch query should not clutter a team's shared library, but a
> proven investigation should promote cleanly from personal to team to the whole organization.
> DataTug already sits on Sneat's production Firebase project and could reuse Sneat's space
> membership/role model directly, but instead keeps "its own parallel `datatug_projects` membership
> model rather than Sneat's space-bounded convention" (current-state audit, ecosystem-claims
> findings). The founder's standing rule is that "almost everything in Sneat is space bounded" — a
> new entity's records belong under `/spaces/{spaceID}/ext/<ext>/` and get their ownership from that
> space's member-role model, not a hand-rolled owner field. `datatug_projects`' flat `userIDs[]`
> (`sneat-firebase/firestore.rules:71-114`) is exactly the pattern that rule warns against, and the
> `spaceID` field DataTug already declared on `ICreateProjectData` sits unused
> (`project.service.ts:193-197`).

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| Personal / team / project / organization scopes | planned | 4 | strategic |
| Project membership through Sneat space membership and roles | planned | 4 | enabler |
| Hosted team workspace on datatug.app | partial | 4 | strategic |

#### REQ: personal-team-project-org-scopes

Proposed future promotion model (Draft, not a founder ruling): a knowledge asset may be assigned to personal,
team, project or organization visibility, and promoted between them per `executable-
knowledge-library` REQ:knowledge-promotion-lifecycle. Scope MUST determine visibility and
discoverability, never execution rights — access policy (P10) remains the separate, authoritative
gate on what a principal can actually read or run.

#### REQ: project-membership-through-sneat-spaces

Founder words: "DataTug project is standalone entity that can be shared (granted)
between multiple spaces (teams)". Project identity MUST remain standalone; sharing
with another Space MUST NOT duplicate the project.

The later founder clarification supersedes automatic Space-role inheritance:
"5 contacts linked to project with some role ... viewer or contributor ... Up to 5
can be linked ... not about space members ... users who were assigned to project
as space contact with one of few predefined roles"; "We have linkage system".
Shared/protected nonpublic project assignments MUST reuse existing Sneat contact-to-project linkage carrying
predefined project roles, not a new `userIDs[]` or parallel membership list. A
Space member without an assigned project role MUST NOT gain shared/protected
nonpublic project access from membership alone. Public-project access and
owner-only private policy MUST remain unchanged; the Pro contact cap MUST NOT
limit public projects or add contact assignment to their access requirements. Project roles MUST matter to authorization while current
contact status, authenticated-user binding and platform guards are rechecked.
The confirmed role permissions are: Viewer may browse, run queries and edit
locally, but MUST NOT save shared changes; Contributor may do all Viewer
actions and save shared changes; Owner may do all Contributor actions and
manage project assignments. These permissions remain subject to current
contact-to-authenticated-user binding, platform guards and applicable plan
entitlements. An Owner role MUST NOT bypass a platform or paid-plan guard; when
a plan ends, shared writes and hosted AI remain denied while local query
execution remains possible.
The [pricing feature](../pricing-and-billing/README.md) governs Pro's five distinct
linked-contact cap across roles, owner included, including unregistered contacts;
this is not a five-member Space limit or a Business MAU count. Exact persisted
linkage adaptation, owner-contact binding, cross-Space contact identity mapping
and legacy spaceID migration remain design questions.

#### REQ: hosted-team-workspace

For shared/protected nonpublic projects, a hosted team workspace on `datatug.app`
MUST let users whose current contact binding and
assigned project role authorize the operation open, browse and run shared DataTug queries
using their existing Sneat identity — the nearest
monetizable wedge per the backstage B2B strategy (2026-07-01). It builds on
REQ:project-membership-through-sneat-spaces; it MUST NOT ship a separate, parallel authentication
or membership system.

## Acceptance Criteria

### AC: scope-assignment-controls-visibility (verifies REQ:personal-team-project-org-scopes)

**Given** a query captured at personal scope
**When** a teammate in the same project browses the shared library
**Then** the query is not visible to them until it is promoted to team scope or higher.

### AC: promotion-does-not-bypass-acl (verifies REQ:personal-team-project-org-scopes)

**Given** a query promoted to organization scope
**When** a principal without read access to its underlying source opens it
**Then** the query is discoverable but its run is still refused by server-enforced access control (P10), proving scope and ACL are separate gates.

### AC: project-membership-follows-space-role (verifies REQ:project-membership-through-sneat-spaces)

**Given** a shared/protected nonpublic project with contacts assigned Viewer, Contributor and Owner through existing linkage, plus a Space administrator without a project assignment
**When** their authenticated users browse project data, run a query, edit it locally, save a shared change or manage assignments
**Then** Viewer can browse, run and edit locally but cannot save shared changes; Contributor can also save shared changes; Owner can also manage assignments. These actions require the current contact-to-authenticated-user binding and all applicable platform and plan guards; an ended plan denies shared writes and hosted AI while allowing local query execution, including for Owner. A Space role alone neither supplies a project role nor grants project access, and no separate `userIDs[]` membership list is introduced. Viewing metadata remains non-billable under the pricing feature's activity rule. This criterion supersedes the earlier Space-role inheritance wording; its retained identifier is historical.

### AC: spaceid-field-actually-used (verifies REQ:project-membership-through-sneat-spaces)

**Given** a shared/protected nonpublic project created with `ICreateProjectData.spaceID` set
**When** the project record is written
**Then** the project remains a standalone entity with explicit sharing/linkage to that Space; only contacts whose assigned project role and current binding authorize reading may read it. Contacts linked through a second Space access the same project identity without a copy; a spaceID field or Space membership alone is not a project assignment.

### AC: hosted-workspace-uses-sneat-identity (verifies REQ:hosted-team-workspace)

**Given** a user already signed in to a Sneat space on `datatug.app`
**When** they open a shared/protected nonpublic project in the hosted team workspace
**Then** no separate DataTug-specific sign-in is required, and access requires their current contact-to-user binding and assigned project role with platform guards; signing in or being a Space member alone does not grant project access. An unregistered contact can count as an assignment without thereby gaining authenticated access. Public projects retain their existing access policy and unlimited-user edit policy; owner-only private projects retain owner-only access.

### AC: reciprocal-linkage-enforces-project-policy (verifies REQ:project-membership-through-sneat-spaces)

**Given** a protected project and both project-side and contact-side role-link mutation paths, including generic linkage endpoints
**When** an assignment exceeds its plan cap, changes a role without authority or bypasses a paid-end guard
**Then** the same policy refuses on either side without a partial reciprocal update. Count each distinct assigned contact once across roles and registration status; still-linked inactive or missing contacts do not silently release capacity. The existing generic endpoint's availability does not prove these product guards are implemented.

## Not Doing / Out of Scope

- A generic workflow engine or app builder: sharing scopes govern visibility of knowledge assets,
  not application logic.
- Client-direct database access for protected sources: scope changes visibility, never bypasses
  server-enforced access control.
- A data catalog's governance/approval workflow: promotion changes discoverability, not a
  certification status.
- Inventing a parallel membership list or linkage engine instead of reusing Sneat's contact/project role links; no automatic project assignment or Space membership creation is implied by a link.

## Open Questions

- Confirm the personal/team/project/organization promotion ladder and its relationship
  to project grants before implementing promotion; current founder words establish
  standalone multi-space project sharing, not this complete visibility model.
- Define owner-contact binding within the cap, cross-Space contact identity mapping,
  persisted linkage adaptation and legacy spaceID migration against existing contracts.
  The confirmed role permissions do not grant a cap exemption.

- Select the authenticated hosting/deployment profile before the collaboration slice;
  keep DataTug projects standalone and compose explicit contact/project role links
  with current platform guards.
- Define conflict/revocation and contact archive/delete/unlink cleanup with the existing
  linkage system. Missing or inactive contacts that remain assigned still count; access
  must reread current contact status and authenticated-user binding. Prove product guards
  for generic mutations from both sides before any protected-sharing activation.

---
*This document follows the https://specscore.md/feature-specification*
