---
format: https://specscore.md/plan-specification
status: Approved
---
# Plan: Phase 2 — Library depth and provenance

**Status:** Approved
**Source Feature:** executable-knowledge-library
**Related Features:** provenance, meaning-inference
**Prerequisite Plans:** 2026-09-09-phase-1-core-investigation-loop
**Date:** 2026-09-09
**Owner:** alex
**Supersedes:** —

## Summary

Deliver a bounded executable-library quality journey after Phase 1 is accepted.
The [roadmap](../../ROADMAP.md) and owning Features define capture, provenance and
mapping review. The earlier plan bundled HTTP credentials, TUI/agent parity, rich types
and time-window semantics; those remain explicit deferred work below, not prerequisites
for this smaller useful outcome. This plan is not authority to start Phase 2 before
Phase 1 acceptance or to treat draft feature scope as implicitly approved implementation.

### Current implementation reservation

On 9 September the founder asked Codex to pick up something more substantial,
possibly from Phase 2, using a faster-model implementation agent. Codex selected the
storage foundation of task 2 as independent preparatory work: existing query-store
reuse, safe atomic persistence, revision conflicts, and reload/failure tests. This
selection was the implementation lead's scope decision, made before Phase 1
acceptance, not approval of every draft Feature. No new server write endpoint or
browser capture flow is activated by this reservation alone. Phase 1 is now
Implemented — Task 17 of its plan records the acceptance evidence — so this plan's
Phase 1 prerequisite is satisfied.

Task 2's Note records the responsible Codex task and shared WB claim. Maintain status
with SpecScore and publish changes; file status does not automatically follow agent
activity. Claim ownership before implementation and retain explicit predecessor gates
when reporting partial delivery. Storage tests cannot close the capture journey AC.

## Journey

**Start, nobody else acts:** an authorized local/Git project user opens the accepted
Phase 1 loop and gets the same useful data. Good result: no navigation, policy or
context regression.

**Capture:** after an ad-hoc lookup, choose Save as project query, review its name,
purpose, source and semantic parameters, and save. Good result: validated query assets
appear as an ordinary uncommitted Git diff, with no credentials or cached result rows.
A denied/failed save keeps the draft and reports the actual error.

**Understand:** open the captured query's provenance and a suggested mapping. The
user can inspect source/query/mapping origin, accept or correct a mapping explicitly,
and rerun the query. Good result: current values are bound through the same server;
a correction becomes a declared mapping without silently altering unrelated entities.

**Reuse by another investigator:** another user with access to the same project files
opens the saved query, supplies their own context and runs it under their own serving
principal. Good result: the query is discoverable and useful without reconstructing it;
project-file sharing is not advertised as hosted team grants or an identity system.

**Epilogue A — close/reopen:** persisted assets and their provenance reload; no claim
is made that unsaved context survives. **Epilogue B — share/replay:** share project/query
identifiers through the established Git/local workflow; the recipient reauthorizes
and reruns. Raw context/results are not embedded in a URL by default.

## Approach

- Prerequisite: Phase 1's corrected transport, source, typed-context and protected-read
  contract plus passing real journeys on recorded revisions. Landed libraries alone
  do not satisfy this prerequisite.
- Reuse the same query model, filestore, semantic engine, policy evaluator and UI
  components. No second collection model, execution endpoint family or permission engine.
- Capture/mapping writes require server-authorized project-write capability, path/ID
  validation, atomic persistence and revision-conflict refusal. A read-only serving
  principal cannot save. Disable submission while pending; retain draft on failure.
- Source/query/mode and applied policy provenance already belong to Phase 1; extend
  with capture/declaration history rather than reintroducing them as missing foundations.
- Collections use existing folders plus README purpose/context first. Additional owner/
  domain schemas require evidence and an owning Feature decision.
- Integration order: capture/provenance/mapping server work; browser adoption; real
  persistence and two-user reuse acceptance. Independent work gets isolated worktrees.

## Tasks

The original task IDs and planning statuses are retained. Verifies names the owning library outcome; Supporting ACs names mandatory cross-capability acceptance, not optional checks. **Activation** marks which
work is in this bounded Phase 2 candidate; deferred tasks require a separate plan before
execution and are not counted in this plan's completion claim.

### Task 1: Semantic types and environment/time kinds

