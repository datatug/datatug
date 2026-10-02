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
**Depends-On:** openvaultdb-project-store, dalgo-project-store
**Grade:** B

## Summary

A secrets vault stores the connection strings and credentials DataTug needs in order to
execute queries, as data in OpenVaultDB storage, held separately from project files and
never readable by a user or any other client.

Store bindings live on the **project record** as **role-tagged** bindings. Every project MUST
have exactly one binding with role `project` — its own items store, in which the project
record itself lives, which is what makes the bootstrap well-founded. Each credentials store is
a binding with role `secrets`, and there MAY be ZERO of them — either the secrets live in the
project's own items store, or the project has none yet. A project MAY bind more than one
secrets store, and one store MAY hold several DataTug projects. The reference arrangement, and
the reason the secrets role is plural, is blast radius: **one vault for DEV
and UAT and a separate vault for PROD**. Which vault holds which environment's secrets is an
explicit owner choice — DataTug applies no default — and the flow selects among the vaults
the **host already exposes**. DataTug does not create or provision vaults: the wallet at
`openvaultdb.com/my/vaults` only adds pointers to vaults a host already exposes. A DataTug
vault is a **reference** to an OpenVaultDB-registered vault (host base URL + vault id). What
DataTug owns is the role-`secrets` semantics, the server-to-vault assignment a `secrets`
binding carries, and the credential envelope; the binding record's field set is owned by
`openvaultdb-project-store` REQ:project-store-binding, which this Feature cites rather than
restating.

The typical arrangement is two stores with separate host ACLs: an **items store** for project
assets and a **credentials store** for secrets. The browser reads the items store directly
with the reader's own credential, and is NEVER given a credential for the credentials store.
The server reads credentials through dalgo/OpenVaultDB uniformly — it asks for a credential
for a given environment from a given store, and dalgo routes the request — and is unaware of
the backing storage (GitHub, Firestore, or a managed vault — an upstream/blocked backend,
REQ:vault-storage-is-pluggable).

Credentials are addressed by the owning project and the endpoint's server at the location
`dalgo-project-store` declares: collection `credentials` as a subcollection of the `projects`
record, key path `ext/datatug/projects/<project-id>/credentials/<server-id>`, file
`ext/datatug/projects/<project-id>/credentials/<server-id>/<server-id>.credential.json`, suffix
`credential` (REQ:vault-mirrors-project-namespace) — so a vault that serves several projects
keeps them apart, and a vault's structure is readable by reading files exactly as a project is
— by the store natively or through the server, never by the client, since the credentials
inside it are not client-readable (REQ:client-cannot-read-credentials).

Credentials are available to the **DataTug runtime** for every user who has access to the
project. The runtime reaches a vault because the DataTug **server** holds an OpenVaultDB
grant for it — the app grant the user gave when connecting the vault
(REQ:server-authorizes-vault-access); a user does not need direct access to a vault and is
not given one. A project refers to connections by logical name, and the runtime resolves
those names against the secrets stores bound to that project. Access to credentials therefore
follows access to the project **on top of** the OpenVaultDB app grant, and there is no
separate user-facing vault permission to grant or revoke.

The per-user registry at `users/{uid}/ext/datatug` is a **discovery index** only; a project's
bindings do not live there.

This Feature extends, and does not weaken, `git-native-projects`
REQ:secrets-never-in-project-files and the shared guard in `heterogeneous-sources`
REQ:secrets-outside-project-files. Those requirements say secrets must never be *in project
files*; this Feature supplies the storage they belong in.

## Problem

The audit found secrets modelling to be a live gap: `QueryDefTarget.Credentials` can be
committed into a git-tracked query file. The audit's accompanying claim that
`Project.Validate()` never validates `Queries` at all is **false**: `Project.Validate()` calls
`p.Queries.Folders.Validate()` and `p.Queries.Items.Validate()`
(`datatug-core/pkg/datatug/project.go:123-133`), which walks every `QueryDef`
(`pkg/datatug/query.go:74-81`), and `QueryDefTarget.Validate` (`query.go:176-189`) rejects
credential material. The residual gap is real but narrower: the resulting error does not name
the offending FILE, so an author has to hunt for it. At the same time a DataTug project is
useless without credentials — executing a query against a real environment requires them — so
"keep secrets out of the project" only works if there is somewhere for them to go.

