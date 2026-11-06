---
weight: 3020
title: "Admin Model"
description: "A guide to the declarative model for Central Set's main management database."
icon: schema
date: 2025-12-16T01:04:15+00:00
lastmod: 2026-11-06
draft: false
images: []
---

The **ADMIN model** defines the schema and management metadata for Central Set's main database. It describes how Central Set manages users, roles, applications, menus, permissions, validations, CRUD behavior, integrations, and UI configuration.

The source of truth is [`assets/models/admin_model.md`](../../../../assets/models/admin_model.md). It is a declarative Markdown file containing YAML blocks for model settings, application menus, table definitions, and configured records.

For an overview of ADMIN's responsibilities, see [Admin Schema (the DNA)](admin-schema/).

## What the model contains

The model has several related layers:

1. **Model settings** identify the model, its database connection, and metadata behavior.
2. **Application and menu configuration** lists the menus and tables Central Set exposes.
3. **Table definitions** describe columns, relationships, and database defaults.
4. **Metadata** describes permissions, UI behavior, validation, actions, query interception, and APIs.
5. **Configured records** in `data` blocks provide initial or example entries.

The file also contains additional model directives and examples. Check each section's `active` setting and operation flags: a definition in the file is not necessarily enabled.

## Model settings

The file begins with the ADMIN model header:

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

The connection uses `DB_DRIVER_NAME` and `DB_DSN` from the environment, keeping environment-specific connection details outside the model. Configure these values for your installation before applying the model; the [Quickstart](../quickstart/) describes the setup.

## Applications and menus

The `cs_app` block declares the application menus and the tables listed in each one. The ADMIN model includes groups such as Admin, APIs, Jobs Scheduling, Dashboards, Arrow Flight, and Quack. Menu definitions can include an icon, order, active flag, configuration, and table list.

For example, the model configures the Dashboards menu to include the `dashboard` table. The corresponding `app`, `menu`, and `menu_table` tables hold the application and navigation metadata used by Central Set.

## Table and column definitions

Each table block identifies a table and describes its columns. For example, the `lang` table includes a primary key, a unique language code, display hints, and an initial English-language record:

```yaml {linenos=table}
table: lang
comment: Languages
columns:
  lang_id:   { type: integer, pk: true, autoincrement: true }
  lang:      { type: varchar, len: 4, unique: true, nullable: false, form_display: true, table_display: true }
  lang_desc: { type: varchar, len: 200, form_display: true, table_display: true }
  excluded:  { type: boolean, default: false }
data:
  - { lang_id: 1, lang: en, lang_desc: English, excluded: false }
```

Common column properties include:

| Property | Purpose |
|---|---|
| `type`, `len` | Database type and, where applicable, length |
| `pk`, `autoincrement` | Primary-key behavior |
| `fk` | Foreign-key reference, such as `role.role_id` |
| `nullable`, `default`, `unique` | Database constraints and defaults |
| `comment` | Human-readable field description |

Table blocks can also include `form_layout`, `table_layout`, and `table_extra_options` to describe presentation and table-specific UI configuration.

## UI metadata

Column properties can describe how fields are shown and edited. The model uses hints such as `form_display`, `table_display`, `form_size`, `order`, `form_code`, `form_regex_val`, and `form_val_msg`. For example, `form_code: markdown` selects a Markdown-oriented editor for a text field; `form_code: sql` indicates SQL-oriented editing.

`form_layout` can organize fields into tabs, set form sizing, and configure related subforms. `table_layout` can configure ordering and table actions. `table_extra_options` can attach custom UI components to a table.

The model uses UI component intercept flags such as `intercept_r` in `table_extra_options`. These describe component/UI behavior; they are distinct from the database query rules declared in `crud_intercept`.

## Permissions and security metadata

ADMIN stores access rules at multiple levels:

- `role_app` associates roles with applications.
- `role_app_menu` controls role access to menus.
- `role_app_menu_table` stores create, read, update, delete, and share permissions for tables in menus.
- `column_level_access` stores column-level operation permissions.
- `row_level_access` and `role_row_level_access` describe user- and role-based permissions on individual rows.

Tables and menu-table entries can also be marked `requires_rla` when row-level access is required. The model includes roles and users, as well as examples of access configuration. Review the active state and access entries when changing security behavior.

## Validation rules

The `validation` table associates a validation with a database, table, and CRUD operation. A rule can specify an SQL check, an error message, severity, reaction, and active state. Related `validation_data` records can provide additional inputs, and `validation_logs` store outcomes.

The model includes an example rule that checks whether a user's email already exists when creating a user. This demonstrates how validation configuration is associated with a target table and operation.

## Actions triggered by CRUD

The `crud_action` table describes operations associated with create, read, update, or delete events. A rule identifies its target database and table, the operations that trigger it, its action type, and whether it is active or user-triggered.

The model provides fields and examples for SQL, email, internal or external API calls, PDF generation, and ETLX workflows. Actions may also include `action_data` inputs and `action_trigger_action` entries for chained behavior. The bundled examples illustrate patterns; some are intentionally marked inactive.

## Query interception

`crud_intercept` defines rules for processing a query before execution. Intercepts are scoped to a database, table, and CRUD operation, and can include a condition and SQL or rewrite configuration. The model declares intercept types for encapsulating a read query and rewriting a DML execution list. Related data and log tables hold additional configuration and records.

This is different from both a CRUD action—which is associated with an operation—and a UI component intercept flag such as `intercept_r`.

## API definitions and calls

