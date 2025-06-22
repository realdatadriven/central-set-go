package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/realdatadriven/etlx"
)

func (app *application) etlxMdParse(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	config := make(map[string]any)
	etlxlib := &etlx.ETLX{Config: config}
	_data, ok := params["data"].(map[string]any)
	if !ok {
		return map[string]interface{}{
			"success": false,
			"msg":     "Check the data passed, possible mal-formated!",
		}
	}
	_conf, ok := _data["conf"].(string)
	if !ok {
		return map[string]interface{}{
			"success": false,
			"msg":     "Please validate the configutration, should be mardown string!",
		}
	}
	err := etlxlib.ConfigFromMDText(_conf)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%v", err),
		}
	}
	if _, ok := etlxlib.Config["REQUIRES"]; ok {
		_logs, err := etlxlib.LoadREQUIRES(nil)
		if err != nil {
			return map[string]interface{}{
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
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
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

func (app *application) etlxRun(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	_data, ok := params["data"].(map[string]any)
	if !ok {
		return map[string]interface{}{
			"success": false,
			"msg":     "Check the data passed, possible mal-formated!",
		}
	}
	config := make(map[string]any)
	etlxlib := &etlx.ETLX{Config: config}
	config, ok = _data["conf"].(map[string]any)
	if !ok {
		_conf, ok := _data["conf"].(string)
		if !ok {
			return map[string]interface{}{
				"success": false,
				"msg":     "Please validate the configutration, should be mardown string!",
			}
		}
		err := etlxlib.ConfigFromMDText(_conf)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"msg":     fmt.Sprintf("%v", err),
			}
		}
	} else {
		etlxlib.Config = config
	}
	// DATE REF
	var _dateRef interface{}
	if _, ok := _data["date_ref"]; ok {
		_dateRef = _data["date_ref"]
	}
	var dateRef []time.Time
	switch _type := _dateRef.(type) {
	case string:
		_dt, _ := time.Parse("2006-01-02", _dateRef.(string))
		dateRef = append(dateRef, _dt)
	case []interface{}:
		for _, _dt := range _dateRef.([]interface{}) {
			_dt, _ := time.Parse("2006-01-02", _dt.(string))
			dateRef = append(dateRef, _dt)
		}
	default:
		fmt.Println("Unable to parse date ref: ", _type, _dateRef)
	}
	// EXTRA CONFIG
	extraConf := map[string]any{}
	if ok {
		extraConf = map[string]any{
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
	logs := []map[string]any{}
	data := map[string]any{}
	// RUN ETL
	if _, ok := etlxlib.Config["ETL"]; ok {
		_logs, err := etlxlib.RunETL(dateRef, nil, extraConf)
		if err != nil {
			data["ETL"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("%v", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["ETL"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	// DATA_QUALITY
	if _, ok := etlxlib.Config["DATA_QUALITY"]; ok {
		_logs, err := etlxlib.RunDATA_QUALITY(dateRef, nil, extraConf)
		if err != nil {
			data["DATA_QUALITY"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("DATA_QUALITY ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["DATA_QUALITY"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	// EXPORTS
	if _, ok := etlxlib.Config["EXPORTS"]; ok {
		_logs, err := etlxlib.RunEXPORTS(dateRef, nil, extraConf)
		if err != nil {
			data["EXPORTS"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("EXPORTS ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["EXPORTS"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	// MULTI_QUERIES
	if _, ok := etlxlib.Config["MULTI_QUERIES"]; ok {
		_logs, _data, err := etlxlib.RunMULTI_QUERIES(dateRef, nil, extraConf)
		if err != nil {
			data["MULTI_QUERIES"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("MULTI_QUERIES ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["MULTI_QUERIES"] = map[string]any{
				"success": true,
				"data":    _data,
				"logs":    _logs,
			}
		}
	}
	// SCRIPTS
	if _, ok := etlxlib.Config["SCRIPTS"]; ok {
		_logs, err := etlxlib.RunSCRIPTS(dateRef, nil, extraConf)
		if err != nil {
			data["SCRIPTS"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("SCRIPTS ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["SCRIPTS"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	// ACTIONS
	if _, ok := etlxlib.Config["ACTIONS"]; ok {
		_logs, err := etlxlib.RunACTIONS(dateRef, nil, extraConf)
		if err != nil {
			data["ACTIONS"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("ACTIONS ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["ACTIONS"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	// LOGS
	if _, ok := etlxlib.Config["LOGS"]; ok {
		_logs, err := etlxlib.RunLOGS(dateRef, nil, logs)
		if err != nil {
			data["LOGS"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("LOGS ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
		}
	}
	// NOTIFY
	if _, ok := etlxlib.Config["NOTIFY"]; ok {
		_logs, err := etlxlib.RunNOTIFY(dateRef, nil, extraConf)
		if err != nil {
			data["NOTIFY"] = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("NOTIFY ERR: %v!", err),
			}
		} else {
			// LOGS
			if _, ok := etlxlib.Config["AUTO_LOGS"]; ok && len(_logs) > 0 {
				_, err := etlxlib.RunLOGS(dateRef, nil, _logs, "AUTO_LOGS")
				if err != nil {
					fmt.Printf("INCREMENTAL AUTOLOGS ERR: %v\n", err)
				}
			}
			logs = append(logs, _logs...)
			data["NOTIFY"] = map[string]any{
				"success": true,
				"logs":    _logs,
			}
		}
	}
	_keys := []any{"NOTIFY", "LOGS", "SCRIPTS", "MULTI_QUERIES", "EXPORTS", "DATA_QUALITY", "ETL", "ACTIONS", "AUTO_LOGS"}
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
		hasOrderedKeys = true
	}
	// fmt.Println("LEVEL 1 H:", __order, len(__order))
	if !hasOrderedKeys {
	} else if len(__order) > 0 {
		//fmt.Print("LEVEL 1 H:", __order)
		for _, key := range __order {
			if !app.contains(_keys, any(key)) {
				_key_conf, ok := etlxlib.Config[key].(map[string]any)
				if !ok {
					continue
				}
				_key_conf_metadata, ok := _key_conf["metadata"].(map[string]any)
				if !ok {
					continue
				}
				if runs_as, ok := _key_conf_metadata["runs_as"]; ok {
					fmt.Printf("%s RUN AS %s:\n", key, runs_as)
					if app.contains(_keys, runs_as) {
						switch runs_as {
						case "ETL":
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
								data[key] = map[string]any{
									"success": true,
									"runs_as": runs_as,
									"logs":    _logs,
								}
							}
						case "DATA_QUALITY":
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
								data[key] = map[string]any{
									"success": true,
									"runs_as": runs_as,
									"logs":    _logs,
								}
							}
						case "MULTI_QUERIES":
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
								data[key] = map[string]any{
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
								data[key] = map[string]any{
									"success": true,
									"runs_as": runs_as,
									"logs":    _logs,
								}
							}
						case "NOTIFY":
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
								data[key] = map[string]any{
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
								data[key] = map[string]any{
									"success": true,
									"runs_as": runs_as,
									"logs":    _logs,
								}
							}
						case "LOGS":
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
								data[key] = map[string]any{
									"success": true,
									"runs_as": runs_as,
									"logs":    _logs,
								}
							}
						default:
							//
						}
					}
				}
			}
		}
	}
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"logs":    logs,
		"data":    data,
	}
}

func (app *application) etlxParseRun(params map[string]interface{}) map[string]interface{} {
	res := app.etlxMdParse(params)
	if res["success"].(bool) {
		params["data"].(map[string]any)["conf"] = res["data"]
		return app.etlxRun(params)
	}
	return res
}