Two forces make a single flat credential store the wrong shape:

- **Blast radius.** Production credentials are not equivalent to development ones. Keeping
  them in one store means one compromise, one over-broad read, or one careless grant
  exposes production. The DEV/UAT-versus-PROD separation exists to bound that, which is why a
  project's secrets bindings are plural rather than one store per user.
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

A vault is provisioned by **OpenVaultDB**, not by DataTug. A DataTug vault MUST be a
**reference** to an OpenVaultDB-registered vault, identified by its host base URL and its
vault id — the pointer OpenVaultDB's wallet persists under `/users/{uid}/vaults` and surfaces
at `openvaultdb.com/my/vaults` (`openvaultdb-com/public/js/my-vaults.js`, `wallet-store.js`).
An OpenVaultDB vault **is** an OpenVaultDB database, so the vault id and the database id are
the same identifier; the pair that identifies a bound store is therefore its host base URL and
its vault id.
That page is a **wallet** surface, not a vault-creation surface: it selects among the
vaults a host already exposes (it lists them with `GET {baseUrl}/vaults` in `fetchVaults`) or
registers an existing repository as a vault, and `wallet-store.js` writes only pointers. It
does **not** create or provision a vault, and DataTug MUST NOT claim that it does. DataTug
MUST NOT own vault registration or vault provisioning, MUST NOT introduce a vault registry of
its own, and MUST consume the vaults the host already exposes.

Decision `openvaultdb-com/spec/decisions/0003-host-vault-namespace-model` (**status: In
Review**) makes the vault the unit the user owns, names, and grants apps access to, and
explicitly defers vault *provisioning*: the wallet "uses the vaults a host already exposes (a
default vault, or ones set up out of band)", and "creating named vaults from the wallet needs
a create-vault capability on the host — a later step". DataTug therefore MUST NOT assume a
create-vault capability exists, and MUST state out-of-band provisioning as an explicit
prerequisite the owner satisfies on the host before DataTug can bind the vault it produces.
What DataTug owns is only the **role-`secrets` semantics and the server-to-vault assignment**
(`secrets-vault` REQ:project-vault-binding) and the **credential envelope**
(REQ:credential-envelope) stored inside the referenced vault; the binding record's field set is
owned by `openvaultdb-project-store` REQ:project-store-binding, and this Feature cites it
rather than restating it.

The reference MUST be shared, not private: a single registered vault MUST be able to serve
several projects, and a vault MUST NOT be a per-user singleton. Access to a vault's contents
MUST be derived from the OpenVaultDB **app grant** the DataTug server holds
(REQ:server-authorizes-vault-access) plus access to the projects that bind it; DataTug MUST
NOT require, and MUST NOT introduce, a separate user-facing vault grant list. Founder
ruling, 2026-09-16: the relation is **projects ↔ vaults**, not users ↔ vaults; users relate
to vaults only transitively through projects.

#### REQ: project-vault-binding

A project MUST record its store bindings **on the project record**, as **role-tagged**
bindings. A project MUST have exactly one binding whose role is `project` — its own items
store, in which the project record itself lives — and the `project`-role store MUST be chosen
before the project's first write. A secrets store is a binding whose role is `secrets`; there
MAY be ZERO of them, and a project MAY carry more than one — the reference arrangement is one
for DEV and UAT and one for PROD. Binding zero `secrets` stores is legal: the project's
secrets then either live in its own items store or do not exist yet. An unknown role
value MUST be rejected. The bindings are additive: adding a role-`secrets` binding MUST NOT
create a second role-`project` binding, MUST NOT duplicate the items-store binding, and MUST
NOT change which store holds the project's assets. A single store MAY in turn hold several
DataTug projects.

Bindings MUST be **single-sourced on the project record**, so that a SHARED project's
bindings cannot disagree between users. The client-readable registry at
`users/{uid}/ext/datatug` MUST remain a per-user **discovery index** only — its shipped shape
is `UserExt.Stores map[string]*StoreBrief` keyed by store id
(`backend/models4datatug/project.go:79-92`) — and bindings MUST NOT live there.

