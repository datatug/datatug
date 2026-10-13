---
format: https://specscore.md/feature-specification
status: Implementing
---

# Feature: Detective and reusable Checks across web and terminal

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/detective-checks?op=request-change) |
**Status:** Implementing
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
approved implementation and required local and cloud storage from the ground up,
including a full cloud demo without starting a local DataTug CLI. This document captures the product requirements; a
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

### Interaction design

The Case is an evidence timeline with an always-reachable context editor, rather than
an empty dashboard. Its header answers: what question is being investigated, which
logical environment is selected, what is known, and what remains uncertain. A finding
may be a hypothesis even when a Check passed; closing a Case records a human decision
and never implies that linked incidents resolved. Use restrained chart/navigation
accents and plain labels such as “Evidence”, “Needs input” and “Unknown”.

The web layout offers a Case list, timeline/results area and context/evidence inspector.
The terminal adapts those into keyboard-driven list/detail/inspector views with clear
focus and back navigation. Both expose next actions adjacent to the evidence that
motivates them, retain unsaved drafts across inspection, and announce busy/errors.
Small displays must retain a usable editor and a path to every result, not clip forms.

The Checks detail page shows inputs, dependencies/used-by, Validators and history.
Before Run, show selected environment, confirmed bindings, affected sources and work
limits. Results show a summary plus separate row and set sections. Users can filter
failed/unknown rows and open a nested path while retaining the query execution identity.
A result badge must never depend on color alone. JSON may be an advanced editor but
must not be the sole way to author typed inputs, Validators or ordinary Case context.

## Problem

Existing investigation-context and dataset-value validation pieces do not provide
this complete workflow. The web app has no persistent Detective Case or reusable
Check catalog. The existing Checks specification describes an earlier single-query
form with diagnostic and CI behavior. This feature adds composite positive-expectation
manual Checks and distinct Case workflows; it preserves legacy dataset checks and
does not turn earlier roadmap descriptions into shipped capabilities.

For this delivery, this Feature supersedes conflicting planned requirements in
[Executable checks](../checks/README.md): a Check may compose multiple saved queries
and global Check references; diagnostic/negative expectations and CI triggering are
future; outcomes use the three independent dimensions below; AI drafting is MVP;
and the shared Detective transport is the initial adapter. The old single-query
`CheckSpec` and proposed fourth dataset-check kind are not a second mandatory engine.
Existing shipped dataset validation and its CI runner retain their behavior. Earlier
safety requirements still apply, especially query-policy limitations, reviewable
assets, human review, execution evidence and denied-resource non-disclosure.

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

Local and cloud stores are both required MVP adapters over common domain contracts,
execution semantics and lifecycle services. The local adapter MUST persist definitions
as reviewable project assets and private Cases/runs/audit outside tracked project
content by default. The cloud adapter MUST persist definitions, private Case context,
findings, Scenario applications, runs and audit through the existing authorized backend.
Both UIs using the same project/store MUST resume the same records. An authenticated
user MUST also resume cloud records from another device without a local agent. Storage
identity, access and reachability MUST remain explicit; no implicit machine-local or
browser-local fallback may discard or fork hosted records.

### REQ: hosted-demo-without-local-cli

A user MUST be able to open a cloud demo project in the web app and exercise the full
MVP without installing or starting `datatug-cli`, opening a loopback port, or connecting
to the demonstrator's machine. Authoritative hosted Check runs MUST use a server-owned
saved-query executor; browser exploration results are not trusted run evidence. The
hosted executor MUST preserve the same authorization, read-only limits, environment
selection, nested/per-row/set-level evaluation and evidence semantics as local execution.
Cloud storage alone does not satisfy this requirement. The server resolves saved-query
IDs and revisions, schema metadata, logical environment and admitted source mappings;
clients cannot submit arbitrary source URLs or fabricated rows/outcomes as a Check
receipt. Preserve source-policy restrictions and time/row/byte limits in evidence.
An ordinary browser DemoDB query is useful exploration, but acceptance requires an
immutable server-issued run referencing its actual query/dependency revisions.

