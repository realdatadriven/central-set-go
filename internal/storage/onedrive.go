package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type OneDriveStorage struct {
	client      *http.Client
	accessToken string
	driveID     string
	root        string
}

func NewOneDriveStorageFromEnv(ctx context.Context) (*OneDriveStorage, error) {
	token := env("ONEDRIVE_ACCESS_TOKEN")
	if token == "" {
		var err error
		token, err = oneDriveToken(ctx)
		if err != nil {
			return nil, err
		}
	}
	root := env("ONEDRIVE_ROOT")
	if root == "" {
		root = "/uploads"
	}
	return &OneDriveStorage{client: http.DefaultClient, accessToken: token, driveID: env("ONEDRIVE_DRIVE_ID"), root: "/" + strings.Trim(root, "/")}, nil
}

func oneDriveToken(ctx context.Context) (string, error) {
	tenant, err := requiredEnv("ONEDRIVE_TENANT_ID")
	if err != nil {
		return "", err
	}
	clientID, err := requiredEnv("ONEDRIVE_CLIENT_ID")
	if err != nil {
		return "", err
	}
	secret, err := requiredEnv("ONEDRIVE_CLIENT_SECRET")
	if err != nil {
		return "", err
	}
	form := url.Values{"client_id": {clientID}, "client_secret": {secret}, "scope": {"https://graph.microsoft.com/.default"}, "grant_type": {"client_credentials"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://login.microsoftonline.com/"+url.PathEscape(tenant)+"/oauth2/v2.0/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request OneDrive token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("request OneDrive token returned HTTP %d: %s", resp.StatusCode, readErrorBody(resp.Body))
	}
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil || tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("OneDrive token response did not contain access_token")
	}
	return tokenResponse.AccessToken, nil
}

func (s *OneDriveStorage) itemURL(name string) string {
	itemPath := url.PathEscape(strings.TrimPrefix(path.Join(s.root, strings.TrimLeft(name, "/")), "/"))
	if s.driveID != "" {
		return "https://graph.microsoft.com/v1.0/drives/" + url.PathEscape(s.driveID) + "/root:/" + itemPath
	}
	return "https://graph.microsoft.com/v1.0/me/drive/root:/" + itemPath
}

func (s *OneDriveStorage) request(ctx context.Context, method, requestURL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	return s.client.Do(req)
}

func (s *OneDriveStorage) Upload(ctx context.Context, r io.Reader, name string) error {
	resp, err := s.request(ctx, http.MethodPut, s.itemURL(name)+":/content", r)
	if err != nil {
		return fmt.Errorf("onedrive upload %q: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("onedrive upload %q returned HTTP %d: %s", name, resp.StatusCode, readErrorBody(resp.Body))
	}
	return nil
}

func (s *OneDriveStorage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	resp, err := s.request(ctx, http.MethodGet, s.itemURL(name)+":/content", nil)
	if err != nil {
		return nil, fmt.Errorf("onedrive download %q: %w", name, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("onedrive download %q returned HTTP %d: %s", name, resp.StatusCode, readErrorBody(resp.Body))
	}
	return resp.Body, nil
}

func (s *OneDriveStorage) Delete(ctx context.Context, name string) error {
	resp, err := s.request(ctx, http.MethodDelete, s.itemURL(name), nil)
	if err != nil {
		return fmt.Errorf("onedrive delete %q: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("onedrive delete %q returned HTTP %d: %s", name, resp.StatusCode, readErrorBody(resp.Body))
	}
	return nil
}

func (s *OneDriveStorage) Exists(ctx context.Context, name string) (bool, error) {
	resp, err := s.request(ctx, http.MethodGet, s.itemURL(name), nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}
	return false, fmt.Errorf("onedrive exists %q returned HTTP %d: %s", name, resp.StatusCode, readErrorBody(resp.Body))
}

func readErrorBody(r io.Reader) string {
	var response struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(r).Decode(&response); err == nil && response.Error.Message != "" {
		return response.Error.Message
	}
	return "unknown error"
}

var _ Storage = (*OneDriveStorage)(nil)
