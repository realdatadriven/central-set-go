package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/etlx"
)

type Dict = map[string]any

func AddForeignKeyToCreateStmt(createStmt, fkString string) string {
	// Trim semicolon if exists
	createStmt = strings.TrimSuffix(createStmt, ";")
	// Prepare regex to match the last closing parenthesis before the semicolon
	re := regexp.MustCompile(`(?i)(?s)(\))\s*$`) // matches last ')'
	if !re.MatchString(createStmt) {
		return createStmt // fallback: no match
	}
	// Insert FK string before the last ')', with a comma
	return re.ReplaceAllString(createStmt, ",\n    "+fkString+"\n)")
}

func (app *application) ScanRowToMap(rows *sql.Rows) (Dict, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))
	for i := range values {
		valuePointers[i] = &values[i]
	}
	if err := rows.Scan(valuePointers...); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	rowMap := make(Dict)
	for i, colName := range columns {
		rowMap[colName] = values[i]
	}
	return rowMap, nil
}

func (app *application) Buckup(params Dict) Dict {
	dsn, admin_db, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	_, adm_dsn, _ := app.ParseConnection(dsn)
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
	if _, ok := params["path"].(string); ok {
		embed_dbs_dir = params["path"].(string)
	} else if os.Getenv("DB_EMBEDED_DIR") != "" {
		embed_dbs_dir = os.Getenv("DB_EMBEDED_DIR")
	}
	admin_db_tables := strings.Split(env.GetString("EXPORT_ADMIN_DB_TABLES", ""), ",")
	etlx_obj := &etlx.ETLX{Config: Dict{}}
	//fmt.Println("APPS:", *apps)
	memDB, _ := etlx.GetDB("duckdb:")
	defer memDB.Close()
	for _, _app := range *apps {
		fmt.Printf("Backup Start: %s -> %v\n", _app["app"], time.Now())
		memDB.ExecuteQuery(`create sequence query_id_seq start 1`)
		sql := `create or replace table "queries" (
			"id" bigint primary key default nextval('query_id_seq'),
			"query" text null,
    		"created_at" timestamp default current_timestamp
		)`
		memDB.ExecuteQuery(sql)
		sql = `create or replace table "adm_query" (
			"id" bigint primary key default nextval('query_id_seq'),
			"query" text null,
    		"created_at" timestamp default current_timestamp
		)`
		memDB.ExecuteQuery(sql)
		sql = `create or replace table "app_query" (
			"id" bigint primary key default nextval('query_id_seq'),
			"query" text null,
    		"created_at" timestamp default current_timestamp
		)`
		memDB.ExecuteQuery(sql)
		err := app.InsertData(memDB, "memory.queries", Dict{"query": "BEGIN TRANSACTION;"})
		if err != nil {
			fmt.Printf("Error executing query %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error executing query %s: %s!", _app["app"], err),
			}
		}
		dsn, dbname, _ := app.GetDBNameFromParams(Dict{"db": _app["db"]})
		appDBCon, err := etlx.GetDB(dsn)
		if err != nil {
			fmt.Printf("Error getting the app DB %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error getting the app DB %s: %s!", _app["app"], err),
			}
		}
		defer appDBCon.Close()
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
		if _app["db"].(string) != admin_db {
			attach := fmt.Sprintf(`attach if not exists '%s' as %s %s`, adm_dsn, admin_db, _type)
			app.InsertData(memDB, "memory.queries", Dict{"query": attach})
			memDB.ExecuteQuery(attach)
			for _, adm_tbl := range admin_db_tables {
				if adm_tbl == "" {
					continue
				}
				sql = fmt.Sprintf(`select * from %s."%s" where "db" = ?`, admin_db, adm_tbl)
				result, _, err := memDB.QueryMultiRows(sql, []any{dbname}...)
				if err != nil {
					fmt.Printf("Error getting the data from %s->%s: %s!", admin_db, adm_tbl, err)
				}
				sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into %s."%s" (":columns") values`, admin_db, adm_tbl), *result)
				app.InsertData(memDB, "memory.queries", Dict{"query": sqls})
				app.InsertData(memDB, "memory.adm_query", Dict{"query": sqls})
			}
			app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`detach %s`, admin_db)})
			memDB.ExecuteQuery(fmt.Sprintf(`detach %s`, admin_db))
		}
		attach := fmt.Sprintf(`attach if not exists '%s' as %s %s`, dsn2, dbname, _type)
		memDB.ExecuteQuery(attach)
		app.InsertData(memDB, "memory.queries", Dict{"query": attach})
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, dbname))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`use %s`, dbname)})
		sql = `select * from duckdb_tables() where database_name = ?`
		tables, _, err := memDB.QueryMultiRows(sql, []any{dbname}...)
		if err != nil {
			fmt.Printf("Error getting the table %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error getting the tables from %s: %s!", _app["app"], err),
			}
		}
		for _, table := range *tables {
			if table["table_name"] == "sqlite_sequence" || table["table_name"] == "sqlite_stat" {
				continue
			}
			sql = table["sql"].(string)
			/*extra_conf := Dict{"driverName": app.config.db.driverName, "dsn": app.config.db.dsn}
			schema, _, err := appDBCon.TableSchema(params, table["table_name"].(string), _app["db"].(string), extra_conf)
			if err != nil {
				fmt.Printf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err)
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err),
				}
			}
			fks := []string{}
			for _, col := range *schema {
				if col["fk"].(bool) {
					//fmt.Println(col)
					fks = append(fks, fmt.Sprintf(`FOREIGN KEY("%s") REFERENCES "%s"("%s")`, col["field"], col["referred_table"], col["referred_column"]))
				} else if col["pk"].(bool) {

				}
			}
			if len(fks) > 0 {
				sql = AddForeignKeyToCreateStmt(sql, app.joinSlice(app.sliceStrs2SliceInterfaces(fks), ","))
			}*/
			app.InsertData(memDB, "memory.queries", Dict{"query": sql})
			sql = fmt.Sprintf(`select * from "%s"`, table["table_name"])
			db_filter := []any{}
			if _app["db"].(string) == admin_db && app.contains(app.sliceStrs2SliceInterfaces(admin_db_tables), table["table_name"]) {
				sql = fmt.Sprintf(`select * from "%s" where "db" = ?`, table["table_name"])
				db_filter = []any{admin_db}
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3600)*time.Second)
			defer cancel()
			rows, err := memDB.QueryRows(ctx, sql, db_filter...)
			if err != nil {
				fmt.Printf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err)
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err),
				}
			}
			defer rows.Close()
			chunk_size := 500
			i := 0
			var result []Dict
			for rows.Next() {
				i += 1
				row, _ := app.ScanRowToMap(rows)
				result = append(result, row)
				if i >= chunk_size {
					i = 0
					sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into "%s" (":columns") values`, table["table_name"]), result)
					app.InsertData(memDB, "memory.queries", Dict{"query": sqls})
					app.InsertData(memDB, "memory.app_query", Dict{"query": sqls})
					result = []Dict{} //result[:0]
				}
			}
			if err := rows.Err(); err != nil {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err),
				}
			}
			if len(result) > 0 {
				sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into "%s" (":columns") values`, table["table_name"]), result)
				app.InsertData(memDB, "memory.queries", Dict{"query": sqls})
				app.InsertData(memDB, "memory.app_query", Dict{"query": sqls})
			}
		}
		app.InsertData(memDB, "memory.queries", Dict{"query": "COMMIT;"})
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, "memory"))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`use %s`, "memory")})
		memDB.ExecuteQuery(fmt.Sprintf(`detach %s`, dbname))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`detach %s`, dbname)})
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
