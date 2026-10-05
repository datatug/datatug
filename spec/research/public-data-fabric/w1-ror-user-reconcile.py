#!/usr/bin/env python3
"""Offline replay of the dedicated W1 decision; not runtime admission tooling."""
import argparse
from collections import Counter
import hashlib
import json
from pathlib import Path
import re
import sqlite3
import subprocess

APPS = 'e7362033ec79c6663d7dbe0483b62fab01f7b9cd'
PROVIDER = 'bbbec903248680caea04e68f94b9a957b6efc55b'
CORE = '982916d73f0a35ff2558b0062f58aa3ac4f24d97'
PREFIX = 'libs/datatug/main/src/lib/queries/fixtures/public-data-fabric/'
PINS = {'affiliations.modelspec.json': 'cfe958f12d9a70bb7533499b76bf01c9447f81f0862852e460b095bb75e1b9cf',
        'affiliations.json': '44a30dd260c74f43b2138934cca0bccd493c7b804e124c5ecbded3c41226e543'}
DB_HASH = 'afdf978130ee9899fedb698f0e9d18399efdea40ea92f01f26d6055ef1ed5d38'
ALPHABET = '0123456789abcdefghjkmnpqrstvwxyz'


def require(condition, message):
    if not condition:
        raise ValueError(message)


def blob(repo, revision, path):
    return subprocess.check_output(['git', '-C', str(repo), 'show', revision + ':' + path])


def sha(data):
    return hashlib.sha256(data).hexdigest()


def valid_url(value):
    if not isinstance(value, str) or not re.fullmatch(r'https://ror\.org/0[0-9a-hjkmnp-tv-z]{6}[0-9]{2}', value):
        return False
    token = value.rsplit('/', 1)[1]
    n = 0
    for char in token[1:7]:
        n = n * 32 + ALPHABET.index(char)
    return int(token[-2:]) == 98 - (n * 100) % 97


