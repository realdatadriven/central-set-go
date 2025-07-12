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

func (app *application) CronJobs() error {
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
