#!/usr/bin/env python3
"""Verify the superseding metadata decision and retained historical evidence offline."""
import argparse
import hashlib
import json
import shutil
import subprocess
import tempfile
from pathlib import Path

from verify_s2_scope_decision import SUCCESSORS as S2_SUCCESSORS, predecessor_bytes as s2_predecessor, check as check_s2_scope, self_test as s2_self_test

BASE = Path(__file__).resolve().parent
ROOT = BASE.parents[2]
PORTFOLIO = 'spec/research/public-data-fabric/evidence/p3-portfolio.json'
RECEIPT = BASE / 'evidence/no-copy-access-decision.json'
SUPERSEDED_FILES = frozenset(['spec/research/public-data-fabric/p3-portfolio.md', 'spec/research/public-data-fabric/evidence/p3-portfolio.json', 'spec/research/public-data-fabric/evidence/w1-followup-decisions/verify.py', 'spec/research/public-data-fabric/evidence/w1-transport-app-fx/verify.py', '.github/workflows/public-data-fabric.yml', 'spec/research/public-data-fabric/README.md', 'spec/features/public-data-fabric/README.md', 'spec/research/public-data-fabric/p2-reconciliation.md'])
APPROVED_EVIDENCE_ADDITIONS = {
    'bigquery-catalogue-corrections-20261006': '4d7da4bf0c36c71ee2caba31d7876fe833407920af55cfda88bd3df4571d7ca3',
    'bigquery-catalogue-discovery-20261006': '767265e46e8c234afaa9048460ceb81246787177651fdea104bcb224c98cd512',
    'bigquery-catalogue-independent-r1-20261006': '142799a479783bbd6d5767edaa0ece8077d9f7ec61b95e51666006a04bd0021c',
}


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


CANDIDATE_RECEIPT = BASE / 'evidence/bigquery-candidate-additions.json'
CANDIDATE_BASE = '298397fcce9963750fadbd9afd3589b974d50a63'
CANDIDATE_ADDITIONS_SHA256 = '2a13a14c2d27c02a47c0a1285fb8269087a95b4d4ff114ce7892a7e2293d6219'
CANDIDATE_VERIFIERS = frozenset(['spec/research/public-data-fabric/verify_no_copy_decision.py', 'spec/research/public-data-fabric/evidence/w1-followup-decisions/verify.py', 'spec/research/public-data-fabric/evidence/w1-transport-app-fx/verify.py', 'spec/research/public-data-fabric/evidence/bigquery-dual-runtime/verify.py'])
COLLECTIONS = ('candidates', 'assertions', 'evidence', 'questions', 'implementation', 'gates')
# Pin successor bytes independently of their editable declaration. The verifier's
# own digest normalizes only this one value to avoid a circular raw-file hash.
CANDIDATE_SELF_SHA256 = '63b9e74c1455572453b23d664d9d4b4cecee50c4a5023f19a4655268238d8a57'
CANDIDATE_SUCCESSOR_SHA256 = {
    'spec/research/public-data-fabric/evidence/w1-followup-decisions/verify.py': 'ec29ec51ac3f659f019b214d0b398dfb098f82376784486d9267edb2b782da2a',
    'spec/research/public-data-fabric/evidence/w1-transport-app-fx/verify.py': 'a7dc9a78299af75360388b063a77baff33a90c5422dacea0b0d5242b4e947a65',
    'spec/research/public-data-fabric/evidence/bigquery-dual-runtime/verify.py': 'f697f059883f379584bf3781d56277dba2d2d23ebf2cb59120f0a22b7475ee88',
}


def check_verifier_successor(name, raw):
    if name == 'spec/research/public-data-fabric/verify_no_copy_decision.py':
        assignment = ("CANDIDATE_SELF_SHA256 = '" + CANDIDATE_SELF_SHA256 + "'").encode()
        normalized = ("CANDIDATE_SELF_SHA256 = '" + '0' * 64 + "'").encode()
        assert raw.count(assignment) == 1, 'self verifier pin changed'
        assert digest(raw.replace(assignment, normalized)) == CANDIDATE_SELF_SHA256, 'self verifier bytes changed'
    else:
        assert digest(raw) == CANDIDATE_SUCCESSOR_SHA256[name], ('unreviewed verifier successor', name)



