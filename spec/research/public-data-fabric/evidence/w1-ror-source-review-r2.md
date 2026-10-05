# [review r2 ror-ingitdb#1]

Reviewed-Head: f5f6c0fbc54d792d1062dbdf1cbce1b690f0dd95
Reviewed-Base: 156fa46e53b04f4d0b8aa54e31cca2d86af63038
Prior-Reviewed-Head: 056741bf4600862685fb5081ebfcbc58b112e6d9
PR: https://github.com/ingitdb/ror-ingitdb/pull/1
Reviewer: independent local Codex/macOS /root/w1_ror_review, 2026-10-05 UTC.

No remaining source-code/model acceptance findings. Both r1 majors are resolved; no source edits, merge or deployment performed.

## Resolution evidence

- R1: recursive duplicate-key hook, explicit non-JSON constant rejection, finite-float overflow guard and `allow_nan=False` now reject the original duplicate-status, NaN and1e400 fixtures. Tests also exercise nested duplicates, Infinity/-Infinity, negative overflow, nested overflow, serializer refusal and clean stage removal after a valid prefix. This is bounded source JSON validation, not a general normalization framework.
- R2: organization and relationship IDs require exact native HTTPS ROR URL shape and checksum; projected scalar/container types are checked before SQLite insertion. Original empty URL, query suffix and relationship numeric/bool coercion fixtures now fail. The checksum follows the publisher generator at https://github.com/ror-community/ror-api/blob/bd040a0d2558a478c06a89118a29eeb9b6142710/rorapi/management/commands/generaterorid.py, independently read. Positive controls retain a valid absent target, withdrawn status, null location ID, repeated GeoNames locations/ordinals, empty label, exact whitespace/Unicode, finite floats and an integer above2^53. No failed output or stage remains.
- Independent fixtures/scripts and machine-readable results live under `r2/`: probes.py, negative-probes.json, positive.py, positive.json. All six prior malformed-input reproductions reject for the intended validation reason.

## Current verification

GitHub PR and `git ls-remote` independently confirm exact head and base above; OPEN/non-draft. Actions run37252007403/job111581395461 (pull_request) and run37252004367/job111581386495 (push) independently report status completed/conclusion success/head_sha equal Reviewed-Head.

Ten repository tests pass. Pinned ModelSpec publish lint, HCL/JSON twin, MeaningGraph checks and Git whitespace checks pass. Tool build provenance independently rechecked: ModelSpeca8de98b4565ae669b6fba00c334b4649c9d40b11 and MeaningGrapha58a506d9e6f76660cc8586478521787976111d0, both vcs.modified=false. Model, meaning, licence and publication-boundary files are unchanged from r1.

Independent full-source audit repeats against the corrected author's bundle, using review-owned formatted-record framing without importer code. All full raw objects/statuses, locations/details/ordinals and relationships/types/targets/labels agree; no extra rows. Counts remain141528organizations/141722locations/72747relationships,138organizations with multiple locations and19729distinct GeoNamesIDs; active138289/inactive1809/withdrawn1430. Archive/member hashes and publisherMD5 are unchanged and independently recomputed. SQLite integrity and declared FK checks pass.

Both corrected author bundles and a fresh reviewer full rebuild have identical SQLite, snapshot, attribution and ID-list hashes. SQLite238563328B SHA256afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38 remains identical to r1; snapshot SHA25680f845f8c652408fc7e5e61bf17f09ffe75bed5087f7b3d7657696ae972d3e96 correctly reflects new parser SHA2562ff1bc52de2d0fd10376a26972c45c12800bdadd2a16576d28c23a45b73aaea3. Committed source/validation.json snapshot exactly matches actual corrected snapshot. Reviewer rebuild measured28.519911seconds/61440000B peakRSS/276264145B archive+bundle temporary bytes on Python3.14.7/SQLite3.53.4, within the per-build ceilings. Timing differs from author runs and is not a deterministic claim. An early verification attempt before this asynchronous build completed found no output; after successful completion all artifact comparisons passed.

## Scope and unresolved dependencies

**Source-only landability accepted.** The source, model, meaning and attribution artifacts can be landed independently after root's supported lifecycle checks. This does not authorize user-ROR eligibility or public deployment. The exact user schema/namespace needs the separately assigned fresh reconciler. Native-id runtime adapter, immutable chunks/hosting, registries/Directory, GeoNames exception closure, combined W1 budgets, B1OPEN4/5/6 and real browser journey remain separate prerequisites. OVDB.md publish:[] and template-only manifest preserve that boundary.

Root reported a separate WB lifecycle problem: the author prematurely finalized its work-log, preventing normal report refresh against a terminal projection. This reviewer did not inspect or repair that lifecycle and makes no completion claim for it. Root retains lifecycle resolution and sole landing ownership. It is not an unresolved code finding and does not imply permission to bypass WB.

VERDICT: blockers=0 majors=0 minors=0 land=yes
