# DataTug / Incidentius / Dashboardius

## Incident Resolution as a Living Organisational Knowledge System

**Status:** Product vision, architecture, and engineering handoff\
**Primary MVP target:** the existing DataTug application/codebase\
**Product profiles:** DataTug · Incidentius · Dashboardius\
**Incidentius positioning:** **An incident-resolution workspace**

------------------------------------------------------------------------

# 1. Executive summary

Incidentius is **not an incident-management platform**. It should not
try to replace paging, on-call, ticketing, status-page, escalation, or
ITSM products.

Those systems can answer:

> Who should respond, and what process should we follow?

Incidentius answers:

> **What is happening, what have we learned, what should we investigate
> next, and how do we resolve it?**

Its job starts when somebody says:

> **"Houston, we've got a problem."**

Its purpose is to help troubleshooters resolve incidents faster by
turning collaborative investigation into **shared, structured,
executable knowledge**.

The deeper platform thesis is:

> **Capture technical reasoning while people work, convert it into
> reusable organisational knowledge, and bring the right knowledge back
> when and where it matters.**

DataTug is the technical backbone. Dashboardius is the
dashboard/visualisation experience. Incidentius is the
incident-resolution experience. They should be different product
profiles of one application, not separately implemented products.

``` text
                 ONE APPLICATION / REPOSITORY
                         │
          shared domain model + capabilities
                         │
             product profile / hostname
                         │
       ┌─────────────────┼─────────────────┐
       ▼                 ▼                 ▼
    DataTug          Incidentius       Dashboardius
       │                 │                 │
       ▼                 ▼                 ▼
 Investigation        Incidents         Dashboards
 developer home         home              home
```

For MVP, **do not build three applications**. Build Incidents inside
DataTug and reuse Dashboardius capabilities throughout incident
investigation, metrics, recovery monitoring, and replay.

The engineering entry point is DataTug. The go-to-market/storytelling
entry point may eventually be Incidentius because incident resolution is
easier to explain than a generic data/investigation platform.

------------------------------------------------------------------------

# 2. Product thesis

During an incident people need to know:

-   What do we actually know?
-   Which facts are confirmed and which are assumptions?
-   Who checked this already?
-   What query did they run and what did it show at that time?
-   Which business entities are affected?
-   Which comparable entities are healthy?
-   What differs between good and bad cases?
-   Which hypotheses were rejected, and why?
-   Are people duplicating work?
-   Which deployment changed this?
-   Where is the relevant code?
-   Who owns or understands this artefact?
-   Have we seen this before?
-   Which checks can run automatically?
-   How bad is the problem?
-   Is it recovering?
-   How do we prove the fix worked?
-   What should we watch for the next five hours?

Incidentius should become the **shared external working memory of the
troubleshooting team**.

``` text
DataTug
  "Interrogate the system and data."
        │
        ▼
Incidentius
  "Why are we investigating, what have we learned,
   and what should happen next?"
        │
        ▼
Dashboardius
  "Show the evidence, impact, progress, recovery,
   and useful investigation views."
```

# 3. Foundational principles

## 3.1 Incident resolution, not incident management

Integrate with established incident-management systems. Do not compete
on paging, on-call schedules, severity bureaucracy, status pages, ticket
workflow, or SLA machinery.

Ask of every feature:

> **Does this help somebody understand or resolve the problem faster?**

## 3.2 Automated bureaucracy

> **Troubleshooters produce investigation knowledge. Incidentius
> produces communication and bureaucracy.**

``` text
engineers investigate naturally
      ↓
queries / voice / messages / annotations / hypotheses
      ↓
structured shared investigation state
      ↓
AI continuously derives:
  - technical status
  - support update
  - operational overview
  - management update
  - timeline
  - handover
  - RCA/postmortem
  - reusable knowledge
```

Encourage troubleshooters to **think aloud and be verbose**. Humans may
be messy; AI creates structure.

## 3.3 Wrong theories are useful

A rejected hypothesis narrows uncertainty. The wrong answer can contain
correct knowledge.

> **Negative evidence is positive progress.**

## 3.4 Incidents are learning moments

When everything works, people rarely document hidden relationships.
During incidents they explain how the system really behaves.

> **Every incident should leave the organisation smarter than it found
> it.**

## 3.5 Every incident should make the next one more boring

A recurrence should not cost the same investigation effort as the first
occurrence.

> **Every incident should make the next one more boring.**

## 3.6 Integrate deeply, depend loosely

Sibling products should amplify the experience without becoming hard
requirements. Cross-promotion should happen at the moment of
demonstrated value.

# 4. MVP architecture: DataTug first

For MVP:

``` text
DataTug application
├── existing DataTug capabilities
├── Dashboardius capabilities / dashboards
└── Incidents        ← new
```

Do not create a separate Incidentius frontend/repository.

Longer term:

``` text
profile: datatug
  home = developer/investigation workspace

profile: incidentius
  home = incidents

profile: dashboardius
  home = dashboards
```

All profiles retain the same underlying capabilities. Engineering should
choose the simplest hostname/configuration mechanism.

> **Build capabilities once. Expose them through different product
> shells.**

# 5. DataTug owns Investigation Context

This is a firm boundary. Incidentius attaches an incident to DataTug
Investigation Context; it does not reimplement it.

``` text
Incidentius Incident
│
├── identity / lifecycle
├── participants / coordination
├── hypotheses / questions
├── evidence relationships
├── recurrence relationships
├── resolution / outcome
├── incident metrics
│
└── attached DataTug Investigation Context
```

A normal DataTug investigation may later become an incident without
losing its context.

# 6. DataTug is an operational map

DataTug is not merely a query library.

``` text
                   DataTug
                      │
       ┌──────────────┼───────────────┐
       ▼              ▼               ▼
  Environments   Infrastructure    Data sources
       │              │               │
 DEV / UAT / PROD   servers          databases
                    jobs             APIs
                    queues           logs
                    services         files
       │              │               │
       └──────────────┼───────────────┘
                      ▼
               Queries / Checks
                      │
                      ▼
            Investigation Context
```

Resources may belong to one environment or be shared between
environments.

# 7. DALgo and heterogeneous sources

DALgo is the data-access abstraction. Logs are first-class DataTug data
sources, not a separate silo.

