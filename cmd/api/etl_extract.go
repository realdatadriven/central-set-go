package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"
	// "github.com/google/uuid"
)

func (app *application) ETLExtract(params map[string]interface{}, db_conf map[string]interface{}, _input map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
	fname := ""
	if _, ok := _input["file"]; ok {
		fname = _input["file"].(string)
	}
	is_tmp := false
	if _, ok := _input["save_only_temp"]; ok {
		if app.contains([]interface{}{true, 1, "1", "true", "True", "TRUE"}, _input["save_only_temp"]) {
			is_tmp = true
		}
	}
	//var _path *os.File
	_path := fmt.Sprintf("static/uploads/%s", fname)
	if is_tmp {
		_path = fmt.Sprintf("%s/%s", os.TempDir(), fname)
	}
	if fname != "" && !app.fileExists(_path) {
		msg, _ := app.i18n.T("file-not-founded", map[string]interface{}{"fname": fname})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	_extract_type := ""
	if _, ok := _conf["type"]; ok {
		_extract_type = _conf["type"].(string)
	}
	if fname != "" && app.fileExists(_path) {
		// HANDLE SEND TO DUCKDB
		return app._extract(params, _path, _input, _conf, _etlrb, _conf_etlrb, db_conf, _step)
	} else if _extract_type == "odbc-csv-duckdb" {
		// EXTRACT ODBC 2 CSV THEN SEND FILE TO DUCKDB
		return app._odbc_csv_duckdb(params, _input, _conf, _etlrb, _conf_etlrb, db_conf, _step)
	} else if _extract_type == "duckdb" {
		// SEND TO DUCK
		return app._extract(params, "", _input, _conf, _etlrb, _conf_etlrb, db_conf, _step)
	} else if fname != "" && !app.fileExists(_path) {
		msg, _ := app.i18n.T("file-not-founded", map[string]interface{}{"fname": fname})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	data := map[string]interface{}{}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"data":    data,
	}
}

func (app *application) duckdb_start(db *etlx.DuckDB, _duck_conf map[string]interface{}, _driver string, _database string) {
	_db_base := filepath.Base(_database)
	_db_ext := filepath.Ext(_database)
	_db_name_no_ext := _db_base[:len(_db_base)-len(_db_ext)]
	extensions := []interface{}{}
	if _, ok := _duck_conf["extensions"]; ok {
		extensions = _duck_conf["extensions"].([]interface{})
	}
	// fmt.Println(extensions)
	if len(extensions) > 0 {
		for _, ext := range extensions {
			_sql := fmt.Sprintf(`LOAD %s`, ext.(string))
			_, err := db.ExecuteQuery(_sql)
			if err != nil {
				fmt.Println("EXTENSION", err, _sql)
			}
		}
	}
	pragmas_config_sql_start := []interface{}{}
	if _, ok := _duck_conf["pragmas_config_sql_start"]; ok {
		pragmas_config_sql_start = _duck_conf["pragmas_config_sql_start"].([]interface{})
	}
	if len(pragmas_config_sql_start) > 0 {
		for _, _sql := range pragmas_config_sql_start {
			//fmt.Println(_sql)
			_, err := db.ExecuteQuery(_sql.(string))
			if err != nil {
				fmt.Println("START", err, _sql)
			}
		}
	}
	if _driver == "duckdb" && _database != "" {
		if _db_name_no_ext == "" {
			_database = fmt.Sprintf(`%s.duckdb`, _database)
		}
		_sql := fmt.Sprintf(`ATTACH 'database\%s' AS %s`, _database, _db_name_no_ext)
		_, err := db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("ATTACH", err, _sql)
		}
		_sql = fmt.Sprintf(`USE %s`, _db_name_no_ext)
		_, err = db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("USE", err, _sql)
		}
	} else if app.contains([]interface{}{"sqlite", "sqlite3"}, _driver) && _database != "" {
		if _db_name_no_ext == "" {
			_database = fmt.Sprintf(`%s.db`, _database)
		}
		_sql := fmt.Sprintf(`ATTACH 'database\%s' AS %s (TYPE SQLITE)`, _database, _db_name_no_ext)
		_, err := db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("ATTACH", err, _sql)
		}
		_sql = fmt.Sprintf(`USE %s`, _db_name_no_ext)
		_, err = db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("USE", err, _sql)
		}
	}
}

