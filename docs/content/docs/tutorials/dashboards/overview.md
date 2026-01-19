---
weight: 7081
title: "Overview - Analytics & Dashboards"
description: "Overview - Analytics & Dashboards"
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Overview

Central Set provides a **first-class analytics and dashboarding layer** by embedding and extending
[evidence.dev](https://docs.evidence.dev).

This allows you to build **interactive analytical dashboards** directly on top of:

* ETLX pipelines
* Exported datasets
* Parquet files
* Curated analytical tables

All dashboards live **inside Central Set**, inherit authentication and authorization, and are fully
integrated with the platform’s data model.

---

## Why Evidence.dev?

Evidence.dev was chosen because it is:

* SQL-first
* Component-based
* Open-source
* Developer-friendly
* Designed for analytical storytelling

However, Central Set **does not use Evidence in its default standalone mode**.

Instead, Evidence is **embedded as a dashboard frontend**, while Central Set controls:

* Configuration
* Data access
* Execution context
* Security
* Dataset lifecycle

In Central Set, Evidence configuration is rendered **in real time in the browser**, rather than being
compiled into a static website.

Conceptually, Evidence behaves as a **Svelte-based rendering layer** that intercepts normal CRUD-based
UI rendering and transforms database-backed records into analytical dashboards.

Some Evidence components are not yet integrated or are slightly modified to better fit Central Set,
but **the core visualization components (what really matters) are fully compatible**, including:

* [If / Else](https://docs.evidence.dev/core-concepts/if-else/)
* [Loops](https://docs.evidence.dev/core-concepts/loops/)
* [Formating](https://docs.evidence.dev/core-concepts/formatting/)
* [Value](https://docs.evidence.dev/components/data/value/)
* [Big Value](https://docs.evidence.dev/components/data/big-value/)
* [Data Table](https://docs.evidence.dev/components/data/data-table/)
* [Delta](https://docs.evidence.dev/components/data/delta/)
* [Charts](https://docs.evidence.dev/components/charts/area-chart/)

Charts are based on [Apache ECharts](https://echarts.apache.org/en/index.html) and can be fully customized using:

* [ECharts Extra Options](https://docs.evidence.dev/components/charts/echarts-options/)
* [Custom Charts](https://docs.evidence.dev/components/charts/custom-echarts/)

This allows the full ECharts configuration object to be passed directly, enabling virtually any
visualization supported by ECharts.

---

## What’s Different from Vanilla Evidence

| Area            | Evidence.dev                | Central Set             |
| --------------- | --------------------------- | ----------------------- |
| Data Sources    | Inline SQL / DB connections | Central Set APIs        |
| Authentication  | None / external             | Native Central Set auth |
| Dataset Storage | Ad-hoc queries              | ETLX datasets & Parquet |
| Runtime         | Static site / dev server    | Embedded application    |

In short:

> **Evidence provides the rendering engine.
> Central Set provides the platform.**

---

## Data Flow: ETLX → Dashboards

The dashboard layer is designed to sit **downstream of ETLX**.

### Typical Flow

```
ETLX Pipeline
   ↓
Dataset / Parquet Export
   ↓
Central Set Dataset Registry
   ↓
Evidence Dashboard Components
```

ETLX pipelines produce:

* Parquet files
* Materialized datasets
* Analytical views

These outputs are **registered and managed** by Central Set and exposed to dashboards as **named datasets**.

Dashboards never access raw databases directly.

---

## Dataset-Centric Architecture

Instead of writing SQL directly inside dashboards:

* Datasets are **defined once**
* Produced by ETLX
* Versioned and traceable
* Reused across dashboards

This ensures:

* Reproducibility
* Performance
* Clear lineage
* Separation of concerns

> Dashboards consume data.
> Pipelines produce data.

---

## Embedded Evidence Frontend

Central Set embeds the Evidence dashboard frontend directly into the UI.

This means:

* Dashboards appear as native Central Set pages
* Navigation is unified
* Authentication is shared
* Permissions apply automatically

From the user’s perspective, dashboards are simply **another application menu**.

---

## Extended Components

In addition to standard Evidence components, Central Set adds its own components, such as `<Stats>`.

Most of these are **Svelte components built with**
[DaisyUI](https://daisyui.com/components/), since Central Set itself is built on the
[DaisyUI](https://daisyui.com) design system.

This ensures visual consistency between dashboards and the rest of the platform UI.

---

## Configuration Model

Evidence dashboards normally rely on a local configuration block.

In Central Set:

* Configuration is **metadata-driven**
* Stored and versioned centrally
* Linked to applications and menus
* Managed from the ADMIN app

This allows dashboards to be:

* Activated or deactivated
* Versioned
* Assigned to profiles
* Shipped as part of an application

---

## Security & Governance

Because dashboards are embedded:

* Access is controlled by Central Set profiles
* Users only see datasets they are authorized to access
* Dashboards can be restricted by role or tenant

This makes dashboards suitable for:

* Internal analytics
* Multi-tenant environments
* Customer-facing data products

---

## Analytics as a Product

Dashboards are not just visualizations.

They are:

* Versioned artifacts
* Backed by pipelines
* Governed by metadata
* Deployable across environments

This aligns with Central Set’s core philosophy:

> **Everything is configuration.
> Everything is data-driven.**

---

## How This Fits Together

| Layer      | Responsibility                   |
| ---------- | -------------------------------- |
| ETLX       | Data pipelines & transformations |
| Datasets   | Canonical analytical outputs     |
| ADMIN      | Metadata, permissions, structure |
| Dashboards | Exploration & insight            |
| Evidence   | Rendering engine                 |
