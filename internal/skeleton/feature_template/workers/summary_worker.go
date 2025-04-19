package workers

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
)

// Work performs a single processing cycle of the worker
func (w *Worker) Work() error {
	// Use background context for overall operation
	baseCtx := context.Background()

	// Debug log worker state at the beginning of each cycle
	slog.Debug("Worker state",
		"feature", "petrock_example_feature_name",
		"lastProcessedID", w.wState.lastProcessedID,
		"pendingSummaries", len(w.wState.pendingSummaries))

	// 1. Process any new messages since last run
	lastIDNum := w.wState.lastProcessedID

	slog.Debug("Checking for new messages after ID",
		"feature", "petrock_example_feature_name",
		"afterID", lastIDNum)

	messageCount := 0
	// Create a separate short-lived context for database operations
	readCtx, readCancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer readCancel()

	for msg := range w.log.After(readCtx, lastIDNum) {
		messageCount++

		// Log message being processed
		slog.Debug("Processing message",
			"feature", "petrock_example_feature_name",
			"messageID", msg.ID,
			"commandType", fmt.Sprintf("%T", msg.DecodedPayload))

		// Process the message with a separate context
		cmdCtx, cmdCancel := context.WithTimeout(baseCtx, 5*time.Second)
		w.processMessage(cmdCtx, msg)
		cmdCancel()

		// Update the last processed ID
		w.wState.lastProcessedID = msg.ID
		slog.Debug("Updated lastProcessedID",
			"feature", "petrock_example_feature_name",
			"lastProcessedID", w.wState.lastProcessedID)
	}

	if messageCount > 0 {
		slog.Debug("Processed new messages",
			"feature", "petrock_example_feature_name",
			"count", messageCount)
	}

	// 2. Process any pending summaries with a fresh context
	summaryCtx, summaryCancel := context.WithTimeout(baseCtx, 10*time.Second)
	defer summaryCancel()
	return w.processPendingSummaries(summaryCtx)
}

// handleCreateCommand processes new item creation commands
func (w *Worker) handleCreateCommand(ctx context.Context, cmd interface{}) {
	// Type assertion for pointer type
	createCmd, ok := cmd.(*CreateCommand)
	if !ok {
		slog.Warn("Expected *CreateCommand but got different type",
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

	// Use a separate context with longer timeout for command execution
	execCtx, execCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer execCancel()

	if err := w.executor.Execute(execCtx, summarizeCmd); err != nil {
		slog.Error("Failed to request summary generation",
			"feature", "petrock_example_feature_name",
			"itemID", createCmd.Name,
			"error", err)
	}
}

// handleSummaryRequestCommand tracks summary generation requests
func (w *Worker) handleSummaryRequestCommand(ctx context.Context, cmd interface{}) {
	// Type assertion for pointer type
	requestCmd, ok := cmd.(*RequestSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected *RequestSummaryGenerationCommand but got different type",
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
func (w *Worker) handleSummaryFailCommand(ctx context.Context, cmd interface{}) {
	// Type assertion for pointer type
	failCmd, ok := cmd.(*FailSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected *FailSummaryGenerationCommand but got different type",
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
func (w *Worker) handleSummarySetCommand(ctx context.Context, cmd interface{}) {
	// Type assertion for pointer type
	setCmd, ok := cmd.(*SetGeneratedSummaryCommand)
	if !ok {
		slog.Warn("Expected *SetGeneratedSummaryCommand but got different type",
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
			failCmd := &FailSummaryGenerationCommand{
				ID:        summary.ItemID,
				RequestID: summary.RequestID,
				Reason:    "timeout",
			}
			// Use a fresh context for command execution
			execCtx, execCancel := context.WithTimeout(context.Background(), 15*time.Second)
			if err := w.executor.Execute(execCtx, failCmd); err != nil {
				slog.Error("Failed to record summary failure",
					"feature", "petrock_example_feature_name",
					"itemID", summary.ItemID,
					"error", err)
			}
			execCancel()
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

	// Use a fresh context for command execution
	execCtx, execCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer execCancel()

	if err := w.executor.Execute(execCtx, setCmd); err != nil {
		return fmt.Errorf("failed to update item with summary: %w", err)
	}

	slog.Info("Successfully generated summary",
		"feature", "petrock_example_feature_name",
		"itemID", summary.ItemID,
		"requestID", summary.RequestID)

	return nil
}
