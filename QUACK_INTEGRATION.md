# Quack Server Pool Integration Guide

## Environment Variables

Add these to your `.env` file:

```bash
# Enable/disable Quack protocol server
QUACK_ENABLED=false

# Default Quack server port (can be overridden per server in admin_model)
QUACK_PORT=8779
```

## Application Integration

### 1. Initialize Quack Manager in main.go

In your `run()` function, after the application is created, initialize the Quack manager:

```go
// Initialize Quack pool if enabled
app.quackEnabled = cfg.quackEnabled
app.quackPool = make(map[int]etlx.DBInterface)
if app.quackEnabled {
    app.quackManager = NewQuackManager(app.db)
    
    // Optionally: Start active Quack servers from config
    // This can be deferred to API calls or done at startup
    // Example:
    // activeServers := app.getActiveQuackServers()
    // for _, serverID := range activeServers {
    //     ctx := context.Background()
    //     if err := app.quackManager.StartQuackServer(ctx, serverID); err != nil {
    //         logger.Error("Failed to start Quack server", "id", serverID, "error", err)
    //     }
    // }
}

// Gracefully shutdown on app close
app.wg.Add(1)
go func() {
    defer app.wg.Done()
    <-ctx.Done()
    if app.quackEnabled && app.quackManager != nil {
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        app.quackManager.StopAllServers(shutdownCtx)
    }
}()
```

### 2. API Endpoints for Quack Management

Add these methods to `api.go` or a new `quack_api.go`:

```go
// POST /api/quack/start/:id
func (app *application) quackStart(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    serverID := app.toInt(id)
    
    if !app.quackEnabled {
        app.errorResponse(w, r, 400, "Quack is not enabled")
        return
    }
    
    ctx := r.Context()
    if err := app.quackManager.StartQuackServer(ctx, serverID); err != nil {
        app.errorResponse(w, r, 400, err.Error())
        return
    }
    
    app.writeJSON(w, http.StatusOK, map[string]any{
        "success": true,
        "message": fmt.Sprintf("Quack server %d started", serverID),
    })
}

// POST /api/quack/stop/:id
func (app *application) quackStop(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    serverID := app.toInt(id)
    
    if !app.quackEnabled {
        app.errorResponse(w, r, 400, "Quack is not enabled")
        return
    }
    
    ctx := r.Context()
    if err := app.quackManager.StopQuackServer(ctx, serverID); err != nil {
        app.errorResponse(w, r, 400, err.Error())
        return
    }
    
    app.writeJSON(w, http.StatusOK, map[string]any{
        "success": true,
        "message": fmt.Sprintf("Quack server %d stopped", serverID),
    })
}

// POST /api/quack/restart/:id
func (app *application) quackRestart(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    serverID := app.toInt(id)
    
    if !app.quackEnabled {
        app.errorResponse(w, r, 400, "Quack is not enabled")
        return
    }
    
    ctx := r.Context()
    if err := app.quackManager.RestartQuackServer(ctx, serverID); err != nil {
        app.errorResponse(w, r, 400, err.Error())
        return
    }
    
    app.writeJSON(w, http.StatusOK, map[string]any{
        "success": true,
        "message": fmt.Sprintf("Quack server %d restarted", serverID),
    })
}

// GET /api/quack/status/:id
func (app *application) quackStatus(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    serverID := app.toInt(id)
    
    if !app.quackEnabled {
        app.errorResponse(w, r, 400, "Quack is not enabled")
        return
    }
    
    isRunning := app.quackManager.IsServerRunning(serverID)
    
    app.writeJSON(w, http.StatusOK, map[string]any{
        "running": isRunning,
        "server_id": serverID,
    })
}

// GET /api/quack/list
func (app *application) quackList(w http.ResponseWriter, r *http.Request) {
    if !app.quackEnabled {
        app.errorResponse(w, r, 400, "Quack is not enabled")
        return
    }
    
    running := app.quackManager.ListRunningServers()
    
    app.writeJSON(w, http.StatusOK, map[string]any{
        "running_servers": running,
    })
}
```

### 3. Extend Security UDF for Your Needs

The `registerSecurityUDF()` method in `quack.go` currently creates a token table. Extend it like this:

