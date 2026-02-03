/*This component takes the Arrow Flight configuration defined in the Admin database and uses **airport-go** to create a Flight server backed by **DuckDB**.

At runtime, it initializes a DuckDB instance and serves the tables defined in the configuration using DuckDB’s **Arrow integration**, exposing them through the Arrow Flight protocol.

The lifecycle of each endpoint is fully driven by configuration:

* `startup_sql` is executed to load extensions and initialize dependencies
* `main_sql` attaches the underlying data sources and exposes tables
* `shutdown_sql` is executed to clean up resources (for example, detaching databases)

All requests are authenticated using the `validateToken` function, ensuring that Arrow Flight follows the same security and access rules as the rest of the platform.

The original design aimed to reuse a **global `duckdb.Connector`** so multiple Flight clients could share the same DuckDB instance. However, due to concurrency issues when sharing connectors across goroutines, the current implementation creates a **new DuckDB connector per scan request**.

While this approach introduces some overhead, real-world testing shows that performance remains acceptable for analytical workloads—especially when compared to the OData v4 API, which is better suited for smaller or transactional queries.

Further optimizations around connector reuse and resource management are planned to improve performance and efficiency over time.

*/

package flight

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/duckdb/duckdb-go/v2"
	"github.com/realdatadriven/etlx"
	"google.golang.org/grpc"

	airport "github.com/hugr-lab/airport-go"
	"github.com/hugr-lab/airport-go/catalog"
	"github.com/hugr-lab/airport-go/filter"

	//duckarrow "github.com/duckdb/duckdb-go/v2/arrow"

	"google.golang.org/grpc/credentials"
)

// FlightManager is the interface the server uses to start/stop the FlightSQL server.
type FlightManager interface {
	Start(listenAddr string) error
	Stop(ctx context.Context) error
}

// AirportAdapter implements FlightManager using hugr-lab/airport-go.
type AirportAdapter struct {
	validateToken func(token string) (string, error)
	table_access  func(params map[string]any, tables []any) map[string]any               // checks is a user has access to a specifc table in this case arrow_flight_table mapping tables
	rla_access    func(params map[string]any, tables []any, row_id []any) map[string]any // checks table in arrow_flight_table is accessible to the user
	grpcSrv       *grpc.Server
	listener      net.Listener
	mem           memory.Allocator
	catalog       catalog.Catalog
	cfg           []map[string]any
	shutdownc     chan struct{}
}

// NewAirportAdapter constructs the adapter with the provided DDB.
func NewAirportAdapter(
	config []map[string]any,
	validateToken func(token string) (string, error),
	table_access func(params map[string]any, tables []any) map[string]any,
	rla_access func(params map[string]any, tables []any, row_id []any) map[string]any,
) *AirportAdapter {
	return &AirportAdapter{
		validateToken: validateToken,
		table_access:  table_access,
		rla_access:    rla_access,
		mem:           memory.DefaultAllocator,
		cfg:           config,
		shutdownc:     make(chan struct{}),
	}
}

