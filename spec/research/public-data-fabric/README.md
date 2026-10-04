# Public Data Fabric programme

Snapshot: 2026-10-04. Owner: main orchestrator. [Owning Feature](../../features/public-data-fabric/README.md).

**Phase 0 is pending. Broad public-data research has not started.** No mappings have been reconciled, no portfolio has been selected, and no public-data launch integration is claimed. The candidate ledger contains only six existing demo demand fixtures plus the preserved TPC-H exclusion. Empty assertion and implementation collections are intentional.

The product journey is connect a supported source → inspect meanings/identifiers → receive explained compatible public suggestions → execute bounded joins/lookups/enrichment → reuse across different physical demo schemas. Each stage must produce observable results and the final acceptance must exercise the real DataTug product path. Component/schema tests alone cannot establish that journey.

## Foundation snapshot and next gate

The three specialist review snapshots found the currently supported semantic and OVDB Chinook fixture sound to the tested extent. Semantic checks pass for Chinook, Northwind, Pubs, Sakila, AdventureWorks and Employees. Local SQLite queries, embedded-grid data and targeted DataTug Go packages pass. The complete MeaningGraph → ModelSpec → OVDB → Directory → DataTug discovery/enrichment journey is not yet implemented or independently verified.

Initial 19:26 UTC snapshot: canonical Chat, canonical data URL, corpus publication and fresh independent verification were pending; cloud PR 13 was observed merged and MeaningGraph registry PR 15 was pending. The original evidence records remain retained.

Updated independent production acceptance, 20:07–20:09 UTC (corpus snapshot 20:02–20:06 UTC): canonical Chat passes at deployed `bddd27a9077b37e3d1e42af8874dd04963afde2b`, seeding all eleven real Chinook tables, attaching scoped context and rendering three ordered Customer rows. Only the AI reply is mocked. Production-origin canonical Invoice HTTPS read passes with 412 rows and SHA256 `88eb7faede360988e9c0f8f8d107e5093db5e712141de14d868f8947b43c373c`. Canonical allowlist/catalog/address contract and negative tests pass; no production catalog UI ingestion path exists yet.

The sole substantive blocking finding is six-corpus publication: AdventureWorks/Employees lack runtime metadata mounts and reciprocal registry/Directory publication; Sakila/Employees website HCL pins are stale. All six hostnames and eighteen model/meaning URLs now return 200. Sakila source/model/meaning registry BSD-3-Clause consistency and Employees full pinned-source rebuild are independently resolved. Declared Node 24 validation is confirmed by exact-main CI and actual runtime logs. Overall independent acceptance remains pending the corpus receipt and root decision; no scope escalation has been approved.

Programme PR 77 landed at `59e8f036ba08103cb33aa7920d8fd98ca9bca179` with exact-main ledger CI success. App PR 189 landed at the deployed commit above. Both have independent source approvals and WB landing/synchronization/cleanup receipts. Compact sanitized details are in `evidence/phase0-hardening-r3.json`. Independent infrastructure review r4 confirms the redundant Cloudflare Workers Builds Git integration is disconnected, preserving the declared GitHub deployment owner and unchanged active version at 100% traffic. Its historical token failure remains Failed; it is not called green.

Anonymous public OVDB access remains surface-specific: browser federation supports optional bearer credentials; CLI parity and chosen browser federation capabilities are later scenario prerequisites. Canonical MeaningGraph/ModelSpec ingestion, explained discovery and integration of the existing HTTPS catalog contract remain Phase 4 prerequisites. Identifier equivalence and ambiguity require dedicated reconciliation. Current independent verdict is one blocker, three Phase 4 majors, zero minors.

Pass Phase 0 only after fresh verification of exact revised/deployed artifacts and resolution or explicit evidence-backed escalation of each blocker. The verifier must distinguish supported foundation functionality from missing future semantic ingestion. Then reread the programme requirements, analyse actual demo representations, and begin bidirectional discovery. Later gates cover dedicated reconciliation, separate portfolio architecture, canonical implementation waves, real scenarios and independent launch review.

## Authority and metadata ownership

This repository owns research and programme state. Implementation plans follow the hub's existing workbench route. MeaningGraph owns meaning definitions/bindings; ModelSpec owns concrete models/fields/keys/references; provider manifests own source revisions, licences, generation provenance and deployment facts; runtime OVDB configuration owns storage/ACL. Git publication registries pin reviewed canonical artifacts, while their generated indexes and websites remain rebuildable.

Research candidates are not published registry entries. Do not add lifecycle/confidence/freshness fields to closed canonical provider schemas just to mirror this ledger. Existing canonical source/notes mechanisms suffice until a selected dynamic integration demonstrates a concrete schema need. A promotion retains research evidence and records the actual canonical address/repository/commit and publication/query receipts.

## Machine-readable collections

