#!/usr/bin/env python3
"""Reproduce bounded reconciliation projections; never fetch or implement sources.

The explicit decisions below are human-readable research judgements. Computed
overlap is evidence for those decisions, never a rule accepting future values.
Run from any directory; --output-dir defaults to a private caller-selected path.
"""
import argparse
from collections import Counter
import csv
import hashlib
import json
from pathlib import Path
import sqlite3

BASE = Path(__file__).resolve().parent
RECEIPT = 'f3r-dedicated-reconciliation'
INPUT = '2d355dabf8ef14d80e366f74f8646b9159cd7bf0'
OBSERVED = '2026-10-05'
LIMITS = ('Research acceptance only; use the enumerated tuples at pinned revisions. '
          'Preserve raw source identity and NULL, empty, unmatched and ambiguous rows. '
          'Fail closed on duplicate targets or changed tuples; no fuzzy match, current-ISO '
          'certification, sovereignty or territorial-boundary equivalence is implied. '
          'Refresh requires review. Canonical publication, source rights and DataTug '
          'execution remain separately gated; B1 stays open at phases 4/5/6.')


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def unique(rows, key):
    index = {}
    for row in rows:
        value = row[key]
        if value is None or value == '' or value in index:
            raise ValueError('NULL/empty/duplicate target key: ' + str(value))
        index[value] = row
    return index


def account(rows, key, targets):
    counts = dict(total=len(rows), null=0, empty=0, matched=0, unmatched=0, ambiguous=0)
    exceptions = []
    for row in rows:
        value = row[key]
        state = 'null' if value is None else 'empty' if value == '' else (
            'matched' if len(targets.get(value, [])) == 1 else
            'ambiguous' if len(targets.get(value, [])) > 1 else 'unmatched')
        counts[state] += 1
        if state != 'matched':
            exceptions.append({'state': state, 'row': row})
    return {'counts': counts, 'exceptions': exceptions}


def edge(id, title, source, target, fields, target_fields, transformation, reason, cardinality,
         evidence, source_namespace, target_namespace):
    return id, dict(title=title, kind='edge', status='reconciled', from_id=source,
                    to_id=target, from_fields=fields, to_fields=target_fields,
                    compatibility='transformable', confidence='high',
                    assumptions=LIMITS, transformation=transformation, reason=reason,
                    expected_cardinality=cardinality, physical_role='explicit scoped reference/context lookup',
                    from_identifier_namespace=source_namespace, to_identifier_namespace=target_namespace,
                    evidence_ids=evidence + [RECEIPT])


