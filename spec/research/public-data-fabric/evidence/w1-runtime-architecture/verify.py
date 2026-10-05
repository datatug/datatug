#!/usr/bin/env python3
"""Bounded offline integrity/preservation checks for accepted architecture evidence."""
import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

HERE = Path(__file__).resolve().parent
BASE = HERE.parents[1]
HUB = BASE.parents[2]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def git(*args):
    return subprocess.check_output(['git', '-C', str(HUB), *args])


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', action='append', default=[], metavar='OWNER/REPOSITORY=LOCAL_CHECKOUT')
    parser.add_argument('--original-root', action='append', default=[], metavar='architecture|review|identity=DIRECTORY')
    parser.add_argument('--original-file', action='append', default=[], metavar='mission|approved-plan=FILE')
    args = parser.parse_args()
    repos = dict(x.split('=', 1) for x in args.repo)
    roots = {k: Path(v) for k, v in (x.split('=', 1) for x in args.original_root)}
    files = {k: Path(v) for k, v in (x.split('=', 1) for x in args.original_file)}
    proof = json.loads((HERE / 'evidence.json').read_text())
    assert proof['kind'] == 'w1-accepted-runtime-architecture-evidence'
    assert proof['acceptance']['status'] == 'accepted-for-implementation'
    assert proof['acceptance']['packet_review'] == 'pending separate adversarial PR review'
    assert [r['verdict']['majors'] for r in proof['review_history']] == [3, 1, 0]
    assert proof['review_history'][-1]['verdict']['land'] is True
    assert [r['id'] for r in proof['accepted_finding_resolutions']] == ['M1', 'M2', 'M3', 'M4']
    assert proof['preservation']['b1'] == 'open' and proof['preservation']['w1_complete'] is False
    assert proof['architecture']['capacity_gate']['status'] == 'not measured'
    assert proof['architecture']['publication']['fresh_geo_r3_promoted'] is False
    originals = {}
    for label, expected in proof['original_documents'].items():
        assert re.fullmatch('[0-9a-f]{64}', expected['sha256']) and expected['bytes'] > 0
        parts = label.split('/', 1)
        path = roots[parts[0]] / parts[1] if len(parts) == 2 and parts[0] in roots else files.get(label)
        if path:
            raw = path.read_bytes()
            assert digest(raw) == expected['sha256'] and len(raw) == expected['bytes'], label
            originals[label] = raw
    final_label = 'architecture/proposal-r3.json'
    if final_label in originals:
        final = json.loads(originals[final_label])
        normalized = {k: final[k] for k in proof['architecture']}
        normalized['publication']['public_identity']['status'] = 'accepted architecture allocation only; live publication proof required'
        assert normalized == proof['architecture']
    approval_label = 'architecture/root-technical-approval-r3.json'
    if approval_label in originals:
        approval = json.loads(originals[approval_label])
        for obj in approval['artifacts']:
            label = ('review/' if obj['name'].startswith('review') else 'architecture/') + obj['name']
            assert proof['original_documents'][label]['sha256'] == obj['sha256']
    checked = 0
    for obj in proof['frozen_source_observations']['immutable_objects']:
        assert re.fullmatch('[0-9a-f]{40}', obj['revision'])
        assert re.fullmatch('[0-9a-f]{64}', obj['sha256']) and obj['bytes'] > 0
        root = repos.get(obj['repository'])
        if root:
            entry = subprocess.check_output(['git', '-C', root, 'ls-tree', obj['revision'], '--', obj['path']]).split()
            assert entry[0] in [b'100644', b'100755']
            if 'git_blob' in obj:
                assert entry[2].decode() == obj['git_blob']
            raw = subprocess.check_output(['git', '-C', root, 'show', f"{obj['revision']}:{obj['path']}"])
            assert digest(raw) == obj['sha256'] and len(raw) == obj['bytes'], obj
            checked += 1
    prior = git('ls-tree', '-r', proof['hub_base'])
    assert digest(prior) == proof['preservation']['prior_git_entries_sha256']
    assert len(prior.decode().splitlines()) == proof['preservation']['prior_git_entries']
    def entries(raw):
        return {x.split('\t')[1]: x.split('\t')[0] for x in raw.decode().splitlines()}
    old, new = entries(prior), entries(git('ls-tree', '-r', 'HEAD'))
    assert all(new.get(path) == identity for path, identity in old.items())
    # Include unstaged/staged changes so an uncommitted prior-file edit cannot pass.
    changed = git('diff', '--name-only', proof['hub_base']).decode().splitlines()
    assert not set(changed).intersection(old), changed
    records = list((BASE / 'implementation/$records').glob('w1-*.json'))
    assert len(records) == 5
    for path in records:
        obj = json.loads(path.read_text())
        assert obj['status'] == 'implementing'
        assert obj['publication_evidence_ids'] == [] and obj['query_evidence_ids'] == []
    assert json.loads((BASE / 'questions/$records/corpus-publication.json').read_text())['status'] == 'open'
    for name in ['phase4', 'phase5', 'phase6']:
        assert json.loads((BASE / f'gates/$records/{name}.json').read_text())['status'] == 'pending'
    for path in list(HERE.glob('*.json')) + list(HERE.glob('*.md')) + [BASE / 'evidence/$records/w1-runtime-architecture.json']:
        assert not re.search(rb'/(?:private/tmp|Users/[^/]+)/', path.read_bytes()), path.name
    print(f"preserved {len(old)} prior Git entries, research and 5 implementing records/B1/gates; rehashed {len(originals)} originals and {checked} immutable source blobs; architecture only")


if __name__ == '__main__':
    main()
