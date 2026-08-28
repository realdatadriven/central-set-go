package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sync"
	"time"

	"github.com/duckdb/duckdb-go/v2"
	"github.com/realdatadriven/etlx"
)

// QuackManager manages the lifecycle of in-memory DuckDB instances keyed by server ID
type QuackManager struct {
	pool          map[int]etlx.DBInterface // Keyed by quack_server_id
	mux           sync.RWMutex
	adminDB       etlx.DBInterface // Reference to admin DB for config fetch
	validateToken func(string) (Dict, error)
	quackConfigs  map[int]Dict // Cache of server configs
	configMu      sync.RWMutex
}

// NewQuackManager creates a new Quack manager instance
func NewQuackManager(adminDB etlx.DBInterface, validateToken func(string) (Dict, error)) *QuackManager {
	return &QuackManager{
		pool:          make(map[int]etlx.DBInterface),
		adminDB:       adminDB,
		validateToken: validateToken,
		quackConfigs:  make(map[int]Dict),
	}
}

// StartQuackServer initializes an in-memory DuckDB instance and runs startup/attach SQL
func (qm *QuackManager) StartQuackServer(ctx context.Context, quackServerID int, config Dict) error {
	qm.mux.Lock()
	defer qm.mux.Unlock()

	// Check if already running
	if _, exists := qm.pool[quackServerID]; exists {
		return fmt.Errorf("quack server %d is already running", quackServerID)
	}

	/*/ Fetch config from admin DB
	config, err := qm.fetchQuackServerConfig(quackServerID)
	if err != nil {
		return fmt.Errorf("failed to fetch quack server config: %w", err)
	}*/

	// Create in-memory DuckDB connection (DSN: "duckdb:")
	conn, err := etlx.GetDB("duckdb:")
	if err != nil {
		return fmt.Errorf("failed to create in-memory duckdb connection: %w", err)
	}
	_, err = conn.ExecuteQuery("INSTALL quack", []any{}...)
	if err != nil {
		return fmt.Errorf("INSTALL quack: %s", err)
	}
	host, ok := config["host"]
	if !ok {
		host = "localhost"
	}
	port, ok := config["port"]
	if !ok {
		port = "9494"
	}
	token, ok := config["token"]
	allowOtherHostnames, ok := config["allow_other_hostnames"]
	if !ok {
		allowOtherHostnames = "false"
	}
	disableSSL, ok := config["disable_ssl"]
	if !ok {
		disableSSL, ok = config["DISABLE_SSL"]
		if !ok {
			disableSSL = "false"
		}
	}
	// port , token
	sql := fmt.Sprintf("CALL quack_serve('quack:%s:%s', token => '%s', allow_other_hostname => %s, disable_ssl => %s);", host, port, token, allowOtherHostnames, disableSSL)
	if _, err := conn.ExecuteQuery(sql); err != nil {
		conn.Close()
		return fmt.Errorf("Quack start failed: %w", err)
	}
	// Execute startup SQL (e.g., "INSTALL SQLITE; LOAD SQLITE;")
	if config["startup_sql"] != "" {
		if err := qm.executeSQL(conn, config["startup_sql"].(string)); err != nil {
			conn.Close()
			return fmt.Errorf("startup sql failed: %w", err)
		}
	}
	// Execute main SQL (e.g., "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE); USE adm;")
	if config["main_sql"] != "" {
		if err := qm.executeSQL(conn, config["main_sql"].(string)); err != nil {
			conn.Close()
			return fmt.Errorf("attach sql failed: %w", err)
		}
	}
	/*/ Register the token validation UDF for this server
	if err := qm.registerCheckTokenUDF(conn, config); err != nil {
		conn.Close()
		return fmt.Errorf("failed to register check token udf: %w", err)
	}*/
	// Store in pool
	qm.pool[quackServerID] = conn
	// Cache config
	qm.configMu.Lock()
	qm.quackConfigs[quackServerID] = config
	qm.configMu.Unlock()
	// Log startup in quack_logs table
	// qm.logQuackEvent(quackServerID, "startup", "online", int64(config["port"].(float64))), "Server started successfully", true)
	return nil
}

