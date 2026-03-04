---
weight: 7084
title: "Dynamic Dataset"
description: "Dashboard With Dynamic Dataset For Each User / Profile"
icon: description
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Dashboard With Dynamic Dataset

### Overview

In scenarios where datasets must be:

- Scoped by **user**
- Scoped by **tenant**
- Scoped by **department**
- Scoped by **vendor**
- Or any Row-Level Access (RLA) domain

…and you still want:

- Pre-cached Parquet files  
- Browser-side DuckDB execution  
- High performance  
- Low query cost  

This pattern is ideal.

Instead of running every query through `crud/read` in the backend, you:

1. Generate a **custom dataset per user** and Export it as Parquet via ETLX
2. Cache it (local disk, S3, blob storage, CDN)
3. Load it in the browser via DuckDB WASM

This allows:

> Secure scoping + high-performance local analytics.

---

## Why This Pattern Exists

You could scope everything using `cs` blocks (secured CRUD).

But that means:

- Every query runs in backend
- Every filter hits the database
- Costs may increase (DB compute / cloud egress)

For datasets that:

- Change once per day
- Update overnight
- Have compute cost
- Are expensive to query repeatedly

It is better to:

> Generate a scoped dataset once  
> Cache it  
> Serve it many times  

---

## Architecture

```
User → Dashboard → Click "Update Dataset"
│
▼
ETLX Backend
│
▼
Generate User-Scoped Parquet
│
▼
Store + Log File
│
▼
Dashboard Loads Parquet in Browser
```

---

## Step 1 — Add Update Button

Add a button to trigger dataset regeneration:

```html {linenos=table}
<GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
    <Button 
        tooltip="Update Dashboard Data" 
        name="my_custom_ds_ex" 
        action="update_custom_ds"
        icon="cloud-arrow-down" 
        _icon="arrow-path-rounded-square" 
        _class='btn-sm btn-gost' 
    />
</GridItem>
```

### Important

* `name="my_custom_ds_ex"`
  → Must match:

  * Config block key
  * ETLX markdown block ID

* `action="update_custom_ds"`
  → Triggers backend ETLX execution

---

## Step 2 — Dashboard Configuration

In your dashboard config:

```config {linenos=table}
"all_query_run_locally_in_ddb_wasm": true,
"pre_prepared_parquets_logs_table": null,
"pre_prepared_parquets_logs_sql": "with _logs as (select *  from dynamic_ds_logs  where fname is not null and user_id = [dash.user.user_id] and (user_id, table_name, created_at) in ( select user_id, table_name, max(created_at) from dynamic_ds_logs group by user_id, table_name ) ) select user_id, table_name as name, replace(fname, 'tmp/', '') as file from _logs",
"pre_prepared_parquets_logs_db": "sqlite3:database/logs_for_dyn_gen_ds.db",
"pre_prepared_parquets_for_ddb_wasm": {
  "ORDERS": "orders.parquet"
},
"update_custom_ds": {
  "my_custom_ds_ex": {
    "etlx_id": 1,
    "app": { "app_id": 2, "app": "ETLX", "db": "ETLX"}
  }
}
```

---

## What This Configuration Does

### 1️⃣ Detect Latest File Per User

Uses:

```
dynamic_ds_logs
```

Filtered by:

```
user_id = [dash.user.user_id]
```

So:

* Each user loads only their dataset
* Fully isolated
* No cross-tenant leakage

---

### 2️⃣ Map File to Logical Table

```json
"pre_prepared_parquets_for_ddb_wasm": {
    "ORDERS": "orders.parquet"
}
```

This makes:

```sql
FROM "ORDERS"
```

Available inside DuckDB WASM.

---

## Step 3 — ETLX Dataset Generator

Inside your dashboard markdown:

`````markdown {linenos=table}
<!-- ETLX CODE BLOCK - EXPORT DATASET -->

````markdown my_custom_ds_ex
# GENERATE_DATA_SETS
```yaml
name: GenerateDSs
description: Exports custom ds
runs_as: EXPORTS
connection: "duckdb:"
path: static/uploads/ # in case of a s3endpoint it can be passed directlly in the export query
active: true
```
## ORDERS
```yaml
name: GenerateORDERSData
description: Exports custom ORDERS data
connection: "duckdb:"
before_sql:
  - INSTALL erpl_web FROM community
  - LOAD erpl_web
  - create_api_auth_secrete
  - attach_odata_endpoint_with_users_copes
  - attach_ex_ecomerce_datalake
  - attach_logs_db
