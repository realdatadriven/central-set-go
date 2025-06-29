package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/realdatadriven/etlx"
)

func (app *application) table_access(params map[string]interface{}, tables []interface{}) map[string]interface{} {
	var user_id int
	if _, ok := params["user"].(map[string]interface{})["user_id"]; ok {
		user_id = int(params["user"].(map[string]interface{})["user_id"].(float64))
	}
	var role_id int
	if _, ok := params["app"].(map[string]interface{})["role_id"]; ok {
		role_id = int(params["app"].(map[string]interface{})["role_id"].(float64))
	}
	var app_id int
	if _, ok := params["app"].(map[string]interface{})["app_id"]; ok {
		app_id = int(params["app"].(map[string]interface{})["app_id"].(float64))
	}
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	allTables := false
	if app.IsEmpty(tables) {
		tables = []interface{}{}
		if !app.IsEmpty(params["data"].(map[string]interface{})["table"]) {
			value := params["data"].(map[string]interface{})["table"]
			switch value.(type) {
			case nil:
				// pass
			case string:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			case []interface{}:
				_tables := params["data"].(map[string]interface{})["table"].([]interface{})
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			case map[interface{}]interface{}:
				// pass
			default:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			}
		} else if !app.IsEmpty(params["data"].(map[string]interface{})["tables"]) {
			value := params["data"].(map[string]interface{})["tables"]
			switch value.(type) {
			case string:
				tables = append(tables, params["data"].(map[string]interface{})["tables"].(string))
			case []interface{}:
				_tables := params["data"].(map[string]interface{})["tables"].([]interface{})
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			default:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			}
		}
		if app.IsEmpty(tables) {
			// fmt.Println("GET ALL TABLES!")
			result, _, err := newDB.AllTables(params, _extra_conf)
			if err != nil {
				return map[string]interface{}{
					"success": false,
					"msg":     fmt.Sprintf("%s", err),
				}
			}
			for _, row := range *result {
				tables = append(tables, string(row["name"].(string)))
			}
			allTables = true
		}
	}
	data := map[string]interface{}{}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
			"tables":  tables,
		}
	} else {
		//fmt.Println(user_id, role_id)
		query := `SELECT DISTINCT user_role.role_id
		FROM user_role
		JOIN role ON user_role.role_id = role.role_id
		WHERE user_role.user_id = $1
			AND user_role.excluded = FALSE
			AND role.excluded = FALSE`
		var queryParams []interface{}
		queryParams = append(queryParams, user_id)
		result, _, err := app.db.QueryMultiRows(query, queryParams...)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		roles := []interface{}{}
		roles = append(roles, role_id)
		for _, row := range *result {
			roles = append(roles, int(row["role_id"].(float64)))
		}
		queryParams = []interface{}{app_id}
		queryParams = append(queryParams, roles)
		query = `SELECT role_app_menu_table.*, "table"."table"
		FROM role_app_menu_table
		JOIN "table" ON "table".table_id = role_app_menu_table.table_id
		WHERE role_app_menu_table.app_id = ?
			AND role_app_menu_table.role_id IN (?)
			AND "table"."table" IN (?)
			AND role_app_menu_table.excluded = FALSE
			AND "table".excluded = FALSE`
		if allTables {
			query = `SELECT role_app_menu_table.*, "table"."table"
			FROM role_app_menu_table
			JOIN "table" ON "table".table_id = role_app_menu_table.table_id
			WHERE role_app_menu_table.app_id = ?
				AND role_app_menu_table.role_id IN (?)
				AND role_app_menu_table.excluded = FALSE
				AND "table".excluded = FALSE`
		} else {
			queryParams = append(queryParams, tables)
		}
		query, args, err := sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		result, _, err = app.db.QueryMultiRows(query, args...)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		for _, row := range *result {
			data[row["table"].(string)] = row
		}
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}