def lookup(db, value):
    if value is None:
        return {'state': 'null'}
    if value == '':
        return {'state': 'empty'}
    if not valid_url(value):
        return {'state': 'invalid'}
    rows = db.execute('SELECT id,status FROM organizations WHERE id = ? COLLATE BINARY', (value,)).fetchall()
    if len(rows) > 1:
        return {'state': 'ambiguous'}
    if not rows:
        return {'state': 'unmatched'}
    key, status = rows[0]
    return {'state': 'matched', 'native_id': key, 'status': status,
            'locations': db.execute('SELECT ordinal,geonames_id FROM locations WHERE organization_id = ? COLLATE BINARY ORDER BY ordinal', (key,)).fetchall()}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--apps-repo', type=Path, required=True)
    parser.add_argument('--provider-repo', type=Path, required=True)
    parser.add_argument('--sqlite', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    objects = {}
    for name, expected in PINS.items():
        data = blob(args.apps_repo, APPS, PREFIX + name)
        require(sha(data) == expected, 'fixture hash mismatch: ' + name)
        objects[name] = json.loads(data)
    schema = objects['affiliations.modelspec.json']
    require(schema['module']['name'] == 'affiliations' and schema['entities']['Affiliation']['properties']['ror_id'] == {'type': 'string', 'required': False}, 'user field contract changed')
    require(schema['entities']['Affiliation']['key'] == ['affiliation_id'], 'affiliation grain changed')
    files = {p: sha(blob(args.provider_repo, PROVIDER, p)) for p in ['model/ror.modelspec.json', 'model/ror.modelspec.hcl', 'model/ror.meaning.yaml', 'source/ror-v2.13.json', 'source/validation.json']}
    target = json.loads(blob(args.provider_repo, PROVIDER, 'model/ror.modelspec.json'))
    require(target['entities']['organizations']['key'] == ['id'] and target['entities']['organizations']['properties']['id'] == {'type': 'string', 'required': True}, 'target identity changed')
    with args.sqlite.open('rb') as stream:
        require(hashlib.file_digest(stream, 'sha256').hexdigest() == DB_HASH, 'SQLite hash mismatch')
    with sqlite3.connect(args.sqlite.resolve().as_uri() + '?mode=ro', uri=True) as db:
        require(db.execute('PRAGMA integrity_check').fetchone()[0] == 'ok', 'SQLite integrity failure')
        total, unique = db.execute('SELECT count(*), count(DISTINCT id) FROM organizations').fetchone()
        require(total == unique == 141528, 'target uniqueness/count changed')
        statuses = dict(db.execute('SELECT status,count(*) FROM organizations GROUP BY status'))
        require(statuses == {'active': 138289, 'inactive': 1809, 'withdrawn': 1430}, 'status counts changed')
        rows = objects['affiliations.json']
        require(len({r['affiliation_id'] for r in rows}) == len(rows) == 8, 'user grain/count changed')
        results = [{'affiliation_id': r['affiliation_id'], 'raw': r['ror_id'], **lookup(db, r['ror_id'])} for r in rows]
        observed = Counter(r['state'] for r in results)
        counts = {state: observed[state] for state in ['matched', 'null', 'empty', 'invalid', 'unmatched', 'ambiguous']}
        require(counts == {'matched': 3, 'null': 1, 'empty': 1, 'invalid': 3, 'unmatched': 0, 'ambiguous': 0}, 'fixture outcomes changed')
        controls = {}
        for value in ['000025p04', 'http://ror.org/000025p04', 'https://ror.org/000025p04/', 'https://ror.org/000025p05', 'https://ror.org/000025p04?x=1', 25]:
            require(lookup(db, value)['state'] == 'invalid', 'invalid representation accepted')
        controls['invalid_representations'] = 6
        absent = 'https://ror.org/000000098'
        require(valid_url(absent) and lookup(db, absent)['state'] == 'unmatched', 'valid unpublished control failed')
        controls['valid_but_absent'] = absent
        controls['status_examples'] = {}
        for status in statuses:
            key = db.execute('SELECT id FROM organizations WHERE status=? ORDER BY id LIMIT 1', (status,)).fetchone()[0]
            result = lookup(db, key)
            require(result['state'] == 'matched' and result['status'] == status and result['native_id'] == key, 'status suppressed or identity rewritten')
            controls['status_examples'][status] = result
        multi = db.execute('SELECT organization_id FROM locations GROUP BY organization_id HAVING count(*) > 1 ORDER BY organization_id LIMIT 1').fetchone()[0]
        details = lookup(db, multi)
        require(len(details['locations']) > 1, 'multiple locations lost')
        controls['multiple_location_example'] = details
        expanded, affiliations, organizations = db.execute('WITH inputs(affiliation_id,ror_id) AS (VALUES ("control-a",?),("control-b",?)) SELECT count(*),count(DISTINCT i.affiliation_id),count(DISTINCT o.id) FROM inputs i JOIN organizations o ON i.ror_id=o.id COLLATE BINARY JOIN locations l ON l.organization_id=o.id COLLATE BINARY', (multi,multi)).fetchone()
        require(expanded == 2 * len(details['locations']) and affiliations == 2 and organizations == 1, 'location expansion corrupted input grain')
        multiple_count = db.execute('SELECT count(*) FROM (SELECT organization_id FROM locations GROUP BY organization_id HAVING count(*) > 1)').fetchone()[0]
        require(multiple_count == 138, 'multiple-location source count changed')
        controls['repeated_affiliations_with_multiple_locations'] = {'input_affiliations': affiliations, 'distinct_organizations': organizations, 'flat_location_rows': expanded, 'locations_per_reference': len(details['locations']), 'warning': 'flat_location_rows is not an affiliation count'}
        controls['source_organizations_with_multiple_locations'] = multiple_count
        with sqlite3.connect(':memory:') as duplicate:
            duplicate.execute('CREATE TABLE organizations(id TEXT,status TEXT)')
            duplicate.executemany('INSERT INTO organizations VALUES (?,?)', [(multi, 'active'), (multi, 'active')])
            require(lookup(duplicate, multi)['state'] == 'ambiguous', 'duplicate target accepted')
        controls['duplicate_target'] = 'ambiguous; no target selected'
    decision = {
        'decision_id': 'w1-ror-user-native-url', 'disposition': 'accepted',
        'scope': 'Exact committed hypothetical fixture and field; native ROR URL snapshot-reference lookup only. Independent adversarial review and root landing remain required.',
        'source': {'repository': 'https://github.com/datatug/datatug-apps', 'revision': APPS,
            'schema_path': PREFIX + 'affiliations.modelspec.json', 'schema_sha256': PINS['affiliations.modelspec.json'],
            'data_path': PREFIX + 'affiliations.json', 'data_sha256': PINS['affiliations.json'],
            'declaration_path': PREFIX + 'README.md', 'declaration_sha256': sha(blob(args.apps_repo, APPS, PREFIX + 'README.md')),
            'module': 'affiliations', 'module_id': schema['module']['id'], 'entity': 'Affiliation', 'property': 'ror_id',
            'datatype': 'string', 'required': False, 'null_policy': 'Retain NULL as a separate no-reference exception; optional schema does not make NULL an identifier.',
            'identifier_namespace': 'ROR:URL', 'role': 'organization-reference value on an affiliation row; not affiliation identity or proof of affiliation'},
        'target': {'repository': 'https://github.com/ingitdb/ror-ingitdb', 'revision': PROVIDER,
            'schema_path': 'model/ror.modelspec.json', 'module': 'ror', 'module_id': target['module']['id'],
            'entity': 'organizations', 'property': 'id', 'datatype': 'string', 'identifier_namespace': 'ROR:URL',
            'role': 'native organization identifier', 'binding_path': 'model/ror.meaning.yaml', 'concept': 'research-organization',
            'core_parent': 'meaning://github.com/meaninggraph/core/organization?ref=' + CORE, 'files_sha256': files,
            'sqlite_sha256': DB_HASH, 'native_identity_distinct_from_serving_id': True},
        'policy': {'allowed_transform': 'identity', 'equality': 'UTF-8 byte-exact; case-and-space-sensitive',
            'validity': 'Exact https://ror.org/ prefix, lowercase Crockford nine-character token and MOD97-10 checksum; format/checksum never proves issued membership. Match only exact native IDs in the pinned full snapshot.',
            'per_input_cardinality': 'zero-or-one organization reference; duplicate target keys fail closed as ambiguous',
            'status': 'Retain active/inactive/withdrawn on every match. Inactive is historical reference context; withdrawn is an erroneous/duplicate/out-of-scope record reference and must be visibly flagged, not certified as a valid organization.',
            'time': 'Metadata and status describe ROR release v2.13 dated 2026-09-22 only. No user affiliation date exists; no current/historical affiliation truth or validity interval is inferred. Refresh or changed source/schema/revision requires a new reviewed pin.',
            'multiplicity': 'affiliation_id defines input grain; repeated ror_id values remain distinct affiliations. Locations are ordered child details, never independent affiliations or extra unique organizations.',
            'exceptions': ['null', 'empty', 'invalid', 'unmatched', 'ambiguous'],
            'no_implicit_transforms': ['trim', 'case-fold', 'bare-ID expansion', 'URL rewriting', 'successor substitution', 'name matching', 'serving-ID substitution'],
            'not_accepted': ['affiliation truth, employment, membership or temporal validity', 'legal-entity or LEI equivalence', 'native ROR fields in any of the six demos', 'new source fields/revisions by name similarity', 'production eligibility or runtime success']},
        'source_release': {'release': 'v2.13', 'date': '2026-09-22', 'archive_sha256': '476988542f1bdd185f779386f90149371f58cb4ebdba326db3b7b428a468d4a5', 'member_sha256': '6d032ff473f771e015da0837fe28b881a8f85822ee7815effefe11dc49c4346c'},
        'measurement': {'source_rows': 8, 'counts': counts, 'distinct_matched_organizations': len({r['native_id'] for r in results if r['state'] == 'matched'}), 'target_organizations': total, 'target_statuses': statuses, 'results': results},
        'counterexamples': controls,
        'authority': ['https://ror.readme.io/docs/identifier', 'https://ror.readme.io/docs/ror-data-structure', 'https://github.com/ror-community/ror-api/blob/bd040a0d2558a478c06a89118a29eeb9b6142710/rorapi/management/commands/generaterorid.py'],
        'open_dependencies': ['independent review and root landing of this decision', 'landed canonical provider model/meaning registrations and reviewed native-identity attachment', 'source-faithful native-id runtime and Directory publication', 'bounded app integration and independent real browser journey', 'B1 and programme phases4/5/6 remain open']}
    args.output.write_text(json.dumps(decision, indent=2, ensure_ascii=False) + '\n')
    print(json.dumps({'counts': counts, 'source_rows': 8, 'output_sha256': sha(args.output.read_bytes())}))


if __name__ == '__main__':
    main()
