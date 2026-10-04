# D1 dated source evidence

Bounded inspections only. Licence pages were consulted separately from download accessibility. Counts describe the inspected edition or response, not guaranteed current production state. Documentation-only fields remain explicitly unverified.

## BLS CPI/OEWS source rights and dated-series semantics

Record: `d1-bls`. Primary BLS copyright page says BLS-authored publications public domain except previously copyrighted photos/illustrations. Numeric CPI/OEWS candidates support copy/cache/commercial/derivative use under that statement; cite source and revision. CPI flatfile docs define series_id/year/period/value/footnotes and series metadata; seasonality,population,base and annualM13 remain distinct. CUUR0000SA0 is a candidate requiring exact metadata verification. OEWS page currentMay2025 with national/state/metro/industryTXT/XLSX and historical files; SOC/area/NAICS vintage,wage units and suppression matter. No observation sample executed. Employees1985–2002 lacks country/currency, so no automatic US deflation or wage benchmark.

Primary URL: https://www.bls.gov/opub/copyright-information.htm

## Census geography/CBP/NAICS screened with observed key gate

Record: `d1-census`. Primary CBP API/datasets,NAICS homepage,GEOID guide and OpenData policy consulted. 2023 variable table uses NAICS2017,not2022;28 variables includeESTAB,EMP,PAYANN($1000),LFO,EMPSZES,YEAR,GEOCOMP and noise/flags. Geography county key isSTATE+COUNTY,2+3 digit strings; ZIP and ZCTA differ and vintages matter. Annual US/state/county snapshots start1986,ZIP1994. Corrected bounded API query returned HTML Missing Key; access requires review. MetadataJSON hit500k bound and was stopped. OpenData policy did not establish specific redistribution grant; copyright URL inaccessible. Needs-review on rights and credentials, no blanket federal-law inference accepted.

Primary URL: https://www.census.gov/data/developers/data-sets/cbp-zbp/cbp-api.html

## CLDR48 bounded production reference observations

Record: `d1-cldr`. Unicode LicenseV3 consulted: use,copy,modify,distribute and sell with copyright/permission notice in copies or documentation. CLDR supplemental charts and LDML docs consulted. Production latest directory resolves to48 (34MB archive), whereas main/charts49 exists; pin an immutable release before import. Bounded release-48 XML:309 territory mappings,266 currency regions with507 entries,1589 language-population entries,258 territory summaries; units155 conversions/143 preferences; windowsZones500 mapping elements. DE/DEM legal-tender interval ends2002-02-28,FR/FRF2002-02-17; not SIX withdrawal-month semantics. Locale language estimates are not annual official census. Region aggregates remain distinct.

Primary URL: https://www.unicode.org/license.txt

## UN Comtrade trade graph retained with publication restrictions

Record: `d1-comtrade`. Primary use/re-dissemination FAQ consulted. Original redistribution and paid-access conditions depend on subscription/licence/record volume; transformed-output discussion still states active-premium requirements for producers. LicenceAgreement opened as emptyJS; restrictive indexed text retained. No free commercial cache/snapshot permission claimed. Proposed trade graph reporter,partner,period,flow,commodity classification/version,quantity/unit,value/flags; actual response schema/count/coverage uninspected. Statistical territories and HS versions are not demo countries/SKUs/NAICS. Blocked publication,remote access needs review.

Primary URL: https://uncomtrade.org/docs/faqs-on-use-and-re-dissemination/

## D1 actual-six demand constraints

Record: `d1-demo-roles`. Read actual six-demo demand profiles and date ranges. Chinook country labels include USA and Czech Republic, with invoices 2021–2025. Northwind orders span 2016–2018 and supplier Sweden outer space is fictional. Pubs sales span 1992–1994 and stores/authors lack country. Sakila geography uses local foreign keys and film language is English. AdventureWorks has padded unit/culture codes and obsolete currency targets DEM, FRF and VEB in 2022–2025 FX fixtures. Employees salaries span 1985–2002 with 9999 end-date sentinel and no country or currency.

