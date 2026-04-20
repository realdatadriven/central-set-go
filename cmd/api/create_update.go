package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/password"

	"github.com/realdatadriven/etlx"
)

//const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (app *application) randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (app *application) CrudCreateUpdte(params Dict, table string, db etlx.DBInterface) Dict {
	var user_id int
	if _, ok := params["user"].(Dict)["user_id"]; ok {
		user_id = int(params["user"].(Dict)["user_id"].(float64))
	}
	var role_id int
	if _, ok := params["user"].(Dict)["role_id"]; ok {
		role_id = int(params["user"].(Dict)["role_id"].(float64))
	}
	/*var app_id int
	if _, ok := params["app"].(Dict)["app_id"]; ok {
		app_id = int(params["app"].(Dict)["app_id"].(float64))
	}*/
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	//fmt.Println(user_id, role_id, app_id)
	_schema := Dict{}
	if _, ok := params["schema"]; ok {
		_schema = params["schema"].(Dict)
	}
	_permissions := Dict{}
	if _, ok := params["permissions"]; ok {
		_permissions = params["permissions"].(Dict)
	}
	// fmt.Println("PK", _schema["pk"])
	pk := ""
	if _, ok := _schema["pk"]; ok {
		pk = _schema["pk"].(string)
	}
	crud_aciton := "create"
	_data := Dict{}
	if _, ok := params["data"].(Dict)["data"]; ok {
		_data = params["data"].(Dict)["data"].(Dict)
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
			queryParams := []any{_data[pk]}
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
	roles := []any{role_id}
	if !app.contains(roles, 1) {
		//fmt.Println("3 PERMISSIONS:", _permissions)
		if _, ok := _permissions[crud_aciton]; !ok {
			msg, _ := app.i18n.T("no-table-access", Dict{
				"table": table,
			})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		} else if !app.contains([]any{true, 1}, _permissions[crud_aciton]) {
			msg, _ := app.i18n.T("no-table-action-access", Dict{
				"table":  table,
				"action": strings.ToUpper(crud_aciton),
			})
			return Dict{
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
	if _, ok := _schema["fields"].(Dict); ok {
		for field, field_data := range _schema["fields"].(Dict) {
			_type := field_data.(Dict)["type"].(string)
			_nullable := true
			if null, ok := field_data.(Dict)["nullable"]; ok {
				if app.contains([]any{0, false, "0", "false", "False", "FALSE"}, null) {
					_nullable = false
				}
			}
			_type = strings.ToLower(_type)
			_value := _data[field]
			if app.contains([]any{"datetime", "date"}, _type) {
				// TREAT DATE AND TIME TYPES
			}
			enable_user := []any{}
			for _, t := range strings.Split(app.config.enable_user, ",") {
				enable_user = append(enable_user, t)
			}
			if app.contains([]any{"created_at", "updated_at"}, field) {
				if _, ok := _data[pk]; ok && field == "created_at" && crud_aciton != "create" {
				} else {
					_data[field] = time.Now()
				}
			} else if app.contains([]any{"excluded"}, field) {
				if _, ok := _data[pk]; ok {
				} else {
					_data[field] = false
				}
			} else if app.contains([]any{"password", "pass"}, field) {
				if _, ok := _data[field]; !ok {
					continue
				} else if _, ok := _data[pk]; !ok || crud_aciton == "create" {
					hashedPassword, err := password.Hash(_data[field].(string))
					if err != nil {
						return Dict{
							"success": true,
							"msg":     "Error hashing password!",
						}
					}
					_data[field] = hashedPassword
				} else if len(_data[field].(string)) < 20 {
					hashedPassword, err := password.Hash(_data[field].(string))
					if err != nil {
						return Dict{
							"success": true,
							"msg":     "Error hashing password!",
						}
					}
					_data[field] = hashedPassword
				}
			} else if app.contains([]any{"app", "app_id"}, field) && !app.contains([]any{"app", "users", "role_app", "role_app_menu", "role_app_menu_table"}, table) {
				if _, ok := _data[field]; !ok && crud_aciton == "create" {
					_data[field] = params["app"].(Dict)[field]
				}
			} else if app.contains([]any{"user", "username", "user_id"}, field) && !app.contains([]any{"user", "users", "user_role", "column_level_access", "row_level_access"}, table) && !app.contains(enable_user, table) {
				if _, ok := _data[field]; !ok && crud_aciton == "create" {
					_data[field] = params["user"].(Dict)[field]
				}
			} else if !_nullable && field != pk && crud_aciton != "delete" {
				if !app.IsEmpty(_data[field]) {
				} else if field == "lang" {
					_data[field] = lang
				} else if app.contains([]any{"db", "database"}, field) {
					_data[field] = params["app"].(Dict)["db"]
				} else {
					msg, _ := app.i18n.T("field-required", Dict{"field": field})
					_errs = append(_errs, msg)
				}
			} else {
				switch _value.(type) {
				case Dict:
					_json, err := json.Marshal(_value)
					if err != nil {
						fmt.Println(field, "unable to convert to JSON!", err)
					}
					_data[field] = _json
				case []Dict:
					_json, err := json.Marshal(_value)
					if err != nil {
						fmt.Println(field, "unable to convert to JSON!", err)
					}
					_data[field] = _json
				case []any:
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
		msg, _ := app.i18n.T("validation-errors", Dict{"n": len(_errs)})
		return Dict{
			"success": false,
			"msg":     msg,
			"errors":  _errs,
		}
	}
	// CHECK ROW LEVEL ACCESS
	if !app.contains(roles, 1) {
		fk_tables_added := []any{}
		fk_tables_pk := Dict{}
		if _, ok := _schema["fields"]; ok {
			for _, field_data := range _schema["fields"].(Dict) {
				if _, ok := field_data.(Dict)["fk"]; !ok {
				} else if field_data.(Dict)["fk"].(bool) || field_data.(Dict)["fk"] == any(1) {
					referred_table := ""
					if _, ok := field_data.(Dict)["referred_table"]; ok {
						referred_table = field_data.(Dict)["referred_table"].(string)
					}
					referred_column := field_data.(Dict)["referred_column"]
					fk_tables_pk[referred_table] = referred_column
					fk_tables_added = append(fk_tables_added, referred_table)
				}
			}
		}
		fk_tables_pk[table] = pk
		_row_level_tables := []string{}
		rla_tables_ids := Dict{}
		if _, ok := params["row_level_tables"]; ok {
			_row_level_tables = params["row_level_tables"].([]string)
			_tables_to_chk := []any{}
			for _, v := range append(fk_tables_added, table) {
				if app.contains(app.sliceStrs2SliceInterfaces(_row_level_tables), v) {
					_tables_to_chk = append(_tables_to_chk, v)
				}
			}
			// fmt.Println("row_level_tables:", _row_level_tables, fk_tables_added, "_tables_to_chk:", _tables_to_chk)
			if len(_tables_to_chk) > 0 {
				params["fk_tables_pk"] = fk_tables_pk
				params["get_users_owned_id"] = true
				rla_access := app.row_level_access(params, _tables_to_chk, []any{})
				params["get_users_owned_id"] = false
				users_owned_ids := Dict{}
				if _, ok := rla_access["data"].(Dict); ok {
					users_owned_ids = rla_access["data"].(Dict)
				}
				//fmt.Println("rla_access:", _tables_to_chk, rla_access)
				if !rla_access["success"].(bool) {
					return rla_access
				} else if _rla_access_data, ok := rla_access["data"].(Dict); ok {
					for key, val := range _rla_access_data {
						//fmt.Println(key, val.([]Dict))
						my_ids := []any{}
						if _, ok := users_owned_ids[key].([]any); ok {
							my_ids = users_owned_ids[key].([]any)
						}
						rla_tables_ids[key] = app.getRLAIds(val.([]Dict), key, crud_aciton, my_ids)
					}
				} else {
					fmt.Println("DEBUG THIS SOMETHING WORONG WITH RLA(READ):", rla_access)
				}
				for _, t := range _tables_to_chk {
					if _, ok := rla_tables_ids[t.(string)].([]any); !ok {
						rla_tables_ids[t.(string)] = []any{}
					}
				}
				//fmt.Println(fk_tables_pk, rla_tables_ids)
				if _, ok := params["data"].(Dict)["filters"]; !ok {
					params["data"].(Dict)["filters"] = []any{}
				}
				fk_tables_pk[table] = pk
				for rla_t, ids := range rla_tables_ids {
					row_id_field := fk_tables_pk[rla_t].(string)
					if _, ok := _data[pk]; rla_t == table && ok && any(user_id) == _data["user_id"] { // VERIFY LATER ON
						continue
					}
					if row_id, ok := _data[row_id_field]; ok {
						if !app.containsInt(ids.([]any), row_id) {
							msg, _ := app.i18n.T("no-row-level-access", Dict{
								"table":  rla_t,
								"row_id": row_id,
								"action": crud_aciton,
							})
							return Dict{
								"success": false,
								"msg":     msg,
							}
						}
					}
				}
			}
		}
	}
	// VALIDATIONS
	/*table: validation*/
	etlx_engine := &etlx.ETLX{}
	_, database, _ := app.GetDBNameFromParams(params)
	//crud_aciton, table
	validation_data := []any{database, table}
	get_validations_sql := fmt.Sprintf(`SELECT * FROM "validation" WHERE "active" = TRUE AND "db" = ? AND "table" = ? AND "%" IS TRUE`, crud_aciton)
	validation_rows, err := app.AdminGetRowsByFilter(get_validations_sql, validation_data)
	if err != nil {
		fmt.Printf("Error occurred while fetching validations: %v", err, get_validations_sql)
		/*return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Error occurred while fetching validations: %v", err),
		}*/
	} else if len(validation_rows) > 0 {
		for _, validation := range validation_rows {
			if _, ok := validation["sql"]; ok {
				sql_rule := validation["sql"].(string)
				valid_reaction_id := validation["valid_reaction_id"]
				validation_log := Dict{
					"validation_id":   validation["validation_id"],
					"validation_code": validation["validation_code"],
					"validation":      validation["validation"],
					"table":           table,
					"db":              database,
					"action":          crud_aciton,
					"user_id":         user_id,
					"app_id":          params["app"].(Dict)["app_id"],
				}
				insert_validation_log_sql := `INSERT INTO "validation_logs" ("validation_id", "validation_code", "validation", "table", "db", "action", "success", "log_message", "user_id", "app_id", "executed_at", "created_at", "updated_at") 
				VALUES (:validation_id, :validation_code, :validation, :table, :db, :action, :success, :log_message, :user_id, :app_id, :executed_at, :created_at, :updated_at)`
				//fmt.Println("VALIDATION SQL:", sql_rule)
				var valid bool
				sql, _filters_opts, _ := etlx_engine.NamedToPositional(sql_rule, _data)
				res, err := app.AdminGetRowsByFilter(sql, _filters_opts)
				if err != nil {
					fmt.Printf("Error executing validation SQL for validation_id %v: %v", validation["validation_id"], err)
					valid = false
					validation_log["success"] = valid
					validation_log["log_message"] = fmt.Sprintf("Error executing validation SQL: %v", err)
					validation_log["executed_at"] = time.Now()
					validation_log["created_at"] = time.Now()
					validation_log["updated_at"] = time.Now()
					_, err = app.db.ExecuteNamedQuery(insert_validation_log_sql, validation_log)
					if err != nil {
						fmt.Printf("Error inserting validation log for validation_id %v: %v", validation["validation_id"], err)
					}
					return Dict{
						"success": false,
						"msg":     fmt.Sprintf("Error executing validation SQL: %v", err),
					}
				} else {
					// valid_reaction_id = 1 means throw erroe if len(*res) > 0 and valid_reaction_id = 2 means throw error if len(*res) == 0
					if app.toInt(valid_reaction_id) == 1 {
						valid = len(res) == 0
					} else if app.toInt(valid_reaction_id) == 2 {
						valid = len(res) > 0
					}
					if !valid {
						msg, err := etlx_engine.RenderTemplate(validation["err_msg"].(string), _data)
						if err != nil {
							msg = validation["err_msg"].(string)
							fmt.Println("Error rendering validation error message template for validation_id", validation["validation_id"], ":", err, validation["err_msg"])
						}
						validation_log["success"] = valid
						validation_log["log_message"] = fmt.Sprintf("Validation %s executed with result: %v. Message: %s", validation["validation_code"], valid, msg)
						validation_log["executed_at"] = time.Now()
						validation_log["created_at"] = time.Now()
						validation_log["updated_at"] = time.Now()
						_, err = app.db.ExecuteNamedQuery(insert_validation_log_sql, validation_log)
						if err != nil {
							fmt.Printf("Error inserting validation log for validation_id %v: %v", validation["validation_id"], err)
						}
						return Dict{
							"success": false,
							"msg":     msg,
						}
					}
				}
				validation_log["success"] = valid
				validation_log["log_message"] = fmt.Sprintf("Validation %s executed with result: %v", validation["validation_code"], valid)
				validation_log["executed_at"] = time.Now()
				validation_log["created_at"] = time.Now()
				validation_log["updated_at"] = time.Now()
				_, err = app.db.ExecuteNamedQuery(insert_validation_log_sql, validation_log)
				if err != nil {
					fmt.Printf("Error inserting validation log for validation_id %v: %v", validation["validation_id"], err)
				}
			}
		}
	}
	// fmt.Println(crud_aciton)
	// REMOVE FIELDS THAT IS NOT IN THE TABLE SCHEMA
	_aux_data := _data
	for key := range _aux_data {
		if _, ok := _schema["fields"].(Dict); ok {
			if _, ok := _schema["fields"].(Dict)[key]; !ok {
				delete(_data, key)
			}
		}
	}
	// CREATE | UPDATE | DELETE
	var keys []any
	for key := range _data {
		keys = append(keys, key)
	}
	cols := app.joinSlice(keys, `", "`)
	vals := app.joinSlice(keys, `, :`)
	_pg_returning := ""
	if db.GetDriverName() == "postgres" && pk != "" {
		_pg_returning = fmt.Sprintf(` RETURNING "%s"`, pk)
	}
	//
	query := fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES (:%s)%s`, table, cols, vals, _pg_returning)
	if crud_aciton != "create" {
		keys = []any{}
		for key := range _data {
			keys = append(keys, fmt.Sprintf(`"%s" = :%s`, key, key))
		}
		cols := app.joinSlice(keys, `, `)
		query = fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = :%s`, table, cols, pk, pk)
		if crud_aciton == "delete" {
			permanently := false
			if _, ok := _aux_data["permanently"]; ok {
				if app.contains([]any{true, 1, "true", "True", "TRUE"}, _aux_data["permanently"]) {
					permanently = true
				}
			}
			if _, ok := _schema["fields"].(Dict)["excluded"]; ok && !permanently {
				query = fmt.Sprintf(`UPDATE "%s" SET "excluded" = TRUE WHERE "%s" = :%s`, table, pk, pk)
			} else {
				query = fmt.Sprintf(`DELETE FROM "%s" WHERE "%s" = :%s`, table, pk, pk)
			}
		}
	}
	//fmt.Println(crud_aciton, pk, _data[pk], query)
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
					return Dict{
						"success": false,
						"table":   table,
						"pk":      pk,
						"msg":     fmt.Sprintf("%s", err),
					}
				}
			} else {
				return Dict{
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
		//fmt.Println(query)
		if err != nil {
			fmt.Println(crud_aciton, pk, _data[pk], query, err)
			return Dict{
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
	if os.Getenv("DYN_LOGIN_TABLE_MAP_TO_USERS") == "true" && crud_aciton == "create" {
		login_table := os.Getenv("DYN_LOGIN_TABLE")
		if login_table == table {
			user_id_field := os.Getenv("DYN_LOGIN_USER_ID_FIELD")
			dyn_login_role_id := os.Getenv("DYN_LOGIN_ROLE_ID")
			username_field := os.Getenv("DYN_LOGIN_USERNAME_FIELD")
			email_field := os.Getenv("DYN_LOGIN_EMAIL_FIELD")
			//password_field := os.Getenv("DYN_LOGIN_PASSWORD_FIELD")
			active_field := os.Getenv("DYN_LOGIN_ACTIVE_FIELD")
			if login_table == "" {
				msg, _ := app.i18n.T("login-table-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
			if user_id_field == "" {
				msg, _ := app.i18n.T("user-id-field-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
			if username_field == "" {
				msg, _ := app.i18n.T("username-field-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
			/*if password_field == "" {
				msg, _ := app.i18n.T("password-field-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}*/
			if email_field == "" {
				msg, _ := app.i18n.T("email-field-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
			if dyn_login_role_id == "" {
				msg, _ := app.i18n.T("role-id-required", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
			// generate random password
			pass, err := password.Hash(app.randomString(8))
			if err != nil {
				msg, _ := app.i18n.T("password-hash-error", Dict{})
				return Dict{"success": false, "msg": msg}
			}
			_data_user := Dict{
				"username":   _data[username_field],
				"first_name": _data[username_field],
				"last_name":  _data[username_field],
				"email":      _data[email_field],
				"role_id":    dyn_login_role_id,
				"password":   pass,
				"active":     _data[active_field],
				"created_at": time.Now(),
				"updated_at": time.Now(),
				"excluded":   false,
			}
			query_user := fmt.Sprintf(`INSERT INTO "users" ("username", "first_name", "last_name", "email", "role_id", "password", "active", "created_at", "updated_at", "excluded") 
			VALUES (:username, :first_name, :last_name, :email, :role_id, :password, :active, :created_at, :updated_at, :excluded)`)
			_, err = app.db.ExecuteNamedQuery(query_user, _data_user)
			if err != nil {
				fmt.Println("Error creating user for login table mapping:", table, err)
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error creating user for login: %s", err),
				}
			} else {
				params["data"] = _data_user
				return app.confirm_emmail(params)
			}
		}
	}
	// CRUD ACTIONS
	get_crud_actions_sql := fmt.Sprintf(`SELECT * FROM "crud_action" WHERE "active" = TRUE AND "db" = ? AND "table" = ? AND "%" IS TRUE`, crud_aciton)
	crud_action_rows, err := app.AdminGetRowsByFilter(get_crud_actions_sql, validation_data)
	if err != nil {
		fmt.Printf("Error occurred while fetching crud_actions: %v", err)
		/*return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Error occurred while fetching crud_actions: %v", err),
		}*/
	} else if len(crud_action_rows) > 0 {
		actionRunner := func(c_action) error {			
			action_type_id := c_action["action_type_id"]
			_, okSql := c_action["sql"]
			_, okEmail := c_action["email_template"]
			api, okAPI := c_action["api"]
			//  register crud_action_logs
			success := true
			crud_action_log := Dict{
				"crud_action_id":   c_action["crud_action_id"],
				"crud_action_code": c_action["crud_action_code"],
				"crud_action":      c_action["crud_action"],
				"table":            table,
				"db":               database,
				"id":               id,
				"action":           crud_aciton,
				"user_id":          user_id,
				"app_id":           params["app"].(Dict)["app_id"],
			}
			insert_crud_action_log_sql := `INSERT INTO "crud_action_logs" ("crud_action_id", "crud_action_code", "crud_action", "table", "db", "id", "action", "user_id", "app_id", "executed_at", "created_at", "updated_at") 
			VALUES (:crud_action_id, :crud_action_code, :crud_action, :table, :db, :id, :action, :user_id, :app_id, :executed_at, :created_at, :updated_at)`
			msg := ""
			if app.toInt(action_type_id) == 1 && okSql { // ExecuteQuery
				if _, ok := c_action["sql"]; ok {
					sql_rule := c_action["sql"].(string)
					err := app.AdminExecuteQuery(sql_rule, _data)
					if err != nil {
						success = false
						msg, _ = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
						crud_action_log["success"] = success
						crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action ExecuteQuery: %v. Message: %s", err, msg)
						crud_action_log["executed_at"] = time.Now()
						crud_action_log["created_at"] = time.Now()
						crud_action_log["updated_at"] = time.Now()
						_, err2 = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
						if err2 != nil {
							fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
						}
						return fmt.Error(msg)
					}
				}
			} else if app.toInt(action_type_id) == 2 && okEmail { // SendEmail
				if _, ok := c_action["email_template"]; ok {
					// Process email template
					_to, _ := c_action["email_to"].(string)
					to := strings.Split(_to, ";")
					emailParams := Dict{
						"to":      to,
						"subject": fmt.Sprintf("%s - %s", c_action["crud_action_code"], c_action["crud_action"]),
						"body":    c_action["email_template"],
						"data": Dict{
							"data": _data,
							"user": params["user"],
						},
					}
					err = etlx_engine.SendEmail(emailParams)
					if err != nil {
						success = false
						crud_action_log["success"] = success
						crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action SendEmail: %v. Message: %s", err, err.Error())
						crud_action_log["executed_at"] = time.Now()
						crud_action_log["created_at"] = time.Now()
						crud_action_log["updated_at"] = time.Now()
						_, err2 = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
						if err2 != nil {
							fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
						}
						return err
					}
				}
			} else if app.toInt(action_type_id) == 3 && okAPI { // CallAPI
				_, err := app.CronRunEndPoint(Dict{"api": api})
				if err != nil {
					success = false
					crud_action_log["success"] = success
					crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action API: %v. Message: %s", err, err.Error())
					crud_action_log["executed_at"] = time.Now()
					crud_action_log["created_at"] = time.Now()
					crud_action_log["updated_at"] = time.Now()
					_, err2 = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
					if err2 != nil {
						fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
					}
					return err
				}
			} else {
				fmt.Println("Unknown action_type_id for crud_action:", action_type_id)
				crud_action_log["success"] = success
				crud_action_log["log_message"] = fmt.Sprintf("Unknown action_type_id: %v", action_type_id)
				crud_action_log["executed_at"] = time.Now()
				crud_action_log["created_at"] = time.Now()
				crud_action_log["updated_at"] = time.Now()
				_, err = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
				if err != nil {
					fmt.Printf("Error inserting crud action log for unknown action_type_id for crud_action_id %v: %v", c_action["crud_action_id"], err)
				}
			}
			crud_action_log["success"] = success
			if msg == "" {
				msg = fmt.Sprintf("CRUD Action %s executed successfully", c_action["crud_action_code"])
			}
			crud_action_log["log_message"] = msg
			crud_action_log["executed_at"] = time.Now()
			crud_action_log["created_at"] = time.Now()
			crud_action_log["updated_at"] = time.Now()
			_, err = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err)
				// Not returning error to avoid interrupting the main CRUD operation flow
			}
			return nil
		}
		for _, c_action := range crud_action_rows {
			parallel, ok := c_action["parallel"].(bool)
			if parallel && ok {
				go func() {
					err := actionRunner(c_action)
					if err != nil {
						fmt.Printf("Error runing the action: %v\n", err)
					}
				}()
			} else {
				err := actionRunner(c_action)
			}
		}
	}
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
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
