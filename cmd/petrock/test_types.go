package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

// TestStep represents a single step in the test execution process
type TestStep interface {
	Name() string                           // Human-readable step name
	Execute(ctx *TestContext) *StepResult   // Execute the step
}

// StepResult contains the outcome of executing a test step
type StepResult struct {
	StepName  string        // Step identifier
	Success   bool          // Whether step passed
	Error     error         // Error if failed
	Duration  time.Duration // Execution time
	Logs      []string      // Step-specific log messages
	StartTime time.Time     // When step started
}

// TestContext holds shared state and resources during test execution
type TestContext struct {
	TempDir     string          // Base temporary directory
	ProjectDir  string          // Generated project directory
	ProjectName string          // Name of generated project
	ModulePath  string          // Go module path
	ServerCmd   *exec.Cmd       // Running server process
	ServerPort  string          // Server port number
	Cleanup     []func() error  // Cleanup functions
}

// StepFunc is a function signature for step implementations
type StepFunc func(ctx *TestContext) *StepResult

// NewStepResult creates a new step result with the given step name
func NewStepResult(stepName string) *StepResult {
	return &StepResult{
		StepName:  stepName,
		Success:   false,
		StartTime: time.Now(),
		Logs:      make([]string, 0),
	}
}

// MarkSuccess marks the step result as successful
func (sr *StepResult) MarkSuccess() *StepResult {
	sr.Success = true
	sr.Duration = time.Since(sr.StartTime)
	return sr
}

// MarkFailure marks the step result as failed with the given error
func (sr *StepResult) MarkFailure(err error) *StepResult {
	sr.Success = false
	sr.Error = err
	sr.Duration = time.Since(sr.StartTime)
	return sr
}

// AddLog adds a log message to the step result
func (sr *StepResult) AddLog(format string, args ...interface{}) {
	sr.Logs = append(sr.Logs, fmt.Sprintf(format, args...))
}

// NewTestContext creates a new test context
func NewTestContext() *TestContext {
	return &TestContext{
		Cleanup: make([]func() error, 0),
	}
}

// AddCleanup adds a cleanup function to be executed when the test finishes
func (ctx *TestContext) AddCleanup(cleanupFunc func() error) {
	ctx.Cleanup = append(ctx.Cleanup, cleanupFunc)
}

// SetTempDir sets the temporary directory for the test
func (ctx *TestContext) SetTempDir(path string) {
	ctx.TempDir = path
}

// RunCleanup executes all cleanup functions
func (ctx *TestContext) RunCleanup() {
	for i := len(ctx.Cleanup) - 1; i >= 0; i-- {
		if err := ctx.Cleanup[i](); err != nil {
			slog.Error("Cleanup function failed", "error", err)
		}
	}
}

// TestRunner executes test steps and manages their results
type TestRunner struct {
	steps   []TestStep
	results []*StepResult
	ctx     *TestContext
}

// NewTestRunner creates a new test runner with the given context
func NewTestRunner(ctx *TestContext) *TestRunner {
	return &TestRunner{
		steps:   make([]TestStep, 0),
		results: make([]*StepResult, 0),
		ctx:     ctx,
	}
}

// AddStep adds a step to the test runner
func (tr *TestRunner) AddStep(step TestStep) {
	tr.steps = append(tr.steps, step)
}

// RunStep executes a single step and returns its result
func (tr *TestRunner) RunStep(step TestStep) *StepResult {
	slog.Info("Starting test step", "step", step.Name())
	
	result := step.Execute(tr.ctx)
	tr.results = append(tr.results, result)
	
	if result.Success {
		slog.Info("Test step completed successfully", 
			"step", step.Name(), 
			"duration", result.Duration)
	} else {
		slog.Error("Test step failed", 
			"step", step.Name(), 
			"duration", result.Duration,
			"error", result.Error)
	}
	
	return result
}

// RunAllSteps executes all registered steps in order
func (tr *TestRunner) RunAllSteps(ctx context.Context) error {
	slog.Info("Starting test execution", "stepCount", len(tr.steps))
	
	for _, step := range tr.steps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			result := tr.RunStep(step)
			if !result.Success {
				return fmt.Errorf("step '%s' failed: %w", step.Name(), result.Error)
			}
		}
	}
	
	slog.Info("All test steps completed successfully")
	return nil
}

// GetResults returns all step results
func (tr *TestRunner) GetResults() []*StepResult {
	return tr.results
}
