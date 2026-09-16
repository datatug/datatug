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
`datatug/projects` collection namespace, server-mediated writes, and a browser read path.
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

A project MUST bind exactly one OpenVaultDB store, identified by a **user-supplied
OpenVaultDB host base URL** and a database id. The binding MUST be recorded in the
authoritative registry at `users/{uid}/ext/datatug` as an entry carrying a `role` field,
stored as a list even while only one binding exists, so that additional bindings are
additive and never require a data migration. The legal `role` values are `project`, the
project store, and `secrets`, reserved for a future vault binding; an unknown role value
MUST be rejected. The OpenVaultDB credential DataTug uses to reach that store MUST be held
server-side and MUST NOT be carried by the registry entry
(REQ:ovdb-credential-is-server-only). DataTug MUST hold an OpenVaultDB credential
for this purpose and MUST NOT hold a GitHub credential for the project store.

Vault bindings are a separate relation owned by `secrets-vault` and MUST NOT be modelled
as project-store bindings: a project may bind several vaults (for example one for DEV and
UAT and one for PROD) while still holding exactly one project-store binding, because a
vault owns its own storage.

#### REQ: project-namespace-layout

Project assets MUST live in the DataTug application namespace expressed by the **collection
name** `datatug/projects`, declared in the bound OpenVaultDB manifest, and DataTug MUST
address a project's assets within that collection by the project's id. This Feature MUST NOT
mandate an on-disk directory tree: the exact layout a backend derives from that declaration
is the backend's business, and no requirement here may assert a path a named backend does
not produce.

The two named backends realise the same declaration differently:

- The local inGitDB backend resolves the mapping from a root
  `/.ingitdb/root-collections.yaml` entry, whose `data_dir` is resolved relative to the
  parent of `.collections/`.
- The OpenVaultDB GitHub backend maps the declared collection name directly to its directory
  (`DirPath: name` in `buildInGitDBDefinition`,
  `openvaultdb-go/pkg/mount/ingitdb_github.go`), so the same declaration yields the
  `datatug/projects/…` namespace, and the collection's schema lives in the bound manifest's
  `schemas.collections`. DataTug MUST NOT require that repository to also carry
  `/.ingitdb/root-collections.yaml` or a `datatug/.collections/` schema directory; on that
  backend those clauses are no-ops.

Physical nesting is an upstream prerequisite, not an assumption. `schemas.collections`
entries carry only `fields` — there is no path, subtree or `data_dir` field on a collection
(`openvaultdb-go/pkg/schema/schema.go`) — and `buildInGitDBDefinition` fixes each
collection's directory to its name. A manifest declaration therefore cannot by itself place
a collection in a nested subtree or give each project id its own directory. DataTug MUST NOT
assume a project's assets physically occupy `datatug/projects/<project-id>/`, and MUST NOT
assert such a path until `openvaultdb-go` supplies **per-collection path support honoured by
`buildInGitDBDefinition`**; that support is an explicit upstream prerequisite of this
Feature, shared with `secrets-vault` REQ:vault-mirrors-project-namespace. Until it lands, the
namespace is satisfied by the `datatug/projects` collection name plus the project's id as its
record key.

A single store (repository) MAY hold several DataTug projects, and DataTug MUST preserve
files it does not manage, byte for byte. `secrets-vault` addresses its credentials through
the same rule — a declared collection name carrying the static prefix, plus identity segments
as record keys — so the two Features MUST describe one layout rule, not two mechanisms. The
layout MUST be documented so a person or an agent can read a project by reading files.

### Read and write paths

#### REQ: server-mediated-writes

Every mutation of a project store MUST go through the DataTug server
(`/v0/datatug/…`, the prefix `datatug/backend` registers in
`api4datatug/routes.go`) and write through `dalgo2openvaultdb`. A browser or AI agent MUST
NOT hold a write credential for the project store, and DataTug MUST NOT write project files
by any path that bypasses the server.

