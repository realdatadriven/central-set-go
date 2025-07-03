package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/realdatadriven/etlx"
	"github.com/robfig/cron/v3"
)

func (app *application) AdminInsertData(table string, data map[string]any) error {
	var keys []any
	for key := range data {
		keys = append(keys, key)
	}
	cols := app.joinSlice(keys, `", "`)
	vals := app.joinSlice(keys, `, :`)
	sql := fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES (:%s)`, table, cols, vals)
	dsn, _, _ := app.GetDBNameFromParams(map[string]any{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return err
	} else {
		defer db.Close()
		_, err = db.ExecuteNamedQuery(sql, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) CronJobs() error {
	dsn, _, _ := app.GetDBNameFromParams(map[string]any{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return fmt.Errorf("error geting the db connection: %w", err)
	}
	defer db.Close()
	sql := `select * from "cron" where active = true and excluded = false`
	jobs, _, err := db.QueryMultiRows(sql, []any{}...)
	if err != nil {
		return fmt.Errorf("error geting the cron jobs: %w", err)
	}
	// Start cron
	c := cron.New()
	for _, job := range *jobs {
		fmt.Printf("1: %T, %v\n", job, job)
		_, err := c.AddFunc(job["cron"].(string), func() {
			fmt.Printf("2: %T, %v\n", job, job)
			data := job
			delete(data, "active")
			data["start_at"] = time.Now()
			endpoint := fmt.Sprintf(`%s/dyn_api/%s`, app.config.baseURL, job["api"].(string))
			fmt.Println("Running cron job:", data["cron_desc"], endpoint, data["start_at"])
			resp, err := http.Get(endpoint)
			if err != nil {
				// log.Printf("Job %s failed: %v\n", j.Name, err)
				return
			}
			defer resp.Body.Close()
			data["end_at"] = time.Now()
			var res_json map[string]any
			// Parse JSON into map
			err = json.NewDecoder(resp.Body).Decode(&res_json)
			if err != nil {
				fmt.Printf("%v", resp.Body)
				data["cron_msg"] = fmt.Sprintf("Error decoding %s response (%v): %v", endpoint, resp.Status, err)
				data["success"] = false
			} else {
				data["cron_msg"] = res_json["msg"]
				data["success"] = res_json["success"]
			}
			data["created_at"] = time.Now()
			data["updated_at"] = time.Now()
			data["excluded"] = false
			fmt.Printf("cron job %s finished %v", endpoint, data["start_at"])
			err = app.AdminInsertData("cron_log", data)
			if err != nil {
				fmt.Printf("Error saving the cron job log: %v\n", err)
			}
		})
		if err != nil {
			fmt.Printf("Error adding the cron %s: %v\n", job["cron_desc"], err)
			data := job
			delete(data, "active")
			data["start_at"] = time.Now()
			data["end_at"] = time.Now()
			data["cron_msg"] = fmt.Sprintf("Error adding the cron: %v", err)
			data["success"] = false
			data["created_at"] = time.Now()
			data["updated_at"] = time.Now()
			data["excluded"] = false
			err = app.AdminInsertData("cron_log", data)
			if err != nil {
				fmt.Printf("Error saving the cron job log: %v\n", err)
			}
		}
	}
	c.Start()
	return nil
}
