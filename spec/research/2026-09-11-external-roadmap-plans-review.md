# External review of the DataTug roadmap and plans (2026-09-11)

Run by the founder with an external frontier model, using `roadmap-plans-review-prompt.md`. A second reviewer, Sol, checked traceability and challenged the findings. The text below is the report exactly as the founder relayed it, followed by the founder's answers, verbatim.

---

**Yes—there are foundation issues to fix before Incidentius execution.** The product boundaries are clearer, but the evidence model, cohort contract and acceptance coverage still disagree in important places.

I reviewed freshly fetched hub `00a53a6`, Dashboardius `253db35`, CLI `d34fd18`, and the specified Workbench branch `2f04e4a`. All requested files were accessible. A separate Sol reviewer checked traceability and challenged the findings. This is a document review, not a fresh certification of Phase 1's running software. No repository files were changed.

**Verdicts and sizing**

Sizes describe the planned work, not remaining work.

| Roadmap / plan | Verdict | Size / critical path | Reason |
|---|---|---|---|
| Roadmap | **Fix then execute** | L / L | Direction is coherent, but stale gates and the unplanned C3 cutover distort sequencing. |
| Phase 1 | **Ready as the documented baseline** | L / L | Task 17 records candidate, journey and publication evidence; later phases must retain that boundary. |
| Phase 2 | **Fix then execute** | M / M | Capture is appropriately bounded; authorization dependencies and stale approval text need correction. |
| Query-storage foundation | **Fix then execute** | M / M | Strong recovery and compatibility design; its component evidence must not close browser ACs. |
| Incidentius MVP | **Rework before execution** | L / L | Several promised outcomes cannot yet be implemented consistently from the contracts. |
| Query-backed boards | **Fix then execute** | L / L | Remove the already-settled decision gate and specify migration/conflict behavior. |
| Incident dashboards | **Fix then execute** | L / L | MVP versus NEXT is clearly separated, but external prerequisites and renderer destination need tightening. |

**Findings, ranked**

Paths below use **H** = `datatug/datatug`, **D** = `sneat-co/dashboardius`, **C** = `datatug/datatug-cli`, **W** = `sneat-co/workbench`.

**F1 — Blocker · CONFIRMED: metric history cannot produce its promised values.**
Location: H:spec/features/evidence-records/README.md, execution-record schema and series requirement (L187); Incidentius tasks 2/6; incident-dashboard task 2.

> "A numeric projection over a series — `{column, aggregate}`"

Records contain `rowCount` and a fingerprint, but no arbitrary column values; snapshots are optional and expire. A recorded `SUM(total)` cannot subsequently be recovered from that fingerprint. Series grouping also names only query ID and bindings, allowing different environments, revisions or policy scopes to become one misleading line.

**Fix:** define retained, policy-governed scalar measurements at execution time, including projection identity, source scope, query revision and completeness. Specify unavailable/expired evidence explicitly. Keep unrelated scopes separate. The proposed "zero for 30 minutes" criterion also needs an observation-window rule; one zero measurement cannot prove it.

**F2 — Blocker · PLAUSIBLE: target authorization does not protect a note's other contents.**
Location: H:spec/features/annotations/README.md, provenance and visibility (L166); incidents AC `events-redacted-under-policy`.

> "visible only if the reader could retrieve that value from at least one source"

An annotation targeting visible `Customer.ID` can contain hidden `Customer.Email` in its body or copied context facts. Target-only authorization would expose it, while the incident AC promises that a note quoting hidden email is withheld.

**Fix:** define sensitivity and authorization for referenced facts and free-text bodies, including search, projections and streams. Add the exact public-ID/hidden-email negative scenario. This is a contract gap, not a demonstrated implementation leak.

**F3 — Blocker · CONFIRMED: cohort and cross-store evidence contracts remain incomplete.**
Location: H:spec/features/compare/README.md, Open Questions (L291), `SideSpec` and `CompareRequest`; Incidentius tasks 1/4/8.

> "`Binding.value` is a single `TypedValue` and the transport appendix defines no set binding."

The journey binds entire cohorts, but the wire contract cannot express the promised expansion. Separately, record references are plain IDs despite project-unique execution IDs and independently routed evidence stores; compare also lacks an incident reference in its request.

