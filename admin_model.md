<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD031 -->
# ADMMIN_MODEL
```yaml
name: ADMIN
description: CS ADMIN Model
runs_as: MODEL
conn: '@DB_DRIVER_NAME:@DB_DSN'
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
  Admin:
    menu_icon: web
    menu_order: 2
    active: true
    tables:
      - app
      - menu
      - role
      - users
      - access_key
      - env
      - {table: arrow_flight, requires_rla: true, active: true}
      - {table: arrow_flight_table, active: false}
      - {table: arrow_flight_table_field, active: false}
      - {table: arrow_flight_table_scope, active: false}
      - user_log
      - custom_table
      - custom_form
      - table
      - table_schema
  Params:
    menu_icon: adjustments
    menu_order: 3
    active: true
    tables:
      - lang
  Jobs Scheduling:
    menu_icon: clock
    menu_order: 4
    active: true
    tables:
      - cron
      - {table: cron_log, active: false}
```

## LANG
```yaml
table: lang
comment: Languages
columns:
  lang_id:     { type: integer, pk: true, autoincrement: true, comment: "Lang ID" }
  lang:        { type: varchar(4), unique: true, nullable: false, comment: "Language", form_display: true, table_display: true }
  lang_desc:   { type: varchar(200), comment: "Description", form_display: true, table_display: true }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {lang_id: 1, lang: en, lang_desc: English, excluded: false}
form_layout:
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
```

## ROLE
```yaml
table: role
comment: Roles
columns:
  role_id:     { type: integer, pk: true, autoincrement: true, comment: "Role ID" }
  role:        { type: varchar(20), nullable: false, unique: true, comment: "Role", form_display: true, table_display: true }
  role_desc:   { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true }
  config:      { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {role_id: 1, role: root, role_desc: "Root role", excluded: false}
  - {role_id: 2, role: no-role, role_desc: "No role set", excluded: false}
  - {role_id: 3, role: tenant, role_desc: "Tenant Role", excluded: false}
form_layout:
  tabs_steps: tabs
  size: 8
  allow_in_subform: {}
  tabs_steps_conf: []
table_extra_options:
  - {size: 12, component: AdminApps, label: permissions, data: '{ "profile": true, "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh", "action": null }, { "type": "btn", "icon": "save", "name": "SAVE", "class": "btn-sm text-info", "label": "crud.save", "action": null } ] }', icon: key, pop_up: true, main: true}
```

## USERS
```yaml
table: users
comment: Users
columns:
  user_id:              { type: integer, pk: true, autoincrement: true, comment: "User ID" }
  username:             { type: varchar(50), unique: true, nullable: false, comment: "Username", form_display: true, table_display: true, form_size: 4 }
  first_name:           { type: varchar(50), nullable: false, comment: "First Name", form_display: true, table_display: true, form_size: 4 }
  last_name:            { type: varchar(50), comment: "Last Name", form_display: true, table_display: true, form_size: 4 }
  email:                { type: varchar(50), unique: true, nullable: false, comment: "Email", form_display: true, table_display: true, form_size: 9 }
  phone:                { type: varchar(50), unique: true, comment: "Phone", form_display: true, table_display: true, form_size: 3 }
  password:             { type: varchar(200), nullable: false, comment: "Password", form_display: true, form_size: 3 }
  role_id:              { type: integer, fk: "role.role_id", comment: "Default Role ID", form_display: true, table_display: true, form_size: 6 }
  lang_id:              { type: integer, fk: "lang.lang_id", comment: "Lang ID", form_display: true, table_display: true, form_size: 3 }
  timezone:             { type: varchar(50), comment: "Timezone", form_display: true, table_display: true, form_size: 12 }
  attach_profile_pic:   { type: varchar(200), comment: "Profile Picture", form_display: true, table_display: true, form_size: 9, form_att: true }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  alter_pass_nxt_login: { type: boolean, default: false, comment: "Alter Password on next login", form_display: true }
  enable_2f_auth:       { type: boolean, default: false, comment: "Enable Two Factor Auth.", form_display: true }
  nxt_code_2f_auth:     { type: varchar(200), comment: "Next Two Factor Code" }
  code_2f_expires_at:   { type: datetime, comment: "2F Code Expires" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
data:
  - {user_id: 1, username: root, password: '*****', first_name: Super, last_name: Admin, email: real.datadriven@gmail.com, role_id: 1, lang_id: 1, active: true, alter_pass_nxt_login: true, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
  allow_in_subform: {user_role: true}
table_extra_options:
  - {size: 12, component: AdminApps, label: permissions, data: '{ "profile": true, "user_rla": true, "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh", "action": null }, { "type": "btn", "icon": "save", "name": "SAVE", "class": "btn-sm text-info", "label": "crud.save", "action": null } ] }', pop_up: true, main: true, icon: key}
```

