package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

type CommandRunner interface {
	Run(ctx context.Context, cmd string) error
	RunOutput(ctx context.Context, cmd string) (string, error)
	Close() error
}

type LocalRunner struct{}

func NewLocal() *LocalRunner {
	return &LocalRunner{}
}

func (r *LocalRunner) Close() error {
	return nil
}

func (r *LocalRunner) Run(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func (r *LocalRunner) RunOutput(ctx context.Context, cmd string) (string, error) {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)

	var out bytes.Buffer

	c.Stdout = &out
	c.Stderr = &out

	err := c.Run()

	return out.String(), err
}

type ServiceManager struct {
	runner CommandRunner
}

func NewServiceManager(r CommandRunner) *ServiceManager {
	return &ServiceManager{
		runner: r,
	}
}

func (s *ServiceManager) systemctl(ctx context.Context, action, service string) error {
	return s.runner.Run(
		ctx,
		fmt.Sprintf("systemctl --user %s %s", action, service),
	)
}

func (s *ServiceManager) Start(ctx context.Context, service string) error {
	return s.systemctl(ctx, "start", service)
}

func (s *ServiceManager) Stop(ctx context.Context, service string) error {
	return s.systemctl(ctx, "stop", service)
}

func (s *ServiceManager) Restart(ctx context.Context, service string) error {
	return s.systemctl(ctx, "restart", service)
}

func (s *ServiceManager) Status(ctx context.Context, service string) (string, error) {
	return s.runner.RunOutput(
		ctx,
		fmt.Sprintf("systemctl --user status %s --no-pager", service),
	)
}

func (s *ServiceManager) Run(ctx context.Context, cmd string) error {
	return s.runner.Run(ctx, cmd)
}

func (s *ServiceManager) RunOutput(ctx context.Context, cmd string) (string, error) {
	return s.runner.RunOutput(ctx, cmd)
}

func (s *ServiceManager) Enable(ctx context.Context, service string) error {
	return s.systemctl(ctx, "enable", service)
}

func (s *ServiceManager) Disable(ctx context.Context, service string) error {
	return s.systemctl(ctx, "disable", service)
}

func (s *ServiceManager) Logs(ctx context.Context, service string, lines int) (string, error) {
	return s.RunOutput(
		ctx,
		fmt.Sprintf(
			"journalctl --user -u %s -n %d --no-pager",
			service,
			lines,
		),
	)
}
