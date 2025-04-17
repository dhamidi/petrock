package petrock_example_feature_name

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// PendingSummary tracks a content item waiting for summarization
type PendingSummary struct {
	RequestID string
	ItemID    string
	Content   string
	CreatedAt time.Time
}

// WorkerState holds worker-specific state
type WorkerState struct {
	lastProcessedID  string
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
			lastProcessedID:  "",
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

// Start initializes the worker and rebuilds its internal state
// by processing existing messages from the message log
func (w *Worker) Start(ctx context.Context) error {
	slog.Info("Starting worker", "feature", "petrock_example_feature_name")

	// Create a timeout context for initialization to avoid hanging
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get current highest message ID
	version, err := w.log.Version(timeoutCtx)
	if err != nil {
		return fmt.Errorf("failed to get message log version: %w", err)
	}

	// Set initial position - we'll perform a full replay in the background
	// after startup to avoid blocking the application
	w.wState.lastProcessedID = fmt.Sprintf("%d", version)

	// Start a background task to process existing messages
	go w.replayExistingMessages(ctx)

	slog.Info("Worker initialization complete", "feature", "petrock_example_feature_name")
	return nil
}

// replayExistingMessages processes all existing messages to rebuild worker state
// This runs in the background to avoid blocking application startup
func (w *Worker) replayExistingMessages(ctx context.Context) {
	slog.Debug("Starting replay of existing messages", "feature", "petrock_example_feature_name")
	
	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Process all existing messages to rebuild internal state
	messageCount := 0
	for msg := range w.log.After(timeoutCtx, 0) {
		messageCount++
		
		// Process the message to update internal state
		w.processMessage(ctx, msg)
		
		// Update the last processed ID
		w.wState.lastProcessedID = fmt.Sprintf("%d", msg.ID)
	}
	
	slog.Info("Background replay completed", 
		"feature", "petrock_example_feature_name", 
		"messages_processed", messageCount,
		"pending_summaries", len(w.wState.pendingSummaries))
}

// Stop gracefully shuts down the worker
func (w *Worker) Stop(ctx context.Context) error {
	slog.Info("Stopping worker", "feature", "petrock_example_feature_name")

	// Clean up any resources if needed
	// For example, close any open connections to external services

	return nil
}

// Work performs a single processing cycle of the worker
func (w *Worker) Work() error {
	ctx := context.Background()
	
	// 1. Process any new messages since last run
	lastID := w.wState.lastProcessedID
	lastIDNum := uint64(0)
	if lastID != "" {
		// Convert lastID to uint64 - handle error in real code
		fmt.Sscanf(lastID, "%d", &lastIDNum)
	}
	
	// Create a timeout context for message processing
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	
	messageCount := 0
	for msg := range w.log.After(timeoutCtx, lastIDNum) {
		messageCount++
		
		// Process the message
		w.processMessage(ctx, msg)
		
		// Update the last processed ID
		w.wState.lastProcessedID = fmt.Sprintf("%d", msg.ID)
	}

	if messageCount > 0 {
		slog.Debug("Processed new messages",
			"feature", "petrock_example_feature_name",
			"count", messageCount)
	}

	// 2. Process any pending summaries
	return w.processPendingSummaries(ctx)
}

// processMessage updates worker state based on message type
func (w *Worker) processMessage(ctx context.Context, msg core.PersistedMessage) {
	// Skip if not a command
	cmd, ok := msg.DecodedPayload.(core.Command)
	if !ok {
		return
	}

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

// handleCreateCommand processes new item creation commands
func (w *Worker) handleCreateCommand(ctx context.Context, cmd core.Command) {
	createCmd, ok := cmd.(CreateCommand)
	if !ok {
		slog.Warn("Expected CreateCommand but got different type",
			"feature", "petrock_example_feature_name",
			"type", fmt.Sprintf("%T", cmd))
		return
	}

	// Request summarization for the new item's content
	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
	summarizeCmd := &RequestSummaryGenerationCommand{
		ID:        createCmd.Name, // Using name as ID from our CreateCommand
		RequestID: requestID,
	}

	if err := w.executor.Execute(ctx, summarizeCmd); err != nil {
		slog.Error("Failed to request summary generation",
			"feature", "petrock_example_feature_name",
			"itemID", createCmd.Name,
			"error", err)
	}
}

// handleSummaryRequestCommand tracks summary generation requests
func (w *Worker) handleSummaryRequestCommand(ctx context.Context, cmd core.Command) {
	requestCmd, ok := cmd.(RequestSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected RequestSummaryGenerationCommand but got different type",
			"feature", "petrock_example_feature_name",
			"type", fmt.Sprintf("%T", cmd))
		return
	}

	// Retrieve the content to summarize from state
	item, found := w.state.GetItem(requestCmd.ID)
	if !found {
		slog.Error("Cannot find item to summarize",
			"feature", "petrock_example_feature_name",
			"itemID", requestCmd.ID)
		return
	}

	// Add to pending summaries
	w.wState.pendingSummaries[requestCmd.RequestID] = PendingSummary{
		RequestID: requestCmd.RequestID,
		ItemID:    requestCmd.ID,
		Content:   item.Content,
		CreatedAt: time.Now(),
	}

	slog.Info("Added content to pending summarization queue",
		"feature", "petrock_example_feature_name",
		"itemID", requestCmd.ID,
		"requestID", requestCmd.RequestID)
}

// handleSummaryFailCommand removes failed summary requests from pending
func (w *Worker) handleSummaryFailCommand(ctx context.Context, cmd core.Command) {
	failCmd, ok := cmd.(FailSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected FailSummaryGenerationCommand but got different type",
			"feature", "petrock_example_feature_name",
			"type", fmt.Sprintf("%T", cmd))
		return
	}

	// Remove from pending summaries
	delete(w.wState.pendingSummaries, failCmd.RequestID)

	slog.Info("Removed failed summary request from queue",
		"feature", "petrock_example_feature_name",
		"itemID", failCmd.ID,
		"requestID", failCmd.RequestID,
		"reason", failCmd.Reason)
}

