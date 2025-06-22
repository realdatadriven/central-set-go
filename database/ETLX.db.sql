BEGIN TRANSACTION;
DROP TABLE IF EXISTS "dashboard_comment";
CREATE TABLE IF NOT EXISTS "dashboard_comment" (
	"dashboard_comment_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"dashboard_comment"	TEXT,
	"dashboard"	VARCHAR(200),
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "dashboard";
CREATE TABLE IF NOT EXISTS "dashboard" (
	"dashboard_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"dashboard"	VARCHAR(200),
	"dashboard_desc"	TEXT,
	"dashboard_conf"	TEXT NOT NULL,
	"order"	INTEGER,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "etlx_conf";
CREATE TABLE IF NOT EXISTS "etlx_conf" (
	"etlx_conf_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etlx_conf"	VARCHAR(200) NOT NULL UNIQUE,
	"etlx_conf_desc"	TEXT,
	"etlx_extra_conf"	TEXT,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "etlx";
CREATE TABLE IF NOT EXISTS "etlx" (
	"etlx_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl"	VARCHAR(200) NOT NULL UNIQUE,
	"etl_desc"	TEXT,
	"attach_etlx_conf"	VARCHAR(200),
	"etlx_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
INSERT INTO "etlx" ("etlx_id","etl","etl_desc","attach_etlx_conf","etlx_conf","active","user_id","app_id","created_at","updated_at","excluded") VALUES (1,'S3_EXTRACT','Example extrating from S3 to a local sqlite3 file',NULL,'# ETL

The [`httpfs`](https://duckdb.org/docs/extensions/httpfs/s3api, "httpfs") extension supports reading/writing/globbing files on object storage servers using the S3 API. S3 offers a standard API to read and write to remote files (while regular http servers, predating S3, do not offer a common write API). DuckDB conforms to the S3 API, that is now common among industry storage providers.
The preferred way to configure and authenticate to S3 endpoints is to use secrets. Multiple secret providers are available

```yaml metadata
name: S3_EXTRACT
description: "Example extrating from S3 to a local sqlite3 file"
connection: "duckdb:"
active: true
```

## VERSION

```yaml metadata
name: VERSION
description: "DDB Version"
table: VERSION
load_conn: "duckdb:"
load_before_sql: "ATTACH ''database/S3_EXTRACT.db'' AS DB (TYPE SQLITE)"
load_sql: ''CREATE OR REPLACE TABLE DB."<table>" AS SELECT version() AS "VERSION";''
load_after_sql: "DETACH DB;"
rows_sql: ''SELECT COUNT(*) AS "nrows" FROM DB."<table>"''
active: true
```

## train_services

```yaml metadata
name: train_services
description: "train_services"
table: train_services
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
INSTALL httpfs;
LOAD httpfs;
```

```sql
-- attach_db
ATTACH ''database/S3_EXTRACT.db'' AS "DB" (TYPE SQLITE)
```

```sql
-- detach_db
DETACH "DB";
```

```sql
-- load_query
CREATE OR REPLACE TABLE "DB"."<table>" AS
FROM ''s3://duckdb-blobs/train_services.parquet'';
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

## S3_EXTRACT

```yaml metadata
name: S3_EXTRACT
description: "Example extrating from S3 to a local sqlite3 file"
table: S3_EXTRACT
load_conn: "duckdb:"
load_before_sql:
  - load_extentions
  - attach_db
  - create_S3_token
load_sql: load_query
load_after_sql: detach_db
drop_sql: drop_sql
clean_sql: clean_sql
rows_sql: nrows
active: false
```

```sql
-- load_extentions
INSTALL httpfs;
LOAD httpfs;
```

```sql
-- attach_db
ATTACH ''database/S3_EXTRACT.db'' AS "DB" (TYPE SQLITE)
```

Example with a [Minio](https://min.io/) local instance

```sql
-- create_S3_token
CREATE SECRET S3_token (
   TYPE S3,
   KEY_ID ''@S3_KEY_ID'',
   SECRET ''@S3_SECRET'',
   ENDPOINT ''127.0.0.1:3000'',
   URL_STYLE ''path''
);
```

```sql
-- detach_db
DETACH "DB";
```

```sql
-- load_query
CREATE OR REPLACE TABLE "DB"."<table>" AS
SELECT * 
FROM ''s3://uploads/flights.csv'';
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
active: true
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
ATTACH ''database/S3_EXTRACT.db'' AS "DB" (TYPE SQLITE)
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
',1,1,3,'2025-03-05T19:47:27.74833913-01:00','2025-03-18 19:03:08.49863581-01:00',0),
 (2,'HTTP_EXTRACT','Example extrating from web to a local sqlite3 file',NULL,'# ETL

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
active: true
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
active: true
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
path: "static/uploads/tmp"
active: false
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
tmp_prefix: ''tmp''
active: true
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
tmp_prefix: ''tmp''
active: true
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
active: true
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
<b>Good Morning</b><br /><br />
This email is gebnerated by ETLX automatically!<br />
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
            <td>{{ .start_at }}</td>
            <td>{{ .end_at }}</td>
            <td>{{ .duration }}</td>
            <td>{{ .success }}</td>
            <td>{{ .msg }}</td>
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
',1,1,3,'2025-03-18T18:51:29.465995658-01:00','2025-03-29 20:46:44.157247059-01:00',0);
COMMIT;
