package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDirectoryStructure(t *testing.T) {
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "output")

	err := EnsureDirectoryStructure(baseDir)
	if err != nil {
		t.Fatalf("EnsureDirectoryStructure() error = %v", err)
	}

	// Verify all expected directories were created
	expectedDirs := []string{
		filepath.Join(baseDir, "ratios"),
		filepath.Join(baseDir, "targets"),
		filepath.Join(baseDir, "edge-cases"),
	}

	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("Expected directory %q not found: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%q is not a directory", dir)
		}
	}
}

func TestEnsureDirectoryStructure_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "output")

	// Call twice - should not error on second call
	err := EnsureDirectoryStructure(baseDir)
	if err != nil {
		t.Fatalf("First EnsureDirectoryStructure() error = %v", err)
	}

	err = EnsureDirectoryStructure(baseDir)
	if err != nil {
		t.Fatalf("Second EnsureDirectoryStructure() error = %v", err)
	}
}

func TestCleanDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some files and subdirectories
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	subFile := filepath.Join(subDir, "nested.txt")
	if err := os.WriteFile(subFile, []byte("nested"), 0644); err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	// Clean the directory
	err := CleanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("CleanDirectory() error = %v", err)
	}

	// Verify directory still exists but is empty
	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Directory was deleted instead of cleaned: %v", err)
	}
	if !info.IsDir() {
		t.Error("Path is no longer a directory")
	}

	// Verify contents were removed
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Directory should be empty, has %d entries", len(entries))
	}
}

func TestCleanDirectory_NonExistent(t *testing.T) {
	// Should not error for non-existent directory
	err := CleanDirectory("/nonexistent/path/that/does/not/exist")
	if err != nil {
		t.Errorf("CleanDirectory() for non-existent path should not error, got: %v", err)
	}
}

func TestGetFileSize(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with known content
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("hello world") // 11 bytes
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	size, err := GetFileSize(testFile)
	if err != nil {
		t.Fatalf("GetFileSize() error = %v", err)
	}

	if size != int64(len(content)) {
		t.Errorf("GetFileSize() = %d, want %d", size, len(content))
	}
}

func TestGetFileSize_NonExistent(t *testing.T) {
	_, err := GetFileSize("/nonexistent/file.txt")
	if err == nil {
		t.Error("GetFileSize() for non-existent file should error")
	}
}