``` text
                   DALgo
                     │
     ┌───────────────┼────────────────┐
     ▼               ▼                ▼
 Databases          Logs             APIs
     │               │                │
     └───────────────┼────────────────┘
                     ▼
               DataTug queries
                     │
                     ▼
             record-set semantics
```

Cross-source correlation should be natural:

``` text
Orders WHERE status = WAITING
        ↓
affected_orders
        ├── JOIN order_details
        ├── JOIN invoices
        ├── JOIN customers
        └── JOIN logs
              ON Order.id / trace_id / correlation_id
```

The system should be able to ask:

> **What log patterns distinguish affected orders from healthy orders?**

Log patterns should themselves become knowledge objects linked to
features, code, teams, queries, checks, and incidents.

# 8. Logical data model: semantic backbone

DataTug models business entities independently of physical storage.

``` text
                 LOGICAL MODEL

                    Order
                   /     \
             Customer    Invoice
                │           │
                └── Shipment
```

Physical mappings can differ completely:

``` text
Order.id
 ├── PROD PostgreSQL: ord_hdr.order_no
 ├── UAT PostgreSQL:  orders.id
 ├── Legacy DB:       tx_order.oid
 ├── Logs:            fields.orderId
 ├── Shipping API:    merchantReference
 └── Events:          payload.aggregateId
```

Investigation Context parameters should strongly prefer logical
parameters:

``` text
Order.id = 72841
Invoice.id = INV-9912
Customer.id = C-17
```

Physical parameters remain an escape hatch when the logical model is
incomplete.

# 9. Improve the logical model while investigating

Example:

``` text
Discovered during INC-2841

shipping_requests.merchant_ref
          appears to map to
              Order.id

Evidence:
  47/47 investigated records matched

[Confirm mapping]
```

Once validated, every future developer, tester and incident benefits.

> **Use AI to reason over explicit organisational knowledge, not as a
> substitute for modelling organisational knowledge.**

# 10. Data flow and lineage

Capture how logical entities flow through systems:

``` text
Order
  created in Web App
        ↓
  orders-db
        ↓
  OrderCreated event
        ↓
  allocation-worker
        ↓
  shipment request
        ↓
  FedEx API
        ↓
  carrier acknowledgement
        ↓
  shipment_events
        ↓
  Order.status update
```

Then ask:

> **Where did affected records diverge from healthy records?**

``` text
Healthy O21:
created ✓ → allocated ✓ → shipment ✓ → ACK ✓ → READY ✓

Affected O11:
created ✓ → allocated ✓ → shipment ✓ → ACK ? → READY ✗
                                      ↑
                              investigate here
```

# 11. Context semantics: cohorts, roles and overlays

A context parameter has both a logical type and an investigation role.

``` text
type: Order.id

role:
  affected
  healthy_control
  suspected
  excluded
  recovered
```

Troubleshooting often means:

> **What differs between the bad set and the good set?**

``` text
Compare:
  affected_orders
vs
  healthy_control_orders

Across:
  customer_type
  server
  deployment
  status transitions
  processing time
  worker
  feature flags
  payment method
  external agent
```

## 11.1 Layered/scoped context

Context is not one flat bag.

``` text
INCIDENT CONTEXT
│
├── Canonical / global
│   affected_orders = [O1 ... O10]
│   start_time      = 14:32
│   environment     = production
│
├── Hypothesis H17
│   "These five orders are affected too"
│   affected_orders += [O11 ... O15]
│
├── Payments team
│   suspicious_payments = [...]
│
├── John's investigation
│   candidate_workers = [...]
│
└── Question Q12
    "Are enterprise orders affected?"
    candidate_customers = [...]
```

Everything can be visible without becoming canonical truth.

If H17 is confirmed:

``` text
H17 confirmed
O11...O15
    │
    └── PROMOTE
           ↓
Canonical affected_orders
O1...O15
```

If rejected, the candidate values remain historically attached to H17
but do not contaminate canonical context.

> **Useful speculation should be immediately usable without becoming
> truth.**

# 12. Source-agnostic contextual record sets and inGitDB

Investigation Context may reference record sets from any
DALgo-accessible source:

-   database query;
-   log query;
-   API result;
-   inGitDB;
-   literals;
-   other sources.

inGitDB is a convenient, powerful place to persist/version temporary or
reusable record sets, but it is not a special requirement of
Investigation Context.

Example progression:

``` text
ad-hoc values
   ↓
temporary record set
   ↓
used successfully in investigation
   ↓
cleaned / documented
   ↓
persistent reusable inGitDB dataset
   ↓
development + testing + future incidents
```

# 13. Cross-environment comparison

A central troubleshooting question:

> **Why does this work in UAT but fail in PROD?**

``` text
            Same DataTug query/check
                    │
          ┌─────────┴─────────┐
          ▼                   ▼
         UAT                 PROD
          │                   │
          └─────────┬─────────┘
                    ▼
                   DIFF
```

Comparison should eventually include data, configuration, deployments,
reference data, external-provider configuration, checks, and relevant
infrastructure.

The broader primitive is:

> **Compare something that works with something that does not.**

Examples:

-   bad vs good records;
-   PROD vs UAT;
-   before vs after deployment;
-   FedEx vs UPS;
-   affected vs unaffected customers;
-   server-17 vs other servers;
-   current incident vs previous occurrence.

# 14. External agents / counterparties

Avoid hard-coding only Customer and Vendor. Model a generic external
agent/counterparty:

``` text
ExternalAgent
├── Customer
├── Shipping provider
├── Payment processor
├── Bank
├── SaaS/API provider
├── Partner
└── Other external party
```

Agent **type** gives generic knowledge. Agent **instance** gives
operational knowledge.

``` text
ShippingProvider
│
├── FedEx
│   ├── protocol / endpoints
│   ├── tables / log locations
│   ├── health checks
│   ├── DataTug queries
│   ├── failure patterns
│   └── previous incidents
│
└── UPS
    ├── different protocol
    ├── different tables
    ├── different checks
    └── different failure semantics
```

Agent identity is another cohort dimension.

Example:

``` text
Affected orders
  FedEx   93%
  UPS      4%
  DHL      3%

Healthy control
  FedEx   31%
  UPS     42%
  DHL     27%
```

