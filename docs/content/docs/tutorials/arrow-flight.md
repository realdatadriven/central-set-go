---
weight: 7091
title: "Arrow Flight"
description: "Apache Arrow Flight Support API"
icon: arrow
date: 2025-12-16T01:04:15+00:00
lastmod: 2026-01-29T00:00:00+00:00
draft: false
images: []
---

## Arrow Flight Support

{{% alert context="warning" text="**Important** - Central-Set is production-ready, but this documentation is still under active development. Parts of the Arrow Flight subsystem are configuration-driven and evolve as access-control and governance features mature." /%}}

### Overview

Central-Set provides **Apache Arrow Flight** support through  
[airport-go](https://github.com/hugr-lab/airport-go), exposing analytical datasets via a **high-performance, governed, and dynamically scoped interface**.

Arrow Flight is primarily designed to serve:

- **ETLX outputs**
- **DuckDB-backed analytical views**
- **Externally attached datasources**

While enforcing the **same security, access control, and multi-tenant rules** used throughout Central-Set.

Unlike static Flight servers, Central-Set builds Arrow Flight endpoints **entirely from configuration stored in the Admin database**, allowing schemas, tables, fields, and scopes to be enabled, restricted, or revoked at runtime — without redeploying the service.

---

## Architecture Summary

At runtime, Central-Set:

1. Loads Arrow Flight configuration from the Admin database
2. Initializes an **in-memory DuckDB instance**
3. Executes lifecycle SQL blocks:
   - `startup_sql`
   - `main_sql`
   - `shutdown_sql`
4. Exposes **only authorized tables, fields, and scopes**
5. Serves data via the Arrow Flight protocol

Each request is authenticated and authorized using the **same JWT and access-key system** as the REST API.

---

## Arrow Flight Configuration Model

Arrow Flight exposure is defined using **three core entities**, allowing fine-grained governance.

---

### ArrowFlightTable

Defines which **tables** may be exposed via Arrow Flight.

| Field | Description |
|-----|------------|
| `arrow_flight_table` | Logical table name exposed to clients |
| `arrow_flight_table_desc` | Description |
| `arrow_flight_id` | Parent Arrow Flight schema |
| `active` | Enables / disables the table |
| `user_id` | Owner |
| `app_id` | Application scope |
| `excluded` | Soft-delete flag |

If table-level access is defined, **only tables explicitly granted to the role or access token are visible**.

---

### ArrowFlightTableField

Defines **field-level visibility** per table.

| Field | Description |
|-----|------------|
| `arrow_flight_table_field` | Column name |
| `arrow_flight_table_field_desc` | Description |
| `arrow_flight_table_id` | Parent table |
| `arrow_flight_id` | Arrow Flight schema |
| `active` | Enables / disables the field |
| `excluded` | Soft-delete flag |

#### Field-Level Access Behavior

- If **field access rules exist**:
  - Fields **without access are still present**
  - But their values are returned as **NULL / empty**
- This preserves:
  - Schema compatibility
  - Stable BI / analytical queries
  - Controlled data masking

This design is intentional and avoids breaking downstream consumers.

---

### ArrowFlightTableScope

Defines **data scopes** using SQL predicates.

| Field | Description |
|-----|------------|
| `arrow_flight_table_scope` | Scope name |
| `arrow_flight_table_scope_desc` | Description |
| `arrow_flight_table_scope_sql` | SQL condition |
| `arrow_flight_table_id` | Target table |
| `arrow_flight_id` | Arrow Flight schema |
| `active` | Enables / disables scope |
| `excluded` | Soft-delete flag |

#### Scope Enforcement Rules

- If **no scopes are defined** → all rows are eligible <!--(subject to RLA)-->
- If **one or more scopes exist**:
  - The token **must have access to at least one scope**
  - Otherwise **no data is returned**
- Multiple scopes are **AND-combined**

Scopes act as a **hard gate** for data visibility.

---

## Security & Access Control

Arrow Flight follows **exactly the same security model as the REST API**.

### Authentication

- ✅ `Authorization: Bearer <token>` is **mandatory**
- ✅ Tokens are created via:
  - **Admin → Admin → Access Keys**
- ✅ Tokens may belong to:
  - A user
  - A service account
  - An automation pipeline

---

### Authorization Layers

Arrow Flight access is evaluated in the following order:

1. **Token validity**
2. **App access**
3. **Schema access**
4. **Table access**
5. **Field access (masking)**
6. **Scope access (gating)**
<!--7. **Row-Level Access (RLA)**-->

All layers must pass for data to be returned.

This makes Arrow Flight safe for:

- Multi-tenant analytics
- External BI tools
- Cross-team data sharing
- Zero-trust environments

---

## TLS / Secure Transport

Arrow Flight can run with or without TLS.

### Required Environment Variables

```env
ENABLE_TLS=false

TLS_CERT_FILE=ssl/server-cert.pem
TLS_KEY_FILE=ssl/server-key.pem
TLS_CA_CERT_FILE=ssl/ca-cert.pem
````

* When `ENABLE_TLS=true`, Arrow Flight serves **gRPC over TLS**
* Clients must trust the configured CA
* Strongly recommended for production and remote access

---

## Enabling Arrow Flight

```env
ENABLE_ARROW_FLIGHT=true
ARROW_FLIGHT_ADDR=0.0.0.0:50051
```

Arrow Flight runs **inside the same binary** as the REST API and shares:

* Authentication
* Configuration
* Access control
* Application context

---

## Defining an Arrow Flight Schema

Create a schema via:

**Admin → Expose Arrow Flight**

Example:

```yaml {linenos=table}
name: my_schema
description: Example analytical schema
db_schema: main

startup_sql: |
  INSTALL SQLITE;
  LOAD SQLITE;

main_sql: |
  ATTACH 'database/test.db' AS my_schema (TYPE SQLITE);
  USE my_schema;

shutdown_sql: |
  USE memory;
  DETACH my_schema;
```

Each schema represents **one logical Arrow Flight endpoint**.

---

## Client Access (DuckDB - Recommended)

Using DuckDB's **airport** extension:

```sql {linenos=table}
INSTALL airport FROM community;
LOAD airport;

CREATE OR REPLACE [PERSISTENT] SECRET airport_auth_secret (
    TYPE airport,
    AUTH_TOKEN 'your_access_token_here',
    SCOPE 'grpc://127.0.0.1:50051'
);

ATTACH '' AS my_server (
    TYPE AIRPORT,
    LOCATION 'grpc://127.0.0.1:50051'
);

SELECT *
FROM my_server.my_schema.orders
LIMIT 10;
```

### What Happens Internally

* Token is validated
* Accessible tables are resolved
* Unauthorized fields are masked
* Scopes are applied
<!--* RLA filters rows-->
* Data is streamed as Arrow batches

---

## Current Limitations & Roadmap

Current focus is **read-optimized analytical access**.

Planned improvements:

* Explicit DML support (INSERT / UPDATE / DELETE)
* Scope composition strategies
* Better schema introspection
* Cached connector reuse
* Declarative exposure policies

---

## Why Arrow Flight in Central-Set?

Arrow Flight allows Central-Set to function as a **governed data serving layer**:

* ETLX outputs become instantly queryable
* No file exports
* No duplication
* Strong access control
* Works with modern analytics stacks

It bridges **data engineering, governance, and analytics** — cleanly and safely.

## Application-Aware Mode (`arrow_flight_conf`)

Arrow Flight can optionally operate in **Application-Aware Mode** when the `arrow_flight_conf` field is defined.

### Example Configuration

```json
{
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

When this configuration is present:

* All tables are served **as if accessed through the `crud/read` API**
* The request is executed within the context of the declared application
* Only the database defined in the application (`db`) is eligible for exposure
* An additional governance layer is applied

---

## What Changes in Application-Aware Mode?

Instead of exposing tables directly from attached DuckDB sources, Arrow Flight:

1. Resolves the application (`app_id`)
2. Switches context to the application’s declared database
3. Applies the same internal logic used by:

```
/dyn_api/crud/read
```

This means:

* ✅ Full CRUD-layer access rules apply
* ✅ Row-Level Access (RLA) is enforced
* ✅ Field-level restrictions are enforced
* ✅ Application permissions are respected
* ✅ Business rules embedded in the CRUD layer are preserved

---

## Security Implications

Application-Aware Mode adds **another layer of security**:

| Layer                      | Direct DuckDB Mode  | Application-Aware Mode |
| -------------------------- | ------------------- | ---------------------- |
| Token validation           | ✅                  | ✅                    |
| Schema access              | ✅                  | ✅                    |
| Table access               | ✅                  | ✅                    |
| Field masking              | ✅                  | ✅                    |
| Scope filtering            | ✅                  | ✅                    |
| Row-Level Access (RLA)     | ❌                  | ✅                    |
| App-level permission model | ❌                  | ✅                    |
| CRUD business rules        | ❌                  | ✅                    |

Because only the database declared in the application configuration is accessible at this level:

* Cross-database access is prevented
* Exposure is limited to the **application’s data domain**
* Fine-grained governance is preserved

---

## When to Use Application-Aware Mode

Use this mode when:

* You want Arrow Flight to behave like a secure analytical gateway over your application database
* You rely heavily on:

  * RLA
  * App-specific permissions
  * CRUD-layer logic
* You need governance parity between REST and Flight

---

## When to Use Direct DuckDB Mode

Use standard (non-application-aware) mode when:

* Serving ETLX outputs
* Serving analytical datasets
* Attaching external data sources
* Building cross-database analytical layers

This mode is more flexible but bypasses CRUD-level business logic.

---

## Summary

If `arrow_flight_conf.app` is defined:

Arrow Flight behaves as a **high-performance analytical interface over the application’s database**, with full CRUD-level security enforcement.

If not defined:

Arrow Flight behaves as a **governed DuckDB analytical server**, controlled by schema, table, field, and scope configuration.

