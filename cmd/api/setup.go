package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/realdatadriven/central-set-go/assets"
	"github.com/realdatadriven/etlx"
)

// Read SQL file and execute each query delimited by semicolon
func (app *application) setupDB(filename string, dbname string, embedded bool) error {
	var content []byte
	var err error
	fmt.Printf(`database/%s`, filename)
	content, err = os.ReadFile(fmt.Sprintf(`database/%s`, filename))
	if embedded && err != nil {
		content, err = assets.EmbeddedFiles.ReadFile(fmt.Sprintf(`setup/%s`, filename))
	}
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	queries := strings.Split(string(content), ";")
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": dbname})
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return fmt.Errorf("geting the connection to %s: %w", dbname, err)
	}
	defer newDB.Close()
	/*_, err = newDB.ExecuteQuery(string(content))
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}*/
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue // Skip empty queries
		}
		// Execute the query
		err := app.executeSQLQuery(trimmedQuery, newDB)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	csapp := fmt.Sprintf(`database/%s.%s.csapp`, dbname, app.config.db.driverName)
	if app.fileExists(csapp) {
		ddb, _ := etlx.GetDB("duckdb:")
		defer ddb.Close()
		sql := fmt.Sprintf(`select * from read_parquet('%s')`, csapp)
		fmt.Println(sql)
		res, _, err := ddb.QueryMultiRows(sql)
		if err != nil {
			return fmt.Errorf("failed to load data file %s: %w", csapp, err)
		}
		for _, d := range *res {
			fmt.Println(d["query"].(string))
			_, err := ddb.ExecuteQuery(d["query"].(string))
			if err != nil {
				return fmt.Errorf("failed execute data loading query %s: %w", d["query"], err)
			}
		}
	}
	return nil
}

// Execute a single SQL query
func (app *application) executeSQLQuery(query string, db etlx.DBInterface) error {
	if strings.HasPrefix(query, "PRAGMA") {
		println(query)
	}
	_, err := db.ExecuteQuery(query)
	if err != nil {
		println(query)
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}
