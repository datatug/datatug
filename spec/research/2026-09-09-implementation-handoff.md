# Implementation handoff to Fable: corrected DataTug and Dashboardius specifications

**Date:** 9 September 2026. **Audience:** Fable and implementation agents continuing
DataTug/Dashboardius work. **Scope:** documentation corrections following the
[independent review](2026-09-09-independent-review.md), authorized by the founder.

The founder requested: “Can you fix all inconsistencies and improve specs and plans?”
and then “commit and merge to remote and local main once ready.” This report describes
the resulting specification changes. It is not evidence that the implementation has
been repaired, that a release exists, or that the product journeys pass. Read current
remote refs and the accompanying delivery report before choosing an implementation base.

## What Fable should do first

1. Read the corrected [Phase 1 Feature](../features/core-investigation-loop/README.md),
   its [normative API/context contract](../features/core-investigation-loop/api-contract.md),
   and [Phase 1 corrective tasks 12–17](../plans/2026-09-09-phase-1-core-investigation-loop.md).
2. Refresh the review's Mac/VM reconciliation table against current remote refs and
   working-tree status. Preserve staged CLI follow-ups and pending ACL/provider work.
   Do not reset, clean or blindly merge a branch based on its old name or reported status.
3. Reconcile existing components to the corrected contracts. Do not restart the
   implementation, redo competitor research, or begin Phase 2 on the strength of
   completed component tasks alone.
4. Verify every applicable AC, including Supporting ACs in the Phase 2 plan; these
   are mandatory cross-capability acceptance, not optional notes. Distinguish schema
   fixtures/unit tests from the real browser journey.

## The important execution changes

| Previous inconsistency | Corrected requirement | Implementation consequence |
|---|---|---|
| Browser and server each implemented different shapes | One normative transport appendix: exact envelopes, method/body placement, errors, values and source identity | Inventory every caller; cut over server and clients together. Do not support both old shapes indefinitely or add another undocumented adapter. |
| SQL execution echoed parameters without applying them | Returned values must be execution-confirmed; native SQL uses driver-bound arguments | Extend the pending CLI argument work through the server execution path. Test two different parameter values produce the correct different rows. |
| Client-origin claims appeared to be server provenance | Binding origins explicitly distinguish client-reported facts from server defaults | Preserve useful origin labels, but do not describe a client-supplied selection/context origin as server-attested evidence. |
| HTTP adapter completion implied browser support | All three approved sources participate in the real Phase 1 browser journey | Wire HTTP saved-query dispatch and source resolution without fabricated SQL database fields. |
| Semantic endpoints used a separate filesystem convention | One authorized resolver adapts existing project source/environment/target records | No second persisted registry or hidden demo-only source path. |
| Multi-target queries required a choice but exposed no choices | Candidate carries authorized target options and needs-target state; execution validates the chosen source | Render a real target selector and structured TARGET_REQUIRED feedback without leaking hidden targets. |
| Session context was unscoped and first-match binding could win | Typed context is keyed by agent/project/environment/securityContextId, with explicit conflicts and required-input rules | React to store changes, reject stale responses, clear old-principal facts, and do not coerce number/string identity. |
| A warning substituted for row/column protection on native SQL | Protected execution refuses unconstrainable native SQL; explicit privileged opaque execution is a different capability | Preserve the stronger pending ACL behavior. Source read permission alone is not an opaque-query grant. |
| Principal ID was treated as its role | Fixed identity includes explicit roles/groups; agent-info reports effective identity and capabilities | Run the support fixture with `--as support --role support`; do not infer membership from its name. Result execution profile describes the actual executor, not all session capabilities. |
| Customer filtering was assumed to secure Invoice/support-note data | Every source/table used by the restricted journey has an explicit policy | Correct fixture policies and negative tests. Do not expect Invoice rows to hide a Customer field that was never selected. |
| A selected customer ID magically supplied country/currency | Country must be selected or retrieved; currency must come from an actual response | Repair J2 and expected counts using the pinned real dataset. No hard-coded CAD inference from CustomerId. |
| HTTP failures silently became successful snapshots | Live failure is honest; a snapshot requires explicit separate selection and retains its timestamp | Keep deterministic offline tests, but separately prove live HTTP. Enforce address/redirect/time/size limits and provider capabilities. |
| Lookup IDs/counts could imply authority | Revalidate each lookup in current scope and policy; counts are bounded and authorized or unavailable | Test forged/cross-scope handles and hidden-row count leakage. |
| J4 was disabled and missing setup could skip required CI | Required setup failures fail the job; all real journeys run from their actual entry points | Complete J1–J4 without intercepted responses or deep-link shortcuts and record companion revisions. |
| Task 10 and task 17 could publish at different stages | Task 10 produces journey evidence; task 17 alone owns final candidate/publication/receipt checks | Do not deploy before epilogue and integrated-candidate acceptance. |

