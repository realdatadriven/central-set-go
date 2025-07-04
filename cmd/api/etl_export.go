package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/realdatadriven/etlx"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	// "github.com/google/uuid"
)

func (app *application) _export(params map[string]any, _item map[string]any, _conf map[string]any, _etlrb map[string]any, _conf_etlrb map[string]any, db_conf map[string]any, _step map[string]any) map[string]any {
	// IN MEMORY DUCKDB CONN
	db, err := etlx.NewDuckDB("")
	if err != nil {
		return map[string]any{
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
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := _etlrb["database"]; ok {
		_database = _etlrb["database"].(string)
		if _database != "" {
			_ext := filepath.Ext(_database)
			if app.contains([]any{".duckdb", ".ddb"}, _ext) {
				_driver = "duckdb"
			}
		}
	} else if _, ok := db_conf["dsn"]; ok {
		_database = db_conf["dsn"].(string)
	}
	// ATTACH DBS TO THE DUCKDB IN MEM CONN
	_duck_conf := map[string]any{}
	if _, ok := _conf["duckdb"]; ok {
		_duck_conf = _conf["duckdb"].(map[string]any)
	}
	if _, ok := _duck_conf["extensions"]; !ok {
		_duck_conf["extensions"] = []any{}
	}
	//fmt.Println("extensions:", _duck_conf["extensions"])
	app.duckdb_start(db, _duck_conf, _driver, _database)
	// DATE REF
	var _date_ref any
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
	case []any:
		for _, _dt := range _date_ref.([]any) {
			_dt, _ := time.Parse("2006-01-02", _dt.(string))
			date_ref = append(date_ref, _dt)
		}
	default:
		fmt.Println("default:", _type)
	}
	// GET THE OUTPUT FIELDS
	_fields_table := ""
	if _, ok := _step["detail_table"]; ok {
		_fields_table = _step["detail_table"].(string)
	} else {
		_fields_table = "etl_rb_exp_dtail"
	}
	etl_rbase_export_id := any(-1)
	if _, ok := _item["etl_rbase_export_id"]; ok {
		etl_rbase_export_id = _item["etl_rbase_export_id"]
	}
	params["data"].(map[string]any)["table"] = _fields_table
	params["data"].(map[string]any)["limit"] = any(-1.0)
	params["data"].(map[string]any)["offset"] = any(0.0)
	params["data"].(map[string]any)["filters"] = []any{}
	params["data"].(map[string]any)["filters"] = append(
		params["data"].(map[string]any)["filters"].([]any),
		map[string]any{
			"field": "etl_rbase_export_id",
			"cond":  "=",
			"value": etl_rbase_export_id,
		},
	)
	params["data"].(map[string]any)["order_by"] = []any{}
	params["data"].(map[string]any)["order_by"] = append(
		params["data"].(map[string]any)["order_by"].([]any),
		map[string]any{
			"field": "etl_rb_exp_dtail_id",
			"order": "ASC",
		},
	)
	_export_details := []map[string]any{}
	_aux_export_details := app.read(params)
	if _, ok := _aux_export_details["success"]; !ok {
		return _aux_export_details
	} else if _aux_export_details["success"].(bool) {
		_export_details = _aux_export_details["data"].([]map[string]any)
	}
	_exps := []map[string]any{}
	_logs_exports := []map[string]any{}
	for _, _dt := range date_ref {
		//fmt.Println(_dt)
		/*for _, _detail := range _export_details {
			fmt.Println(_dt, _detail)
		}*/
		file := app.setQueryDate(_item["attach_file_template"].(string), _dt)
		template := _item["attach_file_template"].(string)
		//fmt.Println(file, template)
		_aux := app._run_export(_dt, _item, _conf, _etlrb, _conf_etlrb, _export_details, file, template, db)
		_exps = append(_exps, _aux)
	}
	//fmt.Println(_exps, _logs_exports)
	_has_err := false
	for _, _exp := range _exps {
		if _, ok := _exp["success"]; !ok {
			_has_err = true
			continue
		} else if _, ok := _exp["success"].(bool); !ok {
			_has_err = true
			continue
		} else if !_exp["success"].(bool) {
			_has_err = true
			continue
		}
		_logs_exports = append(_logs_exports, _exp["aux"].(map[string]any))
	}
	if len(_logs_exports) > 0 {
		// APPEND LOGS OF EXPORTS
		// logs_exports to csv
		_csv_file := fmt.Sprintf(`%s/logs_exports.csv`, os.TempDir())
		cols := []string{}
		_, err := app.SliceToCSV(_logs_exports, cols, _csv_file)
		if err != nil {
			fmt.Println(err)
		}
		_input := map[string]any{
			"file":              filepath.Base(_csv_file),
			"etl_rbase_input":   "logs_exports",
			"database":          _item["database"],
			"save_only_temp":    true,
			"destination_table": "logs_exports",
			/*"etl_rbase_input_conf": map[string]any{
				"type": "file-duckdb",
				"duckdb": map[string]any{
					"sql":  "INSERT INTO \"<table>\" BY NAME SELECT * FROM '<file>'",
				},
			},*/
		}
		_conf = map[string]any{
			"type": "file-duckdb",
			"duckdb": map[string]any{
				"sql": "INSERT INTO \"<table>\" BY NAME SELECT * FROM '<file>'",
			},
		}
		app.duckdb_end(db, _duck_conf, _driver, _database, "")
		_extrt := app.ETLExtract(params, db_conf, _input, _conf, _etlrb, _conf_etlrb, _step)
		fmt.Println(_extrt)
	}
	if _has_err {

	}
	// DETACH DBS TO THE DUCKDB IN MEM CONN
	app.duckdb_end(db, _duck_conf, _driver, _database, "")
	//data := map[string]any{}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    _exps,
	}
}

