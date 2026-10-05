# Bounded semantic follow-up to Phase 3 source gaps

Fresh specialist `fresh-semantic-gap-followup`, Codex on local macOS, 5 October 2026 UTC. Input hub `2d355dabf8ef14d80e366f74f8646b9159cd7bf0`; WB task `fabric-phase3-semantic-followup`. This author did not conduct the original discovery or source-gap research. The entire mission, [Phase 2 reconciliation](p2-reconciliation.md), its relevant existing assertions, [source observations](f3-gap-evidence.md), pinned projections, native Sakila rows and private hash-verified Ember annual capture were reviewed. Core geography and binding conventions were read at the Phase 2 core pin `982916d73f0a35ff2558b0062f58aa3ac4f24d97`; existing ModelSpec entity/property versus physical recordset/column distinctions remain authoritative.

This is a dedicated reconciliation decision, **pending separate adversarial review and root landing**. It changes no phase gate, portfolio question, canonical model, registry or runtime. Root alone owns hub landing. Source observations establish what the publishers supplied; the scoped judgements below establish what the reconciler accepts. They are distinct from implementation, product acceptance and permission to publish an integration.

Six new ordinary `fabric.assertions` records carry the accepted scopes. [The derived decision receipt](evidence/f3r-decisions.json) enumerates the reviewed tuples and exclusions; it is evidence in the existing research index, not a second metadata registry. All 1,170 pre-existing assertions remain byte-for-byte unchanged. All existing candidates remain unchanged except `d2-meaning-electricity`: its previous complete record is retained in the receipt, its old reason and evidence stay in the record, and its reusable definition is now narrowed explicitly. The original broad unresolved proposals and rejected World Bank whole-column mapping remain visible; a query must select the new scoped assertion, not reinterpret the old one as accepted.

## Decisions and measured denominators

| New assertion | Scoped decision | Pinned measurement |
| --- | --- | --- |
| `f3r-edge-geo-wb-context` | Accept enumerated GeoNames area → WB reporting-economy context | 252 source rows: 216 matched, 36 unmatched; reverse 217 nonaggregate economies: 216 matched, CHI unmatched |
| `f3r-edge-ember-geo-area` | Accept enumerated annual Ember reporting-area → GeoNames context | 209 projected areas, all matched; reverse GeoNames 252: 209 matched, 43 unmatched; 97,937 fact rows preserve their grain |
| `f3r-edge-sakila-reviewed-country` | Accept exact 108-label partial reference bridge | 109 native rows: 108 matched, Yugoslavia unmatched |
| `f3r-edge-aw-cldr-regular` | Accept 236 reviewed code/name tuples → CLDR48 encoding rows | 238 native rows: 236 accepted, AN/CS excluded as historical |
| `f3r-edge-aw-geo-regular` | Accept those same 236 tuples → GeoNames reference context | Same source denominator and exclusions; this is another target lookup, not another set of source countries |
| `f3r-map-ember-total-generation` | Accept filtered annual total electrical energy produced, TWh | 104,203 annual rows: 5,915 selected values, 98,288 excluded; unique `(ISO 3 code, Year)` after selecting one series |

All listed source-code/label probes have zero NULL, empty and ambiguous rows in these snapshots. All target lookup keys are unique. These zero observations are not permission to discard such rows later: retain separate NULL, empty, unmatched and ambiguous outputs, and fail closed on duplicate targets. The executable counterexamples exercise all five result states and duplicate refusal. Never present an inner-join survivor count as complete source coverage.

## Geography and reporting-economy context

The reconciler reviewed all 216 paired GeoNames/WB labels, not just their matching three-character strings. The accepted target projection is `$[1][*]` from the pinned full World Bank metadata response, with `region.id != "NA"`. It excludes 78 source aggregate rows and retains 217 source economies. Namibia's `iso2Code == "NA"` is a different field and remains included. Region, adminregion, incomeLevel and lendingType identifiers are separate namespaces.

The accepted bridge is the **enumerated tuple list** in `geo_wb.accepted_tuples`: GeoNames ISO, ISO3, name and feature ID, with the WB economy code, label and category. Exact `ISO3 = id` executes that reviewed list; it cannot admit future equal strings automatically. Target uniqueness gives a zero-or-one lookup, with 36 GeoNames misses preserved. `CHI/JG` represents the Channel Islands source economy and cannot be replaced with Jersey or Guernsey. Retained historical `SCG`/`ANT`, territories missing from WB and all other misses remain unlinked. No current official ISO inventory has been certified. In particular, `XK/XKX` remains the publishers' Kosovo representation.

