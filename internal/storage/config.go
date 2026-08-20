package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func env(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func requiredEnv(key string) (string, error) {
	value := env(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return value, nil
}

// NewFromEnv creates the backend selected by STORAGE_DRIVER.
func NewFromEnv(ctx context.Context) (Storage, error) {
	switch strings.ToLower(env("STORAGE_DRIVER")) {
	case "", "local":
		return NewLocalStorageFromEnv()
	case "s3":
		return NewS3StorageFromEnv(ctx)
	case "google_drive", "googledrive", "drive":
		return NewGoogleDriveStorageFromEnv(ctx)
	case "dropbox":
		return NewDropboxStorageFromEnv()
	case "onedrive", "one_drive":
		return NewOneDriveStorageFromEnv(ctx)
	default:
		return nil, fmt.Errorf("unsupported STORAGE_DRIVER %q; supported values: local, s3, google_drive, dropbox, onedrive", env("STORAGE_DRIVER"))
	}
}
