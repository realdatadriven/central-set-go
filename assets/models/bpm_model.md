<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD047 -->
# BPM_MODEL
```yaml
name: BPM
description: Business Process Mapping and Execution Model
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dashboard:
    menu_icon: document-report
    menu_order: 1
    active: true
    menu_config: '{"label": "dashboard", "tooltip": "dashboard_desc", "load_items": {"table": "dashboard", "tables": ["dashboard"]}}'
    tables:
      - dashboard
  Process Management:
    menu_icon: rectangle-group
    menu_order: 2
    active: true
    tables:
      - process_node
      - stakeholder_focal_point
      - process_activity_actor
      - process_activity_schema
      - process_activity_schema_option
      - process_action
      - process_flow
      - process_edge
      - {table: data_type, active: false}
      - {table: input_type, active: false}
      - {table: size, active: false}
  Process Execution:
    menu_icon: clipboard-document-check
    menu_order: 3
    active: true
    menu_config: '{"label": "process", "tooltip": "process_desc", "load_items": {"table": "process_node", "tables": ["process_node"], "filter": {"element_type": "BP", "process_status": "published", "active": true}}}'
    tables:
      - {table: process_instance, requires_rla: true, active: true}
      - process_instance_step
      - process_instance_data
      - process_instance_log
      - process_notification
```

## MODEL NOTES

- `process_node` is the root entity for process definitions. BP, BA, ST, SP, EP, and DP are all rows in this table, distinguished by `element_type`; there are no separate business-process or process-version tables.
- A BP node owns its process definition and child nodes. `parent_node_id` forms the hierarchy: a top-level BP has no parent, a nested BP points to its parent BP, and BA, ST, SP, DP, and EP nodes point to their owning BP.
- BP identity is the combination of `node_code` and `version`. Publishing a changed version creates a new BP node with the same `node_code` and a new `version`; its child nodes and flows belong to that new BP node, leaving prior versions and their running instances unchanged.
- Subforms are conditional on the parent node's `element_type`: BP shows its child nodes and flows, BA shows actors, schema fields, and actions, ST shows focal points, and DP shows its outgoing decision edges. The form configuration declares both the condition and relationship used to load each subform.
- Field-level `form_hide_cond` and `allow_in_subform` conditions are declarative form metadata; the form renderer must apply the condition before displaying fields or loading child tables. Process execution lists only active, published BP nodes.
- `process_flow` stores source flow expressions. `process_edge` stores the expanded node-to-node edges used to render and execute the graph. Conditional DP edges retain their condition and branch value.
- The BA data schema below follows `workflow_step_schema` from `workflow_model.md`. Its `workflow_step_id` reference is represented by `activity_node_id`; all other field properties retain the workflow schema's meaning.
- `data_type`, `input_type`, `size`, and `dashboard` are shared lookup/dashboard tables from the existing application model, not new BPM-specific definitions.
- Decision conditions are data, not executable code. Define and validate a restricted condition grammar in the execution layer before evaluating them.

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