Reusable demo definitions and datasets may be shared as templates. Mutated definitions,
Cases, pasted context, findings and run history MUST have an authorized user/workspace
scope, so independent visitors cannot read or overwrite one another's private evidence.
Use existing authentication and Personal Space/project-access authority. On the first
persistent action, create or resume a private demo project from a pinned public template,
with a distinct server-issued demo admission permitting only admitted demo sources.
No paid purchase is required, and paid project admission rules remain enforced. Sign-in
may occur at submit; guest provisioning must be explicitly implemented if offered.
Personal Space identity and demo admission MUST be resolved by the server, never from
client assertions of ownership, entitlement or demo mode. Admission pins the template
revision and allowed source identities; requesting a different project/source MUST NOT
turn demo access into paid or arbitrary-source access. Retried provisioning returns the
same private project without recopying over visitor edits. Refresh/restart must preserve
saved progress; the UI MUST explain sign-in or capability requirements.
The demonstration includes authoring, AI proposal review, execution, evidence, Scenario
progress, Incidentius/knowledge linkage and Case lifecycle, not just a read-only tour.

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
Transitions are active → closed (reason required), closed → active (reopened), and
closed → archived. Archived Cases remain readable; editing/execution is disabled.
A second closure requires a new confirmation/reason, even if an earlier reason is
retained in history. Do not represent lifecycle events as client-asserted history.
Client notes have an explicit note kind; only the server emits execution, lifecycle
and confirmation receipts. Evidence references must resolve within authorized scope;
adding an existing run is an explicit link event, not a fabricated execution event.

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
and sufficient snapshots to interpret outcomes after catalog changes. This includes
transitive Check dependencies, promoted Validator versions and query-local Validators
for every entry path, including a run started from an ordinary query result. A nested
Check evaluates the same assertions regardless of its launch surface. A saved query's
canonical ID, immutable revision/fingerprint and execution receipt identify the query;
an editable query ID alone is insufficient historical provenance.

Validate the entire resolved graph before dispatch. Missing Check/Validator/query
references and invalid bindings are definition errors, not unavailable remote systems.
Cycles introduced by editing a promoted or query-local Validator are rejected too.
Draft authoring can expose unresolved references, but cannot save them as runnable.
Deletion/rename must refuse dangling references; title edits preserve resource IDs.

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
The MVP performs manual execution only, with history and drilldown. Apply a whole-run
deadline and cancellation checks between nodes, not only per-query timeouts. Reuse keys
include logical environment, effective authorized execution profile, query revision and
typed bindings. No cache or in-flight result may cross a principal/project boundary.

A row result is tied to its original execution receipt and row ordinal, with a safe
identity reference where available. Sorting/filtering cannot move badges to other rows.
Re-running a saved query creates a new result set and must not annotate an older grid
as though it had been validated. Run records retain no full rows by default; row-level
drilldown uses evidence references and policy-safe explanations.

### REQ: status-outcome-severity

Execution status is `pending|running|completed|error|skipped`; validation outcome is
`pass|fail|unknown`; severity is `critical|error|warning|info`. Unavailable dependency
policy is `block|fail|warn|skip`. Preserve the underlying execution error in all cases.
Required incomplete evidence MUST prevent a parent `pass`; warning/skipping MUST
remain visible. Zero assertions or omitted children MUST NOT manufacture healthy data.

Policy changes the handling decision, never the observed execution fact. Preserve the
underlying safe error code and unknown evidence separately from the policy's effective
outcome. Omitted policy means `block`. Block stops dependent/later work and explicitly
marks unattempted nodes skipped/unknown; fail contributes failure; warn contributes
unknown with warning severity; skip contributes unknown with a visible policy decision.
None can transform a failed attempted query into an assertion that it never executed.
The same semantics apply to nested Check steps and per-row Check Validators.

Parent outcome is fail if any effective child fails, otherwise unknown if any required
child is incomplete/unknown, otherwise pass only when at least one assertion ran.
Execution status reflects actual execution errors independently of validation failure.
Severity is the highest severity across all non-passing children, independent of child
order. Thus an error-severity failure followed by a critical unknown remains critical.

Query restrictions and truncation are part of the evidence contract. Preserve result
limitations in the run tree; row-count, uniqueness and aggregate assertions over a
policy-filtered or truncated set are unknown. Observed row failures may remain failures,
but absence of visible failures cannot establish a whole-population pass. Hidden/missing
columns are unavailable evidence. Reject invalid kind/scope combinations and malformed
recordsets before evaluation; never panic or treat a missing assertion as passing.

### REQ: explicit-environment

The user MUST select the current or an explicit logical environment. A query result
may supply it only when unambiguous. A physical source shared by QA and UAT MUST NOT
choose one silently. Carry the selected environment to every child and validate source
mapping. Unsupported/ambiguous mappings MUST have diagnostic errors before execution.

