# Demo corpus demand analysis — Phase 0.5 r1

Observation: 2026-10-04. **Proposed research, pending independent review and root gate acceptance.** No semantic reconciliation, canonical meanings, selected portfolio, accepted cross-demo joins or live publication completion is claimed.

The desired journey remains connect a demo → recognize a contextual field role → suggest explained public references → execute a measured lookup → reuse that knowledge across different schemas. This analysis supplies evidence for recognition and reference demand. Source selection, fresh reconciliation, canonical ingestion, actual joins and product validation remain later work.

## Actual inventory and reproduction

Freshly fetched `demo-db/websites` at `27967f26fa87e67c8d7a7b3d2784a0cc4c1bfdf6`, `config/databases.json`, names exactly these six fixtures. Each catalogue provider pin matches the freshly fetched provider `origin/main`; each contract hash matches its immutable Git blob. TPC-H remains excluded under the preserved permission decision. No renewed permission research was performed.

| Fixture | Native tables | Executable views | Physical rows | Data / model / existing meaning licence |
|---|---:|---:|---:|---|
| Chinook | 11 | 0 | 15,607 | MIT / MIT / CC0-1.0 |
| Northwind | 13 | 17 | 3,310 | MIT / MIT / CC0-1.0 |
| Pubs | 11 | 1 | 255 | MIT / MIT / CC0-1.0 |
| Sakila | 16 | 7 | 47,268 | BSD-3-Clause / BSD-3-Clause / CC0-1.0 |
| AdventureWorks | 71 | 11 | 759,240 | MIT / MIT / CC0-1.0 |
| Employees | 6 | 2 | 13,584 | CC-BY-SA-3.0 / CC-BY-SA-3.0 / CC-BY-SA-3.0 |

[The input plan](evidence/demand-plan.json) records current provider pins, native SQLite hashes, upstream revisions, conversion/subset limits and the separate licences. [Native observations](evidence/demo-demand.json) contain all 1,120 physical/derived column definitions, ordered keys, declared local FK paths and indexes, with bounded value profiles for 865 distinct fields. Unprofiled derived-view fields retain schema evidence; no all-column value-profiling claim is made. The 128 physical tables and 38 executable views are represented as 166 research recordsets. Nine AdventureWorks source-only views remain outside executable coverage; their preservation is documented in the pinned provider provenance. Native source data is the authority here: its SQLite differs from AdventureWorks' delivery artifact with appended generated `id` fields.

The old Phase 0 extraction and `snapshot-pins.json` remain immutable historical evidence. Sakila and Employees current pins differ from those older records; new observations carry their exact current pins rather than silently rewriting history. This analysis makes no public deployment or current registry acceptance claim. B1 `corpus-publication` remains open at Phase 4/5/6 with the existing DemoDB delivery owner.

`inspect_demand.py` verifies every native SQLite SHA256 before opening `mode=ro&immutable=1` with `query_only`. Provide a JSON object mapping demo IDs to local pinned native SQLite paths. For AdventureWorks/Employees, decode their pinned `data-source/source.sqlite.gz` into a temporary file; the pinned decompressed digest is mandatory. Other native inputs are `data-source/source.sqlite`. No artifact download, browser, service or mutation is required.

```sh
python3 spec/research/public-data-fabric/inspect_demand.py \
  --fixtures /path/to/local-fixture-paths.json --output /tmp/demo-demand-reproduced.json
cmp /tmp/demo-demand-reproduced.json spec/research/public-data-fabric/evidence/demo-demand.json
```

Profiles report declared and actual storage types; BINARY distinct counts; NULL/empty, leading numeric-zero text, outer ASCII-space counts, ASCII case collisions; bounded literal values and raw bounds. Country/currency/language/region/unit/payment-term fields include all distinct values when at most 256. Long text uses a labelled prefix/length/hash, BLOBs use size/hash, and authentication/card/identity examples are suppressed. Date bounds are literal lexical bounds, not normalized timestamps or historical coverage guarantees. SQLite `lower()` only folds ASCII; no Unicode or timezone normalization is claimed.

## Demand and counterexamples

The index retains **162 proposed candidate groups and 937 unique field-role hypotheses**. These are descriptive research counts. [The coverage matrix](evidence/demo-demand-coverage.json) includes every hypothesis group, all six demo cells, native fields, role, assertion IDs, recordset IDs, row counts and nonnull availability. Empty/all-null fields and derived views remain distinguishable. `to_fields: ["proposed-role"]` in mapping assertions denotes a research role hypothesis, not a published canonical target property. All compatibility values remain `unresolved`.

Country and postal-code hypotheses each occur in five demos; employee and person-name roles in six; occupation in five; customer in four. Shared occurrence motivates reference research, without establishing entity identity or semantic equivalence. Related researcher terms such as city/locality, address/postal-address, region/geographic-region and language/culture-locale remain separate. Broad source-specific attribute groups have lower priority and confidence; they remain available to the reconciler.

