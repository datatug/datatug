# Layered-ACL chain reconciliation: inventory and port plan for Task 13

**Date:** 9 September 2026. **Audience:** whoever picks up Phase 1 Task 13 ("Converge
protected execution and actual parameter binding"). **Scope:** a read-only inventory of
the pending layered-ACL provider/SDK/server/UI chain against current `main` on every
repo it touches, plus a port plan for reconciling it with the corrected contracts. No
code was changed to produce this note; the chain itself is landed by a separate
implementer (codex) agent, not by this note's author.

**dal-go/dalgo PR #158 is merged** (`8dece8b02bdba3fc0fdaced8f1a0441270991706`,
2026-09-09T16:59:33Z) — the provider foundation of this chain is in. **The remaining
provider/SDK branches and the datatug-cli/datatug-apps ACL branches are being landed by
the implementer (codex) agent that has been driving this chain, not by this task or its
author.** This note is inventory and a port plan only; it does not itself land, rebase,
or open a PR against any of the branches below.

## The provider chain (dependency order: provider-first)

| Repo | Branch | Head SHA (`git ls-remote`/`gh pr view`) | PR | State | What it adds |
|---|---|---|---|---|---|
| `dal-go/dalgo` | `layered-acl-query` | `8dece8b0` (merge); branch tip `65c4824b` | #158 | **MERGED** | ~10,177 additions: `access/` (portable ACL policy loading, protected-assessment coordinator, scoped field/collection masks, typed-principal-to-owner-identity binding, fail-closed admission) and a new `dtql/authorization/` subpackage (DTQL query authorization results, wire-contract coverage). This is the whole of what "dtql `layered-acl-design`" refers to — not a separate repo. `datatug/dtql`, `datatug/dtql-go`, `datatug/dtql-org`, `datatug/dtql-backstage` were each checked directly; none has a `layered-acl-design` branch, only `main`. |
| `dal-go/dalgo2sql` | `layered-acl-query` | `45fe36ce` | none | pushed only | Adapts `database.go`, `database_options.go`, `reader_base.go`, `reader_records.go`, `reader_recordset.go`, `sqlite_emit.go`, `transaction.go` to the new `access`/authorization APIs. Was 2 ahead / 2 behind its own `main` at check time; `main` has since moved further (`98b9d6f7`, past the v0.11.7 `TextQuery` binding fix). |
| `ingitdb/dalgo2ingitdb` | `layered-acl-query` | `de6f64c5` | none | pushed only | 3 ahead / 0 behind its own `main` — clean. `README.md`, `access_manifest.go` (+test, new), `database.go`, `query.go`, `recordset.go`, `tx_readonly.go`. |
| `openvaultdb/openvaultdb-go` | `layered-acl-query` | `79988958` | none | pushed only | 5 ahead / 1 behind its own `main`. Server/SDK layer: `pkg/auth/*`, `pkg/core/*` (DTQL execution), `pkg/manifest/*`, `pkg/mount/*`, `pkg/server/*` (dtql handler, identity, tokens). |

