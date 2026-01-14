---
weight: 100
date: "2026-01-03T10:00:00+00:00"
draft: false
title: "Overview"
icon: "circle"
toc: true
description: "Central Set (CS) is a configuration-driven, database-reflective platform with a built-in admin UI, secure APIs, data pipelines, and analytics."
publishdate: "2026-01-03T10:00:00+00:00"
tags: ["Beginners", "Admin", "Data Platform", "Databases", "ETLX"]
categories: ["Concepts"]
twitter:
    card: "summary"
    title: "What is Central Set?"
    description: "A database-reflective data platform with built-in admin UI, APIs, and pipelines"
    image: ""
---

---

## Welcome to the **Central Set (CS) documentation**

This documentation explains **what CS is**, **how it works**, and **how to think about it** before diving into applications, databases, security, APIs, pipelines, and dashboards.

If you want to get running quickly, start here:
👉 [Quickstart]({{% relref "quickstart" %}})

---

## What is Central Set?

**Central Set (CS)** is an **open-source, configuration-driven data platform** built with **Golang**.

At its core, CS provides a **built-in admin application** that turns your existing database into:

* A fully functional **Admin UI UI**
* A consistent, **secure API layer**
* A foundation for **data pipelines and analytics**

CS does **not** generate schemas and does **not** hide or abstract your data model.

Instead, it treats the **database itself as the source of truth**.

Your tables *are* the data model — and CS reflects them faithfully into UI, APIs, pipelines, and analytics, with configuration defining behavior rather than structure.

---

## Database as the Data Model

In CS, the **database schema is the canonical model**.

When a database is connected:

* Tables are discovered directly from the database
* Each table becomes manageable through the admin UI
* Data can be browsed, created, updated, and deleted
* APIs and UI share the exact same data contract

There is no duplication of models, DTOs, or schemas.

If it exists in the database, CS can manage it.

---

## Built-In Admin UI (Auto-Generated, Yet Fully Customizable)

CS ships with a **built-in admin application** — no UI builders, no code generation, no separate frontend project.

### Automatic UI generation

For every discovered table, CS automatically provides:

* Data tables for browsing records
* CRUD operations (create, read, update, delete)
* Auto-generated forms for create and edit

UI generation is **schema-driven**:

* Column comments are used as human-friendly labels
* Column names remain stable field identifiers behind the scenes
* Data types drive widgets, formatting, and validation
* Foreign keys define relations, selectors, and lookups

This guarantees that:

* The UI always matches the real schema
* APIs and UI never drift apart
* Schema refactors remain predictable and safe

---

## Admin Configuration Model (Apps → Menus → Tables)

CS uses its **own admin database** to define how your data is presented and managed.

Core admin concepts include:

* **Apps** – logical applications or domains
* **Menus** – navigation and grouping
* **Tables** – database tables exposed in the UI

These relationships:

```
Apps → Menus → Tables
```

Define how CRUD interfaces are generated for your database.

Through the admin UI you can configure:

* Which tables are exposed
* Menu structure and navigation
* Table views, filters, and ordering
* Form layouts (tabs, steps, sections, subforms)
* Relationship behavior (inline edits, selectors, popups)
* Permissions and access rules

All configuration is:

* Stored in the admin database
* Managed through the same CS UI
* Applied instantly, without redeploys

---

## Multi-Database Support (via `.env`)

CS supports **multiple database engines out of the box**.

The application database is configured using a simple `.env` file:

```env
DB_DRIVER=postgres
DB_DSN=postgres://user:password@host:5432/appdb
```

### Key points

* The driver and DSN fully define the application database
* Any database supported by **`sqlx`** can be used
* Switching databases requires **no code changes**
* Different apps can point to different databases

### Commonly used engines

* PostgreSQL
* SQLite
* MySQL / MariaDB
* SQL Server
* Any SQL database supported by `sqlx` drivers

This makes CS suitable for:

* Local development
* Embedded deployments
* Internal tools
* SaaS and multi-tenant platforms

---

## Secure API Access (Keys & Tokens)

Everything in CS is **API-first**.

* The admin UI is just another API consumer
* Access keys and tokens can be generated for users and services
* Permissions are enforced via role-based access control (RBAC)
* External systems interact safely with the same APIs

This enables:

* Automation and scripting
* CI/CD integrations
* Service-to-service communication
* Headless usage of CS as a backend

---

## Data Pipelines with ETLX

CS embeds **ETLX** as its data engineering engine.

With ETLX you can:

* Define metadata-driven pipelines
* Visualize pipelines as workflows
* Execute them manually, via API, or on a schedule
* Capture execution logs, metadata, and results

Pipelines are managed through the same admin UI and governed by the same security and permission model.

---

## Analytics, Notebooks & Dashboards

CS also includes tools for **data exploration and analytics**:

* Built-in SQL notebooks for exploration
* Configuration-driven dashboards (Markdown-based)
* DuckDB (including WASM) for analytical workloads
* Arrow Flight for sharing datasets with external tools

This connects **data engineering outputs** directly to **analysis and consumption**.

---

## How CS Is Different

CS is **not**:

* A CRUD code generator
* A schema designer
* A low-code UI builder
* A collection of loosely coupled tools

CS **is**:

* A **database-reflective admin platform**
* A **configuration-driven control plane**
* An **API-first backend with secure tokens**
* A **data engineering platform powered by ETLX**
* An **analytics foundation built on DuckDB**

> **The database schema is the source of truth.**
> **Configuration defines behavior — not structure.**

---

## Who Is CS For?

CS is designed for:

* Teams building **internal tools and Admin UIs**
* SaaS platforms needing **secure, multi-database management**
* Data engineers wanting **pipelines, governance, and UI in one place**
* Platform teams tired of stitching admin tools, APIs, and ETL together

---

## Where to Go Next

* 👉 [Quickstart]({{% relref "quickstart" %}}) — run CS locally
* 👉 [Core Concepts]({{% relref "quickstart" %}})
* 👉 [Admin & UI]({{% relref "quickstart" %}})
* 👉 [Security & API Access]({{% relref "quickstart" %}})
* 👉 [ETLX & Pipelines]({{% relref "tutorials/etlx" %}})
* 👉 [Dashboards & Analytics]({{% relref "tutorials/dashboards" %}})

---

CS is open source and evolving.

If you believe the **database should be the source of truth** — and that admin, APIs, data pipelines, and analytics should all align around it — you’re in the right place.
