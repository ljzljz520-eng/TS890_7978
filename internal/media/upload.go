package media

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type UploadManager struct {
	Root     string
	MaxBytes int64
}

func NewUploadManager(root string, maxBytes int64) (UploadManager, error) {
	if root == "" {
		return UploadManager{}, errors.New("upload root is required")
	}
	if maxBytes <= 0 {
		return UploadManager{}, errors.New("upload limit must be positive")
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return UploadManager{}, err
	}
	return UploadManager{Root: root, MaxBytes: maxBytes}, nil
}

func (m UploadManager) Destination(clipID, sourceName string) (string, error) {
	if clipID == "" {
		return "", errors.New("clip id is required")
	}
	extension := strings.ToLower(filepath.Ext(sourceName))
	if extension != ".mp4" && extension != ".mov" && extension != ".webm" {
		return "", errors.New("unsupported extension")
	}
	return filepath.Join(m.Root, clipID+extension), nil
}

func (m UploadManager) Save(clipID, sourceName string, reader io.Reader) (string, string, int64, error) {
	destination, err := m.Destination(clipID, sourceName)
	if err != nil {
		return "", "", 0, err
	}
	temporary := destination + ".part"
	file, err := os.OpenFile(temporary, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return "", "", 0, err
	}
	hash := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(file, hash), io.LimitReader(reader, m.MaxBytes+1))
	closeErr := file.Close()
	if copyErr != nil {
		os.Remove(temporary)
		return "", "", 0, copyErr
	}
	if closeErr != nil {
		os.Remove(temporary)
		return "", "", 0, closeErr
	}
	if written > m.MaxBytes {
		os.Remove(temporary)
		return "", "", 0, errors.New("upload exceeds limit")
	}
	if written == 0 {
		os.Remove(temporary)
		return "", "", 0, errors.New("upload is empty")
	}
	if err := os.Rename(temporary, destination); err != nil {
		os.Remove(temporary)
		return "", "", 0, err
	}
	return destination, hex.EncodeToString(hash.Sum(nil)), written, nil
}

func (m UploadManager) Remove(path string) error {
	clean := filepath.Clean(path)
	root, err := filepath.Abs(m.Root)
	if err != nil {
		return err
	}
	absolute, err := filepath.Abs(clean)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return errors.New("path escapes upload root")
	}
	err = os.Remove(absolute)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func DetectMediaType(name string) string {
	value := mime.TypeByExtension(strings.ToLower(filepath.Ext(name)))
	if value == "" {
		return "application/octet-stream"
	}
	if index := strings.Index(value, ";"); index >= 0 {
		return value[:index]
	}
	return value
}

func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 0, errors.New("media path is a directory")
	}
	return info.Size(), nil
}
