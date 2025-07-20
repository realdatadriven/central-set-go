package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/realdatadriven/etlx"
)

func (app *application) nbRun(params Dict) Dict {
	return app.nbRunCells(params)
}

func (app *application) nbRunCells(params Dict) Dict {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{
			"success": true,
			"msg":     msg,
		}
	}
	_data, ok := params["data"].(Dict)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "Check the data passed, possible mal-formated!",
		}
	}
	// fmt.Println(_data)
	cells, ok := _data["cells"].([]any)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "No notebook cells identified!",
		}
	}
	// DATE REF
	var _dateRef any
	if _, ok := _data["date_ref"]; ok {
		_dateRef = _data["date_ref"]
	}
	var dateRef []time.Time
	switch _type := _dateRef.(type) {
	case string:
		_dt, _ := time.Parse("2006-01-02", _dateRef.(string))
		dateRef = append(dateRef, _dt)
	case []any:
		for _, _dt := range _dateRef.([]any) {
			_dt, _ := time.Parse("2006-01-02", _dt.(string))
			dateRef = append(dateRef, _dt)
		}
	default:
		dateRef = append(dateRef, time.Now().AddDate(0, 0, -1))
		fmt.Println("Unable to parse date ref: ", _type, _dateRef)
	}
	/*/ EXTRA CONFIG
	extraConf := Dict{}
	if ok {
		extraConf = Dict{
			"clean": false,
		}
	}*/
	ddbInstance, _ := etlx.GetDB("duckdb:")
	defer ddbInstance.Close()
	etlxlib := &etlx.ETLX{}
	data := Dict{}
	msg, _ := app.i18n.T("success", Dict{})
	start_gbl := time.Now()
	for _, _cell := range cells {
		cell := _cell.(Dict)
		_sql := cell["code"].(string)
		_sql = etlxlib.ReplaceEnvVariable(_sql)
		_sql = etlxlib.ReplaceQueryStringDate(_sql, dateRef)
		_id := cell["id"].(string)
		data[_id] = Dict{}
		_reg_ex := `CREATE.*TABLE|UPDATE.*TABLE|DROP.*|INSERT.*INTO|DELETE|ALTER.*TABLE|UPSERT.*|ATTACH.*|DETACH.*|SECRET.*`
		patt := regexp.MustCompile(_reg_ex)
		_match := patt.FindAllString(_sql, -1)
		patt2 := regexp.MustCompile("SELECT|FROM")
		_match2 := patt2.FindAllString(_sql, -1)
		start := time.Now()
		if len(_match) > 0 {
			_, err := ddbInstance.ExecuteQuery(_sql, []any{}...)
			//msg, _ := app.i18n.T("query-not-allowed", Dict{"query": _sql, "match": app.joinSlice(app.sliceStrs2SliceInterfaces(_match), ";")})
			if err != nil {
				data[_id] = Dict{
					"success": false,
					"msg":     fmt.Sprintf("%s", err),
					"start":   start,
					"end":     time.Now(),
				}
			} else {
				data[_id] = Dict{
					"success": true,
					"start":   start,
					"end":     time.Now(),
					"msg":     msg,
				}
			}
		} else if len(_match2) > 0 {
			query_n_rows := fmt.Sprintf(`SELECT COUNT(*) AS "n_rows" FROM (%s) AS "T"`, _sql)
			patt = regexp.MustCompile(`LIMIT`)
			_match = patt.FindAllString(_sql, -1)
			if len(_match) == 0 {
				limit := 10
				if _, ok := params["data"].(Dict)["limit"]; ok {
					limit = int(params["data"].(Dict)["limit"].(float64))
				}
				offset := 0
				if _, ok := params["data"].(Dict)["offset"]; ok {
					offset = int(params["data"].(Dict)["offset"].(float64))
				}
				if limit != -1 {
					_sql = fmt.Sprintf(`%s LIMIT %d OFFSET %d`, _sql, limit, offset)
				}
			}
			re := regexp.MustCompile(`;+$`)
			_sql := re.ReplaceAllString(_sql, "")
			//fmt.Println(_sql)
			results, cols, _, err := ddbInstance.QueryMultiRowsWithCols(_sql, []any{}...)
			if err != nil {
				data[_id] = Dict{
					"success": false,
					"msg":     fmt.Sprintf("%s", err),
					"start":   start,
					"end":     time.Now(),
				}
			} else {
				n_rows, _, err := ddbInstance.QuerySingleRow(query_n_rows, []any{}...)
				if err != nil {
					data[_id] = Dict{
						"success": false,
						"msg":     fmt.Sprintf("returning query total rows: %s", err),
						"start":   start,
						"end":     time.Now(),
					}
				} else {
					total := 0
					total = int((*n_rows)["n_rows"].(int64))
					data[_id] = Dict{
						"success": true,
						"msg":     msg,
						"data":    *results,
						"nrows":   total,
						"columns": cols,
						"start":   start,
						"end":     time.Now(),
					}
				}
			}
		} else {
			data[_id] = Dict{
				"success": false,
				"msg":     "This query is not in the allowed syntax list, can you check with the system administrator?",
				"start":   start,
				"end":     time.Now(),
			}
		}
	}
	return Dict{
		"success": true,
		"msg":     msg,
		"data":    data,
		"start":   start_gbl,
		"end":     time.Now(),
	}
}