func (app *application) duckdb_end(db *etlx.DuckDB, _duck_conf map[string]interface{}, _driver string, _database string, fileExt string) {
	_db_base := filepath.Base(_database)
	_db_ext := filepath.Ext(_database)
	_db_name_no_ext := _db_base[:len(_db_base)-len(_db_ext)]
	pragmas_config_sql_end := []interface{}{}
	if _, ok := _duck_conf["pragmas_config_sql_end"]; ok {
		pragmas_config_sql_end = _duck_conf["pragmas_config_sql_end"].([]interface{})
	}
	if len(pragmas_config_sql_end) > 0 {
		for _, _sql := range pragmas_config_sql_end {
			_, err := db.ExecuteQuery(_sql.(string))
			if err != nil {
				fmt.Println("END", err, _sql)
			}
		}
	}
	if _driver == "duckdb" && _database != "" {
		db.ExecuteQuery(`USE memory`)
		_sql := fmt.Sprintf(`DETACH %s`, _db_name_no_ext)
		_, err := db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("DETACH", err, _sql)
		}
	} else if app.contains([]interface{}{"sqlite", "sqlite3"}, _driver) && _database != "" {
		db.ExecuteQuery(`USE memory`)
		_sql := fmt.Sprintf(`DETACH %s`, _db_name_no_ext)
		_, err := db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("DETACH", err, _sql)
		}
	}
	if app.contains([]interface{}{".duckdb", ".ddb"}, fileExt) {
		_, err := db.ExecuteQuery(`DETACH FILE`)
		if err != nil {
			fmt.Println("DETACH FILE", err)
		}
	}
	err := db.Close()
	if err != nil {
		fmt.Println("CLOSING:", err)
	}
}

func (app *application) get_ref_from_file(file string) time.Time {
	basename := file
	fileRefPats := []struct {
		patt *regexp.Regexp
		fmrt string
	}{
		{patt: regexp.MustCompile(`\d{8}`), fmrt: "20060102"}, // (\d{8})(?!.*\d+)
		{patt: regexp.MustCompile(`\d{6}`), fmrt: "200601"},   // (\d{6})(?!.*\d+)
		{patt: regexp.MustCompile(`\d{4}`), fmrt: "0601"},     // (\d{4})(?!.*\d+)
	}
	// This will hold the final file_ref value
	var fileRef time.Time
	// Loop through the patterns and try to match
	for _, patt := range fileRefPats {
		// Find all matches for the current pattern
		matches := patt.patt.FindAllString(basename, -1)
		if len(matches) > 0 {
			// If a match is found, attempt to parse it into a date
			matchStr := matches[0]
			dt, err := time.Parse(patt.fmrt, matchStr)
			if err != nil {
				// Handle parse error
				fmt.Println("Error parsing date:", err)
				break
			}
			if patt.fmrt == "200601" || patt.fmrt == "0601" {
				// Calculate the last day of the month for the parsed date
				year, month := dt.Year(), dt.Month()
				// Find the last day of the month
				lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
				// Create a new date with the last day of the month
				fileRef = time.Date(year, month, lastDay, 0, 0, 0, 0, time.UTC)
			} else {
				fileRef = dt
			}
			// Break the loop once a match is found and processed
			break
		}
	}
	return fileRef
}