A project MUST let its owner choose, **explicitly**, which vault holds secrets, and the
configuration flow MAY present that choice in environment terms. DataTug MUST NOT apply a
default grouping and MUST NOT silently pick a vault: the owner selects one from the vaults
the host already exposes, and DataTug MUST NOT present itself as creating a vault. Where the
owner needs a vault that does not exist yet, the flow MUST state the out-of-band
OpenVaultDB provisioning prerequisite explicitly (REQ:vault-is-a-first-class-entity) and MUST
direct the owner to the host's wallet surface to add a vault the host already exposes; it
MUST NOT claim that this page creates or provisions a vault, and MUST NOT create or
provision one itself. The choice MUST resolve to the server-to-vault mapping below. The
reference arrangement MUST be supported: **one vault for DEV and UAT, and a separate vault
for PROD**. Founder ruling, 2026-09-16, verbatim: *"project can have 1 vault for DEV and UAT
and 1 for PROD env and that credentials will be available to datatug runtime for all users
that have access to the project"*; and on the mechanism: *"explicit grouping - you choose a
vault where to store env secrets. and you can create a new vault or select from existing
ones"* — the *create* half of that ruling is delivered by the owner provisioning the vault
out of band and then adding it in the host's wallet, because provisioning is deferred and the
wallet only adds pointers (REQ:vault-is-a-first-class-entity).

Credentials belong to a **server or endpoint**, not to an environment, and several
environments MAY share one server and therefore one credential. The persisted assignment is
therefore a **server-to-vault mapping**, and an environment MUST inherit the vault of the
server it points at. Its carrier is explicit: the **servers field of the role-`secrets`
binding** whose field set is defined once by `openvaultdb-project-store`
REQ:project-store-binding — with one vault serving DEV and UAT, that binding's servers field
holds both, and nothing else records the assignment. DataTug MUST NOT store an independent
environment-to-vault mapping, so that a credential has exactly one home. A choice that would
place one server in two different vaults MUST be rejected, naming the shared server, rather
than silently duplicating the credential.

#### REQ: vault-mirrors-project-namespace

The credential's location and filename are declared, with the rest of the canonical
project-item layout and its naming policy, by `dalgo-project-store`: collection `credentials`
as a subcollection of the `projects` record, key path
`ext/datatug/projects/<project-id>/credentials/<server-id>`, file
`ext/datatug/projects/<project-id>/credentials/<server-id>/<server-id>.credential.json`, suffix
`credential`. This Feature names only that address and MUST NOT define the naming policy —
owned by `dalgo-project-store` REQ:canonical-project-layout,
REQ:type-suffix-is-self-identifying and REQ:hierarchy-is-a-key-path — and MUST NOT restate the
directory shape, the `records_dir` setting, the record format or the type-suffix grammar.

In that address, `projects` is the parent collection, `<project-id>` is the parent record key,
`credentials` is its subcollection, and `<server-id>` is the credential's record key. A
credential MUST be its own record, so that a store holding several projects keeps their
credentials apart and a person or agent can read the vault by reading files. The owning
project's id MUST be part of the address: a credentials store MAY hold several DataTug
projects, so a credential addressed by its endpoint alone cannot be kept apart from another
project's credential for that endpoint.

The vault store's `ext/datatug/projects/<project-id>` is a **key-path scoping parent**, not a copy of the
project: the project's own record lives in the items store, the vault store is a different
store (possibly a different repository), and the parent record need not exist in the vault
store. The scoping comes from the key path, not from a materialised parent.

The credential key MUST NOT carry an environment segment. Because several environments
may share one endpoint and therefore one credential, an environment-keyed key would
store that credential once per environment and let the copies drift out of step — which is
precisely the failure the endpoint key exists to prevent.