**Id:** task-1
**Verifies:** executable-knowledge-library#ac:parameters-typed-consistently
**Supporting ACs:** semantic-layer#ac:semantic-type-visible-not-coercive
**Depends-On:** —
**Status:** planning
**Activation:** Deferred: semantic-layer extension

Preserve this work as a separately planned semantic-types extension. Phase 1 typed transport/isolation is mandatory already; richer business type inheritance and time semantics must not be confused with that baseline.

### Task 2: Capture as project query

**Id:** task-2
**Verifies:** executable-knowledge-library#ac:capture-lands-as-git-diff, executable-knowledge-library#ac:query-pair-storage-guards-writes
**Depends-On:** —
**Status:** queued
**Note:** Codex task 01a08675-7b77-7d41-9477-ad8168d80a19 reserves the query-capture storage foundation; shared WB claim datatug-phase2-task2. Full capture journey awaits Phase1 acceptance.
**Evidence:** wb-claim:datatug-phase2-task2
**Activation:** Phase 2

Implement authorized atomic query capture from existing executed lookups. Persist <id>.query.json and <id>.query.dtql with purpose and explicit typed parameters. Refuse path traversal, credentials, unauthorized writes and stale revisions; preserve browser draft on failure. No automatic commit or hosted team grant.

### Task 3: Useful folder descriptions

**Id:** task-3
**Verifies:** executable-knowledge-library#ac:folder-recursion-works
**Depends-On:** 2
**Status:** planning
**Activation:** Phase 2

Reuse recursive folders and README descriptions. Make the captured query reachable and explain its purpose. Do not introduce mandatory domain tags, an owner schema or rich collection metadata without a Feature decision.

### Task 4: Capture and result provenance on server

**Id:** task-4
**Verifies:** executable-knowledge-library#ac:captured-query-explains-provenance
**Supporting ACs:** provenance#ac:result-shows-full-chain
**Depends-On:** 2
**Status:** planning
**Activation:** Phase 2

Extend existing Phase 1 source/query/mode/policy metadata with capture lineage where applicable. Preserve structured errors and authorized metadata visibility; never derive security facts by parsing English. A capture links to its source/query definition, not stored protected rows.

### Task 5: Provenance in browser

**Id:** task-5
**Verifies:** executable-knowledge-library#ac:captured-query-explains-provenance
**Supporting ACs:** provenance#ac:result-shows-full-chain
**Depends-On:** 4
**Status:** planning
**Activation:** Phase 2

Show the source, query, mapping origin and capture lineage from actual server responses. Preserve honest live/snapshot and limitation labels. Test old-scope response rejection and empty/denied provenance without intercepted fake data.

### Task 6: Mapping accept/correct persistence

**Id:** task-6
**Verifies:** executable-knowledge-library#ac:captured-query-explains-provenance
**Supporting ACs:** meaning-inference#ac:accept-promotes-to-declared, meaning-inference#ac:correction-overrides-suggestion, provenance#ac:mapping-shows-declarer-and-time
**Depends-On:** 2
**Status:** planning
**Activation:** Phase 2

Use the existing mapping model and explicit project-write boundary. Record declared/inferred/manual provenance, actor identity appropriate to the fixed server principal, and time. Validate physical references and refuse stale revisions; do not create a probabilistic inference engine.

### Task 7: Mapping review UI

**Id:** task-7
**Verifies:** executable-knowledge-library#ac:captured-query-explains-provenance
**Supporting ACs:** meaning-inference#ac:accept-promotes-to-declared
**Depends-On:** 6
**Status:** planning
**Activation:** Phase 2

Expose suggestion reason, accept and correction with explicit confirmation. Wait for persisted success before updating the view; preserve draft on denial/conflict/failure. No claim that an inferred meaning is authoritative before acceptance.

### Task 8: Authenticated HTTP and broader descriptors

**Id:** task-8
**Verifies:** executable-knowledge-library#ac:parameters-typed-consistently
**Supporting ACs:** heterogeneous-sources#ac:http-query-executes
**Depends-On:** —
**Status:** planning
**Activation:** Deferred: heterogeneous-sources extension

Resolve the owning Feature credential-authority and descriptor policy questions before execution. Keyless HTTP GET browser execution and SSRF/snapshot boundaries belong to Phase 1 and cannot be deferred here.

