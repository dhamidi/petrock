package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// SetupTempDirStep creates and sets up the temporary directory for testing
type SetupTempDirStep struct {
	baseDir string
}

// NewSetupTempDirStep creates a new setup temp dir step
func NewSetupTempDirStep(baseDir string) *SetupTempDirStep {
	return &SetupTempDirStep{
		baseDir: baseDir,
	}
}

// Name returns the step name
func (s *SetupTempDirStep) Name() string {
	return "Setup Temporary Directory"
}

// Execute creates the temporary directory and sets it in the context
func (s *SetupTempDirStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// 1. Ensure the ./tmp directory exists and create the temporary test directory within it
	if err := os.MkdirAll(s.baseDir, 0755); err != nil {
		return result.MarkFailure(fmt.Errorf("failed to create base temporary directory %s: %w", s.baseDir, err))
	}
	
	tempDir, err := os.MkdirTemp(s.baseDir, "petrock-test-*")
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to create temporary directory in %s: %w", s.baseDir, err))
	}
	
	slog.Info("Testing in temporary directory", "path", tempDir)
	result.AddLog("Created temporary directory: %s", tempDir)
	
	// Explicitly set permissions to 0755 as MkdirTemp defaults to 0700
	slog.Debug("Setting temporary directory permissions to 0755", "path", tempDir)
	if err := os.Chmod(tempDir, 0755); err != nil {
		// Attempt cleanup even if chmod fails
		_ = os.RemoveAll(tempDir)
		return result.MarkFailure(fmt.Errorf("failed to set permissions on temporary directory %s: %w", tempDir, err))
	}
	
	// Set the temp directory in the context
	ctx.SetTempDir(tempDir)
	
	// Add cleanup function
	ctx.AddCleanup(func() error {
		slog.Debug("Cleaning up temporary directory", "path", tempDir)
		if err := os.RemoveAll(tempDir); err != nil {
			slog.Error("Failed to remove temporary directory", "path", tempDir, "error", err)
			return err
		}
		return nil
	})
	
	return result.MarkSuccess()
}

// CreateProjectStep creates a new petrock project
type CreateProjectStep struct {
	projectName string
	modulePath  string
}

// NewCreateProjectStep creates a new create project step
func NewCreateProjectStep(projectName, modulePath string) *CreateProjectStep {
	return &CreateProjectStep{
		projectName: projectName,
		modulePath:  modulePath,
	}
}

// Name returns the step name
func (s *CreateProjectStep) Name() string {
	return "Create New Project"
}

// Execute creates a new petrock project
func (s *CreateProjectStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Store original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to get current working directory: %w", err))
	}
	
	// Add cleanup function to change back to original directory
	ctx.AddCleanup(func() error {
		if err := os.Chdir(originalWd); err != nil {
			slog.Error("Failed to change back to original directory", "path", originalWd, "error", err)
			return err
		}
		return nil
	})
	
	// Change into the temporary directory
	if err := os.Chdir(ctx.TempDir); err != nil {
		return result.MarkFailure(fmt.Errorf("failed to change directory to %s: %w", ctx.TempDir, err))
	}
	slog.Debug("Changed working directory", "path", ctx.TempDir)
	
	// Run 'petrock new' command
	slog.Info("Running 'petrock new'", "project", s.projectName, "module", s.modulePath)
	result.AddLog("Creating project: %s with module path: %s", s.projectName, s.modulePath)
	
	// Execute the 'new' command logic directly
	newArgs := []string{s.projectName, s.modulePath}
	if err := runNew(newCmd, newArgs); err != nil {
		return result.MarkFailure(fmt.Errorf("'petrock new' command failed during test: %w", err))
	}
	
	slog.Info("'petrock new' completed successfully")
	result.AddLog("Project created successfully")
	
	// Set project info in context
	ctx.ProjectName = s.projectName
	ctx.ModulePath = s.modulePath
	ctx.ProjectDir = filepath.Join(ctx.TempDir, s.projectName)
	
	return result.MarkSuccess()
}

// AddFeatureStep adds a feature to the project
type AddFeatureStep struct {
	featureName string
}

