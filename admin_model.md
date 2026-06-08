<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
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
      - {table: user_role, active: false}
      - access_key
      - env
      - validation
      - {table: valid_reaction, active: false}
      - user_log
      - custom_table
      - custom_form
      - table
      - table_schema
      - crud_action
      - {table: crud_action_logs, active: false}
  Arrow Flight:
    menu_icon: paper-airplane
    menu_order: 3
    active: true
    #menu_config: '{"label": "flight_catalog","tooltip": "flight_catalog_desc", "load_items": {"table": "flight_catalog","tables": ["flight_catalog"]}}'
    menu_config: |
      {
          "label": "flight_catalog",
          "tooltip": "flight_catalog_desc",
          "load_items": {
              "table": "flight_catalog",
              "label": "flight_catalog",
              "tooltip": "flight_catalog_desc",
              "detail": true,
              "load_items": {
                  "table": "flight_schema",
                  "label": "flight_schema",
                  "tooltip": "flight_schema_desc",
                  "detail": false
              }, 
              "tables": ["flight_catalog", "flight_schema"]
          }
      }
    tables:
      - {table: flight_catalog, requires_rla: true, active: true}
      - {table: flight_schema, requires_rla: true, active: false}
      - {table: flight_schema_table, active: false}
      - {table: flight_schema_table_field, active: false}
      - {table: flight_schema_table_scope, active: false}
  Quack:
    menu_icon: bolt
    menu_order: 4
    active: false
    tables:
      - {table: quack_server, requires_rla: true, active: true}
      - {table: quack_logs, active: false}
  Jobs Scheduling:
    menu_icon: clock
    menu_order: 5
    active: true
    tables:
      - cron
      - {table: cron_log, active: false}  
  APIs:
    menu_icon: share
    menu_order: 6
    active: true
    tables:
      - {table: api_type, active: false}
      - {table: http_request_type, active: false}
      - api
      - {table: api_header, active: false}
      - {table: api_call_log, active: false}
  Params:
    menu_icon: adjustments
    menu_order: 7
    active: true
    tables:
      - lang
