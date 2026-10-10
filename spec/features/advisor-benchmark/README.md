---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Advisor Benchmark

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/advisor-benchmark?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/advisor-benchmark?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/advisor-benchmark?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/advisor-benchmark?op=request-change) |
**Status:** Draft
**Date:** 2026-10-10
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Kind:** Capability (implementation in `datatug-cli`)
**Pillars:** P4 meaning inference · P5 related-data navigation

## Summary

Specifies, and does not yet build, the harness that measures the Semantic Advisor's fixed rules on six databases with their foreign keys removed, against the pass marks the owner approved.

This is the specification of the last part of what the ModelSpec design proposal calls "Phase 5, DataTug writes, with fixed rules only": "The benchmark harness for A and B as a no-model run in CI, with the registry masked and both renamed variants" (L2808). Its one large criterion, AC:done-when, is that phase's "Done when" line. It measures the sibling feature [`semantic-advisor`](../semantic-advisor/README.md) and has no screen. Nothing here is built, and none of it is part of the launch planned for 2026-10-25. The status is Draft. Only the owner approves it. No AI model is called by anything in this feature.

How to read it: "What a person sees" is one run. "Questions for the owner" holds the four answers it needs; the first reports that the approved pass mark cannot be met on AdventureWorks as it is worded. "The author's choices" lists 14 things the author chose. The rest is for whoever builds it.

## What a person sees

Nobody opens a page. A developer changes a rule of the Advisor in `datatug-cli` and opens a pull request. A job named `advisor-benchmark` runs. It takes the six databases, removes their foreign keys, lets the Advisor work on each as if nobody had described it, and plays a person who knows every answer. For Chinook the job's log reads:

```
advisor-benchmark · no model · datatug-cli 4f1c2ab · seed 20261010
chinook · original names · 11 hidden references
  references   batch 9 lines: 9 true, 0 false     questions 0     stated by the scripted user 2
  meanings     batch 8 lines: 8 true, 0 false     new concepts proposed 0
  actions      reference actions 3 · meaning actions 19 · keep 1 · total 23
  complete     yes: 11 keys, 11 references, 25 bound columns, 10 table concepts
chinook · opaque names · 11 hidden references
  references   batch 0 lines                       questions 0     stated by the scripted user 11
  …
pass marks
  1  no false line in any batch            18 of 18 runs     PASS
  2  recall not below the baseline         151 references    PASS  (batch 119, question 20, stated 12)
  3  no invented concept                   0                 PASS
  4  questions with the true answer        20 of 22          reported, not a pass mark in this run
```

The reference numbers in this picture, and the meaning lines, are what the author's calculation gives for the rules of `semantic-advisor` on the pinned files (see "What the rules give on the six databases"). The number of meaning actions is the author's estimate from Chinook's published meaning file, which binds 10 tables and 25 columns, and is there to show the form. The commit is made up. Nothing has been run through DataTug.

If a rule change adds a line that the truth does not hold, the job fails and names it. For example:

```
FAIL pass mark 1: adventureworks · original names · references batch
  Production.Document.Owner refers to Person.Person      not a hidden reference, and no verdict in truth/adventureworks.references.yaml
```

The developer sees which line, in which database, under which names. Nothing is published, and nothing is sent anywhere.

### One fact, followed: an invoice belongs to a customer

| Where | The real lines |
|---|---|
| The fixture, demo-db/chinook@88f080d `data-source/source.sqlite` | SQLite reports a foreign key of `Invoice`: `CustomerId` to `Customer.CustomerId` |
| The truth of benchmark A, read from the fixture before it is copied | `Invoice.CustomerId → Customer` |
| The stripped copy, original names | `CREATE TABLE "Invoice" ("InvoiceId" INTEGER NOT NULL, "CustomerId" INTEGER NOT NULL, …, PRIMARY KEY ("InvoiceId"))`, and no `FOREIGN KEY` |
| The stripped copy, renamed with the seed 20261010 | the same table under names such as `Lantern`, `LanternId`, `AnvilId`: every word but `id` replaced, the same word always by the same word |
| The stripped copy, opaque names | `t6` with the columns `c1`, `c2`, … |
| The Advisor's batch, original names | `Invoice.CustomerId refers to Customer` · 59 of 59 |
| The harness | the line is in the truth: true. The scripted user accepts the batch: one action for nine lines. |
| The kept model | `field "CustomerId" { record = "Customer" … }` |
| The report | `chinook · original · reached at the level batch: Invoice.CustomerId` |

## Words used here

