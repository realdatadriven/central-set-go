package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"

	"github.com/pascaldekloe/jwt"
)

func (app *application) run_etlx_run_by_name(w http.ResponseWriter, r *http.Request) {
	params := map[string]any{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	fmt.Println(name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = map[string]any{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = map[string]any{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	params["user"] = *(contextGetAuthenticatedUser(r))
	var data map[string]any
	if !token["success"].(bool) {
		data = token
	} else {
		params["data"] = map[string]any{"name": name}
		data = app.etlxRunByName(params)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) dyn_api(w http.ResponseWriter, r *http.Request) {
	var params map[string]any
	ctrl := r.PathValue("ctrl")
	act := r.PathValue("act")
	err := request.DecodeJSON(w, r, &params)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = map[string]any{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = map[string]any{}
	}
	err = app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	//user := *(contextGetAuthenticatedUser(r))
	params["user"] = *(contextGetAuthenticatedUser(r))
	//fmt.Println(params["user"].(map[string]any)["username"].(string), "->", int(params["user"].(map[string]any)["user_id"].(float64)), "->", int(params["user"].(map[string]any)["role_id"].(float64)))
	var data map[string]any
	_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
	_log := map[string]any{
		"user_id": params["user"].(map[string]any)["user_id"],
		"action":  fmt.Sprintf("%s/%s", ctrl, act),
		"req_ip":  _ip,
		"res_at":  time.Now(),
	}
	switch ctrl {
	case "login":
		switch act {
		case "login":
			//app.login(w, r)
			data = app._login(params)
		case "chk_token":
			data = app.verifyToken(r)
		case "alter_pass":
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.alter_pass(params)
			}
		default:
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "admin":
		if act == "apps" {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.apps(params)
			}
		} else if act == "tables" {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.tables(params, []any{})
			}
		} else if act == "menu" {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.menu(params)
			}
		} else if app.contains([]any{"save_table_schema", "create_table_schema", "create_table", "add_table"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.save_table_schema(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "access":
		if app.contains([]any{"tables", "table_access", "permissions"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.table_access(params, []any{})
			}
		} else if app.contains([]any{"row_level_access", "row_level", "rla"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.row_level_access(params, []any{}, []any{})
			}
		} else if app.contains([]any{"row_level_tables", "rla_tables"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.row_level_tables(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "crud":
		if app.contains([]any{"read", "r", "R"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.read(params)
			}
		} else if app.contains([]any{"create", "c", "C", "update", "u", "U", "delete", "d", "D", "create_update"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.create_update(params)
			}
		} else if app.contains([]any{"query", "queries", "q", "Q"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.query(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "export":
		if app.contains([]any{"query", "q", "Query", "Q", "QUERY"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.export_query(params)
			}
		} else if app.contains([]any{"read", "r", "Read", "R", "READ"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.export_read(params)
			}
		} else if app.contains([]any{"dump_file_2_object", "file_2_object", "get_file_content", "file_contet", "file_data"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.dump_file_2_object(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "etl":
		if app.contains([]any{"extract", "Extract", "EXTRACT", "input", "Input", "e", "E", "i", "I"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.extract(params)
			}
		} else if app.contains([]any{"nrows", "n_rows", "rows", "NROWS", "N_ROWS", "ROWS"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.n_rows(params)
			}
		} else if app.contains([]any{"delete", "del", "d", "Delete", "Del", "D", "DELETE", "DEL"}, act) {
			if !token["success"].(bool) {
				err = response.JSON(w, http.StatusOK, token)
				if err != nil {
					app.serverError(w, r, err)
				}
				return
			}
			data = app.delete(params)
		} else if app.contains([]any{"output", "transform", "t", "Output", "Transform", "T", "OUTPUT", "TRANSFORM"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.transform(params)
			}
		} else if app.contains([]any{"export", "load", "E", "L", "e", "l", "Export", "Load", "EXPORT", "LOAD"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.export(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "etlx":
		if app.contains([]any{"config", "parse", "parse_config", "conf", "parse_conf", "parse_md", "get_config"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.etlxMdParse(params)
			}
		} else if app.contains([]any{"run", "exec", "execute", "start", "init"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.etlxRun(params)
			}
		} else if app.contains([]any{"parserun", "parse_run", "parse&run"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.etlxParseRun(params)
			}
		} else {
			data = map[string]any{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	default:
		data = map[string]any{
			"success": false,
			"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			"data":    params,
			"ctrl":    ctrl,
			"act":     act,
		}
	}
	// LOGS
	actions_not_to_log := app.sliceStrs2SliceInterfaces(strings.Split(app.config.actions_not_to_log, ","))
	if !app.contains(actions_not_to_log, act) {
		_log["res_type"] = "success"
		if _, ok := data["success"]; !ok {
			_log["res_type"] = "error"
		} else if _, ok := data["success"].(bool); !ok {
			_log["res_type"] = "error"
		} else if success, ok := data["success"].(bool); ok {
			if success {
				_log["res_type"] = "success"
			}
		}
		_log["res_msg"] = data["msg"]
		_log["row_id"] = data["inserted_primary_key"]
		_log["table"] = params["data"].(map[string]any)["table"]
		_log["db"] = ""
		if _, ok := params["data"].(map[string]any)["database"]; ok {
			_log["db"] = params["data"].(map[string]any)["database"]
		} else if _, ok := params["data"].(map[string]any)["db"]; ok {
			_log["db"] = params["data"].(map[string]any)["db"]
		} else if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(map[string]any)["db"]; ok {
			_log["db"] = params["app"].(map[string]any)["db"]
		}
		if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(map[string]any)["app_id"]; ok {
			_log["app_id"] = params["app"].(map[string]any)["app_id"]
		}
		_log["excluded"] = false
		//fmt.Println(_log)
		_log_params := map[string]any{
			"data": map[string]any{
				"data":  _log,
				"table": "user_log",
				"db":    app.config.db.dsn,
			},
			"app": map[string]any{
				"app_id": any(1.0),
				"db":     filepath.Base(app.config.db.dsn),
			},
			"user": map[string]any{
				"user_id": any(1.0),
				"role_id": any(1.0),
			},
		}
		res := app.create_update(_log_params)
		if _, ok := res["success"]; !ok {
			fmt.Println("Err processing logs:", res)
		} else if _, ok := res["success"].(bool); !ok {
			fmt.Println("Err processing logs:", res["msg"])
		} else if !res["success"].(bool) {
			fmt.Println("Err processing logs:", res["msg"])
		}
	}
	// BROADCAST CHAGE WS
	broadcast_changes := app.sliceStrs2SliceInterfaces(strings.Split(app.config.broadcast_changes, ","))
	if app.contains(broadcast_changes, act) {
		if _, ok := data["success"]; !ok {
		} else if _, ok := data["success"].(bool); !ok {
		} else if success, ok := data["success"].(bool); ok {
			fmt.Println("BROADCAST CHAGE WS:", act, broadcast_changes)
			if success {
				manager := app.NewConnectionManager()
				app.broadcastTableChange(manager, map[string]any{
					"type":     "data_change",
					"database": _log["db"],
					"table":    _log["table"],
				})
			}
		}
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) verifyToken(r *http.Request) map[string]any {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader != "" {
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			token := headerParts[1]
			claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
			if err != nil {
				return map[string]any{
					"success": false,
					"msg":     "Error validating token!",
				}
			}
			if !claims.Valid(time.Now()) {
				return map[string]any{
					"success": false,
					"msg":     "Token has expired!",
				}
			}
			if claims.Issuer != app.config.baseURL {
				return map[string]any{
					"success": false,
					"msg":     "Token is invalid",
				}
			}
			if !claims.AcceptAudience(app.config.baseURL) {
				return map[string]any{
					"success": false,
					"msg":     "Token is invalid!",
				}
			}
			var user map[string]any
			//print(1, " ", claims.Subject, "\n")
			err2 := json.Unmarshal([]byte(claims.Subject), &user)
			if err2 == nil {
				//print(2, " ", user["username"].(string), "\n")
				contextSetAuthenticatedUser(r, &user)
			}
			return map[string]any{
				"success": true,
				"msg":     "Token validated!",
			}
		} else {
			return map[string]any{
				"success": false,
				"msg":     "Token is invalid!",
			}
		}
	}
	return map[string]any{
		"success": false,
		"msg":     "No token received!",
	}
}
