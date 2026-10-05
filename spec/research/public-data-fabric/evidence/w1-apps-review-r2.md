# Apps PR192 independent review r2

Reviewed-Head: b7b4471c068e72265edf99c5fc2e786fd82258e4
Reviewed-Base: 5520dab90e9920afbea4da6a0f71d4bbd05dedc4
Prior-Reviewed-Head: 8fd4a55fed5abf2903274d4234ea0a723e3a23fa
PR: https://github.com/datatug/datatug-apps/pull/192
Reviewer: /root/w1_apps_independent_review
Runtime/host: local Codex on macOS
Scope: all six repair-delta files plus independent transport, regression, fixture integrity, browser evidence and CI verification; original full staged PR assessment retained below. Read-only; no repository edits, merge or finalization.

## Verdict and landing condition

No unresolved source findings. Both r1 findings are fixed. Source review approves this **staged preparation change only**, conditional on all required CI checks succeeding for exactly `b7b4471c068e72265edf99c5fc2e786fd82258e4` before root lands it. This is not current permission to merge while checks are pending and is not whole-W1, schema2/native execution, production admission or live journey acceptance. Root remains sole landing owner.

At freeze, fresh authoritative GitHub results show build SUCCESS; test, e2e and real datatug-cli journey IN_PROGRESS. Run37259662171, jobs111604034745,111604034514,111604379372,111604379360. Receipt: `/private/tmp/data-fabric-w1-apps-review/r2-ci.json`. No earlier-head green result is substituted for these pending checks. A later failure or changed head invalidates the landing condition until resolved and appropriately reverified.

## Closed findings

### R1-M1 — ROR leading-zero syntax: fixed

