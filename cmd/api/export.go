package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/realdatadriven/etlx"
	"github.com/xuri/excelize/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (app *application) generateRandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func (app *application) export_read(params map[string]any) map[string]any {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	_conf := map[string]any{}
	if _, ok := _data["_conf"]; ok {
		_conf = _data["_conf"].(map[string]any)
	}
	if _, ok := _conf["records"]; !ok {
	} else if _, ok := _conf["records"].(string); !ok {
	} else if _conf["records"].(string) == "all_records" {
		params["data"].(map[string]any)["limit"] = any(-1.0)
	}
	_read_data := []map[string]any{}
	_aux_read_data := app.read(params)
	if _, ok := _aux_read_data["success"]; !ok {
		return _aux_read_data
	} else if _aux_read_data["success"].(bool) {
		_read_data = _aux_read_data["data"].([]map[string]any)
	}
	_csv_file := fmt.Sprintf(`%s/%s.csv`, os.TempDir(), app.generateRandomString(40))
	cols := []string{}
	if _, ok := _aux_read_data["cols"]; ok {
		cols = app.sliceInterfaces2SliceStrs(_aux_read_data["cols"].([]any))
	}
	if _, ok := _conf["display_fields"]; !ok {
	} else if _, ok := _conf["display_fields"].(string); !ok {
	} else if _conf["display_fields"].(string) == "interface_fields" {
		if _, ok := _data["_fields"]; !ok {
		} else if _, ok := _data["_fields"].([]any); ok {
			_cols := app.filterInterface(_data["_fields"].([]any), func(r map[string]any) bool {
				return r["display"].(bool)
			})
			if len(_cols) > 0 {
				cols = []string{}
				for _, r := range _cols {
					cols = append(cols, r["name"].(string))
				}
			}
		}
	}
	_, err := app.SliceToCSV(_read_data, cols, _csv_file)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(_csv_file)
	fname := "export"
	if _, ok := _conf["name"]; ok {
		fname = _conf["name"].(string)
	} else if _, ok := _conf["file"]; ok {
		fname = _conf["file"].(string)
	} else if _, ok := _conf["fname"]; ok {
		fname = _conf["fname"].(string)
	} else if _, ok := _data["name"]; ok {
		fname = _data["name"].(string)
	} else if _, ok := _data["file"]; ok {
		fname = _data["file"].(string)
	} else if _, ok := _data["fname"]; ok {
		fname = _data["fname"].(string)
	}
	_format := ".csv"
	if _, ok := _conf["format"]; ok {
		_format = _conf["format"].(string)
	} else if _, ok := _data["format"]; ok {
		_format = _data["format"].(string)
	}
	compress := false
	if _, ok := _conf["compress"]; ok {
		compress = _conf["compress"].(bool)
	} else if _, ok := _data["compress"]; ok {
		compress = _data["compress"].(bool)
	}
	compress_format := ""
	if _, ok := _conf["compress_format"]; ok && compress {
		compress_format = _conf["compress_format"].(string)
	} else if _, ok := _data["compress_format"]; ok && compress {
		compress_format = _data["compress_format"].(string)
	}
	ext := filepath.Ext(fname)
	if ext != "" {
		_format = ext
	} else {
		_dot := "."
		if strings.HasPrefix(_format, ".") {
			_dot = ""
		}
		fname = fmt.Sprintf(`%s%s%s`, fname, _dot, _format)
	}
	_path := fmt.Sprintf(`%s/tmp/%s%s`, app.config.upload_path, fname, compress_format)
	_sql := "SELECT * FROM AUX"
	patt := `COPY.+?\(.+\).+TO+.\'.+\'`
	match := regexp.MustCompile(patt).Match([]byte(strings.ReplaceAll(_sql, "\n", " ")))
	if !match {
		if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, _format) {
			_sql = fmt.Sprintf(`COPY (%s) TO '<fname>' WITH (FORMAT xlsx, HEADER true)`, _sql)
		} else {
			_sql = fmt.Sprintf(`COPY (%s) TO '<fname>'`, _sql)
		}
	}
	_sql = app.set_sql_file_table("file", _sql, _path)
	// DB
	_driver := ""
	_database := ""
	if _, ok := _data["database"]; ok {
		_database = _data["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := _data["db"]; ok {
		_database = _data["db"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	}
	//fmt.Println(_path, _database, _sql)
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("DDB Conn: %s", err),
		}
	}
	defer db.Close()
	// ATTACH DBS TO THE DUCKDB IN MEM CONN
	_duck_conf := map[string]any{}
	if _, ok := _duck_conf["extensions"]; !ok {
		_duck_conf["extensions"] = []any{}
	}
	if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, _format) {
		if !app.contains(_duck_conf["extensions"].([]any), "excel") {
			_duck_conf["extensions"] = append(_duck_conf["extensions"].([]any), "excel")
		}
	}
	app.duckdb_start(db, _duck_conf, _driver, "")
	// ADD THE DATA EXPORTED TO DDB
	_sql_aux := fmt.Sprintf(`CREATE TABLE AUX AS SELECT * FROM '%s'`, _csv_file)
	_, err = db.ExecuteQuery(_sql_aux)
	if err != nil {
		fmt.Println("CREATING AUX:", err, _sql_aux)
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("Err: %s", err),
		}
	}
	_, err = db.ExecuteQuery(_sql)
	if err != nil {
		fmt.Println("EXPORT:", err, _sql)
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("Err: %s", err),
		}
	}
	app.duckdb_end(db, _duck_conf, _driver, "", "")
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"fname":   fmt.Sprintf(`tmp/%s`, filepath.Base(_path)),
	}
}

