# Phase 3 scenario and dataset portfolio architecture

Author proposal, 5 October 2026. Fresh separate specialist `fresh-phase3-portfolio-architect`, local Codex on macOS. Read the complete mission and full persistent inventory at hub `05eecd54e441d7512cd3dc57f8a888595e152984`; semantic decisions are those of the separate Phase 2 source `c26c3ee5d8e5e0edce629b66fd15d6d3cf323192`. No new semantic decisions, canonical artifacts, deployments or phase pass are made here. Independent adversarial review and root acceptance remain pending.

## User journey and observable results

1. I connect my supported source. I can inspect physical schema, declared meanings and pins. If I do nothing else, no user rows are uploaded or queried for metadata discovery.
2. I browse compatible public opportunities. Each explains a specific accepted field/identifier path, source licence, snapshot, coverage and remaining ambiguity. Unsupported paths are shown as unavailable with a reason, not offered as working joins.
3. I choose one bounded enrichment and review its query and projection. I see raw values, join direction, cardinality, time/unit policy and source provenance before execution.
4. I run it. The existing browser OVDB federation/DALgo path returns rows plus explicit matched, unmatched, NULL, ambiguous and multiplied counts; failures and stale capabilities remain visible.
5. I repeat a supported journey on another actual demo representation. The same canonical knowledge is reused, while source-specific bridges and exclusions remain explicit.
6. With no further action, unselected research and declined/gapped paths remain queryable in the same collections. A canonical registry mutation changes later suggestions after documented cache refresh without a hardcoded public catalogue edit.

## Recommendation and counts

Recommend seven **conditional target scenarios**, sequenced by shared foundations and bounded readiness. There are **37 unique public recordset candidates, zero implemented integrations and zero launch-ready scenarios** in this proposal. This is deliberately below the 40–60 planning guidance. Do not inflate counts with individual indicators, aliases, native demo FKs, binding records, representations, empty feed files or repeated scenario membership. All 57 original scenarios/bundles have a disposition; all 682 original candidate IDs/history and 1,170 assertions are preserved. The rejected TPC-H decision remains intact.

[Machine-readable portfolio, source modes, exact assertion graphs, 57-row lineage and 682-row disposition receipt](evidence/p3-portfolio.json) are derived evidence linked by `p3-portfolio-architecture`. Seven new `p3-scenario-*` candidate records in `fabric.candidates` use existing member/dependency/evidence/value fields and stay `needs-review`. No collection/schema extension or implementation record was needed. The existing literal `query_index.py` still cannot infer semantic aliases, parse every namespace spelling or prove suitability from zero matches.

| Target | Practical /100 | Showcase /5 | Current readiness and reason |
|---|---:|---:|---|
| Market geography and dated economic/energy context | 87.4 | 4 | Conditional; Highest cross-demo demand and reusable country/place foundation. Electricity is consolidated into market context because it reuses the same country/year bridge; monthly electricity remains backlog rather than a second launch scenario. |
| Comparable international catalogue, quantities and prices | 81.2 | 3 | Conditional; Consolidates overlapping money, unit and locale proposals into one international catalogue journey. It avoids redundant SIX/UCUM imports; currency, unit and language remain distinct meanings and query steps. |
| Local operating calendars with explicit place and time | 75.4 | 4 | Conditional; Calendar/time is a cross-scenario hub for market years, FX dates and environmental observations. High demand but low current accepted demo coverage keeps it behind initial geography/catalogue bridges. |
| Real organization and affiliation exposure | 75.8 | 4 | Conditional; Adds strong portable real identifiers beyond DemoDB while reusing geography. Consolidates affiliation and supplier proposals without pretending ROR and LEI are equivalent. DOI/ORCID output discovery stays backlog. |
| Premises and transport access context | 51.2 | 5 | Conditional; Compelling connected civic/mobility story but highest bridge/right risk. Rank last for implementation; keep it as a conditional seventh target rather than force a feed-only demo into launch. |
| Asset observations and environmental context | 72.6 | 5 | Conditional; Reuses unit/time foundations and connects reference metadata to observations. Strong showcase after real typed query value; earthquakes remain optional reserve rather than displacing daily operational utility. |
| Operational logs and vulnerability identifiers | 73.4 | 3 | Conditional; Counterbalances demo bias with common real operational schemas and high identifier quality. Two related user journeys share service scope, but are not counted as a fabricated public-data connected component. SPDX rights remain blocked and excluded. |

## Scoring and portfolio tradeoffs

Each dimension is scored 0–5; weighted practical score is the sum of weight × score / 5. Higher cost/risk score means cheaper/safer delivery. Source and legal gates are vetoes, not quantities that a high aggregate can outweigh. Scores compare opportunities after missing work is explicit; they are judgment estimates, not measured market-share forecasts. The JSON supplies every raw score and dimension rationale. Showcase is separately scored and does not enter the practical score.

| Dimension | Weight |
|---|---:|
| real world joinability | 20 |
| identifier quality | 15 |
| semantic connectivity | 10 |
| cross scenario reuse | 10 |
| legal confidence | 10 |
| source authority | 5 |
| schema prevalence | 6 |
| actual demo coverage | 7 |
| representation diversity | 4 |
| quality cadence operations | 8 |
| implementation cost risk | 5 |

Geography wins for real coverage across three accepted demo label paths and eventual structurally different country tables/code chains. Catalogue follows for actual padded unit/language bridges even though public price/quantity comparisons still need work. An initial ROR reference is pulled into W1 to complete an accepted public graph, while LEI and economic extensions stay conditional; calendar follows catalogue because its time policy is shared by later observations; operational technology moves before environment/mobility because identifiers are precise and source rights can be bounded. Mobility stays last despite showcase value 5: the missing station crosswalk, FSA identity/spatial contract and ORR rights are consequential. These tradeoffs are intentional rather than a sort by score.

Use GeoNames for the first country/place graph; CLDR supplements code/context/unit/locale roles, not a competing country authority. World Bank supplies dated observations; undated GeoNames population is not an alternative. Ember annual adds a different question after category/period acceptance; monthly is deferred. IANA tzdb supplies civil rules; current GeoNames offset tables are redundant for that purpose. SIX, UCUM, UN M49, OECD and Wikidata stay useful alternatives with distinct rights/coverage or evidential costs. ROR and LEI remain separate identity branches. NTA native GTFS connectivity does not overcome unresolved rights. No weak seventh scenario must ship merely to satisfy a count.

## Scenario contracts

### p3-scenario-markets: Market geography and dated economic/energy context

I connect customers, order destinations or publisher countries; inspect the country bridge; choose a year and indicator; request a bounded left enrichment; compare business concentration with dated population/economic/electricity context.

**Useful questions.** Which customer or shipping countries remain unlinked, and why? How does order/customer concentration differ from population or GDP per capita in the same year? How does country electricity mix vary across our market footprint, without claiming our purchases caused or emitted it?

