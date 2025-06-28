package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/password"

	"github.com/realdatadriven/etlx"
)

func (app *application) CrudCreateUpdte(params map[string]interface{}, table string, db etlx.DBInterface) map[string]interface{} {
	/*var user_id int
	if _, ok := params["user"].(map[string]interface{})["user_id"]; ok {
		user_id = int(params["user"].(map[string]interface{})["user_id"].(float64))
	}*/
	var role_id int
	if _, ok := params["user"].(map[string]interface{})["role_id"]; ok {
		role_id = int(params["user"].(map[string]interface{})["role_id"].(float64))
	}
	/*var app_id int
	if _, ok := params["app"].(map[string]interface{})["app_id"]; ok {
		app_id = int(params["app"].(map[string]interface{})["app_id"].(float64))
	}*/
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	//fmt.Println(user_id, role_id, app_id)
	_schema := map[string]interface{}{}
	if _, ok := params["schema"]; ok {
		_schema = params["schema"].(map[string]interface{})
	}
	_permissions := map[string]interface{}{}
	if _, ok := params["permissions"]; ok {
		_permissions = params["permissions"].(map[string]interface{})
	}
	pk := ""
	if _, ok := _schema["pk"]; ok {
		pk = _schema["pk"].(string)
	}
	crud_aciton := "create"
	_data := map[string]interface{}{}
	if _, ok := params["data"].(map[string]interface{})["data"]; ok {
		_data = params["data"].(map[string]interface{})["data"].(map[string]interface{})
	}
	if _, ok := _data[pk]; ok {
		_to_delete := false
		if _, ok := _data["_to_delete"]; ok {
			_to_delete = _data["_to_delete"].(bool)
		}
		excluded := false
		if _, ok := _data["excluded"]; ok {
			excluded = _data["excluded"].(bool)
		}
		if _to_delete {
			crud_aciton = "delete"
		} else if excluded {
			crud_aciton = "delete"
		} else {
			query := fmt.Sprintf(`SELECT "%s" FROM "%s" WHERE "%s" = ?`, pk, table, pk)
			queryParams := []interface{}{_data[pk]}
			_pk_exists, _, err := db.QuerySingleRow(query, queryParams...)
			if err != nil {
				fmt.Println(0, query, err)
			} else if _, ok := (*_pk_exists)[pk]; ok {
				// fmt.Println(1, query, (*_pk_exists))
				crud_aciton = "update"
			} else {
				// fmt.Println(2, query, "NO RESULTS!")
			}
		}
	}
	roles := []interface{}{role_id}
	if !app.contains(roles, 1) {
		if _, ok := _permissions["read"]; !ok {
			msg, _ := app.i18n.T("no-table-access", map[string]interface{}{
				"table": table,
			})
			return map[string]interface{}{
				"success": false,
				"msg":     msg,
			}
		} else if !app.contains([]interface{}{true, 1}, _permissions["read"]) {
			msg, _ := app.i18n.T("no-table-action-access", map[string]interface{}{
				"table":  table,
				"action": strings.ToUpper(crud_aciton),
			})
			return map[string]interface{}{
				"success": false,
				"msg":     msg,
			}
		}
	}
	_errs := []string{}
	/*_row_level_tables := []string{}
	if _, ok := params["row_level_tables"]; ok {
		_row_level_tables = params["row_level_tables"].([]string)
	}*/
	// FIELDS
	if _, ok := _schema["fields"].(map[string]interface{}); ok {
		for field, field_data := range _schema["fields"].(map[string]interface{}) {
			_type := field_data.(map[string]interface{})["type"].(string)
			_nullable := true
			if null, ok := field_data.(map[string]interface{})["nullable"]; ok {
				if app.contains([]interface{}{0, false, "0", "false", "False", "FALSE"}, null) {
					_nullable = false
				}
			}
			_type = strings.ToLower(_type)
			_value := _data[field]
			if app.contains([]interface{}{"datetime", "date"}, _type) {
				// TREAT DATE AND TIME TYPES
			}
			enable_user := []interface{}{}
			for _, t := range strings.Split(app.config.enable_user, ",") {
				enable_user = append(enable_user, t)
			}
			if app.contains([]interface{}{"created_at", "updated_at"}, field) {
				if _, ok := _data[pk]; ok && field == "created_at" && crud_aciton != "create" {
				} else {
					_data[field] = time.Now()
				}
			} else if app.contains([]interface{}{"excluded"}, field) {
				if _, ok := _data[pk]; ok {
				} else {
					_data[field] = false
				}
			} else if app.contains([]interface{}{"password", "pass"}, field) {
				if _, ok := _data[field]; !ok {
					continue
				} else if _, ok := _data[pk]; !ok || crud_aciton == "create" {
					hashedPassword, err := password.Hash(_data[field].(string))
					if err != nil {
						return map[string]interface{}{
							"success": true,
							"msg":     "Error hashing password!",
						}
					}
					_data[field] = hashedPassword
				} else if len(_data[field].(string)) < 20 {
					hashedPassword, err := password.Hash(_data[field].(string))
					if err != nil {
						return map[string]interface{}{
							"success": true,
							"msg":     "Error hashing password!",
						}
					}
					_data[field] = hashedPassword
				}
			} else if app.contains([]interface{}{"app", "app_id"}, field) && !app.contains([]interface{}{"app", "role_app", "role_app_menu", "role_app_menu_table"}, table) {
				if _, ok := _data[field]; !ok && crud_aciton == "create" {
					_data[field] = params["app"].(map[string]interface{})[field]
				}
			} else if app.contains([]interface{}{"user", "user_id"}, field) && !app.contains([]interface{}{"user", "user_role", "column_level_access", "row_level_access"}, table) && !app.contains(enable_user, table) {
				if _, ok := _data[field]; !ok && crud_aciton == "create" {
					_data[field] = params["user"].(map[string]interface{})[field]
				}
			} else if !_nullable && field != pk && crud_aciton != "delete" {
				if !app.IsEmpty(_data[field]) {
				} else if field == "lang" {
					_data[field] = lang
				} else if app.contains([]interface{}{"db", "database"}, field) {
					_data[field] = params["app"].(map[string]interface{})["db"]
				} else {
					msg, _ := app.i18n.T("field-required", map[string]interface{}{"field": field})
					_errs = append(_errs, msg)
				}
			} else {
				switch _value.(type) {
				case map[string]interface{}:
					_json, err := json.Marshal(_value)
					if err != nil {
						fmt.Println(field, "unable to convert to JSON!", err)
					}
					_data[field] = _json
				case []map[string]interface{}:
					_json, err := json.Marshal(_value)
					if err != nil {
						fmt.Println(field, "unable to convert to JSON!", err)
					}
					_data[field] = _json
				case []interface{}:
					_json, err := json.Marshal(_value)
					if err != nil {
						fmt.Println(field, "unable to convert to JSON!", err)
					}
					_data[field] = _json
				default:
					//
				}
			}
			//fmt.Println(field, _type, _value)
		}
	}
	if len(_errs) > 0 {
		msg, _ := app.i18n.T("validation-errors", map[string]interface{}{"n": len(_errs)})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
			"errors":  _errs,
		}
	}
	// fmt.Println(crud_aciton)
	// REMOVE FIELDS THAT IS NOT IN THE TABLE SCHEMA
	_aux_data := _data
	for key := range _aux_data {
		if _, ok := _schema["fields"].(map[string]interface{}); ok {
			if _, ok := _schema["fields"].(map[string]interface{})[key]; !ok {
				delete(_data, key)
			}
		}
	}
	// CREATE | UPDATE | DELETE
	var keys []interface{}
	for key := range _data {
		keys = append(keys, key)
	}
	cols := app.joinSlice(keys, `", "`)
	vals := app.joinSlice(keys, `, :`)
	_pg_returning := ""
	if db.GetDriverName() == "postgres" && pk != "" {
		_pg_returning = fmt.Sprintf(` RETURNING "%s"`, pk)
	}
	query := fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES (:%s)%s`, table, cols, vals, _pg_returning)
	if crud_aciton != "create" {
		keys = []interface{}{}
		for key := range _data {
			keys = append(keys, fmt.Sprintf(`"%s" = :%s`, key, key))
		}
		cols := app.joinSlice(keys, `, `)
		query = fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = :%s`, table, cols, pk, pk)
		if crud_aciton == "delete" {
			permanently := false
			if _, ok := _aux_data["permanently"]; ok {
				if app.contains([]interface{}{true, 1, "true", "True", "TRUE"}, _aux_data["permanently"]) {
					permanently = true
				}
			}
			if _, ok := _schema["fields"].(map[string]interface{})["excluded"]; ok && !permanently {
				query = fmt.Sprintf(`UPDATE "%s" SET "excluded" = TRUE WHERE "%s" = :%s`, table, pk, pk)
			} else {
				query = fmt.Sprintf(`DELETE FROM "%s" WHERE "%s" = :%s`, table, pk, pk)
			}
		}
	}
	id := 0
	if db.GetDriverName() == "postgres" && strings.HasPrefix(query, "INSERT") {
		_id, err := db.ExecuteQueryPGInsertWithLastInsertId(query, _data)
		//fmt.Println("ExecuteQueryPGInsertWithLastInsertId", id, query)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				_sql := fmt.Sprintf(`SELECT SETVAL(PG_GET_SERIAL_SEQUENCE('%s', '%s'), NEXTVAL(PG_GET_SERIAL_SEQUENCE('%s', '%s')), FALSE)`, table, pk, table, pk)
				fmt.Println("PG_GET_SERIAL_SEQUENCE:", _sql)
				_, err2 := db.ExecuteQuery(_sql)
				if err2 != nil {
					fmt.Println("Err tring to increment pg id: ", err2)
				}
				_id, err = db.ExecuteQueryPGInsertWithLastInsertId(query, _data)
				if err != nil {
					return map[string]interface{}{
						"success": false,
						"table":   table,
						"pk":      pk,
						"msg":     fmt.Sprintf("%s", err),
					}
				}
			} else {
				return map[string]interface{}{
					"success": false,
					"table":   table,
					"pk":      pk,
					"msg":     fmt.Sprintf("%s", err),
				}
			}
		}
		id = _id
	} else {
		_id, err := db.ExecuteNamedQuery(query, _data)
		// fmt.Println(query)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"table":   table,
				"pk":      pk,
				//"data":    _data,
				//"sql":     query,
				"msg": fmt.Sprintf("%s", err),
			}
		}
		id = _id
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success":              true,
		"msg":                  msg,
		"pk":                   pk,
		"table":                table,
		"id":                   id,
		"inserted_primary_key": id,
		"data":                 _data,
		"sql":                  query,
	}
}
