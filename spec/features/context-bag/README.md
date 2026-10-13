---
format: https://specscore.md/feature-specification
status: Implementing
---

# Feature: Session Context Bags

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/context-bag?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/context-bag?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/context-bag?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/context-bag?op=request-change) |
**Status:** Implementing
**Source Ideas:** —

## Summary

Session-owned Context Bags store typed, environment-scoped working context in OVDB, with
explicit Apply, session-wide workspace binding and reproducible chat inputs. M0 selects
the local store profile from evidence; hosted sessions use the existing OVDB Firestore mount.

## Problem

An investigator should enter Order 123 once, follow its lines and invoices, and keep
the evidence. Today typed context, chat sessions and Detective Cases have separate
lifetimes. Scaling to 300 suspect orders should introduce cohorts, not a new product.
Context Bag is the reusable working context of a session; it is not a query language,
a tag list, or an incident workflow.

Owner approval recorded 2026-10-13 15:46 UTC authorizes the first implementation phase,
M0–M3. The approved real-app layout uses top Ionic `ion-tabs`, `ion-tab-bar` and
`ion-tab-button` session navigation, with a hybrid bag summary and drawer and Ionic
`ion-item` headers and rows. Ionic with AG Grid is the baseline; PrimeNG is not part of
this phase. Local and hosted OVDB profiles are included. M0 selects the store profile
from evidence; inGitDB remains a hypothesis until then. M4 cohort execution, reviewed
extraction, snapshot branching and direct annotations are deferred with the reasons in
the implementation plan. This approval does not decide a sync protocol or authorize
publishing, private-data migration or deployment.

## Behavior

The proposed contract is detailed in [Product and UX](product-ux.md),
[Operational architecture](architecture.md), and [Delivery and reconciliation](delivery.md).
[Discovery evidence](../../research/context-bag/inventory.md) separates main, WIP and gaps.
The [interactive review](../../research/context-bag/review.html) compares four placements
using the same synthetic investigation. It is a disposable prototype, not app code.

### REQ: session-owned-context

Each durable session MUST own one or more bags and an explicit OVDB database identity.
Each bag MUST be associated with an environment and apply only to that environment.
This multi-bag/environment model was explicitly approved by the owner in this chat. Multiple
sessions MUST remain independently resumable. Start/stop/archive/delete MUST be
distinct from focus, running jobs and Detective Case closure. Switching sessions MUST
not transfer values, drafts, permissions or late query responses.

### REQ: explicit-application

Entry checkboxes and the master checkbox MUST edit a draft. Apply Context MUST create
an immutable applied selection referencing exact entry revisions. Compatible Explorer
views use that selection; saved queries and Checks retain explicit Run. Chat history
MUST remain immutable. Disabling context MUST preserve entry selections.

### REQ: multiple-environment-bags

A session MUST begin with one default bag and support additional named bags. Each bag
has immutable environment/identity scope; cross-environment copy creates new identity
and requires review. Apply updates only that environment's applied-bag slot. Views MUST
match environment before considering entries. Compare pins two independent manifests
and executes each side under its own authorization; it MUST NOT union raw IDs across
environments. A visible Clear context button plus confirmation MUST clear the selected
bag in two clicks, without a third Apply and without deleting historical evidence.

### REQ: parameterized-dashboards

Applied context MUST be session-wide across Explorer, queries, dashboards and Checks.
Dashboard definitions MUST declare ordinary inputs and map them to widget-query inputs.
A dashboard run MUST resolve one pinned input manifest, with separate widget receipts.
A cohort MUST NOT silently bind its first member to a scalar input. At least one real
chart and a result grid MUST be included in the first order dashboard. Apply marks an
existing run stale; explicit Refresh creates a new run. Missing or failed widgets MUST
remain visibly missing/failed, never silently appear as zero.

### REQ: portable-bindings

