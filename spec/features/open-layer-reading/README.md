---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Open-Layer Reading

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/datatug/datatug/spec/features/open-layer-reading?op=explore) | [Edit](https://specscore.studio/app/github.com/datatug/datatug/spec/features/open-layer-reading?op=edit) | [Ask question](https://specscore.studio/app/github.com/datatug/datatug/spec/features/open-layer-reading?op=ask) | [Request change](https://specscore.studio/app/github.com/datatug/datatug/spec/features/open-layer-reading?op=request-change) |
**Status:** Draft
**Date:** 2026-10-10
**Owner:** alex
**Source Ideas:** —
**Supersedes:** —
**Kind:** Capability (implementations in `datatug-core`, `datatug-cli` and `datatug-apps`)
**Pillars:** P3 semantic layer · P5 related-data navigation · P12 provenance · P13 infrastructure inventory

## Summary

Specifies, and does not yet build, how DataTug will show what a database's columns mean and which rows are related from three open descriptions: a ModelSpec model, its MeaningGraph meaning file and the database's own description. Nobody writes an entity file by hand for it, and every meaning, key and reference says how it is known.

This is the specification of what the ModelSpec design proposal calls "Phase 4, DataTug reads". Nothing in it is built, and none of it is part of the launch planned for 2026-10-25. The status is Draft. Only the owner approves it. Until he does, nothing here may be built. The proposal places this phase "after the foundation, and after the launch" (line 2796 of its HTML file, written L2796 below). Whether building may start before the launch, once this is approved, is not decided here.

How to read it: "What a person sees" is the whole feature in four moments. "Questions for the owner" holds the six answers it needs. "The author's choices" lists 26 things the author chose and the owner may overrule. The rest, and the appendix [formats.md](formats.md), is for whoever builds and tests it.

## What a person sees

Three moments, all on Chinook, and a fourth that is not a step of the journey. Chinook is a music shop's database with 11 tables, 64 columns, 59 customers and 412 invoices. The three moments are steps 1, 2 and 9 of the journey in the proposal's section 14 (L2749, L2750, L2757).

**Step 1. "I point DataTug at my database."** Maria has a SQLite copy of Chinook. She runs one command:

```
datatug scan --directory ./shop --driver sqlite3 --path chinook.sqlite --db chinook --env local
```

She opens the project. The database page lists her 11 tables. Each has a mark that says how many of its columns DataTug can say something about:

```
Invoice      2 of 9
Customer     2 of 13
```

The database file is byte for byte what it was. DataTug opened it read-only.

**Step 2. "I do nothing."** She opens the table `Invoice` and its Keys tab. Without any action of hers it reads:

```
Primary key    InvoiceId                from the database
Foreign key    CustomerId ⇒ Customer    from the database
```

The database page has one more line:

```
Registered model chinook matches this database: 11 of 11 tables, 64 of 64 columns.
```

DataTug says so and does nothing more. No file was written. No column has a meaning yet.

**Between step 2 and step 9: the one thing Maria does by hand.** Question 2 asks whether it should be one click instead. She runs

```
datatug db match --directory ./shop --env local --db chinook --source-file modelspec://github.com/demo-db/chinook/chinook > shop/dbmodels/chinook/chinook.source.yaml
```

The file has 19 lines. It names the model at a commit and lists the 11 tables:

```yaml
format: datatug-source/draft-1
id: chinook
model:
  address: modelspec://github.com/demo-db/chinook/chinook?ref=3e7bb316d8a438eb7e5930067997a4e4543a275c
meaning:
  graph:
    address: meaning://github.com/demo-db/chinook
recordsets:
  - Album
  # … and the other ten tables, one line each
```

She commits the project.

**Step 9. "My colleague opens the same database and does nothing."** Tom has the same database file. He clones the project, starts his agent on it and opens `Invoice`. He scans nothing and edits no file.

- The header of the column `CustomerId` has a marker. Its tooltip reads `customer.id · published`: this column holds the identifier of a customer, and a registered model says so.
- The Keys tab still says `CustomerId ⇒ Customer`, from the database, and now adds that the model says so too: `also published @ 3e7bb31`.
- He clicks the cell `CustomerId = 5`. The panel offers `Customer` (1 row). In the table `Customer` he clicks `CustomerId = 5`, and the panel offers `Invoice` (7 rows), once.
- The marks have filled: `Invoice 5 of 9`, `Customer 3 of 13`.
- The database page says where it all comes from: `Model: modelspec://github.com/demo-db/chinook/chinook @ 3e7bb31 · published`, and under it `Model and database agree.`
- The project has no `entities` folder. Nobody wrote an entity file.

**Later. The database changes.** Someone drops the column `Total` from `Invoice` in the local copy, and Maria scans again. Under the model line the database page now reads:

```
Model and database differ in 1 place.
```

She opens that line, or runs `datatug db drift --directory ./shop --env local --db chinook`, and reads:

```
error  field-without-column  Invoice.Total  The model's field Total has no column in the table Invoice.
```

DataTug changes nothing. She decides which side to change: the model, the source file or the database.

### One fact, followed: an invoice belongs to a customer

| Where the fact is | Who wrote it | The real lines | How DataTug labels it |
|---|---|---|---|
| The database | Chinook's authors | `FOREIGN KEY ([CustomerId]) REFERENCES [Customer] ([CustomerId])` (demo-db/chinook@c7535f4 `data-source/sql/Chinook_Sqlite.sql:144`) | |
| The shape file, `environments/local/catalogs/chinook/chinook.shape.json` | `datatug scan` | `{ "column": "CustomerId", "table": "Customer", "referencedColumn": "CustomerId", "constraint": 0, "position": 0 }` | `observed` |
| The source file, `dbmodels/chinook/chinook.source.yaml` | Maria | `- Invoice`: the table `Invoice` holds Invoice rows, and each column is the field of the same name | |
| The model, read at commit `3e7bb31` | Chinook's publisher | `field "CustomerId" { record = "Customer" required = true }` (demo-db/chinook@3e7bb31 `model/chinook.modelspec.hcl:340-343`) | `published` |
| The meaning file, read at the same commit | Chinook's publisher | `- model: modelspec:///chinook.Customer` / `role: entity`, and `- model: modelspec:///chinook.Invoice` / `property: CustomerId` / `role: foreign-key` / `note: The customer who made a purchase.` (`model/chinook.meaning.yaml:201-206`) | `published` |
| DataTug's runtime index, in memory | the importer | the entity field `customer.id` is mapped to `Customer.CustomerId` and to `Invoice.CustomerId` | |
| The screen | | marker `customer.id · published`; Keys tab `CustomerId ⇒ Customer`, from the database; related rows `Invoice` (7) | |

Chinook's meaning file is still written in `meaning/draft-1` (`model/chinook.meaning.yaml:1`), so three words in those lines are the earlier ones: `property:` is today's `field:` and names the model's field; `role: entity` is today's `role: instances`; `role: foreign-key` is today's `role: reference` (meaninggraph/core@65510ae `FORMAT.md:376-379`). The importer reads both spellings.

## Words used here

One name for one thing. Where a word is used at two levels, its row says so. The right-hand column says what today's code and features call it.

| Word here | Means | Elsewhere |
|---|---|---|
| database | What a project knows under one id across its environments. The id is the `dbModel` that a catalog file names, or the catalog's own id when it names none. | DataTug's wire key `source` takes this id. |
| catalog | One environment's copy of a database: the folder `environments/<env>/catalogs/<db>/`. `<db>` is the id given to `datatug scan --db`, and may differ from the database's id. | The wire key `source` takes this id too, as a second name. |
| table, column | What the database's engine holds. A view counts as a table unless a line says otherwise. | OpenVaultDB: recordset, column. DataTug's wire: `collection`, `column`. |
| model | A ModelSpec model: one module of record types. Never an AI model, and never DataTug's `dbmodels/` folder. | |
| record type, field | ModelSpec's words for a kind of row and one member of it. | ModelSpec before `1.0-draft-2`: entity, property. |
| key, reference | Used at both levels. A key is what identifies a row: the database's primary key, the record type's `key`. A reference says that a column, or a field, holds the key of another table's, or record type's, rows. "Primary key" and "foreign key" always mean what the database declares. A field that is a reference is written "reference field". | The name of a JSON or YAML entry is also called a key; it is then written with its name, as in "the key `field`". MeaningGraph's role `reference` is written in code type. |
| meaning file, concept, binding | MeaningGraph's words. A binding is a line that says a field holds a concept. | |
| entity concept, property concept | MeaningGraph's kinds of concept: `entity`, a thing with instances; `property`, a property of an entity. `dimension` and `measure` are two more. Always written with "concept", so that "entity" alone stays DataTug's. | meaninggraph/core@65510ae `FORMAT.md:134-136`. `meaning/draft-1` writes `kind: attribute` for `property`. |
| link | One fact "this record type, or this field, holds this concept in this role". Written when a binding line says it, derived when MeaningGraph's rule computes it from the model. | meaninggraph/core@65510ae `FORMAT.md:307-366` |
| entity, entity field | DataTug's existing runtime types, `datatug.Entity` and `EntityField`, as in `Customer.ID`. Always two words for the second: "field" alone is ModelSpec's. | datatug-core@eb3075e `pkg/datatug/entities.go:41-47`, `:107-117` |
| entity file | A hand-written `entities/<id>/<id>.entity.json`. | |
| open layer | One of the three open descriptions this feature reads: the model, the meaning file, and the description of the database, which here is the shape file. | Proposal D8: "the three open layers". |
| shape file | What a scan found, written into the project in the shape of OpenVaultDB's descriptor. | The proposal: "the scan's output". "Descriptor" is avoided: in `datatug-cli` it names a connection file. |
| source file | The hand-kept project file that names a database's model and holds its mapping. | The proposal: "source file" (L1567, and the card of D5 question 3, L2935). The question he answered said "a file inside the DataTug project". |
| mapping | The lines that say which table holds which record type and which column is which field. | |
| runtime index | The two in-memory inputs of the resolver: the entities, and each table's keys and foreign keys. It is not a file. | The proposal: "runtime index". |
| resolver | The three pure functions of `datatug-core` `pkg/semantic`: `Resolve`, `RelatedLookups`, `Applicable`. | The proposal uses "resolver" for a second thing in its Phase 7; not here. |
| importer | The function that builds the runtime index from the shape, the mapping, the model and the links. | |
| label | The word that says how a fact is known: `observed`, `published`, `declared` or `inferred`. | Today's wire keys `provenance` (a column) and `mapping` (a fact); `known` in the three new routes. It is not MeaningGraph's `labels`, which are a concept's names, and not the wire key `label`, which is a related item's display name. |
| mark | The "n of m" beside a table: how many of its m columns are described. A column is described when it has a meaning, or is part of a key, or is a reference. | The proposal: "a mark for how much is described". |
| model line | The one line on the database page that says which model the database follows. | |
| registered model | A model listed in the index of `modelspec-org/registry` at a commit. | |
| match | A registered model is fully present in a database under the same names. | Not the feature `compare`. |
| drift report | The list of differences between a model and a database as last scanned. | The proposal: "a check of a model against the database". Not a `Check` of the feature `checks`. |

The proposal is `/Users/alex/projects/modelspec-conceptual-design-proposal.html`: a local file in no repository, md5 `fb7f6813691fcba8984760c51b3b725c`, cited with its HTML line as `L…`. The phase is in its section 14. Its Phase 4 is not Phase 4 of this repository's `ROADMAP.md` ("project grants and distribution", `ROADMAP.md:137`), not the Public Data Fabric programme's Phase 4, and not pillar P4. "Proposal D5", "proposal D7" and "proposal D8" are the proposal's decision cards, not the decisions D5, D7 and D8 of `spec/product/product-family.md`.

## Problem

A person who connects a database to DataTug today gets its meanings only by writing entity files by hand, even when a model and a meaning file for that database are already published.

- **Meanings come from hand-written files only.** The agent builds the resolver's entities from `entities/**/*.entity.json` and nothing else (datatug-cli@3a1c276 `pkg/server/endpoints/semantic_project.go:27-47`). The demo project carries eight such files (datatug-demo-project@deeee32 `demo-project-1/entities/`). No Go code in `datatug-core` or `datatug-cli` reads a ModelSpec model or a meaning file.
- **A scan and the meanings are not connected.** A scan stores no foreign keys (datatug-cli@3a1c276 `pkg/api/scan_layout.go:23-27`: "Foreign keys are not stored."). The scanners read them, and the save leaves them out (`pkg/schemers/dalgoschema/provider.go:519-547` for PostgreSQL). The meanings routes ask the live database instead, and leave a SQL table's own foreign keys empty (`pkg/server/endpoints/semantic_schema.go:177-184`).
- **A database that is not public has nowhere to say which model it follows.** The publisher manifest needs a public address. Proposal D5 question 3 gives such a database a source file in the DataTug project; no specification for it exists.
- **Two labels cannot say what the owner approved five for.** The wire carries `declared` and `inferred` only (datatug-core@eb3075e `pkg/apicontract/responses.go:8-11`, `:33`). A fact that a publisher registered would have to be shown as `declared`, as if someone in the project had said it.
- **Nothing compares a model with a database a person has connected.** The three format checkers read files at a commit (the proposal, L1586). Chinook's own CI compares its model's keys with its fixture file (demo-db/chinook@c7535f4 `scripts/test-data.mjs:104-122`). That is one publisher's test of its own database, not something a person can run on theirs.

## Decided by the owner, and where

The owner's words are in the column "His words", and nowhere else in this section. Every other quotation here is the text of a question, a card, a record or the proposal, written by its author.

| Decision | What was put to him | His words | Where it is recorded |
|---|---|---|---|
| Proposal D5, question 1 | "**Q1. One line per table in a public database's manifest.**" with Northwind's lines before and after (decision 0012 `:121-136`) | "Q1 - approved." (2026-10-08, 21:34 UTC) | openvaultdb/openvaultdb@37d2389 `spec/decisions/0012-publisher-manifest-maps-record-types-and-fields-in-the-list-of-recordsets.md:162`, `:421` |
| Proposal D5, question 2 | "**Q2. Columns whose names differ from the field.** Nothing can say this today. Proposed, as a made-up example:" with the column written bare, `email_address: email` (`:138-145`); then the proposal author's reply with the bare form and a map with `field:` side by side (`:180-192`) | "Q2 - how about making it more expandable in future by using column attributes, something like", with his example; then "approved the "field" attribute of a column" (2026-10-08, 21:34 and 21:36 UTC) | the same record, `:164-175`, `:202`, `:422` |
| Proposal D5, question 3 | "**Q3. A private database.** The manifest requires a public address, so an internal database has nowhere to write Q1 or Q2. Proposed: the same lines in a file inside the DataTug project:" with the six-line example below (the same record, `:147-156`). The proposal's card restates it after his answer as "It carries the same lines in a source file inside the DataTug project." (L2935-L2941); the record says the card is "not the wording he answered" (`:207-209`). | "Q3 - approved" (2026-10-08, 21:34 UTC) | the same record, `:178`, `:423`; the proposal's card, L2944 |
| Proposal D7 | "Knowledge that is not yet certain stays out of the published formats. DataTug keeps signals and candidates on the user's machine and answers in the project. Five states carry the labels `observed`, `inferred`, `suggested`, `declared` and `published`. Whether a candidate is offered is decided by checks that can rule it out and pointers that name a target, not by a score." (proposal L2961) | "Approve" (2026-10-08) | the proposal's card, L2964. The card is the proposal author's transcription. The original message was not read for this specification. |
| Proposal D8 | "DataTug's own entity files are replaced by the three open layers, in stages: first an importer that feeds the existing runtime, then DataTug writes the open formats directly. Retiring the entity files and types is a separate decision later." (proposal L2969) | "Approve" (2026-10-08) | the proposal's card, L2972. The card is the proposal author's transcription. The original message was not read for this specification. |
| Proposal D6, part a | "The link from a property to a concept is derived from the model's reference or key, by one rule written into MeaningGraph's format. A binding line that only restates the model becomes optional." | "D6a: Approve" | meaninggraph/core@65510ae `spec/decisions/0001-derived-links-and-two-role-names.md:47`, `:70`. In that card "property" is the model's field: the proposal's second revision writes "field" there (same record, `:54-57`). |
| This specification | A list of what was left after Phase 3, among it "Phases 4 to 8 … Each needs a specification you approve first." | "Do cleanup, specify phase 4 and if make sense others" (2026-10-10) | the session record of that day; it is in no repository |

The example he approved for question 3:

```yaml
id: payments
model: { hcl: model/payments.modelspec.hcl }
recordsets:
  - name: payments
    record_type: Payment
    columns: { order_ref: OrderId }
```

What each allows, in the proposal's words, which are its author's and not the owner's:

- Proposal D5 question 3: "a specification for the project's source file; building it is Phase 4." (L2943)
- Proposal D7: "a feature specification for the ledger and the Advisor. Building them is Phase 5." (L2963)
- Proposal D8: "specifications for the importer (Phase 4) and for writing the open formats (Phase 5). It does not authorise removing anything." (L2971)
- The table of phases, for Phase 4: it "Needs D5 question 3, D7 for the labels, D8. All approved.", allows "Writing the specification for the importer and the project's source file" and does not allow "Building them before that specification is approved". (L2836)

What he did not decide, so that nobody reads it into the rows above:

- The list of phases, what Phase 4 holds and its "Done when" line are the proposal author's plan (proposal section 14). The owner approved cards, not that section.
- What each label means, the words on screen, the rules of the mapping beyond the key `field`, the mark, the match and the drift report are the proposal author's text (its sections 7, 8 and 9). This specification takes some of it and changes some of it. The larger points are in the next section as the author's choices. Everything in the appendix is the author's too, as its first paragraph says.
- The exact wording of MeaningGraph's rule for derived links is not his either (meaninggraph/core@65510ae `spec/decisions/0001-…:143-145`, point N33).
- His instruction of 2026-10-10 allows writing this specification. It approves nothing in it.

## The author's choices

Each is the choice of this specification's author. None is the owner's ruling. Each can be changed by one word from him.

| # | The choice | Why | The alternative |
|---|---|---|---|
| 1 | One feature with an appendix, not a parent with sub-features. | The seven parts share one journey and one end-to-end gate, as `core-investigation-loop` does with its `api-contract.md`. | A parent and three sub-features: the two files; the importer and labels; match and drift. |
| 2 | The names "open-layer reading", "match" and "drift report". | "Phase 4", "check" and "compare" each already mean something else in this repository. | Keep the proposal's words and qualify them every time. |
| 3 | The shape file holds the descriptor's list of recordsets and nothing else of the descriptor, under its own format name. | The descriptor's top level requires a public address, a publisher and licences (demo-db/chinook@c7535f4 `schemas/ovdb-database-draft-1.schema.json`). A private database has none. | Ask OpenVaultDB for a private variant of the descriptor first. |
| 4 | One shape file per environment's copy of the database; one source file per database. | What an engine reports can differ between `local` and `prod`. What the tables mean does not. | One shape file per database, replaced by whichever environment was scanned last. |
| 5 | The two files sit in folders a scan already makes: `environments/<env>/catalogs/<db>/` and `dbmodels/<id>/`. | No new folder; they travel with the files they describe. | Two new collections, `dbshapes/` and `dbsources/`. |
| 6 | `rowCount` is `null` and `description` is `""` in the shape file. | A scan runs no count query and has no description. Counting rows is a data query. | Count rows in the scan. |
| 7 | The source file is YAML, with a `format` line and a `meaning` key that the approved example does not have. | "The same lines" as a manifest are YAML. Without `meaning` no meaning can be found. Without `format` a later change cannot be told from a mistake. | No `format` line, as in the approved example; JSON, as `dalgo-project-store` asks of project items. |
| 8 | A registered model is named by address and commit and read from its repository at that commit, into a cache outside the project. It is not copied into the project. | The proposal says DataTug reads "each by address at a commit". A copy could be edited and still be called published. | A command that copies the model and meaning file into the project. |
| 9 | The importer is Go and runs in the agent, on the reference Go libraries where they can read the format. | The resolver it feeds is Go and runs there. The web app's reader of these formats serves another feature and another lane. It follows that a project opened in the hosted web app with no agent shows no meanings (see "Not Doing"). | Reuse the TypeScript reader of the web app's public-data route. |
| 10 | An imported entity takes the id of its entity concept; its entity fields are `id`, `name`, or the id of a property concept. So `customer.id`, not `Customer.ID`. | A concept's id is its public, stable name. Two record types may hold one concept. | The model's names, as in `Customer.CustomerId`. |
| 11 | One meaning for one column, by a fixed order (appendix, section 3). | The header marker shows one. `Customer.SupportRepId` has two links, and `employee.id` is the one a person can follow. | Show every link of a column. |
| 12 | An entity file's declared mapping wins over an imported one. Imported entities are not listed by `datatug entity list`. Example: with a source file added to the demo project, `Customer.CustomerId` still reads `Customer.ID · declared`. | "Declared mappings are authoritative" is a Stable requirement. An imported entity is not a file. A column that an entity file reached by a name pattern moves to the imported entity. A saved query tagged with the entity file's id is then not offered from that column. In the demo, `Invoice.CustomerId` reads `Customer.ID · inferred` today and would read `customer.id · published`; the saved queries tagged `Customer.ID` would no longer be offered from that cell. | Imported meanings win; or list both; or let an entity file's name pattern win too. |
| 13 | `observed` is a label of keys and references only, never of a meaning. A reference that both the database and the model state shows `observed`, with the model named beside it. Example: `CustomerId ⇒ Customer` reads `from the database · also published @ 3e7bb31`. | A database does not say what a column means. About itself, the database is the first witness. | One label by "who said it last"; or the model's label first. |
| 14 | A related row does not say how it is known. Its wire key `label` is its display name and stays. A reference's label is shown where the reference is listed, on the Keys tab. Example: the panel's `Invoice` (7 rows) does not say how it is known; the line `CustomerId ⇒ Customer` on the Keys tab does. | It leaves two of the four existing responses untouched in shape. | A new key on every related item. |
| 15 | A column counts as described when it has a meaning, or is part of a key, or is a reference. Example: `Invoice 2 of 9` with no model: `InvoiceId` is the key and `CustomerId` is a reference. | It is what the proposal's mock-up shows under "What is known". | Count meanings only. |
| 16 | A registered model matches when all of it is present in the database under the same names. Extra tables and columns in the database do not prevent a match. Partial matches are not reported. Example: Chinook with an added table `audit_log` still matches, `11 of 12 tables`. Chinook without `Invoice.Total` does not. | A full match is the one case where a source file of bare names is enough. A guess at a partial match is a suggestion, and suggestions are a later phase. | Report every model above a share of matching names. |
| 17 | The drift report compares the model with the shape file, not with the live database. It reports no value counts. Example: a column dropped today is reported after the next scan, not before. | It needs no connection and no query, and a colleague gets the same answer. Value counts are the next phase's. The proposal reads the other way in one place: "A check that opens the database and compares it with the model is proposed in section 14." (L1884). It also lists, with no phase, a check "that every reference's values are present in its target, with the count shown" (L1586); that one is left out here. A person gets the live comparison by scanning first. | Open the database for every report. |
| 18 | Three new read-only routes and three commands under `datatug db`. | `datatug db` exists (`copy`, `export`). `datatug check` and `datatug compare` are taken. | New top-level commands. |
| 19 | The importer shows no link of a meaning file that fails MeaningGraph's check. | The format says the result is "not defined" for such a graph (meaninggraph/core@65510ae `FORMAT.md:330-331`). | Show written links and hold back derived ones. |
| 20 | Foreign keys live in the shape file. The Draft file `<model>.refs.json` of `datatug-core` is neither written nor read. | One place for one fact. `datatug-core` has had a reader and a writer of that file since v0.43.0 (2026-10-05, `pkg/storage/filestore/model_foreign_keys_file.go:88-97`); no command calls them. | Write foreign keys to `refs.json` and the rest elsewhere. |
| 21 | SQL databases only: `sqlite3`, `postgres`, `sqlserver`. | They are what `datatug scan` reads. The descriptor's shape is tables and columns. | Include inGitDB collections and HTTP sources. |
| 22 | The shape file is committed with the project. | Step 9 needs it: a colleague who clones the project scans nothing. Today's scan files are committed too. | Keep it out of Git and make it again when a colleague connects. |
| 23 | A fact is `published` when its model is named by a registered address, at whatever commit the source file names. | The source file pins what a project reads; the registry's own pin moves on. | `published` only at the commit the registry lists; another commit is refused. |
| 24 | The source file refuses the key `modules`. | Pairing two databases is the proposal's Phase 7. A file that carries it must not be half-read. | Accept the key and ignore it. |
| 25 | Each registry's index is read at the head of its `main` branch. Two settings name another address or a local file. Example: the end-to-end test sets them to the index files of `modelspec-org/registry@a7325ea` and `meaninggraph/registry@49f472f`, so that its model line reads `3e7bb31` on any day. | The index is the one list that has to be current, or a newly registered model is never found. A test needs a fixed one, and a setting is not an intercepted request. The two settings are an independent reviewer's suggestion that the author took. | Ship a copy of each index with every DataTug release. |
| 26 | The three new routes list and count a table only when the person asking may read all of it. Example: in the demo project the role `support` may read `Customer` for Canada only and without `Email` (datatug-demo-project@deeee32 `demo-project-1/policies/customers.yaml`); for that role `Customer` is in no list and no count. | No route that describes tables asks the access policy today (datatug-cli@3a1c276 `pkg/server/endpoints/semantic_columns.go:33-97`), so the rule is new. Hiding is the safe side: a count must not tell what a policy hides. The rule is an independent reviewer's suggestion that the author took. | List a narrowed table with the columns that person may read. |

## Questions for the owner

Six questions. Each changes what is built. The requirements below are written for the recommended answer, and each one that would change names its question.

### Question 1. May a column in the source file be written as a bare string?

Your approved example writes it bare:

```yaml
columns: { order_ref: OrderId }
```

Question 3 was put to you as "the same lines" as questions 1 and 2. In the message that approved it you also asked for question 2's column to become a map (decision 0012, `:147`, `:162-178`).

A publisher manifest must write it as a map, and its reader refuses the bare form with `column "order_ref" must be a map with field:`:

```yaml
columns:
  order_ref:
    field: OrderId
```

That was the coordinating session's decision of 9 October, not yours. Its record says the bare form "falls due again when the file of question 3 is specified" (openvaultdb/openvaultdb@37d2389, decision 0012, point N5, `:482-492`). It is due now.

- **A. Map only (recommended).** One grammar and one reader for both files. Your example is then written `order_ref: { field: OrderId }`. The bare form can be added later without breaking any file.
- B. Both forms in the source file, map only in a manifest. Then the two files are not "the same lines".
- C. Both forms in both files. This reopens the manifest's rule.

### Question 2. When a registered model matches, does DataTug offer to adopt it in this phase?

After a scan of Chinook, DataTug says:

```
Registered model chinook matches this database: 11 of 11 tables, 64 of 64 columns.
```

- **A. It says so, and prints the source file (recommended).** `datatug db match --source-file modelspec://github.com/demo-db/chinook/chinook` prints the 19 lines; a person saves them as `dbmodels/chinook/chinook.source.yaml`. DataTug writes nothing. The proposal's own order is "readers before writers" (L2745), and it places "DataTug writing the model, the meaning file and the source file" in Phase 5 (L2807).
- B. It also offers one confirmation, a button and a command, that writes exactly that file. The proposal's walk-through has this (L1569: "It offers to adopt it. With that one confirmation the model, its references and its published meanings are in place"). It adds a write route to the agent while another lane is building the project write contract.

The proposal's rule for this phase's test is: "Each phase from 4 onward ends with an end-to-end test that walks its part of this journey without a reload or a manual step between stages." (L2761). With A the walk has one step done by hand, saving the printed file, so A does not meet that sentence to the letter; REQ:journey-test names that step as its one exception. With B the step is one confirmation inside DataTug, as in the proposal's walk-through (L1569). The rule is the proposal author's text, not your ruling.

### Question 3. May the label of a Stable contract gain the value `published`?

Today the agent answers, for the demo project:

```json
{ "column": "CustomerId", "entity": "Customer", "field": "ID", "provenance": "declared" }
```

With this feature it would answer, for a database that follows the registered model:

```json
{ "column": "CustomerId", "entity": "customer", "field": "id", "provenance": "published" }
```

`provenance` is a closed list of two values in `core-investigation-loop`, which is Stable (`api-contract.md:72`, `:175`), checked on the server and in the web app. The proposal asks for the new labels and for unchanged screens in one breath (L2800, L2802). Both cannot hold to the letter.

- **A. Yes (recommended).** `provenance` and a fact's `mapping` gain `published`. `core-investigation-loop` is amended when this is built: Stable, Amending, Stable. Two more things change in the four routes, in behaviour and not in shape: `POST semantic/related` never lists one lookup twice (REQ:related-rows-once), and for a project that has a shape file it takes keys and foreign keys from that file (REQ:routes-unchanged).
- B. No. Imported meanings are sent as `declared`, and only the database page says they are published. A registered fact then reads as if someone in the project had declared it.

### Question 4. Does a PostgreSQL scan write its foreign keys into the shape file in this phase?

A PostgreSQL scan already reads foreign keys with their target tables and columns. It does not save them (datatug-cli@3a1c276 `pkg/schemers/dalgoschema/provider.go:519-547`; `apps/datatugapp/commands/cmd_scan_db.go:322`: "The scan does not save foreign keys or indexes."). Writing them into the shape file needs no new reader. If they are not written, step 2 on a PostgreSQL copy of Chinook shows the primary key and no foreign key:

```
Primary key    InvoiceId    from the database
```

If they are written, it shows what a SQLite copy shows:

```
Primary key    InvoiceId                from the database
Foreign key    CustomerId ⇒ Customer    from the database
```

- **A. Yes (recommended).** The shape file holds every foreign key the scan read. A key that points into another schema, or at a table the role may not see, is left out by today's scan with a line in its log; the shape file then says that some were left out (`"foreignKeys": "partly-read"`). The author expects most real projects to be on PostgreSQL; that is an expectation, not a measurement.
- B. No. SQLite only. A PostgreSQL shape file says `"foreignKeys": "not-read"`, although the scan read them. PostgreSQL shows keys from the database and references from a model, when one is named. Its foreign keys follow later. This saves the one criterion that needs a PostgreSQL database loaded with Chinook (AC:foreign-keys-read-or-not).

The SQL Server scan reads foreign keys with their target columns too (`pkg/schemers/mssqlschema/constraints.go:32`; read in code, not run). The proposal's "Done when" names Chinook and no engine (L2803). The end-to-end test here uses SQLite either way.

### Question 5. Does the scan write the shape file beside today's files, or in place of them?

Today a scan of Chinook writes eleven files like `dbmodels/chinook/main/tables/Invoice/main.Invoice.columns.json`: columns and primary keys, no foreign keys. The proposal says the scan's output "uses OVDB's descriptor shape" (L1567).

```
dbmodels/chinook/main/tables/Invoice/main.Invoice.columns.json    today; kept under A, gone under B
environments/local/catalogs/chinook/chinook.shape.json            new
```

- **A. Beside (recommended).** The scan writes both from one result, so they cannot disagree. Every reader of today's files keeps working untouched: the CLI's catalogue readers, chat, saved queries and the web app's GitHub reader. Retiring the columns files is a later change with its own decision.
- B. In place. The columns files go, and every one of those readers moves to the shape file in this phase. It is one shape for one fact at once, and it reaches into code the launch lane is working in.

### Question 6. Is "Chinook with no hand-written entity file" a second project, or the demo project changed?

The demo project shows `Customer.ID · declared` from `entities/Customer/Customer.entity.json`. Its saved queries bind on that id (`"meta": {"entity": "Customer", "field": "ID"}`), a Stable requirement asks for those files (`core-investigation-loop` REQ:demo-knowledge-project), and the required journey tests look for the text `Customer.ID = ` followed by the customer's id (datatug-apps@098db23 `apps/datatug-app/e2e/journey/journey.spec.ts:430`). An imported meaning reads `customer.id · published`.

- **A. A second project (recommended).** The end-to-end test makes its own project by scanning Chinook into an empty folder. `demo-project-1`, its eight entity files, its queries and journeys J1 to J4 are not touched. Chinook's model could not replace all of them anyway: `Country` and `Person` are not in it, nor is the mapping of `Customer.ID` to the `support-notes` collection.
- B. The demo project changed. Its Chinook entity files are removed in this phase. That amends three clauses of a Stable feature, re-tags the saved queries and rewrites the J1 to J4 assertions. The card of proposal D8 says its approval "does not authorise removing anything" (L2971). B would need your word for the removal itself.

Retiring entity files in general is not asked here. Proposal D8 keeps that for a later decision.

## Behavior

Requirements use MUST, MUST NOT, SHOULD and MAY as fixed terms. "Appendix" means [formats.md](formats.md).

### The interface that stays fixed

#### REQ: resolver-unchanged

The exported functions `Resolve`, `RelatedLookups` and `Applicable` of `datatug-core` `pkg/semantic`, and the types they take and return, MUST keep their signatures and their behaviour (datatug-core@eb3075e `pkg/semantic/resolve.go:62`, `related.go:99`, `applicable.go:50`). The importer MUST feed them through their existing inputs only: a list of `*datatug.Entity`, and a map from `semantic.SchemaKey` to `semantic.TableSchema`. The change that builds this feature MUST NOT alter a non-test file of `pkg/semantic`.

#### REQ: routes-unchanged

The four routes `GET semantic/columns`, `POST semantic/related`, `POST semantic/related/rows` and `POST queries/applicable` MUST keep their paths, their request shapes and the keys of their responses. Three things change. The value `published` is allowed in `provenance` and in a fact's `mapping` (Question 3). `POST semantic/related` lists one lookup once (REQ:related-rows-once). For a project that has a shape file, `POST semantic/related` reads keys and foreign keys from that file, in both directions, where today it asks the database and leaves a table's own foreign keys empty. For a project that has no shape file and no source file, each of the four MUST answer as it does today: the same status, keys and values, apart from the time in `provenance.observedAt` and apart from a lookup that today is listed twice.

#### REQ: no-new-screen

The web app MUST gain no page and no address. It gains exactly what the group "What the screens show" lists and the label in an applicable query's explanation (REQ:label-on-the-wire), each on a page that exists today. A project with no shape file and no source file MUST look as it does today.

### Scan and the shape file

#### REQ: scan-writes-shape-file

Every `datatug scan` of a SQL database (`sqlite3`, `postgres`, `sqlserver`) MUST write the shape file `environments/<env>/catalogs/<db>/<db>.shape.json`, whole, in the form of the appendix, section 1. It MUST be built from the same in-memory result as the other files of that scan.

#### REQ: shape-is-descriptor-recordsets

Each item of the shape file's `recordsets` MUST carry exactly the seven required keys of a recordset in `ovdb-database/draft-1`, with the values the appendix gives: `name`, `kind`, `description` (`""`), `rowCount` (`null`), `columns`, `primaryKey`, `foreignKeys`. The file MUST NOT carry the descriptor's other keys, a time, a host, a path, a user name or a password.

#### REQ: shape-rescan-rules

A second scan of an unchanged database MUST write the same bytes. A scan MUST replace the shape file of the environment it scans and MUST NOT touch the shape file of another environment. The file is never edited by hand: DataTug MUST NOT offer to edit it, and a scan overwrites whatever is there.

#### REQ: scan-keeps-existing-files

A scan MUST keep writing every file that `datatug-cli`'s scan specification lists under "Project layout written by a scan", with the same bytes as today (Question 5).

#### REQ: scan-reads-foreign-keys

A scan MUST write into the shape file every foreign key its scanner read, with the columns it points at. The scanners of `sqlite3`, `postgres` and `sqlserver` all read them today (Question 4 for `postgres`). A scan MUST say in `scan.foreignKeys` how complete the list is: `read`; `partly-read` when its scanner left at least one foreign key out; `not-read` when its scanner reads none. For `partly-read` and `not-read`, DataTug MUST treat a foreign key that is not in the file as unknown, never as absent.

#### REQ: database-never-written

A scan MUST open the database read-only and create nothing in it. The importer, the marks, the match and the drift report MUST NOT open the database at all: they read project files.

### The source file

#### REQ: source-file-place

A database's source file is `dbmodels/<id>/<id>.source.yaml`, where `<id>` is the database's id. There is at most one per database, and it serves every environment. In this feature DataTug MUST NOT create, change, move or delete it; a scan MUST leave it byte for byte as it was, also when it removes a dropped table's folder beside it.

#### REQ: source-file-grammar

The source file MUST be read by the rules of the appendix, section 2: the five top-level keys and no other; the mapping lines of a publisher manifest, with the two differences its table marks "Differs"; names matched exactly, with their case. A column is written as a map with the one key `field` (Question 1).

#### REQ: source-file-names-model

The source file MUST name its model in one of two forms, and its meaning file in the form that goes with it: a registered model by `modelspec://` address with a commit of 40 hex characters, or a model that is a file of the project by path. The two forms MUST NOT be mixed. Facts from the first form carry the label `published`; facts from the second carry `declared`.

#### REQ: source-file-absent

A database with no source file MUST work: its tables, keys and foreign keys are shown from the shape file as `observed`; its entity files, if any, work as today; no meaning comes from a model; and the model line says whether a registered model matches. A database that has a source file and no shape file in the request's environment imports nothing, and the model line reads `Model: <the model> · not scanned in <env>`.

#### REQ: source-file-invalid

A source file that breaks a rule MUST be refused whole. DataTug MUST then show no meaning and no reference from a model for that database, MUST list each broken rule in words that name the key or the line, and MUST keep everything else working as with no source file. It MUST NOT use the part of the file that it could read.

### The importer

#### REQ: importer-inputs

The importer MUST be a pure function of four inputs: the shape, the mapping, the model and the links. It runs in the agent, in the process that calls the resolver. A change to any of the four MUST be visible on the next request, with no restart of the agent.

#### REQ: files-of-a-request

For a request that names `source` S in environment E, the catalog MUST be the one `api.ResolveSource` resolves S to in E, as today (datatug-cli@3a1c276 `pkg/api/resolver.go:35-49`). The shape file is that catalog's. The source file is the one of that catalog's database id. Imported mappings carry the database id as their `source`, as entity files do. The `--db` of the three commands is the catalog id, as in `datatug scan`.

#### REQ: importer-reads-both-spellings

The importer MUST read a model in ModelSpec JSON `1.0-draft-2` and in HCL with `record` and `field`, and, for a model read at a commit, the earlier spelling as well (specscore/modelspec@3c149d2 `spec/core-model.md:353-355`). It MUST read a meaning file in `meaning/draft-1` and in `meaning/draft-2`, treating each earlier word as the one that replaced it (meaninggraph/core@65510ae `FORMAT.md:368-379`).

#### REQ: importer-derives-links

The importer MUST apply MeaningGraph's rule for derived links itself, exactly as `FORMAT.md` section "Derived links" states it, and MUST NOT take derived links from another tool's index. It MUST first run MeaningGraph's check on the graph. For a graph that fails the check it MUST show no link at all, written or derived, and MUST report the failure.

#### REQ: importer-builds-entities

From the links and the mapping the importer MUST build entities and entity fields by the table and the eight rules of the appendix, section 3: the entity field each role gives; one meaning for one column, in the stated order; one column of a table for one entity field; no import for a key of several fields; the type table; no name patterns.

#### REQ: importer-builds-table-schemas

For each table the importer MUST fill the resolver's `TableSchema` from the shape file: the primary key, each foreign key with the columns it points at, and the foreign keys of other tables that point here. A reference that the model states and the database does not declare MUST be added to the same lists in memory, so that the resolver offers its related rows, and MUST NOT be written to any file. `GET semantic/columns` keeps reading a table's list of columns from the database, as today; only the entities and the table schemas change where they come from.

#### REQ: entity-files-win

When a project has both entity files and a source file, both feed the runtime index. A column that an entity file maps with a declared mapping MUST keep that mapping and its label `declared`. An imported mapping wins over an entity file's name pattern. An imported entity whose id equals an entity file's entity id MUST be left out whole.

#### REQ: imported-entities-are-not-entity-files

The importer MUST write no entity file and no other project file. Imported entities MUST NOT appear in `datatug entity list`, in `GET entities/all_entities`, or on the web app's entities pages, which list entity files as today.

#### REQ: not-imported-is-reported

Every link and every reference of the model that the rules leave out MUST be listed with its reason, in `GET openlayers/database` and in `datatug db status --format json`. Nothing is dropped silently.

#### REQ: fetched-files-are-pinned-and-rebuildable

A model or meaning file named by address MUST be read from its repository at the commit the source file names, at a path found in the registry's index. A graph that the meaning file pins, and that the check needs, MUST be read whole at its pin: every `*.meaning.yaml` in the root of that repository at that commit. Its list of files comes from the repository at that commit, not from a registry's index. Files read this way are kept in a cache outside the project. The cache is one folder, named by the setting `DATATUG_CACHE_DIR` or else under the user's cache directory. Deleting the cache MUST lose nothing. When a file cannot be read and is not cached, DataTug MUST say so on the model line and show no fact from it. DataTug MUST NOT read such a file at a branch head.

#### REQ: registry-index-source

DataTug MUST read each registry's index from one address: `index.json` at the head of `main` of `modelspec-org/registry` and of `meaninggraph/registry`. An index is the one thing it reads at a branch head. Two settings, `DATATUG_MODELSPEC_REGISTRY_INDEX` and `DATATUG_MEANINGGRAPH_REGISTRY_INDEX`, name another address or a local file. A test that needs a fixed index sets them; that is a setting, not an intercepted request.

### Labels

#### REQ: label-definitions

DataTug MUST use four labels with the meanings of the appendix, section 4: `observed`, `published`, `declared`, `inferred`. It MUST NOT use the proposal's fifth label, `suggested`, in this feature.

#### REQ: label-per-fact

Each fact MUST carry the label the appendix's second table of section 4 gives it. `observed` MUST NOT be the label of a meaning. A fact from a registered model MUST NOT be shown as `declared`. A key or reference that both the shape file and the model state MUST show `observed` and name the model's label beside it, as in `from the database · also published @ 3e7bb31`.

#### REQ: label-on-the-wire

`GET semantic/columns` MUST answer `published` for an imported mapping whose model was named by address, and `declared` for one whose model is a project file. The resolver does not know the difference: the route's handler sets the label from the importer's origins. A fact built from a column with an imported meaning MUST carry that column's label in `mapping`, and the explanation of an applicable query MUST name it, as in `CustomerId → customer.id (published) → parameter CustomerId`. A fact built from a column that an entity file maps is sent as it is today. The explanation is written inside `Applicable` from the value's `Provenance` (datatug-core@eb3075e `pkg/semantic/applicable.go:71`). The handler passes the fact's `mapping` there as it is; `pkg/semantic` gains no constant. The web app sends `mapping` with a grid selection, and accepts `published` in a release no later than the first `datatug-cli` release that sends it.

### What the screens show

#### REQ: table-list-mark

The database page's list of tables MUST show, for each table, a mark "n of m": m is the number of its columns in the shape file; n is the number of those that are described. A column is described when it has a meaning under any label, or is part of a key, or is a reference.

#### REQ: keys-tab-labels

On the table page's Keys tab, the primary key and each reference MUST show how it is known, in the words of the appendix, section 4: "from the database", "published @ " and the first seven characters of the commit, or "in the model". A reference that only the model states is listed with the foreign keys. The tab shows no key for any table today: the page never loads them (datatug-apps@098db23 `libs/datatug/main/src/lib/pages/signed-in/env-db-table/env-db-table.page.ts:299-301`, by its own comment; read, not run). In this feature it is filled from `GET openlayers/database`, for a database that has a shape file. For one that has none it stays as it is.

#### REQ: header-marker-published

A column whose meaning is `published` MUST show the header marker with the icon of a declared mapping and the tooltip `<entity>.<entity field> · published`. It MUST NOT show the style of an inferred mapping.

#### REQ: model-line

The database page MUST show one line that says which model the database follows, in one of the seven forms of AC:model-line-states, and under it whether the model and the database agree. When they agree it reads `Model and database agree.` When they differ it reads `Model and database differ in N places.` (`in 1 place.` for one), where N counts the findings of severity `error` and `warning`. Opening it lists every finding on the same page, one line each, as `datatug db drift` prints them (appendix, section 6).

#### REQ: related-rows-once

The response of `POST semantic/related` MUST NOT list one lookup twice. Where a foreign key and a shared entity field lead to the same database, table and column, one item is returned: the foreign key's. The removal is done by the route's handler, after the resolver.

### Match with registered models

#### REQ: match-definition

A registered model matches a database when, for every record type of the model, the shape file has a table of exactly that name, and that table has a column of exactly the name of each of the record type's fields. Views are not counted. A model with no record type never matches. The comparison reads the shape file and the index of `modelspec-org/registry` and nothing else. Printing a source file also reads the index of `meaninggraph/registry`.

#### REQ: match-is-local-and-read-only

The comparison MUST run on the person's machine. DataTug MUST send no name of a table, a column or a database to any network address: it downloads the two registries' indexes and compares locally. A match MUST NOT create or change any file, and MUST NOT make any statement of the matching model appear on a screen (Question 2).

#### REQ: match-shown

For a database with no source file, the model line MUST name every registered model that matches, sorted by address, with the counts "t of T tables, c of C columns": the model's record types and fields against all tables and columns of the database. `datatug db match --source-file <address>` MUST print a complete source file that names that model. The printed file's `meaning` block names the graph of the MeaningGraph registry's index whose `repository` is the model's `repository` and whose `model_files` list the model's file. With no such graph, or more than one, the file is printed without a `meaning` key and the command says so on its error output. When the index of models cannot be read, DataTug MUST say that, and everything else MUST work.

### Drift report

#### REQ: drift-compares

The drift report MUST compare the model, through the mapping, with the shape file of one environment: names, keys, references, types and whether a required field allows null. It MUST also report every broken rule of the source file. It reads no row and counts no value.

#### REQ: drift-findings

Each difference MUST be one finding of a kind and a severity from the appendix, section 7, with the table, column, record type and field it concerns. Types are compared by the appendix's tables of section 7, which read a `postgres` shape file by its own list. When the shape file says foreign keys were not read, or were partly read, no finding may say that a foreign key is missing.

#### REQ: drift-reports-only

The drift report MUST change nothing: no project file, no model, no database. It says what differs; a person decides which side to change.

### Commands and routes

#### REQ: cli-verbs

`datatug db status`, `datatug db match` and `datatug db drift` MUST exist with the flags, output and exit codes of the appendix, section 6. Each MUST work with the database unreachable.

#### REQ: read-only-routes

The agent MUST serve `GET openlayers/database`, `GET openlayers/match` and `GET openlayers/drift` with the bodies of the appendix, section 5. Each MUST apply the serving session's access policy: a table is listed and counted only when the principal may read it whole (`secureread.Executor.CanReadWholeCollection`, datatug-cli@3a1c276 `pkg/secureread/whole_collection.go:15`). A table that the policy hides, or narrows by rows or by fields, is in no list and no count, and a reference to it from a listed table is left out. Each MUST refuse every method but GET. The three commands read the project folder as the person who runs them and apply no policy; `--format json` prints the body the route returns when no policy narrows it.

### End-to-end gate

#### REQ: journey-test

One automated test MUST walk steps 1, 2 and 9 of the journey on Chinook, against a real agent and with no intercepted request, in four stages: (1) scan a SQLite Chinook into an empty folder and see the 11 tables with their marks; (2) with no action, see keys and foreign keys as known from the database, and the line that the registered model matches; (3) save the file that `datatug db match --source-file` prints, and commit; (4) as a second person, on a fresh clone with a second agent, open the same database and see the meanings, the references and how each is known. No stage may need a restart of the agent, a second scan, or a file edited by hand other than the one of stage 3 (Question 2). Stage 3 is one manual step more than the proposal's rule for this test allows (L2761). It falls away if Question 2 is answered B.

#### REQ: journey-project-has-no-entity-file

The project of that test MUST have no `entities` folder at any stage. The test MUST NOT change `datatug/datatug-demo-project`, and journeys J1 to J4 of `core-investigation-loop` MUST pass unmodified beside it (Question 6).

#### REQ: journey-cli-pin

The test MUST run against a `datatug-cli` release that contains this feature. Moving the pin of the web app's `journey` job to such a release is part of the work. The pin is v0.46.0 today (datatug-apps@098db23 `.github/workflows/ci.yml:168`) and the latest release is v0.68.0. Move it to the latest release in a change of its own first, so that a red J1 to J4 is not laid at this feature's door.

## What must not change

| What | How to test that it did not |
|---|---|
| The resolver | The existing tests of `datatug-core` `pkg/semantic` pass unmodified, and the implementing change touches no non-test file there (AC:fixed-interface). |
| The four routes | The existing contract tests of `datatug-cli` `pkg/server/endpoints` pass unmodified on their fixture, which has entity files and no source file, and the bodies recorded from `datatug-cli` v0.68.0 are equal after the change, the time in `observedAt` apart (AC:fixed-interface). |
| The screens of a project that has neither new file | Journeys J1 to J4 pass unmodified on `demo-project-1` (AC:fixed-interface). |
| The launch lane's work | No file of `datatug/datatug-demo-project` changes; no file of the web app's `public-data` route changes; the agent gains no write route (AC:journey-walk, AC:fixed-interface, AC:routes-are-read-only). |
| The database | The SQLite file has the same SHA-256 before and after the whole walk, and the three commands answer with the file moved away (AC:database-untouched). |
| The registries | DataTug only reads: an index, a file at a commit, or a repository's list of files at a commit; it opens no pull request and writes to no registry (AC:match-sends-nothing). |
| Entity files | `datatug entity list` prints the same before and after on `demo-project-1`, and nothing under `entities/` is created, changed or removed (AC:imported-entities-not-listed). |
| The source file | A scan leaves its bytes as they were (AC:scan-leaves-source-file). |

## Existing features

Line numbers of this repository's own pages, here and above, are as at datatug/datatug@33dd78c. "Depends on" means this specification relies on the feature as it is. "Extends" means it adds to what the feature covers; the feature's own page is not edited here and is listed for a later amendment by its owner. "Leaves" means no change either way.

| Feature | Status | Relation | What follows |
|---|---|---|---|
| `core-investigation-loop` | Stable | Depends on; extends | It owns the resolver, the four routes and the marker. When this is built it needs an amendment of `README.md:118-123` and `api-contract.md:72`, `:175`: `published` beside `declared` and `inferred` (Question 3), and a line that `POST semantic/related` lists one lookup once (REQ:related-rows-once). Its REQ:demo-knowledge-project and J1 to J4 stand (Question 6). Read in code and not run: a grid selection is sent with no `mapping` today (datatug-apps@098db23 `libs/datatug/semantic/src/contract/adapt.ts:90-108`), so its explanation would say "(inferred)" for a declared column. That gap is this feature's to avoid for imported columns (REQ:label-on-the-wire) and `core-investigation-loop`'s to close for the rest. |
| `semantic-layer` | Draft | Depends on | Entities and entity fields stay as defined. Its AC:employee-entity-gap-recorded stays true because imported entities are not listed (REQ:imported-entities-are-not-entity-files). |
| `related-data-navigation` | Draft | Depends on; extends | Both directions of a foreign key now come from the shape file. It names relationship labels of its own, `declared`, `fk`, `semantic`, `discovered` (`:130-134`); `fk` and this page's `observed` are the same fact. Listed for its owner to settle on one word. |
| `meaning-inference` | Draft | Leaves | It owns `inferred` and suggestions. This feature infers nothing and suggests nothing. |
| `provenance` | Draft | Depends on; extends | Its REQ:mapping-provenance names "declared, inferred, corrected". For a `published` mapping, "who or what produced it" is the model's address and commit. Listed for a later amendment. |
| `heterogeneous-sources` | Draft | Depends on | It owns scanning. No connection secret enters either new file. |
| `infrastructure-inventory` | Draft | Depends on | `datatug scan` stays the source of truth for what an environment holds. |
| `git-native-projects` | Implementing | Depends on | Both new files are readable project files with no secret. This feature adds no write through its mutation contract. |
| `dalgo-project-store` | Approved | Depends on; extends | Two new file kinds, suffixes `shape` and `source`, beside existing records. Its Open Question 3 asks whether the suffixes `columns`, `objects` and `refs`, which "have no collection in `ProjectStore`", "get a home, or are they retired" (`:1067-1068`); `shape` and `source` join that question. Its founder decision of 2026-09-17, "writes should be through dalgo" (`:530-531`), will bind the scan's write of the shape file once that feature is built. This specification does not say how the file is written. The source file is YAML and is not written by DataTug in this feature. |
| `public-data-fabric` | Draft | Leaves | Another lane. Its ledger calls canonical ingestion a later prerequisite. This feature does not claim to complete it, and edits nothing under it. Its TypeScript reader in the web app is not used and not changed. |
| `query-workspace` | Approved | Leaves | It lists "foreign key, declared semantic mapping, inferred" for join suggestions (`:122-124`). A fourth word there is its owner's to add. |
| `checks`, `compare` | Approved | Leaves | The drift report is not a `Check`, and the match is not a Compare. Neither namespace is used. |
| `agent-context` | Draft | Leaves | Agents get imported meanings through the same resolver and routes. |
| `pricing-and-billing` | Draft | Leaves | The three routes and commands read metadata and run no data query. |
| `datatug-cli` `spec/features/cli/scan` | Implementing | Extends | One more file per scan, which holds the foreign keys the scan already reads (Questions 4 and 5). To be amended there when built. |
| `datatug-core` `spec/features/model-foreign-keys-file` | Draft | Overlaps | Not used (choice 20), nor is the reader and writer of that file that `datatug-core` has had since v0.43.0. Its owner should mark it superseded if this is approved. |

## Dependencies

- core-investigation-loop
- semantic-layer
- related-data-navigation
- provenance
- heterogeneous-sources
- git-native-projects
- dalgo-project-store

## Blocked Clauses

- **REQ:importer-reads-both-spellings and REQ:importer-derives-links** need a Go reader of `meaning/draft-2` and of derived links. The released library `github.com/meaninggraph/cli/pkg/meaning` v0.3.0 reads `meaning/draft-1` only and derives nothing (meaninggraph/cli@e5d4d8f; its one mention of `draft-2` is a test that expects it refused, `pkg/meaning/schema_test.go:62`). Either that library gains both first, or the importer carries its own reader and passes the conformance cases of `meaninggraph/core`. Chinook's files are `meaning/draft-1` with the earlier role names, which v0.3.0 reads and checks, so the first half of AC:chinook-links needs only the derivation. Its second half, and any file that writes `instances` or `reference`, needs the new reader: v0.3.0's schema lists the roles `entity`, `identifier`, `display-name`, `foreign-key` and `value` only (`pkg/meaning/meaning.schema.json:94`). The conformance cases are tests inside a Node script (meaninggraph/core@65510ae `scripts/test-conformance.mjs`, cases D-01 onward), not files another language can load, so they would be written again in Go.
- Nothing blocks **REQ:scan-reads-foreign-keys**. `dal-go/dalgo` v0.93.0 and `dalgo2postgres` v0.9.1, which `datatug-cli` already requires, report a foreign key's target table and columns.

## Acceptance Criteria

Unless a criterion says otherwise, "the Chinook project" is a folder made by `datatug scan --directory ./shop --driver sqlite3 --path chinook.sqlite --db chinook --env local`, and "Example A" is the source file of that name in the appendix, saved as `dbmodels/chinook/chinook.source.yaml`. `chinook.sqlite` is the file the journey job pins (datatug-apps@098db23 `.github/workflows/ci.yml:232-234`, SHA-256 `f82efedb6c5c40734609e168bc5be5616a2eca6b90ed0048451a8674625e03a3`). It is kept outside the project folder, so that the catalog file records a path a clone on the same machine can open. Where a criterion reads a registry's index, the two settings of REQ:registry-index-source name the index files of `modelspec-org/registry@a7325ea` and `meaninggraph/registry@49f472f`.

### AC: fixed-interface (verifies REQ:resolver-unchanged, REQ:routes-unchanged, REQ:no-new-screen)

**Given** `demo-project-1` as it is at datatug-demo-project@deeee32, with eight entity files, no shape file and no source file, and an agent and web app built with this feature
**When** the existing tests of `pkg/semantic`, the existing contract tests of the four routes and the journey tests J1 to J4 run unmodified
**Then** all pass; the bodies of the four routes for the requests of `semantic_contract_test.go`, recorded from `datatug-cli` v0.68.0, are equal after the change, apart from `observedAt` and from a lookup that was listed twice; the implementing change alters no non-test file of `pkg/semantic` and no file under the web app's `libs/datatug/main/src/lib/queries/public-data/`; the web app's router has no new address; and the agent's routes have grown by exactly three, all GET.

### AC: scan-writes-shape (verifies REQ:scan-writes-shape-file, REQ:shape-is-descriptor-recordsets)

**Given** a SQLite Chinook file and an empty folder `./shop`
**When** `datatug scan --directory ./shop --driver sqlite3 --path chinook.sqlite --db chinook --env local` runs
**Then** `./shop/environments/local/catalogs/chinook/chinook.shape.json` exists with `"format": "datatug-shape/draft-1"`, 11 recordsets, 64 columns and 11 foreign keys in all; its `Invoice` has `"rowCount": null`, `"primaryKey": [{"column": "InvoiceId", "position": 1}]` and one foreign key `{"column": "CustomerId", "table": "Customer", "referencedColumn": "CustomerId", "constraint": 0, "position": 0}`; and the file contains no path, host or time.

### AC: shape-rescan (verifies REQ:shape-rescan-rules, REQ:scan-keeps-existing-files)

**Given** the Chinook project
**When** the same scan runs again, and then a scan with `--env prod` of the same file
**Then** `environments/local/catalogs/chinook/chinook.shape.json` has the same bytes after each, `dbmodels/chinook/main/tables/Invoice/main.Invoice.columns.json` exists with the bytes today's scan writes, and a second shape file exists under `environments/prod/`.

### AC: foreign-keys-read-or-not (verifies REQ:scan-reads-foreign-keys)

**Given** a PostgreSQL database loaded with Chinook; a SQLite file in which one table has a foreign key to a table that is not in the file; and a third scan run with a test scanner that reads no foreign keys
**When** each is scanned
**Then** the first shape file says `"foreignKeys": "read"` and its `Invoice` has the foreign key `CustomerId` to `Customer`; the second says `"foreignKeys": "partly-read"`; the third says `"foreignKeys": "not-read"` and its Keys tab shows no reference as "from the database"; no drift report of the second or the third has a finding `reference-not-declared`; and the drift report of the PostgreSQL Chinook with Example A has no finding `type-differs`.

### AC: database-untouched (verifies REQ:database-never-written)

**Given** the SHA-256 of `chinook.sqlite` taken before the scan
**When** the scan, `datatug db status`, `datatug db match`, `datatug db drift` and the four routes have all been used, and then the SQLite file is moved away and the three `db` commands run again
**Then** the file's SHA-256 is unchanged, and the three commands answer as before.

### AC: scan-leaves-source-file (verifies REQ:source-file-place)

**Given** the Chinook project with Example A
**When** the database loses its table `Playlist` and is scanned again
**Then** the folder of `Playlist` under `dbmodels/chinook/` is gone as today, and `dbmodels/chinook/chinook.source.yaml` has the same bytes.

### AC: source-file-refusals (verifies REQ:source-file-grammar, REQ:source-file-invalid)

**Given** the Chinook project with Example A changed in one of three ways: the line `modules: {}` added; `?ref=3e7bb316…` removed from the address; or an item `- name: Invoice` with `columns: { Total: Amount }`
**When** `GET openlayers/database` and `GET semantic/columns` for `Invoice` are called, and `datatug db drift` runs
**Then** `model.state` is `unreadable` with one problem that names the key (`modules`), the address, or `column "Total" must be a map with field:`; `semantic/columns` returns no column; the Keys tab still shows `CustomerId ⇒ Customer`, from the database; and `datatug db drift` exits 1.

### AC: source-file-absent (verifies REQ:source-file-absent)

**Given** the Chinook project with no source file
**When** `GET openlayers/database` and `GET semantic/columns` for `Invoice` are called
**Then** `model.state` is `none`; `Invoice` has `key` `["InvoiceId"]` and one reference to `Customer`, both `known: "observed"`; and `semantic/columns` returns no column.

### AC: project-model-is-declared (verifies REQ:source-file-names-model, REQ:label-per-fact)

**Given** a project with a scanned table `payments` that has a column `order_ref`, the files `model/payments.modelspec.hcl` (a record type `Payment` with a string field `OrderId`) and `model/payments.meaning.yaml` (a concept `payment` of kind entity bound to `Payment`, and a concept `order-reference` of kind property, `of: payment`, bound to `Payment.OrderId` with role `value`), and the source file of Example B
**When** `GET semantic/columns` for `payments` and `GET openlayers/database` are called
**Then** the first returns `{"column": "order_ref", "entity": "payment", "field": "order-reference", "provenance": "declared"}`, and the second has `model.known` `"declared"` and `model.files` naming the two paths.

### AC: chinook-links (verifies REQ:importer-derives-links, REQ:importer-reads-both-spellings)

**Given** Chinook's model and meaning file at commit `3e7bb316d8a438eb7e5930067997a4e4543a275c`, and a copy of that meaning file rewritten in `meaning/draft-2` with `kind: property`, `field:`, `role: instances` and `role: reference`
**When** the importer computes the links of each
**Then** each gives the same 48 links, 36 written and 12 derived, equal to what the reference `deriveLinks` of meaninggraph/core@65510ae gives; among them `customer` on `Invoice.CustomerId` with role `reference`, written, with its note, and `customer` on `Customer.CustomerId` with role `identifier`, derived.

### AC: failing-graph-shows-nothing (verifies REQ:importer-derives-links)

**Given** a project whose meaning file binds a concept to a field that its model does not have
**When** `GET semantic/columns` is called for any table of that database
**Then** it returns no column from that meaning file, and `GET openlayers/database` reports the check's finding under `model.problems` with `model.state` `unreadable`.

### AC: chinook-index (verifies REQ:importer-builds-entities, REQ:not-imported-is-reported, REQ:importer-inputs)

**Given** the Chinook project with a running agent and no source file
**When** Example A is saved while the agent runs, and `GET semantic/columns` is called for `Customer` and for `Employee`, and `GET openlayers/database`
**Then**, with no restart, `Customer` returns exactly `CustomerId` as `customer.id`, `Country` as `customer.customer-country` and `SupportRepId` as `employee.id`; `Employee` returns exactly `EmployeeId` as `employee.id` and `ReportsTo` as `employee.manager`; and `notImported` has exactly three entries: `music-sales` on `Invoice.Total` (`measure`), `support-rep` on `Customer.SupportRepId` (`second-meaning`) and `employee` on `Employee.ReportsTo` (`entity-field-taken`).

### AC: related-rows-from-shape (verifies REQ:importer-builds-table-schemas, REQ:related-rows-once)

**Given** the Chinook project with Example A
**When** `POST semantic/related` is called for `customer.id = 5` selected on `Customer.CustomerId`, and again selected on `Invoice.CustomerId`
**Then** the first returns exactly one item, collection `Invoice`, count 7; the second returns exactly one item, collection `Customer`, count 1.

### AC: model-reference-without-foreign-key (verifies REQ:importer-builds-table-schemas, REQ:label-per-fact, REQ:keys-tab-labels)

**Given** a SQLite copy of Chinook made with no foreign key constraints, scanned, with Example A
**When** the Keys tab of `Invoice` is opened, `POST semantic/related` is called for `customer.id = 5` on `Customer.CustomerId`, and `datatug db drift` runs
**Then** the tab lists `CustomerId ⇒ Customer` with "published @ 3e7bb31"; the related response has one item, `Invoice`, count 7; the drift report has 11 findings `reference-not-declared` of severity `note` and exits 0; and no project file has changed.

### AC: entity-file-wins (verifies REQ:entity-files-win, REQ:files-of-a-request)

**Given** a copy of `demo-project-1` in which `datatug scan --directory <the copy> --driver sqlite3 --path <the Chinook file> --db chinook --env local` has run, and Example A saved as `dbmodels/chinook/chinook.source.yaml`
**When** `GET semantic/columns` is called with `source=chinook` for `Customer`, and then for `Invoice`
**Then** for `Customer`, `CustomerId` is `{"entity": "Customer", "field": "ID", "provenance": "declared"}` as today, `SupportRepId` is `{"entity": "employee", "field": "id", "provenance": "published"}`, and `notImported` lists `customer` on `Customer.CustomerId` with the reason `entity-file-maps-column`; for `Invoice`, `CustomerId` is `{"entity": "customer", "field": "id", "provenance": "published"}`, where today it is `Customer.ID`, `inferred`.

### AC: imported-entities-not-listed (verifies REQ:imported-entities-are-not-entity-files)

**Given** the Chinook project with Example A, committed, and `demo-project-1`
**When** the four routes and the three new routes have been called on the first, and `datatug entity list` runs on each
**Then** it prints no entity for the first and the same eight as today for the second; the first project has no `entities` folder; and `git status` reports no change in either.

### AC: fetched-files (verifies REQ:fetched-files-are-pinned-and-rebuildable)

**Given** the Chinook project with Example A, and `DATATUG_CACHE_DIR` naming an empty folder
**When** `GET semantic/columns` for `Invoice` is called online, again after that folder is deleted, again offline with the cache warm, and again offline with the folder deleted
**Then** the first three return the same five columns, and the agent's log shows `model/chinook.modelspec.json` and `model/chinook.meaning.yaml` read from `github.com/demo-db/chinook` at `3e7bb316d8a438eb7e5930067997a4e4543a275c`; the fourth returns no column and the model line says the model's files could not be read; and the project folder holds no new file after any of them.

### AC: labels-on-columns (verifies REQ:label-definitions, REQ:label-on-the-wire, REQ:header-marker-published)

**Given** the Chinook project with Example A
**When** `GET semantic/columns` is called for `Invoice` and the table page is opened
**Then** the response has exactly five items, `InvoiceId` as `invoice.id`, `CustomerId` as `customer.id`, `InvoiceDate` as `invoice.invoice-date`, `BillingCountry` as `invoice.billing-country` and `Total` as `invoice.invoice-total`, each with `"provenance": "published"`; no response of any route contains `observed` as a `provenance` or the word `suggested`; and the header of `CustomerId` shows the declared icon, without the inferred style, with the tooltip `customer.id · published`.

### AC: explanation-names-the-label (verifies REQ:label-on-the-wire)

**Given** the Chinook project with Example A and a saved query whose required parameter `CustomerId` is tagged `{"entity": "customer", "field": "id"}`
**When** the cell `CustomerId = 5` of `Invoice` is selected
**Then** the panel lists that query as applicable with the explanation `CustomerId → customer.id (published) → parameter CustomerId`.

### AC: marks (verifies REQ:table-list-mark)

**Given** the Chinook project, first with no source file and then with Example A
**When** the database page is opened each time
**Then** the marks read `Invoice 2 of 9` and `Customer 2 of 13` the first time, 21 of 64 over all tables, and `Invoice 5 of 9` and `Customer 3 of 13` the second time, 35 of 64 over all, as the appendix's table gives for each of the 11 tables.

### AC: keys-tab (verifies REQ:keys-tab-labels)

**Given** the Chinook project with Example A
**When** the Keys tab of `Invoice` is opened
**Then** it shows the primary key `InvoiceId` and the foreign key `CustomerId ⇒ Customer`, each with `from the database · also published @ 3e7bb31`.

### AC: model-line-states (verifies REQ:model-line, REQ:match-shown, REQ:source-file-absent)

**Given** seven projects: Chinook with no source file; a database no registered model matches; Chinook with the registry unreachable; Chinook with Example A; the payments project of AC:project-model-is-declared; Chinook with a source file that has the key `modules`; Chinook with Example A and its shape file removed
**When** the database page of each is opened
**Then** the line reads, in that order: `Registered model chinook matches this database: 11 of 11 tables, 64 of 64 columns.`; `No model. No registered model matches this database.`; `No model. The registry of models could not be read.`; `Model: modelspec://github.com/demo-db/chinook/chinook @ 3e7bb31 · published` followed by `Model and database agree.`; `Model: model/payments.modelspec.hcl · declared`; `Model: the source file cannot be used (1 problem).`; `Model: modelspec://github.com/demo-db/chinook/chinook @ 3e7bb31 · not scanned in local`

### AC: match-chinook (verifies REQ:match-definition, REQ:match-shown, REQ:registry-index-source)

**Given** the Chinook project with no source file, and the two settings of REQ:registry-index-source naming the index files of `modelspec-org/registry@a7325ea`, which lists nine models, and of `meaninggraph/registry@49f472f`
**When** `datatug db match` runs, and then `datatug db match --source-file modelspec://github.com/demo-db/chinook/chinook`
**Then** the first prints one match, `chinook`, "11 of 11 tables, 64 of 64 columns", and none of the other eight; the second prints exactly the text of Example A; both exit 0.

### AC: match-with-extra-and-missing (verifies REQ:match-definition)

**Given** two copies of Chinook, scanned: one with an added table `audit_log` of three columns, one with the column `Invoice.Total` dropped
**When** `datatug db match` runs on each
**Then** the first prints `chinook`, "11 of 12 tables, 64 of 67 columns"; the second prints that no registered model matches.

### AC: match-sends-nothing (verifies REQ:match-is-local-and-read-only)

**Given** the Chinook project with no source file, the two settings of REQ:registry-index-source unset, and every outgoing request of the agent recorded
**When** the database page is opened and `datatug db match` runs
**Then** every recorded request is a GET of a registry's `index.json` or of a file at a commit; no URL, header or body contains `Invoice`, `CustomerId` or `chinook.sqlite`; no project file was created or changed; and `GET semantic/columns` for `Invoice` still returns no column.

### AC: match-offline (verifies REQ:match-shown)

**Given** the Chinook project with no source file, no network, and the two settings of REQ:registry-index-source unset
**When** the database page is opened and `datatug db match` runs
**Then** the model line reads `No model. The registry of models could not be read.`, the marks and the Keys tab are as in AC:marks and AC:source-file-absent, and the command says the same and exits 0.

### AC: drift-chinook-agrees (verifies REQ:drift-compares, REQ:drift-findings)

**Given** the Chinook project with Example A
**When** `datatug db drift` runs
**Then** it prints `Model and database agree: 11 tables, 64 columns, 11 keys, 11 references.`, lists no finding and exits 0.

### AC: drift-finds-differences (verifies REQ:drift-compares, REQ:drift-findings, REQ:drift-reports-only, REQ:model-line)

**Given** the Chinook project with Example A, after the database lost the column `Invoice.Total` and the primary key of `Genre`, and was scanned again
**When** `datatug db drift --format json` runs
**Then** it lists `field-without-column` for `Invoice.Total` and `key-differs` for `Genre`, both of severity `error`, exits 1, and `git status` shows no change other than the scan's; and the database page reads `Model and database differ in 2 places.`, and opening it lists both findings.

### AC: routes-are-read-only (verifies REQ:read-only-routes, REQ:cli-verbs)

**Given** the Chinook project with Example A and a policy that lets the role `admin` read every table and the role `support` read `Customer` and `Invoice` whole and nothing else, served by two agents: one started `--as boss --role admin`, one `--as agent --role support`
**When** the three `openlayers` routes are called with GET on each agent and with POST on the first, and the three commands run with `--format json`
**Then** no response of the second agent names or counts a table other than `Customer` and `Invoice`; each POST is refused; and each command prints the body the first agent returns.

### AC: journey-walk (verifies REQ:journey-test, REQ:journey-project-has-no-entity-file, REQ:journey-cli-pin, REQ:registry-index-source)

**Given** an empty folder, the SQLite Chinook file that the journey job pins, kept outside that folder, a `datatug-cli` release that contains this feature, and the two settings of REQ:registry-index-source naming the index files of `modelspec-org/registry@a7325ea` and `meaninggraph/registry@49f472f`
**When** the test runs its four stages in order with no agent restart and no second scan: scan and open the database page; open the Keys tab of `Invoice`; save the output of `datatug db match --source-file modelspec://github.com/demo-db/chinook/chinook` as the source file and commit; clone the project, start a second agent on the clone and open `Invoice` and `Customer` in a fresh browser
**Then** stage 1 shows 11 tables with `Invoice 2 of 9`; stage 2 shows `CustomerId ⇒ Customer`, from the database, and the line that the registered model `chinook` matches; stage 4 shows the tooltip `customer.id · published`, `Invoice 5 of 9`, the related rows `Invoice` (7) for `CustomerId = 5` of `Customer`, and the model line with the first seven characters of the commit that index lists for `chinook`; the project has no `entities` folder at any stage; the second agent's log shows no scan; the SQLite file's SHA-256 is unchanged; `datatug/datatug-demo-project` has no change; and J1 to J4 pass beside it.

## Rehearse Integration

Every criterion has a surface that a test can observe: a file's bytes, a route's body, a command's output and exit code, a line on a page. This Feature spans three repositories, and each criterion is exercised against a real implementation: the importer and the type tables in `datatug-core`, the scan, the routes and the commands in `datatug-cli`, the pages and the end-to-end walk in `datatug-apps`. Stub scaffolding is therefore deferred to those implementations, as `openvaultdb-project-store` and `ai-query-builder` do. AC:fixed-interface and AC:journey-walk run existing suites unmodified and need no stub of their own.

## Not Doing / Out of Scope

- Anything of the proposal's Phases 5 to 8: suggestions, the label `suggested`, accept and reject, a ledger or an answers file, value counts, inference, DataTug writing a model or a meaning file, publishing, pull requests to a registry.
- Retiring entity files or DataTug's entity types. Proposal D8 keeps that for a later, separate decision. No entity file is removed, and the `datatug entity` commands and the entities pages stay as they are.
- Any change to `demo-project-1`.
- A registry, an index or a service of DataTug's own. DataTug reads the three public registries' files and keeps a cache that can be deleted.
- Meanings in the hosted web app for a project opened from GitHub. Such a page has no agent today (datatug-apps@098db23 `apps/datatug-app/src/app/agent-base-url.ts:13-19`), and this feature does not give it one.
- inGitDB collections and HTTP sources. They get no shape file here.
- Reading the database for the drift report: rows, value counts, "every reference's values are present in its target".

## Open Questions

The questions that change what is built are in [Questions for the owner](#questions-for-the-owner). The points below are left to a later phase by this specification's author. The proposal assigns only some of them to a phase. They are listed so that nobody assumes an answer.

For later phases:

1. Adopting a matching model with one confirmation, if Question 2 is answered A. The proposal's Phase 5 writes the source file.
2. What an unadopted matching model's statements look like before a person keeps them. The proposal answers with `suggested`, which is its Phase 5.
3. A second meaning on one column. `Customer.SupportRepId` is both "an employee" and "the support rep of a customer"; this feature shows the first and reports the second.
4. Measures. `music-sales` on `Invoice.Total` is not imported.
5. A reference to a record type whose key has several fields, and an entity field for such a key. AdventureWorks has two such references.
6. A path into a component, such as `field: Amount.Minor`. The grammar accepts it; nothing resolves it yet.
7. References between two models, the pairing of two databases and the source file's `modules` line. The proposal's Phase 7.
8. A database listed in the OpenVaultDB Directory and opened by its public address. Its mapping is in its publisher's manifest, and the proposal says a project cannot override it. This feature reads the source file only.
9. Engines that fold the case of names. Matching is exact here; a folded name needs a `columns` line each, or a later rule.
10. Retiring the columns files of a scan, if Question 5 is answered A, and the Draft `refs.json`.
11. How DataTug itself writes the source file through the project store, which fixes JSON for project items.
12. Searching MeaningGraph for a name when no model matches (proposal L1569), and value sets.
13. A shape for a private database inside OpenVaultDB's own descriptor format, and one shared generator for it.

Not about what is built:

14. When building may start. The proposal places this phase "after the foundation, and after the launch" (L2796). This specification does not decide whether building may start before the launch of 2026-10-25, once it is approved. When it is built, it changes three places that the launch also runs on: the scan's writer (`datatug-cli` `pkg/api/scan_layout.go`), the closed lists of the Stable contract (`datatug-core` `pkg/apicontract`, `datatug-apps` `libs/datatug/semantic`) and the pin of the journey job (`datatug-apps` `.github/workflows/ci.yml`). The author's suggestion is to start after the launch, as the proposal places it.

---
*This document follows the https://specscore.md/feature-specification*
