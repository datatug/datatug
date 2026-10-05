# Review r1 — BigQuery Go and JavaScript architecture

Reviewed proposal r2 at `/private/tmp/data-fabric-bigquery-browser-architecture/proposal-r2.md`, SHA-256 `9a073ff10388c3fb0f6a6d27c38bbfe9a79b4bf595dad61a4e5cc7e4dad18dd0`, and its JSON companion. Local Codex, macOS, 2026-10-05. Report-only authority; root owns acceptance and landing. No repository changes, credentials, cloud changes, queries or builds. W1 and ECB FX remain intact. This review makes no dataset admission or semantic-equivalence decisions.

## Required changes

### M1 — Validate raw cells before generated Go decoding erases missing-value evidence

**Major. Proposal lines 46, 123–131, 135–148.** The promised distinction between an omitted/malformed wire cell and SQL NULL cannot be recovered from the selected client's decoded model alone. At the pinned generated-client revision, `TableCell.V` is an `interface{}` with no presence field (`google-bigquery-gen.go:10219–10231`); query responses are decoded into generated models at lines 14049–14060. Both `{}` and `{"v":null}` leave `V == nil`. Treating either as a nullable SQL value silently admits malformed responses. Testing only normalized result fixtures would miss this.

**Minimum amendment:** require bounded raw JSON shape/presence/type validation before generated result decoding, or an adapter-owned lossless response decoder. Retain the generated client for supported request/auth operations. Include actual HTTP response fixtures for omitted `v`, explicit `v:null`, wrong cell count, wrong wire value type and valid SQL NULL/JSON null. Run those bytes through each runtime's real decoding path; do not start only from normalized objects. Invalid wire cells must fail without being exported as NULL. No DALgo core change is needed.

