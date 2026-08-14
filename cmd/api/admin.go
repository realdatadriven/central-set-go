package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"

	"github.com/jmoiron/sqlx"
)

func (app *application) apps(params Dict) Dict {
	//fmt.Println("APPS:", params)
	user_id := app.toInt(params["user"].(Dict)["user_id"])
	role_id := app.toInt(params["user"].(Dict)["role_id"])
	//fmt.Println(user_id, role_id)
	query := `SELECT DISTINCT user_role.role_id
	FROM user_role
	JOIN role ON user_role.role_id = role.role_id
	WHERE user_role.user_id = $1
		AND user_role.excluded = FALSE
		AND role.excluded = FALSE`
	var queryParams []any
	queryParams = append(queryParams, user_id)
	result, _, err := app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		//fmt.Println(1, query, fmt.Sprintf("%s", err))
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	roles := []any{}
	roles = append(roles, role_id)
	for _, row := range *result {
		roles = append(roles, app.toInt(row["role_id"]))
	}
	query = `SELECT *
	FROM app
	WHERE app.excluded = FALSE`
	// fmt.Println("ROLES: ", roles)
	queryParams = []any{}
	if !app.contains(roles, 1) {
		query = `SELECT app.*
		FROM app
		JOIN role_app ON role_app.app_id = app.app_id
		WHERE role_app.role_id IN (?)
			AND role_app.access = TRUE
			AND role_app.excluded = FALSE
			AND app.excluded = FALSE`
		//fmt.Println(app.joinSlice(roles, ","))
		queryParams = append(queryParams, roles)
	}
	query, args, err := sqlx.In(query, queryParams...)
	result, _, err = app.db.QueryMultiRows(query, args...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
		"data":    *result,
	}
}

func (app *application) menu(params Dict) Dict {
	//fmt.Println(params)
	user_id := app.toInt(params["user"].(Dict)["user_id"])
	role_id := app.toInt(params["user"].(Dict)["role_id"])
	var app_id int
	if _, ok := params["app"].(Dict)["app_id"]; ok {
		app_id = app.toInt(params["app"].(Dict)["app_id"])
	}
	//fmt.Println(user_id, role_id)
	query := `SELECT DISTINCT user_role.role_id
	FROM user_role
	JOIN role ON user_role.role_id = role.role_id
	WHERE user_role.user_id = $1
		AND user_role.excluded = FALSE
		AND role.excluded = FALSE`
	var queryParams []any
	queryParams = append(queryParams, user_id)
	result, _, err := app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	roles := []any{}
	roles = append(roles, role_id)
	for _, row := range *result {
		roles = append(roles, app.toInt(row["role_id"]))
	}
	// MENU
	query = `SELECT *
	FROM menu
	WHERE app_id = ?
		AND excluded = FALSE
		AND active = TRUE
	ORDER BY menu_order ASC, menu_id ASC`
	//fmt.Println("ROLES: ", roles)
	queryParams = []any{app_id}
	if !app.contains(roles, 1) {
		query = `SELECT DISTINCT menu.*
		FROM menu
		JOIN role_app_menu ON (
			role_app_menu.menu_id = menu.menu_id 
			AND role_app_menu.app_id = menu.app_id
		)
		WHERE menu.app_id = ?
			AND role_app_menu.role_id IN (?)
			AND role_app_menu.access = TRUE
			AND role_app_menu.excluded = FALSE
			AND menu.excluded = FALSE
			AND menu.active = TRUE
			AND (role_app_menu.menu_id, role_app_menu.app_id, role_app_menu.updated_at) IN (
				SELECT menu_id, app_id, MAX(updated_at)
				FROM role_app_menu
				WHERE access = True
					AND role_id IN (?)
					AND excluded = FALSE
				GROUP BY menu_id, app_id
			)
		ORDER BY menu.menu_order ASC, menu.menu_id ASC`
		//fmt.Println(app.joinSlice(roles, ","))
		queryParams = append(queryParams, roles, roles)
	}
	query, args, err := sqlx.In(query, queryParams...)
	//fmt.Println("query:", 1, query, args, len(args))
	if err != nil {
		println("Error geting the table query:", err)
	}
	_menu, _, err := app.db.QueryMultiRows(query, args...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	// MENU TABLES
	query = `SELECT *
	FROM menu_table
	WHERE app_id = ?
		AND excluded = FALSE
	ORDER BY menu_table_id ASC`
	queryParams = []any{app_id}
	if !app.contains(roles, 1) {
		query = `SELECT menu_table.*
		FROM menu_table
		JOIN role_app_menu_table ON (
			role_app_menu_table.menu_id = menu_table.menu_id 
			AND role_app_menu_table.table_id = menu_table.table_id 
			AND role_app_menu_table.app_id = menu_table.app_id
		)
		WHERE menu_table.app_id = ?
			AND role_app_menu_table.role_id IN (?)
			AND (
				role_app_menu_table."read" = TRUE
				OR role_app_menu_table."create" = TRUE
			)
			AND role_app_menu_table.excluded = FALSE
			AND menu_table.excluded = FALSE
			AND (role_app_menu_table.table_id, role_app_menu_table.menu_id, role_app_menu_table.app_id, role_app_menu_table.updated_at) IN (
				SELECT table_id, menu_id, app_id, MAX(updated_at)
				FROM role_app_menu_table
				WHERE (role_app_menu_table."read" = TRUE OR role_app_menu_table."create" = TRUE)
					AND role_id IN (?)
					AND excluded = FALSE
				GROUP BY table_id, menu_id, app_id
			)
		ORDER BY menu_table_id ASC`
		queryParams = append(queryParams, roles, roles)
	}
	query, args, err = sqlx.In(query, queryParams...)
	//fmt.Println("query:", 1, query, args, len(args))
	if err != nil {
		println("Error geting the table query:", err)
	}
	_menu_table, _, err := app.db.QueryMultiRows(query, args...)
	if err != nil {

		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Error getting menus: %s", err),
		}
	}
	// TABLES
	_tables := app.tables(params, []any{})
	_table_by_id := map[int64]any{}
	if _, ok := _tables["table_by_id"]; ok {
		_table_by_id = _tables["table_by_id"].(map[int64]any)
		//fmt.Println(_table_by_id)
	}
	if _, ok := _tables["data"]; ok {
		_tables = _tables["data"].(Dict)
	}
	// MENUS
	menus := []Dict{}
	for _, mn := range *_menu {
		_aux := mn
		_aux["children"] = []Dict{}
		for _, mnt := range *_menu_table {
			if _, ok := mnt["menu_id"]; !ok {
			} else if _, ok := mn["menu_id"]; !ok {
			} else if app.toInt(mnt["menu_id"]) == app.toInt(mn["menu_id"]) {
				_mnt := mnt
				//fmt.Println(1, _table_by_id[mnt["table_id"].(int64)].(Dict))
				if _, ok := _table_by_id[mnt["table_id"].(int64)].(Dict); ok {
					_mnt["table"] = _table_by_id[mnt["table_id"].(int64)].(Dict)["table"].(string)
					//fmt.Println(2, _table_by_id[mnt["table_id"].(int64)].(Dict))
				}
				_mnt["menu"] = mn["menu"]
				_aux["children"] = append(_aux["children"].([]Dict), _mnt)
			}
		}
		menus = append(menus, _aux)
	}
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
		"data": Dict{
			"menu":   menus,
			"tables": _tables,
		},
	}
}