func (app *application) set_sql_file_table(key string, sql string, file_table string) string {
	pats := map[string]*regexp.Regexp{
		"file":  regexp.MustCompile(`<file>|<filename>|<fname>|<file_name>|{file}|{filename}|{fname}|{file_name}`), // (?i)
		"table": regexp.MustCompile(`<table>|<table_name>|<tablename>|{table}|{table_name}|{tablename}`),
		"tmp":   regexp.MustCompile(`<tmp_path>|<tmp>|{tmp_path}|{tmp}`), // (?i)
	}
	re := pats[key]
	return re.ReplaceAllString(sql, file_table)
}

func getDtFmrt(format string) string {
	go_fmrt := format
	formats := []struct {
		frmt    string
		go_fmrt string
	}{
		{`YYYY|AAAA`, "2006"},
		{`YY|AA`, "06"},
		{`MM`, "01"},
		{`DD`, "02"},
		{`HH`, "15"},
		{`mm`, "04"},
		{`SS`, "05"},
		{`TSTAMP|STAMP`, "20060102150405"},
	}
	for _, f := range formats {
		re := regexp.MustCompile(f.frmt)
		go_fmrt = re.ReplaceAllString(go_fmrt, f.go_fmrt)
	}
	return go_fmrt
}

// setQueryDate formats the query string by inserting the given date reference in place of placeholders
func (app *application) setQueryDate(query string, dateRef interface{}) string {
	patt := regexp.MustCompile(`(["]?\w+["]?\.\w+\s?=\s?'\{.*?\}'|["]?\w+["]?\s?=\s?'\{.*?\}')`)
	matches := patt.FindAllString(query, -1)
	if len(matches) == 0 {
		patt = regexp.MustCompile(`["]?\w+["]?\s?=\s?'\{.*?\}'`)
		matches = patt.FindAllString(query, -1)
	}
	if len(matches) > 0 {
		patt2 := regexp.MustCompile(`'\{.*?\}'`)
		for _, m := range matches {
			format := patt2.FindString(m)
			if format != "" {
				frmtFinal := getDtFmrt(format)
				frmtFinal = strings.ReplaceAll(frmtFinal, "{", "")
				frmtFinal = strings.ReplaceAll(frmtFinal, "}", "")
				var procc string
				if dates, ok := dateRef.([]time.Time); ok {
					dts := []string{}
					for _, dt := range dates {
						dts = append(dts, dt.Format(frmtFinal))
					}
					procc = regexp.MustCompile(patt2.String()).ReplaceAllString(m, fmt.Sprintf("(%s)", strings.Join(dts, ",")))
					patt3 := regexp.MustCompile(`\s?=\s?`)
					procc = patt3.ReplaceAllString(procc, " IN ")
				} else if dt, ok := dateRef.(time.Time); ok {
					procc = regexp.MustCompile(patt2.String()).ReplaceAllString(m, dt.Format(frmtFinal))
				}
				patt = regexp.MustCompile(regexp.QuoteMeta(m))
				query = patt.ReplaceAllString(query, procc)
			}
		}
	}
	// Replace remaining date placeholders
	patt = regexp.MustCompile(`'?\{.*?\}'?`)
	matches = patt.FindAllString(query, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			frmtFinal := getDtFmrt(m)
			frmtFinal = strings.ReplaceAll(frmtFinal, "{", "")
			frmtFinal = strings.ReplaceAll(frmtFinal, "}", "")
			var procc string
			if dates, ok := dateRef.([]time.Time); ok {
				procc = regexp.MustCompile(patt.String()).ReplaceAllString(m, dates[0].Format(frmtFinal))
			} else if dt, ok := dateRef.(time.Time); ok {
				procc = regexp.MustCompile(patt.String()).ReplaceAllString(m, dt.Format(frmtFinal))
			}
			patt = regexp.MustCompile(regexp.QuoteMeta(m))
			query = patt.ReplaceAllString(query, procc)
		}
	}
	// Handle cases for temporary tables with date extensions
	patt = regexp.MustCompile(
		`YYYY.?MM.?DD|AAAA.?MM.?DD|YY.?MM.?DD|AA.?MM.?DD|YYYY.?MM|AAAA.?MM|YY.?MM|AA.?MM|MM.?DD|DD.?MM.?YYYY|DD.?MM.?AAAA|DD.?MM.?YY|DD.?MM.?AA`,
	)
	matches = patt.FindAllString(query, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			frmtFinal := getDtFmrt(m)
			var procc string
			if dates, ok := dateRef.([]time.Time); ok {
				procc = regexp.MustCompile(patt.String()).ReplaceAllString(m, dates[0].Format(frmtFinal))
			} else if dt, ok := dateRef.(time.Time); ok {
				procc = regexp.MustCompile(patt.String()).ReplaceAllString(m, dt.Format(frmtFinal))
			}
			patt = regexp.MustCompile(regexp.QuoteMeta(m))
			query = patt.ReplaceAllString(query, procc)
		}
	}
	return query
}

