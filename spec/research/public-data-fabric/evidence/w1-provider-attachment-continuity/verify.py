#!/usr/bin/env python3
"""Offline continuity checks; optional original receipts and immutable Git objects."""
import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

HERE = Path(__file__).resolve().parent
BASE = HERE.parents[1]


def sha256(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', action='append', default=[], metavar='GITHUB_URL=LOCAL_CHECKOUT',
                        help='optionally rehash immutable provider and decision objects; no fetching')
    parser.add_argument('--review-originals', type=Path, help='optional original provider-review directory')
    parser.add_argument('--landing-originals', type=Path, help='optional original provider landing receipt directory')
    args = parser.parse_args()
    repos = dict(value.split('=', 1) for value in args.repo)
    proof = json.loads((HERE / 'evidence.json').read_text())
    record = json.loads((BASE / 'evidence/$records/w1-provider-attachment-continuity.json').read_text())
    assert proof['kind'] == 'w1-landed-provider-attachment-continuity'
    assert proof['packet_review'] == 'pending separate adversarial review'
    assert record['independent'] is False and record['outcome'] == 'observation'
    assert record['path'] == str((HERE / 'evidence.json').relative_to(BASE.parents[2]))
    assert proof['verification']['released_harness_rerun'] is False
    assert proof['review_history']['r2']['released_reader']['version'] == 'v0.27.1'
    assert proof['review_history']['r2']['released_reader']['replacement'] is False
    assert proof['review_history']['r2']['released_reader']['fallback_resolver_used'] is False
    assert proof['review_history']['r2']['results']['geo']['source_data_stage'] == 'outside_scope'
    assert proof['review_history']['r2']['results']['ror']['source_data_refs'] == 1
    originals = 0
    for name, entry in proof['original_documents'].items():
        assert re.fullmatch('[0-9a-f]{64}', entry['sha256']) and entry['bytes'] > 0
        root = args.review_originals if entry['original_group'] == 'provider-review' else args.landing_originals
        if root:
            raw = (root / name).read_bytes()
            assert sha256(raw) == entry['sha256'] and len(raw) == entry['bytes'], name
            originals += 1
    for round_name in ['r1', 'r2']:
        for ext in ['md', 'json']:
            assert proof['review_history'][round_name][ext + '_sha256'] == proof['original_documents'][f'review-{round_name}.{ext}']['sha256']
    objects = 0

    def blob(repository, revision, path, expected):
        nonlocal objects
        root = repos.get(repository)
        if not root:
            return None
        entry = subprocess.check_output(['git', '-C', root, 'ls-tree', revision, '--', path]).split()
        assert entry[0] in [b'100644', b'100755'], path
        data = subprocess.check_output(['git', '-C', root, 'show', f'{revision}:{path}'])
        assert sha256(data) == expected, path
        objects += 1
        return data

    def references(value, repository, revision):
        if isinstance(value, dict):
            if {'path', 'sha256'} <= value.keys():
                blob(value.get('repository', repository), value.get('revision', value.get('commit', revision)), value['path'], value['sha256'])
            for child in value.values():
                references(child, repository, revision)
        elif isinstance(value, list):
            for child in value:
                references(child, repository, revision)

    assert len(proof['providers']) == 2
    for provider, count in zip(proof['providers'], [4, 1]):
        assert len(provider['contracts']) == count
        assert provider['reviewed_landed_tree_equal'] is True
        assert provider['remote_observation']['remote_main'] == provider['landed_commit']
        assert all(c['conclusion'] in ['success', 'skipped'] for c in provider['remote_observation']['checks'])
        root = repos.get(provider['repository'])
        if root:
            changed = subprocess.check_output(['git', '-C', root, 'diff', '--name-only', provider['prior_base'], provider['landed_commit']]).decode().splitlines()
            assert sorted(changed) == sorted(provider['changed_paths']) and len(changed) == 7
            def entries(commit):
                return {line.split('\t')[1]: line.split('\t')[0] for line in subprocess.check_output(['git', '-C', root, 'ls-tree', '-r', commit]).decode().splitlines()}
            old, new = entries(provider['prior_base']), entries(provider['landed_commit'])
            unchanged = [path for path in old if path not in changed]
            assert len(unchanged) == provider['unchanged_prior_entries']
            assert all(new[path] == old[path] for path in unchanged)
            for commit in [provider['reviewed_commit'], provider['landed_commit']]:
                actual = subprocess.check_output(['git', '-C', root, 'rev-parse', commit + '^{tree}']).decode().strip()
                assert actual == provider['tree'], commit
            raw = blob(provider['repository'], provider['landed_commit'], provider['attachment']['path'], provider['attachment']['sha256'])
            attachment = json.loads(raw)
            if count == 4:
                assert [{key: contract[key] for key in ['source', 'decision', 'bridge']} for contract in attachment['contracts']] == provider['contracts']
                assert all(contract['target'] == provider['shared_target'] for contract in attachment['contracts'])
                assert all(contract['policy'] == provider['scope_policy'] for contract in attachment['contracts'])
            else:
                assert attachment['contracts'] == provider['contracts']
            env = json.loads(blob(provider['repository'], provider['landed_commit'], provider['envelope']['path'], provider['envelope']['sha256']))
            assert {'path': env['path'], 'sha256': env['sha256']} == provider['attachment']
        # Native dataset path is decoded/virtual; verify its committed ordered chunks only.
        for chunk in provider['original_native_dataset']['chunks']:
            blob(provider['repository'], provider['landed_commit'], chunk['path'], chunk['sha256'])
        for pin in provider['unchanged_model_native_binding_pins']:
            raw = blob(provider['repository'], provider['landed_commit'], pin['path'], pin['sha256'])
            if raw is not None and pin['path'] == 'source/artifact-snapshot.json':
                assert json.loads(raw)['sqlite'] == provider['original_native_dataset']
        for contract in provider['contracts']:
            references(contract['source'], provider['repository'], provider['landed_commit'])
            references(contract['decision'], provider['repository'], provider['landed_commit'])
            if 'bridge' in contract:
                references(contract['bridge'], provider['repository'], provider['landed_commit'])
        if 'original_user_declaration' in provider:
            references(provider['original_user_declaration'], provider['repository'], provider['landed_commit'])
    references(proof['prior_wrapper_continuity'], None, None)
    records = list((BASE / 'implementation/$records').glob('w1-*.json'))
    assert len(records) == 5
    for path in records:
        state = json.loads(path.read_text())
        assert state['status'] == 'implementing'
        assert state['publication_evidence_ids'] == [] and state['query_evidence_ids'] == []
    for path in list(HERE.glob('*.json')) + list(HERE.glob('*.md')) + [BASE / 'evidence/$records/w1-provider-attachment-continuity.json']:
        assert not re.search(rb'/(?:private/tmp|Users/[^/]+)/', path.read_bytes()), path.name
    print(f'checked scope, review-hash declarations and 5 unchanged implementation states; rehashed {originals} original receipts and {objects} immutable Git blobs; no released-harness execution')


if __name__ == '__main__':
    main()
