package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
  api_id:                { type: integer, pk: true, autoincrement: true, comment: "API ID" }
  api_name:              { type: varchar, len: 100, nullable: false, comment: "API Name", form_display: true, table_display: true, form_size: 6, order: 1 }
  api_type_id:           { type: integer, fk: "api_type.api_type_id", nullable: false, comment: "API Type", form_display: true, table_display: true, form_size: 3, order: 2 }
  http_request_type_id:  { type: integer, fk: "http_request_type.http_request_type_id", nullable: false, comment: "HTTP Request Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  api_description:       { type: text, comment: "Description", form_display: true, form_long_text: true, table_display: true, order: 4 }
  endpoint:              { type: varchar, len: 500, nullable: false, comment: "API Endpoint", form_display: true, table_display: true, form_size: 9, order: 5 }
  request_body_template: { type: text, comment: "Request Body Template (Go template)", form_display: true, form_long_text: true, form_code: json, order: 6 }
  num_retries:           { type: integer, default: 3, comment: "Number of Retries", form_display: true, table_display: true, form_size: 3, order: 7 }
  timeout_seconds:       { type: integer, default: 30, comment: "Timeout (seconds)", form_display: true, table_display: true, form_size: 3, order: 8 }
  headers_template:      { type: text, comment: "Headers Template (use @VAR_NAME for environment variables)", form_display: true, form_long_text: true, form_code: json, order: 9 }
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
  request_at:       { type: datetime, nullable: false, comment: "Request DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", form_use_label: true, order: 3 }
  response_at:      { type: datetime, comment: "Response DateTime", form_display: true, table_display: true, form_date_format: "YY/MM/DD HH:mm:ss", form_use_label: true, order: 4 }
  request_body:     { type: text, comment: "Request Body", form_display: true, form_long_text: true, form_code: json, order: 5 }
  response_body:    { type: text, comment: "Response Body", form_display: true, form_long_text: true, form_code: json, order: 6 }
  response_status:  { type: integer, comment: "Response Status Code", form_display: true, table_display: true, order: 7 }
  response_message: { type: varchar, len: 500, comment: "Response Message", form_display: true, table_display: true, order: 8 }
  crud_trggrd_db:       { type: varchar, len: 50, comment: "Crud Triggered DB", form_display: true, table_display: true, order: 9, form_size: 3 }
  crud_trggrd_table:    { type: varchar, len: 50, comment: "Crud Triggered Table", form_display: true, table_display: true, order: 9, form_size: 3 }
  crud_trggrd_pk_field: { type: varchar, len: 50, comment: "Crud Triggered FK Field", form_display: true, table_display: true, order: 9, form_size: 3 }
  crud_trggrd_row_id:   { type: varchar, len: 50, comment: "Crud Triggered Row ID", form_display: true, table_display: true, order: 9, form_size: 3 }
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

func (app *application) runAPI(params Dict) Dict {
	//fmt.Println("DATA:", params["data"])
	var user_id int
	if _, ok := params["user"].(Dict)["user_id"]; ok {
		user_id = app.toInt(params["user"].(Dict)["user_id"])
	}
	var app_id int
	if _, ok := params["app"].(Dict)["app_id"]; ok {
		app_id = app.toInt(params["app"].(Dict)["app_id"])
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
	endpoint := any(nil)
	if _, ok := params["data"].(Dict)["endpoint"]; ok {
		endpoint = params["data"].(Dict)["endpoint"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		endpoint = params["data"].(Dict)["data"].(Dict)["endpoint"]
	} else if _, ok := params["data"].(Dict)["api"].(Dict); ok {
		endpoint = params["data"].(Dict)["api"].(Dict)["endpoint"]
	}
	// fmt.Println(1, api_id, api_name, endpoint)
	var api Dict
	var err error
	if api_id == nil || api_id == any(nil) {
		if api_name != nil && api_name != any(nil) {
			_sql := "select * from api where api_name = ? and excluded = false and active = true"
			api, err = app.AdminGetRowByFilter(_sql, []any{api_name})
			if err != nil {
				return Dict{
					"success": false,
					"msg":     "API not found with the provided name!",
				}
			}
		} else if endpoint != nil && endpoint != any(nil) {
			_sql := "select * from api where endpoint = ? and excluded = false and active = true"
			api, err = app.AdminGetRowByFilter(_sql, []any{endpoint})
			if err != nil {
				return Dict{
					"success": false,
					"msg":     "API not found with the provided endpoint!",
				}
			}
		} else {
			return Dict{
				"success": false,
				"msg":     "No API ID or API name found!",
			}
		}
	} else {
		sql := "select * from api where api_id = ? and excluded = false and active = true"
		api, err = app.AdminGetRowByID(sql, app.toInt(api_id))
		if err != nil {
			return Dict{
				"success": false,
				"msg":     "No API ID or API name found!",
			}
		}
	}
	_sql := "select * from api_header where api_id = ? and excluded = false and active = true"
	api_headers, err := app.AdminGetRowsByFilter(_sql, []any{api_id})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     "Failed to fetch API headers!",
		}
	}
	_data, ok := params["data"].(Dict)
	if !ok {
		_data = Dict{}
	}
	data, ok := _data["data"].(Dict)
	if !ok {
		data = Dict{}
	}
	// prepair api call by each api_type in api_details
	api_type_id := app.toInt(api["api_type_id"])
	http_request_type_id := app.toInt(api["http_request_type_id"])
	fmt.Println("HTTP Request Type ID:", http_request_type_id)
	sql := "select * from http_request_type where http_request_type_id = ? and excluded = false"
	http_request_type, err := app.AdminGetRowByID(sql, http_request_type_id)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     "HTTP Request Type not found!",
		}
	}
	endpoint, ok = api["endpoint"].(string)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "API endpoint is required for API call!",
		}
	}
	request_body_template, _ := api["request_body_template"].(string)
	//num_retries := app.toInt(api["num_retries"])
	//timeout_seconds := app.toInt(api["timeout_seconds"])
	var request_body string
	if request_body_template != "" {
		// Here you can implement logic to render the request_body_template with the appropriate data
		request_body, err = app.RenderTemplate(request_body_template, data)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     "Failed to render request body template!",
			}
		}
	}
	// fmt.Println("Rendered Request Body:", request_body)
	// Make API call based on api_type_id
	api_logs := Dict{
		"api_id":               api["api_id"],
		"api_name":             api["api_name"],
		"request_at":           time.Now(),
		"request_body":         request_body,
		"crud_trggrd_db":       data["crud_trggrd_db"],
		"crud_trggrd_table":    data["crud_trggrd_table"],
		"crud_trggrd_pk_field": data["crud_trggrd_pk_field"],
		"crud_trggrd_row_id":   data["crud_trggrd_row_id"],
		"user_id":              user_id,
		"app_id":               app_id,
		"excluded":             false,
	}
	switch int(api_type_id) {
	case 1: // REST
		if endpoint == "" {
			return Dict{
				"success": false,
				"msg":     "API endpoint is required for REST API!",
			}
		}
		// Implement REST API call logic here using endpoint, request_body, headers, num_retries, and timeout_seconds
		method := "GET"
		if _method, ok := http_request_type["http_request_type"].(string); ok && method != "" {
			method = _method
		}
		req, err := http.NewRequest(method, endpoint.(string), bytes.NewBuffer([]byte(request_body))) // bytes.NewBuffer(jsonBody)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     "Failed to create HTTP request!",
			}
		}
		// Set headers
		for _, value := range api_headers {
			header_name, _ := value["header_name"].(string)
			header_value, _ := value["header_value"].(string)
			if header_name != "" && header_value != "" {
				continue
			}
			req.Header.Set(header_name, header_value)
		}
		//req.Header.Set("Content-Type", "application/json")
		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     "Failed to make HTTP request!",
			}
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     "Failed to read HTTP response body!",
			}
		}
		api_logs["response_at"] = time.Now()
		api_logs["response_status"] = resp.StatusCode
		api_logs["response_message"] = resp.Status
		api_logs["response_body"] = string(bodyBytes)
		api_logs["created_at"] = time.Now()
		api_logs["updated_at"] = time.Now()
	case 2: // SOAP
		// Implement SOAP API call logic here using endpoint, request_body, headers, num_retries, and timeout_seconds
		return Dict{
			"success": false,
			"msg":     "SOAP API type is not implemented yet!",
		}
	case 3: // gRPC
		// Implement gRPC API call logic here using endpoint, request_body, headers, num_retries, and timeout_seconds
		return Dict{
			"success": false,
			"msg":     "gRPC API type is not implemented yet!",
		}
	case 4: // GraphQL
		// Implement GraphQL API call logic here using endpoint, request_body, headers, num_retries, and timeout_seconds
		return Dict{
			"success": false,
			"msg":     "GraphQL API type is not implemented yet!",
		}
	default:
		return Dict{
			"success": false,
			"msg":     "Unsupported API type!",
		}
	}
	// Save API call logs to database
	insert_query := `insert into api_call_log (api_id, api_name, request_at, response_at, request_body, response_body, response_status, response_message, user_id, app_id, created_at, updated_at, excluded)
	values (:api_id, :api_name, :request_at, :response_at, :request_body, :response_body, :response_status, :response_message, :user_id, :app_id, :created_at, :updated_at, :excluded)`
	_, err2 := app.db.ExecuteNamedQuery(insert_query, api_logs)
	if err2 != nil {
		fmt.Printf("Error inserting API call log for api_id %v: %v", api_logs["api_id"], err2)
	}
	return Dict{
		"success": true,
		"msg":     "API call executed successfully!",
		"data":    api_logs["response_body"],
	}
}
