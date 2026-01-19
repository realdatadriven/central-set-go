---
weight: 7071
title: "Overview"
description: "Configuration-driven data pipelines powered by ETLX"
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Overview

Central Set was born primarily from the need to **manage data engineering workflows**, but real-world data systems require much more than just pipelines.

To operate reliable data platforms, you need:

* User and role management
* Database and schema management
* Forms and datatables to manage configurations
* Secure access keys and API tokens
* UI-driven configuration
* Observability over pipelines
* A way to **run data workflows reliably**
* Dashboards for monitoring and insight

What started as a data-engineering support layer evolved into something broader.

Because Central Set is **fully configuration-driven**, database-backed, and UI-agnostic, it can act as a **general-purpose backend platform** for small to medium-scale applications — with a UI that adapts to your data model.

That data-engineering execution layer is **ETLX**.

> **ETLX is the data pipeline engine behind Central Set.**

Central Set provides the **UI, configuration, security, and orchestration layer**.
ETLX provides the **runtime, execution model, and pipeline semantics**.

---

## What is ETLX?

**ETLX** is a lightweight, specification-driven data pipeline engine.

It focuses on:

* Declarative pipeline definitions
* Explicit data movement
* Deterministic execution
* Simple primitives instead of heavy abstractions

ETLX pipelines are defined as **specifications**, not frameworks or monolithic codebases.

📘 Full ETLX documentation
👉 [https://realdatadriven.github.io/etlxdocs](https://realdatadriven.github.io/etlxdocs)

📦 ETLX repository
👉 [https://github.com/realdatadriven/etlx](https://github.com/realdatadriven/etlx)

---

## Central Set + ETLX

Central Set ships with ETLX **enabled by default**.

Together, they form a layered system:

```
┌────────────────────────────┐
│ Central Set UI             │
│ - Users & Roles            │
│ - Databases                │
│ - Apps & Menus             │
│ - Tables & Forms           │
│ - Access Keys              │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ Configuration & Metadata   │
│ (stored in databases)      │
└──────────────┬─────────────┘
               │
               ▼
┌────────────────────────────┐
│ ETLX Runtime               │
│ - Pipelines                │
│ - Extract / Transform      │
│ - Load / Write             │
│ - File & DB connectors     │
└────────────────────────────┘
```

The **UI never hardcodes pipelines**.

Everything ETLX runs is:

* Defined via database-backed specifications
* Triggered through APIs
* Observable and manageable from the platform

---

## Why ETLX Exists

Most data tooling ecosystems assume:

* Heavy orchestration frameworks
* YAML sprawl
* Code-first pipelines
* Tight coupling to infrastructure

ETLX takes a different approach:

* Pipelines are **data-first**
* Configuration lives in **databases**
* Execution is **explicit and deterministic**
* The UI is **just another client**

This makes ETLX well suited for:

* Embedded data platforms
* Multi-tenant systems
* Admin-driven data workflows
* Headless or UI-driven execution

---

## ETLX as a First-Class Application

In Central Set, ETLX is treated like any other **application**:

* It can define:

  * Menus
  * Tables
  * Permissions
* Pipelines can be:

  * Configured via UI
  * Triggered via API
  * Observed through dashboards

This allows you to build:

* Data ingestion systems
* Transformation pipelines
* Validation workflows
* Automation jobs
* Internal data products

All without writing frontend code.

---

## Running Pipelines via API

ETLX pipelines can be executed **headlessly** using the API.

### Run Pipeline

**Endpoint**

```
GET /etlx/run/{name}
```

**Description**

Executes an ETLX pipeline by name.

* `{name}` is the pipeline identifier
* Execution context is resolved at runtime
* Permissions are enforced via access keys or user tokens

**Example**

```
GET /etlx/run/daily_sales_load
```

This enables:

* CI/CD triggers
* Cron-based execution
* External system integration
* Event-driven workflows

---

## Observability & Health (Conceptual)

ETLX pipelines are designed to expose:

* Execution status
* Errors
* Logs
* Metadata about runs

Central Set can leverage this data to build:

* Pipeline health dashboards
* Execution history views
* Error inspection tools
* Retry and recovery flows

> ETLX does not impose a visualization model.
> Central Set provides the building blocks to create one.

---

## Key Principles

* **Data apps need more than pipelines**
* **Pipelines need configuration, not frameworks**
* **The database is the control plane**
* **The UI reflects metadata**
* **ETLX executes what Central Set describes**