### Task 9: Agent-readable context

**Id:** task-9
**Verifies:** executable-knowledge-library#ac:applicability-updates-with-library-growth
**Supporting ACs:** agent-context#ac:context-verb-matches-web-panel, agent-context#ac:agent-context-blocked-by-policy
**Depends-On:** —
**Status:** planning
**Activation:** Deferred: agent-context parity plan

Expose the same semantic/policy contract later; no agent-only unfiltered view. Scope a separate parity task after the browser journey rather than hiding it inside capture acceptance.

### Task 10: TUI handoff and command hygiene

**Id:** task-10
**Verifies:** executable-knowledge-library#ac:parameters-typed-consistently
**Supporting ACs:** surfaces#ac:handoff-preserves-selection, surfaces#ac:tui-browses-real-schema
**Depends-On:** —
**Status:** planning
**Activation:** Deferred: surfaces parity plan

Preserve as independent maintenance/parity work. Inventory existing commands and respect safe explicit state transfer; no broad CLI deletion in this library-quality candidate.

### Task 11: Time-window auto-narrowing

**Id:** task-11
**Verifies:** executable-knowledge-library#ac:applicability-updates-with-library-growth
**Supporting ACs:** investigation-context#ac:environment-narrows-applicable-list
**Depends-On:** 1
**Status:** planning
**Activation:** Deferred: investigation-context extension

Typed scope and conflict handling are mandatory Phase 1 behavior. Rich time semantics and additional applicability narrowing require explicit types and product evidence; never silently filter an already-open query.

### Task 12: Capture, reopen and second-user reuse acceptance

**Id:** task-12
**Verifies:** executable-knowledge-library#ac:capture-lands-as-git-diff, executable-knowledge-library#ac:captured-query-explains-provenance, executable-knowledge-library#ac:semantic-parameter-declared, executable-knowledge-library#ac:parameters-typed-consistently, executable-knowledge-library#ac:chain-explains-every-entry, executable-knowledge-library#ac:applicability-updates-with-library-growth
**Supporting ACs:** executable-knowledge-library#ac:capture-lands-as-git-diff, provenance#ac:result-shows-full-chain, meaning-inference#ac:accept-promotes-to-declared
**Depends-On:** 2, 3, 5, 7
**Status:** planning
**Activation:** Phase 2

Run the complete browser journey from actual Phase 1 entry with real server/store: capture → reopen → rerun, including denied saves, conflicts, failure draft retention, closing/reopening in the browser and another authorized investigator reusing the saved query. Required setup failures fail CI. Record exact dependency revisions; do not claim hosted sharing or deferred capability acceptance. Completing this task hands the DataTug core to the roadmap's Core coherent-experience review gate, which decides what comes next.

## Deferred AC Coverage

- executable-knowledge-library#ac:playbook-guides-without-automating — Phase 3 guided investigation plan.
- executable-knowledge-library#ac:promotion-changes-visibility-not-content — Phase 4 project grants/promotion plan.
- executable-knowledge-library#ac:public-pack-import-is-explicit — separately planned public packs experiment.
- executable-knowledge-library#ac:useful-query-becomes-check — Track B1, plan 2026-09-11-incidentius-mvp task 5, after this plan's task 2 write path.

Tasks 1, 8, 9, 10 and 11 remain traceable backlog for their owning Features and are
not included in the Phase 2 completion gate. They are deferred for independent scope
and dependency decisions, not to satisfy a task-count cap. No Phase 1 AC is deferred.

## Risks

- Capture may become an unauthorized write endpoint: validate the server capability
  and negative tests before connecting the UI.
- File/Git sharing may be mistaken for hosted team access: use Save as project query
  and explicit recipient authorization, not an unimplemented promotion claim.
- Rich schemas can grow without demonstrated need: start with folders/README and
  existing query metadata, measure second-engineer reuse and authoring overhead.
- A result can leak through saved state: persist definitions, never cached restricted
  rows or credentials; recipient runs under their own policy.

## Open Questions

No new owner decision blocks drafting this bounded plan. Before implementation,
confirm the Phase 1 acceptance prerequisite and activate the relevant draft Features
through SpecScore. Deferred product choices stay in their owning Features.

---
*This document follows the https://specscore.md/plan-specification*
