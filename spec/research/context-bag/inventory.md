# Context Bag repository evidence

2026-10-13. Source inspection, not deployment certification. Local paths are under
`/Users/alex/projects`. Fresh-fetch rows are distinguished from cached references and
active WIP. The full source vision is the user-provided
`/Users/alex/projects/datatug-context-bag-orchestrator-vision.md`; later conversation
steering adds multi-bag/environment behavior, dashboards, Clear and public-site work.

## Main and existing components

| Repository / inspected revision | Exact path and line | Finding |
|---|---|---|
| datatug/datatug, fetched c679a6e71061219eede31d02219ffb3ca7db7f24 | spec/features/investigation-context/README.md | Prior scope/conditions/layers contract; immediate disable conflicts with new Apply |
| datatug/datatug-apps, pulled a35c58e9da555eab9b1ed1cab102bfd14844180b | libs/datatug/semantic/src/lib/services/investigation-context.service.ts:28,219,260,415 | Scope lacks session; value-derived identity; sessionStorage; immediate enabled changes |
| Same | libs/datatug/semantic/src/lib/services/binding-resolver.ts:1,65 | Explicit binding/no silent rebind, unsupported shapes reported |
| Same | libs/datatug/main/src/lib/chat/chat-session.service.ts:11,33,127 | DALgo IndexedDB sessions/turns/results; multiple sessions already exist; no start/stop/archive state |
| Same | libs/datatug/main/src/lib/chat/investigation/investigation-storage.ts:27 | Separate investigation IndexedDB; deliberate migration isolation |
| Same | libs/datatug/main/src/lib/pages/signed-in/env-db-table/env-db-table.page.ts:648 | Table load passes source/table/limit without context bindings |
| Same | libs/datatug/main/src/lib/queries/query/page/query-page.component.html:625,1379 | Parameters/Sources/Related sidebar and execution-details accordion |
| Same | libs/datatug/main/src/lib/queries/query/page/query-page.component.ts:848,1478,1825 | Parameter recovery; Use as filter changes draft query with Undo; save state distinct |
| Same | libs/datatug/main/src/lib/queries/query/page/query-workspace-layout.component.html:11 | Stacked editor/results and mobile Jump/Back navigation |
| Same | libs/datatug/board-models/src/lib/board.ts:23; widget-parameter.ts:22; widget-sql.ts:18 | Existing board input definitions and widget query mappings |
| datatug/datatug-core, fetched 7ea9a517ff131c8b5c52ce3571847954e048385c | pkg/apicontract/typed_value_set.go:10 | Inline set capped at 500: 5,300 reference path is new work |
| Same | pkg/detective/repository.go:37; pkg/investigation/typed_value.go:15 | Optimistic typed repository/receipt/audit pattern; lossless scalar encoding |
| Same | pkg/semantic/README.md:3,53 | Existing semantic Resolve/RelatedLookups/Applicable and composite relations |
| Same | pkg/datatug/boards.go:24,56,265 | Board Parameters/RequiredParams and query-input mapping; nested required-group semantics still TODO |
| datatug/datatug-cli, fetched 24b3e431e74f228145c18f2a3a2dde63fb0801d3 | pkg/chat/workspace.go:13,75; pkg/chat/session_chat.go:41 | Terminal references/attachments and durable session switching restrictions |
| Same | pkg/executionstore/store.go:1,57,283 | Private SQLite evidence snapshots with limits; not proof of proposed session cache |
| Same | apps/datatugapp/datatugui/dtproject/project.go:43 | Dashboard list placeholder |
| Dashboardius fetched 1270dce21b9eb7fbf0bd1d1366f0f3d65e196bc4 | web/src/app/dashboard/model/board.model.ts:26; ui/cards/query-widget.component.ts:129 | Separate UI uses embedded data, not established generic context-execution pipeline |
| Same | spec/features/query-backed-board-persistence/README.md:24 | Describes scripted in-memory prototype; no production dashboard claim |

Apps was initially inspected at 52c6e94e, then refreshed on the owner's request. The clean
canonical main was already a35c58e9; `git pull --ff-only origin main` returned Already up
to date. The added query-author changes were explicitly re-read and reconciled.

## Operational storage and semantic dependencies

These are inspected cached refs unless noted; sandbox fetch failures are access limits,
not evidence that upstream lacks newer work. Revalidate pinned versions before implementation.

