---
weight: 3020
title: "Admin Model"
description: "Understanding the Central Set Admin Model and how model-driven architecture defines applications, permissions, APIs, and metadata."
icon: schema
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

{{% alert context="warning" text="**Caution** — Central Set is **production-ready and actively used**, but this documentation is still **under active development**. Large parts of the docs are **auto-generated**, evolving alongside the platform, and some sections may be incomplete, rough around the edges, or change frequently."  /%}}

The **Admin Model** defines the core metadata and system structure used by Central Set.

It is written as a **Markdown + YAML model**, which allows the platform to:

- Initialize databases
- Generate metadata
- Configure applications
- Manage permissions
- Define dashboards
- Control APIs and integrations

All of this is done **from a single model file**.

---

# 🧠 Model Philosophy

Central Set uses a **model-driven architecture**.

Instead of manually creating:

- tables
- menus
- permissions
- APIs
- dashboards

everything is defined inside a **model file**.

This model becomes the **single source of truth** for the application.

Advantages include:

- self-documenting database schema
- reproducible environments
- consistent deployments
- Git-versioned infrastructure
- LLM-friendly structure

The model can initialize or update the database with:

```bash
central-set --init --model admin_model.md
````

---

# 📦 Model Structure

A model is composed of several parts.

High-level structure:

```
Model Metadata
Applications (cs_app)
Tables
Security
Metadata Tables
Integrations
Dashboards
Jobs / Scheduling
```

Each section is written in **YAML blocks inside Markdown**, which allows the model to remain both **human readable** and **machine executable**.

---

# ⚙️ Model Metadata

The first section defines the **model configuration**.

Example:

```yaml
name: ADMIN
description: CS ADMIN Model
runs_as: MODEL
conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
```

Important fields:

| Field                   | Description              |
| ----------------------- | ------------------------ |
| `name`                  | Model name               |
| `conn`                  | Database connection      |
| `runs_as`               | Execution mode           |
| `create_all`            | Controls schema creation |
| `update_table_metadata` | Updates metadata tables  |
| `active`                | Enables the model        |

The connection usually references **environment variables**.

Example:

```
@DB_DRIVER_NAME
@DB_DSN
```

This keeps credentials **outside the model file**.

---

# 🧩 Application Definition (`cs_app`)

The `cs_app` section defines the **applications and menus** that appear inside the Central Set UI.

Example:

```yaml
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    tables:
      - dashboard
```

Applications can define:

* menus
* menu icons
* table mappings
* metadata configuration

Example structure:

```
App
 ├── Menu
 │   ├── Tables
 │   └── Configuration
```

From this, Central Set automatically generates:

* UI navigation
* CRUD APIs
* access control layers

---

# 🗄️ Table Definitions

Each table is defined using a YAML block.

Example:

```yaml
table: users
comment: Users
columns:
  user_id:
    type: integer
    pk: true
    autoincrement: true
```

Columns can define:

| Property   | Description       |
| ---------- | ----------------- |
| `type`     | SQL type          |
| `pk`       | Primary key       |
| `fk`       | Foreign key       |
| `nullable` | Allow null values |
| `default`  | Default value     |
| `unique`   | Unique constraint |
| `comment`  | Field description |

Example FK definition:

```
fk: "role.role_id"
```

This automatically generates relational constraints.

---

# 🔐 Security Model

Central Set implements **multi-layer security** directly inside the model.

Security layers include:

### Role Based Access

```
role
role_app
role_app_menu
role_app_menu_table
```

Controls:

* application access
* menu access
* table access
* CRUD permissions

---

### Row Level Access (RLA)

Tables:

```
row_level_access
role_row_level_access
```

RLA allows restricting records by:

* tenant
* department
* region
* customer

Example concept:

```
User A → only sees rows belonging to Tenant A
User B → only sees rows belonging to Tenant B
```

---

### Column Level Access

Table:

```
column_level_access
```

Controls which fields a user can:

* read
* create
* update

---

# 📊 Dashboard System

The model defines the **dashboard engine**.

Tables:

```
dashboard
dashboard_comment
```

Dashboards are defined using **Markdown configuration stored in the database**.

Example field:

```
dashboard_conf
```

This configuration can contain:

* SQL queries
* API calls
* ETLX pipelines
* visualization configuration

Dashboards therefore become **fully dynamic and programmable**.

---

# 🚀 Arrow Flight Integration

The Admin model also defines the **Arrow Flight server configuration**.

Tables:

```
arrow_flight
arrow_flight_table
arrow_flight_table_field
arrow_flight_table_scope
```

These allow exposing datasets through the **Arrow Flight protocol**.

Capabilities include:

* schema definition
* table discovery
* scoped datasets
* dynamic SQL execution

Example configuration fields:

| Field                 | Purpose               |
| --------------------- | --------------------- |
| `startup_sql`         | initialization SQL    |
| `main_sql`            | main execution query  |
| `table_discover_sql`  | dynamic table listing |
| `table_scan_tmpl_sql` | table scan template   |
| `shutdown_sql`        | cleanup logic         |

This allows Central Set to act as a **data service layer**.

---

# 🔑 Access Keys

Table:

```
access_key
```

Access keys allow:

* service integrations
* external APIs
* programmatic access

Each key includes:

* token
* expiration
* user association
* activation flag

These keys can be used in the API:

```
Authorization: Bearer <access_token>
```

---

# ⚙️ Environment Variables

Table:

```
env
```

Stores environment variables managed by Central Set.

Example usage:

```
@ENV_VAR_NAME
```

These variables can be used in:

* ETLX pipelines
* database connections
* integrations
* API configuration

---

# ⏱ Job Scheduling

Tables:

```
cron
cron_log
```

This allows Central Set to run scheduled jobs.

Example jobs:

* ETLX pipelines
* backups
* environment sync
* maintenance tasks

Example cron expression:

```
0 0 * * *
```

Jobs execute **Central Set APIs**.

Example:

```
etlx/name/[pipeline_name]
```

---

# 📜 System Metadata

Several tables store metadata used internally by Central Set:

```
table
table_schema
menu
menu_table
custom_table
custom_form
translate_table
translate_table_field
```

These allow the system to dynamically:

* generate CRUD APIs
* build UI forms
* create dashboards
* manage translations
* extend metadata

---

# 🧠 Why Everything Lives in the Model

By defining everything in a **single Markdown model**, Central Set achieves:

* Infrastructure as Code
* Self-documenting architecture
* Automated API generation
* Built-in security layers
* Dynamic UI generation

The model becomes the **foundation of the entire platform**.

---

# 🚀 Summary

The **Admin Model** defines the entire core system of Central Set, including:

* users
* roles
* applications
* menus
* security layers
* dashboards
* data services
* integrations
* job scheduling

All of this is controlled from a **single Markdown model file**.

This approach enables Central Set to function as a **metadata-driven backend platform** capable of dynamically generating:

* APIs
* dashboards
* data pipelines
* integrations