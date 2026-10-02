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
**Depends-On:** git-native-projects, server-enforced-access, dalgo-project-store
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

This Feature covers the store bindings a project carries — its required `project`-role
binding plus zero or more `secrets`-role bindings, all **role-tagged and recorded on the
project record itself** — server-mediated writes, and a browser
read path. Credentials are explicitly **not** part of it: the browser is never given a
credential for a credentials store, and project files carry logical pointers only. The
canonical project layout and the DALgo persistence architecture it rests on are owned by
`dalgo-project-store`; this Feature references that layout rather than restating it.

## Problem

The initial audit recorded that the local filestore is the only fully working `Store`
backend, that GitHub-hosted projects can write but not read back
(`GhStorage.FileExists`/`OpenFile` panic), and that the store abstraction has a single
real implementation despite four named backends. The DataTug cloud store has since become a
working store for project records: `datatug/backend` creates a project in Firestore through
DALgo (`api4datatug/http_project_create.go`, `?store=firestore`; `facade4datatug/project_create.go`),
which is the source store AC:store-switch-is-an-explicit-migration moves from. Meanwhile the knowledge a DataTug user
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

### Store binding

The canonical layout of a project store — every rule about the collection and record
directory shape, the record filename and its type suffix, the records directory, a record's
format and type, and the collection/subcollection key path — is owned by
`dalgo-project-store` (REQ:canonical-project-layout, REQ:type-suffix-is-self-identifying,
REQ:hierarchy-is-a-key-path, REQ:readme-per-collection-and-record). This Feature defers to
that layout and does not restate it, and no requirement here asserts the shape of a
collection, a record key or a path. The remote (OpenVaultDB) path's upstream prerequisite —
letting a declared collection express its subcollection names and the full per-collection `record_file` configuration (`format`, `name`, `records_dir`, `type`), tracked upstream as `openvaultdb/openvaultdb-go#26` — is
likewise owned and recorded by `dalgo-project-store`
REQ:remote-nesting-prerequisite-is-recorded, and this Feature does not claim that the remote
nested layout works today.

#### REQ: project-store-binding

A project MUST have exactly one binding whose `role` is `project` — its own items store — and
MUST record every binding **on the project record**. That required `project`-role store is what
makes the bootstrap well-founded: bindings live on the project record, and the project record
lives in the `project`-role store, so a project with no `project`-role binding would have
nowhere to live. A project MAY additionally bind **zero or more** stores whose `role` is
`secrets`, and binding zero of those is legal. Every binding MUST carry a `role` field naming
what the binding is for. Recording bindings on the project record, single-sourced, is required
so that a **shared** project's bindings cannot disagree between users: no per-user copy of a
binding is authoritative, and DataTug MUST NOT treat the per-user discovery index at
`users/{uid}/ext/datatug` as the record of a binding. The legal `role` values are:

- `project` — the project's own items store, holding the project's assets. A project MUST have
  exactly one `project`-role binding, and that store MUST be chosen before the project's first
  write (REQ:store-choice-is-explicit).
- `secrets` — a store holding the project's credentials. There MAY be **zero or more**
  `secrets` bindings: zero is legal either because the credentials live in the project's own
  items store, or because the project has none yet. That placement in the items store is the
  legitimate zero-`secrets`-binding case, and it is subject to the colocation policy owned by
  `secrets-vault` REQ:client-cannot-read-credentials: refused by default, with the specific
  exposed secrets enumerated, and permitted only through the explicit informed confirmation
  (REQ:secrets-colocation-defaults-to-reject). The `secrets` role is live, not reserved.

The rule's trigger is a **credentials placement** in a store a client read credential can
reach — never the presence of ordinary project item records, and never the mere fact that the
`project`-role items store is client-readable, which it is by construction because the browser
read path reaches it (REQ:browser-read-path). The mandatory `project`-role binding MUST NOT be
rejected merely because that store is client-readable; only a credentials placement there
triggers the rule.

**The zero-`secrets`-binding override carrier.** With zero `secrets` bindings there is no
`secrets` binding for validation to reject and nowhere on one to record an override, so the
carrier is the mandatory `project`-role binding: that binding's validation performs the
refuse-by-default, the enumeration of the specific secrets by name, and the informed
confirmation, and the recorded override is attached to that binding as its explicit statement
that those named secrets are client-readable. Where a `secrets` binding exists, the carrier is
that binding, as `secrets-vault` REQ:client-cannot-read-credentials states.

