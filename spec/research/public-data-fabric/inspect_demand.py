#!/usr/bin/env python3
"""Reproduce bounded demo demand profiles; never normalize or mutate fixtures.

Inputs are research observations, not canonical mappings. Supply the six local
SQLite paths as a JSON object with --fixtures; every file must match the pinned
SHA256 in demand-plan.json. The output has no execution timestamp or local paths
so a reviewer can compare it byte for byte with evidence/demo-demand.json.
"""
import argparse
import hashlib
import json
from pathlib import Path
import re
import sqlite3

BASE = Path(__file__).resolve().parent


def digest(path):
    with path.open('rb') as stream:
        return hashlib.file_digest(stream, 'sha256').hexdigest()


def quoted(name):
    return '"' + name.replace('"', '""') + '"'


def scalar(value):
    if isinstance(value, bytes):
        return {'blob_bytes': len(value), 'sha256': hashlib.sha256(value).hexdigest()}
    if isinstance(value, str) and len(value) > 160:
        return {'text_prefix': value[:160], 'characters': len(value),
                'sha256_utf8': hashlib.sha256(value.encode()).hexdigest()}
    return value


def profile(conn, table, field, enumerate_values=False):
    t, f = quoted(table), quoted(field)
    counts = conn.execute(f'SELECT count(*), count({f}), count(DISTINCT {f} COLLATE BINARY), '
                          f'sum(CASE WHEN typeof({f})=\'text\' AND {f}=\'\' THEN 1 ELSE 0 END), '
                          f'sum(CASE WHEN typeof({f})=\'text\' AND trim({f})!={f} THEN 1 ELSE 0 END) '
                          f'FROM {t}').fetchone()
    types = dict(conn.execute(f'SELECT typeof({f}),count(*) FROM {t} GROUP BY typeof({f})'))
    samples = [scalar(r[0]) for r in conn.execute(
        f'SELECT DISTINCT {f} FROM {t} WHERE {f} IS NOT NULL ORDER BY {f} COLLATE BINARY LIMIT 8')]
    # Python regex avoids SQLite GLOB's wildcard semantics. Count exact text
    # numeric-leading-zero values without rewriting either values or identity.
    leading = []
    leading_count = 0
    for value, n in conn.execute(f'SELECT {f},count(*) FROM {t} WHERE typeof({f})=\'text\' GROUP BY {f}'):
        if re.fullmatch(r'0[0-9]+', value):
            leading_count += n
            leading.append(value)
    bounds = conn.execute(f'SELECT min({f}),max({f}) FROM {t}').fetchone()
    # SQLite lower() only folds ASCII. Case collisions are diagnostic evidence,
    # never a case-insensitive compatibility or Unicode normalization claim.
    collisions = [r[0] for r in conn.execute(
        f'SELECT lower({f}) FROM {t} WHERE typeof({f})=\'text\' '
        f'GROUP BY lower({f}) HAVING count(DISTINCT {f} COLLATE BINARY)>1 ORDER BY lower({f}) LIMIT 8')]
    sensitive = bool(re.search(r'password|salt|cardnumber|nationalid|loginid|username', field, re.I))
    result = {'rows': counts[0], 'nulls': counts[0]-counts[1],
            'distinct_nonnull': counts[2], 'empty_text': counts[3] or 0,
            'outer_ascii_spaces': counts[4] or 0, 'sqlite_storage_types': types,
            'binary_order_samples': [] if sensitive else samples,
            'values_redacted': sensitive, 'numeric_leading_zero_rows': leading_count,
            'numeric_leading_zero_examples': [] if sensitive else sorted(leading)[:8],
            'raw_min': None if sensitive else scalar(bounds[0]),
            'raw_max': None if sensitive else scalar(bounds[1]),
            'ascii_case_collisions': collisions}
    if enumerate_values and counts[2] <= 256 and not sensitive:
        result['all_distinct_values_with_counts'] = [
            [scalar(r[0]), r[1]] for r in conn.execute(
                f'SELECT {f},count(*) FROM {t} WHERE {f} IS NOT NULL '
                f'GROUP BY {f} COLLATE BINARY ORDER BY {f} COLLATE BINARY')]
    return result