```

## LANG
```yaml
table: lang
comment: Languages
columns:
  lang_id:     { type: integer, pk: true, autoincrement: true, comment: "Lang ID" }
  lang:        { type: varchar, len: 4, unique: true, nullable: false, comment: "Language", form_display: true, table_display: true, order: 1 }
  lang_desc:   { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {lang_id: 1, lang: en, lang_desc: English, excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## ROLE
```yaml
table: role
comment: Roles
columns:
  role_id:     { type: integer, pk: true, autoincrement: true, comment: "Role ID" }
  role:        { type: varchar, len: 20, nullable: false, unique: true, comment: "Role", form_display: true, table_display: true, order: 1 }
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
  size: 4
table_extra_options:
  - {size: 12, component: AdminApps, label: permissions, data: '{ "profile": true, "actions": [ { "type": "btn", "icon": "refresh", "name": "REFRESH", "class": "btn-sm text-info", "label": "crud.refresh", "action": null }, { "type": "btn", "icon": "save", "name": "SAVE", "class": "btn-sm text-info", "label": "crud.save", "action": null } ] }', icon: key, pop_up: true, main: true}
```

## USERS
```yaml
table: users
comment: Users
columns:
  user_id:              { type: integer, pk: true, autoincrement: true, comment: "User ID" }
  username:             { type: varchar, len: 50, unique: true, nullable: false, comment: "Username", form_display: true, table_display: true, form_size: 4, order: 1 }
  first_name:           { type: varchar, len: 50, nullable: false, comment: "First Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  last_name:            { type: varchar, len: 50, comment: "Last Name", form_display: true, table_display: true, form_size: 4, order: 3 }
  password:             { type: varchar, len: 200, nullable: false, comment: "Password", form_display: true, form_use_label: true, form_size: 4, order: 4 }
  email:                { type: varchar, len: 50, unique: true, nullable: false, comment: "Email", form_display: true, table_display: true, form_size: 4, order: 5 }
  phone:                { type: varchar, len: 50, unique: false, comment: "Phone", form_display: true, table_display: true, form_size: 4, order: 6 }
  role_id:              { type: integer, fk: "role.role_id", comment: "Default Role ID", form_display: true, table_display: true, form_size: 4, order: 7 }
  lang_id:              { type: integer, fk: "lang.lang_id", comment: "Lang ID", form_display: true, table_display: true, form_size: 4, order: 8 }
  timezone:             { type: varchar, len: 50, comment: "Timezone", form_display: true, form_size: 4, order: 9 }
  attach_profile_pic:   { type: varchar, len: 200, comment: "Profile Picture", form_display: true, table_display: true, form_size: 3, form_att: true, order: 10 }
  failed_login_attmpt:  { type: integer, comment: "# Failed Login Attempts", form_display: true, table_display: false, form_size: 3, form_att: true, order: 10 }
  last_failed_login:    { type: datetime, comment: "Last Failed Login Attempts", form_display: true, table_display: false, form_size: 3, form_att: true, order: 10 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 11 }
  alter_pass_nxt_login: { type: boolean, default: false, comment: "Alter Password on next login", form_display: true, order: 12, form_size: 4 }
  enable_2f_auth:       { type: boolean, default: false, comment: "Enable Two Factor Auth.", form_display: true, order: 13, form_size: 3 }
  nxt_code_2f_auth:     { type: varchar, len: 200, comment: "Next Two Factor Code", order: 14 }
  code_2f_expires_at:   { type: datetime, comment: "2F Code Expires", order: 18 }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
data:
  - {user_id: 1, username: root, password: '$2b$12$tfPUUvgU9eHTIvAy/kZo1eW2lrh2rfsX0Qx8YqomZKREoX7sUsbS6', first_name: Super, last_name: Admin, email: admin@domain.com, role_id: 1, lang_id: 1, active: true, alter_pass_nxt_login: true, excluded: false}
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
  user_id:       { type: integer, fk: "users.user_id", nullable: false, comment: "User", table_display: true, form_display: true, order: 1 }
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
  app:         { type: varchar, len: 20, unique: true, nullable: false, comment: "App Name", form_display: true, table_display: true, form_size: 9, order: 1 }
  app_desc:    { type: text, comment: "Description", form_display: true, form_long_text: true, form_code: markdown, form_rendermd: true, table_display: true, order: 3 }
  version:     { type: varchar, len: 10, nullable: false, comment: "Version", form_display: true, table_display: true, form_size: 3, order: 2 }
  email:       { type: varchar, len: 200, comment: "Email", form_display: true, table_display: true, form_size: 3, form_regex_val: "^\\w+([\\.-]?\\w+)*@\\w+([\\.-]?\\w+)*(\\.\\w{2,3})+$", order: 4 }
  db:          { type: varchar, len: 20, nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 3, order: 5 }
  # conn_string: { type: varchar, len: 200, nullable: false, comment: "Conn String", form_display: true, table_display: true, form_size: 3, order: 5 }
  attach_logo: { type: varchar, len: 200, comment: "Logo", form_display: true, table_display: true, form_size: 3, form_att: true, order: 6 }
  category:    { type: varchar, len: 200, comment: "Category", form_display: true, table_display: true, form_size: 3, order: 6 }
  config:      { type: text, comment: "Config", form_long_text: true, form_code: json }
  user_id:     { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {app_id: 1, app: ADMIN, app_desc: Admin, version: 1.0.0, db: ADMIN, category: Admin, user_id: 1}
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
  menu:          { type: varchar, len: 20, nullable: false, comment: "Menu", form_display: true, table_display: true, form_size: 3, order: 1 }
  menu_desc:     { type: text, comment: "Description", form_display: true, form_size: 9, table_display: true, order: 2 }
  menu_icon:     { type: varchar, len: 20, comment: "Icon", form_display: true, table_display: true, form_size: 4, order: 3 }
  menu_order:    { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 4, order: 4 }
  menu_config:   { type: text, comment: "Menu Config", form_display: true, form_long_text: true, form_code: json, table_display: true, form_use_label: true, order: 6 }
  app_id:        { type: integer, fk: "app.app_id", comment: "App ID" }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, form_size: 4, table_display: true, order: 5 }
  user_id:       { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
  allow_in_subform: {menu_table: true}
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
  table_id:     { type: integer, pk: true, autoincrement: true, comment: "Table ID" }
  table:        { type: varchar, len: 100, unique: false, nullable: false, comment: "Table Name", form_display: true, table_display: true, form_size: 12, order: 1 }
  table_desc:   { type: text, comment: "Description", form_display: true, table_display: true, form_size: 12, order: 2 }
  db:           { type: varchar, len: 50, comment: "Database / Schema", form_display: true, table_display: true, form_size: 12, order: 3 }
  requires_rla: { type: boolean, default: false, comment: "Requires Row Level Access (RLA)", form_display: true, table_display: true, order: 4 }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created/Updated by", order: 5 }
  app_id:       { type: integer, fk: "app.app_id", comment: "Application", order: 6 }
  created_at:   { type: datetime, comment: "Created at", order: 7 }
  updated_at:   { type: datetime, comment: "Updated at", order: 8 }
  excluded:     { type: boolean, default: false, comment: "Excluded", order: 9 }
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
  menu_id:        { type: integer, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 1, form_size: 8 }
  table_id:       { type: integer, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true, order: 2, form_size: 4 }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 3, form_size: 4 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5, form_size: 4 }
  requires_rla:   { type: boolean, default: false, comment: "Requires Row Level Access", form_display: true, table_display: true, order: 6, form_size: 4 }
  menu_table_cnf: { type: text, comment: "Config", form_display: true, table_display: true, form_long_text: true, form_code: json, order: 7 }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID"  }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 5
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
  action:      { type: varchar, len: 200, nullable: false, comment: "Action", form_display: true, table_display: true, order: 2 }
  req_ip:      { type: varchar, len: 200, comment: "Request IP", form_display: true, table_display: true, form_use_label: true, order: 3 }
  req_at:      { type: datetime, comment: "Request at", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm", form_use_label: true, order: 4 }
  req_data:    { type: text, comment: "Request Data", form_long_text: true, form_use_label: true, order: 5 }
  res_at:      { type: datetime, comment: "Response at", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm", form_use_label: true, order: 6 }
  res_type:    { type: varchar, len: 200, comment: "Response Type", form_display: true, table_display: true, form_use_label: true, order: 7 }
  res_msg:     { type: varchar, len: 500, comment: "Response Message", form_display: true, table_display: true, form_use_label: true, order: 8 }
  res_data:    { type: text, comment: "Request Data", form_long_text: true, form_use_label: true, order: 9 }
  table:       { type: varchar, len: 200, comment: "Table", form_display: true, table_display: true, order: 10 }
  db:          { type: varchar, len: 200, comment: "Database", form_display: true, table_display: true, order: 11 }
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
  table:           { type: varchar, len: 200, comment: "Table", form_display: true, table_display: true, form_size: 6, order: 1 }
  db:              { type: varchar, len: 200, comment: "Database", form_display: true, table_display: true, form_size: 6, order: 2 }
  config:          { type: text, comment: "Config", form_display: true, form_long_text: true, form_code: json, table_display: true, order: 3 }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID", order: 4 }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  created_at:      { type: datetime, comment: "Created at", order: 6 }
  updated_at:      { type: datetime, comment: "Updated at", order: 7 }
  excluded:        { type: boolean, default: false, comment: "Excluded", order: 8 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
```

## CUSTOM_FORM
```yaml
table: custom_form
comment: Custom Form
columns:
  custom_form_id: { type: integer, pk: true, autoincrement: true, comment: "Custom Form ID" }
  table:          { type: varchar, len: 200, comment: "Table", form_display: true, table_display: true, form_size: 6, order: 1 }
  db:             { type: varchar, len: 200, comment: "Database", form_display: true, table_display: true, form_size: 6, order: 2 }
  config:         { type: text, comment: "Config", form_display: true, form_long_text: true, form_code: json, table_display: true, order: 3 }
  app_id:         { type: integer, fk: "app.app_id", comment: "App ID", order: 4 }
  user_id:        { type: integer, fk: "users.user_id", comment: "User ID", order: 5 }
  created_at:     { type: datetime, comment: "Created at", order: 6 }
  updated_at:     { type: datetime, comment: "Updated at", order: 7 }
  excluded:       { type: boolean, default: false, comment: "Excluded", order: 8 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
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
  table:                    { type: varchar, len: 200, nullable: false, comment: "Table", order: 4 }
  db:                       { type: varchar, len: 200, nullable: false, comment: "Database", order: 5 }
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
  table:                  { type: varchar, len: 200, nullable: false, comment: "Table", order: 3 }
  db:                     { type: varchar, len: 200, nullable: false, comment: "Database", order: 4 }
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
  table:               { type: varchar, len: 200, nullable: false, comment: "Table", order: 3 }
  db:                  { type: varchar, len: 200, nullable: false, comment: "Database", order: 4 }
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
  table_org_desc:    { type: varchar, len: 200, nullable: false, comment: "Table Org. Desc", form_display: true, table_display: true, order: 1 }
  table_transl_desc: { type: varchar, len: 200, nullable: false, comment: "Table Transl. Desc", form_display: true, table_display: true, order: 2 }
  table_tooltip:     { type: varchar, len: 500, comment: "Table Tooltip", form_display: true, table_display: true, order: 3 }
  table:             { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
  db:                { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 5 }
  lang:              { type: varchar, len: 5, nullable: false, comment: "Lang", form_display: true, table_display: true, order: 6 }
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
  field_org_desc:      { type: varchar, len: 200, nullable: false, comment: "Field Org. Desc", form_display: true, table_display: true, order: 1 }
  field_transl_desc:   { type: varchar, len: 200, nullable: false, comment: "Field Transl. Desc", form_display: true, table_display: true, order: 2 }
  field_tooltip:       { type: varchar, len: 500, comment: "Field Tooltip", form_display: true, table_display: true, order: 3 }
  field:               { type: varchar, len: 200, nullable: false, comment: "Field", form_display: true, table_display: true, order: 4 }
  table:               { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 5 }
  db:                  { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 6 }
  lang:                { type: varchar, len: 5, nullable: false, comment: "Lang", form_display: true, table_display: true, order: 7 }
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
  db:              { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, form_size: 3, table_display: true, order: 1 }
  table:           { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, form_size: 3, table_display: true, order: 2 }
  field:           { type: varchar, len: 200, nullable: false, comment: "Field", form_display: true, form_size: 3, table_display: true, order: 3 }
  type:            { type: varchar, len: 200, nullable: false, comment: "Type", form_display: true, form_size: 3, table_display: true, order: 4 }
  comment:         { type: varchar, len: 200, comment: "Comment", form_display: true, form_size: 3, table_display: true, order: 5 }
  pk:              { type: boolean, default: false, comment: "Primary Key", form_display: true, form_size: 3, table_display: true, order: 6 }
  autoincrement:   { type: boolean, default: false, comment: "Auto Increment", form_display: true, form_size: 3, table_display: true, order: 7 }
  nullable:        { type: boolean, default: false, comment: "Nullable", form_display: true, form_size: 3, table_display: true, order: 8 }
  computed:        { type: boolean, default: false, comment: "Computed", form_display: true, form_size: 3, table_display: true, order: 9 }
  default:         { type: varchar, len: 200, comment: "Default", form_display: true, form_size: 3, table_display: true, order: 10 }
  fk:              { type: boolean, default: false, comment: "Foreign Key", form_display: true, form_size: 3, table_display: true, order: 11 }
  referred_table:  { type: varchar, len: 200, comment: "Ref. Table", form_display: true, form_size: 3, table_display: true, order: 12 }
  referred_column: { type: varchar, len: 200, comment: "Ref. Column", form_display: true, form_size: 3, table_display: true, order: 13 }
  field_order:     { type: integer, comment: "Field Order", form_display: true, form_size: 3, table_display: true, order: 14 }
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
  cron_id:       { type: integer, pk: true, autoincrement: true, comment: "Cron ID" }
  cron:          { type: varchar, len: 100, unique: false, nullable: false, comment: "Cron Name", form_display: true, table_display: true, form_size: 3, order: 1 }
  cron_desc:     { type: text, nullable: false, comment: "Description", form_display: true, table_display: true, form_size: 9, order: 2 }
  api:           { type: varchar, len: 200, nullable: false, comment: "API Endpoint / Action", form_display: true, table_display: true, form_size: 12, order: 3 }
  db:            { type: varchar, len: 50, comment: "Database (if applicable)", order: 4, form_display: true, table_display: true, form_size: 4 }
  table:         { type: varchar, len: 100, comment: "Table (if applicable)", order: 5, form_display: true, table_display: true, form_size: 4 }
  app_id:        { type: integer, fk: "app.app_id", comment: "Application ID", order: 6, form_display: true, table_display: true, form_size: 4 }
  active:        { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 7 }
  run_only_once: { type: boolean, default: false, comment: "Run Once", form_display: true, table_display: true, form_size: 2, order: 9 }
  last_run:      { type: datetime, comment: "Last Run", form_display: true, table_display: true, form_size: 4, order: 10 }
  user_id:       { type: integer, fk: "users.user_id", comment: "Created/Updated by", order: 8 }
  created_at:    { type: datetime, comment: "Created at", order: 9 }
  updated_at:    { type: datetime, comment: "Updated at", order: 10 }
  excluded:      { type: boolean, default: false, comment: "Excluded", order: 11 }
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
  cron:        { type: varchar, len: 50, nullable: false, comment: "Cron", order: 2, form_display: true, table_display: true, form_size: 3 }
  cron_desc:   { type: varchar, len: 200, nullable: false, comment: "Decription", order: 3, form_display: true, table_display: true, form_size: 9 }
  api:         { type: varchar, len: 200, nullable: false, comment: "API", order: 4, form_display: true, table_display: true, form_size: 12 }
  start_at:    { type: datetime, comment: "Job Start", order: 5, form_display: true, table_display: true, form_size: 4 }
  end_at:      { type: datetime, comment: "Job End", order: 6, form_display: true, table_display: true, form_size: 4 }
  success:     { type: boolean, default: true, comment: "Success", order: 7, form_display: true, table_display: true, form_size: 4 }
  cron_msg:    { type: text, nullable: false, comment: "Message", order: 8, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  app_id:      { type: integer, fk: "app.app_id", comment: "App ID" }
  db:          { type: varchar, len: 200, comment: "Database" }
  table:       { type: varchar, len: 50, comment: "Table" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
table_layout:
  default_order:
    - { field: cron_log_id, order: DESC } 
```

## ACCESS_KEY
```yaml
table: access_key
comment: API / Access Tokens & Keys
columns:
  access_key_id:    { type: integer, pk: true, autoincrement: true, comment: "Access Key ID" }
  access_key_desc:  { type: varchar, len: 200, nullable: false, comment: "Description", form_display: true, table_display: true, form_size: 12, order: 1 }
  access_token:     { type: text, nullable: false, comment: "Token / Secret", form_display: true, table_display: true, form_long_text: true, form_code: txt, order: 2, form_copy: true, table_ellipsis: 90 }
  expires_at:       { type: datetime, comment: "Expires at", form_display: true, table_display: true, form_size: 4, order: 3 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4, order: 4 }
  for_user_id:      { type: integer, fk: "users.user_id", comment: "Assigned to User", form_display: true, table_display: true, form_size: 4, order: 5, form_use_label: true, table_use_label: true }
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
  - component: AccessKey
    label: accessKey
    icon: key
    data: '{ "actions": [ {"type": "btn", "icon": "refresh", "name": "refresh", "class": "btn-sm text-info", "label": "crud.refresh", "action": null}, {"type": "btn", "icon": "save", "name": "save", "class": "btn-sm text-info", "label": "crud.save", "action": null }, {"type": "icon", "icon": "cog-8-tooth", "name": "form_customization", "label": "crud.form_customization", "action": null} ] }'
    size: 6
    intercept_c: true
    intercept_u: true
    main: true
```

## ENV
```yaml
table: env
comment: Envariomental Variables
columns:
  env_id:       { type: integer, pk: true, autoincrement: true, comment: "env ID" }
  env_name:     { type: varchar, len: 200, unique: false, nullable: false, comment: "Env Name", order: 1, form_display: true, table_display: true, form_size: 6 }
  env_value:    { type: text, nullable: false, comment: "Env Value", order: 4, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  on_srv_start: { type: boolean, default: true, comment: "Set On Server Start", order: 3, form_display: true, table_display: true, form_size: 3 }
  active:       { type: boolean, default: true, comment: "Active", order: 3, form_display: true, table_display: true, form_size: 3 }
  user_id:      { type: integer, fk: "users.user_id", comment: "Created BY" }
  created_at:   { type: datetime, comment: "Created at" }
  updated_at:   { type: datetime, comment: "Updated at" }
  excluded:     { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
```

## FLIGHT_CATALOG
```yaml
table: flight_catalog
comment: Arrow Flight Catalogs
columns:
  flight_catalog_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  flight_catalog:      { type: varchar, len: 200, unique: true, nullable: false, comment: "Catalog Name", form_display: true, table_display: true, order: 1, form_size: 10, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!"}
  flight_catalog_desc: { type: text, comment: "Description", form_display: true, table_display: true, order: 3, form_size: 12, form_long_text: true }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_order: 2, form_size: 2 }
  flight_catalog_conf: { type: text, comment: "Config", form_display: false, order: 3, form_size: 12, form_long_text: true, form_code: json }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {flight_catalog_id: 1, flight_catalog: admin, flight_catalog_desc: "Default Arrow Flight catalog", active: true, app_id: 1, user_id: 1, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {flight_schema: true}
table_layout:
  default_order:
    - { field: flight_catalog_id, order: ASC }
  allow_import: false
```

## FLIGHT_SCHEMA
```yaml
table: flight_schema
comment: Arrow Flight Schemas
columns:
  flight_schema_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  flight_schema:       { type: varchar, len: 200, unique: true, nullable: false, comment: "Schema", form_display: true, table_display: true, order: 1, form_size: 5, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  flight_schema_desc:  { type: text, comment: "Description", form_display: true, table_display: true, order: 4, form_text_long: true }
  startup_sql:         { type: text, comment: "Startup SQL", form_display: true, table_display: true, order: 5, form_long_text: true, form_code: sql }
  main_sql:            { type: text, nullable: false, comment: "Main SQL", form_display: true, table_display: true, order: 6, form_long_text: true, form_code: sql }
  table_discover_sql:  { type: text, comment: "Table Discover SQL", form_display: true, order: 7, form_long_text: true, form_code: sql }
  table_scan_tmpl_sql: { type: text, comment: "Table Scan Template SQL", form_display: true, order: 8, form_long_text: true, form_code: sql }
  shutdown_sql:        { type: text, comment: "Shutdown SQL", form_display: true, table_display: true, order: 9, form_long_text: true, form_code: sql }
  flight_schema_conf:  { type: text, comment: "Configuration", form_display: true, order: 10, form_long_text: true, form_code: json}
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_order: 3, form_size: 3 }
  flight_catalog_id:   { type: integer, fk: "flight_catalog.flight_catalog_id", nullable: false, comment: "Flight Catalog", form_display: true, table_display: true, order: 2, form_size: 4 }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {flight_schema_id: 1, flight_catalog_id: 1, flight_schema: adm, flight_schema_desc: "Ex. Arrow Flight Schema using ADMIN app", startup_sql: "INSTALL SQLITE;LOAD SQLITE;", main_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE);USE adm;", shutdown_sql: "USE memory;DETACH adm;", active: false, app_id: 1, user_id: 1, excluded: false}
form_layout:
  tabs_steps: tabs
  tabs_steps_conf:
    - {label: Schema, fields: [flight_schema, flight_catalog, active, flight_schema_desc]}
    - {label: Managment SQLs, fields: [startup_sql, main_sql, table_discover_sql, table_scan_tmpl_sql, shutdown_sql]}
    - {label: Config, fields: [flight_schema_conf]}
  form_in_popup: false
  size: 10
  sub_form_size: 10
  allow_in_subform: {flight_schema_table: true}
table_layout:
  default_order:
    - { field: flight_schema_id, order: ASC }
```

## FLIGHT_SCHEMA_TABLE
```yaml
table: flight_schema_table
comment: Flight Schema Tables
columns:
  flight_schema_table_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  flight_schema_id:         { type: integer, fk: "flight_schema.flight_schema_id", nullable: false, comment: "Flight Schema", form_display: true, table_display: true, order: 1, form_size: 3 }
  table_name:               { type: varchar, len: 200, nullable: false, comment: "Table Name", form_display: true, table_display: true, order: 2, form_size: 3 }
  table_desc:               { type: text, comment: "Description", form_display: true, table_display: true, order: 6, form_long_text: true, form_code: markdown }
  order:                    { type: integer, comment: "Order", form_display: true, table_display: true, order: 3, form_size: 3 }
  active:                   { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 4, form_size: 3 }
  flight_schema_table_conf: { type: text, comment: "Configuration", form_display: true, order: 6, form_size: 12, form_long_text: true, form_code: txt }
  user_id:                  { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                   { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:               { type: datetime, comment: "Created at" }
  updated_at:               { type: datetime, comment: "Updated at" }
  excluded:                 { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  sub_form_size: 8
  allow_in_subform:
    flight_schema_table_field: true
    flight_schema_table_scope: true
table_layout:
  default_order:
    - { field: order, order: ASC }
  allow_import: false
```

## FLIGHT_SCHEMA_TABLE_FIELD
```yaml
table: flight_schema_table_field
comment: Flight Schema Table Fields
columns:
  flight_schema_table_field_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  flight_schema_table_field:      { type: varchar, len: 200, nullable: false, comment: "Field Name", form_display: true, table_display: true, order: 1, form_size: 3 }
  flight_schema_table_field_desc: { type: text, comment: "Field Description", form_display: true, table_display: true, order: 2, form_size: 9 }
  flight_schema_table_id:         { type: integer, fk: "flight_schema_table.flight_schema_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true, order: 3, form_size: 3 }
  flight_schema_id:               { type: integer, fk: "flight_schema.flight_schema_id", comment: "Arrow Flight ID", order: 4, form_display: true, table_display: true, form_size: 4 }
  active:                         { type: boolean, default: true, comment: "Active", order: 5, form_display: true, table_display: true, form_size: 4 }
  user_id:                        { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                         { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:                     { type: datetime, comment: "Created at" }
  updated_at:                     { type: datetime, comment: "Updated at" }
  excluded:                       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 5
```

## FLIGHT_SCHEMA_TABLE_SCOPE
```yaml
table: flight_schema_table_scope
comment: Flight Schema Table Scopes
columns:
  flight_schema_table_scope_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  flight_schema_table_scope:      { type: varchar, len: 200, unique: true, nullable: false, comment: "Scope Name", form_display: true, table_display: true, order: 1, form_size: 4 }
  flight_schema_table_scope_desc: { type: text, comment: "Scope Description", form_display: true, table_display: true, order: 2, form_size: 8 }
  flight_schema_table_scope_sql:  { type: text, nullable: false, comment: "Scope SQL", form_display: true, order: 6, form_long_text: true, form_code: sql }
  flight_schema_table_id:         { type: integer, fk: "flight_schema_table.flight_schema_table_id", comment: "Arrow Flight Table ID", form_display: true, table_display: true, order: 3, form_size: 4 }
  flight_schema_id:               { type: integer, fk: "flight_schema.flight_schema_id", comment: "Arrow Flight ID", form_display: true, table_display: true, order: 4, form_size: 4 }
  active:                         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5, form_size: 4 }
  user_id:                        { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                         { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:                     { type: datetime, comment: "Created at" }
  updated_at:                     { type: datetime, comment: "Updated at" }
  excluded:                       { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## QUACK_SERVER
```yaml
table: quack_server
comment: DuckDB Quack Server
columns:
  quack_server_id: { type: integer, pk: true, autoincrement: true, comment: "ID" }
  quack_name:      { type: varchar, len: 200, unique: true, nullable: false, comment: "Name", form_display: true, table_display: true, order: 1, form_size: 3 }
  quack_desc:      { type: text, comment: "Description", form_display: true, table_display: true, order: 2, form_size: 9 }
  port:            { type: integer, nullable: false, comment: "Port", form_display: true, table_display: true, order: 3, form_size: 2 }
  token:           { type: varchar, len: 200, comment: "Access Token", form_display: true, table_display: false, order: 4, form_size: 4 }
  protocol:        { type: varchar, len: 20, default: quack, comment: "Protocol", form_display: true, table_display: true, order: 5, form_size: 2 }
  startup_sql:     { type: text, comment: "Startup SQL", form_display: true, order: 8, form_long_text: true, form_code: sql }
  main_sql:        { type: text, comment: "Main SQL", form_display: true, order: 9, form_long_text: true, form_code: sql }
  shutdown_sql:    { type: text, comment: "Shutdown SQL", form_display: true, order: 10, form_long_text: true, form_code: sql }
  status:          { type: varchar, len: 50, default: offline, comment: "Status", form_display: true, table_display: true, order: 6, form_size: 2 }
  quack_conf:      { type: text, comment: "Configuration", form_display: true, order: 11, form_long_text: true, form_code: json }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 17}
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {quack_server_id: 1, quack_name: "Quack Admin DB", quack_desc: "Expose ADMIN DB via DuckDB Quack", port: 8779, token: "replace-me", protocol: quack, startup_sql: "INSTALL SQLITE; LOAD SQLITE;", main_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE); USE adm;", shutdown_sql: "USE memory; DETACH adm;", status: offline, active: false, app_id: 1, user_id: 1, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  sub_form_size: 10
table_layout:
  default_order:
    - { field: quack_name, order: ASC }  
  exec_button: 
    - {callApi: true, method: POST, api: quack/start,   tooltip: Start Quack Server,   icon: play,       active: true}
    - {callApi: true, method: POST, api: quack/stop,    tooltip: Stop Quack Server,    icon: stop,       active: true}
    - {callApi: true, method: POST, api: quack/restart, tooltip: Restart Quack Server, icon: arrow-path, active: true}
```

## QUACK_LOGS
```yaml
table: quack_logs
comment: Quack Server Activity Logs
columns:
  quack_log_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  quack_server_id: { type: integer, fk: "quack_server.quack_server_id", nullable: false, comment: "Quack Server", form_display: true, table_display: true, order: 1, form_size: 3 }
  event:           { type: varchar, len: 100, nullable: false, comment: "Event", form_display: true, table_display: true, order: 2, form_size: 3 }
  status:          { type: varchar, len: 50, comment: "Status", form_display: true, table_display: true, order: 3, form_size: 2 }
  port:            { type: integer, comment: "Port", form_display: true, table_display: true, order: 4, form_size: 2 }
  message:         { type: text, comment: "Message", form_display: true, table_display: true, order: 5, form_long_text: true, form_code: txt }
  log_time:        { type: datetime, nullable: false, comment: "Log Time", form_display: true, table_display: true, order: 6, form_size: 3 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 7, form_size: 2 }
  user_id:         { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  allow_in_submenu: {}
  default_order:
    - { field: log_time, order: DESC }
  allow_import: false
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar, len: 200, comment: "Dashboard", form_display: true, table_display: true, form_size: 8, order: 1 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, order: 5 }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 2, order: 2 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 3 }
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
  dashboard_comment:    { type: text, comment: "Comments", order: 1 }
  dashboard:            { type: varchar, len: 200, comment: "Dashboard", order: 2 }
  active:               { type: boolean, default: true, comment: "Active", order: 3 }
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID", order: 4 }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID", order: 5 }
  created_at:           { type: datetime, comment: "Created at", order: 6 }
  updated_at:           { type: datetime, comment: "Updated at", order: 7 }
  excluded:             { type: boolean, default: false, comment: "Excluded", order: 8 }
