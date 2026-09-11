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

`Scope = {project: string, environment: string, securityContextId: string}`. Project and environment IDs are
required, nonempty, and resolved within this server's registered projects. A
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
type Binding = {
  parameterId: string; value: TypedValue;
  origin: 'selection' | 'context' | 'manual' | 'default';
  originEvidence: 'server-default' | 'client-reported';
  factId?: string;
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
| POST `exec/run_query` | ExecutionRequest below | `Result` |

Related operations use POST so semantic values are not automatically copied into
URLs, browser history or access logs. Source/project metadata GET queries contain
identifiers only. Server logs record a correlation ID and operation, not raw facts.

```typescript
type Candidate = {
  queryId: string;
  targets: {source:string;label:string}[]; // authorized eligible targets only
  selectedSource?: string; // present only when one eligible target remains
  bindings: Binding[];
  chain: {parameterId:string; factId?:string; explanation:string}[];
  missing: string[]; // parameter IDs, including non-semantic required parameters
  ambiguous: {parameterId:string; factIds:string[]}[];
  state: 'runnable' | 'needs-input' | 'needs-target' | 'source-unavailable';
};
type ExecutionRequest = {
  project: string; environment: string; securityContextId: string; source?: string;
  queryId?: string; dtql?: string; // exactly one required
  parameters: {[parameterId:string]: TypedValue};
  bindingOrigins: {parameterId:string; origin:Binding['origin']; factId?:string}[];
  mode: 'live' | 'snapshot'; snapshotId?: string;
  limit?: number;
};
```

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

Protected execution refuses opaque native SQL before dispatch. A separate explicit
opaque-query grant may allow read-only native SQL for a privileged principal, but
the actual native SQL execution result must say `opaque-privileged` and the UI must state row/column enforcement
does not apply. The support demo has no such grant. Browsers cannot choose or upgrade
profiles. agent-info reports allowed capabilities; Result.provenance.executionProfile
reports the actual executor used for that run. A privileged principal running protected
DTQL still receives protected provenance. All legacy routes obey the same rule.

Errors use `{error:{code:string,message:string,field?:string,requestId:string,targets?:{source:string,label:string}[]}}`.
Only TARGET_REQUIRED may include authorized target options.
HTTP 400 covers INVALID_REQUEST, TYPE_MISMATCH, MISSING_PARAMETER, AMBIGUOUS_BINDING
and TARGET_REQUIRED; 401 UNAUTHENTICATED; 403 ACCESS_DENIED or
UNSUPPORTED_PROTECTED_EXECUTION; 404 NOT_FOUND (including invisible resources);
409 STALE_CONTEXT; 413 RESPONSE_TOO_LARGE; 503 SOURCE_UNAVAILABLE; 504 TIMEOUT.
Messages and IDs must not reveal protected values or credentials. Tests assert both
status and code, not English wording.

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
