package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type TerraformRun struct {
	Config string
	State  json.RawMessage // store as JSON in DB
}

func (app *application) DeployTerraformForTenant(tenantID string, run *TerraformRun) (map[string]string, error) {
	workDir, _ := os.MkdirTemp("", "tf-*")
	os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644)

	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, err
	}
	tf, _ := tfexec.NewTerraform(workDir, tfPath)

	// 👇 Load tenant-specific env vars
	//tenantEnv, _ := getTenantEnvVars(tenantID)
	sql := `select * from "env" where "costomer_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.AdminGetRowsByFilter(sql, []any{tenantID})
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
	tf.SetEnv(mergedEnv)

	// Normal Terraform lifecycle
	if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
		return nil, err
	}
	if err := tf.Apply(context.Background()); err != nil {
		return nil, err
	}

	outputs, _ := tf.Output(context.Background())
	result := map[string]string{}
	for k, v := range outputs {
		result[k] = fmt.Sprintf("%v", v.Value)
	}

	stateBytes, _ := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	tfexec.CleanEnv(mergedEnv)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)
	//fmt.Printf("Terraform outputs: %v\n", result)
	//fmt.Printf("Terraform state: %s\n", string(run.State))
	return result, nil
}

func (app *application) DestroyTerraform(tenantID string, run *TerraformRun) error {
	// 1. Create temp dir
	workDir, err := os.MkdirTemp("", "tf-demo-*")
	if err != nil {
		return err
	}

	// 2. Write config
	if err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(run.Config), 0644); err != nil {
		return err
	}

	// 3. Restore state
	if err := os.WriteFile(filepath.Join(workDir, "terraform.tfstate"), run.State, 0644); err != nil {
		return err
	}

	// 4. Run terraform destroy
	tfPath, _ := exec.LookPath("terraform")
	tf, _ := tfexec.NewTerraform(workDir, tfPath)

	sql := `select * from "env" where "costomer_id" = ? and "active" = true and "excluded" = false`
	tenantEnv, err := app.AdminGetRowsByFilter(sql, []any{tenantID})
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
	tf.SetEnv(mergedEnv)

	if err := tf.Destroy(context.Background()); err != nil {
		return fmt.Errorf("destroy failed: %w", err)
	}

	tfexec.CleanEnv(mergedEnv)

	stateBytes, _ := os.ReadFile(filepath.Join(workDir, "terraform.tfstate"))
	run.State = json.RawMessage(stateBytes)
	// Clean up temp directory after use
	defer os.RemoveAll(workDir)

	return nil
}
