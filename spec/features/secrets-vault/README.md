---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Secrets vault

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/secrets-vault?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/secrets-vault?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/secrets-vault?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/secrets-vault?op=request-change) |
**Status:** Approved
**Date:** 2026-09-16
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Depends-On:** openvaultdb-project-store
**Grade:** B

## Summary

A secrets vault stores the connection strings and credentials DataTug needs in order to
execute queries, as data in OpenVaultDB storage, held separately from project files and
never readable by a user or any other client.

Vaults are bound to **projects**, and that relation is many-to-many: a project may bind
more than one vault, and one vault may serve several projects. The reference arrangement,
and the reason the relation is plural, is blast radius: **one vault for DEV and UAT and a
separate vault for PROD**. Which vault holds which environment's secrets is an explicit
owner choice — DataTug applies no default — and the flow selects among the vaults the owner
has already registered with OpenVaultDB. DataTug does not create vaults: registration and
provisioning belong to OpenVaultDB, and a DataTug vault is a **reference** to an
OpenVaultDB-registered vault (host base URL + vault id). What DataTug owns is the
project→vault mapping and the credential envelope.

Credentials are addressed with the same namespace rule as a project — a collection whose name
carries the static prefix and whose record keys carry the identity segments, declared in the
vault store's bound OpenVaultDB manifest (REQ:vault-mirrors-project-namespace) — so one vault
serves several projects without collision, and a vault's structure is readable by reading
files exactly as a project is — by the store natively or through the server, never by the
client, since the credentials inside it are not client-readable (REQ:client-cannot-read-credentials).

Credentials are available to the **DataTug runtime** for every user who has access to the
project. The runtime reaches a vault because the DataTug **server** holds an OpenVaultDB
grant for it — the app grant the user gave when connecting the vault
(REQ:server-authorizes-vault-access); a user does not need direct access to a vault and is
not given one. A project refers to connections by logical name, and the runtime resolves
those names against the vaults bound to that project. Access to credentials therefore
follows access to the project **on top of** the OpenVaultDB app grant, and there is no
separate user-facing vault permission to grant or revoke.

This Feature extends, and does not weaken, `git-native-projects`
REQ:secrets-never-in-project-files and the shared guard in `heterogeneous-sources`
REQ:secrets-outside-project-files. Those requirements say secrets must never be *in project
files*; this Feature supplies the storage they belong in.

## Problem

The audit found secrets modelling to be a live gap: `QueryDefTarget.Credentials` can be
committed into a git-tracked query file, and `Project.Validate()` never validates
`Queries` at all. At the same time a DataTug project is useless without credentials —
executing a query against a real environment requires them — so "keep secrets out of the
project" only works if there is somewhere for them to go.

Two forces make a single flat credential store the wrong shape:

- **Blast radius.** Production credentials are not equivalent to development ones. Keeping
  them in one store means one compromise, one over-broad read, or one careless grant
  exposes production. The DEV/UAT-versus-PROD separation exists to bound that, which is why the
  vault relation is per project and plural rather than one store per user.
- **Identity of the reader.** Credentials must be usable by whoever may run the project's
  queries, without each of them being handed the credential. That points at the runtime as
  the reader — the hosted or self-hosted backend — rather than at individual users.

The runtimes still differ in trust position, and the model is explicit about it: a hosted
backend holds credentials on the user's behalf; a self-hosted backend does so on
infrastructure the user controls; and in the CLI the user *is* the runtime, so claiming
that credentials are hidden from them there would be false. The CLI keeps resolving from
local, non-committed storage for now, and feeding a local CLI from a central store is out
of scope.

## Behavior

### Vault model

#### REQ: vault-is-a-first-class-entity

