<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
# ETLX_MODEL
```yaml
name: ETLX
description: ETLX Model
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    active: true
    menu_config: '{"label": "dashboard","tooltip": "dashboard_desc","load_items": {"table": "dashboard","tables": ["dashboard"]}}'
    tables:
      - dashboard
  ETLX:
    menu_icon: circle-stack
    menu_order: 2
    active: true
    tables:
      - etlx
      - etlx_conf
      - manage_query
  Notebook:
    menu_icon: book-open
    menu_order: 3
    active: true
    menu_config: '{"label": "notebook","tooltip": "notebook_desc","load_items": {"table": "notebook","tables": ["notebook"]}}'
    tables:
      - notebook
```

## ETLX
```yaml
table: etlx
comment: ETLX
columns:
  etlx_id:          { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etl:              { type: varchar, length: 200, unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 3, order: 1 }
  etl_desc:         { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: false, form_size: 6, order: 2 }
  attach_etlx_conf: { type: varchar, length: 200, comment: "Config File" }
  etlx_conf:        { type: text, comment: "Config Text", form_display: true, form_long_text: true, form_code: markdown, order: 4 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 3 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
table_layout:
  default_order: [{field: etlx_id, order: DESC}]
table_extra_options:
  - {component: ETLX, label: etlx, icon: play, size: 11, main: true, data: '{ "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh" }, { "type": "btn", "icon": "bolt", "name": "RUN_ALL", "class": "btn-sm text-info", "label": "crud.run_all" } ] }'}
```

## ETLX_CONF
```yaml
table: etlx_conf
comment: ETLX Extra Config
columns:
  etlx_conf_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etlx_conf:       { type: varchar, length: 200, unique: true, nullable: false, comment: "Name", form_display: true, table_display: true }
  etlx_conf_desc:  { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  etlx_extra_conf: { type: text, comment: "Config Text", form_display: true, table_display: true, form_long_text: true, form_code: markdown }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: etlx_conf_id, order: DESC}]
```

## MANAGE_QUERY
```yaml
table: manage_query
comment: Queries
columns:
  manage_query_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  manage_query:      { type: varchar, length: 200, nullable: false, comment: "Query Desc", form_display: true, table_display: true, form_size: 6, order: 1 }
  database:          { type: varchar, length: 200, nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 3, order: 2 }
  manage_query_conf: { type: text, comment: "Query Config", form_display: true, form_long_text: true, form_code: json, order: 4 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 3 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
table_layout:
  default_order: [{field: manage_query_id, order: DESC}]
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar, length: 200, comment: "Dashboard", form_display: true, table_display: true, form_size: 3, order: 1 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: false, form_size: 9, order: 2 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, order: 5 }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 3, order: 3 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 4 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - { component: EvidenceDash, label: dashboard, intercept_r: true, size: 12 }
```

## DASHBOARD_COMMENT
```yaml
table: dashboard_comment
comment: Dashboards Comments
columns:
  dashboard_comment_id: { type: integer, pk: true, autoincrement: true, comment: "Comment ID" }
  dashboard_comment:    { type: text, comment: "Comments", form_display: true, table_display: true, form_long_text: true, form_code: markdown }
  dashboard:            { type: varchar, length: 200, comment: "Dashboard", form_display: true, table_display: true, form_size: 6 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
table_layout:
  default_order: [{field: dashboard_comment_id, order: DESC}]
```

## NOTEBOOK
```yaml
table: notebook
comment: Notebooks
columns:
  notebook_id:   { type: integer, pk: true, autoincrement: true, comment: "Notebook ID" }
  notebook:      { type: varchar, length: 200, comment: "Name", form_display: true, table_display: true, form_size: 9, order: 1 }
  notebook_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  notebook_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, order: 4 }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 2 }
  user_id:       { type: integer, comment: "User ID" }
  app_id:        { type: integer, comment: "App ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
table_layout:
  default_order: [{ field: notebook_id, order: DESC }]
table_extra_options:
- { component: Notebook, label: notebook, icon: book-open,size: 12, intercept_r: true }
```

# DATA
```yaml
name: DATA
description: DATA Model ETLX
database: ETLX
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## DASHBOARD_LOGS
```yaml
table: dashboard
description: Add default Logs Dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id:   1
  dashboard:      Logs
  dashboard_desc: Logs Example
  dashboard_conf: FileContent(examples/logs_dashboard.md)
  order:          1
  active:         true
  user_id:        1
  app_id:        appId()
  created_at:    Now()
  updated_at:    Now()
  excluded:      false
```

## DASHBOARD_EX_TCPH_DS
```yaml
table: dashboard
description: Add Ex TCP-H DataSet
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id:   2
  dashboard:      Ex TCP-H DataSet
  dashboard_desc: This Exemple use the TCP-h Dataset generated in the ETLX Ducklake Exemple
  dashboard_conf: FileContent(examples/tcp-h-dataset-dashboard.md)
  order:          1
  active:         true
  user_id:        1
  app_id:        appId()
  created_at:    Now()
  updated_at:    Now()
  excluded:      false
```

## DASHBOARD_EX_DYN_DS
```yaml
table: dashboard
description: Add Dynamic Dataset Dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id:   3
  dashboard:      Ex Dynamic DS
  dashboard_desc: This Exemple shows how to use dynamic dataset and scoping
  dashboard_conf: FileContent(examples/dyn_ds_dashboard.md)
  order:          1
  active:         true
  user_id:        1
  app_id:        appId()
  created_at:    Now()
  updated_at:    Now()
  excluded:      false
```

## ETLX_HTTP_EX
```yaml
table: etlx
description: Add ETLX HTTP Example
cond: 'WHERE etlx_id = :etlx_id AND excluded = false'
data:
  etlx_id:    1
  etl:        HTTP_EX
  etl_desc:   HTTP Examples
  etlx_conf:  FileContent(examples/http.md)
  active:     true
  user_id:    1
  app_id:     appId()
  created_at: Now()
  updated_at: Now()
  excluded:   false
```

## ETLX_PG_EX
```yaml
table: etlx
description: Add ETLX Postgres Example
cond: 'WHERE etlx_id = :etlx_id AND excluded = false'
data:
  etlx_id:    2
  etl:        PG_EX
  etl_desc:   Postgres Examples
  etlx_conf:  FileContent(examples/pg.md)
  active:     true
  user_id:    1
  app_id:     appId()
  created_at: Now()
  updated_at: Now()
  excluded:   false
```

## ETLX_DUCKLAKE_EX
```yaml
table: etlx
description: Add ETLX ducklake Example
cond: 'WHERE etlx_id = :etlx_id AND excluded = false'
data:
  etlx_id:    3
  etl:        DUCKLAKE_EX
  etl_desc:   Ducklake Examples
  etlx_conf:  FileContent(examples/ducklake.md)
  active:     true
  user_id:    1
  app_id:     appId()
  created_at: Now()
  updated_at: Now()
  excluded:   false
```

## ETLX_SQLITE_EX
```yaml
table: etlx
description: Add SQLite Default Example
cond: 'WHERE etlx_id = :etlx_id AND excluded = false'
data:
  etlx_id:    4
  etl:        SQLITE_EX
  etl_desc:   SQLite Example
  etlx_conf:  FileContent(examples/tmpl.sqlite.md)
  active:     true
  user_id:    1
  app_id:     appId()
  created_at: Now()
  updated_at: Now()
  excluded:   false
```