// Start builds an airport-go catalog from the DuckDB schemas and tables discovered via the manager.
// It then creates a gRPC server and registers the airport server.
func (a *AirportAdapter) Start(listenAddr string) error {
	// Build catalog using airport.NewCatalogBuilder()
	//builder := airport.NewCatalogBuilder()
	builder := NewCatalogBuilder().Dynamic()
	db, err := duckdb.NewConnector("", nil)
	if err != nil {
		return err
	}
	defer db.Close()
	conn := sql.OpenDB(db)
	defer conn.Close()
	_etlx := etlx.ETLX{}
	// For each schema defined in config, discover its tables and add them as SimpleTable entries.
	for _, s := range a.cfg {
		schemaName := s["flight_schema"].(string)
		// create a schema builder for this schema
		sb := builder.Schema(schemaName)
				.Comment("Main application schema")
		// execute startup_sql
		var main_sql string
		if startup_sql, ok := s["startup_sql"].(string); ok {
			startup_sql = _etlx.ReplaceEnvVariable(startup_sql)
			//fmt.Printf("%s: %s\n", s["arrow_flight"], startup_sql)
			_, err := conn.ExecContext(context.Background(), startup_sql)
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], startup_sql, err)
			}
			main_sql = _etlx.ReplaceEnvVariable(s["main_sql"].(string))
			//fmt.Printf("%s: %s\n", s["arrow_flight"], main_sql)
			_, err = conn.ExecContext(context.Background(), main_sql)
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], main_sql, err)
			}
		}
		// check if main_sql has use <schema>; is not do use <schema>;
		/*/ fmt.Println(strings.Contains(strings.ToLower(main_sql), fmt.Sprintf("use %s;", strings.ToLower(schemaName))), main_sql)
		if !strings.Contains(strings.ToLower(main_sql), fmt.Sprintf("use %s;", strings.ToLower(schemaName))) {
			_, err = conn.ExecContext(context.Background(), fmt.Sprintf("USE %s;", schemaName))
			if err != nil {
				return fmt.Errorf("use schema %s: %w", schemaName, err)
			}
		}*/
		//q := `SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name`
		table_discover_sql := `select table_name from duckdb_tables`
		if _table_discover_sql, ok := s["table_discover_sql"].(string); ok {
			table_discover_sql = _etlx.ReplaceEnvVariable(_table_discover_sql)
			//fmt.Println("table_discover_sql:", table_discover_sql)
		}
		rows, err := conn.QueryContext(context.Background(), table_discover_sql, schemaName)
		if err != nil {
			return fmt.Errorf("query tables: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var tname string
			if err := rows.Scan(&tname); err != nil {
				return fmt.Errorf("scan table name: %w", err)
			}
			//fmt.Println(tname)
			// if s["tables"].(map[string]any) exists and it length > 0 and tname not in s["tables"].(map[string]any) skip
			if tables, ok := s["tables"].(map[string]any); ok && len(tables) > 0 {
				found := false
				if _, ok := tables[tname]; ok {
					found = true
					//rla = a.rla_tables(map[string]any{}, []any{tname})
				}
				if !found {
					continue
				}
			}
			// // if s["tables"].(map[string]any) exists and it length > 0 and tname is in s["tables"].(map[string]any), and s["tables"].(map[string]any)[tname].(map[string]any)[fields] has length > 0, filter arrowSchema to only include those fields
			_fields := []string{}
			if tables, ok := s["tables"].(map[string]any); ok && len(tables) > 0 {
				if tableConf, ok := tables[tname].(map[string]any); ok {
					if fields, ok := tableConf["fields"].(map[string]any); ok && len(fields) > 0 {
						for field, _ := range fields {
							_fields = append(_fields, fmt.Sprintf(`"%s"`, field))
						}
					}
				}
			}
			if len(_fields) == 0 {
				_fields = []string{"*"}
			}
			query := fmt.Sprintf(`SELECT %s FROM %s."%s" LIMIT 0`, strings.Join(_fields, ","), schemaName, tname)
			// if s["table_scan_tmpl_sql"] exists use it to build the query
			if table_scan_tmpl_sql, ok := s["table_scan_tmpl_sql"].(string); ok {
				table_scan_tmpl_sql = _etlx.ReplaceEnvVariable(table_scan_tmpl_sql)
				query = strings.ReplaceAll(table_scan_tmpl_sql, "{{table_name}}", tname)
				query = strings.ReplaceAll(query, "{{schema_name}}", schemaName)
				query = strings.ReplaceAll(query, "{{fields}}", strings.Join(_fields, ","))
				query = fmt.Sprintf(query+" LIMIT 0", strings.Join(_fields, ","), schemaName, tname)
				//fmt.Println("table_scan_tmpl_sql query:", query)
			}
			conn, err := db.Connect(context.Background())
			if err != nil {
				return err
			}
			defer conn.Close()
			arrow, err := duckdb.NewArrowFromConn(conn)
			if err != nil {
				return err
			}
			rdr, err := arrow.QueryContext(context.Background(), query)
			if err != nil {
				return err
			}
			defer rdr.Release()
			arrowSchema := rdr.Schema()
			//fmt.Println(tname, arrowSchema)
			scanFn := a.makeScanFunc(a.mem, schemaName, tname, arrowSchema, s)
			// register simple table under current schema builder
			sb.SimpleTable(airport.SimpleTableDef{
				Name:     tname,
				Comment:  tname,
				Schema:   arrowSchema,
				ScanFunc: scanFn,
			})
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("tables rows err: %w", err)
		}
		// execute shutdown_sql
		if shutdown_sql, ok := s["shutdown_sql"].(string); ok {
			shutdown_sql := _etlx.ReplaceEnvVariable(shutdown_sql)
			//fmt.Printf("%s: %s\n", s["arrow_flight"], shutdown_sql)
			_, err = conn.ExecContext(context.Background(), shutdown_sql)
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], shutdown_sql, err)
			}
		}
		// conn.ExecContext(context.Background(), "USE memory;")
	}
	conn.Close()
	db.Close()
	cat, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build catalog: %w", err)
	}
	a.catalog = cat
	// Create grpc server and register airport server
	debugLevel := slog.LevelInfo
	if os.Getenv("ARROW_FLIGHT_LOG_LEVEL") == "LevelInfo" {
		debugLevel = slog.LevelInfo
	} else if os.Getenv("ARROW_FLIGHT_LOG_LEVEL") == "LevelWarn" {
		debugLevel = slog.LevelWarn
	} else if os.Getenv("ARROW_FLIGHT_LOG_LEVEL") == "LevelError" {
		debugLevel = slog.LevelError
	} else if os.Getenv("ARROW_FLIGHT_LOG_LEVEL") == "LevelDebug" {
		debugLevel = slog.LevelDebug
	}
	config := airport.ServerConfig{
		Catalog:  cat,
		Auth:     airport.BearerAuth(a.validateToken),
		Address:  listenAddr,
		LogLevel: &debugLevel,
	}
	// https://github.com/hugr-lab/airport-go/blob/main/examples/tls/main.go
	// Load TLS credentials
	creds, err := loadTLSCredentialsV2() //loadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	opts := airport.ServerOptions(config)
	if creds != nil {
		// fmt.Println("TLS CREDS:", creds)
		opts = append(opts, grpc.Creds(creds))
	}
	a.grpcSrv = grpc.NewServer(opts...)
	if err := airport.NewServer(a.grpcSrv, config); err != nil {
		return fmt.Errorf("airport.NewServer failed: %w", err)
	}
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", listenAddr, err)
	}
	a.listener = lis
	// Serve in a goroutine
	go func() {
		log.Printf("[flight] Airport server listening on %s", listenAddr)
		if err := a.grpcSrv.Serve(lis); err != nil {
			log.Printf("[flight] grpc serve error: %v", err)
		}
		close(a.shutdownc)
	}()
	return nil
}

