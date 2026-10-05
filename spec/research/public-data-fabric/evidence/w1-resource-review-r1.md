# Independent resource review r1

**Conditional acceptance: cached-input combined generation is supported after one minor wording correction. Full capture/build acceptance remains OPEN.** No generation rerun is needed for the corrected scope. No repository edits, source acquisition, services, browsers or publication were performed by this review.

## Finding R1-M1 — minor: RSS lifetime wording exceeds the capture point

`measure.py:126` captures worker `ru_maxrss` after verification, before result construction, JSON serialization, file writing, final stdout and process exit. `measure.py:182` captures observer RSS before its own serialization/exit. Calling these full-lifetime maxima is unsupported.

Correct wording: **worker peak through post-verification capture point: 35,323,904 B; conservative payload RSS bound through that point: 164,986,880 B, excluding later receipt serialization/stdout/teardown and WB/observer. Observer peak through receipt assembly: 26,836,992 B.** Generator wait4 lifetime maxima remain valid. Preserve raw receipts and add a scope annotation rather than rewriting historical evidence. A wording repair suffices for cached-generation acceptance; a new externally reaped run is needed only to claim full worker/observer lifetimes.

## Independently executed checks

The WB-governed audit script exited 0. It independently recomputed source and artifact SHA256s, streamed all Geo source rows, hashed the decompressed ROR member, queried read-only SQLite integrity and tables, recalculated status/missing-reference counts, decoded Geo gzip chunks, compared original Geo files and snapshot deltas, and recalculated every sample's member sum and all maxima. No author measurement functions were imported or executed.

| Check | Recomputed result |
|---|---:|
| Samples | 478 |
| Actual sample interval min / median / max | 0.267174 / 0.276944 / 0.319108 s |
| Simultaneous sampled WB + worker + generator RSS | 181,075,968 B |
| Sampled logical / allocated disk | 1,513,037,456 / 1,532,022,784 B |
| Six unique source inodes | 670,249,165 B |
| Geo / ROR sidecars | 42,238 / 160,088 B |
| Geo SQLite | 75,304,960 B |
| ROR SQLite | 238,563,328 B |

The retained observer elapsed receipt is 132.308684583 s; sample timestamps span 0.021884083 to 132.047946666 s and corroborate that duration. Historical monotonic timings and wait4 values cannot be independently remeasured without another run. ROR and Geo phase receipts are consistent with internal measurements and nonoverlapping sampled PID presence. The generation code is sequential and launches no subprocesses. The numerical payload-bound calculation is correct within the corrected scope.

Geo SQLite and decoded chunks: `cb53af14c67daadfc20ecbfe304d631d651022179793a77aa11413244d2f456d`. ROR SQLite: `afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38`. Both pass `PRAGMA integrity_check`. All 16 non-snapshot Geo files match original r1-a; snapshot changes are exactly the two revision fields.

Counts match: 252 countries, 3,865 admin1 rows, 75,300 places, 844,831 alternate names; ROR 141,528 organizations, 141,722 locations, 138 multi-location organizations, 19,729 distinct GeoNames references; active/inactive/withdrawn 138,289/1,809/1,430. Independently recomputed missing-reference counts are 9 ROR places, 234 admin1 references and 0 countries. Full source line counts agree, including 13,472,243 allCountries rows and 19,219,256 alternateNamesV2 rows.

Live read-only `ls-remote` confirmed Geo main `6f4cf1269bc393048f6b204069135a62f0bb6c02` and ROR main `dc78c1e929f1f10018f8c39e690351059a50c71f`. Source pins `fc53a3537304dd77a4e97162608454b2da2ae241` and `bbbec903248680caea04e68f94b9a957b6efc55b` are ancestors; all eight exported code/asset files equal both their source-pin Git objects and current canonical-head objects. Packaging was outside this measurement.

## Limits retained

- Original official capture elapsed time and RSS remain absent and cannot be recovered by this replay.
- No acquisition was included. Zero bulk GET is supported by the cached command path and retained author receipt, not independently captured historical network traffic. The six HEAD responses are retained author evidence; this review did not repeat them.
- Sampled process RSS excludes observer/ps and may miss transients. The process list is an OS sample, not an atomic all-process physical-memory measurement.
- Disk peaks cover paths beneath the private measurement root plus symlinked inputs, inode deduplicated. They exclude unlinked SQLite files and external prior comparison fixtures. They are sampled observations, not complete filesystem high-water or quota enforcement.
- All generation code/assets match both pinned source commits and the current canonical heads, but packaging was not run in the measured operation. No CI, deployment, Directory, browser, runtime concurrency or B1 conclusion is made.
- Reproduction relies on original absolute cache paths, original ROR closure and original r1-a comparison fixtures hard-coded in measure.py. A fresh directory alone is insufficient on a machine lacking those retained inputs.
- Independent source review checks byte hashes, Geo row counts, SQL projection counts/statuses and missing references. It does not redo semantic reconciliation or licensing decisions.

## Smallest responsible next operational validation

1. Freeze the reviewed source manifest and require a designated owner to identify an authorized immutable transport location for all six exact input objects. There is no executable network acquisition of the frozen GeoNames bytes available in the supplied evidence today; do not manufacture a mirror URL or use changed mutable bytes.

