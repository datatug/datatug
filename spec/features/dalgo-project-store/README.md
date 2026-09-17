---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Project store on DALgo rails

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dalgo-project-store?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dalgo-project-store?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dalgo-project-store?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/dalgo-project-store?op=request-change) |
**Status:** Approved
**Date:** 2026-09-16
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Depends-On:** git-native-projects
**Grade:** A

## Summary

`datatug-core`'s **project persistence** moves off the bespoke `pkg/storage/filestore`
filesystem code and onto DALgo, with one canonical on-disk structure. A project stored
locally and a project reached through OpenVaultDB become the *same* code path with a
different driver, and the stored form stays human-browsable — which is the product
("own your data"). This **partially** implements `git-native-projects`
REQ:project-store-pluggable-backends: it commits to the driver-agnostic code path — no driver
branch above the DALgo boundary — and to the local-directory backend concretely, and it owns
the canonical layout and layout policy every backend rides on. The backends it does not deliver
are owned follow-on work of this tree, landed incrementally and named per backend in
Not Doing / Out of Scope: the remote OpenVaultDB binding by `openvaultdb-project-store`,
GitHub-native read parity via the GitHub-backed driver, and DataTug Cloud via
`dalgo2firestore`.

This Feature does **not** own the OpenVaultDB remote binding (see
`openvaultdb-project-store`), secrets vaults (see `secrets-vault`), or user-data
semantics. It owns the DALgo boundary for the project store, the canonical layout and its
naming policy, the retirement of the filestore's own transaction machinery, and the
recording of the upstream prerequisites the remote and local paths depend on.

Driver facts this Feature relies on, stated once: the local driver is `dalgo2ingitdb`
v0.6.1, whose `FORMAT.md` is the cross-language contract with `ingitdb-ts`, validated by
shared golden fixtures; `record_file` there supports
`format: yaml|yml|json|markdown|csv|jsonl|toml|ingr` and a `record_file.name` template
containing `{key}`. Omitting the `$records` segment is the job of
`record_file.records_dir: '.'`, and that property belongs to `ingitdb-go/ingitdb >= v0.6.0`
(`RecordFileDef.RecordsDir`, honoured by `RecordFileDef.RecordsBasePath()`). The local
prerequisite is a **consumer pin, not a driver bump**: `dalgo2ingitdb` v0.6.1's own `go.mod`
declares only a *minimum* of `ingitdb-go/ingitdb` v0.5.2, so the resolved version is decided
by the consuming module — a consumer that does not require `ingitdb-go/ingitdb >= v0.6.0`
directly resolves v0.5.2, whose `RecordFileDef` has no `RecordsDir` field, and the strict
`KnownFields(true)` reader (`ingitdb-go/ingitdb@v0.5.2/validator/def_validator.go:35`; the
same reader is at the identical path and line in v0.6.0 and v0.6.1, so the pin is the only difference)
then rejects a definition carrying `record_file.records_dir` as an unknown field. A consumer that
resolves v0.6.0 or later accepts the key: `RecordsDir` is declared at `record_file_def.go:29`,
validated at `:61`, and `RecordsBasePath()` returns `""` for `"."` at `:131-137`. This is
already satisfied where it matters today — `datatug-cli/go.mod:39` directly requires
`github.com/ingitdb/ingitdb-go/ingitdb` at v0.6.0 or later, `go list -m` in `datatug-cli` resolves
v0.6.0 or later, there is no `go.work`, no `vendor/` and no `replace`, and
`datatug-cli/pkg/dbcopy/url.go:359` constructs the local driver as
`dalgo2ingitdb.NewDatabase(r.Path, validator.NewCollectionsReader(), opts...)`, so the
`validator` in that build is that version's. That pin belongs to whichever module **constructs the
driver** — the CLI, which already satisfies it, and the backend as the other edge — and not
to `datatug-core`, which by REQ:project-store-is-a-dalgo-database MUST NOT declare any DALgo
driver. Record paths are resolved by `db.reader.ReadDefinition`
(`dalgo2ingitdb@v0.6.1/database.go:321`) and expanded through
`colDef.RecordFile.RecordsBasePath()` (`query.go:587`, `record_io.go:118`); the definition
reader is strict (`KnownFields(true)`), and `schema_reader.go:105` is the `DescribeCollection`
path only, a non-strict `yaml.Unmarshal` that does not touch record resolution. The remote
driver is
`dalgo2openvaultdb` (the build resolves `v0.3.2-0.20260909113030-cacd422564fc`; its API was
verified in the cached `v0.3.1`), offering `Get`/`Exists`/`GetMulti`, read-only and read-write
transactions, `Set`/`Insert`/`Delete`/field-level `Update`, all `*Multi` forms, auto-ID
insert and structured query, over the key-path shape
`collection/id/subcollection/id`. Nesting is resolved by `resolveScopedCollection` in
`dalgo2ingitdb4github/scoped_collection.go` on the remote path and by the same-shaped
`resolveScopedCollection` in the pinned local driver's `dalgo2ingitdb/scoped_collection.go`.

**Dependency policy: latest releases.** Every module in this tree requires the latest release of
the storage stack, and an older version still pinned is a defect to fix, not a baseline to
design around. As of 2026-09-17 the latest releases are `dal-go/dalgo` v0.80.5, `dal-go/record`
v0.1.3, `ingitdb-go/ingitdb` v0.6.1, `dalgo2ingitdb` v0.6.1, `dalgo2ingitdb4github` v0.2.8,
`openvaultdb-go` v0.6.2 and `dalgo2openvaultdb` v0.3.1, with `openvaultdb-go` v0.6.2 itself
requiring the newer `dalgo2openvaultdb` pseudo-version above. The facts cited in this Feature
were re-checked at those releases. `dalgo2ingitdb` v0.6.1 still declares only an `ingitdb-go/ingitdb`
v0.5.2 floor, so the consumer pin below remains necessary under this policy. It is satisfied
by requiring the latest `ingitdb-go/ingitdb` directly.

## Problem

A knowledge asset that cannot be reviewed, diffed and moved between local, hosted and
Git-backed stores is a liability, not an asset — the reason `git-native-projects` commits
to Git-native storage. Today the only fully working backend is the bespoke filestore, so
"a project wherever its owner keeps it" is not reachable without writing a second
implementation.

The cost is concrete and localised in `datatug-core/pkg/storage/filestore`:

- **The layout and its naming grammar are hand-written.** `JsonFileName(id, suffix)`
  decides file names, and each store re-derives its own directory shape; four item types
  (entity, board, db model, env db catalog) accept **two** on-disk layouts with a
  `nested wins; mismatch is a load error` rule that exists only in this code — a flat
  `<collection>/<id>.<suffix>.json` file versus a per-record directory, with board's nested
  legacy filename being the fixed `board.json` rather than an id-prefixed name
  (`boards_store.go:45-47`). `environment` is a **competing-filename** case, not a
  flat-vs-nested one: its two accepted on-disk names are `environment-summary.json` and
  `<id>.env.json`, both at the same depth inside the same per-environment directory
  (`environments_store.go:25-31`).
- **Transaction and locking machinery is bespoke**: `query_txn.go` (~48 KB) and
  `query_lock.go`, plus a reserved on-disk `.dt-query-txn/` namespace that is part of the
  project layout and therefore leaks an implementation detail into the user's repository.
- **The same item has two homes.** Recordset definitions are read by the loader from the
  project-relative `data/recordsets/...` path while a flat
  `recordsets/<id>.recordset.json` also exists — a duplication the move collapses.
- **Defects worth not carrying forward:** in the demo repo
  `datatug-project.json`'s `id` disagrees with its directory name, and
  `newDirProjectItemsStore`'s lister reads `.datatug-<itemFileSuffix>.json`
  (`project_items_store.go:144`) while its saver writes `summaryFileName` (`:184`), so
  listing a directory-stored item cannot find what the saver wrote.

