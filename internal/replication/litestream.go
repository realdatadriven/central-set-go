package replication

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benbjohnson/litestream"
	lss3 "github.com/benbjohnson/litestream/s3"
)

// Manager owns the Litestream store and lets the app shut it down cleanly.
type Manager struct {
	store *litestream.Store
}

// Enabled reports whether USE_LITESTREAM=true.
func Enabled() bool {
	return os.Getenv("USE_LITESTREAM") == "true"
}

// Start builds DB/Replica pairs from LITESTREAM_CONFIG_FILE, opens the
// store (which starts its own internal replication/compaction goroutines),
// and additionally launches a supervisor goroutine that logs sync status.
// It returns immediately; call Wait or just keep the returned Manager
// around until shutdown.
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

	var dbs []*litestream.DB
	for _, dbCfg := range cfg.DBs {
		db := litestream.NewDB(dbCfg.Path)

		// internal/replication/litestream.go — inside Start()
		for _, dbCfg := range cfg.DBs {
			fmt.Println(dbCfg)
			db := litestream.NewDB(dbCfg.Path)
			if dbCfg.Replica.Type != "s3" && dbCfg.Replica.Type != "" {
				return nil, fmt.Errorf("unsupported replica type %q for %s", dbCfg.Replica.Type, dbCfg.Path)
			}
			client := lss3.NewReplicaClient()
			client.Bucket = dbCfg.Replica.Bucket
			client.Path = dbCfg.Replica.Path
			client.Region = dbCfg.Replica.Region
			client.AccessKeyID = accessKey
			client.SecretAccessKey = secretKey
			replica := litestream.NewReplicaWithClient(db, client)
			db.Replica = replica // <-- this must be set before store.Open()
			if err := db.EnsureExists(ctx); err != nil {
				return nil, fmt.Errorf("restoring %s from replica: %w", dbCfg.Path, err)
			}
			dbs = append(dbs, db)
		}
		dbs = append(dbs, db)
	}
	levels := litestream.CompactionLevels{
		{Level: 0},
		{Level: 1, Interval: 10 * time.Second},
	}
	store := litestream.NewStore(dbs, levels)
	if err := store.Open(ctx); err != nil {
		return nil, fmt.Errorf("opening litestream store: %w", err)
	}
	m := &Manager{store: store}
	// Supervisor goroutine: periodic sync-status logging inside the binary,
	// as you asked for. The actual replication/compaction loops are already
	// running inside store.Open() above.
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, db := range dbs {
					status, err := db.SyncStatus(ctx)
					if err != nil {
						log.Printf("litestream: sync status error for %s: %v", db.Path(), err)
						continue
					}
					log.Printf("litestream: %s local=%d remote=%d in_sync=%t", db.Path(), status.LocalTXID, status.RemoteTXID, status.InSync)
				}
			}
		}
	}()

	return m, nil
}

// Close flushes and stops replication. Call this on graceful shutdown,
// after your app has closed its own DB connections.
func (m *Manager) Close(ctx context.Context) error {
	if m == nil || m.store == nil {
		return nil
	}
	return m.store.Close(ctx)
}
