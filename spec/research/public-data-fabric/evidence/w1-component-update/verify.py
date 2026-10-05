#!/usr/bin/env python3
"""Offline integrity and scope checks for the additive public component packet."""
import hashlib
import json
from pathlib import Path
import re

HERE = Path(__file__).resolve().parent
BASE = HERE.parents[1]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def main():
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
    assert len(receipts['receipts']) == manifest['root_receipt_count']
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
          f"{len(receipts['receipts'])} receipt summaries and 5 implementing states; private originals not rehashed")


if __name__ == '__main__':
    main()