### REQ: audit-and-privacy

Audit is project infrastructure with shared event contracts and both local JSONL and
cloud transactional storage implementations. Store actor, scope, time, action, source/
environment IDs, query identity/fingerprint, duration/row counts, outcome, run/definition/
correlation references
and redacted input metadata. Do not dump rows, pasted text, credentials or parameter
values into audit by default. Private persisted Case context follows explicit storage
and access rules. Audit write failure MUST be visible; retry/durability semantics
MUST NOT silently discard it. Runs and Case activity remain separate linked records.
Local files require an explicit recovery protocol for partial resource/audit writes;
cloud resource changes, operation receipts and audit commit in one transaction. A failed
transaction publishes none of them. If commit acknowledgement is lost, report uncertainty
and reconcile by operation identity. Query billing/activity counters are not substitutes
for execution evidence or the immutable Detective audit.

Retention/export/redaction and secret safeguards MUST be documented. Future object
store adapters batch/partition objects; they MUST NOT assume atomic appends to S3/GCS.
Recorded provenance explains an original run; it does not promise deterministic replay
of external systems.

### REQ: scenarios-ai-and-integrations

Reusable Scenarios MUST declare compatible context/inputs and manual query/Check steps,
allow selection after a Case begins, and retain their applied version/progress. Store
a separate application identity, Scenario revision/snapshot and ordered step receipts
on the Case; a bare Scenario ID cannot establish what guidance was followed. Applying
a Scenario offers its first manual action and never starts all queries automatically. Cases
MUST support linking/creating Incidentius incidents and starting independent Cases from
an incident with reviewed safe context. Avoid one-to-one ownership or lifecycle coupling.
Both stores require a real Incidentius adapter through its owning extension contract: create/link an incident from a Case and
start another independent Case from that incident. Authorize both resources independently;
a link does not confer access, reveal hidden context or synchronize lifecycle. Cross-store
or cross-extension writes use idempotent linked operations with visible partial progress
and reconciliation; do not promise a transaction across unrelated stores.

Meaningful UI actions SHOULD be expressible through the existing AI action harness:
context extraction, discovery, query/Check/Validator drafts, Case edits, manual runs and
failure explanation. Draft persistence/execution requires review when appropriate;
the AI uses the same deterministic authorization/action routes and MUST NOT bypass
permissions. Do not invent a parallel AI framework or label heuristic parsing as AI.
Extend the existing harness's typed action registry and dispatcher where DataTug actions
are absent, through the owning AI extension contracts and host wiring. A real provider
must support the hosted demo; a local browser-only interpreter is not cloud parity.
Treat pasted/source/model text as untrusted proposal content, with visible extraction
provenance and confirmation before deterministic mutations or manual execution.

Persist knowledge proposals as private, typed DataTug resources with the source Case,
finding/evidence references, hypothesis/confirmed distinction, intended MeaningGraph or
ModelSpec target identity/revision, and review status. Both UIs expose review/edit and
an explicit target-specific handoff/export or owning-extension submission with a receipt.
A proposal is not published knowledge. Promotion requires the destination's authorization
and review contract; neither public catalog overwrite nor arbitrary timeline text counts.
No private Case context is copied into a public knowledge catalog implicitly.

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
| Local and cloud storage | Required | Required | Same domain/version/history semantics |
| Full hosted demo without local CLI | Required | Cloud access required | Isolated persistent demo progress and hosted execution |

### Transport and implementation boundary

