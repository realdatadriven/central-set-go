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
	// Read the file content
	if embedded && err != nil {
		content, err = assets.EmbeddedFiles.ReadFile(fmt.Sprintf(`setup/%s`, filename))
	}
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	// Convert the content to a string and split by semicolon to get individual queries
	queries := strings.Split(string(content), ";")
	// Loop over each query, trimming and executing
	dsn, _, _ := app.GetDBNameFromParams(map[string]any{"db": dbname})
	newDB, err := etlx.GetDB(dsn)
	if err != nil {
		return fmt.Errorf("geting the connection to %s: %w", dbname, err)
	}
	defer newDB.Close()
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
	return nil
}

// Execute a single SQL query
func (app *application) executeSQLQuery(query string, db etlx.DBInterface) error {
	_, err := db.ExecuteQuery(query)
	if err != nil {
		println(query)
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}