An unknown `role` value MUST be rejected. **This Feature owns a binding record's one field
set, defined here and not restated by any other Feature.** A binding records, on the project
record:

- its `role`;
- the bound store's identity — a **user-supplied OpenVaultDB host base URL** and a **vault
  id** (`secrets-vault` REQ:vault-is-a-first-class-entity);
- for a `secrets` binding, the **servers it serves** — the persisted server-to-vault
  assignment whose content `secrets-vault` REQ:project-vault-binding owns, so that with one
  vault serving DEV and UAT the assignment is recorded on that binding and environments
  inherit their server's vault; and
- where an informed override under REQ:secrets-colocation-defaults-to-reject has been
  recorded, the explicit statement naming the secrets that are client-readable, attached to
  the binding that carried the confirmation — the `secrets` binding, or the `project`-role
  binding in the zero-`secrets`-binding case.

A binding MUST NOT carry the OpenVaultDB credential DataTug uses to reach the store: that
credential MUST be held server-side (REQ:ovdb-credential-is-server-only). DataTug MUST hold an
OpenVaultDB credential for this purpose and MUST NOT hold a GitHub credential for a
`project`-role store.

The selection, visibility and validation of the `secrets` role are owned by `secrets-vault`
(REQ:project-vault-binding, REQ:client-cannot-read-credentials), which cites this requirement
for the binding record's field set instead of restating it; this Feature owns that field set
and the `project` role.

### Read and write paths

#### REQ: server-mediated-writes

Every mutation of a project store MUST go through the DataTug server
(`/v0/datatug/…`, the prefix `datatug/backend` registers in
`api4datatug/routes.go`) and write through `dalgo2openvaultdb`. A browser or AI agent MUST
NOT hold a write credential for the project store, and DataTug MUST NOT write project files
by any path that bypasses the server.

#### REQ: browser-read-path

A browser MUST be able to read a project's **items store** directly from the bound store's
**OpenVaultDB host**, over that host's `/v1` data API, without the DataTug server, through a
reader module shipped in `datatug-apps`. The reader addresses the host by the binding's
recorded identity alone — its host base URL and vault id (REQ:project-store-binding) — so it
needs **no backing-repository coordinates**, never talks to the git host behind the vault, and
resolves nothing from the registry, which MUST NOT record a project's bindings or its backing
repository (REQ:registry-is-authoritative). The host's `/.well-known/openvaultdb` document is
used only to validate that the base URL is an OpenVaultDB host; it carries no repository
coordinates (`openvaultdb-com/interface/main.tsp` `ServerInfo`), and the reader MUST NOT depend
on it for any.

The reader's credential is an OpenVaultDB **application grant** obtained through the host's
Connect flow (consent → one-time code → `/token`). It is **read-only by construction** because
DataTug requests only the `records:read` capability, bound to the vault's database id. The `/token`
response does not echo the granted capabilities (`openvaultdb-go/pkg/server/connect.go:160-166`),
so read-only rests on what DataTug requested at consent, not on inspecting the token. In
`openvaultdb-go` v0.6.2 the binding's **vault id is the OpenVaultDB database id** (`/v1/databases`).
The published contract names the same thing a vault (`openvaultdb-com/interface/main.tsp`).

Grants are scoped per **root collection** of the key (`pkg/server/records.go:56-58`). Every project
item lives under `projects/...`, so the narrowest grant is `records:read:projects` on that database.
It reaches every project in the database, and any credentials placed there. The browser is NEVER
given a credential for a credentials store, and a credentials placement in a database such a
grant can reach is exactly what `secrets-vault` REQ:client-cannot-read-credentials refuses by
default.

The host MUST admit the DataTug app origin through its CORS configuration
(`openvaultdb-go/pkg/server/cors.go`). When it does not, the reader MUST fail with an error naming
the host and the CORS cause rather than an empty project. Grants expire after one hour with no
refresh (`pkg/auth/store.go`), so the reader re-runs consent when the grant expires. A
browser-held Connect flow is a public client, so it waits on upstream PKCE and `redirect_uri`
binding (`openvaultdb/openvaultdb-go#24`); see Blocked Clauses.