## PROCESS_NODE
```yaml
table: process_node
comment: "Process Nodes"
tooltip: "A BP, BA, ST, SP, EP, or DP element and root record for BPM definitions"
columns:
  process_node_id:   { type: integer, pk: true, autoincrement: true, comment: "Process Node ID", tooltip: "Unique identifier of the process element" }
  parent_node_id:    { type: integer, fk: "process_node.process_node_id", comment: "Parent Node ID", tooltip: "Containing BP node; null for a top-level BP", form_display: true, table_display: true, form_size: 6 }
  node_code:         { type: varchar, len: 50, nullable: false, comment: "Node Code", tooltip: "Stable element identifier, for example BP001, BA001, or ST001", form_display: true, table_display: true, form_size: 3 }
  element_type:      { type: varchar, len: 2, nullable: false, comment: "Element Type", tooltip: "Element prefix: BP, BA, ST, SP, EP, or DP", form_display: true, table_display: true, form_size: 3 }
  version:           { type: varchar, len: 50, nullable: true, comment: "Version", tooltip: "Definition version; required for BP nodes and copied forward when publishing a change", form_display: true, table_display: true, form_size: 3, form_hide_cond: 'data.element_type!=="BP"' }
  name:              { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Display name of the process element", form_display: true, table_display: true, form_size: 6 }
  description:       { type: text, comment: "Description", tooltip: "Description of the process element", form_display: true, form_long_text: true, form_size: 12 }
  order_index:       { type: integer, comment: "Order Index", tooltip: "Display order within the parent process", form_display: true, table_display: true, form_size: 3 }
  process_status:    { type: varchar, len: 30, default: "draft", comment: "Process Status", tooltip: "BP lifecycle status: draft, published, or retired", form_display: true, table_display: true, form_size: 3, form_hide_cond: 'data.element_type!=="BP"' }
  steps_orientation: { type: varchar, len: 20, comment: "Steps Orientation", tooltip: "Preferred horizontal or vertical diagram orientation for a BP", form_display: true, form_size: 4, form_hide_cond: 'data.element_type!=="BP"' }
  definition_source: { type: text, comment: "Definition Source", tooltip: "Markdown/YAML source for a BP definition", form_display: true, form_code: markdown, form_long_text: true, form_size: 12, form_hide_cond: 'data.element_type!=="BP"' }
  execution_mode:    { type: varchar, len: 20, default: "manual", comment: "Execution Mode", tooltip: "Manual, automated, or mixed execution for a BA", form_display: true, table_display: true, form_size: 4, form_hide_cond: 'data.element_type!=="BA"' }
  schedule:          { type: varchar, len: 200, comment: "Schedule", tooltip: "Optional cron schedule for an automated BA", form_display: true, table_display: true, form_size: 4, form_hide_cond: 'data.element_type!=="BA"' }
  subprocess_node_id: { type: integer, fk: "process_node.process_node_id", comment: "Sub-process Node ID", tooltip: "BP node launched by this activity, when applicable", form_display: true, table_display: true, form_size: 6, form_hide_cond: 'data.element_type!=="BA"' }
  layout_x:          { type: integer, comment: "Layout X", tooltip: "Optional diagram position on the horizontal axis" }
  layout_y:          { type: integer, comment: "Layout Y", tooltip: "Optional diagram position on the vertical axis" }
  attachment:        { type: text, comment: "Attachments", tooltip: "References to documents related to this element", form_display: true, form_long_text: true, form_size: 12 }
  active:            { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the element is active", form_display: true, table_display: true, form_size: 3 }
  user_id:           { type: integer, comment: "User ID", tooltip: "User responsible for this element" }
  app_id:            { type: integer, comment: "App ID", tooltip: "Application context" }
  created_at:        { type: datetime, comment: "Created At", tooltip: "When the element was created" }
  updated_at:        { type: datetime, comment: "Updated At", tooltip: "When the element was last updated" }
  excluded:          { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the element is excluded from active use" }
unique_together:
  - [node_code, version]
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  allow_in_subform:
    process_flow: {show: true, when: 'data.element_type==="BP"'}
    process_activity_actor: {show: true, when: 'data.element_type==="BA"'}
    process_activity_schema: {show: true, when: 'data.element_type==="BA"'}
    process_action: {show: true, when: 'data.element_type==="BA"'}
    stakeholder_focal_point: {show: true, when: 'data.element_type==="ST"'}
    process_edge: {show: true, when: 'data.element_type==="DP"'}
  tabs_steps_conf:
    - {label: Node, fields: [parent_node_id, node_code, element_type, version, name, order_index, active, description]}
    - {label: Diagram, fields: [layout_x, layout_y, attachment]}
table_layout:
  default_order: [{field: node_code, order: ASC}]
```

## STAKEHOLDER_FOCAL_POINT
```yaml
table: stakeholder_focal_point
comment: "Stakeholder Focal Point"
tooltip: "Contact person for a stakeholder element, optionally linked to an application user"
columns:
  stakeholder_focal_point_id: { type: integer, pk: true, autoincrement: true, comment: "Focal Point ID", tooltip: "Unique identifier of the stakeholder contact" }
  stakeholder_node_id:        { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Stakeholder Node ID", tooltip: "ST element this contact represents", form_display: true, table_display: true, form_size: 6 }
  name:                       { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the focal point", form_display: true, table_display: true, form_size: 6 }
  email:                      { type: varchar, len: 255, nullable: false, comment: "Email", tooltip: "Contact email; match the app user email to enable user notifications", form_display: true, table_display: true, form_size: 6 }
  role:                       { type: varchar, len: 100, comment: "Role", tooltip: "Role of the focal point", form_display: true, table_display: true, form_size: 6 }
  user_id:                    { type: integer, comment: "User ID", tooltip: "Linked application user, when available", form_display: true, table_display: true, form_size: 6 }
  active:                     { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether this focal point is active", form_display: true, table_display: true, form_size: 3 }
  created_at:                 { type: datetime, comment: "Created At", tooltip: "When the focal point was created" }
  updated_at:                 { type: datetime, comment: "Updated At", tooltip: "When the focal point was last updated" }
  excluded:                   { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the focal point is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 8
table_layout:
  default_order: [{field: name, order: ASC}]
```

