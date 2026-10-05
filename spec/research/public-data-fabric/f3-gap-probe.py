#!/usr/bin/env python3
"""Bounded source observations only; no semantic acceptance or canonical publication."""
import argparse
import collections
import csv
import datetime
import hashlib
import html
import io
import json
from pathlib import Path
import re
import sqlite3
import time
import urllib.request
import xml.etree.ElementTree as ET

BASE = Path(__file__).resolve().parent
LABELS = {
    'Congo, The Democratic Republic of the': ('CD', 'geo-CD'),
    'Holy See (Vatican City State)': ('VA', 'geo-VA'),
    'Kazakstan': ('KZ', 'geo-kz'),
    'Russian Federation': ('RU', 'geo-RU'),
    'Réunion': ('RE', 'geo-RE'),
    'Virgin Islands, U.S.': ('VI', 'geo-VI'),
}


def digest(body):
    return hashlib.sha256(body).hexdigest()


def table_rows(body):
    """HTML whitespace repair is display extraction, never native-label repair."""
    return [[re.sub(r'\s+', ' ', html.unescape(re.sub('<[^>]+>', ' ', cell))).strip()
             for cell in re.findall(r'<td\b[^>]*>(.*?)</td>', row, re.S)]
            for row in re.findall(r'<tr\b[^>]*>.*?</tr>', body, re.S)]


def join_probe(left, right, left_field, right_field):
    index = collections.defaultdict(list)
    for row in right:
        key = row.get(right_field)
        if key not in (None, ''):
            index[key].append(row)
    count = collections.Counter()
    unmatched = []
    ambiguous = []
    for row in left:
        count['total'] += 1
        key = row.get(left_field)
        if key is None:
            count['null'] += 1
        elif key == '':
            count['empty'] += 1
        elif len(index.get(key, [])) == 1:
            count['matched'] += 1
        elif len(index.get(key, [])) > 1:
            count['ambiguous'] += 1
            ambiguous.append(row)
        else:
            count['unmatched'] += 1
            unmatched.append(row)
    return {
        'comparison': 'Exact raw string equality; source probe only, no identity approval',
        'from_field': left_field, 'to_field': right_field,
        'counts': {k: count[k] for k in ('total', 'null', 'empty', 'matched', 'unmatched', 'ambiguous')},
        'target_rows': len(right),
        'target_null_key_rows': sum(r.get(right_field) is None for r in right),
        'target_empty_key_rows': sum(r.get(right_field) == '' for r in right),
        'target_distinct_nonempty_keys': len(index),
        'target_duplicate_key_groups': sum(len(v) > 1 for v in index.values()),
        'max_target_multiplicity': max(map(len, index.values()), default=0),
        'unmatched': unmatched, 'ambiguous': ambiguous,
    }


def capture(directory, existing):
    """Refresh only the recorded URLs, privately, preserving size/time limits."""
    directory.mkdir(parents=True, exist_ok=True)
    receipts = {}
    for key, old in existing.items():
        receipt = {k: old[k] for k in ('url', 'max_bytes')}
        receipt['observed_at'] = datetime.datetime.now(datetime.timezone.utc).isoformat()
        start = time.monotonic()
        try:
            with urllib.request.urlopen(old['url'], timeout=25) as response:
                if response.url != old['url']:
                    raise ValueError('Unexpected redirect; source scope requires review')
                body = bytearray()
                while True:
                    part = response.read(65536)
                    if not part:
                        break
                    body.extend(part)
                    if len(body) > old['max_bytes']:
                        body = body[:old['max_bytes'] + 1]
                        break
                    if time.monotonic() - start > 90:
                        raise TimeoutError('90 second capture limit')
                receipt.update(status=response.status, final_url=response.url,
                               bytes=len(body), sha256=digest(body),
                               complete=len(body) <= old['max_bytes'],
                               headers={k: response.headers[k] for k in
                                        ('Last-Modified', 'ETag', 'Content-Length', 'Content-Type')
                                        if k in response.headers})
                if key == 'ember-annual' and not receipt['complete']:
                    raise ValueError('Annual CSV exceeded 17MB limit')
                (directory / key).write_bytes(body)
        except Exception as error:
            receipt['error'] = str(error)
        receipts[key] = receipt
    (directory / 'receipts.json').write_text(json.dumps(receipts, indent=2) + '\n')


def cldr_validity(path):
    result = {}
    for element in ET.parse(path).findall('./idValidity/id'):
        codes = []
        for token in element.text.split():
            if '~' not in token:
                codes.append(token)
                continue
            start, suffix = token.split('~')
            end = start[:-len(suffix)] + suffix
            assert start[:-1] == end[:-1], 'Unreviewed CLDR range syntax'
            codes.extend(start[:-1] + chr(c) for c in range(ord(start[-1]), ord(end[-1]) + 1))
        result[element.attrib['idStatus']] = sorted(codes)
    return result


