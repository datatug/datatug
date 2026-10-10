# Open-Layer Reading: file and wire formats

This is the formats appendix of the [owning Feature](README.md). It is not a separate
feature. It holds the exact shapes that the requirements there point at: two project
files, the rule that turns links into DataTug's runtime index, the labels, three
read-only routes, three commands and the drift findings.

Nothing here is built. Everything here is the specification author's proposal unless a
line says it is the owner's. The words are those of the Feature's section
[Words used here](README.md#words-used-here).

Every example uses Chinook as registered: `demo-db/chinook` at commit
`3e7bb316d8a438eb7e5930067997a4e4543a275c`, written `3e7bb31` below.

## 1. The shape file

| | |
|---|---|
| Path | `environments/<env>/catalogs/<db>/<db>.shape.json`, beside `<db>.db.json` |
| Written by | every `datatug scan` of that database in that environment, whole, each time |
| Edited by a person | never; the next scan replaces it |
| Read by | the importer, the marks, the match and the drift report |
| Committed | yes, like the other files a scan writes; DataTug stages it and the person commits |

The file for Chinook, cut to one table and two of its nine columns:

```json
{
  "format": "datatug-shape/draft-1",
  "id": "chinook",
  "scan": { "driver": "sqlite3", "foreignKeys": "read" },
  "recordsets": [
    {
      "name": "Invoice",
      "kind": "table",
      "description": "",
      "rowCount": null,
      "columns": [
        { "name": "InvoiceId", "type": "INTEGER", "nullable": false,
          "primaryKey": true, "primaryKeyPosition": 1, "defaultValue": null },
        { "name": "CustomerId", "type": "INTEGER", "nullable": false,
          "primaryKey": false, "primaryKeyPosition": null, "defaultValue": null }
      ],
      "primaryKey": [ { "column": "InvoiceId", "position": 1 } ],
      "foreignKeys": [
        { "column": "CustomerId", "table": "Customer",
          "referencedColumn": "CustomerId", "constraint": 0, "position": 0 }
      ]
    }
  ]
}
```

| Key | Value | Rule |
|---|---|---|
| `format` | `datatug-shape/draft-1` | Any other value: the file is not read, and the model line says so. |
| `id` | the catalog's id, `<db>` | Must equal the folder name. |
| `scan.driver` | `sqlite3`, `postgres` or `sqlserver` | The driver the scan used. |
| `scan.foreignKeys` | `read`, `partly-read` or `not-read` | `partly-read`: the scanner read foreign keys and left at least one out. `not-read`: the scanner of this engine does not read foreign keys. In both cases a foreign key that is not in the list is unknown, and an empty `foreignKeys` list says nothing about the database. |
| `recordsets[]` | one item per table and per view | Sorted by `name`, comparing bytes. |

One item of `recordsets` has exactly the seven required keys of a recordset in OpenVaultDB's
database descriptor `ovdb-database/draft-1`
(demo-db/chinook@c7535f4 `schemas/ovdb-database-draft-1.schema.json`), with the same
meaning and the same key names inside `columns`, `primaryKey` and `foreignKeys`:

| Key | What the scan writes |
|---|---|
| `name` | The table's own name. A table outside the engine's default schema (`main`, `public`, `dbo`) is written `<schema>.<table>`. |
| `kind` | `table` or `view`. |
| `description` | Always `""`. A scan has no description to give. |
| `rowCount` | Always `null`. A scan counts no rows. |
| `columns[]` | `name`; `type`, the type text the scan holds: the engine's own for `sqlite3` and `sqlserver`, and DALgo's type name for `postgres` (`int`, `string`, `bool`, `float`, `decimal`, `time`, `bytes`, `null`, `unknown`), which is all that scan has today; `nullable`; `primaryKey`; `primaryKeyPosition`, 1 or more, or `null`; `defaultValue`, the engine's text or `null`. In the engine's column order. |
| `primaryKey[]` | `column`, `position`. Empty when the table has no primary key. |
| `foreignKeys[]` | `column`, `table`, `referencedColumn`, `constraint`, `position`. A foreign key of several columns is several items with one `constraint` number, in `position` order. Constraints are numbered from 0 in the order the engine reports them. When the engine names no target column, `referencedColumn` is the target's primary key column at the same position. |

The descriptor's optional key `modelEntity` is never written: the mapping lives in the
source file. No other key is written. In particular the file holds no time, no host, no
path, no user name and no password, so that two scans of an unchanged database write the
same bytes.

## 2. The source file

| | |
|---|---|
| Path | `dbmodels/<id>/<id>.source.yaml`, beside `<id>.dbmodel.json` |
| Written by | a person. In this feature DataTug never creates, changes or deletes it. |
| Read by | the importer and the drift report |
| One per | database: the id that the four `semantic/*` and `queries/applicable` routes call `source` |

`<id>` is the database's id. In today's files it is the `dbModel` that an environment's
`<db>.db.json` names, for example `chinook` for the catalogs `chinook-local` and
`chinook-prod` (datatug-demo-project@deeee32
`demo-project-1/environments/local/catalogs/chinook-local/chinook-local.db.json`), or the
catalog's own id when its file names none. One source file therefore serves every
environment of the database. Each catalog has its own shape file. Which two files a
request reads is REQ:files-of-a-request of the Feature.

### Example A: a database that follows a registered model

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
  - Artist
  - Customer
  - Employee
  - Genre
  - Invoice
  - InvoiceLine
  - MediaType
  - Playlist
  - PlaylistTrack
  - Track
```

### Example B: an internal database with its own model

The owner's approved example (see the Feature, "Decided by the owner, and where"), with
the three additions this specification makes: the `format` line, the `meaning` key, and
the column written as a map (Question 1 of the Feature).

```yaml
format: datatug-source/draft-1
id: payments
model:
  hcl: model/payments.modelspec.hcl
meaning:
  file: model/payments.meaning.yaml
recordsets:
  - name: payments
    record_type: Payment
    columns:
      order_ref:
        field: OrderId
```

### Keys

| Key | Required | Value | Rule |
|---|---|---|---|
| `format` | yes | `datatug-source/draft-1` | Any other value is refused. |
| `id` | yes | the database's id | Must equal the folder name `<id>`. |
| `model` | yes | a map in one of two forms | See "How the model and the meaning file are named". |
| `meaning` | no | a map in the form that goes with `model` | Left out: references and keys come from the model, and no column gets a meaning. |
| `recordsets` | yes | a list | Each item is a table's name, or a map. See "The mapping lines". |

Any other top-level key is refused, `modules` included: pairing two databases belongs to a
later phase, and a file that carries it must not be half-read by this one.

### How the model and the meaning file are named

| Form | `model` | `meaning` | Where the files are read | Label |
|---|---|---|---|---|
| Registered | `address:` a `modelspec://` address with `?ref=` and a commit of 40 hex characters | `graph: { address: }` a `meaning://` address of the same repository, with no `?ref=` | The model file that the ModelSpec registry's index lists for the address (its JSON file when it lists one, else its HCL file), and the meaning files that the MeaningGraph registry's index lists for the graph, each read from that repository at the commit in `?ref=` | `published` |
| In the project | `hcl:` and, or, `modelspec:` a path inside the project folder | `file:` a path inside the project folder | The project's own files | `declared` |

Rules for both forms:

- The two forms are not mixed: an address for one and a path for the other is refused.
  MeaningGraph allows a binding to a model of the same repository only
  (meaninggraph/core@65510ae `FORMAT.md:290`).
- An address with no `?ref=`, or a `?ref=` that is not 40 hex characters, is refused.
- A path is relative to the project folder, uses `/`, and may not leave the folder.
- With both `hcl:` and `modelspec:` given, the JSON file is read. The HCL file is the
  one a person edits.
- The module that the model names must be a key of the meaning file's `models` map. If it
  is not, no meaning is imported and the problem is reported.

### The mapping lines

These are the publisher manifest's lines, read by the publisher manifest's rules
(openvaultdb/openvaultdb@37d2389 `spec/decisions/0012-…`, the points named N5 to N13
there). Only the differences are DataTug's.

| Line | Means | Same as the manifest? |
|---|---|---|
| `- Invoice` | Table `Invoice` holds rows of record type `Invoice`. Each column is the field of the same name. | Yes. |
| `- name: payments`<br>`  record_type: Payment` | Table `payments` holds rows of record type `Payment`. | Yes. `record_type` may be left out when it equals `name`. |
| `columns:`<br>`  order_ref:`<br>`    field: OrderId` | Column `order_ref` is the field `OrderId`. A column that is not listed is the field of its own name. | Yes: a map with the one key `field`. A bare string in its place is refused, as in a manifest (Question 1). |
| a table that no item names | The table is outside the model. Its keys and foreign keys are still shown as known from the database. | Differs: a manifest lists every recordset it publishes. |
| two items with one `record_type` | Allowed. Both tables hold rows of that record type. | Differs: a manifest refuses it. |
| two columns of one table with one `field` | Refused. | Yes. |
| an item's key other than `name`, `record_type`, `columns`; a column's key other than `field` | Refused. | Yes. |
| names | Matched exactly, with their case. No folding. | Yes. |
| `field: Contact.Email`, a path into a component | Accepted by the grammar. The column gets no meaning in this feature and is reported as not imported (reason `component-path`). | Yes: the manifest's readers do not resolve such a path either. |

## 3. From links to DataTug's runtime index

The importer is one pure function. It reads no file and no network itself: its caller
loads the four inputs.

| Input | Comes from |
|---|---|
| the shape | the shape file of the request's catalog |
| the mapping | the source file's `recordsets` |
| the model | the file that the source file's `model` names |
| the links | the meaning files that the source file's `meaning` names, checked, then put through MeaningGraph's rule for derived links (meaninggraph/core@65510ae `FORMAT.md:307-366`) |

| Output | Type, as it is today in datatug-core@eb3075e |
|---|---|
| imported entities | `[]*datatug.Entity` (`pkg/datatug/entities.go:41-47`, `:107-117`) |
| table schemas | `map[semantic.SchemaKey]semantic.TableSchema` (`pkg/semantic/related.go:28-40`) |
| origins | for each imported mapping and each key or reference: the label it carries |
| not imported | each link or reference that the rules below leave out, with the reason |

### Which entity field a link becomes

A link is (concept, record type, field, role). The role `instances` gives no entity field:
it only says which record type holds a concept's rows.

| Role | Kind of the concept | Entity id | Entity field id | Key field |
|---|---|---|---|---|
| `identifier` | entity | the concept's id | `id` | yes |
| `reference` | entity | the concept's id | `id` | yes |
| `display-name` | entity | the concept's id | `name` | no |
| `value` | property or dimension, with `of:` naming an entity concept of the same graph | the id of that entity concept | the concept's id | no |
| `reference` | property or dimension, with `of:` naming an entity concept of the same graph | the id of that entity concept | the concept's id | no |
| `value` | measure | not imported (reason `measure`) | | |

Rules:

1. **From field to column.** The link's record type and field are turned into a table and
   a column by the mapping. A field with no column in this database is reported as not
   imported (reason `no-column`).
2. **One meaning for one column.** When several links reach one column, one is taken, in
   this order: `identifier`; `reference` where the concept's kind is entity;
   `display-name`; `value`; `reference` where the concept's kind is property or
   dimension. Within one rank the lowest concept id wins, comparing bytes.
   The others are reported as not imported, with the reason `second-meaning`.
3. **One column of a table for one entity field.** DataTug allows an entity field one
   mapping in a table (`pkg/datatug/entities.go:165-178`). Columns that hold an
   `identifier` are placed first. A later column of the same table that would take the
   same entity field takes its next link in the order of rule 2, or nothing. The link
   passed over is reported with the reason `entity-field-taken`.
4. **Keys of several fields.** An `identifier` on a record type whose key has several
   fields is not imported (reason `several-field-key`). A reference to such a record type
   is not imported either (reason `reference-to-several-field-key`).
5. **Types.** An entity field's type comes from the model's field: `string` → `string`,
   `int` → `integer`, `float` → `float`, `bool` → `boolean`, `decimal` → `decimal`,
   `uuid` → `UUID`, `date` → `date`, `time` → `time`, `datetime` → `datetime`,
   `document`, `json` and `any` → `JSON`. A reference field has the type of its target's
   key field. The field `id` has the type of the `identifier` column.
6. **Titles.** An entity's title is the concept's English name (`labels.en`) when it has
   one.
7. **No name patterns.** An imported entity field has no `NamePatterns`. The importer
   states; it does not guess.
8. **Entity files win.** A column that an entity file maps with a declared mapping gets no
   imported mapping (reason `entity-file-maps-column`). An imported entity whose id equals
   the id of an entity file's entity, comparing bytes, is left out whole (reason
   `entity-file-has-id`).

