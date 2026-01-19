---
weight: 7082
title: "Example Dashboard"
description: "Analytics & Dashboards — Logs Dashboard Example"
icon: auto_awesome
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## Analytics & Dashboards — Logs Dashboard Example

This page walks through a real dashboard definition and explains **how configuration, filters, metrics, charts, and tables work together** to produce the final analytics UI.

This dashboard uses:

* **Evidence.dev components** (Charts, Stats, DataTable, Inputs)
* **CentralSet runtime** instead of Evidence’s native data layer
* **ETLX-generated datasets and Parquet files**
* **DuckDB WASM** for in-browser analytics

To start, by deafult an empty dashboard is display and it has an edit button:

[Empty Dashboard](images/screenshots/empty-dash-light.png)

On clicking in the edit button, this form is openned, where you can set your configuration ([more about configuration](config/))

[Empty Dashboard](images/screenshots/edit-dash-light.png)

...

---

## 1. Dashboard Configuration Block

Every dashboard starts with a `config` block.
This is where CentralSet **overrides Evidence defaults** and wires the dashboard to ETLX outputs.

```json config
"auto_refresh_every_n_seconds": null,  
"all_query_run_locally_in_ddb_wasm": true,
"pre_prepared_parquets_logs_table": null,
"pre_prepared_parquets_logs_sql": "with _logs as (select * from etlx_logs where fname is not null and (fname, end_at) in (select fname, max(end_at) from etlx_logs group by fname)) select item_key as name, replace(fname, 'tmp/', '') as file FROM _logs",
"pre_prepared_parquets_logs_db": "sqlite3:database/sqlite_ex.db",
"pre_prepared_parquets_for_ddb_wasm": {
    "LOGS": "hist_logs.parquet"
},
"replace_source_name_in_sql": {
    "ds_name": "ds_x|ds_y"
},
"query_datasource": {
    "ds_name": "ds_"
}
```

### What this configuration does

#### Execution model

* `all_query_run_locally_in_ddb_wasm: true`
  → All SQL runs **inside the browser** using DuckDB WASM.

#### Data source binding

* `pre_prepared_parquets_for_ddb_wasm`
  Maps logical table names (`LOGS`) to **ETLX-exported Parquet files**.

This means:

```sql
FROM "LOGS"
```

actually reads from `hist_logs.parquet`.

#### Pre-prepared metadata

* `pre_prepared_parquets_logs_sql`
  Used to **discover available log files** and metadata before loading the dashboard.

#### Source name normalization

* `replace_source_name_in_sql` and `query_datasource`
  Allow dashboards to stay **portable**, even when datasets are renamed or versioned.

📸 **What you see in the UI**
The dashboard loads instantly with no backend round-trips, powered entirely by Parquet + WASM.

---

## 2. Filters Section (User Controls)

Filters define **how users interact with the data**.
Each filter is backed by SQL queries that populate dropdowns and inputs.

---

### 2.1 Reference Date Selector

```sql _dts
select distinct strftime("ref"::date, '%Y-%m-%d') "ref"
from "LOGS"
where "ref" is not null
order by "ref" desc
```

**Purpose**

* Extracts all available reference dates from logs
* Feeds the date picker suggestions

📸 **UI result**
A date input that **suggests only valid dates present in the data**.

---

### 2.2 Process Hierarchy Filters

#### Main Process

```sql main_process_query
select distinct "key" as "main_process"
from "LOGS"
```

#### Sub Process

```sql sub_process_query
select distinct "item_key" as "sub_process"
from "LOGS"
where "item_key" is not null
```

These queries populate hierarchical dropdowns:

* **Main Process** → high-level ETLX step
* **Sub Process** → specific task within the step

---

### 2.3 Success / Failure Filter

```sql query_success
select *
from (values
    ('%', 'ALL'),
    ('true', 'TRUE'),
    ('false', 'FALSE')
) t ("val", "desc")
```

Allows filtering logs by:

* All
* Successful runs
* Failed runs

---

### 2.4 Filter Layout (Grid)

```html
<Grid>
    <GridItem>
        <Input type=date ... input_label="Reference Date" />
    </GridItem>

    <GridItem>
        <Input type=date ... input_label="Reference Date N-1" />
    </GridItem>

    <GridItem>
        <Dropdown ... input_label='Main Process'/>
    </GridItem>

    <GridItem>
        <Dropdown ... input_label='Sub Process'/>
    </GridItem>

    <GridItem>
        <RadioButtons input_label="Status" />
    </GridItem>

    <GridItem class='grow'/>

    <GridItem>
        <Button icon="pencil" tooltip="Edit"/>
    </GridItem>

    <GridItem>
        <Button icon="refresh" tooltip="Update"/>
    </GridItem>
</Grid>
```

