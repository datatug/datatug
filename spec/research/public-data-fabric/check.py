#!/usr/bin/env python3
"""Research-only consistency constraints supplementing `ingitdb validate`."""
import argparse
import json
from pathlib import Path
import re
import shutil
import tempfile

BASE = Path(__file__).resolve().parent
TABLES = ('candidates', 'assertions', 'evidence', 'questions', 'implementation', 'gates')


def load(base):
    return {table: {p.stem: json.loads(p.read_text())
                    for p in sorted((base / table / '$records').glob('*.json'))}
            for table in TABLES}


def check(base=BASE):
    rows = load(base)
    errors = []

    def require(condition, label):
        if not condition:
            errors.append(label)

    for table in TABLES:
        schema = json.loads((base / table / '.collection/definition.yaml').read_text())['columns']
        for key, row in rows[table].items():
            label = f'{table}/{key}'
            require(bool(re.fullmatch(r'[a-z0-9][a-z0-9-]*', key)), label + ': invalid stable ID')
            require(not (set(row) - set(schema)), label + ': undeclared columns')
            # inGitDB owns scalar/array types, enums and lengths. These checks
            # also make dangling links fail before any status interpretation.
            for name, field in schema.items():
                require(not field.get('required') or name in row, label + ': missing ' + name)
                if 'foreign_key' in field and name in row:
                    target = field['foreign_key'].removeprefix('fabric.')
                    values = row[name] if field['type'] == '[]string' else [row[name]]
                    for value in values:
                        require(value in rows[target], label + ': dangling ' + name + '/' + value)
            if table != 'evidence' and row.get('status') in ('blocked', 'rejected'):
                require(bool(row.get('reason')), label + ': decision needs reason')
                require(bool(row.get('evidence_ids')), label + ': decision needs evidence')
            if table in ('candidates', 'assertions', 'implementation'):
                require(bool(row.get('evidence_ids')), label + ': lifecycle needs evidence')
            if row.get('status') == 'validated':
                require(any(rows['evidence'].get(e, {}).get('outcome') == 'pass'
                            for e in row.get('evidence_ids', [])), label + ': validated needs passing evidence')

    candidates = rows['candidates']

    def kind(id, allowed):
        return candidates.get(id, {}).get('kind') in allowed

    for key, row in candidates.items():
        if 'source_id' in row:
            require(kind(row['source_id'], ['source']), key + ': source_id must name source')
        if 'dataset_id' in row:
            require(kind(row['dataset_id'], ['dataset']), key + ': dataset_id must name dataset')
        if row['kind'] == 'recordset':
            require('dataset_id' in row, key + ': recordset needs parent dataset')
        if 'member_ids' in row:
            require(row['kind'] in ('scenario', 'bundle'), key + ': only scenarios/bundles have members')
        for id in row.get('demo_ids', []):
            require(kind(id, ['dataset']), key + ': demo_ids must name datasets')

    for key, row in rows['assertions'].items():
        if row.get('kind') == 'mapping':
            require(kind(row.get('from_id'), ['recordset']), key + ': mapping starts at recordset')
            require(kind(row.get('to_id'), ['meaning']), key + ': mapping targets meaning')
        if row.get('kind') == 'edge':
            require(kind(row.get('from_id'), ['recordset']) and kind(row.get('to_id'), ['recordset']),
                    key + ': join edge needs two recordsets')
        if 'meaning_id' in row:
            require(kind(row['meaning_id'], ['meaning']), key + ': meaning_id must name meaning')
        if row.get('compatibility') == 'transformable':
            require(bool(row.get('transformation')), key + ': transformable needs explicit transformation')
        if row.get('status') in ('reconciled', 'selected', 'implementing', 'implemented', 'validated'):
            require(row.get('compatibility') not in ('unresolved', 'incompatible'), key + ': unreconciled compatibility')
            require(any(rows['evidence'].get(e, {}).get('review_stage') == 'reconciliation'
                        and rows['evidence'].get(e, {}).get('independent') is True
                        and rows['evidence'].get(e, {}).get('outcome') == 'pass'
                        for e in row.get('evidence_ids', [])), key + ': needs dedicated reconciliation receipt')

    for key, row in rows['implementation'].items():
        if row.get('status') in ('implemented', 'validated'):
            require(bool(row.get('canonical_address')) and bool(row.get('repository'))
                    and bool(re.fullmatch('[0-9a-f]{40}', row.get('commit', ''))),
                    key + ': canonical publication needs address and exact repository pin')
            for name in ('publication_evidence_ids', 'query_evidence_ids'):
                require(bool(row.get(name)) and all(rows['evidence'].get(e, {}).get('outcome') == 'pass'
                        for e in row.get(name, [])), key + ': publication/query needs passing receipts')

    for key, row in rows['questions'].items():
        if row.get('status') in ('resolved', 'declined'):
            require(bool(row.get('resolution')) and bool(row.get('evidence_ids')), key + ': closure needs rationale/evidence')

    for key, row in rows['gates'].items():
        if row.get('status') == 'passed':
            require(all(rows['gates'].get(p, {}).get('status') == 'passed'
                        for p in row.get('prerequisite_ids', [])), key + ': prerequisite not passed')
            require(all(rows['questions'].get(q, {}).get('status') in ('resolved', 'declined')
                        for q in row.get('blocker_ids', [])), key + ': unresolved blocker')
            evidence = [rows['evidence'].get(e, {}) for e in row.get('evidence_ids', [])]
            require(bool(evidence) and all(e.get('outcome') == 'pass' for e in evidence), key + ': passing receipts required')
            require(any(e.get('independent') is True and e.get('reviewer')
                        for e in evidence), key + ': independent reviewer receipt required')

    def acyclic(table, field):
        visiting, visited = set(), set()

        def visit(key):
            if key in visiting:
                errors.append(table + ': cyclic ' + field)
                return
            if key in visited:
                return
            visiting.add(key)
            for dep in rows[table].get(key, {}).get(field, []):
                visit(dep)
            visiting.remove(key)
            visited.add(key)
        for key in rows[table]:
            visit(key)

    acyclic('gates', 'prerequisite_ids')
    acyclic('candidates', 'dependency_ids')
    snapshot = json.loads((base / 'evidence/corpus-models-and-bindings.json').read_text())
    pins = json.loads((base / 'evidence/snapshot-pins.json').read_text())
    for corpus in snapshot['corpora']:
        key = corpus['id']
        require(kind(key, ['dataset']), 'snapshot fixture missing candidate ' + key)
        require(candidates.get(key, {}).get('commit') == pins[key]['commit'], key + ': snapshot pin differs')
    return errors


