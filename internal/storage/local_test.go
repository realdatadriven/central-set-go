package storage

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestLocalStorageRoundTrip(t *testing.T) {
	backend, err := NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := backend.Upload(ctx, strings.NewReader("hello"), "nested/file.txt"); err != nil {
		t.Fatal(err)
	}
	if exists, err := backend.Exists(ctx, "nested/file.txt"); err != nil || !exists {
		t.Fatalf("Exists() = %v, %v", exists, err)
	}
	reader, err := backend.Download(ctx, "nested/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	contents, err := io.ReadAll(reader)
	if err != nil || string(contents) != "hello" {
		t.Fatalf("Download() = %q, %v", contents, err)
	}
	if err := backend.Delete(ctx, "nested/file.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestLocalStorageRejectsPathTraversal(t *testing.T) {
	backend, err := NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Upload(context.Background(), strings.NewReader("x"), "../outside"); err == nil {
		t.Fatal("Upload() accepted a path outside the storage root")
	}
	if err := backend.Delete(context.Background(), "../outside"); err == nil {
		t.Fatal("Delete() accepted a path outside the storage root")
	}
}
