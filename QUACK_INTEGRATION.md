# Quack Server Pool Integration Guide

## Environment Variables

Add these to your `.env` file:

```bash
# Enable/disable Quack protocol server
QUACK_ENABLED=false

# Default Quack server port (can be overridden per server in admin_model)
QUACK_PORT=8779
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
    main_sql: "ATTACH 'database/ADMIN.db' AS adm (TYPE SQLITE); USE adm;", 
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

## QuackAPI Services

QuackAPI services are configured separately in the `quack_api` table; the
existing `quack_server` configuration and Quack lifecycle remain unchanged.
Starting a service runs its `startup_sql` followed by `main_sql`. Put route
definitions and the `quackapi_serve(...)` call in `main_sql`. Stopping the
service explicitly runs `shutdown_sql` before closing its DuckDB connection.
The controls use `/quack_api/start/{name}`, `/quack_api/stop/{name}`, and
`/quack_api/restart/{name}`.

For example:

```yaml
startup_sql: "INSTALL quackapi FROM community; LOAD quackapi;"
main_sql: "CREATE ROUTE hello GET '/hello' AS SELECT 'world' AS msg; SELECT * FROM quackapi_serve(9495);"
shutdown_sql: "SELECT * FROM quackapi_stop(9495);"
```

QuackAPI is gated by the existing `QUACK_ENABLED` setting. Use QuackAPI auth
and route policies to control access to shared data.

## Next Steps

- Extend `registerSecurityUDF()` with your actual token validation logic
- Implement admin UI for starting/stopping/restarting servers
- Add metrics/monitoring for server health
- Consider connection pooling within each DuckDB instance for high concurrency
