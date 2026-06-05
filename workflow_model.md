<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# WORKFLOW_MODEL
```yaml
name: WORKFLOW
description: Dynamic Workflow and Process Management Model
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
    menu_config: '{"label": "dashboard","tooltip": "dashboard_desc", "load_items": {"table": "dashboard", "tables": ["dashboard"]}}'
    tables:
      - dashboard
  Define Workflow:
    menu_icon: rectangle-group
    menu_order: 2
    active: true
    tables:
      - {table: workflow, requires_rla: true, active: true}
      - {table: workflow_sla, active: false}
      - {table: workflow_step, active: false}
      - {table: workflow_dependence, active: false}
      - {table: workflow_step_cond, active: false}
      - {table: workflow_step_sla, active: false}
      - {table: input_type, active: false}
      - {table: data_type, active: false}
      - {table: size, active: false}
      - {table: workflow_step_schema, active: false}
      - {table: workflow_step_schema_option, active: false}
      - {table: role, active: false}
      - {table: workflow_step_responsible, active: false}
      - {table: subscriber_type, active: false}
      - {table: workflow_step_subscriber, active: false}
      - {table: department, requires_rla: true, active: true}
      - {table: workflow_step_department, active: false}
  Execute Workflow:
    menu_icon: clipboard-document-check
    menu_order: 3
    active: true
    menu_config: '{"label": "workflow", "tooltip": "workflow_desc", "load_items": {"table": "workflow", "tables": ["workflow"]}}'
    tables:
      - {table: workflow, requires_rla: true, active: true}
      - {table: workflow_instance, requires_rla: false, active: false}
      - {table: workflow_instance_step, active: false}
      - {table: workflow_data, active: false}
      - {table: workflow_log, active: false}
      - {table: workflow_notification, active: false}
```

<!--WORKFLOW DEFINITION-->

## WORKFLOW
```yaml
table: workflow
comment: "Workflow"
tooltip: "Defines workflow processes"
columns:
  workflow_id:       { type: integer, pk: true, autoincrement: true, comment: "Workflow ID", tooltip: "Unique identifier of the workflow" }
  workflow:          { type: varchar, len: 200, unique: true, nullable: false, comment: "Workflow", tooltip: "Name of the workflow", form_display: true, table_display: true, form_size: 6, form_order: 1 }
  workflow_desc:     { type: text, comment: "Workflow Desc", tooltip: "Description of the workflow", form_display: true, table_display: true, form_long_text: true, form_order: 5 }
  order:             { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 2, form_order: 2 }
  version:           { type: varchar, len: 200, default: 'v1.0.0', comment: "Version", tooltip: "Version number of the workflow", form_display: true, table_display: true, form_size: 2, form_order: 3 }
  active:            { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the workflow is active", form_display: true, table_display: true, form_size: 2, form_order: 4 }  
  schedule:          { type: varchar, len: 200, comment: "Cron Schedule", tooltip: "Cron Representation of when it runs, if so", form_display: true, table_display: true, form_size: 4, form_order: 8 }
  steps_orientation: { type: varchar, len: 200, comment: "Step Orientation", tooltip: "Vertical / Horizontal", form_display: true, table_display: true, form_size: 4, form_order: 9 }
  workflow_icon:     { type: varchar, len: 200, comment: "Icon", tooltip: "Workflow Icon - Hero Icon", form_display: true, table_display: true, form_size: 4, form_order: 10 }
  email_template:    { type: text, comment: "Email Template", tooltip: "Email", form_display: true, form_long_text: true, form_code: html, form_order: 11 }
  user_id:           { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the workflow" }
  app_id:            { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:        { type: datetime, comment: "Created AT", tooltip: "Date and time when the workflow was created" }
  updated_at:        { type: datetime, comment: "Updated AT", tooltip: "Date and time when the workflow was last updated" }
  excluded:          { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the workflow is excluded from active use" }
form_layout: 
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  allow_in_subform: {workflow_step: true, workflow_dependence: true, workflow_sla: false}
  tabs_steps_conf:
    - {label: Workflow, fields: [workflow, order, version, step_color, active, workflow_desc, schedule, steps_orientation, workflow_icon]}
    - {label: Template, fields: [email_template]}
table_layout:
  default_order: [{field: order, order: ASC}]
```

## WORKFLOW_STEP
```yaml
table: workflow_step
comment: "Workflow Step"
tooltip: "Defines the steps of a workflow"
columns:
  workflow_step_id:    { type: integer, pk: true, autoincrement: true, comment: "Workflow Step ID", tooltip: "Unique identifier of the workflow step" }
  step:                { type: varchar, len: 200, nullable: false, comment: "Step", tooltip: "Name of the step", form_display: true, table_display: true, form_size: 4, form_order: 1 }
  step_desc:           { type: text, comment: "Step Desc", tooltip: "Description of the step", form_display: true, form_long_text: true, form_order: 6 }
  step_order:          { type: integer, comment: "Step Order", tooltip: "Order of execution of the step", form_display: true, table_display: true, form_size: 2, form_order: 2 }
  workflow_id:         { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow to which the step belongs", form_display: true, table_display: true, form_size: 4, form_order: 7 }
  step_icon:           { type: varchar, len: 200, comment: "Icon", tooltip: "Step Icon", form_display: true, table_display: true, form_size: 2, form_order: 3 }
  step_color:          { type: varchar, len: 200, comment: "Color", tooltip: "Step Color", form_display: true, table_display: true, form_size: 2, form_order: 4 }
  active:              { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step is active", form_display: true, form_size: 2, form_order: 5 }
  step_email_template: { type: text, comment: "Email Template", tooltip: "Email", form_display: true, form_code: html, form_long_text: true, form_order: 8 }
  document_template:   { type: text, comment: "Doc Template", tooltip: "In case the step is suposed to generate some kind of document, here will be the template, and it will be a gostatus templat tha has access to all the data from the previous step, current date, user, and the processes itself", form_display: true, form_code: html, form_long_text: true, form_order: 9 }
  api:                 { type: varchar, len: 255, comment: "Trigers API", tooltip: "API that is called", form_display: true, table_display: false, form_size: 8, form_order: 7 }
  user_id:             { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the step definition" }
  app_id:              { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:          { type: datetime, comment: "Created AT", tooltip: "Date and time when the step was created" }
  updated_at:          { type: datetime, comment: "Updated AT", tooltip: "Date and time when the step was last updated" }
  excluded:            { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the step is excluded from active use" }
form_layout: 
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  allow_in_subform:
    workflow_step_schema: true
    workflow_step_cond: true
    workflow_step_responsible: true
    workflow_step_subscriber: true
    workflow_step_department: true
    workflow_step_sla: false
  tabs_steps_conf:
    - {label: Step, fields: [step, step_order, step_icon, step_color, active, step_desc]}
    - {label: Conf, fields: [workflow_id, api, step_email_template, document_template]}
table_layout:
  default_order: [{field: step_order, order: ASC}]
```

