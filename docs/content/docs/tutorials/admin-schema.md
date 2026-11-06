---
weight: 7030
title: "Admin Schema (the DNA)"
description: "How the ADMIN database manages Central Set users, applications, access control, validations, CRUD behavior, and APIs."
icon: schema
date: 2025-12-16T01:04:15+00:00
lastmod: 2026-11-06
draft: false
images: []
---

## Overview

The **ADMIN database is Central Set's main management database**. It holds the definitions and records Central Set uses to manage users, applications, menus, permissions, integrations, and data operations across the platform.

The declarative source for its schema and initial configuration is [`assets/models/admin_model.md`](../../../../assets/models/admin_model.md). The model defines database tables, their relationships and UI metadata, application menus, and initial records in one place.

Application databases hold their own domain data. ADMIN registers those applications and manages the shared users, roles, access rules, and behavior associated with them.

## What ADMIN manages

| Area | What the model provides |
|---|---|
| **Users and roles** | User accounts and profiles, role assignments, login-related fields, and user activity logs. |
| **Applications and menus** | Application registry, menu structure, menu-to-table assignments, ordering, icons, and configuration. |
| **Credentials and environment** | Access-key/token records and environment-variable records, including activation and expiry/configuration fields. |
| **Permissions** | Role access to applications and menus, CRUD permissions for tables, column-level permissions, and user- or role-based row-level access. |
| **Validation** | Rules associated with a database, table, and CRUD operation, with SQL conditions, messages, severity/reaction settings, and validation logs. |
| **CRUD actions** | Configurable operations associated with create, read, update, and delete events, including SQL, email, API calls, PDF generation, and ETLX workflows. |
| **CRUD interception** | Rules for handling or rewriting queries before execution, with conditions and optional follow-up SQL. |
| **APIs and integrations** | API definitions, request methods, endpoints, headers, templates, retries, timeouts, data mappings, and call logs. |
| **Other platform features** | Dashboards, scheduled jobs, Arrow Flight and Quack configuration, translations, and supporting metadata. |

These are not just labels in a model: the database includes corresponding tables such as `users`, `role`, `app`, `menu`, `access_key`, `role_app_menu_table`, `row_level_access`, `validation`, `crud_action`, `crud_intercept`, and `api`.

## Applications and navigation

The model's `cs_app` section declares the applications' menu groups, and the `app`, `menu`, and `menu_table` tables store the application and navigation metadata. A menu groups tables and can define its icon, order, active state, and configuration.

For example, the ADMIN model declares menu groups for Admin, APIs, Jobs Scheduling, and other platform areas. Each menu lists the tables it exposes. This metadata drives the platform's navigation and determines which tables are available through that application context.

Table and column definitions also carry UI metadata: field labels, form and table visibility, ordering, layouts, and specialized components. This lets Central Set build data interfaces from the model and metadata instead of defining a separate interface for every table.

## Layered access control

ADMIN represents permissions at several levels:

1. **Application access** — `role_app` associates roles with applications.
2. **Menu access** — `role_app_menu` controls access to an application's menus.
3. **Table CRUD access** — `role_app_menu_table` grants create, read, update, delete, and share permissions for a table in a menu.
4. **Column-level access** — `column_level_access` records operation permissions for a column.
5. **Row-level access (RLA)** — `row_level_access` and `role_row_level_access` describe per-row permissions for a user or role. Tables and menu-table entries can also be marked `requires_rla`.

This makes access decisions part of the central application metadata, rather than a collection of unrelated per-screen settings.

## Validation rules

The `validation` table associates a rule with a database, table, and operation such as create, read, update, or delete. A rule can specify an SQL check, an error message, whether it is active, and its severity and reaction. `validation_data` can provide additional rule inputs, while `validation_logs` records validation outcomes.

The model includes a seeded example that checks whether a user's email already exists on user creation. It illustrates how a rule can be declared with its target table and operation instead of being embedded in an individual form.

## Actions on CRUD operations

The `crud_action` table describes work associated with CRUD events. A rule identifies its table and database, selects which operations trigger it, and can be activated or exposed as a user-triggered action. Depending on its action type, the model provides fields for SQL, email templates, internal or external API calls, PDF generation, and ETLX workflows. Actions can also have related `action_data` and chained `action_trigger_action` records.

For example, the model includes illustrative rules for sending an email or writing a user log after a user record is created, and an API-call action that uses values from the operation context. Examples marked inactive are configuration examples, not necessarily enabled behavior.

## Query interception

CRUD actions run in response to configured operations; `crud_intercept` describes a separate concern: processing or rewriting a query before it executes. An intercept is scoped by database, table, and CRUD operation, and can include a condition and SQL template. The declared intercept types include encapsulating a read query when its condition matches and rewriting a DML execution list. `after_sql` can specify follow-up SQL, and interception logs are represented separately.

## APIs that actions and UI controls can call

The `api` table stores API definitions, including the endpoint, request method, headers and body templates, retry/timeout settings, and optional follow-up SQL. Related tables hold headers, data mappings, and API-call logs.

An API can be used by a CRUD action or exposed as an action/button on a table or menu. The model's examples show an action configured to call an internal API using operation values such as the action name and user ID; another example defines an external API and passes request data into its endpoint template. The API table's `exec_button` configuration is one example of a UI control that invokes an API. Templates can use the relevant operation or record data, so a call can act on the data that triggered it.

## The declarative model

The model file brings together:

- The database and model settings
- Application menus and table assignments
- Table columns, relationships, and UI configuration
- Security, validation, action, interception, and API definitions
- Initial records declared in `data` blocks

To initialize the configured database from a repository checkout:

```bash
central-set --init --model assets/models/admin_model.md
```

To apply the model to an existing database:

```bash
central-set --update --model assets/models/admin_model.md
```

The model is the declarative source for the ADMIN database structure and its configured metadata. Active flags and operation-specific settings determine which declared features are enabled.

## Extending ADMIN

Update [`assets/models/admin_model.md`](../../../../assets/models/admin_model.md) when changing Central Set's shared management model—for example, adding a shared permission or integration capability. Keep domain-specific business tables in their own application model and database, then register the application in ADMIN. See [Creating a New App / Data Model](new-app-model/).

## Key takeaways

- ADMIN is the main database for managing Central Set and its applications.
- It manages users, roles, applications, menus, credentials, and layered access control.
- It also holds declarative validation, CRUD action, query interception, and API configuration.
- Application databases continue to hold their own domain data.
- [`admin_model.md`](../../../../assets/models/admin_model.md) is the declarative source for ADMIN's schema and initial configuration.
