# Independent review r1 — W1 canonical registrations

Reviewed-Head: 535c98308058650840d66bb336f05ecd10ec690f
Reviewed-Head: 603f0837c56aad96e05207a7336ad1eeb513e465

**Verdict: approve content; zero blockers, majors or minors.** Root remains the sole landing owner. Required exact-head CI must pass before landing; MeaningGraph CI was still IN_PROGRESS with an empty conclusion at final observation, while ModelSpec CI was SUCCESS. This review is not a merge, publication or programme gate receipt.

## Exact reviewed scope

- meaninggraph/registry#18 at `535c98308058650840d66bb336f05ecd10ec690f`, base `2d92bbdd45fc2552f231f8f993a5b8b8d5b7fa29`.
- modelspec-org/registry#16 at `603f0837c56aad96e05207a7336ad1eeb513e465`, base `48b30b250a61385d94d46a70968677870750ea35`.
- Both PR API headRefOid and live branch refs matched; both PRs OPEN and non-draft. Saved remote/run JSON binds evidence to those exact heads.
- Compared all changes against base Git objects. Read the governing DataTug mission, approved W1 Workbench plan, applicable rules, existing collection/builder/validator/test contracts and landed provider model/binding/license objects.

## Findings

None. No declined findings.

## Verified content

The existing canonical record homes are used. MG adds two graphs and two explicit core dependency records; ModelSpec adds two source models. The registered GeoNames commit is `6f4cf1269bc393048f6b204069135a62f0bb6c02`; ROR is `dc78c1e929f1f10018f8c39e690351059a50c71f`; core remains `982916d73f0a35ff2558b0062f58aa3ac4f24d97`. Canonical validators verified paths, immutable commits on default-branch histories, declared licenses, meaning/model linkage and dependency closure. Model JSON and binding SHA256 values were independently recomputed from pinned Git objects and match the root landed-tree receipts (see static-evidence.json). Both model HCL/JSON pairs validate at these revisions.

Source-qualified addresses/modules and recordset grain are retained: four logical GeoNames recordsets with auxiliary source-scoped country bridges, and one logical ROR organization recordset with ordered location/relationship projections. The records grant no user identifier mapping, raw-label compatibility, name-to-identity inference or runtime eligibility. Provider bindings retain their existing accepted meaning scope; no meanings or reconciliation decisions are introduced.

The new entries declare CC0 model/meaning metadata while distinguishing GeoNames CC-BY-4.0 data and ROR's embedded GeoNames attribution. The pinned ROR meaning source and DATA-LICENSE.md declare source ROR metadata CC0. Both provider code licenses are separately declared CC0. The existing trakhimenok maintainer record is reused; live GitHub provider commit author identities are trakhimenok, consistent with actual provider contributions and root assignment. No maintainer was invented.

MG count changes 7->9 and ModelSpec 6->8. Every prior index entry and other top-level metadata remains equal; exactly geonames and ror are added. Both existing builders reproduce the committed bytes. MG index SHA256: `42931545086beb7c69d803767b8d50e524baff2ff37a2ec9eb8526eee16aa4a5`. MS index SHA256: `302a04f9775aec8627ddab0361030156f592abb3e446088ba43bd813aa527576`.

The MG test change adjusts the exact graph count and the test's anchored explicit repository regex by adding only ingitdb/(geo-ingitdb|ror-ingitdb). The assertion retains nonempty requests, HTTPS/github.com, exact path/end anchoring and the existing provider alternatives. No production allowlist, safety validator or failure assertion is weakened. The changed test passes.

## Independent checks

Executed through WB from private Git-archive snapshots with private copies of registry caches; author fixtures/caches and worktrees were not edited. Node 24.20.0, local inGitDB 0.67.1, pinned SpecScore 0.54.2.

- MeaningGraph focused suite: 5/5. Canonical validation covers all 9 graphs; negative controls reject wrong dependency pins, missing/unused dependency records, omitted model declarations and refused homepage URLs while allowing only explicit repository fetches.
- MeaningGraph inGitDB controls: 34/34.
- ModelSpec offline semantic/security suite: 54/54.
- ModelSpec inGitDB controls: 32/32.
- ModelSpec canonical validation: 8 models; pinned HCL lint: 8/8.
- Both index reproductions and both exact diff whitespace checks pass. Both author worktrees remain clean.

The first MG inGitDB invocation lacked INGITDB_CLI and failed setup; the corrected environment passed 34/34, with both logs retained. The full ~9-minute MG51 suite was deliberately not repeated: unchanged security implementation, five relevant full-registry/negative controls, inGitDB controls and required CI supply appropriately bounded coverage. CI uses pinned inGitDB0.68.2; local0.67.1 evidence is not represented as that release.

## Landing boundary

Root must obtain terminal MeaningGraph exact-head CI success before landing and retain normal main-CI/tree/sync/cleanup receipts. MS exact-head run37261670562 succeeds; MG run37262128150 remained IN_PROGRESS at this review. This review approves registrations only. User-contract acceptance, representation publication, runtime/Directory eligibility, resource/source-generation proof, real-browser journey acceptance, B1 and all broader programme gates are outside this verdict and remain governed by their owners.

VERDICT: blockers=0 majors=0 minors=0 land=yes
