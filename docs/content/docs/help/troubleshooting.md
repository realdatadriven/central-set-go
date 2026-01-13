---
weight: 2115
title: "Troubleshooting"
icon: "tools_wrench"
description: "Solutions to common problems."
lead: "Solutions to common problems."
date: 2025-11-12T15:22:20+01:00
lastmod: 2025-11-12T15:22:20+01:00
draft: false
images: []
toc: true
---

{{% alert context="warning" text="**Caution** - This documentation is in progress" /%}}

## Troubleshooting

This section provides solutions to common issues you may encounter while using **ETLX**.

Because ETLX is **declarative**, **SQL-first**, and **configuration-driven**, most problems are not runtime crashes but **execution decisions**: blocks not running, steps being skipped, or conditions preventing execution.

> **Important mental model:**  
> If ETLX executed *anything*, it was logged.  
> If it did *nothing*, the reason is almost always visible in logs or terminal output.

If you don’t find a solution here, please open an issue or start a discussion.

## Common Issues and Solutions

### Issue 1: You execute a pipeline, but nothing happens

**Symptoms**

* The command finishes immediately
* No errors are shown
* No data is created or modified

**Possible Causes & Solutions**



#### 1️⃣ All blocks or steps are inactive

ETLX only executes blocks and steps marked as `active: true`.

* If `active` is missing → it is considered **active by default**
* If `active: false` → the block (and all children) are skipped

**Check**

```yaml
active: true
```

If a parent block is inactive, all nested blocks are ignored.

#### 2️⃣ The pipeline does not match the execution mode

Some blocks require a specific `runs_as` mode (e.g. `ETL`, `QUERY_DOC`, `EXPORTS`, `SCRIPTS`).

If you run:

```bash
etlx --config pipeline.md
```

But the top-level block has:

```yaml
runs_as: ETL
```

Then ETLX may skip execution depending on flags and configuration.

✔ Ensure the `runs_as` value matches your intent and the execution context.

#### 3️⃣ No executable SQL was defined

ETLX executes **explicit SQL references** such as:

* `load_sql`
* `query_sql`
* `export_sql`
* `script_sql`

If a step only contains metadata and no executable SQL, it will be skipped **by design**.

This is intentional — metadata alone is not execution.

#### 4️⃣ Conditions prevented execution

Some steps may include conditional execution (e.g. validations or `<step>_condition`).

If a condition evaluates to false, the step will not run.

✔ Check validation rules and condition SQL.
✔ The evaluation result is logged.

### Issue 2: You can’t find the logs for debugging

**Symptoms**

* No visible logs in the terminal
* You expect SQL-level output but see nothing

### How logging works in ETLX (important)

✔ **Everything is logged**
✔ Logs are generated for:

* Execution steps
* SQL execution
* Errors
* Validation results
* Skipped steps and reasons

If logs are **not explicitly configured to be stored somewhere**, ETLX will:

> 👉 Save them to the **operating system’s default temporary directory**

This guarantees logs always exist, even if you did not define a logging sink.

#### 1️⃣ Default logging behavior

Depending on configuration, logs may be:

* Printed to the terminal (stdout)
* Stored in the OS temp directory
* Written to files
* Persisted in a database
* Exported as structured artifacts

If you did not configure a destination, check:

* The terminal output where ETLX was executed
* The OS temporary folder

#### 2️⃣ The terminal output matters

ETLX is explicit by design.

The terminal where execution happens often tells you in in debugg mode `ETLX_DEBUG_QUERY=true`:

* Which blocks were detected
* Which steps were skipped
* Why something didn’t run
* Which SQL failed

👉 **Always read the terminal output first**.

#### 3️⃣ Use EXPORTS for observability

ETLX encourages **observable pipelines**.

You can define an `EXPORTS` block to generate:

* Execution summaries
* Logs in tabular form
* Validation results
* Data quality reports
* Audit trails

Logs are not only runtime artifacts — they can be **persisted, queried, and audited** like any other dataset.

### Issue 3: A step is silently skipped

**Symptoms**

* Some steps run, others don’t
* No errors are raised

**Possible Causes**

* `active: false` on the step
* Parent block inactive
* Validation rule failed
* Condition evaluated to false

**Tip**

Think of ETLX as a **declarative execution graph**:

> If the configuration says “don’t run”, ETLX won’t — and it will log why.

### Issue 4: SQL runs fine in my database, but fails in ETLX

**Symptoms**

* SQL works in psql / mysql / SSMS
* Fails or behaves differently in ETLX

**Possible Causes**

#### 1️⃣ Different execution engine

ETLX may execute SQL via:

* DuckDB
* Another attached engine

Some SQL syntax, functions, or data types are engine-specific.

✔ Verify:

* SQL dialect
* Extensions
* Function availability

#### 2️⃣ Missing ATTACH or extension

If your SQL references external data:

```sql
FROM SRC.my_table
```

Make sure the source is attached:

```yaml
before_sql:
  - ATTACH 'postgres:@PG_DSN' AS SRC (TYPE POSTGRES)
```

If it’s not attached, ETLX cannot resolve it — and this will be logged.

### Issue 5: Data is duplicated or partially loaded

**Symptoms**

* Re-running a pipeline duplicates data
* Unexpected extra rows

**Explanation**

ETLX does **not** assume idempotency by default.

You must explicitly define:

* Validation rules
* Conditional execution
* Cleanup logic

**Best Practices**

* Use `load_validation`
* Use date-based or watermark checks
* Fail fast if data already exists

### Issue 6: Exports or templates generate empty output

**Symptoms**

* Export file exists but is empty
* Template renders without content

**Possible Causes**

* `data_sql` returned no rows
* Template logic does not match the data structure
* You are iterating over a key that does not exist

**Tip**

> ETLX templates can consume **any data**, including the parsed configuration itself (`.conf`).

If something exists in the config, it can be exported.

### Issue 7: I expected documentation or governance artifacts, but nothing was generated

**Explanation**

ETLX only generates artifacts that are **explicitly declared**.

Documentation, data dictionaries, and audits are **pipeline outputs**, not side effects.

✔ Ensure you have an `EXPORTS` block that:

* References metadata
* Uses a text or HTML template
* Targets a file or destination

## General Debugging Tips

* Read the terminal output
* Assume logs exist — find where they are stored
* Start simple: one block, one step
* Verify `active` flags
* Check execution order
* Inspect parsed configuration via exports
* Prefer explicit behavior over assumptions

## Still stuck?

If the issue persists:

* Open an issue on GitHub
  👉 [https://github.com/realdatadriven/etlx/issues](https://github.com/realdatadriven/etlx/issues)

* Include:

  * Minimal pipeline example
  * Expected vs actual behavior
  * Execution engine and version
  * Relevant logs or terminal output

ETLX favors **transparency and clarity** — most issues can be understood by inspecting configuration, logs, and execution output.