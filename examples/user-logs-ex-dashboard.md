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
            defaultValue={config?.moment()?.subtract(1, 'day').format('YYYY-MM-DD')}
            _class='input input-sm input-bordered'
            name=date_start
            input_label="Reference Date N-1:"
        />
    </GridItem>
    <GridItem width='w-auto' _type='auto' _class='p-1'>
        <Input type=date
            defaultValue={config?.moment()?.subtract(0, 'day').format('YYYY-MM-DD')}
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

<!--- 1. BIG NUMBERS -->

```sql big_numbers
WITH logs AS (
    SELECT l.*
        , CAST((julianday(l.res_at) - julianday(l.req_at)) * 86400000 AS REAL) AS duration_miliseconds
    FROM "ADMIN.db"."user_log" l
    LEFT JOIN "ADMIN.db"."users" u ON u."user_id" = l."user_id"
    WHERE COALESCE(l."action", '') LIKE 'inputs.action.value'
        AND COALESCE(l."res_type", '') LIKE 'inputs.status.value'
        AND COALESCE(l."db", '') LIKE 'inputs.db.value'
        AND COALESCE(l."table", '') LIKE 'inputs.table.value'
),
logs_daily_hist AS (
    SELECT STRFTIME('%Y-%m-%d', "req_at") as date_ref
        , count(user_log_id) as total_day
        , count(user_log_id) filter(where "res_type" = 'success') as total_success_day
        , count(user_log_id) filter(where "res_type" = 'error') as total_error_day
        , avg(CAST((julianday(res_at) - julianday(req_at)) * 86400000 AS REAL)) as avg_duration_miliseconds_day
    FROM logs
    WHERE "req_at" < 'inputs.date_start.value'
    GROUP BY STRFTIME('%Y-%m-%d', "req_at")
),
avgs_daily_hist AS (
    SELECT avg(total_day) as avg_total
        , avg(total_success_day) as avg_total_success
        , avg(total_error_day) as avg_total_error
        , avg(avg_duration_miliseconds_day) as avg_daily_duration
    FROM logs_daily_hist
),
agg_logs AS (
    SELECT count(l.user_log_id) filter (where l."req_at" BETWEEN 'inputs.date_start.value' AND 'inputs.date_end.value') total
        , count(l.user_log_id) filter(where l."req_at" BETWEEN 'inputs.date_start.value' AND 'inputs.date_end.value' AND l."res_type" = 'success') total_success
        , count(l.user_log_id) filter(where l."req_at" BETWEEN 'inputs.date_start.value' AND 'inputs.date_end.value' AND l."res_type" = 'error') total_error
        , avg(l.duration_miliseconds) filter (where l."req_at" BETWEEN 'inputs.date_start.value' AND 'inputs.date_end.value') as avg_duration_miliseconds
    FROM logs l
)
SELECT l.total
    , case when lavg.avg_total = 0 then 0 else (l.total - lavg.avg_total) / lavg.avg_total end as total_delta
    , l.total_success
    , case when lavg.avg_total_success = 0 then 0 else (l.total_success - lavg.avg_total_success) / lavg.avg_total_success end as success_delta
    , l.total_error
    , case when lavg.avg_total = 0 then 0 else (l.total_error - lavg.avg_total_error) / lavg.avg_total_error end as error_delta
    , l.avg_duration_miliseconds
    , case when lavg.avg_daily_duration = 0 then 0 else (l.avg_duration_miliseconds - lavg.avg_daily_duration) / lavg.avg_daily_duration end as duration_delta
FROM agg_logs l
CROSS JOIN avgs_daily_hist AS lavg 
```

<Div _class="w-full p-2">
    <Stats _class='shadow' name=big_numbers_select>
        <Stat name=total parent_name=big_numbers_select bg_selected='bg-base-200'>
            <StatFigure _class='text-info p-0 w-14 h-14' icon='document-text' />
            <StatTitle _class='text-info font-bold'># Total Logs Entry</StatTitle>
            <StatValue _class=''
                data={big_numbers}
                value=total
                fmt=num0
                name=total
            />
            <StatDesc _class=''
                data={big_numbers}
                value=total_delta
                fmt=pct2
                title='vs daily average'
            />
        </Stat>
        <Stat name=total_success parent_name=big_numbers_select bg_selected='bg-base-200'>
            <StatFigure _class='text-success p-0 w-14 h-14' icon='check-circle'/>
            <StatTitle _class='text-success font-bold'># Success</StatTitle>
            <StatValue _class=''
                data={big_numbers}
                value=total_success
                fmt=num0
                name=total_success
            />
            <StatDesc _class=''
                data={big_numbers}
                value=success_delta
                fmt=pct2
                title='vs daily average'
            />
        </Stat>
        <Stat name=total_error parent_name=big_numbers_select bg_selected='bg-base-200'>
            <StatFigure _class='text-error p-0 w-14 h-14' icon='x-circle'/>
            <StatTitle _class='text-error font-bold'># Error</StatTitle>
            <StatValue _class=''
                data={big_numbers}
                value=total_error
                fmt=num0
                name=total_error
            />
            <StatDesc _class=''
                data={big_numbers}
                value=error_delta
                fmt=pct2
                title='vs daily average'
            />
        </Stat>
        <Stat name=avg_duration parent_name=big_numbers_select bg_selected='bg-base-200'>
            <StatFigure _class='text-success p-0 w-14 h-14' icon='check'/>
            <StatTitle _class='text-success font-bold'># Average Duration Response (ms)</StatTitle>
            <StatValue _class=''
                data={big_numbers}
                value=avg_duration
                fmt=num2
                name=avg_duration
            />
            <StatDesc _class='text-error'
                data={big_numbers}
                value=duration_delta
                fmt=pct2
                title='vs daily average'
            />
        </Stat>
    </Stats>
</Div>