The `api` table describes API integrations, including request method, endpoint, templates, retries, timeouts, and optional follow-up SQL. Related tables store headers, data configuration, and call logs.

An API can be invoked from a CRUD action, with templates using values from the operation context, or exposed as a UI action/button through table configuration such as `exec_button`. The model includes examples of both configured API records and actions that call APIs.

## Trace a CRUD action through the model

The model's `CRUD_ACTION_API_CALL` example is a compact illustration of how ADMIN metadata connects a database operation to an API:

```yaml {linenos=table}
table: crud_action
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: CRUD_ACTION_API_CALL
  action_type_id: 3
  table: users
  db: ADMIN
  active: true
  create: true
  update: true
  delete: true
  api: 'run/process_user/{{.action}}/{{.user_id}}'
```

Read this declaration as a rule attached to create, update, and delete operations on `ADMIN.users`. When enabled for an operation, the configured internal API path is rendered with values from that operation's context. The action definition also has an error message and `parallel` setting in the full model.

The `LOG_IP_GEODATA` example shows a longer chain:

1. A user-triggered CRUD action on `user_log` refers to the API named `IP_GEODATA`.
2. The API endpoint template uses the request IP from its input data.
3. An `action_trigger_action` entry connects successful processing to `HANDLE_IP_GEODATA`, which is configured as an ETLX workflow action.

These are model examples, not a recommendation to send user activity data to an external service. Before enabling or adapting such a chain, review what data it sends, which endpoint receives it, and how failures and logs should be handled.

## Schema definitions and data instructions

There are two related kinds of configuration in the file:

- A table definition can include a `data` list beside `columns`. For example, `lang` declares its table structure and the initial `en` language record together.
- Later `MODEL_DATA` entries use `table`, `cond`, and `data` to describe configuration records for existing tables. The bundled examples update app configuration, add a dashboard, and declare validation and CRUD-action records. A `children` block describes related records, such as a validation's `validation_data`.

The `cond` value identifies a row using a SQL `WHERE` condition and named values from `data`; it is not the validation SQL or an API request. This separation is useful when reading the file: table blocks describe the database and defaults, while data instructions configure records that use those tables.

## Other model modes in the file

After the core ADMIN model and its data examples, the file also contains blocks with other `runs_as` values, including `ROLE`, `ROLE_USERS`, and `MODEL_SQL`. These describe role access, user-to-role assignment, or SQL-oriented model work rather than new tables in the core ADMIN schema.

Treat each block as its own declaration and inspect its `active` setting before relying on it. For example, the included `ROLE_ACCESS` and `ROLE_USERS` headers are inactive. A SQL/procedure example later in the file is a separate concern from applying the ADMIN table definitions.

## Other ADMIN capabilities

The model also defines metadata for:

- **Access keys and environment values** — `access_key` and `env`.
- **Dashboards** — `dashboard` and `dashboard_comment`; dashboard configuration is stored in `dashboard_conf`.
- **Scheduled jobs** — `cron` and `cron_log`.
- **Arrow Flight and Quack** — catalogs, schemas, table scopes, server/API configuration, and logs.
- **Translations and table metadata** — translation tables, table registry, and schema metadata.

The tables are part of the central management model; application-specific business data belongs in that application's database and model.

## How the management layers fit together

The model describes linked records rather than isolated features:

```text {linenos=table}
users ── user_role ── role
                         └── role_app ── app
                              role_app_menu ── menu
                                   role_app_menu_table ── table
                                                         ├── CRUD permissions
                                                         ├── column_level_access
                                                         └── row_level_access

table + operation ── validation rules
                  ├── query intercept rules
                  └── CRUD actions ── action data / API / chained actions
```

The access tables answer different questions: which applications can a role enter, which menus can it see, which CRUD operations can it perform on a table, and whether a particular column or row is accessible. `requires_rla` on a table or menu-table mapping indicates that row-level rules are relevant for that table; it is not itself a row permission.

For data operations, validations, query interception, and CRUD actions are separate configuration systems. A validation declares a check and response for an operation; an intercept can transform or wrap a query before execution; an action performs configured work associated with an operation. Their matching table/database and operation fields determine their scope.

## Review before deployment

The model is executable configuration, not just reference documentation. Before applying changes to a shared or production database:

- Review table changes and any model data instructions, including their `cond` selectors and `children`.
- Review users, roles, seeded access, and all application/menu/table/column/row permissions.
- Change or remove bootstrap/demo accounts and confirm password-reset expectations.
- Keep access tokens, environment values, and external API credentials controlled; do not put real secrets in a model committed to source control.
- Check `active`, CRUD operation flags, `requires_rla`, and user-trigger settings rather than assuming every declared feature is enabled.
- Check SQL conditions/templates, external endpoints, data sent to those endpoints, and follow-up SQL.
- Back up an existing database before applying a model change and verify the resulting access and behavior.

## Apply the model

From a repository checkout, initialize the configured database with:

```bash
central-set --init --model assets/models/admin_model.md
```

To apply the model to an existing database, use:

```bash
central-set --update --model assets/models/admin_model.md
```

The packaged application can also load its embedded model by file name:

```bash
central-set --init --model admin_model.md
```

Review the model and database configuration before applying changes, especially when working with an existing database. To add a separate business application, see [Creating a New App / Data Model](new-app-model/).

## Key points

- ADMIN is Central Set's main management database; application databases hold their own domain data.
- The model declares the ADMIN schema, application/menu metadata, security rules, and platform behavior in one version-controlled file.
- Table definitions can include both database properties and UI metadata.
- A declared rule or example is not necessarily active; check its activation and operation flags.
