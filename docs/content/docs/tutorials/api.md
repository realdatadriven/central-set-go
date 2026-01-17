---
weight: 7080
title: "Dynamic API Reference"
description: "D"
icon: code
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

# Dynamic API Reference

Central Set Go exposes a **dynamic, context-driven API**.
All operations — including read-only ones — use `POST` requests to allow rich payloads, permissions, and future extensibility.

> ⚠️ **Design note**
> Even read operations (`apps`, `menus`, `tables`, `crud.read`) use `POST` by design.

---

## Authentication

### Login

Authenticates a user and returns a JWT token used by all subsequent requests.

**Endpoint**

```
POST /dyn_api/login/login
```

**Headers**

```
Content-Type: application/json
```

**Request**

```json
{
  "lang": "pt",
  "data": {
    "username": "root",
    "password": "1234"
  }
}
```

**Response (simplified)**

```json
{
  "success": true,
  "token": "<JWT_TOKEN>",
  "user": {
    "user_id": 1,
    "username": "root",
    "role_id": 1
  }
}
```

---

### Check Token

Validates an existing JWT token.

**Endpoint**

```
POST /dyn_api/login/chk_token
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en"
}
```

---

### Change Password

Allows an authenticated user to change their password.

**Endpoint**

```
POST /dyn_api/login/alter_pass
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en",
  "controller": "login",
  "action": "alter_pass",
  "data": {
    "username": "root",
    "password": 12344654,
    "new_password": 1234
  }
}
```

---

## Applications & Navigation

### Load Applications

Returns the list of applications available to the user.

**Endpoint**

```
POST /dyn_api/admin/apps
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en"
}
```

---

### Load Menus

Loads the menus for a given application.

**Endpoint**

```
POST /dyn_api/admin/menu
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en",
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "app_desc": "Admin",
    "version": "1.0.0",
    "db": "ADMIN"
  }
}
```

Menus define:

* Sidebar navigation
* Icons
* Table grouping

---

### Load Tables (Metadata)

Returns metadata about tables available in the application.

**Endpoint**

```
POST /dyn_api/admin/tables
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en",
  "data": {
    "table": "app"
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

## CRUD Operations

All CRUD operations share the same **dynamic payload structure**.

### Read

**Endpoint**

```
POST /dyn_api/crud/read
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request**

```json
{
  "lang": "en",
  "data": {
    "table": "menu",
    "limit": 11,
    "offset": 0,
    "filters--": [
      { "field": "lang_id", "cond": "=", "value": 3 }
    ],
    "order_by": [
      { "field": "lang_id2", "order": "DESC" }
    ]
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

### Create

**Endpoint**

```
POST /dyn_api/crud/create
```
**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```
**Request**

```json
{
  "lang": "en",
  "data": {
    "table": "lang",
    "data": {
      "lang": "es4",
      "lang_desc": "Português"
    }
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

### Update

**Endpoint**

```
POST /dyn_api/crud/update
```
**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```
**Request**

```json
{
  "lang": "pt",
  "data": {
    "table": "lang",
    "data": {
      "lang_id": 3,
      "lang": "pt_BR",
      "lang_desc": "Português Brasil"
    }
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

### Delete

Supports soft delete and permanent delete.

**Endpoint**

```
POST /dyn_api/crud/delete
```
**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```
**Request**

```json
{
  "lang": "en",
  "data": {
    "table": "lang",
    "data": {
      "lang_id": 3,
      "exclude": true,
      "permanently--": true
    }
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

### Query (Raw SQL)

Executes a raw SQL query.

**Endpoint**

```
POST /dyn_api/crud/query
```
**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```
**Request**

```json
{
  "lang": "en",
  "data": {
    "db": "DB.duckdb",
    "query": "select * from \"table_name\"",
    "limit": 10,
    "offset": 0
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

---

## File Uploads

Uploads files for temporary or permanent storage.

**Endpoint**

```
POST /upload
```

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: multipart/form-data
```

**Form Data**

| Field  | Description      |
| ------ | ---------------- |
| `lang` | Language         |
| `tmp`  | Temporary upload |
| `path` | Destination path |

---

## ETL

### Extract from File + Attach Database

**Endpoint**

```
POST /dyn_api/etlx/run
```

**Description**

Extracts data from files (CSV, XLSX, DuckDB) and optionally attaches databases for transformation.

This endpoint powers **ETLX pipelines** in the UI.

**Headers**

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request (simplified)**

```json
{
  "lang": "en",
  "data": {
    "conf": "# ..."
  },
  "app": {
    "app_id": 2,
    "app": "ETLX",
    "db": "ETLX"
  }
}
```


---

## Conceptual Summary

| Area           | Covered By              |
| -------------- | ----------------------- |
| Authentication | `/dyn_api/login/*`      |
| Applications   | `/dyn_api/admin/apps`   |
| Menus          | `/dyn_api/admin/menu`   |
| Metadata       | `/dyn_api/admin/tables` |
| CRUD           | `/dyn_api/crud/*`       |
| Files          | `/upload`               |
| ETL            | `/dyn_api/etl/*`        |

