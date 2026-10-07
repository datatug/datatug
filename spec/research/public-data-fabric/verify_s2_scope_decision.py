#!/usr/bin/env python3
"""Check the bounded root S2 scope decision and exact metadata predecessor chain."""
import argparse
import hashlib
import json
import shutil
import subprocess
import tempfile
from pathlib import Path

BASE_COMMIT = '5c2ca879ef5fc1a3e8ebb20b4c4b847fd161febe'
PREFIX = 'spec/research/public-data-fabric/'
ROOT = Path(__file__).resolve().parents[3]
DECISION = PREFIX + 'evidence/s2-initial-scope-decision.json'
RECEIPT = PREFIX + 'evidence/s2-scope-successors.json'
DECISION_SHA256 = 'ef31c7626b200b2777dd3ca2e4264067349b1d794d87df996ed4adf3661e7f91'
# Independent pins: an editable receipt cannot authorize arbitrary successors.
SUCCESSORS = {'spec/research/public-data-fabric/candidates/$records/p3-scenario-catalogue.json': {'original_sha256': '97f35e5c3ecde4d77e4bd8568e6f2b174d936f2af0e875a6e5f83052436b45fc',
                                                                                     'successor_sha256': '085f05fdf66f732e8d37e8d06291fcf2d8d59783278151a57e040d7339fe4913'},
 'spec/research/public-data-fabric/questions/$records/p3-followup-2.json': {'original_sha256': '678c59c91dfa986ab6b97eee51cc0e24933556f7f87287830add057856caf489',
                                                                            'successor_sha256': '4e70e5add7964fff6068fd7c87e9e9e477dda5e9c27c2e9c622a4cdf6993d0a8'},
 'spec/research/public-data-fabric/evidence/p3-portfolio.json': {'original_sha256': '03bcc2a1f2bec0fafe0f9791b171dc000d1ffa67fc28db36b19a1d155972e0fa',
                                                                 'successor_sha256': 'a6dc005b6054e56ee758dc08280dd982ce941e2bac2ef6ff5e25f502dcb118e1'},
 'spec/research/public-data-fabric/p3-portfolio.md': {'original_sha256': 'ffe97bcb6d319ea4820679a02e66b02f3e169f6e1c95154e63687613aa296b9e',
                                                      'successor_sha256': '4d9dec978cefe7f304c1d1187f244572b91ff2f77f9ae5192c069257b238eb51'},
 'spec/research/public-data-fabric/verify_no_copy_decision.py': {'original_sha256': '26fe40ef5cb81f9a8ab03b013fbd27a81ebd418a15af7eda35a542a836f7c662',
                                                                 'successor_sha256': 'fd65c35e9dc65a8a55783a405311ede129488b370c234ac26cdb28c138ed1e06'},
 'spec/research/public-data-fabric/evidence/w1-followup-decisions/verify.py': {'original_sha256': 'ec29ec51ac3f659f019b214d0b398dfb098f82376784486d9267edb2b782da2a',
                                                                               'successor_sha256': '64e966596822d44b9ef7c319b4fa33c4e9bbd6c32a33a306717b1c5f16549c5d'},
 'spec/research/public-data-fabric/evidence/bigquery-dual-runtime/verify.py': {'original_sha256': 'f697f059883f379584bf3781d56277dba2d2d23ebf2cb59120f0a22b7475ee88',
                                                                               'successor_sha256': '12d692407e23d631573677789cc330f42664e6e0077c24b9f1550058f828fd60'}}


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def predecessor_bytes(root, name, current):
    """Unwrap only the seven reviewed metadata successors; all other bytes pass through."""
    receipt_path = root / RECEIPT
    if name not in SUCCESSORS or not receipt_path.is_file():
        return current
    receipt = json.loads(receipt_path.read_text())
    assert receipt['format'] == 'data-fabric-s2-scope-successors/1'
    assert receipt['base_commit'] == BASE_COMMIT
    assert receipt['decision_path'] == DECISION
    assert receipt['decision_sha256'] == DECISION_SHA256
    assert receipt['successors'] == SUCCESSORS, 'unknown S2 metadata successor'
    assert digest((root / DECISION).read_bytes()) == DECISION_SHA256, 'S2 scope decision changed'
    pin = SUCCESSORS[name]
    assert digest(current) == pin['successor_sha256'], ('unreviewed S2 successor', name)
    previous = subprocess.check_output(['git', '-C', str(root), 'show', BASE_COMMIT + ':' + name])
    assert digest(previous) == pin['original_sha256'], ('S2 predecessor changed', name)
    return previous