func (app *application) setStrEnv(input string) string {
	re := regexp.MustCompile(`@ENV\.\w+`)
	matches := re.FindAllString(input, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			envVar := strings.TrimPrefix(match, "@ENV.")
			envValue := os.Getenv(envVar)
			if envValue != "" {
				input = strings.ReplaceAll(input, match, envValue)
			}
		}
	}
	return input
}

func (app *application) _extract(params map[string]interface{}, file string, _input map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, db_conf map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
	// FILE HANDLE
	//fileName := filepath.Base(file)
	fileExt := filepath.Ext(file)
	//fileNameNoExt := fileName[:len(fileName)-len(fileExt)]
	//fmt.Println(file, fileName, fileNameNoExt, fileExt)
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
	if _, ok := _input["database"]; ok {
		_database = _input["database"].(string)
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
	if fileExt == ".xlsx" {
		if !app.contains(_duck_conf["extensions"].([]interface{}), "excel") {
			_duck_conf["extensions"] = append(_duck_conf["extensions"].([]interface{}), "excel")
		}
	}
	if app.contains([]interface{}{"sqlite", "sqlite3"}, _driver) {
		if !app.contains(_duck_conf["extensions"].([]interface{}), "sqlite") {
			_duck_conf["extensions"] = append(_duck_conf["extensions"].([]interface{}), "sqlite")
		}
	}
	if app.contains([]interface{}{".duckdb", ".ddb"}, fileExt) {
		if !app.contains(_duck_conf["extensions"].([]interface{}), "sqlite") {
			_duck_conf["extensions"] = append(_duck_conf["extensions"].([]interface{}), "sqlite")
		}
	}
	//fmt.Println("extensions:", _duck_conf["extensions"])
	app.duckdb_start(db, _duck_conf, _driver, _database)
	// RUN EXTRACT HERE
	_db_base := filepath.Base(_database)
	_db_ext := filepath.Ext(_database)
	_db_name_no_ext := _db_base[:len(_db_base)-len(_db_ext)]
	//fmt.Println(_driver, _database, _db_base, _db_ext, _db_name_no_ext)
	destination_table := ""
	if _, ok := _input["destination_table"]; ok {
		destination_table = _input["destination_table"].(string)
	}
	// DATE REF
	var _date_ref interface{}
	if _, ok := _input["date_ref"]; ok {
		_date_ref = _input["date_ref"]
	} else if _, ok := _step["dates_refs"]; ok {
		_date_ref = _step["dates_refs"]
	}
	var date_ref []time.Time
	switch _date_ref.(type) {
	case string:
		_dt, _ := time.Parse("2006-01-02", _date_ref.(string))
		date_ref = append(date_ref, _dt)
	case []interface{}:
		for _, _dt := range _date_ref.([]interface{}) {
			_dt, _ := time.Parse("2006-01-02", _dt.(string))
			date_ref = append(date_ref, _dt)
		}
	default:
		// fmt.Println("default:", _type)
	}
	//fmt.Println(date_ref)
	// FILE REF
	file_ref := app.get_ref_from_file(file)
	//fmt.Println(file_ref)
	// CHECK DATE
	check_ref_date := false
	if _, ok := _input["check_ref_date"]; !ok {
	} else if app.contains([]interface{}{true, 1, "1", "true", "True", "TRUE"}, _input["check_ref_date"]) {
		check_ref_date = true
	}
	ref_date_field := ""
	if _, ok := _input["ref_date_field"]; !ok {
	} else if _, ok := _input["ref_date_field"].(string); ok {
		ref_date_field = _input["ref_date_field"].(string)
	}
	date_format_org := "YYYYMMDD"
	if _, ok := _input["date_format_org"]; !ok {
	} else if _, ok := _input["date_format_org"].(string); ok {
		date_format_org = _input["date_format_org"].(string)
	}
	//fmt.Println(check_ref_date, ref_date_field)
	if !check_ref_date || ref_date_field == "" {
	} else if date_ref[0].Format("20060102") != file_ref.Format("20060102") && app.contains([]interface{}{"file_ref", "FILEREF", "REF"}, ref_date_field) {
		app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
		msg, _ := app.i18n.T("dt-file-no-match-dt-form", map[string]interface{}{
			"date_ref": date_ref[0].Format("2006-01-02"),
			"file_ref": file_ref.Format("2006-01-02"),
		})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	} else if !app.contains([]interface{}{"file_ref", "FILEREF", "REF"}, ref_date_field) {
		if _, ok := _duck_conf["valid"]; !ok {
			_duck_conf["valid"] = []interface{}{}
		}
		if app.IsEmpty(_duck_conf["valid"]) {
			_duck_conf["valid"] = append(_duck_conf["valid"].([]interface{}), map[string]interface{}{
				"sql":  fmt.Sprintf(`SELECT * FROM main.\"<table>\" WHERE "%s"='{%s}' LIMIT 10`, ref_date_field, date_format_org),
				"rule": "throw_if_not_empty",
				"msg":  "The table (<table>) already has the data from the date YYYY-MM-DD!",
			})
			_duck_conf["valid"] = append(_duck_conf["valid"].([]interface{}), map[string]interface{}{
				"sql":  fmt.Sprintf(`SELECT * FROM \"<file>\" WHERE "%s"='{%s}' LIMIT 10`, ref_date_field, date_format_org),
				"rule": "throw_if_empty",
				"msg":  "The file (<file>) has no data from the date \"YYYY-MM-DD\"!",
			})
		}
	}
	// CHECK VALIDATION QUERIES
	_sql := fmt.Sprintf(`CREATE OR REPLACE TABLE %s."%s" AS SELECT * FROM '%s'`, _db_name_no_ext, destination_table, file)
	if app.contains([]interface{}{".xls", ".xlsx"}, fileExt) {
		_sql = fmt.Sprintf(`CREATE OR REPLACE TABLE %s."%s" AS SELECT * FROM ST_READ('%s')`, _db_name_no_ext, destination_table, file)
	}
	if app.contains([]interface{}{".duckdb", ".ddb"}, fileExt) && file != "" {
		_, err = db.ExecuteQuery(fmt.Sprintf(`ATTACH '%s' AS FILE`, file))
		if err != nil {
			fmt.Println("ATTACH:", fileExt, err, _sql)
			app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("Err Attachin: %s -> %s", file, err),
			}
		}
		_sql = fmt.Sprintf(`CREATE OR REPLACE TABLE %s."%s" AS SELECT * FROM FILE."%s"`, _db_name_no_ext, destination_table, destination_table)
	}
	// VALIDATIONS
	if _, ok := _duck_conf["valid"]; !ok {
	} else if app.IsEmpty(_duck_conf["valid"]) {
	} else {
		_valid := _duck_conf["valid"]
		switch _type := _valid.(type) {
		case []interface{}:
			for _, v := range _valid.([]interface{}) {
				_query := ""
				if _, ok := v.(map[string]interface{})["sql"]; ok {
					_query = v.(map[string]interface{})["sql"].(string)
				}
				if _, ok := v.(map[string]interface{})["query"]; ok {
					_query = v.(map[string]interface{})["query"].(string)
				}
				_query = app.set_sql_file_table("table", _query, destination_table)
				_query = app.set_sql_file_table("file", _query, file)
				_query = app.setQueryDate(_query, date_ref)
				_msg := ""
				if _, ok := v.(map[string]interface{})["msg"]; ok {
					_msg = v.(map[string]interface{})["msg"].(string)
				}
				_msg = app.set_sql_file_table("table", _msg, destination_table)
				_msg = app.set_sql_file_table("file", _msg, file)
				_msg = app.setQueryDate(_msg, date_ref)
				_rule := ""
				if _, ok := v.(map[string]interface{})["rule"]; ok {
					_rule = v.(map[string]interface{})["rule"].(string)
				}
				res, _, err := db.QueryMultiRows(_query, []interface{}{}...)
				if err != nil {
					fmt.Println("Err:", _query, err)
				} else {
					//fmt.Println((*res)[0])
					if len((*res)) > 0 && _rule == "throw_if_not_empty" {
						app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
						return map[string]interface{}{
							"success": false,
							"msg":     _msg,
						}
					} else if len((*res)) == 0 && _rule == "throw_if_empty" {
						app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
						return map[string]interface{}{
							"success": false,
							"msg":     _msg,
						}
					}
				}
				//fmt.Println("VALIDATION QUERY:", _query, _msg, _rule)
			}
		default:
			fmt.Println(_type)
		}
	}
	// EXECUTE THE MAIN QUERY
	if _, ok := _duck_conf["query"]; ok {
		_sql = _duck_conf["query"].(string)
	} else if _, ok := _duck_conf["sql"]; ok {
		_sql = _duck_conf["sql"].(string)
	}
	_sql = app.set_sql_file_table("table", _sql, destination_table)
	_sql = app.set_sql_file_table("file", _sql, file)
	_sql = app.setQueryDate(_sql, date_ref)
	_sql = app.setStrEnv(_sql)
	fmt.Println(_sql)
	// CHECK IF SQL MATCHS INSERT PATT
	//patt := `(?i)INSERT.+?INTO\s*["|\` + "`" + `|\[]?\w+["|\` + "`" + `|\]]?\.\s*["|\` + "`" + `|\[]?\w+["|\` + "`" + `|\]]?|INSERT.+?INTO\s*["|\` + "`" + `|\[]?\s*.\w+["|\` + "`" + `|\]]?`
	patt := `(?i)INSERT\s+INTO\s+(?:"?(\w+)"?\.)?"?(\w+)"?`
	re := regexp.MustCompile(patt)
	match := re.FindStringSubmatch(_sql)
	if len(match) > 0 {
		res := re.FindAllSubmatch([]byte(_sql), -1)
		// fmt.Println(match[0], res[0][1], res[0][2])
		_db_name := ""
		_tbl_name := ""
		if len(res) > 0 {
			_db_name = string(res[0][1])
			_tbl_name = string(res[0][2])
		}
		sql_table_chk := fmt.Sprintf(`SELECT * 
		FROM "information_schema"."tables" 
		WHERE "table_name" = '%s'`, _tbl_name)
		if _db_name != "" {
			if _db_name == "main" {
				_db_name = _db_name_no_ext
			}
			sql_table_chk = fmt.Sprintf(`SELECT * 
			FROM "information_schema"."tables" 
			WHERE "table_name" = '%s'
				AND LOWER("table_catalog") = LOWER('%s')`, _tbl_name, _db_name)
		}
		//fmt.Println(sql_table_chk, _db_name, _tbl_name)
		chk_tbl, _, err := db.QueryMultiRows(sql_table_chk, []interface{}{}...)
		_table_already_created := true
		if err != nil {
			fmt.Println("Err:", sql_table_chk, err)
		} else {
			if len((*chk_tbl)) > 0 {
			} else {
				_table_already_created = false
				if _, ok := _duck_conf["sql_create_if_not_exists"]; ok {
					_sql = _duck_conf["sql_create_if_not_exists"].(string)
				} else {
					patt = `(?i)INSERT\s+INTO\s+(?:"?(\w+)"?\.)?"?(\w+)"?`
					_replc_str := `CREATE TABLE "${1}"."${2}" AS`
					if _db_name == "" {
						_replc_str = `CREATE TABLE "${2}" AS`
					}
					_sql = regexp.MustCompile(patt).ReplaceAllString(_sql, _replc_str)
					_sql = regexp.MustCompile(`\bAS\s+AS\b`).ReplaceAllString(_sql, "AS")
					_sql = regexp.MustCompile(`BY.+NAME`).ReplaceAllString(_sql, "")
				}
				fmt.Println(_table_already_created, _sql)
			}
		}
		// CHECK IF IT MATCHES SELECT * FROM THE FILE / TABLE
		patt = `(?is)SELECT\s+\*\s+FROM\s+["']?([\w./-]+)["']?(?:\s*\.\s*["']?([\w./-]+)["'])?`
		patt = `(?is)SELECT\s+\*\s+FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?`
		patt = `(?is)SELECT\s+\*\s+FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?|\b([\w./-]+)\s*\(\s*["']?([\w./-]+)["']?(?:\s*,\s*[\w\s=]*)*\s*\)`
		patt = `(?is)SELECT\s+\*\s+FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?|SELECT\s+\*\s+FROM\s+([\w./-]+)\s*\(\s*["']?([\w./-]+)["']?(?:\s*,\s*[\w\s=]*)*\s*\)`
		patt = `(?is)SELECT\s+\*\s+FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?|SELECT\s+\*\s+FROM\s+([\w./-]+)\s*\(([^)]+)\)`
		re = regexp.MustCompile(patt)
		select_star := re.FindString(_sql)
		patt = `(?i)SELECT\s+\*\s+FROM\s+(["'\w]+(?:\.[\w]+)?)\s*(\((.*?)\))?`
		re = regexp.MustCompile(patt)
		matches := re.FindStringSubmatch(_sql)
		if len(matches) > 0 {
			select_star = matches[0]
		}
		//fmt.Println("SELECT STAR:", select_star)
		if _table_already_created && len(select_star) > 0 {
			_file_columns_list := []string{}
			_file_columns_types := map[string]string{}
			desc_select_star := fmt.Sprintf(`DESC %s`, select_star)
			desc_org, _, err := db.QueryMultiRows(desc_select_star, []interface{}{}...)
			if err != nil {
				fmt.Println("Err:", desc_select_star, err)
			} else {
				for _, row := range *desc_org {
					if _, ok := row["column_name"]; ok {
						_file_columns_list = append(_file_columns_list, row["column_name"].(string))
						_file_columns_types[row["column_name"].(string)] = row["column_type"].(string)
					}
				}
			}
			//fmt.Println(select_star, _file_columns_list, _file_columns_types)
			_db_columns_list := []string{}
			desc_dest_query := fmt.Sprintf(`DESC %s`, destination_table)
			desc_dest, _, err := db.QueryMultiRows(desc_dest_query, []interface{}{}...)
			if err != nil {
				fmt.Println("Err:", desc_dest_query, err)
			} else {
				for _, row := range *desc_dest {
					if _, ok := row["column_name"]; ok {
						_db_columns_list = append(_db_columns_list, row["column_name"].(string))
					}
				}
				for _, col := range _file_columns_list {
					if !app.contains(app.sliceStrs2SliceInterfaces(_db_columns_list), col) {
						_sql_add_new_field := fmt.Sprintf(`ALTER TABLE "%s" ADD "%s" %s NULL`, destination_table, col, _file_columns_types[col])
						_, err := db.ExecuteQuery(_sql_add_new_field, []interface{}{}...)
						if err != nil {
							fmt.Println("ERR ADDING NEW FIELDS:", _sql_add_new_field, err)
						}
					}
				}
			}
		}
	}
	_, err = db.ExecuteQuery(_sql)
	if err != nil {
		fmt.Println("EXTRACT:", err, _sql)
		app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("Err Importing: %s", err),
		}
	}
	patt = `(?is)SELECT.+FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?`
	patt = `(?is)SELECT.+?FROM\s+(?:["']?([\w./-]+)["']?\.)?["']?([\w./-]+)["']?.*`
	re = regexp.MustCompile(patt)
	_select := re.FindString(_sql)
	_sql = fmt.Sprintf(`SELECT COUNT(*) AS "n_rows" FROM (%s) AS "T"`, _select)
	fmt.Println(_sql)
	n_rows_res, _, err := db.QuerySingleRow(_sql, []interface{}{}...)
	if err != nil {
		fmt.Println("EXTRACT NROWS:", err, _sql)
		app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("Err Importing: %s", err),
		}
	}
	n_rows := (*n_rows_res)["n_rows"]
	// DETACH DBS TO THE DUCKDB IN MEM CONN
	app.duckdb_end(db, _duck_conf, _driver, _database, fileExt)
	//data := map[string]interface{}{}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"n_rows":  n_rows,
	}
}

