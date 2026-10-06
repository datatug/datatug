#!/usr/bin/env python3
"""Verify the superseding metadata decision and retained historical evidence offline."""
import argparse
import hashlib
import json
import subprocess
import tempfile
from pathlib import Path

BASE = Path(__file__).resolve().parent
ROOT = BASE.parents[2]
PORTFOLIO = 'spec/research/public-data-fabric/evidence/p3-portfolio.json'
RECEIPT = BASE / 'evidence/no-copy-access-decision.json'
SUPERSEDED_FILES = frozenset(['spec/research/public-data-fabric/p3-portfolio.md', 'spec/research/public-data-fabric/evidence/p3-portfolio.json', 'spec/research/public-data-fabric/evidence/w1-followup-decisions/verify.py', 'spec/research/public-data-fabric/evidence/w1-transport-app-fx/verify.py', '.github/workflows/public-data-fabric.yml', 'spec/research/public-data-fabric/README.md', 'spec/features/public-data-fabric/README.md'])
APPROVED_EVIDENCE_ADDITIONS = {
    'bigquery-catalogue-corrections-20261006': '4d7da4bf0c36c71ee2caba31d7876fe833407920af55cfda88bd3df4571d7ca3',
    'bigquery-catalogue-discovery-20261006': '767265e46e8c234afaa9048460ceb81246787177651fdea104bcb224c98cd512',
    'bigquery-catalogue-independent-r1-20261006': '142799a479783bbd6d5767edaa0ece8077d9f7ec61b95e51666006a04bd0021c',
}


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def check_evidence_additions(records, original_keys, approved):
    current_keys = {p.stem for p in records.glob('*.json') if p.name != 'README.json'}
    assert not set(approved) & original_keys, 'addition overlaps original evidence'
    assert current_keys == original_keys | set(approved), 'missing or unknown evidence record'
    for key, expected in approved.items():
        assert digest((records / (key + '.json')).read_bytes()) == expected, ('added evidence tampered', key)
    return set(approved)


def approved_evidence_additions(root):
    receipt = json.loads((root / RECEIPT.relative_to(ROOT)).read_text())
    directory = 'spec/research/public-data-fabric/evidence/$records'
    names = subprocess.check_output(['git', '-C', str(root), 'ls-tree', '-r', '--name-only', receipt['base_commit'], '--', directory], text=True).splitlines()
    original_keys = {Path(name).stem for name in names if name.endswith('.json') and Path(name).name != 'README.json'}
    assert len(original_keys) == 183, 'original evidence count'
    return check_evidence_additions(root / directory, original_keys, APPROVED_EVIDENCE_ADDITIONS)


def preserved_bytes(root, name, expected):
    """Keep old receipt hashes exact; only reviewed, named successor files differ."""
    current = (root / name).read_bytes()
    if digest(current) == expected:
        return current
    receipt = json.loads((root / RECEIPT.relative_to(ROOT)).read_text())
    assert name in SUPERSEDED_FILES, ('record or unrelated-file mutation', name)
    change = receipt['superseded_files'].get(name)
    assert change is not None, ('unrecorded successor', name)
    assert change['original_sha256'] == expected, ('old evidence changed', name)
    assert digest(current) == change['successor_sha256'], ('successor changed', name)
    original = subprocess.check_output(['git', '-C', str(root), 'show', receipt['base_commit'] + ':' + name])
    assert digest(original) == expected, ('historical object changed', name)
    return original


