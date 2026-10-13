---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Detective and reusable Checks across web and terminal

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Detective is DataTug's investigation mode: follow the evidence, navigate uncertainty,
and find the cause. A persistent Case carries the context, evidence, findings and
history of an investigation. Checks are independently useful, reusable validation
workflows; Validators are their atomic query/result assertions.

Both Detective and Checks/Validators are required in the `datatug-apps` web UI and
the `datatug-cli` interactive terminal UI. CLI commands and JSON output supplement
the terminal UI; they do not replace it. Layouts adapt to their surface while the
underlying definitions, semantics, permissions and stored history stay consistent.

**Source context:** the owner's Detective + Checks brief, including its idea/vibe
addendum and the explicit web/terminal clarification on 2026-10-13. The owner then
requested autonomous work. This document captures the product requirements; a
scaffold, component, focused test or pending branch is not delivery evidence.

**Plan:** [Central implementation plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/detective-checks/README.md).

## User journey

1. Open a project and choose Detective or Checks in either UI. The catalogs load,
   or the UI explains the missing agent/store permission without a fake empty list.
2. Create a Case from a question, selected query values, pasted ticket or Incidentius
   reference. Review extracted proposals; save confirmed typed context. Five order IDs
   and seven invoice IDs remain distinct lists, not one concatenated string.
3. See compatible queries, Checks and Scenarios with the matching/missing inputs
   explained. Investigate freely or apply a Scenario without conflating it with a Check.
4. Inspect a Check's globally reusable dependencies, select a logical environment,
   review bound inputs and manually run. A shared QA/UAT source requires an explicit
   choice. No query executes merely because the page opens or text is pasted.
5. Inspect execution status, validation outcome and severity separately, down to
   nested queries/Checks and per-row Validators. Missing CRM evidence stays unknown
   or failed according to the selected policy; it never becomes a clean green pass.
6. Keep a finding and its evidence references in the Case, optionally link an incident
   or propose knowledge for human review. Close with a reason, then reopen if necessary.
7. Open the same project/store in the other UI. The saved Case, confirmed context,
   history and Check run references are present after refresh/restart without another
   action. Local-only stores disclose their machine/reachability limits.

## Problem

Existing investigation-context and dataset-value validation pieces do not provide
this complete workflow. The web app has no persistent Detective Case or reusable
Check catalog. The existing Checks specification describes an earlier single-query
form with diagnostic and CI behavior. This feature adds composite positive-expectation
manual Checks and distinct Case workflows; it preserves legacy dataset checks and
does not turn earlier roadmap descriptions into shipped capabilities.

Cases and Incidentius incidents solve different problems. A Case may have no incident,
one incident may have several independent Cases, and a Case may link several incidents.
Scenarios guide investigation; Checks validate expectations; neither owns the other.

## Behavior

### REQ: shared-surfaces-and-store

Both UIs MUST expose the journeys and the capability matrix below. Shared server-side
contracts and services own execution, lifecycle transitions and persistence. The web
client MUST NOT write Firestore directly or reimplement the Check evaluator. Terminal
work MUST reuse existing Bubble Tea/navigation conventions and preserve normal typing
for search and editing. Every create flow lands on the created resource; cold links
reload from its stable ID, never rely on router state alone.

The initial local agent adapter MUST persist definitions as reviewable project assets
and private Cases/runs/audit outside tracked project content by default. Both UIs using
the same agent/store MUST resume the same records. Cross-device resume is conditional
on an authenticated reachable shared store; a loopback-only agent MUST NOT imply
cloud or cross-device availability. Unsupported hosted storage is explicit and is
not silently replaced with browser-local persistence.

### REQ: stable-identity-and-versions

New resources MUST use human-readable slug IDs with uniqueness checks and safe path
validation. Rename MUST preserve the referenced ID unless references are migrated
explicitly. Edits MUST use optimistic versions and refuse stale writes. Legacy IDs
and dataset `options`/`regexp`/`sql` checks MUST remain valid without bulk migration.

### REQ: case-lifecycle-and-evidence

A Case MUST have a title, typed context, status `active|closed|archived`, findings,
evidence/run references, incident links and ordered activity. Close requires a reason;
reopening is a transition/event, preserving previous closure reasons. Meaningful
context, lifecycle, query, Check, finding, link and reviewed AI actions MUST be recorded.
Updates MUST survive process restart and concurrent access without losing history.

Findings MUST distinguish hypothesis from confirmed observation. Proposals to
MeaningGraph/ModelSpec MUST retain provenance and human review; uncertain findings
MUST NOT silently become authoritative knowledge.

### REQ: typed-context-and-recommendations

Context and Check inputs MUST represent named typed scalars, lists and date ranges.
Reuse the canonical investigation scalar value encoding. Proposed values MUST be
editable/confirmable; ambiguous extraction MUST show its evidence and uncertainty.
Check-constrained extraction uses the declared input schema. Broad intake ranks
compatible Checks/Scenarios with an explanation and missing inputs. Case scope
changes MUST NOT carry values or late responses into another project/principal.

