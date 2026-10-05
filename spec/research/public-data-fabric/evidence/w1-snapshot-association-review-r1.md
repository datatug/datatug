# [review r1 ovdb#44] Native snapshot association implementation

Reviewed-Head: 839566d2c9e16d1a5aba439de519298214ccf888
Base: d166a1f8a38a644de9d11f381ad406dad3f89b36
PR: https://github.com/openvaultdb/ovdb/pull/44
Runtime/host: local Codex, macOS. Read-only repository review. Root remains sole landing owner.

## Verdict and scope

No blockers, majors or minors found. Recommend landing this bounded provenance-reader extension after the root owner completes the normal draft/landing transition. The PR is currently OPEN and draft; this is a code-review recommendation, not a merge receipt. No repository changes, merge, WB finalization or provider artifact edits were performed.

The exact remote main and PR head were read with git ls-remote and match the base/head above. The author worktree was clean at the reviewed head before and after tests. The entire handoff and accepted architecture review were read, then the full changed implementation and tests were inspected against their obligations.

This review grants no native Geo semantic acceptance, source-to-native mapping decision, canonical Directory admission, runtime adapter, app support, production eligibility, deployment, B1, or browser journey completion. The hypothetical native Geo fixture input/decision remains an interchange fixture. Existing label-bridge authority is unchanged. Provider bulk reconstruction/key uniqueness is covered by the separate Geo artifact review, not claimed here.

## Independent verification

- Fresh reviewer-owned run: `go test -count=1 -race -covermode=atomic -coverprofile=/private/tmp/data-fabric-w1-native-association-review/full-cover.out ./internal/publisher/... ./publisher/representation ./internal/covergate/... ./cmd/covergate/...` passed. All eight packages reached 100%.
- Reviewer ran `go run ./cmd/covergate` on that new profile with the exact eight CI-gated packages: **2397 of 2397 statements, 8 packages**. This is independent of the author's same-count profile.
- A separate public-API harness, `probes.go`, ran under `go run -race`, with 19 passing checks. It rehashes changed originals, receipts and metadata so semantic rejection is reached rather than merely failing an outer checksum. Cases include 1 versus 1.0, adjacent integers beyond float64 exactness, negative-zero token mismatch, long-s Unicode alias of sha256, null file, absent hash, wrong property/binding/namespace, nonzero duplicates, fractional/overflow counts, accepted exact int64 maximum, unrelated array/scalar preservation, escaped duplicate key, lone surrogate, invalid UTF8, excessive depth, and exact immutable native lookup forwarding. Rejection cases assert the expected failure reason. No dataset, keyset or chunk resolver read is allowed by the harness.
- The main native fixture test explicitly asserts native-identifier execution, preserves the original 7941-byte SHA256 d564721b80809537424d17b2bf2276697593a3cb380ea020f33c070e1ff0e7ab, and retains its five-element chunks array. The complete fixture is 42906 metadata bytes with eight resolver reads. The reviewer independently recalculated this total.
- Both actual native-a and native-b packages' country receipt, original snapshot, later metadata snapshot, model and binding match the committed fixture byte-for-byte. Only those small metadata files were read; no SQLite, compressed chunk or key corpus was loaded.
- Both contract schemas, default publisher packages, workflow and parity artifacts have zero diff from the base. Schema hashes independently checked: schema.json c75eb5cae7201a9c10432239ece5ffd93a750f191c9b04129c89169f4449d844; schema2.json fa281400fbe063d1abbcbf5c7578dd23932dd798ee91810e0906ffa2431dc0bc. Fresh tests preserve format1, format2 bridge, real ROR native absent-association, and real Geo label-bridge behavior.
- Exact-head CI run 37259206683 has authoritative conclusion SUCCESS. All nine applicable jobs succeeded: lint, build/test, web/build/embed/browser, Linux/macOS/Windows lifecycle, tidy, exact coverage, parity. Version bump is intentionally disabled and SKIPPED, matching unchanged workflow configuration. Saved ci.json and pr.json record the result.

Initial sandboxed Go test startup could not access the shared Go build cache. The approved rerun completed successfully; the failed startup produced no test result and is not counted as verification.

## Contract and negative obligation audit

Presence of the closed association selects one literal bounded ASCII output key. Absence preserves filename association; null/malformed explicit association rejects, including a test against otherwise valid ROR evidence. Source has only path/hash, is resolved in immutable provider context, and must match exact later metadata artifact closure. Unsafe/external/mutable/unresolved references, wrong hashes, byte excess and direct receipt/later-metadata/dataset cycles reject.

The original is strict parsed before use. Duplicate decoded keys, invalid Unicode, excessive nesting and malformed roots reject. Numeric tokens are preserved through json.Number for parsed equality, so neither float64 rounding nor 1-to-1.0 normalization can manufacture original/embedded equality. The initial existing strict receipt parse conservatively rejects overflow float64 values; this does not admit mismatches. Only the literal selected descriptor supplies file/hash, so unrelated chunks arrays survive and matching decoys cannot substitute. Consumed protocol aliases reject; entity names remain literal dynamic keys.

The selected file and hash must match the unchanged native dataset. Entity count must exist and be a nonnegative int64 integer token equal to receipt records. Receipt model, binding, entity, property, namespace and dataset must match contract scope; claimed duplicates must be zero. Missing/wrong selectors, selected arrays, missing/wrong file/hash, fractional/exponent/negative/overflow counts, mismatches, closure omission, malformed sources and no-fallback cases are covered in the implementation tests and reviewed alongside the independent probes.

LookupNative preserves exact repository/revision/model/snapshot/dataset and Limit 2. Metadata validation never resolves the native dataset. Resolver identity, regular-file and symlink/redirect behavior remain the pre-existing explicit Context resolver obligation; this change does not add a transport implementation.

## Evidence pins

- Helper: 839566d2c9e16d1a5aba439de519298214ccf888.
- Reviewed original Geo source provider: fc53a3537304dd77a4e97162608454b2da2ae241 (architecture reference, not freshly fetched by this reviewer).
- Fixture packaging generator: 7c0223df8e59b80df8ddf71e775197a153cf3e30, explicitly a fixture dependency pin, not a claim of published generated artifacts at that commit.
- Model SHA256: 609fe73e07b84db65bdb1a904b9f84dd5c43bcbbc75eeb2d75f9a69ed3e8af81.
- Binding SHA256: 21bb35aecf1f1214908f245cfeb2d0c5e7c6caf3915a5ffd604254642b3b4b57.
- Country proof SHA256: 4f52bf1f9bbb417fc3ca1a67b0fc96e0af9d8baf271fd66cfe9eaaaa5b8e4412.
- Later metadata SHA256: 5bb94bc67a326da337bd1a8be978dd5a384095e60d47b2669afc579ad86d396e.
- Dataset SHA256 cited by the verified metadata: cb53af14c67daadfc20ecbfe304d631d651022179793a77aa11413244d2f456d; bulk bytes were not read here.
- MeaningGraph core fixture pin: 982916d73f0a35ff2558b0062f58aa3ac4f24d97.

## Findings

None. No severity/reproduction entry is required because no defect was found. This conclusion would be falsified by a rehashed invalid original/descriptor/closure admitted through Check, a regression of the absent-association route, or a bulk resolver read during the explicit native fixture check. The reviewed tests and probes did not produce those outcomes.

VERDICT: blockers=0 majors=0 minors=0 land=yes
