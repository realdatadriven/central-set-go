package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

type DropboxStorage struct {
	client      *http.Client
	accessToken string
	root        string
}

func NewDropboxStorageFromEnv() (*DropboxStorage, error) {
	token, err := requiredEnv("DROPBOX_ACCESS_TOKEN")
	if err != nil {
		return nil, err
	}
	root := env("DROPBOX_ROOT")
	if root == "" {
		root = "/uploads"
	}
	return &DropboxStorage{client: http.DefaultClient, accessToken: token, root: "/" + strings.Trim(root, "/")}, nil
}

func (s *DropboxStorage) fullPath(name string) string {
	return path.Join(s.root, strings.TrimLeft(name, "/"))
}

func (s *DropboxStorage) request(ctx context.Context, endpoint string, arg any, content io.Reader) (*http.Response, error) {
	var body io.Reader
	if arg != nil {
		encoded, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(encoded))
	}
	if content != nil {
		body = content
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.dropboxapi.com/2/"+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	if content == nil {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	return s.client.Do(req)
}

func dropboxError(resp *http.Response) error {
	defer resp.Body.Close()
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, readErrorBody(resp.Body))
}

func (s *DropboxStorage) Upload(ctx context.Context, r io.Reader, name string) error {
	arg := map[string]any{"path": s.fullPath(name), "mode": "overwrite", "autorename": false, "mute": true}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://content.dropboxapi.com/2/files/upload", r)
	if err != nil {
		return err
	}
	encoded, _ := json.Marshal(arg)
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", string(encoded))
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("dropbox upload %q: %w", name, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dropbox upload %q: %w", name, dropboxError(resp))
	}
	resp.Body.Close()
	return nil
}

func (s *DropboxStorage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	resp, err := s.request(ctx, "files/download", map[string]string{"path": s.fullPath(name)}, nil)
	if err != nil {
		return nil, fmt.Errorf("dropbox download %q: %w", name, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, dropboxError(resp)
	}
	return resp.Body, nil
}

func (s *DropboxStorage) Delete(ctx context.Context, name string) error {
	resp, err := s.request(ctx, "files/delete_v2", map[string]string{"path": s.fullPath(name)}, nil)
	if err != nil {
		return fmt.Errorf("dropbox delete %q: %w", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return dropboxError(resp)
	}
	resp.Body.Close()
	return nil
}

func (s *DropboxStorage) Exists(ctx context.Context, name string) (bool, error) {
	resp, err := s.request(ctx, "files/get_metadata", map[string]string{"path": s.fullPath(name)}, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, dropboxError(resp)
	}
	resp.Body.Close()
	return true, nil
}

var _ Storage = (*DropboxStorage)(nil)