**Actual schema and representation boundary.** Three demos have accepted public label paths: Chinook Customer.Country 59/59; Northwind Customers.Country 91/93 total, 91/91 nonnull, two NULL; Orders.ShipCountry 830/830; Pubs publishers.country 8/8. Sakila address→city→country uses integer local FKs but public bridge remains unresolved (102/109 observed; seven unresolved). AdventureWorks Address→StateProvince→CountryRegion has 238 code-shaped values; category bridge unresolved. ISO alpha-2/alpha-3/zero-padded numeric remain different representations. Employees has no country. Pubs authors/stores have no country. Invoice billing country is not automatically the accepted Customer.Country bridge.

**Concrete dependencies.** `d1-geo-countries`, `d1-geo-places`, `d1-geo-admin1`, `d1-geo-alternates`, `d1-wb-economies`, `d1-wb-indicators`, `d1-wb-observations`, `d2-rs-ember-year`, `d1-cldr-codes`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: `demand-edge-7bccf3b6f28faf7e`, `demand-edge-efccdb52286fa787`, `demand-edge-f9fdae95ad9ea754`, `demand-edge-288524a1a0717baa`, `demand-edge-d504803cd9e2b6ef`, `demand-edge-bb25e8134a3b4384`, `demand-edge-eec245783bf06ed4`, `demand-edge-03d175dc3c91c038`, `demand-edge-50dc7acba519e1ab`, `demand-edge-c4e0b5551a3db181`.
- Accepted public/source assertions: `d1-edge-chinook-country`, `d1-edge-northwind-country`, `d1-edge-northwind-shipping-country`, `d1-edge-pubs-publisher-country`, `d1-edge-place-country`, `d1-edge-place-admin1`, `d1-edge-alternate-place`, `d1-edge-wb-economy-observations`, `d1-edge-wb-indicator-observations`.
- Explicit unresolved assertion dependencies: `d1-edge-geo-wb-economy`, `d1-edge-aw-countrycode`, `d1-edge-cldr-geonames-country`, `d1-edge-sakila-country`, `d2-edge-invoice-electricity`, `d2-map-ember-code`, `d2-map-ember-electricity`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** GeoNames is current geography, WDI annual definitions/revisions and Ember country/aggregate/fuel grain need explicit category filters. No country→economy bridge accepted yet. Chinook 2021–2025, Northwind 2016–2018 and Pubs 1992–1994 require independently measured per-indicator year coverage. No undated GeoNames population denominator, no aggregate economies, no missing=zero; source field mapping rejection stays rejected.

**Reproducible acceptance plan (not yet executed).**
1. Reproduce the four accepted bridge denominators from pinned native bytes; retain every NULL/unmatched row.
2. Have a dedicated reconciler accept/reject exact ISO3→economy and Ember projections; until then show country lookup only and do not label it the full scenario.
3. After acceptance query one pinned population and one economic indicator for explicit source/year plus one non-overlapping Ember metric; assert unique observation grain before aggregation.
4. Negative fixtures: numeric 032 versus 32, unreviewed Sweden-space alias, WB aggregate, historical Yugoslavia and missing year must remain visible.
5. Browser journey on Chinook and Northwind first; add Pubs after B1, then Sakila/AdventureWorks only when their bridges and deployments pass.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Highest cross-demo demand and reusable country/place foundation. Electricity is consolidated into market context because it reuses the same country/year bridge; monthly electricity remains backlog rather than a second launch scenario. Showcase score 4/5. Ranked market tables and aligned year trends; country geometry is deferred, so no boundary map claim.

### p3-scenario-catalogue: Comparable international catalogue, quantities and prices

I connect an international product catalogue with raw units, localized descriptions and explicit currency/date prices; inspect each supported representation; request a comparable-unit/price view with unavailable conversions left blank.

**Useful questions.** Which products carry accepted mass units and which are packaging or unsupported? Which descriptions use the six accepted language-component codes? Where do fixture FX quotes lack a same-day ECB comparison or use obsolete currencies?

**Actual schema and representation boundary.** AdventureWorks joins Product.WeightUnitMeasureCode→UnitMeasure and ProductModelProductDescriptionCulture→Culture using native keys; only exact raw KG-space/LB-space and six padded ar/en/es/fr/he/th values have accepted public bridges. ProductNumber is a local alternate-key lookup, never GTIN. Chinook Track.UnitPrice is catalogue price under explicit USD fixture modelling, not a native currency column. Sakila language.name bridge is unresolved; all original_language_id values NULL. Three demos show relevant structure but one currently has accepted public language/unit bridges. No Employees salary currency or Northwind/Pubs/Sakila amount currency invented.

**Concrete dependencies.** `d1-cldr-units`, `d1-cldr-names`, `d1-cldr-currencies`, `d1-cldr-money-digits`, `d1-ecb-fx`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: `demand-edge-5e4048d05f4a3b4c`, `demand-edge-1ee03f7a65db8683`, `demand-edge-d6684de4716206f5`, `demand-edge-0f035076da34cf64`, `demand-edge-7491d98703cf2f06`, `demand-edge-2bb766331ac76bb4`, `demand-edge-4cf787cba484d389`.
- Accepted public/source assertions: `p2-edge-aw-mass-units`, `p2-edge-aw-language-cultures`, `p2-edge-aw-alternate-product-number`, `p2-edge-aw-alternate-fx-pair-date`, `d1-edge-currency-code-digits`, `d1-mapping-ecb-rate`, `p2-map-chinook-catalogue-price`.
- Explicit unresolved assertion dependencies: `d1-edge-aw-fx-currentref`, `d1-edge-ecb-currency-validity`, `d1-edge-currency-interval-country`, `d1-edge-aw-culture-locale`, `d1-edge-sakila-language`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** ECB indicative quote units/EUR differ from commercial rates; USD→q requires two same-day legs and explicit missing-business-day policy. Recent AW DEM/FRF/VEB quotes remain unsupported. CLDR regional currency intervals can multiply rows; fractions DEFAULT/cash-rounding differ. Conversion requires weight+unit and decimal factor, not unit text alone. Language component is not locale or customer spoken language. SIX rights unresolved and excluded from selected recordsets; no current/historical ISO completeness claim.

**Reproducible acceptance plan (not yet executed).**
1. Prove exact KG-space/LB-space mapping and six-language-domain projection; show raw values beside result; invariant and zh-cht stay unresolved.
2. Measure product→unit native FK and nullable weight; request decimal conversion only after reconciler approves factor/dimension/quantity projection.
3. Have reconciler decide exact date/currency/rate-kind projection; query same-day ECB legs and test inverse direction, absent date, obsolete code and zero/NULL behavior.
4. Choose one CLDR display locale and alt policy so language/currency names do not multiply product rows.
5. Trace one AdventureWorks product through description/culture and weight; exercise a price query separately and retain provenance per column.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Consolidates overlapping money, unit and locale proposals into one international catalogue journey. It avoids redundant SIX/UCUM imports; currency, unit and language remain distinct meanings and query steps. Showcase score 3/5. A side-by-side raw/normalized catalogue with exception counts, currency dates and attributed localized labels.

### p3-scenario-calendar: Local operating calendars with explicit place and time

I connect dated activity plus an explicitly declared jurisdiction/timezone; review date-only versus instant semantics; choose a local-calendar enrichment; see coverage gaps before comparing activity on working days or holidays.