## USER_ROLE
```yaml
table: user_role
comment: User-Role Assignments
columns:
  user_role_id:  { type: integer, pk: true, autoincrement: true, comment: "User-Role Assignment ID" }
  user_id:       { type: integer, fk: "users.user_id", nullable: false, comment: "User", form_display: false, table_display: true }
  role_id:       { type: integer, fk: "role.role_id", nullable: false, comment: "Role / Profile", form_display: true, table_display: true }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 12 }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 4
```

## APP
```yaml
table: app
comment: Applications
columns:
  app_id:      { type: integer, pk: true, autoincrement: true, comment: "App ID" }
  app:         { type: varchar(20), unique: true, nullable: false, comment: "App Name", form_display: true, table_display: true, form_size: 9 }
  app_desc:    { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true }
  version:     { type: varchar(10), nullable: false, comment: "Version", form_display: true, table_display: true, form_size: 3 }
  email:       { type: varchar(200), comment: "Email", form_display: true, table_display: true, form_size: 6, form_regex_val: "^\\w+([\\.-]?\\w+)*@\\w+([\\.-]?\\w+)*(\\.\\w{2,3})+$" }
  db:          { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 3 }
  attach_logo: { type: varchar(200), comment: "Logo", form_display: true, table_display: true, form_size: 3, form_att: true }
  config:      { type: text, comment: "Config", form_display: false, form_long_text: true }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {app_id: 1, app: ADMIN, app_desc: Admin, version: 1.0.0, user_id: 1, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {menu: true, role_app: true}
  tabs_steps_conf: []
  sub_form_size: 9
form_extra_options: []
table_layout:
  allow_in_submenu: {menu: true}
  default_order: [{field: app_id, order: DESC}]
  allow_import: false
table_extra_options:
  - {size: 12, component: AdminApps, label: admin_apps, data: '{ "profile": false, "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh", "action": null }, { "type": "btn", "icon": "save", "name": "SAVE", "class": "btn-sm text-info", "label": "crud.save", "action": null } ] }', icon: cog, pop_up: false, main: true}
```

## MENU
```yaml
table: menu
comment: Menu Items
columns:
  menu_id:       { type: integer, pk: true, autoincrement: true, comment: "Menu ID" }
  menu:          { type: varchar(20), unique: true, nullable: false, comment: "Menu", form_display: true, table_display: true, form_size: 12 }
  menu_desc:     { type: text, comment: "Description", form_display: true, table_display: true }
  menu_icon:     { type: varchar(20), comment: "Icon", form_display: true, table_display: true, form_size: 6 }
  menu_order:    { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 6 }
  menu_config:   { type: text, comment: "Menu Config", form_display: true, form_long_text: true, table_display: true }
  config:        { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true }
  app_id:        { type: integer, fk: "app.app_id", comment: "App ID", form_display: false, table_display: true }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true }
  user_id:       { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
  allow_in_subform: {}
  tabs_steps_conf: []
table_layout:
  allow_in_submenu: {menu_table: true}
  default_order: [{field: menu_order, order: ASC}]
  allow_import: false
```

## TABLE
```yaml
table: table
comment: Tables
columns:
  table_id:         { type: integer, pk: true, autoincrement: true, comment: "Table ID" }
  table:            { type: varchar(100), unique: true, nullable: false, comment: "Table Name", form_display: true, table_display: true, form_size: 12 }
  table_desc:       { type: text, comment: "Description", form_display: true, table_display: true, form_size: 12 }
  db:               { type: varchar(50), comment: "Database / Schema", form_display: true, table_display: true, form_size: 12 }
  requires_rla:     { type: boolean, default: false, comment: "Requires Row Level Access (RLA)", form_display: true, table_display: true }
  user_id:          { type: integer, fk: "users.user_id", comment: "Created/Updated by" }
  app_id:           { type: integer, fk: "app.app_id", comment: "Application" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 12
```

