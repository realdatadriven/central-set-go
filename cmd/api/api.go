package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"
)

func parsePath(s string) (db string, table string, query url.Values, err error) {
	parts := strings.SplitN(s, "?", 2)
	path := parts[0]
	if len(parts) == 2 {
		query, err = url.ParseQuery(parts[1])
		if err != nil {
			return "", "", nil, err
		}
	} else {
		query = url.Values{}
	}
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) != 2 {
		return "", "", nil, fmt.Errorf("invalid path format: expected db/table")
	}
	db = pathParts[0]
	table = pathParts[1]
	return
}

func (app *application) GetAPIData(params Dict, api_data_res []Dict, _data Dict) (Dict, error) {
	// loop api_data_res
	etlx_engine := &etlx.ETLX{}
	res := Dict{}
	for _, api_data := range api_data_res {
		//fmt.Println(api_data)
		name := api_data["api_data"].(string)
		if app.toInt(api_data["api_data_type_id"]) == 3 { // ODATA
			odata_path, ok := api_data["odata_path"].(string)
			if !ok {
				return nil, fmt.Errorf("Error, odata_path is not set!")
			}
			odata_path, err := etlx_engine.RenderTemplate(odata_path, _data)
			if err != nil {
				odata_path = api_data["odata_path"].(string)
			} else {
				odata_path = etlx_engine.ReplaceEnvVariable(odata_path)
			}
			fmt.Println(api_data["odata_path"], odata_path)
			sigle_row_obj := app.toBool(api_data["sigle_row_obj"])
			db, table, query, err := parsePath(odata_path)
			fmt.Println(db, table, query)
			if err != nil {
				return nil, err
			}
			results, err := app.OData2C7Read(params, db, table, query)
			if err != nil {
				return nil, err
			}
			if sigle_row_obj {
				aux := Dict{}
				if len(results) > 0 {
					aux = results[0]
				}
				res[name] = aux
			} else {
				res[name] = results
			}
		} else {
			return nil, fmt.Errorf("Error: API Data Type %d is not implemented yet!", app.toInt(api_data["api_data_type_id"]))
		}
	}
	return res, nil
}

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
	if _, ok := params["data"].(Dict)["api_endpoint"]; ok {
		endpoint = params["data"].(Dict)["api_endpoint"]
	} else if _, ok := params["data"].(Dict)["endpoint"]; ok {
		endpoint = params["data"].(Dict)["endpoint"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		endpoint = params["data"].(Dict)["data"].(Dict)["endpoint"]
	} else if _, ok := params["data"].(Dict)["api"].(Dict); ok {
		endpoint = params["data"].(Dict)["api"].(Dict)["endpoint"]
	}
	//fmt.Println(1, api_id, api_name, endpoint)
	var api Dict
	var err error
	if api_id == nil && api_id == any(nil) {
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
		if len(api) == 0 {
			return Dict{
				"success": false,
				"msg":     "API not found with the provided details!",
			}
		}
		api_id = api["api_id"]
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
	fmt.Println(1, api_id, api_name, endpoint)
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
	// fmt.Println("HTTP Request Type ID:", http_request_type_id)
	sql := "select * from http_request_type where http_request_type_id = ? and excluded = false"
	http_request_type, err := app.AdminGetRowByID(sql, http_request_type_id)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     "HTTP Request Type not found!",
		}
	}
	// API DATA TO HELP BUILD THE TEPLATE
	sql = "select * from api_data where api_id = ? and excluded = false"
	api_data_res, err := app.AdminGetRowsByFilter(sql, []any{api_id})
	if err != nil {
		fmt.Println("Error getting API Data:", err)
		return Dict{
			"success": false,
			"msg":     "Error getting API Data!",
		}
	}
	// fmt.Println("API DATA:", api_data_res)
	api_data, err := app.GetAPIData(params, api_data_res, data)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Error getting API Data: %v", err.Error()),
		}
	}
	for key, val := range api_data {
		_data[key] = val
	}
	api_endpoint, ok := api["endpoint"].(string)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "API endpoint is required for API call!",
		}
	}
	api_endpoint, err = app.RenderTemplate(api_endpoint, _data)
	keys := []any{}
	for key := range _data {
		keys = append(keys, key)
	}
	fmt.Println("API ENDPOINT:", api["endpoint"], api_endpoint, keys) // _data["data"]
	if err != nil {
		api_endpoint = api["endpoint"].(string)
	}
	request_body_template, _ := api["request_body_template"].(string)
	//num_retries := app.toInt(api["num_retries"])
	//timeout_seconds := app.toInt(api["timeout_seconds"])
	var request_body string
	if request_body_template != "" {
		// Here you can implement logic to render the request_body_template with the appropriate data
		request_body, err = app.RenderTemplate(request_body_template, _data)
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
		if api_endpoint == "" {
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
		req, err := http.NewRequest(method, api_endpoint, bytes.NewBuffer([]byte(request_body))) // bytes.NewBuffer(jsonBody)
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
		//fmt.Println("HTTP Response Status:", resp.Status)
		//fmt.Println("HTTP Response Body:", string(bodyBytes))
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
	_data["logs"] = api_logs
	after_sql, okAfterSQL := api["after_sql"].(string)
	if okAfterSQL && after_sql != "" {
		after_sql, err = etlx_engine.RenderTemplate(after_sql, _data)
		if err != nil {
			return fmt.Errorf("Error rendering API after sql %v", err.Error())
		}
		after_sql = etlx_engine.ReplaceEnvVariable(after_sql)
		err = app.ExecuteQuery(after_sql, params, _data)
		if err != nil {
			success = false
			/*crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Errorf("Error executing after SQL %s!", err.Error())
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)*/
			if err != nil {
				fmt.Printf("Error inserting crud action log for unknown action_type_id for crud_action_id %v: %v", c_action["crud_action_id"], err)
			}
			return fmt.Errorf("Error executing after SQL %s!", err.Error())
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
