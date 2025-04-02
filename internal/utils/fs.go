package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir creates a directory if it doesn't exist.
// It's similar to `mkdir -p`.
func EnsureDir(path string) error {
	err := os.MkdirAll(path, 0755) // 0755 is standard permission for directories
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// WriteFile writes byte content to a file, ensuring the parent directory exists.
func WriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to ensure directory %s for file %s: %w", dir, path, err)
	}

	err := os.WriteFile(path, content, 0644) // 0644 is standard permission for files
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}