The remote (OpenVaultDB) path's ability to express this nesting is an upstream prerequisite
recorded by `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded: a declared
collection must be able to express its subcollection names and the full per-collection `record_file` configuration (`format`, the `name` template, `records_dir` and `type`, not `format` and `name` alone), which `openvaultdb-go`'s
`buildInGitDBDefinition` cannot do today. This Feature MUST NOT claim the remote nested
credentials layout works, and no AC clause that depends on it MAY be scaffolded as passing.
Only the OpenVaultDB leg waits on it: the local `dalgo2ingitdb` driver already resolves the
nesting (`dalgo2ingitdb/scoped_collection.go`), so the same key path is exercisable today
wherever credentials are placed in a local store — the zero-`secrets`-binding placement in a
local project items store that holds several projects — and that leg MUST NOT be blocked
along with the OpenVaultDB one.

Where a credentials store shares its storage with client-readable project data, binding
validation MUST reject the placement by default under REQ:client-cannot-read-credentials; the
canonical key path is what a path-based rule can protect when an owner overrides that default
with informed confirmation, and what physical separation keeps out of a client-readable store.

#### REQ: vault-storage-is-pluggable

A vault's storage MUST be selectable like any other OpenVaultDB storage, and DataTug MUST
NOT require a particular backing store: GitHub and Firestore are candidates, realised as
OpenVaultDB storage engines behind the host (`openvaultdb-go/pkg/mount/ingitdb_github.go`,
`pkg/mount/firestore.go`) and reached by the server through the one `dalgo2openvaultdb`
driver, so supporting them adds no DataTug-side driver and is not the DataTug Cloud or
GitHub-native **project-item store** follow-on work named by `dalgo-project-store`; and a **managed
vault** is an upstream/blocked capability — no DataTug artifact defines what provides it or
what interface it exposes, so this Feature MUST NOT claim it exists, and no AC clause may be
recorded as passing against it until an upstream definition lands. The server MUST read
credentials through dalgo/OpenVaultDB uniformly and MUST be unaware of the backing store
(REQ:server-reads-credentials-uniformly).

The typical arrangement is two stores with separate host ACLs — an **items store** holding
project assets and a **credentials store** holding secrets. The browser reads the items store
directly with its own read-only credential and MUST NEVER be given a credential for the
credentials store. Whether a given store may hold credentials depends on that store's own
enforcement, and the rule admits exactly one path: client-readable credentials are never
acceptable EXCEPT through a recorded informed override under
REQ:client-cannot-read-credentials, and any binding so recorded MUST state which named secrets
are client-readable. A store that cannot deny the client access by its own rules MUST NOT hold
credentials anywhere that client credential can reach except under that recorded override;
where an owner proposes such a placement, binding validation MUST reject it by default
(REQ:client-cannot-read-credentials).

#### REQ: client-cannot-read-credentials

A credential stored in a vault MUST NOT be readable by the browser, by any other client
DataTug ships, or by a user acting directly against the vault's storage. The browser MUST
NEVER be given a credential for the credentials store. Enforcement is the storage's
obligation, by access rules or by separation, and **binding validation is DataTug's single
enforcement owner**: this Feature's binding validation, and no second component.
`openvaultdb-project-store` is a conformance consumer of this rule only and MUST NOT define a
second copy of it.

**The trigger is a credentials placement, not a client-readable store.** The rule's trigger is
a **credentials placement** in a store a client read credential can reach — never the presence
of ordinary project item records in such a store. The mandatory `project`-role binding MUST
NOT be rejected merely because its items store is client-readable: that store is
client-readable by construction, because the browser read path reaches it, and only a
credentials placement there triggers the rule.

Binding validation MUST implement the rule in three parts:

- (a) **Reject by default.** A binding whose storage a client read credential can reach and
  that places credentials there MUST be rejected by default.
- (b) **Enumerate the exposure.** The refusal MUST enumerate the SPECIFIC secrets that would
  be exposed, by name, rather than merely asserting that credentials would be
  client-readable.
- (c) **Informed override.** The user MAY override the rejection only through an explicit
  confirmation step — an "Are you sure?"-style alert showing that list of named secrets. If
  overridden, the binding MUST be recorded with an explicit statement that those named
  secrets are client-readable, and the credentials MUST NOT be represented as
  client-inaccessible anywhere in DataTug.

**The override carrier is named, so the rule has no binding-absent gap.** For a `secrets`
binding, the carrier of (a)–(c) and of the recorded override is that `secrets` binding, as
`openvaultdb-project-store` REQ:project-store-binding records. For a project with **ZERO**
`secrets` bindings whose credentials live in its own items store there is no `secrets` binding
to reject and nowhere on one to record the override, so the carrier is the mandatory
`project`-role binding: that binding's validation performs the refuse-by-default, the
enumeration of the specific secrets, and the informed confirmation, and the recorded override
is attached to that binding. The trigger remains the credentials placement — the `project`-role
binding is not rejected merely because its items store is client-readable — so the policy runs
wherever credentials are placed and is never skipped for want of a `secrets` binding.

The credential DataTug uses on the browser read path MUST itself be read-only by construction
(REQ:client-holds-no-write-credential), so a client cannot write to a vault it can reach.

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

The two browser read paths satisfy this in different ways:
- **An OpenVaultDB project store:** the credential is an OpenVaultDB Connect grant for which DataTug
  requested only `records:read` (`openvaultdb-project-store` REQ:browser-read-path).
- **A git-native project store** (`git-native-projects`), which the client reads straight from
  GitHub: a DataTug-issued GitHub App installation token with `contents:read`, or a user-supplied
  token accepted only when its read-only scope is provable.

#### REQ: credential-envelope

Each stored credential MUST be held in an opaque envelope carrying an algorithm or version
marker, so that encryption can be introduced later as a per-vault setting without a data
migration. The initial implementation writes the plaintext variant of that envelope.

### Resolution

#### REQ: runtime-resolves-credentials

A runtime MUST resolve a project's logical connection names by following
environment → server → vault: the environment's server determines the vault, and the
credential is looked up there. The server MUST resolve from the project's bound secrets store
using the OpenVaultDB authorization it holds (REQ:server-authorizes-vault-access), and MUST
ALSO check that the requesting user has access to the project: project access is an additional
check on top of the OpenVaultDB app grant, never a substitute for it. It MUST fail closed
when a name cannot be resolved rather than substituting a default or an empty credential,
and MUST fail closed when the OpenVaultDB authorization is absent or rejected. Resolution
MUST be available for every user who has access to the project, and MUST NOT require that
user to hold access to the vault. The local CLI MUST resolve from local, non-committed
storage for now.

#### REQ: server-reads-credentials-uniformly

The DataTug server MUST read credentials through dalgo/OpenVaultDB uniformly and MUST be
unaware of the backing storage, whether it is GitHub, Firestore, or a managed vault (an
upstream/blocked backend, REQ:vault-storage-is-pluggable). It MUST
ask for a credential for a given environment from a given store, and dalgo MUST route that
request to the store's backing implementation. All storage-specific handling MUST live below
that interface, so adding a backing store MUST NOT add a storage branch to the server's
resolution path. The server MUST present the authorization the store requires
(REQ:server-authorizes-vault-access) and MUST fail closed when it is missing.

#### REQ: vault-binding-is-visible

DataTug MUST show a project's role-tagged store bindings and which environments each
role-`secrets` store serves, and MUST make the two-layer authorization explicit: the
OpenVaultDB **app grant** for a vault is held by the DataTug **server** (the app), and a
user's own ability to use credentials follows access to the **project**, which is an
additional DataTug-side check on top of that grant rather than a substitute for it. DataTug
MUST also make clear that a user is given no access to the vault itself, and MUST show any
binding recorded through an informed override as client-readable, naming the secrets it
exposes. Under the initial plaintext envelope, anyone able to read a vault's storage can
read its credentials, so the vault's storage choice is the security decision; DataTug MUST
state this where the binding is configured rather than implying a protection the storage does
not provide.

## Acceptance Criteria

### AC: explicit-vault-assignment (verifies REQ:project-vault-binding)

**Given** a project whose DEV and UAT environments point at one server, whose PROD environment points at a different server, and with no secrets binding yet
**When** the owner configures where secrets are stored
**Then** DataTug proposes no default vault, the flow offers the vaults the host already exposes for selection and, when a suitable vault does not exist, states the out-of-band OpenVaultDB provisioning prerequisite and links to the host's wallet surface to add a vault the host already exposes instead of creating one itself, the choice is persisted on the project record as role-`secrets` bindings whose servers field carries the server-to-vault assignment that places the DEV and UAT server in one vault and the PROD server in another, no independent environment-to-vault mapping is stored, the environments resolve to their server's vault, and the project's role-`project` binding is unchanged.

### AC: new-vault-is-registered-with-openvaultdb (verifies REQ:vault-is-a-first-class-entity, REQ:project-vault-binding)

**Given** a project owner with no suitable vault that the host already exposes
**When** the owner needs a new vault while configuring where secrets are stored
**Then** DataTug does not create or provision a vault itself and does not claim the wallet page creates one, it directs the owner to the host's wallet surface to add a vault the host already exposes, it states the out-of-band OpenVaultDB provisioning prerequisite explicitly, and the binding is recorded only once the owner has provisioned the vault out of band and it appears among the vaults the host exposes.

### AC: a-vault-serves-several-projects (verifies REQ:vault-is-a-first-class-entity)

**Given** one OpenVaultDB-registered vault referenced by two different projects
**When** either project executes a query against an environment the vault serves
**Then** both projects resolve credentials from that vault, and the vault is owned by its OpenVaultDB host and is not owned by or private to either project.

### AC: vault-is-laid-out-per-project (verifies REQ:vault-mirrors-project-namespace)

**Given** one vault bound by two projects
**When** credentials are written for both
**Then** each credential is its own record addressed as `ext/datatug/projects/<project-id>/credentials/<server-id>` and written as `ext/datatug/projects/<project-id>/credentials/<server-id>/<server-id>.credential.json` under the canonical project-item layout owned by `dalgo-project-store` REQ:canonical-project-layout, the vault store's `ext/datatug/projects/<project-id>` acts as the scoping parent record without the project's own record existing in that store, neither project's records expose the other's credentials, and no record key contains a path separator. The clause that both projects' credentials are distinguishable as nested records in one shared **OpenVaultDB** vault cannot be exercised until the upstream prerequisite recorded in `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded lands, and this Feature does not claim that layout works there today; the same clause is exercisable now against a local `dalgo2ingitdb` store holding both projects and MUST NOT be blocked with the OpenVaultDB leg.

