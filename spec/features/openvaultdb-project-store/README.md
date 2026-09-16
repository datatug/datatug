---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: OpenVaultDB project store

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/openvaultdb-project-store?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/openvaultdb-project-store?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/openvaultdb-project-store?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/openvaultdb-project-store?op=request-change) |
**Status:** Approved
**Date:** 2026-09-16
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Depends-On:** git-native-projects, server-enforced-access, secrets-vault
**Grade:** A

## Summary

Implement `git-native-projects` REQ:project-store-pluggable-backends by wiring an
OpenVaultDB-backed project store, so a DataTug project can live in storage its owner
controls — a GitHub repository reached through OpenVaultDB — while the hosted web app
keeps working against it.

Founder instruction, 2026-09-09, verbatim: *"It should be possible to store and use
datatug project from remote OpenVaultDB."* The pieces already exist on both sides:
`dalgo2openvaultdb` is a working DALgo driver over OpenVaultDB's `/v1` API, and DataTug's
`storage.Store`/`ProjectStore` interfaces exist. What has never been done is wiring them
together for a real project.

This Feature covers the store itself: one repository binding per project, the
`datatug/projects/<id>/` namespace, server-mediated writes, and a browser read path.
Credentials are explicitly **not** part of it — they belong to the `secrets-vault`
Feature, and project files carry logical pointers only.

## Problem

The initial audit recorded that the local filestore is the only fully working `Store`
backend, that GitHub-hosted projects can write but not read back
(`GhStorage.FileExists`/`OpenFile` panic), and that the store abstraction has a single
real implementation despite four named backends. Meanwhile the knowledge a DataTug user
accumulates is exactly the asset they would most want to keep in storage they own.

Two constraints shape the solution and are inherited rather than invented here:

- `git-native-projects` already forbids client-direct writes to a protected project
  store: *"writes to a hosted or remote store still go through the trusted server"*. So
  writes are server-mediated. Reads have no such constraint, and the web app needs them
  to be cheap.
- `server-enforced-access` already owns project/source/row/column enforcement through
  DALgo, and OpenVaultDB's own policies compose with it rather than replacing it.

The consequence is an asymmetric path: writes traverse the trusted server, reads go from
the browser straight to the store. That asymmetry is deliberate and its limits are stated
in REQ:browser-read-path rather than left implicit.

## Behavior

### Store binding and layout

#### REQ: project-store-binding

A project MUST bind exactly one OpenVaultDB store, identified by an OpenVaultDB endpoint,
a database id, and a credential held server-side. The binding MUST be recorded in the
authoritative registry at `users/{uid}/ext/datatug` as an entry carrying a `role` field,
stored as a list even while only one binding exists, so that additional bindings are
additive and never require a data migration. The legal `role` values are `project`, the
project store, and `secrets`, reserved for a future vault binding; an unknown role value
MUST be rejected. DataTug MUST hold an OpenVaultDB credential
for this purpose and MUST NOT hold a GitHub credential for the project store.

Vault bindings are a separate relation owned by `secrets-vault` and MUST NOT be modelled
as project-store bindings: a project may bind several vaults (for example one for DEV and
UAT and one for PROD) while still holding exactly one project-store binding, because a
vault owns its own storage.

#### REQ: project-namespace-layout

Project assets MUST live under `/datatug/projects/<project-id>/…` inside the bound store,
using OpenVaultDB's application-namespace profile: a root `/.ingitdb/root-collections.yaml`
mapping a `projects` collection under `datatug/.collections/`, with `data_dir` resolved
relative to the parent of `.collections/` so the record subtree is `datatug/projects/`.
A single store (repository) MAY hold several DataTug projects, and DataTug MUST preserve
files it does not manage, byte for byte. This is the same resolution rule `secrets-vault`
relies on for `datatug/projects/<project-id>/credentials/<server-id>/`; the two Features
MUST describe one layout, not two. The layout MUST be documented so a person or an
agent can read a project by reading files.

### Read and write paths

#### REQ: server-mediated-writes

Every mutation of a project store MUST go through the DataTug server
(`/datatug/v0/…`) and write through `dalgo2openvaultdb`. A browser or AI agent MUST NOT
hold a write credential for the project store, and DataTug MUST NOT write project files
by any path that bypasses the server.

#### REQ: browser-read-path

