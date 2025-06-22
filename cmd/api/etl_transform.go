package main

import (
	"fmt"
	"path/filepath"
	"time"

	querybuilder "github.com/realdatadriven/central-set-go/querydoc"

	"github.com/realdatadriven/etlx"
	// "github.com/google/uuid"
)

func (app *application) _transform(params map[string]interface{}, _item map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, db_conf map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
	// IN MEMORY DUCKDB CONN
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("DDB Conn: %s", err),
		}
	}
	defer db.Close()
	// DESTNATION DB
	_driver := ""
	if _, ok := db_conf["driverName"]; ok {
		_driver = db_conf["driverName"].(string)
	}
	_database := ""
	if _, ok := _item["database"]; ok {
		_database = _item["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]interface{}{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := _etlrb["database"]; ok {
		_database = _etlrb["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]interface{}{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := db_conf["dsn"]; ok {
		_database = db_conf["dsn"].(string)
	}
	// ATTACH DBS TO THE DUCKDB IN MEM CONN
	_duck_conf := map[string]interface{}{}
	if _, ok := _conf["duckdb"]; ok {
		_duck_conf = _conf["duckdb"].(map[string]interface{})
	}
	if _, ok := _duck_conf["extensions"]; !ok {
		_duck_conf["extensions"] = []interface{}{}
	}
	//fmt.Println("extensions:", _duck_conf["extensions"])
	app.duckdb_start(db, _duck_conf, _driver, _database)
	/*/ RUN EXTRACT HERE
	_db_base := filepath.Base(_database)
	_db_ext := filepath.Ext(_database)
	_db_name_no_ext := _db_base[:len(_db_base)-len(_db_ext)]
	//fmt.Println(_driver, _database, _db_base, _db_ext, _db_name_no_ext)/*/
	destination_table := ""
	if _, ok := _item["destination_table"]; ok {
		destination_table = _item["destination_table"].(string)
	}
	// DATE REF
	var _date_ref interface{}
	if _, ok := _item["date_ref"]; ok {
		_date_ref = _item["date_ref"]
	} else if _, ok := _step["dates_refs"]; ok {
		_date_ref = _step["dates_refs"]
	}
	var date_ref []time.Time
	switch _type := _date_ref.(type) {
	case string:
		_dt, _ := time.Parse("2006-01-02", _date_ref.(string))
		date_ref = append(date_ref, _dt)
	case []interface{}:
		for _, _dt := range _date_ref.([]interface{}) {
			_dt, _ := time.Parse("2006-01-02", _dt.(string))
			date_ref = append(date_ref, _dt)
		}
	default:
		fmt.Println("default:", _type)
	}
	// DATE FIELDS
	ref_date_field := ""
	if _, ok := _item["ref_date_field"]; !ok {
	} else if _, ok := _item["ref_date_field"].(string); ok {
		ref_date_field = _item["ref_date_field"].(string)
	}
	if _, ok := _item["date_field"].(string); ok {
		ref_date_field = _item["date_field"].(string)
	}
	date_format_org := "YYYYMMDD"
	if _, ok := _item["date_format_org"]; !ok {
	} else if _, ok := _item["date_format_org"].(string); ok {
		date_format_org = _item["date_format_org"].(string)
	}
	if _, ok := _item["date_field_format"]; !ok {
	} else if _, ok := _item["date_field_format"].(string); ok {
		date_format_org = _item["date_field_format"].(string)
	}
	// GET THE OUTPUT FIELDS
	_fields_table := ""
	if _, ok := _step["detail_table"]; ok {
		_fields_table = _step["detail_table"].(string)
	} else {
		_fields_table = "etl_rb_output_field"
	}
	etl_rbase_output_id := interface{}(-1)
	if _, ok := _item["etl_rbase_output_id"]; ok {
		etl_rbase_output_id = _item["etl_rbase_output_id"]
	}
	params["data"].(map[string]interface{})["table"] = _fields_table
	params["data"].(map[string]interface{})["limit"] = interface{}(-1.0)
	params["data"].(map[string]interface{})["offset"] = interface{}(0.0)
	params["data"].(map[string]interface{})["filters"] = []interface{}{}
	params["data"].(map[string]interface{})["filters"] = append(
		params["data"].(map[string]interface{})["filters"].([]interface{}),
		map[string]interface{}{
			"field": "etl_rbase_output_id",
			"cond":  "=",
			"value": etl_rbase_output_id,
		},
	)
	params["data"].(map[string]interface{})["order_by"] = []interface{}{}
	params["data"].(map[string]interface{})["order_by"] = append(
		params["data"].(map[string]interface{})["order_by"].([]interface{}),
		map[string]interface{}{
			"field": "field_order",
			"order": "ASC",
		},
	)
	params["data"].(map[string]interface{})["order_by"] = append(
		params["data"].(map[string]interface{})["order_by"].([]interface{}),
		map[string]interface{}{
			"field": "etl_rb_output_field_id",
			"order": "ASC",
		},
	)
	query_parts := map[string]interface{}{}
	_fields_order := []string{}
	_aux_fields := app.read(params)
	if _, ok := _aux_fields["success"]; !ok {
		return _aux_fields
	} else if _aux_fields["success"].(bool) {
		//fmt.Println(_aux_fields["sql"])
		for _, _field := range _aux_fields["data"].([]map[string]interface{}) {
			_fields_order = append(_fields_order, _field["etl_rb_output_field"].(string))
			query_parts[_field["etl_rb_output_field"].(string)] = map[string]interface{}{
				"name":     _field["etl_rb_output_field"],
				"desc":     _field["etl_rb_output_field_desc"],
				"select":   _field["sql_select"],
				"from":     _field["sql_from"],
				"join":     _field["sql_join"],
				"where":    _field["sql_where"],
				"group_by": _field["sql_group_by"],
				"order_by": _field["sql_order_by"],
				"having":   _field["sql_having"],
				"window":   _field["sql_window"],
				"active":   _field["active"],
			}
		}
	}
	//fmt.Println("AUX FIELDS:", len(query_parts))
	qd := querybuilder.QueryDoc{
		QueryParts:  make(map[string]querybuilder.Field),
		FieldOrders: _fields_order,
	}
	err = qd.SetQueryPartsFromMap(query_parts)
	if err != nil {
		fmt.Println("Error setting field:", err)
	} else {
		fmt.Println("Field set successfully:")
	}
	_sql := qd.GetQuerySQLFromMap()
	_sql = app.setQueryDate(_sql, date_ref)
	//fmt.Println(_sql[:100])
	//sql_drop := ""
	output_type_id := 1
	if _, ok := _item["output_type_id"]; !ok {
	} else if _, ok := _item["output_type_id"].(int); ok {
		output_type_id = _item["output_type_id"].(int)
	}
	sql_bak := _sql
	if output_type_id == 1 {
		//sql_drop = fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, destination_table)
		if app.contains([]interface{}{1, true, "1", "True", "TRUE", "T", "1"}, _item["append_it"]) {
			_sql = fmt.Sprintf(`CREATE IF NOT EXISTS TABLE "%s" AS %s`, destination_table, _sql)
		} else {
			_sql = fmt.Sprintf(`CREATE OR REPLACE TABLE "%s" AS %s`, destination_table, _sql)
		}
	} else if output_type_id == 2 {
		//sql_drop = fmt.Sprintf(`DROP VIEW IF EXISTS "%s"`, destination_table)
		_sql = fmt.Sprintf(`CREATE OR REPLACE VIEW "%s" AS %s`, destination_table, _sql)
	} else {
		msg, _ := app.i18n.T("output-type-not-suported", map[string]interface{}{"name": _item["etl_rbase_output"]})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
		}
	}
	if app.contains([]interface{}{1, true, "1", "True", "TRUE", "T", "1"}, _item["append_it"]) {
		_, err := db.ExecuteQuery(_sql, []interface{}{}...)
		if err != nil {
			fmt.Println("TRANSFORM:", err)
			app.duckdb_end(db, _duck_conf, _driver, _database, "")
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("Err Importing: %s", err),
			}
		}
		_sql_del := fmt.Sprintf(`DELETE FROM "%s"`, destination_table)
		if ref_date_field != "" {
			_sql_del = fmt.Sprintf(`%s WHERE "%s" = '%s'`, _sql, ref_date_field, date_format_org)
		}
		n_rows, err := db.ExecuteQueryRowsAffected(_sql_del, []interface{}{}...)
		if err != nil {
			fmt.Println("TRANSFORM CLEANING ERR:", err)
		}
		fmt.Printf("CLEANING BEFORE APPENDING ON %s: %d ROWS", destination_table, n_rows)
		_sql = fmt.Sprintf(`INSERT INTO "%s" VALUES %s`, destination_table, sql_bak)
	}
	_file, err := app.tempFIle(_sql, fmt.Sprintf("query.%s.*.sql", destination_table))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(_file)
	n_rows, err := db.ExecuteQueryRowsAffected(_sql, []interface{}{}...)
	if err != nil {
		fmt.Println("TRANSFORM:", err)
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("Err Importing: %s", err),
		}
	}
	// DETACH DBS TO THE DUCKDB IN MEM CONN
	app.duckdb_end(db, _duck_conf, _driver, _database, "")
	//data := map[string]interface{}{}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"n_rows":  n_rows,
	}
}
