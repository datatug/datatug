# Core investigation loop: normative transport contract

This is the transport appendix of the [owning Feature](README.md), not a separate
feature or evidence that these endpoints already conform. The 9 September review
found incompatible implementations. Corrective plan tasks must migrate both sides
and validate actual serialized messages against this contract before marking the
journey accepted. These are engineering decisions for the approved scope, not
additional founder rulings.

## Scope and identity

All routes are relative to `/datatug`. JSON uses the exact field names below.
GET parameters go in the URL query; POST parameters go only in the JSON body.
All scoped calls also carry the current securityContextId; the server validates it
before execution and responds STALE_CONTEXT after principal/policy-session changes.
An ID is a staleness check, never authentication. agent-info obtains the initial ID.
Repeated identity in two locations, duplicate JSON keys, unknown security-relevant
fields, and a client-supplied principal or role are rejected, never reconciled by
precedence. Legacy routes must delegate to the same authorized executor or fail
closed; a migration must inventory and update every browser/CLI caller.

`Scope = {storeId?: string, project: string, environment: string, securityContextId: string}`.
Project and environment IDs are required, nonempty, and resolved within this server's
registered projects. `storeId` is additive and optional for existing single-store callers;
when absent, the server resolves the project only within its configured primary project store.
Cross-store requests and every persisted cross-artefact reference MUST name it explicitly. A
project without multiple environments declares one explicit `default` environment;
the server does not silently switch environments. Browser state additionally keys
scope by the normalized agent URL and server-issued `securityContextId` from
agent-info. That opaque ID changes on principal or policy-session changes and
confers no authority. It must not expose credentials or identify another session.

`SourceRef = {source: string, collection: string}`. Source is a stable project-local
ID, resolved through the project registry for the requested environment. The
registry may adapt existing environment/catalog/QueryDefTarget records; it must
not introduce a second competing persisted source store. A registry entry identifies
the adapter and configured target. Semantic discovery and execution use the same
resolver. Neither filesystem paths, arbitrary HTTP URLs nor credentials are accepted
from the browser as a source reference.

Saved query targets resolve to registered sources. With exactly one eligible target,
the server chooses it and reports it. With multiple targets the user selects a source
explicitly; the server validates membership in that query's allowed targets. Zero
eligible targets returns SOURCE_UNAVAILABLE. This works for SQLite, inGitDB and HTTP;
an HTTP request is not forced to invent a database name.

## Shared JSON types

```typescript
type ExecutionRef = {storeId: string; projectId: string; executionId: string};
type IncidentRef = {storeId: string; incidentId: string};
type TypedValue =
  | {type: 'string'; value: string}
  | {type: 'number'; value: number} // finite, exactly representable JSON number
  | {type: 'integer'; value: string} // canonical decimal integer, no leading +/zeros
  | {type: 'decimal'; value: string} // canonical decimal; preserve precision
  | {type: 'boolean'; value: boolean}
  | {type: 'date'; value: string} // YYYY-MM-DD, valid calendar date
  | {type: 'datetime'; value: string} // RFC3339 normalized to UTC
  | {type: 'null'; value: null};
type PhysicalRef = {source: string; collection: string; column: string};
type Fact = {
  id: string; entity: string; field: string; value: TypedValue;
  condition?: '==' | '!=' | '>' | '>=' | '<' | '<='; // OPTIONAL; absent means '=='
  origin: 'selection' | 'context' | 'manual';
  role?: 'affected' | 'healthy_control' | 'suspected' | 'excluded' | 'recovered'; // OPTIONAL, see below
  layer?: string; // OPTIONAL: 'canonical' (default when absent) or an overlay id such as 'hypothesis:H17'
  physical?: PhysicalRef;
  mapping?: 'declared' | 'inferred';
  enabled: boolean;
};
type Limitation = {
  policy: string; rowsFiltered: boolean; hiddenColumns: string[];
};
type BindingBase = {
  parameterId: string; value: TypedValue | TypedValueSet;
  origin: 'selection' | 'context' | 'manual' | 'default';
  originEvidence: 'server-default' | 'client-reported';
};
type Binding = BindingBase & (
  | {value: TypedValue; factId?: string; valueFactIds?: never}
  | {value: TypedValueSet; factId?: never; valueFactIds?: string[][]}
);
type TypedValueSet = {
  type: 'set';
  values: TypedValue[]; // 1..500 unique, non-null values of one scalar type
};
type Result = {
  recordset: {
    columns: {name: string; type: string}[];
    rows: TypedValue[][];
  };
  limitations: Limitation[];
  bindingsApplied: Binding[];
  provenance: {
    source: string; collection?: string; queryId?: string;
    mode: 'live' | 'snapshot'; snapshotId?: string; observedAt: string;
    executionProfile: 'protected' | 'opaque-privileged';
  };
  truncated: boolean;
  execution?: ExecutionRef; // present only when the server recorded this run
};
```

