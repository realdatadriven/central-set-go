---
weight: 870
title: "Go API & Programmatic Usage"
description: "Use ETLX as a Go library to execute pipelines, run specific stages, or embed ETLX into your own applications."
icon: integration_instructions
tags: ["golang", "api", "embedding", "runtime"]
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
toc: true
---

## Go API & Programmatic Usage

ETLX is **not only a CLI tool** — it is also a **Go library** that you can embed directly into your own applications.

The official `etlx` binary is built using **the same public Go APIs** exposed by the project.
This page documents those APIs in a **ready-to-use** way, based on how `cmd/main.go` constructs and executes pipelines.

## 🔧 CLI Internals (Important Context)

The ETLX CLI:

* Uses **Go’s standard `flag` package**
* Does **not** rely on Cobra or subcommands
* Internally calls the same APIs you can call from Go

This means:

> 👉 **Anything the CLI can do, you can do programmatically**

## 📦 Installing as a Go Dependency

```bash
go get github.com/realdatadriven/etlx
```

## 🧠 Core Concepts (Go Perspective)

At runtime, ETLX works in three phases:

1. **Parse Markdown configuration**
2. **Resolve execution scope** (date, only, skip, steps, etc.)
3. **Execute pipeline + collect logs**

The Go API mirrors this exactly.

## 🏗 Creating an ETLX Engine

```go {linenos=table}
import "github.com/realdatadriven/etlx"
```

```go {linenos=table}
engine := &etlx.ETLX{}
```

This `ETLX` instance holds:

* Parsed pipeline configuration
* Execution options
* Runtime state
* Execution logs

## 📄 Loading a Pipeline Configuration

### Load from file

```go {linenos=table}
err := engine.ConfigFromFile("pipeline.md")
if err != nil {
    panic(err)
}
```

### Load from in-memory Markdown

Useful for dynamic or generated pipelines:

```go {linenos=table}
md := `# ETLX Pipeline...`

err := engine.ConfigFromMDText(md)
if err != nil {
    panic(err)
}
```

## 📆 Running a Pipeline (Equivalent to CLI)

Equivalent CLI:

```bash
etlx --config pipeline.md --date 2025-01-01
```

### Minimal execution

```go {linenos=table}
dates := []time.Time{
    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
}

logs, err := engine.RunETL(dates, nil, nil)
if err != nil {
    panic(err)
}
```

## 🎯 Execution Options (`--only`, `--skip`, `--steps`)

The CLI flags are translated into a simple options map.

### Run only specific keys

```go {linenos=table}
opts := map[string]any{
    "only": []string{"sales", "customers"},
}

logs, err := engine.RunETL(dates, nil, opts)
```

Equivalent CLI:

```bash
etlx --config pipeline.md --only sales,customers
```

### Skip specific keys

```go {linenos=table}
opts := map[string]any{
    "skip": []string{"debug_tables"},
}
```

### Run specific lifecycle steps

```go {linenos=table}
opts := map[string]any{
    "steps": []string{"extract", "load"},
}
```

Equivalent CLI:

```bash
etlx --config pipeline.md --steps extract,load
```

## 🧹 Clean & Drop Operations

Equivalent CLI:

```bash
etlx --config pipeline.md --clean
etlx --config pipeline.md --drop
```

### Clean SQL blocks

```go {linenos=table}
opts := map[string]any{
    "clean": true,
}

engine.RunETL(dates, nil, opts)
```

### Drop SQL blocks

```go {linenos=table}
opts := map[string]any{
    "drop": true,
}

engine.RunETL(dates, nil, opts)
```

## 📊 Accessing Execution Logs

Every execution returns **structured logs**.

```go {linenos=table}
logs, _ := engine.RunETL(dates, nil, nil)

for _, log := range logs {
    fmt.Println(
        log.Name,
        log.Ref,
        log.Duration,
        log.Success,
        log.Msg,
    )
}
```

Each log entry typically contains:

| Field      | Description        |
| ---------- | ------------------ |
| `name`     | Step or block name |
| `ref`      | Reference date     |
| `start`    | Start timestamp    |
| `end`      | End timestamp      |
| `duration` | Execution time     |
| `success`  | Success / failure  |
| `msg`      | Message or error   |

## 🧪 Running Specific Sections

ETLX exposes **explicit execution entry points** for each pipeline capability.

All functions share the same execution model:

```go {linenos=table}
func(dates []time.Time, cfg any, opts map[string]any) ([]etlx.Log, error)
```

The difference is **which Markdown sections are executed**.

## ▶️ RunETL

Executes the full ETL / ELT lifecycle see [ETL / ELT]({{% relref "etl-elt" %}})

```go {linenos=table}
logs, err := engine.RunETL(dates, nil, nil)
```

## 🧪 RunDATA_QUALITY

Runs operational or auxiliary scripts see [Data Quality]({{% relref "data-quality" %}})

```go {linenos=table}
logs, err := engine.RunDATA_QUALITY(dates, nil, nil)
```

## 📤 RunEXPORTS

Executes export definitions see [Exports]({{% relref "exports" %}})

```go
logs, err := engine.RunEXPORTS(dates, nil, nil)
```

## 📜 RunSCRIPTS

Runs operational or auxiliary scripts see [Scripts]({{% relref "scripts" %}})

```go {linenos=table}
logs, err := engine.RunSCRIPTS(dates, nil, nil)
```

## ⚙️ RunACTIONS

Runs post-processing or orchestration actions see [Actions]({{% relref "actions" %}})

```go {linenos=table}
logs, err := engine.RunACTIONS(dates, nil, nil)
```

## 🧾 RunLOGS

Processes logging or audit blocks see [Logs / Observability]({{% relref "logs" %}})

```go {linenos=table}
logs, err := engine.RunLOGS(dates, nil, nil)
```

## 🔔 RunNOTIFY

Executes notification definitions see [Notificatios]({{% relref "notify" %}})

```go {linenos=table}
logs, err := engine.RunNOTIFY(dates, nil, nil)
```

## 🧵 Using ETLX Inside Long-Running Services

```go {linenos=table}
func runPipeline(ref time.Time) error {
    engine := &etlx.ETLX{}

    if err := engine.ConfigFromFile("pipeline.md"); err != nil {
        return err
    }

    _, err := engine.RunETL([]time.Time{ref}, nil, nil)
    return err
}
```

## 🧠 Why This Matters

Because ETLX is:

* **A runtime**
* **An open specification**
* **A Go library**

You are not locked into:

* A closed binary
* A hidden execution model
* A proprietary DSL

The **Markdown pipeline *is* the documentation**, and the **Go API is the execution layer**.

## 🏁 Summary

✔ CLI uses **standard Go flags**
✔ Go API mirrors the CLI exactly
✔ Pipelines can be executed from files or memory
✔ Each lifecycle has its own entry point
✔ Structured logs for audit & observability
✔ Ideal for embedding, automation, and regulated workloads
