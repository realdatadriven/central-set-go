package storage

import (
	"context"
	"io"
)

// Storage is the common contract implemented by every file backend.
type Storage interface {
	Upload(ctx context.Context, r io.Reader, path string) error
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
}
