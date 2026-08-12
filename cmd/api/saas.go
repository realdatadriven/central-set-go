package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	texttemplate "text/template"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	cetypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/realdatadriven/central-set-go/assets"
	"github.com/realdatadriven/central-set-go/internal/env"
)

func GenerateSecret() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type TerraformRun struct {
	Config string
	State  json.RawMessage // store as JSON in DB
	Lock   string
}

// func that checks if the link has not prefix http:// or https:// and adds http:// if missing
func ensureHTTPPrefix(url string) string {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		return "http://" + url
	}
	return url
}

// turn json.RawMessage like '"text"' into string text
func rawMessageToString(raw json.RawMessage) string {
	var str string
	if err := json.Unmarshal(raw, &str); err != nil {
		return ""
	}
	return strings.ReplaceAll(str, `"`, "")
}

func terraformCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve home dir: %w", err)
	}
	cacheDir := filepath.Join(home, ".cs-terraform-cache")
	// Create directory if not exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create terraform cache dir: %w", err)
	}
	// Add Terraform plugin cache
	//env = append(env, "TF_PLUGIN_CACHE_DIR="+cacheDir)
	return cacheDir, nil
}

func opentofuCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve home dir: %w", err)
	}
	cacheDir := filepath.Join(home, ".cs-opentofu-cache")
	// Create directory if not exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create opentofu cache dir: %w", err)
	}
	// Add OpenTofu plugin cache
	//env = append(env, "TF_PLUGIN_CACHE_DIR="+cacheDir)
	return cacheDir, nil
}

