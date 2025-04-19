package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

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

	// Apply the change using the state's Apply method
	if err := e.state.Apply(cmd, msg); err != nil {
		slog.Error("State Apply failed for SetGeneratedSummaryCommand", "error", err, "id", cmd.ID)
		return fmt.Errorf("state.Apply failed for SetGeneratedSummaryCommand: %w", err)
	}

	slog.Debug("State change applied successfully for SetGeneratedSummaryCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}
