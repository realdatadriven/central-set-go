<!---DASHBOARD-->

```config
    "auto_refresh_every_n_seconds": null,  
    "all_query_run_locally_in_ddb_wasm": true,
    "pre_prepared_parquets_logs_table": null,
    "pre_prepared_parquets_logs_sql": "with _logs as (select *  from dynamic_ds_logs  where fname is not null and user_id = [dash.user.user_id] and (user_id, table_name, created_at) in ( select user_id, table_name, max(created_at) from dynamic_ds_logs group by user_id, table_name ) ) select user_id, table_name as name, replace(fname, 'tmp/', '') as file, created_at as lst_dt from _logs",
    "pre_prepared_parquets_logs_db": "sqlite3:database/logs_for_dyn_gen_ds.db",
    "pre_prepared_parquets_for_ddb_wasm": {
        "ORDERS": "orders.parquet"
    },
    "replace_source_name_in_sql": {
        "ds_name": "ds_x|ds_y"
    },
    "query_datasource": {
        "ds_name": "ds_"
    },
    "update_custom_ds": {
        "my_custom_ds_ex": {
            "etlx_id": 1,
            "app": { "app_id": 2, "app": "ETLX", "db": "ETLX"}
        }
    }
```

<!---FILTERS SECTION-->

<!-- SUCCESS -->
```sql query_success
select *
from (values
    ('%', 'ALL')
    , ('true', 'TRUE')
    , ('false', 'FALSE')
) t ("val", "desc")
```
<Grid>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            list=dates 
            data={_dts} 
            name=date_ref 
            options_value=ref 
            options_label=ref
            input_label="Reference Date"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(2, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_ref_n1
            input_label="Reference Date N-1"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1 grow'/>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Edit" name = "edit" action = "edit" label="" icon = "pencil" _class='btn-sm btn-gost' />
    </GridItem>
    {/if}
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Clone" name = "duplicate" action = "duplicate" label="" icon = "document-duplicate" _class='btn-sm btn-gost' />
    </GridItem>
    {/if}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Update Dashboard Data" name="my_custom_ds_ex" action="update_custom_ds" icon="cloud-arrow-down" _icon="arrow-path-rounded-square" _class='btn-sm btn-gost' 
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Print" name = "print" action = "print" label="" icon = "printer" _class='btn-sm btn-gost' />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Update" name = "refresh" action = "refresh" label="" icon = "refresh" _class='btn-sm btn-gost' />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Details" name = "details" action = "details" label="" icon = "ellipsis-vertical" _class='btn-sm btn-gost' />
    </GridItem>
     {/if}
</Grid>

<!---DASHBOARD SECTION-->

```sql total_orders
select count(*) as total
from "ORDERS"
```

<!-- ETLX CODE BLOCK - EXPORT DATASET -->

````markdown my_custom_ds_ex
# GENERATE_DATA_SETS
```yaml
name: GenerateDSs
description: Exports custom ds
runs_as: EXPORTS
connection: "duckdb:"
path: static/uploads # in case of a s3endpoint it can be passed directlly in the export query
active: true
```
## ORDERS
```yaml
name: GenerateORDERSData
description: Exports custom ORDERS data
connection: "duckdb:"
before_sql:
  - INSTALL erpl_web FROM community
  - LOAD erpl_web
  - create_api_auth_secrete
  - attach_odata_endpoint_with_users_copes
  - attach_ex_ecomerce_datalake
  - attach_logs_db
export_sql: 
  - generate_my_orders_data
  - create_logs_table_if_not_exists
  - insert_generated_file_into_logs
after_sql:
  - DETACH scopes
  - DETACH dl
  - DETACH logs
path: tmp/orders.[dash.user.user_id].{YYYYMMDD}.{TSTAMP}.parquet
tmp_prefix: tmp
active: true
```
```sql
-- create_api_auth_secrete
CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '[dash.user.jwt_token]',
  SCOPE 'http://localhost:4444/'
);
```
```sql
-- attach_odata_endpoint_with_users_copes
ATTACH IF NOT EXISTS 'http://localhost:4444/odata/ETLX' AS scopes (TYPE ODATA);
```
```sql
-- attach_ex_ecomerce_datalake
ATTACH 'ducklake:sqlite:database/dl_metadata.sqlite' AS dl (DATA_PATH 'database/dl/');
```
```sql
-- attach_logs_db
ATTACH 'database/logs_for_dyn_gen_ds.db' AS logs (TYPE SQLITE);
```
```sql
-- generate_my_orders_data
COPY (
  SELECT orders.*
  FROM dl.orders
  /*WHERE orders.department_id IN (
    -- OData API Will use CS crud/read by the user in the <JWT_TOKEN> given by its session
    SELECT department_id
    FROM scopes.department
  )*/
) TO '<fname>';
```
```sql
-- create_logs_table_if_not_exists
CREATE TABLE IF NOT EXISTS logs.dynamic_ds_logs (
    --id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    table_name VARCHAR NOT NULL,
    fname      VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
```sql
-- insert_generated_file_into_logs
INSERT INTO logs.dynamic_ds_logs (user_id, table_name, fname) VALUES ([dash.user.user_id], 'ORDERS', PARSE_FILENAME('<fname>'));
```
```sql  x
with _logs as (
  select * 
  from dynamic_ds_logs 
  where fname is not null
    and user_id = [dash.user.user_id]
    and (user_id, table_name, created_at) in (
      select user_id, table_name, max(created_at) 
      from dynamic_ds_logs 
      group by user_id, table_name
    )
) 
select user_id, table_name as name, replace(fname, 'tmp/', '') as file, created_at as lst_dt
FROM _logs
```
````