This establishes a **provider-qualified reporting context reference**, not equality of sovereignty, legal jurisdiction, statistical populations or geographic boundaries. Labels such as GeoNames Palestinian Territory and WB West Bank and Gaza must both remain visible; the assertion authorizes identifying related country/economy context, not a claim that their polygons or every statistic's coverage coincide. The existing WB economy → observation source reference can then locate candidate indicator rows. Before comparing or dividing measures, separately establish the indicator's area definition, source, period, unit, revision and rights. Catalogue observation time is not an economic observation year, and undated GeoNames population supplies no annual denominator. No WDI observations were fetched in this pass.

## Ember area, series and quantity are different dimensions

The accepted annual area projection is `DISTINCT (Area, ISO 3 code, Area type)` after `Area type == "Country or economy"` and nonempty `ISO 3 code`. Its 209 code/name tuples were independently reviewed against GeoNames. Fifteen Region areas and their 6,266 source facts are excluded and retained in the evidence. Do not deduplicate fact rows to construct a fact table: they vary by Year and Electricity source. The full country/economy fact key `(ISO 3 code, Year, Electricity source)` is unique in the source observation. Area-reference lookup is many facts to one reference row. The nine Ember areas without a WB economy remain explicit: COK, FLK, GUF, GLP, MTQ, MSR, REU, SPM and TWN.

