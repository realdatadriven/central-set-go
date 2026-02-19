---
weight: 7090
title: "OData v4 API"
description: "OData v4 API"
icon: code
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

# OData v4 API

{{% alert context="warning" text="**Documentation status** — The application itself is mature and production-ready, but this documentation is still under active development. Large parts are auto-generated and continuously evolving." /%}}

The application exposes a **standards-compliant OData v4 API** that allows controlled, secure access to relational data using HTTP-based query semantics.

This API is designed to be:

* Consumable by **data science and analytics tools**
* Easy to integrate with **DuckDB, Python, R, BI tools**
* Governed by the **same access-control rules** as the UI
* Composable with **ETLX workflows and metadata-driven pipelines**

---

## Endpoint Overview

```http
GET /odata/{db}/{table}
```

### Example

```http
GET http://localhost:4444/odata/ADMIN/menu?$filter=app_id eq 1&$format=json
```

* `db` → logical database name (e.g. `ADMIN`)
* `table` → exposed table or view
* Query options follow **OData v4** conventions

---

## Supported OData Query Options

The API supports the most commonly used OData v4 operators:

| Feature         | Example                 |
| --------------- | ----------------------- |
| Filtering       | `$filter=app_id eq 1`   |
| Comparison      | `$filter=app_id gt 1`   |
| Field selection | `$select=app_id,app,db` |
| Pagination      | `$top=50&$skip=100`     |
| Sorting         | `$orderby=app_id desc`  |
| Format          | `$format=json`          |

---

## Authentication & Authorization

> **All access rules are enforced. No exceptions.**

### 🔐 Authorization Header (Required)

Every request **must** include an access token:

```http
Authorization: Bearer <ACCESS_TOKEN>
```

### Token Management

Access tokens are created and managed in the application UI:

```
Admin → Menu → Access Keys
```

### Token Scope & Permissions

* Tokens are **always associated with a user**
* The user **must have access** to:

  * The database
  * The table or resource
* If **row-level access rules** are defined, they are automatically applied

This allows OData to be safely used for:

* Tenant isolation
* Scoped analytics
* Per-user or per-role data access

---

## Row-Level Security (RLS)

If row-level access rules exist in the application:

* They are **automatically enforced**
* They apply identically to:

  * UI access
  * OData access
  * ETLX pipelines (when applicable)

### Example Use Cases

* Restrict data by `tenant_id`
* Scope results by `user_id`
* Filter records by ownership or domain

This makes OData ideal for **secure downstream consumption**.

---

## DuckDB Integration (Recommended)

The OData API integrates seamlessly with DuckDB using the `erpl_web` extension.

### Setup

```sql {linenos=table}
INSTALL erpl_web FROM community;
LOAD erpl_web;

-- Enable tracing (optional, useful for debugging)
SET erpl_trace_enabled = TRUE;
SET erpl_trace_level = 'DEBUG';

-- Create API secret
CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '<JWT_TOKEN>',
  SCOPE 'http://localhost:4444/'
);
```

### Querying via OData

```sql {linenos=table}
-- Direct read
FROM HTTP_GET('http://localhost:4444/odata/ADMIN/app');

-- OData query
FROM ODATA_READ('http://localhost:4444/odata/ADMIN/app?$filter=app_id gt 1');

-- Attach as a database
ATTACH IF NOT EXISTS 'http://localhost:4444/odata/ADMIN' AS admin (TYPE odata);

SELECT app_id, app, db, excluded
FROM admin.app
WHERE app_id = 1;
```

## ETLX & Metadata Integration

The OData layer is intentionally compatible with **ETLX pipelines**:

* OData endpoints can be used as **data sources**
* Metadata can be merged into ETLX configs
* Enables:

  * Governance-aware ingestion
  * Reproducible analytics
  * API-driven pipelines

> At this point, **the only limitation is imagination**.


## Intended Use Cases

* Data science APIs
* Ad-hoc analytics
* BI & reporting tools
* ETLX-driven workflows
* Secure data sharing
* Debugging & inspection