```go
// In quack.go, extend registerSecurityUDF:
func (qm *QuackManager) registerSecurityUDFExtended(conn etlx.DBInterface, quackServerID int) error {
    // Fetch token
    qm.configMu.RLock()
    config := qm.quackConfigs[quackServerID]
    qm.configMu.RUnlock()
    
    // Create security context table
    securitySQL := fmt.Sprintf(`
        CREATE TEMPORARY TABLE quack_security (
            server_id INTEGER,
            token VARCHAR,
            max_connections INTEGER,
            allow_write BOOLEAN
        );
        INSERT INTO quack_security VALUES 
            (%d, '%s', 100, FALSE);
            
        CREATE FUNCTION validate_token(t VARCHAR) RETURNS BOOLEAN AS 
            'SELECT EXISTS(SELECT 1 FROM quack_security WHERE token = t)';
            
        CREATE FUNCTION get_user_filters() RETURNS VARCHAR AS
            'SELECT current_user';
    `, quackServerID, config.Token)
    
    return qm.executeSQL(conn, securitySQL)
}
```

### 4. Use from Access Layer

In your access control layer (similar to Arrow Flight):

```go
// Add this to your access.go for RLA checks on Quack servers
func (app *application) quack_table_access(params map[string]any, tables []any) map[string]any {
    var quackServerID int
    if _, ok := params["quack_server"].(map[string]any)["quack_server_id"]; ok {
        quackServerID = app.toInt(params["quack_server"].(map[string]any)["quack_server_id"])
    }
    
    // Get the in-memory connector
    conn, err := app.quackManager.GetQuackConnector(quackServerID)
    if err != nil {
        return map[string]any{
            "success": false,
            "msg":     err.Error(),
        }
    }
    
    // Use conn to query data with table access controls
    // Same RLA pattern as Arrow Flight...
    
    return map[string]any{
        "success": true,
        "msg":     "access granted",
        "data":    data,
    }
}
```

## Database Startup Hook

To auto-start Quack servers marked as active:

```go
// In your database/app initialization
func (app *application) initializeQuackServers() error {
    if !app.quackEnabled {
        return nil
    }
    
    query := `
        SELECT quack_server_id FROM quack_server 
        WHERE active = TRUE AND excluded = FALSE
    `
    
    rows, _, err := app.db.QueryMultiRows(query, []any{}...)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    for _, row := range *rows {
        serverID := app.toInt(row["quack_server_id"])
        if err := app.quackManager.StartQuackServer(ctx, serverID); err != nil {
            app.logger.Error("Failed to start Quack server", "id", serverID, "error", err)
        }
    }
    
    return nil
}

// Call this during app startup
func (app *application) run() error {
    // ... existing setup ...
    
    if app.quackEnabled {
        if err := app.initializeQuackServers(); err != nil {
            app.logger.Error("Failed to initialize Quack servers", "error", err)
        }
    }
    
    // ... rest of run ...
}
```

## Example Quack Server Configuration

From `admin_model.md`:

```yaml
- {
    quack_server_id: 1, 
    quack_name: "Quack Admin DB", 
    quack_desc: "Expose ADMIN DB via DuckDB Quack", 
    port: 8779, 
    token: "your-secret-token",
    protocol: quack, 
    startup_sql: "INSTALL SQLITE; LOAD SQLITE;", 
    attach_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE); USE adm;", 
    shutdown_sql: "USE memory; DETACH adm;", 
    status: offline, 
    active: false, 
    app_id: 1, 
    user_id: 1, 
    excluded: false
  }
```

## Key Design Points

1. **In-Memory Instances**: Each Quack server runs its own in-memory DuckDB instance
2. **Pool Keyed by Server ID**: Easy lookup and lifecycle management
3. **Thread-Safe**: Uses `sync.RWMutex` for concurrent access
4. **Extensible Security**: UDF framework ready for custom token validation
5. **Audit Trail**: All lifecycle events logged to `quack_logs` table
6. **Graceful Shutdown**: Runs `shutdown_sql` before closing connections
7. **No Global State**: Each instance is isolated and can be independently started/stopped

## Next Steps

- Extend `registerSecurityUDF()` with your actual token validation logic
- Implement admin UI for starting/stopping/restarting servers
- Add metrics/monitoring for server health
- Consider connection pooling within each DuckDB instance for high concurrency
