package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	root string
}

func NewLocalStorage(root string) (*LocalStorage, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("local storage root cannot be empty")
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	return &LocalStorage{root: root}, nil
}

func NewLocalStorageFromEnv() (*LocalStorage, error) {
	root := env("STORAGE_LOCAL_PATH")
	if root == "" {
		root = env("UPLOAD")
		if root == "" {
			root = "static/uploads"
		}
	}
	return NewLocalStorage(root)
}

func (s *LocalStorage) fullPath(name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("invalid storage path: %q", name)
	}

	clean := filepath.Clean(name)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage path: %q", name)
	}

	return filepath.Join(s.root, clean), nil
}

func contextErr(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *LocalStorage) Upload(ctx context.Context, r io.Reader, name string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}

	fullPath, err := s.fullPath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(fullPath), ".storage-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)

	if _, err := io.Copy(temporary, r); err != nil {
		temporary.Close()
		return fmt.Errorf("write file: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if err := os.Rename(temporaryName, fullPath); err != nil {
		return fmt.Errorf("replace file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	fullPath, err := s.fullPath(name)
	if err != nil {
		return nil, err
	}

	return os.Open(fullPath)
}

func (s *LocalStorage) Delete(ctx context.Context, name string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	fullPath, err := s.fullPath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, name string) (bool, error) {
	if err := contextErr(ctx); err != nil {
		return false, err
	}
	fullPath, err := s.fullPath(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

var _ Storage = (*LocalStorage)(nil)
