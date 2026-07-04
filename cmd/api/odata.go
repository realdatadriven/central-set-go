package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseODataQuery(q url.Values /*r *http.Request*/) (Dict, error) {
	//q := r.URL.Query()
	out := make(Dict)
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
func ODataToCentralParams(odata Dict) (Dict, error) {
	/*out := Dict{
		"filters":  []any{},
		"fields":   []any{},
		"offset":   any(0),
		"limit":    any(10),
		"order_by": []any{},
	}*/
	out := Dict{}
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
	// $schema = true
	if v, ok := odata["$schema"]; ok {
		if b, err := strconv.ParseBool(v.(string)); err == nil {
			out["include_schema"] = b
		}
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
func parseFunctionFilter(expr string) (Dict, error) {
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
	return Dict{
		"field": field,
		"cond":  cond,
		"value": value,
	}, nil
}
func parseSimpleFilter(expr string) (Dict, error) {
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
		return Dict{
			"field": field,
			"cond":  cond,
			"value": strings.Join(parts, ","),
		}, nil
	case "between":
		if len(tokens) < 5 || strings.ToLower(tokens[3]) != "and" {
			return nil, fmt.Errorf("invalid between syntax: %s", expr)
		}
		v1 := stripSingleQuotes(tokens[2])
		v2 := stripSingleQuotes(tokens[4])
		return Dict{
			"field": field,
			"cond":  cond,
			"value": v1 + "," + v2,
		}, nil
	default:
		value := stripSingleQuotes(strings.Join(tokens[2:], " "))
		return Dict{
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
		item := Dict{
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
func sqlToEdmType(sqlType string) string {
	t := strings.ToUpper(sqlType)
	switch {
	case strings.Contains(t, "INT"):
		return "Edm.Int32"
	case strings.Contains(t, "BIGINT"):
		return "Edm.Int64"
	case strings.Contains(t, "BOOLEAN"):
		return "Edm.Boolean"
	case strings.Contains(t, "DATETIME"), strings.Contains(t, "TIMESTAMP"):
		return "Edm.DateTimeOffset"
	case strings.Contains(t, "DATE"):
		return "Edm.Date"
	case strings.Contains(t, "DECIMAL"), strings.Contains(t, "NUMERIC"):
		return "Edm.Decimal"
	case strings.Contains(t, "FLOAT"), strings.Contains(t, "DOUBLE"):
		return "Edm.Double"
	default:
		return "Edm.String"
	}
}
func BuildODataMetadata(rows []Dict) (string, error) {
	type Column struct {
		Name     string
		Type     string
		Nullable bool
		IsPK     bool
	}
	type Entity struct {
		Name    string
		Columns []Column
		PKs     []string
	}
	schemas := map[string]map[string]*Entity{}
	for _, r := range rows {
		if r["excluded"] == true {
			continue
		}
		db := r["db"].(string)
		table := r["table"].(string)
		field := r["field"].(string)
		sqlType := r["type"].(string)
		nullable := r["nullable"] != false
		isPK := r["pk"] == true || r["pk"] == 1
		if schemas[db] == nil {
			schemas[db] = map[string]*Entity{}
		}
		if schemas[db][table] == nil {
			schemas[db][table] = &Entity{Name: table}
		}
		col := Column{
			Name:     field,
			Type:     sqlToEdmType(sqlType),
			Nullable: nullable,
			IsPK:     isPK,
		}
		e := schemas[db][table]
		e.Columns = append(e.Columns, col)
		if isPK {
			e.PKs = append(e.PKs, field)
		}
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	b.WriteString(`
<edmx:Edmx Version="4.0"
 xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx"
 xmlns="http://docs.oasis-open.org/odata/ns/edm">
 <edmx:DataServices>`)
	for db, tables := range schemas {
		b.WriteString(fmt.Sprintf(`<Schema Namespace="%s">`, db))
		for _, e := range tables {
			b.WriteString(fmt.Sprintf(`<EntityType Name="%s">`, e.Name))
			if len(e.PKs) > 0 {
				b.WriteString(`<Key>`)
				for _, pk := range e.PKs {
					b.WriteString(fmt.Sprintf(`<PropertyRef Name="%s"/>`, pk))
				}
				b.WriteString(`</Key>`)
			}
			for _, c := range e.Columns {
				b.WriteString(fmt.Sprintf(
					`<Property Name="%s" Type="%s" Nullable="%t"/>`,
					c.Name, c.Type, c.Nullable,
				))
			}
			b.WriteString(`</EntityType>`)
		}
		b.WriteString(`<EntityContainer Name="Container">`)
		for _, e := range tables {
			b.WriteString(fmt.Sprintf(
				`<EntitySet Name="%s" EntityType="%s.%s"/>`,
				e.Name, db, e.Name,
			))
		}
		b.WriteString(`</EntityContainer>`)
		b.WriteString(`</Schema>`)
	}
	b.WriteString(`
 </edmx:DataServices>
</edmx:Edmx>`)
	return b.String(), nil
}
func BuildODataMetadataJSON(rows []Dict) (Dict, error) {
	type Entity struct {
		Props Dict
		Keys  []string
	}
	model := Dict{
		"@odata.context": "$metadata",
	}
	namespaces := map[string]map[string]*Entity{}
	for _, r := range rows {
		if r["excluded"] == true {
			continue
		}
		db := r["db"].(string)
		table := r["table"].(string)
		field := r["field"].(string)
		sqlType := r["type"].(string)
		nullable := r["nullable"] != false
		isPK := r["pk"] == true || r["pk"] == 1
		if namespaces[db] == nil {
			namespaces[db] = map[string]*Entity{}
		}
		if namespaces[db][table] == nil {
			namespaces[db][table] = &Entity{
				Props: Dict{
					"$Kind": "EntityType",
				},
			}
		}
		prop := Dict{
			"$Type": sqlToEdmType(sqlType),
		}
		if !nullable {
			prop["$Nullable"] = false
		}
		namespaces[db][table].Props[field] = prop
		if isPK {
			namespaces[db][table].Keys = append(
				namespaces[db][table].Keys,
				field,
			)
		}
	}
	for db, tables := range namespaces {
		ns := Dict{}
		// entity types
		for name, e := range tables {
			if len(e.Keys) > 0 {
				e.Props["$Key"] = e.Keys
			}
			ns[name] = e.Props
		}
		// entity container
		container := Dict{
			"$Kind": "EntityContainer",
		}
		for name := range tables {
			container[name] = Dict{
				"$Collection": true,
				"$Type":       db + "." + name,
			}
		}
		ns["Container"] = container
		model[db] = ns
	}
	return model, nil
}
func (app *application) odata_api_metadata(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("METADATA ONLY!")
	db := r.PathValue("db")
	w.Header().Set("OData-Version", "4.0")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	sql := `select * from table_schema where db = ? and excluded = false  order by field_order`
	_table_schema, err := app.AdminGetRowsByFilter(sql, []any{db})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_res := Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": "MetadataErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(_res)
		return
	}
	xml, err := BuildODataMetadata(_table_schema)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_res := Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": "BuildMetadataErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(_res)
		return
	}
	//fmt.Println(xml)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml))
	/*model, err := BuildODataMetadataJSON(_table_schema)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_res := Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": "BuildMetadataErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(_res)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)*/
}
func isXMLMetadataRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/xml") ||
		strings.Contains(accept, "application/atom+xml")
}
func (app *application) OData2C7Read(params Dict, db, table string, odata_path url.Values) ([]Dict, error) {
	sql := `select * from app where (app = ? or db = ?) and excluded = false`
	_app, err := app.AdminGetRowByFilter(sql, []any{db, table})
	if err != nil {
		return nil, err
	}
	odata_params, err := ParseODataQuery(odata_path)
	if err != nil {
		return nil, err
	}
	csParams, err := ODataToCentralParams(odata_params)
	if err != nil {
		return nil, err
	}
	var data Dict
	params["app"] = _app
	params["data"].(Dict)["db"] = db
	params["data"].(Dict)["table"] = table
	for key, val := range csParams {
		params["data"].(Dict)[key] = val
	}
	data = app.read(params)
	if !data["success"].(bool) {
		return nil, fmt.Errorf("%s", data["msg"])
	}
	if _, ok := data["data"].([]Dict); !ok {
		return nil, fmt.Errorf("%s", data["msg"])
	}
	return data["data"].([]Dict), nil
}
func (app *application) odata_api(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	w.Header().Set("OData-Version", "4.0")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("Content-Type", "application/json")
	if table == "$metadata" {
		sql := `select * from table_schema where db = ? and excluded = false  order by field_order`
		_table_schema, err := app.AdminGetRowsByFilter(sql, []any{db, table})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_res := Dict{
				"error": Dict{
					"code":    "GeneralError",
					"message": "MetadataErr: " + err.Error(),
				},
			}
			json.NewEncoder(w).Encode(_res)
			return
		}
		xml, err := BuildODataMetadata(_table_schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_res := Dict{
				"error": Dict{
					"code":    "GeneralError",
					"message": "BuildDBMetadataErr: " + err.Error(),
				},
			}
			json.NewEncoder(w).Encode(_res)
			return
		}
		//fmt.Println(xml)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(xml))
		/*model, err := BuildODataMetadataJSON(_table_schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_res := Dict{
				"error": Dict{
					"code":    "GeneralError",
					"message": "BuildMetadataErr: " + err.Error(),
				},
			}
			json.NewEncoder(w).Encode(_res)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model)*/
		return
	} else if isXMLMetadataRequest(r) {
		sql := `select * from "table_schema" where "db" = ? and "table" = ? and "excluded" = false`
		_table_schema, err := app.AdminGetRowsByFilter(sql, []any{db, table})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_res := Dict{
				"error": Dict{
					"code":    "GeneralError",
					"message": "MetadataErr: " + err.Error(),
				},
			}
			json.NewEncoder(w).Encode(_res)
			return
		}
		xml, err := BuildODataMetadata(_table_schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_res := Dict{
				"error": Dict{
					"code":    "GeneralError",
					"message": "BuildTableMetadataErr: " + err.Error(),
				},
			}
			json.NewEncoder(w).Encode(_res)
			return
		}
		/*/fmt.Println(xml)
		url := "http://" + r.Host + strings.TrimSuffix(r.URL.Path, r.URL.RawQuery)
		xml = fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom"
		      xmlns:d="http://docs.oasis-open.org/odata/ns/data"
		      xmlns:m="http://docs.oasis-open.org/odata/ns/metadata">
		  <title type="text">%s</title>
		  <id>%s</id>
		</feed>`, table, url)*/
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(xml))
		return
	}
	sql := `select * from app where (app = ? or db = ?) and excluded = false`
	_app, err := app.AdminGetRowByFilter(sql, []any{db, table})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_res := Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": "AppErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(_res)
		return
	}
	//fmt.Println(_app)
	q := r.URL.Query()
	odata_params, err := ParseODataQuery(q)
	response := Dict{
		"error": Dict{
			"code":    "BadRequest",
			"message": "Invalid filter expression",
			"target":  "filter",
			"details": []Dict{
				{
					"code":    "SyntaxError",
					"message": "Expected 'eq' or 'ne' operator",
					"target":  "$filter",
				},
			},
			"innererror": Dict{
				"trace":   "...",
				"context": "...",
			},
		},
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = Dict{
			"error": Dict{
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
		response = Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": "ODataToCentralSetParamsErr: " + err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := app.verifyToken(r)
	//fmt.Println(params["user"].(Dict)["username"].(string), "->", app.toInt(params["user"].(Dict)["user_id"]), "->", app.toInt(params["user"].(Dict)["role_id"]))
	var data Dict
	params := Dict{
		"lang": "en",
		"app":  _app,
		"data": Dict{
			"db":    db,
			"table": table,
			//"odata_params": odata_params,
			//"cs_params":    csParams,
		},
	}
	for key, val := range csParams {
		//fmt.Println(key, val)
		params["data"].(Dict)[key] = val
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
	//fmt.Println(params["user"])
	if !token["success"].(bool) {
		w.WriteHeader(http.StatusForbidden)
		response = Dict{
			"error": Dict{
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
		response = Dict{
			"error": Dict{
				"code":    "GeneralError",
				"message": data["msg"],
				//"params":  params,
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response = Dict{
		"@odata.context": fmt.Sprintf("%s/odata/%s/$metadata#%s", r.Host, db, table), //"http://" + r.Host + strings.TrimSuffix(r.URL.Path, r.URL.RawQuery),
		"@odata.count":   data["total"],                                              // optional: total count (when $count=true or $inlinecount)
		"value":          data["data"],
		//"params":         params,
		//"@odata.nextLink": "http://" + r.Host + strings.TrimSuffix(r.URL.Path, r.URL.RawQuery), // optional: for paging
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

/**
INSTALL erpl_web FROM community;
LOAD erpl_web;

-- Usefull for debugging
SET erpl_trace_enabled = TRUE;
SET erpl_trace_level = 'DEBUG';

-- Create Secret
CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '...',
  SCOPE 'http://localhost:4444/'
);

FROM HTTP_GET('http://localhost:4444/odata/ADMIN/app');

FROM ODATA_READ('http://localhost:4444/odata/ADMIN/app?$filter=app_id gt 1');

ATTACH IF NOT EXISTS 'http://localhost:4444/odata/ADMIN' AS admin (TYPE odata);
SELECT app_id, app, db, excluded FROM admin.app WHERE app_id = 1;
**/
