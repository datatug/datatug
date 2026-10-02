# Specifications

SpecScore-formatted specifications for DataTug — the product-level hub. Per-surface
implementation specs live in `datatug-cli/spec` and `datatug-apps/spec` and link back
here through `**Implements:**` references.

## Reading order

1. [Product thesis and pillars](product/README.md) — what DataTug is, its fundamental abstraction, the fifteen pillars, non-goals, tested hypotheses.
2. [Feature hierarchy](product/hierarchy.md) — pillars → capabilities with audited state, roadmap phase and kind.
3. [Roadmap](../ROADMAP.md) — phases 0–5, the no-brainers → bets ladder, Dashboardius contract, UI-kit decision, founder questions.
4. [Phase 1 Feature: core investigation loop](features/core-investigation-loop/README.md) and its [plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-09-phase-1-core-investigation-loop/README.md) — implementation-ready.
5. [Product family](product/product-family.md) — DataTug · Incidentius · Dashboardius as one application family: shared capabilities, profiles, ownership boundaries, ecosystem map, decisions.
6. [Incidents Feature](features/incidents/README.md) and its [MVP plan](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/2026-09-11-incidentius-mvp/README.md) — the Incidentius core; siblings [evidence-records](features/evidence-records/README.md), [checks](features/checks/README.md), [annotations](features/annotations/README.md), [compare](features/compare/README.md), [product-profiles](features/product-profiles/README.md).
7. [Research](research/README.md) — current-state and Web UI audits, the incident-resolution vision and its reconciliation (evidence for everything above); competitor research is in the private `datatug/backstage` repository (`competitors/`).

## Contents

| Directory | Purpose |
|---|---|
| [`product/`](product/README.md) | Thesis, pillars, hierarchy |
| [`features/`](features/README.md) | Feature specifications — one Capability per pillar, plus cross-surface Capabilities |
| [`ideas/`](ideas/README.md) | Pre-spec one-pagers exploring problem-direction-MVP |
| [`plans/`](https://github.com/sneat-co/workbench/blob/main/spec/plans/github.com/datatug/datatug/README.md) | Implementation plans linked to Features, stored in `sneat-co/workbench` |
| [`research/`](research/README.md) | Audits and research records |

## Open Questions

None at this time.
