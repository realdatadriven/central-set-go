---
weight: 7090
title: "Arrow Flight"
description: "Apache Arrow Flight Support API"
icon: arrow
date: 2025-12-16T01:04:15+00:00
lastmod: 2026-01-29T00:00:00+00:00
draft: false
images: []
---

## Arrow Flight Support

{{% alert context="warning" text="**Important** – Central-Set is production-ready, but this documentation is still under active development. Parts of the Arrow Flight subsystem are configuration-driven and evolve as access-control and governance features mature." /%}}

### Overview

Central-Set provides **Apache Arrow Flight** support through  
[airport-go](https://github.com/hugr-lab/airport-go), exposing analytical datasets via a **high-performance, governed, and dynamically scoped interface**.

Arrow Flight is primarily designed to serve:

- **ETLX outputs**
- **DuckDB-backed analytical views**
- **Externally attached datasources**

…while enforcing the **same security, access control, and multi-tenant rules** used throughout Central-Set.

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

- If **no scopes are defined** → all rows are eligible (subject to RLA)
- If **one or more scopes exist**:
  - The token **must have access to at least one scope**
  - Otherwise **no data is returned**
- Multiple scopes are **OR-combined**

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
7. **Row-Level Access (RLA)**

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
