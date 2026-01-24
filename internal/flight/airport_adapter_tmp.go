package flight

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/grpc"

	airport "github.com/hugr-lab/airport-go"
	"github.com/hugr-lab/airport-go/catalog"
	//duckarrow "github.com/duckdb/duckdb-go/v2/arrow"
)

// FlightManagerTmp is the interface the server uses to start/stop the FlightSQL server.
type FlightManagerTmp interface {
	Start(listenAddr string) error
	Stop(ctx context.Context) error
}

// AirportAdapterTmp implements FlightManagerTmp using hugr-lab/airport-go.
type AirportAdapterTmp struct {
	manager       *sql.DB
	validateToken func(token string) (string, error)
	grpcSrv       *grpc.Server
	listener      net.Listener
	mem           memory.Allocator
	catalog       catalog.Catalog
	cfg           []map[string]any
	shutdownc     chan struct{}
}

// NewAirportAdapterTmp constructs the adapter with the provided DDB.
func NewAirportAdapterTmp(manager *sql.DB, config []map[string]any, validateToken func(token string) (string, error)) *AirportAdapterTmp {
	return &AirportAdapterTmp{
		manager:       manager,
		validateToken: validateToken,
		mem:           memory.DefaultAllocator,
		cfg:           config,
		shutdownc:     make(chan struct{}),
	}
}

