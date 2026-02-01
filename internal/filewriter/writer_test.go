package filewriter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	writer := NewWriter()
	config := FileConfig{
		Path:  filePath,
		Mode:  0644,
		Owner: -1,
		Group: -1,
	}

	content := "test content"
	if err := writer.WriteFile(config, content); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected '%s', got '%s'", content, string(data))
	}
}

func TestWriteFile_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "test.txt")

	writer := NewWriter()
	config := FileConfig{
		Path:  filePath,
		Mode:  0600,
		Owner: -1,
		Group: -1,
	}

	if err := writer.WriteFile(config, "content"); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("file was not created")
	}
}

func TestWriteFile_Permissions(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	writer := NewWriter()
	config := FileConfig{
		Path:  filePath,
		Mode:  0600,
		Owner: -1,
		Group: -1,
	}

	if err := writer.WriteFile(config, "content"); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %o", info.Mode().Perm())
	}
}

func TestWriteFile_Atomic(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	_ = os.WriteFile(filePath, []byte("original"), 0644)

	writer := NewWriter()
	config := FileConfig{
		Path:  filePath,
		Mode:  0644,
		Owner: -1,
		Group: -1,
	}

	if err := writer.WriteFile(config, "updated"); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != "updated" {
		t.Errorf("expected 'updated', got '%s'", string(data))
	}

	tmpFile := filePath + ".tmp"
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("temp file was not cleaned up")
	}
}

func TestParseMode_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected os.FileMode
	}{
		{"0600", 0600},
		{"0644", 0644},
		{"0755", 0755},
		{"", 0600},
	}

	for _, tt := range tests {
		mode, err := ParseMode(tt.input)
		if err != nil {
			t.Errorf("ParseMode(%s) failed: %v", tt.input, err)
		}
		if mode != tt.expected {
			t.Errorf("ParseMode(%s) = %o, expected %o", tt.input, mode, tt.expected)
		}
	}
}

func TestParseMode_Invalid(t *testing.T) {
	_, err := ParseMode("invalid")
	if err == nil {
		t.Error("expected error for invalid mode, got nil")
	}
}

func TestParseOwner_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1000", 1000},
		{"0", 0},
		{"", -1},
	}

	for _, tt := range tests {
		owner, err := ParseOwner(tt.input)
		if err != nil {
			t.Errorf("ParseOwner(%s) failed: %v", tt.input, err)
		}
		if owner != tt.expected {
			t.Errorf("ParseOwner(%s) = %d, expected %d", tt.input, owner, tt.expected)
		}
	}
}

func TestParseOwner_Invalid(t *testing.T) {
	_, err := ParseOwner("invalid")
	if err == nil {
		t.Error("expected error for invalid owner, got nil")
	}
}

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	_ = os.WriteFile(filePath, []byte("content"), 0644)

	mode, uid, gid, err := GetFileInfo(filePath)
	if err != nil {
		t.Fatalf("failed to get file info: %v", err)
	}

	if mode.Perm() != 0644 {
		t.Errorf("expected mode 0644, got %o", mode.Perm())
	}

	if uid < 0 {
		t.Error("expected valid uid")
	}

	if gid < 0 {
		t.Error("expected valid gid")
	}
}

func TestGetFileInfo_NonExistent(t *testing.T) {
	_, _, _, err := GetFileInfo("/nonexistent/file")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
