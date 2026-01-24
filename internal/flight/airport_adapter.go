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
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"log/slog"
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
	validateToken func(token string) (string, error)
	grpcSrv       *grpc.Server
	listener      net.Listener
	mem           memory.Allocator
	catalog       catalog.Catalog
	cfg           []map[string]any
	shutdownc     chan struct{}
}

// NewAirportAdapter constructs the adapter with the provided DDB.
func NewAirportAdapter(config []map[string]any, validateToken func(token string) (string, error)) *AirportAdapter {
	return &AirportAdapter{
		validateToken: validateToken,
		mem:           memory.DefaultAllocator,
		cfg:           config,
		shutdownc:     make(chan struct{}),
	}
}

// Start builds an airport-go catalog from the DuckDB schemas and tables discovered via the manager.
// It then creates a gRPC server and registers the airport server.
func (a *AirportAdapter) Start(listenAddr string) error {
	// Build catalog using airport.NewCatalogBuilder()
	builder := airport.NewCatalogBuilder()
	db, err := duckdb.NewConnector("", nil)
	if err != nil {
		return err
	}
	defer db.Close()
	conn := sql.OpenDB(db)
	defer conn.Close()
	// For each schema defined in config, discover its tables and add them as SimpleTable entries.
	for _, s := range a.cfg {
		schemaName := s["flight_schema"].(string)
		// create a schema builder for this schema
		sb := builder.Schema(schemaName)
		// execute startup_sql
		if _, ok := s["startup_sql"].(string); ok {
			//fmt.Printf("%s: %s\n", s["arrow_flight"], s["startup_sql"])
			_, err := conn.ExecContext(context.Background(), s["startup_sql"].(string))
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], s["startup_sql"], err)
			}
			//fmt.Printf("%s: %s\n", s["arrow_flight"], s["main_sql"])
			_, err = conn.ExecContext(context.Background(), s["main_sql"].(string))
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], s["main_sql"], err)
			}
		}
		_, err = conn.ExecContext(context.Background(), fmt.Sprintf("USE %s;", schemaName))
		if err != nil {
			return fmt.Errorf("use schema %s: %w", schemaName, err)
		}
		//conn, _ := db.Connect(context.Background())
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
			conn, err := db.Connect(context.Background())
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
			fmt.Println(tname, arrowSchema)
			scanFn := makeScanFunc( /*a.manager,*/ a.mem, schemaName, tname, arrowSchema, s)
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
		if _, ok := s["shutdown_sql"].(string); ok {
			// fmt.Printf("%s: %s\n", s["arrow_flight"], s["shutdown_sql"])
			_, err = conn.ExecContext(context.Background(), s["shutdown_sql"].(string))
			if err != nil {
				fmt.Printf("%s: %s: %s\n", s["arrow_flight"], s["shutdown_sql"], err)
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

// The idea is the use the global duckdb.Connector to create connections
// for each scan request, so we can share the same DB across multiple
// connections/flight clients.
// but im im passing the all configuration and making the connection and initializing the sturup and shutdown
// all over agin because of some problems with the shared connector across goroutines.
// This needs to be improved later for performance and resource usage.
func makeScanFunc( /*db *duckdb.Connector, */ mem memory.Allocator, schemaName, tableName string, aSchema *arrow.Schema, conf map[string]any) func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	return func(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
		query := fmt.Sprintf("SELECT * FROM \"%s\".\"%s\"", schemaName, tableName)
		fmt.Println(query)
		//fmt.Println(_sql, fligths)
		db, err := duckdb.NewConnector("", nil)
		if err != nil {
			return nil, err
		}
		defer func() { // EXECUTE SHUTDOWN SQL ON EXIT
			conn := sql.OpenDB(db)
			defer conn.Close()
			if _, ok := conf["shutdown_sql"].(string); ok {
				fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["shutdown_sql"])
				_, err := conn.ExecContext(context.Background(), conf["shutdown_sql"].(string))
				if err != nil {
					fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], conf["shutdown_sql"], err)
				}
			}
			db.Close()
		}()
		// EXECUTE STARTUP SQL
		conn := sql.OpenDB(db)
		defer conn.Close()
		if _, ok := conf["startup_sql"].(string); ok {
			//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["startup_sql"])
			_, err := conn.ExecContext(context.Background(), conf["startup_sql"].(string))
			if err != nil {
				fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], conf["startup_sql"], err)
			}
			//fmt.Printf("%s: %s\n", conf["arrow_flight"], conf["main_sql"])
			_, err = conn.ExecContext(context.Background(), conf["main_sql"].(string))
			if err != nil {
				fmt.Printf("Err %s: %s: %s\n", conf["arrow_flight"], conf["main_sql"], err)
			}
		}
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
			conn:         conn2,
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
