---
weight: 7080
title: "Dynamic API Reference"
description: "Dynamic API Reference"
icon: code
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

# Dynamic API Reference

{{% alert context="warning" %}}
**Important** — Central-Set uses a **dynamic API router**. Endpoints are not statically defined at compile time like traditional REST APIs. This document represents the **authoritative contract** for all supported dynamic routes and is derived from the internal `dyn_api.go` dispatcher.
{{% /alert %}}

---

## Overview

Central-Set exposes a single dynamic entrypoint:

```
POST /dyn_api/{ctrl}/{act}
```

Where:

| Parameter | Description                                          |
| --------- | ---------------------------------------------------- |
| `ctrl`    | Logical namespace (similar to a class or controller) |
| `act`     | Action or method executed within that namespace      |

Internally, each `{ctrl}/{act}` pair is mapped to a Go function via a `switch` statement. All requests and responses follow a **strict, uniform envelope**, enabling:

* Dynamic routing
* Consistent authentication & authorization
* Language-agnostic clients
* Easier UI auto-generation

---

## Authentication & Context Injection

All dynamic API endpoints:

* 🔐 **Require Authorization** (JWT access token)
* 📌 Use the `Authorization: Bearer <token>` header
* 👤 Automatically inject the authenticated **user identity**
* 🧠 Automatically resolve the **active app context**

If no app is explicitly selected, Central-Set defaults to:

```json
{
  "app_id": 1,
  "app": "ADMIN",
  "db": "ADMIN"
}
```

Row-level access, column-level access, and role-based permissions are enforced transparently.

---

## Request Envelope (Standard)

```json
{
  "lang": "en",
  "app": "ADMIN",
  "data": { }
}
```

| Field  | Description                    |
| ------ | ------------------------------ |
| `lang` | Response language (i18n aware) |
| `app`  | Target application context     |
| `data` | Action-specific payload        |

---

## Response Envelope (Standard)

```json
{
  "success": true,
  "msg": "Operation completed",
  "data": { }
}
```

| Field     | Description                      |
| --------- | -------------------------------- |
| `success` | Indicates request outcome        |
| `msg`     | Human-readable status message    |
| `data`    | Action-specific response payload |

---

# Auth Controller

## POST /dyn_api/auth/login

Authenticates a user and returns an access token.

### Request `data`

```json
{
  "username": "admin",
  "password": "secret"
}
```

### Behavior

* Validates credentials
* Issues a JWT token
* Binds token to user identity and roles

### Response `data`

```json
{
  "token": "<jwt>",
  "user": {
    "id": 1,
    "username": "admin",
    "roles": ["admin"]
  }
}
```

---

## POST /dyn_api/auth/refresh_token

Refreshes an existing access token.

### Request `data`

```json
{
  "token": "<expired-or-near-expiry-token>"
}
```

### Behavior

* Validates refresh eligibility
* Issues a new token

### Response `data`

```json
{
  "token": "<new-jwt>"
}
```

---

# User Controller

## POST /dyn_api/user/list

Returns a paginated list of users visible to the caller.

### Request `data`

```json
{
  "page": 1,
  "page_size": 50,
  "filters": {
    "active": true
  }
}
```

### Behavior

* Applies role-based visibility
* Applies row-level access rules

### Response `data`

```json
{
  "items": [ { "id": 1, "username": "admin" } ],
  "total": 1
}
```

---

## POST /dyn_api/user/create

Creates a new user.

### Request `data`

```json
{
  "username": "new_user",
  "password": "secret",
  "roles": ["viewer"]
}
```

### Behavior

* Hashes password
* Assigns roles
* Validates permissions

### Response `data`

```json
{
  "id": 42
}
```

---

# Table Controller

## POST /dyn_api/table/list

Lists registered tables for the active app.

### Request `data`

```json
{
  "schema": "main"
}
```

### Behavior

* Applies table-level permissions
* Filters hidden or restricted tables

### Response `data`

```json
{
  "tables": ["orders", "customers"]
}
```

---

## POST /dyn_api/table/query

Executes a controlled query against a table.

### Request `data`

```json
{
  "table": "orders",
  "limit": 100
}
```

### Behavior

* Enforces column-level access
* Enforces row-level access
* Executes via DuckDB

### Response `data`

```json
{
  "rows": [ { "order_id": 1 } ]
}
```

---

# ETLX Controller

## POST /dyn_api/etlx/run

Executes an ETLX pipeline.

### Request `data`

```json
{
  "pipeline_id": 7,
  "params": {
    "date": "2025-01-01"
  }
}
```

### Behavior

* Resolves pipeline definition
* Injects secrets via environment variables
* Executes inside DuckDB

### Response `data`

```json
{
  "run_id": "etl-20250101-001"
}
```

---

# Arrow Flight Controller

## POST /dyn_api/arrow_flight/list

Lists exposed Arrow Flight schemas.

### Response `data`

```json
{
  "schemas": ["analytics", "sales"]
}
```

---

## POST /dyn_api/arrow_flight/expose

Creates or updates an Arrow Flight schema exposure.

### Request `data`

```json
{
  "name": "analytics",
  "startup_sql": "INSTALL parquet;",
  "main_sql": "ATTACH 'data.db';",
  "shutdown_sql": "DETACH analytics;"
}
```

### Behavior

* Validates SQL lifecycle
* Registers schema

### Response `data`

```json
{
  "success": true
}
```

---

## Notes on Dynamic API Design

* All endpoints are **POST-only** by design
* URL defines behavior, body defines intent
* Strong typing is enforced at runtime
* The same contract powers:

  * REST clients
  * UI auto-forms
  * CLI tooling

This design allows Central-Set to evolve rapidly **without breaking clients**, while maintaining strong governance and security guarantees.
