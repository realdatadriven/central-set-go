BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS arrow_flight (
	arrow_flight_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	arrow_flight VARCHAR(200) NOT NULL, 
	arrow_flight_desc TEXT, 
	arrow_flight_conf TEXT, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (arrow_flight)
);
CREATE TABLE IF NOT EXISTS dashboard (
	dashboard_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	dashboard VARCHAR(200), 
	dashboard_desc TEXT, 
	dashboard_conf TEXT NOT NULL, 
	"order" INTEGER, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS dashboard_comment (
	dashboard_comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	dashboard_comment TEXT, 
	dashboard VARCHAR(200), 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS etlx (
	etlx_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	etl VARCHAR(200) NOT NULL, 
	etl_desc TEXT, 
	attach_etlx_conf VARCHAR(200), 
	etlx_conf TEXT, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (etl)
);
CREATE TABLE IF NOT EXISTS etlx_conf (
	etlx_conf_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	etlx_conf VARCHAR(200) NOT NULL, 
	etlx_conf_desc TEXT, 
	etlx_extra_conf TEXT, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (etlx_conf)
);
CREATE TABLE IF NOT EXISTS manage_query (
	manage_query_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	manage_query VARCHAR(200) NOT NULL, 
	"database" VARCHAR(200) NOT NULL, 
	manage_query_conf TEXT, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS notebook (
	notebook_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	notebook VARCHAR(200), 
	notebook_desc TEXT, 
	notebook_conf TEXT NOT NULL, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
INSERT INTO "dashboard" VALUES (1,'Logs ETL','Logs','<!---DASHBOARD-->

```config
    "auto_refresh_every_n_seconds": null,  
    "all_query_run_locally_in_ddb_wasm": true,
    "pre_prepared_parquets_logs_table": null,
    "pre_prepared_parquets_logs_sql": "with _logs as (select * from etlx_logs where fname is not null and (fname, end_at) in (select fname, max(end_at) from etlx_logs group by fname)) select item_key as name, replace(fname, ''tmp/'', '''') as file FROM _logs",
    "pre_prepared_parquets_logs_db": "postgres:user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable",
    "pre_prepared_parquets_for_ddb_wasm": {
        "LOGS": "hist_logs.parquet"
    },
    "replace_source_name_in_sql": {
        "ds_name": "ds_x|ds_y"
    },
    "query_datasource": {
        "ds_name": "ds_"
    },
    "disable_md": false,
    "disable_code": false
```

<!---FILTERS SECTION-->

```sql query_sample
select ''Option 1'' as "label", 1 as "value"
union
select ''Option 2'' as "label", 2 as "value"
```

<Grid>
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1''>
        <Dropdown data={query_sample} name=input_1 value=value label=label defaultValue="%" --no-input_label=''Dropdown Exemple''>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1 grow''/>
    <GridItem width=''w-auto !align-bottom'' _type=''auto'' extra_cls=''p-1''>
        <RadioButtons 
            data={query_sample} 
            name=input_2
            value=value
            label=label
            defaultValue=nth_1
            extra_cls=''btn-sm''
        />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1 text-left''>
        <Button tooltip="Editar" name = "edit" action = "edit" label="" icon = "pencil" extra_cls=''btn-sm btn-gost'' />
    </GridItem>
    {/if}
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1 text-left''>
        <Button tooltip="Duplicar" name = "duplicate" action = "duplicate" label="" icon = "document-duplicate" extra_cls=''btn-sm btn-gost'' />
    </GridItem>
    {/if}
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1 text-left''>
        <Button tooltip="Atualizar" name = "refresh" action = "refresh" label="" icon = "refresh" extra_cls=''btn-sm btn-gost'' />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width=''w-auto'' _type=''auto'' extra_cls=''p-1 text-left''>
        <Button tooltip="Detalhes" name = "details" action = "details" label="" icon = "ellipsis-vertical" extra_cls=''btn-sm btn-gost'' />
    </GridItem>
     {/if}
</Grid>

<!---DASHBOARD SECTION-->

```sql _logs
select * 
from "LOGS"
```
',1,1,1,2,'2025-11-19T19:03:10.163309779-01:00','2025-11-19 20:03:43.485166468-01:00',0);
INSERT INTO "etlx" VALUES (1,'HTTP_EX','HTTP Exemples',NULL,'# ETL

<https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page>

```yaml metadata
name: HTTP_EXTRACT
description: "Example extrating from web to a local sqlite3 file"
connection: "duckdb:"
database: HTTP_EXTRACT.db
active: true
```

## VERSION

```yaml metadata
name: VERSION
description: "DDB Version"
table: VERSION
load_conn: "duckdb:"
load_before_sql: "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
load_sql: ''CREATE OR REPLACE TABLE DB."<table>" AS SELECT version() AS "VERSION";''
load_after_sql: "DETACH DB;"
rows_sql: ''SELECT COUNT(*) AS "nrows" FROM DB."<table>"''
active: true
```

## NYC_TAXI

```yaml metadata
name: NYC_TAXI
description: "Example extrating from web to a local sqlite3 file"
table: NYC_TAXI
load_conn: "duckdb:"
load_before_sql:
  - load_extentions
  - attach_db
load_sql: load_query
load_after_sql: detach_db
drop_sql: drop_sql
clean_sql: clean_sql
rows_sql: nrows
active: false
```

```sql
-- load_extentions
INSTALL sqlite;
LOAD sqlite;
```

```sql
-- attach_db
ATTACH ''database/HTTP_EXTRACT.db'' AS "DB" (TYPE SQLITE)
```

```sql
-- detach_db
DETACH "DB";
```

```sql
-- load_query
CREATE OR REPLACE TABLE "DB"."<table>" AS
SELECT * 
FROM ''https://d37ci6vzurychx.cloudfront.net/trip-data/yellow_tripdata_2024-01.parquet'';
```

```sql
-- drop_sql
DROP TABLE IF EXISTS "DB"."<table>";
```

```sql
-- clean_sql
DELETE FROM "DB"."<table>";
```

```sql
-- nrows
SELECT COUNT(*) AS "nrows" FROM "DB"."<table>"
```

# DATA_QUALITY

```yaml
description: "Runs some queries to check quality / validate."
active: false
```

## Rule0001

```yaml
name: Rule0001
description: "Check if the field trip_distance from the NYC_TAXI is missing or zero"
connection: "duckdb:"
before_sql:
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS \"DB\" (TYPE SQLITE)"
query: quality_check_query
fix_quality_err: fix_quality_err_query
column: total_reg_with_err # Defaults to ''total''.
check_only: true
fix_only: false 
after_sql: "DETACH DB"
active: true
```

```sql
-- quality_check_query
SELECT COUNT(*) AS "total_reg_with_err"
FROM "DB"."NYC_TAXI"
WHERE "trip_distance" IS NULL
  OR "trip_distance" = 0;
```

```sql
-- fix_quality_err_query
UPDATE "DB"."NYC_TAXI"
  SET "trip_distance" = "trip_distance"
WHERE "trip_distance" IS NULL
  OR "trip_distance" = 0;
```

# MULTI_QUERIES

```yaml
description: "Define multiple structured queries combined with UNION."
connection: "duckdb:"
before_sql:
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS \"DB\" (TYPE SQLITE)"
save_sql: save_mult_query_res
save_on_err_patt: ''(?i)table.+with.+name.+(\w+).+does.+not.+exist''
save_on_err_sql: create_mult_query_res
after_sql: "DETACH DB"
union_key: "UNION ALL\n" # Defaults to UNION.
active: false
```

```sql
-- save_mult_query_res
INSERT INTO "DB"."MULTI_QUERY" BY NAME
[[final_query]]
```

```sql
-- create_mult_query_res
CREATE OR REPLACE TABLE "DB"."MULTI_QUERY" AS
[[final_query]]
```

## Row1

```yaml
name: Row1
description: "Row 1"
query: row_query
active: true
```

```sql
-- row_query
SELECT ''# number of rows'' AS "variable", COUNT(*) AS "value"
FROM "DB"."NYC_TAXI"
```

## Row2

```yaml
name: Row2
description: "Row 2"
query: row_query
active: true
```

```sql
-- row_query
SELECT ''total revenue'' AS "variable", SUM("total_amount") AS "value"
FROM "DB"."NYC_TAXI"
```

## Row3

```yaml
name: Row3
description: "Row 3"
query: row_query
active: true
```

```sql
-- row_query
SELECT *
FROM (
  SELECT "DOLocationID" AS "variable", SUM("total_amount") AS "value"
  FROM "DB"."NYC_TAXI"
  GROUP BY "DOLocationID"
  ORDER BY "DOLocationID"
) AS "T"
```

# EXPORTS

Exports data to files.

```yaml metadata
name: DailyReports
description: "Daily reports"
connection: "duckdb:"
path: "examples"
active: true
```

## CSV_EXPORT

```yaml metadata
name: CSV_EXPORT
description: "Export data to CSV"
connection: "duckdb:"
before_sql:
  - "INSTALL sqlite"
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
export_sql: export
after_sql: "DETACH DB"
path: ''nyc_taxy_YYYYMMDD.csv''
tmp_prefix: null
active: false
```

```sql
-- export
COPY (
    SELECT *
    FROM "DB"."NYC_TAXI"
    WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
    LIMIT 100
) TO ''<fname>'' (FORMAT ''csv'', HEADER TRUE);
```

## XLSX_EXPORT

```yaml metadata
name: XLSX_EXPORT
description: "Export data to Excel file"
connection: "duckdb:"
before_sql:
  - "INSTALL sqlite"
  - "LOAD sqlite"
  - "INSTALL excel"
  - "LOAD excel"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
export_sql: xl_export
after_sql: "DETACH DB"
path: ''nyc_taxy_YYYYMMDD.xlsx''
tmp_prefix: null
active: false
```

```sql
-- xl_export
COPY (
    SELECT *
    FROM "DB"."NYC_TAXI"
    WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
    LIMIT 100
) TO ''<fname>'' (FORMAT XLSX, HEADER TRUE, SHEET ''NYC'');
```

## XLSX_TMPL

```yaml metadata
name: XLSX_TMPL
description: "Export data to Excel template"
connection: "duckdb:"
before_sql:
  - "INSTALL sqlite"
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
after_sql: "DETACH DB"
tmp_prefix: null
template: "nyc_taxy_YYYYMMDD.xlsx"
path: "nyc_taxy_YYYYMMDD.xlsx"
mapping:
  - sheet: resume
    range: A2
    sql: resume
    type: value
    key: total
  - sheet: detail
    range: A1
    sql: detail
    type: range
    header: true
active: false
```

```sql
-- resume
SELECT COUNT(*) AS "total"
FROM "DB"."NYC_TAXI"
WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
```

```sql
-- detail_old
SELECT *
FROM "DB"."NYC_TAXI"
WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
LIMIT 100
```

```sql
-- detail
pivot (select * from "DB"."NYC_TAXI") as t
on strftime("tpep_pickup_datetime"::datetime, ''%d'')
using sum(total_amount) AS total, count(*) AS total_trips
group by PULocationID
```

```sql
-- data_to_export
SELECT *
FROM "DB"."NYC_TAXI"
WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
LIMIT 100
```

## TEXT_TMPL

```yaml metadata
name: TEXT_TMPL
description: "Export data to text base template"
connection: "duckdb:"
before_sql:
  - "INSTALL sqlite"
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
data_sql:
  - logs
#  - data
after_sql: "DETACH DB"
tmp_prefix: null
text_template: true
template: template
return_content: false #if true the template text content is returned, useful for integration
path: "nyc_taxy_YYYYMMDD.html"
active: true
```

```sql
-- data
SELECT *
FROM "DB"."NYC_TAXI"
WHERE "tpep_pickup_datetime"::DATETIME <= ''{YYYY-MM-DD}''
LIMIT 100
```

```sql
-- logs
SELECT *
FROM "DB"."etlx_logs"
/*WHERE "ref" = ''{YYYY-MM-DD}''*/
ORDER BY start_at DESC
```

```html template
<style>
  table {
    border-collapse: collapse;
    width: 100%;
    font-family: sans-serif;
    font-size: 14px;
  }
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }
  th {
    background-color: #f2f2f2;
    font-weight: bold;
  }
  tr:nth-child(even) {
    background-color: #f9f9f9;
  }
  tr:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    background-color: #eef6ff;
  }
</style>
<b>ETLX Text Template</b><br /><br />
This is gebnerated by ETLX automatically!<br />
{{ with .logs }}
    {{ if eq .success true }}
      <table>
        <tr>
            <th>Name</th>
            <th>Ref</th>
            <th>Start</th>
            <th>End</th>
            <th>Duration</th>
            <th>Success</th>
            <th>Message</th>
        </tr>
        {{ range .data }}
        <tr>
            <td>{{ .name }}</td>
            <td>{{ .ref }}</td>
            <td>{{ .start_at | date "2006-01-02 15:04:05" }}</td>
            <td>{{ .end_at | date "2006-01-02 15:04:05" }}</td>
            <td style="text-align: right">{{ divf .duration 1000000000 | printf "%.4fs" }}</td>
            <td>{{ .success }}</td>
            <td><span title="{{ .msg }}">{{ .msg | toString | abbrev 30}}</span></td>
        </tr>
        {{ else }}
        <tr>
          <td colspan="7">No items available</td>
        </tr>
        {{ end }}
      </table>
    {{ else }}
      <p>{{.msg}}</p>
    {{ end }}
{{ else }}
<p>Logs information missing.</p>
{{ end }}
```

# LOGS

```yaml metadata
name: LOGS
description: "Example saving logs"
table: etlx_logs
connection: "duckdb:"
before_sql:
  - load_extentions
  - attach_db
  - ''USE DB;''
save_log_sql: load_logs
save_on_err_patt: ''(?i)table.+does.+not.+exist|does.+not.+have.+column.+with.+name''
save_on_err_sql:
  - create_logs
  - get_dyn_queries[create_columns_missing]
  - load_logs
after_sql:
  - ''USE memory;''
  - detach_db
tmp_dir: /tmp
active: false
```

```sql
-- load_extentions
INSTALL Sqlite;
LOAD Sqlite;
INSTALL json;
LOAD json;
```

```sql
-- attach_db
ATTACH ''database/HTTP_EXTRACT.db'' AS "DB" (TYPE SQLITE)
```

```sql
-- detach_db
DETACH "DB";
```

```sql
-- load_logs
INSERT INTO "DB"."<table>" BY NAME
SELECT * 
FROM read_json(''<fname>'');
```

```sql
-- create_logs
CREATE TABLE IF NOT EXISTS "DB"."<table>" AS
SELECT * 
FROM read_json(''<fname>'');
```

```sql
-- create_columns_missing
WITH source_columns AS (
    SELECT column_name, column_type 
    FROM (DESCRIBE SELECT * FROM read_json(''<fname>''))
),
destination_columns AS (
    SELECT column_name, data_type as column_type
    FROM duckdb_columns 
    WHERE table_name = ''<table>''
),
missing_columns AS (
    SELECT s.column_name, s.column_type
    FROM source_columns s
    LEFT JOIN destination_columns d ON s.column_name = d.column_name
    WHERE d.column_name IS NULL
)
SELECT ''ALTER TABLE "DB"."<table>" ADD COLUMN "'' || column_name || ''" '' || column_type || '';'' AS query
FROM missing_columns;
```

# NOTIFY

```yaml metadata
name: Notefication
description: "Notefication"
connection: "duckdb:"
path: "examples"
active: false
```

## ETL_STATUS

```yaml metadata
name: ETL_STATUS
description: "ETL Satus"
connection: "duckdb:"
before_sql:
  - "INSTALL sqlite"
  - "LOAD sqlite"
  - "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
data_sql:
  - logs
after_sql: "DETACH DB"
to:
  - real.datadriven@gmail.com
cc: null
bcc: null
subject: ''ETLX YYYYMMDD''
body: body_tml
attachments_:
  - hf.md
  - http.md
active: true
```

```html body_tml
<style>
  table {
    border-collapse: collapse;
    width: 100%;
    font-family: sans-serif;
    font-size: 14px;
  }
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }
  th {
    background-color: #f2f2f2;
    font-weight: bold;
  }
  tr:nth-child(even) {
    background-color: #f9f9f9;
  }
  tr:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    background-color: #eef6ff;
  }
</style>
<b>Good Morning!</b><br /><br />
This email was gebnerated by ETLX automatically!<br />
LOGS:<br />
{{ with .logs }}
    {{ if eq .success true }}
      <table>
        <tr>
            <th>Name</th>
            <th>Ref</th>
            <th>Start</th>
            <th>End</th>
            <th>Duration</th>
            <th>Success</th>
            <th>Message</th>
        </tr>
        {{ range .data }}
        <tr>
            <td>{{ .name }}</td>
            <td>{{ .ref }}</td>
            <td>{{ .start_at | date "2006-01-02 15:04:05" }}</td>
            <td>{{ .end_at | date "2006-01-02 15:04:05" }}</td>
            <td>{{ divf .duration 1000000000 | printf "%.4fs" }}</td>
            <td>{{ .success }}</td>
            <td><span title="{{ .msg }}">{{ .msg | toString | abbrev 30}}</span></td>
        </tr>
        {{ else }}
        <tr>
          <td colspan="7">No items available</td>
        </tr>
        {{ end }}
      </table>
    {{ else }}
      <p>{{.msg}}</p>
    {{ end }}
{{ else }}
<p>Logs information missing.</p>
{{ end }}
```

```sql
-- logs
SELECT *
FROM "DB"."etlx_logs"
WHERE "ref" = ''{YYYY-MM-DD}''
```

# ETL_LAST

<https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page>

```yaml metadata
name: HTTP_EXTRACT
description: "Example extrating from web to a local sqlite3 file"
connection: "duckdb:"
database: HTTP_EXTRACT.db
runs_as: ETL
active: true
```

## VERSION2

```yaml metadata
name: VERSION2
description: "DDB Version"
table: VERSION2
load_conn: "duckdb:"
load_before_sql: "ATTACH ''database/HTTP_EXTRACT.db'' AS DB (TYPE SQLITE)"
load_sql: ''CREATE OR REPLACE TABLE DB."<table>" AS SELECT version() AS "VERSION";''
load_after_sql: "DETACH DB;"
rows_sql: ''SELECT COUNT(*) AS "nrows" FROM DB."<table>"''
active: true
```
',1,1,2,'2025-11-05T19:12:42.281477582-01:00','2025-11-05 19:12:55.543131515-01:00',0);
INSERT INTO "etlx" VALUES (2,'PG_EX','PG',NULL,'<!-- markdownlint-disable MD025 -->

# ETL

<https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page>

```yaml metadata
name: HTTP_EXTRACT
description: "Example extrating from web to a local postgres file"
connection: "postgres:user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable"
database: "postgres:user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable"
active: true
```

## VERSION

```yaml metadata
name: VERSION
description: "DDB Version"
table: VERSION
load_conn: "duckdb:"
load_before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
load_sql: ''CREATE OR REPLACE TABLE DB."<table>" AS SELECT version() AS "VERSION";''
load_after_sql: "DETACH DB;"
rows_sql: ''SELECT COUNT(*) AS "nrows" FROM DB."<table>"''
active: true
```

## NYC_TAXI

```yaml metadata
name: NYC_TAXI
description: "Example extrating from web to a local postgres file"
table: NYC_TAXI
load_conn: "duckdb:"
load_before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
load_sql: load_query
load_after_sql: DETACH "DB"
drop_sql: DROP TABLE IF EXISTS "DB"."<table>"
clean_sql: DELETE FROM "DB"."<table>"
rows_sql: SELECT COUNT(*) AS "nrows" FROM "DB"."<table>"
active: false
```

```sql
-- load_query
CREATE OR REPLACE TABLE "DB"."<table>" AS
SELECT * 
FROM ''https://d37ci6vzurychx.cloudfront.net/trip-data/yellow_tripdata_2024-01.parquet'';
```

## PeadkHours

```yaml metadata
name: PeadkHours
description: Peask Hours Analysis
table: PeadkHours
transform_conn: "duckdb:"
transform_before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
transform_sql: preadk_hours_load_query
transform_after_sql: DETACH "DB"
drop_sql: DROP TABLE IF EXISTS "DB"."<table>"
clean_sql: DELETE FROM "DB"."<table>"
rows_sql: SELECT COUNT(*) AS "nrows" FROM "DB"."<table>"
active: true
```

```sql
-- preadk_hours_load_query
CREATE OR REPLACE TABLE "DB"."<table>" AS
[[PeakHoursAnalysis]]
```

## DailyRevenueTripVolume

```yaml metadata
name: DailyRevenueTripVolume
description: Daily Revenue and Trip Volume
has_placeholders: true
schema: TRF
database: "postgres:user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable search_path=TRF"
table: DailyRevenueTripVolume
transform_conn: "duckdb:"
transform_before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
transform_sql:
  - CREATE SCHEMA IF NOT EXISTS "DB"."<schema>"
  - DailyRevenueTripVolume
transform_after_sql: DETACH "DB"
drop_sql: DROP TABLE IF EXISTS "DB"."<schema>"."<table>"
clean_sql: DELETE FROM "DB"."<schema>"."<table>"
rows_sql: SELECT COUNT(*) AS "nrows" FROM "DB"."<schema>"."<table>"
active: true
```

```sql
-- DailyRevenueTripVolume
CREATE OR REPLACE TABLE "DB"."TRF"."<table>" AS
SELECT CAST(tpep_pickup_datetime AS DATE) AS trip_date,
    COUNT(*) AS total_trips,
    ROUND(SUM(total_amount), 2) AS total_revenue,
    ROUND(AVG(total_amount), 2) AS avg_revenue_per_trip,
    ROUND(SUM(trip_distance), 2) AS total_miles,
    ROUND(AVG(trip_distance), 2) AS avg_trip_distance
FROM DB.NYC_TAXI
GROUP BY trip_date
ORDER BY trip_date
```

# EXPORTS

Exports data to files.

```yaml metadata
name: EXPORTS
description: Exports Examples
connection: "duckdb:"
path: "static/uploads/tmp"
active: true
```

## ExportDailyRevenueTripVolume

```yaml metadata
name: ExportDailyRevenueTripVolume
description: "Export data to CSV"
connection: "duckdb:"
before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
export_sql: export
after_sql: "DETACH DB"
path: ''ExportDailyRevenueTripVolume_{YYYYMMDD}.{TSTAMP}.parquet''
tmp_prefix: ''tmp''
active: true
```

```sql
-- export
COPY (
    SELECT CAST(tpep_pickup_datetime AS DATE) AS trip_date,
        COUNT(*) AS total_trips,
        ROUND(SUM(total_amount), 2) AS total_revenue,
        ROUND(AVG(total_amount), 2) AS avg_revenue_per_trip,
        ROUND(SUM(trip_distance), 2) AS total_miles,
        ROUND(AVG(trip_distance), 2) AS avg_trip_distance
    FROM DB.NYC_TAXI
    GROUP BY trip_date
    ORDER BY trip_date
) TO ''<fname>''
```

## hist_logs

```yaml metadata
name: hist_logs
description: "Export data to CSV"
connection: "duckdb:"
before_sql: "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
export_sql: export
after_sql: "DETACH DB"
path: ''hist_logs_{YYYYMMDD}.{TSTAMP}.parquet''
tmp_prefix: ''tmp''
active: true
```

```sql
-- export
COPY (
    SELECT *
    FROM DB.etlx_logs
) TO ''<fname>''
```

# LOGS

```yaml metadata
name: LOGS
description: "Example saving logs"
table: etlx_logs
connection: "duckdb:"
before_sql:
  - "ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES)"
  - ''USE DB;''
  - LOAD json
  - "get_dyn_queries[create_columns_missing](ATTACH ''user=postgres password=1234 dbname=ETLX_DATA host=localhost port=5432 sslmode=disable'' AS DB (TYPE POSTGRES), DETACH DB)"
save_log_sql: load_logs
save_on_err_patt: ''(?i)table.+does.+not.+exist''
save_on_err_sql: create_logs
after_sql:
  - ''USE memory;''
  - DETACH "DB"
active: true
```

```sql
-- load_logs
INSERT INTO "DB"."<table>" BY NAME
SELECT * 
FROM read_json(''<fname>'');
```

```sql
-- create_logs
CREATE TABLE IF NOT EXISTS "DB"."<table>" AS
SELECT * 
FROM read_json(''<fname>'');
```

```sql
-- create_columns_missing
WITH source_columns AS (
    SELECT column_name, column_type 
    FROM (DESCRIBE SELECT * FROM read_json(''<fname>''))
),
destination_columns AS (
    SELECT column_name, data_type as column_type
    FROM duckdb_columns 
    WHERE table_name = ''<table>''
),
missing_columns AS (
    SELECT s.column_name, s.column_type
    FROM source_columns s
    LEFT JOIN destination_columns d ON s.column_name = d.column_name
    WHERE d.column_name IS NULL
)
SELECT ''ALTER TABLE "DB"."<table>" ADD COLUMN "'' || column_name || ''" '' || column_type || '';'' AS query
FROM missing_columns
WHERE (SELECT COUNT(*) FROM destination_columns) > 0;
```

# REQUIRES

```yaml metadata
name: REQUIRES
description: "Example requires"
active: true
```

## PeakHoursAnalysis

```yaml metadata
name: PeakHoursAnalysis
description: "Analyze peak hours for NYC Yellow Taxi rides"
path: examples/PeakHoursAnalysis.sql          
```
',1,1,2,'2025-11-16T13:38:33.933754728-01:00','2026-01-11 08:16:13.748448333-01:00',0);
COMMIT;
