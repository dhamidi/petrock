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

	// Explicitly set permissions to 0755 as MkdirTemp defaults to 0700
	slog.Debug("Setting temporary directory permissions to 0755", "path", tempDir) // Keep this debug log for now
	if err := os.Chmod(tempDir, 0755); err != nil {
		// Attempt cleanup even if chmod fails
		_ = os.RemoveAll(tempDir)
		return fmt.Errorf("failed to set permissions on temporary directory %s: %w", tempDir, err)
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

	// 4. Change into the newly created project directory (`./selftest`)
	// Since the current working directory is tempDir, we just need to chdir into projectName.
	if err := os.Chdir(projectName); err != nil {
		// Get current WD for better error message
		currentWd, _ := os.Getwd()
		return fmt.Errorf("failed to change directory from %s to %s: %w", currentWd, projectName, err)
	}
	// Get the absolute path for logging clarity
	projectAbsDir, _ := filepath.Abs(projectName)
	slog.Debug("Changed working directory", "path", projectAbsDir)


	// 5. Run `go build ./...`
	slog.Info("Running 'go build ./...'")
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Stdout = os.Stdout // Pipe output to user
	buildCmd.Stderr = os.Stderr // Pipe errors to user
	// No need to set buildCmd.Dir, as we are already in the correct directory

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("'go build ./...' failed in %s: %w", projectAbsDir, err) // Use projectAbsDir here
	}

	slog.Info("Self-test completed successfully!")
	fmt.Println("\nSuccess! The generated project builds correctly.")
	return nil
}
