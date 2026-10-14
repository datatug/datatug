# Context Bag delivery, reconciliation and owner review

2026-10-13 · M0–M3 implementation approved 2026-10-13 15:46 UTC

## Journey first

| Investigator journey | Observable good result | Work package |
|---|---|---|
| Create session; add Order 123 | Durable session/default bag, explicit environment; no mandatory alias | M1 |
| Apply; open Orders then OrderDetails | Visible pinned bindings; related rows appear without another Apply | M2 |
| Open Order overview dashboard; Run | Required input filled; result grid and 1–5 charts share manifest | M3 |
| Attach bag entry; ask a question | Compact chip; immutable send snapshot and execution receipt | M3 |
| Create Good/Bad and UAT/PROD bags | Named independent contexts; mismatched environment never binds | M2/M4 |
| Compare bags with same dashboard | Explicit independent scopes and results; no raw-ID cross-env union | M4 |
| Clear context → Confirm | Selected bag emptied in two clicks; other bags and history intact | M2 |
| Keep finding in Case; stop session | Case evidence links preserve inputs; Stop doesn't close Case | M3/M4 |
| Restart/reopen from authorized store | Context/history restored without rerunning anything | M1/M3 |
| Work offline with unavailable remote source | Context readable/editable per profile; query says unavailable | M1/M5 |

## Phases and work packages

Effort ranges are planning estimates, not delivery promises. They assume existing
Detective contracts remain usable and a small team with one integration owner. First
usable slice targets roughly 3–5 engineering weeks after approval; OVDB durability or
hosted admission gaps may change that. Each milestone ends in an exercised journey,
not just a library release. Keep the critical path narrow and re-estimate after M0.

| Phase / milestone | Scope and dependencies | Acceptance / exit | Estimate |
|---|---|---|---|
| 0 / M0: choose store profile | Benchmark OVDB inGitDB vs SQLite; conditional-write and crash-recovery spike; inspect latest Detective contracts | 1/300/5,300 values, 100 sessions, repeated revisions; restart/race/timeout/private-hosted probes; record chosen authority | 2–4 days |
| 1 / M1: sessions and durable bags | Session lifecycle; default bag; typed scalar/set; private OVDB records; local + hosted profiles; browser cache/draft boundaries | Two sessions survive restart; versions/receipts reject lost updates; no local CLI for hosted web; migrate old context explicitly | 4–7 days |
| 1 / M2: context across workspace | Top session tabs + All sessions; hybrid bag summary/drawer; Ionic-first row/header composition; multi-bag environment scope; stable IDs; Apply/master; scoped Clear; semantic direct relation; latest Query Author integration | One ID → Orders → OrderDetails; UAT cannot affect PROD; selected bag Clear is two clicks; close/reopen preserves session; keyboard/narrow layouts | 4–6 days |
| 1 / M3: complete useful slice | Composer references/send pins; parameterized Order dashboard with grid + at least one chart (two recommended); Case evidence bridge | Whole journey above through stop/restart in web and terminal; dashboard widget receipts share inputs; history doesn't change | 4–7 days |
| 2 / M4: investigate cohorts | 5,300 reference-backed inputs; compare bag dashboards; grid-to-bag; reviewed extraction; snapshot branching; direct annotations; ranges | All 5,300 used without prompt dump/truncation; compare authorized sides; source retention explicit; overlap explained | 1–2 weeks |
| 2 / M5: offline and history polish | Installable web shell/cache; local draft import; indexed indirect annotation discovery; retention/export | Cold offline reopen previously downloaded session; unavailable source distinct; index watermark; idempotent draft import | 1–2 weeks |
| 3: publish and manual handoff | Explicit selected-context publish; cross-device hosted use already exists from M1; manual linked refresh/diff; propose/apply | Permission preview, revocation/read checks; parent revisions don't silently mutate child evidence | Re-estimate after real usage |
| 4: collaboration and automation | Multi-writer offline replay/conflicts, optional live links, classification rules, scheduled checks | Only after observed demand and explicit conflict/permission policy | Deferred |
| Launch / L1: explain and promote | Website copy, feature page, homepage/navigation, painted marine art, live demo entry | Delivered behavior demonstrated and availability claims reviewed; no premature sync claims | After useful slice, near launch |

