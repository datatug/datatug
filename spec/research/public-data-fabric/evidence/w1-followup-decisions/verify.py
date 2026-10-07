#!/usr/bin/env python3
"""Rehash the curated decision packet and preserved ledger; no fetch or admission."""
import argparse
import sys
import hashlib
import json
import re
import subprocess
from pathlib import Path

HERE = Path(__file__).resolve().parent
BASE = HERE.parents[1]
ROOT = HERE.parents[4]
sys.path.insert(0, str(BASE))
from verify_no_copy_decision import approved_evidence_additions, approved_record_additions, preserved_bytes


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def read(name):
    return json.loads((HERE / name).read_text())


def normalize(text):
    """Only replace execution-local paths with non-authoritative archive labels."""
    text = re.sub(r'/Users/alex/projects/(?:\.worktrees/[^/]+/github\.com/)?([a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+)', r'https://github.com/\1', text)
    text = re.sub(r'/private/tmp/([a-zA-Z0-9_.-]+)(/[^\s\"\'`<>),;]+)?', lambda m: 'artifact:' + m.group(1) + (m.group(2) or ''), text)
    return re.sub(r'/Users/[^\s\"\'`<>),;]+', lambda m: 'private-reference:' + Path(m.group(0)).name, text)


def omit_samples(value):
    if isinstance(value, dict):
        return {k: omit_samples(v) for k, v in value.items() if k != 'sample_rows'}
    if isinstance(value, list):
        return [omit_samples(v) for v in value]
    return value


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--original-index', type=Path, help='optional private label-to-original-path mapping; rehash originals and reproduce transcriptions')
    parser.add_argument('--repo', action='append', default=[], metavar='GITHUB_URL=LOCAL_CHECKOUT', help='optional immutable source-object verification; repeat per repository')
    args = parser.parse_args()
    manifest = read('source-manifest.json')
    proof = read('evidence.json')
    decision = read('semantic-decision.json')
    documents = read('reviewed-records.json')['documents']
    declarations = manifest['original_document_hash_declarations']
    assert set(documents) == set(declarations)
    for name, expected in manifest['packet_files'].items():
        assert digest((HERE / name).read_bytes()) == expected, name
    for label, declaration in declarations.items():
        content = documents[label]['content'].encode()
        assert digest(content) == declaration['public_hash'], label
        assert len(content) == declaration['public_bytes'], label
        assert re.fullmatch('[0-9a-f]{64}', declaration['source_hash']), label
        assert declaration['source_bytes'] > 0, label
        if documents[label]['transformation'] == 'exact UTF-8 bytes':
            assert digest(content) == declaration['source_hash'], label

    def document(label):
        return json.loads(documents[label]['content'])

    # Associate root acceptance and independent review with original authored bytes.
    checked_refs = 0
    artifact_labels = {entry['artifact_label']: label for label, entry in declarations.items() if 'artifact_label' in entry}
    for scope, acceptance_label in [('semantic', 'root-acceptance-r1.json'), ('graph', 'root-acceptance-r2.json'), ('fx', 'root-acceptance-r2.json')]:
        for path, expected in document(scope + '/' + acceptance_label)['inputs'].items():
            label = artifact_labels[path]
            assert declarations[label]['source_hash'] == expected, label
            checked_refs += 1
    selection = document('selection/root-technical-approval-r2.json')
    for field, prefix in [('proposal_sha256', 'selection'), ('review_sha256', 'selection-review')]:
        for name, expected in selection[field].items():
            assert declarations[prefix + '/' + name]['source_hash'] == expected, name
            checked_refs += 1
    for scope, names in [('semantic-review', ['review-r1.json']), ('graph-review', ['review-r2.json']), ('fx-review', ['review-r2.json'])]:
        review = document(scope + '/' + names[0])
        counts = review.get('counts', review.get('verdict', {}))
        assert all(counts[key] == 0 for key in ['blockers', 'majors', 'minors']), scope
    authored = document('semantic/proposal.json')
    assert digest((HERE / 'semantic-proposal.md').read_bytes()) == declarations['semantic/proposal.md']['source_hash']
    assert decision['authored_proposal']['sha256'] == declarations['semantic/proposal.md']['source_hash']
    assert decision['authored_proposal']['source_json_hash'] == declarations['semantic/proposal.json']['source_hash']
    for field in ['projections', 'graph', 'bounds', 'authority_partition', 'required_tests']:
        assert decision[field] == authored[field], field
    assert not decision['execution_admitted'] and not decision['publication_admitted']
    assert document('graph-review/review-r1.json')['counts']['majors'] == 3
    assert set(proof['graph_architecture']['resolution']) == {'M1', 'M2', 'M3'}
    closure = proof['graph_architecture']['metadata_measurement']
    assert closure['measured_total'] == closure['known_bytes'] + closure['provisional_bytes'] == 516371
    assert closure['final_closure_complete'] is False and len(closure['holes']) == 8
    assert proof['graph_architecture']['limits']['metadata_aggregate_bytes'] == 2097152
    assert proof['selection']['native_physical_tables'] == 16
    assert proof['selection']['selected_public_physical_tables'] == 11
    assert proof['selection']['logical_public_recordsets'] == 5
    fx = proof['fx_readiness']
    assert fx['canonical_custom_terms_publication_gate']['status'] == 'open'
    assert not any(fx[key] for key in ['selected', 'implemented', 'published', 'new_semantic_acceptance', 'raw_downloads_republished'])
    correction = document('fx/repair-r2a-evidence.json')
    assert correction['independent_review_r1_sha256'] == declarations['fx-review/review-r1.md']['source_hash']
    app = proof['app_component']
    assert app['reviewed_head'] == app['receipt_head'] == '65ebc741fb5f4ec21a8d9b436a3a5e6998759fa4'
    assert app['merge_sha'] == app['main_ci_head'] == 'fe89a041b99037d9d651d37aec4117da63b0d6d0'
    assert app['reviewed_tree'] == app['merge_tree'] == 'a8301cc5f76469ffc721898b2f3a8f94d5ea5cfe'
    assert app['main_ci_status'] == 'passed' and len(app['main_checks']) == 5
    assert all(check['conclusion'] == 'success' for check in app['main_checks'])
    assert app['branch_deleted'] and app['worktree_absent'] and not app['task5_full_live_acceptance']
    assert app['reviewed_head'] in documents['app-review/root-reviewed-head-r2.md']['content']
    assert not proof['programme']['w1_complete'] and proof['programme']['b1_open']
    assert not any(proof['programme'][key] for key in ['phase4_passed', 'phase5_passed', 'phase6_passed'])

    # Preserve every old assertion, gate, implementation, question and evidence byte.
    for name, expected in manifest['preserved_research_files'].items():
        assert digest(preserved_bytes(ROOT, name, expected)) == expected, name
    for key, original in manifest['candidate_originals'].items():
        current = json.loads((BASE / 'candidates/$records' / (key + '.json')).read_text())
        assert set(current) == set(original), key
        assert current['reason'].startswith(original['reason']), key
        assert current['evidence_ids'] == original['evidence_ids'] + ['w1-fx-priority-readiness'], key
        for field in set(original) - set(manifest['allowed_candidate_changes']):
            assert current[field] == original[field], (key, field)
    # Explicit later additions do not alter this packet's original records.
    # Original bytes remain checked above; the BigQuery packet separately
    # verifies its additions and the complete original candidate/assertion set.
    later_records = {
        'candidates': {
            'bq-source-public-data', 'bq-data-citibike', 'bq-data-noaa-gsod',
            'bq-data-nyc-tlc', 'bq-data-noaa-icoads', 'bq-data-storm-events',
            'bq-smoke-sample-shakespeare', 'bq-opportunity-dual-runtime',
        },
        'evidence': {
            'w1-p1-wire-semantic-addendum', 'bigquery-architecture-review-r1',
            'bigquery-discovery-r1', 'bigquery-dual-runtime-r3',
            'w1-fx-semantic-plan', 'w1-transport-app-observation',
        },
        'questions': {'bigquery-dual-runtime-prerequisites'},
    }
    later_records['evidence'].update(approved_evidence_additions(ROOT))
    for table, additions in approved_record_additions(ROOT).items():
        later_records.setdefault(table, set()).update(additions)
    for table, count in manifest['baseline_record_counts'].items():
        records = BASE / table / '$records'
        additions = later_records.get(table, set())
        assert all((records / (key + '.json')).is_file() for key in additions), table
        actual = len(list(records.glob('*.json')))
        expected = count + len(additions) + (len(manifest['new_evidence_ids']) if table == 'evidence' else 0)
        assert actual == expected, table
    assert manifest['baseline_record_counts']['candidates'] == 691
    assert manifest['baseline_record_counts']['assertions'] == 1177
    for path in (BASE / 'implementation/$records').glob('w1-*.json'):
        row = json.loads(path.read_text())
        assert row['status'] == 'implementing'
        assert row['publication_evidence_ids'] == row['query_evidence_ids'] == []
    for path in HERE.glob('*'):
        if path.suffix in {'.md', '.json'}:
            assert not re.search(rb'/(?:private/tmp|Users/[^/]+)/', path.read_bytes()), path.name

    originals = 0
    if args.original_index:
        index = json.loads(args.original_index.read_text())
        for label, declaration in declarations.items():
            raw = Path(index[label]).read_bytes()
            assert digest(raw) == declaration['source_hash'] and len(raw) == declaration['source_bytes'], label
            content = normalize(raw.decode())
            if 'retained_fields' in declaration:
                original = json.loads(raw)
                content = normalize(json.dumps({k: original[k] for k in declaration['retained_fields']}, indent=2, ensure_ascii=False) + '\n')
            elif label == 'fx/probe.json':
                content = json.dumps(omit_samples(json.loads(content)), indent=2, ensure_ascii=False) + '\n'
            assert content == documents[label]['content'], label
            originals += 1
        for label, declaration in manifest['unrepublished_artifact_hash_declarations'].items():
            raw = Path(index[label]).read_bytes()
            assert digest(raw) == declaration['source_hash'] and len(raw) == declaration['source_bytes'], label
            originals += 1
    repositories = dict(value.split('=', 1) for value in args.repo)
    objects = 0
    for obj in manifest['immutable_git_objects']:
        repo = repositories.get(obj['repository'])
        if repo:
            raw = subprocess.check_output(['git', '-C', repo, 'show', obj['revision'] + ':' + obj['path']])
            assert digest(raw) == obj['sha256'], obj
            if 'bytes' in obj:
                assert len(raw) == obj['bytes'], obj
            objects += 1
    print(json.dumps({'packet_hashes': len(manifest['packet_files']), 'source_transcriptions': len(documents), 'root_input_hash_associations': checked_refs, 'preserved_files': len(manifest['preserved_research_files']), 'candidates': 691, 'assertions': 1177, 'private_originals_rehashed': originals, 'immutable_git_objects_rehashed': objects, 'w1_complete': False}))


if __name__ == '__main__':
    main()
