package main

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sync/atomic"
	"time"
)

var ErrQuotaExceeded = errors.New("directory size quota exceeded, writes disabled")

type SizeGuard struct {
	dir        string
	maxBytes   int64
	current    atomic.Int64
	overQuota  atomic.Bool
	stopCh     chan struct{}
}

// NewSizeGuard starts a background monitor that recomputes dir size
// every interval and flips overQuota when it exceeds maxGB.
func NewSizeGuard(dir string, maxGB float64, interval time.Duration) *SizeGuard {
	g := &SizeGuard{
		dir:      dir,
		maxBytes: int64(maxGB * 1024 * 1024 * 1024),
		stopCh:   make(chan struct{}),
	}
	go g.loop(interval)
	return g
}

func (g *SizeGuard) loop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	g.check() // initial check on startup
	for {
		select {
		case <-ticker.C:
			g.check()
		case <-g.stopCh:
			return
		}
	}
}

func (g *SizeGuard) check() {
	var total int64
	filepath.WalkDir(g.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries, don't abort the whole walk
		}
		if !d.IsDir() {
			if info, err := d.Info(); err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	g.current.Store(total)
	g.overQuota.Store(total >= g.maxBytes)
}

// AllowWrite should be checked before any write operation.
func (g *SizeGuard) AllowWrite() error {
	if g.overQuota.Load() {
		return ErrQuotaExceeded
	}
	return nil
}

func (g *SizeGuard) CurrentSize() int64 {
	return g.current.Load()
}

func (g *SizeGuard) Stop() {
	close(g.stopCh)
}

/*
guard := sizeguard.NewSizeGuard("/opt/myapp/data", 10.0, 30*time.Second) // 10GB, checked every 30s
defer guard.Stop()

func writeFile(guard *sizeguard.SizeGuard, path string, data []byte) error {
	if err := guard.AllowWrite(); err != nil {
		return err // reject before touching disk
	}
	return os.WriteFile(path, data, 0644)
}
type SizeGuard struct {
	maxBytes int64
	current  atomic.Int64
	// ...
}

// Call this after every successful write, with the delta in bytes.
func (g *SizeGuard) RecordWrite(deltaBytes int64) {
	g.current.Add(deltaBytes)
}

func (g *SizeGuard) AllowWrite(estimatedBytes int64) error {
	if g.current.Load()+estimatedBytes > g.maxBytes {
		return ErrQuotaExceeded
	}
	return nil
}*/