`public-data-scenario.ts` now requires literal `0` and six lowercase Crockford characters before the two checksum digits, retaining the MOD97 check. This aligns with [ROR's primary identifier documentation](https://ror.readme.io/docs/identifier), examined in r1; it introduces no new semantic mapping or normalization.

The reviewer reran the same independent request-path probe against newly transpiled, unchanged source files from this exact head and the actual installed `@dalgo/core` parser. Only the HTTP responses are intercepted fixtures. Evidence: `r2/probe.mjs` and `r2/probe-output.txt` under this review directory.

- `https://ror.org/100000096`: validator false; one driver request; **zero target requests**; raw value retained; denominator1, invalid1, unmatched0.
- `https://ror.org/000000098`: validator true; driver then target `organizations.id In [raw]` request; raw value retained; denominator1, invalid0, unmatched1.

The new regression test also exercises both controls through `runFederatedQuery`, not only the standalone validator. The accepted source schema, data and declaration remain byte-identical to apps commit `e7362033ec79c6663d7dbe0483b62fab01f7b9cd`, with their exact pinned hashes verified independently. Receipt: `r2/fixture-hashes.json`. The accepted hub reconciliation scope is unchanged.

### R1-m1 — Product provenance hierarchy: fixed

Independently inspected the source changes, browser receipt `/private/tmp/fabric-w1-apps-r1-browser-receipt.json`, and both actual captured images `/private/tmp/fabric-w1-apps-r1-primary.jpg` and `/private/tmp/fabric-w1-apps-r1-provenance.jpg`. The primary view now identifies the selected source, truthfully says release date/coverage unknown, and shows unavailable lookup with disabled Run. Full hashes and byte budgets are collapsed. Expanding canonical provenance retains all three repositories, immutable revisions, file paths, SHA256 values and 264743/2097152-byte accounting. Selected-field and saved-query provenance also remain inspectable under details; changed-revision acknowledgement stays visible. No field picker was replaced by manual JSON entry.

These are independently inspected author-captured browser images, not a newly repeated live browser session. Their exact-head receipt covers the unavailable state only. No network/live W1 query acceptance is inferred.

## Independent verification

- Remote head/base and OPEN state verified twice through `gh pr view`; local HEAD matches and worktree remains clean.
- Resolved `datatug-main` Nx configuration obtained via `NX_DAEMON=false NX_ISOLATE_PLUGINS=false wb run -- pnpm nx show project datatug-main --json`; receipt `r2-nx-config.json`.
- Executed `NX_DAEMON=false NX_ISOLATE_PLUGINS=false wb run -- pnpm nx run datatug-main:test --testFiles=bounded-federation.spec.ts,public-data-page.component.spec.ts,query-page.component.spec.ts --watch=false`: **66 tests passed across three files, exit0**. Log `r2-tests.log`.
- Reviewed the entire six-file repair diff: regex, bounded regression tests, primary/saved templates, UI disclosure regression and README scope corrections. Independent r1 checks of unchanged discovery/transport code are retained; no broader implementation is newly claimed.
- Original accepted fixture hashes: schema `cfe958f12d9a70bb7533499b76bf01c9447f81f0862852e460b095bb75e1b9cf`; data `44a30dd260c74f43b2138934cca0bccd493c7b804e124c5ecbded3c41226e543`; declaration `ee5dd5d06a9e6db3848e3fc3db71afe2567f6e3e56b7bf8f9477eda5bca54575`.

## Retained F1-F6 and J1-J6 assessment

| Requirement | Staged result | Acceptance still needed |
|---|---|---|
| F1 canonical discovery | Pinned MG/ModelSpec/Directory reader and rebuildable cache; strict bounded admission. Configured provider/entity/property selection uses explicit coordinates. | Landed usable scoped companions and mutation-to-suggestion proof. |
| F2 explained planning | Schema1 exact bridge/property/binding validation and visible unavailable reasons. | Native schema2, fresh exact decision admission, usable source/release/coverage/confidence explanation. |
| F3 execution | Existing worker/DALgo path gains source-first keyed transport with no static fallback; focused fixture tests pass. | Published runtime filter/snapshot/cancel/deadline/pagination acceptance. |
| F4 exceptions | Separate NULL/empty/invalid/unmatched/ambiguous/multiplied counts and fixture ROR status/grain behavior. | ROR syntax defect fixed in r2; prove real global/provider orphan/location behavior and accepted denominators. |
| F5 provenance/reuse | Exact saved coordinates/rights and changed-canonical acknowledgement UI; serialized eligibility cannot execute. | Real store preservation, full fresh revalidation/re-execution, freshness/licence/coverage admission. |
| F6 product evidence | Actual unavailable-state browser receipt plus separate fixture journeys. | Independent integrated real browser/network journey at exact app/provider/runtime pins. |
| J1 connect/inspect | Actual configured source selection observed; exact field selection exercised with fixtures. | Real supported field mapping/provenance at published pins. |
| J2 inspect suggestion | Current production visibly unavailable; fixture scoped explanation. | Admitted canonical compatible suggestion. |
| J3 explicit bounded query | Run is explicit and production disabled; bounded transport fixture evidence. | Actual bounded published requests. |
| J4 inspect exceptions | Template/transport fixture evidence; r2 verifies malformed versus valid-but-absent syntax distinction. | Real data exceptions and attribution. |
| J5 save/reopen | Existing QueriesService shape reused and fixture template/query-page checks. | Real saved scenario revalidation and explicit reuse. |
| J6 no further action | No automatic query/refresh in reviewed path; fixture no-action assertions. | Actual browser network proof with retained snapshot. |

## Remaining implementation and publication prerequisites

- Canonical production admission remains hardcoded false in discoverRepresentations and all saved PublicDataScenario execution is refused in the executor/query page. This is a staged preparation review only.
- The app consumes reviewed schema1 label-bridge metadata only. Native schema2 execution/reader, exact ROR user decision admission and source-faithful native ID/location execution are unimplemented.
- Publish reviewed Directory attachments and immutable provider/model/meaning/runtime companions, then implement fresh admission without trusting serialized eligible flags. A newly landed ROR artifact wrapper alone does not meet this prerequisite.
- Complete product explanations with actual source release/freshness/coverage/confidence and licence/attribution scope, plus stale last-known-good handling, at the admitted exact pins.
- Prove canonical revision mutation changes suggestions after cache invalidation; current independent tests prove cache refetch and exact scope rejection but do not establish the full admitted mutation-to-suggestion journey.
- Prove real configured field selection, actual keyed OVDB requests, native malformed/absent/status/multiplicity behavior, orphan admin1 visibility and the accepted real demo denominators. Fixture tests do not establish these published outcomes.
- Validate exact runtime pagination/snapshot identity, cancellation and deadline behavior plus boundary/over-limit counts/bytes against the published endpoint. Current fixture test evidence is narrower than full live limits acceptance.
- Complete exact saved-plan revalidation and acknowledged re-execution after canonical changes, preserving original result snapshot and all provenance through the real query store.
- Obtain independent actual browser network evidence for J1-J6, including no extra query after no action. B1/runtime/Directory ownership and programme gates remain open.

The README now explicitly distinguishes these app implementation gaps from companion publication prerequisites. Native schema2 helper publication does not implement its app consumer, and the ROR artifact wrapper does not grant canonical decision/runtime admission. Mutation-to-suggestion proof, saved revalidation, freshness/coverage and the independent live W1 journey remain open. No unsupported native receipt normalization, semantic widening or source-data rewrite is authorized by this review.

VERDICT: blockers=0 majors=0 minors=0 land=yes
