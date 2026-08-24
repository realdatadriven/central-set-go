package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type Connection struct {
	Valid  bool
	Driver string
	DSN    string
	Name   string
	Attach string
}

var supportedDrivers = map[string]bool{
	"duckdb":   true,
	"sqlite":   true,
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"odbc":     true,
	"mssql":    true,
}

func ParseETLXConn(s string) Connection {
	s = strings.TrimSpace(s)

	driver, dsn, ok := strings.Cut(s, ":")
	if !ok {
		return Connection{}
	}

	driver = strings.ToLower(strings.TrimSpace(driver))
	dsn = strings.TrimSpace(dsn)

	if driver == "" || dsn == "" || !supportedDrivers[driver] {
		return Connection{}
	}

	// Normalize aliases.
	if driver == "sqlite3" {
		driver = "sqlite"
	}

	if !validDSN(driver, dsn) {
		return Connection{}
	}

	name := connectionName(driver, dsn)
	if name == "" {
		return Connection{}
	}

	attach, err := duckDBAttach(driver, dsn, name)
	if err != nil {
		return Connection{}
	}

	return Connection{
		Valid:  true,
		Driver: driver,
		DSN:    dsn,
		Name:   name,
		Attach: attach,
	}
}

func validDSN(driver, dsn string) bool {
	switch driver {
	case "sqlite", "duckdb":
		return dsn != ""

	case "postgres":
		return strings.HasPrefix(dsn, "postgres://") ||
			strings.HasPrefix(dsn, "postgresql://") ||
			strings.Contains(dsn, "=")

	case "mysql":
		return strings.Contains(dsn, "@") ||
			strings.Contains(dsn, "tcp(")

	case "odbc":
		return strings.Contains(dsn, "=")

	case "mssql":
		return strings.HasPrefix(dsn, "sqlserver://") ||
			strings.Contains(dsn, "=")

	default:
		return false
	}
}

func connectionName(driver, dsn string) string {
	switch driver {
	case "sqlite", "duckdb":
		// database/db.sqlite -> db
		// database/mydb.duckdb -> mydb
		// /var/data/foo.db -> foo
		base := filepath.Base(dsn)

		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)

	default:
		// For network databases, use the database name
		// where it can be extracted.
		return databaseName(dsn)
	}
}

func databaseName(dsn string) string {
	// PostgreSQL / MySQL / MSSQL URLs
	if strings.Contains(dsn, "://") {
		if u, err := url.Parse(dsn); err == nil {
			name := strings.Trim(u.Path, "/")
			if name != "" {
				return name
			}
		}
	}

	// Look for Database=xxx / database=xxx / DBName=xxx
	for _, part := range strings.FieldsFunc(dsn, func(r rune) bool {
		return r == ';' || r == '&'
	}) {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(k)) {
		case "database", "dbname", "db", "initial catalog":
			return strings.TrimSpace(v)
		}
	}

	return ""
}

func duckDBAttach(driver, dsn, name string) (string, error) {
	var typ string

	switch driver {
	case "sqlite":
		typ = "SQLITE"
	case "duckdb":
		typ = "DUCKDB"
	case "postgres":
		typ = "POSTGRES"
	case "mysql":
		typ = "MYSQL"
	case "odbc":
		typ = "ODBC"
	case "mssql":
		typ = "MSSQL"
	default:
		return "", fmt.Errorf("unsupported driver %q", driver)
	}

	return fmt.Sprintf(
		"ATTACH IF NOT EXISTS '%s' AS %s (TYPE %s)",
		escapeDuckDBString(dsn),
		escapeDuckDBIdentifier(name),
		typ,
	), nil
}

func escapeDuckDBString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func escapeDuckDBIdentifier(s string) string {
	// DuckDB identifiers can be double quoted.
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
