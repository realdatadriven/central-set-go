# Central Set (CS)

**Central Set** is a **configuration-driven data platform & Admin UI** written in **Golang**.

It unifies **application management, databases, users, roles, ETL workflows, dashboards, and analytics**
into a single system — driven by metadata instead of custom code.


---

## What is Central Set?

Central Set acts as a **control plane** for data-driven applications:

- Manage applications, schemas, and databases
- Configure CRUD UIs and admin workflows
- Build and run [ETLX](https://github.com/realdatadriven/etlx) data pipelines
- Expose data via APIs (REST-ish, OData v4, Arrow Flight)
- Enforce access control, row-level security, and governance
- [Create dashboards](https://realdatadriven.github.io/central-set-go/docs/tutorials/dashboards/) and reports using Markdown + DuckDB-WASM ([evidence.dev](https://docs.evidence.dev) style)

All of this is defined using **configuration**, not hard-coded logic.

---

## Key Capabilities

- ⚙️ **Configuration-driven platform**
- 🧩 **Built-in Admin UI**
- 🗄️ **Multi-database support** (SQLite, PostgreSQL, MySQL, SQL Server)
- 🌳 [**ETLX-powered data workflows**](https://realdatadriven.github.io/central-set-go/docs/tutorials/etlx/)
- 📊 **Markdown-based dashboards (DuckDB-WASM)**
- 🔌 **API-first design**
- 🔑 **Token & role-based access control**
- 🔍 **OData v4 read-only query API**
- 🚀 **Arrow Flight data sharing**

---

## Documentation

📘 **Full documentation lives here:**

👉 **https://realdatadriven.github.io/central-set-go/**

> 🚧 **Note**  
> Central Set is production-ready, but documentation is still evolving and partially auto-generated, under evaluation and may change.

The docs cover:

- Architecture & concepts
- [Quickstart guides](https://realdatadriven.github.io/central-set-go/docs/quickstart/)
- Configuration reference
- [ETLX workflows](https://realdatadriven.github.io/central-set-go/docs/tutorials/etlx/)
- API reference
- [Arrow Flight Supporte](https://realdatadriven.github.io/central-set-go/docs/tutorials/arrow-flight/)
- [OData v4 API](https://realdatadriven.github.io/central-set-go/docs/tutorials/odata-v4/)
- Security & access control
- Dashboards & analytics

---

## Quick Links

- 📖 Docs: https://realdatadriven.github.io/central-set-go/
- 💻 Source: https://github.com/realdatadriven/central-set-go
- 🧪 ETLX: https://github.com/realdatadriven/etlx

---

## Philosophy

> **Configuration is the product.**  
> Code executes it.

Central Set is designed to scale **governance, analytics, and automation**
without multiplying services or frameworks.

---

## License

MIT License