def check(root=ROOT):
    assert (root / RECEIPT).is_file(), 'S2 scope receipt missing'
    decision = json.loads((root / DECISION).read_text())
    assert digest((root / DECISION).read_bytes()) == DECISION_SHA256
    originals = {}
    for name in SUCCESSORS:
        originals[name] = predecessor_bytes(root, name, (root / name).read_bytes())
    assert decision['id'] == 's2-initial-non-fx-scope-20261007'
    assert decision['base_commit'] == BASE_COMMIT
    assert decision['question_status'] == 'open'
    assert decision['question_classification'] == 'blocking'
    for flag in ['source_admission_changed', 'semantic_assertions_changed', 'candidate_statuses_changed',
                 'phase_gate_statuses_changed', 'implementation_approved', 'source_rights_accepted',
                 'provider_reads_authorized', 'retention_authorized']:
        assert decision[flag] is False, flag
    assert decision['initial_assertion_ids'] == ['p2-edge-aw-mass-units', 'p2-edge-aw-language-cultures',
                                               'p2-edge-aw-alternate-product-number', 'p2-map-chinook-catalogue-price']
    assert decision['deferred_assertion_ids'] == ['d1-edge-aw-fx-currentref']
    assert decision['retained_candidate_ids'] == ['d1-ecb-fx', 'd1-scenario-fx-audit']
    for id in decision['initial_assertion_ids']:
        path = root / PREFIX / 'assertions/$records' / (id + '.json')
        assertion = json.loads(path.read_text())
        assert assertion['status'] == 'reconciled', ('unaccepted initial assertion', id)
        assert assertion['compatibility'] != 'unresolved', id
    # Narrow only descriptive fields, never lifecycle or relationships.
    candidate = PREFIX + 'candidates/$records/p3-scenario-catalogue.json'
    before = json.loads(originals[candidate]); after = json.loads((root / candidate).read_text())
    changed = {key for key in before.keys() | after.keys() if before.get(key) != after.get(key)}
    assert changed == {'title', 'reason', 'practical_value'}
    assert after['reason'].startswith(before['reason']), 'historical scenario prose lost'
    question = PREFIX + 'questions/$records/p3-followup-2.json'
    before = json.loads(originals[question]); after = json.loads((root / question).read_text())
    assert {key for key in before.keys() | after.keys() if before.get(key) != after.get(key)} == {'resolution'}
    assert after['status'] == 'open' and after['classification'] == 'blocking'
    assert after['resolution'].startswith(before['resolution']), 'historical question requirements lost'
    portfolio = PREFIX + 'evidence/p3-portfolio.json'
    before = json.loads(originals[portfolio]); after = json.loads((root / portfolio).read_text())
    assert after.pop('superseding_s2_scope_decision') == decision
    assert after == before, 'historical portfolio fields changed'
    print('S2 scope: four accepted non-FX paths; historical FX deferred; composite question open; exact predecessors retained')


def self_test():
    """Mutation checks copy authored metadata only; never source rows or provider payloads."""
    with tempfile.TemporaryDirectory(prefix='fabric-s2-scope-metadata-') as directory:
        root = Path(directory)
        # Git queries still address the existing repository; no new Git snapshot.
        (root / '.git').write_text('gitdir: ' + subprocess.check_output(
            ['git', '-C', str(ROOT), 'rev-parse', '--absolute-git-dir'], text=True).strip() + '\n')
        names = list(SUCCESSORS) + [DECISION, RECEIPT]
        decision = json.loads((ROOT / DECISION).read_text())
        names += [PREFIX + 'assertions/$records/' + id + '.json' for id in decision['initial_assertion_ids']]
        for name in names:
            target = root / name; target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(ROOT / name, target)
        check(root)
        cases = [(PREFIX + 'questions/$records/p3-followup-2.json', lambda d: d.update(status='resolved')),
                 (PREFIX + 'candidates/$records/p3-scenario-catalogue.json', lambda d: d.update(status='validated')),
                 (DECISION, lambda d: d.update(retention_authorized=True)),
                 (RECEIPT, lambda d: d['successors'].update({'unrelated.json': {'original_sha256': '0'*64, 'successor_sha256': '0'*64}}))]
        for name, mutate in cases:
            path = root / name; original = path.read_bytes(); value = json.loads(original)
            mutate(value); path.write_text(json.dumps(value))
            try:
                check(root)
            except AssertionError:
                pass
            else:
                raise AssertionError('S2 invalid mutation accepted: ' + name)
            path.write_bytes(original)
        # A matching edited receipt must not launder modified scope/status bytes.
        name = PREFIX + 'questions/$records/p3-followup-2.json'
        path = root / name; original = path.read_bytes(); value = json.loads(original)
        value['status'] = 'resolved'; path.write_text(json.dumps(value))
        receipt = json.loads((root / RECEIPT).read_text())
        receipt['successors'][name]['successor_sha256'] = digest(path.read_bytes())
        (root / RECEIPT).write_text(json.dumps(receipt))
        try:
            predecessor_bytes(root, name, path.read_bytes())
        except AssertionError:
            pass
        else:
            raise AssertionError('S2 matching receipt laundered question closure')
    print('S2 scope: five negative authored-metadata mutations refused')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--self-test', action='store_true')
    args = parser.parse_args()
    check()
    if args.self_test:
        self_test()