func (app *application) export_query(params map[string]any) map[string]any {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	_conf := map[string]any{}
	if _, ok := _data["_conf"]; ok {
		_conf = _data["_conf"].(map[string]any)
	}
	fname := "export"
	if _, ok := _conf["name"]; ok {
		fname = _conf["name"].(string)
	} else if _, ok := _conf["file"]; ok {
		fname = _conf["file"].(string)
	} else if _, ok := _conf["fname"]; ok {
		fname = _conf["fname"].(string)
	} else if _, ok := _data["name"]; ok {
		fname = _data["name"].(string)
	} else if _, ok := _data["file"]; ok {
		fname = _data["file"].(string)
	} else if _, ok := _data["fname"]; ok {
		fname = _data["fname"].(string)
	}
	_format := ".csv"
	if _, ok := _conf["format"]; ok {
		_format = _conf["format"].(string)
	} else if _, ok := _data["format"]; ok {
		_format = _data["format"].(string)
	}
	compress := false
	if _, ok := _conf["compress"]; ok {
		compress = _conf["compress"].(bool)
	} else if _, ok := _data["compress"]; ok {
		compress = _data["compress"].(bool)
	}
	compress_format := ""
	if _, ok := _conf["compress_format"]; ok && compress {
		compress_format = _conf["compress_format"].(string)
	} else if _, ok := _data["compress_format"]; ok && compress {
		compress_format = _data["compress_format"].(string)
	}
	ext := filepath.Ext(fname)
	if ext != "" {
		_format = ext
	} else {
		_dot := "."
		if strings.HasPrefix(_format, ".") {
			_dot = ""
		}
		fname = fmt.Sprintf(`%s%s%s`, fname, _dot, _format)
	}
	_path := fmt.Sprintf(`%s/tmp/%s%s`, app.config.upload_path, fname, compress_format)
	_sql := ""
	if _, ok := _data["sql"]; ok {
		_sql = _data["sql"].(string)
	} else if _, ok := _data["query"]; ok {
		_sql = _data["query"].(string)
	}
	//fmt.Println(1, _data)
	patt := `COPY.+?\(.+\).+TO+.\'.+\'`
	match := regexp.MustCompile(patt).Match([]byte(strings.ReplaceAll(_sql, "\n", " ")))
	if !match {
		if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, _format) {
			_sql = fmt.Sprintf(`COPY (%s) TO '<fname>' WITH (FORMAT ([]any), "excel"), DRIVER 'xlsx')`, _sql)
		} else {
			_sql = fmt.Sprintf(`COPY (%s) TO '<fname>'`, _sql)
		}
	}
	_sql = app.set_sql_file_table("file", _sql, _path)
	// DB
	_driver := ""
	_database := ""
	if _, ok := _data["database"]; ok {
		_database = _data["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := _data["db"]; ok {
		_database = _data["db"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	}
	//fmt.Println(dl)
	//fmt.Println(_path, _database, _sql)
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("DDB Conn: %s", err),
		}
	}
	defer db.Close()
	dl := etlx.NewDuckLakeParser().Parse(_database)
	if dl.IsDuckLake {
		attach := fmt.Sprintf(`ATTACH %s`, dl.DSN)
		if dl.HasAttach {
			attach = _database
		}
		//fmt.Println(1, attach)
		_, err = db.ExecuteQuery(attach)
		if err != nil {
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("Err Exporting: %s", err),
			}
		}
		if dl.DuckLakeName != "" {
			db.ExecuteQuery(fmt.Sprintf(`USE %s`, dl.DuckLakeName))
		}
		//fmt.Println(2, fmt.Sprintf(`USE %s`, dl.DuckLakeName))
		_, err = db.ExecuteQuery(_sql)
		if err != nil {
			if dl.DuckLakeName != "" {
				db.ExecuteQuery("USE memory")
				db.ExecuteQuery(fmt.Sprintf(`DETACH %s`, dl.DuckLakeName))
			}
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("Err Exporting: %s", err),
			}
		}
		if dl.DuckLakeName != "" {
			db.ExecuteQuery("USE memory")
			db.ExecuteQuery(fmt.Sprintf(`DETACH %s`, dl.DuckLakeName))
		}
	} else {
		// ATTACH DBS TO THE DUCKDB IN MEM CONN
		_duck_conf := map[string]any{}
		if _, ok := _duck_conf["extensions"]; !ok {
			_duck_conf["extensions"] = []any{}
		}
		if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, _format) {
			if !app.contains(_duck_conf["extensions"].([]any), "excel") {
				_duck_conf["extensions"] = append(_duck_conf["extensions"].([]any), "excel")
			}
		}
		app.duckdb_start(db, _duck_conf, _driver, _database)
		_, err = db.ExecuteQuery(_sql)
		if err != nil {
			fmt.Println("EXPORT:", err, _sql)
			app.duckdb_end(db, _duck_conf, _driver, _database, "")
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("Err Emporting: %s", err),
			}
		}
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"fname":   fmt.Sprintf(`tmp/%s`, filepath.Base(_path)),
	}
}