The words of [`semantic-advisor`](../semantic-advisor/README.md#words-used-here) and of [`open-layer-reading`](../open-layer-reading/README.md#words-used-here) are used with their meanings there. The words below are new.

| Word here | Means | Elsewhere |
|---|---|---|
| harness | The program of this feature: it builds the copies, runs the Advisor, plays the scripted user, compares with the truth and writes the report. | The proposal: "the benchmark harness". |
| benchmark A, benchmark B | A: references, with the foreign keys hidden. B: meanings, with the database's own meaning file withheld. | The proposal's section 12 (L2516, L2517). Its benchmark D and its tests C and E belong to its Phase 7. |
| fixture | One of the six SQLite files, at a commit, with its SHA-256. | |
| stripped copy | A copy of a fixture with every foreign key removed, and its primary keys kept. | The proposal: "A copy of the database with every foreign-key constraint removed. Primary keys stay." (L2516) |
| variant | One of three ways a stripped copy is named: original, renamed, opaque. | The proposal: "Two renamed variants" (L2527). |
| run | One database in one variant, from the scan to the kept model. There are 18. | |
| hidden reference | A foreign key of one column that the fixture declares. There are 151. | The proposal: "the 151 references that are one column each" (L2539). |
| truth | What the harness holds right: the hidden references, the bindings of the withheld meaning file, and the verdicts a person added. | The proposal: "hidden truth". |
| verdict | One line a person adds to the truth about a statement the withheld files do not settle: true, false, or a concept. | The proposal: "Truth is reviewed before it is used." (L2529) |
| false line | A line of a batch whose statement the truth says is false. | Proposal D12: "no false line in any batch". |
| scripted user | The part of the harness that knows the truth and acts by one fixed policy, through the commands of `datatug advisor`. | The proposal, L2528. |
| action | One call of a command that changes the answers file, or of `keep`. | The proposal: "Each of those is one action." (L2528) |
| complete model | The state in which every key, every hidden reference and every binding of the truth is in the kept files. | The proposal, L2545. |
| pass mark | One of the four bars of proposal D12, as this feature applies it. | The proposal: "D12's bars as the harness's pass marks" (L2837). |
| baseline | The committed record of which hidden references the no-model run reaches, and at which level. | The proposal: "the no-model run, which is rungs 1 and 2 alone" (L2540). |
| masked index | A registry's index file with the six databases' entries removed. | The proposal: "A masked registry" (L2526). |

The proposal is `/Users/alex/projects/modelspec-conceptual-design-proposal.html`: a local file in no repository, md5 `fb7f6813691fcba8984760c51b3b725c`, cited with its HTML line as `L…`. "Proposal D12" is its decision card, not decision D12 of `spec/product/product-family.md` ("Assertion kinds separate fact, claim, hypothesis, inference", `:227`). The job of this feature is not the job `ledger` of the workflow "Public Data Fabric ledger".

## Problem

DataTug is about to suggest what columns mean and refer to, and nothing measures whether a suggestion is right.

- **A wrong suggestion that a person accepts with one click becomes a file in their project.** The Advisor offers a batch for one action. The owner approved a bar for it: "no false line in any batch in any of the six databases" (proposal D12). No program checks a batch against anything.
- **The one measurement that exists is a probe of one database.** The proposal ran simple rules on Chinook and reports 9 candidates, 9 true (L1550). It ran nothing on the other five, on renamed names or on meanings (L2560: "Not built").
- **The answers are published.** The models and meaning files of all six databases are registered. An Advisor that reads the registry would find the whole truth in one index entry and seem to infer it.
- **The house already asks for this.** `ROADMAP.md:257` says of relationship and meaning suggestions: "measure precision and correction cost before UI". `meaning-inference` says a suggestion "MUST be validated for precision on the demo project before it reaches any UI" (`:67-68`), and `related-data-navigation` that discovery "MUST stay an experiment validated for precision" (`:94`). Nothing does it.
- **A later phase needs a number to be held to.** Proposal D12's second bar compares a run that calls a model with "the run with no model". That run does not exist.

## Decided by the owner, and where

The owner's words are in the column "His words", and nowhere else in this section.

| Decision | What was put to him | His words | Where it is recorded |
|---|---|---|---|
| Proposal D12 | "Nothing the Advisor learns from a model reaches a user until the benchmark of section 12 has measured it against a bar you set. My suggested bars: no false line in any batch in any of the six databases; recall never below the run with no model; no invented concept where one existed; and nine in ten of a model's questions with the true answer among the buttons." (L3001) | "D12: Approve these bars" (2026-10-08, 21:18 UTC) | the proposal's card, L3004; his message in the session record of that day, which is in no repository |
| Proposal D7 | "… Whether a candidate is offered is decided by checks that can rule it out and pointers that name a target, not by a score." (L2961) | "D7: Approve" (the same message) | the proposal's card, L2964; the same record |
| This specification | A list of what was left after Phase 3, among it "Phases 4 to 8 … Each needs a specification you approve first." | "Do cleanup, specify phase 4 and if make sense others" (2026-10-10) | the session record of that day; it is in no repository |

What the approval of proposal D12 allows, in the proposal's words, which are its author's: "these bars as the gate for Phases 5 and 6, and a specification for the benchmark harness. Building the harness is part of Phase 5, after that specification is approved." (L3003)

The four bars, each with its number:

| # | His approved words | The number |
|---|---|---|
| 1 | "no false line in any batch in any of the six databases" | zero false lines; six databases |
| 2 | "recall never below the run with no model" | not below the no-model run's recall |
| 3 | "no invented concept where one existed" | zero |
| 4 | "nine in ten of a model's questions with the true answer among the buttons" | nine in ten |

What he did not decide, so that nobody reads it into the rows above:

- The card lists four bars. The table of the proposal's section 12 has more cells headed "Suggested bar": for meanings "Precision as for references" (L2542), for a candidate shown with no pointer "None, by rule" (L2544), for actions "Reported" (L2545), for a rung "A rung that adds no recall at equal precision is removed" (L2547). Its own note calls every bar "My assumption" (L2549). Whether "these bars" takes in those rows is not in his words. This specification treats the first two as its own checks and says so in "The author's choices".
- What is true. The proposal says the truth of benchmark A is "The engine's own declared foreign keys" (L2516). That sentence is its author's, and Question 1 shows where it and bar 1 cannot both hold.
- Whether the renamed copies are inside a pass mark. The proposal says "All three results are reported." (L2527) and no more.
- The scripted user's policy, what counts as an action, and the seven rules of "What keeps the result honest" (L2526 to L2532) are the proposal author's text.
- The number of actions. The proposal sets none: "Reported. Section 7 assumes about seven for Chinook's structure." (L2545), and marks that seven "My assumption" (L1597).

## What the rules give on the six databases

This is the author's calculation, not an implementation and not a run of the harness. A script applied the naming rule, the three gates and the level rule of `semantic-advisor`, and its concept rule, to the six pinned fixtures, read-only, with each fixture's foreign keys treated as hidden and with the original names. The script is not in any repository.

| Database | Hidden references | Value counts run | Batch lines | … that are hidden references | Questions | … with the true answer among the buttons | Reached by a question | Left for the scripted user to state |
|---|---|---|---|---|---|---|---|---|
| chinook | 11 | 9 | 9 | 9 | 0 | 0 | 0 | 2 |
| northwind | 13 | 11 | 9 | 9 | 0 | 0 | 0 | 4 |
| sakila | 22 | 27 | 18 | 18 | 4 | 3 | 3 | 1 |
| pubs | 10 | 13 | 7 | 7 | 3 | 3 | 3 | 0 |
| employees | 6 | 6 | 6 | 6 | 0 | 0 | 0 | 0 |
| adventureworks | 89, and one constraint of two columns | 197 | 74 | 70 | 15 | 14 | 14 | 5 |
| All six | 151 | 263 | 123 | 119 | 22 | 20 | 20 | 12 |

What stands out:

- In five of the six, every batch line is a hidden reference.
- **In AdventureWorks, four batch lines are not hidden references.** Two are real references that the database does not declare: `Production.TransactionHistoryArchive.ProductID` and `Production.WorkOrderRouting.ProductID`, both to `Production.Product`. Two are the columns of the one foreign key of two columns: `Sales.SalesOrderDetail.ProductID` and `.SpecialOfferID`. Question 1 is about the first pair; choice 6 about the second.
- The references left to state are the ones no name points at (`Customer.SupportRepId` and `Employee.ReportsTo` of Chinook, `Employees.ReportsTo` and `Orders.ShipVia` of Northwind, `Production.Document.Owner` and the two columns of `Production.BillOfMaterials` of AdventureWorks), the ones with no value to count (Northwind's `CustomerCustomerDemo` is empty; sakila's `film.original_language_id` is null in every row), and two columns of AdventureWorks with five tables to choose from.
- The one wrong question in sakila asks which of `film.film_id` and `film_text.film_id` refers to the other. Neither does. The one in AdventureWorks is about `Person.BusinessEntity.BusinessEntityID`.
- With opaque names nothing points: no batch line and no question in any of the six, and all 151 are stated.
- For meanings, on the original names: 17 table lines and 5 column lines in all six. Against the published meaning files, 18 are true, none is false, and four have no binding to compare with: the tables `city` and `country` of sakila, `Sales.Currency` of AdventureWorks, and `Employees.ReportsTo` of Northwind. Question 3 is about those.

The renamed variant was not calculated.

## The author's choices

Each is the choice of this specification's author. None is the owner's ruling. Each can be changed by one word from him.

| # | The choice | Why | The alternative |
|---|---|---|---|
| 1 | A feature of its own, beside `semantic-advisor`. | Proposal D12's card allows "a specification for the benchmark harness". It has no screen, its own inputs and outputs, and the proposal's Phase 6 runs it again with a model. | An appendix of `semantic-advisor`. |
| 2 | The six fixtures are the files of the `demo-db` repositories at the commits and checksums of REQ:fixtures-pinned. Chinook is demo-db/chinook's file, not the file that the web app's journey job and `open-layer-reading` use. | The proposal names "each a repository in the `demo-db` organisation" (L2522). The two Chinook files have the same 11 tables, 64 columns and 11 foreign keys and different bytes (SHA-256 `7651ba37…` against `f82efedb…`). | One Chinook for everything. |
| 3 | The copies are built on every run and are never committed or uploaded. | The proposal: "The data licence of each must be confirmed before a modified copy is redistributed." (L2522). The employees data is CC-BY-SA-3.0. A copy that is never stored is never redistributed. | Commit the 18 copies. |
| 4 | A copy holds tables, columns, primary keys and rows. It holds no view, index, trigger, check or unique constraint. A copy is pinned by a digest of its content, not of its bytes. | The 38 views of the six would have to be rewritten for each renaming, and no registered model has a record type for a view. Two versions of SQLite can write different bytes for the same content. | Keep the views in the original variant. |
| 5 | The renamed variant replaces words, not names: every word but `id` gets another word, drawn by the seed from a list that holds no word of a searched concept and no word of the database. A plural keeps its `s`. | A hand-made list of 128 table names and 816 column names would be the largest artefact of the harness and a place for mistakes. Word for word, `Customer` and `CustomerId` still agree, which is what the proposal asks: "the naming signal survives and the memorised schema does not" (L2527). Its own example, `Buyer`, is a synonym of the universal concept `customer`; a list that allows it would let the concept rule find the answer. | A hand-written list of plausible names. |
| 6 | A batch line on a column of a hidden foreign key of several columns is counted apart: it is neither true nor false. | The proposal: that constraint "is scored separately, because the format cannot yet state it" (L2516). | Count it false. |
| 7 | The harness removes all six databases' entries from both indexes in every run, not only the entry of the database under test. | It is simpler than deciding which entries were "made from the same fixture" (L2526), and nothing in a no-model run reads another database's entry. | Remove one entry per run. |
| 8 | The universal concepts are given at `meaninggraph/core@dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7`, from a folder the harness fills before the run. During a run nothing is requested from a network. | Benchmark B: "The universal concepts are given" (L2517). That commit is what the MeaningGraph registry lists for `core` (meaninggraph/registry@49f472f). A run with no network is repeatable. | Fetch during the run. |
| 9 | An action is a call of `accept`, `reject`, `answer`, `state` or `keep`. `scan`, `suggest`, `review` and `follow` are not actions. | The proposal lists four kinds (L2528) and elsewhere also counts connecting and approving the diff (L1597). A scan is the same single step in every run; a Keep is the person's decision and is counted. | Count the scan and the review too. |
| 10 | Two more things are checked in every run and fail it: a candidate shown with no pointer, and a batch of meanings with a false line. | They are rows of the proposal's section 12 (L2544, L2542). The first is a rule of `semantic-advisor`; the second follows from the words of bar 1, "any batch". Neither is on the card. Question 3 asks about the second. | Report both and fail on neither. |
| 11 | Bar 2 is applied to this run as a committed baseline: the no-model run on a later commit must reach every hidden reference that the baseline lists, at the same level or a better one, unless the same change edits the baseline. | As worded, the bar compares a later run with this one and cannot fail here. A baseline makes a rule change that loses a reference visible in the pull request that causes it. | Only record the recall. |
| 12 | The scripted user is a program that calls the commands of `datatug advisor` and nothing else. | "Agents and people use the same verbs." (L2504). An action is then something a person could do. | Call the Advisor's functions directly. |
| 13 | SQLite only. | The six fixtures are SQLite files. The Advisor's rules read the shape file and two counts, which do not depend on the engine. | Load each into PostgreSQL as well. |
| 14 | The job runs on a pull request and on a push to the main branch of `datatug-cli`, when the change touches the Advisor or the harness. It has no schedule. | The measurement belongs with the change that moves it. | A nightly run. |

## Questions for the owner

Four questions. Each changes what is built or what "passed" means.

### Question 1. Is a real reference that the database does not declare a false line?

AdventureWorks declares a foreign key on `Production.TransactionHistory.ProductID`, to `Production.Product`. Its archive table has the same column and no foreign key. With the rules of `semantic-advisor` the batch for AdventureWorks holds this line:

```
Production.TransactionHistoryArchive.ProductID refers to Production.Product
    ✓ named after Product's key ProductID
    ✓ Production.Product.ProductID is the primary key
    ✓ all 497 distinct values occur there
```

`Production.WorkOrderRouting.ProductID` is a second one, with 149 of 149. By the proposal's sentence the truth is "The engine's own declared foreign keys" (L2516), so both lines are false, and the pass mark you approved, "no false line in any batch in any of the six databases", is not met. No fixed rule can tell these two from the 70 lines beside them that are declared.

- **A. No: a person's verdict settles it (recommended).** The truth is the declared foreign keys and, for each batch line they do not contain, one verdict of a person, written into a file of the harness with a reason and frozen. A line with the verdict "true" is true. A line with no verdict fails the run, so nothing is waved through. The proposal does the same for meanings: "Truth is reviewed before it is used." (L2529). By the author's calculation that is two verdicts for the original names.
- B. Yes, by the letter. Then the rules of `semantic-advisor` do not pass on AdventureWorks, and Phase 5 is not done until either a rule is found that leaves those two out and keeps the 70, or the bar is changed.
- C. AdventureWorks is left out of pass mark 1 and only reported. Then the bar holds for five databases, and 91 of the 153 references are outside it.

### Question 2. Are the renamed copies inside the pass marks?

The proposal gives a model's memory as the reason for them: "A generative model has very likely seen Chinook's schema in training" (L2527). This run calls no model. The proposal still lists "both renamed variants" for Phase 5 (L2808) and says "All three results are reported." With opaque names the run looks like this for every database:

```
chinook · opaque names     batch 0 lines     stated by the scripted user 11
```

- **A. Yes, all three (recommended).** Pass marks 1 and 3 hold in all 18 runs, and each of the 18 reaches a complete model. It costs nothing, since no rule should make a false line under any name, and it catches a rule that leans on a spelling. The number of actions is reported for each run and has no limit.
- B. Only the original names are inside the pass marks. The other 12 runs are reported.

### Question 3. Who says what a column means when the published file is silent, and is a false meaning a false line?

The withheld meaning files are thin. Chinook's binds 25 of 64 fields; sakila's binds none. The Advisor's batch for sakila holds:

```
city is a City          ✓ the table's name is the label of the universal concept city
```

Sakila's published meaning file binds no concept to the table `city`. So the harness cannot say whether the line is true. The proposal asks that "a person marks every remaining column as either bindable to a named concept or deliberately unbound" before any run (L2529). For six databases that is 762 of 816 columns.

- **A. A person gives a verdict only for what a batch names, and a false meaning fails the run (recommended).** The harness starts from the published bindings. A batch line about a table or a column they do not bind needs one verdict in a file of the harness: a concept, or "unbound". Without it the run fails and names the line. By the author's calculation that is four verdicts for the original names. You, or someone you name, give them. Pass mark 1 then covers batches of meanings too: its words are "any batch".
- B. The full review the proposal asks for, before the first run: 762 columns.
- C. Meanings are reported and are outside pass mark 1 until a full review exists.

### Question 4. Are the questions that fixed rules ask held to nine in ten?

Bar 4 is worded for "a model's questions". Fixed rules ask questions too, when two tables fit one column. In pubs:

```
titles.pub_id: which one?     [publishers]  [pub_info]  [None of these]
```

The true answer, `publishers`, is among the buttons. By the author's calculation the six databases get 22 such questions, 20 of them with the true answer among the buttons: nine in ten, and no more.

- **A. Reported, not a pass mark in this run (recommended).** The card's words are about a model. The share is in every report, and the proposal's Phase 6 holds a model's questions to it.
- B. Held to nine in ten now, over the six databases together. One more wrong question would then fail the run.

## Behavior

Requirements use MUST, MUST NOT, SHOULD and MAY as fixed terms.

### Fixtures and copies

#### REQ: fixtures-pinned

The harness MUST use exactly these six files, each read from its repository at the commit given, and MUST stop before any run when a file's SHA-256 is not the one given.

| Database | Repository and commit | Path | SHA-256 of the SQLite file |
|---|---|---|---|
| chinook | demo-db/chinook@88f080d72205034e1161bb80594acb7d77fc58ef | `data-source/source.sqlite` | `7651ba378ac2fcd0dfc3c66fb101f7a7eed3ba39a612ec642b96e20702061f15` |
| northwind | demo-db/northwind@c07b4666734a3445b2b0c2fbaf84305dcdb1bf45 | `data-source/source.sqlite` | `279b34136771aee75d802094b3329515a2b01da65d2a20a4a9e3b58c29b4fd20` |
| sakila | demo-db/sakila@113e54ad83c3e003a4fd195c893f2b19a921435c | `data-source/source.sqlite` | `9bebeee50fecb1fee115c206a4f380f2f6d1f009e7796b18b8a8cfba7c810454` |
| pubs | demo-db/pubs@e99ca33330a043e23b2c5a665356fb8ad8d0b508 | `data-source/source.sqlite` | `b45b08b7c06441cc0b138b031ea5bd32bd5866aebf874c89d2ebf7fc6abb3691` |
| employees | demo-db/employees@add7744ddc5cefc535a7cbb09b3dffee41811ade | `data-source/source.sqlite.gz` | `46b49dd0e141cd8db66d57680a10febfccc15e6f54dc7b8f3c2316b55e11c5f0`, unpacked |
| adventureworks | demo-db/adventureworks@c91e06e9c12380e7709f66ab1e69cf7f858f67c1 | `data-source/source.sqlite.gz` | `6a105e1982becfe003fc7a307d7cad9d738390cedd79d166c2167fd817168109`, unpacked |

Each database's published meaning file is `model/<id>.meaning.yaml` at the same commit. A fixture is opened read-only and is never changed.

#### REQ: stripped-copy

For each fixture the harness MUST build a stripped copy with a script, never by hand: a new SQLite file with every table of the fixture, its columns in their order with the same type text, the same nullability and the same default, its primary key, and every row. The copy MUST hold no foreign key, and no view, index, trigger, check or unique constraint. The harness MUST compute a digest of each copy's content, tables in the order of their names and rows in the order of their keys, and MUST stop when it is not the digest committed with the harness for that fixture, variant and seed.

#### REQ: renamed-variants

Each stripped copy MUST be built in three variants. `original`: the fixture's names. `renamed`: every word of every table and column name, split as `semantic-advisor` splits a name, is replaced by another word, except the word `id`; the same word is always replaced by the same word; a word and its singular share one replacement; the letter case and the separators of the name are kept. The words are drawn by a seeded generator from a list committed with the harness, from which every word of a searched concept and every word of the database is removed first. `opaque`: tables are `t1`, `t2`, … in the order of their original names, and the columns of each are `c1`, `c2`, … in their order. The seed is one number, committed with the harness. The same seed MUST give the same copies.

#### REQ: copies-not-committed

A stripped copy MUST exist only in the working folder of a run. The harness MUST NOT commit one, attach one to a job or upload one. The repository of the harness holds the script, the seed, the word list, the digests, the truth files and the baseline.

#### REQ: masked-registry

Every run MUST set the two settings of `open-layer-reading` REQ:registry-index-source to masked indexes: the index files of `modelspec-org/registry` and of `meaninggraph/registry`, each at a commit committed with the harness, with the items of the six databases removed and the index's checksum computed again. The entry of `core` stays. Every run MUST start with a cache folder that holds the files of `core` at `dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7` and nothing else, an empty Advisor ledger and an empty project folder.

#### REQ: no-model-no-network

A run MUST call no AI model and no decision engine, and MUST make no request to any network address from its scan to its report. The harness MUST fail a run in which one is made. Fetching the fixtures, the indexes and `core` happens before the first run.

### Truth

#### REQ: truth-references

The truth of benchmark A for a database is: every foreign key of one column that the fixture declares, read from the fixture before it is copied, as a `refers-to` from that column to its target table; and the verdicts of `truth/<id>.references.yaml` (Question 1). A batch line is true when it is a hidden reference or has the verdict `true`; false when it has the verdict `false`; counted apart when its column belongs to a foreign key of several columns; and otherwise unreviewed. In the renamed and opaque variants the truth is the same statements under the copy's names.

#### REQ: truth-meanings

The truth of benchmark B for a database is: every binding of its published meaning file, as the table or column it binds and the concept, with every universal concept that concept extends; and the verdicts of `truth/<id>.meanings.yaml` (Question 3). A meaning line is true when its concept is the concept the truth has there or a universal concept that one extends; false when the truth has another concept there or the verdict `unbound`; and otherwise unreviewed. A table or a column that the truth binds twice is matched by either.

#### REQ: verdicts-are-a-persons

A verdict MUST carry who gave it, when and one line of why, and is changed only by a person in a commit. The harness MUST NOT write a verdict. A run that meets an unreviewed batch line MUST fail and name the database, the variant and the line.

### The scripted user

#### REQ: scripted-user

The scripted user MUST act only through the commands of `datatug advisor`, by this policy, in this order: (1) `suggest` for the whole database; (2) for each batch, `reject` each false line, then `accept` the batch; (3) for each question, `answer` with the true alternative, or with `--none` when none is true; (4) repeat from 1 until no card is left; (5) `state refers-to` for each hidden reference that is not yet a fact; (6) `state is-a` for each table the truth binds that has no table concept yet, with the universal concept when the truth's concept extends one and with `--new` and the truth's own label and description when it does not; (7) `state means` for each column the truth binds that has no matching link yet, in the same way; (8) `review`, then `keep`. It MUST read the truth and MUST NOT read the fixture, the published files or the Advisor ledger.

#### REQ: action-count

Each call of `accept`, `reject`, `answer`, `state` or `keep` is one action, whatever the number of lines it covers. The harness MUST count the actions of every run by kind and in total. `scan`, `suggest`, `review` and `follow` are not counted.

#### REQ: complete-model

A run has reached a complete model when, in the files that the Keep wrote: every table of the copy has a record type whose key is the table's primary key; every hidden reference is a reference field to its target's record type; every table that the truth binds has a concept bound with the role `instances` that matches the truth; every column that the truth binds has a written or a derived link to a concept that matches the truth; and ModelSpec's lint, MeaningGraph's check and the drift report of `open-layer-reading` report no error. A column of a foreign key of several columns is not required to be a reference.

### What is measured

#### REQ: false-line

Before the scripted user changes anything in a batch, the harness MUST compare every line of that batch with the truth and record it as true, false, counted apart or unreviewed. This is done for every batch that any `suggest` of the run returns, of references and of meanings.

#### REQ: measures-reported

For every run the report MUST give: the hidden references; the batch lines of each kind and how many are true, false and counted apart; the questions asked, how many have the true answer among their buttons, and how many are about a pair that is not a reference; the hidden references reached at each of the levels batch, question and stated; of the tables and columns the truth binds, how many got a suggestion; the concepts with no `extends` that a card proposed; the candidates shown with no pointer; the actions by kind and in total; the number of value counts and the time each table took, as a median and a 95th percentile; and whether the model is complete. The report MUST also give each measure for the six databases together, per variant.

#### REQ: pass-marks

The harness MUST fail when one of these does not hold, and MUST name the run and the line:

1. "no false line in any batch in any of the six databases": no batch of any of the 18 runs holds a false line or an unreviewed line (Questions 1, 2 and 3).
2. "recall never below the run with no model": every hidden reference that the committed baseline lists as reached is reached, at the same level or a better one (REQ:baseline).
3. "no invented concept where one existed": no card of any run proposes a concept that is not one of those searched, and no kept meaning file holds a concept with no `extends` for a table or a column whose truth extends a universal concept.
4. "nine in ten of a model's questions with the true answer among the buttons": no model asks a question in this run. The share for the questions that fixed rules ask is reported and fails nothing (Question 4).

It MUST also fail when a card shows a candidate that no pointer names, and when any run does not reach a complete model.

#### REQ: baseline

The harness MUST keep one committed file that lists, for each database and variant, every hidden reference and the level at which the no-model run reached it: batch, question or stated. A run compares itself with that file. A change that lowers a level MUST fail unless it edits the file in the same pull request, so that the loss is seen and agreed. This file is the "run with no model" that the proposal's Phase 6 is held to.

#### REQ: run-record

A report MUST record the commit of `datatug-cli`, the versions of the two checkers, the commits of the six fixtures, of the two indexes and of `core`, the seed, the digests of the 18 copies, and the full list of candidates of every `suggest` with their levels and reasons. Two runs on the same commits MUST give the same report, apart from times.

### Where it runs

#### REQ: ci-job

The harness MUST run as one job, `advisor-benchmark`, in the CI of `datatug-cli`, on every pull request and every push to its main branch that changes a file of the Advisor or of the harness. It MUST write the report as a file of the job and print the lines of "What a person sees". It has no schedule. A developer MUST be able to run the same harness on their machine with one command and get the same report.

## What must not change

| What | How to test that it did not |
|---|---|
| The six fixtures and their repositories | Each file's SHA-256 is checked before and after the 18 runs; the harness clones or downloads and never pushes (AC:fixtures-are-verified). |
| The registries | The indexes are read at pinned commits before the run; during a run no request is made at all (AC:run-is-offline). |
| That nothing leaves the machine, and no model is called | A run with any outgoing request fails (AC:run-is-offline). |
| That no modified copy of a database is redistributed | After a run, the repository and the job's files hold no SQLite file (AC:copies-are-not-kept). |
| The launch lane's work | The job is new, runs only for changes to the Advisor or the harness, and changes no existing job, test or fixture of `datatug-cli`, `datatug-apps` or `datatug/datatug-demo-project` (AC:job-runs-on-advisor-changes). |
| What `semantic-advisor` and `open-layer-reading` specify | The harness adds no rule, no setting and no command to either; the scripted user uses the commands as a person would (AC:scripted-chinook). |

## Existing features

Line numbers of this repository's own pages are as at datatug/datatug@4aa2b66.

| Feature | Status | Relation | What follows |
|---|---|---|---|
| `semantic-advisor` | Draft | Depends on | It is what is measured. Its REQ:done-when-is-the-benchmark points here. This feature cannot be built before its commands exist. |
| `open-layer-reading` | Draft | Depends on | Its two index settings, its cache folder and its drift report are used as they are. |
| `meaning-inference` | Draft | Extends | Its condition "validated for precision on the demo project before it reaches any UI" (`:67-68`) is met here, on six databases and not one. Listed for its owner. |
| `related-data-navigation` | Draft | Extends | The same, for its REQ:automatic-relationship-discovery (`:94`). |
| `core-investigation-loop` | Stable | Leaves | Its REQ:journey-e2e is the house precedent for a walk "against the real agent … with no network interception". This harness drives commands, not a browser. |
| `checks`, `compare` | Approved | Leaves | A pass mark is not a `Check`, and comparing a kept model with the truth is not a Compare. |
| `public-data-fabric` | Draft | Leaves | Another lane. Its workflow's job `ledger` stays the one required check of this repository. Nothing under it is read or edited. |

## Dependencies

- semantic-advisor
- open-layer-reading

## Blocked Clauses

- **Everything here** needs the commands of `semantic-advisor`, which is a Draft and unbuilt.
- **REQ:complete-model** runs MeaningGraph's check on files in the format `semantic-advisor` writes. Under its Question 3, answer A, that needs a released checker for `meaning/draft-2`.
- **REQ:truth-references and REQ:truth-meanings** need the verdicts of Questions 1 and 3 before pass mark 1 can be met: by the author's calculation two and four, for the original names. The renamed variant may need the same ones under other names; the harness translates a verdict as it translates the truth.

## Acceptance Criteria

"The harness" is the program at the commit under test, run with the committed seed. "A run" is as in "Words used here".

### AC: fixtures-are-verified (verifies REQ:fixtures-pinned)

**Given** the six fixtures fetched, and one more fetch in which one byte of `chinook`'s file is changed
**When** the harness starts on each
**Then** on the first it prints the six SHA-256 values of REQ:fixtures-pinned and goes on; on the second it stops before any scan and names `chinook` and the two values; and after a full run the six files have the same SHA-256 as before.

### AC: stripped-chinook (verifies REQ:stripped-copy)

**Given** the fixture of Chinook
**When** the harness builds its stripped copy in the variant `original`, twice
**Then** each copy has 11 tables, 64 columns, 11 primary keys, no foreign key and no view; `Invoice` has 412 rows and `Customer` 59; `PRAGMA foreign_key_list` returns nothing for any table; both builds have the committed digest; and a copy in which one row of `Invoice` is changed by the test is refused for its digest.

### AC: variants-are-built (verifies REQ:renamed-variants)

**Given** the fixtures of Chinook and of Northwind
**When** the harness builds the variants `renamed` and `opaque` of each, twice with the committed seed and once with another
**Then** in `renamed`, the table that was `Customer` and the column that was `CustomerId` share their first word, the column still ends with `Id`, no name is a name of the fixture, and no word of a name is a label or a synonym of a concept of `core`; the table that was `Customers` in Northwind ends with `s`, and its name without that `s` is the first word of the column that was `CustomerID`; in `opaque`, Chinook's tables are `t1` to `t11` and the table that was `Invoice` is `t6` with the columns `c1` to `c9`; the two builds with one seed have equal digests; and the build with another seed has other names.

### AC: copies-are-not-kept (verifies REQ:copies-not-committed)

**Given** a clone of the harness's repository and one finished job
**When** the repository's files, the job's files and the job's log are searched
**Then** no SQLite file and no dump of one is found; the repository holds the script, the seed, the word list, 18 digests, the truth files and the baseline.

### AC: registry-is-masked (verifies REQ:masked-registry)

**Given** the stripped Chinook in the variant `original`, scanned, once with the masked indexes and once with the same index files unmasked
**When** `datatug db match` and `datatug advisor suggest --format json` run on each
**Then** with the masked indexes no registered model matches, and no reason of any card names a registered model; with the unmasked ones the model `chinook` matches and nine lines name it, which shows what the mask hides; and each masked index lists `core` and none of `chinook`, `northwind`, `sakila`, `pubs`, `employees` and `adventureworks`, with a checksum that fits its content.

### AC: run-is-offline (verifies REQ:no-model-no-network)

**Given** the fixtures, the indexes and `core` fetched, and every outgoing request of the harness recorded from then on
**When** the 18 runs are done; and, in a second pass, a test build of the Advisor that makes one request during `suggest` is used
**Then** the first pass records no request at all; the second fails its run and names the address.

### AC: truth-of-chinook (verifies REQ:truth-references, REQ:truth-meanings)

**Given** the fixture of Chinook and its published meaning file at the pinned commit
**When** the harness builds the truth for the variants `original` and `opaque`
**Then** the references are 11, among them `Invoice.CustomerId → Customer`, `Customer.SupportRepId → Employee` and `Employee.ReportsTo → Employee`, and in `opaque` the same 11 under names such as `t6.c2 → t3`; the meanings hold a concept for 10 tables and for 25 columns, among them `Invoice.Total` with `invoice-total` and the universal concept it extends, and `Track.UnitPrice` with `list-price` and the universal concept `price`.

### AC: unreviewed-line-fails (verifies REQ:verdicts-are-a-persons, REQ:false-line, REQ:truth-references)

**Given** the harness with an empty `truth/adventureworks.references.yaml`, and then with two verdicts `true` in it, for `Production.TransactionHistoryArchive.ProductID` and `Production.WorkOrderRouting.ProductID`, each with a name, a date and a reason
**When** the run of AdventureWorks in the variant `original` is done each time
**Then** the first fails pass mark 1 and names both lines as unreviewed, before the scripted user has acted on the batch; the second records 72 lines as true and two as counted apart, `Sales.SalesOrderDetail.ProductID` and `Sales.SalesOrderDetail.SpecialOfferID`; and the harness has written no line to either truth file.

### AC: false-line-fails (verifies REQ:false-line, REQ:pass-marks)

**Given** the harness, and a copy of its truth in which the test has removed `Invoice.CustomerId → Customer` from Chinook's hidden references and added the verdict `false` for it
**When** the run of Chinook in the variant `original` is done
**Then** it fails pass mark 1 and names `Invoice.CustomerId refers to Customer` in the references batch of `chinook · original`; the report records that the scripted user rejected that line before accepting the other eight; and the job's result is a failure.

### AC: scripted-chinook (verifies REQ:scripted-user, REQ:action-count, REQ:complete-model)

**Given** the stripped Chinook in the variant `original`, with the truth of AC:truth-of-chinook
**When** the run is done
**Then** the recorded commands begin with `suggest`; the actions about references are exactly three, one `accept --batch references` of nine lines and two `state refers-to`, for `Customer.SupportRepId` and `Employee.ReportsTo`; one `accept --batch meanings` covers the eight lines of `semantic-advisor` AC:chinook-meanings; the last action is one `keep`; no command other than those of `datatug advisor` and `datatug scan` was called; the report gives a count for each kind of action and their total; and the kept files hold 11 record types, 11 reference fields, and a matching concept for each of the 10 tables and 25 columns of the truth, with no error from ModelSpec's lint, MeaningGraph's check or the drift report.

### AC: scripted-opaque (verifies REQ:scripted-user, REQ:complete-model, REQ:measures-reported)

**Given** the stripped Chinook in the variant `opaque`
**When** the run is done
**Then** no `suggest` returns a card about a reference or a meaning; the scripted user states all 11 references, one action each; the report says batch 0, question 0, stated 11, and no value count for a card; and the model is complete.

### AC: report-has-every-measure (verifies REQ:measures-reported, REQ:pass-marks, REQ:run-record)

**Given** two full passes of the harness on the same commits
**When** the two reports are compared, and the first is read
**Then** they are equal apart from times; the report quotes the four bars in the owner's approved words, each with its value and `PASS`, `FAIL` or `reported`; for sakila in the variant `original` it lists the question about `film.film_id` and `film_text.film_id` as one about a pair that is not a reference; it gives every measure of REQ:measures-reported for each of the 18 runs and for the six databases together per variant; and it records the commits, the versions, the seed, the 18 digests and every candidate with its level and reasons.

### AC: baseline-holds-recall (verifies REQ:baseline, REQ:pass-marks)

**Given** the committed baseline, and a test build of the Advisor whose naming rule no longer fires for `InvoiceLine.TrackId`
**When** the harness runs with that build, and again after the baseline's line for `InvoiceLine.TrackId` is changed from `batch` to `stated` in the same change
**Then** the first fails pass mark 2 and names `chinook · original · InvoiceLine.TrackId: batch in the baseline, stated in this run`; the second passes, and the change to the baseline is a line of the pull request's difference.

### AC: invented-concept-fails (verifies REQ:pass-marks)

**Given** a test build of the Advisor in which a card proposes a new concept `invoice` with no `extends` for the table `Invoice`
**When** the run of Chinook in the variant `original` is done
**Then** it fails pass mark 3 and names the card; and with the real build the report says `0` for concepts with no `extends` that a card proposed, in all 18 runs.

### AC: job-runs-on-advisor-changes (verifies REQ:ci-job)

**Given** the CI of `datatug-cli` with this feature, and three pull requests: one that changes a naming rule of the Advisor, one that changes a file of the harness, and one that changes only the terminal chat
**When** each is opened
**Then** the job `advisor-benchmark` runs on the first two and not on the third; its files hold the report; its log holds the lines of "What a person sees"; no existing job of that repository has changed its name, its steps or its result; and the same command on a developer's machine gives a report equal to the job's, apart from times.

### AC: done-when (verifies REQ:pass-marks, REQ:complete-model, REQ:action-count, REQ:masked-registry, REQ:renamed-variants)

**Given** the six fixtures of REQ:fixtures-pinned, the masked indexes, the committed seed, the truth with its verdicts, and builds of `datatug-cli` that contain `semantic-advisor` and this feature
**When** the harness does its 18 runs: each of the six databases, without its foreign keys, in the variants `original`, `renamed` and `opaque`
**Then** every run reaches a complete model; the report gives the number of scripted actions of each run; no batch of any run, of references or of meanings, holds a false line or an unreviewed one; no card of any run proposes a concept that is not one of those searched; no request left the machine and no model was called; pass marks 1, 2 and 3 say `PASS` and pass mark 4 says `reported`; and for Chinook in the variant `original` the references batch has nine lines, all true.

## Rehearse Integration

Every criterion has a surface that a test can observe: a SHA-256, a digest, a line of the report, the list of commands the scripted user called, the job's result. The harness is itself a test program, and each criterion is exercised against the real commands of `datatug-cli`. Stub scaffolding is therefore deferred to that implementation, as `open-layer-reading` does.

## Not Doing / Out of Scope

- The proposal's benchmark D and its tests C and E. They need code and two databases, which are its Phase 7.
- A run with Jev or with a generative model, running anything three times, and the measure "what each rung adds". They are its Phase 6, which reuses this harness and this baseline.
- A limit on the number of actions. It is counted and reported.
- More databases than six. The proposal notes that six "can demonstrate a batch precision of about 0.975 at most" (L3002) and that 91 of the 153 references sit in one schema (L2530); a larger set is one of its open questions (L3122).
- PostgreSQL and SQL Server.
- A page, a dashboard or a badge for the results.
- Any change to the six `demo-db` repositories, to a registry or to `meaninggraph/core`.

## Open Questions

The questions that change what is built are in [Questions for the owner](#questions-for-the-owner). The points below are left open by this specification's author.

For later phases:

1. The same harness with a model: which rung adds what, and three runs each (the proposal, L2532, L2547).
2. A full review of every column of the six databases for benchmark B, if Question 3 is answered A.
3. A larger set of hard cases than six databases hold.
4. A foreign key of several columns as a hidden reference, once a model can state one.
5. The same run on PostgreSQL.

Not about what is built:

6. The data licences. The harness redistributes no copy. Whether the truth files, which name tables and columns of the six databases and no row, need a licence note of their own was not checked.
7. When building may start. The author's suggestion is to build this before the Advisor's pages and after its commands, so that no card is seen before a run has passed.

---
*This document follows the https://specscore.md/feature-specification*