2. Implement and independently review a small streaming acquisition controller using that manifest. Refuse redirects, stream to fresh files while hashing, enforce an aggregate 1073741824-byte cap and one 1800-second deadline covering capture, input verification, serial ROR then Geo generation, and output verification. Capture ROR once. Require exact lengths/hashes before generation; fail closed and retain diagnostic receipts on any mismatch.

3. Run that controller once through WB in a fresh private directory. Record acquisition URLs/timestamps and measured network payload separately from original official provenance. Reuse the exact pinned generation commands in the existing worker receipt and require the exact native SQLite hashes, Geo decoded-chunk hash, retained SQL counts/statuses and missing-reference counts.

4. If strict RSS/disk ceilings are part of operational acceptance, measure the controller and all child/wrapper lifetimes externally (or use an explicit enforceable memory ceiling) and use a bounded filesystem/quota that counts open unlinked SQLite temporary allocations. External instrumentation must itself have explicit scope; a supervisor exiting after serialization cannot prove its own lifetime peak.

5. Label the result new immutable-mirror acquisition plus build. It cannot repair the missing historical official-capture measurement. If the requirement instead mandates official upstream capture, separately review the changed monthly candidate and its provenance/semantics before acquisition/promotion; do not silently repin this snapshot.

This is a proposed validation design with an unresolved immutable-transport prerequisite, not a tested runbook or an already executable acquisition command.


## Selected next test: separately pinned official operational candidate

The orchestrator selected this option. Before any new download:

1. Repair cached-evidence RSS capture-point wording and retain original raw receipts plus their hashes.
2. Create a distinct candidate ID/output root and manifest; lock generator code pins, official URL allowlist and the existing ROR release pin. Accepted artifact/publication pins remain unchanged.
3. Independently review the acquisition/measurement controller before download: redirect refusal, streaming hash/length validation, response/timeout handling, one aggregate 1 GiB transferred-payload limit and one 1800-second end-to-end deadline, fail-closed cleanup/retained diagnostics, serial generation, external process accounting, explicit disk scope and handling of open unlinked SQLite files.
4. Agree the prospective operational acceptance scope: what RSS/disk maxima are bounded or merely sampled, whether instrumentation is included, per-output 256 MiB SQLite and 2 MiB metadata ceilings, and 512 MiB RSS / 8 GiB disk limits. An unmeasured or sampled-only component must remain disclosed rather than silently counted as proven.
5. Do not pre-approve unknown new Geo hashes as accepted corpus data. Capture only to the separate candidate, freeze observed hashes, then perform independent source-delta/projection/resource review before operational acceptance; publication promotion remains a later gate.

Acquire current official bytes once into a distinct immutable candidate manifest, then run the same pinned serial generators within one measured envelope. This can satisfy a prospective operational capture/build budget gate without a frozen mirror. It does not demonstrate the accepted historical snapshot capture.

Before acquisition, explicitly authorize the new candidate and source scope. After capture, freeze exact source hashes/lengths/timestamps and independently review changed rows/counts, ROR closure, missing references, schema and projection invariants, candidate hashes and resource receipts. Never require different source bytes to reproduce the old output hashes; establish and review new candidate outputs while retaining old artifacts unchanged.

Resource acceptance may be limited to the reviewed same-scope candidate operation. No publication or republishing all dependent artifacts is required merely to measure that operation. A major source-size/schema/scope change must not be treated as representative without review.

Only later promotion of this candidate into the corpus requires regenerating and revalidating every source-dependent artifact, provenance/manifest hash, registry/Directory dependency and deployed journey that changes. Existing published artifacts stay on their original pins until that promotion passes.

Missing original capture telemetry remains an evidence limitation, not necessarily a permanent plan gate. The orchestrator may close a prospective reproducible-operation requirement with this explicitly reviewed new operation; it must not claim retrospective proof for the old exact bytes.

## Exact evidence binding

Every read file, exported asset, source symlink target and artifact SHA256 is bound in `review-r1.json:evidence_sha256`; independent machine results are in `audit-results.json`, reproduced by `audit.py`. Primary reviewed evidence hashes:

- `report.md`: `86305a41eec6bfba4e8d06e84f9420cab6699389f65f7aff3362911b82baa8c6`
- `summary.json`: `eb96a8cd3395bc24e2d8cc558bd2aeb3ff307139b3cd7bc472b056fab31449df`
- `measure.py`: `d5b4703109c4b01dd66b29b21bca862718c4e14168ed968f804ab624c9e01120`
- `inspection-network.json`: `302e16eeeddf78eac1b2c232c979d9b05a1122223498fc834415169e92b32cc9`
- `observer-result.json`: `ffe93673705d42fbaa4f6f3fcf65d18443849923ba45423021df2051f0fdc62d`
- `worker-result.json`: `a7da360d4bfe23377414bf380569cd6f212d4c909650d4a7fab8315438ebf6d6`
- `samples.jsonl`: `b102e6c788840f7b14919878665880c2e666875e39c6d97776b6c48c09dd3119`

VERDICT: blockers=0 majors=0 minors=1 land=no
