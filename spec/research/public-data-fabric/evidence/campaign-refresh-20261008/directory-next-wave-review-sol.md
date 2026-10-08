# [review r1 directory-next-wave] Independent review

Verdict: **GO for bounded next-wave preparation**, with the two wording corrections below. This does not approve source execution, direct-browser readiness, deployment, semantic extensions or activation. No blocking or major finding in that limited scope.

I independently queried each repository's remote `refs/heads/main` with `git ls-remote` and read the matching `origin/main` objects. All three SHA receipts match the memo: DataTug `3973ae1a5bcd1c2959fef3bea7f8aafa54a156e5`, Directory `480f369bedf4c3c8c6d1d67c32a083c6c820de54`, website `3a789ea00d65e329c0ee64f51474529ac6c6e185`. Website main ends at PR34. No provider API/resource call, source-row copy or repository edit was performed.

## Findings

1. **Minor — IANA deployment overclaim** (`directory-next-wave-audit-sol.md:17`). Replace “deployed-code preparation” with “implemented transport preparation.” The authoritative `openvaultdb/directory` source `sources/$records/iana-http-status-codes.yaml:34–51` says delivered implementation and explicitly says synthetic checks do not establish live execution, production configuration, deployed sink behavior or a browser journey. Its linked implementation commits and default-off composition are useful evidence, but no deployment receipt follows from them.

2. **Minor — GLEIF status scope is broader than the current native proposal** (`directory-next-wave-audit-sol.md:14`). Replace “initially identifier/status scoped” with “initially identifier scoped; preserve and separately admit registration/entity status context.” `sources/$records/gleif-lei-entities.yaml` proposes only `lei` under `data[].attributes`. `d1-mapping-lei-id.json` requires preservation of registration/entity status, but reconciles the identifier only. The exact-key route, single-resource envelope and status fields remain unaccepted. The memo already lists those admission gaps, so this is wording rather than an execution authorization defect.

## Checks and limits

- GOV.UK is a defensible first preparation target: small historical response and reconciled division/date/title event meaning. Its enclosing division, notes, multiplicity and date policy are correctly preserved. `p3-portfolio.md:140–161` additionally requires explicit jurisdiction/zone for any later transaction enrichment; Northwind/Pubs are coverage negatives and Chinook needs supplied, reviewed division context. The memo promises only a Directory source journey, so it does not claim these joins.
- GLEIF and CISA provide plausible real identifier utility while the memo correctly disclaims demo keys, name/product matching, ROR/LEI identity and beneficial ownership. CISA CVE mapping and the exact KEV/NVD namespace edge are reconciled; NVD rights and full vulnerability journey remain blocked. CLDR and IANA remain raw proposals, not inferred currency validity, rounding or scalar status mappings.
- All shortlist HTTP definitions are inactive with proposed access routes and unverified GET/CORS. Documentary licences do not establish exact returned-material/intended-use admission. Historical observations and CLDR metadata pins are accurately distinguished from current response evidence.
- WDI has one exact `country_summary` schema/time/US observation; query activation is blocked. Website code defaults `OVDB_WDI_METADATA_ENABLED` to false and separates sign-in from GIS/metadata access. The memo's instruction to leave metadata and queries disabled in this wave is sound; it makes no live deployment claim about configuration.
- No-copy policy, source authority, frozen GeoNames/ROR activation, original candidates/dispositions, pending phases 4/5/6 and five implementing W1 records are preserved. Original-provider identity remains separate from any future gateway identity.

For an implementation brief, carry forward the complete admission checklist from the authoritative portfolio: finite row/byte/time budgets before reads, a complete sink inventory and owned no-store enforcement including browser/service worker, cancellation/navigation/query-change teardown, exact source/model/decoder/profile pins, local/public browser receipts, and paid-use applicability/notices or an enforced free-only boundary. Synthetic tests prepare those gates; they cannot satisfy live source acceptance. The memo is suitable for preparation prioritization, not a substitute for that source-specific brief or gate closure.

VERDICT: blockers=0 majors=0 minors=2 land=yes
