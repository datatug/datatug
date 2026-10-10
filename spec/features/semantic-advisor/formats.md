# Semantic Advisor: rules, files and wire formats

This is the formats appendix of the [owning Feature](README.md). It is not a separate
feature. It holds the exact shapes that the requirements there point at: the Advisor
ledger, the statements and their ids, the naming rule, the concept rule, the gates and the
levels, the answers file, the three files a Keep writes, the review, six routes and the
commands.

Nothing here is built. Everything here is the specification author's proposal unless a
line says it is the owner's. The words are those of the Feature's section
[Words used here](README.md#words-used-here) and of
[`open-layer-reading`](../open-layer-reading/README.md#words-used-here).

Every example uses "the stripped Chinook": the copy of Chinook with no foreign key that
[`advisor-benchmark`](../advisor-benchmark/README.md) builds from
demo-db/chinook@88f080d `data-source/source.sqlite`, scanned as the database `chinook` in
the environment `local`. Universal concepts are read at
`meaninggraph/core@dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7`, written `dd5ce32` below,
which is the commit the MeaningGraph registry lists for `core`
(meaninggraph/registry@49f472f `index.json`).

## 1. The Advisor ledger

| | |
|---|---|
| Folder | `<personal folder>/<project id>/advisor/<database id>/`. The personal folder is the one personal queries use: the setting `DATATUG_PERSONAL_DIR`, or else `~/.datatug/projects` (datatug-cli@3a1c276 `pkg/personalqueries/dir.go:19-25`). Question 1 of the Feature. |
| Files | `signals.jsonl` and `follows.jsonl`: one JSON object per line |
| Written by | the agent and the `datatug advisor` commands, never by a person |
| Written how | a line is appended. No line is ever changed or removed. |
| Holds | names, counts, origins and times. Never a value of a row. |
| Committed | never. It is outside the project. |
| If deleted | every signal is made again by the next request. The lines of `follows.jsonl` are lost. No answer is lost. |

### A signal

One thing one rule found. Two real lines:

```json
{"id":"sig_c0d9450fb0f82f32","kind":"metadata","job":"pointer","rule":"same-name-as-key","strength":"strong","statement":"refers-to","from":"chinook/Invoice.CustomerId","to":"chinook/Customer","origin":"environments/local/catalogs/chinook/chinook.shape.json","at":"2026-10-12T09:00:04Z"}
{"id":"sig_aab0d8bb5309b10e","kind":"statistics","job":"gate","rule":"values-present","result":"passed","statement":"refers-to","from":"chinook/Invoice.CustomerId","to":"chinook/Customer","counts":{"distinct":59,"present":59},"origin":"env:local","at":"2026-10-12T09:00:05Z"}
```

| Key | Value |
|---|---|
| `id` | `sig_` and the first 16 hex characters of the SHA-256 of the text `<rule>\n<statement>\n<from>\n<to>\n<origin>` |
| `kind` | `metadata`, `statistics` or `published`: three of the proposal's seven kinds (L2249). The other four (`queries`, `code`, `ai`, `person`) make no signal in this feature. |
| `job` | `pointer` or `gate` |
| `rule` | pointers: `declared-foreign-key`, `same-name-as-key`, `table-name-then-key`, `concept-names-table`, `concept-names-column`, `registered-model-states`. Gates: `types-comparable`, `target-is-key`, `values-present`, `concept-free`. |
| `strength` | pointers only: `decides`, `strong` or `weak` |
| `result` | gates only: `passed`, `partial`, `failed` or `not-counted` |
| `why` | with `not-counted` only: `no-values`, `time-limit` or `unreachable` |
| `counts` | `values-present` only: `distinct` and `present` (section 5) |
| `statement`, `from`, `to` | the statement the signal is about (section 2) |
| `origin` | where it was seen: the shape file's path for a rule that reads the shape; `env:<env>` for a count; the graph's address with its pin for a concept; the model's address with its commit for a registered model |
| `gone` | `true` on a line that says the rule no longer finds the signal |
| `at` | the time the line was written, UTC |

A line is appended when a signal is new, or when its `strength`, `result`, `counts` or
`gone` differ from the last line with that `id`. The current state of a signal is its last
line. So a second scan of an unchanged database appends nothing.

### A follow

One line per Follow. It names the link and never the value that was followed:

```json
{"candidate":"cand_d5bc057fcdd3654d","statement":"refers-to","from":"chinook/Customer.SupportRepId","to":"chinook/Employee","env":"local","by":"tom","rows":1,"at":"2026-10-13T16:18:40Z"}
```

No rule reads `follows.jsonl`. A Follow is never a signal.

## 2. Statements, candidates and their ids

A statement is one thing that can be said about a table or a column. There are four.

| Statement | `from` | `to` | A person reads | Where a Keep writes it |
|---|---|---|---|---|
| `refers-to` | a column, `chinook/Invoice.CustomerId` | a table, `chinook/Customer` | "CustomerId refers to Customer" | the model: the field becomes a reference field |
| `is-a` | a table, `chinook/Invoice` | a concept | "Invoice is an Invoice" | the meaning file: a concept bound to the record type with the role `instances` |
| `means` | a column, `chinook/Invoice.Total` | a concept | "Total is an Invoice total" | the meaning file: a concept bound to the field |
| `in-model` | a table, `chinook/Track` | empty | "Track is in the model" | the model: a record type; the source file: a line |

A column is written `<database id>/<Table>.<Column>` and a table `<database id>/<Table>`,
with the names the shape file has. A concept is written by its address without a pin,
`meaning://github.com/meaninggraph/core/invoice-total`, or `local:<concept id>` for a
concept of the database's own meaning file.

The proposal's figure names four statement words, `refers-to`, `means`, `maps-to` and
`same-values-as` (L2224). `maps-to` and `same-values-as` are not used here; `is-a` and
`in-model` are the author's.

A candidate is a statement that at least one signal is about and that no person has
answered. It is computed from the signals, the shape file and the answers on every
request, and is stored nowhere.

| Key | Value |
|---|---|
| `id` | `cand_` and the first 16 hex characters of the SHA-256 of `<statement>\n<from>\n<to>`. `refers-to\nchinook/Invoice.CustomerId\nchinook/Customer` gives `cand_85b83a04f850461e` on every machine. |
| `statement`, `from`, `to` | as above |
| `level` | `recorded`, `dropped`, `held`, `batch` or `question` (section 5) |
| `label` | `inferred` for `batch`, `suggested` for `question`, none for `dropped` and `held` |
| `reasons` | in the order pointers, gates: one for the naming rules together; one for each other pointer; one for each gate, `types-comparable` only when it fails; and for a `refers-to` a last one that says whether the database declares a foreign key. Each has `text`, `kind`, `job` and `mark`: `points`, `allows`, `in-part`, `rules-out` or `says-nothing` |
| `group` | `references` or `meanings` for a batch line; the question's id for an alternative |

## 3. The naming rule

**A name's words.** Every character that is not a letter or a digit becomes a break. A
break is added before an upper-case letter that follows a lower-case letter or a digit,
and before the last upper-case letter of a run that a lower-case letter follows. Digits
stay with the letters before them. Words are compared in lower case.

| Name | Words |
|---|---|
| `CustomerId` | customer, id |
| `customer_id` | customer, id |
| `InvoiceLineId` | invoice, line, id |
| `BusinessEntityID` | business, entity, id |
| `c2`, `Address2` | c2; address2 |
| `Sales.SalesOrderDetail` (a table) | sales, order, detail: a table's words are those of its name after the last dot |

**A table's singular form.** The table's words with the last word changed by the first
rule that fits: `ies` at the end becomes `y`; `es` after `s`, `x`, `z`, `ch` or `sh` is
removed; a single `s` at the end is removed. `Customers` gives customer; `Categories`
gives category; `Address` stays.

**The two rules.** Let A be a column, T a table with a primary key of exactly one column
B, and A not B itself. T may be A's own table.

| Rule | Fires when | Strength |
|---|---|---|
| `same-name-as-key` | A's words equal B's words, and B has two words or more | strong |
| `table-name-then-key` | A's words end with T's words, or T's singular form, followed by `id`; or followed by B's words, without T's words at their start when B begins with them | strong |

Found on the stripped Chinook, and nothing else is: nine pairs of a column and a table,
each by both rules. `Album.ArtistId → Artist`, `Invoice.CustomerId → Customer`,
`InvoiceLine.InvoiceId → Invoice`, `InvoiceLine.TrackId → Track`,
`PlaylistTrack.PlaylistId → Playlist`, `PlaylistTrack.TrackId → Track`,
`Track.AlbumId → Album`, `Track.GenreId → Genre`, `Track.MediaTypeId → MediaType`. They
are the nine that the proposal's probe ends with (L1550).

Must not be found:

| Column | Table | Why not |
|---|---|---|
| `InvoiceLine.InvoiceLineId` | `Invoice` | Its words end with line, id. Neither rule fires. The proposal's probe matched it by "contains the table's name" (L3142) and needed the value count to remove it. |
| `Customer.SupportRepId`, `Employee.ReportsTo` | `Employee` | No word of the table or its key is in the name. |
| a column `id` | a table whose key is `id` | The key has one word. Without this, every table of a database whose keys are all called `id` would point at every other. |
| `t3.c1` | `t5`, whose key is `c1` | One word again. In a database with opaque names nothing points. |
| `Customer.CustomerId` | `Customer` | It is the key itself. |

The rule reads the shape file only. It reads no row and calls nothing.

## 4. The concept rule

**The concepts that are searched.** The concepts of every graph that the MeaningGraph
registry's index lists with `kind: universal`, read at the commit the index gives, and the
concepts of the database's own meaning file when it has one. Today the index lists one
universal graph, `core`, with 28 concepts. No other graph is read, and nothing is searched
over a network: the files are fetched at their commit into the cache of
`open-layer-reading` REQ:fetched-files-are-pinned-and-rebuildable and matched on the
machine. When the database's meaning file pins `core`, the concepts are read at that pin
instead.

**A concept's words** are its labels and synonyms in every language, in lower case, as
MeaningGraph compares them (meaninggraph/core@25242a6 `FORMAT.md:137-139`).

| Rule | Fires when | Strength | Gives |
|---|---|---|---|
| `concept-names-table` | A table's words, or its singular form, joined by one space, are a word of exactly one concept of kind `entity`. | strong | `is-a` |
| `concept-names-column` | The column's table has a table concept, or `concept-names-table` points at one. A concept of kind `property` or `dimension` names, with `of`, that concept or one it extends. And the column's words, alone or after the table's words, are a word of that concept. Exactly one concept fits. | strong | `means` |

A table concept is the entity concept a table holds: the one bound to its record type with
the role `instances`, or the one an accepted `is-a` names.

When more than one concept fits a table or a column, each gets a `weak` pointer, and they
are the alternatives of one question. In `core` at `dd5ce32` no English word names two
concepts, so this does not happen on the six benchmark databases.

Found on the stripped Chinook, and nothing else is:

| Table or column | Concept of `core` | By |
|---|---|---|
| `Customer`, `Employee`, `Invoice`, `InvoiceLine` | `customer`, `employee`, `invoice`, `invoice-line` | the label |
| `Invoice.Total` | `invoice-total`, `of: invoice` | the synonym "total", and the label "Invoice total" |
| `InvoiceLine.UnitPrice`, `InvoiceLine.Quantity` | `unit-price`, `quantity`, both `of: invoice-line` | the label |
| `Employee.ReportsTo` | `manager`, `of: employee` | the synonym "reports to" |

Chinook's published meaning file binds all eight to a concept that extends the one named
(demo-db/chinook@88f080d `model/chinook.meaning.yaml`; the author's comparison).

Must not be found:

| Column | Concept the name equals | Why not |
|---|---|---|
| `Track.UnitPrice` | `unit-price` | `unit-price` is `of: invoice-line`, and `Track` has no table concept. Chinook's published file binds this column to `list-price`, which extends `price`. A match on the column's name alone would be a false line. |
| `Customer.Country`, `Employee.Country` | `country` | `country` is an entity concept, not a property of a customer. Chinook's published concept for `Customer.Country` takes its values from `country` and does not extend it. |
| `Customer.City`, `Customer.Company` | `city`; `organization` by its synonym "company" | Entity concepts again. |
| a column `Price` of any table | `price` | `price` has no `of`. |

A concept that a search finds by resemblance is never a pointer here. It appears only in
the list a person opens with "Means…".

## 5. Gates, the value count and the levels

| Gate | Passes when | Otherwise |
|---|---|---|
| `types-comparable` | Both columns have the same type family, by the tables of `open-layer-reading` (its appendix, section 7), and the family is known. | `failed` |
| `target-is-key` | The target column is the whole primary key of its table. | `failed`. A target that is only unique is not known to the shape file. |
| `values-present` | See below. | |
| `concept-free` (`is-a`) | The table is not a view and has no table concept, and the concept is the table concept of no other table of the database. | `failed` |
| `concept-free` (`means`) | The column has no written binding and no accepted `means`, and no entity file maps it. | `failed` |

**The value count.** For a column A and a target key B, in one environment:

- `distinct` is the number of distinct values of A that are not null, over all rows.
- `present` is how many of those values occur in B.
- `passed`: `distinct` is 1 or more and `present` equals it. `partial`: `present` is at
  least 95% of `distinct` and less than all. `failed`: less than 95%.
- `not-counted`: `distinct` is 0 (`no-values`); or the count did not finish in 10 seconds
  (`time-limit`); or the database could not be opened (`unreachable`).

`Invoice.CustomerId` against `Customer.CustomerId`: 59 and 59, `passed`.
`InvoiceLine.InvoiceLineId` against `Invoice.InvoiceId`: 2,240 and 412, `failed`.

It is two aggregate queries at most, in a read-only session, through the agent's read
path. It returns two numbers. No value of a row leaves the query.

**The levels.** Each candidate has exactly one. The first that applies wins.

| # | Level | When | A person sees |
|---|---|---|---|
| 1 | `recorded` | A pointer of strength `decides` names it: a declared foreign key. Or the statement is already known: from the database, from the model the database follows, or from an accepted answer. | The fact, as `open-layer-reading` shows it. No card. |
| 2 | `dropped` | A gate has `failed`. | Nothing. |
| 3 | `held` | No pointer; or `values-present` is `not-counted`. | No card. It is listed under "Refers to…". |
| 4 | `batch` | Every gate has `passed`; a `strong` pointer names it; no other candidate of the same statement kind and the same `from` is at level 4 or 5; and it is not half of a mutual pair. | A line of the batch card of its kind. Label `inferred`. |
| 5 | `question` | Everything else. | An alternative of a question card. Label `suggested`. |

A mutual pair is two `refers-to` candidates, each from a table's key column to the other
table, that both pass every gate: the two key columns hold the same values. They make one
question with three buttons. In pubs, `pub_info.pub_id` and `publishers.pub_id` are one.
When either column has another live candidate as well, each column keeps a question of
its own.

A question has at most four alternatives. Candidates of one `from` that would make more
are `held` instead.

## 6. The answers file

| | |
|---|---|
| Path | `dbmodels/<id>/<id>.answers.jsonl`, beside the source file. One per database. |
| Format | JSON Lines. Every line carries `"format":"datatug-answers/draft-1"`. |
| Written by | the agent and the `datatug advisor` commands, one line per action, at once |
| Written how | a line is appended. No line is ever changed or removed. |
| Committed | yes. DataTug stages it; the person commits. |
| Holds | statements, reasons, who and when. Never a value of a row. |

The four lines of the Feature's walk, then a colleague's note:

```json
{"format":"datatug-answers/draft-1","id":"ans_28b988b46e87a6d1","statement":"refers-to","from":"chinook/Invoice.CustomerId","to":"chinook/Customer","answer":"accepted","origin":"suggestion","reasons":["named after Customer's key CustomerId","Customer.CustomerId is the primary key","all 59 distinct values occur there (local)"],"by":"maria","at":"2026-10-12T09:02:11Z"}
{"format":"datatug-answers/draft-1","id":"ans_7fe05cc385a1ef9e","statement":"is-a","from":"chinook/Invoice","to":"meaning://github.com/meaninggraph/core/invoice","answer":"accepted","origin":"suggestion","reasons":["the table's name is the label of the universal concept invoice"],"by":"maria","at":"2026-10-12T09:02:40Z"}
{"format":"datatug-answers/draft-1","id":"ans_6172fa3333ec94f5","statement":"means","from":"chinook/Invoice.Total","to":"meaning://github.com/meaninggraph/core/invoice-total","answer":"accepted","origin":"suggestion","reasons":["\"total\" is a synonym of invoice-total, a property of invoice"],"by":"maria","at":"2026-10-12T09:02:40Z"}
{"format":"datatug-answers/draft-1","id":"ans_35d100b2bd3c07a4","statement":"refers-to","from":"chinook/Customer.SupportRepId","to":"chinook/Employee","answer":"accepted","origin":"stated","reasons":["stated by maria","all 3 distinct values occur in Employee.EmployeeId (local)"],"by":"maria","at":"2026-10-12T09:04:03Z"}
{"format":"datatug-answers/draft-1","id":"ans_dd059f3cda5eba93","statement":"refers-to","from":"chinook/Customer.SupportRepId","to":"chinook/Employee","answer":"open","origin":"suggestion","reasons":["the registered model chinook states it (not adopted)"],"note":"looks like a reference to Employee; not sure","by":"tom","at":"2026-10-13T16:20:00Z"}
```

| Key | Value |
|---|---|
| `id` | `ans_` and the first 16 hex characters of the SHA-256 of `<statement>\n<from>\n<to>\n<answer>\n<by>\n<at>` |
| `statement`, `from`, `to` | the statement, complete in itself (section 2) |
| `answer` | `accepted`, `rejected`, `open` or `withdrawn` |
| `origin` | `suggestion` when the person answered a card; `stated` when the person said it unasked |
| `reasons` | the reasons shown when the person answered, as text. A colleague reads them without the ledger. |
| `note` | required with `open`, optional otherwise; one line of the person's words |
| `concept` | with a `to` that is a new concept only: `kind`, `label` and `description`, as the person typed them |
| `by` | the principal the agent serves, or the operating system's user name for a command. Never a name the request carries. |
| `at` | the time of the action, UTC, to the second |

**The current answer** to a statement is the last line with that `statement`, `from` and
`to`. `withdrawn` means there is none. The proposal has three answer values (L2234);
`withdrawn` is the author's, and is how an undo is written without changing a line.

**Read again.** The file is read on every request. A line whose table or column is no
longer in the shape file stays in the file and is reported as a finding; it is used for
nothing until the name is back.

## 7. The three files a Keep writes

### The model

| | |
|---|---|
| Path, when DataTug creates it | `model/<module>/<module>.modelspec.hcl` |
| Path, when the source file already names a model in the project | that file, wherever it is |
| Form | ModelSpec HCL with `record` and `field`, the words of `1.0-draft-2`. No JSON copy is written. |
| Module name | the database's id, with every character outside `A-Z a-z 0-9 _` replaced by `_`, and `db_` put in front when it would start with a digit. `chinook-nofk` gives `chinook_nofk`. |

A generated record type, as it is written for `Invoice` (ModelSpec's released CLI v0.3.0
lints the whole generated file of 331 lines with no finding, under both its profiles; the
author ran it):

```hcl
record "Invoice" {
  key = ["InvoiceId"]

  field "InvoiceId" {
    type     = "int"
    required = true
  }

  field "CustomerId" {
    record   = "Customer"
    required = true
  }

  # … seven more fields
}
```

| What | Rule |
|---|---|
| Record types | one per table that is in the model, sorted by name, comparing bytes. A view gets none. |
| Record type name | the table's name with every character outside `A-Z a-z 0-9 _` replaced by `_`, and `_` put in front when it would start with a digit. `Order Details` gives `Order_Details`; `Sales.SalesOrderDetail` gives `Sales_SalesOrderDetail`. Two tables that would get one name: neither is generated, and both are reported. |
| Fields | one per column, in the engine's order, named by the same rule |
| `key` | the primary key's columns in their order; left out when the table has none |
| `type` | by the column's type family: integer `int`; text `string`; decimal `decimal`; float `float`; boolean `bool`; date `date`; time `time`; datetime and temporal `datetime`; uuid `uuid`; json `json`; binary `any`; a type text in no family `any`. These are the pairs that the drift report of `open-layer-reading` accepts. |
| `required` | `true` when the column refuses null; left out otherwise |
| A reference field | `record = "<record type>"` in place of `type`, for a declared foreign key of one column whose target is a table's whole primary key, and for an accepted `refers-to`. A field has one of the two, never both. |
| A foreign key of several columns | not written as a reference; reported with the reason `several-column-foreign-key` |
| Layout | two spaces; one blank line before each `field`; `=` signs of one field in one column |
| An existing file | DataTug changes one line to make a field a reference, and adds a `record` block at the end for a new record type. Every other line, comments included, stays as it was. |

### The meaning file

| | |
|---|---|
| Path, when DataTug creates it | `model/<module>/<module>.meaning.yaml`, beside the model, so that one folder is one graph |
| Format | `meaning/draft-2` (Question 3 of the Feature) |
| Pin | Every `extends` carries one pin: the pin the file already has, or, for a new file, the commit the registry's index lists for that graph. |

The whole file after the Feature's walk, 31 lines. The reference checker of
meaninggraph/core@25242a6 reports no problem for it and derives a third link, `invoice` on
`Invoice.InvoiceId` with the role `identifier`. The same file in the earlier spelling
passes the released checker v0.3.0. The author ran both.

```yaml
format: meaning/draft-2
id: chinook
name: chinook
description: What the database chinook means. Written by DataTug from the answers in dbmodels/chinook/chinook.answers.jsonl.
models:
  chinook: chinook.modelspec.hcl
concepts:
  - id: invoice
    kind: entity
    extends: meaning://github.com/meaninggraph/core/invoice?ref=dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7
    labels:
      en: Invoice
    description: >-
      A request for payment for goods or services supplied to a customer: what
      was supplied, when, where it is billed, and the amount charged. Not an
      order (a request to supply) and not a receipt (proof of payment). A
      dataset whose invoices are purchases says so with its own synonyms.
    bindings:
      - model: modelspec:///chinook.Invoice
        role: instances
  - id: invoice-total
    kind: property
    of: invoice
    extends: meaning://github.com/meaninggraph/core/invoice-total?ref=dd5ce32c4554c0833c59925bbea1e7a7e5ed01e7
    labels:
      en: Invoice total
    description: The amount charged on an invoice, normally the sum of its lines.
    bindings:
      - model: modelspec:///chinook.Invoice
        field: Total
        role: value
```

| What | Rule |
|---|---|
| A universal concept | It cannot be bound directly: a binding names a model of the same repository only (meaninggraph/core@25242a6 `FORMAT.md:305-309`). DataTug writes a concept of the project that `extends` it and binds that one. |
| Its `id` | the universal concept's id; when that id is taken, the table's words joined by `-`, a `-`, and the id |
| Its `kind`, `labels.en`, `description` | copied from the universal concept. The format requires all three, and the universal concept's own words are the only text at hand. |
| Its `of` and `values-of` | the concept of the project that extends the one the universal concept names |
| A new concept | written only for a statement whose `to` is new, with the label and the description the person typed, `kind: entity` for a table and `kind: property` for a column, and no `extends`. Its `id` is the label's words joined by `-`. A new property has `of`: its table's concept. The review marks it "new". |
| A new property on a reference field | The format lets such a field carry only the role `reference`, and that role needs `values-of`. DataTug writes `values-of`: the target table's concept. When the target has none, the statement cannot be kept yet, and the review says "Say what an Employee is first." |
| The role of `is-a` | `instances` |
| The role of `means` | `display-name` when the concept is the column's own table concept. Otherwise `value`; or `reference` when the field is a reference field, as the format requires (`FORMAT.md:297-304`). |
| A `refers-to` | writes no line here. MeaningGraph derives the link from the model once the target has a table concept. |
| The earlier spelling | `meaning/draft-1` writes `kind: attribute`, `property:`, `role: entity` and `role: foreign-key` for the four words above. |

### The source file

Written in the form "In the project" of `open-layer-reading` (its appendix, section 2).
The file a Keep creates for the walk, 18 lines:

```yaml
format: datatug-source/draft-1
id: chinook
model:
  hcl: model/chinook/chinook.modelspec.hcl
meaning:
  file: model/chinook/chinook.meaning.yaml
recordsets:
  - Album
  # … and the other ten tables, one line each
```

| What | Rule |
|---|---|
| A table whose record type and fields have its own names | a bare line, `- Invoice` |
| Any other table | a map: `name`, `record_type`, and under `columns` one map with the key `field` for each column whose field has another name |
| `meaning` | present only when a meaning file exists |
| An existing file | DataTug adds an item to `recordsets`, or the `meaning` key. Every other line stays as it was. |
| A registered model | Adopt writes the form "Registered": exactly the text that `datatug db match --source-file <address>` prints. |

## 8. The review

```ts
type Review = {
  waiting: number;                                   // statements the plan writes
  layers: {
    layer: 'meaning' | 'model' | 'database';
    heading: string;                                 // "Model: new, 11 record types; 2 references"
    sentences: string[];                             // "Invoice.CustomerId now refers to Customer."
    files: { path: string; change: 'new' | 'changed'; diff: string;   // a unified diff
             revision: string | null }[];            // the file's revision the plan was made from
  }[];
  checks: { id: 'model-lint' | 'meaning-check' | 'source-file' | 'model-and-database' | 'values-present';
            passed: boolean; blocks: boolean; line: string;
            findings: { file?: string; line?: number; message: string }[] }[];
  cannotBeKept: { answer: string; statement: string; from: string; to: string; finding: string }[];
  canKeep: boolean;                                  // every check that blocks has passed
};
```

| Check | What runs | Blocks a Keep |
|---|---|---|
| `model-lint` | ModelSpec's lint, default profile, on the planned model | yes |
| `meaning-check` | MeaningGraph's check on the planned folder, with each pinned graph supplied from the cache | yes |
| `source-file` | the source file's own reader of `open-layer-reading` on the planned file | yes |
| `model-and-database` | the drift report of `open-layer-reading` on the planned files: no finding of severity `error` | yes |
| `values-present` | the value count of each `refers-to` the plan writes: `59 of 59 distinct CustomerId values found in Customer` | no. A person's statement stands; the count is shown. |

OpenVaultDB's `ovdb publisher check` is not run: it reads a publisher's manifest in a
public repository, not this file (the proposal, L1870).

## 9. Routes

Relative to `/datatug`. `Scope` and `source` are as
[`core-investigation-loop/api-contract.md`](../core-investigation-loop/api-contract.md)
defines them. A table that the serving principal may not read whole is in no response.

```ts
type Reason = { text: string; kind: 'metadata' | 'statistics' | 'published';
                job: 'pointer' | 'gate';
                mark: 'points' | 'allows' | 'in-part' | 'rules-out' | 'says-nothing' };
type Line = { candidate: string; statement: Statement; from: string; to: string;
              sentence: string; reasons: Reason[]; notOffered?: string };
type Card =
  | { kind: 'batch'; group: 'references' | 'meanings'; label: 'inferred'; lines: Line[] }
  | { kind: 'question'; id: string; label: 'suggested'; title: string;
      alternatives: Line[]; note?: { by: string; text: string } };

// GET advisor/cards?source=…&table=…      table left out: the whole database
type Cards = { cards: Card[];              // at most three
               more: number;               // cards not shown
               waiting: number;
               canWrite: boolean; whyNot?: string;
               problems: string[] };       // "The database cannot be reached; references are not suggested."

// POST advisor/answers                    needs the write capability
type AnswerRequest = {
  operationId: string;
  answers: { statement: Statement; from: string; to: string;
             answer: 'accepted' | 'rejected' | 'open' | 'withdrawn';
             origin: 'suggestion' | 'stated'; note?: string;
             concept?: { label: string; description: string };
             ifCurrent: string | null }[]; // the id of the current answer the client saw
};

// POST advisor/follow                     needs no write capability
type FollowRequest = { candidate: string; value: Value };
type FollowResult = { via: 'candidate'; confirmed: false; notice: 'via a candidate link, not confirmed';
                      result: Result };    // the Result of POST semantic/related/rows

// GET  advisor/review?source=…            the type Review of section 8
// POST advisor/keep                       needs the write capability
type KeepRequest = { operationId: string; revisions: Record<string, string | null> };
type KeepResult = { written: string[]; staged: boolean; waiting: number };

// POST advisor/adopt                      needs the write capability
type AdoptRequest = { operationId: string; model: string };   // a modelspec:// address
```

No response has a key named `score` or `confidence`, and no number that is not a count.
The four routes of `core-investigation-loop` change in nothing.

## 10. Commands

```
datatug advisor suggest --directory <path> | --project <id>  --env <env> --db <db>  [--table <t>] [--recount] [--format text|json]
datatug advisor accept  <candidate>… | --batch references|meanings [--except <candidate>…]
datatug advisor reject  <candidate>…
datatug advisor answer  <question> --choose <candidate> | --none
datatug advisor state   refers-to <Table.Column> <Table>
datatug advisor state   is-a <Table> <concept> | --new <label> --description <text>
datatug advisor state   means <Table.Column> <concept> | --new <label> --description <text>
datatug advisor state   in-model <Table> | --all
datatug advisor note    <statement> <from> <to> --text <one line>
datatug advisor undo    <answer>
datatug advisor follow  <candidate> --value <v>
datatug advisor review  [--format text|json]
datatug advisor keep
datatug advisor adopt   <model address>
```

Every command takes `--directory` or `--project`, `--env` and `--db` as `suggest` does;
they are left out above. `--db` is the catalog's id, as in `datatug scan`. The commands
read and write the project folder as the person who runs them and apply no access policy,
as the `datatug db` commands of `open-layer-reading` do. `--format json` prints the body
of the route of the same name.

| Command | Changes | Exit code |
|---|---|---|
| `suggest`, `review` | the Advisor ledger only | 0; 2 for a bad flag or an unreadable project |
| `follow` | the Advisor ledger only | 0; 1 when the candidate is not offered |
| `accept`, `reject`, `answer`, `state`, `note`, `undo` | the answers file, one line for each statement, staged | 0; 1 when the statement cannot be answered, with the reason |
| `keep` | the model, the meaning file and the source file, staged | 0; 1 when a check that blocks has failed or a file changed since the review, with nothing written |
| `adopt` | the source file, staged | 0; 1 when the model does not match or a source file exists |
