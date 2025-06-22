package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/realdatadriven/etlx"
	// "github.com/google/uuid"
)

func (app *application) _delete(_item map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, db_conf map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
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
	// RUN EXTRACT HERE
	_db_base := filepath.Base(_database)
	_db_ext := filepath.Ext(_database)
	_db_name_no_ext := _db_base[:len(_db_base)-len(_db_ext)]
	//fmt.Println(_driver, _database, _db_base, _db_ext, _db_name_no_ext)
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
	//fmt.Println(date_ref)
	// CHECK DATE
	check_ref_date := false
	if _, ok := _item["check_ref_date"]; !ok {
	} else if app.contains([]interface{}{true, 1, "1", "true", "True", "TRUE"}, _item["check_ref_date"]) {
		check_ref_date = true
	}
	ref_date_field := ""
	if _, ok := _item["ref_date_field"]; !ok {
	} else if _, ok := _item["ref_date_field"].(string); ok {
		ref_date_field = _item["ref_date_field"].(string)
	}
	if _, ok := _item["date_field"].(string); ok {
		ref_date_field = _item["date_field"].(string)
		check_ref_date = true
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
	_sql := fmt.Sprintf(`DELETE FROM %s."%s"`, _db_name_no_ext, destination_table)
	if check_ref_date && ref_date_field != "" {
		_sql = fmt.Sprintf(`%s WHERE "%s" = '%s'`, _sql, ref_date_field, date_format_org)
	}
	_sql = app.setQueryDate(_sql, date_ref)
	_sql = app.setStrEnv(_sql)
	//fmt.Println(_sql)
	n_rows, err := db.ExecuteQueryRowsAffected(_sql, []interface{}{}...)
	if err != nil {
		fmt.Println("DELETE:", err, _sql)
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("Err Importing: %s", err),
		}
	}
	/*n_rows_res, _, err := db.QuerySingleRow(_sql, []interface{}{}...)
	if err != nil {
		fmt.Println("EXTRACT NROWS:", err, _sql)
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("Err Importing: %s", err),
		}
	}
	n_rows := (*n_rows_res)["n_rows"]*/
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