def candidate_receipt(root):
    receipt = json.loads((root / CANDIDATE_RECEIPT.relative_to(ROOT)).read_text())
    assert receipt['format'] == 'data-fabric-candidate-additions/1'
    assert receipt['base_commit'] == CANDIDATE_BASE, 'candidate baseline changed'
    assert receipt['reviewed_proposal_sha256'] == '1c46a70dbc7736cec093adfbae7db74d91022caf0e35b80e24646bf788c7b116'
    raw = json.dumps(receipt['additions'], sort_keys=True, separators=(',', ':')).encode()
    assert digest(raw) == CANDIDATE_ADDITIONS_SHA256, 'reviewed addition hashes changed'
    assert set(receipt['verifier_successors']) == CANDIDATE_VERIFIERS, 'unknown verifier successor'
    for name, change in receipt['verifier_successors'].items():
        original = subprocess.check_output(['git', '-C', str(root), 'show', CANDIDATE_BASE + ':' + name])
        assert digest(original) == change['original_sha256'], ('verifier original changed', name)
        current = s2_predecessor(root, name, (root / name).read_bytes())
        check_verifier_successor(name, current)
        assert digest(current) == change['successor_sha256'], ('verifier successor changed', name)
    return receipt


def approved_record_additions(root):
    """Allow only the reviewed 40 additions; retain all 2,121 original blobs."""
    receipt = candidate_receipt(root)
    prefix = str(BASE.relative_to(ROOT)) + '/'
    rows = subprocess.check_output(['git', '-C', str(root), 'ls-tree', '-r', CANDIDATE_BASE, '--', prefix], text=True).splitlines()
    baseline = {}
    for row in rows:
        meta, name = row.split('\t', 1)
        relative = Path(name).relative_to(prefix)
        if len(relative.parts) == 3 and relative.parts[0] in COLLECTIONS and relative.parts[1] == '$records' and relative.suffix == '.json':
            baseline[name] = meta.split()[2]
    assert len(baseline) == 2121, 'candidate original inventory changed'
    names = sorted(baseline)
    actual = subprocess.check_output(['git', '-C', str(root), 'hash-object', '--stdin-paths'], input='\n'.join(names) + '\n', text=True, stderr=subprocess.PIPE).splitlines()
    for i, name in enumerate(names):
        if name in S2_SUCCESSORS:
            previous = s2_predecessor(root, name, (root / name).read_bytes())
            actual[i] = subprocess.check_output(['git', '-C', str(root), 'hash-object', '--stdin'], input=previous).decode().strip()
    assert actual == [baseline[name] for name in names], 'original ledger record changed'
    result = {}
    for table in COLLECTIONS:
        records = root / prefix / table / '$records'
        original_keys = {Path(name).stem for name in baseline if name.startswith(prefix + table + '/')}
        approved = receipt['additions'].get(table, {})
        result[table] = check_evidence_additions(records, original_keys, approved)
    assert {t: len(ids) for t, ids in result.items()} == {'candidates': 15, 'assertions': 10, 'evidence': 10, 'questions': 5, 'implementation': 0, 'gates': 0}
    return result


def verifier_predecessor(root, name, current):
    """Walk only hash-bound verifier successors back to their historical bytes."""
    if name not in CANDIDATE_VERIFIERS or not (root / CANDIDATE_RECEIPT.relative_to(ROOT)).is_file():
        return current
    receipt = candidate_receipt(root)
    return subprocess.check_output(['git', '-C', str(root), 'show', receipt['base_commit'] + ':' + name])


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
    approved = {**APPROVED_EVIDENCE_ADDITIONS}
    candidate = candidate_receipt(root)['additions']['evidence']
    assert not set(approved) & set(candidate), 'overlapping evidence additions'
    approved.update(candidate)
    return check_evidence_additions(root / directory, original_keys, approved)


