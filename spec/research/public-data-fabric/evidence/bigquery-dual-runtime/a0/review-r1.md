# Independent BigQuery A0 review — round 1

Architecture acceptable: **0 blockers, 0 majors, 2 minors**. Root should adopt or explicitly decline the following small clarifications before the immutable corpus is frozen. Existing ownership/provisioning/release prerequisites remain binding. No implementation or live acceptance is claimed.

## Exact reviewed inputs

- `proposal-r1.md`: `08e088b6d89016d40acfcd669df2158f83033d89a7a2bbdabe3e4aff1357d7b6`
- `proposal-r1.json`: `87010fa5d0a239847baf72af8fb76a846b4e13e854ceb6a720b808a5e981f379`
- `hash-index.json`: `7e4fb164e8c4100421001f589e054013a5e9902b230152d498abb599d534bd38`

Canonical r3: `datatug/datatug@0d71452a89d967cb0b10b931fe45af92aa9a3dff`. Reviewer runtime/host: local Codex/macOS.

## Findings

### m1 — minor: Freeze the canonical encoder and explicit digest payloads

Locations: `proposal-r1.md:53`, `proposal-r1.md:56`, `proposal-r1.md:60`, `proposal-r1.md:65`, `proposal-r1.md:104`.

Sorted keys and no whitespace do not uniquely specify UTF-8 JSON bytes. Go default encoding escapes HTML and U+2028/U+2029; JavaScript JSON.stringify does not make the same choices. Round-trippable FLOAT64 text is also not unique. ReadPlan contains its own digest, and the contract needs an explicit hash projection rather than hashing the self-containing object. The mandated corpus can close this before implementation; this is not a reason to add a core interface.

Exact proposed amendment: Use RFC 8785 JCS for the adapter-owned canonical JSON payload: valid Unicode scalar strings, no Unicode normalization, UTF-16 code-unit property order, its exact string escaping and UTF-8 bytes. Reject unpaired surrogates. All numeric JSON tokens in these payloads are bounded exact safe integers; warehouse integer/decimal values remain canonical strings. Finite FLOAT64 strings use ECMAScript NumberToString (including its exponent thresholds), normalizing negative zero to "0", with the same formatter semantics in Go. Define named hash payloads: ReadPlan hashes exactly version/sourceDigest/projection/where/order/limit/parameters/sql, omitting only its top-level digest; SourceProfile hashes all declared fields; the observation binding hashes only exact table identity/location/native type/config/schema; approval hashes the declared binding fields (effective plan, observation and policy digests, execution principal including generation, project/location, cap/budget/bounds and estimate), excluding its own digest and the explicitly unbound times/nonce. Do not recursively delete properties named digest from arbitrary values. Publish these exact payload schemas and canonical bytes in the frozen corpus.

Conformance cases:

- Go and JS produce identical bytes and SHA-256 for strings containing quotes, backslashes, control characters, <>&, U+2028/U+2029 and supplementary characters.
- Non-BMP property ordering follows UTF-16; lone surrogate is rejected; composed/decomposed Unicode strings remain distinct.
- FLOAT64 cases -0, 1e-7, 1e-6, 1e20 and 1e21 have identical strings and digests.
- Changing a top-level digest field alone leaves its named underlying digest calculation unchanged; changing a bound parameter or cap changes approval; excluded observation/preview timestamps do not.

### m2 — minor: Make the run deadline persistent across Resume

Locations: `proposal-r1.md:58`, `proposal-r1.md:61`, `proposal-r1.md:63`, `proposal-r1.md:77`, `proposal-r1.md:79`.

The contract persists counters and bounds but does not name when wallMs starts or whether paused/process-exit time counts. A new Resume context could unintentionally restart 120 seconds. Status and cancel need their stated recovery role without giving paging a fresh execution budget.

