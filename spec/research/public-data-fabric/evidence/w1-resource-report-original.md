# W1 combined source resource replay — candidate evidence for independent review

A new cached-input serial ROR + GeoNames replay completed successfully. It reproduces the reviewed native SQLite bytes and all retained source counts. It does not measure the original capture or close full capture+build, runtime, hosting, Directory, browser or B1 gates.

All repository operations were reads. All new files are below this private directory. No source acquisition, service, browser, packaging build, publication or repository edit ran. The six upstream HEAD requests transferred zero source payload bytes; mutable source changes stopped bulk acquisition.

| Metric | New observed result / bound |
|---|---:|
| Complete observed cached operation, including WB, input checks, serial generators and output verification | 132.308685 s |
| ROR generator wall / lifetime peak RSS | 5.557266 s / 47120384 B |
| GeoNames generator wall / lifetime peak RSS | 124.162462 s / 129662976 B |
| Sampled simultaneous WB + worker + generators RSS | 181075968 B |
| Conservative payload RSS upper bound: worker lifetime peak + max(serial generator lifetime peaks) | 164986880 B |
| External observer lifetime peak RSS, separate scope | 26836992 B |
| Sampled path/inode temporary disk, logical / allocated | 1513037456 / 1532022784 B |
| Unique frozen compressed inputs, shared ROR counted once | 670249165 B |
| New bulk network acquisition | 0 B |
| Geo / ROR offline sidecar metadata | 42238 / 160088 B |

RSS lifetime peaks come from macOS wait4 on each generator, which launches no subprocesses. The worker peak comes from getrusage(RUSAGE_SELF). The conservative payload bound is a bound, not a sampled concurrent peak; it excludes WB and observer. The external process-tree samples include WB and all workload descendants and were taken at 0.25-second minimum intervals with actual timestamps retained (478 samples). External observer/ps are excluded from that tree. Disk counts follow input symlinks and deduplicate inodes. Disk is sampled, not quota-enforced; immediately unlinked SQLite temporary files may be absent from the path scan. These limits remain explicit and prevent a strict unsampled aggregate acceptance claim.

Pins: GeoNames source generator `fc53a3537304dd77a4e97162608454b2da2ae241`, importer SHA256 `4808c935cfce0c4b61fb1ab71ebfedcd4a8a5ec279c8ab4230db89233aab0f20`; ROR original source generator `bbbec903248680caea04e68f94b9a957b6efc55b`, parser SHA256 `2ff1bc52de2d0fd10376a26972c45c12800bdadd2a16576d28c23a45b73aaea3`. Source code and generation assets were exported using git archive at those exact commits. Root-supplied ROR canonical artifact pin is `dc78c1e929f1f10018f8c39e690351059a50c71f`; Geo packaging candidate `49d600f5bbddf31900ade5639d760c3dffea20ef` and ROR packaging were not exercised. This receipt makes no fresh CI/landing claim.

Environment: macOS 26.6.2 arm64, Python 3.14.7, SQLite 3.53.4, WB 0.182.0. Exact executable/version strings are in inspection-network.json.

Native data: Geo SQLite 75304960 B, SHA256 `cb53af14c67daadfc20ecbfe304d631d651022179793a77aa11413244d2f456d`; ROR SQLite 238563328 B, SHA256 `afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38`. Ordered Geo gzip chunks independently decoded to that same 75304960 B/hash. SQLite integrity checks passed. All 16 Geo files other than snapshot.json match original r1-a bytes. Original Geo snapshot SHA256 `d564721b80809537424d17b2bf2276697593a3cb380ea020f33c070e1ff0e7ab`; new snapshot `8b2cb98e097ef9f3135179784c27c4b9a54b43f6605020912721510911149f09` changes only generator/importer revision to the canonical pin. Resource receipts contain naturally new measurements and remain separate from native deterministic data.