// Stop gracefully stops the airport-go server.
func (a *AirportAdapter) Stop(ctx context.Context) error {
	if a.grpcSrv != nil {
		a.grpcSrv.GracefulStop()
		// wait until serve goroutine exits or context times out
		select {
		case <-a.shutdownc:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (a *AirportAdapter) contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
func IdentityJSONToMap(identity string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(identity), &m); err != nil {
		return nil, err
	}
	return m, nil
}
func (a *AirportAdapter) containsInt(slice []any, element any) bool {
	for _, v := range slice {
		if v.(int64) == element.(int64) {
			return true
		}
	}
	return false
}
func (a *AirportAdapter) getRLAIds(rla_access []map[string]any, table, access_type string) []any {
	data := []any{}
	for _, v := range rla_access {
		_access := false
		switch v[access_type].(type) {
		case bool:
			_access = v[access_type].(bool)
		case int:
			_access = v[access_type].(int) == 1
		case float32:
			_access = v[access_type].(float32) == 1
		case float64:
			_access = v[access_type].(float64) == 1
		case int64:
			_access = v[access_type].(int64) == 1
		case int32:
			_access = v[access_type].(int32) == 1
		case string:
			i, err := strconv.Atoi(v[access_type].(string))
			if err != nil {
				_access = false
			}
			_access = i == 1
		default:
			_access = false // or handle error
		}
		if v["table"].(string) == table && _access {
			data = append(data, v["row_id"])
		}
	}
	return data
}

// The idea is the use the global duckdb.Connector to create connections
// for each scan request, so we can share the same DB across multiple
// connections/flight clients.
// but im im passing the all configuration and making the connection and initializing the sturup and shutdown
// all over again because of some problems with the shared connector across goroutines.
// This needs to be improved later for performance and resource usage.
func (a *AirportAdapter) makeScanFunc(mem memory.Allocator, schemaName, tableName string, aSchema *arrow.Schema, conf map[string]any) func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	return func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
		// CHECK IF TABLE IS ALLOWED MEANING IF THERE IS TABLE MAPPING, ONLLY THOSE ARE EXPOSED
		arrow_flight_id := conf["arrow_flight_id"]
		var arrow_flight_table_id any
		if tables, ok := conf["tables"].(map[string]any); ok && len(tables) > 0 {
			found := false
			if _, ok := tables[tableName].(map[string]any); ok {
				found = true
				arrow_flight_table_id = tables[tableName].(map[string]any)["arrow_flight_table_id"]
			}
			if !found {
				return nil, fmt.Errorf("table %s not allowed", tableName)
			}
		}
		// CHECK IF USER HAS ACCESS TO SPECIFC TABLE
		// fmt.Println("arrow_flight_id:", arrow_flight_id, "arrow_flight_table_id:", arrow_flight_table_id)
		user, err := IdentityJSONToMap(airport.IdentityFromContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("Error getting Identity From Context %v", err)
		}
		//fmt.Println(user)
		// CHECK IF USER HAS ASSCESS TO THE SCHEMA
		rla_tables := []string{}
		if _, ok := conf["rla_tables"].([]string); ok {
			rla_tables = conf["rla_tables"].([]string)
		}
		// check if arrow_flight is in rla_tables
		//fmt.Printf("%T", conf["rla_tables"])
		//fmt.Println(conf["rla_tables"], rla_tables, " CHECK CONTAINS ", "arrow_flight")
		//schema_permissions := map[string]any{}
		//schema_table_permissions := map[string]any{}
		scopes_access := map[string]any{}
		fields_access := map[string]any{}
		if user["role_id"] != any(1) && (a.contains(rla_tables, "arrow_flight") ||
			a.contains(rla_tables, "arrow_flight_table") ||
			a.contains(rla_tables, "arrow_flight_table_field") ||
			a.contains(rla_tables, "arrow_flight_table_scope")) {
			/*schema_access := a.table_access(
				map[string]any{
					"app":  map[string]any{"app_id": 1},
					"data": map[string]any{},
					"user": user,
				}, []any{"arrow_flight", "arrow_flight_table"})
			if !schema_access["success"].(bool) {
				fmt.Println("schema_access:", tableName, schema_access["msg"])
			} else if _, ok := schema_access["data"]; ok {
				_permissions := schema_access["data"].(map[string]any)
				if _, ok := _permissions["arrow_flight"]; ok {
					schema_permissions = _permissions["arrow_flight"].(map[string]any)
				} else {
					return nil, fmt.Errorf("Access denied access to the schema %s!", schemaName)
				}
				if _, ok := _permissions["arrow_flight_table"]; ok {
					schema_table_permissions = _permissions["arrow_flight_table"].(map[string]any)
				} else {
					return nil, fmt.Errorf("Access denied access to the table %s from schema %s!", schemaName, tableName)
				}
			}*/
			// setup for fields access
			arrow_flight_table_field_ids := []any{}
			_field_id_map := map[string]any{}
			if a.contains(rla_tables, "arrow_flight_table_field") {
				if tables, ok := conf["tables"].(map[string]any); ok && len(tables) > 0 {
					if tableConf, ok := tables[tableName].(map[string]any); ok {
						if fields, ok := tableConf["fields"].(map[string]any); ok && len(fields) > 0 {
							for field := range fields {
								arrow_flight_table_field_ids = append(arrow_flight_table_field_ids, fields[field].(map[string]any)["arrow_flight_table_field_id"])
								_field_id_map[field] = fields[field].(map[string]any)["arrow_flight_table_field_id"]
							}
						}
					}
				}
			}
			// SCOPE ACCESS
			arrow_flight_table_scope_ids := []any{}
			_scope_id_map := map[string]any{}
			if a.contains(rla_tables, "arrow_flight_table_scope") {
				if tables, ok := conf["tables"].(map[string]any); ok && len(tables) > 0 {
					if tableConf, ok := tables[tableName].(map[string]any); ok {
						if scopes, ok := tableConf["scopes"].(map[string]any); ok && len(scopes) > 0 {
							for scope := range scopes {
								arrow_flight_table_scope_ids = append(arrow_flight_table_scope_ids, scopes[scope].(map[string]any)["arrow_flight_table_scope_id"])
								_scope_id_map[scope] = scopes[scope].(map[string]any)["arrow_flight_table_scope_id"]
							}
							fmt.Println("SCOPES:", arrow_flight_table_scope_ids, _scope_id_map)
						}
					}
				}
			}
			// RUN THE ACCESSS
			_ids := append([]any{arrow_flight_id, arrow_flight_table_id}, arrow_flight_table_field_ids...)
			_ids = append(_ids, arrow_flight_table_scope_ids...)
			rla_access := a.rla_access(
				map[string]any{
					"app":  map[string]any{"app_id": 1},
					"data": map[string]any{},
					"user": user,
				},
				[]any{"arrow_flight", "arrow_flight_table", "arrow_flight_table_field", "arrow_flight_table_scope"},
				_ids,
			)
			//fmt.Println("rla_access:", tableName, rla_access)
			if !rla_access["success"].(bool) {
				fmt.Println("rla_access:", tableName, rla_access["msg"])
			} else if _, ok := rla_access["data"]; ok {
				_permissions := rla_access["data"].(map[string]any)
				if _, ok := _permissions["arrow_flight"]; ok {
					_schema := _permissions["arrow_flight"].([]map[string]any)
					ids := a.getRLAIds(_schema, "arrow_flight", "read")
					//fmt.Println("arrow_flight row ids:", ids)
					if !a.containsInt(ids, arrow_flight_id) {
						return nil, fmt.Errorf("Access denied to the schema \"%s\"!", schemaName)
					}
				} else if a.contains(rla_tables, "arrow_flight") {
					return nil, fmt.Errorf("Access denied to the schema \"%s\"!", schemaName)
				}
				if _, ok := _permissions["arrow_flight_table"]; ok {
					_schema_table := _permissions["arrow_flight_table"].([]map[string]any)
					ids := a.getRLAIds(_schema_table, "arrow_flight_table", "read")
					//fmt.Println("arrow_flight_table row ids:", _schema_table, ids)
					if !a.containsInt(ids, arrow_flight_table_id) {
						return nil, fmt.Errorf("Access denied to the table \"%s\" from schema \"%s\"!", tableName, schemaName)
					}
				} else if a.contains(rla_tables, "arrow_flight_table") {
					return nil, fmt.Errorf("Access denied to the table \"%s\" from schema \"%s\"!", tableName, schemaName)
				}
				// FIELDS ACCESS
				if _, ok := _permissions["arrow_flight_table_field"]; ok {
					_schema_table_field := _permissions["arrow_flight_table_field"].([]map[string]any)
					ids := a.getRLAIds(_schema_table_field, "arrow_flight_table_field", "read")
					if len(ids) == 0 {
						return nil, fmt.Errorf("Access denied, fileds level access are required on table \"%s\" from schema \"%s\", and you have access to none!", tableName, schemaName)
					}
					//fmt.Println("arrow_flight_table_field row ids:", _schema_table, ids)
					for field, arrow_flight_table_field_id := range _field_id_map {
						if a.containsInt(ids, arrow_flight_table_field_id) {
							fields_access[field] = true
							//return nil, fmt.Errorf("Access denied to the field \"%s\" in the table \"%s\" from schema \"%s\"!", field, tableName, schemaName)
						} else {
							fields_access[field] = false
						}
					}
				} else if a.contains(rla_tables, "arrow_flight_table_field") {
					for field := range _field_id_map {
						fields_access[field] = false
					}
					return nil, fmt.Errorf("Access denied to the fields on table \"%s\" from schema \"%s\"!", tableName, schemaName)
				}
				// SCOPE ACCESS
				if _, ok := _permissions["arrow_flight_table_scope"]; ok {
					_schema_table_scope := _permissions["arrow_flight_table_scope"].([]map[string]any)
					ids := a.getRLAIds(_schema_table_scope, "arrow_flight_table_scope", "read")
					if len(ids) == 0 {
						return nil, fmt.Errorf("Access denied, scopes are required on table \"%s\" from schema \"%s\", and you have access to none!", tableName, schemaName)
					}
					for scope, arrow_flight_table_scope_id := range _scope_id_map {
						if a.containsInt(ids, arrow_flight_table_scope_id) {
							scopes_access[scope] = true
						} else {
							scopes_access[scope] = false
						}
					}
				} else if a.contains(rla_tables, "arrow_flight_table_scope") {
					for scope := range _scope_id_map {
						scopes_access[scope] = false
					}
					return nil, fmt.Errorf("Access denied, scopes are required on table \"%s\" from schema \"%s\", and you have access to none!", tableName, schemaName)
				}
			}
			//fmt.Println(schema_permissions, schema_table_permissions)
		}
		// CONNECTION
		_etlx := etlx.ETLX{}
		db, err := duckdb.NewConnector("", nil)
		if err != nil {
			return nil, err
		}
		conn := sql.OpenDB(db)
		defer conn.Close()
		//fmt.Println(_sql, fligths)
		defer func() { // EXECUTE SHUTDOWN SQL ON EXIT
			conn := sql.OpenDB(db)
			defer conn.Close()
			if shutdown_sql, ok := conf["shutdown_sql"].(string); ok {
				//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["shutdown_sql"])
				shutdown_sql = _etlx.ReplaceEnvVariable(shutdown_sql)
				_, err := conn.ExecContext(context.Background(), shutdown_sql)
				if err != nil {
					fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], shutdown_sql, err)
				}
			}
			db.Close()
		}()
		// EXECUTE STARTUP SQL
		if startup_sql, ok := conf["startup_sql"].(string); ok {
			//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["startup_sql"])
			startup_sql = _etlx.ReplaceEnvVariable(startup_sql)
			_, err := conn.ExecContext(context.Background(), startup_sql)
			if err != nil {
				fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], startup_sql, err)
			}
			//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["main_sql"])
			main_sql := _etlx.ReplaceEnvVariable(conf["main_sql"].(string))
			_, err = conn.ExecContext(context.Background(), main_sql)
			if err != nil {
				fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], main_sql, err)
			}
		}
		// MAP FIELD TYPES
		_fields_type := map[string]any{}
		if len(fields_access) > 0 {
			query := fmt.Sprintf("DESC %s.\"%s\"", schemaName, tableName)
			query = fmt.Sprintf(`SELECT column_name, data_type FROM duckdb_columns WHERE database_name = '%s' AND table_name = '%s'`, schemaName, tableName)
			//query = `SELECT column_name, data_type FROM duckdb_columns WHERE database_name = ? AND table_name = ?`
			//fmt.Println(query)
			rows, err := conn.QueryContext(context.Background(), query, schemaName, tableName)
			if err != nil {
				return nil, fmt.Errorf("query tables: %w", err)
			}
			defer rows.Close()
			for rows.Next() {
				var name, dtype string
				if err := rows.Scan(&name, &dtype); err != nil {
					return nil, fmt.Errorf("scan column: %w", err)
				}
				_fields_type[name] = dtype
			}
			if err := rows.Err(); err != nil {
				return nil, fmt.Errorf("tables rows err: %w", err)
			}
		}
		// PREPARE THE FIELDS TO USE
		_fields := []string{}
		if tables, ok := conf["tables"].(map[string]any); ok && len(tables) > 0 {
			if tableConf, ok := tables[tableName].(map[string]any); ok {
				if fields, ok := tableConf["fields"].(map[string]any); ok && len(fields) > 0 {
					for field := range fields {
						if len(fields_access) > 0 {
							if _access, ok := fields_access[field].(bool); ok {
								if _access {
									_fields = append(_fields, fmt.Sprintf(`"%s"`, field))
								} else if _, ok := _fields_type[field]; ok {
									_fields = append(_fields, fmt.Sprintf(`NULL::%s AS "%s"`, _fields_type[field], field))
								} else {
									_fields = append(_fields, fmt.Sprintf(`NULL AS "%s"`, field))
								}
							} else if _, ok := _fields_type[field]; ok {
								_fields = append(_fields, fmt.Sprintf(`NULL::%s AS "%s"`, _fields_type[field], field))
							} else {
								_fields = append(_fields, fmt.Sprintf(`NULL AS "%s"`, field))
							}
						} else {
							_fields = append(_fields, fmt.Sprintf(`"%s"`, field))
						}

					}
				}
			}
		}
		if len(_fields) == 0 {
			_fields = []string{"*"}
		}
		// PREPARE THE scopes TO USE
		_scopes := []string{}
		if tables, ok := conf["tables"].(map[string]any); ok && len(tables) > 0 {
			if tableConf, ok := tables[tableName].(map[string]any); ok {
				if scopes, ok := tableConf["scopes"].(map[string]any); ok && len(scopes) > 0 {
					for scope := range scopes {
						scope_sql := scopes[scope].(map[string]any)["arrow_flight_table_scope_sql"].(string)
						scope_sql = _etlx.ReplaceEnvVariable(scope_sql)
						if len(scopes_access) > 0 {
							if _access, ok := scopes_access[scope].(bool); ok {
								if _access {
									_scopes = append(_scopes, fmt.Sprintf(`%s`, scope_sql))
								}
							}
						} else {
							_scopes = append(_scopes, fmt.Sprintf(`%s`, scope_sql))
						}
					}
				}
			}
		}
		// BUILD THE FINAL QUERY
		query := fmt.Sprintf("SELECT %s FROM %s.\"%s\"", strings.Join(_fields, ","), schemaName, tableName)
		if table_scan_tmpl_sql, ok := conf["table_scan_tmpl_sql"].(string); ok {
			table_scan_tmpl_sql = _etlx.ReplaceEnvVariable(table_scan_tmpl_sql)
			query = strings.ReplaceAll(table_scan_tmpl_sql, "{{table_name}}", tableName)
			query = strings.ReplaceAll(query, "{{schema_name}}", schemaName)
			query = strings.ReplaceAll(query, "{{fields}}", strings.Join(_fields, ","))
			query = fmt.Sprintf(query, strings.Join(_fields, ","), schemaName, tableName)
			//fmt.Println("table_scan_tmpl_sql query:", query)
		}
		hasFilters := false
		if opts.Filter != nil {
			// Parse filter JSON
			fp, err := filter.Parse(opts.Filter)
			if err != nil {
				return nil, err
			}
			// Encode to SQL WHERE clause
			enc := filter.NewDuckDBEncoder(nil)
			whereClause := enc.EncodeFilters(fp)
			//fmt.Printf("Filter applied on %s.%s: %s\n", schemaName, tableName, whereClause)
			if whereClause != "" {
				hasFilters = true
				query = fmt.Sprintf("%s WHERE (%s)", query, whereClause)
			}
			// Use whereClause with your database query
		}
		if len(_scopes) > 0 {
			_scopes_cond := strings.Join(_scopes, " AND ")
			if !hasFilters {
				query = fmt.Sprintf("%s WHERE (%s)", query, _scopes_cond)
			} else {
				query = fmt.Sprintf("%s AND (%s)", query, _scopes_cond)
			}
		}
		fmt.Println(query, scopes_access)
		conn2, err := db.Connect(context.Background())
		if err != nil {
			return nil, err
		}
		//defer conn.Close()
		arrow, err := duckdb.NewArrowFromConn(conn2)
		if err != nil {
			return nil, err
		}
		rdr, err := arrow.QueryContext(context.Background(), query)
		if err != nil {
			conn2.Close()
			return nil, err
		}
		/*/defer rdr.Release()
		for rdr.Next() {
			rec := rdr.Record()
			// rec is an Arrow RecordBatch
			fmt.Println("rows:", rec.NumRows())
		}*/
		/*if shutdown_sql, ok := conf["shutdown_sql"].(string); ok {
			//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["shutdown_sql"])
			shutdown_sql = _etlx.ReplaceEnvVariable(shutdown_sql)
			_, err := conn.ExecContext(context.Background(), shutdown_sql)
			if err != nil {
				fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], shutdown_sql, err)
			}
		}*/
		//return nil, nil
		return &connBoundRecordReader{
			RecordReader: rdr,
			conn:         conn2,
		}, nil
	}
}

