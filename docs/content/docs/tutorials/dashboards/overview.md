---
weight: 7081
title: "Overview"
description: "Overview"
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Overview

Central Set provides a **first-class analytics and dashboarding layer** by embedding and extending
[evidence.dev](https://docs.evidence.dev).

This allows you to build **interactive analytical dashboards** directly on top of:

* ETLX pipelines
* Exported datasets
* Parquet files
* Curated analytical tables

All dashboards live **inside Central Set**, inherit authentication and authorization, and are fully
integrated with the platform's data model.

---

## Why Evidence.dev?

Evidence.dev was chosen because it is:

* SQL-first
* Component-based
* Open-source
* Developer-friendly
* Designed for analytical storytelling

However, Central Set **does not use Evidence in its default standalone mode**.

Instead, Evidence is **embedded as a dashboard frontend**, while Central Set controls:

* Configuration
* Data access
* Execution context
* Security
* Dataset lifecycle

In Central Set, Evidence configuration is rendered **in real time in the browser**, rather than being
compiled into a static website.

Conceptually, Evidence behaves as a **Svelte-based rendering layer** that intercepts normal CRUD-based
UI rendering and transforms database-backed records into analytical dashboards.

Some Evidence components are not yet integrated or are slightly modified to better fit Central Set,
but **the core visualization components (what really matters) are fully compatible**, including:

* [If / Else](https://docs.evidence.dev/core-concepts/if-else/)
* [Loops](https://docs.evidence.dev/core-concepts/loops/)
* [Formating](https://docs.evidence.dev/core-concepts/formatting/)
* [Value](https://docs.evidence.dev/components/data/value/)
* [Big Value](https://docs.evidence.dev/components/data/big-value/)
* [Data Table](https://docs.evidence.dev/components/data/data-table/)
* [Delta](https://docs.evidence.dev/components/data/delta/)
* [Charts](https://docs.evidence.dev/components/charts/area-chart/)

Charts are based on [Apache ECharts](https://echarts.apache.org/en/index.html) and can be fully customized using:

* [ECharts Extra Options](https://docs.evidence.dev/components/charts/echarts-options/)
* [Custom Charts](https://docs.evidence.dev/components/charts/custom-echarts/)

This allows the full ECharts configuration object to be passed directly, enabling virtually any
visualization supported by ECharts.

---

## What's Different from Vanilla Evidence

| Area            | Evidence.dev                | Central Set             |
| --------------- | --------------------------- | ----------------------- |
| Data Sources    | Inline SQL / DB connections | Central Set APIs        |
| Authentication  | None / external             | Native Central Set auth |
| Dataset Storage | Ad-hoc queries              | ETLX datasets & Parquet |
| Runtime         | Static site / dev server    | Embedded application    |

In short:

> **Evidence provides the rendering engine.
> Central Set provides the platform.**

---

## Data Flow: ETLX → Dashboards

The dashboard layer is designed to sit **downstream of ETLX**.

### Typical Flow

```
ETLX Pipeline
   ↓
Dataset / Parquet Export
   ↓
Central Set Dataset Registry
   ↓
Evidence Dashboard Components
```

ETLX pipelines produce:

* Parquet files
* Materialized datasets
* Analytical views

These outputs are **registered and managed** by Central Set and exposed to dashboards as **named datasets**.

Dashboards never access raw databases directly.

---

## Dataset-Centric Architecture

Instead of writing SQL directly inside dashboards:

* Datasets are **defined once**
* Produced by ETLX
* Versioned and traceable
* Reused across dashboards

This ensures:

* Reproducibility
* Performance
* Clear lineage
* Separation of concerns

> Dashboards consume data.
> Pipelines produce data.

---

## Embedded Evidence Frontend

Central Set embeds the Evidence dashboard frontend directly into the UI.

This means:

* Dashboards appear as native Central Set menu
* Navigation is unified
* Authentication is shared
* Permissions apply automatically

From the user's perspective, dashboards are simply **another application menu**.

---

## Extended Components

In addition to standard Evidence components, Central Set adds its own components, such as `<Stats>`.

Most of these are **Svelte components built with**
[DaisyUI](https://daisyui.com/components/), since Central Set itself is built on the
[DaisyUI](https://daisyui.com) design system.

This ensures visual consistency between dashboards and the rest of the platform UI.

---

## Configuration Model

Evidence dashboards normally rely on a local configuration block.

In Central Set:

* Configuration is **metadata-driven**
* Stored and versioned centrally
* Linked to applications and menus
* Managed from the ADMIN app

This allows dashboards to be:

* Activated or deactivated
* Versioned
* Assigned to profiles
* Shipped as part of an application

---

## Security & Governance

Because dashboards are embedded:

* Access is controlled by Central Set profiles
* Users only see datasets they are authorized to access
* Dashboards can be restricted by role or tenant

This makes dashboards suitable for:

* Internal analytics
* Multi-tenant environments
* Customer-facing data products

---

## Analytics as a Product

Dashboards are not just visualizations.

They are:

* Versioned artifacts
* Backed by pipelines
* Governed by metadata
* Deployable across environments

This aligns with Central Set's core philosophy:

> **Everything is configuration.
> Everything is data-driven.**

---

## How This Fits Together

| Layer      | Responsibility                   |
| ---------- | -------------------------------- |
| ETLX       | Data pipelines & transformations |
| Datasets   | Canonical analytical outputs     |
| ADMIN      | Metadata, permissions, structure |
| Dashboards | Exploration & insight            |
| Evidence   | Rendering engine                 |


# Data Sources in Dashboards

Dashboards in Central Set can load data in **two different ways**, depending on your needs:

1. **SQL queries (`sql`)**
2. **Central Set secured queries (`cs`)**

Both approaches can coexist in the same dashboard.

---

## 1️⃣ SQL Blocks (Analytics Mode)

The most common method is using SQL blocks:

````markdown
```sql my_query
SELECT *
FROM "LOGS"
WHERE "ref" = 'inputs.date_ref.value'
```
````

### How it works

* Executed:

  * In DuckDB WASM (browser)
  * Or mapped database
* Reads from:

  * Parquet files
  * ETLX outputs
  * Registered data sources
* Designed for:

  * Analytics
  * Aggregations
  * KPIs
  * Charts
  * Time series

This is ideal when:

* You need fast analytics
* You control the dataset
* Data is already scoped at source level

---

## 2️⃣ `cs` Blocks (CRUD / Secured Mode)

Dashboards can also execute **direct `crud/read` operations** using a `cs` block:

````markdown
```cs departments
  "table": "departments",
  "limit": -1,
  "order_by": ["name asc"]
```
````

This does **not** execute raw SQL.

Instead, it:

* Calls the internal `crud/read` API
* Applies:

  * Table permissions
  * Role permissions
  * Row Level Access (RLA)
  * Tenant scoping or Department scoping
  * App/database isolation
* Returns only authorized records

---

## Why `cs` Blocks Matter

Unlike `sql` blocks, `cs` blocks are:

> Fully governed by Central Set security rules.

This means:

* Row-level security is enforced automatically
* Multi-tenant isolation is respected
* Scoped domains (e.g. departments) are enforced
* Access keys and user tokens are validated

You cannot bypass security using `cs`.

---

## When to Use `cs` Instead of SQL

Use `cs` blocks when:

✔ You need data scoped by tenant
✔ You need department-level filtering
✔ You need row-level security applied
✔ You want filters driven by secured domain tables
✔ You are building dashboards for multiple roles

---

## Example: Secure Department Filter

Instead of:

````markdown
```sql departments
SELECT id, name FROM departments
```
````

You should use:

````markdown
```cs departments
  "table": "departments",
  "limit": -1
```
````

Now:

* A user from Department A only sees the A departments
* A manager sees only their department scope
* An admin sees all departments

No extra SQL or dashboard required.

---

## Mixing SQL + CS in the Same Dashboard

This is a powerful pattern:

### Use `cs` for:

* Filter domains
* Secured lookup tables
* Tenant-scoped dropdowns
* User-specific data

### Use `sql` for:

* Aggregations
* Metrics
* Historical analysis
* Parquet-based analytics

Could all be filtered by the `cs` result

Example:

````markdown
<!-- getting the departments to be added to a multi-select filter, where selected values will be pushed to all the sql blocks with `inputs.departments.value` -->
```cs departments
"table": "departments",
"fields": ["department_id", "department"]
"limit": -1
```

<!-- this way each user will see exactlly the department that he/shee can access -->
```sql total_by_department
SELECT department_id, COUNT(*) total
FROM "LOGS"
WHERE department_id IN (inputs.departments.value)
GROUP BY department_id
```
````

You can then bind:

* `departments` → Dropdown
* `total_by_department` → Chart

Security + Analytics combined.

---

## Execution Model Comparison

| Feature            | `sql` Block              | `cs` Block                 |
| ------------------ | ------------------------ | -------------------------- |
| Execution Engine   | DuckDB / DB              | Central Set API            |
| Row Level Security | ❌ Only if encoded in SQL | ✅ Always enforced          |
| Tenant Isolation   | ❌ Manual                 | ✅ Automatic                |
| Performance        | Optimized for analytics  | Optimized for secured CRUD |
| Best For           | KPIs / Charts            | Secured filters / domains  |

---

## Security Recommendation

If your dashboard is:

* Multi-tenant
* Embedded in external apps
* Used by different roles

Then:

> Use `cs` blocks for all domain filters and scoped datasets.

Let SQL handle analytics,
and let `crud/read` handle security.

---

## Advanced Use Case: Regulated Environments

In regulated systems:

* Departments
* Cost centers
* Business units
* Tenants
* Regulatory domains

Should always be retrieved via `cs`.

This guarantees:

* Compliance with RLA
* No accidental overexposure
* No manual WHERE tenant_id = ...

Security remains centralized in Central Set — not scattered across SQL.

---

## Architectural Principle

Dashboards are not just analytics views.

They are:

> A controlled execution layer that can combine analytics performance with secured operational data.

By supporting both `sql` and `cs`, Central Set enables:

* Secure operational dashboards
* Tenant-aware analytics
* Embedded BI
* Regulated data products

All in the same markdown definition.
