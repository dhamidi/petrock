package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall" // Import for umask

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs a self-test by creating and building a new project",
	Long: `Creates a temporary directory, runs 'petrock new selftest ...' within it,
and then attempts to build the newly generated project using 'go build ./...'.
This helps verify that the templates and the 'new' command are working correctly.`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	// Log the current umask
	originalUmask := syscall.Umask(0) // Get current umask
	syscall.Umask(originalUmask)      // Set it back immediately
	slog.Info("Current umask", "mask", fmt.Sprintf("%04o", originalUmask))

	// 1. Ensure the ./tmp directory exists and create the temporary test directory within it
	tmpBaseDir := "./tmp"
	if err := os.MkdirAll(tmpBaseDir, 0755); err != nil { // Use explicit 0755 permission
		return fmt.Errorf("failed to create base temporary directory %s: %w", tmpBaseDir, err)
	}
	tempDir, err := os.MkdirTemp(tmpBaseDir, "petrock-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory in %s: %w", tmpBaseDir, err)
	}
	slog.Info("Testing in temporary directory", "path", tempDir)

	// Log permissions of the created temp directory
	fileInfo, err := os.Stat(tempDir)
	if err != nil {
		slog.Warn("Could not stat temporary directory to check permissions", "path", tempDir, "error", err)
	} else {
		slog.Info("Temporary directory permissions", "path", tempDir, "permissions", fileInfo.Mode().String())
	}

	// Explicitly set permissions to 0755 as MkdirTemp defaults to 0700
	slog.Debug("Setting temporary directory permissions to 0755", "path", tempDir)
	if err := os.Chmod(tempDir, 0755); err != nil {
		// Attempt cleanup even if chmod fails
		_ = os.RemoveAll(tempDir)
		return fmt.Errorf("failed to set permissions on temporary directory %s: %w", tempDir, err)
	}

	// Verify permissions *after* chmod
	postChmodInfo, err := os.Stat(tempDir)
	if err != nil {
		slog.Warn("Could not stat temporary directory after chmod", "path", tempDir, "error", err)
	} else {
		slog.Info("Temporary directory permissions *after* chmod", "path", tempDir, "permissions", postChmodInfo.Mode().String())
	}


	// Ensure the temporary directory is cleaned up afterwards
	defer func() {
		slog.Debug("Cleaning up temporary directory", "path", tempDir)
		if err := os.RemoveAll(tempDir); err != nil {
			slog.Error("Failed to remove temporary directory", "path", tempDir, "error", err)
		}
	}()

	// Store original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	// Change back to original directory at the end
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			slog.Error("Failed to change back to original directory", "path", originalWd, "error", err)
		}
	}()

	// 2. Change into that directory
	if err := os.Chdir(tempDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", tempDir, err)
	}
	slog.Debug("Changed working directory", "path", tempDir)

	// 3. Run `petrock new selftest github.com/petrock/selftest`
	projectName := "selftest"
	modulePath := "github.com/petrock/selftest"
	slog.Info("Running 'petrock new'", "project", projectName, "module", modulePath)

	// We can execute the 'new' command's logic directly.
	// Ensure 'newCmd' is accessible (it should be if defined in the same package).
	// We need to simulate the command-line arguments for runNew.
	newArgs := []string{projectName, modulePath}
	if err := runNew(newCmd, newArgs); err != nil { // Assuming newCmd is accessible
		return fmt.Errorf("'petrock new' command failed during test: %w", err)
	}
	slog.Info("'petrock new' completed successfully")

	// 4. Change into `./selftest`
	projectDir := filepath.Join(tempDir, projectName)
	if err := os.Chdir(projectDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", projectDir, err)
	}
	slog.Debug("Changed working directory", "path", projectDir)

	// 5. Run `go build ./...`
	slog.Info("Running 'go build ./...'")
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Stdout = os.Stdout // Pipe output to user
	buildCmd.Stderr = os.Stderr // Pipe errors to user
	buildCmd.Dir = projectDir   // Ensure command runs in the project directory

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("'go build ./...' failed in %s: %w", projectDir, err)
	}

	slog.Info("Self-test completed successfully!")
	fmt.Println("\nSuccess! The generated project builds correctly.")
	return nil
}
