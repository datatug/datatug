---
format: https://specscore.md/feature-specification
status: Approved
---
# Feature: Annotations

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/annotations?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/annotations?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/annotations?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/annotations?op=request-change) |
**Status:** Approved
**Source Ideas:** —

## Summary

An annotation is a note attached to the data it is actually about — an entity value, an
execution record, later a row set or a relationship — so that what one investigator learned
resurfaces automatically the next time that same data appears, in any authorized query.

An annotation carries its own provenance (who, when, in which incident, from which execution
and bindings) and it is visible only to someone who can see the data it annotates.

## Problem

Knowledge found during an investigation is written in the wrong place. "Order 72841 retried
correctly — probably not the retry bug" lands in a chat thread, a ticket comment or nobody's
notes, addressed by nothing. The next person hits the same order through a different query and
starts from zero. Nothing exists in any DataTug repository today (reconciliation §1,
"Annotations — absent").

The founder's requirement (vision §16) is annotation "at the level where evidence actually
exists", with a directional rule:

> Prefer **data-addressable** annotations where possible. If an order is annotated through one
> query and appears later through a different query, the relevant annotation should be
> discoverable again.

Data-addressable is also the hard part: an annotation on `Customer.email = a@b.com` is itself a
disclosure of a value. Resurfacing must therefore ride on the same policy evaluator that
decides whether the value may be shown at all.

```text
      grid cell: Order.ID = 72841
              │ "Annotate"
              ▼
      Annotation { target: entity-value(Order.ID = 72841), body, provenance }
              │
   stored as a project asset through the server (revisioned, no rows, no secrets)
              │
              ├─► the SAME value appears in a DIFFERENT query's result later
              │         └─► marker in the grid + entry in the context panel
              │
              └─► visibility = "may this reader see Order.ID = 72841 right now?"
                        no  ─►  the annotation does not exist for them
```

## Behavior

### Capabilities and phases

| Capability | State | Tier | Kind |
|---|---|---|---|
| Annotate an entity value (`entity.field = value`) | planned | MVP (B1) | strategic |
| Annotate an execution record | planned | MVP (B1) | strategic |
| Provenance captured automatically | planned | MVP (B1) | strategic |
| Resurfacing in any authorized result (grid marker + context panel) | planned | MVP (B1) | strategic |
| Stored as a revisioned project asset through the server | planned | MVP (B1) | enabler |
| Visibility follows the target's current policy | planned | MVP (B1) | strategic |
| `datatug annotation` CLI and `/datatug/annotations` API | planned | MVP (B1) | surface |
| Annotate a snapshot row by key facts | planned | NEXT | strategic |
| Annotate a row set (record + selection) | planned | NEXT | strategic |
| Annotate a field, a relationship or a query | planned | NEXT | strategic |
| Annotation shown on a board card | planned | NEXT | surface |
| AI-suggested annotations from an incident's events | planned | FUTURE | strategic |

#### REQ: annotation-targets

An annotation MUST address exactly one target. The target kinds are:

| Kind | Address | Tier |
|---|---|---|
| `entity-value` | `{entity, field, value: TypedValue}` — the `entity.field = value` form | MVP |
| `execution-record` | `{recordId}` | MVP |
| `snapshot-row` | `{recordId, key: Fact[]}` — a row identified by its key facts, not by row index | NEXT |
| `row-set` | `{recordId, selection}` — a record plus an explicit selection | NEXT |
| `field` | `{entity, field}` — the field itself, no value | NEXT |
| `relationship` | `{lookupId-equivalent declared relationship}` | NEXT |
| `query` | `{queryId}` | NEXT |

MVP MUST implement `entity-value` and `execution-record` only; the later kinds MUST be
representable in the stored shape from the start (a `target.kind` discriminator) so adding one
is not a migration. A target MUST be semantic, never physical: an annotation addresses
`Order.ID = 72841`, not `orders.id` in `prod`, so it survives a mapping change and resurfaces
across sources exactly as `semantic-layer` mappings intend. A `snapshot-row` target MUST be
keyed by facts rather than by ordinal position, because a snapshot's row order is not a stable
identity.

#### REQ: provenance-auto-captured

An annotation MUST capture its provenance automatically, with no form for the author to fill
in: author (the current server-resolved principal, never a client-supplied name), creation
time, project and environment scope, the execution record it was made from, the incident it
was made under if any, the selected facts active at the time, and the query id and parameters
of that execution. This is vision §16's list ("incident, author, timestamp, query, parameters,
source, snapshot and selected entities/values") reduced to references: provenance MUST be
stored as ids into [evidence-records](../evidence-records/README.md) and
[incidents](../incidents/README.md), not as a copy of their contents.

