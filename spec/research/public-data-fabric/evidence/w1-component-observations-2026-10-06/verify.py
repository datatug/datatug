"""Offline frozen packet integrity, source associations and additive preservation."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import subprocess

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[4]


def require(ok, message):
    if not ok:
        raise ValueError(message)


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load(path):
    return json.loads(path.read_text())


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--original-map', type=Path)
    args = parser.parse_args()
    manifest = load(HERE / 'source-manifest.json')
    evidence = load(HERE / 'evidence.json')
    receipts = load(HERE / 'receipts.json')
    for name, expected in manifest['packet_files'].items():
        require(digest(HERE / name) == expected, f'Packet changed: {name}')
    require(set(manifest['packet_files']) == {'README.md', 'evidence.json', 'receipts.json', 'verify.py'}, 'Packet file set')
    base = manifest['hub_base']
    require(base == evidence['hub_base'], 'Base mismatch')
    require(re.fullmatch('[0-9a-f]{40}', base), 'Base identity')
    baseline = subprocess.check_output(['git', '-C', str(ROOT), 'ls-tree', '-r', '--name-only', base, 'spec/research/public-data-fabric']).decode().splitlines()
    require(set(baseline) == set(manifest['preserved_files']), 'Prior research inventory mismatch')
    objects = subprocess.check_output(['git', '-C', str(ROOT), 'cat-file', '--batch'], input=''.join(f'{base}:{name}\n' for name in baseline).encode())
    cursor = 0
    for name in baseline:
        end = objects.index(b'\n', cursor)
        header = objects[cursor:end].split()
        require(len(header) == 3 and header[1] == b'blob', f'Baseline object: {name}')
        size = int(header[2])
        original = objects[end + 1:end + 1 + size]
        cursor = end + 2 + size
        expected = manifest['preserved_files'][name]
        require(digest(ROOT / name) == expected, f'Prior research changed: {name}')
        require(hashlib.sha256(original).hexdigest() == expected, f'Baseline hash mismatch: {name}')
    counts = {g: sum(f'/{g}/$records/' in n and n.endswith('.json') for n in baseline)
              for g in ['candidates', 'assertions', 'questions', 'implementation', 'gates']}
    require(counts == manifest['preserved_records'] == {'candidates': 699, 'assertions': 1177, 'questions': 46, 'implementation': 5, 'gates': 8}, 'Preserved counts')
    require(evidence['programme'] == {'B1': 'open', 'phase4': 'active', 'phase5': 'open', 'phase6': 'open', 'w1_complete': False, 'cold_cloud_test': 'deferred by user', 'cold_cloud_capacity_accepted': False, 'public_geo_ror_admission': False, 'transport_accepted': False, 'memory_limit_bytes': 536870912, 'both_go_and_js_required': True}, 'Acceptance boundary drift')
    require(load(ROOT / 'spec/research/public-data-fabric/questions/$records/corpus-publication.json')['status'] == 'open', 'B1 source status')
    originals = load(args.original_map) if args.original_map else None
    if originals:
        require(set(originals) == set(manifest['original_documents']), 'Original map labels')
    for label, declared in manifest['original_documents'].items():
        require(re.fullmatch('[0-9a-f]{64}', declared['sha256']) and declared['bytes'] > 0, f'Invalid original declaration: {label}')
        if originals:
            path = Path(originals[label])
            require(digest(path) == declared['sha256'] and path.stat().st_size == declared['bytes'], f'Original mismatch: {label}')
    observations = {}
    for receipt in receipts['receipts']:
        label = receipt['id']
        require(label not in observations, 'Duplicate receipt')
        observations[label] = receipt['observation']
        d = manifest['original_documents'][label]
        require(receipt['source_sha256'] == d['sha256'] and receipt['source_bytes'] == d['bytes'], f'Receipt hash association: {label}')
        if originals:
            raw = load(Path(originals[label]))
            fields = receipt['fields']
            selected = raw if fields is None else {k: raw[k] for k in fields if k in raw}
            require(selected == receipt['observation'], f'Original field association: {label}')
    capacity = observations['capacity-report']
    require(capacity['source_head'] == capacity['expected_source_head'] == '889d722791e36c08226069002b2429883c0c3b4a', 'Capacity source identity')
    require(capacity['outcome'] == 'hosted_linux_candidate_pass' and capacity['actual_corpus'] is True, 'Actual candidate outcome')
    require(all(capacity[k] is False for k in ['cold_cloud_capacity_accepted', 'public_provider_admission', 'transport_accepted']), 'Candidate promotion')
    artifact = observations['capacity-artifact-metadata']['artifacts'][0]
    require(artifact['id'] == 11383528006 and artifact['workflow_run']['id'] == 37397632919 and artifact['workflow_run']['head_sha'] == capacity['source_head'], 'Artifact binding')
    require(artifact['digest'] == 'sha256:1fdb372096eba5a6dd19574a89376a68dc35300980cf3e1b7de01a6d23ca86e4', 'Archive digest')
    require(manifest['original_documents']['capacity-report']['sha256'] == '135048a2318ef29bdfd184621071afc8d7966c76f9c75970397c99846ba1ce68', 'Report digest')
    prep = receipts['capacity_preparation_projection']
    require(prep['native_tables'] == 16 and prep['selected_tables'] == 11 and prep['keys_scanned'] == 1280314, 'Complete corpus counts')
    workload = receipts['capacity_workload_projection']
    require(workload['final']['memory.max'] == 536870912 and workload['final']['memory.peak'] == 316518400, 'Retained memory observations')
    require(all(workload['final']['memory.events'][k] == 0 for k in ['max', 'oom', 'oom_kill']), 'Memory event observations')
    if originals:
        raw = load(Path(originals['capacity-report']))
        for field, section in [('capacity_workload_projection', 'workload'), ('capacity_preparation_projection', 'preparation')]:
            require(receipts[field] == {k: raw[section][k] for k in receipts[field]}, f'Capacity projection association: {field}')
    summary = observations['capacity-land-summary-r1']
    require(summary['head'] == capacity['source_head'] and summary['merge_sha'] == '775c40c237b8196971b70264c5b491dfaa96845b', 'Root landing summary identity')
    require(summary['postmain_ci']['head'] == summary['merge_sha'] and summary['cloud_run_test'] == 'deferred by user', 'Root postmain summary scope')
    cors = observations['cors-root-completion-r1']
    require(cors['reviewed_head'] == '87453befe2d57dd31251f45f7e002b1102bae375' and cors['merge'] == 'c10e6587b8519be6b1736bbda27bab09f6116ff7', 'CORS binding')
    ci = observations['cors-main-ci-r1']
    require(ci['head'] == ci['observed_head'] == ci['observed_target_head'] == cors['merge'], 'CORS exact-main CI')
    require(ci['status'] == 'passed' and ci['checks'] and all(c['conclusion'] in ['success', 'skipped', 'neutral'] for c in ci['checks']), 'CORS terminal checks')
    smoke = observations['cors-live-smoke-r1']
    require(smoke['source_head'] == cors['merge'] and len(smoke['results']) == 6, 'CORS smoke identity')
    require(sum(len(p['checks']) for p in smoke['results']) == 36 and all(p['outcome'] == 'passed' for p in smoke['results']), 'CORS bounded checks')
    for name, declared in cors['evidence_sha256'].items():
        labels = {'land-r1.json': 'cors-land-r1', 'main-ci-r1.json': 'cors-main-ci-r1', 'cloud-deploy-r1.log': 'cors-cloud-deploy-r1', 'worker-deploy-r1.log': 'cors-worker-deploy-r1', 'live-smoke-r1.json': 'cors-live-smoke-r1'}
        require(manifest['original_documents'][labels[name]]['sha256'] == declared, f'CORS root chain: {name}')
    mg_run = receipts['meaninggraph_deploy_projection']
    require(mg_run['id'] == evidence['meaninggraph_correction']['successful_deploy_run'] and mg_run['head_sha'] == evidence['meaninggraph_correction']['site_head'] and mg_run['conclusion'] == 'success', 'MG deploy observation')
    if originals:
        source_runs = load(Path(originals['b1-mg-site-runs.json']))['workflow_runs']
        original_run = next(r for r in source_runs if r['id'] == mg_run['id'])
        require(mg_run == {k: original_run[k] for k in mg_run}, 'MG run source association')
    matrix = observations['b1-provider-matrix.json']
    require(len(matrix) == 6 and sum(len(p['files']) for p in matrix) == 18, 'B1 matrix counts')
    require(all(p['registry_pins_agree'] and p['runtime_http_status'] == p['directory_http_status'] == p['meaninggraph_graph_http_status'] == '200' and all(f['live_equals_pinned'] for f in p['files']) for p in matrix), 'B1 bounded observations')
    require(evidence['reviews'][2]['verdict'] == {'blockers': 2, 'majors': 0, 'minors': 1, 'land': False}, 'B1 closure verdict retained')
    for path in HERE.iterdir():
        if path.is_file():
            require(not re.search(r'/(?:Users|private/tmp|home)/', path.read_text()), f'Private path published: {path.name}')
    print(json.dumps({'outcome': 'passed', 'preserved_files': len(baseline), 'preserved_records': counts, 'original_declarations': len(manifest['original_documents']), 'originals_rehashed': len(originals) if originals else 0, 'scope': 'offline packet integrity and additive preservation only'}, sort_keys=True))


if __name__ == '__main__':
    main()
