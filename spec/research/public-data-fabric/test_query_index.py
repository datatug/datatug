#!/usr/bin/env python3
"""Small counterexample fixtures for graph provenance and truthful research gaps."""
from datetime import datetime, timezone
import unittest
import sys
sys.dont_write_bytecode = True
from query_index import query


class ResearchQueries(unittest.TestCase):
    def setUp(self):
        def row(kind, **kw):
            return dict(kind=kind, title='Fixture', status='promising', evidence_ids=['e'], **kw)
        self.c = {'demo': row('dataset'), 'demo-r': row('recordset', dataset_id='demo'),
                  'source': row('source', freshness='needs-review', observed_at='2026-01-01T00:00:00Z'),
                  'dataset': row('dataset', source_id='source', research_depth='shallow'),
                  'r': row('recordset', dataset_id='dataset'), 'm': row('meaning'),
                  'unlinked': row('meaning'),
                  'scenario': row('scenario', member_ids=['r'], demo_ids=['demo'])}
        self.a = {'map': dict(kind='mapping', from_id='r', to_id='m', meaning_id='m',
                              status='researched', compatibility='unresolved', evidence_ids=['e']),
                  'join': dict(kind='edge', from_id='demo-r', to_id='r',
                               status='researched', compatibility='unresolved', evidence_ids=['e'])}
        self.e = {'e': dict(outcome='observation')}
        self.coverage = dict(catalogue_demo_ids=['demo'], rows=[dict(
            meaning_id='demand', structural_demo_count=4, demo_ids=['demo'], observed_demo_ids=[],
            cells={'demo': [dict(has_observed_nonnull_values=False)]})])

    def run_query(self):
        return query(dict(candidates=self.c, assertions=self.a, evidence=self.e), self.coverage,
                     datetime(2026, 10, 4, tzinfo=timezone.utc), ('m', 'Geo.Country'))

    def test_graph_and_unresolved_are_not_acceptance(self):
        r = self.run_query()
        self.assertEqual(r['meaning_datasets']['m']['accepted_datasets'], [])
        self.assertEqual(r['meaning_datasets']['m']['associations'][0]['datasets'], ['dataset'])
        self.assertFalse(r['meaning_datasets']['Geo.Country']['candidate_exists'])
        self.assertEqual(r['scenario_dependencies']['scenario']['nodes'], ['dataset', 'r', 'scenario', 'source'])
        self.assertEqual({x['route'] for x in r['demo_enrichments']['demo']},
                         {'demo-recordset assertion', 'qualified scenario proposal'})
        self.assertEqual(r['dataset_mapping_gaps'][0]['classification'], 'unresolved-only')
        self.assertEqual(r['four_plus_demo_demand'][0]['cells']['demo']['observed_nonnull_cells'], 0)
        gap = next(x for x in r['meaning_dataset_gaps'] if x['meaning_id'] == 'unlinked')
        self.assertEqual((gap['linkage'], gap['suitability']),
                         ('no-explicit-public-dataset-link', 'undetermined'))

    def test_missing_schedule_inheritance_and_accepted_receipt(self):
        r = self.run_query()
        self.assertIn('r', r['freshness']['groups']['unscheduled'])
        self.assertIn('r', r['freshness']['groups']['stale_or_needs_review'])
        self.assertEqual(r['freshness']['records']['r']['fields']['freshness']['from_id'], 'source')
        self.assertEqual(r['research_depth']['promising_shallow_ids'], ['dataset'])
        self.c['dataset']['revalidate_after'] = '2026-01-02T00:00:00Z'
        self.a['map'].update(status='reconciled', compatibility='exact')
        # Lifecycle alone cannot make a mapping accepted.
        self.assertEqual(self.run_query()['meaning_datasets']['m']['accepted_datasets'], [])
        self.e['e'].update(outcome='pass', independent=True, review_stage='reconciliation')
        r = self.run_query()
        self.assertEqual(r['meaning_datasets']['m']['accepted_datasets'], ['dataset'])
        self.assertEqual(r['dataset_mapping_gaps'][0]['classification'], 'accepted-reconciliation-backed')
        self.assertIn('r', r['freshness']['groups']['due'])
        self.assertNotIn('r', r['freshness']['groups']['unscheduled'])


if __name__ == '__main__':
    unittest.main()