## PROCESS_ACTIVITY_ACTOR
```yaml
table: process_activity_actor
comment: "Process Activity Actor"
tooltip: "Associates an activity with the stakeholder responsible for performing it"
columns:
  process_activity_actor_id: { type: integer, pk: true, autoincrement: true, comment: "Activity Actor ID", tooltip: "Unique identifier of the activity assignment" }
  activity_node_id:          { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Activity Node ID", tooltip: "BA element assigned to the stakeholder", form_display: true, table_display: true, form_size: 6 }
  stakeholder_node_id:       { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Stakeholder Node ID", tooltip: "ST element performing the activity", form_display: true, table_display: true, form_size: 6 }
  responsibility:            { type: varchar, len: 50, default: "assignee", comment: "Responsibility", tooltip: "Assignment responsibility, such as assignee or reviewer", form_display: true, table_display: true, form_size: 6 }
  active:                    { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether this assignment is active", form_display: true, table_display: true, form_size: 3 }
  created_at:                { type: datetime, comment: "Created At", tooltip: "When the assignment was created" }
  updated_at:                { type: datetime, comment: "Updated At", tooltip: "When the assignment was last updated" }
  excluded:                  { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the assignment is excluded from active use" }
table_layout:
  default_order: [{field: process_activity_actor_id, order: DESC}]
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

## PROCESS_ACTIVITY_SCHEMA
```yaml
table: process_activity_schema
comment: "Process Activity Schema"
tooltip: "Defines the data fields collected for a BA activity; mirrors workflow_step_schema"
columns:
  process_activity_schema_id: { type: integer, pk: true, autoincrement: true, comment: "Workflow Step Schema ID", tooltip: "Unique identifier of the schema field" }
  activity_node_id:           { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Activity Node ID", tooltip: "Identifier of the BA activity where the field is collected", form_display: true, table_display: true, form_size: 6 }
  field:                      { type: varchar, len: 200, nullable: false, comment: "Field", tooltip: "Technical identifier of the field", form_display: true, table_display: true, form_size: 4 }
  label:                      { type: varchar, len: 200, nullable: false, comment: "Label", tooltip: "Display name of the field", form_display: true, table_display: true, form_size: 4 }
  data_type:                  { type: varchar, fk: "data_type.data_type", nullable: false, comment: "Data Type", tooltip: "Type of data stored in the field", form_display: true, table_display: true, form_size: 4 }
  nullable:                   { type: boolean, default: true, comment: "Nullable", tooltip: "Indicates whether the field can be empty", form_display: true, table_display: true, form_size: 3 }
  default_value:              { type: varchar, len: 200, comment: "Default Value", tooltip: "Default value assigned to the field", form_display: true, table_display: false, form_size: 3 }
  validation_rule:            { type: varchar, len: 200, comment: "Validation Rule", tooltip: "Regex validation rule for the field", form_display: true, table_display: false, form_size: 3 }
  order_index:                { type: integer, comment: "Order Index", tooltip: "Position of the field within the step", form_display: true, table_display: true, form_size: 3 }
  format:                     { type: varchar, len: 200, comment: "Format", tooltip: "Format intl.Format", form_display: true, form_size: 4 }
  size:                       { type: integer, fk: "size.size", comment: "Size", tooltip: "1 - 12 size that will be shown in form", form_display: true, form_size: 2 }
  elipsis:                    { type: integer, comment: "Elipsis", tooltip: "Text elipsis", form_display: true, form_size: 2 }
  input_type:                 { type: varchar, fk: "input_type.input_type", comment: "Options Input Type", tooltip: "Combobox,Checkbox or Radio", form_display: true, table_display: true, form_size: 2 }
  active:                     { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the field is active", form_display: true, table_display: true, form_size: 2 }
  options:                    { type: text, comment: "Options", tooltip: "JSON Array of string or array of objects{label,value}", form_display: true, form_long_text: true, form_size: 12 }
  user_id:                    { type: integer, comment: "User ID", tooltip: "Identifier of the user responsible for the field definition" }
  app_id:                     { type: integer, comment: "App ID", tooltip: "Identifier of the application context" }
  created_at:                 { type: datetime, comment: "Created AT", tooltip: "Date and time when the field was created" }
  updated_at:                 { type: datetime, comment: "Updated AT", tooltip: "Date and time when the field was last updated" }
  excluded:                   { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the field is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 6
  allow_in_subform: {process_activity_schema_option: true}
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## PROCESS_ACTIVITY_SCHEMA_OPTION
```yaml
table: process_activity_schema_option
comment: "Process Activity Schema Option"
tooltip: "Stores an option for a BA schema field when options are managed as rows"
columns:
  process_activity_schema_option_id: { type: integer, pk: true, autoincrement: true, comment: "Schema Option ID", tooltip: "Unique identifier of the field option" }
  process_activity_schema_id:        { type: integer, nullable: false, fk: "process_activity_schema.process_activity_schema_id", comment: "Schema ID", tooltip: "Field definition that owns this option", form_display: true, table_display: true, form_size: 6 }
  label:                             { type: varchar, len: 200, nullable: false, comment: "Label", tooltip: "Option display label", form_display: true, table_display: true, form_size: 6 }
  value:                             { type: varchar, len: 200, nullable: false, comment: "Value", tooltip: "Value stored when the option is selected", form_display: true, table_display: true, form_size: 6 }
  order_index:                       { type: integer, comment: "Order Index", tooltip: "Position of the option", form_display: true, table_display: true, form_size: 3 }
  active:                            { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the option is available", form_display: true, table_display: true, form_size: 3 }
  created_at:                        { type: datetime, comment: "Created At", tooltip: "When the option was created" }
  updated_at:                        { type: datetime, comment: "Updated At", tooltip: "When the option was last updated" }
  excluded:                          { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the option is excluded from active use" }
table_layout:
  default_order: [{field: order_index, order: ASC}]
```

## PROCESS_ACTION
```yaml
table: process_action
comment: "Process Action"
tooltip: "Defines an API or ETL action available from a BA activity"
columns:
  process_action_id: { type: integer, pk: true, autoincrement: true, comment: "Process Action ID", tooltip: "Unique identifier of the action" }
  activity_node_id:  { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Activity Node ID", tooltip: "BA activity that exposes the action", form_display: true, table_display: true, form_size: 6 }
  name:              { type: varchar, len: 100, nullable: false, comment: "Name", tooltip: "Technical name of the action", form_display: true, table_display: true, form_size: 4 }
  label:             { type: varchar, len: 200, nullable: false, comment: "Label", tooltip: "Button label shown to the actor", form_display: true, table_display: true, form_size: 4 }
  action_type:       { type: varchar, len: 30, nullable: false, comment: "Action Type", tooltip: "API call or ETL run", form_display: true, table_display: true, form_size: 4 }
  endpoint:          { type: varchar, len: 500, comment: "Endpoint", tooltip: "Backend API or ETL identifier", form_display: true, table_display: true, form_size: 8 }
  configuration:     { type: text, comment: "Configuration", tooltip: "Action parameters and execution settings as JSON", form_display: true, form_long_text: true, form_size: 12 }
  active:            { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the action is available", form_display: true, table_display: true, form_size: 3 }
  user_id:           { type: integer, comment: "User ID", tooltip: "User responsible for the action definition" }
  app_id:            { type: integer, comment: "App ID", tooltip: "Application context" }
  created_at:        { type: datetime, comment: "Created At", tooltip: "When the action was created" }
  updated_at:        { type: datetime, comment: "Updated At", tooltip: "When the action was last updated" }
  excluded:          { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the action is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 8
table_layout:
  default_order: [{field: process_action_id, order: DESC}]
```

## PROCESS_FLOW
```yaml
table: process_flow
comment: "Process Flow"
tooltip: "Stores source flow expressions before expansion into executable graph edges"
columns:
  process_flow_id:  { type: integer, pk: true, autoincrement: true, comment: "Process Flow ID", tooltip: "Unique identifier of the source flow" }
  process_node_id:  { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Business Process Node ID", tooltip: "BP node that owns this flow", form_display: true, table_display: true, form_size: 6 }
  flow_code:        { type: varchar, len: 50, nullable: false, comment: "Flow Code", tooltip: "Stable flow identifier, for example FL001", form_display: true, table_display: true, form_size: 3 }
  expression:       { type: text, nullable: false, comment: "Expression", tooltip: "Source sequence using -> and decision expressions", form_display: true, form_long_text: true, form_code: yaml, form_size: 12 }
  active:           { type: boolean, default: true, comment: "Active", tooltip: "Indicates whether the flow is active", form_display: true, table_display: true, form_size: 3 }
  created_at:       { type: datetime, comment: "Created At", tooltip: "When the flow was created" }
  updated_at:       { type: datetime, comment: "Updated At", tooltip: "When the flow was last updated" }
  excluded:         { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the flow is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: true
  size: 8
table_layout:
  default_order: [{field: flow_code, order: ASC}]
```

## PROCESS_EDGE
```yaml
table: process_edge
comment: "Process Edge"
tooltip: "Expanded connection between process nodes; conditional edges represent DP branches"
columns:
  process_edge_id: { type: integer, pk: true, autoincrement: true, comment: "Process Edge ID", tooltip: "Unique identifier of the graph edge" }
  process_flow_id: { type: integer, nullable: false, fk: "process_flow.process_flow_id", comment: "Process Flow ID", tooltip: "Source flow expression", form_display: true, table_display: true, form_size: 6 }
  from_node_id:    { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "From Node ID", tooltip: "Source process element", form_display: true, table_display: true, form_size: 6 }
  to_node_id:      { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "To Node ID", tooltip: "Destination process element", form_display: true, table_display: true, form_size: 6 }
  edge_order:      { type: integer, comment: "Edge Order", tooltip: "Order of this connection in the source flow", form_display: true, table_display: true, form_size: 3 }
  is_decision:     { type: boolean, default: false, comment: "Is Decision", tooltip: "Indicates that this is a conditional decision branch", form_display: true, table_display: true, form_size: 3 }
  condition:       { type: varchar, len: 500, comment: "Condition", tooltip: "Restricted expression evaluated by the decision engine", form_display: true, table_display: false, form_size: 8 }
  branch_value:    { type: boolean, comment: "Branch Value", tooltip: "Expected condition result for this branch", form_display: true, table_display: true, form_size: 4 }
  created_at:      { type: datetime, comment: "Created At", tooltip: "When the edge was generated" }
  updated_at:      { type: datetime, comment: "Updated At", tooltip: "When the edge was last updated" }
  excluded:        { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the edge is excluded from active use" }
table_layout:
  default_order: [{field: edge_order, order: ASC}]
```

## PROCESS_INSTANCE
```yaml
table: process_instance
comment: "Process Instance"
tooltip: "A running or completed execution of a published BP node"
columns:
  process_instance_id: { type: integer, pk: true, autoincrement: true, comment: "Process Instance ID", tooltip: "Unique identifier of this process execution" }
  process_node_id:     { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Business Process Node ID", tooltip: "BP node being executed", form_display: true, table_display: true, form_size: 6 }
  definition_snapshot: { type: text, comment: "Definition Snapshot", tooltip: "Definition state captured when execution starts, preserving its audit context without a separate version record" }
  parent_instance_id:  { type: integer, fk: "process_instance.process_instance_id", comment: "Parent Instance ID", tooltip: "Parent execution when this instance is a sub-process", form_display: true, table_display: true, form_size: 6 }
  parent_instance_step_id: { type: integer, _fk: "process_instance_step.process_instance_step_id", comment: "Parent Instance Step ID", tooltip: "Parent activity that launched this sub-process", form_display: true, table_display: true, form_size: 6 }
  instance_code:       { type: varchar, len: 100, unique: true, comment: "Instance Code", tooltip: "Human-readable execution reference", form_display: true, table_display: true, form_size: 6 }
  status:              { type: varchar, len: 30, default: "pending", nullable: false, comment: "Status", tooltip: "Pending, in progress, completed, cancelled, or failed", form_display: true, table_display: true, form_size: 4 }
  current_node_id:     { type: integer, fk: "process_node.process_node_id", comment: "Current Node ID", tooltip: "Current process element, when execution is active", form_display: true, table_display: true, form_size: 4 }
  current_user_id:     { type: integer, comment: "Current User ID", tooltip: "Application user currently responsible for the next action", form_display: true, table_display: true, form_size: 4 }
  started_by_user_id:  { type: integer, comment: "Started By User ID", tooltip: "Application user who started the process", form_display: true, table_display: true, form_size: 6 }
  started_at:          { type: datetime, comment: "Started At", tooltip: "When execution started", form_display: true, table_display: true, form_size: 6 }
  completed_at:        { type: datetime, comment: "Completed At", tooltip: "When execution completed", form_display: true, table_display: true, form_size: 6 }
  created_at:          { type: datetime, comment: "Created At", tooltip: "When the instance record was created" }
  updated_at:          { type: datetime, comment: "Updated At", tooltip: "When the instance was last updated" }
  excluded:            { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the instance is excluded from active use" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order: [{field: started_at, order: DESC}]
```

## PROCESS_INSTANCE_STEP
```yaml
table: process_instance_step
comment: "Process Instance Step"
tooltip: "Tracks each activity or process element reached during an instance"
columns:
  process_instance_step_id: { type: integer, pk: true, autoincrement: true, comment: "Instance Step ID", tooltip: "Unique identifier of this execution step" }
  process_instance_id:      { type: integer, nullable: false, fk: "process_instance.process_instance_id", comment: "Process Instance ID", tooltip: "Execution containing this step", form_display: true, table_display: true, form_size: 6 }
  process_node_id:          { type: integer, nullable: false, fk: "process_node.process_node_id", comment: "Process Node ID", tooltip: "Definition element reached by this step", form_display: true, table_display: true, form_size: 6 }
  status:                   { type: varchar, len: 30, default: "pending", nullable: false, comment: "Status", tooltip: "Pending, in progress, completed, sent back, skipped, or failed", form_display: true, table_display: true, form_size: 4 }
  assigned_user_id:         { type: integer, comment: "Assigned User ID", tooltip: "User responsible for this step", form_display: true, table_display: true, form_size: 4 }
  started_at:               { type: datetime, comment: "Started At", tooltip: "When work on the step started", form_display: true, table_display: true, form_size: 4 }
  completed_at:             { type: datetime, comment: "Completed At", tooltip: "When the step was completed", form_display: true, table_display: true, form_size: 4 }
  sent_back_to_step_id:     { type: integer, fk: "process_instance_step.process_instance_step_id", comment: "Sent Back To Step ID", tooltip: "Earlier execution step when this step is returned for rework", form_display: true, table_display: true, form_size: 6 }
  created_at:               { type: datetime, comment: "Created At", tooltip: "When the execution step was created" }
  updated_at:               { type: datetime, comment: "Updated At", tooltip: "When the execution step was last updated" }
  excluded:                 { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the execution step is excluded from active use" }
table_layout:
  default_order: [{field: process_instance_step_id, order: ASC}]
```

## PROCESS_INSTANCE_DATA
```yaml
table: process_instance_data
comment: "Process Instance Data"
tooltip: "Stores values collected for schema fields during process execution"
columns:
  process_instance_data_id: { type: integer, pk: true, autoincrement: true, comment: "Instance Data ID", tooltip: "Unique identifier of the collected value" }
  process_instance_id:      { type: integer, nullable: false, fk: "process_instance.process_instance_id", comment: "Process Instance ID", tooltip: "Execution this value belongs to", form_display: true, table_display: true, form_size: 6 }
  process_instance_step_id: { type: integer, nullable: false, fk: "process_instance_step.process_instance_step_id", comment: "Instance Step ID", tooltip: "Step where the value was collected", form_display: true, table_display: true, form_size: 6 }
  process_activity_schema_id: { type: integer, nullable: false, fk: "process_activity_schema.process_activity_schema_id", comment: "Activity Schema ID", tooltip: "Field definition used to collect the value", form_display: true, table_display: true, form_size: 6 }
  collected_by_user_id:     { type: integer, comment: "Collected By User ID", tooltip: "Application user who entered or changed the value", form_display: true, table_display: true, form_size: 6 }
  value:                    { type: text, comment: "Value", tooltip: "Collected field value, serialized according to its data type", form_display: true, form_long_text: true, form_size: 12 }
  created_at:               { type: datetime, comment: "Created At", tooltip: "When the value was collected" }
  updated_at:               { type: datetime, comment: "Updated At", tooltip: "When the value was last updated" }
  excluded:                 { type: boolean, default: false, comment: "Excluded", tooltip: "Indicates whether the value is excluded from active use" }
table_layout:
  default_order: [{field: process_instance_data_id, order: ASC}]
```

## PROCESS_INSTANCE_LOG
```yaml
table: process_instance_log
comment: "Process Instance Log"
tooltip: "Audit trail for process and activity execution"
columns:
  process_instance_log_id: { type: integer, pk: true, autoincrement: true, comment: "Instance Log ID", tooltip: "Unique identifier of the audit event" }
  process_instance_id:     { type: integer, nullable: false, fk: "process_instance.process_instance_id", comment: "Process Instance ID", tooltip: "Execution associated with this event", form_display: true, table_display: true, form_size: 6 }
  process_instance_step_id: { type: integer, fk: "process_instance_step.process_instance_step_id", comment: "Instance Step ID", tooltip: "Execution step associated with this event", form_display: true, table_display: true, form_size: 6 }
  event_type:              { type: varchar, len: 50, nullable: false, comment: "Event Type", tooltip: "Event such as started, assigned, completed, sent back, action run, or decision taken", form_display: true, table_display: true, form_size: 6 }
  details:                 { type: text, comment: "Details", tooltip: "Event context as text or JSON", form_display: true, form_long_text: true, form_size: 12 }
  actor_user_id:           { type: integer, comment: "Actor User ID", tooltip: "Application user who caused the event", form_display: true, table_display: true, form_size: 6 }
  created_at:              { type: datetime, comment: "Created At", tooltip: "When the event occurred", form_display: true, table_display: true, form_size: 6 }
table_layout:
  default_order: [{field: created_at, order: DESC}]
```

## PROCESS_NOTIFICATION
```yaml
table: process_notification
comment: "Process Notification"
tooltip: "Tracks notifications sent to focal points and assigned application users"
columns:
  process_notification_id: { type: integer, pk: true, autoincrement: true, comment: "Notification ID", tooltip: "Unique identifier of the notification" }
  process_instance_id:    { type: integer, nullable: false, fk: "process_instance.process_instance_id", comment: "Process Instance ID", tooltip: "Execution associated with the notification", form_display: true, table_display: true, form_size: 6 }
  process_instance_step_id: { type: integer, fk: "process_instance_step.process_instance_step_id", comment: "Instance Step ID", tooltip: "Step that triggered the notification", form_display: true, table_display: true, form_size: 6 }
  recipient_email:        { type: varchar, len: 255, nullable: false, comment: "Recipient Email", tooltip: "Recipient contact address", form_display: true, table_display: true, form_size: 6 }
  recipient_user_id:      { type: integer, comment: "Recipient User ID", tooltip: "Linked application user, when available", form_display: true, table_display: true, form_size: 6 }
  channel:                { type: varchar, len: 30, default: "email", comment: "Channel", tooltip: "Notification delivery channel", form_display: true, table_display: true, form_size: 4 }
  status:                 { type: varchar, len: 30, default: "pending", comment: "Status", tooltip: "Pending, sent, or failed", form_display: true, table_display: true, form_size: 4 }
  sent_at:                { type: datetime, comment: "Sent At", tooltip: "When delivery was attempted", form_display: true, table_display: true, form_size: 4 }
  error:                  { type: text, comment: "Error", tooltip: "Delivery error details, if sending failed", form_display: true, form_long_text: true, form_size: 12 }
  created_at:             { type: datetime, comment: "Created At", tooltip: "When the notification was queued" }
  updated_at:             { type: datetime, comment: "Updated At", tooltip: "When the notification record was last updated" }
table_layout:
  default_order: [{field: created_at, order: DESC}]
```
