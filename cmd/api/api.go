package main

/** -- SCHEMA DEFINITIONS --
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

## API
```yaml
table: api
comment: API Integrations
columns:
  api_id:                { type: integer, pk: true, autoincrement: true, comment: "API ID" }
  api_name:              { type: varchar, len: 100, nullable: false, comment: "API Name", form_display: true, table_display: true, form_size: 6, order: 1 }
  api_type_id:           { type: integer, fk: "api_type.api_type_id", nullable: false, comment: "API Type", form_display: true, table_display: true, form_size: 3, order: 2 }
  api_description:       { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 3 }
  endpoint:              { type: varchar, len: 500, nullable: false, comment: "API Endpoint", form_display: true, table_display: true, form_size: 9, order: 4 }
  request_body_template: { type: text, comment: "Request Body Template (Go template)", form_display: true, form_long_text: true, form_code: json, order: 5 }
  num_retries:           { type: integer, default: 3, comment: "Number of Retries", form_display: true, table_display: true, form_size: 3, order: 6 }
  timeout_seconds:       { type: integer, default: 30, comment: "Timeout (seconds)", form_display: true, table_display: true, form_size: 3, order: 7 }
  headers_template:      { type: text, comment: "Headers Template (use @VAR_NAME for environment variables)", form_display: true, form_long_text: true, form_code: json, order: 8 }
  active:                { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:               { type: integer, fk: "users.user_id", comment: "Created by", order: 10 }
  app_id:                { type: integer, fk: "app.app_id", comment: "App ID", order: 11 }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {api_header: true, api_call_log: true}
  sub_form_size: 9
```

## API_HEADER
```yaml
table: api_header
comment: API Headers
columns:
  api_header_id:   { type: integer, pk: true, autoincrement: true, comment: "Header ID" }
  api_id:          { type: integer, fk: "api.api_id", nullable: false, comment: "API", form_display: true, table_display: true, form_size: 4, order: 1 }
  header_name:     { type: varchar, len: 100, nullable: false, comment: "Header Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  header_value:    { type: text, nullable: false, comment: "Header Value (supports @VAR_NAME for env variables)", form_display: true, form_long_text: true, table_display: true, form_size: 8, order: 3 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:         { type: integer, fk: "users.user_id", comment: "Created by", order: 10 }
  app_id:          { type: integer, fk: "app.app_id", comment: "App ID", order: 11 }
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
  api_call_log_id:  { type: integer, pk: true, autoincrement: true, comment: "Log ID" }
  api_id:           { type: integer, fk: "api.api_id", nullable: false, comment: "API", form_display: true, table_display: true, form_size: 4, order: 1 }
  api_name:         { type: varchar, len: 100, comment: "API Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  request_datetime:  { type: datetime, nullable: false, comment: "Request DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", form_use_label: true, order: 3 }
  response_datetime: { type: datetime, comment: "Response DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", form_use_label: true, order: 4 }
  request_body:     { type: text, comment: "Request Body", form_display: true, form_long_text: true, form_code: json, order: 5 }
  response_body:    { type: text, comment: "Response Body", form_display: true, form_long_text: true, form_code: json, order: 6 }
  response_status:  { type: integer, comment: "Response Status Code", form_display: true, table_display: true, order: 7 }
  response_message: { type: varchar, len: 500, comment: "Response Message", form_display: true, table_display: true, order: 8 }
  user_id:          { type: integer, fk: "users.user_id", comment: "User ID", order: 9 }
  app_id:           { type: integer, fk: "app.app_id", comment: "App ID", order: 10 }
  created_at:       { type: datetime, comment: "Created at", order: 11 }
  updated_at:       { type: datetime, comment: "Updated at", order: 12 }
  excluded:         { type: boolean, default: false, comment: "Excluded", order: 13 }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
```
**/

func (app *application) getAPI(params Dict) Dict {
	table := "api"
	if _, ok := params["table"].(string); ok {
		table, _ = params["table"].(string)
	} else if _, ok := params["data"].(Dict)["table"].(string); ok {
		table, _ = params["data"].(Dict)["table"].(string)
	}
	database := "ETLX"
	if _, ok := params["db"].(string); ok {
		database, _ = params["db"].(string)
	} else if _, ok := params["data"].(Dict)["db"].(string); ok {
		database, _ = params["data"].(Dict)["db"].(string)
	} else if _, ok := params["database"].(string); ok {
		database, _ = params["database"].(string)
	} else if _, ok := params["data"].(Dict)["database"].(string); ok {
		database, _ = params["data"].(Dict)["database"].(string)
	}
	api_id := any(nil)
	if _, ok := params["data"].(Dict)["api_id"]; ok {
		api_id = params["data"].(Dict)["api_id"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		api_id = params["data"].(Dict)["data"].(Dict)["api_id"]
	} else if _, ok := params["data"].(Dict)["api"].(Dict); ok {
		api_id = params["data"].(Dict)["api"].(Dict)["api_id"]
	}
	api_name := any(nil)
	if _, ok := params["data"].(Dict)["api_name"]; ok {
		api_name = params["data"].(Dict)["api_name"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		api_name = params["data"].(Dict)["data"].(Dict)["api_name"]
	} else if _, ok := params["data"].(Dict)["api"].(Dict); ok {
		api_name = params["data"].(Dict)["api"].(Dict)["api_name"]
	}
	_aux_params := params
	_aux_params["data"].(Dict)["table"] = []any{table, "api_header"}
	_aux_params["data"].(Dict)["db"] = database
	_aux_params["data"].(Dict)["limit"] = any(-1.0)
	_aux_params["data"].(Dict)["offset"] = any(0.0)
	// If api_id is not provided, try to find it using api_name
	if api_id == nil || api_id == any(nil) {
		if api_name != nil && api_name != any(nil) {
			_aux_params["data"].(Dict)["filters"] = []any{Dict{
				"field": "api_name",
				"cond":  "=",
				"value": api_name,
			}}
		}
	} else {
		_aux_params["data"].(Dict)["filters"] = []any{Dict{
			"field": "api_id",
			"cond":  "=",
			"value": api_id,
		}}
	}
	if (api_id == nil || api_id == any(nil)) && (api_name == nil || api_name == any(nil)) {
		return Dict{
			"success": false,
			"msg":     "No API ID or API name found!",
		}
	}
	return app.read(_aux_params)
}

func (app *application) runAPI(params Dict) Dict {
	api_data := app.getAPI(params)
	if success, ok := api_data["success"].(bool); !ok || !success {
		return Dict{
			"success": false,
			"msg":     "API not found!",
		}
	}
	api_records, ok := api_data["data"].(Dict)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "Invalid API data format!",
		}
	}
	api, ok := api_records["api"].(Dict)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "Invalid API data format!",
		}
	}
	if data, ok : 

}