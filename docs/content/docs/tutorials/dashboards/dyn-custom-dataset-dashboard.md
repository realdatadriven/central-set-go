---
weight: 7084
title: "Exporting Dashboard Data to Excel Templates"
description: "Generate spreadsheet files from dashboard filters and queries using ETLX export templates."
icon: description
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---
1. Add a buttom update my ds with action `update_custom_ds` and give it a name that will be used in the cofig for metadata and in the markdwon block as its id like this:
```html {linenos=table}
<GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
    <Button 
        tooltip="Update Dashboard Data" 
        name="my_custom_ds_ex" 
        action="update_custom_ds"
        icon="document-arrow-down" 
        _class='btn-sm btn-gost' 
    />
</GridItem>
```
2. In config custom ds points to a update_custom_ds like:
```config
"update_custom_ds": {
  "my_custom_ds_ex": {
    "etlx_id": 1,
    "app": {"app_id": 2,"app": "ETLX","db": "ETLX"},
    "map_files_2_tables": {
        "SALES": "sales_by_dep.parquet"
    }
  }
}
```
3. The config itsel that will be some traformation and also generates the files for the dashboard the will be custom to each user role, this is to be used where dashboard will ned to be user scoped, or any other scope like
`````markdown {linenos=table}
<!-- ETLX CODE BLOCK - EXPORT DATASET -->

````markdown my_custom_ds_ex
# GENERATE_DATA_SETS
```yaml
name: GenerateDSs
description: Exports custom ds
runs_as: EXPORTS
connection: "duckdb:"
path: static/uploads/ # in case of a s3endpoint it can be passed directlly in the export query
active: true
```
## SALES
```yaml
name: GenerateSALESData
description: Exports custom SALES data
connection: "duckdb:"
before_sql:
  - INSTALL erpl_web FROM community
  - LOAD erpl_web
  - create_api_auth_secrete
  - attach_odata_endpoint_with_users_copes
  - attach_sales_datalake
  - attach_logs_db
export_sql: 
  - generate_my_sales_data
  - create_logs_table_if_not_exists
  - insert_generated_file_into_logs
after_sql:
  - DETACH scopes
  - DETACH dl
  - DETACH logs
path: tmp/sales_by_dep.[dash.user.user_id].{YYYYMMDD}.{TSTAMP}.parquet'
tmp_prefix: 'tmp'
active: true
```
```sql
-- create_api_auth_secrete
CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '<JWT_TOKEN>',
  SCOPE 'http://localhost:4444/'
);
```
```sql
-- attach_odata_endpoint_with_users_copes
ATTACH IF NOT EXISTS 'http://localhost:4444/odata/rla_tables' AS scopes (TYPE ODATA);
```
```sql
-- attach_ex_sales_datalake
ATTACH 'ducklake:sqlite:database/dl_metadata.sqlite' AS dl (DATA_PATH 'dl/');
```
```sql
-- attach_logs_db
ATTACH 'database/logs_for_dyn_gen_ds.db' AS logs (TYPE SQLITE);
```
```sql
-- generate_my_sales_data
COPY (
  SELECT sales_by_dep.*
  FROM dl.sales.sales_by_dep
  WHERE sales_by_dep.department_id IN (
    -- OData API Will use CS crud/read by the user in the <JWT_TOKEN> given by its session
    SELECT department_id
    FROM scopes.department
  )
) TO '<fname>';
```
```sql
-- create_logs_table_if_not_exists
CREATE TABLE IF NOT EXISTS logs.dynamic_ds_logs (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    fname      VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
```sql
-- insert_generated_file_into_logs
INSERT INTO logs.dynamic_ds_logs (user_id, fname) VALUES ([dash.user.user_id], '<fname>');
```
````
`````

in this exemple the etlx block genertes a parquet file (coud be many) and return the generated file name that is maped to a table in the global configuration, making that table avaliable for nomal dashboard query, mut that dataset was dynamically genearted by users onlly access, this is the best way to scope data by departments, vendors, tenants when needed to be done
as you ca se the datalake has ulfiltered data nad by itsel it does not apply row level security, but in the CS ETLX app, evry one of my users is scoped departments  one or more, and the query generate_my_sales_data