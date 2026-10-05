---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: BigQuery Dual Runtime

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric/bigquery-dual-runtime?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric/bigquery-dual-runtime?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric/bigquery-dual-runtime?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric/bigquery-dual-runtime?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

User-prioritized BigQuery access through required Go and preferred browser-direct JavaScript paths, with retained source and runtime gates.

## Problem

Users should inspect and query useful BigQuery public-data candidates in DataTug and Go consumers while choosing their own execution identity/project and bounded query budget. A public provider catalogue, Google login or transport prototype alone establishes neither source rights nor semantic compatibility. This priority must preserve W1 GeoNames/ROR, ECB FX and all useful unselected research.

## Behavior

### User journey

1. Browse candidate metadata without login, with research/admission and rights/freshness limits visible.
2. At Run, preserve the draft through explicit DataTug Google sign-in/signup and separate BigQuery consent. Verify the execution account independently with stable subject identity; select a user-owned project.
3. Validate the reviewed source identity, schema, native-table configuration and location; preview the exact plan and bytes cap. A changed plan or execution context invalidates approval.
4. Run the bounded read through browser-direct JavaScript or the required Go consumer. Preserve exact values and the same job identity across at least two pages. Stop waiting and consented remote cancellation remain distinct.
5. Promote only after source-specific rights/metadata/access-route review, required canonical-profile publication and independent actual runtime acceptance.

### REQ: both-runtimes

BigQuery MUST be accessible from both Go and JavaScript. Browser-direct execution is preferred for DataTug. Reuse the maintained `dal-go/dalgo-http-adapters/packages/bigquery`; a Go sibling module remains proposed after current ownership/reuse admission. The existing DataTug CLI query/run family is a candidate real consumer, whose exact bounded analytical input route remains A0. A local Go server harness is also required. The [canonical amended architecture r3](../../../research/public-data-fabric/evidence/bigquery-dual-runtime/README.md) MUST pass fresh review before implementation scoping; approved plans follow the established Workbench SpecScore route.

### REQ: bounded-source-and-identity

Both runtimes MUST validate bounded raw HTTP cell presence/count/type before any lossy generated Go decoding or result export. First-slice sources MUST be reviewed trusted native TABLEs with exact datasets.get/tables.get location/identity/schema/configuration checks, rejection of unsupported or conflicting kinds/configurations, and a pre-submit recheck bound to approval. The non-atomic upstream replacement race MUST remain explicit. GIS MUST request openid plus optional email with BigQuery authorization, use same-token fixed trusted Google UserInfo, and bind execution identity to stable sub; missing email is allowed and missing sub/partial grants refuse. Firebase account identity MUST NOT substitute for execution-account proof. Mandatory dry-run/cap, same-job paging, exact values, no ambiguous resubmission and distinct local stop/remote cancellation apply.

### REQ: preserve-lanes-and-admission

The BigQuery priority MUST preserve W1 GeoNames/ROR, ECB FX and original research/semantic assertions. Access routes MUST stay unreconciled until specialist review; public availability and browser-direct transport MUST NOT imply licence clearance, source admission, immutable source contents or query success. Google registration at Run MUST be explicit and independent of cloud consent, IAM, project setup and approved spending. No automatic cloud project/billing setup or shared-spend proxy is authorized by this feature.

## Acceptance Criteria

### AC: required-go-and-js (verifies REQ:both-runtimes)

**Given** reviewed architecture, accepted A0 consumer/SDK contracts and owner-provisioned runtime contexts
**When** BigQuery delivery is claimed
**Then** both the actual Go consumer plus server harness and the deployed-browser path have bounded metadata → dry-run → capped execution → two-page acceptance receipts, shared immutable raw-wire parity fixtures and exact reviewed source/CI/deployment/landing evidence; neither runtime is silently deferred.

### AC: reject-unapproved-execution (verifies REQ:bounded-source-and-identity)

**Given** omitted or malformed wire cells, unsupported source kinds/configuration, source metadata changed since preview, or partial identity grants/missing stable subject
**When** either runtime validates a response or attempts execution
**Then** malformed cells cannot become SQL NULL, unsupported/changed source approvals prevent submission, identity failures require recovery, missing email alone preserves verified sub, and no etag check is claimed to remove the residual replacement race.

### AC: preserve-unreconciled-research (verifies REQ:preserve-lanes-and-admission)

**Given** the prioritized BigQuery route candidates and existing W1/ECB FX work
**When** research is curated or an access route is evaluated
**Then** all original691 candidates and1177 assertions remain available unchanged, new route candidates retain specialist-review/rights/metadata prerequisites, and no semantic equivalence, admission, licence or live-query success is inferred from this Draft feature.

## Open Questions

Fresh edited-r3 review and A0 exact bounded Go consumer route remain pending. Source-specific rights/metadata/access-route reconciliation, canonical source-access profile and both live runtime gates remain open. See [fabric.questions](../../../research/public-data-fabric/questions/$records/bigquery-dual-runtime-prerequisites.json).

---
*This document follows the https://specscore.md/feature-specification*
