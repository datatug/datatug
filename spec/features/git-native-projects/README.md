---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Git-native projects

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/git-native-projects?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P9 is Git-native projects: every knowledge asset is a readable, versioned, diffable, portable file — no secrets in project files — and the project store itself, wherever it lives, is explorable inside DataTug.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

It is the pillar that makes P7's compounding library real rather than aspirational: a captured query is worthless as an accumulating asset if it cannot be reviewed, diffed and moved between local, hosted and Git-backed stores. Initial audit snapshot (before later 9 September landings): the local filestore is the only fully working `Store` backend; `--git=stage` only covers `entity add`; secrets can still land in committed query files because `Project.Validate()` never validates `Queries` at all; GitHub- hosted projects can write but not reliably read back (`GhStorage.FileExists`/`OpenFile` panic). Two founder additions from 2026-09-09 extend this pillar: a remote OpenVaultDB project store, and making the project store itself explorable as a source inside DataTug. Phase 0–1 add the secrets guard; Phase 2–3 build the pluggable store backends and self-exploration; Phase 4 adds diff/review affordances. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`) link to it.

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

## Not Doing / Out of Scope

- `--git=commit` from DataTug itself: DataTug stages changes; committing remains a deliberate,
  separate human or CI action.
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
