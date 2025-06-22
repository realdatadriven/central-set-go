package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/realdatadriven/central-set-go/assets"
)

// Read SQL file and execute each query delimited by semicolon
func (app *application) setupDB(filename string, embedded bool) error {
	var content []byte
	var err error
	// Read the file content
	if embedded {
		content, err = assets.EmbeddedFiles.ReadFile(fmt.Sprintf(`setup/%s`, filename))
	} else {
		content, err = os.ReadFile(fmt.Sprintf(`database/%s`, filename))
	}
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	// Convert the content to a string and split by semicolon to get individual queries
	queries := strings.Split(string(content), ";")
	// Loop over each query, trimming and executing
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue // Skip empty queries
		}
		// Execute the query
		err := app.executeSQLQuery(trimmedQuery)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return nil
}

// Execute a single SQL query
func (app *application) executeSQLQuery(query string) error {
	//println(query)
	_, err := app.db.ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}