func (app *application) _run_export(date_ref time.Time, _export map[string]any, _conf map[string]any, _etlrb map[string]any, _conf_etlrb map[string]any, _export_details []map[string]any, file string, template string, db *etlx.DuckDB) map[string]any {
	//tmp := "tmp"
	ext := filepath.Ext(file)
	if app.contains([]any{".xls", ".xlsx", ".XLS", ".XLSX"}, ext) {
		_, err := db.ExecuteQuery(`INSTALL Excel`)
		if err != nil {
			_err_msg := fmt.Sprintf(`Err: %s`, err)
			fmt.Println(_err_msg)
		}
		_, err = db.ExecuteQuery(`LOAD Excel`)
		if err != nil {
			_err_msg := fmt.Sprintf(`Err: %s`, err)
			fmt.Println(_err_msg)
		}
	}
	_exps := []map[string]any{}
	for _, _detail := range _export_details {
		//fmt.Println(date_ref, _detail)
		_dtail_conf := map[string]any{}
		if _, ok := _detail["etl_rb_exp_dtail_conf"].(map[string]any); ok {
			_dtail_conf = _detail["etl_rb_exp_dtail_conf"].(map[string]any)
		} else if _, ok := _detail["conf"].(map[string]any); ok {
			_dtail_conf = _detail["conf"].(map[string]any)
		} else if _, ok := _detail["_conf"].(map[string]any); ok {
			_dtail_conf = _detail["_conf"].(map[string]any)
		}
		_sql := ""
		if _, ok := _detail["sql_export_query"]; ok {
			_sql = _detail["sql_export_query"].(string)
		} else if _, ok := _detail["sql"]; ok {
			_sql = _detail["sql"].(string)
		} else if _, ok := _detail["query"]; ok {
			_sql = _detail["query"].(string)
		}
		_sql = app.setQueryDate(_sql, date_ref)
		patt := `COPY.+?\(.+\).+TO+.\'.+\'`
		match := regexp.MustCompile(patt).Match([]byte(strings.ReplaceAll(_sql, "\n", " ")))
		if !match {
			_sql = fmt.Sprintf(`COPY (%s) TO '<fname>'`, _sql)
		}
		details_to_tmp := false
		if _, ok := _dtail_conf["details_to_tmp"]; ok {
			if app.contains([]any{true, 1, "True", "TRUE", "T", "1"}, _dtail_conf["details_to_tmp"]) {
				details_to_tmp = true
			}
		}
		each_details_on_its_own_file := false
		if _, ok := _etlrb["each_details_on_its_own_file"]; ok {
			if app.contains([]any{true, 1, "True", "TRUE", "T", "1"}, _etlrb["each_details_on_its_own_file"]) {
				each_details_on_its_own_file = true
			}
		}
		_export_full_path := fmt.Sprintf(`%s/%s`, app.config.upload_path, file)
		if each_details_on_its_own_file {
			//f'{_detail.get("etl_rb_exp_dtail")}_YYYYMMDD{_format}'
			_file := fmt.Sprintf(`%s_YYYYMMDD%s`, _detail["etl_rb_exp_dtail"].(string), ext)
			_export_full_path = fmt.Sprintf(`%s/%s`, app.config.upload_path, _file)
			if details_to_tmp {
				_export_full_path = fmt.Sprintf(`%s/%s`, os.TempDir(), _file)
			}
		}
		_export_full_path = app.setQueryDate(_export_full_path, date_ref)
		//fmt.Println(_export_full_path)
		_sql = app.set_sql_file_table("file", _sql, _export_full_path)
		//fmt.Println(_sql)
		_, err := db.ExecuteQuery(_sql)
		if err != nil {
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf(`Err: %s`, err),
			}
		}
		_aux := map[string]any{
			"name":  _detail["etl_rb_exp_dtail"],
			"ref":   app.setQueryDate("{YYYY-MM-DD}", date_ref),
			"stamp": time.Now(),
			"file":  filepath.Base(_export_full_path),
		}
		msg, _ := app.i18n.T("success", map[string]any{})
		_exps = append(_exps, map[string]any{
			"success": true,
			"msg":     msg,
			"fname":   _export_full_path,
			"aux":     _aux,
		})
		//fmt.Println(_sql)
		// ZIP ALL THE DETAILS WITH THE EMPLATE
		if each_details_on_its_own_file && len(_export_details) > 1 {
			ext = filepath.Ext(template)
			basename := app.setQueryDate(template[:len(template)-len(ext)], date_ref)
			zped_path := fmt.Sprintf(`%s/tmp/%s.zip`, app.config.upload_path, basename)
			files := []string{zped_path}
			_auxs := []map[string]any{}
			for _, exp := range _exps {
				files = append(files, exp["fname"].(string))
				if _, ok := exp["aux"]; ok {
					_auxs = append(_auxs, exp["aux"].(map[string]any))
				}
			}
			err := app.createZipFile(zped_path, files)
			if err != nil {
				return map[string]any{
					"success": false,
					"msg":     fmt.Sprintf(`Err Zipping: %s`, err),
				}
			}
			msg, _ := app.i18n.T("success", map[string]any{})
			return map[string]any{
				"success": true,
				"msg":     msg,
				"fname":   fmt.Sprintf(`tmp/%s.zip`, basename),
				"aux":     _auxs,
			}
		} else {
			_aux2 := _exps[0]
			_aux2["fname"] = fmt.Sprintf(`tmp/%s`, filepath.Base(_aux2["fname"].(string)))
			return _aux2
		}
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		//"n_rows":  n_rows,
	}
}

