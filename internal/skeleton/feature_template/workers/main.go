package workers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// WorkerState holds worker-specific state
type WorkerState struct {
	lastProcessedID  uint64                    // ID of last processed message
	pendingSummaries map[string]PendingSummary // keyed by RequestID
}

// Worker implements background processing for the feature
type Worker struct {
	app      *core.App
	executor *core.Executor
	state    *State
	log      *core.MessageLog

	// Worker's internal state
	wState *WorkerState

	// Configuration for external service
	apiURL string
	apiKey string
	client *http.Client
}

// NewWorker creates a new worker instance with its dependencies
func NewWorker(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) *Worker {
	return &Worker{
		app:      app,
		executor: executor,
		state:    state,
		log:      log,
		wState: &WorkerState{
			lastProcessedID:  0,
			pendingSummaries: make(map[string]PendingSummary),
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// These values should be configured from environment variables in a real application
		apiURL: "https://api.example.com/summarize",
		apiKey: "YOUR_API_KEY",
	}
}

// Start initializes the worker
// Message replay is now handled directly in App.StartWorkers
func (w *Worker) Start(ctx context.Context) error {
	slog.Info("Starting worker", "feature", "petrock_example_feature_name")

	// Initialize worker state
	w.wState.lastProcessedID = 0 // Start from the beginning of the message log
	w.wState.pendingSummaries = make(map[string]PendingSummary)

	slog.Info("Worker initialization complete", "feature", "petrock_example_feature_name")
	return nil
}

// Stop gracefully shuts down the worker
func (w *Worker) Stop(ctx context.Context) error {
	slog.Info("Stopping worker", "feature", "petrock_example_feature_name")

	// Clean up any resources if needed
	// For example, close any open connections to external services

	return nil
}

// WorkerInfo provides self-description information for introspection
func (w *Worker) WorkerInfo() *core.WorkerInfo {
	return &core.WorkerInfo{
		Name:        "petrock_example_feature_name Worker",
		Description: "Handles background processing for the petrock_example_feature_name feature, including content summarization",
	}
}

// processMessage updates worker state based on message type
func (w *Worker) processMessage(ctx context.Context, msg core.PersistedMessage) {
	// Skip if not a command
	cmd, ok := msg.DecodedPayload.(core.Command)
	if !ok {
		slog.Debug("Skipping non-command message",
			"feature", "petrock_example_feature_name",
			"messageID", msg.ID,
			"type", fmt.Sprintf("%T", msg.DecodedPayload))
		return
	}

	slog.Debug("Processing command",
		"feature", "petrock_example_feature_name",
		"messageID", msg.ID,
		"commandName", cmd.CommandName())

	switch cmd.CommandName() {
	case "petrock_example_feature_name/create":
		// When a new item is created, request summarization if needed
		w.handleCreateCommand(ctx, cmd)

	case "petrock_example_feature_name/request-summary-generation":
		// Track summary generation requests
		w.handleSummaryRequestCommand(ctx, cmd)

	case "petrock_example_feature_name/fail-summary-generation":
		// Handle failed summary generation
		w.handleSummaryFailCommand(ctx, cmd)

	case "petrock_example_feature_name/set-generated-summary":
		// Remove from pending when summary is set
		w.handleSummarySetCommand(ctx, cmd)
	}
}
