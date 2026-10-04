# [review r1 datatug/datatug#86] Independent Phase 2 semantic review

Reviewed head: `0b514301e2faa434bfe855d55396cedc6a3291e8`; accepted input/base: `70791b1f05cfabe7815304ec810022c98776958a`. Reviewer: fresh local Codex on macOS, distinct from discovery and dedicated reconciliation. Review is bounded Phase 2 correctness/completeness, not portfolio, publication or DataTug execution acceptance.

## Findings

1. **Major — an external bridge is incorrectly required for already-supported local identifier roles.** `assertions/$records/demand-map-03f2614b38097665.json:26` leaves Chinook Employee.EmployeeId → employee-identifier unresolved because no canonical binding and external representation bridge exists. Its meaning candidate explicitly reuses core employee, while `demand-map-67e7d31645d6b0b1.json:25-26` already accepts the same source field's local identifier role. The native/ModelSpec key and canonical entity scope support that bounded local meaning; external identity is a different question. The meaning file binds Employee as an entity, not EmployeeId explicitly as identifier: canonical identifier-role publication remains a genuine implementation prerequisite. Resolve the duplicate decision explicitly and inspect the bounded employee/customer/order identifier pattern; do not turn absence of portable identity into a blanket local-semantic evidence gap. Preserve original hypotheses and distinguish duplicate mappings in counts. The author has identified 16 proposals for a bounded correction; those corrections are not accepted by this exact-head review yet.

2. **Minor — country match coverage omits null source rows.** `p2-reconciliation.md:39` says Northwind Customers 91/91. The pinned native Customers table contains 93 rows, including two NULL Country values. Independent reconstruction matches all 91 nonnull values, while `evidence/p2-bridge-probes.json` omits null accounting and `assertions/$records/d1-edge-northwind-country.json:15,29` says zero unmatched/ambiguous. State 91/91 nonnull, 91/93 total, two NULL; preserve those rows under an explicit lookup policy. This does not invalidate any reviewed nonnull label → ISO mapping.

3. **Minor — current role contradicts scoped acceptance.** Eight reconciled D1 edges retain `physical_role: Proposed contextual join; not executed or reconciled`, including `d1-edge-northwind-country.json:27`, four country bridges, ROR→GeoNames, places→admin1, alternate names→places and places→country. Replace current role with the actual scoped bridge/reference, retaining non-execution in DataTug and the old text as historical evidence.

## Independently verified evidence