def preserved_bytes(root, name, expected):
    """Keep old receipt hashes exact; only reviewed, named successor files differ."""
    current = s2_predecessor(root, name, (root / name).read_bytes())
    if digest(current) == expected:
        return current
    current = verifier_predecessor(root, name, current)
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
    check_s2_scope(ROOT)
    candidate_additions = approved_record_additions(ROOT)
    receipt = json.loads(RECEIPT.read_text())
    assert set(receipt['superseded_files']) == SUPERSEDED_FILES
    original = json.loads(subprocess.check_output(['git', '-C', str(ROOT), 'show', receipt['base_commit'] + ':' + PORTFOLIO]))
    current = json.loads(s2_predecessor(ROOT, PORTFOLIO, (ROOT / PORTFOLIO).read_bytes()))
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
    hashes = subprocess.check_output(['git', '-C', str(ROOT), 'hash-object', '--stdin-paths'], input='\n'.join(names) + '\n', text=True, stderr=subprocess.PIPE).splitlines()
    assert len(hashes) == len(rows)
    for row, actual in zip(rows, hashes):
        meta, name = row.split('\t', 1)
        if name in receipt['superseded_files']:
            continue
        if name in CANDIDATE_VERIFIERS:
            historical = subprocess.check_output(['git', '-C', str(ROOT), 'show', receipt['base_commit'] + ':' + name])
            assert preserved_bytes(ROOT, name, digest(historical)) == historical
            continue
        if name in S2_SUCCESSORS:
            previous = s2_predecessor(ROOT, name, (ROOT / name).read_bytes())
            actual = subprocess.check_output(['git', '-C', str(ROOT), 'hash-object', '--stdin'], input=previous).decode().strip()
        assert actual == meta.split()[2], ('research bytes changed', name)
        for collection in ['candidates', 'assertions', 'implementation', 'gates', 'questions', 'evidence']:
            if name.startswith('spec/research/public-data-fabric/' + collection + '/$records/') and name.endswith('.json') and Path(name).name != 'README.json':
                counts[collection] = counts.get(collection, 0) + 1
    # Catalogue evidence and the exact reviewed candidate additions are hash-bound.
    # Unknown additions and all original-record changes remain forbidden.
    additions = approved_evidence_additions(ROOT)
    for collection, count in counts.items():
        current_count = sum(p.name != 'README.json' for p in (BASE / collection / '$records').glob('*.json'))
        expected_count = count + (len(additions) if collection == 'evidence' else len(candidate_additions[collection]))
        assert current_count == expected_count, (collection, current_count, expected_count)
    assert counts['candidates'] == 699 and counts['assertions'] == 1177
    assert counts['evidence'] == 183 and counts['evidence'] + len(additions) == 196
    print(json.dumps({'unchanged_original_records': counts, 'hash_bound_evidence_additions': sorted(additions), 'total_evidence': 196, 'hash_bound_candidate_additions': {t: sorted(ids) for t, ids in candidate_additions.items()}, 'original_portfolio_fields_unchanged': True, 'scenarios': 7, 'unique_recordsets': 37, 'original_dispositions': 57, 'no_activation_or_copy_authorization': True}))


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