func (app *application) row_level_access(params map[string]interface{}, tables []interface{}, row_id []interface{}) map[string]interface{} {
	var user_id int
	if _, ok := params["user"].(map[string]interface{})["user_id"]; ok {
		user_id = int(params["user"].(map[string]interface{})["user_id"].(float64))
	}
	var role_id int
	if _, ok := params["app"].(map[string]interface{})["role_id"]; ok {
		role_id = int(params["app"].(map[string]interface{})["role_id"].(float64))
	}
	var app_id int
	if _, ok := params["app"].(map[string]interface{})["app_id"]; ok {
		app_id = int(params["app"].(map[string]interface{})["app_id"].(float64))
	}
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	allTables := false
	if app.IsEmpty(tables) {
		tables = []interface{}{}
		if !app.IsEmpty(params["data"].(map[string]interface{})["table"]) {
			value := params["data"].(map[string]interface{})["table"]
			switch value.(type) {
			case nil:
				// pass
			case string:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			case []interface{}:
				_tables := params["data"].(map[string]interface{})["table"].([]interface{})
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			case map[interface{}]interface{}:
				// pass
			default:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			}
		} else if !app.IsEmpty(params["data"].(map[string]interface{})["tables"]) {
			value := params["data"].(map[string]interface{})["tables"]
			switch value.(type) {
			case string:
				tables = append(tables, params["data"].(map[string]interface{})["tables"].(string))
			case []interface{}:
				_tables := params["data"].(map[string]interface{})["tables"].([]interface{})
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			default:
				tables = append(tables, params["data"].(map[string]interface{})["table"].(string))
			}
		}
		if app.IsEmpty(tables) {
			// fmt.Println("GET ALL TABLES!")
			result, _, err := newDB.AllTables(params, _extra_conf)
			if err != nil {
				return map[string]interface{}{
					"success": false,
					"msg":     fmt.Sprintf("%s", err),
				}
			}
			defer newDB.Close()
			for _, row := range *result {
				tables = append(tables, string(row["name"].(string)))
			}
			allTables = true
		}
	}
	data := map[string]interface{}{}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
			"tables":  tables,
		}
	} else {
		if app.IsEmpty(role_id) {
			row_id = []interface{}{}
			if !app.IsEmpty(params["data"].(map[string]interface{})["row_id"]) {
				value := params["data"].(map[string]interface{})["row_id"]
				switch value.(type) {
				case string:
					row_id = append(row_id, params["data"].(map[string]interface{})["row_id"].(string))
				case []interface{}:
					_row_ids := params["data"].(map[string]interface{})["row_id"].([]interface{})
					for t := 0; t < len(_row_ids); t++ {
						row_id = append(row_id, _row_ids[t])
					}
				default:
					tables = append(tables, params["data"].(map[string]interface{})["row_id"].(string))
				}
			}
		}
		//fmt.Println(user_id, role_id)
		query := `SELECT DISTINCT user_role.role_id
		FROM user_role
		JOIN role ON user_role.role_id = role.role_id
		WHERE user_role.user_id = $1
			AND user_role.excluded = FALSE
			AND role.excluded = FALSE`
		var queryParams []interface{}
		queryParams = append(queryParams, user_id)
		result, _, err := app.db.QueryMultiRows(query, queryParams...)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		roles := []interface{}{}
		roles = append(roles, role_id)
		for _, row := range *result {
			roles = append(roles, int(row["role_id"].(float64)))
		}
		queryParams = []interface{}{app_id}
		queryParams = append(queryParams, roles)
		_get_table_lists := ""
		if allTables {
			_get_table_lists = `AND role_row_level_access."table" IN (?)`
			queryParams = append(queryParams, tables)
		}
		_get_row_id_lists := ""
		if !app.IsEmpty(row_id) {
			_get_row_id_lists = `AND role_row_level_access.row_id IN (?)`
			queryParams = append(queryParams, row_id)
		}
		query = fmt.Sprintf(`SELECT role_row_level_access.*
		FROM role_row_level_access
		JOIN (
			SELECT "table", role_id, app_id, row_id, MAX("updated_at") AS "max_updated_at"
			FROM role_row_level_access
			GROUP BY "table", role_id, app_id, row_id
		) AS "T" ON (
			"T"."table" = role_row_level_access."table"
			AND "T"."role_id" = role_row_level_access."role_id"
			AND "T"."app_id" = role_row_level_access."app_id"
			AND "T"."row_id" = role_row_level_access."row_id"
			AND "T"."max_updated_at" = role_row_level_access."updated_at"
		)
		WHERE role_row_level_access.app_id = ?
			AND role_row_level_access.role_id IN (?) %s %s 
			AND role_row_level_access.excluded = FALSE`, _get_table_lists, _get_row_id_lists)
		query, args, err := sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		result, _, err = app.db.QueryMultiRows(query, args...)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		for _, row := range *result {
			if _, ok := data[row["table"].(string)]; !ok {
				data[row["table"].(string)] = []map[string]interface{}{}
			}
			_aux := row
			_aux["org"] = "role_row_level_access"
			data[row["table"].(string)] = append(data[row["table"].(string)].([]map[string]interface{}), _aux)
		}
		// row_level_access
		queryParams = []interface{}{app_id, user_id}
		if allTables {
			queryParams = append(queryParams, tables)
		}
		if !app.IsEmpty(row_id) {
			queryParams = append(queryParams, row_id)
		}
		query = fmt.Sprintf(`SELECT row_level_access.*
		FROM row_level_access
		JOIN (
			SELECT "table", user_id, app_id, row_id, MAX("updated_at") AS "max_updated_at"
			FROM row_level_access
			GROUP BY "table", user_id, app_id, row_id
		) AS "T" ON (
			"T"."table" = row_level_access."table"
			AND "T"."user_id" = row_level_access."user_id"
			AND "T"."app_id" = row_level_access."app_id"
			AND "T"."row_id" = row_level_access."row_id"
			AND "T"."max_updated_at" = row_level_access."updated_at"
		)
		WHERE row_level_access.app_id = ?
			AND row_level_access.user_id = ? %s %s 
			AND row_level_access.excluded = FALSE`, _get_table_lists, _get_row_id_lists)
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		result, _, err = app.db.QueryMultiRows(query, args...)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		for _, row := range *result {
			if _, ok := data[row["table"].(string)]; !ok {
				data[row["table"].(string)] = []map[string]interface{}{}
			}
			_aux := row
			_aux["org"] = "row_level_access"
			if !app.IsEmpty(data[row["table"].(string)]) {
				_aux2 := app.filter(
					data[row["table"].(string)].([]map[string]interface{}),
					func(r map[string]interface{}) bool {
						_r_id := r["row_id"] == _aux["row_id"]
						_tbl := r["table"] == _aux["table"]
						_app_id := r["app_id"] == _aux["app_id"]
						_org := r["org"] == "role_row_level_access"
						return _r_id && _tbl && _app_id && _org
					},
				)
				fmt.Println(_aux2)
			}
			data[row["table"].(string)] = append(data[row["table"].(string)].([]map[string]interface{}), _aux)
		}
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}

func (app *application) row_level_tables(params map[string]interface{}) map[string]interface{} {
	var app_id int
	if _, ok := params["app"].(map[string]interface{})["app_id"]; ok {
		app_id = int(params["app"].(map[string]interface{})["app_id"].(float64))
	}
	tables := []string{}
	queryParams := []interface{}{app_id}
	query := `SELECT "table"."table"
	FROM menu_table
	JOIN "table" ON "table".table_id = menu_table.table_id
	WHERE menu_table.app_id = $1
		AND menu_table.requires_rla = TRUE
		AND menu_table.excluded = FALSE
		AND "table".excluded = FALSE
	UNION
	SELECT "table"."table"
	FROM menu_table
	JOIN "table" ON "table".table_id = menu_table.table_id
	WHERE menu_table.requires_rla = TRUE
		AND menu_table.excluded = FALSE
		AND "table".excluded = FALSE`
	result, _, err := app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	for _, row := range *result {
		tables = append(tables, row["table"].(string))
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"tables":  tables,
	}
}
