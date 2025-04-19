package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*SetGeneratedSummaryCommand)(nil)

// SetGeneratedSummaryCommand sets the generated summary for an item
type SetGeneratedSummaryCommand struct {
	ID        string `json:"id"`         // ID of the item
	RequestID string `json:"request_id"` // References the original request
	Summary   string `json:"summary"`    // The generated summary text
}

// CommandName returns the unique kebab-case name for this command type
func (c *SetGeneratedSummaryCommand) CommandName() string {
	return "petrock_example_feature_name/set-generated-summary"
}

// HandleSetGeneratedSummary applies state changes for SetGeneratedSummaryCommand.
func (e *Executor) HandleSetGeneratedSummary(ctx context.Context, command core.Command, msg *core.Message) error {
	// Type assertion for pointer type
	cmd, ok := command.(*SetGeneratedSummaryCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleSetGeneratedSummary, expected *SetGeneratedSummaryCommand", command)
		slog.Error("Type assertion failed in HandleSetGeneratedSummary", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for SetGeneratedSummaryCommand", "feature", "petrock_example_feature_name", "id", cmd.ID, "requestID", cmd.RequestID)

	// Get the existing item
	existingItem, found := e.state.GetItem(cmd.ID)
	if !found {
		err := fmt.Errorf("item with ID %s not found for summary update", cmd.ID)
		slog.Error("Set summary failed", "error", err, "id", cmd.ID)
		return err
	}

	// Set the summary
	existingItem.Summary = cmd.Summary
	existingItem.UpdatedAt = getTimestamp(msg)
	existingItem.Version++

	// Save the updated item
	if err := e.state.UpdateItem(existingItem); err != nil {
		slog.Error("Failed to update item summary in state", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to update item summary in state: %w", err)
	}

	slog.Debug("State change applied successfully for SetGeneratedSummaryCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}
