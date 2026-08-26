package replication

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBs []DBConfig `yaml:"dbs"`
}

type DBConfig struct {
	Path      string        `yaml:"path"`      // single-file mode (optional)
	Dir       string        `yaml:"dir"`        // directory mode
	Pattern   string        `yaml:"pattern"`
	Recursive bool          `yaml:"recursive"`
	Watch     bool          `yaml:"watch"`
	MetaDir   string        `yaml:"meta-dir"`
	Replica   ReplicaConfig `yaml:"replica"`
}

type ReplicaConfig struct {
	Type           string   `yaml:"type"`
	Bucket         string   `yaml:"bucket"`
	Path           string   `yaml:"path"`
	Region         string   `yaml:"region"`
	Endpoint       string   `yaml:"endpoint"`
	ForcePathStyle bool     `yaml:"force-path-style"`
	SyncInterval   Duration `yaml:"sync-interval"`
}

// Duration lets us parse YAML strings like "1s" into a time.Duration,
// since yaml.v3 won't do this for a plain time.Duration field.
type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("parsing duration %q: %w", s, err)
	}
	*d = Duration(parsed)
	return nil
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