**Useful questions.** Which transactions fall on declared regional bank holidays? How does activity differ by local weekday after a supported timezone conversion? Which historical dates or ambiguous local instants cannot be classified?

**Actual schema and representation boundary.** All six contain date-related structure; this is six-schema demand, zero accepted public demo calendar bridges. Chinook InvoiceDate is date-only 2021–2025 with no inferred timezone; Northwind 2016–2018 and Pubs 1992–1994 fall outside observed GOVUK 2019–2028. Sakila rental/payment timestamps and AW modified/shift dates have distinct roles. Employees interval sentinel 9999-01-01 is not a day to enrich. No country→single-zone or UK→division assumption.

**Concrete dependencies.** `d1-geo-places`, `d1-tz-zones`, `d1-tz-rules`, `d1-cldr-week`, `d1-cldr-windows-zones`, `d1-govuk-holidays`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: none.
- Accepted public/source assertions: `d1-mapping-holiday`.
- Explicit unresolved assertion dependencies: `d1-edge-place-zone`, `d1-edge-windows-zone`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** GOVUK has three divisions, not one UK calendar; region mapping and historical coverage are missing. tzdb civil rules must be evaluated by maintained runtime library; source-rule rows are not fixed offsets. CLDR Windows mapping may be one-to-many. Browser local timezone is not source timezone. Existing historical demos cannot be made compatible by moving their dates.

**Reproducible acceptance plan (not yet executed).**
1. Start with explicit user-schema jurisdiction and IANA zone, then reconcile typed date/instant/interval mappings; no metadata discovery-triggered user row query.
2. Use dated tzdb and two actual offset transitions to test missing/repeated local time; leave unresolved until user supplies policy.
3. Validate division/date holiday keys and explicit coverage; show Northwind/Pubs as out-of-coverage negative journeys.
4. Whole journey on Chinook only if a declared division scope is supplied and reviewed; disclose supplied scenario context as annotation, never native fixture fact.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Calendar/time is a cross-scenario hub for market years, FX dates and environmental observations. High demand but low current accepted demo coverage keeps it behind initial geography/catalogue bridges. Showcase score 4/5. A dated heatmap with explicit coverage gaps and DST counterexamples.

### p3-scenario-organizations: Real organization and affiliation exposure

I connect affiliations carrying ROR IDs or counterparties carrying LEIs; choose an exact namespace lookup; inspect source identity, geographic exposure and explicitly typed consolidation links, with reporting exceptions visible.

**Useful questions.** Where are identified research affiliations located, and what dated market context surrounds them? Which declared accounting-consolidation parents exist for our LEI counterparties? Which parent relationships are unknown, excepted or stale?

**Actual schema and representation boundary.** No actual demo has ROR or LEI keys. This is an explicit broader user-schema opportunity. Synthetic Northwind suppliers, Chinook employers and AW stores are negative examples: names and local IDs never resolve to real entities. ROR→GeoNames→country is accepted public-public connectivity; user ROR/LEI field mapping still needs declared contract and dedicated acceptance. ROR and LEI are two branches connected only by a user-supplied, reviewed organization-to-key relation, not a universal crosswalk.

**Concrete dependencies.** `d1-ror-organisations`, `d1-geo-places`, `d1-geo-countries`, `d1-wb-economies`, `d1-wb-indicators`, `d1-wb-observations`, `d1-lei-entities`, `d1-lei-relationships`, `d1-lei-exceptions`, `d2-rs-ember-year`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: none.
- Accepted public/source assertions: `d1-mapping-ror-id`, `d1-edge-ror-geonames`, `d1-edge-place-country`, `d1-mapping-lei-id`.
- Explicit unresolved assertion dependencies: `d1-edge-lei-child-identity`, `d1-edge-lei-parent-identity`, `d1-edge-geo-wb-economy`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** ROR has multiple locations/relationships; avoid multiplying affiliation counts. LEI relationship/exception CDF rows are documented, not sampled; legal entity and organization roles differ. Accounting consolidation is not all ownership or credit risk. Country energy/economy context remains conditional on S1 bridges; current affiliation location does not reconstruct history.

**Reproducible acceptance plan (not yet executed).**
1. Sample bounded current ROR and LEI schemas/rows, including multi-location and reporting exception; reconcile relation types/time/grain before canonical modelling.
2. Use a transparent hypothetical user table containing real published ROR IDs, with independently justified key mappings; never manufacture demo identity.
3. Retain location multiplicity and count distinct user affiliations separately; show missing IDs, stale revision and unknown parent.
4. Query ROR→place→country first; only accepted subsequent country/year projection enables economic/energy columns.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Adds strong portable real identifiers beyond DemoDB while reusing geography. Consolidates affiliation and supplier proposals without pretending ROR and LEI are equivalent. DOI/ORCID output discovery stays backlog. Showcase score 4/5. An affiliation/place map or typed consolidation graph, never a fabricated universal ownership network.

### p3-scenario-mobility: Premises and transport access context

I connect premises with explicit coordinates or known ATCO/locality IDs; request a scoped transport-stop lookup; inspect locality and dated usage/inspection context, distinguishing proximity from actual route accessibility.

**Useful questions.** Which declared premises are near transport stops within an explicitly chosen distance? Which station-usage observations apply after an audited authority crosswalk? How old are the inspection and usage observations, and where are public addresses suppressed?

**Actual schema and representation boundary.** No actual demo has ATCO/NPTG/ORR/FHRS IDs or validated coordinates; no direct demo join is accepted. Chinook City/Country and Northwind supplier names/postcodes stay negative cases. NaPTAN ATCO strings and NPTG locality keys differ; area010 sample uses admin code009. The accepted 1679/1679 locality probe is scoped to that snapshot, not national accuracy.

**Concrete dependencies.** `d2-rs-naptan`, `d2-rs-nptg`, `d2-rs-orr`, `d2-rs-fsa`, `d1-geo-countries`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: none.
- Accepted public/source assertions: `d2-edge-stop-locality`, `d2-map-naptan-atco`, `d2-map-nptg-id`.
- Explicit unresolved assertion dependencies: `d2-edge-fsa-stops-spatial`, `d2-map-fsa-id`, `d2-edge-customer-locality`, `d2-edge-supplier-ratings`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** ORR NLC/TLC→NaPTAN bridge and OGL scope unresolved; FSA FHRS/FHIS schemes, source ID reuse, private location suppression and freshness conditions matter. NaPTAN↔FSA requires independently approved CRS/axis/radius projection. Approximate distance is not route access or identity. NTA feed rights conflict; do not publish despite five accepted native GTFS edges. Irish GTFS and BODS/TfL are preserved alternate branches, not selected counts.