**Fix:** freeze bounded set-binding semantics and source refusal behavior before browser work. Define qualified execution references and propagate store/project/incident scope through checks, annotations and compare. Prove two stores with colliding short IDs cannot cross-resolve evidence.

**F4 — Blocker · PLAUSIBLE: incident durability assumes fewer writers than the CLI permits.**
Location: H:spec/plans/2026-09-11-incidentius-mvp.md, Risks (L374); C:spec/features/cli/incident/README.md, server-backed execution (L103).

> "Appends are serialised by the server per incident; the MVP is single-server."

The CLI may invoke the server package in-process when no server runs. Two CLI processes can therefore compete for sequence allocation and append. Retry after a lost response, partial append, and interruption midway through merging two incidents lack recovery contracts.

**Fix:** require one store writer or cross-process exclusion, idempotent mutation IDs, durable append/recovery rules and recoverable merge semantics. Test simultaneous CLI mutations and interrupted merge. Imported events must not retrospectively alter pre-merge belief or reactivate copied grants.

**F5 — Should-fix · CONFIRMED: immutable records are later mutated by retention.**
Location: H:spec/features/evidence-records/README.md, immutability and retention (L75); Incidentius task 2.

> "no field may be edited or removed"
> "with `snapshotRef` cleared"

Both cannot hold. Annotation history similarly relies on Git history even though project writes deliberately do not commit automatically.

**Fix:** keep immutable execution receipts and separate snapshot availability/expiry state. Preserve annotation revision bodies explicitly, or require a real history mechanism; an uncommitted overwrite is not history.

**F6 — Should-fix · CONFIRMED: the browser journey has an escape hatch into CLI-only acceptance.**
Location: H:spec/plans/2026-09-11-incidentius-mvp.md, task 8 (L242).

> "Compare UI is NEXT"

The journey promises affected/control comparison, but no task owns its MVP web interaction. The whole-journey AC is loose enough to pass using CLI comparison.

**Fix:** include a minimal web compare action and honest result table in B1; defer elaborate grids and charts. Name the browser action in task 12's acceptance.

**F7 — Should-fix · CONFIRMED: dependency metadata contradicts the intended parallel work.**
Location: H:Incidentius plan, prerequisites and dependency explanation (L10).

> "**Prerequisite Plans:** … phase-2-library-depth-and-provenance"

The prose requires only Phase 2 task 2 for checks/annotations, but metadata gates the entire Incidentius plan on all Phase 2.

**Fix:** represent the bounded external dependency. Also encode:

- Phase 2 task 6 → shared authorization/write foundation from task 2.
- Incidentius task 4 → task 3's event API.
- Incidentius task 8 → the resolved cohort contract.
- Incidentius tasks 5/7 → Phase 2 task 2.
- Incident-dashboard tasks → the specific evidence, metrics, board API and widget-contract providers.

Preparation may remain parallel; dependent acceptance cannot.

**F8 — Should-fix · CONFIRMED: C3 has a cutover inventory but no executable owner plan.**
Location: H:spec/features/product-profiles/README.md, cutover and three-profile AC (L243); Incidentius task 9; incident-dashboard task 1.

> "all three are served from the same build artefact"

Task 9 implements two profiles while claiming that three-deployment AC. C3's deployment transfer, old-shell removal and domain smoke have no task chain. The renderer plan also builds an intermediate library in the old repository before relocating it.

**Fix:** give C3 an explicit plan with deployment-before-removal acceptance. Separate B1's two-profile milestone from family cutover acceptance. Build the shared renderer directly in its final Nx workspace unless a demonstrated constraint requires the intermediate move.

**F9 — Should-fix · CONFIRMED: valid AC identifiers hide partial proof.**
Location: W:query-storage-foundation/README.md, Deferred AC Coverage (L247).

> "This plan cannot complete … the Feature acceptance criterion."

Nevertheless, all five storage tasks say they verify the browser capture AC. Incidentius tasks 8/11 similarly cite the whole journey while delivering comparison infrastructure or fixtures.

**Fix:** distinguish component contribution from final verification, or add component ACs. Apply the same correction to Incidentius tasks 2/7/13 and the board plan's decision task. Keep integrated proof on the final journey tasks.

