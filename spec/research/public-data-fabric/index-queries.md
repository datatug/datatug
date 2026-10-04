# Deterministic research index queries

`query_index.py` is a fixed read-only view of the existing candidate, assertion and evidence collections plus the retained actual-six demand coverage matrix. It creates no registry, semantic alias table, service or mutable index. Results are deterministic for the same input files, explicit meaning IDs and explicit UTC as-of time. Run from repository root:

```sh
python3 spec/research/public-data-fabric/query_index.py --as-of 2026-10-04T23:00:00Z --meaning-id Geo.Country --meaning-id d1-meaning-country-area --meaning-id m1-country-area
python3 spec/research/public-data-fabric/test_query_index.py
```

Redirect the first command to a local JSON file when inspecting full results. Exact IDs are mandatory: `Geo.Country` is a literal requested identifier, not an alias for any research or canonical meaning address. A zero accepted result before reconciliation is truthful. Provisional associations appear separately; no M1/D1/D2 equivalence is inferred. Existing inGitDB scalar queries remain available in the programme README; this helper reads the same authoritative record files to handle graph traversal and missing values consistently.

| Mission question | Result key and interpretation |
|---|---|
| Which datasets expose a meaning? | `meaning_datasets` shows exact-ID public-dataset associations and accepted datasets separately. `provisional_meaning_associations` retains every explicitly associated research ID, assertion status, compatibility, fields and evidence. Recordset → dataset is required; demo fixture datasets are excluded from public enrichments. |
| Which use ISO alpha-2? | `alpha2_representations` reports literal alpha2/alpha-2 and ISO2/iso2Code spellings under its documented regex policy. These are proposed representation evidence; neither code equivalence nor field compatibility is inferred. |
| Which meanings occur in 4+ demos? | `four_plus_demo_demand` uses the retained structural count, observed demo IDs and six cell summaries. Structural cells, observed nonnull cells and unobserved cells remain distinct. Zero observed values do not erase structural presence. |
| Which datasets can enrich each demo? | `demo_enrichments` returns qualified proposal routes: direct dataset demand annotations, demo-recordset assertions into public recordsets, and scenario membership. It makes no accepted-join or match-rate claim. |
| Which scenarios depend on a dataset? | `scenario_dependencies` traverses both member_ids and dependency_ids, then dataset/source parents. `scenario_dependents_by_dataset` provides the inverse lookup, including zero dependents. Node lists preserve source/recordset/meaning gaps and cycles are visited once. Conceptual bundles remain distinct from executable scenarios. |
| Which candidates were rejected for licensing? | `licensing_decisions` retains blocked and rejected records whose recorded reason/licence notes mention permission constraints, including TPC-H. This is a recorded-decision query, not an automated legal ruling; needs-review notes remain in source records. |
| Which sources need freshness review? | `freshness.groups` separately returns due, stale-or-needs-review, unknown-or-missing, missing observation and unscheduled IDs. Groups overlap. `freshness.records` distinguishes demo fixtures and record kinds, and retains effective values and provenance for each source/dataset/recordset; recordsets may inherit each field from dataset then source. No review schedule is fabricated. Inspection time is not upstream release/import time. |
| Which promising candidates are shallow? | `research_depth.promising_shallow_ids` uses the optional typed field, independently of status. All shallow and unknown/missing IDs are also exposed. |
| Which meanings have no suitable dataset? | `meaning_dataset_gaps` answers the narrower observable question: no explicitly linked public dataset candidate recorded, or explicit candidates with blocked/rejected and accepted-mapping branches. Suitability remains undetermined until fresh Phase2/3 work. Standards sources are not available datasets. |
| Which datasets lack reconciled mappings? | `dataset_mapping_gaps` distinguishes no mapping recorded, unresolved-only, incompatible/rejected, mixed unreconciled and accepted-reconciliation-backed. Edges are retained separately and do not substitute for semantic mappings. Accepted status requires a passing independent reconciliation receipt. |

## Research depth definitions

`fabric.candidates.research_depth` is optional and research-only:

- `shallow`: credible lead or explicitly shallow concept/source pass, without sufficient source-specific structural investigation.
- `documented`: primary source or retained schema/role evidence inspected and constraints described; source data records need not have been sampled.
- `sampled`: bounded actual source payload records or structural counts/fields inspected for this candidate. This does not claim exhaustive coverage, successful joins, fresh upstream data, rights clearance or accepted semantics.
- `unknown`: insufficient candidate-specific evidence to assign another depth. Missing field is treated as unknown by the helper.

Population is conservative. M1's explicit shallow concepts retain shallow depth; focused primary-standard concepts, standards sources and documented conceptual proposals use documented. D1 provisional meanings use unknown; documentation-only recordsets and dataset/provider summaries use documented, bounded source recordsets use sampled, and retained shallow leads use shallow. D2 normalization follows its source-specific depth evidence after landing. Depth is never derived from promising/researched lifecycle. Dated source observations and original reason prose remain intact. Depth does not automatically inherit from another recordset in the same dataset.

Two counterexample tests cover graph routes, conceptual member dependencies, structural-but-unobserved demand, missing canonical IDs, unresolved mapping acceptance, missing schedules, inheritance provenance, and reconciliation-receipt requirements. Synthetic accepted data appears only in tests and does not mutate the ledger. The existing ledger checker remains unchanged; CI adds the query fixture check.