A vault is registered and provisioned by **OpenVaultDB**, not by DataTug. A DataTug vault
MUST be a **reference** to an OpenVaultDB-registered vault, identified by its host base URL
and its vault id — the pointer OpenVaultDB's wallet persists under `/users/{uid}/vaults` and
surfaces at `openvaultdb.com/my/vaults` (`openvaultdb-com/public/js/my-vaults.js`,
`wallet-store.js`). DataTug MUST NOT own vault registration or vault provisioning, MUST NOT
introduce a vault registry of its own, and MUST consume the vaults the host already exposes.
Decision `openvaultdb-com/spec/decisions/0003-host-vault-namespace-model` makes the vault the
unit the user owns, names, and grants apps access to, and explicitly defers vault
*provisioning*, so DataTug MUST NOT assume a create-vault capability exists. What DataTug
owns is only the **project → vault mapping** (REQ:project-vault-binding) and the
**credential envelope** (REQ:credential-envelope) stored inside the referenced vault.

The reference MUST be shared, not private: a single registered vault MUST be able to serve
several projects, and a vault MUST NOT be a per-user singleton. Access to a vault's contents
MUST be derived from the OpenVaultDB **app grant** the DataTug server holds
(REQ:server-authorizes-vault-access) plus access to the projects that bind it; DataTug MUST
NOT require, and MUST NOT introduce, a separate user-facing vault grant list. Founder
ruling, 2026-09-16: the relation is **projects ↔ vaults**, not users ↔ vaults; users relate
to vaults only transitively through projects.

#### REQ: project-vault-binding

A project MUST let its owner choose, **explicitly**, which vault holds secrets, and the
configuration flow MAY present that choice in environment terms. DataTug MUST NOT apply a
default grouping and MUST NOT silently pick a vault: the owner selects one from the vaults
OpenVaultDB has already registered for the owner, and DataTug MUST NOT present itself as
creating a vault. Where the owner needs a vault that does not exist yet, the flow MUST
delegate creation to OpenVaultDB's vault-registration surface — a deep link to the host's
registration page, for example `openvaultdb.com/my/vaults` — and MUST state the OpenVaultDB
prerequisite explicitly rather than creating or provisioning a vault itself. The choice MUST
resolve to the server-to-vault mapping below. The reference arrangement MUST be supported:
**one vault for DEV and UAT, and a separate vault for PROD**. Founder ruling, 2026-09-16,
verbatim: *"project can have 1 vault for DEV and UAT and 1 for PROD env and that credentials
will be available to datatug runtime for all users that have access to the project"*; and on
the mechanism: *"explicit grouping - you choose a vault where to store env secrets. and you
can create a new vault or select from existing ones"* — the *create* half of that ruling is
delivered by delegating to OpenVaultDB registration, because registration and provisioning
are OpenVaultDB's (REQ:vault-is-a-first-class-entity).

Credentials belong to a **server or endpoint**, not to an environment, and several
environments MAY share one server and therefore one credential. The persisted assignment is
therefore a **server-to-vault mapping**, and an environment MUST inherit the vault of the
server it points at. DataTug MUST NOT store an independent environment-to-vault mapping, so
that a credential has exactly one home. A choice that would place one server in two
different vaults MUST be rejected, naming the shared server, rather than silently
duplicating the credential.

The assignment MUST be project-level configuration recorded with the project, and MUST NOT
require a second project-store binding — a vault owns its own storage, so binding more
vaults MUST NOT multiply the project's store bindings.

#### REQ: vault-mirrors-project-namespace

A vault MUST address credentials through the same application-namespace rule the project
store uses: the subtree is expressed by the **collection name** in the vault store's bound
OpenVaultDB manifest, and the identity segments are record keys within that collection. The
vault's manifest MUST declare the collection `datatug/credentials`, and DataTug MUST address a
credential within it by the project's id and the endpoint's server id
(`<project-id>/<server-id>`), so that one vault can serve several projects without collisions
and a person or agent can read the vault by reading files. This Feature MUST NOT mandate an
on-disk directory tree: the exact layout a backend derives from the declaration is the
backend's business, and no requirement here may assert a path a named backend does not
produce.

