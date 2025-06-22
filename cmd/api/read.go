package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/realdatadriven/etlx"
)

func (app *application) CrudRead(params map[string]any, table string, db *etlx.DB) map[string]any {
	/*var user_id int
	if _, ok := params["user"].(map[string]any)["user_id"]; ok {
		user_id = int(params["user"].(map[string]any)["user_id"].(float64))
	}*/
	var role_id int
	if _, ok := params["user"].(map[string]any)["role_id"]; ok {
		role_id = int(params["user"].(map[string]any)["role_id"].(float64))
	}
	/*var app_id int
	if _, ok := params["app"].(map[string]any)["app_id"]; ok {
		app_id = int(params["app"].(map[string]any)["app_id"].(float64))
	}*/
	//fmt.Println(user_id, role_id, app_id)
	_schema := map[string]any{}
	if _, ok := params["schema"]; ok {
		_schema = params["schema"].(map[string]any)
	}
	//fmt.Println("READ SCHEMA:", _schema["fields"])
	_permissions := map[string]any{}
	if _, ok := params["permissions"]; ok {
		_permissions = params["permissions"].(map[string]any)
	}
	roles := []any{role_id}
	if !app.contains(roles, 1) {
		if _, ok := _permissions["read"]; !ok {
			msg, _ := app.i18n.T("no-table-access", map[string]any{
				"table": table,
			})
			return map[string]any{
				"success": false,
				"msg":     msg,
			}
		} else if !app.contains([]any{true, 1}, _permissions["read"]) {
			msg, _ := app.i18n.T("no-table-action-access", map[string]any{
				"table":  table,
				"action": "READ",
			})
			return map[string]any{
				"success": false,
				"msg":     msg,
			}
		}
	}
	/*_row_level_tables := []string{}
	if _, ok := params["row_level_tables"]; ok {
		_row_level_tables = params["row_level_tables"].([]string)
	}*/
	limit := 10
	if _, ok := params["data"].(map[string]any)["limit"].(float64); ok {
		limit = int(params["data"].(map[string]any)["limit"].(float64))
	}
	offset := 0
	if _, ok := params["data"].(map[string]any)["offset"].(float64); ok {
		offset = int(params["data"].(map[string]any)["offset"].(float64))
	}
	// FIELDS
	_flds := []any{fmt.Sprintf(`"%s".*`, table)}
	//fields := []any{}
	if _, ok := params["data"].(map[string]any)["fields"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["fields"].([]any); ok {
		fields := params["data"].(map[string]any)["fields"].([]any)
		_flds = app._map2(fields, func(m any) any {
			m = fmt.Sprintf(`"%s"."%s"`, table, m.(string))
			return m
		})
	}
	// JOINS
	joins := []any{}
	join := ""
	if _, ok := params["data"].(map[string]any)["join"]; ok {
		join = params["data"].(map[string]any)["join"].(string)
	}
	if _, ok := params["data"].(map[string]any)["join_overwrite"]; ok {
		join_overwrite := params["data"].(map[string]any)["join_overwrite"].(map[string]any)
		if _, ok := join_overwrite[table]; ok {
			join = join_overwrite[table].(string)
		}
	}
	fk_tables_fields := map[string]any{}
	if _, ok := _schema["fields"]; !ok {
		// pass
		_schema["fields"] = map[string]any{}
	} else if join == "none" {
		// pass
	} else if join == "all" {
		for field, field_data := range _schema["fields"].(map[string]any) {
			if _, ok := field_data.(map[string]any)["fk"]; !ok {
			} else if field_data.(map[string]any)["fk"].(bool) {
				referred_table := ""
				if _, ok := field_data.(map[string]any)["referred_table"]; ok {
					referred_table = field_data.(map[string]any)["referred_table"].(string)
				}
				referred_column := field_data.(map[string]any)["referred_column"]
				if _, ok := field_data.(map[string]any)["referred_column"]; ok {
					referred_column = field_data.(map[string]any)["referred_column"].(string)
				}
				_referred_table_schema := map[string]any{}
				if _, ok := params["schemas"].(map[string]any)[referred_table]; ok {
					_referred_table_schema = params["schemas"].(map[string]any)[referred_table].(map[string]any)
				} else {
					_schemas := app.tables(params, []any{referred_table})
					if !_schemas["success"].(bool) {
						return _schemas
					}
					if _, ok := _schemas["data"]; ok {
						_schemas = _schemas["data"].(map[string]any)
						if _, ok := _schemas[referred_table]; ok {
							_referred_table_schema = _schemas[referred_table].(map[string]any)
						}
						params["schemas"].(map[string]any)[referred_table] = _referred_table_schema
					}
				}
				if referred_table != "" && referred_column != "" && len(_referred_table_schema) > 0 {
					joins = append(joins, fmt.Sprintf(`LEFT OUTER JOIN "%s" ON "%s"."%s" = "%s"."%s"`, referred_table, referred_table, referred_column, table, field))
					for key := range _referred_table_schema["fields"].(map[string]any) {
						//if _, ok := _schema["fields"].(map[string]any)[key]; !ok {
						_flds = append(_flds, fmt.Sprintf(`"%s"."%s" AS "%s_%s"`, referred_table, key, referred_table, key))
						//}
					}
					fk_tables_fields[referred_table] = _referred_table_schema
				}
				// fmt.Println(field, referred_table, referred_column)
			}
		}
	} else {
		for field, field_data := range _schema["fields"].(map[string]any) {
			if _, ok := field_data.(map[string]any)["fk"]; !ok {
			} else if field_data.(map[string]any)["fk"].(bool) {
				referred_table := ""
				if _, ok := field_data.(map[string]any)["referred_table"]; ok {
					referred_table = field_data.(map[string]any)["referred_table"].(string)
				}
				referred_column := field_data.(map[string]any)["referred_column"]
				if _, ok := field_data.(map[string]any)["referred_column"]; ok {
					referred_column = field_data.(map[string]any)["referred_column"].(string)
				}
				_referred_table_schema := map[string]any{}
				if _, ok := params["schemas"].(map[string]any)[referred_table]; ok {
					_referred_table_schema = params["schemas"].(map[string]any)[referred_table].(map[string]any)
				} else {
					_schemas := app.tables(params, []any{referred_table})
					if !_schemas["success"].(bool) {
						return _schemas
					}
					if _, ok := _schemas["data"]; ok {
						_schemas = _schemas["data"].(map[string]any)
						if _, ok := _schemas[referred_table]; ok {
							_referred_table_schema = _schemas[referred_table].(map[string]any)
						}
						params["schemas"].(map[string]any)[referred_table] = _referred_table_schema
					}
				}
				if referred_table != "" && referred_column != "" && len(_referred_table_schema) > 0 {
					if _, ok := _schema["fields"].(map[string]any)[referred_column.(string)]; !ok {
						continue
					}
					joins = append(joins, fmt.Sprintf(`LEFT OUTER JOIN "%s" ON "%s"."%s" = "%s"."%s"`, referred_table, referred_table, referred_column, table, field))
					keys := make([]any, len(_referred_table_schema["fields"].(map[string]any)))
					if _, ok := _referred_table_schema["fields_order"]; ok {
						keys = _referred_table_schema["fields_order"].([]any)
					} else {
						for key := range _referred_table_schema["fields"].(map[string]any) {
							keys = append(keys, key)
						}
					}
					if len(keys) > 1 {
						// fmt.Println(keys[1], keys)
						if _, ok := _schema["fields"].(map[string]any)[keys[1].(string)]; !ok {
							_flds = append(_flds, fmt.Sprintf(`"%s"."%s"`, referred_table, keys[1].(string)))
						}
					}
					fk_tables_fields[referred_table] = _referred_table_schema
				}
				// fmt.Println(field, referred_table, referred_column, len(_referred_table_schema), joins)
			}
		}
	}
	// FILTERS
	queryParams := []any{}
	filters := []any{}
	if _, ok := _schema["fields"].(map[string]any)["excluded"]; ok {
		filters = []any{fmt.Sprintf(`"%s"."excluded" IS FALSE`, table)}
	}
	if _, ok := params["data"].(map[string]any)["filters"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["filters"].([]any); ok {
		_filters := params["data"].(map[string]any)["filters"].([]any)
		for _, filter := range _filters {
			_field := filter.(map[string]any)["field"].(string)
			_cond := "="
			if _, ok := filter.(map[string]any)["cond"]; ok {
				_cond = filter.(map[string]any)["cond"].(string)
			}
			var _value any
			if _, ok := filter.(map[string]any)["value"]; ok {
				_value = filter.(map[string]any)["value"]
			}
			if _, ok := params["data"].(map[string]any)["ignore_filter"]; ok {
				ignore_filter := params["data"].(map[string]any)["ignore_filter"] //.(map[string]any)
				switch _type := ignore_filter.(type) {
				case nil:
					_ = _type
				case string:
					if ignore_filter.(string) == table {
						continue
					}
				case map[string]any:
					if _, ok := ignore_filter.(map[string]any)[table]; ok {
						continue
					}
				case []any:
					if app.contains(ignore_filter.([]any), table) {
						continue
					}
				default:
					_ = _type
				}
			}
			if _, ok := params["data"].(map[string]any)["apply_only_to"]; ok {
				apply_only_to := params["data"].(map[string]any)["apply_only_to"] //.(map[string]any)
				switch _type := apply_only_to.(type) {
				case nil:
					_ = _type
				case string:
					if apply_only_to.(string) != table {
						continue
					}
				case map[string]any:
					if _, ok := apply_only_to.(map[string]any)[table]; !ok {
						continue
					}
				case []any:
					if !app.contains(apply_only_to.([]any), table) {
						continue
					}
				default:
					_ = _type
				}
			}
			_table := table
			// allow fields that in the join tables to be passed as filters
			is_in_fk_fields := false
			if len(fk_tables_fields) > 0 {
				for _tbl, _tbl_fields := range fk_tables_fields {
					if _, ok := _tbl_fields.(map[string]any)[_field]; ok {
						is_in_fk_fields = true
						_table = _tbl
						break
					}
				}
			}
			if _, ok := _schema["fields"].(map[string]any)[_field]; !ok && !is_in_fk_fields {
				// pass fm
			} else if app.contains([]any{"=", "!=", ">", "<", ">=", "<="}, _cond) {
				filters = append(filters, fmt.Sprintf(`"%s"."%s" %s ?`, _table, _field, _cond))
				queryParams = append(queryParams, _value)
			} else if app.contains([]any{"in", "not in"}, strings.ToLower(_cond)) {
				filters = append(filters, fmt.Sprintf(`"%s"."%s" %s (?)`, _table, _field, _cond))
				queryParams = append(queryParams, strings.Split(_value.(string), ","))
			} else if app.contains([]any{"between", "not between"}, strings.ToLower(_cond)) {
				filters = append(filters, fmt.Sprintf(`"%s"."%s" %s ? AND ?`, _table, _field, _cond))
				queryParams = append(queryParams, strings.Split(_value.(string), ","))
			} else if app.contains([]any{"like", "not like"}, strings.ToLower(_cond)) {
				//lwildc := "%"
				//rwildc := "%"
				if strings.Contains(_value.(string), "%") {
					//lwildc = ""
					//rwildc = ""
				} else {
					_value = fmt.Sprintf(`%%%s%%`, _value)
				}
				//filters = append(filters, fmt.Sprintf(`"%s"."%s" %s '%s?%s'`, _table, _field, _cond, lwildc, rwildc))
				filters = append(filters, fmt.Sprintf(`"%s"."%s" %s ?`, _table, _field, _cond))
				queryParams = append(queryParams, _value)
			} else if app.contains([]any{"is true", "is false", "is null", "is not null"}, strings.ToLower(_cond)) {
				filters = append(filters, fmt.Sprintf(`"%s"."%s" %s`, _table, _field, _cond))
				//queryParams = append(queryParams, _value)
			}
		}
	}
	// ORDER BY
	orderBy := []any{}
	if _, ok := params["data"].(map[string]any)["order_by"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["order_by"].([]any); ok {
		_order_by := params["data"].(map[string]any)["order_by"].([]any)
		for _, oby := range _order_by {
			_field := oby.(map[string]any)["field"].(string)
			_order := "ASC"
			if _, ok := oby.(map[string]any)["order"]; ok {
				_order = oby.(map[string]any)["order"].(string)
			}
			if _, ok := _schema["fields"].(map[string]any)[_field]; !ok {
				// pass
			} else {
				orderBy = append(orderBy, fmt.Sprintf(`"%s"."%s" %s`, table, _field, _order))
			}
		}
	}
	distinct := ""
	if _, ok := params["data"].(map[string]any)["distinct"]; ok {
		if params["data"].(map[string]any)["distinct"].(bool) {
			distinct = "DISTINCT"
		}
	}
	search_patt := []any{}
	if _, ok := params["data"].(map[string]any)["pattern"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["pattern"].(string); !ok {
	} else if _pattern, ok := params["data"].(map[string]any)["pattern"].(string); ok && _pattern != "" {
		//_pattern := params["data"].(map[string]any)["pattern"].(string)
		key := "%" + _pattern + "%"
		// _split_pattern = re.compile(r'\||\;')
		re := regexp.MustCompile(`[|;]`)
		_splited_keys := re.Split(_pattern, -1)
		// for
		if strings.Contains(_pattern, "%") {
			for _field, field_data := range _schema["fields"].(map[string]any) {
				_type := field_data.(map[string]any)["type"].(string)
				_type = strings.ToLower(_type)
				if app.contains([]any{"bool", "boolean", "bool", "boolean", "int", "integer", "float", "real", "decimal"}, _type) {
					continue
				}
				search_patt = append(search_patt, fmt.Sprintf(`CAST("%s"."%s" AS VARCHAR) LIKE ?`, table, _field))
				queryParams = append(queryParams, _pattern)
			}
		} else if len(_splited_keys) <= 1 {
			for _field, field_data := range _schema["fields"].(map[string]any) {
				_type := field_data.(map[string]any)["type"].(string)
				_type = strings.ToLower(_type)
				if app.contains([]any{"bool", "boolean", "int", "integer", "float", "real", "decimal"}, _type) {
					continue
				}
				search_patt = append(search_patt, fmt.Sprintf(`CAST("%s"."%s" AS VARCHAR) LIKE ?`, table, _field))
				queryParams = append(queryParams, key)
			}
		} else {
			for _, k := range _splited_keys {
				if strings.Contains(k, "%") {
					for _field, field_data := range _schema["fields"].(map[string]any) {
						_type := field_data.(map[string]any)["type"].(string)
						_type = strings.ToLower(_type)
						if app.contains([]any{"bool", "boolean", "int", "integer", "float", "real", "decimal"}, _type) {
							continue
						}
						search_patt = append(search_patt, fmt.Sprintf(`CAST("%s"."%s" AS VARCHAR) LIKE ?`, table, _field))
						queryParams = append(queryParams, k)
					}
				} else {
					for _field, field_data := range _schema["fields"].(map[string]any) {
						_type := field_data.(map[string]any)["type"].(string)
						_type = strings.ToLower(_type)
						if app.contains([]any{"bool", "boolean", "int", "integer", "float", "real", "decimal"}, _type) {
							continue
						}
						search_patt = append(search_patt, fmt.Sprintf(`CAST("%s"."%s" AS VARCHAR) LIKE ?`, table, _field))
						queryParams = append(queryParams, k)
					}
				}
			}
		}
	}
	query := fmt.Sprintf(`SELECT %s %s FROM "%s"`, distinct, app.joinSlice(_flds, ","), table)
	if len(joins) > 0 {
		query = fmt.Sprintf(`%s %s`, query, app.joinSlice(joins, "\n"))
	}
	if len(filters) > 0 {
		query = fmt.Sprintf(`%s WHERE (%s)`, query, app.joinSlice(filters, " AND "))
	}
	if len(search_patt) > 0 {
		_where := " WHERE"
		if len(filters) > 0 {
			_where = " AND"
		}
		query = fmt.Sprintf(`%s%s (%s)`, query, _where, app.joinSlice(search_patt, " OR "))
	}
	if len(orderBy) > 0 {
		query = fmt.Sprintf(`%s ORDER BY %s`, query, app.joinSlice(orderBy, ", "))
	}
	if limit != -1 {
		query = fmt.Sprintf(`%s LIMIT %d OFFSET %d`, query, limit, offset)
	}
	query, args, err := sqlx.In(query, queryParams...)
	if err != nil {
		println("Error geting the table query:", err)
	}
	//fmt.Println(query, args)
	results := make([]map[string]any, 0)
	data, _, err := db.QueryMultiRows(query, args...)
	if err != nil {
		fmt.Println("READ ERR:", args, query, err)
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
			"sql":     query,
		}
	} else if *data == nil || len(*data) == 0 {
	} else {
		results = *data
	}
	total := 0
	query_total := fmt.Sprintf(`SELECT COUNT(*) AS "n_rows" FROM (%s) AS "T"`, query)
	//fmt.Println(query_total)
	trows, _, err := db.QuerySingleRow(query_total, args...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	// fmt.Println((*trows))
	total = int((*trows)["n_rows"].(int64))
	//fmt.Println(app_id, user_id, pk, args, query)
	//data := map[string]any{}
	msg, _ := app.i18n.T("success", map[string]any{})
	return map[string]any{
		"success": true,
		"msg":     msg,
		"data":    results,
		"total":   total,
		"cols":    _schema["fields_order"],
		//"schema":           _schema,
		"permissions": _permissions,
		//"row_level_tables": _row_level_tables,
		"sql": query,
	}
}