## WORKFLOW_DEPENDENCE
```yaml
table: workflow_dependence
comment: "Workflow dependencies"
tooltip: "Defines workflow dependencies / relations"
columns:
  workflow_depend_id:    { type: integer, pk: true, autoincrement: true, comment: " ID" }
  workflow_depend:       { type: varchar, len: 200, nullable: false, comment: "Relation", form_display: true, table_display: true, form_size: 6, form_order: 3 }
  workflow_depend_desc:  { type: text, comment: "Relation Description", form_display: true, table_display: true, form_long_text: true, form_code: markdown, form_order: 6 }
  workflow_id:           { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Main Workflow ID", tooltip: "Current workflow", form_label: "Current Workflow", form_use_label: true, form_display: true, table_display: true, form_size: 6, form_order: 1 }
  depends_on:            { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Depends Workflow ID", form_label: "Depends On Workflow", form_use_label: true, form_display: true, table_display: true, form_size: 6, form_order: 2 }
  depend_order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 3, form_order: 4}
  active:                { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, form_order: 5 }  
  user_id:               { type: integer, comment: "User ID" }
  app_id:                { type: integer, comment: "App ID" }
  created_at:            { type: datetime, comment: "Created AT" }
  updated_at:            { type: datetime, comment: "Updated AT" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout: 
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order: [{field: depend_order, order: ASC}]
```

## WORKFLOW_SLA
```yaml
table: workflow_sla
comment: "Workflow SLA"
tooltip: "Defines Service Level Agreement rules for workflows"
columns:
  workflow_sla_id:  { type: integer, pk: true, autoincrement: true, comment: "Workflow SLA ID", tooltip: "Unique identifier of the workflow SLA rule" }
  workflow_id:      { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow", form_display: true, table_display: true, order: 5, form_size: 6 }
  name:             { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the SLA rule", form_display: true, table_display: true, order: 1, form_size: 6 }
  description:      { type: text, comment: "Description", tooltip: "Description of the SLA rule", form_display: true, form_long_text: true, form_order: 4 }
  duration_hours:   { type: integer, nullable: false, comment: "Duration Hours", tooltip: "SLA duration in hours", form_display: true, table_display: true, order: 6, form_size: 3 }
  escalation_hours: { type: integer, comment: "Escalation Hours", tooltip: "Hours before escalation is triggered", form_display: true, table_display: true, order: 7, form_size: 3 }
  priority:         { type: varchar, len: 50, comment: "Priority", tooltip: "Priority level for the SLA", table_display: true, form_display: true, order: 2, form_size: 3 }
  active:           { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the SLA rule is active", table_display: true, form_display: true, order: 3, form_size: 3 }
  user_id:          { type: integer, comment: "User ID", tooltip: "Identifier of the user who created the SLA rule" }
  created_at:       { type: datetime, comment: "Created AT", tooltip: "Date and time when the SLA rule was created" }
  updated_at:       { type: datetime, comment: "Updated AT", tooltip: "Date and time when the SLA rule was last updated" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout: 
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order: [{field: workflow_sla_id, order: DESC}]
```

## WORKFLOW_STEP_SLA
```yaml
table: workflow_step_sla
comment: "Step SLA"
tooltip: "Defines Service Level Agreement rules for workflow steps"
columns:
  workflow_step_sla_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step SLA ID", tooltip: "Unique identifier of the step SLA rule" }
  workflow_step_id:     { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  name:                 { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the step SLA rule", form_display: true, table_display: true }
  description:          { type: text, comment: "Description", tooltip: "Description of the step SLA rule", form_display: true }
  duration_hours:       { type: integer, nullable: false, comment: "Duration Hours", tooltip: "SLA duration in hours", form_display: true, table_display: true }
  escalation_hours:     { type: integer, comment: "Escalation Hours", tooltip: "Hours before escalation is triggered", form_display: true }
  priority:             { type: varchar, len: 50, comment: "Priority", tooltip: "Priority level for the step SLA", form_display: true }
  active:               { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step SLA rule is active" }
  user_id:              { type: integer, comment: "User ID", tooltip: "Identifier of the user who created the step SLA rule" }
  created_at:           { type: datetime, comment: "Created AT", tooltip: "Date and time when the step SLA rule was created" }
  updated_at:           { type: datetime, comment: "Updated AT", tooltip: "Date and time when the step SLA rule was last updated" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
table_layout:
  default_order: [{field: workflow_step_sla_id, order: DESC}]
```

## DEPARTMENT
```yaml
table: department
comment: "Department"
tooltip: "Represents an organizational department"
columns:
  department_id: { type: integer, pk: true, autoincrement: true, comment: "Department ID", tooltip: "Unique identifier of the department" }
  name:          { type: varchar, len: 200, nullable: false, comment: "Department Name", tooltip: "Name of the department", form_display: true, table_display: true }
  description:   { type: text, comment: "Description", tooltip: "Description of the department", form_display: true, table_display: true }
  active:        { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the department is active" }
  created_at:    { type: datetime, comment: "Created AT", tooltip: "Date and time when the department was created" }
  updated_at:    { type: datetime, comment: "Updated AT", tooltip: "Date and time when the department was last updated" }
  excluded:      { type: boolean, default: false, comment: "Excluded" }
table_layout:
  default_order: [{field: department_id, order: DESC}]
```

## WORKFLOW_STEP_DEPARTMENT
```yaml
table: workflow_step_department
comment: "Step Department"
tooltip: "Associates departments with workflow steps"
columns:
  workflow_step_department_id: { type: integer, pk: true, autoincrement: true, comment: "Department Workflow Step ID", tooltip: "Unique identifier of the relation" }
  department_id:               { type: integer, nullable: false, fk: "department.department_id", comment: "Department ID", tooltip: "Identifier of the department", form_display: true, table_display: true }
  workflow_step_id:            { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  active:                      { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the relation is active" }
  created_at:                  { type: datetime, comment: "Created AT", tooltip: "Date and time when the relation was created" }
  updated_at:                  { type: datetime, comment: "Updated AT", tooltip: "Date and time when the relation was last updated" }
  excluded:                    { type: boolean, default: false, comment: "Excluded" }
table_layout:
  default_order: [{field: workflow_step_department_id, order: DESC}]
```

## INPUT_TYPE
```yaml
table: input_type
comment: InputType
columns:
  input_type:   { type: varchar, len: 50, pk: true, comment: "ID" }
  input_type_desc: { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {input_type: text,     input_type_desc: text,     excluded: false}
  - {input_type: textarea, input_type_desc: textarea, excluded: false}
  - {input_type: password, input_type_desc: password, excluded: false}
  - {input_type: checkbox, input_type_desc: checkbox, excluded: false}
  - {input_type: radio,    input_type_desc: radio,    excluded: false}
  - {input_type: date,     input_type_desc: date,     excluded: false}
  - {input_type: datetime, input_type_desc: datetime, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
```