### REQ: global-checks-local-validators

A Check MUST have a stable ID, description, inputs, version, provenance and one or
more saved-query or Check-reference steps. Queries use the existing library rather
than duplicating text. Nested Checks MUST remain independently discoverable/runnable
and show dependency and used-by relationships. Reject cycles and missing references.

Validators are local to a query definition by default. They MUST support field/row,
whole-set, metadata and reusable Check-reference assertions. A user MAY explicitly
promote a Validator into a global reusable catalog; references and used-by relations
MUST survive edits. Historical runs MUST retain the evaluated graph/definition versions
and sufficient snapshots to interpret outcomes after catalog changes.

### REQ: bounded-authorized-execution

All Check query execution MUST use the existing authorized read-only query path,
including cross-source OVDB/DALgo/TugQL execution and metadata queries. Project-level
permissions are reused; no new per-resource ACL scheme. Missing ability to authorize
a resource MUST fail closed. Database reads remain constrained by supported execution
profiles; HTTP method/trust limitations MUST be explicit, without claiming an opaque
endpoint cannot have side effects.

Nested/per-row work MUST have cycle/depth/fanout/result/time bounds, cancellation,
bounded concurrency and reuse of repeated compatible reads. The UI MUST disclose
estimated fanout or an execution limit rather than silently launch unbounded N+1 work.
The MVP performs manual execution only, with history and drilldown.

### REQ: status-outcome-severity

Execution status is `pending|running|completed|error|skipped`; validation outcome is
`pass|fail|unknown`; severity is `critical|error|warning|info`. Unavailable dependency
policy is `block|fail|warn|skip`. Preserve the underlying execution error in all cases.
Required incomplete evidence MUST prevent a parent `pass`; warning/skipping MUST
remain visible. Zero assertions or omitted children MUST NOT manufacture healthy data.

### REQ: explicit-environment

The user MUST select the current or an explicit logical environment. A query result
may supply it only when unambiguous. A physical source shared by QA and UAT MUST NOT
choose one silently. Carry the selected environment to every child and validate source
mapping. Unsupported/ambiguous mappings MUST have diagnostic errors before execution.

### REQ: audit-and-privacy

Audit is project infrastructure with a pluggable append interface and an initial local
JSONL implementation. Store actor, scope, time, action, source/environment IDs, query
identity/fingerprint, duration/row counts, outcome, run/definition/correlation references
and redacted input metadata. Do not dump rows, pasted text, credentials or parameter
values into audit by default. Private persisted Case context follows explicit storage
and access rules. Audit write failure MUST be visible; retry/durability semantics
MUST NOT silently discard it. Runs and Case activity remain separate linked records.

Retention/export/redaction and secret safeguards MUST be documented. Future object
store adapters batch/partition objects; they MUST NOT assume atomic appends to S3/GCS.
Recorded provenance explains an original run; it does not promise deterministic replay
of external systems.

### REQ: scenarios-ai-and-integrations

Reusable Scenarios MUST declare compatible context/inputs and manual query/Check steps,
allow selection after a Case begins, and retain their applied version/progress. Cases
MUST support linking/creating Incidentius incidents and starting independent Cases from
an incident with reviewed safe context. Avoid one-to-one ownership or lifecycle coupling.

Meaningful UI actions SHOULD be expressible through the existing AI action harness:
context extraction, discovery, query/Check/Validator drafts, Case edits, manual runs and
failure explanation. Draft persistence/execution requires review when appropriate;
the AI uses the same deterministic authorization/action routes and MUST NOT bypass
permissions. Do not invent a parallel AI framework or label heuristic parsing as AI.

### REQ: product-story-and-deferred-scope

Dedicated Detective and Checks/Validators feature pages MUST use the existing DataTug.app
feature structure, navigation and accessibility/SEO patterns, explain the user problem
and actual shipped journey, and distinguish illustrations from product screenshots.
The marine investigator identity is calm, technical and evidence-led: fog/uncertainty,
charts/navigation, rough seas during incidents, safe harbor after understanding. Avoid
police noir and gamification.

Deferred: scheduling, CI/deployment/event-trigger execution, automatic closure/polling,
negative/diagnostic expectations, threshold UX, complex ACLs and autonomous historical
Scenario optimization. Preserve future resolution criteria (feed arrived/invoice paid),
ownership, evidence, notifications, human confirmation and reopen/sync considerations
as design topics. Do not implement those engines in the MVP.

### Capability matrix

| Journey | Web UI | Terminal UI | Shared behavior |
|---|---|---|---|
| Standalone Checks/catalog/editor | Required | Required | Same IDs, versions, queries/dependencies |
| Local/promoted Validators and query-row entry | Required | Required | Same assertions and promotion rules |
| Environment/typed inputs/manual run | Required | Required | Same authorized executor and bounds |
| Nested/row/set results and run history | Required | Required | Same immutable run tree/provenance |
| Case list/create/context/evidence/findings | Required | Required | Same persistent Case store |
| Close/reopen/archive/timeline | Required | Required | Same transitions and audit |
| Paste review/AI actions/Scenario discovery | Required | Required | Same confirmation and permission rules |
| Incidentius/knowledge linkage | Required | Required | Separate resources, reviewed provenance |
| Cross-UI resume | Required | Required | Same reachable authorized store |

