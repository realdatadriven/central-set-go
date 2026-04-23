package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/terraform-exec/tfexec"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	cetypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/realdatadriven/central-set-go/assets"
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
	sql := `select * from "plan" where "plan_id" = ? and "active" = true and "excluded" = false`
	// fmt.Println(sql, _data["plan_id"])
	plan, err := app.GetRowByFilter(sql, params, []any{_data["plan_id"]})
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
	if _, ok := plan["terraform_template"].(string); !ok {
		msg, _ := app.i18n.T("unable-to-match-plan-id", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	_tmpl_data := map[string]any{"tenant_id": tenantID, "env": tenantEnv, "envKV": tenantEnvKeyPair, "plan": plan, "tenant": tenant, "data": _data}
	parsedTmpl, err := app.RenderTemplate(plan["terraform_template"].(string), _tmpl_data)
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
		res, err = app.DeployOpenTofuForTenant(params, tenantID, run)
		_json_out, _ := json.Marshal(res)
		_data["tf_public_ip"] = rawMessageToString(json.RawMessage(res["public_ip"]))
		_data["tf_public_dns"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["public_dns"])))
		_data["tf_public_url"] = ensureHTTPPrefix(rawMessageToString(json.RawMessage(res["url"])))
		_data["terraform_outputs"] = string(_json_out)
		_data["deployed"] = true
	case "cancel", "destroy":
		err = app.DestroyOpenTofu(params, tenantID, run)
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

func (app *application) DeployOpenTofuForTenant(params Dict, tenantID any, run *TerraformRun) (map[string]string, error) {
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
	//fmt.Printf("OpenTofu outputs: %v\n", result)
	//fmt.Printf("OpenTofu state: %s\n", string(run.State))
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

func (app *application) DestroyOpenTofu(params Dict, tenantID any, run *TerraformRun) error {
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

	tmpl, err := app.RenderTemplate(string(content), customer)

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
