package main

import (
	"fmt"
	"os"
	"time"

	"github.com/realdatadriven/etlx"
)

type Dict = map[string]any

func (app *application) Buckup(params Dict) Dict {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	// fmt.Println(dsn)
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error geting the db connection: %s", err),
		}
	}
	defer db.Close()
	sql := `select * from "app" where excluded = false and "app" like ?`
	// fmt.Println(sql)
	_app := "%"
	if a, ok := params["data"].(Dict)["name"].(string); ok && a != "" {
		_app = params["data"].(Dict)["name"].(string)
	}
	apps, _, err := db.QueryMultiRows(sql, []any{_app}...)
	if err != nil {
		fmt.Printf("error geting the apps: %s\n", err)
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error geting the apps: %s", err),
		}
	}
	embed_dbs_dir := "database"
	if os.Getenv("DB_EMBEDED_DIR") != "" {
		embed_dbs_dir = os.Getenv("DB_EMBEDED_DIR")
	}
	//fmt.Println("APPS:", *apps)
	memDB, _ := etlx.GetDB("duckdb:")
	defer memDB.Close()
	for _, _app := range *apps {
		fmt.Printf("Backup Start: %s -> %v\n", _app["app"], time.Now())
		memDB.ExecuteQuery(`CREATE OR REPLACE TABLE "queries" ("query" TEXT NULL)`)
		err := app.InsertData(memDB, "memory.queries", Dict{"query": "BEGIN TRANSACTION;"})
		if err != nil {
			fmt.Printf("Error executing query %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error executing query %s: %s!", _app["app"], err),
			}
		}
		dsn, dbname, _ := app.GetDBNameFromParams(Dict{"db": _app["app"]})
		_, dsn2, _ := app.ParseConnection(dsn)
		_type := ""
		if db.GetDriverName() == "sqlite3" || db.GetDriverName() == "sqlite" {
			_type = "(type sqlite)"
		} else if db.GetDriverName() == "postgres" {
			_type = "(type postgres)"
		} else if db.GetDriverName() == "mysql" {
			_type = "(type mysql)"
		} else if db.GetDriverName() == "odbc" {
			_type = "(type odbc)"
		} else if db.GetDriverName() == "duckdb" {
			_type = ""
		}
		attach := fmt.Sprintf(`attach '%s' as %s %s`, dsn2, dbname, _type)
		//fmt.Println(1, attach)
		memDB.ExecuteQuery(attach)
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, dbname))
		// EXPORT TABLE SQL
		sql = `SELECT * FROM duckdb_tables() where database_name = ?`
		tables, _, err := memDB.QueryMultiRows(sql, []any{dbname}...)
		if err != nil {
			fmt.Printf("Error getting the table %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error getting the table %s: %s!", _app["app"], err),
			}
		}
		for _, table := range *tables {
			if table["table_name"] == "sqlite_sequence" {
				continue
			}
			app.InsertData(memDB, "memory.queries", Dict{"query": table["sql"]})
			// TABLE DATA adapt from etlx db2db
		}
		app.InsertData(memDB, "memory.queries", Dict{"query": "COMMIT;"})
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, "memory"))
		memDB.ExecuteQuery(fmt.Sprintf(`detach %s`, dbname))
		_sql := fmt.Sprintf(`copy memory."queries" to '%s/%s.%s.csapp' (format parquet)`, embed_dbs_dir, _app["app"], app.config.db.driverName)
		_, err = memDB.ExecuteQuery(_sql)
		if err != nil {
			fmt.Printf("Error exporting the app %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error exporting the app %s: %s!", _app["app"], err),
			}
		}
		fmt.Printf("Backup End: %s -> %v\n", _app["app"], time.Now())
	}
	msg, _ := app.i18n.T("success", Dict{})
	data := Dict{
		"success": true,
		"msg":     msg,
	}
	return data
}

func (app *application) InsertData(db etlx.DBInterface, table string, data Dict) error {
	var columns []any
	var placeholders []any
	var values []any
	for key, val := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}
	cols := app.joinSlice(columns, `", "`)
	plch := app.joinSlice(placeholders, `, `)
	sql := fmt.Sprintf(`insert into %s ("%s") values (%s)`, table, cols, plch)
	//fmt.Println(sql, values)
	_, err := db.ExecuteQuery(sql, values...)
	if err != nil {
		return err
	}
	return nil
}
