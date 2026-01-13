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

### What is ETLX?

**ETLX** is an **open-source, SQL-first data workflow engine** and an **evolving open specification** for describing complete data workflows as **executable documentation**.

ETLX lets you define ETL, analytics, reporting, exports, etc using **Markdown + YAML + SQL**, where the pipeline itself becomes the documentation, governance artifact, and execution plan.



### What problem does ETLX solve?

ETLX addresses common problems in modern data workflows:

* Hidden logic spread across code, SQL, and orchestration tools
* Poor documentation and weak data governance
* Tight coupling between pipelines and execution engines
* Difficult auditing, reproducibility, and compliance

ETLX makes **all logic explicit**, **self-documenting**, and **auditable by design**.



### Who is ETLX for?

ETLX is designed for:

* **Data engineers** building ETL / ELT pipelines
* **Analytics engineers** working SQL-first
* **Data scientists** needing reproducible / automated data workflows
* **Data analysts** creating reports and dashboards
* **Platform and data architects**
* **Governance, compliance, and data quality teams**
* **Organizations** that value transparency and reproducibility

If you believe *SQL should be the transformation language and documentation should not be an afterthought*, ETLX is for you.



### How do I get started?

Start with the documentation:

👉 **Getting Started & Concepts**
[https://realdatadriven.github.io/etlxdocs](https://realdatadriven.github.io/etlxdocs)

You can:

* Run ETLX via the CLI
* Embed it as a Go library
* Start with a single Markdown file and grow from there



## Technical Questions

### What technologies does ETLX use?

ETLX is built primarily in **Go** and is:

* **SQL-first**
* **Markdown-driven**
* **DBMS Engine-agnostic**

It is powered by **DuckDB** for in-process analytics, but also supports:

* PostgreSQL
* SQLite
* MySQL
* SQL Server
* ODBC-compatible databases

ETLX does **not** introduce a proprietary DSL.



### Is ETLX tied to DuckDB?

No.

DuckDB is the default and recommended engine due to its:

* In-process execution
* Multi-source querying
* Performance

However, ETLX is designed to support **multiple SQL engines** and execution backends. DuckDB is a strength — not a lock-in.



### Is ETLX just a runtime?

No.

ETLX is:

* A **runtime**
* A **configuration model**
* An **specification**

The same configuration can be used to:

* Execute workflows
* Generate documentation
* Produce governance artifacts like data dictionaries, data lineage, data quality rules, reports, ...



### How does ETLX handle documentation and governance?

ETLX parses Markdown into a structured model (`map[string]any` in Go), which includes:

* Dataset metadata
* Column-level definitions
* Ownership and lineage
* Validation rules
* Execution semantics

This enables automatic generation of:

* Data dictionaries
* Governance documentation
* Audit reports
* Quality checks

Without duplicating logic.



### How can I contribute?

ETLX is community-driven.

You can contribute by:

* Improving documentation
* Adding examples
* Submitting bug reports
* Proposing specification improvements
* Contributing code

👉 **Contribution Guide**
[https://realdatadriven.github.io/etlxdocs/docs/contributing/](https://realdatadriven.github.io/etlxdocs/docs/contributing/)



### Where is the source code?

The source code is hosted on GitHub:

👉 [https://github.com/realdatadriven/etlx](https://github.com/realdatadriven/etlx)



## Support Questions

### How can I get help?

You can get support by:

* Reading the documentation
* Opening GitHub issues
* Participating in discussions

The community is encouraged to help shape ETLX’s evolution.



### Are there tutorials or examples?

Yes.

The documentation includes:

* Quickstart guides
* Core concepts
* Advanced examples
* Real-world patterns

👉 [https://realdatadriven.github.io/etlxdocs](https://realdatadriven.github.io/etlxdocs)



### How do I report a bug?

Please open an issue on GitHub:

👉 [https://github.com/realdatadriven/etlx/issues](https://github.com/realdatadriven/etlx/issues)

Include:

* A minimal example
* The ETLX version
* Expected vs actual behavior



## Licensing Questions

### What license does ETLX use?

ETLX is licensed under the **Apache License 2.0**.

This allows:

* Commercial use
* Modification
* Distribution

With proper attribution.



### Can ETLX be used commercially?

Yes.

ETLX is explicitly designed to be usable in **commercial and enterprise environments**.



### How should ETLX be attributed?

Please reference the project as:

**ETLX** — [https://github.com/realdatadriven/etlx](https://github.com/realdatadriven/etlx)
by **RealDataDriven**



## Future & Roadmap Questions

### Is ETLX stable?

ETLX is actively developed.

The **core ideas are stable**, while the **specification is evolving** with community feedback.



### Are there upcoming features?

Planned and ongoing work includes:

* Expanded open specification
* More advanced examples
* Additional export and governance templates
* Improved observability and validation primitives



### How can I suggest features?

Feature suggestions are welcome.

You can:

* Open a GitHub issue
* Start a discussion
* Propose changes to the specification



### Will ETLX receive regular updates?

Yes.

ETLX follows an incremental, transparent development approach.
Changes and improvements are tracked openly in GitHub. 