Incidentius should not conclude that FedEx is broken, but should surface
the overrepresentation and immediately offer FedEx-specific diagnostics.

# 15. Query-result snapshots

Query-result snapshots should be first-class immutable evidence.

A query execution associated with an incident should preserve enough
information to understand what was actually observed:

``` text
QueryExecution
  query/version
  parameters
  incident
  user / agent
  datasource + environment
  executed_at
  duration
  permissions/grants used
  sanitisation policy
  row count
  result snapshot
  result fingerprint
```

Benefits:

1.  **Historical evidence:** production data changes; the snapshot
    preserves what the engineer saw.
2.  **Shared cache:** expensive/long-running queries do not need to be
    rerun by every investigator.
3.  **Collaboration:** annotations attach to the same result set.
4.  **Comparison over time:** snapshots show recovery.

Example:

``` text
Snapshot 14:43
        ↓
Snapshot 14:51
        ↓
Snapshot 15:02

37 affected orders → 19 → 3
```

Historical snapshots must not become a backdoor around current
permissions. Retention, masking and access semantics need explicit
design.

# 16. Row/value/entity-level annotations

A major differentiator is annotation at the level where evidence
actually exists.

Support annotations on:

-   entity;
-   row;
-   row set;
-   field;
-   value;
-   relationship;
-   query/result;
-   snapshot.

Examples:

> `Order/72841` --- "Retry happened correctly; probably not the retry
> bug."

> `status=PROCESSING` on these 37 orders --- "All were created on
> server-17 between 14:31--14:36."

> `customer_type=enterprise` --- "Is this actually limited to enterprise
> customers?"

Annotations should automatically carry provenance: incident, author,
timestamp, query, parameters, source, snapshot and selected
entities/values.

Prefer **data-addressable** annotations where possible. If an order is
annotated through one query and appears later through a different query,
the relevant annotation should be discoverable again.

# 17. Hypotheses, evidence and uncertainty reduction

Hypotheses should be first-class.

Possible states:

``` text
Hypothesis
├── proposed
├── investigating
├── supported
├── weakened
├── rejected
└── confirmed
```

Example:

``` text
H12: Retry mechanism malfunction
status: rejected

because:
  E31: retries completed for affected orders
  E34: retry latency within normal range

scope:
  affected Order.id set

investigation:
  queries: [Q19, Q22]

confidence: high
```

Do not throw away failed paths.

A new engineer joining later should see **what was checked, against
which records, using which queries, what was observed, and why the
theory was rejected**.

# 18. From human troubleshooting to automated checks

Structured incident history makes automated hypothesis validation
possible.

The progression should be:

``` text
human investigation
      ↓
AI-assisted suggestion
      ↓
repeatable investigation pattern
      ↓
AI proposes executable check
      ↓
historical testing + human review
      ↓
deterministic check
      ↓
automatic future investigation
```

A generic DataTug `Check` abstraction may support purposes such as:

``` text
Check
 ├─ health
 ├─ diagnostic
 ├─ hypothesis-validation
 ├─ remediation-verification
 └─ regression
```

Each check may have:

-   inputs;
-   queries;
-   expected conditions;
-   evaluation;
-   scope;
-   cost;
-   risk;
-   required privileges;
-   historical executions.

AI can discover/propose a check, but mature checks should become
deterministic where practical.

# 19. DataTug system health checks

Existing/defined health checks are immediate executable knowledge.

An incident should not start from a blank page.

``` text
INC-2841 created
"Enterprise orders stuck processing"
          ↓
Relevant system/feature identified
          ↓
DataTug health checks

✓ Database connectivity
✓ Payment settlement
✗ Finalizer backlog abnormal
✓ Retry processing
? Worker distribution
          ↓
Investigation starts with
several branches already narrowed
```

A useful investigation query can evolve into a diagnostic or permanent
health check.

# 20. Incident metrics: measure the problem and the recovery

Incident metrics are not generic observability dashboards. They are
metrics with meaning in the context of the incident.

Example:

``` text
IMPACT
  orders.status=WAITING       847
  affected customers           93
  oldest waiting order        47m

DIAGNOSTIC
  shipping queue depth          0
  finalizer queue depth     1,283
  finalizer throughput       12/min
  CPU                         43%

RECOVERY
  stuck orders                847 → 312 → 41 → 0
  processing latency          31m → 8m → 2m

WATCH
  new stuck orders              0
  finalizer queue depth        23
  error rate                 0.2%
```

Resolution criteria should be explicit but captured with minimal
bureaucracy:

``` text
✓ No new stuck orders for 30 minutes
✓ Existing backlog < 10
✓ Oldest WAITING order < 5 minutes
✓ Finalizer throughput within expected range
✓ Error rate back within baseline
```

If an engineer says:

> "Once the backlog reaches zero and stays there for half an hour, we're
> good."

AI should be able to turn that into a proposed resolution criterion.

# 21. Dashboardius during incidents

Dashboardius is a first-class part of the story.

During investigation, useful DataTug queries and metrics should be easy
to assemble into an incident-specific dashboard.

Over time Incidentius/DataTug may recognise recurring collections of
useful queries and propose reusable dashboards for particular incident
classes.

Example:

``` text
Order Processing Incident Dashboard

- affected orders
- processing latency
- failed jobs
- retry queue
- errors by server
- recent deployments
- customer impact
- similar incidents
```

Dashboardius should support:

-   live impact dashboards;
-   investigation dashboards;
-   affected-vs-control comparisons;
-   environment comparisons;
-   recovery dashboards;
-   post-resolution watch dashboards;
-   historical/replay visualisations.

A normal dashboard is designed beforehand. An **investigation-generated
dashboard** emerges from the actual work that proved useful.

# 22. Recovery and watch period

Resolution is not necessarily the end of useful monitoring.

Conceptually:

``` text
Investigating
     ↓
Mitigating
     ↓
Recovering
     ↓
Resolved
     ↓
Watching
     ↓
Closed
```

These do not all need to become bureaucratic mandatory statuses, but the
semantics are useful.

Example timeline:

``` text
14:32   47 stuck orders
14:41  183
14:52  847  ← peak impact

14:57  mitigation deployed

15:03  612
15:12  311
15:24   41
15:31    0

16:01  resolution criteria satisfied

16:01–21:01
       observation period
       no recurrence
```

