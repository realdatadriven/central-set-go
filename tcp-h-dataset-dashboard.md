<!---DASHBOARD - exemple using tcp-h generated dataset to simulate a dashboard from a e-comerce-->

```config
    "auto_refresh_every_n_seconds": null,
    "replace_source_name_in_sql": {
        "ds_name": "ds_x|ds_y"
    },
    "query_datasource": {
        "ds_name": "ds_"
    }
```

<!---FILTERS SECTION-->
```sql test
select *
from "sample.duckdb"."lineitem"
limit 10
```

<!-- REGION -->
```sql r_region
select * from (values ('%', 'ALL')) t ("val", "desc")
union
select r_name as "val", r_name as "desc"
from "sample.duckdb"."region"
```

<!-- SHIP MODE -->
```sql l_shipmode
select l_shipmode  as "val", l_shipmode  as "desc"
from "sample.duckdb"."lineitem"
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
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <RadioButtons 
            data={r_region} 
            name=region
            value=val
            label=desc
            defaultValue=nth_0
            _class='btn-sm'
            input_label="Region"
        />
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <Dropdown data={l_shipmode} name=ship_mode value=val label=desc defaultValue="%" input_label='Ship Mode'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
        <!--<RadioButtons 
            data={l_shipmode} 
            name=ship_mode
            value=val
            label=desc
            defaultValue=nth_0
            _class='btn-sm'
            input_label="Ship Mode"
        />-->
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1 grow'/>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Edit" name = "edit" action = "edit" label="" icon = "pencil" _class='btn-sm btn-gost' />
    </GridItem>
    {/if}
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Clone" name = "create" action = "create" label="" icon = "plus" _class='btn-sm btn-gost' />
    </GridItem>
    {/if}
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Clone" name = "duplicate" action = "duplicate" label="" icon = "document-duplicate" _class='btn-sm btn-gost' />
    </GridItem>
    {/if}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Update" name = "refresh" action = "refresh" label="" icon = "refresh" _class='btn-sm btn-gost' />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Details" name = "details" action = "details" label="" icon = "ellipsis-vertical" _class='btn-sm btn-gost' />
    </GridItem>
     {/if}
</Grid>

