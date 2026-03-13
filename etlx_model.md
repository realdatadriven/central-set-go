<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD031 -->
# ETLX_MODEL
```yaml
name: ETLX
description: ETLX Model
runs_as: MODEL
conn: 'sqlite3:database/ETLX.db'
admin_conn: 'sqlite3:database/ADMIN.db'
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
  etl:              { type: varchar(200), unique: true, nullable: false, comment: "Name" }
  etl_desc:         { type: text, comment: "Description" }
  attach_etlx_conf: { type: varchar(200), comment: "Config File" }
  etlx_conf:        { type: text, comment: "Config Text" }
  active:           { type: boolean, default: true, comment: "Active" }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
```

## ETLX_CONF
```yaml
table: etlx_conf
comment: ETLX Extra Cofig
columns:
  etlx_conf_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etlx_conf:       { type: varchar(200), unique: true, nullable: false, comment: "Name" }
  etlx_conf_desc:  { type: text, comment: "Description" }
  etlx_extra_conf: { type: text, comment: "Config Text" }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
```

## MANAGE_QUERY
```yaml
table: manage_query
comment: Queries
columns:
  manage_query_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  manage_query:      { type: varchar(200), nullable: false, comment: "Query Desc" }
  database:          { type: varchar(200), nullable: false, comment: "Database" }
  manage_query_conf: { type: text, comment: "Query Config" }
  active:            { type: boolean, default: true, comment: "Active" }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar(200), comment: "Dashboard" }
  dashboard_desc: { type: text, comment: "Description" }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params" }
  order:          { type: integer, comment: "Order" }
  active:         { type: boolean, default: true, comment: "Active" }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
```

## DASHBOARD_COMMENT
```yaml
table: dashboard_comment
comment: Dashboards Comments
columns:
  dashboard_comment_id: { type: integer, pk: true, autoincrement: true, comment: "Comment ID" }
  dashboard_comment:    { type: text, comment: "Comments" }
  dashboard:            { type: varchar(200), comment: "Dashboard" }
  active:               { type: boolean, default: true, comment: "Active" }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
```

## NOTEBOOK
```yaml
table: notebook
comment: Notebooks
columns:
  notebook_id:   { type: integer, pk: true, autoincrement: true, comment: "Notebook ID" }
  notebook:      { type: varchar(200), comment: "Name" }
  notebook_desc: { type: text, comment: "Description" }
  notebook_conf: { type: text, nullable: false, comment: "Conf / Params" }
  active:        { type: boolean, default: true, comment: "Active" }
  user_id:       { type: integer, comment: "User ID" }
  app_id:        { type: integer, comment: "App ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
```


# ETLX
```yaml
table: etlx
comment: ETLX
columns:
  etlx_id:          { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etl:              { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 6 }
  etl_desc:         { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: false }
  attach_etlx_conf: { type: varchar(200), comment: "Config File", form_display: true, table_display: true, form_att: true, form_size: 3 }
  etlx_conf:        { type: text, comment: "Config Text", form_display: true, form_long_text: true, form_code: markdown }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: etlx_id, order: DESC}]
table_extra_options:
  - - {size: 11, component: ETLX, label: etlx, icon: play, main: true, data: '{ "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh" }, { "type": "btn", "icon": "bolt", "name": "RUN_ALL", "class": "btn-sm text-info", "label": "crud.run_all" } ] }'}
```

# ETLX_CONF
```yaml
table: etlx_conf
comment: ETLX Extra Config
columns:
  etlx_conf_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  etlx_conf:       { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 6 }
  etlx_conf_desc:  { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  etlx_extra_conf: { type: text, comment: "Config Text", form_display: true, form_long_text: true, form_code: json }
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

# MANAGE_QUERY
```yaml
table: manage_query
comment: Queries
columns:
  manage_query_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  manage_query:      { type: varchar(200), nullable: false, comment: "Query Desc", form_display: true, table_display: true, form_size: 6 }
  database:          { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 6 }
  manage_query_conf: { type: text, comment: "Query Config", form_display: true, form_long_text: true, form_code: json }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: manage_query_id, order: DESC}]
```

# DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true, form_size: 6 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 3 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - {size: 12, component: EvidenceDash, label: dashboard, intercept_r: true}
```

# DASHBOARD_COMMENT
```yaml
table: dashboard_comment
comment: Dashboards Comments
columns:
  dashboard_comment_id: { type: integer, pk: true, autoincrement: true, comment: "Comment ID" }
  dashboard_comment:    { type: text, comment: "Comments", form_display: true, table_display: true, form_long_text: true }
  dashboard:            { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true, form_size: 6 }
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
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: dashboard_comment_id, order: DESC}]
```

# NOTEBOOK
```yaml
table: notebook
comment: Notebooks
columns:
  notebook_id:   { type: integer, pk: true, autoincrement: true, comment: "Notebook ID" }
  notebook:      { type: varchar(200), comment: "Name", form_display: true, table_display: true, form_size: 6 }
  notebook_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true }
  notebook_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:       { type: integer, comment: "User ID" }
  app_id:        { type: integer, comment: "App ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf: []
  sub_form_size: 9
table_layout:
  default_order: [{field: notebook_id, order: DESC}]
table_extra_options:
- {size: 12, component: Notebook, label: notebook, icon: book-open, intercept_r: true}
```