### Transport and implementation boundary

The initial adapter exposes `GET /datatug/detective/state?project=<id>` with the
existing security-context mechanism, and `POST /datatug/detective/action?project=<id>`
with a typed discriminated action request. State/action payloads MUST be bounded,
validated and scoped to the served project/principal. Actor and timestamps are
server-issued. The shared wire model lives with the agent implementation; browser
and terminal clients MUST exercise the same serialized contract in integration tests.
Resource capability errors MUST be rendered distinctly from an empty catalog.

## Acceptance Criteria

### AC: checks-both-surfaces (verifies REQ:shared-surfaces-and-store)

**Given** a served writable project, **When** a user opens either UI, creates/edits
and manually runs `validate-invoice`, **Then** the saved definition, dependency tree,
version and run appear in both UIs, including after cold reload.

### AC: nested-invoice (verifies REQ:global-checks-local-validators)

**Given** `validate-invoice` references `validate-invoice-lines` and `validate-crm-customer`,
**When** an invoice is validated across its sources, **Then** non-null customer ID,
invoice/line total equality and CRM existence have separate results and nested provenance.

### AC: row-and-set-results (verifies REQ:bounded-authorized-execution)

**Given** a query returning ten invoices with a per-row `validate-invoice` Validator
and whole-set row-count/aggregate assertions, **When** manually run from either UI,
**Then** each row is inspectable/filterable, set results stay distinct, and bounded
fanout/repeated-read reuse is verified without bypassing the executor.

### AC: paste-lists-and-ranges (verifies REQ:typed-context-and-recommendations)

**Given** a ticket naming five orders, seven invoices and a date range, **When** pasted
in either UI, **Then** proposals retain types/cardinality/evidence, require confirmation,
and explain compatible reconciliation Checks/Scenarios without automatic execution.

### AC: environment-ambiguity (verifies REQ:explicit-environment)

**Given** one physical source mapped to QA and UAT, **When** a user launches a Check
from an unscoped query result, **Then** no source executes until logical scope is
resolved and every child receives the chosen environment.

### AC: unavailable-evidence (verifies REQ:status-outcome-severity)

**Given** an unavailable CRM dependency, **When** each unavailable policy is exercised,
**Then** both UIs retain the execution error, unknown evidence and policy effect,
and the parent cannot claim that every required assertion passed.

### AC: metadata-and-access (verifies REQ:bounded-authorized-execution)

**Given** schema metadata and ModelSpec expectations plus a restricted principal,
**When** a Check compares them, **Then** metadata uses the query machinery and all
denied/unsupported reads/writes fail closed without exposing hidden Case/input data.

### AC: case-lifecycle-and-resume (verifies REQ:case-lifecycle-and-evidence)

**Given** multiple Cases linked to one Incidentius incident, **When** one is closed,
reopened and opened in the other UI after server restart, **Then** its context, closure
and reopen history, evidence/run references and independent lifecycle are intact.

### AC: catalogs-and-history (verifies REQ:stable-identity-and-versions)

**Given** reused nested Checks and a promoted Validator, **When** definitions change,
**Then** used-by links remain correct, cycles/stale edits/collisions are rejected,
and old runs retain every evaluated dependency version/snapshot.

### AC: reviewed-ai-and-knowledge (verifies REQ:scenarios-ai-and-integrations)

**Given** a request to validate invoice totals, **When** AI drafts a query/Validators
or proposes a shared discovery, **Then** both UIs allow review/edit before save/run/
promotion, record confirmation and provenance, and enforce the ordinary action permissions.

### AC: audit-without-data-dumps (verifies REQ:audit-and-privacy)

**Given** query/Check runs and Case changes, **When** inspecting/restarting the local
JSONL audit backend or injecting a write failure, **Then** events are durable, bounded,
linked and redacted with visible failures and no result-set/secret/paste dumps.

### AC: complete-user-journey (verifies REQ:shared-surfaces-and-store)

**Given** a real demo project and agent, **When** automated web and terminal journeys
walk create → context confirmation → manual nested Check → drilldown → evidence →
close/reopen → resume from the other UI, **Then** no hidden setup/reload/tap is needed
and the test asserts persisted mechanisms and authorized execution, not only labels.

### AC: feature-pages-and-honest-scope (verifies REQ:product-story-and-deferred-scope)

**Given** the landed MVP, **When** visiting the linked Detective and Checks feature
pages on desktop/mobile, **Then** the actual problem/workflow is clear, navigation and
accessibility checks pass, and future capabilities are explicitly future.

## Open Questions

No new owner decision is required for the local-agent implementation direction.
Hosted adapters, AI/provider availability and integration acceptance remain explicit
engineering verification items; they do not reduce the required two-UI workflows.

---
*This document follows the https://specscore.md/feature-specification*
