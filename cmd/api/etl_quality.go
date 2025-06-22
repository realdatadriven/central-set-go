package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/realdatadriven/etlx"
	// "github.com/google/uuid"
)

func (app *application) _quality(params map[string]interface{}, _item map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, db_conf map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
	// IN MEMORY DUCKDB CONN
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("DDB Conn: %s", err),
		}
	}
	defer db.Close()
	// DB
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
	// GET THE OUTPUT FIELDS
	_table := ""
	if _, ok := _step["table"]; ok {
		_table = _step["table"].(string)
	} else {
		_table = "etl_rbase_quality"
	}
	etl_report_base_id := interface{}(-1)
	if _, ok := _item["etl_report_base_id"]; ok {
		etl_report_base_id = _item["etl_report_base_id"]
	}
	params["data"].(map[string]interface{})["table"] = _table
	params["data"].(map[string]interface{})["limit"] = interface{}(-1.0)
	params["data"].(map[string]interface{})["offset"] = interface{}(0.0)
	params["data"].(map[string]interface{})["filters"] = []interface{}{}
	params["data"].(map[string]interface{})["filters"] = append(
		params["data"].(map[string]interface{})["filters"].([]interface{}),
		map[string]interface{}{
			"field": "etl_report_base_id",
			"cond":  "=",
			"value": etl_report_base_id,
		},
	)
	params["data"].(map[string]interface{})["order_by"] = []interface{}{}
	params["data"].(map[string]interface{})["order_by"] = append(
		params["data"].(map[string]interface{})["order_by"].([]interface{}),
		map[string]interface{}{
			"field": "rule_order",
			"order": "ASC",
		},
	)
	params["data"].(map[string]interface{})["order_by"] = append(
		params["data"].(map[string]interface{})["order_by"].([]interface{}),
		map[string]interface{}{
			"field": "etl_rbase_quality_id",
			"order": "ASC",
		},
	)
	_queries := map[string]interface{}{}
	_qualities := app.read(params)
	if _, ok := _qualities["success"]; !ok {
		return _qualities
	} else if _qualities["success"].(bool) {
		//fmt.Println(_qualities["sql"])
		for _, _query := range _qualities["data"].([]map[string]interface{}) {
			if _, ok := _query["sql_quality_check"]; !ok {
			} else if _active, ok := _query["active"]; ok {
				if app.contains([]interface{}{true, 1, "1", "true", "True", "TRUE"}, _active) {
					_sql := _query["sql_quality_check"].(string)
					_sql = app.setQueryDate(_sql, date_ref)
					_queries[_query["etl_rbase_quality_id"].(string)] = _sql
				}
			}
		}
	}
	params["data"].(map[string]interface{})["query"] = _queries
	params["data"].(map[string]interface{})["database"] = _database
	_exec_queries := app.query(params)
	fmt.Println(_exec_queries)
	// DETACH DBS TO THE DUCKDB IN MEM CONN
	app.duckdb_end(db, _duck_conf, _driver, _database, "")
	//data := map[string]interface{}{}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		// "n_rows":  n_rows,
	}
}