| Repository / revision | Evidence | Design consequence |
|---|---|---|
| openvaultdb/openvaultdb-go 25185bdb75a9a40274e69559297fe02d40a05e67 | pkg/manifest/manifest.go:24,61,149; pkg/server/server.go:377 | Logical DB and multiple engines; batch/query routes exist |
| Same | pkg/core/core.go:532,651 | Batch/push boundary; publication can fail after persistence |
| ingitdb/dalgo2ingitdb, OVDB pin v0.6.2 / 609f28ce20db250c93f50463938434e81c431da1 | database.go:343; git_commit.go:84; query.go:19,53 | Writer lock, in-memory rollback, expected Git head; scan query; crash recovery still needs proof |
| dal-go/dalgo2openvaultdb 4c79be8b98ccca3c929e89b96b4a79dd95bac572 | tx.go:60,230; updates.go:23; query.go:80 | Buffered remote batch; preconditions/cursors unsupported |
| dal-go/dalgo2ovdb-js 0580a5ef3d689c669b1af5613b7a501e1f6e24ed | src/database.ts:118,167 | Similar HTTP reads/batch; not a remote transactional readset |
| ingitdb/dalgo2ingitdb4github 2e32527ba5a79db40c93a3cfbb11c2cd40d5f507 | batching.go:15,76; tree_writer.go:78 | One Git commit per buffered batch; no earlier-read assertion |
| ingitdb/ingitdb-ts f931267f840d164a61d9dba67f285d065fd1a73c | packages/client-github/src/cache/idb-cache.ts:1; changes/pending-changes.ts:147 | Browser cache exists; file-by-file commitAll isn't atomic domain transaction |
| dal-go/dalgo-http-adapters app pin 519665be7bab21bae3c388eb2395afb8b6df1be9 | packages/indexeddb/src/query.ts:187 | Existing browser adapter scans then filters; large-history index is a separate capability |
| github.com/meaninggraph/core 25242a6c84d1b8d4799c718ef237944a429835ba | FORMAT.md:21,70; commerce.meaning.yaml | Meaning/qualified references, draft-reader caveat; order and invoice semantics distinct |
| specscore/modelspec local 3c149d20c5a0779cf668bb646042362c33577c82 | README.md:19,44,68 | Shape uses record/field; meaning belongs to MeaningGraph; not permission policy |
| modelspec-org/registry 6e954755ee79b758213911dc2d97687808ed0f84 | scripts/fixtures/chinookdb/model/chinook.modelspec.json | Key arrays and references support composite representation |
| datatug/backstage fetched abc212d23b0b46ea7fcf9ecce6df3addfaca25be | docs/vision-and-inventory.md:32 | Jev decision narrowing candidate; named emulator isn't proof of TypeSafe provider integration |

Do not confuse `ingitdb/dalgo2ingitdb4local` legacy adapter (several direct-method panics)
with the `dalgo2ingitdb` version pinned by OVDB. Likewise `openvaultdb/ovdb` CLI/publisher,
`openvaultdb/openvaultdb` specifications and `openvaultdb/openvaultdb-go` runtime are
different repositories. The proposal uses runtime evidence, not names alone.

## Detective active work and coordination

Inspected clean local hub WIP HEAD `23a43cb6ebac70ba51e903c193a91bf66f6babc8` at
`.worktrees/detective-checks/github.com/datatug/datatug/spec/features/detective-checks/README.md`:
lines 24–28 and 100–130 require local+cloud and hosted execution; 158 defines Case
lifecycle; 269–288 defines private audit and transaction/recovery; 357 defines core/adapter
ownership. It is active adjacent work, not a blanket completion receipt.

Apps WIP `c01fdc740bbd9fa3d6d2c737704589a54d012029`:
`libs/datatug/main/src/lib/detective/models.ts:112` stores Case typed context;
`cases-page.component.html:53` says Case evidence spans sessions;
`cases-page.component.ts:292,591` reviews imports from current basket.
CLI WIP `ab4d8effd9ff2b6142d92932ac6fb727a57f4915` has private store paths in
`pkg/detective/store.go:38`. Backend WIP at `41052a3a5cbbf02f68f02e47d1230059fca3d116`
included untracked `facade4datatug/detective_action.go:20` transaction work when inspected.
These files may change while their owner works; no modifications were made to them.

Under the user's explicit authorization, coordination messages were sent to
**Detective and checks in UI** about session lifecycle, applied revisions and OVDB.
The proposal preserves that lane's implementation authority. No inferred agreement or
claimed delivery is based on merely sending those messages.

## Website location and visual sources

Current local website repository is `github.com/datatug/websites`, origin
`git@github.com:datatug/websites.git`, inspected HEAD
`73729f399b809924d0191f698e7b0ce89fd2e5dc`. Feature routes and content live at
`apps/app/src/pages/en/features/[slug].astro`, `apps/app/src/content/features.ts`.
Painted site references actually viewed: `apps/app/public/static/illustrations/`
`pricing-active-users-harbor-960.webp` and `pricing-shared-knowledge-960.webp`, plus
`apps/app/public/static/tugboat-coast.png`. Direct live-site fetch failed in this tool;
the style match is grounded in these repository assets, not a claim of live visual QA.

## Validation boundary

This discovery used source reads, fresh targeted fetches where available and isolated
design artifacts. No product implementation, source query, private data migration or
deployment was performed. Prototype uses synthetic data and memory-only state. Its
JavaScript behavior can be checked locally; rendered browser review was blocked by the
in-app browser file-URL policy. No alternate browser/server workaround was attempted.


## Late integration refresh

While this proposal was being prepared, the Detective feature specification landed in
DataTug hub main at `feb88bd9ddb8ac9f4dfcb94d38d77a43571cd2a0` (PR162). This was
verified by the canonical repository commit/tree, not inferred from a lane report. Its
feature remains Implementing; specification landing is not product completion. The
Context Bag worktree began at c679a6e. SpecScore integration validation used an exported
feb88bd tree overlaid with these new artifacts and merged index rows: lint returned `[]`
and exit 0. No active Detective worktree was modified.
