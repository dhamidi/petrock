package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*FailSummaryGenerationCommand)(nil)

// FailSummaryGenerationCommand indicates a summary generation request failed
type FailSummaryGenerationCommand struct {
	ID        string `json:"id"`         // ID of the item
	RequestID string `json:"request_id"` // References the original request
	Reason    string `json:"reason"`     // Reason for failure
}

// CommandName returns the unique kebab-case name for this command type
func (c *FailSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/fail-summary-generation"
}

// HandleFailSummaryGeneration applies state changes for FailSummaryGenerationCommand.
func (e *Executor) HandleFailSummaryGeneration(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	cmd, ok := command.(*FailSummaryGenerationCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleFailSummaryGeneration, expected *FailSummaryGenerationCommand", command)
		slog.Error("Type assertion failed in HandleFailSummaryGeneration", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for FailSummaryGenerationCommand", "feature", "petrock_example_feature_name", "id", cmd.ID, "requestID", cmd.RequestID)
	// Nothing to do in state - this just tells the worker to stop trying
	return nil
}
