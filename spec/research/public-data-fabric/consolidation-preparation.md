# Phase1 consolidation preparation

All three bounded discovery streams are independently reviewed and WB-landed. This preparation integrates them on main `e031937c4e52c48e3103da4b346c44ee20e8de06`, retaining original observations and all provisional semantics. Fresh independent combined review passed with zero findings, and root accepts bounded Phase1 discovery at source `6be819afab1483a6c5592f1d7b6ad4d41c95ac6d`. The passed gate/receipt delta still requires final independent review, exact-head CI and WB landing. Dedicated Phase2 reconciliation and Phase3 portfolio architecture have not started. B1 remains open at Phase4/5/6; TPC-H remains excluded.

The only collection extension is optional `research_depth` (shallow/documented/sampled/unknown). Population changes only that field in M1/D1/D2 records. No accepted mapping, canonical metadata, publication or selected portfolio is added. [Query definitions and recipes](index-queries.md) document exact-ID lookup, literal representation policy, graph traversal, depth and freshness provenance. [Dated deterministic results](evidence/consolidation-query-2026-10-04.json) retain the explicit as-of time, reviewed stream heads and input digest; they are derived query output, not a second catalogue.

Executed results:

```json
{
  "collections": {
    "candidates": 681,
    "assertions": 1162,
    "evidence": 128
  },
  "public_datasets": 41,
  "depth": {
    "shallow": 26,
    "documented": 232,
    "sampled": 47,
    "unknown": 376
  },
  "explicit_meaning_associations": 31,
  "canonical_request_accepted": [],
  "four_plus_structural_groups": 6,
  "alpha2_dataset_projection": [
    "d1-dataset-un-m49",
    "d1-dataset-worldbank"
  ],
  "demo_proposal_routes": {
    "adventureworks": 10,
    "chinook": 10,
    "employees": 1,
    "northwind": 7,
    "pubs": 4,
    "sakila": 5
  },
  "freshness_groups": {
    "due": 0,
    "stale_or_needs_review": 254,
    "unknown_or_missing": 80,
    "missing_observation": 47,
    "unscheduled": 307
  },
  "dataset_mapping_groups": {
    "no-mapping-recorded": 23,
    "unresolved-only": 18
  },
  "unlinked_meanings": 232,
  "licensing_decisions": [
    "d1-dataset-comtrade",
    "d1-dataset-un-m49",
    "d1-m49-areas",
    "d1-source-worldbank-microdata",
    "d1-trade-classifications",
    "d1-trade-observations",
    "tpch"
  ]
}
```

A zero accepted Geo.Country result is expected before reconciliation. Explicit provisional associations remain separate. An unlinked meaning means no public-dataset association is recorded; it does not mean no suitable dataset exists. Likewise no mapping recorded, unresolved-only and incompatible/rejected remain distinct. Demo enrichment routes are qualified proposals, not verified cross-demo joins. The alpha2 query only searches literal candidate representations, retains dataset/source provenance and cannot derive dataset use from conceptual mention or differently spelled assertion fields.

Freshness groups overlap and retain provenance for each inherited field. Inspection timestamps are not last upstream release/import timestamps. Missing schedules remain visible rather than becoming invented dates. Depth samples are bounded and candidate-specific; empty inspected GTFS headers do not establish nonnull observations, and documented GHCND daily/crime schemas are not sampled rows.

Durable `consolidation-{m1,d1,d2}-review-r{1,2}` evidence retains earlier findings and exact correction dispositions with source-review hashes. Separate landing records retain reviewed heads, merge revisions, exact-main CI and governed cleanup. Author `consolidation-preparation` evidence proves preparation checks only and cannot substitute for independent acceptance. The Phase1 gate candidate is passed by root decision using the distinct independent discovery receipt; its receiving publication obligations remain unchanged. The dated query result is pinned to reviewed source `6be819afab1483a6c5592f1d7b6ad4d41c95ac6d`. Its input digest excludes all later appended review/decision records and must be reproduced at that source revision, not against this later receipt/gate delta.

Validation passed: inGitDB, existing ledger with all22 mutation checks, SpecScore, whitespace diff and three materially distinct query counterexample fixtures. CI includes those fixtures. Root owns independent combined acceptance, final gate decision and WB landing; this preparation is reviewable research tooling only.