### AC: credentials-are-keyed-by-endpoint (verifies REQ:vault-mirrors-project-namespace)

**Given** a project whose DEV and UAT environments both point at the same server
**When** a credential is written for that server
**Then** exactly one credential record exists for the owning project and that server at the key path `ext/datatug/projects/<project-id>/credentials/<server-id>`, both environments resolve it, no environment-keyed duplicate exists, and no record key contains a path separator. The clause that the record sits under the owning project's key-path parent in an **OpenVaultDB** vault shared with other projects cannot be exercised until the upstream prerequisite recorded in `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded lands; the same clause is exercisable now against a local `dalgo2ingitdb` store shared with other projects and MUST NOT be blocked with the OpenVaultDB leg.

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

**Given** a vault stored in an OpenVaultDB database on the Firestore storage engine, in a database separate from the one the client's project-store read grant is bound to, while the host's Firestore project also holds non-secret project data
**When** a client attempts to read a credential with its read grant, or with any grant not issued to the DataTug server for that vault
**Then** OpenVaultDB denies the read, because a grant bound to another database never matches the vault's database (`openvaultdb-go/pkg/auth/auth.go` `Allows`), while the server, presenting the OpenVaultDB app grant it holds for the vault and no client credential, reads the same credential successfully. Direct access to the host's Firestore project is the host operator's concern, not a DataTug client's, and is not what this AC asserts.