The two named backends realise the same declaration differently, exactly as
`openvaultdb-project-store` REQ:project-namespace-layout describes: the local inGitDB backend
resolves it from a root `/.ingitdb/root-collections.yaml` entry whose `data_dir` is resolved
relative to the parent of `.collections/`, while the OpenVaultDB GitHub backend maps the
declared collection name directly to its directory (`DirPath: name` in
`buildInGitDBDefinition`, `openvaultdb-go/pkg/mount/ingitdb_github.go`) and takes the
collection's schema from the bound manifest's `schemas.collections`. DataTug MUST NOT require
that repository to carry `/.ingitdb/root-collections.yaml` or a `datatug/.collections/`
schema directory; on that backend those clauses are no-ops.

Physical nesting is an upstream prerequisite, not an assumption. `schemas.collections`
entries carry only `fields` — there is no path, subtree or `data_dir` field on a collection
(`openvaultdb-go/pkg/schema/schema.go`) — and `buildInGitDBDefinition` fixes each
collection's directory to its name, so a declaration cannot by itself place the `credentials`
collection under a per-project subtree or give each server id its own directory. DataTug MUST
NOT assume credentials physically occupy
`/datatug/projects/<project-id>/credentials/<server-id>/`, and MUST NOT assert that path
until `openvaultdb-go` supplies **per-collection path support honoured by
`buildInGitDBDefinition`**; that support is an explicit upstream prerequisite of this
Feature, shared with `openvaultdb-project-store` REQ:project-namespace-layout. Until it
lands, the namespace is satisfied by the `datatug/credentials` collection name plus the
project id and server id as record keys.

The credential key MUST NOT carry an environment segment. Because several environments
may share one endpoint and therefore one credential, an environment-keyed key would
store that credential once per environment and let the copies drift out of step — which is
precisely the failure the endpoint key exists to prevent.

Where a vault shares a store with project data, the store's rules MUST deny clients the
credential collection, and the client read path MUST NOT request it. This namespace is the
boundary in both cases: it is what a path-based rule can protect, and what physical
separation keeps out of a client-readable store.

#### REQ: vault-storage-is-pluggable

A vault's storage MUST be selectable like any other OpenVaultDB storage, and DataTug MUST
NOT require a particular backing store. Whether a given store may hold a vault's
credentials depends on that store's own enforcement. A store that can deny the client
access by its own rules MAY hold them alongside other data, even in the same database — a
server-mediated store such as Firestore with security rules is the reference case. A store
the client reads directly with a user-supplied credential — a git repository reached over
its REST API — has no such rules and cannot protect an individual path, so it MUST NOT hold
credentials anywhere that credential can reach; credentials for such a project MUST live in
a separate store. The rule is unconditional in each case: there is no store in which
client-readable credentials are acceptable.

#### REQ: client-cannot-read-credentials

A credential stored in a vault MUST NOT be readable by the browser, by any other client
DataTug ships, or by a user acting directly against the vault's storage. Enforcement is
the storage's obligation, by access rules or by separation, and DataTug MUST detect and
refuse or warn when a chosen binding cannot provide it — in particular when a client read
credential can reach the vault's storage. The credential DataTug uses on the browser read
path MUST itself be read-only by construction (REQ:client-holds-no-write-credential), so a
client cannot write to a vault it can reach.

#### REQ: server-authorizes-vault-access

The DataTug server MUST hold an explicit authorization to read a vault before it resolves any
credential from it. The authorization MUST be obtained through one of: an OpenVaultDB
Connect/app token scoped to the vault (the OVDB Connect consent flow, then that host's
`/authorize`), a GitHub App installation covering the vault repository, or a user-supplied
personal access token. DataTug MUST NOT read a vault on the strength of a user's project
access alone.

Whatever form the authorization takes, the credential or token that carries it MUST be held
in a **server-only** store and MUST NEVER be written to the client-readable registry at
`users/{uid}/ext/datatug` or into any client bundle. A missing or rejected authorization MUST
fail closed, naming the missing authorization, rather than falling back to an anonymous or
default read.

