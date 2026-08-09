<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
<!-- markdownlint-disable MD024 -->

# SURVEY_MODEL
```yaml
name: SURVEY
description: Custom Survey Builder Model
runs_as: MODEL
database: SURVEY
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
  Survey Builder:
    menu_icon: clipboard-document-list
    menu_order: 2
    active: true
    menu_config: '{"label": "survey","tooltip": "description","load_items": {"table": "survey","tables": ["survey"]}}'
    tables:
      - survey
      - {table: page, active: false}
      - {table: question, active: false}
      - {table: choice, active: false}
      - {table: condition, active: false}
      - {table: question_type, active: false}
  Survey JSON:
    menu_icon: question-mark-circle
    menu_order: 3
    active: true
    #menu_config: '{"label": "survey","tooltip": "description","load_items": {"table": "survey_json","tables": ["survey_json"]}}'
    tables:
      - survey_json
  Survey Response:
    menu_icon: circle-stack
    menu_order: 4
    active: true
    tables:
      - response
```

## QUESTION_TYPE
```yaml
table: question_type
comment: Question Types
columns:
  question_type_id: { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  type_name:        { type: varchar(50),  unique: true, nullable: false, comment: "Type Name", form_display: true, table_display: true, form_size: 4 }
  label:            { type: varchar(100), comment: "Display Label",     form_display: true, table_display: true, form_size: 4 }
  icon:             { type: varchar(50),  comment: "Icon",              form_display: true, form_size: 4 }
  settings_schema:  { type: text,         comment: "Allowed settings (JSON schema) for this input type", form_display: true, form_long_text: true, form_code: json }
  active:           { type: boolean,      default: true, comment: "Active", form_display: true, table_display: true, form_size: 3 }
  user_id:          { type: integer,      comment: "Created by" }
  created_at:       { type: datetime,     comment: "Created at" }
  updated_at:       { type: datetime,     comment: "Updated at" }
  excluded:         { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 9
table_layout:
  default_order: [{field: type_name, order: ASC}]
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

## SURVEY
```yaml
table: survey
comment: Surveys
columns:
  survey_id:   { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  title:       { type: varchar(255), nullable: false, comment: "Title",       form_display: true, table_display: true, form_size: 6 }
  description: { type: text,         comment: "Description",                 form_display: true, table_display: true, form_long_text: true, form_size: 9 }
  status:      { type: varchar(50),  default: draft, comment: "Status: draft | published | archived", form_display: true, table_display: true, form_size: 3 }
  user_id:     { type: integer,      comment: "Created by" }
  created_at:  { type: datetime,     comment: "Created at",   table_display: true }
  updated_at:  { type: datetime,     comment: "Updated at" }
  excluded:    { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {page: true }
table_layout:
  default_order: [{field: updated_at, order: DESC}]
table_extra_options:
  - {size: 12, component: SurveyPreview, label: preview, icon: eye, pop_up: true}
```

## PAGE
```yaml
table: page
comment: Survey Pages
columns:
  page_id:     { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  survey_id:   { type: integer,      fk: "survey.survey_id", nullable: false, comment: "Survey", form_display: true, table_display: true, form_size: 6 }
  page_order:  { type: integer,      nullable: false, default: 1, comment: "Order", form_display: true, table_display: true, form_size: 3 }
  title:       { type: varchar(255), comment: "Title",       form_display: true, table_display: true, form_size: 6 }
  description: { type: text,         comment: "Description / instructions", form_display: true, form_long_text: true, form_size: 9 }
  user_id:     { type: integer,      comment: "Created by" }
  created_at:  { type: datetime,     comment: "Created at" }
  updated_at:  { type: datetime,     comment: "Updated at" }
  excluded:    { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {page: true }
table_layout:
  default_order: [{field: survey_id, order: ASC}, {field: page_order, order: ASC}]
```

## QUESTION
```yaml
table: question
comment: Survey Questions
columns:
  question_id:      { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  page_id:          { type: integer,      fk: "page.page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 4 }
  question_order:   { type: integer,      nullable: false, default: 1, comment: "Order", form_display: true, table_display: true, form_size: 2 }
  text:             { type: text,         nullable: false, comment: "Question text", form_display: true, table_display: true, form_long_text: true, form_size: 9 }
  question_type_id: { type: integer,      fk: "question_type.question_type_id", nullable: false, comment: "Input Type", form_display: true, table_display: true, form_size: 3 }
  is_required:      { type: boolean,      default: false, comment: "Required", form_display: true, table_display: true, form_size: 2 }
  placeholder:      { type: varchar(255), comment: "Placeholder", form_display: true, form_size: 5 }
  default_value:    { type: text,         comment: "Default value", form_display: true, form_size: 5 }
  settings_json:    { type: text,         comment: "Type-specific settings (min/max, rows, regex, choices flags, ...)", form_display: true, form_long_text: true, form_code: json }
  user_id:          { type: integer,      comment: "Created by" }
  created_at:       { type: datetime,     comment: "Created at" }
  updated_at:       { type: datetime,     comment: "Updated at" }
  excluded:         { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
  allow_in_subform: {page: true }
table_layout:
  default_order: [{field: page_id, order: ASC}, {field: question_order, order: ASC}]
```

## CHOICE
```yaml
table: choice
comment: Question Choices (radio / checkbox / dropdown)
columns:
  choice_id:    { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  question_id:  { type: integer,      fk: "question.question_id", nullable: false, comment: "Question", form_display: true, table_display: true, form_size: 5 }
  value:        { type: varchar(255), nullable: false, comment: "Internal value", form_display: true, table_display: true, form_size: 3 }
  text:         { type: varchar(255), nullable: false, comment: "Display text",  form_display: true, table_display: true, form_size: 4 }
  choice_order: { type: integer,      nullable: false, default: 1, comment: "Order", form_display: true, table_display: true, form_size: 2 }
  user_id:      { type: integer,      comment: "Created by" }
  created_at:   { type: datetime,     comment: "Created at" }
  updated_at:   { type: datetime,     comment: "Updated at" }
  excluded:     { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
  sub_form_size: 9
table_layout:
  default_order: [{field: question_id, order: ASC}, {field: choice_order, order: ASC}]
```

## CONDITION
```yaml
table: condition
comment: Conditional Logic (show / hide / enable / disable / set value)
columns:
  condition_id:        { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  target_type:         { type: varchar(50),  nullable: false, comment: "Target: question | page",  form_display: true, table_display: true, form_size: 3 }
  target_id:           { type: integer,      nullable: false, comment: "Target ID (question_id or page_id, per target_type)", form_display: true, table_display: true, form_size: 3 }
  source_question_id:  { type: integer,      fk: "question.question_id", nullable: false, comment: "Source question", form_display: true, table_display: true, form_size: 6 }
  operator:            { type: varchar(50),  nullable: false, comment: "= | != | > | < | contains", form_display: true, table_display: true, form_size: 3 }
  value:               { type: text,         nullable: false, comment: "Value to compare against", form_display: true, table_display: true, form_size: 4 }
  action:              { type: varchar(50),  nullable: false, comment: "show | hide | enable | disable | set_value | set_required", form_display: true, table_display: true, form_size: 3 }
  action_value:        { type: text,         comment: "Value to set when action = set_value", form_display: true, form_size: 8 }
  logic_json:          { type: text,         comment: "Optional nested AND/OR rule tree (overrides operator/value/source_question_id when set)", form_display: true, form_long_text: true, form_code: json }
  user_id:             { type: integer,      comment: "Created by" }
  created_at:          { type: datetime,     comment: "Created at" }
  updated_at:          { type: datetime,     comment: "Updated at" }
  excluded:            { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 8
  sub_form_size: 9
table_layout:
  default_order: [{field: source_question_id, order: ASC}]
```

## SURVEY_JSON
```yaml
table: survey_json
comment: Surveys JSON
columns:
  survey_id:    { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  title:        { type: varchar(255), nullable: false, comment: "Title",       form_display: true, table_display: true, form_size: 7, form_order: 1 }
  description:  { type: text,         comment: "Description",                 form_display: true, table_display: true, form_long_text: true, form_order: 4 }
  status:       { type: varchar(50),  default: draft, comment: "Status: draft | published | archived", form_display: true, table_display: true, form_size: 3, form_order: 3 }
  survey_json:  { type: text,         comment: "JSON Spec", form_display: true, table_display: true, form_long_text: true, form_code: json }
  survey_theme: { type: text,         comment: "JSON Theme", form_display: true, table_display: true, form_long_text: true, form_code: json }
  is_public:    { type: boolean,      default: false, comment: "Is Public", form_display: true, table_display: true, form_size: 2, form_order: 2 }
  user_id:      { type: integer,      comment: "Created by" }
  created_at:   { type: datetime,     comment: "Created at" }
  updated_at:   { type: datetime,     comment: "Updated at" }
  excluded:     { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf:
    - {label: Survey, fields: [title, is_public, status, description]}
    - {label: Survey JSON, fields: [survey_json]}
    - {label: Theme JSON, fields: [survey_theme]}
table_layout:
  default_order: [{field: survey_id, order: DESC}]
```

## RESPONSE
```yaml
table: response
comment: Survey Responses (raw submissions)
columns:
  response_id: { type: integer,      pk: true, autoincrement: true, comment: "ID" }
  survey_id:   { type: integer,      fk: "survey_json.survey_id", nullable: false, comment: "Survey", form_display: true, table_display: true, form_size: 4 }
  ip:          { type: varchar(45),  comment: "Submitter IP",    form_display: true, table_display: true, form_size: 3 }
  email:       { type: varchar(255), comment: "Submitter email", form_display: true, table_display: true, form_size: 4 }
  raw_json:    { type: text,         comment: "Raw response payload, as submitted", form_display: true, form_long_text: true, form_code: json }
  user_id:     { type: integer,      comment: "Created by" }
  created_at:  { type: datetime,     comment: "Submitted at", table_display: true }
  updated_at:  { type: datetime,     comment: "Updated at" }
  excluded:    { type: boolean,      default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  sub_form_size: 12
table_layout:
  default_order: [{field: created_at, order: DESC}]
```

# DATA
```yaml
name: DATA
description: DATA Model Survey Builder - seed default question types
runs_as: MODEL_DATA
database: SURVEY
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## QUESTION_TYPES
```yaml
table: question_type
description: Seed input types
cond: 'WHERE question_type_id = :question_type_id'
data:
  - question_type_id: 1
    type_name: text
    label: Single-line text
    icon: bars-3-bottom-left
    settings_schema: '{"min_length": "int", "max_length": "int", "regex_pattern": "string"}'
    active: true
  - question_type_id: 2
    type_name: textarea
    label: Multi-line text
    icon: bars-3-bottom-left
    settings_schema: '{"min_length": "int", "max_length": "int", "rows": "int", "cols": "int"}'
    active: true
  - question_type_id: 3
    type_name: number
    label: Number
    icon: hashtag
    settings_schema: '{"min_value": "number", "max_value": "number", "step": "number"}'
    active: true
  - question_type_id: 4
    type_name: radio
    label: Single choice
    icon: check-circle
    settings_schema: '{"choices": "array"}'
    active: true
  - question_type_id: 5
    type_name: checkbox
    label: Multiple choice
    icon: check
    settings_schema: '{"choices": "array", "min_selected": "int", "max_selected": "int"}'
    active: true
  - question_type_id: 6
    type_name: dropdown
    label: Dropdown
    icon: chevron-down
    settings_schema: '{"choices": "array", "is_searchable": "bool"}'
    active: true
  - question_type_id: 7
    type_name: date
    label: Date
    icon: calendar
    settings_schema: '{"min_date": "date", "max_date": "date", "date_format": "string"}'
    active: true
  - question_type_id: 8
    type_name: email
    label: Email
    icon: at-symbol
    settings_schema: '{"regex_pattern": "string"}'
    active: true
  - question_type_id: 9
    type_name: url
    label: URL
    icon: link
    settings_schema: '{"regex_pattern": "string"}'
    active: true
  - question_type_id: 10
    type_name: file
    label: File upload
    icon: paper-clip
    settings_schema: '{"allowed_file_types": "array", "max_file_size": "int", "max_files": "int"}'
    active: true
  - question_type_id: 11
    type_name: rating
    label: Rating
    icon: star
    settings_schema: '{"min_rating": "int", "max_rating": "int", "step": "int"}'
    active: true
  - question_type_id: 12
    type_name: slider
    label: Range slider
    icon: adjustments-horizontal
    settings_schema: '{"min_value": "number", "max_value": "number", "step": "number"}'
    active: true
```

## SURVEY_JSON
```yaml
table: survey_json
description: Survey JSON Example
cond: 'WHERE title = :title'
data:
  title: HOTEL BY THE SEA
  description: 1901 Thornridge Cir. Shiloh, Hawaii 81063 +1 (808) 555-0111
  status: dreaft
  survey_json: FileContent(assets/survey/survey.json)  
  survey_theme: FileContent(assets/survey/theme.json)
```

# ROLE_ACCESS
```yaml
name: ROLE_ACCESS
description: SURVEY main roles creation
database: SURVEY
runs_as: ROLE
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
```

## ANONYMOUS
```yaml
name: anonymous
description: Anonymous Role
access:
  - SURVEY:
    - Survey JSON:
      - {table: survey_json, read: true}
    - Survey Response:
      - {table: response, create: true, read: false, update: true}
active: true
```
