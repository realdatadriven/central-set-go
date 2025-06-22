package main

import (
	"fmt"
	"regexp"

	"github.com/realdatadriven/etlx"
)

func (app *application) CrudRunQuery(params map[string]interface{}, query string, db etlx.DBInterface) map[string]interface{} {
	patt := regexp.MustCompile(`CREATE.*TABLE|UPDATE.*TABLE|DROP.*|INSERT.*INTO|DELETE|ALTER.*TABLE|UPSERT.*`)
	_match := patt.FindAllString(query, -1)
	if len(_match) > 0 {
		msg, _ := app.i18n.T("query-not-allowed", map[string]interface{}{"query": query, "match": app.joinSlice(app.sliceStrs2SliceInterfaces(_match), ";")})
		return map[string]interface{}{
			"success": false,
			"msg":     msg,
		}
	}
	query_n_rows := fmt.Sprintf(`SELECT COUNT(*) AS "n_rows" FROM (%s) AS "T"`, query)
	patt = regexp.MustCompile(`LIMIT`)
	_match = patt.FindAllString(query, -1)
	if len(_match) == 0 {
		limit := 10
		if _, ok := params["data"].(map[string]interface{})["limit"]; ok {
			limit = int(params["data"].(map[string]interface{})["limit"].(float64))
		}
		offset := 0
		if _, ok := params["data"].(map[string]interface{})["offset"]; ok {
			offset = int(params["data"].(map[string]interface{})["offset"].(float64))
		}
		if limit != -1 {
			query = fmt.Sprintf(`%s LIMIT %d OFFSET %d`, query, limit, offset)
		}
	}
	results, cols, _, err := db.QueryMultiRowsWithCols(query, []interface{}{}...)
	//fmt.Println(query, (*results))
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	n_rows, _, err := db.QuerySingleRow(query_n_rows, []interface{}{}...)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"msg":     fmt.Sprintf("%s", err),
		}
	}
	total := 0
	total = int((*n_rows)["n_rows"].(int64))
	msg, _ := app.i18n.T("success", map[string]interface{}{})
	return map[string]interface{}{
		"success": true,
		"msg":     msg,
		"sql":     query,
		"data":    *results,
		"n_rows":  total,
		"cols":    cols,
	}
}