func (app *application) _odbc_csv_duckdb(params map[string]interface{}, _input map[string]interface{}, _conf map[string]interface{}, _etlrb map[string]interface{}, _conf_etlrb map[string]interface{}, db_conf map[string]interface{}, _step map[string]interface{}) map[string]interface{} {
	odbc_conn := ""
	if _, ok := _conf["odbc_conn"]; ok {
		odbc_conn = _conf["odbc_conn"].(string)
	} else if _, ok := _conf["params"]; !ok {
	} else if _, ok := _conf["params"].(map[string]interface{})["odbc_conn"]; ok {
		odbc_conn = _conf["params"].(map[string]interface{})["odbc_conn"].(string)
	}
	if odbc_conn == "" {
		msg, _ := app.i18n.T("no-odbc-conn-str", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
		}
	}
	//fmt.Println(1, odbc_conn)
	odbc_conn = app.setStrEnv(odbc_conn)
	//fmt.Println(2, odbc_conn)
	query := ""
	if _, ok := _conf["query"]; ok {
		query = _conf["query"].(string)
	} else if _, ok := _conf["sql"]; !ok {
		query = _conf["sql"].(string)
	}
	if query == "" {
		msg, _ := app.i18n.T("no-odbc-query", map[string]interface{}{})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
		}
	}
	//fmt.Println(3, query)
	// DATE REF
	var _date_ref interface{}
	if _, ok := _input["date_ref"]; ok {
		_date_ref = _input["date_ref"]
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
	//fmt.Println("date_ref:", date_ref)
	query = app.setStrEnv(query)
	query = app.setQueryDate(query, date_ref)
	//fmt.Println(4, query)
	destination_table := ""
	if _, ok := _input["destination_table"]; ok {
		destination_table = _input["destination_table"].(string)
	}
	//fmt.Println(5, destination_table)
	//new_odbc, err := odbc.New(odbc_conn)
	new_odbc, err := etlx.NewODBC(odbc_conn)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("ODBC Conn: %s", err),
		}
	}
	defer new_odbc.Close()
	_csv_path := fmt.Sprintf(`%s/%s_YYYYMMDD.csv`, os.TempDir(), destination_table)
	_csv_path = app.setQueryDate(_csv_path, date_ref)
	//fmt.Println(_csv_path)
	_, err = new_odbc.Query2CSV(query, _csv_path)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("ODBC CSV: %s", err),
		}
	}
	return app._extract(params, _csv_path, _input, _conf, _etlrb, _conf_etlrb, db_conf, _step)
}