// NewAddFeatureStep creates a new add feature step
func NewAddFeatureStep(featureName string) *AddFeatureStep {
	return &AddFeatureStep{
		featureName: featureName,
	}
}

// Name returns the step name
func (s *AddFeatureStep) Name() string {
	return "Add Feature"
}

// Execute adds a feature to the project
func (s *AddFeatureStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Change into the newly created project directory
	if err := os.Chdir(ctx.ProjectName); err != nil {
		currentWd, _ := os.Getwd()
		return result.MarkFailure(fmt.Errorf("failed to change directory from %s to %s: %w", currentWd, ctx.ProjectName, err))
	}
	
	// Get the absolute path for logging clarity
	projectAbsDir, _ := filepath.Abs(ctx.ProjectName)
	slog.Debug("Changed working directory", "path", projectAbsDir)
	result.AddLog("Changed to project directory: %s", projectAbsDir)
	
	// Run 'petrock feature' command
	slog.Info("Running 'petrock feature'", "feature", s.featureName)
	result.AddLog("Adding feature: %s", s.featureName)
	
	featureArgs := []string{s.featureName}
	if err := runFeature(featureCmd, featureArgs); err != nil {
		return result.MarkFailure(fmt.Errorf("'petrock feature %s' command failed during test: %w", s.featureName, err))
	}
	
	slog.Info("'petrock feature' completed successfully")
	result.AddLog("Feature added successfully")
	
	return result.MarkSuccess()
}

// BuildProjectStep builds the project using go build
type BuildProjectStep struct{}

// NewBuildProjectStep creates a new build project step
func NewBuildProjectStep() *BuildProjectStep {
	return &BuildProjectStep{}
}

// Name returns the step name
func (s *BuildProjectStep) Name() string {
	return "Build Project"
}

// Execute builds the project
func (s *BuildProjectStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Run 'go build ./...' to ensure the project builds after adding the feature
	slog.Info("Running 'go build ./...' (after adding feature)")
	result.AddLog("Building project with 'go build ./...'")
	
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Stdout = os.Stdout // Pipe output to user
	buildCmd.Stderr = os.Stderr // Pipe errors to user
	// No need to set buildCmd.Dir, as we are already in the correct directory
	
	if err := buildCmd.Run(); err != nil {
		projectAbsDir, _ := filepath.Abs(".")
		return result.MarkFailure(fmt.Errorf("'go build ./...' failed in %s: %w", projectAbsDir, err))
	}
	
	result.AddLog("Project built successfully")
	
	// Also build the server binary for later use
	slog.Info("Building server binary for integration test")
	result.AddLog("Building server binary")
	
	buildServerCmd := exec.Command("go", "build", "-o", ctx.ProjectName+"-server", "./cmd/"+ctx.ProjectName)
	if err := buildServerCmd.Run(); err != nil {
		return result.MarkFailure(fmt.Errorf("failed to build server binary: %w", err))
	}
	
	result.AddLog("Server binary built successfully")
	return result.MarkSuccess()
}

// StartServerStep starts the web server for testing
type StartServerStep struct {
	port string
}

// NewStartServerStep creates a new start server step
func NewStartServerStep(port string) *StartServerStep {
	return &StartServerStep{
		port: port,
	}
}

// Name returns the step name
func (s *StartServerStep) Name() string {
	return "Start Web Server"
}

// Execute starts the web server
func (s *StartServerStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Start the server directly (no go run)
	slog.Info("Starting web server for integration test")
	result.AddLog("Starting server on port %s", s.port)
	
	serverCmd := exec.Command("./"+ctx.ProjectName+"-server", "serve", "--port", s.port)
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr
	
	if err := serverCmd.Start(); err != nil {
		return result.MarkFailure(fmt.Errorf("failed to start web server: %w", err))
	}
	
	// Store the server process in context
	ctx.ServerCmd = serverCmd
	ctx.ServerPort = s.port
	
	// Add cleanup function to terminate server
	ctx.AddCleanup(func() error {
		if ctx.ServerCmd != nil && ctx.ServerCmd.Process != nil {
			slog.Info("Terminating test web server")
			if err := ctx.ServerCmd.Process.Kill(); err != nil {
				slog.Error("Failed to kill web server process", "error", err)
				return err
			}
			// Wait for the process to exit
			_ = ctx.ServerCmd.Wait()
		}
		return nil
	})
	
	// Wait a moment for the server to initialize
	slog.Info("Waiting for server to initialize...")
	result.AddLog("Waiting for server to initialize...")
	time.Sleep(2 * time.Second)
	
	result.AddLog("Server started successfully")
	return result.MarkSuccess()
}