ASCII charts/diagrams like the following should be used throughout
product documentation and demos:

``` text
Queue
depth
  ▲
  │       ╭────╮
  │     ╭─╯    ╰╮
  │   ╭─╯       ╰──╮
  │───╯              ╰──────
  └──────────────────────────► time
        ↑          ↑
     incident   mitigation
      start      applied
```

# 23. Incident lifecycle and outcomes

Start with:

> **Houston, we've got a problem**

Keep terminal outcomes deliberately small:

-   **Resolved** --- real problem, sufficiently fixed/mitigated.
-   **False alarm** --- no actual problem.
-   **Accepted / known issue** --- understood but deliberately not
    resolved now.
-   **Handed off** --- resolution responsibility moved elsewhere.
-   **Unresolved / investigation stopped** --- investigation ended
    without resolution.

Do not use duplicate as a terminal outcome. True duplicates should
merge.

A recurrence remains a separate incident and is linked to the previous
occurrence.

# 24. Recurrence as a first-class relationship

``` text
INC-104  Payment finalizer stopped
   │
   ├── recurrence → INC-187
   │                    │
   │                    └── recurrence → INC-241
   │
   └── related → INC-263
```

Track **Time to Recurrence Identification (TTRI)**:

> Time from incident start/detection until the team correctly identifies
> it as a recurrence.

Track both:

``` text
incident detected
      ↓
AI suggests possible recurrence    ← AI TTRI
      ↓
human confirms recurrence          ← Confirmed TTRI
```

Also measure whether recurrence recognition reduces resolution time.

# 25. Search/discovery

Search is crucial for recurrence identification.

"Have we seen anything like this before?" should use more than semantic
text similarity.

Possible signals:

-   symptoms;
-   logical entities;
-   context parameters;
-   affected/control cohorts;
-   log patterns;
-   external agents;
-   failed checks;
-   environments;
-   infrastructure;
-   code/features;
-   deployments;
-   metric shapes;
-   hypotheses;
-   resolutions;
-   DataTug queries used.

Long-term knowledge must understand **decay**. Queries, schemas,
mappings and architecture change. Knowledge should be versioned, scoped,
deprecated or marked stale rather than treated as timeless truth.

# 26. Voice as a first-class investigation input

Voice recognition should not be a convenience feature. It should be an
**investigation capture channel**.

People say things during incidents they would never bother typing:

> "This feels similar to that February outage, except those orders were
> all on EU-2... I'm not sure it matters, but these ones also seem to be
> older enterprise accounts."

Possible extracted signals:

``` text
possible recurrence → February incident
difference          → not EU-2
candidate factor    → account age
candidate cohort    → enterprise customers
confidence          → speculative
```

Preserve the raw transcript and provenance. AI interpretation must not
silently turn speculation into fact.

Voice, typed messages, Slack/Teams messages and pasted text should all
feed the same incident model.

# 27. Slack/Teams/chat connectors

A Slack bot attached to an incident thread can:

-   capture messages and reactions;
-   extract typed parameters and hypotheses;
-   link DataTug queries/results;
-   answer "has anyone checked X?";
-   post concise coordinator updates;
-   surface previous incidents;
-   detect duplicated investigation;
-   accept commands such as promoting candidate entities into canonical
    context.

Connectors are alternate capture surfaces, not separate incident
systems.

# 28. Parsing emails/messages/log snippets into typed context

Unstructured text should be parsed into typed logical parameters.

Example:

> Customer ACME says order 72841 and invoice INV-9912 failed around
> 14:32. Trace id `4f8a...`.

Possible extraction:

``` text
Order.id       = 72841
Invoice.id     = INV-9912
started_at     = ~14:32
trace_id       = 4f8a...
Customer       = ACME
```

Use a hybrid approach:

``` text
email / message / voice / logs
            ↓
 deterministic extractors
       + AI extraction
            ↓
 typed logical parameters
```

Regex/pattern extraction remains valuable for known formats. AI handles
flexible/implicit language.

Every extracted value needs provenance, extraction method and
confidence.

# 29. Self-learning extraction skills

AI can initially infer a value and later propose deterministic
extraction.

``` text
1. AI infers "Order 72841" → Order.id
2. Pattern repeats.
3. Incidentius proposes deterministic extractor.
4. Test against historical examples.
5. Human approves.
6. Future extraction becomes deterministic.
7. AI remains fallback for ambiguity.
```

Treat extractors as versioned skills with:

-   output logical type;
-   patterns/parsers;
-   positive/negative examples;
-   source systems;
-   confidence rules;
-   ownership;
-   false-positive statistics;
-   provenance.

General progression:

``` text
Extraction skill
   ↓
Query-selection skill
   ↓
Diagnostic skill
   ↓
Hypothesis-validation skill
   ↓
Resolution-verification skill
```

Self-learning should mean: **learn autonomously, propose autonomously,
test autonomously, promote cautiously.**

# 30. Shared investigation and AI coordination

Incidentius should distinguish an **AI Troubleshooter** from an **AI
Coordinator**.

The troubleshooter investigates technical hypotheses.

The coordinator maintains the team's shared mental model:

``` text
Who is investigating what?
What do we currently know?
What has been ruled out?
What remains unanswered?
Are two people duplicating work?
What can be investigated in parallel?
Who has relevant expertise?
What evidence contradicts a new hypothesis?
What useful work has no owner?
```

Example:

> **John:** Maybe orders got stuck waiting for a truck because the
> shipping API is down?

> **AI Coordinator:** Jack checked the shipping queue three minutes ago.
> It was empty and processing about 10 messages/minute, which is within
> its expected range. A shipping API outage therefore looks unlikely.
> Jack is checking carrier acknowledgements. He suggested checking
> whether affected orders received a `ready_for_dispatch` event in
> parallel.

The coordinator should intervene selectively. Constant AI chatter would
itself become a distraction.

# 31. Automatically generated status views

The same structured incident state should generate different
projections.

For troubleshooters:

``` text
47 orders affected.
Shipping queue healthy.
Shipping API outage largely ruled out.
Jack checking carrier acknowledgements.
Maria comparing affected orders against deployment a81f2c.
Nobody currently checking ready_for_dispatch event generation.
```

For leadership:

``` text
Customer order processing is degraded.
47 affected orders identified.
No new affected orders in the last 8 minutes.
Shipping infrastructure appears healthy.
Investigation is focused on event generation and a recent deployment.
No confirmed root cause yet.
```

For support:

``` text
Some orders are experiencing delayed processing.
The issue is under active investigation.
Existing orders remain recorded.
No evidence of data loss has been found.
```

AI must carefully distinguish **facts, hypotheses and inference**.
Generated statements should retain provenance.

# 32. Deployment, code and SpecScore correlation

SpecScore can define DataTug queries/health checks for features, and
features can link to source code.

Incidentius should exploit deployment information:

``` text
Incident evidence:
  Order.status = PROCESSING
        ↓
DataTug:
  logical field + checks
        ↓
SpecScore:
  feature + expected behaviour
        ↓
deployment:
  service / commit range
        ↓
code:
  changed source related to Order.status
```

Example:

``` text
Deployment: order-service v2.41.7
Deployed: 14:17
Incident started: 14:32

Changed code related to Order.status:
  order_finalizer.go
  changed transition PROCESSING → READY
  commit d1f904a
```

Relevance should combine temporal and semantic signals, not merely
"deployment happened recently".

# 33. Codegrapher integration

Codegrapher should be an optional amplifier, not a hard dependency.

Given evidence such as:

-   database table;
-   column;
-   stored procedure;
-   log pattern;
-   stack trace;
-   feature;

Codegrapher can surface:

-   repositories;
-   source files;
-   symbols;
-   callers;
-   recent changes;
-   lint/analysis findings;
-   ownership.

If Codegrapher is not configured, the investigation must continue
normally.

# 34. Ownership and expertise discovery

Ownership should be modelled around technical artefacts:

``` text
Person / Team
  ↕
Repository
Source path / symbol
Database
Table / column
Stored procedure
Job
Queue/topic
API endpoint
Feature
Spec
Health check
Stack-trace pattern
```

Distinguish:

-   **Owner** --- accountable.
-   **Maintainer** --- routinely works on it.
-   **Expert** --- inferred from history/Git/incidents.
-   **Current investigator** --- working on it now.

Git history is evidence, not authority.

Example responder suggestions:

``` text
Payments Team — owns payments table
Alice — recent changes to settle_payment
Carol — resolved 3 incidents with this stack-trace pattern
Platform Team — owns PaymentRetryWorker
```

# 35. Incident-scoped privileged production access

One of the strongest DataTug capabilities is controlled, audited,
incident-scoped production access.

Normally:

``` text
Developer → Production
           DENIED
```

During an assigned incident:

``` text
INC-2841
  affected: Payments / Customer X
  assigned: Alice, Bob
        ↓
Temporary investigation grant
        ↓
DataTug
  ✓ payments: SELECT
  ✓ customers: SELECT
  ✓ orders: SELECT
  ✗ sensitive card data
  masked: email, phone, address
  expires: incident resolved / time limit
```

Engineers should not need raw production credentials. They receive
permission to perform specific operations through DataTug.

Sensitive-field escalation should be narrow, justified and temporary.

# 36. Controlled production modifications

Writes deserve stricter treatment than reads.

Example:

``` text
Production modification

Incident: INC-2841
Environment: PROD
Operation: UPDATE
Affected rows: 1

Before
status = PROCESSING

After
status = READY

Requested by: Alice
Reason: unblock affected order

[Request approval]
```

Record:

``` text
who
when
incident
query/action
parameters
approval
affected entities
before state
after state
result
```

This should answer not only:

> Who changed Order 72841?

but also:

> Why was it changed?

because the action remains connected to incident evidence and
resolution.

# 37. Masking and AI-visible data

Masking must happen in the trusted execution path before prohibited data
reaches the client or AI.

Potential policy dimensions:

``` text
human-visible
AI-visible
queryable
aggregatable
exportable
modifiable
```

An AI agent might be allowed to reason that 17 affected customers share
a domain without receiving full email addresses.

Security architecture deserves its own design track and may extend
beyond MVP.

# 38. RCA is not merely a document

The durable artefact should be a structured **Incident Resolution
Record**.

RCA is one generated view.

``` text
Incident Resolution Record
        │
        ├── Live incident workspace
        ├── Dashboardius views
        ├── Interactive timeline
        ├── RCA
        ├── Postmortem
        ├── Executive summary
        ├── Support summary
        ├── Training replay
        ├── Troubleshooting guide
        ├── Automated checks
        └── Video/story
```

The record should be a mini-site with hyperlinks and attached artefacts.

Users should be able to inspect historical DataTug snapshots and, where
safe/appropriate, rerun queries against current environments.

# 39. Timeline replay and executable learning

A future engineer should be able to replay an incident:

``` text
14:31 — First symptom
14:36 — Incident started
14:39 — Shipping hypothesis proposed
14:41 — Shipping health check run
14:43 — Shipping hypothesis weakened
14:52 — Deployment correlation found
14:57 — Mitigation applied
15:31 — Backlog reaches zero
16:01 — Resolution criteria satisfied
```

Technical replay should expose:

-   queries;
-   parameters;
-   snapshots;
-   annotations;
-   hypotheses;
-   code diffs;
-   metrics;
-   decisions.

Training mode may pause:

> At this point the team knew A, B and C. What would you investigate
> next?

Historical incidents can become interactive debugging exercises without
somebody authoring a separate course.

# 40. Generated videos and storytelling

Once the incident graph is structured, AI can generate serious or
humorous summaries.

Example title:

> **INC-2841 --- The Curious Case of the 37 Waiting Orders**

The humour should come from recognisable troubleshooting situations, not
from making customer impact itself a joke.

A generated incident video can become useful for:

-   engineering all-hands;
-   onboarding;
-   training;
-   marketing;
-   cultural storytelling.

This is another form of automated bureaucracy: nobody should spend hours
manually creating the presentation.

# 41. Appreciation: the organisation should know its heroes

Resolved incidents should visibly credit contributors.

Recognition should not go only to the person who found the root cause.

Useful contributions include:

-   identifying affected cohorts;
-   eliminating a plausible hypothesis;
-   providing the critical customer/order ID;
-   remembering a previous incident;
-   creating a reusable check;
-   implementing mitigation;
-   verifying recovery;
-   coordinating effectively.