func (app *application) ParseConnection(conn string) (string, string, error) {
	parts := strings.SplitN(conn, ":", 2)
	if len(parts) < 2 {
		return "", conn, nil
	}
	dl := etlx.NewDuckLakeParser().Parse(conn)
	if dl.IsDuckLake {
		return "ducklake", conn, nil
	}
	return parts[0], parts[1], nil
}

// ExtractDBName extracts the database name (dbname) from various connection string formats.
func (app *application) ExtractURLDBName(dsn string) (string, error) {
	// First, try parsing as a URL (handles URL-style connection strings)
	if strings.Contains(dsn, "://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("invalid URL format: %w", err)
		}
		// In URL-style DSNs, the path usually starts with "/", so trim it
		dbname := strings.TrimPrefix(u.Path, "/")
		if dbname != "" {
			return dbname, nil
		}
	}
	// Fallback: try parsing key-value style (e.g. user=... dbname=... port=...)
	re := regexp.MustCompile(`(?i)\bdbname\s*=\s*([^\s]+)`)
	match := re.FindStringSubmatch(dsn)
	if len(match) >= 2 {
		return match[1], nil
	}

	return "", fmt.Errorf("could not find dbname in dsn")
}

func (app *application) GetDBNameFromParams(params Dict) (string, string, error) {
	// fmt.Println(params)
	var _database any
	if !app.IsEmpty(params["db"]) {
		_database = params["db"]
	} else if !app.IsEmpty(params["database"]) {
		_database = params["database"]
	} else if _, ok := params["data"].(Dict); !ok {
	} else if !app.IsEmpty(params["data"].(Dict)["db"]) {
		_database = params["data"].(Dict)["db"]
	} else if !app.IsEmpty(params["data"].(Dict)["database"]) {
		_database = params["data"].(Dict)["database"]
	} else if _, ok := params["app"].(Dict); !ok {
	} else if !app.IsEmpty(params["app"].(Dict)["db"]) {
		_database = params["app"].(Dict)["db"]
	} else if !app.IsEmpty(params["app"].(Dict)["db"]) {
		_database = params["app"].(Dict)["db"]
	}
	//_not_embed_dbs := []any{"postgres", "postgresql", "pg", "pgql", "mysql"}
	_embed_dbs := []any{"sqlite", "sqlite3", "duckdb", "ducklake"}
	_embed_dbs_ext := []any{".db", ".duckdb", ".ddb", ".sqlite", ".ducklake"}
	//fmt.Println(1, _database)
	switch _type := _database.(type) {
	case nil:
		//fmt.Println("IS NIL:", _database, _type)
		_db := ""
		fileName := filepath.Base(app.config.db.dsn)
		fileExt := filepath.Ext(app.config.db.dsn)
		if app.fileExists(app.config.db.dsn) || (fileName != "" && fileName != "." && fileExt != "") {
			_db = fileName[:len(fileName)-len(fileExt)]
		} else {
			dbname, err := app.ExtractURLDBName(app.config.db.dsn)
			if err == nil && dbname != "" {
				_db = dbname
			}
		}
		return app.config.db.dsn, _db, nil
	case string:
		_dsn := _database.(string)
		_driver, dsn, err := app.ParseConnection(_dsn)
		if _driver == "ducklake" {
			return dsn, _dsn, nil
		}
		//fmt.Println(_dsn, _driver, dsn)
		dirName := filepath.Dir(dsn)
		fileName := filepath.Base(dsn)
		fileExt := filepath.Ext(dsn)
		if err != nil {
			dsn = _dsn
		}
		if _driver == "" {
			if app.contains([]any{".duckdb", ".ddb"}, fileExt) {
				_driver = "duckdb"
			} else if app.contains([]any{".db", ".sqlite"}, fileExt) {
				_driver = "sqlite3"
			} else {
				_driver = app.config.db.driverName
			}
		}
		if app.contains(_embed_dbs, _driver) || app.contains(_embed_dbs_ext, fileExt) {
			embed_dbs_dir := "database"
			if os.Getenv("DB_EMBEDED_DIR") != "" {
				embed_dbs_dir = os.Getenv("DB_EMBEDED_DIR")
			}
			//fmt.Println("dirName: ", dirName, "fileName: ", fileName, "fileExt: ", fileExt)
			if filepath.Base(dsn) == fileName || dirName == "" {
				dsn = fmt.Sprintf("%s:%s/%s", _driver, embed_dbs_dir, fileName)
			}
			if fileExt == "" {
				_embed_dbs = []any{"sqlite", "sqlite3"}
				if _driver == "duckdb" {
					dsn = fmt.Sprintf("%s:%s/%s.duckdb", _driver, embed_dbs_dir, fileName)
				} else if app.contains(_embed_dbs, _driver) {
					dsn = fmt.Sprintf("%s:%s/%s.db", _driver, embed_dbs_dir, fileName)
				}
			}
			fileName = filepath.Base(_database.(string))
			fileExt = filepath.Ext(_database.(string))
			if app.fileExists(_database.(string)) || (fileName != "" && fileName != "." && fileExt != "") {
				_database = fileName[:len(fileName)-len(fileExt)]
			}
		} else {
			new_dsn, err := etlx.ReplaceDBName(app.config.db.dsn, dsn)
			if err != nil {
				fmt.Println("Errr getting the DSN for ", dsn)
			}
			dsn = fmt.Sprintf("%s:%s", _driver, new_dsn)
			if strings.HasPrefix(new_dsn, fmt.Sprintf("%s:", _driver)) {
				dsn = new_dsn
			}
			dbname, err := app.ExtractURLDBName(dsn)
			if err == nil && dbname != "" {
				_database = dbname
			}
			//fmt.Println(1, "ExtractURLDBName:", dbname, _database.(string), err)
		}
		//fmt.Println(2, "ExtractURLDBName:", _database.(string))
		return dsn, _database.(string), nil
	case []any:
		fmt.Println("IS []any:", _database, _type)
		return "", "", errors.New("database conf is of type []any")
	default:
		return _database.(string), _database.(string), nil
	}
}

