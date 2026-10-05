# [review r1 datatug#91] W1 ROR user-field decision

Reviewed-Head: f8093950122e000aaccffefb7c478adeaf275fbf
Base: 2694ed58c8bb82277324030cf82a596ee1695db5
Reviewer: w1_ror_user_review, independent local Codex/macOS analytical review, 2026-10-05.
Scope: source acceptability and additive research-ledger decision only. No source edits or landing performed; root remains sole landing owner.

## Verdict

No blocking, major or minor findings. The exact hypothetical optional string `affiliations.Affiliation.ror_id` is acceptable as an exact native ROR URL reference into the pinned snapshot under this decision's explicit limits. It conveys neither affiliation truth nor a validity interval. Acceptance does not generalize to other sources, revisions, fields, representations or identifiers.

## Evidence examined

Read the full programme brief and approved W1 plan, Phase2 reconciliation and ledger conventions, all eight added files, exact app fixture schema/data/declaration at `e7362033ec79c6663d7dbe0483b62fab01f7b9cd`, provider model/meaning/source/validation objects at `bbbec903248680caea04e68f94b9a957b6efc55b`, and core organization definition at `982916d73f0a35ff2558b0062f58aa3ac4f24d97`.

Git blob comparisons prove all pre-existing repository files, including all 2,171 prior research files, unchanged. Only eight new files exist in the diff. Existing checkers, workflows, gate state, B1, unresolved/rejected research and implementation state remain untouched. The new ordinary edge resolves to the existing source-recordset and meaning candidates through existing ledger conventions.

The user schema's actual key is `affiliation_id`; `ror_id` is an optional scalar string and is not declared as affiliation identity. The target model's actual key is `organizations.id`, required string; the provider meaning binds that property as the identifier of the source-local research-organization concept extending pinned core organization. Locations have the ordered `(organization_id, ordinal)` grain. No serving ID, organization name, legal entity or LEI is substituted for the ROR URL.

[Official identifier documentation](https://ror.readme.io/docs/identifier) and [pinned ROR generator](https://github.com/ror-community/ror-api/blob/bd040a0d2558a478c06a89118a29eeb9b6142710/rorapi/management/commands/generaterorid.py) support the native full URL and checksum interpretation. The decision correctly narrows representation support beyond forms accepted by the ROR API. [Official status definitions](https://ror.readme.io/docs/ror-data-structure#status) support historical inactive context and erroneous-record withdrawn context, without automatic successor substitution. [The exact release record](https://zenodo.org/records/22902037) corroborates ROR release v2.13, 2026-09-22, 141,528 organizations, CC0 metadata and embedded GeoNames CC-BY-4.0. The release-specific licence statement resolves the generic data-structure page's older CC-BY-3.0 wording; this review changes no inherited source licence decision.

## Reproduced checks

- Author replay reproduced the decision byte-for-byte: SHA256 `a4caa4b6461230a3fa2ae8548121344cb286cfb98a9e19761a207c9f74479595`. Fixture counts are matched3 / NULL1 / empty1 / invalid3 / unmatched0 / ambiguous0, with two distinct organizations among three distinct input affiliations.
- Independently verified SQLite size238563328 and SHA256 `afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38`, integrity, all141528 unique organization keys, all-status counts and141722 locations. Every full-snapshot native key satisfies the checked representation/checksum.
- Independently selected different real active/inactive/withdrawn and multi-location records from the author's controls, cross-checked their results against retained raw JSON, and reproduced4 expanded location rows for2 affiliation inputs referencing1 organization.
- Additional leading/trailing space, hostname/token case, short/HTTP/slash/query, non-string and checksum controls stayed invalid without raw-value repair. A checksum-valid absent key remained unmatched. Duplicate exact target IDs with conflicting statuses produced ambiguity with no chosen target.
- Schema-byte, fixture-byte and SQLite-byte mutations failed with the appropriate hash mismatch before output creation. Wrong source/revision/schema/entity/property/namespace are explicitly outside the accepted immutable tuple; this review does not claim a production eligibility engine was exercised.
- All24 existing ledger mutation checks,3 query tests,2 native-probe tests, inGitDB validation, SpecScore lint and exact-diff whitespace checks passed.
- Extracted exact Git objects into the review's private directory: user/provider ModelSpec publish lint, provider HCL/JSON equivalence, and MeaningGraph binding checks passed. `go version -m` independently confirmed ModelSpec CLI `a8de98b4565ae669b6fba00c334b4649c9d40b11` and MeaningGraph CLI `a58a506d9e6f76660cc8586478521787976111d0`, both `vcs.modified=false`.
- Authoritative GitHub API confirmed PR91's exact head/base, run37253879652 and job111586820013 completed with conclusion `success` at that head, provider main at the exact accepted pin, and apps PR192 still draft at the stated fixture commit. Provider main CI37252466209 is also successful at the accepted source commit.

Private reproductions: `independent.py`, `independent-probes.json`, `pin-controls.py`, `pin-controls.json`, `ci.json`, `pr.json`, `source-ci.json` in this directory. Run scripts through `wb run -- python3 -B` from the author worktree. Author replay output is `/private/tmp/w1-ror-user-independent-replay.json`.

## Limits and downstream requirements

This review authorizes only the narrow decision to proceed to root landing. App PR192 is a pinned hypothetical fixture, not released canonical application authority. No new global ROR bridge is required or justified. A native-identity attachment must reference the landed decision and preserve every source/target pin, property role, raw value, exception and status boundary; its independent consumer tests must reject wrong source/revision/property/namespace/hash and prevent location-based affiliation multiplication.

Canonical registrations, reviewed native-identity attachment, source-faithful runtime, Directory publication, bounded application integration and real browser journey remain separate required gates. B1 and programme phases4/5/6 remain open. No current-affiliation, successor, legal-entity, production-eligibility, published-registry or full-launch claim is accepted. The already independently accepted source build was not downloaded or rebuilt again. Final root merge/main CI/cleanup are not performed by this reviewer.

VERDICT: blockers=0 majors=0 minors=0 land=yes