```typescript
type Annotation = {
  id: string;
  target: {kind: 'entity-value'|'execution-record'|'snapshot-row'|'row-set'|'field'|'relationship'|'query'; /* address per REQ:annotation-targets */};
  body: string;                       // plain text, human-authored
  assertion: 'claim' | 'question';    // see REQ:annotation-is-a-claim-not-an-observation
  provenance: {
    author: {id: string};
    createdAt: string;                // RFC3339, UTC, server clock
    scope: {project: string; environment: string};
    recordId?: string;
    incidentId?: string;
    facts: Fact[];                    // the active Investigation Context at authoring time
    queryId?: string;
    parameters?: {[parameterId: string]: TypedValue};
  };
  revisions: {at: string; author: {id: string}}[];
};
```

`facts` MUST be the enabled facts of the authoring context as `investigation-context` defines
them, so a reader can see the question that was being asked, not only the answer that was
annotated. The server MUST reject a client-supplied `author`, `createdAt` or `scope`, the same
way the transport contract rejects a client-supplied principal.

#### REQ: data-addressable-resurfacing

An annotation on `entity.field = value` MUST be surfaced wherever that value appears in any
result the reader is authorized to see — not only in the query it was created from. Two
surfaces are required in MVP: a **marker in the grid** on the annotated cell, and an **entry in
the context panel** listing the annotations relevant to the current result.

Matching MUST be typed equality on the semantic address, using the same `TypedValue` semantics
the binding contract already fixes (`core-investigation-loop/api-contract.md`): no
string/number coercion, `false`/`0`/empty string are values, and null matches only null.
Resolution MUST run server-side over the result the server just produced, so no client ever
asks "which annotations exist for this value?" with a value the server did not authorize it to
hold. A result with no annotated values MUST return an empty list rather than a suppressed
one, so an empty context panel is honest.

Resurfacing MUST NOT change the query, the bindings or the row set. An annotation is a marker
on data, never a filter — the same rule `investigation-context` REQ:no-hidden-filters holds for
context values.

#### REQ: annotation-storage

Annotations MUST be project assets under `annotations/`, written through the server like every
other project write, readable and diffable as files (`git-native-projects`
REQ:readable-file-layout). An annotation file MUST NOT contain result rows, snapshot contents,
query text or credentials; it holds a target address, a body, and references. Editing an
annotation MUST create a revision rather than replacing history: `revisions` grows, prior
bodies remain in the file's Git history, and a deletion is a tombstone the timeline can still
render.

File layout (lead assumption): `annotations/<id>.json`, one annotation per file, so two
investigators annotating different values in the same minute do not collide on one file. This
depends on the Phase 2 project-write path; until it exists, `POST /datatug/annotations` MUST
fail with an explicit "project is read-only" error rather than holding the annotation in
memory.

#### REQ: annotation-visibility-follows-target

A reader who cannot see the target value under the **current** policy MUST NOT see the
annotation, its body, its author or its existence. Visibility is evaluated per read, for the
current principal, by the same evaluator that authorizes the data
(`server-enforced-access` REQ:dalgo-access-policies-external) — an annotation is never granted
separately and confers no access to anything.

Concretely: an annotation whose target is `entity-value` is visible only if the reader could
retrieve that value from at least one source they are authorized to read; an annotation whose
target is `execution-record` is visible only if the reader may read that record under
`evidence-records` REQ:snapshot-access-rechecked. Beyond the target check, an annotation is a
project asset and follows the project's grants to spaces (`sharing-scopes`
REQ:project-membership-through-sneat-spaces) — there is no annotation-level owner or ACL.

Denials MUST be non-disclosing: no count of hidden annotations, no author names, no "3 notes
you cannot see" affordance.

#### REQ: annotation-is-a-claim-not-an-observation

An annotation body is a human statement about data, not a measurement of it. When an
annotation is made under an incident it MUST enter the incident event stream with
`assertion.kind = claim` — or `question`, a lead assumption adding the "Is this actually
limited to enterprise customers?" case from vision §16 — never `observation` or
`deterministic-result`, per
the family's assertion-kind rule (product-family §5, D12; vision §12, §31). An AI-generated
annotation, when that arrives, MUST be `inference` and MUST carry the events it was derived
from. The measurement it comments on is the execution record; the annotation is the opinion.

#### REQ: api-and-cli

The server MUST expose, under `/datatug`, following the existing scope, staleness and error
envelope rules:

| Method/path | Request | Success envelope |
|---|---|---|
| POST `annotations` | Scope + `{target, body, assertion?, incidentId?, recordId?, facts?}` | `{annotation: Annotation}` |
| GET `annotations` | Scope + one of `target`, `recordId`, `incidentId`, plus `limit` | `{annotations: Annotation[], truncated: boolean}` |
| POST `annotations/for_result` | Scope + `{recordId}` or the just-returned result identity | `{annotations: Annotation[]}` — the resurfacing lookup |