def native(path, expected_hash, sql):
    assert digest(path.read_bytes()) == expected_hash, 'Native input hash differs: ' + str(path)
    with sqlite3.connect(path.as_uri() + '?mode=ro&immutable=1', uri=True) as conn:
        conn.row_factory = sqlite3.Row
        return [dict(row) for row in conn.execute(sql)]


def observe(args):
    p = args.captures
    geo = []
    for line in (p / 'geo').read_text().splitlines():
        if line and not line.startswith('#'):
            fields = line.split('\t')
            geo.append(dict(zip(('ISO', 'ISO3', 'ISO-Numeric', 'fips', 'Country'), fields[:5])) |
                       {'geonameid': fields[16]})
    wb_body = json.loads((p / 'wb').read_text())
    assert wb_body[0]['pages'] == 1 and wb_body[0]['total'] == len(wb_body[1])
    wb = [{k: row[k] for k in ('id', 'iso2Code', 'name', 'region', 'adminregion', 'incomeLevel', 'lendingType')}
          for row in wb_body[1]]
    wb_nonaggregate = [r for r in wb if r['region']['id'] != 'NA']
    annual = list(csv.DictReader(io.StringIO((p / 'ember-annual').read_text(encoding='utf-8-sig'))))
    area_keys = sorted({(r['Area'], r['ISO 3 code'], r['Area type']) for r in annual})
    area_projection = []
    for area, iso3, category in area_keys:
        group = [r for r in annual if r['Area'] == area]
        area_projection.append({'Area': area, 'ISO 3 code': iso3, 'Area type': category,
                                'rows': len(group), 'min_year': min(int(r['Year']) for r in group),
                                'max_year': max(int(r['Year']) for r in group)})
    countries = [a for a in area_projection if a['Area type'] == 'Country or economy']
    validity = cldr_validity(p / 'cldr-validity')
    statuses = {code: category for category, codes in validity.items() for code in codes}
    supp = ET.parse(p / 'cldr-supp')
    meta = ET.parse(p / 'cldr-meta')
    code_map = {r.attrib['type']: r.attrib for r in supp.findall('.//territoryCodes')}
    aliases = {r.attrib['type']: r.attrib for r in meta.findall('.//territoryAlias')}
    territory_info = {r.attrib['type'] for r in supp.findall('.//territoryInfo/territory')}
    cldr_labels = [(r.attrib, r.text) for r in ET.parse(p / 'cldr-en').findall('.//territories/territory')]
    table = re.search(r'<table id\s*=\s*"downloadTableEN".*?</table>', (p / 'un-m49').read_text(), re.S)
    assert table, 'Complete English UN table missing from prefix'
    un = table_rows(table.group())
    sakila = native(args.sakila, '9bebeee50fecb1fee115c206a4f380f2f6d1f009e7796b18b8a8cfba7c810454',
                    'SELECT country_id, country, last_update FROM country ORDER BY country_id')
    aw = native(args.adventureworks, 'e0f352f0e3a28ff15158130c237a91d58f1b60fffa7f30d2c412b92df4065dfd',
                'SELECT CountryRegionCode, Name, ModifiedDate FROM "Person.CountryRegion" ORDER BY CountryRegionCode')
    geo_iso = {r['ISO']: r for r in geo}
    six = []
    for raw, (iso, page) in LABELS.items():
        source = next(r for r in sakila if r['country'] == raw)
        rows = table_rows((p / page).read_text())
        exact = [r[:4] for r in rows if len(r) >= 4 and r[2] == raw]
        english = [r[:4] for r in rows if len(r) >= 4 and r[0] == 'en']
        six.append({'native': source, 'candidate_target_only': geo_iso[iso],
                    'geonames_page_source_key': page, 'exact_alternate_name_rows': exact,
                    'english_alternate_name_rows': english,
                    'un_code_record': [dict(zip(('Country or Area', 'M49 Code', 'ISO-alpha2 Code', 'ISO-alpha3 Code'), r[8:12]))
                                       for r in un if len(r) > 11 and r[10] == iso],
                    'cldr_english_labels': [{'attributes': a, 'text': v} for a, v in cldr_labels if a['type'] == iso],
                    'disposition': 'observation; dedicated reconciler must decide representation equivalence'})
    aw_projection = []
    for row in aw:
        code = row['CountryRegionCode']
        aw_projection.append(row | {'cldr_validity_status': statuses.get(code), 'cldr_territory_codes': code_map.get(code),
                                     'cldr_territory_info_present': code in territory_info,
                                     'cldr_alias': aliases.get(code),
                                     'geonames_codes': {k: geo_iso[code][k] for k in ('ISO', 'ISO3', 'ISO-Numeric', 'geonameid')} if code in geo_iso else None})
    return {
        'geo_wb': {'worldbank_response_metadata': wb_body[0], 'geonames_projection': geo,
                   'worldbank_projection': [
                       {k: r[k] for k in ('id', 'iso2Code', 'name')} |
                       {field + '.id': r[field]['id'] for field in ('region', 'adminregion', 'incomeLevel', 'lendingType')}
                       for r in wb],
                   'category_universes': {field: [dict(zip(('id', 'iso2code', 'value'), key)) | {'rows': count}
                                                 for key, count in sorted(collections.Counter(
                                                     tuple(r[field][f] for f in ('id', 'iso2code', 'value')) for r in wb).items())]
                                          for field in ('region', 'adminregion', 'incomeLevel', 'lendingType')},
                   'geo_to_all_wb': join_probe(geo, wb, 'ISO3', 'id'),
                   'geo_to_nonaggregate_wb': join_probe(geo, wb_nonaggregate, 'ISO3', 'id'),
                   'nonaggregate_wb_to_geo': join_probe(wb_nonaggregate, geo, 'id', 'ISO3')},
        'ember': {'header': list(annual[0]), 'source_rows': len(annual),
                  'country_fact_rows_per_iso3': {k: v for k, v in sorted(collections.Counter(r['ISO 3 code'] for r in annual if r['Area type'] == 'Country or economy').items())},
                  'area_projection': area_projection,
                  'row_categories': dict(collections.Counter(r['Area type'] for r in annual)),
                  'country_fact_to_geo': join_probe([{'ISO 3 code': r['ISO 3 code']} for r in annual if r['Area type'] == 'Country or economy'], geo, 'ISO 3 code', 'ISO3'),
                  'country_fact_composite_key': ['ISO 3 code', 'Year', 'Electricity source'],
                  'duplicate_country_iso3_year_source_keys': sum(v > 1 for v in collections.Counter((r['ISO 3 code'], r['Year'], r['Electricity source']) for r in annual if r['Area type'] == 'Country or economy').values()),
                  'area_categories': dict(collections.Counter(a['Area type'] for a in area_projection)),
                  'series_aggregate_universe': [dict(zip(('Electricity source', 'Is aggregated source'), v))
                                               for v in sorted({(r['Electricity source'], r['Is aggregated source']) for r in annual})],
                  'duplicate_area_year_source_keys': sum(v > 1 for v in collections.Counter(
                      (r['Area'], r['Year'], r['Electricity source']) for r in annual).values()),
                  'geo_to_ember_country_projection': join_probe(geo, countries, 'ISO3', 'ISO 3 code'),
                  'ember_country_to_geo': join_probe(countries, geo, 'ISO 3 code', 'ISO3'),
                  'ember_country_to_nonaggregate_wb': join_probe(countries, wb_nonaggregate, 'ISO 3 code', 'id'),
                  'monthly_scope': 'No complete monthly universe fetched; annual findings do not establish monthly coverage'},
        'sakila': {'table': 'country', 'native_key': ['country_id'], 'native_label': 'country',
                   'six_followups': six, 'historical_separate': next(r for r in sakila if r['country'] == 'Yugoslavia'),
                   'cldr_yu_alias_not_identity': aliases['YU'],
                   'exact_whole_label_to_complete_geonames_countryinfo': join_probe(sakila, geo, 'country', 'Country'),
                   'exact_raw_alternate_labels_in_six_assigned_pages': sum(bool(r['exact_alternate_name_rows']) for r in six)},
        'adventureworks': {'table': 'Person.CountryRegion', 'native_key': ['CountryRegionCode'],
                          'source_rows': len(aw), 'projection': aw_projection,
                          'cldr_validity_path': '/supplementalData/idValidity/id[@type=region][@idStatus]',
                          'cldr_validity_universe': validity,
                          'aw_category_counts': {category: sum(statuses.get(r['CountryRegionCode']) == category for r in aw) for category in validity},
                          'aw_missing_validity': [r['CountryRegionCode'] for r in aw if r['CountryRegionCode'] not in statuses],
                          'category_caution': 'CLDR regular is not current ISO-assigned: it also contains provider extensions such as XK, AC, DG, EA, IC, TA. Do not conflate regular with sovereignty.',
                          'modified_date_distinct': sorted({r['ModifiedDate'] for r in aw}),
                          'aw_to_geonames': join_probe(aw, geo, 'CountryRegionCode', 'ISO'),
                          'aw_to_cldr_codes': join_probe(aw, list(code_map.values()), 'CountryRegionCode', 'type'),
                          'cldr_deprecated': [r for r in aw_projection if r['cldr_alias']],
                          'missing_cldr_territory_info': [r['CountryRegionCode'] for r in aw_projection if not r['cldr_territory_info_present']]},
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--captures', type=Path, required=True)
    parser.add_argument('--sakila', type=Path, required=True)
    parser.add_argument('--adventureworks', type=Path, required=True)
    parser.add_argument('--capture', action='store_true')
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    pinned = json.loads((BASE / 'evidence/f3-gap-probes.json').read_text())
    if args.capture:
        capture(args.captures, pinned['source_receipts'])
    for key, receipt in pinned['source_receipts'].items():
        if receipt.get('sha256'):
            assert digest((args.captures / key).read_bytes()) == receipt['sha256'], 'Source revision differs: ' + key
    result = observe(args)
    args.output.write_text(json.dumps(result, indent=2, ensure_ascii=False) + '\n')
    print(json.dumps({'reproduced_sections': list(result), 'output': str(args.output)}))


if __name__ == '__main__':
    main()