**Reproducible acceptance plan (not yet executed).**
1. Reproduce area010→NPTG count and namespace counterexample; validate every selected snapshot key/multiplicity.
2. Bounded researcher establishes ORR rights and actual NLC/TLC/ATCO bridge; dedicated reconciler accepts only proven scope or wave excludes ORR and remains incomplete.
3. Reconcile FSA identity validity plus explicit spatial association contract; confirm CRS/axis/order and max-distance behavior with negative points.
4. Use explicit coordinate-bearing hypothetical premises; show dated inspection rating, missing/suppressed location, association distance and no claimed transport travel time.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Compelling connected civic/mobility story but highest bridge/right risk. Rank last for implementation; keep it as a conditional seventh target rather than force a feed-only demo into launch. Showcase score 5/5. Premises/stop maps labeled with distance, dates and missing bridge/identity limitations.

### p3-scenario-environment: Asset observations and environmental context

I connect assets with explicit coordinates or EA measure URIs; inspect typed observations; choose a bounded date range; compare observed station values with separate grid weather context after approved temporal/spatial aggregation.

**Useful questions.** Which readings and units belong to the selected asset/measure and time window? Do valid weather periods coincide with unusual measured activity? Where do missing values, latency or different datums prevent comparison?

**Actual schema and representation boundary.** No validated demo coordinates or EA URIs; no direct demo enrichment accepted. AW spatial BLOBs are opaque. EA reading.measure→measure.@id→station.@id is accepted; source station identity is not transport-stop identity or POWER grid cell identity. URI transport fetch still follows allowlist/redirect policy. Local time and UTC date support remain separate.

**Concrete dependencies.** `d2-rs-ea-stations`, `d2-rs-ea-measures`, `d2-rs-ea-readings`, `d2-rs-power`, `d1-cldr-units`, `d1-tz-zones`, `d1-tz-rules`, `d2-rs-naptan`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: none.
- Accepted public/source assertions: `d2-edge-ea-reading-measure`, `d2-edge-ea-measure-station`, `d2-map-ea-time`.
- Explicit unresolved assertion dependencies: `d2-edge-ea-power-context`, `d2-map-power-temp`, `d2-map-power-rain`, `d2-edge-stops-water-spatial`, `d2-edge-customer-weather`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** EA beta service/no SLA, 15-minute periods, datum/qualifier and negative stage; never use rejected whole readings.value→stage mapping. POWER date-keyed parameter objects require typed unpivot, units and fill -999; observed bounded three days do not prove historical coverage. Nearest station/grid is contextual association, not flood or safety prediction; no causal claim. NOAA global rights and Open-Meteo free-service commercial restriction prevent shortcut substitution.

**Reproducible acceptance plan (not yet executed).**
1. Reproduce EA exact resource references; first implement typed measure selection and independently accepted scalar projection before value meaning.
2. Research/reconcile POWER parameter/time/unit/unpivot, then query explicit coordinate/day; test -999 versus valid negative values and missing day.
3. Compare only approved common temporal support; retain station/datum, grid coordinate, provider revision, latency and aggregation formula.
4. Show no automatic weather query from user data discovery; after explicit request verify bounded network query and per-column provenance.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Reuses unit/time foundations and connects reference metadata to observations. Strong showcase after real typed query value; earthquakes remain optional reserve rather than displacing daily operational utility. Showcase score 5/5. Typed station/weather time series with missing periods, valid negative stage and different temporal support visible.

### p3-scenario-technology: Operational logs and vulnerability identifiers

I connect service/event logs with port+transport, IP protocol, HTTP status or exact media tokens and a vulnerability inventory with CVE IDs; inspect source namespaces; request independent bounded enrichments linked through my declared service key.

**Useful questions.** What do these transport/port and HTTP values mean, including unassigned or ranged entries? Which exact CVEs are listed in KEV and what current NVD analysis is available? Which service rows remain unclassified without guessing from names?

**Actual schema and representation boundary.** No actual DemoDB exposes these portable identifiers; use explicit user-schema contracts. Port+transport is a compound lookup; only tcp/udp→TCP/UDP bridge accepted. HTTP/media branches connect through a user-owned service/event key; there is no invented IANA→CVE public edge. CVE exact identity is accepted public-public; user field contracts still require reconciliation.

**Concrete dependencies.** `d2-rs-iana-ports`, `d2-rs-iana-ip`, `d2-rs-iana-http`, `d2-rs-iana-media`, `d2-rs-cisa`, `d2-rs-nvd`. Dataset/source parents and precise source URLs, rights, cadence and coverage are listed per recordset in the JSON; shared recordsets count once.

**Graph and semantic boundary.**
- Accepted native references: none.
- Accepted public/source assertions: `d2-edge-ports-ip-protocol`, `d2-edge-cisa-nvd`, `d2-map-cisa-cve`, `d2-map-nvd-cve`.
- Explicit unresolved assertion dependencies: `d2-map-iana-port`, `d2-map-iana-ip`, `d2-map-iana-status`, `d2-map-iana-media`.

The JSON `graph_edges` contains exact endpoints/fields, compatibility, transformations and scoped reconciliation reasons for every assertion above. Accepted source field mappings are not counted as executed joins. Hypothetical user fields require explicit scoped contracts and dedicated reconciliation; new missing edges cannot be inferred from path similarity.

**Coverage, rights and operating risks.** NVD reuse/rate policy unresolved blocks NVD publication and complete two-public-source vulnerability journey. CISA CC0 covers catalogue, not linked content/logos. IANA protocol-specific CC0 excludes linked RFC texts. Scalar/range/reserved projection and full media type completeness remain unresolved. KEV absence does not mean safe; dueDate is not a user deadline; names never identify affected products.

**Reproducible acceptance plan (not yet executed).**
1. Reconcile exact scalar/range projection and compound transport-port key; preserve all assignment alternatives without selecting one arbitrarily.
2. Use declared CVE user field, verify exact KEV→NVD match and source revision; rights must pass before NVD import/query.
3. Test missing/invalid CVE, a reserved/ranged port or protocol, duplicate service descriptions and an incomplete media template.
4. Run whole service-log and vulnerability subjourneys separately before a user-owned service join; assert no public cross-edge between unrelated identifiers.

**Expected DataTug behavior.** Apply F1–F6 and the full user journey above: discovery from canonical metadata only; explain the named assertion chain and all gaps; require user-triggered bounded query; display source/version/licence and raw-to-projected values; retain unmatched/ambiguous/multiplied rows. A rejected or unresolved assertion must never acquire a success badge because an adjacent edge works.

**Why and showcase.** Counterbalances demo bias with common real operational schemas and high identifier quality. Two related user journeys share service scope, but are not counted as a fabricated public-data connected component. SPDX rights remain blocked and excluded. Showcase score 3/5. A drillable service/identifier table and KEV timelines with absence and provenance labels.

## Shared foundations and canonical ownership

F1 — **Canonical metadata ingestion.** DataTug still lacks ingestion of the pinned MeaningGraph/ModelSpec/Directory contracts and integration of the HTTPS catalog contract. Implement a rebuildable discovery view, preserve exact canonical pins/addresses, and prove a registry change affects suggestions after invalidation. No new authoritative registry, mutable semantic service or hardcoded launch catalogue.

