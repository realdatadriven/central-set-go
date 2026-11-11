<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD047 -->

# UI_FOR_BPM
```yaml
name: UI_FOR_BPM
description: Dynamic UI dashboard for published BPM processes and process instances
database: UI
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## BPM_UI
```yaml
table: ui
description: Add the BPM process dashboard
cond: 'WHERE ui_slug = :ui_slug AND excluded = false'
data:
  ui_slug: bpm
  ui_name: BPM Process Dashboard
  ui_desc: Start published processes and view their activities, process maps, and execution history
  default_locale: en
  active: true
  children:
    - table: ui_partial
      description: Add the process instance results partial
      cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
      data:
        ui_id: ui_id()
        ui_partial: instance_results
        ui_partial_desc: Search and filter process instances with HTMX
        partial_template: FileContent(assets/bpm_ui/instance_results.html)
        active: true
    - table: ui_partial
      description: Add the process instance detail partial
      cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
      data:
        ui_id: ui_id()
        ui_partial: instance_detail
        ui_partial_desc: Process instance steps, activity data, process map, and subprocesses
        partial_template: FileContent(assets/bpm_ui/instance_detail.html)
        active: true
    - table: ui_page
      description: Add the BPM dashboard page
      cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
      data:
        ui_id: ui_id()
        page_key: dashboard
        page_title: BPM Process Dashboard
        meta_description: Start and monitor business process instances
        page_template: FileContent(assets/bpm_ui/dashboard.html)
        cache_seconds: 0
        default_page: true
        login_required: true
        active: true
    - table: ui_page
      description: Add process instance search and detail page
      cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
      data:
        ui_id: ui_id()
        page_key: instances
        page_title: Process Instances
        meta_description: Search process instances and inspect activity history and values
        page_template: FileContent(assets/bpm_ui/instances.html)
        cache_seconds: 0
        default_page: false
        login_required: true
        active: true
```

# BPM_EMAIL_ACTIONS
```yaml
name: BPM_EMAIL_ACTIONS
description: Email notifications for BPM activity step starts and completions
database: ADMIN
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## BPM_BA_STEP_COMPLETED_EMAIL
```yaml
table: crud_action
description: Email the assigned user when a BA activity step is completed
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: BPM_BA_STEP_COMPLETED_EMAIL
  crud_action: Notify assignee when BPM activity completes
  action_type_id: 2
  err_msg: 'Could not send the BPM activity completion email.'
  table: process_instance_step
  db: BPM
  active: true
  create: false
  read: false
  update: true
  delete: false
  user_trigger: false
  parallel: false
  sql_condition: |
    SELECT '{{.process_node.element_type}}' = 'BA'
      AND lower('{{replace (toString .status) "'" "''"}}') = 'completed'
      AND '{{replace (default "" .assignee.email) "'" "''"}}' <> '' AS cond
  email_subject: 'Completed activity: {{.process_node.name}}'
  email_to: '{{.assignee.email}}'
  email_template: |
    <p>Hello {{.assignee.username}},</p>
    <p>The activity <strong>{{.process_node.name}}</strong> has been completed.</p>
    <p>Process instance step: {{.process_instance_step_id}}</p>
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      - action_data: process_node
        action_data_desc: Activity node details for the BPM email
        action_data_type_id: 3
        crud_action_id: crud_action_id()
        odata_path: 'BPM/process_node?$filter=process_node_id eq {{.process_node_id}}'
        single_row_obj: true
        active: true
      - action_data: assignee
        action_data_desc: Assigned user details for the BPM email
        action_data_type_id: 3
        crud_action_id: crud_action_id()
        odata_path: 'ADMIN/users?$filter=user_id eq {{.assigned_user_id}}'
        single_row_obj: true
        active: true
```

## BPM_BA_STEP_STARTED_EMAIL
```yaml
table: crud_action
description: Email the assigned user when a BA activity step begins
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: BPM_BA_STEP_STARTED_EMAIL
  crud_action: Notify assignee when BPM activity begins
  action_type_id: 2
  err_msg: 'Could not send the BPM activity start email.'
  table: process_instance_step
  db: BPM
  active: true
  create: true
  read: false
  update: true
  delete: false
  user_trigger: false
  parallel: false
  sql_condition: |
    SELECT '{{.process_node.element_type}}' = 'BA'
      AND lower('{{replace (toString .status) "'" "''"}}') = 'in_progress'
      AND '{{replace (default "" .assignee.email) "'" "''"}}' <> '' AS cond
  email_subject: 'Activity started: {{.process_node.name}}'
  email_to: '{{.assignee.email}}'
  email_template: |
    <p>Hello {{.assignee.username}},</p>
    <p>The activity <strong>{{.process_node.name}}</strong> has started and is assigned to you.</p>
    <p>Process instance step: {{.process_instance_step_id}}</p>
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      - action_data: process_node
        action_data_desc: Activity node details for the BPM email
        action_data_type_id: 3
        crud_action_id: crud_action_id()
        odata_path: 'BPM/process_node?$filter=process_node_id eq {{.process_node_id}}'
        single_row_obj: true
        active: true
      - action_data: assignee
        action_data_desc: Assigned user details for the BPM email
        action_data_type_id: 3
        crud_action_id: crud_action_id()
        odata_path: 'ADMIN/users?$filter=user_id eq {{.assigned_user_id}}'
        single_row_obj: true
        active: true
```

The completion rule sends on a qualifying `process_instance_step` update; the start rule runs on create or update so it also catches a step created directly as `in_progress`. Both require the record to include `process_node_id` and `assigned_user_id`, the referenced node to be a BA, and a resolvable assignee email. These are CRUD-operation triggers rather than status-transition detectors: if an already completed or in-progress step is updated again, its matching email can be sent again. Marking a step complete does not itself evaluate edges, create the next step, or start a subprocess; the BPM execution layer must perform that transition and set the next BA step to `in_progress` for its notification to run.
