package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/etlx"
)

func (app *application) etlxMdParse(params Dict) Dict {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{
			"success": true,
			"msg":     msg,
		}
	}
	config := make(Dict)
	etlxlib := &etlx.ETLX{Config: config}
	_data, ok := params["data"].(Dict)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "Check the data passed, possible mal-formated!",
		}
	}
	_conf, ok := _data["conf"].(string)
	if !ok {
		return Dict{
			"success": false,
			"msg":     "Please validate the configutration, should be mardown string!",
		}
	}
	err := etlxlib.ConfigFromMDText(_conf)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%v", err),
		}
	}
	if _, ok := etlxlib.Config["REQUIRES"]; ok {
		_logs, err := etlxlib.LoadREQUIRES(nil)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("REQUIRES ERR: %v", err),
				"logs":    _logs,
			}
		}
	}
	// Print the parsed configuration
	if os.Getenv("ETLX_DEBUG_QUERY") == "true" {
		etlxlib.PrintConfigAsJSON(etlxlib.Config)
	}
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
		"data":    etlxlib.Config,
	}
}

func anyToStrings(input []any) []string {
	result := make([]string, 0, len(input))
	for _, v := range input {
		if str, ok := v.(string); ok {
			result = append(result, str)
		} else {
			result = append(result, fmt.Sprintf("%v", v)) // Convert non-string values to string
		}
	}
	return result
}

