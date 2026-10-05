# GeoNames W1 independent adversarial review r2

Reviewed-Head: bb603ab40fb8d1478e6ab496e84f6bf02c1e8c3c
Base: 514b55042645c2bb63aa5a7f85c6171580f3926c
PR: https://github.com/ingitdb/geo-ingitdb/pull/6
Runtime: local Codex/macOS. Root is sole landing owner.

**Source-only landing accepted. Both r1 majors are fixed; no unresolved findings.**

## Resolved controls

**GEO-R1-1 — resolved.** Independent probes replay the original duplicate `locations` and nested `geonames_id` cases, NaN, numeric overflow, and non-JSON whitespace. All reject at read chunk sizes 1, 7, 65,535, 65,536 and 65,537. The parser now rejects duplicate object members recursively and nonfinite values, recognizes JSON whitespace only, and applies the UTF-8 record-byte limit even after a successful decode. Independently constructed ASCII and multibyte records at exactly 1MiB pass; +1 byte fails at the three large chunk boundaries, with a following valid object retained. Original duplicate-JSON full build writes a blocked receipt and leaves no output/stage. Complete pinned ROR archive revalidation passes under the repaired reader: 141,528 organizations, 141,722 locations, 19,729 IDs, original archive SHA256 `476988542f1bdd185f779386f90149371f58cb4ebdba326db3b7b428a468d4a5`.

**GEO-R1-2 — resolved.** Original non-ROR city ID99 with disagreeing allCountries fields now fails, writes blocked receipt and leaves no output/stage. Matching overlap retains exactly one row. The comparison now covers all initially retained city IDs independently of ROR membership; global duplicate-key indexing still applies. Documentation explicitly describes separately fetched upstream files and failure on conflicts, without an atomic-source assumption. R1's complete pinned-source audit found zero actual retained-row conflicts.

## Reverified evidence

- Authenticated GitHub PR and check-run APIs independently confirm this exact remote head and unchanged base. Validator, Build & test, Lint, offline_projection and wb-import100% checks succeeded; three expected publication/version/reuse checks skipped. All checks terminal at this head.
- Offline suite 15/15 passes. Pinned ModelSpec publish lint/export and MeaningGraph check pass with zero warnings/errors.
- Delta comprises importer, tests, importer documentation and durable validation receipt only. Accepted bridge dictionaries, model/bindings, native namespace, canonical pins and all legacy importer/collection behavior remain unchanged.
- Repaired `r1-a` and `r1-b` are byte-identical across all17files. All16files other than snapshot equal the previously independently audited `corrected-a`. Only snapshot generator/importer revisions and importer SHA changed; importer SHA matches the reviewed script. New snapshot SHA256: `d564721b80809537424d17b2bf2276697593a3cb380ea020f33c070e1ff0e7ab`.
- Ordered chunks reconstruct 75,304,960-byte SQLite SHA256 `cb53af14c67daadfc20ecbfe304d631d651022179793a77aa11413244d2f456d`; SQLite integrity and every declared physical table count pass. Metadata is42,238bytes. Counts remain252countries,3,865admin1,75,300places,844,831aliases,19,729requiredROR IDs,9missingROR IDs and234missingadmin1 references. Four bridges remain24/21/21/3 and country keys retain `GeoNames:countryInfoISO2`. Source-field fidelity, exceptions, exact P2 dictionaries/denominators and model semantics inherit the complete r1 row/provenance audit via byte identity.

## Measurements and boundaries

Inspected external author receipts report169.33/169.95seconds and129,482,752/129,761,280-byte peakRSS. Maximum sampled logical input/staging disk is3,774,135,231bytes; captured source bytes remain670,249,165. These fit the source-generation budgets; output SQLite/metadata fit256MiB/2MiB. Measurements are inspected author observations, not independently recreated peaks. Initial capture whole elapsed/RSS was not instrumented; combined-wave resource acceptance remains open.

The reviewer re-ran complete ROR validation, not another full GeoNames build or another global uniqueness index. Full r1 source-row evidence is preserved and transferred only through verified byte identity. No source edits, merges, browser/port launches or publication occurred. Production representation contracts, registry/Directory publication, runtime native-id handling, hosting, ROR user-contract acceptance, B1 and real-browser whole-W1 acceptance remain separate downstream gates. This verdict authorizes root to consider source-only landing, not completion of W1.

Private reproductions: `r2-probes.py` / `r2-probes.json`, `r2-bundle.py` / `r2-bundle-audit.json`, alongside frozen `review-r1.*` and full raw-source evidence. Run the scripts through `wb run -- python3` from the author worktree.

VERDICT: blockers=0 majors=0 minors=0 land=yes