func (app *application) dump_file_2_object(params map[string]any) map[string]any {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	//fmt.Println(_csv_file)
	fname := "export"
	if _, ok := _data["name"]; ok {
		fname = _data["name"].(string)
	} else if _, ok := _data["file"]; ok {
		fname = _data["file"].(string)
	} else if _, ok := _data["fname"]; ok {
		fname = _data["fname"].(string)
	}
	_format := ".csv"
	if _, ok := _data["format"]; ok {
		_format = _data["format"].(string)
	}
	ext := filepath.Ext(fname)
	if ext != "" {
		_format = ext
	} else {
		_dot := "."
		if strings.HasPrefix(_format, ".") {
			_dot = ""
		}
		fname = fmt.Sprintf(`%s%s%s`, fname, _dot, _format)
	}
	is_tmp := true
	//fmt.Println(_data)
	if _, ok := _data["tmp"]; !ok {
	} else if _, ok := _data["tmp"].(bool); ok {
		is_tmp = _data["tmp"].(bool)
	}
	_path := fmt.Sprintf(`%s/%s`, app.config.upload_path, fname)
	if is_tmp {
		_path = fmt.Sprintf(`%s/%s`, os.TempDir(), fname)
	}
	_sql := fmt.Sprintf(`SELECT * FROM '%s'`, _path)
	if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, ext) {
		_sql = fmt.Sprintf(`SELECT * FROM READ_XLSX('%s', HEADER = TRUE)`, _path)
	}
	//fmt.Println(_path, _sql)
	// DB
	_driver := ""
	_database := ""
	if _, ok := _data["database"]; ok {
		_database = _data["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := _data["db"]; ok {
		_database = _data["db"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	}
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("DDB Conn: %s", err),
		}
	}
	defer db.Close()
	// ATTACH DBS TO THE DUCKDB IN MEM CONN
	_duck_conf := map[string]any{}
	if _, ok := _duck_conf["extensions"]; !ok {
		_duck_conf["extensions"] = []any{}
	}
	if app.contains([]any{".xlsx", "xlsx", ".XLSX", "XLSX"}, _format) {
		if !app.contains(_duck_conf["extensions"].([]any), "excel") {
			_duck_conf["extensions"] = append(_duck_conf["extensions"].([]any), "excel")
		}
	}
	app.duckdb_start(db, _duck_conf, _driver, "")
	data, _, err := db.QueryMultiRows(_sql, []any{}...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("Err: %s", err),
		}
	}
	app.duckdb_end(db, _duck_conf, _driver, "", "")
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    *data,
	}
}

func (app *application) dump_2_html(params map[string]any) map[string]any {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	}
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	tmpl := ""
	if _, ok := _data["html"]; ok {
		tmpl = _data["html"].(string)
	} else if _, ok := _data["htmlstr"]; ok {
		tmpl = _data["htmlstr"].(string)
	}
	// Parse the template
	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		return map[string]any{
			"success": true,
			"msg":     fmt.Sprintf("Error parsing template: %v", err),
		}
	}
	// Execute the template to a string
	var result bytes.Buffer
	if err := t.Execute(&result, _data); err != nil {
		return map[string]any{
			"success": true,
			"msg":     fmt.Sprintf("Error executing template: %v", err),
		}
	}
	// Output the processed template as a string
	output := result.String()
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"html":    output,
	}
}

