# Review r1: configured user-source admission

Scope: read-only architecture review on local Codex/macOS. This recommends an implementation boundary for root approval. It grants no new semantic acceptance, production eligibility, deployment authority, or programme gate closure. No repository files were changed. The app author retains the isolated `fabric-w1-native-revalidation` lane.

## Recommendation

Approve a **declared user source adapter in the existing project catalog path**, plus a **bounded driver-source stage in the existing federation executor**. Preserve the current canonical public-provider adapter separately. A user's schema can be verified directly against its immutable repository/revision/path/hash and the exact accepted representation contract; it need not be a globally registered public model. Public targets continue to require canonical ModelSpec, MeaningGraph and Directory publication. Source declarations, structural compatibility, accepted semantic scope, transport readiness and executable admission remain separately represented.

This is implementation infrastructure missing from the current app. The exact accepted Affiliation fixture requires no new semantic reconciliation when its source/data/schema coordinates and representation remain unchanged. Extending acceptance to another source/revision/property, changing its raw representation or applying its decision to arbitrary user data would require a distinct reconciler. General private-auth and local mutable-schema support require additional technical identity/access work; they are not established by this recommendation.

## Verified evidence

Authoritative `git ls-remote ... refs/heads/main` observations in this review:

- apps `e6569c5d606af0c4094b9e9d7db6ca1c5964a62b`
- hub `17263dbacabdfe95e53fc3c6980177416bdab941`
- ModelSpec registry `34be21159e522e85b6e5c110f8bfb8baaa8035a8`
- MeaningGraph registry `3d85ae6fa06b7f38f59cfd2893e07ba801856b96`
- OVDB helper `bb535a658f06a8752688ba36cb4b76ecf0483254`
- ROR `dc78c1e929f1f10018f8c39e690351059a50c71f`
- Geo `6f4cf1269bc393048f6b204069135a62f0bb6c02`

Read full mission brief and approved Workbench W1 plan, app AGENTS/Nx skill, relevant working-agreement rules, exact app main source and author's explicitly uncommitted native worktree/question. Local canonical app and OVDB working trees were clean when read, and exact Git objects were used for the decision, fixture, registry and helper facts. This review does not report CI or deployment results. A read-only `NX_DAEMON=false pnpm nx show project datatug-main --json` attempt failed before Nx startup because the sandbox could not fetch/verify the project-selected pnpm release; no package-manager safety bypass was attempted. Resolved Nx targets are therefore unverified by this reviewer and remain the implementation owner's check. No service or browser was started.

Recomputed Git-object hashes:

- hub decision `spec/research/public-data-fabric/evidence/w1-ror-user-decision.json`: `a4caa4b6461230a3fa2ae8548121344cb286cfb98a9e19761a207c9f74479595` at hub `17263dbacabdfe95e53fc3c6980177416bdab941`.
- apps schema `libs/datatug/main/src/lib/queries/fixtures/public-data-fabric/affiliations.modelspec.json`: `cfe958f12d9a70bb7533499b76bf01c9447f81f0862852e460b095bb75e1b9cf` at `e7362033ec79c6663d7dbe0483b62fab01f7b9cd`.
- apps data sibling `affiliations.json`: `44a30dd260c74f43b2138934cca0bccd493c7b804e124c5ecbded3c41226e543` at that same fixture commit.
- ModelSpec registry exact index contains no `https://github.com/datatug/datatug-apps` source entry.

The helper's `publisher/representation/contract.go:204` already accepts an explicit external immutable source schema via `Context.Resolve`, verifies bytes and exact ModelSpec property, and does not require source membership in a global registry. `internal/publisher/repo/representation.go` remains opt-in and expressly does not grant semantic admission. The recommended source distinction aligns with this existing helper boundary.

## Findings

### M1 — Public identity is incorrectly mandatory for configured user discovery

At app main `queries/public-data/representation-discovery.ts`, `verifyContract` finds a source ModelSpec registry record before reading the already pinned schema. `configured-source.ts:49-85` independently requires a Directory provider and globally registered model. `configured-public-data-sources.service.ts:81-131` retains at most upstream repository/revision and table names, with no independently declared source schema or physical mapping. Consequently the accepted fixture cannot reach ordinary configured-source inspection without a fabricated public registration. Removing just `findRecord` leaves this second barrier and fails to establish that the selected connection actually contains the declared source.