Example:

``` text
Incident resolved in 47 minutes

John
  proposed key hypotheses and narrowed the affected cohort

Jack
  eliminated the largest false lead in six minutes
  and created a reusable diagnostic check

Maria
  correlated the failures with deployment a81f2c

Anna
  implemented and verified the mitigation

Sam from Support
  provided the customer evidence that unlocked the investigation
```

Avoid simplistic competitive rankings that create bad incentives.

Use specific, evidence-backed appreciation.

Recognition is not merely social polish. It encourages the behaviours
the knowledge system needs: thinking aloud, sharing evidence,
documenting reasoning and creating reusable knowledge.

# 42. Personas

Separate **persona**, **incident role**, and **permission**.

## Primary personas

-   Developer / Engineer
-   Tester / QA Engineer
-   Support Specialist
-   SRE / Operations Engineer
-   Data / DB Engineer
-   Incident Coordinator / Manager
-   Product Owner / Domain Expert
-   Engineering Manager / Team Lead
-   Security / Compliance
-   Executive / CEO
-   Business User
-   Observer / Learner

## External participants

Use the generic External Agent/Counterparty model for customers,
providers, carriers, banks, partners, etc.

## Business User terminology

Use **Business User** for someone who uses the application/system being
troubleshot.

Reserve plain **User** for a DataTug / Incidentius / Dashboardius user.

A Business User may provide uniquely valuable evidence:

> "I approved these three invoices. The screen said Approved. Then I
> refreshed and they were back in Pending."

That can produce affected/control cohorts and reproduction context.

## Tester / QA Engineer

QA is a major DataTug persona, not an afterthought.

Testers need to:

-   inspect data state;
-   compare expected vs actual;
-   verify jobs/events;
-   create datasets/cohorts;
-   reuse incident failure modes;
-   create regression checks;
-   validate fixes.

# 43. Persona vs incident role vs permission

Example:

``` text
Persona
  Developer

Incident role in INC-2841
  Investigator

Permissions
  production SELECT
  PII masked
  no production writes

Temporary grant
  selected sensitive field access for 30m
```

A developer may be an investigator in one incident and coordinator in
another.

Do not encode CEO, Product Owner, etc. directly as RBAC roles.

# 44. Knowledge flywheel: incidents ↔ development ↔ testing

Knowledge flows in both directions.

``` text
              DEVELOPMENT
                  │
        features / bug fixes
        queries / reasoning
        implementation knowledge
                  │
                  ▼
               DATATUG
          ORGANISATIONAL KNOWLEDGE
             ▲           ▲
             │           │
       TESTING            │
   scenarios / checks     │
   edge cases / datasets  │
             │            │
             └─────┬──────┘
                   │
                   ▼
          INCIDENT RESOLUTION
       evidence / hypotheses
       failures / resolutions
                   │
                   └──────────► back into knowledge
```

A query written at 02:00 during an incident may later become:

-   a developer debugging query;
-   a QA validation query;
-   a regression check;
-   a health check;
-   a troubleshooting scenario;
-   a data-quality assertion;
-   a Dashboardius card;
-   an automated hypothesis validation.

The reverse is equally important.

A developer's query or tester's observation created months before an
incident may prevent a wrong investigation branch in production.

> **Incidentius should not start learning when an incident starts.**

# 45. SpecScore integration

SpecScore can link:

-   features;
-   requirements;
-   expected behaviour;
-   DataTug queries/checks;
-   implementation/source.

Incident findings can reveal:

-   missing specification;
-   ambiguous expected behaviour;
-   required change;
-   regression scenario.

Optional workflow:

``` text
Incident reveals undocumented behaviour
        ↓
Incidentius records finding
        ↓
SpecScore available?
        ↓
Create/update requirement
linked back to incident
```

SpecScore should amplify the workflow without becoming a hard
requirement.

# 46. OpenVaultDB and inGitDB integration

Useful contextual/reference datasets may be stored and versioned in
inGitDB.

OpenVaultDB can provide richer project/reference-data management using
inGitDB as its backbone.

Again: optional amplification, not mandatory coupling.

The UI should surface the integration at the moment value becomes
obvious.

# 47. Pricing and growth strategy

Initial commercial principle:

> **Individual investigation should be free. Collaborative incident
> resolution is the commercial product.**

Even stronger:

> **Knowledge spreads for free. Collaboration creates the bill.**

Free/personal use should remain genuinely useful.

Viewing incident knowledge should be free or extremely generous:

-   Resolution Record;
-   queries;
-   snapshots;
-   dashboards;
-   timeline;
-   annotations;
-   hypotheses;
-   replay;
-   learning materials.

This is a growth mechanism. Early adopters can invite colleagues to real
incidents and let them experience value without procurement friction.

Consider monthly active contributor pricing rather than assigned seats.

Possible philosophy:

``` text
12 people viewed incidents
 7 actively contributed
 5 viewers remain free

bill based on active contributors
```

Contribution should be allowed before payment. Do not put a paywall in
front of somebody joining an urgent incident.

If an organisation does not pay, use organisation-level
warnings/grace/read-only restrictions later. Never make billing the
reason an outage lasts longer.

# 48. Show value alongside the invoice

Where measurements are factual, show the organisation what happened:

``` text
September

7 active troubleshooters
18 free viewers

6 incidents resolved
2 recurrences identified

Median recurrence identification
  21 min → 3 min

11 previous queries reused
23 hypotheses eliminated
4 reusable checks created
```

Be conservative with estimated savings. Prefer directly measured
metrics.

# 49. Product effectiveness metrics

Potential metrics include:

-   MTTR;
-   time to mitigation;
-   time to useful first hypothesis;
-   time to first relevant previous incident;
-   TTRI;
-   recurrence rate;
-   resolution-time improvement across recurrences;
-   percentage of incidents reusing historical queries/checks;
-   number of rejected hypotheses that prevented duplicated work;
-   reusable checks created;
-   investigation knowledge promoted into deterministic skills;
-   time spent requesting status updates;
-   time spent preparing postmortems;
-   recovery verification time.

Do not optimise vanity metrics such as number of comments or number of
queries.

# 50. DataTug CLI: incidents as first-class agent surface

Incidents must be first-class in the DataTug CLI/API, not only the web
UI.

