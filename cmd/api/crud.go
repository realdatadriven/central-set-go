package main

import (
	"fmt"
	"path/filepath"

	"github.com/realdatadriven/etlx"
)

func (app *application) read(params map[string]interface{}) map[string]interface{} {
	// DATABASE
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	newDB, _, _, err := app.db.FromParams(params, _extra_conf)
	//fmt.Println("FromParams:", _driver, _database)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	if newDB != nil {
		defer newDB.Close()
	}
	tables := []interface{}{}
	if !app.IsEmpty(params["data"].(map[string]interface{})["table"]) {
		value := params["data"].(map[string]interface{})["table"]
		switch value.(type) {
		case nil:
			_ = true
		case string:
			tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
		case []interface{}:
			_tables := params["data"].(map[string]interface{})["table"].([]interface{})
			fmt.Println(_tables)
			for t := 0; t < len(_tables); t++ {
				if _, ok := _tables[t].(string); ok {
					tables = append(tables, _tables[t].(string))
				}
			}
		case map[interface{}]interface{}:
			// pass
		default:
			_ = true
		}
	} else if !app.IsEmpty(params["data"].(map[string]interface{})["tables"]) {
		value := params["data"].(map[string]interface{})["tables"]
		switch value.(type) {
		case string:
			tables = append(tables, params["data"].(map[string]interface{})["tables"].(string))
		case []interface{}:
			_tables := params["data"].(map[string]interface{})["tables"].([]interface{})
			for t := 0; t < len(_tables); t++ {
				tables = append(tables, _tables[t].(string))
			}
		default:
			_ = true
		}
	}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", map[string]interface{}{})
		return map[string]interface{}{
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
		_schemas = _schemas["data"].(map[string]interface{})
		params["schemas"] = _schemas
	}
	_permissions := app.table_access(params, tables)
	if !_permissions["success"].(bool) {
		return _permissions
	}
	if _, ok := _permissions["data"]; ok {
		_permissions = _permissions["data"].(map[string]interface{})
	} else {
		_permissions = map[string]interface{}{}
	}
	_row_level_tables := app.row_level_tables(params)
	if !_row_level_tables["success"].(bool) {
		return _row_level_tables
	}
	params["row_level_tables"] = []string{}
	if _, ok := _row_level_tables["tables"]; ok {
		params["row_level_tables"] = _row_level_tables["tables"].([]string)
	}
	data := map[string]interface{}{}
	for _, table := range tables {
		params["schema"] = map[string]interface{}{}
		if _, ok := _schemas[table.(string)]; ok {
			params["schema"] = _schemas[table.(string)].(map[string]interface{})
		}
		params["permissions"] = map[string]interface{}{}
		if _, ok := _permissions[table.(string)]; ok {
			params["permissions"] = _permissions[table.(string)].(map[string]interface{})
		}
		if len(tables) == 1 {
			return app.CrudRead(params, table.(string), newDB)
		}
		data[table.(string)] = app.CrudRead(params, table.(string), newDB)
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
		//"tables":       tables,
		//"_schemas":     _schemas,
		//"permissions": _permissions,
		//"sql": query,
	}
}

