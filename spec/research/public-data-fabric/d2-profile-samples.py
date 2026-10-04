#!/usr/bin/env python3
"""Reprobe private bounded Phase1 captures; never fetch, publish, or accept mappings."""
import argparse, collections, csv, hashlib, io, json, pathlib, re, zipfile


def csv_rows(path):
    reader = csv.DictReader(io.StringIO(path.read_text(encoding='utf-8-sig')))
    return reader.fieldnames, list(reader)


def profiles(columns, rows):
    result = {}
    for column in columns:
        values = [row.get(column, '') for row in rows]
        result[column] = {
            'rows': len(values),
            'empty': sum(value == '' for value in values),
            'distinct_including_empty': len(set(values)),
            'outer_whitespace': sum(value != value.strip() for value in values),
            'leading_zero_digit_strings': sum(len(value) > 1 and value.startswith('0') and value.isdigit() for value in values),
        }
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--samples', required=True, type=pathlib.Path)
    parser.add_argument('--output', required=True, type=pathlib.Path)
    args = parser.parse_args()
    p = args.samples
    verified = {}
    for receipt_path in sorted(p.glob('*.receipt.json')):
        receipt = json.loads(receipt_path.read_text())
        raw = p / receipt_path.name.replace('.receipt.json', '.raw')
        if not raw.exists():
            continue
        digest = hashlib.sha256(raw.read_bytes()).hexdigest()
        if digest != receipt['sha256_read_bytes']:
            raise ValueError('Capture hash mismatch: ' + raw.name)
        verified[raw.name] = {'sha256_read_bytes': digest, 'observed_at': receipt['observed_at']}
    nc, n = csv_rows(p / 'naptan-010.raw')
    gc, g = csv_rows(p / 'nptg-full.raw')
    target = collections.Counter(row['NptgLocalityCode'] for row in g)
    z = zipfile.ZipFile(p / 'nta-feed.raw')
    feed = {}
    for name in z.namelist():
        reader = csv.DictReader(io.StringIO(z.read(name).decode('utf-8-sig')))
        rows = list(reader)
        feed[name] = {'columns': reader.fieldnames, 'rows': rows}
    checks = []
    for left, right, field in [
        ('routes.txt', 'agency.txt', 'agency_id'),
        ('trips.txt', 'routes.txt', 'route_id'),
        ('trips.txt', 'calendar.txt', 'service_id'),
        ('stop_times.txt', 'trips.txt', 'trip_id'),
        ('stop_times.txt', 'stops.txt', 'stop_id'),
    ]:
        values = {row[field] for row in feed[right]['rows']}
        checks.append({'from': left, 'to': right, 'field': field,
                       'checked_rows': len(feed[left]['rows']),
                       'unmatched': sum(row[field] not in values for row in feed[left]['rows'])})
    kev = json.loads((p / 'cisa.raw').read_bytes())
    nvd = json.loads((p / 'nvd.raw').read_bytes())
    cve = nvd['vulnerabilities'][0]['cve']['id']
    result = {
        'scope': 'Private cached captures only. A successful probe does not accept semantics, legal use, publication, or DataTug execution.',
        'verified_read_byte_hashes': verified,
        'naptan_area': '010',
        'naptan_rows': len(n),
        'nptg_rows': len(g),
        'nptg_duplicate_locality_ids': sum(count > 1 for count in target.values()),
        'stop_locality_unmatched': sum(row['NptgLocalityCode'] not in target for row in n),
        'naptan_profiles': profiles(nc, n),
        'nptg_profiles': profiles(gc, g),
        'nta_profiles': {name: profiles(data['columns'], data['rows']) for name, data in feed.items()},
        'nta_native_fk_checks': checks,
        'kev_rows': len(kev['vulnerabilities']),
        'kev_unique_cveIDs': len({row['cveID'] for row in kev['vulnerabilities']}),
        'kev_nvd_single_CVE': cve,
        'kev_matches_for_single_CVE': sum(row['cveID'] == cve for row in kev['vulnerabilities']),
        'police_polygon_rows': {name: len(json.loads((p / (name + '.raw')).read_bytes())) for name in ['police', 'police-small']},
        'fsa_complete_nodes_in_prefix': len(re.findall(r'<EstablishmentDetail>.*?</EstablishmentDetail>', (p / 'fsa-bulk.raw').read_text(), re.S)),
    }
    args.output.write_text(json.dumps(result, indent=2) + '\n')
    print(json.dumps({key: value for key, value in result.items() if key not in ['verified_read_byte_hashes', 'naptan_profiles', 'nptg_profiles', 'nta_profiles']}, indent=2))


if __name__ == '__main__':
    main()