The resurfacing lookup is a POST so semantic values are never copied into URLs, browser history
or access logs, exactly as `semantic/related` already is. The CLI follows the singular-namespace
convention: `datatug annotation add --entity Order --field ID --value 72841 --body "…"` and
`datatug annotation list [--record <id>] [--incident <id>] [--json]`. CLI and web use the same
endpoints; neither reads or writes an annotation file directly.

## Acceptance Criteria

### AC: annotate-value-resurfaces-in-other-query (verifies REQ:data-addressable-resurfacing)

**Given** `Order.ID = 72841` annotated with "Retry happened correctly; probably not the retry bug" from the *Stuck orders* query
**When** the same authorized user later opens *Orders by customer*, a different query whose result also contains `Order.ID = 72841`
**Then** the cell carries an annotation marker and the context panel lists that annotation with its author and time, without the user searching for it and without the second query's rows or bindings changing.

### AC: annotation-carries-provenance (verifies REQ:provenance-auto-captured)

**Given** an investigator working under incident `INC-2841` with `Customer.ID = 5` active in the Investigation Context and a recorded execution of *Stuck orders*
**When** they annotate a value in that result and supply only the body text
**Then** the stored annotation carries the server-resolved author, server timestamp, project/environment scope, the execution record id, `INC-2841`, the active facts and the query id with its parameters — and a request that tries to set `author` or `createdAt` itself is rejected.

### AC: annotation-hidden-when-target-denied (verifies REQ:annotation-visibility-follows-target)

**Given** an annotation on `Customer.email = a@b.com` and the `support` principal whose policy hides the `email` column on every source that maps it
**When** `support` reads the annotated result and calls the annotation list and resurfacing endpoints
**Then** no annotation is returned, no count, author or placeholder is shown, and the responses are indistinguishable from a project with no annotation on that value.

### AC: annotation-on-execution-record (verifies REQ:annotation-targets, REQ:annotation-is-a-claim-not-an-observation)

**Given** execution record `exec-7f3` captured during `INC-2841`
**When** a participant annotates the record itself with "This is the 14:52 peak measurement"
**Then** the annotation's target is `{kind: 'execution-record', recordId: 'exec-7f3'}`, it appears when the record is opened, and the incident event it produces has `assertion.kind = claim`, not `observation`.

### AC: annotation-never-stores-restricted-rows (verifies REQ:annotation-storage)

**Given** an annotation created from a result containing columns the project masks
**When** the written `annotations/<id>.json` file is inspected on disk
**Then** it contains a target address, body and references only — no result rows, no snapshot content, no query text and no credentials — and editing the body appends a revision entry instead of replacing the previous one.

## Not Doing / Out of Scope

- A discussion system: threads, mentions, reactions, notifications and assignment belong to the
  incident's participants and event stream, not to a note on a value.
- A second knowledge store: an annotation is a project file written through the server, beside
  queries, boards and incidents — not a database of its own.
- Annotation-level permissions, sharing links or per-annotation owners. Visibility is derived
  from the target and the project's grants; there is nothing to configure.
- Rewriting or enriching the annotated data. An annotation never changes a value, a query, a
  binding or a row set.
- Free-text search across annotation bodies in MVP; incident search
  ([incidents](../incidents/README.md)) owns discovery, and deterministic fact overlap comes
  first.
- Automatic extraction of annotations from chat, email or voice (vision §26–§28) — those are
  capture channels into the incident stream, specified elsewhere.

## Open Questions

- **Cross-environment resurfacing.** Should an annotation made on `Order.ID = 72841` in `uat`
  surface on the same id in `prod`? Lead assumption: no — provenance records the environment and
  resurfacing is environment-scoped by default, because the same key can mean different records
  in different environments. A per-entity "ids are global" declaration would relax it; no
  evidence yet that it is needed.
- **High-cardinality resurfacing cost.** Matching every cell of a 500-row result against every
  annotation is cheap at demo scale and unbounded at project scale. Lead assumption: index by
  `(entity, field, canonical value)` server-side and cap the lookup at the same 500-row result
  limit; revisit if a project ever holds enough annotations to make that visible.
- **Annotating a value the reader can see but the author could not** (and the reverse, after a
  policy change): the annotation's provenance says the author saw something the reader cannot
  reconstruct. MVP shows the annotation when the target is visible now and says nothing about
  the author's access then. Confirm that is acceptable before hosted multi-user.

---
*This document follows the https://specscore.md/feature-specification*
