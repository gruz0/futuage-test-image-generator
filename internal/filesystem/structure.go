package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirectoryStructure creates the required output directory structure
func EnsureDirectoryStructure(baseDir string) error {
	// Create main directories
	dirs := []string{
		filepath.Join(baseDir, "ratios"),
		filepath.Join(baseDir, "targets"),
		filepath.Join(baseDir, "edge-cases"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// CleanDirectory removes all contents of a directory but keeps the directory
func CleanDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, nothing to clean
		}
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
