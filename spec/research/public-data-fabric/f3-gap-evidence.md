# Phase 3 targeted geography source evidence

Source researcher `portfolio-gap-evidence`, local Codex on macOS, observed 5 October 2026 UTC. WB task `fabric-phase3-gap-evidence`, hub input `05eecd54e441d7512cd3dc57f8a888595e152984`. This is a bounded follow-up to the accepted Phase 2 gaps, not a reconciliation receipt or portfolio decision. Root owns independent source review and landing; the dedicated reconciler owns equivalence, transformations and exclusions. Existing candidate/assertion/gate files are unchanged.

Machine-readable evidence is [f3-gap-probes.json](evidence/f3-gap-probes.json); its `informs` section enumerates exact existing candidate/assertion IDs for each observation. Three ordinary `fabric.evidence` observation records point to it: `f3-geo-economy-source-probes`, `f3-ember-area-source-probes`, `f3-demo-country-representations`. The file retains complete compact GeoNames/WB area-code projections and WB category universes, Ember annual area/category projection and per-area year coverage, all 238 AdventureWorks code observations, the full small CLDR region-validity universe, original Sakila rows, every unmatched probe row, source URLs/hashes/UTC receipts and rights limits. Raw captures remain private.

## GeoNames and World Bank

GeoNames `countryInfo.txt` is complete: 252 rows, SHA256 `93bafc525813f22e4711ff9ed6d626343094ce48c26388dc7c49189b3d7d5512`, unchanged from Phase 2. `ISO`, `ISO3` and `geonameid` are unique in this snapshot. These include territories, retained historical rows and a temporary Kosovo code, rather than a declared sovereign-state-only set. Original comments identify `CS` as no longer existing, `AN` as dissolved on 2010-10-10, `UK` as reserved while `GB` is official, and temporary `XK` for Kosovo. `ISO-Numeric` remains text; `XK` has `0`, not a normal three-digit ISO numeric representation. FIPS is an alternate namespace, never silently compared to ISO.

The World Bank single-page API response has 295 rows: 78 have `$[1][*].region.id == "NA"` and `region.value == "Aggregates"`; 217 remain under the observed nonaggregate filter. This filter establishes a source category, not a current-ISO country set. The projection preserves economy `id`/`iso2Code` separately from `region`, `adminregion`, `incomeLevel` and `lendingType` category IDs; complete category objects and counts retain exact values and whitespace. Do not confuse Namibia's economy `iso2Code == "NA"` with aggregate `region.id == "NA"`.

| Exact source probe | Total | NULL | Empty | Matched | Unmatched | Ambiguous |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| GeoNames `ISO3` → WB all economy `id` | 252 | 0 | 0 | 216 | 36 | 0 |
| GeoNames `ISO3` → WB nonaggregate `id` | 252 | 0 | 0 | 216 | 36 | 0 |
| WB nonaggregate `id` → GeoNames `ISO3` | 217 | 0 | 0 | 216 | 1 | 0 |

Target duplicate-key groups are zero; max multiplicity is one. Both WB projections have unique keys in the probe. The WB-only nonaggregate is `CHI` / `JG`, Channel Islands; GeoNames separately represents Guernsey and Jersey and has no identical `CHI` key. The 36 GeoNames misses include `TWN`, `REU`, `VAT`, `SCG`, `ANT` and multiple territories. They remain individually enumerated; no territory is made invalid by its absence from WB. WB metadata does not provide historical validity intervals here. Catalogue freshness is not an economic observation year, and this pass fetched no WDI observations.

## Ember annual country/economy projection

The official `files.ember-energy.org` public annual CSV was read once privately with a 17,000,000-byte / 90-second cap, following root's bounded-capture authorization. Actual size 16,084,305 bytes, SHA256 `ea214963f4a98b26f52aaf541736d4310349a905e64f83e63425bfc3ab6255d7`, `Last-Modified: 2026-09-22`, ETag in the machine receipt. No raw CSV is published.

There are 104,203 facts: 97,937 `Country or economy` rows and 6,266 `Region` rows. The derived distinct `(Area, ISO 3 code, Area type)` projection has 224 areas: 209 country/economy areas with unique nonempty ISO3 and 15 region areas with empty ISO3. Region labels are ASEAN, AU, Africa, Asia, EU, Europe, G20, G7, Latin America and Caribbean, MENA, Middle East, North America, OECD, Oceania and World. The annual file spans 1985–2025 overall; every area's min/max year is retained, without claiming complete data for every year or metric.

