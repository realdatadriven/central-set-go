---
weight: 100
date: "2026-01-03T10:00:00+00:00"
draft: false
title: "Overview"
icon: "circle"
toc: true
description: "ETLX is a modern, composable, metadata-driven ETL framework for data engineers."
publishdate: "2026-01-03T10:00:00+00:00"
tags: ["Beginners", "ETL", "DuckDB"]
categories: ["Concepts"]

twitter:
  card: "summary"
  title: "What is ETLX?"
  description: "A modern, composable ETL framework built for data engineers"
  image: ""
---

Welcome to the **ETLX documentation**.

This guide explains **what ETLX is**, **why it exists**, and **how to think about it** before diving into pipelines, configuration, and execution details.

If you are looking for a fast way to get started, jump directly to the [Quickstart]({{% relref "quickstart" %}}).

---

## What is ETLX?

**ETLX** is an **open-source, developer-first ETL framework** designed to make data pipelines:

- **Simpler**
- **More transparent**
- **Easier to reason about**
- **Fully observable and reproducible**

Instead of relying on heavyweight orchestration platforms or opaque runtime behavior, ETLX embraces a **declarative, metadata-first approach**.

At its core, ETLX lets you define *what should happen* — and then executes it deterministically, while capturing **rich execution metadata** along the way.

---

## Why ETLX?

Modern data stacks are powerful — but often **overengineered**.

Common pain points ETLX aims to solve:

- Pipelines that are hard to debug
- Logic scattered across code, configs, and orchestration tools
- Poor visibility into *what actually happened* during execution
- Documentation that drifts away from reality
- Vendor lock-in and engine-specific implementations

ETLX takes a different path.

<!-- markdownlint-disable MD026 -->

### ETLX is built on these principles:

- **Code-first, configuration-driven**
- **Database-centric**, powered by DuckDB
- **Composable pipelines**, not monolithic DAGs
- **Local-first**, but production-ready
- **Metadata as a first-class citizen**

---

## Core Capabilities

<!-- markdownlint-disable MD033 -->

<div class="row flex-xl-wrap pb-4">

  <!-- Declarative Pipelines -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/etl-elt#example">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">schema</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">Declarative Pipelines</p>
          <p class="para card-text mb-0">
            Define what should happen, not how.
            Pipelines are structured, readable, and reproducible.
          </p>
        </div>
      </div>
    </a>
  </div>

  <!-- DuckDB at the Core -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/ddb-at-the-core">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">database</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">DuckDB at the Core</p>
          <p class="para card-text mb-0">
            SQL-first transformations and in-process analytics
            powered by DuckDB.
          </p>
        </div>
      </div>
    </a>
  </div>

  <!-- Multi-Engine Execution -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/multi-engine-execution">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">hub</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">Multi-Engine Execution</p>
          <p class="para card-text mb-0">
            Run pipelines on DuckDB, PostgreSQL, SQLite,
            MySQL, SQL Server, and ODBC sources.
          </p>
        </div>
      </div>
    </a>
  </div>

  <!-- Execution & Observability -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/logs">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">visibility</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">Full Observability</p>
          <p class="para card-text mb-0">
            Every execution captures timings, validations,
            warnings, memory usage, and failure context.
          </p>
        </div>
      </div>
    </a>
  </div>

  <!-- Metadata & Lineage -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/query-documentation#documenting-fields-with-metadata">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">description</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">Metadata & Governance Doc</p>
          <p class="para card-text mb-0">
            Pipelines double as metadata documents,
            enabling lineage, dictionaries, and governance.
          </p>
        </div>
      </div>
    </a>
  </div>

  <!-- Beyond ETL -->
  <div id="list-item" class="col-md-4 col-12 py-2">
    <a class="text-decoration-none text-reset" href="../features/beyond-etl">
      <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
        <span class="h1 icon-color">
          <i class="material-icons align-middle">analytics</i>
        </span>
        <div class="card-body p-0 content">
          <p class="fs-5 fw-semibold card-title mb-1">Beyond ETL</p>
          <p class="para card-text mb-0">
            Use ETLX for reporting, document generation,
            structured exports, and regulatory workflows.
          </p>
        </div>
      </div>
    </a>
  </div>

</div>

---

## How ETLX Is Different

ETLX is **not**:

- A workflow scheduler
- A GUI-driven orchestration platform
- A black-box ETL tool

ETLX **is**:

- A deterministic execution engine
- A metadata-driven pipeline framework
- A foundation for documentation, lineage, and governance
- A tool that scales from local development to production

> **Your pipeline configuration is the source of truth.**

---

## Who Is ETLX For?

ETLX is designed for:

- Data engineers who prefer **clarity**
- Analytics engineers building SQL-first pipelines
- Teams that want **observable, auditable execution**
- Organizations that care about **governance and lineage**
- Anyone tired of pipelines that are hard to explain six months later

---

## Where to Go Next

- 👉 [Quickstart]({{% relref "quickstart" %}}) — run your first pipeline
- 👉 [Core Concepts]({{% relref "features/beyond-etl" %}}) — understand pipelines and metadata
- 👉 [Execution & Observability]({{% relref "features/logs" %}}) — see what ETLX captures
- 👉 [Configuration Reference]({{% relref "features/etl-elt#example" %}}) — full schema details

---

ETLX is open source and evolving.

If you want to understand **what your pipelines are doing, why they ran, and how to trust them**, you’re in the right place.