Primary URL: https://github.com/datatug/datatug/tree/90cd16e8ad42131f66145b643697cc97bae8a295/spec/research/public-data-fabric

## ECB FX dated reference observations and reuse conditions

Record: `d1-ecb`. Primary ECB copyright conditions consulted: accurate redistribution cites ECB; label modifications; sold documents notify buyers data is freely available. FX page and bounded dailyXML consulted:1547 bytes,date2026-10-02,USD1.1225 units perEUR. Working-day updates around16CET exceptTARGET closures; historicalCSV/XML/SDMX,1999 onward depending currency. Indicative information only, not transaction quotes. RUB suspended after2022-03-01; AW DEM/FRF and other historical targets lack comparable2022–2025 quotes. Missing dates and rate direction need explicit policy.

Primary URL: https://www.ecb.europa.eu/services/using-our-site/disclaimer/html/index.en.html

## GeoNames postal subset licence and quality limits

Record: `d1-geonames-postal`. Primary postal readme consulted. It declares CC-BY-4.0 while linking CC-BY3.0, and GB_full adds Royal Mail2022 copyright/database rights. Blanket redistribution remains needs-review. Near100 countries; Irish/Chilean/Maltese/Chinese codes can be truncated, Brazil mostly major codes. Coordinates can be calculated averages with accuracy flags. country+postcode can name multiple places. No postal dump downloaded.

Primary URL: https://download.geonames.org/export/zip/readme.txt

## GeoNames gazetteer documentation and bounded samples

Record: `d1-geonames`. Primary extract readme and export terms consulted. Gazetteer TSV is CC-BY-4.0: attribution and change notices allow redistribution, caching, commercial and derivative use. CountryInfo bounded sample: 31,678 bytes,252 rows. Population is undated (US327167434), not annual census evidence. timeZones sample has419 lines including header,418 records. GeoNames IDs connect country/place/admin/aliases/hierarchy; admin1 commonly FIPS-derived with ISO exceptions, never universal ISO3166-2. Daily dumps/deltas; cities5000 documented about50k. Source fields and observations are retained in the linked candidates and dated evidence note.

Primary URL: https://download.geonames.org/export/dump/readme.txt

## GLEIF LEI identity and accounting-parent observations

Record: `d1-gleif`. Primary LEI/LE-RD termsCC0 consulted; supported-access/no-circumvention/no-false-affiliation service conditions remain. Daily concatenated2026-10-04:3,451,839 entity records/518.79MB;672,251 relationships/35.96MB;6,230,385 reporting exceptions/47.98MB. Bounded API page1 size1:2987 bytes,goldenCopy2026-10-04T08:00Z,total3,451,530 differs from concatenated pipeline snapshot. Entity fields and direct/ultimate-parent links observed. Relationships describe accounting consolidation, not complete beneficial ownership; unknown/reporting exceptions distinct from no parent. Demos supply no LEI; real user LEI or authority+registration number required.

Primary URL: https://www.gleif.org/en/meta/lei-data-terms-of-use

## GOVUK regional bank holiday feed permissions and actualcoverage

Record: `d1-govuk`. Primary GOVUK terms and NationalArchivesOGL3 consulted (web open of OGL failed,bounded HTTP succeeded). Commercial copy/cache/adapt/redistribute allowed with source attribution/licence link,no endorsement,and third-party exceptions. Actual JSON22207 bytes:England&Wales83,Scotland94,NorthernIreland103 events; all range2019-01-01–2028-12-26. Fields division,title,date,notes,bunting. Fits Chinook dates but not Northwind2016–2018 or Pubs1992–1994. UK label alone cannot pick division; naive timestamp timezone remains unknown. Rolling feed cadence unspecified,not a permanent archive.

Primary URL: https://www.gov.uk/help/terms-conditions

## IANA tzdb2026e source licence and zone semantics

