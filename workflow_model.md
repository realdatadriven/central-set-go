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
      - workflow_sla
      - workflow_step
      - workflow_step_sla
      - workflow_step_schema
      - workflow_step_schema_option
      - workflow_step_responsible
      - department
      - department_workflow_step
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
comment: "Workflow"
tooltip: "Defines workflow processes"
columns:
  workflow_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow ID", tooltip: "Unique identifier of the workflow"  }
  workflow: { type: varchar(200), unique: true, nullable: false, comment: "Workflow", tooltip: "Name of the workflow", form_display: true, table_display: true, form_size: 6  }
  workflow_desc: { type: text, comment: "Workflow Desc", tooltip: "Description of the workflow", form_display: true, table_display: true  }
  version: { type: integer, default: 1, comment: "Version", tooltip: "Version number of the workflow", form_display: true, table_display: true, form_size: 3  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the workflow is active", form_display: true, table_display: true, form_size: 3  }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the workflow"  }
  app_id: { type: integer, comment: "App ID", tooltip: "Identifier of the application context"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the workflow was created"  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the workflow was last updated"  }
  excluded: { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the workflow is excluded from active use"  }
table_layout:
  default_order: [{field: workflow_id, order: DESC}]
```

## WORKFLOW_STEP
```yaml
table: workflow_step
comment: "Workflow Step"
tooltip: "Defines the steps of a workflow"
columns:
  workflow_step_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step ID", tooltip: "Unique identifier of the workflow step"  }
  workflow_id: { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow to which the step belongs", form_display: true, table_display: true  }
  step: { type: varchar(200), nullable: false, comment: "Step", tooltip: "Name of the step", form_display: true, table_display: true  }
  step_desc: { type: text, comment: "Step Desc", tooltip: "Description of the step", form_display: true  }
  step_order: { type: integer, comment: "Step Order", tooltip: "Order of execution of the step", form_display: true, table_display: true  }
  is_final: { type: boolean, default: false, comment: "Is Final", tooltip: "Indicates whether the step is the final step", form_display: true, table_display: true  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step is active"  }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the step definition"  }
  app_id: { type: integer, comment: "App ID", tooltip: "Identifier of the application context"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the step was created"  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the step was last updated"  }
  excluded: { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the step is excluded from active use"  }
table_layout:
  default_order: [{field: step_order, order: ASC}]
```

## WORKFLOW_SLA
```yaml
table: workflow_sla
comment: "Workflow SLA"
tooltip: "Defines Service Level Agreement rules for workflows"
columns:
  workflow_sla_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow SLA ID", tooltip: "Unique identifier of the workflow SLA rule" }
  workflow_id: { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow", form_display: true, table_display: true }
  name: { type: varchar(200), nullable: false, comment: "Name", tooltip: "Name of the SLA rule", form_display: true, table_display: true }
  description: { type: text, comment: "Description", tooltip: "Description of the SLA rule", form_display: true }
  duration_hours: { type: integer, nullable: false, comment: "Duration Hours", tooltip: "SLA duration in hours", form_display: true, table_display: true }
  escalation_hours: { type: integer, comment: "Escalation Hours", tooltip: "Hours before escalation is triggered", form_display: true }
  priority: { type: varchar(50), comment: "Priority", tooltip: "Priority level for the SLA", form_display: true }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the SLA rule is active" }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user who created the SLA rule" }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the SLA rule was created" }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the SLA rule was last updated" }
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
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  name: { type: varchar(200), nullable: false, comment: "Name", tooltip: "Name of the step SLA rule", form_display: true, table_display: true }
  description: { type: text, comment: "Description", tooltip: "Description of the step SLA rule", form_display: true }
  duration_hours: { type: integer, nullable: false, comment: "Duration Hours", tooltip: "SLA duration in hours", form_display: true, table_display: true }
  escalation_hours: { type: integer, comment: "Escalation Hours", tooltip: "Hours before escalation is triggered", form_display: true }
  priority: { type: varchar(50), comment: "Priority", tooltip: "Priority level for the step SLA", form_display: true }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step SLA rule is active" }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user who created the step SLA rule" }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the step SLA rule was created" }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the step SLA rule was last updated" }
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
  name: { type: varchar(200), nullable: false, comment: "Name", tooltip: "Name of the department", form_display: true, table_display: true }
  description: { type: text, comment: "Description", tooltip: "Description of the department", form_display: true, table_display: true }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the department is active" }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the department was created" }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the department was last updated" }