Arrays are ordered; rows have exactly one value per returned column. Result values
are typed at the transport boundary; UI adapters may unwrap them for a grid. Missing
and null differ: omit an unbound parameter; null satisfies only a nullable parameter.
False, zero and empty string are present values. No automatic string/number coercion.
Type names for parameters/columns map explicitly to existing core definitions;
unsupported types fail validation rather than being converted silently.

Facts are suggestions, not access credentials. The server revalidates physical
references and mapping declarations against authorized project metadata. Manual
facts are labeled manual and cannot impersonate a selected protected record.

`Fact.condition` (OPTIONAL; lead assumption 2026-09-10, pending founder confirmation)
names the comparison the fact expresses against `Fact.value`: one of `==`, `!=`, `>`,
`>=`, `<`, `<=`, following `investigation-context` REQ:context-variable-conditions.
Absent `condition` means `==`, matching every fact produced before this field existed,
so no existing caller or fixture changes shape. An agent that does not implement
conditions MUST NOT apply a non-`==` fact as equality; it MUST leave the parameter
unbound and report it unbound in the response (`Candidate.missing`, and a `chain`
entry whose `explanation` says the fact was skipped for an unsupported condition),
never silently narrow or misstate the query. This keeps the "why this query is
constrained" view (`investigation-context` REQ:no-hidden-filters) truthful even when
client and server disagree on condition support.

`Fact.role` and `Fact.layer` (OPTIONAL; lead assumption 2026-09-11, pending founder
confirmation) carry the cohort role and the context layer defined by `investigation-context`
REQ:cohort-roles and REQ:context-layers-and-promotion. Absent `layer` means `canonical`.
Neither field changes binding precedence, typing or authorization: a fact in an overlay layer
binds only when the request names that layer, a role is display and cohort metadata, and an
agent that does not implement them MUST ignore `role` and treat a non-canonical `layer` as
"do not bind automatically" rather than binding it as canonical. No existing caller or fixture
changes shape.

`TypedValueSet` is the bounded cohort binding. It is produced only when a caller explicitly
selects a cohort: all enabled facts must name the same `entity.field`, have equality condition
(`condition` absent or `==`), and contain 1..500 non-null values of one compatible scalar type.
The server deduplicates by canonical `TypedValue` equality and sorts the retained values by
their canonical JSON bytes before execution and recording, so retries are byte-identical. A
query parameter must declare set support; its executor maps the set to one source-native
set/array/`IN` binding. A source that cannot bind the set safely returns
`SOURCE_UNAVAILABLE` for that candidate/run. The server MUST NOT issue N queries, interpolate
values into query text, silently keep the first fact, or coerce mixed types. Empty, oversized,
mixed-type, null-containing or non-equality cohorts are `INVALID_REQUEST`.

Limitations report applied restrictions, not the number of rejected rows or their
values. Only names visible in authorized metadata may appear in `hiddenColumns` or
`policy`; use a generic policy label and an empty name list where names are protected.
An empty limitation list does not authorize opaque execution. Denied requests return
no recordset, sample, true count, SQL text or hidden value.

## Endpoint table