func GenerateExcelExport(header map[string]any, details []map[string]any, db *sql.DB) (string, error) {
	templateFile, ok := header["attach_file_template"].(string)
	if !ok || templateFile == "" {
		return "", fmt.Errorf("attach_file_template missing or invalid")
	}
	// Check for supported spreadsheet extensions
	ext := filepath.Ext(templateFile)
	if ext != ".xlsx" && ext != ".xls" && ext != ".xlsm" {
		return "", fmt.Errorf("unsupported template file extension: %s", ext)
	}
	// Open or create a new workbook
	var file *excelize.File
	var err error
	if fileExists(templateFile) {
		file, err = excelize.OpenFile(templateFile)
		if err != nil {
			return "", fmt.Errorf("failed to open template file: %w", err)
		}
	} else {
		file = excelize.NewFile()
	}
	// Iterate over each detail entry to handle data insertion
	for _, detail := range details {
		// Skip inactive entries
		if active, _ := detail["active"].(bool); !active {
			continue
		}
		destSheetName := detail["dest_sheet_name"].(string)
		destTableName := detail["dest_table_name"].(string)
		sqlQuery := detail["sql_export_query"].(string)
		// Check or create the destination sheet
		sheetIndex, err := file.GetSheetIndex(destSheetName)
		if sheetIndex == -1 || err != nil {
			file.NewSheet(destSheetName)
		} else {
			file.DeleteSheet(destSheetName)
			file.NewSheet(destSheetName)
		}
		// Fetch data from the database using the provided SQL query
		rows, err := db.Query(sqlQuery)
		if err != nil {
			return "", fmt.Errorf("error executing query for detail ID %v: %w", detail["etl_rb_exp_dtail_id"], err)
		}
		defer rows.Close()
		// Fetch column names
		columns, err := rows.Columns()
		if err != nil {
			return "", fmt.Errorf("failed to get columns: %w", err)
		}
		// Write column headers
		for colIdx, colName := range columns {
			cell, err := excelize.JoinCellName(string('A'+colIdx), 1)
			if err != nil {
				return "", fmt.Errorf("failed to set columns: %w", err)
			}
			file.SetCellValue(destSheetName, cell, colName)
		}
		// Write data rows
		rowIdx := 2
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(values))
			for i := range pointers {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				return "", fmt.Errorf("failed to scan row data: %w", err)
			}
			for colIdx, value := range values {
				cell, err := excelize.JoinCellName(string('A'+colIdx), rowIdx)
				if err != nil {
					return "", fmt.Errorf("failed to set columns: %w", err)
				}
				file.SetCellValue(destSheetName, cell, value)
			}
			rowIdx++
		}
		// Create Excel table if `dest_table_name` is specified
		if destTableName != "" {
			cell, err := excelize.JoinCellName(string('A'+len(columns)-1), rowIdx-1)
			if err != nil {
				return "", fmt.Errorf("failed to set columns: %w", err)
			}
			tableRange := fmt.Sprintf("A1:%s", cell)
			err = file.AddTable(destSheetName, &excelize.Table{
				Name:            destTableName,
				Range:           tableRange,
				StyleName:       "TableStyleMedium9",
				ShowFirstColumn: false,
				ShowLastColumn:  false,
				//ShowRowStripes:    true,
				ShowColumnStripes: false,
			})
			if err != nil {
				return "", fmt.Errorf("failed to create table %s on sheet %s: %w", destTableName, destSheetName, err)
			}
		}
	}
	/*/ Enable pivot table updates on open if template has pivot tables
	for _, sheet := range file.GetSheetList() {
		pivotTables, err := file.GetPivotTable(sheet)
		if err == nil {
			for _, pivot := range pivotTables {
				err = file.SetPivotTable(sheet, pivot.Name, excelize.PivotTableOptions{UpdateOnOpen: true})
				if err != nil {
					log.Printf("failed to set pivot update on open for %s: %v", pivot.Name, err)
				}
			}
		}
	}*/
	// Save the file to a temporary location
	outputFile := filepath.Join(os.TempDir(), fmt.Sprintf("export_%d%s", header["etl_rbase_export_id"], ext))
	err = file.SaveAs(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	return outputFile, nil
}

// Utility function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
