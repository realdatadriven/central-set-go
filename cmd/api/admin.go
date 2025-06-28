package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"

	"github.com/jmoiron/sqlx"
)

func (app *application) apps(params map[string]any) map[string]any {
	//fmt.Println("APPS:", params)
	user_id := int(params["user"].(map[string]any)["user_id"].(float64))
	role_id := int(params["user"].(map[string]any)["role_id"].(float64))
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
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	roles := []any{}
	roles = append(roles, role_id)
	for _, row := range *result {
		roles = append(roles, int(row["role_id"].(float64)))
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
		WHERE role_app.role_id IN ($1)
			AND role_app.access = TRUE
			AND role_app.excluded = FALSE
			AND app.excluded = FALSE`
		//fmt.Println(app.joinSlice(roles, ","))
		queryParams = append(queryParams, app.joinSlice(roles, ","))
	}
	result, _, err = app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    *result,
	}
}

func (app *application) menu(params map[string]any) map[string]any {
	//fmt.Println(params)
	user_id := int(params["user"].(map[string]any)["user_id"].(float64))
	role_id := int(params["user"].(map[string]any)["role_id"].(float64))
	var app_id int
	if _, ok := params["app"].(map[string]any)["app_id"]; ok {
		app_id = int(params["app"].(map[string]any)["app_id"].(float64))
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
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	roles := []any{}
	roles = append(roles, role_id)
	for _, row := range *result {
		roles = append(roles, int(row["role_id"].(float64)))
	}
	// MENU
	query = `SELECT *
	FROM menu
	WHERE app_id = $1
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
		WHERE menu.app_id = $1
			AND role_app_menu.role_id IN ($2)
			AND role_app_menu.access = TRUE
			AND role_app_menu.excluded = FALSE
			AND menu.excluded = FALSE
			AND menu.active = TRUE
			AND (role_app_menu.menu_id, role_app_menu.app_id, role_app_menu.updated_at) IN (
				SELECT menu_id, app_id, MAX(updated_at)
				FROM role_app_menu
				WHERE access = True
					AND role_id IN ($2)
					AND excluded = FALSE
				GROUP BY menu_id, app_id
			)
		ORDER BY menu.menu_order ASC, menu.menu_id ASC`
		//fmt.Println(app.joinSlice(roles, ","))
		queryParams = append(queryParams, app.joinSlice(roles, ","))
	}
	_menu, _, err := app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	// MENU TABLES
	query = `SELECT *
	FROM menu_table
	WHERE app_id = $1
		AND excluded = FALSE
	ORDER BY menu_table_id ASC`
	queryParams = []any{app_id}
	if !app.contains(roles, 1) {
		query = `SELECT menu_table.*
		FROM menu_table
		JOIN role_app_menu_table ON (
			role_app_menu_table.menu_id = menu_table.menu_id 
			role_app_menu_table.table_id = menu_table.table_id 
			AND role_app_menu_table.app_id = menu_table.app_id
		)
		WHERE menu_table.app_id = $1
			AND role_app_menu_table.role_id IN ($2)
			AND (
				role_app_menu_table.read = TRUE
				OR role_app_menu_table.create = TRUE
			)
			AND role_app_menu_table.excluded = FALSE
			AND menu_table.excluded = FALSE
			AND (role_app_menu_table.table_id, role_app_menu_table.menu_id, role_app_menu_table.app_id, role_app_menu_table.updated_at) IN (
				SELECT table_id, menu_id, app_id, MAX(updated_at)
				FROM role_app_menu_table
				WHERE access = True
					AND role_id IN ($2)
					AND excluded = FALSE
				GROUP BY table_id, menu_id, app_id
			)
		ORDER BY menu_table_id ASC`
		queryParams = append(queryParams, app.joinSlice(roles, ","))
	}
	_menu_table, _, err := app.db.QueryMultiRows(query, queryParams...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
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
		_tables = _tables["data"].(map[string]any)
	}
	// MENUS
	menus := []map[string]any{}
	for _, mn := range *_menu {
		_aux := mn
		_aux["children"] = []map[string]any{}
		for _, mnt := range *_menu_table {
			if _, ok := mnt["menu_id"].(any); !ok {
			} else if _, ok := mn["menu_id"].(any); !ok {
			} else if int(mnt["menu_id"].(int64)) == int(mn["menu_id"].(int64)) {
				_mnt := mnt
				//fmt.Println(1, _table_by_id[mnt["table_id"].(int64)].(map[string]any))
				if _, ok := _table_by_id[mnt["table_id"].(int64)].(map[string]any); ok {
					_mnt["table"] = _table_by_id[mnt["table_id"].(int64)].(map[string]any)["table"].(string)
					//fmt.Println(2, _table_by_id[mnt["table_id"].(int64)].(map[string]any))
				}
				_mnt["menu"] = mn["menu"]
				_aux["children"] = append(_aux["children"].([]map[string]any), _mnt)
			}
		}
		menus = append(menus, _aux)
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data": map[string]any{
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

func (app *application) GetDBNameFromParams(params map[string]any) (string, string, error) {
	var _database any
	if !app.IsEmpty(params["db"]) {
		_database = params["db"]
	} else if !app.IsEmpty(params["data"].(map[string]any)["db"]) {
		_database = params["data"].(map[string]any)["db"]
	} else if !app.IsEmpty(params["data"].(map[string]any)["database"]) {
		_database = params["data"].(map[string]any)["database"]
	} else if !app.IsEmpty(params["app"].(map[string]any)["db"]) {
		_database = params["app"].(map[string]any)["db"]
	} else if !app.IsEmpty(params["app"].(map[string]any)["db"]) {
		_database = params["app"].(map[string]any)["db"]
	}
	//_not_embed_dbs := []any{"postgres", "postgresql", "pg", "pgql", "mysql"}
	_embed_dbs := []any{"sqlite", "sqlite3", "duckdb", "ducklake"}
	_embed_dbs_ext := []any{".db", ".duckdb", ".ddb", ".sqlite", ".ducklake"}
	//fmt.Println(1, _database)
	switch _type := _database.(type) {
	case nil:
		return app.config.db.dsn, "", nil
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
			if err == nil {
				_database = dbname
			}
		}
		return dsn, _database.(string), nil
	case []any:
		fmt.Println("IS []any:", _database, _type)
		return "", "", errors.New("database conf is of type []any")
	default:
		return _database.(string), _database.(string), nil
	}
}

func (app *application) tables(params map[string]any, tables []any) map[string]any {
	//fmt.Println(params)
	var user_id int
	if _, ok := params["user"].(map[string]any)["user_id"]; ok {
		user_id = int(params["user"].(map[string]any)["user_id"].(float64))
	}
	var app_id int
	if _, ok := params["app"].(map[string]any)["app_id"]; ok {
		app_id = int(params["app"].(map[string]any)["app_id"].(float64))
	}
	// DATABASE
	_extra_conf := map[string]any{
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
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	allTables := false
	if app.IsEmpty(tables) {
		tables = []any{}
		if !app.IsEmpty(params["data"].(map[string]any)["table"]) {
			value := params["data"].(map[string]any)["table"]
			switch value.(type) {
			case nil:
				// pass
			case string:
				tables = append(tables, params["data"].(map[string]any)["table"].(string))
			case []any:
				_tables := params["data"].(map[string]any)["table"].([]any)
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			case map[any]any:
				// pass
			default:
				tables = append(tables, params["data"].(map[string]any)["table"].(string))
			}
		} else if !app.IsEmpty(params["data"].(map[string]any)["tables"]) {
			value := params["data"].(map[string]any)["tables"]
			switch value.(type) {
			case string:
				tables = append(tables, params["data"].(map[string]any)["tables"].(string))
			case []any:
				_tables := params["data"].(map[string]any)["tables"].([]any)
				for t := 0; t < len(_tables); t++ {
					tables = append(tables, _tables[t])
				}
			default:
				tables = append(tables, params["data"].(map[string]any)["table"].(string))
			}
		}
		if app.IsEmpty(tables) {
			// fmt.Println("GET ALL TABLES!")
			result, _, err := newDB.AllTables(params, _extra_conf)
			if err != nil {
				return map[string]any{
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
	//fmt.Println(dsn, _database, tables, allTables)
	data := map[string]any{}
	table_by_id := map[int64]any{}
	if app.IsEmpty(tables) {
		msg, _ := app.i18n.T("no-table", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
			"tables":  tables,
		}
	} else {
		// GET THE TABLES DATA IN table
		query := `SELECT * FROM "table" WHERE db = ? AND "table" IN (?) AND excluded = FALSE`
		queryParams := []any{_database}
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
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		// fmt.Println(_table)
		if allTables {
			tables_in_table := []any{}
			for _, row := range *_table {
				tables_in_table = append(tables_in_table, row["table"].(string))
			}
			// fmt.Println(tables_in_table)
			results := []map[string]any{}
			for _, table := range tables {
				if !app.contains(tables_in_table, table) {
					//fmt.Println("ADD TABLE:", table)
					results = append(results, map[string]any{
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
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		translate_table := map[string]any{}
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
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		translate_table_field := map[string]any{}
		for _, row := range *results {
			if _, ok := translate_table_field[row["table"].(string)]; !ok {
				translate_table_field[row["table"].(string)] = map[string]any{}
			}
			/*if _, ok := translate_table_field[row["table"].(string)].(map[string]any)["fields"]; !ok {
				translate_table_field[row["table"].(string)].(map[string]any)["fields"] = map[string]any{}
			}
			translate_table_field[row["table"].(string)].(map[string]any)["fields"].(map[string]any)[row["field"].(string)] = row*/
			translate_table_field[row["table"].(string)].(map[string]any)[row["field"].(string)] = row
		}
		// fmt.Println(translate_table_field)
		// GET THE TABLES DATA IN table_schema
		query = `SELECT * FROM table_schema WHERE db = ? AND "table" IN (?) AND excluded = FALSE`
		queryParams = []any{_database}
		if allTables {
			query = `SELECT * FROM table_schema WHERE db = ? AND excluded = FALSE`
		} else {
			queryParams = append(queryParams, tables)
		}
		//queryParams = append(queryParams, app.joinSlice(tables, "','"))
		query, args, err = sqlx.In(query, queryParams...)
		if err != nil {
			println("Error geting the table query:", err)
		}
		//fmt.Println(query, args, queryParams)
		table_schema := map[string]any{}
		_table_schema, _, err := app.db.QueryMultiRows(query, args...)
		if err != nil {
			fmt.Println("TABLE SCHEMA:", query, err)
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		//fmt.Println(*_table_schema)
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
							//fmt.Println(query)
							for _, row := range results {
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
					return map[string]any{
						"success": false,
						"msg":     fmt.Sprintf("%s", err),
					}
				}
			}
			//fmt.Println("tables_not_in_schema:", tables_not_in_schema)
		}
		table_fields := map[string]any{}
		for _, row := range *_table_schema {
			if _, ok := table_schema[row["table"].(string)]; !ok {
				table_schema[row["table"].(string)] = map[string]any{}
			}
			if _, ok := table_fields[row["table"].(string)]; !ok {
				table_fields[row["table"].(string)] = []any{}
			}
			_row := row
			/*if _, ok := table_schema[row["table"].(string)].(map[string]any)["fields"]; !ok {
				table_schema[row["table"].(string)].(map[string]any)["fields"] = map[string]any{}
			}
			table_schema[row["table"].(string)].(map[string]any)["fields"].(map[string]any)[row["field"].(string)] = row*/
			comment := _row["comment"]
			if _, ok := translate_table_field[row["table"].(string)]; !ok {
			} else if _, ok := translate_table_field[row["table"].(string)].(map[string]any)[row["field"].(string)]; !ok {
			} else if _, ok := translate_table_field[row["table"].(string)].(map[string]any)[row["field"].(string)].(map[string]any)["field_transl_desc"]; ok {
				comment = translate_table_field[row["table"].(string)].(map[string]any)[row["field"].(string)].(map[string]any)["field_transl_desc"]
			}
			_row["comment"] = comment
			_row["name"] = _row["field"]
			if _, ok := _row["fk"]; !ok {
			} else if app.contains([]any{1, true, "true", "True", "TRUE", "T", "1"}, _row["fk"]) {
				referred_columns_desc := ""
				if _, ok := table_fields[row["referred_table"].(string)].([]any); ok {
					referred_columns_desc = table_fields[row["referred_table"].(string)].([]any)[1].(string)
				}
				_row["ref"] = map[string]any{
					"referred_table":        _row["referred_table"],
					"referred_column":       _row["referred_column"],
					"referred_columns_desc": referred_columns_desc,
				}
			}
			table_schema[row["table"].(string)].(map[string]any)[row["field"].(string)] = _row
			table_fields[row["table"].(string)] = append(table_fields[row["table"].(string)].([]any), row["field"])
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
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		custom_form := map[string]any{}
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
			fmt.Println("custom_table:", query, err)
			return map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%s", err),
			}
		}
		custom_table := map[string]any{}
		for _, row := range *results {
			custom_table[row["table"].(string)] = row
		}
		// return
		for _, row := range *_table {
			comment := row["table_desc"]
			if _, ok := translate_table[row["table"].(string)]; ok {
				comment = translate_table[row["table"].(string)].(map[string]any)["table_transl_desc"]
			}
			var pk string
			if _, ok := table_schema[row["table"].(string)]; ok {
				for key, value := range table_schema[row["table"].(string)].(map[string]any) {
					if properties, ok := value.(map[string]any); ok {
						// Check if the "pk" field exists and is true
						if _pk, found := properties["pk"]; found && _pk == true {
							pk = key
							break
						}
					}
				}
			}
			table_by_id[row["table_id"].(int64)] = row
			data[row["table"].(string)] = map[string]any{
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
				"pk":                    pk,
				"fields_order":          table_fields[row["table"].(string)],
			}
		}
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success":     true,
		"msg":         msg,
		"data":        data,
		"table_by_id": table_by_id,
	}
}

// Generates CREATE TABLE SQL statements with comments, adapting to SQL dialects
func generateCreateTableSQL(driver, tableName, tableComment string, fields []map[string]any) string {
	var schema strings.Builder
	schema.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	// Collect foreign keys and comments
	var foreignKeys []string
	var columnComments []string

	// Generate column definitions based on the driver
	for _, field := range fields {
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
func getColumnType(driver string, field map[string]any) string {
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
func getPrimaryKey(driver string, field map[string]any) string {
	if pk, ok := field["primary_key"].(bool); ok && pk {
		if driver == "mysql" || driver == "sqlserver" || driver == "mssql" {
			return " PRIMARY KEY"
		}
	}
	return ""
}

// Autoincrement syntax adjustments per driver
func getAutoIncrement(driver string, field map[string]any) string {
	if field["autoincrement"] == true {
		if driver == "sqlite3" {
			return " AUTOINCREMENT"
		}
	}
	return ""
}

// Nullable syntax adjustments per driver
func getNullable(driver string, field map[string]any) string {
	if nullable, ok := field["nullable"].(bool); ok && !nullable {
		return " NOT NULL"
	}
	return ""
}

// Unique constraint syntax adjustments
func getUnique(driver string, field map[string]any) string {
	if unique, ok := field["unique"].(bool); ok && unique {
		return " UNIQUE"
	}
	return ""
}

// Default value handling based on driver
func getDefaultValue(driver string, field map[string]any) string {
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

func (app *application) save_table_schema(params map[string]any) map[string]any {
	//fmt.Println(params)
	//user_id := int(params["user"].(map[string]any)["user_id"].(float64))
	//role_id := int(params["user"].(map[string]any)["role_id"].(float64))
	//var app_id int
	//if _, ok := params["app"].(map[string]any)["app_id"]; ok {
	//	app_id = int(params["app"].(map[string]any)["app_id"].(float64))
	//}
	// DATABASE
	//fmt.Println(lang)
	dsn, _, _ := app.GetDBNameFromParams(params)
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	defer newDB.Close()
	_data := map[string]any{}
	if _, ok := params["data"]; !ok {
		msg, _ := app.i18n.T("no_data", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := params["data"].(map[string]any); ok {
		_data = params["data"].(map[string]any)
	}
	table_metadata := map[string]any{}
	if _, ok := _data["table_metadata"]; !ok {
		msg, _ := app.i18n.T("no_table_metadata", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := _data["table_metadata"].(map[string]any); !ok {
		msg, _ := app.i18n.T("no_table_metadata", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	table_metadata = _data["table_metadata"].(map[string]any)
	name := ""
	if _, ok := table_metadata["name"]; !ok {
		msg, _ := app.i18n.T("no_table_name", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := table_metadata["name"].(string); !ok {
		msg, _ := app.i18n.T("no_table_name", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	name = table_metadata["name"].(string)
	comment := ""
	if _, ok := table_metadata["comment"]; !ok {
		msg, _ := app.i18n.T("no_table_comment", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := table_metadata["comment"].(string); !ok {
		msg, _ := app.i18n.T("no_table_comment", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	comment = table_metadata["comment"].(string)
	fields := []map[string]any{}
	if _, ok := table_metadata["fields"]; !ok {
		msg, _ := app.i18n.T("no_fields", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	} else if _, ok := table_metadata["fields"].([]map[string]any); !ok {
		msg, _ := app.i18n.T("no_fields", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	fields = table_metadata["fields"].([]map[string]any)
	if len(fields) < 2 {
		msg, _ := app.i18n.T("table_must_have_2_or_more_fields", map[string]any{})
		return map[string]any{
			"success": false,
			"msg":     msg,
		}
	}
	core_tables := app.sliceStrs2SliceInterfaces(strings.Split(app.config.core_tables, ","))
	if app.contains(core_tables, name) {
		msg, _ := app.i18n.T("change_core_tables_not_allowed", map[string]any{})
		return map[string]any{
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
	schema := generateCreateTableSQL(newDB.GetDriverName(), name, comment, fields)
	fmt.Println(schema)
	/*/ Map for SQLAlchemy types to SQL types
	var saTypesToSQL = map[string]string{
		"Integer":  "INTEGER",
		"String":   "VARCHAR",
		"Text":     "TEXT",
		"Date":     "DATE",
		"DateTime": "DATETIME",
		"Time":     "TIME",
		"Float":    "DECIMAL",
		"Boolean":  "BOOLEAN",
	}*/
	msg, _ := app.i18n.T("sql-generated-to-be validated", map[string]any{})
	return map[string]any{
		"success": false,
		"msg":     msg,
	}
}