F2 — **Representation-safe explanation and planning.** Use core meanings already accepted at `982916d73f0a35ff2558b0062f58aa3ac4f24d97` (23 concepts across six files). Missing reusable timezone/instant/interval/unit/observation/identifier concepts require accepted dedicated reconciliation and narrowed definitions before canonical creation. Same-meaning fields are not automatically equal keys. Display identifier namespace, raw role, transforms, confidence and license. Compile selected explicit bridges into ordinary query/provider projections; no universal transform language.

F3 — **Existing execution path.** Browser OVDB source resolution plus DALgo joins remains the default. Use selected source mounts or normal queryable projections and existing allowlist/redirect controls; cap rows, bytes, execution time and remote refresh. Do not design CLI/multiserver support into the critical path without an actual selected journey needing it. Source main `openvaultdb/openvaultdb-go@2ed8748df9c1dfe5ced8127c29e53251c2594c66` adds database-as-source capability, but the B1 refresh shows deployed Cloud pins still older; this is not production availability. Refresh capability/consumer contracts before implementation.

F4 — **Mismatch and cardinality.** Keep raw source rows and original identifiers. Measure total, nonnull, NULL, empty, matched, unmatched, ambiguous and multiplied counts separately. Grain must be fixed before aggregate comparison. Do not silently trim, case/diacritic fold, fill dates, assign currencies, drop NULL countries or select an arbitrary location/zone/parent.

F5 — **Rights, freshness and provenance.** Before each import/query, verify the specific source material, licence, attribution, commercial access, caching/snapshot/adaptation/redistribution scope, excluded third parties and endpoint conditions. A linked URL or successful source probe is not permission. Version source, parser and projection; record fetch time/hash, upstream observation/release time separately, temporal/geographic coverage and revalidation cadence. Last-known-good snapshots must be labeled stale; a failed refresh never silently becomes current.

F6 — **Real product evidence.** Verify whole browser connect → explained suggestion → explicit bounded OVDB query → visible exceptions/provenance → reuse journey. Inspect actual network execution, not just source probes. Golden positive paths plus negative controls must run against the exact reviewed app/provider/runtime revisions. Failure leaves the scenario conditional; structural validators alone cannot pass Phase5.

| Canonical responsibility | Existing repository / mechanism | Verification basis and next action |
|---|---|---|
| Universal meanings and contract | `meaninggraph/core`; graph validator `meaninggraph/cli` | Reuse Phase2 exact pins; new concept only after dedicated decision and normal canonical review. |
| Dataset-local bindings, concrete models, source manifests/provenance | Existing provider repositories, including `demo-db/{chinook,northwind,pubs,sakila,adventureworks,employees}` | Exact source/native hashes are in `p2-inputs-and-tools.json`; do not choose a new provider repository from a made-up name. Resolve publisher ownership for each new source during the approved implementation plan. |
| Meaning publication registry | `meaninggraph/registry`, `graphs/$records`, generated `index.json` | Fresh fetched `origin/main` `2d92bbdd45fc2552f231f8f993a5b8b8d5b7fa29` index confirms all six providers. README table is stale; follow-up retained. |
| Model publication registry | `modelspec-org/registry`, model records and generated index | Fresh primary README confirms Git authority; Phase2 registry pin `48b30b250a61385d94d46a70968677870750ea35` has all six. Validate source HCL and exported JSON with `modelspec-org/cli`; specification belongs to `specscore/modelspec`. |
| OVDB canonical database directory | `openvaultdb/directory`, `databases/$records`, generated `index.json` | Fresh canonical fetch `02db362144d7924c6081dd6768cc0d7187ac9bc3`; website is derived. Local misnamed `/openvaultdb/registry` has the same origin but is a legacy/divergent clone, not another authority or edit target. |
| Runtime query/mount/config | `openvaultdb/openvaultdb-go`, `openvaultdb/cloud`, DALgo adapters | Source capabilities and deployed consumer pins differ; B1 owner must publish/repin/deploy and independently validate. |
| User discovery/explanation/query | `datatug/datatug-apps`, shared contracts as needed in `datatug/datatug-core` | F1–F6 are real Phase4 work; current app capability is not implied by research. CLI is conditional. |
| Portfolio and implementation-plan ownership | Hub feature and these existing research collections; formal implementation plans follow `sneat-co/workbench/spec/plans/github.com/datatug/datatug/` | This is an architecture proposal linked from the existing feature. No parallel `spec/plans` tree or unapproved implementation dispatch is created. |

## Source integration modes and freshness contract

The JSON has one `recordset_integration_modes` entry for each of the 37 unique candidates, retaining the exact candidate/source IDs and recorded licence text. Those are research observations from 4 October, not a new legal ruling or current-data claim. Planned handling is reproducible attributed snapshots/refresh through OVDB where terms permit. GeoNames, WDI, Ember and POWER use their recorded CC BY 4.0 scopes; Unicode data needs its copyright/permission notice; ROR metadata and GLEIF/IANA protocol registry material use their recorded CC0 scopes; the selected ROR release explicitly identifies embedded GeoNames location data as CC BY 4.0, so its attribution must also be retained; tzdb is public domain; GOVUK, NaPTAN/NPTG, EA and FSA use source-specific OGL terms. ECB needs source credit, accurate reproduction, modification marking and its recorded sold-document condition. FSA requires a current rating or explicit update date and excludes protected imagery. NVD and ORR stay blocked, including assumed remote query access, until exact terms are accepted. All terms must be rechecked for the actual release/artifact before implementation.

Use the existing manifest/source/notes ownership; do not add lifecycle/confidence fields to closed canonical schemas merely to copy research metadata. Daily or volatile sources need bounded refresh with hashes, structural/key/grain diffs and a last-good pin; annual observations need edition/revision checks and explicit missing periods; versioned reference sources need release pins and migration review. Source-specific cadence remains in the ledger (unknown remains unknown). Implementation chooses concrete schedules only after payload and rate limits are measured, not by turning observation dates into freshness guarantees.

## W1 operational scope: useful global reference coverage

The smaller Colombia/ROR fetch below is only a plumbing proof. It is not the selected coverage and cannot pass a launch scenario. Selected W1 uses global country/admin references plus a deterministic global place subset with complete ROR-reference closure. All proposed budgets below are architecture limits to measure during implementation, not demonstrated runtime performance.