#### REQ: browser-read-path

A browser MUST be able to read a project store directly from the store's git host over its
REST API, using a read credential the user supplies, through a reader module shipped in
`datatug-apps`. The reader MUST be conformance-tested against the writer's output, so that
the Go writer and the browser reader cannot drift apart silently. Because this path
bypasses the server, DataTug MUST document its limits plainly: a reader credential that
can reach the store's `datatug/projects` collection can reach anything else that credential
covers, so
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

#### REQ: ovdb-credential-is-server-only

The OpenVaultDB credential DataTug uses to reach a project's bound store MUST be held in a
**server-only** credential store and MUST NEVER be written to the client-readable registry at
`users/{uid}/ext/datatug` or into any client bundle. The registry entry there records the
binding — host base URL, database id and `role` — and nothing else; a credential in that document
is a defect, because the client reads the registry to list projects. A missing or rejected
credential MUST fail closed, naming the missing credential, rather than falling back to
anonymous or default access. This mirrors, for the project store's data plane, the
prohibition `secrets-vault` REQ:server-authorizes-vault-access states for vault
authorization, so the two Features are consistent.

#### REQ: no-vault-inside-a-client-readable-store

A project store MUST NOT be configured so that a client read credential can reach a bound
vault's storage. Because a vault addresses its credentials through the same collection-name
namespace rule as the project store and keys them per project and per endpoint
(`secrets-vault` REQ:vault-mirrors-project-namespace), colocating a bound vault in a
repository the client reads directly would silently defeat `secrets-vault`
REQ:client-cannot-read-credentials. Binding validation MUST reject such a placement rather
than warn and continue.

## Acceptance Criteria

### AC: binding-is-recorded-as-a-roled-list (verifies REQ:project-store-binding)

**Given** a user creating a project against a GitHub repository through OpenVaultDB
**When** the project is created
**Then** the registry at `users/{uid}/ext/datatug` holds a binding entry naming the OpenVaultDB host base URL, database and `role: project`, the entry is stored as a list, the entry carries no OpenVaultDB credential, and no GitHub credential is stored by the server for that project.

### AC: assets-land-under-the-project-namespace (verifies REQ:project-namespace-layout)

**Given** a project bound to an OpenVaultDB store whose repository already contains unrelated files
**When** an entity and a query are added
**Then** they are written as records of the collection named `datatug/projects` declared in the bound manifest and keyed by the project's id — declared in the manifest's `schemas.collections` on the OpenVaultDB GitHub backend, and mapped by a `/.ingitdb/root-collections.yaml` `data_dir` entry on the local inGitDB backend — and the pre-existing unrelated files are unchanged byte for byte. On the OpenVaultDB GitHub backend, no `/.ingitdb/root-collections.yaml` and no `datatug/.collections/` schema directory is required in the repository.

### AC: two-projects-share-one-store (verifies REQ:project-namespace-layout)

**Given** two DataTug projects bound to the same repository
**When** both projects are written
**Then** each project's assets are records of the `datatug/projects` collection keyed by that project's id, the two projects' records are disjoint, and neither write alters the other project's files.

### AC: a-client-cannot-write-the-store (verifies REQ:server-mediated-writes)

**Given** a signed-in browser session for a project bound to an OpenVaultDB store
**When** the client attempts to mutate project files without going through `/v0/datatug/…`
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

**Given** a project bound to an OpenVaultDB store whose OpenVaultDB credential is held only in the server-only credential store and not in the registry
**When** the server reads the project store and then mutates it
**Then** both operations authenticate with that server-only credential, succeed while it is valid, and fail closed when it is absent or rejected, with no GitHub credential involved in either.

### AC: ovdb-credential-never-appears-in-the-registry (verifies REQ:ovdb-credential-is-server-only)