Dependencies: M0→M1→M2→M3 is the first journey; dashboard definitions/fixture queries can
be prepared with M1 after approval. M4's large-set executor depends on a reference-backed
input capability; it does not block a useful one-ID dashboard. Nevertheless 5,300 storage
and transport feasibility is tested in M0 so the chosen schema cannot paint us into a
corner. Local and hosted deployment are both first-slice goals, not postponed to Phase 3.
Cloud access does not mean collaboration/sync. Browser installable-shell work and richer
offline editing can follow while CLI local OVDB supports real offline sessions early.

The owner approved implementation of M0–M3 on 2026-10-13 15:46 UTC. The corresponding
Feature-sourced SpecScore plan belongs in the configured central plan store and must
link every acceptance criterion to an M0–M3 task or a concrete deferral. Root owns
implementation coordination and landing.

## First-slice dependency and acceptance supplement

Prepared during design review; this refines the approved M0–M3 scope.
Fresh fetches on 2026-10-13 confirm apps `a35c58e9da555eab9b1ed1cab102bfd14844180b`,
specifications `feb88bd9ddb8ac9f4dfcb94d38d77a43571cd2a0`, and Core
`f54cb46430d61dd26ebb5a35d93f623e1942b932`. Core now includes reviewed AI/knowledge
contracts from PR #362. Source availability is not evidence of deployed adapter support.

| Dependency / responsible area | Required before | Concrete proof after approval |
|---|---|---|
| OVDB command owner + local/hosted adapters | M1 write path | Same expected-version conflict and lost-acknowledgement recovery fixture passes against both authorities; bytes survive process restart |
| Session access and hosted admission | M1 cloud journey | Fresh browser creates a private session without a local process; unauthorized project member cannot list/read it; revocation blocks replay |
| Semantic binding + explorer/editor integration | M2 | Order identifier binds related OrderDetails through verified relation metadata; ambiguous tenant/order key blocks only affected consumer |
| Dashboard definition/runtime owner | M3 | One required order input; grid and two recommended charts consume the same pinned binding; changing UAT leaves an open PROD run fresh |
| Detective adapter + session input snapshot | M3 Case bridge | Private input snapshot precedes dispatch; duplicate retry finds the same accepted run; audit and run list omit raw values |
| Detective native Case linkage + session outbox | M3 Case bridge | Case-bound run and Case link finalize through Detective's transaction; kill the separate input-evidence handoff worker, then restart it without another source execution; denied handoff remains visibly blocked |
| Existing chat/investigation import | M1/M3 restore journey | Import twice yields one new session; old records remain recoverable; first reopen causes no query or Check execution |
| Client lifecycle/result routing | M3 whole journey | Start slow PROD run, switch session, Stop original; response stays with original receipt and does not overwrite the focused session |

The first release must exercise the single-order journey on local and hosted profiles
through web and terminal, including **do nothing after reopen**: context and history
appear without a hidden refresh, and no source query executes. Restore a session after
its source becomes unreachable; retained context remains readable and new execution is
explicitly unavailable. A successful client demo alone cannot substitute for durable
server receipts, current permissions and restart recovery.

M0 should explicitly resolve three possible schedule multipliers: whether hosted OVDB
can run commands at its transaction owner; whether the admitted source has actual
order/category rows for the dashboard; and whether Detective's accepted operation can
be looked up after an acknowledgement is lost. Until these are proved, the 3–5 week
range remains a hypothesis. Missing one does not justify dropping the paired local/cloud
journey; revise the scope sequencing and estimate for review.

Reuse the newly landed `detectiveai` extraction contracts in M4 only after the bag-target
command and source/quote-retention contract are explicit. Neither reviewed AI proposals
nor knowledge publication is required to get the first manual one-order journey working.
See the architecture's Detective input-evidence section for the proposed private
snapshot and cross-store recovery contract. No automatic Detective store migration is
included in this supplement.

Supplement review r1: independent source review found no blockers, majors or minors.
The principal separately checked the referenced Core contracts and corrected the
Case-bound finalization boundary before review. This is design validation only;
the acceptance probes above remain implementation work after approval.

## Storage spike must answer, not merely benchmark speed