func (app *application) HandleService(params Dict, action string) Dict {
	var user_id int
	if _, ok := params["user"].(Dict)["user_id"]; ok {
		user_id = app.toInt(params["user"].(Dict)["user_id"])
	}
	var loc *time.Location
	if _, ok := params["location"].(*time.Location); ok {
		loc = params["location"].(*time.Location)
	} else {
		loc = time.Local
	}
	_data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		if _, ok := params["data"].(Dict)["data"]; !ok {
			msg, _ := app.i18n.T("no-data", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		_data = params["data"].(Dict)["data"].(Dict)
	} else {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	_table := params["data"].(Dict)["table"]
	sql := `select * from "subs_server" where "subscription_id" = ? and "active" = true and "excluded" = false`
	subsServer, err := app.GetRowsByFilter(sql, params, []any{_data["subscription_id"]})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	for _, server := range subsServer {
		sql = `select * from "subs_server_service" where subs_server_id = ? AND "subscription_id" = ? and "active" = true and "excluded" = false`
		subsServerService, err := app.GetRowsByFilter(sql, params, []any{server["subs_server_id"], _data["subscription_id"]})
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		var svc *ServiceManager
		if server["subs_server"].(string) == "localhost" || strings.HasPrefix(server["subs_server"].(string), "127.") {
			local := NewLocal()
			svc = NewServiceManager(local)
		} else {
			sshIntance, err := NewSSH(server["subs_server"].(string), server["subs_server_user"].(string), server["subs_server_key"].(string), server["subs_server_host_key"].(string))
			if err != nil {
				return Dict{
					"success": false,
					"msg":     err.Error(),
				}
			}
			defer sshIntance.Close()
			svc = NewServiceManager(sshIntance)
		}
		serverLogsData := ""
		logs_ins_sql := `INSERT INTO srv_logs_check_hist
			(srv_logs_check_hist, srv_logs, subs_server_service_id, subs_server_id, subscription_id, tenant_id, resquested_at, response_at, err_response_message, user_id, created_at, updated_at, excluded) VALUES
			(:srv_logs_check_hist, :srv_logs, :subs_server_service_id, :subs_server_id, :subscription_id, :tenant_id, :resquested_at, :response_at, :err_response_message, :user_id, :created_at, :updated_at, :excluded)
		`
		statusData := ""
		status_ins_sql := `INSERT INTO srv_status_chk_hist
			(srv_status_chk_hist, srv_status, subs_server_service_id, subs_server_id, subscription_id, tenant_id, resquested_at, response_at, err_response_message, user_id, created_at, updated_at, excluded) VALUES
			(:srv_status_chk_hist, :srv_status, :subs_server_service_id, :subs_server_id, :subscription_id, :tenant_id, :resquested_at, :response_at, :err_response_message, :user_id, :created_at, :updated_at, :excluded)
		`
		for _, service := range subsServerService {
			switch action {
			case "start":
				if _, ok := service["service_start_cmd"].(string); !ok || service["service_start_cmd"].(string) == "" {
					err = svc.Start(context.Background(), service["subs_server_service"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				} else {
					err = svc.Run(context.Background(), service["service_start_cmd"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				}
				_data["deployed"] = true
				_data["status"] = "started"
			case "stop":
				if _, ok := service["service_stop_cmd"].(string); !ok || service["service_stop_cmd"].(string) == "" {
					err = svc.Stop(context.Background(), service["subs_server_service"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				} else {
					err = svc.Run(context.Background(), service["service_stop_cmd"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				}
				_data["deployed"] = true
				_data["status"] = "stoped"
			case "restart":
				if _, ok := service["service_restart_cmd"].(string); !ok || service["service_restart_cmd"].(string) == "" {
					err = svc.Restart(context.Background(), service["subs_server_service"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				} else {
					err = svc.Run(context.Background(), service["service_restart_cmd"].(string))
					if err != nil {
						return Dict{
							"success": false,
							"msg":     err.Error(),
						}
					}
				}
				_data["deployed"] = true
				_data["status"] = "started"
			case "status":
				statusData = ""
				now := time.Now().In(loc)
				logs := service
				logs["user_id"] = user_id
				logs["resquested_at"] = now
				logs["srv_status_chk_hist"] = fmt.Sprintf("%s %s %s", service["subs_server_service"], action, now.Format("2006-01-02 15:04:05"))
				logs["err_response_message"] = nil
				if _, ok := service["service_status_cmd"].(string); !ok || service["service_status_cmd"].(string) == "" {
					statusData, err = svc.Status(context.Background(), service["subs_server_service"].(string))
					if err != nil {
						logs["err_response_message"] = err.Error()
					}
				} else {
					statusData, err = svc.RunOutput(context.Background(), service["service_status_cmd"].(string))
					if err != nil {
						logs["err_response_message"] = err.Error()
					}
				}
				logs["srv_status"] = statusData
				logs["response_at"] = time.Now().In(loc)
				logs["created_at"] = time.Now().In(loc)
				logs["updated_at"] = time.Now().In(loc)
				logs["excluded"] = false
				if logs["err_response_message"] != nil {
					return Dict{
						"success": false,
						"msg":     logs["err_response_message"],
					}
				}
				err := app.ExecuteQuery(status_ins_sql, params, logs)
				if err != nil {
					fmt.Printf("Err inserting status check history: %s\n", err.Error())
				}
			case "logs":
				serverLogsData = ""
				now := time.Now().In(loc)
				logs := service
				logs["user_id"] = user_id
				logs["resquested_at"] = now
				logs["srv_logs_check_hist"] = fmt.Sprintf("%s %s %s", service["subs_server_service"], action, now.Format("2006-01-02 15:04:05"))
				logs["err_response_message"] = nil
				if _, ok := service["service_logs_cmd"].(string); !ok || service["service_logs_cmd"].(string) == "" {
					lines := 100
					if _, ok := params["lines"]; ok {
						lines = app.toInt(params["lines"])
					}
					serverLogsData, err = svc.Logs(context.Background(), service["subs_server_service"].(string), lines)
					if err != nil {
						logs["err_response_message"] = err.Error()
					}
				} else {
					serverLogsData, err = svc.RunOutput(context.Background(), service["service_logs_cmd"].(string))
					if err != nil {
						logs["err_response_message"] = err.Error()
					}
				}
				logs["srv_logs"] = serverLogsData
				logs["response_at"] = time.Now().In(loc)
				logs["created_at"] = time.Now().In(loc)
				logs["updated_at"] = time.Now().In(loc)
				logs["excluded"] = false
				if logs["err_response_message"] != nil {
					return Dict{
						"success": false,
						"msg":     logs["err_response_message"],
					}
				}
				err := app.ExecuteQuery(logs_ins_sql, params, logs)
				if err != nil {
					fmt.Printf("Err inserting logs check history: %s\n", err.Error())
				}
			}
		}
		if statusData != "" {
			_data["server_status"] = statusData
		}
		if serverLogsData != "" {
			_data["server_logs"] = serverLogsData
		}
	}
	params["data"].(Dict)["table"] = _table
	params["data"].(Dict)["data"] = _data
	upsert := app.create_update(params)
	// fmt.Println(upsert)
	if _, ok := upsert["success"]; !ok {
		return upsert
	} else if _, ok := upsert["success"].(bool); !ok {
		return upsert
	} else if ok, _ := upsert["success"].(bool); !ok {
		return upsert
	}
	_, db, _ := app.GetDBNameFromParams(params)
	bdc := Dict{
		"type":     "data_change",
		"database": db,
		"table":    params["data"].(Dict)["table"],
	}
	app.BroadCastChange(bdc)
	//fmt.Println("BOADCASTED:", bdc)
	return Dict{
		"success": true,
		"msg":     "Action completed successfully",
	}
}

func DomainName(name string) string {
	name = strings.ToLower(name)
	// Decompose accented characters
	name = norm.NFD.String(name)
	var b strings.Builder
	for _, r := range name {
		// Remove accents
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case unicode.IsSpace(r) || r == '_' || r == '-':
			b.WriteRune('_')
		default:
			// remove special characters
		}
	}
	// Collapse multiple underscores
	result := b.String()
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}

func (app *application) GetSubdomain(plan, tenant, subscription Dict) string {
	subdomain := ""
	// SAAS_SUBDOMAIN_TMPL='{{.plan.plan}}{{.subscription.subscription_id}}.{{.tenant.slug}}'
	// SAAS_SUBDOMAIN_TMPL='{{.plan.plan}}{{.subscription.subscription_id}}.{{if .tenant.slug}}{{.tenant.slug}}{{else}}{{.tenant.tenant}}{{end}}'
	if os.Getenv("SAAS_SUBDOMAIN_TMPL") != "" {
		subtml, err := app.RenderTextTemplate(os.Getenv("SAAS_SUBDOMAIN_TMPL"), Dict{"plan": plan, "tenant": tenant, "subscription": subscription})
		if err != nil {
			fmt.Println("Error rendering SAAS_SUBDOMAIN_TMPL template", err.Error())
		} else {
			return DomainName(subtml)
		}
	}
	if plan_name, ok := plan["plan"].(string); ok {
		// fmt.Println(plan_name, subscription["subscription_id"])
		subdomain = fmt.Sprintf("%s%d", DomainName(plan_name), toInt(subscription["subscription_id"]))
	}
	if slug, ok := tenant["slug"].(string); ok {
		subdomain = fmt.Sprintf("%s%s", subdomain, DomainName(slug))
	} else {
		subdomain = fmt.Sprintf("%s%s", subdomain, DomainName(tenant["tenant"].(string)))
	}
	return subdomain
}

// type ReadViaOData struct {params Dict}
func (app *application) ODataGetRow(params Dict, table, field string, id any) (Dict, error) {
	_, odb, _ := app.GetDBNameFromParams(params)
	odata_path := fmt.Sprintf(`%s/%s?$filter=%s eq %d`, odb, table, field, app.toInt(id))
	// fmt.Println("ODATA Path", odata_path)
	data := app.ODataRead(params, odata_path)
	if !data["success"].(bool) {
		return nil, fmt.Errorf("%s", data["msg"])
	}
	if _, ok := data["data"].([]Dict); !ok {
		return nil, fmt.Errorf("No data returned for %s id %d", table, id)
	}
	if len(data["data"].([]Dict)) == 0 {
		return nil, fmt.Errorf("No data returned for %s id %d", table, id)
	}
	return data["data"].([]Dict)[0], nil
}

func (app *application) RunDeploy(params Dict) Dict {
	_data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		if _, ok := params["data"].(Dict)["data"]; !ok {
			msg, _ := app.i18n.T("no-data", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		_data = params["data"].(Dict)["data"].(Dict)
	} else {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	_table := params["data"].(Dict)["table"]
	action := params["data"].(Dict)["action"].(string)
	// sql := `select * from "subscription" where "subscription_id" = ? and "active" = true and "excluded" = false`
	subs, err := app.ODataGetRow(params, "subscription", "subscription_id", _data["subscription_id"]) //app.GetRowByFilter(sql, params, []any{_data["subscription_id"]})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	if len(subs) > 0 {
		_data = subs
	}
	// sql := `select * from "plan" where "plan_id" = ? and "active" = true and "excluded" = false`
	plan, err := app.ODataGetRow(params, "plan", "plan_id", _data["plan_id"]) //app.GetRowByFilter(sql, params, []any{_data["plan_id"]})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	// data["data"].([]Dict)
	sql := `select * from "plan_service" where "plan_id" = ? and "active" = true and "excluded" = false`
	plan_service, err := app.GetRowsByFilter(sql, params, []any{plan["plan_id"]})
	if err != nil {
		fmt.Println("env err:", err)
	}
	tenantID := _data["tenant_id"]
	subscriptionID := _data["subscription_id"]
	// sql = `select * from "tenant" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenant, err := app.ODataGetRow(params, "tenant", "tenant_id", _data["tenant_id"]) // app.GetRowByFilter(sql, params, []any{tenantID})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	sql = `select * from "env" where "tenant_id" = ? and "on_srv_start" = true and "active" = true and "excluded" = false`
	tenantEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	if err != nil {
		fmt.Println("env err:", err)
	}
	tenantEnvKeyPair := map[string]any{}
	for _, v := range tenantEnv {
		tenantEnvKeyPair[v["env_name"].(string)] = v["env_value"].(string)
	}
	sql = `select * from "sys_env" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenantSysEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	if err != nil {
		fmt.Println("sys_env err:", err)
	}
	tenantSysEnvKeyPair := map[string]any{}
	for _, v := range tenantSysEnv {
		tenantSysEnvKeyPair[v["env_name"].(string)] = v["env_value"].(string)
	}
	//fmt.Println(tenant, tenantEnv)
	if _, ok := plan["terraform_template"].(string); !ok {
		msg, _ := app.i18n.T("unable-to-match-plan-id", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	secret1, _ := GenerateSecret()
	secret2, _ := GenerateSecret()
	secret3, _ := GenerateSecret()
	_tmpl_data := map[string]any{
		"secret_key":      secret1,
		"secret_key2":     secret2,
		"secret_key3":     secret3,
		"tenant_id":       tenantID,
		"subscription_id": subscriptionID,
		"subdomain":       app.GetSubdomain(plan, tenant, _data),
		"env":             tenantEnv,
		"envKV":           tenantEnvKeyPair,
		"sysEnv":          tenantSysEnv,
		"sysEnvKV":        tenantSysEnvKeyPair,
		"plan":            plan,
		"tenant":          tenant,
		"data":            _data,
	}
	// plan_file
	sql = `select * from "plan_file" where "plan_id" = ? and "active" = true and "excluded" = false`
	// fmt.Println(sql, _data["plan_id"])
	plan_file, err := app.GetRowsByFilter(sql, params, []any{_data["plan_id"]})
	_files := map[string]string{}
	if err != nil {
		fmt.Println("Error getting the plan files", err.Error())
	} else {
		for _, _file := range plan_file {
			fname := _file["plan_file"].(string)
			// file name replace the last _ with . and remove the first part before the last
			name := strings.Replace(fname, "_", ".", -1)
			base := filepath.Base(name)
			ext := filepath.Ext(base)
			base_no_ext := strings.Replace(base, ext, "", 1)
			temptex, err := os.CreateTemp("", fmt.Sprintf("%s-*%s", base_no_ext, ext))
			if err != nil {
				fmt.Println("Error creating temporary file", fname, temptex.Name(), err.Error())
			}
			//fmt.Println("FNAME:", fname, temptex.Name())
			filetml, err := app.RenderTextTemplate(_file["file_template"].(string), _tmpl_data)
			if err != nil {
				fmt.Println("Error rendering file template", err.Error())
				filetml = _file["file_template"].(string)
			}
			// fmt.Println(temptex.Name(), output_path)
			// defer os.Remove(temptex.Name())
			defer temptex.Close()
			_, err = temptex.WriteString(filetml)
			if err != nil {
				fmt.Println("Error writing to temporary file", err.Error())
			}
			temptex.Close()
			_files[fname] = temptex.Name()
		}
		_tmpl_data["plan_files"] = _files
	}
	parsedTmpl, err := app.RenderTextTemplate(plan["terraform_template"].(string), _tmpl_data)
	if err != nil {
		fmt.Print("Terraform Template Parse Err:", err.Error())
		return Dict{
			"success": false,
			"msg":     "terraform_template Err:" + err.Error(),
		}
	}
	// fmt.Printf(parsedTmpl)
	run := &TerraformRun{Config: parsedTmpl}
	if terraform_state, ok := _data["terraform_state"].(string); ok && terraform_state != "" {
		// fmt.Println("State:", _data["terraform_state"])
		run.State = json.RawMessage([]byte(terraform_state))
	}
	if terraform_lock, ok := _data["terraform_lock"].(string); ok && terraform_lock != "" {
		// fmt.Println("Lock:", __data["terraform_state"])
		run.Lock = terraform_lock
	}
	err = EnsureDefaultCredentials(env.GetString("SSH_KEY", ""))
	if err != nil {
		fmt.Printf("Err: %s!\n", err.Error())
		// return Dict{"success": false, "msg": fmt.Sprintf("Err: %s!\n", err.Error())}
	}
	var res map[string]string
	switch action {
	case "deploy", "install":
		//res, err = app.DeployTerraformForTenant(params, tenantID, run)
		res, err = app.DeployOpenTofuForTenant(params, tenantID, run, action)
		_json_out, _ := json.Marshal(res)
		_data["tf_public_ip"] = rawMessageToString(json.RawMessage(res["public_ip"]))
		_data["tf_public_dns"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["public_dns"])))
		_data["tf_public_url"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["url"])))
		_data["terraform_outputs"] = string(_json_out)
		if err == nil {
			_data["deployed"] = true
			_data["status"] = "started"
		} else {
			_data["deployed"] = false
			_data["status"] = "has_error"
		}
	/*case "start", "boot", "starup", "restart":
		res, err = app.DeployOpenTofuForTenant(params, tenantID, run, "restart")
		if err == nil {
			_data["deployed"] = true
			_data["status"] = "started"
		}
	case "stop", "shutdown", "sp", "stp", "pause":
		res, err = app.DeployOpenTofuForTenant(params, tenantID, run, "stop")
		if err == nil {
			_data["deployed"] = true
			_data["status"] = "stoped"
		}*/
	case "cancel", "destroy":
		//err = app.DestroyTerraform(params, tenantID, run)
		err = app.DestroyOpenTofu(params, tenantID, run, action)
		if err != nil {
			_data["deployed"] = true
			_data["status"] = "has_error"
		} else {
			_data["deployed"] = false
			_data["status"] = "destroyed"
		}
	}
	msg, _ := app.i18n.T("success", Dict{})
	_data["tf_err_msg"] = nil
	success := true
	if err != nil {
		_data["tf_err_msg"] = fmt.Sprintf("%s: %s", action, err.Error())
		msg = err.Error()
		success = false
	}
	_data["terraform_state"] = string(run.State)
	_data["terraform_lock"] = string(run.Lock)
	err2 := AddKnownHost(_data["tf_public_ip"].(string), 22, env.GetString("SSH_HOST_KEY", "$HOME/.ssh/known_hosts"))
	if err2 != nil {
		fmt.Printf("Err adding to nkown host: %s!\n", err2.Error())
	}
	params["data"].(Dict)["table"] = _table
	//fmt.Println(1, "terraform_state", len(_data["terraform_state"].(string)), "terraform_lock", len(_data["terraform_lock"].(string)), len(run.Lock))
	params["data"].(Dict)["data"] = _data
	upsert := app.create_update(params)
	//fmt.Println(upsert)
	if _, ok := upsert["success"]; !ok {
		return upsert
	} else if _, ok := upsert["success"].(bool); !ok {
		return upsert
	} else if ok, _ := upsert["success"].(bool); !ok {
		return upsert
	}
	_, db, _ := app.GetDBNameFromParams(params)
	bdc := Dict{
		"type":     "data_change",
		"database": db,
		"table":    params["data"].(Dict)["table"],
	}
	app.BroadCastChange(bdc)
	if len(plan_service) > 0 {
		sql = `select * from "subs_server" where "subscription_id" = ? and "active" = true and "excluded" = false`
		subsServer, err := app.GetRowsByFilter(sql, params, []any{_data["subscription_id"]})
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		if len(subsServer) == 0 {
			params["data"].(Dict)["table"] = "subs_server"
			params["data"].(Dict)["data"] = Dict{
				"subs_server":          rawMessageToString(json.RawMessage(res["service_server"])),
				"subs_server_desc":     rawMessageToString(json.RawMessage(res["service_server"])),
				"subs_server_user":     rawMessageToString(json.RawMessage(res["service_user"])),
				"subs_server_key":      rawMessageToString(json.RawMessage(res["service_key"])),
				"subs_server_host_key": rawMessageToString(json.RawMessage(res["service_host_key"])),
				"subscription_id":      _data["subscription_id"],
				"tenant_id":            _data["tenant_id"],
				"active":               true,
			}
			results := app.create_update(params)
			//fmt.Println(upsert)
			if _, ok := upsert["success"]; !ok {
				// return upsert
			} else if _, ok := upsert["success"].(bool); !ok {
				//return upsert
			} else if ok, _ := upsert["success"].(bool); !ok {
				// return upsert
			} else {
				_id := results["inserted_primary_key"]
				params["data"].(Dict)["table"] = "subs_server_service"
				auxSrv := []Dict{}
				for _, _srv := range plan_service {
					srv, err := app.RenderTextTemplate(_srv["plan_service"].(string), _tmpl_data)
					if err != nil {
						// fmt.Print("Terraform Template Parse Err:", err.Error())
						return Dict{
							"success": false,
							"msg":     "plan_service template Err:" + err.Error(),
						}
					}
					auxSrv = append(auxSrv, Dict{
						"subs_server_service": srv,
						"subs_server_id":      _id,
						"subscription_id":     _data["subscription_id"],
						"tenant_id":           _data["tenant_id"],
						"active":              true,
					})
				}
				params["data"].(Dict)["data"] = auxSrv
				results = app.create_update(params)
			}
		}
	}
	// fmt.Println("BOADCASTED:", bdc)
	// plan_service
	//println(subsServer, plan_service)
	return Dict{
		"success": success,
		"msg":     msg,
		"data":    res,
	}
}

func (app *application) DeployTerraformForTenant(params Dict, tenantID any, run *TerraformRun, action string) (map[string]string, error) {
	workDir, _ := os.MkdirTemp("", "tf-*")
	fmt.Println("Temp TF Dir:", workDir)
	os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644)
	//return nil, fmt.Errorf("template_test: %s", workDir)
	if len(run.State) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, "terraform.tfstate"), run.State, 0644); err != nil {
			return nil, err
		}
	}
	if len(run.Lock) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, ".terraform.lock.hcl"), []byte(run.Lock), 0644); err != nil {
			return nil, err
		}
	}
	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, err
	}
	tf, _ := tfexec.NewTerraform(workDir, tfPath)
	// 👇 Load tenant-specific env vars
	//tenantEnv, _ := getTenantEnvVars(tenantID)
	sql := `select * from "env" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	// Merge with system env so PATH, etc. still exists
	baseEnv := os.Environ()
	mergedEnv := map[string]string{}
	for _, e := range baseEnv {
		//fmt.Println(e)
		_env := strings.Split(e, "=")
		mergedEnv[_env[0]] = _env[1]
	}
	for _, v := range tenantEnv {
		mergedEnv[v["env_name"].(string)] = v["env_value"].(string)
		//fmt.Printf("Setting env var for tenant %s: %s=%s\n", tenantID, k, v)
	}
	_path, err := terraformCachePath()
	//fmt.Println("Temp TF CASH Dir:", _path)
	if err != nil {
		fmt.Println("Temp TF CASH Dir Err:", _path, err)
	}
	mergedEnv["TF_PLUGIN_CACHE_DIR"] = _path
	tf.SetEnv(mergedEnv)
	var stateBytes []byte
	var lockBytes []byte
	// Normal Terraform lifecycle
	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		return nil, err
	}
	if len(run.State) > 0 {
		if err := tf.Refresh(context.Background()); err != nil {
			stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
			if _err == nil {
				run.State = json.RawMessage(stateBytes)
			}
			lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
			if _err == nil {
				run.Lock = string(lockBytes)
			}
			fmt.Printf("resfreh failed: %s", err)
			return nil, fmt.Errorf("resfreh failed: %w", err)
		}
	}
	//fmt.Println("Init Pass")
	opts := []tfexec.ApplyOption{
		tfexec.Var(fmt.Sprintf("action=%s", action)),
		//tfexec.Var("region=eu-west-1"),
		//tfexec.Var("instance_count=3"),
		//tfexec.Var(fmt.Sprintf("environment=%s", "staging")),
		//tfexec.VarFile("prod.tfvars"), // optional, if you also have a vars file
		//tfexec.Parallelism(5),
		//tfexec.Refresh(true),
	}
	if err := tf.Apply(context.Background(), opts...); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Apply Failed:", err)
		return nil, err
	}
	//fmt.Println("Apply Pass")
	outputs, err := tf.Output(context.Background())
	if err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Output Failed:", err)
		return nil, err
	}
	//fmt.Println("Output Pass")
	result := map[string]string{}
	for k, v := range outputs {
		//if v.Value != nil {
		//fmt.Println(k, string(v.Value))
		result[k] = fmt.Sprintf("%s", v.Value)
		//}
	}
	stateBytes, _ = os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	lockBytes, _ = os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
	run.Lock = string(lockBytes)
	tfexec.CleanEnv(mergedEnv)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)
	//fmt.Printf("Terraform outputs: %v\n", result)
	//fmt.Printf("Terraform state: %s\n", string(run.State))
	return result, nil
}

func (app *application) DeployOpenTofuForTenant(params Dict, tenantID any, run *TerraformRun, action string) (map[string]string, error) {
	workDir, _ := os.MkdirTemp("", "tofu-*")
	fmt.Println("Temp Tofu Dir:", workDir)
	os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644)
	//return nil, fmt.Errorf("template_test: %s", workDir)
	if len(run.State) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, "terraform.tfstate"), run.State, 0644); err != nil {
			return nil, err
		}
	}
	if len(run.Lock) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, ".terraform.lock.hcl"), []byte(run.Lock), 0644); err != nil {
			return nil, err
		}
	}
	tfPath, err := exec.LookPath("tofu")
	if err != nil {
		return nil, err
	}
	tf, _ := tfexec.NewTerraform(workDir, tfPath)
	// 👇 Load tenant-specific env vars
	//tenantEnv, _ := getTenantEnvVars(tenantID)
	sql := `select * from "env" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	// Merge with system env so PATH, etc. still exists
	baseEnv := os.Environ()
	mergedEnv := map[string]string{}
	for _, e := range baseEnv {
		//fmt.Println(e)
		_env := strings.Split(e, "=")
		mergedEnv[_env[0]] = _env[1]
	}
	for _, v := range tenantEnv {
		mergedEnv[v["env_name"].(string)] = v["env_value"].(string)
		//fmt.Printf("Setting env var for tenant %s: %s=%s\n", tenantID, k, v)
	}
	_path, err := opentofuCachePath()
	//fmt.Println("Temp Tofu CASH Dir:", _path)
	if err != nil {
		fmt.Println("Temp Tofu CASH Dir Err:", _path, err)
	}
	mergedEnv["TF_PLUGIN_CACHE_DIR"] = _path
	tf.SetEnv(mergedEnv)
	var stateBytes []byte
	var lockBytes []byte
	// Normal OpenTofu lifecycle
	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		return nil, err
	}
	if len(run.State) > 0 {
		if err := tf.Refresh(context.Background()); err != nil {
			stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
			if _err == nil {
				run.State = json.RawMessage(stateBytes)
			}
			lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
			if _err == nil {
				run.Lock = string(lockBytes)
			}
			fmt.Printf("refresh failed: %s", err)
			return nil, fmt.Errorf("refresh failed: %w", err)
		}
	}
	//fmt.Println("Init Pass")
	opts := []tfexec.ApplyOption{
		tfexec.Var(fmt.Sprintf("action=%s", action)),
		//tfexec.Var("region=eu-west-1"),
		//tfexec.Var("instance_count=3"),
		//tfexec.Var(fmt.Sprintf("environment=%s", "staging")),
		//tfexec.VarFile("prod.tfvars"), // optional, if you also have a vars file
		//tfexec.Parallelism(5),
		//tfexec.Refresh(true),
	}
	if err := tf.Apply(context.Background(), opts...); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Apply Failed:", err)
		return nil, err
	}
	//fmt.Println("Apply Pass")
	outputs, err := tf.Output(context.Background())
	if err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Output Failed:", err)
		return nil, err
	}
	//fmt.Println("Output Pass")
	result := map[string]string{}
	for k, v := range outputs {
		//if v.Value != nil {
		//fmt.Println(k, string(v.Value))
		result[k] = fmt.Sprintf("%s", v.Value)
		//}
	}
	stateBytes, _ = os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	lockBytes, _ = os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
	run.Lock = string(lockBytes)
	tfexec.CleanEnv(mergedEnv)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)
	//fmt.Printf("OpenTofu outputs: %v\n", result)
	//fmt.Printf("OpenTofu state: %s\n", string(run.State))
	return result, nil
}

func (app *application) RenderTemplate(tmplStr string, data map[string]any) (string, error) {
	// Create a FuncMap with some common functions
	// funcMap := sprig.FuncMap()
	tmpl, err := template.New("tmpl").Funcs(sprig.FuncMap()).Parse(tmplStr)
	//tmpl, err := template.New("email").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	//fmt.Println(buf.String())
	return buf.String(), nil
}

func (app *application) RenderTextTemplate(tmplStr string, data map[string]any) (string, error) {
	// Create a FuncMap with some common functions
	// funcMap := sprig.FuncMap()
	tmpl, err := texttemplate.New("tmpl").Funcs(sprig.FuncMap()).Parse(tmplStr)
	//tmpl, err := template.New("email").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	//fmt.Println(buf.String())
	return buf.String(), nil
}

func (app *application) DestroyTerraform(params Dict, tenantID any, run *TerraformRun, action string) error {
	// 1. Create temp dir
	workDir, _ := os.MkdirTemp("", "tf-*")
	fmt.Println("Temp TF Dir:", workDir)
	// 2. Write config
	if err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644); err != nil {
		return err
	}
	// 3. Restore state
	if err := os.WriteFile(filepath.Join(workDir, "terraform.tfstate"), run.State, 0644); err != nil {
		return err
	}
	if len(run.Lock) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, ".terraform.lock.hcl"), []byte(run.Lock), 0644); err != nil {
			return err
		}
	}
	// 4. Run terraform destroy
	tfPath, _ := exec.LookPath("terraform")
	tf, _ := tfexec.NewTerraform(workDir, tfPath)
	sql := `select * from "env" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	// Merge with system env so PATH, etc. still exists
	baseEnv := os.Environ()
	mergedEnv := map[string]string{}
	for _, e := range baseEnv {
		//fmt.Println(e)
		_env := strings.Split(e, "=")
		mergedEnv[_env[0]] = _env[1]
	}
	for _, v := range tenantEnv {
		mergedEnv[v["env_name"].(string)] = v["env_value"].(string)
		//fmt.Printf("Setting env var for tenant %s: %s=%s\n", tenantID, k, v)
	}
	_path, err := terraformCachePath()
	//fmt.Println("Temp TF CASH Dir:", _path)
	if err != nil {
		fmt.Println("Temp TF CASH Dir Err:", _path, err)
	}
	mergedEnv["TF_PLUGIN_CACHE_DIR"] = _path
	tf.SetEnv(mergedEnv)
	var stateBytes []byte
	var lockBytes []byte
	// Normal Terraform lifecycle
	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Init Failed:", err)
		return err
	}

	if err := tf.Refresh(context.Background()); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Printf("resfreh failed: %s", err)
		return fmt.Errorf("resfreh failed: %w", err)
	}
	opts := []tfexec.DestroyOption{
		tfexec.Var(fmt.Sprintf("action=%s", action)),
		//tfexec.Var("region=eu-west-1"),
		//tfexec.Var("instance_count=3"),
		//tfexec.Var(fmt.Sprintf("environment=%s", "staging")),
		//tfexec.VarFile("prod.tfvars"), // optional, if you also have a vars file
		//tfexec.Parallelism(5),
		//tfexec.Refresh(true),
	}
	if err := tf.Destroy(context.Background(), opts...); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Printf("destroy failed: %s", err)
		return fmt.Errorf("destroy failed: %w", err)
	}
	tfexec.CleanEnv(mergedEnv)
	stateBytes, _ = os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	lockBytes, _ = os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
	run.Lock = string(lockBytes)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)
	return nil
}

