package main

import "fmt"

/*## PROCESS_TYPE
```yaml
table: process_type
comment: Proccess
columns:
  process_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Proccess ID" }
  process_type:      { type: varchar(20), nullable: false, unique: true, comment: "Proccess", form_display: true, table_display: true, order: 1 }
  process_type_desc: { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 2 }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
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
  batch_process:       { type: varchar(200), nullable: false, comment: "Process", form_display: true, table_display: true, order: 2, form_size: 9 }
  batch_process_code:  { type: varchar(200), nullable: false, comment: "Code", form_display: true, table_display: true, order: 1, form_size: 2 }
  batch_process_desc:  { type: text, comment: "Description", form_display: true, table_display: true, order: 1, form_size: 2 }
  cron:                { type: varchar(200), nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4 }
  batch_process_order: { type: integer, comment: "Proccess ID", order: 3, form_size: 2 }
  process_type_id:     { type: integer, fk: "process_type.process_type_id", comment: "Proccess ID", order: 3, form_size: 2 }
  err_msg:             { type: varchar(200), nullable: false, comment: "Error Message", form_display: true, table_display: true, order: 4 }
  db:                  { type: varchar(200), nullable: false, comment: "Table", form_display: true, table_display: true, order: 4 }
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
#data:
#  - {batch_process_id: 1, batch_process: Validate user Email existance, batch_process_code: USR01, process_type_id: 2, err_msg: "User {{.email}} already exists!", table: users, db: ADMIN, sql: "select * from users where email = :email", app_id: 1, create: true, user_id: 1}
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
  batch_process_code:   { type: varchar(200), comment: "Proccess Code", order: 2 }
  batch_process:        { type: varchar(200), comment: "Proccess Name", order: 3 }
  db:                   { type: varchar(200), comment: "Database", order: 5 }
  process_type:         { type: varchar(20), comment: "Action Type", order: 7 }
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
```*/

func (app *application) batch_processs(params Dict) Dict {
	valParams := []any{}
	batch_processs_sql := fmt.Sprintf(`SELECT * FROM "batch_process" WHERE "active" = TRUE`)
	batch_process_rows, err := app.AdminGetRowsByFilter(batch_processs_sql, valParams)
	if err != nil {
		fmt.Printf("Error occurred while fetching batch_processs: %v", err)
	} else if len(batch_process_rows) > 0 {

	}
	return Dict{}
}

/*ALL this can be solved with cron and etlx*/
