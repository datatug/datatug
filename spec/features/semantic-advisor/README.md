---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Semantic Advisor

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-advisor?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-advisor?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-advisor?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/semantic-advisor?op=request-change) |
**Status:** Draft
**Date:** 2026-10-10
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Kind:** Capability (implementations in `datatug-core`, `datatug-cli` and `datatug-apps`)
**Pillars:** P4 meaning inference · P3 semantic layer · P5 related-data navigation · P9 Git-native projects · P12 provenance

## Summary

Specifies, and does not yet build, how DataTug suggests what a column means and what it refers to from fixed rules, records a person's answers, and writes the model, the meaning file and the source file.

This is the specification of what the ModelSpec design proposal calls "Phase 5, DataTug writes, with fixed rules only", without its benchmark, which is the sibling feature [`advisor-benchmark`](../advisor-benchmark/README.md). It stands on [`open-layer-reading`](../open-layer-reading/README.md), the proposal's Phase 4, which is itself a Draft that the owner has not approved. Nothing here is built, and none of it is part of the launch planned for 2026-10-25. The status is Draft. Only the owner approves it. Until he does, nothing here may be built. No AI model is called by anything in this feature, and nothing a person does in it leaves their machine.

How to read it: "What a person sees" is the whole feature in six moments. "Questions for the owner" holds the six answers it needs. "The author's choices" lists 33 things the author chose and the owner may overrule. "If the questions of open-layer-reading are answered otherwise" says what moves. The rest, and the appendix [formats.md](formats.md), is for whoever builds and tests it.

## What a person sees

Six moments, all on a copy of Chinook whose foreign keys were removed. Chinook is a music shop's database with 11 tables, 64 columns, 59 customers and 412 invoices. Nothing in this copy says that an invoice belongs to a customer. The moments are steps 3 to 7 and 10 of the journey in the proposal's section 14 (L2751 to L2755, L2758).

Maria knows this database. She has scanned it into the project `./shop` in a Git repository, as the database `chinook`. For the first five moments no registered model matches it.

**Step 3. "I open a table I know."** She opens `Invoice`. Its mark reads `Invoice 1 of 9`: only the key is described. Beside the table the Advisor shows two cards, and no more:

```
References · 1 line · inferred
  CustomerId refers to Customer
    ✓ named after Customer's key CustomerId            metadata, points
    ✓ Customer.CustomerId is the primary key           metadata, allows
    ✓ all 59 distinct values occur there               statistics, allows
    – no foreign key is declared                       metadata
  [Accept]  [Not a relationship]  [Follow]  [Not now]

Meanings · 2 lines · inferred
  Invoice is an Invoice
    ✓ the table's name is the label of the universal concept invoice
  Total is an Invoice total
    ✓ "total" is a synonym of invoice-total, a property of invoice
  [Accept 2]  [Not now]                                 each line can be removed first
```

Every reason is something DataTug read or counted. There is no score.

**Step 4. "I accept what is right and answer what is asked."** She presses both Accept buttons. At once:

- the mark reads `Invoice 3 of 9`, and the header of `Total` shows `invoice.invoice-total · declared`, "confirmed by maria";
- a line under the cards reads `3 confirmed changes are waiting · Review changes`;
- three lines were added to `dbmodels/chinook/chinook.answers.jsonl` in her project;
- nothing was sent anywhere. The database was read, and never written.

**Step 5. "I say something it did not ask."** Nothing pointed at the support representative of a customer. She opens `Customer`, clicks the Advisor control in the header of `SupportRepId`, and sees a list that begins "Refers to": ten tables whose key its three values fit, each with `3 of 3`. She clicks `Employee`. Two clicks. The line under the cards reads `4 confirmed changes are waiting`. From here on DataTug treats her statement exactly as it treats the suggestion she accepted.

**Step 6. "I read one summary of what will change."** She opens "Review changes":

```
Meaning: 2 concepts, 2 bindings, nothing new
  Invoice is an Invoice, the universal concept.
  Invoice.Total is an Invoice total, the universal concept.
    + model/chinook/chinook.meaning.yaml                 new, 31 lines

Model: new, generated from the scan of local: 11 record types, 64 fields, 11 keys; 2 references
  Invoice.CustomerId now refers to Customer.
  Customer.SupportRepId now refers to Employee.
    + model/chinook/chinook.modelspec.hcl                new, 331 lines

Database: 11 tables map to the record type of the same name
    + dbmodels/chinook/chinook.source.yaml               new, 18 lines

Checks
  ✓ the model passes ModelSpec's lint
  ✓ the meaning file passes MeaningGraph's check
  ✓ the source file is read without a problem
  ✓ model and database agree
  ✓ 59 of 59 distinct CustomerId values found in Customer; 3 of 3 distinct SupportRepId values found in Employee

[Keep]
```

Each file line opens to the exact lines that will be written.

**Step 7. "I keep it."** She presses Keep. The three files are written and staged. `git status` lists them and the answers file. She commits. The model, the meaning file, the source file and her answers are files in her project's repository. The line under the cards is gone, and every mark and marker reads as it did before the Keep.

**Step 10. "I open a database I do not know."** Tom opens another copy of the same database in a project of his own. This time the registry lists Chinook's model, and the database page says so:

```
Registered model chinook matches this database: 11 of 11 tables, 64 of 64 columns.   [Adopt]
```

He has adopted nothing yet. The keys are in place, "from the database". On `Customer` the Advisor shows what the registered model states and the database does not, as a question:

```
SupportRepId refers to Employee?  · suggested
    ✓ the registered model chinook states it (not adopted)   published, points
    ✓ all 3 distinct values occur in Employee.EmployeeId      statistics, allows
  [Employee]  [None of these]  [Follow]  [Not sure…]  [Not now]
```

He selects the cell `SupportRepId = 3` and presses Follow. The panel shows the one row of `Employee` whose `EmployeeId` is 3, under the line `via a candidate link, not confirmed`. He is not sure. He presses "Not sure…" and types "looks like a reference to Employee; not sure". One line with `"answer":"open"` is added to the answers file. Whoever knows the database finds the question, his note and its reasons when they open the table.

### One fact, followed: an invoice belongs to a customer

| Where the fact is | Who wrote it | The real lines | How DataTug labels it |
|---|---|---|---|
| The database | nobody | `[CustomerId] INTEGER NOT NULL`, and no `FOREIGN KEY` line | |
| The shape file, `environments/local/catalogs/chinook/chinook.shape.json` | `datatug scan` | `"foreignKeys": []` in the item `Invoice` | |
| The Advisor ledger, `~/.datatug/projects/shop/advisor/chinook/signals.jsonl` | the Advisor | `{"id":"sig_c0d9450fb0f82f32","kind":"metadata","job":"pointer","rule":"same-name-as-key","strength":"strong",…}` and `{"id":"sig_aab0d8bb5309b10e","kind":"statistics","job":"gate","rule":"values-present","result":"passed",…,"counts":{"distinct":59,"present":59},…}` | none: a signal is not a statement |
| A candidate, computed and not stored | the Advisor | `cand_85b83a04f850461e`: `refers-to`, from `chinook/Invoice.CustomerId`, to `chinook/Customer`, level `batch` | `inferred` |
| The answers file, `dbmodels/chinook/chinook.answers.jsonl` | Maria, by pressing Accept | `{"format":"datatug-answers/draft-1","id":"ans_28b988b46e87a6d1","statement":"refers-to","from":"chinook/Invoice.CustomerId","to":"chinook/Customer","answer":"accepted","origin":"suggestion","reasons":[…],"by":"maria","at":"2026-10-12T09:02:11Z"}` | `declared`, "confirmed by maria" |
| The review | the Advisor | `Invoice.CustomerId now refers to Customer.` | |
| The model, `model/chinook/chinook.modelspec.hcl` | DataTug, when Maria pressed Keep | `field "CustomerId" {` / `record   = "Customer"` / `required = true` / `}` | `declared`, "in the model" |
| The meaning file, `model/chinook/chinook.meaning.yaml` | DataTug | no line for this fact. MeaningGraph derives the link from the model once `Customer` has a concept. | |
| The source file, `dbmodels/chinook/chinook.source.yaml` | DataTug | `- Invoice` | |
| The screen, for Maria and for a colleague who clones the project | | Keys tab `CustomerId ⇒ Customer · in the model`; from the cell `CustomerId = 5` of `Invoice`, the panel offers `Customer` (1 row) | |

## Words used here