func (app *application) toInt(v any) int {
	switch val := v.(type) {
	case float32:
		return int(val)
	case float64:
		return int(val)
	case int64:
		return int(val)
	case int32:
		return int(val)
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0 // Handle invalid strings as needed
		}
		return i
	case int:
		return val
	default:
		return 0 // or handle error
	}
}

func (app *application) toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float32:
		return bool(val == 1)
	case float64:
		return bool(val == 1)
	case int64:
		return bool(val == 1)
	case int32:
		return bool(val == 1)
	case int:
		return bool(val == 1)
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false // Handle invalid strings as needed
		}
		return b
	default:
		return false // or handle error
	}
}

func (app *application) tables(params Dict, tables []any) Dict {
	//fmt.Println(1, params)
	var user_id int
	if _, ok := params["user"].(Dict)["user_id"]; ok {
		user_id = app.toInt(params["user"].(Dict)["user_id"])
	}
	var app_id int
	if _, ok := params["app"].(Dict)["app_id"]; ok {
		app_id = app.toInt(params["app"].(Dict)["app_id"])
	}
	// DATABASE
	_extra_conf := Dict{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	dsn, _database, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	allTables := false
	if app.IsEmpty(tables) {
		tables = []any{}
		if !app.IsEmpty(params["data"].(Dict)["table"]) {
			value := params["data"].(Dict)["table"]
			switch value.(type) {
			case nil:
				// pass
			case string:
				tables = append(tables, params["data"].(Dict)["table"].(string))
			case []any:
				_tables := params["data"].(Dict)["table"].([]any)
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			case map[any]any:
				// pass
			default:
				tables = append(tables, params["data"].(Dict)["table"].(string))
			}
		} else if !app.IsEmpty(params["data"].(Dict)["tables"]) {
			value := params["data"].(Dict)["tables"]
			switch value.(type) {
			case string:
				tables = append(tables, params["data"].(Dict)["tables"].(string))
			case []any:
				_tables := params["data"].(Dict)["tables"].([]any)
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			default:
				tables = append(tables, params["data"].(Dict)["table"].(string))
			}
		}
		//fmt.Println("TABLES:", tables)
		if app.IsEmpty(tables) {
			// fmt.Println("GET ALL TABLES!")
			result, _, err := newDB.AllTables(params, _extra_conf)
			if err != nil {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("%s", err),
				}
			}
			for _, row := range *result {
				//fmt.Println(row)
				if _, ok := row["name"]; !ok {
				} else if _, ok := row["name"].(string); ok {
					tables = append(tables, string(row["name"].(string)))
				} else if _, ok := row["name"].([]byte); ok {
					tables = append(tables, string(row["name"].([]byte)))
				}
			}
			allTables = true
		}
	}
	//fmt.Println(2, dsn, _database, tables, allTables)
	data := Dict{}
	table_by_id := map[int64]any{}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
			"tables":  tables,
		}
	} else {
		// GET THE TABLES DATA IN table
		query := `SELECT * FROM "table" WHERE db = ? AND "table" IN (?) AND excluded = FALSE`
		queryParams := []any{_database}
		// fmt.Println("DATABASE:", _database)
		if allTables {
			query = `SELECT * FROM "table" WHERE db = ? AND excluded = FALSE`
		} else {
			queryParams = append(queryParams, tables)
		}
		//queryParams = append(queryParams, app.joinSlice(tables, "','"))
		query, args, err := sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query: ", err)
		}
		//fmt.Println(query, args, queryParams)
		_table, _, err := app.db.QueryMultiRows(query, args...)
		if err != nil {
			fmt.Println("TABLES: ", query, args, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		//fmt.Println(_table)
		if allTables {
			tables_in_table := []any{}
			for _, row := range *_table {
				tables_in_table = append(tables_in_table, row["table"].(string))
			}
			// fmt.Println(tables_in_table)
			results := []Dict{}
			for _, table := range tables {
				if !app.contains(tables_in_table, table) {
					//fmt.Println("ADD TABLE:", table)
					results = append(results, Dict{
						"table":        table,
						"table_desc":   table,
						"db":           _database,
						"requires_rla": false,
						"user_id":      user_id,
						"created_at":   time.Now(),
						"updated_at":   time.Now(),
						"excluded":     false,
					})
				}
			}
			if len(results) > 0 {
				//fmt.Println(results[0])
				var keys []any
				//var prms []any
				i := 0
				for key := range results[0] {
					i++
					keys = append(keys, key)
					//prms = append(prms, fmt.Sprintf("$%d", i))
				}
				// CHECK IF DUCKDB USE SOME OTHER WAY
				cols := app.joinSlice(keys, `", "`)
				vals := app.joinSlice(keys, `, :`)
				/*if driver == "duckdb" {
					vals = app.joinSlice(prms, `,`)
				} else {
					//vals = fmt.Sprintf(":%s", vals)
					vals = app.joinSlice(prms, `,`)
				}*/
				query := fmt.Sprintf(`INSERT INTO "table" ("%s") VALUES (:%s)`, cols, vals)
				/*_, err := app.db.ExecuteNamedQuery(query, results)
				if err != nil {
					fmt.Println("Error inserting table:", err)
				}*/
				//fmt.Println(query)
				for _, row := range results {
					_, err := app.db.ExecuteNamedQuery(query, row)
					/*values := []any{}
					for _, value := range row {
						values = append(values, value)
					}
					println(values)
					_, err := newDB.ExecuteQuery(query, values...)*/
					if err != nil {
						fmt.Println("Error inserting table:", err)
					}
				}
			}
		}
		// table comments / translations translate_table
		query = `SELECT * FROM translate_table WHERE db = ? AND lang = ? AND "table" IN (?) AND excluded = FALSE`
		queryParams = []any{_database, lang}
		if allTables {
			query = `SELECT * FROM translate_table WHERE db = ? AND lang = ? AND excluded = FALSE`
		} else {
			queryParams = append(queryParams, tables)
		}
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		results, _, err := app.db.QueryMultiRows(query, args...)
		if err != nil {
			fmt.Println("TABLES TRANSL:", query, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		translate_table := Dict{}
		for _, row := range *results {
			translate_table[row["table"].(string)] = row
		}
		//fmt.Println(translate_table)
		// fields comments / translations translate_table_field
		query = `SELECT * FROM translate_table_field WHERE db = ? AND lang = ? AND "table" IN (?) AND excluded = FALSE`
		queryParams = []any{_database, lang}
		if allTables {
			query = `SELECT * FROM translate_table_field WHERE db = ? AND lang = ? AND excluded = FALSE`
		} else {
			queryParams = append(queryParams, tables)
		}
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		results, _, err = app.db.QueryMultiRows(query, args...)
		if err != nil {
			fmt.Println("TARNSL FIELDS:", query, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		translate_table_field := Dict{}
		for _, row := range *results {
			if _, ok := translate_table_field[row["table"].(string)]; !ok {
				translate_table_field[row["table"].(string)] = Dict{}
			}
			/*if _, ok := translate_table_field[row["table"].(string)].(Dict)["fields"]; !ok {
				translate_table_field[row["table"].(string)].(Dict)["fields"] = Dict{}
			}
			translate_table_field[row["table"].(string)].(Dict)["fields"].(Dict)[row["field"].(string)] = row*/
			translate_table_field[row["table"].(string)].(Dict)[row["field"].(string)] = row
		}
		// fmt.Println(translate_table_field)
		// GET THE TABLES DATA IN table_schema
		query = `select ts.* 
		from table_schema ts
		left join "table" t on ts.db = t.db and ts."table" = t."table" and t.excluded = false
		where ts.db = ? 
			and ts."table" in (?) 
			and ts.excluded = false 
		order by t.table_id, ts.field_order`
		query = `SELECT * FROM table_schema WHERE db = ? AND "table" IN (?) AND excluded = FALSE order by field_order`
		queryParams = []any{_database}
		if allTables {
			query = `select ts.* 
			from table_schema ts
			left join "table" t on ts.db = t.db and ts."table" = t."table" and t.excluded = false
			where ts.db = ?
				and ts.excluded = false 
			order by t.table_id, ts.field_order`
			query = `SELECT * FROM table_schema WHERE db = ? AND excluded = FALSE order by field_order`
		} else {
			queryParams = append(queryParams, tables)
		}
		//queryParams = append(queryParams, app.joinSlice(tables, "','"))
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		//fmt.Println(allTables, query, args, queryParams, _database)
		table_schema := Dict{}
		_table_schema, _, err := app.db.QueryMultiRows(query, args...)
		if err != nil {
			fmt.Println("TABLE SCHEMA:", query, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		//fmt.Println(3, *_table_schema)
		// POPULATE table_schema WITH THOSE WHO ARE NOT IN table_schema
		if allTables {
			tables_not_in_schema := []any{}
			tables_in_schema := []any{}
			for _, row := range *_table_schema {
				tables_in_schema = append(tables_in_schema, row["table"].(string))
			}
			//fmt.Println("tables_in_schema:", tables_in_schema)
			for _, table := range tables {
				if !app.contains(tables_in_schema, table) {
					tables_not_in_schema = append(tables_not_in_schema, table)
					//fmt.Printf("Index: %d, Name: %s\n", _, table)
					// PUT IT
					res, _, err := newDB.TableSchema(params, table.(string), _database, _extra_conf)
					if err != nil {
						fmt.Printf("%s\n", err)
					} else {
						if len(*res) > 0 {
							results := *res
							//fmt.Println(results[0])
							var keys []any
							// Iterate over the map and collect the keys
							for key := range results[0] {
								keys = append(keys, key)
							}
							cols := app.joinSlice(keys, `", "`)
							vals := app.joinSlice(keys, `, :`)
							// Loop through the slice of maps and insert each record
							_ins_query := fmt.Sprintf(`INSERT INTO table_schema ("%s") VALUES (:%s)`, cols, vals)
							//fmt.Println(_ins_query)
							for _, row := range results {
								if table, ok := row["table"].(string); ok {
									if strings.HasPrefix(table, "sqlite_") {
										continue
									}
								}
								_, err := app.db.ExecuteNamedQuery(_ins_query, row)
								if err != nil {
									fmt.Println("Error inserting table_schema:", err)
								}
							}
						}
					}
				}
			}
			if len(tables_not_in_schema) > 0 {
				_table_schema, _, err = app.db.QueryMultiRows(query, args...)
				if err != nil {
					fmt.Println("TABLE SCHEMA CREATED:", query, err)
					return Dict{
						"success": false,
						"msg":     fmt.Sprintf("%s", err),
					}
				}
			}
			//fmt.Println("tables_not_in_schema:", tables_not_in_schema)
		}
		table_fields := Dict{}
		fk_tables_added := []any{}
		for _, _row := range *_table_schema {
			if _, ok := table_schema[_row["table"].(string)]; !ok {
				table_schema[_row["table"].(string)] = Dict{}
			}
			if _, ok := table_fields[_row["table"].(string)]; !ok {
				table_fields[_row["table"].(string)] = []any{}
			}
			//_row := row
			/*if _, ok := table_schema[row["table"].(string)].(Dict)["fields"]; !ok {
				table_schema[row["table"].(string)].(Dict)["fields"] = Dict{}
			}
			table_schema[row["table"].(string)].(Dict)["fields"].(Dict)[row["field"].(string)] = row*/
			comment := _row["comment"]
			if _, ok := translate_table_field[_row["table"].(string)]; !ok {
			} else if _, ok := translate_table_field[_row["table"].(string)].(Dict)[_row["field"].(string)]; !ok {
			} else if _, ok := translate_table_field[_row["table"].(string)].(Dict)[_row["field"].(string)].(Dict)["field_transl_desc"]; ok {
				comment = translate_table_field[_row["table"].(string)].(Dict)[_row["field"].(string)].(Dict)["field_transl_desc"]
			}
			_row["comment"] = comment
			_row["name"] = _row["field"]
			if _, ok := _row["fk"]; !ok {
			} else if app.contains([]any{1, true, "true", "True", "TRUE", "T", "1"}, _row["fk"]) || app.toBool(_row["fk"]) {
				//fmt.Println(_row["field"], _row["table"], _row["referred_table"], _row["referred_column"])
				referred_columns_desc := ""
				if _, ok := table_fields[_row["referred_table"].(string)].([]any); ok {
					if len(table_fields[_row["referred_table"].(string)].([]any)) > 1 {
						referred_columns_desc = table_fields[_row["referred_table"].(string)].([]any)[1].(string)
					}
				}
				fk_tables_added = append(fk_tables_added, Dict{"table": _row["table"], "referred_table": _row["referred_table"]})
				acorr := app.filterAny(fk_tables_added, func(r any) bool {
					return r.(Dict)["table"].(string) == _row["table"].(string) && r.(Dict)["referred_table"].(string) == _row["referred_table"].(string)
				})
				referred_column := _row["referred_table"]
				referred_columns_desc_org := referred_columns_desc
				if len(acorr) > 1 {
					//referred_column = fmt.Sprintf("%s%d", referred_column, len(acorr))
					referred_columns_desc = fmt.Sprintf("%s%d", referred_columns_desc, len(acorr))
				}
				_row["ref"] = Dict{
					"referred_table":            _row["referred_table"],
					"referred_column":           referred_column,
					"referred_columns_desc_org": referred_columns_desc_org,
					"referred_columns_desc":     referred_columns_desc,
				}
				// fmt.Println("REF:", _row["field"], _row["referred_table"].(string), _row["ref"])
			}
			table_schema[_row["table"].(string)].(Dict)[_row["field"].(string)] = _row
			table_fields[_row["table"].(string)] = append(table_fields[_row["table"].(string)].([]any), _row["field"])
			//_row = nil
		}
		// table form customizations custom_form
		query = `SELECT * 
		FROM custom_form
		WHERE db = ?
			AND (user_id = ? OR user_id = 1)
			AND app_id = ?
			AND "table" IN (?) 
			AND excluded = FALSE
		ORDER BY user_id DESC, custom_form_id DESC`
		queryParams = []any{_database, user_id, app_id}
		if allTables {
			query = `SELECT * 
			FROM custom_form
			WHERE db = ?
				AND (user_id = ? OR user_id = 1)
				AND app_id = ?
				AND excluded = FALSE
			ORDER BY user_id DESC, custom_form_id DESC`
		} else {
			queryParams = append(queryParams, tables)
		}
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		results, _, err = app.db.QueryMultiRows(query, args...)
		// fmt.Println("custom_form:", queryParams, results)
		if err != nil {
			fmt.Println("custom_form:", query, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		custom_form := Dict{}
		for _, row := range *results {
			// fmt.Println("custom_form:", row["table"].(string))
			custom_form[row["table"].(string)] = row
		}
		// table table customizations custom_table
		query = `SELECT * 
		FROM custom_table
		WHERE db = ?
			AND (user_id = ? OR user_id = 1)
			AND app_id = ?
			AND "table" IN (?)
			AND excluded = FALSE
		ORDER BY user_id DESC, custom_table_id DESC`
		queryParams = []any{_database, user_id, app_id}
		if allTables {
			query = `SELECT * 
			FROM custom_table
			WHERE db = ?
				AND (user_id = ? OR user_id = 1)
				AND app_id = ?
				AND excluded = FALSE
			ORDER BY user_id DESC, custom_table_id DESC`
		} else {
			queryParams = append(queryParams, tables)
		}
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		results, _, err = app.db.QueryMultiRows(query, args...)
		if err != nil {
			//fmt.Println("custom_table:", query, err)
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		custom_table := Dict{}
		for _, row := range *results {
			custom_table[row["table"].(string)] = row
		}
		// get crud actions for tables
		crud_actions := app.getTableCrudActions(app.db, _database, tables)
		if _, ok := crud_actions["data"].(Dict); ok {
			crud_actions = crud_actions["data"].(Dict)
		} else {
			fmt.Println("VALIDATIONS ERR:", crud_actions["msg"])
		}
		// get crud validations for tables
		crud_validations := app.getTableCrudValidations(app.db, _database, tables)
		if _, ok := crud_validations["data"].(Dict); ok {
			crud_validations = crud_validations["data"].(Dict)
		} else {
			fmt.Println("VALIDATIONS ERR:", crud_validations["msg"])
		}
		// return
		for _, row := range *_table {
			comment := row["table_desc"]
			if _, ok := translate_table[row["table"].(string)]; ok {
				comment = translate_table[row["table"].(string)].(Dict)["table_transl_desc"]
			}
			var pk string
			if _, ok := table_schema[row["table"].(string)]; ok {
				for key, value := range table_schema[row["table"].(string)].(Dict) {
					if properties, ok := value.(Dict); ok {
						// Check if the "pk" field exists and is true
						if _pk, found := properties["pk"]; found && _pk == true {
							pk = key
							break
						}
					}
				}
			}
			table_by_id[row["table_id"].(int64)] = row
			data[row["table"].(string)] = Dict{
				"table_id":              row["table_id"],
				"table":                 row["table"],
				"comment":               comment,
				"database":              row["db"],
				"_table":                row,
				"fields":                table_schema[row["table"].(string)],
				"custom_table":          custom_table[row["table"].(string)],
				"custom_form":           custom_form[row["table"].(string)],
				"translate_table":       translate_table[row["table"].(string)],
				"translate_table_field": translate_table_field[row["table"].(string)],
				"crud_actions":          crud_actions[row["table"].(string)],
				"validations":           crud_validations[row["table"].(string)],
				"pk":                    pk,
				"fields_order":          table_fields[row["table"].(string)],
			}
		}
	}
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success":     true,
		"msg":         msg,
		"data":        data,
		"table_by_id": table_by_id,
	}
}

func (app *application) getTableCrudActions(dbCon etlx.DBInterface, database string, tables []any) Dict {
	query := `SELECT * FROM crud_action WHERE db = ? AND "table" IN (?) AND excluded = FALSE`
	queryParams := []any{database, tables}
	query, args, err := sqlx.In(query, queryParams...)
	if err != nil {
		println("Error geting the table query: ", err)
	}
	//fmt.Println(query, args, queryParams)
	res, _, err := app.db.QueryMultiRows(query, args...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	data := Dict{}
	for _, row := range *res {
		if _, ok := data[row["table"].(string)]; !ok {
			data[row["table"].(string)] = []any{}
		}
		data[row["table"].(string)] = append(data[row["table"].(string)].([]any), row)
	}
	// fmt.Println("getTableCrudActions:", data)
	return Dict{
		"success": true,
		"msg":     "success",
		"data":    data,
	}
}

func (app *application) getTableCrudValidations(dbCon etlx.DBInterface, database string, tables []any) Dict {
	query := `SELECT * FROM validation WHERE db = ? AND "table" IN (?) AND excluded = FALSE`
	queryParams := []any{database, tables}
	query, args, err := sqlx.In(query, queryParams...)
	if err != nil {
		println("Error geting the table query: ", err)
	}
	//fmt.Println(query, args, queryParams)
	res, _, err := app.db.QueryMultiRows(query, args...)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	data := Dict{}
	for _, row := range *res {
		if _, ok := data[row["table"].(string)]; !ok {
			data[row["table"].(string)] = []any{}
		}
		data[row["table"].(string)] = append(data[row["table"].(string)].([]any), row)
	}
	//fmt.Println("getTableCrudValidations:", data)
	return Dict{
		"success": true,
		"msg":     "success",
		"data":    data,
	}
}

// Generates CREATE TABLE SQL statements with comments, adapting to SQL dialects
func generateCreateTableSQL(driver, tableName, tableComment string, fields []any) string {
	var schema strings.Builder
	schema.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	// Collect foreign keys and comments
	var foreignKeys []string
	var columnComments []string

	// Generate column definitions based on the driver
	for _, _field := range fields {
		field, ok := _field.(Dict)
		if !ok {
			continue
		}
		name := field["name"].(string)
		columnType := getColumnType(driver, field)

		// Primary key, autoincrement, nullable, and unique adjustments
		primaryKey := getPrimaryKey(driver, field)
		autoincrement := getAutoIncrement(driver, field)
		nullable := getNullable(driver, field)
		unique := getUnique(driver, field)

		// Handle default values
		defaultValue := getDefaultValue(driver, field)

		// Handle foreign keys
		if fk, ok := field["foreign_key"].(string); ok {
			foreignKeys = append(foreignKeys, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s", name, fk))
		}

		// Collect comments for columns
		if cmt, ok := field["comment"].(string); ok {
			columnComments = append(columnComments, getColumnComment(driver, tableName, name, cmt))
		}

		// Build the column definition string
		columnDef := fmt.Sprintf("    %s %s%s%s%s%s%s", name, columnType, primaryKey, autoincrement, nullable, unique, defaultValue)
		schema.WriteString(columnDef + ",\n")
	}

	// Add foreign key constraints
	for _, fk := range foreignKeys {
		schema.WriteString("    " + fk + ",\n")
	}

	// Trim the trailing comma and add closing parenthesis
	schemaStr := strings.TrimRight(schema.String(), ",\n") + "\n);\n"

	// Add table comment and column comments if supported by the driver
	if driver == "postgres" || driver == "mysql" {
		if tableComment != "" {
			schemaStr += getTableComment(driver, tableName, tableComment)
		}
		for _, colComment := range columnComments {
			schemaStr += colComment + "\n"
		}
	}

	return schemaStr
}

// Returns the appropriate SQL column type based on driver and field type
func getColumnType(driver string, field Dict) string {
	columnType := field["type"].(string)
	if nchar, ok := field["nchar"].(int); ok {
		columnType += fmt.Sprintf("(%d)", nchar)
	}
	// Map SQL types per dialect
	switch driver {
	case "postgres":
		if columnType == "INTEGER" && field["autoincrement"] == true {
			return "SERIAL "
		}
	case "mysql":
		if columnType == "INTEGER" && field["autoincrement"] == true {
			return "INT AUTO_INCREMENT"
		}
	case "sqlserver", "mssql":
		if columnType == "INTEGER" && field["autoincrement"] == true {
			return "INT IDENTITY(1,1)"
		}
	}
	return columnType
}

// Primary key syntax adjustments
func getPrimaryKey(driver string, field Dict) string {
	if pk, ok := field["primary_key"].(bool); ok && pk {
		if driver == "mysql" || driver == "sqlserver" || driver == "mssql" {
			return " PRIMARY KEY"
		}
	}
	return ""
}

// Autoincrement syntax adjustments per driver
func getAutoIncrement(driver string, field Dict) string {
	if field["autoincrement"] == true {
		if driver == "sqlite3" {
			return " AUTOINCREMENT"
		}
	}
	return ""
}

// Nullable syntax adjustments per driver
func getNullable(driver string, field Dict) string {
	if nullable, ok := field["nullable"].(bool); ok && !nullable {
		return " NOT NULL"
	}
	return ""
}

// Unique constraint syntax adjustments
func getUnique(driver string, field Dict) string {
	if unique, ok := field["unique"].(bool); ok && unique {
		return " UNIQUE"
	}
	return ""
}

// Default value handling based on driver
func getDefaultValue(driver string, field Dict) string {
	if defaultVal, ok := field["default"]; ok {
		switch v := defaultVal.(type) {
		case bool:
			return fmt.Sprintf(" DEFAULT %t", v)
		case string:
			return fmt.Sprintf(" DEFAULT '%s'", v)
		case int, float64:
			return fmt.Sprintf(" DEFAULT %v", v)
		}
	}
	return ""
}

// Generate column comment if supported by the driver
func getColumnComment(driver, tableName, columnName, comment string) string {
	switch driver {
	case "postgres", "mysql":
		return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", tableName, columnName, comment)
	}
	return ""
}

// Generate table comment if supported by the driver
func getTableComment(driver, tableName, comment string) string {
	switch driver {
	case "postgres", "mysql":
		return fmt.Sprintf("COMMENT ON TABLE %s IS '%s';\n", tableName, comment)
	}
	return ""
}

func generateModelYAML(tableName, tableComment string, fields []any) string {
	var schema strings.Builder
	schema.WriteString(fmt.Sprintf("table: %s\n", tableName))
	if tableComment != "" {
		schema.WriteString(fmt.Sprintf("comment: %s\n", tableComment))
	}
	schema.WriteString("columns:\n")
	for i, _field := range fields {
		field, ok := _field.(Dict)
		if !ok {
			continue
		}
		name := field["name"].(string)
		fmt.Println(i, name)
		var parts []string
		// type
		if t, ok := field["type"].(string); ok {
			if nchar, ok := field["nchar"].(int); ok {
				parts = append(parts, fmt.Sprintf("type: %s(%d)", strings.ToLower("VARCHAR"), nchar))
			} else {
				parts = append(parts, fmt.Sprintf("type: %s", strings.ToLower(t)))
			}
		}
		// primary key
		if pk, ok := field["primary_key"].(bool); ok && pk {
			parts = append(parts, "pk: true")
		}
		// autoincrement
		if ai, ok := field["autoincrement"].(bool); ok && ai {
			parts = append(parts, "autoincrement: true")
		}
		// nullable
		if nullable, ok := field["nullable"].(bool); ok && !nullable {
			parts = append(parts, "nullable: false")
		}
		// unique
		if unique, ok := field["unique"].(bool); ok && unique {
			parts = append(parts, "unique: true")
		}
		// default
		if def, ok := field["default"]; ok {
			switch v := def.(type) {
			case nil:
				// pass
			case string:
				parts = append(parts, fmt.Sprintf("default: \"%s\"", v))
			default:
				parts = append(parts, fmt.Sprintf("default: %v", v))
			}
		}
		// foreign key
		if fk, ok := field["foreign_key"].(string); ok {
			parts = append(parts, fmt.Sprintf("fk: \"%s\"", fk))
		}
		// comment
		if cmt, ok := field["comment"].(string); ok {
			parts = append(parts, fmt.Sprintf("comment: \"%s\"", cmt))
		}
		schema.WriteString(fmt.Sprintf("  %s: { %s }\n", name, strings.Join(parts, ", ")))
	}
	return schema.String()
}

func (app *application) save_table_schema(params Dict) Dict {
	//fmt.Println(params)
	//user_id := app.toInt(params["user"].(Dict)["user_id"])
	//role_id := app.toInt(params["user"].(Dict)["role_id"])
	//var app_id int
	//if _, ok := params["app"].(Dict)["app_id"]; ok {
	//	app_id = intapp.toInt(params["app"].(Dict)["app_id"])
	//}
	// DATABASE
	//fmt.Println(lang)
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	_data := Dict{}
	if _, ok := params["data"]; !ok {
		msg, _ := app.i18n.T("no_data", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := params["data"].(Dict); ok {
		_data = params["data"].(Dict)
	}
	table_metadata := Dict{}
	if _, ok := _data["table_metadata"]; !ok {
		msg, _ := app.i18n.T("no_table_metadata", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := _data["table_metadata"].(Dict); !ok {
		msg, _ := app.i18n.T("no_table_metadata", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	table_metadata = _data["table_metadata"].(Dict)
	name := ""
	if _, ok := table_metadata["name"]; !ok {
		msg, _ := app.i18n.T("no_table_name", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := table_metadata["name"].(string); !ok {
		msg, _ := app.i18n.T("no_table_name", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	name = table_metadata["name"].(string)
	comment := ""
	if _, ok := table_metadata["comment"]; !ok {
		msg, _ := app.i18n.T("no_table_comment", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := table_metadata["comment"].(string); !ok {
		msg, _ := app.i18n.T("no_table_comment", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	comment = table_metadata["comment"].(string)
	// fmt.Println("COLUMNS:", _data["fields"])
	fields := []any{}
	if _, ok := _data["fields"]; !ok {
		msg, _ := app.i18n.T("no_fields", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := _data["fields"].([]any); !ok {
		msg, _ := app.i18n.T("no_fields", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	fields = _data["fields"].([]any)
	if len(fields) < 2 {
		msg, _ := app.i18n.T("table_must_have_2_or_more_fields", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	core_tables := app.sliceStrs2SliceInterfaces(strings.Split(app.config.core_tables, ","))
	if app.contains(core_tables, name) {
		msg, _ := app.i18n.T("change_core_tables_not_allowed", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	/*table_id := any(nil)
	if _, ok := table_metadata["table_id"]; ok {
		table_id = table_metadata["table_id"]
	}
	table_org_name := ""
	if _, ok := table_metadata["table_org_name"]; ok {
		table_org_name = table_metadata["table_org_name"].(string)
	} else {
		table_org_name = name
	}*/
	//schema := generateCreateTableSQL(newDB.GetDriverName(), name, comment, fields)
	//fmt.Println(schema)
	schema_yaml := generateModelYAML(name, comment, fields)
	//fmt.Println(schema_yaml)
	model := getMDModel(name, schema_yaml, dsn)
	//fmt.Println(model)
	p := Dict{
		"db": dsn,
		"data": Dict{
			"order_metadata": any(true),
			"conf":           any(model),
		},
	}
	return app.etlxRun(p, true)
}

func getMDModel(tableName, yamlContent, dsn string) string {
	var out strings.Builder

	// Model header
	out.WriteString("# ADMMIN_MODEL\n")
	out.WriteString("```yaml\n")
	out.WriteString(fmt.Sprintf("name: %s\n", tableName))
	out.WriteString(fmt.Sprintf("description: %s\n", tableName))
	out.WriteString("runs_as: MODEL\n")
	out.WriteString(fmt.Sprintf("conn: '%s'\n", dsn))
	out.WriteString("```\n\n")

	// Table section
	out.WriteString(fmt.Sprintf("## %s\n", strings.ToUpper(tableName)))
	out.WriteString("```yaml\n")
	out.WriteString(yamlContent)
	out.WriteString("```\n\n")
	return out.String()
}