### AC: binding-is-refused-when-credentials-would-be-client-readable (verifies REQ:client-cannot-read-credentials)

**Given** a proposed vault binding whose storage a client read credential can reach
**When** the binding is saved
**Then** this Feature's binding validation refuses the credentials placement by default, the refusal enumerates the specific secrets that would be exposed by name rather than merely asserting that credentials would be client-readable, the user may override only through an explicit confirmation step that shows that list, and if the user overrides, the carrier binding — the `secrets` binding, or the `project`-role binding when the project has zero `secrets` bindings — is recorded with an explicit statement that those named secrets are client-readable and is not represented as client-inaccessible anywhere; a client credential that is not read-only by construction is refused outright and never recorded, and neither a client-readable `project`-role items store as such nor an ordinary project item record in it triggers the refusal.

### AC: zero-secrets-binding-refusal-rides-the-project-role-binding (verifies REQ:client-cannot-read-credentials, REQ:project-vault-binding)

**Given** a project with zero `role: secrets` bindings whose credentials live in its own items store, and a client read credential that can reach that store
**When** the project's bindings are validated and saved
**Then** there is no `secrets` binding to validate, so validation performs the refuse-by-default, the enumeration of the specific secrets by name, and the informed confirmation on the mandatory `project`-role binding, that binding is saved only after the user confirms through the explicit "Are you sure?"-style alert showing that list, the recorded override is attached to that `project`-role binding as the explicit statement that those named secrets are client-readable, and the binding is not rejected merely because the items store is client-readable nor because it holds ordinary project item records.

