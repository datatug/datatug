# [review r2 ovdb#43] Schema repair delta

Reviewed-Head: 191cf7d36618ee529696ce758ad007f0db61309d
Reviewed-Previous-Head: a22c7cd5be916108e4aedc0ca8afcfe438447003
Reviewed-Base: aec07d17be1e074210acefddb858363ce7777704
Runtime/host: local Codex, macOS. Author worktree read-only and clean.

Source verdict: PASS. R1 and R2 are resolved; no new findings. Root may land this staged helper only after all applicable CI jobs at this unchanged Reviewed-Head succeed. The CI-only condition is open at freeze: run37256175837 is in progress, with web/browser and Windows lifecycle still running. Every other applicable job has succeeded, including exact publisher coverage, parity goldens/real-Git publisher checks, Ubuntu/macOS lifecycle, build/test, lint and tidy. Version bump is intentionally skipped. Authoritative PR is OPEN/draft with this exact head and Reviewed-Base. Receipts: pr-r2.json and ci-r2-final.json.

## Delta and reproductions

Only schema2.json, README.md and new schema_registry_test.go changed since r1. No production Go code changed.

- R1 resolved: format 2 declares `https://openvaultdb.com/schemas/representation-contract-2.json`. Private `schema_registry_r2.go` registers both exported schemas using their own declared IDs in the real jsonschema/v6 compiler, compiles both, accepts real ROR and GeoNames fixtures under format 2, and rejects them under format 1. Output: schema-registry-r2.log. Author regression also validates the format 1 fixture and wrong-version controls in the same registry.
- R2 resolved: path regex admits ordinary safe components and literal `$records` components only. The unchanged r1 private probe now rejects `$HOME/x` at Parse. Independent race run executes both exported-schema and Parse tests for embedded dollar, `$HOME`, `$recordsx`, suffixes, shell forms, percent encoding, URL/absolute/traversal/empty-component negatives. Valid canonical Geo decision `$records` path and repeated literal components pass. Output: probe-r2.log and race-r2.log.
- README explicitly requires reviewed source/model/binding/native-data continuity for later wrapper revisions and dedicated specialist review for semantic change. Structural Check cannot automatically inherit or widen the pinned decision.

## Verification

1. Independent `wb run -- go test -count=1 -race ./publisher/representation ./internal/publisher/repo` passed both packages at Reviewed-Head.
2. Independently reran exact covergate over the author's fresh single full profile `/private/tmp/data-fabric-w1-native-contract/repair-r1-full-cover.out`: 2324/2324 statements across eight packages. This is explicitly an author-generated profile, not a second independent full-suite run; independent full-profile verification at the previous head is preserved in r1 and the repair alters no statement-bearing production Go. Current exact-head CI independently completed its full race/profile and gate successfully.
3. Schema2 SHA256: fa281400fbe063d1abbcbf5c7578dd23932dd798ee91810e0906ffa2431dc0bc, independently matched.
4. Format1 schema bytes, default manifest/rules, parity reference files and workflow remain identical to Reviewed-Base. No second publisher profile/default admission was introduced. Final local HEAD exact and git status clean.
5. The original r1 metadata/native SQLite and actual small companion compatibility evidence remains applicable: the repaired schema accepts existing canonical fixture paths and changes no native linkage/reader implementation. Separate ROR provider PR3 and its new packaging head are not accepted by this helper review.

No production eligibility, chunk reconstruction, runtime hosting, canonical companion admission, app user-contract execution, B1 or whole-W1 acceptance is granted. Native logical dataset identity remains linked without bulk metadata reads; provider/runtime review independently proves decoded chunk bytes. Later metadata-wrapper pin carry-forward remains an explicit reviewed obligation.

Reproduction commands (run from the author worktree):

    wb run -- go run /private/tmp/data-fabric-w1-native-contract-review/probe.go
    wb run -- go run /private/tmp/data-fabric-w1-native-contract-review/schema_registry_r2.go
    wb run -- go test -count=1 -race ./publisher/representation ./internal/publisher/repo

VERDICT: blockers=0 majors=0 minors=0 land=yes
