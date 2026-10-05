# W1 hypothetical ROR user-field reconciliation

Dedicated fresh specialist `w1_ror_user_reconciliation`, Codex / Astra6 High on local macOS, 5 October 2026 UTC. Input hub `2694ed58c8bb82277324030cf82a596ee1695db5`. The full mission, approved W1 plan, Phase2 source decisions, existing ledger conventions, exact committed user declaration, provider models/bindings and complete hash-verified SQLite informed this decision. This specialist authored neither the user fixture nor the ROR provider. Root remains sole landing owner.

**Accept the exact hypothetical `affiliations.Affiliation.ror_id` field as an optional native ROR URL reference to the pinned ROR snapshot.** Acceptance is deliberately limited to reference lookup. It proves neither that an affiliation is true nor when it held. This semantic decision awaits separate adversarial review and root landing; production admission still requires canonical companions, the reviewed native-identity attachment, registries, runtime, Directory and app journey. B1 and phases4/5/6 remain open.

[The ordinary assertion](assertions/$records/w1-ror-user-native-url.json) uses the existing source recordset and meaning candidate IDs. [The exact decision and replay evidence](evidence/w1-ror-user-decision.json) bind repository/revision/schema/entity/property, datatype, identifier namespace, role, representation, allowed transformation, status, time and multiplicity. Two new candidate records describe the hypothetical dataset/recordset. All previous candidates, assertions, evidence, questions, gates and implementation states remain byte-for-byte unchanged. This adds no registry, canonical meaning, general transform language or admission callback.

## Accepted tuple and its limits

| Surface | Exact scope |
| --- | --- |
| User source | `datatug/datatug-apps@e7362033ec79c6663d7dbe0483b62fab01f7b9cd` |
| User schema | `libs/datatug/main/src/lib/queries/fixtures/public-data-fabric/affiliations.modelspec.json`; SHA256 `cfe958f12d9a70bb7533499b76bf01c9447f81f0862852e460b095bb75e1b9cf` |
| User field | module `affiliations`, entity `Affiliation`, property `ror_id`, optional string; NULL is retained as a no-reference exception |
| User role | Organization reference on an affiliation row; row identity is `affiliation_id`, not `ror_id` |
| Namespace/transform | `ROR:URL`; native full HTTPS URL; identity only; UTF-8 byte-exact, case-and-space-sensitive equality |
| Target | `ingitdb/ror-ingitdb@bbbec903248680caea04e68f94b9a957b6efc55b`, `model/ror.modelspec.json`, `ror.organizations.id`, required string/native identifier |
| Canonical binding | `model/ror.meaning.yaml`, `research-organization`, identifier binding to native `id`; extends core `organization@982916d73f0a35ff2558b0062f58aa3ac4f24d97` |
| Snapshot | ROR release v2.13, 2026-09-22; source and output hashes recorded in the decision |
| Cardinality | zero-or-one organization reference per input; repeated references remain distinct input affiliations; locations retain their ordered child grain |

[Official ROR identifier documentation](https://ror.readme.io/docs/identifier) prefers the complete URL, while its API also accepts shorter forms. This narrower user contract accepts only the complete native representation. Rejected short, case, space, query and slash variants are unsupported representations here, not assertions that every such form is universally invalid in ROR. The token follows Crockford encoding and the checksum from [ROR's pinned generator](https://github.com/ror-community/ror-api/blob/bd040a0d2558a478c06a89118a29eeb9b6142710/rorapi/management/commands/generaterorid.py). Passing format/checksum does not prove issuance; actual snapshot membership distinguishes matched from unmatched.

[Official ROR status policy](https://ror.readme.io/docs/ror-data-structure#status) distinguishes active, inactive and withdrawn. Every matched record keeps its status. Inactive records supply historical reference context; withdrawn records supply visibly flagged erroneous/duplicate/out-of-scope record context and cannot certify a valid organization. No automatic successor substitution is accepted. The status, name, relationships and locations describe the pinned release, not the user's affiliation time. No user date is supplied and no validity interval is inferred. Changes to source revision/schema or the provider snapshot require reviewed repinning.

This preserves the distinction between native `organizations.id` and runtime serving identity. ROR organization references do not imply LEI/legal-entity equivalence. No person identity, membership, employment, affiliation truth, location-specific affiliation or six-demo native ROR field is accepted. A changed property, namespace or revision is outside this decision even if it has the same spelling or values. A declaration alone and the earlier native public-source decision do not establish that broader mapping.

## Measured outcomes and counterexamples

The committed eight-row fixture yields **3 matched affiliation references to 2 distinct active organizations, 1 NULL, 1 empty and 3 invalid representations**. Unmatched and ambiguous counts are each zero for these eight rows. `a1` and `a3` remain distinct affiliations to `https://ror.org/000025p04`; `a2` references `https://ror.org/00002d369`. Each happens to have one source location; that observation does not prove multi-location behavior.

Separate replay controls therefore exercise one real source record with multiple locations, retained examples of all three statuses, a checksum-valid absent URL, malformed checksum/short/HTTP/trailing-slash/query/non-string representations and a duplicate-target collision. The duplicate produces ambiguity without selecting a target; malformed values never receive normalization. The additional source controls test policy boundaries, not acceptance of additional user fixture rows. NULL, empty, invalid, unmatched and ambiguous must remain distinct, with original values retained.

The complete source contains 141,528 organizations (138,289 active; 1,809 inactive; 1,430 withdrawn). Its SQLite hash is `afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38`; the replay verifies that hash, integrity and unique native identity before measuring. It reads exact Git objects rather than mutable author files. The live GitHub API independently confirmed apps PR192's exact draft head and provider main on this review date. Draft fixture publication is not canonical application release acceptance.

## Reproduction and downstream use

```sh
wb run -- python3 -B spec/research/public-data-fabric/w1-ror-user-reconcile.py \
  --apps-repo /Users/alex/projects/datatug/datatug-apps \
  --provider-repo /Users/alex/projects/ingitdb/ror-ingitdb \
  --sqlite /private/tmp/data-fabric-w1-ror/build-final-a/ror.sqlite \
  --output /private/tmp/w1-ror-user-reproduced.json
```

This performs no network fetch and reproduces the decision JSON byte-for-byte. The replay is research evidence, not a production eligibility engine. Canonical attachment consumers must reference this scoped decision at its eventual landed hub commit and exact file hash, preserve the source tuple, and validate their own execution linkage without widening the decision. Native exact identity needs no manufactured label bridge and no global browser key download. The tooling owner must independently review native-identity support; fitting an older attachment shape is not a reason to change these semantics.

[Author validation](evidence/w1-ror-user-author-validation.json) records actual replay, hash rejection, pinned ModelSpec/MeaningGraph checks, all 24 existing ledger mutations, query fixtures, native probes, inGitDB validation and SpecScore. Independent review and exact-head CI belong to separate root-controlled publication receipts. ROR metadata's CC0 and embedded GeoNames CC-BY-4.0 attribution boundary from the accepted provider remain intact; no source dump is committed or republished here.
