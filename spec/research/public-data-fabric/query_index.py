#!/usr/bin/env python3
"""Fixed read-only research queries. No equivalence inference or portfolio selection."""
import argparse
from datetime import datetime, timezone
import json
from pathlib import Path
import re

BASE = Path(__file__).resolve().parent
ACCEPTED = {'reconciled', 'selected', 'implementing', 'implemented', 'validated'}


def load(base):
    tables = {name: {p.stem: json.loads(p.read_text()) for p in
                    sorted((base / name / '$records').glob('*.json'))}
              for name in ('candidates', 'assertions', 'evidence')}
    coverage = json.loads((base / 'evidence/demo-demand-coverage.json').read_text())
    return tables, coverage


def query(tables, coverage, as_of, meaning_ids=('Geo.Country',)):
    c, assertions, evidence = (tables[k] for k in ('candidates', 'assertions', 'evidence'))
    demos = set(coverage['catalogue_demo_ids'])
    public = {k for k, v in c.items() if v['kind'] == 'dataset' and k not in demos}
    meanings = {k for k, v in c.items() if v['kind'] == 'meaning'}
    recordsets = {k for k, v in c.items() if v['kind'] == 'recordset'}

    def accepted(a):
        return (a.get('kind') == 'mapping' and a.get('status') in ACCEPTED and
                a.get('compatibility') not in (None, 'unresolved', 'incompatible') and
                any(evidence.get(e, {}).get('outcome') == 'pass' and
                    evidence.get(e, {}).get('independent') is True and
                    evidence.get(e, {}).get('review_stage') == 'reconciliation'
                    for e in a.get('evidence_ids', [])))

    def assertion_receipt(k, a):
        return {name: value for name, value in dict(id=k, **a).items()
                if name in ('id', 'kind', 'status', 'compatibility', 'evidence_ids',
                            'from_id', 'to_id', 'meaning_id', 'from_fields', 'to_fields')}

    def datasets_for(a):
        return sorted({c[r]['dataset_id'] for r in (a.get('from_id'), a.get('to_id'))
                       if r in recordsets and c[r]['dataset_id'] in public})

    associations = {m: [] for m in meanings}
    dataset_assertions = {d: [] for d in public}
    for k, a in sorted(assertions.items()):
        ds = datasets_for(a)
        for d in ds:
            dataset_assertions[d].append((k, a))
        m = a.get('meaning_id') or (a.get('to_id') if a.get('kind') == 'mapping' else None)
        if m in meanings and ds:
            associations[m].append(dict(datasets=ds, assertion=assertion_receipt(k, a),
                                        accepted=accepted(a)))

    def closure(start):
        seen, pending = set(), [start]
        while pending:
            k = pending.pop()
            if k in seen or k not in c:
                continue
            seen.add(k)
            row = c[k]
            pending.extend(row.get('member_ids', []) + row.get('dependency_ids', []))
            pending.extend(row[f] for f in ('dataset_id', 'source_id') if row.get(f))
        return sorted(seen)

    scenarios = {}
    for k, v in sorted(c.items()):
        if v['kind'] in ('scenario', 'bundle'):
            nodes = closure(k)
            scenarios[k] = dict(status=v['status'], demo_ids=v.get('demo_ids', []),
                                nodes=nodes, datasets=sorted(public.intersection(nodes)),
                                evidence_ids=v['evidence_ids'])

    enrichments = {d: [] for d in sorted(demos)}
    for dataset in sorted(public):
        for demo in sorted(set(c[dataset].get('demo_ids', [])).intersection(demos)):
            enrichments[demo].append(dict(route='direct dataset demand annotation', datasets=[dataset],
                                          status=c[dataset]['status'], evidence_ids=c[dataset]['evidence_ids']))
    for k, a in sorted(assertions.items()):
        ds = datasets_for(a)
        if not ds:
            continue
        origin = {c[r]['dataset_id'] for r in (a.get('from_id'), a.get('to_id'))
                  if r in recordsets and c[r]['dataset_id'] in demos}
        for demo in sorted(origin):
            enrichments[demo].append(dict(route='demo-recordset assertion', datasets=ds,
                                          assertion=assertion_receipt(k, a)))
    for k, s in scenarios.items():
        for demo in sorted(set(s['demo_ids']).intersection(demos)):
            if s['datasets']:
                enrichments[demo].append(dict(route='qualified scenario proposal',
                                              scenario=k, datasets=s['datasets'],
                                              status=s['status'], evidence_ids=s['evidence_ids']))

    prevalence = []
    for row in coverage['rows']:
        if row['structural_demo_count'] >= 4:
            prevalence.append({k: row[k] for k in ('meaning_id', 'structural_demo_count',
                                                    'demo_ids', 'observed_demo_ids')} |
                              {'cells': {d: {'structural_cells': len(cells),
                                             'observed_nonnull_cells': sum(bool(x.get('has_observed_nonnull_values'))
                                                                           for x in cells),
                                             'unobserved_cells': sum(not bool(x.get('has_observed_nonnull_values'))
                                                                    for x in cells)}
                                         for d, cells in sorted(row['cells'].items())}})

    def effective(k, field):
        row = c[k]
        order = [k]
        dataset = row.get('dataset_id')
        if dataset:
            order.append(dataset)
        source = row.get('source_id') or c.get(dataset, {}).get('source_id')
        if source:
            order.append(source)
        for origin in order:
            value = c.get(origin, {}).get(field)
            if value is not None:
                return {'value': value, 'from_id': origin, 'inherited': origin != k}
        return {'value': None, 'from_id': None, 'inherited': False}

    freshness = {group: [] for group in ('due', 'stale_or_needs_review',
                                        'unknown_or_missing', 'missing_observation', 'unscheduled')}
    freshness_records = {}
    for k, row in sorted(c.items()):
        if row['kind'] not in ('source', 'dataset', 'recordset'):
            continue
        fields = {f: effective(k, f) for f in ('freshness', 'observed_at', 'revalidate_after', 'update_frequency')}
        freshness_records[k] = dict(kind=row['kind'], demo_fixture=(k in demos or row.get('dataset_id') in demos), fields=fields)
        due = fields['revalidate_after']['value']
        if due and datetime.fromisoformat(due.replace('Z', '+00:00')) <= as_of:
            freshness['due'].append(k)
        if fields['freshness']['value'] in ('stale', 'needs-review'):
            freshness['stale_or_needs_review'].append(k)
        if fields['freshness']['value'] in (None, 'unknown'):
            freshness['unknown_or_missing'].append(k)
        if not fields['observed_at']['value']:
            freshness['missing_observation'].append(k)
        if not due:
            freshness['unscheduled'].append(k)

    gaps = []
    for m in sorted(meanings):
        links = associations[m]
        ds = sorted({d for link in links for d in link['datasets']})
        gaps.append(dict(meaning_id=m, datasets=ds,
                         linkage='explicit-public-dataset-candidates' if ds else 'no-explicit-public-dataset-link',
                         blocked_or_rejected=[d for d in ds if c[d]['status'] in ('blocked', 'rejected')],
                         accepted_mapping_recorded=any(x['accepted'] for x in links),
                         suitability='undetermined', evidence_ids=c[m]['evidence_ids']))

    mapping_gaps = []
    for d, links in sorted(dataset_assertions.items()):
        mappings = [(k, a) for k, a in links if a.get('kind') == 'mapping']
        state = ('accepted-reconciliation-backed' if any(accepted(a) for _, a in mappings) else
                 'no-mapping-recorded' if not mappings else
                 'unresolved-only' if all(a.get('compatibility') == 'unresolved' for _, a in mappings) else
                 'incompatible-or-rejected' if all(a.get('compatibility') == 'incompatible' or
                                                  a.get('status') == 'rejected' for _, a in mappings) else
                 'unreconciled-mixed')
        mapping_gaps.append(dict(dataset_id=d, classification=state,
                                 mapping_ids=[k for k, _ in mappings],
                                 edge_ids=[k for k, a in links if a.get('kind') == 'edge']))

    alpha2 = re.compile(r'alpha[-_ ]?2|\biso[-_ ]?2(?:code)?\b', re.I)
    representation = [dict(id=k, kind=v['kind'], status=v['status'],
                           representations=[x for x in v.get('identifier_representations', []) if alpha2.search(x)],
                           evidence_ids=v['evidence_ids'],
                           dataset_id=v.get('dataset_id') if v['kind'] == 'recordset' else None,
                           source_id=v.get('source_id') if v['kind'] in ('recordset', 'dataset') else None) for k, v in sorted(c.items())
                      if any(alpha2.search(x) for x in v.get('identifier_representations', []))]
    licensing = [dict(id=k, status=v['status'], licence=v.get('licence'),
                      licence_notes=v.get('licence_notes'), reason=v.get('reason'),
                      evidence_ids=v['evidence_ids']) for k, v in sorted(c.items())
                 if v['status'] in ('blocked', 'rejected') and
                 re.search(r'licen[cs]|permission|redistribut|copyright',
                           ' '.join(str(v.get(f, '')) for f in ('reason', 'licence', 'licence_notes')), re.I)]
    return dict(
        notice='Research query results only: no equivalence, suitable-dataset verdict or portfolio selection.',
        as_of=as_of.isoformat(), counts={name: len(rows) for name, rows in tables.items()},
        meaning_datasets={m: dict(candidate_exists=m in c, associations=associations.get(m, []),
                                  accepted_datasets=sorted({d for x in associations.get(m, []) if x['accepted']
                                                            for d in x['datasets']})) for m in meaning_ids},
        provisional_meaning_associations={m: links for m, links in sorted(associations.items()) if links},
        alpha2_representations=dict(policy='Literal case-insensitive alpha2/alpha-2/alpha_2/alpha 2 or whole ISO2/ISO-2/ISO_2/ISO 2, optionally followed by Code; no inferred code equivalence.', matches=representation,
                                    public_dataset_projection=sorted({x['dataset_id'] for x in representation
                                                                      if x['kind'] == 'recordset' and x['dataset_id'] in public}),
                                    scope='Searches candidate identifier_representations only, not assertions. Conceptual matches do not imply dataset usage; differently spelled fields such as GeoNames ISO remain outside this literal query.'),
        four_plus_demo_demand=prevalence, demo_enrichments=enrichments,
        scenario_dependencies=scenarios,
        scenario_dependents_by_dataset={d: sorted(k for k, v in scenarios.items() if d in v['datasets'])
                                        for d in sorted(public)}, licensing_decisions=licensing,
        freshness=dict(groups=freshness, records=freshness_records,
                       note='Inspection time is not upstream release/import freshness. Groups overlap; missing schedules do not mean no review is needed.'),
        research_depth=dict(counts={depth: sum(v.get('research_depth', 'unknown') == depth for v in c.values())
                                   for depth in ('shallow', 'documented', 'sampled', 'unknown')},
                            promising_shallow_ids=sorted(k for k, v in c.items() if v['status'] == 'promising' and v.get('research_depth') == 'shallow'),
                            shallow_ids=sorted(k for k, v in c.items() if v.get('research_depth') == 'shallow'),
                            unknown_or_missing_ids=sorted(k for k, v in c.items() if v.get('research_depth', 'unknown') == 'unknown')),
        meaning_dataset_gaps=gaps, dataset_mapping_gaps=mapping_gaps)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--base', type=Path, default=BASE)
    parser.add_argument('--as-of', required=True, help='Explicit ISO date/time; UTC if no offset')
    parser.add_argument('--meaning-id', action='append', help='Exact candidate ID; no alias reconciliation')
    args = parser.parse_args()
    moment = datetime.fromisoformat(args.as_of.replace('Z', '+00:00'))
    if moment.tzinfo is None:
        moment = moment.replace(tzinfo=timezone.utc)
    tables, coverage = load(args.base)
    print(json.dumps(query(tables, coverage, moment, args.meaning_id or ('Geo.Country',)),
                     indent=2, ensure_ascii=False, sort_keys=True))


if __name__ == '__main__':
    main()
