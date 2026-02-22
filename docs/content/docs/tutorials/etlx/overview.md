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


## Secure Execution Configuration

By default, ETLX is flexible and allows execution of configurations coming from different contexts (UI, API, CLI-style payloads).

However, in hardened deployments — especially when Central Set is exposed:

- Outside a private LAN
- To untrusted networks
- Through public APIs
- Or in multi-tenant environments

You may want to **restrict how pipeline configurations are executed**.

To support this, Central Set provides additional security flags.

---

## Execution Security Modes

### 1️⃣ `ETLX_ALLOW_CLI_CONFIG`

```env
ETLX_ALLOW_CLI_CONFIG=true
````

**Default behavior when NOT enabled:**

Only this endpoint is allowed:

```
GET /etlx/run/{name}
```

In this mode:

* The pipeline configuration is located **exclusively in the database**
* The server loads the stored specification
* No client-submitted configuration is executed
* The client cannot alter execution semantics

This is the safest and most restrictive mode.

---

### 2️⃣ `ETLX_VALIDATE_ETLX_ACCESS`

```env
ETLX_VALIDATE_ETLX_ACCESS=true
```

When enabled:

* The user must have explicit access to the ETLX record stored in the database
* The server checks:

  * Application ownership
  * Role permissions
  * ETLX configuration visibility
* Execution is denied if access is not granted

This prevents:

* Users from triggering pipelines they should not see
* Cross-app execution
* Unauthorized automation

Recommended for:

* Multi-tenant environments
* Shared infrastructure
* Regulated systems

---

### 3️⃣ `ETLX_VALIDATE_CLI_CONF_WITH_DB`

```env
ETLX_VALIDATE_CLI_CONF_WITH_DB=true
```

Used **in combination with**:

```env
ETLX_ALLOW_CLI_CONFIG=true
ETLX_VALIDATE_ETLX_ACCESS=true
```

This option ensures that:

* If a configuration is submitted from the client
* It must match the server-side database configuration
* Queries and structural definitions are validated

This prevents:

* Injecting modified SQL
* Altering execution order
* Changing connections dynamically
* Bypassing governance rules

In other words:

> The client may request execution,
> but cannot redefine execution.

---

## Recommended Production Configuration

For secure environments exposed beyond a trusted network:

```env
ETLX_ALLOW_CLI_CONFIG=true
ETLX_VALIDATE_ETLX_ACCESS=true
ETLX_VALIDATE_CLI_CONF_WITH_DB=true
```

This enforces:

* Server-side configuration authority
* Strict permission validation
* No execution of arbitrary client-submitted specs

---

## Security Philosophy

ETLX is intentionally flexible because:

* Central Set can operate as a controlled internal platform
* Some deployments require dynamic execution
* Some environments are trusted

But security-sensitive environments require stricter guarantees.

These flags allow you to choose between:

| Mode            | Behavior                                 |
| --------------- | ---------------------------------------- |
| Open / Flexible | Client-submitted configuration allowed   |
| Controlled      | Server validates configuration ownership |
| Hardened        | Only DB-defined pipelines executable     |

---

## Why This Matters

Without validation:

* A malicious or compromised client could modify SQL
* Change output paths
* Execute unauthorized connections
* Trigger unintended workloads

With validation enabled:

* The database becomes the single source of truth
* The server becomes the execution authority
* Clients become requesters, not spec authors

This aligns with Central Set's principle:

> The database is the control plane.
> The runtime executes what the platform defines.

---

## Summary

These flags allow ETLX to operate in three security tiers:

1. Development / Internal
2. Controlled Production
3. Hardened / Public-facing

You can choose the appropriate level based on:

* Network exposure
* Tenant isolation needs
* Regulatory requirements
* Governance strictness

ETLX remains flexible —
but security is always configurable.

## Exposing ETLX to External Clients (Security Recommendation)

Allowing external clients (outside your LAN or trusted network) to execute pipelines introduces additional risk.

Even with:

* `ETLX_VALIDATE_ETLX_ACCESS=true`
* `ETLX_VALIDATE_CLI_CONF_WITH_DB=true`
* Strict token validation
* Role-based access control

There is still increased attack surface:

* Increased brute-force attempts
* Token leakage risk
* Misconfiguration exposure
* Workload abuse (DoS-style heavy pipeline execution)
* Query abuse if misvalidated

For this reason, **the recommended architecture is isolation**.

---

## Recommended Architecture: Public Execution Instance

If external execution is required:

### 🟢 Deploy a Dedicated Public Instance

This instance should:

* Only connect to **public or non-sensitive datasets**
* Not have access to:

  * Internal databases
  * Admin schemas
  * ETLX internal metadata beyond what is necessary
* Use a restricted database user
* Run with limited network access

Think of it as a **data-serving node**, not your control plane.

---

## Separation of Concerns

| Instance Type     | Purpose                                    | Data Scope       |
| ----------------- | ------------------------------------------ | ---------------- |
| Internal Instance | Admin UI, configuration, private pipelines | Full access      |
| Public Instance   | External API / public pipelines            | Public-only data |

This ensures that even if:

* A token is leaked
* An endpoint is abused
* A configuration bypass is discovered

The blast radius is limited to **public data only**.

---

## Defense-in-Depth Strategy

For hardened deployments:

1. Disable CLI config if not required:

```env
ETLX_ALLOW_CLI_CONFIG=false
```

2. Require access validation:

```env
ETLX_VALIDATE_ETLX_ACCESS=true
```

3. Validate client config against DB:

```env
ETLX_VALIDATE_CLI_CONF_WITH_DB=true
```

4. Deploy public-facing workloads on isolated infrastructure.

5. Optionally:

   * Add rate limiting
   * Add request quotas
   * Use API gateways
   * Monitor execution frequency
   * Restrict maximum pipeline runtime

---

## Why Isolation Is Stronger Than Validation

Validation protects logic.
Isolation protects infrastructure.

Even perfectly validated logic can still:

* Execute expensive queries
* Consume compute
* Trigger heavy I/O
* Expose performance-sensitive datasets

Isolation ensures:

> A compromise affects only the intended exposure layer — never your core platform.

---

## Production Principle

If a pipeline can be executed by:

* Anonymous users
* External partners
* Public systems
* Third-party automation

Then that pipeline should execute in an environment that assumes:

> It may eventually be abused.

Design for containment, not perfection.

---

## Final Recommendation

If you must allow outsiders to execute pipelines:

✔ Use a dedicated instance
✔ Restrict it to public data only
✔ Apply strict ETLX validation flags
✔ Monitor usage
✔ Treat it as an exposed surface

Central Set + ETLX are flexible by design —
but secure architecture is your responsibility.

Isolation is your strongest guarantee.

