<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
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
    menu_icon: user-group
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
  lang:        { type: varchar(4), unique: true, nullable: false, comment: "Language", form_display: true, table_display: true, order: 1 }
  lang_desc:   { type: varchar(200), comment: "Description", form_display: true, table_display: true, order: 2 }
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
  role:        { type: varchar(20), nullable: false, unique: true, comment: "Role", form_display: true, table_display: true, order: 1 }
  role_desc:   { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  config:      { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true, order: 3 }
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
  username:             { type: varchar(50), unique: true, nullable: false, comment: "Username", form_display: true, table_display: true, form_size: 4, order: 1 }
  first_name:           { type: varchar(50), nullable: false, comment: "First Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  last_name:            { type: varchar(50), comment: "Last Name", form_display: true, table_display: true, form_size: 4, order: 3 }
  password:             { type: varchar(200), nullable: false, comment: "Password", form_display: true, form_use_label: true, form_size: 3, order: 4 }
  email:                { type: varchar(50), unique: true, nullable: false, comment: "Email", form_display: true, table_display: true, form_size: 9, order: 5 }
  phone:                { type: varchar(50), unique: false, comment: "Phone", form_display: true, table_display: true, form_size: 3, order: 6 }
  role_id:              { type: integer, fk: "role.role_id", comment: "Default Role ID", form_display: true, table_display: true, form_size: 6, order: 7 }
  lang_id:              { type: integer, fk: "lang.lang_id", comment: "Lang ID", form_display: true, table_display: true, form_size: 3, order: 8 }
  timezone:             { type: varchar(50), comment: "Timezone", form_display: true, table_display: true, form_size: 12, order: 9 }
  attach_profile_pic:   { type: varchar(200), comment: "Profile Picture", form_display: true, table_display: true, form_size: 9, form_att: true, order: 10 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 11 }
  alter_pass_nxt_login: { type: boolean, default: false, comment: "Alter Password on next login", form_display: true, order: 15 }
  enable_2f_auth:       { type: boolean, default: false, comment: "Enable Two Factor Auth.", form_display: true, order: 16 }
  nxt_code_2f_auth:     { type: varchar(200), comment: "Next Two Factor Code", order: 17 }
  code_2f_expires_at:   { type: datetime, comment: "2F Code Expires", order: 18 }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
data:
  - {user_id: 1, username: root, password: '$2b$12$tfPUUvgU9eHTIvAy/kZo1eW2lrh2rfsX0Qx8YqomZKREoX7sUsbS6', first_name: Super, last_name: Admin, email: real.datadriven@gmail.com, role_id: 1, lang_id: 1, active: true, alter_pass_nxt_login: true, excluded: false}
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
  user_id:       { type: integer, fk: "users.user_id", nullable: false, comment: "User", table_display: true, order: 1 }
  role_id:       { type: integer, fk: "role.role_id", nullable: false, comment: "Role / Profile", form_display: true, table_display: true, order: 2 }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 12, order: 3 }
  created_at:    { type: datetime, comment: "Created at", order: 4 }
  updated_at:    { type: datetime, comment: "Updated at", order: 5 }
  excluded:      { type: boolean, default: false, comment: "Excluded", order: 6 }
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
  app:         { type: varchar(20), unique: true, nullable: false, comment: "App Name", form_display: true, table_display: true, form_size: 9, order: 1 }
  app_desc:    { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 3 }
  version:     { type: varchar(10), nullable: false, comment: "Version", form_display: true, table_display: true, form_size: 3, order: 2 }
  email:       { type: varchar(200), comment: "Email", form_display: true, table_display: true, form_size: 6, form_regex_val: "^\\w+([\\.-]?\\w+)*@\\w+([\\.-]?\\w+)*(\\.\\w{2,3})+$", order: 4 }
  db:          { type: varchar(20), nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 3, order: 5 }
  attach_logo: { type: varchar(200), comment: "Logo", form_display: true, table_display: true, form_size: 3, form_att: true, order: 6 }
  config:      { type: text, comment: "Config" }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {app_id: 1, app: ADMIN, app_desc: Admin, version: 1.0.0, db: ADMIN, user_id: 1}
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
  menu:          { type: varchar(20), unique: false, nullable: false, comment: "Menu", form_display: true, table_display: true, form_size: 12, order: 1 }
  menu_desc:     { type: text, comment: "Description", form_display: true, table_display: true, order: 2 }
  menu_icon:     { type: varchar(20), comment: "Icon", form_display: true, table_display: true, form_size: 6, order: 3 }
  menu_order:    { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 6, order: 4 }
  menu_config:   { type: text, comment: "Menu Config", form_display: true, form_long_text: true, table_display: true, form_use_label: true, order: 5 }
  config:        { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true, order: 6 }
  app_id:        { type: integer, fk: "app.app_id", comment: "App ID" }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 8 }
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
  table:            { type: varchar(100), unique: false, nullable: false, comment: "Table Name", form_display: true, table_display: true, form_size: 12, order: 1 }
  table_desc:       { type: text, comment: "Description", form_display: true, table_display: true, form_size: 12, order: 2 }
  db:               { type: varchar(50), comment: "Database / Schema", form_display: true, table_display: true, form_size: 12, order: 3 }
  requires_rla:     { type: boolean, default: false, comment: "Requires Row Level Access (RLA)", form_display: true, table_display: true, order: 4 }
  user_id:          { type: integer, fk: "users.user_id", comment: "Created/Updated by", order: 5 }
  app_id:           { type: integer, fk: "app.app_id", comment: "Application", order: 6 }
  created_at:       { type: datetime, comment: "Created at", order: 7 }
  updated_at:       { type: datetime, comment: "Updated at", order: 8 }
  excluded:         { type: boolean, default: false, comment: "Excluded", order: 9 }
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
  menu_id:        { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 1 }
  table_id:       { type: integer, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true, order: 2 }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 3 }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID", order: 4 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5 }
  requires_rla:   { type: boolean, default: false, comment: "Requires Row Level Access", form_display: true, table_display: true, order: 6 }
  menu_table_cnf: { type: text, comment: "Config", form_display: true, table_display: true, form_long_text: true, order: 7 }
  created_at:     { type: datetime, comment: "Created at", order: 8 }
  updated_at:     { type: datetime, comment: "Updated at", order: 9 }
  excluded:       { type: boolean, default: false, comment: "Excluded", order: 10 }
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
  role_id:     { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  access:      { type: boolean, default: true, comment: "Access", form_display: true, table_display: true, order: 3 }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID", order: 4 }
  created_at:  { type: datetime, comment: "Created at", order: 5 }
  updated_at:  { type: datetime, comment: "Updated at", order: 6 }
  excluded:    { type: boolean, default: false, comment: "Excluded", order: 7 }
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
  role_id:          { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:           { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  menu_id:          { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 3 }
  access:           { type: boolean, default: true, comment: "Access", form_display: true, table_display: true, order: 4 }
  user_id:          { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  created_at:       { type: datetime, comment: "Created at", order: 6 }
  updated_at:       { type: datetime, comment: "Updated at", order: 7 }
  excluded:         { type: boolean, default: false, comment: "Excluded", order: 8 }
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
  role_id:                { type: integer, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:                 { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  menu_id:                { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 3 }
  table_id:               { type: integer, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true, order: 4 }
  create:                 { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 5 }
  read:                   { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 6 }
  update:                 { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 7 }
  delete:                 { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 8 }
  share:                  { type: boolean, default: false, comment: "Share", form_display: true, table_display: true, order: 9 }
  user_id:                { type: integer, fk: "users.user_id", comment: "User ID", order: 10 }
  created_at:             { type: datetime, comment: "Created at", order: 11 }
  updated_at:             { type: datetime, comment: "Updated at", order: 12 }
  excluded:               { type: boolean, default: false, comment: "Excluded", order: 13 }
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
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID", form_display: true, table_display: true, order: 1 }
  action:      { type: varchar(200), nullable: false, comment: "Action", form_display: true, table_display: true, order: 2 }
  req_ip:      { type: varchar(200), comment: "Request IP", form_display: true, table_display: true, form_use_label: true, order: 3 }
  req_at:      { type: datetime, comment: "Request at", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm", form_use_label: true, order: 4 }
  req_data:    { type: text, comment: "Request Data", form_long_text: true, form_use_label: true, order: 5 }
  res_at:      { type: datetime, comment: "Response at", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm", form_use_label: true, order: 6 }
  res_type:    { type: varchar(200), comment: "Response Type", form_display: true, table_display: true, form_use_label: true, order: 7 }
  res_msg:     { type: varchar(500), comment: "Response Message", form_display: true, table_display: true, form_use_label: true, order: 8 }
  res_data:    { type: text, comment: "Request Data", form_long_text: true, form_use_label: true, order: 9 }
  table:       { type: varchar(200), comment: "Table", form_display: true, table_display: true, order: 10 }
  db:          { type: varchar(200), comment: "Database", form_display: true, table_display: true, order: 11 }
  row_id:      { type: integer, comment: "Database", order: 12 }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID", order: 13 }
  old_data:    { type: text, comment: "Old Data", form_long_text: true, order: 14 }
  new_data:    { type: text, comment: "New Data", form_long_text: true, order: 15 }
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
  table:           { type: varchar(200), comment: "Table", form_display: true, table_display: true, order: 1 }
  db:              { type: varchar(200), comment: "Database", form_display: true, table_display: true, order: 2 }
  config:          { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true, order: 3 }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID", order: 4 }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  created_at:      { type: datetime, comment: "Created at", order: 6 }
  updated_at:      { type: datetime, comment: "Updated at", order: 7 }
  excluded:        { type: boolean, default: false, comment: "Excluded", order: 8 }
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
  table:          { type: varchar(200), comment: "Table", form_display: true, table_display: true, order: 1 }
  db:             { type: varchar(200), comment: "Database", form_display: true, table_display: true, order: 2 }
  config:         { type: text, comment: "Config", form_display: true, form_long_text: true, table_display: true, order: 3 }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", order: 4 }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  created_at:     { type: datetime, comment: "Created at", order: 6 }
  updated_at:     { type: datetime, comment: "Updated at", order: 7 }
  excluded:       { type: boolean, default: false, comment: "Excluded", order: 8 }
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
  role_id:                  { type: integer, fk: "role.role_id", comment: "Role ID", order: 1 }
  row_id:                   { type: integer, nullable: false, comment: "Row ID", order: 2 }
  table_id:                 { type: integer, fk: "table.table_id", comment: "Table ID", order: 3 }
  table:                    { type: varchar(200), nullable: false, comment: "Table", order: 4 }
  db:                       { type: varchar(200), nullable: false, comment: "Database", order: 5 }
  user_id:                  { type: integer, fk: "users.user_id", comment: "User ID", order: 6 }
  app_id:                   { type: integer, fk: "app.app_id", comment: "App ID", order: 7 }
  read:                     { type: boolean, default: false, comment: "Read", order: 8 }
  update:                   { type: boolean, default: false, comment: "Update", order: 9 }
  delete:                   { type: boolean, default: false, comment: "Delete", order: 10 }
  share:                    { type: boolean, default: false, comment: "Share", order: 11 }
  created_at:               { type: datetime, comment: "Created at", order: 12 }
  updated_at:               { type: datetime, comment: "Updated at", order: 13 }
  excluded:                 { type: boolean, default: false, comment: "Excluded", order: 14 }
```

## COLUMN_LEVEL_ACCESS
```yaml
table: column_level_access
comment: Column Level Access
columns:
  column_level_access_id: { type: integer, pk: true, autoincrement: true, comment: "Column Level Access ID" }
  column:                 { type: integer, nullable: false, comment: "Column", order: 1 }
  table_id:               { type: integer, fk: "table.table_id", comment: "Table ID", order: 2 }
  table:                  { type: varchar(200), nullable: false, comment: "Table", order: 3 }
  db:                     { type: varchar(200), nullable: false, comment: "Database", order: 4 }
  user_id:                { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  app_id:                 { type: integer, fk: "app.app_id", comment: "App ID", order: 6 }
  create:                 { type: boolean, default: false, comment: "Create", order: 7 }
  read:                   { type: boolean, default: false, comment: "Read", order: 8 }
  update:                 { type: boolean, default: false, comment: "Update", order: 9 }
  created_at:             { type: datetime, comment: "Created at", order: 10 }
  updated_at:             { type: datetime, comment: "Updated at", order: 11 }
  excluded:               { type: boolean, default: false, comment: "Excluded", order: 12 }
```

## ROW_LEVEL_ACCESS
```yaml
table: row_level_access
comment: Row Level Access
columns:
  row_level_access_id: { type: integer, pk: true, autoincrement: true, comment: "Row Level Access ID" }
  row_id:              { type: integer, nullable: false, comment: "Row ID", order: 1 }
  table_id:            { type: integer, fk: "table.table_id", comment: "Table ID", order: 2 }
  table:               { type: varchar(200), nullable: false, comment: "Table", order: 3 }
  db:                  { type: varchar(200), nullable: false, comment: "Database", order: 4 }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID", order: 6 }
  read:                { type: boolean, default: false, comment: "Read", order: 7 }
  update:              { type: boolean, default: false, comment: "Update", order: 8 }
  delete:              { type: boolean, default: false, comment: "Delete", order: 9 }
  share:               { type: boolean, default: false, comment: "Share", order: 10 }
  created_at:          { type: datetime, comment: "Created at", order: 11 }
  updated_at:          { type: datetime, comment: "Updated at", order: 12 }
  excluded:            { type: boolean, default: false, comment: "Excluded", order: 13 }
```

## TRANSLATE_TABLE
```yaml
table: translate_table
comment: Translate Table
columns:
  transl_tbl_id:     { type: integer, pk: true, autoincrement: true, comment: "Translate Table ID" }
  table_org_desc:    { type: varchar(200), nullable: false, comment: "Table Org. Desc", form_display: true, table_display: true, order: 1 }
  table_transl_desc: { type: varchar(200), nullable: false, comment: "Table Transl. Desc", form_display: true, table_display: true, order: 2 }
  table_tooltip:     { type: varchar(500), comment: "Table Tooltip", form_display: true, table_display: true, order: 3 }
  table:             { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
  db:                { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, order: 5 }
  lang:              { type: varchar(5), nullable: false, comment: "Lang", form_display: true, table_display: true, order: 6 }
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
  field_org_desc:      { type: varchar(200), nullable: false, comment: "Field Org. Desc", form_display: true, table_display: true, order: 1 }
  field_transl_desc:   { type: varchar(200), nullable: false, comment: "Field Transl. Desc", form_display: true, table_display: true, order: 2 }
  field_tooltip:       { type: varchar(500), comment: "Field Tooltip", form_display: true, table_display: true, order: 3 }
  field:               { type: varchar(200), nullable: false, comment: "Field", form_display: true, table_display: true, order: 4 }
  table:               { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 5 }
  db:                  { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, order: 6 }
  lang:                { type: varchar(5), nullable: false, comment: "Lang", form_display: true, table_display: true, order: 7 }
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
  db:              { type: varchar(200), nullable: false, comment: "Database", form_display: true, table_display: true, order: 1 }
  table:           { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 2 }
  field:           { type: varchar(200), nullable: false, comment: "Field", form_display: true, table_display: true, order: 3 }
  type:            { type: varchar(200), nullable: false, comment: "Type", form_display: true, table_display: true, order: 4 }
  comment:         { type: varchar(200), comment: "Comment", form_display: true, table_display: true, order: 5 }
  pk:              { type: boolean, default: false, comment: "Primary Key", form_display: true, table_display: true, order: 6 }
  autoincrement:   { type: boolean, default: false, comment: "Auto Increment", form_display: true, table_display: true, order: 7 }
  nullable:        { type: boolean, default: false, comment: "Nullable", form_display: true, table_display: true, order: 8 }
  computed:        { type: boolean, default: false, comment: "Computed", form_display: true, table_display: true, order: 9 }
  default:         { type: varchar(200), comment: "Default", form_display: true, table_display: true, order: 10 }
  fk:              { type: boolean, default: false, comment: "Foreign Key", form_display: true, table_display: true, order: 11 }
  referred_table:  { type: varchar(200), comment: "Ref. Table.", form_display: true, table_display: true, order: 12 }
  referred_column: { type: varchar(200), comment: "Ref. Column", form_display: true, table_display: true, order: 13 }
  field_order:     { type: integer, comment: "Field Order", form_display: true, table_display: true, order: 14 }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID", order: 15 }
  created_at:      { type: datetime, comment: "Created at", order: 16 }
  updated_at:      { type: datetime, comment: "Updated at", order: 17 }
  excluded:        { type: boolean, default: false, comment: "Excluded", order: 18 }
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
  cron:         { type: varchar(100), unique: false, nullable: false, comment: "Cron Name", form_display: true, table_display: true, form_size: 3, order: 1 }
  cron_desc:    { type: text, nullable: false, comment: "Description", form_display: true, table_display: true, form_size: 9, order: 2 }
  api:          { type: varchar(200), nullable: false, comment: "API Endpoint / Action", form_display: true, table_display: true, form_size: 10, order: 3 }
  db:           { type: varchar(50), comment: "Database (if applicable)", order: 4 }
  table:        { type: varchar(100), comment: "Table (if applicable)", order: 5 }
  app_id:       { type: integer, fk: "app.app_id", comment: "Application ID", order: 6 }
  active:       { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 7 }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created/Updated by", order: 8 }
  created_at:   { type: datetime, comment: "Created at", order: 9 }
  updated_at:   { type: datetime, comment: "Updated at", order: 10 }
  excluded:     { type: boolean, default: false, comment: "Excluded", order: 11 }
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
  cron_id:     { type: integer, fk: "cron.cron_id", comment: "Cron ID", order: 1 }
  cron:        { type: varchar(50), nullable: false, comment: "Cron", order: 2 }
  cron_desc:   { type: varchar(200), nullable: false, comment: "Decription", order: 3 }
  api:         { type: varchar(200), nullable: false, comment: "API", order: 4 }
  start_at:    { type: datetime, comment: "Job Start", order: 5 }
  end_at:      { type: datetime, comment: "Job End", order: 6 }
  success:     { type: boolean, default: true, comment: "Success", order: 7 }
  cron_msg:    { type: text, nullable: false, comment: "Message", order: 8 }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID", order: 9 }
  db:          { type: varchar(200), comment: "Database", order: 10 }
  table:       { type: varchar(50), comment: "Table", order: 11 }
  created_at:  { type: datetime, comment: "Created at", order: 12 }
  updated_at:  { type: datetime, comment: "Updated at", order: 13 }
  excluded:    { type: boolean, default: false, comment: "Excluded", order: 14 }
```

## ACCESS_KEY
```yaml
table: access_key
comment: API / Access Tokens & Keys
columns:
  access_key_id:    { type: integer, pk: true, autoincrement: true, comment: "Access Key ID" }
  access_key_desc:  { type: varchar(200), nullable: false, comment: "Description", form_display: true, table_display: true, form_size: 12, order: 1 }
  access_token:     { type: text, nullable: false, comment: "Token / Secret", form_display: true, table_display: true, form_long_text: true, order: 2 }
  expires_at:       { type: datetime, comment: "Expires at", form_display: true, table_display: true, form_size: 4, order: 3 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4, order: 4 }
  for_user_id:      { type: integer, fk: "users.user_id", comment: "Assigned to User", form_display: true, table_display: true, form_size: 4, order: 5 }
  user_id:          { type: integer, fk: "users.user_id", comment: "Created by", order: 6 }
  app_id:           { type: integer, fk: "app.app_id", comment: "Application", order: 7 }
  created_at:       { type: datetime, comment: "Created at", order: 8 }
  updated_at:       { type: datetime, comment: "Updated at", order: 9 }
  excluded:         { type: boolean, default: false, comment: "Excluded", order: 10 }
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
  env_name:     { type: varchar(200), unique: false, nullable: false, comment: "Env Name", order: 1 }
  env_value:    { type: text, nullable: false, comment: "Env Value", order: 2 }
  on_srv_start: { type: boolean, default: true, comment: "Set On Server Start", order: 3 }
  active:       { type: boolean, default: true, comment: "Active", order: 4 }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created BY", order: 5 }
  created_at:   { type: datetime, comment: "Created at", order: 6 }
  updated_at:   { type: datetime, comment: "Updated at", order: 7 }
  excluded:     { type: boolean, default: false, comment: "Excluded", order: 8 }
```

## ARROW_FLIGHT
```yaml
table: arrow_flight
comment: Expose Arrow Flight
columns:
  arrow_flight_id:     { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight:        { type: varchar(200), unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, order: 1 }
  arrow_flight_desc:   { type: text, comment: "Description", form_display: true, table_display: true, order: 2 }
  flight_schema:       { type: varchar(200), unique: true, nullable: false, comment: "Schema Name", form_display: true, table_display: true, order: 3 }
  startup_sql:         { type: text, comment: "Startup SQL", form_display: true, order: 4 }
  main_sql:            { type: text, nullable: false, comment: "Main SQL", form_display: true, order: 5 }
  table_discover_sql:  { type: text, comment: "Table Discover SQL", form_display: true, order: 6 }
  table_scan_tmpl_sql: { type: text, comment: "Table Scan Template SQL", form_display: true, order: 7 }
  shutdown_sql:        { type: text, comment: "Shutdown SQL", form_display: true, order: 8 }
  arrow_flight_conf:   { type: text, comment: "Configuration", form_display: true, order: 9 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 10 }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID", order: 11 }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID", order: 12 }
  created_at:          { type: datetime, comment: "Created at", order: 13 }
  updated_at:          { type: datetime, comment: "Updated at", order: 14 }
  excluded:            { type: boolean, default: false, comment: "Excluded", order: 15 }
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
  arrow_flight_table_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  arrow_flight_id:         { type: integer, fk: "arrow_flight.arrow_flight_id", nullable: false, comment: "Arrow Flight", form_display: true, table_display: true, order: 1 }
  table_name:              { type: varchar(200), nullable: false, comment: "Table Name", form_display: true, table_display: true, order: 2 }
  table_desc:              { type: text, comment: "Description", form_display: true, table_display: true, order: 3 }
  order:                   { type: integer, comment: "Order", form_display: true, table_display: true, order: 4 }
  active:                  { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5 }
  arrow_flight_table_conf: { type: text, comment: "Configuration", form_display: true, order: 6 }
  user_id:                 { type: integer, fk: "users.user_id", comment: "User ID", order: 7 }
  app_id:                  { type: integer, fk: "app.app_id", comment: "App ID", order: 8 }
  created_at:              { type: datetime, comment: "Created at", order: 9 }
  updated_at:              { type: datetime, comment: "Updated at", order: 10 }
  excluded:                { type: boolean, default: false, comment: "Excluded", order: 11 }
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
  arrow_flight_table_field:      { type: varchar(200), nullable: false, comment: "Field Name", form_display: true, table_display: true, order: 1 }
  arrow_flight_table_field_desc: { type: text, comment: "Field Description", form_display: true, table_display: true, order: 2 }
  arrow_flight_table_id:         { type: integer, fk: "arrow_flight_table.arrow_flight_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true, order: 3 }
  arrow_flight_id:               { type: integer, fk: "arrow_flight.arrow_flight_id", comment: "Arrow Flight ID", order: 4 }
  active:                        { type: boolean, default: true, comment: "Active", order: 5 }
  user_id:                       { type: integer, fk: "users.user_id", comment: "User ID", order: 6 }
  app_id:                        { type: integer, fk: "app.app_id", comment: "App ID", order: 7 }
  created_at:                    { type: datetime, comment: "Created at", order: 8 }
  updated_at:                    { type: datetime, comment: "Updated at", order: 9 }
  excluded:                      { type: boolean, default: false, comment: "Excluded", order: 10 }
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
  arrow_flight_table_scope:      { type: varchar(200), unique: true, nullable: false, comment: "Scope Name", form_display: true, table_display: true, order: 1 }
  arrow_flight_table_scope_desc: { type: text, comment: "Scope Description", form_display: true, table_display: true, order: 2 }
  arrow_flight_table_scope_sql:  { type: text, nullable: false, comment: "Scope SQL", form_display: true, order: 3 }
  arrow_flight_table_id:         { type: integer, fk: "arrow_flight_table.arrow_flight_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true, order: 4 }
  arrow_flight_id:               { type: integer, fk: "arrow_flight.arrow_flight_id", comment: "Arrow Flight ID", form_display: true, table_display: true, order: 5 }
  active:                        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 6 }
  user_id:                       { type: integer, fk: "users.user_id", comment: "User ID", order: 7 }
  app_id:                        { type: integer, fk: "app.app_id", comment: "App ID", order: 8 }
  created_at:                    { type: datetime, comment: "Created at", order: 9 }
  updated_at:                    { type: datetime, comment: "Updated at", order: 10 }
  excluded:                      { type: boolean, default: false, comment: "Excluded", order: 11 }
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
  dashboard:      { type: varchar(200), comment: "Dashboard", form_display: true, table_display: true, order: 1 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, order: 2 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, table_display: true, order: 3 }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, order: 4 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5 }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID", order: 6 }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", order: 7 }
  created_at:     { type: datetime, comment: "Created at", order: 8 }
  updated_at:     { type: datetime, comment: "Updated at", order: 9 }
  excluded:       { type: boolean, default: false, comment: "Excluded", order: 10 }
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
  dashboard_comment:    { type: text, comment: "Comments", order: 1 }
  dashboard:            { type: varchar(200), comment: "Dashboard", order: 2 }
  active:               { type: boolean, default: true, comment: "Active", order: 3 }
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID", order: 4 }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID", order: 5 }
  created_at:           { type: datetime, comment: "Created at", order: 6 }
  updated_at:           { type: datetime, comment: "Updated at", order: 7 }
  excluded:             { type: boolean, default: false, comment: "Excluded", order: 8 }
```

## THROW_ERR
```yaml
table: throw_err
comment: Throw Error
columns:
  throw_err_id:     { type: integer, pk: true, autoincrement: true, comment: "Role ID" }
  throw_err:        { type: varchar(20), nullable: false, unique: true, comment: "Role", form_display: true, table_display: true, order: 1 }
  throw_err_desc:   { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
data:
  - {throw_err_id: 1, throw_err: if_empty, throw_err_desc: Throw Error if Empty, excluded: false}
  - {throw_err_id: 2, throw_err: if_not_empty, throw_err_desc: Throw Error if not Empty, excluded: false}
form_layout:
  size: 4
```

## VALIDATIONS
```yaml
table: validation
comment: Validation Roles
columns:
  validation_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  validation:      { type: varchar(200), nullable: false, comment: "Validation", form_display: true, table_display: true, order: 2, form_size: 9 }
  validation_code: { type: varchar(200), nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 2 }
  throw_err_id:    { type: integer, fk: "throw_err.throw_err_id", comment: "Throw Error ID", order: 3, form_size: 2 }
  err_msg:         { type: varchar(200), nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4 }
  table:           { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
  db:              { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
  create:          { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 5 }
  read:            { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 6 }
  update:          { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 7 }
  delete:          { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 8 }
  sql:             { type: text, nullable: false, comment: "SQL Rule", form_display: true, table_display: true, order: 4, form_long_text: true, form_code: sql }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID", order: 10 }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  created_at:      { type: datetime, comment: "Created at", order: 11 }
  updated_at:      { type: datetime, comment: "Updated at", order: 12 }
  excluded:        { type: boolean, default: false, comment: "Excluded", order: 13 }
data:
  - {validation: Validate user Email existance, validation_code: USR01, throw_err_id: 2, err_msg: "User {{.email}} already exists!", table: users, db: ADMIN, sql: "select * from users where email = :email", app_id: 1, create: true, user_id: 1}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```