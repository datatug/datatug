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

### AC: canonical-publication-and-query (verifies REQ:canonical-waves)

**Given** a selected wave
**When** implementation is marked implemented or validated
**Then** exact canonical addresses/revisions, passing publication receipts and actual query receipts exist.

### AC: whole-journey-launch (verifies REQ:launch-evidence)

**Given** selected interconnected public data and at least two structurally different demo databases
**When** a user follows the complete connect → explained suggestion → query → cross-demo reuse journey
**Then** joins return correct bounded results, transformations are visible, no unresolved blocking review finding remains, and the backlog survives the recorded launch decision.

## Open Questions

Current foundation blockers and later surface-specific integration prerequisites are preserved in [fabric.questions](../../research/public-data-fabric/README.md). Phase0 and Phase0.5 passed for bounded research; broad discovery is underway with individually reviewed streams. Root has accepted bounded Phase1 combined discovery; the gate/receipt delta awaits final review, CI and landing; dedicated reconciliation, portfolio selection and launch implementation have not started. B1 remains open at Phase4/5/6.

---
*This document follows the https://specscore.md/feature-specification*