def candidate_additions_self_test():
    # Reuse local Git objects and copy ledger metadata only; no network/providers.
    with tempfile.TemporaryDirectory(prefix='candidate-preservation-') as directory:
        root = Path(directory) / 'ledger'
        subprocess.run(['git', 'clone', '--quiet', '--shared', '--no-checkout', str(ROOT), str(root)], check=True)
        for table in COLLECTIONS:
            relative = BASE.relative_to(ROOT) / table / '$records'
            shutil.copytree(ROOT / relative, root / relative)
        for name in CANDIDATE_VERIFIERS | set(S2_SUCCESSORS) | {
                'spec/research/public-data-fabric/evidence/s2-initial-scope-decision.json',
                'spec/research/public-data-fabric/evidence/s2-scope-successors.json',
                str(CANDIDATE_RECEIPT.relative_to(ROOT)), str(RECEIPT.relative_to(ROOT))}:
            path = root / name
            path.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(ROOT / name, path)
        assert sum(map(len, approved_record_additions(root).values())) == 40
        packet = root / CANDIDATE_RECEIPT.relative_to(ROOT)
        original = root / BASE.relative_to(ROOT) / 'gates/$records/phase4.json'
        receipt = candidate_receipt(root)
        key = next(iter(receipt['additions']['candidates']))
        added = root / BASE.relative_to(ROOT) / 'candidates/$records' / (key + '.json')
        unknown = added.with_name('synthetic-unknown-addition.json')
        verifier = root / next(iter(CANDIDATE_VERIFIERS))
        checks = 0

        def refuse(label, path, mutation, call=None):
            nonlocal checks
            before = path.read_bytes() if path.exists() else None
            try:
                mutation(path)
                try:
                    (call or (lambda: approved_record_additions(root)))()
                except (AssertionError, FileNotFoundError, subprocess.CalledProcessError):
                    checks += 1
                else:
                    raise AssertionError('accepted invalid candidate mutation: ' + label)
            finally:
                if before is None:
                    path.unlink(missing_ok=True)
                else:
                    path.write_bytes(before)

        refuse('original removal', original, lambda p: p.unlink())
        refuse('original gate changed', original, lambda p: p.write_bytes(b'{}\n'))
        refuse('approved addition removal', added, lambda p: p.unlink())
        refuse('approved addition changed', added, lambda p: p.write_bytes(b'{}\n'))
        refuse('unknown addition', unknown, lambda p: p.write_bytes(b'{}\n'))
        refuse('verifier changed', verifier, lambda p: p.write_bytes(b'# altered verifier\n'))

        def change_receipt(path, mutate):
            value = json.loads(path.read_text())
            mutate(value)
            path.write_text(json.dumps(value))

        refuse('forged addition hash', packet, lambda p: change_receipt(p, lambda r: r['additions']['candidates'].update({key: '0' * 64})))
        refuse('baseline changed', packet, lambda p: change_receipt(p, lambda r: r.update(base_commit='0' * 40)))
        refuse('unknown verifier', packet, lambda p: change_receipt(p, lambda r: r['verifier_successors'].update({'unknown.py': {}})))
        refuse('forged verifier hash', packet, lambda p: change_receipt(p, lambda r: r['verifier_successors'][str(verifier.relative_to(root))].update(successor_sha256='0' * 64)))
        name = str(verifier.relative_to(root))
        expected = receipt['verifier_successors'][name]['original_sha256']
        assert digest(preserved_bytes(root, name, expected)) == expected
        refuse('successor tampering on historical lookup', verifier, lambda p: p.write_bytes(b'# altered verifier\n'), lambda: preserved_bytes(root, name, expected))
        # A declaration edited alongside its verifier must never authorize new
        # bytes. Exercise all four successor files through both approval paths.
        for name in sorted(CANDIDATE_VERIFIERS):
            successor = root / name
            old_successor = successor.read_bytes()
            old_packet = packet.read_bytes()
            try:
                successor.write_bytes(b'# arbitrary unreviewed successor\n')
                value = json.loads(old_packet)
                value['verifier_successors'][name]['successor_sha256'] = digest(successor.read_bytes())
                packet.write_text(json.dumps(value))
                for call in [lambda: approved_record_additions(root),
                             lambda: preserved_bytes(root, name, value['verifier_successors'][name]['original_sha256'])]:
                    try:
                        call()
                    except AssertionError:
                        checks += 1
                    else:
                        raise AssertionError('accepted paired verifier/declaration mutation: ' + name)
            finally:
                successor.write_bytes(old_successor)
                packet.write_bytes(old_packet)
        assert sum(map(len, approved_record_additions(root).values())) == 40
        print(str(checks) + ' candidate addition/original/successor tamper refusals passed')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--self-test', action='store_true')
    args = parser.parse_args()
    main()
    if args.self_test:
        self_test()
        candidate_additions_self_test()
        s2_self_test()
