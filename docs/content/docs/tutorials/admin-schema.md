---
weight: 7030
title: "Admin Schema (the DNA)"
description: "How Central Set defines applications, UI, permissions, and automation using database metadata."
icon: cogs
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Overview

{{% alert context="warning" text="**Caution** — Central Set is **production-ready and actively used**, but this documentation is still **under active development**. Large parts of the docs are **auto-generated**, evolving alongside the platform, and some sections may be incomplete, rough around the edges, or change frequently."  /%}}


The **ADMIN schema** is the **DNA of Central Set**.

Everything you see in the UI — applications, menus, tables, forms, permissions, and even automation — is defined as **data stored in the ADMIN database**.

There is no hardcoded admin UI.
There are no per-app APIs.
There are no generated frontends.

Instead, Central Set uses a **metadata-driven control plane** that describes *what exists*, *who can access it*, and *how it behaves*.

> **If it exists in the ADMIN schema, it can exist in the UI and APIs.**

---

## The ADMIN Database as a Control Plane

The ADMIN database is not just a configuration store.

It acts as a **control plane** that defines:

* Applications and their databases
* Navigation structure (menus)
* Tables and schemas
* CRUD behavior
* Permissions and access rules
* UI layout and forms
* Automation hooks
* Integration with ETLX pipelines

Central Set reads this metadata at runtime and uses it to **generate UI behavior and API responses dynamically**.

---

## A Metadata-First Architecture

Central Set follows a simple but powerful rule:

> **The database is the specification.**

Instead of encoding logic in code or YAML files, Central Set stores intent as structured metadata.

This enables:

* Runtime UI generation
* Dynamic APIs
* Multi-tenant setups
* Safe extensibility
* Zero frontend rebuilds

---

## Conceptual Model (No Table Names)

At a high level, the ADMIN schema models the system in layers:

```

Users & Roles
↓
Applications
↓
Menus
↓
Tables
↓
Permissions & UI Behavior
↓
Automation & Pipelines

```

Each layer builds on the previous one, and each layer is defined as data.

---

## Applications, Menus, and Tables

### Applications

An **Application** represents a logical admin interface.

Each application:

* Is linked to one database
* Defines a context for menus and tables
* Can be isolated by permissions
* Can represent a product, domain, or system

The built-in **ADMIN** app is just one application — not a special case.

---

### Menus

Menus define **navigation**, not logic.

They group tables and features into meaningful sections and control how users move through the application.

Menus are:

* Ordered
* Permission-aware
* Fully configurable via metadata

---

### Tables

Tables are the heart of the system.

A table definition tells Central Set:

* Which database table to expose
* How CRUD operations behave
* Which columns are visible
* How data is labeled
* How forms are rendered
* What permissions apply

> Adding a table to the ADMIN schema is enough to make it appear in the UI.

---

## Schema Discovery & Introspection

Central Set dynamically inspects connected databases to understand:

* Tables
* Columns
* Types
* Constraints
* Relationships

This information is cached and enriched with metadata stored in the ADMIN schema.

This allows Central Set to:

* Render forms automatically
* Validate input
* Build datatables
* Enforce permissions
* Stay in sync with the database

---

## Permissions Are Data

Permissions in Central Set are **not hardcoded rules**.

They are metadata-driven and can apply at multiple levels:

* Application access
* Menu visibility
* Table access
* CRUD operations
* Column visibility
* Row-level filtering

Because permissions live in the database:

* They can be managed via the UI
* They apply consistently to UI and APIs
* They scale naturally in multi-tenant systems

---

## UI Customization via Metadata

The ADMIN schema also defines **how things look and behave**.

This includes:

* Column labels
* Field order
* Form layouts
* Widgets and input types
* Translations
* Read-only or computed fields

This makes Central Set suitable not just for admins, but for **data-facing internal tools and apps**.

---

## Automation & ETLX Integration

The ADMIN schema also acts as the **bridge to automation**.

It defines:

* Which workflows exist
* How they are triggered
* What data they operate on
* How they integrate with ETLX pipelines

Central Set handles:

* Configuration
* Security
* Triggers
* Observability hooks

ETLX handles:

* Execution
* Data movement
* Pipeline semantics

Together, they form a complete data platform.

👉 See **ETLX & Data Pipelines** for more details.

---

## The Full Schema (Source of Truth)

This documentation explains **concepts**, not every table and column.