Record: `d1-iana`. Primary tzdb LICENSE consulted: data public domain; optional date.c,newstrftime.3,strftime.c haveBSD3 code terms separately. Primary IANA page and bounded version file identify2026e,release2026-09-29. zone1970.tab17592 bytes,312 rows: overlappingISO2 country list,principal point inISO6709,ZoneID,comments. Country↔zone is many-to-many; catalogue points are not polygons. Historical rules and links are downloadable; reliable cross-zone equivalence targets since1970,earlier dates limited. GeoNames current-year fixed offsets cannot substitute historical rules or prove naive demo timestamp zones.

Primary URL: https://data.iana.org/time-zones/tzdb/LICENSE

## Natural Earth map geometry: reference rights and category limits

Record: `d1-naturalearth`. Primary terms grant public-domain map use,copy,modify and commercial redistribution; attribution optional. Admin0 countries110m page advertises version5.1.1,210.08KB shapefileZIP and de-facto geometry. Country,sovereign-state,map-unit categories differ; map units needed for overseas regions. Its broad258-country prose is not a measured110m feature count. Exact identifier fields and sentinel values uninspected; no ISO join accepted. Map thematic population lacks comparable dated observation provenance. No archive downloaded.

Primary URL: https://www.naturalearthdata.com/about/terms-of-use/

## OECD economic SDMX screened terms distinct from publication licences

Record: `d1-oecd`. Primary terms section3Data consulted: download,copy,adapt,distribute and commercial use with dataset citation unless specific or third-party restrictions apply. Written-contentCC-BY4 policy is not a blanket data licence. Retained shallow economic SDMX opportunity; no actual dataflow,DSD,codelist,schema,count or coverage inspected. Identify one differentiated flow and its owner/unit/frequency rights before expansion.

Primary URL: https://www.oecd.org/en/about/terms-conditions.html

## ROR research affiliations with explicit GeoNames bridge

Record: `d1-ror`. Primary ROR terms explicitlyCC0 for IDs/metadata; website proseCC-BY4 separately. APIv2 sample02mhbdp94:1500 bytes,fields id,names,domains,types,status,external_ids,locations,relationships. Actual locations.geonames_id3688689 (Bogotá,countryCO) and child RORURL observed. MonthlyJSON dumps and rolling curation documented in FAQ; full dump count/bytes uninspected. ROR organizations need not be legal entities. No measured match to synthetic Pubs publishers; real affiliation/ROR inputs enable the geography bridge.

Primary URL: https://ror.org/about/terms/

## SIX ISO4217 current/historical lists: permission unresolved

Record: `d1-six`. Official ISO4217 maintenance agency lists and primary SIX terms overview consulted; free download is not an identified redistribution/cache/derivative grant. Needs-review,no publication rights claimed. Bounded ListOne47491 bytes,280 country-currency entries,published2026-09-17; ListThree30638 bytes,169 entries,published2026-01-01. CtryNm,CcyNm,Ccy,CcyNbr,CcyMnrUnts; historical WthdrwlDt has month precision. AFA withdrawal2003-01,FRF France2002-03. Preserve numeric leading zeros,funds,code reuse and country duplicates; do not treat a list row as a unique currency.

Primary URL: https://www.six-group.com/en/products-services/financial-information/market-reference-data/data-standards.html

## UCUM licence materially constrains transformed recordsets

Record: `d1-ucum`. Primary UCUM licence1.1June2024 consulted. Unmodified reproduction/distribution and commercial interoperating software are allowed with copyright/licence/disclaimer notices; modifications,derivative works and competing standards prohibited. Developer essence XML and case-sensitive/case-insensitive artifacts documented. Bounded artifact request exceeded500k and stopped,no unit count verified. Needs-review for flattening/reformatting into OVDB; unmodified reference interoperation is a separate potential mode. CLDR conversion data is a separate permissive candidate,not automatic UCUM equivalence. AW codes are local; packaging is not physical-unit conversion.

Primary URL: https://ucum.org/license

## UN M49 statutory area reference with restrictive generic terms

