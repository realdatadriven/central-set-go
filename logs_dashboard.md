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
```sql main_process_query
select distinct "key" as "main_process"
from "LOGS"
```

<!-- LEVEL 2 - ITEMS -->
```sql sub_process_query
select distinct "item_key" as "sub_process"
from "LOGS"
where "item_key" is not null
```

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
    <GridItem width='w-auto' _type='auto' extra_cls='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            extra_cls='input input-sm input-bordered'
            list=dates 
            data={_dts} 
            name=date_ref 
            options_value=ref 
            options_label=ref
            input_label="Reference Date"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1'>
        <Dropdown data={main_process_query} name=main_process value=main_process label=main_process defaultValue="%" input_label='Main Process'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1'>
        <Dropdown data={sub_process_query} name=sub_process value=sub_process label=sub_process defaultValue="%" input_label='Sub Proccess'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' extra_cls='p-1'>
        <RadioButtons 
            data={query_success} 
            name=success
            value=val
            label=desc
            defaultValue=nth_0
            extra_cls='btn-sm'
            input_label="Status"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' extra_cls='p-1 grow'/>
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

<!--{inputs?.date_ref?.value}-->

<!---DASHBOARD SECTION-->

<!-- BIG NUMBERS SECTION -->



<!-- LOG DETAILS -->

```sql _logs
select *
from "LOGS"
where "ref" = 'inputs.date_ref.value'
    and "key" like 'inputs.main_process.value'
    and "item_key" like 'inputs.sub_process.value'
    and case 
        when "success" is true or "success" = 1 then 'true' 
        else 'false' 
    end like 'inputs.success.value'
order by "start_at" asc
```

<DataTable data={_logs}> 
    <Column id=key title="Proccess"/> 
	<Column id=item_key  title="Sub Proccess"/> 
	<Column id=start_at title=Start/> 
	<Column id=end_at title=End/> 
	<Column id=duration title=Duration/> 
	<Column id=msg title=Message/> 
	<Column id=success title=Success/>
</DataTable>