The same move also settles the **user data** format, which the demo currently stores
inconsistently with the target: `demo-project-1/data/ingitdb/support-notes` writes each
support note as an individual `format: yaml` / `type: map[string]any` file under that
collection's `$records/`, whereas the target keeps one inGitDB database root per project at
`projects/<project-id>/data/ingitdb/`, with each dataset at
`projects/<project-id>/data/ingitdb/<dataset>/<dataset>.ingr`. The target deliberately chooses
`format: ingr` with a
`ListOfRecords` or `MapOfRecords` type instead — a format that is inherently a multi-record
container (`RecordFileDef.Validate()` rejects `SingleRecord` for `RecordFormatINGR`) and
therefore keeps a dataset in one named file, which is the compact, recursion-free shape user
data needs. That choice is a stated decision in REQ:user-data-uses-a-compact-ingr-file, not an
unrecorded side effect of the item layout.

DALgo already supplies what the bespoke code hand-rolls — collection definitions, a
`record_file` template, transactions, and a key-path hierarchy — on both the local and the
remote driver. What is missing is the decision that the project store *is* a DALgo
database, and a single canonical shape that both drivers can serve.

## Behavior

### DALgo boundary

The project store is a DALgo database. `datatug-core` addresses a project through
`dal-go/dalgo` interfaces only; a driver is constructed at the edge, and both the local
`dalgo2ingitdb` stack and the remote `dalgo2openvaultdb` stack present the same interface
to core. The target shape already holds today — `datatug-core/go.mod` depends on
`dal-go/dalgo` (v0.80.4 at `28e57ff`; latest release v0.80.5) with **no driver** — and this Feature keeps it that way rather than
reintroducing a concrete backend into core.

#### REQ: project-store-is-a-dalgo-database

Project persistence MUST be reached through the `dal-go/dalgo` interfaces. `datatug-core`
MUST depend on the DALgo interfaces and MUST NOT depend on any driver; the driver is
selected at the edge (CLI, backend), and no code above DALgo may branch on which storage
backs the project. The testable outcome is the **absence of a driver-type branch above the
DALgo boundary**, not the existence of two backends: one DALgo-facing code path MUST serve a
project under different local configurations, and nothing above the boundary may be able to
tell which concrete driver is in use.

### Canonical layout

A **project item** is one of a project's own metadata/definition records: the project record
itself and the records of the item types listed in the canonical layout below. It explicitly
**excludes user data** — dataset records are the one deliberate exception (see User data is
the deliberate exception) — and "project item" is used consistently throughout this Feature,
never as a synonym for "asset".

One shape for every project item, whether the bytes live on disk or behind OpenVaultDB.
`projects` is not a special case: a project record's directory **is** the project
directory, and that same directory holds its subcollections — exactly `FORMAT.md`'s
interleaving of parent-record ids and subcollection names.

#### REQ: canonical-project-layout

There MUST be one collection per project-item type and one directory per record; the
record file MUST be `<key>/<key>.<type-suffix>.json`; `record_file.records_dir` MUST be `.`
so no `$records` segment appears; `format` MUST be `json`; and `type` MUST be
`map[string]any`. A collection's `record_file` is:

```yaml
record_file:
  name: '{key}/{key}.<type-suffix>.json'
  format: json
  type: map[string]any
  records_dir: '.'
```

`record_file.records_dir: '.'` is the mechanism that omits the `$records` segment, and it is
a property of `ingitdb-go/ingitdb >= v0.6.0`; `RecordsDir` is a field of
`RecordFileDef`, so the key is nested inside `record_file` and a sibling `records_dir` is
rejected by the strict definition reader (see the driver facts above). The layout is
producible on the local path today, provided the consuming module resolves
`ingitdb-go/ingitdb >= v0.6.0` — the prerequisite recorded in
REQ:local-records-dir-prerequisite-is-recorded — because `RecordsBasePath()` then returns no
`$records` segment for `records_dir: '.'` (`record_file_def.go:131-137`).

The canonical layout this produces is one tree, shared by every backend:

```
<store root>/
  projects/
    <project-id>/
      <project-id>.datatug-project.json
      credentials/<server-id>/<server-id>.credential.json
      queries/<query-id>/<query-id>.query.json
      entities/<entity-id>/<entity-id>.entity.json
      environments/<env-id>/<env-id>.env.json
      environments/<env-id>/servers/<server-id>/<server-id>.server.json
      environments/<env-id>/catalogs/<catalog-id>/<catalog-id>.db.json
      dbmodels/<model-id>/<model-id>.dbmodel.json
      boards/<board-id>/<board-id>.board.json
      recordsets/<set-id>/<set-id>.recordset.json
      folders/<folder-id>/<folder-id>.folder.json
      dbdrivers/<driver-id>/
        <driver-id>.dbdriver.json
        dbservers/<server-id>/<server-id>.dbserver.json
      data/ingitdb/<dataset>/<dataset>.ingr
    <project-id-2>/
      ...
```

The inGitDB on-disk contract also requires schema metadata the demo carries:
`<store root>/.ingitdb/root-collections.yaml` and a per-collection
`.collection/definition.yaml`. Those are backend-owned metadata — the local driver's own
schema files, not records — and are outside the record layout this Feature declares; they are
neither project items nor part of the per-record directory shape above.

The `projects` collection MAY hold records for more than one project: a store is **not**
per-project, and a second project written to the same store MUST be addressed under its own
`projects/<project-id>` subtree rather than colliding with the first.

The credential location is part of this canonical layout. Collection `credentials` is a
subcollection of the `projects` record, key path
`projects/<project-id>/credentials/<server-id>`, file
`projects/<project-id>/credentials/<server-id>/<server-id>.credential.json`, suffix
`credential`. This Feature owns that naming policy so `secrets-vault` needs only to name the
address; credential semantics remain owned by `secrets-vault`
REQ:vault-mirrors-project-namespace. Which store realises the subcollection is decided by the
**server-to-vault assignment carried on the servers field of the role-`secrets` binding that
covers the credential's server** (`secrets-vault` REQ:project-vault-binding): a **vault
store**, or, when the project has **zero** `secrets` bindings, the project's own items store.
That second placement is the legitimate zero-binding case and is subject to the colocation
policy owned by `secrets-vault` REQ:client-cannot-read-credentials. A credential record MUST
NOT be stored inside another project-item type's collection directory — `credentials` is its
own collection, so no credential file may sit under the `queries/`, `entities/`,
`environments/`, `dbmodels/`, `boards/`, `recordsets/`, `folders/`, `dbdrivers/` or `data/`
directories.

#### REQ: type-suffix-is-self-identifying

Every project-item filename MUST carry a type-specific type suffix, so a file is
identifiable when detached from its directory (an email attachment, a log line), and no
collection MAY use a bare `{key}.json`. What MUST be distinct is each collection's **dotted
suffix as it appears in its `record_file.name` template** — the trailing `.<suffix>.json`
segment. No collection's template MAY end with another collection's dotted suffix. A bare
substring overlap between two suffixes, such as `server` inside `dbserver`, is explicitly
harmless and is not a violation.

That is the rule the mechanism actually enforces: the discriminator is a regex built from
the `record_file.name` template (`dalgo2ingitdb4github@v0.2.8/query.go:255-259`), and it
matches the **dotted** ending exactly — so `x.dbserver.json` does not match
`\.server\.json$`, and the `server`/`dbserver` substring overlap cannot cross-match because
the preceding dot makes the two patterns disjoint. A distinct dotted suffix per collection is
what keeps the scheme unambiguous on the recursive GitHub reader: on the remote path the
reader lists the collection directory recursively and filters only by that template-derived
regex, so the per-collection dotted suffix is what prevents a collection from matching its
descendants' files. Locally the same property is belt-and-braces: the template becomes a glob
(`strings.ReplaceAll(name, "{key}", "*")`), so glob depth equals template depth and
`{key}/{key}.query.json` matches exactly one level and never `README.md` or a deeper
subcollection.

#### REQ: hierarchy-is-a-key-path