**Given** a server bound to a project's OpenVaultDB store and holding that store's OpenVaultDB credential
**When** the client-readable registry at `users/{uid}/ext/datatug` and every client bundle are inspected
**Then** neither carries the OpenVaultDB credential, which lives only in the server-only credential store, and a read attempted against the registry fails closed rather than yielding the credential.

### AC: project-store-reads-are-the-only-exception (verifies REQ:project-store-read-exception)

**Given** a protected data source and a project store
**When** a client attempts a direct read of each, and a write to each
**Then** the direct read of the protected data source is refused, the direct read of the project store succeeds, and both writes are refused unless they arrive through the trusted server.

### AC: vault-cannot-be-colocated-in-a-client-readable-store (verifies REQ:no-vault-inside-a-client-readable-store)

**Given** a project store in a repository the client reads with a user-supplied credential
**When** a vault binding is configured so that a credential the client holds can reach that vault's `credentials` collection
**Then** DataTug rejects the binding, no credential is written to that location, and the rejection names the colocation as the reason.

## Rehearse Integration

Every AC has a concrete behavioural surface — registry records, store records, HTTP
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
and AC:two-projects-share-one-store on the store's collection records; AC:project-list-comes-from-the-registry
on the Firestore registry; AC:browser-reads-a-project-without-the-server on the
`datatug-apps` reader plus its conformance test; AC:validator-rejects-credentials-in-project-assets
on `datatug-cli` exit behaviour; and AC:vault-cannot-be-colocated-in-a-client-readable-store
on binding-validation rejection. One obligation is explicitly **not** this Feature's:
`git-native-projects` AC:github-project-reads-back-what-it-wrote depends on `GhStorage`
read methods that still panic, so no stub authored here can satisfy it and it must not be
attached to this Feature's plan as a Supporting AC.

## Not Doing / Out of Scope

- Per-environment repositories and more than one store binding per project. The binding is
  a list so this is additive later, but it is not built here.
- Credentials, connection strings and any other secret: owned by `secrets-vault`, never by
  this Feature's store.
- A platform-owned hosted OpenVaultDB endpoint and its platform-wide `/ext/v0/…` path
  convention. `/ext/v0` appears nowhere in the OpenVaultDB repositories and is a platform
  convention for platform-owned endpoints, not part of the OpenVaultDB API; a user-supplied
  OpenVaultDB host base URL validated by `/.well-known/openvaultdb` discovery does not depend
  on it. The platform endpoint is tracked as a platform Feature in the platform repo rather
  than duplicated here.
- A full `openvaultdb-ts` SDK. The browser reader is deliberately minimal and disposable.
- Read parity in the GitHub `GhStorage` backend: `git-native-projects`
  REQ:github-hosted-projects-read-write requires `GhStorage.FileExists`/`OpenFile` to return
  written content, but both still `panic("implement me")` in `datatug-cli`
  (`pkg/dtgithub/github_storage.go`). `git-native-projects`
  AC:github-project-reads-back-what-it-wrote is therefore **undeliverable as scoped here** and
  MUST NOT be recorded as a Supporting AC of this Feature's plan: this Feature's browser read
  path reaches the repository over its REST API and does not use `GhStorage`, and adding
  `GhStorage` read support is a separate task owned by `git-native-projects` — an explicit
  prerequisite for that AC, not an obligation of this Feature.
- Client-direct writes to the project store, which `git-native-projects` already excludes.
- Direct browser-to-database query execution.

## Open Questions

- The binding endpoint shape is settled, not provisional: it is a **user-supplied
  OpenVaultDB host base URL**, validated by the `/.well-known/openvaultdb` discovery document
  OpenVaultDB serves (`openvaultdb-go/pkg/server/server.go`), with the data plane addressed
  under that base at `/v1`. `/ext/v0` is a platform convention for a **platform-owned** hosted
  endpoint, which is out of scope (see Not Doing / Out of Scope); it is not required of, and
  does not constrain, a user-supplied host.

---

*This document follows the https://specscore.md/feature-specification*
