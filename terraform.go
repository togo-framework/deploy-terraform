// Package terraform is a generic Terraform (IaC) driver for togo deploy. It runs
// `terraform init/plan/apply` (and `destroy`) in a Terraform directory, so any of
// the cloud deploy plugins can express infrastructure as Terraform. Select with
// DEPLOY_PROVIDER=terraform; set the dir via spec.Options["tf_dir"], TF_DIR, or
// <project>/infra.
package terraform

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/togo-framework/deploy"
	"github.com/togo-framework/togo"
)

func init() {
	deploy.RegisterDriver("terraform", func(k *togo.Kernel) (deploy.Deployer, error) {
		return &driver{dir: os.Getenv("TF_DIR")}, nil
	})
}

type driver struct{ dir string }

func (d *driver) tfDir(spec deploy.Spec) string {
	if v, ok := spec.Options["tf_dir"].(string); ok && v != "" {
		return v
	}
	if d.dir != "" {
		return d.dir
	}
	base := spec.Dir
	if base == "" {
		base = "."
	}
	return filepath.Join(base, "infra")
}

func tf(ctx context.Context, dir string, args ...string) (string, error) {
	if _, err := exec.LookPath("terraform"); err != nil {
		return "", fmt.Errorf("deploy-terraform: terraform not found on PATH")
	}
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	return out.String(), cmd.Run()
}

func varArgs(spec deploy.Spec) []string {
	var a []string
	for k, v := range spec.Env {
		a = append(a, "-var", k+"="+v)
	}
	if vs, ok := spec.Options["vars"].(map[string]string); ok {
		for k, v := range vs {
			a = append(a, "-var", k+"="+v)
		}
	}
	return a
}

func (d *driver) Provision(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	dir := d.tfDir(spec)
	if out, err := tf(ctx, dir, "init", "-input=false"); err != nil {
		return nil, fmt.Errorf("terraform init: %w\n%s", err, out)
	}
	args := append([]string{"apply", "-auto-approve", "-input=false"}, varArgs(spec)...)
	out, err := tf(ctx, dir, args...)
	if err != nil {
		return nil, fmt.Errorf("terraform apply: %w\n%s", err, out)
	}
	return &deploy.Result{Message: "terraform applied (" + dir + ")", Raw: map[string]any{"output": out}}, nil
}

// Deploy = Provision for the IaC driver (apply brings infra to desired state).
func (d *driver) Deploy(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	return d.Provision(ctx, spec)
}

func (d *driver) Destroy(ctx context.Context, spec deploy.Spec) error {
	args := append([]string{"destroy", "-auto-approve", "-input=false"}, varArgs(spec)...)
	out, err := tf(ctx, d.tfDir(spec), args...)
	if err != nil {
		return fmt.Errorf("terraform destroy: %w\n%s", err, out)
	}
	return nil
}

func (d *driver) Status(ctx context.Context, spec deploy.Spec) (*deploy.Status, error) {
	out, err := tf(ctx, d.tfDir(spec), "state", "list")
	return &deploy.Status{Healthy: err == nil, Detail: out}, nil
}
