# W1 and BigQuery delivery observations — 2026-10-06

This additive snapshot begins at hub `baa0ba7faeaf34ed4a4728c76106da5399de7237` and Workbench `9fd225742a6eeb6cf439ef935977a8c6c2dba550`. [Normalized original receipts](publication-receipts.json), [scope and corrections](evidence.json), and the [source hash manifest](source-manifest.json) preserve dated component observations. Every prior research file except the index README link remains byte-identical, including all 699 candidates and 1177 assertions, FX research, implementation records, questions and gates. Prior packets remain historical observations.

| Component | Public change | Landed main | Observation at freeze |
| --- | --- | --- | --- |
| OVDB Publisher | [PR54](https://github.com/openvaultdb/ovdb/pull/54) | `39d4214bcc128768c403f6f994e0af4cbad5b49c` | OVDB CLI v0.31.0 tag and exact-main CI observed |
| Guarded app graph/history | [PR196](https://github.com/datatug/datatug-apps/pull/196) | `d327f5403b020c64d37fe75c3a9c867e66ad3007` | Exact-main CI and deployment run 37376377336 observed; guarded browser surface only |
| Selected runtime | [PR19](https://github.com/openvaultdb/cloud/pull/19) | `3a11a8523d4f8ef3baab3e8f68382124b0d65aab` | Scoped selected-runtime component receipt |
| Protected image runtime | [PR20](https://github.com/openvaultdb/cloud/pull/20) | `d776994d5b8d512c387db14157af4272fbd26c24` | Exact-main CI, Linux component acceptance and deployment receipts; full actual W1 capacity remains open |
| Preparation fix | [PR21](https://github.com/openvaultdb/cloud/pull/21) | `383125f6e4f31d07822ed7e4a7cb12aae7fb50ee` | Exact-main CI; Cloud Run `chinook-ovdb-00015-cxx`, Worker `b4ced494-4c04-4ba6-9e06-ccaf4c65f053`; six existing gateway provider checks |
| Geo candidate wrapper | [PR9](https://github.com/ingitdb/geo-ingitdb/pull/9) | `d53aa273448e99e792899340eb4c162e67834b69` | Exact-main CI terminal; availability flags remain false |
| ROR candidate wrapper | [PR5](https://github.com/ingitdb/ror-ingitdb/pull/5) | `49bc19225c1f2c4e96fbf414bf2359b4550e1f17` | Exact-main CI terminal; availability flags remain false |
| Go BigQuery execution | [PR3](https://github.com/dal-go/dalgo2bigquery/pull/3) | `0a5596742d43b8be171ce2dba6b35b6872146145` | dalgo2bigquery v0.1.0 public release and downloaded module |
| Go BigQuery Snapshot | [PR4](https://github.com/dal-go/dalgo2bigquery/pull/4) | `7df3d802f2c5188025c56fd8e41d27f01206a303` | dalgo2bigquery v0.1.1 public release/tag and downloaded module |
| Go driver documentation | [PR5](https://github.com/dal-go/dalgo2bigquery/pull/5) | `5290181f7fec4adf6de8dedffb8eca8f5921dc8f` | Exact-main CI terminal |
| DataTug CLI BigQuery | [PR343](https://github.com/datatug/datatug-cli/pull/343) | `3a2d617f98e3762c0d3c26ce21ff69bae51dd11f` | Exact-main CI terminal; public DataTug CLI v0.60.0 tag points at this merge; official archive checksum, strict host code signing, exact binary version and six subcommand help surfaces root-verified |

Landing observations retain original `head_sha`, `reviewed_head` when present, merge commit, branch deletion and canonical synchronization fields. CI observations retain explicit original `head` and per-check conclusions. No routing index head is inferred to be an observed CI head. Failed landing attempts and earlier pending CI remain in the receipt history. A terminal skipped check is preserved as skipped. Public tag/ref observations and Go module origin hashes are separate from CI conclusions. The binary receipt retains sandbox codesign exit 1 and successful strict host codesign exit 0 separately. Package-manager publication is a separate workflow; no Homebrew tap update is claimed.

## Preparation evidence correction

The macOS local prerequisite proof recorded 16 physical tables, 11 selected serving tables, five logical public recordsets and 1,280,314 scanned keys. Its native/key proof remains scoped local evidence. Its generator-pass claim is explicitly invalidated: both historical generator logs reported `metadata/checksums.json` mismatch and generator exit 1; a later `git diff` exit 0 was wrongly used as the generator result. Original receipts and logs are retained by hash, without rewriting the failure.

The corrected checksum-only successors are Geo `afd2b879979078f18db4d11c509ae5d7a2f85dfd` and ROR `228d4b235dff0561785a026f411d3295fc656d34`. Fresh independent wrapper review verified every checksum entry, standalone generator exit 0 and unabridged OVDB Publisher check exit 0 at those exact heads. Only their successful successor checks supersede the historical generator failure. The original blocked heads remain rejected. The wrapper review's local finding named B1 is distinct from programme B1; it grants no programme closure, runtime capacity acceptance or public source admission.

## Remaining gates

Phase 4 remains active. Programme B1, phases 5/6, full actual W1 Linux capacity/startup, 128MiB spool/concurrency proof, DataTug CORS and transport, public Geo/ROR wrapper/Directory/App admission, and the real J1–J6 journey remain open. The six existing gateway provider checks do not establish DataTug-origin CORS or public W1 admission. Source-generation and local key proof do not establish deployed startup/capacity.

Go and JavaScript BigQuery delivery are both mandatory. JavaScript foundation publication remains blocked by automatic approval review requiring exact human approval, and independent execution review is platform-blocked. Operator-controlled execution identity/project and explicit positive spending cap are unresolved; credentials and paid queries have no authorization from this packet. Actual Go CLI/local-server and browser live acceptance, canonical source-access profile and source-specific rights/admission gates remain open. Neither Go release nor fixture success closes those gates.

## Verification

Run `python3 spec/research/public-data-fabric/evidence/w1-delivery-catchup/verify.py` from the hub root. It checks packet hashes, original field associations, exact-main head/check coherence, release-to-tag/module bindings, history/correction boundaries and every preserved research file. It makes no network request. Optional `--original-map` accepts a private label-to-file map to rehash original receipts and compare normalized fields directly; the private map is never committed. Default verification validates the frozen public packet, not the private originals or current remote state. Receipts are dated observations; this packet is not a new semantic or phase acceptance decision.
