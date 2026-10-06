"""Check frozen receipt provenance and immutable research; optional private originals."""
import argparse
import hashlib
import json
from pathlib import Path
import re


def require(condition, message):
    if not condition:
        raise ValueError(message)


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def check_ci(observation):
    if observation.get('status') != 'passed' or 'checks' not in observation:
        return
    head = observation.get('head')
    require(head and observation.get('observed_head') == head, 'CI observed head mismatch')
    require(observation.get('observed_target_head') == head, 'CI target head mismatch')
    checks = observation['checks']
    require(checks, 'Empty successful CI observation')
    require(all(c.get('conclusion') in ('success', 'skipped', 'neutral') for c in checks),
            'Successful CI contains failed or pending check')


def controls():
    for bad in [
        {'status': 'passed', 'head': 'a', 'observed_head': 'b', 'checks': [{}]},
        {'status': 'passed', 'head': 'a', 'observed_head': 'a', 'observed_target_head': 'b', 'checks': [{}]},
        {'status': 'passed', 'head': 'a', 'observed_head': 'a', 'observed_target_head': 'a', 'checks': [{'conclusion': 'failure'}]},
        {'status': 'passed', 'head': 'a', 'observed_head': 'a', 'observed_target_head': 'a', 'checks': []},
    ]:
        try:
            check_ci(bad)
        except ValueError:
            continue
        raise ValueError('Provenance negative control accepted')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--original-map', type=Path)
    args = parser.parse_args()
    packet = Path(__file__).resolve().parent
    root = packet.parents[4]
    manifest = json.loads((packet / 'source-manifest.json').read_text())
    receipts = json.loads((packet / 'publication-receipts.json').read_text())
    evidence = json.loads((packet / 'evidence.json').read_text())
    for name, sha in manifest['packet_files'].items():
        require(digest(packet / name) == sha, f'Packet hash mismatch: {name}')
    for name, sha in manifest['preserved_files'].items():
        require(digest(root / name) == sha, f'Preserved file changed: {name}')
    counts = {g: sum(f'/{g}/$records/' in n and n.endswith('.json') for n in manifest['preserved_files'])
              for g in ['candidates', 'assertions', 'questions', 'implementation', 'gates']}
    require(counts == evidence['preserved_record_counts'] == manifest['preserved_records'], 'Record count drift')
    require(counts['candidates'] == 699 and counts['assertions'] == 1177, 'Candidate/assertion preservation drift')
    originals = json.loads(args.original_map.read_text()) if args.original_map else None
    by_id = {}
    for receipt in receipts['receipts']:
        label = receipt['id']
        require(label not in by_id, 'Duplicate receipt id')
        by_id[label] = receipt['observation']
        declaration = manifest['original_documents'][label]
        require(receipt['source_sha256'] == declaration['sha256'] and receipt['source_bytes'] == declaration['bytes'],
                f'Original declaration mismatch: {label}')
        observation = receipt['observation']
        check_ci(observation)
        if originals:
            source = Path(originals[label])
            require(digest(source) == declaration['sha256'] and source.stat().st_size == declaration['bytes'],
                    f'Original source changed: {label}')
            try:
                raw = json.loads(source.read_text())
            except json.JSONDecodeError:
                raw = {}
            for key, value in observation.items():
                if key in ('premerge_checks', 'checks', 'workflow_runs', 'tag_object', 'module_origin', 'generator_commands', 'wrapper_commands'):
                    if key == 'generator_commands':
                        require(value == [{'name': c['name'], 'exit_code': c['result']['exit_code']} for c in raw], f'Generator association mismatch: {label}')
                    elif key == 'wrapper_commands':
                        require(value == [{k: c[k] for k in ['repository', 'head', 'exit_code']} for c in raw], f'Wrapper association mismatch: {label}')
                    elif key == 'premerge_checks':
                        for nested_key, nested_value in value.items():
                            if nested_key == 'checks':
                                require(nested_value == [{k: c[k] for k in ['name', 'bucket', 'conclusion', 'link', 'check_run_id'] if k in c} for c in raw['checks']['checks']], f'Premerge check association mismatch: {label}')
                            else:
                                require(raw['checks'].get(nested_key) == nested_value, f'Premerge field mismatch: {label}')
                    elif key == 'checks':
                        require(value == [{k: c[k] for k in ['name', 'bucket', 'conclusion', 'link', 'check_run_id'] if k in c}
                                          for c in raw['checks'] if isinstance(c, dict)], f'Check association mismatch: {label}')
                    elif key == 'workflow_runs':
                        require(value == [{k: c[k] for k in ['databaseId', 'headSha', 'conclusion', 'status', 'name'] if k in c}
                                          for c in raw if isinstance(c, dict) and 'headSha' in c], f'Workflow association mismatch: {label}')
                    elif key == 'tag_object':
                        require(value == {k: raw['object'][k] for k in ['sha', 'type'] if k in raw['object']}, f'Tag association mismatch: {label}')
                    elif key == 'module_origin':
                        require(value == {k: raw['Origin'][k] for k in ['VCS', 'URL', 'Hash', 'Ref'] if k in raw['Origin']}, f'Module association mismatch: {label}')
                    continue
                require(isinstance(raw, dict) and raw.get(key) == value, f'Explicit original field mismatch: {label}/{key}')
    require(set(by_id) == set(manifest['original_documents']), 'Source manifest receipt set mismatch')
    for prefix, head, version in [
        ('go-snapshot', '7df3d802f2c5188025c56fd8e41d27f01206a303', 'v0.1.1'),
        ('cli', '3a2d617f98e3762c0d3c26ce21ff69bae51dd11f', 'v0.60.0'),
    ]:
        release = by_id[prefix + '-release-r1']
        tag = by_id[prefix + '-release-tag-r1']
        require(release['isDraft'] is False and release['isPrerelease'] is False and release['publishedAt'], 'Unpublished release')
        require(release['tagName'] == version and tag['ref'] == 'refs/tags/' + version, 'Release/tag version mismatch')
        require(tag['tag_object'] == {'sha': head, 'type': 'commit'}, 'Release tag head mismatch')
        require(by_id[prefix + '-land-r1']['merge_sha'] == head, 'Release landing mismatch')
    require(by_id['go-snapshot-module-download-r1']['module_origin']['Hash'] == by_id['go-snapshot-release-tag-r1']['tag_object']['sha'],
            'Module/release head mismatch')
    require(all(c['exit_code'] == 0 for c in by_id['preparation-wrapper-generator-results-r3']['generator_commands']), 'Successor generator failed')
    require(all(c['exit_code'] == 0 for c in by_id['preparation-wrapper-check-results-r3']['wrapper_commands']), 'Successor wrapper failed')
    smoke = by_id['cli-release-smoke-r1']
    require(smoke['codesign']['sandbox_attempt_exit'] == 1 and smoke['codesign']['host_attempt_exit'] == 0, 'Signature observation history lost')
    require(smoke['checksum']['matches_published_checksum'] is True and smoke['version']['exit_code'] == 0 and smoke['credential_or_query_use'] is False, 'Binary surface proof drift')
    require(by_id['cli-main-ci-r1']['status'] == 'pending' and by_id['cli-main-ci-r2']['status'] == 'passed', 'Pending history lost')
    correction = evidence['correction']
    require(correction['historical_generator_exit'] == 1 and correction['successor_generator_exit'] == 0, 'Generator correction lost')
    require(evidence['programme'] == {'phase4': 'active', 'phase5': 'open', 'phase6': 'open', 'B1': 'open',
                                     'public_geo_ror_admission': False, 'both_go_and_js_required': True}, 'Programme boundary drift')
    for path in packet.iterdir():
        if path.is_file():
            content = path.read_text()
            require(not re.search(r'/(?:Users|private/tmp|home)/', content), f'Private path published: {path.name}')
    controls()
    print(f'PASS: {len(by_id)} original declarations, {len(manifest["preserved_files"])} unchanged research files; provenance controls; private originals={bool(originals)}')


if __name__ == '__main__':
    main()