func (app *application) create_update(params map[string]interface{}) map[string]interface{} {
	// DATABASE
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	newDB, _, _, err := app.db.FromParams(params, _extra_conf)
	//fmt.Println("FromParams:", _driver, _database)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	if newDB != nil {
		defer newDB.Close()
	}
	table := ""
	if _, ok := params["data"].(map[string]interface{})["table"]; ok {
		table = params["data"].(map[string]interface{})["table"].(string)
	}
	if table == "" {
		msg, _ := app.i18n.T("no-table", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	tables := []interface{}{table}
	var _data interface{}
	if _, ok := params["data"].(map[string]interface{})["data"]; ok {
		_data = params["data"].(map[string]interface{})["data"]
	}
	switch _data.(type) {
	case []map[string]interface{}:
		for _, d := range _data.([]map[string]interface{}) {
			if _, ok := d["_table"]; ok {
				if app.contains(tables, d["_table"].(string)) {
					tables = append(tables, d["_table"].(string))
				}
			}
		}
	default:
		_ = ""
	}
	_schemas := app.tables(params, tables)
	if !_schemas["success"].(bool) {
		return _schemas
	}
	if _, ok := _schemas["data"]; ok {
		_schemas = _schemas["data"].(map[string]interface{})
		params["schemas"] = _schemas
	}
	// fmt.Println("TABLES TO CREATE:", tables, _schemas)
	_permissions := app.table_access(params, tables)
	if !_permissions["success"].(bool) {
		return _permissions
	}
	if _, ok := _permissions["data"]; ok {
		_permissions = _permissions["data"].(map[string]interface{})
	} else {
		_permissions = map[string]interface{}{}
	}
	_row_level_tables := app.row_level_tables(params)
	if !_row_level_tables["success"].(bool) {
		return _row_level_tables
	}
	params["row_level_tables"] = []string{}
	if _, ok := _row_level_tables["tables"]; ok {
		params["row_level_tables"] = _row_level_tables["tables"].([]string)
	}
	data := map[string]interface{}{}
	switch _data.(type) {
	case []interface{}:
		for i, d := range _data.([]interface{}) {
			tbl := table
			if _, ok := d.(map[string]interface{})["_table"]; ok {
				tbl = d.(map[string]interface{})["_table"].(string)
			}
			params["schema"] = map[string]interface{}{}
			if _, ok := _schemas[tbl]; ok {
				params["schema"] = _schemas[tbl].(map[string]interface{})
			}
			params["permissions"] = map[string]interface{}{}
			if _, ok := _permissions[tbl]; ok {
				params["permissions"] = _permissions[tbl].(map[string]interface{})
			}
			//fmt.Println(i, tbl, d)
			params["data"].(map[string]interface{})["data"] = d
			data[fmt.Sprintf(`row-%s-%d`, tbl, i)] = app.CrudCreateUpdte(params, tbl, newDB)
		}
	case []map[string]interface{}:
		for i, d := range _data.([]map[string]interface{}) {
			tbl := table
			if _, ok := d["_table"]; ok {
				tbl = d["_table"].(string)
			}
			params["schema"] = map[string]interface{}{}
			if _, ok := _schemas[tbl]; ok {
				params["schema"] = _schemas[tbl].(map[string]interface{})
			}
			params["permissions"] = map[string]interface{}{}
			if _, ok := _permissions[tbl]; ok {
				params["permissions"] = _permissions[tbl].(map[string]interface{})
			}
			params["data"].(map[string]interface{})["data"] = d
			data[fmt.Sprintf(`row-%s-%d`, tbl, i)] = app.CrudCreateUpdte(params, tbl, newDB)
			fmt.Println(fmt.Sprintf(`row-%s-%d`, tbl, i), data[fmt.Sprintf(`row-%s-%d`, tbl, i)].(map[string]interface{})["msg"])
		}
	case map[string]interface{}:
		params["schema"] = map[string]interface{}{}
		if _, ok := _schemas[table]; ok {
			params["schema"] = _schemas[table].(map[string]interface{})
		}
		params["permissions"] = map[string]interface{}{}
		if _, ok := _permissions[table]; ok {
			params["permissions"] = _permissions[table].(map[string]interface{})
		}
		return app.CrudCreateUpdte(params, table, newDB)
	default:
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
		}
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}

func (app *application) query(params map[string]interface{}) map[string]interface{} {
	// DATABASE
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	var newDB etlx.DBInterface
	newDB, driver, _database, err := app.db.FromParams(params, _extra_conf)
	//fmt.Println("FromParams DB:", driver, _database)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	if newDB != nil {
		defer newDB.Close()
	}
	if driver == "duckdb" {
		_db_ext := filepath.Ext(_database)
		// fmt.Println(_database, _db_ext)
		if _db_ext != "" {
			_db_ext = ""
		}
		newDB, err = etlx.NewDuckDB(fmt.Sprintf(`database/%s%s`, _database, _db_ext))
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		defer newDB.Close()
	}
	_data := map[string]interface{}{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]interface{})
	}
	if app.IsEmpty(_data) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	var _query interface{}
	if _, ok := _data["query"]; ok {
		_query = _data["query"]
	}
	if app.IsEmpty(_query) {
		msg, _ := app.i18n.T("no-query", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	data := map[string]interface{}{}
	switch _type := _query.(type) {
	case string:
		return app.CrudRunQuery(params, _query.(string), newDB)
	case []interface{}:
		for i, query := range _query.([]interface{}) {
			data[fmt.Sprintf(`query-%d`, i)] = app.CrudRunQuery(params, query.(string), newDB)
		}
	case []string:
		for i, query := range _query.([]string) {
			data[fmt.Sprintf(`query-%d`, i)] = app.CrudRunQuery(params, query, newDB)
		}
	case map[string]interface{}:
		for key, query := range _query.(map[string]interface{}) {
			data[key] = app.CrudRunQuery(params, query.(string), newDB)
		}
	case map[string]string:
		for key, query := range _query.(map[string]string) {
			data[key] = app.CrudRunQuery(params, query, newDB)
		}
	default:
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
			"type":    _type,
		}
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}
