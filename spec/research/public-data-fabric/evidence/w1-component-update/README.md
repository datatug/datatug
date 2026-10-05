# W1 component publication update — 2026-10-05

This additive snapshot starts at hub `d5b58603d4142330e2dcdd504efa5d16a3dbbbd6`. The [machine record](evidence.json), [normalized receipt snapshot](publication-receipts.json) and [hash manifest](source-manifest.json) retain public coordinates and scoped results from independently reviewed components. Original private reports and machine envelopes remain private. The earlier [follow-up packet](../w1-followup/README.md), accepted decisions, candidate indexes, questions and gates remain unchanged.

Six components have root-verified exact reviewed/landed tree equality, canonical sync and source branch/worktree cleanup:

| Component | PR | Landed main | Exact-main CI at freeze |
| --- | --- | --- | --- |
| OVDB exact source input | [45](https://github.com/openvaultdb/ovdb/pull/45) | `e64f8adcabaf1d64f844757dd2d27e917011731a` | Passed; OVDB release v0.27.0 points here |
| Chinook offline precheck | [4](https://github.com/demo-db/chinook/pull/4) | `1d15e4ae9b3323db9b575df35ab4b97fdeaff6cf` | Passed |
| Declared demo source | [6](https://github.com/datatug/chinook-demo/pull/6) | `4ed08a43e166b41c2f2cde0d3665eef78a1e8d86` | Passed |
| Apps exact source input | [194](https://github.com/datatug/datatug-apps/pull/194) | `413bb772e04f139210adec92ad0cb3fce2b07883` | Passed |
| Chinook licence grammar | [5](https://github.com/demo-db/chinook/pull/5) | `39e1830b1ad1f190e62cbfa9f3eec96ae28e92dc` | Passed |
| OVDB literal records reader | [46](https://github.com/openvaultdb/ovdb/pull/46) | `27f2664782e01feb4fb5938f28b8f8eae0a868a8` | Pending; release pending |

Apps r1 found M1: source checksum verification occurred after decoding. The repair verifies bounded raw bytes before fatal UTF8 decoding, strict JSON parsing and page selection. Independent r2 passed 107 targeted tests and eight ordering probes with zero remaining findings; the author reported 2,253 tests. A fresh root receipt verifies all five exact-main checks passed; the earlier pending observation and failed WB receipt-mismatch attempt remain in the normalized receipt history. The exact-main [Cloudflare deployment workflow](https://github.com/datatug/datatug-apps/actions/runs/37274522297) succeeded. That workflow does not establish live W1 query or browser acceptance.

Declared-demo r1 independently exercised released DataTug CLI v0.55.0 at `8f408b3b6158ccac808771fd78152f6b6ff08d8f`, validated the project and matched the unchanged real hero golden: 412 invoices and 24 countries. It verified the accepted eight-row source's schema/data and immutable catalog identity. This closes the stale CLI-release preparation blocker for that configuration. It grants no changed-input semantics, public-target admission or browser acceptance.

Chinook precheck r1's five major and one minor findings are retained in the machine record with their r2 resolutions: metadata root shape, control-escape decoding, order-independent duplicate source identity, publisher meaning-file association, repository dot-component refusal and bounded malformed-row findings. Independent r2 reran original integration probes and 464 crossproduct assertions. Independent r3 verified three additive schema checksum entries and unchanged preexisting entries; it did not rerun full corpus generation. Local precheck strength remains partial: external source bytes and canonical admission are unresolved.

The separately reviewed licence architecture and root technical approval permit only `licences.data` conjunctions of two through four distinct known IDs, separated by literal ` AND `, at most 64 ASCII bytes, using the existing 18 Publisher atoms with no normalization. Independent architecture r1 had one minor clarification, resolved by root: new operand rejection rules apply only to conjunctions. Every existing single-token profile remains unchanged, including Directory shaped unknown IDs, plus/LicenseRef/case behavior; Publisher single-token allowlist and model/meaning rules remain unchanged. This is a bounded technical grammar decision, without blanket legal compatibility or provider admission. Chinook's independently reviewed implementation passed eight focused tests and 1,748 probes. Directory/Go parity remains pending.

The source rights remain scoped: [ROR release](https://zenodo.org/records/22902037), [GeoNames source notice](https://download.geonames.org/export/dump/readme.txt), [CC BY attribution conditions](https://creativecommons.org/licenses/by/4.0/legalcode.en#s3), and [SPDX conjunction syntax](https://spdx.github.io/spdx-spec/v2.3/SPDX-license-expressions/#d43-conjunctive-and-operator). Existing hashed `DATA-LICENSE.md` must be retained, and actual app/runtime/download attribution delivery requires separate verification before opt-in.

The OVDB records-reader repair separates literal `$records` path enumeration from repository/ref grammar enforcement. Independent review passed focused, race, schema and real committed provider metadata probes. Its production harness reproduced failure with OVDB v0.27.0 and success with the reviewed fix, exercising four contracts and eight negative controls. Source-data execution was outside that harness scope. The reviewed provider candidate was not admitted by this reader test.

Geo PR8 at `28b8295d3ebe15a7d81ce6d05c2f1a1997eff463` and ROR PR4 at `b531fd8458226b402f8e08170777b7aa4da99f90` remain candidates pending fresh dedicated semantic/packaging reviews. Accepted provider pins and original semantic decisions are unchanged. The external publication owner reports authorization and six passing API checks; root has not verified that report in this snapshot. B1 remains open pending publication, cleanup and live review. All five implementation records remain `implementing` with empty publication/query evidence arrays; phase 4 is active, phases 5/6 and the real J1–J6 browser journey have not passed. W1 remains incomplete.

Run `python3 spec/research/public-data-fabric/evidence/w1-component-update/verify.py` from the repository root. It rehashes packet files and preserved research groups, checks hash declarations and component boundaries, and scans this public packet for private paths. It does not rehash private original documents, fetch source objects, perform reconciliation or admit execution. Receipt outcomes are dated observations; empty/pending attempts are not promoted to successful CI or publication.

## Open Questions

Authoritative questions remain in `fabric.questions`. Exact-main CI/release, Directory/Go licence parity, scoped attribution delivery, fresh provider admission, B1 and the live J1–J6 product journey retain their existing owners and acceptance gates.