| Exact source probe | Total | NULL | Empty | Matched | Unmatched | Ambiguous |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| GeoNames ISO3 → Ember country/economy projection | 252 | 0 | 0 | 209 | 43 | 0 |
| Ember country/economy ISO3 → GeoNames | 209 | 0 | 0 | 209 | 0 | 0 |
| Ember country/economy ISO3 → WB nonaggregates | 209 | 0 | 0 | 200 | 9 | 0 |
| Ember country/economy fact ISO3 → GeoNames | 97,937 | 0 | 0 | 97,937 | 0 | 0 |

Every target projection has zero duplicate code groups and max multiplicity one. Fact `ISO 3 code` repeats across years and electricity sources; per-code row counts are retained. The nine WB misses are `COK`, `FLK`, `GUF`, `GLP`, `MTQ`, `MSR`, `REU`, `SPM`, `TWN`. All are preserved with exact Ember area labels.

CSV country/economy filter: `Area type == "Country or economy"` plus nonempty `ISO 3 code`. OpenAPI 1.2.0 yearly/monthly parameter `is_aggregate_entity=false` concerns area aggregation. `is_aggregate_series` and CSV `Is aggregated source` concern fuel/series aggregation and cannot substitute for the country filter. The 17 actual electricity-source/aggregation pairs are retained. `(Area, Year, Electricity source)` and the filtered country `(ISO3, Year, Electricity source)` have zero duplicate tuples; this does not authorize summing Total generation/Clean/Fossil alongside constituent fuels. Missing metric cells, method revisions and same-period alignment remain future integration requirements.

Unauthenticated options and one-year data probes returned 403. No API key/account was created. No complete monthly universe was fetched; annual coverage and code matches cannot establish monthly coverage. The earlier monthly prefix remains the accepted limited receipt.

## Sakila representation follow-up

Native source SHA256 matches Phase 2; fields are actual `country.country_id`, `country.country`, `country.last_update`. All seven gap rows have `last_update = "2006-02-15 04:44:00"`; this operational update timestamp does not establish the intended country's historical boundary or validity date.

| Native ID and raw label | Target candidate only | Exact raw alternate-name evidence |
| --- | --- | --- |
| 25 `Congo, The Democratic Republic of the` | CD / COD / 180 | No exact occurrence in targeted GeoNames page. English alternate is `Congo, The Democratic Republic Of` (missing final `the`); UN has `Democratic Republic of the Congo`. Keep explicit representation-review gap. |
| 41 `Holy See (Vatican City State)` | VA / VAT / 336 | Exact English GeoNames alternate name; UN names `Holy See`, CLDR48 labels `Vatican City`. Keep territory/authority scope review. |
| 51 `Kazakstan` | KZ / KAZ / 398 | Exact untagged and multiple language-tagged GeoNames alternate rows; no English tagging is inferred. UN/CLDR use `Kazakhstan`. |
| 80 `Russian Federation` | RU / RUS / 643 | Exact English GeoNames alternate and UN name; CLDR48 labels `Russia`. |
| 79 `Réunion` | RE / REU / 638 | Exact English and multiple language-tagged GeoNames alternates; exact UN/CLDR label. No global accent folding inferred. |
| 106 `Virgin Islands, U.S.` | VI / VIR / 850 | No exact occurrence in targeted GeoNames page. GeoNames English labels `U.S. Virgin Islands` and `Virgin Islands (US)`; UN has `United States Virgin Islands`. Do not confuse British Virgin Islands VG/VGB. |

The six candidate codes come from source observations, not approved aliases. Exact occurrences were searched in six specifically targeted feature pages and are not a global alternate-name uniqueness test. Multiple language rows for one feature do not become multiple target identities. Original Phase 2 102/109 bridge coverage remains unchanged; this pack publishes no revised accepted 108/109 result.

`country_id = 108`, raw `Yugoslavia`, stays separate. Official ISO 3166-3 Newsletter I-3 (2003-07-23) records the coding/name event YU/YUG/891, validity 1974/2003, former-name code YUCS. It does not prove which historical Yugoslav boundary Sakila intended. CLDR48's deprecated YU alias lists RS/ME, which is retained as a lossy source alias observation only. No modern successor is substituted and no date/boundary issue is closed.

## AdventureWorks code categories and history

The pinned 238-row `Person.CountryRegion` fixture reproduces exactly; `CountryRegionCode` is its source-native primary key. All 238 codes occur in GeoNames and CLDR48 `territoryCodes`, with no NULL/empty codes, no unmatched codes and no duplicate target key groups. Full row observations retain `Name` and `ModifiedDate`, uniformly `2019-04-30 00:00:00.000`. That timestamp cannot make obsolete reference codes current.

