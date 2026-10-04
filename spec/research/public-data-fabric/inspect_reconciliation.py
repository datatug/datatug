#!/usr/bin/env python3
"""Read-only Phase 2 checks against hash-pinned native demo fixtures.

This measures local references and current core-country matching only. It does
not decide semantic equivalence, infer global identifiers, or publish data.
"""
import argparse
import hashlib
import json
from pathlib import Path
import sqlite3

BASE = Path(__file__).resolve().parent


def q(value):
    return '"' + value.replace('"', '""') + '"'


def country_counts(values, matched, unmatched, ambiguous):
    """Retain every source row; empty strings remain nonnull unmatched values."""
    total = sum(count for _, count in values)
    nulls = sum(count for value, count in values if value is None)
    missing = sum(row['rows'] for row in unmatched)
    uncertain = sum(row['rows'] for row in ambiguous)
    assert matched + missing + uncertain == total - nulls
    return dict(source_rows=total, nonnull_source_rows=total-nulls,
                null_source_rows=nulls,
                empty_string_rows=sum(count for value, count in values if value == ''),
                unmatched_nonnull_rows=missing, ambiguous_nonnull_rows=uncertain,
                unlinked_source_rows=nulls+missing+uncertain)


def inspect(fixtures, countries):
    demand = json.loads((BASE / 'evidence/demo-demand.json').read_text())['corpora']
    candidates = {p.stem: json.loads(p.read_text()) for p in (BASE / 'candidates/$records').glob('*.json')}
    assertions = {p.stem: json.loads(p.read_text()) for p in (BASE / 'assertions/$records').glob('demand-*.json')}
    conns = {}
    for demo, info in demand.items():
        path = Path(fixtures[demo]).resolve()
        with path.open('rb') as stream:
            assert hashlib.file_digest(stream, 'sha256').hexdigest() == info['input']['sqlite_sha256'], demo
        conns[demo] = sqlite3.connect(path.as_uri() + '?mode=ro&immutable=1', uri=True)
        conns[demo].execute('PRAGMA query_only=ON')
    result = {'scope': 'Native pinned sources; local-reference quality and starter core-country coverage only. No external source or semantic compatibility acceptance.',
              'input_hashes': {d: r['input']['sqlite_sha256'] for d, r in demand.items()},
              'local_references': {}, 'country_representations': {}, 'alternate_keys': {}}
    for key, row in sorted(assertions.items()):
        source = candidates[row['from_id']]
        demo = source['dataset_id']
        table = source['title'][len(demo) + 1:]
        conn = conns[demo]
        if row['kind'] == 'edge':
            target = candidates[row['to_id']]
            assert target['dataset_id'] == demo, key
            target_table = target['title'][len(demo) + 1:]
            fields, targets = row['from_fields'], row['to_fields']
            assert len(fields) == len(targets), key
            nonnull = ' AND '.join('s.' + q(f) + ' IS NOT NULL' for f in fields)
            condition = ' AND '.join('s.' + q(f) + '=t.' + q(t) for f, t in zip(fields, targets))
            count = conn.execute(f'SELECT count(*) FROM {q(table)} s WHERE {nonnull}').fetchone()[0]
            missing = conn.execute(f'SELECT count(*) FROM {q(table)} s WHERE {nonnull} AND NOT EXISTS '
                                   f'(SELECT 1 FROM {q(target_table)} t WHERE {condition})').fetchone()[0]
            keys = ','.join(q(f) for f in targets)
            target_nonnull = ' AND '.join(q(f) + ' IS NOT NULL' for f in targets)
            duplicate_groups = conn.execute(f'SELECT count(*) FROM (SELECT {keys} FROM {q(target_table)} '
                                             f'WHERE {target_nonnull} GROUP BY {keys} HAVING count(*)>1)').fetchone()[0]
            total = conn.execute(f'SELECT count(*) FROM {q(table)}').fetchone()[0]
            result['local_references'][key] = dict(demo=demo, source_table=table, target_table=target_table,
                from_fields=fields, to_fields=targets, source_rows=total, nonnull_source_rows=count,
                null_source_rows=total-count, unmatched_nonnull_rows=missing, target_duplicate_key_groups=duplicate_groups)
        elif row['to_id'] == 'demand-meaning-country':
            field = row['from_fields'][0]
            if field == 'country_id':
                continue  # A Sakila local key is not an ISO numeric code.
            values = list(conn.execute(f'SELECT {q(field)},count(*) FROM {q(table)} GROUP BY {q(field)}'))
            mode = 'codes.alpha2' if field == 'CountryRegionCode' else 'labels'
            matched, unmatched, ambiguous = 0, [], []
            for value, count in values:
                if value is None:
                    continue
                ids = []
                for country in countries:
                    labels = list(country['labels'].values())
                    labels += [v for words in country.get('aliases', {}).values() for v in words]
                    match = (str(value) == country['codes']['alpha2'] if mode == 'codes.alpha2'
                             else str(value).lower() in [v.lower() for v in labels])
                    if match:
                        ids.append(country['id'])
                if len(ids) == 1:
                    matched += count
                elif ids:
                    ambiguous.append({'value': value, 'rows': count, 'ids': ids})
                else:
                    unmatched.append({'value': value, 'rows': count})
            result['country_representations'][key] = dict(demo=demo, table=table, field=field,
                match=mode, matched_rows=matched, unmatched=unmatched, ambiguous=ambiguous,
                **country_counts(values, matched, unmatched, ambiguous))
    for demo, info in sorted(demand.items()):
        for table, rs in sorted(info['recordsets'].items()):
            unique = [i for i in rs['indexes'] if i['unique'] and i['origin'] != 'pk']
            if unique:
                result['alternate_keys'][demo + '.' + table] = unique
    for conn in conns.values():
        conn.close()
    return result


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fixtures', type=Path, required=True)
    parser.add_argument('--core-country-values', type=Path, required=True,
                        help='JSON values extracted from pinned core geo.meaning.yaml country concept')
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    result = inspect(json.loads(args.fixtures.read_text()), json.loads(args.core_country_values.read_text()))
    args.output.write_text(json.dumps(result, indent=2, ensure_ascii=False) + '\n')
    print(f"Measured {len(result['local_references'])} local references and {len(result['country_representations'])} country representations")
