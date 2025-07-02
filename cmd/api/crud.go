package main

import (
	"fmt"

	"github.com/realdatadriven/etlx"
)

func (app *application) read(params map[string]any) map[string]any {
	// DATABASE
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	tables := []any{}
	if !app.IsEmpty(params["data"].(map[string]any)["table"]) {
		value := params["data"].(map[string]any)["table"]
		switch value.(type) {
		case nil:
			_ = true
		case string:
			tables = append(tables, params["data"].(map[string]any)["table"].(string))
		case []any:
			_tables := params["data"].(map[string]any)["table"].([]any)
			fmt.Println(_tables)
			for t := 0; t < len(_tables); t++ {
				if _, ok := _tables[t].(string); ok {
					tables = append(tables, _tables[t].(string))
				}
			}
		case map[any]any:
			// pass
		default:
			_ = true
		}
	} else if !app.IsEmpty(params["data"].(map[string]any)["tables"]) {
		value := params["data"].(map[string]any)["tables"]
		switch value.(type) {
		case string:
			tables = append(tables, params["data"].(map[string]any)["tables"].(string))
		case []any:
			_tables := params["data"].(map[string]any)["tables"].([]any)
			for t := 0; t < len(_tables); t++ {
				tables = append(tables, _tables[t].(string))
			}
		default:
			_ = true
		}
	}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	//fmt.Println("TABLES TO READ:", tables)
	_schemas := app.tables(params, tables)
	if !_schemas["success"].(bool) {
		return _schemas
	}
	if _, ok := _schemas["data"]; ok {
		_schemas = _schemas["data"].(map[string]any)
		params["schemas"] = _schemas
	}
	_permissions := app.table_access(params, tables)
	if !_permissions["success"].(bool) {
		return _permissions
	}
	if _, ok := _permissions["data"]; ok {
		_permissions = _permissions["data"].(map[string]any)
	} else {
		_permissions = map[string]any{}
	}
	_row_level_tables := app.row_level_tables(params)
	if !_row_level_tables["success"].(bool) {
		return _row_level_tables
	}
	params["row_level_tables"] = []string{}
	if _, ok := _row_level_tables["tables"]; ok {
		params["row_level_tables"] = _row_level_tables["tables"].([]string)
	}
	data := map[string]any{}
	for _, table := range tables {
		params["schema"] = map[string]any{}
		if _, ok := _schemas[table.(string)]; ok {
			params["schema"] = _schemas[table.(string)].(map[string]any)
		}
		params["permissions"] = map[string]any{}
		if _, ok := _permissions[table.(string)]; ok {
			params["permissions"] = _permissions[table.(string)].(map[string]any)
		}
		if len(tables) == 1 {
			return app.CrudRead(params, table.(string), newDB)
		}
		data[table.(string)] = app.CrudRead(params, table.(string), newDB)
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    data,
		//"tables":       tables,
		//"_schemas":     _schemas,
		//"permissions": _permissions,
		//"sql": query,
	}
}

