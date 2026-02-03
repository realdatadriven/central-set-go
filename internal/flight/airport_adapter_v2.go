package flight

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/duckdb/duckdb-go/v2"
	"github.com/realdatadriven/etlx"
	"google.golang.org/grpc"

	airport "github.com/hugr-lab/airport-go"
	"github.com/hugr-lab/airport-go/catalog"
	"github.com/hugr-lab/airport-go/filter"
)

// =============================================================================
// Custom Catalog, Schema, and Table Implementations
// =============================================================================

// AFCatalog implements catalog.DynamicCatalog for both DDL and DML operations.
// It integrates with DuckDB for actual schema/table discovery and data scanning,
// but DDL/DML operations are currently stubbed.
type AFCatalog struct {
	mu      sync.RWMutex
	schemas map[string]*AFSchema
	config  []map[string]any // Configuration for DuckDB connections and queries
	mem     memory.Allocator
}

func NewAFCatalog(config []map[string]any) *AFCatalog {
	return &AFCatalog{
		schemas: make(map[string]*AFSchema),
		config:  config,
		mem:     memory.DefaultAllocator,
	}
}

// GetOrCreateSchema returns an existing schema or creates a new one.
// This is primarily for internal management; external DDL is stubbed.
func (c *AFCatalog) GetOrCreateSchema(name string) *AFSchema {
	c.mu.Lock()
	defer c.mu.Unlock()

	if s, ok := c.schemas[name]; ok {
		return s
	}

	s := NewAFSchema(name, "", c.config, c.mem)
	//s.Tables()
	c.schemas[name] = s
	return s
}

// Schemas implements catalog.Catalog.
func (c *AFCatalog) Schemas(ctx context.Context) ([]catalog.Schema, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// In a real scenario, this would dynamically discover schemas from DuckDB
	// For this example, we'll populate based on the provided config.
	var result []catalog.Schema
	for _, cfg := range c.config {
		schemaName := cfg["flight_schema"].(string)
		// Ensure schema exists in our internal map
		schema := c.GetOrCreateSchema(schemaName)
		//schema.Tables()
		result = append(result, schema)
	}
	return result, nil
}

// Schema implements catalog.Catalog.
func (c *AFCatalog) Schema(ctx context.Context, name string) (catalog.Schema, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if s, ok := c.schemas[name]; ok {
		return s, nil
	}
	// Attempt to discover from DuckDB if not already in memory
	for _, cfg := range c.config {
		if cfg["flight_schema"].(string) == name {
			schema := c.GetOrCreateSchema(name)
			return schema, nil
		}
	}
	return nil, nil
}

// CreateSchema implements catalog.DynamicCatalog.
func (c *AFCatalog) CreateSchema(_ context.Context, name string, opts catalog.CreateSchemaOptions) (catalog.Schema, error) {
	log.Printf("[DDL Not Supported] Attempted to CREATE SCHEMA: %s (Comment: %s)", name, opts.Comment)
	return nil, catalog.ErrNotFound
}

// DropSchema implements catalog.DynamicCatalog.
func (c *AFCatalog) DropSchema(_ context.Context, name string, opts catalog.DropSchemaOptions) error {
	log.Printf("[DDL Not Supported] Attempted to DROP SCHEMA: %s (IgnoreNotFound: %t)", name, opts.IgnoreNotFound)
	return catalog.ErrNotFound
}

// AFSchema implements catalog.DynamicSchema for table management operations.
type AFSchema struct {
	mu      sync.RWMutex
	name    string
	comment string
	tables  map[string]*AFTable
	config  []map[string]any // Configuration for DuckDB connections and queries
	mem     memory.Allocator
}

func NewAFSchema(name, comment string, config []map[string]any, mem memory.Allocator) *AFSchema {
	return &AFSchema{
		name:    name,
		comment: comment,
		tables:  make(map[string]*AFTable),
		config:  config,
		mem:     mem,
	}
}

func (s *AFSchema) Name() string    { return s.name }
func (s *AFSchema) Comment() string { return s.comment }

