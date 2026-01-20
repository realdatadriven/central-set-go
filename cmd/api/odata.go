package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseODataQuery(r *http.Request) (map[string]any, error) {
	q := r.URL.Query()
	out := make(map[string]any)
	for key, values := range q {
		if !strings.HasPrefix(key, "$") {
			continue
		}
		if len(values) == 0 {
			continue
		}
		value := values[0]
		switch key {
		case "$filter":
			out["$filter"] = value
		case "$top", "$skip":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid %s: %s", key, value)
			}
			out[key] = v
		case "$select":
			out["$select"] = splitAndTrim(value)
		case "$orderby":
			out["$orderby"] = value
		case "$expand":
			out["$expand"] = splitAndTrim(value)
			return nil, fmt.Errorf("%s: %s not suported yet!", key, value)
		case "$count":
			b, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("invalid $count: %s", value)
			}
			out["$count"] = b
			// return nil, fmt.Errorf("%s: %s not suported yet!", key, value)
		default:
			// keep unknown $params for future support
			out[key] = value
		}
	}
	return out, nil
}
func ODataToCentralParams(odata map[string]any) (map[string]any, error) {
	/*out := map[string]any{
		"filters":  []any{},
		"fields":   []any{},
		"offset":   any(0),
		"limit":    any(10),
		"order_by": []any{},
	}*/
	out := map[string]any{}
	// $select -> fields
	if v, ok := odata["$select"]; ok {
		out["fields"] = toAnySlice(v.([]string))
	}
	// $top -> limit
	if v, ok := odata["$top"]; ok {
		out["limit"] = v
	}
	// $skip -> offset
	if v, ok := odata["$skip"]; ok {
		out["offset"] = v
	}
	// $orderby -> order_by
	if v, ok := odata["$orderby"]; ok {
		orderBy, err := convertOrderBy(v.(string))
		if err != nil {
			return nil, err
		}
		out["order_by"] = orderBy
	}
	// $filter -> filters
	if v, ok := odata["$filter"]; ok {
		filters, err := convertFilter(v.(string))
		if err != nil {
			return nil, err
		}
		out["filters"] = filters
	}
	return out, nil
}

var odataOpToCond = map[string]string{
	"eq":      "=",
	"ne":      "!=",
	"gt":      ">",
	"ge":      ">=",
	"lt":      "<",
	"le":      "<=",
	"like":    "LIKE",
	"in":      "IN",
	"between": "BETWEEN",
	// string functions
	"startswith": "LIKE",
	"endswith":   "LIKE",
	"contains":   "LIKE",
}

