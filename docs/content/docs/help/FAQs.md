---
weight: 2120
title: "FAQ"
icon: "quiz"
description: "Answers to frequently asked questions."
lead: "Answers to frequently asked questions."
date: 2025-10-06T08:49:31+00:00
lastmod: 2025-10-06T08:49:31+00:00
draft: false
images: []
toc: true
---

## General Questions

### What is Central-Set?

**Central-Set** is a **dynamic, SQL-driven application platform and data governance layer** designed to build secure data APIs, manage metadata, enforce access control, and expose analytical interfaces — all from configuration.

It combines:

- Dynamic REST APIs  
- ETLX-powered data workflows  
- OData v4 endpoints  
- Apache Arrow Flight support  
- Fine-grained access control (RLA / field-level security)  

All managed through a unified Admin interface.



### What problem does Central-Set solve?

Central-Set addresses common challenges in modern data platforms:

* Hard-coded APIs that are difficult to evolve  
* Poorly governed data access  
* Fragmented authentication and authorization  
* Limited metadata visibility  
* Complex integration between transformation engines and APIs  
* Inconsistent security enforcement across services  

Central-Set provides:

* A **dynamic API layer**
* Centralized **governance and access control**
* Built-in **multi-tenant security**
* Unified analytical data exposure



### Who is Central-Set for?

Central-Set is designed for:

* **Data platform engineers**
* **Backend engineers building data APIs**
* **Analytics engineers**
* **Data architects**
* **SaaS builders**
* **Governance and compliance teams**
* **Organizations building internal data platforms**

If you want to build **secure, metadata-driven, SQL-first APIs without hardcoding every endpoint**, Central-Set is for you.



### How do I get started?

Start with the documentation:

👉 https://realdatadriven.github.io/central-set-go/

You can:

* Run it as a standalone server
* Use SQLite for quick local setup
* Connect PostgreSQL or other databases
* Enable ETLX workflows
* Expose data via OData or Arrow Flight



## Technical Questions

### What technologies does Central-Set use?

Central-Set is built primarily in **Go** and uses:

* SQL databases (SQLite, PostgreSQL, MySQL, SQL Server)
* DuckDB (for ETLX and analytical workloads)
* Apache Arrow Flight
* OData v4
* JWT-based authentication
* Dynamic routing engine (`/dyn_api/{ctrl}/{act}`)

It is designed to be:

* Engine-agnostic
* Metadata-driven
* Configurable via environment variables



### What is the Dynamic API (`/dyn_api`)?

Central-Set exposes a dynamic endpoint:

```

POST /dyn_api/{ctrl}/{act}

````

Where:

- `ctrl` acts like a namespace (similar to a class)
- `act` acts like a method

Requests follow a standard structure:

```json
{
  "lang": "en",
  "app": { ... },
  "data": { ... }
}
````

Responses follow a standard format:

```json
{
  "success": true,
  "msg": "message",
  "data": { ... }
}
```

This allows new capabilities to be added without hardcoding new REST routes.

### Does Central-Set support fine-grained access control?

Yes.

Central-Set supports:

* Role-based access control
* Application-level access
* Table-level access
* Field-level access
* Row-Level Access (RLA)
* Scope-based filtering
* Token-based access keys

These rules apply consistently across:

* REST APIs
* OData v4
* Arrow Flight

### How does Central-Set integrate with ETLX?

Central-Set uses **ETLX** as its data integration engine.

ETLX allows:

* SQL-first data transformations
* Multi-source integration
* Metadata-driven workflows
* Governance artifact generation

ETLX outputs can be:

* Queried internally
* Exposed via OData
* Exposed via Arrow Flight
* Governed by Central-Set access rules

### Can Central-Set expose analytical data?

Yes.

Central-Set supports:

#### OData v4

* HTTP-based
* Filterable (`$filter`)
* Good for transactional and moderate workloads

#### Apache Arrow Flight

* High-performance
* Columnar
* Ideal for BI and data science tools
* Backed by DuckDB
* Supports schema/table/field/scope-level restrictions

### Is Arrow Flight secure?

Yes.

Arrow Flight requires:

* Authorization header
* Valid access token
* Proper role permissions
* Optional Row-Level Access

It can also run in:

* Direct DuckDB mode
* Application-aware mode (CRUD-layer enforced)

TLS is supported via:

```env
ENABLE_TLS=true
TLS_CERT_FILE=ssl/server-cert.pem
TLS_KEY_FILE=ssl/server-key.pem
TLS_CA_CERT_FILE=ssl/ca-cert.pem
```

### Can Central-Set be multi-tenant?

Yes.

Multi-tenancy is supported through:

* App isolation
* Database isolation
* Role-based permissions
* Row-Level Access
* Scoped tokens

### Is Central-Set just an API server?

No.

Central-Set is:

* A dynamic API engine
* A governance layer
* A metadata manager
* A data integration orchestrator (via ETLX)
* An analytical data server (via Arrow Flight)
* An OData provider

## Support Questions

### How can I get help?

You can get support by:

* Reading the documentation
* Opening GitHub issues
* Participating in discussions

👉 [https://github.com/realdatadriven/central-set-go](https://github.com/realdatadriven/central-set-go)

### Are there examples?

Yes.

The documentation includes:

* Configuration examples
* OData usage examples
* Arrow Flight examples
* ETLX integration examples
* Dynamic API templates

👉 [https://realdatadriven.github.io/central-set-go/](https://realdatadriven.github.io/central-set-go/)

### How do I report a bug?

Please open an issue:

👉 [https://github.com/realdatadriven/central-set-go/issues](https://github.com/realdatadriven/central-set-go/issues)

Include:

* Version
* Environment details
* Minimal reproduction
* Expected vs actual behavior

## Licensing Questions

### What license does Central-Set use?

Central-Set is licensed under the **Apache License 2.0**.

This allows:

* Commercial use
* Modification
* Distribution

With proper attribution.

### Can Central-Set be used commercially?

Yes.

Central-Set is designed for:

* SaaS platforms
* Internal enterprise data platforms
* Data product architectures
* Regulated environments

### How should Central-Set be attributed?

Please reference the project as:

**Central-Set** — [https://github.com/realdatadriven/central-set-go](https://github.com/realdatadriven/central-set-go)
by **RealDataDriven**

## Future & Roadmap Questions

### Is Central-Set stable?

Yes.

Central-Set is **production-ready and actively used**, although the documentation is still evolving.

### What is currently evolving?

Ongoing work includes:

* Enhanced Arrow Flight governance
* Improved schema introspection
* Expanded ETLX integration
* Performance optimizations
* Better dynamic API documentation generation

### How can I suggest features?

You can:

* Open a GitHub issue
* Start a discussion
* Propose architectural improvements
* Contribute documentation

### Will Central-Set receive regular updates?

Yes.

Central-Set follows an incremental, transparent development approach.
Changes and improvements are tracked publicly in GitHub.

