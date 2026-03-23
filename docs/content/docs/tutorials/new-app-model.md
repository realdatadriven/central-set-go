---
weight: 3022
title: "Creating a New App / Data Model"
description: "How to add a new application or domain-specific data model in Central Set — with automatic support for UI, dashboards, ETLX pipelines, cron jobs, role-based access, row-level security, and more."
icon: schema
date: 2026-03-20
lastmod: 2026-03-20
draft: false
images: []
---

{{% alert context="warning" text="**Central Set is production-ready** and this pattern is actively used in real deployments. The documentation evolves together with the platform — expect refinements as new features (ETLX enhancements, better notebook support, etc.) are added." /%}}

Central Set is built around **model-driven architecture**.

You can (and usually should) create **separate models** for each business domain or application instead of putting everything into the `ADMIN` model.

Benefits of dedicated app models:

- Clear **separation of concerns**  
- Independent lifecycle & versioning  
- Easier team collaboration (different teams own different models)  
- Automatic generation of:  
  • CRUD UI + forms + datatables  
  • Role-based access (RBAC)  
  • Row-level security (RLA)  
  • Dashboards  
  • ETLX data pipelines  
  • Scheduled jobs (cron)  
  • Notebook-style analysis  
  • Arrow Flight exposure (optional)

---

## Recommended Pattern — One Model per App / Domain

Create a new Markdown file (e.g. `sales_model.md`, `inventory_model.md`, `finance_model.md`) with this structure:

1. Model metadata header  
2. `cs_app` block → defines the menu group(s)  
3. Domain-specific tables  
4. (Optional) Re-use or reference shared tables (dashboards, notebooks, cron, …)  
5. (Optional) Seed data / default records  

Typical file name convention:  
`{UPPERCASE_APP_NAME}_MODEL.md`  
(e.g. `SALES_MODEL.md`, `HR_MODEL.md`, `ETLX_MODEL.md`)

After placing this file in the models directory and running

```bash
central-set --init --model sales_model.md
```

You immediately get:

- New sidebar menu group **Sales**  
- CRUD interfaces for all tables  
- Automatic form & datatable generation  
- Permission management (via ADMIN → roles)  
- Ability to add row-level filters per role  
- Dashboard tab (reusing the shared `dashboard` table)  
- Possibility to create ETLX pipelines that read/write SALES tables  
- Cron jobs that run on SALES data  
- Notebooks for ad-hoc analysis

---

## Re-using Shared Features (Dashboards • ETLX • Cron • Notebooks)

Most models **do not redefine** `dashboard`, `etlx`, `cron`, `notebook` tables.

Instead they:

1. Add their own tables to an existing menu group (e.g. Dashboards, ETLX)  
2. Or create a new menu group that references shared tables

Now users with access to the **ETLX** menu can create pipelines like:

- Daily sales aggregation → `SALES.sales_kpi`  
- Order sync from external API → `SALES.order`  
- Customer segmentation ETL

Same logic applies to **cron**, **notebook**, **dashboard**.

---

## When to Extend ADMIN vs Create New Model

| Goal                                 | Recommendation                  | Why? |
|--------------------------------------|----------------------------------|------|
| Add global user/role/menu/permission | Extend **ADMIN**                | Single source of truth for auth |
| Add new business domain (sales, hr, inventory, …) | Create **new model** | Separation of concerns, independent lifecycle |
| Add new ETLX pipelines / notebooks / dashboards for a domain | New model + reference shared tables | Clean & focused |
| Add very generic / cross-cutting feature | Consider extending ADMIN (rare) | Avoid ADMIN bloat |

**Rule of thumb**:  
If the tables belong to a **specific business area** → new model.  
If they affect **global system behavior** → ADMIN.

---

## Quick Checklist — New App Model

- [ ] Unique `name:` (e.g. SALES, INVENTORY, CRM)  
- [ ] Own SQLite/PostgreSQL/etc. file or schema  
- [ ] `admin_conn` pointing to ADMIN.db (for users/roles)  
- [ ] At least one `cs_app` group with your tables  
- [ ] `requires_rla: true` on tables that need per-row security  
- [ ] Use `form_display`, `table_display`, `form_size`, `form_code: json/markdown/sql` for nice UI  
- [ ] Optional: `table_extra_options` with custom components  
- [ ] Optional: seed data block at the end