// StopQuackServer gracefully shuts down a Quack server instance
func (qm *QuackManager) StopQuackServer(ctx context.Context, quackServerID int) error {
	qm.mux.Lock()
	defer qm.mux.Unlock()
	conn, exists := qm.pool[quackServerID]
	if !exists {
		return fmt.Errorf("quack server %d is not running", quackServerID)
	}
	// Fetch config
	qm.configMu.RLock()
	config, ok := qm.quackConfigs[quackServerID]
	qm.configMu.RUnlock()
	if ok && config["shutdown_sql"] != "" {
		// Execute shutdown SQL (e.g., "USE memory; DETACH adm;")
		if err := qm.executeSQL(conn, config["shutdown_sql"].(string)); err != nil {
			// Log but don't fail shutdown
			// // qm.logQuackEvent(quackServerID, "shutdown", "error", int64(config["port"].(float64))), fmt.Sprintf("Shutdown error: %v", err), false)
		}
	}
	config = qm.quackConfigs[quackServerID]
	host, ok := config["host"]
	if !ok {
		host = "localhost"
	}
	port, ok := config["port"]
	if !ok {
		port = "9494"
	}
	sql := fmt.Sprintf("CALL quack_stop('quack:%s:%s');", host, port)
	if _, err := conn.ExecuteQuery(sql); err != nil {
		conn.Close()
		return fmt.Errorf("Stop quack failed: %w", err)
	}
	// Close connection
	conn.Close()
	// Remove from pool
	delete(qm.pool, quackServerID)
	// Remove config
	qm.configMu.Lock()
	delete(qm.quackConfigs, quackServerID)
	qm.configMu.Unlock()
	// Log shutdown
	if ok {
		// qm.logQuackEvent(quackServerID, "shutdown", "offline", int64(config.Port), "Server stopped", true)
	}
	return nil
}

// RestartQuackServer restarts a Quack server (stop then start)
func (qm *QuackManager) RestartQuackServer(ctx context.Context, quackServerID int, config Dict) error {
	// Stop
	if err := qm.StopQuackServer(ctx, quackServerID); err != nil {
		return fmt.Errorf("failed to stop server during restart: %w", err)
	}
	// Start
	if err := qm.StartQuackServer(ctx, quackServerID, config); err != nil {
		return fmt.Errorf("failed to start server during restart: %w", err)
	}
	return nil
}

// GetQuackConnector retrieves the in-memory connector for a running server
func (qm *QuackManager) GetQuackConnector(quackServerID int) (etlx.DBInterface, error) {
	qm.mux.RLock()
	defer qm.mux.RUnlock()
	conn, exists := qm.pool[quackServerID]
	if !exists {
		return nil, fmt.Errorf("quack server %d is not running", quackServerID)
	}
	return conn, nil
}

// executeSQL runs a potentially multi-statement SQL script
func (qm *QuackManager) executeSQL(conn etlx.DBInterface, sql string) error {
	// For now, execute the entire script as-is
	// In production, you may want to split by semicolons and handle errors per statement
	rows, _, err := conn.QueryMultiRows(sql, []any{}...)
	if err != nil {
		return err
	}
	// Consume rows (if any)
	if rows != nil && len(*rows) > 0 {
		// Just consume them, don't need to do anything
	}
	return nil
}

type quackCheckTokenUDF struct {
	validateToken func(string) (Dict, error)
}

func (q *quackCheckTokenUDF) Config() duckdb.ScalarFuncConfig {
	inputTypeInfo, err := duckdb.NewTypeInfo(duckdb.TYPE_VARCHAR)
	if err != nil {
		panic(fmt.Errorf("failed to create input type info: %w", err))
	}
	resultTypeInfo, err := duckdb.NewTypeInfo(duckdb.TYPE_BOOLEAN)
	if err != nil {
		panic(fmt.Errorf("failed to create result type info: %w", err))
	}
	return duckdb.ScalarFuncConfig{
		InputTypeInfos: []duckdb.TypeInfo{inputTypeInfo},
		ResultTypeInfo: resultTypeInfo,
	}
}

func (q *quackCheckTokenUDF) Executor() duckdb.ScalarFuncExecutor {
	return duckdb.ScalarFuncExecutor{RowExecutor: q.exec}
}

func (q *quackCheckTokenUDF) exec(values []driver.Value) (any, error) {
	if len(values) == 0 {
		return false, fmt.Errorf("quack_check_token requires one argument")
	}
	if values[0] == nil {
		return false, nil
	}
	token, ok := values[0].(string)
	fmt.Println("Validating token in quack_check_token UDF:", token)
	if !ok {
		return false, fmt.Errorf("invalid token type for quack_check_token")
	}
	if q.validateToken == nil {
		return false, fmt.Errorf("token validation function is not configured")
	}
	user, err := q.validateToken(token)
	if err != nil {
		fmt.Println("Token validation failed:", err)
		return false, nil
	}
	fmt.Println("TOKEN USER DATA:", user)
	return true, nil
}