// https://github.com/hugr-lab/airport-go/blob/main/examples/tls/main.go
// loadTLSCredentials loads TLS credentials from files.
// In production, use proper certificate management.
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server certificate and key
	enableTLS := strings.ToLower(os.Getenv("ENABLE_TLS")) == "true"
	if enableTLS {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")
		caFile := os.Getenv("TLS_CA_CERT_FILE")
		if certFile == "" || keyFile == "" {
			return nil, fmt.Errorf("ENABLE_TLS is true but TLS_CERT_FILE or TLS_KEY_FILE or TLS_CA_CERT_FILE is not set %s", "")
		}
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load server cert: %w", err)
		}
		// Load CA certificate for mutual TLS (optional)
		//if caFile != "" {
		certPool := x509.NewCertPool()
		if caCert, err := os.ReadFile(caFile); err == nil {
			if !certPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to add CA cert to pool")
			}
		}
		//}
		// Configure TLS
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.NoClientCert, // Change to tls.RequireAndVerifyClientCert for mTLS
			ClientCAs:    certPool,
			MinVersion:   tls.VersionTLS12,
		}
		return credentials.NewTLS(tlsConfig), nil
	}
	return nil, nil
}

func loadTLSCredentialsV2() (credentials.TransportCredentials, error) {
	enableTLS := strings.ToLower(os.Getenv("ENABLE_TLS")) == "true"
	if enableTLS {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")
		caFile := os.Getenv("TLS_CA_CERT_FILE")
		if certFile == "" || keyFile == "" || caFile == "" {
			return nil, fmt.Errorf("ENABLE_TLS is true but TLS_CERT_FILE or TLS_KEY_FILE or TLS_CA_CERT_FILE is not set %s", "")
		}
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("load server cert: %w", err)
		}
		// CA pool (for mTLS or client verification)
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read CA cert: %w", err)
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("append CA cert")
		}
		return credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.NoClientCert, //tls.RequestClientCert, ////tls.RequireAndVerifyClientCert, // Enable mTLS
			// ClientCAs:    certPool,
			MinVersion: tls.VersionTLS13, // Use TLS 1.3
		}), nil
	}
	return nil, nil
}

type connBoundRecordReader struct {
	array.RecordReader
	conn driver.Conn
}

func (r *connBoundRecordReader) Release() {
	r.RecordReader.Release()
	_ = r.conn.Close()
}
