<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
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
      - {table: workflow_dependence, actve: false}
      - {table: workflow_step_cond, active: false}
      - {table: workflow_step_sla, active: false}
      - {table: input_type, active: false}
      - {table: data_type, active: false}
      - {table: size, active: false}
      - {table: workflow_step_schema, active: false}
      - {table: workflow_step_schema_option, active: false}
      - {table: workflow_step_responsible, active: false}
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
  allow_in_subform: {workflow_step: true, workflow_dependence: true, workflow_sla: true}
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
    workflow_step_sla: true
  tabs_steps_conf:
    - {label: Main, fields: [step, step_order, step_icon, step_color, active, step_desc]}
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
  workflow_ depend:      { type: varchar, len: 200, nullable: false, comment: "Relation", form_display: true, table_display: true, form_size: 6, form_order: 3 }
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
  size: 8
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
  workflow_id:      { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow", form_display: true, table_display: true }
  name:             { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the SLA rule", form_display: true, table_display: true }
  description:      { type: text, comment: "Description", tooltip: "Description of the SLA rule", form_display: true }
  duration_hours:   { type: integer, nullable: false, comment: "Duration Hours", tooltip: "SLA duration in hours", form_display: true, table_display: true }
  escalation_hours: { type: integer, comment: "Escalation Hours", tooltip: "Hours before escalation is triggered", form_display: true }
  priority:         { type: varchar, len: 50, comment: "Priority", tooltip: "Priority level for the SLA", form_display: true }
  active:           { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the SLA rule is active" }
  user_id:          { type: integer, comment: "User ID", tooltip: "Identifier of the user who created the SLA rule" }
  created_at:       { type: datetime, comment: "Created AT", tooltip: "Date and time when the SLA rule was created" }
  updated_at:       { type: datetime, comment: "Updated AT", tooltip: "Date and time when the SLA rule was last updated" }
table_layout:
  default_order: [{field: workflow_sla_id, order: DESC}]
```

## WORKFLOW_STEP_SLA
```yaml
table: workflow_step_sla
comment: "Workflow Step SLA"
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
  name:          { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the department", form_display: true, table_display: true }
  description:   { type: text, comment: "Description", tooltip: "Description of the department", form_display: true, table_display: true }
  active:        { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the department is active" }
  created_at:    { type: datetime, comment: "Created AT", tooltip: "Date and time when the department was created" }
  updated_at:    { type: datetime, comment: "Updated AT", tooltip: "Date and time when the department was last updated" }
table_layout:
  default_order: [{field: department_id, order: DESC}]
```

## WORKFLOW_STEP_DEPARTMENT
```yaml
table: workflow_step_department
comment: "Department Workflow Step"
tooltip: "Associates departments with workflow steps"
columns:
  workflow_step_department_id: { type: integer, pk: true, autoincrement: true, comment: "Department Workflow Step ID", tooltip: "Unique identifier of the relation" }
  department_id:               { type: integer, nullable: false, fk: "department.department_id", comment: "Department ID", tooltip: "Identifier of the department", form_display: true, table_display: true }
  workflow_step_id:            { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  active:                      { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the relation is active" }
  created_at:                  { type: datetime, comment: "Created AT", tooltip: "Date and time when the relation was created" }
table_layout:
  default_order: [{field: workflow_step_department_id, order: DESC}]
```

## INPUT_TYPE
```yaml
table: input_type
comment: InputType
columns:
  input_type_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  input_type:      { type: varchar, len: 4, unique: true, nullable: false, comment: "Input Type", form_display: true, table_display: true, order: 1 }
  input_type_desc: { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {input_type_id: 1, input_type: text, excluded: false}
  - {input_type_id: 2, input_type: textarea, excluded: false}
  - {input_type_id: 3, input_type: password, excluded: false}
  - {input_type_id: 4, input_type: checkbox, excluded: false}
  - {input_type_id: 5, input_type: radio, excluded: false}
  - {input_type_id: 6, input_type: date, excluded: false}
  - {input_type_id: 7, input_type: datetime, excluded: false}
form_layout:
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
```

## DATA_TYPE
```yaml
table: data_type
comment: Input Type
columns:
  data_type_id:   { type: integer, pk: true, autoincrement: true, comment: "ID" }
  data_type:      { type: varchar, len: 4, unique: true, nullable: false, comment: "Data Type", form_display: true, table_display: true, order: 1 }
  data_type_desc: { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
data:
  - {data_type_id: 1, data_type: text, excluded: false}
  - {data_type_id: 2, data_type: varchar, excluded: false}
  - {data_type_id: 3, data_type: boolean, excluded: false}
  - {data_type_id: 4, data_type: integer, excluded: false}
  - {data_type_id: 5, data_type: decimal, excluded: false}
  - {data_type_id: 6, data_type: date, excluded: false}
  - {data_type_id: 7, data_type: datetime, excluded: false}
form_layout:
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
```

## SIZE
```yaml
table: size
comment: Size
columns:
  size_id:     { type: integer, pk: true, autoincrement: true, comment: "ID" }
  size:        { type: varchar, len: 4, unique: true, nullable: false, comment: "Data Type", form_display: true, table_display: true, order: 1 }
  size_desc:   { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, order: 2 }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {size_id: 1, size: "1/12 - 8.33%", excluded: false}
  - {size_id: 2, size: "2/12 - 16.67%", excluded: false}
  - {size_id: 3, size: "3/12 - 25%", excluded: false}
  - {size_id: 4, size: "4/12 - 33.33%", excluded: false}
  - {size_id: 5, size: "5/12 - 41.67%", excluded: false}
  - {size_id: 6, size: "6/12 - 50%", excluded: false}
  - {size_id: 7, size: "7/12 - 58.33%", excluded: false}
  - {size_id: 8, size: "8/12 - 66.67%", excluded: false}
  - {size_id: 9, size: "9/12 - 75%", excluded: false}
  - {size_id: 10, size: "10/12 - 83.33%", excluded: false}
  - {size_id: 11, size: "11/12 - 91.67%", excluded: false}
  - {size_id: 12, size: "12/12 - 100%", excluded: false}
form_layout:
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
```

## WORKFLOW_STEP_SCHEMA
```yaml
table: workflow_step_schema
comment: "Workflow Step Schema"
tooltip: "Defines the data structure required for each step of a workflow"
columns:
  workflow_step_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field" }
  workflow_id:             { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow associated with the field", form_display: true, table_display: true  }
  workflow_step_id:        { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the field is collected", form_display: true, table_display: true  }
  field:                   { type: varchar, len: 200, nullable: false, comment: "Field", tooltip: "Technical identifier of the field", form_display: true, table_display: true  }
  label:                   { type: varchar, len: 200, nullable: false, comment: "Label", tooltip: "Display name of the field", form_display: true, table_display: true  }
  data_type_id:            { type: integer, fk: "data_type.data_type_id", nullable: false, comment: "Data Type", tooltip: "Type of data stored in the field", form_display: true  }
  nullable:                { type: boolean, default: true, comment: "Nullable", tooltip: "Indicates whether the field can be empty", form_display: true  }
  default_value:           { type: varchar, len: 200, comment: "Default Value", tooltip: "Default value assigned to the field" }
  validation_rule:         { type: varchar, len: 200, comment: "Validation Rule", tooltip: "Regex validation rule for the field" }
  order_index:             { type: integer, comment: "Order Index", tooltip: "Position of the field within the step", form_display: true, table_display: true  }
  format:                  { type: varchar, len: 200, comment: "Format", tooltip: "Format intl.Format" }
  size_id:                 { type: integer, fk: "size.size_id", comment: "Size", tooltip: "1 - 12 size that will be shown in form" }
  elipsis:                 { type: integer, comment: "Elipsis", tooltip: "Text elipsis" }
  options:                 { type: text, comment: "Options", tooltip: "JSON Array of string or array of objects{label,value}" }
  input_type_id:           { type: integer, fk: "input_type.input_type_id", comment: "Options Input Type", tooltip: "Combobox,Checkbox or Radio" }
  active:                  { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active" }
  user_id:                 { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition" }
  app_id:                  { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:              { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created" }
  updated_at:              { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated" }
  excluded:                { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use" }
form_layout:
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_COND
```yaml
table: workflow_step_cond
comment: "Workflow Step Conditions"
tooltip: "Defines the data structure required for each to create workflow condition"
columns:
  workflow_step_cond_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field" }
  workflow_id:           { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow associated with the field", form_display: true, table_display: true  }
  workflow_step_id:      { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the field is collected", form_display: true, table_display: true  }
  cond_description:      { type: text, nullable: false, comment: "Description", tooltip: "Cndition Description", form_display: true, table_display: true, form_long_text: true, form_code: markdown }
  cond_trigger:          { type: text, nullable: false, comment: "Condition Trigger", tooltip: "JS Rule that when matched triger", form_display: true, table_display: true, form_long_text: true, form_code: js }
  cond_action:           { type: text, nullable: false, comment: "Condition Action", tooltip: "JS Rule run on triggered", form_display: true, table_display: true, form_long_text: true, form_code: js }
  active:                { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active" }
  user_id:               { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition" }
  app_id:                { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:            { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created" }
  updated_at:            { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated" }
  excluded:              { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use" }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_SCHEMA_OPTION
```yaml
table: workflow_step_schema_option
comment: "Workflow Step Schema Option"
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
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_RESPONSIBLE
```yaml
table: workflow_step_responsible
comment: "Workflow Step Responsible"
tooltip: "Defines the responsible users for each workflow step"
columns:
  workflow_step_responsible_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Responsible ID", tooltip: "Unique identifier of the assignment" }
  workflow_step_id:             { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step associated with the assignment", table_display: true  }
  user_id:                      { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the step" }
  email:                        { type: varchar, len: 100, comment: "Email", tooltip: "Email associated with the responsibility" }
  first_name:                   { type: varchar, len: 50, nullable: false, comment: "First Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  last_name:                    { type: varchar, len: 50, comment: "Last Name", form_display: true, table_display: true, form_size: 4, order: 3 }
  department_id:                { type: integer, comment: "Department ID", tooltip: "Identifier of the department responsible for the step" }
  role:                         { type: varchar, len: 100, comment: "Role", tooltip: "Role associated with the responsibility" }
  active:                       { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the assignment is active" }
  created_at:                   { type: datetime, comment: "Created AT", tooltip: "Date and time when the assignment was created" }
  updated_at:                   { type: datetime, comment: "Updated AT", tooltip: "Date and time when the option was last updated" }
```

## WORKFLOW_STEP_SUBSCRIBER
```yaml
table: workflow_step_subscriber
comment: "Workflow Step Subscriber"
tooltip: "Tracks interested parties and stakeholders for workflow steps"
columns:
  workflow_step_subscriber_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Subscriber ID", tooltip: "Unique identifier of the subscription" }
  workflow_step_id:            { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  user_id:                     { type: integer, nullable: false, comment: "User ID", tooltip: "Identifier of the user interested in the step", form_display: true, table_display: true }
  subscriber_type:             { type: varchar, len: 50, comment: "Subscriber Type", tooltip: "Type of subscriber (responsible, observer, stakeholder, etc.)", form_display: true, table_display: true }
  notify_on_start:             { type: boolean, default: true, comment: "Notify On Start", tooltip: "Send notification when step starts" }
  notify_on_complete:          { type: boolean, default: true, comment: "Notify On Complete", tooltip: "Send notification when step completes" }
  notify_on_escalation:        { type: boolean, default: false, comment: "Notify On Escalation", tooltip: "Send notification on SLA escalation" }
  active:                      { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the subscription is active" }
  created_at:                  { type: datetime, comment: "Created AT", tooltip: "Date and time when the subscription was created" }
  updated_at:                  { type: datetime, comment: "Updated AT", tooltip: "Date and time when the subscription was last updated" }
table_layout:
  default_order: [{field: workflow_step_subscriber_id, order: DESC}]
```

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
  tabs_steps: deactivate
  form_in_popup: true
  size: 6
```

<!--WORKFLOW EXECUTION-->

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
description: Exemple of a workflow
runs_as: WORKFLOW
database: WORKFLOW
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
```

## STEP 1
```yaml
name: WORKFLOW_1
description: Exemple of a workflow
active: true
attributes:
  field1: {label: field 1, data_type: text, options: '["A", "B", "C"]'}
```