def self_test():
    """Mutation checks for the guards that collection definitions cannot express."""
    mutations = {
        'unresolved gate': lambda r: r['gates']['phase0'].update(status='passed'),
        'phase ordering': lambda r: r['gates']['phase1'].update(status='passed'),
        'missing rejection reason': lambda r: r['candidates']['tpch'].pop('reason'),
        'dangling evidence': lambda r: r['candidates']['chinook'].update(evidence_ids=['absent']),
        'wrong link kind': lambda r: r['candidates']['chinook'].update(source_id='northwind'),
        'false validation': lambda r: r['candidates']['tpch'].update(status='validated'),
        'cyclic gate': lambda r: r['gates']['phase0'].update(prerequisite_ids=['phase1']),
        'lost fixture': lambda r: r['candidates'].pop('employees'),
        'unsubstantiated publication': lambda r: r['implementation'].update(fake={
            'title': 'fake', 'candidate_id': 'chinook', 'status': 'implemented',
            'owner': 'test', 'evidence_ids': ['corpus-chinook']}),
    }
    assert not check(), check()
    for label, mutate in mutations.items():
        with tempfile.TemporaryDirectory(prefix='fabric-ledger-test-') as tmp:
            base = Path(tmp) / 'fabric'
            shutil.copytree(BASE, base)
            rows = load(base)
            mutate(rows)
            for table in TABLES:
                for p in (base / table / '$records').glob('*.json'):
                    p.unlink()
                for key, row in rows[table].items():
                    (base / table / '$records' / (key + '.json')).write_text(json.dumps(row))
            assert check(base), 'accepted invalid mutation: ' + label
    print(f'{len(mutations)} ledger mutation checks passed')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--self-test', action='store_true')
    args = parser.parse_args()
    if args.self_test:
        self_test()
    errors = check()
    if errors:
        print('\n'.join(errors))
        raise SystemExit(1)
    print('Public Data Fabric ledger relationships, decisions and phase gates valid')