For each profile report: command p50/p95; opening/page latency; total bytes and file count;
repeated Apply/Send overhead; crash recovery at staging/head/receipt/publication boundaries;
two stale writers; process restart; idempotent replay; clear/delete referential effects;
permitted/denied reads; cache grants and revoked authorization; source unreachable.
The 5,300 fixture has typed duplicate keys and composite tenant/order keys. Input cardinality,
payload cardinality and executor-consumed cardinality must agree. Never call a syntactic
query success proof of correct whole-set execution.

Provisional targets on a documented dev laptop: cached bag open p95 under 200 ms,
Apply local p95 under 300 ms excluding source execution, hosted command p95 under 1 s
excluding source execution. These are hypotheses for comparison, not product guarantees.
Failure of inGitDB at required durability is a blocker; slower speed triggers comparison
with OVDB SQLite. Index/SQLite-cache proposals must state which bytes are authoritative.

## Reconciliation with existing specs and WIP

| Existing item | Conflict / reuse | Proposed resolution |
|---|---|---|
| Investigation Context scope/sessionStorage | No session ID; value-derived IDs; immediate checkbox mutation | Add session/bag identity and revisions; import old facts; Apply gates new consumers |
| Query Author a35c58e9 | Parameters/Sources/Related, mobile result navigation, row filter Undo | Reuse shell and receipt accordions; separate Add to Bag from query text filters |
| Core inline sets | Maximum 500 | Keep bounded inline path; introduce authorized collection reference for larger sets |
| Existing chat sessions | Multiple durable sessions, no requested lifecycle | Add lifecycle and bags; do not pretend sessions themselves are entirely new |
| Detective/Checks WIP | Case typed context, local/cloud repository and authoritative runs | Session working context; Case evidence snapshots; reuse domain semantics; no parallel evaluator |
| Cases vs sessions vs Incidentius | Different lifecycles | Session Stop never closes Case or resolves incident; explicit links and receipts |
| Board models | Existing inputs/mappings; incomplete requiredParams semantics | Reuse input types, specify required-group semantics before use, add runtime binding/run receipts |
| Annotations spec | Project-asset storage; semantic identity can miss environment/tenant | Private OVDB observations with explicit identity domain and provenance; revise owning spec after approval |
| MeaningGraph / ModelSpec | Meaning vs shape; mixed legacy format support | Pin qualified identities/format/revision; no guessed joins or automatic knowledge promotion |
| OVDB remote transactions | Buffered batch lacks read-version assertion | Co-located command transaction or conditional-batch extension; conformance tests required |
| inGitDB default | Good provenance, scan reads and crash proof gap | Benchmark default hypothesis; SQLite engine/index/cache remain distinct options |
| Offline hosted sessions | No inspected browser OVDB/inGitDB runtime | Named cache plus explicit local draft/import package; no fake remote save or audit |
| Single bag original vision | Owner subsequently approved multiple bags/environment gating | Session owns bag collection; applied context set has environment-scoped slots |
| Master checkbox with multiple bags | “Disable all” must remain understandable | Session-wide enabled gate, staged until Apply; individual slots/entries retained |
| Clear action | Must take exactly two clicks | Clear selected bag + confirmation atomically clears active slot/chips; historical refs survive |

For approval, these are proposed amendments, not claims that adjacent feature documents
were already changed. Preserve WIP implementation and coordinate the next integration
against exact current heads. The Detective chat was messaged under the owner's explicit
coordination authorization; no reply is assumed as acceptance of this design.

## Release and rollback

Use additive contracts/capability negotiation. Release the domain model and conformance
fixtures first, then local/hosted repository profiles, then web/terminal consumers. Gate
the new workspace flow until an entire supported profile passes the journey. Never
dual-write old/global and new/session context as competing authorities. Migration imports
are idempotent, retain originals until confirmation, and preserve typed values.

Rollback removes the new navigation/entry points and freezes unsupported new writes;
it retains export/read access to stored revisions. It does not run an old writer over
new reference payloads or delete newly captured investigations. Website availability
flags and demo links roll back with product availability. No automatic cloud relocation,
history rewrite or background sync migration is part of initial release.