Severity: major; normal J1/J2 user path is missing, while current fail-closed behavior prevents unsafe execution.

### M2 — Source transport is incorrectly coupled to public Directory and a shared OVDB base

At app main `public-data.service.ts:112-154`, scenario construction requires a Directory record for the source repository+revision and requires both source and target routes under one allowlisted OVDB base. `bounded-federation.ts` then assumes all source stages are OVDB requests. This is independent from model compatibility and cannot be satisfied merely by declaring a schema. Test helpers currently invent a user Directory runtime entry to exercise the factory.

Severity: major; there is no actual configured fixture transport through this path. This is a missing driver integration, not evidence that the semantic mapping is unresolved.

## Smallest concrete integration

### 1. Normal project configuration and explicit physical mapping

Reuse `ProjectService`, `EnvironmentService`, `GithubProjectReaderService`, the existing `environments/<environment>/catalogs/<catalog>/<catalog>.db.json` location, and `https-json-catalog.ts`/its JSON schema. The latter already defines `driver`, `urlTemplate`, `keys`, table SHA256, optional upstream, address checks and bounded file parsing. Its closed key list means source-model support must be a reviewed additive schema/parser change, not an unchecked extra field.

Add one typed optional source-model declaration to this catalog format (field name finalized by the author). It should contain an immutable schema reference, exact module name/id, and explicit mappings from physical schema/table/field to model entity/property with declared raw datatype and identifier namespace. Reuse `ImmutableFile` and the representation source coordinates; do not create a second semantic mapping catalogue. Namespace is a declaration cross-checked against the accepted contract, never inferred from a property name or auto-filled as acceptance. Record the exact project configuration commit/path/hash separately from the upstream schema/data pins.

A concrete W1 configuration uses the **existing** https-json driver with physical table `affiliations`, key `affiliation_id`, and `urlTemplate` `https://raw.githubusercontent.com/datatug/datatug-apps/e7362033ec79c6663d7dbe0483b62fab01f7b9cd/libs/datatug/main/src/lib/queries/fixtures/public-data-fabric/{table}.json`. `sha256.affiliations` is the exact data hash above. Its explicit model mapping is empty physical schema / `affiliations` -> `affiliations.Affiliation`; `affiliation_id` and `ror_id` map byte-for-byte to those properties, with `ror_id` declared nullable string and namespace `ROR:URL`. The exact schema reference stays at the accepted old commit even if the containing project configuration is published later. The module id is `github.com/datatug/datatug-apps/affiliations`.

This is a normal project-catalog fixture, not a new product picker special case. Put parser/UI fixtures under the existing project-files test conventions, then publish actual project configuration through an existing project owner before claiming a connect journey. No new repository or global schema registration is needed. No filename, research-ID or namespace switch may select this fixture in application logic. Schema and table inspection should show the configured declaration and its unavailable/compatible status through the existing configured-source UI. Do not expose an arbitrary descriptor-text box as the supposed normal flow.

### 2. Verify source identity without global publication

Extend `ConfiguredPublicSource`/`ConfiguredFieldChoice` with a discriminated resolved source context: the existing canonical provider branch, and a declared project-source branch. The latter must be constructed from validated selected project configuration plus fetched schema bytes, not from an arbitrary suggestion's source fields. It carries connection identity, configuration pin, immutable source schema ref, exact physical-to-model mapping, and the driver descriptor.

Refactor `verifyContract` to separate source verification from canonical public-target verification. For declared user sources, verify the schema bytes directly against the reference supplied independently by the selected connection and compare all repository/revision/path/SHA256/module/entity/property/type/namespace coordinates to the contract. Check module id, optionality and key/grain when present in the accepted source context; do not make NULL an identifier. No trimming, case folding, entity-name inference or substitution of same-name models. Preserve the existing public/demo provider adapter's registry and Directory linkage.

All target ModelSpec record/file/entity/property checks, MeaningGraph binding/role/core/dependency pins, snapshot/native proof/attachment hashes and exact decision-reference checks remain mandatory and unchanged. Do not derive acceptance from `disposition`, scope prose, labels, research IDs, or the declaration. Structural compatibility may be shown while publication/semantic admission remains unavailable. The exact accepted decision remains referenced; this adapter does not reissue it.

