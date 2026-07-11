package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"
)

func (app *application) RunCrudAction(params, c_action, _data Dict) error {
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
	action_type_id := c_action["action_type_id"]
	_, okSql := c_action["sql"]
	_, okEmail := c_action["email_template"]
	api, okAPI := c_action["api"]
	api_id, okAPIID := c_action["api_id"]
	api_name, okAPIName := c_action["api_name"]
	api_endpoint, okAPIEndpoint := c_action["api_endpoint"]
	pdf_path, okPDFPath := c_action["pdf_path"]
	pdf_template, okPDFTmpl := c_action["pdf_template"]
	pdf_tex_template, okPDFTexTmpl := c_action["pdf_tex_template"]
	after_sql, okAfterSQL := c_action["after_sql"].(string)
	etlx_md_template, okETLX := c_action["etlx_md_template"]
	//  register crud_action_logs
	success := true
	etlx_engine := &etlx.ETLX{}
	crud_action_log := Dict{
		"crud_action_id":   c_action["crud_action_id"],
		"crud_action_code": c_action["crud_action_code"],
		"crud_action":      c_action["crud_action"],
		"table":            table,
		"db":               database,
		"pk_field":         pk,
		"id":               id,
		"action":           crud_aciton,
		"action_type":      c_action["action_type_id"],
		"user_id":          user_id,
		"app_id":           params["app"].(Dict)["app_id"],
		"started_at":       time.Now().In(loc),
		"log_data":         "",
	}
	insert_crud_action_log_sql := `INSERT INTO "crud_action_logs" ("crud_action_id", "crud_action_code", "crud_action", "table", "db", "pk_field", "id", "action", "action_type", "success", "log_message", "log_data", "user_id", "app_id", "executed_at", "created_at", "updated_at") 
			VALUES (:crud_action_id, :crud_action_code, :crud_action, :table, :db, :pk_field, :id, :action, :action_type, :success, :log_message, :log_data, :user_id, :app_id, :executed_at, :created_at, :updated_at)`
	// ACTION DATA TO HELP BUILD THE TEPLATE
	sql := "select * from action_data where crud_action_id = ? and excluded = false"
	valid_data_res, err := app.AdminGetRowsByFilter(sql, []any{c_action["crud_action_id"]})
	if err != nil {
		return fmt.Errorf("Error getting ACTION Data: %s", err)
	}
	// fmt.Println("ACTION DATA:", valid_data_res)
	action_data, err := app.GetActionData(params, valid_data_res, _data)
	if err != nil {
		return fmt.Errorf("Error getting ACTION Data: %s", err)
	}
	for key, val := range action_data {
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
			return fmt.Errorf("Error preparing sql_condition %s -> %s: %v", c_action["crud_action_code"], sql_condition, err.Error())
		}
		res, _, err := app.db.QuerySingleRow(query, args...)
		if err != nil {
			return fmt.Errorf("Error executing sql_condition %s -> %s: %v", c_action["crud_action_code"], sql_condition, err.Error())
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
			msg = fmt.Sprintf("SQL Condition not met for CRUD Action: %v", c_action["crud_action_code"])
			crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Sprintf("SQL Condition not met for CRUD Action: %v", c_action["crud_action_code"])
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err)
			}
			return fmt.Errorf(msg)
		}
	}
	// fmt.Println("PASSED COND!")
	if app.toInt(action_type_id) == 1 && okSql { // ExecuteQuery
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
				_data["err"] = err.Error()
				success = false
				msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
				if err != nil {
					msg = c_action["err_msg"].(string)
				}
				crud_action_log["success"] = success
				crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action ExecuteQuery: %v. Message: %s", _data["err"], msg)
				crud_action_log["executed_at"] = time.Now().In(loc)
				crud_action_log["created_at"] = time.Now().In(loc)
				crud_action_log["updated_at"] = time.Now().In(loc)
				_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
				if err2 != nil {
					fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
				}
				msg = fmt.Sprintf("Error executing query: %s", err.Error())
				return fmt.Errorf("Err %s", msg)
			}
		}
	} else if app.toInt(action_type_id) == 2 && okEmail { // SendEmail
		if _, ok := c_action["email_template"]; ok {
			// Process email template
			_to, _ := c_action["email_to"].(string)
			_to_tmpl, _ := etlx_engine.RenderTemplate(_to, _data)
			_subject, ok := c_action["email_subject"].(string)
			if !ok {
				_subject = fmt.Sprintf("%s - %s", c_action["crud_action_code"], c_action["crud_action"])
			}
			subject, _ := etlx_engine.RenderTemplate(_subject, _data)
			//fmt.Println("Rendered email_to template:", _to, _to_tmpl, subject)
			to := strings.Split(_to_tmpl, ";")
			emailParams := Dict{
				"to":      app.sliceStrs2SliceInterfaces(to),
				"subject": subject,
				"body":    c_action["email_template"],
				"data": Dict{
					"data": _data,
					"user": params["user"],
				},
			}
			err := etlx_engine.SendEmail(emailParams)
			if err != nil {
				_data["err"] = err.Error()
				success = false
				crud_action_log["success"] = success
				crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action SendEmail: %v. Message: %s", err, err.Error())
				crud_action_log["executed_at"] = time.Now().In(loc)
				crud_action_log["created_at"] = time.Now().In(loc)
				crud_action_log["updated_at"] = time.Now().In(loc)
				_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
				if err2 != nil {
					fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
				}
				msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
				if err != nil {
					msg = c_action["err_msg"].(string)
				}
				return fmt.Errorf("Err %w", msg)
			}
		}
	} else if app.toInt(action_type_id) == 3 && okAPI { // CallAPI
		_api, err := etlx_engine.RenderTemplate(api.(string), _data)
		if err == nil {
			api = _api
		}
		_, err = app.CronRunEndPoint(Dict{"api": api, "data": _data})
		if err != nil {
			_data["err"] = err.Error()
			success = false
			crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action API: %s", _data["err"])
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err2 != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
			}
			msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
			if err != nil {
				msg = c_action["err_msg"].(string)
			}
			return err
		}
	} else if app.toInt(action_type_id) == 4 && (okAPIID || okAPIName || okAPIEndpoint) { // CallExternalAPI
		_params := Dict{
			"user": params["user"],
			"app":  params["app"],
			"lang": params["lang"],
			"data": Dict{
				"api_id":         api_id,
				"api_name":       api_name,
				"api_endpoint":   api_endpoint,
				"data":           _data,
				"action":         c_action,
				"db":             database,
				"table":          table,
				"table_pk_field": pk,
				"row_id":         id,
				"user":           params["user"],
			},
		}
		res := app.runAPI(_params)
		_, ok := res["success"].(bool)
		//fmt.Println("API RES:", res["data"].(string))
		if !ok || !res["success"].(bool) {
			success = false
			msg := ""
			if _, ok := res["msg"].(string); ok {
				msg = res["msg"].(string)
			}
			crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Sprintf("Error executing CRUD Action External API. Message: %s. API Response: %v", msg, res)
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err2 != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
			}
			return fmt.Errorf("Error executing external API: %s", msg)
		} else if log_data, ok := res["data"].(string); ok {
			crud_action_log["log_data"] = log_data
			_data["api_response"] = log_data
		}
	} else if app.toInt(action_type_id) == 5 && (okPDFPath && (okPDFTmpl || okPDFTexTmpl)) { // GeneratePDF
		// fmt.Println("GeneratePDF:", pdf_path, pdf_template)
		output_path, err := etlx_engine.RenderTemplate(pdf_path.(string), _data)
		if err != nil {
			output_path = pdf_path.(string)
		}
		use_latex := app.toBool(c_action["use_latex"])
		output_path = etlx_engine.ReplaceEnvVariable(output_path)
		_data["fname"] = output_path
		if use_latex {
			latex := pdf_tex_template.(string)
			//latex = LatexEscape(latex)
			latex, err = etlx_engine.RenderTemplate(latex, _data)
			latex = etlx_engine.ReplaceEnvVariable(latex)
			if err != nil {
				return err
			}
			err = app.GenPDFFromLatex(latex, output_path)
		} else {
			html, err := etlx_engine.RenderTemplate(pdf_template.(string), _data)
			html = etlx_engine.ReplaceEnvVariable(html)
			if err != nil {
				return err
			}
			err = app.GenPDFFromHTML(html, output_path)
		}
		if err != nil {
			_data["err"] = err.Error()
			fmt.Println("GeneratePDF:", err)
			success = false
			msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
			if err != nil {
				msg = c_action["err_msg"].(string)
			}
			crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Sprintf("Error generating PDF: %v", err)
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err2 != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
			}
			return fmt.Errorf("Error generating PDF: %s", msg)
		}
	} else if app.toInt(action_type_id) == 6 && okETLX { // ETLX
		etlx_md_template, err := app.RenderTextTemplate(etlx_md_template.(string), _data)
		if err != nil {
			return err
		}
		//fmt.Println(etlx_md_template)
		params := Dict{
			"db": app.config.db.dsn,
			"data": Dict{
				"order_metadata": any(true),
				"params":         _data,
				"conf":           any(string(etlx_md_template)),
			},
		}
		res := app.etlxRun(params, true)
		if res["success"].(bool) != true {
			success = false
			_data["error"] = res["msg"]
			msg, err = etlx_engine.RenderTemplate(c_action["err_msg"].(string), _data)
			if err != nil {
				msg = c_action["err_msg"].(string)
			}
			crud_action_log["success"] = success
			crud_action_log["log_message"] = msg
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err2 := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err2 != nil {
				fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err2)
			}
			return fmt.Errorf("Error executing external ETLX: %s", msg)
		}
	} else {
		success = false
		fmt.Println("Unknown action_type_id for crud_action:", action_type_id)
		msg = fmt.Sprintf("Unknown action type for CRUD Action: %v", action_type_id)
		crud_action_log["success"] = success
		crud_action_log["log_message"] = fmt.Sprintf("Unknown action_type_id: %v", action_type_id)
		crud_action_log["executed_at"] = time.Now().In(loc)
		crud_action_log["created_at"] = time.Now().In(loc)
		crud_action_log["updated_at"] = time.Now().In(loc)
		_, err := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
		if err != nil {
			fmt.Printf("Error inserting crud action log for unknown action_type_id for crud_action_id %v: %v", c_action["crud_action_id"], err)
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
			crud_action_log["success"] = success
			crud_action_log["log_message"] = fmt.Errorf("Error executing after SQL %s!", err.Error())
			crud_action_log["executed_at"] = time.Now().In(loc)
			crud_action_log["created_at"] = time.Now().In(loc)
			crud_action_log["updated_at"] = time.Now().In(loc)
			_, err := app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
			if err != nil {
				fmt.Printf("Error inserting crud action log for unknown action_type_id for crud_action_id %v: %v", c_action["crud_action_id"], err)
			}
			return fmt.Errorf("Error executing after SQL %s!", err.Error())
		}
	}
	crud_action_log["success"] = success
	if msg == "" || success {
		msg = fmt.Sprintf("CRUD Action %s executed successfully", c_action["crud_action_code"])
	}
	crud_action_log["log_message"] = msg
	crud_action_log["executed_at"] = time.Now().In(loc)
	crud_action_log["created_at"] = time.Now().In(loc)
	crud_action_log["updated_at"] = time.Now().In(loc)
	_, err = app.db.ExecuteNamedQuery(insert_crud_action_log_sql, crud_action_log)
	if err != nil {
		fmt.Printf("Error inserting crud action log for crud_action_id %v: %v", c_action["crud_action_id"], err)
		// Not returning error to avoid interrupting the main CRUD operation flow
	}
	if !success {
		return fmt.Errorf(msg)
	} else {
		// CRUD ACTIONS
		get_crud_actions_sql := `SELECT ca.*, at.action_trigger_action_id, ca2.crud_action_code AS main_action
FROM crud_action ca
JOIN action_trigger_action at ON at.action_trigger_code = ca.crud_action_code
JOIN crud_action ca2 ON ca2.crud_action_id = at.crud_action_id
WHERE at.crud_action_id = ? AND at.excluded = false AND ca.excluded = false
-- ORDER BY at."trigger_order" ASC`
		crud_action_rows, err := app.AdminGetRowsByFilter(get_crud_actions_sql, []any{c_action["crud_action_id"]})
		if err != nil {
			fmt.Printf("Error occurred while fetching crud_actions: %v", err)
		} else if len(crud_action_rows) > 0 {
			fmt.Println("HAS TRIGGERED ACTIONS!")
			for _, c_ation_trigger := range crud_action_rows {
				fmt.Println("HAS TRIGGERED ACTION:", c_ation_trigger["crud_action_code"])
				parallel, ok := c_ation_trigger["parallel"].(bool)
				if parallel && ok {
					go func() {
						err := app.RunCrudAction(params, c_ation_trigger, _data) //actionRunner(c_ation_trigger)
						if err != nil {
							fmt.Printf("Error runing the action: %s triggered by %s -> %v\n", c_ation_trigger["crud_action_code"], c_ation_trigger["main_action"], err.Error())
						}
					}()
				} else {
					err := app.RunCrudAction(params, c_ation_trigger, _data) // actionRunner(c_ation_trigger)
					if err != nil {
						fmt.Printf("Error runing the action: %s triggered by %s -> %v\n", c_ation_trigger["crud_action_code"], c_ation_trigger["main_action"], err.Error())
					}
				}
			}
		}
	}
	return nil
}

