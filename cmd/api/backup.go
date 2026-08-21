package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

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

func (app *application) _Buckup(params Dict) Dict {
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
	//admin_db_tables := strings.Split(env.GetString("EXPORT_ADMIN_DB_TABLES", ""), ",")
	etlx_obj := &etlx.ETLX{Config: Dict{}}
	//fmt.Println("APPS:", *apps)
	memDB, err := etlx.GetDB("duckdb:")
	if err != nil {
		fmt.Println("Buckup Err:", err)
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Buckup Err connecting to in memory duckdb: %s", err.Error()),
		}
	}
	defer memDB.Close()
	for _, _app := range *apps {
		fmt.Printf("Backup Start: %s -> %v\n", _app["app"], time.Now())
		memDB.ExecuteQuery(`create sequence query_id_seq start 1`)
		sql := `create or replace table "queries" (
			"id" bigint primary key default nextval('query_id_seq'),
			"query" text null,
			"admin" boolean null,
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
		_driver := ""
		if db.GetDriverName() == "sqlite3" || db.GetDriverName() == "sqlite" {
			_type = "(type sqlite)"
			_driver = "sqlite"
		} else if db.GetDriverName() == "postgres" {
			_type = "(type postgres)"
			_driver = "postgres"
		} else if db.GetDriverName() == "mysql" {
			_type = "(type mysql)"
			_driver = "mysql"
		} else if db.GetDriverName() == "odbc" {
			_type = "(type odbc)"
			_driver = "odbc"
		} else if db.GetDriverName() == "duckdb" {
			_type = ""
			_driver = "duckdb"
		}
		//fmt.Println("CHK ADM vs APP DB:", _app["db"].(string) != admin_db, _app["db"].(string), admin_db)
		if _app["db"].(string) != admin_db {
			attach := fmt.Sprintf(`attach if not exists '%s' as %s %s`, adm_dsn, admin_db, _type)
			//fmt.Println(attach)
			memDB.ExecuteQuery(attach)
			memDB.ExecuteQuery(fmt.Sprintf(`use %s`, admin_db))
			//fmt.Println(attach)
			app.InsertData(memDB, "memory.queries", Dict{"query": attach})
			if _driver == "sqlite" {
				//app.InsertData(memDB, "memory.queries", Dict{"query": "PRAGMA foreign_keys = OFF;", "admin": false})
			}
			sql = `show tables`
			res_adm_tbl, _, err := memDB.QueryMultiRows(sql, []any{}...)
			if err != nil {
				fmt.Printf("Error getting tables from the admin app: %s: %s!", admin_db, err)
				continue
			}
			for _, _adm_tbl := range *res_adm_tbl {
				adm_tbl, _ := _adm_tbl["name"].(string)
				fmt.Println("ADM TABLES:", adm_tbl)
				if adm_tbl == "" {
					continue
				}
				_sql := `select * from duckdb_columns() where table_name = ? and column_name = ?`
				result, _, err := memDB.QueryMultiRows(_sql, []any{adm_tbl, "app_id"}...)
				if err != nil {
					fmt.Printf("Error checking if table has app_id: %s->%s: %s!", admin_db, adm_tbl, err)
					continue
				}
				sql = ""
				_filter := []any{}
				if len(*result) > 0 {
					sql = fmt.Sprintf(`select * from %s."%s" where "app_id" = ?`, admin_db, adm_tbl)
					//fmt.Println("TABLE HAS APP:", admin_db, adm_tbl, sql)
					_filter = append(_filter, _app["app_id"])
				} else {
					result, _, err := memDB.QueryMultiRows(_sql, []any{adm_tbl, "db"}...)
					if err != nil {
						fmt.Printf("Error checking if table has app_id: %s->%s: %s!", admin_db, adm_tbl, err)
						continue
					}
					if len(*result) > 0 {
						sql = fmt.Sprintf(`select * from %s."%s" where "db" = ?`, admin_db, adm_tbl)
						//fmt.Println("TABLE HAS DB:", admin_db, adm_tbl, sql)
						_filter = append(_filter, dbname)
					} else {
						continue
					}
				}
				//fmt.Println("ADM:", sql)
				result, _, err = memDB.QueryMultiRows(sql, _filter...)
				if err != nil {
					fmt.Printf("Error getting the data from %s->%s: %s!", admin_db, adm_tbl, err)
					continue
				} else if len(*result) == 0 {
					continue
				}
				sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into %s."%s" (":columns") values`, admin_db, adm_tbl), *result)
				app.InsertData(memDB, "memory.queries", Dict{"query": sqls, "admin": true})
				sqls, _ = etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into "%s" (":columns") values`, adm_tbl), *result)
				app.InsertData(memDB, "memory.adm_query", Dict{"query": sqls})
			}
			if _driver == "sqlite" {
				//app.InsertData(memDB, "memory.queries", Dict{"query": "PRAGMA foreign_keys = ON;", "admin": false})
			}
			app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`detach %s`, admin_db)})
			memDB.ExecuteQuery(fmt.Sprintf(`use %s`, "memory"))
			memDB.ExecuteQuery(fmt.Sprintf(`detach %s`, admin_db))
		}
		attach := fmt.Sprintf(`attach if not exists '%s' as %s %s`, dsn2, dbname, _type)
		memDB.ExecuteQuery(attach)
		app.InsertData(memDB, "memory.queries", Dict{"query": attach})
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, dbname))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`use %s`, dbname)})
		sql = `select * from duckdb_tables() where database_name = ?`
		switch _driver {
		case "sqlite":
			//app.InsertData(memDB, "memory.queries", Dict{"query": "PRAGMA foreign_keys = OFF;", "admin": false})
		case "postgres":
			sql = `from duckdb_tables() where database_name = ? and schema_name = 'public'`
		case "mysql", "odbc", "mssql":
			sql = `from duckdb_tables() where database_name = ? and schema_name = 'default'`
		}
		tables, _, err := memDB.QueryMultiRows(sql, []any{dbname}...)
		if err != nil {
			fmt.Printf("Error getting the table %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error getting the tables from %s: %s!", _app["app"], err),
			}
		}
		for _, table := range *tables {
			if strings.HasPrefix(table["table_name"].(string), "sqlite_") {
				continue
			}
			sensitive_tables := strings.Split(os.Getenv("SENSITIVE_TABLES"), ",")
			if app.contains(app.sliceStrs2SliceInterfaces(sensitive_tables), table) {
				continue
			}
			_filter := []any{}
			sql = fmt.Sprintf(`select * from "%s"`, table["table_name"])
			//fmt.Println(table)
			_sql := `select * from duckdb_columns() where table_name = ? and column_name = ?`
			result, _, err := memDB.QueryMultiRows(_sql, []any{table["table_name"], "app_id"}...)
			if err != nil {
				fmt.Printf("Error checking if table has app_id: %s->%s: %s!", _app["db"].(string), table["table_name"], err)
				//continue
			}
			if len(*result) > 0 {
				sql = fmt.Sprintf(`select * from %s."%s" where "app_id" = ?`, _app["db"].(string), table["table_name"])
				//fmt.Println("TABLE HAS APP:", _app["db"].(string), table["table_name"], sql)
				_filter = append(_filter, _app["app_id"])
			} else {
				result, _, err := memDB.QueryMultiRows(_sql, []any{table["table_name"], "db"}...)
				if err != nil {
					fmt.Printf("Error checking if table has app_id: %s->%s: %s!", _app["db"].(string), table["table_name"], err)
					//continue
				}
				if len(*result) > 0 {
					sql = fmt.Sprintf(`select * from %s."%s" where "db" = ?`, _app["db"].(string), table["table_name"])
					_filter = append(_filter, _app["db"].(string))
					//fmt.Println("TABLE HAS DB:", _app["db"].(string), table["table_name"], sql)
				} else {
					//continue
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3600)*time.Second)
			defer cancel()
			rows, err := memDB.QueryRows(ctx, sql, _filter...)
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
			var result2 []Dict
			for rows.Next() {
				i += 1
				row, _ := app.ScanRowToMap(rows)
				result2 = append(result2, row)
				if i >= chunk_size {
					i = 0
					sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into "%s" (":columns") values`, table["table_name"]), result2)
					app.InsertData(memDB, "memory.queries", Dict{"query": sqls, "admin": true})
					app.InsertData(memDB, "memory.app_query", Dict{"query": sqls})
					result2 = []Dict{} //result[:0]
				}
			}
			if err := rows.Err(); err != nil {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting the data from %s->%s: %s!", _app["app"], table["table_name"], err),
				}
			}
			if len(result2) > 0 {
				sqls, _ := etlx_obj.BuildInsertSQL(fmt.Sprintf(`insert into "%s" (":columns") values`, table["table_name"]), result2)
				app.InsertData(memDB, "memory.queries", Dict{"query": sqls})
				app.InsertData(memDB, "memory.app_query", Dict{"query": sqls})
			}
		}
		if _driver == "sqlite" {
			//app.InsertData(memDB, "memory.queries", Dict{"query": "PRAGMA foreign_keys = ON;", "admin": false})
		}
		app.InsertData(memDB, "memory.queries", Dict{"query": "COMMIT;"})
		memDB.ExecuteQuery(fmt.Sprintf(`use %s`, "memory"))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`use %s`, "memory")})
		memDB.ExecuteQuery(fmt.Sprintf(`detach %s`, dbname))
		app.InsertData(memDB, "memory.queries", Dict{"query": fmt.Sprintf(`detach %s`, dbname)})
		/*_sql := fmt.Sprintf(`copy memory."queries" to '%s/%s.%s.csapppq' (format parquet)`, embed_dbs_dir, _app["app"], app.config.db.driverName)
		_, err = memDB.ExecuteQuery(_sql)
		if err != nil {
			fmt.Printf("Error exporting the app %s: %s!", _app["app"], err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error exporting the app %s: %s!", _app["app"], err),
			}
		}*/
		filename := fmt.Sprintf(`%s/%s.%s.csapp`, embed_dbs_dir, _app["app"], app.config.db.driverName)
		if err := os.Remove(filename); err != nil {
			fmt.Printf("could not delete %s: %v", filename, err)
			//continue
		}
		// upload to s3
		if app.config.useS3 {
			// Upload to S3
			file, err := os.Open(filename)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				defer file.Close()
				fname, err := app.uploadToS3(file, fmt.Sprintf("%s.%s.csapp", _app["app"], app.config.db.driverName))
				if err != nil {
					fmt.Println(err.Error())
					/*return Dict{
						"success": false,
						"msg":     "Failed to upload to S3: " + err.Error(),
					}*/
				}
				fmt.Printf("Uploaded to S3: %s\n", fname)
			}
		}
		attch := fmt.Sprintf(`attach '%s' as %s`, filename, _app["app"])
		memDB.ExecuteQuery(attch)
		memDB.ExecuteQuery(fmt.Sprintf(`copy from database memory to %s`, _app["app"]))
		memDB.ExecuteQuery(fmt.Sprintf(`DETACH %s`, _app["app"]))
		fmt.Printf("Backup End: %s -> %v\n", _app["app"], time.Now())
	}
	msg, _ := app.i18n.T("success", Dict{})
	data := Dict{
		"success": true,
		"msg":     msg,
	}
	return data
}

