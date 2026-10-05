#!/usr/bin/env python3
"""Offline integrity and scope checks for the additive public component packet."""
import argparse
import hashlib
import json
from pathlib import Path
import re

HERE = Path(__file__).resolve().parent
BASE = HERE.parents[1]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def receipt_coordinates(original):
    """Extract only explicit original-document head fields, never index context."""
    if isinstance(original, list):
        return {'workflow_heads': [
            {'run_id': row['databaseId'], 'receipt_head': row['headSha']}
            for row in original]}
    coordinates = {}
    for field in ('head', 'head_sha'):
        if original.get(field) is not None:
            coordinates['receipt_head'] = original[field]
            coordinates['receipt_head_field'] = field
            break
    if original.get('reviewed_head') is not None:
        coordinates['reviewed_head'] = original['reviewed_head']
    return coordinates


def normalized_entry(index_entry, original):
    label = Path(index_entry['path']).name
    if label == 'receipt.json':
        label = 'canonical-sync-recovery-receipt.json'
    entry = {key: index_entry[key] for key in ('sha256', 'bytes')}
    entry['label'] = label
    for key in ('outcome', 'merge_sha', 'release'):
        if index_entry.get(key) is not None:
            entry[key] = index_entry[key]
    if index_entry.get('head') is not None:
        entry['context_head'] = index_entry['head']
    entry.update(receipt_coordinates(original))
    return entry


def check_original_associations(index_path, manifest, public_receipts):
    raw_index = index_path.read_bytes()
    assert digest(raw_index) == manifest['root_index_original_sha256']
    assert len(raw_index) == manifest['root_index_original_bytes']
    index = json.loads(raw_index)
    assert len(index['receipts']) == len(public_receipts) == 76
    for indexed, public in zip(index['receipts'], public_receipts):
        raw = Path(indexed['path']).read_bytes()
        assert digest(raw) == indexed['sha256'] and len(raw) == indexed['bytes']
        original = json.loads(raw) if raw else {}
        assert public == normalized_entry(indexed, original), public['label']
    return len(public_receipts)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--original-index', type=Path,
                        help='optional private frozen index; rehash and associate every original receipt')
    args = parser.parse_args()
    manifest = json.loads((HERE / 'source-manifest.json').read_text())
    proof = json.loads((HERE / 'evidence.json').read_text())
    receipts = json.loads((HERE / 'publication-receipts.json').read_text())
    for name, expected in manifest['packet_files'].items():
        assert digest((HERE / name).read_bytes()) == expected, name
    for name, expected in manifest['preserved_research_files'].items():
        folder = BASE / name
        if folder.is_file():
            count, actual = 1, digest(folder.read_bytes())
        else:
            entries = [(str(p.relative_to(folder)), digest(p.read_bytes()))
                       for p in sorted(folder.rglob('*'))
                       if p.is_file() and '__pycache__' not in p.parts]
            count = len(entries)
            actual = digest(json.dumps(entries, separators=(',', ':')).encode())
        assert count == expected['files'] and actual == expected['sha256'], name
    assert len(receipts['receipts']) == manifest['root_receipt_count'] == 76
    for entry in receipts['receipts']:
        assert 'head' not in entry, entry['label']
        for field in ('context_head', 'receipt_head', 'reviewed_head'):
            if field in entry:
                assert re.fullmatch('[0-9a-f]{40}', entry[field]), entry['label']
        if 'receipt_head' in entry:
            assert entry['receipt_head_field'] in {'head', 'head_sha'}
        for workflow in entry.get('workflow_heads', []):
            assert isinstance(workflow['run_id'], int)
            assert re.fullmatch('[0-9a-f]{40}', workflow['receipt_head'])
    # A context mismatch must remain visible rather than replacing the actual head.
    control = {'path': 'control.json', 'bytes': 1, 'sha256': '0' * 64,
               'head': 'a' * 40, 'outcome': 'passed'}
    expected = normalized_entry(control, {'head': 'b' * 40})
    assert expected['context_head'] == 'a' * 40 and expected['receipt_head'] == 'b' * 40
    assert normalized_entry(control, {'head': 'a' * 40}) != expected
    assert receipt_coordinates([{'databaseId': 1, 'headSha': 'b' * 40}]) == {
        'workflow_heads': [{'run_id': 1, 'receipt_head': 'b' * 40}]}
    original_count = 0
    if args.original_index:
        original_count = check_original_associations(args.original_index, manifest, receipts['receipts'])
    for entry in list(manifest['original_document_hash_declarations'].values()) + receipts['receipts']:
        assert re.fullmatch('[0-9a-f]{64}', entry['sha256'])
        assert isinstance(entry['bytes'], int) and entry['bytes'] >= 0
    assert proof['kind'] == 'w1-component-publication-update'
    assert len(proof['components']) == 6
    for component in proof['components']:
        assert component['reviewed_tree'] == component['merge_tree']
        assert component['canonical_sync_clean'] and component['source_branch_worktree_absent']
        assert component['main_ci'] in {'pending', 'passed'}
    assert proof['programme']['w1_complete'] is False
    assert proof['programme']['b1_open'] is True
    assert not proof['programme']['phase5_passed'] and not proof['programme']['phase6_passed']
    for candidate in proof['provider_candidates']:
        assert candidate['admitted'] is False
    implementation = list((BASE / 'implementation/$records').glob('w1-*.json'))
    assert len(implementation) == 5
    for path in implementation:
        row = json.loads(path.read_text())
        assert row['status'] == 'implementing'
        assert row['publication_evidence_ids'] == [] and row['query_evidence_ids'] == []
    for path in HERE.iterdir():
        if path.suffix in {'.json', '.md'}:
            assert not re.search(rb'/(?:private/tmp|Users/[^/]+)/', path.read_bytes()), path.name
    print(f"verified 2 packet hashes, {len(manifest['preserved_research_files'])} preserved groups/files, "
          f"{len(receipts['receipts'])} receipt summaries and 5 implementing states; "
          f"{original_count} private originals rehashed and head-associated")


if __name__ == '__main__':
    main()