The [GeoNames listing](https://download.geonames.org/export/dump/) currently advertises countryInfo about 31 KB, admin1 about 148 KB, cities5000 about 5.4 MB compressed, allCountries about 402 MB and alternateNamesV2 about 195 MB. The [pinned ROR release 2.13](https://zenodo.org/records/22902037), dated 22 September 2026, reports 141,528 organizations and 37.5 MB ZIP; the exact file is `v2.13-2026-09-22-ror-data.zip`. This release says ROR metadata is CC0 and embedded GeoNames locations are CC BY4. JSON is the source representation, with all source statuses retained.

- **Countries:** complete `countryInfo.txt`; fresh probe 252 unique ISO rows / 31,678 bytes. This preserves all four accepted demo country-label paths independently of place coverage.
- **Admin1:** complete `admin1CodesASCII.txt`; fresh probe 3,865 unique composite codes / 151,583 bytes. Keep missing admin references; a gazetteer can contain valid source values absent from this reference snapshot.
- **Organizations:** complete pinned ROR JSON release. Preserve inactive/withdrawn organizations, all location references and original relationship arrays. The exact ID lookup remains useful across real user affiliation schemas; no names are converted to identity.
- **Places:** union of source-defined global `cities5000` (about 50,000 urban/admin-seat features) and every nonnull GeoNames location ID referenced by that full ROR release. Stream `allCountries.zip` once to extract this closure. A referenced ID absent upstream stays an explicit dangling reference with its ROR row preserved. This is not a complete global gazetteer or accepted demo geocoder.
- **Aliases:** stream `alternateNamesV2.zip` and retain all rows whose GeoNames ID belongs to the included place set. Preserve stable alias IDs, language, preferred/short/historical flags and from/to; sort output deterministically. The accepted alias→place reference does not accept arbitrary user-text→place matching.

The input plan is approximately 640 MB compressed per complete rebuild. Proposed limits: 1 GiB download total, 8 GiB temporary disk, 512 MiB process RSS, 30 minutes. Use streamed parsing and disk-backed key indexes/ZIP member reads; do not load the global TSVs or full JSON into memory. Generated provider database target ≤256 MiB with indexes on exact keys; metadata ≤2 MiB. No raw large artifacts belong in the source repository or browser. Stop and rerank if actual measurements exceed these limits; no silent row removal or budget increase. Snapshot hosting uses the existing OVDB mechanisms and does not assume a new paid API, account or service.

A first browser query is user-triggered and limited to 1,000 user rows, 50 distinct ROR IDs or 100 place IDs, 5,000 result rows, 5 MiB response and 10 seconds. These are proposed acceptance budgets, to be implemented on both query planning and execution. Pagination needs an explicit next request; truncation and unmatched/multiplied rows remain visible. Querying never triggers a global rebuild. Initial refresh is a reviewed monthly snapshot/release rebuild, labeled with its GeoNames date; do not call it daily-current. Upstream version/header monitoring can be configured in implementation, not created by this phase.

Freeze the captured GeoNames 2026-10-04 batch by SHA256 and HTTP metadata, plus exact RORrelease 2.13. If mutable source bytes have changed, freeze a new coherent batch and rerun validation; never invent a provider version. Record compressed and extracted-member hashes, importer revision, coverage predicate and upstream terms. Rebuilding must preserve native keys and deterministic output; changed/deleted IDs, source-status changes and newly missing references are review inputs before repinning the registry.

[The bounded operational probe](evidence/p3-w1-operational-probe.json) fetched full Colombia archives (2,248,873/659,217 compressed bytes) and one known ROR record, with zero redirects and a 4 MiB per-response cap. It observed 82,094 unique places,62,005 unique 10-column alias rows, zero alias→place misses and the existing seed `02mhbdp94→3688689→CO.34→CO`. It also found 475 nonempty place→admin1 references absent from the complete admin1 file. Keep that as an expected mismatch class; do not publish a zero-orphan claim. The probe did not download the global archives or full ROR release and establishes neither global closure nor semantics.

Before W1 completion: measure full ROR closure and every retained/missing key; prove the four actual demo country denominators; independently check the seed and at least two deterministic additional-country examples; verify source licence/attribution, resource bounds, immutable generated artifact integrity and actual browser query/provenance. All F1–F6 and B1 gates still apply. If the global build is not economical, a smaller global country-only step is allowed with explicit reduced scope; five recordsets is guidance, not a reason to turn the one-organization plumbing fixture into launch coverage.

## Proposed waves and acceptance gates

| Wave | New distinct public recordsets | Scenario payoff and condition |
|---|---:|---|
| W1 Country/place and keyed affiliation foundation | 5: `d1-geo-countries`, `d1-geo-places`, `d1-geo-admin1`, `d1-geo-alternates`, `d1-ror-organisations` | Five connected recordsets with accepted provider references/bridges. Delivers country lookup and exact ROR→place→country subjourneys; not the complete dated market or LEI scenario. Canonical source-specific binding fidelity, user-key contracts, rights revalidation, F1–F6 and B1 remain required. |
| W2 Dated market and organization exposure | 8: `d1-wb-economies`, `d1-wb-indicators`, `d1-wb-observations`, `d2-rs-ember-year`, `d1-cldr-codes`, `d1-lei-entities`, `d1-lei-relationships`, `d1-lei-exceptions` | Two bounded packets: five market/category recordsets after dedicated country/economy/Ember acceptance; three LEI recordsets after CDF/relationship acceptance. Reuses W1 geography/ROR; no ROR↔LEI equivalence. |
| W3 Comparable international catalogue | 5: `d1-cldr-units`, `d1-cldr-names`, `d1-cldr-currencies`, `d1-cldr-money-digits`, `d1-ecb-fx` | Accepted padded unit/language subset first; complete quantity/FX projections still need reconciliation. |
| W4 Place and local calendar | 5: `d1-tz-zones`, `d1-tz-rules`, `d1-cldr-week`, `d1-cldr-windows-zones`, `d1-govuk-holidays` | Place/zone/division and historical coverage must be explicit; do not declare six-demo calendar readiness. |
| W5 Operational technology | 6: `d2-rs-iana-ports`, `d2-rs-iana-ip`, `d2-rs-iana-http`, `d2-rs-iana-media`, `d2-rs-cisa`, `d2-rs-nvd` | NVD terms and scalar/range projections gate the full scenario; no service→CVE public identity inferred. |
| W6 Typed environmental observation | 4: `d2-rs-ea-stations`, `d2-rs-ea-measures`, `d2-rs-ea-readings`, `d2-rs-power` | Four natural connected recordsets; no quota filler. Requires accepted typed projections and spatial/time context. |
| W7 Premises and mobility | 4: `d2-rs-naptan`, `d2-rs-nptg`, `d2-rs-orr`, `d2-rs-fsa` | Four natural recordsets; do not replace missing ORR/FSA bridges with rights-blocked GTFS or fabricate reachability. |

W1 is the five-recordset connected foundation supported by accepted Phase2 public assertions. It makes no claim that full market or organization scenarios are complete. W2–W4 reuse geography/time/unit hubs; rerank after W1 and split the two W2 packets if only one closes its prerequisites. W5 is independent enough to move earlier if its rights close while spatial work stalls. W6 and W7 each naturally contain four new recordsets; this is an explicit exception to the approximate 5–10 guidance, avoiding filler. Their dependency on earlier reference sets makes the full scenario larger. The W6 optional stop-context extension waits for W7 and is not necessary for the core EA→POWER journey. Every wave is conditional and remains unimplemented.

For each wave: (a) freeze exact recordset scope and rights; (b) dedicated reconciler accepts all required field/bridge/projection semantics; (c) canonical provider/MG/ModelSpec artifacts validate and registries pin the reviewed commits; (d) runtime mount/HTTPS capability and source freshness pass; (e) explicit DataTug positive/negative journey returns measured results with source attribution; (f) fresh independent reviewer checks semantics, rights, cardinality, publication and product path; (g) root accepts/rejects, records receipts and reranks. No phase gate passes while B1 or another assigned blocker is open. Before phase-gate acceptance, root must bind any selected unresolved integration conditions to the existing questions/gate mechanism rather than treating this prose as satisfied.

**B1 remains OPEN/BLOCKING for phases 4/5/6.** The separate DemoDB delivery owner retains an unanswered publication authorization for `dal-go/dalgo2sql`, `openvaultdb/directory`, releases and consumer deployment. This lane does not publish, replace the owner or message another chat. Latest retained refresh has AW/Employees metadata 404, Sakila/AW/Employees live Directory 404, stale Pubs live pin, Sakila advertised DTQL 404/capability mismatch, older Cloud consumers and final immutable-native acceptance open. All-six canonical registrations do not close those runtime/live obligations. The canonical Chinook host remains `chinook.demodb.dev`.

## Targeted follow-up and preserved backlog

Questions `p3-followup-1` through `p3-followup-6` name bounded work and payoff. Their classification is nonblocking to Phase3 architecture; affected implementation is expressly gated. Root has commissioned separate `f3-*` evidence work for country/economy/Ember and demo representation gaps; those future receipts are not silently incorporated or treated as semantic acceptance here. Dedicated semantic follow-up must judge them separately.

- Priority 1: Bounded exact category/representation evidence, followed by separate semantic specialist decisions. Payoff: Unlock S1/W1, S4 exposure and actual additional demos. Bound: Complete small code/category snapshots and seven Sakila exceptions; no bulk economic import; no blanket bridge. IDs: `d1-edge-geo-wb-economy`, `d2-map-ember-code`, `d2-map-ember-electricity`, `d1-edge-sakila-country`, `d1-edge-aw-countrycode`.
- Priority 1: Which date/rate-kind/currency projections and same-day two-leg comparison are supported? Payoff: Complete S2 explicit-price comparisons without guessed currency. Bound: 14 AW target currencies, one supported pair/day, one missing business day and DEM/FRF/VEB counterexamples; SIX rights separate. IDs: `d1-edge-aw-fx-currentref`, `d1-edge-ecb-currency-validity`.
- Priority 1: Can one declared jurisdiction and zone support dated activity, and which historical calendars are licensed for Northwind/Pubs periods? Payoff: S3, future workforce and mobility. Bound: One supported division/zone journey; no global holiday scrape. Preserve out-of-coverage negatives. IDs: `d1-edge-place-zone`, `d1-edge-windows-zone`, `d1-govuk-holidays`.
- Priority 1: Actual LEI relation/exception CDF grain and NVD source-specific reuse/rate terms? Payoff: Complete independently useful S4 and S5 branches. Bound: One relation, one exception and official NVD terms; no legal inference from successful API response. IDs: `d1-lei-relationships`, `d1-lei-exceptions`, `d2-rs-nvd`.
- Priority 1: Typed EA measure and POWER scalar unpivot with compatible daily support? Payoff: S6 four-recordset connected observations. Bound: One EA parameter/datum, two POWER parameters and three UTC days; replacement assertion for rejected unqualified stage mapping. IDs: `d2-map-power-temp`, `d2-map-power-rain`, `d2-map-ea-stage`, `d2-edge-ea-power-context`.
- Priority 2: ORR licence and NLC/TLC-to-ATCO bridge; source-valid FSA ID and declared spatial association? Payoff: S7 service-access context. Bound: One rail area, one annual period, one local authority, no nationwide geocoder or identity from premises names. IDs: `d2-data-orr`, `d2-rs-fsa`, `d2-edge-fsa-stops-spatial`.
- Priority 2: Historical country/area/industry/occupation editions and salary units for real branches? Payoff: Strong next portfolio workforce/industry scenario; better use of Employees only if context exists. Bound: One explicit user schema and one licensed, accessible series; do not manufacture currency/country for fabricated employees. IDs: `d1-dataset-bls`, `d1-dataset-census`, `d1-cpi-series`, `d1-cpi-observations`, `d1-naics-codes`, `d1-oews-observations`, `d1-cbp-observations`.
- Priority 2: Authoritative terms and feed/usage namespaces before extending mobility? Payoff: Alternative mobility supply/demand branches. Bound: Resolve documented NTA restrictions first; accepted native GTFS references do not grant rights; missing/empty GTFS files are not integrations. IDs: `d2-data-nta`, `d2-data-bods`, `d2-data-tfl`.
- Priority 2: Do polygons/postcodes/admin2/hierarchy/language distribution/monthly electricity add measurable marginal utility? Payoff: Map, locality and historical-resolution extensions after core journeys. Bound: Inspect actual geometry/key/grain or source-specific postal rights; current-year timezone offsets never replace tzdb; language population is not spoken-language identity. IDs: `d1-dataset-naturalearth`, `d1-dataset-geonames-postal`, `d1-geo-admin2`, `d1-geo-hierarchy`, `d1-geo-zones`, `d1-cldr-languages`, `d2-data-ember-month`.
- Priority 3: Resolve source-specific publication, transformed-data and commercial API constraints before considering implementation. Payoff: Useful specialist extensions; none bypasses present rights gate. Bound: One source decision at a time; remote linking/query is not automatically permitted; no TPC-H re-research. IDs: `d1-dataset-six`, `d1-dataset-ucum`, `d1-dataset-un-m49`, `d1-dataset-comtrade`, `d2-data-spdx`, `d2-data-ghcnd`, `d2-data-openaq`, `d2-data-openmeteo`, `d2-data-gias`.
- Priority 3: Can these shallow/showcase/quality-limited candidates beat current marginal utility? Payoff: Alternative economics, identifier bridges, geophysical/civic maps. Bound: OECD duplicates WDI unless new required coverage; Wikidata claims need provenance; USGS third-party contributor rights; BTS stale/truncated; empty police sample proves no observation coverage. IDs: `d1-dataset-oecd`, `d1-dataset-wikidata`, `d2-data-usgs`, `d2-data-bts`, `d2-data-police`.

All 57 original scenario/bundle rows, including the 21 conceptual relationship bundles, are retained below and in JSON. Consolidation means overlapping intent was routed into a stronger opportunity; it does not silently include every original member in launch scope. Missing DOI/ORCID, GTIN/classification, workforce, health, postal/geocoder, GTFS or historic-calendar branches remain backlog. Original lifecycle/evidence and every assertion byte are unchanged. Three ROR candidate licence fields receive a scope-qualified update and appended evidence; their original text is retained in `preservation.rights_clarifications`. All other original candidates have a separate planning disposition in `candidate_dispositions`; original semantic disposition and source rights remain authoritative.

| Original ID | Disposition | Target / reason |
|---|---|---|
| `d1-scenario-branch-benchmarks` | prioritized-backlog | P1 Historical US branch/labour benchmarking is valuable but no sampled observation/portable demo GEOID, NAICS/SOC or salary currency; rights/access and editions first. |
| `d1-scenario-calendar-context` | consolidated-proposal | p3-scenario-calendar |
| `d1-scenario-fx-audit` | consolidated-proposal | p3-scenario-catalogue |
| `d1-scenario-locale-units` | consolidated-proposal | p3-scenario-catalogue |
| `d1-scenario-market-geography` | consolidated-proposal | p3-scenario-markets |
| `d1-scenario-real-organisation` | consolidated-proposal | p3-scenario-organizations |
| `d2-bundle-civic` | consolidated-proposal | p3-scenario-mobility |
| `d2-bundle-environment` | consolidated-proposal | p3-scenario-environment |
| `d2-bundle-mobility` | consolidated-proposal | p3-scenario-mobility |
| `d2-bundle-technology` | consolidated-proposal | p3-scenario-technology |
| `d2-scenario-civic-services` | consolidated-proposal | p3-scenario-mobility |
| `d2-scenario-electricity` | consolidated-proposal | p3-scenario-markets |
| `d2-scenario-geophysical` | consolidated-proposal | p3-scenario-environment |
| `d2-scenario-internet` | consolidated-proposal | p3-scenario-technology |
| `d2-scenario-irish` | consolidated-proposal | p3-scenario-mobility |
| `d2-scenario-transport` | consolidated-proposal | p3-scenario-mobility |
| `d2-scenario-vulnerabilities` | consolidated-proposal | p3-scenario-technology |
| `d2-scenario-water-weather` | consolidated-proposal | p3-scenario-environment |
| `demand-scenario-amount-currency-units` | consolidated-proposal | p3-scenario-catalogue |
| `demand-scenario-employment-validity` | prioritized-backlog | P1 Employees interval/occupation demand retained, but currency/country/classification assignments and historical sources are absent. |
| `demand-scenario-geography-context` | consolidated-proposal | p3-scenario-markets |
| `demand-scenario-language-locale` | consolidated-proposal | p3-scenario-catalogue |
| `demand-scenario-media-classifications` | consolidated-proposal | p3-scenario-catalogue |
| `demand-scenario-product-units-inventory` | consolidated-proposal | p3-scenario-catalogue |
| `demand-scenario-transaction-calendar` | consolidated-proposal | p3-scenario-calendar |
| `m1-edge-address-site` | consolidated-proposal | p3-scenario-organizations |
| `m1-edge-class-instance` | prioritized-backlog | P2 Conceptual item-instance relation retained; no portable instance source/identifier assignment. |
| `m1-edge-class-product` | prioritized-backlog | P2 Product classifier edition and assignment evidence absent; local SKU/category cannot substitute. |
| `m1-edge-country-region` | consolidated-proposal | p3-scenario-markets |
| `m1-edge-energy-product-flow` | consolidated-proposal | p3-scenario-environment |
| `m1-edge-entity-parent` | consolidated-proposal | p3-scenario-organizations |
| `m1-edge-health-population` | prioritized-backlog | P2 Health-rate denominator/geography/cohort relation requires new evidence and reconciliation. |
| `m1-edge-instant-zone` | consolidated-proposal | p3-scenario-calendar |
| `m1-edge-membership-unit` | prioritized-backlog | P1 Native employee membership is local; workforce public benchmark links and periods remain missing. |
| `m1-edge-population-area` | consolidated-proposal | p3-scenario-markets |
| `m1-edge-postcode-locality` | consolidated-proposal | p3-scenario-markets |
| `m1-edge-rate-currency` | consolidated-proposal | p3-scenario-catalogue |
| `m1-edge-registration-entity` | consolidated-proposal | p3-scenario-organizations |
| `m1-edge-ror-domain` | consolidated-proposal | p3-scenario-organizations |
| `m1-edge-station-observation` | consolidated-proposal | p3-scenario-environment |
| `m1-edge-subdivision-country` | consolidated-proposal | p3-scenario-markets |
| `m1-edge-tag-language` | consolidated-proposal | p3-scenario-catalogue |
| `m1-edge-transport-city` | consolidated-proposal | p3-scenario-mobility |
| `m1-edge-trip-route` | consolidated-proposal | p3-scenario-mobility |
| `m1-edge-trip-stop` | consolidated-proposal | p3-scenario-mobility |
| `m1-edge-unit-property` | consolidated-proposal | p3-scenario-catalogue |
| `m1-scenario-calendar-demand` | consolidated-proposal | p3-scenario-calendar |
| `m1-scenario-civic-health` | prioritized-backlog | P2 Health denominator/cohort/boundary meanings and suitable licensed observation sources absent; do not repurpose country population. |
| `m1-scenario-energy-environment` | consolidated-proposal | p3-scenario-environment |
| `m1-scenario-international-content` | consolidated-proposal | p3-scenario-catalogue |
| `m1-scenario-market-context` | consolidated-proposal | p3-scenario-markets |
| `m1-scenario-mobility-context` | consolidated-proposal | p3-scenario-mobility |
| `m1-scenario-price-comparison` | consolidated-proposal | p3-scenario-catalogue |
| `m1-scenario-product-logistics` | prioritized-backlog | P2 GTIN/item/packaging/shipment identifiers absent from demos; classification assignments and source inventory are shallow. |
| `m1-scenario-research-outputs` | consolidated-proposal | p3-scenario-organizations |
| `m1-scenario-vendor-network` | consolidated-proposal | p3-scenario-organizations |
| `m1-scenario-workforce` | prioritized-backlog | P1 Meaning-rich workforce opportunity awaits occupation/area/unit assignments and bounded historical BLS/Census access. |

## Reranking and stopping

After each accepted wave, recompute marginal benefit using actually usable new paths, measured match/multiplicity/coverage and maintenance cost. Reuse earns value only if a required bridge has passed; hypothetical demo column prevalence does not. Prioritize one bounded follow-up when it unlocks at least two scenarios or an additional structurally different actual demo; stop/replan that branch if the evidence rejects the needed identity, rights remain unclear, coverage misses the intended time/place, or cost exceeds the next ready scenario. Keep the rejected/blocked evidence and reroute to a narrower truthful journey.

Stop initial launch when root has five to seven excellent independently validated journeys if evidence supports them, foundational cross-demo meanings work through several physical representations, canonical changes affect product suggestions, all selected rights/provenance/query/publication gates pass, and remaining backlog has lower marginal value. Do not claim launch completion merely because 37 recordsets were imported or seven proposal cards exist. If fewer scenarios are defensible, root records a scope decision with evidence. Any move toward about 75 actual implemented recordsets requires explicit portfolio review; growth from physical splits also counts.

## Author verification and phase boundary

The separate author validation receipt is `evidence/p3-author-validation.json`. Run the existing inGitDB validation, 24 ledger mutation cases, three query fixtures, two native-probe counterexamples, SpecScore lint and whitespace check. A conservation audit verifies all original 1,170 assertion hashes and 679 unchanged candidate hashes, plus the exact three ROR rights-only deltas, all 57 dispositions, every referenced accepted/unresolved assertion, exact 37-set union and per-wave uniqueness. This checks references and preservation, not correctness of architecture or semantics.

No gate record is modified and no status transitions to selected/implemented/validated occur. The existing owning feature links this proposal while remaining Draft. Root must obtain a fresh independent adversarial review of the exact PR head, resolve findings, and decide Phase3 acceptance separately. Source probes, author checks and this separate-architect receipt cannot substitute for that review.
