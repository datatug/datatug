# Independent adversarial review r2 — OVDB PR42

Reviewed-Head: 18cfe991dc4dd47fece2ec354dda95146e7d260a
Base: 73bea463f2dfcfe610358d004707260b5ceeda31
Previous reviewed head: 983e95eef01f2e8815de611e496c466caf053a9d
Runtime/host: local Codex on macOS. Root remains sole landing owner. No source worktree edits or merges.

## Source result

Accepted as a staged offline helper implementation: all three r1 majors are closed, with no new blockers/majors/minors found. This does not freeze production attachment admission or approve semantic scope, Directory/app publication, the ROR user contract, or whole W1. The exact seven-file repair diff was reviewed. Default publisher, parity fixtures/reference pins, workflow/gate inclusion and runtime source are unchanged in this repair.

R1: exactKeys now rejects case-folded aliases of consumed JSON fields before struct decoding. Bridge/key-index fields are closed; unrelated canonical fields and case-distinct dynamic entity/property identifiers are preserved. Original bridge and ModelSpec alias exploit probes now reject.

R2: scalarEscapes rejects unpaired high/low surrogate escapes before Go can replace them. Original lone-surrogate probe now rejects; valid supplementary scalar pairs, literal/escaped U+FFFD, escaped backslashes and quotes retain correct lookup identity.

R3: singleYAML requires EOF after exactly one YAML document for binding, canonical meaning and manifest attachment parsing. Original trailing-document probe rejects. Empty/malformed additional documents reject; document terminators and trailing comments still pass.

## Independent execution

A private archive of the exact corrected commit is under `r2-source/`. The original r1 reviewer probes were copied unchanged except renaming their function to avoid collision with the author's regression test; no source worktree is touched. Private numeric/type controls were adapted only to the renamed exactClosed helper. An additional reviewer-only test exercises every property in the actual corrected GeoNames ModelSpec, its complete YAML binding document and snapshot bridge/key-index closure.

Commands (cwd `/private/tmp/data-fabric-w1-contract-review/r2-source`):

```sh
wb run -- go test -count=1 -race -covermode=atomic -coverprofile=/private/tmp/data-fabric-w1-contract-review/r2-cover.out ./internal/publisher/... ./publisher/representation ./internal/covergate/... ./cmd/covergate/...
wb run -- go run ./cmd/covergate /private/tmp/data-fabric-w1-contract-review/r2-cover.out ./internal/publisher/exitcode ./internal/publisher/manifest ./internal/publisher/repo ./internal/publisher/rules ./internal/publisher/checkcmd ./publisher/representation ./internal/covergate ./cmd/covergate
wb run -- go test -count=1 ./publisher/representation -run 'TestIndependent|TestReviewControls|TestUnicodeScalarIdentity|TestOneCompleteYAMLDocument|TestExactConsumedJSONKeys' -v
```

All pass. The independently executed race gate covers **2214/2214 statements in exactly 8 packages**. Full focused tests include the unchanged default-publisher rejection/parity controls. NaN/Infinity/overflow, duplicate JSON keys and numeric/boolean-to-string coercion controls still reject. Evidence: `r2-tests.txt`, `r2-cover.out`, `r2-probes.txt` and private tests in `r2-source/publisher/representation/independent*_test.go`.

Both corrected GeoNames bundles independently verified `GeoNames:countryInfoISO2`, key-index SHA256 `c57b21174d0f0d2badd575cf1c2605840bc862e6cebd52064ceffd9c266ee9c8`, exact matching snapshot entry, and identical snapshot SHA256 `3fce69153378a0af22f049af8bcd507e3e79ba34b2020ac767cee6c0ee4977bc`. This closes the earlier namespace observation only; native source truth/whole bundle publication remain their separate reviews.

## Remote and CI receipt

Read-only authenticated PR API verified exact corrected head and unchanged base, open draft. Commit check-runs API captured workflow run37252803923 and job conclusions in `r2-check-runs.jsonl`. The previous Linux port-collision failure does not recur: Linux job111583708213 completed with SUCCESS, without unrelated runtime edits. Exact coverage, Build & test, Lint and tidy also succeeded.

At receipt freeze these checks remain in progress: Lifecycle and paths (windows-latest), Parity goldens are current (publisher check packages). There are no observed failed checks. Source review is accepted, but landing remains pending these exact-head conclusions. Root can close this CI-only gate from authoritative check-run success without another source review if the head does not change.

VERDICT: blockers=0 majors=0 minors=0 land=no