---

The **ETLX model** is a perfect real-world example:

- it manages ETLX data pipelines  
- it includes notebook-style analysis  
- it re-uses shared dashboards  
- it lives in its own database (`ETLX.db`)  
- it integrates tightly with the central ADMIN model (users, roles, permissions)

This same pattern can be used to quickly build:

- inventory management  
- CRM  
- financial ledger  
- IoT telemetry backend  
- customer analytics warehouse  
- … or just a personal scratchpad database with nice UI

---

## ETLX — A Complete, Self-contained App Model

File: `etlx_model.md`

````markdown {linenos=table}
<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
# ETLX_MODEL
```yaml
name: ETLX
description: ETLX Model
runs_as: MODEL
conn: 'sqlite3:database/ETLX.db'
admin_conn: 'sqlite3:database/ADMIN.db'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    active: true
    menu_config: '{"label": "dashboard","tooltip": "dashboard_desc","load_items": {"table": "dashboard","tables": ["dashboard"]}}'
    tables:
      - dashboard
  ETLX:
    menu_icon: circle-stack
    menu_order: 2
    active: true
    tables:
      - etlx
      - etlx_conf
      - manage_query
  Notebook:
    menu_icon: book-open
    menu_order: 3
    active: true
    menu_config: '{"label": "notebook","tooltip": "notebook_desc","load_items": {"table": "notebook","tables": ["notebook"]}}'
    tables:
      - notebook
```

## ETLX
```yaml
table: etlx
comment: ETLX
columns:
  etlx_id:          { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etl:              { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 3 }
  etl_desc:         { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: false, form_size: 9 }
  attach_etlx_conf: { type: varchar(200), comment: "Config File" }
  etlx_conf:        { type: text, comment: "Config Text", form_display: true, form_long_text: true, form_code: markdown }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: etlx_id, order: DESC}]
table_extra_options:
  - - {size: 11, component: ETLX, label: etlx, icon: play, main: true, data: '{ "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh" }, { "type": "btn", "icon": "bolt", "name": "RUN_ALL", "class": "btn-sm text-info", "label": "crud.run_all" } ] }'}
```

## ETLX_CONF
```yaml
table: etlx_conf
comment: ETLX Extra Config
columns:
  etlx_conf_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etlx_conf:       { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 6 }
  etlx_conf_desc:  { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  etlx_extra_conf: { type: text, comment: "Config Text", form_display: true, form_long_text: true, form_code: json }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: etlx_conf_id, order: DESC}]
```

## MANAGE_QUERY
```yaml
table: manage_query
comment: Queries
columns:
  manage_query_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  manage_query:      { type: varchar(200), nullable: false, comment: "Query Desc", form_display: true, table_display: true, form_size: 6 }
  database:          { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 6 }
  manage_query_conf: { type: text, comment: "Query Config", form_display: true, form_long_text: true, form_code: json }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: manage_query_id, order: DESC}]
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true, form_size: 3 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 9 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 3 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - {size: 12, component: EvidenceDash, label: dashboard, intercept_r: true}
```

## DASHBOARD_COMMENT
```yaml
table: dashboard_comment
comment: Dashboards Comments
columns:
  dashboard_comment_id: { type: integer, pk: true, autoincrement: true, comment: "Comment ID" }
  dashboard_comment:    { type: text, comment: "Comments", form_display: true, table_display: true, form_long_text: true, form_code: markdown }
  dashboard:            { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true, form_size: 6 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: dashboard_comment_id, order: DESC}]
```

## NOTEBOOK
```yaml
table: notebook
comment: Notebooks
columns:
  notebook_id:   { type: integer, pk: true, autoincrement: true, comment: "Notebook ID" }
  notebook:      { type: varchar(200), comment: "Name", form_display: true, table_display: true, form_size: 6 }
  notebook_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  notebook_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:       { type: integer, comment: "User ID" }
  app_id:        { type: integer, comment: "App ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: notebook_id, order: DESC}]
table_extra_options:
- {size: 12, component: Notebook, label: notebook, icon: book-open, intercept_r: true}
```