Evidence: [pinned generated TableCell](https://github.com/googleapis/google-api-go-client/blob/267f7a8c8e539f21e64b696686bb9476bf0c5ef0/bigquery/v2/bigquery-gen.go#L10219), [query decoding](https://github.com/googleapis/google-api-go-client/blob/267f7a8c8e539f21e64b696686bb9476bf0c5ef0/bigquery/v2/bigquery-gen.go#L14049).

### M2 — Make the forbidden external-execution boundary transitive through source tables

**Major. Proposal lines 30–34, 66, 78, 139–140.** The proposal excludes external connections, remote functions and external-cost pathways, but the descriptor/preflight contract never names permitted BigQuery table kinds or dependency checks. Restricting the emitted statement to `SELECT <columns> FROM <reviewed identifier>` does not establish that boundary: BigQuery table resources include SQL views and external tables. A logical view can hide SQL calls/dependencies absent from the emitted query. Google documents additional Cloud Run costs for remote functions. A valid schema, location and dry-run estimate therefore do not establish the promised execution profile.

**Minimum amendment:** define the first slice as reviewed native BigQuery tables and reject unsupported table kinds/configurations before dry run or execution. In particular, reject VIEW, EXTERNAL and unreviewed MATERIALIZED_VIEW/BigLake/external-storage configurations rather than attempting dependency support now. Record source kind/configuration with metadata observation and invalidate approval if the observed execution-relevant metadata changes. Include a native-table positive fixture and view/external-configuration negative fixtures in both runtimes, asserting no query submission. Future dependency analysis can remain deferred. This is a transport eligibility guard, not a dataset-rights or semantic decision.

**Exact initial contract:** `datasets.get` must establish the selected dataset location; `tables.get` must match the expected fully qualified table identity, observed schema and admitted native `TABLE` profile. Reject external-data, view, materialized-view and BigLake configurations even if another field looks acceptable; unknown types/configurations fail closed. `SNAPSHOT` and other kinds remain unsupported until explicitly added. Recheck this execution-relevant observation before submission, bind it into the approved plan, and invalidate approval on relevant changes. Test kind/configuration/schema/location changes between preview and submission. This preflight is not atomic with `jobs.query`: an upstream owner can replace a source after the check. Record that residual race truthfully and restrict the first slice to reviewed trusted publishers; do not claim an immutable source or an adversarial-source security guarantee from an etag check. A hard atomic guarantee would require a separately designed source-pinning mechanism, outside this bounded slice.

Evidence: [BigQuery table resource types and external configuration](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/tables), [remote-function execution and additional costs](https://docs.cloud.google.com/bigquery/docs/remote-functions). The hidden-dependency consequence is an inference from those documented resource and SQL-execution contracts.

### m1 — Name the execution-account identity grant and lookup

**Minor, required A0 clarification before implementing account matching. Proposal lines 9, 54–56, 154.** The fail-closed instruction is correct, but the normal successful path is unspecified: GIS `TokenResponse` has no `sub` or email, and repeating consent with only `bigquery.readonly` does not supply either. Firebase identity and account hints cannot fill the gap.

**Minimum amendment:** request the identity scopes needed by the chosen identity lookup alongside BigQuery read-only consent; use `openid` for subject identity and add `email` when an email display label is wanted, then call Google's UserInfo endpoint with the same access token. Bind the execution principal to the returned stable `sub`; email is only an optional display label, not a linkage key or a substitute for a missing subject. Keep this separate from Firebase account creation/linking. Revalidate after token/account changes and add same-account, different-account, partial-grant, missing-subject and missing-email cases (missing email must not invalidate a verified subject). Use the Google-owned identity endpoint discovered from its fixed trusted discovery document; public source metadata must not supply an identity URL. Verify the exact GIS/UserInfo flow at the deployed-browser gate. No new human product decision is needed.

Evidence: [GIS TokenResponse fields](https://developers.google.com/identity/oauth2/web/reference/js-reference), [Google identity scopes, UserInfo and stable subject](https://developers.google.com/identity/openid-connect/openid-connect).

## Accepted architecture

- Generated REST `jobs.query` is a justified default for the required read-only scope. Official REST docs accept `bigquery.readonly`; `jobs.insert` and cancellation do not. Pinned high-level Go SDK source confirms `Query.Read` can fall back to `Query.Run`/insert. The proposal correctly distinguishes IAM, OAuth, Stop waiting and actual cancellation. [Query scopes](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/jobs/query), [insert scopes](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/jobs/insert), [cancel semantics/scopes](https://docs.cloud.google.com/bigquery/docs/reference/rest/v2/jobs/cancel).
- Exact DALgo Go source supports keyless `QueryExecutor`/`RecordsetReader` and source `DbType`; the direct executor plus full-query guard avoids generic fallback without requiring a core-interface change. Optional keyed access remains properly gated.
- Reusing the maintained JS monorepo, adapting its obsolete peer/imports, testing both exact core revisions and preserving the W1 app pin is sound. This accepts the migration design, not current package compatibility.
- Partial-page cursor offset or boundary-only cursor issuance, full job identity, omitted-schema reuse, no ambiguous resubmission, decimal strings, temporal precision, shared immutable fixture revision and both-runtime parity are appropriate. M1 strengthens the fixture entry point.
- Canonical OVDB source-access extension is honestly described as required new schema work; it does not misrepresent an existing mount engine or immutable proof. Independent source-specific admission remains intact.
- Required real Go consumer selection before coding, local server harness, deployed-browser and Go live journeys, and two-page receipts are explicit. The root's pending concrete consumer selection is an A0 prerequisite, not a reason to invent a broader server product or block this bounded architecture.

## Evidence and limits

Both evidence indexes' referenced local snapshot hashes verified without mismatches. Reviewed pinned DALgo contracts, generated/high-level Go SDK sources, current/pinned JS package contracts, current adapter code and closed OVDB schema/mount evidence. Independently read DALgo recordset rows/columns from the exact Git object. Rechecked the cited Google documentation live. No CI, deployment, real OAuth/CORS, live source availability, source licensing or runtime success is asserted. Missing implementation/live proof is already a delivery gate and is not counted again as an architecture defect.

Accept the direction; amend these concrete A0 contracts before landing the proposal as implementation-ready. No other blockers or major defects found in the scoped review.

VERDICT: blockers=0 majors=2 minors=1 land=no