A credentials store is reached only by the server, which reads credentials through
`dalgo`/OpenVaultDB uniformly, unaware of whether the backing storage is GitHub, Firestore, or a
managed vault. The reader MUST be conformance-tested against the writer's output, so that the Go
writer and the browser reader cannot drift apart silently. Because this path bypasses the server,
DataTug MUST document its limits plainly: whatever the grant covers (every project in the
database) is readable by the browser, so read-scoped restrictions expressed in DataTug MUST NOT be
advertised as enforced on this path. This is the one declared exception to `server-enforced-access`'s rule that reads
reach a source through the trusted server; see REQ:project-store-read-exception.

#### REQ: project-store-read-exception

Browser-direct reads of the project store are a declared, narrowly scoped exception to
`server-enforced-access`, which forbids client-direct access to protected data sources. The
exception MUST cover reads of the project store and nothing else: every write MUST still go
through the trusted server, and every read of a protected data source MUST still go through
it. `server-enforced-access` has already been amended so that its prohibition is scoped to
protected data sources rather than to every client connection, and it already points at this
exception (its Not Doing / Out of Scope section). This Feature MUST NOT be read as widening
client-direct access anywhere else.

### Registry and store choice

#### REQ: registry-is-authoritative

The Firestore registry at `users/{uid}/ext/datatug` MUST remain the authoritative list of
a user's DataTug projects, and DataTug MUST NOT maintain a derived project index inside
the OpenVaultDB store. That registry is a per-user **discovery index only**: its shipped
shape is `UserExt.Stores map[string]*StoreBrief` keyed by store id, holding each store's
projects (store → projects) (`backend/models4datatug/project.go:79-92`). It MUST NOT record a
project's bindings and MUST NOT carry a `role` field; a binding exists only on the project
record (REQ:project-store-binding). Listing projects MUST NOT depend on a read-modify-write of
shared store state, because OpenVaultDB's MVP offers no cross-transaction isolation and no
update preconditions.

#### REQ: store-choice-is-explicit

A project's `project`-role store MUST be chosen before the project's first write and recorded
with the project. Storing one project's items in two stores at once MUST NOT be possible, and
moving a project between stores MUST be an explicit migration operation rather than an
implicit dual write.

### Credential hygiene

#### REQ: no-inline-credentials

A project asset MUST reference a connection by logical name and MUST NOT carry a
credential value, and no binding MUST carry one. `datatug validate` MUST reject a project
where any asset carries a credential. `Project.Validate()` already validates `Queries` — it
calls `p.Queries.Folders.Validate()` and `p.Queries.Items.Validate()`
(`datatug-core/pkg/datatug/project.go:123-133`), which walks every `QueryDef`
(`pkg/datatug/query.go:74-81`), and `QueryDefTarget.Validate` (`query.go:176-189`) rejects
credential material. The residual gap is that the resulting error does not name the offending
FILE, so a credential is reported without pointing at the file that carries it. This guard is
owned by `git-native-projects` REQ:secrets-never-in-project-files, with which
`heterogeneous-sources` REQ:secrets-outside-project-files shares one implementation; closing
the file-naming gap is part of that owning requirement and is NOT closed here. This Feature is
a CONSUMER of that shared guard: it MUST enforce the shared guard on the surfaces it owns and
MUST NOT claim to close the file-naming gap itself and MUST NOT define a third, independent
copy of it. Credentials themselves are stored by the `secrets-vault` Feature.

#### REQ: ovdb-credential-is-server-only

Credentials are **server-only**. A typical arrangement is two stores — an items store holding
the project's assets and a credentials store holding the project's credentials — with separate
host ACLs, but the rule does not depend on that arrangement: whatever a bound store holds, the
credential DataTug uses to reach it MUST be held in a **server-only** credential store and MUST
NEVER be written to a client-readable surface — the registry at `users/{uid}/ext/datatug`, the
project record or any other binding, or any client bundle. The browser reads the **items** store
directly with its read-only OpenVaultDB grant (REQ:browser-read-path) and is NEVER given a
credential for a credentials store; the server reads credentials through `dalgo`/OpenVaultDB
uniformly, unaware of whether the backing storage is GitHub, Firestore, or a managed vault. A
binding records exactly the field set defined once in REQ:project-store-binding — `role`, host
base URL, vault id, a `secrets` binding's servers, and any recorded informed override; a
credential on a binding or in that document is a defect, because the client reads both. A
missing or rejected credential MUST fail closed, naming the missing credential, rather than
falling back to anonymous or default access. This mirrors, for the project store's data plane,
the prohibition `secrets-vault` REQ:server-authorizes-vault-access states for vault
authorization, so the two Features are consistent.

