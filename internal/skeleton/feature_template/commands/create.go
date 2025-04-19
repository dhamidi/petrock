package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*CreateCommand)(nil)
var _ Validator = (*CreateCommand)(nil)

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"` // e.g., User ID
	CreatedAt   time.Time `json:"created_at"` // Timestamp when created
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create" // Removed suffix
}

// Validate implements the Validator interface for CreateCommand.
// It performs validation checks, potentially using the current state.
func (c *CreateCommand) Validate(state *State) error {
	// Trim all string fields
	trimmedName := strings.TrimSpace(c.Name)
	trimmedDescription := strings.TrimSpace(c.Description)

	// Basic stateless validation
	if trimmedName == "" {
		return errors.New("item name cannot be empty")
	}

	if trimmedDescription == "" {
		return errors.New("item description cannot be empty")
	}

	// Example stateful validation: Check if an item with the same name already exists
	// This simplistic approach just looks for items with the same name
	items, _ := state.ListItems(1, 1000, "")
	for _, item := range items {
		if item.Name == trimmedName {
			return fmt.Errorf("item with name %q already exists", trimmedName)
		}
	}

	// Add other validation rules...
	return nil
}

// HandleCreate applies state changes for CreateCommand.
func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message) error {
	// Type assertion for pointer type
	cmd, ok := command.(*CreateCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleCreate, expected *CreateCommand", command)
		slog.Error("Type assertion failed in HandleCreate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)

	// Apply the change using the state's Apply method
	if err := e.state.Apply(cmd, msg); err != nil {
		// Log the error, but return it to trigger panic in core.Executor
		slog.Error("State Apply failed for CreateCommand", "error", err, "name", cmd.Name)
		return fmt.Errorf("state.Apply failed for CreateCommand: %w", err)
	}

	slog.Debug("State change applied successfully for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)
	return nil
}
