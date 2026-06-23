---
weight: 300
date: "2026-01-03T10:00:00+00:00"
draft: false
title: "Quickstart"
icon: "rocket_launch"
description: "Get Central Set (CS) running in minutes with a connected database admin UI and optional ETLX support."
publishdate: "2026-01-03T10:00:00+00:00"
tags: ["Beginners", "Quickstart", "Admin", "Databases"]
categories: ["Getting Started"]

twitter:
  card: "summary"
  title: "CS Quickstart"
  description: "Initialize and run Central Set locally in minutes"
---


## 🚀 Quickstart

This guide will help you **run Central Set (CS)** locally in minutes — initialize the admin database, access the admin UI, and optionally set up the ETLX subsystem.

---

## ✅ Requirements

### Minimum

* **Linux, macOS, or Windows**
* **A SQL database engine** (SQLite, PostgreSQL, MySQL, etc.)
* CS treats the **database itself as the data model**

### Optional (for ETLX / ODBC features)

* **unixODBC** (Linux/macOS) for ODBC-based ETL sources

---

## 🐙 Installation

You can run CS in **three different ways**:

1. Precompiled binary (recommended)
2. Build from source
3. Docker

Choose the option that best fits your workflow.

---

### ▶️ Option 1: Download a Precompiled Binary (Recommended)

#### Recommended: One-Line Installer (Linux & macOS)

The fastest way to install the latest release is:

```bash
curl -fsSL https://realdatadriven.github.io/central-set-go/install.sh | sh
```

The installer will:

* Detect your operating system and architecture
* Download the latest Central Set release
* Install the `c7` binary
* Verify the installation

Verify the installation:

```bash
./c7 --help
```

> 💡 The installer is the recommended method for most users because it always installs the latest stable release.

#### Manual download

Download the **latest CS release** for your platform:

👉 https://github.com/realdatadriven/central-set-go/releases/latest

Make it executable (Linux/macOS):

```bash
chmod +x central-set
```

Run it:

```bash
./central-set --help
```

---

### 🛠️ Option 2: Build from Source (Git Clone)

If you want to build CS yourself:

### Requirements

* **Go ≥ 1.26**
* **git**

### Clone and build

```bash
git clone https://github.com/realdatadriven/central-set-go.git
cd central-set-go
go mod tidy
go build -tags="duckdb_arrow" -o central-set ./cmd/api
```

Run it:

```bash
./central-set --help
```

This produces the same binary as the official releases.

---

## 🐳 Option 3: Run with Docker

CS can be run entirely via Docker.

### Build the image

```bash
docker build -t central-set-go:latest .
```

### Run CS

```bash
docker run --rm -it \
  -p 4444:4444 \
  -v $(pwd)/database:/app/database \
  central-set-go:latest
```

To initialize the admin database:

```bash
docker run --rm -it \
  -v $(pwd)/database:/app/database \
  central-set-go:latest --init
```

> 💡 Mounting the `database/` directory ensures your admin DB persists between runs.

---

## 🗄️ Initialize the **Admin Database**

Before first use, initialize the admin database:

```bash
./central-set --init --model admin_model.md
```

This will:

* Create the internal **admin database**
* Create default configuration tables (`apps`, `menus`, `tables`, etc.) for the admin ui
* Create a **default user**
* Print credentials:

```
Username: root
Password: 1234
```

### 🏦 Model-Based Initialization

Central Set follows a **model-driven database approach**, inspired by patterns commonly used.

Instead of manually creating tables or running raw SQL migrations, the recommended workflow is to:

1. Define your **application model** in a Markdown model file
2. Describe entities, fields, relationships, and metadata
3. Let Central Set **generate and manage the database structure**

This approach provides several advantages:

* **Consistent database structures**
* **Self-documented schema**
* **Reproducible environments**
* **Tighter integration with Central Set APIs and metadata**
* **Safer schema evolution**

Although the example above initializes the **Admin database**, this approach is **not limited to it**.

Any application database can — and ideally **should** — be initialized using the same **model-based workflow**.

> In fact, new backend projects built with Central Set are **recommended to start with a model file**, allowing the platform to manage the database schema from the beginning.

A dedicated page will explore the **Model system and architecture** in more detail.

---

## ▶️ Start CS

```bash
./central-set
```

CS will start a web server and expose the admin UI.

---

## 🖥️ Open the Admin UI

👉 [http://localhost:4444](http://localhost:4444)

Login using:

```
Username: root
Password: 1234
```

---

## 🧬 Optional: Initialize ETLX Support

To enable **ETLX pipelines, notebooks, and SQL tools**, initialize an additional app database:

```bash
./central-set --init --model etlx_model.md
```

This creates an ETLX-powered app that integrates with:

* Pipeline execution
* Observability
* Dashboards

---

## ⚙️ Configure with `.env`

CS ships with a sample environment file:

```bash
cp dot-env-example.txt .env
```

Edit `.env` to configure:

* Database driver & DSN
* HTTP port
* Authentication & security
* ETLX options

### Example

```env
# Admin database
DB_DRIVER_NAME=sqlite3
DB_DSN=database/ADMIN.db

# HTTP server
HTTP_PORT=4444
```

> Any database supported by **sqlx** can be used by simply changing the driver and DSN.

---

## 🧩 What Just Happened

CS has now:

* Initialized its **admin control database**
* Reflected database tables into:

  * Auto-generated data tables
  * CRUD forms
  * APIs
* Enabled **RBAC at table level**
* Exposed everything through a **secure, API-first backend**
* Optionally enabled **ETLX-powered pipelines and analytics**

You can now:

* Manage applications, menus, and tables
* Customize UI layouts and forms
* Generate API keys and tokens
* Schedule ETL jobs
* Build dashboards and notebooks

---

## Steps

* 👉 **Admin & UI** — How Apps → Menus → Tables define the UI
* 👉 **Security & API Access** — Users, roles, access keys
* 👉 **ETLX & Pipelines** — Data workflows and scheduling
* 👉 **Dashboards & Analytics** — DuckDB-powered insights

---

CS is **MIT-licensed**, **open source**, and designed to make the **database the single source of truth**.

If you believe admin UIs, APIs, pipelines, and analytics should all align around the schema — welcome aboard 🚀
