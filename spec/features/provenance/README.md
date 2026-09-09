---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Provenance

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/provenance?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/provenance?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/provenance?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/provenance?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P12 is provenance: every value and every mapping can say where it came from — source, table or endpoint, field, query, transformation, mapping and any policy limitation applied.

It is what makes DataTug's inspectability claims verifiable rather than asserted: a "declared" mapping or a "limited by policy" result only means something if the user can trace exactly why. Audited state on 2026-09-09: not implemented beyond schema docs, though two seeds exist — Phase 1's `limitation-visible` requirement (every result carries policy limitations) and the layered-ACL design's Explain-Access reason-code/`blockers[]` taxonomy, which is specified but not shipped; the CLI's own current access-denial explanation is "implemented-but-poorly-exposed," deriving attribution by splitting an English sentence on `" via "` rather than returning structured data (`pkg/accesspolicies/explain.go`) — a gap the layered-ACL design packet itself flags. Phase 2 builds result provenance on top of Phase 1's limitation seed; mapping provenance ties into P4's review/ correct workflow. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`, `dal-go/dalgo`) link to it.

## Problem

"Where did this value come from" and "why can't I see this column" are the two questions a
troubleshooting tool must answer precisely, or the user stops trusting the tool the first time an
answer feels arbitrary. DataTug already produces the raw material for both answers — every result
Phase 1 returns carries `limitations` (policy name, rows filtered, hidden columns; `core-
investigation-loop` REQ:limitation-visible) and the access-policy evaluator already computes a
decision — but the explanation itself is currently "derived by splitting the English
`Decision.Explanation` string on `' via '`" (`pkg/accesspolicies/explain.go`), which the layered-ACL
design's own review flagged as a gap to close with structured reason codes
(`vm/layered-acl-design` `04-authorization-contract.md`, specified/documented, not shipped). Mapping
provenance has even less today: nothing records who or what declared or inferred a `Customer.ID`
mapping, so a wrong mapping cannot be traced back to its source once it has propagated through the
library.

## Behavior

### Capabilities and phases

| Capability | State | Phase | Kind |
|---|---|---|---|
| Result provenance: source, table/collection/endpoint, query, policy limitations | planned | 2 | strategic |
| Mapping provenance: who/what declared or inferred a mapping | planned | 2 | strategic |

#### REQ: result-provenance

Every result MUST be able to state, on demand, its full provenance chain: the source, the table/
collection/endpoint, the query or lookup that produced it, any transformation applied, and any
policy limitation. This builds directly on `core-investigation-loop` REQ:limitation-visible, which
already requires the policy-limitation half of this chain; this Feature extends it to the source/
query/transformation half that Phase 1 does not cover.

#### REQ: mapping-provenance

Every semantic mapping (declared, inferred, corrected) MUST record who or what produced it and
when. This is the record `meaning-inference` REQ:review-accept-correct-provenance writes to; this
Feature owns the provenance record's shape and guarantees it is queryable — "show me every mapping
this user corrected" MUST be answerable without grepping commit history by hand.

#### REQ: structured-access-explanation

An access-policy decision MUST be explainable with structured reason codes (matching the intent of
the layered-ACL design's Explain-Access `blockers[]` taxonomy), not by parsing an English sentence.
This Feature does not redesign that taxonomy — it depends on the layered-ACL track (see
`server-enforced-access` REQ:layered-acl-track) landing task 21 before result provenance can show a
fully structured "why was I refused" explanation instead of today's string-derived approximation.

## Acceptance Criteria

### AC: result-shows-full-chain (verifies REQ:result-provenance)

**Given** a related-record lookup for `Customer.ID = 5` returning `Invoices` from `chinook.sqlite`
**When** the user asks "where did this come from"
**Then** the response names the source (`chinook`), the collection (`Invoice`), the lookup that produced it, and any policy limitation applied — the same shape `core-investigation-loop` AC:restricted-rows-and-columns-server already returns for limitations, extended to source/query.

### AC: mapping-shows-declarer-and-time (verifies REQ:mapping-provenance)

**Given** `Customer.ID → chinook.Customer.CustomerId` was accepted from an inferred suggestion
**When** the mapping's provenance is inspected
**Then** it names the user who accepted it and the timestamp, matching `meaning-inference` AC:accept-promotes-to-declared.

### AC: access-refusal-carries-reason-code (verifies REQ:structured-access-explanation)

**Given** the `support` principal requests `Customer.Email`
**When** the server refuses the request
**Then** the refusal carries a structured reason (policy name, blocker type), not only an English sentence a client would need to parse — available once the layered-ACL task 21 lands; today's `ACCESS_DENIED` response (`core-investigation-loop` AC:hidden-column-refused) is the interim, less structured form.

## Not Doing / Out of Scope

- A data catalog's lineage diagram as a governance artifact: provenance exists to answer "why" and
  "where from" during an investigation, not to produce a certified lineage report.
- Redesigning the Explain-Access reason-code taxonomy: that belongs to the layered-ACL plan this
  Feature consumes once its integration task lands.
- Provenance for values never read through DataTug's server (a value pasted in manually, for
  example) — provenance only covers what the server itself produced or observed.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