Reasons not named above: `measure`; `no-owning-entity`, for a property or dimension whose
`of:` is absent or names no entity concept in the same graph; `component-path`, for a
field written as a path into a component; `no-column`, for a field that has no column in
this database, because no line maps its record type to a table or because the table has
no such column.

### Table schemas

For each table of the shape, `semantic.TableSchema` is filled from the shape file:
`PrimaryKey` from `primaryKey`; `ForeignKeys` from `foreignKeys`, one per `constraint`,
with `Columns`, `RefTable` and `RefColumns` in `position` order; `ReferencedBy` from the
foreign keys of every table that name this one, the table itself included
(`Employee.ReportsTo` points at `Employee`).

A reference that the model states and the database does not declare is added to the same
two lists, under the name `model:<RecordType>.<Field>`, pointing at the target's key
column. It is not written to any file. A reference whose target has a key of several
fields is not added.

### Chinook, worked through

Computed by the specification's author with a script that applies the rules above to
Chinook's four files at `3e7bb31`, using the reference `deriveLinks` of
meaninggraph/core@65510ae. It is a calculation, not an implementation.

- 48 links: 36 written, 12 derived.
- 10 imported entities, 25 entity fields, 35 of 64 columns with a meaning.
- 10 links have the role `instances`. 35 give a column its meaning. 3 are reported as not
  imported: `music-sales` on `Invoice.Total` (`measure`); `support-rep` on
  `Customer.SupportRepId` (`second-meaning`); `employee` on `Employee.ReportsTo`
  (`entity-field-taken`).