Exact proposed amendment: At the start of run reauthorization/pre-submit preparation, persist runStartedAt and executionDeadline = runStartedAt + wallMs in the trusted ledger. The execution deadline is absolute; stopped/offline time counts. All run metadata/dry-run/recheck/submission/poll/page operations, including Resume, use the earlier of caller deadline, per-HTTP timeout and this original executionDeadline. Resume never changes the deadline or resets row/page/byte counters. Once expired, result reads return local_stopped with a deadline reason and the existing receipt without network dispatch. Explicit status/cancel may use a fresh control context of at most 15 seconds after that deadline, but cannot read rows, rerun a query, renew approval, reset cumulative byte counters or release a reservation without authoritative terminal billing evidence. They retain bounded responses and existing cumulative byte accounting; if those bytes are exhausted, return the bound error and retain the reservation. Preview has its own bounded operation and the existing five-minute approval TTL; time awaiting user approval is not charged to a run that has not started.

Conformance cases:

- With a fake clock, stop at 119 seconds and resume at 121: zero page/POST requests, original deadline and counters unchanged, partial receipt retained.
- Resume before deadline permits only the remaining time; GET retries and backoff consume that time.
- Status/cancel after execution deadline use fresh <=15-second control contexts; they never reset execution limits or submit queries.
- Exhausted cumulative bytes reject further control responses safely and retain unresolved reservation; completed status with valid billing reconciles only the existing job.

## Verified conclusions

- The concrete `datatug query bigquery` group is a viable real consumer. The actual existing query group has one registration point; its saved/ad-hoc keyed routes justify keeping the selected bounded route separate. No `dal.DB`, core interface or generic paid leaf scan is needed.
- Released DALgo recordset authorization rewrites residuals and explicit projection before calling the underlying recordset method; the private sentinel propagates correctly. Whole-query guards before and after capture, plus reauthorization on run/resume, are appropriately scoped.
- Generated SDK methods and exact cap/useLegacySql behavior compile. The nil-collapse problem is real; raw bounded validation remains mandatory. Additional tests confirm `NullFields`/`ForceSendFields` encode null, empty scalar and empty array exactly. A custom in-memory transport returns one attempted submission on ambiguous failure. This is compile/protocol evidence, not production decoder acceptance.
- Native TABLE/configuration/location rechecks, full source/schema binding and explicit residual replacement race meet M2. A closed metadata classification may reject a new Google field until classified; this is the selected safe limitation, not a reason to broaden A0.
- Exact string values, SQL NULL/JSON-null separation, schema-before-rows, partial-page offset/digest and same-job cursors have no identified architecture defect. Runtime fixtures must still prove these contracts.
- Cap reservations, persistent private ledger, one-shot approval and no automatic submission replay are coherent. Reconnection is explicit and requires the same verified execution subject. Browser same-token UserInfo and separate app identity preserve m1; Go user ADC cannot acquire missing grants through SDK options, and trusted workload providers need no invented OIDC sub.
- Source connection tests remain distinct from admission; W1 core pin, existing source/semantic gates and ECB FX are preserved. Real CLI and local server acceptance plus deployed browser/two-page gates remain mandatory.
- Go module provisioning and JS release-file coordination are explicit preconditions. Fresh WB output remains incomplete and shows Northwind and W1 overlaps. It does not authorize exclusive ownership.
- Scope is proportionate. A narrow private capture session avoids core work. The shared raw corpus and raw validation duplicate some decoding intentionally to preserve evidence lost by generated models. No extra registry, shared spender, source reconciliation, raw SQL or general query engine is needed.

## Verification receipts

- `hash-verification.json`: all 45 indexed files match.
- `compile-rerun.log`: original five tests passed fresh using Go 1.27.1.
- `independent-tests.log`: seven tests passed, including two additional independently authored tests.
- Tests used in-memory transports and cached released modules with network module fetching disabled. Default build-cache access was sandbox-denied; a private cache succeeded.
- Exact canonical r3 and CLI blobs were checked against their recorded commits; released module source was read directly. JS source/release snapshots were checked at recorded revisions.
- `wb-active-dal-go.json` and `wb-active-datatug.json`: each command returned exit 1, explicitly incomplete. No repository mutation/ownership claim followed.
- No jobs, credentials, OAuth setup, source edits, servers or background processes remain owned by this review.

Primary references: [jobs.query](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/jobs/query), [getQueryResults](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/jobs/getQueryResults), [QueryParameter](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/QueryParameter), [Google OIDC](https://developers.google.com/identity/openid-connect/openid-connect), [JCS](https://www.rfc-editor.org/rfc/rfc8785.html).

VERDICT: blockers=0 majors=0 minors=2 land=yes