📸 **UI result**

* A compact **control bar**
* Filters on the left
* Action buttons (edit, refresh, duplicate) on the right

---

## 3. Big Numbers (KPIs)

The **Big Numbers section** summarizes system health at a glance.

---

### 3.1 KPI Query

```sql big_numbers_query
select
    count(*) filter("ref" = 'inputs.date_ref.value') as "total",
    count(*) filter("ref" = 'inputs.date_ref.value' and success) as "total_success",
    count(*) filter("ref" = 'inputs.date_ref.value' and not success) as "total_fail",
    total_success / total as "success_delta",
    total_fail / total as "fail_delta",
    count(*) filter("ref" = 'inputs.date_ref_n1.value') as "total_n1",
    ((total - total_n1) / total_n1) as "total_delta"
from "LOGS"
where "key" like 'inputs.main_process.value'
    and "item_key" like 'inputs.sub_process.value'
```

This query:

* Compares **today vs N-1**
* Computes **success/failure ratios**
* Powers KPI deltas

---

### 3.2 KPI Rendering

```html
<Stats>
    <Stat>
        <StatFigure icon='document-text'/>
        <StatTitle># Total Logs Entry</StatTitle>
        <StatValue value=total />
        <StatDesc value=total_delta fmt=pct2 />
    </Stat>

    <Stat>
        <StatFigure icon='check-circle'/>
        <StatTitle># Success</StatTitle>
        <StatValue value=total_success />
    </Stat>

    <Stat>
        <StatFigure icon='x-circle'/>
        <StatTitle># Fail</StatTitle>
        <StatValue value=total_fail />
    </Stat>
</Stats>
```

📸 **UI result**

* Large numeric tiles
* Color-coded status
* Percentage deltas vs previous period

---

## 4. Charts Section

This section explains **why things changed**, not just **what changed**.

---

### 4.1 Logs by Process (Pie Chart)

```sql total_by_process_query
select "key" as "name", count(*) as "value"
from "LOGS"
where "ref" = 'inputs.date_ref.value'
group by "key"
```

```html
<ECharts
    config={{
        series: [{
            type: 'pie',
            radius: ['40%', '70%'],
            data: queries.total_by_process_query.data
        }]
    }}
/>
```

📸 **UI result**

* Donut chart showing **log volume by process**
* Quickly highlights bottlenecks

---

### 4.2 Logs Over Time (Area Chart)

```sql total_by_ref_query
select
    "ref"::varchar as "dt",
    case when success then 'Success' else 'Error' end as "category",
    count(*) as "total"
from "LOGS"
group by "ref", category
order by "ref"
```

```html
<AreaChart
    data={total_by_ref_query}
    x=dt
    y=total
    series=category
/>
```

📸 **UI result**

* Time series of successes vs errors
* Clear trend visibility

---

## 5. Log Details Table

The final section lets users **drill into raw data**.

---

### 5.1 Logs Query

```sql _logs
select *
from "LOGS"
where "ref" = 'inputs.date_ref.value'
order by "start_at"
```

### 5.2 DataTable

```html
<DataTable data={_logs} rows=20 search=true>
    <Column id=key title="Process"/>
    <Column id=item_key title="Sub Process"/>
    <Column id=start_at title="Start"/>
    <Column id=end_at title="End"/>
    <Column id=duration title="Duration"/>
    <Column id=msg title="Message"/>
    <Column id=success title="Success"/>
</DataTable>
```

📸 **UI result**

* Paginated
* Searchable
* Ideal for debugging failed ETLX runs

---

## How This Differs From Native Evidence

| Feature     | Evidence      | CentralSet              |
| ----------- | ------------- | ----------------------- |
| Data Source | SQL files     | ETLX Parquet outputs    |
| Runtime     | Node / DB     | DuckDB WASM             |
| Config      | `sources.yml` | Embedded `config` block |
| Deployment  | Evidence app  | Embedded dashboard      |
| Purpose     | BI            | BI / Operational analytics   |

---

## Summary

This dashboard demonstrates how **ETLX + CentralSet** turns operational logs into:

* Instant analytics
* Zero-backend dashboards
* Fully portable markdown definitions

It combines:

* SQL for logic
* Evidence components for UI
* Parquet for performance