| Entity | Entity field | Columns |
|---|---|---|
| `customer` | `id` | `Customer.CustomerId`, `Invoice.CustomerId` |
| `customer` | `customer-country` | `Customer.Country` |
| `employee` | `id` | `Employee.EmployeeId`, `Customer.SupportRepId` |
| `employee` | `manager` | `Employee.ReportsTo` |
| `invoice` | `id` | `Invoice.InvoiceId`, `InvoiceLine.InvoiceId` |
| `invoice` | `invoice-date`, `billing-country`, `invoice-total` | `Invoice.InvoiceDate`, `Invoice.BillingCountry`, `Invoice.Total` |
| `track` | `id` | `Track.TrackId`, `InvoiceLine.TrackId`, `PlaylistTrack.TrackId` |

Two columns show the rules at work. `Customer.SupportRepId` has two links, `employee` (a
derived reference) and `support-rep` (a written link with the role `reference`, of a
property concept): rule 2 takes
`employee.id`. `Employee.ReportsTo` has `employee` and `manager`: `employee.id` is already
held in that table by `EmployeeId`, so rule 3 gives it `employee.manager`.

| Table | Mark with no model | Mark with the model and its meanings |
|---|---|---|
| Album | 2 of 3 | 3 of 3 |
| Artist | 1 of 2 | 2 of 2 |
| Customer | 2 of 13 | 3 of 13 |
| Employee | 2 of 15 | 2 of 15 |
| Genre | 1 of 2 | 2 of 2 |
| Invoice | 2 of 9 | 5 of 9 |
| InvoiceLine | 3 of 5 | 5 of 5 |
| MediaType | 1 of 2 | 2 of 2 |
| Playlist | 1 of 2 | 2 of 2 |
| PlaylistTrack | 2 of 2 | 2 of 2 |
| Track | 4 of 9 | 7 of 9 |
| All | 21 of 64 | 35 of 64 |