The root `.ingitdb/root-collections.yaml` registers six collections. Definitions reuse the ecosystem's `.collection/definition.yaml` + `$records` conventions. Definition files use JSON syntax, which is valid YAML, to make the additional consistency checker dependency-free. Stable IDs are record filenames; references contain keys within their declared target collection, not display titles.

| Collection | Purpose and queries |
|---|---|
| `fabric.candidates` | `kind` selects meaning/source/dataset/recordset/scenario/bundle indexes. Status selects rejected/blocked/low-priority and research depth. Source/dataset/member/dependency/demo IDs connect them. Identifier representations, canonical pins, source/licence/coverage/freshness and separate practical/showcase value are typed optional evidence fields. |
| `fabric.assertions` | Mapping or join edge: recordset/meaning endpoints, actual field names, identifier namespaces, physical role, expected cardinality, compatibility kind, assumptions, confidence, transformation, evidence and rejection reason. `meaning_id` permits joins by meaning; candidate `demo_ids` preserves reuse demand. |
| `fabric.evidence` | Dated observations, decisions, exact commands/revisions, outcomes and reviewer receipts. `review_stage` separates dedicated reconciliation from ordinary observations. `independent` alone is never a passing outcome. |
| `fabric.questions` | Classified blocking/nonblocking/unrelated findings with owners, linked candidates, `gate_ids` membership and evidence-backed closure/escalation. Every blocking question must name its gate(s); each gate's blocker list must exactly match assigned blocking questions. |
| `fabric.implementation` | Candidate promotion/status with exact canonical addresses/pins and distinct publication/query receipts. |
| `fabric.gates` | Requirements, prerequisites, blocker questions, evidence and orchestrator decision for phases 0 through 6. |

Records require evidence for lifecycle decisions. Rejected/blocked records require reasons. Mapping/edge endpoints must have the correct kinds, transformations must be explicit, reconciled assertions require a fresh dedicated reconciliation receipt, implemented integrations require actual canonical/query/publication receipts, and passed gates require closed/escalated assigned blockers, the fixed programme prerequisite graph and passing named independent-review evidence for that gate's stage. All required gate records must remain present. Phase 3 also requires a portfolio architect distinct from Phase 2's reconciliation reviewer. Gate evidence can link new acceptance receipts while question history retains older failing/partial observations. Checks validate declared receipts and relationships; they do not prove the truth of a human review or replace live verification.

Useful lookup procedures: find meaning IDs, filter mapping assertions by `to_id`, then follow recordsets to datasets; filter `identifier_representations` for a representation; inspect candidate `demo_ids` for cross-demo coverage; follow scenario/bundle `member_ids` and `dependency_ids` for prerequisites; filter rejected status and licence reason; compare `revalidate_after` against current time. Zero matches remain a useful gap, not a reason to invent data. inGitDB scalar filters handle simple cases; structured JSON output supports array membership and graph traversal.

```sh
ingitdb list collections --path .
ingitdb select --path . --from fabric.gates --where 'status==pending' --format json
ingitdb select --path . --from fabric.candidates --where 'kind==meaning' --format json
ingitdb select --path . --from fabric.candidates --where 'status==rejected' --format json
ingitdb select --path . --from fabric.assertions --where 'to_id==YOUR-MEANING-ID' --format json
ingitdb select --path . --from fabric.questions --where 'status==open' --format json
```

Run from repository root before committing:

```sh
ingitdb validate --path .
python3 spec/research/public-data-fabric/check.py --self-test
specscore spec lint
git diff --check
```

CI runs inGitDB plus the consistency/mutation checks on every PR and main push. No runtime service, index database or new canonical registry is introduced.

## Evidence retained

- `evidence/snapshot-pins.json`: fetched exact origin/main semantic/provider revisions.
- `evidence/corpus-checks.json`: original executed argv, exit codes and durations for all eighteen provider semantic checks. Temporary paths identify the original review execution environment; re-run equivalent tools against the pinned source. Usage-error retries are not recorded as successful checks.
- `evidence/corpus-models-and-bindings.json`: literal existing six-corpus extraction for later demand analysis, explicitly not new mappings/reconciliation. Legacy URLs remain historical source content.
- `evidence/ovdb-checks.json`: concise exact commands/revisions/outcomes from the OVDB lane. Incomplete Employees rebuild and bounded Directory test run are not full-suite passes.
- `evidence/datatug-findings.json` and `evidence/datatug-checks.json`: source-bearing defect/prerequisite findings and exact targeted commands/outcomes; commands that merely inspect an Nx target do not claim tests ran.
- `evidence/architecture-decisions.json`: canonical ingestion/namespace/execution ownership decisions; this review performed no reconciliation.
- Evidence records and questions retain useful decisions without private prompts, raw browser images, large exports or full temporary logs.

TPC-H remains rejected under the previous 2026-10-04 user-confirmed permission exclusion documented in the existing delivery record. Its newer README-only proposal is not superseding permission evidence. No new licensing research was performed here.

## Open Questions

The authoritative open/closed questions are records in `fabric.questions`; query examples above preserve owners and evidence.
