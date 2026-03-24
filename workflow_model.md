<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
# WORKFLOW_MODEL
```yaml
name: WORKFLOW
description: Dynamic Workflow and Process Management Model
runs_as: MODEL
conn: 'sqlite3:database/WORKFLOW.db'
create_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Workflow:
    menu_icon: arrows-right-left
    menu_order: 1
    active: true
    tables:
      - workflow
      - workflow_step
      - workflow_step_schema
      - workflow_step_schema_option
      - workflow_step_responsible
  Execution:
    menu_icon: play
    menu_order: 2
    active: true
    tables:
      - workflow_instance
      - workflow_instance_step
      - workflow_data
      - workflow_log
```

## WORKFLOW
```yaml
table: workflow
comment: Defines workflow processes
columns:
  workflow_id:   { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the workflow" }
  workflow:      { type: varchar(200), unique: true, nullable: false, comment: "Name of the workflow", form_display: true, table_display: true, form_size: 6 }
  workflow_desc: { type: text, comment: "Description of the workflow", form_display: true, table_display: true }
  version:       { type: integer, default: 1, comment: "Version number of the workflow", form_display: true, table_display: true, form_size: 3 }
  active:        { type: boolean, default: true, comment: "Indicates whether the workflow is active", form_display: true, table_display: true, form_size: 3 }
  user_id:       { type: integer, comment: "Identifier of the user responsible for the workflow" }
  app_id:        { type: integer, comment: "Identifier of the application context" }
  created_at:    { type: datetime, comment: "Date and time when the workflow was created" }
  updated_at:    { type: datetime, comment: "Date and time when the workflow was last updated" }
  excluded:      { type: boolean, default: false, comment: "Indicates whether the workflow is excluded from active use" }
table_layout:
  default_order: [{field: workflow_id, order: DESC}]
```

## WORKFLOW_STEP
```yaml
table: workflow_step
comment: Defines the steps of a workflow
columns:
  workflow_step_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the workflow step" }
  workflow_id:      { type: integer, nullable: false, comment: "Identifier of the workflow to which the step belongs", form_display: true, table_display: true }
  step:             { type: varchar(200), nullable: false, comment: "Name of the step", form_display: true, table_display: true }
  step_desc:        { type: text, comment: "Description of the step", form_display: true }
  step_order:       { type: integer, comment: "Order of execution of the step", form_display: true, table_display: true }
  is_final:         { type: boolean, default: false, comment: "Indicates whether the step is the final step", form_display: true, table_display: true }
  active:           { type: boolean, default: true, comment: "Indicates whether the step is active" }
  user_id:          { type: integer, comment: "Identifier of the user responsible for the step definition" }
  app_id:           { type: integer, comment: "Identifier of the application context" }
  created_at:       { type: datetime, comment: "Date and time when the step was created" }
  updated_at:       { type: datetime, comment: "Date and time when the step was last updated" }
  excluded:         { type: boolean, default: false, comment: "Indicates whether the step is excluded from active use" }
table_layout:
  default_order: [{field: step_order, order: ASC}]
```

## WORKFLOW_STEP_SCHEMA
```yaml
table: workflow_step_schema
comment: Defines the data structure required for each step of a workflow
columns:
  workflow_step_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the schema field" }
  workflow_id:             { type: integer, nullable: false, comment: "Identifier of the workflow associated with the field", form_display: true, table_display: true }
  workflow_step_id:        { type: integer, nullable: false, comment: "Identifier of the step where the field is collected", form_display: true, table_display: true }
  field:                   { type: varchar(200), nullable: false, comment: "Technical identifier of the field", form_display: true, table_display: true }
  label:                   { type: varchar(200), nullable: false, comment: "Display name of the field", form_display: true, table_display: true }
  data_type:               { type: varchar(50), nullable: false, comment: "Type of data stored in the field", form_display: true }
  nullable:                { type: boolean, default: true, comment: "Indicates whether the field can be empty", form_display: true }
  default_value:           { type: text, comment: "Default value assigned to the field" }
  validation_rule:         { type: text, comment: "Rule used to validate the field value" }
  order_index:             { type: integer, comment: "Position of the field within the step", form_display: true, table_display: true }
  active:                  { type: boolean, default: true, comment: "Indicates whether the field is active" }
  user_id:                 { type: integer, comment: "Identifier of the user responsible for the field definition" }
  app_id:                  { type: integer, comment: "Identifier of the application context" }
  created_at:              { type: datetime, comment: "Date and time when the field was created" }
  updated_at:              { type: datetime, comment: "Date and time when the field was last updated" }
  excluded:                { type: boolean, default: false, comment: "Indicates whether the field is excluded from active use" }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_SCHEMA_OPTION
```yaml
table: workflow_step_schema_option
comment: Defines the possible values for fields with predefined options
columns:
  workflow_step_schema_option_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the option" }
  workflow_step_schema_id:        { type: integer, nullable: false, comment: "Identifier of the field to which the option belongs", table_display: true }
  option_value:                   { type: varchar(200), nullable: false, comment: "Stored value representing the option", form_display: true, table_display: true }
  option_label:                   { type: varchar(200), nullable: false, comment: "Display label of the option", form_display: true, table_display: true }
  order_index:                    { type: integer, comment: "Position of the option within the list", form_display: true, table_display: true }
  active:                         { type: boolean, default: true, comment: "Indicates whether the option is active" }
  created_at:                     { type: datetime, comment: "Date and time when the option was created" }
  updated_at:                     { type: datetime, comment: "Date and time when the option was last updated" }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