Record: `d1-un-m49`. Primary M49 methodology page consulted: statistical countries/areas,world/region/subregion hierarchy,three-digitM49 andISO2/ISO3 bridge; current online and printed editions1970–1999. Primary UN TermsOfUse search-indexed text consulted,direct open failed: generic personal noncommercial scope excludes redistribution/resale/derivatives absent more-specific permission. M49 snapshot redistribution blocked pending specific grant/legal review. CLDR offers differently licensed containment metadata; its categories remain unreconciled.

Primary URL: https://unstats.un.org/unsd/methodology/m49/

## Wikidata secondary multilingual and identifier bridge

Record: `d1-wikidata`. Primary licensing page consulted: structured main/property/lexeme/entity-schema dataCC0; other proseCC-BY-SA4. Retained shallow secondary QID and multilingual alias/identifier bridge candidate. Suggested P297/P298/P299 country codes,P1566GeoNames claims need exact definitions,samples,rank/qualifier/reference and multiplicity review. No accepted authoritative identity. Bulk snapshots or bounded entityJSON/SPARQL access possible; source operational limits uninspected.

Primary URL: https://www.wikidata.org/wiki/Wikidata:Licensing

## WDI API and indicator-specific permissions

Record: `d1-worldbank`. Primary Data Catalog licence policy and SP.POP.TOTL,NY.GDP.PCAP.CD,FP.CPI.TOTL.ZG pages consulted. Each indicator explicitly CC-BY4.0: attribution/changes and mandatory World Bank dispute terms; recheck any other indicator. API metadata US:id=USA,iso2Code=US; observation country.id=US,countryiso3code=USA. Two population rows2023/2024 have sourceid2,lastupdated2026-07-13 and nullable value/unit/obs_status/decimal fields. Pages advertise1960–2025 but country/indicator completeness varies. Annual observations undergo revisions; JSON/XML API and CSV/XLS snapshots. Region/income/lending and economy aggregates are not countries.

Primary URL: https://datacatalog.worldbank.org/public-licenses

## Bounded response measurements

