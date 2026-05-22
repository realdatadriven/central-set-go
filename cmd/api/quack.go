package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/realdatadriven/etlx"
)

// QuackManager manages the lifecycle of in-memory DuckDB instances keyed by server ID
type QuackManager struct {
	pool         map[int]etlx.DBInterface // Keyed by quack_server_id
	mux          sync.RWMutex
	adminDB      etlx.DBInterface // Reference to admin DB for config fetch
	quackConfigs map[int]Dict     // Cache of server configs
	configMu     sync.RWMutex
}

// NewQuackManager creates a new Quack manager instance
func NewQuackManager(adminDB etlx.DBInterface) *QuackManager {
	return &QuackManager{
		pool:         make(map[int]etlx.DBInterface),
		adminDB:      adminDB,
		quackConfigs: make(map[int]Dict),
	}
}

// StartQuackServer initializes an in-memory DuckDB instance and runs startup/attach SQL
func (qm *QuackManager) StartQuackServer(ctx context.Context, quackServerID int) error {
	qm.mux.Lock()
	defer qm.mux.Unlock()

	// Check if already running
	if _, exists := qm.pool[quackServerID]; exists {
		return fmt.Errorf("quack server %d is already running", quackServerID)
	}

	// Fetch config from admin DB
	config, err := qm.fetchQuackServerConfig(quackServerID)
	if err != nil {
		return fmt.Errorf("failed to fetch quack server config: %w", err)
	}

	// Create in-memory DuckDB connection (DSN: "duckdb:")
	conn, err := etlx.GetDB("duckdb:")
	if err != nil {
		return fmt.Errorf("failed to create in-memory duckdb connection: %w", err)
	}

	// Execute startup SQL (e.g., "INSTALL SQLITE; LOAD SQLITE;")
	if config["startup_sql"] != "" {
		if err := qm.executeSQL(conn, config["startup_sql"].(string)); err != nil {
			conn.Close()
			return fmt.Errorf("startup sql failed: %w", err)
		}
	}

	// Execute attach SQL (e.g., "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE); USE adm;")
	if config["attach_sql"] != "" {
		if err := qm.executeSQL(conn, config["attach_sql"].(string)); err != nil {
			conn.Close()
			return fmt.Errorf("attach sql failed: %w", err)
		}
	}

	// Register the security UDF for this server (token validation)
	if err := qm.registerSecurityUDF(conn, quackServerID); err != nil {
		conn.Close()
		return fmt.Errorf("failed to register security udf: %w", err)
	}

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
func (qm *QuackManager) RestartQuackServer(ctx context.Context, quackServerID int) error {
	// Stop
	if err := qm.StopQuackServer(ctx, quackServerID); err != nil {
		return fmt.Errorf("failed to stop server during restart: %w", err)
	}

	// Start
	if err := qm.StartQuackServer(ctx, quackServerID); err != nil {
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

// registerSecurityUDF registers a UDF for token validation
// This can be extended with actual security logic
func (qm *QuackManager) registerSecurityUDF(conn etlx.DBInterface, quackServerID int) error {
	// Fetch the token for this server
	qm.configMu.RLock()
	config, ok := qm.quackConfigs[quackServerID]
	qm.configMu.RUnlock()

	if !ok {
		return fmt.Errorf("server config not found for security udf registration")
	}

	// Create a simple macro or function that validates the token
	// Example: CREATE FUNCTION validate_quack_token(token VARCHAR) RETURNS BOOLEAN AS ...
	// For now, we'll create a simple placeholder that can be extended
	//udfSQL := fmt.Sprintf(`CREATE FUNCTION validate_quack_token(token VARCHAR) RETURNS BOOLEAN AS 'return token == "%s"' LANGUAGE python;`, config.Token)

	// Note: This is pseudo-code. Actual UDF registration depends on DuckDB's capabilities
	// DuckDB supports SQL UDFs but Python UDFs require the python extension
	// For simplicity, we can use a SQL-based check or store tokens in a table

	// Alternative: Store token in a table and query it
	tokenTableSQL := fmt.Sprintf(`
		CREATE TEMPORARY TABLE quack_tokens (server_id INTEGER, token VARCHAR);
		INSERT INTO quack_tokens VALUES (%d, '%s');
	`, quackServerID, config["token"].(string))

	if err := qm.executeSQL(conn, tokenTableSQL); err != nil {
		return fmt.Errorf("failed to create token validation table: %w", err)
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