The words of [`open-layer-reading`](../open-layer-reading/README.md#words-used-here) are used with the meanings it gives them: database, catalog, table, column, model, record type, field, key, reference, meaning file, concept, binding, link, entity file, shape file, source file, mapping, runtime index, resolver, importer, label, mark, registered model, match, drift report. The words below are new.

| Word here | Means | Elsewhere |
|---|---|---|
| Advisor | The part of DataTug that suggests, asks, records answers and writes the three files. It is a section of pages that exist, and one review page. | The proposal: "the Semantic Advisor" (its section 9). |
| statement | One thing that can be said about a table or a column: `refers-to`, `is-a`, `means` or `in-model` (appendix, section 2). | The proposal's figure has `refers-to`, `means`, `maps-to`, `same-values-as` (L2224). |
| fact | A statement that DataTug treats as known. | As `open-layer-reading` uses it. Not the wire type `Fact` of `core-investigation-loop`, which is one value a person carries between queries. |
| signal | One thing one rule found: a name that fits, a count, a concept's word. Not a statement. | The proposal's glossary, L653. |
| pointer | A signal that names a target. It is `decides`, `strong` or `weak`. | The proposal, L655, L2176. |
| gate | A signal that can only rule a statement out: `passed`, `partial`, `failed` or `not-counted`. | The proposal, L655. Not a delivery gate of `ROADMAP.md`, not a reviewer gate of `specscore.yaml`, and not a `Check` of the feature `checks`. |
| value count | The two numbers of the gate `values-present`: how many distinct values a column has, and how many of them occur in the target's key. | The proposal: "a value count per column pair" (L2806). |
| candidate | A statement that a signal is about and that nobody has answered. It is computed, never stored. | The proposal's glossary, L654. Not `apicontract.Candidate`, which is a saved query's runnability. |
| level | What happens to a candidate: `recorded`, `dropped`, `held`, `batch` or `question`. | The proposal, L2195 to L2203. |
| batch, line | A batch is all candidates of one kind, references or meanings, at the level `batch` in one database. A line is one of them. | The proposal's glossary, L656. |
| question, alternative | One card that asks which of at most four statements about one column or table is true. | The proposal's glossary, L656. |
| card | What the Advisor shows: the lines of a batch that belong to the page, or a question. | |
| Advisor ledger | Two append-only files of signals and follows on the person's machine, outside the project. | The proposal: "the ledger". "Ledger" alone already names the Public Data Fabric's programme ledger, its required check, and a billing ledger of `pricing-and-billing`. |
| answer | One line that says what a person decided about a statement: `accepted`, `rejected`, `open` or `withdrawn`. | The proposal's glossary, L657. |
| answers file | The project file that holds a database's answers. | The proposal: "the answers file". |
| note | An answer `open` with the person's words. | The proposal, L1618. Not an annotation of the feature `annotations`. |
| waiting change | An accepted statement that the project's three files do not state yet. | The proposal's mock-up: "3 confirmed changes are waiting" (L2040). |
| plan | The three files as they would be after a Keep, built in memory. | |
| review | The one page that shows the plan as sentences, file changes and checks. | The proposal: "the review diff". |
| Keep | The action that writes the plan's files into the project and stages them. | The proposal: "Keep private" (L2150). |
| checker | ModelSpec's lint or MeaningGraph's check, run on a file. | The proposal: "the existing checkers" (L2807). |
| universal concept | A concept of a graph that the MeaningGraph registry lists with `kind: universal`. Today that is `meaninggraph/core`. | |
| table concept | The entity concept a table holds. | MeaningGraph: the concept bound to the record type with the role `instances`. |
| adopt | To write a source file that names a registered model. | The proposal, L1569, L1613. |

Two labels are added to the four of `open-layer-reading`, with the meanings of proposal D7's card: `inferred` on a candidate that fixed rules alone support, and `suggested` on one that needs a person. Both appear in the Advisor's own routes and nowhere else. On the four routes of `core-investigation-loop`, `inferred` keeps the one meaning it has today: a name pattern of an entity file matched.

The proposal is `/Users/alex/projects/modelspec-conceptual-design-proposal.html`: a local file in no repository, md5 `fb7f6813691fcba8984760c51b3b725c`, cited with its HTML line as `L…`. Its Phase 5 is not Phase 5 of this repository's `ROADMAP.md` ("Dashboardius real persistence", `ROADMAP.md:153`), and not the "Phase 5" that `meaning-inference` names for probabilistic suggestion (`:22`, `:49`). "Proposal D5", "proposal D7", "proposal D8", "proposal D11" and "proposal D12" are the proposal's decision cards, not the decisions of those numbers in `spec/product/product-family.md`.

## Problem

A person who knows what a column means can tell DataTug only by writing a file by hand, and DataTug has nothing to say about a database that declares no foreign keys.

- **There is no way to say one thing about one column.** Meanings come from entity files, written whole (`datatug entity add -f`). No command and no screen sets what one column means or refers to, and the web app cannot save a change to an existing entity (datatug-apps@098db23 `libs/datatug/main/src/lib/pages/signed-in/entity-edit/entity-edit-page.component.ts:161-163`: `// TODO: Implement entity changes save logic`). `open-layer-reading` adds a source file, a model and a meaning file, and a person writes all three.
- **A database without foreign keys shows no related row.** Related rows come from foreign keys and from entity files. Nothing counts values, compares two columns or proposes a reference (a search of every non-test Go file of datatug-core@eb3075e and datatug-cli@3a1c276 finds no count of distinct values pushed to a database and no comparison of two columns).
- **The one rule that guesses today is checked against nothing.** A name pattern of an entity file matches a column's name in any table. In the demo project the field `Album.ID` carries the pattern `ArtistId` (datatug-demo-project@1007ab0 `demo-project-1/entities/Album/Album.entity.json:7-12`), so every `ArtistId` column reads as an album's id. The reason for a match is computed and never sent (`pkg/semantic/resolve.go:35-38`; the response has four keys, datatug-core@eb3075e `pkg/apicontract/responses.go:16-21`).
- **A person's answer has nowhere to go.** Nothing can be accepted, rejected or left for a colleague: the web app has no such control, and the agent has no such route.
- **Nothing writes the open formats.** Proposal D8, which the owner approved, has DataTug write them "directly" in its second stage. No specification for that exists, and neither does one for the ledger and the Advisor that proposal D7 allows.

## Decided by the owner, and where

The owner's words are in the column "His words", and nowhere else in this section. Every other quotation here is the text of a question, a card or the proposal, written by its author.

| Decision | What was put to him | His words | Where it is recorded |
|---|---|---|---|
| Proposal D5, question 3 | "**Q3. A private database.** The manifest requires a public address, so an internal database has nowhere to write Q1 or Q2. Proposed: the same lines in a file inside the DataTug project:" with a six-line example | "Q3 - approved" (2026-10-08, 21:34 UTC) | openvaultdb/openvaultdb@37d2389 `spec/decisions/0012-publisher-manifest-maps-record-types-and-fields-in-the-list-of-recordsets.md:147-156`, `:178`, `:423` |
| Proposal D7 | "Knowledge that is not yet certain stays out of the published formats. DataTug keeps signals and candidates on the user's machine and answers in the project. Five states carry the labels `observed`, `inferred`, `suggested`, `declared` and `published`. Whether a candidate is offered is decided by checks that can rule it out and pointers that name a target, not by a score." (L2961) | "D7: Approve" (2026-10-08, 21:18 UTC) | the proposal's card, L2964; his message in the session record of that day, which is in no repository |
| Proposal D8 | "DataTug's own entity files are replaced by the three open layers, in stages: first an importer that feeds the existing runtime, then DataTug writes the open formats directly. Retiring the entity files and types is a separate decision later." (L2969) | "D8: Approve" (the same message) | the proposal's card, L2972; the same record |
| Proposal D11 | "Five rungs for the Advisor: rules, search and counting, Jev, a generative model, a person. Jev sees names and concept text only. A model's answer can raise a question and can never fill a one-click batch or confirm anything." (L2993) | "D11: Approve" (the same message) | the proposal's card, L2996; the same record |
| Proposal D12 | "Nothing the Advisor learns from a model reaches a user until the benchmark of section 12 has measured it against a bar you set. My suggested bars: no false line in any batch in any of the six databases; recall never below the run with no model; no invented concept where one existed; and nine in ten of a model's questions with the true answer among the buttons." (L3001) | "D12: Approve these bars" (the same message) | the proposal's card, L3004; the same record |
| This specification | A list of what was left after Phase 3, among it "Phases 4 to 8 … Each needs a specification you approve first." | "Do cleanup, specify phase 4 and if make sense others" (2026-10-10) | the session record of that day; it is in no repository |

What each allows, in the proposal's words, which are its author's and not the owner's:

- Proposal D7: "a feature specification for the ledger and the Advisor. Building them is Phase 5." (L2963)
- Proposal D8: "specifications for the importer (Phase 4) and for writing the open formats (Phase 5). It does not authorise removing anything." (L2971)
- Proposal D11: "a specification for the skill and its tools. Building anything that calls a model is Phase 6, and also needs D12." (L2995). It is not among Phase 5's needs. It is here for what it forbids.
- Proposal D12: "these bars as the gate for Phases 5 and 6, and a specification for the benchmark harness. Building the harness is part of Phase 5, after that specification is approved." (L3003). That specification is `advisor-benchmark`.
- The table of phases, for Phase 5: it "Needs D5 question 3, D7, D8, D12. All approved.", allows "Writing the specifications for the ledger, the Advisor and the benchmark harness, with D12's bars as the harness's pass marks" and does not allow "Building them before that specification is approved; anything that calls a model". (L2837)

What he did not decide, so that nobody reads it into the rows above:

- The list of phases, what Phase 5 holds, its "Done when" line and the journey are the proposal author's plan (its section 14). The owner approved cards, not that section.
- The three gates, the ten pointers, the five levels and their order are marked "Proposed" in the proposal (L2206). The 95% line and "at most three cards" are marked "My assumption" (L2206, L2253). The shapes of a signal, a candidate and an answer, the actions and what each writes, and the words on screen are the proposal author's text (its sections 7 and 9). This specification takes some of it and changes some of it; the next section says where.
- That he judged it sensible to specify Phase 5 now. His instruction says "if make sense". The judgment that it does was the coordinating session's.
- His instruction of 2026-10-10 allows writing this specification. It approves nothing in it.

## The author's choices

Each is the choice of this specification's author. None is the owner's ruling. Each can be changed by one word from him. Choices that depart from a sentence of the proposal say so.

| # | The choice | Why | The alternative |
|---|---|---|---|
| 1 | Two features: this one, and `advisor-benchmark` for the harness. Not one, and not three. | Proposal D7's card allows "a feature specification for the ledger and the Advisor", one; proposal D12's allows "a specification for the benchmark harness", another. The ledger, the cards, the answers, the review and the writing share one journey and one end-to-end test. The harness has no screen, runs in CI and is reused by the proposal's Phase 6. The house cuts a programme the same way (`incidents`, `evidence-records`, `annotations`, `checks`, `compare`). | One feature with two appendices, as `open-layer-reading` is one. |
| 2 | The names "Semantic Advisor", "Advisor ledger", "gate", "checker", "Keep". | "Phase 5", "ledger", "check" and "candidate" each already mean something else in this repository. | The proposal's words, qualified every time. |
| 3 | The Advisor ledger holds signals and follows. Candidates are computed from them on every request and stored nowhere. | One source of truth, and nothing to go stale. The proposal's figure stores candidates (L2222); proposal D7's card says only that they stay "on the user's machine", which a computation there does. | A third file of candidates. |
| 4 | A candidate's id is the first 16 hex characters of the SHA-256 of its statement, written as three lines. | The proposal says the id "is computed from the statement, so it is the same on every machine" (L2222) and not how. | A registry of ids. |
| 5 | One answers file per database, `dbmodels/<id>/<id>.answers.jsonl`: JSON Lines, append-only. The last line for a statement is its answer; an undo is a line `withdrawn`. | Two colleagues who each add lines merge without a conflict in the middle of a record. Nothing is lost by an undo. `incidents` keeps `events.jsonl` the same way (`incidents/README.md:98`). | One JSON document that is rewritten, as the proposal's figure may suggest. |
| 6 | Four statements: `refers-to`, `is-a`, `means`, `in-model`. "Takes its values from…", "Holds the same values as…" and the bulk import of the proposal (L1604, L1606) are left out. | They are what a complete model of the six benchmark databases needs. Value sets are also left out by `open-layer-reading` (its open question 12). | All five unprompted statements of the proposal. |
| 7 | "Declared keys" (L2806) is read both ways. A declared foreign key is a pointer that decides by itself: it is shown as `observed`, as today, it is written into a generated model as a reference without a question, and no candidate is made for its column. A table's primary key is what the naming rule and the gate `target-is-key` read. | The proposal does not say which of the two it means, and both are needed. | Only the first. |
| 8 | The naming rule compares whole words, in two forms (appendix, section 3). It is not the probe's "contains the table's name". A key of one word, such as `id`, is never matched by name alone. A table's name is also tried in the singular. | The probe's third clause matched `InvoiceLine.InvoiceLineId` to `Invoice` and leaned on the value count to remove it (L1557). On the six databases the word rule points at no such pair on Chinook, and the one-word guard keeps a database whose keys are all `id`, or all `c1`, from pointing everywhere (the author's calculation, in `advisor-benchmark`). | The probe's three clauses as they are. |
| 9 | A target must be a table's whole primary key of one column. "Or is unique" (L2182) is left for later. | The shape file holds primary keys and no unique index. | Read unique indexes in the scan first. |
| 10 | The value count is exact, over distinct values that are not null, over all rows. It runs only for a pair that a pointer names and that passes the two gates that read no data, and for the list "Refers to…" when a person opens it. One count at a time, each with a limit of 10 seconds. A count is reused until the next scan of that catalog or for 24 hours. | The proposal states the count three ways (59 distinct values, L2007; 412 rows, L2149; a sample, L2336) and its probe compares sets of distinct values (L3136). Counting only pointed pairs is 9 counts on Chinook where the probe made 552 (the author's calculation). 10 seconds is the server's default limit for a query (`core-investigation-loop/api-contract.md:344`). The mock-up's line "The same values also occur in 5 other primary keys" (L2010) is not on the card, because it needs counts for pairs that nothing points at; the list "Refers to…" shows them. | Count rows; or sample; or count every pair. |
| 11 | The 95% line (L2183) is fixed in this feature. | The proposal calls it "a project setting" (L2206). No project file holds settings of this kind today, and adding one is a format of its own. | A setting, once it has a home. |
| 12 | A pair whose values cannot be counted (no value, the time limit, the database unreachable) gets no card. It is listed first under "Refers to…", with the reason. | Then nothing is offered that a count has not allowed, and a name alone never makes a suggestion. `query-workspace` says "Name similarity alone grants no authority" (`:146`). | A question card for it. |
| 13 | Two key columns that hold the same values and point at each other make one question, not two lines. A column with more than four live targets gets no card. | Fixed rules cannot tell the direction. On pubs, `pub_info.pub_id` and `publishers.pub_id` are such a pair. "A question shows at most four buttons." (L2528) | Offer the direction with fewer rows. |
| 14 | The concept rule (appendix, section 4): a table by its name; a column only through a property of its table's concept. Only universal graphs and the database's own meaning file are searched. | The proposal words the rule three ways (L2207, L1576, L2401). A column's name alone finds `unit-price` on `Track.UnitPrice`, where Chinook's published file binds `list-price`: a false line. The rule here finds eight meanings on Chinook, all true by that file, and not `country`, which the proposal's list names (L1576) and its own measure would score wrong (L2542). | The column's name alone; or every registered graph. |
| 15 | A column's meaning needs its table's concept first. DataTug makes no concept that a person did not ask for. | The importer of `open-layer-reading` shows a property only through its owning entity concept (its appendix, section 3). The proposal: "It creates a new concept only when you ask for one" (L1582). Bar 3 of proposal D12 is "no invented concept where one existed". | Generate a concept for the table silently. |
| 16 | A batch is one card per kind. A table's page shows the lines that belong to that table; the database page shows whole batches. A reference belongs to the table of the column that refers. | "Suggestions belong to the table, column or query you opened." (L2253). The proposal's walk accepts "the group of 9" in one action (L1579). | One card per line. |
| 17 | A card about a meaning has Accept, "Not this" and "Not now". It has no Follow. | A meaning is not a link, and there is nothing to follow. The proposal gives four actions to a card "about a link" (L1614) and two buttons to its meaning card (L2024). | Four buttons everywhere. |
| 18 | Follow needs a selected value of the column, and has its own route. | The related-rows contract of `core-investigation-loop` is Stable and has no place to say "not confirmed". Choice 14 of `open-layer-reading` stays as it is. | A new key in the Stable contract. |
| 19 | A note is an answer `open`. It is not an annotation. | The proposal says where it is kept: "an open answer in the project's shared answers file" (L1618). An annotation's target "MUST be semantic, never physical" (`annotations/README.md:89`), and a column of a database nobody has described is physical. | A new target kind in `annotations`, which is Approved. |
| 20 | Waiting changes are derived: the accepted statements that the three files do not state. There is no store of staged changes. | An interrupted Keep then cannot lose or repeat a change: what is not in the files is still waiting. | A staging file. |
| 21 | One function builds the plan, and it is used three times: for what the pages show after an accept, for the review, and for the Keep. | "It is treated exactly like an accepted suggestion" and "the marks fill" before anything is kept (L2752, L2753) then hold by construction, and a Keep changes no mark. | An overlay for the pages and a separate writer. |
| 22 | The model is written as HCL only, with no JSON copy. | HCL is "the intended authored source format" (specscore/modelspec@3c149d2 `spec/hcl-authoring.md:7-10`). The JSON form needs a module id and a version that a private project does not have, and two files can disagree: ModelSpec's lint only warns of a stale copy. `open-layer-reading` reads the HCL when it is the only one named. The proposal lists the model "and its JSON" for a publisher's repository (L1624). | Write both and run `modelspec export --check`. |
| 23 | Files that DataTug creates go to `model/<module>/`. The module's name is made from the database's id. A view gets no record type. Names that are not identifiers get `_`. | One folder is one graph for MeaningGraph's check, so two databases of one project cannot collide on a concept id such as `customer`. A meaning file can name only identifiers (meaninggraph/core@25242a6 `meaning.draft-2.schema.json`), and a database's id may have a hyphen. No registered model has a record type for a view. | The proposal's flat `model/<module>.…` (L1624). |
| 24 | In an existing model or source file DataTug changes only the lines of the change, and refuses the Keep when it cannot. | A person's comments and order are theirs. No released library writes ModelSpec HCL (modelspec-org/cli@21c6314 `pkg/modelspec`). | Rewrite the file in a fixed layout. |
| 25 | A concept that extends a universal one copies its label and its description. | MeaningGraph requires both on every concept (`FORMAT.md:112-126`). The proposal's mock-up of the written lines has neither (L2124 to L2129). | Ask the person for a description each time. |
| 26 | Keep writes and stages. It never commits. | The house rule: "DataTug stages, the user commits" (`git-native-projects/README.md:81`). The journey's "files in my project's repository" (L2755) holds for files in the working tree. | Commit through the reviewed commit operation of `git-native-projects`, once it exists. |
| 27 | The review is one new page. Its file changes are text differences that DataTug computes. | It is the first form of the diff view that `git-native-projects` plans for "an edited mapping" (`:120-125`). It adds sentences and checks, and versions nothing itself. | A section of the database page. |
| 28 | Four checks block a Keep. The fifth, the value count of each reference, is shown and blocks nothing. | "A gate that fails beside it is reported as a finding; it does not overrule the database or a person." (L2199) | Block on all five. |
| 29 | Every action that writes needs the agent's write capability. Without it the buttons are shown disabled with the reason. The person named in an answer is the principal the agent serves. | The Stable contract refuses project writes without it (`core-investigation-loop/api-contract.md:392-394`). `checks` asks for "absent or explicitly disabled with the reason shown" (`:170-171`); `annotations` for "never a client-supplied name" (`:98`). | Hide the Advisor on a read-only agent. |
| 30 | A table that the serving principal may not read whole gets no count and no suggestion, and no reason names it. | The same side as choice 26 of `open-layer-reading`: "a count must not tell what a policy hides". Under a row policy the count a person may see is not the count that exists. | Count what the principal may see. |
| 31 | An accept never writes an entity file. A column that an entity file maps gets no meaning suggestion and no meaning statement. Name patterns are not turned into rules. | Proposal D8: DataTug "writes the open formats directly", and its approval "does not authorise removing anything". The card's reason says the name patterns "become fixed pointer rules" (L2970); that belongs with retiring entity files, which is a later decision. | Accept into the entity file, as `meaning-inference` AC:accept-promotes-to-declared says. |
| 32 | The references of a matching registered model that nobody adopted are shown as `suggested`, read from the registry's index. Adopt is one confirmation that writes the source file. | Journey step 10: "what the database and any matching model state is in place" (L2758), and L1613. The index already lists each model's references, so nothing more is fetched. Its meanings are not in the index and come with adoption. | Show nothing of the model until it is adopted. |
| 33 | Commands under `datatug advisor`, and six routes under `advisor/`. | "Each tool is a CLI command first" (L2504). The benchmark's scripted user needs verbs to call. The name is free. | More verbs under `datatug db`. |

## Questions for the owner

Six questions. Each changes what is built. The requirements below are written for the recommended answer, and each one that would change names its question.

### Question 1. Where is the Advisor ledger kept?

Proposal D7's card, which you approved, says signals and candidates stay "on the user's machine". The proposal's text says more: "in a ledger folder inside the project that Git ignores" (L2240). The two places:

```
~/.datatug/projects/shop/advisor/chinook/signals.jsonl      A: outside the project, beside personal queries
./shop/.datatug/advisor/chinook/signals.jsonl                B: inside the project, with a line added to .gitignore
```

- **A. Outside the project (recommended).** It is where personal queries already live, by your answer of 2026-09-11 "I'm Ok with your suggestion for now" (`executable-knowledge-library/README.md:109-131`), under the same setting `DATATUG_PERSONAL_DIR`. DataTug then never edits a project's `.gitignore`, and a mistake cannot commit a count. It departs from the sentence at L2240 and keeps the card's words.
- B. Inside the project, ignored by Git. As the proposal's text says. DataTug must then write `.gitignore`, and a project folder that is not a Git repository has no such protection.

### Question 2. May a suggestion be used before a person accepts it?

The line `CustomerId refers to Customer` in step 3 is `inferred`: fixed rules alone support it. The proposal lets such a candidate drive a lookup or a join at once, "always with its label, as DataTug does today" (L2168). Then, before Maria accepts anything, a click on the cell `CustomerId = 5` of `Invoice` would already offer:

```
Customer (1 row) · inferred
```

- **A. No (recommended).** Until a person accepts it, a candidate lives in the Advisor's cards only. It is in no marker, no related row and no join, and the four routes of the Stable contract never carry it. A person can use it once, through Follow, and the result says "not confirmed". `query-workspace`, which you approved on 2026-10-09, says "name matches alone are not relationships" (`:162`); a batch line is more than a name, but one rule for everything unaccepted is the simplest thing to explain and to test.
- B. Yes, for `inferred` only. The panel offers the row with the label. The closed label lists of the Stable contract stay as they are, since `inferred` is in them, but its meaning there widens from "an entity file's name pattern matched" to "any fixed rule". A false line then shows as a related row before anyone has looked at it.

### Question 3. In which format does DataTug write a meaning file?

You approved `meaning/draft-2` (decisions 0001 to 0003 of `meaninggraph/core`). No released checker reads it yet: "no released reader does yet … Published graphs stay in `meaning/draft-1` until the readers are released" (meaninggraph/core@25242a6 `FORMAT.md:21-27`). The reader is an open pull request (meaninggraph/cli#11, open on 2026-10-10). One binding in each:

```yaml
kind: property          # meaning/draft-2
field: Total
role: value
```
```yaml
kind: attribute         # meaning/draft-1
property: Total
role: value
```

- **A. `meaning/draft-2` (recommended).** DataTug writes the words you approved. Building this feature waits for a released checker that reads them. A graph in one format may pin a graph in the other (`FORMAT.md:63-66`), so the file can extend today's `core`.
- B. `meaning/draft-1` now. The released checker v0.3.0 checks it today. Every file DataTug writes then carries the words you approved replacing, and is rewritten later.

### Question 4. When the database follows a registered model, where does an accepted statement go?

Tom adopts the registered model. The source file names it by address, and its facts are `published`. He then accepts `SupportRepId refers to Employee`. Suppose the model did not state it. The model is its publisher's, not his, and the source file of `open-layer-reading` cannot name a registered model and a file of the project together (its REQ:source-file-names-model).

- **A. It stays an answer (recommended).** The line is in the project's answers file. Inside the project it is a fact, `declared`, "confirmed by tom · not in the model". No file is written for it, and it is never a waiting change. The proposal says the same: "DataTug keeps it as your answer and uses it for navigation and joins inside your project" (L1908). Meanings cannot be accepted or stated for such a database in this feature: a person can leave a note. Saying more about a published model needs a file of the project beside it, which is the `modules` line of the proposal's Phase 7.
- B. DataTug offers to copy the registered model and its meaning file into the project and to follow the copy. Then everything can be written, the facts turn from `published` to `declared`, and the copy no longer follows its publisher. Choice 8 of `open-layer-reading` argues against a copy.

### Question 5. May the three files keep their own formats, in folders of their own?

`dalgo-project-store`, which is Approved, fixes one layout for project items, `<key>/<key>.<type-suffix>.json`, and records your decision "writes should be through dalgo" (`:172`, `:531`). The files of this feature:

```
model/chinook/chinook.modelspec.hcl        ModelSpec's format, HCL
model/chinook/chinook.meaning.yaml         MeaningGraph's format, YAML
dbmodels/chinook/chinook.source.yaml       the source file of open-layer-reading, YAML
dbmodels/chinook/chinook.answers.jsonl     DataTug's own, JSON Lines
```

- **A. Yes (recommended).** The model and the meaning file are files of two other standards, and their checkers read them only as they are. The source file is what `open-layer-reading` reads. They are stored as they are, as dataset records already are ("the one deliberate exception", `:160-161`). The answers file is a log and stays JSON Lines. All four are still written only by the server, staged and never committed.
- B. Each is wrapped in a JSON record of the canonical layout. Then no other tool can read the model or the meaning file from the project, and every Keep also has to unwrap them for the checkers.
- C. Only the answers file follows the canonical layout, one record per answer.

### Question 6. May DataTug count values when a person opens a table?

Step 3 needs the count `59 of 59` before the card can be shown. Opening `Invoice` then runs one count in the database. Opening a table of AdventureWorks runs up to a few. `pricing-and-billing` records your words "Viewing metadata is free. Editing or executing data queries is active usage." (`:113`) and counts a query run "without requiring an explicit Run button" as activity (`:75`).

- **A. Yes, with the limits of choice 10 (recommended).** Only pairs that a name points at, one at a time, 10 seconds each, reused for 24 hours. The journey asks for it: "I open a table I know … at most three suggestions" (L2751). Whether such a count is billable activity is that feature's to say; this one only names it.
- B. Only when the person presses "Look for references". Then step 3 shows the meanings card and a button, and the references card after one more click.

## If the questions of open-layer-reading are answered otherwise

`open-layer-reading` puts six questions to the owner. None is answered. This feature is written for the recommended answer to each.

| Its question | The other answer | What changes here |
|---|---|---|
| 1. A column in the source file as a bare string | B or C: the bare form is read | Nothing. REQ:source-file-written writes the map form, which every answer reads, and never rewrites a line a person wrote. |
| 2. Offer to adopt a matching model | B: Phase 4 has the confirmation and a write route | REQ:adopt is then built there. Here it shrinks to "the Advisor's Adopt button calls it", and AC:unadopted-model-and-adopt loses its last clause. REQ:pointer-unadopted-model stays. |
| 3. `published` on the Stable contract's label lists | B: imported meanings are sent as `declared` | Nothing in a requirement. This feature sends `declared` for an accepted statement on those routes and nothing else new. A fact from a registered model and a fact a person accepted then read alike there, and only the Advisor's routes and the database page tell them apart. |
| 4. PostgreSQL foreign keys in the shape file | B: a PostgreSQL shape file says `not-read` | REQ:pointer-declared-key has nothing to read on PostgreSQL. A reference the database declares is then offered as a suggestion, and its card says "foreign keys were not read for this database" where it says "no foreign key is declared". A generated model has no reference from a foreign key. The end-to-end test and the benchmark use SQLite either way. |
| 5. The shape file beside or in place of today's files | B: in place | Nothing. This feature reads the shape file only. |
| 6. A second project, or the demo project changed | B: the demo project loses its Chinook entity files | REQ:entity-file-columns matters less, since the demo has no entity file for Chinook then. Nothing changes in a requirement. This feature's tests make their own projects either way. |

## Behavior

Requirements use MUST, MUST NOT, SHOULD and MAY as fixed terms. "Appendix" means [formats.md](formats.md).

### What stays fixed

#### REQ: reading-unchanged

Every requirement of `open-layer-reading` MUST keep holding, with the changes the table "What this changes in open-layer-reading" names and no other. The change that builds this feature MUST NOT alter a non-test file of `datatug-core` `pkg/semantic`. The four routes `GET semantic/columns`, `POST semantic/related`, `POST semantic/related/rows` and `POST queries/applicable` MUST keep their paths, their request shapes, the keys of their responses and their closed lists of values. A project with no answers file and no Advisor ledger MUST look and answer as it does without this feature, apart from the Advisor's cards.

#### REQ: no-model-no-network

Nothing in this feature MUST call an AI model, a decision engine or a search service. The only requests that leave the machine are those `open-layer-reading` already allows: a GET of a registry's `index.json`, and a GET of a file at a commit. No request MUST carry the name of a table, a column, a database or a concept that a person chose, nor a count.

#### REQ: database-read-only

The value count is the one thing in this feature that opens the database. It MUST open it read-only and create nothing in it. No value of a row MUST be written to the Advisor ledger, to the answers file, to a log or to any response of an Advisor route other than `POST advisor/follow`, whose result is rows the person asked for.

### Signals and the Advisor ledger

#### REQ: ledger-place

Signals and follows MUST be kept in the Advisor ledger, in the folder the appendix, section 1, gives: under the personal folder, by project id and database id, outside the project (Question 1). DataTug MUST NOT write a signal, a candidate or a count into any project file, and MUST NOT change a project's `.gitignore`.

#### REQ: ledger-append-only

The two files of the Advisor ledger MUST be written by appending lines of the shapes of the appendix, section 1. No line is changed or removed. A signal's line is appended only when the signal is new or its state differs from its last line, so a second run on an unchanged database appends nothing. Deleting the folder MUST lose no answer and no fact: the next request makes the signals again.

#### REQ: signals-refreshed

The signals of the rules that read the shape file, and of the concept rule, MUST be computed again whenever the shape file, the registries' indexes or the database's meaning file have changed since they were computed. A value count MUST be reused until the next scan of that catalog or for 24 hours, whichever comes first, and counted again on `datatug advisor suggest --recount`.

### Pointers

#### REQ: pointer-declared-key

A foreign key of one column that the shape file holds is a pointer of strength `decides`. Its statement is at the level `recorded`, is shown as `observed` exactly as `open-layer-reading` shows it, and is never a card. No `refers-to` candidate MUST be made for a column that already has a reference from the database, from the model the database follows or from an accepted answer. When the shape file says `"foreignKeys": "not-read"` or `"partly-read"`, a column with no foreign key in the file is treated as having none known: candidates are made, and their cards say that foreign keys were not read.

#### REQ: pointer-names

A `refers-to` candidate from a column to a table gets a `strong` pointer when one of the two naming rules of the appendix, section 3, fires: `same-name-as-key` or `table-name-then-key`. The rules read the shape file only, compare whole words, and fire only toward a table whose primary key is one column. No other comparison of names is a pointer.

#### REQ: pointer-exact-concept

An `is-a` candidate for a table and a `means` candidate for a column get a `strong` pointer by the two concept rules of the appendix, section 4, and by no other: `concept-names-table` and `concept-names-column`. The concepts searched are those of the universal graphs of the MeaningGraph registry's index, each at its listed commit or at the pin the database's meaning file carries, and those of the database's own meaning file. They are read from files on the machine. A rule that finds two concepts points at neither with strength `strong`; the two are the alternatives of one question.

#### REQ: pointer-unadopted-model

When a registered model matches the database, by REQ:match-definition of `open-layer-reading`, and the database has no source file, each reference that the registry's index lists for that model is a `weak` pointer for the `refers-to` statement it names, with the model's address and commit as its origin. Nothing else of the model is used, and no file of the model is fetched for it.

### Gates

#### REQ: gate-types-comparable

A `refers-to` candidate MUST be dropped unless the referring column and the target column have the same known type family, by the type tables of `open-layer-reading`.

#### REQ: gate-target-is-key

A `refers-to` candidate MUST be dropped unless its target column is the whole primary key of the target table.

#### REQ: gate-values-present

The gate `values-present` MUST be decided by the value count of the appendix, section 5: `passed` when every distinct value of the referring column that is not null occurs in the target's key and there is at least one; `partial` when at least 95% and not all occur; `failed` below 95%; `not-counted` when there is no value, when the count did not finish within its limit, or when the database could not be opened. The two numbers MUST be shown wherever the gate is a reason.

#### REQ: value-count-cost

A value count MUST run only for a candidate that has a pointer and has passed the two gates that read no data, and for the tables of the list "Refers to…" when a person opens it (Question 6). Counts MUST run one at a time for one database, each in a read-only session through the agent's read path, each with a limit of 10 seconds. A count MUST NOT run for a table that the serving principal may not read whole.

#### REQ: gates-for-meanings

An `is-a` or `means` candidate MUST be dropped when its gate `concept-free` fails (appendix, section 5): the table already has a table concept, the concept is another table's, the column already has a written binding or an accepted meaning, or an entity file maps the column. The concepts searched carry no data type, so no gate compares a column's type with a concept in this feature.

### Candidates and levels

#### REQ: candidate-identity

A candidate MUST be computed, on every request, from the current signals, the shape file and the answers file, and MUST NOT be stored. Its id MUST be computed from its statement as the appendix, section 2, gives it, so that the same statement has the same id on every machine and after every scan.

#### REQ: level-rule

Every candidate MUST get exactly one level by the table of the appendix, section 5, applied in its order: `recorded`, `dropped`, `held`, `batch`, `question`. A candidate is at the level `batch` only when every gate has passed in full, a `strong` pointer names it, it is the only live candidate of its kind for its column or table, and it is not half of a mutual pair. No number of `weak` pointers makes a candidate a batch line. A question has at most four alternatives.

#### REQ: no-score

DataTug MUST NOT compute, store, send or show a score, a confidence or a probability for a candidate. What a person sees is the level, the label and the reasons. The only numbers are counts.

#### REQ: unaccepted-stays-in-advisor

A candidate that no person has accepted MUST NOT enter the runtime index, and MUST NOT appear in a response of the four routes of REQ:reading-unchanged, in a header marker, in a mark, in the Keys tab, in related rows, in a join suggestion or in an applicable query's explanation (Question 2). It appears in the Advisor's cards and routes, with its label, and can be used through Follow.

### What a person sees

#### REQ: cards

The table page and the database page MUST show an Advisor section with at most three cards. A table's page shows the batch lines and the questions whose `from` is that table or one of its columns; the database page shows the whole batches and every question. Order: a question before a batch; then a card about the selected column; then references before meanings; then by name. Each line MUST show its sentence and every reason, each reason with its kind and its job, in the words of the appendix, section 2.

#### REQ: batch-card

A batch card MUST let a person accept all its lines with one action, remove any line first, and reject a line. Accepting a `means` line whose table concept is itself a line of the same batch accepts that line too, and the card says so before the action.

#### REQ: question-card

A question card MUST show its alternatives as buttons, at most four, and the buttons "None of these" and "Not now". Choosing an alternative accepts that statement. "None of these" rejects every alternative. A question whose statement has an answer `open` MUST show the note and who left it.

#### REQ: quiet

The Advisor MUST show nothing about a table or a database that the person has not opened. It MUST NOT raise a notification or a modal, and MUST NOT show a count of suggestions. The one count it shows is the number of waiting changes.

#### REQ: policy-hides-suggestions

A candidate whose table, or whose target table, the serving principal may not read whole MUST be in no card and no response, and no reason and no list MUST name such a table or a count from it. The test is the one `open-layer-reading` uses for its routes.

### Actions and answers

#### REQ: answers-file

Every answer MUST be one line appended to `dbmodels/<id>/<id>.answers.jsonl`, with the keys of the appendix, section 6. Each line is complete in itself: it repeats the statement and the reasons that were shown. DataTug MUST NOT change or remove a line. The current answer to a statement is the last line that names it. `by` is the principal the agent serves; a request that carries an author or a time MUST be refused.

#### REQ: answer-recognised-again

An answer applies to the statement it names, by the statement's three parts. After any number of scans, on any machine and with an empty Advisor ledger, a candidate for the same statement MUST get that answer: an accepted one is a fact, a rejected one is not offered, an open one is a question with its note. An answer whose table or column the shape file no longer has MUST stay in the file, be used for nothing, and be reported as a finding on the database page.

#### REQ: accept

Accepting a line or choosing an alternative MUST append one line with `"answer":"accepted"` and `"origin":"suggestion"`, with the reasons then shown, at once and before the response.

#### REQ: reject

Rejecting MUST append one line with `"answer":"rejected"`. A candidate whose current answer is `rejected` MUST NOT be offered in that project, by any rule, until the answer is withdrawn.

#### REQ: not-now

"Not now" MUST write nothing: no line in the answers file and none in the Advisor ledger. The card is hidden until the page is opened again.

#### REQ: follow

Follow is offered on a line or a question about a `refers-to` candidate, and needs a selected cell of the referring column. It MUST run one lookup of the target table's rows whose key equals that value, through the same read path, policy and limits as `POST semantic/related/rows`, and MUST show the rows under the line "via a candidate link, not confirmed". It MUST write nothing in the project, MUST append one line to `follows.jsonl`, and MUST NOT be read by any rule as a signal.

#### REQ: state-unprompted

A person MUST be able to make each of the four statements without a suggestion. On a column: the first click opens the Advisor control in the column's header, which lists at once "Refers to" with every table whose key the column's type fits, ordered by the share of values present and then by name, each with its count, and "Means" with the concepts the concept rule allows and a search of the concepts of REQ:pointer-exact-concept. The second click chooses. For a column whose table has no table concept, "Means" first asks what the table is. On a table: "Is a" with the same search, and "Add to the model". A stated statement MUST append one line with `"answer":"accepted"` and `"origin":"stated"`, and from then on MUST be treated exactly as an accepted suggestion: the same fact, the same label, the same waiting change. A new concept is offered only at the end of a search, takes a label and a one-line description from the person, and is written by no other path.

#### REQ: note

On any card, and on any row of the lists of REQ:state-unprompted, a person MUST be able to leave a note: one line of text, appended as an answer `open` with the statement and its reasons. A statement whose current answer is `open` MUST be shown to everyone who opens the project as a question with the note, whether or not their own rules point at it.

#### REQ: undo

Every answer MUST be undoable from the column's or the table's Advisor control, which lists its answers. An undo appends a line `withdrawn`. Withdrawing an accepted statement that the files already state makes a waiting change that removes it. A Keep is undone by withdrawing statements, or with Git; DataTug offers nothing else.

#### REQ: write-capability

Accept, reject, an answer to a question, a statement, a note, an undo, Keep and Adopt are project writes. Each MUST be refused, with nothing written, unless the agent has its write capability and the serving principal is allowed the write, by the rules the Stable contract gives for a query write and against the resource `/datatug_projects/{project}/advisor/{database id}`. Each request carries an operation id, so that a retry gives the first result and writes nothing twice, and names the answer it saw as current; when another answer has become current it is refused with `REVISION_CONFLICT`. On an agent without the write capability the buttons MUST be shown disabled with the reason, and the cards, Follow and "Not now" MUST work.

#### REQ: accepted-is-declared-at-once

From the moment its line is appended, an accepted statement MUST be a fact with the label `declared`: the marks count it, the Keys tab lists an accepted reference, related rows follow it, and a column with an accepted meaning shows its marker. Until a Keep the words are "confirmed by" and the person's name; after it, the words of `open-layer-reading`. Nothing a person sees, other than those words and the line of waiting changes, MUST change at the Keep.

#### REQ: entity-file-columns

A column that an entity file maps with a declared mapping MUST get no `means` candidate, and a `means` statement for it MUST be refused with a sentence that names the entity file. No action of this feature MUST create, change or remove an entity file. `refers-to` and `is-a` are not affected.

#### REQ: conflict-shown

When a rejected `refers-to` is later declared by the database as a foreign key, DataTug MUST show the foreign key as `observed`, as for any other, and one finding that names the rejection and who made it. When a gate fails for a statement a person accepted, DataTug MUST keep the statement and show one finding with the gate and its count. It MUST NOT remove or change an answer in either case.

### Waiting changes, the plan and the review

#### REQ: waiting-changes-derived

The waiting changes of a database MUST be computed, not stored: the statements whose current answer is `accepted` and that the project's model, meaning file and source file do not state, and the statements those files state because of an answer that is now `withdrawn`. Their number MUST be shown under the cards as "N confirmed changes are waiting", with the link to the review.

#### REQ: plan-function

One function MUST build the plan: from the project's three files as they are, the shape file and the waiting changes, the three files as they would be. It writes nothing. The pages MUST show what the importer of `open-layer-reading` reads from the plan; the review MUST show the difference between the files and the plan; the Keep MUST write the plan. A waiting change that would make a blocking check fail is left out of the plan and listed as "cannot be kept yet" with the finding.

#### REQ: review-summary

The review MUST be one page that shows, for one database: the changes as sentences grouped under the three layers Meaning, Model and Database, in that order; under each layer the files that change, each opening to its exact difference; the checks, each with its result; and one button, Keep. A layer with no change says so in one line. A concept that extends no other is marked "new".

#### REQ: review-checks

The review MUST run the five checks of the appendix, section 8, on the plan: ModelSpec's lint on the model; MeaningGraph's check on the meaning file's folder, with each pinned graph supplied from the cache; the source file's own reader; the drift report; and the value count of every reference the plan writes. The first two run as the reference Go libraries, in the agent, with no network. The first four block a Keep when they fail. The fifth is shown with its counts and blocks nothing.

#### REQ: keep

Keep MUST write the plan's files and nothing else, as one operation: each file whole, by replacing it, in an order in which the project passes the four blocking checks after each replacement. Each file is written only if it is still at the revision the review was made from, or still absent. The files written MUST be staged in the project's Git repository, when it is in one, and MUST NOT be committed. After a Keep that was interrupted at any point, every file is whole, the project passes the four checks, and what was not written is still a waiting change. A file that an interrupted Keep left at a path the source file does not name yet is replaced by the next Keep.

#### REQ: keep-refused

When a blocking check fails, or a file has changed since the review was made, Keep MUST write nothing and MUST say which check or which file. The answers stay as they are, and the changes stay waiting.

### What a Keep writes

#### REQ: model-file

The model MUST be ModelSpec HCL in the words `record` and `field`, by the rules of the appendix, section 7. When the database has no model, the plan generates one from the shape file: a record type for each table that is in the model, a field for each column, the key, the types of the appendix's table, and a reference field for each declared foreign key of one column and each accepted `refers-to`. The first waiting change of a database with no model puts every table in the model. In a model that exists, DataTug changes only the lines of the change and keeps every other byte; when it cannot make the change that way, the change is "cannot be kept yet" and the review shows the lines to write by hand.

#### REQ: meaning-file

The meaning file MUST be written in `meaning/draft-2` (Question 3) by the rules of the appendix, section 7: a concept of the project that extends the universal concept and is bound to the record type or the field; the label and the description copied from the universal concept; one pin for every reference to another graph; the role each statement takes. A concept with no `extends` is written only for a statement whose `to` is a new concept. A `refers-to` writes no line in the meaning file.

#### REQ: source-file-written

The source file MUST be written in the form "In the project" of `open-layer-reading`, with a column as a map with the key `field`. When the database has none, the plan creates it; when it has one, DataTug adds only an item or the `meaning` key and keeps every other line. A scan MUST still leave the file byte for byte as it was.

#### REQ: written-files-read-back

The files a Keep writes MUST be read by the importer of `open-layer-reading` with no restart of the agent, and give, for every kept statement, the same entity field, the same reference and the same label `declared` as the plan gave before the Keep. A colleague who clones the project and scans nothing MUST see the same.

#### REQ: registered-model-database

For a database whose source file names a registered model, an accepted or stated `refers-to` MUST be a fact inside the project, `declared`, with the words "not in the model"; it is never a waiting change and no file is written for it. `is-a` and `means` MUST NOT be offered or accepted for such a database; a note can be left (Question 4).

#### REQ: adopt

When a registered model matches and the database has no source file, the model line and the Advisor MUST offer Adopt for it. Adopt writes the source file with exactly the text that `datatug db match --source-file <address>` prints, stages it, and writes nothing else. An accepted `refers-to` that the adopted model also states stays in the answers file and is shown once, as `published`.

### Commands and routes

#### REQ: cli-verbs

The commands of the appendix, section 10, MUST exist under `datatug advisor`, with those flags, outputs and exit codes. Each action a person can take on a page has a command, and a command writes what the page's action writes.

#### REQ: routes

The agent MUST serve the six routes of the appendix, section 9, with those bodies: `GET advisor/cards`, `POST advisor/answers`, `POST advisor/follow`, `GET advisor/review`, `POST advisor/keep` and `POST advisor/adopt`. The two GET routes and `POST advisor/follow` need no write capability. No response has a key or a value that is a score.

### End-to-end gate

#### REQ: journey-test

One automated test MUST walk steps 3 to 7 and 10 of the journey on the stripped Chinook, against a real agent started with its write capability and with no intercepted request, in six stages, in the order and with the results of "What a person sees": open `Invoice` and see the two cards; accept both; state `SupportRepId refers to Employee` in two clicks; read the review; Keep, commit, and open the project as a second person on a fresh clone; and, in a second project in which the registered model matches, see the question, Follow it and leave a note. No stage may need a reload of the page, a restart of the agent, a second scan or a file edited by hand. `git commit` and `git clone` are run by the test as the two people.

#### REQ: done-when-is-the-benchmark

This feature MUST NOT be called built until `advisor-benchmark` AC:done-when passes on the same commits. That criterion is the proposal's "Done when" line for its Phase 5 (L2809).

## What this changes in open-layer-reading

`open-layer-reading` is a Draft. Its page is not edited here. If both are approved, these clauses of it read as follows.

| Its clause | Says | With this feature |
|---|---|---|
| REQ:source-file-place | "In this feature DataTug MUST NOT create, change, move or delete it" | Keep and Adopt create it and add to it. A scan still leaves it byte for byte. |
| REQ:match-is-local-and-read-only | A match "MUST NOT make any statement of the matching model appear on a screen" | Its references appear in the Advisor's cards as `suggested`, and nowhere else. A match still creates and changes no file. |
| REQ:label-definitions | "It MUST NOT use the proposal's fifth label, `suggested`, in this feature." | Used here, in the Advisor's routes only. AC:labels-on-columns stays true of every other route. |
| REQ:database-never-written | "The importer, the marks, the match and the drift report MUST NOT open the database at all" | Still true of those four. The value count is a fifth thing and opens the database read-only. |
| REQ:importer-inputs | The importer is "a pure function of four inputs" | Unchanged. It is called on the plan's files where a database has waiting changes. |
| REQ:no-new-screen | "The web app MUST gain no page and no address." | One page more, the review. |
| AC:fixed-interface, REQ:read-only-routes | "the agent's routes have grown by exactly three, all GET" | By six more. Three of them write, behind the write capability. |
| Choice 17 | "Value counts are the next phase's." | They are here. The drift report still counts nothing. |

## What must not change

| What | How to test that it did not |
|---|---|
| The database | Its SHA-256 is the same before and after the whole walk; every count runs in a read-only session (AC:counts-leave-the-database-alone). |
| That nothing leaves the machine | With every outgoing request recorded: only GETs of an index or of a file at a commit, and none carries a table's, a column's or a chosen concept's name (AC:nothing-leaves-the-machine). |
| The registries | DataTug opens no pull request and writes to no registry; the same record shows no request other than a GET. |
| What `open-layer-reading` specifies | Its criteria pass unmodified on a project with no answers file, and the four routes answer the same bodies there (AC:reading-stays-as-it-is). |
| The resolver and the Stable contract | The implementing change touches no non-test file of `pkg/semantic`, and no closed list of the four routes gains a value (AC:reading-stays-as-it-is). |
| Entity files | Nothing under `entities/` is created, changed or removed by any action, and `datatug entity list` prints the same (AC:entity-file-column). |
| The launch lane's work | No file of `datatug/datatug-demo-project` changes; no file of the web app's `public-data` route changes; journeys J1 to J4 pass unmodified beside the walk (AC:journey-walk). This feature asks for no change to `git-native-projects` or to `query-workspace`. |
| A project on a read-only agent | Every write is refused and nothing is written (AC:read-only-agent). |

## Existing features

Line numbers of this repository's own pages are as at datatug/datatug@4aa2b66. "Depends on" means this specification relies on the feature as it is. "Extends" means it adds to what the feature covers; the feature's own page is not edited here and is listed for a later amendment by its owner. "Leaves" means no change either way.

| Feature | Status | Relation | What follows |
|---|---|---|---|
| `open-layer-reading` | Draft | Depends on; extends | See the table above. This feature cannot be approved before it. |
| `core-investigation-loop` | Stable | Depends on | The four routes and their lists stay as they are. Its write rules (`api-contract.md:392-405`) are the model for REQ:write-capability; they name query resources only, so the resource here is new. Its `agent-info` does not say whether writes are allowed (`:543-548`); `GET advisor/cards` says it for the Advisor. |
| `meaning-inference` | Draft | Extends | It owns suggestions. Two of its clauses do not fit what the owner approved in proposal D7 and D8, and are listed for its owner: "Every such suggestion MUST carry a confidence signal" (`:65`), where D7 says "not by a score"; and AC:accept-promotes-to-declared, "the entity file is updated" (`:101`), where D8 has DataTug write the open formats. This feature builds none of its REQ:probabilistic-meaning-suggestions. |
| Plan "Phase 2, library depth and provenance", Tasks 6 and 7 (sneat-co/workbench@1602d0d `spec/plans/github.com/datatug/datatug/2026-09-09-phase-2-library-depth-and-provenance/README.md:145-165`) | Approved; both tasks planning | Overlaps | They plan accept and correct into entity files, and a review UI for it. If this feature is approved they are its owner's to retire or to re-aim. |
| `related-data-navigation` | Draft | Extends | Its REQ:automatic-relationship-discovery asks that discovery by value overlap "MUST stay an experiment validated for precision" (`:94`) and be told apart from a declared relationship. Here the value count rules out and never proposes; the benchmark is the validation; and nothing unaccepted is a relationship. Its label `discovered` (`:134`) and this page's `suggested` and `inferred` are listed for its owner to settle on one word, beside `fk` and `observed`. |
| `provenance` | Draft | Extends | An answer line is the record of "who or what produced it and when" for a mapping (`:60-65`). "Every mapping this user corrected" is a filter of the answers file by `by` and `answer`. Listed for a later amendment that names the file. |
| `semantic-layer` | Draft | Depends on | Entities and entity files stay. After this feature a project has two places where a person can declare a meaning. A new statement goes to the open formats, by proposal D8. |
| `git-native-projects` | Implementing | Depends on | REQ:git-stage-mutations: stage, never commit. REQ:one-project-mutation-contract covers "project creation and saved queries" today (`:140`). This feature needs a write of up to three files that are not queries as one operation; it states its own rules (REQ:keep) and asks that feature for nothing. If its contract has grown to other assets when this is built, the Keep goes through it. |
| `dalgo-project-store`, `openvaultdb-project-store` | Approved | Depends on; extends | Every write goes through the server. Four file kinds are added (Question 5). |
| `server-enforced-access` | Draft | Depends on | Counts are reads of rows and pass the same policy path. REQ:policy-hides-suggestions is the rule for what a policy narrows. |
| `executable-knowledge-library` | Draft | Depends on | Its personal folder, and the owner's recorded answers about it (`:109-131`), are where the Advisor ledger sits (Question 1). |
| `annotations` | Approved | Leaves | A note here is an answer, not an annotation (choice 19). |
| `checks`, `compare` | Approved | Leaves | A gate, a checker and the review's checks are not a `Check`. `datatug check` and `datatug compare` are not used. |
| `query-workspace` | Approved | Leaves | "Name similarity alone grants no authority" (`:146`) holds: a name only makes a candidate, and no candidate is a relationship until a person accepts it. An accepted reference is then a declared reference of the project. |
| `evidence-records` | Approved | Leaves | A value count is not recorded as an execution record by this feature. Whether a project policy "record always" should cover it is that feature's to say. |
| `pricing-and-billing` | Draft | Leaves | A value count executes a data query without a Run button (Question 6), and an answer is one of the "other project mutations" that feature leaves undecided (`:113`). This page decides neither. |
| `surfaces` | Draft | Leaves | "No surface holds private state the others cannot see" (`:14`): the Advisor ledger is the agent's, reached by the web app and the commands alike, and holds nothing that a scan cannot make again. |
| `agent-context` | Draft | Leaves | An AI agent gets accepted facts through the same resolver. It gets no candidate, no skill and no new tool here; those are the proposal's Phase 6. |
| `public-data-fabric` | Draft | Leaves | Another lane. "Candidate state MUST stay separate from published canonical registries" (`:42`) is proposal D7's first sentence. Nothing under it is edited. |

No page under `spec/features` other than `open-layer-reading` names ModelSpec's record types or fields, so none is listed here for the earlier words `entity` and `property`.

## Dependencies

- open-layer-reading
- core-investigation-loop
- meaning-inference
- related-data-navigation
- provenance
- git-native-projects
- dalgo-project-store
- server-enforced-access
- executable-knowledge-library

## Blocked Clauses

- **REQ:meaning-file and the second check of REQ:review-checks** need a released Go checker that reads `meaning/draft-2`. `github.com/meaninggraph/cli/pkg/meaning` v0.3.0 reads `meaning/draft-1` only (meaninggraph/cli@e5d4d8f `pkg/meaning/schema.go:24`); the change is the open pull request meaninggraph/cli#11. Under answer B of Question 3 nothing blocks.
- **REQ:written-files-read-back** needs what blocks `open-layer-reading`: a Go reader that derives links.
- **REQ:value-count-cost** asks that a count run through the agent's policy-enforced read path. `dal-go/dalgo` v0.93.0, which `datatug-cli` requires, has `CountDistinctAs` (`dal/q_functions.go:196`). Whether one structured query can count the values of one column that occur in another table, on SQLite and on PostgreSQL, was read and not run.
- Nothing blocks **REQ:model-file**: `github.com/modelspec-org/cli/pkg/modelspec` v0.3.0 lints `1.0-draft-2` as a library and touches no network. It writes no HCL, so the writer is DataTug's own.

## Acceptance Criteria

Unless a criterion says otherwise: "the stripped Chinook" is the SQLite file that `advisor-benchmark` REQ:stripped-copy builds from demo-db/chinook@88f080d72205034e1161bb80594acb7d77fc58ef `data-source/source.sqlite` (SHA-256 `7651ba378ac2fcd0dfc3c66fb101f7a7eed3ba39a612ec642b96e20702061f15`), with its original names. "The project" is a folder in a Git repository with one commit, made by `datatug scan --directory ./shop --driver sqlite3 --path chinook-nofk.sqlite --db chinook --env local`, served by an agent started with its write capability as the principal `maria`. The two settings of `open-layer-reading` REQ:registry-index-source name index files in which no model and no graph of Chinook is listed ("the masked indexes" of `advisor-benchmark`), and `core` is listed at `dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7`. "After the walk" means after steps 3 to 5 of "What a person sees": three suggestions accepted and one reference stated.

### AC: reading-stays-as-it-is (verifies REQ:reading-unchanged, REQ:unaccepted-stays-in-advisor)

**Given** the project with nothing accepted, and the Chinook project of `open-layer-reading` with its Example A
**When** the four routes are called for `Invoice` and `Customer` on the first, and every criterion of `open-layer-reading` runs on the second
**Then** on the first, `GET semantic/columns` returns no column, `POST semantic/related` for `CustomerId = 5` of `Invoice` returns no item, and no response contains `suggested`, `inferred` or the id of a candidate; on the second every criterion passes unmodified; and the implementing change alters no non-test file of `pkg/semantic`.

### AC: nothing-leaves-the-machine (verifies REQ:no-model-no-network)

**Given** the project, with every outgoing request of the agent and of the commands recorded
**When** the whole of "What a person sees", steps 3 to 7, is done
**Then** every recorded request is a GET of a registry's `index.json` or of a file at a commit; no URL, header or body contains `Invoice`, `CustomerId`, `SupportRepId`, `59` as a count or `invoice-total`; and no request went to a model, a decision engine or a search service.

### AC: counts-leave-the-database-alone (verifies REQ:database-read-only, REQ:value-count-cost)

**Given** the SHA-256 of the stripped Chinook taken before the scan, and the agent's log of the queries it sent to the database
**When** `Invoice` and then each of the other ten tables is opened
**Then** the file's SHA-256 is unchanged; the log shows nine counts in all, each in a read-only session, none running while another runs, and none for a pair that no naming rule points at; and neither the Advisor ledger nor any response of `GET advisor/cards` contains a value of a row.

### AC: ledger-lines (verifies REQ:ledger-place, REQ:ledger-append-only, REQ:signals-refreshed)

**Given** the project, with `DATATUG_PERSONAL_DIR` naming an empty folder
**When** `datatug advisor suggest` runs, then runs again, then the folder is deleted and it runs a third time
**Then** after the first run `<that folder>/shop/advisor/chinook/signals.jsonl` holds, among others, the two lines of the appendix, section 1, apart from `at`; the second run appends no line and sends no count to the database; the third run writes the same signals again and prints the same cards; and `git status` in the project shows no change after any of them, `.gitignore` included.

### AC: chinook-references (verifies REQ:pointer-names, REQ:gate-types-comparable, REQ:gate-target-is-key, REQ:gate-values-present, REQ:level-rule, REQ:candidate-identity)

**Given** the project
**When** `datatug advisor suggest --format json` runs
**Then** the batch `references` has exactly nine lines: `Album.ArtistId` to `Artist`, `Invoice.CustomerId` to `Customer`, `InvoiceLine.InvoiceId` to `Invoice`, `InvoiceLine.TrackId` to `Track`, `PlaylistTrack.PlaylistId` to `Playlist`, `PlaylistTrack.TrackId` to `Track`, `Track.AlbumId` to `Album`, `Track.GenreId` to `Genre` and `Track.MediaTypeId` to `MediaType`; the line for `Invoice.CustomerId` has the id `cand_85b83a04f850461e`, the label `inferred` and the reason "all 59 distinct values occur there"; there is no question; and no card names `Customer.SupportRepId`, `Employee.ReportsTo` or `InvoiceLine.InvoiceLineId`.

### AC: names-that-must-not-point (verifies REQ:pointer-names)

**Given** three scanned databases: the stripped Chinook; a copy of it in which every table is named `t1` to `t11` and every column `c1`, `c2` and so on within its table; and a database of two tables `orders` and `customers` whose keys are both named `id`, where `orders` also has a column `customer_id`
**When** `datatug advisor suggest --format json` runs on each
**Then** the first has no signal of a naming rule for `InvoiceLine.InvoiceLineId`; the second has no signal of a naming rule at all and no card about a reference; and the third has exactly one, `table-name-then-key`, from `orders.customer_id` to `customers`, and none between the two `id` columns.

### AC: declared-key-decides (verifies REQ:pointer-declared-key)

**Given** a project made by scanning Chinook with its eleven foreign keys, and a second made by scanning the stripped Chinook with a test scanner that reads no foreign keys, so that its shape file says `"foreignKeys": "not-read"`
**When** `datatug advisor suggest` runs on each, and `datatug advisor state in-model --all` and `datatug advisor review` run on the first
**Then** the first has no card about a reference, and its planned model has eleven reference fields, among them `record = "Customer"` on `Invoice.CustomerId`, with no answer about a reference in the answers file; the second has the nine lines of AC:chinook-references, each with the reason "foreign keys were not read for this database".

### AC: values-in-part-and-missing (verifies REQ:gate-values-present, REQ:level-rule)

**Given** two more stripped copies of Chinook, built by the test: in one, two of the 59 customers that invoices name are removed from `Customer`; in the other, four
**When** `datatug advisor suggest --format json` runs on each
**Then** in the first, `Invoice.CustomerId` to `Customer` is a question with the label `suggested`, one alternative and the reason "57 of 59 distinct values occur there", and it is not a line of the batch; in the second it is in no card, and its signal `values-present` says `failed` with the counts 59 and 55.

### AC: count-not-possible (verifies REQ:gate-values-present, REQ:value-count-cost)

**Given** the stripped Northwind of `advisor-benchmark`, whose table `CustomerCustomerDemo` has no row, scanned
**When** `datatug advisor suggest --format json` runs, and the Advisor control of `CustomerCustomerDemo.CustomerID` is opened
**Then** no card names `CustomerCustomerDemo.CustomerID`; its signal `values-present` says `not-counted` with `"why":"no-values"`; and the list "Refers to" shows `Customers` first with the words "no values to count".

### AC: one-column-several-targets (verifies REQ:level-rule, REQ:question-card)

**Given** the stripped pubs of `advisor-benchmark`, scanned
**When** `datatug advisor suggest --format json` runs
**Then** `titles.pub_id` is one question with the two alternatives `pub_info` and `publishers` and is in no batch; `pub_info.pub_id` and `publishers.pub_id` are one question with three buttons, "pub_info refers to publishers", "publishers refers to pub_info" and "None of these"; and no question has more than four alternatives.

### AC: chinook-meanings (verifies REQ:pointer-exact-concept, REQ:gates-for-meanings)

**Given** the project
**When** `datatug advisor suggest --format json` runs
**Then** the batch `meanings` has exactly eight lines: `is-a` for `Customer`, `Employee`, `Invoice` and `InvoiceLine` with the concepts `customer`, `employee`, `invoice` and `invoice-line` of `core`, and `means` for `Invoice.Total` (`invoice-total`), `InvoiceLine.UnitPrice` (`unit-price`), `InvoiceLine.Quantity` (`quantity`) and `Employee.ReportsTo` (`manager`); and no card names `Track.UnitPrice`, `Customer.Country`, `Employee.Country`, `Customer.City` or `Customer.Company`.

### AC: no-score-anywhere (verifies REQ:no-score, REQ:routes)

**Given** the project
**When** `GET advisor/cards` is called for the database and for `Invoice`, and `GET advisor/review` after the walk
**Then** no response has a key named `score`, `confidence`, `probability` or `rank`; every number in a response is a count of values, lines, files or waiting changes; and the page shows no percentage.

### AC: cards-on-invoice (verifies REQ:cards, REQ:batch-card, REQ:quiet)

**Given** the project, opened in the web app
**When** the table `Invoice` is opened, and then the database page
**Then** the table page shows exactly two cards: `References · 1 line · inferred` with the line `CustomerId refers to Customer`, its four reasons each with kind and job, and four buttons; and `Meanings · 2 lines · inferred` with `Invoice is an Invoice` and `Total is an Invoice total`; the database page shows two cards, of nine and eight lines; no notification, modal or count of suggestions appears on any page; and a page of another table shows nothing about `Invoice`.

### AC: policy-hides (verifies REQ:policy-hides-suggestions)

**Given** the project with a policy that lets the role `support` read `Invoice` whole and `Customer` for one country only, served by a second agent started `--as agent --role support`
**When** `GET advisor/cards` is called for `Invoice` on that agent, and the Advisor control of `Invoice.CustomerId` is opened
**Then** no card and no list names `Customer`, no count from `Customer` is in any response, and the agent's log shows no count that reads `Customer`.

### AC: accept-is-a-fact-at-once (verifies REQ:accept, REQ:answers-file, REQ:accepted-is-declared-at-once)

**Given** the project with the two cards of `Invoice` shown
**When** Accept is pressed on both
**Then** `dbmodels/chinook/chinook.answers.jsonl` has exactly the first three lines of the appendix, section 6, apart from `at` and the ids that depend on it, and is staged; with no reload, the mark reads `Invoice 3 of 9`, the header of `Total` shows `invoice.invoice-total · declared` with "confirmed by maria", the Keys tab lists `CustomerId ⇒ Customer` with "confirmed by maria", `POST semantic/related` for `CustomerId = 5` of `Invoice` returns one item, `Customer`, count 1, and `GET semantic/columns` for `Invoice` returns `InvoiceId` as `invoice.id` and `Total` as `invoice.invoice-total`, both `"provenance": "declared"`; and no file under `model/` exists.

### AC: read-only-agent (verifies REQ:write-capability)

**Given** the project served by an agent started without its write capability
**When** `Invoice` is opened, Follow is used on the reference with the cell `CustomerId = 5` selected, and `POST advisor/answers`, `POST advisor/keep` and `POST advisor/adopt` are called directly
**Then** the two cards are shown with Accept and "Not a relationship" disabled and the words "this agent was started without --allow-writes"; Follow shows one row of `Customer`; each of the three POSTs answers `ACCESS_DENIED`; and no file of the project has changed.

### AC: answers-carry-no-claimed-author (verifies REQ:answers-file, REQ:write-capability)

**Given** the project
**When** `POST advisor/answers` is called with a body that carries `"by":"someone"`, then with the same operation id twice for one valid accept, then with `ifCurrent: null` for a statement that already has an answer
**Then** the first is refused with nothing written; the second writes one line, and both calls return that line's id; the third answers `REVISION_CONFLICT` with nothing written.

### AC: reject-and-undo (verifies REQ:reject, REQ:undo, REQ:answer-recognised-again)

**Given** the project
**When** "Not a relationship" is pressed on `CustomerId refers to Customer`; the Advisor ledger is deleted and the database is scanned again; then the rejection is undone from the column's Advisor control
**Then** after the first, the answers file has one line `rejected` and the card is gone; after the second, `datatug advisor suggest` still does not offer that line, and offers the other eight; after the third, the file has a second line, `withdrawn`, the first line is unchanged, and the line is offered again with the same id.

### AC: not-now-writes-nothing (verifies REQ:not-now)

**Given** the project, with the bytes of the answers file and of the Advisor ledger's two files noted
**When** "Not now" is pressed on both cards of `Invoice`, and the page is opened again
**Then** all three files have the same bytes, or are still absent; the cards were hidden after the press and are shown again.

### AC: follow-says-unconfirmed (verifies REQ:follow)

**Given** the project with nothing accepted, and the cell `CustomerId = 5` of `Invoice` selected
**When** Follow is pressed on `CustomerId refers to Customer`, and `datatug advisor suggest` runs afterwards
**Then** the panel shows one row of `Customer` under the line "via a candidate link, not confirmed", and the body of `POST advisor/follow` has `"via":"candidate"` and `"confirmed":false`; `follows.jsonl` has one new line that does not contain `5`; no file of the project has changed; and the line's reasons are the four it had before.

### AC: state-in-two-clicks (verifies REQ:state-unprompted)

**Given** the project with `Customer` open
**When** the Advisor control in the header of `SupportRepId` is clicked, and then `Employee` in its list
**Then** after the first click the list "Refers to" shows ten tables, each with `3 of 3`, in the order of their names; after the second, the answers file has the fourth line of the appendix, section 6, with `"origin":"stated"`; the mark reads `Customer 2 of 13`; `POST semantic/related` for `SupportRepId = 3` of `Customer` returns `Employee`, count 1; the waiting changes have grown by one; and `datatug advisor state refers-to Customer.SupportRepId Employee` on a fresh copy of the project writes the same line.

### AC: note-is-an-open-answer (verifies REQ:note, REQ:answer-recognised-again)

**Given** the project, served as the principal `tom`, who leaves the note "looks like a reference to Employee; not sure" on the row `Employee` of the list of `Customer.SupportRepId`, and the change is committed
**When** a second person clones the project, starts an agent with an empty Advisor ledger and opens `Customer`
**Then** the answers file has one line with `"answer":"open"` and that note; the second person sees a question `SupportRepId refers to Employee?` with the label `suggested`, the note and "tom"; and the statement is not a fact: the mark reads `Customer 1 of 13`.

### AC: entity-file-column (verifies REQ:entity-file-columns)

**Given** a copy of `demo-project-1` in which the stripped Chinook was scanned as the catalog `chinook`
**When** `datatug advisor suggest --format json` runs, and `datatug advisor state means Customer.CustomerId customer` is tried
**Then** no line of the batch `meanings` names `Customer.CustomerId`; the command exits 1 with a sentence that names `entities/Customer/Customer.entity.json`; `datatug entity list` prints the same eight entities as before; and nothing under `entities/` has changed.

### AC: conflict-is-shown-not-settled (verifies REQ:conflict-shown)

**Given** the project in which `CustomerId refers to Customer` was rejected, and a second project in which `Customer.SupportRepId refers to Artist` was stated
**When** the first database is replaced by Chinook with its foreign keys and scanned again; and in the second, the rows of `Artist` whose key is 3, 4 or 5 are removed and the count is made again
**Then** in the first, the Keys tab lists `CustomerId ⇒ Customer`, "from the database", and the database page has one finding that names the rejection and "maria"; in the second, the statement is still a fact, with one finding "0 of 3 distinct values occur in Artist.ArtistId"; and neither answers file has lost or changed a line.

### AC: waiting-and-plan (verifies REQ:waiting-changes-derived, REQ:plan-function)

**Given** the project after the walk
**When** `GET advisor/cards` and `GET advisor/review` are called, then the stated reference is undone and both are called again
**Then** the first pair says `"waiting": 4`, and the review's planned model has `record = "Employee"` on `Customer.SupportRepId`; the second pair says `"waiting": 3`, and that field has `type = "int"`; no file other than the answers file was written; and no file or key anywhere holds a list of waiting changes.

### AC: review-of-the-walk (verifies REQ:review-summary, REQ:review-checks)

**Given** the project after the walk
**When** the review is opened
**Then** it shows three layers in the order Meaning, Model, Database with the headings and sentences of "What a person sees", step 6; three files marked new, of 31, 331 and 18 lines, each opening to its text; five checks, all passed, the fifth reading "59 of 59 distinct CustomerId values found in Customer; 3 of 3 distinct SupportRepId values found in Employee"; no concept marked "new"; and one button, Keep, enabled.

### AC: keep-writes-three-files (verifies REQ:keep, REQ:model-file, REQ:meaning-file, REQ:source-file-written, REQ:written-files-read-back)

**Given** the project after the walk
**When** Keep is pressed, and then, with no restart, the four routes are called as in AC:accept-is-a-fact-at-once; and then the project is committed, cloned, and opened by a second agent that scans nothing
**Then** `model/chinook/chinook.modelspec.hcl` has 11 record types and 64 fields, with `record = "Customer"` on `Invoice.CustomerId` and `record = "Employee"` on `Customer.SupportRepId` and no other reference; `model/chinook/chinook.meaning.yaml` is the 31 lines of the appendix, section 7; `dbmodels/chinook/chinook.source.yaml` is the 18 lines of the appendix; all three are staged and none is committed by DataTug; the released `modelspec lint` exits 0 on the model's folder; the four routes answer exactly as before the Keep; the waiting changes are 0; and the second agent shows `Invoice 3 of 9`, the marker `invoice.invoice-total · declared` and `CustomerId ⇒ Customer · in the model`.

### AC: keep-refused-on-a-failed-check (verifies REQ:keep-refused, REQ:review-checks)

**Given** the project after the walk and a Keep, in which a person then edits `model/chinook/chinook.meaning.yaml` by hand so that its binding names the field `Totl`, and accepts `Customer is a Customer`
**When** the review is opened and Keep is pressed
**Then** the check of the meaning file is shown failed, with the checker's finding, the file and the line; Keep is disabled, and `POST advisor/keep` called directly answers with that finding; no file is written; and the statement is still an answer `accepted` and still waiting.

### AC: kept-only-what-can-be (verifies REQ:plan-function, REQ:meaning-file)

**Given** the project after the walk, in which `Customer is a Customer` is then accepted and a person states that `Customer.SupportRepId` means a new concept, with the label "Support representative" and a description
**When** the review is opened and Keep is pressed; then `Employee is an Employee` is accepted and the review is opened again
**Then** the first review lists the new concept under "cannot be kept yet" with the sentence "Say what an Employee is first.", passes its checks without it, and Keep writes the other five changes and leaves `"waiting": 1`; the second review plans a concept `support-representative` of kind `property` with `of: customer`, `values-of: employee` and no `extends`, bound to `Customer.SupportRepId` with the role `reference`, and marks it "new".

### AC: keep-meets-a-changed-file (verifies REQ:keep, REQ:keep-refused)

**Given** the project after the walk and a first Keep, with one more reference and one more meaning accepted, `Album.ArtistId refers to Artist` and `Customer is a Customer`, and the review open
**When** a line is added by hand to `model/chinook/chinook.modelspec.hcl` and Keep is pressed; and, in a second run, the agent is stopped after it has replaced the model and before it has replaced anything else
**Then** in the first run Keep answers `REVISION_CONFLICT` naming that file, nothing is written, and a new review shows the person's line kept; in the second run, after a restart, every file is whole, the released `modelspec lint` exits 0, and `GET advisor/cards` reports exactly the changes that were not written as waiting.

### AC: edit-keeps-a-persons-lines (verifies REQ:model-file, REQ:source-file-written)

**Given** a project with a model `model/shop.modelspec.hcl` of two record types, `Customer` and `Invoice`, written by hand with comments, in which `Invoice.CustomerId` has `type = "int"`, and a hand-written source file that names it
**When** `Invoice.CustomerId refers to Customer` is stated and kept
**Then** the difference of the model is exactly one line, `type = "int"` replaced by `record = "Customer"` with the `=` signs of that field in one column; every comment and every other byte is as it was; the source file is unchanged; and the released `modelspec lint` exits 0.

### AC: names-and-types-generated (verifies REQ:model-file, REQ:source-file-written)

**Given** the stripped Northwind of `advisor-benchmark`, scanned, with its 17 views
**When** `datatug advisor state in-model --all` and `datatug advisor keep` run
**Then** the model has 13 record types and none for a view; the table `Order Details` is the record type `Order_Details`, and the source file has for it a map with `name: Order Details` and `record_type: Order_Details`, while `Orders` is a bare line; `Categories.Picture`, a `BLOB`, is a field of type `any`; the drift report has no finding of severity `error` or `warning`; and the released `modelspec lint` exits 0.

### AC: new-concept-only-when-asked (verifies REQ:meaning-file, REQ:state-unprompted, REQ:review-summary)

**Given** the project with `Track` open
**When** "Is a" is opened, the search for "track" returns nothing, and the person chooses "New concept", types the label "Track" and the description "One recording that the shop sells."; and the review is opened
**Then** the answers file has one line `is-a` whose `to` is `local:track` and whose `concept` has that label and description; the review's Meaning layer has the heading "1 concept, 1 binding, 1 new" and marks `track` "new"; the planned meaning file has a concept `track` of kind `entity` with no `extends`; and in no run of `datatug advisor suggest` on any fixture of `advisor-benchmark` does a card propose a concept that is not one of those searched.

### AC: registered-model-keeps-an-answer (verifies REQ:registered-model-database)

**Given** a copy of the stripped Chinook to which the test has added a column `Customer.ReferredBy`, an integer that holds the `CustomerId` of another customer in ten rows, scanned, with Example A of `open-layer-reading` as its source file, served with the write capability as `tom`
**When** `datatug advisor state refers-to Customer.ReferredBy Customer` runs, and then `datatug advisor state is-a Track --new Track --description x`
**Then** the first writes one line to the answers file; the Keys tab of `Customer` lists `ReferredBy ⇒ Customer` with "confirmed by tom · not in the model"; `GET advisor/cards` says `"waiting": 0`; no file other than the answers file is written; and the second exits 1 with a sentence that says the database follows a registered model.

### AC: unadopted-model-and-adopt (verifies REQ:pointer-unadopted-model, REQ:adopt)

**Given** a project made by scanning the stripped Chinook, with the two settings naming the index files of `modelspec-org/registry@a7325ea` and `meaninggraph/registry@49f472f`, which list Chinook
**When** `datatug advisor suggest --format json` runs, and then `datatug advisor adopt modelspec://github.com/demo-db/chinook/chinook`
**Then** before the adoption the batch `references` has the nine lines of AC:chinook-references, each with one more reason, "the registered model chinook states it (not adopted)", and there are two questions, `Customer.SupportRepId refers to Employee?` and `Employee.ReportsTo refers to Employee?`, each with the label `suggested` and that reason; no file of the model was fetched; after it, `dbmodels/chinook/chinook.source.yaml` has exactly the text of Example A of `open-layer-reading`, is staged, and no card about a reference is left.

### AC: commands-do-what-pages-do (verifies REQ:cli-verbs, REQ:routes)

**Given** two fresh copies of the project
**When** the walk is done in the web app on one, and on the other by `datatug advisor accept --batch references --except` the eight other lines, `accept` of the two meaning lines of `Invoice`, `state refers-to Customer.SupportRepId Employee`, `review --format json` and `keep`
**Then** the two answers files are equal apart from `at`, `by` and the ids that depend on them; the three kept files are byte for byte equal; `review --format json` printed the body of `GET advisor/review`; and each command exited 0.

### AC: journey-walk (verifies REQ:journey-test, REQ:done-when-is-the-benchmark)

**Given** an empty Git repository, the stripped Chinook kept outside it, a `datatug-cli` release that contains this feature, the masked indexes for the first project and the indexes of `modelspec-org/registry@a7325ea` and `meaninggraph/registry@49f472f` for the second
**When** the test runs its six stages in order, with no reload of a page, no restart of an agent, no second scan and no file edited by hand
**Then** stage 1 shows the two cards of AC:cards-on-invoice; stage 2 shows `Invoice 3 of 9` and "3 confirmed changes are waiting"; stage 3 takes two clicks and shows "4 confirmed changes are waiting"; stage 4 shows the review of AC:review-of-the-walk; stage 5 leaves the four files in the repository, and the second person on a fresh clone sees what AC:keep-writes-three-files says; stage 6 shows the model line with Adopt, the question about `SupportRepId`, one row of `Employee` under "via a candidate link, not confirmed", and one line `open` in the second project's answers file; the SHA-256 of the SQLite file is unchanged; the recorded requests are as in AC:nothing-leaves-the-machine; journeys J1 to J4 pass beside it; and the test's report names the run of `advisor-benchmark` AC:done-when on the same commits and its result.

## Rehearse Integration

Every criterion has a surface that a test can observe: a file's bytes, a line of JSON, a route's body, a command's output and exit code, a line on a page, a query in the agent's log. This Feature spans three repositories, and each criterion is exercised against a real implementation: the rules, the plan and the writers in `datatug-core`, the counts, the routes and the commands in `datatug-cli`, the cards, the review and the end-to-end walk in `datatug-apps`. Stub scaffolding is therefore deferred to those implementations, as `open-layer-reading` and `openvaultdb-project-store` do.

## Not Doing / Out of Scope

- Anything of the proposal's Phases 6 to 8: a saved query as a pointer, query lineage, Jev, a generative model, the skill, an ORM mapping or code as a signal, links between two databases, publishing, a pull request to a registry, a publisher's manifest.
- A pointer from a concept found by resemblance, from a query, or from anything a model said.
- The panel "What this query shows" and the chat card of the proposal's section 9. Both need what Phase 6 brings.
- Retiring entity files, turning their name patterns into rules, or writing an entity file. Proposal D8 keeps the first for a later decision.
- A commit by DataTug, a branch, a merge or a pull request for the project's own repository.
- The hosted web app for a project opened with no agent. It has no Advisor, as it has no meanings in `open-layer-reading`.
- inGitDB collections and HTTP sources.
- A record type for a view; a reference to a key of several columns; a target that is unique and not the primary key.
- A score.

## Open Questions

The questions that change what is built are in [Questions for the owner](#questions-for-the-owner). The points below are left to a later phase by this specification's author. They are listed so that nobody assumes an answer.

For later phases:

1. "Takes its values from…", "Holds the same values as…" and value sets (the proposal, L1604), and the bulk import of what a person already has (L1606).
2. A saved query that joins two columns as a pointer (L2187). It calls no model, and the proposal still places query lineage in its Phase 6 (L2812).
3. A reference to a key of several columns. AdventureWorks has one such foreign key.
4. A target that is unique and not the primary key, once a scan reads unique indexes.
5. Meanings for a database that follows a registered model, and references between two models: the `modules` line of the proposal's Phase 7.
6. The 95% line as a project setting, and where a project keeps such settings.
7. Turning the name patterns of entity files into pointer rules, with the decision to retire entity files.
8. A second meaning on one column, and measures. `open-layer-reading` shows one meaning and imports no measure; a person can state a second here and it is written, and only the first is shown.
9. Engines that fold the case of names. The rules here compare words in lower case; the source file's mapping stays exact.
10. Committing the kept files through the reviewed commit operation of `git-native-projects`, once it exists.
11. Showing a candidate to an AI agent, and letting an agent answer one. The proposal's Phase 6.
12. Whether a value count is an execution record, and whether it is billable activity. Each is another feature's to say.

Not about what is built:

13. When building may start. The proposal places its product phases "after the foundation, and after the launch" (L2796). This specification does not decide it. Building it changes the agent's route table, the table page and the journey job, which the launch also runs on. The author's suggestion is to start after the launch, after `open-layer-reading` is built, and with `advisor-benchmark` first, since `ROADMAP.md:257` asks to "measure precision and correction cost before UI".

---
*This document follows the https://specscore.md/feature-specification*