// Start builds an airport-go catalog from the DuckDB schemas and tables discovered via the manager.
// It then creates a gRPC server and registers the airport server.
func (a *AirportAdapterTmp) Start(listenAddr string) error {
	// Build catalog using airport.NewCatalogBuilder()
	builder := airport.NewCatalogBuilder()

	// For each schema defined in config, discover its tables and add them as SimpleTable entries.
	for _, s := range a.cfg {
		schemaName := s["flight_schema"].(string)
		// create a schema builder for this schema
		sb := builder.Schema(schemaName)

		// discover tables for this schema using information_schema
		tables, err := discoverTables(a.manager, schemaName)
		if err != nil {
			log.Printf("[flight] warning: failed to discover tables for schema %s: %v", schemaName, err)
			continue
		}

		for _, t := range tables {
			arrowSchema := buildArrowSchemaFromColumns(t.Columns)
			// fmt.Println(schemaName, t.Name)
			// create scan func closure capturing schema/table/columns
			scanFn := makeScanFuncTmp(a.manager, a.mem, schemaName, t.Name, arrowSchema, t.Columns)

			// register simple table under current schema builder
			sb.SimpleTable(airport.SimpleTableDef{
				Name:     t.Name,
				Comment:  t.Name,
				Schema:   arrowSchema,
				ScanFunc: scanFn,
			})
		}
	}
	cat, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build catalog: %w", err)
	}
	a.catalog = cat
	// Create grpc server and register airport server
	debugLevel := slog.LevelDebug
	config := airport.ServerConfig{
		Catalog:  cat,
		Auth:     airport.BearerAuth(a.validateToken),
		Address:  listenAddr,
		LogLevel: &debugLevel,
	}
	opts := airport.ServerOptions(config)
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
func (a *AirportAdapterTmp) Stop(ctx context.Context) error {
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

// Internal helper types
type tableMeta struct {
	Schema  string
	Name    string
	Columns []columnMeta
}

type columnMeta struct {
	Name       string
	DuckDBType string
}

// discoverTables uses DuckDB's information_schema to list tables and columns for a given schema.
func discoverTables(db *sql.DB, schema string) ([]tableMeta, error) {
	fmt.Println("discoverTables for schema:", schema)
	q := `SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name`
	q = `select table_name from duckdb_tables`
	rows, err := db.QueryContext(context.Background(), q, schema)
	if err != nil {
		return nil, fmt.Errorf("query tables: %w", err)
	}
	defer rows.Close()
	var out []tableMeta
	for rows.Next() {
		var tname string
		if err := rows.Scan(&tname); err != nil {
			return nil, fmt.Errorf("scan table name: %w", err)
		}
		// fmt.Println("TABLE:", tname)
		cols, err := discoverColumns(db, schema, tname)
		if err != nil {
			return nil, fmt.Errorf("discover columns for %s.%s: %w", schema, tname, err)
		}
		out = append(out, tableMeta{
			Schema:  schema,
			Name:    tname,
			Columns: cols,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tables rows err: %w", err)
	}
	return out, nil
}

func discoverColumns(db *sql.DB, schema, table string) ([]columnMeta, error) {
	q := `SELECT column_name, data_type FROM duckdb_columns WHERE /*schema_name = ? AND*/ table_name = ? ORDER BY numeric_precision`
	//rows, err := db.Query(q, schema, table)d
	rows, err := db.QueryContext(context.Background(), q, table)
	if err != nil {
		return nil, fmt.Errorf("query columns: %s %w", schema, err)
	}
	defer rows.Close()
	var cols []columnMeta
	for rows.Next() {
		var name, dtype string
		if err := rows.Scan(&name, &dtype); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		// fmt.Println("  COLUMN:", name, dtype)
		cols = append(cols, columnMeta{
			Name:       name,
			DuckDBType: strings.ToUpper(strings.TrimSpace(dtype)),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("columns rows err: %w", err)
	}
	return cols, nil
}

// buildArrowSchemaFromColumns converts discovered columns into an arrow.Schema following provided mapping.
func buildArrowSchemaFromColumns(cols []columnMeta) *arrow.Schema {
	fields := make([]arrow.Field, 0, len(cols))
	for _, c := range cols {
		dt := mapDuckTypeToArrow(c.DuckDBType)
		fields = append(fields, arrow.Field{Name: c.Name, Type: dt, Nullable: true})
	}
	return arrow.NewSchema(fields, nil)
}

// mapDuckTypeToArrow implements the mapping rules you provided.
func mapDuckTypeToArrow(duck string) arrow.DataType {
	// Accept already uppercase strings (discoverColumns uppercases)
	switch {
	case duck == "BOOLEAN":
		return arrow.FixedWidthTypes.Boolean
	case strings.HasPrefix(duck, "VARCHAR") || strings.HasPrefix(duck, "TEXT") || duck == "STRING":
		return arrow.BinaryTypes.String
	case strings.HasPrefix(duck, "INT8") || duck == "TINYINT":
		return arrow.PrimitiveTypes.Int8
	case strings.HasPrefix(duck, "INT16") || duck == "SMALLINT":
		return arrow.PrimitiveTypes.Int16
	case strings.HasPrefix(duck, "INT32") || duck == "INTEGER" || duck == "INT":
		return arrow.PrimitiveTypes.Int32
	case strings.HasPrefix(duck, "INT64") || duck == "BIGINT":
		return arrow.PrimitiveTypes.Int64
	case duck == "UTINYINT" || duck == "UINTEGER8":
		return arrow.PrimitiveTypes.Uint8
	case duck == "USMALLINT" || duck == "UINTEGER16":
		return arrow.PrimitiveTypes.Uint16
	case duck == "UINTEGER" || duck == "UINT32":
		return arrow.PrimitiveTypes.Uint32
	case duck == "UBIGINT" || duck == "UINT64":
		return arrow.PrimitiveTypes.Uint64
	case strings.HasPrefix(duck, "FLOAT") || duck == "REAL":
		return arrow.PrimitiveTypes.Float32
	case strings.HasPrefix(duck, "DOUBLE") || strings.HasPrefix(duck, "DOUBLE PRECISION"):
		return arrow.PrimitiveTypes.Float64
	case strings.HasPrefix(duck, "DECIMAL"):
		// parse DECIMAL(precision,scale)
		var p int32 = 38
		var s int32 = 6
		if i := strings.Index(duck, "("); i != -1 {
			// tolerate formats like DECIMAL(p,s)
			fmt.Sscanf(duck, "DECIMAL(%d,%d)", &p, &s)
		}
		// create a Decimal128Type
		if f, _ := arrow.NewDecimalType(arrow.DECIMAL256, 38, 6); f != nil {
			return f
		}
		return arrow.BinaryTypes.String
	case strings.HasPrefix(duck, "TIMESTAMP"):
		// use microsecond precision timestamp
		return arrow.FixedWidthTypes.Timestamp_us
	case strings.HasPrefix(duck, "DATE"):
		return arrow.FixedWidthTypes.Date32
	case strings.HasPrefix(duck, "TIME"):
		return arrow.FixedWidthTypes.Time64us
	case strings.HasPrefix(duck, "BLOB") || strings.HasPrefix(duck, "BYTEA"):
		return arrow.BinaryTypes.Binary
	case strings.HasPrefix(duck, "UUID"):
		return arrow.BinaryTypes.String
	case strings.HasPrefix(duck, "LIST"):
		// Default LIST<T> → LIST<STRING>
		return arrow.ListOf(arrow.BinaryTypes.String)
	case strings.HasPrefix(duck, "STRUCT"):
		// Default STRUCT<...> → empty struct
		return arrow.StructOf()
	default:
		// fallback → string
		return arrow.BinaryTypes.String
	}
}

// makeScanFunc returns a catalog.ScanFunc that queries DuckDB and converts results to Arrow RecordReader
func makeScanFuncTmp(db *sql.DB, mem memory.Allocator, schemaName, tableName string, aSchema *arrow.Schema, cols []columnMeta) func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	return func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
		// Build query: SELECT * FROM schema.table
		query := fmt.Sprintf("SELECT * FROM \"%s\".\"%s\"", schemaName, tableName)
		if opts != nil && opts.Limit > 0 {
			query = fmt.Sprintf("%s LIMIT %d", query, opts.Limit)
		}
		fmt.Println(query)
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query %s: %w", query, err)
		}
		defer rows.Close()

		// We'll build a single in-memory Record (for simplicity).
		builder := array.NewRecordBuilder(mem, aSchema)
		defer func() {
			// ensure builder is released if the function exits before we return a RecordReader
			//_ = builder.Release()
		}()

		colCount := len(cols)
		for rows.Next() {
			// Scan into []interface{}
			values := make([]interface{}, colCount)
			valuePtrs := make([]interface{}, colCount)
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				return nil, fmt.Errorf("scan row: %w", err)
			}

			// append each column into builder
			for i := 0; i < colCount; i++ {
				v := values[i]
				// handle nulls
				if v == nil {
					builder.Field(i).AppendNull()
					continue
				}

				switch aSchema.Field(i).Type.(type) {
				case *arrow.Decimal128Type:
					// Decimal parsing is complex; append null for now (TODO: implement robust parsing)
					builder.Field(i).AppendNull()
				case *arrow.Int8Type:
					switch vv := v.(type) {
					case int8:
						builder.Field(i).(*array.Int8Builder).Append(vv)
					case int16:
						builder.Field(i).(*array.Int8Builder).Append(int8(vv))
					case int32:
						builder.Field(i).(*array.Int8Builder).Append(int8(vv))
					case int64:
						builder.Field(i).(*array.Int8Builder).Append(int8(vv))
					case int:
						builder.Field(i).(*array.Int8Builder).Append(int8(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Int16Type:
					switch vv := v.(type) {
					case int16:
						builder.Field(i).(*array.Int16Builder).Append(vv)
					case int32:
						builder.Field(i).(*array.Int16Builder).Append(int16(vv))
					case int64:
						builder.Field(i).(*array.Int16Builder).Append(int16(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Int32Type:
					switch vv := v.(type) {
					case int32:
						builder.Field(i).(*array.Int32Builder).Append(vv)
					case int64:
						builder.Field(i).(*array.Int32Builder).Append(int32(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Int64Type:
					switch vv := v.(type) {
					case int64:
						builder.Field(i).(*array.Int64Builder).Append(vv)
					case int32:
						builder.Field(i).(*array.Int64Builder).Append(int64(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Uint8Type:
					switch vv := v.(type) {
					case int64:
						builder.Field(i).(*array.Uint8Builder).Append(uint8(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Uint16Type:
					switch vv := v.(type) {
					case int64:
						builder.Field(i).(*array.Uint16Builder).Append(uint16(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Uint32Type:
					switch vv := v.(type) {
					case int64:
						builder.Field(i).(*array.Uint32Builder).Append(uint32(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Uint64Type:
					switch vv := v.(type) {
					case int64:
						builder.Field(i).(*array.Uint64Builder).Append(uint64(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Float32Type:
					switch vv := v.(type) {
					case float32:
						builder.Field(i).(*array.Float32Builder).Append(vv)
					case float64:
						builder.Field(i).(*array.Float32Builder).Append(float32(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Float64Type:
					switch vv := v.(type) {
					case float64:
						builder.Field(i).(*array.Float64Builder).Append(vv)
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.StringType:
					switch vv := v.(type) {
					case []byte:
						builder.Field(i).(*array.StringBuilder).Append(string(vv))
					case string:
						builder.Field(i).(*array.StringBuilder).Append(vv)
					default:
						builder.Field(i).(*array.StringBuilder).Append(fmt.Sprint(vv))
					}
				case *arrow.BinaryType:
					switch vv := v.(type) {
					case []byte:
						builder.Field(i).(*array.BinaryBuilder).Append(vv)
					case string:
						builder.Field(i).(*array.BinaryBuilder).Append([]byte(vv))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.TimestampType:
					switch vv := v.(type) {
					case time.Time:
						us := vv.UnixNano() / 1_000
						builder.Field(i).(*array.TimestampBuilder).Append(arrow.Timestamp(us))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Date32Type:
					switch vv := v.(type) {
					case time.Time:
						// Date32 is days since epoch (UTC)
						days := int32(vv.UTC().Truncate(24*time.Hour).Unix() / 86400)
						builder.Field(i).(*array.Date32Builder).Append(arrow.Date32(days))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.Time64Type:
					switch vv := v.(type) {
					case time.Time:
						us := vv.UnixNano() / 1_000
						builder.Field(i).(*array.Time64Builder).Append(arrow.Time64(us))
					default:
						builder.Field(i).AppendNull()
					}
				case *arrow.ListType, *arrow.StructType:
					// Complex types are not implemented: append null
					builder.Field(i).AppendNull()
				default:
					// fallback: append string representation
					switch vv := v.(type) {
					case []byte:
						builder.Field(i).(*array.StringBuilder).Append(string(vv))
					case string:
						builder.Field(i).(*array.StringBuilder).Append(vv)
					default:
						builder.Field(i).(*array.StringBuilder).Append(fmt.Sprint(vv))
					}
				}
			}
		}

		// check rows.Err()
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		// Build a single record
		rec := builder.NewRecordBatch()
		// release builder now that record is created
		builder.Release()

		// Return a RecordReader with a single record. The airport-go code expects the caller to call reader.Release().
		//return array.NewRecordReader(rec.Schema(), []arrow.RecordBatch{rec}), nil
		return array.NewRecordReader(rec.Schema(), []arrow.RecordBatch{rec})
	}
}