Use the Sneat ecosystem convention of singular command namespaces:

``` bash
datatug incident list
datatug incident search
datatug incident show INC-2841
datatug incident similar INC-2841
datatug incident create
datatug incident resolve INC-2841
```

Exact commands should be designed against existing CLI conventions, but
capabilities should include:

-   list/search;
-   inspect;
-   create/open;
-   update;
-   attach/query Investigation Context;
-   inspect/add hypotheses;
-   attach evidence;
-   search similar incidents;
-   link recurrence;
-   resolve;
-   consume structured event streams.

Important architectural rule:

> **Every important incident-resolution capability should be API-first
> and, where appropriate, exposed through the DataTug CLI so humans,
> scripts and AI agents operate on the same model.**

# 51. `datatug incident watch`

Support both:

``` bash
datatug incident watch
```

Watch for new/reopened/escalated incident activity accessible in the
current project/workspace.

And:

``` bash
datatug incident watch INC-2841
```

Follow a specific incident.

Human-readable example:

``` text
10:42:17  CONTEXT     +5 suspected Order.id
10:43:02  HYPOTHESIS  H17 created: "FedEx acknowledgement failure"
10:43:31  QUERY       Jack ran fedex-ack-health
10:43:36  EVIDENCE    E31 attached to H17
10:43:51  HYPOTHESIS  H17 → REJECTED
10:44:03  METRIC      affected_orders 47 → 39
10:44:19  ASSIGNMENT  Maria → investigate deployment d81fa2
```

Structured output is essential for AI agents.

Potential future semantics:

``` bash
datatug incident watch INC-2841 --json
datatug incident watch INC-2841 --since <cursor>
```

This implies an underlying structured incident event/change stream.

The same stream can power:

``` text
live watch
    │
    ├── AI coordination
    ├── Slack connector
    ├── notifications
    ├── status generation
    └── timeline replay
```

# 52. Autonomous/agent-oriented incident workflow

A persistent agent can watch for new incidents:

``` text
datatug incident watch
        │
        ▼
   NEW INC-2841
        │
        ├── read Investigation Context
        ├── normalise/extract parameters
        ├── search similar incidents
        ├── identify possible recurrence
        ├── find relevant checks
        ├── run permitted safe checks
        ├── compare environments
        ├── inspect recent deployments
        ├── identify owners/experts
        └── propose initial hypotheses
```

By the time a human opens the incident, useful work may already be
complete.

# 53. Canonical demo incident: required deliverable

Create a canonical demo incident scenario, but **do not prematurely lock
its exact story in this document**.

It must serve simultaneously as:

``` text
E2E fixture
   ↓
interactive demo
   ↓
onboarding/training
   ↓
marketing story
   ↓
regression suite
```

Requirements:

-   understandable without specialist domain knowledge;
-   technically rich;
-   complex but rewardingly challenging;
-   multiple plausible hypotheses;
-   convincing red herrings;
-   collaborative investigation;
-   useful wrong theories;
-   DataTug queries and snapshots;
-   Dashboardius metrics/dashboards;
-   annotations;
-   context extraction;
-   affected/control cohorts;
-   environment comparison;
-   external agent involvement;
-   deployment/code correlation;
-   ownership/expertise;
-   optional privileged access;
-   mitigation and verification;
-   appreciation/credits;
-   interactive replay.

The tone should feel like a **detective dramedy**. Humour should come
from recognisable troubleshooting situations rather than trivialising
customer impact.

# 54. The crucial post-credits recurrence

The canonical demo needs a second occurrence.

The first incident is the movie.

Then the credits roll.

``` text
                 RESOLVED

                  CREDITS

      John ............ Shipping hypothesis
      Maria ........... Deployment correlation
      Jack ............ Diagnostic checks
      Anna ............ The fix
      Incidentius ..... Excessive paperwork

                     ...

             THREE MONTHS LATER
```

Same symptom appears.

No dramatic music.

``` text
Possible recurrence of INC-2841 detected.
4/5 characteristic signals match.
Running previously validated diagnostics…

✓ Shipping healthy
✓ Database healthy
✗ Finalizer worker assignment
✓ Customer configuration normal

Likely recurrence: INC-2841
```

Example result:

``` text
Detected                 09:17
Recurrence suggested     09:17:21
Confirmed                09:18
Mitigation               09:21
Recovered                09:24

TTRI                     1m
Time to resolution       7m
```

> **The P.S. is boring because the product worked.**

This contrast should be a key marketing device.

# 55. Why Incidentius may be the easiest product story to sell

Although DataTug is the technical backbone, incident resolution may be
the easiest commercial narrative.

Incidents connect:

-   developers;
-   testers;
-   support;
-   SRE/operations;
-   DB/data engineers;
-   incident managers;
-   product/domain experts;
-   engineering managers;
-   executives;
-   business users;
-   external counterparties.

Everyone shares one measurable goal:

> **Resolve the issue ASAP.**

A canonical incident can demonstrate the whole stack without first
explaining why somebody needs a "modern data investigation platform".

``` text
"Houston, we've got a problem"
            ↓
customer/business symptom
            ↓
typed logical context
            ↓
DataTug queries + logs
            ↓
affected vs healthy comparison
            ↓
Dashboardius metrics
            ↓
hypotheses + annotations
            ↓
health checks
            ↓
external-agent correlation
            ↓
UAT vs PROD comparison
            ↓
deployment → feature → code
            ↓
ownership → right people
            ↓
controlled production action
            ↓
recovery dashboard
            ↓
resolution + replay
            ↓
reusable checks / knowledge
            ↓
boring recurrence
```

This demonstrates DataTug rather than merely describing it.

# 56. Technical backbone ≠ go-to-market entry point

Preserve this distinction:

**Engineering strategy**

> Build the single DataTug application. Add Incidents. Reuse
> Dashboardius. Do not create three implementations.

**Potential GTM strategy**

> Lead with Incidentius because the urgent problem and ROI are obvious.

Incidentius can become the wedge through which organisations discover
the broader DataTug / Dashboardius capabilities.

# 57. Architecture and design questions engineering should resolve

Do not treat the brainstorm as a final schema.

The engineering/design exercise should explicitly resolve:

1.  Exact incident domain model.
2.  Immutable event/provenance model vs mutable current state.
3.  How replay reconstructs "what was believed at 14:42".
4.  How Incidentius references DataTug Investigation Context.
5.  How context overlays/promotions/rejections are represented.
6.  How result snapshots are stored, retained and permissioned.
7.  How generic source-agnostic record-set references work.
8.  How Dashboardius incident dashboards are represented/reused.
9.  How incident search/similarity is indexed and scored.
10. How knowledge applicability/version/staleness is represented.
11. Event-stream mechanics behind `datatug incident watch`.
12. AI trust boundaries: fact vs human statement vs hypothesis vs AI
    inference vs deterministic result.
13. Production-access and approval architecture.
14. Product-profile/hostname mechanism.
15. Existing capabilities that can be reused rather than rebuilt.

# 58. Required audit before design

Before implementing new abstractions, inspect the existing DataTug
code/specs for:

-   Investigation Context;
-   logical model;
-   environments;
-   infrastructure;
-   DALgo;
-   data-source model;
-   logs;
-   query library;
-   health checks;
-   snapshots/caching;
-   annotations;
-   inGitDB;
-   Dashboardius integration;
-   access control;
-   auditing;
-   CLI conventions;
-   existing AI/agent surfaces.

Do not redesign capabilities that already exist.

Record gaps as explicit requirements/extensions.

# 59. MVP vs later vision

The document contains a broad vision. Engineering must ruthlessly
classify capabilities:

``` text
MVP
NEXT
FUTURE
```

MVP should prove the core thesis with minimal scope inside DataTug.

A reasonable MVP direction to evaluate:

-   Incident model/list/detail.
-   "Houston, we've got a problem" creation.
-   Attach/create DataTug Investigation Context.
-   Incident timeline/activity.
-   Hypotheses with support/reject outcomes.
-   Evidence links to DataTug queries/results.
-   Basic annotations/context linkage using existing capabilities.
-   Search/link related and recurrence incidents.
-   Basic incident metrics using Dashboardius.
-   CLI/API incident access.
-   `datatug incident watch`.
-   AI-generated status summary from structured state if existing AI
    infrastructure makes it cheap.
-   Canonical demo/E2E scenario covering implemented capabilities.

Do **not** force advanced production privilege escalation, autonomous
self-learning, generated videos, full Slack/Teams ingestion, or
sophisticated similarity ranking into the first MVP unless existing
infrastructure makes them unexpectedly cheap.

# 60. Final product vision

The long-term system should feel like this:

``` text
                  ORGANISATION WORKS
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
     Development       Testing       Incidents
          │              │              │
          └──────────────┼──────────────┘
                         ▼
                       DataTug
        logical model / queries / DALgo / context
        environments / infrastructure / knowledge
                         │
              ┌──────────┴──────────┐
              ▼                     ▼
          Incidentius           Dashboardius
       resolution reasoning     visual evidence
       coordination/replay      metrics/dashboards
              │                     │
              └──────────┬──────────┘
                         ▼
             ORGANISATIONAL KNOWLEDGE
                         │
       ┌─────────────────┼─────────────────┐
       ▼                 ▼                 ▼
 reusable queries    deterministic      better logical
 and datasets        checks/skills      system model
       │                 │                 │
       └─────────────────┼─────────────────┘
                         ▼
                 NEXT INVESTIGATION
                    starts smarter
```

The system should not merely remember incidents.

It should learn how the organisation investigates its own systems.

It should turn:

``` text
Observation
    ↓
Reasoning
    ↓
Hypothesis
    ↓
Investigation action
    ↓
Evidence
    ↓
Conclusion
```

into reusable, progressively more executable organisational knowledge.

The end state is not "AI writes a better postmortem".

The end state is:

> **A new incident starts with the accumulated investigative capability
> of every developer, tester, support specialist, domain expert and
> troubleshooter who has worked on the system before.**

And after the incident:

> **Nobody has to stop and write the story from scratch. The story,
> evidence, timeline, dashboards, RCA, replay, lessons, checks and
> credits already exist because Incidentius captured the work while it
> happened.**

------------------------------------------------------------------------

# 61. Core statements to preserve

These statements capture the intended philosophy and should survive
future rewriting:

> **Incidentius is an incident-resolution workspace.**

> **Troubleshooters produce investigation knowledge. Incidentius
> produces communication and bureaucracy.**

> **Negative evidence is positive progress because it reduces
> uncertainty.**

> **The wrong answer can contain correct knowledge.**

> **Every incident should leave the organisation smarter than it found
> it.**

> **Every incident should make the next one more boring.**

> **Incidentius should not start learning when an incident starts.**

> **Knowledge spreads for free. Collaboration creates the bill.**

> **Integrate deeply, depend loosely.**

> **Use AI to reason over explicit organisational knowledge, not as a
> substitute for modelling organisational knowledge.**

> **Build capabilities once. Expose them through DataTug, Incidentius
> and Dashboardius product profiles.**

------------------------------------------------------------------------

# 62. Engineering handoff instruction

Treat this document as a product vision and design input, **not as
permission to implement every idea immediately**.

The engineering team/frontier model should:

1.  Audit the existing DataTug / Dashboardius architecture and
    specifications first.
2.  Map every concept in this document to existing capability, required
    extension, or new capability.
3.  Preserve DataTug ownership of Investigation Context and DALgo-based
    data access.
4.  Preserve the single-app/single-repository direction for DataTug /
    Incidentius / Dashboardius.
5.  Produce a clear MVP boundary before implementation.
6.  Design the incident/event/provenance model with replay and AI-agent
    use in mind.
7.  Keep AI inference distinguishable from facts and deterministic
    evidence.
8.  Design API/CLI surfaces alongside UI surfaces.
9.  Create the canonical detective-story demo/E2E scenario after
    understanding the architecture.
10. Include ASCII diagrams and narrative scenarios in specifications
    where they explain intent better than isolated requirements.
11. Optimise primarily for quality of design and long-term coherence
    while keeping implementation time, complexity and token/cost usage
    in mind.
12. Avoid unnecessary questions: proceed autonomously on non-blocked
    areas and stop only for genuinely blocking decisions.

The objective is not to create another incident-management suite.

The objective is to make troubleshooting **faster, calmer, more
collaborative, more explainable, more secure, and progressively smarter
every time the organisation uses it**.
