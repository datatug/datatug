---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Sharing scopes on Sneat spaces

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/sharing-scopes?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P8 is sharing scopes: personal → team → project → organization scopes for knowledge assets, reusing Sneat's existing space membership and role model rather than inventing a parallel one.

It is the enabler P7's knowledge-promotion lifecycle depends on (hypothesis H5, "confirmed for the model, not for the code" — Sneat's space membership/roles exist, DataTug does not use them yet). Audited state on 2026-09-09: not implemented. `datatug_projects` uses a flat `userIDs[]` field rather than Sneat's space-bounded convention, `@sneat/space-models`/`@sneat/space-services` are declared dependencies with zero import sites, and `ICreateProjectData.spaceID` is a typed field that is never populated or read. Phase 3 builds the scope model and project-membership wiring; it is also the nearest monetizable wedge (a hosted team workspace) per the backstage B2B strategy. This is a Capability; per-repo Implementations (`datatug-apps`, `sneat-go`) link to it.

## Problem

A knowledge library only compounds if the right people can find what other people captured, at the
right granularity — a personal scratch query should not clutter a team's shared library, but a
proven investigation should promote cleanly from personal to team to the whole organization.
DataTug already sits on Sneat's production Firebase project and could reuse Sneat's space
membership/role model directly, but instead keeps "its own parallel `datatug_projects` membership
model rather than Sneat's space-bounded convention" (current-state audit, ecosystem-claims
findings). The founder's standing rule is that "almost everything in Sneat is space bounded" — a
new entity's records belong under `/spaces/{spaceID}/ext/<ext>/` and get their ownership from that
space's member-role model, not a hand-rolled owner field. `datatug_projects`' flat `userIDs[]`
(`sneat-firebase/firestore.rules:71-114`) is exactly the pattern that rule warns against, and the
`spaceID` field DataTug already declared on `ICreateProjectData` sits unused
(`project.service.ts:193-197`).

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Personal / team / project / organization scopes | planned | 3 | strategic |
| Project membership through Sneat space membership and roles | planned | 3 | enabler |
| Hosted team workspace on datatug.app | partial | 3 | strategic |

#### REQ: personal-team-project-org-scopes

A knowledge asset (query, collection, playbook) MUST be assignable to exactly one of personal,
team, project or organization scope at a time, and promotable between them per `executable-
knowledge-library` REQ:knowledge-promotion-lifecycle. Scope MUST determine visibility and
discoverability, never execution rights — access policy (P10) remains the separate, authoritative
gate on what a principal can actually read or run.

#### REQ: project-membership-through-sneat-spaces

Project membership MUST be derived from Sneat space membership and roles, not a hand-rolled
`userIDs[]` field. `datatug_projects` records MUST move under `/spaces/{spaceID}/ext/datatug/` per
the space-bounded convention, with `ICreateProjectData.spaceID` actually read and enforced rather
than declared and ignored. This is a founder decision to confirm before implementation (see Open
Questions).

#### REQ: hosted-team-workspace

A hosted team workspace on `datatug.app` MUST let a Sneat space's members open, browse and run
queries against shared DataTug projects using their existing Sneat identity — the nearest
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

**Given** a Sneat space with an admin and a viewer role
**When** each opens the same DataTug project under that space
**Then** their DataTug project permissions match their Sneat space role, with no separate `userIDs[]` entry required.

### AC: spaceid-field-actually-used (verifies REQ:project-membership-through-sneat-spaces)

**Given** a new project created with `ICreateProjectData.spaceID` set
**When** the project record is written
**Then** it is stored under `/spaces/{spaceID}/ext/datatug/` and readable by that space's members, not as a flat top-level `datatug_projects` document.

### AC: hosted-workspace-uses-sneat-identity (verifies REQ:hosted-team-workspace)

**Given** a user already signed in to a Sneat space on `datatug.app`
**When** they open a shared project in the hosted team workspace
**Then** no separate DataTug-specific sign-in is required, and their access matches their space role.

## Not Doing / Out of Scope

- A generic workflow engine or app builder: sharing scopes govern visibility of knowledge assets,
  not application logic.
- Client-direct database access for protected sources: scope changes visibility, never bypasses
  server-enforced access control.
- A data catalog's governance/approval workflow: promotion changes discoverability, not a
  certification status.
- Inventing a DataTug-specific membership or role system in parallel with Sneat's space model.

## Open Questions

- Should `datatug_projects` move to the space-bounded model
  (`/spaces/{spaceID}/ext/datatug/`), or keep a project-level `userIDs[]` with a space link? This
  is a founder decision (per the space-bounded standing rule, the recommended answer is the
  space-bounded model, but it is not yet a founder ruling for this repository).

---
*This document follows the https://specscore.md/feature-specification*
