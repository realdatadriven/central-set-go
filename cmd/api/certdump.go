package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// internal/certdump/CertDumpConfig.go
type CertDumpConfig struct {
	Enabled   bool
	AcmeJSON  string
	Domains   []string
	OutputDir string
}

func CertDumpLoadConfig() CertDumpConfig {
	return CertDumpConfig{
		Enabled:   parseBool(os.Getenv("ENABLE_TRAEFIK_CERT_DUMP")),
		AcmeJSON:  os.Getenv("TRAEFIK_ACME_JSON_PATH"),
		Domains:   splitCSV(os.Getenv("TRAEFIK_DOMAINS")),
		OutputDir: os.Getenv("TRAEFIK_CERT_DUMPS"),
	}
}
func parseBool(v string) bool {
	return strings.EqualFold(v, "true")
}
func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	items := strings.Split(v, ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}
	return items
}

// internal/certdump/dump.go

type acmeStorage struct {
	Resolvers map[string]resolver `json:"-"`
}
type resolver struct {
	Certificates []certificate `json:"Certificates"`
}

type certificate struct {
	Domain struct {
		Main string `json:"main"`
	} `json:"domain"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

func DumpCertificates(acmeJSON string, domains []string, output string) error {
	data, err := os.ReadFile(acmeJSON)
	if err != nil {
		return err
	}
	var raw map[string]resolver
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	domainSet := map[string]bool{}
	for _, d := range domains {
		domainSet[d] = true
	}
	for _, resolver := range raw {
		for _, cert := range resolver.Certificates {
			domain := cert.Domain.Main
			if !domainSet[domain] {
				continue
			}
			certBytes, err := base64.StdEncoding.DecodeString(cert.Certificate)
			if err != nil {
				return err
			}
			keyBytes, err := base64.StdEncoding.DecodeString(cert.Key)
			if err != nil {
				return err
			}
			dir := filepath.Join(output, domain)
			if err := os.MkdirAll(dir, 0700); err != nil {
				return err
			}
			if err := atomicWrite(filepath.Join(dir, "cert.pem"), certBytes, 0644); err != nil {
				return err
			}
			if err := atomicWrite(filepath.Join(dir, "key.pem"), keyBytes, 0600); err != nil {
				return err
			}
			fmt.Println("certificate exported:", domain)
		}
	}
	return nil
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// internal/certdump/watcher.go
func CertDumpWatch(ctx context.Context, cfg CertDumpConfig) error {
	if err := DumpCertificates(cfg.AcmeJSON, cfg.Domains, cfg.OutputDir); err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	dir := filepath.Dir(cfg.AcmeJSON)
	if err := watcher.Add(dir); err != nil {
		return err
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Name != cfg.AcmeJSON {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				time.Sleep(time.Second)
				DumpCertificates(cfg.AcmeJSON, cfg.Domains, cfg.OutputDir)
			}
		case err := <-watcher.Errors:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}