Counts retained: countries 252; admin1 3865; places 75300; alternate names 844831; ROR organizations 141528 across all active 138289/inactive 1809/withdrawn 1430 statuses; locations 141722; organizations with multiple locations 138; unique ROR GeoNames references 19729; missing ROR places 9; missing admin1 references 234; missing country references 0. Scanned allCountries 13472243 and alternateNamesV2 19219256 rows. No new source hash, narrower scope, inferred identity or grain change was accepted.

Upstream mismatch observed at 2026-10-05T03:47:13Z: allCountries 422005927 B vs frozen 422005409; alternateNamesV2 204812561 vs 204810908; cities5000 5709192 vs 5708626. All five GeoNames Last-Modified timestamps advanced from Oct 4 to Oct 5. ROR archive length remains 37540961 B; a matching HEAD length is not a new hash verification. inspection-network.json retains exact expected hashes and sanitized primary headers; no bulk GET began.

Exact executed commands:

```sh
python3 /private/tmp/data-fabric-w1-combined-resource-validation/measure.py observe
wb run -- /opt/homebrew/opt/python@3.14/bin/python3.14 /private/tmp/data-fabric-w1-combined-resource-validation/measure.py worker
/opt/homebrew/opt/python@3.14/bin/python3.14 /private/tmp/data-fabric-w1-combined-resource-validation/code/ror/scripts/import_ror.py --archive /private/tmp/data-fabric-w1-combined-resource-validation/inputs-ror/v2.13-2026-09-22-ror-data.zip --pin /private/tmp/data-fabric-w1-combined-resource-validation/code/ror/source/ror-v2.13.json --out /private/tmp/data-fabric-w1-combined-resource-validation/ror-output
/opt/homebrew/opt/python@3.14/bin/python3.14 /private/tmp/data-fabric-w1-combined-resource-validation/code/geo/scripts/geonames_w1.py build --inputs /private/tmp/data-fabric-w1-combined-resource-validation/inputs-geo/inputs.json --ror-closure /private/tmp/data-fabric-w1-combined-resource-validation/closure/ror-closure.json --out /private/tmp/data-fabric-w1-combined-resource-validation/geo-output --tool-revision fc53a3537304dd77a4e97162608454b2da2ae241 --receipt /private/tmp/data-fabric-w1-combined-resource-validation/geo-generation-internal.json
```

Independent reproduction: inspect measure.py, exported scripts/assets, frozen input manifest, inspection-network.json, each *-external.json, generator internal receipts, samples.jsonl and worker-result.json. Recompute sample maxima and simultaneous member sums; validate wait4 maxima against internal peak receipts while preserving their different exit timing. Verify original-source hashes, full SQLite hashes/counts, chunk decode identity and the two-field snapshot delta. For a new bounded replay, create a fresh private directory, export the same two code pins/assets and copy measure.py into it; ROOT follows the script path. Execute `python3 <fresh>/measure.py observe` with permission to read process statistics. Existing output paths fail closed. Do not rerun or overwrite this evidence directory.

Recommendation: accept this as newly measured frozen-input combined source-generation evidence after independent review. Keep original full-capture elapsed/RSS and full captured W1 acceptance OPEN. To test a new acquisition without changing the reviewed data, a designated owner can supply authorized immutable hash-addressed URLs for the five frozen GeoNames files plus shared ROR archive, then run an independently reviewed streaming transport wrapper with exact hashes, redirect refusal, aggregate 1 GiB and 1800 s caps followed by the same serial generators. Preserve original source provenance alongside a new explicit mirror transport/fetch receipt and require identical native outputs/counts. Such a test measures a new mirror acquisition, never the historical official capture. Official mutable acquisition instead requires a separately reviewed new monthly candidate; it cannot silently become the old snapshot. If strict unsampled total RSS/disk acceptance is required, add lifetime peak bounds for every wrapper process and quota/open-file disk measurement. No new architecture decision or source promotion was made here.

Independent review is pending with root. This report is an evidence handoff, not a programme acceptance verdict.