func (app *application) DestroyOpenTofu(params Dict, tenantID any, run *TerraformRun, action string) error {
	// 1. Create temp dir
	workDir, _ := os.MkdirTemp("", "tofu-*")
	fmt.Println("Temp Tofu Dir:", workDir)
	// 2. Write config
	if err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644); err != nil {
		return err
	}
	// 3. Restore state
	if err := os.WriteFile(filepath.Join(workDir, "terraform.tfstate"), run.State, 0644); err != nil {
		return err
	}
	if len(run.Lock) > 0 {
		if err := os.WriteFile(filepath.Join(workDir, ".terraform.lock.hcl"), []byte(run.Lock), 0644); err != nil {
			return err
		}
	}
	// 4. Run tofu destroy
	tfPath, _ := exec.LookPath("tofu")
	tf, _ := tfexec.NewTerraform(workDir, tfPath)
	sql := `select * from "env" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	// Merge with system env so PATH, etc. still exists
	baseEnv := os.Environ()
	mergedEnv := map[string]string{}
	for _, e := range baseEnv {
		//fmt.Println(e)
		_env := strings.Split(e, "=")
		mergedEnv[_env[0]] = _env[1]
	}
	for _, v := range tenantEnv {
		mergedEnv[v["env_name"].(string)] = v["env_value"].(string)
		//fmt.Printf("Setting env var for tenant %s: %s=%s\n", tenantID, k, v)
	}
	_path, err := opentofuCachePath()
	//fmt.Println("Temp Tofu CASH Dir:", _path)
	if err != nil {
		fmt.Println("Temp Tofu CASH Dir Err:", _path, err)
	}
	mergedEnv["TF_PLUGIN_CACHE_DIR"] = _path
	tf.SetEnv(mergedEnv)
	var stateBytes []byte
	var lockBytes []byte
	// Normal OpenTofu lifecycle
	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Println("Init Failed:", err)
		return err
	}

	if err := tf.Refresh(context.Background()); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Printf("refresh failed: %s", err)
		return fmt.Errorf("refresh failed: %w", err)
	}
	opts := []tfexec.DestroyOption{
		tfexec.Var(fmt.Sprintf("action=%s", action)),
		//tfexec.Var("region=eu-west-1"),
		//tfexec.Var("instance_count=3"),
		//tfexec.Var(fmt.Sprintf("environment=%s", "staging")),
		//tfexec.VarFile("prod.tfvars"), // optional, if you also have a vars file
		//tfexec.Parallelism(5),
		//tfexec.Refresh(true),
	}
	if err := tf.Destroy(context.Background(), opts...); err != nil {
		stateBytes, _err := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
		if _err == nil {
			run.State = json.RawMessage(stateBytes)
		}
		lockBytes, _err := os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
		if _err == nil {
			run.Lock = string(lockBytes)
		}
		fmt.Printf("destroy failed: %s", err)
		return fmt.Errorf("destroy failed: %w", err)
	}
	tfexec.CleanEnv(mergedEnv)
	stateBytes, _ = os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	lockBytes, _ = os.ReadFile(filepath.Join(workDir, ".terraform.lock.hcl"))
	run.Lock = string(lockBytes)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)
	return nil
}

type InstanceMetrics struct {
	InstanceID       string
	InstanceType     string
	State            string
	LaunchTime       time.Time
	CPUPercent       *float64 // CloudWatch
	DiskSizeGB       int32    // EBS size
	EstimatedCostUSD *float64 // Cost Explorer (best effort)
	Limitations      []string
}

type TFState struct {
	Resources []struct {
		Type      string `json:"type"`
		Instances []struct {
			Attributes struct {
				ID        string   `json:"id"`
				VolumeIDs []string `json:"volume_ids"`
			} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

// ResourceMetrics - Common fields for any AWS resource + EC2-specific
type ResourceMetrics struct {
	Type                 string // e.g., "aws_instance", "aws_s3_bucket"
	ID                   string
	State                string     // if available (e.g., EC2 state)
	LaunchTime           *time.Time // if applicable
	CPUUtilizationAvgPct *float64   // EC2 only, last hour
	DiskSizeGB           int64      // EC2 only
	EstimatedCostUSD     *float64   // last ~30 days, best effort
	Limitations          []string
}

func GetEC2MetricsFromTFState(
	ctx context.Context,
	tfstateJSON []byte,
	awsCfg aws.Config,
) ([]InstanceMetrics, error) {
	var state TFState
	if err := json.Unmarshal(tfstateJSON, &state); err != nil {
		return nil, err
	}
	ec2c := ec2.NewFromConfig(awsCfg)
	cwc := cloudwatch.NewFromConfig(awsCfg)
	cec := costexplorer.NewFromConfig(awsCfg)
	var results []InstanceMetrics
	for _, res := range state.Resources {
		if res.Type != "aws_instance" {
			continue
		}
		for _, inst := range res.Instances {
			instanceID := inst.Attributes.ID
			// --------------------
			// EC2 info
			// --------------------
			di, err := ec2c.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
				InstanceIds: []string{instanceID},
			})
			if err != nil {
				continue
			}
			ec2inst := di.Reservations[0].Instances[0]
			m := InstanceMetrics{
				InstanceID:   instanceID,
				InstanceType: string(ec2inst.InstanceType),
				State:        string(ec2inst.State.Name),
				LaunchTime:   *ec2inst.LaunchTime,
				Limitations:  []string{},
			}
			// --------------------
			// CPU (CloudWatch)
			// --------------------
			cpu, err := cwc.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
				Namespace:  aws.String("AWS/EC2"),
				MetricName: aws.String("CPUUtilization"),
				Dimensions: []cwtypes.Dimension{
					{Name: aws.String("InstanceId"), Value: aws.String(instanceID)},
				},
				StartTime: aws.Time(time.Now().Add(-1 * time.Hour)),
				EndTime:   aws.Time(time.Now()),
				Period:    aws.Int32(300),
				Statistics: []cwtypes.Statistic{
					cwtypes.StatisticAverage,
				},
			})
			if err == nil && len(cpu.Datapoints) > 0 {
				val := *cpu.Datapoints[0].Average
				m.CPUPercent = &val
			} else {
				m.Limitations = append(m.Limitations, "CPU data unavailable")
			}
			// --------------------
			// Disk size (EBS)
			// --------------------
			var totalDisk int32
			for _, bd := range ec2inst.BlockDeviceMappings {
				if bd.Ebs == nil {
					continue
				}
				vol, err := ec2c.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
					VolumeIds: []string{*bd.Ebs.VolumeId},
				})
				if err == nil && len(vol.Volumes) > 0 {
					totalDisk += *vol.Volumes[0].Size
				}
			}
			m.DiskSizeGB = totalDisk
			// --------------------
			// Cost (best effort)
			// --------------------
			cost, err := cec.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
				TimePeriod: &cetypes.DateInterval{
					Start: aws.String(time.Now().AddDate(0, 0, -30).Format("2006-01-02")),
					End:   aws.String(time.Now().Format("2006-01-02")),
				},
				Granularity: cetypes.GranularityMonthly,
				Metrics:     []string{"UnblendedCost"},
				Filter: &cetypes.Expression{
					Dimensions: &cetypes.DimensionValues{
						Key:    cetypes.DimensionResourceId,
						Values: []string{instanceID},
					},
				},
			})
			if err == nil && len(cost.ResultsByTime) > 0 {
				amountStr := cost.ResultsByTime[0].Total["UnblendedCost"].Amount
				if v, err := strconv.ParseFloat(*amountStr, 64); err == nil {
					m.EstimatedCostUSD = &v
				}
			} else {
				m.Limitations = append(m.Limitations, "Cost data delayed or unavailable")
			}
			// --------------------
			// Known hard limitations
			// --------------------
			m.Limitations = append(m.Limitations,
				"Memory usage unavailable without CloudWatch Agent",
				"Filesystem usage unavailable without SSH or SSM",
			)
			results = append(results, m)
		}
	}
	return results, nil
}

// DYNAMIC ROUTE

/*func GetResourceMetricsFromTFStateV2(
	ctx context.Context,
	tfstateJSON []byte,
	cfg aws.Config,
) ([]ResourceMetrics, error) {
	var state TFState
	if err := json.Unmarshal(tfstateJSON, &state); err != nil {
		return nil, fmt.Errorf("failed to parse terraform state: %w", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	cwClient := cloudwatch.NewFromConfig(cfg)
	ceClient := costexplorer.NewFromConfig(cfg)

	var results []ResourceMetrics

	now := time.Now().UTC()
	ceStart := now.AddDate(0, -1, 0).Format("2006-01-02") // ~last 30 days
	ceEnd := now.Format("2006-01-02")

	for _, res := range state.Resources {
		if !strings.HasPrefix(res.Type, "aws_") {
			continue // Only AWS resources
		}

		for _, inst := range res.Instances {
			attrs := inst.Attributes
			idVal, ok := attrs.ID
			if !ok || idVal == "" {
				continue // No usable ID
			}

			m := ResourceMetrics{
				Type:        res.Type,
				ID:          idVal,
				Limitations: []string{},
			}

			// Common: Attempt cost for ANY resource
			costOut, err := ceClient.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
				TimePeriod: &cetypes.DateInterval{
					Start: aws.String(ceStart),
					End:   aws.String(ceEnd),
				},
				Granularity: cetypes.GranularityMonthly,
				Metrics:     []string{"UnblendedCost"},
				Filter: &cetypes.Expression{
					Dimensions: &cetypes.DimensionValues{
						Key:    cetypes.DimensionResourceId,
						Values: []string{idVal},
					},
				},
			})
			if err == nil && len(costOut.ResultsByTime) > 0 && len(costOut.ResultsByTime[0].Total) > 0 {
				amountStr := costOut.ResultsByTime[0].Total["UnblendedCost"].Amount
				if amountStr != nil {
					if v, err := strconv.ParseFloat(*amountStr, 64); err == nil && v > 0 {
						m.EstimatedCostUSD = &v
					}
				}
			}
			if m.EstimatedCostUSD == nil {
				m.Limitations = append(m.Limitations, "Cost data unavailable (enable resource-level granularity in Cost Explorer prefs, or data delayed/no usage)")
			}

			// EC2-specific metrics if applicable
			if res.Type == "aws_instance" {
				// DescribeInstances
				di, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
					InstanceIds: []string{idVal},
				})
				if err == nil && len(di.Reservations) > 0 && len(di.Reservations[0].Instances) > 0 {
					ec2Inst := di.Reservations[0].Instances[0]
					m.State = string(ec2Inst.State.Name)
					if ec2Inst.LaunchTime != nil {
						lt := *ec2Inst.LaunchTime
						m.LaunchTime = &lt
					}

					// CPU
					start := now.Add(-1 * time.Hour)
					cpuStats, err := cwClient.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("CPUUtilization"),
						Dimensions: []cwtypes.Dimension{{Name: aws.String("InstanceId"), Value: aws.String(idVal)}},
						StartTime:  &start,
						EndTime:    &now,
						Period:     aws.Int32(300),
						Statistics: []cwtypes.Statistic{cwtypes.StatisticAverage},
					})
					if err == nil && len(cpuStats.Datapoints) > 0 {
						avg := *cpuStats.Datapoints[len(cpuStats.Datapoints)-1].Average
						m.CPUUtilizationAvgPct = &avg
					} else {
						m.Limitations = append(m.Limitations, "CPU data unavailable")
					}

					// Disk (EBS total size)
					var totalDisk int32
					for _, bd := range ec2Inst.BlockDeviceMappings {
						if bd.Ebs == nil || bd.Ebs.VolumeId == nil {
							continue
						}
						volOut, err := ec2Client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
							VolumeIds: []string{*bd.Ebs.VolumeId},
						})
						if err == nil && len(volOut.Volumes) > 0 {
							totalDisk += *volOut.Volumes[0].Size
						}
					}
					m.DiskSizeGB = totalDisk
				} else {
					m.Limitations = append(m.Limitations, "EC2 details unavailable")
				}
			}

			// Always add common limitations
			m.Limitations = append(m.Limitations,
				"Costs are best-effort; enable resource-level data in Cost Explorer for better accuracy",
				"Use tags or CUR+Athena for precise per-resource breakdown (especially S3)",
				"Memory/disk usage % requires agents/SSM",
			)

			results = append(results, m)
		}
	}

	return results, nil
}*/

// In your traefik.yml (or equivalent static config):
/*```YAML
providers:
  file:
    directory: "/etc/traefik/dynamic/tenants"   # or any path you like
    watch: true