def run(args):
    source_path = BASE / 'evidence/f3-gap-probes.json'
    assert digest(source_path) == '4954ab39ad016fe6ba9dbb5c3b212ba8db8ed939145a3ab66e9e89973326a70d'
    source = json.loads(source_path.read_text())
    p2_path = BASE / 'evidence/p2-bridge-probes.json'
    assert digest(p2_path) == '392a7d310113223db974f004e7791e0cab71fb8601d805691b1978d800708e5f'
    p2 = json.loads(p2_path.read_text())
    obs = source['observations']
    geo = obs['geo_wb']['geonames_projection']
    geo2, geo3 = unique(geo, 'ISO'), unique(geo, 'ISO3')
    wb_all = obs['geo_wb']['worldbank_projection']
    wb_rows = [r for r in wb_all if r['region.id'] != 'NA']
    wb = unique(wb_rows, 'id')
    economies = [{'GeoNames': r, 'WorldBank': wb[r['ISO3']]} for r in geo if r['ISO3'] in wb]
    assert len(economies) == 216 and len(wb_rows) == 217
    assert wb['NAM']['iso2Code'] == 'NA' and wb['CHI']['iso2Code'] == 'JG'
    areas = [r for r in obs['ember']['area_projection'] if r['Area type'] == 'Country or economy' and r['ISO 3 code']]
    unique(areas, 'ISO 3 code')
    ember = [{'Ember': r, 'GeoNames': geo3[r['ISO 3 code']]} for r in areas]
    assert len(ember) == 209
    aw = obs['adventureworks']['projection']
    regular = [r for r in aw if r['cldr_validity_status'] == 'regular']
    historical = [r for r in aw if r['cldr_validity_status'] == 'deprecated']
    assert len(regular) == 236 and {r['CountryRegionCode'] for r in historical} == {'AN', 'CS'}
    for r in regular:
        assert r['CountryRegionCode'] == r['cldr_territory_codes']['type'] == r['geonames_codes']['ISO']
        assert r['cldr_territory_codes']['alpha3'] == r['geonames_codes']['ISO3']
        assert r['cldr_territory_codes']['numeric'] == r['geonames_codes']['ISO-Numeric']

    alias_reason = {
        'CD': 'The full Democratic Republic qualifier identifies CD rather than neighbouring CG. Reordered phrase and final article are accepted for this exact raw Sakila label only; an exact GeoNames alternate occurrence is not claimed.',
        'VA': 'Explicit parenthetical Vatican City State and exact GeoNames alternate support the territorial lookup VA. This does not equate the Holy See as a legal/religious authority with Vatican territory.',
        'KZ': 'Exact GeoNames alternate and concordant KZ/KAZ/398 identify Kazakhstan. No English-language tagging or general spelling correction is inferred.',
        'RU': 'Exact English alternate and UN label support RU/RUS/643. This is reference context, not a boundary definition for disputed territory or historic Russian entities.',
        'RE': 'Exact accented label and RE/REU/638 support the separate territory. Preserve the accent and do not replace Reunion with France.',
        'VI': 'U.S. qualifier and reordered Virgin Islands label support VI/VIR/850. No exact alternate occurrence is claimed; British Virgin Islands VG is explicitly excluded.'}
    additions = [{'native': r['native'], 'target': r['candidate_target_only'],
                  'decision': 'accepted exact source-row label representation',
                  'reason': alias_reason[r['candidate_target_only']['ISO']]}
                 for r in obs['sakila']['six_followups']]
    labels = [{'raw': r['raw'], 'ISO': r['ISO']} for r in p2['countries']['d1-edge-sakila-country']['reviewed_values']]
    labels += [{'raw': r['native']['country'], 'ISO': r['target']['ISO']} for r in additions]
    labels.sort(key=lambda r: r['raw'])
    unique(labels, 'raw')
    assert len(labels) == 108
    assert not any(r['raw'] == 'Yugoslavia' for r in labels)
    native_path = args.sakila
    assert digest(native_path) == source['source_contracts']['native']['sakila']['sha256']
    con = sqlite3.connect(f'file:{native_path}?mode=ro', uri=True)
    con.row_factory = sqlite3.Row
    native = [dict(r) for r in con.execute('SELECT country_id,country,last_update FROM country ORDER BY country_id')]
    con.close()
    label_targets = {r['raw']: [geo2[r['ISO']]] for r in labels}
    sakila = account(native, 'country', label_targets)
    assert sakila['counts'] == dict(total=109, null=0, empty=0, matched=108, unmatched=1, ambiguous=0)

    annual_path = args.ember_annual
    assert digest(annual_path) == source['source_receipts']['ember-annual']['sha256']
    with annual_path.open() as f:
        facts = list(csv.DictReader(f))
    totals = [r for r in facts if r['Area type'] == 'Country or economy' and r['ISO 3 code']
              and r['Electricity source'] == 'Total generation']
    keys = Counter((r['ISO 3 code'], r['Year']) for r in totals)
    assert max(keys.values()) == 1 and len(totals) == 5915
    assert {r['Is aggregated source'] for r in totals} == {'True'}
    assert all(r['Generation (TWh)'] != '' for r in totals)
    annual = dict(source_rows=len(facts), projected_rows=len(totals), excluded_rows=len(facts)-len(totals),
                  missing_generation_cells=sum(r['Generation (TWh)'] == '' for r in totals),
                  duplicate_key_groups=sum(n > 1 for n in keys.values()),
                  key=['ISO 3 code', 'Year'], min_year=min(r['Year'] for r in totals),
                  max_year=max(r['Year'] for r in totals),
                  per_area=[{'ISO3': code, 'years': sorted(int(y) for c,y in keys if c == code)}
                            for code in sorted({c for c,y in keys})])

    # Synthetic counterexamples make the denominator and refusal contract observable.
    negative = account([{'k': None}, {'k': ''}, {'k': 'na'}, {'k': 'NAM'}, {'k': 'X'}],
                       'k', {'NAM': [1], 'X': [1, 2]})
    assert negative['counts'] == dict(total=5, null=1, empty=1, matched=1, unmatched=1, ambiguous=1)
    duplicate_refused = False
    try:
        unique([{'k': 'A'}, {'k': 'A'}], 'k')
    except ValueError:
        duplicate_refused = True
    assert duplicate_refused
    assert 'Yugoslavia' not in label_targets and 'Virgin Islands' not in label_targets
    assert 'Réunion ' not in label_targets and 'Congo' not in label_targets
    assert 'AN' not in {r['CountryRegionCode'] for r in regular}
    assert 'Net imports' != 'Total generation' and 'Demand' != 'Total generation'

    assertions = dict([
        edge('f3r-edge-geo-wb-context', 'Reviewed GeoNames areas to World Bank nonaggregate economy context',
             'd1-geo-countries', 'd1-wb-economies', ['ISO3'], ['id'],
             'Use only the 216 GeoNames/WorldBank tuples in f3r-decisions.json.geo_wb.accepted_tuples. Project $[1][*] WHERE region.id != "NA"; exact ISO3=id lookup after tuple and unique-key validation. LEFT lookup, retain all 252 source rows; 36 unmatched. Never filter economy.iso2Code != "NA".',
             'Accepted provider-area to reporting-economy reference context after reviewing both labels, code namespaces and source categories. It is not unconditional ISO identity or identical statistical/geographic boundaries. CHI/JG is not Jersey or Guernsey alone; absent territories remain unmatched. Broad d1-mapping-wb-economy stays rejected; d1-edge-geo-wb-economy remains the unaccepted broad proposal.',
             '252 GeoNames rows -> 216 single targets + 36 unmatched; reverse 217 economies -> 216 targets + CHI unmatched; no NULL/empty/ambiguous/duplicate keys at the pinned snapshot',
             ['f3-geo-economy-source-probes'], 'GeoNames:countryInfo pinned ISO3', 'WorldBank:v2 pinned economy.id'),
        edge('f3r-edge-ember-geo-area', 'Reviewed annual Ember reporting areas to GeoNames reference context',
             'd2-rs-ember-year', 'd1-geo-countries', ['ISO 3 code'], ['ISO3'],
             'Filter Area type == "Country or economy" AND ISO 3 code is nonempty. Use only 209 reviewed tuples in f3r-decisions.json.ember.accepted_tuples. DISTINCT (Area, ISO 3 code, Area type) forms the area projection; exact tuple/code lookup to GeoNames, never fact deduplication. Preserve Year and Electricity source on facts.',
             'Accepted annual source-area context after full label/code review. XKX remains a provider extension, not an official ISO certificate. Fifteen Region areas are excluded with their raw rows retained. Historical observation geography follows Ember methodology, not a timeless GeoNames boundary. No monthly acceptance follows.',
             '209 projected areas -> 209 unique GeoNames targets; 97,937 annual facts -> many:one reference lookup. Reverse GeoNames 252 -> 209 matched + 43 unmatched; 200 of209 areas reach WB and9 remain unmatched',
             ['f3-ember-area-source-probes','f3-geo-economy-source-probes'], 'Ember:annual country/economy ISO3', 'GeoNames:countryInfo pinned ISO3'),
        edge('f3r-edge-sakila-reviewed-country', 'Pinned Sakila 108-label partial country reference bridge',
             'demand-sakila-country', 'd1-geo-countries', ['country'], ['ISO'],
             'Use only exact raw label -> ISO pairs in f3r-decisions.json.sakila.accepted_labels at the native SHA. Preserve country_id and country. No runtime case folding, trim, accent stripping or fuzzy expansion. LEFT lookup to unique GeoNames ISO, output explicit unmatched/ambiguous states. Yugoslavia remains unmatched.',
             'Accept the 102 previously measured labels after independent review plus six individually reasoned aliases. This partial bridge is separate from the original whole-field proposal. Sakila IDs are local; last_update is operational metadata and supplies no historical validity. Territorial lookup is neither a legal-authority identity nor proof of a dated boundary.',
             '109 source rows ->108 single targets +1 historical unresolved; NULL0 empty0 ambiguous0; target ISO unique',
             ['p2-bridge-probes','f3-demo-country-representations'], 'Sakila:pinned country.country raw label', 'GeoNames:countryInfo pinned ISO'),
        edge('f3r-edge-aw-cldr-regular', 'AdventureWorks regular-code subset to CLDR48 encoding rows',
             'demand-adventureworks-person-countryregion','d1-cldr-codes',['CountryRegionCode'],['type'],
             'Use only 236 exact (CountryRegionCode,Name) rows in f3r-decisions.json.adventureworks.regular_tuples whose CLDR48 validity is regular; lookup codeMappings/territoryCodes@type and retain alpha3/numeric strings. Leave AN and CS in the source with historical-unresolved disposition.',
             'Accept source-code representation lookup into the pinned CLDR encoding table after code/name/category and alpha3/numeric corroboration. CLDR regular is not certification of current ISO assignment or sovereignty. Do not apply territoryAlias replacements; broad 238-row proposal remains unaccepted.',
             '238 input rows ->236 accepted single targets +2 historically excluded; NULL0 empty0 ambiguous0; raw encoding-only probe238/238 is not accepted coverage',
             ['f3-demo-country-representations'], 'AdventureWorks:pinned Person.CountryRegion.CountryRegionCode','CLDR48:codeMappings/territoryCodes@type'),
        edge('f3r-edge-aw-geo-regular', 'AdventureWorks reviewed 236-code country reference context',
             'demand-adventureworks-person-countryregion','d1-geo-countries',['CountryRegionCode'],['ISO'],
             'Use the same 236 reviewed (CountryRegionCode,Name) tuples as f3r-edge-aw-cldr-regular, with exact GeoNames ISO and corroborating CLDR alpha3/numeric. LEFT lookup; AN/CS explicitly historical-unresolved; ModifiedDate does not supply validity.',
             'Accepted reference-area context only. Historic spellings and geopolitical labels are preserved. Exact administrative extent at any transaction date remains unproved; no boundary or successor substitution. Source labels distinguish CS Serbia and Montenegro from earlier reuse, but do not establish applicable observation periods.',
             '238 input rows ->236 accepted single targets +AN/CS excluded; NULL0 empty0 ambiguous0; GeoNames target ISO unique',
             ['f3-demo-country-representations'], 'AdventureWorks:pinned Person.CountryRegion.CountryRegionCode','GeoNames:countryInfo pinned ISO'),
    ])
    assertions['f3r-map-ember-total-generation'] = dict(
        title='Filtered annual total-generation TWh implements electricity generation amount', kind='mapping',
        status='reconciled', from_id='d2-rs-ember-year', to_id='d2-meaning-electricity',
        meaning_id='d2-meaning-electricity', from_fields=['Generation (TWh)'],
        to_fields=['electricity-generation-amount'], compatibility='transformable', confidence='high',
        assumptions=LIMITS, evidence_ids=['f3-ember-area-source-probes',RECEIPT],
        physical_role='source-defined annual quantity value, not a derived fleet sum',
        expected_cardinality='5,915 annual values at unique (ISO 3 code,Year) for selected Total generation series; 209 areas, years1985-2025 with actual per-area year sets retained',
        transformation='Filter Area type == "Country or economy", nonempty ISO 3 code, Electricity source == "Total generation", and validate Is aggregated source == "True". Project ISO3,Year,raw Generation (TWh), source revision. Preserve missing values as NULL, never zero. No aggregation across fuels, territories or periods is authorized.',
        reason='Accepted electrical energy produced over a stated reporting area/period, expressed in TWh under Ember source methodology. Demand, net imports, capacity GW, percentages, intensity and emissions are distinct. Selecting is_aggregate_series=false would discard this accepted total; nonaggregate Net imports also proves that flag is not a generation-semantic filter. Gross/net conventions and boundary/method revisions remain provider-qualified.')

    previous = json.loads((BASE/'candidates/$records/d2-meaning-electricity.json').read_text())
    if '\n\nF3R ACCEPTED:' in previous['reason']:
        previous['reason'] = previous['reason'].split('\n\nF3R ACCEPTED:')[0]
        previous['semantic_disposition'] = 'unresolved'
        previous['evidence_ids'] = [e for e in previous['evidence_ids'] if e != RECEIPT]
    meaning = dict(previous)
    meaning['semantic_disposition'] = 'accepted'
    meaning['evidence_ids'] = previous['evidence_ids'] + [RECEIPT]
    meaning['reason'] += ('\n\nF3R ACCEPTED: Narrow reusable definition: electrical energy produced over a specified reporting area and period, qualified by production source and methodology; unit TWh is energy, not power. Default aggregation none because source totals and fuel groups overlap. Only f3r-map-ember-total-generation has an accepted source binding here. Demand, net imports, capacity, emissions, intensity and shares remain distinct; no universal provider/method/boundary equality. Prefer canonical reuse after definition search; no canonical artifact is created by this research decision.')
    evidence = dict(scope='Fresh dedicated semantic reconciliation; derived evidence, not another registry, implementation, phase gate or independent adversarial review.',
                    hub_input=INPUT, observed_at=OBSERVED, source_evidence_sha256=digest(source_path),
                    p2_bridge_evidence_sha256=digest(p2_path),
                    prior_records_preserved='All pre-existing assertions/candidates unchanged except the explicitly narrowed electricity meaning; its exact prior state follows.',
                    previous_electricity_meaning=previous,
                    geo_wb={'accepted_tuples':economies,'forward':obs['geo_wb']['geo_to_nonaggregate_wb'],
                            'reverse':obs['geo_wb']['nonaggregate_wb_to_geo'],
                            'excluded_aggregate_ids':[r['id'] for r in wb_all if r['region.id']=='NA']},
                    ember={'accepted_tuples':ember,'area_forward':obs['ember']['ember_country_to_geo'],
                           'area_reverse':obs['ember']['geo_to_ember_country_projection'],
                           'wb_lookup':obs['ember']['ember_country_to_nonaggregate_wb'],
                           'excluded_regions':[r for r in obs['ember']['area_projection'] if r['Area type']=='Region'],
                           'annual_total_generation':annual,
                           'meaning_source':{'url':'https://ember-energy.org/data/yearly-electricity-data/',
                            'observed_at':OBSERVED,'access':'web provider; no byte digest claimed',
                            'bounded_facts':'Series include fuels, overlapping aggregates, demand and net imports; metrics are separate columns. Landing-page month wording is inconsistent with annual Year CSV; annual grain follows actual file and yearly endpoint.'}},
                    sakila={'accepted_labels':labels,'six_alias_decisions':additions,'measured':sakila,
                            'historical_unresolved':obs['sakila']['historical_separate'],
                            'rejected':'YU alias RS/ME is not an identity-preserving substitution; intended historic entity/date/boundary unresolved'},
                    adventureworks={'regular_tuples':regular,'historical_excluded':historical,
                        'accepted_counts':dict(total=238,null=0,empty=0,matched=236,unmatched=2,ambiguous=0),
                        'rejected':'No AN -> CW/SX/BQ or CS -> RS/ME substitution. CLDR encoding availability does not make deprecated rows current; ModifiedDate is not validity.'},
                    unresolved=['World Bank indicator-specific area/time/method alignment and rights; catalogue context does not authorize per-capita/comparative arithmetic.',
                        'Ember monthly universe and complete month/year comparison; d2-edge-ember-annual-monthly remains unresolved.',
                        'Chinook invoice-to-electricity execution and invoice denominator checks; d2-edge-invoice-electricity remains unresolved. Use annual area/year aggregate context only, never infer invoice energy, footprint or currency.',
                        'Sakila Yugoslavia historical identity; AW AN/CS historic interval applications; full current ISO/canonical country bindings.'],
                    counterexamples={'denominators':negative,'duplicate_target_refused':duplicate_refused,
                        'exclusions':['Namibia retained while WB region aggregate NA excluded','CHI not GG/JE','YU not RS/ME','AN/CS excluded','bare Virgin Islands and Congo unresolved','no accent/whitespace normalization','Total generation chosen explicitly; Demand/Net imports excluded']})
    out = args.output_dir
    def write(relative, value):
        path=out/relative;path.parent.mkdir(parents=True,exist_ok=True)
        path.write_text(json.dumps(value,ensure_ascii=False,indent=2)+'\n')
    write('evidence/f3r-decisions.json', evidence)
    write('evidence/$records/'+RECEIPT+'.json', dict(title='Fresh bounded Phase3 semantic follow-up',
        observed_at=OBSERVED,outcome='pass',summary='Dedicated fresh reconciler completed scoped decisions with explicit rejection and unresolved cases. Not adversarial review, a gate, implementation, source publication or DataTug acceptance.',
        reviewer='fresh-semantic-gap-followup', independent=True,review_stage='reconciliation',
        path='spec/research/public-data-fabric/f3r-reconciliation.md'))
    for id,row in assertions.items():write('assertions/$records/'+id+'.json',row)
    write('candidates/$records/d2-meaning-electricity.json',meaning)
    print(json.dumps({'accepted_assertions':len(assertions),'geo_wb':216,'ember_areas':209,
        'annual_total_values':len(totals),'sakila':sakila['counts'],'aw':236,'output':str(out)}))


if __name__ == '__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--sakila',type=Path,required=True)
    parser.add_argument('--ember-annual',type=Path,required=True)
    parser.add_argument('--output-dir',type=Path,required=True)
    run(parser.parse_args())
