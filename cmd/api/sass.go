package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type TerraformRun struct {
	Config string
	State  json.RawMessage // store as JSON in DB
	Lock   string
}

// func that checks if the link has not prefix http:// or https:// and adds http:// if missing
func ensureHTTPPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
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
	return str
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
	action := params["data"].(Dict)["action"].(string)
	sql := `select * from "deployment" where "deployment_id" = ? and "active" = true and "excluded" = false`
	// fmt.Println(sql, _data["deployment_id"])
	deployment, err := app.GetRowByFilter(sql, params, []any{_data["deployment_id"]})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	tenantID := _data["tenant_id"]
	sql = `select * from "tenant" where "tenant_id" = ? and "active" = true and "excluded" = false`
	tenant, err := app.GetRowsByFilter(sql, params, []any{tenantID})
	if err != nil {
		fmt.Println("tenant err:", err)
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
	//fmt.Println(tenant, tenantEnv)
	if _, ok := deployment["terraform_template"].(string); !ok {
		msg, _ := app.i18n.T("unable-to-match-deployment-id", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	_tmpl_data := map[string]any{"tenant_id": tenantID, "env": tenantEnv, "envKV": tenantEnvKeyPair, "deployment": deployment, "tenant": tenant, "data": _data}
	parsedTmpl, err := app.RenderTemplate(deployment["terraform_template"].(string), _tmpl_data)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     "terraform_template" + err.Error(),
		}
	}
	// fmt.Printf(parsedTmpl)
	run := &TerraformRun{Config: parsedTmpl}
	if terraform_state, ok := _data["terraform_state"].(string); ok && terraform_state != "" {
		// fmt.Println("State:", _data["terraform_state"])
		run.State = json.RawMessage([]byte(terraform_state))
	}
	if terraform_lock, ok := _data["terraform_lock"].(string); ok && terraform_lock != "" {
		// fmt.Println("Lock:", _data["terraform_state"])
		run.Lock = terraform_lock
	}
	var res map[string]string
	switch action {
	case "deploy":
		res, err = app.DeployTerraformForTenant(params, tenantID, run)
		_json_out, _ := json.Marshal(res)
		_data["tf_public_ip"] = rawMessageToString(json.RawMessage(res["public_ip"]))
		_data["tf_public_dns"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["public_dns"])))
		_data["tf_public_url"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["url"])))
		_data["terraform_outputs"] = string(_json_out)
		_data["deployed"] = true
	case "cancel", "destroy":
		err = app.DestroyTerraform(params, tenantID, run)
		_data["deployed"] = false
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
	//fmt.Println(1, "terraform_state", len(_data["terraform_state"].(string)), "terraform_lock", len(_data["terraform_lock"].(string)), len(run.Lock))
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
	return Dict{
		"success": success,
		"msg":     msg,
		"data":    res,
	}
}

func (app *application) DeployTerraformForTenant(params Dict, tenantID any, run *TerraformRun) (map[string]string, error) {
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
	if err := tf.Apply(context.Background()); err != nil {
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

func (etlx *application) RenderTemplate(tmplStr string, data map[string]any) (string, error) {
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

func (app *application) DestroyTerraform(params Dict, tenantID any, run *TerraformRun) error {
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
	if err := tf.Destroy(context.Background()); err != nil {
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