def main():
    receipt = json.loads(RECEIPT.read_text())
    assert set(receipt['superseded_files']) == SUPERSEDED_FILES
    original = json.loads(subprocess.check_output(['git', '-C', str(ROOT), 'show', receipt['base_commit'] + ':' + PORTFOLIO]))
    current = json.loads((ROOT / PORTFOLIO).read_text())
    decision = current['superseding_access_decision']
    assert {k: v for k, v in current.items() if k != 'superseding_access_decision'} == original
    assert decision['id'] == 'public-data-original-provider-no-copy-20261006'
    assert decision['retention'].startswith('none unless specifically user-authorized')
    assert decision['w1']['copied_provider_activation'] == 'frozen; authorization or working original-resource alternative pending'
    assert decision['ecb']['public_activation'].startswith('NO-GO')
    assert decision['bigquery']['priority'] == 'parallel A0 lane preserved'
    assert set(decision['recordset_ids']) == set(original['unique_public_recordsets'])
    assert len(original['unique_public_recordsets']) == 37 and len(original['original_scenario_dispositions']) == 57
    for name, change in receipt['superseded_files'].items():
        preserved_bytes(ROOT, name, change['original_sha256'])
    # A single Git tree listing verifies every original research/evidence byte.
    rows = subprocess.check_output(['git', '-C', str(ROOT), 'ls-tree', '-r', receipt['base_commit'], '--', 'spec/research/public-data-fabric'], text=True).splitlines()
    counts = {}
    names = [row.split('\t', 1)[1] for row in rows]
    hashes = subprocess.check_output(['git', '-C', str(ROOT), 'hash-object', '--stdin-paths'], input='\n'.join(names) + '\n', text=True).splitlines()
    assert len(hashes) == len(rows)
    for row, actual in zip(rows, hashes):
        meta, name = row.split('\t', 1)
        if name in receipt['superseded_files']:
            continue
        assert actual == meta.split()[2], ('research bytes changed', name)
        for collection in ['candidates', 'assertions', 'implementation', 'gates', 'questions', 'evidence']:
            if name.startswith('spec/research/public-data-fabric/' + collection + '/$records/') and name.endswith('.json') and Path(name).name != 'README.json':
                counts[collection] = counts.get(collection, 0) + 1
    # Only the three separately authorized, hash-bound catalogue evidence records
    # extend the unchanged original ledger. No other addition is permitted.
    additions = approved_evidence_additions(ROOT)
    for collection, count in counts.items():
        current_count = sum(p.name != 'README.json' for p in (BASE / collection / '$records').glob('*.json'))
        expected_count = count + (len(additions) if collection == 'evidence' else 0)
        assert current_count == expected_count, (collection, current_count, expected_count)
    assert counts['candidates'] == 699 and counts['assertions'] == 1177
    assert counts['evidence'] == 183 and counts['evidence'] + len(additions) == 186
    print(json.dumps({'unchanged_original_records': counts, 'hash_bound_evidence_additions': sorted(additions), 'total_evidence': 186, 'original_portfolio_fields_unchanged': True, 'scenarios': 7, 'unique_recordsets': 37, 'original_dispositions': 57, 'no_activation_or_copy_authorization': True}))


def self_test():
    # Fabricated metadata only: no provider values, body copies or live access.
    with tempfile.TemporaryDirectory(prefix='no-copy-preservation-') as directory:
        root = Path(directory)
        receipt_path = root / RECEIPT.relative_to(ROOT)
        receipt_path.parent.mkdir(parents=True)
        receipt_path.write_text(json.dumps({'superseded_files': {}}))
        checks = 0
        for collection in ['candidates', 'assertions', 'implementation', 'gates']:
            name = 'spec/research/public-data-fabric/' + collection + '/$records/synthetic.json'
            path = root / name
            path.parent.mkdir(parents=True)
            before = b'{"id":"synthetic","status":"needs-review"}\n'
            path.write_bytes(before)
            assert preserved_bytes(root, name, digest(before)) == before
            path.write_bytes(b'{"id":"synthetic","status":"validated"}\n')
            try:
                preserved_bytes(root, name, digest(before))
            except AssertionError:
                checks += 1
            else:
                raise AssertionError(('status mutation accepted', collection))
        original = b'{"historical":"synthetic"}\n'
        successor = b'{"historical":"synthetic","superseding_access_decision":{"retention":"none"}}\n'
        path = root / PORTFOLIO
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(b'{"superseding_access_decision":{"retention":"authorized"}}\n')
        receipt_path.write_text(json.dumps({'superseded_files': {PORTFOLIO: {'original_sha256': digest(original), 'successor_sha256': digest(successor)}}}))
        try:
            preserved_bytes(root, PORTFOLIO, digest(original))
        except AssertionError:
            checks += 1
        else:
            raise AssertionError('superseding decision tampering accepted')
        records = root / 'synthetic-evidence'
        records.mkdir()
        (records / 'original.json').write_bytes(b'{"historical":"synthetic"}\n')
        approved = {'addition-' + str(i): digest(b'{"synthetic":"documentation"}\n') for i in range(3)}
        for key in approved:
            (records / (key + '.json')).write_bytes(b'{"synthetic":"documentation"}\n')
        assert check_evidence_additions(records, {'original'}, approved) == set(approved)
        for key in approved:
            path = records / (key + '.json')
            path.write_bytes(b'{"synthetic":"changed"}\n')
            try:
                check_evidence_additions(records, {'original'}, approved)
            except AssertionError:
                checks += 1
            else:
                raise AssertionError(('added evidence tampering accepted', key))
            path.write_bytes(b'{"synthetic":"documentation"}\n')
        (records / 'unknown-extra.json').write_bytes(b'{}\n')
        try:
            check_evidence_additions(records, {'original'}, approved)
        except AssertionError:
            checks += 1
        else:
            raise AssertionError('unknown extra evidence accepted')
        (records / 'unknown-extra.json').unlink()
        (records / 'original.json').unlink()
        try:
            check_evidence_additions(records, {'original'}, approved)
        except AssertionError:
            checks += 1
        else:
            raise AssertionError('original evidence removal accepted')
        print(str(checks) + ' synthetic preservation/tamper refusals passed')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--self-test', action='store_true')
    args = parser.parse_args()
    main()
    if args.self_test:
        self_test()
