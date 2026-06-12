package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage abstracts where attachment bytes live. The application depends only on
// this interface, so swapping local disk for S3/GCS later requires no changes
// outside this package.
type Storage interface {
	Save(storedName string, src io.Reader) error
	Open(storedName string) (io.ReadCloser, error)
	Delete(storedName string) error
}

// LocalStorage persists files on the local filesystem under BaseDir.
type LocalStorage struct {
	BaseDir string
}

// NewLocalStorage creates the base directory if needed and returns a LocalStorage.
func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0o750); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &LocalStorage{BaseDir: baseDir}, nil
}

func (s *LocalStorage) path(storedName string) string {
	// filepath.Base strips any path components, preventing directory traversal
	// via a crafted stored name.
	return filepath.Join(s.BaseDir, filepath.Base(storedName))
}

// Save writes the reader's contents to a new file.
func (s *LocalStorage) Save(storedName string, src io.Reader) error {
	dst, err := os.Create(s.path(storedName))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

// Open returns a reader for a stored file.
func (s *LocalStorage) Open(storedName string) (io.ReadCloser, error) {
	return os.Open(s.path(storedName))
}

// Delete removes a stored file. A missing file is not treated as an error.
func (s *LocalStorage) Delete(storedName string) error {
	err := os.Remove(s.path(storedName))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