Telemetry contains safe event types/count buckets/durations: session created/resumed,
Apply/clear, binding ambiguity resolved, context not applicable, dashboard opened/run,
time-to-first-related-result, persistence error/retry, offline copy/import, late-result
routed correctly. Never raw IDs, labels, pasted text, query values or sensitive titles.
Success metrics: fewer repeated parameter entries, fast first correct result, retained
context after restart, no unnoticed cross-environment binding, repeated session reuse.

## Website and marine-art work package (deferred toward the end)

Owner explicitly requested `datatug.app` landing-page promotion and painted marine
illustrations in the current site's style. Keep this as a dedicated subagent lane near
launch, after useful behavior is demonstrable; do not spend the core-design phase polishing
art further. The concept sheet is a reference, not an approved production asset.

Repository: current `github.com/datatug/websites`, not the older renamed checkout.
Starting files: `apps/app/src/content/features.ts`, `apps/app/src/pages/en/features/[slug].astro`,
feature index and localized homepage/navigation. Discover locale conventions again before
editing. Proposed page `/en/features/context-bag/`; confirm slug and reciprocal links then.

Draft message: **Carry your investigation from query to query.** Enter an order once,
open related data and dashboards, compare named cohorts, and keep UAT/PROD contexts
separate. Explain session versus bags with a brief diagram, explicit Apply, independent
chat attachments and recorded history. Feature card on homepage and feature index;
cross-link query/exploration, Detective and dashboards. Demo CTA opens a supported
session flow; no dead-end standalone feature page.

Marine direction: one tugboat/net/fish for one record; two nets/fish-and-crabs for cohorts;
separate waters for environment boundaries. Match existing hand-painted watercolor/gouache
navy/cream boats, brass details, soft paper edges and warm coastal light. Reference assets
are `pricing-active-users-harbor-960.webp` and `pricing-shared-knowledge-960.webp` in
`apps/app/public/static/illustrations/`. Keep metaphors out of core parameter terminology.
Art captions must say Illustration; images cannot masquerade as product screenshots.
Use real UI only after acceptance. Supply alt text, responsive sizes, image compression,
contrast and motion restraint. Avoid equating UAT/PROD records by matching fish color.

Launch lane acceptance: homepage→feature→demo path works; navigation and all supported
locales are coherent; charts/large-set/offline claims match capability flags; check laptop,
mobile and established locale breakpoints; no collaborative/live-sync promise before
delivery. A separate reviewer checks copy against exact implementation receipts.

## Owner review

Owner approval recorded 2026-10-13 15:46 UTC authorizes the first implementation phase,
M0–M3. The approved real-app direction is top Ionic session tabs (`ion-tabs`, top
`ion-tab-bar`, `ion-tab-button`) and the hybrid Context Bag summary/drawer layout, with
Ionic `ion-item` headers and rows and AG Grid for data grids. PrimeNG is outside the
baseline. Local and hosted OVDB profiles remain in the first slice; M0 chooses the
storage profile from evidence. Website promotion and painted marine illustration are
approved direction, with publication after the delivered behavior can be demonstrated.
M4 is deferred: 5,300-member execution awaits an authorized reference-backed executor;
reviewed extraction awaits an explicit bag-target and source/quote-retention contract;
snapshot branching and direct annotations are outside the first useful one-order
journey. This approval does not decide a sync protocol or authorize publishing, private-
data migration or deployment.

Implementation decisions carried into M0–M3:

1. Validate top Ionic session tabs and the hybrid bag layout in the actual app, including
   tab overflow, close/reopen, mobile editing space and keyboard focus.
2. M0 compares OVDB inGitDB and SQLite capabilities. No sync protocol or engine rewrite is
   preapproved.
3. Use explicit dashboard Run/Refresh and immutable historical runs; Explorer follows Apply.
4. Keep snapshot branching, browser offline draft imports, live shared bags and automatic
   merging out of M0–M3.

The approval authorizes M0–M3 implementation, not unreviewed publishing, storage
migration or deployment of private data. Review findings and their dispositions are recorded in the
[review log](../../research/context-bag/review-log.md).

## Implementation evidence

Snapshot dated 2026-10-14. These receipts document delivered slices and exercised
behavior; they do not declare M0–M3 or the full paired journey complete.