func convertFilter(expr string) ([]any, error) {
	clauses := splitOnAnd(expr)
	out := make([]any, 0, len(clauses))
	for _, c := range clauses {
		filter, err := parseSimpleFilter(c)
		if err != nil {
			return nil, err
		}
		out = append(out, filter)
	}
	return out, nil
}
func stripSingleQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}
func parseFunctionFilter(expr string) (map[string]any, error) {
	open := strings.Index(expr, "(")
	close := strings.LastIndex(expr, ")")
	if open < 0 || close <= open {
		return nil, fmt.Errorf("invalid function syntax: %s", expr)
	}
	fn := strings.ToLower(strings.TrimSpace(expr[:open]))
	args := expr[open+1 : close]
	parts := strings.SplitN(args, ",", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid function args: %s", expr)
	}
	field := strings.TrimSpace(parts[0])
	value := stripSingleQuotes(strings.TrimSpace(parts[1]))
	cond, ok := odataOpToCond[fn]
	if !ok {
		return nil, fmt.Errorf("unsupported function: %s", fn)
	}
	switch fn {
	case "startswith":
		value = value + "%"
	case "endswith":
		value = "%" + value
	case "contains":
		value = "%" + value + "%"
	}
	return map[string]any{
		"field": field,
		"cond":  cond,
		"value": value,
	}, nil
}
func parseSimpleFilter(expr string) (map[string]any, error) {
	expr = strings.TrimSpace(expr)
	// function-style
	if strings.Contains(expr, "(") && strings.HasSuffix(expr, ")") {
		open := strings.Index(expr, "(")
		fn := strings.TrimSpace(expr[:open])
		fmt.Println("FUNCTION: ", fn)
		if strings.Contains(fn, " ") { // HAS ( BUT NOT A FUNCTION
		} else {
			return parseFunctionFilter(expr)
		}
	}
	tokens := strings.Fields(expr)
	if len(tokens) < 3 {
		return nil, fmt.Errorf("invalid filter: %s", expr)
	}
	field := tokens[0]
	op := strings.ToLower(tokens[1])
	cond, ok := odataOpToCond[op]
	if !ok {
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}
	switch op {
	case "in":
		value := strings.Join(tokens[2:], " ")
		value = strings.Trim(value, "()")
		value = strings.ReplaceAll(value, ",", ";")
		// strip quotes from each item
		parts := strings.Split(value, ";")
		for i := range parts {
			parts[i] = stripSingleQuotes(strings.TrimSpace(parts[i]))
		}
		return map[string]any{
			"field": field,
			"cond":  cond,
			"value": strings.Join(parts, ";"),
		}, nil
	case "between":
		if len(tokens) < 5 || strings.ToLower(tokens[3]) != "and" {
			return nil, fmt.Errorf("invalid between syntax: %s", expr)
		}
		v1 := stripSingleQuotes(tokens[2])
		v2 := stripSingleQuotes(tokens[4])
		return map[string]any{
			"field": field,
			"cond":  cond,
			"value": v1 + ";" + v2,
		}, nil
	default:
		value := stripSingleQuotes(strings.Join(tokens[2:], " "))
		return map[string]any{
			"field": field,
			"cond":  cond,
			"value": value,
		}, nil
	}
}
func splitOnAnd(expr string) []string {
	parts := strings.Split(expr, " and ")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
func toAnySlice[T any](in []T) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
func convertOrderBy(expr string) ([]any, error) {
	parts := strings.Split(expr, ",")
	out := make([]any, 0, len(parts))
	for _, p := range parts {
		fields := strings.Fields(strings.TrimSpace(p))
		if len(fields) == 0 {
			continue
		}
		item := map[string]any{
			"field": fields[0],
			"order": "ASC",
		}
		if len(fields) > 1 {
			order := strings.ToUpper(fields[1])
			if order != "ASC" && order != "DESC" {
				return nil, fmt.Errorf("invalid order: %s", fields[1])
			}
			item["order"] = order
		}
		out = append(out, item)
	}
	return out, nil
}
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
func (app *application) odata_api(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	w.Header().Set("Content-Type", "application/json")
	sql := `select * from app where (app = ? or db = ?) and excluded = false`
	_app, err := app.AdminGetRowByFilter(sql, []any{db, table})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_res := map[string]any{
			"error": map[string]any{
				"code":    "GeneralError",
				"message": "AppErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(_res)
		return
	}
	//fmt.Println(_app)
	odata_params, err := ParseODataQuery(r)
	response := map[string]any{
		"error": map[string]any{
			"code":    "BadRequest",
			"message": "Invalid filter expression",
			"target":  "filter",
			"details": []map[string]any{
				{
					"code":    "SyntaxError",
					"message": "Expected 'eq' or 'ne' operator",
					"target":  "$filter",
				},
			},
			"innererror": map[string]any{
				"trace":   "...",
				"context": "...",
			},
		},
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]any{
			"error": map[string]any{
				"code":    "GeneralError",
				"message": "ParseODataQueryErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	csParams, err := ODataToCentralParams(odata_params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]any{
			"error": map[string]any{
				"code":    "GeneralError",
				"message": "ODataToCentralSetParamsErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := app.verifyToken(r)
	//fmt.Println(params["user"].(Dict)["username"].(string), "->", int(params["user"].(Dict)["user_id"].(float64)), "->", int(params["user"].(Dict)["role_id"].(float64)))
	var data Dict
	params := map[string]any{
		"lang": "en",
		"app":  _app,
		"data": map[string]any{
			"db":    db,
			"table": table,
			//"odata_params": odata_params,
			//"cs_params":    csParams,
		},
	}
	for key, val := range csParams {
		//fmt.Println(key, val)
		params["data"].(map[string]any)[key] = val
	}
	_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
	_log := Dict{
		"action": fmt.Sprintf("OData %s/%s", db, table),
		"req_ip": _ip,
		"res_at": time.Now(),
	}
	if token["success"].(bool) {
		params["user"] = *(contextGetAuthenticatedUser(r))
		_log["user_id"] = params["user"].(Dict)["user_id"]
	}
	if !token["success"].(bool) {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]any{
			"error": map[string]any{
				"code":    "GeneralError",
				"message": token["msg"],
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		data = app.read(params)
	}
	// LOGS
	actions_not_to_log := app.sliceStrs2SliceInterfaces(strings.Split(app.config.actions_not_to_log, ","))
	if !app.contains(actions_not_to_log, "odata") {
		_log["res_type"] = "success"
		if _, ok := data["success"]; !ok {
			_log["res_type"] = "error"
		} else if _, ok := data["success"].(bool); !ok {
			_log["res_type"] = "error"
		} else if success, ok := data["success"].(bool); ok {
			if success {
				_log["res_type"] = "success"
			}
		}
		_log["res_msg"] = data["msg"]
		_log["row_id"] = data["inserted_primary_key"]
		_log["table"] = params["data"].(Dict)["table"]
		_log["db"] = ""
		if _, ok := params["data"].(Dict)["database"]; ok {
			_log["db"] = params["data"].(Dict)["database"]
		} else if _, ok := params["data"].(Dict)["db"]; ok {
			_log["db"] = params["data"].(Dict)["db"]
		} else if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(Dict)["db"]; ok {
			_log["db"] = params["app"].(Dict)["db"]
		}
		if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(Dict)["app_id"]; ok {
			_log["app_id"] = params["app"].(Dict)["app_id"]
		}
		_log["excluded"] = false
		//fmt.Println(_log)
		_log_params := Dict{
			"data": Dict{
				"data":  _log,
				"table": "user_log",
				"db":    app.config.db.dsn,
			},
			"app": Dict{
				"app_id": any(1.0),
				"db":     filepath.Base(app.config.db.dsn),
			},
			"user": Dict{
				"user_id": any(1.0),
				"role_id": any(1.0),
			},
		}
		res := app.create_update(_log_params)
		if _, ok := res["success"]; !ok {
			fmt.Println("Err processing logs:", res)
		} else if _, ok := res["success"].(bool); !ok {
			fmt.Println("Err processing logs:", res["msg"])
		} else if !res["success"].(bool) {
			fmt.Println("Err processing logs:", res["msg"])
		}
	}
	if !data["success"].(bool) {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]any{
			"error": map[string]any{
				"code":    "GeneralError",
				"message": data["msg"],
				//"params":  params,
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response = map[string]any{
		"@odata.context": "http://" + r.Host + strings.TrimSuffix(r.URL.Path, r.URL.RawQuery),
		"@odata.count":   data["total"], // optional: total count (when $count=true or $inlinecount)
		"value":          data["data"],
		//"params":         params,
		//"@odata.nextLink": "http://" + r.Host + strings.TrimSuffix(r.URL.Path, r.URL.RawQuery), // optional: for paging
	}
	json.NewEncoder(w).Encode(response)
}