func (app *application) GetActionData(params Dict, valid_data_res []Dict, _data Dict) (Dict, error) {
	// loop valid_data_res
	etlx_engine := &etlx.ETLX{}
	res := Dict{}
	for _, action_data := range valid_data_res {
		name := action_data["action_data"].(string)
		if app.toInt(action_data["action_data_type_id"]) == 3 { // ODATA
			odata_path, ok := action_data["odata_path"].(string)
			if !ok {
				return nil, fmt.Errorf("Error, odata_path is not set!")
			}
			odata_path, err := etlx_engine.RenderTemplate(odata_path, _data)
			if err != nil {
				odata_path = action_data["odata_path"].(string)
			} else {
				odata_path = etlx_engine.ReplaceEnvVariable(odata_path)
			}
			// fmt.Println(action_data["odata_path"], odata_path)
			sigle_row_obj := app.toBool(action_data["sigle_row_obj"])
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
			return nil, fmt.Errorf("Error: API Data Type %d is not implemented yet!", app.toInt(action_data["action_data_type_id"]))
		}
	}
	return res, nil
}

func (app *application) GetValidationData(params Dict, valid_data_res []Dict, _data Dict) (Dict, error) {
	// loop valid_data_res
	etlx_engine := &etlx.ETLX{}
	res := Dict{}
	for _, validation_data := range valid_data_res {
		name := validation_data["validation_data"].(string)
		odata_path, ok := validation_data["odata_path"].(string)
		if !ok {
			return nil, fmt.Errorf("Error, odata_path is not set!")
		}
		odata_path, err := etlx_engine.RenderTemplate(odata_path, _data)
		if err != nil {
			odata_path = validation_data["odata_path"].(string)
		} else {
			odata_path = etlx_engine.ReplaceEnvVariable(odata_path)
		}
		// fmt.Println(validation_data["odata_path"], odata_path)
		sigle_row_obj := app.toBool(validation_data["sigle_row_obj"])
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
	}
	return res, nil
}