// Tables implements catalog.Schema.
func (s *AFSchema) Tables(ctx context.Context) ([]catalog.Table, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []catalog.Table
	// Iterate through config to find tables for this schema
	for _, cfg := range s.config {
		if cfg["flight_schema"].(string) == s.name {
			db, err := duckdb.NewConnector("", nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create DuckDB connector: %w", err)
			}
			conn := sql.OpenDB(db)
			defer conn.Close()

			_etlx := etlx.ETLX{}
			// Execute startup_sql and main_sql if present
			if startup_sql, ok := cfg["startup_sql"].(string); ok {
				startup_sql = _etlx.ReplaceEnvVariable(startup_sql)
				_, err := conn.ExecContext(context.Background(), startup_sql)
				if err != nil {
					log.Printf("Error executing startup_sql for schema %s: %v", s.name, err)
				}
			}
			if main_sql, ok := cfg["main_sql"].(string); ok {
				main_sql = _etlx.ReplaceEnvVariable(main_sql)
				_, err := conn.ExecContext(context.Background(), main_sql)
				if err != nil {
					log.Printf("Error executing main_sql for schema %s: %v", s.name, err)
				}
			}

			table_discover_sql := `select table_name from duckdb_tables`
			if _table_discover_sql, ok := cfg["table_discover_sql"].(string); ok {
				table_discover_sql = _etlx.ReplaceEnvVariable(_table_discover_sql)
			}
			rows, err := conn.QueryContext(context.Background(), table_discover_sql, s.name)
			if err != nil {
				return nil, fmt.Errorf("query tables for schema %s: %w", s.name, err)
			}
			defer rows.Close()

			for rows.Next() {
				var tname string
				if err := rows.Scan(&tname); err != nil {
					return nil, fmt.Errorf("scan table name for schema %s: %w", s.name, err)
				}
				// Check if table is explicitly allowed in config
				if tablesConfig, ok := cfg["tables"].(map[string]any); ok && len(tablesConfig) > 0 {
					if _, found := tablesConfig[tname]; !found {
						continue // Skip if not explicitly listed
					}
				}

				// Get Arrow Schema for the table
				duckdbConn, err := db.Connect(context.Background())
				if err != nil {
					return nil, fmt.Errorf("failed to connect to DuckDB for schema discovery: %w", err)
				}
				arrowReader, err := duckdb.NewArrowFromConn(duckdbConn)
				if err != nil {
					duckdbConn.Close()
					return nil, fmt.Errorf("failed to create Arrow reader from DuckDB connection: %w", err)
				}

				_fields := []string{"*"}
				if tablesConfig, ok := cfg["tables"].(map[string]any); ok {
					if tableConf, ok := tablesConfig[tname].(map[string]any); ok {
						if fields, ok := tableConf["fields"].(map[string]any); ok && len(fields) > 0 {
							_fields = make([]string, 0, len(fields))
							for fieldName := range fields {
								_fields = append(_fields, fmt.Sprintf(`"%s"`, fieldName))
							}
						}
					}
				}
				query := fmt.Sprintf(`SELECT %s FROM %s."%s" LIMIT 0`, strings.Join(_fields, ","), s.name, tname)
				if table_scan_tmpl_sql, ok := cfg["table_scan_tmpl_sql"].(string); ok {
					table_scan_tmpl_sql = _etlx.ReplaceEnvVariable(table_scan_tmpl_sql)
					query = strings.ReplaceAll(table_scan_tmpl_sql, "{{table_name}}", tname)
					query = strings.ReplaceAll(query, "{{schema_name}}", s.name)
					query = strings.ReplaceAll(query, "{{fields}}", strings.Join(_fields, ","))
					query = fmt.Sprintf(query, strings.Join(_fields, ","), s.name, tname)
				}

				rdr, err := arrowReader.QueryContext(context.Background(), query)
				if err != nil {
					duckdbConn.Close()
					return nil, fmt.Errorf("failed to query schema for table %s.%s: %w", s.name, tname, err)
				}
				arrowSchema := rdr.Schema()
				rdr.Release()
				duckdbConn.Close()

				// Create and store AFTable
				table := NewAFTable(s.name, tname, arrowSchema, cfg, s.mem)
				s.tables[tname] = table
				result = append(result, table)
			}
		}
	}
	return result, nil
}

// Table implements catalog.Schema.
func (s *AFSchema) Table(ctx context.Context, name string) (catalog.Table, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if t, ok := s.tables[name]; ok {
		return t, nil
	}
	// Attempt to discover the table if not already loaded
	_, err := s.Tables(ctx) // This will populate s.tables if not already done
	if err != nil {
		return nil, err
	}
	return s.tables[name], nil
}

// ScalarFunctions, TableFunctions, TableFunctionsInOut are not supported in this example.
func (s *AFSchema) ScalarFunctions(ctx context.Context) ([]catalog.ScalarFunction, error) {
	return nil, nil
}
func (s *AFSchema) TableFunctions(ctx context.Context) ([]catalog.TableFunction, error) {
	return nil, nil
}
func (s *AFSchema) TableFunctionsInOut(ctx context.Context) ([]catalog.TableFunctionInOut, error) {
	return nil, nil
}

// CreateTable implements catalog.DynamicSchema.
func (s *AFSchema) CreateTable(_ context.Context, name string, schema *arrow.Schema, opts catalog.CreateTableOptions) (catalog.Table, error) {
	log.Printf("[DDL Not Supported] Attempted to CREATE TABLE %s.%s (Comment: %s, OnConflict: %v)", s.name, name, opts.Comment, opts.OnConflict)
	return nil, catalog.ErrNotFound
}