### AC: git-backed-vault-must-be-out-of-client-reach (verifies REQ:vault-storage-is-pluggable)

**Given** a project store in a git repository that the client reads directly with a user-supplied credential
**When** a vault binding is configured such that the same credential can reach the vault's storage
**Then** DataTug rejects the binding by default, naming the specific secrets that would be exposed, and does not claim the credentials are client-inaccessible unless the owner records an informed override, and the user-supplied credential is required to be read-only so the client cannot write to anything that credential reaches.

### AC: a-project-carries-no-credential (verifies REQ:runtime-resolves-credentials)

**Given** a project whose environments point at logical connection names backed by two vaults
**When** the project's files are inspected
**Then** no file contains a credential value, the project records its role-tagged bindings on the project record — one role `project` for its items store and one role `secrets` per vault, with each environment inheriting the vault of the server it points at — the per-user registry at `users/{uid}/ext/datatug` carries no binding, and connection names resolve only at runtime.

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
**Then** the project's role-tagged bindings are shown — the role `project` items store and each role `secrets` vault — with the servers each vault holds and the environments those servers serve, any binding recorded through an informed override is shown as client-readable with the secrets it exposes named, the storage choice is presented as the security decision with the plaintext-envelope limitation stated where the binding is configured, and the user is told that the DataTug server holds the OpenVaultDB app grant for each vault and that the user's own access to credentials follows access to the project on top of that grant.

### AC: envelope-is-version-tagged (verifies REQ:credential-envelope)

**Given** a credential written by the initial implementation
**When** the stored value is inspected
**Then** it carries an algorithm or version marker identifying the plaintext variant, so a later encrypted variant is distinguishable without re-reading the secret itself.

### AC: server-reads-credentials-without-knowing-the-backend (verifies REQ:server-reads-credentials-uniformly)

**Given** credentials held in an OpenVaultDB vault backed by GitHub and one backed by Firestore
**When** the server resolves a credential for a given environment from each store
**Then** it asks dalgo for that environment's credential from that store and dalgo routes to the backing implementation, the server's resolution code path is identical across the two, and no storage-specific branch exists above dalgo. The managed-vault leg of this AC is not defined by any DataTug artifact: it is an upstream/blocked capability (REQ:vault-storage-is-pluggable), so no clause exercising it may be recorded as passing until that definition lands.

### AC: blocked-clauses-are-not-recorded-as-passing (verifies REQ:vault-storage-is-pluggable, REQ:vault-mirrors-project-namespace)

**Given** the AC clauses that depend on the remote OpenVaultDB nesting prerequisite and on the managed-vault backend
**When** this Feature's status is recorded
**Then** every such clause is enumerated in the `## Blocked Clauses` list below together with the upstream owner it waits on, and none of them is recorded as passing before that owner lands.

## Blocked Clauses

Two upstream capabilities have not landed (OpenVaultDB manifest nesting, which also blocks every GitHub-engine write under `projects/`, and the managed-vault definition), so the following AC clauses cannot be exercised.
This list is the carrier AC:blocked-clauses-are-not-recorded-as-passing points at; until the
named upstream owner lands, none of these clauses may be scaffolded or recorded as passing.

1. **AC:vault-is-laid-out-per-project** — the OpenVaultDB leg of the clause that both
   projects' credentials are distinguishable as nested records in one shared store. It waits on
   `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded: `openvaultdb-go`'s manifest
   schema cannot yet declare subcollection names and the full per-collection `record_file` configuration (`format`, the `name` template, `records_dir` and `type`, not `format` and `name` alone), so no OpenVaultDB vault produces that nested layout
   today. The local `dalgo2ingitdb` leg of the same clause is not blocked.
2. **AC:credentials-are-keyed-by-endpoint** — the OpenVaultDB leg of the clause that the record
   sits under the owning project's key-path parent in a store shared with other projects. It
   waits on the same `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded; the local
   `dalgo2ingitdb` leg is not blocked.