- **Contextual geography is the strongest practical bridge.** Chinook and Northwind carry country labels, Northwind employees use `UK`/`USA`, Pubs publishers carry `France`/`Germany`/`USA`, Sakila needs `address.city_id → city.country_id → country.country`, and AdventureWorks carries source country codes plus `Address.StateProvinceID → StateProvince.CountryRegionCode`. These paths preserve local FK namespaces. AdventureWorks' 238 country codes include historical/territorial ambiguity; ISO-like shape alone does not establish modern ISO identity.
- **Raw dirt is part of the demand.** Chinook has `Edinburgh ` in one customer and seven invoices. Northwind has `CustomerID='Val2 '` and `Suppliers.Country='Sweden '`. AdventureWorks has 176 padded state codes, 23 padded unit codes and seven padded culture codes, including whitespace-only invariant culture. Pubs has literal `M�nchen` replacement characters. Postal text retains leading zeros; Sakila has four empty postcodes. No trimming, repair, translation or approximate join has been accepted.
- **Place roles require context.** Northwind customer/shipping `Region` uses broad European areas; supplier regions mix those with `NSW`/`Victoria`. Its separate Regions/Territories are sales assignments, and postal-shaped TerritoryID keys are not proven postcodes. Pubs authors/stores have no Country column. Customer, employee, billing, shipping, supplier and publisher locations stay separate. AdventureWorks geography/hierarchy BLOBs are opaque, not decoded coordinates or paths.
- **Calendar demand spans schemas but events differ.** Chinook invoices span 2021-01-01–2025-12-22; Northwind orders 2016-07-04–2018-05-06; Pubs sales 1992-06-15–1994-09-14. Pubs publication maximum 2000-01-01 includes two frozen omitted defaults. Employees `9999-01-01` is a current-validity sentinel, not an observation in year 9999. Employment, order, shipment, rental, payment, publication, effective-period and audit dates are separately qualified. Missing timezone and fixture subset limits survive.
- **Currency/units are a valuable hub with restricted immediate coverage.** AdventureWorks explicitly contains 105 currency codes and 13,532 dated FX rows (2022-05-30–2025-05-30; USD source currency, 14 target currencies), including historical codes such as AFA/DEM. This does not give every money field a currency. Chinook's existing meaning declaration says USD while raw data has no currency field; the declaration and raw absence remain distinct. Northwind/Pubs/Sakila/Employees also have no raw currency field. Pubs discount percentages differ from Northwind fractional discounts; source NUMERIC conversions do not guarantee decimal precision.
- **People and organisations need role-aware matching.** Customer, supplier, publisher, store, employee, author, performer, contact and business entity do not share an ID namespace. Names/phone/email/SKUs can motivate reference or classification lookup but do not establish global identity. Employees has 1,024 selected fabricated people in an FK-closed subset, not the 300,024-person upstream population. Compensation, occupations, departments and shifts need date/role/unit qualifiers.
- **Media, commerce and classification remain reusable research.** Chinook genre/media format/track duration, Pubs publication type/authorship/royalty/payment terms, Sakila film rating/SET-valued features/language, and AdventureWorks product units/category/BOM/inventory/tax/shipment/offer data expose useful hubs. Film primary language is English for all 1,000 films; original language is NULL throughout. Payment terms are not payment transactions; rentals, invoices, orders and sales lines remain separate. Primary and original language, and language versus culture/locale, remain distinct.

Native composite keys include Northwind `Order Details(OrderID,ProductID)`, Chinook playlist membership, Pubs sale/authorship and Employees history/assignment keys. Pubs discounts and royalty schedules remain keyless: their whole-row hypotheses list actual columns solely as evidence, without constructing an identifier. Local FK assertions document declared structure only; measured match rates, multiplicity and portable bridges remain scenario validation work.

## Research handoff

Seven proposed scenario candidates preserve contextual geography, transaction calendars, employment validity, amounts/currency/rates, media classifications, product units/inventory and language/locale lookups. Their public sources are still gaps, their dependencies are research hypotheses, and none is selected or executed. The four `demand-*` open questions retain geographic, temporal, amount/unit and identity/classification uncertainties for the fresh Phase 2 reconciler.

Phase 1 should prioritize sources that can explain country representations and contextual postcodes/localities across five demos, calendar reference reuse with event/validity distinctions, and explicit currency/unit/classification bridges. Currency, culture and unit datasets also have strong graph value despite narrower observed demo coverage. The same public reference may be reusable after reconciliation; coverage counts must never substitute for identifier quality, licence suitability or measured match rates.

Source/model/meaning licences are not interchangeable. Employees-derived observations retain its share-alike scope; mixed-source candidate metadata does not assert a blanket CC0 licence. Source selection and any redistribution remain subject to the programme licensing gate. Nothing here promotes the source fixtures or candidate groups into a competing canonical registry.

Review fixes r1/r2 retain [original hypotheses and reasons](evidence/demo-demand-corrections.json): phone/address types are classifications, the region-scope flag is not a name, customer person/store references are separate relationships, and keyless entities have no invented key. Seven identical surviving proposals after these source-role corrections were deduplicated. Final independent exact-head review remains pending.