## Revised task order and completion rules

Historical completed tasks 1–11 retain their lifecycle records and implementation
provenance. They describe component delivery, not acceptance of the corrected Feature.
SpecScore completed tasks are terminal; the corrections are explicit new tasks rather
than silently reopening or erasing history.

| Corrective task | Work | Dependencies |
|---|---|---|
| 12 | Freeze transport and unify source resolution | Landed model task 5 |
| 13 | Converge protected execution and real parameter binding | 12 |
| 14 | Complete bounded HTTP browser execution | 12, 13 |
| 15 | Converge browser semantics and isolated typed context | 12 |
| 16 | Correct demo policies and real transport acceptance | 13, 14, 15 |
| 17 | Prove complete candidate, epilogues and publication | Existing harness/UI tasks plus 16; exact dependencies in plan |

Tasks 13 and 15 may run concurrently in isolated worktrees. Task 14 follows task 13.
One landing owner coordinates each repository/target. A provider wave tested with
local GOWORK links is not published dependency convergence. Record local validation,
candidate CI, remote landing, CLI release, app deployment and smoke evidence separately.

## Code and pending work to preserve

- Shared semantic resolver/model work, query folders and declared mappings.
- The DALgo HTTP adapter, configured descriptors and explicitly labeled snapshots.
- The external datatug-core module; never resurrect the deleted vendored copy.
- Staged CLI saved-query/native-argument follow-ups and semantic demo test deltas.
- The pending VM layered ACL provider/SDK/server/UI chain. Port old vendored-core
  configuration references upstream and reconcile overlaps with current server/auth code.
- Existing grid, semantic panel, query form and context components; adapt their wire
  contract rather than rebuilding the product.
- The already-landed Go board-parameter correction.
- Dashboardius reducers, widgets, stable metric slots, table alternatives, Undo/Reset
  and honest scripted-demo labels.