```*/

// Create a template file (e.g. templates/customer-router.yaml.tmpl):

/*
```yaml
# templates/customer-router.yaml.tmpl

http:
  routers:
    {{ .Slug }}-router:
      rule: "Host(`{{ .Slug }}.yourdomain.com`)"   # or Host(`yourdomain.com`) && PathPrefix(`/{{ .Slug }}`)
      service: {{ .Slug }}-service
      entryPoints:
        - websecure
      tls: {}                                      # or tls: { certResolver: myresolver } for Let's Encrypt

  services:
    {{ .Slug }}-service:
      loadBalancer:
        servers:
          - url: "{{ .CloudURL }}"                 # full https://... from your DB
```
*/

// Customer holds the data for templating
type Customer struct {
	Slug     string
	CloudURL string
	// add more fields as needed: Domain, Middlewares, etc.
}

func (app *application) generateCustomerConfig(customer Dict) error {
	tmplPath := "templates/customer-router.yaml"
	if os.Getenv("TRAEFIK_TMPL_PATH") != "" {
		tmplPath = os.Getenv("TRAEFIK_TMPL_PATH")
	}
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		content, err = assets.EmbeddedFiles.ReadFile(tmplPath)
	}
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	outputDir := os.Getenv("TRAEFIK_DYNAMIC_DIR") //"/etc/traefik/dynamic/tenants" // must be writable by your app
	// or use os.Getenv("TRAEFIK_DYNAMIC_DIR") for flexibility
	if outputDir == "" {
		return fmt.Errorf("No TRAEFIK_DYNAMIC_DIR found in your enviromental variables|")
	}
	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", customer["slug"]))

	tmpl, err := app.RenderTextTemplate(string(content), customer)

	// Parse template
	//tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Create / overwrite file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	//if err := tmpl.Execute(file, customer); err != nil {
	//	return fmt.Errorf("execute template: %w", err)
	//}
	_, err = file.WriteString(tmpl)
	if err != nil {
		return fmt.Errorf("write file %s: %w", outputPath, err)
	}
	fmt.Printf("Generated Traefik config for %s → %s\n", customer["slug"], outputPath)
	return nil
}

func (app *application) deleteCustomerConfig(slug string) error {
	outputDir := os.Getenv("TRAEFIK_DYNAMIC_DIR") // "/etc/traefik/dynamic/tenants"
	if outputDir == "" {
		return fmt.Errorf("No TRAEFIK_DYNAMIC_DIR found in your enviromental variables|")
	}
	filePath := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", slug))
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Printf("Removed Traefik config for %s\n", slug)
	return nil
}