// registerCheckTokenUDF registers a Go scalar UDF for token validation.
func (qm *QuackManager) registerCheckTokenUDF(conn etlx.DBInterface, config Dict) error {
	if qm.validateToken == nil {
		return fmt.Errorf("no token validation function configured for quack manager")
	}

	ddb, ok := conn.(*etlx.DuckDB)
	if !ok {
		return fmt.Errorf("check token udf registration requires a DuckDB connection")
	}

	sqlConn, err := ddb.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("failed to obtain duckdb connection: %w", err)
	}
	defer sqlConn.Close()

	udf := &quackCheckTokenUDF{validateToken: qm.validateToken}
	if err := duckdb.RegisterScalarUDF(sqlConn, "quack_check_token", udf); err != nil {
		return fmt.Errorf("failed to register quack_check_token UDF: %w", err)
	}

	macroSQL := `CREATE MACRO IF NOT EXISTS check_token(session_id, client_token, server_token) AS (quack_check_token(client_token));`
	if _, err := conn.ExecuteQuery(macroSQL); err != nil {
		return fmt.Errorf("failed to create check_token macro: %w", err)
	}
	sql := `SET GLOBAL quack_authentication_function = 'check_token';`
	if _, err := conn.ExecuteQuery(sql); err != nil {
		return fmt.Errorf("failed to set quack_authentication_function: %w", err)
	}

	return nil
}

// fetchQuackServerConfig retrieves Quack server config from admin DB
func (qm *QuackManager) fetchQuackServerConfig(quackServerID int) (Dict, error) {
	query := `
		SELECT *
		FROM quack_server
		WHERE quack_server_id = $1 AND active = TRUE AND excluded = FALSE
	`
	rows, _, err := qm.adminDB.QueryMultiRows(query, []any{quackServerID}...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	if rows == nil || len(*rows) == 0 {
		return nil, fmt.Errorf("quack server %d not found", quackServerID)
	}
	return (*rows)[0], nil
}

// logQuackEvent logs server lifecycle events to quack_logs table
func (qm *QuackManager) logQuackEvent(quackServerID int, event string, status string, port int64, message string, success bool) {
	// This runs asynchronously; any errors are non-fatal
	go func() {
		insertQuery := `
			INSERT INTO quack_logs (quack_server_id, event, status, port, message, log_time, active, excluded)
			VALUES ($1, $2, $3, $4, $5, $6, TRUE, FALSE)
		`
		_, err := qm.adminDB.ExecuteQuery(insertQuery, quackServerID, event, status, port, message, time.Now())
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to log quack event: %v\n", err)
		}
	}()
}

// ListRunningServers returns IDs of all running Quack servers
func (qm *QuackManager) ListRunningServers() []int {
	qm.mux.RLock()
	defer qm.mux.RUnlock()
	ids := make([]int, 0, len(qm.pool))
	for id := range qm.pool {
		ids = append(ids, id)
	}
	return ids
}

// IsServerRunning checks if a specific Quack server is running
func (qm *QuackManager) IsServerRunning(quackServerID int) bool {
	qm.mux.RLock()
	defer qm.mux.RUnlock()
	_, exists := qm.pool[quackServerID]
	return exists
}

// StopAllServers gracefully shuts down all running Quack servers
func (qm *QuackManager) StopAllServers(ctx context.Context) error {
	qm.mux.Lock()
	ids := make([]int, 0, len(qm.pool))
	for id := range qm.pool {
		ids = append(ids, id)
	}
	qm.mux.Unlock()
	for _, id := range ids {
		if err := qm.StopQuackServer(ctx, id); err != nil {
			fmt.Printf("Error stopping quack server %d: %v\n", id, err)
		}
	}
	return nil
}

// Helper functions
func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	default:
		return 0
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (app *application) quackValidateToken(token string) (Dict, error) {
	// fmt.Println("TOKEN:", token)
	identity, err := app.verifyTokenString(token)
	if err != nil {
		fmt.Printf("Err validating token: %s\n", err)
		return nil, err
	}
	return identity, nil
}