| Method/path | Request | Exact success envelope |
|---|---|---|
| GET `agent-info` | None | `{version:string, principal:{id:string,roles:string[],groups:string[]}, securityContextId:string, projects:[{id:string}], capabilities:{protectedQueries:boolean,opaqueReadOnly:boolean}}` |
| GET `semantic/columns` | Scope + SourceRef | `{columns:[{column:string, entity:string, field:string, provenance:'declared'\|'inferred'}]}`; unmapped columns omitted |
| POST `semantic/related` | Scope + `{fact:Fact,limit?:number}` | `{related:[{lookupId:string,label:string,source:string,collection:string,count:number\|null}],truncated:boolean}` |
| POST `semantic/related/rows` | Scope + `{lookupId:string,value:TypedValue,limit?:number}` | `Result` |
| POST `queries/applicable` | Scope + `{values:Fact[]}` | `{applicable:Candidate[],notYet:Candidate[]}` |
| GET `queries/all_queries` | `{project:string, root?:'shared'\|'personal'}` | `{id:string,title?:string,folders:QueriesFolder[],items:QueryDef[]}` (datatug-core's existing recursive shape, unchanged); `root` missing/blank defaults to `'shared'` |
| POST `exec/run_query` | ExecutionRequest below | `Result` |
| POST `queries/capture` | CaptureQueryRequest below | `CaptureQueryResponse` below; `201 Created` for a create, `200 OK` for an update |

Related operations use POST so semantic values are not automatically copied into
URLs, browser history or access logs. Source/project metadata GET queries contain
identifiers only. Server logs record a correlation ID and operation, not raw facts.

`queries/all_queries` lists the project's saved-query tree rooted at `~` (shared; the
default, and byte-identical to the pre-existing response when no `root` param is
present at all) or at the serving principal's own `user:<principalId>` root
(`root=personal`, additive — implemented by datatug-cli PR #230). An anonymous
principal (no `--as` at `datatug serve` startup) gets an empty personal folder, never
an error; a read never creates that directory; the personal tree is never merged into
the shared tree. A `root` value that is present, non-blank and neither `shared` nor
`personal` is `400 INVALID_REQUEST` naming field `root` — the error envelope defined
in Security and errors below. An agent built before this parameter existed ignores
the unrecognized query-string key and always returns the shared tree regardless of
`root`; a client MUST treat any response whose top-level `id` is `~` as the shared
tree even when it requested `root=personal`, never infer `root` support from the
request alone. Where personal queries live on disk (not merely how a client asks for
them) is `executable-knowledge-library` REQ:personal-queries-storage; this endpoint
defines only the wire parameter.

```typescript
type Candidate = {
  queryId: string;
  targets: {source:string;label:string}[]; // authorized eligible targets only
  selectedSource?: string; // present only when one eligible target remains
  bindings: Binding[];
  chain: ({parameterId:string; explanation:string} & (
    | {factId?:string; valueFactIds?:never}
    | {factId?:never; valueFactIds?:string[][]}
  ))[];
  missing: string[]; // parameter IDs, including non-semantic required parameters
  ambiguous: {parameterId:string; factIds:string[]}[];
  state: 'runnable' | 'needs-input' | 'needs-target' | 'source-unavailable';
};
type ExecutionRequest = {
  storeId?: string; project: string; environment: string; securityContextId: string; source?: string;
  queryId?: string; dtql?: string; // exactly one required
  parameters: {[parameterId:string]: TypedValue | TypedValueSet};
  bindingOrigins: ({parameterId:string; origin:Binding['origin']} & (
    | {factId?:string; valueFactIds?:never}
    | {factId?:never; valueFactIds?:string[][]}
  ))[];
  mode: 'live' | 'snapshot'; snapshotId?: string;
  limit?: number;
  incident?: IncidentRef; // {storeId, incidentId}, the incidents Feature's type. OPTIONAL; lead assumption 2026-09-11: a read made under an incident, so incident-scoped grants apply and the execution record carries the incident
};
```

An omitted `ExecutionRequest.storeId` is accepted only when `project` resolves uniquely in the
configured primary project store. The server normalizes it to that explicit store id before
authorization, execution, recording or incident attachment. If more than one store could
match, or the request names a non-primary store only implicitly, it is `INVALID_REQUEST`; no
search-order fallback is allowed.

`applicable` contains only runnable candidates. `notYet` contains authorized query
metadata for needs-input/needs-target/source-unavailable candidates. Multiple eligible targets
produce needs-target, with target options exposed on the candidate; selecting one
passes its source ID to execution. TARGET_REQUIRED errors return the same authorized
target options in error.targets, never hidden source IDs. Unauthorized queries and
targets are omitted entirely; direct execution is independently authorized. Candidates
sort by queryId for deterministic tests. A `queryId` is the saved query's folder-qualified id
(`customers/customer-invoices`), the same id `get_query` and `run_query` accept; a bare
id is accepted on input only when it resolves to exactly one query across folders,
otherwise INVALID_REQUEST names the candidates, and a missing query is NOT_FOUND. Chain explanations distinguish declared,
inferred and manual inputs; a missing parameter still has an explanation.

For ad-hoc DTQL, source is required. For saved queries, source obeys target resolution
above. Unknown parameters and type mismatches are rejected. The execution request
carries ordinary chosen values; the server never silently binds from stored browser
context. It validates defaults too and returns only parameters actually applied by
the executor. `bindingOrigins` supplies display provenance, is checked for exactly
the submitted keys, and never influences authorization. SQL values use driver-bound
arguments; never textual substitution. The returned parameter IDs/values are execution-confirmed; their selection/context/manual
origin remains explicitly client-reported (`originEvidence:client-reported`). Server
default origins use `originEvidence:server-default` only when validated against the
query definition. The UI must not present client origins as server-attested provenance.
A response cannot echo supplied parameter values as applied before execution confirms
them. Suggested Candidate bindings carry the same provenance distinction.

A saved query parameter declares `acceptsSet: true` to accept `TypedValueSet`; absence means
scalar-only. Set values are never accepted for ad-hoc DTQL or an undeclared parameter. For a
set derived from facts, `valueFactIds` has exactly one group for each canonical
`TypedValueSet.values` entry. Each group contains every fact id that contributed that value,
sorted and deduplicated by id; duplicate-valued facts therefore remain attributable without
duplicating the bound value. A scalar selected from facts carries exactly one `factId`.
Manual scalar and set bindings carry neither field, as do server defaults; a manual set is
already fully represented by its canonical values and MUST NOT invent fact provenance.
`Binding`, `Candidate.chain` and `ExecutionRequest.bindingOrigins` use this same shape. The
server rejects an empty fact-id group, a group/value count mismatch, a fact id attached to a
different canonical value, or any object carrying both `factId` and `valueFactIds`.

```typescript
type CaptureQueryRequest = {
  project: string; environment: string; securityContextId: string;
  ifNoneMatch?: true; ifMatch?: string; // exactly one required: create vs. update
  query: CapturedQuery;
};
type CapturedQuery = {
  folderPath: string; // "/"-separated folder under the project's queries root; "" is the root itself
  id: string; title: string; purpose: string;
  source: string; // SourceRef.source the captured lookup read
  dtql: string;
  parameters: CapturedParameter[];
  bindingOrigins: {parameterId: string; origin: 'selection' | 'context' | 'manual'}[];
};
type CapturedParameter = {
  id: string;
  type: 'string' | 'number' | 'integer' | 'decimal' | 'boolean' | 'date' | 'datetime';
  title?: string; isRequired: boolean;
  meta?: {entity: string; field: string}; // present only for a semantic parameter
};
type CaptureQueryResponse = {
  queryId: string; // canonical, folder-qualified id; the id queries/applicable and exec/run_query accept
  revision: string; // opaque; send back as ifMatch to replace exactly this stored revision
  query: CapturedQuery;
  provenance: {author?: string; environment: string; collection: string};
};
```

`queries/capture` ("Save as project query") persists a lookup the browser already ran
as an ordinary project query — a `<id>.query.json` / `<id>.query.dtql` pair left as an
uncommitted Git change — with its purpose, typed parameters and source/binding
provenance. Gated behind the server's write capability (`--allow-writes`); see Security
and errors for the full authorization rule shared with the legacy write routes. The
request flattens `project`, `environment` and `securityContextId` to the top level
alongside exactly one optimistic-concurrency condition: `ifNoneMatch:true` creates only
when nothing is stored at the query's location, `ifMatch:"<revision>"` replaces only
that exact revision. The request carries no result rows, no bound or default parameter
values, no session-local fact IDs, no claimed author and no credentials; strict
decoding rejects any of them, along with duplicate keys and trailing data.
`CapturedParameter` has no default value: a captured default would persist a possibly
protected value into a git-tracked project file. The success envelope is sent only once
the store has durably persisted the pair — `201 Created` for a create, `200 OK` for an
update. `queryId` is the stored query's canonical, folder-qualified ID; `revision` is
opaque and compared for equality only; `provenance` is what the server itself
recorded — `author` (the serving principal, omitted with no identified principal), the
`environment` the source was resolved in, and the `collection` the DTQL reads. Nothing
is committed to Git and no hosted team grant is implied.

## Binding and context behavior

For each parameter: an explicit user edit wins; otherwise one compatible selected
fact wins; otherwise one compatible enabled context fact wins; otherwise a declared
default applies. More than one distinct typed value within the highest eligible
automatic tier is ambiguous and blocks Run. Selection overriding a different context
value is shown as an explicit conflict requiring confirmation; an already-open query
never rebinds silently. Disabled facts are ignored. Defaults never hide an ambiguity.
The user can clear a value; a cleared required value blocks Run until supplied.

Context is isolated by agent URL/project/environment/securityContextId. Switching
scope clears active bindings and opens that scope's own empty or retained local
context; it never imports facts automatically. Sign-out or a changed securityContextId
clears prior principal data. Requests/responses carry the requested scope internally;
late responses from another scope are discarded. A new tab starts without another
tab's context unless the user explicitly imports approved shared state.

## Bounded lookups and HTTP

Default result limit is 100 and maximum is 500. Related discovery returns at most 50
targets; reaching the cap sets `truncated`. Counts are optional: return null if an
exact authorized count cannot be obtained within a 2-second budget. Never compute or
return the unrestricted count. Server execution has a 10-second default timeout,
with a configured upper bound of 30 seconds; cancellation releases the query session.

`lookupId` is an opaque handle to a server-validated relationship, scoped to project,
environment and security context, not authority. Each rows request revalidates the
relationship, value type, current source policy and limit. A forged/stale/cross-scope
handle fails without metadata or row leakage. Do not trust a handle to retain an old
authorization after policy change.

Phase 1 HTTP uses configured public keyless JSON GET descriptors through dalgo2http.
The descriptor fixes the HTTPS host/path and allowed parameters; browser-provided
URLs/headers are rejected. Deny private, loopback, link-local and metadata-service
addresses, including DNS resolution/rebinding and redirects to them. Redirects are
disabled for the demo. Bound response bytes to 2 MiB and validate JSON/row shape.
Provider capabilities must be checked before dispatch: if a requested protected
predicate/projection cannot be enforced safely, reject it rather than fetching an
unrestricted result and claiming enforcement. Never send secrets in query strings.

Live failure returns SOURCE_UNAVAILABLE. The UI may offer a labeled recorded snapshot;
the user must explicitly choose it, producing a second request with `mode:snapshot`
and a configured snapshotId. Snapshot responses retain recorded time and identity,
and pass the same policy boundary. No silent fallback or production-acceptance claim
based on a fixture is allowed.

Snapshot discovery (lead assumption 2026-09-10, pending founder confirmation): a
SOURCE_UNAVAILABLE response for a query with a recorded fixture carries a sibling
top-level `details` key alongside `error` — `{error:{...}, details:{availableSnapshots:
[{snapshotId, recordedAt}]}}`, never nested inside `error` itself. This shape was
forced by an implementation constraint, not chosen freely: datatug-core's
`apicontract.ErrorBody` has no generic extension field, and its own validation
restricts `targets` to TARGET_REQUIRED only, so neither existing field can carry
snapshot metadata without a datatug-core release. `snapshotId` is
`<queryId>@<recordedAt RFC3339>` (e.g. `country-facts@2026-09-09T00:00:00Z`), derived
from the project's own `fixtures/http/<queryId>.json` file, never invented client-side.
No other field belongs in `details`. `datatug serve --http-offline` (lead assumption
2026-09-10, pending founder confirmation) is an operator-facing switch, documented in
`serve --help`, that makes every live HTTP fetch fail SOURCE_UNAVAILABLE without ever
touching the network — for offline demos and for deterministic CI proof of "network
disabled → honest failure → explicit snapshot"; it is not a test-only hook and never
substitutes a snapshot silently.

## Security and errors

Phase 1 serve binds loopback, fixes the principal from `--as` with explicit role/group
flags, and displays it. With a policy set, startup requires an explicit principal;
there is no implicit admin. In a no-policy local project the server may use a clearly
labeled local owner profile. Non-loopback operation requires an authenticated server
profile and is outside this local fixed-principal slice. CORS allowlisting is defense
in depth, not identity. Project writes must be refused unless the server has an
explicit write capability; this read journey must not expose an unauthenticated
mutation endpoint as a side effect.

Every project query write — `queries/capture` and the legacy `queries/create_query`,
`queries/update_query` and `queries/delete_query` — is authorized against the resource
`/datatug_projects/{project}/queries/{queryId}`, deny by default. Without the server's
write capability the write is refused regardless of principal. With it: the
`--no-policies` local owner profile may write unconditionally; otherwise at least one
policy must be loaded, and every loaded policy must allow the write to the serving
principal. A grant that holds only under a row condition, a check or a field list does
not count — a project file has no rows or fields to evaluate it against. A denial names
the deciding policy (when one decided), the operation, the resource and a reason, never
the policy file's server path.

Project-write authorization compares a saved query's id in one canonical spelling on
both sides — the resource being written, and every query id a loaded policy's rule
paths name — full Unicode case folding applied after canonical decomposition, with
every invisible (default-ignorable) character dropped. Two ids that name the same file
on any file system DataTug serves from (`Revenue`, `revenue`, `REVENUE`; NFC and NFD
`café`) always get the same write decision, whichever spelling the policy's rule or the
request's query id uses. A loaded policy that names a rule no query write could ever
be decided by — one scoped below a query id (a folder rule such as
`.../queries/reports/**`, which can only ever match a query literally named `reports`,
never the queries inside a folder named `reports`), one whose nested scopes cannot be
reassembled into `/datatug_projects/{project}/queries/{queryId}`, or one naming a
literal query id no write could ever produce — refuses every project query write while
that policy is loaded, naming the unusable rule in the denial. Reads are unaffected;
they still run through the policy as loaded.

None of that decision's detail reaches the client. `queries/capture` answers a refused
write with 403 ACCESS_DENIED and a fixed message naming neither policy, rule, role nor
resource; the one exception is a refusal no policy decided at all (no configured serve
session, no `--allow-writes`, or no policy loaded), whose message says only which of
those three applies, never anything about the query. The legacy
`queries/create_query`, `queries/update_query` and `queries/delete_query` routes
answer a permission denial the same way, with 403 ACCESS_DENIED and their own fixed
message naming neither policy, rule, role nor resource. On these legacy routes a write
that cannot take the project store's write lock before a bounded deadline ends 504
TIMEOUT, and a write that keeps losing to another write of the same query ends 409
REVISION_CONFLICT once retries are exhausted — both with nothing written; any other
unclassified failure is a 500. These four outcomes share one fixed
`{error:{code,message,requestId}}` body, the message naming no server path, store
text, policy, rule or role, and a `requestId` the agent log carries the full detail
under. An ordinary validation failure on a legacy route (a missing parameter, an
invalid query) stays a plain 400 whose free-text message, not a structured `code` or
`field`, names what was wrong.

A query id and each folder-path segment refuse the same names on every write path,
since all of them become a file or directory name in the project's queries tree: an
invisible (default-ignorable) character such as a zero-width joiner, byte-order mark
or variation selector; an unassigned, private-use or line/paragraph-separator code
point; and a name shaped like a Windows 8.3 short name (`REPORT~1`), which NTFS can
resolve to a different long name than the one it spells. These stand alongside the
existing refusal of control and bidirectional-text control characters, a path
separator, `.` or `..`, a leading `.`, a trailing `.` or space, and the
Windows-reserved device names. `queries/capture` refuses with 400 INVALID_REQUEST
naming `query.id` or `query.folderPath`; the legacy `create_query`, `update_query` and
`delete_query` routes refuse with a plain 400 whose message names the offending
segment.

Every write path also refuses an embedded credential — a password or token inside a
URL or connection-string-shaped value, a secret-named key/value pair (`password=`,
`token=`, `api_key=`, an Azure
`AccountKey=`/`SharedAccessKey=`/`SharedAccessSignature=`, ...), a JSON member such as
`{"password":"..."}`, or a literal HTTP credential header (`Authorization:`,
`X-API-Key:`, ...) — in the query's title, its query text and every parameter's
declared default value, on `queries/capture` and the legacy
`create_query`/`update_query` routes alike. A value that is empty, or a placeholder
rather than a literal secret, is allowed: `?`, a positional bind such as `$1`, a named
bind such as `:new_password`, a bare or braced environment reference such as
`$PGPASSWORD` or `${PGPASSWORD}`, a double-brace template such as `{{token}}`, or an
angle-bracket placeholder such as `<password>`. `queries/capture` refuses with 400
INVALID_REQUEST naming `query.<field>`; the legacy routes refuse with a plain 400
whose message names the field.

A query's capture provenance — that it was saved through `queries/capture`, and by
whom, in which environment, reading which collection, with which binding origins — is
the server's own record on the legacy routes too, never a client's claim.
`create_query` refuses a request whose query carries a capture block at all, with a
400 whose message names the field. `update_query` accepts one only when it is exactly
the stored query's own (its bindings compared regardless of order); a capture block
that differs, or that appears where the stored query has none, is refused the same
way, with 400. Omitting the block on `update_query` keeps the stored capture unchanged
— omission never clears it. Because the decision is made against whatever is currently
stored, a write that raced another write to the same query retries against the freshly
stored value, and can still end in 409 REVISION_CONFLICT with nothing written once
retries are exhausted.

A folder-qualified write on the legacy `create_query`, `update_query` or
`delete_query` routes (a `folderPath` other than the queries root, or an id containing
`/`) is refused with 400 against the currently released datatug-core: its project
store writes and deletes every query at the queries root regardless of the folder a
request names, so authorizing the folder the caller asked for would authorize a
decision the store does not act on, silently touching the root query of the same name
instead. The legacy routes accept and delete root queries only in this build. This is
a version-dependent limitation, not a shape the contract endorses: it is lifted, with
no wire-format change, once the CLI is built against a datatug-core release whose
project store honours `folderPath`.

Protected execution refuses opaque native SQL before dispatch. A separate explicit
opaque-query grant may allow read-only native SQL for a privileged principal, but
the actual native SQL execution result must say `opaque-privileged` and the UI must state row/column enforcement
does not apply. The support demo has no such grant. Browsers cannot choose or upgrade
profiles. agent-info reports allowed capabilities; Result.provenance.executionProfile
reports the actual executor used for that run. A privileged principal running protected
DTQL still receives protected provenance. All legacy routes obey the same rule.

`queries/capture` refuses in this order, writing nothing before the store is reached:
missing `project`, `environment` or `securityContextId` (400 MISSING_PARAMETER); a
`securityContextId` that does not match the current session (409 STALE_CONTEXT); not
exactly one of `ifNoneMatch`/`ifMatch` (400 INVALID_REQUEST on `ifNoneMatch/ifMatch`);
an invalid `query` — unsafe folder/id/source name, oversized or control-character
title/purpose, a DTQL that fails to parse or does not read exactly its declared typed
parameters, an unrecognized binding origin, or an embedded credential (400
INVALID_REQUEST naming `query.<field>`); an unknown project (404 NOT_FOUND); the
serving principal not authorized for this write (403 ACCESS_DENIED, the rule above); a
`source` that is not a catalog source of `environment` a saved query can target (400
INVALID_REQUEST on `query.source`). Only then is the conditional store write attempted,
bounded by a fixed write deadline; it fails 400 INVALID_REQUEST for an unsafe stored
location or an incomplete stored pair, 409 REVISION_CONFLICT for a failed
`ifNoneMatch`/`ifMatch` condition, or 504 TIMEOUT — nothing written — when the store's
write lock cannot be taken before the deadline. Any other store failure is a 500 whose
message names no server path or store text and carries the response's own `requestId`,
outside this appendix's closed code set.

Errors use `{error:{code:string,message:string,field?:string,requestId:string,targets?:{source:string,label:string}[]}}`.
Only TARGET_REQUIRED may include authorized target options.
HTTP 400 covers INVALID_REQUEST, TYPE_MISMATCH, MISSING_PARAMETER, AMBIGUOUS_BINDING
and TARGET_REQUIRED; 401 UNAUTHENTICATED; 403 ACCESS_DENIED or
UNSUPPORTED_PROTECTED_EXECUTION; 404 NOT_FOUND (including invisible resources);
409 STALE_CONTEXT or REVISION_CONFLICT; 413 RESPONSE_TOO_LARGE; 503 SOURCE_UNAVAILABLE;
504 TIMEOUT.
Messages and IDs must not reveal protected values or credentials. Tests assert both
status and code, not English wording.

REVISION_CONFLICT (409) refuses a project write whose optimistic-concurrency condition
failed: an `ifNoneMatch` create found a query already stored at the location (field
`query.id`), or an `ifMatch` update named a revision that is stale or no longer exists
(field `ifMatch`). Nothing is written. Recovery is to reload the query's current
revision, then retry with `ifMatch`. It is distinct from STALE_CONTEXT: STALE_CONTEXT
means the caller's `securityContextId` is itself stale and is recovered by calling
`agent-info` again — retrying with a fresher revision can never fix a stale context,
and calling `agent-info` again can never fix a failed write condition.

## Known gaps

- Reading a query's current revision: not yet provided. No endpoint returns it, so a
  REVISION_CONFLICT recovery on `ifMatch` must reload the query through another route
  to learn its current revision before retrying.
- Whether writes are allowed: not yet provided. `agent-info`'s `capabilities` do not
  say whether the server was started with `--allow-writes`, so a client cannot tell
  before attempting a write.
- A stored query file that does not parse: `create_query` and `update_query` read the
  stored query first, to decide its capture provenance, and a parse failure there is
  indistinguishable from any other store failure, so it answers 500 like any other
  unclassified failure, not a distinct 4xx.
- A query id the newer name rules refuse (one stored before those rules took effect)
  can still be read; it cannot be updated or deleted through the legacy routes, which
  apply the same id rule to every write.

## Acceptance and migration

The implementation task must freeze schema fixtures beside the owning code contract
and consume the same fixtures on server and client. Fixtures test shapes, not fake
browser acceptance. Cover empty arrays, null/false/zero/large integer, ambiguous and
missing inputs, omitted unauthorized metadata, HTTP failures, all envelopes and errors.
Then call the real server from the real client and prove actual parameter effects.

Inventory existing callers of renamed methods and payloads. Cut over them all in the
same candidate; delete superseded client adapters. If an external compatibility
consumer is found, record its concrete rollout choice in the owning Feature before
retaining a dual protocol. Do not reinterpret whichever implementation currently
exists as an alternative normative contract.

## Open Questions

None for this transport slice. Numeric bounds above are conservative engineering
defaults, configurable within documented server policy; they are not founder rulings.
