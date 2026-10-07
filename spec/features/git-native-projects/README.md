---
format: https://specscore.md/feature-specification
status: Implementing
---

# Feature: Git-native projects

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=request-change) |
**Status:** Implementing
**Source Ideas:** —

## Summary

P9 is Git-native projects: every knowledge asset is a readable, versioned, diffable, portable file — no secrets in project files — and the project store itself, wherever it lives, is explorable inside DataTug.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

It is the pillar that makes P7's compounding library real rather than aspirational: a captured query is worthless as an accumulating asset if it cannot be reviewed, diffed and moved between local, hosted and Git-backed stores. Initial audit snapshot (before later 9 September landings): the local filestore is the only fully working `Store` backend; `--git=stage` only covers `entity add`; secrets can still land in committed query files because `Project.Validate()` never validates `Queries` at all; GitHub- hosted projects can write but not reliably read back (`GhStorage.FileExists`/`OpenFile` panic). Two founder additions from 2026-09-09 extend this pillar: a remote OpenVaultDB project store, and making the project store itself explorable as a source inside DataTug. Phase 0–1 add the secrets guard; Phase 2–3 build the pluggable store backends and self-exploration; Phase 4 adds diff/review affordances. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`) link to it.

The immediate browser journey is a signed-in user choosing a GitHub repository
they own, creating or selecting a working branch, saving several project queries
through DataTug Cloud, reviewing the diff, and merging to a chosen target. They
can reload and run a saved query against a hosted DemoDB source in the browser
without installing the CLI. [Public Data Fabric](../public-data-fabric/README.md)
owns browser OVDB federation and row execution; this Feature owns where the
query definition is read and written.
DataTug Cloud's Firestore project store accepts the same project mutation
contract but has no Git branch capability. The CLI-served API accepts that
contract for a local clone, with local pull/push operations on its current
branch. These are requirements, not a claim that all paths work today.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> A knowledge asset that cannot be reviewed before it changes team behaviour is a liability, not an
> asset — this is why DataTug commits to Git-native storage over a database-backed catalog
> (hypothesis H6, "supported as a hypothesis by competitor behaviour... Bruno, DataGrip, dbt, Grafana-as-code and
> Lightdash all moved toward files"). Today the promise is only partly kept: the filestore layout is
> real and the demo project lives on GitHub as proof (`datatug-core/pkg/storage/filestore`), but
> `--git=commit` is explicitly rejected (exit 2) and only `entity add` stages its own changes
> (`git_integration.go`). Secrets modeling is a live gap: `QueryDefTarget.Credentials` can be
> committed to a git-tracked query file and `Project.Validate()` never checks `Queries` at all
> (current-state audit, datatug-core findings). GitHub-hosted projects are "write-only" — the audit
> found `GhStorage.FileExists`/`OpenFile` both `panic("implement me")`
> (`pkg/dtgithub/github_storage.go:40-48`), so a write path exists that cannot read its own output
> back through the same interface.

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| Readable file layout for every asset | usable | — | strategic |
| Mutating commands stage their own files (`--git=stage`) | usable (entity add only) | 2 | enabler |
| Secrets never in project files (validation + guidance) | planned | 0–1 | hygiene |
| GitHub-hosted projects: read and write | partial (write-only) | 2 | strategic |
| Project store abstraction with pluggable backends (local files · GitHub · **remote OpenVaultDB** · DataTug Cloud) | partial (interface exists; filestore only) | 3 | strategic |
| **Project store explorable inside DataTug** (the project's own entities, queries, mappings, policies browsable and queryable as a source) | planned | 2–3 | strategic |
| Diff/review affordances for knowledge changes | planned | 4 | strategic |

#### REQ: readable-file-layout

Every project asset — entities, queries, environments, servers, policies — MUST be stored as a
readable file in a stable, documented layout, so a person or an AI agent can read a project by
reading files, with no database required. This capability is already usable; this REQ pins it as
the non-negotiable baseline every new asset type must follow.

#### REQ: git-stage-mutations

A command that mutates a project MUST stage the files it changed (`--git=stage`) so the change is
visible as a pending Git diff before it is committed. Today only `entity add` does this; every
mutating command added going forward MUST follow the same pattern rather than writing files with no
staging at all. `--git=commit` remains deliberately out of scope — DataTug stages, the user commits.

#### REQ: secrets-never-in-project-files

`datatug validate` MUST reject a project where a `QueryDefTarget` (or any other asset) carries a
credential value. `Project.Validate()` already validates `Queries`: `datatug-core/pkg/datatug/project.go:123-133`
calls `p.Queries.Folders.Validate()` and `p.Queries.Items.Validate()`, and `pkg/datatug/query.go:74-81`
walks every `QueryDef` with `QueryDefTarget.Validate` (`query.go:176-189`) rejecting credential
material. The residual gap is that the resulting error does not name the offending file; closing
that gap is required, and this Feature MUST NOT be read as asserting that `Queries` goes
unvalidated today. Documentation MUST point to a local, non-committed credential store (for example
`~/.datatug`) as the correct place for secrets. This is shared ground with `heterogeneous-sources`
REQ:secrets-outside-project-files; the two Features share one guard, not two.

#### REQ: github-hosted-projects-read-write

A GitHub-hosted project MUST support both write and read through the same `Store` interface other
backends use — today only write works, because `FileExists`/`OpenFile` panic. Read parity is
required before a GitHub-hosted project can be treated as equivalent to a local one.

#### REQ: project-store-pluggable-backends

The project store abstraction MUST support pluggable backends: local directory, GitHub repository,
a remote OpenVaultDB vault, and DataTug Cloud. Founder, 2026-09-09, verbatim: *"It should be
possible to store and use datatug project from remote OpenVaultDB."* The `storage.Store`/
`ProjectStore` interfaces already exist and `dalgo2openvaultdb` already exists as an adapter; this
REQ requires wiring a project store implementation on top of it, with OpenVaultDB's own policies
composing with DataTug's server-enforced access control (see `server-enforced-access`), not
replacing it.

#### REQ: project-store-explorable-inside-datatug

The project store itself MUST be explorable inside DataTug — its own entities, queries, mappings
and policies browsable and queryable as a source, using the same grid, related-data and context
loop the rest of the product provides. Founder, 2026-09-09, verbatim: *"The project store itself
should be available for exploration in DataTug itself."* This falls out naturally once project
records are inGitDB/OpenVaultDB collections: the knowledge graph becomes just another source under
`heterogeneous-sources`, subject to the same navigation loop as Chinook or `support-notes`.

#### REQ: diff-review-affordances

A knowledge change (a new query, an edited mapping, a promoted collection) SHOULD have a
diff/review affordance in DataTug itself, riding on Git rather than reimplementing version control.
This is later-phase UI work; the underlying guarantee (every change is a Git diff) already holds
once REQ:git-stage-mutations covers every mutating command.

#### REQ: one-project-mutation-contract

Creating a project and saving a query MUST use a storage-neutral, versioned
DataTug API contract. Hosted and GitHub stores use the DataTug Cloud API;
local clones use the DataTug CLI-served API. The client selects the API endpoint
and project store reference, not a provider-specific mutation shape. Existing
`create_query` and `update_query` callers remain compatible while an additive
conditional `save_query` operation accepts the existing query definition,
including federation metadata and its body, and returns an opaque query revision.
The same project reference, actor authorization, validation, concurrency,
operation receipt and store-adapter pipeline MUST also govern existing editable
project assets (project metadata, query folders, boards and entities as their
operations are migrated). The first delivery slice covers project creation
and saved queries; it MUST NOT be described as complete asset parity.
Exactly one create-only or update-if-revision-matches condition is required.
The pair of query metadata and type-derived body sidecar MUST be one logical
write; API readers and crash recovery cannot expose a mixed pair after an
interrupted write. This does not require two physical filesystem renames to be
atomic together. Both APIs apply the same
Core query validation and secret screen before the write. A client-supplied
operation ID makes a retry after a committed write or index failure recover
the original result and repair the index, without a duplicate commit. The
receipt binds the verified actor, store, project/repository, branch,
operation kind and canonical request-payload digest. Reusing an operation ID
with different scope or payload returns a typed conflict without a write;
even a matching replay rechecks current actor and provider authorization.
This rule applies to project create, asset save, branch, explicit commit and
sync mutations, not only query saves.

#### REQ: provider-capabilities-and-branch-workflow

Every project API MUST report its supported capabilities. Git stores expose
list, select, create, diff and merge branches; selection is client state and
every mutation names an explicit selected branch and its observed head. A
branch list/create request before project creation may address a prospective
project reference with repository and folder, but still checks actor and
repository authority; no manifest or index is required to select a branch.
The project manifest is then created on that selected branch and observed
head, including when the default branch is protected. A
GitHub save atomically commits its pair on that branch, then conditionally
advances the branch only when both query revision and branch head still match.
It never force-updates a ref. Create checks its source head; merge checks both
reviewed source and target heads and current actor permission. Conflicts return
a typed conflict without modifying either branch. Branch protection can result
in a review proposal or pending checks; the response MUST distinguish that from a
completed merge and include the review URL when one exists. A compare URL is
only an invitation to open a PR, not a claim that the API created one; an
API-created PR requires a separately granted Pull requests permission. The Firestore
project store reports branching unsupported and rejects branch fields rather
than pretending a Firestore revision is a Git branch.

#### REQ: local-current-branch-synchronization

The CLI-served API MUST support pull and push on the local clone's current
branch using its configured upstream. Pull defaults to fast-forward only and
must refuse divergence, a stale observed head, or a dirty worktree that would
be overwritten. Push is a normal non-force push, checking observed local and
remote heads; remote rejection is explicit. Git commands receive validated
branch/remote arguments directly, without shell interpretation, and execute
under the same repository lock as local branch and query mutations. Local
branch switching or merge MUST preserve uncommitted changes: use isolated
branch worktrees or fail cleanly on dirty/current-branch mismatch, never stash,
reset or silently commit. `--git=stage` retains its existing behavior. An
explicit reviewed `commit` API operation is available for the user to commit
only selected changed project files on the current branch, with expected head
and file revisions; it MUST never stage or commit unrelated checkout changes.
The user invokes that operation before merging or pushing local changes.

#### REQ: github-actor-authority-and-private-reads

GitHub project writes MUST pass through DataTug Cloud and a dedicated DataTug
GitHub App grant for the selected repository. A Firebase identity and App
installation alone do not grant repository write access. Each operation MUST
use an actor-bound GitHub App user access token, whose access is restricted by
both that user's current GitHub permissions and the App's selected-repository
Contents permission; loss of access must deny the next operation. Tokens and
rotating refresh credentials stay in server-side encrypted storage with
serialized refresh, never in project files, URLs, logs or the browser's durable
storage. The OAuth callback's unavoidable short-lived authorization code and
one-time state are removed from browser history/analytics-sensitive surfaces
before the app continues, never logged or persisted, and bound to the verified
actor (with PKCE where supported). A reused code or cross-actor state is
refused. If a single-use refresh rotates at GitHub but the replacement cannot
be persisted before a crash, the connection fails closed and asks the user to
reconnect; it never reuses an invalid token or guesses a successor. A private
project read uses the authenticated API and the same
repository-scoped authorization; it never falls back to a public mirror or
reuses a private browser cache entry without current actor authorization.
Private cache entries are actor-scoped and cleared on account changes or
revocation. Public
read may retain the existing anonymous path. GitHub repository collaborators
are not inferred to be Sneat contacts or Space members. The first release
selects an existing repository; a user may create one on GitHub and return to
select it. The repository must have an initialized default branch: an empty
repository with no HEAD is shown an explicit initialize-on-GitHub action and
cannot be saved until its branch exists. This avoids claiming a create-ref
compare-and-swap that GitHub does not offer for a branchless repository.
Creating a new GitHub repository through the DataTug API is a
separate capability requiring a separately reviewed repository-creation grant,
not implicit in Contents permission.

#### REQ: project-index-reconciliation

The selected project store is the source of truth for project files. DataTug's
user index records discoverability after the store write; an index failure
leaves a recoverable operation receipt and a retry repairs the index without
rewriting the Git commit. An external GitHub change is visible after read
refresh, and an old head or query revision is rejected for writes. Shared
project entitlements and contact assignments are checked separately by their
own Pro policy; a public project or owner-only private project remains
available on Free, and DemoDB query execution is not paywalled here.

## Acceptance Criteria

### AC: every-asset-is-a-file (verifies REQ:readable-file-layout)

**Given** the demo project after Phase 1's `Customer`/`Invoice` entities and three queries are added
**When** the project directory is listed
**Then** every asset appears as a readable file under a documented path, with no asset requiring a database to inspect.

### AC: mutating-command-stages-its-change (verifies REQ:git-stage-mutations)

**Given** a fresh clone of the demo project
**When** `datatug entity add` runs with `--git=stage`
**Then** `git status` shows the new entity file staged, and `--git=commit` is refused with exit code 2.

### AC: validator-rejects-committed-credentials (verifies REQ:secrets-never-in-project-files)

**Given** a `QueryDefTarget` with an inline password in a query file
**When** `datatug validate` runs
**Then** it fails, naming the file and field, and passes once the credential is removed.

### AC: github-project-reads-back-what-it-wrote (verifies REQ:github-hosted-projects-read-write)

**Given** a project written to a GitHub repository through `GhStorage`
**When** the same repository is opened for reading
**Then** `FileExists`/`OpenFile` return the written content instead of panicking.

### AC: openvaultdb-project-store-round-trips (verifies REQ:project-store-pluggable-backends)

**Given** a project configured to use a remote OpenVaultDB vault as its store
**When** an entity is added and then the project is reopened
**Then** the entity is read back from the OpenVaultDB store, and OpenVaultDB's own access policy is honoured alongside DataTug's.

### AC: project-store-browsable-as-a-source (verifies REQ:project-store-explorable-inside-datatug)

**Given** a project whose records live in an inGitDB-backed store
**When** the user opens "Project store" as a source inside DataTug
**Then** the entities/queries/policy collections appear in the grid, and a value in one (for example a query's `Customer.ID` parameter) is selectable and drives the same related-data and applicable-query loop as any other source.

### AC: knowledge-change-has-a-diff-view (verifies REQ:diff-review-affordances)

**Given** a captured query pending review
**When** a reviewer opens the diff affordance
**Then** they see the exact file-level Git diff the change would produce, before approving it.

### AC: one-query-save-contract-round-trips (verifies REQ:one-project-mutation-contract)

**Given** a query using `federation.ovdbBaseUrl` and `federation.tables` with a DTQL body
**When** it is created and updated through the same `save_query` DTO against a local clone, GitHub project, or Firestore project
**Then** the metadata and `.query.dtql` body round-trip without field loss, each response returns an opaque revision, stale updates return a typed conflict without a torn pair, legacy routes retain their responses, and a repeated operation ID returns the original result without a second Git commit.

### AC: project-create-uses-trusted-api (verifies REQ:one-project-mutation-contract)

**Given** a signed-in user choosing a Firestore project, an authorized GitHub repository with an initialized selected branch, or a local clone
**When** they create a project through the selected DataTug API endpoint
**Then** the store contains a valid project manifest, the appropriate project index resolves it after retryable reconciliation, GitHub creation uses the user's selected repository and branch with current user permission, and the DataTug client sends every project mutation through the selected DataTug API without a direct protected-store write path. A GitHub user may still edit their own repository through other tools. Later editable asset mutations use the same authorization and revision pipeline as their adapters migrate; first-slice acceptance claims only project creation and saved queries.

### AC: branchless-repository-needs-initialization (verifies REQ:github-actor-authority-and-private-reads)

**Given** an authorized GitHub repository with no branch or HEAD
**When** the user tries to create a DataTug project in it
**Then** the API refuses without writing a manifest or index and the browser offers an explicit GitHub initialization path; after the default branch exists, creation uses its observed head and returns a typed conflict if that head races forward, without overwriting the other change.

### AC: remaining-editable-assets-use-common-pipeline (verifies REQ:one-project-mutation-contract)

**Given** project metadata, query folders, boards and entities that the client can edit in a supported store
**When** each asset is changed through its DataTug API operation
**Then** the same project reference, actor authorization, validation, revision conflict and operation replay rules apply, the selected Git branch is pinned for Git-backed stores, and no DataTug browser-direct protected-store mutation path remains. This is a phased parity criterion, separate from the first project-create and saved-query delivery slice.

### AC: working-branch-changes-merge-or-propose (verifies REQ:provider-capabilities-and-branch-workflow)

**Given** a Git project with source and target heads observed by the user
**When** the user creates/selects a branch, makes two saved-query changes, reviews the diff, and requests merge to the chosen target
**Then** each save changes only the selected branch, the reviewed heads are checked again, a successful merge advances the target without force, a protected target returns an explicit proposal or pending state and URL, and a stale head leaves both branches unchanged. The Firestore project advertises and enforces no branch capability.

### AC: local-current-branch-pull-push-is-safe (verifies REQ:local-current-branch-synchronization)

**Given** a local clone served by the CLI API with a configured upstream
**When** the user pulls or pushes the current branch with expected heads
**Then** fast-forward pull and normal push succeed when the heads match; dirty-overwrite, divergence, stale-head and remote-rejection cases return explicit errors without stash, reset, force or loss of local changes.

### AC: local-reviewed-commit-scopes-files (verifies REQ:local-current-branch-synchronization)

**Given** two edited project query pairs and an unrelated modified file in a local clone
**When** the user reviews the project diff and explicitly commits the two selected pairs through the CLI-served API with the observed head and file revisions
**Then** the commit contains exactly those four query files, the unrelated file remains untouched, and a stale head or changed selected file prevents the commit without loss of edits.

### AC: github-revocation-and-private-read (verifies REQ:github-actor-authority-and-private-reads)

**Given** a selected private repository with a DataTug GitHub App installation and two signed-in users, one with GitHub write access and one read-only
**When** each reads and tries to save, and the writer's GitHub access is then revoked
**Then** the read-only user can read only while authorized and cannot save, the former writer's next operation is denied, an unrelated installation or repository is rejected, and no private bytes are served from a public mirror or an earlier actor's browser cache after account switch or revocation.

### AC: committed-retry-repairs-index (verifies REQ:project-index-reconciliation)

**Given** a valid GitHub project mutation whose commit succeeds but its index update or HTTP response fails
**When** the same operation ID is retried
**Then** the API rechecks the current actor's authority, discovers the committed result, repairs the index, returns its original revision and branch head, and creates no second commit. The same operation ID with a changed payload, branch, project, store, operation kind or actor returns a typed conflict without a write. An out-of-band branch edit then appears on refresh and makes an old-head save conflict.

## Not Doing / Out of Scope

- Automatic `--git=commit` for local CLI changes: saving locally stages
  changes and leaves them uncommitted. A separate explicit, reviewed,
  file-scoped API commit is in scope for the user's branch workflow. Saving to
  a remote GitHub branch via the Cloud API necessarily makes an atomic Git
  commit and is also in scope.
- A hosted database as the canonical project store; every backend (local, GitHub, OpenVaultDB,
  DataTug Cloud) exposes the same file-shaped project, never a proprietary schema.
- Client-direct writes to a protected project store; writes to a hosted or remote store still go
  through the trusted server (P10), matching the sneat-go writes-go-through-endpoints rule applied
  to DataTug's own architecture.
- Reimplementing version control inside DataTug; diff/review rides on Git, it does not replace it.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
