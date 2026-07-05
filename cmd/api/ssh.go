package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

type Runner struct {
	client *ssh.Client
}

func NewSSH(host, user, keyFile string) (*Runner, error) {
	key, err := os.ReadFile(keyFile)
	if err != nil {
		key = []byte(keyFile)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host+":22", cfg)
	if err != nil {
		return nil, err
	}
	return &Runner{
		client: client,
	}, nil
}

func (r *Runner) Close() error {
	return r.client.Close()
}

func (r *Runner) Run(ctx context.Context, cmd string) error {
	session, err := r.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	if err := session.Start(cmd); err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	done := make(chan error, 1)
	go func() {
		done <- session.Wait()
	}()
	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (r *Runner) RunOutput(ctx context.Context, cmd string) (string, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var out bytes.Buffer
	session.Stdout = &out
	session.Stderr = &out
	err = session.Run(cmd)
	return out.String(), err
}

func (r *Runner) systemctl(ctx context.Context, action, service string) error {
	return r.Run(ctx, fmt.Sprintf("systemctl --user %s %s", action, service))
}

func (r *Runner) Start(ctx context.Context, service string) error {
	return r.systemctl(ctx, "start", service)
}

func (r *Runner) Stop(ctx context.Context, service string) error {
	return r.systemctl(ctx, "stop", service)
}

func (r *Runner) Restart(ctx context.Context, service string) error {
	return r.systemctl(ctx, "restart", service)
}

func (r *Runner) Enable(ctx context.Context, service string) error {
	return r.systemctl(ctx, "enable", service)
}

func (r *Runner) Disable(ctx context.Context, service string) error {
	return r.systemctl(ctx, "disable", service)
}

func (r *Runner) Status(ctx context.Context, service string) (string, error) {
	return r.RunOutput(ctx, fmt.Sprintf("systemctl --user status %s --no-pager", service))
}

func (r *Runner) Logs(ctx context.Context, service string, lines int) (string, error) {
	return r.RunOutput(ctx, fmt.Sprintf("journalctl --user -u %s -n %d --no-pager", service, lines))
}

/*
	ctx := context.Background()
	ssh, err := NewSSH(
		"csete.online",
		"clovis",
		os.ExpandEnv("$HOME/.ssh/id_ed25519"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer ssh.Close()

	fmt.Println("Stopping...")
	if err := ssh.Stop(ctx, "c7"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting...")
	if err := ssh.Start(ctx, "c7"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Restarting Traefik...")
	if err := ssh.Restart(ctx, "traefik"); err != nil {
		log.Fatal(err)
	}

	status, err := ssh.Status(ctx, "c7")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(status)

	logs, err := ssh.Logs(ctx, "c7", 20)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(logs)

	// Execute any command
	err = ssh.Run(ctx, `
echo "Current directory:"
pwd

echo
echo "Files:"
ls -la
`)
	if err != nil {
		log.Fatal(err)
	}
*/
