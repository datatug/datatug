# Design review and validation record

2026-10-13. The owner approved the design and first-phase implementation in the
conversation. The design reviews below do not establish release, deployment or
completed production browser testing.

## Independent review rounds

| Round | Lens | Findings | Disposition |
|---|---|---|---|
| Product r1/r2 | UX + investigator | Explorer-only terminology; unclear Stop/Resume; global Apply could be blocked by one consumer; composer revision policy; missing dashboard AC; overlap universe | Fixed in product specification and feature requirements |
| Architecture r1 | Offline/storage | Browser copy looked authoritative without OVDB; mutable scope could reinterpret history; dashboard requirement absent; engine wording | Fixed: local draft package/idempotent import; immutable bag scope and pinned manifests; formal dashboard AC; hypothesis wording |
| Prototype fix r1/r2 | Interaction consistency | View all/cohort mode bypassed Apply; global state leaked between sessions; composer lost drafts; no archive restore; Clear absent | Fixed and exercised with state assertions |
| Multiple-bag reconciliation | Product + architecture | One bag no longer enough; environment identity; global disable ambiguity | Owner approved multiple bags/environment gating; per-scope manifests plus explicit global master |
| Architecture r2 | Integrated four documents | No blocker/major; empty-clear wording and contextSet table inconsistency | Both minor findings fixed |
| Product r3 | Real-world investigator + source review | PROD dashboard invalidated by unrelated UAT changes; branches cloned entry IDs | Fixed: compare own scope manifest; new child IDs with origin refs |
| Product r4 | Independent recheck | No remaining substantive contradictions | VERDICT: blockers=0 majors=0 minors=0 land=yes |

Product and storage specialists cross-reviewed applied revisions, composer independence,
offline authority, lifecycle, deletion and environment comparison. Principal designer
reconciled their recommendations rather than concatenating reports. No finding was
silently waived. No review verdict claims the store/executor capabilities are already built.

## Checks performed

- Principal source inspection confirmed session scope lacks session ID, existing board
  input/mapping fields, 500-inline-set limit, Explorer load contract, and inGitDB lock/
  rollback implementation. Apps main pulled on request at a35c58e9 and latest Query Author
  layout explicitly reconciled.
- Prototype author ran syntax and state/source assertions through four fix rounds.
- Principal independently ran [verify-prototype.cjs](verify-prototype.cjs): Apply gating,
  123→4 related fixture rows, no pre-Apply change, 5,300-ID immutable chat receipt,
  UAT Apply leaving PROD intact, selected-bag Clear, composer retention and stopped
  execution passed. This uses a minimal DOM stub; it does not verify browser pixels.
- `git diff --check` passed for tracked design edits.
- SpecScore lint passed (`[]`, exit 0) in an integration validation tree consisting of
  latest DataTug main feb88bd plus the design artifacts and merged index entries. The
  original design worktree base c679a6e predates the concurrently landed Detective
  Feature; its routed plan reference consequently needed the latest Feature for validation.

## Validation limits and deferred work

The in-app browser refused local `file:` navigation under its URL policy. No alternate
browser, server or indirect workaround was attempted. Rendered layout/mobile/focus QA
and investigator usability sessions remain unperformed. The HTML clearly labels itself
as a synthetic, memory-only prototype. Source tests are useful evidence of its state
model, not end-to-end production proof.

OVDB durability/concurrency, 5,300-member real executor semantics, local/hosted parity,
browser offline imports and source authorization require the proposed M0/M1 acceptance
work after overall design approval. Art/landing refinement is explicitly deferred near
launch per the owner; the painted concept sheet is only a reference.

## Publication status

The approved design was preserved in the `context-bag-phase-one` WB task and
published for review in [DataTug PR #163](https://github.com/datatug/datatug/pull/163).
The Feature is Implementing; its routed first-phase Plan is published in
[Workbench PR #69](https://github.com/sneat-co/workbench/pull/69).
The owner approved Ionic session tabs and the hybrid bag layout. First-phase
implementation includes web and CLI paths with local and hosted storage.

Implementation and its separate code review are in progress. The original
prototype evidence above remains prototype evidence. It does not establish
application execution, hosted durability or delivery. Landing-page promotion and
painted marine illustrations remain deferred near launch as requested.
