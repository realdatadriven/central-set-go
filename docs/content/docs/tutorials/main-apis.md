---
weight: 7075
title: "Main APIs"
description: "A practical overview of the main Central Set Go API groups and request patterns."
icon: api
date: 2026-11-07T00:00:00+00:00
lastmod: 2026-11-07T00:00:00+00:00
draft: false
images: []
---

# Main APIs

Central Set Go exposes APIs for authentication, application metadata, data access, file uploads, and ETLX workflows. This guide summarizes the main endpoints represented in the C7 OpenPost collection and the request conventions used by the server.

Set `BASE_URL` to the address of your Central Set Go server. The examples use `http://localhost:4444`; replace it with your deployment URL.

## Request conventions

The dynamic API uses `POST` requests with JSON bodies, including for read-only operations. Most requests share a payload with a language and a `data` object. Data operations also identify the application and its database:

```json
{
  "lang": "en",
  "data": {},
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

Send JSON requests with:

```http
Content-Type: application/json
```

After login, send the returned JWT on authenticated requests:

```http
Authorization: Bearer <JWT>
```

Authentication and password workflows are described in [Authentication & Access](./auth.md).

## Main endpoint groups

| API group | Main endpoints | Purpose |
| --- | --- | --- |
| Authentication | `POST /dyn_api/login/login`, `/dyn_api/login/chk_token`, `/dyn_api/login/alter_pass`, `/dyn_api/login/recover_pass`, `/dyn_api/login/reset_pass`, `/dyn_api/login/signup`, `/dyn_api/login/validate_code` | Sign in, validate a token, and manage account credentials and verification. |
| Applications and metadata | `POST /dyn_api/admin/apps`, `/dyn_api/admin/menu`, `/dyn_api/admin/tables` | List accessible applications, load an application's menus, and inspect table metadata. |
| Dynamic CRUD | `POST /dyn_api/crud/read`, `/dyn_api/crud/create`, `/dyn_api/crud/update`, `/dyn_api/crud/delete`, `/dyn_api/crud/query` | Read or change records and run a database query through the dynamic API. |
| REST CRUD | `GET`, `POST`, `PUT`, `PATCH`, and `DELETE /crud/{db}/{table}[/{id}]` | Use HTTP methods and path parameters for record operations. |
| ETLX | `POST /dyn_api/etlx/config`, `/dyn_api/etlx/run`, `/dyn_api/etlx/parse_run`, `/dyn_api/etlx/run_by_name`, `/dyn_api/etlx/query_md` | Parse ETLX configuration, run a pipeline, or query ETLX markdown. |
| Uploads | `POST /upload` | Upload a file using multipart form data. |
| OData | `GET /odata/{db}`, `/odata/{db}/$metadata`, `/odata/{db}/{table}` | Browse database metadata and read table data using OData-style query options. |

Dynamic API endpoints use `/dyn_api/{controller}/{action}`. For example, the server accepts both `/dyn_api/login/login` and the `/dyn_api/login/index` login alias used in the OpenPost collection. Use `/dyn_api/login/login` in new clients.

## 1. Authenticate

Login does not require an existing JWT. Replace the credential placeholders with an account for your instance:

```sh
curl -X POST "$BASE_URL/dyn_api/login/login" \
  -H "Content-Type: application/json" \
  -d '{
    "lang": "en",
    "data": {
      "username": "<USERNAME>",
      "password": "<PASSWORD>"
    }
  }'