The local adapter exposes `GET /datatug/detective/state?project=<id>` with the
existing security-context mechanism, and `POST /datatug/detective/action?project=<id>`
with a typed discriminated action request. State/action payloads MUST be bounded,
validated and scoped to the served project/principal. Actor and timestamps are
server-issued. The hosted adapter uses existing authenticated project-reference routes
and server-side project authorization. Cloud mutations MUST go through Sneat HTTPS
endpoints; clients never write Firestore directly. Common domain types and evaluation
live in an importable shared module, with local filesystem/agent and cloud persistence/
execution supplied through interfaces. The shared implementation belongs in
`datatug-core`; the cloud adapter belongs in `datatug/backend`, composed by `sneat-go`.
The shared Core owns domain models, validation, evaluation and lifecycle semantics;
the CLI owns local filesystem/JSONL and execution/API adapters; the backend owns hosted
DALgo/Firestore persistence, admission, project-resource authorization and execution.
Core MUST NOT import the CLI, cloud SDKs or extension implementations. Cross-extension
reuse goes through owning contracts; `sneat-go` only composes and configures adapters.
Neither adapter owns a second evaluator or a divergent lifecycle. Browser and terminal
clients MUST exercise the same serialized domain contract and both adapters in integration tests.
Resource capability errors MUST be rendered distinctly from an empty catalog.
A query-specific grant does not grant access to private Case context or the complete
catalog. Protected deployment requires project-resource read/write/run authorization
before loading or returning these records; unrestricted local-owner proof is not proof
of that boundary. Both clients use the same capabilities and fail-closed behavior.
Cloud private resources remain logically owned by the admitted project in its Space.
Physically, all new private Detective records use the server-only extension namespace
`/spaces/{spaceID}/ext/datatug/detectiveProjects/{projectScopeKey}/{collection}/{id}`.
The server derives and verifies `projectScopeKey` from normalized store/project/branch
identity; a client-supplied key is never authority. Separate definitions, Cases, Case
events, runs/nodes, operation receipts and audit into bounded records/collections rather
than rewriting a growing project document. Existing project/query records retain their
paths. Do not create a parallel per-user root.
Use existing active project-contact/role authorization in addition to current admission;
Space membership alone is insufficient. Map read, edit and run actions to explicit
capabilities, rechecking current authority on every request/replay and execution boundary.
Demo admission is a distinct trusted source-restricted path, not an exception disabling
ordinary protected-project checks.
Firestore rules MUST disjointly exclude `detectiveProjects` and every descendant from
every broader client-read grant, preserving existing project/query read behavior;
a narrower deny cannot revoke an overlapping allow. Test direct client access as well
as HTTPS authorization, including Space members lacking project grants. The exclusions
must cover the complete protected Detective namespace, operation records and audit,
including nested descendants and future collections. Admission records retain their
existing private exclusions. Do not rely on parsing recursive wildcard paths to discover
which project descendants happen to contain private Detective data.

Cloud mutations commit scope/actor/payload-bound operation receipts, resource versions
and audit events atomically. A run is accepted durably before external execution, and
its result/evidence/Case linkage completes transactionally. External network execution
MUST remain outside retryable database transactions. Duplicate requests return the
same accepted operation/result rather than execute twice. Scope receipts by principal,
Space/project and action, and bind a canonical payload digest; reject changed payloads
within that identity and never disclose another principal's receipt. Replays reauthorize
current access, even when the original attempt succeeded. Concurrent acceptance requires
transactional uniqueness; workers claim accepted runs so retries do not dispatch them
again. A restart during execution exposes interrupted/uncertain execution with unknown
outcome and a recovery path, never invents success. Do not promise exactly-once external
reads after a crash: reconciliation or an explicitly authorized new attempt must preserve
attempt identity and prevent old workers from overwriting newer finalization. Finish
server evidence, audit, operation status and linked Case activity atomically.

State pages, catalogs and history need bounded reads and stable pagination or direct
resource lookup. Reaching a documented capacity must reject a new item before mutation,
not save an item that makes the next read of the entire project fail. A long-lived Case
must remain readable/closable as its activity grows. An uncertain local write or lost
commit acknowledgement returns an identifiable uncertain outcome plus a reconciliation path; blind retries must not duplicate runs or
lose their Case links. Local durability includes restart and failure-injection tests.

Representative definition and context (query IDs resolve in the project library):

```json
{
  "id": "validate-invoice", "title": "Validate invoice", "version": 0,
  "inputs": [{"name": "invoice_id", "type": "integer", "cardinality": "scalar", "required": true}],
  "steps": [
    {"id": "invoice", "kind": "query", "queryId": "invoice-detail",
     "bindings": [{"input": "invoice_id", "fromInput": "invoice_id"}],
     "validators": [{"id": "customer-present", "scope": "row", "kind": "non-null", "field": "CustomerId", "severity": "error"}]},
    {"id": "lines", "kind": "check", "checkId": "validate-invoice-lines",
     "bindings": [{"input": "invoice_id", "fromInput": "invoice_id"}], "onUnavailable": "block"}
  ]
}
```