- GitHub PR86 is OPEN/non-draft; remote head/base match above. Exact-head ledger run 37244595912 has conclusion SUCCESS.
- Fresh `git ls-remote` confirms every canonical pin below. Existing six collections remain authoritative; coverage is derived, with no second registry introduced.
- Compared all 681 original candidate and 1,162 assertion records to the accepted base: every ID, original reason/assumption substring and earlier evidence ID is retained. Counts independently reproduce 247 reconciled, 19 rejected, 904 unresolved assertions; 63 accepted, 5 rejected, 196 unresolved meanings; additions are one meaning/eight assertions.
- WB-governed independent checks pass: all six inGitDB collections; 24 mutation checks (original 22 plus two new meaning guards); three query fixtures; one composite/representation fixture; SpecScore zero violations. Validator changes are additive and preserve existing lifecycle/licensing/gate requirements. Optional semantic_disposition is restricted to meanings and requires a named independent reconciliation receipt.
- Replayed inspect_reconciliation against all six SHA256-pinned native SQLite fixtures with read-only immutable connections and independently extracted pinned core country values. Output exactly equals committed p2-native-checks.json: 152 local references, 23 country representations, 42 recordsets with alternate indexes. No nonnull orphan/duplicate target tuple; empty/all-null references are distinguished. Separate SQL confirms Pubs ord_num has two duplicate groups and native primary key (stor_id,ord_num,title_id). AdventureWorks state code/current composite both have zero duplicates; current uniqueness alone does not justify removing declared country scope.
- Independently reconstructed all five country maps using fresh GeoNames countryInfo and pinned core aliases: exact reviewed_values equality. Matched nonnull counts 59/91/830/8; Sakila 102 with the exact seven unresolved labels, including historical Yugoslavia. GeoNames has 252 unique ISO rows. Fresh response hash equals author receipt: 93bafc525813f22e4711ff9ed6d626343094ce48c26388dc7c49189b3d7d5512. Null accounting finding above remains.
- Fresh CLDR48 units response matches author SHA256 6ef31425565e76c456295ca9d2b140fed9da27bc67712228f0079dedbdf94614; pound is mass relative to kilogram with lb_to_kg. This supports the two narrow padded-code proposals, not arbitrary packaging/volume conversion.
- Canonical Chinook MeaningGraph check and ModelSpec publish lint independently pass. Inspected core country/customer/order/invoice/price contracts and the binding-role implementation, including exact codes and no generic trim pipeline. Full 18-provider validation receipts were inspected rather than redundantly rerun.
- Inspected risky accepted/rejected proposals: World Bank source economy versus country; CLDR grouping/currency default and dated multiplicity; catalogue price versus invoice-line price; local person/customer/employee versus global identity; ordered composite and alternate keys; GTFS feed/date scope; NaPTAN/NPTG versus rail/station namespaces; EA typed measure references versus universal stage; POWER unpivot/fill assumptions; CVE identity versus affected products; narrow tcp/udp bridge. No false-positive join approval found in those scoped decisions.
- Public query helper accepts field mappings only with reconciliation receipts, retains edges separately, and leaves dataset suitability undetermined. Its demo enrichment routes do not claim accepted joins. Meanings/sources/licensing remain independent.

## Bounded completeness and remaining limits

Most unresolved proposals retain useful source fields and relevant domain distinctions; no arbitrary target-count deficiency is alleged. The local-identifier pattern in finding 1 is a concrete exception where the evidence-gap wording suppresses useful already-supported scope. Other unresolved historical/category/typed-observation/namespace questions remain legitimate implementation or targeted-research inputs.

Primary references checked: [GeoNames countryInfo](https://download.geonames.org/export/dump/countryInfo.txt), [CLDR48 units](https://raw.githubusercontent.com/unicode-org/cldr/release-48/common/supplemental/units.xml), and [World Bank country API](https://datahelpdesk.worldbank.org/knowledgebase/articles/898590-country-api-queries). World Bank explicitly documents use of its own codes where ISO codes are unavailable, supporting refusal to treat the entire economy catalogue as ISO countries.

This review does not remeasure every D1/D2 provider payload or all 1,170 assertions independently. It inspects preserved captures and domain decisions, reproduces the high-value native/country paths, and challenges representative source semantics. It does not assert all public endpoints are currently live. Native/source bridge measurements are not DataTug product executions.

Phase2 remains pending; independent acceptance must follow corrections and root gate decision. B1 remains OPEN/BLOCKING for phases4/5/6, including runtime/Directory/pin/capability/immutable-native obligations. External publication approval remains unanswered. No source changes, canonical publication, portfolio selection, browser-opening tests or merges performed by this reviewer.

## Canonical pins

- MeaningGraph core: 982916d73f0a35ff2558b0062f58aa3ac4f24d97
- MeaningGraph registry: 2d92bbdd45fc2552f231f8f993a5b8b8d5b7fa29
- MeaningGraph CLI: a58a506d9e6f76660cc8586478521787976111d0
- ModelSpec registry: 48b30b250a61385d94d46a70968677870750ea35
- ModelSpec CLI: a8de98b4565ae669b6fba00c334b4649c9d40b11
- ModelSpec spec: ff2393a99e0bd093272c1765a57eefdea01c0a90

Local reproduction artifacts: native.json, reproduced.json and reproduce.py beside this report. Initial Node YAML extraction failed because that checkout has no installed YAML package; Python YAML extraction succeeded before the recorded native replay. Failed preparation was not counted as validation.

VERDICT: blockers=0 majors=1 minors=2 land=no