// handleSummarySetCommand removes completed summary requests from pending
func (w *Worker) handleSummarySetCommand(ctx context.Context, cmd core.Command) {
	setCmd, ok := cmd.(SetGeneratedSummaryCommand)
	if !ok {
		slog.Warn("Expected SetGeneratedSummaryCommand but got different type",
			"feature", "petrock_example_feature_name",
			"type", fmt.Sprintf("%T", cmd))
		return
	}

	// Remove from pending summaries
	delete(w.wState.pendingSummaries, setCmd.RequestID)

	slog.Info("Summary successfully set for item",
		"feature", "petrock_example_feature_name",
		"itemID", setCmd.ID,
		"requestID", setCmd.RequestID)
}

// processPendingSummaries calls external API for pending summaries
func (w *Worker) processPendingSummaries(ctx context.Context) error {
	pendingCount := len(w.wState.pendingSummaries)

	// Make a copy of pending summaries to process
	summariesToProcess := make([]PendingSummary, 0, pendingCount)
	for _, summary := range w.wState.pendingSummaries {
		summariesToProcess = append(summariesToProcess, summary)
	}

	if pendingCount == 0 {
		return nil
	}

	slog.Debug("Processing pending summaries",
		"feature", "petrock_example_feature_name",
		"count", pendingCount)

	// Process each pending summary
	for _, summary := range summariesToProcess {
		// Skip if older than 24 hours (prevent infinite retries)
		if time.Since(summary.CreatedAt) > 24*time.Hour {
			slog.Warn("Abandoning old summary request",
				"feature", "petrock_example_feature_name",
				"itemID", summary.ItemID,
				"requestID", summary.RequestID,
				"age", time.Since(summary.CreatedAt))

			// Send a failure command
			// failCmd := &FailSummaryGenerationCommand{
			// 	ID:        summary.ItemID,
			// 	RequestID: summary.RequestID,
			// 	Reason:    "timeout",
			// }
			// if err := w.executor.Execute(ctx, failCmd); err != nil {
			// 	slog.Error("Failed to record summary failure",
			// 		"feature", "petrock_example_feature_name",
			// 		"itemID", summary.ItemID,
			// 		"error", err)
			// }
			continue
		}

		// Call the summarization API
		if err := w.callSummarizationAPI(ctx, summary); err != nil {
			slog.Error("Failed to call summarization API",
				"feature", "petrock_example_feature_name",
				"itemID", summary.ItemID,
				"error", err)
			continue
		}
	}

	return nil
}

// callSummarizationAPI calls the external API to generate a summary
func (w *Worker) callSummarizationAPI(ctx context.Context, summary PendingSummary) error {
	// This is a mock implementation - in a real application, this would call an actual API

	// Simulate API call with a random delay between 500ms and 1.5s
	delay := time.Duration(500+rand.Intn(1000)) * time.Millisecond
	slog.Debug("Simulating external API call",
		"feature", "petrock_example_feature_name",
		"itemID", summary.ItemID,
		"delay", delay.String())
	time.Sleep(delay)

	// In a real implementation, you would make an actual HTTP request:
	//
	// type SummarizeRequest struct {
	// 	Content string `json:"content"`
	// }
	//
	// type SummarizeResponse struct {
	// 	Summary string `json:"summary"`
	// }
	//
	// payload := SummarizeRequest{Content: summary.Content}
	// jsonData, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal request: %w", err)
	// }
	//
	// req, err := http.NewRequestWithContext(ctx, "POST", w.apiURL, bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return fmt.Errorf("failed to create request: %w", err)
	// }
	//
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+w.apiKey)
	//
	// resp, err := w.client.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("API request failed: %w", err)
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	// }
	//
	// var result SummarizeResponse
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return fmt.Errorf("failed to decode response: %w", err)
	// }

	// For demo purposes, generate a fake summary
	fakeSummary := fmt.Sprintf("This is a concise summary of the content for item '%s'. The original text has been analyzed and condensed to capture the key points while maintaining clarity and context.", summary.ItemID)

	// Update the application state with the summary
	setCmd := &SetGeneratedSummaryCommand{
		ID:        summary.ItemID,
		RequestID: summary.RequestID,
		Summary:   fakeSummary,
	}

	if err := w.executor.Execute(ctx, setCmd); err != nil {
		return fmt.Errorf("failed to update item with summary: %w", err)
	}

	slog.Info("Successfully generated summary",
		"feature", "petrock_example_feature_name",
		"itemID", summary.ItemID,
		"requestID", summary.RequestID)

	return nil
}
