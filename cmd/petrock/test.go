package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

	// 4a. Run `petrock feature posts` inside the new project
	featureName := "posts"
	slog.Info("Running 'petrock feature'", "feature", featureName)
	featureArgs := []string{featureName}
	// Ensure featureCmd is accessible (it should be if defined in the same package)
	if err := runFeature(featureCmd, featureArgs); err != nil {
		return fmt.Errorf("'petrock feature %s' command failed during test: %w", featureName, err)
	}
	slog.Info("'petrock feature' completed successfully")


	// 5. Run `go build ./...` to ensure the project still builds after adding the feature
	slog.Info("Running 'go build ./...' (after adding feature)")
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Stdout = os.Stdout // Pipe output to user
	buildCmd.Stderr = os.Stderr // Pipe errors to user
	// No need to set buildCmd.Dir, as we are already in the correct directory

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("'go build ./...' failed in %s: %w", projectAbsDir, err) // Use projectAbsDir here
	}

	// 6. Build the server binary first
	slog.Info("Building server binary for integration test")
	buildServerCmd := exec.Command("go", "build", "-o", "selftest-server", "./cmd/selftest")
	if err := buildServerCmd.Run(); err != nil {
		return fmt.Errorf("failed to build server binary: %w", err)
	}

	// 7. Start the server directly (no go run)
	slog.Info("Starting web server for integration test")
	serverCmd := exec.Command("./selftest-server", "serve", "--port", "8081")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr
	
	if err := serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	// Ensure server is terminated when we're done
	defer func() {
		slog.Info("Terminating test web server")
		if err := serverCmd.Process.Kill(); err != nil {
			slog.Error("Failed to kill web server process", "error", err)
		}
		// Wait for the process to exit
		_ = serverCmd.Wait()
	}()

	// 8. Wait a moment for the server to initialize
	slog.Info("Waiting for server to initialize...")
	time.Sleep(2 * time.Second)

	// 9. Make an HTTP request to /posts
	slog.Info("Testing HTTP endpoint", "url", "http://localhost:8081/posts")
	resp, err := http.Get("http://localhost:8081/posts")
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 10. Verify the response is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	slog.Info("HTTP endpoint test successful", "status", resp.Status)

	// 11. Test the self inspect command
	slog.Info("Testing 'self inspect' command")
	selfInspectCmd := exec.Command("go", "run", "./cmd/selftest", "self", "inspect")
	
	// Capture the command output
	selfInspectOutput, err := selfInspectCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run 'self inspect' command: %w", err)
	}

	// 12. Verify the output is valid JSON
	slog.Info("Verifying 'self inspect' output is valid JSON")
	var result map[string]interface{}
	if err := json.Unmarshal(selfInspectOutput, &result); err != nil {
		return fmt.Errorf("'self inspect' command did not produce valid JSON: %w", err)
	}

	// 13. Verify the JSON contains the expected keys
	expectedKeys := []string{"commands", "queries", "routes", "features", "workers"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			return fmt.Errorf("'self inspect' output missing expected key: %s", key)
		}
	}

	// 14. Verify workers are included in the output
	workers, ok := result["workers"].([]interface{})
	if !ok {
		return fmt.Errorf("'workers' key doesn't contain an array of workers")
	}
	
	// Check that we have at least one worker
	if len(workers) == 0 {
		return fmt.Errorf("no workers found in self-inspect output")
	}
	
	// Check the first worker has the expected fields
	worker, ok := workers[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("worker is not properly structured")
	}
	
	// Verify worker fields
	workerFields := []string{"name", "type", "methods"}
	for _, field := range workerFields {
		if _, ok := worker[field]; !ok {
			return fmt.Errorf("worker missing expected field: %s", field)
		}
	}
	
	slog.Info("'self inspect' command test successful")

	slog.Info("Self-test completed successfully!")
	fmt.Println("\nSuccess! The generated project builds correctly and serves content.")
	return nil
}