For W1's public Git source, `CanonicalMetadataReader` already provides immutable raw GitHub URL/hash/budget/redirect checks and can be used for source schema bytes with source-specific labeling. A project-store reader abstraction may be injected for future private schemas, but it must prove exact commit/path/hash on returned raw bytes and scope caches to account/project/source. `GithubProjectReaderService.getRawText` preserves text and accepts a project ref, but currently reads anonymously and may try a public mirror. It is **not an authenticated private reader**. Do not claim or implement private credential forwarding through the public reader. An unavailable private source should remain explicit unavailable, without prompting global publication. Supporting private Git access is technical work; adopting local content-hash identities in place of repository/revision is a separate identity-contract review, with reconciliation required if it changes the accepted semantic source scope.

### 3. Generic bounded source transport in the existing query path

Separate query relation identity from remote route identity. Scenario creation should receive the verified configured source context rather than recover every driver from public Directory. Public targets still resolve through Directory to the explicit approved OVDB route. Save source connection/configuration/data/schema/mapping pins in the existing query's public-data provenance; do not save credentials. A logical driver relation does not claim that a `user` database exists at a public OVDB URL.

Introduce a narrow `read driver page` boundary inside the current bounded federation flow. Keep the existing OVDB implementation for canonical provider drivers. Add a generic https-json implementation using the validated catalog's `tableUrls`, declared primary key and **mandatory exact hash for the selected source table**. Feed verified selected records into the same DALgo/IndexedDB/worker execution and exception collector. Source rows retain `affiliation_id` grain, original nullable `ror_id`, duplicates and raw strings. Native serving IDs never replace these values.

`FederatedQueryObserver.fetch` and worker message/extras are existing integration seams, but the safe extension is an explicit driver resolver or dispatch for the exact selected logical driver relation, not an all-requests fallback. Every public target stage must continue through `createBoundedFederationFetch` and a real allowlisted OVDB POST with the same bounds/provenance. Unknown relations reject. No public target may be satisfied from user-source JSON, app bundle, cache-only fake receipt or `static-ovdb-fetch.ts`. That file is demo-only and presently lacks the required hash, byte and scope controls; do not reuse it unchanged or treat its synthetic snapshots as runtime publication.

The accepted user JSON is the data source itself. Reading that small immutable user source is not evidence that ROR/Geo runtime hosting exists. This preserves the plan's prohibition on static public-target fallback.

### 4. Transport controls and explicit limits

- Start row reads only on the user's explicit run; discovery loads bounded metadata only. Opening/saving/revalidating a scenario performs no row fetch or public lookup.
- Require HTTPS, checked expanded table URL, exact immutable Git source/data reference for this W1 path, credential-free URLs, no arbitrary request-method/body/host overrides, `redirect:error`, and rejection of redirected/opaque responses. Browser CORS must succeed naturally; no proxy or mode `no-cors` workaround.
- Reuse streaming `boundedResponseText`, strict duplicate-key JSON parsing, AbortSignal/deadline propagation and aggregate byte accounting. Charge the original downloaded JSON bytes plus every public response/error body to the same <=5MiB run limit; no second uncharged fetch through an observer. Retain <=10second total deadline and cancellation through worker and both transports.
- A plain JSON array cannot support server-side pagination. For this minimal driver, explicitly support only whole source arrays of at most 1000 rows within the byte cap; reject larger arrays without silent truncation. User-selected page size is <=1000 and is applied after verifying the bounded complete file/hash. Offset/page selection is explicit and stable against that immutable data hash. Document that reading any page reads this bounded whole source file. Larger sources need a real paged driver; do not suggest scalability from the eight-row fixture.
- Validate unique nonempty declared row keys; refuse wrong root shape, malformed/duplicate JSON object keys, duplicate row identity, unknown configured relation, missing required key, wrong physical mapping and nonstring/nonnullable values. Preserve optional NULL separately according to the admitted contract; do not coerce values.
- Valid empty array yields zero driver rows and no public request. Empty response/malformed JSON fails. A selected page containing only NULL/empty/invalid identifiers retains each exception and makes no lookup request. A valid but absent target remains unmatched, not invalid or a fabricated row.
- Enforce distinct identifier limit (50 ROR/100 place), <=5000 output rows and all existing target snapshot/expiry/pagination checks. Duplicate affiliations with the same ROR key remain two inputs but one deduplicated lookup key. Locations remain details without multiplying affiliation denominators.
- Separate source credentials from public target credentials and never forward project/private-store headers to raw GitHub or OVDB. W1's source needs no authentication.

