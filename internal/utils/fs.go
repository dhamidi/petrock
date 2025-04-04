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

// CopyDir recursively copies content from an fs.FS (like embed.FS) starting at srcRoot
// to a destination directory on the filesystem (dest).
// It specifically handles renaming a subdirectory named 'cmd/dirPlaceholder'
// under srcRoot to 'cmd/dirReplacement' under dest.
func CopyDir(fsys fs.FS, srcRoot, dest, dirPlaceholder, dirReplacement string) error {
	// Ensure the base destination directory exists
	if err := EnsureDir(dest); err != nil {
		return fmt.Errorf("failed to ensure destination directory %s: %w", dest, err)
	}

	// Use forward slash for FS paths and placeholder matching
	placeholderCmdDir := "cmd/" + dirPlaceholder
	replacementCmdDir := filepath.Join("cmd", dirReplacement) // Use OS separator for dest path

	return fs.WalkDir(fsys, srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Handle potential errors during walk (e.g., permission issues if reading from OS FS)
			return fmt.Errorf("error accessing path %q within source FS: %w", path, err)
		}

		// srcRoot is the starting point in fsys, path is relative to fsys root.
		// We need the path relative to srcRoot for constructing the destination path.
		// If srcRoot is ".", path is already relative.
		relPath := path
		if srcRoot != "." {
			// This check might be needed if srcRoot is not "."
			if !strings.HasPrefix(path, srcRoot) {
				// Should not happen if WalkDir is used correctly, but good sanity check
				return fmt.Errorf("walked path %q does not start with srcRoot %q", path, srcRoot)
			}
			relPath, err = filepath.Rel(srcRoot, path) // Use filepath for potential OS differences if srcRoot had separators
			if err != nil {
				return fmt.Errorf("failed to get relative path for %q from %q: %w", path, srcRoot, err)
			}
		}

		// Convert relPath to use OS-specific separators for joining with dest
		osRelPath := filepath.FromSlash(relPath)

		// Determine the target path, applying the rename logic
		targetPath := filepath.Join(dest, osRelPath)
		// Use forward slash path for matching placeholder prefix
		if strings.HasPrefix(path, placeholderCmdDir) {
			// Check if it's the directory itself or a file/subdir within it
			if path == placeholderCmdDir || strings.HasPrefix(path, placeholderCmdDir+"/") {
				// Replace using forward slash path, then construct OS-specific target path
				newRelPath := strings.Replace(path, placeholderCmdDir, replacementCmdDir, 1)
				targetPath = filepath.Join(dest, filepath.FromSlash(newRelPath)) // Ensure OS-specific separator
			}
		}

		if d.IsDir() {
			// Create the directory in the destination
			// Skip the source root itself as EnsureDir(dest) already created it
			if path != srcRoot {
				slog.Debug("Creating directory", "path", targetPath)
				// Use original directory permissions from embedded FS if possible
				info, statErr := fs.Stat(fsys, path)
				mode := os.ModeDir | 0755 // Default directory mode
				if statErr == nil {
					mode = info.Mode() | os.ModeDir // Ensure it's marked as dir
				} else {
					slog.Warn("Could not stat source directory in FS, using default permissions", "path", path, "error", statErr)
				}
				if err := os.MkdirAll(targetPath, mode); err != nil { // Use MkdirAll with mode
					return fmt.Errorf("failed to create target directory %s: %w", targetPath, err)
				}
			}
		} else {
			// Copy the file
			slog.Debug("Copying file", "from", path, "to", targetPath)
			if err := copyFileFromFS(fsys, path, targetPath); err != nil {
				return fmt.Errorf("failed to copy file from FS path %q to %q: %w", path, targetPath, err)
			}
		}
		return nil
	})
}

// copyFileFromFS copies a single file from an fs.FS to a destination path.
func copyFileFromFS(fsys fs.FS, srcPath, destPath string) error {
	sourceFile, err := fsys.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %q in FS: %w", srcPath, err)
	}
	defer sourceFile.Close()

	// Get source file info for permissions
	srcInfo, err := fs.Stat(fsys, srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat source file %q in FS: %w", srcPath, err)
	}
	mode := srcInfo.Mode()

	// Ensure destination directory exists before creating file
	destDir := filepath.Dir(destPath)
	if err := EnsureDir(destDir); err != nil {
		return fmt.Errorf("failed to ensure destination directory %q for file %q: %w", destDir, destPath, err)
	}


	// Create destination file with source permissions
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("failed to create destination file %q: %w", destPath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content for %q: %w", destPath, err)
	}

	// Chmod might be redundant if OpenFile worked correctly, but can be a safeguard
	// if umask affected the creation mode.
	// return os.Chmod(destPath, mode)
	return nil // Permissions set by OpenFile
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