Measurements are observation receipts, not redistributed datasets.
- `cldr`: {"url": "https://raw.githubusercontent.com/unicode-org/cldr/release-48/common/supplemental/supplementalData.xml", "bytes": 399951, "sha256": "fb054c9f9a95f0c4da9dd1e5e0ef6dd208b6f9a73f45303011be247277d81bef"}
- `cldr-metadata`: {"url": "https://raw.githubusercontent.com/unicode-org/cldr/release-48/common/supplemental/supplementalMetadata.xml", "bytes": 180141, "sha256": "36e807ce72b15304dd993f132216f408351be1c4376edaa2fe9e9547e2efcac1"}
- `geonames-country`: {"url": "https://download.geonames.org/export/dump/countryInfo.txt", "bytes": 31678, "sha256": "93bafc525813f22e4711ff9ed6d626343094ce48c26388dc7c49189b3d7d5512"}
- `geonames-zones`: {"url": "https://download.geonames.org/export/dump/timeZones.txt", "bytes": 13896, "sha256": "ea6f8bdcc259c21c562e8f7e7e0b0457cb89403bed60c76aac49ccee9a9ed18c"}
- `worldbank-country`: {"url": "https://api.worldbank.org/v2/country/US?format=json", "bytes": 421, "sha256": "6d9e77d07055daa0f034dc43ee5a2c1abeaa15d1ee7bacd413a9fd2e1510f51f"}
- `worldbank-population`: {"url": "https://api.worldbank.org/v2/country/US/indicator/SP.POP.TOTL?date=2023:2024&format=json", "bytes": 496, "sha256": "f48a8d48a3053cd51760fac2ee406ab57fe64e9eab011c30d6211ee97772231c"}
- `ecb`: {"url": "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml", "bytes": 1547, "sha256": "d16ef240d4d97a953b925cbf81451c812c559cd33f2ebe07b08b51e4085709e2"}
- `six-current`: {"url": "https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml", "bytes": 47491, "sha256": "33139b438657d1cee116ba737807ea71d19d6de4b90f799a09c56f0cc6a1b0ff"}
- `six-historical`: {"url": "https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-three.xml", "bytes": 30638, "sha256": "98fde2423cdb916dd59dcf5fe96222edad8fa198d865c1c83dbc464b9cc52387"}
- `govuk-holidays`: {"url": "https://www.gov.uk/bank-holidays.json", "bytes": 22207, "sha256": "538b3482c28b85ecd2db606a0d5ae6ad17248900b6498700ce0a48d26a3ecde6"}
- `gleif`: {"url": "https://api.gleif.org/api/v1/lei-records?page[size]=1", "bytes": 2987}
- `ror`: {"url": "https://api.ror.org/v2/organizations/02mhbdp94", "bytes": 1500, "status": "active"}
- `cldr-week`: {"url": "https://raw.githubusercontent.com/unicode-org/cldr/release-48/common/supplemental/windowsZones.xml", "bytes": 49378}
- `cldr-units`: {"url": "https://raw.githubusercontent.com/unicode-org/cldr/release-48/common/supplemental/units.xml", "bytes": 45651}
- `census-cbp`: {"url": "https://api.census.gov/data/2023/cbp?get=NAME,NAICS2022,ESTAB,EMP&for=county:001&in=state:01&NAICS2022=00", "error": "Expecting value: line 1 column 1 (char 0)"}
- `census-naics`: {"url": "https://api.census.gov/data/2023/cbp/variables.json"}
- `iana-zones`: {"url": "https://data.iana.org/time-zones/tzdb/zone1970.tab", "bytes": 17592}
- `iana-news`: {"url": "https://data.iana.org/time-zones/tzdb/version", "bytes": 6, "version": "2026e"}
- `ucum`: {"url": "https://ucum.org/ucum-essence.xml"}
- `ogl`: {"url": "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/", "bytes": 10450}
- `six-page`: {"url": "https://www.six-group.com/en/products-services/financial-information/market-reference-data/data-standards.html", "bytes": 213936}

## Primary rights pages consulted

Public download access was assessed separately from publication permission. Generic policy and unavailable grant cases remain blocked or needs-review as described above.

- [GeoNames](https://www.geonames.org/export/)
- [GeoNames postal scope](https://download.geonames.org/export/zip/readme.txt)
- [Unicode CLDR](https://www.unicode.org/license.txt)
- [World Bank WDI](https://datacatalog.worldbank.org/public-licenses)
- [ECB](https://www.ecb.europa.eu/services/using-our-site/disclaimer/html/index.en.html)
- [SIX](https://www.six-group.com/en/services/legal/terms-and-conditions.html)
- [IANA data and optional code exceptions](https://data.iana.org/time-zones/tzdb/LICENSE)
- [GLEIF](https://www.gleif.org/en/meta/lei-data-terms-of-use)
- [ROR](https://ror.org/about/terms/)
- [Natural Earth](https://www.naturalearthdata.com/about/terms-of-use/)
- [Census policy (no specific redistribution grant established)](https://www.census.gov/topics/research/research-transparency-public-access/open-data.html)
- [BLS](https://www.bls.gov/opub/copyright-information.htm)
- [UCUM](https://ucum.org/license)
- [UN generic terms (M49-specific grant unestablished)](https://www.un.org/en/about-us/terms-of-use)
- [Comtrade dissemination FAQ](https://uncomtrade.org/docs/faqs-on-use-and-re-dissemination/)
- [Comtrade licence (JS page/search indexed terms)](https://comtradeplus.un.org/LicenseAgreement)
- [OECD section 3 data terms](https://www.oecd.org/en/about/terms-conditions.html)
- [Wikidata structured data](https://www.wikidata.org/wiki/Wikidata:Licensing)
- [GOV.UK](https://www.gov.uk/help/terms-conditions)
- [Open Government Licence v3](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/)
