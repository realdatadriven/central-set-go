package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/etlx"
)

func DeepCompare(input, config map[string]any, path string) (bool, string) {
	ignoreKeys := map[string]bool{
		"file":           true,
		"date":           true,
		"date_ref":       true,
		"loading":        true,
		"open":           true,
		"__order":        true,
		"__lakes":        true,
		"num_rows":       true,
		"skip_extract":   true,
		"skip_load":      true,
		"skip_transform": true,
		"tmp":            true,
		"temp":           true,
	}
	// Check all keys in input against config
	for key, inputVal := range input {
		if ignoreKeys[key] {
			continue
		}
		fullPath := key
		if path != "" {
			fullPath = path + "." + key
		}
		configVal, exists := config[key]
		if !exists {
			// fmt.Println(key, inputVal, configVal)
			return false, fmt.Sprintf("Key %s missing in config", fullPath)
		}
		// Compare based on type
		switch iv := inputVal.(type) {
		case string:
			if cv, ok := configVal.(string); !ok || iv != cv {
				return false, fmt.Sprintf("Mismatch at %s: input '%v' != config '%v'", fullPath, iv, configVal)
			}
		case bool:
			/*if cv, ok := configVal.(bool); !ok || iv != cv {
				return false, fmt.Sprintf("Mismatch at %s: input %v != config %v", fullPath, iv, configVal)
			}*/
			continue
		case map[string]any:
			if cv, ok := configVal.(map[string]any); !ok {
				return false, fmt.Sprintf("Type mismatch at %s: expected map", fullPath)
			} else {
				equal, msg := DeepCompare(iv, cv, fullPath)
				if !equal {
					return false, msg
				}
			}
		case []any:
			if cv, ok := configVal.([]any); !ok {
				return false, fmt.Sprintf("Type mismatch at %s: expected array", fullPath)
			} else if len(iv) != len(cv) {
				return false, fmt.Sprintf("Array length mismatch at %s: input %d != config %d", fullPath, len(iv), len(cv))
			} else {
				for i := range iv {
					arrayPath := fmt.Sprintf("%s[%d]", fullPath, i)
					equal, msg := deepCompareAny(iv[i], cv[i], arrayPath)
					if !equal {
						return false, msg
					}
				}
			}
		default:
			// For other types, use reflect.DeepEqual as fallback
			if !reflect.DeepEqual(inputVal, configVal) {
				return false, fmt.Sprintf("Mismatch at %s: input %v != config %v", fullPath, inputVal, configVal)
			}
		}
	}
	/*/ Check for extra keys in config that are not in input (ignoring ignored keys)
	for key := range config {
		if ignoreKeys[key] {
			continue
		}
		if _, exists := input[key]; !exists {
			fullPath := key
			if path != "" {
				fullPath = path + "." + key
			}
			return false, fmt.Sprintf("Extra key %s in config", fullPath)
		}
	}*/
	return true, ""
}

// Helper function to compare any types (for array elements)
func deepCompareAny(a, b any, path string) (bool, string) {
	switch av := a.(type) {
	case string:
		if bv, ok := b.(string); !ok || av != bv {
			//fmt.Println("rdeepCompareAny string:", a, b)
			return false, fmt.Sprintf("Mismatch at %s: input '%v' != config '%v'", path, av, b)
		}
	case bool:
		if bv, ok := b.(bool); !ok || av != bv {
			return false, fmt.Sprintf("Mismatch at %s: input %v != config %v", path, av, b)
		}
	case map[string]any:
		if bv, ok := b.(map[string]any); !ok {
			return false, fmt.Sprintf("Type mismatch at %s: expected map", path)
		} else {
			return DeepCompare(av, bv, path)
		}
	case []any:
		if bv, ok := b.([]any); !ok {
			return false, fmt.Sprintf("Type mismatch at %s: expected array", path)
		} else if len(av) != len(bv) {
			return false, fmt.Sprintf("Array length mismatch at %s: input %d != config %d", path, len(av), len(bv))
		} else {
			for i := range av {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				equal, msg := deepCompareAny(av[i], bv[i], arrayPath)
				if !equal {
					return false, msg
				}
			}
		}
	default:
		// Fallback for other types
		if !reflect.DeepEqual(a, b) {
			//fmt.Println("reflect.DeepEqual:", a, b)
			return false, fmt.Sprintf("Mismatch at %s: input %v != config %v", path, a, b)
		}
	}
	return true, ""
}

