package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/duckdb/duckdb-go/v2"
)

func main() {
	c, err := duckdb.NewConnector("", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	conn := sql.OpenDB(c)
	defer conn.Close()
	_, err = conn.ExecContext(context.Background(), "LOAD SQLITE;")
	if err != nil {
		fmt.Println(err)
	}
	_, err = conn.ExecContext(context.Background(), "ATTACH 'database/ADMIN.db' AS my_schema (TYPE SQLITE);")
	if err != nil {
		fmt.Println(err)
	}
	conn2, err := c.Connect(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	defer conn2.Close()

	// Obtain the Arrow from the connection.
	arrow, err := duckdb.NewArrowFromConn(conn2)
	if err != nil {
		fmt.Println(err)
	}
	q := `select table_name from duckdb_tables`
	rows, err := conn.QueryContext(context.Background(), q)
	if err != nil {
		fmt.Errorf("query tables: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tname string
		if err := rows.Scan(&tname); err != nil {
			fmt.Errorf("scan table name: %w", err)
		}
		rdr, err := arrow.QueryContext(context.Background(), "SELECT * FROM my_schema."+tname+" LIMIT 0")
		if err != nil {
			fmt.Println(err)
		}
		defer rdr.Release()
		fmt.Println(tname, rdr.Schema())
	}
	if err := rows.Err(); err != nil {
		fmt.Errorf("tables rows err: %w", err)
	}

	_, err = conn.ExecContext(context.Background(), "DETACH my_schema;")
	if err != nil {
		fmt.Println(err)
	}
	// Print the Arrow record batches.
	//go run -tags="duckdb_arrow" test-ddb-arrow.go
}