```

## VALID_REACTION
```yaml
table: valid_reaction
comment: Validation Reaction
columns:
  valid_reaction_id:   { type: integer, pk: true, autoincrement: true, comment: "Validation Reaction ID" }
  valid_reaction:      { type: varchar, len: 20, nullable: false, unique: true, comment: "Validation Reaction", form_display: true, table_display: true, order: 1 }
  valid_reaction_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {valid_reaction_id: 1, valid_reaction: if_empty, valid_reaction_desc: Validation Reaction if Empty, excluded: false}
  - {valid_reaction_id: 2, valid_reaction: if_not_empty, valid_reaction_desc: Validation Reaction if not Empty, excluded: false}
form_layout:
  size: 4
```

## VALIDATION
```yaml
table: validation
comment: Validation Rules
columns:
  validation_id:     { type: integer, pk: true, autoincrement: true, comment: "ID" }
  validation:        { type: varchar, len: 200, nullable: false, comment: "Validation", form_display: true, table_display: true, order: 2, form_size: 9 }
  validation_code:   { type: varchar, len: 200, nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 3 }
  valid_reaction_id: { type: integer, fk: "valid_reaction.valid_reaction_id", comment: "Validation Reaction ID", form_display: true, table_display: true, order: 3, form_size: 2 }
  err_msg:           { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4, form_size: 6 }
  table:             { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 5, form_size: 2 }
  db:                { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 6, form_size: 2 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 7, form_size: 2 }
  create:            { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 8, form_size: 2 }
  read:              { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 9, form_size: 2 }
  update:            { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 10, form_size: 2 }
  delete:            { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 11, form_size: 2 }
  sql:               { type: text, nullable: false, comment: "SQL Rule", form_display: true, order: 12, form_long_text: true, form_code: sql }
  user_id:           { type: integer, fk: "users.user_id", comment: "User ID", order: 10 }
  app_id:            { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at", order: 11 }
  updated_at:        { type: datetime, comment: "Updated at", order: 12 }
  excluded:          { type: boolean, default: false, comment: "Excluded", order: 13 }
data:
  - {validation_id: 1, validation: Validate user Email existance, validation_code: USR01, valid_reaction_id: 2, err_msg: "User {{.email}} already exists!", table: users, db: ADMIN, sql: "select * from users where email = :email", app_id: 1, create: true, user_id: 1}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  allow_in_subform: {validation_logs: true}
  size: 8
```

## VALIDATION_LOGS
```yaml
table: validation_logs
comment: Validation Logs
columns:
  validation_log_id: { type: integer, pk: true, autoincrement: true, comment: "Validation Log ID" }
  validation_id:     { type: integer, fk: "validation.validation_id", comment: "Validation ID", order: 1 }
  validation_code:   { type: varchar, len: 200, comment: "Validation Code", order: 2, form_display: true, table_display: true, form_size: 4 }
  validation:        { type: varchar, len: 200, comment: "Validation Name", order: 3, form_display: true, table_display: true, form_size: 8 }
  table:             { type: varchar, len: 200, comment: "Table", order: 4, form_display: true, table_display: true, form_size: 3 }
  db:                { type: varchar, len: 200, comment: "Database", order: 5, form_display: true, table_display: true, form_size: 3 }
  action:            { type: varchar, len: 10, comment: "Action", order: 6, form_display: true, table_display: true, form_size: 3 }
  success:           { type: boolean, default: true, comment: "Success", order: 9, form_display: true, table_display: true, form_size: 3 }
  log_message:       { type: text, comment: "Log Message", order: 10, form_display: true, table_display: true, form_size: 12, form_long_text: true, form_code: txt }
  user_id:           { type: integer, fk: "users.user_id", comment: "User ID", order: 7 }
  app_id:            { type: integer, fk: "app.app_id", comment: "App ID", order: 8 }
  executed_at:       { type: datetime, comment: "Executed At", order: 11 }
  created_at:        { type: datetime, comment: "Created at", order: 12 }
  updated_at:        { type: datetime, comment: "Updated at", order: 13 }
  excluded:          { type: boolean, default: false, comment: "Excluded", order: 14 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order:
    - { field: validation_log_id, order: DESC }  
```

## ACTION_TYPE
```yaml
table: action_type
comment: CRUD Action Reaction
columns:
  action_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Type ID" }
  action_type:      { type: varchar, len: 20, nullable: false, unique: true, comment: "Type", form_display: true, table_display: true, order: 1 }
  action_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
data:
  - {action_type_id: 1, action_type: ExecuteQuery, action_type_desc: Execute Query, excluded: false}
  - {action_type_id: 2, action_type: SendEmail, action_type_desc: Send Email, excluded: false}
  - {action_type_id: 3, action_type: InternalAPICall, action_type_desc: Internal API Call, excluded: false}
  - {action_type_id: 4, action_type: ExternalAPICall, action_type_desc: External API Call, excluded: false}
form_layout:
  size: 4
```

## CRUD_ACTION
```yaml
table: crud_action
comment: CRUD Action Rules
tooltip: Dispaches some actions after a crud operation
columns:
  crud_action_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  crud_action:       { type: varchar, len: 200, nullable: false, comment: "Action", form_display: true, table_display: true, order: 2, form_size: 9 }
  crud_action_code:  { type: varchar, len: 200, nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 3 }
  action_type_id:    { type: integer, fk: "action_type.action_type_id", comment: "Type ID", form_display: true, table_display: true, order: 3, form_size: 2 }
  err_msg:           { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4, form_size: 6}
  table:             { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 5, form_size: 2 }
  db:                { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 6, form_size: 2 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 7, form_size: 2 }
  create:            { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 8, form_size: 2 }
  read:              { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 9, form_size: 2 }
  update:            { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 10, form_size: 2 }
  delete:            { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 11, form_size: 2 }
  sql:               { type: text, nullable: false, comment: "SQL Rule", form_display: true, table_display: true, order: 12, form_long_text: true, form_code: sql, form_hide_cond: "data?.action_type_id !== 1"}
  email_template:    { type: text, nullable: false, comment: "Email Template", form_display: true, table_display: true, order: 13, form_long_text: true, form_code: html, form_hide_cond: "data?.action_type_id !== 2" }
  email_to:          { type: text, nullable: false, comment: "Email To", tooltip: "Email list separated with semicolon", form_display: true, table_display: true, order: 14, form_long_text: true, form_code: text, form_hide_cond: "data?.action_type_id !== 2" }
  api:               { type: varchar, len: 200, comment: "Call API", form_display: true, table_display: true, order: 15, form_size: 9, form_hide_cond: "data?.action_type_id !== 3" }
  api_id:            { type: integer, comment: "API ID", form_display: true, table_display: true, order: 16, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  api_name:          { type: varchar, len: 200, comment: "API Name", form_display: true, table_display: true, order: 17, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  api_endpoint:      { type: varchar, len: 255, comment: "API Endpoint", form_display: true, table_display: true, order: 18, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  parallel:          { type: boolean, default: false, comment: "Run Parallel", form_display: true, table_display: true, order: 19, form_size: 3 }
  user_id:           { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:            { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  allow_in_subform: {crud_action_logs: true}
  size: 10

```

## CRUD_ACTION_LOGS
```yaml
table: crud_action_logs
comment: CRUD Action Logs
columns:
  crud_action_log_id: { type: integer, pk: true, autoincrement: true, comment: "Log ID" }
  crud_action_id:     { type: integer, fk: "crud_action.crud_action_id", comment: "ID", order: 1 }
  crud_action_code:   { type: varchar, len: 200, comment: "Code", order: 2, form_display: true, table_display: true, form_size: 4 }
  crud_action:        { type: varchar, len: 200, comment: "Name", order: 3, form_display: true, table_display: true, form_size: 8 }
  table:              { type: varchar, len: 200, comment: "Table", order: 4, form_display: true, table_display: true, form_size: 4 }
  db:                 { type: varchar, len: 200, comment: "Database", order: 5, form_display: true, table_display: true, form_size: 4 }
  id:                 { type: integer, comment: "ID", order: 5, form_display: true, table_display: true, form_size: 4 }
  action:             { type: varchar, len: 10, comment: "Action", order: 6, form_display: true, table_display: true, form_size: 4 }
  action_type:        { type: varchar, len: 20, comment: "Action Type", order: 7, form_display: true, table_display: true, form_size: 4 }
  success:            { type: boolean, default: true, comment: "Success", order: 10, form_display: true, table_display: true, form_size: 4 }
  log_message:        { type: text, comment: "Log Message", order: 11, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  user_id:            { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:             { type: integer, fk: "app.app_id", comment: "App ID" }
  executed_at:        { type: datetime, comment: "Executed At" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order:
    - { field: crud_action_log_id, order: DESC } 
```

## PROCESS_TYPE
```yaml
table: process_type
comment: Proccess
columns:
  process_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Proccess ID" }
  process_type:      { type: varchar, len: 20, nullable: false, unique: true, comment: "Proccess", form_display: true, table_display: true, order: 1 }
  process_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
data:
  - {process_type_id: 1, process_type: ExecuteQuery, process_type_desc: Execute Query, excluded: false}
  - {process_type_id: 2, process_type: SendEmail, process_type_desc: Send Email, excluded: false}
form_layout:
  size: 4
```

## BATCH_PROCESS
```yaml
table: batch_process
comment: Processes
tooltip: Processes that runs automatically in the background
columns:
  batch_process_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  batch_process:       { type: varchar, len: 200, nullable: false, comment: "Process", form_display: true, table_display: true, order: 2, form_size: 9 }
  batch_process_code:  { type: varchar, len: 200, nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 2 }
  batch_process_desc:  { type: text, comment: "Description", form_display: true, table_display: true, order: 1, form_size: 2 }
  cron:                { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4 }
  batch_process_order: { type: integer, comment: "Proccess ID", order: 3, form_size: 2 }
  process_type_id:     { type: integer, fk: "process_type.process_type_id", comment: "Proccess ID", order: 3, form_size: 2 }
  err_msg:             { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4 }
  db:                  { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 5 }
  create:              { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 5 }
  read:                { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 6 }
  update:              { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 7 }
  delete:              { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 8 }
  sql:                 { type: text, nullable: false, comment: "SQL Rule", form_display: true, table_display: true, order: 4, form_long_text: true, form_code: sql }
  email_template:      { type: text, nullable: false, comment: "Email Template", form_display: true, table_display: true, order: 4, form_long_text: true, form_code: html }
  email_to:            { type: text, nullable: false, comment: "Email To", tooltip: "Email list separated with semicolon", form_display: true, table_display: true, order: 5, form_long_text: true, form_code: text }
  user_id:             { type: integer, fk: "users.user_id", comment: "User ID", order: 10 }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  created_at:          { type: datetime, comment: "Created at", order: 11 }
  updated_at:          { type: datetime, comment: "Updated at", order: 12 }
  excluded:            { type: boolean, default: false, comment: "Excluded", order: 13 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## BATCH_PROCESS_LOGS
```yaml
table: batch_process_logs
comment: Proccess Logs
columns:
  batch_process_log_id: { type: integer, pk: true, autoincrement: true, comment: "Proccess Log ID" }
  batch_process_id:     { type: integer, fk: "batch_process.batch_process_id", comment: "Proccess ID", order: 1 }
  batch_process_code:   { type: varchar, len: 200, comment: "Proccess Code", order: 2 }
  batch_process:        { type: varchar, len: 200, comment: "Proccess Name", order: 3 }
  db:                   { type: varchar, len: 200, comment: "Database", order: 5 }
  process_type:         { type: varchar, len: 20, comment: "Action Type", order: 7 }
  success:              { type: boolean, default: true, comment: "Success", order: 10 }
  log_message:          { type: text, comment: "Log Message", order: 11 }
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID", order: 8 }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID", order: 9 }
  executed_at:          { type: datetime, comment: "Executed At", order: 12 }
  created_at:           { type: datetime, comment: "Created at", order: 13 }
  updated_at:           { type: datetime, comment: "Updated at", order: 14 }
  excluded:             { type: boolean, default: false, comment: "Excluded", order: 15 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## API_TYPE
```yaml
table: api_type
comment: API Types
columns:
  api_type_id:   { type: integer, pk: true, autoincrement: true, comment: "API Type ID" }
  api_type:      { type: varchar, len: 50, unique: true, nullable: false, comment: "API Type", form_display: true, table_display: true, order: 1 }
  api_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
data:
  - {api_type_id: 1, api_type: REST, api_type_desc: RESTful API, excluded: false}
  - {api_type_id: 2, api_type: SOAP, api_type_desc: SOAP Web Service, excluded: false}
  - {api_type_id: 3, api_type: gRPC, api_type_desc: gRPC Protocol, excluded: false}
  - {api_type_id: 4, api_type: GraphQL, api_type_desc: GraphQL API, excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## HTTP_REQUEST_TYPE
```yaml
table: http_request_type
comment: HTTP Request Types
columns:
  http_request_type_id:   { type: integer, pk: true, autoincrement: true, comment: "HTTP Request Type ID" }
  http_request_type:      { type: varchar, len: 20, unique: true, nullable: false, comment: "HTTP Request Type", form_display: true, table_display: true, order: 1 }
  http_request_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:             { type: datetime, comment: "Created at" }
  updated_at:             { type: datetime, comment: "Updated at" }
  excluded:               { type: boolean, default: false, comment: "Excluded" }
data:
  - {http_request_type_id: 1, http_request_type: GET, http_request_type_desc: "HTTP GET method", excluded: false}
  - {http_request_type_id: 2, http_request_type: POST, http_request_type_desc: "HTTP POST method", excluded: false}
  - {http_request_type_id: 3, http_request_type: PUT, http_request_type_desc: "HTTP PUT method", excluded: false}
  - {http_request_type_id: 4, http_request_type: DELETE, http_request_type_desc: "HTTP DELETE method", excluded: false}
  - {http_request_type_id: 5, http_request_type: PATCH, http_request_type_desc: "HTTP PATCH method", excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## API
```yaml
table: api
comment: API Integrations
columns:
  api_id:                { type: integer, pk: true, autoincrement: true, comment: "API ID" }
  api_name:              { type: varchar, len: 100, nullable: false, comment: "API Name", form_display: true, table_display: true, form_size: 6, order: 1 }
  api_type_id:           { type: integer, fk: "api_type.api_type_id", nullable: false, comment: "API Type", form_display: true, table_display: true, form_size: 3, order: 2 }
  http_request_type_id:  { type: integer, fk: "http_request_type.http_request_type_id", nullable: false, comment: "HTTP Request Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  api_description:       { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 4 }
  endpoint:              { type: varchar, len: 500, nullable: false, comment: "API Endpoint", form_display: true, table_display: true, form_size: 9, order: 5 }
  request_body_template: { type: text, comment: "Request Body Template (Go template)", form_display: true, form_long_text: true, form_code: json, order: 6 }
  num_retries:           { type: integer, default: 3, comment: "Number of Retries", form_display: true, table_display: true, form_size: 3, order: 7 }
  timeout_seconds:       { type: integer, default: 30, comment: "Timeout (seconds)", form_display: true, table_display: true, form_size: 3, order: 8 }
  headers_template:      { type: text, comment: "Headers Template (use @VAR_NAME for environment variables)", form_display: true, form_long_text: true, form_code: json, order: 9 }
  active:                { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:               { type: integer, fk: "users.user_id", comment: "Created by" }
  app_id:                { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  size: 9
  allow_in_subform: {api_call_log: true, api_header: true}
  tabs_steps_conf:
    - {label: API Data, fields: [api_id, api_name, api_type, http_request_type, endpoint, active, num_retries, timeout_seconds, api_description]}
    - {label: Templates, fields: [request_body_template, headers_template]}
  sub_form_size: 9
table_layout:
  default_order:
    - { field: api_id, order: ASC }  
  exec_button: 
    - {callApi: true, method: POST, api: api/run, tooltip: Get Public IP, icon: play, active: true}
data:
  - {api_id: 1, api_name: My public ip, api_type_id: 1, http_request_type_id: 1, api_description: Get my public ip, endpoint: "https://api.ipify.org/", active: true, user_id: 1, app_id: 1, excluded: false}
```

## API_HEADER
```yaml
table: api_header
comment: API Headers
columns:
  api_header_id:   { type: integer, pk: true, autoincrement: true, comment: "Header ID" }
  header_name:     { type: varchar, len: 100, nullable: false, comment: "Header Name", form_display: true, table_display: true, form_size: 5, order: 2 }
  header_value:    { type: text, nullable: false, comment: "Header Value", tootip: "Header Value (supports @VAR_NAME for env variables)", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 4 }
  api_id:          { type: integer, fk: "api.api_id", nullable: false, comment: "API", form_display: true, table_display: true, form_size: 4, order: 1 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 2 }
  user_id:         { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## API_CALL_LOG
```yaml
table: api_call_log
comment: API Call Logs
columns:
  api_call_log_id:      { type: integer, pk: true, autoincrement: true, comment: "Log ID" }
  api_id:               { type: integer, fk: "api.api_id", nullable: false, comment: "API", form_display: true, table_display: true, form_size: 4, order: 1 }
  api_name:             { type: varchar, len: 100, comment: "API Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  request_at:           { type: datetime, nullable: false, comment: "Request DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", order: 3, form_size: 3  }
  response_at:          { type: datetime, comment: "Response DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", order: 4, form_size: 3 }
  request_body:         { type: text, comment: "Request Body", form_display: true, form_long_text: true, form_code: txt, order: 7 }
  response_body:        { type: text, comment: "Response Body", form_display: true, form_long_text: true, form_code: txt, order: 8 }
  response_status:      { type: integer, comment: "Response Status Code", form_display: true, table_display: true, order: 5, form_size: 2 }
  response_message:     { type: text, comment: "Response Message", form_display: true, table_display: true, order: 6, form_long_text: true }
  crud_trggrd_db:       { type: varchar, len: 50, comment: "Crud Triggered DB", form_display: true, order: 9, form_size: 3 }
  crud_trggrd_table:    { type: varchar, len: 50, comment: "Crud Triggered Table", form_display: true, order: 10, form_size: 3 }
  crud_trggrd_pk_field: { type: varchar, len: 50, comment: "Crud Triggered FK Field", form_display: true, order: 11, form_size: 3 }
  crud_trggrd_row_id:   { type: varchar, len: 50, comment: "Crud Triggered Row ID", form_display: true, order: 12, form_size: 3 }
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded"}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order:
    - { field: crud_action_log_id, order: DESC } 
```

# DATA
```yaml
name: DATA
description: DATA Model ADMIN
database: ADMIN
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## DASHBOARD_LOGS
```yaml
table: dashboard
description: Add default Amin User Logs Dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id:   1
  dashboard:      User Logs
  dashboard_desc: User Logs Example
  dashboard_conf: FileContent(examples/user-logs-ex-dashboard.md)
  order:          1
  active:         true
  user_id:        1
  app_id:         appId()
  created_at:     Now()
  updated_at:     Now()
  excluded:       false
```

## EX_CRUD_ACTION_CREATE
```yaml
table: crud_action
description: Add Example of CRUD Actions
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: USER_CREATED_EMAIL
  crud_action: User Created Email
  action_type_id: 2
  err_msg: 'Error sending the email to {{.data.email}}!'
  table: users
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  email_to: '{{.data.email}}'
  email_template: |
    <p>Hi {{.data.first_name}},</p>
    <p>Your account has been created.</p>
    <p><strong>Username:</strong> {{.data.username}}</p>
    <p>If you did not request this, contact support.</p>
```

## EX_CRUD_ACTION_UPDATE
```yaml
table: crud_action
description: Add Example of CRUD Actions
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: USER_UPDATED_EMAIL
  crud_action: User Updated Email
  action_type_id: 2
  err_msg: 'Error sending the email to {{.data.email}}!'
  table: users
  db: ADMIN
  active: true
  create: false
  update: true
  delete: false
  email_to: '{{.data.email}}'
  email_template: |
    <p>Hi {{.data.first_name}},</p>
    <p>Your account has been updated.</p>
    <p>If you did not request this change, contact support.</p>
```

## EX_CRUD_ACTION_DELETE
```yaml
table: crud_action
description: Add Example of CRUD Actions
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: USER_DELETED_EMAIL
  crud_action: User Deleted Email
  action_type_id: 2
  err_msg: 'Error sending the email to {{.data.email}}!'
  table: users
  db: ADMIN
  active: true
  create: false
  update: false
  delete: true
  email_to: '{{.data.email}}'
  email_template: |
    <p>Hi {{.data.first_name}},</p>
    <p>Your account has been deleted.</p>
    <p>If you have questions, contact support.</p>
```

# RUN_ESPECIFC_SQL
```yaml
name: Run domain specifc sql
runs_as: MODEL_SQL
description: Run domain specifc sql
connection: '@DB_DRIVER_NAME:@DB_DSN'
active: false
```

## PROCEDURE_1
```yaml
name: PROCEDURE_1
description: Create store procedure with domain specifc sql
connection: '@DB_DRIVER_NAME:@DB_DSN'
script_sql: pg_procedure
active: true
```

```sql
-- pg_procedure
CREATE OR REPLACE PROCEDURE X AS $$

$$
```