```

On success, the response includes a `token` and user details. Use that token as a Bearer token in subsequent requests. Login, signup, password recovery, reset, and verification have their own request details in [Authentication & Access](./auth.md).

## 2. Load applications and metadata

List applications available to the authenticated user:

```sh
curl -X POST "$BASE_URL/dyn_api/admin/apps" \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"lang":"en"}'
```

Use the selected application's context to request its menus or table metadata:

```json
{
  "lang": "en",
  "data": {
    "table": "menu"
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

Send this body to `POST /dyn_api/admin/menu` to load menus; use `POST /dyn_api/admin/tables` to load table metadata. The `app` values should describe an application the current user can access.

## 3. Read and change records

Dynamic CRUD requests use `data.table` to select a table. A read may include pagination, filters, selected fields, and sort order:

```json
{
  "lang": "en",
  "data": {
    "table": "menu",
    "limit": 25,
    "offset": 0,
    "filters": [
      {
        "field": "app_id",
        "cond": "=",
        "value": 1
      }
    ],
    "order_by": [
      {
        "field": "menu_id",
        "order": "ASC"
      }
    ]
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

Send the body to `POST /dyn_api/crud/read`. The canonical filter key is `filters`; use `fields` for field selection and `order_by` for sorting.

Create and update requests use the same outer structure, with the record fields nested under `data.data`. For example:

```json
{
  "lang": "en",
  "data": {
    "table": "lang",
    "data": {
      "lang": "fr",
      "lang_desc": "Français"
    }
  },
  "app": {
    "app_id": 1,
    "app": "ADMIN",
    "db": "ADMIN"
  }
}
```

Send that body to `POST /dyn_api/crud/create`. For `POST /dyn_api/crud/update`, include the record's primary key in the nested `data.data` object. `POST /dyn_api/crud/delete` accepts the table and record identifier in the same nested structure; consult your application's deletion policy before using permanent-delete options.

The REST CRUD routes provide an alternative for clients that prefer HTTP methods:

| Operation | Request |
| --- | --- |
| List records | `GET /crud/{db}/{table}` |
| Get a record | `GET /crud/{db}/{table}/{id}` |
| Create a record | `POST /crud/{db}/{table}` |
| Replace or update a record | `PUT /crud/{db}/{table}/{id}` |
| Partially update a record | `PATCH /crud/{db}/{table}/{id}` |
| Delete a record | `DELETE /crud/{db}/{table}/{id}` |

`POST /dyn_api/crud/query` is also available for database queries. It is a privileged operation: only send trusted query text and databases that the caller is authorized to access.

For the full dynamic request reference, see [Dynamic API Reference](./api.md).

## 4. Run ETLX workflows

The ETLX dynamic endpoints accept JSON payloads containing a `data` object. A configuration parse request, for example, supplies the ETLX markdown as `data.conf`:

```json
{
  "lang": "en",
  "data": {
    "conf": "<ETLX markdown configuration>"
  },
  "params": {}
}
```

Send the body to `POST /dyn_api/etlx/config`. Use `/dyn_api/etlx/run` to execute a supplied configuration, `/dyn_api/etlx/parse_run` to parse and execute it, or `/dyn_api/etlx/run_by_name` to run a saved workflow. Execution of client-supplied configurations is controlled by the server's `ETLX_ALLOW_CLI_CONFIG` setting; ask the instance administrator if these calls are disabled.

## 5. Upload files

`POST /upload` uses `multipart/form-data` rather than a JSON body. Include a `file` part; `lang` and `tmp` are optional form fields. Set `tmp` to `true` for a temporary upload:

```sh
curl -X POST "$BASE_URL/upload" \
  -H "Authorization: Bearer <JWT>" \
  -F "file=@./data.csv" \
  -F "lang=en" \
  -F "tmp=true"
```

The response reports whether the upload succeeded and includes the resulting file name when available.

## 6. Read data with OData

OData routes use `GET` and query-string options rather than the dynamic JSON envelope:

```http
GET /odata/ADMIN/menu?$filter=app_id%20eq%201&$select=menu_id,menu_name&$top=25
Authorization: Bearer <JWT>
```

The implementation supports common options including `$filter`, `$select`, `$top`, `$skip`, and `$orderby`. OData endpoints are available when `ENABLE_ODATA=true` is set on the server. See [OData v4 API](./odata-v4.md) for additional examples and metadata routes.

## Responses

Dynamic endpoints return JSON with a `success` field and, when relevant, a `msg` or result data. Check both the HTTP status and the JSON `success` value when handling a response; some application-level errors are represented in the JSON body.
