package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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
	
	// Execute the steps and handle cleanup
	defer ctx.RunCleanup()
	
	// Execute all steps
	testCtx := context.Background()
	if err := runner.RunAllSteps(testCtx); err != nil {
		reportFinalResults(runner, false)
		return err
	}
	
	reportFinalResults(runner, true)
	cmdCtx.UI.ShowSuccess(cmdCtx.Ctx, "\nSuccess! The generated project builds correctly and serves content.\n")
	return nil
}

// setupTestSteps configures all the test steps in the correct order
func setupTestSteps(runner *TestRunner) {
	// Environment setup steps
	runner.AddStep(NewSetupTempDirStep("./tmp"))
	runner.AddStep(NewCreateProjectStep("selftest", "github.com/petrock/selftest"))
	runner.AddStep(NewAddFeatureStep("posts"))
	
	// Build steps
	runner.AddStep(NewBuildProjectStep())
	
	// Server steps
	runner.AddStep(NewStartServerStep("8081"))
	
	// HTTP testing steps
	runner.AddStep(NewHTTPGetStep("http://localhost:8081/posts", http.StatusOK))
	
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
