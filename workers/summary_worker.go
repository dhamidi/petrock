package workers

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/dhamidi/petrock/core" // Placeholder for target project's core package
)

// handleCreateCommand processes new item creation commands
func handleCreateCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	createCmd, ok := cmd.(*CreateCommand)
	if !ok {
		slog.Warn("Expected *CreateCommand but got different type",
			"feature", "posts",
			"type", fmt.Sprintf("%T", cmd))
		return fmt.Errorf("unexpected command type: %T", cmd)
	}

	// Skip side effects during replay
	if pctx.IsReplay {
		return nil
	}

	// Request summarization for the new item's content
	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
	summarizeCmd := &RequestSummaryGenerationCommand{
		ID:        createCmd.Name, // Using name as ID from our CreateCommand
		RequestID: requestID,
	}

	// Use a separate context with longer timeout for command execution
	execCtx, execCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer execCancel()

	if err := workerState.executor.Execute(execCtx, summarizeCmd); err != nil {
		slog.Error("Failed to request summary generation",
			"feature", "posts",
			"itemID", createCmd.Name,
			"error", err)
		return err
	}

	return nil
}

// handleSummaryRequestCommand tracks summary generation requests
func handleSummaryRequestCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	requestCmd, ok := cmd.(*RequestSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected *RequestSummaryGenerationCommand but got different type",
			"feature", "posts",
			"type", fmt.Sprintf("%T", cmd))
		return fmt.Errorf("unexpected command type: %T", cmd)
	}

	// Retrieve the content to summarize from state
	item, found := workerState.state.GetItem(requestCmd.ID)
	if !found {
		slog.Error("Cannot find item to summarize",
			"feature", "posts",
			"itemID", requestCmd.ID)
		return fmt.Errorf("item not found: %s", requestCmd.ID)
	}

	// ALWAYS update internal state (both replay and normal)
	workerState.pendingSummaries[requestCmd.RequestID] = PendingSummary{
		RequestID: requestCmd.RequestID,
		ItemID:    requestCmd.ID,
		Content:   item.Content,
		CreatedAt: time.Now(),
	}

	// Skip side effects during replay
	if pctx.IsReplay {
		return nil
	}

	slog.Info("Added content to pending summarization queue",
		"feature", "posts",
		"itemID", requestCmd.ID,
		"requestID", requestCmd.RequestID)

	return nil
}

// handleSummaryFailCommand removes failed summary requests from pending
func handleSummaryFailCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	failCmd, ok := cmd.(*FailSummaryGenerationCommand)
	if !ok {
		slog.Warn("Expected *FailSummaryGenerationCommand but got different type",
			"feature", "posts",
			"type", fmt.Sprintf("%T", cmd))
		return fmt.Errorf("unexpected command type: %T", cmd)
	}

	// ALWAYS update internal state (both replay and normal)
	delete(workerState.pendingSummaries, failCmd.RequestID)

	// Skip side effects during replay
	if pctx.IsReplay {
		return nil
	}

	slog.Info("Removed failed summary request from queue",
		"feature", "posts",
		"itemID", failCmd.ID,
		"requestID", failCmd.RequestID,
		"reason", failCmd.Reason)

	return nil
}

// handleSummarySetCommand removes completed summary requests from pending
func handleSummarySetCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	setCmd, ok := cmd.(*SetGeneratedSummaryCommand)
	if !ok {
		slog.Warn("Expected *SetGeneratedSummaryCommand but got different type",
			"feature", "posts",
			"type", fmt.Sprintf("%T", cmd))
		return fmt.Errorf("unexpected command type: %T", cmd)
	}

	// ALWAYS update internal state (both replay and normal)
	delete(workerState.pendingSummaries, setCmd.RequestID)

	// No side effects for this command - it's a state update only

	return nil
}

// processPendingSummaries calls external API for pending summaries
func processPendingSummaries(ctx context.Context, workerState *WorkerState) error {
	pendingCount := len(workerState.pendingSummaries)

	// Make a copy of pending summaries to process
	summariesToProcess := make([]PendingSummary, 0, pendingCount)
	for _, summary := range workerState.pendingSummaries {
		summariesToProcess = append(summariesToProcess, summary)
	}

	if pendingCount == 0 {
		return nil
	}

	slog.Debug("Processing pending summaries",
		"feature", "posts",
		"count", pendingCount)

	// Process each pending summary
	for _, summary := range summariesToProcess {
		// Skip if older than 24 hours (prevent infinite retries)
		if time.Since(summary.CreatedAt) > 24*time.Hour {
			slog.Warn("Abandoning old summary request",
				"feature", "posts",
				"itemID", summary.ItemID,
				"requestID", summary.RequestID,
				"age", time.Since(summary.CreatedAt))

			// Send a failure command
			failCmd := &FailSummaryGenerationCommand{
				ID:        summary.ItemID,
				RequestID: summary.RequestID,
				Reason:    "timeout",
			}
			// Use a fresh context with longer timeout for command execution
			execCtx, execCancel := context.WithTimeout(context.Background(), 60*time.Second)
			if err := workerState.executor.Execute(execCtx, failCmd); err != nil {
				slog.Error("Failed to record summary failure",
					"feature", "posts",
					"itemID", summary.ItemID,
					"error", err)
			}
			execCancel()
			continue
		}

		// Call the summarization API
		if err := callSummarizationAPI(ctx, workerState, summary); err != nil {
			slog.Error("Failed to call summarization API",
				"feature", "posts",
				"itemID", summary.ItemID,
				"error", err)
			continue
		}
	}

	return nil
}

// callSummarizationAPI calls the external API to generate a summary
func callSummarizationAPI(ctx context.Context, workerState *WorkerState, summary PendingSummary) error {
	// This is a mock implementation - in a real application, this would call an actual API

	// Simulate API call with a random delay between 500ms and 1.5s
	delay := time.Duration(500+rand.Intn(1000)) * time.Millisecond
	slog.Debug("Simulating external API call",
		"feature", "posts",
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
	// req, err := http.NewRequestWithContext(ctx, "POST", workerState.apiURL, bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return fmt.Errorf("failed to create request: %w", err)
	// }
	//
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+workerState.apiKey)
	//
	// resp, err := workerState.client.Do(req)
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

	// Use a fresh context with longer timeout for command execution
	execCtx, execCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer execCancel()

	if err := workerState.executor.Execute(execCtx, setCmd); err != nil {
		return fmt.Errorf("failed to update item with summary: %w", err)
	}

	slog.Info("Successfully generated summary",
		"feature", "posts",
		"itemID", summary.ItemID,
		"requestID", summary.RequestID)

	return nil
}
