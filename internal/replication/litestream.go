package replication

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benbjohnson/litestream"
	lss3 "github.com/benbjohnson/litestream/s3"
	"github.com/fsnotify/fsnotify"
)

type Manager struct {
	store   *litestream.Store
	watcher *fsnotify.Watcher
}

func Enabled() bool {
	return os.Getenv("USE_LITESTREAM") == "true"
}

// Start reads LITESTREAM_CONFIG_FILE, opens an (initially empty) Store, then
// scans + registers every matching database. If watch: true, it also starts
// an fsnotify goroutine that registers new files and unregisters deleted ones
// as they happen — this is the "goroutine inside the binary" doing the work.
func Start(ctx context.Context) (*Manager, error) {
	cfgPath := os.Getenv("LITESTREAM_CONFIG_FILE")
	if cfgPath == "" {
		return nil, fmt.Errorf("LITESTREAM_CONFIG_FILE is not set")
	}
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	accessKey := os.Getenv("LITESTREAM_AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("LITESTREAM_AWS_SECRET_ACCESS_KEY")

	levels := litestream.CompactionLevels{
		{Level: 0},
		{Level: 1, Interval: 10 * time.Second},
	}
	// Start with no DBs — we register them ourselves, below and dynamically.
	store := litestream.NewStore(nil, levels)
	if err := store.Open(ctx); err != nil {
		return nil, fmt.Errorf("opening litestream store: %w", err)
	}

	m := &Manager{store: store}

	for _, dbCfg := range cfg.DBs {
		if dbCfg.Dir == "" {
			return nil, fmt.Errorf("only directory-mode ('dir') configs are supported by this loader")
		}
		if err := m.scanAndRegisterDir(ctx, dbCfg, accessKey, secretKey); err != nil {
			return nil, err
		}
		if dbCfg.Watch {
			if err := m.watchDir(ctx, dbCfg, accessKey, secretKey); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m *Manager) scanAndRegisterDir(ctx context.Context, dbCfg DBConfig, accessKey, secretKey string) error {
	return filepath.WalkDir(dbCfg.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !dbCfg.Recursive && path != dbCfg.Dir {
				return filepath.SkipDir
			}
			return nil
		}
		matched, _ := filepath.Match(dbCfg.Pattern, filepath.Base(path))
		if !matched || !isSQLiteFile(path) {
			return nil
		}
		return m.registerDB(ctx, dbCfg, path, accessKey, secretKey)
	})
}

func (m *Manager) registerDB(ctx context.Context, dbCfg DBConfig, path, accessKey, secretKey string) error {
	if m.store.FindDB(path) != nil {
		return nil // already registered
	}

	relPath, err := filepath.Rel(dbCfg.Dir, path)
	if err != nil {
		return err
	}
	relPath = filepath.ToSlash(relPath)

	db := litestream.NewDB(path)
	if dbCfg.MetaDir != "" {
		db.SetMetaPath(filepath.Join(dbCfg.MetaDir, relPath+"-litestream"))
	}

	client := lss3.NewReplicaClient()
	client.Bucket = dbCfg.Replica.Bucket
	client.Path = strings.TrimSuffix(dbCfg.Replica.Path, "/") + "/" + relPath
	client.Region = dbCfg.Replica.Region
	client.Endpoint = dbCfg.Replica.Endpoint
	client.ForcePathStyle = dbCfg.Replica.ForcePathStyle
	client.AccessKeyID = accessKey
	client.SecretAccessKey = secretKey

	replica := litestream.NewReplicaWithClient(db, client)
	if dbCfg.Replica.SyncInterval > 0 {
		replica.SyncInterval = time.Duration(dbCfg.Replica.SyncInterval)
	}
	db.Replica = replica // must be set before EnsureExists/RegisterDB

	if err := db.EnsureExists(ctx); err != nil {
		return fmt.Errorf("restoring %s from replica: %w", path, err)
	}
	if err := m.store.RegisterDB(db); err != nil {
		return fmt.Errorf("registering %s: %w", path, err)
	}

	log.Printf("litestream: replicating %s -> s3://%s/%s", path, client.Bucket, client.Path)
	return nil
}

func (m *Manager) watchDir(ctx context.Context, dbCfg DBConfig, accessKey, secretKey string) error {
	if m.watcher == nil {
		w, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("creating fsnotify watcher: %w", err)
		}
		m.watcher = w
	}

	if err := m.watcher.Add(dbCfg.Dir); err != nil {
		return fmt.Errorf("watching %s: %w", dbCfg.Dir, err)
	}
	if dbCfg.Recursive {
		_ = filepath.WalkDir(dbCfg.Dir, func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				_ = m.watcher.Add(path)
			}
			return nil
		})
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-m.watcher.Events:
				if !ok {
					return
				}
				switch {
				case event.Op&fsnotify.Create != 0:
					matched, _ := filepath.Match(dbCfg.Pattern, filepath.Base(event.Name))
					if !matched {
						continue
					}
					time.Sleep(200 * time.Millisecond) // let the writer finish
					if isSQLiteFile(event.Name) {
						if err := m.registerDB(ctx, dbCfg, event.Name, accessKey, secretKey); err != nil {
							log.Printf("litestream: failed to register %s: %v", event.Name, err)
						}
					}
				case event.Op&fsnotify.Remove != 0:
					if m.store.FindDB(event.Name) != nil {
						if err := m.store.UnregisterDB(ctx, event.Name); err != nil {
							log.Printf("litestream: failed to unregister %s: %v", event.Name, err)
						}
					}
				}
			case err, ok := <-m.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("litestream: watcher error: %v", err)
			}
		}
	}()

	return nil
}

func isSQLiteFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	header := make([]byte, 16)
	if _, err := f.Read(header); err != nil {
		return false
	}
	return string(header) == "SQLite format 3\x00"
}

func (m *Manager) Close(ctx context.Context) error {
	if m.watcher != nil {
		_ = m.watcher.Close()
	}
	if m.store == nil {
		return nil
	}
	return m.store.Close(ctx)
}