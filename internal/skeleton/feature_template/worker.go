package petrock_example_feature_name

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
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
	mu              sync.Mutex
	lastProcessedID string
	pendingSummaries map[string]PendingSummary // keyed by RequestID
}

// Worker implements background processing for the feature
type Worker struct {
	app      *core.App
	executor *core.Executor
	state    *State
	log      *core.MessageLog
	
	// Worker's internal state
	wState   *WorkerState
	
	// Configuration for external service
	apiURL   string
	apiKey   string
	client   *http.Client
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
	
	// Process all existing messages to rebuild internal state
	messageCount := 0
	for msg := range w.log.After(ctx, 0) {
		messageCount++
		
		// Process the message to update internal state
		w.processMessage(ctx, msg)
		
		// Update the last processed ID
		w.wState.mu.Lock()
		w.wState.lastProcessedID = fmt.Sprintf("%d", msg.ID)
		w.wState.mu.Unlock()
	}
	
	slog.Info("Worker initialization complete", 
		"feature", "petrock_example_feature_name", 
		"messages_processed", messageCount,
		"pending_summaries", len(w.wState.pendingSummaries))
	return nil
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
	w.wState.mu.Lock()
	lastID := w.wState.lastProcessedID
	lastIDNum := uint64(0)
	if lastID != "" {
		// Convert lastID to uint64 - handle error in real code
		fmt.Sscanf(lastID, "%d", &lastIDNum)
	}
	w.wState.mu.Unlock()
	
	messageCount := 0
	for msg := range w.log.After(ctx, lastIDNum) {
		messageCount++
		
		// Process the message
		w.processMessage(ctx, msg)
		
		// Update the last processed ID
		w.wState.mu.Lock()
		w.wState.lastProcessedID = fmt.Sprintf("%d", msg.ID)
		w.wState.mu.Unlock()
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
	// Example implementation - would need to be adapted to actual command structure
	// createCmd, ok := cmd.(*CreateCommand)
	// if !ok {
	// 	return
	// }
	// 
	// // Request summarization for the new item's content
	// requestID := uuid.New().String()
	// summarizeCmd := &RequestSummaryGenerationCommand{
	// 	ID:        createCmd.ID,
	// 	RequestID: requestID,
	// }
	// 
	// if err := w.executor.Execute(ctx, summarizeCmd); err != nil {
	// 	slog.Error("Failed to request summary generation", 
	// 		"feature", "petrock_example_feature_name",
	// 		"itemID", createCmd.ID, 
	// 		"error", err)
	// }
}

// handleSummaryRequestCommand tracks summary generation requests
func (w *Worker) handleSummaryRequestCommand(ctx context.Context, cmd core.Command) {
	// Example implementation - would need to be adapted to actual command structure
	// requestCmd, ok := cmd.(*RequestSummaryGenerationCommand)
	// if !ok {
	// 	return
	// }
	// 
	// // Retrieve the content to summarize from state
	// item, found := w.state.GetByID(requestCmd.ID)
	// if !found {
	// 	slog.Error("Cannot find item to summarize", 
	// 		"feature", "petrock_example_feature_name",
	// 		"itemID", requestCmd.ID)
	// 	return
	// }
	// 
	// // Add to pending summaries
	// w.wState.mu.Lock()
	// w.wState.pendingSummaries[requestCmd.RequestID] = PendingSummary{
	// 	RequestID: requestCmd.RequestID,
	// 	ItemID:    requestCmd.ID,
	// 	Content:   item.Content,
	// 	CreatedAt: time.Now(),
	// }
	// w.wState.mu.Unlock()
}

// handleSummaryFailCommand removes failed summary requests from pending
func (w *Worker) handleSummaryFailCommand(ctx context.Context, cmd core.Command) {
	// Example implementation - would need to be adapted to actual command structure
	// failCmd, ok := cmd.(*FailSummaryGenerationCommand)
	// if !ok {
	// 	return
	// }
	// 
	// // Remove from pending summaries
	// w.wState.mu.Lock()
	// delete(w.wState.pendingSummaries, failCmd.RequestID)
	// w.wState.mu.Unlock()
}

// handleSummarySetCommand removes completed summary requests from pending
func (w *Worker) handleSummarySetCommand(ctx context.Context, cmd core.Command) {
	// Example implementation - would need to be adapted to actual command structure
	// setCmd, ok := cmd.(*SetGeneratedSummaryCommand)
	// if !ok {
	// 	return
	// }
	// 
	// // Remove from pending summaries
	// w.wState.mu.Lock()
	// delete(w.wState.pendingSummaries, setCmd.RequestID)
	// w.wState.mu.Unlock()
}

// processPendingSummaries calls external API for pending summaries
func (w *Worker) processPendingSummaries(ctx context.Context) error {
	w.wState.mu.Lock()
	pendingCount := len(w.wState.pendingSummaries)
	
	// Make a copy of pending summaries to process
	summariesToProcess := make([]PendingSummary, 0, pendingCount)
	for _, summary := range w.wState.pendingSummaries {
		summariesToProcess = append(summariesToProcess, summary)
	}
	w.wState.mu.Unlock()
	
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
	
	// Simulate API call with a short delay
	time.Sleep(100 * time.Millisecond)
	
	// Example API request preparation
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
	// Update the application state with the summary
	// setCmd := &SetGeneratedSummaryCommand{
	// 	ID:        summary.ItemID,
	// 	RequestID: summary.RequestID,
	// 	Summary:   fmt.Sprintf("This is a generated summary for content ID %s", summary.ItemID),
	// }
	// 
	// if err := w.executor.Execute(ctx, setCmd); err != nil {
	// 	return fmt.Errorf("failed to update item with summary: %w", err)
	// }
	
	slog.Info("Successfully generated summary", 
		"feature", "petrock_example_feature_name",
		"itemID", summary.ItemID, 
		"requestID", summary.RequestID)
	
	return nil
}