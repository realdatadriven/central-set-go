package main

import (
	"fmt"

	"github.com/realdatadriven/etlx"
)

func (app *application) Buckup(params map[string]any) map[string]any {
	dsn, _, _ := app.GetDBNameFromParams(map[string]any{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return map[string]any{
			"success": false,
			"msg":     fmt.Errorf("error geting the db connection: %w", err),
		}
	}
	defer db.Close()
	sql := `select * from "app" where "app" like ? excluded = false`
	_app := "%"
	if _, ok := params["data"].(map[string]any)["name"].(string); ok {
		_app = params["data"].(map[string]any)["name"].(string)
	}
	apps, _, err := db.QueryMultiRows(sql, []any{_app}...)
	if err != nil {
		return map[string]any{
			"success": true,
			"msg":     fmt.Errorf("error geting the apps: %w", err),
		}
	}
	for _, _app := range *apps {
		fmt.Printf("1: %v\n", _app)
	}
	msg, _ := app.i18n.T("success", map[string]any{})
	data := map[string]any{
		"success": true,
		"msg":     msg,
	}
	return data
}
