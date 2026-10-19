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

This feature is in implementation, not a claim that hosted access works today. The existing local demo uses application default credentials (ADC) and a fixed project behind a loopback single-owner OVDB server. The delegated adapter, server contracts, encrypted credential storage, grant/connection registry, bounded reader, and feature-gated browser UI are merged. The DataTug Cloud host has not wired the routes, verifier, quotas or cleanup worker; production OAuth configuration, deployment and a live user-grant pilot are also outstanding. The production and Business Preview UI flags remain off.

### Delivery status (2026-10-19)

| Slice | Evidence | State |
| --- | --- | --- |
| Direct delegated DALGO reads of one fixed view | [`dalgo2logs` PRs #6–#8](https://github.com/dal-go/dalgo2logs/pull/8) | Merged; synthetic-reader and bounded-query tests, no entry persistence |
| OAuth/service contract, encrypted grant registry, direct reader | [`backend` PRs #59–#62](https://github.com/datatug/backend/pull/62) | Merged; emulator, race, lint, coverage and security CI passed |
| Connection and evidence UI | [`datatug-apps` PR #396](https://github.com/datatug/datatug-apps/pull/396) | Merged; local and CI browser E2E passed with emulators; production flag off |
| Hosted route composition and worker | `datatug-cloud` | Open; repository requires Claude Code subagents, and that runtime is not authenticated on the founder Mac in this session |
| Real Google grant and view pilot | Disposable nonproduction Google project | Open; no production credential or IAM changes made |

The current browser test uses an emulator and synthetic API responses. It does not prove a real Google OAuth consent, IAM grant or Cloud Logging read. The merged UI does not yet provide hosted log→record→log navigation; the local Northwind demo does.

## Problem

A developer can inspect their own Cloud project locally, but a hosted DataTug user cannot authorize a specific Cloud Logging resource and see related entries as themselves. Sharing the current ADC credential or owner token among hosted users would make one principal's permissions stand in for another's. Signing in to DataTug with Google establishes a DataTug identity; it does not itself grant Cloud Logging access.

## Behavior

### User journey

1. A signed-in user opens **My connections → Google Cloud logs**, chooses **Connect Google Cloud**, and sees that DataTug will read entries from one selected log view without retaining copies. DataTug begins a separate Google authorization-code flow for `openid email https://www.googleapis.com/auth/logging.read`; `openid email` identifies the connected Google account, while `logging.read` is the Cloud Logging API scope. A different Google account from DataTug sign-in is allowed and clearly shown.
2. After consent, the user enters a project ID and full view name `projects/{project}/locations/{location}/buckets/{bucket}/views/{view}`. The server validates that the view belongs to the selected project and performs a bounded read probe with the user's token. It never requests project-list access just to populate a picker. The UI explains how a project administrator can grant `roles/logging.viewAccessor` on a narrow view. Permission failures are actionable without echoing Google error bodies or log content.
3. The user chooses this personal connection in a DataTug project and filters by bounded time, severity and an order/correlation ID. For each query the server resolves the opaque connection ID **and verified DataTug UID**, fixes the resource to the stored view, validates query shape and limits, and reads Cloud Logging directly. The browser never receives a Google token. The initial Cloud view disables the free-text control with an explanation because the current adapter has no message/text predicate.
4. DataTug renders timestamp, severity, source, message, raw evidence and source reference in the existing log view. Opening a related record and returning to logs preserves connection ID and query criteria. Any investigation preservation built for this hosted release stores **references, criteria and narrative only**; the current URL/context bridge is not a persisted case. Exact replay re-queries the source by a versioned locator; a broad criteria rerun must never be mistaken for proof that a specific entry disappeared.
5. **Remove view** disables that one connection immediately and evicts any associated in-memory result. Other views under the same Google grant remain usable. **Disconnect Google account** explicitly disables all of that account's connections before requesting Google grant revocation. The UI warns that Google revocation invalidates every scope and client under the OAuth project, so this integration uses an isolated OAuth project. The disabled encrypted grant remains in a bounded revocation tombstone until cleanup succeeds. If provider revocation remains unconfirmed after a finite retry window, destroy every stored/cached token, permanently fence and drain that grant generation's refresh/query/revoke jobs, then retire the tombstone to a credential-free terminal audit event. Reconnect for the same Google `sub` is blocked across **all** DataTug UIDs until this cleanup barrier completes. A later fresh consent, verified Google subject and read probe may then reconnect; warn that Google-side revocation was unconfirmed and direct the user to Google's access controls. An already-sent provider operation may still take effect later, so reconnection can fail and must never be advertised as guaranteed. Independent Google revocation fails queries closed and prompts reconnection.

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

Use an isolated Google OAuth **project** and dedicated web client for this integration, with server-side authorization-code exchange, exact redirect URI, state bound to the DataTug session/UID and one-time nonce, PKCE, and offline access only after explicit consent. Do not reuse the DataTug sign-in OAuth project or attach other product scopes to this one. Validate ID-token issuer, audience, expiry and nonce before storing `sub`; do not treat the Google ID token as DataTug authentication. Store refresh tokens with envelope encryption under a managed key, apart from configuration and project data; restrict decrypt/use to the query service identity and audit it. Keep short-lived access tokens only in bounded server memory. If Google returns no refresh token, show a reconnect path; never substitute server ADC. Token refresh and the read probe use the same user grant. Never put credentials in URLs, browser storage, project files, logs, tracing or error responses.

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
| Revoke Google grant | Verified UID; account-level confirmation; disable every child view first | Explicit project-wide OAuth revocation and bounded encrypted cleanup tombstone |
| Discover/query | Verified UID; owner equality; active grant; fixed view; bounded supported query | Current entries only, never persisted |
| Reopen saved reference | Re-authorize current viewer independently, resolve versioned locator exactly | Current evidence, unavailable or ambiguous marker |

The current DataTug Cloud API maps `/v0/datatug/` GET, HEAD and OPTIONS to `datatug:projects:read`, and other methods to `datatug:projects:write`. A read-only log query is POST: add an explicit path/method policy requiring a **read** grant (`datatug:projects:read` in the initial browser/CLI model) plus the personal-connection check. Unrelated POST routes retain their write policy. Personal connect/disconnect is browser-session-only initially; future CLI management needs a separate grant rather than pretending it is a project write. Firebase verification remains the UID and active DataTug OAuth-grant authority. Caller-supplied UID, owner, project, view and token never override server state.

Google IAM decides which entries the delegated user can see; `logging.read` alone grants no project access. Require a project-local log view for hosted v1, avoiding folder/organization or whole-project fan-out. Show the exact view and Google account before enabling the connection. DataTug cannot assume a log payload contains no secrets or personal data: display raw details as inert text, set `Cache-Control: no-store`, disable public payload links and automatic AI processing, redact connection IDs/project names/log references/identifiers/content from analytics and diagnostics, and audit only connection events plus query metadata/counts (never filter text or entries). **Workspace sharing is deferred.** If reference-only investigation sharing is later enabled, each viewer needs their own Google grant and current view IAM, and a denied viewer never sees the author's payload.

Defaults should inherit 24 hours, 200 entries and 30 seconds; retain seven-day, 1,000-entry and five-minute adapter hard ceilings only if hosted quotas are tighter. Set one view per connection, finite per-user concurrency/rate, response-byte cap, cancellation on client disconnect and finite pagination. Distinguish no matches, cap reached, more possibly available, expired grant, IAM denial and unsupported filter. Source references are pointers, not durable copies; Cloud retention, deletion and IAM changes can make them unresolved.

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
2. The browser receives no Google credential. Registry is server-only; refresh tokens are encrypted under a managed key; view removal disables its queries, and account disconnect disables all child views before provider-wide revocation cleanup completes.
3. Discovery/query uses ordinary read-only DALGO collections and direct bounded `entries.list` calls. No log payload/result is persistently stored in a database, cache, analytics event or investigation; responses set `Cache-Control: no-store`.
4. An unmodified OVDB client discovers/queries an authorized generic mount using the existing public contract; an OVDB protocol diff is empty.
5. DataTug completes log→record→log navigation against current source data. The Cloud UI exposes only the stated filter whitelist, labels its ASCII-only case-insensitive identifier equality, and disables free text with an explanation. Consent, denied, expired, empty, capped, error, narrow/mobile and evidence-unavailable states are tested. If exact reference replay is offered, it uses the bounded locator resolver and does not treat a broad rerun as exact replay.
6. Negative tests cover IDOR/cross-UID, concurrent OAuth callbacks for the same Google `sub` under two DataTug UIDs, generic duplicate-owner denial without account disclosure, relink only after explicit old-owner release and cleanup barrier, OAuth state/PKCE replay, missing refresh token, token disclosure, view mismatch, IAM revocation, unsupported query, empty pages with page token, throttling, cancellation, large payload, removal with other views still linked, provider-wide revoke failure/tombstone cleanup, reconnect blocked across UIDs during pending revocation, late old-generation refresh/revoke jobs unable to affect a new grant, and exact replay after newer entries push the target beyond the normal result cap. Use Google client test doubles, then a disposable nonproduction project with synthetic logs for a gated live pilot. No production credentials in CI.
7. Security review verifies scope, storage, redaction, retention, user disclosure, Google consent-screen requirements and least-privilege IAM before public rollout. Pilot measures API cost and latency rather than inferring a guarantee from local benchmarks.

## Implementation plan and approval gates

1. **Owner decision:** Approved personal delegated OAuth and mandatory fixed log view on 2026-10-19. This approves architecture only, not linking a production account. The feature is Implementing.
2. **Adapter:** Accept trusted user credentials and fixed log-view resource; prove exact filter translation, pagination, limits and read-only behavior. Keep OVDB public API unchanged and test a normal client.
3. **DataTug server:** Add provider grant plus child-view registry, managed-key token lifecycle, OAuth start/callback, view removal and account-wide disconnect with revocation tombstone, route-specific read policy, UID/IAM checks, quotas, sanitized errors and audit. Review threat model before a real-account pilot.
4. **UI:** Connect/reconnect, remove-view and revoke-account states, project/view selection and probe, existing log view and record navigation. A reference-only save/replay bridge is separate, and exact replay must meet the locator contract above before it is advertised. Use existing DemoDB synthetic fixture for E2E and a disposable Google project for live pilot.
5. **Release:** Verify actual consent screen/scopes, IAM view grant, revocation, denial after grant loss, exact-head CI, cost bounds and browser journey before enabling opt-in hosted access.

## Open Questions

- Workspace-managed service identities and project-wide search remain separate decisions.
- Select the exact KMS/secret-store service and ownership of the isolated OAuth project/client against the deployment environment; do not create production credentials before a threat model and operator runbook.
- Confirm Google consent-screen verification classification and Cloud Logging billing/quota behavior in the actual Cloud project before public launch.

## Source and contract evidence

- [`dalgo2logs/cloudlogging/cloud.go`](https://github.com/dal-go/dalgo2logs/blob/62c8d1efd058915c82fc726cd8ce201aa602faeb/cloudlogging/cloud.go) and [`personal.go`](https://github.com/dal-go/dalgo2logs/blob/62c8d1efd058915c82fc726cd8ce201aa602faeb/cloudlogging/personal.go): local ADC and request-scoped delegated constructors with bounded reads.
- [`dalgo2logs/docs/operations.md`](https://github.com/dal-go/dalgo2logs/blob/40002199025ec819550bbe7f0e764a5e2a83afdc/docs/operations.md): loopback owner-token boundary and hosted gap.
- [`datatug-cloud/internal/accesspolicy/accesspolicy.go`](https://github.com/datatug/datatug-cloud/blob/d860d8067d1d4b615b494ef43da9fa383908b8d2/internal/accesspolicy/accesspolicy.go): method-based policy that treats query POST as write.
- [`backend/host4datatug/auth.go`](https://github.com/datatug/backend/blob/2e194840d886191d674e93b524e93a9fde6b52ed/host4datatug/auth.go): verified DataTug UID and active OAuth-grant path.
- [Google `entries.list`](https://docs.cloud.google.com/logging/docs/reference/v2/rest/v2/entries/list): view resource names, `logging.read`, ordering and pagination.
- [Cloud Logging query language](https://docs.cloud.google.com/logging/docs/view/logging-query-language): case-insensitive equality versus token-based `SEARCH` and substring semantics.
- [Cloud Logging IAM](https://docs.cloud.google.com/logging/docs/access-control) and [log views](https://docs.cloud.google.com/logging/docs/logs-views): view authorization.
- [Google OAuth web-server flow](https://developers.google.com/identity/protocols/oauth2/web-server) and [best practices](https://developers.google.com/identity/protocols/oauth2/resources/best-practices): consent, offline tokens, state, storage and revocation.
- [Google OpenID Connect](https://developers.google.com/identity/openid-connect/openid-connect): the stable `sub` identifier and optional email claim.

---
*This document follows the https://specscore.md/feature-specification*