Hierarchy MUST be expressed as collection/subcollection key paths
(`projects/<project-id>/queries/<query-id>`), and no record key MAY contain a path
separator. The relation is `collection/id/subcollection/id`, the shape
`dalgo2openvaultdb` addresses and the shape `resolveScopedCollection` resolves to a
concrete in-repo directory on both driver legs — `dalgo2ingitdb4github/scoped_collection.go`
on the remote/GitHub path and `dalgo2ingitdb/scoped_collection.go` on the pinned local path,
which interleaves parent-record ids with subcollection names the same way.

The canonical key paths this Feature owns are:

- `projects/<project-id>` — the project record
- `projects/<project-id>/queries/<query-id>`
- `projects/<project-id>/entities/<entity-id>`
- `projects/<project-id>/environments/<env-id>`
- `projects/<project-id>/environments/<env-id>/servers/<server-id>`
- `projects/<project-id>/environments/<env-id>/catalogs/<catalog-id>`
- `projects/<project-id>/dbmodels/<model-id>`
- `projects/<project-id>/boards/<board-id>`
- `projects/<project-id>/recordsets/<set-id>`
- `projects/<project-id>/folders/<folder-id>`
- `projects/<project-id>/dbdrivers/<driver-id>`
- `projects/<project-id>/dbdrivers/<driver-id>/dbservers/<server-id>` — a db server is
  placed under its db driver, with record file `.../<server-id>.dbserver.json`
- `projects/<project-id>/credentials/<server-id>` — realised in the vault store when the
  project has a `secrets` binding, and in the project's own items store only when the project
  has **zero** `secrets` bindings, subject to the colocation policy owned by `secrets-vault`

#### REQ: readme-per-collection-and-record

A `README.md` MUST be present at each collection directory and inside each record
directory, and its presence MUST NOT affect record reads.

#### REQ: existing-filenames-stay-recognisable

Type suffixes MUST be inherited from the existing `JsonFileName` constant set
(`datatug-core/pkg/storage/file_names.go`) wherever one exists, so migrated repositories
keep recognisable filenames.

| Type | Suffix | Origin |
|---|---|---|
| project | `datatug-project` | new — matches the name the existing creator already labels |
| query | `query` | inherited |
| entity | `entity` | inherited |
| environment | `env` | inherited |
| env server | `server` | inherited |
| env catalog | `db` | inherited |
| db model | `dbmodel` | inherited |
| board | `board` | inherited |
| recordset | `recordset` | inherited |
| project db server | `dbserver` | inherited |
| folder | `folder` | new |
| db driver | `dbdriver` | new |
| credential | `credential` | new — the suffix and its naming policy are owned by this Feature; how a credential is realised and stored is owned by `secrets-vault` (a vault store when the project has a `secrets` binding, otherwise the project's own items store) |

`credential` is a project item by the shared rule but is realised in a **vault store**
whenever the project has a `secrets` binding — the vault the server-to-vault assignment on
that binding's servers field selects (`secrets-vault` REQ:project-vault-binding) — and in the
project's own items store only in the zero-`secrets`-binding case, subject to the colocation
policy owned by `secrets-vault`: its collection and key path are declared in the canonical
layout above (`projects/<project-id>/credentials/<server-id>`, file
`.../<server-id>.credential.json`), so this Feature owns the suffix and naming policy while
`secrets-vault` REQ:vault-mirrors-project-namespace owns credential semantics.

`dbserver` is not homeless: it is placed under its db driver in the layout above
(`projects/<project-id>/dbdrivers/<driver-id>/dbservers/<server-id>`, file
`.../<server-id>.dbserver.json`).

`columns`, `objects` and `refs` are inherited suffixes with no collection in the current
`ProjectStore` interface; they need a home or retirement (see Open Questions).

### Store-level guarantees

#### REQ: store-preserves-foreign-files

A store is a shared container, not a per-project sandbox, and every write this Feature performs
MUST honour two rules about what it does *not* touch:

- **Foreign-file preservation.** Pre-existing unrelated files in the store MUST be unchanged
  byte for byte.
- **Cross-project non-interference.** Two projects written to the same store MUST have
  disjoint records and disjoint user data — each project's records and user data live only
  under its own `projects/<project-id>` subtree — and neither write MUST alter the other
  project's files.

Both rules were previously carried by Acceptance Criteria of the sibling
`openvaultdb-project-store` Feature — foreign-file preservation by its
`AC:assets-land-under-the-project-namespace`, and cross-project non-interference by its
`AC:two-projects-share-one-store`. That Feature's decomposition removed those ids and
re-homed the canonical layout here, which left the two rules owned by nobody. They are
re-homed here so the ownership is explicit and the provenance is traceable rather than
assumed.

### User data is the deliberate exception

User data (dataset records) is **not** a project item and does not use the canonical item
layout. It stays compact: one multi-record file per dataset rather than a directory per
record. `RecordFileDef.Validate()` rejects `SingleRecord` for `RecordFormatINGR` and
requires `ListOfRecords` or `MapOfRecords` — the format is inherently a multi-record
container, which also avoids the recursion question entirely because these types read one
named file. Because a dataset belongs to a project while a store MAY hold several projects,
a project's user data lives in one inGitDB database root scoped per project at
`projects/<project-id>/data/ingitdb/`, matching the legacy project-relative root
(`demo-project-1/data/ingitdb/`, which holds `.ingitdb/root-collections.yaml` and one
directory per collection). Each dataset is a collection inside it, so the dataset file is
`projects/<project-id>/data/ingitdb/<dataset>/<dataset>.ingr`.

#### REQ: user-data-is-not-a-project-item

User data MUST NOT be stored as a project item: the per-project `data/` subtree —
`projects/<project-id>/data/`, whose inGitDB database root is
`projects/<project-id>/data/ingitdb/` — MUST be reserved for that project's user data, so
one project's user data is never reachable by, or able to collide with, another project's;
and no project item MAY be stored under a project's `data/` subtree.

#### REQ: user-data-uses-a-compact-ingr-file

A dataset MUST be stored compactly as one multi-record file per dataset under the project's
user-data database root: `projects/<project-id>/data/ingitdb/<dataset>/<dataset>.ingr` MUST
carry `format: ingr` and `type` MUST be `ListOfRecords` (`'[]map[string]any'`) or
`MapOfRecords` (`'map[$record_id]map[$field_name]any'`), so a dataset is
one named file rather than a directory per record. The dataset collection's `record_file`
MUST declare `name: '<dataset>.ingr'` with **no** `{key}` placeholder, so
`RecordFileDef.RecordsBasePath()` adds no `$records` segment for a user-data collection and the
dataset file keeps the asserted path; a `{key}` template would reintroduce a `$records/`
directory for that collection and break it, which is the opposite of the compact shape this
requirement mandates. This is a deliberate decision, not a consequence of the item layout:
the demo stores the same dataset as individual
`format: yaml` / `type: map[string]any` files under the collection's `$records/`, while
`ingr` with a multi-record type is chosen because `RecordFileDef.Validate()` rejects
`SingleRecord` for `RecordFormatINGR` and the multi-record container reads one named file.

### Recordsets

#### REQ: recordset-definitions-are-a-project-item

Recordset definitions MUST be project items at
`projects/<project-id>/recordsets/<set-id>/<set-id>.recordset.json`; the two existing
project-relative locations — the loader's `data/recordsets/...` path
(`loader_recordsets.go:22`) and the flat `recordsets/<id>.recordset.json` — MUST be removed
by the move, collapsing the duplication.

### Retiring the bespoke machinery

#### REQ: bespoke-transaction-machinery-is-retired

The filestore's own transaction/locking machinery (`query_txn.go`, `query_txn_slots.go`,
`query_lock.go`) MUST be replaced by DALgo transactions,
and the reserved on-disk transaction namespace (`.dt-query-txn/`) MUST no longer be part
of the project layout.

### Hard cut-over

Founder decision, 2026-09-17: *"hard cut-over - no backward compatability is needed"* and
*"writes should be through dalgo"*. The move is a full cut-over. No compatibility shim, dual
read or feature flag keeps the old layout alive.

The consumer inventory, taken 2026-09-17 across every repository in the `datatug` organisation,
has two parts.

