package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"
	"github.com/realdatadriven/etlx"
)

// QuackAPIManager manages DuckDB instances running configured QuackAPI SQL.
type QuackAPIManager struct {
	pool    map[int]etlx.DBInterface
	configs map[int]Dict
	mux     sync.RWMutex
}

func NewQuackAPIManager() *QuackAPIManager {
	return &QuackAPIManager{
		pool:    make(map[int]etlx.DBInterface),
		configs: make(map[int]Dict),
	}
}

func RenderTextTemplate(tmplStr string, data map[string]any) (string, error) {
	funcMap := sprig.FuncMap()
	tmpl, err := texttemplate.New("tmpl").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	return buf.String(), nil
}

// Start runs startup_sql followed by main_sql on a new in-memory DuckDB.
func (qm *QuackAPIManager) Start(ctx context.Context, quackAPIID int, config Dict) error {
	if !env.GetBool("QUACK_API_ENABLED", false) {
		return fmt.Errorf("quack API is not enabled")
	}
	qm.mux.Lock()
	defer qm.mux.Unlock()
	if _, exists := qm.pool[quackAPIID]; exists {
		return fmt.Errorf("quack API %d is already running", quackAPIID)
	}
	conn, err := etlx.GetDB("duckdb:")
	if err != nil {
		return fmt.Errorf("failed to create in-memory duckdb connection: %w", err)
	}
	if startupSQL := toString(config["startup_sql"]); startupSQL != "" {
		_, err := conn.ExecuteQuery(startupSQL, Dict{})
		if err != nil {
			conn.Close()
			return fmt.Errorf("startup SQL failed: %w", err)
		}
	}
	if mainSQL := toString(config["main_sql"]); mainSQL != "" {
		_, err := conn.ExecuteQuery(mainSQL, Dict{})
		if err != nil {
			conn.Close()
			return fmt.Errorf("main SQL failed: %w", err)
		}
	}
	// get QUACK_API_RPROXY_PATH & QUACK_API_RPROXY_TMPL to put httpserver route for subdomain
	path := env.GetString("QUACK_API_RPROXY_PATH", "")
	tmpl := env.GetString("QUACK_API_RPROXY_TMPL", "")
	if path != "" && tmpl != "" {
		pth_cnf, err := RenderTextTemplate(path, config)
		if err != nil {
			fmt.Printf("failed to render proxy config template path: %w\n", err)
		}
		prx_cnf, err := RenderTextTemplate(tmpl, config)
		if err != nil {
			fmt.Printf("failed to render proxy config template: %w\n", err)
		} else if !fileExists(pth_cnf) {
			err := os.WriteFile(pth_cnf, []byte(prx_cnf), 0644)
			if err != nil {
				fmt.Printf("Failed to write file: %s", err)
			}
		}
	}
	qm.pool[quackAPIID] = conn
	qm.configs[quackAPIID] = config
	return nil
}

// Stop runs shutdown_sql only for an explicit stop, then closes its DuckDB.
func (qm *QuackAPIManager) Stop(ctx context.Context, quackAPIID int) error {
	if !env.GetBool("QUACK_API_ENABLED", false) {
		return fmt.Errorf("quack API is not enabled")
	}
	qm.mux.Lock()
	defer qm.mux.Unlock()
	conn, exists := qm.pool[quackAPIID]
	if !exists {
		return fmt.Errorf("quack API %d is not running", quackAPIID)
	}
	var shutdownErr error
	if shutdownSQL := toString(qm.configs[quackAPIID]["shutdown_sql"]); shutdownSQL != "" {
		_, err := conn.ExecuteQuery(shutdownSQL, Dict{})
		if err != nil {
			shutdownErr = fmt.Errorf("shutdown SQL failed: %w", err)
		}
	}
	/*path := env.GetString("QUACK_API_RPROXY_PATH", "")
	tmpl := env.GetString("QUACK_API_RPROXY_TMPL", "")
	if path != "" && tmpl != "" {
		pth_cnf, err := RenderTextTemplate(path, qm.configs[quackAPIID])
		if err != nil {
			fmt.Printf("failed to render proxy path template: %w\n", err)
		} else if !fileExists(pth_cnf) {
			err := os.Remove(pth_cnf)
			if err != nil {
				fmt.Printf("Failed to remove file: %s", err)
			}
		}
	}*/
	closeErr := conn.Close()
	delete(qm.pool, quackAPIID)
	delete(qm.configs, quackAPIID)
	if shutdownErr != nil {
		return shutdownErr
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close quack API %d connection: %w", quackAPIID, closeErr)
	}
	return nil
}

func (qm *QuackAPIManager) Restart(ctx context.Context, quackAPIID int, config Dict) error {
	if !env.GetBool("QUACK_API_ENABLED", false) {
		return fmt.Errorf("quack API is not enabled")
	}
	if err := qm.Stop(ctx, quackAPIID); err != nil {
		return fmt.Errorf("failed to stop quack API during restart: %w", err)
	}
	if err := qm.Start(ctx, quackAPIID, config); err != nil {
		return fmt.Errorf("failed to start quack API during restart: %w", err)
	}
	return nil
}

