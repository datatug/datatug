# [review r1 govuk-holidays-admission-sol] Independent adversarial review

Date: 2026-10-08. Runtime: Codex Sol, local host. Reviewed `govuk-holidays-admission-sol.md`. Read-only repository inspection and documentation/source-code web reads; no provider dataset GET/HEAD, payload, fixture, source snapshot or repository edit.

## Verdict

The central decision is supported: GO to propose a bounded direct-browser seam, NO-GO for source admission, Directory activation or claims of working multi-demo joins. The research is not implementation or activation authority. One major implementation-policy omission should be corrected before presenting the seam as ready for approval; one minor validation clarification remains. No contradiction was found in the endpoint, CORS-code, licence, attribution, freshness or division-scoped semantic evidence.

## Findings

### Major M1 — The proposed fetch policy omits redirect rejection

Location: report lines 24, 37 and 51. The proposal fixes the initial HTTPS URL and declines broad redirect/path admission, but does not specify `redirect: "error"`. Browser Fetch defaults to following redirects. Checking a final URL or reviewing a redirect chain after the request cannot prevent the browser from having contacted an unadmitted destination; a one-time clean chain also cannot establish future behavior. The [WHATWG Fetch Standard](https://fetch.spec.whatwg.org/#concept-request-redirect-mode) distinguishes follow, error and manual; error refuses redirects, while manual supplies an opaque redirect response rather than a script-readable chain.

Required correction: explicitly propose `mode: "cors"`, `credentials: "omit"`, `redirect: "error"`, fixed exact URL, no custom headers, no-store and no-referrer. The critical missing option is redirect error; the other two make the existing prose explicit. Describe the activation check as proving the exact endpoint succeeds without a redirect under that policy. If a redirect is encountered, fail closed and seek a separately reviewed canonical resource; do not follow it to discover a replacement. Optional metadata-only network inspection may diagnose a redirect, within task authority, but does not expand runtime admission.

This is a major for the proposed implementation seam, not evidence that the documented endpoint currently redirects. No endpoint request was made.

### Minor m1 — Separate browser readability from header-inspection receipts

Location: report line 51. Verifying the actual CORS header and content readability are different checks. A script can read a successful CORS response without being able to enumerate `Access-Control-Allow-Origin` or arbitrary response headers. The [Fetch Standard's CORS filtering rules](https://fetch.spec.whatwg.org/#concept-filtered-response-cors) restrict script-visible headers. Similarly, absence of a script-readable source-revision or Retry-After header need not establish absence at the network layer.

Suggested correction: specify a metadata-only browser network-header receipt for actual CORS/expiry/revision/retry headers, and a separate real-app-origin readable fetch proof using the proposed options. Do not capture a HAR or tool envelope containing bodies. State that one successful GET proves the successful path only; cancellation, oversize/malformed refusal and error handling need separate bounded implementation proof, while deliberately provoking provider throttling is unnecessary. The report already retains those gates, so this is clarification rather than an activation blocker newly introduced by review.

## Independently checked evidence

- The [official GDS catalogue](https://www.api.gov.uk/gds/bank-holidays/) names the exact full-feed endpoint and its England/Wales, Scotland and Northern Ireland remit. Its documentation points to alphagov/frontend. No endpoint-specific quota or service guarantee was established.
- [GOV.UK reuse guidance](https://www.gov.uk/help/reuse-govuk-content) explicitly lists the JSON Bank Holidays API under OGL reuse. Its rate-limit language concerns scraping generally; the report correctly avoids inventing an API quota.
- The [calendar controller](https://raw.githubusercontent.com/alphagov/frontend/main/app/controllers/calendar_controller.rb) sets wildcard CORS for JSON requests and separate full-calendar/division expiry values. This supports a browser candidate; it supplies neither deployed CDN headers nor an authenticated runtime proof. The [routes](https://raw.githubusercontent.com/alphagov/frontend/main/config/routes.rb) support GET and include the stated NI shorthand redirect.
- The [calendar serializer](https://raw.githubusercontent.com/alphagov/frontend/main/app/models/calendar.rb) emits divisions without its internal last_updated property. The [division serializer](https://raw.githubusercontent.com/alphagov/frontend/main/app/models/calendar/division.rb) emits parent division plus flattened events; the [event serializer](https://raw.githubusercontent.com/alphagov/frontend/main/app/models/calendar/event.rb) exposes stored attributes after date parsing and localization. The report correctly distinguishes this code from a current wire-schema admission. Its civil-date, duplicate/multiplicity, missing/null/empty, derived-context and bunting gates are appropriate.
- [Official calendar documentation](https://raw.githubusercontent.com/alphagov/frontend/main/docs/calendars.md) documents internal storage and publishing behavior, not an update cadence or permanent historical archive. Fetch time is therefore not a source-revision promise.
- The official [National Archives OGL 3.0 page](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/) was accessible through official-domain indexed text, while direct open returned 403. The indexed text supports commercial combination and attribution duties, exclusions and non-endorsement. The report accurately discloses this limitation and preserves exact-resource review. Permission to reuse is correctly kept separate from programme no-retention authority.
- [GOV.UK privacy notice](https://www.gov.uk/help/privacy-notice) supports the provider IP/browser/access-log disclosure. No-referrer and credential omission do not eliminate provider access logs. The report does not claim end-to-end provider no-store.

## Repository and product constraints

Read via `git show origin/main:<path>`, independently from the exact locally available refs:

- datatug/datatug `3973ae1a5bcd1c2959fef3bea7f8aafa54a156e5`: candidate `spec/research/public-data-fabric/candidates/$records/d1-govuk-holidays.json`, assertion `assertions/$records/d1-mapping-holiday.json`, scenario `candidates/$records/d1-scenario-calendar-context.json`, evidence `evidence/$records/d1-govuk.json` and question `questions/$records/d1-historical-context.json`.
- openvaultdb/directory `480f369bedf4c3c8c6d1d67c32a083c6c820de54`: `sources/$records/govuk-bank-holidays.yaml` remains inactive, live-http-via-ovdb and retention none, with execution, rights, schema/binding, retention, freshness and demo-context gates.

Both `git ls-remote origin refs/heads/main` checks failed because the shell could not resolve github.com. These are exact available origin/main object receipts, not refreshed remote-head claims. This matches the report's explicit limitation; no live-current repository assertion is accepted from them.

The report preserves accepted bank-holiday event scope and rejects universal event identity, global calendar and business-closure inference. It preserves old Northwind/Pubs coverage mismatches as dated negative evidence, does not turn Chinook date overlap into admitted jurisdiction, keeps Sakila/date-location uncertainty, treats Employees' validity sentinel separately and labels AdventureWorks as unaccepted. Multi-demo compatibility remains unresolved. Observed min/max dates must not be used later as proof of a complete calendar or an unbroken authoritative history.

The GO is limited to a proposed seam and any separately authorized bounded implementation work. All live provider reads, canonical schema/binding admission, Directory vocabulary changes, product availability and activation remain separately gated. Fixing M1/m1 would make this research recommendation reviewable; it would not satisfy those gates.

R1 VERDICT: blockers=0 majors=1 minors=1 land=no

## [review r2 govuk-holidays-admission-sol] Correction recheck

Re-read the revised note on 2026-10-08, local Codex Sol runtime. This recheck used the revised note and the existing R1 evidence only; no provider request, repository edit or new remote-head claim.

- M1 resolved at line 37: the exact fixed HTTPS URL now requires `mode: "cors"`, `credentials: "omit"` and `redirect: "error"`, retaining no-store/no-referrer and rejecting destination discovery through redirects. Line 51 now requires successful real-origin readability with these exact options and no redirect. This closes the omission without expanding admission.
- m1 resolved at line 51: real-origin readable fetch proof is separate from metadata-only network-header receipts. Script-header filtering, absent-header interpretation, body-bearing HAR/tool-envelope exclusion and separate successful/error-path proof are explicit. The revision also explicitly declines provider-throttling provocation and inference of completeness from date extrema.
- No finding declined. No new finding introduced by these revisions. The existing exact-resource/runtime/schema/binding/rights/retention/freshness/demo gates remain intact.

The research recommendation passes independent review. GO remains limited to proposing the bounded implementation seam; source admission, Directory activation and multi-demo compatibility remain NO-GO. `land=yes` here means the revised research recommendation may be presented, not that any repository or source may be activated or landed.

VERDICT: blockers=0 majors=0 minors=0 land=yes