#### REQ: secrets-colocation-defaults-to-reject

This Feature is a store-side **CONFORMANCE CONSUMER** of the colocation rule whose single
normative owner and enforcement is `secrets-vault` REQ:client-cannot-read-credentials, and
nothing else. That owning requirement defines the rule: binding validation rejects BY DEFAULT
a configuration in which a client read credential can reach a bound credentials store's
storage, the refusal enumerates the SPECIFIC secrets that would be exposed by name, and the
placement MAY be recorded only through an explicit "Are you sure?"-style confirmation step
showing that list, in which case the binding MUST state that those named secrets are
client-readable.

This Feature's obligation is conformance, not definition: the bindings it records on the
project record (REQ:project-store-binding) MUST expose the store-side facts the owner's
enforcement inspects — the bound store's identity, the credentials store's placement, and,
when an override is recorded, the statement naming the client-readable secrets — and this
Feature MUST NOT restate the rule's normative content, MUST NOT define a second, independent
enforcement, and MUST NOT be read as closing this rule on its own. Because a credentials store
lays its credentials out under the canonical project-item layout (`dalgo-project-store`
REQ:canonical-project-layout) and keys them per project and per endpoint (`secrets-vault`
REQ:vault-mirrors-project-namespace), colocating a credentials store in a repository the
client reads directly would defeat `secrets-vault` REQ:client-cannot-read-credentials; that is
the store-side consequence this Feature must make visible, and its enforcement stays with the
owner.

## Acceptance Criteria

### AC: binding-is-recorded-on-the-project (verifies REQ:project-store-binding)

**Given** a user creating a project against a GitHub repository through OpenVaultDB
**When** the project is created
**Then** the project record holds a binding naming the OpenVaultDB host base URL, vault id and `role: project`, the binding carries no OpenVaultDB credential, no GitHub credential is stored by the server for that project, and the per-user discovery index at `users/{uid}/ext/datatug` remains a store → projects brief carrying no authoritative binding and no `role` field.

### AC: bindings-are-single-sourced-on-a-shared-project (verifies REQ:project-store-binding)

**Given** a new project whose `project`-role store has not yet been chosen, and a second project shared between two users whose credentials live in its own items store
**When** each project is saved and each user reads the shared project's bindings
**Then** the first project has its `project`-role store chosen before its first write and saves with exactly one `role: project` binding and zero `role: secrets` bindings, the second records a `role: project` binding and zero `role: secrets` bindings whose validation refuses the credentials placement in its own items store by default and enumerates the specific secrets by name, the second is saved only after the user confirms through the explicit "Are you sure?"-style alert showing that list and its `project`-role binding records the statement that those named secrets are client-readable, neither project's mandatory `project`-role binding is rejected merely because its items store is client-readable, both users read the same binding set from the project record, and an unknown `role` value is rejected.

### AC: binding-record-carries-its-field-set (verifies REQ:project-store-binding)

**Given** a project with one `project`-role binding and one `secrets` binding whose vault serves that project's DEV and UAT servers
**When** the project record's bindings are inspected
**Then** each binding records its `role` and the bound store's OpenVaultDB host base URL and vault id, the `secrets` binding alone records the DEV and UAT servers it serves as the server-to-vault assignment with no independent environment-to-vault mapping stored anywhere, no binding carries an OpenVaultDB credential, and the per-user registry at `users/{uid}/ext/datatug` records no binding and no backing-repository coordinate.

### AC: a-client-cannot-write-the-store (verifies REQ:server-mediated-writes)

**Given** a signed-in browser session for a project bound to an OpenVaultDB store
**When** the client attempts to mutate project files without going through `/v0/datatug/…`
**Then** the mutation does not reach the store, and the same mutation submitted through the server succeeds.

### AC: browser-reads-a-project-without-the-server (verifies REQ:browser-read-path)

