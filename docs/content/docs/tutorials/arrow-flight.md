---
weight: 7090
title: "Arrow Flight"
description: "Apache Arrow Flight Support API"
icon: arrow
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Arrow Flight Support

{{% alert context="warning" text="**Important** - Central-Set is production-ready, but this documentation is still under active development. Large parts are auto-generated and may change as features evolve / get documented." /%}}

### Overview

Central-Set provides [**Apache Arrow Flight**](https://arrow.apache.org/docs/format/Flight.html) support trought [airport-go](https://github.com/hugr-lab/airport-go) to expose analytical datasets through a **unified, high-performance interface**, primarily targeting **ETLX outputs**, but capable of serving **any datasource supported by DuckDB**.

The goal is to make transformed data **immediately consumable by analytical engines, BI tools, and data science workflows**, without exporting files or duplicating storage.

Arrow Flight endpoints are dynamically defined and managed through the **Admin UI**, backed by **in-memory DuckDB instances** that attach external datasources on demand.

It takes the Arrow Flight configuration defined in the Admin database and uses **airport-go** to create a Flight server backed by **DuckDB**.

At runtime, it initializes a DuckDB instance and serves the tables defined in the configuration using DuckDB's **Arrow integration**, exposing them through the Arrow Flight protocol.

The lifecycle of each endpoint is fully driven by configuration:

* `startup_sql` is executed to load extensions and initialize dependencies
* `main_sql` attaches the underlying data sources and exposes tables
* `shutdown_sql` is executed to clean up resources (for example, detaching databases)

All requests are authenticated using the `validateToken` function, ensuring that Arrow Flight follows the same security and access rules as the rest of the platform.

The original design aimed to reuse a **global `duckdb.Connector`** so multiple Flight clients could share the same DuckDB instance. However, due to concurrency issues when sharing connectors across goroutines, the current implementation creates a **new DuckDB connector per scan request**.

While this approach introduces some overhead, real-world testing shows that performance remains acceptable for analytical workloads—especially when compared to the OData v4 API, which is better suited for smaller or transactional queries.

Further optimizations around connector reuse and resource management are planned to improve performance and efficiency over time.

---

### Security & Access Control

All Arrow Flight access follows the **same security model as the REST API**:

* ✅ **Authorization is mandatory**
* ✅ An **access token must be provided** via the `Authorization` header
* ✅ Tokens are created in:

  * **Admin → Admin → Access Keys**
* ✅ The token **must belong to a user with access** to the requested resources
* ✅ **Row-Level Access (RLA)** rules apply automatically when defined

This makes Arrow Flight ideal for:

* Scoped analytical access
* Secure data sharing
* Multi-tenant environments
* Fine-grained filtering of exposed datasets

---

### How Arrow Flight Works in Central-Set

Each Arrow Flight endpoint is defined as a **schema configuration** stored in the **Admin database**.

At runtime:

1. An **in-memory DuckDB** instance is created
2. `startup_sql` is executed (load extensions, dependencies, etc.)
3. `main_sql` attaches external datasources and exposes tables
4. The schema is made available via Arrow Flight
5. On disconnect, `shutdown_sql` performs cleanup

This design allows:

* Dynamic datasource attachment
* Zero persistence inside the Flight server
* Clean lifecycle management per connection

---

### Enabling Arrow Flight

To start Arrow Flight alongside the REST API, set the following environment variables:

```env
ENABLE_ARROW_FLIGHT=true
ARROW_FLIGHT_ADDR=0.0.0.0:50051
```

Once enabled, Arrow Flight runs in the **same binary** as the REST API.

---

### Exposing a Schema

Create a new schema via:

**Admin → Expose Arrow Flight**

Example configuration:

```yaml
name: my_schema
description: schema_description
db_schema: main
startup_sql: |
  INSTALL SQLITE;
  LOAD SQLITE;
main_sql: |
  ATTACH 'database/test.db' AS my_schema (TYPE SQLITE);
  USE my_schema;
shutdown_sql: |
  USE memory;
  DETACH my_schema;
```

Each entry represents **one Arrow Flight schema**.

---

### Example: SQLite + TPC-H

You can generate a SQLite database with TPC-H data using DuckDB:

```bash
duckdb -c "
INSTALL tpch;
LOAD tpch;
CALL dbgen(sf=1);
ATTACH 'database/test.db' (TYPE SQLITE);
COPY FROM DATABASE memory TO test;
DETACH test;
"
```

(Alternatively, use any existing SQLite database.)

---

### Connecting from DuckDB (Recommended)

Use DuckDB's **airport community extension**:

```sql {linenos=table}
INSTALL airport FROM community;
LOAD airport;

CREATE OR REPLACE PERSISTENT SECRET airport_auth_secret (
    TYPE airport,
    AUTH_TOKEN 'your_access_token_here',
    SCOPE 'grpc://127.0.0.1:50051'
);

ATTACH '' AS my_server (TYPE AIRPORT, LOCATION 'grpc://127.0.0.1:50051');

SELECT *
FROM my_server.my_schema.orders
LIMIT 10;
```

This provides:

* Native Arrow Flight performance
* Zero-copy data access
* Seamless SQL integration

---

### Current Limitations & Roadmap

Arrow Flight support is functional and stable, but **still evolving**.

Planned improvements include:

* Explicit control over **which tables are exposed**
* Field-level visibility controls
* Optional **DML support** (INSERT / UPDATE / DELETE)
* Better schema introspection

The current focus is on **read-optimized analytical workloads**.

---

### Why Arrow Flight in Central-Set?

Arrow Flight enables Central-Set to act as a **data serving layer**, not just a transformation engine:

* ETLX outputs become instantly queryable
* No file exports or data duplication
* Works with modern analytical tooling
* Secure, scoped, multi-tenant by design

This makes it especially useful for **analytics, BI, and data science** use cases outside Central-Set, while keeping governance centralized.
