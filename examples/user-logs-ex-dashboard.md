<!---DASHBOARD - exemple using tcp-h generated dataset to simulate a dashboard from a e-comerce-->
<!-- markdownlint-disable MD041 -->
<!-- markdownlint-disable MD033 -->
<!-- markdownlint-disable MD009 -->

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

<!-- ACTIONS -->
```sql actions
select distinct action as "val", action as "desc"
from "ADMIN.db"."user_log"
```

<!-- SUCCESS -->
```sql status
select column1 as "val", column2 as "desc"
from (values 
    ('%', 'ALL'),
    ('success', 'Success'),
    ('error', 'Error')
) -- t ("val", "desc")
```

<!-- USERS -->
```sql users
select distinct "user_log"."user_id"  as "val", "users"."username"  as "desc"
from "ADMIN.db"."user_log"
join "ADMIN.db"."users" ON "users"."user_id" = "user_log"."user_id"
```

<!-- DATABASES -->
```sql dbs
select distinct "db" as "val", "db" as "desc"
from "ADMIN.db"."user_log"
where "db" is not null
```

<!-- TABLES -->
```sql tables
select distinct "table" as "val", "table" as "desc"
from "ADMIN.db"."user_log"
where "table" is not null
```

<Grid>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(7, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_start
            input_label="Reference Date N-1:"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_end
            input_label="Reference Date:"
        />
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <RadioButtons 
            data={status} 
            name=status
            value=val
            label=desc
            defaultValue=nth_0
            _class='btn-sm'
            input_label="Status:"
        />
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <!--<RadioButtons 
            data={actions} 
            name=action
            value=val
            label=desc
            defaultValue=nth_0
            _class='btn-sm'
            input_label="Actions:"
        />-->
        <Dropdown data={actions} name=action value=val label=desc defaultValue="%" input_label='Actions:'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <Dropdown data={users} name=user value=val label=desc defaultValue="%" input_label='Users:'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <Dropdown data={dbs} name=db value=val label=desc defaultValue="%" input_label='Databases:'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
    </GridItem>
    <GridItem width='w-auto !align-bottom' _type='auto' _class='p-1'>
        <Dropdown data={tables} name=table value=val label=desc defaultValue="%" input_label='Tables:'>
            <DropdownOption value="%" valueLabel="All"/>
        </Dropdown>
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

<!--- DASHBOARD CONTENT SECTION -->

```sql big_nubers
WITH logs AS (
    SELECT l.*
        , CAST((julianday(l.res_at) - julianday(l.req_at)) * 86400 AS INT) AS duration_seconds
    FROM "ADMIN.db"."user_log" l
    LEFT JOIN "ADMIN.db"."users" u ON u."user_id" = l."user_id"
    WHERE l."req_at" BETWEEN 'inputs.date_start.value' AND 'inputs.date_end.value'
        AND COALESCE(l."action", '') LIKE 'inputs.action.value'
        AND COALESCE(l."res_type", '') LIKE 'inputs.status.value'
        AND COALESCE(l."db", '') LIKE 'inputs.db.value'
        AND COALESCE(l."table", '') LIKE 'inputs.table.value'
)
SELECT count(user_log_id) total
    , count(user_log_id) filter(where res_type = 'success') total_success
    , count(user_log_id) filter(where res_type = 'error') total_error
    , avg(duration_seconds) as avg_duration_seconds
FROM logs
```