`validate-invoice-lines` and `validate-crm-customer` are separately saved catalog assets.
The complete fixture adds the CRM step and an invoice-total/line-sum assertion. Five
orders and seven invoices use two named typed sets; a date range retains explicit start
and end values, with query bindings selecting each boundary. Extraction evidence and
confidence belong to the proposal, not to the confirmed scalar value representation.
Reviewable JSON assets retain stable IDs on title changes; new draft writes use version
0 and edits echo the last server version. A conflict reloads the latest resource and
preserves the user's draft for reconciliation.

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
and the parent cannot claim that every required assertion passed. Test explicit and
omitted policy, nested and per-row invocations, unattempted-node records, and permutations
of a failure plus a higher-severity unknown. Policy-filtered/truncated row-count and
aggregate results must remain unknown even when the visible subset meets the expectation.

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
and old runs retain every evaluated dependency version/snapshot. Repeat with a
query-local Validator invoking a Check whose own query-local Validator fails, plus a
grandchild and promoted Validator; entry paths must agree before and after edits.

### AC: reviewed-ai-and-knowledge (verifies REQ:scenarios-ai-and-integrations)

**Given** a request to validate invoice totals, **When** AI drafts a query/Validators
or proposes a shared discovery, **Then** both UIs allow review/edit before save/run/
promotion, record confirmation and provenance, and enforce the ordinary action permissions.

### AC: audit-without-data-dumps (verifies REQ:audit-and-privacy)

**Given** query/Check runs and Case changes on each adapter, **When** inspecting/restarting
its audit store or injecting a resource/audit/acknowledgement failure, **Then** events are
durable, bounded, linked and redacted with visible failures and no result-set/secret/paste
dumps. Cloud transaction rollback exposes no partial mutation; lost acknowledgements and
local partial writes reconcile by operation identity. Concurrent duplicates, revoked
replay access and worker interruption cannot fabricate or duplicate a completed run.

### AC: complete-user-journey (verifies REQ:shared-surfaces-and-store)

**Given** real demo projects served by the local adapter and separately by the hosted
backend, **When** automated web and terminal journeys on each adapter
walk create → context confirmation → manual nested Check → drilldown → evidence →
close/reopen → resume from the other UI, **Then** no hidden setup/reload/tap is needed
and the test asserts persisted mechanisms and authorized execution, not only labels.

### AC: local-and-cloud-store-parity (verifies REQ:shared-surfaces-and-store)

**Given** equivalent local and cloud projects, **When** web and terminal clients save,
edit and run the same definitions and resume Cases after restart, **Then** IDs, versions,
typed context, Scenario progress, evidence and lifecycle semantics agree. Cloud requests
authenticate and enforce project-resource access before reading or mutating records;
cross-user/project probes cannot disclose or alter another investigation. Include an
otherwise valid Space member without an active project assignment, revoked permission
on replay, and direct Firestore probes into nested Case/run/receipt/audit records.

### AC: full-cloud-demo (verifies REQ:hosted-demo-without-local-cli)

**Given** a fresh browser with no local DataTug process or agent connection, **When** a
visitor opens the cloud demo and follows the complete MVP journey, **Then** they can
author/edit Checks and Validators, confirm pasted context and AI drafts, manually run
nested invoice/CRM Checks, inspect per-row/set-level evidence, apply a Scenario, capture
findings and knowledge/Incidentius links, close/reopen the Case and resume saved progress.
Browse the public template, sign in at the first persistent action without a purchase,
and prove server-issued source-restricted demo admission without changing paid guards.
Use real admitted invoice and CRM sources, including the ten-row failure fixture, and
inspect authoritative server execution/audit receipts rather than client-reported results.
Prove a second authorized browser/device and the terminal resume the saved cloud Case
after service restart, and a different visitor cannot access it. Fail the test on any
loopback request, agent setup prompt, in-memory-only history, simulated result, or unimplemented required action.

### AC: feature-pages-and-honest-scope (verifies REQ:product-story-and-deferred-scope)

**Given** the landed MVP, **When** visiting the linked Detective and Checks feature
pages on desktop/mobile, **Then** the actual problem/workflow is clear, navigation and
accessibility checks pass, and future capabilities are explicitly future.

## Open Questions

The owner approved implementation on 2026-10-13 and explicitly required both local and
cloud storage from the ground up, with full cloud demos requiring no local CLI. This
supersedes the earlier local-agent-first delivery assumption. Hosted execution, existing
identity integration and real provider/Incidentius/knowledge acceptance are engineering
work within the approved scope; absence of an adapter does not defer cloud capability.

---
*This document follows the https://specscore.md/feature-specification*
