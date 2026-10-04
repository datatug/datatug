# Geography, commerce and economic research stream

Observed 2026-10-04. All candidates and joins remain proposed and unresolved. This stream adds 98 candidate records, 36 assertion records, 19 evidence receipts and five open review questions. Source evidence distinguishes inspected records, documentation-only fields and shallow follow-ups. No public data files are published.

The most promising connected scenario is customer or delivery geography → reviewed country/area identity → dated WDI population, GDP per capita and inflation. Chinook, Northwind, Sakila, AdventureWorks and Pubs expose different location structures. GeoNames supplies feature and alias IDs; CLDR supplies territory context and historical currency assignments; WDI supplies dated observations. Map geometry is a separate Natural Earth candidate requiring actual field inspection. Questions include sales concentration, market context and destination differences. Suitable views are ranked tables, time series and carefully scoped maps. Undated GeoNames population must never replace dated WDI observations.

AdventureWorks FX provides a concrete audit scenario: dates 2022–2025 include obsolete DEM, FRF and VEB targets. SIX current/historical ISO 4217 lookup and CLDR assignment intervals can flag inconsistencies. ECB indicative EUR quotes support a proposed same-day cross-rate comparison for supported currencies, with quote direction, business-day absence and rate kind preserved. SIX redistribution remains unverified. Country does not identify transaction currency.

Civil time and calendar context connects GeoNames locations, IANA civil zones and CLDR Windows mappings. A country can have multiple zones; coordinates and zone identifiers require explicit lookup. GOV.UK holidays cover 2019–2028 in the inspected snapshot: useful for relevant Chinook dates, outside Northwind 2016–2018 and Pubs 1992–1994. UK is insufficient to select England/Wales, Scotland or Northern Ireland. Views could show local timestamps and holiday annotations without causal sales claims.

Locale and unit enrichment connects padded AdventureWorks culture/unit codes and Sakila language to CLDR. Legacy zh-cht, the invariant culture and packaging BOX require explicit decisions. UCUM is a distinct namespace with restrictions on modification and derivative standards; OVDB transformation needs legal review. A unit lookup cannot supply product conversion quantities.

US business benchmarks offer a conditional user scenario: a real branch schema with verified country, county GEOID, industry edition and wage units can join Census CBP and BLS CPI/OEWS. Existing demos do not supply sufficient country, NAICS, SOC or monetary-unit evidence. The 2023 CBP metadata uses NAICS2017; a bounded API query returned Missing Key. Pubs and Employees dates require historical releases and explicit user-provided context. Benchmark views include establishments and employment by county/industry, and inflation-adjusted salary trends after units are confirmed.

A new real-organisation scenario uses an affiliation schema with ROR ID or an entity schema with LEI. ROR directly exposes a sampled GeoNames place ID, enabling organisation → place → country → economy joins. GLEIF provides legal entities and accounting-consolidation parent relationships, with exception records indicating unknown rather than no parent. No demo local organisation IDs or names are accepted as global identifiers. Network views and affiliation geography are plausible; ROR and LEI identity equivalence remains unresolved.

| Provider | Lifecycle | Licence captured | Primary source |
|---|---|---|---|
| US Bureau of Labor Statistics | researched | public-domain | https://www.bls.gov/cpi/data.htm |
| US Census Bureau | researched | needs-review | https://www.census.gov/data/developers/data-sets/cbp-zbp/cbp-api.html |
| Unicode Consortium | researched | Unicode-3.0 | https://www.unicode.org/Public/cldr/48/ |
| UN Statistics Division | researched | needs-review | https://comtradeplus.un.org/ |
| European Central Bank | researched | ECB reuse conditions | https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml |
| GeoNames | researched | needs-review | https://download.geonames.org/export/zip/ |
| GeoNames | researched | CC-BY-4.0 | https://download.geonames.org/export/dump/ |
| Global Legal Entity Identifier Foundation | researched | CC0-1.0 | https://www.gleif.org/en/lei-data/gleif-concatenated-file/download-the-concatenated-file |
| UK Government Digital Service | researched | OGL-UK-3.0 | https://www.gov.uk/bank-holidays |
| IANA tzdb project | researched | public-domain | https://www.iana.org/time-zones |
| Natural Earth/NACIS contributors | researched | public-domain | https://www.naturalearthdata.com/downloads/110m-cultural-vectors/110m-admin-0-countries/ |
| OECD | researched | OECD data-use terms | https://data-explorer.oecd.org/ |
| Research Organization Registry | researched | CC0-1.0 | https://api.ror.org/v2/organizations |
| SIX ISO4217 maintenance agency | researched | needs-review | https://www.six-group.com/en/products-services/financial-information/market-reference-data/data-standards.html |
| Regenstrief Institute | researched | UCUM licence1.1 | https://ucum.org/docs/artifacts |
| UN Statistics Division | researched | needs-review | https://unstats.un.org/unsd/methodology/m49/ |
| Wikimedia/Wikidata community | researched | CC0-1.0 | https://www.wikidata.org/wiki/Wikidata:Data_access |
| World Bank Microdata Library restrictive-access class | blocked | Microdata Research License (many datasets; varies) | https://datacatalog.worldbank.org/public-licenses |
| World Bank | researched | CC-BY-4.0 | https://api.worldbank.org/v2/ |

UN M49 and UN Comtrade remain blocked for redistribution pending source-specific permission/access plans. GeoNames postal data, SIX, UCUM and Census remain needs-review. OECD and Wikidata are credible shallow follow-ups, not verified dataflow or entity bridges. World Bank microdata is retained as restricted-access research rather than fabricated survey recordsets.

Five questions preserve licensing, temporal/category equivalence, demo compatibility and operational gaps. Fresh Phase2 semantic review must decide equivalence and cardinality. Phase1 remains pending; corpus publication blocker B1 is owned elsewhere. No TPC-H or competing DemoDB publication work was performed.
