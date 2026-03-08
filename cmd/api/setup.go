package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/realdatadriven/central-set-go/assets"
	"github.com/realdatadriven/etlx"
)

func (app *application) setupWithModel(model string) error {
	var content []byte
	content, err := os.ReadFile(fmt.Sprintf(`%s`, model))
	if err != nil {
		content, err = os.ReadFile(fmt.Sprintf(`database/%s`, model))
		if err != nil {
			content, err = assets.EmbeddedFiles.ReadFile(fmt.Sprintf(`setup/%s`, model))
			if err != nil {
				return err
			}
		}
	}
	// Process the model content as needed
	params := Dict{
		"db": app.config.db.dsn,
		"data": Dict{
			"order_metadata": any(true),
			"config":         string(content),
		},
	}
	res := app.etlxRun(params, true)
	if res["success"].(bool) != true {
		return fmt.Errorf("failed to setup with model: %s", res["msg"])
	}
	return nil
}

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
	/*/ PARQUET STYLE
	if app.fileExists(csapp) {
		ddb, _ := etlx.GetDB("duckdb:")
		defer ddb.Close()
		sql := fmt.Sprintf(`select * from read_parquet('%s')`, csapp)
		// fmt.Println(sql)
		res, _, err := ddb.QueryMultiRows(sql)
		if err != nil {
			return fmt.Errorf("failed to load data file %s: %w", csapp, err)
		}
		for _, d := range *res {
			// fmt.Println(d["query"].(string))
			_, err := ddb.ExecuteQuery(d["query"].(string))
			if err != nil {
				return fmt.Errorf("failed execute data loading query %s: %w", d["query"], err)
			}
		}
	}*/
	// DUCKDB STYLE
	//fmt.Println("Data File:", app.fileExists(csapp), csapp)
	if app.fileExists(csapp) {
		//fmt.Printf(`duckdb:%s`, csapp)
		ddb, err := etlx.GetDB(fmt.Sprintf(`duckdb:%s`, csapp))
		if err != nil {
			return err
		}
		defer ddb.Close()
		// ADMIN
		sql := `select * from "adm_query"`
		res, _, err := ddb.QueryMultiRows(sql)
		if err != nil {
			return fmt.Errorf("failed to load data file %s: %w", csapp, err)
		}
		dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
		//fmt.Println(dsn)
		admDB, err := etlx.GetDB(dsn)
		if err != nil {
			return fmt.Errorf("geting the connection to %s: %w", app.config.db.dsn, err)
		}
		if admDB.GetDriverName() == "sqlite3" {
			admDB.ExecuteQuery("PRAGMA foreign_keys = OFF")
		}
		defer admDB.Close()
		for _, d := range *res {
			_, err := admDB.ExecuteQuery(d["query"].(string))
			if err != nil {
				fmt.Printf("failed execute data loading query %s: %s", d["query"], err)
				//return fmt.Errorf("failed execute data loading query %s: %w", d["query"], err)
			}
		}
		if admDB.GetDriverName() == "sqlite3" {
			admDB.ExecuteQuery("PRAGMA foreign_keys = ON")
		}
		// APP
		sql = `select * from "app_query"`
		res, _, err = ddb.QueryMultiRows(sql)
		if err != nil {
			return fmt.Errorf("failed to load data file %s: %w", csapp, err)
		}
		if newDB.GetDriverName() == "sqlite3" {
			newDB.ExecuteQuery("PRAGMA foreign_keys = OFF")
		}
		for _, d := range *res {
			_, err := newDB.ExecuteQuery(d["query"].(string))
			if err != nil {
				fmt.Printf("failed execute data loading query %s: %w", d["query"], err)
				// return fmt.Errorf("failed execute data loading query %s: %w", d["query"], err)
			}
		}
		if newDB.GetDriverName() == "sqlite3" {
			newDB.ExecuteQuery("PRAGMA foreign_keys = ON")
		}
	}
	return nil
}

// Execute a single SQL query
func (app *application) executeSQLQuery(query string, db etlx.DBInterface) error {
	//fmt.Println("Executing query...", query)
	_, err := db.ExecuteQuery(query)
	if err != nil {
		//fmt.Println(query)
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}
