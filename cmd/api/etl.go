package main

import (
	"encoding/json"
	"fmt"
)

func (app *application) extract(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	conf_key := "etl_rbase_input_conf"
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, conf_key)
	return app.ETLExtract(params, _extra_conf, _item, _conf, _etlrb, _conf_etlrb, _step)
}

func (app *application) n_rows(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	conf_key := ""
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, conf_key)
	return app._n_rows(_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
}

func (app *application) delete(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	conf_key := ""
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, conf_key)
	//_n_rows := app._n_rows(_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
	_del_res := app._delete(_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
	//_del_res["n_rows"] = _n_rows["n_rows"]
	return _del_res
}

func (app *application) transform(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, "etl_rbase_output_conf")
	return app._transform(params, _item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
}

func (app *application) export(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, "etl_rbase_export_conf")
	return app._export(params, _item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
}

func (app *application) quality(params map[string]interface{}) map[string]interface{} {
	if app.IsEmpty(params["data"]) {
		msg, _ := app.i18n.T("no-data", map[string]interface{}{})
		return map[string]interface{}{
			"success": true,
			"msg":     msg,
		}
	}
	_item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step := app.etl_commons(params, "etl_rbase_quality_conf")
	return app._quality(params, _item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step)
}

func (app *application) etl_commons(params map[string]interface{}, conf_key string) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	_data := map[string]interface{}{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]interface{})
	}
	_item := map[string]interface{}{}
	if _, ok := _data["data"]; ok {
		_item = _data["data"].(map[string]interface{})
	}
	_step := map[string]interface{}{}
	if _, ok := _data["step"]; ok {
		_step = _data["step"].(map[string]interface{})
	}
	if conf_key == "" {
		if _, ok := _step["table"]; ok {
			conf_key = fmt.Sprintf(`%s_conf`, _step["table"])
		}
	}
	// fmt.Println("conf_key:", _step["table"], conf_key)
	_etlrb := map[string]interface{}{}
	_conf := map[string]interface{}{}
	_conf_etlrb := map[string]interface{}{}
	if _, ok := _item[conf_key]; ok {
		_aux_cnf := _item[conf_key]
		switch _type := _aux_cnf.(type) {
		case string:
			err := json.Unmarshal([]byte(_aux_cnf.(string)), &_conf)
			if err != nil {
				fmt.Println("conf err", _aux_cnf, err)
			}
		case map[string]interface{}:
			_conf = _aux_cnf.(map[string]interface{})
		default:
			fmt.Println(_type)
		}
	}
	if _, ok := _etlrb["etl_report_base_conf"]; ok {
		_aux_cnf := _etlrb["etl_report_base_conf"]
		switch _type := _aux_cnf.(type) {
		case string:
			err := json.Unmarshal([]byte(_aux_cnf.(string)), &_conf_etlrb)
			if err != nil {
				fmt.Println("conf err", _aux_cnf, err)
			}
		case map[string]interface{}:
			_conf_etlrb = _aux_cnf.(map[string]interface{})
		default:
			fmt.Println(_type)
		}
	}
	// DATABASE
	_extra_conf := map[string]interface{}{
		"driverName": app.config.db.driverName,
		"dsn":        app.config.db.dsn,
	}
	return _item, _conf, _etlrb, _conf_etlrb, _extra_conf, _step
}