`is_aggregate_entity` / CSV Area type concern area aggregation. `is_aggregate_series` / CSV Is aggregated source concern series aggregation. They are never interchangeable filters. The 17 observed series include overlapping totals as well as Demand and Net imports; even a nonaggregate series is not necessarily generation. [Ember's yearly documentation](https://ember-energy.org/data/yearly-electricity-data/) corroborates those distinctions. Its download-format paragraph says month despite being on the yearly page; the actual pinned CSV Year column, measured unique keys and yearly endpoint establish the annual grain here.

The narrow production-energy projection selects all of:

```text
Area type = "Country or economy"
ISO 3 code is nonempty
Electricity source = "Total generation"
Is aggregated source = "True"  # validate the observed flag, do not invert it
```

Project area/code, Year, raw `Generation (TWh)` and source revision. There are 5,915 rows with no duplicate `(ISO3, Year)` and no empty generation cell in this pinned capture. The evidence retains the actual year set for every area, not merely a min/max range. The total range is 1985–2025; it says nothing about completeness of monthly data. On refresh, absent rows and missing metric cells remain missing, never zero.

`d2-meaning-electricity` now means electrical energy produced over a specified reporting area and period, qualified by production source and methodology. TWh is energy; Capacity GW, demand, net imports, emissions, intensity and percentage shares are distinct concepts. Default aggregation is none: no sum across fuel totals/constituents, territorial overlaps or periods is authorized here. Source-specific gross/net treatment, estimates, method revisions and historical boundary treatment must remain attached before any cross-provider comparison. No claim is made that the generic core already has this exact quantity meaning; canonical reuse/definition checks precede publication. The accepted binding is specifically the filtered total-generation projection, not every value in the CSV column.

For W2 this supplies an implementable area/year context table. Join an independently reviewed demo-country bridge to the area reference, aggregate demo facts to their intended grain first, and look up one annual context row. Never sum the same contextual national value once per invoice, infer invoice electricity use or carbon footprint, or assign invoice currency from geography. `d2-edge-invoice-electricity` remains unresolved pending its measured invoice denominator and execution policy. `d2-edge-ember-annual-monthly` remains unresolved: no complete monthly universe or full-year comparability check exists.

## Sakila: six explicit representations, one unresolved history

The 102 earlier measured labels were independently reviewed with their GeoNames target labels. They are accepted only as the exact raw values enumerated in `sakila.accepted_labels`; this scope deliberately does not enable generic case folding or broad aliases. The six additions are judgement-based representations supported by the recorded primary-source observations:

| Native ID and exact raw label | Accepted target | Reason and limit |
| --- | --- | --- |
| 25 `Congo, The Democratic Republic of the` | CD / COD / 180 | The Democratic Republic qualifier and concordant code observations distinguish CD from CG. The reordered phrase and final article are explicitly reviewed; an exact GeoNames alternate occurrence is not claimed. |
| 41 `Holy See (Vatican City State)` | VA / VAT / 336 | Exact alternate plus the explicit territorial parenthetical supports VA reference context. The Holy See as a legal/religious authority is not equated with a territory. |
| 51 `Kazakstan` | KZ / KAZ / 398 | Exact alternate and concordant codes support this spelling. No English tagging or general spelling correction is inferred. |
| 80 `Russian Federation` | RU / RUS / 643 | Exact alternate and UN label support reference identity. No historical or disputed-boundary assertion follows. |
| 79 `Réunion` | RE / REU / 638 | Exact accented label supports the separate territory. No global accent stripping or France substitution. |
| 106 `Virgin Islands, U.S.` | VI / VIR / 850 | Explicit U.S. qualifier supports VI; the word-order change is reviewed. British Virgin Islands VG and unqualified Virgin Islands are excluded. No exact alternate occurrence is claimed. |

The lookup retains `country_id` as source-local identity. Multiple alternate-language rows on one GeoNames page are not multiple target identities; targeted-page observation is not proof that a label is globally unique in all gazetteers. The accepted bridge is deliberately source-row scoped. Re-executing it against the pinned native fixture measures 108/109, with the complete unresolved row retained.

`country_id=108`, `Yugoslavia`, remains unresolved. ISO's YU/YUG history and CLDR's RS/ME alias do not determine the intended entity or boundary. `last_update=2006-02-15 04:44:00` is not a validity date. Successor substitution is rejected, including joining this label to Serbia, Montenegro or their sum. Historical enrichment requires a separate entity/interval decision.

## AdventureWorks: encoding lookup with explicit historical exclusions

The full 238-row native code/name projection was reviewed. Each of the 236 regular rows has concordant CLDR48 alpha-2/alpha-3/numeric and GeoNames representations. Accept only those exact native code/name tuples for reference context and encoding lookup. Preserve zeros in numeric strings and all source labels, including older spellings. The CLDR regular category does not itself mean current official ISO assignment; this pass makes no such claim.

AN/Netherlands Antilles and CS/Serbia and Montenegro remain in the source and outside accepted current-reference context. Their occurrence in CLDR encoding rows and retained GeoNames history explains the observed raw 238/238 code overlap without changing accepted coverage to 238/238. The names disambiguate the intended stored labels, but neither `ModifiedDate=2019-04-30` nor a code lookup determines historical transaction validity or a dated statistical boundary. CS reuse and the corrected former-name code are preserved in the source evidence. No AN→CW/SX/BQ or CS→RS/ME substitution is accepted. A historical encoding-only display may show those source records with their deprecated category; it must not present them as linked current territories.

## Reproduction, preservation and publication boundary

Use the exact source captures from the prior bounded research; this command performs no network access:

```sh
wb run -- python3 -B spec/research/public-data-fabric/f3r-reconcile.py \
  --sakila /private/tmp/fabric-phase2-native/sakila.sqlite \
  --ember-annual /private/tmp/fabric-f3-gap/ember-annual \
  --output-dir /private/tmp/data-fabric-semantic-gap-followup/reproduced
```

The generator pins both committed input-evidence hashes and verifies native/CSV hashes before producing output. Changed upstream values fail visibly instead of silently receiving the same semantic decision. Compare its nine generated files to the committed research artifacts. Exact prior electricity meaning state is retained, and the new rationale is appended to the old one. The generator produces records in the existing six-collection format; no schema or validator guard changes are needed.

[Author validation](evidence/f3r-author-validation.json) records the executed checks after refreshing through WB to portfolio base `0eac8db20e78f984367ed767108b2d06ccc2dbf2`: all 24 ledger mutation checks, three query fixtures, two native counterexamples, all six inGitDB collections, SpecScore with zero violations, and whitespace checks pass. Nine generated files reproduce byte-for-byte. The source-input preservation check verified all 1,855 other original assertion/candidate records before the portfolio refresh; root's separately landed portfolio/ROR changes are inherited without alteration. Exact-head CI and adversarial review belong to the separate publication/landing receipts.

GeoNames-derived facts require GeoNames CC-BY-4.0 attribution and change indication; Ember facts require Ember attribution under CC-BY-4.0. CLDR's Unicode License V3 notice remains in the source pack. The World Bank metadata context does not clear every indicator's rights. Native Sakila BSD-3-Clause and AdventureWorks MIT pins stay unchanged. ROR's CC0 metadata versus GeoNames CC-BY location boundary is unaffected. No raw download is committed.

B1 and all Phase 4/5/6 conditions remain open. External DemoDB publication remains with its existing owner and authorization hold. TPC-H stays excluded; `chinook.demodb.dev` remains canonical. The new decisions can address the narrowly matching semantic follow-ups, but root must separately review and update portfolio questions; this lane changes none of them.