// Function to add a file to the zip archive
func (app *application) addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	// Create a file in the zip archive based on the original file's name
	w, err := zipWriter.Create(filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %w", err)
	}
	// Copy the file's content into the zip file
	if _, err = io.Copy(w, file); err != nil {
		return fmt.Errorf("failed to write file to zip: %w", err)
	}
	return nil
}

// Main function to create a zip file and add multiple files
func (app *application) createZipFile(zipFileName string, filesToAdd []string) error {
	// Create the zip file
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()
	// Initialize zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	// Add files to the zip archive
	for _, file := range filesToAdd {
		if err := app.addFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("failed to add file to zip: %w", err)
		}
	}
	return nil
}

func isUTF8(s string) bool {
	return utf8.ValidString(s)
}

func convertToUTF8(isoStr string) (string, error) {
	if isUTF8(isoStr) {
		return isoStr, nil
	}
	reader := strings.NewReader(isoStr)
	transformer := charmap.ISO8859_1.NewDecoder()
	utf8Bytes, err := io.ReadAll(transform.NewReader(reader, transformer))
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

func hasDecimalPlace(v any) (bool, error) {
	// Try to cast v to float64
	floatVal, ok := v.(float64)
	if !ok {
		return false, fmt.Errorf("value is not a float64, it is %v", reflect.TypeOf(v))
	}

	// Check if the float has a decimal part
	if floatVal != float64(int(floatVal)) {
		return true, nil
	}
	return false, nil
}

func (app *application) SliceToCSV(data []map[string]any, cols []string, csv_path string) (bool, error) {
	csvFile, err := os.Create(csv_path)
	if err != nil {
		return false, fmt.Errorf("error creating CSV file: %w", err)
	}
	defer csvFile.Close()
	// CSV
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()
	// Get column names
	columns := []string{}
	if len(cols) > 0 {
		columns = cols
	} else {
		for col := range data[0] {
			columns = append(columns, col)
		}
	}
	// Write column names to CSV header
	csvWriter.Write(columns)
	for _, row := range data {
		var rowData []string
		//for _, value := range row {
		for _, col := range columns {
			value := row[col]
			//rowData = append(rowData, fmt.Sprintf("%v", value))
			switch v := value.(type) {
			case nil:
				// Format integer types
				rowData = append(rowData, "")
			case int, int8, int16, int32, int64:
				// Format integer types
				rowData = append(rowData, fmt.Sprintf("%d", v))
			case float64, float32:
				//fmt.Println(col, v)
				// Format large numbers without scientific notation
				hasDec, err := hasDecimalPlace(v)
				if err != nil {
					fmt.Println(err)
					rowData = append(rowData, fmt.Sprintf("%v", value))
				} else if hasDec {
					rowData = append(rowData, fmt.Sprintf("%f", v))
				} else {
					rowData = append(rowData, fmt.Sprintf("%.f", v))
				}
			case []byte:
				// Convert byte slice (UTF-8 data) to a string
				utf8Str, err := convertToUTF8(string(v))
				if err != nil {
					fmt.Println("Failed to convert to UTF-8:", v, err)
				}
				rowData = append(rowData, strings.TrimSpace(string(utf8Str)))
			default:
				// Default formatting for other types
				rowData = append(rowData, fmt.Sprintf("%v", value))
			}
		}
		csvWriter.Write(rowData)
	}
	return true, nil
}