---

## WORKFLOW_STEP_RESPONSIBLE

```yaml
table: workflow_step_responsible
comment: Defines the responsible users for each workflow step
columns:
  workflow_step_responsible_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the assignment" }
  workflow_step_id:             { type: integer, nullable: false, comment: "Identifier of the step associated with the assignment", table_display: true }
  user_id:                      { type: integer, comment: "Identifier of the user responsible for the step" }
  role:                         { type: varchar(100), comment: "Role associated with the responsibility" }
  active:                       { type: boolean, default: true, comment: "Indicates whether the assignment is active" }
  created_at:                   { type: datetime, comment: "Date and time when the assignment was created" }
```

## WORKFLOW_INSTANCE
```yaml
table: workflow_instance
comment: Represents an execution instance of a workflow
columns:
  workflow_instance_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the workflow instance" }
  workflow_id:          { type: integer, nullable: false, comment: "Identifier of the workflow being executed", table_display: true }
  status:               { type: varchar(50), comment: "Current status of the workflow instance", form_display: true, table_display: true }
  current_step_id:      { type: integer, comment: "Identifier of the current step in execution", table_display: true }
  started_by:           { type: integer, comment: "Identifier of the user who started the workflow" }
  active:               { type: boolean, default: true, comment: "Indicates whether the instance is active" }
  created_at:           { type: datetime, comment: "Date and time when the instance was created" }
  updated_at:           { type: datetime, comment: "Date and time when the instance was last updated" }
table_layout:
  default_order: [{field: workflow_instance_id, order: DESC}]
```

## WORKFLOW_INSTANCE_STEP
```yaml
table: workflow_instance_step
comment: Represents the execution of each step within a workflow instance
columns:
  workflow_instance_step_id: { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the instance step" }
  workflow_instance_id:      { type: integer, nullable: false, comment: "Identifier of the workflow instance", table_display: true }
  workflow_step_id:          { type: integer, nullable: false, comment: "Identifier of the step being executed", table_display: true }
  status:                    { type: varchar(50), comment: "Current status of the step", form_display: true, table_display: true }
  assigned_to:               { type: integer, comment: "Identifier of the user assigned to the step" }
  started_at:                { type: datetime, comment: "Date and time when the step execution started" }
  completed_at:              { type: datetime, comment: "Date and time when the step execution was completed" }
  active:                    { type: boolean, default: true, comment: "Indicates whether the step execution is active" }
  created_at:                { type: datetime, comment: "Date and time when the record was created" }
table_layout:
  default_order: [{field: workflow_instance_step_id, order: DESC}]
```

## WORKFLOW_DATA
```yaml
table: workflow_data
comment: Stores values collected during workflow step execution
columns:
  workflow_data_id:        { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the stored value" }
  workflow_instance_id:    { type: integer, nullable: false, comment: "Identifier of the workflow instance", table_display: true }
  workflow_step_id:        { type: integer, nullable: false, comment: "Identifier of the step where the value was collected", table_display: true }
  workflow_step_schema_id: { type: integer, nullable: false, comment: "Identifier of the field definition", table_display: true }
  field:                   { type: varchar(200), nullable: false, comment: "Technical identifier of the field", table_display: true }
  value:                   { type: text, comment: "Value provided for the field" }
  is_latest:               { type: boolean, default: true, comment: "Indicates whether the record is the most recent value for the field" }
  created_at:              { type: datetime, comment: "Date and time when the value was recorded", table_display: true }
  updated_at:              { type: datetime, comment: "Date and time when the value was last updated" }
table_layout:
  default_order: [{field: workflow_data_id, order: DESC}]
```

## WORKFLOW_LOG
```yaml
table: workflow_log
comment: Records all actions and state changes during workflow execution
columns:
  workflow_log_id:      { type: integer, pk: true, autoincrement: true, comment: "Unique identifier of the log entry" }
  workflow_instance_id: { type: integer, nullable: false, comment: "Identifier of the workflow instance", table_display: true }
  workflow_step_id:     { type: integer, comment: "Identifier of the step associated with the action", table_display: true }
  action:               { type: varchar(50), nullable: false, comment: "Type of action performed" }
  status_from:          { type: varchar(50), comment: "Status before the action" }
  status_to:            { type: varchar(50), comment: "Status after the action", table_display: true }
  obs:                  { type: text, comment: "Additional information describing the action", form_long_text: true }
  performed_by:         { type: integer, comment: "Identifier of the user who performed the action" }
  created_at:           { type: datetime, comment: "Date and time when the action was recorded", table_display: true }
table_layout:
  default_order: [{field: workflow_log_id, order: DESC}]
```