**F10 — Should-fix · CONFIRMED: compare's requirement and test demand different event kinds.**
Location: H:spec/features/compare/README.md, evidence requirement and AC (L267).

> "`assertion.kind = observation`"

The requirement and Incidentius task 8 specify `deterministic-result`.

**Fix:** correct the AC and pin one shared event fixture. Also distinguish live sides, which create records, from snapshot sides, which reuse existing records without re-execution.

**Traceability results**

All **64 tasks** contain `Verifies`; every referenced ID and all **40 Supporting AC references** exist. No local dependency cycles were found. No task literally verifies nothing, but F9 identifies tasks whose declared proof exceeds their output.

Every source-feature AC is assigned or explicitly deferred. Remaining coverage gaps are:

- **No direct AC:** personal-query storage; Resolution Record; API/CLI requirements in evidence, checks, annotations and compare; preserving schema comparison; profile workspace reuse and brand/build-info requirements. Some have task prose—this is missing acceptance coverage, not proof of missing implementation.
- **Unassigned sibling ACs:** the two C3 cutover ACs; `surfaces#profile-does-not-hide-a-surface`; the eleven `cli/incident` ACs have no explicit cross-repository task references.
- Existing context, access-control, provenance and board sibling ACs often rely on equivalent Phase 1/product ACs without explicit equivalence links. Record those mappings rather than creating duplicate implementation tasks.
- Personal-query storage needs its own explicit acceptance owner; it must not silently expand Phase 2's **project-query-only** scope.

**Critical path to real value**

The first real-user value is already Phase 1's useful result and authorized follow-up lookup. The next valuable moment is **saving an investigation finding so it survives reopening**.

1. Complete storage-foundation tasks 1–5.
2. Complete Phase 2 task 2's authorized endpoint and browser capture.
3. Deliver folder discovery and capture provenance—tasks 3–5—and an early capture/reopen/rerun gate.
4. In parallel, Incidentius tasks 1–3 establish safe events and evidence.
5. Bring a minimal task 9 incident page forward: create → run an existing check → attach evidence → reopen.
6. Add cohorts, annotations, metrics, access workflow and recurrence; finish task 12's full acceptance.

Author the demo story early, before fixture implementation. Neither mapping-review polish nor the Dashboardius shell cutover should delay that first incident pilot.

**Top three cuts**

1. **Move Phase 2 mapping correction behind first capture/reuse delivery.** Saves a separate mutation model and review UI from the first-value path.
2. **Use explicit recurrence links before broad similarity ranking in the first pilot.** Saves indexing, ranking and explanation work; retain ranking in a subsequent approved slice.
3. **Remove the renderer's intermediate repository relocation.** Saves packaging, consumer rewiring and duplicate migration review.

These are proposed scope/sequencing amendments, not permission to silently drop current ACs.

**Founder decisions needed**

- **May production snapshots ever enter Git?** Options: never in MVP **(recommended)**; explicit opt-in; unrestricted. Blocks Incidentius task 2's storage policy.
- **Which next pilot comes first?** Options: a narrow Incidentius case that survives reopen **(recommended)**; complete library refinement; complete the entire B1 bundle before user feedback. This selects sequencing, not product boundaries.

Store-key spelling, lock choice, revision tokens and fingerprint mechanics are engineering decisions. Branding need not block a configuration-selected Incidentius pilot.

**Largest schedule risk:** making historical evidence trustworthy under current permissions across stores, retries and policy changes—not building the incident screens.

**Nits, together:** Phase 2 still says "Draft" inside an Approved plan; ownership notes still await already-recorded Phase 1 acceptance; the roadmap and board task 1 still reopen the settled SQL-binding decision. Retire `datatug-core/docs/ROADMAP.md` as an active roadmap and replace it with a historical label plus hub link—it contains only the old decoupling bullets.

---

## Founder answers (verbatim, 2026-09-11)

- May production snapshots ever enter Git? "Not in MVP"
- Which next pilot comes first? "complete library refinement - e.g. DataTug improvements first. complete Phase 2 and review the core as a coherent experience before prioritizing Incidentius."
- Earlier the same day, on the minimal web compare in the Incidentius MVP: "compare - A"
