---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Pricing and billing

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/pricing-and-billing?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/pricing-and-billing?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/pricing-and-billing?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/pricing-and-billing?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

DataTug has four advertised plans: Free, Pro, Business and Enterprise. Business has a fixed base including seven billable monthly active users (MAU), €20 for each additional MAU, and optional €192-per-unit prepaid additional capacity in independently paid twelve-month lots. The plan shape comes from the founder's 2026-10-06 strategy, copied verbatim in the private DataTug backstage pricing documents. His later explicit choice is EUR-only for now, with the same numeric prices; the attachment's dollar signs and an earlier three-currency discussion do not authorize a USD display, estimate or selector. This Draft distinguishes decided amounts from proposed activity/lifecycle semantics and open implementation proof. The shared money/capacity mechanics must be configurable for other Sneat products; DataTug alone defines which DataTug operations qualify as activity.

This supersedes the Team/Business/Company tier ladder, including the old Business $199/20-contributor/6,000-AI-question meaning. The founder authorized implementation of sufficiently specified Business pieces; this artifact does not itself prove or arm **live** Business checkout or usage billing. Existing EUR Pro subscribers, price identities, entitlements and invoices remain valid until a separately reviewed migration; equal numeric EUR display does not itself authorize changing their provider prices or terms.

The dated [engineering notes](engineering-notes.md) record provider and source boundaries. Their common-anniversary/next-boundary prepaid suggestion is superseded by the founder decisions below. The [Business implementation plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/business-usage-pricing/README.md) records the current contracts, recommendations and activation gates.

## Problem

Seats charge for invited but dormant colleagues. The former six-card page also presents Team/Company and the previous Business tier after the founder replaced them with one usage-based Business plan. A buyer needs a predictable bill; an operator needs to explain each counted person, covered unit and overage without counting logins or client telemetry.

