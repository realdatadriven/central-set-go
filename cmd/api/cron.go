package main

import (
	"fmt"

	"github.com/realdatadriven/etlx"
	"github.com/robfig/cron/v3"
)

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
		_, err := c.AddFunc(job["cron"].(string), func() {
			/*resp, err := http.Get("http://localhost:8080" + j.Endpoint)
			if err != nil {
				log.Printf("Job %s failed: %v\n", j.Name, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Job %s ran: %s\n", job["cron_desc"].(string), resp.Status)*/
			fmt.Println("Running cron jobe", job)
		})
		if err != nil {
			//return fmt.Errorf("error geting the cron jobs: %w", err)
			fmt.Printf("failed to schedule job %s: %v\n", job["cron_desc"].(string), err)
		}
	}
	c.Start()
	return nil
}