### 5. Save/reopen and proof

Save immutable source data and schema refs, configuration pin, selected physical mapping, source grain, driver kind, selected bounds/page, exact decision ref and all current canonical/public snapshot pins. Reopen does not run. Fresh validation re-resolves the selected configured connection and compares every saved coordinate; source config/data/schema/mapping changes, target changes and missing publication/transport all remain explicit. Rebuild an acknowledged plan only from freshly verified objects; a saved boolean or a hand-edited query never grants admission.

Tests should remove the hypothetical global user ModelSpec/Directory entries from the positive declared-source setup. Keep hypothetical **target publication envelopes** labeled as consumer structural tests while owner-controlled public admission is closed. Use exact accepted data/schema/decision bytes, plus mutation controls. No generated test envelope counts as actual target publication.

## Required negative controls

1. Exact source works structurally with no source entry in public ModelSpec or Directory; public target removal still rejects.
2. Independently mutate source repository/revision/path/hash/module/module id/entity/property/datatype/namespace/configured table/field/key: source unavailable, no target query. Same named field in another project/connection is not accepted by name.
3. Missing/ambiguous mapping and changed catalog/config pin reject; source declared only by a target attachment cannot connect a user's table. Different source data hash at the same schema cannot reuse the exact fixture acceptance.
4. Alter decision bytes/reference, proof scope, attachment, target binding/meaning/snapshot pins: closed. A JSON `disposition:accepted` or accepted-looking prose has no effect.
5. HTTP failures, CORS failures, redirect, unapproved URL, missing/wrong hash, invalid UTF8, duplicate JSON keys, oversized response/array, duplicate row keys, cancellation and timeout produce failure with zero success result; byte guard includes error bodies.
6. Row boundary1000/1001, identifier50/51 or100/101, result5000/5001, bytes exact/over, zero/all-invalid input, explicit page offset and no automatic next-page tests. Hash/cancel failures must not send public queries from partially read rows.
7. Verify target requests use real network transport in integration tests and cannot be intercepted by the user-driver adapter. Discovery/save/reopen/null action sends none.
8. Saved stale config/schema/data/physical mapping/decision/public revision cannot execute or silently refresh. Production admission remains closed despite structurally valid input.

## Remaining prerequisites and decision classification

Root may approve this technical boundary now without reconciling semantics. The app owner can implement the project-catalog declaration, declared-source verification, bounded driver page adapter and saved provenance with current public execution still closed. A small independently reviewed existing-project configuration publication is needed for a real connect flow. Private auth/local identity support is a separate technical scope and should not be represented as complete by W1's public immutable fixture.

Owner-controlled default publisher/Directory representation attachment admission and decision linkage remain absent prerequisites. The opt-in helper is only a validator. Exact native-id runtime mount, recordset/model mapping, provider release/binary/service-origin receipts, target CORS/read-only/snapshot behavior and B1 remain with designated owners. Live J1-J6 and production closure still require those publications and independent real-browser proof; this review cannot authorize them.

Falsifier of the infrastructure conclusion: an existing exact app-main project declaration adapter plus bounded driver route serving this fixture without fabricated registry/Directory entries would eliminate M1/M2. Scoped source inspection found neither; the current source-only helper provides the intended identity mechanism, while the user runtime stage must be implemented.

Rejected options: register the fixture globally merely to satisfy the consumer; delete source checks with no configured context; infer scope from source names; interpret free-prose decisions; invent a new acceptance ledger; repurpose the public static demo fallback; manufacture a source OVDB database or live target receipt.

Architecture verdict: recommend root approval of this bounded implementation design. Current normal user-source path retains two major implementation gaps; no unsafe execution bypass was found or recommended. `land=no` is scoped to declaring this path complete before implementation/review, not a programme gate decision.

VERDICT: blockers=0 majors=2 minors=0 land=no