The **canonical definition** of the ADMIN schema lives in the repository as a living document:

📓 **ADMIN schema notebook**  
👉 https://github.com/realdatadriven/central-set-go/blob/main/db.ipynb

This notebook documents:

* Tables
* Relationships
* Intended usage
* Evolution of the schema

If you want to extend Central Set at a deep level, this is the place to start.

---

## Key Takeaways

* ADMIN is fully metadata-driven
* The database is the control plane
* UI and APIs share the same contract
* Permissions are data, not logic
* ETLX executes what ADMIN describes

> **Central Set does not define applications in code.  
> It defines them in data.**

## Adding a New Application (How ADMIN Grows)

Central Set ships with a **fully configured ADMIN application**, defined entirely in the ADMIN database.

Adding a new application — such as **ETLX** — does **not** require frontend work or custom backend code.
Instead, you **extend the metadata** that already powers ADMIN.

The ETLX application included in the repository is the **canonical example** of how this works.

---

## The High-Level Flow

Adding a new app follows a **deterministic, repeatable process**:

```
1. Create a database schema for the app
2. Register the app in ADMIN
3. Discover tables from the app database
4. Attach tables to menus
5. Customize UI behavior (optional)
```

This entire flow is **data-driven** and can be replayed on any machine.

---

## 1. Define the Application Database

Each application points to **one database**.

For ETLX, this database is defined in:

📓 `app.ETLX.ipynb` (repository root)

This notebook:

* Defines ETLX domain tables (pipelines, dashboards, notebooks, etc.)
* Uses SQLAlchemy only as a **schema authoring tool**
* Creates the physical database if it does not exist
* Adds rich **table and column comments**

> These comments later become **UI labels and tooltips**.

At this stage, Central Set does **not yet know** about the app — only the database exists.

---

## 2. Register the Application in ADMIN

Once the ETLX database exists, it is registered as an **Application** inside the ADMIN schema.

This is done by inserting a record into the `app` table:

* App name (e.g. `ETLX`)
* Description
* Version
* Database name
* Ownership metadata

From this moment on:

> **Central Set recognizes the database as an application context.**

---

## 3. Automatic Schema Discovery

After the app is registered, Central Set **introspects the database**:

* Tables
* Columns
* Types
* Primary keys
* Foreign keys
* Defaults
* Comments

This discovery process populates ADMIN metadata tables such as:

* Table registry
* Column registry
* Translations
* Relationships

No manual mapping is required.

> If a table exists in the app database, it becomes *eligible* for the UI.

---

## 4. Create Menus for the App

Applications do not expose tables directly.

They expose **menus**, which define navigation and grouping.

For ETLX, menus include:

* Dashboards
* ETLX
* Notebook

Each menu is:

* Created in the ADMIN database
* Linked to the application
* Assigned an icon and order
* Marked active or inactive

Menus are **pure metadata** — they do not contain logic.

---

## 5. Attach Tables to Menus

Once tables are discovered, they can be attached to menus.

This step defines:

* Which tables appear in which menu
* Whether the table is active
* Special behaviors (e.g. row-level access)

For example, in ETLX:

* `etlx`, `etlx_conf`, `manage_query` → **ETLX menu**
* `dashboard` → **Dashboards menu**
* `notebook` → **Notebook menu**

At this point:

> **The UI becomes fully functional for the new app.**

Tables appear as datatables.
Forms are auto-generated.
CRUD APIs are enabled.

---

## 6. UI Behavior Comes from Metadata

The ETLX notebook also demonstrates how Central Set uses metadata to control UI behavior:

* Table comments → table titles
* Column comments → form labels
* JSON configs → menu behavior
* Flags → visibility and activation

Because this configuration is stored in ADMIN tables, it can be:

* Modified from the UI itself
* Exported
* Reapplied on another installation

---

## Portability & Reproducibility

A key design goal is **repeatable installations**.

The ETLX setup shows how Central Set supports:

* Rebuilding apps on new machines
* Syncing UI customizations
* Shipping opinionated defaults
* Treating the database as the source of truth

This allows Central Set to behave like a **data-driven platform**, not a static admin panel.

---

## What This Enables

By following this pattern, you can:

* Add new domains without frontend work
* Build internal tools quickly
* Expose databases safely
* Treat data models as products
* Give ETLX (or any engine) a full UI

> **ADMIN defines the language.
> Applications define the vocabulary.
> Databases define the meaning.**