At website base `6d96742`, `apps/io/src/data/pricing.mjs` still lists `team`, `business` and `company` as fixed-price plans and applies the launch offer to all four fixed paid cards. The existing Pro Checkout is a single-plan EUR quote path. No accepted evidence here proves a Business MAU invoice, annual prepaid component, correction or final payment journey. These remain activation gates. The private [current pricing policy](https://github.com/datatug/backstage/blob/main/docs/pricing/README.md) records the earlier Pro/launch/refund decisions; its 2026-10-05 ladder sections are being superseded by this feature and the [2026-10-06 founder strategy](https://github.com/datatug/backstage/blob/main/docs/pricing/founder-strategy-2026-10-06.md). [Stripe's mixed-interval guide](https://docs.stripe.com/billing/subscriptions/mixed-interval) and [meter guide](https://docs.stripe.com/billing/subscriptions/usage-based/recording-usage-api) inform engineering options, not a claim of provider readiness.

## Behavior

### Buyer journey and observable stages

1. A visitor sees Free €0, Pro €19/month, Business from €99/month and Enterprise by contact. Desktop compares Free/Pro/Business in one row with a full-width Enterprise card underneath; mobile stacks them in that order. The page distinguishes working features from coming ones and does not offer Business purchase before billing is ready. No USD estimate or currency selector is part of this decision.
2. A Pro buyer retains the account-bound Pro purchase and entitlement journey. The public €19 Pro buy action is enabled only when a verified EUR quote/catalogue matches the displayed period, amount and offer terms; otherwise purchase stays unavailable. A mismatch refuses before provider session creation. Existing EUR subscribers continue under their recorded price identities until an explicit migration.
3. A future Business buyer sees €99 monthly base including seven billable MAU, unlimited invitations, and €20 for each uncovered MAU in the monthly period. Annual base follows the two-months-free direction (€990 from current stated numbers); separately, €192 prepaid additional capacity buys one extra monthly covered MAU for its own twelve-month term. Customers may purchase additional lots at any time and combine overlapping lots; each paid lot preserves its own start/end unless the customer explicitly accepts an extension quote. Only Pro/Business **base subscriptions** receive the 30% launch discount. Usage at €20 and prepaid capacity at €192 are never launch-discounted.
4. A successful, authenticated, user-initiated hosted DataTug operation is attributed to one verified payer and period. The first qualifying operation for a user/payer/period counts once. Passive or failed events count zero. The monthly statement identifies the evidence, included seven, applicable committed units, uncovered units and charge.
5. The bill can be replayed and corrected without duplicate charges. Provider invoice finalization, late-event correction, monthly/annual base renewal and independently paid capacity purchases are separately proved before Business purchase opens. Mixed-interval subscriptions remain an optional composition, not a prerequisite.
6. A Business Space can turn on SSO during signup or later in Space settings. If SSO is already active for that Space through another Sneat product, the DataTug purchase uses that entitlement without a second €99 charge. Ending DataTug Business must not cancel Space SSO while another product still uses it. A paid activation, not a mere identity-provider sign-in, governs the SSO charge.

#### REQ:four-plan-offer

The public page MUST present exactly Free, Pro, Business and Enterprise as the current new-offer plans, with the founder's layout and honest availability. Retired Team/Company names MAY remain only in attributed history or historical subscription records.

#### REQ:pro-compatibility

EUR-only display and catalogue changes MUST preserve existing Pro price identity, invoicing and entitlements. A Pro buy action MUST remain disabled until its verified EUR quote matches the displayed period, amount and offer terms; any mismatch MUST refuse before a new payable session. No USD estimate, FX conversion or currency selector is required by this decision.

#### REQ:plan-end-access

When a plan ends, the affected cloud projects MUST become read-only and unable to use hosted AI while local query execution remains possible; data and nonpublic visibility MUST remain intact. The exact affected-project cohort on fallback to Free is unresolved and MUST be proved before enforcement.

#### REQ:business-components

Business MUST calculate uncovered monthly units as max(0, counted MAU minus seven included minus prepaid capacity applicable to that window), at the configured €20 unit price. Each confirmed prepaid purchase MUST preserve its independently paid twelve-month term, add to overlapping active paid lots, and reset unused monthly capacity without banking. A later purchase MUST NOT reset or shorten an earlier paid term. An optional accepted alignment extension MUST preserve the earlier paid term and charge only its added duration. Exact mid-window coverage and partial-month extension pricing remain the plan's labeled recommendations until approved.

#### REQ:discount-components

The 30% launch grant MUST apply only to eligible Pro and Business base subscription components. Usage and annual prepaid capacity MUST NOT receive that launch discount. Invitation and renewal transitions MUST preserve the same component boundary.

#### REQ:billable-mau

DataTug MUST derive billable activity from verified successful server operations, attribute it to one authorized payer, deduplicate user/payer/period evidence, and exclude passive or unsuccessful activity. Exact trigger choices below are an engineering proposal awaiting implementation proof.

#### REQ:reusable-billing

Shared billing logic MUST take an immutable versioned per-product configuration, isolate products and payers, and persist auditable dedup/capacity/quote history without DataTug-specific prices or query semantics embedded in its money engine. Minor-unit discount arithmetic MUST reuse the existing exactness/refusal policy rather than silently truncate an inexact percentage result; overflow and invalid configuration MUST refuse before a payable quote.

#### REQ:space-sso

Sneat Space SSO MUST use existing Sneat OIDC/SAML rails as the starting capability, with one €99 monthly Space-owned service/charge across products, enableable during Business signup or later in Space settings. An already enabled Space MUST NOT be charged again; a DataTug product change MUST NOT cancel SSO used by another product. A scheduled SSO cancellation MUST keep sign-in available through the proved paid period and deny new SSO sign-in after its end while retaining independently authenticated administrator recovery.

### Confirmed commercial and existing product rules

| Component | New advertised EUR amount | Meaning |
|---|---:|---|
| Free | €0 | Evaluate and use DataTug personally; previous Free project/AI limits remain only where not superseded. |
| Pro base | €19/month; €190/year | Individual professional use; preserve compatible checkout and owner-scoped included AI. |
| Business base | €99/month; €990/year | Includes seven billable MAU; no conventional seat purchase. |
| Business overage | €20/additional MAU/month | No initial volume tiers. |
| Business prepaid additional capacity | €192/unit/twelve-month term | Independently paid additive lots; no banking. Optional explicitly quoted expiry alignment preserves existing paid time. |
| Separate Sneat Space SSO service | €99/month per Space | One independently paid service across Sneat products; no DataTug launch discount. |
| Enterprise | Custom/contact | No speculative self-service price or entitlements. |

Business with 7, 8 and 10 MAU and no prepay costs €99, €119 and €159 respectively. With ten prepaid units, 17 MAU are covered; 21 MAU adds €80 overage; 14 MAU leaves three unused units that do not carry forward. Annual Business base at ten monthly payments is €990; annual prepaid capacity is a distinct commitment. The founder confirmed base-only launch-discount scope for Pro and Business. Invitation treatment of multi-component Business invoices still needs proof; the existing largest-active-discount rule must not leak onto usage or prepay.

Without prepaid capacity, the same arithmetic gives 20 MAU = €359, 50 MAU = €959 and 100 MAU = €1,959. These are arithmetic examples, not claims of a live metered invoice.

Every plan retains unlimited owner-only private projects and unlimited public projects. Pro adds five protected projects with at most five authenticated users each, owner included. Sneat Space permissions govern cloud edits and GitHub permissions govern GitHub-hosted edits; a plan limit grants no access. A full refund immediately ends the affected Pro entitlement and paid claim, without ending another qualifying subscription or proving provider renewal cancellation. These rules are not replaced by Business pricing.

When any plan ends for any reason, affected projects become read-only and cannot use hosted AI; the user may still run queries locally. Retain projects and data; never automatically make a private or protected project public. The precise affected-project cohort during a paid-to-Free fallback still needs an authoritative ownership mapping before enforcement.

### Proposed billable-MAU contract (not a founder quotation or shipped hook)

Use a UTC half-open monthly **subscription-anchored** window from the billed item, including when base billing is annual. For DataTug Business, one authenticated Sneat user counts once per verified payer space and period after a successful, user-initiated hosted query/exploration with a server result, committed project-resource mutation, or completed hosted-AI request including server-delivered cached completion. Exclude mere sign-in, membership, invitations, notifications, passive opens, denied/failed/partial operations, replays/no-ops, local-only and GitHub-side work. Verify the payer space against the standalone project's grant: a member of two spaces is not automatically billed twice, but may count for two payers after separately qualifying work. Current endpoints have not yet proved this complete activity contract.

DataTug's adapter emits stable idempotent operation identity and non-sensitive explanatory evidence. A product-neutral Paymentus/shared-Sneat owner holds the payer/user/period dedup ledger, capacity arithmetic, quote-component boundaries, corrections and invoice handoff. Each product supplies a reviewed versioned configuration snapshot: currency, validated monthly and annual base amounts, included monthly units, overage unit amount, annual capacity unit amount, and component discount eligibility. Do not hardcode universal two-month or 20% formulas in shared logic. Snapshot the config version on a frozen quote and accepted period; a later deploy cannot rewrite historical prices. Other Sneat products can use their own qualifying-activity adapters. Cross-product payer/config/event isolation is required. No admin UI is needed; config changes are reviewed, committed and deployed.

### Proposed prepaid lifecycle and provider boundary

The founder confirmed that a customer can buy extra capacity at any time and combine purchases. Each purchase has its own twelve-month term after confirmed payment: ten extra units followed four months later by five extra units gives fifteen extra units for eight months, then five extra units for four months. The seven base-included MAU remain separate. These are capacity rights that refresh monthly, not money credits consumed once.

The founder also proposed an optional extension to align an existing lot to a new lot's expiry. Offer this as a separate explicit quote, retaining purchase lineage: in the example, the existing ten units retain their eight already-paid months; adding four months at the stated €16 monthly equivalent costs €640, and the new five-unit year costs €960, for €1,600 due. The customer can instead keep independent expiry dates. Alignment MUST NOT silently reprice, reset or discard paid time. This optional design direction is recorded for review; no provider flow or exact partial-month policy is thereby proved.

Payment verification establishes the purchase start, without imposing a next-month activation delay. A twelve-month lot bought mid-window can intersect thirteen subscription windows. The plan explicitly recommends a treatment for those partial windows and exact extension-date arithmetic, separately from the decided independent-lot purchase rule. Renewal, decrease, base cancellation and component refunds also have concrete labeled recommendations in that plan.

Current Stripe mixed-interval subscriptions can be an engineering option for annual base/prepay plus monthly usage, but Checkout Sessions cannot create that mixed subscription directly. This is not a working Business checkout. A reviewed Paymentus create/reconcile path, authoritative activity ledger, invoice/correction fixtures and founder-run provider tests gate activation. Preserve Pro's single-item Checkout. Enterprise examples are possibilities, not shipped features.

The founder also decided on **€99/month SSO for a Sneat Space**, reusing existing Sneat OIDC/SAML capability. The Space can enable it during Business signup or later from an **Enable SSO** setting in Business/Space settings. The SSO entitlement and single recurring charge belong to the Space across Sneat products, not to a DataTug product subscription; buying another Sneat product for the same Space must not charge SSO again. This requested end-user setting differs from an unrequested billing-configuration/admin UI. Existing SSO adapters are reuse evidence, not proof that this Space billing/enable journey or a DataTug Business purchase is deployed.

The Space-owned SSO service must not disappear merely because a DataTug plan changes or ends while another Sneat product still uses that Space. An existing Space-level service/subscription is the preferred composition unless reviewed ownership evidence proves an equivalent safe unification with product checkout.

SSO-only payment is not a DataTug base subscription: it cannot claim or retain a DataTug launch-discount place, or make its payer a DataTug active paying buyer. Conversely, a failed SSO payment or another product's payment must not revoke an independently paid DataTug Pro/Business entitlement. These are inferred architecture safeguards, not additional founder quotations.

## Acceptance Criteria

### AC: four-plan-page (verifies REQ:four-plan-offer)

**Given** the public pricing page **When** a visitor opens desktop and mobile views **Then** only Free, Pro, Business and Enterprise appear in the specified layout, with Business base/included/overage/prepaid amounts distinct, no Team/Company card and no claim that Business checkout works today.

### AC: pro-compatibility-and-currency (verifies REQ:pro-compatibility)

**Given** an existing EUR Pro quote/subscription **When** the revised EUR-only page appears **Then** the prior provider price identity, invoice and entitlement remain unchanged, and the €19 buy action is unavailable until its verified EUR quote matches the displayed period, amount and offer terms. A mismatch is refused before provider session creation; no USD estimate or selector appears.

### AC: ended-plan-access (verifies REQ:plan-end-access)

**Given** a verified ended plan and an affected cloud project **When** a user attempts a shared save, hosted AI and local query execution **Then** the shared save and hosted AI are denied, local query execution remains possible, and project data/private visibility are retained. No automatic project-cohort guess is permitted before ownership mapping is proved.

### AC: activity-qualification (verifies REQ:billable-mau)

**Given** one user who signs in, passively opens a result, completes one qualifying hosted operation, replays it and completes another for the same payer/period **When** activity is reconciled **Then** that user contributes one MAU with auditable evidence; denied, failed, partial, local and no-op events contribute none.

### AC: payer-and-period-isolation (verifies REQ:billable-mau)

**Given** a user in two payer spaces **When** qualifying work occurs for only one verified payer **Then** only that payer gains one MAU in the current UTC subscription-anchored period, and a new half-open period starts a new count.

### AC: business-math (verifies REQ:business-components)

**Given** Business with seven included MAU **When** a period closes with 7, 8 or 10 MAU and no prepaid units **Then** overage is €0, €20 or €60 respectively, on top of the €99 base, without a volume tier.

### AC: nonbankable-prepaid (verifies REQ:business-components)

**Given** ten prepaid units effective for a window **When** 14, 17 or 21 MAU qualify **Then** uncovered usage is 0, 0 or 4 units; unused committed capacity does not enlarge a later window.

### AC: base-only-discount (verifies REQ:discount-components)

**Given** a valid 30% launch grant on a Business invoice with base, overage and prepaid components **When** the amount is calculated **Then** only the eligible base is discounted; €20 usage and €192 prepaid unit prices stay intact. Monthly/annual base periods and invitations obey the same component boundary.

### AC: shared-config-isolation (verifies REQ:reusable-billing)

**Given** two Sneat products with different versioned currency, rates and qualifying-event adapters **When** identically named users emit activity **Then** payer ledgers and configurations remain separate; neither inherits DataTug's seven/€20/€192 constants or query trigger by default. An inexact fractional-minor-unit discount or arithmetic overflow refuses without a quote rather than truncating a payable amount.

### AC: one-space-sso-charge (verifies REQ:space-sso)

**Given** a Sneat Space using DataTug Business and another Sneat product **When** it enables SSO during signup or later in Space settings **Then** the Space has one €99 monthly SSO service/charge. A signup against an already SSO-enabled Space makes no duplicate purchase, and a DataTug plan change does not cancel SSO needed by the other product. An ordinary sign-in alone creates no SSO purchase; SSO-only payment does not consume a DataTug launch place or active-buyer count, and an SSO/other-product payment failure does not end an independently paid DataTug plan.

### AC: sso-cancellation-paid-end (verifies REQ:space-sso)

**Given** a Space with verified paid SSO that is canceled for renewal **When** sign-in starts or an in-flight identity-provider callback or token exchange completes before the proved paid-through instant **Then** SSO remains available. At and after that instant those new sign-in paths refuse, while an independently authenticated current Space administrator can still access billing and configuration recovery. Existing application sessions are not automatically revoked. This behavior is a requirement; the payment proof, migration and recovery path remain implementation gates.

### AC: frozen-replay-and-correction (verifies REQ:reusable-billing)

**Given** a duplicate/late event, a changed product configuration and an accepted historical quote **When** reconciliation runs **Then** its original version is retained, charges are idempotent, a correction is explainable, and unresolved provider/finalization evidence refuses rather than guesses.

## Not Doing / Out of Scope

- This Draft does not by itself arm live Business checkout, create live provider prices, or claim a verified MAU invoice; fake-backed implementation may proceed within the founder's authorized scope.
- No new Team, Startup or Company self-service tier; historical subscribers need a separate reviewed migration.
- No Enterprise implementation, SCIM or generic billing-configuration admin UI; the requested Space-level Enable SSO setting remains in scope.
- No inherited Business AI allowance: the former 6,000 questions belonged to the superseded $199 tier.
- No client telemetry as billable authority or charge to every project grant.

## Open Questions

- **OQ: prepaid-edge-and-service-lifecycle.** Independently additive twelve-month lots bought at any time are decided; optional explicit alignment is proposed by the founder. Review the plan's partial-window coverage, exact alignment-pricing, renewal/decrease, failed-payment, cancellation and component-refund recommendations before arming charges; common-anniversary or next-month activation is not the default.
- **OQ: activity-trigger-proof.** Which hosted endpoints emit complete server-side qualifying-operation receipts, including cached AI and committed mutations? The proposed semantic set needs implementation/contract review.
- **OQ: invoice-close-and-corrections.** Prove event cutoff, meter timing, duplicate/late correction and credit-note flow against test-mode provider evidence before Business billing opens.
- **OQ: Business entitlements.** The former $199 Business contributor cap, 6,000-AI allowance and protected-project assumptions are obsolete. Define the €99 plan's included AI, protected access and real hosted sharing availability before paid checkout.
- **OQ: invitations.** Prove time-limited invitation transitions against Business base-only component scope on monthly and annual invoices.
- **OQ: SSO integration proof.** The €99 monthly Space charge, cross-product entitlement and cancellation-through-paid-end behavior are decided. Verify OIDC/SAML reuse, the during-signup/later-enable flows, one-charge-per-Space behavior, paid-through enforcement and administrator recovery; full-refund and dunning behavior, existing active-config migration and provider/infrastructure cost remain open before arming. Provider sign-in MAU is distinct from DataTug billable MAU.
- **OQ: Business component refunds.** A refund of only a usage charge or prepaid capacity purchase must not blindly apply Pro's whole-subscription terminal-refund rule. Define which complete Business base refund ends entitlement and releases a buyer's active claim, and what cancellation does to remaining prepaid capacity when mixed provider items end.
- **OQ: tax/price migration.** Prove EUR new-offer tax/legal terms and any deliberate migration of existing EUR Pro provider prices or terms; matching currency and numeric display do not perform that migration.

---
*This document follows the https://specscore.md/feature-specification*