// StopServerStep stops the web server
type StopServerStep struct{}

// NewStopServerStep creates a new stop server step
func NewStopServerStep() *StopServerStep {
	return &StopServerStep{}
}

// Name returns the step name
func (s *StopServerStep) Name() string {
	return "Stop Web Server"
}

// Execute stops the web server
func (s *StopServerStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	if ctx.ServerCmd == nil || ctx.ServerCmd.Process == nil {
		return result.MarkFailure(fmt.Errorf("no server process to stop"))
	}
	
	slog.Info("Stopping web server")
	result.AddLog("Stopping web server")
	
	if err := ctx.ServerCmd.Process.Kill(); err != nil {
		return result.MarkFailure(fmt.Errorf("failed to kill web server process: %w", err))
	}
	
	// Wait for the process to exit
	_ = ctx.ServerCmd.Wait()
	ctx.ServerCmd = nil
	
	result.AddLog("Server stopped successfully")
	return result.MarkSuccess()
}

// HTTPGetStep performs an HTTP GET request
type HTTPGetStep struct {
	url            string
	expectedStatus int
}

// NewHTTPGetStep creates a new HTTP GET step
func NewHTTPGetStep(url string, expectedStatus int) *HTTPGetStep {
	return &HTTPGetStep{
		url:            url,
		expectedStatus: expectedStatus,
	}
}

// Name returns the step name
func (s *HTTPGetStep) Name() string {
	return "HTTP GET Request"
}

// Execute performs the HTTP GET request
func (s *HTTPGetStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Make an HTTP request
	slog.Info("Testing HTTP endpoint", "url", s.url)
	result.AddLog("Making HTTP GET request to: %s", s.url)
	
	resp, err := http.Get(s.url)
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to make HTTP request: %w", err))
	}
	defer resp.Body.Close()
	
	// Verify the response status
	if resp.StatusCode != s.expectedStatus {
		return result.MarkFailure(fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, s.expectedStatus))
	}
	
	slog.Info("HTTP endpoint test successful", "status", resp.Status)
	result.AddLog("HTTP GET request successful, status: %s", resp.Status)
	
	return result.MarkSuccess()
}

// HTTPPostStep performs an HTTP POST request
type HTTPPostStep struct {
	url            string
	formData       url.Values
	expectedStatus int
	shouldFail     bool
}

// NewHTTPPostStep creates a new HTTP POST step
func NewHTTPPostStep(url string, formData url.Values, expectedStatus int, shouldFail bool) *HTTPPostStep {
	return &HTTPPostStep{
		url:            url,
		formData:       formData,
		expectedStatus: expectedStatus,
		shouldFail:     shouldFail,
	}
}

// Name returns the step name
func (s *HTTPPostStep) Name() string {
	return "HTTP POST Request"
}

// Execute performs the HTTP POST request
func (s *HTTPPostStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	if s.shouldFail {
		slog.Info("Testing invalid POST request with empty fields")
		result.AddLog("Testing POST request expected to fail with status %d", s.expectedStatus)
	} else {
		slog.Info("Testing valid POST request")
		result.AddLog("Testing POST request expected to succeed with status %d", s.expectedStatus)
	}
	
	postResp, err := http.PostForm(s.url, s.formData)
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to make POST request: %w", err))
	}
	defer postResp.Body.Close()
	
	// Verify the response status
	if postResp.StatusCode != s.expectedStatus {
		return result.MarkFailure(fmt.Errorf("unexpected status code: got %d, want %d", postResp.StatusCode, s.expectedStatus))
	}
	
	if s.shouldFail {
		slog.Info("Invalid POST request test successful", "status", postResp.Status)
		result.AddLog("POST request correctly failed with status: %s", postResp.Status)
	} else {
		slog.Info("Valid POST request test successful", "status", postResp.Status)
		result.AddLog("POST request successful with status: %s", postResp.Status)
	}
	
	return result.MarkSuccess()
}

