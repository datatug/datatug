---
format: https://specscore.md/feature-specification
status: Implementing
---

# Feature: Hosted Google Cloud log evidence

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/hosted-google-cloud-log-evidence?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/hosted-google-cloud-log-evidence?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/hosted-google-cloud-log-evidence?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/hosted-google-cloud-log-evidence?op=request-change) |
**Status:** Implementing
**Source Ideas:** —

## Summary

Give a signed-in DataTug user a personal, direct-read connection to one Google Cloud Logging log view.

**Owner-approved direction (2026-10-19):** The user separately consents to the narrow `logging.read` OAuth API scope and selects one project-local **log view**. DataTug uses that user's grant for bounded, direct `entries.list` queries through `dalgo2logs`. A server-side DataTug gateway authenticates and authorizes every request. There is no log ingestion, database snapshot, shared credential, aggregation system, or public OVDB protocol change. Workspace sharing needs a later authorization decision.

This feature remains in implementation. The direct adapter, server contracts,
encrypted grant registry, bounded reader and feature-gated browser UI are
merged. The local Northwind demo reads committed synthetic files through a
loopback OVDB server and displays results in AG Grid; it is not a hosted Google
user-grant pilot. Host route composition, identity checks, quotas, task-based
cleanup and the pilot UID allowlist are merged in `datatug-cloud` PR #12.
Task OIDC authenticates cleanup; the UID allowlist applies only to verified
user management/query routes. Generic secret bindings and paused OIDC Scheduler
deployment support are merged in `sneat-co/cicd` PR #44; its header-removal
follow-up PR #45 is also merged. Host PR #25 merged at
`9d055b12fae5ae74793e8f452bd7fc8714e22d33`; exact-main CI passed 10/10 at
`13:00:44Z`, and deployment `39426215088` succeeded. Provider readback at
`13:07:48Z` found Belgium production `datatug-cloud-00010-bjt` and preview
`datatug-cloud-preview-00009-2fx` at 100% on the merged image. Production
references OAuth Secret Manager version `2` and state-key version `1`, with
master feature and user routes both `1`; preview and Frankfurt remain without
pilot configuration. All four health checks returned HTTP 200.
Pre-activation checks on 2026-10-20 found all 35 retained DataTug revisions
across Belgium and Frankfurt and all 7 current project Cloud Run services
without personal Cloud Logs settings. All 9 credential roots were empty at
`11:36:36Z`. The deployed Firestore ruleset matched reviewed origin main
`a78abcd7e55258b93ba996b7bfe1e48149cfe029` byte for byte (SHA-256
`9038e32d45753b5ebe62ac078743da5b0cd7eadb40b9dfedae723c8bf245d8d2`).
This supports the pre-activation old-writer check, not blanket shutdown proof.
The private OIDC control returned authenticated status and rejected unauthenticated
status; at the earlier checkpoint, pause at expected epoch `0` persisted as paused
epoch `1`. The later approved resume to epoch `2` is recorded below. The owner
approved the existing recurring cleanup, which was enabled at
`12:21:37.755516Z`. Its scheduled target returned HTTP 200 at
`12:22:00.796489Z` on revision `datatug-cloud-00007-lt8`; the OIDC target was
correct and `X-Cron-Secret` was absent. The earlier single approved target
request returned HTTP 200 at `12:15:20.352921Z` on revision
`datatug-cloud-00006-fc2`. Later scheduled target requests returned HTTP 200
at `13:06:08Z`, `13:07:07Z` and `13:08:04Z` on production revision
`datatug-cloud-00010-bjt`; the OIDC target was correct and
`X-Cron-Secret` was absent. These responses do not establish which records
were deleted. An approved CAS resume from epoch `1` to `2` returned HTTP 200 at
`13:09:15.798Z`; durable readback confirmed epoch `2`, `paused=false`. The
disposable Scheduler probe was deleted and returned `NOT_FOUND` at
`13:10:25Z`. All four
Belgium/Frankfurt production/preview health checks returned HTTP 200 with the
expected service header; preview and both Frankfurt services had no personal
Cloud Logs settings. The
restricted test identity
has a view-only binding on the synthetic log view, but positive view-read and
negative broader-project-read checks remain pending. No Cloud Logging user
grant exists. The original Google OAuth client
credential is retired; this did not revoke a Google user token. PR #25's
Firebase `auth_time` fix is deployed. No Google Cloud Logging user grant
exists. Current receipts and
remaining gates are in the
[pilot guide](https://github.com/datatug/datatug-demo-project/blob/main/docs/logs/google-cloud-pilot.md).
[Apps PR #407](https://github.com/datatug/datatug-apps/pull/407) merged
at `a63cfe2883d1078f9a3c1058483197126ea372de`; its production UI flag is on
behind the fail-closed, server-confirmed UID access gate, and Business Preview
remains off. The post-merge [CI run 39292521711](https://github.com/datatug/datatug-apps/actions/runs/39292521711)
passed preflight, build, test, Cloud Logs E2E and deploy-gate checks.
Cloudflare workflow [39293867514](https://github.com/datatug/datatug-apps/actions/runs/39293867514)
succeeded, and public build-info reports SHA
`a63cfe2883d1078f9a3c1058483197126ea372de` with version `0.1.0+1201`. No
post-deploy CI smoke was run. Fresh Google basic sign-in reached the
authenticated Cloud Logs connections page, which showed no accounts or views.
Authenticated GET connections returned HTTP 200 at `13:09:41.557Z` and
`13:09:41.993Z`; POST OAuth start returned HTTP 200 at `13:10:08.380Z` on
production `00010-bjt`. The Logging chooser did not recognize the selected
provider account, so no consent, grant or query was completed. The same-origin
`createAuthUri` probe returned HTTP 400 `INVALID_CONTINUE_URI`; its callback-
domain ownership issue remains unresolved.

### Delivery status (2026-10-19)

| Slice | Evidence | State |
| --- | --- | --- |
| Direct delegated DALGO reads of one fixed view | [`dalgo2logs` PRs #6–#8](https://github.com/dal-go/dalgo2logs/pull/8) | Merged; synthetic-reader and bounded-query tests, no entry persistence |
| OAuth/service contract, encrypted grant registry, direct reader | [`backend` PRs #59–#64](https://github.com/datatug/backend/pull/64), [`#66` alignment receipt](https://github.com/datatug/backend/pull/66) | Implementation merged; cancellation cleanup and form POST callback included; emulator, race, lint, coverage and security CI passed. PR #66 merged at `9d9c79ddbdbf9d979b8057a0482b03822bd03bc7` |
| Reversible activation backend | [backend PR #67](https://github.com/datatug/backend/pull/67) | Merged at `d3112f646accffedcc7279995ea7362c390f8e8c`; reviewed candidate `7a52784dbdf7e923941fc66f202ad8092e15be65` passed [Backend CI 39308503102](https://github.com/datatug/backend/actions/runs/39308503102) and [Security 39308502042](https://github.com/datatug/backend/actions/runs/39308502042). The full local Firestore adapter suite passed with `-race`, 100% statement coverage and zero uncovered statements in 14.416s. It checks genuine overlapping transactions, a separate committed-pause rejection, frozen admission across injected callback replay, held exchange, and final-inventory resume fencing. No forced commit order or observed SDK retry is claimed. Published `v0.38.0` points to that merge commit; post-merge [Backend CI 39309725492](https://github.com/datatug/backend/actions/runs/39309725492) and [Security 39309724286](https://github.com/datatug/backend/actions/runs/39309724286) passed, with WB verifying the exact target at 19:25:05 UTC. No live activation proof exists |
| Reversible activation host | [datatug-cloud PR #17](https://github.com/datatug/datatug-cloud/pull/17) | Merged at `669728342d2b9c36d7b37e254c821ba19b4162bf`, pinning backend `v0.38.0`. Reviewed candidate `a67028d5f689d0ff0f73719d535cd638836780de` passed [Go CI 39313747372](https://github.com/datatug/datatug-cloud/actions/runs/39313747372): build/test, lint, and the actual Firestore-emulator HTTP journey, which passed with `-race` in 8.446s. The local full host race suite had 100% statement coverage and zero uncovered statements. The fixed-pilot-UID control validates task OIDC, bounds bodies and deadlines, and preserves the existing cleanup route. Its synthetic journey uses fake OAuth/KMS/OIDC and does not prove live authenticated task delivery or Google consent/read/revoke. Post-merge [Go CI 39314930149](https://github.com/datatug/datatug-cloud/actions/runs/39314930149) and [deployment 39316404362](https://github.com/datatug/datatug-cloud/actions/runs/39316404362) passed; WB verified all ten exact-target checks at 20:09:41 UTC. Provider readback at 20:09:47 UTC found production `datatug-cloud-00015-vdq` and preview `datatug-cloud-preview-00016-p9d` at 100% traffic, both configured with the merge SHA image tag. Their immutable image digests are respectively `dfc28e7b327c4950d8ec1d1955851975b08a0aa3a4c2a32659d905b18d773115` and `6af0f4fc286dce15a5fad73c42213750d73d50daf60ebd36f7e67bdab7adaaa9`. Both had no personal Cloud Logs environment variables or OAuth secret bindings, and no pilot Scheduler job existed. GET health checks at 20:10:10 UTC returned 200 with `X-Sneat-Service: datatug-cloud` for both services. The code is deployed with the master feature off; no live lifecycle invocation or hosted Google journey is claimed |
| Connection and evidence UI | [`datatug-apps` PRs #396](https://github.com/datatug/datatug-apps/pull/396) and [#398](https://github.com/datatug/datatug-apps/pull/398), [`#403` alignment receipt](https://github.com/datatug/datatug-apps/pull/403) | Implementation merged; hosted connection UI is gated; local evidence rows use AG Grid; local and CI browser E2E passed. PR #403 merged at `0fc93214ff275edabd86d23881e5f566a1a4bddb`; exact candidate build, test and Cloud Logs E2E passed |
| Pilot UI access gate and sign-in recovery | [Apps PR #407](https://github.com/datatug/datatug-apps/pull/407) | Merged at `a63cfe2883d1078f9a3c1058483197126ea372de`; post-merge [CI run 39292521711](https://github.com/datatug/datatug-apps/actions/runs/39292521711) passed preflight, build, test, Cloud Logs E2E and deploy-gate checks. Cloudflare workflow [39293867514](https://github.com/datatug/datatug-apps/actions/runs/39293867514) succeeded; public build-info reports the same SHA and version `0.1.0+1201`. Production UI is on behind the verified-UID server gate; Business Preview is off. No post-deploy CI smoke was run. Fresh basic Google sign-in reached the authenticated connections page, which showed no accounts or views. Authenticated GET connections returned HTTP 200 at `13:09:41.557Z` and `13:09:41.993Z`; POST OAuth start returned HTTP 200 at `13:10:08.380Z` on production `00010-bjt`. Google did not recognize the selected provider account in the Logging chooser; no consent, grant or query was completed |
| Hosted route composition and task cleanup | [datatug-cloud PR #9](https://github.com/datatug/datatug-cloud/pull/9), [PR #12](https://github.com/datatug/datatug-cloud/pull/12), [PR #14](https://github.com/datatug/datatug-cloud/pull/14), and [PR #15](https://github.com/datatug/datatug-cloud/pull/15) | PR #9 merged at `4d180b91`; reviewed head, Build & test and Lint passed. PR #12 merged at `1d2822f2e0d853352738a01bb3bb39846d5efc98`; deployment [39272895816](https://github.com/datatug/datatug-cloud/actions/runs/39272895816) succeeded at revision `00008-tw6` with GET health 200. PR #14 deployment [39288374923](https://github.com/datatug/datatug-cloud/actions/runs/39288374923) succeeded; provider readback found `00012-stx` at 100% traffic with GET health 200. PR #15 merged at `e5925566ab0c81f95a7bbe46b73fd9836f1c024f`; its private user-routes switch omits all eight user handlers while retaining cleanup registration, credentials and task OIDC. Post-merge Go CI [39298598295](https://github.com/datatug/datatug-cloud/actions/runs/39298598295) and deployment workflow [39298871944](https://github.com/datatug/datatug-cloud/actions/runs/39298871944) passed. Its 18:05:25 UTC provider snapshot found production `00013-hjk` and preview `datatug-cloud-preview-00014-np9`, both at 100% traffic on image `e5925566ab0c81f95a7bbe46b73fd9836f1c024f`, with no Cloud Logs environment variables, OAuth secret bindings or pilot Scheduler job. |
| Host README-only release follow-up | [datatug-cloud PR #16](https://github.com/datatug/datatug-cloud/pull/16) | Merged at `b5c6c8444fd3c98c177d441b47f691cbdfae1d5b`; candidate build/lint [39299058079](https://github.com/datatug/datatug-cloud/actions/runs/39299058079), post-merge Go CI [39299663820](https://github.com/datatug/datatug-cloud/actions/runs/39299663820), and deployment [39299949750](https://github.com/datatug/datatug-cloud/actions/runs/39299949750) passed. Provider readback at 18:11:59 UTC found production `datatug-cloud-00014-lhg` and preview `datatug-cloud-preview-00015-6w6`, both at 100% traffic on image `b5c6c8444fd3c98c177d441b47f691cbdfae1d5b`; both had no Cloud Logs environment variables or OAuth secret bindings, and no pilot Scheduler job exists. GET health returned 200 at 18:12:01 UTC. No grant, live cleanup or shutdown proof exists |
| Initial opt-in production deployment wiring (2026-10-19) | [`cicd` PR #44](https://github.com/sneat-co/cicd/pull/44), [`datatug-cloud` PR #14](https://github.com/datatug/datatug-cloud/pull/14) | Merged at `9571a3c78a448598ccce184a2aa7625ab07d6a53` and `4868105ee730ec3c2ae4914f78970214da0b9676`; independent review and exact-head CI passed. Generic CICD step tests: 133 passed. Post-merge host Go CI and deployment passed; direct provider readback and health are recorded above. This is the earlier deployment state; current setup is recorded in the summary and pilot guide |
| Scheduler header fix and current pilot setup (2026-10-20) | [`cicd` PR #45](https://github.com/sneat-co/cicd/pull/45), [host auth fix PR #25](https://github.com/datatug/datatug-cloud/pull/25), [pilot execution receipt](https://github.com/datatug/datatug-demo-project/blob/main/docs/logs/google-cloud-pilot.md) | PR #45 merged at `3d89e25a11d042b4c855765e58ec876838780679`; 133 local fake-provider checks passed, candidate GitHub CI showed 8 checks, and post-merge main CI passed 6/6. OIDC updates clear the workflow-owned custom-header map. Host PR #25 merged at `9d055b12fae5ae74793e8f452bd7fc8714e22d33`; exact-main CI passed 10/10 at `13:00:44Z`, and deployment `39426215088` succeeded. Provider readback at `13:07:48Z` found production `datatug-cloud-00010-bjt` and preview `datatug-cloud-preview-00009-2fx` at 100% on the merged image. Production has master feature/routes `1`, secret references `2`/`1`, and allowlisted UID match; the approved resume persisted epoch `2`, `paused=false`. All four Belgium/Frankfurt health checks returned 200. Preview and Frankfurt services have no personal Cloud Logs settings. The approved recurring job is enabled; scheduled OIDC target requests returned HTTP 200 at `13:06:08Z`, `13:07:07Z` and `13:08:04Z` without `X-Cron-Secret`. They do not establish which records were deleted. The disposable Scheduler probe was deleted and returned `NOT_FOUND`. Fresh Google basic sign-in reached the authenticated connections page; GET connections and POST OAuth start returned HTTP 200 on production `00010-bjt`. The page showed no accounts or views. The separate Logging chooser did not recognize the selected provider account; no consent or grant was completed. Positive/negative view reads, Cloud Logging consent and a real log query remain unverified. No Google user grant was added |
| Firebase server-only rules | [`sneat-firebase` PR #61](https://github.com/sneat-co/sneat-firebase/pull/61) | Merged at `a78abcd7e55258b93ba996b7bfe1e48149cfe029`; 263 tests across 20 suites; live deploy run [39272536444](https://github.com/sneat-co/sneat-firebase/actions/runs/39272536444) succeeded; exact byte readback SHA `9038e32d45753b5ebe62ac078743da5b0cd7eadb40b9dfedae723c8bf245d8d2`, ruleset `bbb016cf-1cea-4cc5-adb8-79bae63d9c2b` |
| Reproducible local AG Grid investigation | [`datatug-demo-project` PR #55](https://github.com/datatug/datatug-demo-project/pull/55) | Merged; real app screenshots and standalone illustrative HTML; synthetic logs only |
| Narrow synthetic Google log view and direct DALGO read | [`datatug-demo-project` PR #57](https://github.com/datatug/datatug-demo-project/pull/57) | Merged; eight clearly marked synthetic entries in `sneat-eur3-1`, fixed view read verified. A restricted test identity now has a view-only binding; positive/negative reads, hosted consent and browser log journey remain open |

The hosted browser test uses an emulator and synthetic API responses. It does
not prove a real hosted Google OAuth consent, positive/negative reads through
the restricted test identity or browser log journey. The narrow synthetic view
and direct DALGO read in `sneat-eur3-1` are recorded separately in PR #57; the
source read used an Editor identity, so it does not prove least-privilege access.
The hosted UI does not yet
provide log→record→log navigation; the local Northwind demo does. That local
demo reads files directly, without any log database copy.

## Problem

A developer can inspect their own Cloud project locally, but a hosted DataTug user cannot authorize a specific Cloud Logging resource and see related entries as themselves. Sharing the current ADC credential or owner token among hosted users would make one principal's permissions stand in for another's. Signing in to DataTug with Google establishes a DataTug identity; it does not itself grant Cloud Logging access.

## Behavior

### User journey

1. A signed-in user opens **My connections → Google Cloud logs**, chooses **Connect Google Cloud**, and sees that DataTug will read entries from one selected log view without retaining copies. DataTug starts an authorization-code flow with a dedicated Cloud Logging web client in the existing shared Google Cloud project `sneat-eur3-1`, requesting `openid email https://www.googleapis.com/auth/logging.read`. The verified Firebase session UID remains the DataTug identity; the delegated Google grant is linked to that same UID after the user completes consent. The user can choose which Google account to authorize, and DataTug shows it clearly.
2. After consent, the user enters a project ID and full view name `projects/{project}/locations/{location}/buckets/{bucket}/views/{view}`. The server validates that the view belongs to the selected project and performs a bounded read probe with the user's token. It never requests project-list access just to populate a picker. The UI explains how a project administrator can grant `roles/logging.viewAccessor` on a narrow view. Permission failures are actionable without echoing Google error bodies or log content.
3. The user chooses this personal connection in a DataTug project and filters by bounded time, severity and an order/correlation ID. For each query the server resolves the opaque connection ID **and verified DataTug UID**, fixes the resource to the stored view, validates query shape and limits, and reads Cloud Logging directly. The browser never receives a Google token. The initial Cloud view disables the free-text control with an explanation because the current adapter has no message/text predicate.
4. DataTug renders timestamp, severity, source, message, raw evidence and source reference in the existing log view. Opening a related record and returning to logs preserves connection ID and query criteria. Any investigation preservation built for this hosted release stores **references, criteria and narrative only**; the current URL/context bridge is not a persisted case. Exact replay re-queries the source by a versioned locator; a broad criteria rerun must never be mistaken for proof that a specific entry disappeared.
5. **Remove view** disables only that connection and evicts any associated in-memory result. It does not revoke Google permissions; other views under the same grant remain usable. **Revoke Google access** is a separate, explicit account-level action that first disables all of that account's DataTug log views and then asks Google to revoke the account's permissions for the shared Sneat app. The confirmation warns that other Sneat or DataTug clients using this Google authorization may ask for consent again; the Sneat/Firebase user account remains. The disabled encrypted grant remains in a bounded revocation tombstone until cleanup succeeds. If provider revocation remains unconfirmed after a finite retry window, destroy every stored/cached token, permanently fence and drain that grant generation's refresh/query/revoke jobs, then retire the tombstone to a credential-free terminal audit event. Reconnect for the same Google `sub` is blocked across **all** DataTug UIDs until this cleanup barrier completes. A later fresh consent, verified Google subject and read probe may then reconnect; warn that Google-side revocation was unconfirmed and direct the user to Google's access controls. An already-sent provider operation may still take effect later, so reconnection can fail and must never be advertised as guaranteed. Independent Google revocation fails queries closed and prompts reconnection.

### Architecture and contracts

```text
DataTug browser (Firebase session)
  → DataTug Cloud API (verified UID, route policy, CSRF/origin, per-UID quotas)
  → personal registry (UID + fixed project/view + encrypted OAuth grant)
  → request-scoped DALGO cloudlogging adapter (user token, fixed view, bounded query)
  → Google Cloud Logging entries.list (no export, ingestion or shared index)
  → ordinary DALGO records and DataTug log presentation
```

The registry models one provider grant per DataTug UID + Google `sub` + OAuth project and one or more view connections under it. A server-enforced unique constraint on `(OAuth project, Google sub)` permits **one active or revocation-pending DataTug UID globally**, not merely one grant per UID; callback insertion is atomic so simultaneous authorizations cannot create two owners. If that Google account is already linked elsewhere, return a generic conflict without revealing the other DataTug account. Linking it to a different UID requires the old owner's explicit disconnect and completed cleanup barrier, including destruction of credentials and old-generation jobs; if the old DataTug account is inaccessible, use an audited support recovery with independent identity proof, never an automatic timeout takeover. The grant stores verified stable `sub`, display email, state, timestamps, key version and encrypted refresh token; each view connection stores opaque ID, owner UID and fixed project/view. Email is display-only, never an ownership key. The registry stores **no entries, raw payloads, query results, query text or business identifiers**. It is server-only, separate from browser-writeable project connection declarations. Project membership cannot grant access to a personal token.

Use a trusted dedicated Cloud Logging OAuth **web client** in the existing shared Google Cloud project `sneat-eur3-1`. Keep the Firebase login identity and OAuth application shared across Sneat and DataTug; link the separately consented `logging.read` grant to the verified DataTug UID from the same Firebase session. Use server-side authorization-code exchange, an exact redirect URI, state bound to the DataTug session/UID and one-time nonce, PKCE, and offline access only after explicit consent. Keep this client's requested scopes narrow. Validate ID-token issuer, audience, expiry and nonce before storing `sub`; do not treat the Google ID token as DataTug authentication. Store refresh tokens encrypted under a managed key, apart from configuration and project data; restrict decrypt/use to the query service identity and audit it. Keep short-lived access tokens only in bounded server memory. If Google returns no refresh token, show a reconnect path; never substitute server ADC. Token refresh and the read probe use the same user grant. Never put credentials in URLs, browser storage, project files, logs, tracing or error responses.

**What KMS protects:** The current vault calls Cloud KMS directly to encrypt each
user's refresh token; Firestore stores only its ciphertext and grant metadata.
Authenticated encryption binds the ciphertext to the UID, OAuth project, Google
subject, grant generation and grant ID. Decryption requires those same values
and permission on the configured key. The encryption key remains managed by
KMS. Secret Manager separately holds the shared OAuth client secret and a
32-byte AES-GCM key that encrypts short-lived pending PKCE verifiers and nonces,
binding them to the state hash and UID. KMS does not store log entries; this
feature reads logs directly from Google Cloud Logging without copying them
into Firestore.

The production app at `datatug.app` calls `api.sneat.cloud`, so a SameSite cookie set by the cross-site OAuth-start API request cannot reliably bind the Google callback. The implementation uses a two-step completion instead: the authenticated start returns an authorization URL and random app binding secret held only in the redirect tab's session storage; the server callback exchanges the code and seals a short-lived pending grant, then redirects to a fixed frontend route with an opaque one-time receipt. The signed-in browser submits that receipt and binding secret to a final API route; the server checks the original UID, binding hash, expiry and one-time consumption before claiming the Google account. The browser scrubs the receipt from its URL before telemetry or rendering. A bounded recovery worker destroys expired pending credentials. Google code and tokens never reach the frontend.

The adapter now accepts a trusted request-scoped user token source and one validated project-local log-view resource, as well as the earlier local ADC mode. It builds `entries.list.resourceNames` from that fixed view; no DALGO query parameter can select a project or view. Time-window, result, byte, timeout and page/scan budgets are enforced, including empty pages with a nonempty `nextPageToken`. Unsupported predicates are rejected. DALGO remains read-only with ordinary records/schema. OVDB can mount this DB through its generic mount; no OVDB core branch inspects the adapter name or introduces a log-specific public field or endpoint.

The hosted v1 query whitelist is inclusive UTC `timestamp >=` and `timestamp <=` (both required), `severity =` for documented enum values, and equality for `correlation_id`, `order_id` and `service` against the adapter's mapped Google fields. The initial identifier equality contract accepts **ASCII operands and mapped values only**, compared without ASCII letter case. A non-ASCII operand or candidate returns an explicit unsupported-query result. Google Logging's broader string `=` comparison uses Unicode NFKC_CF; the adapter does not silently claim that full behavior. The UI must label the ASCII limit. The existing adapter has no message/free-text predicate. Regex, arbitrary SQL, full-text, substring, OR groups, offset, keyed `Get` and unsupported ordering are rejected or disabled, never approximated. A later text feature needs a separately tested contract for Google `SEARCH` token semantics and must not be presented as substring search.

The hosted API owns authorization. It may call DALGO directly or delegate through a private request-scoped generic OVDB mount; use the simpler measured integration. Neither path exposes the local owner token as a hosted credential. A normal OVDB client remains able to discover/query an authorized generic mount using the unchanged public protocol.

The current Cloud adapter's `source_ref` is a hash of an SDK reconstruction, `Get` is unsupported, and its ordinary query does not filter by reference. For exact replay, add a **DataTug-internal** bounded resolver, without adding a public OVDB operation: save a versioned, payload-free locator containing connection ID, fixed view identity, `logName`, timestamp, `insertId` when present, and an entry digest/algorithm version. Re-fetch within that view using the narrowest supported `entries.list` filter, inspect all candidates within a strict page/byte/time budget, recompute the digest and return exactly one match. If the provider omits `insertId`, the timestamp is absent, candidates collide, the page budget is exhausted or the digest differs, return **unavailable or ambiguous**, never another entry. An adapter change to canonicalization must retain the old digest algorithm for existing locator versions or explicitly mark them unresolvable. Test replay after newer entries push the original past the normal result cap. A general OVDB keyed `Get` remains unsupported.

| Route intent | Required server decision | Allowed result |
| --- | --- | --- |
| OAuth start | Valid browser DataTug session; server-bound UID, one-time state, nonce and PKCE | Authorization URL and browser-held finish binding; no Google credential |
| OAuth callback | One-time state, nonce and PKCE; verified Google identity and scope | Short-lived server-sealed completion and opaque receipt; no active connection yet |
| OAuth finish | Valid browser DataTug session; original UID, receipt, finish binding and expiry; atomic subject claim | Personal grant metadata; no browser credential |
| List/create/remove `me/log-connections` | Verified UID; owner equality; CSRF/origin checks on mutations | Personal view metadata only |
| Revoke Google access | Verified UID; explicit shared-project warning; disable every child view first | Request Google revocation for the shared Sneat app and retain a bounded encrypted cleanup tombstone; other clients may require consent again |
| Discover/query | Verified UID; owner equality; active grant; fixed view; bounded supported query | Current entries only, never persisted |
| Reopen saved reference | Re-authorize current viewer independently, resolve versioned locator exactly | Current evidence, unavailable or ambiguous marker |

The current DataTug Cloud API maps `/v0/datatug/` GET, HEAD and OPTIONS to `datatug:projects:read`, and other methods to `datatug:projects:write`. A read-only log query is POST: add an explicit path/method policy requiring a **read** grant (`datatug:projects:read` in the initial browser/CLI model) plus the personal-connection check. Unrelated POST routes retain their write policy. Personal connect/disconnect is browser-session-only initially; future CLI management needs a separate grant rather than pretending it is a project write. Firebase verification remains the UID and active DataTug OAuth-grant authority. Caller-supplied UID, owner, project, view and token never override server state.

Google IAM decides which entries the delegated user can see; `logging.read` alone grants no project access. Require a project-local log view for hosted v1, avoiding folder/organization or whole-project fan-out. Show the exact view and Google account before enabling the connection. DataTug cannot assume a log payload contains no secrets or personal data: display raw details as inert text, set `Cache-Control: no-store`, disable public payload links and automatic AI processing, redact connection IDs/project names/log references/identifiers/content from analytics and diagnostics, and audit only connection events plus query metadata/counts (never filter text or entries). **Workspace sharing is deferred.** If reference-only investigation sharing is later enabled, each viewer needs their own Google grant and current view IAM, and a denied viewer never sees the author's payload.

The implemented personal hosted query limit is an inclusive range of at most 24 hours, 200 entries and 30 seconds; the UI must state and enforce that 24-hour maximum before submitting. Standalone adapter ceilings do not expand the hosted contract. Set one view per connection, finite per-user concurrency/rate, response-byte cap, cancellation on client disconnect and finite pagination. Distinguish no matches, cap reached, more possibly available, expired grant, IAM denial and unsupported filter. Source references are pointers, not durable copies; Cloud retention, deletion and IAM changes can make them unresolved.

### Alternatives

| Option | Decision | Why |
| --- | --- | --- |
| Shared DataTug service account | Reject for personal MVP | Central credential magnifies cross-tenant risk and makes DataTug model customer IAM. Design workspace-admin access separately. |
| Reuse Firebase Google login token | Reject | Sign-in is not Cloud Logging consent and cannot be assumed to carry `logging.read` or a refresh grant. |
| Browser calls Google directly | Reject | Cannot reliably enforce server-side fixed view, quota and replay ownership. |
| User OAuth with project-wide reads | Defer | Less adapter work but broader access and regional fan-out; retain current ADC project path for local single-owner use. |
| User OAuth with fixed log view | **Recommend** | Google view IAM and delegated identity remain the data boundary. |
| Copy logs into DataTug DB/index | Reject | Conflicts with the direct-evidence scope and creates a retention/aggregation product. |

## Acceptance Criteria

1. A personal connection cannot query until separate Google consent and a live bounded probe of one project-local view. Wrong UID, view/project mismatch, changed IAM and revoked grants fail closed.
2. The browser receives no Google credential. Registry is server-only; refresh tokens are encrypted under a managed key; removing one view disables only its queries, and explicit Google access revocation disables all child views before provider revocation cleanup completes.
3. Discovery/query uses ordinary read-only DALGO collections and direct bounded `entries.list` calls. No log payload/result is persistently stored in a database, cache, analytics event or investigation; responses set `Cache-Control: no-store`.
4. An unmodified OVDB client discovers/queries an authorized generic mount using the existing public contract; an OVDB protocol diff is empty.
5. DataTug completes log→record→log navigation against current source data. The Cloud UI exposes only the stated filter whitelist, labels its ASCII-only case-insensitive identifier equality, and disables free text with an explanation. Consent, denied, expired, empty, capped, error, narrow/mobile and evidence-unavailable states are tested. If exact reference replay is offered, it uses the bounded locator resolver and does not treat a broad rerun as exact replay.
6. Negative tests cover IDOR/cross-UID, concurrent OAuth callbacks for the same Google `sub` under two DataTug UIDs, generic duplicate-owner denial without account disclosure, relink only after explicit old-owner release and cleanup barrier, OAuth state/PKCE replay, missing refresh token, token disclosure, view mismatch, IAM revocation, unsupported query, empty pages with page token, throttling, cancellation, large payload, removal with other views still linked, provider-wide revoke failure/tombstone cleanup, reconnect blocked across UIDs during pending revocation, late old-generation refresh/revoke jobs unable to affect a new grant, and exact replay after newer entries push the target beyond the normal result cap. Use Google client test doubles, then the existing narrow synthetic view in `sneat-eur3-1` for the gated live read pilot. Verify view-only IAM with a low-privilege identity before claiming least-privilege access. No production credentials in CI.
7. Security review verifies scope, storage, redaction, retention, user disclosure, Google consent-screen requirements and least-privilege IAM before public rollout. Pilot measures API cost and latency rather than inferring a guarantee from local benchmarks.
8. New personal Google OAuth activation is fail-closed until a private, fixed-pilot-UID operator action initializes its pause epoch. Pause/resume use expected-epoch compare-and-set; missing, stale or future epochs conflict without mutation, same-state retries are idempotent, and overflow fails closed. A real transition increments the epoch; zero is never an open-state wildcard.
9. Pending admission, pending consumption, completion persistence and grant activation each transactionally verify the same initialized, open UID lifecycle epoch. A callback consumed before pause cannot save or activate after pause/resume; old pending rows and completions remain invalid after resume. Existing user grants/connections and lifecycle fields survive every user-row write. `Grant.Generation` remains a separate fence.
10. A synthetic host-plus-Firestore-emulator journey overlaps pause and completion activation, accepts either valid serialization, and checks that the final paused epoch and grant inventory agree with the activation outcome. A separate post-committed-pause control verifies rejection; no forced Firestore commit order or automatic SDK retry is claimed. The journey lists and explicitly disconnects accepted grants, exercises bounded cleanup/reconciliation, resumes and rejects old state before fresh consent succeeds. Pause, view removal, expiry and local purge never call Google revoke. Tests cover CAS replay, initialization/legacy zero, overflow, UID isolation, row preservation, held callbacks, bounded reconciliation failures and OIDC claim/body rejection. Cross-revision and live-provider gates below remain separate.

## Implementation plan and approval gates

1. **Owner decision:** On 2026-10-19, approved personal delegated OAuth with a mandatory fixed log view and a dedicated web client in the existing shared Google Cloud project `sneat-eur3-1`. Sneat/Firebase login remains shared, and the delegated grant is linked to the verified UID from that session. This approves architecture only, not linking a production account. The feature is Implementing.
2. **Adapter:** Accept trusted user credentials and fixed log-view resource; prove exact filter translation, pagination, limits and read-only behavior. Keep OVDB public API unchanged and test a normal client.
3. **DataTug server:** Add provider grant plus child-view registry, managed-key token lifecycle, OAuth start/callback, view removal and account-wide Google access revocation with a revocation tombstone, route-specific read policy, UID/IAM checks, quotas, sanitized errors and audit. Review threat model before a real-account pilot.
4. **UI:** Connect/reconnect, remove-view and revoke-access states, project/view selection and probe, existing log view and record navigation. A reference-only save/replay bridge is separate, and exact replay must meet the locator contract above before it is advertised. Use the existing synthetic log view in `sneat-eur3-1` for the source pilot and keep its hosted OAuth/browser gaps explicit.
5. **Release:** Verify actual consent screen/scopes, IAM view grant, revocation, denial after grant loss, exact-head CI, cost bounds and browser journey before enabling opt-in hosted access.
6. **Reversible activation barrier:** Backend epoch state and transaction guards are merged in PR #67; the host-private OIDC control and synthetic Firestore journey are merged in host PR #17. Post-merge CI, deployment and provider readback passed with the master feature off. The pilot is not activated. Complete the cross-revision and live-provider gates below before considering pilot activation.

### Reversible OAuth activation pause

The barrier controls only admission of new personal Google OAuth attempts for the one configured pilot UID. Firebase identity, existing grant/query authorization, the public user experience and the OVDB contract do not change. Missing lifecycle state fails closed; an operator must initialize the private state before any new pilot consent. Do not describe this as public onboarding or as revoking a user's Google access.

Store `ActivationPaused` and `ActivationEpoch` on the existing private
`datatugPersonalLogUsers` row, separate from each grant's generation. Lifecycle
initialization and ordinary user-row mutations must preserve existing grants
and connections.

| Lifecycle state | Request | Result |
| --- | --- | --- |
| Missing lifecycle fields or row | `status` | No write; report uninitialized, paused, epoch `0` |
| Missing lifecycle fields or row | `pause(expectedEpoch=0)` | Initialize/preserve the UID row as paused at epoch `1` |
| Initialized paused at epoch `1` | `resume(expectedEpoch=1)` | Open epoch `2` |
| Initialized at epoch `e` | Transition with `expectedEpoch=e` | Change paused state and increment to `e+1` |
| Initialized at epoch `e` | Same-state request with `expectedEpoch=e` | Idempotent no-op; epoch stays `e` |
| Initialized | Missing, zero-wildcard, stale or future expected epoch | Conflict/invalid; no mutation |
| Maximum epoch | Real transition | Fail closed before write; never wrap or reset |

SavePending transactionally requires initialized/open state and stores the current epoch, retaining the first admitted epoch across callback retries. ConsumePending reads the pending record and the same UID row in one transaction, consumes once and returns that epoch internally. Completion carries the consumed epoch through provider exchange/sealing; SaveCompletion rechecks it rather than stamping a newer epoch. ActivateCompletion rechecks the completion epoch against the same open lifecycle row. These transaction guards serialize pause with persistence and activation: activation committed first appears in the post-pause inventory; a committed pause prevents stale activation. Legacy epoch-zero pending/completion records cannot activate after upgrade. Preserve `GrantIDs`/connections while initializing the UID row and preserve pause/epoch through activation, create/remove view, disconnect and final cleanup writes. Keep grant-generation/operation fences intact.

List and the verified user's explicit Disconnect remain available while paused. Pause/resume and view removal, expiry, reconciliation or local purge do not call Google revoke or force-disconnect an active grant; existing grant queries keep their current authorization. Reopening admits only fresh pending attempts and completions. Recommended synthetic journey: initialize with pause CAS `0` → epoch `1`, resume CAS `1` → epoch `2`, create a grant and race a held finish against pause CAS `2` → epoch `3`; list every activation-first grant, explicitly Disconnect it through normal fenced cleanup, reconcile unclaimed rows while paused, then resume CAS `3` → epoch `4`. Verify old epoch-2 callbacks/receipts fail after resume, fresh consent succeeds and stale commands cannot change epoch `4`.

The host control is private and bounded to `status`, `pause`, `resume` and `reconcile`, with expected epoch on mutations. It derives its sole UID from the configured pilot allowlist and accepts no UID/provider subject/grant selector. Before parsing the action, validate the Google-signed OIDC token's exact audience, Google issuer, verified email and exact approved invoker service account; this activation authority requires owner review. Keep the control present when user routes are stopped but the master service/key flag remains on; omit it in preview and under master-off. Return only paused/epoch, bounded flags/counts and `durableReconciled`, with `Cache-Control: no-store`; do not log or return identities, receipts, keys or tokens. Taskshost v0.7.1 returns only empty 200/500 and cannot provide this status/CAS contract, so use a small host-local handler with injected token validation rather than pretending the task receiver supports it. This adds no public endpoint, Firebase identity change or OVDB contract.

Reconciliation is a durable bounded result, not proof of transient callback quiescence or Google-side revocation. Check remaining pending/completion rows and active/revoking inventory with bounded existence queries, then revalidate the same paused epoch before `durableReconciled=true`. No pending/completion rows or active/revoking grants may remain for that result; active grants require explicit Disconnect, so a paused status with remaining active grants is not reconciled. A concurrent resume invalidates the result; unknown or incomplete visibility must return remaining/unknown, never zero. Forced lifecycle reconciliation purge of unclaimed rows, including unexpired rows, requires the paused epoch; ordinary expiry cleanup remains allowed while open, and queued cleanup remains generation-fenced. Keep terminal revocation outcomes and operation fences; preserve `PreviousRevocationUnconfirmed` on reconnect, and retain a false provider-revocation outcome as an honest warning without retrying it against a future grant.

Before initializing or trusting the barrier in production, deploy guard-aware writers everywhere and retire traffic, tags/direct revision access and scheduled paths to binaries that ignore the epoch or can erase lifecycle fields through whole-row writes. Prove their in-flight writes ended; do not roll back to an unguarded binary while the master flag is on. Reversible pause does not make incompatible rollback safe. Code/source checks and emulator tests do not prove deployed OIDC invocation, old-revision exclusion, transient callback quiescence or real shutdown. Replacement OAuth credentials, consent, low-privilege view IAM, provider read, and separately verified live cleanup/shutdown remain manual release gates.

## Open Questions

- Workspace-managed service identities and project-wide search remain separate decisions.
- The narrow synthetic view and direct DALGO read in `sneat-eur3-1` are verified
  by [demo-project PR #57](https://github.com/datatug/datatug-demo-project/pull/57).
  A dedicated OAuth web client has since been created in that shared project
  with client ID `588648831063-rr6titvvsu52j9n7gqc07uagqoru3lfd.apps.googleusercontent.com`
  and callback `https://api.sneat.cloud/v0/datatug/me/google-log-grants/oauth/callback`.
  The verified Firebase identity at `auth.sneat.co` remains shared with Sneat.
  No Google user grant has been completed. Following explicit owner approval,
  the client secret was stored/read-back verified as version `1` of
  `datatug-cloud-logs-oauth-client-secret` in Secret Manager (`europe-west3`),
  and the temporary JSON was deleted. A browser recovery tool emitted the secret
  into private tool output despite suppression. On 2026-10-20 the owner created
  a replacement using **Add secret**. The replacement was privately verified
  to differ from version `1`, stored as enabled version `2`, read back for
  equality, and its temporary JSON deleted. The GitHub deployment selector
  `DATATUG_PERSONAL_CLOUD_LOGS_CLIENT_SECRET_VERSION` is now `2`; stored version
  `1` is disabled. The owner confirmed retirement of the original credential
  in Google's OAuth client. Runtime configuration now references client
  version `2` and state-key version `1`; the master feature and user routes are
  `1`, and the allowlisted lifecycle was later resumed to epoch `2` with durable
  readback `paused=false`. No Google user-token revocation occurred because no Cloud
  Logging user grant exists. A view-only binding was added for the restricted
  test identity on the synthetic view; positive and negative reads remain
  unverified. The earlier source read used an Editor identity.
- Following explicit owner approval, KMS ring `datatug-cloud-logs` and key
  `google-oauth-grants` were created in `sneat-eur3-1/europe-west3`; primary
  software version `1` is enabled for `ENCRYPT_DECRYPT`. After automatic approval
  review required exact permission, the owner separately approved the current
  shared runtime account `588648831063-compute@developer.gserviceaccount.com`
  receiving key-level `roles/cloudkms.cryptoKeyEncrypterDecrypter`, and creation
  of `datatug-cloud-logs-oauth-state-key`. Both were applied. State-secret version
  `1` contains the standard-base64 encoding of 32 random bytes, generated and
  transferred in memory with read-back equality verified; replication is only
  `europe-west3`. Existing Secret Manager access is unchanged. Other workloads
  using the shared runtime identity also gain access to the new KMS key, an
  explicitly approved pilot limitation.
  Generic deployment wiring is merged in [CICD PR #44](https://github.com/sneat-co/cicd/pull/44)
  and [host PR #14](https://github.com/datatug/datatug-cloud/pull/14); the
  Scheduler header fix is merged in [CICD PR #45](https://github.com/sneat-co/cicd/pull/45).
  Its fake-provider regression rejects remove-header tombstones and verifies
  header clearing in both OIDC update paths. The hosted user routes are enabled
  under current production configuration; the allowlisted lifecycle is active
  at epoch `2`. Bounded
  cleanup is merged in host PR #12 and is invoked through the existing
  `/_tasks` taskshost route, authenticated by task OIDC. The pilot UID allowlist
  applies only to verified user management/query routes. Production readback
  at `13:07:48Z` found Belgium production `datatug-cloud-00010-bjt` and
  preview `datatug-cloud-preview-00009-2fx` at 100% on the merged PR #25 image
  `9d055b12fae5ae74793e8f452bd7fc8714e22d33`. Production master feature
  and user routes are `1`, secret references are versions `2`/`1`, and
  the allowlisted lifecycle is active at epoch `2`. The private OIDC
  control returned 200 for authenticated status and 401 without
  authentication. The managed recurring job is enabled; scheduled OIDC targets returned HTTP 200 at
  `13:06:08Z`, `13:07:07Z` and `13:08:04Z` on production revision
  `datatug-cloud-00010-bjt`, without `X-Cron-Secret`. These responses do not
  establish which records were deleted. The disposable probe was deleted and
  returned `NOT_FOUND` at `13:10:25Z`. All four Belgium/Frankfurt health checks returned 200;
  preview and Frankfurt services have no personal Cloud Logs settings.
  The existing task invoker already has `run.invoker`; no new invoker IAM
  grant was added. No Google user grant exists. Host PR #25's Firebase
  `auth_time` fix is merged and deployed; exact-main CI passed 10/10 and the
  deployment succeeded. Earlier 2026-10-20 provider snapshots above are historical;
  current receipts and open checks are in the [pilot guide](https://github.com/datatug/datatug-demo-project/blob/main/docs/logs/google-cloud-pilot.md).
  Source review identified a shutdown race: an in-flight `finish` can activate
  a completion after the initial revoke inventory, removing that completion
  and creating an active grant; cleanup only handles expiring completions and
  queued revocations. Before turning the master feature flag off after
  activation, prove a no-further-activation barrier across all serving and old
  revisions and in-flight start/callback/finish requests, then take a final
  active-grant inventory. Any late active grant needs a separately verified
  revocation path, followed by a generation-fenced drain of queued revocations
  and purge of all completion records while the task route and keys remain
  available. The shutdown sequence has no live proof. PR #15 only stops user
  routes while preserving cleanup; it does not provide that barrier, revoke
  grants, or clear the activation gate. Turning the master feature flag to `0`
  removes cleanup registration and its credential bindings.
- Fresh Google basic sign-in reached the authenticated Cloud Logs connections
  page; GET connections and POST OAuth start returned HTTP 200 on production
  `00010-bjt`. The page showed no accounts or views. The separate Logging
  chooser did not recognize the selected provider account. No Logging consent, grant or query is verified. Public
  helper checks returned HTTP 200, the shared helper's `init.json` identified
  project `sneat-eur3-1` and auth domain `auth.sneat.co`, and its two handler
  scripts matched the canonical Firebase helper by SHA-256. An earlier same-
  origin `createAuthUri` probe returned HTTP 400 `INVALID_CONTINUE_URI`; its
  callback-domain ownership issue remains unresolved. The successful basic
  login does not establish a Logging OAuth grant. No identity, Firebase project
  or auth-domain change has been made.
- The pilot view is temporary and its receipt records a 2026-11-18 review date.
  Confirm Google consent-screen verification classification and Cloud Logging
  billing/quota behavior in the actual Cloud project before public launch.

## Source and contract evidence

- [`dalgo2logs/cloudlogging/cloud.go`](https://github.com/dal-go/dalgo2logs/blob/62c8d1efd058915c82fc726cd8ce201aa602faeb/cloudlogging/cloud.go) and [`personal.go`](https://github.com/dal-go/dalgo2logs/blob/62c8d1efd058915c82fc726cd8ce201aa602faeb/cloudlogging/personal.go): local ADC and request-scoped delegated constructors with bounded reads.
- [`dalgo2logs/docs/operations.md`](https://github.com/dal-go/dalgo2logs/blob/62c8d1efd058915c82fc726cd8ce201aa602faeb/docs/operations.md): the separate loopback owner-token configuration and its local-only boundary.
- [`datatug-cloud/internal/composition/personal_cloud_logs.go`](https://github.com/datatug/datatug-cloud/blob/1d2822f2e0d853352738a01bb3bb39846d5efc98/internal/composition/personal_cloud_logs.go): opt-in direct-reader host composition, bounded task cleanup and pilot UID allowlist.
- [`datatug-cloud/internal/accesspolicy/accesspolicy.go`](https://github.com/datatug/datatug-cloud/blob/caaf347b7a3c848207f34e5341cce23d3c7ca6ad/internal/accesspolicy/accesspolicy.go): read-only POST query route policy; the exact form POST callback exception is in host composition above.
- [`backend` PR #64](https://github.com/datatug/backend/pull/64): OAuth callback code and state arrive in a bounded form POST body, avoiding callback query parameters in Cloud Run ingress URLs.
- [`backend/host4datatug/auth.go`](https://github.com/datatug/backend/blob/c918eef314e2a2e15d7e8cbf7e277246049b8376/host4datatug/auth.go): verified DataTug UID and active OAuth-grant path.
- [Google `entries.list`](https://docs.cloud.google.com/logging/docs/reference/v2/rest/v2/entries/list): view resource names, `logging.read`, ordering and pagination.
- [Cloud Logging query language](https://docs.cloud.google.com/logging/docs/view/logging-query-language): case-insensitive equality versus token-based `SEARCH` and substring semantics.
- [Cloud Logging IAM](https://docs.cloud.google.com/logging/docs/access-control) and [log views](https://docs.cloud.google.com/logging/docs/logs-views): view authorization.
- [Google OAuth web-server flow](https://developers.google.com/identity/protocols/oauth2/web-server) and [best practices](https://developers.google.com/identity/protocols/oauth2/resources/best-practices): consent, offline tokens, state, storage and revocation.
- [Google OpenID Connect](https://developers.google.com/identity/openid-connect/openid-connect): the stable `sub` identifier and optional email claim.

---
*This document follows the https://specscore.md/feature-specification*
