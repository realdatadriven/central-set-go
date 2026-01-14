---
weight: 300
date: "2026-01-03T10:00:00+00:00"
draft: false
title: "Quickstart"
icon: "rocket_launch"
description: "Get Central Set Go (CSGO) running in minutes with a connected database admin UI and optional ETLX support."
publishdate: "2026-01-03T10:00:00+00:00"
tags: ["Beginners", "Quickstart", "Admin", "Databases"]
categories: ["Getting Started"]

twitter:
  card: "summary"
  title: "CSGO Quickstart"
  description: "Initialize and run Central Set Go locally in minutes"
---


## 🚀 Quickstart

This guide will help you **run Central Set Go (CSGO)** locally in minutes — initialize the admin database, access the admin UI, and optionally set up the ETLX subsystem.

---

## ✅ Requirements

### Minimum

* **Linux, macOS, or Windows**
* **A SQL database engine** (SQLite, PostgreSQL, MySQL, etc.)
* CSGO treats the **database itself as the data model**

### Optional (for ETLX / ODBC features)

* **unixODBC** (Linux/macOS) for ODBC-based ETL sources

---

## 🐙 Installation

You can run CSGO in **three different ways**:

1. Precompiled binary (recommended)
2. Build from source
3. Docker

Choose the option that best fits your workflow.

---

## ▶️ Option 1: Download a Precompiled Binary (Recommended)

Download the **latest CSGO release** for your platform:

👉 https://github.com/realdatadriven/central-set-go/releases

Make it executable (Linux/macOS):

```bash
chmod +x central-set
```

---

## 🛠️ Option 2: Build from Source (Git Clone)

If you want to build CSGO yourself:

### Requirements

* **Go ≥ 1.21**
* **git**

### Clone and build

```bash
git clone https://github.com/realdatadriven/central-set-go.git
cd central-set-go
go build -o central-set ./cmd
```

Run it:

```bash
./central-set --help
```

This produces the same binary as the official releases.

---

## 🐳 Option 3: Run with Docker

CSGO can be run entirely via Docker.

### Build the image

```bash
docker build -t central-set-go:latest .
```

### Run CSGO

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
./central-set --init
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

---

## ▶️ Start CSGO

```bash
./central-set
```

CSGO will start a web server and expose the admin UI.

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
./central-set --init --dbname ETLX
```

This creates an ETLX-powered app that integrates with:

* Pipeline execution
* Observability
* Dashboards

---

## ⚙️ Configure with `.env`

CSGO ships with a sample environment file:

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

> 🧠 Any database supported by **sqlx** can be used by simply changing the driver and DSN.

---

## 🧩 What Just Happened

CSGO has now:

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

## 🧠 Next Steps

* 👉 **Admin & UI** — How Apps → Menus → Tables define the UI
* 👉 **Security & API Access** — Users, roles, access keys
* 👉 **ETLX & Pipelines** — Data workflows and scheduling
* 👉 **Dashboards & Analytics** — DuckDB-powered insights

---

CSGO is **MIT-licensed**, **open source**, and designed to make the **database the single source of truth**.

If you believe admin UIs, APIs, pipelines, and analytics should all align around the schema — welcome aboard 🚀