## 4. Labels

| Label | A person reads | Means | Carries |
|---|---|---|---|
| `observed` | "from the database" | The database itself declares it: a primary key or a foreign key that a scan read. | nothing more |
| `published` | "published @ 3e7bb31" | A model, or its meaning file, registered and read at a commit, states it. | the address and the commit, on the model line |
| `declared` | "declared" on a column; "in the model" on a key or a reference | The project states it: an entity file's mapping, or a model or meaning file that is a file of the project. | the file's path, on the model line or the entity page |
| `inferred` | "inferred" | A name pattern of an entity file matched the column's name. As today. | the rule, as today |

Which fact carries which label:

| Fact | Stated by | Label |
|---|---|---|
| A column's meaning | an entity file's declared mapping | `declared` |
| A column's meaning | a model and meaning file named by address | `published` |
| A column's meaning | a model and meaning file that are project files | `declared` |
| A column's meaning | an entity file's name pattern | `inferred` |
| A key or a reference | the shape file | `observed` |
| A key or a reference | the model only: the database declares none | the model's label, `published` or `declared` |
| A key or a reference | both the shape file and the model | `observed`, and the model's label is listed beside it: `from the database · also published @ 3e7bb31` |

## 5. Routes

All three are `GET`, relative to `/datatug`, take `Scope` and `source` as
[`core-investigation-loop/api-contract.md`](../core-investigation-loop/api-contract.md)
defines them, read project metadata only, and never write. Each applies the serving
session's access policy: a table is in a list and in a count only when the principal may
read it whole (REQ:read-only-routes of the Feature).