| Area | Evidence | Remaining |
|---|---|---|
| M0 storage decision | SQLite was selected for the local OVDB profile. The inGitDB crash fixture observed inconsistent head/receipt/audit states and remains excluded under [dalgo2ingitdb issue 15](https://github.com/ingitdb/dalgo2ingitdb/issues/15). Actual SQLite and hosted Firestore transaction-owner conformance was exercised. | Keep inGitDB deferred until its crash boundary is repaired and retested. |
| Core contracts | [Core PR 364](https://github.com/datatug/datatug-core/pull/364) (`v0.51.0`), [PR 365](https://github.com/datatug/datatug-core/pull/365) (`v0.51.1`), and [PR 367](https://github.com/datatug/datatug-core/pull/367) (`v0.52.0`) are merged to main; all observed PR and main checks passed. | This establishes the domain/storage contracts, not completion of the app journeys. |
| Hosted API and host | [Backend PR 42](https://github.com/datatug/backend/pull/42), merge `80df77f00def597da860c20454307e4a671726ce`, landed with all ten observed main checks passing. Host PR 1255 is deployed at revision `sneat-app-00582-lkc` with 100% traffic; the activation receipt reports `session_context_enabled=true`. Health returned 200, unauthenticated access returned 401, and the Firestore write probe passed. | No authenticated end-to-end hosted journey has been exercised. Host activation and an unauthenticated 401 do not prove that journey. |
| Browser fixture and lifecycle | Chinook `Invoice.InvoiceId = 1` showed one invoice row in Explorer; its details query returned two invoice-line rows, and the shared-input genre bar reports Rock totaling $1.98. Saved/reloaded state was checked under run `run_bd2f1ba2c84941d5bffb40f7e3b9578e`. Browser checks also exercised two distinct session tabs, a UAT session and second PROD bag, stop/archive/restore/resume, scoped Clear+Confirm that preserved text/history, and explicit resolution of two `outcome-unknown` results without redispatch; the same four run IDs remained and pending count returned to zero. After stopping and restarting the local API over the same SQLite data, the CLI reopened browser-created session `s_f8727a35ec9947c1be43a3f36f457da6` and all four runs; before/after session and run snapshots were deeply equal. Reloading the browser restored both session tabs and the unsent draft “Compare this invoice with the next one.” | The approved invoice fixture and this local cross-surface restart are exercised; broader relationships, the full hosted web/terminal journey, composer execution, and native Detective bridge remain unverified. |
| CLI transport | [CLI PR 374](https://github.com/datatug/datatug-cli/pull/374), merge `f5c9e8c79ca7804623f38f52060f721a7a700eaf`, released as `v0.75.0`; all 18 observed main and release checks passed and the worktree was cleaned. The real CLI journey passed. | The paired authenticated hosted journey remains unverified. |
| Apps and browser delivery | [Apps PR 297](https://github.com/datatug/datatug-apps/pull/297) merged at `be0dc1555444d752f0b856e979d137fb82bc9bd7`. All four required PR checks (build, test, e2e, and real CLI journey) passed for exact head `8f18bb047d1766509608ae46e8d71d9617b3894b`. | Main CI/deployment verification is still in progress, and the authenticated hosted journey remains unverified. |
| Detective bridge | Detective PR 302 provides Case-bound workflows and trusted query receipts, but not a generic session bag bridge. Backend PR 45 merged at `80693b37e183a938ff9f6f03e0cb9b955109fb83` on 2026-10-14 and is published as tag `v0.29.1`; Detective owner is wiring the host to that version, but hosted Detective activation has not been proved. The `detectiveQuerySnapshots` TTL policy is ACTIVE and its ascending, descending, and array-contains `expiresAt` indexes are READY. Device-personal semantic copies require browser execution and suppression of project-only validators as an integration constraint; this behavior is not established as shipped. | Native Context Bag Case/run finalization and the session outbox evidence link remain outstanding. M0–M3 are incomplete. |
| Public explanation and homepage | The [OVDB session storage article](https://datatug.app/en/under-the-hood/ovdb-session-storage/) is live. Websites PR 102 landed at merge `b718302f5145e6cbacbb46dd6879dc55bcb6937f`; all 20 stable main checks passed, including build and deploy. Its task worktree and local/remote branches were later cleaned up. | Canonical sync remains blocked by seven pre-existing modified files, which were preserved; do not claim canonical sync is complete. Homepage publication does not imply the unfinished hosted and Detective flows are available. |