#### REQ: client-holds-no-write-credential

DataTug MUST issue no write credential to a client for a vault or for a project store that
can reach a vault. The credential DataTug supplies to, or accepts from, a client for the
browser read path MUST be read-only by construction: DataTug MUST NOT hand a client a
`contents:write` git token or an OpenVaultDB **owner** token, and MUST NOT accept one as the
browser read credential. A write-capable credential on the client defeats
REQ:client-cannot-read-credentials and `openvaultdb-project-store`
REQ:server-mediated-writes, so presenting one is a binding defect, not a user preference.

#### REQ: credential-envelope

Each stored credential MUST be held in an opaque envelope carrying an algorithm or version
marker, so that encryption can be introduced later as a per-vault setting without a data
migration. The initial implementation writes the plaintext variant of that envelope.

### Resolution

#### REQ: runtime-resolves-credentials

A runtime MUST resolve a project's logical connection names by following
environment → server → vault: the environment's server determines the vault, and the
credential is looked up there. The server MUST resolve from the project's bound vault using
the OpenVaultDB authorization it holds (REQ:server-authorizes-vault-access), and MUST ALSO
check that the requesting user has access to the project: project access is an additional
check on top of the OpenVaultDB app grant, never a substitute for it. It MUST fail closed
when a name cannot be resolved rather than substituting a default or an empty credential,
and MUST fail closed when the OpenVaultDB authorization is absent or rejected. Resolution
MUST be available for every user who has access to the project, and MUST NOT require that
user to hold access to the vault. The local CLI MUST resolve from local, non-committed
storage for now.

#### REQ: vault-binding-is-visible

DataTug MUST show which vaults a project binds and which environments each vault serves,
and MUST make the two-layer authorization explicit: the OpenVaultDB **app grant** for a
vault is held by the DataTug **server** (the app), and a user's own ability to use
credentials follows access to the **project**, which is an additional DataTug-side check on
top of that grant rather than a substitute for it. DataTug MUST also make clear that a user
is given no access to the vault itself. Under the initial plaintext envelope, anyone able to
read a vault's storage can read its credentials, so the vault's storage choice is the
security decision; DataTug MUST state this where the binding is configured rather than
implying a protection the storage does not provide.

## Acceptance Criteria

### AC: explicit-vault-assignment (verifies REQ:project-vault-binding)

**Given** a project whose DEV and UAT environments point at one server, whose PROD environment points at a different server, and with no vault assignment yet
**When** the owner configures where secrets are stored
**Then** DataTug proposes no default vault, the flow offers the owner's OpenVaultDB-registered vaults for selection and, when a suitable vault does not exist, links to OpenVaultDB's vault-registration surface instead of creating one itself, the choice is persisted as a server-to-vault mapping that places the DEV and UAT server in one vault and the PROD server in another, no independent environment-to-vault mapping is stored, the environments resolve to their server's vault, and the project still holds exactly one project-store binding.

### AC: new-vault-is-registered-with-openvaultdb (verifies REQ:vault-is-a-first-class-entity, REQ:project-vault-binding)

**Given** a project owner with no suitable OpenVaultDB-registered vault
**When** the owner needs a new vault while configuring where secrets are stored
**Then** DataTug does not create or provision a vault itself, it directs the owner to OpenVaultDB's vault-registration surface and states the OpenVaultDB prerequisite, and the binding is recorded only once the vault appears in the owner's OpenVaultDB-registered vault list.

### AC: a-vault-serves-several-projects (verifies REQ:vault-is-a-first-class-entity)

**Given** one OpenVaultDB-registered vault referenced by two different projects
**When** either project executes a query against an environment the vault serves
**Then** both projects resolve credentials from that vault, and the vault is owned by its OpenVaultDB host and is not owned by or private to either project.

### AC: vault-is-laid-out-per-project (verifies REQ:vault-mirrors-project-namespace)

