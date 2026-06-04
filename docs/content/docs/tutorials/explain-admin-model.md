---
weight: 3021
title: "Explaining ADMIN Model"
description: "Central Set's core administrative model — defines users, roles, permissions, menus, applications, security layers, dashboards, scheduled jobs, Arrow Flight exposure, and dynamic UI generation."
icon: schema
date: 2026-03-20
lastmod: 2026-03-20
draft: false
images: []
---

{{% alert context="warning" text="**Note** — Central Set is production-ready and under active development. This documentation page is **auto-generation friendly** and will evolve together with the platform. Some sections may still be expanding or receive refinements." /%}}

The **ADMIN_MODEL** is the **foundational metadata model** of Central Set.

It defines:

- the security model (users • roles • permissions • row/column access)
- application & menu structure
- UI generation rules (forms • datatables • layouts • custom components)
- dashboards
- scheduled jobs
- environment variables
- Arrow Flight / data service exposure
- access tokens
- translation & custom metadata

Everything lives inside **one Markdown + YAML file** — the single source of truth.

```bash
central-set --init --model admin_model.md
```

or

```bash
central-set --update-metadata --model admin_model.md
```

---

## 🧠 Philosophy of the ADMIN Model

Central Set follows a **strong model-driven architecture**.

Instead of hand-writing:

- database tables
- CRUD endpoints
- RBAC rules
- forms & datagrids
- navigation menus
- dashboards
- background jobs
- data service endpoints

… you declare the desired system in **declarative YAML blocks** inside Markdown.

The runtime then:

- creates / updates schema
- generates metadata entries
- builds REST-like APIs
- expose every table as OData v4 API
- renders admin UI automatically
- enforces multi-layer security
- exposes Arrow Flight endpoints (if configured)

---

## 📦 High-Level Structure

```text
ADMIN_MODEL.md
├── Model metadata (name, connection, behavior flags)
├── cs_app                → applications & menu groups
├── Core tables           → lang, role, users, user_role, app, menu, ...
├── Security tables       → role_app, role_app_menu, role_app_menu_table, row_level_access, ...
├── UI & metadata tables  → table, table_schema, menu_table, custom_table, custom_form, ...
├── Dashboard tables      → dashboard, dashboard_comment
├── Arrow Flight          → flight_schema, flight_schema_table, flight_schema_table_scope, ...
├── Jobs & logging        → cron, cron_log, user_log
├── Integration & utils   → access_key, env, translate_table, translate_table_field
```

---

## ⚙️ Model Header

```yaml {linenos=table}
name: ADMIN
description: CS ADMIN Model
runs_as: MODEL
conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
```

Most important flags:

| Field                   | Typical value     | Meaning                                      |
|-------------------------|-------------------|----------------------------------------------|
| `conn`                  | `@DB_…`           | Database connection string (env vars)        |
| `create_all`            | `checkfirst`      | Create tables if missing                     |
| `update_table_metadata` | `true`            | Refresh `table` & `table_schema` entries     |
| `active`                | `true`            | Load this model on startup                   |

---

## 🧩 Application & Menu Structure (`cs_app`)

```yaml {linenos=table}
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    active: true
    tables:
      - dashboard
  Admin:
    menu_icon: user-group
    menu_order: 2
    tables:
      - app
      - menu
      - role
      - users
      ## …
```

This block automatically creates:

- sidebar menu groups
- icons & ordering
- associated tables (shown as sub-items or default views)

---

## 🔐 Security Layers

Central Set implements **four levels** of access control:

1. **Application access** — `role_app`
2. **Menu access**      — `role_app_menu`
3. **Table / CRUD access** — `role_app_menu_table` (create/read/update/delete/share)
4. **Row Level Access (RLA)** — `role_row_level_access` / `row_level_access`

Typical flow:

```
User → Role(s) → App permission → Menu permission → Table CRUD permission → Row filter
```

---

## 🖥️ Automatic UI Generation

Every table can carry **UI hints**:

```yaml {linenos=table}
columns:
  username:
    type: varchar(50)
    unique: true
    nullable: false
    form_display: true
    table_display: true
    form_size: 4
    order: 1
```

Common UI properties:

| Property             | Effect                                         |
|----------------------|------------------------------------------------|
| `form_display`       | Show in create/update forms                    |
| `table_display`      | Show as column in list view                    |
| `form_size`          | Bootstrap column width (1–12)                  |
| `form_long_text`     | Use textarea instead of input                  |
| `form_code`          | markdown, json, sql, … → code editor mode      |
| `form_att`           | File / attachment upload                       |

Form layout control:

```yaml {linenos=table}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
```

---

## 🎛️ Custom UI Extensions

Via `table_extra_options`:

```yaml {linenos=table}
table_extra_options:
  - {size: 12, component: EvidenceDash, label: dashboard, intercept_r: true}
  - {size: 12, component: AdminApps, label: permissions, icon: key, pop_up: true}
```

Popular built-in components:

- `EvidenceDash` — renders markdown dashboards
- `AdminApps`   — permission matrix / role editor
- `AccessKey`   — token generation & display

---

## 📊 Dashboards

Dashboards are stored as markdown + configuration:

```yaml {linenos=table}
dashboard_conf: { type: text, form_code: markdown }
```

The `EvidenceDash` component interprets this field and can embed:

- SQL query results
- charts
- ETLX results
- external API data
- markdown content

---

## 🚀 Arrow Flight / Data Service Exposure

```yaml {linenos=table}
flight_schema:
  flight_schema: adm
  startup_sql:   "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE);"
  main_sql:      "USE adm;"
```

Allows Central Set to act as an **Arrow Flight server**, exposing selected tables/scopes to:

- Python (pyarrow)
- R
- BI tools with Arrow Flight support
- Datafusion, Polars, …

---

## ⏱️ Scheduled Jobs (cron)

```yaml {linenos=table}
cron:
  cron: "0 0 * * *"
  api:  "etlx/name/daily-backup"
  active: true
```

Central Set comes with a built-in cron runner that calls internal APIs.

---

## 🔑 Access Keys & Environment Variables

- `access_key` → bearer tokens for machine-to-machine access
- `env` → centrally managed environment variables (usable via `@ENV_NAME`)

---

## Summary — Why one big model?

By keeping **everything** in one declarative Markdown/YAML file you get:

- Git-versioned infrastructure
- Reproducible environments
- Self-documenting schema
- Automatic security & UI generation
- Easy auditing & review
- LLM-assisted modifications possible

The **ADMIN_MODEL** is effectively the **constitution** of your Central Set instance.

Feel free to extend it — every new table, role, menu group or dashboard added here immediately becomes part of the live system.

[More Details](https://github.com/realdatadriven/central-set-go/blob/main/admin_model.md)
