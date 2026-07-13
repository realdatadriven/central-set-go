package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

type Runner struct {
	client *ssh.Client
}

/*func NewSSHnoKnownHost(host, user, keyFile string) (*Runner, error) {
	key, err := os.ReadFile(os.ExpandEnv(keyFile))
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
}*/

func NewSSH(host, user, keyFile, hostKey string) (*Runner, error) {
	// fmt.Println(host, user, keyFile, hostKey)
	key, err := os.ReadFile(os.ExpandEnv(keyFile))
	if err != nil {
		key = []byte(keyFile)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	callback, err := knownhosts.New(os.ExpandEnv(hostKey))
	if err != nil {
		return nil, err
	}
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: callback,
	}
	/*
		hostKeyBytes, err := os.ReadFile(os.ExpandEnv(hostKey))
		if err != nil {
			hostKeyBytes = []byte(hostKey)
		}
		// fmt.Println(hostKey, string(hostKeyBytes))
		hostPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(hostKeyBytes)
		if err != nil {
			return nil, err
		}
		cfg := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.FixedHostKey(hostPublicKey),
		}
	*/
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

// upload file
func (r *Runner) Upload(ctx context.Context, localPath, remotePath string) error {
	client, err := sftp.NewClient(r.client)
	if err != nil {
		return fmt.Errorf("SFTP client creation failed: %w", err)
	}
	defer client.Close()
	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer srcFile.Close()
	dstFile, err := client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("could not create remote file: %w", err)
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	return nil
}

func (r *Runner) Download(ctx context.Context, localPath, remotePath string) error {
	client, err := sftp.NewClient(r.client)
	if err != nil {
		return fmt.Errorf("SFTP client creation failed: %w", err)
	}
	defer client.Close()
	srcFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("could not open remote file: %w", err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("could not create local file: %w", err)
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	return nil
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

type Credentials struct {
	PrivateKeyPEM string // PKCS#8 PEM
	PublicKeySSH  string // OpenSSH authorized_keys format
}

func GenerateCredentials() (*Credentials, error) {
	// Generate Ed25519 key pair.
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	// Encode the private key as PKCS#8.
	pkcs8, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8,
	})
	// Encode the public key in OpenSSH format.
	pubKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, err
	}
	publicSSH := ssh.MarshalAuthorizedKey(pubKey)
	return &Credentials{
		PrivateKeyPEM: string(privatePEM),
		PublicKeySSH:  string(publicSSH),
	}, nil
}

func AddKnownHost(host string, port int, knownHosts string) error {
	cmd := exec.Command(
		"ssh-keyscan",
		"-H",
		"-p",
		fmt.Sprint(port),
		host,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	f, err := os.OpenFile(
		knownHosts,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(out.Bytes())
	return err
}

/*cred, err := sshutil.GenerateCredentials()
if err != nil {
	log.Fatal(err)
}

fmt.Println("PRIVATE KEY:")
fmt.Println(cred.PrivateKeyPEM)

fmt.Println("PUBLIC KEY:")
fmt.Println(cred.PublicKeySSH)

err = sshutil.AddKnownHost(
	"csete.online",
	22,
	"./deployment/demo/known_hosts",
)
if err != nil {
	log.Fatal(err)
}
*/