**Given** one vault bound by two projects
**When** credentials are written for both
**Then** each project's credentials are records of the `datatug/credentials` collection declared in the vault store's bound manifest, keyed by that project's id and the server id, neither project's records expose the other's credentials, and the collection is declared in the manifest's `schemas.collections` on the OpenVaultDB GitHub backend and mapped by a `/.ingitdb/root-collections.yaml` `data_dir` entry on the local inGitDB backend. On the OpenVaultDB GitHub backend, no `/.ingitdb/root-collections.yaml` and no `datatug/.collections/` schema directory is required in the repository.

### AC: credentials-are-keyed-by-endpoint (verifies REQ:vault-mirrors-project-namespace)

**Given** a project whose DEV and UAT environments both point at the same server
**When** a credential is written for that server
**Then** exactly one credential record exists for the project-and-server key in the `datatug/credentials` collection, both environments resolve it, and no environment-keyed duplicate exists.

### AC: a-shared-server-has-one-vault-home (verifies REQ:project-vault-binding)

**Given** one server referenced by both a DEV environment and a PROD environment
**When** the owner assigns those environments to two different vaults
**Then** DataTug rejects the assignment, names the shared server as the conflict, and no credential is written to either vault until the conflict is resolved.

### AC: access-follows-project-not-vault (verifies REQ:vault-is-a-first-class-entity, REQ:vault-binding-is-visible)

**Given** a user with access to a project that binds a vault, and no DataTug vault-level grant
**When** the user's access to credentials is evaluated
**Then** the server first presents the OpenVaultDB app grant it holds for the vault, then checks that the user has access to the project, the user can run queries that use the vault's credentials through the runtime only when both checks pass, and no separate DataTug vault grant exists to list, grant or revoke.

### AC: server-authorizes-against-openvaultdb (verifies REQ:server-authorizes-vault-access)

**Given** a project bound to an OpenVaultDB-registered vault and a DataTug server that holds no OpenVaultDB authorization for it
**When** the server attempts to resolve a credential from that vault
**Then** resolution fails closed naming the missing authorization, and succeeds only after the server has obtained authorization through an OVDB Connect/app token, a GitHub App installation or a user-supplied token.

### AC: vault-authorization-stays-server-side (verifies REQ:server-authorizes-vault-access)

**Given** a server that has obtained authorization to read a vault
**When** the client-readable registry at `users/{uid}/ext/datatug` and every client bundle are inspected
**Then** neither carries the vault-read credential or token, which lives only in the server-side store.

### AC: client-credential-is-read-only (verifies REQ:client-holds-no-write-credential)

**Given** a client holding the credential DataTug supplies to, or accepts from, it for the browser read path
**When** the client attempts a write to anything that credential reaches, including a colocated vault's `credentials` collection
**Then** the write is refused because the credential is read-only by construction, and DataTug has issued no `contents:write` git token and no OpenVaultDB owner token to the client.

### AC: firestore-backed-vault-is-client-inaccessible (verifies REQ:vault-storage-is-pluggable, REQ:client-cannot-read-credentials)

**Given** a vault stored in a server-mediated store with security rules, sharing a database with non-secret project data
**When** a client or a user acting directly against that database attempts to read a credential
**Then** the store denies the read, while the server, presenting the OpenVaultDB authorization it holds for the vault and no client credential, reads the same credential successfully.

### AC: binding-is-refused-when-credentials-would-be-client-readable (verifies REQ:client-cannot-read-credentials)

**Given** a proposed vault binding whose storage a client read credential can reach
**When** the binding is saved
**Then** DataTug either refuses it outright or records it with an explicit statement that the credentials would be client-readable, records which of the two outcomes occurred, and in neither case represents the credentials as client-inaccessible; a client credential that is not read-only by construction is refused, never recorded.

### AC: git-backed-vault-must-be-out-of-client-reach (verifies REQ:vault-storage-is-pluggable)

