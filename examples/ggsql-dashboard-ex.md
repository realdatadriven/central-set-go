<!---DASHBOARD - exemple using ggsql to generate vegga-lite spec from ddb and render -->

```config
    "auto_refresh_every_n_seconds": null,
    "replace_source_name_in_sql": {
        "ds_name": "ds_x|ds_y"
    },
    "query_datasource": {
        "ds_name": "ds_"
    },
    "duckdb": {
        "startup": "INSTALL ggsql FROM community;LOAD ggsql;SET ggsql_output = 'spec';",
        "shutdown": null
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
select distinct r_name as "val", r_name as "desc"
from "sample.duckdb"."region"
```

<!-- SHIP MODE -->
```sql l_shipmode
select * from (values ('%', 'ALL')) t ("val", "desc")
union
select distinct l_shipmode  as "val", l_shipmode  as "desc"
from "sample.duckdb"."lineitem"
limit 10
```

<Grid>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_start
            input_label="Reference Date"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(2, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_end
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
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Export" name = "fill-template" action = "fill-template" label="" icon = "document-arrow-down" _class='btn-sm btn-gost' />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Print" name = "print" action = "print" label="" icon = "printer" _class='btn-sm btn-gost' />
    </GridItem>
    {#if $_global?.tables?.[$$props?.table]?.permissions?.create === true || !$_global?.tables?.[$$props?.table]?.permissions}
    <GridItem width='w-auto' _type='auto' _class='p-1 text-left'>
        <Button tooltip="Details" name = "details" action = "details" label="" icon = "ellipsis-vertical" _class='btn-sm btn-gost' />
    </GridItem>
     {/if}
</Grid>

<!--- DASHBOARD CONTENT SECTION --->

```sql get_the_raw_spec
SELECT x, x*x AS y
FROM range(10) t(x) 
VISUALISE x, y DRAW line;
```

<VegaEmbed data={get_the_raw_spec} 
    column=plot 
    _class=""
    _style=""
/>