def inspect(plan, fixtures):
    result = {'notice': 'Pinned native fixture observations and proposed demand only. '
              'No semantic reconciliation, compatibility acceptance, public publication or live query claim.',
              'catalogue': plan['catalogue'], 'observed_at': plan['observed_at'], 'corpora': {}}
    for demo, source in sorted(plan['inputs'].items()):
        path = Path(fixtures[demo]).resolve()
        actual = digest(path)
        if actual != source['sqlite_sha256']:
            raise ValueError(f'{demo}: fixture digest {actual} differs from pinned input')
        conn = sqlite3.connect(path.as_uri() + '?mode=ro&immutable=1', uri=True)
        conn.execute('PRAGMA query_only=ON')
        tables = {}
        for table, kind in conn.execute("SELECT name,type FROM sqlite_master WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY name"):
            columns = [{'name': r[1], 'declared_type': r[2], 'not_null': bool(r[3]),
                        'default_sql': r[4], 'primary_key_position': r[5]}
                       for r in conn.execute('PRAGMA table_info(' + quoted(table) + ')')]
            fks = [{'constraint_id': r[0], 'position': r[1], 'target_table': r[2],
                    'from_field': r[3], 'target_field': r[4], 'on_update': r[5], 'on_delete': r[6]}
                   for r in conn.execute('PRAGMA foreign_key_list(' + quoted(table) + ')')]
            indexes = [{'name': r[1], 'unique': bool(r[2]), 'origin': r[3],
                        'fields': [i[2] for i in conn.execute('PRAGMA index_info(' + quoted(r[1]) + ')')]}
                       for r in conn.execute('PRAGMA index_list(' + quoted(table) + ')')]
            tables[table] = {'object_kind': kind,
                             'rows': conn.execute('SELECT count(*) FROM ' + quoted(table)).fetchone()[0],
                             'columns': columns, 'foreign_keys': fks, 'indexes': indexes, 'profiles': {}}
        for mapping in plan['mappings']:
            if mapping['demo'] != demo:
                continue
            table = mapping['table']
            if table not in tables:
                raise ValueError(f'{demo}: missing native recordset {table}')
            names = {c['name'] for c in tables[table]['columns']}
            for field in mapping['fields']:
                if field not in names:
                    raise ValueError(f'{demo}.{table}: missing native field {field}')
                if field not in tables[table]['profiles']:
                    detailed = mapping['meaning_slug'] in (
                        'country', 'currency', 'currency-code', 'language', 'language-name',
                        'culture-locale', 'subnational-region', 'geographic-region', 'region',
                        'sales-region', 'unit-of-measure', 'payment-terms')
                    tables[table]['profiles'][field] = profile(conn, table, field, detailed)
        views = [r[0] for r in conn.execute("SELECT name FROM sqlite_master WHERE type='view' ORDER BY name")]
        result['corpora'][demo] = {'input': source,
                                 'table_count': sum(t['object_kind']=='table' for t in tables.values()),
                                 'view_count': len(views),
                                 'table_rows': sum(t['rows'] for t in tables.values() if t['object_kind']=='table'),
                                 'view_names': views, 'recordsets': tables}
        conn.close()
    return result


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fixtures', type=Path, required=True,
                        help='JSON object demo ID -> local SQLite file path')
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    plan = json.loads((BASE / 'evidence/demand-plan.json').read_text())
    result = inspect(plan, json.loads(args.fixtures.read_text()))
    args.output.write_text(json.dumps(result, indent=2, ensure_ascii=False) + '\n')
    print('Profiled ' + str(len(result['corpora'])) + ' pinned native fixtures')