**Given** a project store in a git repository that the client reads directly with a user-supplied credential
**When** a vault binding is configured such that the same credential can reach the vault's storage
**Then** DataTug refuses or warns on the binding and does not claim the credentials are client-inaccessible, and the user-supplied credential is required to be read-only so the client cannot write to anything that credential reaches.

### AC: a-project-carries-no-credential (verifies REQ:runtime-resolves-credentials)

**Given** a project whose environments point at logical connection names backed by two vaults
**When** the project's files are inspected
**Then** no file contains a credential value, the project records one store binding and a server-to-vault mapping with each environment inheriting the vault of the server it points at, and connection names resolve only at runtime.

### AC: unresolvable-connection-fails-closed (verifies REQ:runtime-resolves-credentials)

**Given** a project referencing a logical connection name that the vault reached through the target environment's server does not define
**When** a query is executed against that environment
**Then** execution fails naming the unresolved connection and no query reaches any environment.

### AC: local-cli-uses-local-secrets (verifies REQ:runtime-resolves-credentials)

**Given** the CLI running in a user's terminal with a project bound to a remote OpenVaultDB store and to a vault
**When** a query executes
**Then** the credential is resolved from local, non-committed storage and nothing is written into the project store.

### AC: bindings-are-visible (verifies REQ:vault-binding-is-visible)

**Given** a project bound to a DEV/UAT vault and a PROD vault
**When** the owner opens the project's credential configuration
**Then** both vaults are shown with the servers each holds and the environments those servers serve, the storage choice is presented as the security decision with the plaintext-envelope limitation stated where the binding is configured, and the user is told that the DataTug server holds the OpenVaultDB app grant for each vault and that the user's own access to credentials follows access to the project on top of that grant.

### AC: envelope-is-version-tagged (verifies REQ:credential-envelope)

**Given** a credential written by the initial implementation
**When** the stored value is inspected
**Then** it carries an algorithm or version marker identifying the plaintext variant, so a later encrypted variant is distinguishable without re-reading the secret itself.

## Rehearse Integration

Every AC has a concrete behavioural surface — vault bindings, store rule outcomes,
resolution failures, CLI resolution, and the stored envelope's shape — and all are
testable in principle. The capability is realised across `datatug/backend` (vault bindings
and runtime resolution), `datatug-cli` (local resolution), `datatug-apps` (binding
visibility) and the chosen store's own rules, so stub scaffolding is deferred to those
implementations rather than authored here, consistent with `ai-query-builder` and
`serve-brokered-query-builder`. The security claims deserve stubs early because they are the
claims a reviewer must be able to falsify rather than conveniences:
AC:firestore-backed-vault-is-client-inaccessible,
AC:git-backed-vault-must-be-out-of-client-reach, AC:access-follows-project-not-vault,
AC:server-authorizes-against-openvaultdb, AC:vault-authorization-stays-server-side and
AC:client-credential-is-read-only.

## Not Doing / Out of Scope

- Encryption of vault contents, and the key-custody design it needs. The envelope in
  REQ:credential-envelope makes it additive; it is not built here.
- Distributing centrally stored secrets to a locally running CLI. The CLI keeps local
  secrets for now.
- Secrets in project files, in any store, under any circumstances.
- Vault registration and provisioning: owned by OpenVaultDB, not by DataTug. DataTug
  references vaults the host already exposes and never creates one; the registration
  surface is OpenVaultDB's, reached by delegation.
- User-facing vault grants: access derives from project access on top of the OpenVaultDB app
  grant the server holds, so there is nothing to grant per user.
- Per-environment vault granularity finer than the server-to-vault mapping (one vault per
  DEV, UAT and PROD separately), and per-project vaults.
- A new permission model for vaults; access reuses `sharing-scopes` and
  `server-enforced-access`.

## Open Questions

- The key-custody story for the deferred encrypted envelope, which differs per runtime and
  cannot be settled without the hosted/self-hosted/CLI split being final.

---

*This document follows the https://specscore.md/feature-specification*
