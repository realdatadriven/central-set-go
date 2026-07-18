package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"
)

func (app *application) RunCrudIntercept(params, c_action, _data Dict) error {
	var loc *time.Location
	if _, ok := params["location"].(*time.Location); ok {
		loc = params["location"].(*time.Location)
	} else {
		loc = time.Local
	}
	table := params["table"]
	database := params["database"]
	id := params["id"]
	crud_aciton := params["crud_aciton"]
	user_id := params["user_id"]
	pk := params["pk"]
	intercept_type_id := c_action["intercept_type_id"]
	_, okSql := c_action["sql"]
	_, okFieldList := c_action["field_list"]
	after_sql, okAfterSQL := c_action["after_sql"].(string)
	//  register crud_intercept_logs
	success := true
	etlx_engine := &etlx.ETLX{}
	crud_intercept_log := Dict{
		"crud_intercept_id":   c_action["crud_intercept_id"],
		"crud_intercept_code": c_action["crud_intercept_code"],
		"crud_action":         c_action["crud_action"],
		"table":               table,
		"db":                  database,
		"pk_field":            pk,
		"id":                  id,
		"action":              crud_aciton,
		"intercept_type":      c_action["intercept_type_id"],
		"user_id":             user_id,
		"app_id":              params["app"].(Dict)["app_id"],
		"started_at":          time.Now().In(loc),
		"log_data":            "",
	}
	insert_crud_intercept_log_sql := `INSERT INTO "crud_intercept_logs" ("crud_intercept_id", "crud_intercept_code", "crud_action", "table", "db", "pk_field", "id", "action", "intercept_type", "success", "log_message", "log_data", "user_id", "app_id", "executed_at", "created_at", "updated_at") 
			VALUES (:crud_intercept_id, :crud_intercept_code, :crud_action, :table, :db, :pk_field, :id, :action, :intercept_type, :success, :log_message, :log_data, :user_id, :app_id, :executed_at, :created_at, :updated_at)`
	// ACTION DATA TO HELP BUILD THE TEPLATE
	sql := "select * from intercept_data where crud_intercept_id = ? and excluded = false"
	valid_data_res, err := app.AdminGetRowsByFilter(sql, []any{c_action["crud_intercept_id"]})
	if err != nil {
		return fmt.Errorf("Error getting ACTION Data: %s", err)
	}
	// fmt.Println("ACTION DATA:", valid_data_res)
	intercept_data, err := app.GetInterceptData(params, valid_data_res, _data)
	if err != nil {
		return fmt.Errorf("Error getting ACTION Data: %s", err)
	}
	for key, val := range intercept_data {
		_data[key] = val
	}
	msg := ""
	// CHECK CONFITION SQL CONDITION
	sql_condition, sqlConditionOk := c_action["sql_condition"].(string)
	if sqlConditionOk && sql_condition != "" {
		sql_condition, err = etlx_engine.RenderTemplate(sql_condition, _data)
		if err != nil {
			return fmt.Errorf("Error rendering sql_condition %v", err.Error())
		}
		sql_condition = etlx_engine.ReplaceEnvVariable(sql_condition)
		// fmt.Println("SQL COND:", sql_condition)
		query, args, err := etlx_engine.NamedToPositional(sql_condition, _data)
		if err != nil {
			return fmt.Errorf("Error preparing sql_condition %s -> %s: %v", c_action["crud_intercept_code"], sql_condition, err.Error())
		}
		res, _, err := app.db.QuerySingleRow(query, args...)
		if err != nil {
			return fmt.Errorf("Error executing sql_condition %s -> %s: %v", c_action["crud_intercept_code"], sql_condition, err.Error())
		}
		cond := false
		if len(*res) > 0 {
			if _, ok := (*res)["cond"]; !ok {
				return fmt.Errorf("Error executing sql_condition, 'cond' column not found in result set")
			} else {
				cond = app.toBool((*res)["cond"])
			}
		}
		if !cond {
			success = false
			msg = fmt.Sprintf("SQL Condition not met for CRUD Action: %v", c_action["crud_intercept_code"])
			crud_intercept_log["success"] = success
			crud_intercept_log["log_message"] = fmt.Sprintf("SQL Condition not met for CRUD Action: %v", c_action["crud_intercept_code"])
			crud_intercept_log["executed_at"] = time.Now().In(loc)
			crud_intercept_log["created_at"] = time.Now().In(loc)
			crud_intercept_log["updated_at"] = time.Now().In(loc)
			_, err := app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
			if err != nil {
				fmt.Printf("Error inserting crud action log for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err)
			}
			return fmt.Errorf(msg)
		}
	}
	// fmt.Println("PASSED COND!")
	if app.toInt(intercept_type_id) == 1 && okSql { // EncapReadQueryBeforeExec
		if _, ok := c_action["sql"]; ok {
			// fmt.Println(c_action["sql"])
			sql_rule, err := etlx_engine.RenderTemplate(c_action["sql"].(string), _data)
			if err != nil {
				return fmt.Errorf("Error rendering sql_condition %v", err.Error())
			}
			sql_rule = etlx_engine.ReplaceEnvVariable(sql_rule)
			// fmt.Println(sql_rule)
			err = app.ExecuteQuery(sql_rule, params, _data)
			if err != nil {
				fmt.Println("SQL ERR:", sql_rule, err.Error())
				_data["err"] = err.Error()
				success = false
				msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
				if err != nil {
					msg = c_action["err_msg"].(string)
				}
				crud_intercept_log["success"] = success
				crud_intercept_log["log_message"] = fmt.Sprintf("Error executing CRUD Action ExecuteQuery: %v. Message: %s", _data["err"], msg)
				crud_intercept_log["executed_at"] = time.Now().In(loc)
				crud_intercept_log["created_at"] = time.Now().In(loc)
				crud_intercept_log["updated_at"] = time.Now().In(loc)
				_, err2 := app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
				if err2 != nil {
					fmt.Printf("Error inserting crud action log for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err2)
				}
				msg = fmt.Sprintf("Error executing query: %s", err.Error())
				return fmt.Errorf("Err %s", msg)
			}
		}
	} else if app.toInt(intercept_type_id) == 2 && okFieldList { // RemoveField
		if _, ok := c_action["field_list"]; ok {
			// Process email template
			_to, _ := c_action["email_to"].(string)
			_to_tmpl, _ := etlx_engine.RenderTemplate(_to, _data)
			_subject, ok := c_action["email_subject"].(string)
			if !ok {
				_subject = fmt.Sprintf("%s - %s", c_action["crud_intercept_code"], c_action["crud_action"])
			}
			subject, _ := etlx_engine.RenderTemplate(_subject, _data)
			//fmt.Println("Rendered email_to template:", _to, _to_tmpl, subject)
			to := strings.Split(_to_tmpl, ";")
			emailParams := Dict{
				"to":      app.sliceStrs2SliceInterfaces(to),
				"subject": subject,
				"body":    c_action["field_list"],
				"data": Dict{
					"data": _data,
					"user": params["user"],
				},
			}
			err := etlx_engine.SendEmail(emailParams)
			if err != nil {
				_data["err"] = err.Error()
				success = false
				crud_intercept_log["success"] = success
				crud_intercept_log["log_message"] = fmt.Sprintf("Error executing CRUD Action SendEmail: %v. Message: %s", err, err.Error())
				crud_intercept_log["executed_at"] = time.Now().In(loc)
				crud_intercept_log["created_at"] = time.Now().In(loc)
				crud_intercept_log["updated_at"] = time.Now().In(loc)
				_, err2 := app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
				if err2 != nil {
					fmt.Printf("Error inserting crud action log for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err2)
				}
				msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
				if err != nil {
					msg = c_action["err_msg"].(string)
				}
				return fmt.Errorf("Err %w", msg)
			}
		}
	} else {
		success = false
		fmt.Println("Unknown intercept_type_id for crud_action:", intercept_type_id)
		msg = fmt.Sprintf("Unknown action type for CRUD Action: %v", intercept_type_id)
		crud_intercept_log["success"] = success
		crud_intercept_log["log_message"] = fmt.Sprintf("Unknown intercept_type_id: %v", intercept_type_id)
		crud_intercept_log["executed_at"] = time.Now().In(loc)
		crud_intercept_log["created_at"] = time.Now().In(loc)
		crud_intercept_log["updated_at"] = time.Now().In(loc)
		_, err := app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
		if err != nil {
			fmt.Printf("Error inserting crud action log for unknown intercept_type_id for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err)
		}
		return fmt.Errorf(msg)
	}
	if okAfterSQL && after_sql != "" {
		after_sql, err = etlx_engine.RenderTemplate(after_sql, _data)
		if err != nil {
			return fmt.Errorf("Error rendering after sql %v", err.Error())
		}
		after_sql = etlx_engine.ReplaceEnvVariable(after_sql)
		err = app.ExecuteQuery(after_sql, params, _data)
		if err != nil {
			success = false
			crud_intercept_log["success"] = success
			crud_intercept_log["log_message"] = fmt.Errorf("Error executing after SQL %s!", err.Error())
			crud_intercept_log["executed_at"] = time.Now().In(loc)
			crud_intercept_log["created_at"] = time.Now().In(loc)
			crud_intercept_log["updated_at"] = time.Now().In(loc)
			_, err := app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
			if err != nil {
				fmt.Printf("Error inserting crud action log for unknown intercept_type_id for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err)
			}
			return fmt.Errorf("Error executing after SQL %s!", err.Error())
		}
	}
	crud_intercept_log["success"] = success
	if msg == "" || success {
		msg = fmt.Sprintf("CRUD Action %s executed successfully", c_action["crud_intercept_code"])
	}
	crud_intercept_log["log_message"] = msg
	crud_intercept_log["executed_at"] = time.Now().In(loc)
	crud_intercept_log["created_at"] = time.Now().In(loc)
	crud_intercept_log["updated_at"] = time.Now().In(loc)
	_, err = app.db.ExecuteNamedQuery(insert_crud_intercept_log_sql, crud_intercept_log)
	if err != nil {
		fmt.Printf("Error inserting crud action log for crud_intercept_id %v: %v", c_action["crud_intercept_id"], err)
		// Not returning error to avoid interrupting the main CRUD operation flow
	}
	if !success {
		return fmt.Errorf(msg)
	} else {
		// CRUD ACTIONS
		get_crud_actions_sql := `SELECT ca.*, at.intercept_trigger_intercept_id, ca2.crud_intercept_code AS main_action
FROM crud_action ca
JOIN intercept_trigger_action at ON at.intercept_trigger_code = ca.crud_intercept_code
JOIN crud_action ca2 ON ca2.crud_intercept_id = at.crud_intercept_id
WHERE at.crud_intercept_id = ? AND at.excluded = false AND ca.excluded = false
-- ORDER BY at."trigger_order" ASC`
		crud_intercept_rows, err := app.AdminGetRowsByFilter(get_crud_actions_sql, []any{c_action["crud_intercept_id"]})
		if err != nil {
			fmt.Printf("Error occurred while fetching crud_actions: %v", err)
		} else if len(crud_intercept_rows) > 0 {
			fmt.Println("HAS TRIGGERED ACTIONS!")
			for _, c_ation_trigger := range crud_intercept_rows {
				fmt.Println("HAS TRIGGERED ACTION:", c_ation_trigger["crud_intercept_code"])
				parallel, ok := c_ation_trigger["parallel"].(bool)
				if parallel && ok {
					go func() {
						err := app.RunCrudIntercept(params, c_ation_trigger, _data) //actionRunner(c_ation_trigger)
						if err != nil {
							fmt.Printf("Error runing the action: %s triggered by %s -> %v\n", c_ation_trigger["crud_intercept_code"], c_ation_trigger["main_action"], err.Error())
						}
					}()
				} else {
					err := app.RunCrudIntercept(params, c_ation_trigger, _data) // actionRunner(c_ation_trigger)
					if err != nil {
						fmt.Printf("Error runing the action: %s triggered by %s -> %v\n", c_ation_trigger["crud_intercept_code"], c_ation_trigger["main_action"], err.Error())
					}
				}
			}
		}
	}
	return nil
}

func (app *application) GetInterceptData(params Dict, valid_data_res []Dict, _data Dict) (Dict, error) {
	// loop valid_data_res
	etlx_engine := &etlx.ETLX{}
	res := Dict{}
	for _, intercept_data := range valid_data_res {
		name := intercept_data["intercept_data"].(string)
		if app.toInt(intercept_data["intercept_data_type_id"]) == 3 { // ODATA
			odata_path, ok := intercept_data["odata_path"].(string)
			if !ok {
				return nil, fmt.Errorf("Error, odata_path is not set!")
			}
			odata_path, err := etlx_engine.RenderTemplate(odata_path, _data)
			if err != nil {
				odata_path = intercept_data["odata_path"].(string)
			} else {
				odata_path = etlx_engine.ReplaceEnvVariable(odata_path)
			}
			// fmt.Println(intercept_data["odata_path"], odata_path)
			sigle_row_obj := app.toBool(intercept_data["sigle_row_obj"])
			// fmt.Println(db, table, query)
			if err != nil {
				return nil, err
			}
			results, err := app.OData2C7Read(params, odata_path)
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
			return nil, fmt.Errorf("Error: Intercept Data Type %d is not implemented yet!", app.toInt(intercept_data["intercept_data_type_id"]))
		}
	}
	return res, nil
}