// DropTable implements catalog.DynamicSchema.
func (s *AFSchema) DropTable(_ context.Context, name string, opts catalog.DropTableOptions) error {
	log.Printf("[DDL Not Supported] Attempted to DROP TABLE %s.%s (IgnoreNotFound: %t)", s.name, name, opts.IgnoreNotFound)
	return catalog.ErrNotFound
}

// RenameTable implements catalog.DynamicSchema.
func (s *AFSchema) RenameTable(_ context.Context, oldName, newName string, opts catalog.RenameTableOptions) error {
	log.Printf("[DDL Not Supported] Attempted to RENAME TABLE %s.%s to %s (IgnoreNotFound: %t)", s.name, oldName, newName, opts.IgnoreNotFound)
	return catalog.ErrNotFound
}

// AFTable implements catalog.DynamicTable, catalog.InsertableTable, catalog.UpdatableTable, catalog.DeletableTable.
type AFTable struct {
	mu          sync.RWMutex
	schemaName  string
	name        string
	arrowSchema *arrow.Schema
	config      map[string]any // Specific config for this table's schema
	mem         memory.Allocator
}

func NewAFTable(schemaName, name string, arrowSchema *arrow.Schema, config map[string]any, mem memory.Allocator) *AFTable {
	return &AFTable{
		schemaName:  schemaName,
		name:        name,
		arrowSchema: arrowSchema,
		config:      config,
		mem:         mem,
	}
}

func (t *AFTable) Name() string { return t.name }
func (t *AFTable) Comment() string {
	return fmt.Sprintf("Table %s.%s from DuckDB", t.schemaName, t.name)
}
func (t *AFTable) ArrowSchema(columns []string) *arrow.Schema {
	return catalog.ProjectSchema(t.arrowSchema, columns)
}

// Scan implements catalog.Table.
func (t *AFTable) Scan(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	log.Printf("[Scan Request] Scanning table %s.%s with options: %+v", t.schemaName, t.name, opts)

	_etlx := etlx.ETLX{}

	db, err := duckdb.NewConnector("", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DuckDB connector for scan: %w", err)
	}
	conn := sql.OpenDB(db)
	// conn will be closed by connBoundRecordReader.Release()

	// Execute startup_sql and main_sql if present
	if startup_sql, ok := t.config["startup_sql"].(string); ok {
		startup_sql = _etlx.ReplaceEnvVariable(startup_sql)
		_, err := conn.ExecContext(context.Background(), startup_sql)
		if err != nil {
			log.Printf("Error executing startup_sql for scan %s.%s: %v", t.schemaName, t.name, err)
		}
	}
	if main_sql, ok := t.config["main_sql"].(string); ok {
		main_sql = _etlx.ReplaceEnvVariable(main_sql)
		_, err := conn.ExecContext(context.Background(), main_sql)
		if err != nil {
			log.Printf("Error executing main_sql for scan %s.%s: %v", t.schemaName, t.name, err)
		}
	}

	// Prepare fields for query
	_fields := []string{"*"}
	// This part needs to be more robust, considering actual field access logic from the user's provided code
	// For simplicity, we'll assume all fields are accessible for now.

	query := fmt.Sprintf("SELECT %s FROM %s.\"%s\"", strings.Join(_fields, ","), t.schemaName, t.name)
	// Apply table_scan_tmpl_sql if available
	if table_scan_tmpl_sql, ok := t.config["table_scan_tmpl_sql"].(string); ok {
		table_scan_tmpl_sql = _etlx.ReplaceEnvVariable(table_scan_tmpl_sql)
		query = strings.ReplaceAll(table_scan_tmpl_sql, "{{table_name}}", t.name)
		query = strings.ReplaceAll(query, "{{schema_name}}", t.schemaName)
		query = strings.ReplaceAll(query, "{{fields}}", strings.Join(_fields, ","))
		query = fmt.Sprintf(query, strings.Join(_fields, ","), t.schemaName, t.name)
	}

	//hasFilters := false
	if opts.Filter != nil {
		fp, err := filter.Parse(opts.Filter)
		if err != nil {
			return nil, fmt.Errorf("failed to parse filter: %w", err)
		}
		enc := filter.NewDuckDBEncoder(nil)
		whereClause := enc.EncodeFilters(fp)
		if whereClause != "" {
			//hasFilters = true
			query = fmt.Sprintf("%s WHERE (%s)", query, whereClause)
		}
	}

	// Add scope filtering if needed (simplified for this example)
	// This would involve integrating the user's `rla_access` logic.

	duckdbConn, err := db.Connect(context.Background())
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to connect to DuckDB for query: %w", err)
	}
	arrowReader, err := duckdb.NewArrowFromConn(duckdbConn)
	if err != nil {
		_ = duckdbConn.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create Arrow reader for query: %w", err)
	}

	rdr, err := arrowReader.QueryContext(ctx, query)
	if err != nil {
		_ = duckdbConn.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to execute query for scan %s.%s: %w", t.schemaName, t.name, err)
	}

	return &connBoundRecordReader{
		RecordReader: rdr,
		conn:         duckdbConn,
	}, nil
}