CLDR48 `common/validity/region.xml` classifies 236 source codes `regular`, two `deprecated` (`AN`, `CS`), and zero `special`, `macroregion`, `reserved`, `private_use`, `unknown` or missing. The full small category universe and exact structural path are retained. `regular` is CLDR's category; it also includes extensions such as XK, AC, DG, EA, IC, TA and therefore does not prove current official ISO assignment or sovereignty.

`AN` and `CS` are absent from CLDR `territoryInfo` but present in its encoding crosswalk. CLDR deprecated alias AN has **three** replacements CW/SX/BQ; CS has **two** RS/ME. Official ISO notices establish AN/ANT/530 deletion (Newsletter VI-8, 2010-12-15) and former CS/SCG/891 validity 2003/2006 with ISO 3166-3 code CSXX (Newsletter I-5, 2006-12-01). That latter notice corrects CSHH because it previously referred to Czechoslovakia, demonstrating code reuse/history ambiguity. Successor lists are not safe single-row identity repairs.

CLDR release-48 resolves to commit `acd6d88ae493633240e19a87a721076a8a75c310`; supplemental-data content matches the exact Phase 2 hash. Direct ISO HTTP fetches returned 403. Official PDF text was consulted through the web provider; the pack preserves document dates/pages/URLs and explicitly unavailable byte hashes. No complete current official ISO/reserved-code inventory is claimed, and no restricted ISO PDF is redistributed.

## Rights and reproducibility

GeoNames data is CC-BY-4.0: credit GeoNames and indicate changes. Ember's `info.license` / `info.description` explicitly licenses data CC-BY-4.0 for any purpose with attribution and no additional restrictions. World Bank's public-license policy defaults to CC-BY-4.0 with extra mediation/arbitration terms unless a dataset states otherwise; this small metadata probe is not licence clearance for every indicator. CLDR's Unicode License V3 copyright/permission notice is retained in [f3-unicode-LICENSE.txt](evidence/f3-unicode-LICENSE.txt). Source-native Sakila BSD-3-Clause and AdventureWorks MIT scopes remain attached to their pinned evidence. UN contributes six bounded factual code observations only; no blanket table-republication rights are inferred. The pack publishes derived facts/projections, source attribution and links rather than raw source documents or exports.

Reproduce with the durable source probe, using the existing private captures and hash-verified native fixtures:

```sh
python3 spec/research/public-data-fabric/f3-gap-probe.py \
  --captures /private/tmp/fabric-f3-gap \
  --sakila /private/tmp/fabric-phase2-native/sakila.sqlite \
  --adventureworks /private/tmp/fabric-phase2-native/adventureworks.sqlite \
  --output /private/tmp/fabric-f3-gap/reproduced.json
```

The output equals `f3-gap-probes.json.observations` exactly. `--capture` refreshes only the recorded official URLs privately with bounds; all hashed bodies must match the pinned receipts before reproduction. Changed upstream revisions fail visibly and require a new dated observation, rather than silently changing this snapshot. ISO/web-provider derived history facts require manual source inspection because the direct 403 prevents byte-level reproduction. UN's bounded prefix is incomplete as a document but contains the complete English table; no complete all-language document is claimed.

The researcher ran inGitDB, ledger checks, SpecScore and diff checks through WB as recorded in the PR. Source observations do not close B1, select a portfolio, authorize canonical publication or establish a DataTug journey. TPC-H remains excluded; the canonical Chinook hostname stays `chinook.demodb.dev`.

## Additional bounded ROR metadata corroboration

At the portfolio architect's request, an independent read of [ROR release metadata](https://zenodo.org/records/22902037) confirms ROR data release v2.13, 2026-09-22, publisher-reported 141,528 organisations, schema version 2, archive `v2.13-2026-09-22-ror-data.zip`, displayed 37.5 MB and advertised MD5 `a935217de21e7983b977f5419f1a6b20`. The archive was not downloaded or checksum-verified. The primary description explicitly separates ROR metadata CC0 from GeoNames location data CC-BY-4.0; JSON is the record format, while CSV is a subset with flattened fields. [Bounded receipt](evidence/f3-ror-release-observation.json) and `f3-ror-release-rights` preserve the independently hashed metadata page and existing IDs informed. This does not certify a global import, whole-dump CC0 rights, semantic compatibility or portfolio scope.
