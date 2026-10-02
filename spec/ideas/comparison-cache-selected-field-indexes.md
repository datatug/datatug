---
format: https://specscore.md/idea-specification
status: Approved
---

# Idea: Selected-field indexes for local comparison caches

**Status:** Approved
**Phase:** FUTURE
**Date:** 2026-09-14
**Owner:** alex
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we accelerate explicitly requested field filters over schemaless local comparison-cache rows without turning arbitrary source columns into the fixed SQLite envelope?

## Context

Task 15 caches post-policy comparison rows as canonical minified JSON TEXT in a schema-versioned local SQLite envelope indexed only by comparison id, state and canonical sort key. Some future large-cache journeys may need a user-selected field filter.

## Recommended Direction

On explicit user or project configuration, lazily extract only the selected JSON fields into separately governed index structures. Keep the JSON payload authoritative for the cache row, version migrations, and re-check current authorization before every read.

## Alternatives Considered

- Index arbitrary JSON paths up front — rejected because source schemas are unbounded and most
  cached comparisons will never use a field filter.
- Add dynamic source columns to `comparison_rows` — rejected because it couples the cache envelope
  and its migrations to every compared source schema.
- Use a columnar spill format — rejected until measured cache size and scan latency show that the
  simpler SQLite JSON scan plus selected indexes is inadequate.

## MVP Scope

After measured demand, support one explicitly selected typed field, build and remove its local index lazily, and prove identical filtered results against a full JSON scan.

## Not Doing (and Why)

- Arbitrary JSON-field indexes — unbounded storage and migration cost without demand
- Dynamic source columns in comparison_rows — breaks the fixed cross-source envelope
- Compression, SQLite JSONB, codec negotiation or columnar spill — unrelated to selected-field lookup

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Users need repeated filtering on a small, explicitly chosen set of fields | Measure Task 15 cache-page filter requests and require repeated use before building an index |
| Should-be-true | A typed selected-field index materially improves filtered page latency over a full JSON scan | Benchmark representative large caches with and without the index |
| Might-be-true | Lazy local index creation and eviction fit the existing cache lifecycle | Prototype one field and verify schema-version migration, authorization and eviction tests |


## SpecScore Integration

- **New Features this would create:** TBD at design time
- **Existing Features affected:** none
- **Dependencies:** local-comparison-cache

## Open Questions

None at this time.
