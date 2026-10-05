# Independent resource scope correction review r2

Reviewer: w1_canonical_evidence, local Codex/macOS, distinct from the resource author and original r1 reviewer.
Scope: additive report/summary correction only. No generation, source acquisition, bulk-data audit, service, browser, deployment or gate closure was performed.

Accept the corrected evidence for its cached-input combined-generation scope. R1-M1 is resolved; no new findings. This verdict relies on the retained independent r1 native/source/resource audit for unchanged measurements and verifies only the correction and original evidence identity.

## Exact reviewed artifacts

- report-r2.md SHA256 `46d1ec9efc9b2b003b1d57fb4b3e9201c1ffeb26012c4ccf00f9fd2d649dc7a2`.
- summary-r2.json SHA256 `356312583066e7270f23f154fa2910fd030fb0d74363dbd02b59235dc81e5098`.
- Original independent review-r1.md SHA256 `5993ff0238599946e73b6f866ebaf2313b2a76917acaced9d3fea5dc904a9fba`; review-r1.json SHA256 `eb8c380c469d01e6f1034ad8d8246287ba36f513625481bd9f8d1dfe58139bde`.

Read the full original review, original and corrected report/summary, and actual measure.py capture sites. Recomputed all six original hashes named in summary-r2.preserved_original_sha256; each matches both retained bytes and the original review evidence map. Also hash-checked the 44 original receipt/code/metadata files at most2MiB in that map, including inspection-network.json and generator phase receipts. Large native/source/chunk byte hashes inherit the original independent audit; they were not redundantly recomputed in this narrow review.

## Resolved correction

At measure.py126 the worker captures RSS after verification and before result construction, serialization, stdout and exit. At182 the observer captures RSS during receipt assembly before its own serialization/exit. Both corrected documents now describe those exact capture points. Only externally reaped generator wait4 values are called full-lifetime maxima. The corrected payload bound excludes later worker serialization/stdout/teardown and WB/observer. The unchanged measurements are worker35,323,904B, bound164,986,880B and observer26,836,992B.

A structural comparison confirms every metric value is unchanged after exactly three capture-scope key renames. Other top-level summary changes are only measurement-limit annotations, explicit reproduction prerequisites, original hash preservation, review binding and revision label. Source pins, native hashes/counts, timings, sampled disk/process values, network observations and open gates are unchanged. The report delta is the matching scope correction plus review history and reproduction caveat.

The corrected report and summary explicitly require retained absolute cached-input paths, the original ROR closure and r1-a comparison fixtures. A fresh directory alone is not claimed sufficient on another host. Original raw receipts and measure.py are unchanged, including their historical overbroad wording; the additive corrected evidence supplies the narrower interpretation instead of rewriting measurement history.

## Accepted limits

Acceptance covers frozen cached inputs through serial full ROR/Geo generation and retained output verification under disclosed measurements. It does not establish the original official capture elapsed/RSS, complete acquisition+build, unsampled all-process RSS, unlinked-file disk high-water, packaging, runtime concurrency, hosting, Directory, browser or B1 acceptance. The separately pinned official candidate remains prospective and needs controller review before bulk acquisition and independent candidate/source/resource review afterward. Original published data and pins remain unchanged; no semantic decision or eligibility is granted.

VERDICT: blockers=0 majors=0 minors=0 land=yes