## DATA_TYPE
```yaml
table: data_type
comment: Input Type
columns:
  data_type:      { type: varchar, len: 50, pk: true, comment: "ID" }
  data_type_desc: { type: varchar, len: 50, unique: true, nullable: false, comment: "Data Type", form_display: true, table_display: true, order: 1 }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
data:
  - {data_type: text,     data_type_desc: text,     excluded: false}
  - {data_type: varchar,  data_type_desc: varchar,  excluded: false}
  - {data_type: boolean,  data_type_desc: boolean,  excluded: false}
  - {data_type: integer,  data_type_desc: integer,  excluded: false}
  - {data_type: decimal,  data_type_desc: decimal,  excluded: false}
  - {data_type: date,     data_type_desc: date,     excluded: false}
  - {data_type: datetime, data_type_desc: datetime, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
```

## SIZE
```yaml
table: size
comment: Size
columns:
  size:       { type: integer, pk: true, autoincrement: true, comment: "ID" }
  size_desc:  { type: varchar, len: 20, unique: true, nullable: false, comment: "Size", form_display: true, table_display: true, order: 1 }
  created_at: { type: datetime, comment: "Created at" }
  updated_at: { type: datetime, comment: "Updated at" }
  excluded:   { type: boolean, default: false, comment: "Excluded" }
data:
  - {size: 1,  size_desc: "1/12 - 8.33%",    excluded: false}
  - {size: 2,  size_desc: "2/12 - 16.67%",   excluded: false}
  - {size: 3,  size_desc: "3/12 - 25%",      excluded: false}
  - {size: 4,  size_desc: "4/12 - 33.33%",   excluded: false}
  - {size: 5,  size_desc: "5/12 - 41.67%",   excluded: false}
  - {size: 6,  size_desc: "6/12 - 50%",      excluded: false}
  - {size: 7,  size_desc: "7/12 - 58.33%",   excluded: false}
  - {size: 8,  size_desc: "8/12 - 66.67%",   excluded: false}
  - {size: 9,  size_desc: "9/12 - 75%",      excluded: false}
  - {size: 10, size_desc: "10/12 - 83.33%",  excluded: false}
  - {size: 11, size_desc: "11/12 - 91.67%",  excluded: false}
  - {size: 12, size_desc: "12/12 - 100%",    excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
```