func (app *application) startQuackServer(params Dict) Dict {
	quack_server_id := any(nil)
	if _, ok := params["data"].(Dict)["quack_server_id"]; ok {
		quack_server_id = params["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["quack"].(Dict)["quack_server_id"]
	}

	quack_server := any(nil)
	if _, ok := params["data"].(Dict)["name"]; ok {
		quack_server = params["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server = params["data"].(Dict)["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server = params["data"].(Dict)["quack"].(Dict)["name"]
	}
	if quack_server == nil || quack_server == any(nil) {
		if _, ok := params["data"].(Dict)["quack_server"]; ok {
			quack_server = params["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
			quack_server = params["data"].(Dict)["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
			quack_server = params["data"].(Dict)["quack"].(Dict)["quack_server"]
		}
	}
	if (quack_server_id == nil || quack_server_id == any(nil)) && (quack_server == nil || quack_server == any(nil)) {
		return Dict{
			"success": false,
			"msg":     "No Quack Server not especified in the requested!",
		}
	}
	accessQuackData := Dict{}
	if quack_server_id != nil && quack_server_id != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server_id)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	} else if quack_server != nil && quack_server != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	quackInfo := Dict{}
	if quack, ok := accessQuackData["data"].([]Dict); ok {
		if len(quack) > 0 {
			quackInfo = quack[0]
			quack_server_id = quackInfo["quack_server_id"]
			// pass
		} else {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	if !app.quackInstanciated || app.quackManager == nil {
		app.quackManager = NewQuackManager(app.db, app.verifyTokenString)
		app.quackInstanciated = true
	}
	err := app.quackManager.StartQuackServer(context.Background(), app.toInt(quack_server_id), quackInfo)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Failed to start Quack Server: %v", err),
		}
	}
	return Dict{
		"success": true,
		"msg":     "Quack Server started successfully!",
	}
}

func (app *application) stopQuackServer(params Dict) Dict {
	quack_server_id := any(nil)
	if _, ok := params["data"].(Dict)["quack_server_id"]; ok {
		quack_server_id = params["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["quack"].(Dict)["quack_server_id"]
	}

	quack_server := any(nil)
	if _, ok := params["data"].(Dict)["name"]; ok {
		quack_server = params["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server = params["data"].(Dict)["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server = params["data"].(Dict)["quack"].(Dict)["name"]
	}
	if quack_server == nil || quack_server == any(nil) {
		if _, ok := params["data"].(Dict)["quack_server"]; ok {
			quack_server = params["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
			quack_server = params["data"].(Dict)["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
			quack_server = params["data"].(Dict)["quack"].(Dict)["quack_server"]
		}
	}
	if (quack_server_id == nil || quack_server_id == any(nil)) && (quack_server == nil || quack_server == any(nil)) {
		return Dict{
			"success": false,
			"msg":     "No Quack Server not especified in the requested!",
		}
	}
	accessQuackData := Dict{}
	if quack_server_id != nil && quack_server_id != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server_id)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	} else if quack_server != nil && quack_server != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	quackInfo := Dict{}
	if quack, ok := accessQuackData["data"].([]Dict); ok {
		if len(quack) > 0 {
			quackInfo = quack[0]
			quack_server_id = quackInfo["quack_server_id"]
			// pass
		} else {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	if !app.quackInstanciated || app.quackManager == nil {
		app.quackManager = NewQuackManager(app.db, app.verifyTokenString)
		app.quackInstanciated = true
	}
	err := app.quackManager.StopQuackServer(context.Background(), app.toInt(quack_server_id))
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Failed to stop Quack Server: %v", err),
		}
	}
	return Dict{
		"success": true,
		"msg":     "Quack Server stopped successfully!",
	}
}

func (app *application) restartQuackServer(params Dict) Dict {
	quack_server_id := any(nil)
	if _, ok := params["data"].(Dict)["quack_server_id"]; ok {
		quack_server_id = params["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["data"].(Dict)["quack_server_id"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server_id = params["data"].(Dict)["quack"].(Dict)["quack_server_id"]
	}

	quack_server := any(nil)
	if _, ok := params["data"].(Dict)["name"]; ok {
		quack_server = params["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
		quack_server = params["data"].(Dict)["data"].(Dict)["name"]
	} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
		quack_server = params["data"].(Dict)["quack"].(Dict)["name"]
	}
	if quack_server == nil || quack_server == any(nil) {
		if _, ok := params["data"].(Dict)["quack_server"]; ok {
			quack_server = params["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["data"].(Dict); ok {
			quack_server = params["data"].(Dict)["data"].(Dict)["quack_server"]
		} else if _, ok := params["data"].(Dict)["quack"].(Dict); ok {
			quack_server = params["data"].(Dict)["quack"].(Dict)["quack_server"]
		}
	}
	if (quack_server_id == nil || quack_server_id == any(nil)) && (quack_server == nil || quack_server == any(nil)) {
		return Dict{
			"success": false,
			"msg":     "No Quack Server not especified in the requested!",
		}
	}
	accessQuackData := Dict{}
	if quack_server_id != nil && quack_server_id != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server_id)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	} else if quack_server != nil && quack_server != any(nil) {
		accessQuackData = app.getQuackByID(params, quack_server)
		if accessQuackData["success"] != true {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	quackInfo := Dict{}
	if quack, ok := accessQuackData["data"].([]Dict); ok {
		if len(quack) > 0 {
			quackInfo = quack[0]
			quack_server_id = quackInfo["quack_server_id"]
			// pass
		} else {
			return Dict{
				"success": false,
				"msg":     "Quack Server not found or inactive!",
			}
		}
	}
	if !app.quackInstanciated || app.quackManager == nil {
		app.quackManager = NewQuackManager(app.db, app.verifyTokenString)
		app.quackInstanciated = true
	}
	err := app.quackManager.RestartQuackServer(context.Background(), app.toInt(quack_server_id), quackInfo)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("Failed to restart Quack Server: %v", err),
		}
	}
	return Dict{
		"success": true,
		"msg":     "Quack Server restarted successfully!",
	}
}

func (app *application) getQuackByID(params Dict, quack_server_id any) Dict {
	table := "quack_server"
	if _, ok := params["table"].(string); ok {
		table, _ = params["table"].(string)
	} else if _, ok := params["data"].(Dict)["table"].(string); ok {
		table, _ = params["data"].(Dict)["table"].(string)
	}
	database := "ADMIN"
	if _, ok := params["db"].(string); ok {
		database, _ = params["db"].(string)
	} else if _, ok := params["data"].(Dict)["db"].(string); ok {
		database, _ = params["data"].(Dict)["db"].(string)
	} else if _, ok := params["database"].(string); ok {
		database, _ = params["database"].(string)
	} else if _, ok := params["data"].(Dict)["database"].(string); ok {
		database, _ = params["data"].(Dict)["database"].(string)
	}
	if quack_server_id == nil || quack_server_id == any(nil) {
		return Dict{
			"success": false,
			"msg":     "No Quack Server ID found!",
		}
	}
	_aux_params := params
	_aux_params["data"].(Dict)["table"] = table
	_aux_params["data"].(Dict)["db"] = database
	_aux_params["data"].(Dict)["limit"] = any(1.0)
	_aux_params["data"].(Dict)["offset"] = any(0.0)
	_aux_params["data"].(Dict)["filters"] = []any{Dict{
		"field": "quack_server_id",
		"cond":  "=",
		"value": quack_server_id,
	}}
	res := app.read(_aux_params)
	res["quack_server_id"] = quack_server_id
	return res
}

func (app *application) getQuackByName(params Dict, quack_server any) Dict {
	table := "quack_server"
	if _, ok := params["table"].(string); ok {
		table, _ = params["table"].(string)
	} else if _, ok := params["data"].(Dict)["table"].(string); ok {
		table, _ = params["data"].(Dict)["table"].(string)
	}
	database := "ADMIN"
	if _, ok := params["db"].(string); ok {
		database, _ = params["db"].(string)
	} else if _, ok := params["data"].(Dict)["db"].(string); ok {
		database, _ = params["data"].(Dict)["db"].(string)
	} else if _, ok := params["database"].(string); ok {
		database, _ = params["database"].(string)
	} else if _, ok := params["data"].(Dict)["database"].(string); ok {
		database, _ = params["data"].(Dict)["database"].(string)
	}
	if quack_server == nil || quack_server == any(nil) {
		return Dict{
			"success": false,
			"msg":     "No Quack Server ID found!",
		}
	}
	_aux_params := params
	_aux_params["data"].(Dict)["table"] = table
	_aux_params["data"].(Dict)["db"] = database
	_aux_params["data"].(Dict)["limit"] = any(1.0)
	_aux_params["data"].(Dict)["offset"] = any(0.0)
	_aux_params["data"].(Dict)["filters"] = []any{Dict{
		"field": "quack_server",
		"cond":  "=",
		"value": quack_server,
	}}
	res := app.read(_aux_params)
	res["quack_server"] = quack_server
	return res
}
