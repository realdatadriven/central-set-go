package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/realdatadriven/etlx"
)

func (app *application) containsInt(slice []any, element any) bool {
	for _, v := range slice {
		if v.(int64) == element.(int64) {
			return true
		}
	}
	return false
}

func (app *application) isSafeSQLIdentifier(value string) bool {
	if value == "" {
		return false
	}
	ok, err := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_]*$`, value)
	return err == nil && ok
}

func (app *application) getRLAIds(rla_access []map[string]any, table, access_type string, my_ids []any) []any {
	data := []any{}
	for _, v := range rla_access {
		_access := false
		switch v[access_type].(type) {
		case bool:
			_access = v[access_type].(bool)
		case int:
			_access = v[access_type].(int) == 1
		case float32:
			_access = v[access_type].(float32) == 1
		case float64:
			_access = v[access_type].(float64) == 1
		case int64:
			_access = v[access_type].(int64) == 1
		case int32:
			_access = v[access_type].(int32) == 1
		case string:
			i, err := strconv.Atoi(v[access_type].(string))
			if err != nil {
				_access = false
			}
			_access = i == 1
		default:
			_access = false // or handle error
		}
		if v["table"].(string) == table && _access {
			data = append(data, v["row_id"])
		}
	}
	if len(my_ids) > 0 {
		for _, v := range my_ids {
			data = append(data, v)
		}
	}
	return data
}
func (app *application) CrudRead(params map[string]any, table string, db etlx.DBInterface) map[string]any {
	var user_id int
	if _, ok := params["user"].(map[string]any)["user_id"]; ok {
		user_id = app.toInt(params["user"].(map[string]any)["user_id"])
	}
	var role_id int
	if _, ok := params["user"].(map[string]any)["role_id"]; ok {
		role_id = app.toInt(params["user"].(map[string]any)["role_id"])
	}
	var loc *time.Location
	if _, ok := params["location"].(*time.Location); ok {
		loc = params["location"].(*time.Location)
	} else {
		loc = time.Local
	}
	/*var app_id int
	if _, ok := params["app"].(map[string]any)["app_id"]; ok {
		app_id = app.toInt(params["app"].(map[string]any)["app_id"])
	}*/
	//fmt.Println(user_id, role_id, app_id)
	_schema := map[string]any{}
	if _, ok := params["schema"]; ok {
		_schema = params["schema"].(map[string]any)
	}
	pk := ""
	if _, ok := _schema["pk"]; ok {
		pk = _schema["pk"].(string)
	}
	//fmt.Println("READ SCHEMA:", _schema["fields"])
	_permissions := map[string]any{}
	if _, ok := params["permissions"]; ok {
		_permissions = params["permissions"].(map[string]any)
	}
	roles := []any{role_id}
	if !app.contains(roles, 1) {
		//fmt.Println(_permissions)
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
	limit := 10
	if _, ok := params["data"].(map[string]any)["limit"]; ok {
		limit = app.toInt(params["data"].(map[string]any)["limit"])
	} /*else if _, ok := params["data"].(map[string]any)["limit"]; ok {
		limit = app.toInt(params["data"].(map[string]any)["limit"])
	} else if _, ok := params["data"].(map[string]any)["limit"].(int); ok {
		limit = params["data"].(map[string]any)["limit"].(int)
	}*/
	offset := 0
	if _, ok := params["data"].(map[string]any)["offset"]; ok {
		offset = app.toInt(params["data"].(map[string]any)["offset"])
	}
	table_schema := ""
	if _, ok := params["data"].(Dict)["schema"]; ok {
		table_schema = params["data"].(Dict)["schema"].(string)
	}
	if !app.isSafeSQLIdentifier(table) {
		return map[string]any{
			"success": false,
			"msg":     "Invalid table identifier",
		}
	}
	schm := ""
	if table_schema != "" {
		if !app.isSafeSQLIdentifier(table_schema) {
			return map[string]any{
				"success": false,
				"msg":     "Invalid schema identifier",
			}
		}
		schm = fmt.Sprintf(`%s.`, table_schema)
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
	// isRefToSameTable := false
	fk_tables_fields := map[string]any{}
	fk_tables_added := []any{}
	fk_tables_pk := Dict{}
	if _, ok := _schema["fields"]; ok {
		for _, field_data := range _schema["fields"].(map[string]any) {
			//fmt.Println(field_data.(map[string]any)["name"], field_data.(map[string]any)["fk"], field_data.(map[string]any)["fk"] == any(true), (field_data.(map[string]any)["fk"] == any(1) || field_data.(map[string]any)["fk"] == any(1.0)))
			_fk, ok := field_data.(map[string]any)["fk"]
			if !ok {
				//} else if _fk == any(true) || (_fk == any(1) || _fk == any(1.0)) {
			} else if app.toBool(_fk) {
				referred_table := ""
				if _, ok := field_data.(map[string]any)["referred_table"]; ok {
					referred_table = field_data.(map[string]any)["referred_table"].(string)
				}
				referred_column := field_data.(map[string]any)["referred_column"]
				fk_tables_pk[referred_table] = referred_column
				//fmt.Println("referred_table:", referred_table, "referred_column:", referred_column)
				if referred_table == table { // TOVAOID OVERWRITING MAIN TABLE COLUMNS
					fk_tables_added = append(fk_tables_added, referred_table)
				}
			}
		}
	} else {
		fmt.Println(table, "NO SCHEMA FIELDS (READ)")
	}
	if _, ok := _schema["fields"]; !ok {
		// pass
		_schema["fields"] = map[string]any{}
	} else if join == "none" {
		// pass
	} else if join == "all" {
		//fmt.Println(table, 1)
		for field, field_data := range _schema["fields"].(map[string]any) {
			if _, ok := field_data.(map[string]any)["fk"]; !ok {
				//fmt.Println(field, 1)
				//fmt.Println(table, field, 2, field_data.(map[string]any)["fk"])
				//} else if field_data.(map[string]any)["fk"] == any(true) || (field_data.(map[string]any)["fk"] == any(1) || field_data.(map[string]any)["fk"] == any(1.0)) {
			} else if app.toBool(field_data.(map[string]any)["fk"]) {
				// fmt.Printf("%s -> %s %v %T\n", table, field, field_data.(map[string]any)["fk"], field_data.(map[string]any)["fk"])
				referred_table := ""
				level := 1
				if _, ok := field_data.(map[string]any)["referred_table"]; ok {
					referred_table = field_data.(map[string]any)["referred_table"].(string)
					fk_tables_added = append(fk_tables_added, referred_table)
					acorr := app.filterAny(fk_tables_added, func(r any) bool { return r.(string) == referred_table })
					if len(acorr) > 1 {
						level = len(acorr)
					}
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
				alias := referred_table
				field_sufix := ""
				if level > 1 {
					alias = fmt.Sprintf("%s%d", alias, level)
					field_sufix = fmt.Sprintf("%d", level)
				}
				if referred_table != "" && referred_column != "" && len(_referred_table_schema) > 0 {
					joins = append(joins, fmt.Sprintf(`LEFT OUTER JOIN %s"%s" "%s" ON "%s"."%s" = "%s"."%s"`, schm, referred_table, alias, alias, referred_column, table, field))
					for key := range _referred_table_schema["fields"].(map[string]any) {
						//if _, ok := _schema["fields"].(map[string]any)[key]; !ok {
						_flds = append(_flds, fmt.Sprintf(`"%s"."%s" AS "%s_%s%s"`, alias, key, referred_table, key, field_sufix))
						//}
					}
					fk_tables_fields[alias] = _referred_table_schema["fields"].(map[string]any)
				}
				// fmt.Println(field, referred_table, referred_column)
			} else {
				// fmt.Printf("%s -> %s %v %T\n", table, field, field_data.(map[string]any)["fk"], field_data.(map[string]any)["fk"])
			}
		}
	} else {
		for field, field_data := range _schema["fields"].(map[string]any) {
			if _, ok := field_data.(map[string]any)["fk"]; !ok {
				//} else if field_data.(map[string]any)["fk"] == any(true) || (field_data.(map[string]any)["fk"] == any(1) || field_data.(map[string]any)["fk"] == any(1.0)) {
			} else if app.toBool(field_data.(map[string]any)["fk"]) {
				referred_table := ""
				level := 1
				if _, ok := field_data.(map[string]any)["referred_table"]; ok {
					referred_table = field_data.(map[string]any)["referred_table"].(string)
					fk_tables_added = append(fk_tables_added, referred_table)
					acorr := app.filterAny(fk_tables_added, func(r any) bool { return r.(string) == referred_table })
					if len(acorr) > 1 {
						level = len(acorr)
					}
					// fmt.Println(1, "REF TABLE", referred_table, "OCORRR", acorr, fk_tables_added, "LEVEL", level)
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
				alias := referred_table
				field_sufix := ""
				if level > 1 {
					alias = fmt.Sprintf("%s%d", alias, level)
					field_sufix = fmt.Sprintf("%d", level)
				}
				if referred_table != "" && referred_column != "" && len(_referred_table_schema) > 0 {
					if _, ok := _schema["fields"].(map[string]any)[referred_column.(string)]; !ok {
						continue
					}
					joins = append(joins, fmt.Sprintf(`LEFT OUTER JOIN %s"%s" "%s" ON "%s"."%s" = "%s"."%s"`, schm, referred_table, alias, alias, referred_column, table, field))
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
							_flds = append(_flds, fmt.Sprintf(`"%s"."%s" AS "%s%s"`, alias, keys[1].(string), keys[1].(string), field_sufix))
							if field_sufix != "" {
								// fmt.Printf(`"%s"."%s" AS "%s%s"\n`, alias, keys[1].(string), keys[1].(string), field_sufix)
							}
						}
					}
					fk_tables_fields[alias] = _referred_table_schema["fields"].(map[string]any) //_schema["fields"].(map[string]any)
				}
				// fmt.Println(field, referred_table, referred_column, len(_referred_table_schema), joins)
			}
		}
	}

	if table == "menu_table" {
		//fmt.Printf("%s: %s -> %v\n", table, join, _flds)
	}
	// CHECK ROW LEVEL ACCESS, IF SO UPDATE FILTERS field_id in (?) ? = row_id allowed
	_row_level_tables := []string{}
	rla_tables_ids := Dict{}
	if !app.contains(roles, 1) {
		if _, ok := params["row_level_tables"]; ok {
			_row_level_tables = params["row_level_tables"].([]string)
			_tables_to_chk := []any{}
			for _, v := range append(fk_tables_added, table) {
				if app.contains(app.sliceStrs2SliceInterfaces(_row_level_tables), v) {
					_tables_to_chk = append(_tables_to_chk, v)
				}
			}
			//fmt.Println("row_level_tables:", _row_level_tables, fk_tables_added, "_tables_to_chk:", _tables_to_chk)
			if len(_tables_to_chk) > 0 {
				rla_access := app.row_level_access(params, _tables_to_chk, []any{})
				//fmt.Println("rla_access:", _tables_to_chk, rla_access)
				if !rla_access["success"].(bool) {
					return rla_access
				} else if _rla_access_data, ok := rla_access["data"].(map[string]any); ok {
					for key, val := range _rla_access_data {
						//fmt.Println(key, val.([]map[string]any))
						rla_tables_ids[key] = app.getRLAIds(val.([]map[string]any), key, "read", []any{})
					}
				} else {
					fmt.Println("DEBUG THIS SOMETHING WORONG WITH RLA(READ):", rla_access)
				}
				for _, t := range _tables_to_chk {
					if _, ok := rla_tables_ids[t.(string)].([]any); !ok {
						rla_tables_ids[t.(string)] = []any{}
					}
				}
				//fmt.Println(fk_tables_pk, rla_tables_ids)
				if _, ok := params["data"].(map[string]any)["filters"]; !ok {
					params["data"].(map[string]any)["filters"] = []any{}
				}
				fk_tables_pk[table] = pk
				_filters := params["data"].(map[string]any)["filters"].([]any)
				for rla_t, ids := range rla_tables_ids {
					/*/fmt.Println(rla_t, fk_tables_pk[rla_t])
					if rla_t != table {
						_aux := map[string]any{"field": fk_tables_pk[rla_t], "cond": "IN", "value": app.joinSlice(ids.([]any), ",")}
						_filters = append(_filters, _aux)
					} else {*/
					_aux := map[string]any{
						"field":     fk_tables_pk[rla_t],
						"cond":      "IN",
						"value":     app.joinSlice(ids.([]any), ","),
						"glue_cond": "OR",
						"field2":    "user_id",
						"cond2":     "=",
						"value2":    user_id,
					}
					_filters = append(_filters, _aux)
					//}
				}
				//fmt.Println("RLA Filters Added:", _filters)
				params["data"].(map[string]any)["filters"] = _filters
			}
		}
	}
	// FILTERS
	queryParams := []any{}
	filters := []any{}
	if _, ok := _schema["fields"].(map[string]any); !ok {
	} else if _, ok := _schema["fields"].(map[string]any)["excluded"]; ok {
		filters = []any{fmt.Sprintf(`"%s"."excluded" IS FALSE`, table)}
	}
	if _, ok := params["data"].(map[string]any)["filters"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["filters"].([]any); ok {
		_filters := params["data"].(map[string]any)["filters"].([]any)
		for _, filter := range _filters {
			_field := filter.(map[string]any)["field"].(string)
			_field2 := ""
			if _, ok := filter.(map[string]any)["field2"]; ok {
				_field2 = filter.(map[string]any)["field2"].(string)
			}
			_cond := "="
			if _, ok := filter.(map[string]any)["cond"]; ok {
				_cond = filter.(map[string]any)["cond"].(string)
			}
			_glue_cond := "OR"
			if _, ok := filter.(map[string]any)["glue_cond"]; ok {
				_glue_cond = filter.(map[string]any)["glue_cond"].(string)
			}
			_cond2 := "="
			if _, ok := filter.(map[string]any)["cond2"]; ok {
				_cond2 = filter.(map[string]any)["cond2"].(string)
			}
			var _value any
			if _, ok := filter.(map[string]any)["value"]; ok {
				_value = filter.(map[string]any)["value"]
			}
			var _value2 any
			if _, ok := filter.(map[string]any)["value2"]; ok {
				_value2 = filter.(map[string]any)["value2"]
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
						// fmt.Println(_tbl, _tbl_fields)
						is_in_fk_fields = true
						_table = _tbl
						// fmt.Println("FILTER FIELD IN FK TABLE:", _table, _field)
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
				_aux2 := ""
				if _glue_cond != "" && _cond2 != "" && _field2 != "" && _value2 != "" {
					_aux2 = fmt.Sprintf(` %s "%s"."%s" %s ?`, _glue_cond, _table, _field2, _cond2)
				}
				filters = append(filters, fmt.Sprintf(`("%s"."%s" %s (?)%s)`, _table, _field, _cond, _aux2))
				queryParams = append(queryParams, strings.Split(_value.(string), ","))
				if _aux2 != "" && _value2 != "" && strings.ToLower(_cond2) != "in" {
					queryParams = append(queryParams, _value2)
				} else if _aux2 != "" && _value2 != "" && strings.ToLower(_cond2) == "in" {
					queryParams = append(queryParams, strings.Split(_value2.(string), ","))
				}
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
	query := fmt.Sprintf(`SELECT %s %s FROM %s"%s"`, distinct, app.joinSlice(_flds, ","), schm, table)
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
	query_total := fmt.Sprintf(`SELECT COUNT(*) AS "n_rows" FROM (%s) AS "T"`, query)
	if len(orderBy) > 0 {
		query = fmt.Sprintf(`%s ORDER BY %s`, query, app.joinSlice(orderBy, ", "))
	}
	if limit != -1 {
		query = fmt.Sprintf(`%s LIMIT %d OFFSET %d`, query, limit, offset)
	}
	// INTERCEPT READ QUERY
	_, database, _ := app.GetDBNameFromParams(params)
	//crud_aciton, table
	crud_aciton := "read"
	intercept_data := []any{database, table}
	get_crud_intercepts_sql := fmt.Sprintf(`SELECT * FROM "crud_intercept" WHERE "intercept_type_id" = 1 AND "active" = TRUE AND "db" = ? AND "table" = ? AND "%s" IS TRUE ORDER BY "db", "table", "order", "crud_intercept_id"`, crud_aciton)
	crud_intercept_rows, err := app.AdminGetRowsByFilter(get_crud_intercepts_sql, intercept_data)
	if err != nil {
		fmt.Printf("Error occurred while fetching crud_intercepts: %v", err)
	} else if len(crud_intercept_rows) > 0 {
		params["loc"] = loc
		params["table"] = table
		params["database"] = database
		params["pk"] = pk
		params["id"] = nil
		params["crud_aciton"] = crud_aciton
		params["user_id"] = user_id
		query_org := query
		_data := Dict{
			"read_sql": query,
			"_user":    params["user"],
		}
		for _, c_intercept := range crud_intercept_rows {
			ires, err := app.RunCrudIntercept(params, c_intercept, _data)
			if err != nil {
				return map[string]any{
					"success": false,
					"msg":     fmt.Sprintf("Error runing the interception: %s -> %s\n", c_intercept["crud_intercept_code"], err.Error()),
				}
			}
			if q, ok := ires["output"].(string); ok {
				query = q
				_data["read_sql"] = query
				fmt.Println("QUERY INTERCEPT:", query_org, query)
			}
		}
	}
	//fmt.Println(query)
	query, args, err := sqlx.In(query, queryParams...)
	if err != nil {
		println("Error geting the table query:", err)
	}
	sqlOnly := false
	if _, ok := params["data"].(map[string]any)["sql_only"]; !ok {
	} else if _, ok := params["data"].(map[string]any)["sql_only"].(bool); ok {
		sqlOnly = params["data"].(map[string]any)["sql_only"].(bool)
	}
	if sqlOnly {
		//fmt.Println("SQL ONLY!")
		msg, _ := app.i18n.T("success", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
			"data":    Dict{},
			"sql":     query,
			"args":    args,
		}
	}
	if table == "menu_table" {
		// fmt.Printf("%s: %s -> %s %v", table, join, query, args)
	}
	results := make([]map[string]any, 0)
	data, _, err := db.QueryMultiRows(query, args...)
	if err != nil {
		fmt.Println("READ ERR:", args, query, err)
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
			"sql":     query,
		}
	} else if /* *data == nil || */ len(*data) == 0 {
	} else {
		results = *data
	}
	total := 0
	//fmt.Println(query_total)
	trows, _, err := db.QuerySingleRow(query_total, args...)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	total = app.toInt((*trows)["n_rows"])
	// fmt.Println(table, user_id, pk, args, total, query)
	//data := map[string]any{}
	msg, _ := app.i18n.T("success", map[string]any{})
	res := map[string]any{
		"success": true,
		"msg":     msg,
		"data":    results,
		"total":   total,
		"cols":    _schema["fields_order"],
		// "schema":      _schema,
		"permissions": _permissions,
		//"row_level_tables": _row_level_tables,
		"sql": query,
	}
	if include_schema, ok := params["data"].(Dict)["include_schema"].(bool); ok && include_schema {
		res["schema"] = _schema
	}
	return res
}
