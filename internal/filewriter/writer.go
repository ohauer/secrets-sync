package filewriter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// FileConfig holds file writing configuration
type FileConfig struct {
	Path  string
	Mode  os.FileMode
	Owner int
	Group int
}

// Writer handles atomic file writing
type Writer struct{}

// NewWriter creates a new file writer
func NewWriter() *Writer {
	return &Writer{}
}

// WriteFile writes content to a file atomically
func (w *Writer) WriteFile(config FileConfig, content string) error {
	if err := w.ensureDir(filepath.Dir(config.Path)); err != nil {
		return err
	}

	tmpFile := config.Path + ".tmp"

	if err := os.WriteFile(tmpFile, []byte(content), config.Mode); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if config.Owner >= 0 || config.Group >= 0 {
		uid := config.Owner
		gid := config.Group
		if uid < 0 {
			uid = -1
		}
		if gid < 0 {
			gid = -1
		}
		if err := os.Chown(tmpFile, uid, gid); err != nil {
			_ = os.Remove(tmpFile)
			return fmt.Errorf("failed to set ownership: %w", err)
		}
	}

	if err := os.Rename(tmpFile, config.Path); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (w *Writer) ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ParseMode parses a mode string (e.g., "0600") to os.FileMode
func ParseMode(mode string) (os.FileMode, error) {
	if mode == "" {
		return 0600, nil
	}

	m, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid mode: %w", err)
	}

	return os.FileMode(m), nil
}

// ParseOwner parses owner/group string to int
func ParseOwner(owner string) (int, error) {
	if owner == "" {
		return -1, nil
	}

	o, err := strconv.Atoi(owner)
	if err != nil {
		return -1, fmt.Errorf("invalid owner: %w", err)
	}

	return o, nil
}

// GetFileInfo returns file information
func GetFileInfo(path string) (os.FileMode, int, int, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, -1, -1, err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.Mode(), -1, -1, nil
	}

	return info.Mode(), int(stat.Uid), int(stat.Gid), nil
}
