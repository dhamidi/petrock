package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/dhamidi/petrock/internal/ui"
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
	// Create test context and runner
	ctx := NewTestContext()
	runner := NewTestRunner(ctx)
	
	// Setup all test steps
	setupTestSteps(runner)
	
	// Execute all steps
	testCtx := context.Background()
	testSuccess := true
	var testErr error
	
	if err := runner.RunAllSteps(testCtx); err != nil {
		testSuccess = false
		testErr = err
		reportFinalResults(runner, false)
	} else {
		reportFinalResults(runner, true)
		cmdCtx.UI.ShowSuccess(cmdCtx.Ctx, "\nSuccess! The generated project builds correctly and serves content.\n")
	}
	
	// Handle cleanup based on test results
	if testSuccess {
		// Clean up on success
		ctx.RunCleanup()
	} else {
		// On failure, skip cleanup and show project location for debugging
		cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "\n🔍 Test failed - project preserved for debugging:")
		if ctx.TempDir != "" {
			cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "   Base directory: %s\n", ctx.TempDir)
		}
		if ctx.ProjectName != "" && ctx.TempDir != "" {
			projectPath := filepath.Join(ctx.TempDir, ctx.ProjectName)
			cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "   Project directory: %s\n", projectPath)
		}
		cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "\nTo clean up manually: rm -rf %s\n", ctx.TempDir)
	}
	
	return testErr
}

// setupTestSteps configures all the test steps in the correct order
func setupTestSteps(runner *TestRunner) {
	// Environment setup steps
	runner.AddStep(NewSetupTempDirStep("./tmp"))
	runner.AddStep(NewCreateProjectStep("selftest", "github.com/petrock/selftest"))
	runner.AddStep(NewAddFeatureStep("posts"))
	
	// Build steps
	runner.AddStep(NewBuildProjectStep())
	
	// Test command generation and registration
	runner.AddStep(NewGenerateCommandStep("posts", "schedule-publication"))
	runner.AddStep(NewBuildProjectStep()) // Rebuild after generating command
	
	// Server steps
	runner.AddStep(NewStartServerStep("8081"))
	
	// HTTP testing steps
	runner.AddStep(NewHTTPGetStep("http://localhost:8081/posts", http.StatusOK))
	
	// Test command registration - verify generated command appears in command list
	runner.AddStep(NewVerifyCommandRegistrationStep("http://localhost:8081/commands", "posts/schedule-publication"))
	
	// Test invalid POST request (empty fields should fail validation)
	invalidFormData := url.Values{}
	invalidFormData.Set("name", "")
	invalidFormData.Set("description", "")
	runner.AddStep(NewHTTPPostStep("http://localhost:8081/posts/new", invalidFormData, http.StatusOK, true))
	
	// Test invalid command API request
	commandPayload := map[string]interface{}{
		"type": "posts/create",
		"payload": map[string]interface{}{
			"name":        "",
			"description": "",
		},
	}
	runner.AddStep(NewCommandAPIStep("http://localhost:8081/commands", commandPayload, http.StatusBadRequest))
	
	// Self-inspection step
	runner.AddStep(NewSelfInspectStep())
	
	// Server cleanup step
	runner.AddStep(NewStopServerStep())
	
	// MCP Server testing steps
	runner.AddStep(NewStartMCPServerStep())
	runner.AddStep(NewMCPInitializeStep())
	runner.AddStep(NewMCPListToolsStep())
	runner.AddStep(NewMCPGenerateCommandStep())
	runner.AddStep(NewStopMCPServerStep())
}

// reportFinalResults provides a summary of test execution
func reportFinalResults(runner *TestRunner, success bool) {
	results := runner.GetResults()
	
	if success {
		cmdCtx.UI.ShowSuccess(cmdCtx.Ctx, "\n✅ Test Suite Completed Successfully\n")
		cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "   %d steps executed\n", len(results))
	} else {
		cmdCtx.UI.ShowError(cmdCtx.Ctx, fmt.Errorf("❌ Test Suite Failed"))
		failedCount := 0
		for _, result := range results {
			if !result.Success {
				failedCount++
			}
		}
		cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeError, "   %d/%d steps failed\n", failedCount, len(results))
	}
	
	// Show step summary
	cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "\nStep Summary:\n")
	for _, result := range results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		msgType := ui.MessageTypeInfo
		if !result.Success {
			msgType = ui.MessageTypeError
		}
		cmdCtx.UI.Present(cmdCtx.Ctx, msgType, "  %s %s (%v)\n", status, result.StepName, result.Duration.Round(10_000_000)) // Round to 10ms
	}
}
