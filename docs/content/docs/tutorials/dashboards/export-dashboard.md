---
weight: 7083
title: "Exporting Dashboard Data to Excel Templates"
description: "Generate spreadsheet files from dashboard filters and queries using ETLX export templates."
icon: description
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Exporting Dashboard Data to Excel Templates

Central-Set dashboards are not limited to visualization.

They can also:

- 📊 Render analytics interactively
- 📁 Generate structured Excel files
- 🏛️ Produce regulatory reports
- 🔁 Reuse the exact same dashboard filters and queries
- 🔐 Enforce the same access control and row-level security

This is achieved by combining:

- Dashboard compiled queries
- ETLX export templates
- A dashboard action button (`fill-template`)
- Optional backend validation against ETLX app configuration

---

### Why This Matters

Many organizations must:

- Fill regulator-mandated Excel templates
- Generate monthly operational reports
- Export filtered data for compliance
- Deliver structured files to partners

Instead of:

- Rewriting queries
- Maintaining duplicate export logic
- Hardcoding spreadsheets

Central-Set allows you to:

> Use the same dashboard state and SQL logic  
> and generate a structured Excel file automatically.

---

## How It Works

### Step 1 — Add an Export Button

Inside your dashboard filter bar:

```html {linenos=table}
<GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
    <Button 
        tooltip="Export" 
        name="export_template" 
        action="fill-template" 
        label="" 
        icon="document-arrow-down" 
        _class='btn-sm btn-gost' 
    />
</GridItem>
```

#### What `action="fill-template"` does

When clicked:

1. Captures current dashboard state (filters, compiled SQL)
2. Replaces ETLX placeholders
3. Sends configuration to backend
4. Backend runs ETLX export
5. Returns generated file path
6. File becomes downloadable

---

## Step 2 — Define the ETLX Export Block

Inside the same dashboard Markdown, add:

`````markdown {linenos=table}
<!-- ETLX CODE BLOCK - EXPORT TEMPLATE -->

````markdown export_template
# EXPORT_TEMPLATE
```yaml
name: ExportLogsToXlsxTempl
description: Exports logs to xlsx template
runs_as: EXPORTS
connection: "duckdb:"
path: static/uploads/
active: true
```
## TEMPLATE
```yaml
name: ExportLogsToXlsxTempl
description: Exports logs to xlsx template
connection: "duckdb:"
before_sql: add_logs_table
template: logs_template.xlsx
path: tmp/logs_template_YYYYMMDD.xlsx
mapping:
  - sheet: summary
    range: A1
    sql: big_numbers
    type: value
    key: total
  - sheet: details
    range: B2
    sql: details
    type: range
    table: details
    table_style: TableStyleLight1
    header: true
    if_exists: delete
after_sql: null
active: true
```
```sql
-- add_logs_table
CREATE OR REPLACE TABLE "LOGS" AS 
SELECT * 
FROM 'static/uploads/tmp/[dash.pre_prepared_parquets.LOGS]'
```
```sql
-- details
[dash.compiled_queries._logs]
```
```sql
-- big_numbers
[dash.compiled_queries.big_numbers_query]
```
````
`````
---

## Understanding the Dynamic Placeholders

These special tokens are replaced at runtime:

### `dash.pre_prepared_parquets.LOGS`

Resolved to the actual parquet file used by the dashboard.

### `dash.compiled_queries._logs`

Resolved to the fully compiled SQL including:

* Selected filters
* Current dashboard inputs
* Main process selection
* Success filter
* Date filter
* Row-Level Access constraints (if any)

### `dash.compiled_queries.big_numbers_query`

Uses the same compiled query powering KPI tiles.

This guarantees:

> The exported Excel matches exactly what the user sees.

---

## Backend Execution Flow

When the export button is clicked:

1. Dashboard sends:

   * Dashboard config
   * Current filter values
   * Compiled SQL
   * Export block
2. Backend validates:

   * User token
   * App permissions
   * ETLX permission
3. ETLX runs in server context
4. XLSX file is generated
5. File path is returned

---

## Security Model

Exports follow the same security model as:

* `/dyn_api`
* OData
* Arrow Flight

Meaning:

* Authorization header required
* Token must belong to authorized user
* Row-Level Access is applied
* Field-level restrictions apply
* App-level isolation enforced

---

## Optional: Restricting Export to ETLX App

Client-side execution of ETLX configuration may be restricted.

In this case, define at dashboard config level:

```json {linenos=table}
"export_template": {
    "etlx_id": 1,
    "app": {
        "app_id": 2,
        "app": "ETLX",
        "db": "ETLX"
    }
}
```

This ensures:

* Export template must exist in ETLX app
* Backend checks:

  * Template ID
  * App ownership
  * Administrator-defined configuration
* Prevents arbitrary template execution

This adds an additional governance layer.

---

## Regulatory Reporting Use Case

Example scenario:

You must fill a regulator-provided Excel file:

```
Regulator_Template.xlsx
```

With:

* KPI summary in Sheet 1
* Transaction breakdown in Sheet 2
* Aggregations in Sheet 3

Instead of manually:

* Copy/pasting data
* Running scripts
* Maintaining duplicate SQL

You:

1. Build dashboard
2. Validate logic visually
3. Attach ETLX template mapping
4. Export final regulated file

Same logic.
Same filters.
Same permissions.
Same dataset.

---

## Supported Mapping Types

| Type    | Description                       |
| ------- | --------------------------------- |
| value   | Single scalar value               |
| range   | Multi-row dataset                 |
| table   | Structured Excel table with style |
| formula | (future support)                  |

---

## Limitations

* Client-side ETLX execution may be restricted
* Template must exist server-side if governed
* Large exports depend on backend resource limits
* XLSX template must exist in configured path

---

## Recommended Architecture

For production environments:

* Store templates in ETLX app
* Use `etlx_id` binding
* Enforce permission checks
* Keep dashboard export block minimal
* Version-control ETLX templates

---

## Full Flow Summary

Dashboard UI
↓
User sets filters
↓
Dashboard compiles SQL
↓
User clicks Export
↓
Placeholders replaced
↓
Backend validates permissions
↓
ETLX runs template
↓
Excel file generated
↓
Download link returned

---

## Why This Is Powerful

This turns dashboards into:

* Operational reporting tools
* Regulatory submission generators
* Controlled export engines
* Governed Excel automation pipelines

All using:

* Markdown
* SQL
* ETLX
* Central-Set security model

No duplicated logic.
No hidden scripts.
No inconsistent exports.

---

## Final Thoughts

With this pattern, Central-Set evolves from:

📊 Visualization platform
to
📊 + 📁 Enterprise Reporting Engine

And because ETLX is SQL-first and specification-driven:

Your dashboard becomes:

* Documentation
* Execution plan
* Export configuration
* Governance artifact

All in one place.