// CommandAPIStep tests the command API endpoint
type CommandAPIStep struct {
	url         string
	payload     map[string]interface{}
	expectedStatus int
}

// NewCommandAPIStep creates a new command API step
func NewCommandAPIStep(url string, payload map[string]interface{}, expectedStatus int) *CommandAPIStep {
	return &CommandAPIStep{
		url:         url,
		payload:     payload,
		expectedStatus: expectedStatus,
	}
}

// Name returns the step name
func (s *CommandAPIStep) Name() string {
	return "Command API Request"
}

// Execute performs the command API request
func (s *CommandAPIStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	slog.Info("Testing command API request with validation error")
	result.AddLog("Testing command API request expected to return status %d", s.expectedStatus)
	
	jsonData, err := json.Marshal(s.payload)
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to marshal command payload: %w", err))
	}
	
	cmdResp, err := http.Post(s.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to make command API request: %w", err))
	}
	defer cmdResp.Body.Close()
	
	// Should get the expected status code
	if cmdResp.StatusCode != s.expectedStatus {
		return result.MarkFailure(fmt.Errorf("unexpected status for command validation error: got %d, want %d", cmdResp.StatusCode, s.expectedStatus))
	}
	
	// Verify the response is JSON with validation error details (if expecting 400)
	if s.expectedStatus == http.StatusBadRequest {
		var cmdErrorResp map[string]interface{}
		if err := json.NewDecoder(cmdResp.Body).Decode(&cmdErrorResp); err != nil {
			return result.MarkFailure(fmt.Errorf("failed to decode command error response: %w", err))
		}
		
		if cmdErrorResp["error"] != "Validation failed" {
			return result.MarkFailure(fmt.Errorf("unexpected error message: got %v, want 'Validation failed'", cmdErrorResp["error"]))
		}
	}
	
	slog.Info("Command API request test successful", "status", cmdResp.Status)
	result.AddLog("Command API request successful with status: %s", cmdResp.Status)
	
	return result.MarkSuccess()
}

// SelfInspectStep tests the self inspect command
type SelfInspectStep struct{}

// NewSelfInspectStep creates a new self inspect step
func NewSelfInspectStep() *SelfInspectStep {
	return &SelfInspectStep{}
}

// Name returns the step name
func (s *SelfInspectStep) Name() string {
	return "Self Inspect Command"
}

// Execute runs the self inspect command and validates output
func (s *SelfInspectStep) Execute(ctx *TestContext) *StepResult {
	result := NewStepResult(s.Name())
	
	// Test the self inspect command
	slog.Info("Testing 'self inspect' command")
	result.AddLog("Running 'go run ./cmd/%s self inspect'", ctx.ProjectName)
	
	selfInspectCmd := exec.Command("go", "run", "./cmd/"+ctx.ProjectName, "self", "inspect")
	
	// Capture the command output
	selfInspectOutput, err := selfInspectCmd.Output()
	if err != nil {
		return result.MarkFailure(fmt.Errorf("failed to run 'self inspect' command: %w", err))
	}
	
	// Verify the output is valid JSON
	slog.Info("Verifying 'self inspect' output is valid JSON")
	result.AddLog("Validating JSON output from self inspect command")
	
	if err := validateInspectJSON(selfInspectOutput); err != nil {
		return result.MarkFailure(err)
	}
	
	slog.Info("'self inspect' command test successful")
	result.AddLog("Self inspect command validation successful")
	
	return result.MarkSuccess()
}

// validateInspectJSON validates the JSON output from self inspect
func validateInspectJSON(jsonData []byte) error {
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return fmt.Errorf("'self inspect' command did not produce valid JSON: %w", err)
	}
	
	// Verify the JSON contains the expected keys
	expectedKeys := []string{"commands", "queries", "routes", "features", "workers"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			return fmt.Errorf("'self inspect' output missing expected key: %s", key)
		}
	}
	
	// Verify workers are included in the output
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
	
	return nil
}
