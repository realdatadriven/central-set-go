package flight

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/duckdb/duckdb-go/v2"
	"google.golang.org/grpc"

	airport "github.com/hugr-lab/airport-go"
	"github.com/hugr-lab/airport-go/catalog"
	//duckarrow "github.com/duckdb/duckdb-go/v2/arrow"
)

// FlightManager is the interface the server uses to start/stop the FlightSQL server.
type FlightManager interface {
	Start(listenAddr string) error
	Stop(ctx context.Context) error
}

// AirportAdapter implements FlightManager using hugr-lab/airport-go.
type AirportAdapter struct {
	manager   *duckdb.Connector
	grpcSrv   *grpc.Server
	listener  net.Listener
	mem       memory.Allocator
	catalog   catalog.Catalog
	cfg       []map[string]any
	shutdownc chan struct{}
}

// NewAirportAdapter constructs the adapter with the provided DDB.
func NewAirportAdapter(manager *duckdb.Connector, config []map[string]any) *AirportAdapter {
	return &AirportAdapter{
		manager:   manager,
		mem:       memory.DefaultAllocator,
		cfg:       config,
		shutdownc: make(chan struct{}),
	}
}

// Start builds an airport-go catalog from the DuckDB schemas and tables discovered via the manager.
// It then creates a gRPC server and registers the airport server.
func (a *AirportAdapter) Start(listenAddr string) error {
	// Build catalog using airport.NewCatalogBuilder()
	builder := airport.NewCatalogBuilder()
	// For each schema defined in config, discover its tables and add them as SimpleTable entries.
	for _, s := range a.cfg {
		schemaName := s["flight_schema"].(string)
		// create a schema builder for this schema
		sb := builder.Schema(schemaName)
		//conn, _ := db.Connect(context.Background())
		conn := sql.OpenDB(a.manager)
		defer conn.Close()
		_, err := conn.ExecContext(context.Background(), fmt.Sprintf("USE %s;", schemaName))
		if err != nil {
			return fmt.Errorf("use schema %s: %w", schemaName, err)
		}
		q := `SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name`
		q = `select table_name from duckdb_tables`
		rows, err := conn.QueryContext(context.Background(), q, schemaName)
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
			q = fmt.Sprintf(`SELECT * FROM "%s"."%s" LIMIT 0`, schemaName, tname)
			conn, err := a.manager.Connect(context.Background())
			if err != nil {
				return err
			}
			defer conn.Close()
			arrow, err := duckdb.NewArrowFromConn(conn)
			if err != nil {
				return err
			}
			rdr, err := arrow.QueryContext(context.Background(), q)
			if err != nil {
				return err
			}
			defer rdr.Release()
			arrowSchema := rdr.Schema()
			//fmt.Println(arrowSchema)
			scanFn := makeScanFunc(a.manager, a.mem, schemaName, tname, arrowSchema)
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
		conn.ExecContext(context.Background(), "USE memory;")
	}
	cat, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build catalog: %w", err)
	}
	a.catalog = cat
	// Create grpc server and register airport server
	a.grpcSrv = grpc.NewServer()
	if err := airport.NewServer(a.grpcSrv, airport.ServerConfig{
		Catalog: cat,
		Auth: airport.BearerAuth(func(token string) (string, error) {
			fmt.Println("token:", token)
			if token == "secret-api-key" {
				return "user1", nil
			}
			return "", airport.ErrUnauthorized
		}),
		Address: listenAddr,
	}); err != nil {
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

// Good: Stream batches as you read
func makeScanFunc(db *duckdb.Connector, mem memory.Allocator, schemaName, tableName string, aSchema *arrow.Schema) func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	return func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
		query := fmt.Sprintf("SELECT * FROM \"%s\".\"%s\"", schemaName, tableName)
		fmt.Println(query)
		conn, err := db.Connect(context.Background())
		if err != nil {
			return nil, err
		}
		//defer conn.Close()
		arrow, err := duckdb.NewArrowFromConn(conn)
		if err != nil {
			return nil, err
		}
		rdr, err := arrow.QueryContext(context.Background(), query)
		if err != nil {
			conn.Close()
			return nil, err
		}
		/*if aSchema != nil && !arrow.SchemaEqual(aSchema, rdr.Schema()) {
			rdr.Release()
			conn.Close()
			return nil, fmt.Errorf("arrow schema mismatch")
		}*/
		/*/defer rdr.Release()
		for rdr.Next() {
			rec := rdr.Record()
			// rec is an Arrow RecordBatch
			fmt.Println("rows:", rec.NumRows())
		}*/
		//return nil, nil
		return &connBoundRecordReader{
			RecordReader: rdr,
			conn:         conn,
		}, nil
	}
}

type connBoundRecordReader struct {
	array.RecordReader
	conn driver.Conn
}

func (r *connBoundRecordReader) Release() {
	r.RecordReader.Release()
	_ = r.conn.Close()
}