Queries MUST declare ordinary typed parameters, never bag addresses. Runtime binding
MUST validate semantic identity, scope, cardinality, types and a supported relationship
path; ambiguities require a choice. A large set MUST reach the executor by an authorized
collection handle rather than through LLM prompt text. Empty sets MUST NOT mean all rows.

### REQ: durable-revisions

Stable IDs MUST be independent of names and contents. Composite keys and new value
kinds MUST fit versioned envelopes. Send-time attachments MUST pin revisions.
Operational changes MUST preserve version checks, idempotent receipts and metadata-only
audit through existing authorized domain boundaries. OVDB engine capability MUST be
demonstrated before relying on transactional or offline guarantees.

### REQ: privacy-and-knowledge

Extraction MUST review candidates before saving; retaining source text is a separate
choice, initially Always ask. Entity annotations MUST distinguish their target from
their originating Case/session, respect current access and identity domain, and keep
operational observations distinct from authoritative MeaningGraph definitions.

## Acceptance Criteria

### AC: one-order (verifies REQ:session-owned-context, REQ:explicit-application)

Create a session, add typed Order.ID 123, Apply, then open Orders and OrderDetails.
The relevant record and related rows appear with visible binding provenance. Opening
OrderDetails needs no second context action. A pending checkbox change changes neither
view; Apply commits it. A declared ambiguous relationship requires selection.

### AC: isolated-sessions (verifies REQ:session-owned-context)

Start two sessions with different bags and composer drafts. Switch while a query runs;
the result remains in its originating session. Stop prevents new work, preserves state,
and reports cancellation uncertainty. Archive is reversible and read-only. Delete shows
referenced evidence impact. None of these actions closes or resolves a Case/incident.

### AC: large-cohort-history (verifies REQ:portable-bindings, REQ:durable-revisions)

Attach 5,300 Order IDs, send a related-row count question, then edit the cohort and Apply.
The prompt contains metadata and an authorized tool handle, not 5,300 values. Execution
uses every member under bounded admission. The original turn still names its original
revision and result. Unsupported engines fail visibly rather than truncate.

### AC: environment-bags-and-clear (verifies REQ:multiple-environment-bags)

Create UAT affected and PROD affected bags in one session. Applying UAT changes only
UAT consumers; applying PROD leaves UAT applied state intact. Compare shows independently
pinned sides and never treats UAT Order 123 as PROD Order 123. Clear context → Confirm
empties only the selected bag, removes its active environment binding if any and its
unsent attachments, preserving other bags, sent turns and saved dashboard runs.

### AC: order-dashboard (verifies REQ:parameterized-dashboards)

Given applied Order.ID 123, opening Order overview prefills its required scalar input.
Run produces a details grid and at least one chart, all using one pinned input manifest
and individually recorded widget outcomes. A 300-order cohort requires explicit scalar
selection or a declared cohort dashboard. Apply another ID: previous widgets remain
labeled with their old run until Refresh. A failed category query shows an error rather
than a zero-valued pie chart. Master disable suspends automatic binding in every consumer.

### AC: operational-parity (verifies REQ:durable-revisions)

On the chosen local OVDB engine and separately hosted OVDB, retry a timed-out command,
race two updates, revoke access, and restart during commit. Verify no lost update,
duplicate semantic event, half-published revision or fabricated success. Reopen the same
hosted session from another device and the terminal, with no local CLI on the web device.

### AC: privacy-review (verifies REQ:privacy-and-knowledge)

Paste synthetic email with mixed IDs. Review and correct proposed types, decline source
retention, then save context. Raw text is absent from persisted prompts/audit/logs.
Show external AI disclosure before transmission. Historic entity discovery distinguishes
direct notes from collection membership and never reveals denied counts or snippets.

## Open Questions

The approved layout and first-phase boundary are recorded in
[Delivery](delivery.md#owner-review). M0 selects the OVDB engine profile from evidence;
cohort combination and offline-copy behavior remain outside M0–M3. Existing Investigation
Context and Annotation specifications are not silently rewritten or declared superseded.

---
*This document follows the https://specscore.md/feature-specification*
