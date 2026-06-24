package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"

	"github.com/joho/godotenv"
	"github.com/pascaldekloe/jwt"
)

func (app *application) run_backup(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	fmt.Println("run_backup:", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	params["user"] = *(contextGetAuthenticatedUser(r))
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else {
		params["data"] = Dict{"name": name}
		data = app.Buckup(params)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) refreshEnv(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	params["user"] = *(contextGetAuthenticatedUser(r))
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else {
		data = app.syncEnv()
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) run_notebook(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	fmt.Println("run_notebook:", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	params["user"] = *(contextGetAuthenticatedUser(r))
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else {
		params["data"] = Dict{"name": name}
		data = app.nbRun(params)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) run_etlx_run_by_name(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	// fmt.Println("run_etlx_run_by_name: ", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	user := *(contextGetAuthenticatedUser(r))
	params["user"] = user
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else {
		params["data"] = Dict{"name": name}
		if _, ok := user["params"]; ok {
			params["params"] = user["params"]
		} else if _, ok := user["data"]; ok {
			params["params"] = user["data"]
		}
		data = app.etlxRunByName(params)
		//fmt.Println(data)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) startQuack(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	// fmt.Println("run_etlx_run_by_name: ", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	user := *(contextGetAuthenticatedUser(r))
	params["user"] = user
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else if !app.quackEnabled {
		data = Dict{
			"success": false,
			"msg":     "Quack is not enabled",
		}
	} else {
		params["data"] = Dict{"name": name}
		if _, ok := user["params"]; ok {
			params["params"] = user["params"]
		} else if _, ok := user["data"]; ok {
			params["params"] = user["data"]
		}
		data = app.startQuackServer(params)
		//fmt.Println(data)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) stopQuack(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	// fmt.Println("run_etlx_run_by_name: ", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	user := *(contextGetAuthenticatedUser(r))
	params["user"] = user
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else if !app.quackEnabled {
		data = Dict{
			"success": false,
			"msg":     "Quack is not enabled",
		}
	} else {
		params["data"] = Dict{"name": name}
		if _, ok := user["params"]; ok {
			params["params"] = user["params"]
		} else if _, ok := user["data"]; ok {
			params["params"] = user["data"]
		}
		data = app.stopQuackServer(params)
		//fmt.Println(data)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) restartQuack(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := r.PathValue("name")
	// fmt.Println("run_etlx_run_by_name: ", name)
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err := app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	user := *(contextGetAuthenticatedUser(r))
	params["user"] = user
	var data Dict
	if !token["success"].(bool) {
		data = token
	} else if !app.quackEnabled {
		data = Dict{
			"success": false,
			"msg":     "Quack is not enabled",
		}
	} else {
		params["data"] = Dict{"name": name}
		if _, ok := user["params"]; ok {
			params["params"] = user["params"]
		} else if _, ok := user["data"]; ok {
			params["params"] = user["data"]
		}
		data = app.restartQuackServer(params)
		//fmt.Println(data)
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) getLocationFromRequest(r *http.Request, params Dict) *time.Location {
	// timezone
	var loc *time.Location
	var err error
	if tz, ok := params["timezone"].(string); ok {
		// Do something with the timezone parameter
		loc, err = time.LoadLocation(tz)
		if err != nil {
			tz := r.Header.Get("X-Timezone")
			if tz != "" {
				loc, err := time.LoadLocation(tz)
				if err != nil {
					fmt.Println("Error loading timezone from header:", err)
					loc = time.Local
				} else {
					loc = loc
				}
			} else {
				loc = time.Local
			}
		} else {
			loc = time.Local
		}
	} else {
		// const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone in javascript to get the timezone of the user and send it in the header X-Timezone
		tz := r.Header.Get("X-Timezone")
		if tz != "" {
			loc, err = time.LoadLocation(tz)
			if err != nil {
				fmt.Println("Error loading timezone from header:", err)
				loc = time.Local
			} else {
				loc = time.Local
			}
		} else {
			loc = time.Local
		}
	}
	return loc
}

func (app *application) dyn_api(w http.ResponseWriter, r *http.Request) {
	var params Dict
	ctrl := r.PathValue("ctrl")
	act := r.PathValue("act")
	fmt.Println(ctrl, act)
	err := request.DecodeJSON(w, r, &params)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	// timezone
	loc := app.getLocationFromRequest(r, params)
	params["location"] = loc
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	if _, ok := params["data"]; !ok {
		params["data"] = Dict{}
	}
	if _, ok := params["app"]; !ok {
		params["app"] = Dict{}
	}
	err = app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	params["token"] = r.Header.Get("Authorization")
	token := app.verifyToken(r)
	//fmt.Println(params["user"].(Dict)["username"].(string), "->", app.toInt(params["user"].(Dict)["user_id"].(float64)), "->", app.toInt(params["user"].(Dict)["role_id"].(float64)))
	var data Dict
	_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
	_log := Dict{
		"action": fmt.Sprintf("%s/%s", ctrl, act),
		"req_ip": _ip,
		"req_at": time.Now().In(loc),
	}
	fmt.Println(token, params)
	if token["success"].(bool) {
		//user := *(contextGetAuthenticatedUser(r))
		params["user"] = *(contextGetAuthenticatedUser(r))
		_log["user_id"] = params["user"].(Dict)["user_id"]
	}
	//check if app.appType is community, licensor, or licensee, if it is licensee check the enviromental varibales CS_LICENCOR_TOKEN and CS_LICENCOR_URL
	//use the those to make a post request to the CS_LICENCOR_URL/dyn_api/license/verify_license endpoint with the token in the header Authorization
	if app.appType == "licensee" {
		// also i want the verifcation to be done only once per app.licenceVerificationPeriodicity in the app.lastLicenseValidation timestamp
		if time.Since(app.lastLicenseValidation) >= app.licenceVerificationPeriodicity {
			token = app.validateLicense()
			// update the app.lastLicenseValidation timestamp if success
			if token["success"].(bool) {
				app.lastLicenseValidation = time.Now().In(loc) // when this became older the app.licenceVerificationPeriodicity it will check again
			}
		} // else it will be verifyed untill the licence gets validated again
	}
	// ROUTES
	switch ctrl {
	// handle ctrl = license and act = verify_license that baically just checks the token sent in the header Authorization
	case "license", "lic":
		switch act {
		case "verify_license", "verify_lic", "verify": // this route just checks the license token validity and returns the token data in a licensee licensor setup
			if !token["success"].(bool) {
				data = Dict{
					"success": false,
					"msg":     "Lience token validation faild please contact the admin!",
				}
			} else {
				data = token
			}
		default:
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "env", "environment":
		switch act {
		case "refresh", "update", "sync":
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.syncEnv()
			}
		default:
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "login":
		switch act {
		case "login", "sign_in", "signin", "auth", "authenticate", "log_in", "logon", "log_on", "index":
			// fmt.Println("LOGIN:")
			data = app._login(params)
			// fmt.Println("LOGIN DATA:", data)
		case "validate_code", "validate_2f_code", "valid_code", "valid_2f_code", "2f_code", "two_factor_code":
			data = app.two_factor_code_valid(params)
		// dynamic_login
		case "dynamic_login", "dynamic_auth", "dynamic_authenticate":
			if os.Getenv("DYN_LOGIN_TABLE_MAP_TO_USERS") == "true" {
				data = Dict{
					"success": false,
					"msg":     "Dynamic login is not allowed, contatt the admin to enable it if you think it is needed for your use case!",
				}
			} else {
				params["login_table"] = os.Getenv("DYN_LOGIN_TABLE")
				params["user_id_field"] = os.Getenv("DYN_LOGIN_USER_ID_FIELD")
				params["dyn_login_role_id"] = os.Getenv("DYN_LOGIN_ROLE_ID")
				params["username_field"] = os.Getenv("DYN_LOGIN_USERNAME_FIELD")
				params["email_field"] = os.Getenv("DYN_LOGIN_EMAIL_FIELD")
				params["password_field"] = os.Getenv("DYN_LOGIN_PASSWORD_FIELD")
				params["active_field"] = os.Getenv("DYN_LOGIN_ACTIVE_FIELD")
				data = app.dynamic_login(params)
			}
		case "signup", "sign_up", "dynamic_signup", "dyn_signup", "dynamic_sign_up", "dynamic_register", "dyn_register", "dynamic_registration":
			params["login_table"] = os.Getenv("DYN_LOGIN_TABLE")
			params["user_id_field"] = os.Getenv("DYN_LOGIN_USER_ID_FIELD")
			params["dyn_login_role_id"] = os.Getenv("DYN_LOGIN_ROLE_ID")
			params["username_field"] = os.Getenv("DYN_LOGIN_USERNAME_FIELD")
			params["email_field"] = os.Getenv("DYN_LOGIN_EMAIL_FIELD")
			params["password_field"] = os.Getenv("DYN_LOGIN_PASSWORD_FIELD")
			params["active_field"] = os.Getenv("DYN_LOGIN_ACTIVE_FIELD")
			data = app.dynamic_signup(params)
		// social_login
		//case "social_login", "social_auth", "social_authenticate":
		//	data = app.social_login(params)
		// recover_pass
		case "recover_pass", "recover_password":
			//fmt.Println(ctrl, act, params)
			data = app.recover_pass(params)
			// fmt.Println("recover_pass data:", data)
		// reset_pass
		case "reset_pass", "reset_password":
			data = app.reset_pass(params)
		// verify_token
		case "chk_token", "verify_token":
			data = app.verifyToken(r)
		case "alter_pass":
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.alter_pass(params)
			}
		case "user", "user-from-token", "token-user", "usr":
			if !token["success"].(bool) {
				data = token
			} else {
				msg, _ := app.i18n.T("success", Dict{})
				data = Dict{
					"success": true,
					"msg":     msg,
					"data":    params["user"],
				}
			}
		case "access_key", "access_token", "credentials":
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.access_key(params)
			}
		default:
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "admin", "adm":
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
			data = Dict{
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
			data = Dict{
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
			data = Dict{
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
			data = Dict{
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
			data = Dict{
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
				if !env.GetBool("ETLX_ALLOW_CLI_CONFIG", false) {
					data = Dict{
						"success": false,
						"msg":     "Executing ETLX with cliente config is not allowed, ETLX_ALLOW_CLI_CONFIG must be true, please contact the Admin!",
					}
				} else {
					data = app.etlxRun(params, false)
				}
			}
		} else if app.contains([]any{"parserun", "parse_run", "parse&run"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				if !env.GetBool("ETLX_ALLOW_CLI_CONFIG", false) {
					data = Dict{
						"success": false,
						"msg":     "Executing ETLX with cliente config is not allowed, ETLX_ALLOW_CLI_CONFIG must be true, please contact the Admin!",
					}
				} else {
					data = app.etlxParseRun(params)
				}
			}
		} else if app.contains([]any{"run_by_name", "run_name", "name", "by_name", "byName", "byname"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.etlxRunByName(params)
			}
		} else if app.contains([]any{"query_md", "query-md", "query-markdown", "query_markdown"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.queryETLXMD(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "nb", "notebook", "noteb", "nbook":
		if app.contains([]any{"run_cell", "cell", "c", "run_cells", "cells", "cs"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.nbRunCells(params)
			}
		} else if app.contains([]any{"run_by_name", "run_name", "name", "by_name", "byName", "byname"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.nbRunByName(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "api", "API", "external_api", "external", "ext_api", "extapi":
		if app.contains([]any{"run", "execute", "exec", "exe", "x"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.runAPI(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "crud_actions", "run_actions", "user_triggered_actions", "actions":
		if app.contains([]any{"run", "execute", "exec", "exe", "x"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.UserTriggeredCrudAction(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "quack", "quack-protocol", "quack-serve", "quack-server", "quackserve", "quackserver":
		if app.contains([]any{"start", "startup", "run", "open", "index"}, act) {
			if !token["success"].(bool) {
				data = token
			} else if !app.quackEnabled {
				data = Dict{
					"success": false,
					"msg":     "Quack is not enabled",
				}
			} else {
				data = app.startQuackServer(params)
			}
		} else if app.contains([]any{"restart", "reboot", "rerun", "reopen"}, act) {
			if !token["success"].(bool) {
				data = token
			} else if !app.quackEnabled {
				data = Dict{
					"success": false,
					"msg":     "Quack is not enabled",
				}
			} else {
				data = app.restartQuackServer(params)
			}
		} else if app.contains([]any{"stop", "shutdown", "break"}, act) {
			if !token["success"].(bool) {
				data = token
			} else if !app.quackEnabled {
				data = Dict{
					"success": false,
					"msg":     "Quack is not enabled",
				}
			} else {
				data = app.stopQuackServer(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "cron":
		if app.contains([]any{"run", "r", "execute", "exec", "e"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				_aux, ok := params["data"].(Dict)
				_data := Dict{}
				if ok {
					_data = _aux["data"].(Dict)
				}
				_jwt, _ := app.getToken(r)
				_data["token"] = _jwt
				res, err := app.CronRunEndPoint(_data)
				if err != nil {
					data = Dict{"success": false, "msg": fmt.Sprintf("Error %s", err)}
				} else {
					_, ok1 := res["success"].(bool)
					_, ok2 := res["msg"].(string)
					if ok1 && ok2 {
						data = res
					} else {
						msg, _ := app.i18n.T("success", Dict{})
						data = Dict{"success": true, "msg": msg}
					}
				}
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "saas", "software_as_a_service", "software_service":
		if app.contains([]any{"deploy", "dp"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.RunDeploy(params)
			}
		} else if app.contains([]any{"cancel", "destroy", "quit", "drop"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				data = app.RunDeploy(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	case "pay", "payment", "stripe":
		if app.contains([]any{"syncprod", "syncproduct", "sync_prod", "sync_product"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				app.PaymentInit()
				// fmt.Println("SyncOrCreateProduct:", params["data"])
				data = app.SyncOrCreateProduct(params)
			}
		} else if app.contains([]any{"synccust", "synccustomer", "sync_cust", "sync_customer"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				app.PaymentInit()
				data = app.CreateOrSyncCustomer(params)
			}
		} else if app.contains([]any{"syncsub", "syncsubscription", "sync_subs", "sync_subscription"}, act) {
			if !token["success"].(bool) {
				data = token
			} else {
				app.PaymentInit()
				data = app.CreateOrUpdateSubscription(params)
			}
		} else {
			data = Dict{
				"success": false,
				"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			}
		}
	default:
		data = Dict{
			"success": false,
			"msg":     fmt.Sprintf("No route %s/%s exists yet!", ctrl, act),
			"data":    params,
			"ctrl":    ctrl,
			"act":     act,
		}
	}
	// LOGS
	actions_not_to_log := app.sliceStrs2SliceInterfaces(strings.Split(app.config.actions_not_to_log, ","))
	// fmt.Println(actions_not_to_log)
	if !app.contains(actions_not_to_log, act) {
		_log["res_at"] = time.Now().In(loc)
		_log["res_type"] = "success"
		if _, ok := data["success"].(bool); !ok {
			_log["res_type"] = "error"
		} else if success, ok := data["success"].(bool); ok {
			if success {
				_log["res_type"] = "success"
			} else {
				_log["res_type"] = "error"
			}
		}
		_log["res_msg"] = data["msg"]
		_log["row_id"] = data["inserted_primary_key"]
		_log["table"] = params["data"].(Dict)["table"]
		_log["db"] = ""
		if _, ok := params["data"].(Dict)["database"]; ok {
			_log["db"] = params["data"].(Dict)["database"]
		} else if _, ok := params["data"].(Dict)["db"]; ok {
			_log["db"] = params["data"].(Dict)["db"]
		} else if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(Dict)["db"]; ok {
			_log["db"] = params["app"].(Dict)["db"]
		}
		if _, ok := params["app"]; !ok {
		} else if _, ok := params["app"].(Dict)["app_id"]; ok {
			_log["app_id"] = params["app"].(Dict)["app_id"]
		}
		_log["excluded"] = false
		// fmt.Println("LOGS:", _log)
		_log_params := Dict{
			"data": Dict{
				"data":  _log,
				"table": "user_log",
				"db":    app.config.db.dsn,
			},
			"app": Dict{
				"app_id": any(1.0),
				"db":     filepath.Base(app.config.db.dsn),
			},
			"user": Dict{
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
			// fmt.Println("BROADCAST CHAGE WS:", success, act, broadcast_changes, _log)
			if success {
				_data := Dict{
					"type":     "data_change",
					"database": _log["db"],
					"table":    _log["table"],
				}
				// WS
				if app.WS_ConnectionManager != nil {
					//manager := app.NewConnectionManager()
					app.broadcastTableChange(app.WS_ConnectionManager, _data)
				}
				// SSE
				if env.GetBool("SSE_ENABLE", false) {
					if app.SSE_Broker != nil {
						app.SSE_Broker.NotifyAll(_data)
					}
				}
			}
		}
	}
	// fmt.Println(data)
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// syncEnv loads environment variables from the database and sets them in the OS environment
func (app *application) syncEnv() Dict {
	// load the .env file
	_err := godotenv.Load()
	if _err != nil {
		fmt.Println("Error loading .env file")
	}
	sql := `select * from "env" where "active" = ? /*and "on_srv_start" = ?*/ and "excluded" = ?`
	tenantEnv, err := app.AdminGetRowsByFilter(sql, []any{true, false})
	if err != nil {
		fmt.Printf("Error fetching tenant env vars: %v\n", err)
	} else {
		for _, v := range tenantEnv {
			os.Setenv(v["env_name"].(string), v["env_value"].(string))
			fmt.Printf("Setting env var for admin %s=%s\n", v["env_name"], v["env_value"])
		}
	}
	return Dict{
		"success": true,
		"msg":     "Environment variables synchronized successfully",
	}
}
func (app *application) validateLicense() Dict {
	token := Dict{}
	licensorToken := os.Getenv("CS_LICENSOR_TOKEN")
	licensorURL := os.Getenv("CS_LICENSOR_URL")
	if licensorToken != "" && licensorURL != "" {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/dyn_api/license/verify_license", licensorURL), nil)
		if err != nil {
			token = Dict{
				"success": false,
				"msg":     fmt.Sprintf("Error creating request to licensor, please contact admin: %s", err),
			}
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", licensorToken))
			// is a aplication/json request
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				token = Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error making request to licensor, please contact admin: %s", err),
				}
			} else {
				defer resp.Body.Close()
				var licensorResp Dict
				err = json.NewDecoder(resp.Body).Decode(&licensorResp)
				if err != nil {
					token = Dict{
						"success": false,
						"msg":     fmt.Sprintf("Error decoding licensor response, please contact admin: %s", err),
					}
				} else {
					if success, ok := licensorResp["success"].(bool); ok {
						if success {
							token = licensorResp
						} else {
							token = Dict{
								"success": false,
								"msg":     "License validation failed from licensor, please contact admin!",
							}
						}
					} else {
						token = Dict{
							"success": false,
							"msg":     "Invalid response from licensor, please contact admin!",
						}
					}
				}
			}
		}
	} else {
		token = Dict{
			"success": false,
			"msg":     "Licensor token or URL not set, please contact admin!",
		}
	}
	return token
}
func (app *application) getToken(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	//fmt.Println("getToken:", authorizationHeader)
	if authorizationHeader != "" {
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			return headerParts[1], nil

		} else {
			return "", fmt.Errorf("token is invalid")
		}
	}
	return "", fmt.Errorf("No token received")
}
func (app *application) verifyToken(r *http.Request) Dict {
	authorizationHeader := r.Header.Get("Authorization")
	/*/ USE COOKIE PROVIDED IN THE OAUTH AS IF Authorization HEADER
	cookie, err := r.Cookie("session")
	if err == nil && authorizationHeader == "" {
		authorizationHeader = "Bearer " + cookie.Value
	}*/
	loc := app.getLocationFromRequest(r, Dict{})
	//fmt.Println("verifyToken:", authorizationHeader)
	if authorizationHeader != "" {
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			token := headerParts[1]
			claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
			if err != nil {
				return Dict{
					"success": false,
					"msg":     "Error validating token!",
				}
			}
			if !claims.Valid(time.Now().In(loc)) {
				return Dict{
					"success": false,
					"msg":     "Token has expired!",
				}
			}
			if claims.Issuer != app.config.baseURL {
				return Dict{
					"success": false,
					"msg":     "Token is invalid",
				}
			}
			if !claims.AcceptAudience(app.config.baseURL) {
				return Dict{
					"success": false,
					"msg":     "Token is invalid!",
				}
			}
			var user Dict
			//print(1, " ", claims.Subject, "\n")
			err2 := json.Unmarshal([]byte(claims.Subject), &user)
			if err2 == nil {
				//print(2, " ", user["username"].(string), "\n")
				contextSetAuthenticatedUser(r, &user)
			}
			return Dict{
				"success": true,
				"msg":     "Token validated!",
			}
		} else {
			return Dict{
				"success": false,
				"msg":     "Token is invalid!",
			}
		}
	}
	return Dict{
		"success": false,
		"msg":     "No token received!",
	}
}
func (app *application) verifyTokenString(authorizationHeader string) (Dict, error) {
	//fmt.Println("verifyTokenString:", authorizationHeader)
	if authorizationHeader != "" {
		token := authorizationHeader
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
		if err != nil {
			return nil, fmt.Errorf("Error validating token: %w", err)
		}
		if !claims.Valid(time.Now()) {
			return nil, fmt.Errorf("Token has expired: %w", err)
		}
		if claims.Issuer != app.config.baseURL {
			return nil, fmt.Errorf("Token Issuer is invalid: %w", app.config.baseURL)
		}
		if !claims.AcceptAudience(app.config.baseURL) {
			return nil, fmt.Errorf("Token AcceptAudience is invalid: %w", app.config.baseURL)
		}
		var user Dict
		// fmt.Println(claims.Subject)
		err2 := json.Unmarshal([]byte(claims.Subject), &user)
		if err2 != nil {
			return nil, fmt.Errorf("Error unmarshaling token subject: %w", err)
		}
		return user, nil
	}
	return nil, fmt.Errorf("No token received: %w", "")
}

// stripe payment webhook handler
/*func (app *application) stripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading request body: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	endpointSecret := os.Getenv("STRIPE_ENDPOINT_SECRET")
	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := stripe.WebhookConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(w, "Error verifying webhook signature: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		paymentIntent := event.Data.Object.(*stripe.PaymentIntent)
		fmt.Printf("PaymentIntent was successful! %s\n", paymentIntent.ID)
		// Then define and call a method to handle the successful payment intent.
		// handlePaymentIntentSucceeded(paymentIntent)
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}
	w.WriteHeader(http.StatusOK)
}*/

// stripe create a product and price endpoint
/*func (app *application) createStripeProductAndPrice(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	request.DecodeJSON(w, r, &params)
	name := params["name"].(string)
	description := params["description"].(string)
	unitAmount := int64(params["unit_amount"].(float64)) // in cents
	currency := params["currency"].(string)
	// Create a new product
	productParams := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}
	product, err := product.New(productParams)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Create a new price for the product
	priceParams := &stripe.PriceParams{
		Product:    stripe.String(product.ID),
		UnitAmount: stripe.Int64(unitAmount),
		Currency:   stripe.String(currency),
	}
	price, err := price.New(priceParams)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := Dict{
		"success": true,
		"msg":     "Product and price created successfully",
		"data": Dict{
			"product": product,
			"price":   price,
		},
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}*/
