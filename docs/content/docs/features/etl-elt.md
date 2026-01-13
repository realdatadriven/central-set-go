---
weight: 710
title: "ETL | ELT"
description: "How ETLX models, executes, and observes ETL and ELT pipelines using declarative, metadata-driven configurations."
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---
---

## ETL | ELT in ETLX

ETLX supports both **ETL (Extract–Transform–Load)** and **ELT (Extract–Load–Transform)** execution models.

Rather than enforcing a specific pattern, ETLX lets you **declare the intent of the pipeline**, and then executes it deterministically based on metadata.

At the highest level, an ETL or ELT pipeline is defined by:

* A **root block** describing the pipeline
* One or more **execution units** (inputs, transformations, outputs)
* Explicit **SQL blocks** defining the logic
* Optional **error handling and lifecycle hooks**

---

## Defining an ETL or ELT Pipeline

A pipeline is declared using a **level-one Markdown heading** (`#`) combined with metadata.

The execution mode is defined using the `runs_as` key:

* `runs_as: ETL`
* `runs_as: ELT`

### Example

````md {linenos=table}
# INPUTS
```yaml
name: INPUTS
description: Extracts data from source and loads it into the target
runs_as: ETL
active: true
```

## INPUT_1
```yaml
name: INPUT_1
description: Input 1 from an ODBC source
table: INPUT_1
load_conn: "duckdb:"
load_before_sql:
  - "ATTACH 'ducklake:@DL_DSN_URL' AS DL (DATA_PATH 's3://dl-bucket...')"
  - "ATTACH '@OLTP_DSN_URL' AS PG (TYPE POSTGRES)"
load_sql: load_input_in_dl
load_on_err_match_patt: '(?i)table.+with.+name.+(\w+).+does.+not.+exist'
load_on_err_match_sql: create_input_in_dl
load_after_sql:
  - DETACH DL
  - DETACH PG
active: true
```

```sql
-- load_input_in_dl
INSERT INTO DL.INPUT_1 BY NAME
SELECT * FROM PG.INPUT_1
```

```sql
-- create_input_in_dl
CREATE TABLE DL.INPUT_1 AS
SELECT * FROM PG.INPUT_1
```
````

> `@DL_DSN_URL` (e.g. `mysql:db=ducklake_catalog host=your_mysql_host`) and `@OLTP_DSN_URL` (e.g. `postgres:dbname=erpdb host=your_postgres_host user=postgres password=your_pass`) are **environment variables** used to define database connection strings.
>They can be provided through a `.env` file located at the root of the project and are automatically loaded at runtime.
>These variables allow ETLX to connect to different data sources without hardcoding credentials, making configurations portable, secure, and environment-agnostic.

## Execution Model

### 1. Pipeline Initialization

Execution starts at the **root pipeline block** (`# INPUTS`).

From this block, ETLX extracts:

* **Pipeline metadata**

  * `name`
  * `description`
  * `runs_as` (ETL or ELT)
* **Global execution context**

  * Default connection
  * Execution timestamp
  * Runtime metadata
* **Activation state**

  * Pipelines marked as `active: false` are skipped

This metadata becomes part of the **execution trace** and **observability layer**.


### 2. Iteration Over Execution Units

Each **level-two heading** (`## INPUT_1`, `## TRANSFORM_X`, etc.) represents an **execution unit**.

For each unit, ETLX:

1. Reads its metadata
2. Resolves connections
3. Determines which steps apply
4. Executes steps in a deterministic order

Inactive units (`active: false`) are skipped but still recorded in metadata.


## ETL / ELT Steps

Each execution unit may define one or more **steps**, depending on the execution model:

* `extract`
* `transform`
* `load`

Each step supports a consistent lifecycle:

### Step Lifecycle Hooks

For any step `<step>` (e.g. `load`):

| Hook                | Purpose                                        |
| ------------------- | ---------------------------------------------- |
| `<step>_before_sql` | Setup (e.g. attach databases, prepare schemas) |
| `<step>_sql`        | Main execution logic                           |
| `<step>_after_sql`  | Cleanup (detach, finalize, release resources)  |


### SQL Resolution Rules

Each SQL hook can be defined as:

* `null` → step is skipped
* `string` → reference to a named SQL block
* inline SQL string
* list/array → executed sequentially

Named SQL blocks are resolved from:

* `sql [query_name]` blocks
* Or SQL comments: `-- query_name`

This allows **clear separation of metadata and logic**.

## Connection Handling

Each step can specify a connection using `<step>_conn`.

* If defined → used for that step
* If `null` or omitted → falls back to the pipeline’s default connection

This enables **multi-engine pipelines** (DuckDB, Postgres, ODBC, etc.) within a single execution.

## Error Handling & Recovery

ETLX provides **pattern-based error handling**, allowing pipelines to recover dynamically from known failure conditions.

For any step `<step>`:

* `<step>_on_err_match_patt`
  A regular expression matched against the database error message
* `<step>_on_err_match_sql`
  SQL executed when the pattern matches

This is especially useful for:

* Creating missing tables
* Initializing schemas
* Handling first-run scenarios
* Working with evolving datasets

The same mechanism applies to:

* `<step>_before_on_err_match_*`

## Observability & Execution Metadata

Every ETLX run is **fully observable by design**.

For each pipeline, step, and sub-step, ETLX records:

* Start and end timestamps
* Execution duration
* Connection used
* SQL executed
* Rows affected (when available)
* Errors, warnings, and retries
* Memory and resource usage

This metadata can be:

* Logged
* Stored
* Queried via SQL
* Used to generate reports and documentation

## Design Principles

ETLX treats ETL and ELT pipelines as:

* **Declarative execution plans**
* **Structured metadata documents**
* **Executable documentation**

This ensures pipelines are:

* Reproducible
* Inspectable
* Auditable
* Portable across environments
* Easy to maintain and evolve

## Summary

In ETLX:

* ETL and ELT are **modes**, not rigid architectures
* Configuration defines *intent*, not orchestration
* SQL remains first-class
* Metadata powers execution, observability, and documentation
* Pipelines are self-describing and deterministic

> **Your pipeline configuration is your source of truth.**