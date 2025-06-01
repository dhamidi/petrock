package workers

import (
	"context"
	"net/http"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// WorkerState holds worker-specific state
type WorkerState struct {
	pendingSummaries map[string]PendingSummary // keyed by RequestID
	state            *State                     // Reference to application state
	executor         *core.Executor             // Reference to command executor
	apiURL           string                     // Configuration for external service
	apiKey           string
	client           *http.Client
}

// NewWorker creates a new worker instance using the core worker infrastructure
func NewWorker(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) core.Worker {
	workerState := &WorkerState{
		pendingSummaries: make(map[string]PendingSummary),
		state:            state,
		executor:         executor,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// These values should be configured from environment variables in a real application
		apiURL: "https://api.example.com/summarize",
		apiKey: "YOUR_API_KEY",
	}

	worker := core.NewWorker(
		"petrock_example_feature_name Worker",
		"Handles background processing for the petrock_example_feature_name feature, including content summarization",
		workerState,
	)

	// Set dependencies
	worker.SetDependencies(log, executor)

	// Register command handlers with closures that capture worker state
	worker.OnCommand("petrock_example_feature_name/create", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
		return handleCreateCommand(ctx, cmd, msg, workerState)
	})
	worker.OnCommand("petrock_example_feature_name/request-summary-generation", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
		return handleSummaryRequestCommand(ctx, cmd, msg, workerState)
	})
	worker.OnCommand("petrock_example_feature_name/fail-summary-generation", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
		return handleSummaryFailCommand(ctx, cmd, msg, workerState)
	})
	worker.OnCommand("petrock_example_feature_name/set-generated-summary", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
		return handleSummarySetCommand(ctx, cmd, msg, workerState)
	})

	// Set periodic work
	worker.SetPeriodicWork(func(ctx context.Context) error {
		return processPendingSummaries(ctx, worker.State().(*WorkerState))
	})

	return worker
}