A browser MUST be able to read a project store directly from the store's git host over its
REST API, using a read credential the user supplies, through a reader module shipped in
`datatug-apps`. The reader MUST be conformance-tested against the writer's output, so that
the Go writer and the browser reader cannot drift apart silently. Because this path
bypasses the server, DataTug MUST document its limits plainly: a reader credential that
can reach `/datatug/projects/<id>/…` can reach anything else that credential covers, so
read-scoped restrictions expressed in DataTug MUST NOT be advertised as enforced on this
path. This is the one declared exception to `server-enforced-access`'s rule that reads
reach a source through the trusted server; see REQ:project-store-read-exception.

#### REQ: project-store-read-exception

Browser-direct reads of the project store are a declared, narrowly scoped exception to
`server-enforced-access`, which forbids client-direct access to protected data sources. The
exception MUST cover reads of the project store and nothing else: every write MUST still go
through the trusted server, and every read of a protected data source MUST still go through
it. `server-enforced-access` MUST be amended so that its prohibition is scoped to protected
data sources rather than to every client connection, and MUST point at this exception. This
Feature MUST NOT be read as widening client-direct access anywhere else.

### Registry and store choice

#### REQ: registry-is-authoritative

The Firestore registry at `users/{uid}/ext/datatug` MUST remain the authoritative list of
a user's DataTug projects, and DataTug MUST NOT maintain a derived project index inside
the OpenVaultDB store. Listing projects MUST NOT depend on a read-modify-write of shared
store state, because OpenVaultDB's MVP offers no cross-transaction isolation and no update
preconditions.

#### REQ: store-choice-is-explicit

A project's store type MUST be chosen when the project is created and recorded with the
project. Storing one project in two stores at once MUST NOT be possible, and moving a
project between stores MUST be an explicit migration operation rather than an implicit
dual write.

### Credential hygiene

#### REQ: no-inline-credentials

A project asset MUST reference a connection by logical name and MUST NOT carry a
credential value. `datatug validate` MUST reject a project where any asset carries a
credential, and `Project.Validate()` MUST validate `Queries` — which the audit found it
never does. This guard is owned by `git-native-projects`
REQ:secrets-never-in-project-files, with which `heterogeneous-sources`
REQ:secrets-outside-project-files shares one implementation; this Feature MUST enforce that
shared guard and MUST NOT define a third, independent copy of it. Credentials themselves
are stored by the `secrets-vault` Feature.

#### REQ: no-vault-inside-a-client-readable-store

A project store MUST NOT be configured so that a client read credential can reach a bound
vault's storage. Because a vault lays its credentials out under
`/datatug/projects/<project-id>/credentials/<server-id>/…` using the same namespace profile
as the project store, colocating a bound vault in a repository the client reads directly
would silently defeat `secrets-vault` REQ:client-cannot-read-credentials. Binding
validation MUST reject such a placement rather than warn and continue.

## Acceptance Criteria

### AC: binding-is-recorded-as-a-roled-list (verifies REQ:project-store-binding)

**Given** a user creating a project against a GitHub repository through OpenVaultDB
**When** the project is created
**Then** the registry at `users/{uid}/ext/datatug` holds a binding entry naming the OpenVaultDB endpoint, database and `role: project`, the entry is stored as a list, and no GitHub credential is stored by the server for that project.

### AC: assets-land-under-the-project-namespace (verifies REQ:project-namespace-layout)

**Given** a project bound to an OpenVaultDB store whose repository already contains unrelated files
**When** an entity and a query are added
**Then** they appear under `datatug/projects/<project-id>/`, `/.ingitdb/root-collections.yaml` maps the `projects` collection under `datatug/.collections/`, and the pre-existing unrelated files are unchanged byte for byte.

### AC: two-projects-share-one-store (verifies REQ:project-namespace-layout)

**Given** two DataTug projects bound to the same repository
**When** both projects are written
**Then** each project's assets appear under its own `datatug/projects/<project-id>/` subtree and neither write alters the other project's files.

### AC: a-client-cannot-write-the-store (verifies REQ:server-mediated-writes)

**Given** a signed-in browser session for a project bound to an OpenVaultDB store
**When** the client attempts to mutate project files without going through `/datatug/v0/…`
**Then** the mutation does not reach the store, and the same mutation submitted through the server succeeds.

### AC: browser-reads-a-project-without-the-server (verifies REQ:browser-read-path)

