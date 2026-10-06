#!/usr/bin/env python3
"""Check frozen packet, source associations and unchanged ledger; no admission/fetch."""
import argparse
import hashlib
import json
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[4]
BASE = HERE.parents[1]

def digest(raw):
    return hashlib.sha256(raw).hexdigest()

def read(path):
    return json.loads(path.read_text())

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--original-map', type=Path)
    args = parser.parse_args()
    manifest = read(HERE / 'source-manifest.json')
    proof = read(HERE / 'evidence.json')
    docs = read(HERE / 'reviewed-records.json')['documents']
    for name, expected in manifest['packet_files'].items():
        assert digest((HERE / name).read_bytes()) == expected, name
    for label, doc in docs.items():
        assert digest(doc['content'].encode()) == doc['public_sha256'], label
        assert len(doc['source_sha256']) == 64 and doc['source_bytes'] > 0, label
    document = lambda label: json.loads(docs[label]['content'])
    review = document('fx-review/review-r1.json')
    for name, expected in review['reviewed_packet_sha256'].items():
        label = 'fx/' + name
        if label in docs:
            assert docs[label]['source_sha256'] == expected, label
    for name, expected in review['receipts'].items():
        label = 'fx-review/' + name
        if label in docs:
            assert docs[label]['source_sha256'] == expected, label
    assert all(review[k] == 0 for k in ['blockers', 'majors', 'minors'])
    assert document('fx-review/verification.json')['all_checks_pass']
    acceptance = document('fx/root-acceptance-r1.json')
    assert not any(acceptance[k] for k in ['canonical_semantic_edits_approved', 'source_selected', 'ovdb_publication_approved'])
    assert len(acceptance['remaining_gates']) == 5
    for scope, label in [('transport', 'cloud-pr24'), ('app', 'app-pr197')]:
        pr = document('remote/' + label + '.json')
        assert pr['state'] == 'MERGED'
        assert pr['headRefOid'] == proof[scope]['head']
        assert pr['mergeCommit']['oid'] == proof[scope]['merge']
    options, get = document('transport/options-r1.json'), document('transport/get-r1.json')
    assert options['pass'] and get['exposure_pass']
    assert len(options['rows']) == len(get['rows']) == 6
    for row in options['rows']:
        assert row['status'] == 204 and row['acao'] == 'https://datatug.app' and not row['missing_pins']
    for row in get['rows']:
        assert row['status'] == 200 and row['acao'] == 'https://datatug.app'
        assert not row['missing_exposed_pins'] and not row['present_pin_values']
    assert document('transport/negative-r1.json')['acao'] is None
    assert not proof['transport']['immutable_runtime_accepted']
    assert proof['app']['exact_main_ci'] == 'pending'
    assert not proof['app']['public_sources_admitted']
    assert not any(proof['fx'][k] for k in ['selected', 'implemented', 'published', 'canonical_assertion_edits', 'new_core_concept_approved', 'implementation_policy_approved'])
    assert proof['programme']['b1_open'] and proof['programme']['phases4_5_6_open']
    assert not any(proof['programme'][k] for k in ['w1_complete', 'geo_ror_admission', 'availability_changed', 'paid_query', 'cloud_run_test'])
    for name, expected in manifest['preserved_research_files'].items():
        assert digest((ROOT / name).read_bytes()) == expected, name
    assert len(list((BASE / 'candidates/$records').glob('*.json'))) == 699
    assert len(list((BASE / 'assertions/$records').glob('*.json'))) == 1177
    originals = 0
    if args.original_map:
        paths = read(args.original_map)
        assert set(paths) == set(docs)
        for label, path in paths.items():
            raw = Path(path).read_bytes()
            assert digest(raw) == docs[label]['source_sha256'], label
            assert len(raw) == docs[label]['source_bytes'], label
            originals += 1
    print(json.dumps({'packet_files': len(manifest['packet_files']), 'source_documents': len(docs), 'private_originals': originals, 'preserved_research_files': len(manifest['preserved_research_files']), 'candidates': 699, 'assertions': 1177, 'availability_changed': False, 'app_main_ci': proof['app']['exact_main_ci']}))

if __name__ == '__main__':
    main()