```ts
type Known = 'observed' | 'published' | 'declared';

// GET openlayers/database
type DatabaseDescription = {
  shape: null | { foreignKeys: 'read' | 'partly-read' | 'not-read' };  // null: never scanned
  model:
    | { state: 'none' }                                      // no source file
    | { state: 'named'; known: 'published' | 'declared';
        address?: string; commit?: string;                   // the Registered form
        files?: string[];                                    // the In-the-project form
        meanings: boolean }                                  // false: no `meaning` key
    | { state: 'unreadable'; problems: string[] };
  tables: {
    name: string; kind: 'table' | 'view';
    columns: number; described: number;
    key: null | { columns: string[]; known: Known; alsoStatedBy: Known[] };
    references: { columns: string[]; table: string; referencedColumns: string[];
                  known: Known; alsoStatedBy: Known[] }[];
  }[];
  notImported: { concept?: string; recordType: string; field?: string;
                 role?: string; reason: string }[];
};

// GET openlayers/match
type MatchReport = {
  registry: { read: boolean; problem?: string };
  matches: { id: string; address: string; commit: string;
             tablesMatched: number; tables: number;
             columnsMatched: number; columns: number;
             sourceFile: string }[];                         // the text of a source file that names it
};

// GET openlayers/drift
type DriftReport = {
  compared: boolean; reason?: string;                        // false: no model, or no shape
  findings: { kind: string; severity: 'error' | 'warning' | 'note';
              table?: string; column?: string; recordType?: string; field?: string;
              message: string }[];
};
```

The shapes of the four existing routes change in one place only: the values of
`provenance` in the response of `GET semantic/columns` and of `mapping` in a `Fact` gain
`published` (Question 3 of the Feature). Their behaviour changes in two more places, both
in `POST semantic/related`: one lookup is listed once, and for a project that has a shape
file the keys and foreign keys come from that file (REQ:routes-unchanged of the Feature).

## 6. Commands

```
datatug db status --project <id> | --directory <path>  --env <env> --db <db>  [--format text|json]
datatug db match  --project <id> | --directory <path>  --env <env> --db <db>  [--format text|json]
                  [--source-file <model address>]
datatug db drift  --project <id> | --directory <path>  --env <env> --db <db>  [--format text|json]
```

`--format json` prints the body of the route of the same name. None of the three opens the
database: each reads project files, and `match` reads the registries' indexes. `--db` is
the catalog's id, as in `datatug scan`. The three apply no access policy: they read the
project folder as the person who runs them.

| Command | Text output on Chinook | Exit code |
|---|---|---|
| `db status`, no source file | `Database chinook (local): 11 tables, 64 columns. Foreign keys read.`<br>`No model.`<br>then one line per table, `Invoice        2 of 9` | 0 |
| `db status`, with Example A | the same, with `Model: modelspec://github.com/demo-db/chinook/chinook @ 3e7bb31 · published` and `Invoice        5 of 9` | 0; 1 when the source file is unreadable |
| `db match` | `Registered model chinook matches this database: 11 of 11 tables, 64 of 64 columns.`<br>`  modelspec://github.com/demo-db/chinook/chinook @ 3e7bb316d8a438eb7e5930067997a4e4543a275c` | 0, whether or not a model matches |
| `db match --source-file modelspec://github.com/demo-db/chinook/chinook` | exactly the text of Example A, and nothing else | 0; 1 when that model does not match |
| `db drift`, with Example A | `Model and database agree: 11 tables, 64 columns, 11 keys, 11 references.` | 0 with no finding of severity `error`; 1 otherwise |
| `db drift`, after `Invoice.Total` was dropped and the database scanned again | `Model and database differ in 1 place.`<br>`error  field-without-column  Invoice.Total  The model's field Total has no column in the table Invoice.`<br>One line per finding: severity, kind, where, message. Errors first, then warnings, then notes. | 1 |
| any, with a bad flag or an unreadable project | the usage or the error | 2 |