3. **AC:server-reads-credentials-without-knowing-the-backend** — the GitHub-engine leg waits on
   `openvaultdb/openvaultdb-go#26`: on the GitHub engine no key under `ext/datatug/projects/<project-id>/`
   can be written at all (`dalgo2ingitdb4github/scoped_collection.go:86-99`). The identical-code-path clause is exercisable across the Firestore engine and the local `ingitdb` engine until the GitHub engine joins. Separately, the managed-vault leg is blocked as follows. It
   waits on an upstream definition of what provides a managed vault and what interface it
   exposes, owned outside DataTug and named as an upstream/blocked capability by
   REQ:vault-storage-is-pluggable; no DataTug artifact defines it.

Every clause named above stays blocked until its upstream owner lands; the rest of each AC is
exercisable on its own and MUST NOT be blocked along with the named clause.

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
AC:server-authorizes-against-openvaultdb, AC:vault-authorization-stays-server-side,
AC:client-credential-is-read-only and
AC:server-reads-credentials-without-knowing-the-backend. One clause is explicitly **not**
rehearsable yet: the nested per-project addressing half of AC:vault-is-laid-out-per-project and
AC:credentials-are-keyed-by-endpoint cannot be exercised until the upstream prerequisite
recorded in `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded lands, so a stub
for its OpenVaultDB leg would assert a layout no OpenVaultDB vault produces; its local
`dalgo2ingitdb` leg is rehearsable now. The layout policy itself is rehearsed
under `dalgo-project-store`, not here. One AC is documentary and has no runtime observer, so
the Rehearse heuristic skips it: AC:blocked-clauses-are-not-recorded-as-passing asserts what
this document records in `## Blocked Clauses`, not what code does. One further AC lands on the
binding-validation surface in `datatug/backend`:
AC:zero-secrets-binding-refusal-rides-the-project-role-binding, which exercises the
zero-`secrets`-binding carrier of the refusal and its informed override.

## Not Doing / Out of Scope

- Encryption of vault contents, and the key-custody design it needs. The envelope in
  REQ:credential-envelope makes it additive; it is not built here.
- Distributing centrally stored secrets to a locally running CLI. The CLI keeps local
  secrets for now.
- Secrets in project files, in any store, under any circumstances.
- Vault registration and provisioning: owned by OpenVaultDB, not by DataTug. DataTug
  references vaults the host already exposes and never creates one; the wallet surface is
  OpenVaultDB's and only adds pointers, and provisioning is an out-of-band prerequisite the
  owner satisfies on the host.
- User-facing vault grants: access derives from project access on top of the OpenVaultDB app
  grant the server holds, so there is nothing to grant per user.
- Per-environment vault granularity finer than the server-to-vault mapping (one vault per
  DEV, UAT and PROD separately), and per-project vaults.
- The canonical project-item layout and its naming policy: owned by `dalgo-project-store`.
  This Feature names only the credential's own address and filename within that layout —
  collection `credentials` as a subcollection of the `projects` record, key path
  `ext/datatug/projects/<project-id>/credentials/<server-id>`, file
  `ext/datatug/projects/<project-id>/credentials/<server-id>/<server-id>.credential.json`, suffix
  `credential` — and MUST NOT restate the general directory shape, `records_dir`, the record
  format or the type-suffix grammar. The remote nested
  path is an upstream prerequisite recorded by `dalgo-project-store`
  REQ:remote-nesting-prerequisite-is-recorded, not something this Feature may require before
  it lands.
- A new permission model for vaults; access reuses `sharing-scopes` and
  `server-enforced-access`.

## Open Questions

- The key-custody story for the deferred encrypted envelope, which differs per runtime and
  cannot be settled without the hosted/self-hosted/CLI split being final.
- The remote nested path's upstream prerequisite — a declared OpenVaultDB collection that can
  express subcollection names and the full per-collection `record_file` configuration (`format`, the `name` template, `records_dir` and `type`, not `format` and `name` alone) — is owned and recorded by
  `dalgo-project-store` REQ:remote-nesting-prerequisite-is-recorded. This Feature has no
  separate open question on the credential layout: the key path is fixed
  (`ext/datatug/projects/<project-id>/credentials/<server-id>`) and only its remote realisation is
  pending.

---

*This document follows the https://specscore.md/feature-specification*
