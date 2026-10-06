---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Public Data Fabric

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/public-data-fabric?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

Explainable public-data enrichment across MeaningGraph, ModelSpec, OpenVaultDB and DataTug, governed by evidence and phase gates.

## Contents

| Child | Description |
|---|---|
| [bigquery-dual-runtime](bigquery-dual-runtime/README.md) | Required Go and browser-direct JavaScript BigQuery access with retained source and runtime gates. |

## Problem

Users connecting operational data should discover useful public joins, lookups and enrichments with an explanation of identifier and semantic compatibility. Similar column names alone are insufficient. Useful unselected or rejected research must survive later implementation waves.

## Behavior

### User journey

1. A user connects a supported source. Observable result: physical fields and declared mappings remain inspectable.
2. DataTug identifies evidence-backed compatible public recordsets. Observable result: each suggestion explains meanings, identifiers, transformations, confidence, provenance and licence limits; ambiguous matches remain visibly unresolved.
3. The user chooses a join/lookup/enrichment. Observable result: a bounded query executes through the existing supported source/federation path and returns usable rows.
4. The user reuses the scenario across differently structured demo databases. Observable result: compatibility follows semantics and documented representations rather than shared names.
5. With no further user action, unselected opportunities remain in the research index. Observable result: later agents can recover dependencies, decisions and evidence.

### REQ: foundation-gate

Broad public-data research MUST wait for Phase 0 foundation review/hardening and a fresh independent verifier. The Chinook golden path MUST work to the supported extent; missing canonical semantic ingestion MUST remain an explicit later integration prerequisite rather than a completion claim. Every blocking finding MUST be resolved or explicitly escalated with rationale/evidence. Gate state lives in the [programme ledger](../../research/public-data-fabric/README.md).

### REQ: preserve-research

Meaning and public-data discovery MUST proceed bidirectionally after the gate. All credible candidate meanings, sources, datasets/recordsets, identifiers, mappings, edges, scenarios/bundles, rejected/blocked decisions and questions MUST have stable IDs, lifecycle state and evidence. Source authority, licence/attribution/restrictions, coverage, freshness and separate practical/showcase value MUST be captured when applicable. Selection is not rejection. Candidate state MUST stay separate from published canonical registries.

### REQ: semantic-reconciliation

A fresh dedicated reconciler MUST review mappings and preserve accepted, rejected and ambiguous decisions. Distinguish semantic equivalence, identifier equivalence, explicit transformations and approximate matching; never infer equivalence from names. A separate fresh portfolio architect MUST select interconnected scenario bundles, identify bridges/hubs and retain unselected opportunities.

### REQ: canonical-waves

Selected integrations MUST publish actual MeaningGraph, ModelSpec, OVDB, registry and Directory artifacts through their existing owners and validators. Preserve exact canonical pins, provenance/licence and query/publication receipts. Implement roughly 5–10 interconnected integrations per wave, then rerank. Current supported browser OVDB federation and CLI anonymous-source parity have different capability boundaries.

### REQ: original-provider-and-retention