## 7. Drift findings

| Kind | Severity | Says |
|---|---|---|
| `source-file-invalid` | error | The source file breaks a rule of section 2. One finding per broken rule. |
| `table-not-in-database` | error | The source file names a table that the shape does not have. |
| `record-type-not-in-model` | error | The source file names a record type that the model does not have. |
| `column-not-in-database` | error | The source file names a column that the table does not have. |
| `field-not-in-model` | error | The source file names a field that the record type does not have. |
| `field-without-column` | error | A field of a mapped record type has no column in its table. |
| `key-differs` | error | The record type's key, turned into columns, is not the table's primary key, as an ordered list. |
| `reference-target-differs` | error | The model's reference and the database's foreign key on one column name different tables. |
| `type-differs` | warning | The column's type family is not one the field's type accepts (table below). |
| `required-but-nullable` | warning | The field is `required` and the column allows null. |
| `foreign-key-not-in-model` | warning | The database declares a foreign key on a mapped column, and the field is not a reference. |
| `reference-not-declared` | note | The model states a reference, and the database declares no foreign key for it. Not reported when `scan.foreignKeys` is `not-read` or `partly-read`. |
| `column-without-field` | note | A column of a mapped table is not a field of its record type: it is outside the model. |
| `table-outside-model` | note | A table that no item of `recordsets` names. |
| `type-not-compared` | note | The column's type text is not in the table below. One finding per distinct type text. |

A test compares a finding's `kind` and names, not its `message`.

Type families. The family of a column is found from its type text, without what follows
the first `(`, ignoring case:

| Family | Type texts |
|---|---|
| integer | `INTEGER`, `INT`, `BIGINT`, `SMALLINT`, `TINYINT`, `INT2`, `INT4`, `INT8`, `SERIAL`, `BIGSERIAL` |
| text | `TEXT`, `CHAR`, `NCHAR`, `VARCHAR`, `NVARCHAR`, `CHARACTER`, `CHARACTER VARYING`, `CLOB`, `NTEXT` |
| decimal | `NUMERIC`, `DECIMAL`, `MONEY` |
| float | `REAL`, `FLOAT`, `DOUBLE`, `DOUBLE PRECISION`, `FLOAT4`, `FLOAT8` |
| boolean | `BOOLEAN`, `BOOL`, `BIT` |
| date | `DATE` |
| time | `TIME`, `TIME WITHOUT TIME ZONE`, `TIME WITH TIME ZONE` |
| datetime | `DATETIME`, `DATETIME2`, `SMALLDATETIME`, `TIMESTAMP`, `TIMESTAMPTZ`, `TIMESTAMP WITHOUT TIME ZONE`, `TIMESTAMP WITH TIME ZONE` |
| uuid | `UUID`, `UNIQUEIDENTIFIER` |
| json | `JSON`, `JSONB` |
| binary | `BLOB`, `BYTEA`, `BINARY`, `VARBINARY` |

A shape file whose `scan.driver` is `postgres` is read by this list instead: `int` is
integer; `string` is text; `bool` is boolean; `float` is float; `decimal` is decimal;
`bytes` is binary; `time` is the family temporal, which the field types `date`, `time` and
`datetime` all accept, because DALgo does not tell the three apart; `null` and `unknown`
are not compared (`type-not-compared`). PostgreSQL's `uuid` and `json` columns arrive as
`string` (dalgo2postgres v0.9.1 `type_mapping.go:62-65`), which the field types `uuid`,
`document` and `json` accept as text.

| The field's type | Families it accepts |
|---|---|
| `string` | text |
| `int` | integer |
| `float` | float, decimal |
| `decimal` | decimal |
| `bool` | boolean, integer |
| `uuid` | uuid, text |
| `date` | date, temporal, text |
| `time` | time, temporal, text |
| `datetime` | datetime, temporal, text |
| `document`, `json` | json, text |
| `any` | every family |
| a reference | the families that the type of the target's key field accepts |

On Chinook as a SQLite scan reads it, every one of the 64 columns passes: the pairs are `int` with `INTEGER` (13
columns, and 11 more through a reference), `string` with `NVARCHAR(n)` (34), `datetime`
with `DATETIME` (3) and `decimal` with `NUMERIC(10,2)` (3). Each of the 30 required fields
has a column that refuses null. Computed by the author from Chinook's model and descriptor
at `3e7bb31`.