Order is provider-first: `dalgo` (landed) → `dalgo2sql`/`dalgo2ingitdb` (adapters
consuming dalgo's new APIs) → `openvaultdb-go` (SDK/server on the adapters) →
`datatug-cli` (`pkg/openvaultdb` + a proxy) → `datatug-apps` (UI on the proxy). None of
the three still-pending provider repos has an open PR yet.

## VM worktrees: datatug-cli and datatug-apps

**`datatug-cli/.worktrees/datatug-layered-acl-query`** (`a7e1f62`, two commits) vs.
current `main` (`e6fc6f1`): 12 files, +1111/−11.

- `pkg/datatug-core/dtconfig/*` — the **vendored path**, fully deleted on `main`
  (confirmed: `git show origin/main:pkg/datatug-core/dtconfig/settings.go` → file not
  found). The external module `github.com/datatug/datatug-core@v0.23.0`'s own
  `pkg/dtconfig/settings.go` already has a near-identical `Settings` struct, missing
  only this branch's new `OpenVaultDB` field.
  **Not a mechanical port**: `Settings.ResolveOpenVaultDBTargets()` imports
  `github.com/datatug/datatug-cli/pkg/openvaultdb` — a datatug-cli-owned package.
  Moving it into the external `datatug-core` module (which datatug-cli imports *from*,
  never the reverse) would invert the dependency direction. Either the field shape moves
  upstream without the CLI-specific resolver, or the whole extension stays a CLI-local
  sidecar — a design decision, not a rebase outcome.
- `pkg/server/openvaultdb_proxy.go` (+361, new) — **the core reconciliation problem**.
  References none of `secureread`, `accesspolicies`, `resolveServeSession`, or
  `pkg/api`. It is a fully separate authenticated-session/authorization path, built
  directly on dalgo's new `access`/`dtql/authorization` types, with zero overlap with
  the `pkg/secureread` + `pkg/accesspolicies` + `resolveServeSession` machinery that has
  since become the real, landed policy-enforcement path on `main`. This is what "reconcile
  overlaps with current server/auth code" means in practice: two competing
  implementations of the same concept, not a merge conflict.
- `pkg/server/http_server.go` — additive and low mechanical risk (a new
  `NewHttpServerWithOpenVaultDB` constructor, one extra route-registration call), but
  sits beside code the same-day PRs #199/#205/#209 also touched.
- `apps/datatugapp/commands/cmd_serve.go` (+47/−4) — **not mechanical**. The diff's
  context still shows `dtconfig.GetSettings()` called directly and a log line reading
  "parsed but not yet enforced (plan task 6 wires access policies)" — this predates
  `resolveServeSession`/`secureread.Session` entirely. Current `main`'s
  `serveCommandAction` is a substantially different function body.

**`datatug-apps/.worktrees/datatug-apps-layered-acl-query`** (`5ebb264`, four commits)
vs. current `main` (`c5cf301`): 17 files, +1853/−6, mostly new files under
`apps/datatug-app/src/app/openvaultdb/` plus e2e coverage. Only two existing files
touched (`datatug-app-routes.ts` +14, `main.ts` +5) — low mechanical risk.
**Design-change flag**: the new routes (`pwa/repo/:repo/agent/:agentId`,
`agent/:agentId`) register a second, competing URL convention alongside current
`main`'s `serveAgentURLs`, which prints `store/http-<host>:<port>`. Task 15 (browser
semantics/isolated context) will need to pick one.

Older worktrees not independently re-diffed here (time-boxed; already assessed in the
[independent review](2026-09-09-independent-review.md#mac-versus-vm-merged-versus-pending-reconciliation)):
`datatug-cli` WB `datatug-query-policies` (`f6d6d17`), the `datatug-apps` integration
worktree (`e45adcd`).

## Already satisfied on `datatug-cli` `main`

| Task 13 requirement | Satisfied by | Verified |
|---|---|---|
| Driver-bound parameters on the server path | PR #209 (merged 2026-09-09T16:45:49Z) | `dal-go/dalgo2sql` bumped to v0.11.7 (real `sql.Named`/positional binding); `pkg/secureread.RunNativeSQL` takes `args ...dal.QueryArg`; the CLI's own string-substitution workaround removed. |
| Explicit privileged opaque grants remain separate | PR #199 + PR #205 | `pkg/secureread`'s `LimitationNativeSQL`/`access.OpaqueQueryScope` pattern (an `opaqueQuery: true` rule set, distinct from ordinary path-scoped rules). |
| Policies per source | `datatug-demo-projects` commit `09407c38` "fix(demo): enforce support policy per source" | Authored 2026-09-09T17:27:51+01:00. |

What only the ACL chain still provides: (1) protected proxy execution
(session-token-gated, capability-checked reads — no equivalent in `secureread`), (2)
richer typed-principal/scoped-mask capability checks than `accesspolicies`'s current
YAML-rule model, (3) structured DTQL authorization at the wire level versus
`accesspolicies.Explain`'s text-line reporting.

## Conflicts against Task 12's transport freeze

1. **Envelope shape** — Task 12 freezes `pkg/apicontract` as the schema authority.
   `openvaultdb_proxy.go` defines and returns its own response shapes entirely outside it.
2. **Session/principal model** — `secureread.Session` (project + `--as`/`--role`/
   `--group`, YAML policy files) versus the ACL chain's URL-fragment session-token flow
   are two different answers to "who is asking."
3. **Source targeting** — Task 12's resolver maps `SourceRef.source` through the
   project's own environment/catalog records. `OpenVaultDBTarget{BaseURL, DatabaseID,
   TokenEnv}` is a second, independent, `Settings`-driven target registry.
4. **URL/routing convention** — `pwa/repo/:repo/agent/:agentId` versus
   `store/http-<host>:<port>`; Task 15 needs one, not two.
5. **Config ownership** — `OpenVaultDBConfig`/`WebUIOrigin` as `Settings` extensions
   collide with "schema authority = `datatug-core` module"; these fields cannot silently
   move upstream, nor silently stay local, without an explicit decision.

## Assumption

"This lane does only the part that does not depend on Task 12" (Task 13's own
dependency, verbatim in the plan) is read here as: inventory and port-planning is
in scope; resolving any of the five conflicts above is not, since that resolution
needs Task 12's frozen contract to exist first.