The [6 October superseding access decision](../../research/public-data-fabric/p3-portfolio.md#superseding-access-decision--6-october-2026) MUST govern over historical snapshot plans. OVDB definitions MUST identify original HTTP/native-provider resources; DALgo JS direct access is preferred where verified CORS/auth/security/capability/cost permits, with OVDB Go proxy otherwise. Any copy/snapshot of real source or result data MUST have specific user authorization; source terms, query consent and operator configuration MUST NOT substitute. Without it, all owned paths MUST enforce transient bounded processing, reject retained paging before reads/slots/spools, bypass source/result/history/recovery/export/log/cache stores and enforce no-store, with complete sink and teardown tests. Metadata pins MUST distinguish authored definitions from future live bytes and original provider from executor.

Existing GeoNames/ROR copied-provider activation MUST stay frozen pending explicit authorization or a working original-resource alternative. Existing data and ECB bootstrap artifacts MUST remain inactive and preserved until a working replacement and owner-authorized disposition. Public admission MUST wait for exact semantic/rights/read bindings, no-retention/no-store proof and actual local/public journeys; affected paid paths MUST wait for applicability/notices or an enforced free-only boundary. BigQuery A0 remains a priority parallel lane with its separate identity, cost, both-runtime and provider-retention gates. This requirement grants no source-data copy, deletion, paid query or deployment authority.

### REQ: launch-evidence

Each wave MUST be independently reviewed for semantic correctness, model fidelity, provenance/licensing, freshness, queryability, publication and false-positive join risk. Real DataTug journeys MUST prove multi-public-data and cross-demo enrichment despite different representations. Aim for roughly 5–7 excellent scenarios if supported by research, rather than counts as quotas; roughly 40–60 recordsets is guidance and about 75 triggers explicit scope review. The orchestrator MUST record the launch stopping decision and preserve the remaining backlog.

## Dependencies

- semantic-layer
- meaning-inference
- related-data-navigation
- provenance

## Acceptance Criteria

### AC: phase0-controls-discovery (verifies REQ:foundation-gate)

**Given** an unresolved blocking finding still assigned to Phase 0 or a missing independent Phase 0 receipt
**When** an agent tries to mark Phase 0 passed
**Then** ledger validation fails and broad research remains prohibited.

### AC: recoverable-index (verifies REQ:preserve-research)

**Given** candidates at different lifecycle stages
**When** a later agent queries meaning/identifier compatibility, demo coverage, scenario dependencies, rejected licensing decisions or due freshness checks
**Then** stable IDs, links, reasons and evidence are recoverable without consulting private session history.

### AC: separate-specialist-decisions (verifies REQ:semantic-reconciliation)

**Given** proposed mappings and scenario bundles
**When** reconciliation and portfolio passes complete
**Then** their separate reviewer receipts and accepted/rejected/ambiguous decisions remain linked, and unreconciled assertions cannot be promoted as reconciled.

### AC: canonical-publication-and-query (verifies REQ:canonical-waves, REQ:original-provider-and-retention)

**Given** a selected wave
**When** implementation is marked implemented or validated
**Then** exact canonical addresses/revisions, passing publication receipts and actual query receipts exist. For sources without explicit retention authorization, retained paging is refused before source I/O/slots/spooling; synthetic sink/teardown tests show no source/result bytes enter history/recovery/cache/export/log stores on success/empty/error/cancel/navigation/query change; owned routes enforce no-store; actual direct/proxy journeys bind original provider, executor, terms and read observations; existing copied-provider activation remains frozen pending its separate conditions.

### AC: whole-journey-launch (verifies REQ:launch-evidence)

**Given** selected interconnected public data and at least two structurally different demo databases
**When** a user follows the complete connect → explained suggestion → query → cross-demo reuse journey
**Then** joins return correct bounded results, transformations are visible, no unresolved blocking review finding remains, and the backlog survives the recorded launch decision.

## BigQuery priority extension

The [Draft BigQuery dual-runtime feature](bigquery-dual-runtime/README.md) records the parallel user priority and required Go/JavaScript acceptance boundary. Its amended architecture and source/runtime prerequisites remain pending. W1 candidate/scenario value and retained evidence remain, while the 6 October decision supersedes its copied-provider access/refresh plan. New implementation planning follows the same Workbench SpecScore route after A0 and fresh review.

## Open Questions

Current blockers and integration prerequisites are preserved in [fabric.questions](../../research/public-data-fabric/README.md). Phases 0, 0.5, 1 and 2 have bounded accepted gates. The separate [Phase 3 portfolio proposal](../../research/public-data-fabric/p3-portfolio.md) recommends seven conditional scenarios and 37 distinct public recordsets, preserving all research and semantic decisions. Independent r2 review passed and root accepted bounded Phase 3 architecture at `caa9e75d9b24c8ddd95266bc46d97670d492849f`; the [acceptance receipt](../../research/public-data-fabric/evidence/p3-phase-acceptance.json) preserves the exact source and limits. Final receipt-delta review/CI/root landing remain pending. No implementation has started; all scenarios remain conditional. B1 and the selected integration conditions remain open/blocking at Phases 4/5/6. Approved implementation plans follow the existing [workbench plan route](https://github.com/sneat-co/workbench/tree/main/spec/plans/github.com/datatug/datatug), not a parallel hub plan tree.

---
*This document follows the https://specscore.md/feature-specification*
