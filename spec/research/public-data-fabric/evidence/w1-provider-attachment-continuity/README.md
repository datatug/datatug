# W1 landed provider attachment continuity

[The packet](evidence.json) binds Geo PR 8 and ROR PR 4 to their exact reviewed heads, merged trees and `model/representations.json` hashes. A separate continuity check observed both merges and remote main identities, equal reviewed/merged trees, 125 immutable blob hashes, seven ordered native chunks, and 9,857 Geo / 17 ROR unchanged prior tree entries. Geo merge checks had five successes and three skips; ROR had one success. Receipt timestamps describe that observation, rather than a continuing CI guarantee.

The original four byte-exact Geo country dictionaries and original eight-row ROR affiliation input remain pinned to their independently reviewed decisions, schemas, models, bindings and native snapshots. Original declaration/role and affiliation grain remain unchanged. Future schema-conforming inputs receive no acceptance here. Original review round 1, including the OVDB v0.27.0 reader failure, remains preserved; round 2 accepted provider metadata landing after the actual released OVDB v0.27.1 reader passed. This continuity check verified those receipt hashes and equal trees; it did not rerun the released harness.

`original_documents` associates small public labels, raw SHA256 hashes and byte counts with the original private review, approval, released-reader and root landing receipts. These are provenance declarations, not copied raw tool envelopes. This packet still awaits a separate adversarial review; its `fabric.evidence` observation has `independent: false`. The prior semantic decisions retain their original authority. No publication/acceptance arrays, implementation state, root manifest opt-in, Directory admission, runtime eligibility or launch status change.

Run the bounded offline check from the repository root:

```sh
python3 -B spec/research/public-data-fabric/evidence/w1-provider-attachment-continuity/verify.py
```

Like the existing W1 packet verifiers, optional `--repo GITHUB_URL=LOCAL_CHECKOUT` arguments rehash immutable Git objects without fetching; `--review-originals DIRECTORY` and `--landing-originals DIRECTORY` rehash the associated raw originals when available. The output distinguishes declared receipt hashes from actually rehashed originals and Git objects. This command checks continuity; it does not execute data builds, semantic reconciliation, the released-reader harness or a product journey.