The [review reconciliation table](2026-09-09-independent-review.md#mac-versus-vm-merged-versus-pending-reconciliation)
contains the dated worktree paths and source SHAs. Recheck them; they are discovery
leads, not permission to discard a branch or a promise that it is still unmerged.

## Roadmap and Phase 2 changes

The [roadmap](../../ROADMAP.md) now sequences bounded observable outcomes:

1. Phase 1: converge and prove the existing investigation.
2. Phase 2: capture as **project query**, useful folder descriptions, provenance and
   bounded mapping review, then close/reopen and second-investigator reuse.
3. Phase 3: human-guided query sequences, with explicit branches and terminal outcomes.
4. Phase 4: standalone project grants and distribution, including revocation/isolation.
5. Phase 5: real Dashboardius query → card → save → reload → authorized colleague.

Phase 2 has an explicit SpecScore prerequisite on Phase 1; readiness must remain
blocked until that predecessor completes. Its former HTTP credentials, richer types,
time-window semantics, agent/TUI parity and command cleanup tasks remain named deferred
work requiring separate activation/plans. They were not removed to satisfy a task cap.
No Phase 1 source, isolation rule or security AC was deferred.

Save as project query means an authorized atomic project-file write and a visible Git
diff, not an implicit commit or unimplemented hosted team grant. Failed, denied or
conflicting writes preserve the user's draft. Capture stores definitions and provenance,
not cached restricted result rows or credentials.

## Dashboardius changes and remaining decisions

The hub [dashboards Feature](../features/dashboards/README.md) owns the DataTug boundary.
Dashboardius now has its own [persistence Feature](https://github.com/sneat-co/dashboardius/blob/main/spec/features/query-backed-board-persistence/README.md)
and [journey plan](https://github.com/sneat-co/dashboardius/blob/main/spec/plans/2026-09-09-query-backed-board-persistence.md).

DataTug owns project/board contracts, storage, query execution and policy. Dashboardius
owns presentation and editing. Its prototype TypeScript model is not a field-for-field
durable Go contract. Reconcile supported fields and row identity explicitly through a
real Go/store/browser round trip. Keep runtime results and presentation view state
separate from the persisted asset.

Loading board JSON authorizes a project asset; it does not perform row filtering.
Each widget re-execution independently applies the recipient's source/row/column policy.
No sender-authorized cached result is restored from the saved board.

The **SQL-widget binding choice remains open** in the hub Feature: library-query
reference, board-level target, or widget-level database with view-time environment.
The reference approach is recommended, not founder-selected. Do not invent wire fields
or choose an option while implementing unrelated Phase 1 work. The permanent ban on
embedding charts in DataTug was removed: Phase 1 excludes boards, while later embedding
may reuse Dashboardius presentation without creating a second renderer.

The standalone multi-space project rule remains explicit. Exact grant representation,
role composition and the wider personal/team/project/organization promotion ladder
remain draft design choices; they are not extra founder rulings.

## Research, positioning and hierarchy changes

The hierarchy retains P1–P15 IDs but distinguishes user capabilities, enablers,
surfaces and experiments. P6 context and P7 executable knowledge organize the primary
value. Existing CLI ideas and overlapping AI surface Features are mapped to owners
instead of disappearing behind a new hierarchy. Initial audit findings are labeled
historical; they no longer claim that later-landed components do not exist.

The marketing repository now documents disjoint scoring weights, the exact formula,
Unknown/N/A handling and portfolio overlap. Historical scores are reproducible from
retained cells. Runme and dbt/MetricFlow are targeted, explicitly unscored comparators,
not a claim that all market research was independently redone. Absolute white-space
and confirmed-demand language is replaced by evidence-qualified hypotheses.

The next demand test is another engineer actually reusing a captured investigation;
measure time, copying, bad bindings, authoring overhead and repeat use against the
team's current tools. No invented customer outcomes or success thresholds were added.
Licensing remains DataTug CLI Apache 2.0 and private web UI for now; future web licensing
is undecided. Keep ahead-of-code capabilities labeled planned.

## Verification of this documentation change

Before handoff, the edited hub passed SpecScore structural/traceability lint, internal
file-link checks and task dependency-cycle checks. Phase 2 readiness correctly reports
the executing Phase 1 prerequisite. A separate consistency review checked the new
contracts and its six actionable findings were addressed. Dashboardius documentation
passed SpecScore; WB additionally ran its existing test/build checks for landing.
The marketing lane checked changed links and reproduced all 84 matrix scores.

These are documentation and repository-validation results. **They do not close any
product journey AC.** Implementation must still prove the corrected contracts and
publish the actual product candidate under the normal authorization/protection rules.

## Open Questions

No new question blocks the corrected local Phase 1 contract. The owning Features retain
broader board binding, grant representation, hosted identity, HTTP credentials, embedding
and future web-license decisions. Do not substitute an implementation assumption for a
founder answer; continue independent work while an actual dependent decision is pending.