func (app *application) create_update(params map[string]any) map[string]any {
	// DATABASE
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	if newDB != nil {
		defer newDB.Close()
	}
	table := ""
	if _, ok := params["data"].(map[string]any)["table"]; ok {
		table = params["data"].(map[string]any)["table"].(string)
	}
	if table == "" {
		msg, _ := app.i18n.T("no-table", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	tables := []any{table}
	var _data any
	if _, ok := params["data"].(map[string]any)["data"]; ok {
		_data = params["data"].(map[string]any)["data"]
	}
	switch _data.(type) {
	case []map[string]any:
		for _, d := range _data.([]map[string]any) {
			if _, ok := d["_table"]; ok {
				if app.contains(tables, d["_table"].(string)) {
					tables = append(tables, d["_table"].(string))
				}
			}
		}
	default:
		_ = ""
	}
	//fmt.Println("TABLES:", tables)
	_schemas := app.tables(params, tables)
	if !_schemas["success"].(bool) {
		return _schemas
	}
	if _, ok := _schemas["data"]; ok {
		_schemas = _schemas["data"].(map[string]any)
		params["schemas"] = _schemas
	}
	// fmt.Println("TABLES TO CREATE:", tables, _schemas)
	_permissions := app.table_access(params, tables)
	if !_permissions["success"].(bool) {
		return _permissions
	}
	if _, ok := _permissions["data"]; ok {
		_permissions = _permissions["data"].(map[string]any)
	} else {
		_permissions = map[string]any{}
	}
	_row_level_tables := app.row_level_tables(params)
	if !_row_level_tables["success"].(bool) {
		return _row_level_tables
	}
	params["row_level_tables"] = []string{}
	if _, ok := _row_level_tables["tables"]; ok {
		params["row_level_tables"] = _row_level_tables["tables"].([]string)
	}
	data := map[string]any{}
	switch _data.(type) {
	case []any:
		for i, d := range _data.([]any) {
			tbl := table
			if _, ok := d.(map[string]any)["_table"]; ok {
				tbl = d.(map[string]any)["_table"].(string)
			}
			params["schema"] = map[string]any{}
			if _, ok := _schemas[tbl]; ok {
				params["schema"] = _schemas[tbl].(map[string]any)
			}
			params["permissions"] = map[string]any{}
			if _, ok := _permissions[tbl]; ok {
				params["permissions"] = _permissions[tbl].(map[string]any)
			}
			//fmt.Println(i, tbl, d)
			params["data"].(map[string]any)["data"] = d
			data[fmt.Sprintf(`row-%s-%d`, tbl, i)] = app.CrudCreateUpdte(params, tbl, newDB)
		}
	case []map[string]any:
		for i, d := range _data.([]map[string]any) {
			tbl := table
			if _, ok := d["_table"]; ok {
				tbl = d["_table"].(string)
			}
			params["schema"] = map[string]any{}
			if _, ok := _schemas[tbl]; ok {
				params["schema"] = _schemas[tbl].(map[string]any)
			}
			params["permissions"] = map[string]any{}
			if _, ok := _permissions[tbl]; ok {
				params["permissions"] = _permissions[tbl].(map[string]any)
			}
			params["data"].(map[string]any)["data"] = d
			data[fmt.Sprintf(`row-%s-%d`, tbl, i)] = app.CrudCreateUpdte(params, tbl, newDB)
			fmt.Println(fmt.Sprintf(`row-%s-%d`, tbl, i), data[fmt.Sprintf(`row-%s-%d`, tbl, i)].(map[string]any)["msg"])
		}
	case map[string]any:
		params["schema"] = map[string]any{}
		if _, ok := _schemas[table]; ok {
			params["schema"] = _schemas[table].(map[string]any)
		}
		params["permissions"] = map[string]any{}
		if _, ok := _permissions[table]; ok {
			params["permissions"] = _permissions[table].(map[string]any)
		}
		return app.CrudCreateUpdte(params, table, newDB)
	default:
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}

func (app *application) query(params map[string]any) map[string]any {
	// DATABASE
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	if app.IsEmpty(_data) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	var _query any
	if _, ok := _data["query"]; ok {
		_query = _data["query"]
	}
	if app.IsEmpty(_query) {
		msg, _ := app.i18n.T("no-query", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	data := map[string]any{}
	switch _type := _query.(type) {
	case string:
		return app.CrudRunQuery(params, _query.(string), newDB)
	case []any:
		for i, query := range _query.([]any) {
			data[fmt.Sprintf(`query-%d`, i)] = app.CrudRunQuery(params, query.(string), newDB)
		}
	case []string:
		for i, query := range _query.([]string) {
			data[fmt.Sprintf(`query-%d`, i)] = app.CrudRunQuery(params, query, newDB)
		}
	case map[string]any:
		for key, query := range _query.(map[string]any) {
			data[key] = app.CrudRunQuery(params, query.(string), newDB)
		}
	case map[string]string:
		for key, query := range _query.(map[string]string) {
			data[key] = app.CrudRunQuery(params, query, newDB)
		}
	default:
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
			"type":    _type,
		}
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}
