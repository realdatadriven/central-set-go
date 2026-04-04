package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/realdatadriven/etlx"
	"github.com/robfig/cron/v3"
)

func (app *application) AdminGetJWT(user Dict) (string, error) {
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	claims.Subject = string(json_user)
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return "", err
	}
	return string(jwtBytes), nil
}

func (app *application) AdminInsertData(table string, data Dict) error {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return err
	} else {
		defer db.Close()
		var keys []any
		for key := range data {
			keys = append(keys, key)
		}
		cols := app.joinSlice(keys, `", "`)
		vals := app.joinSlice(keys, `, :`)
		sql := fmt.Sprintf(`insert into "%s" ("%s") values (:%s)`, table, cols, vals)
		_, err = db.ExecuteNamedQuery(sql, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) AdminExecuteQuery(query string, data Dict) error {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return err
	} else {
		defer db.Close()
		_, err = db.ExecuteNamedQuery(query, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) AdminGetRowByID(sql string, id any) (Dict, error) {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		res, _, err := db.QuerySingleRow(sql, []any{id}...)
		if err != nil {
			return nil, err
		}
		return *res, nil
	}
}

func (app *application) AdminGetRowByFilter(sql string, params []any) (Dict, error) {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		res, _, err := db.QuerySingleRow(sql, params...)
		if err != nil {
			return nil, err
		}
		return *res, nil
	}
}

func (app *application) AdminGetRowsByFilter(sql string, params []any) ([]Dict, error) {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		res, _, err := db.QueryMultiRows(sql, params...)
		if err != nil {
			return nil, err
		}
		return *res, nil
	}
}

func (app *application) GetRowByFilter(sql string, params Dict, filters []any) (Dict, error) {
	dsn, _, _ := app.GetDBNameFromParams(params)
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		res, _, err := db.QuerySingleRow(sql, filters...)
		if err != nil {
			return nil, err
		}
		return *res, nil
	}
}

func (app *application) GetRowsByFilter(sql string, params Dict, filters []any) ([]Dict, error) {
	dsn, _, _ := app.GetDBNameFromParams(params)
	//fmt.Println("GetRowsByFilter:", dsn)
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		res, _, err := db.QueryMultiRows(sql, filters...)
		if err != nil {
			return nil, err
		}
		return *res, nil
	}
}

