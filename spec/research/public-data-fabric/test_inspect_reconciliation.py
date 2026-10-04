#!/usr/bin/env python3
"""Counterexamples for measured FK/country evidence, not semantic verdicts."""
import hashlib
import json
from pathlib import Path
import sqlite3
import tempfile
import unittest
import inspect_reconciliation as probe


class EvidenceCounterexamples(unittest.TestCase):
    def test_composite_orphan_null_duplicate_and_exact_country(self):
        with tempfile.TemporaryDirectory() as tmp:
            base = Path(tmp)
            for name in ['evidence', 'candidates/$records', 'assertions/$records']:
                (base / name).mkdir(parents=True)
            db = base / 'fixture.sqlite'
            conn = sqlite3.connect(db)
            conn.executescript('CREATE TABLE target(a TEXT,b TEXT);'
                              'CREATE TABLE source(a TEXT,b TEXT,country TEXT);'
                              "INSERT INTO target VALUES('01','x'),('01','x'),('01','y');"
                              "INSERT INTO source VALUES('01','x','United States'),"
                              "('01','z','United States '),(NULL,'x','032'),('01',NULL,'32');")
            conn.close()
            def save(path, value):
                (base / path).write_text(json.dumps(value))
            save('evidence/demo-demand.json', {'corpora': {'test': {
                'input': {'sqlite_sha256': hashlib.sha256(db.read_bytes()).hexdigest()},
                'recordsets': {}}}})
            for name in ['source', 'target']:
                save('candidates/$records/' + name + '.json',
                     {'dataset_id': 'test', 'title': 'test.' + name})
            save('assertions/$records/demand-edge-test.json', {
                'kind': 'edge', 'from_id': 'source', 'to_id': 'target',
                'from_fields': ['a', 'b'], 'to_fields': ['a', 'b']})
            save('assertions/$records/demand-map-test.json', {
                'kind': 'mapping', 'from_id': 'source', 'to_id': 'demand-meaning-country',
                'from_fields': ['country']})
            prior = probe.BASE
            try:
                probe.BASE = base
                result = probe.inspect({'test': str(db)}, [{
                    'id': 'us', 'labels': {'en': 'United States'}, 'codes': {'alpha2': 'US'}}])
            finally:
                probe.BASE = prior
            edge = result['local_references']['demand-edge-test']
            self.assertEqual(2, edge['nonnull_source_rows'])
            self.assertEqual(2, edge['null_source_rows'])
            self.assertEqual(1, edge['unmatched_nonnull_rows'])
            self.assertEqual(1, edge['target_duplicate_key_groups'])
            country = result['country_representations']['demand-map-test']
            self.assertEqual(1, country['matched_rows'])
            self.assertEqual({'United States ', '032', '32'},
                             {r['value'] for r in country['unmatched']})


if __name__ == '__main__':
    unittest.main()
