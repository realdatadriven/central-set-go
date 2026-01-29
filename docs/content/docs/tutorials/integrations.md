---
weight: 7050
title: "Integrations"
description: "Learn how to connect Central Set with multiple databases, data sources, and data platforms using ETLX to integrate, transform, and securely expose data."
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

{{% alert context="warning" text="**Caution** — Central Set is **production-ready and actively used**, but this documentation is still **under active development**. Large parts of the docs are **auto-generated**, evolving alongside the platform, and some sections may be incomplete, rough around the edges, or change frequently." /%}}

## Overview

Central Set uses **ETLX** as its data integration and transformation engine.

When used as recommended — with **DuckDB as the query engine** — ETLX gives Central Set the ability to integrate data from **a vast ecosystem of data sources**, normalize and transform it in a secure and reproducible way, and expose it through **standardized consumption protocols** such as **Apache Arrow Flight** or **OData v4**.

The result is a **unified, secure, analytics-ready data layer** that can sit in front of operational systems, data lakes, files, APIs, and cloud platforms.

---

## ETLX + DuckDB: The Integration Core

ETLX is designed around **declarative, SQL-first data integration**, using DuckDB as its execution engine.

DuckDB provides:
- High-performance analytical execution
- Native support for Apache Arrow
- A rich extension ecosystem
- Embedded, serverless-friendly deployment

ETLX builds on top of this by:
- Managing connections and credentials securely
- Orchestrating ingestion, transformation, and exposure
- Making integrations reproducible and environment-safe

🔗 **ETLX**: https://github.com/realdatadriven/etlx  
🔗 **DuckDB**: https://duckdb.org

---

## Supported Data Sources

Through DuckDB core extensions and community extensions, Central Set can integrate data from **many (and growing) sources**.

### File-Based Sources

Native or extension-based support for:
- **Parquet**  
- **CSV**
- **JSON**
- **Avro**
- **Spreadsheets (Excel / Sheets-like formats)**
- **Vortex**

🔗 DuckDB File Formats:  
https://duckdb.org/docs/data/overview

---

### Relational Databases

Supported via native connectors and extensions:
- **PostgreSQL**
- **MySQL**
- **SQLite**
- **ODBC** (planned / evolving)

🔗 DuckDB PostgreSQL Extension:  
https://duckdb.org/docs/extensions/postgres  

🔗 DuckDB MySQL Extension:  
https://duckdb.org/docs/extensions/mysql  

🔗 DuckDB SQLite Scanner:  
https://duckdb.org/docs/extensions/sqlite

---

### Open Table & Lakehouse Formats

Modern analytics and lakehouse integrations:
- **Apache Iceberg**
- **Delta Lake**
- **Unity Catalog**

🔗 DuckDB Iceberg Extension:  
https://duckdb.org/docs/extensions/iceberg  

🔗 DuckDB Delta Extension:  
https://duckdb.org/docs/extensions/delta

---

### Arrow, Flight & ADBC Ecosystem

Deep integration with the Arrow ecosystem:
- **Apache Arrow**
- **Apache Arrow Flight**
- **ADBC Drivers**

These are foundational for how Central Set exposes data efficiently and safely.

🔗 Apache Arrow:  
https://arrow.apache.org  

🔗 Arrow Flight:  
https://arrow.apache.org/docs/format/Flight.html  

🔗 ADBC:  
https://arrow.apache.org/adbc/

---

### Community & Advanced Integrations

Through DuckDB community extensions and ETLX adapters, the ecosystem expands almost endlessly:

- **Web APIs**
- **OData**
- **BigQuery**
- **Snowflake**
- **Google Sheets**
- **Microsoft SQL Server**
- **SSH-based sources**
- **Web scraping**
- **XML / HTML**
- **Custom connectors**

This allows Central Set to act as a **data unification layer**, even across heterogeneous, non-traditional sources.

🔗 DuckDB Community Extensions:  
https://duckdb.org/docs/extensions/overview

---

## Secure Credential Management

All integrations are designed with **security-first principles**.

ETLX supports:
- Environment variable injection (`@ENV_VAR_NAME`)
- DuckDB `CREATE SECRET` syntax
- Token- and key-based authentication
- No hardcoded credentials in configs or SQL

This makes configurations safe to:
- Commit to version control
- Share across environments
- Deploy in containerized or cloud setups

---

## Transform Once, Expose Anywhere

Once data is integrated and transformed, Central Set can **expose it in multiple standardized ways**.

### Apache Arrow Flight

- High-performance, columnar transport
- Language-agnostic (Python, Go, Java, Rust, etc.)
- Ideal for analytics, ML, and large-scale consumers
- Secure, token-authenticated access

Arrow Flight is the **primary and recommended exposure protocol**.

---

### OData v4

- REST-based, standards-compliant API
- Ideal for BI tools, dashboards, and lightweight consumers
- Supports filtering, projection, and pagination
- Best suited for smaller or interactive workloads

---

## Access Control & Data Scoping

All exposed data respects **Central Set’s access rules**:

- **Authorization header is required**
- Access tokens are created in:
  - **Admin → Access Keys**
- Tokens are scoped to users
- Users must have access to the underlying resources
- **Row-level access rules**, when defined, are enforced

This allows Central Set to:
- Securely expose shared datasets
- Scope data per tenant, team, or user
- Reuse the same integrations safely across consumers

---

## Why This Matters

By combining:
- ETLX for integration and orchestration
- DuckDB for execution
- Arrow Flight and OData for exposure
- Strong access control and security

Central Set becomes a **single, unified gateway** between raw data sources and final consumers — without forcing data duplication, fragile pipelines, or vendor lock-in.

---

👉 For hands-on examples, configuration files, and real-world setups, see the documentation site:  
**https://realdatadriven.github.io/central-set-go/**
