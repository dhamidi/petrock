package utils

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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

// CopyDir recursively copies a directory from src to dest.
// It specifically handles renaming a subdirectory named 'cmd/dirPlaceholder'
// under src to 'cmd/dirReplacement' under dest.
func CopyDir(src, dest, dirPlaceholder, dirReplacement string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory %s: %w", src, err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source %s is not a directory", src)
	}

	// Ensure the base destination directory exists
	if err := EnsureDir(dest); err != nil {
		return fmt.Errorf("failed to ensure destination directory %s: %w", dest, err)
	}

	placeholderCmdDir := filepath.Join("cmd", dirPlaceholder)
	replacementCmdDir := filepath.Join("cmd", dirReplacement)

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s during walk: %w", path, err)
		}

		// Calculate the relative path from the source root
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		// Determine the target path, applying the rename logic
		targetPath := filepath.Join(dest, relPath)
		if strings.HasPrefix(relPath, placeholderCmdDir) {
			// Check if it's the directory itself or a file/subdir within it
			if relPath == placeholderCmdDir || strings.HasPrefix(relPath, placeholderCmdDir+string(filepath.Separator)) {
				newRelPath := strings.Replace(relPath, placeholderCmdDir, replacementCmdDir, 1)
				targetPath = filepath.Join(dest, newRelPath)
			}
		}

		if d.IsDir() {
			// Create the directory in the destination
			// Skip the source root itself as EnsureDir(dest) already created it
			if path != src {
				slog.Debug("Creating directory", "path", targetPath)
				if err := EnsureDir(targetPath); err != nil {
					return fmt.Errorf("failed to create target directory %s: %w", targetPath, err)
				}
			}
		} else {
			// Copy the file
			slog.Debug("Copying file", "from", path, "to", targetPath)
			if err := copyFile(path, targetPath); err != nil {
				return fmt.Errorf("failed to copy file from %s to %s: %w", path, targetPath, err)
			}
		}
		return nil
	})
}

// copyFile copies a single file from src to dest.
func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Preserve permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcInfo.Mode())
}

// ReplaceInFiles walks through all files in rootDir and replaces occurrences
// of keys in the replacements map with their corresponding values.
func ReplaceInFiles(rootDir string, replacements map[string]string) error {
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s during replacement walk: %w", path, err)
		}

		// Skip directories and non-regular files
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		slog.Debug("Processing file for replacements", "path", path)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s for replacement: %w", path, err)
		}

		originalContent := string(content)
		modifiedContent := originalContent

		for placeholder, value := range replacements {
			modifiedContent = strings.ReplaceAll(modifiedContent, placeholder, value)
		}

		// Only write back if content actually changed
		if modifiedContent != originalContent {
			slog.Debug("Writing modified content", "path", path)
			// Get original file permissions
			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to get file info for %s: %w", path, err)
			}
			// Write back with original permissions
			err = os.WriteFile(path, []byte(modifiedContent), info.Mode())
			if err != nil {
				return fmt.Errorf("failed to write modified file %s: %w", path, err)
			}
		}

		return nil
	})
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