func (app *application) etlxRun(params Dict, ignore bool) Dict {
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
	if env.GetBool("ETLX_VALIDATE_ETLX_ACCESS", false) && !ignore {
		x := app.getEtlxByID(params)
		if _, ok := x["success"]; !ok {
			return x
		} else if !x["success"].(bool) {
			return x
		} else if len(x["data"].([]Dict)) == 0 {
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("ETLX ID %s does not exists or you don have access to it!", x["etlx_id"]),
			}
		}
	}
	config := make(Dict)
	etlxlib := &etlx.ETLX{Config: config}
	config, ok = _data["conf"].(Dict)
	if !ok {
		_conf, ok := _data["conf"].(string)
		if !ok {
			return Dict{
				"success": false,
				"msg":     "Please validate the configutration, should be mardown string!",
			}
		}
		err := etlxlib.ConfigFromMDText(_conf)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     fmt.Sprintf("%v", err),
			}
		}
	} else {
		etlxlib.Config = config
	}
	if _, ok := etlxlib.Config["REQUIRES"]; ok {
		_logs, err := etlxlib.LoadREQUIRES(nil)
		if err != nil {
			fmt.Printf("REQUIRES ERR: %v %v", err, _logs)
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
	// EXTRA CONFIG
	extraConf := Dict{}
	if ok {
		extraConf = Dict{
			"clean": false,
			"drop":  false,
			"rows":  false,
		}
		if clean, ok := _data["clean"].(bool); ok {
			extraConf["clean"] = clean
		}
		if drop, ok := _data["drop"].(bool); ok {
			extraConf["drop"] = drop
		}
		if rows, ok := _data["rows"].(bool); ok {
			extraConf["rows"] = rows
		}
		if file, ok := _data["file"].(string); ok {
			extraConf["file"] = file
		}
		if only, ok := _data["clean"].(string); ok {
			extraConf["only"] = strings.Split(only, ",")
		} else if only, ok := _data["only"].([]string); ok {
			extraConf["only"] = only
		} else if only, ok := _data["only"].([]any); ok {
			extraConf["only"] = anyToStrings(only)
		}
		if skip, ok := _data["skip"].(string); ok {
			extraConf["skip"] = strings.Split(skip, ",")
		} else if skip, ok := _data["skip"].([]string); ok {
			extraConf["skip"] = skip
		} else if skip, ok := _data["skip"].([]any); ok {
			extraConf["skip"] = anyToStrings(skip)
		}
		if steps, ok := _data["steps"].(string); ok {
			extraConf["steps"] = strings.Split(steps, ",")
		} else if steps, ok := _data["steps"].([]string); ok {
			extraConf["steps"] = steps
		} else if steps, ok := _data["steps"].([]any); ok {
			extraConf["steps"] = anyToStrings(steps)
		}
	}
	//fmt.Println("extraConf:", extraConf)
	logs := []Dict{}
	data := Dict{}
	_keys := []any{"NOTIFY", "LOGS", "SCRIPTS", "MULTI_QUERIES", "EXPORTS", "DATA_QUALITY", "ETL", "ELT", "ACTIONS", "AUTO_LOGS", "REQUIRES"}
	__order, ok := etlxlib.Config["__order"].([]string)
	hasOrderedKeys := false
	if !ok {
		__order2, ok := etlxlib.Config["__order"].([]any)
		if ok {
			hasOrderedKeys = true
			__order = []string{}
			for _, key := range __order2 {
				__order = append(__order, key.(string))
			}
		}
	} else {
		etlxlib.Config["__order"] = []any{}
		for key := range etlxlib.Config {
			etlxlib.Config["__order"] = append(etlxlib.Config["__order"].([]any), key)
		}
		//hasOrderedKeys = true
	}
	// fmt.Println("LEVEL 1 H:", __order, len(__order))
	if !hasOrderedKeys {
	} else if len(__order) > 0 {
		//fmt.Print("LEVEL 1 H:", __order)
		for _, key := range __order {
			if key == "metadata" || key == "__order" || key == "order" {
				continue
			}
			//if !app.contains(_keys, any(key)) {
			_key_conf, ok := etlxlib.Config[key].(Dict)
			if !ok {
				continue
			}
			_key_conf_metadata, ok := _key_conf["metadata"].(Dict)
			if !ok {
				continue
			}
			if _, ok := _key_conf_metadata["runs_as"]; !ok {
				_key_conf_metadata["runs_as"] = strings.ToUpper(key)
			}
			if runs_as, ok := _key_conf_metadata["runs_as"]; ok {
				fmt.Printf("%s RUN AS %s:\n", key, runs_as)
				if app.contains(_keys, runs_as) {
					switch runs_as {
					case "ETL", "ELT":
						_logs, err := etlxlib.RunETL(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "DATA_QUALITY", "DATAQUALITY", "QUALITY":
						_logs, err := etlxlib.RunDATA_QUALITY(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "MULTI_QUERIES", "STACKED_QUERIES":
						_logs, _, err := etlxlib.RunMULTI_QUERIES(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "EXPORTS":
						_logs, err := etlxlib.RunEXPORTS(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "NOTIFY", "NOTIFICATION":
						_logs, err := etlxlib.RunNOTIFY(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
						}
					case "ACTIONS":
						_logs, err := etlxlib.RunACTIONS(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "SCRIPTS":
						_logs, err := etlxlib.RunSCRIPTS(dateRef, nil, extraConf, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "LOGS", "OBSERVABILITY":
						_logs, err := etlxlib.RunLOGS(dateRef, nil, logs, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
							data[key] = Dict{
								"success": true,
								"runs_as": runs_as,
								"logs":    _logs,
							}
						}
					case "REQUIRES", "IMPORTS":
						_logs, err := etlxlib.LoadREQUIRES(nil, key)
						if err != nil {
							fmt.Printf("%s AS %s ERR: %v\n", key, runs_as, err)
						} else {
							if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
								_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
								if err != nil {
									fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
								}
							}
							logs = append(logs, _logs...)
						}
					default:
						//
					}
				}
			}
			//}
		}
	}
	msg, _ := app.i18n.T("success", Dict{})
	//fmt.Println(msg)
	return Dict{
		"success": true,
		"msg":     msg,
		"logs":    logs,
		"data":    data,
	}
}

func (app *application) etlxParseRun(params Dict) Dict {
	res := app.etlxMdParse(params)
	if res["success"].(bool) {
		params["data"].(Dict)["conf"] = res["data"]
		return app.etlxRun(params, false)
	}
	return res
}

func (app *application) getEtlxByID(params Dict) Dict {
	table := "etlx"
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
	etlx_id := any(nil)
	if _, ok := params["data"].(Dict)["etlx_id"]; ok {
		etlx_id = params["data"].(Dict)["etlx_id"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		etlx_id = params["data"].(Dict)["data"].(Dict)["etlx_id"]
	} else if _, ok := params["data"].(Dict)["etlx"].(Dict); ok {
		etlx_id = params["data"].(Dict)["etlx"].(Dict)["etlx_id"]
	}
	if etlx_id == nil || etlx_id == any(nil) {
		return Dict{
			"success": false,
			"msg":     "No ETLX ID found!",
		}
	}
	_aux_params := params
	_aux_params["data"].(Dict)["table"] = table
	_aux_params["data"].(Dict)["db"] = database
	_aux_params["data"].(Dict)["limit"] = any(1.0)
	_aux_params["data"].(Dict)["offset"] = any(0.0)
	_aux_params["data"].(Dict)["filters"] = []any{Dict{
		"field": "etlx_id",
		"cond":  "=",
		"value": etlx_id,
	}}
	res := app.read(_aux_params)
	res["etlx_id"] = etlx_id
	return res
}

func (app *application) etlxRunByName(params Dict) Dict {
	name, ok := params["name"].(string)
	if !ok {
		name, _ = params["data"].(Dict)["name"].(string)
	}
	table := "etlx"
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
	_aux_params["data"].(Dict)["filters"] = []any{Dict{
		"field": "etl",
		"cond":  "=",
		"value": name,
	}}
	etlx_get_conf := app.read(_aux_params)
	//fmt.Println(len(etlx_get_conf["data"].([]Dict)), etlx_get_conf["data"])
	if _, ok := etlx_get_conf["success"]; !ok {
		return etlx_get_conf
	} else if !etlx_get_conf["success"].(bool) {
		return etlx_get_conf
	} else if len(etlx_get_conf["data"].([]Dict)) == 0 {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("ETL %s does not exists", name),
		}
	}
	params["data"].(Dict)["conf"] = etlx_get_conf["data"].([]Dict)[0]["etlx_conf"]
	res := app.etlxParseRun(params)
	if res["success"].(bool) {
		params["data"].(Dict)["conf"] = res["data"]
		return app.etlxRun(params, true)
	}
	return res
}