## WORKFLOW_STEP_SCHEMA
```yaml
table: workflow_step_schema
comment: "Step Schema"
tooltip: "Defines the data structure required for each step of a workflow"
columns:
  workflow_step_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field" }
  workflow_id:             { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow associated with the field", form_display: true, table_display: true, form_size: 6 }
  workflow_step_id:        { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the field is collected", form_display: true, table_display: true, form_size: 6 }
  field:                   { type: varchar, len: 200, nullable: false, comment: "Field", tooltip: "Technical identifier of the field", form_display: true, table_display: true, form_size: 4 }
  label:                   { type: varchar, len: 200, nullable: false, comment: "Label", tooltip: "Display name of the field", form_display: true, table_display: true, form_size: 4 }
  data_type:               { type: varchar, fk: "data_type.data_type", nullable: false, comment: "Data Type", tooltip: "Type of data stored in the field", form_display: true, table_display: true, form_size: 4 }
  nullable:                { type: boolean, default: true, comment: "Nullable", tooltip: "Indicates whether the field can be empty", form_display: true, table_display: true, form_size: 3 }
  default_value:           { type: varchar, len: 200, comment: "Default Value", tooltip: "Default value assigned to the field", form_display: true, table_display: false, form_size: 3 }
  validation_rule:         { type: varchar, len: 200, comment: "Validation Rule", tooltip: "Regex validation rule for the field", form_display: true, table_display: false, form_size: 3 }
  order_index:             { type: integer, comment: "Order Index", tooltip: "Position of the field within the step", form_display: true, table_display: true, form_size: 3}
  format:                  { type: varchar, len: 200, comment: "Format", tooltip: "Format intl.Format", form_display: true, form_size: 4 }
  size:                    { type: integer, fk: "size.size", comment: "Size", tooltip: "1 - 12 size that will be shown in form", form_display: true, form_size: 2 }
  elipsis:                 { type: integer, comment: "Elipsis", tooltip: "Text elipsis", form_display: true, form_size: 2 }
  input_type:              { type: varchar, fk: "input_type.input_type", comment: "Options Input Type", tooltip: "Combobox,Checkbox or Radio", form_display: true, table_display: true, form_size: 2 }
  active:                  { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active" , form_display: true, table_display: true, form_size: 2 }
  options:                 { type: text, comment: "Options", tooltip: "JSON Array of string or array of objects{label,value}", form_display: true, form_long_text: true, form_size: 12 }
  user_id:                 { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition" }
  app_id:                  { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:              { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created" }
  updated_at:              { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated" }
  excluded:                { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_COND
```yaml
table: workflow_step_cond
comment: "Step Conditions"
tooltip: "Defines the data structure required for each to create workflow condition"
columns:
  workflow_step_cond_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field" }
  workflow_id:           { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow associated with the field", form_display: true, table_display: true, form_size: 5, form_order: 1 }
  workflow_step_id:      { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the field is collected", form_display: true, table_display: true, form_size: 5, form_order: 2  }
  cond_description:      { type: text, nullable: false, comment: "Description", tooltip: "Cndition Description", form_display: true, table_display: true, form_long_text: true, form_code: markdown, form_order: 4 }
  cond_trigger:          { type: text, nullable: false, comment: "Condition Trigger", tooltip: "JS Rule that when matched triger", form_display: true, table_display: true, form_long_text: true, form_code: js, form_order: 5 }
  cond_action:           { type: text, nullable: false, comment: "Condition Action", tooltip: "JS Rule run on triggered", form_display: true, table_display: true, form_long_text: true, form_code: js, form_order: 6 }
  active:                { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active", form_size: 2, form_order: 3 }
  user_id:               { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition" }
  app_id:                { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:            { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created" }
  updated_at:            { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated" }
  excluded:              { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 8
table_layout:
  default_order: [{field: workflow_step_cond_id, order: ASC}]
```

## WORKFLOW_STEP_SCHEMA_OPTION
```yaml
table: workflow_step_schema_option
comment: "Step Schema Option"
tooltip: "Defines the possible values for fields with predefined options"
columns:
  workflow_step_schema_option_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema Option ID", tooltip: "Unique identifier of the option" }
  workflow_step_schema_id:        { type: integer, nullable: false, fk: "workflow_step_schema.workflow_step_schema_id", comment: "Workflow Step Schema ID", tooltip: "Identifier of the field to which the option belongs", table_display: true  }
  option_value:                   { type: varchar, len: 200, nullable: false, comment: "Option Value", tooltip: "Stored value representing the option", form_display: true, table_display: true  }
  option_label:                   { type: varchar, len: 200, nullable: false, comment: "Option Label", tooltip: "Display label of the option", form_display: true, table_display: true  }
  order_index:                    { type: integer, comment: "Order Index", tooltip: "Position of the option within the list", form_display: true, table_display: true  }
  active:                         { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the option is active" }
  created_at:                     { type: datetime, comment: "Created AT", tooltip: "Date and time when the option was created" }
  updated_at:                     { type: datetime, comment: "Updated AT", tooltip: "Date and time when the option was last updated" }
  excluded:                       { type: boolean, default: false, comment: "Excluded" }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## ROLE
```yaml
table: role
comment: Role Type
columns:
  role:        { type: varchar, len: 50, pk: true, comment: "ID" }
  role_desc:   { type: varchar, len: 200, comment: "Role", form_display: true, table_display: true, order: 2 }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {role: Owner,      role_desc: Owner,      excluded: false}
  - {role: Observer,   role_desc: Observer,   excluded: false}
  - {role: Supervisor, role_desc: Supervisor, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 4
```

## WORKFLOW_STEP_RESPONSIBLE
```yaml
table: workflow_step_responsible
comment: "Step Responsible"
tooltip: "Defines the responsible users for each workflow step"
columns:
  workflow_step_responsible_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Responsible ID", tooltip: "Unique identifier of the assignment" }
  email:                        { type: varchar, len: 100, nullable: false, comment: "Email", tooltip: "Email associated with the responsibility", form_display: true, table_display: true, form_size: 4, order: 1 }
  first_name:                   { type: varchar, len: 50, comment: "First Name", form_display: true, table_display: true, form_size: 3, order: 2 }
  last_name:                    { type: varchar, len: 50, comment: "Last Name", form_display: true, table_display: true, form_size: 3, order: 3 }
  department_id:                { type: integer, comment: "Department ID", tooltip: "Identifier of the department responsible for the step", form_display: true, table_display: true, form_size: 4, order: 5 }
  role:                         { type: varchar, len: 100, comment: "Role", fk: "role.role", tooltip: "Role associated with the responsibility", form_display: true, table_display: true, form_size: 4, order: 6 }
  workflow_step_id:             { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step associated with the assignment", table_display: true, form_size: 4, order: 7 }
  active:                       { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the assignment is active", form_display: true, table_display: true, form_size: 2, order: 4 }
  responsible_email_template:   { type: text, comment: "Email Template", form_display: true, form_long_text: true, form_code: html}
  user_id:                      { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the step" }
  created_at:                   { type: datetime, comment: "Created AT", tooltip: "Date and time when the assignment was created" }
  updated_at:                   { type: datetime, comment: "Updated AT", tooltip: "Date and time when the option was last updated" }
  excluded:                     { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: workflow_step_responsible_id, order: ASC}]
```

## SUBSCRIBER_TYPE
```yaml
table: subscriber_type
comment: Subscriber Types
columns:
  subscriber_type:      { type: varchar, len: 50, pk: true, comment: "ID" }
  subscriber_type_desc: { type: varchar, len: 200, comment: "Subscriber Type", form_display: true, table_display: true, order: 2 }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
data:
  - {subscriber_type: Owner,      subscriber_type_desc: Owner,      excluded: false}
  - {subscriber_type: Observer,   subscriber_type_desc: Observer,   excluded: false}
  - {subscriber_type: Supervisor, subscriber_type_desc: Supervisor, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 4
```

## WORKFLOW_STEP_SUBSCRIBER
```yaml
table: workflow_step_subscriber
comment: "Step Subscriber"
tooltip: "Tracks interested parties and stakeholders for workflow steps"
columns:
  workflow_step_subscriber_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Subscriber ID", tooltip: "Unique identifier of the subscription" }
  email:                       { type: varchar, len: 100, nullable: false, comment: "Email", tooltip: "Email associated with the responsibility", form_display: true, table_display: true, form_size: 4, order: 1 }
  first_name:                  { type: varchar, len: 50, comment: "First Name", form_display: true, table_display: true, form_size: 3, order: 2 }
  last_name:                   { type: varchar, len: 50, comment: "Last Name", form_display: true, table_display: true, form_size: 3, order: 3 }
  subscriber_type:             { type: varchar, len: 50, comment: "Subscriber Type", fk: "subscriber_type.subscriber_type", tooltip: "Type of subscriber (responsible, observer, stakeholder, etc.)", form_display: true, table_display: true, form_size: 3, order: 5 }
  notify_on_start:             { type: boolean, default: true, comment: "Notify On Start", tooltip: "Send notification when step starts", form_display: true, table_display: true, form_size: 3, order: 6}
  notify_on_complete:          { type: boolean, default: true, comment: "Notify On Complete", tooltip: "Send notification when step completes", form_display: true, table_display: true, form_size: 3, order: 7}
  notify_on_escalation:        { type: boolean, default: false, comment: "Notify On Escalation", tooltip: "Send notification on SLA escalation", form_display: true, table_display: true, form_size: 3, order: 8}
  workflow_step_id:            { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true, form_size: 4, order: 9 }
  subscriber_email_template:   { type: text, comment: "Email Template", form_display: true, form_long_text: true, form_code: html}
  active:                      { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the assignment is active", form_display: true, table_display: true, form_size: 2, order: 4 }
  user_id:                     { type: integer, comment: "User ID", tooltip: "Identifier of the user interested in the step" }
  created_at:                  { type: datetime, comment: "Created AT", tooltip: "Date and time when the subscription was created" }
  updated_at:                  { type: datetime, comment: "Updated AT", tooltip: "Date and time when the subscription was last updated" }
  excluded:                    { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: workflow_step_subscriber_id, order: ASC}]
```

<!--WORKFLOW EXECUTION-->

## STATUS
```yaml
table: status
comment: Status
columns:
  status_id:   { type: integer, pk: true, autoincrement: true, comment: "Status ID" }
  status:      { type: varchar, len: 4, unique: true, nullable: false, comment: "Status", form_display: true, table_display: true, order: 1 }
  status_desc: { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {status_id: 1, status: Asigned, excluded: false}
  - {status_id: 2, status: Started, excluded: false}
  - {status_id: 3, status: Stabd By, excluded: false}
  - {status_id: 4, status: returned, excluded: false}
  - {status_id: 5, status: Conlcuded, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
```

## WORKFLOW_INSTANCE
```yaml
table: workflow_instance
comment: "Workflow Instance"
tooltip: "Represents an execution instance of a workflow"
columns:
  workflow_instance_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Instance ID", tooltip: "Unique identifier of the workflow instance" }
  workflow_id:          { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow being executed", table_display: true  }
  start_dt:             { type: datetime, nullable: false, comment: "Started AT", tooltip: "Date and time when the instance was started", form_display: true, table_display: true    }
  workflow_desc:        { type: text, nullable: false, comment: "Workflow Desc", tooltip: "Description of the workflow", form_display: true, table_display: true  }
  status_id:            { type: integer, fk: "status.status_id", comment: "Status", tooltip: "Current status of the workflow instance", form_display: true, table_display: true  }
  current_step_id:      { type: integer, comment: "Current Step ID", fk: "workflow_step.workflow_step_id", tooltip: "Identifier of the current step in execution", table_display: true  }
  active:               { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the instance is active" }
  created_at:           { type: datetime, comment: "Created AT", tooltip: "Date and time when the instance was created" }
  updated_at:           { type: datetime, comment: "Updated AT", tooltip: "Date and time when the instance was last updated" }
table_layout:
  default_order: [{field: start_dt, order: DESC}]
table_extra_options:
  - {component: Workflow, label: workflow, icon: play, size: 12, intercept_c: true, intercept_u: true}
```

## WORKFLOW_INSTANCE_STEP
```yaml
table: workflow_instance_step
comment: "Workflow Instance Step"
tooltip: "Represents the execution of each step within a workflow instance"
columns:
  workflow_instance_step_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Instance Step ID", tooltip: "Unique identifier of the instance step" }
  workflow_instance_id:      { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id:          { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step being executed", table_display: true  }
  workflow_step_status_id:   { type: integer, fk: "status.status_id", comment: "Status", tooltip: "Current status of the step", form_display: true, table_display: true  }
  assigned_to:               { type: integer, comment: "Assigned To", tooltip: "Identifier of the user assigned to the step" }
  started_at:                { type: datetime, comment: "Started AT", tooltip: "Date and time when the step execution started" }
  completed_at:              { type: datetime, comment: "Completed AT", tooltip: "Date and time when the step execution was completed" }
  active:                    { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step execution is active" }
  created_at:                { type: datetime, comment: "Created AT", tooltip: "Date and time when the record was created" }
  updated_at:                { type: datetime, comment: "Updated AT", tooltip: "Date and time when the instance was last updated" }
table_layout:
  default_order: [{field: workflow_instance_step_id, order: DESC}]
```

## WORKFLOW_DATA
```yaml
table: workflow_data
comment: "Workflow Data"
tooltip: "Stores values collected during workflow step execution"
columns:
  workflow_data_id:        { type: integer, pk: true, autoincrement: true, comment: "Workflow Data ID", tooltip: "Unique identifier of the stored value" }
  workflow_instance_id:    { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id:        { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the value was collected", table_display: true  }
  workflow_step_schema_id: { type: integer, nullable: false, fk: "workflow_step_schema.workflow_step_schema_id", comment: "Workflow Step Schema ID", tooltip: "Identifier of the field definition", table_display: true  }
  field:                   { type: varchar, len: 200, nullable: false, comment: "Field", tooltip: "Technical identifier of the field", table_display: true  }
  value:                   { type: text, comment: "Value", tooltip: "Value provided for the field" }
  is_latest:               { type: boolean, default: true, comment: "Is Latest", tooltip: "Indicates whether the record is the most recent value for the field" }
  created_at:              { type: datetime, comment: "Created AT", tooltip: "Date and time when the value was recorded", table_display: true  }
  updated_at:              { type: datetime, comment: "Updated AT", tooltip: "Date and time when the value was last updated" }
table_layout:
  default_order: [{field: workflow_data_id, order: DESC}]
```

## WORKFLOW_LOG
```yaml
table: workflow_log
comment: "Workflow Log"
tooltip: "Records all actions and state changes during workflow execution"
columns:
  workflow_log_id:      { type: integer, pk: true, autoincrement: true, comment: "Workflow Log ID", tooltip: "Unique identifier of the log entry" }
  workflow_instance_id: { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id:     { type: integer, comment: "Workflow Step ID", fk: "workflow_step.workflow_step_id", tooltip: "Identifier of the step associated with the action", table_display: true  }
  action:               { type: varchar, len: 50, nullable: false, comment: "Action", tooltip: "Type of action performed" }
  status_from:          { type: varchar, len: 50, comment: "Status From", tooltip: "Status before the action" }
  status_to:            { type: varchar, len: 50, comment: "Status To", tooltip: "Status after the action", table_display: true  }
  obs:                  { type: text, comment: "Obs", tooltip: "Additional information describing the action", form_long_text: true  }
  performed_by:         { type: integer, comment: "Performed By", tooltip: "Identifier of the user who performed the action" }
  created_at:           { type: datetime, comment: "Created AT", tooltip: "Date and time when the action was recorded", table_display: true  }
table_layout:
  default_order: [{field: workflow_log_id, order: DESC}]
```

## WORKFLOW_NOTIFICATION
```yaml
table: workflow_notification
comment: "Workflow Notification"
tooltip: "Tracks email and message notifications sent during workflow execution"
columns:
  workflow_notification_id:  { type: integer, pk: true, autoincrement: true, comment: "Workflow Notification ID", tooltip: "Unique identifier of the notification record" }
  workflow_instance_id:      { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true }
  workflow_instance_step_id: { type: integer, comment: "Workflow Instance Step ID", fk: "workflow_instance_step.workflow_instance_step_id", tooltip: "Identifier of the step execution, if applicable", table_display: true }
  recipient_user_id:         { type: integer, nullable: false, comment: "Recipient User ID", tooltip: "Identifier of the recipient user", form_display: true, table_display: true }
  notification_type:         { type: varchar, len: 50, nullable: false, comment: "Notification Type", tooltip: "Type of notification (step_started, step_completed, escalation, etc.)", form_display: true, table_display: true }
  subject:                   { type: varchar, len: 255, comment: "Subject", tooltip: "Email subject or notification title", form_display: true, table_display: true }
  message:                   { type: text, comment: "Message", tooltip: "Email body or notification message content", form_display: true, form_long_text: true }
  delivery_status:           { type: varchar, len: 50, default: "pending", comment: "Delivery Status", tooltip: "Status of delivery (pending, sent, failed, bounced)", form_display: true, table_display: true }
  delivery_attempts:         { type: integer, default: 0, comment: "Delivery Attempts", tooltip: "Number of delivery attempts made" }
  sent_at:                   { type: datetime, comment: "Sent AT", tooltip: "Date and time when the notification was sent", table_display: true }
  error_message:             { type: text, comment: "Error Message", tooltip: "Error details if delivery failed", form_long_text: true }
  created_at:                { type: datetime, comment: "Created AT", tooltip: "Date and time when the notification was created", table_display: true }
table_layout:
  default_order: [{field: workflow_notification_id, order: DESC}]
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

# WORKFLOW 1
```yaml
name: WORKFLOW_1
table: workflow
runs_as: WORKFLOW
description: Exemple of a workflow
icon: rectangle-group
order: 1
version: v1.0.0
orientation: vertical
database: WORKFLOW
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
```

## STEP 1
```yaml
name: STEP_1
table: workflow_step
description: Step 1
order: 1
icon: plus
color: green
active: true
schema:
  - {field: field1, label: field 1, data_type: text, input_type: radio, nullable: false, size: 3, options: '["A", "B", "C"]'}
responsibles:
  - {email: real.datadriven@gmail.com, first_name: real, last_name: datadriven, role: owner}
subscribers:
  - {email: real.datadriven@gmail.com, first_name: real, last_name: datadriven, start: true, complete: true}
```

## STEP 2
```yaml
name: STEP_2
table: workflow_step
description: Step 2
order: 2
icon: plus
color: yellow
active: true
schema:
  - {field: field1, label: field 1, data_type: text, input_type: radio, nullable: false, size: 3, options: '["A", "B", "C"]'}
responsibles:
  - {email: real.datadriven@gmail.com, first_name: real, last_name: datadriven, role: owner}
subscribers:
  - {email: real.datadriven@gmail.com, first_name: real, last_name: datadriven, start: true, complete: true}
```

---

# HELPDESK SUPPORT WORKFLOW

**Workflow Overview:** Complete IT Help Desk Support Process with Ticket Opening, Assignment, Resolution, Quality Verification, and Customer Approval.

```yaml
name: HELPDESK_SUPPORT
table: workflow
runs_as: WORKFLOW
description: Comprehensive IT Help Desk Support Workflow with Multi-Step Approval and Resolution Process
icon: lifebuoy
order: 2
version: v1.0.0
orientation: vertical
database: HELPDESK
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
email_template: |
  <h1>Support Ticket Update</h1>
  <p>Your support ticket #{{ticket_number}} has been updated.</p>
  <p><strong>Status:</strong> {{current_status}}</p>
  <p><strong>Current Step:</strong> {{current_step}}</p>
  <p><strong>Details:</strong> {{step_description}}</p>
```

## HELPDESK SLA RULES

```yaml
name: HELPDESK_SLA_CRITICAL
table: workflow_sla
workflow_name: HELPDESK_SUPPORT
sla_name: "Critical Priority SLA"
description: "SLA for critical issues - 1 hour response, 4 hours resolution"
duration_hours: 4
escalation_hours: 1
priority: critical
active: true
```

```yaml
name: HELPDESK_SLA_HIGH
table: workflow_sla
workflow_name: HELPDESK_SUPPORT
sla_name: "High Priority SLA"
description: "SLA for high priority issues - 4 hours response, 8 hours resolution"
duration_hours: 8
escalation_hours: 4
priority: high
active: true
```

```yaml
name: HELPDESK_SLA_NORMAL
table: workflow_sla
workflow_name: HELPDESK_SUPPORT
sla_name: "Normal Priority SLA"
description: "SLA for normal priority issues - 8 hours response, 24 hours resolution"
duration_hours: 24
escalation_hours: 12
priority: normal
active: true
```

## STEP 1: TICKET OPENING

**Description:** Customer submits a new support ticket with issue details.

```yaml
name: STEP_1_TICKET_OPENING
table: workflow_step
workflow_name: HELPDESK_SUPPORT
step: "Ticket Opening"
step_desc: "Customer initiates a support ticket by providing issue details, affected system, and impact assessment"
step_order: 1
step_icon: ticket-plus
step_color: blue
active: true
api: /api/helpdesk/ticket/create
step_email_template: |
  <h2>New Support Ticket Received</h2>
  <p>Your ticket <strong>#{{ticket_number}}</strong> has been successfully created.</p>
  <p><strong>Title:</strong> {{issue_title}}</p>
  <p><strong>Priority:</strong> {{priority_level}}</p>
  <p><strong>Created:</strong> {{created_at}}</p>

# Rich Schema for Ticket Opening
schema:
  - {field: customer_email, label: "Customer Email", data_type: varchar, input_type: text, nullable: false, size: 6, order_index: 1, validation_rule: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$', default_value: null}
  - {field: customer_name, label: "Customer Full Name", data_type: varchar, input_type: text, nullable: false, size: 6, order_index: 2}
  - {field: phone_number, label: "Contact Phone Number", data_type: varchar, input_type: text, nullable: true, size: 4, order_index: 3, validation_rule: '^\+?[0-9]{7,15}$'}
  - {field: issue_title, label: "Issue Title", data_type: varchar, input_type: text, nullable: false, size: 12, order_index: 4}
  - {field: issue_description, label: "Detailed Issue Description", data_type: text, input_type: textarea, nullable: false, size: 12, order_index: 5}
  - {field: affected_system, label: "Affected System/Application", data_type: varchar, input_type: radio, nullable: false, size: 6, order_index: 6, options: '[{"label": "Email System", "value": "email"}, {"label": "VPN", "value": "vpn"}, {"label": "File Server", "value": "fileserver"}, {"label": "Database", "value": "database"}, {"label": "Web Application", "value": "webapp"}, {"label": "Other", "value": "other"}]'}
  - {field: priority_level, label: "Priority Level", data_type: varchar, input_type: radio, nullable: false, size: 6, order_index: 7, options: '[{"label": "Critical - System Down", "value": "critical"}, {"label": "High - Significant Impact", "value": "high"}, {"label": "Normal - Minor Impact", "value": "normal"}, {"label": "Low - Documentation", "value": "low"}]', default_value: normal}
  - {field: impact_users_count, label: "Number of Users Affected", data_type: integer, input_type: text, nullable: false, size: 3, order_index: 8}
  - {field: has_workaround, label: "Is there a Workaround?", data_type: boolean, input_type: checkbox, nullable: true, size: 3, order_index: 9}
  - {field: workaround_description, label: "Workaround Details (if applicable)", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 10}
  - {field: ticket_attachment_ref, label: "Attachment Reference (URL or Path)", data_type: varchar, input_type: text, nullable: true, size: 12, order_index: 11}

# Step Conditions
conditions:
  - cond_description: "Auto-escalate if critical and no response within SLA"
    cond_trigger: |
      priority_level === 'critical' && 
      (new Date() - started_at) > 3600000
    cond_action: |
      escalation_flag = true;
      notify_manager = true;

  - cond_description: "Auto-assign low priority to self-service knowledge base"
    cond_trigger: |
      priority_level === 'low'
    cond_action: |
      assigned_to_team = 'knowledge_base';
      auto_response = true;

# Responsibles for Step
responsibles:
  - {email: support@company.com, first_name: Support, last_name: Team, role: owner, department_id: 3}

# Subscribers
subscribers:
  - {email: support-manager@company.com, first_name: Support, last_name: Manager, subscriber_type: supervisor, notify_on_start: true, notify_on_complete: true, notify_on_escalation: true}

# Step SLA
step_sla:
  - {name: "Ticket Opening SLA", description: "Time allowed for customer to complete ticket submission", duration_hours: 1, escalation_hours: 0.5, priority: normal, active: true}
```

## STEP 2: TICKET REVIEW & ASSIGNMENT

**Description:** Support team lead reviews ticket and assigns to appropriate technician based on expertise and workload.

```yaml
name: STEP_2_TICKET_REVIEW_ASSIGNMENT
table: workflow_step
workflow_name: HELPDESK_SUPPORT
step: "Ticket Review & Assignment"
step_desc: "Support supervisor reviews ticket details, determines urgency, and assigns to qualified technician"
step_order: 2
step_icon: user-check
step_color: purple
active: true
api: /api/helpdesk/ticket/assign
step_email_template: |
  <h2>Ticket Assignment Notice</h2>
  <p>Ticket <strong>#{{ticket_number}}</strong> has been assigned to you.</p>
  <p><strong>Title:</strong> {{issue_title}}</p>
  <p><strong>Priority:</strong> {{priority_level}}</p>
  <p><strong>Assigned By:</strong> {{assigned_by_name}}</p>
  <p><strong>SLA End Time:</strong> {{sla_deadline}}</p>

schema:
  - {field: ticket_number, label: "Ticket Number", data_type: varchar, input_type: text, nullable: false, size: 4, order_index: 1, default_value: "AUTO-GENERATED"}
  - {field: review_notes, label: "Supervisor Review Notes", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 2}
  - {field: assigned_technician, label: "Assign to Technician", data_type: varchar, input_type: radio, nullable: false, size: 6, order_index: 3, options: '[{"label": "John Smith - Network Expert", "value": "john.smith@company.com"}, {"label": "Sarah Johnson - Database DBA", "value": "sarah.johnson@company.com"}, {"label": "Mike Chen - Application Support", "value": "mike.chen@company.com"}, {"label": "Lisa Rodriguez - Infrastructure", "value": "lisa.rodriguez@company.com"}]'}
  - {field: urgency_assessment, label: "Reassess Priority", data_type: varchar, input_type: radio, nullable: false, size: 6, order_index: 4, options: '[{"label": "Confirmed Critical", "value": "critical"}, {"label": "High Priority", "value": "high"}, {"label": "Normal Priority", "value": "normal"}, {"label": "Can Wait", "value": "low"}]'}
  - {field: expected_resolution_time, label: "Expected Resolution Time (hours)", data_type: integer, input_type: text, nullable: true, size: 3, order_index: 5}
  - {field: requires_manager_approval, label: "Requires Manager Approval?", data_type: boolean, input_type: checkbox, nullable: true, size: 3, order_index: 6}
  - {field: approval_reason, label: "Reason for Manager Approval (if needed)", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 7}

conditions:
  - cond_description: "Route critical tickets requiring manager approval"
    cond_trigger: |
      requires_manager_approval === true &&
      urgency_assessment === 'critical'
    cond_action: |
      approval_required = true;
      route_to_manager = true;
      escalation_level = 'manager';

responsibles:
  - {email: support-supervisor@company.com, first_name: Support, last_name: Supervisor, role: owner, department_id: 3}

subscribers:
  - {email: it-director@company.com, first_name: IT, last_name: Director, subscriber_type: supervisor, notify_on_complete: true, notify_on_escalation: true}
  - {email: support-manager@company.com, first_name: Support, last_name: Manager, subscriber_type: observer, notify_on_start: true}
```

## STEP 3: ISSUE RESOLUTION

**Description:** Assigned technician works on resolving the issue, updating status and adding notes.

```yaml
name: STEP_3_ISSUE_RESOLUTION
table: workflow_step
workflow_name: HELPDESK_SUPPORT
step: "Issue Resolution"
step_desc: "Technician diagnoses and resolves the issue, maintaining detailed resolution logs"
step_order: 3
step_icon: wrench
step_color: yellow
active: true
api: /api/helpdesk/ticket/resolve
step_email_template: |
  <h2>Working on Your Ticket</h2>
  <p>Our technician is actively working on ticket <strong>#{{ticket_number}}</strong>.</p>
  <p><strong>Current Status:</strong> {{resolution_status}}</p>
  <p><strong>Progress:</strong> {{progress_percentage}}%</p>
  <p><strong>Last Update:</strong> {{last_update_time}}</p>

schema:
  - {field: resolution_status, label: "Current Resolution Status", data_type: varchar, input_type: radio, nullable: false, size: 6, order_index: 1, options: '[{"label": "Investigating", "value": "investigating"}, {"label": "Working on Fix", "value": "working"}, {"label": "Testing Solution", "value": "testing"}, {"label": "Ready for Customer Test", "value": "ready_test"}]', default_value: investigating}
  - {field: work_log, label: "Detailed Work Log", data_type: text, input_type: textarea, nullable: false, size: 12, order_index: 2}
  - {field: root_cause, label: "Root Cause Analysis", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 3}
  - {field: resolution_steps, label: "Resolution Steps Taken", data_type: text, input_type: textarea, nullable: false, size: 12, order_index: 4}
  - {field: prevention_measures, label: "Prevention Measures for Future", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 5}
  - {field: requires_customer_action, label: "Requires Customer Action?", data_type: boolean, input_type: checkbox, nullable: true, size: 3, order_index: 6}
  - {field: customer_action_details, label: "Customer Action Details", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 7}
  - {field: resolution_documentation_url, label: "Documentation/KB Article URL", data_type: varchar, input_type: text, nullable: true, size: 12, order_index: 8}
  - {field: time_spent_hours, label: "Time Spent (hours)", data_type: decimal, input_type: text, nullable: true, size: 3, order_index: 9}
  - {field: additional_resources_used, label: "Additional Resources/Tools Used", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 10}

conditions:
  - cond_description: "Alert if resolution exceeds SLA"
    cond_trigger: |
      (new Date() - started_at) > 
      (priority_level === 'critical' ? 14400000 : 
       priority_level === 'high' ? 28800000 : 86400000)
    cond_action: |
      sla_breach = true;
      notify_escalation = true;
      escalation_level = 'manager';

  - cond_description: "Route to QA if customer action required"
    cond_trigger: |
      requires_customer_action === true
    cond_action: |
      qa_review_required = true;
      step_comment = 'Awaiting customer verification';

responsibles:
  - {email: john.smith@company.com, first_name: John, last_name: Smith, role: owner, department_id: 3}
  - {email: mike.chen@company.com, first_name: Mike, last_name: Chen, role: owner, department_id: 3}

subscribers:
  - {email: support-manager@company.com, first_name: Support, last_name: Manager, subscriber_type: observer, notify_on_start: true, notify_on_complete: true, notify_on_escalation: true}

step_sla:
  - {name: "Resolution SLA", description: "Time allowed for issue resolution based on priority", duration_hours: 8, escalation_hours: 4, priority: normal, active: true}
```

## STEP 4: QUALITY CHECK & VERIFICATION

**Description:** QA supervisor verifies resolution quality before customer approval.

```yaml
name: STEP_4_QUALITY_CHECK
table: workflow_step
workflow_name: HELPDESK_SUPPORT
step: "Quality Check & Verification"
step_desc: "QA Team Lead verifies the solution quality, tests functionality, and ensures documentation completeness"
step_order: 4
step_icon: check-circle
step_color: green
active: true
api: /api/helpdesk/ticket/verify
step_email_template: |
  <h2>Quality Assurance Review</h2>
  <p>Ticket <strong>#{{ticket_number}}</strong> has passed quality review.</p>
  <p><strong>Resolution Quality:</strong> {{quality_score}}/10</p>
  <p><strong>QA Notes:</strong> {{qa_notes}}</p>
  <p><strong>Status:</strong> Ready for customer approval</p>

schema:
  - {field: qa_tested, label: "Have you tested the solution?", data_type: boolean, input_type: checkbox, nullable: false, size: 3, order_index: 1}
  - {field: qa_test_results, label: "Test Results and Findings", data_type: text, input_type: textarea, nullable: false, size: 12, order_index: 2}
  - {field: solution_quality_score, label: "Solution Quality Score (1-10)", data_type: integer, input_type: text, nullable: false, size: 2, order_index: 3, validation_rule: '^[1-9]$|^10$'}
  - {field: documentation_complete, label: "Documentation Complete?", data_type: boolean, input_type: checkbox, nullable: false, size: 3, order_index: 4}
  - {field: qa_recommendations, label: "QA Recommendations", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 5}
  - {field: quality_issues, label: "Any Quality Issues Found?", data_type: boolean, input_type: checkbox, nullable: true, size: 3, order_index: 6}
  - {field: issue_details, label: "Issue Details (if any)", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 7}
  - {field: rejection_reason, label: "Reason for Rejection (if applicable)", data_type: varchar, input_type: radio, nullable: true, size: 6, order_index: 8, options: '[{"label": "Incomplete Solution", "value": "incomplete"}, {"label": "Poor Documentation", "value": "poor_doc"}, {"label": "Test Failed", "value": "test_failed"}, {"label": "Does Not Meet Requirements", "value": "not_req"}]'}
  - {field: qa_approved, label: "QA Approval", data_type: boolean, input_type: checkbox, nullable: false, size: 3, order_index: 9}

conditions:
  - cond_description: "Reject and send back to technician if quality issues found"
    cond_trigger: |
      quality_issues === true &&
      solution_quality_score < 7
    cond_action: |
      qa_approved = false;
      send_back_to_technician = true;
      rejection_notification = true;
      step_status = 'rejected';

  - cond_description: "Auto-approve high quality solutions"
    cond_trigger: |
      solution_quality_score >= 9 &&
      documentation_complete === true
    cond_action: |
      qa_approved = true;
      expedite_customer_approval = true;

responsibles:
  - {email: qa-lead@company.com, first_name: QA, last_name: Lead, role: owner, department_id: 4}
  - {email: qa-team@company.com, first_name: QA, last_name: Team, role: owner, department_id: 4}

subscribers:
  - {email: support-manager@company.com, first_name: Support, last_name: Manager, subscriber_type: supervisor, notify_on_complete: true}
  - {email: john.smith@company.com, first_name: John, last_name: Smith, subscriber_type: observer, notify_on_start: true, notify_on_complete: true}
```

## STEP 5: CUSTOMER APPROVAL & CLOSURE

**Description:** Customer verifies resolution works, provides approval, and case is closed.

```yaml
name: STEP_5_CUSTOMER_APPROVAL_CLOSURE
table: workflow_step
workflow_name: HELPDESK_SUPPORT
step: "Customer Approval & Closure"
step_desc: "Customer verifies the solution works as expected and approves ticket closure"
step_order: 5
step_icon: thumbs-up
step_color: lime
active: true
api: /api/helpdesk/ticket/close
step_email_template: |
  <h2>Please Verify Your Resolution</h2>
  <p>Dear {{customer_name}},</p>
  <p>Your support ticket <strong>#{{ticket_number}}</strong> has been resolved.</p>
  <p><strong>Please test the solution</strong> and confirm it works for you.</p>
  <p>Click the link below to approve or request additional work:</p>
  <p><a href="{{approval_link}}">Review and Approve Resolution</a></p>
  <p>This ticket will auto-close in 48 hours if we don't hear back.</p>

schema:
  - {field: solution_works, label: "Does the solution work for you?", data_type: boolean, input_type: checkbox, nullable: false, size: 4, order_index: 1}
  - {field: solution_effectiveness, label: "Solution Effectiveness Rating (1-10)", data_type: integer, input_type: text, nullable: false, size: 4, order_index: 2, validation_rule: '^[1-9]$|^10$'}
  - {field: additional_issues, label: "Any Additional Issues?", data_type: boolean, input_type: checkbox, nullable: true, size: 3, order_index: 3}
  - {field: additional_issues_description, label: "Describe Additional Issues", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 4}
  - {field: customer_feedback, label: "Customer Feedback", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 5}
  - {field: support_satisfaction, label: "Support Team Performance Rating (1-10)", data_type: integer, input_type: text, nullable: true, size: 4, order_index: 6}
  - {field: satisfaction_comments, label: "Satisfaction Comments", data_type: text, input_type: textarea, nullable: true, size: 12, order_index: 7}
  - {field: would_recommend, label: "Would you recommend our support?", data_type: boolean, input_type: checkbox, nullable: true, size: 4, order_index: 8}
  - {field: customer_approval, label: "I approve this resolution and authorize ticket closure", data_type: boolean, input_type: checkbox, nullable: false, size: 12, order_index: 9}

conditions:
  - cond_description: "Auto-close ticket on customer approval"
    cond_trigger: |
      customer_approval === true &&
      solution_works === true
    cond_action: |
      ticket_status = 'closed';
      closure_date = new Date();
      send_closure_notification = true;
      archive_ticket = true;

  - cond_description: "Reopen ticket if additional issues found"
    cond_trigger: |
      additional_issues === true
    cond_action: |
      reopen_ticket = true;
      route_back_to_technician = true;
      new_priority = 'high';
      notification_sent = true;

  - cond_description: "Flag low satisfaction for review"
    cond_trigger: |
      support_satisfaction < 6
    cond_action: |
      flag_for_management_review = true;
      send_to_manager = true;
      quality_improvement_flag = true;

responsibles:
  - {email: support@company.com, first_name: Support, last_name: Team, role: owner, department_id: 3}

subscribers:
  - {email: support-manager@company.com, first_name: Support, last_name: Manager, subscriber_type: supervisor, notify_on_start: true, notify_on_complete: true}
  - {email: customer-success@company.com, first_name: Customer, last_name: Success, subscriber_type: observer, notify_on_complete: true}

step_sla:
  - {name: "Customer Approval SLA", description: "Time allowed for customer to approve or request changes", duration_hours: 48, escalation_hours: 24, priority: normal, active: true}
```

---

## HELPDESK WORKFLOW DEPENDENCIES

```yaml
name: HELPDESK_DEPENDENCIES
table: workflow_dependence
description: "Step Sequence for Helpdesk Workflow"
workflow_id: HELPDESK_SUPPORT
dependencies:
  - depends_on: STEP_1_TICKET_OPENING
    description: "Ticket must be opened before review"
    order: 1

  - depends_on: STEP_2_TICKET_REVIEW_ASSIGNMENT
    description: "Ticket must be reviewed and assigned before resolution"
    order: 2

  - depends_on: STEP_3_ISSUE_RESOLUTION
    description: "Issue must be resolved before QA verification"
    order: 3

  - depends_on: STEP_4_QUALITY_CHECK
    description: "QA verification required before customer approval"
    order: 4

  - depends_on: STEP_5_CUSTOMER_APPROVAL_CLOSURE
    description: "Customer approval required for final closure"
    order: 5
```
