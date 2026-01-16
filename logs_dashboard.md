<!---DASHBOARD-->

```config
    "auto_refresh_every_n_seconds": null,  
    "all_query_run_locally_in_ddb_wasm": true,
    "pre_prepared_parquets_logs_table": null,
    "pre_prepared_parquets_logs_sql": "with _logs as (select * from etlx_logs where fname is not null and (fname, end_at) in (select fname, max(end_at) from etlx_logs group by fname)) select item_key as name, replace(fname, 'tmp/', '') as file FROM _logs",
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
select 'Option 1' as "label", 1 as "value"
union
select 'Option 2' as "label", 2 as "value"
```

<!-- DATA (add option to sugest dates that alread in the database to the compenent) --> 
```sql _dts
select distinct strftime("ref"::date, '%Y-%m-%d') "ref"
from "LOGS"
where "ref" is not null
order by "ref" desc
```

<!-- LEVEL 1 - PROCESS -->
```sql lvl1_proc
select distinct "key"
from "LOGS"
```

<!-- LEVEL 2 - ITEMS -->

<!-- SUCCESS -->

<Grid>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            extra_cls='input input-sm input-bordered'
            list=dates 
            data={_dts} 
            name=date_ref 
            options_value=ref 
            options_label=ref
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1'>
        <Dropdown data={query_sample} name=input_1 value=value label=label defaultValue="%" --no-input_label='Dropdown Exemple'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 grow'/>
    <GridItem width='w-auto !align-bottom' _type='auto' extra_cls='p-1'>
        <RadioButtons 
            data={query_sample} 
            name=input_2
            value=value
            label=label
            defaultValue=nth_1
            extra_cls='btn-sm'
        />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 text-left'>
        <Button tooltip="Editar" name = "edit" action = "edit" label="" icon = "pencil" extra_cls='btn-sm btn-gost' />
    </GridItem>
    {/if}
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 text-left'>
        <Button tooltip="Duplicar" name = "duplicate" action = "duplicate" label="" icon = "document-duplicate" extra_cls='btn-sm btn-gost' />
    </GridItem>
    {/if}
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 text-left'>
        <Button tooltip="Atualizar" name = "refresh" action = "refresh" label="" icon = "refresh" extra_cls='btn-sm btn-gost' />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 text-left'>
        <Button tooltip="Detalhes" name = "details" action = "details" label="" icon = "ellipsis-vertical" extra_cls='btn-sm btn-gost' />
    </GridItem>
     {/if}
</Grid>

<!---DASHBOARD SECTION-->

```sql _logs
select * 
from "LOGS"
```

{inputs?.date_ref?.value}