## MENU_TABLE
```yaml
table: menu_table
comment: Menu Tables
columns:
  menu_table_id:  { type: integer, pk: true, autoincrement: true, comment: "Menu Table ID" }
  menu_id:        { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true }
  table_id:       { type: integer, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID" }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true }
  requires_rla:   { type: boolean, default: false, comment: "Requires Row Level Access", form_display: true, table_display: true }
  menu_table_cnf: { type: text, comment: "Config", form_display: true, table_display: true, form_long_text: true }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## ROLE_APP
```yaml
table: role_app
comment: Role Apps
columns:
  role_app_id: { type: integer, pk: true, autoincrement: true, comment: "Role App ID" }
  role_id:     { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true }
  access:      { type: boolean, default: true, comment: "Access", form_display: true, table_display: true }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## ROLE_APP_MENU
```yaml
table: role_app_menu
comment: Role App Menus
columns:
  role_app_menu_id: { type: integer, pk: true, autoincrement: true, comment: "Role App Menu ID" }
  role_id:          { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true }
  app_id:           { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true }
  menu_id:          { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true }
  access:           { type: boolean, default: true, comment: "Access", form_display: true, table_display: true }
  user_id:          { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## ROLE_APP_MENU_TABLE
```yaml
table: role_app_menu_table
comment: Role App Menu Tables
columns:
  role_app_menu_table_id: { type: integer, pk: true, autoincrement: true, comment: "Role App Menu Table ID" }
  role_id:                { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true }
  app_id:                 { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true }
  menu_id:                { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true }
  table_id:               { type: integer, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true }
  create:                 { type: boolean, default: false, comment: "Create", form_display: true, table_display: true }
  read:                   { type: boolean, default: false, comment: "Read", form_display: true, table_display: true }
  update:                 { type: boolean, default: false, comment: "Update", form_display: true, table_display: true }
  delete:                 { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true }
  share:                  { type: boolean, default: false, comment: "Share", form_display: true, table_display: true }
  user_id:                { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:             { type: datetime, comment: "Created at" }
  updated_at:             { type: datetime, comment: "Updated at" }
  excluded:               { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## USER_LOG
```yaml
table: user_log
comment: User Logs
columns:
  user_log_id: { type: integer, pk: true, autoincrement: true, comment: "User Log ID" }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID", form_display: true, table_display: true }
  action:      { type: varchar(200), nullable: false, comment: "Action", form_display: true, table_display: true }
  req_ip:      { type: varchar(200), comment: "Request IP", form_display: true, table_display: true }
  req_at:      { type: datetime, comment: "Request at", form_display: true, table_display: true }
  req_data:    { type: text, comment: "Request Data" }
  res_at:      { type: datetime, comment: "Response at", form_display: true, table_display: true }
  res_type:    { type: varchar(200), comment: "Response Type", form_display: true, table_display: true }
  res_msg:     { type: varchar(500), comment: "Response Message", form_display: true, table_display: true }
  res_data:    { type: text, comment: "Request Data" }
  table:       { type: varchar(200), comment: "Table", form_display: true, table_display: true }
  db:          { type: varchar(200), comment: "Database", form_display: true, table_display: true }
  row_id:      { type: integer, comment: "Database" }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID" }
  old_data:    { type: text, comment: "Old Data" }
  new_data:    { type: text, comment: "New Data" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## CUSTOM_TABLE
```yaml
table: custom_table
comment: Custom Table
columns:
  custom_table_id: { type: integer, pk: true, autoincrement: true, comment: "Custom Table ID" }
  table:           { type: varchar(200), comment: "Table", form_display: true, table_display: true }
  db:              { type: varchar(200), comment: "Database", form_display: true, table_display: true }
  config:          { type: text, comment: "Config", form_display: true, form_long_text: true, form_code: "json", table_display: true }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID" }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

## CUSTOM_FORM
```yaml
table: custom_form
comment: Custom Form
columns:
  custom_form_id: { type: integer, pk: true, autoincrement: true, comment: "Custom Form ID" }
  table:          { type: varchar(200), comment: "Table", form_display: true, table_display: true }
  db:             { type: varchar(200), comment: "Database", form_display: true, table_display: true }
  config:         { type: text, comment: "Config", form_display: true, form_long_text: true, form_code: "json", table_display: true  }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID" }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

## ROLE_ROW_LEVEL_ACCESS
```yaml
table: role_row_level_access
comment: Role Row Level Access
columns:
  role_row_level_access_id: { type: integer, pk: true, autoincrement: true, comment: "Role Row Level Access ID" }
  role_id:                  { type: integer, fk: "role.role_id", comment: "Role ID" }
  row_id:                   { type: integer, nullable: false, comment: "Row ID" }
  table_id:                 { type: integer, fk: "table.table_id", comment: "Table ID" }
  table:                    { type: varchar(200), nullable: false, comment: "Table" }
  db:                       { type: varchar(200), nullable: false, comment: "Database" }
  user_id:                  { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                   { type: integer, fk: "app.app_id", comment: "App ID" }
  read:                     { type: boolean, default: false, comment: "Read" }
  update:                   { type: boolean, default: false, comment: "Update" }
  delete:                   { type: boolean, default: false, comment: "Delete" }
  share:                    { type: boolean, default: false, comment: "Share" }
  created_at:               { type: datetime, comment: "Created at" }
  updated_at:               { type: datetime, comment: "Updated at" }
  excluded:                 { type: boolean, default: false, comment: "Excluded" }
```

## COLUMN_LEVEL_ACCESS
```yaml
table: column_level_access
comment: Column Level Access
columns:
  column_level_access_id: { type: integer, pk: true, autoincrement: true, comment: "Column Level Access ID" }
  column:                 { type: integer, nullable: false, comment: "Column" }
  table_id:               { type: integer, fk: "table.table_id", comment: "Table ID" }
  table:                  { type: varchar(200), nullable: false, comment: "Table" }
  db:                     { type: varchar(200), nullable: false, comment: "Database" }
  user_id:                { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                 { type: integer, fk: "app.app_id", comment: "App ID" }
  create:                 { type: boolean, default: false, comment: "Create" }
  read:                   { type: boolean, default: false, comment: "Read" }
  update:                 { type: boolean, default: false, comment: "Update" }
  created_at:             { type: datetime, comment: "Created at" }
  updated_at:             { type: datetime, comment: "Updated at" }
  excluded:               { type: boolean, default: false, comment: "Excluded" }
```

## ROW_LEVEL_ACCESS
```yaml
table: row_level_access
comment: Row Level Access
columns:
  row_level_access_id: { type: integer, pk: true, autoincrement: true, comment: "Row Level Access ID" }
  row_id:              { type: integer, nullable: false, comment: "Row ID" }
  table_id:            { type: integer, fk: "table.table_id", comment: "Table ID" }
  table:               { type: varchar(200), nullable: false, comment: "Table" }
  db:                  { type: varchar(200), nullable: false, comment: "Database" }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  read:                { type: boolean, default: false, comment: "Read" }
  update:              { type: boolean, default: false, comment: "Update" }
  delete:              { type: boolean, default: false, comment: "Delete" }
  share:               { type: boolean, default: false, comment: "Share" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
```

## TRANSLATE_TABLE
```yaml
table: translate_table
comment: Translate Table
columns:
  transl_tbl_id:     { type: integer, pk: true, autoincrement: true, comment: "Translate Table ID" }
  table_org_desc:    { type: varchar(200), nullable: false, comment: "Table Org. Desc", form_display: true, table_display: true }
  table_transl_desc: { type: varchar(200), nullable: false, comment: "Table Transl. Desc", form_display: true, table_display: true }
  table:             { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true }
  db:                { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true }
  lang:              { type: varchar(5), nullable: false, comment: "Lang", form_display: true, table_display: true }
  user_id:           { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:            { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

## TRANSLATE_TABLE_FIELD
```yaml
table: translate_table_field
comment: Translate Table Fields
columns:
  transl_tbl_field_id: { type: integer, pk: true, autoincrement: true, comment: "Translate Table Field ID" }
  field_org_desc:      { type: varchar(200), nullable: false, comment: "Field Org. Desc", form_display: true, table_display: true }
  field_transl_desc:   { type: varchar(200), nullable: false, comment: "Field Transl. Desc", form_display: true, table_display: true }
  field:               { type: varchar(200), nullable: false, comment: "Field", form_display: true, table_display: true }
  table:               { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true }
  db:                  { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true }
  lang:                { type: varchar(5), nullable: false, comment: "Lang", form_display: true, table_display: true }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

## TABLE_SCHEMA
```yaml
table: table_schema
comment: Table Schema
columns:
  table_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Table field ID" }
  db:              { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true }
  table:           { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true }
  field:           { type: varchar(200), nullable: false, comment: "Field", form_display: true, table_display: true }
  type:            { type: varchar(200), nullable: false, comment: "Type", form_display: true, table_display: true }
  comment:         { type: varchar(200), comment: "Comment", form_display: true, table_display: true }
  pk:              { type: boolean, default: false, comment: "Primary Key", form_display: true, table_display: true }
  autoincrement:   { type: boolean, default: false, comment: "Auto Increment", form_display: true, table_display: true }
  nullable:        { type: boolean, default: false, comment: "Nullable", form_display: true, table_display: true }
  computed:        { type: boolean, default: false, comment: "Computed", form_display: true, table_display: true }
  default:         { type: varchar(200), comment: "Default", form_display: true, table_display: true }
  fk:              { type: boolean, default: false, comment: "Foreign Key", form_display: true, table_display: true }
  referred_table:  { type: varchar(200), comment: "Ref. Table.", form_display: true, table_display: true }
  referred_column: { type: varchar(200), comment: "Ref. Column", form_display: true, table_display: true }
  field_order:     { type: integer, comment: "Field Order", form_display: true, table_display: true }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
```

## CRON
```yaml
table: cron
comment: Jobs scheduling
columns:
  cron_id:      { type: integer, pk: true, autoincrement: true, comment: "Cron ID" }
  cron:         { type: varchar(100), unique: true, nullable: false, comment: "Cron Name", form_display: true, table_display: true, form_size: 3 }
  cron_desc:    { type: text, comment: "Description", form_display: true, table_display: true, form_size: 9 }
  api:          { type: varchar(200), nullable: false, comment: "API Endpoint / Action", form_display: true, table_display: true, form_size: 10 }
  db:           { type: varchar(50), comment: "Database (if applicable)" }
  table:        { type: varchar(100), comment: "Table (if applicable)" }
  app_id:       { type: integer, fk: "app.app_id", comment: "Application ID" }
  active:       { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2 }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created/Updated by" }
  created_at:   { type: datetime, comment: "Created at" }
  updated_at:   { type: datetime, comment: "Updated at" }
  excluded:     { type: boolean, default: false, comment: "Excluded" }
data:
  - {cron_id: 1, cron: "0 0 * * *", cron_desc: Backup, api: buckup, app_id: 1, db: ADMIN, active: false, excluded: false}
  - {cron_id: 2, cron: "0 0 * * *", cron_desc: "Update Env", api: "env/sync", app_id: 1, db: ADMIN, active: false, excluded: false}
  - {cron_id: 3, cron: "0 0 * * *", cron_desc: "ETLX Example", api: "etlx/name/[etlx_name]", app_id: 1, db: ADMIN, active: false, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  sub_form_size: null
  sub_form_limit: 10
  allow_in_subform: {cron_log: true}
  tabs_steps_conf: []
table_layout:
  allow_in_submenu: {}
  default_order: []
  allow_import: false
  exec_button: {callApi: true, api: cron/run, tooltip: Run Backup, icon: play}
```

## CRON_LOG
```yaml
table: cron_log
comment: Jobs scheduling logs
columns:
  cron_log_id: { type: integer, pk: true, autoincrement: true, comment: "Cron Log ID" }
  cron_id:     { type: integer, fk: "cron.cron_id", comment: "Cron ID" }
  cron:        { type: varchar(50), nullable: false, comment: "Cron" }
  cron_desc:   { type: varchar(200), nullable: false, comment: "Decription" }
  api:         { type: varchar(200), nullable: false, comment: "API" }
  start_at:    { type: datetime, comment: "Job Start" }
  end_at:      { type: datetime, comment: "Job End" }
  success:     { type: boolean, default: true, comment: "Success" }
  cron_msg:    { type: text, nullable: false, comment: "Message" }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID" }
  db:          { type: varchar(200), comment: "Database" }
  table:       { type: varchar(50), comment: "Table" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
```

## ACCESS_KEY
```yaml
table: access_key
comment: API / Access Tokens & Keys
columns:
  access_key_id:    { type: integer, pk: true, autoincrement: true, comment: "Access Key ID" }
  access_key_desc:  { type: varchar(200), nullable: false, comment: "Description", form_display: true, table_display: true, form_size: 12 }
  access_token:     { type: text, nullable: false, comment: "Token / Secret", form_display: true, table_display: true, form_long_text: true, form_readonly: true }
  expires_at:       { type: datetime, comment: "Expires at", form_display: true, table_display: true, form_size: 4 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4 }
  for_user_id:      { type: integer, fk: "users.user_id", comment: "Assigned to User", form_display: true, table_display: true, form_size: 4, form_fk_label: "username" }
  user_id:          { type: integer, fk: "users.user_id", comment: "Created by" }
  app_id:           { type: integer, fk: "app.app_id", comment: "Application" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
  sub_form_size: 6
  sub_form_limit: 5
table_extra_options:
  - size: 6
    component: AccessKey
    label: accessKey
    intercept_c: true
    intercept_u: true
    data: '{ "actions": [ {"type": "btn", "icon": "refresh", "name": "refresh", "class": "btn-sm text-info", "label": "crud.refresh", "action": null}, {"type": "btn", "icon": "save", "name": "save", "class": "btn-sm text-info", "label": "crud.save", "action": null }, {"type": "icon", "icon": "cog-8-tooth", "name": "form_customization", "label": "crud.form_customization", "action": null} ] }'
    main: true
    icon: key
```

## ENV
```yaml
table: env
comment: Envariomental Variables
columns:
  env_id:       { type: integer, pk: true, autoincrement: true, comment: "env ID" }
  env_name:     { type: varchar(200), unique: true, nullable: false, comment: "Env Name" }
  env_value:    { type: text, nullable: false, comment: "Env Value" }
  on_srv_start: { type: boolean, default: true, comment: "Set On Server Start" }
  active:       { type: boolean, default: true, comment: "Active" }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created BY" }
  created_at:   { type: datetime, comment: "Created at" }
  updated_at:   { type: datetime, comment: "Updated at" }
  excluded:     { type: boolean, default: false, comment: "Excluded" }
```

## ARROW_FLIGHT
```yaml
table: arrow_flight
comment: Expose Arrow Flight
columns:
  arrow_flight_id:     { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight:        { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true }
  arrow_flight_desc:   { type: text, comment: "Description", form_display: true, table_display: true, form_code: markdown }
  flight_schema:       { type: varchar(200), unique: true, nullable: false, comment: "Schema Name", form_display: true, table_display: true }
  startup_sql:         { type: text, comment: "Startup SQL", form_display: true, form_code: sql }
  main_sql:            { type: text, nullable: false, comment: "Main SQL", form_display: true, form_code: sql }
  table_discover_sql:  { type: text, comment: "Table Discover SQL", form_display: true, form_code: sql }
  table_scan_tmpl_sql: { type: text, comment: "Table Scan Template SQL", form_display: true, form_code: sql }
  shutdown_sql:        { type: text, comment: "Shutdown SQL", form_display: true, form_code: sql }
  arrow_flight_conf:   { type: text, comment: "Configuration", form_display: true, form_code: json }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_sizelg: 3, form_sizexl: 3 }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {arrow_flight_id: 1, arrow_flight: "Expose Admin DB", arrow_flight_desc: "Ex. Arrow Flight Schema using ADMIN app", flight_schema: adm, startup_sql: "INSTALL SQLITE;LOAD SQLITE;", main_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE);USE adm;", shutdown_sql: "USE memory;DETACH adm;", active: false, app_id: 1, user_id: 1, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  sub_form_size: 9
  allow_in_subform: {arrow_flight_table: true, arrow_flight_scope: true}
  tabs_steps_conf: []
form_extra_options: []
table_layout:
  allow_in_submenu: {}
  default_order:
    - { field: arrow_flight, order: ASC }
```

## ARROW_FLIGHT_TABLE
```yaml
table: arrow_flight_table
comment: Arrow Flight Tables
columns:
  arrow_flight_table_id: { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight_id:       { type: integer, fk: "arrow_flight.arrow_flight_id", nullable: false, comment: "Arrow Flight", form_display: true, table_display: true }
  table_name:            { type: varchar(200), nullable: false, comment: "Table Name", form_display: true, table_display: true }
  table_desc:            { type: text, comment: "Description", form_display: true, table_display: true, form_code: markdown }
  order:                 { type: integer, comment: "Order", form_display: true, table_display: true, form_sizelg: 3, form_sizexl: 3 }
  active:                { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_sizelg: 3, form_sizexl: 3 }
  arrow_flight_table_conf:{ type: text, comment: "Configuration", form_display: true, form_code: json }
  user_id:               { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  sub_form_size: 8
  allow_in_subform:
    arrow_flight_table_field: true
table_layout:
  allow_in_submenu: {}
  default_order:
    - { field: order, order: ASC }
  allow_import: false
```

## ARROW_FLIGHT_TABLE_FIELD
```yaml
table: arrow_flight_table_field
comment: Arrow Flight - Tables Fields
columns:
  arrow_flight_table_field_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight_table_field:      { type: varchar(200), nullable: false, comment: "Field Name", form_display: true, table_display: true }
  arrow_flight_table_field_desc: { type: text, comment: "Field Description", form_display: true, table_display: true, form_code: markdown }
  arrow_flight_table_id:         { type: integer, fk: "arrow_flight_table.arrow_flight_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true }
  arrow_flight_id:               { type: integer, fk: "arrow_flight.arrow_flight_id", comment: "Arrow Flight ID" }
  active:                        { type: boolean, default: true, comment: "Active", form_sizelg: 3, form_sizexl: 3 }
  user_id:                       { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                        { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:                    { type: datetime, comment: "Created at" }
  updated_at:                    { type: datetime, comment: "Updated at" }
  excluded:                      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  sub_form_size: null
  allow_in_subform: {}
  tabs_steps_conf: []
table_layout:
  allow_in_submenu: {}
  default_order:
    - { field: arrow_flight_table_field, order: ASC }
  allow_import: false
```

## ARROW_FLIGHT_TABLE_SCOPE
```yaml
table: arrow_flight_table_scope
comment: Arrow Flight - Tables Scopes
columns:
  arrow_flight_table_scope_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight_table_scope:      { type: varchar(200), unique: true, nullable: false, comment: "Scope Name", form_display: true, table_display: true }
  arrow_flight_table_scope_desc: { type: text, comment: "Scope Description", form_display: true, table_display: true, form_code: markdown }
  arrow_flight_table_scope_sql:  { type: text, nullable: false, comment: "Scope SQL", form_display: true, form_code: sql }
  arrow_flight_table_id:         { type: integer, fk: "arrow_flight_table.arrow_flight_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true }
  arrow_flight_id:               { type: integer, fk: "arrow_flight.arrow_flight_id", comment: "Arrow Flight ID", form_display: true, table_display: true }
  active:                        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_sizelg: 3, form_sizexl: 3 }
  user_id:                       { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                        { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:                    { type: datetime, comment: "Created at" }
  updated_at:                    { type: datetime, comment: "Updated at" }
  excluded:                      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  allow_in_submenu: {}
  default_order:
    - { field: arrow_flight_table_scope, order: ASC }
  allow_import: false
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_code: markdown, table_display: true }
  order:          { type: integer, comment: "Order", form_display: true, form_sizelg: 3, form_sizexl: 3, table_display: true }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, form_sizelg: 3, form_sizexl: 3, table_display: true }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
table_extra_options:
  - {size: 12, component: EvidenceDash, label: dashboard, intercept_r: true}
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
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
```
