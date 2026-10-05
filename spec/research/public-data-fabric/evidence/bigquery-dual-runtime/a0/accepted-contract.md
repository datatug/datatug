# Accepted BigQuery A0 contract

Root accepted A0 architecture with the two mandatory amendments below. The exact original proposal, independent review and root decision remain immutable. This amendment document overrides only the corresponding ambiguous proposal-r1 wording; canonical r3 remains binding. Acceptance covers architecture and implementation scope only. Both Go and JavaScript, the actual CLI and local server, the deployed-browser live gate, source-access profile and specialist rights/admission gates remain mandatory. No implementation, provisioning, source admission or runtime success is recorded.

The selected real Go consumer is `datatug query bigquery`; released DALgo v0.89.6 and Google API SDK v0.296.0 compile with Go 1.27. The proposed driver home is `github.com/dal-go/dalgo2bigquery`; JavaScript reuses `dalgo-http-adapters/packages/bigquery` as `@dalgo/bigquery`. Provisioning and shared release-file ownership remain prerequisites. JCS applies only to adapter-owned digest payloads; the raw BigQuery lossless bounded parser remains required.

## m1 — Freeze the canonical encoder and explicit digest payloads

Use RFC 8785 JCS for the adapter-owned canonical JSON payload: valid Unicode scalar strings, no Unicode normalization, UTF-16 code-unit property order, its exact string escaping and UTF-8 bytes. Reject unpaired surrogates. All numeric JSON tokens in these payloads are bounded exact safe integers; warehouse integer/decimal values remain canonical strings. Finite FLOAT64 strings use ECMAScript NumberToString (including its exponent thresholds), normalizing negative zero to "0", with the same formatter semantics in Go. Define named hash payloads: ReadPlan hashes exactly version/sourceDigest/projection/where/order/limit/parameters/sql, omitting only its top-level digest; SourceProfile hashes all declared fields; the observation binding hashes only exact table identity/location/native type/config/schema; approval hashes the declared binding fields (effective plan, observation and policy digests, execution principal including generation, project/location, cap/budget/bounds and estimate), excluding its own digest and the explicitly unbound times/nonce. Do not recursively delete properties named digest from arbitrary values. Publish these exact payload schemas and canonical bytes in the frozen corpus.

Required conformance cases in both runtimes:

- Go and JS produce identical bytes and SHA-256 for strings containing quotes, backslashes, control characters, <>&, U+2028/U+2029 and supplementary characters.
- Non-BMP property ordering follows UTF-16; lone surrogate is rejected; composed/decomposed Unicode strings remain distinct.
- FLOAT64 cases -0, 1e-7, 1e-6, 1e20 and 1e21 have identical strings and digests.
- Changing a top-level digest field alone leaves its named underlying digest calculation unchanged; changing a bound parameter or cap changes approval; excluded observation/preview timestamps do not.

Disposition: adopted verbatim; required in canonical contract and both runtime conformance corpus.

## m2 — Make the run deadline persistent across Resume

At the start of run reauthorization/pre-submit preparation, persist runStartedAt and executionDeadline = runStartedAt + wallMs in the trusted ledger. The execution deadline is absolute; stopped/offline time counts. All run metadata/dry-run/recheck/submission/poll/page operations, including Resume, use the earlier of caller deadline, per-HTTP timeout and this original executionDeadline. Resume never changes the deadline or resets row/page/byte counters. Once expired, result reads return local_stopped with a deadline reason and the existing receipt without network dispatch. Explicit status/cancel may use a fresh control context of at most 15 seconds after that deadline, but cannot read rows, rerun a query, renew approval, reset cumulative byte counters or release a reservation without authoritative terminal billing evidence. They retain bounded responses and existing cumulative byte accounting; if those bytes are exhausted, return the bound error and retain the reservation. Preview has its own bounded operation and the existing five-minute approval TTL; time awaiting user approval is not charged to a run that has not started.

Required conformance cases in both runtimes:

- With a fake clock, stop at 119 seconds and resume at 121: zero page/POST requests, original deadline and counters unchanged, partial receipt retained.
- Resume before deadline permits only the remaining time; GET retries and backoff consume that time.
- Status/cancel after execution deadline use fresh <=15-second control contexts; they never reset execution limits or submit queries.
- Exhausted cumulative bytes reject further control responses safely and retain unresolved reservation; completed status with valid billing reconciles only the existing job.

Disposition: adopted verbatim; required in canonical contract and both runtime conformance corpus.

## Review and evidence scope

Independent review r1 found zero blockers, zero majors and two minors; both are adopted verbatim above. Compile evidence proves the released SDK calls, serialization, ambiguous single submission, nil collapse and DALgo policy capture in memory. It is not adapter or live acceptance. Original index entries outside this curated folder identify exact historical source blobs and private ownership snapshots; those large snapshots are not republished. `source-index.json` retains immutable repository revisions/object paths. The local author rehashed all 45 original indexed files before copying this focused evidence. The packet verifier checks all curated bytes and the review/root input bindings.

Implementation plan: [existing Workbench DataTug namespace](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/public-data-fabric-bigquery/README.md). Root remains sole landing owner. The plan must pass fresh independent review before landing and preserves W1 GeoNames/ROR, ECB FX and the complete candidate backlog.