**Given** a project written to a repository through the server, and a user-supplied read credential for that repository
**When** the browser opens the project with the server unavailable
**Then** the project's assets are read from the repository over its REST API, and the parsed result matches what the conformance test derives from the Go writer's output.

### AC: project-list-comes-from-the-registry (verifies REQ:registry-is-authoritative)

**Given** a store containing projects written by more than one user
**When** a user opens their project list
**Then** the list is served from the registry at `users/{uid}/ext/datatug`, listing performs no write to the OpenVaultDB store, and no derived project index exists inside that store to be read or updated.

### AC: store-switch-is-an-explicit-migration (verifies REQ:store-choice-is-explicit)

**Given** a project bound to the DataTug cloud store
**When** the owner moves it to an OpenVaultDB store
**Then** the move is performed as a single explicit migration, the project is bound to exactly one store afterwards, and no write is ever directed at both stores.

### AC: validator-rejects-credentials-in-project-assets (verifies REQ:no-inline-credentials)

**Given** a query file in a project store that carries an inline password
**When** `datatug validate` runs
**Then** it fails naming the file and field, and passes once the credential is replaced by a logical connection name.

### AC: server-authenticates-with-the-recorded-credential (verifies REQ:project-store-binding)

**Given** a project bound to an OpenVaultDB store with a server-side credential recorded in the registry
**When** the server reads the project store and then mutates it
**Then** both operations authenticate with the recorded OpenVaultDB credential, succeed while it is valid, and fail closed when it is absent or rejected, with no GitHub credential involved in either.

### AC: project-store-reads-are-the-only-exception (verifies REQ:project-store-read-exception)

**Given** a protected data source and a project store
**When** a client attempts a direct read of each, and a write to each
**Then** the direct read of the protected data source is refused, the direct read of the project store succeeds, and both writes are refused unless they arrive through the trusted server.

### AC: vault-cannot-be-colocated-in-a-client-readable-store (verifies REQ:no-vault-inside-a-client-readable-store)

**Given** a project store in a repository the client reads with a user-supplied credential
**When** a vault binding is configured so that a credential the client holds can reach that vault's `credentials` collection
**Then** DataTug rejects the binding, no credential is written to that location, and the rejection names the colocation as the reason.

## Rehearse Integration

Every AC has a concrete behavioural surface — registry records, store paths, HTTP
requests, validator exit behaviour, and a reader-versus-writer comparison — and all are
testable in principle. This Feature defines a capability spanning several repositories,
and each AC is exercised against a concrete implementation: server-mediated writes and the
store binding in `datatug/backend`, the browser read path in `datatug-apps`, the validator
in `datatug-cli`, and the OpenVaultDB data plane in `openvaultdb-go`. Stub scaffolding is
therefore deferred to those implementations, so the stubs track real runnable surfaces
rather than an abstract capability — consistent with `ai-query-builder` and
`serve-brokered-query-builder`. One AC carries an obligation that must survive that
deferral: AC:browser-reads-a-project-without-the-server is a conformance test between the
Go writer and the browser reader, and it must be authored where both sides are runnable. Concretely: the binding and credential ACs
land on `datatug/backend` HTTP and credential-resolution surfaces; AC:assets-land-under-the-project-namespace
and AC:two-projects-share-one-store on the store's file system; AC:project-list-comes-from-the-registry
on the Firestore registry; AC:browser-reads-a-project-without-the-server on the
`datatug-apps` reader plus its conformance test; AC:validator-rejects-credentials-in-project-assets
on `datatug-cli` exit behaviour; and AC:vault-cannot-be-colocated-in-a-client-readable-store
on binding-validation rejection.

## Not Doing / Out of Scope

- Per-environment repositories and more than one store binding per project. The binding is
  a list so this is additive later, but it is not built here.
- Credentials, connection strings and any other secret: owned by `secrets-vault`, never by
  this Feature's store.
- The platform-wide `/ext/v0/…` API path convention. The OpenVaultDB data plane's public
  URL depends on it, so it is tracked as a platform Feature in the platform repo rather
  than duplicated here.
- A full `openvaultdb-ts` SDK. The browser reader is deliberately minimal and disposable.
- Client-direct writes to the project store, which `git-native-projects` already excludes.
- Direct browser-to-database query execution.

## Open Questions

- The canonical public URL of the OpenVaultDB data plane depends on the platform path
  convention Feature; until that lands, the endpoint shape used by
  REQ:project-store-binding is provisional.

---

*This document follows the https://specscore.md/feature-specification*