func (app *application) nbRunByName(params Dict) Dict {
	name, ok := params["name"].(string)
	if !ok {
		name, _ = params["data"].(Dict)["name"].(string)
	}
	table := "notebook"
	if _, ok := params["table"].(string); ok {
		table, _ = params["table"].(string)
	} else if _, ok := params["data"].(Dict)["table"].(string); ok {
		table, _ = params["data"].(Dict)["table"].(string)
	}
	database := "ETLX"
	if _, ok := params["db"].(string); ok {
		database, _ = params["db"].(string)
	} else if _, ok := params["data"].(Dict)["db"].(string); ok {
		database, _ = params["data"].(Dict)["db"].(string)
	} else if _, ok := params["database"].(string); ok {
		database, _ = params["database"].(string)
	} else if _, ok := params["data"].(Dict)["database"].(string); ok {
		database, _ = params["data"].(Dict)["database"].(string)
	}
	_aux_params := params
	_aux_params["data"].(Dict)["table"] = table
	_aux_params["data"].(Dict)["db"] = database
	_aux_params["data"].(Dict)["limit"] = any(1.0)
	_aux_params["data"].(Dict)["offset"] = any(0.0)
	_aux_params["data"].(Dict)["filters"] = []any{}
	_aux_params["data"].(Dict)["filters"] = append(
		_aux_params["data"].(Dict)["filters"].([]any),
		Dict{
			"field": "notebook",
			"cond":  "=",
			"value": name,
		},
	)
	_aux_params["data"].(Dict)["order_by"] = []any{}
	_aux_params["data"].(Dict)["order_by"] = append(
		_aux_params["data"].(Dict)["order_by"].([]any),
		Dict{
			"field": "notebook_id",
			"order": "desc",
		},
	)
	nb_get_conf := app.read(_aux_params)
	//fmt.Println(len(nb_get_conf["data"].([]Dict)), nb_get_conf["data"])
	if _, ok := nb_get_conf["success"]; !ok {
		return nb_get_conf
	} else if !nb_get_conf["success"].(bool) {
		return nb_get_conf
	} else if len(nb_get_conf["data"].([]Dict)) == 0 {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Notebook %s does not exists", name),
		}
	}
	var data Dict
	err := json.Unmarshal([]byte(nb_get_conf["data"].([]Dict)[0]["notebook_conf"].(string)), &data)
	if err != nil {
		panic(err)
	}
	params["data"] = data
	return app.nbRun(params)
}