**Go importers of `filestore`:**
- `datatug-cli` `866b6b0`:
  - `apps/datatugapp/commands`: `cmd_validate.go`, `project_base_command.go`, `cmd_scan_db.go`,
    `cmd_serve.go`, `cmd_init_project.go`, `cmd_render.go` and their tests
  - `apps/datatugapp/datatugui/dtproject`: `create_project_screen.go`, `datatug_demo_project.go`,
    `projects_ui.go`, `add_to_github_repo.go`
  - `pkg/server`: `http_server.go` and endpoint tests
  - `pkg/api` tests
- `datatug-project-validator`: `main.go`
- `datatug-demo-projects`: `tests/demo_project_1_test.go`

**Legacy-layout readers and describers outside Go:**
- `datatug-apps`, `libs/datatug/main/src/lib/services`:
  - `repo/datatug-store.service.github.ts`
  - `repo/github/github-project-reader.service.ts`
  - `repo/github/github-project-create.service.ts`
  - `unsorted/environment.service.ts` (`/environment-summary`)
  - their specs
- `ai-plugin`: `skills/datatug-init/SKILL.md`

`datatug/backend` imports no filestore. The validator and the demo tests pin an older
`datatug-core`, so "no callers" is judged at the latest `datatug-core` across all of these, not
at their pins. The DataTug-owned repository still stored in the legacy layout is
`datatug/datatug-demo-projects`.

#### REQ: filestore-is-cut-over