# DATA
```yaml
name: DATA
description: DATA Model ETLX 
runs_as: MODEL_DATA
conn: 'sqlite3:database/ETLX.db'
```

## dashboard-1
```yaml
table: dashboard
description: Add default Logs Dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND dashboard = :dashboard AND excluded = false'
data:
  dashboard_id:   1
  dashboard:      Logs
  dashboard_desc: Logs Example
  dashboard_conf: FileContent(logs_dashboard.md)
  order:          1
  active:         true
  user_id:        1
  app_id:         2
  created_at:    Now()
  updated_at:    Now()
  excluded:      false
```
````

After placing the file in your models folder and running

```bash
central-set --init --model /path/etlx_model.md
```

it will create the database in any of supported database the supported and update admin model with all the necessary metadata needed for a new app.

You instantly get:

- A new **ETLX** menu group with full CRUD for pipelines  
- Custom **Run** and **Run All** buttons (via `table_extra_options` + `component: ETLX`)  
- Nice form layouts + code editors (markdown for pipeline config, JSON for extra config)  
- Reusable **dashboards** tab with EvidenceDash renderer  
- **Notebook** section for ad-hoc analysis  
- Automatic RBAC (who can see/edit/run pipelines)  
- Optional row-level security (add `requires_rla: true` to any table)  
- Possibility to schedule pipeline runs via the central `cron` table

---

## Why This Pattern Is So Powerful

| Feature you get "for free"          | How it appears in ETLX model                     |
|-------------------------------------|--------------------------------------------------|
| Own database / schema               | `conn: 'sqlite3:database/ETLX.db'` (or postgres://…, clickhouse://…, etc.) |
| Sidebar menu group                  | `cs_app: ETLX: …`                                |
| Beautiful CRUD UI + forms           | `form_display: true`, `form_size`, `form_code: markdown/json` |
| Custom buttons & components         | `table_extra_options` + `component: ETLX`        |
| Dashboards                          | Reference shared `dashboard` table               |
| Notebooks                           | Reference shared `notebook` table                |
| Data pipelines & queries            | `etlx`, `etlx_conf`, `manage_query` tables       |
| Role-based access                   | Controlled from ADMIN → roles → role_app_menu_table |
| Row-level security (optional)       | Add `requires_rla: true` to any table            |
| Scheduled execution                 | Link to central `cron` table                     |
| Easy prototyping → production       | Start with SQLite → later change conn to prod DB |


## Quick Variations — Use the Same Pattern for Anything

**Analytics / BI backend**  
Change menu to “Analytics”, tables to `kpi_daily`, `report`, `metric`

**Inventory app**  
`conn: 'postgresql:…'`  
tables: `product`, `warehouse`, `stock_movement`, `supplier`

**Personal finance tracker**  
`conn: 'sqlite3:finance.db'`  
tables: `transaction`, `category`, `account`, `budget`
---

You decide the database, the tables, the menu name — Central Set generates the rest.

1. Create one `.md` file per logical application/domain  
2. Use SQLite for fast local prototyping  
3. Switch to PostgreSQL/MSSQL/MYSQL/etc. when ready (just change `conn`)  
4. Reuse shared tables (`dashboard`, `notebook`, `cron`, `etlx`) across models  
5. Control access centrally from the ADMIN model

This gives you:

- instant admin UI  
- version-controlled schema & config  
- powerful data tools (pipelines, notebooks, dashboards)  
- strong security model  
- easy transition from prototype → production

The ETLX model is one of the best starting points — copy it, rename it, change the tables, and you have a brand-new app in minutes.

## Summary

Central Set shines when you **compose many small, focused models**.

Each new app/model automatically inherits powerful features:

- Full admin UI (forms + grids)  
- Security (RBAC + RLA)  
- Data pipelines (ETLX)  
- Scheduling (cron)  
- Interactive analysis (notebooks)  
- Dynamic dashboards  
- Extensibility (custom components, Arrow Flight, …)

This turns Central Set into a **meta-platform** — you build new applications just by writing declarative Markdown/YAML models.

Start small → one new model for your most important domain → iterate!

Questions / want a concrete model for HR / Inventory / Finance? Just describe the domain — we can generate a starter template.