func (qm *QuackAPIManager) GetConnector(quackAPIID int) (etlx.DBInterface, error) {
	qm.mux.RLock()
	defer qm.mux.RUnlock()
	conn, exists := qm.pool[quackAPIID]
	if !exists {
		return nil, fmt.Errorf("quack API %d is not running", quackAPIID)
	}
	return conn, nil
}

func (app *application) quackAPIManagerInstance() *QuackAPIManager {
	app.quackAPIOnce.Do(func() {
		app.quackAPIManager = NewQuackAPIManager()
	})
	return app.quackAPIManager
}

func (app *application) startQuackAPI(w http.ResponseWriter, r *http.Request) {
	app.quackAPIAction(w, r, "start")
}

func (app *application) stopQuackAPI(w http.ResponseWriter, r *http.Request) {
	app.quackAPIAction(w, r, "stop")
}

func (app *application) restartQuackAPI(w http.ResponseWriter, r *http.Request) {
	app.quackAPIAction(w, r, "restart")
}

func (app *application) quackAPIAction(w http.ResponseWriter, r *http.Request, action string) {
	params := Dict{}
	if err := request.DecodeJSON(w, r, &params); err != nil &&
		(err.Error() != "body must not be empty" || r.PathValue("name") == "") {
		app.badRequest(w, r, err)
		return
	}
	if _, ok := params["data"].(Dict); !ok {
		params["data"] = Dict{}
	}
	params["quack_api_path_name"] = r.PathValue("name")
	lang := "en"
	if value, ok := params["lang"].(string); ok {
		lang = value
	}
	if err := app.i18n.ChangeLanguage(lang); err != nil {
		fmt.Println(err)
	}
	token := app.verifyToken(r)
	user := *(contextGetAuthenticatedUser(r))
	params["user"] = user
	var result Dict
	switch {
	case !token["success"].(bool):
		result = token
	case !app.quackEnabled:
		result = Dict{"success": false, "msg": "Quack is not enabled"}
	default:
		if userParams, ok := user["params"]; ok {
			params["params"] = userParams
		} else if userData, ok := user["data"]; ok {
			params["params"] = userData
		}
		result = app.runQuackAPIAction(params, action)
	}
	if err := response.JSON(w, http.StatusOK, result); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) runQuackAPIAction(params Dict, action string) Dict {
	config, err := app.getQuackAPIConfig(params)
	if err != nil {
		return Dict{"success": false, "msg": err.Error()}
	}
	apiID := app.toInt(config["quack_api_id"])
	manager := app.quackAPIManagerInstance()
	switch action {
	case "start":
		err = manager.Start(context.Background(), apiID, config)
	case "stop":
		err = manager.Stop(context.Background(), apiID)
	case "restart":
		err = manager.Restart(context.Background(), apiID, config)
	default:
		return Dict{"success": false, "msg": "Unknown QuackAPI action"}
	}
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("Failed to %s QuackAPI: %v", action, err)}
	}
	return Dict{"success": true, "msg": fmt.Sprintf("QuackAPI %s successful", action)}
}

func (app *application) getQuackAPIConfig(params Dict) (Dict, error) {
	data, _ := params["data"].(Dict)
	if nested, ok := data["data"].(Dict); ok {
		data = nested
	}
	if nested, ok := data["quack_api"].(Dict); ok {
		data = nested
	}
	apiName, _ := data["quack_api"].(string)
	var field string
	var value any
	switch {
	case data["quack_api_id"] != nil:
		field, value = "quack_api_id", data["quack_api_id"]
	case apiName != "":
		field, value = "quack_api", apiName
	case data["name"] != nil:
		field, value = "quack_api", data["name"]
	case params["quack_api_id"] != nil:
		field, value = "quack_api_id", params["quack_api_id"]
	case params["quack_api"] != nil:
		field, value = "quack_api", params["quack_api"]
	case params["name"] != nil:
		field, value = "quack_api", params["name"]
	case params["quack_api_path_name"] != nil && params["quack_api_path_name"] != "":
		field, value = "quack_api", params["quack_api_path_name"]
	default:
		return nil, fmt.Errorf("no QuackAPI service was specified")
	}
	readParams := make(Dict, len(params))
	for key, item := range params {
		readParams[key] = item
	}
	readParams["db"] = "ADMIN"
	readParams["data"] = Dict{
		"table":  "quack_api",
		"db":     "ADMIN",
		"limit":  1.0,
		"offset": 0.0,
		"filters": []any{Dict{
			"field": field,
			"cond":  "=",
			"value": value,
		}},
	}
	result := app.read(readParams)
	if result["success"] != true {
		return nil, fmt.Errorf("QuackAPI service not found or inactive")
	}
	apis, ok := result["data"].([]Dict)
	if !ok || len(apis) == 0 {
		return nil, fmt.Errorf("QuackAPI service not found or inactive")
	}
	return apis[0], nil
}