func (app *application) ExecuteQuery(query string, params, data Dict) error {
	dsn, _, _ := app.GetDBNameFromParams(params)
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return err
	} else {
		defer db.Close()
		_, err = db.ExecuteNamedQuery(query, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) CronRunEndPoint(data Dict) (Dict, error) {
	api, ok := data["api"].(string)
	if !ok {
		api, _ = data["endpoint"].(string)
	}
	endpoint := fmt.Sprintf(`%s/%s`, app.config.baseURL, api)
	_jwt, ok := data["token"].(string)
	if !ok {
		_jwt, _ = app.AdminGetJWT(Dict{"user_id": 1, "username": "root", "role_id": 1, "active": true, "excluded": false})
	}
	req, _ := http.NewRequest("GET", endpoint, nil) // bytes.NewBuffer(jsonBody)
	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", _jwt))
	//req.Header.Set("Content-Type", "application/json")
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res_json Dict
	// Parse JSON into map
	err = json.NewDecoder(resp.Body).Decode(&res_json)
	if err != nil {
		return nil, err
	}
	//fmt.Println(1, res_json)
	return res_json, nil
}

func (app *application) CronJobsOLD() error {
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
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
	c := cron.New()
	for _, job := range *jobs {
		//fmt.Printf("1: %T, %v\n", job, job)
		_, err := c.AddFunc(job["cron"].(string), func() {
			//fmt.Printf("2: %T, %v\n", job, job)
			sql := `select * from "cron" where "cron_id" = ? and "cron" = ? and "active" = true and "excluded" = false`
			data, err := app.AdminGetRowByFilter(sql, []any{job["cron_id"], job["cron"]})
			last_run, last_run_ok := data["last_run"]
			if err != nil {
				data = job
				delete(data, "active")
				data["start_at"] = time.Now()
				data["end_at"] = time.Now()
				data["cron_msg"] = fmt.Sprintf("Error geting update version of %s->%s: %v", job["cron"], job["api"], err)
				data["success"] = false
				data["created_at"] = time.Now()
				data["updated_at"] = time.Now()
				data["excluded"] = false
				fmt.Printf("Error geting update version of %s: %v\n", job["api"], err)
				err = app.AdminInsertData("cron_log", data)
				if err != nil {
					fmt.Printf("Error saving the cron job log: %v\n", err)
				}
			} else if len(data) == 0 {
				data = job
				delete(data, "active")
				data["start_at"] = time.Now()
				data["end_at"] = time.Now()
				data["cron_msg"] = fmt.Sprintf("Error geting update version of %s->%s", job["cron"], job["api"])
				data["success"] = false
				data["created_at"] = time.Now()
				data["updated_at"] = time.Now()
				data["excluded"] = false
				fmt.Printf("Error geting update version of %s: %v\n", job["api"], err)
				err = app.AdminInsertData("cron_log", data)
				if err != nil {
					fmt.Printf("Error saving the cron job log: %v\n", err)
				}
			} else if run_only_once, run_only_once_ok := data["run_only_once"]; run_only_once_ok && app.toBool(run_only_once) && last_run_ok && last_run != nil {
				//
			} else {
				delete(data, "active")
				data["start_at"] = time.Now()
				endpoint := fmt.Sprintf(`%s/%s`, app.config.baseURL, data["api"].(string))
				fmt.Println("Running cron job:", data["cron_desc"], endpoint, data["start_at"])
				res_json, err := app.CronRunEndPoint(data)
				if err != nil {
					data["cron_msg"] = fmt.Sprintf("Error making %s request: %v", endpoint, err)
					data["success"] = false
					data["created_at"] = time.Now()
					data["updated_at"] = time.Now()
					data["excluded"] = false
					fmt.Printf("cron job %s finished %v", endpoint, data["end_at"])
					err = app.AdminInsertData("cron_log", data)
					if err != nil {
						fmt.Printf("Error saving the cron job log: %v\n", err)
					}
				} else {
					// update add last_run date
					data["updated_at"] = time.Now()
					data["last_run"] = time.Now()
					sql := `update cron set last_run = :last_run where cron_id = :cron_id`
					err = app.AdminExecuteQuery(sql, data)
					if err != nil {
						fmt.Printf("Body: %v -> %v\n", res_json, data)
						fmt.Printf("Error updating the cron last_run: %v\n", err)
					}
					data["created_at"] = time.Now()
					data["updated_at"] = time.Now()
					data["excluded"] = false
					fmt.Printf("cron job %s finished %v", endpoint, data["end_at"])
					err = app.AdminInsertData("cron_log", data)
					if err != nil {
						fmt.Printf("Body: %v -> %v\n", res_json, data)
						fmt.Printf("Error saving the cron job log: %v\n", err)
					}
				}
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

// RegisterCronJob registers (or re-registers) a single job on the global scheduler.
// If the job was previously registered, its old entry is removed first.
func (app *application) RegisterCronJob(job Dict) {
	cronID := job["cron_id"]
	app.cronEntriesMu.Lock()
	// Remove the old entry if it exists (handles schedule changes)
	if entryID, exists := app.cronEntries[cronID]; exists {
		app.cronScheduler.Remove(entryID)
		delete(app.cronEntries, cronID)
	}
	app.cronEntriesMu.Unlock()
	schedule, ok := job["cron"].(string)
	if !ok || schedule == "" {
		fmt.Printf("Skipping job %v: missing cron schedule\n", cronID)
		return
	}
	entryID, err := app.cronScheduler.AddFunc(schedule, func() {
		// Always re-fetch the job so we get the latest state from DB
		sql := `select * from "cron" where "cron_id" = ? and "cron" = ? and "active" = true and "excluded" = false`
		data, err := app.AdminGetRowByFilter(sql, []any{job["cron_id"], job["cron"]})
		logEntry := func(d Dict, msg string, success bool) {
			d["start_at"] = time.Now()
			d["end_at"] = time.Now()
			d["cron_msg"] = msg
			d["success"] = success
			d["created_at"] = time.Now()
			d["updated_at"] = time.Now()
			d["excluded"] = false
			delete(d, "active")
			if logErr := app.AdminInsertData("cron_log", d); logErr != nil {
				fmt.Printf("Error saving the cron job log: %v\n", logErr)
			}
		}
		if err != nil {
			logEntry(job, fmt.Sprintf("Error fetching updated job %v->%v: %v", job["cron"], job["api"], err), false)
			return
		}
		if len(data) == 0 {
			logEntry(job, fmt.Sprintf("Job no longer active: %v->%v", job["cron"], job["api"]), false)
			return
		}
		// run_only_once guard: skip if already ran
		if runOnce, ok := data["run_only_once"]; ok && app.toBool(runOnce) {
			if lastRun, ok := data["last_run"]; ok && lastRun != nil {
				fmt.Printf("Skipping run_only_once job %v, already ran at %v\n", data["api"], lastRun)
				return
			}
		}
		delete(data, "active")
		data["start_at"] = time.Now()
		endpoint := fmt.Sprintf(`%s/%s`, app.config.baseURL, data["api"].(string))
		fmt.Println("Running cron job:", data["cron_desc"], endpoint, data["start_at"])
		res_json, err := app.CronRunEndPoint(data)
		if err != nil {
			logEntry(data, fmt.Sprintf("Error calling %s: %v", endpoint, err), false)
			return
		}
		// Update last_run
		data["last_run"] = time.Now()
		data["updated_at"] = time.Now()
		updateSQL := `update cron set last_run = :last_run where cron_id = :cron_id`
		if updateErr := app.AdminExecuteQuery(updateSQL, data); updateErr != nil {
			fmt.Printf("Error updating last_run for %v: %v\n", data["api"], updateErr)
		}
		data["end_at"] = time.Now()
		data["created_at"] = time.Now()
		data["excluded"] = false
		fmt.Printf("cron job %s finished %v\n", endpoint, data["end_at"])
		_ = res_json
		logEntry(data, "OK", true)
	})
	if err != nil {
		fmt.Printf("Error adding cron %v (%v): %v\n", job["cron_desc"], schedule, err)
		logData := job
		delete(logData, "active")
		logData["start_at"] = time.Now()
		logData["end_at"] = time.Now()
		logData["cron_msg"] = fmt.Sprintf("Error adding the cron: %v", err)
		logData["success"] = false
		logData["created_at"] = time.Now()
		logData["updated_at"] = time.Now()
		logData["excluded"] = false
		if logErr := app.AdminInsertData("cron_log", logData); logErr != nil {
			fmt.Printf("Error saving the cron job log: %v\n", logErr)
		}
		return
	}
	app.cronEntriesMu.Lock()
	app.cronEntries[cronID] = entryID
	app.cronEntriesMu.Unlock()
	fmt.Printf("Registered cron job %v [%v] as entry %v\n", job["cron_desc"], schedule, entryID)
}

// CronJobs initialises the global scheduler, loads all existing active jobs,
// and adds a watcher that picks up new/updated rows every minute.
func (app *application) CronJobs() error {
	// Initialise shared state if not already done
	if app.cronEntries == nil {
		app.cronEntries = make(map[any]cron.EntryID)
	}
	if app.cronScheduler == nil {
		app.cronScheduler = cron.New()
	}
	// ── 1. Load all currently active jobs ───────────────────────────────────
	dsn, _, _ := app.GetDBNameFromParams(Dict{"db": app.config.db.dsn})
	db, err := etlx.GetDB(dsn)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	sql := `select * from "cron" where active = true and excluded = false`
	jobs, _, err := db.QueryMultiRows(sql, []any{}...)
	db.Close()
	if err != nil {
		return fmt.Errorf("error getting cron jobs: %w", err)
	}
	for _, job := range *jobs {
		app.RegisterCronJob(job)
	}
	// ── 2. Mark the baseline check time ─────────────────────────────────────
	app.lastCronCheck = time.Now()
	// ── 3. Watcher: every minute look for rows newer than lastCronCheck ──────
	_, err = app.cronScheduler.AddFunc("@every 1m", func() {
		checkFrom := app.lastCronCheck
		sql := `select * from "cron"
		        where active = true and excluded = false
		          and (created_at >= ? or updated_at >= ?)`
		newJobs, err := app.AdminGetRowsByFilter(sql, []any{checkFrom, checkFrom})
		app.lastCronCheck = time.Now() // advance the window before querying
		if err != nil {
			fmt.Printf("Watcher: error querying new cron jobs: %v\n", err)
			return
		}
		for _, job := range newJobs {
			fmt.Printf("Watcher: detected new/updated job %v, re-registering\n", job["cron_desc"])
			app.RegisterCronJob(job)
		}
	})
	if err != nil {
		return fmt.Errorf("error adding watcher cron: %w", err)
	}
	app.cronScheduler.Start()
	fmt.Println("Cron scheduler started")
	return nil
}