func (app *application) etlxMdParse(params Dict) Dict {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{
			"success": true,
			"msg":     msg,
		}
	}
	_params := Dict{}
	if _, ok := params["params"].(Dict); ok {
		_params = params["params"].(Dict)
	}
	config := make(Dict)
	etlxlib := &etlx.ETLX{Config: config, Params: _params}
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
		// etlxlib.PrintConfigAsJSON(etlxlib.Config)
	}
	mdData, _ := etlxlib.QueryETLXMD("")
	nodes, ok := mdData["nodes"]
	if !ok {
		fmt.Println("No nodes data found")
	}
	edges, ok := mdData["edges"]
	if !ok {
		fmt.Println("No edges data found")
	}
	if len(nodes) == 0 {
		nodes, ok = mdData["nodes_est"]
		if !ok {
			fmt.Println("No nodes data found")
		}
	}
	if len(edges) == 0 {
		edges, ok = mdData["edges_est"]
		if !ok {
			fmt.Println("No edges data found")
		}
	}
	flow := etlxlib.GenerateMermaidFlowchart(nodes, edges)
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
		"data":    etlxlib.Config,
		"mdData":  mdData,
		"flow":    flow,
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

// run query using INSTALL markdown FROM community on the md config
func (app *application) queryETLXMD(params Dict) Dict {
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
	_conf, ok := x["data"].([]Dict)[0]["etlx_conf"].(string)
	if !ok {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("ETLX ID %s does not have configuration!", x["etlx_id"]),
		}
	}
	_params := Dict{}
	if _, ok := params["params"].(Dict); ok {
		_params = params["params"].(Dict)
	}
	config := make(Dict)
	etlxlib := &etlx.ETLX{Config: config, Params: _params}
	err := etlxlib.ConfigFromMDText(_conf)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("%v", err),
		}
	}
	res, err := etlxlib.QueryETLXMD("")
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": err == nil,
		"msg":     msg,
		"data":    res,
	}
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
	_params := Dict{}
	if _, ok := params["params"].(Dict); ok {
		_params = params["params"].(Dict)
	}
	var loc *time.Location
	if _, ok := params["location"].(*time.Location); ok {
		loc = params["location"].(*time.Location)
	} else {
		loc = time.Local
	}
	// CONFIG
	config := make(Dict)
	etlxlib := &etlx.ETLX{Config: config, Params: _params, TimeZone: loc}
	etlxlib.MetadataOrder = false
	if order_metadata, ok := _data["order_metadata"].(bool); ok {
		fmt.Println("order_metadata:", order_metadata)
		etlxlib.MetadataOrder = order_metadata
	}
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
	// VALIDATE
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
		if env.GetBool("ETLX_VALIDATE_CLI_CONF_WITH_DB", true) {
			// CONF
			_conf, _ := x["data"].([]Dict)[0]["etlx_conf"]
			res := app.etlxMdParse(Dict{"data": Dict{"conf": _conf}})
			if !res["success"].(bool) {
				return res
			}
			db_cnf, ok := res["data"].(Dict)
			if !ok {
				return Dict{
					"success": false,
					"msg":     "Unable to parse the database config!",
				}
			}
			//etlxlib.PrintConfigAsJSON(config)
			//etlxlib.PrintConfigAsJSON(db_cnf)
			// VALIDATE
			equal, msg := DeepCompare(config, db_cnf, "")
			if !equal {
				fmt.Println("DeepCompare:", msg)
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
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
	logs, data, err := etlxlib.RunETLX(extraConf, dateRef)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
			"logs":    logs,
			"data":    data,
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
