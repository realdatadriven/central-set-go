package replication

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// internal/replication/config.go
type Config struct {
	DBs []DBConfig `yaml:"dbs"`
}

type DBConfig struct {
	Path    string        `yaml:"path"`
	Replica ReplicaConfig `yaml:"replica"` // singular now
}

type ReplicaConfig struct {
	Type   string `yaml:"type"`
	Bucket string `yaml:"bucket"`
	Path   string `yaml:"path"`
	Region string `yaml:"region"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading litestream config %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parsing litestream config %q: %w", path, err)
	}
	if len(cfg.DBs) == 0 {
		return nil, fmt.Errorf("litestream config %q defines no databases", path)
	}
	return &cfg, nil
}