---
weight: 3020
title: "Admin Model"
description: "Understanding the Central Set Admin Model and how model-driven architecture defines applications, permissions, APIs, metadata, and UI."
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
- Generate **UI forms and datatables**

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
- forms
- datatables

everything is defined inside a **model file**.

This model becomes the **single source of truth** for the application.

Advantages include:

- self-documenting database schema
- reproducible environments
- consistent deployments
- Git-versioned infrastructure
- LLM-friendly structure
- automatic UI generation

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
UI Definitions
Integrations
Dashboards
Jobs / Scheduling
```

Each section is written in **YAML blocks inside Markdown**, which allows the model to remain both **human readable** and **machine executable**.

---

# ⚙️ Model Metadata

The first section defines the **model configuration**.

Example:

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

```yaml {linenos=table}
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

```yaml {linenos=table}
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

# 🧾 UI Generation from the Model

Central Set models now also define **UI behavior**, including:

* form fields
* datatable columns
* layout configuration
* component hooks
* custom UI extensions

This allows the platform to automatically generate **admin interfaces directly from the schema**.

---

# 🧩 Field UI Configuration

Fields can define how they appear in the UI.

Example:

```yaml {linenos=table}
dashboard: { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true }
```

Available options include:

| Property        | Description                            |
| --------------- | -------------------------------------- |
| `form_display`  | Display field in forms                 |
| `table_display` | Display field in datatable             |
| `form_code`     | Enables code editor mode               |
| `form_sizelg`   | Form layout size (large screens)       |
| `form_sizexl`   | Form layout size (extra large screens) |

Example with code editor:

```yaml {linenos=table}
dashboard_conf:
  type: text
  form_display: true
  form_code: markdown
```

This enables a **Markdown editor** for the dashboard configuration.

---

# 🧩 Form Layout Configuration

Forms can define layout options using `form_layout`.

Example:

```yaml {linenos=table}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

Options include:

| Property        | Description             |
| --------------- | ----------------------- |
| `tabs_steps`    | Enables tab-based forms |
| `form_in_popup` | Allows popup forms      |
| `size`          | Form container width    |

This allows the UI to automatically build **structured form layouts**.

---

# 📊 Datatable Configuration

Datatable behavior can also be controlled.

Example options:

* visible columns
* ordering
* component extensions

These configurations are automatically used by the frontend datatable components.

---

# 🧩 Extra UI Options and Component Hooks

Tables can also define **extra UI behavior** through `table_extra_options`.

Example from the **Dashboard table**:

```yaml {linenos=table}
table_extra_options:
  - {size: 12, component: EvidenceDash, label: dashboard, intercept_r: true}
```

This allows injecting **custom UI components** into the interface.

Example capabilities include:

* custom visual components
* dashboard renderers
* data visualizations
* intercepting CRUD operations

---

# 🔌 CRUD Interceptors

Components can intercept CRUD actions.

Example:

```
intercept_r
```

This means the custom component will intercept **read operations**.

Possible interception hooks include:

| Hook          | Purpose                   |
| ------------- | ------------------------- |
| `intercept_r` | intercept read operations |
| `intercept_c` | intercept create          |
| `intercept_u` | intercept update          |
| `intercept_d` | intercept delete          |

This enables **advanced UI extensions without modifying backend logic**.

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
flight_schema
flight_schema_table
flight_schema_table_field
flight_schema_table_scope
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
* Automatic UI generation
* Dynamic dashboard configuration

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
* UI forms
* datatables
* data services
* integrations
* job scheduling

All of this is controlled from a **single Markdown model file**.

This approach enables Central Set to function as a **fully metadata-driven backend platform** capable of dynamically generating:

* APIs
* dashboards
* data pipelines
* UI applications
* integrations

