---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Infrastructure inventory

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/infrastructure-inventory?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/infrastructure-inventory?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/infrastructure-inventory?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/infrastructure-inventory?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

P13 is infrastructure inventory: environments, servers, databases, endpoints and services linked to sources, so an investigation can climb from a record to the infrastructure that produced it — "this order came from `prod`'s `orders-service` on `db-primary-3`." It is a later-phase extension of the navigation loop P5 establishes for data records, applied one level up to the systems that host the data.

**Evidence boundary:** initial-audit observations below are historical. Current integration
status is summarized in the [hierarchy](../../product/hierarchy.md) and the
[independent review](../../research/2026-09-09-independent-review.md); no historical
state label is a present completion or absence claim. The [roadmap](../../../ROADMAP.md)
owns phase sequencing; capability tables below retain original audit grouping.

Initial audit snapshot (before later 9 September landings): environments, servers and DB catalogs are usable and scanned (`environment.go`, `server.go`, `dbcatalog.go`, populated by `datatug scan`); services and endpoints are "prototype (vestigial `ApiService`/`ApiEndpoint`)," not persisted by any store; the climb from a record to its infrastructure does not exist at all. Phase 5 (experiments) is where services/ endpoints and the climb are planned — this pillar deliberately trails the core loop rather than competing with it for Phase 1 attention. This is a Capability; per-repo Implementations (`datatug-core`, `datatug-cli`) link to it.

## Problem

The following initial-audit diagnosis explains why this capability was proposed.
It is historical evidence, not a claim about current code; consult the capability
hierarchy and revised Phase 1 plan for current integration work.

> Once P5 lets a user follow a `Customer.ID` to its invoices and support notes, the natural next
> question in a real troubleshooting session is infrastructural: which environment, which server,
> which service produced this record, and is that service healthy. DataTug already models and scans
> environments, servers and catalogs (`environment.go`, `server.go`, `dbcatalog.go`; `pkg/schemer`),
> so a source is never disconnected from its hosting context in the data model — but the audit found
> the services/endpoints layer is "vestigial" (`ApiService`/`ApiEndpoint` types exist, "not persisted
> by any store") and there is no code anywhere that climbs from a record to the server or service that
> produced it. This pillar is explicitly a Phase 5 experiment, not a Phase 1 gap: the core
> investigation loop does not depend on it, and building it early would pull attention from the
> pillars the demo needs to prove first.

## Behavior

### Initial audit capability inventory (historical)

| Capability | State | Phase | Kind |
|---|---|---|---|
| Environments, servers, catalogs | usable | — | enabler |
| Services and endpoints linked to sources | prototype (vestigial `ApiService`/`ApiEndpoint`) | 5 | strategic |
| Climb from a record to its infrastructure | planned | 5 | strategic |

#### REQ: environments-servers-catalogs

An environment, its servers, and each server's database catalogs MUST remain scannable and
inspectable as they are today, with `datatug scan` as the source of truth for what a given
environment actually contains. This capability is already usable; this REQ pins it as the standing
foundation the rest of P13 builds on.

#### REQ: services-endpoints-inventory

A source SHOULD be linkable to the service and endpoint that expose it (an HTTP API, a specific
database instance behind a load balancer), persisted as a real project asset rather than the
current vestigial, unpersisted `ApiService`/`ApiEndpoint` types. This inventory MUST stay
inspectable data — DataTug does not monitor service health, it records the relationship a human
declares.

#### REQ: climb-to-infrastructure

Given a record (an order, a support note), the user SHOULD be able to climb from it to its source,
then to that source's service, environment and server, in one guided path — "order → source →
service → environment → server." This is an extension of P5's related-data navigation one level up
the stack; it MUST reuse that navigation pattern rather than introduce a separate mechanism.

## Acceptance Criteria

### AC: environment-scan-lists-servers-and-catalogs (verifies REQ:environments-servers-catalogs)

**Given** the demo project's configured environments
**When** `datatug scan` runs
**Then** each environment lists its servers and each server lists its database catalogs, matching what is actually reachable.

### AC: service-endpoint-recorded-and-linked (verifies REQ:services-endpoints-inventory)

**Given** a source backing the `chinook` database
**When** a service and endpoint are declared for it
**Then** the link is persisted as a project file and survives a project reload, unlike today's unpersisted `ApiService`/`ApiEndpoint` types.

### AC: climb-from-record-to-server (verifies REQ:climb-to-infrastructure)

**Given** an `Invoice` record from the `chinook` source
**When** the user chooses "show infrastructure" on it
**Then** the path shown is `Invoice → chinook source → (service, if declared) → environment → server`, using the same navigation interaction pattern as a P5 related-record lookup.

## Not Doing / Out of Scope

- Live infrastructure monitoring (health checks, uptime, alerting): DataTug records the declared
  relationship between a source and its infrastructure, it does not become an observability tool.
- A generic CMDB or asset-management product: the inventory exists to let an investigation climb
  from a record to its context, not to catalog infrastructure for its own sake.
- Automatic infrastructure discovery from network scanning; links between a source and its service/
  server are declared, not inferred from traffic.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