// AddField implements catalog.DynamicTable.
func (t *AFTable) AddField(_ context.Context, _ *arrow.Schema, _ catalog.AddFieldOptions) error {
	log.Printf("[DDL Not Supported] Attempted to ADD FIELD to table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// RemoveField implements catalog.DynamicTable.
func (t *AFTable) RemoveField(_ context.Context, _ []string, _ catalog.RemoveFieldOptions) error {
	log.Printf("[DDL Not Supported] Attempted to REMOVE FIELD from table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// RenameField implements catalog.DynamicTable.
func (t *AFTable) RenameField(_ context.Context, _ []string, _ string, _ catalog.RenameFieldOptions) error {
	log.Printf("[DDL Not Supported] Attempted to RENAME FIELD in table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// ChangeColumnType implements catalog.DynamicTable.
func (t *AFTable) ChangeColumnType(_ context.Context, _ *arrow.Schema, _ string, _ catalog.ChangeColumnTypeOptions) error {
	log.Printf("[DDL Not Supported] Attempted to CHANGE COLUMN TYPE in table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// SetNotNull implements catalog.DynamicTable.
func (t *AFTable) SetNotNull(_ context.Context, _ string, _ catalog.SetNotNullOptions) error {
	log.Printf("[DDL Not Supported] Attempted to SET NOT NULL on column in table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// DropNotNull implements catalog.DynamicTable.
func (t *AFTable) DropNotNull(_ context.Context, _ string, _ catalog.DropNotNullOptions) error {
	log.Printf("[DDL Not Supported] Attempted to DROP NOT NULL on column in table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// SetDefault implements catalog.DynamicTable.
func (t *AFTable) SetDefault(_ context.Context, _ string, _ string, _ catalog.SetDefaultOptions) error {
	log.Printf("[DDL Not Supported] Attempted to SET DEFAULT on column in table %s.%s", t.schemaName, t.name)
	return catalog.ErrNotFound
}

// Insert implements catalog.InsertableTable.
func (t *AFTable) Insert(ctx context.Context, rows array.RecordReader, opts *catalog.DMLOptions) (*catalog.DMLResult, error) {
	log.Printf("[DML Not Supported] Attempted to INSERT into table %s.%s", t.schemaName, t.name)
	return nil, catalog.ErrNotFound
}

// Update implements catalog.UpdatableTable.
func (t *AFTable) Update(ctx context.Context, rowIDs []int64, rows array.RecordReader, opts *catalog.DMLOptions) (*catalog.DMLResult, error) {
	log.Printf("[DML Not Supported] Attempted to UPDATE table %s.%s", t.schemaName, t.name)
	return nil, catalog.ErrNotFound
}

// Delete implements catalog.DeletableTable.
func (t *AFTable) Delete(ctx context.Context, rowIDs []int64, opts *catalog.DMLOptions) (*catalog.DMLResult, error) {
	log.Printf("[DML Not Supported] Attempted to DELETE from table %s.%s", t.schemaName, t.name)
	return nil, catalog.ErrNotFound
}

// =============================================================================
// Main Server Setup (adapted from user's provided code)
// =============================================================================

type AFAirportAdapter struct {
	validateToken func(token string) (string, error)
	table_access  func(params map[string]any, tables []any) map[string]any
	rla_access    func(params map[string]any, tables []any, row_id []any) map[string]any
	grpcSrv       *grpc.Server
	listener      net.Listener
	mem           memory.Allocator
	catalog       catalog.Catalog
	cfg           []map[string]any
	shutdownc     chan struct{}
}

func NewAFAirportAdapter(
	config []map[string]any,
	validateToken func(token string) (string, error),
	table_access func(params map[string]any, tables []any) map[string]any,
	rla_access func(params map[string]any, tables []any, row_id []any) map[string]any,
) *AFAirportAdapter {
	return &AFAirportAdapter{
		validateToken: validateToken,
		table_access:  table_access,
		rla_access:    rla_access,
		mem:           memory.DefaultAllocator,
		cfg:           config,
		shutdownc:     make(chan struct{}),
	}
}

func (a *AFAirportAdapter) Start(listenAddr string) error {
	// Use our custom AFCatalog
	Catalog := NewAFCatalog(a.cfg)
	a.catalog = Catalog

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
		Catalog:  Catalog,
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

func (a *AFAirportAdapter) Stop(ctx context.Context) error {
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