**Given** a project written to an OpenVaultDB store through the server, and a read-only OpenVaultDB grant (`records:read` only, bound to that vault's database) obtained through the host's Connect flow
**When** the browser opens the project with the DataTug server unavailable
**Then** the project's assets are read from the OpenVaultDB host's `/v1` data API addressed by the bound store's recorded host base URL and vault id alone, no backing-repository coordinate is resolved or needed and none is carried by the registry, the grant DataTug requested carries only `records:read`, no credential for a credentials store is ever accepted or used by the browser, and the parsed result matches what the conformance test derives from the Go writer's output.

### AC: project-list-comes-from-the-registry (verifies REQ:registry-is-authoritative)

**Given** a store containing projects written by more than one user
**When** a user opens their project list
**Then** the list is served from the per-user discovery index at `users/{uid}/ext/datatug`, listing performs no write to the OpenVaultDB store, and no derived project index exists inside that store to be read or updated.

### AC: store-switch-is-an-explicit-migration (verifies REQ:store-choice-is-explicit)

**Given** a project bound to the DataTug cloud store
**When** the owner moves it to an OpenVaultDB store
**Then** the move is performed as a single explicit migration, the project holds exactly one `project`-role binding afterwards, and no write is ever directed at both stores.

### AC: validator-rejects-credentials-in-project-assets (verifies REQ:no-inline-credentials)

**Given** a query file in a project store that carries an inline password
**When** `datatug validate` runs
**Then** it fails naming the file and field, and passes once the credential is replaced by a logical connection name. This is a consumer-side AC: the file-naming behaviour it asserts is delivered by the guard owned by `git-native-projects` REQ:secrets-never-in-project-files, not by this Feature.

### AC: server-authenticates-with-the-recorded-credential (verifies REQ:project-store-binding)

**Given** a project bound to an OpenVaultDB store whose OpenVaultDB credential is held only in the server-only credential store and not on the project record or in the registry
**When** the server reads the project store and then mutates it
**Then** both operations authenticate with that server-only credential, succeed while it is valid, and fail closed when it is absent or rejected, with no GitHub credential involved in either.

### AC: ovdb-credential-never-appears-in-the-registry (verifies REQ:ovdb-credential-is-server-only)

**Given** a server bound to a project's OpenVaultDB store and holding that store's OpenVaultDB credential
**When** the client-readable registry at `users/{uid}/ext/datatug`, the project record, every binding, and every client bundle are inspected
**Then** none carries the OpenVaultDB credential, which lives only in the server-only credential store, and a read attempted against the registry or a binding fails closed rather than yielding the credential.

### AC: project-store-reads-are-the-only-exception (verifies REQ:project-store-read-exception)

**Given** a protected data source and a project store
**When** a client attempts a direct read of each, and a write to each
**Then** the direct read of the protected data source is refused, the direct read of the project store succeeds, and both writes are refused unless they arrive through the trusted server.

### AC: colocation-refused-by-default-with-informed-override (verifies REQ:secrets-colocation-defaults-to-reject)

**Given** a project store the client reads directly with its read credential, and a proposed credentials-store binding whose storage that credential can reach
**When** the binding is saved
**Then** binding validation refuses it by default, the refusal enumerates the specific secrets that would be exposed rather than merely asserting that credentials would be client-readable, and the binding is accepted only after the user confirms through an explicit "Are you sure?"-style alert showing that list; if overridden, the recorded binding states that those named secrets are client-readable and does not represent them as client-inaccessible. This is a store-side conformance AC: the refusal and its informed override are enforced by `secrets-vault` REQ:client-cannot-read-credentials, so this AC observes that outcome and cannot be verified independently of that enforcement — it MUST NOT be recorded as passing before that enforcement exists, a dependency enumerated in `## Conformance Dependencies` and carried by AC:colocation-conformance-depends-on-secrets-vault, which also covers the zero-`secrets`-binding case in which the `project`-role binding is the carrier.

### AC: colocation-conformance-depends-on-secrets-vault (verifies REQ:secrets-colocation-defaults-to-reject)

**Given** the AC clauses of this Feature that assert the colocation refusal, its named-secret enumeration and its informed override, including the zero-`secrets`-binding case carried by the `project`-role binding
**When** this Feature's status is recorded
**Then** every such clause — the whole of AC:colocation-refused-by-default-with-informed-override and the zero-`secrets`-binding clause of AC:bindings-are-single-sourced-on-a-shared-project — is enumerated in the `## Conformance Dependencies` list below together with the owning requirement it waits on, and none of them is recorded as passing before that enforcement exists.

### AC: blocked-clauses-are-not-recorded-as-passing (verifies REQ:browser-read-path, REQ:project-store-binding)

**Given** the AC clauses of this Feature that wait on an upstream `openvaultdb-go` change
**When** this Feature's status is recorded
**Then** every such clause is enumerated in the `## Blocked Clauses` list below together with its upstream owner, and none of them is recorded as passing before that owner lands.

## Blocked Clauses

Two upstream `openvaultdb-go` changes have not landed. This list is the carrier
AC:blocked-clauses-are-not-recorded-as-passing points at. Until each owner lands, none of these
clauses may be scaffolded or recorded as passing. The rest of each AC is exercisable and MUST NOT
be blocked along with the named clause.

1. **GitHub-engine legs** of AC:server-authenticates-with-the-recorded-credential,
   AC:a-client-cannot-write-the-store, AC:browser-reads-a-project-without-the-server,
   AC:store-switch-is-an-explicit-migration and AC:project-store-reads-are-the-only-exception wait on
   `openvaultdb/openvaultdb-go#26`. On the GitHub engine no key under `ext/datatug/projects/<project-id>/` can be
   written, because the manifest-built definition has no subcollections
   (`dalgo2ingitdb4github/scoped_collection.go:86-99`). These ACs are exercised on an OpenVaultDB
   host using its local `ingitdb` engine, once `dalgo-project-store` proves that engine serves the
   canonical layout.
2. **The browser-held Connect flow** of AC:browser-reads-a-project-without-the-server waits on
   `openvaultdb/openvaultdb-go#24` (PKCE, and binding `redirect_uri` at code exchange). A public
   browser client must not rely on an authorization code that anyone who intercepts it can
   exchange. The reader and its conformance test against the Go writer's golden fixture are
   exercisable with a grant obtained out of band.

## Conformance Dependencies

This Feature is a store-side **conformance consumer** of the colocation rule whose single
normative owner and enforcement is `secrets-vault` REQ:client-cannot-read-credentials, and of
nothing else. `secrets-vault` depends on this Feature and this Feature MUST NOT depend on
`secrets-vault` — that edge would close a cycle — so the enforcement this Feature waits on is
**not reachable by a dependency edge from here**; it is recorded as an explicit conformance
dependency instead. The following clauses are what the carrier
AC:colocation-conformance-depends-on-secrets-vault points at, and until the named enforcement
exists it may not be scaffolded or recorded as passing.

1. **AC:colocation-refused-by-default-with-informed-override** — waits on `secrets-vault`
   REQ:client-cannot-read-credentials: the refuse-by-default rejection of a credentials
   placement a client read credential can reach, the enumeration of the specific secrets by
   name, and the informed override recorded on the binding that carried the confirmation.
   This Feature supplies only the store-side carrier of that override on the project record —
   the `secrets` binding, or the `project`-role binding when the project has zero `secrets`
   bindings (REQ:project-store-binding) — and MUST NOT define a second enforcement of the rule
   while it waits.
2. **AC:bindings-are-single-sourced-on-a-shared-project** — the zero-`secrets`-binding clause:
   that the second project's credentials placement in its own items store is refused by
   default with the specific secrets enumerated by name, is saved only after the explicit
   "Are you sure?"-style confirmation, and records the client-readable statement on its
   `project`-role binding. It waits on the same `secrets-vault`
   REQ:client-cannot-read-credentials enforcement. The rest of that AC — one `role: project`
   binding chosen before the first write, zero `role: secrets` bindings, both users reading one
   binding set from the project record, the `project`-role binding not rejected merely for being
   client-readable, and an unknown `role` rejected — is exercisable on its own and MUST NOT be
   gated along with the named clause.

This Feature MUST NOT record the clauses above as passing before that enforcement lands, and
MUST NOT restate the rule's normative content in the meantime. A verifier that reads the
clause as independently verifiable is reading it wrongly: the dependency is real, named and
one-directional, and it is recorded here rather than assumed.

## Rehearse Integration

Every AC has a concrete behavioural surface — project-record bindings, registry briefs,
HTTP requests, validator exit behaviour, and a reader-versus-writer comparison — and
all are testable in principle. This Feature defines a capability spanning several repositories,
and each AC is exercised against a concrete implementation: server-mediated writes and the
store binding in `datatug/backend`, the browser read path in `datatug-apps`, the validator
in `datatug-cli`, and the OpenVaultDB data plane in `openvaultdb-go`. Stub scaffolding is
therefore deferred to those implementations, so the stubs track real runnable surfaces
rather than an abstract capability — consistent with `ai-query-builder` and
`serve-brokered-query-builder`. One AC carries an obligation that must survive that
deferral: AC:browser-reads-a-project-without-the-server is a conformance test between the
Go writer and the browser reader, and it must be authored where both sides are runnable. The
canonical layout itself — including the remote nested path and its recorded upstream
prerequisite — is rehearsed under `dalgo-project-store` and MUST NOT be scaffolded as an
obligation of this Feature. Concretely: the binding and credential ACs
land on `datatug/backend` HTTP and credential-resolution surfaces;
AC:project-list-comes-from-the-registry
on the Firestore registry; AC:browser-reads-a-project-without-the-server on the
`datatug-apps` reader plus its conformance test; AC:validator-rejects-credentials-in-project-assets
on `datatug-cli` exit behaviour as a consumer of the guard owned by `git-native-projects`; and
AC:colocation-refused-by-default-with-informed-override on the store-side conformance surface
of the refusal and informed override enforced by `secrets-vault`, which cannot be verified
until that enforcement exists; and AC:binding-record-carries-its-field-set on that same
`datatug/backend` binding surface. One AC is documentary and has no runtime observer, so the
Rehearse heuristic skips it: AC:colocation-conformance-depends-on-secrets-vault asserts what
this document records in `## Conformance Dependencies`, not what code does. AC:blocked-clauses-are-not-recorded-as-passing is documentary in the same way: it asserts what
`## Blocked Clauses` records. One obligation is explicitly **not**
this Feature's: `git-native-projects` AC:github-project-reads-back-what-it-wrote depends on
`GhStorage` read methods that still panic, so no stub authored here can satisfy it and it must
not be attached to this Feature's plan as a Supporting AC.

## Not Doing / Out of Scope

- More than one `project`-role store binding per project, and per-environment repositories. A
  project holds exactly one `project`-role binding; `secrets`-role bindings are additive and
  legal (REQ:project-store-binding), but more than one items store per project is not built
  here.
- The canonical project layout and its naming policy, and the DALgo persistence architecture
  they rest on: owned by `dalgo-project-store`. This Feature references that layout and MUST
  NOT define its own rule for the collection shape, the record filename, `records_dir`, the
  type-suffix grammar or the key path.
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
  path reads through the OpenVaultDB host's `/v1` API and does not use `GhStorage`, and adding
  `GhStorage` read support is a separate task owned by `git-native-projects` — an explicit
  prerequisite for that AC, not an obligation of this Feature.
- Client-direct writes to the project store, which `git-native-projects` already excludes.
- Direct browser-to-database query execution.

## Open Questions

- How is a `project`-role binding to the **DataTug cloud store** recorded? The binding field set
  (host base URL plus vault id) describes an OpenVaultDB store. A cloud-store project today is
  identified by the backend's `firestore` store id (`models4datatug.FirestoreStoreID`). Until this
  is settled, AC:store-switch-is-an-explicit-migration treats the source as that store id and the
  target as the one OpenVaultDB `project`-role binding.

- The binding endpoint shape is settled, not provisional: it is a **user-supplied
  OpenVaultDB host base URL**, validated by the `/.well-known/openvaultdb` discovery document
  OpenVaultDB serves (`openvaultdb-go/pkg/server/server.go`), with the data plane addressed
  under that base at `/v1`. `/ext/v0` is a platform convention for a **platform-owned** hosted
  endpoint, which is out of scope (see Not Doing / Out of Scope); it is not required of, and
  does not constrain, a user-supplied host.
- The remote (OpenVaultDB) nested layout is not an open question of this Feature: its upstream
  prerequisite — a declared collection able to express subcollection names and the full per-collection `record_file` configuration (`format`, `name`, `records_dir`, `type`), tracked upstream as `openvaultdb/openvaultdb-go#26`
  — is owned and recorded by `dalgo-project-store`
  REQ:remote-nesting-prerequisite-is-recorded, and the canonical layout it must realise is
  owned there too.

---

*This document follows the https://specscore.md/feature-specification*
