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