func (app *application) Backup(params Dict) Dict {
	// -------------------------------------------------------------------------
	// Get the admin database
	// -------------------------------------------------------------------------
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	// _, adminDSN, _ := app.ParseConnection(dsn)
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error getting the db connection: %s", err),
		}
	}
	defer db.Close()
	// -------------------------------------------------------------------------
	// Get applications to backup
	// -------------------------------------------------------------------------
	sql := `select * from "app" where excluded = false and lower("app") like ?`
	appName := "%"
	if data, ok := params["data"].(Dict); ok {
		if name, ok := data["name"].(string); ok && name != "" {
			appName = name
		}
	}
	apps, _, err := db.QueryMultiRows(sql, []any{strings.ToLower(appName)}...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error getting the apps: %s", err),
		}
	}
	// -------------------------------------------------------------------------
	// Create an in-memory DuckDB connection.
	//
	// This connection is used as the DuckDB coordinator for ATTACH,
	// EXPORT DATABASE and the S3 httpfs extension.
	// -------------------------------------------------------------------------
	memDB, err := etlx.GetDB("duckdb:")
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error connecting to in-memory DuckDB: %s", err),
		}
	}
	defer memDB.Close()
	// -------------------------------------------------------------------------
	// S3 configuration
	// -------------------------------------------------------------------------
	bucket := strings.TrimSpace(os.Getenv("BACKUP_S3_BUCKET"))
	region := strings.TrimSpace(os.Getenv("BACKUP_S3_REGION"))
	endpoint := strings.TrimSpace(os.Getenv("BACKUP_S3_ENDPOINT"))
	accessKey := os.Getenv("BACKUP_S3_ACCESS_KEY_ID")
	secretKey := os.Getenv("BACKUP_S3_SECRET_ACCESS_KEY")
	urlStyle := strings.ToLower(strings.TrimSpace(os.Getenv("BACKUP_S3_URL_STYLE")))
	if bucket == "" {
		return Dict{
			"success": false,
			"msg":     "BACKUP_S3_BUCKET is not configured",
		}
	}
	if accessKey == "" {
		return Dict{
			"success": false,
			"msg":     "BACKUP_S3_ACCESS_KEY_ID is not configured",
		}
	}
	if secretKey == "" {
		return Dict{
			"success": false,
			"msg":     "BACKUP_S3_SECRET_ACCESS_KEY is not configured",
		}
	}
	if region == "" {
		region = "us-east-1"
	}
	// DuckDB accepts "path" and "vhost".
	if urlStyle != "" && urlStyle != "path" && urlStyle != "vhost" {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("invalid BACKUP_S3_URL_STYLE %q, expected path or vhost", urlStyle),
		}
	}
	// -------------------------------------------------------------------------
	// Create DuckDB S3 secret
	// -------------------------------------------------------------------------
	secretSQL := fmt.Sprintf(`
CREATE OR REPLACE SECRET c7_backup_s3 (
	TYPE s3,
	PROVIDER config,
	KEY_ID %s,
	SECRET %s,
	REGION %s
`, sqlStringLiteral(accessKey), sqlStringLiteral(secretKey), sqlStringLiteral(region))
	if endpoint != "" {
		secretSQL += fmt.Sprintf(",\n\tENDPOINT %s\n", sqlStringLiteral(endpoint))
	}
	if urlStyle != "" {
		secretSQL += fmt.Sprintf(",\n\tURL_STYLE %s\n", sqlStringLiteral(urlStyle))
	}
	secretSQL += `);`
	if _, err := memDB.ExecuteQuery(secretSQL); err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("error configuring S3 backup credentials: %s", err),
		}
	}
	// fmt.Println(secretSQL)
	// -------------------------------------------------------------------------
	// One timestamp for the complete backup run.
	//
	// Every application backed up during this invocation gets the same
	// timestamp, making it easy to identify a backup run.
	// -------------------------------------------------------------------------
	backupTS := time.Now().UTC().Format("20060102T150405Z")
	fmt.Printf("Backup started: %s\n", backupTS)
	// -------------------------------------------------------------------------
	// Backup every application
	// -------------------------------------------------------------------------
	for _, application := range *apps {
		appName, ok := application["app"].(string)
		if !ok || appName == "" {
			fmt.Printf("Skipping application with invalid app name: %v\n", application)
			continue
		}
		fmt.Printf("Backup Start: %s -> %v\n", appName, time.Now())
		// -------------------------------------------------------------
		// Get the application's database
		// -------------------------------------------------------------
		appDBName, ok := application["db"].(string)
		if !ok || appDBName == "" {
			fmt.Printf("Skipping application %s: database name is missing\n", appName)
			continue
		}
		appDSN, dbName, _ := app.GetDBNameFromParams(Dict{"db": appDBName})
		appDBCon, err := etlx.GetDB(appDSN)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("error getting the app DB %s: %s", appName, err),
			}
		}
		// Do not defer this inside the loop. Close it after this application
		// has been backed up.
		_, appDSN2, _ := app.ParseConnection(appDSN)
		// -------------------------------------------------------------
		// Determine source database type.
		//
		// This preserves the existing ATTACH behavior.
		// -------------------------------------------------------------
		dbType := "(type sqlite)"
		//driver := "sqlite"
		switch appDBCon.GetDriverName() {
		case "sqlite3", "sqlite":
			dbType = "(type sqlite)"
			//driver = "sqlite"
		case "postgres":
			dbType = "(type postgres)"
			//driver = "postgres"
		case "mysql":
			dbType = "(type mysql)"
			//driver = "mysql"
		case "odbc":
			dbType = "(type odbc)"
			//driver = "odbc"
		case "mssql":
			dbType = "(type mssql)"
			//driver = "mssql"
		case "duckdb":
			dbType = ""
			//driver = "duckdb"
		}
		// -------------------------------------------------------------
		// Attach the application database to DuckDB.
		// -------------------------------------------------------------
		attachSQL := fmt.Sprintf(`ATTACH IF NOT EXISTS %s AS %s %s`, sqlStringLiteral(appDSN2), quoteDuckDBIdentifier(dbName), dbType)
		// fmt.Println(attachSQL)
		if _, err := memDB.ExecuteQuery(attachSQL); err != nil {
			appDBCon.Close()
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("error attaching database %s for application %s: %s", dbName, appName, err),
			}
		}
		// -------------------------------------------------------------
		// Select the application database.
		// -------------------------------------------------------------
		useSQL := fmt.Sprintf(`USE %s`, quoteDuckDBIdentifier(dbName))
		// fmt.Println(useSQL)
		if _, err := memDB.ExecuteQuery(useSQL); err != nil {
			memDB.ExecuteQuery(`USE memory`)
			memDB.ExecuteQuery(fmt.Sprintf(`DETACH %s`, quoteDuckDBIdentifier(dbName)))
			appDBCon.Close()
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("error selecting database %s: %s", dbName, err),
			}
		}
		// -------------------------------------------------------------
		// Build the S3 destination.
		//
		// Example:
		//
		// s3://c7-backups/my_database/20260821T143500Z/
		// -------------------------------------------------------------
		backupPath := fmt.Sprintf("s3://%s/%s/%s", strings.Trim(bucket, "/"), sanitizeS3PathPart(dbName), backupTS)
		fmt.Printf("Exporting database %s -> %s\n", dbName, backupPath)
		// -------------------------------------------------------------
		// Export the complete DuckDB database as Parquet.
		//
		// DuckDB creates its export files under backupPath, including
		// schema.sql, load.sql and the table parquet files.
		// -------------------------------------------------------------
		exportSQL := fmt.Sprintf(`EXPORT DATABASE %s (FORMAT parquet);`, sqlStringLiteral(backupPath))
		// fmt.Println(exportSQL)
		if _, err := memDB.ExecuteQuery(exportSQL); err != nil {
			// Always detach before returning.
			memDB.ExecuteQuery(`USE memory`)
			memDB.ExecuteQuery(fmt.Sprintf(`DETACH %s`, quoteDuckDBIdentifier(dbName)))
			appDBCon.Close()
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("error exporting database %s for application %s to S3: %s", dbName, appName, err),
			}
		}
		fmt.Printf("Backup exported: %s -> %s\n", appName, backupPath)
		// -------------------------------------------------------------
		// Detach application database.
		// -------------------------------------------------------------
		memDB.ExecuteQuery(`USE memory`)
		detachSQL := fmt.Sprintf(`DETACH %s`, quoteDuckDBIdentifier(dbName))
		if _, err := memDB.ExecuteQuery(detachSQL); err != nil {
			appDBCon.Close()
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("error detaching database %s: %s", dbName, err),
			}
		}
		appDBCon.Close()
		fmt.Printf("Backup End: %s -> %v\n", appName, time.Now())
	}
	// -------------------------------------------------------------------------
	// Done
	// -------------------------------------------------------------------------
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
	}
}

func sqlStringLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
func quoteDuckDBIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}
func sanitizeS3PathPart(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "/", "_")
	value = strings.ReplaceAll(value, "\\", "_")
	return value
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