export_sql: 
  - generate_my_sales_data
  - create_logs_table_if_not_exists
  - insert_generated_file_into_logs
after_sql:
  - DETACH scopes
  - DETACH dl
  - DETACH logs
path: tmp/orders.[dash.user.user_id].{YYYYMMDD}.{TSTAMP}.parquet
tmp_prefix: tmp
active: true
```
```sql
-- create_api_auth_secrete
CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '[dash.user.jwt_token]',
  SCOPE 'http://localhost:4444/'
);
```
```sql
-- attach_odata_endpoint_with_users_copes
ATTACH IF NOT EXISTS 'http://localhost:4444/odata/ETLX' AS scopes (TYPE ODATA);
```
```sql
-- attach_ex_ecomerce_datalake
ATTACH 'ducklake:sqlite:database/dl_metadata.sqlite' AS dl (DATA_PATH 'database/dl/');
```
```sql
-- attach_logs_db
ATTACH 'database/logs_for_dyn_gen_ds.db' AS logs (TYPE SQLITE);
```
```sql
-- generate_my_sales_data
COPY (
  SELECT orders.*
  FROM dl.orders
  WHERE orders.department_id IN (
    -- OData API Will use CS crud/read by the user in the <JWT_TOKEN> given by its session
    SELECT department_id
    FROM scopes.department
  )
) TO '<fname>';
```
```sql
-- create_logs_table_if_not_exists
CREATE TABLE IF NOT EXISTS logs.dynamic_ds_logs (
    --id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    table_name VARCHAR NOT NULL,
    fname      VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
```sql
-- insert_generated_file_into_logs
INSERT INTO logs.dynamic_ds_logs (user_id, table_name, fname) VALUES ([dash.user.user_id], 'ORDERS', '<fname>');
```
````
`````

---

## How Security Is Applied

Your datalake:

```
dl.sales.orders
```

Is **unfiltered**.

It has no Row-Level Access.

---

## RLA Is Applied Here

```sql
COPY (
  SELECT orders.*
  FROM dl.sales.orders
  WHERE orders.department_id IN (
    SELECT department_id
    FROM scopes.department
  )
) TO '<fname>';
```

Where:

```
scopes.department
```

Is an OData endpoint backed by:

* Central Set `crud/read`
* User JWT token
* RLA enforcement

Meaning:

> The dataset is generated already scoped by the user's allowed departments.

---

# Logging Generated Files

```sql
CREATE TABLE IF NOT EXISTS logs.dynamic_ds_logs (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    table_name VARCHAR NOT NULL,
    fname      VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

And then:

```sql
INSERT INTO logs.dynamic_ds_logs (user_id, table_name, fname) VALUES ([dash.user.user_id], 'ORDERS', '<fname>');
```

This enables:

* Version control
* File history
* Regeneration tracking
* Cleanup strategies

---

# Final Result

Now your dashboard queries run like normal:

```sql
SELECT department_id, SUM(amount)
FROM "ORDERS"
GROUP BY department_id
```

But:

* Data is already scoped
* Loaded from user-specific Parquet
* Executed in browser
* Zero backend queries for analytics

---

# When To Use This Pattern

Use Dynamic Dataset when:

✔ Data must be scoped by RLA
✔ Dataset is expensive to compute
✔ Dataset updates periodically (daily / monthly)
✔ You want minimal backend calls
✔ You want analytics fully client-side and fast

---

# When NOT To Use This

If:

* Data changes constantly (real-time)
* You need live transactional accuracy
* User scoping changes frequently

Then prefer:

* `cs` blocks
* Direct secured backend queries

---

# Comparison

| Approach        | Security       | Performance | Backend Load | Best For              |
| --------------- | -------------- | ----------- | ------------ | --------------------- |
| `cs` blocks     | Live RLA       | Medium      | High         | Real-time             |
| SQL only        | No RLA         | Very High   | None         | Public data           |
| Dynamic Dataset | Pre-scoped RLA | Very High   | Low          | Periodic secured data |

---

# Key Takeaway

Dynamic Dataset combines:

* ETLX execution power
* Central Set RLA security
* Parquet performance
* DuckDB WASM execution
* User isolation
* Cost control

It is the most scalable way to deliver:

> Secure multi-tenant analytics at scale.

