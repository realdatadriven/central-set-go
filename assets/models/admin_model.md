<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
<!-- markdownlint-disable MD024 -->
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
      - {table: valid_criticity, active: false}
      - {table: validation_data, active: false}
      - {table: action_trigger_action, active: false}
      - user_log
      - custom_table
      - custom_form
      - table
      - table_schema
      - crud_action
      - {table: crud_action_logs, active: false}
      - {table: action_data_type, active: false}
      - {table: action_data, active: false}
      - crud_intercept
      - {table: crud_intercept_logs, active: false}
      - {table: ntercept_data_type, active: false}
      - {table: ntercept_data, active: false}
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
      - {table: api_data_type, active: false}
      - {table: api_data, active: false}
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
  - {role_id: 4, role: anonymous, role_desc: "Anonymous role", excluded: false}
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
  - {user_id: 2, username: anonymous, password: '$2b$12$tfPUUvgU9eHTIvAy/kZo1eW2lrh2rfsX0Qx8YqomZKREoX7sUsbS6', first_name: Anonymous, last_name: User, email: anonymous@domain.com, role_id: 4, lang_id: 1, active: true, alter_pass_nxt_login: false, excluded: false}
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
  app_desc:    { type: text, comment: "Description", form_display: true, form_long_text: true, form_code: markdown, form_rendermd: false, table_display: true, order: 3 }
  version:     { type: varchar, len: 10, nullable: false, comment: "Version", form_display: true, table_display: true, form_size: 3, order: 2 }
  email:       { type: varchar, len: 200, comment: "Email", form_display: true, table_display: true, form_size: 3, form_regex_val: "^\\w+([\\.-]?\\w+)*@\\w+([\\.-]?\\w+)*(\\.\\w{2,3})+$", order: 4 }
  db:          { type: varchar, len: 20, nullable: false, comment: "Database", form_display: true, table_display: true, form_size: 3, order: 5 }
  # conn_string: { type: varchar, len: 200, comment: "Conn String", form_display: true, table_display: true, form_size: 3, order: 5 }
  attach_logo: { type: varchar, len: 200, comment: "Logo", form_display: true, table_display: true, form_size: 3, form_att: true, order: 6 }
  category:    { type: varchar, len: 200, comment: "Category", form_display: true, table_display: true, form_size: 3, order: 6 }
  config:      { type: text, comment: "Config", form_display: true, form_long_text: true, form_code: json }
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
  tabs_steps_conf:
    - {label: App, fields: [app, app_desc, version, email, db, attach_logo, category]}
    - {label: Config, fields: [config]}
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
comment: Menus
columns:
  menu_id:       { type: integer, pk: true, autoincrement: true, comment: "Menu ID" }
  menu:          { type: varchar, len: 20, nullable: false, comment: "Menu", form_display: true, table_display: true, form_size: 3, order: 1 }
  menu_desc:     { type: text, comment: "Description", form_display: true, form_size: 7, table_display: true, order: 2 }
  menu_icon:     { type: varchar, len: 20, comment: "Icon", form_display: true, table_display: true, form_size: 4, order: 4 }
  menu_order:    { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 4, order: 5 }
  menu_config:   { type: text, comment: "Menu Config", form_display: true, form_long_text: true, form_code: json, table_display: true, form_use_label: true, order: 7 }
  app_id:        { type: integer, nullable: false, fk: "app.app_id", comment: "App ID" , form_display: true, table_display: true, form_size: 2, order: 3}
  active:        { type: boolean, default: true, comment: "Active", form_display: true, form_size: 4, table_display: true, order: 6 }
  user_id:       { type: integer, fk: "users.user_id", comment: "User ID" }
  created_at:    { type: datetime, comment: "Created at" }
  updated_at:    { type: datetime, comment: "Updated at" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  allow_in_subform: {menu_table: true}
  tabs_steps_conf: []
table_layout:
  allow_in_submenu: {menu_table: true}
  default_order: [{field: app_id, order: ASC}, {field: menu_order, order: ASC}]
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
  menu_id:        { type: integer, nullable: false, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 1, form_size: 8 }
  table_id:       { type: integer, nullable: false, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true, order: 2, form_size: 4 }
  app_id:         { type: integer, nullable: false, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 3, form_size: 4 }
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
  role_id:     { type: integer, nullable: false, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:      { type: integer, nullable: false, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  access:      { type: boolean, nullable: false, default: true, comment: "Access", form_display: true, table_display: true, order: 3 }
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
  role_id:          { type: integer, nullable: false, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:           { type: integer, nullable: false, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  menu_id:          { type: integer, nullable: false, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 3 }
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
  role_id:                { type: integer, nullable: false, fk: "role.role_id", comment: "Role ID", form_display: true, table_display: true, order: 1 }
  app_id:                 { type: integer, nullable: false, fk: "app.app_id", comment: "App ID", form_display: true, table_display: true, order: 2 }
  menu_id:                { type: integer, nullable: false, fk: "menu.menu_id", comment: "Menu ID", form_display: true, table_display: true, order: 3 }
  table_id:               { type: integer, nullable: false, fk: "table.table_id", comment: "Table ID", form_display: true, table_display: true, order: 4 }
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
  user_id:     { type: integer, nullable: false, fk: "users.user_id", comment: "User ID", form_display: true, table_display: true, order: 1 }
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
  access_token:     { type: text, nullable: false, comment: "Token / Secret", form_display: false, table_display: true, form_long_text: true, form_code: txt, order: 2, form_copy: true, table_copy: true, table_ellipsis: 90 }
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
  - {flight_schema_id: 1, flight_catalog_id: 1, flight_schema: adm, flight_schema_desc: "Ex. Arrow Flight Schema using ADMIN app", startup_sql: "INSTALL SQLITE;LOAD SQLITE;", main_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE);USE adm;", shutdown_sql: "USE memory;DETACH adm;", active: true, app_id: 1, user_id: 1, excluded: false}
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
  - {valid_reaction_id: 1, valid_reaction: trow_err_if_empty, valid_reaction_desc: Validation Reaction if Empty, excluded: false}
  - {valid_reaction_id: 2, valid_reaction: trow_err_if_not_empty, valid_reaction_desc: Validation Reaction if not Empty, excluded: false}
  - {valid_reaction_id: 3, valid_reaction: log_it_as_alert, valid_reaction_desc: Only Log it as Alert or Non-Critical Error, excluded: false}
form_layout:
  size: 4
```

## VALID_CRITICITY
```yaml
table: valid_criticity
comment: Validation Criticity Levels
columns:
  valid_criticity_id: { type: integer, pk: true, autoincrement: true, comment: "Validation Criticity Level ID" }
  criticity_level:        { type: varchar, len: 50, nullable: false, unique: true, comment: "Criticity Level", form_display: true, table_display: true, order: 1 }
  criticity_desc:         { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:             { type: datetime, comment: "Created at" }
  updated_at:             { type: datetime, comment: "Updated at" }
  excluded:               { type: boolean, default: false, comment: "Excluded" }
data:
  - {valid_criticity_id: 1, criticity_level: low, criticity_desc: Low severity, excluded: false}
  - {valid_criticity_id: 2, criticity_level: medium, criticity_desc: Medium severity, excluded: false}
  - {valid_criticity_id: 3, criticity_level: high, criticity_desc: High severity, excluded: false}
  - {valid_criticity_id: 4, criticity_level: critical, criticity_desc: Critical severity, excluded: false}
form_layout:
  size: 4
```

## VALIDATION
```yaml
table: validation
comment: Validation Rules
columns:
  validation_id:      { type: integer, pk: true, autoincrement: true, comment: "ID" }
  validation_code:    { type: varchar, len: 200, nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 3 }
  validation:         { type: varchar, len: 200, nullable: false, comment: "Validation", form_display: true, table_display: true, order: 2, form_size: 9 }
  valid_criticity_id: { type: integer, fk: "valid_criticity.valid_criticity_id", comment: "Validation Criticity Level ID", form_display: true, table_display: true, order: 2, form_size: 2 }
  valid_reaction_id:  { type: integer, fk: "valid_reaction.valid_reaction_id", comment: "Validation Reaction ID", form_display: true, table_display: true, order: 3, form_size: 2 }
  err_msg:            { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: false, order: 4, form_size: 6 }
  table:              { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 5, form_size: 2 }
  db:                 { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 6, form_size: 2 }
  active:             { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 7, form_size: 2 }
  create:             { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 8, form_size: 2 }
  read:               { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 9, form_size: 2 }
  update:             { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 10, form_size: 2 }
  delete:             { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 11, form_size: 2 }
  sql:                { type: text, nullable: false, comment: "SQL Rule", form_display: true, order: 12, form_long_text: true, form_code: sql }
  user_id:            { type: integer, fk: "users.user_id", comment: "User ID", order: 10 }
  app_id:             { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at", order: 11 }
  updated_at:         { type: datetime, comment: "Updated at", order: 12 }
  excluded:           { type: boolean, default: false, comment: "Excluded", order: 13 }
data:
  - {validation_id: 1, validation: Validate user Email existance, validation_code: USR01, valid_criticity_id: 1, valid_reaction_id: 2, err_msg: "User {{.email}} already exists!", table: users, db: ADMIN, sql: "select * from users where email = :email", app_id: 1, create: true, user_id: 1}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  allow_in_subform: {validation_logs: true, validation_data: true}
  size: 8
```

## VALIDATION_DATA
```yaml
table: validation_data
comment: Validation Data
columns:
  validation_data_id:      { type: integer, pk: true, autoincrement: true, comment: "Validation Data ID" }
  validation_data:         { type: varchar, len: 100, nullable: false, comment: "Validation Data", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  validation_data_desc:    { type: text, comment: "Validation Data Desc", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  validation_id:           { type: integer, fk: "validation.validation_id", nullable: false, comment: "Validation", form_display: true, table_display: true, form_size: 4, order: 1 }  
  odata_path:              { type: text, comment: "OData URL", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 9 }
  sigle_row_obj:           { type: boolean, default: false, comment: "Single Row Object", form_display: true, form_size: 4, order: 10 }
  active:                  { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 4 }
  user_id:                 { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:                  { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:              { type: datetime, comment: "Created at" }
  updated_at:              { type: datetime, comment: "Updated at" }
  excluded:                { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## VALIDATION_LOGS
```yaml
table: validation_logs
comment: Validation Logs
columns:
  validation_log_id:  { type: integer, pk: true, autoincrement: true, comment: "Validation Log ID" }
  validation_code:    { type: varchar, len: 200, comment: "Validation Code", order: 2, form_display: true, table_display: true, form_size: 4 }
  validation:         { type: varchar, len: 200, comment: "Validation Name", order: 3, form_display: true, table_display: true, form_size: 8 }
  validation_id:      { type: integer, fk: "validation.validation_id", comment: "Validation ID", order: 1 }
  valid_criticity_id: { type: integer, fk: "valid_criticity.valid_criticity_id", comment: "Validation Criticity Level ID", order: 4, form_display: true, table_display: true, form_size: 3 }
  table:              { type: varchar, len: 200, comment: "Table", order: 5, form_display: true, table_display: true, form_size: 3 }
  db:                 { type: varchar, len: 200, comment: "Database", order: 6, form_display: true, table_display: true, form_size: 3 }
  action:             { type: varchar, len: 10, comment: "Action", order: 7, form_display: true, table_display: true, form_size: 3 }
  executed_at:        { type: datetime, comment: "Executed At", order: 7, form_display: true, table_display: true, form_size: 4}
  success:            { type: boolean, default: true, comment: "Success", order: 9, form_display: true, table_display: true, form_size: 3 }
  log_message:        { type: text, comment: "Log Message", order: 10, form_display: true, table_display: true, form_size: 12, form_long_text: true, form_code: txt }
  user_id:            { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:             { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded", order: 14 }
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
  - {action_type_id: 5, action_type: GeneratePDF, action_type_desc: GeneratePDF, excluded: false}
  - {action_type_id: 6, action_type: RunETLXWorkflow, action_type_desc: Run ETLX Workflow, excluded: false}
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
  user_trigger:      { type: boolean, default: false, comment: "By User", tooltip: "Can be triggered by user", form_display: true, table_display: true, order: 11, form_size: 2 }
  user_trigger_icon: { type: varchar, len: 50, default: false, comment: "By User Icon", tooltip: "User triggered by user ICON", form_display: true, order: 22, form_size: 3, form_hide_cond: "!data?.user_trigger"  }
  sql_condition:     { type: text, comment: "SQL Condition (SELECT BOOLExpression AS cond)", "tooltip": "select '{{.field_to_chrck}}' = 'value to apply to' as cond", form_display: true, order: 11, form_long_text: true, form_code: sql }
  sql:               { type: text, comment: "SQL Rule", form_display: true, order: 12, form_long_text: true, form_code: sql, form_hide_cond: "data?.action_type_id !== 1"}
  email_template:    { type: text, comment: "Email Template", form_display: true, order: 13, form_long_text: true, form_code: html, form_hide_cond: "data?.action_type_id !== 2" }
  email_to:          { type: text, comment: "Email To", tooltip: "Email list separated with semicolon", form_display: true, order: 14, form_size: 6, form_hide_cond: "data?.action_type_id !== 2" }
  email_subject:     { type: text, comment: "Email Subject", form_display: true, order: 14, form_size: 6, form_hide_cond: "data?.action_type_id !== 2" }
  api:               { type: varchar, len: 200, comment: "Call API", form_display: true, order: 15, form_size: 9, form_hide_cond: "data?.action_type_id !== 3" }
  api_id:            { type: integer, comment: "API ID", form_display: true, order: 16, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  api_name:          { type: varchar, len: 200, comment: "API Name", form_display: true, order: 17, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  api_endpoint:      { type: varchar, len: 255, comment: "API Endpoint", form_display: true, order: 18, form_size: 4, form_hide_cond: "data?.action_type_id !== 4" }
  pdf_path:          { type: varchar, len: 200, comment: "PDF Path", form_display: true, order: 19, form_size: 9, form_hide_cond: "data?.action_type_id !== 5" }
  use_latex:         { type: boolean, default: false, comment: "Use Latex", form_display: true, table_display: true, order: 20, form_size: 3, form_hide_cond: "data?.action_type_id !== 5" }
  pdf_tex_template:  { type: text, comment: "PDF LaTex Template", form_display: true, order: 21, form_long_text: true, form_code: latex, form_hide_cond: "data?.action_type_id !== 5 || data?.use_latex === false" }
  pdf_template:      { type: text, comment: "PDF Template", form_display: true, order: 21, form_long_text: true, form_code: html, form_hide_cond: "data?.action_type_id !== 5 || data?.use_latex === true" }
  etlx_md_template:  { type: text, comment: "ETLX Template", form_display: true, order: 21, form_long_text: true, form_code: markdown, form_hide_cond: "data?.action_type_id !== 6" }
  after_sql:         { type: text, comment: "SQL Run After Action", form_display: true, order: 22, form_long_text: true, form_code: sql }
  parallel:          { type: boolean, default: false, comment: "Run Parallel", form_display: true, table_display: true, order: 23, form_size: 3 }
  user_id:           { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:            { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  allow_in_subform: {crud_action_logs: true, action_data: true}
  size: 10
  tabs_steps_conf:
    - {label: Action Def, fields: [crud_action, crud_action_code, action_type, err_msg, table, db, active, create, read, update, delete, user_trigger, user_trigger_icon, parallel]}
    - {label: Config / Templates, fields: [sql_condition, sql, email_template, email_to, email_subject, api, api_name, api_endpoint, pdf_path, use_latex, pdf_template, pdf_tex_template, etlx_md_template, after_sql]}
```

## ACTION_DATA_TYPE
```yaml
table: action_data_type
comment: Action DATA Types
columns:
  action_data_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Action Type ID" }
  action_data_type:      { type: varchar, len: 50, unique: true, nullable: false, comment: "Action DATA Type", form_display: true, table_display: true, order: 1 }
  action_data_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
data:
  - {action_data_type_id: 1, action_data_type: C7SQLQuery, action_data_type_desc: C7 SQL Query, excluded: true}
  - {action_data_type_id: 2, action_data_type: C7Read, action_data_type_desc: C7 Read, excluded: true}
  - {action_data_type_id: 3, action_data_type: C7OData, action_data_type_desc: C7 OData, excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## ACTION_DATA
```yaml
table: action_data
comment: Action Data
columns:
  action_data_id:      { type: integer, pk: true, autoincrement: true, comment: "Action Data ID" }
  action_data:         { type: varchar, len: 100, nullable: false, comment: "Action Data", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  action_data_desc:    { type: text, comment: "Action Data Desc", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  action_data_type_id: { type: integer, fk: "action_data_type.action_data_type_id", nullable: false, comment: "Action Data Type", form_display: true, table_display: true, form_size: 2, order: 3 }
  crud_action_id:      { type: integer, fk: "crud_action.crud_action_id", nullable: false, comment: "Crud Action", form_display: true, table_display: true, form_size: 4, order: 1 }
  action_data_sql:     { type: text, comment: "SQL", form_display: true, form_long_text: true, form_code: sql, form_size: 12, order: 6, form_hide_cond: "data?.action_data_type_id !== 1" }
  read_table:          { type: varchar, len: 50, comment: "Read Table", form_display: true, form_long_text: true, form_code: sql, form_size: 3, order: 7, form_hide_cond: "data?.action_data_type_id !== 2" }
  read_params_json:    { type: text, comment: "Read Params (JSON)", form_display: true, form_long_text: true, form_code: json, form_size: 12, order: 8, form_hide_cond: "data?.action_data_type_id !== 2" }
  odata_path:          { type: text, comment: "OData URL", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 9, form_hide_cond: "data?.action_data_type_id !== 3" }
  sigle_row_obj:       { type: boolean, default: false, comment: "Single Row Object", form_display: true, form_size: 4, order: 10 }
  active:              { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 4 }
  user_id:             { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:              { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## ACTION_TRIGGERS_ACTION
```yaml
table: action_trigger_action
comment: Action Triggers Action
columns:
  action_trigger_action_id:   { type: integer, pk: true, autoincrement: true, comment: "Action Triggers Action ID" }
  action_trigger_action:      { type: varchar, len: 100, nullable: false, comment: "Trigger Code", form_display: true, table_display: true, form_size: 3, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  action_trigger_action_desc: { type: text, comment: "Trigger Desc", form_display: true, form_long_text: true, table_display: true, order: 4 }
  action_trigger_code:        { type: varchar, len: 100, nullable: false, comment: "Code of Action to Trigger", form_display: true, table_display: true, form_size: 3, order: 3, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  crud_action_id:             { type: integer, fk: "crud_action.crud_action_id", nullable: false, comment: "Crud Action", form_display: true, table_display: true, form_size: 4, order: 1 }
  trigger_order:              { type: integer, comment: "Trigger Order", table_display: true, form_display: true, form_size: 2, form_order: 5 }
  active:                     { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 3 }
  user_id:                    { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:                     { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:                 { type: datetime, comment: "Created at" }
  updated_at:                 { type: datetime, comment: "Updated at" }
  excluded:                   { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## CRUD_ACTION_LOGS
```yaml
table: crud_action_logs
comment: CRUD Action Logs
columns:
  crud_action_log_id: { type: integer, pk: true, autoincrement: true, comment: "Log ID" }
  crud_action_id:     { type: integer, fk: "crud_action.crud_action_id", comment: "ID", order: 1 }
  crud_action_code:   { type: varchar, len: 200, comment: "Code", order: 2, form_display: true, table_display: true, form_size: 4 }
  crud_action:        { type: varchar, len: 200, comment: "Name", order: 3, form_display: true, table_display: false, form_size: 8 }
  table:              { type: varchar, len: 200, comment: "Table", order: 4, form_display: true, table_display: true, form_size: 3 }
  db:                 { type: varchar, len: 200, comment: "Database", order: 5, form_display: true, table_display: true, form_size: 3 }
  pk_field:           { type: integer, comment: "PK Filed", order: 5, form_display: true, table_display: false, form_size: 3 }
  id:                 { type: integer, comment: "Row ID", order: 5, form_display: true, table_display: true, form_size: 3 }
  action:             { type: varchar, len: 10, comment: "Action", order: 6, form_display: true, table_display: true, form_size: 3 }
  action_type:        { type: varchar, len: 20, comment: "Action Type", order: 7, form_display: true, table_display: false, form_size: 3 }
  executed_at:        { type: datetime, comment: "Executed At", order: 8, form_display: true, table_display: true, form_size: 4 }
  success:            { type: boolean, default: true, comment: "Success", order: 10, form_display: true, table_display: true, form_size: 2 }
  log_message:        { type: text, comment: "Log Message", order: 11, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  log_data:           { type: text, comment: "Log Data", order: 12, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  user_id:            { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:             { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
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

## INTERCEPT_TYPE
```yaml
table: intercept_type
comment: CRUD intercept Reintercept
columns:
  intercept_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Type ID" }
  intercept_type:      { type: varchar, len: 20, nullable: false, unique: true, comment: "Type", form_display: true, table_display: true, order: 1 }
  intercept_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {intercept_type_id: 1, intercept_type: EncapReadQueryBeforeExec, intercept_type_desc: Encapsulate read query, before execute on condition, excluded: false}
  - {intercept_type_id: 2, intercept_type: ReWriteDMLQueryBeforeExec, intercept_type_desc: Remove escrever Qery Antes de Executar, excluded: false}
form_layout:
  size: 4
```
## CRUD_INTERCEPT
```yaml
table: crud_intercept
comment: CRUD intercept Rules
tooltip: Dispaches some intercepts after a crud operation
columns:
  crud_intercept_id:    { type: integer, pk: true, autoincrement: true, comment: "ID" }
  crud_intercept:       { type: varchar, len: 200, nullable: false, comment: "intercept", form_display: true, table_display: true, order: 2, form_size: 9 }
  crud_intercept_code:  { type: varchar, len: 200, nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 3 }
  intercept_type_id:    { type: integer, fk: "intercept_type.intercept_type_id", comment: "Type ID", form_display: true, table_display: true, order: 3, form_size: 2 }
  err_msg:              { type: varchar, len: 200, nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4, form_size: 6}
  table:                { type: varchar, len: 200, nullable: false, comment: "Table", form_display: true, table_display: true, order: 5, form_size: 2 }
  db:                   { type: varchar, len: 200, nullable: false, comment: "Database", form_display: true, table_display: true, order: 6, form_size: 2 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, order: 7, form_size: 2 }
  create:               { type: boolean, default: false, comment: "Create", form_display: true, table_display: true, order: 8, form_size: 2 }
  read:                 { type: boolean, default: false, comment: "Read", form_display: true, table_display: true, order: 9, form_size: 2 }
  update:               { type: boolean, default: false, comment: "Update", form_display: true, table_display: true, order: 10, form_size: 2 }
  delete:               { type: boolean, default: false, comment: "Delete", form_display: true, table_display: true, order: 11, form_size: 2 }
  sql_condition:        { type: text, comment: "SQL Condition (SELECT BOOLExpression AS cond)", "tooltip": "select '{{.field_to_chrck}}' = 'value to apply to' as cond", form_display: true, order: 11, form_long_text: true, form_code: sql }
  sql:                  { type: text, comment: "SQL Template Encapsulate", form_display: true, order: 12, form_long_text: true, form_code: sql, form_hide_cond: "data?.intercept_type_id !== 1"}
  rewrite_exec_list:    { type: text, comment: "Rewrite Exec List (JSON Array with)", form_display: true, order: 13, form_long_text: true, form_code: json, form_hide_cond: "data?.intercept_type_id !== 2" }
  after_sql:            { type: text, comment: "SQL Run After intercept", form_display: true, order: 22, form_long_text: true, form_code: sql }
  parallel:             { type: boolean, default: false, comment: "Run Parallel", form_display: true, table_display: true, order: 23, form_size: 3 }
  user_id:              { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:               { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  allow_in_subform: {crud_intercept_logs: true, intercept_data: true}
  size: 10
  tabs_steps_conf:
    - {label: Intercept Def, fields: [crud_intercept, crud_intercept_code, intercept_type, err_msg, table, db, active, create, read, update, delete, user_trigger, user_trigger_icon, parallel]}
    - {label: Config / Templates, fields: [sql_condition, sql, field_list, after_sql]}
```

## INTERCEPT_DATA_TYPE
```yaml
table: intercept_data_type
comment: Intercept DATA Types
columns:
  intercept_data_type_id:   { type: integer, pk: true, autoincrement: true, comment: "intercept Type ID" }
  intercept_data_type:      { type: varchar, len: 50, unique: true, nullable: false, comment: "intercept DATA Type", form_display: true, table_display: true, order: 1 }
  intercept_data_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:               { type: datetime, comment: "Created at" }
  updated_at:               { type: datetime, comment: "Updated at" }
  excluded:                 { type: boolean, default: false, comment: "Excluded" }
data:
  - {intercept_data_type_id: 1, intercept_data_type: C7SQLQuery, intercept_data_type_desc: C7 SQL Query, excluded: true}
  - {intercept_data_type_id: 2, intercept_data_type: C7Read, intercept_data_type_desc: C7 Read, excluded: true}
  - {intercept_data_type_id: 3, intercept_data_type: C7OData, intercept_data_type_desc: C7 OData, excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## INTERCEPT_DATA
```yaml
table: intercept_data
comment: Interception Data
columns:
  intercept_data_id:      { type: integer, pk: true, autoincrement: true, comment: "Intercept Data ID" }
  intercept_data:         { type: varchar, len: 100, nullable: false, comment: "Intercept Data", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  intercept_data_desc:    { type: text, comment: "intercept Data Desc", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  intercept_data_type_id: { type: integer, fk: "intercept_data_type.intercept_data_type_id", nullable: false, comment: "intercept Data Type", form_display: true, table_display: true, form_size: 2, order: 3 }
  crud_intercept_id:      { type: integer, fk: "crud_intercept.crud_intercept_id", nullable: false, comment: "Crud intercept", form_display: true, table_display: true, form_size: 4, order: 1 }
  intercept_data_sql:     { type: text, comment: "SQL", form_display: true, form_long_text: true, form_code: sql, form_size: 12, order: 6, form_hide_cond: "data?.intercept_data_type_id !== 1" }
  odata_path:             { type: text, comment: "OData URL", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 9, form_hide_cond: "data?.intercept_data_type_id !== 3" }
  sigle_row_obj:          { type: boolean, default: false, comment: "Single Row Object", form_display: true, form_size: 4, order: 10 }
  active:                 { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 4 }
  user_id:                { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:                 { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:             { type: datetime, comment: "Created at" }
  updated_at:             { type: datetime, comment: "Updated at" }
  excluded:               { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```

## CRUD_INTERCEPT_LOGS
```yaml
table: crud_intercept_logs
comment: CRUD Interception Logs
columns:
  crud_intercept_log_id: { type: integer, pk: true, autoincrement: true, comment: "Log ID" }
  crud_intercept_id:     { type: integer, fk: "crud_intercept.crud_intercept_id", comment: "ID", order: 1 }
  crud_intercept_code:   { type: varchar, len: 200, comment: "Code", order: 2, form_display: true, table_display: true, form_size: 4 }
  crud_intercept:        { type: varchar, len: 200, comment: "Name", order: 3, form_display: true, table_display: true, form_size: 8 }
  table:                 { type: varchar, len: 200, comment: "Table", order: 4, form_display: true, table_display: true, form_size: 3 }
  db:                    { type: varchar, len: 200, comment: "Database", order: 5, form_display: true, table_display: true, form_size: 3 }
  pk_field:              { type: integer, comment: "PK Filed", order: 5, form_display: true, table_display: true, form_size: 3 }
  id:                    { type: integer, comment: "Row ID", order: 5, form_display: true, table_display: true, form_size: 3 }
  intercept:             { type: varchar, len: 10, comment: "intercept", order: 6, form_display: true, table_display: true, form_size: 3 }
  intercept_type:        { type: varchar, len: 20, comment: "intercept Type", order: 7, form_display: true, table_display: true, form_size: 3 }
  executed_at:           { type: datetime, comment: "Executed At", order: 8, form_display: true, table_display: true, form_size: 4 }
  success:               { type: boolean, default: true, comment: "Success", order: 10, form_display: true, table_display: true, form_size: 2 }
  log_message:           { type: text, comment: "Log Message", order: 11, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  log_data:              { type: text, comment: "Log Data", order: 12, form_display: true, table_display: true, form_long_text: true, form_code: txt }
  user_id:               { type: integer, fk: "users.user_id", comment: "User ID" }
  app_id:                { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order:
    - { field: crud_intercept_log_id, order: DESC } 
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
  api_id:                { type: integer, pk: true, autoincrement: true, start: 1000, comment: "API ID" }
  api_name:              { type: varchar, len: 100, nullable: false, comment: "API Name", form_display: true, table_display: true, form_size: 6, order: 1 }
  api_type_id:           { type: integer, fk: "api_type.api_type_id", nullable: false, comment: "API Type", form_display: true, table_display: true, form_size: 3, order: 2 }
  http_request_type_id:  { type: integer, fk: "http_request_type.http_request_type_id", nullable: false, comment: "HTTP Request Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  api_description:       { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 4 }
  endpoint:              { type: varchar, len: 500, nullable: false, comment: "API Endpoint", form_display: true, table_display: true, form_size: 9, order: 5, table_ellipsis: 50 }
  request_body_template: { type: text, comment: "Request Body Template (Go template)", form_display: true, form_long_text: true, form_code: json, order: 6 }
  num_retries:           { type: integer, default: 3, comment: "Number of Retries", form_display: true, table_display: true, form_size: 3, order: 7 }
  timeout_seconds:       { type: integer, default: 30, comment: "Timeout (seconds)", form_display: true, table_display: true, form_size: 3, order: 8 }
  headers_template:      { type: text, comment: "Headers Template (use @VAR_NAME for environment variables)", form_display: true, form_long_text: true, form_code: json, order: 9 }
  after_sql:             { type: text, comment: "SQL Run After", form_display: true, form_long_text: true, form_code: sql, order: 11 }
  active:                { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 11 }
  user_id:               { type: integer, fk: "users.user_id", comment: "Created by" }
  app_id:                { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  size: 9
  allow_in_subform: {api_call_log: true, api_header: true, api_data: true }
  tabs_steps_conf:
    - {label: API Data, fields: [api_id, api_name, api_type, http_request_type, endpoint, active, num_retries, timeout_seconds, api_description]}
    - {label: Templates, fields: [request_body_template, headers_template, after_sql]}
  sub_form_size: 9
table_layout:
  default_order: [{ field: api_id, order: ASC }]
  exec_button: 
    - {callApi: true, method: POST, api: api/run, tooltip: Get Public IP, icon: play, active: true}
data:
  - {api_id: 1, api_name: My public ip, api_type_id: 1, http_request_type_id: 1, api_description: Get my public ip, endpoint: "https://api.ipify.org", active: true, user_id: 1, app_id: 1, excluded: false}
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

## API_DATA_TYPE
```yaml
table: api_data_type
comment: API DATA Types
columns:
  api_data_type_id:   { type: integer, pk: true, autoincrement: true, comment: "API Type ID" }
  api_data_type:      { type: varchar, len: 50, unique: true, nullable: false, comment: "API DATA Type", form_display: true, table_display: true, order: 1 }
  api_data_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
data:
  - {api_data_type_id: 1, api_data_type: C7SQLQuery, api_data_type_desc: C7 SQL Query, excluded: true}
  - {api_data_type_id: 2, api_data_type: C7Read, api_data_type_desc: C7 Read, excluded: true}
  - {api_data_type_id: 3, api_data_type: C7OData, api_data_type_desc: C7 OData, excluded: false}
form_layout:
  form_in_popup: true
  size: 4
```

## API_DATA
```yaml
table: api_data
comment: API Data
columns:
  api_data_id:      { type: integer, pk: true, autoincrement: true, comment: "API Data ID" }
  api_data:         { type: varchar, len: 100, nullable: false, comment: "API Data", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not beging by number, no space or special character!" }
  api_data_desc:    { type: text, comment: "API Data Desc", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  api_data_type_id: { type: integer, fk: "api_data_type.api_data_type_id", nullable: false, comment: "API Data Type", form_display: true, table_display: true, form_size: 2, order: 3 }
  api_id:           { type: integer, fk: "api.api_id", nullable: false, comment: "API", form_display: true, table_display: true, form_size: 4, order: 1 }
  action_data_sql:  { type: text, comment: "SQL", form_display: true, form_long_text: true, form_code: sql, form_size: 12, order: 6, form_hide_cond: "data?.action_data_type_id !== 1" }
  read_table:       { type: varchar, len: 50, comment: "Read Table", form_display: true, form_long_text: true, form_code: sql, form_size: 3, order: 7, form_hide_cond: "data?.action_data_type_id !== 2" }
  read_params_json: { type: text, comment: "Read Params (JSON)", form_display: true, form_long_text: true, form_code: json, form_size: 12, order: 8, form_hide_cond: "data?.action_data_type_id !== 2" }
  odata_path:       { type: text, comment: "OData URL", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 9, form_hide_cond: "data?.action_data_type_id !== 3" }
  sigle_row_obj:    { type: boolean, default: false, comment: "Single Row Object", form_display: true, form_size: 4, order: 10 }
  active:           { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 4 }
  user_id:          { type: integer, fk: "users.user_id", comment: "Created by"  }
  app_id:           { type: integer, fk: "app.app_id", comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
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
  size: 8
table_layout:
  default_order: [{ field: crud_action_log_id, order: DESC }] 
```

# DATA
```yaml
name: DATA
description: DATA Model ADMIN
database: ADMIN
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## APP_SET_CONFIG_CSS_VARS
```yaml
table: app
description: Update apps, add css vars for https://v4.daisyui.com/docs/colors/
cond: 'WHERE app_id = :app_id AND excluded = false'
data:
  app_id: 1
  config: |
    {
      "css_vars": {
        "--p": "55% 0.3 240", 
        "--pc": "98% 0.01 240", 
        "--s": "70% 0.25 200", 
        "--sc": "98% 0.01 200", 
        "--a": "65% 0.25 160", 
        "--ac": "98% 0.01 160"
      }
    }
  updated_at: Now()
  excluded:   false
```

## APP_SET_CONFIG_CSS_VARS_v5
```yaml
table: apps
description: Update apps, add css vars for https://daisyui.com/docs/themes/
cond: 'WHERE app_id = :app_id'
active: false
data:
  app_id: 1
  config: |
    '{
      "css_vars": {
        "--color-base-100": "oklch(98% 0.02 240)", 
        "--color-base-200": "oklch(95% 0.03 240)", 
        "--color-base-300": "oklch(92% 0.04 240)", 
        "--color-base-content": "oklch(20% 0.05 240)", 
        "--color-primary": "oklch(55% 0.3 240)", 
        "--color-primary-content": "oklch(98% 0.01 240)", 
        "--color-secondary": "oklch(70% 0.25 200)", 
        "--color-secondary-content": "oklch(98% 0.01 200)", 
        "--color-accent": "oklch(65% 0.25 160)", 
        "--color-accent-content": "oklch(98% 0.01 160)", 
        "--color-neutral": "oklch(50% 0.05 240)", 
        "--color-neutral-content": "oklch(98% 0.01 240)", 
        "--color-info": "oklch(70% 0.2 220)", 
        "--color-info-content": "oklch(98% 0.01 220)", 
        "--color-success": "oklch(65% 0.25 140)", 
        "--color-success-content": "oklch(98% 0.01 140)", 
        "--color-warning": "oklch(80% 0.25 80)", 
        "--color-warning-content": "oklch(20% 0.05 80)", 
        "--color-error": "oklch(65% 0.3 30)", 
        "--color-error-content": "oklch(98% 0.01 30)"
      }
    }'
  updated_at: Now()
  excluded:   false
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

## APP_MENU_UNIQUENESS
```yaml
table: validation
description: Ensure menu uniqueness for a specific app
cond: 'WHERE validation_code = :validation_code'
data:
  validation:         Ensure menu uniqueness for a specific app
  validation_code:    APP_MENU_UNIQUENESS
  valid_criticity_id: 4
  valid_reaction_id:  2
  err_msg:            'Menu "{{.menu}}" already exists for app "{{.appdata.app}}" (app_id = {{.app_id}})!'
  table:              menu
  db:                 ADMIN
  sql:                "select * from menu where menu = :menu and app_id = :app_id and excluded = false"
  app_id:             1
  create:             true
  user_id:            1
  created_at:         Now()
  updated_at:         Now()
  excluded:           false
  children:
    table: validation_data
    cond: 'WHERE validation_data = :validation_data'
    data:
      validation_data:       appdata
      validation_data_desc:  Associated App Data
      validation_id:         validation_id()
      odata_path:            "ADMIN/app?$filter=app_id eq {{.app_id}}"
      sigle_row_obj:         true
      active:                true
      user_id:               1
      app_id:                appId()
      created_at:            Now()
      updated_at:            Now()
      excluded:              false
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
  err_msg: 'Error sending the email to {{.email}}!'
  table: users
  db: ADMIN
  active: false
  create: true
  update: false
  delete: false
  parallel: true
  email_subject: C7 User Created
  email_to: '{{.email}}'
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
  err_msg: 'Error sending the email to {{.email}}!'
  table: users
  db: ADMIN
  active: false
  create: false
  update: true
  delete: false
  parallel: true
  email_subject: C7 User Updated
  email_to: '{{.email}}'
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
  err_msg: 'Error sending the email to {{.email}}!'
  table: users
  db: ADMIN
  active: false
  create: false
  update: false
  delete: true
  parallel: true
  email_subject: C7 User Deleted
  email_to: '{{.email}}'
  email_template: |
    <p>Hi {{.data.first_name}},</p>
    <p>Your account has been deleted.</p>
    <p>If you have questions, contact support.</p>
```

## EX_CRUD_ACTION_CREATE_LOG
```yaml
table: crud_action
description: Add Example of CRUD Action that logs user creation
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: USER_CREATE_LOG
  crud_action: User Create Log
  action_type_id: 1
  err_msg: 'Error inserting user log for {{.user_id}} {{.email}}'
  table: users
  db: ADMIN
  active: false
  create: true
  update: false
  delete: false
  parallel: true
  sql: |
    INSERT INTO "user_log" ("user_id", "action", "table", "db", "row_id", "created_at", "updated_at", "excluded")
    VALUES (:user_id, 'crud/create from actions', 'users', 'ADMIN', :user_id, :created_at, :updated_at, FALSE)
```

## EX_CRUD_ACTION_API_CALL
```yaml
table: crud_action
description: Add Example of CRUD Action API CALL Example
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: CRUD_ACTION_API_CALL
  crud_action: User Create Log
  action_type_id: 3
  err_msg: 'Error calling the api for {{.user_id}} {{.email}}'
  table: users
  db: ADMIN
  active: true
  create: true
  update: true
  delete: true
  parallel: true
  api: 'run/process_user/{{.action}}/{{.user_id}}'
```

## EX_CRUD_ACTION_EXTERNAL_API_CALL
```yaml
table: crud_action
description: Add Example of CRUD Action External API CALL Example
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: CRUD_ACTION_EXT_API_CALL
  crud_action: User External API Call Example
  action_type_id: 4
  err_msg: 'Error calling the external api for {{.user_id}} {{.email}} to get the external ip'
  table: users
  db: ADMIN
  active: true
  user_trigger: true
  user_trigger_icon: play
  parallel: true
  api_endpoint: 'https://api.ipify.org'
```

## EX_CRUD_ACTION_GEN_PDF_EX
```yaml
table: crud_action
description: Add Example of CRUD Action Generate PDF Example
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code:  EX_CRUD_ACTION_GEN_PDF_EX
  crud_action:       Generate PDF Example
  action_type_id:    5
  err_msg:           'Error generating PDF example {{.user_id}} {{.email}} to get the external ip'
  table:             users
  db:                ADMIN
  active:            true
  user_trigger:      true
  user_trigger_icon: document-arrow-down
  parallel:          false
  pdf_path:          'static/uploads/tmp/user_{{.username}}.pdf'
  pdf_template:      FileContent(examples/exemple_action_template.html)
  after_sql:         "update users set attach_profile_pic = 'tmp/user_{{.username}}.pdf' where user_id = {{.user_id}}"
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      action_data:         main_role
      action_data_desc:    Main Role
      action_data_type_id: 3
      crud_action_id:      crud_action_id()
      odata_path:          "ADMIN/role?$filter=role_id eq {{.role_id}}"
      sigle_row_obj:       true
      active:              true
      user_id:             1
      app_id:              appId()
      created_at:          Now()
      updated_at:          Now()
      excluded:            false
```

## EX_ACTION_GEN_PDF_LATEX_EX
```yaml
table: crud_action
description: Add Example of CRUD Action Generate PDF Example From Latex Tmpl
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code:  EX_ACTION_GEN_PDF_LATEX_EX
  crud_action:       PDF Example From TEX
  action_type_id:    5
  err_msg:           'Error generating PDF From Latex example {{.user_id}} {{.email}} to get the external ip'
  table:             users
  db:                ADMIN
  active:            true
  user_trigger:      true
  user_trigger_icon: document-text
  parallel:          false
  pdf_path:          'static/uploads/tmp/user_{{.username}}.pdf'
  use_latex:         true
  pdf_tex_template:  FileContent(examples/static.invoice.ex.tex)
  after_sql:         "update users set attach_profile_pic = 'tmp/user_{{.username}}.pdf' where user_id = {{.user_id}}"
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      action_data:         main_role2
      action_data_desc:    Main Role 2
      action_data_type_id: 3
      crud_action_id:      crud_action_id()
      odata_path:          "ADMIN/role?$filter=role_id eq {{.role_id}}"
      sigle_row_obj:       true
      active:              true
```

## LOG_IP_GEODATA
```yaml
table: crud_action
description: Example of CRUD Actions Get User log IP GeoData
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: LOG_IP_GEODATA
  crud_action: Get User log IP GeoData
  action_type_id: 4
  err_msg: 'Error calling the external api for {{.ip}} to get the ip geodata'
  table: user_log
  db: ADMIN
  active: true
  user_trigger: true
  user_trigger_icon: globe-alt
  parallel: true
  api_name: IP_GEODATA
  children: # TO BE ABLE TO TRIGGER A NEW ACTION TO PARSE THE RESPONSE api_response THIS WAY ON SUCCESS 
    table: action_trigger_action
    cond: 'WHERE action_trigger_action = :action_trigger_action'
    data:
      crud_action_id:             crud_action_id()
      action_trigger_action:      CALL_HANDLE_IP_GEODATA
      action_trigger_action_desc: On Success it call the HANDLE_IP_GEODATA Action
      action_trigger_code:        HANDLE_IP_GEODATA
      trigger_order:              1
      active:                     true
```

## API_IP_GEODATA
```yaml
table: api
description: Get IP GeoData
cond: 'WHERE api_name = :api_name'
data:
  api_name:              IP_GEODATA
  api_type_id:           1
  http_request_type_id:  1
  api_description:       Get IP Geo Data
  endpoint:              'https://ipinfo.io/{{.data.req_ip}}/json' # api_id api_name db table_pk_field user api_endpoint data action table row_id
  active:                true
```

## HANDLE_IP_GEODATA
```yaml
table: crud_action
description: Example o CRUD Actions Parsing API response
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: HANDLE_IP_GEODATA
  crud_action: 'Handle the body "api_response" in the scope'
  action_type_id: 6
  err_msg: Error hadle the reponse got from API_IP_GEODATA {{.ip}}!
  table: user_log
  db: ADMIN
  active: true
  parallel: true
  etlx_md_template: FileContent(examples/HANDLE_IP_GEODATA.md)
```

# ROLE_ACCESS
```yaml
name: ROLE_ACCESS
description: Role Model Example, Create Role anonymous and gives access to ADMIN/Arrow Flight ...
database: ADMIN #RLA may need to access records on other DBs
runs_as: ROLE
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: false
```

## ANONYMOUS
```yaml
name: anonymous
description: Anonymous Role
access:
  - ADMIN:
    - Arrow Flight:
      - {table: flight_catalog, read: true, rla: [{flight_catalog: admin, read: true, share: true}]}
      - flight_schema
      - flight_schema_table
active: true
```

# ROLE_USERS
```yaml
name: ROLE_USERS
description: Give user access to a role
runs_as: ROLE_USERS
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: false
```

## ANONYMOUS
```yaml
name: anonymous
description: Anonymous Role
users: [root, admin@domain.com, anonymous.user]
active: true
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