table_layout:
  default_order: [{field: department_id, order: DESC}]
```

## DEPARTMENT_WORKFLOW_STEP
```yaml
table: department_workflow_step
comment: "Department Workflow Step"
tooltip: "Associates departments with workflow steps"
columns:
  department_workflow_step_id: { type: integer, pk: true, autoincrement: true, comment: "Department Workflow Step ID", tooltip: "Unique identifier of the relation" }
  department_id: { type: integer, nullable: false, fk: "department.department_id", comment: "Department ID", tooltip: "Identifier of the department", form_display: true, table_display: true }
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the workflow step", form_display: true, table_display: true }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the relation is active" }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the relation was created" }
table_layout:
  default_order: [{field: department_workflow_step_id, order: DESC}]
```

## WORKFLOW_STEP_SCHEMA
```yaml
table: workflow_step_schema
comment: "Workflow Step Schema"
tooltip: "Defines the data structure required for each step of a workflow"
columns:
  workflow_step_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field"  }
  workflow_id: { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow associated with the field", form_display: true, table_display: true  }
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the field is collected", form_display: true, table_display: true  }
  field: { type: varchar(200), nullable: false, comment: "Field", tooltip: "Technical identifier of the field", form_display: true, table_display: true  }
  label: { type: varchar(200), nullable: false, comment: "Label", tooltip: "Display name of the field", form_display: true, table_display: true  }
  data_type: { type: varchar(50), nullable: false, comment: "Data Type", tooltip: "Type of data stored in the field", form_display: true  }
  nullable: { type: boolean, default: true, comment: "Nullable", tooltip: "Indicates whether the field can be empty", form_display: true  }
  default_value: { type: text, comment: "Default Value", tooltip: "Default value assigned to the field"  }
  validation_rule: { type: text, comment: "Validation Rule", tooltip: "Rule used to validate the field value"  }
  order_index: { type: integer, comment: "Order Index", tooltip: "Position of the field within the step", form_display: true, table_display: true  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active"  }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition"  }
  app_id: { type: integer, comment: "App ID", tooltip: "Identifier of the application context"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created"  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated"  }
  excluded: { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use"  }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_SCHEMA_OPTION
```yaml
table: workflow_step_schema_option
comment: "Workflow Step Schema Option"
tooltip: "Defines the possible values for fields with predefined options"
columns:
  workflow_step_schema_option_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema Option ID", tooltip: "Unique identifier of the option"  }
  workflow_step_schema_id: { type: integer, nullable: false, fk: "workflow_step_schema.workflow_step_schema_id", comment: "Workflow Step Schema ID", tooltip: "Identifier of the field to which the option belongs", table_display: true  }
  option_value: { type: varchar(200), nullable: false, comment: "Option Value", tooltip: "Stored value representing the option", form_display: true, table_display: true  }
  option_label: { type: varchar(200), nullable: false, comment: "Option Label", tooltip: "Display label of the option", form_display: true, table_display: true  }
  order_index: { type: integer, comment: "Order Index", tooltip: "Position of the option within the list", form_display: true, table_display: true  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the option is active"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the option was created"  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the option was last updated"  }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## WORKFLOW_STEP_RESPONSIBLE
```yaml
table: workflow_step_responsible
comment: "Workflow Step Responsible"
tooltip: "Defines the responsible users for each workflow step"
columns:
  workflow_step_responsible_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Responsible ID", tooltip: "Unique identifier of the assignment"  }
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step associated with the assignment", table_display: true  }
  user_id: { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the step"  }
  department_id: { type: integer, comment: "Department ID", tooltip: "Identifier of the department responsible for the step"  }
  role: { type: varchar(100), comment: "Role", tooltip: "Role associated with the responsibility"  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the assignment is active"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the assignment was created"  }
```

## WORKFLOW_INSTANCE
```yaml
table: workflow_instance
comment: "Workflow Instance"
tooltip: "Represents an execution instance of a workflow"
columns:
  workflow_instance_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Instance ID", tooltip: "Unique identifier of the workflow instance"  }
  workflow_id: { type: integer, nullable: false, fk: "workflow.workflow_id", comment: "Workflow ID", tooltip: "Identifier of the workflow being executed", table_display: true  }
  status: { type: varchar(50), comment: "Status", tooltip: "Current status of the workflow instance", form_display: true, table_display: true  }
  current_step_id: { type: integer, comment: "Current Step ID", fk: "workflow_step.workflow_step_id", tooltip: "Identifier of the current step in execution", table_display: true  }
  started_by: { type: integer, comment: "Started By", tooltip: "Identifier of the user who started the workflow"  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the instance is active"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the instance was created"  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the instance was last updated"  }
table_layout:
  default_order: [{field: workflow_instance_id, order: DESC}]
```

## WORKFLOW_INSTANCE_STEP
```yaml
table: workflow_instance_step
comment: "Workflow Instance Step"
tooltip: "Represents the execution of each step within a workflow instance"
columns:
  workflow_instance_step_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Instance Step ID", tooltip: "Unique identifier of the instance step"  }
  workflow_instance_id: { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step being executed", table_display: true  }
  status: { type: varchar(50), comment: "Status", tooltip: "Current status of the step", form_display: true, table_display: true  }
  assigned_to: { type: integer, comment: "Assigned To", tooltip: "Identifier of the user assigned to the step"  }
  started_at: { type: datetime, comment: "Started AT", tooltip: "Date and time when the step execution started"  }
  completed_at: { type: datetime, comment: "Completed AT", tooltip: "Date and time when the step execution was completed"  }
  active: { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the step execution is active"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the record was created"  }
table_layout:
  default_order: [{field: workflow_instance_step_id, order: DESC}]
```

## WORKFLOW_DATA
```yaml
table: workflow_data
comment: "Workflow Data"
tooltip: "Stores values collected during workflow step execution"
columns:
  workflow_data_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Data ID", tooltip: "Unique identifier of the stored value"  }
  workflow_instance_id: { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id: { type: integer, nullable: false, fk: "workflow_step.workflow_step_id", comment: "Workflow Step ID", tooltip: "Identifier of the step where the value was collected", table_display: true  }
  workflow_step_schema_id: { type: integer, nullable: false, fk: "workflow_step_schema.workflow_step_schema_id", comment: "Workflow Step Schema ID", tooltip: "Identifier of the field definition", table_display: true  }
  field: { type: varchar(200), nullable: false, comment: "Field", tooltip: "Technical identifier of the field", table_display: true  }
  value: { type: text, comment: "Value", tooltip: "Value provided for the field"  }
  is_latest: { type: boolean, default: true, comment: "Is Latest", tooltip: "Indicates whether the record is the most recent value for the field"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the value was recorded", table_display: true  }
  updated_at: { type: datetime, comment: "Updated AT", tooltip: "Date and time when the value was last updated"  }
table_layout:
  default_order: [{field: workflow_data_id, order: DESC}]
```

## WORKFLOW_LOG
```yaml
table: workflow_log
comment: "Workflow Log"
tooltip: "Records all actions and state changes during workflow execution"
columns:
  workflow_log_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Log ID", tooltip: "Unique identifier of the log entry"  }
  workflow_instance_id: { type: integer, nullable: false, fk: "workflow_instance.workflow_instance_id", comment: "Workflow Instance ID", tooltip: "Identifier of the workflow instance", table_display: true  }
  workflow_step_id: { type: integer, comment: "Workflow Step ID", fk: "workflow_step.workflow_step_id", tooltip: "Identifier of the step associated with the action", table_display: true  }
  action: { type: varchar(50), nullable: false, comment: "Action", tooltip: "Type of action performed"  }
  status_from: { type: varchar(50), comment: "Status From", tooltip: "Status before the action"  }
  status_to: { type: varchar(50), comment: "Status To", tooltip: "Status after the action", table_display: true  }
  obs: { type: text, comment: "Obs", tooltip: "Additional information describing the action", form_long_text: true  }
  performed_by: { type: integer, comment: "Performed By", tooltip: "Identifier of the user who performed the action"  }
  created_at: { type: datetime, comment: "Created AT", tooltip: "Date and time when the action was recorded", table_display: true  }
table_layout:
  default_order: [{field: workflow_log_id, order: DESC}]
```