`datatug-core/pkg/storage/filestore` MUST be deleted, and every consumer in the inventory MUST
read and write projects through the DALgo-backed store (Go consumers) or read and describe only
the canonical layout (`datatug-apps` readers, `ai-plugin` docs). Afterwards no module in this tree imports it. No reader
MAY tolerate a legacy layout: the flat `<collection>/<id>.<suffix>.json` files, the
flat-vs-nested dual layouts, `environment-summary.json`, `data/recordsets/`, `$records/`
directories, `.datatug-<suffix>.json` summaries and legacy `*.sql.json` queries are unsupported.
Opening one MUST fail with a typed error that names the unsupported path. It MUST NOT load
partially or silently. Detection is a **read** of the directory tree, which DALgo cannot list,
so it happens at the edge that opens a project (the CLI's open path) before the driver is
constructed. It is never a write outside DALgo. DataTug-owned repositories are rewritten to the canonical layout in the
same cut-over, by a one-off converter that reads with `filestore` before it is deleted and writes
through the DALgo store, so nothing DataTug ships depends on the old layout. The migration defects named
in Problem die with it: the demo `datatug-project.json` id mismatch is fixed by that rewrite, and
the `newDirProjectItemsStore` lister/saver mismatch is deleted with the package.

#### REQ: query-write-guards-survive-the-cut-over

The query write guarantees that `executable-knowledge-library`
AC:query-pair-storage-guards-writes established on the filestore MUST hold on the DALgo-backed
store, because retiring the machinery MUST NOT retire the behaviour:

- the `datatug.RevisionedQueriesStore` contract (`datatug-core/pkg/datatug/query_revision.go`):
  create only when absent, update or delete only from the revision read, and a typed
  revision-conflict error otherwise;
- typed validation errors, raised before any write, for a query location that escapes the
  project or passes through a symlink, and for content that carries credentials or result rows;
- after an update interrupted mid-write, a freshly opened store reads exactly the old or the new
  complete query, never a mix;
- no Git commit is made.

The query body stays in its `.query.dtql` sidecar (Open Question 4, resolved), carried as the query record's companion body file. Neither carrier provides it on `dalgo2ingitdb`
v0.6.1 today: records are written in place with `os.WriteFile` (`record_io.go:175`, `:223`), not
by temp file plus rename. So the interrupted-update clause waits on crash-atomic record writes
upstream (`ingitdb/dalgo2ingitdb#15`), and a sidecar carrier would additionally need a durable
cross-file journal. Because the filestore is today's only implementation of that guarantee, the
cut-over MUST NOT delete it before the guarantee holds on the DALgo store.

#### REQ: driver-prerequisites-are-recorded

The local path depends on three `dalgo2ingitdb` capabilities that v0.6.1 lacks. They MUST be
recorded, not assumed. The clauses that need them are enumerated under Blocked Clauses and MUST
NOT be recorded as passing before they land:

- **Crash-atomic record writes** (`ingitdb/dalgo2ingitdb#15`), needed by the interrupted-update
  clause of AC:query-write-guards-survive-the-cut-over.
- **Creating a collection with its full `record_file` and subcollections through DALgo.** This
  needs two changes: `dal-go/dalgo#163`, a schema API that can declare subcollections and pass
  driver options, and `ingitdb/dalgo2ingitdb#16`, whose `CreateCollection` today always writes the
  default `{key}.yaml` record file (`schema_modifier.go:405-421`). Founder decision, 2026-09-17:
  *"I'd prefer to fix issues at source upstream instead of doing temporary solutions"*. So there is
  no bootstrap install of the schema, and every clause that needs a created store waits on both.
- **A companion body file** (`ingitdb/ingitdb-go#24`), so a query record's DTQL text lives in its
  `{key}.query.dtql` sidecar under DALgo. Writing the record and its body crash-safely as one
  change is the widened `ingitdb/dalgo2ingitdb#15`.
- **A `README.md` per collection and per record, written by the driver**
  (`ingitdb/dalgo2ingitdb#17`), needed by AC:readmes-are-present-and-do-not-affect-reads,
  because DataTug writes go through DALgo and not beside it.


### Local path and its consumer pin

The local path's prerequisite is a **consumer pin**, not a driver bump.
`record_file.records_dir` is the inGitDB property that omits the `$records` segment, and it
exists from `ingitdb-go/ingitdb` v0.6.0. `dalgo2ingitdb` v0.6.1's own `go.mod` declares only
a *minimum* of `ingitdb-go/ingitdb` v0.5.2, so the resolved version is decided by the
consuming module: a consumer that does not require `ingitdb-go/ingitdb >= v0.6.0` directly
resolves v0.5.2, whose `RecordFileDef` has no `RecordsDir` field, and the strict
`KnownFields(true)` reader (`ingitdb-go/ingitdb@v0.5.2/validator/def_validator.go:35`; the
same reader is at the identical path and line in v0.6.0 and v0.6.1, so the resolved version is the only
difference) then
rejects a definition carrying `record_file.records_dir` as an unknown field. A consumer that
resolves v0.6.0 or later accepts the key, and the canonical `$records`-free paths are producible; the
record path is `db.reader.ReadDefinition` (`dalgo2ingitdb@v0.6.1/database.go:321`) →
`colDef.RecordFile.RecordsBasePath()` (`query.go:587`, `record_io.go:118`), and the non-strict
`yaml.Unmarshal` in `schema_reader.go:105` is `DescribeCollection` only. `datatug-cli` already
satisfies the pin (`datatug-cli/go.mod:39` directly requires
`github.com/ingitdb/ingitdb-go/ingitdb` at v0.6.0 or later), so the canonical layout is producible on the
local path today. The pin belongs to whichever module **constructs the local driver**, not to
`datatug-core`: the CLI already satisfies it, the backend is the other edge that constructs the
driver and MUST satisfy it too, and `datatug-core` MUST NOT declare any DALgo driver
(REQ:project-store-is-a-dalgo-database), so a requirement on core would be vacuous — `go mod
tidy` would prune an unused `require`. This is recorded rather than assumed, exactly as the
remote nesting prerequisite is.

#### REQ: local-records-dir-prerequisite-is-recorded

The local path depends on a consumer-side prerequisite that MUST be recorded, not assumed:
a consuming module MUST resolve `ingitdb-go/ingitdb >= v0.6.0`, so that a collection
definition carrying `record_file.records_dir: '.'` is accepted by the strict
`KnownFields(true)` reader and `RecordFileDef.RecordsBasePath()` returns no `$records` segment
(`RecordsDir` at `record_file_def.go:29`, `:131-137`). The pin is the whole requirement:
`dalgo2ingitdb` v0.6.1 already calls `RecordsBasePath()` on the definition it reads and its
`go.mod` carries only a v0.5.2 floor, so it needs nothing further — the minimum is raised by
the consumer, not the driver, and "propagate `record_file.records_dir` through its schema
reader" is **not** the mechanism: `schema_reader.go:105` is not on the record path and maps
columns for `DescribeCollection` only. The consequence of missing the pin is stated, not
implied: a consumer resolving `ingitdb-go/ingitdb` v0.5.2 has a `RecordFileDef` with no
`RecordsDir` field, so the strict reader rejects a definition carrying
`record_file.records_dir` as an unknown field rather than dropping it. `datatug-cli` already
satisfies the pin by requiring `github.com/ingitdb/ingitdb-go/ingitdb` at v0.6.0 or later directly
(`datatug-cli/go.mod:39`), and the pin is owed by whichever module **constructs the driver** —
the CLI today, the backend at the other edge — not by `datatug-core`, which MUST NOT declare a
DALgo driver at all (REQ:project-store-is-a-dalgo-database). With
the pin satisfied, the canonical `$records`-free paths are producible and the local AC
clauses are exercisable; no local AC clause is blocked by this requirement.

### Remote path and its upstream prerequisite

The remote driver capability to nest already exists, but the `openvaultdb-go` manifest
cannot express it: `buildInGitDBDefinition` (`pkg/mount/ingitdb_github.go:75-85`) builds
only flat `def.Collections[name]` with `DirPath: name` and no `SubCollections`, because
`schema.Collection` (`pkg/schema/schema.go`) carries only `Fields`. It also hardcodes
`format: yaml` and `name: '{key}.yaml'`. What upstream must therefore deliver is a manifest
schema able to declare **subcollection names** and the **full per-collection `record_file`
configuration** — the whole `ingitdb.RecordFileDef`, carrying at minimum its `format`, its
`name` template, its `records_dir` (`record_file_def.go:25-29`) and its `type` (the record
type, `record_file_def.go:16-20`) — surfaced through the OpenVaultDB API, so
that `buildInGitDBDefinition` can populate subcollections instead of only flat
`def.Collections[name]`. All four keys are load-bearing, not a convenience: `format` and
`name` alone cannot produce the asserted layout, because a `{key}`-bearing name with no
`records_dir` makes `RecordFileDef.RecordsBasePath()` return `"$records"`
(`record_file_def.go:131-142`), so a query lands at
`projects/p1/queries/$records/q1/q1.query.json` instead of the asserted
`projects/p1/queries/q1/q1.query.json`; only `records_dir: '.'` yields the asserted path, and
the Feature's own driver facts above already state that mechanism. `type` is required too,
because the user-data dataset needs `ListOfRecords` or `MapOfRecords` remotely and a record
type cannot be expressed by `format` or `name` either. This is an upstream prerequisite, not
something this Feature can work around, so it is recorded rather than assumed; if it does not
land, the remote path degrades to flat single-directory collections (and on the GitHub engine cannot
write a project item at all; see below) and every nested-layout AC
clause stays blocked (see Blocked Clauses).

#### REQ: remote-nesting-prerequisite-is-recorded

The remote (OpenVaultDB) path depends on an upstream prerequisite that MUST be recorded,
not assumed: `openvaultdb-go`'s `buildInGitDBDefinition` currently builds only flat
collections and hardcodes YAML and a `{key}.yaml` name, because its manifest schema has
no subcollections and no `record_file` config. The prerequisite is concrete: `openvaultdb-go`
MUST provide a manifest schema able to declare **subcollection names** and the **full
per-collection `record_file` configuration** — the whole `ingitdb.RecordFileDef`, carrying at
minimum `format`, the `name` template, `records_dir` (`record_file_def.go:25-29`) and `type`
(the record type, `record_file_def.go:16-20`) — surfaced
through the OpenVaultDB API, so that `buildInGitDBDefinition` can populate subcollections
instead of only flat `def.Collections[name]`. A manifest that exposes only `format` and
`name` is **not** sufficient and does not unblock the remote canonical layout: without
`records_dir`, `RecordFileDef.RecordsBasePath()` returns `"$records"` for a `{key}`-bearing
name (`record_file_def.go:131-142`) and the record lands at
`projects/p1/queries/$records/q1/q1.query.json`; without `type`, the user-data dataset's
`ListOfRecords`/`MapOfRecords` shape cannot be declared. Until the whole configuration lands,
the remote path cannot produce the canonical nested layout; the consequence is stated, not
implied — the remote
path degrades to flat single-directory collections (on the GitHub engine it cannot write a project
item at all), so every AC clause that asserts the
nested layout, enumerated under Blocked Clauses, stays blocked. Every AC clause that cannot
be exercised until it lands MUST be stated as such and MUST NOT be scaffolded as passing.

The prerequisite is tracked upstream as `openvaultdb/openvaultdb-go#26`, and its reach depends on
the host's storage engine:

- **GitHub engine.** The manifest-built definition carries no `SubCollections`, and
  `dalgo2ingitdb4github/scoped_collection.go:86-99` (v0.2.8) rejects any subcollection key missing
  from the definition. No project item, all of which live under `projects/<project-id>/`, can be
  written there at all, flat or nested. Every remote clause on the GitHub engine is blocked.
- **Local `ingitdb` engine.** The host opens `dalgo2ingitdb` over its data directory and reads the
  on-disk definitions (`openvaultdb-go/pkg/mount/mount.go` `openInGitDB`). With DataTug's schema installed there and `schema_mode: schemaless` declared, it is expected to serve the canonical nested layout today. Its writes are committed to Git, and a user-supplied host has no API path to install the schema until #26 lands. That expectation
  MUST be proven by a test before any clause relies on it, and until proven it is not recorded as
  passing.

## Acceptance Criteria

### AC: core-depends-on-dalgo-interfaces-only (verifies REQ:project-store-is-a-dalgo-database)

**Given** `datatug-core/go.mod`
**When** its module graph and the packages above the DALgo boundary are inspected
**Then** it depends on the `dal-go/dalgo` interfaces module, declares no DALgo driver module (`dalgo2ingitdb`, `dalgo2openvaultdb`, `dalgo2ingitdb4github`), and no package above the DALgo boundary branches on a concrete driver type.

### AC: driver-is-selected-at-the-edge (verifies REQ:project-store-is-a-dalgo-database)

**Given** the same DALgo-facing project-store code path and two different local configurations of the local `dalgo2ingitdb` driver (for example two different project roots, or two different `dalgo2ingitdb` database definitions)
**When** a project is opened under each local configuration
**Then** both configurations are served by that same code path, and no code above the DALgo boundary branches on a concrete driver type.

### AC: local-and-remote-paths-are-one-code-path (verifies REQ:project-store-is-a-dalgo-database)

**Given** the same DALgo-facing project-store code path and both the local `dalgo2ingitdb` configuration and the remote `dalgo2openvaultdb` backend
**When** a project is opened through each
**Then** both are served by that same code path with no code above the DALgo boundary branching on a concrete driver type, and the remote leg serves that code path over the canonical nested layout, landing a saved record at `projects/p1/queries/q1/q1.query.json`.

### AC: collection-declares-the-canonical-record-file (verifies REQ:canonical-project-layout)

**Given** the authored `definition.yaml` for a project-item collection whose inherited suffix is `query`
**When** that file's text is read
**Then** it declares a `record_file` mapping carrying `name: '{key}/{key}.query.json'`, `format: json`, `type: map[string]any` and a nested `records_dir: '.'`, with `record_file.records_dir` written **inside** `record_file` and not as a sibling key.

### AC: local-write-is-one-directory-per-record (verifies REQ:canonical-project-layout)

**Given** a project `p1` on the local `dalgo2ingitdb` v0.6.1 driver whose collection definition carries the canonical `record_file` (`records_dir: '.'`)
**When** a query with id `q1` is saved
**Then** the record gets its own directory holding `q1.query.json` — the parent directory is the `queries` collection directory, so the file is `projects/p1/queries/q1/q1.query.json` and not a flat `projects/p1/queries/q1.query.json` written directly into `queries`.

### AC: local-write-omits-the-records-segment (verifies REQ:canonical-project-layout)

**Given** a local `dalgo2ingitdb` driver whose `ingitdb-go/ingitdb` dependency resolves `>= v0.6.0` so that `record_file.records_dir` is honoured
**When** a query with id `q1` in project `p1` is saved
**Then** it lands at `projects/p1/queries/q1/q1.query.json` with no `$records` segment in the path.

### AC: a-store-holds-several-projects (verifies REQ:canonical-project-layout)

**Given** a store holding project `p1` and a second project `p2`
**When** `p2` is written to the same store and addressed through DALgo
**Then** `p2` is addressed under its own `projects/p2` subtree, distinct from `projects/p1`, and its records do not collide with `p1`'s.

### AC: remote-write-produces-the-canonical-nested-layout (verifies REQ:canonical-project-layout)

**Given** a project `p1` on the remote `dalgo2openvaultdb` driver
**When** a query with id `q1` is saved
**Then** it lands at `projects/p1/queries/q1/q1.query.json` in the same canonical nested layout the local driver produces.

### AC: credential-address-is-part-of-the-layout (verifies REQ:canonical-project-layout)

**Given** a credential for server `s1` in project `p1`
**When** it is addressed through DALgo
**Then** the key path is `projects/p1/credentials/s1` and the record file is `projects/p1/credentials/s1/s1.credential.json` carrying the `credential` suffix.

### AC: detached-file-is-identified-by-its-suffix (verifies REQ:type-suffix-is-self-identifying)

**Given** one project-item file of each type, copied out of its directory into a scratch directory
**When** each detached filename is inspected without its parent directory
**Then** the type suffix alone identifies the item type.

### AC: no-bare-key-json-and-no-suffix-ends-another (verifies REQ:type-suffix-is-self-identifying)

**Given** the declared collection definitions for all project-item types
**When** their `record_file.name` templates and type suffixes are checked
**Then** no template is a bare `{key}.json`, every template ends with a type-specific dotted `.<suffix>.json` segment, and no template ends with another collection's dotted suffix.

### AC: nested-item-is-addressed-by-a-key-path (verifies REQ:hierarchy-is-a-key-path)

**Given** a query `q1` inside project `p1`
**When** it is addressed through DALgo
**Then** the key path is `projects/p1/queries/q1` in the shape `collection/id/subcollection/id`, and no substituted record key contains a path separator.

### AC: db-server-is-addressed-deeply (verifies REQ:hierarchy-is-a-key-path)

**Given** a db server `s1` belonging to db driver `d1` of project `p1`
**When** it is addressed through DALgo
**Then** the key path is `projects/p1/dbdrivers/d1/dbservers/s1` and its record file is `projects/p1/dbdrivers/d1/dbservers/s1/s1.dbserver.json`.

### AC: readmes-are-present-and-do-not-affect-reads (verifies REQ:readme-per-collection-and-record)

**Given** a project-item collection whose collection directory and whose record directory carry no `README.md`
**When** a record is written to that collection and the collection's records are then listed and read
**Then** a `README.md` is present at the collection directory and inside the record directory, neither `README.md` is returned as a record, and every record still reads.

### AC: user-data-stays-compact (verifies REQ:user-data-uses-a-compact-ingr-file)

**Given** a dataset `support-notes` in project `p1`'s user data
**When** it is persisted
**Then** it is one multi-record file at `projects/p1/data/ingitdb/support-notes/support-notes.ingr` under the project's user-data database root with `format: ingr`, a `type` of `ListOfRecords` (`'[]map[string]any'`) or `MapOfRecords` (`'map[$record_id]map[$field_name]any'`), and a dataset-collection `record_file.name` of `support-notes.ingr` carrying no `{key}` placeholder, and it is not a directory per record.

### AC: data-is-reserved-for-user-data (verifies REQ:user-data-is-not-a-project-item)

**Given** a project `p1` whose project items and user data have both been written
**When** the `projects/p1/data/` subtree, its `projects/p1/data/ingitdb/` database root, and the collection directories are listed
**Then** every project item lies outside `projects/p1/data/`, and `projects/p1/data/` holds only `p1`'s user data — one inGitDB database at `projects/p1/data/ingitdb/`, with each dataset file at `projects/p1/data/ingitdb/<dataset>/<dataset>.ingr`.

### AC: recordset-definitions-collapse-to-one-location (verifies REQ:recordset-definitions-are-a-project-item)

**Given** a recordset definition `support-notes` in project `p1`
**When** it is saved through the DALgo-backed store and then loaded
**Then** it is written to and read from `projects/p1/recordsets/support-notes/support-notes.recordset.json`, and the save writes neither `projects/p1/data/recordsets/support-notes.recordset.json` (the loader's project-relative current path) nor the flat `projects/p1/recordsets/support-notes.recordset.json`.

### AC: suffixes-are-inherited-from-the-jsonfilename-set (verifies REQ:existing-filenames-stay-recognisable)

**Given** the `JsonFileName` suffix constants in `datatug-core/pkg/storage/file_names.go`
**When** each project-item type's canonical suffix is assigned
**Then** every type with an existing constant keeps that exact value (`query`, `entity`, `env`, `server`, `db`, `dbmodel`, `board`, `recordset`, `dbserver`), and only `project`, `folder`, `dbdriver` and `credential` — which have no constant — receive a new suffix.

### AC: foreign-files-survive-a-project-write (verifies REQ:store-preserves-foreign-files)

**Given** a store that already contains unrelated files this Feature does not manage
**When** a project is written into that store
**Then** every pre-existing unrelated file is unchanged byte for byte.

### AC: projects-in-one-store-do-not-interfere (verifies REQ:store-preserves-foreign-files)

**Given** two projects `p1` and `p2` written to the same store
**When** each project is written
**Then** the two projects' records are disjoint, and the write of one alters none of the other project's files.

### AC: changes-go-through-dalgo-transactions (verifies REQ:bespoke-transaction-machinery-is-retired)

**Given** a change touching more than one project-item record
**When** it is applied through the store
**Then** it runs inside a DALgo transaction, no filestore transaction or lock code (`query_txn.go`, `query_txn_slots.go`, `query_lock.go`) is on the write path, and after the mutation a fresh DALgo read in a separate transaction observes the change while no bespoke journal or temporary file remains.

### AC: reserved-transaction-namespace-is-gone (verifies REQ:bespoke-transaction-machinery-is-retired)

**Given** a project directory written by the DALgo-backed store
**When** the repository root is listed and searched repository-wide
**Then** it contains no reserved transaction directory (`.dt-query-txn/`) and no lock file, no on-disk transaction bookkeeping remains part of the project layout, and a search for the reserved namespace name `.dt-query-txn` finds nothing.

### AC: filestore-has-no-callers (verifies REQ:filestore-is-cut-over)

**Given** every Go module in the inventory (`datatug-core`, `datatug-cli`, `datatug/backend`, `datatug-project-validator`, `datatug-demo-projects`) at the cut-over, each requiring the latest `datatug-core`
**When** their packages are searched for `datatug-core/pkg/storage/filestore` and `datatug-core` is listed
**Then** the package does not exist in `datatug-core`, no module imports it, and every consumer in the inventory reads and writes projects through the DALgo-backed store.

### AC: no-legacy-layout-reader-remains (verifies REQ:filestore-is-cut-over)

**Given** `datatug-apps` and `ai-plugin` on `main` after the cut-over
**When** their sources and docs are searched for the legacy names (`datatug-project.json` at a project root, `environment-summary`, `.datatug-`, `data/recordsets`, `$records`)
**Then** no reader requests a legacy-layout path, no doc describes one, and the `datatug-apps` GitHub project reader reads a canonical-layout project from a golden fixture.

### AC: legacy-layout-is-refused (verifies REQ:filestore-is-cut-over)

**Given** a project directory in a legacy layout, one case each for a flat `queries/q1.query.json`, a `data/recordsets/` definition, an `environment-summary.json` and a `$records/` directory
**When** the project is opened through the DALgo-backed store
**Then** each open fails with a typed error naming the unsupported path, and none returns a partially loaded project.

### AC: owned-demo-repository-uses-the-canonical-layout (verifies REQ:filestore-is-cut-over)

**Given** `datatug/datatug-demo-projects` on `main` after the cut-over
**When** every project in it is opened through the DALgo-backed store and its files are listed
**Then** every project loads, every project item is at its canonical `<key>/<key>.<type-suffix>.json` path, and no legacy-layout file remains.

### AC: query-write-guards-survive-the-cut-over (verifies REQ:query-write-guards-survive-the-cut-over)

**Given** a project on the DALgo-backed store whose query `customers/customer-invoices` is at revision `R`
**When** the store receives, in turn, a write whose location escapes the project or passes through a symlink, a write whose content carries credentials or result rows, a create for an id that already exists, an update conditioned on a revision other than `R`, and an update interrupted mid-write
**Then** each unsafe or stale write is refused with a typed validation or revision-conflict error and changes no file, a freshly opened store after the interruption reads exactly the old or the new complete query and never a mix, and nothing is committed to Git.

### AC: remote-nesting-prerequisite-is-stated (verifies REQ:remote-nesting-prerequisite-is-recorded)

**Given** this Feature's recorded prerequisites
**When** the remote path is read for what it depends on
**Then** it states that `openvaultdb-go`'s `buildInGitDBDefinition` builds only flat `def.Collections[name]` with `DirPath: name` and no `SubCollections` because `schema.Collection` carries only `Fields`, that it hardcodes `format: yaml` and `name: '{key}.yaml'`, that the prerequisite requires the manifest to declare **subcollection names** and the **full per-collection `record_file` configuration** — `format`, the `name` template, `records_dir` and `type`, not `format` and `name` alone — surfaced through the OpenVaultDB API, that without `records_dir` a `{key}`-bearing name makes `RecordFileDef.RecordsBasePath()` return `"$records"` (`record_file_def.go:131-142`) so a record lands at `projects/p1/queries/$records/q1/q1.query.json` rather than the canonical `projects/p1/queries/q1/q1.query.json`, and that nesting is resolved by `resolveScopedCollection` in `dalgo2ingitdb4github/scoped_collection.go` on the remote path and by the same-shaped `resolveScopedCollection` in `dalgo2ingitdb/scoped_collection.go` on the local path.

### AC: local-records-dir-prerequisite-is-stated (verifies REQ:local-records-dir-prerequisite-is-recorded)

**Given** this Feature's recorded prerequisites
**When** the local path is read for what it depends on
**Then** it states that `record_file.records_dir: '.'` is the mechanism that omits the `$records` segment and is a property of `ingitdb-go/ingitdb >= v0.6.0` (`RecordFileDef.RecordsDir` declared at `record_file_def.go:29`, validated at `:61`, returning no `$records` segment at `:131-137`), that `dalgo2ingitdb` v0.6.1's `go.mod` declares only a *minimum* of `ingitdb-go/ingitdb` v0.5.2 so the requirement falls on the consuming module to resolve `ingitdb-go/ingitdb >= v0.6.0`, that a consumer resolving v0.5.2 has a `RecordFileDef` with no `RecordsDir` field and the strict `KnownFields(true)` reader then rejects a definition carrying `record_file.records_dir` as an unknown field on the `db.reader.ReadDefinition` → `RecordFileDef.RecordsBasePath()` record path, that `schema_reader.go:105` is `DescribeCollection` only and not on that path, and that `datatug-cli` already satisfies the pin by directly requiring `github.com/ingitdb/ingitdb-go/ingitdb` at v0.6.0 or later (`datatug-cli/go.mod:39`) while the pin is owed by whichever module constructs the driver — the CLI and the backend — and not by `datatug-core`, which MUST NOT declare a DALgo driver.

### AC: blocked-clauses-are-not-recorded-as-passing (verifies REQ:remote-nesting-prerequisite-is-recorded)

**Given** the AC clauses that wait on an upstream prerequisite (the remote nesting prerequisite and the recorded `dalgo2ingitdb` driver prerequisites)
**When** the Feature's status is recorded
**Then** every such clause is enumerated in the `## Blocked Clauses` list below as un-exercisable until the upstream prerequisite lands, and none of them is recorded as passing.

## Blocked Clauses

Upstream prerequisites that have not landed block the following AC clauses. This list is the
carrier AC:blocked-clauses-are-not-recorded-as-passing points at. Until each owner lands, none of
these clauses may be scaffolded or recorded as passing. The rest of each AC is exercisable and
MUST NOT be blocked along with the named clause. The consumer pin in
REQ:local-records-dir-prerequisite-is-recorded blocks nothing.

1. **AC:remote-write-produces-the-canonical-nested-layout** (GitHub-engine leg). Owner:
   REQ:remote-nesting-prerequisite-is-recorded, `openvaultdb/openvaultdb-go#26`. The local-engine
   leg is exercisable once a test proves the expectation stated in that requirement.
2. **AC:local-and-remote-paths-are-one-code-path** (GitHub-engine leg of the remote clause). Owner:
   the same. The rest of the AC is exercisable against an OpenVaultDB host on its local `ingitdb`
   engine: one DALgo-facing code path serves both, with no driver branch above the boundary.
3. **AC:query-write-guards-survive-the-cut-over** (the interrupted-update clause). Owner:
   REQ:driver-prerequisites-are-recorded, `ingitdb/dalgo2ingitdb#15`. A sidecar carrier would also
   need a durable cross-file journal (Open Question 4).
5. **Every AC clause that needs a DALgo-created store**: AC:collection-declares-the-canonical-record-file
   and, through it, every local layout, store, user-data, recordset, transaction, query-guard,
   legacy-refusal and cut-over AC. Owner: REQ:driver-prerequisites-are-recorded, tracked as
   `dal-go/dalgo#163` and `ingitdb/dalgo2ingitdb#16`. AC:core-depends-on-dalgo-interfaces-only, the
   two documentary prerequisite ACs and the carrier AC are not blocked.
6. **AC:query-write-guards-survive-the-cut-over** (the `.query.dtql` pair clauses). Owner:
   REQ:driver-prerequisites-are-recorded, tracked as `ingitdb/ingitdb-go#24` (companion body file).
4. **AC:readmes-are-present-and-do-not-affect-reads** (whole AC). Owner:
   REQ:driver-prerequisites-are-recorded, `ingitdb/dalgo2ingitdb#17`.



No `_tests/` stubs are scaffolded in this authoring pass, and that is a deliberate,
recorded decision rather than an omission.

Most of this Feature's ACs assert facts about a store that does not exist yet: the file
shapes (AC:collection-declares-the-canonical-record-file,
AC:local-write-is-one-directory-per-record,
AC:detached-file-is-identified-by-its-suffix,
AC:no-bare-key-json-and-no-suffix-ends-another,
AC:nested-item-is-addressed-by-a-key-path,
AC:suffixes-are-inherited-from-the-jsonfilename-set,
AC:a-store-holds-several-projects,
AC:credential-address-is-part-of-the-layout, AC:db-server-is-addressed-deeply,
AC:readmes-are-present-and-do-not-affect-reads,
AC:user-data-stays-compact, AC:data-is-reserved-for-user-data,
AC:recordset-definitions-collapse-to-one-location, AC:reserved-transaction-namespace-is-gone),
the module graph (AC:core-depends-on-dalgo-interfaces-only) and store-level round trips
(AC:driver-is-selected-at-the-edge, AC:changes-go-through-dalgo-transactions,
AC:foreign-files-survive-a-project-write, AC:projects-in-one-store-do-not-interfere) and the
cut-over (AC:filestore-has-no-callers, AC:legacy-layout-is-refused,
AC:owned-demo-repository-uses-the-canonical-layout,
AC:query-write-guards-survive-the-cut-over). Their
testable surfaces are real — filesystem layout assertions, `go list -m` over
`datatug-core/go.mod`, and DALgo-backed round trips — so stubs SHOULD be scaffolded when
implementation begins, pinned to the implementation's own package paths. Scaffolding them
now would fix test paths for packages and collection definitions that have not been
created, which is a guess, not a rehearsal.

Three ACs are documentary and are skipped by the Rehearse heuristic's not-testable branch,
because they have no runtime observer:
AC:remote-nesting-prerequisite-is-stated,
AC:local-records-dir-prerequisite-is-stated and
AC:blocked-clauses-are-not-recorded-as-passing assert what this document records,
not what code does.

The clauses enumerated under Blocked Clauses MUST NOT be scaffolded or recorded as passing until
their owners land:

- the GitHub-engine legs of AC:remote-write-produces-the-canonical-nested-layout and
  AC:local-and-remote-paths-are-one-code-path (`openvaultdb/openvaultdb-go#26`, which needs the
  **full** per-collection `record_file`, not a name template);
- the interrupted-update clause of AC:query-write-guards-survive-the-cut-over
  (`ingitdb/dalgo2ingitdb#15`);
- the whole of AC:readmes-are-present-and-do-not-affect-reads (`ingitdb/dalgo2ingitdb#17`).

The local-engine OpenVaultDB legs become rehearsable once a test proves the host serves the
nested layout. The consumer pin in REQ:local-records-dir-prerequisite-is-recorded is satisfied, so
AC:local-write-omits-the-records-segment and AC:collection-declares-the-canonical-record-file are
exercisable and MUST be scaffolded when implementation begins.

## Not Doing / Out of Scope

- **The parent's four backends are delivered here only in part; the residual is owned
  follow-on work of this tree, landed incrementally.** The parent requirement
  `git-native-projects` REQ:project-store-pluggable-backends enumerates four backends: local
  directory, GitHub repository, a remote OpenVaultDB vault, and DataTug Cloud. This Feature
  commits to the
  **driver-agnostic code path** — no driver branch above the DALgo boundary — plus the
  **local** backend concretely, and owns the canonical layout and layout policy every
  backend rides on. The parent requirement is therefore **partially** delivered, with the
  residual named per backend:
  - **Local directory — delivered here.** This Feature owns the canonical layout, the
    suffix grammar, the DALgo boundary and the local-driver ACs.
  - **Remote OpenVaultDB vault — owned by `openvaultdb-project-store`.** The store binding
    is owned there (REQ:project-store-binding), not by this Feature. This Feature owns the
    canonical layout the binding must realise and the recorded upstream prerequisite
    (`openvaultdb-go` subcollections and the full per-collection `record_file` configuration —
    `format`, `name`, `records_dir` and `type`) the remote nested path
    depends on.
  - **GitHub-native read parity — owned follow-on work of this tree, landed incrementally,
    via the GitHub-backed driver.** The GitHub-backed driver (`dalgo2ingitdb4github`) is the
    framing through which this tree reaches the GitHub native paths as follow-on work, not
    orphaned and not deferred to an unnamed future; the outcome it must reach is already named
    by `git-native-projects` REQ:github-hosted-projects-read-write. This Feature does **not**
    discharge that requirement: its `AC:github-project-reads-back-what-it-wrote` is scoped to
    the `Store`/`GhStorage` surface, `GhStorage`'s `FileExists`/`OpenFile` still
    `panic("implement me")` (`datatug-cli/pkg/dtgithub/github_storage.go:40-47`), and adding
    `GhStorage` read support is a separate task owned by `git-native-projects`
    (`openvaultdb-project-store` records the same ownership).
  - **DataTug Cloud — owned follow-on work of this tree, landed incrementally, via
    `dalgo2firestore`.** It is a named residual of the parent requirement owned by this
    tree's DALgo driver work, not orphaned and not deferred to an unnamed future.
- **The OpenVaultDB remote binding itself** — owned by `openvaultdb-project-store`. This
  Feature records the upstream nesting prerequisite the remote path depends on; it does
  not build the binding, the registry, or the browser read path.
- **Secrets vaults, credentials and connection strings** — owned by `secrets-vault`.
  Credentials are never written into a project's own asset files and are never stored inside
  another project-item type's collection directory — `credentials` is its own collection;
  they are held in a **vault store** whenever a `secrets` binding exists, and with zero
  `secrets` bindings they may live in the project's own items store subject to
  `secrets-vault`'s colocation policy. Their canonical address inside the layout is declared
  by REQ:canonical-project-layout, while their semantics belong to `secrets-vault`.
- **User-data semantics** — what a dataset means, how it is queried and how it is
  modelled. `projects/<project-id>/data/`, whose inGitDB database root is
  `projects/<project-id>/data/ingitdb/`, is *reserved* here; the data model is not owned
  here.
- **The deferred `policies/`, `fixtures/`, `datasources/` question.** Whether these belong
  inside the project-item structure is deliberately deferred. Findings recorded so they
  are not lost: `policies/` is live via `datatug-cli` `pkg/accesspolicies` and is not in
  `ProjectStore`; `fixtures/` is live, consumed by the demo repo's own tests, and is not in
  `ProjectStore`; `datasources/` is orphaned — no reader in `datatug-core/pkg`,
  `datatug-cli/pkg` or `datatug-cli/apps`, and a search of `datatug-apps` timed out and is
  unverified. If these are kept, `datasources/` should carry a README noting it may be a
  leftover to be deleted or implemented, and that note removed once resolved.
- **Any backward compatibility with the legacy layouts.** Founder decision, 2026-09-17: a hard
  cut-over. No read-only shim for the dual layouts (entity, board, db model, env db catalog,
  and `environment-summary.json` vs `<id>.env.json`), and no conversion of legacy `*.sql.json`
  queries. REQ:filestore-is-cut-over refuses them all.
- **Migrating repositories DataTug does not own.** Only DataTug-owned repositories are
  rewritten. A user's legacy repository is refused with a typed error, not converted.
- **The repo-root `.datatug.yaml`** is an Open Question, not settled by this Feature. How the
  query body is carried is Open Question 4. The guarantees it must keep are not open
  (REQ:query-write-guards-survive-the-cut-over).
- **A home or retirement for the `columns`, `objects` and `refs` suffixes** — inherited
  suffixes with no collection in the current `ProjectStore` interface; an Open Question.

## Open Questions

1. ~~Read-only compatibility shim for the legacy dual layouts, or hard cut-over?~~ Resolved
   2026-09-17: hard cut-over, no backward compatibility (REQ:filestore-is-cut-over).
2. ~~Legacy `*.sql.json` queries — convert or keep ignoring?~~ Resolved 2026-09-17: neither.
   They are unsupported and refused under the hard cut-over.
3. `columns`, `objects` and `refs` suffixes have no collection in `ProjectStore` — do they
   get a home, or are they retired?
4. ~~Query body carrier: sidecar or in the record?~~ Resolved 2026-09-17 (founder): **keep the
   `.query.dtql` sidecar**. Under DALgo it becomes a companion body file of the query record
   (`ingitdb/ingitdb-go#24`), written crash-safely with the record (`ingitdb/dalgo2ingitdb#15`,
   widened to multi-file). `executable-knowledge-library` AC:query-pair-storage-guards-writes
   therefore stands as written.
5. ~~Bootstrap schema install until #16?~~ Resolved 2026-09-17 (founder): no temporary solution.
   The schema is created through DALgo once `dal-go/dalgo#163` and `ingitdb/dalgo2ingitdb#16` land.
6. Is `.datatug.yaml` still needed at the repo root, given projects are discoverable by
   listing the `projects` collection and nothing in `datatug-core` reads it during
   `LoadProject`?
7. ~~Who owns the `openvaultdb-go` subcollection/full `record_file` configuration
   prerequisite?~~ Owned upstream as `openvaultdb/openvaultdb-go#26`, filed 2026-09-17. When it
   lands is still open. Its absence blocks only the remote *nested* path — the remote path stays
   flat, as stated in REQ:remote-nesting-prerequisite-is-recorded and Blocked Clauses.
8. ~~Do the two migration defects get fixed by this move?~~ Resolved 2026-09-17: yes, by the
   hard cut-over (REQ:filestore-is-cut-over).

---
*This document follows the https://specscore.md/feature-specification*
