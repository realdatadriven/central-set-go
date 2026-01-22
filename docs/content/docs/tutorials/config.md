---
weight: 7091
title: "Configuration"
description: "Configuration options for Central Set"
icon: code
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Configuration

{{% alert context="warning" %}}
**Caution** — Central Set is production-ready, but this documentation is still under active development and partially auto-generated. Configuration options may expand as new capabilities are added.
{{% /alert %}}

---

Central Set is primarily configured using environment variables (`.env`).
This page documents the **core and most commonly required settings** to get the platform running in development or production.

> 💡 Tip: You do **not** need to configure everything. A minimal `.env` is usually enough to start.

---

## Minimal Required Configuration

The following variables are required for **any** Central Set instance (if not set, defaults are used):

```env
BASE_URL=http://localhost:4444
HTTP_PORT=4444

DB_DRIVER_NAME=sqlite3
DB_DSN=database/ADMIN.db

JWT_SECRET_KEY=change-me
COOKIE_SECRET_KEY=change-me
SECRET_KEY=change-me
```

These define:

* Where the API/UI is exposed
* How Central Set connects to its **admin database**
* How authentication tokens and cookies are secured

---

## Server & Base URL

```env
BASE_URL=http://localhost:4444
HTTP_PORT=4444
```

| Variable    | Description                                           |
| ----------- | ----------------------------------------------------- |
| `BASE_URL`  | Public base URL used for redirects, links, and tokens |
| `HTTP_PORT` | Port the HTTP server listens on                       |

---

## Authentication (Admin Login)

Central Set supports **basic authentication** for the initial admin access.

```env
BASIC_AUTH_USERNAME=root
BASIC_AUTH_HASHED_PASSWORD=$2a$10$...
```

* Password **must be bcrypt-hashed**
* Used mainly for bootstrap / admin scenarios

---

## Security & Tokens

```env
JWT_SECRET_KEY=mhaitpm4v3mesosefepyupo6qzpbvidc
COOKIE_SECRET_KEY=f2rkbev2yxhk5viz77ok4rxfip6npjpm
SECRET_KEY=openssl-generated-hex
ALGORITHM=HS256
TOKEN_EXPIRE_HOURS=1440
```

| Variable             | Description                |
| -------------------- | -------------------------- |
| `JWT_SECRET_KEY`     | Signs API access tokens    |
| `COOKIE_SECRET_KEY`  | Encrypts session cookies   |
| `SECRET_KEY`         | Internal cryptographic key |
| `ALGORITHM`          | JWT signing algorithm      |
| `TOKEN_EXPIRE_HOURS` | Token expiration time      |

⚠️ **Always change these values in production.**

---

## Database Configuration

### SQLite (Default)

```env
DB_DRIVER_NAME=sqlite3
DB_DSN=database/ADMIN.db
```

### PostgreSQL Example

```env
DB_DRIVER_NAME=postgres
DB_DSN=user=postgres password=1234 dbname=ADMIN host=localhost port=5432 sslmode=disable
```

```env
DB_AUTOMIGRATE=false
```

| Variable         | Description                                                   |
| ---------------- | ------------------------------------------------------------- |
| `DB_DRIVER_NAME` | Database driver (`sqlite3`, `postgres`, `mysql`, `sqlserver`) |
| `DB_DSN`         | Database connection string                                    |
| `DB_AUTOMIGRATE` | Automatically migrate schema on startup                       |

---

## Static Files & Uploads

```env
STATIC=static
ASSETS=static/assets
UPLOAD=static/uploads
UPLOAD_SIZE=104857600
```

Controls:

* Static UI assets
* File uploads
* Maximum upload size (bytes)

---

## Feature Flags & Core Tables

```env
ENABLE_APP=app,role_app,role_app_menu,role_app_menu_table
ENABLE_USER=user_role,column_level_access,row_level_access

CORE_TABLES=lang,role,user,users,user_role,app,menu,table,menu_table,...
```

These variables control:

* Which **core modules** are enabled
* Which tables are managed by the system

> Useful for lightweight or embedded deployments.

---

## Logging & Auditing

```env
ACTIONS_NOT_TO_LOG=apps,menu,tables,chk_token
BROADCAST_CHANGES=create_update,create,update,delete,c,u,d
```

Controls what actions are:

* Logged
* Broadcast to clients (UI refresh, events)

---

## Email & Notifications (Optional)

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=no.reply@gmail.com
SMTP_PASSWORD=******
SMTP_FROM=APP <no.reply@gmail.com>
```

Used for:

* Notifications
* Password resets
* System alerts

---

## Object Storage (S3 Compatible)

```env
USE_S3_STORAGE=false

S3_BUCKET=uploads
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_REGION=us-east-1
AWS_ENDPOINT=127.0.0.1:3000
S3_FORCE_PATH_STYLE=true
S3_SKIP_SSL_VERIFY=true
```

Supports:

* AWS S3
* MinIO
* Other S3-compatible backends

---

## ETLX & Data Pipelines

```env
ETLX_DEBUG_QUERY=true
```

Enables verbose query logging for ETLX workflows.

---

## Export & Background Jobs

```env
EXPORT_CONN_TIMEOUT=3600
EXPORT_ENC_KEY=change-me
EXPORT_ADMIN_DB_TABLES=app,table,user_log,...
```

Used by:

* Export jobs
* Governance artifacts
* Scheduled ETLX tasks

---

## LDAP Authentication (Optional)

```env
USE_LDAP_AUTH=false
LDAP_URL=ldap://localhost:1389
LDAP_BIND_USER=cn=admin,dc=example,dc=com
LDAP_PASSWORD=admin
LDAP_BASE_DN=dc=example,dc=com
LDAP_SKIP_VERIFY_CERT=true
```

Allows Central Set to authenticate users against LDAP directories.

---

## Arrow Flight (Analytics Sharing)

```env
ENABLE_ARROW_FLIGHT=false
ARROW_FLIGHT_ADDR=:8815
ARROW_FLIGHT_TLS=false
```

When enabled, Central Set exposes datasets via **Apache Arrow Flight** for:

* Data science
* BI tools
* High-performance analytics

---

## Recommended Production Checklist

* [ ] Change all secret keys
* [ ] Enable HTTPS
* [ ] Disable debug flags
* [ ] Use PostgreSQL instead of SQLite
* [ ] Enable row-level access rules
* [ ] Configure backups

---

## Summary

Central Set configuration is:

* Environment-driven
* Explicit and auditable
* Flexible enough for embedded or large deployments

Most installations only require **10–15 variables** to get started.
The rest unlock advanced features when you need them.
