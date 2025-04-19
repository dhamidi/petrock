package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
	"github.com/petrock/example_module_path/petrock_example_feature_name/state" // Import state package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*RequestSummaryGenerationCommand)(nil)
var _ Validator = (*RequestSummaryGenerationCommand)(nil)

// RequestSummaryGenerationCommand requests a summary be generated for an item
type RequestSummaryGenerationCommand struct {
	ID        string `json:"id"`         // ID of the item to summarize
	RequestID string `json:"request_id"` // Unique ID for this summary request
}

// CommandName returns the unique kebab-case name for this command type
func (c *RequestSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/request-summary-generation"
}

// Validate implements the Validator interface
func (c *RequestSummaryGenerationCommand) Validate(state *state.State) error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("item ID cannot be empty")
	}
	if strings.TrimSpace(c.RequestID) == "" {
		return errors.New("request ID cannot be empty")
	}

	// Verify the item exists
	_, found := state.GetItem(c.ID)
	if !found {
		return fmt.Errorf("item with ID %q not found", c.ID)
	}
	return nil
}

// HandleRequestSummaryGeneration applies state changes for RequestSummaryGenerationCommand.
func (e *Executor) HandleRequestSummaryGeneration(ctx context.Context, command core.Command, msg *core.Message) error {
	// Type assertion for pointer type
	cmd, ok := command.(*RequestSummaryGenerationCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleRequestSummaryGeneration, expected *RequestSummaryGenerationCommand", command)
		slog.Error("Type assertion failed in HandleRequestSummaryGeneration", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for RequestSummaryGenerationCommand", "feature", "petrock_example_feature_name", "id", cmd.ID, "requestID", cmd.RequestID)
	// Nothing to do here - worker will handle this asynchronously
	// The command is just a trigger for the worker
	return nil
}
