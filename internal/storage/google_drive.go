package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveStorage struct {
	service  *drive.Service
	folderID string
}

func NewGoogleDriveStorageFromEnv(ctx context.Context) (*GoogleDriveStorage, error) {
	credentialsFile, err := requiredEnv("GOOGLE_DRIVE_CREDENTIALS_FILE")
	if err != nil {
		return nil, err
	}
	folderID, err := requiredEnv("GOOGLE_DRIVE_FOLDER_ID")
	if err != nil {
		return nil, err
	}
	service, err := drive.NewService(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("create Google Drive client: %w", err)
	}
	return &GoogleDriveStorage{service: service, folderID: folderID}, nil
}

func (s *GoogleDriveStorage) child(ctx context.Context, parentID, name, mimeType string) (*drive.File, error) {
	query := fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", strings.ReplaceAll(name, "'", "\\'"), parentID)
	result, err := s.service.Files.List().Q(query).PageSize(10).Fields("files(id,name,mimeType,parents)").SupportsAllDrives(true).IncludeItemsFromAllDrives(true).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	for _, file := range result.Files {
		if file.MimeType == mimeType || mimeType == "" {
			return file, nil
		}
	}
	return nil, nil
}

func (s *GoogleDriveStorage) parent(ctx context.Context, name string, create bool) (string, error) {
	parentID := s.folderID
	parts := strings.Split(strings.Trim(name, "/"), "/")
	for _, part := range parts[:len(parts)-1] {
		folder, err := s.child(ctx, parentID, part, "application/vnd.google-apps.folder")
		if err != nil {
			return "", err
		}
		if folder == nil {
			if !create {
				return "", nil
			}
			folder, err = s.service.Files.Create(&drive.File{Name: part, MimeType: "application/vnd.google-apps.folder", Parents: []string{parentID}}).SupportsAllDrives(true).Context(ctx).Do()
			if err != nil {
				return "", err
			}
		}
		parentID = folder.Id
	}
	return parentID, nil
}

func (s *GoogleDriveStorage) findFile(ctx context.Context, name string) (*drive.File, error) {
	parentID, err := s.parent(ctx, name, false)
	if err != nil || parentID == "" {
		return nil, err
	}
	return s.child(ctx, parentID, path.Base(name), "")
}

func (s *GoogleDriveStorage) Upload(ctx context.Context, r io.Reader, name string) error {
	parentID, err := s.parent(ctx, name, true)
	if err != nil {
		return fmt.Errorf("resolve Google Drive parent: %w", err)
	}
	existing, err := s.child(ctx, parentID, path.Base(name), "")
	if err != nil {
		return fmt.Errorf("find Google Drive file: %w", err)
	}
	if existing != nil {
		_, err = s.service.Files.Update(existing.Id, &drive.File{}).Media(r).SupportsAllDrives(true).Context(ctx).Do()
	} else {
		_, err = s.service.Files.Create(&drive.File{Name: path.Base(name), Parents: []string{parentID}}).Media(r).SupportsAllDrives(true).Context(ctx).Do()
	}
	if err != nil {
		fmt.Println("GoogleDriveStorage:", name, err)
		return fmt.Errorf("upload Google Drive file %q: %w", name, err)
	}
	return nil
}

func (s *GoogleDriveStorage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	fmt.Println("GoogleDriveStorage:", name)
	file, err := s.findFile(ctx, name)
	if err != nil {
		fmt.Printf("find Google Drive file: %s\n", err)
		return nil, fmt.Errorf("find Google Drive file: %w", err)
	}
	if file == nil {
		fmt.Printf("Google Drive file %q not found\n", name)
		return nil, fmt.Errorf("Google Drive file %q not found", name)
	}
	result, err := s.service.Files.Get(file.Id).SupportsAllDrives(true).Context(ctx).Download()
	if err != nil {
		fmt.Printf("download Google Drive file %q: %s\n", name, err)
		return nil, fmt.Errorf("download Google Drive file %q: %w", name, err)
	}
	return result.Body, nil
}

func (s *GoogleDriveStorage) Delete(ctx context.Context, name string) error {
	file, err := s.findFile(ctx, name)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	if err := s.service.Files.Delete(file.Id).SupportsAllDrives(true).Context(ctx).Do(); err != nil {
		return fmt.Errorf("delete Google Drive file %q: %w", name, err)
	}
	return nil
}

func (s *GoogleDriveStorage) Exists(ctx context.Context, name string) (bool, error) {
	file, err := s.findFile(ctx, name)
	return file != nil, err
}

var _ Storage = (*GoogleDriveStorage)(nil)
